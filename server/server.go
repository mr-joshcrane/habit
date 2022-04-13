package server

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"math/rand"
	"net"
	"habit"
	"habit/stores/pbfilestore"
	// "github.com/google/uuid"
	"google.golang.org/grpc"
)

type UserData map[string]*habitpb.Habit

type HabitService struct {
	habitpb.UnimplementedHabitServiceServer
	store habit.Store
}

func ListenAndServe(addr string) error {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	store, err := pbfilestore.Open("store")
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}

	grpc := grpc.NewServer()
	s := &HabitService{
		store: store,
	}
	habitpb.RegisterHabitServiceServer(grpc, s)
	return grpc.Serve(lis)
}

func (s *HabitService) PerformHabit(ctx context.Context, req *habitpb.PerformHabitRequest) (*habitpb.PerformHabitResponse, error) {
	username := habit.Username(req.GetUsername())
	habitID := habit.HabitID(req.GetHabitname())
	streak := s.store.PerformHabit(username, habitID)
	return &habitpb.PerformHabitResponse{
		Streak: int32(streak),
		Ok: true,
	}, nil
}

func (s *HabitService) ListHabits(ctx context.Context, req *habitpb.ListHabitsRequest) (*habitpb.ListHabitsResponse, error) {
	username := habit.Username(req.GetUsername())
	habits := s.store.ListHabits(username)
	resp := &habitpb.ListHabitsResponse{
		Habits: habits,
	}
	return resp, nil
}

func (s *HabitService) RegisterBattle(ctx context.Context, req *habitpb.BattleRequest) (*habitpb.BattleResponse, error) {
	return nil, nil
	// resp, err := s.GetHabit(ctx, req.Habit)
	// h := resp.Habit
	// if err != nil {
	// 	return nil, err
	// }
	// b, ok := s.battles[req.GetCode()]
	// if !ok {
	// 	code := generateBattleCode(6)
	// 	s.battles[code] = &habitpb.Battle{
	// 		HabitOne: h,
	// 		Code:     code,
	// 	}
	// 	b := s.battles[code]
	// 	s.battleAssociations[h] = append(s.battleAssociations[h], b)
	// 	fmt.Println(s.battleAssociations)
	// 	return &habitpb.BattleResponse{
	// 		Battle: &habitpb.Battle{
	// 			Code:     b.GetCode(),
	// 			HabitOne: b.GetHabitOne(),
	// 			HabitTwo: b.GetHabitTwo(),
	// 			Winner:   b.GetWinner(),
	// 		},
	// 		Ok:      true,
	// 		Message: "new battle created",
	// 	}, nil
	// }
	// if b.HabitOne == h || b.HabitTwo == h {
	// 	return &habitpb.BattleResponse{
	// 		Battle:  b,
	// 		Ok:      true,
	// 		Message: "battle has begun!",
	// 	}, nil
	// }
	// if b.HabitTwo != nil {
	// 	return nil, fmt.Errorf("battle with this code already has two participants")
	// }
	// b.HabitTwo = h
	// s.battleAssociations[h] = append(s.battleAssociations[h], b)
	// fmt.Println(s.battleAssociations)
	// return &habitpb.BattleResponse{
	// 	Battle:  b,
	// 	Ok:      true,
	// 	Message: "battle has begun!",
	// }, nil
}

func (s *HabitService) GetBattleAssociations(ctx context.Context, req *habitpb.BattleAssociationsRequest) (*habitpb.BattleAssociationsResponse, error) {
	return nil, nil
	// hreq := &habitpb.GetHabitRequest{
	// 	Habitname: req.GetHabit().GetHabitName(),
	// 	Username: req.GetHabit().GetUser(),
	// }
	// h, err := s.GetHabit(ctx, hreq)
	// if err != nil {
	// 	return nil, err
	// }
	// return &habitpb.BattleAssociationsResponse{
	// 	Habit: h.Habit,
	// 	Battle: s.battleAssociations[h.Habit],
	// }, nil
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateBattleCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

