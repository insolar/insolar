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
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/replica"

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
	Bus            insolar.MessageBus
	JetCoordinator jet.Coordinator
	PCS            insolar.PlatformCryptographyScheme
	BlobAccessor   blob.Accessor
	BlobModifier   blob.Modifier
	RecordAccessor object.RecordAccessor
	RecordModifier object.RecordModifier

	IndexAccessor object.IndexAccessor
	IndexModifier object.IndexModifier

	DropModifier  drop.Modifier
	PulseAccessor pulse.Accessor
	JetModifier   jet.Modifier
	JetAccessor   jet.Accessor
	JetKeeper     replica.JetKeeper
	Sender        bus.Sender

	jetID insolar.JetID
	dep   *proc.Dependencies
}

// New creates a new handler.
func New() *Handler {
	h := &Handler{
		jetID: insolar.ZeroJetID,
	}
	dep := proc.Dependencies{
		PassState: func(p *proc.PassState) {
			p.Dep.Blobs = h.BlobAccessor
			p.Dep.Records = h.RecordAccessor
			p.Dep.Sender = h.Sender
		},
		GetCode: func(p *proc.GetCode) {
			p.Dep.Sender = h.Sender
			p.Dep.RecordAccessor = h.RecordAccessor
			p.Dep.BlobAccessor = h.BlobAccessor
		},
		SendRequests: func(p *proc.SendRequests) {
			p.Dep(h.Sender, h.RecordAccessor, h.IndexAccessor)
		},
	}
	h.dep = &dep
	return h
}

func (h *Handler) Process(msg *watermillMsg.Message) ([]*watermillMsg.Message, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(bus.MetaTraceID))

	for k, v := range msg.Metadata {
		ctx, _ = inslogger.WithField(ctx, k, v)
	}
	logger := inslogger.FromContext(ctx)

	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}

	err = h.handle(ctx, msg)
	if err != nil {
		logger.Error(errors.Wrap(err, "handle error"))
	}

	return nil, nil
}

func (h *Handler) handle(ctx context.Context, msg *watermillMsg.Message) error {
	var err error

	meta := payload.Meta{}
	err = meta.Unmarshal(msg.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}
	ctx, _ = inslogger.WithField(ctx, "msg_type", payloadType.String())

	switch payloadType {
	case payload.TypeGetFilament:
		p := proc.NewSendRequests(meta)
		h.dep.SendRequests(p)
		return p.Proceed(ctx)
	case payload.TypePassState:
		p := proc.NewPassState(meta)
		h.dep.PassState(p)
		err = p.Proceed(ctx)
	case payload.TypeGetCode:
		p := proc.NewGetCode(meta)
		h.dep.GetCode(p)
		err = p.Proceed(ctx)
	case payload.TypePass:
		err = h.handlePass(ctx, meta)
	case payload.TypeError:
		h.handleError(ctx, meta)
	default:
		err = fmt.Errorf("no handler for message type %s", payloadType.String())
	}
	if err != nil {
		h.replyError(ctx, meta, err)
	}
	return err
}

func (h *Handler) handleError(ctx context.Context, msg payload.Meta) {
	pl := payload.Error{}
	err := pl.Unmarshal(msg.Payload)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to unmarshal error"))
		return
	}

	inslogger.FromContext(ctx).Error("received error: ", pl.Text)
}

func (h *Handler) handlePass(ctx context.Context, meta payload.Meta) error {
	pass := payload.Pass{}
	err := pass.Unmarshal(meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal pass payload")
	}

	originMeta := payload.Meta{}
	err = originMeta.Unmarshal(pass.Origin)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal origin message")
	}
	payloadType, err := payload.UnmarshalType(originMeta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload type")
	}

	ctx, _ = inslogger.WithField(ctx, "msg_type_original", payloadType.String())

	switch payloadType { // nolint
	case payload.TypeGetCode:
		p := proc.NewGetCode(originMeta)
		h.dep.GetCode(p)
		err = p.Proceed(ctx)
	default:
		err = fmt.Errorf("no pass handler for message type %s", payload.Type(originMeta.Polymorph).String())
	}
	if err != nil {
		h.replyError(ctx, originMeta, err)
	}
	return err
}

