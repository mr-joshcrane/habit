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
	h, err := s.GetHabit(username, habitID)
	if err != nil {
		t.Fatal("habit should exist, but it does")
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
	h, err := s.GetHabit(username, habitID)
	if err != nil {
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
	h, err := s.GetHabit(username, habitID)
	if err != nil  {
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
	h, err := s.GetHabit(username, habitID)
	if err != nil {
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

func TestCreateChallengeWithNoInputReturnsNewCode(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName:     "habit",
		Streak:        22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username:      "test",
	}
	habit.BattleCodeGenerator = func() habit.BattleCode { return habit.BattleCode("AAAAA") }
	want := habit.BattleCode("AAAAA")
	battle := habit.CreateChallenge(&h1, "")
	got := battle.Code
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCreateChallengeWithCodeReturnsCode(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName:     "habit",
		Streak:        22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username:      "test",
	}
	want := habit.BattleCode("ZINGO")
	battle := habit.CreateChallenge(&h1, "ZINGO")
	got := battle.Code
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestBattlesWithLessThanTwoHabitsArePending(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName:     "dance",
		Streak:        22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username:      "jeff",
	}

	b := habit.Battle{
		HabitOne: &h1,
		Code:     "EXSGYY",
		Winner:   "",
	}
	want := true
	got := b.IsPending()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestBattlesWithTwoHabitsAreNotPending(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName:     "dance",
		Streak:        22,
		LastPerformed: monday(),
		Username:      "jeff",
	}
	h2 := habit.Habit{
		HabitName:     "sing",
		Streak:        11,
		LastPerformed: tuesday(),
		Username:      "pete",
	}

	b := habit.Battle{
		HabitOne: &h1,
		HabitTwo: &h2,
		Code:     "EXSGYY",
		Winner:   "",
	}
	want := false
	got := b.IsPending()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestJoinBattleFailsIfAlreadyEnrolled(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName:     "habit",
		Streak:        22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username:      "test",
	}
	battle := habit.CreateChallenge(&h1, "ZINGO")
	_, err := habit.JoinBattle(&h1, battle)
	if err == nil {
		t.Fatal("expected to fail but did not")
	}
}

func TestJoinBattleFailsIfBattleIsFull(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName:     "habit",
		Streak:        22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username:      "test",
	}
	h2 := habit.Habit{
		HabitName:     "habit two",
		Streak:        22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username:      "test two",
	}
	h3 := habit.Habit{
		HabitName:     "habit three",
		Streak:        22,
		LastPerformed: time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC),
		Username:      "test three",
	}
	battle := habit.CreateChallenge(&h1, "ZINGO")
	_, err := habit.JoinBattle(&h2, battle)
	if err != nil {
		t.Fatalf("did not expect second habit to fail to join challenge")
	}
	_, err = habit.JoinBattle(&h3, battle)
	if err == nil {
		t.Fatal("expected to fail but did not")
	}
}

func TestDetermineWinner(t *testing.T) {
	t.Parallel()
	h1 := habit.Habit{
		HabitName:     "dance",
		Streak:        22,
		LastPerformed: mondayNight(),
		Username:      "jeff",
	}
	h2 := habit.Habit{
		HabitName:     "samba",
		Streak:        22,
		LastPerformed: mondayNight(),
		Username:      "gary",
	}

	b := habit.Battle{
		HabitOne: &h1,
		HabitTwo: &h2,
		Code:     "EXSGYY",
		Winner:   "",
	}

	habit.Now = tuesday
	winner := b.DetermineWinner()
	if winner != "" {
		t.Fatalf("should not have declared a winner as there is still time remaining to perform the habit")
	}
	habit.Now = tuesdayNight
	h1.Perform()

	habit.Now = wednesday
	winner = b.DetermineWinner()
	if winner != "jeff" {
		t.Fatalf("jeff should have won, gary's streak lapsed: winner was %s", winner)
	}
}

func monday() time.Time {
	return time.Date(2020, time.April, 23, 0, 0, 0, 0, time.UTC)
}

func mondayNight() time.Time {
	return time.Date(2020, time.April, 23, 20, 0, 0, 0, time.UTC)
}

func tuesday() time.Time {
	return time.Date(2020, time.April, 24, 0, 0, 0, 0, time.UTC)
}

func tuesdayNight() time.Time {
	return time.Date(2020, time.April, 24, 20, 0, 0, 0, time.UTC)
}

func wednesday() time.Time {
	return time.Date(2020, time.April, 25, 0, 0, 0, 0, time.UTC)
}
