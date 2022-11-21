package dynamodbstore

import (
	"context"
	"fmt"
	"github.com/mr-joshcrane/habit"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBStore struct {
	client    *dynamodb.Client
	tablename string
}

type Record struct {
	SK string
	PK string
}

func Open(addr string, tablename string) (*DynamoDBStore) {
	client := CreateLocalClient()
	CreateTableIfNotExists(client, tablename)
	return &DynamoDBStore{
		client:    client,
		tablename: tablename,
	}
}

func (s *DynamoDBStore) GetHabit(username habit.Username, habitID habit.HabitID) (*habit.Habit, error) {
	SK := fmt.Sprintf("HABIT#%s", habitID)
	command := &dynamodb.GetItemInput{
		TableName: aws.String(s.tablename),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: string(username)},
			"SK": &types.AttributeValueMemberS{Value: SK},
		},
	}

	resp, err := s.client.GetItem(context.TODO(), command)
	if err != nil {
		return nil, err
	}
	if len(resp.Item) == 0 {
		return &habit.Habit{
			HabitName:     string(habitID),
			Username:      string(username),
			Streak:        1,
			LastPerformed: habit.Now(),
		}, err
	}
	var h habit.Habit
	attributevalue.Unmarshal(resp.Item["value"], &h)
	return &h, nil
}

func (s *DynamoDBStore) UpdateHabit(h *habit.Habit) error {
	SK := fmt.Sprintf("HABIT#%s", h.HabitName)
	v, err := attributevalue.Marshal(h)
	if err != nil {
		return err
	}
	itemMap := map[string]types.AttributeValue{
		"PK":    &types.AttributeValueMemberS{Value: h.Username},
		"SK":    &types.AttributeValueMemberS{Value: SK},
		"type":  &types.AttributeValueMemberS{Value: "HABIT"},
		"value": v,
	}

	command := dynamodb.PutItemInput{
		TableName: &s.tablename,
		Item:      itemMap,
	}
	s.client.PutItem(context.TODO(), &command)
	return nil
}

func (s *DynamoDBStore) GetBattle(code habit.BattleCode) (*habit.Battle, error) {
	itemMap := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: string(code)},
		"SK": &types.AttributeValueMemberS{Value: string(code)},
	}
	command := dynamodb.GetItemInput{
		TableName: &s.tablename,
		Key:       itemMap,
	}
	resp, err := s.client.GetItem(context.TODO(), &command)
	if err != nil {
		return nil, err
	}
	if resp.Item == nil {
		return nil, fmt.Errorf("battle not found")
	}
	var b habit.Battle
	attributevalue.Unmarshal(resp.Item["value"], &b)
	return &b, nil
}

func (s *DynamoDBStore) UpdateBattle(b *habit.Battle) error {
	v, err := attributevalue.Marshal(b)
	if err != nil {
		return err
	}
	itemMap := map[string]types.AttributeValue{
		"PK":    &types.AttributeValueMemberS{Value: string(b.Code)},
		"SK":    &types.AttributeValueMemberS{Value: string(b.Code)},
		"type":  &types.AttributeValueMemberS{Value: "BATTLE"},
		"value": v,
	}

	command := dynamodb.PutItemInput{
		TableName: &s.tablename,
		Item:      itemMap,
	}

	s.client.PutItem(context.TODO(), &command)
	if b.HabitOne.HabitName != "" {
		s.UpdateBattleStubs(*b, *b.HabitOne)
	} else {
		if b.HabitTwo.HabitName != "" {
			s.UpdateBattleStubs(*b, *b.HabitTwo)
	}
	
	}
	return nil
}

func (s *DynamoDBStore) UpdateBattleStubs(b habit.Battle, h habit.Habit) error {
	SK := fmt.Sprintf("BATTLESTUB#%s#%s", h.HabitName, b.Code )
	itemMap := map[string]types.AttributeValue{
		"PK":    &types.AttributeValueMemberS{Value: string(h.Username)},
		"SK":    &types.AttributeValueMemberS{Value: SK},
		"type":  &types.AttributeValueMemberS{Value: "BATTLESTUB"},
		"code": &types.AttributeValueMemberS{Value: string(b.Code)},
		"h1": &types.AttributeValueMemberS{Value: string(b.HabitOne.HabitName)},
		"h2": &types.AttributeValueMemberS{Value: string(b.HabitTwo.HabitName)},
	}

	command := dynamodb.PutItemInput{
		TableName: &s.tablename,
		Item:      itemMap,
	}

	s.client.PutItem(context.TODO(), &command)
	return nil
}

func (s *DynamoDBStore) ListBattlesByUser(username habit.Username) ([]*habit.Battle, error) {
	command := &dynamodb.QueryInput{
		TableName:              aws.String(s.tablename),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :b)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: string(username)},
			":b": &types.AttributeValueMemberS{Value: "BATTLESTUB"},
		},
	}
	resp, err := s.client.Query(context.TODO(), command)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	codes := []string{}
	for _, v := range resp.Items {		c := ""
		err := attributevalue.Unmarshal(v["code"], &c)
		if err != nil{ 
			fmt.Println(err)
			return []*habit.Battle{}, err
		}
		codes = append(codes, c)
	}
	results := []*habit.Battle{}
	for _, v := range codes {
		b, err := s.GetBattle(habit.BattleCode(v))
		if err != nil {
			fmt.Println(err)
			return []*habit.Battle{}, err
		}
		results = append(results, b)
	}
	return results, nil
}

func (s *DynamoDBStore) ListHabits(username habit.Username) ([]*habit.Habit, error) {
	results := []*habit.Habit{}
	command := &dynamodb.QueryInput{
		TableName:              aws.String(s.tablename),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :b)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: string(username)},
			":b": &types.AttributeValueMemberS{Value: "HABIT"},
		},
	}
	resp, err := s.client.Query(context.TODO(), command)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range resp.Items {
		var h habit.Habit
		err := attributevalue.Unmarshal(v["value"], &h)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, &h)
	}
	return results, nil
}
