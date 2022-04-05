package server_test

import (
	"habit"
	"log"
	"testing"

	"habit/server"
	"habit/stores/networkstore"

	"time"
)

func TestUpdateHabitReturnsErrorForHabitWithNoName(t *testing.T) {
	t.Parallel()
	addr := "localhost:8010"
	input := &habit.Habit{
		Streak: 1,
		LastPerformed: time.Unix(1648556311, 0),
	}
	go func() {
		err := server.ListenAndServe(addr)
		if err != nil {
			log.Fatalf("unable to start server: %v", err)
		}
	}()
	s, err := networkstore.Open(addr)
	if err != nil {
		t.Fatalf("unable to create connection to server: %sv", err)
	}
	err = s.UpdateHabit(input)
	if err == nil {
		t.Fatal("nil")
	}
}

func TestGetHabitReturnsNotOKForHabitWithNoName(t *testing.T) {
	t.Parallel()
	addr := "localhost:8011"
	go func() {
		err := server.ListenAndServe(addr)
		if err != nil {
			log.Fatalf("unable to start server: %v", err)
		}
	}()
	s, err := networkstore.Open(addr)
	if err != nil {
		t.Fatalf("unable to create connection to server: %sv", err)
	}
	_, ok := s.GetHabit("", "username")
	if ok {
		t.Fatal("ok")
	}
}
func TestGetHabitReturnsNotOKForHabitWithNoUser(t *testing.T) {
	t.Parallel()
	addr := "localhost:8012"
	go func() {
		err := server.ListenAndServe(addr)
		if err != nil {
			log.Fatalf("unable to start server: %v", err)
		}
	}()
	s, err := networkstore.Open(addr)
	if err != nil {
		t.Fatalf("unable to create connection to server: %sv", err)
	}
	_, ok := s.GetHabit("habit", "")
	if ok {
		t.Fatal("ok")
	}
}
