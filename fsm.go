package persephone

import (
	"fmt"
)

type AbstractFSM struct {
	states            *States
	inputs            *Inputs
	rules             *TransitionRules
	inputActions      *InputActions
	transitionActions *TransitionActions
	entryActions      *EntryActions
	exitActions       *ExitActions
	currentState      *State
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
	if qp.currentState == nil {
		qp.currentState = qp.states.GetInitialState()
	}
}

func (qp *AbstractFSM) AddRule(sourceState, input, targetState string, inputAction FSMMethod) {
	src, err := qp.states.GetStateByName(sourceState)
	if err != nil {
		panic("Undefined source state: " + err.Error())
	}

	tgt, err := qp.states.GetStateByName(targetState)
	if err != nil {
		panic("Undefined target state: " + err.Error())
	}

	in, err := qp.inputs.GetInputByName(input)
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
	act := NewFSMAction(action)

	st, err := qp.states.GetStateByName(state)
	if err != nil {
		panic("Undefined state.")
	}

	if qp.entryActions == nil {
		qp.entryActions = newEntryActions()
	}

	if ea, err := qp.entryActions.getActionsByStateName(st.GetName()); err != nil {
		// no entry action for given state yet
		ea := newEntryAction(st.GetName())
		ea.add(act)
		qp.entryActions.add(ea)
	} else {
		ea.add(act)
	}

}

func (qp *AbstractFSM) AddExitAction(state string, action FSMMethod) {
	act := NewFSMAction(action)

	st, err := qp.states.GetStateByName(state)
	if err != nil {
		panic("Undefined state.")
	}

	if qp.exitActions == nil {
		qp.exitActions = newExitActions()
	}

	if ea, err := qp.exitActions.getActionsByStateName(st.GetName()); err != nil {
		// no exit action for given state yet
		ea := newEntryAction(st.GetName())
		ea.add(act)
		qp.exitActions.add(ea)
	} else {
		ea.add(act)
	}

}

func (qp *AbstractFSM) AddInputAction(state, inputSymbol string, action FSMMethod) {
	act := NewFSMAction(action)

	src, err_s := qp.states.GetStateByName(state)
	if err_s != nil {
		panic("Undefined state: " + err_s.Error())
	}

	in, err_i := qp.inputs.GetInputByName(inputSymbol)
	if err_i != nil {
		panic("Undefined input symbol: " + err_i.Error())
	}

	if qp.inputActions == nil {
		qp.inputActions = newInputActions()
	}

	if sia, err := qp.inputActions.getMapByStateName(src.GetName()); err != nil {
		m := newSourceInputActionMap(src)
		m.addMap(in, act)
		qp.inputActions.addMap(m)
	} else if ia, err := sia.getActionsByInputName(in.GetName()); err != nil {
		sia.addMap(in, act)
	} else {
		ia.addAction(act)
	}

}

func (qp *AbstractFSM) AddTransitionAction(sourceState, targetState string, action FSMMethod) {
	act := NewFSMAction(action)

	src, err_s := qp.states.GetStateByName(sourceState)
	if err_s != nil {
		panic("Undefined source state: " + err_s.Error())
	}

	tgt, err_t := qp.states.GetStateByName(targetState)
	if err_t != nil {
		panic("Undefined target state: " + err_t.Error())
	}

	if qp.transitionActions == nil {
		qp.transitionActions = newTransitionActions()
	}

	if sta, err := qp.transitionActions.getMapBySourceName(src.GetName()); err != nil {
		m := newSourceTargetActionMap(src)
		m.addMap(tgt, act)
		qp.transitionActions.addMap(m)
	} else if ta, err := sta.getActionsByTargetName(tgt.GetName()); err != nil {
		sta.addMap(tgt, act)
	} else {
		ta.addAction(act)
	}

}

func (qp *AbstractFSM) Process(input string) error {

	qp.init()

	in, err_i := qp.inputs.GetInputByName(input)

	if err_i != nil {
		panic("Undefined input symbol: " + err_i.Error())
	}

	r, err := qp.rules.getRuleBySourceName(qp.currentState.GetName())

	if err != nil {
		return fmt.Errorf("There are no rules left for current state %s: %s", qp.currentState.GetName(), err.Error())
	}

	ti, err := r.getTransitionByInputName(input)

	if err != nil {
		return fmt.Errorf("There are no rule/s for {%s, %s} pair: %s", qp.currentState.GetName(), input, err.Error())
	}

	targetState := ti.getDestination()

	sourceState := qp.currentState

	// do exit actions before leaving state
	if exa, err := qp.exitActions.getActionsByStateName(sourceState.GetName()); err == nil {
		for _, action := range exa.getActions() {
			err := action.Do()
			if err != nil {
				return err
			}
		}
	}

	// do actions on recognized input
	if sia, err := qp.inputActions.getMapByStateName(sourceState.GetName()); err == nil {
		if iam, err := sia.getActionsByInputName(in.GetName()); err == nil {
			for _, action := range iam.getActions() {
				err := action.Do()
				if err != nil {
					return err
				}
			}
		}
	}

	qp.currentState = targetState

	// do transition actions
	if sta, err := qp.transitionActions.getMapBySourceName(sourceState.GetName()); err == nil {
		if tam, err := sta.getActionsByTargetName(targetState.GetName()); err == nil {
			for _, action := range tam.getActions() {
				err := action.Do()
				if err != nil {
					return err
				}
			}
		}
	}

	// do entry actions on new state
	if ena, err := qp.entryActions.getActionsByStateName(targetState.GetName()); err == nil {
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
	if r, err := qp.rules.getRuleBySourceName(qp.currentState.GetName()); err != nil {
		// no rule for current source yet
		return
	} else if _, err := r.getTransitionByInputName(input); err != nil {
		// no rule for input yet
		return
	}

	return true
}

func (qp *AbstractFSM) GetState() *State {
	return qp.currentState
}

func (qp *AbstractFSM) SetState(state string) {
	if src, err := qp.states.GetStateByName(state); err != nil {
		panic("Undefined state: " + err.Error())
	} else {
		qp.currentState = src
	}
}

func (qp *AbstractFSM) Reset() {
	qp.currentState = qp.states.GetInitialState()
}
