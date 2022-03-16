package pbfilestore

import (
	"habit"
	"habit/proto/habitpb"
	"io"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
)

type PBFileStore struct {
	data map[string]*habit.Habit
	path string
}

func (s *PBFileStore) UpdateHabit(*habit.Habit) error {
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

func (s PBFileStore) GetHabit(name string) (*habit.Habit, bool) {
	h, ok := s.data[name]
	if ok {
		return h, true
	}
	h = &habit.Habit{
		Streak: 1,
	}
	s.data[name] = h
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
	habits := data2.GetStore()
	for k, v := range habits {
		name := k
		streak := int(v.GetStreak())
		lastPerformed := v.GetLastPerformed()
		data[name] = &habit.Habit{
			Streak: streak,
			LastPerformed: time.Unix(lastPerformed,0),			
		}
	}
	return &PBFileStore{
		data: data,
		path: path,
	}, nil
}

func (s PBFileStore) ToProto() *habitpb.Habits {
	h := habitpb.Habits{}
	h.Store = make(map[string]*habitpb.Habits_Habit)
	
	for k, v := range s.data {
		h.Store[k] = &habitpb.Habits_Habit{
			Streak: int32(v.Streak),
			LastPerformed: v.LastPerformed.Unix(),
		}
	}
	return &h
}
