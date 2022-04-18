package habit

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Store interface {
	PerformHabit(Username, HabitID) (int, error)
	ListHabits(Username) []string
	RegisterBattle(string, HabitID) (string, error)
	GetBattleAssociations(HabitID) ([]string, error)
}

type Username string
type HabitID string
type BattleCode string

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
			Code:     BattleCodeGenerator(),
		}
	}
	return &Battle{
		HabitOne: h,
		Code:     code,
	}
}

func JoinBattle(h *Habit, b *Battle) (*Battle, error) {
	if b.HabitOne == h || b.HabitTwo == h {
		return nil, fmt.Errorf("already enrolled in this battle")
	}
	if b.HabitOne != nil && b.HabitTwo != nil {
		return nil, fmt.Errorf("battle already has two participants")
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

func generateBattleCode() BattleCode {
	length := 5
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(length)]
	}
	return BattleCode(b)
}

func RunCLI(s Store) {
	challenge := flag.String("c", "none", "Create or join a new challenge")
	flag.Parse()
	args := flag.Args()
	username, ok := os.LookupEnv("USER")
	if !ok {
		username = "unknown"
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stdout, "Pass the name of the habit you performed today\nExample: %s played violin\n", os.Args[0])
		os.Exit(0)
	}
	habitID := strings.Join(args, " ")
	streak, err := s.PerformHabit(Username(username), HabitID(habitID))
	if err != nil {
		fmt.Fprintf(os.Stderr, "issues performing habit: %v", err)
		os.Exit(1)
	}
	fmt.Println()
	if streak == 1 {
		fmt.Fprintf(os.Stdout, "New streak started for a new habit: %s!\n", habitID)
	} else {
		fmt.Fprintf(os.Stdout, "Well done, you continued working on habit: %s!\n", habitID)
		fmt.Fprintf(os.Stdout, "You've been performing this for a streak of %d day(s)!\n", streak)
	}
	hList := s.ListHabits(Username(username))
	fmt.Fprintf(os.Stdout, "All your current habits: %s!\n", hList)
	if *challenge != "none" {
		code, err := s.RegisterBattle(*challenge, HabitID(habitID))
		if err != nil {
			fmt.Println("register battle")
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		if *challenge == "" {
			fmt.Fprintf(os.Stdout, "New challenge initiated, please give the user the following code: %s\n", code)
		} else {
			fmt.Fprintf(os.Stdout, "Joined challenge: %s\n", code)
		}
	}
	b, err := s.GetBattleAssociations(HabitID(habitID))
	if err != nil {
		fmt.Println("get battle association")

		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "This habit is associated with the following battles: %s!\n", b)
}
