package habit_test

import (
	"habit"
	"testing"
	// "time"

	"github.com/google/go-cmp/cmp"
)

func TestPerformNewHabit(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/" + t.Name()
	s, err := habit.OpenJSONStore(path)
	if err != nil {
		t.Fatal(err)
	}
	tracker := habit.NewTracker(s)
	h, ok := tracker.GetHabit("piano")
	if ok {
		t.Fatal("habit already exists")
	}
	h.Perform()
	h2, ok := tracker.GetHabit("piano")
	if !ok {
		t.Fatal("habit should exist, but does not")
	}
	want := 1
	got := h2.Reps
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestNewHabitIsStreakOfOne(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/" + t.Name()
	s, err := habit.OpenJSONStore(path)
	if err != nil {
		t.Fatal(err)
	}
	tracker := habit.NewTracker(s)
	h, ok := tracker.GetHabit("piano")
	if ok {
		t.Fatal("habit should not, but it does")
	}
	h.Perform()
	want := 1
	got := h.Streak()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// func TestHabitPerformedOnNextDayIsStreakOfTwo(t *testing.T) {
// 	t.Parallel()
// 	path := t.TempDir() + "/" + t.Name()
// 	s, err := habit.OpenJSONStore(path)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tracker := habit.NewTracker(s)
// 	h, ok := tracker.GetHabit("piano")
// 	if ok {
// 		t.Fatal("habit should not, but it does")
// 	}
// 	yesterday := time.Now().AddDate(0, 0, -1)
// 	h.Perform(yesterday)
// 	h.Perform()
// 	want := 2
// 	got := h.Streak()
// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}
// }
