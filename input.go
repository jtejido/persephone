package persephone

type Input uint

type Inputs struct {
	inputs bitSet
}

func (inputs *Inputs) Add(input Input) {
	inputs.inputs.add(uint(input))
}

func (inputs *Inputs) Contains(input Input) bool {
	return inputs.inputs.contains(uint(input))
}
