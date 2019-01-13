# persephone

Persephone is an Abstract Finite State Machine/Transducer, which means, the primary usage is pretty much just embedding it to any form of an FSM struct/implementation (mealy or moore machine) one is creating. For efficiency, it uses bitset for states and inputs.


```
// STATES
const (
	OPENED State = iota
	CLOSED
	CLOSED_AND_LOCKED
)

// INPUTS
const (
	OPEN Input = iota
	CLOSE
	LOCK
	UNLOCK
)

// Sample DoorFSM, a classic one. AbstractFSM will be embedded here (to compensate for Golang's lack of inheritance).
type DoorFSM struct {
	*persephone.AbstractFSM
}

func NewDoorFSM(states States, inputs Inputs) *DoorFSM {
	// pass states and inputs to abstractFSM
	return &DoorFSM{
		AbstractFSM: New(states, inputs),
	}
}

func (fsm *DoorFSM) UnlockAction() error {
	fmt.Println("unlocking action.")
	return nil
}

// Callback on an action, say, actions that affects recognized lexemes
func (fsm *DoorFSM) InputAction() error {
	fmt.Println("input detected.")
	return nil
}

// sample transition callback when transitioning from closed to opened
func (fsm *DoorFSM) OpenAction() error {
	fmt.Println("open action.")
	return nil
}

func (fsm *DoorFSM) CloseEntryAction() error {
	fmt.Println("close entry action.")
	return nil
}

func (fsm *DoorFSM) CloseExitAction() error {
	fmt.Println("close exit action.")
	return nil
}

func main() {
	// initialize states, either do the things below inside the NewDoorFSM, init() or outside.
	// when initial state isn't set, the first added state defaults as the initial state. Please be cautious with state orders when no 
	// INITIAL_STATE is specified.
	var states States
	states.Add(OPENED, INITIAL_STATE)
	states.Add(CLOSED, NORMAL_STATE)
	states.Add(CLOSED_AND_LOCKED, NORMAL_STATE)

	// initialize accepted inputs, either do the things below inside the NewDoorFSM, init() or outside.
	var inputs Inputs
	inputs.Add(OPEN)
	inputs.Add(CLOSE)
	inputs.Add(LOCK)
	inputs.Add(UNLOCK)

	fsm := NewDoorFSM(states, inputs)

	// Either do these inside our NewDoorFSM, an init() or outside, just like below.
	// Call convention pattern is utilized since AbstractFSM is embedded.
	// Add transition rules [entryState, input, targetState, transitionCallback]
	fsm.AddRule(OPENED, CLOSE, CLOSED, nil)
	fsm.AddRule(CLOSED, OPEN, OPENED, nil)
	fsm.AddRule(CLOSED, LOCK, CLOSED_AND_LOCKED, nil)
	fsm.AddRule(CLOSED_AND_LOCKED, UNLOCK, CLOSED, fsm.UnlockAction)

	fsm.AddInputAction(CLOSED_AND_LOCKED, UNLOCK, fsm.UnlockAction)
	fsm.AddTransitionAction(CLOSED, OPENED, fsm.OpenAction)
	fsm.AddEntryAction(CLOSED, fsm.CloseEntryAction)
	fsm.AddExitAction(CLOSED, fsm.CloseExitAction)

	// Now let's process inputs
	fsm.Process(CLOSE)
	fmt.Printf("%d \n\n", fsm.GetState())

	fsm.Process(LOCK)
	fmt.Printf("%d \n\n", fsm.GetState())

	fsm.Process(UNLOCK)
	fmt.Printf("%d \n\n", fsm.GetState())

	fsm.Process(OPEN)
	fmt.Printf("%d \n\n", fsm.GetState())

	// to check if it accepts invalid transitions
	err := fsm.Process(LOCK)
	fmt.Printf("%s \n", err.Error())                // you could stop here now
	fmt.Printf("%d \n\n", fsm.GetState()) // still opened, the door's last state.

	// to check if it accepts further valid transitions after previous invalid
	fsm.Process(CLOSE)
	fmt.Printf("%d \n\n", fsm.GetState())

	// can't do this too when it's just closed, door should be locked first
	err_a := fsm.Process(UNLOCK)
	fmt.Printf("%s \n", err_a.Error()) // you could stop here now
	fmt.Printf("%d \n\n", fsm.GetState()) // still closed

	// Can() checks if you can go to a certain state without tripping the callbacks
	fmt.Println(fsm.Can(OPEN)) // you can open
	fmt.Println(fsm.Can(LOCK)) // or lock
	fmt.Println(fsm.Can(CLOSE)) // but not close
	fmt.Println(fsm.Can(UNLOCK)) // or unlock
}
```


## Example

Go to the sample DoorFSM for a classic FSM example.

More examples to come.

