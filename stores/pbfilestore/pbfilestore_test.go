package pbfilestore_test

import (
	"habit"
	"habit/stores/pbfilestore"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestOpenReturnsEmptyFileStoreIfFileNotExists(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/path_does_not_exist"
	_, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatalf("Open incorrectly errored: %t", err)
	}
}

func TestOpenReturnsErrorIfInsufficientPermissions(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/insufficient_perms"
	_, err := os.Create(path)
	if err != nil {
		t.Fatal("Error creating test file")
	}
	err = os.Chmod(path, 0200)
	if err != nil {
		t.Fatal("Unable to set perms on file")
	}
	_, err = pbfilestore.Open(path)
	if err == nil {
		t.Fatal("Open did not return error")
	}
}

func TestCanRetriveAStoredHabit(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + t.Name()

	username := habit.Username("henryford")
	habitID := habit.HabitID(t.Name())
	store, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatalf("error opening store: %v", err)
	}
	want := &habit.Habit{
		HabitName:     t.Name(),
		Username:      "henryford",
		Streak:        11,
		LastPerformed: aWhileAgo(),
	}

	err = store.UpdateHabit(want)
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
	path := t.TempDir() + t.Name()

	username := habit.Username("username")

	store, err := pbfilestore.Open(path)
	if err != nil {
		t.Fatalf("error opening store: %v", err)
	}
	tracker := habit.NewTracker(store)

	tracker.PerformHabit(username, habit.HabitID("habit1"))
	tracker.PerformHabit(username, habit.HabitID("habit2"))
	tracker.PerformHabit(username, habit.HabitID("habit3"))


	want := []string{"habit1", "habit2", "habit3"}
	got := tracker.DisplayHabits(username)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Equal(want, got))
	}
}

func aWhileAgo() time.Time {
	return time.Date(2000, time.April, 23, 0, 0, 0, 0, time.UTC)
}
