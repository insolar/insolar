package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/pkg/errors"
)

type ReturnReply struct {
	Message bus.Message
	Err     error
	Reply   insolar.Reply
}

func (p *ReturnReply) Proceed(context.Context) {
	p.Message.ReplyTo <- bus.Reply{Reply: p.Reply, Err: p.Err}
}

type FetchJet struct {
	Message bus.Message
	handler *MessageHandler

	Res struct {
		JetID insolar.JetID
		Miss  bool
		Err   error
	}
}

func (p *FetchJet) Proceed(ctx context.Context) {
	jet, err := p.jet(ctx)
	if jet != nil {
		p.Res.JetID = *jet
	}
	p.Res.Err = err
}

func (p *FetchJet) jet(ctx context.Context) (*insolar.JetID, error) {
	parcel := p.Message.Parcel
	msg := parcel.Message()
	if msg.DefaultTarget() == nil {
		return nil, errors.New("unexpected message")
	}

	// Hack to temporary allow any genesis request.
	if p.Message.Parcel.Pulse() <= insolar.FirstPulseNumber {
		return insolar.NewJetID(0, nil), nil
	}

	// Check token jet.
	token := parcel.DelegationToken()
	if token != nil {
		// Calculate jet for target pulse.
		target := *msg.DefaultTarget().Record()
		pulse := target.Pulse()
		switch tm := msg.(type) {
		case *message.GetObject:
			pulse = tm.State.Pulse()
		case *message.GetChildren:
			if tm.FromChild == nil {
				return nil, errors.New("fetching children without child pointer is forbidden")
			}
			pulse = tm.FromChild.Pulse()
		case *message.GetRequest:
			pulse = tm.Request.Pulse()
		}
		jetID, actual := p.handler.JetStorage.ForID(ctx, pulse, target)
		if !actual {
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"msg":   msg.Type().String(),
				"jet":   jetID.DebugString(),
				"pulse": pulse,
			}).Error("jet is not actual")
		}

		return &jetID, nil
	}

	// Calculate jet for current pulse.
	var jetID insolar.ID
	if msg.DefaultTarget().Record().Pulse() == insolar.PulseNumberJet {
		jetID = *msg.DefaultTarget().Record()
	} else {
		j, err := p.handler.jetTreeUpdater.fetchJet(ctx, *msg.DefaultTarget().Record(), parcel.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch jet tree")
		}

		jetID = *j
	}

	// Check if jet is ours.
	node, err := p.handler.JetCoordinator.LightExecutorForJet(ctx, jetID, parcel.Pulse())
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate executor for jet")
	}

	if *node != p.handler.JetCoordinator.Me() {
		p.Res.Miss = true
		j := insolar.JetID(jetID)
		return &j, nil
	}

	ctx = addJetIDToLogger(ctx, jetID)

	j := insolar.JetID(jetID)
	return &j, nil
}

type WaitHot struct {
	Message bus.Message
	JetID   insolar.JetID

	handler *MessageHandler

	Res struct {
		timeout bool
	}
}

func (p *WaitHot) Proceed(ctx context.Context) {
	parcel := p.Message.Parcel
	// Hack is needed for genesis:
	// because we don't have hot data on first pulse and without this we would stale.
	if parcel.Pulse() <= insolar.FirstPulseNumber {
		return
	}

	// If the call is a call in redirect-chain
	// skip waiting for the hot records
	if parcel.DelegationToken() != nil {
		return
	}

	err := p.handler.HotDataWaiter.Wait(ctx, insolar.ID(p.JetID))
	if err != nil {
		p.Res.timeout = true
		return
	}
}

type ProcGetObject struct {
	Message bus.Message
	JetID   insolar.JetID
	Handler *MessageHandler
}

func (p *ProcGetObject) Proceed(ctx context.Context) {
	ctx = contextWithJet(ctx, insolar.ID(p.JetID))
	r := bus.Reply{}
	r.Reply, r.Err = p.handle(ctx, p.Message.Parcel)
	p.Message.ReplyTo <- r
}

func (p *ProcGetObject) handle(
	ctx context.Context, parcel insolar.Parcel,
) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObject)
	jetID := jetFromContext(ctx)
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object": msg.Head.Record().DebugString(),
		"pulse":  parcel.Pulse(),
	})

	p.Handler.RecentStorageProvider.GetIndexStorage(ctx, jetID).AddObject(ctx, *msg.Head.Record())

	p.Handler.IDLocker.Lock(msg.Head.Record())
	defer p.Handler.IDLocker.Unlock(msg.Head.Record())

	// Fetch object index. If not found redirect.
	idx, err := p.Handler.ObjectStorage.GetObjectIndex(ctx, jetID, msg.Head.Record())
	if err == insolar.ErrNotFound {
		logger.Debug("failed to fetch index (fetching from heavy)")
		node, err := p.Handler.JetCoordinator.Heavy(ctx, parcel.Pulse())
		if err != nil {
			return nil, err
		}
		idx, err = p.Handler.saveIndexFromHeavy(ctx, jetID, msg.Head, node)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch index from heavy")
		}
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch object index %s", msg.Head.Record().String())
	}

	// Determine object state id.
	var stateID *insolar.ID
	if msg.State != nil {
		stateID = msg.State
	} else {
		if msg.Approved {
			stateID = idx.LatestStateApproved
		} else {
			stateID = idx.LatestState
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
			ChildPointer: idx.ChildPointer,
			Parent:       idx.Parent,
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
			ChildPointer: idx.ChildPointer,
			Parent:       idx.Parent,
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
	if idx.ChildPointer != nil {
		childPointer = idx.ChildPointer
	}
	rep := reply.Object{
		Head:         msg.Head,
		State:        *stateID,
		Prototype:    state.GetImage(),
		IsPrototype:  state.GetIsPrototype(),
		ChildPointer: childPointer,
		Parent:       idx.Parent,
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
				JetID: insolar.JetID(jetID),
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
