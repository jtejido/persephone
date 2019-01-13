package persephone

import (
	"fmt"
)

type TransitionActions struct {
	actionMap map[State]*TargetActionMap
}

func newTransitionActions() *TransitionActions {
	return &TransitionActions{
		actionMap: make(map[State]*TargetActionMap),
	}
}

func (ta *TransitionActions) addMap(src, tgt State, act Callback) (err error) {
	sta, err := ta.getMapBySource(src)

	if err != nil {
		return
	}

	if sta != nil {
		sta.addAction(tgt, act)
		return
	}

	m := newTargetActionMap()
	m.addAction(tgt, act)
	ta.actionMap[src] = m
	return

}

func (ta *TransitionActions) getMapBySource(src State) (*TargetActionMap, error) {
	v, ok := ta.actionMap[src]

	if ok {
		return v, nil
	}

	return nil, fmt.Errorf("No TargetAction Map found for Source State: %d", src)
}

type TargetActionMap struct {
	actionMap map[State][]Callback
}

func newTargetActionMap() *TargetActionMap {
	return &TargetActionMap{
		actionMap: make(map[State][]Callback),
	}
}

func (ta *TargetActionMap) addAction(state State, action Callback) {
	ta.actionMap[state] = append(ta.actionMap[state], action)
}

func (ta *TargetActionMap) getActionsByTarget(state State) ([]Callback, error) {
	v, ok := ta.actionMap[state]

	if !ok {
		return nil, fmt.Errorf("No TargetAction Map found for Target State: %d", state)
	}

	return v, nil
}
