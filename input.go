package persephone

import (
	"fmt"
)

type Input struct {
	Name  string
	Input interface{}
}

func NewInput(name string, input interface{}) *Input {
	return &Input{
		Name:  name,
		Input: input,
	}
}

func (input Input) GetInput() interface{} {
	return input.Input
}

func (input Input) GetName() string {
	return input.Name
}

type Inputs struct {
	inputs map[string]*Input
}

func NewInputs() *Inputs {
	return &Inputs{
		inputs: make(map[string]*Input),
	}
}

func (inputs *Inputs) Add(input *Input) {
	if inputs.inputs == nil {
		inputs.inputs = make(map[string]*Input)
	}

	inputs.inputs[input.GetName()] = input
}

func (inputs Inputs) GetInput(name string) (*Input, error) {
	v, ok := inputs.inputs[name]

	if !ok {
		return nil, fmt.Errorf("No Input found for name: %s", name)
	}

	return v, nil
}
