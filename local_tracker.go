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

func (t *LocalTracker) RegisterBattle(code BattleCode, username Username, habitID HabitID) (BattleCode, Pending, error) {
	h, err := t.Store.GetHabit(Username(username), HabitID(habitID))
	if err != nil {
		return "", false, err
	}
	if code == "" {
		b := CreateChallenge(h, code)
		t.Store.UpdateBattle(b)
		return b.Code, true, nil
	}
	b, err := t.Store.GetBattle(code)
	if err != nil {
		return "", false, err
	}
	b, err = JoinBattle(h, b)
	if err != nil {
		return "", false, err
	}
	t.Store.UpdateBattle(b)
	return b.Code, Pending(b.IsPending()), nil
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
