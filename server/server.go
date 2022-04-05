package server

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type UserData map[string]*habitpb.Habit
type HabitService struct {
	habitpb.UnimplementedHabitServiceServer
	store map[string]UserData
}

func ListenAndServe(addr string) error {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpc := grpc.NewServer()
	s := &HabitService{
		store: map[string]UserData{},
	}
	habitpb.RegisterHabitServiceServer(grpc, s)
	return grpc.Serve(lis)
}

func (s *HabitService) GetHabit(ctx context.Context, req *habitpb.GetHabitRequest) (*habitpb.GetHabitResponse, error) {
	store, ok := s.store[req.GetUsername()]
	if !ok {
		return &habitpb.GetHabitResponse{
			Habit: nil,
			Ok:    false,
		}, nil
	}	
	h, ok := store[req.GetHabitname()]
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
	fmt.Println(s.store)
	if req.Habit.GetHabitName() == "" {
		return nil, fmt.Errorf("habitname is required")
	}
	if req.Habit.GetUser() == "" {
		return nil, fmt.Errorf("username is required")
	}
	_, ok := s.store[req.Habit.GetUser()]
	if !ok {
		s.store[req.Habit.GetUser()] = map[string]*habitpb.Habit{}
		s.store[req.Habit.GetUser()][req.Habit.HabitName] = req.Habit
		s.store[req.Habit.GetUser()][req.Habit.HabitName].Streak = 1
		return &habitpb.UpdateHabitResponse{
			Message: "New store created. Habit UPSERTED successfully",
		}, nil
	}
	_, ok = s.store[req.Habit.GetUser()][req.Habit.GetHabitName()]
	if !ok {
		s.store[req.Habit.GetUser()][req.Habit.GetHabitName()] = req.Habit
		s.store[req.Habit.GetUser()][req.Habit.GetHabitName()].Streak = 1
		return &habitpb.UpdateHabitResponse{
			Message: "Habit UPSERTED successfully",
		}, nil
	}
	store := s.store[req.Habit.GetUser()]
	habit := store[req.Habit.GetHabitName()]

	habit.Streak ++
	habit.LastPerformed = req.Habit.GetLastPerformed()
	
	return &habitpb.UpdateHabitResponse{
		Message: "Habit UPDATED successfully",
	}, nil
}
