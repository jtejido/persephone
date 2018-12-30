package persephone

import (
	"fmt"
)

type InputActions struct {
	actionMap []*SourceInputActionMap
}

func newInputActions() *InputActions {
	return &InputActions{
		actionMap: make([]*SourceInputActionMap, 0),
	}
}

func (ia *InputActions) addMap(m *SourceInputActionMap) {
	ia.actionMap = append(ia.actionMap, m)
}

func (ia *InputActions) getMapByStateName(name string) (*SourceInputActionMap, error) {
	for _, m := range ia.actionMap {
		if m.getState().GetName() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("No SourceInputAction Map found for State name: %s", name)
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

func (sia *SourceInputActionMap) addMap(input *Input, action *FSMAction) {
	i := newInputActionMap(input)
	i.addAction(action)
	sia.actionMap = append(sia.actionMap, i)
}

func (sia *SourceInputActionMap) getActionsByInputName(name string) (*InputActionMap, error) {
	for _, m := range sia.actionMap {
		if m.getInput().GetName() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("No InputAction Map found for Input name: %s", name)
}

type InputActionMap struct {
	input   *Input
	actions []*FSMAction
}

func newInputActionMap(input *Input) *InputActionMap {
	return &InputActionMap{
		input:   input,
		actions: make([]*FSMAction, 0),
	}
}

func (iam *InputActionMap) getInput() *Input {
	return iam.input
}

func (iam *InputActionMap) getActions() []*FSMAction {
	return iam.actions
}

func (iam *InputActionMap) addAction(action *FSMAction) {
	iam.actions = append(iam.actions, action)
}
