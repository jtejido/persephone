package persephone

import (
	"fmt"
)

type InputActions struct {
	actionMap map[string]*SourceInputActionMap
}

func newInputActions() *InputActions {
	return &InputActions{
		actionMap: make(map[string]*SourceInputActionMap),
	}
}

func (ia *InputActions) addMap(src *State, in *Input, act Callback) {

	if ia.actionMap == nil {
		ia.actionMap = make(map[string]*SourceInputActionMap)
	}

	sia, _ := ia.getMapByState(src.GetName())

	if sia != nil {
		ai, err := sia.getActionsByInput(in.GetName())

		if err != nil {
			sia.addMap(in, act)

			return
		}

		ai.addAction(act)

		return
	}

	m := newSourceInputActionMap(src)
	m.addMap(in, act)
	ia.actionMap[src.GetName()] = m

}

func (ia *InputActions) getMapByState(name string) (*SourceInputActionMap, error) {

	v, ok := ia.actionMap[name]

	if ok {
		return v, nil
	}

	return nil, fmt.Errorf("No SourceInputAction Map found for State: %s", name)
}

type SourceInputActionMap struct {
	state     *State
	actionMap []*InputActionMap
}

func newSourceInputActionMap(state *State) *SourceInputActionMap {
	return &SourceInputActionMap{
		state:     state,
		actionMap: make([]*InputActionMap, 0),
	}
}

func (sia *SourceInputActionMap) getState() *State {
	return sia.state
}

func (sia *SourceInputActionMap) addMap(input *Input, action Callback) {
	i := newInputActionMap(input)
	i.addAction(action)
	sia.actionMap = append(sia.actionMap, i)
}

func (sia *SourceInputActionMap) getActionsByInput(name string) (*InputActionMap, error) {
	for _, m := range sia.actionMap {
		if m.getInput().GetName() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("No InputAction Map found for Input: %s", name)
}

type InputActionMap struct {
	input   *Input
	actions []Callback
}

func newInputActionMap(input *Input) *InputActionMap {
	return &InputActionMap{
		input:   input,
		actions: make([]Callback, 0),
	}
}

func (iam *InputActionMap) getInput() *Input {
	return iam.input
}

func (iam *InputActionMap) getActions() []Callback {
	return iam.actions
}

func (iam *InputActionMap) addAction(action Callback) {
	iam.actions = append(iam.actions, action)
}
