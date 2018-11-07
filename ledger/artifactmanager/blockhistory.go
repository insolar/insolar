package artifactmanager

import (
	"context"
	"errors"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/record"
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
	fromChild  *core.RecordID
	buff       []record.ObjectState
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
	ref := i.buff[i.buffIndex].GetImage()
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
		},
	)
	if err != nil {
		return err
	}
	rep, ok := genericReply.(*reply.History)
	if !ok {
		return errors.New("failed to fetch record")
	}
	if rep.NextFrom == nil {
		i.canFetch = false
	}
	i.buff = rep.Refs
	i.buffIndex = 0
	i.fromChild = rep.NextFrom
	return nil
}

func (i *HistoryIterator) hasInBuffer() bool {
	return i.buffIndex < len(i.buff)
}
