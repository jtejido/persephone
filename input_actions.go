package persephone

import (
	"fmt"
)

type InputActions struct {
	actionMap map[State]*InputActionMap
}

func newInputActions() *InputActions {
	return &InputActions{
		actionMap: make(map[State]*InputActionMap),
	}
}

func (ia *InputActions) addMap(src State, in Input, act Callback) (err error) {

	sia, _ := ia.getMapByState(src)

	if sia != nil {
		sia.addAction(in, act)
		return
	}

	m := newInputActionMap()
	m.addAction(in, act)
	ia.actionMap[src] = m
	return
}

func (ia *InputActions) getMapByState(state State) (*InputActionMap, error) {

	v, ok := ia.actionMap[state]

	if ok {
		return v, nil
	}

	return nil, fmt.Errorf("No InputAction Map found for State: %d", state)
}

type InputActionMap struct {
	actionMap map[Input][]Callback
}

func newInputActionMap() *InputActionMap {
	return &InputActionMap{
		actionMap: make(map[Input][]Callback, 0),
	}
}

func (ia *InputActionMap) addAction(input Input, action Callback) {
	ia.actionMap[input] = append(ia.actionMap[input], action)
}

func (ia *InputActionMap) getActionsByInput(input Input) ([]Callback, error) {

	v, ok := ia.actionMap[input]

	if !ok {
		return nil, fmt.Errorf("No Actions found for Input: %d", input)
	}

	return v, nil

}
