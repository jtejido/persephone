package persephone

import (
	"fmt"
)

type TransitionActions struct {
	actionMap map[string]*SourceTargetActionMap
}

func newTransitionActions() *TransitionActions {
	return &TransitionActions{
		actionMap: make(map[string]*SourceTargetActionMap),
	}
}

func (ta *TransitionActions) addMap(src, tgt *State, act Callback) {

	if ta.actionMap == nil {
		ta.actionMap = make(map[string]*SourceTargetActionMap)
	}

	sta, _ := ta.getMapBySource(src.GetName())

	if sta != nil {
		at, err := sta.getActionsByTarget(tgt.GetName())

		if err != nil {
			sta.addMap(tgt, act)

			return
		}

		at.addAction(act)

		return
	}

	m := newSourceTargetActionMap(src)
	m.addMap(tgt, act)
	ta.actionMap[src.GetName()] = m

}

func (ta *TransitionActions) getMapBySource(name string) (*SourceTargetActionMap, error) {

	v, ok := ta.actionMap[name]

	if ok {
		return v, nil
	}

	return nil, fmt.Errorf("No SourceTargetAction Map found for Source State: %s", name)
}

type SourceTargetActionMap struct {
	source    *State
	actionMap map[string]*TargetActionMap
}

func newSourceTargetActionMap(state *State) *SourceTargetActionMap {
	return &SourceTargetActionMap{
		source:    state,
		actionMap: make(map[string]*TargetActionMap, 0),
	}
}

func (sta *SourceTargetActionMap) getSource() *State {
	return sta.source
}

func (sta *SourceTargetActionMap) addMap(state *State, action Callback) {
	if sta.actionMap == nil {
		sta.actionMap = make(map[string]*TargetActionMap)
	}
	
	ta := newTargetAction(state)
	ta.addAction(action)
	sta.actionMap[state.GetName()] = ta
}

func (sta *SourceTargetActionMap) getActionsByTarget(name string) (*TargetActionMap, error) {
	v, ok := sta.actionMap[name]

	if ok {
		return v, nil
	}

	return nil, fmt.Errorf("No TargetAction Map found for Target State: %s", name)
}

type TargetActionMap struct {
	target  *State
	actions []Callback
}

func newTargetAction(state *State) *TargetActionMap {
	return &TargetActionMap{
		target:  state,
		actions: make([]Callback, 0),
	}
}

func (tam *TargetActionMap) getTarget() *State {
	return tam.target
}

func (tam *TargetActionMap) getActions() []Callback {
	return tam.actions
}

func (tam *TargetActionMap) addAction(action Callback) {
	tam.actions = append(tam.actions, action)
}
