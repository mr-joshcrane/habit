package habit

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"net"

	"google.golang.org/grpc"
	"github.com/phayes/freeport"
)

type Server struct {
	habitpb.UnimplementedHabitServiceServer
	port int
	Tracker
}

func NewServer(tracker Tracker) (*Server, error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	grpc := grpc.NewServer()
	s := &Server{
		Tracker: tracker,
		port: port,
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
	code, err := s.Tracker.RegisterBattle(username, habitID)
	if err != nil {
		return nil, err
	}
	return &habitpb.BattleResponse{
		Ok:      true,
		Message: "Joined successfully",
		Code:    string(code),
	}, nil
}

func (s *Server) JoinBattle(ctx context.Context, req *habitpb.BattleRequest) (*habitpb.BattleResponse, error) {
	habitID := HabitID(req.GetHabitID())
	username := Username(req.GetUsername())
	code, err := s.Tracker.RegisterBattle(username, habitID)
	if err != nil {
		return nil, err
	}
	return &habitpb.BattleResponse{
		Ok:      true,
		Message: "Registered successfully",
		Code:    string(code),
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

func (s *Server) Client() *RPCClient {
	c, err := NewRPCClient(s.port)
	if err != nil {
		panic(err)
	}
	return c
}
