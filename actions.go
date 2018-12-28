package persephone

import (
	"fmt"
)

type Action interface {
	getName() string
	getActions() []*FSMAction
	add(action *FSMAction)
}

type Actions interface {
	getActionsByStateName(string) Action
	add(Action)
}

type EntryAction struct {
	name    string
	actions []*FSMAction
}

func newEntryAction(name string) *EntryAction {
	return &EntryAction{
		name:    name,
		actions: make([]*FSMAction, 0),
	}
}

func (ea *EntryAction) getName() string {
	return ea.name
}

func (ea *EntryAction) add(action *FSMAction) {
	ea.actions = append(ea.actions, action)
}

func (ea *EntryAction) getActions() []*FSMAction {
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

func (eas *EntryActions) getActionsByStateName(name string) (Action, error) {
	for _, action := range eas.entryActions {
		if action.getName() == name {
			return action, nil
		}
	}

	return nil, fmt.Errorf("No EntryAction found for name: %s", name)
}

type ExitAction struct {
	name    string
	actions []*FSMAction
}

func newExitAction(name string) *ExitAction {
	return &ExitAction{
		name:    name,
		actions: make([]*FSMAction, 0),
	}
}

func (ea *ExitAction) getName() string {
	return ea.name
}

func (ea *ExitAction) getActions() []*FSMAction {
	return ea.actions
}

func (ea *ExitAction) add(action *FSMAction) {
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

func (eas *ExitActions) getActionsByStateName(name string) (Action, error) {
	for _, action := range eas.exitActions {
		if action.getName() == name {
			return action, nil
		}
	}

	return nil, fmt.Errorf("No ExitAction found for name: %s", name)
}
