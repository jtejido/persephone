package persephone

type StateType int

const (
	NORMAL_STATE StateType = iota
	INITIAL_STATE
)

type State uint

// Currently, we do not allow removal of a declared state
type States struct {
	states       bitSet
	initialState State
}

// The first state added always has the initial state type, otherwise, the first added state is the initialState, so please be cautious in state orders.
func (s *States) Add(state State, st StateType) {

	if st == INITIAL_STATE || s.states.len() == 0 {
		s.initialState = state
	}

	s.states.add(uint(state))
}

func (s *States) GetInitialState() State {
	return s.initialState
}

func (s *States) Contains(state State) bool {
	return s.states.contains(uint(state))
}
