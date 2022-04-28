package habit

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	habitpb.UnimplementedHabitServiceServer
	Tracker
}

func NewServer(tracker Tracker) (*Server, error) {
	fmt.Println("Starting server")
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	grpc := grpc.NewServer()
	s := &Server{
		Tracker: tracker,
	}
	habitpb.RegisterHabitServiceServer(grpc, s)
	go func() {
		err := grpc.Serve(lis)
		if err != nil {
			log.Print(err)
		}
	}()
	return s, nil
}

func (s *Server) Client() *gRPCClient {
	c, err := Client()
	if err != nil {
		panic(err)
	}
	return c

}

func (s *Server) PerformHabit(ctx context.Context, req *habitpb.PerformHabitRequest) (*habitpb.PerformHabitResponse, error) {
	username := Username(req.GetUsername())
	habitID := HabitID(req.GetHabitname())
	streak, err := s.Tracker.PerformHabit(username, habitID)
	if err != nil {
		return nil, err
	}
	return &habitpb.PerformHabitResponse{
		Streak: int32(streak),
		Ok:     true,
	}, nil
}

func (s *Server) DisplayHabits(ctx context.Context, req *habitpb.ListHabitsRequest) (*habitpb.ListHabitsResponse, error) {
	username := Username(req.GetUsername())
	habits := s.Tracker.DisplayHabits(username)
	h := []string{}
	h = append(h, habits...)

	resp := &habitpb.ListHabitsResponse{
		Habits: h,
	}
	return resp, nil
}

func (s *Server) RegisterBattle(ctx context.Context, req *habitpb.BattleRequest) (*habitpb.BattleResponse, error) {
	habitID := HabitID(req.GetHabitID())
	username := Username(req.GetUsername())
	code := BattleCode(req.GetCode())

	code, pending, err := s.Tracker.RegisterBattle(code, username, habitID)
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

func (s *Server) GetBattleAssociations(ctx context.Context, req *habitpb.BattleAssociationsRequest) (*habitpb.BattleAssociationsResponse, error) {
	username := Username(req.GetUsername())
	habitID := HabitID(req.GetHabitID())
	resp := s.Tracker.GetBattleAssociations(username, habitID)
	var codes []string
	for _, v := range resp {
		codes = append(codes, string(v))
	}
	return &habitpb.BattleAssociationsResponse{
		Codes: codes,
	}, nil
}

// func (t *LocalTracker) PerformHabit(username habit.Username, habitID habit.HabitID) (int, error) {
// 	h, err := t.Store.GetHabit(username, habitID)
// 	if err != nil {
// 		fmt.Println("err")
// 		return 0, err
// 	}

// 	h.Perform()
// 	t.Store.UpdateHabit(h)
// 	return h.Streak, nil
// }

// func (t *LocalTracker) DisplayHabits(username habit.Username) []string {
// 	resp, err := t.Store.ListHabits(username)
// 	if err != nil {
// 		return []string{}
// 	}
// 	results := []string{}
// 	for _, v := range resp {
// 		results = append(results, v.HabitName)
// 	}
// 	return results
// }

// func (t *LocalTracker) RegisterBattle(code habit.BattleCode, username habit.Username, habitID habit.HabitID) (habit.BattleCode, habit.Pending, error) {
// 	h, err := t.Store.GetHabit(habit.Username(username), habit.HabitID(habitID))
// 	if err != nil {
// 		return "", false, err
// 	}
// 	if code == "" {
// 		b := habit.CreateChallenge(h, code)
// 		t.Store.UpdateBattle(b)
// 		return b.Code, true, nil
// 	}
// 	b, err := t.Store.GetBattle(code)
// 	if err != nil {
// 		return "", false, err
// 	}
// 	b, err = habit.JoinBattle(h, b)
// 	if err != nil {
// 		return "", false, err
// 	}
// 	t.Store.UpdateBattle(b)
// 	return b.Code, habit.Pending(b.IsPending()), nil
// }

// func (t *LocalTracker) GetBattleAssociations(username habit.Username, habitID habit.HabitID) ([]habit.BattleID, error) {
// 	ba, err := t.Store.ListBattlesByUser(username)
// 	fmt.Println(len(ba))
// 	if err != nil {
// 		return nil, err
// 	}
// 	associations := []habit.BattleID{}
// 	for _, v := range ba {
// 		fmt.Println(v)
// 		fmt.Println(v.HabitOne)
// 		associations = append(associations, habit.BattleID(v.Code))
// 	}
// 	fmt.Println(associations)
// 	return associations, nil
// }
