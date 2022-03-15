package habit

import (
	"habit/proto/habitpb"
	"io"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
)

type Store struct {
	data map[string]*Habit
	path string
}

func (s *Store) Save() error {
	p := s.ToProto()
	data, err := proto.Marshal(p)
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

func (s Store) GetHabit(name string) (*Habit, bool) {
	habit, ok := s.data[name]
	if ok {
		return habit, true
	}
	h := &Habit{
		Streak: 1,
	}
	s.data[name] = h
	return h, false
}

func OpenJSONStore(path string) (*Store, error) {
	data := map[string]*Habit{}
	data2 := habitpb.Habits{}
	_, err := os.Stat(path)
	if err != nil {
		return &Store{
			data: map[string]*Habit{},
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
	habits := data2.GetHabits()
	for _, v := range habits {
		name := v.GetName()
		streak := int(v.GetStreak())
		lastPerformed := v.GetLastPerformed()
		data[name] = &Habit{
			Streak: streak,
			LastPerformed: time.Unix(lastPerformed,0),			
		}
	}
	return &Store{
		data: data,
		path: path,
	}, nil
}

func (s Store) ToProto() *habitpb.Habits {
	habits := []*habitpb.Habits_Habit{}
	for k, v := range s.data {
		habit := habitpb.Habits_Habit{
			Name: k,
			Streak: int32(v.Streak),
			LastPerformed: v.LastPerformed.Unix(),
		}
		habits = append(habits, &habit)
	}
	
	h := habitpb.Habits{
		Habits: habits,
	}
	return &h
}
