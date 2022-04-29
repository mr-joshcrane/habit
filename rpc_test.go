package habit_test

import (
	"habit"
	"habit/stores/pbfilestore"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClientServerCanRetrieveStoredHabits(t *testing.T) {
	t.Parallel()
	username := habit.Username("user1")
	habitname := habit.HabitID("habit1")

	path := t.TempDir() + t.Name()
	store, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatalf("failed to open filestore: %v", err)
	}
	tracker := habit.NewTracker(store)
	server, err := habit.NewServer(tracker)
	if err != nil {
		t.Fatalf("failed to create server instance: %v", err)
	}
	client := server.Client()

	client.PerformHabit(username, habitname)
	want := []string{"habit1"}
	got := client.DisplayHabits(username)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestClientServerCanRegisterABattle(t *testing.T) {
	t.Parallel()
	username := habit.Username("user1")
	habitname := habit.HabitID("habit1")

	path := t.TempDir() + t.Name()
	store, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatalf("failed to create filestore: %v", err)
	}
	tracker := habit.NewTracker(store)
	server, err := habit.NewServer(tracker)
	if err != nil {
		t.Fatalf("failed to create server instance: %v", err)
	}
	client := server.Client()
	want := habit.BattleCode("EXPCT")
	habit.BattleCodeGenerator = func() habit.BattleCode { return want }
	got, err := client.RegisterBattle(username, habitname)
	if err != nil {
		t.Fatalf("failed to register battle: %v", err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
