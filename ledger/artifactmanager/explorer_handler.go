/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"context"
	"errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
)

func (h *MessageHandler) handleGetHistory(ctx context.Context, pulseNumber core.PulseNumber, genericMsg core.SignedMessage) (core.Reply, error) {
	msg := genericMsg.Message().(*message.GetHistory)

	idx, state, _, err := getObject(ctx, h.db, msg.Object.Record(), nil, false)
	if err != nil {
		return nil, err
	}
	history := []reply.Object{}
	var current *core.RecordID

	if msg.From != nil {
		current = msg.From
	} else {
		current = state
	}

	counter := 0
	for current != nil {
		// We have enough results.
		if counter >= msg.Amount {
			return &reply.ExplorerList{Refs: history, NextFrom: current}, nil
		}
		counter++

		rec, err := h.db.GetRecord(ctx, current)
		if err != nil {
			return nil, errors.New("failed to retrieve object state")
		}
		currentState, ok := rec.(record.ObjectState)

		if !ok {
			return nil, errors.New("Cannot cust to object state, type ")
		}
		current = currentState.PrevStateID()

		// Skip records later than specified pulse.
		// recPulse := current.Pulse()
		// if msg.Pulse != nil && recPulse > *msg.Pulse {
		// 	continue
		// }

		var memory []byte
		if currentState.GetMemory() != nil {
			memory, err = h.db.GetBlob(ctx, currentState.GetMemory())
			if err != nil {
				return nil, err
			}
		}

		history = append(history, reply.Object{
			Head:         msg.Object,
			Prototype:    currentState.GetImage(),
			IsPrototype:  currentState.GetIsPrototype(),
			ChildPointer: currentState.PrevStateID(),
			Parent:       idx.Parent,
			Memory:       memory,
		})
	}
	return &reply.ExplorerList{Refs: history, NextFrom: nil}, nil
}
