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
)

// HistoryIterator is used to iterate object history.
//
// During iteration object refs will be fetched from remote source.
type HistoryIterator struct {
	ctx        context.Context
	messageBus core.MessageBus
	object     core.RecordRef
	chunkSize  int
	fromPulse  *core.PulseNumber
	fromPrev   *core.RecordID
	buff       []reply.Object
	buffIndex  int
	canFetch   bool
}

// History returns object's history references.
func (d *ObjectDescriptor) History(pulse *core.PulseNumber) (core.RefIterator, error) {
	return d.am.GetHistory(d.ctx, d.head, pulse)
}

// NewHistoryIterator creates new history iterator.
func NewHistoryIterator(ctx context.Context, mb core.MessageBus, object core.RecordRef, fromPulse *core.PulseNumber, chunkSize int) (*HistoryIterator, error) {
	iter := HistoryIterator{
		ctx:        ctx,
		messageBus: mb,
		object:     object,
		fromPulse:  fromPulse,
		chunkSize:  chunkSize,
		canFetch:   true,
	}
	err := iter.fetch()
	if err != nil {
		return nil, err
	}
	return &iter, nil
}

// HasNext checks if any elements left in iterator.
func (i *HistoryIterator) HasNext() bool {
	return i.hasInBuffer() || i.canFetch
}

// Next returns next element.
func (i *HistoryIterator) Next() (*core.RecordRef, error) {
	// Get element from buffer.
	if !i.hasInBuffer() && i.canFetch {
		err := i.fetch()
		if err != nil {
			return nil, err
		}
	}
	ref := i.nextFromBuffer()
	if ref == nil {
		return nil, errors.New("failed to fetch record")
	}

	return ref, nil
}

func (i *HistoryIterator) nextFromBuffer() *core.RecordRef {
	if !i.hasInBuffer() {
		return nil
	}
	ref := core.NewRecordRef(
		*i.object.Domain(),
		*i.buff[i.buffIndex].ChildPointer,
	)
	i.buffIndex++
	return ref
}

func (i *HistoryIterator) fetch() error {
	if !i.canFetch {
		return errors.New("failed to fetch record")
	}
	genericReply, err := i.messageBus.Send(
		i.ctx,
		&message.GetHistory{
			Object: i.object,
			Pulse:  i.fromPulse,
			From:   i.fromPrev,
			Amount: i.chunkSize,
		},
	)

	switch rep := genericReply.(type) {
	case *reply.ExplorerList:
		{
			if rep.NextFrom == nil {
				i.canFetch = false
			}
			i.buff = rep.Refs
			i.buffIndex = 0
			i.fromPrev = rep.NextFrom
		}
	case *reply.Error:
		err = rep.Error()
	default:
		err = ErrUnexpectedReply
	}
	return err
}

func (i *HistoryIterator) hasInBuffer() bool {
	return i.buffIndex < len(i.buff)
}
