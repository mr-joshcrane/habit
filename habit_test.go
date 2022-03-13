package habit_test

import (
	"habit"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewHabitPerformedHasAStreakOfOne(t *testing.T) {
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
	habit.Now = monday
	h.Perform()
	want := 1
	got := h.Streak
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPerformingAHabitTwiceOnTheSameDayDoesNotIncreaseStreak(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/" + t.Name()
	s, err := habit.OpenJSONStore(path)
	if err != nil {
		t.Fatal(err)
	}
	tracker := habit.NewTracker(s)
	h, ok := tracker.GetHabit("piano")
	if ok {
		t.Fatal("habit should not exist, but it does")
	}
	habit.Now = monday
	h.Perform()
	h.Perform()
	want := 1
	got := h.Streak
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
		t.Fatal("habit should not exist, but it does")
	}

	habit.Now = monday
	h.Perform()
	habit.Now = tuesday
	h.Perform()
	// duplicate performances on same day should not affect streak
	h.Perform()
	habit.Now = wednesday
	h.Perform()

	want := 3
	got := h.Streak
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMissingADayResetsStreak(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/" + t.Name()
	s, err := habit.OpenJSONStore(path)
	if err != nil {
		t.Fatal(err)
	}
	tracker := habit.NewTracker(s)
	h, ok := tracker.GetHabit("piano")
	if ok {
		t.Fatal("habit should not exist, but it does")
	}
	h.Streak = 100
	h.LastPerformed = monday()
	habit.Now = wednesday
	h.Perform()

	want := 1
	got := h.Streak
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func monday() time.Time {
	return time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC)
}

func tuesday() time.Time {
	return time.Date(2020, time.April, 24, 0, 0, 0, 0, time.UTC)
}

func wednesday() time.Time {
	return time.Date(2020, time.April, 25, 0, 0, 0, 0, time.UTC)
}
