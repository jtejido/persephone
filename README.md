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

// Sample DoorFSM, a classic one. AbstractFSM will be embedded here (to compensate for Golang's lack of inheritance).
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
	 
	// initialize states and inputs. You can do the things below, either inside the NewDoorFSM func, an init() func or outside.

	states := NewStates()
	states.Add(&persephone.State{"opened", OPENED, INITIAL_STATE})
	states.Add(&persephone.State{"closed", CLOSED, NORMAL_STATE})
	states.Add(&persephone.State{"closedAndLocked", CLOSED_AND_LOCKED, NORMAL_STATE})

	inputs := NewInputs()
	inputs.Add(&persephone.Input{"open", OPEN})
	inputs.Add(&persephone.Input{"close", CLOSE})
	inputs.Add(&persephone.Input{"lock", LOCK})
	inputs.Add(&persephone.Input{"unlock", UNLOCK})

	// start the door
	fsm := NewDoorFSM(states, inputs)

	// Call convention pattern is utilized as AbstractFSM is already embedded.
	// add transition rules
	fsm.AddRule("opened", "close", "closed", nil)
	fsm.AddRule("closed", "open", "opened", nil)
	fsm.AddRule("closed", "lock", "closedAndLocked", nil)
	fsm.AddRule("closedAndLocked", "unlock", "closed", fsm.UnlockAction) // fourth is trivial; guarding, syntax error firing and stuff.
	
	// Callback on an action, say, actions that affects recognized input when a certain state is triggered.
	fsm.AddInputAction("closedAndLocked", "unlock", fsm.UnlockAction)
	
	// Callback when a transition from src to dest state occured.
	fsm.AddTransitionAction("closed", "opened", fsm.OpenAction)

	// Entry and Exit actions, methods to run before and after a state.
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
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // opened, no changes in state, proceed as usual.
}
```


## Example

Go to the sample DoorFSM for a classic FSM example.

More examples to come.

