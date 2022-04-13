package habit_test

import (
	"habit"
	"habit/stores/pbfilestore"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewHabitPerformedHasAStreakOfOne(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/" + t.Name()
	username := habit.Username("test")
	habitID := habit.HabitID("piano")
	s, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	h, ok := s.GetHabit(username, habitID)
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
	username := habit.Username("test")
	habitID := habit.HabitID("piano")
	s, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	h, ok := s.GetHabit(username, habitID)
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
	username := habit.Username("test")
	habitID := habit.HabitID("piano")
	s, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	h, ok := s.GetHabit(username, habitID)
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
	username := habit.Username("test")
	habitID := habit.HabitID("piano")
	s, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	h, ok := s.GetHabit(username, habitID)
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

func TestBattlesWithLessThanTwoHabitsAreNotActive(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName: "dance",
		Streak: 22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username: "jeff",
	}
	
	b := habit.Battle{
		HabitOne: h1,
		Code: "EXSGYY",
		Winner: "",
	}
	want := false
	got := b.IsActive()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestBattlesWithTwoHabitsAreActive(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName: "dance",
		Streak: 22,
		LastPerformed: monday(),
		Username: "jeff",
	}
	h2 := habit.Habit{
		HabitName: "sing",
		Streak: 11,
		LastPerformed: tuesday(),
		Username: "pete",
	}
	
	b := habit.Battle{
		HabitOne: h1,
		HabitTwo: h2,
		Code: "EXSGYY",
		Winner: "",
	}
	want := true
	got := b.IsActive()
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

// Testing network store

// Contract for a store, for things where we dont care how, its in habit

// Implimentation tests concrete gotchya pbfilestore
