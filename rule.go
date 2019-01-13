package persephone

import (
	"fmt"
)

// one source state can have multiple input => dest
type TransitionMap struct {
	transitions map[Input]State
}

func newTransitionMap() *TransitionMap {
	return &TransitionMap{
		transitions: make(map[Input]State, 0),
	}
}

func (tr *TransitionMap) addTransition(input Input, dest State) {
	tr.transitions[input] = dest
}

func (tr *TransitionMap) getDestinationByInput(input Input) (State, error) {

	v, ok := tr.transitions[input]

	if !ok {
		return 0, fmt.Errorf("No Transition found for Input: %d", input)
	}

	return v, nil

}

// State transition table
type TransitionRules struct {
	rules map[State]*TransitionMap
}

func newTransitionRules() *TransitionRules {
	return &TransitionRules{
		rules: make(map[State]*TransitionMap, 0),
	}
}

func (rs *TransitionRules) addRule(src State, in Input, tgt State) error {

	if rs.rules == nil {
		rs.rules = make(map[State]*TransitionMap, 0)
	}

	r, _ := rs.getRuleBySource(src)

	if r != nil {
		// means there's rule for source already, just add input=>dest map
		_, err := r.getDestinationByInput(in)

		if err == nil {
			// only one dest by input allowed
			return fmt.Errorf("Rule for {state,input} pair is already defined.")
		}

		// replace whatever's in there
		r.addTransition(in, tgt)
		return nil
	}

	tm := newTransitionMap()
	tm.addTransition(in, tgt)
	rs.rules[src] = tm
	return nil
}

func (rs *TransitionRules) getRuleBySource(src State) (*TransitionMap, error) {

	v, ok := rs.rules[src]

	if !ok {
		return nil, fmt.Errorf("No Rule found at Source State: %d", src)
	}

	return v, nil
}
