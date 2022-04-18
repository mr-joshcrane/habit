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

func (s *NetworkStore) RegisterBattle(code string, habitID habit.HabitID) (string, error) {
	req := habitpb.BattleRequest{
		Code: code,
	}
	resp, err := s.client.RegisterBattle(context.TODO(), &req)
	if err != nil {
		return "", err
	}
	return resp.GetCode(), err
}

func (s *NetworkStore) GetBattleAssociations(habitID habit.HabitID) ([]string, error) {	
	req := habitpb.BattleAssociationsRequest{
		HabitID: string(habitID),
	}
	resp, err := s.client.GetBattleAssociations(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	fmt.Print(resp)
	return resp.Codes, nil
}

func (s *NetworkStore) PerformHabit(username habit.Username, habitID habit.HabitID) (int, error) {
	req := habitpb.PerformHabitRequest{
		Username: string(username),
		Habitname: string(habitID),
	}
	resp, err := s.client.PerformHabit(context.TODO(), &req)
	if err != nil {
		return 0, err
	}
	return int(resp.GetStreak()), nil
}
