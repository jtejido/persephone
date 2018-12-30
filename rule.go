package persephone

import (
	"fmt"
)

// one source state can have multiple input => dest
type TransitionMap struct {
	source         	*State
	transitions 	[]*Transition
}

func newTransitionMap(src *State) *TransitionMap {
	return &TransitionMap{
		source:         src,
		transitions: 	make([]*Transition, 0),
	}
}

func (tr *TransitionMap) addTransition(input *Input, dest *State) {
	tr.transitions = append(tr.transitions, newTransition(input, dest))
}

func (tr *TransitionMap) getSource() *State {
	return tr.source
}

func (tr *TransitionMap) getTransitions() []*Transition {
	return tr.transitions
}

func (tr *TransitionMap) getTransitionByInput(name string) (*Transition, error) {
	for _, transition := range tr.transitions {
		if transition.getInput().GetName() == name {
			return transition, nil
		}
	}

	return nil, fmt.Errorf("No Transition found for Input: %s", name)
}

// only one destination for one input
type Transition struct {
	input *Input
	dest  *State
}

func newTransition(input *Input, dest *State) *Transition {
	return &Transition{
		input: input,
		dest:  dest,
	}
}

func (t *Transition) getInput() *Input {
	return t.input
}

func (t *Transition) getDestination() *State {
	return t.dest
}

// State transition table
type TransitionRules struct {
	rules []*TransitionMap
}

func newTransitionRules() *TransitionRules {
	return &TransitionRules{
		rules: make([]*TransitionMap, 0),
	}
}

func (rs *TransitionRules) addRule(src *State, in *Input, tgt *State) error {

	if rs.rules == nil {
		rs.rules = make([]*TransitionMap, 0)
	}

	r, _ := rs.getRuleBySource(src.GetName())

	if r != nil {
		// means there's rule for source already, just add input=>dest map
		ti, _ := r.getTransitionByInput(in.GetName())

		if ti != nil {
			// only one dest by input allowed
			return fmt.Errorf("Rule for {state,input} pair is already defined.")
		}

		r.addTransition(in, tgt)
		return nil
	}

	tm := newTransitionMap(src)
	tm.addTransition(in, tgt)
	rs.rules = append(rs.rules, tm)
	return nil
}

func (rs *TransitionRules) getRuleBySource(name string) (*TransitionMap, error) {
	for _, rule := range rs.rules {
		if rule.getSource().GetName() == name {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("No Rule found at Source State: %s", name)
}
