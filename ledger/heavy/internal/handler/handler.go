package handler

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
)

type Handler struct {
	Bus            insolar.MessageBus     `inject:""`
	JetCoordinator insolar.JetCoordinator `inject:""`
	HeavySync      insolar.HeavySync      `inject:""`

	// TODO: @imarkin 27.03.2019 - remove it after all new storages integration (INS-2013, etc)
	ObjectStorage storage.ObjectStorage `inject:""`

	object.RecordAccessor `inject:""`

	jetID insolar.JetID
}

// NewMessageHandler creates new handler.
func New() *Handler {
	return &Handler{
		jetID: *insolar.NewJetID(0, nil),
	}
}

func (h *Handler) Init(ctx context.Context) error {
	h.Bus.MustRegister(insolar.TypeHeavyStartStop, h.handleHeavyStartStop)
	h.Bus.MustRegister(insolar.TypeHeavyPayload, h.handleHeavyPayload)

	h.Bus.MustRegister(insolar.TypeGetCode, h.handleGetCode)
	h.Bus.MustRegister(insolar.TypeGetObject, h.handleGetObject)
	h.Bus.MustRegister(insolar.TypeGetDelegate, h.handleGetDelegate)
	h.Bus.MustRegister(insolar.TypeGetChildren, h.handleGetChildren)
	h.Bus.MustRegister(insolar.TypeGetObjectIndex, h.handleGetObjectIndex)
	h.Bus.MustRegister(insolar.TypeGetRequest, h.handleGetRequest)
	return nil
}

func (h *Handler) handleGetCode(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetCode)
	jetID := *insolar.NewJetID(0, nil)

	codeRec, err := h.getCode(ctx, msg.Code.Record())
	if err != nil {
		return nil, err
	}
	code, err := h.ObjectStorage.GetBlob(ctx, insolar.ID(jetID), codeRec.Code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch code blob")
	}

	rep := reply.Code{
		Code:        code,
		MachineType: codeRec.MachineType,
	}

	return &rep, nil
}

func (h *Handler) handleGetObject(
	ctx context.Context, parcel insolar.Parcel,
) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObject)

	// Fetch object index. If not found redirect.
	idx, err := h.ObjectStorage.GetObjectIndex(ctx, insolar.ID(h.jetID), msg.Head.Record())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch object index for %s", msg.Head.Record().DebugString())
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

	// Fetch state record.
	rec, err := h.RecordAccessor.ForID(ctx, *stateID)
	virtRec := rec.Record
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch state %s for %s", stateID.DebugString(), msg.Head.Record()))
	}
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
		rep.Memory, err = h.ObjectStorage.GetBlob(ctx, insolar.ID(h.jetID), state.GetMemory())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch blob")
		}
	}

	return &rep, nil
}

func (h *Handler) handleGetDelegate(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetDelegate)

	idx, err := h.ObjectStorage.GetObjectIndex(ctx, insolar.ID(h.jetID), msg.Head.Record())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch index for %v", msg.Head.Record()))
	}

	delegateRef, ok := idx.Delegates[msg.AsType]
	if !ok {
		return nil, errors.New("the object has no delegate for this type")
	}
	rep := reply.Delegate{
		Head: delegateRef,
	}

	return &rep, nil
}

