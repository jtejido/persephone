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
	/*
	 * The names act as unique identifier, the first added state is the machine's INITIAL_STATE.
	 */
	 
	// initialize states and inputs, either do the things below, either inside the NewDoorFSM func, an init() func or outside.
	
	states := persephone.NewStates()
	states.Add("opened", OPENED)
	states.Add("closed", CLOSED)
	states.Add("closedAndLocked", CLOSED_AND_LOCKED)

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
	fsm.AddRule("closedAndLocked", "unlock", "closed", fsm.UnlockAction) // fourth is trivial; guarding, error firing and stuff.
	
	// Callback on an action, say, actions that affects recognized input when a certain state is triggered.
	fsm.AddInputAction("closedAndLocked", "unlock", fsm.UnlockAction)
	
	// Callback when a transition from src to dest state occured.
	fsm.AddTransitionAction("closed", "opened", fsm.OpenAction)

	// Entry and Exit actions
	fsm.AddEntryAction("closed", fsm.CloseEntryAction)
	fsm.AddExitAction("closed", fsm.CloseExitAction)

	// sample process for the door
	fsm.Process("close")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // closed

	fsm.Process("lock")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // closedAndLocked

	fsm.Process("unlock")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // closed

	fsm.Process("open")
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // opened
	
	err:= fsm.Process("lock")
	fmt.Printf("%s \n", err.Error()) // invalid input, can't lock the door when it's in opened state, this returns error.
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // opened, no changes in state.
}
```


## Example

Go to the sample DoorFSM for a classic FSM example.

More examples to come.
