package main

import (
	"fmt"
	. "github.com/jtejido/persephone"
)

// STATES
const (
	OPENED int = iota
	CLOSED
	CLOSED_AND_LOCKED
)

// INPUTS
const (
	OPEN int = iota
	CLOSE
	LOCK
	UNLOCK
)

// Sample DoorFSM, a classic one, persephone will be embedded here (To compensate for Golang's lack of inheritance).
type DoorFSM struct {
	*AbstractFSM
}

func NewDoorFSM(states *States, inputs *Inputs) *DoorFSM {
	// pass states and inputs to abstractFSM
	return &DoorFSM{
		AbstractFSM: New(states, inputs),
	}
}

func (fsm *DoorFSM) UnlockAction() error {
	fmt.Println("unlocking action.")
	return nil
}

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
	states := NewStates()
	states.Add(&State{"opened", OPENED, INITIAL_STATE})
	states.Add(&State{"closed", CLOSED, NORMAL_STATE})
	states.Add(&State{"closedAndLocked", CLOSED_AND_LOCKED, NORMAL_STATE})

	// initialize accepted inputs, either do the things below inside the NewDoorFSM, init() or outside.
	inputs := NewInputs()
	inputs.Add(&Input{"open", OPEN})
	inputs.Add(&Input{"close", CLOSE})
	inputs.Add(&Input{"lock", LOCK})
	inputs.Add(&Input{"unlock", UNLOCK})

	fsm := NewDoorFSM(states, inputs)

	// Either do the things below inside the NewDoorFSM, init() or outside.
	// Call convention pattern is utilized since AbstractFSM is embedded.
	// Add transition rules [entryState, input, targetState, transitionCallback]
	fsm.AddRule("opened", "close", "closed", nil)
	fsm.AddRule("closed", "open", "opened", nil)
	fsm.AddRule("closed", "lock", "closedAndLocked", nil)
	fsm.AddRule("closedAndLocked", "unlock", "closed", fsm.UnlockAction)

	// Callback on an action, say, actions that affects recognized lexemes
	fsm.AddInputAction("closedAndLocked", "unlock", fsm.UnlockAction)
	fsm.AddTransitionAction("closed", "opened", fsm.OpenAction)

	// Entry and Exit actions
	fsm.AddEntryAction("closed", fsm.CloseEntryAction)
	fsm.AddExitAction("closed", fsm.CloseExitAction)

	fsm.Process("close")
	fmt.Printf("%s \n\n", fsm.GetState().GetName())

	fsm.Process("lock")
	fmt.Printf("%s \n\n", fsm.GetState().GetName())

	fsm.Process("unlock")
	fmt.Printf("%s \n\n", fsm.GetState().GetName())

	fsm.Process("open")
	fmt.Printf("%s \n\n", fsm.GetState().GetName())

	// just to check if it accepts invalid transitions
	err := fsm.Process("lock")
	fmt.Printf("%s \n", err.Error())                // should stop here now
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // still opened, the door's last state.

	// just to check if it accepts further valid transitions after previous invalid
	fsm.Process("close")
	fmt.Printf("%s \n\n", fsm.GetState().GetName())

	// can't do this too when it's just closed, door should be locked first
	err_a := fsm.Process("unlock")
	fmt.Printf("%s \n", err_a.Error()) // should stop here now
	fmt.Printf("%s \n\n", fsm.GetState().GetName()) // still closed

	// Can() checks if you can go to a certain state without tripping the alarm (the callbacks)
	fmt.Println(fsm.Can("open")) // you can open
	fmt.Println(fsm.Can("lock")) // or lock
	fmt.Println(fsm.Can("close")) // but not close
	fmt.Println(fsm.Can("unlock")) // or unlock
}
