package server

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"math/rand"
	"net"
	"time"

	// "github.com/google/uuid"
	"google.golang.org/grpc"
)

type UserData map[string]*habitpb.Habit

type HabitService struct {
	habitpb.UnimplementedHabitServiceServer
	userdata map[string]UserData
	battles  map[string]*habitpb.Battle
	battleAssociations map[*habitpb.Habit][]*habitpb.Battle
}

func ListenAndServe(addr string) error {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpc := grpc.NewServer()
	s := &HabitService{
		userdata: map[string]UserData{},
		battles:  map[string]*habitpb.Battle{},
		battleAssociations: map[*habitpb.Habit][]*habitpb.Battle{},
	}
	habitpb.RegisterHabitServiceServer(grpc, s)
	return grpc.Serve(lis)
}

func (s *HabitService) GetHabit(ctx context.Context, req *habitpb.GetHabitRequest) (*habitpb.GetHabitResponse, error) {
	store, ok := s.userdata[req.GetUsername()]
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

func (s *HabitService) ListHabits(ctx context.Context, req *habitpb.ListHabitsRequest) (*habitpb.ListHabitsResponse, error) {
	habits, ok := s.userdata[req.GetUsername()]
	if !ok {
		return &habitpb.ListHabitsResponse{
			Username: req.Username,
			Habits: &habitpb.Habits{},
		}, nil
	}
	return &habitpb.ListHabitsResponse{
		Username: req.Username,
		Habits: &habitpb.Habits{
			Habits: habits,
		},
	}, nil
}

func (s *HabitService) UpdateHabit(ctx context.Context, req *habitpb.UpdateHabitRequest) (*habitpb.UpdateHabitResponse, error) {
	fmt.Println(s.userdata)
	fmt.Println(s.battles)
	if req.Habit.GetHabitName() == "" {
		return nil, fmt.Errorf("habitname is required")
	}
	if req.Habit.GetUser() == "" {
		return nil, fmt.Errorf("username is required")
	}
	// if req.Habit.Id == "" {
	// 	req.Habit.Id = uuid.New().String()
	// }
	_, ok := s.userdata[req.Habit.GetUser()]
	if !ok {
		s.userdata[req.Habit.GetUser()] = map[string]*habitpb.Habit{}
		s.userdata[req.Habit.GetUser()][req.Habit.HabitName] = req.Habit
		s.userdata[req.Habit.GetUser()][req.Habit.HabitName].Streak = 1
		return &habitpb.UpdateHabitResponse{
			Message: "New store created. Habit UPSERTED successfully",
		}, nil
	}
	_, ok = s.userdata[req.Habit.GetUser()][req.Habit.GetHabitName()]
	if !ok {
		s.userdata[req.Habit.GetUser()][req.Habit.GetHabitName()] = req.Habit
		s.userdata[req.Habit.GetUser()][req.Habit.GetHabitName()].Streak = 1
		return &habitpb.UpdateHabitResponse{
			Message: "Habit UPSERTED successfully",
		}, nil
	}
	store := s.userdata[req.Habit.GetUser()]
	habit := store[req.Habit.GetHabitName()]

	// ???

	requestTimestamp := time.Unix(req.Habit.GetLastPerformed(), 0)
	lastRecordedTimestamp := time.Unix(habit.GetLastPerformed(), 0).Day()
	yesterday := requestTimestamp.AddDate(0, 0, -1).Day()
	if lastRecordedTimestamp == yesterday {
		habit.Streak++
	} else if requestTimestamp.Day() != lastRecordedTimestamp {
		habit.Streak = 1
	}
	habit.LastPerformed = req.Habit.GetLastPerformed()

	return &habitpb.UpdateHabitResponse{
		Message: "Habit UPDATED successfully",
	}, nil
}

func (s *HabitService) RegisterBattle(ctx context.Context, req *habitpb.BattleRequest) (*habitpb.BattleResponse, error) {
	resp, err := s.GetHabit(ctx, req.Habit)
	h := resp.Habit
	if err != nil {
		return nil, err
	}
	b, ok := s.battles[req.GetCode()]
	if !ok {
		code := generateBattleCode(6)
		s.battles[code] = &habitpb.Battle{
			HabitOne: h,
			Code:     code,
		}
		b := s.battles[code]
		s.battleAssociations[h] = append(s.battleAssociations[h], b)
		fmt.Println(s.battleAssociations)
		return &habitpb.BattleResponse{
			Battle: &habitpb.Battle{
				Code:     b.GetCode(),
				HabitOne: b.GetHabitOne(),
				HabitTwo: b.GetHabitTwo(),
				Winner:   b.GetWinner(),
			},
			Ok:      true,
			Message: "new battle created",
		}, nil
	}
	if b.HabitOne == h || b.HabitTwo == h {
		return &habitpb.BattleResponse{
			Battle:  b,
			Ok:      true,
			Message: "battle has begun!",
		}, nil
	}
	if b.HabitTwo != nil {
		return nil, fmt.Errorf("battle with this code already has two participants")
	}
	b.HabitTwo = h
	s.battleAssociations[h] = append(s.battleAssociations[h], b)
	fmt.Println(s.battleAssociations)
	return &habitpb.BattleResponse{
		Battle:  b,
		Ok:      true,
		Message: "battle has begun!",
	}, nil
}

func (s *HabitService) GetBattleAssociations(ctx context.Context, req *habitpb.BattleAssociationsRequest) (*habitpb.BattleAssociationsResponse, error) {
	hreq := &habitpb.GetHabitRequest{
		Habitname: req.GetHabit().GetHabitName(),
		Username: req.GetHabit().GetUser(),
	}
	h, err := s.GetHabit(ctx, hreq)
	if err != nil {
		return nil, err
	}
	return &habitpb.BattleAssociationsResponse{
		Habit: h.Habit,
		Battle: s.battleAssociations[h.Habit],
	}, nil
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateBattleCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
