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

type Habit struct {
	Name      string
	History   []HabitPerformed
	Phase     int
	Procedure []string
}

func (h *Habit) RecordHabit(time time.Time) {
	h.History = append(h.History, HabitPerformed{time})
}

func (h *Habit) UpdateProceedure(steps string) {
	proceedure := strings.Split(steps, ",")
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

func (p *Person) GetOrCreateHabit(name string, phase int) (*Habit, error) {
	for i, v := range p.Habits {
		if v.Name == name {
			return &p.Habits[i], nil
		}
	}
	habit := Habit{
		Name:      name,
		History:   []HabitPerformed{},
		Phase:     phase,
		Procedure: []string{},
	}
	p.Habits = append(p.Habits, habit)
	return &habit, nil
}

func (p *Person) Write() error {
	data, err := json.Marshal(p.Habits)
	if err != nil {
		return err
	}

	_, err = p.Datastore.Write([]byte(data))
	return err
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
	if p.Datastore != nil {

	}

	return p, nil
}

func RunCLI() {
	fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	file := fset.String("f", "store.json", "the name of the file store")
	habit := fset.String("h", "", "name of the habit")
	phase := fset.Int("p", 0, "the phase of the habit")
	err := fset.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	args := fset.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s -h HABIT -p PHASE COMMA, SEPERATED, HABIT, PROCEEDURE\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -h brushteeth -p 3 go to bathroom, put toothpaste on toothbrush, brush evenly for at least 120 seconds \n", os.Args[0])
		os.Exit(1)
	}
	proceedure := strings.Join(args, " ")

	s := OpenFilestore(*file)
	person, err := NewPerson(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	h, err := person.GetOrCreateHabit(*habit, *phase)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	h.UpdateProceedure(proceedure)
	fmt.Println(person)
	fmt.Println(person.Datastore)
	err = person.Write()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
