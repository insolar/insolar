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

package handler

import (
	"context"
	"fmt"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/proc"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
)

// Handler is a base struct for heavy's methods
type Handler struct {
	WmBus                 bus.Bus
	Bus                   insolar.MessageBus
	JetCoordinator        jet.Coordinator
	PCS                   insolar.PlatformCryptographyScheme
	BlobAccessor          blob.Accessor
	BlobModifier          blob.Modifier
	RecordAccessor        object.RecordAccessor
	RecordModifier        object.RecordModifier
	IndexLifelineAccessor object.LifelineAccessor
	IndexBucketModifier   object.IndexBucketModifier
	DropModifier          drop.Modifier
	Sender                bus.Sender

	jetID insolar.JetID
	dep   *proc.Dependencies
}

// New creates a new handler.
func New() *Handler {
	h := &Handler{
		jetID: *insolar.NewJetID(0, nil),
	}
	dep := proc.Dependencies{
		PassState: func(p *proc.PassState) {
			p.Dep.Blobs = h.BlobAccessor
			p.Dep.Records = h.RecordAccessor
			p.Dep.Sender = h.Sender
			p.Dep.Coordinator = h.JetCoordinator
			p.Dep.PCS = h.PCS
		},
	}
	h.dep = &dep
	return h

}

func (h *Handler) Process(msg *watermillMsg.Message) ([]*watermillMsg.Message, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(bus.MetaTraceID))

	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
	ctx, logger := inslogger.WithField(ctx, "pulse", fmt.Sprintf("%d", meta.Pulse))

	err = h.handle(ctx, msg)
	if err != nil {
		logger.Error(errors.Wrap(err, "handle error"))
		errMsg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
		if err != nil {
			logger.Error(errors.Wrap(err, "failed to reply error"))
			return nil, nil
		}
		go h.Sender.Reply(ctx, meta, errMsg)
	}

	return nil, nil
}

func (h *Handler) handle(ctx context.Context, msg *watermillMsg.Message) error {
	pl, err := payload.UnmarshalFromMeta(msg.Payload)
	if err != nil {
		return errors.Wrap(err, "can't deserialize meta payload")
	}
	switch pl.(type) {
	case *payload.PassState:
		p := proc.NewPassState(msg)
		h.dep.PassState(p)
		return p.Proceed(ctx)
	default:
		return fmt.Errorf("no handler for message type #%T", pl)
	}
}

func (h *Handler) Init(ctx context.Context) error {
	h.Bus.MustRegister(insolar.TypeHeavyPayload, h.handleHeavyPayload)

	h.Bus.MustRegister(insolar.TypeGetCode, h.handleGetCode)
	h.Bus.MustRegister(insolar.TypeGetDelegate, h.handleGetDelegate)
	h.Bus.MustRegister(insolar.TypeGetChildren, h.handleGetChildren)
	h.Bus.MustRegister(insolar.TypeGetObjectIndex, h.handleGetObjectIndex)
	h.Bus.MustRegister(insolar.TypeGetRequest, h.handleGetRequest)
	return nil
}

func (h *Handler) handleGetCode(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetCode)

	codeRec, err := h.getCode(ctx, msg.Code.Record())
	if err != nil {
		return nil, err
	}

	code, err := h.BlobAccessor.ForID(ctx, codeRec.Code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch code blob")
	}

	rep := reply.Code{
		Code:        code.Value,
		MachineType: codeRec.MachineType,
	}

	return &rep, nil
}