func (h *Handler) replyError(ctx context.Context, replyTo payload.Meta, err error) {
	errMsg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to reply error"))
	}
	go h.Sender.Reply(ctx, replyTo, errMsg)
}

func (h *Handler) Init(ctx context.Context) error {
	h.Bus.MustRegister(insolar.TypeHeavyPayload, h.handleHeavyPayload)

	h.Bus.MustRegister(insolar.TypeGetDelegate, h.handleGetDelegate)
	h.Bus.MustRegister(insolar.TypeGetChildren, h.handleGetChildren)
	h.Bus.MustRegister(insolar.TypeGetObjectIndex, h.handleGetObjectIndex)
	h.Bus.MustRegister(insolar.TypeGetJet, h.handleGetJet)
	h.Bus.MustRegister(insolar.TypeGetRequest, h.handleGetRequest)
	return nil
}

func (h *Handler) handleGetDelegate(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetDelegate)

	idx, err := h.IndexAccessor.ForID(ctx, parcel.Pulse(), *msg.Head.Record())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch index for %v", msg.Head.Record()))
	}

	delegateRef, ok := idx.Lifeline.DelegateByKey(msg.AsType)
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

	idx, err := h.IndexAccessor.ForID(ctx, parcel.Pulse(), *msg.Parent.Record())
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
		currentChild = idx.Lifeline.ChildPointer
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

func (h *Handler) handleGetJet(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetJet)
	jet, actual := h.JetAccessor.ForID(ctx, msg.Pulse, msg.Object)

	return &reply.Jet{ID: insolar.ID(jet), Actual: actual}, nil
}

func (h *Handler) handleGetObjectIndex(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObjectIndex)

	idx, err := h.IndexAccessor.ForID(ctx, parcel.Pulse(), *msg.Object.Record())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch object index for %v", msg.Object.Record().String())
	}

	buf := object.EncodeLifeline(idx.Lifeline)

	return &reply.ObjectIndex{Index: buf}, nil
}

func (h *Handler) handleHeavyPayload(ctx context.Context, genericMsg insolar.Parcel) (insolar.Reply, error) {
	logger := inslogger.FromContext(ctx)
	msg := genericMsg.Message().(*message.HeavyPayload)

	storeRecords(ctx, h.RecordModifier, h.PCS, msg.PulseNum, msg.Records)
	if err := storeIndexBuckets(ctx, h.IndexModifier, msg.IndexBuckets, msg.PulseNum); err != nil {
		return &reply.HeavyError{Message: err.Error(), JetID: msg.JetID, PulseNum: msg.PulseNum}, nil
	}

	pulse, err := h.PulseAccessor.Latest(ctx)
	if err != nil {
		return &reply.HeavyError{Message: err.Error(), JetID: msg.JetID, PulseNum: msg.PulseNum}, nil
	}
	futurePulse := pulse.NextPulseNumber
	drop, err := storeDrop(ctx, h.DropModifier, msg.Drop)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to store drop"))
		return &reply.HeavyError{Message: err.Error(), JetID: msg.JetID, PulseNum: msg.PulseNum}, nil
	}
	if drop.Split {
		_, _, err = h.JetModifier.Split(ctx, futurePulse, drop.JetID)

	} else {
		err = h.JetModifier.Update(ctx, futurePulse, false, drop.JetID)
	}
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to split/update jet=%v pulse=%v", drop.JetID.DebugString(), futurePulse))
		return &reply.HeavyError{Message: err.Error(), JetID: msg.JetID, PulseNum: msg.PulseNum}, nil
	}

	if err := h.JetKeeper.Add(ctx, drop.Pulse, drop.JetID); err != nil {
		logger.Error(errors.Wrapf(err, "failed to add jet to JetKeeper jet=%v", msg.JetID.DebugString()))
		return &reply.HeavyError{Message: err.Error(), JetID: msg.JetID, PulseNum: msg.PulseNum}, nil
	}

	storeBlobs(ctx, h.BlobModifier, h.PCS, msg.PulseNum, msg.Blobs)

	stats.Record(ctx,
		statReceivedHeavyPayloadCount.M(1),
	)

	return &reply.OK{}, nil
}
