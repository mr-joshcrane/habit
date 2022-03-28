package main

import (
	"context"
	"fmt"
	// "habit"
	"habit/proto/habitpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

var database = map[string]*habitpb.Habit{}

type HabitService struct {
  habitpb.UnimplementedHabitServiceServer
}

func (s *HabitService) GetHabit(ctx context.Context, req *habitpb.GetHabitRequest) (*habitpb.GetHabitResponse, error) {
	h, ok := database[req.Habitname]
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
	fmt.Println(h)
	database[h.Habit.GetHabitName()] = h.Habit
	fmt.Println(database)
	return &habitpb.UpdateHabitResponse{
		Ok: true,
		Message: "Store updated successfully",
	}, nil
}

func main() {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpc := grpc.NewServer()
	s := &HabitService{}
	habitpb.RegisterHabitServiceServer(grpc, s)
	grpc.Serve(lis)
}
