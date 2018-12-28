package persephone

import (
	"fmt"
)

type Input struct {
	name  string
	input interface{}
}

func NewInput(name string, input interface{}) *Input {
	return &Input{
		name:  name,
		input: input,
	}
}

func (input Input) GetInput() interface{} {
	return input.input
}

func (input Input) GetName() string {
	return input.name
}

type Inputs struct {
	inputs []*Input
}

func NewInputs() *Inputs {
	return &Inputs{
		inputs: make([]*Input, 0),
	}
}

func (inputs *Inputs) Add(name string, input interface{}) {
	st := &Input{
		name:  name,
		input: input,
	}

	inputs.inputs = append(inputs.inputs, st)
}

func (inputs Inputs) GetInputByName(name string) (*Input, error) {
	for _, input := range inputs.inputs {
		if input.GetName() == name {
			return input, nil
		}
	}

	return nil, fmt.Errorf("No Input found for name: %s", name)
}
