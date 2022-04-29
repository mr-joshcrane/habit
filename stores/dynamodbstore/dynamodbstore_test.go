package dynamodbstore_test

import (
	"fmt"
	"habit"
	"habit/stores/dynamodbstore"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCanRetriveAStoredHabit(t *testing.T) {
	t.Parallel()
	path := "localhost:9000"

	username := habit.Username("henryford")
	habitID := habit.HabitID(t.Name())
	store := dynamodbstore.Open(path, t.Name())
	want := &habit.Habit{
		HabitName: t.Name(),
		Username: "henryford",
		Streak: 11,
		LastPerformed: aWhileAgo(),
	}
	
	err := store.UpdateHabit(want)
	if err != nil {
		t.Fatalf("failed to update a habit: %v", err)
	}
	got, err := store.GetHabit(username, habitID)
	if err != nil {
		t.Fatalf("failed to get habit: %v", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCanRetrieveAllHabits(t *testing.T) {
	t.Parallel()
	path := "localhost:9000"

	username := habit.Username("username")

	store := dynamodbstore.Open(path, t.Name())
	tracker := habit.NewTracker(store)

	tracker.PerformHabit(username, habit.HabitID("habit1"))
	tracker.PerformHabit(username, habit.HabitID("habit2"))
	tracker.PerformHabit(username, habit.HabitID("habit3"))

	want := []string{"habit1", "habit2", "habit3"}
	got := tracker.DisplayHabits(username)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCanRetriveAStoredBattle(t *testing.T) {
	t.Parallel()
	path := "localhost:9000"
	store := dynamodbstore.Open(path, t.Name())
	code := habit.BattleCode("BATTL")

	h := &habit.Habit{
		HabitName: t.Name(),
		Username: "greg",
		Streak: 11,
		LastPerformed: aWhileAgo(),
	}

	habit.BattleCodeGenerator = func() habit.BattleCode { return code }
	want := habit.CreateChallenge(h)

	
	err := store.UpdateBattle(want)
	if err != nil {
		t.Fatalf("failed to update battle: %v", err)
	}
	got, err := store.GetBattle(code)
	if err != nil {
		t.Fatalf("failed to get battle: %v", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCanRetrieveAllBattlesFromOneUser(t *testing.T) {
	t.Parallel()
	path := "localhost:9000"

	username1 := habit.Username("user1")
	username2 := habit.Username("user2")

	tablename := t.Name() + fmt.Sprint(time.Now().UnixMilli())

	store := dynamodbstore.Open(path, tablename)
	tracker := habit.NewTracker(store)

	tracker.PerformHabit(username1, habit.HabitID("habit1"))
	tracker.PerformHabit(username1, habit.HabitID("habit2"))
	tracker.PerformHabit(username2, habit.HabitID("habit3"))

	_, err := tracker.RegisterBattle(username1, habit.HabitID("habitthefirst"))
	if err != nil {
		t.Fatalf("unable to register battle: %v", err)
	}
	_, err = tracker.RegisterBattle(username1, habit.HabitID("habitthesecond"))
	if err != nil {
		t.Fatalf("unable to register battle: %v", err)
	}
	_, err = tracker.RegisterBattle(username2, habit.HabitID("habit3"))
	if err != nil {
		t.Fatalf("unable to register battle: %v", err)
	}

	want := 2
	got, err := store.ListBattlesByUser(username1)
	if err != nil {
		t.Fatalf("failed to list battles: %v", err)

	}
	if !cmp.Equal(want, len(got)) {
		t.Error(cmp.Diff(want, got))
	}
}


func aWhileAgo() time.Time {
	return time.Date(2000, time.April, 23, 0, 0, 0, 0, time.UTC)
}
