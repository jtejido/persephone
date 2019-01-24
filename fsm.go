package persephone

import (
	"fmt"
	"sync"
)

type AbstractFSM struct {
	states            States
	inputs            Inputs
	rules             *TransitionRules
	inputActions      *InputActions
	transitionActions *TransitionActions
	entryActions      *EntryActions
	exitActions       *ExitActions
	currentState      State
	mu                sync.RWMutex
}

func New(states States, inputs Inputs) *AbstractFSM {
	fsm := &AbstractFSM{
		states:            states,
		inputs:            inputs,
		rules:             newTransitionRules(),
		inputActions:      newInputActions(),
		entryActions:      newEntryActions(),
		transitionActions: newTransitionActions(),
		exitActions:       newExitActions(),
	}

	fsm.init()

	return fsm
}

func (qp *AbstractFSM) init() {
	qp.mu.Lock()
	defer qp.mu.Unlock()
	qp.currentState = qp.states.GetInitialState()
}

func (qp *AbstractFSM) AddRule(sourceState State, input Input, targetState State, inputAction FSMMethod) (err error) {
	qp.mu.Lock()

	if !qp.states.Contains(sourceState) {
		panic("Undefined source state")
	}

	if !qp.states.Contains(targetState) {
		panic("Undefined target state")
	}

	if !qp.inputs.Contains(input) {
		panic("Undefined input")
	}

	if qp.rules == nil {
		qp.rules = newTransitionRules()
	}

	if err = qp.rules.addRule(sourceState, input, targetState); err != nil {
		return
	}

	qp.mu.Unlock()

	if inputAction != nil {
		qp.AddInputAction(sourceState, input, inputAction)
	}

	return

}

func (qp *AbstractFSM) AddEntryAction(state State, action FSMMethod) (err error) {
	qp.mu.Lock()
	defer qp.mu.Unlock()

	if !qp.states.Contains(state) {
		panic("Undefined state.")
	}

	if qp.entryActions == nil {
		qp.entryActions = newEntryActions()
	}

	act := NewFSMAction(action)

	ea, err := qp.entryActions.getActionsByState(state)

	if err != nil {
		// no entry action for given state yet
		ea := newEntryAction(state)
		ea.add(act)
		qp.entryActions.add(ea)
		return
	}

	ea.add(act)

	return

}

func (qp *AbstractFSM) AddExitAction(state State, action FSMMethod) (err error) {
	qp.mu.Lock()
	defer qp.mu.Unlock()

	if !qp.states.Contains(state) {
		panic("Undefined state.")
	}

	if qp.exitActions == nil {
		qp.exitActions = newExitActions()
	}

	act := NewFSMAction(action)

	ea, err := qp.exitActions.getActionsByState(state)

	if err != nil {
		// no exit action for given state yet
		ea := newExitAction(state)
		ea.add(act)
		qp.exitActions.add(ea)
		return
	}

	ea.add(act)

	return

}

func (qp *AbstractFSM) AddInputAction(state State, input Input, action FSMMethod) (err error) {
	qp.mu.Lock()
	defer qp.mu.Unlock()

	if !qp.states.Contains(state) {
		panic("Undefined state.")
	}

	if !qp.inputs.Contains(input) {
		panic("Undefined input symbol. ")
	}

	if qp.inputActions == nil {
		qp.inputActions = newInputActions()
	}

	// do actions on recognized input
	return qp.inputActions.addMap(state, input, NewFSMAction(action))

}

func (qp *AbstractFSM) AddTransitionAction(sourceState, targetState State, action FSMMethod) (err error) {
	qp.mu.Lock()
	defer qp.mu.Unlock()

	if !qp.states.Contains(sourceState) {
		panic("Undefined source state.")
	}

	if !qp.states.Contains(targetState) {
		panic("Undefined target state.")
	}

	if qp.transitionActions == nil {
		qp.transitionActions = newTransitionActions()
	}

	return qp.transitionActions.addMap(sourceState, targetState, NewFSMAction(action))

}

func (qp *AbstractFSM) Process(input Input) error {
	qp.mu.Lock()
	defer qp.mu.Unlock()

	if !qp.inputs.Contains(input) {
		panic("Undefined input symbol.")
	}
	
	r, err := qp.rules.getRuleBySource(qp.currentState)

	if err != nil {
		return fmt.Errorf("There are no rules left for current state %d: %s", qp.currentState, err.Error())
	}

	targetState, err := r.getDestinationByInput(input)

	if err != nil {
		return fmt.Errorf("There are no rule/s for {%d, %d} pair: %s", qp.currentState, input, err.Error())
	}

	sourceState := qp.currentState

	// do exit actions before leaving state
	if exa, err := qp.exitActions.getActionsByState(sourceState); err == nil {
		for _, action := range exa.getActions() {
			err := action.Do()
			if err != nil {
				return err
			}
		}
	}

	// do actions on recognized input
	if sia, err := qp.inputActions.getMapByState(sourceState); err == nil {
		if actions, err := sia.getActionsByInput(input); err == nil {
			for _, action := range actions {
				err := action.Do()
				if err != nil {
					return err
				}
			}
		}
	}

	qp.currentState = targetState

	// do actions after transition
	if sta, err := qp.transitionActions.getMapBySource(sourceState); err == nil {
		if actions, err := sta.getActionsByTarget(targetState); err == nil {
			for _, action := range actions {
				err := action.Do()
				if err != nil {
					return err
				}
			}
		}
	}

	// do entry actions on new state
	if ena, err := qp.entryActions.getActionsByState(targetState); err == nil {
		for _, action := range ena.getActions() {
			err := action.Do()
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (qp *AbstractFSM) Can(input Input) (ok bool) {
	qp.mu.RLock()
	defer qp.mu.RUnlock()

	if r, err := qp.rules.getRuleBySource(qp.currentState); err != nil {
		// no rule for current source yet
		return
	} else if _, err := r.getDestinationByInput(input); err != nil {
		// no rule for input yet
		return
	}

	return true
}

func (qp *AbstractFSM) GetState() State {
	qp.mu.RLock()
	defer qp.mu.RUnlock()
	return qp.currentState
}

func (qp *AbstractFSM) SetState(state State) {
	qp.mu.RLock()
	defer qp.mu.RUnlock()
	if !qp.states.Contains(state) {
		panic("Undefined state.")
	}

	qp.currentState = state

}

func (qp *AbstractFSM) Reset() {
	qp.mu.Lock()
	defer qp.mu.Unlock()
	qp.currentState = qp.states.GetInitialState()
	return
}
