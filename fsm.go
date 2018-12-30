package persephone

import (
	"fmt"
	"sync"
)

type AbstractFSM struct {
	states            	*States
	inputs            	*Inputs
	rules             	*TransitionRules
	inputActions      	*InputActions
	transitionActions 	*TransitionActions
	entryActions      	*EntryActions
	exitActions       	*ExitActions
	currentState      	*State
	stateMu 			sync.RWMutex
	actionMu 			sync.RWMutex
}

func New(states *States, inputs *Inputs) *AbstractFSM {
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
	qp.stateMu.Lock()
	defer qp.stateMu.Unlock()
	if qp.currentState == nil {
		qp.currentState = qp.states.GetInitialState()
	}
	return
}

func (qp *AbstractFSM) AddRule(sourceState, input, targetState string, inputAction FSMMethod) {
	src, err := qp.states.GetState(sourceState)
	if err != nil {
		panic("Undefined source state: " + err.Error())
	}

	tgt, err := qp.states.GetState(targetState)
	if err != nil {
		panic("Undefined target state: " + err.Error())
	}

	in, err := qp.inputs.GetInput(input)
	if err != nil {
		panic("Undefined input: " + err.Error())
	}

	if qp.rules == nil {
		qp.rules = newTransitionRules()
	}

	if err := qp.rules.addRule(src, in, tgt); err != nil {
		panic(err.Error())
	}

	if inputAction != nil {
		qp.AddInputAction(sourceState, input, inputAction)
	}

}

func (qp *AbstractFSM) AddEntryAction(state string, action FSMMethod) {
	qp.actionMu.Lock()
	defer qp.actionMu.Unlock()
	act := NewFSMAction(action)

	st, err := qp.states.GetState(state)
	if err != nil {
		panic("Undefined state.")
	}

	if qp.entryActions == nil {
		qp.entryActions = newEntryActions()
	}

	ea, err := qp.entryActions.getActionsByState(st.GetName())

	if err != nil {
		// no entry action for given state yet
		ea := newEntryAction(st)
		ea.add(act)
		qp.entryActions.add(ea)
		return
	} 

	ea.add(act)

	return

}

func (qp *AbstractFSM) AddExitAction(state string, action FSMMethod) {
	qp.actionMu.Lock()
	defer qp.actionMu.Unlock()
	act := NewFSMAction(action)

	st, err := qp.states.GetState(state)
	if err != nil {
		panic("Undefined state.")
	}

	if qp.exitActions == nil {
		qp.exitActions = newExitActions()
	}

	ea, err := qp.exitActions.getActionsByState(st.GetName())

	if err != nil {
		// no exit action for given state yet
		ea := newEntryAction(st)
		ea.add(act)
		qp.exitActions.add(ea)
		return
	}

	ea.add(act)

	return

}

func (qp *AbstractFSM) AddInputAction(state, input string, action FSMMethod) {
	qp.actionMu.Lock()
	defer qp.actionMu.Unlock()
	act := NewFSMAction(action)

	src, err := qp.states.GetState(state)

	if err != nil {
		panic("Undefined state: " + err.Error())
	}

	in, err := qp.inputs.GetInput(input)

	if err != nil {
		panic("Undefined input symbol: " + err.Error())
	}

	if qp.inputActions == nil {
		qp.inputActions = newInputActions()
	}

	qp.inputActions.addMap(src, in, act)

	return

}

func (qp *AbstractFSM) AddTransitionAction(sourceState, targetState string, action FSMMethod) {
	qp.actionMu.Lock()
	defer qp.actionMu.Unlock()
	act := NewFSMAction(action)

	src, err_s := qp.states.GetState(sourceState)

	if err_s != nil {
		panic("Undefined source state: " + err_s.Error())
	}

	tgt, err_t := qp.states.GetState(targetState)

	if err_t != nil {
		panic("Undefined target state: " + err_t.Error())
	}

	if qp.transitionActions == nil {
		qp.transitionActions = newTransitionActions()
	}

	qp.transitionActions.addMap(src, tgt, act)

	return

}

func (qp *AbstractFSM) Process(input string) error {
	qp.actionMu.RLock()
	defer qp.actionMu.RUnlock()

	qp.stateMu.Lock()
	defer qp.stateMu.Unlock()

	in, err_i := qp.inputs.GetInput(input)

	if err_i != nil {
		panic("Undefined input symbol: " + err_i.Error())
	}

	r, err := qp.rules.getRuleBySource(qp.currentState.GetName())

	if err != nil {
		return fmt.Errorf("There are no rules left for current state %s: %s", qp.currentState.GetName(), err.Error())
	}

	ti, err := r.getTransitionByInput(in.GetName())

	if err != nil {
		return fmt.Errorf("There are no rule/s for {%s, %s} pair: %s", qp.currentState.GetName(), input, err.Error())
	}

	targetState := ti.getDestination()

	sourceState := qp.currentState

	// do exit actions before leaving state
	if exa, err := qp.exitActions.getActionsByState(sourceState.GetName()); err == nil {
		for _, action := range exa.getActions() {
			err := action.Do()
			if err != nil {
				return err
			}
		}
	}

	// do actions on recognized input
	if sia, err := qp.inputActions.getMapByState(sourceState.GetName()); err == nil {
		if iam, err := sia.getActionsByInput(in.GetName()); err == nil {
			for _, action := range iam.getActions() {
				err := action.Do()
				if err != nil {
					return err
				}
			}
		}
	}

	qp.currentState = targetState

	// do actions after transition
	if sta, err := qp.transitionActions.getMapBySource(sourceState.GetName()); err == nil {
		if tam, err := sta.getActionsByTarget(targetState.GetName()); err == nil {
			for _, action := range tam.getActions() {
				err := action.Do()
				if err != nil {
					return err
				}
			}
		}
	}

	// do entry actions on new state
	if ena, err := qp.entryActions.getActionsByState(targetState.GetName()); err == nil {
		for _, action := range ena.getActions() {
			err := action.Do()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (qp *AbstractFSM) Can(input string) (ok bool) {
	qp.stateMu.RLock()
	defer qp.stateMu.RUnlock()
	if r, err := qp.rules.getRuleBySource(qp.currentState.GetName()); err != nil {
		// no rule for current source yet
		return
	} else if _, err := r.getTransitionByInput(input); err != nil {
		// no rule for input yet
		return
	}

	return true
}

func (qp *AbstractFSM) GetState() *State {
	qp.stateMu.RLock()
	defer qp.stateMu.RUnlock()
	return qp.currentState
}

func (qp *AbstractFSM) SetState(state string) {
	qp.stateMu.Lock()
	defer qp.stateMu.Unlock()
	if src, err := qp.states.GetState(state); err != nil {
		panic("Undefined state: " + err.Error())
	} else {
		qp.currentState = src
	}

	return
}

func (qp *AbstractFSM) Reset() {
	qp.stateMu.Lock()
	defer qp.stateMu.Unlock()
	qp.currentState = qp.states.GetInitialState()
	return
}
