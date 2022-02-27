package habit_test

import (
	"habit"
	"testing"
	"time"

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

func TestHabitPerformedTwiceOnSameDayIsStreakOfOne(t *testing.T) {
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
	h.Perform()
	want := 1
	got := h.Streak()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestHabitPerformedOnThreeConsecutiveDaysIsStreakOfThree(t *testing.T) {
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
	yesterday := func() time.Time {
		return time.Now().AddDate(0, 0, -1)
	}
	dayBeforeYesterday := func() time.Time {
		return time.Now().AddDate(0, 0, -2)
	}

	h.Perform(dayBeforeYesterday)
	h.Perform(yesterday)
	// duplicate performances on same day should not affect streak
	h.Perform(yesterday)
	h.Perform()

	want := 3
	got := h.Streak()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
