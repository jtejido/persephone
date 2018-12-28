package persephone

type FSMMethod func() error

type FSMAction struct {
	method FSMMethod
}

func NewFSMAction(method FSMMethod) *FSMAction {
	return &FSMAction{
		method: method,
	}
}

func (fsma *FSMAction) DoAction() error {
	return fsma.method()
}
