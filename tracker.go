package habit

type LocalTracker struct {
	Store Store
}

func NewTracker(s Store) *LocalTracker {
	return &LocalTracker{
		Store: s,
	}
}

func (t *LocalTracker) PerformHabit(username Username, habitID HabitID) (int, error) {
	h, err := t.Store.GetHabit(username, habitID)
	if err != nil {
		return 0, err
	}
	h.Perform()
	err = t.Store.UpdateHabit(h)
	if err != nil {
		return 0, err
	}
	return h.Streak, nil
}

func (t *LocalTracker) DisplayHabits(username Username) []string {
	resp, err := t.Store.ListHabits(username)
	if err != nil {
		return []string{}
	}
	results := []string{}
	for _, v := range resp {
		results = append(results, v.HabitName)
	}
	return results
}

func (t *LocalTracker) RegisterBattle(username Username, habitID HabitID) (BattleCode, error) {
	h, err := t.Store.GetHabit(Username(username), HabitID(habitID))
	if err != nil {
		return "", err
	}
	b := CreateChallenge(h)
	t.Store.UpdateBattle(b)
	return b.Code, nil
}

func (t *LocalTracker) JoinBattle(code BattleCode, username Username, habitID HabitID) error {
	h, err := t.Store.GetHabit(Username(username), HabitID(habitID))
	if err != nil {
		return err
	}
	b, err := t.Store.GetBattle(code)
	if err != nil {
		return err
	}
	b, err = JoinBattle(h, b)
	if err != nil {
		return err
	}
	t.Store.UpdateBattle(b)
	return nil
}

func (t *LocalTracker) GetBattleAssociations(username Username, habitID HabitID) []BattleCode {
	ba, err := t.Store.ListBattlesByUser(username)
	if err != nil {
		return []BattleCode{}
	}
	associations := []BattleCode{}
	for _, v := range ba {
		associations = append(associations, BattleCode(v.Code))
	}
	return associations
}
