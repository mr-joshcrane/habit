package pbfilestore

import (
	"fmt"
	"habit"
	"habit/proto/habitpb"
	"io"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
)

type PBFileStore struct {
	data map[string]*habit.Habit
	path string
}

func (s *PBFileStore) UpdateHabits() error {
	data, err := proto.Marshal(s.ToProto())
	if err != nil {
		return err
	}
	file, err := os.OpenFile(s.path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	file.Write(data)
	return err
}

func (s *PBFileStore) GetHabit(username habit.Username, HabitId habit.HabitID) (*habit.Habit, bool) {
	h, ok := s.data[string(HabitId)]
	if ok {
		return h, true
	}
	h = &habit.Habit{
		Streak: 1,
	}
	s.data[string(HabitId)] = h
	return h, false
}

func Open(path string) (*PBFileStore, error) {
	data := map[string]*habit.Habit{}
	data2 := habitpb.Habits{}
	_, err := os.Stat(path)
	if err != nil {
		return &PBFileStore{
			data: map[string]*habit.Habit{},
			path: path,
		}, nil
	}
	file, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(contents, &data2)
	if err != nil {
		return nil, err
	}
	habits := data2.Habits
	for k, v := range habits {
		name := k
		streak := int(v.GetStreak())
		lastPerformed := v.GetLastPerformed()
		habitName := v.GetHabitName()
		username := v.GetUser()
		data[name] = &habit.Habit{
			Streak: streak,
			LastPerformed: time.Unix(lastPerformed,0),
			HabitName: habitName,
			Username: username,		
		}
	}
	return &PBFileStore{
		data: data,
		path: path,
	}, nil
}

func (s *PBFileStore) ToProto() *habitpb.Habits {
	h := map[string]*habitpb.Habit{}
	
	for k, v := range s.data {
		h[k] = &habitpb.Habit{
			Streak: int32(v.Streak),
			LastPerformed: v.LastPerformed.Unix(),
			HabitName: k,
			User: v.Username,
		}
	}
	return &habitpb.Habits{
		Habits: h,
	}
}

func (s *PBFileStore) RegisterBattle(code string, h *habit.Habit) (string, error) {
	return "not implemented", nil
}

func (s *PBFileStore) GetBattleAssociations(*habit.Habit) []string {
	return nil
}

func (s *PBFileStore) ListHabits(username habit.Username) []string {
	fmt.Println("running list habits")
	fmt.Println(s.data)
	habits := []string{}
	for k, v := range s.data {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println(v.HabitName)
		habits = append(habits,k)
	}
	return habits
}

func (s *PBFileStore) PerformHabit(username habit.Username, habitID habit.HabitID) int {
	h, _ := s.GetHabit(username, habitID)
	h.Perform()
	s.UpdateHabits()
	return h.Streak
}
