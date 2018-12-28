package persephone

import (
	"fmt"
)

type AbstractFSM struct {
	states            *States
	inputs            *Inputs
	rules             *Rules
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
		rules:             newRules(),
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
	src, err_s := qp.states.GetStateByName(sourceState)
	if err_s != nil {
		panic("Undefined source state: " + err_s.Error())
	}

	tgt, err_t := qp.states.GetStateByName(targetState)
	if err_t != nil {
		panic("Undefined target state: " + err_t.Error())
	}

	in, err_i := qp.inputs.GetInputByName(input)
	if err_i != nil {
		panic("Undefined input: " + err_i.Error())
	}

	if qp.rules == nil {
		qp.rules = newRules()
	}

	r, err := qp.rules.getRuleBySourceName(sourceState)

	if err != nil {
		qp.rules.addRule(newRule(src, in, tgt))
	} else if r.isDestinationByInputDefined(input) {
		panic("Rule for {state,input} pair is already defined.")
	} else {
		r.addTransition(in, tgt)
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
	var targetState *State

	in, err_i := qp.inputs.GetInputByName(input)

	if err_i != nil {
		panic("Undefined input symbol: " + err_i.Error())
	}

	if r, err := qp.rules.getRuleBySourceName(qp.currentState.GetName()); err != nil {
		return fmt.Errorf("There are no rules left for current state %s: %s", qp.currentState.GetName(), err.Error())
	} else if !r.isDestinationByInputDefined(input) {
		return fmt.Errorf("There are no other rules for {%s, %s} pair.", qp.currentState.GetName(), input)
	} else {
		targetState = r.getTransitionByInputName(input).getDestination()
	}

	sourceState := qp.currentState

	if exa, err := qp.exitActions.getActionsByStateName(sourceState.GetName()); err == nil {
		if sourceState.GetName() != targetState.GetName() {
			for _, action := range exa.getActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	if sia, err := qp.inputActions.getMapByStateName(sourceState.GetName()); err == nil {
		if iam, err := sia.getActionsByInputName(in.GetName()); err == nil {
			for _, action := range iam.getActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	qp.currentState = targetState

	if sta, err := qp.transitionActions.getMapBySourceName(sourceState.GetName()); err == nil {
		if tam, err := sta.getActionsByTargetName(targetState.GetName()); err == nil {
			for _, action := range tam.getActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	if ena, err := qp.entryActions.getActionsByStateName(targetState.GetName()); err == nil {
		if sourceState.GetName() != targetState.GetName() {
			for _, action := range ena.getActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (qp *AbstractFSM) Can(input string) (ok bool) {
	if r, err := qp.rules.getRuleBySourceName(qp.currentState.GetName()); err != nil {
		return
	} else if !r.isDestinationByInputDefined(input) {
		return
	} else {
		ok = true
	}

	return
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
