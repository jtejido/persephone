package persephone

import (
	"fmt"
)

type Action interface {
	GetName() string
	GetActions() []*FSMAction
	Add(action *FSMAction)
}

type Actions interface {
	GetActionsByStateName(string) Action
	Add(Action)
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

func (ea *EntryAction) GetName() string {
	return ea.name
}

func (ea *EntryAction) Add(action *FSMAction) {
	ea.actions = append(ea.actions, action)
}

func (ea *EntryAction) GetActions() []*FSMAction {
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

func (eas *EntryActions) Add(ea Action) {
	eas.entryActions = append(eas.entryActions, ea)
}

func (eas *EntryActions) GetActionsByStateName(name string) (Action, error) {
	for _, action := range eas.entryActions {
		if action.GetName() == name {
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

func (ea *ExitAction) GetName() string {
	return ea.name
}

func (ea *ExitAction) GetActions() []*FSMAction {
	return ea.actions
}

func (ea *ExitAction) Add(action *FSMAction) {
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

func (eas *ExitActions) Add(ea Action) {
	eas.exitActions = append(eas.exitActions, ea)
}

func (eas *ExitActions) GetActionsByStateName(name string) (Action, error) {
	for _, action := range eas.exitActions {
		if action.GetName() == name {
			return action, nil
		}
	}

	return nil, fmt.Errorf("No ExitAction found for name: %s", name)
}
