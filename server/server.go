package server

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type HabitService struct {
	habitpb.UnimplementedHabitServiceServer
	database map[string]*habitpb.Habit
}

func ListenAndServe(addr string) error {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpc := grpc.NewServer()
	s := &HabitService{
		database: map[string]*habitpb.Habit{},
	}
	habitpb.RegisterHabitServiceServer(grpc, s)
	return grpc.Serve(lis)
}

func (s *HabitService) GetHabit(ctx context.Context, req *habitpb.GetHabitRequest) (*habitpb.GetHabitResponse, error) {
	h, ok := s.database[req.GetHabitname()]
	if !ok {
		return &habitpb.GetHabitResponse{
			Habit: nil,
			Ok:    false,
		}, nil
	}	

	return &habitpb.GetHabitResponse{
		Habit: h,
		Ok:    true,
	}, nil
}

func (s *HabitService) UpdateHabit(ctx context.Context, req *habitpb.UpdateHabitRequest) (*habitpb.UpdateHabitResponse, error) {
	_, ok := s.database[req.Habit.GetHabitName()]
	if !ok {
		return nil, fmt.Errorf("no such habit exists: %s", req.Habit.GetHabitName())
	}

	fmt.Println(req)
	s.database[req.Habit.GetHabitName()] = req.Habit
	fmt.Println(s.database)
	return &habitpb.UpdateHabitResponse{
		Ok:      true,
		Message: "Store updated successfully",
	}, nil
}
