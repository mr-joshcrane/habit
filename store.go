package habit

import (
	"encoding/json"
	"io"
	"os"
)

type JSONStore struct {
	data map[string]*Habit
	path string
}

func (s *JSONStore) Save() error {
	data, err := json.Marshal(s.data)
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

func (s JSONStore) GetHabit(name string) (*Habit, bool) {
	habit, ok := s.data[name]
	if ok {
		return habit, true
	}
	h := &Habit{
		Reps: 0,
	}
	s.data[name] = h
	return h, false
}

func OpenJSONStore(path string) (*JSONStore, error) {
	data := map[string]*Habit{}
	_, err := os.Stat(path)
	if err != nil {
		return &JSONStore{
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
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return nil, err
	}
	return &JSONStore{
		data: data,
		path: path,
	}, nil
}
