# persephone

Persephone is an Abstract Finite State Machine/Transducer, which means, the primary usage is pretty much just embedding it to any form of an FSM struct/implementation (mealy or moore machine) one is creating.


```
const (
	OPENED int = iota
	CLOSED
	CLOSED_AND_LOCKED
)

const (
	OPEN int = iota
	CLOSE
	LOCK
	UNLOCK
)

type DoorFSM struct {
	*persephone.AbstractFSM
}

func NewDoorFSM(states *persephone.States, inputs *persephone.Inputs) *DoorFSM {
	// pass states and inputs to abstractFSM
	return &DoorFSM{
		AbstractFSM: persephone.New(states, inputs),
	}
}

func (fsm *DoorFSM) UnlockAction() error {
	fmt.Println("unlocking the door.")
	return nil
}

func main() {
	// initialize states, either do the things below inside the NewDoorFSM, init() or outside.
	states := persephone.NewStates()
	states.Add("opened", OPENED)
	states.Add("closed", CLOSED)
	states.Add("closedAndLocked", CLOSED_AND_LOCKED)

	// initialize accepted inputs, either do the things below inside the NewDoorFSM, init() or outside.
	inputs := persephone.NewInputs()
	inputs.Add("open", OPEN)
	inputs.Add("close", CLOSE)
	inputs.Add("lock", LOCK)
	inputs.Add("unlock", UNLOCK)

	// start the door
	fsm := NewDoorFSM(states, inputs)

	// add transition rules
	fsm.AddRule("opened", "close", "closed", nil)
	fsm.AddRule("closed", "open", "opened", nil)
	fsm.AddRule("closed", "lock", "closedAndLocked", nil)
	fsm.AddRule("closedAndLocked", "unlock", "closed", fsm.UnlockAction)

	// sample process for the door
	fsm.Process("close")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // closed

	fsm.Process("lock")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // closedAndLocked

	fsm.Process("unlock")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // closed

	fsm.Process("open")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // opened
}
```


## Example

Go to the sample DoorFSM for a classic FSM example.
More examples to come.