func (h *Handler) handleGetChildren(
	ctx context.Context, parcel insolar.Parcel,
) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetChildren)

	idx, err := h.ObjectStorage.GetObjectIndex(ctx, insolar.ID(h.jetID), msg.Parent.Record())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch index for %v", msg.Parent.Record()))
	}

	var (
		refs         []insolar.Reference
		currentChild *insolar.ID
	)

	// Counting from specified child or the latest.
	if msg.FromChild != nil {
		currentChild = msg.FromChild
	} else {
		currentChild = idx.ChildPointer
	}

	// The object has no children.
	if currentChild == nil {
		return &reply.Children{Refs: nil, NextFrom: nil}, nil
	}

	// Try to fetch the first child.
	_, err = h.RecordAccessor.ForID(ctx, *currentChild)
	if err == object.ErrNotFound {
		text := fmt.Sprintf(
			"failed to fetch child %s for %s",
			currentChild.DebugString(),
			msg.Parent.Record().DebugString(),
		)
		return nil, errors.Wrap(err, text)
	}

	counter := 0
	for currentChild != nil {
		// We have enough results.
		if counter >= msg.Amount {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		counter++

		rec, err := h.RecordAccessor.ForID(ctx, *currentChild)
		virtRec := rec.Record
		// We don't have this child reference. Return what was collected.
		if err == object.ErrNotFound {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		if err != nil {
			return nil, errors.New("failed to retrieve children")
		}

		childRec, ok := virtRec.(*object.ChildRecord)
		if !ok {
			return nil, errors.New("failed to retrieve children")
		}
		currentChild = childRec.PrevChild

		// Skip records later than specified pulse.
		recPulse := childRec.Ref.Record().Pulse()
		if msg.FromPulse != nil && recPulse > *msg.FromPulse {
			continue
		}
		refs = append(refs, childRec.Ref)
	}

	return &reply.Children{Refs: refs, NextFrom: nil}, nil
}

func (h *Handler) handleGetRequest(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetRequest)

	rec, err := h.ObjectStorage.GetRecord(ctx, insolar.ID(h.jetID), &msg.Request)
	if err != nil {
		return nil, errors.New("failed to fetch request")
	}

	req, ok := rec.(*object.RequestRecord)
	if !ok {
		return nil, errors.New("failed to decode request")
	}

	rep := reply.Request{
		ID:     msg.Request,
		Record: object.SerializeRecord(req),
	}

	return &rep, nil
}

func (h *Handler) handleGetObjectIndex(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObjectIndex)

	idx, err := h.ObjectStorage.GetObjectIndex(ctx, insolar.ID(h.jetID), msg.Object.Record())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}

	buf := object.EncodeIndex(*idx)

	return &reply.ObjectIndex{Index: buf}, nil
}

func (h *Handler) getCode(ctx context.Context, id *insolar.ID) (*object.CodeRecord, error) {
	jetID := *insolar.NewJetID(0, nil)

	rec, err := h.ObjectStorage.GetRecord(ctx, insolar.ID(jetID), id)
	if err != nil {
		return nil, err
	}
	codeRec, ok := rec.(*object.CodeRecord)
	if !ok {
		return nil, errors.New("failed to retrieve code record")
	}

	return codeRec, nil
}

// TODO: check sender if it was light material in synced pulses:
// sender := genericMsg.GetSender()
// sender.isItWasLMInPulse(pulsenum)
func (h *Handler) handleHeavyPayload(ctx context.Context, genericMsg insolar.Parcel) (insolar.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyPayload)

	if err := h.HeavySync.Store(ctx, insolar.ID(msg.JetID), msg.PulseNum, msg.Records); err != nil {
		return heavyerrreply(err)
	}
	if err := h.HeavySync.StoreDrop(ctx, msg.JetID, msg.Drop); err != nil {
		return heavyerrreply(err)
	}

	return &reply.OK{}, nil
}

func (h *Handler) handleHeavyStartStop(ctx context.Context, genericMsg insolar.Parcel) (insolar.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyStartStop)

	// stop
	if msg.Finished {
		if err := h.HeavySync.Stop(ctx, insolar.ID(msg.JetID), msg.PulseNum); err != nil {
			return nil, err
		}
		return &reply.OK{}, nil
	}
	// start
	if err := h.HeavySync.Start(ctx, insolar.ID(msg.JetID), msg.PulseNum); err != nil {
		return heavyerrreply(err)
	}
	return &reply.OK{}, nil
}

func heavyerrreply(err error) (insolar.Reply, error) {
	if herr, ok := err.(*reply.HeavyError); ok {
		return herr, nil
	}
	return nil, err
}
