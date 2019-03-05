// +build with_generated

package main

import (
	"github.com/insolar/insolar/conveyor/generator/state_machines/sample"
	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"testing"
)

func Test_Generated_State_Machine(t *testing.T) {
	element := slot.NewSlotElementHelperMock(t)
	element.GetInputEventFunc = func() interface{} {
		return sample.Event{}
	}
	element.GetPayloadFunc = func() interface{} {
		return &sample.Payload{}
	}

	machine := sample.SMRHTestStateMachineFactory()

	var state uint32 = 0
	for {
		_, state, _ = machine.GetTransitionHandler(constant.Present, state)(element)
		if state == 0 {
			break
		}
	}
}

