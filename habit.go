package habit

import (
	"os"
	"fmt"
	"strings"
)

type Tracker struct {
	store *JSONStore
}

type Habit struct {
	Reps int
}

func NewTracker(store *JSONStore) *Tracker {
	return &Tracker{
		store: store,
	}
}

func (t *Tracker) GetHabit(name string) (*Habit, bool) {
	return t.store.GetHabit(name)
}

func (h *Habit) Perform() {
	h.Reps++
}

func RunCLI() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintf(os.Stdout, "Pass the name of the habit you performed today\nExample: %s played violin\n", os.Args[0])
		os.Exit(0)
	}
	defaultPath := "habit.json"
	s, err := OpenJSONStore(defaultPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	t := NewTracker(s)
	habit := strings.Join(args, " ")
	h, ok := t.GetHabit(habit)
	h.Perform()
	if !ok {
		fmt.Fprintf(os.Stdout, "Well done, you started the new habit: %s!\n", habit)
	} else {
		fmt.Fprintf(os.Stdout, "Well done, you continued working on habit: %s!\n", habit)
	}
	t.store.Save()
}
