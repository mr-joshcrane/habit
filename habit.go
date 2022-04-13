package habit

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Store interface {
	PerformHabit(Username, HabitID) int
	ListHabits(Username) []string
	RegisterBattle(string, *Habit) (string, error)
	GetBattleAssociations(*Habit) []string
}

type Username string
type HabitID string

type Habit struct {
	HabitName     string
	Streak        int
	LastPerformed time.Time
	Username      string
}

type Battle struct {
	HabitOne Habit
	HabitTwo Habit
	Code     string
	Winner   string
}

type TimeOption func() time.Time

var Now = time.Now

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

func (b *Battle) IsActive() bool {
	if b.HabitOne.HabitName == "" || b.HabitTwo.HabitName == "" {
		return false
	}
	return true
}

func RunCLI(s Store) {
	// challenge := flag.String("c", "none", "Create or join a new challenge")
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
	streak := s.PerformHabit(Username(username), HabitID(habitID))
	// h.Perform()
	fmt.Println()
	if !ok {
		fmt.Fprintf(os.Stdout, "Well done, you started the new habit: %s!\n", habitID)
	} else {
		fmt.Fprintf(os.Stdout, "Well done, you continued working on habit: %s!\n", habitID)
		fmt.Fprintf(os.Stdout, "You've been performing this for a streak of %d day(s)!\n", streak)
	}
	hList := s.ListHabits(Username(username))
	fmt.Fprintf(os.Stdout, "All your current habits: %s!\n", hList)
	// if *challenge != "none" {
	// 	code, err := s.RegisterBattle(*challenge, h)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, err.Error())
	// 		os.Exit(1)
	// 	}
	// 	if *challenge == "" {
	// 		fmt.Fprintf(os.Stdout, "New challenge initiated, please give the user the following code: %s\n", code)
	// 	} else {
	// 		fmt.Fprintf(os.Stdout, "Joined challenge: %s\n", code)
	// 	}
	// }
	// b := s.GetBattleAssociations(h)
	// fmt.Fprintf(os.Stdout, "This habit is associated with the following battles: %s!\n", b)
}