func (h *Handler) handleGetObject(
	ctx context.Context, parcel insolar.Parcel,
) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObject)

	// Fetch object index. If not found redirect.
	idx, err := h.IndexLifelineAccessor.ForID(ctx, parcel.Pulse(), *msg.Head.Record())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch object index for %s", msg.Head.Record().DebugString())
	}

	// Determine object state id.
	var stateID *insolar.ID
	if msg.State != nil {
		stateID = msg.State
	} else {
		stateID = idx.LatestState
	}
	if stateID == nil {
		return &reply.Error{ErrType: reply.ErrStateNotAvailable}, nil
	}

	// Fetch state record.
	rec, err := h.RecordAccessor.ForID(ctx, *stateID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch state %s for %s", stateID.DebugString(), msg.Head.Record()))
	}

	virtRec := rec.Virtual
	concrete := record.Unwrap(virtRec)
	state, ok := concrete.(record.State)
	if !ok {
		return nil, errors.New("invalid object record")
	}
	if state.ID() == record.StateDeactivation {
		return &reply.Error{ErrType: reply.ErrDeactivated}, nil
	}

	rep := reply.Object{
		Head:         msg.Head,
		State:        *stateID,
		Prototype:    state.GetImage(),
		IsPrototype:  state.GetIsPrototype(),
		ChildPointer: idx.ChildPointer,
		Parent:       idx.Parent,
	}

	if state.GetMemory() != nil && state.GetMemory().NotEmpty() {
		b, err := h.BlobAccessor.ForID(ctx, *state.GetMemory())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch blob")
		}
		rep.Memory = b.Value
	}

	return &rep, nil
}

func (h *Handler) handleGetDelegate(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetDelegate)

	idx, err := h.IndexLifelineAccessor.ForID(ctx, parcel.Pulse(), *msg.Head.Record())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch index for %v", msg.Head.Record()))
	}

	delegateRef, ok := idx.DelegateByKey(msg.AsType)
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

	idx, err := h.IndexLifelineAccessor.ForID(ctx, parcel.Pulse(), *msg.Parent.Record())
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

		// We don't have this child reference. Return what was collected.
		if err == object.ErrNotFound {
			return &reply.Children{Refs: refs, NextFrom: currentChild}, nil
		}
		if err != nil {
			return nil, errors.New("failed to retrieve children")
		}

		virtRec := rec.Virtual
		concrete := record.Unwrap(virtRec)
		childRec, ok := concrete.(*record.Child)
		if !ok {
			return nil, errors.New("failed to retrieve children")
		}

		currentChild = &childRec.PrevChild

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

	rec, err := h.RecordAccessor.ForID(ctx, msg.Request)
	if err != nil {
		return nil, errors.New("failed to fetch request")
	}

	virtRec := rec.Virtual
	concrete := record.Unwrap(virtRec)
	_, ok := concrete.(*record.Request)
	if !ok {
		return nil, errors.New("failed to decode request")
	}

	data, err := virtRec.Marshal()
	if err != nil {
		return nil, errors.New("failed to serialize request")
	}

	rep := reply.Request{
		ID:     msg.Request,
		Record: data,
	}

	return &rep, nil
}

func (h *Handler) handleGetObjectIndex(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObjectIndex)

	idx, err := h.IndexLifelineAccessor.ForID(ctx, parcel.Pulse(), *msg.Object.Record())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch object index")
	}

	buf := object.EncodeIndex(idx)

	return &reply.ObjectIndex{Index: buf}, nil
}

func (h *Handler) getCode(ctx context.Context, id *insolar.ID) (record.Code, error) {
	rec, err := h.RecordAccessor.ForID(ctx, *id)
	if err != nil {
		return record.Code{}, errors.Wrap(err, "can't get record from storage")
	}

	virtRec := rec.Virtual
	concrete := record.Unwrap(virtRec)
	codeRec, ok := concrete.(*record.Code)
	if !ok {
		return record.Code{}, errors.New("failed to retrieve code record")
	}

	return *codeRec, nil
}

func (h *Handler) handleHeavyPayload(ctx context.Context, genericMsg insolar.Parcel) (insolar.Reply, error) {
	msg := genericMsg.Message().(*message.HeavyPayload)

	storeRecords(ctx, h.RecordModifier, h.PCS, msg.PulseNum, msg.Records)
	if err := storeIndexBuckets(ctx, h.IndexBucketModifier, msg.IndexBuckets, msg.PulseNum); err != nil {
		return &reply.HeavyError{Message: err.Error(), JetID: msg.JetID, PulseNum: msg.PulseNum}, nil
	}
	if err := storeDrop(ctx, h.DropModifier, msg.Drop); err != nil {
		return &reply.HeavyError{Message: err.Error(), JetID: msg.JetID, PulseNum: msg.PulseNum}, nil
	}
	storeBlobs(ctx, h.BlobModifier, h.PCS, msg.PulseNum, msg.Blobs)

	stats.Record(ctx,
		statReceivedHeavyPayloadCount.M(1),
	)

	return &reply.OK{}, nil
}
