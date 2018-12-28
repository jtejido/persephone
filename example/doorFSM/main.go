package main

import (
	"fmt"
	"github.com/jtejido/persephone"
)

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

// Sample DoorFSM, a classic one, persephone will be embedded here.
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

	fsm := NewDoorFSM(states, inputs)

	// Either do the things below inside the NewDoorFSM, init() or outside.
	// add transition rules [entryState, input, targetState, transitionCallback]
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

	// can't do this too when it's just closed, door should be locked
	err_a := fsm.Process("unlock")
	fmt.Printf("%s \n", err_a.Error()) // should stop here now
	fmt.Printf("%s \n", fsm.GetState().GetName())
}
