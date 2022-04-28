package habit

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
)

type Store interface {
	GetHabit(Username, HabitID) (*Habit, error)
	UpdateHabit(*Habit) error
	ListHabits(Username) ([]*Habit, error)
	GetBattle(BattleCode) (*Battle, error)
	UpdateBattle(*Battle) error
	ListBattlesByUser(Username) ([]*Battle, error)
}

type Tracker interface {
	PerformHabit(Username, HabitID) (int, error)
	DisplayHabits(Username) []string
	RegisterBattle(BattleCode, Username, HabitID) (BattleCode, Pending, error)
	GetBattleAssociations(Username, HabitID) []BattleCode
}

type Username string
type HabitID string
type BattleID string
type BattleCode string
type Pending bool

type Habit struct {
	HabitName     string
	Streak        int
	LastPerformed time.Time
	Username      string
}

type Battle struct {
	HabitOne *Habit
	HabitTwo *Habit
	Code     BattleCode
	Winner   string
}

type TimeOption func() time.Time

var Now = time.Now
var BattleCodeGenerator = generateBattleCode

func (h *Habit) performedPreviousDay(d time.Time) bool {
	previousDay := d.AddDate(0, 0, -1)
	return h.LastPerformed.Day() == previousDay.Day()
}

func (h *Habit) Perform() {
	t := Now()
	if h.performedPreviousDay(t) {
		h.Streak++
	} else if h.LastPerformed.Before(t.AddDate(0, 0, -1)) {
		h.Streak = 1
	}
	h.LastPerformed = t
}

func CreateChallenge(h *Habit, code BattleCode) *Battle {
	if code == "" {
		return &Battle{
			HabitOne: h,
			HabitTwo: &Habit{},
			Code:     BattleCodeGenerator(),
		}
	}
	return &Battle{
		HabitOne: h,
		HabitTwo: &Habit{},
		Code:     code,
	}
}

func JoinBattle(h *Habit, b *Battle) (*Battle, error) {
	if b.HabitOne == h || b.HabitTwo == h {
		return nil, fmt.Errorf("already enrolled in this battle")
	}
	if b.HabitOne.HabitName != "" && b.HabitTwo.HabitName != "" {
		return nil, fmt.Errorf("battle already has two participants")
	}
	if b.HabitOne.Username == h.Username {
		return nil, fmt.Errorf("participant is already registered in this battle")
	}
	b.HabitTwo = h
	return b, nil
}

func (b *Battle) DetermineWinner() string {
	t := Now()
	t1 := b.HabitOne.LastPerformed
	t2 := b.HabitTwo.LastPerformed
	if t1.After(t.AddDate(0, 0, -1)) && t2.After(t.AddDate(0, 0, -1)) {
		return ""
	}
	if t1.After(t2.AddDate(0, 0, -1)) {
		return b.HabitOne.Username
	}
	return b.HabitTwo.Username
}

func (b *Battle) IsPending() bool {
	if b.HabitOne == nil || b.HabitTwo == nil {
		return true
	}
	return false
}

func generateBattleCode() BattleCode {
	length := 5
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(length)]
	}
	return BattleCode(b)
}

func RunCLI(p Tracker) {
	challenge := flag.String("c", "none", "Create or join a new challenge")
	flag.Parse()
	args := flag.Args()
	username, ok := os.LookupEnv("USER")
	if !ok {
		username = "unknown"
	}
	if len(args) == 0 {
		in := fmt.Sprintf("Pass the name of the habit you performed today\nExample: %s played violin\n", os.Args[0])
		out, _ := glamour.Render(in, "dark")
		fmt.Println(out)
		os.Exit(0)
	}
	habitID := strings.Join(args, " ")
	streak, err := p.PerformHabit(Username(username), HabitID(habitID))
	if err != nil {
		in := fmt.Sprintf("issues performing habit: %v", err)
		out, _ := glamour.Render(in, "dark")
		fmt.Fprint(os.Stderr, out)
		os.Exit(1)
	}
	if streak == 1 {
		in := fmt.Sprintf("New streak started for a new habit: %s!\n", habitID)
		out, _ := glamour.Render(in, "dark")
		fmt.Fprint(os.Stderr, out)
	} else {
		in := fmt.Sprintf("Well done, you continued working on habit: %s!\nYou've been performing this for a streak of %d day(s)!\n", habitID, streak)
		out, _ := glamour.Render(in, "dark")
		fmt.Fprint(os.Stdout, out)
	}
	hList := p.DisplayHabits(Username(username))
	in := fmt.Sprintf("All your current habits: %s!\n", hList)
	out, _ := glamour.Render(in, "dark")
	fmt.Fprint(os.Stdout, out)
	if *challenge != "none" {
		code, pending, err := p.RegisterBattle(BattleCode(*challenge), Username(username), HabitID(habitID))
		if err != nil {
			out, _ := glamour.Render(err.Error(), "dark")
			fmt.Fprint(os.Stderr, out)
			os.Exit(1)
		}
		if pending {
			in = fmt.Sprintf("New challenge initiated, please give the user the following code: %s\n", code)
			out, _ = glamour.Render(in, "dark")
			fmt.Fprint(os.Stdout, out)
		} else {
			in = fmt.Sprintf("Joined challenge: %s\n", code)
			out, _ = glamour.Render(in, "dark")
			fmt.Fprint(os.Stdout, out)
			fmt.Fprintf(os.Stdout, "Joined challenge: %s\n", code)
		}
	}
	b := p.GetBattleAssociations(Username(username), HabitID(habitID))
	in = fmt.Sprintf("Your current battles: %s!\n", b)
	out, _ = glamour.Render(in, "dark")
	fmt.Fprint(os.Stdout, out)
}
