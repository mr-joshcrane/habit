package habit_test

import (
	"habit"
	"testing"
	"time"
	"io/ioutil"

	"github.com/google/go-cmp/cmp"
)

func TestCreateHabit(t *testing.T) {
	t.Parallel()
	got, err := habit.NewPerson()
	if err != nil {
		t.Fatal(err)
	}
	got.GetOrCreateHabit("brush teeth", 1)
	got.GetOrCreateHabit("do pushups", 1)

	want := habit.Person{
		Habits: []habit.Habit{
			{
				Name:      "brush teeth",
				History:   []habit.HabitPerformed{},
				Phase:     1,
				Procedure: []string{},
			},
			{
				Name:      "do pushups",
				History:   []habit.HabitPerformed{},
				Phase:     1,
				Procedure: []string{},
			},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestHabitPerformed(t *testing.T) {
	t.Parallel()
	got := habit.Habit{
		Name:      "brush teeth",
		History:   []habit.HabitPerformed{},
		Phase:     1,
		Procedure: []string{},
	}
	now := time.Now()
	testTime := func(x *time.Time) error {
		*x = now
		return nil
	}
	got.RecordHabit(testTime)

	want := habit.Habit{
		Name: "brush teeth",
		History: []habit.HabitPerformed{
			{
				Date: now,
			},
		},
		Phase:     1,
		Procedure: []string{},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestProceedureUpdated(t *testing.T) {
	t.Parallel()
	h := habit.Habit{
		Name:      "brush teeth",
		History:   []habit.HabitPerformed{},
		Phase:     1,
		Procedure: []string{},
	}
	h.UpdateProceedure("go to bathroom, put toothpaste on toothbrush, brush evenly for at least 120 seconds")

	want := []string{
		"go to bathroom",
		"put toothpaste on toothbrush",
		"brush evenly for at least 120 seconds",
	}
	got := h.Procedure

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetHabit(t *testing.T) {
	t.Parallel()
	want := &habit.Habit{
		Name:    "brush teeth",
		History: []habit.HabitPerformed{},
		Phase:   1,
		Procedure: []string{
			"go to bathroom",
			"put toothpaste on toothbrush",
			"brush evenly for at least 120 seconds",
		},
	}
	otherHabit := &habit.Habit{
		Name:    "go for run",
		History: []habit.HabitPerformed{},
		Phase:   1,
		Procedure: []string{
			"put on shoes",
			"walk out door",
			"run 5km",
		},
	}
	h := habit.Person{
		Habits: []habit.Habit{
			*otherHabit,
			*want,
		},
	}

	got := h.GetOrCreateHabit("brush teeth", 1)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWrite( t *testing.T) {
	t.Parallel()
	tempFile, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatalf("unable to create temporary file, %t", err)
	}
	first, err := habit.NewPerson()
	first.Datastore = tempFile
	if err != nil {
		t.Fatalf("unable to create new person, %t", err)
	}
	h := first.GetOrCreateHabit("meditate", 3)
	if err != nil {
		t.Fatalf("unable to create new habit, %t", err)
	}
	h.UpdateProceedure("sit down, meditate")
	first.Write()

	s := habit.OpenFilestore(tempFile.Name())
	if err != nil {
		t.Fatalf("unable to read temporary file, %t", err)
	}
	second, err  := habit.NewPerson(s)
	if err != nil {
		t.Fatalf("unable to create new person, %t", err)
	}

	want := first.Habits
	got := second.Habits

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}