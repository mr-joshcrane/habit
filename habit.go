package habit

import (
	// "encoding/json"
	// "flag"
	// "fmt"
	// "io"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	// "strings"
	// "time"
)

// type HabitPerformed struct {
// 	Date time.Time
// }

// type TimeOption func(*time.Time) error

// type Habit struct {
// 	Name      string
// 	History   []HabitPerformed
// 	Phase     int
// 	Procedure []string
// }

// func (h *Habit) RecordHabit(t ...TimeOption) {
// 	time := time.Now()
// 	for _, opt := range t {
// 		_ = opt(&time)
// 	}
// 	h.History = append(h.History, HabitPerformed{time})
// }

// func (h *Habit) UpdateProceedure(steps string) {
// 	proceedure := strings.Split(steps, ",")
// 	h.Procedure = []string{}
// 	for _, v := range proceedure {
// 		v = strings.Trim(v, " ")
// 		h.Procedure = append(h.Procedure, v)
// 	}
// }

// type Person struct {
// 	Habits    []Habit
// 	Datastore io.ReadWriter
// }

// type Option func(*Person) error

// func (p *Person) GetOrCreateHabit(name string, phase int) *Habit {
// 	for i, v := range p.Habits {
// 		if v.Name == name {
// 			return &p.Habits[i]
// 		}
// 	}
// 	habit := Habit{
// 		Name:      name,
// 		History:   []HabitPerformed{},
// 		Phase:     phase,
// 		Procedure: []string{},
// 	}
// 	p.Habits = append(p.Habits, habit)
// 	return p.GetOrCreateHabit(name, phase)
// }

// 	_, err = p.Datastore.Write([]byte(data))
// 	return err
// }

// func (p *Person) Display() string {
// 	habits := "---Current Habits---\n"
// 	for _, v := range p.Habits {
// 		habits += fmt.Sprintf("Habit name: %s\nPhase: %d\nTimes Performed: %d,\nProceedure: %s\n\n", v.Name, v.Phase, len(v.History), v.Procedure)
// 	}

// 	return habits
// }

// func NewPerson(opts ...Option) (Person, error) {
// 	p := Person{
// 		Habits:    []Habit{},
// 		Datastore: nil,
// 	}
// 	for _, opt := range opts {
// 		err := opt(&p)
// 		if err != nil {
// 			return Person{}, err
// 		}
// 	}
// 	return p, nil
// }

// func RunCLI() {
// 	fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
// 	file := fset.String("f", "store.json", "the name of the file store")
// 	habit := fset.String("name", "", "name of the habit")
// 	phase := fset.Int("phase", 0, "the phase of the habit")
// 	markComplete := fset.Bool("complete", false, "indicating that a particular habit has been completed")
// 	err := fset.Parse(os.Args[1:])
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// 	s := OpenFilestore(*file)
// 	person, err := NewPerson(s)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		person.Write()
// 		os.Exit(1)
// 	}
// 	defer person.Write()
// 	args := fset.Args()
// 	if *habit == "" {
// 		habits := person.Display()
// 		fmt.Fprintln(os.Stdout, habits)
// 		person.Write()
// 		os.Exit(0)
// 	}
// 	h := person.GetOrCreateHabit(*habit, *phase)
// 	if *markComplete {
// 		h.RecordHabit()
// 		fmt.Fprintf(os.Stdout, "Habit %s marked as complete!\n", h.Name)
// 		person.Write()
// 		os.Exit(0)
// 	}
// 	if len(args) < 1 {
// 		fmt.Fprintf(os.Stderr, "Usage: %s -name HABIT -phase PHASE COMMA, SEPERATED, HABIT, PROCEEDURE\n", os.Args[0])
// 		fmt.Fprintf(os.Stderr, "Example: %s -name brushteeth -phase 3 go to bathroom, put toothpaste on toothbrush, brush evenly for at least 120 seconds \n", os.Args[0])
// 		os.Exit(1)
// 	}
// 	proceedure := strings.Join(args, " ")
// 	h.UpdateProceedure(proceedure)

// 	habits := person.Display()
// 	fmt.Fprintln(os.Stdout, habits)
// }
type Store struct {
	data map[string]int
	path string
}

func (s *Store) Save() error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	file, err := os.Create(s.path)
		if err != nil {
			return err
	}
	_, err = file.Write([]byte(data))
	return err
}

func Read(path string) (map[string]int, error) {
	data := map[string]int{}
	file, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(contents) > 0 {
		err = json.Unmarshal(contents, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (s Store) HabitExists(name string) bool {
	_, ok := s.data[name]
	return ok
}

func (s *Store) PerformHabit(name string) {
	s.data[name]++
}

func (s Store) TimesPerformed(name string) int {
	return s.data[name]
}

func OpenStore(name string) (*Store, error) {
	data, err := Read(name)
	if err != nil {
		return nil, err
	}
	return &Store{
		data: data,
		path: name,
	}, err
}

func CreateFileStore(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateTempFileStore(path string) (string, error) {
	file, err := os.CreateTemp("", path)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func RunCLI() {
	args := os.Args[1:]
	defaultPath := "habit.json"
	err := CreateFileStore(defaultPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	s, err := OpenStore(defaultPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	habit := strings.Join(args, " ")
	exists := s.HabitExists(habit)
	s.PerformHabit(habit)
	err = s.Save()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	if !exists {
		fmt.Fprintf(os.Stdout, "Well done, you started the new habit: %s!\n", habit)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stdout, "Well done, you continued working on habit: %s!\n", habit)
}