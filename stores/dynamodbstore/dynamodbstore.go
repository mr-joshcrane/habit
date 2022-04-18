package dynamodbstore

import (
	"context"
	"fmt"
	"habit"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStore struct {
	client *dynamodb.Client
	tablename string
}

type Record struct {
    SK     string
    PK     string
}

func Open(addr string, tablename string) (*DynamoDBStore, error) {
	client := CreateLocalClient()
	CreateTableIfNotExists(client, tablename)
	return &DynamoDBStore{
		client: client,
		tablename: tablename,
	}, nil
}

func (s *DynamoDBStore) ListHabits(username habit.Username) []string {
	results := []string{}
	command := dynamodb.ScanInput{
		TableName: &s.tablename,

	}
	resp, err := s.client.Scan(context.TODO(), &command)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range resp.Items {
		var h habit.Habit
		attributevalue.Unmarshal(v["value"], &h)
		results = append(results, h.HabitName)
	}
	return results
}

func (s *DynamoDBStore) RegisterBattle(code string, habitID habit.HabitID) (string, error) {
	return "", nil
}

func (s *DynamoDBStore) GetBattleAssociations(habitID habit.HabitID) ([]string, error) {
	return []string{}, nil
}

func (s *DynamoDBStore) PerformHabit(username habit.Username, habitID habit.HabitID) (int, error) {
	h, err := s.GetHabit(username, habitID)
	if err != nil {
		return 0, err
	}
	h.Perform()
	s.UpdateHabit(h)
	return h.Streak, nil
}

func (s *DynamoDBStore) GetHabit(username habit.Username, habitID habit.HabitID) (*habit.Habit, error) {
	itemMap := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: string(username)},
		"SK": &types.AttributeValueMemberS{Value: string(habitID)},
	}
	command := dynamodb.GetItemInput{
		TableName: &s.tablename,
		Key: itemMap,
	}
	resp, err := s.client.GetItem(context.TODO(), &command)
	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		fmt.Println("item not found")
		return &habit.Habit{
			Username: string(username),
			HabitName: string(habitID),
			Streak: 1,
		}, nil
	}
	var h habit.Habit
	attributevalue.Unmarshal(resp.Item["value"], &h)
	fmt.Println(h.Streak)
	return &h, nil
}

func (s *DynamoDBStore) UpdateHabit(h *habit.Habit) error {
	v, err := attributevalue.Marshal(h)
	if err != nil {
		return err
	}
	itemMap := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: h.Username},
		"SK": &types.AttributeValueMemberS{Value: h.HabitName},
		"type": &types.AttributeValueMemberS{Value: "HABIT"},
		"value": v,
	}

	command := dynamodb.PutItemInput{
		TableName: &s.tablename,
		Item:      itemMap,
	}
	s.client.PutItem(context.TODO(), &command)
	return nil
}
