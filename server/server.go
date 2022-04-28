package server

import (
	"context"
	"fmt"
	"habit"
	"habit/proto/habitpb"
	"habit/stores/dynamodbstore"

	// "habit/stores/pbfilestore"
	"log"
	"net"

	// "github.com/google/uuid"
	"google.golang.org/grpc"
)

type HabitService struct {
	habitpb.UnimplementedHabitServiceServer
    habit.LocalTracker
}

func ListenAndServe(addr string, tablename string) error {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	store := dynamodbstore.Open(addr, tablename)

	tracker := habit.NewTracker(store)
	grpc := grpc.NewServer()
	s := &HabitService{
		LocalTracker: *tracker,
	}
	habitpb.RegisterHabitServiceServer(grpc, s)
	return grpc.Serve(lis)
}

func (s *HabitService) PerformHabit(ctx context.Context, req *habitpb.PerformHabitRequest) (*habitpb.PerformHabitResponse, error) {
	username := habit.Username(req.GetUsername())
	habitID := habit.HabitID(req.GetHabitname())
	streak, err := s.LocalTracker.PerformHabit(username, habitID)
	if err != nil {
		return nil, err
	}
	return &habitpb.PerformHabitResponse{
		Streak: int32(streak),
		Ok:     true,
	}, nil
}

func (s *HabitService) DisplayHabits(ctx context.Context, req *habitpb.ListHabitsRequest) (*habitpb.ListHabitsResponse, error) {
	username := habit.Username(req.GetUsername())
	habits, err := s.LocalTracker.Store.ListHabits(username)
	if err != nil {
		return nil, err
	}
	h := []string{}
	for _, v := range habits {
		h = append(h, v.HabitName)
	}
	resp := &habitpb.ListHabitsResponse{
		Habits: h,
	}
	return resp, nil
}

func (s *HabitService) RegisterBattle(ctx context.Context, req *habitpb.BattleRequest) (*habitpb.BattleResponse, error) {
	habitID := habit.HabitID(req.GetHabitID())
	username := habit.Username(req.GetUsername())
	code := habit.BattleCode(req.GetCode())

	code, pending, err := s.LocalTracker.RegisterBattle(code, username, habitID)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return &habitpb.BattleResponse{
		Ok:      true,
		Message: "joined successfully",
		Code:    string(code),
		Pending: bool(pending),
	}, nil
}

func (s *HabitService) GetBattleAssociations(ctx context.Context, req *habitpb.BattleAssociationsRequest) (*habitpb.BattleAssociationsResponse, error) {
	username := habit.Username(req.GetUsername())
	habitID := habit.HabitID(req.GetHabitID())
	resp := s.LocalTracker.GetBattleAssociations(username, habitID)
	var codes []string
	for _, v := range resp {
		codes = append(codes, string(v))
	}
	return &habitpb.BattleAssociationsResponse{
		Codes: codes,
	}, nil
}
