package persephone

import (
	"fmt"
)

type StateType int

const (
	NORMAL_STATE StateType = iota
	INITIAL_STATE
)

type State struct {
	Name      string
	State     interface{}
	StateType StateType
}

func NewState(name string, state interface{}, stateType StateType) *State {
	return &State{
		Name:      name,
		State:     state,
		StateType: stateType,
	}
}

func (state *State) GetState() interface{} {
	return state.State
}

func (state *State) GetName() string {
	return state.Name
}

func (state *State) GetType() StateType {
	return state.StateType
}

// Currently, we do not allow removal of a declared state
type States struct {
	states map[string]*State
}

func NewStates() *States {
	return &States{
		states: make(map[string]*State),
	}
}

// The first state added always has the initial state type
func (s *States) Add(state *State) {

	if s.states == nil {
		s.states = make(map[string]*State)
	}

	if state.GetType() == INITIAL_STATE && s.GetInitialState() != nil {
		panic("Initial state already set.")
	}

	s.states[state.GetName()] = state
}

func (s States) GetInitialState() (state *State) {
	for _, state = range s.states {
		if state.GetType() == INITIAL_STATE {
			return
		}
	}

	return nil
}

func (s States) GetState(name string) (*State, error) {

	v, ok := s.states[name]

	if !ok {
		return nil, fmt.Errorf("No State found for name: %s", name)
	}

	return v, nil

}
