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
	name      string
	state     interface{}
	stateType StateType
}

func NewState(name string, state interface{}, stateType StateType) *State {
	return &State{
		name:      name,
		state:     state,
		stateType: stateType,
	}
}

func (state *State) GetState() interface{} {
	return state.state
}

func (state *State) GetName() string {
	return state.name
}

func (state *State) GetType() StateType {
	return state.stateType
}

// Currently, we do not allow removal of a declared state
type States struct {
	states []*State
}

func NewStates() *States {
	return &States{
		states: make([]*State, 0),
	}
}

// The first state added always has the initial state type
func (s *States) Add(name string, state interface{}) {

	var state_type = NORMAL_STATE

	if s.states == nil {
		s.states = make([]*State, 0)
	}

	if len(s.states) <= 0 {
		state_type = INITIAL_STATE
	}

	st := &State{
		name:      name,
		state:     state,
		stateType: state_type,
	}

	s.states = append(s.states, st)
}

func (s States) GetInitialState() (state *State) {
	for _, state = range s.states {
		if state.GetType() == INITIAL_STATE {
			return
		}
	}

	return nil
}

func (s States) GetStateByName(name string) (*State, error) {
	for _, state := range s.states {
		if state.GetName() == name {
			return state, nil
		}
	}

	return nil, fmt.Errorf("No State found for name: %s", name)
}
