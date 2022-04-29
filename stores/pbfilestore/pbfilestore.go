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
			Streak:        streak,
			LastPerformed: time.Unix(lastPerformed, 0),
			HabitName:     habitName,
			Username:      username,
		}
	}
	return &PBFileStore{
		data: data,
		path: path,
	}, nil
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

func (s *PBFileStore) GetHabit(username habit.Username, HabitID habit.HabitID) (*habit.Habit, error) {
	h, ok := s.data[string(HabitID)]
	if ok {
		return h, nil
	}
	h = &habit.Habit{
		Username: string(username),
		HabitName: string(HabitID),
		Streak: 1,
	}
	return h, nil
}

func (s *PBFileStore) UpdateHabit(h *habit.Habit) error {
	s.data[h.HabitName] = h
	return s.UpdateHabits()
}

func (s *PBFileStore) ListHabits(username habit.Username) ([]*habit.Habit, error) {
	fmt.Println(s.data)
	habits := []*habit.Habit{}
	for _, v := range s.data {
		fmt.Println(v)
		habits = append(habits, v)
	}
	return habits, nil
}

// Battles not implemented with local file storage
func (s *PBFileStore) GetBattle(habit.BattleCode) (*habit.Battle, error) {
	return &habit.Battle{}, nil
}

// Battles not implemented with local file storage
func (s *PBFileStore) UpdateBattle(b *habit.Battle) error {
	return nil
}

// Battles not implemented with local file storage
func (s *PBFileStore) ListBattlesByUser(username habit.Username) ([]*habit.Battle, error) {
	return []*habit.Battle{}, nil
}

func (s *PBFileStore) ToProto() *habitpb.Habits {
	h := map[string]*habitpb.Habit{}

	for k, v := range s.data {
		h[k] = &habitpb.Habit{
			Streak:        int32(v.Streak),
			LastPerformed: v.LastPerformed.Unix(),
			HabitName:     k,
			User:          v.Username,
		}
	}
	return &habitpb.Habits{
		Habits: h,
	}
}
