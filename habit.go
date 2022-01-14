package habit

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type HabitPerformed struct {
	Date time.Time
}

type TimeOption func(*time.Time) error

type Habit struct {
	Name      string
	History   []HabitPerformed
	Phase     int
	Procedure []string
}

func (h *Habit) RecordHabit(t ...TimeOption) {
	time := time.Now()
	for _, opt := range t {
		_ = opt(&time)
	}
	h.History = append(h.History, HabitPerformed{time})
}

func (h *Habit) UpdateProceedure(steps string) {
	proceedure := strings.Split(steps, ",")
	h.Procedure = []string{}
	for _, v := range proceedure {
		v = strings.Trim(v, " ")
		h.Procedure = append(h.Procedure, v)
	}
}

type Person struct {
	Habits    []Habit
	Datastore io.ReadWriter
}

type Option func(*Person) error

func (p *Person) GetOrCreateHabit(name string, phase int) *Habit {
	for i, v := range p.Habits {
		if v.Name == name {
			return &p.Habits[i]
		}
	}
	habit := Habit{
		Name:      name,
		History:   []HabitPerformed{},
		Phase:     phase,
		Procedure: []string{},
	}
	p.Habits = append(p.Habits, habit)
	return p.GetOrCreateHabit(name, phase)
}

func (p *Person) Write() error {
	data, err := json.Marshal(p.Habits)
	if err != nil {
		return err
	}

	_, err = p.Datastore.Write([]byte(data))
	return err
}

func (p *Person) Display() string {
	habits := "---Current Habits---\n"
	for _, v := range p.Habits {
		habits += fmt.Sprintf("Habit name: %s\nPhase: %d\nTimes Performed: %d,\nProceedure: %s\n\n", v.Name, v.Phase, len(v.History), v.Procedure)
	}

	return habits
}

func OpenFilestore(path string) Option {
	return func(p *Person) error {
		err := CreateFileStore(path)
		if err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		if len(data) > 0 {
			err = json.Unmarshal(data, &p.Habits)
			if err != nil {
				return err
			}
		}
		f, err = os.Create(path)
		if err != nil {
			return err
		}
		p.Datastore = f
		return nil
	}
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

func NewPerson(opts ...Option) (Person, error) {
	p := Person{
		Habits:    []Habit{},
		Datastore: nil,
	}
	for _, opt := range opts {
		err := opt(&p)
		if err != nil {
			return Person{}, err
		}
	}
	return p, nil
}


func RunCLI() {
	fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	file := fset.String("f", "store.json", "the name of the file store")
	habit := fset.String("name", "", "name of the habit")
	phase := fset.Int("phase", 0, "the phase of the habit")
	markComplete := fset.Bool("complete", false, "indicating that a particular habit has been completed")
	err := fset.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	s := OpenFilestore(*file)
	person, err := NewPerson(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		person.Write()
		os.Exit(1)
	}
	defer person.Write()
	args := fset.Args()
	if *habit == "" {
		habits := person.Display()
		fmt.Fprintln(os.Stdout, habits)
		person.Write()
		os.Exit(0)
	}
	h := person.GetOrCreateHabit(*habit, *phase)
	if *markComplete {
		h.RecordHabit()
		fmt.Fprintf(os.Stdout, "Habit %s marked as complete!\n", h.Name)
		person.Write()
		os.Exit(0)
	}
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s -name HABIT -phase PHASE COMMA, SEPERATED, HABIT, PROCEEDURE\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -name brushteeth -phase 3 go to bathroom, put toothpaste on toothbrush, brush evenly for at least 120 seconds \n", os.Args[0])
		os.Exit(1)
	}
	proceedure := strings.Join(args, " ")
	h.UpdateProceedure(proceedure)

	habits := person.Display()
	fmt.Fprintln(os.Stdout, habits)
}
