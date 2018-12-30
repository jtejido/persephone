package persephone

// What a func should implement, having error in return means that I can stop at any given point during processing.
type FSMMethod func() error

type Callback interface {
	Do() error
}

type FSMAction struct {
	method FSMMethod
}

func NewFSMAction(method FSMMethod) *FSMAction {
	return &FSMAction{
		method: method,
	}
}

func (fsma *FSMAction) Do() error {
	return fsma.method()
}
