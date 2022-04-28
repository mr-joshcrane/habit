package habit

import (
	"context"
	"fmt"
	"habit/proto/habitpb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type gRPCClient struct {
	conn   *grpc.ClientConn
	client habitpb.HabitServiceClient
}

func Client() (*gRPCClient, error) {
	address := "localhost:8080"
	insecure := grpc.WithTransportCredentials(insecure.NewCredentials())
	block := grpc.WithBlock()
	timeout := grpc.WithTimeout(time.Second * 3)
	conn, err := grpc.Dial(address, insecure, block, timeout)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	client := habitpb.NewHabitServiceClient(conn)
	if err != nil {
		return nil, err
	}
	return &gRPCClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *gRPCClient) PerformHabit(username Username, habitID HabitID) (int, error) {
	req := habitpb.PerformHabitRequest{
		Username:  string(username),
		Habitname: string(habitID),
	}
	resp, err := c.client.PerformHabit(context.TODO(), &req)
	if err != nil {
		return 0, err
	}
	return int(resp.GetStreak()), nil
}

func (c *gRPCClient) DisplayHabits(username Username) []string {
	req := habitpb.ListHabitsRequest{
		Username: string(username),
	}
	h, err := c.client.DisplayHabits(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return h.Habits
}

func (c *gRPCClient) RegisterBattle(code BattleCode, username Username, habitID HabitID) (BattleCode, Pending, error) {
	req := habitpb.BattleRequest{
		Code:     string(code),
		HabitID:  string(habitID),
		Username: string(username),
	}
	resp, err := c.client.RegisterBattle(context.TODO(), &req)
	if err != nil {
		return BattleCode(""), Pending(false), err
	}
	return BattleCode(resp.GetCode()), Pending(resp.GetPending()), err
}

func (c *gRPCClient) GetBattleAssociations(username Username, habitID HabitID) []BattleCode {
	req := habitpb.BattleAssociationsRequest{
		Username: string(username),
		HabitID:  string(habitID),
	}
	resp, err := c.client.GetBattleAssociations(context.TODO(), &req)
	if err != nil {
		return nil
	}
	var codes []BattleCode
	for _, v := range resp.Codes {
		codes = append(codes, BattleCode(v))
	}
	return codes
}
