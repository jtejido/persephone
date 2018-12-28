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
	return &AbstractFSM{
		states:            states,
		inputs:            inputs,
		rules:             newRules(),
		inputActions:      newInputActions(),
		entryActions:      newEntryActions(),
		transitionActions: newTransitionActions(),
		exitActions:       newExitActions(),
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

	r, err := qp.rules.GetRuleBySourceName(sourceState)

	if err != nil {
		qp.rules.addRule(NewRule(src, in, tgt))
	} else {
		if r.IsDestinationByInputDefined(input) {
			panic("Rule for {state,input} pair is already defined.")
		} else {
			r.AddTransition(in, tgt)
		}
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

	if ea, err := qp.entryActions.GetActionsByStateName(st.GetName()); err != nil {
		ea := newEntryAction(st.GetName())
		ea.Add(act)
		qp.entryActions.Add(ea)
	} else {
		ea.Add(act)
	}

}

func (qp *AbstractFSM) AddExitAction(state string, action FSMMethod) {
	act := NewFSMAction(action)

	st, err := qp.states.GetStateByName(state)
	if err != nil {
		panic("Undefined state.")
	}

	if ea, err := qp.exitActions.GetActionsByStateName(st.GetName()); err != nil {
		ea := newEntryAction(st.GetName())
		ea.Add(act)
		qp.exitActions.Add(ea)
	} else {
		ea.Add(act)
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

	sia, err := qp.inputActions.GetMapByStateName(src.GetName())

	if err != nil {
		m := newSourceInputActionMap(src)
		m.AddMap(in, act)
		qp.inputActions.AddMap(m)
	} else {
		if inAct, err_in := sia.GetActionsByInputName(in.GetName()); err_in != nil {
			sia.AddMap(in, act)
		} else {
			inAct.AddAction(act)
		}
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

	sta, err := qp.transitionActions.GetMapBySourceName(src.GetName())

	if err != nil {
		m := newSourceTargetActionMap(src)
		m.AddMap(tgt, act)
		qp.transitionActions.AddMap(m)
	} else {
		if tAct, err_in := sta.GetActionsByTargetName(tgt.GetName()); err_in != nil {
			sta.AddMap(tgt, act)
		} else {
			tAct.AddAction(act)
		}
	}

}

func (qp *AbstractFSM) Process(input string) error {
	var targetState *State

	if qp.currentState == nil {
		qp.currentState = qp.states.GetInitialState()
	}

	in, err_i := qp.inputs.GetInputByName(input)

	if err_i != nil {
		panic("Undefined input symbol: " + err_i.Error())
	}

	if r, err := qp.rules.GetRuleBySourceName(qp.currentState.GetName()); err != nil {
		return fmt.Errorf("There are no rules left for current state %s: %s", qp.currentState.GetName(), err.Error())
	} else {
		if !r.IsDestinationByInputDefined(input) {
			return fmt.Errorf("There are no other rules for {%s, %s} pair.", qp.currentState.GetName(), input)
		} else {
			targetState = r.GetTransitionByInputName(input).GetDestination()
		}
	}

	sourceState := qp.currentState

	if exa, err := qp.exitActions.GetActionsByStateName(sourceState.GetName()); err == nil {
		if sourceState.GetName() != targetState.GetName() {
			for _, action := range exa.GetActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	if sia, err := qp.inputActions.GetMapByStateName(sourceState.GetName()); err == nil {
		if iam, err := sia.GetActionsByInputName(in.GetName()); err == nil {
			for _, action := range iam.GetActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	qp.currentState = targetState

	if sta, err := qp.transitionActions.GetMapBySourceName(sourceState.GetName()); err == nil {
		if tam, err := sta.GetActionsByTargetName(targetState.GetName()); err == nil {
			for _, action := range tam.GetActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	if ena, err := qp.entryActions.GetActionsByStateName(targetState.GetName()); err == nil {
		if sourceState.GetName() != targetState.GetName() {
			for _, action := range ena.GetActions() {
				err := action.DoAction()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
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
