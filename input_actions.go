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

func (ia *InputActions) AddMap(m *SourceInputActionMap) {
	ia.actionMap = append(ia.actionMap, m)
}

func (ia *InputActions) GetMapByStateName(name string) (*SourceInputActionMap, error) {
	for _, m := range ia.actionMap {
		if m.GetState().GetName() == name {
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

func (sia *SourceInputActionMap) GetState() *State {
	return sia.state
}

func (sia *SourceInputActionMap) AddMap(input *Input, action *FSMAction) {
	i := newInputAction(input)
	i.AddAction(action)
	sia.actionMap = append(sia.actionMap, i)
}

func (sia *SourceInputActionMap) GetActionsByInputName(name string) (*InputActionMap, error) {
	for _, m := range sia.actionMap {
		if m.GetInput().GetName() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("No InputAction Map found for Input name: %s", name)
}

type InputActionMap struct {
	input   *Input
	actions []*FSMAction
}

func newInputAction(input *Input) *InputActionMap {
	return &InputActionMap{
		input:   input,
		actions: make([]*FSMAction, 0),
	}
}

func (iam *InputActionMap) GetInput() *Input {
	return iam.input
}

func (iam *InputActionMap) GetActions() []*FSMAction {
	return iam.actions
}

func (iam *InputActionMap) AddAction(action *FSMAction) {
	iam.actions = append(iam.actions, action)
}
