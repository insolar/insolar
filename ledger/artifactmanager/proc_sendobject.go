package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/pkg/errors"
)

type SendObject struct {
	Handler *MessageHandler

	Message bus.Message
	Jet     insolar.JetID
	Index   object.Lifeline
}

func (p *SendObject) Proceed(ctx context.Context) error {
	ctx = contextWithJet(ctx, insolar.ID(p.Jet))
	r := bus.Reply{}
	r.Reply, r.Err = p.handle(ctx, p.Message.Parcel)
	p.Message.ReplyTo <- r
	return nil
}

func (p *SendObject) handle(
	ctx context.Context, parcel insolar.Parcel,
) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObject)
	logger := inslogger.FromContext(ctx)

	var stateID *insolar.ID
	if msg.State != nil {
		stateID = msg.State
	} else {
		if msg.Approved {
			stateID = p.Index.LatestStateApproved
		} else {
			stateID = p.Index.LatestState
		}
	}
	if stateID == nil {
		return &reply.Error{ErrType: reply.ErrStateNotAvailable}, nil
	}

	var (
		stateJet *insolar.ID
	)
	onHeavy, err := p.Handler.JetCoordinator.IsBeyondLimit(ctx, parcel.Pulse(), stateID.Pulse())
	if err != nil && err != pulse.ErrNotFound {
		return nil, err
	}
	if onHeavy {
		hNode, err := p.Handler.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		logger.WithFields(map[string]interface{}{
			"state":    stateID.DebugString(),
			"going_to": hNode.String(),
		}).Debug("fetching object (on heavy)")

		obj, err := p.Handler.fetchObject(ctx, msg.Head, *hNode, stateID, parcel.Pulse())
		if err != nil {
			if err == insolar.ErrDeactivated {
				return &reply.Error{ErrType: reply.ErrDeactivated}, nil
			}
			return nil, err
		}

		return &reply.Object{
			Head:         msg.Head,
			State:        *stateID,
			Prototype:    obj.Prototype,
			IsPrototype:  obj.IsPrototype,
			ChildPointer: p.Index.ChildPointer,
			Parent:       p.Index.Parent,
			Memory:       obj.Memory,
		}, nil
	}

	stateJetID, actual := p.Handler.JetStorage.ForID(ctx, stateID.Pulse(), *msg.Head.Record())
	stateJet = (*insolar.ID)(&stateJetID)

	if !actual {
		actualJet, err := p.Handler.jetTreeUpdater.fetchJet(ctx, *msg.Head.Record(), stateID.Pulse())
		if err != nil {
			return nil, err
		}
		stateJet = actualJet
	}

	// Fetch state record.
	rec, err := p.Handler.RecordAccessor.ForID(ctx, *stateID)

	if err == object.ErrNotFound {
		// The record wasn't found on the current suitNode. Return redirect to the node that contains it.
		// We get Jet tree for pulse when given state was added.
		suitNode, err := p.Handler.JetCoordinator.NodeForJet(ctx, *stateJet, parcel.Pulse(), stateID.Pulse())
		if err != nil {
			return nil, err
		}
		logger.WithFields(map[string]interface{}{
			"state":    stateID.DebugString(),
			"going_to": suitNode.String(),
		}).Debug("fetching object (record not found)")

		obj, err := p.Handler.fetchObject(ctx, msg.Head, *suitNode, stateID, parcel.Pulse())
		if err != nil {
			if err == insolar.ErrDeactivated {
				return &reply.Error{ErrType: reply.ErrDeactivated}, nil
			}
			return nil, err
		}

		return &reply.Object{
			Head:         msg.Head,
			State:        *stateID,
			Prototype:    obj.Prototype,
			IsPrototype:  obj.IsPrototype,
			ChildPointer: p.Index.ChildPointer,
			Parent:       p.Index.Parent,
			Memory:       obj.Memory,
		}, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "can't fetch record from storage")
	}

	virtRec := rec.Record
	state, ok := virtRec.(object.State)
	if !ok {
		return nil, errors.New("invalid object record")
	}

	if state.ID() == object.StateDeactivation {
		return &reply.Error{ErrType: reply.ErrDeactivated}, nil
	}

	var childPointer *insolar.ID
	if p.Index.ChildPointer != nil {
		childPointer = p.Index.ChildPointer
	}
	rep := reply.Object{
		Head:         msg.Head,
		State:        *stateID,
		Prototype:    state.GetImage(),
		IsPrototype:  state.GetIsPrototype(),
		ChildPointer: childPointer,
		Parent:       p.Index.Parent,
	}

	if state.GetMemory() != nil {
		b, err := p.Handler.BlobAccessor.ForID(ctx, *state.GetMemory())
		if err == blob.ErrNotFound {
			hNode, err := p.Handler.JetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				return nil, err
			}
			obj, err := p.Handler.fetchObject(ctx, msg.Head, *hNode, stateID, parcel.Pulse())
			if err != nil {
				return nil, err
			}
			err = p.Handler.BlobModifier.Set(ctx, *state.GetMemory(), blob.Blob{
				JetID: insolar.JetID(p.Jet),
				Value: obj.Memory},
			)
			if err != nil {
				return nil, err
			}
			b.Value = obj.Memory
		}
		rep.Memory = b.Value
	}

	return &rep, nil
}
