package dynamodbstore_test

import (
	"habit"
	"habit/stores/dynamodbstore"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCanRetriveAHabitPerformedThreeTimes(t *testing.T) {
	t.Parallel()
	path := "localhost:9000"

	username := habit.Username("username")
	habitID := habit.HabitID("habit")

	store, err := dynamodbstore.Open(path, t.Name())
	if err != nil {
		t.Fatalf("Open incorrectly errored: %t", err)
	}
	habit.Now = aWhileAgo
	store.PerformHabit(username, habitID)
	habit.Now = monday
	store.PerformHabit(username, habitID)
	habit.Now = tuesday
	store.PerformHabit(username, habitID)
	habit.Now = wednesday
	store.PerformHabit(username, habitID)
	want := &habit.Habit{
		HabitName:     "habit",
		Username:      "username",
		LastPerformed: wednesday(),
		Streak:        3,
	}
	got, err := store.GetHabit(username, habitID)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCanRetrieveAllHabits(t *testing.T) {
	t.Parallel()
	path := "localhost:9000"

	username := habit.Username("username")

	store, err := dynamodbstore.Open(path, t.Name())

	if err != nil {
		t.Fatalf("Open incorrectly errored: %t", err)
	}
	store.PerformHabit(username, habit.HabitID("habit1"))
	store.PerformHabit(username, habit.HabitID("habit2"))
	store.PerformHabit(username, habit.HabitID("habit3"))

	want := []string{"habit1", "habit2", "habit3"}
	got := store.ListHabits(username)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func aWhileAgo() time.Time {
	return time.Date(2000, time.April, 23, 0, 0, 0, 0, time.UTC)
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
