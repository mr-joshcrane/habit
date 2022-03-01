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
	dayBeforeYesterday := func() time.Time {
		return time.Now().AddDate(0, 0, -2)
	}
	h.Perform(dayBeforeYesterday)
	h.Perform()

	want := 1
	got := h.Streak
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
