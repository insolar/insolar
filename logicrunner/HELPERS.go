//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package logicrunner

import (
	"context"
	"strings"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type StackingHandle interface {
	SetNext(next StackingHandle)
	GetPresent() flow.Handle
}

func BuildHandles(procedure flow.Procedure, subHandles ...StackingHandle) flow.Handle {
	var termHandle StackingHandle
	termHandle = &HandleAsProcedure{
		procedure: procedure,
	}
	subHandlesLen := len(subHandles)
	currentHandle := termHandle
	for i := range subHandles {
		nextHandle := subHandles[subHandlesLen-i-1]
		nextHandle.SetNext(currentHandle)
		currentHandle = nextHandle
	}

	return currentHandle.GetPresent()
}

type HandleIncorrectMessagePulse struct {
	numRepeats int
	next       flow.Handle
}

func (h *HandleIncorrectMessagePulse) SetNext(next StackingHandle) {
	h.next = next.GetPresent()
}

func (h *HandleIncorrectMessagePulse) GetPresent() flow.Handle {
	return h.Present
}

func (h *HandleIncorrectMessagePulse) Present(ctx context.Context, f flow.Flow) error {
	err := f.Handle(ctx, h.next)
	inslogger.FromContext(ctx).Infof(">>>>>>>>>> 1")
	if err != nil {
		inslogger.FromContext(ctx).Infof(">>>>>>>>>> 2 , %s. Left Repeats: %d", err, h.numRepeats)
		if strings.Contains(err.Error(), "Incorrect message pulse") || err == flow.ErrCancelled {
			h.numRepeats--
			inslogger.FromContext(ctx).Infof(">>>>>>>>>> 3 , %s. Left Repeats: %d", err, h.numRepeats)
			if h.numRepeats > 0 {
				inslogger.FromContext(ctx).Infof("[ HandleIncorrectMessagePulse ] Got '%s'. Try one more time.", err.Error())
				return f.Migrate(ctx, h.Present)
			}

		}
	}
	return err
}

type HandleAsProcedure struct {
	procedure  flow.Procedure
	cancelable bool
}

func (h *HandleAsProcedure) SetNext(next StackingHandle) {
}

func (h *HandleAsProcedure) GetPresent() flow.Handle {
	return h.Present
}

func (h *HandleAsProcedure) Present(ctx context.Context, f flow.Flow) error {
	return f.Procedure(ctx, h.procedure, h.cancelable)
}
