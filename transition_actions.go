package persephone

import (
	"fmt"
)

type TransitionActions struct {
	actionMap []*SourceTargetActionMap
}

func newTransitionActions() *TransitionActions {
	return &TransitionActions{
		actionMap: make([]*SourceTargetActionMap, 0),
	}
}

func (ta *TransitionActions) addMap(m *SourceTargetActionMap) {
	ta.actionMap = append(ta.actionMap, m)
}

func (ta *TransitionActions) getMapBySourceName(name string) (*SourceTargetActionMap, error) {
	for _, m := range ta.actionMap {
		if m.getSource().GetName() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("No SourceTargetAction Map found for Source State name: %s", name)
}

type SourceTargetActionMap struct {
	source    *State
	actionMap []*TargetActionMap
}

func newSourceTargetActionMap(state *State) *SourceTargetActionMap {
	return &SourceTargetActionMap{
		source:    state,
		actionMap: make([]*TargetActionMap, 0),
	}
}

func (sta *SourceTargetActionMap) getSource() *State {
	return sta.source
}

func (sta *SourceTargetActionMap) addMap(state *State, action *FSMAction) {
	ta := newTargetAction(state)
	ta.addAction(action)
	sta.actionMap = append(sta.actionMap, ta)
}

func (sta *SourceTargetActionMap) getActionsByTargetName(name string) (*TargetActionMap, error) {
	for _, m := range sta.actionMap {
		if m.getTarget().GetName() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("No TargetAction Map found for Target State name: %s", name)
}

type TargetActionMap struct {
	target  *State
	actions []*FSMAction
}

func newTargetAction(state *State) *TargetActionMap {
	return &TargetActionMap{
		target:  state,
		actions: make([]*FSMAction, 0),
	}
}

func (tam *TargetActionMap) getTarget() *State {
	return tam.target
}

func (tam *TargetActionMap) getActions() []*FSMAction {
	return tam.actions
}

func (tam *TargetActionMap) addAction(action *FSMAction) {
	tam.actions = append(tam.actions, action)
}
