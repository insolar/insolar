package getobject

import (
	"context"

	"github.com/insolar/insolar/conveyor/fsm"
)

type transitionFunc func(context.Context, fsm.SlotElementHelper) fsm.ElementState
type adapterResponseFunc func(context.Context, fsm.SlotElementHelper, interface{}) fsm.ElementState

type slotElement struct {
	transitions      []transitionFunc
	adapterResponses []adapterResponseFunc
}

func (se *slotElement) Transition(state fsm.ElementState, f transitionFunc, _ ...fsm.ElementState) *slotElement {
	for i := len(se.transitions); i <= int(state); i++ {
		se.transitions[i] = nil
	}
	se.transitions[int(state)] = f
	return se
}

func (se *slotElement) AdapterResponse(state fsm.ElementState, f adapterResponseFunc, _ ...fsm.ElementState) *slotElement {
	for i := len(se.transitions); i <= int(state); i++ {
		se.adapterResponses[i] = nil
	}
	se.adapterResponses[int(state)] = f
	return se
}

func (se *slotElement) GetTransition(state fsm.ElementState) transitionFunc {
	return se.transitions[int(state)]
}
