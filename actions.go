package persephone

import (
	"fmt"
)

type Action interface {
	getState() *State
	getActions() []Callback
	add(Callback)
}

type Actions interface {
	getActionsByState(string) Action
	add(Action)
}

type EntryAction struct {
	state    *State
	actions []Callback
}

func newEntryAction(state *State) *EntryAction {
	return &EntryAction{
		state:    state,
		actions: make([]Callback, 0),
	}
}

func (ea *EntryAction) getState() *State {
	return ea.state
}

func (ea *EntryAction) add(action Callback) {
	ea.actions = append(ea.actions, action)
}

func (ea *EntryAction) getActions() []Callback {
	return ea.actions
}

type EntryActions struct {
	entryActions []Action
}

func newEntryActions() *EntryActions {
	return &EntryActions{
		entryActions: make([]Action, 0),
	}
}

func (eas *EntryActions) add(ea Action) {
	eas.entryActions = append(eas.entryActions, ea)
}

func (eas *EntryActions) getActionsByState(name string) (Action, error) {
	for _, action := range eas.entryActions {
		if action.getState().GetName() == name {
			return action, nil
		}
	}

	return nil, fmt.Errorf("No EntryAction found for state: %s", name)
}

type ExitAction struct {
	state    *State
	actions []Callback
}

func newExitAction(state *State) *ExitAction {
	return &ExitAction{
		state:    state,
		actions: make([]Callback, 0),
	}
}

func (ea *ExitAction) getState() *State {
	return ea.state
}

func (ea *ExitAction) getActions() []Callback {
	return ea.actions
}

func (ea *ExitAction) add(action Callback) {
	ea.actions = append(ea.actions, action)
}

type ExitActions struct {
	exitActions []Action
}

func newExitActions() *ExitActions {
	return &ExitActions{
		exitActions: make([]Action, 0),
	}
}

func (eas *ExitActions) add(ea Action) {
	eas.exitActions = append(eas.exitActions, ea)
}

func (eas *ExitActions) getActionsByState(name string) (Action, error) {
	for _, action := range eas.exitActions {
		if action.getState().GetName() == name {
			return action, nil
		}
	}

	return nil, fmt.Errorf("No ExitAction found for state: %s", name)
}
