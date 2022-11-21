package habit

import (
	"context"
	"fmt"
	"github.com/mr-joshcrane/habit/proto/habitpb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	
)

type RPCClient struct {
	client habitpb.HabitServiceClient
}

func NewRPCClient(p int) (*RPCClient, error) {
	address := fmt.Sprintf("localhost:%d", p)
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
	return &RPCClient{
		client: client,
	}, nil
}

func (c *RPCClient) PerformHabit(username Username, habitID HabitID) (int, error) {
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

func (c *RPCClient) DisplayHabits(username Username) []string {
	req := habitpb.ListHabitsRequest{
		Username: string(username),
	}
	h, err := c.client.DisplayHabits(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return h.Habits
}

func (c *RPCClient) RegisterBattle(username Username, habitID HabitID) (BattleCode, error) {
	req := habitpb.BattleRequest{
		HabitID:  string(habitID),
		Username: string(username),
	}
	resp, err := c.client.RegisterBattle(context.TODO(), &req)
	if err != nil {
		return BattleCode(""), err
	}
	return BattleCode(resp.GetCode()), err
}

func (c *RPCClient) JoinBattle(code BattleCode, username Username, habitID HabitID) error {
	req := habitpb.BattleRequest{
		Code:     string(code),
		HabitID:  string(habitID),
		Username: string(username),
	}
	_, err := c.client.RegisterBattle(context.TODO(), &req)
	if err != nil {
		return err
	}
	return err
}

func (c *RPCClient) GetBattleAssociations(username Username, habitID HabitID) []BattleCode {
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
