package main

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

var database = map[string]*habitpb.Habits_Habit{}

type HabitServer struct {
  habitpb.UnimplementedHabitServiceServer
}

func (s *HabitServer) GetHabit(ctx context.Context, Empty *habitpb.Empty) (*habitpb.Habits, error) {
	return &habitpb.Habits{
		Store: database,
	}, nil
}

func (s *HabitServer) UpdateHabits(ctx context.Context, Habit *habitpb.Habits) (*habitpb.UpdateHabitsResponse, error) {
	database = Habit.Store
	fmt.Println(database)
	return &habitpb.UpdateHabitsResponse{
		Success: true,
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
	s := &HabitServer{}
	habitpb.RegisterHabitServiceServer(grpc, s)
	grpc.Serve(lis)
}
