package server

import (
	"context"
	"errors"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type HabitService struct {
  habitpb.UnimplementedHabitServiceServer
  database  map[string]*habitpb.Habit
}

func ListenAndServe() error {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
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
	h, ok := s.database[req.Habitname]
	if ok {
		return &habitpb.GetHabitResponse{
			Habit: h,
			Ok: true,
		}, nil
	}
	return &habitpb.GetHabitResponse{
		Habit: nil,
		Ok: false,
	}, nil
}

func (s *HabitService) UpdateHabit(ctx context.Context, h *habitpb.UpdateHabitRequest) (*habitpb.UpdateHabitResponse, error) {
	if (h.Habit.HabitName == "" || h.Habit.LastPerformed == 0 || h.Habit.Streak == 0) {
		return nil, errors.New("missing input")
	}

	fmt.Println(h)
	s.database[h.Habit.GetHabitName()] = h.Habit
	fmt.Println(s.database)
	return &habitpb.UpdateHabitResponse{
		Ok: true,
		Message: "Store updated successfully",
	}, nil
}
