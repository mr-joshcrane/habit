package networkstore

import (
	"context"
	"errors"
	"habit"
	"habit/proto/habitpb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NetworkStore struct {
	data map[string]*habit.Habit
	client habitpb.HabitServiceClient
}


func Open(path string) (*NetworkStore, error) {
	insecure := grpc.WithTransportCredentials(insecure.NewCredentials())
	block := grpc.WithBlock()
	timeout := grpc.WithTimeout(time.Second * 3)
	conn, err := grpc.Dial("localhost:8080", insecure, block, timeout)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	defer conn.Close()
	client := habitpb.NewHabitServiceClient(conn)
	habits, err := client.GetHabit(context.TODO(), &habitpb.Empty{})
	if err != nil {
		return nil, err
	}
	data := map[string]*habit.Habit{}
	for k, v := range habits.Store {
		data[k] = &habit.Habit{
			LastPerformed: time.Unix(v.LastPerformed, 0),
			Streak: int(v.Streak),
		}
	}

	return &NetworkStore{
		data: data,
		client: client,
	}, nil
}

func (s *NetworkStore) UpdateHabit(habit *habit.Habit) error {
	data := map[string]*habitpb.Habits_Habit{}
	for k, v := range s.data {
		data[k] = &habitpb.Habits_Habit{
			Streak: int32(v.Streak),
			LastPerformed: v.LastPerformed.Unix(),
		}
	}
	store := habitpb.Habits{
		Store: data,
	}
	response, err := s.client.UpdateHabits(context.TODO(), &store)
	if err != nil {
		return err
	}
	if !response.Success {
		return errors.New(response.Message)
	}

	return nil
}

func (s NetworkStore) GetHabit(name string) (*habit.Habit, bool) {
	h, ok := s.data[name]
	if ok {
		return h, true
	}
	h = &habit.Habit{
		Streak: 1,
	}
	s.data[name] = h
	return h, false
}
