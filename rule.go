package persephone

import (
	"fmt"
)

// one source state can have multiple input => dest
type Rule struct {
	src         *State
	transitions []*Transition
}

func newRule(src *State, input *Input, dest *State) *Rule {
	r := &Rule{
		src:         src,
		transitions: make([]*Transition, 0),
	}

	r.transitions = append(r.transitions, newTransition(input, dest))

	return r
}

func (r *Rule) addTransition(input *Input, dest *State) {
	r.transitions = append(r.transitions, newTransition(input, dest))
}

func (r *Rule) getSource() *State {
	return r.src
}

func (r *Rule) getTransitions() []*Transition {
	return r.transitions
}

func (r *Rule) getTransitionByInputName(name string) *Transition {
	for _, transition := range r.transitions {
		if transition.getInput().GetName() == name {
			return transition
		}
	}

	return nil
}

func (r *Rule) isDestinationByInputDefined(name string) bool {
	for _, transition := range r.transitions {
		if transition.getInput().GetName() == name {
			return true
		}
	}

	return false
}

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
type Rules struct {
	rules []*Rule
}

func newRules() *Rules {
	return &Rules{
		rules: make([]*Rule, 0),
	}
}

func (rs *Rules) addRule(rule *Rule) {
	rs.rules = append(rs.rules, rule)
}

func (rs *Rules) getRuleBySourceName(name string) (*Rule, error) {
	for _, rule := range rs.rules {
		if rule.getSource().GetName() == name {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("No Rule found at SourceState name: %s", name)
}
