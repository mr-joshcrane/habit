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
	conn   *grpc.ClientConn
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

func (s *NetworkStore) ListHabits(username habit.Username) []string {
	req := habitpb.ListHabitsRequest{
		Username: string(username),
	}
	h, err := s.client.ListHabits(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return h.Habits
}

func (s *NetworkStore) RegisterBattle(code string, h *habit.Habit) (string, error) {
	req := habitpb.BattleRequest{
		Code: code,
		Habit: &habitpb.PerformHabitRequest{
			Habitname: h.HabitName,
			Username:  h.Username,
		},
	}
	resp, err := s.client.RegisterBattle(context.TODO(), &req)
	if err != nil {
		return "", err
	}
	return resp.Battle.GetCode(), err
}

func (s *NetworkStore) GetBattleAssociations(h *habit.Habit) []string {
	habit := &habitpb.Habit{
		HabitName:     h.HabitName,
		Streak:        int32(h.Streak),
		LastPerformed: h.LastPerformed.Unix(),
		User:          h.Username,
	}
	req := habitpb.BattleAssociationsRequest{
		Habit: habit,
	}
	associations, err := s.client.GetBattleAssociations(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
	}
	codes := []string{}
	for _, v := range associations.Battle {
		codes = append(codes, v.Code)
	}
	return codes
}

func (s *NetworkStore) PerformHabit(username habit.Username, habitID habit.HabitID) int {
	req := habitpb.PerformHabitRequest{
		Username: string(username),
		Habitname: string(habitID),
	}
	resp, err := s.client.PerformHabit(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return int(resp.GetStreak())
}
