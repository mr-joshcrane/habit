package networkstore

import (
	"context"
	"fmt"
	"habit"
	"habit/proto/habitpb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NetworkStore struct {
	conn *grpc.ClientConn
	client habitpb.HabitServiceClient
}

func Open(addr string) (*NetworkStore, error) {
	insecure := grpc.WithTransportCredentials(insecure.NewCredentials())
	block := grpc.WithBlock()
	timeout := grpc.WithTimeout(time.Second * 3)
	conn, err := grpc.Dial(addr, insecure, block, timeout)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	client := habitpb.NewHabitServiceClient(conn)
	if err != nil {
		return nil, err
	}

	return &NetworkStore{
		client: client,
	}, nil
}

func (s *NetworkStore) Close() {
	s.conn.Close()
}

func (s *NetworkStore) UpdateHabit(habit *habit.Habit) error {
	h := &habitpb.Habit{
		HabitName: habit.HabitName,
		Streak: int32(habit.Streak),
		LastPerformed: habit.LastPerformed.Unix(),
		User: habit.Username,
	}
	req := habitpb.UpdateHabitRequest{
		Habit: h,
	}
	_, err := s.client.UpdateHabit(context.TODO(), &req)
	if err != nil {
		return err
	}
	return nil
}

func (s NetworkStore) GetHabit(habitname, username string) (*habit.Habit, bool) {
	req := habitpb.GetHabitRequest{
		Habitname: habitname,
		Username: username,
	}
	h, err := s.client.GetHabit(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return &habit.Habit{
		Streak: int(h.Habit.GetStreak()),
		LastPerformed: time.Unix(h.Habit.GetLastPerformed(), 0),
		Username: username,
		HabitName: habitname,
	}, h.GetOk()
}
