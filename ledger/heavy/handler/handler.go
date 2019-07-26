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
	"bytes"
	"context"
	"fmt"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"

	"github.com/insolar/insolar/ledger/heavy/proc"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
)

// Handler is a base struct for heavy's methods
type Handler struct {
	cfg configuration.Ledger

	Bus            insolar.MessageBus
	JetCoordinator jet.Coordinator
	PCS            insolar.PlatformCryptographyScheme
	RecordAccessor object.RecordAccessor
	RecordModifier object.RecordModifier

	IndexAccessor object.IndexAccessor
	IndexModifier object.IndexModifier

	DropModifier  drop.Modifier
	PulseAccessor pulse.Accessor
	JetModifier   jet.Modifier
	JetAccessor   jet.Accessor
	JetKeeper     executor.JetKeeper

	Sender bus.Sender

	jetID insolar.JetID
	dep   *proc.Dependencies
}

// New creates a new handler.
func New(cfg configuration.Ledger) *Handler {
	h := &Handler{
		cfg:   cfg,
		jetID: insolar.ZeroJetID,
	}
	dep := proc.Dependencies{
		PassState: func(p *proc.PassState) {
			p.Dep.Records = h.RecordAccessor
			p.Dep.Sender = h.Sender
		},
		GetCode: func(p *proc.GetCode) {
			p.Dep.Sender = h.Sender
			p.Dep.RecordAccessor = h.RecordAccessor
		},
		SendRequests: func(p *proc.SendRequests) {
			p.Dep(h.Sender, h.RecordAccessor, h.IndexAccessor)
		},
		GetRequest: func(p *proc.GetRequest) {
			p.Dep(h.RecordAccessor, h.Sender)
		},
		Replication: func(p *proc.Replication) {
			p.Dep(
				h.RecordModifier,
				h.IndexModifier,
				h.PCS,
				h.PulseAccessor,
				h.DropModifier,
				h.JetModifier,
				h.JetKeeper,
			)
		},
		GetJet: func(p *proc.GetJet) {
			p.Dep(
				h.JetAccessor,
				h.Sender)
		},
	}
	h.dep = &dep
	return h
}

func (h *Handler) Process(msg *watermillMsg.Message) ([]*watermillMsg.Message, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(bus.MetaTraceID))
	parentSpan, err := instracer.Deserialize([]byte(msg.Metadata.Get(bus.MetaSpanData)))
	if err == nil {
		ctx = instracer.WithParentSpan(ctx, parentSpan)
	} else {
		inslogger.FromContext(ctx).Error(err)
	}

	for k, v := range msg.Metadata {
		if k == bus.MetaSpanData || k == bus.MetaTraceID {
			continue
		}
		ctx, _ = inslogger.WithField(ctx, k, v)
	}
	logger := inslogger.FromContext(ctx)

	meta := payload.Meta{}
	err = meta.Unmarshal(msg.Payload)
	if err != nil {
		logger.Error(err)
	}

	err = h.handle(ctx, msg)
	if err != nil {
		logger.Error(errors.Wrap(err, "handle error"))
	}

	return nil, nil
}

func (h *Handler) handleParcel(ctx context.Context, msg *watermillMsg.Message) error {
	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal meta")
	}

	parcel, err := message.DeserializeParcel(bytes.NewBuffer(meta.Payload))
	if err != nil {
		return errors.Wrap(err, "can't deserialize payload to parcel")
	}

	msgType := msg.Metadata.Get(bus.MetaType)
	ctx, _ = inslogger.WithField(ctx, "msg_type", msgType)
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("Present %v", parcel.Message().Type().String()))
	defer span.End()

	var rep insolar.Reply
	switch msgType {
	case insolar.TypeGetChildren.String():
		rep, err = h.handleGetChildren(ctx, parcel)
	case insolar.TypeGetDelegate.String():
		rep, err = h.handleGetDelegate(ctx, parcel)
	case insolar.TypeGetObjectIndex.String():
		rep, err = h.handleGetObjectIndex(ctx, parcel)
	default:
		err = fmt.Errorf("no handler for message type %s", msgType)
	}
	if err != nil {
		h.replyError(ctx, meta, errors.Wrap(err, "error while handle parcel"))
	} else {
		resAsMsg := bus.ReplyAsMessage(ctx, rep)
		h.Sender.Reply(ctx, meta, resAsMsg)
	}
	return err
}

func (h *Handler) handle(ctx context.Context, msg *watermillMsg.Message) error {
	msgType := msg.Metadata.Get(bus.MetaType)
	if msgType != "" {
		return h.handleParcel(ctx, msg)
	}

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
	case payload.TypeGetRequest:
		p := proc.NewGetRequest(meta)
		h.dep.GetRequest(p)
		err = p.Proceed(ctx)
	case payload.TypeGetFilament:
		p := proc.NewSendRequests(meta)
		h.dep.SendRequests(p)
		err = p.Proceed(ctx)
	case payload.TypePassState:
		p := proc.NewPassState(meta)
		h.dep.PassState(p)
		err = p.Proceed(ctx)
	case payload.TypeGetCode:
		p := proc.NewGetCode(meta)
		h.dep.GetCode(p)
		err = p.Proceed(ctx)
	case payload.TypeGetJet:
		p := proc.NewGetJet(meta)
		h.dep.GetJet(p)
		err = p.Proceed(ctx)
	case payload.TypePass:
		err = h.handlePass(ctx, meta)
	case payload.TypeError:
		h.handleError(ctx, meta)
	case payload.TypeGotHotConfirmation:
		h.handleGotHotConfirmation(ctx, meta)
	case payload.TypeReplication:
		p := proc.NewReplication(meta, h.cfg)
		h.dep.Replication(p)
		err = p.Proceed(ctx)
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
	case payload.TypeGetRequest:
		p := proc.NewGetRequest(originMeta)
		h.dep.GetRequest(p)
		err = p.Proceed(ctx)
	default:
		err = fmt.Errorf("no pass handler for message type %s", payloadType.String())
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

func (h *Handler) handleGetObjectIndex(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := parcel.Message().(*message.GetObjectIndex)

	idx, err := h.IndexAccessor.ForID(ctx, parcel.Pulse(), *msg.Object.Record())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch object index for %v", msg.Object.Record().String())
	}

	buf := object.EncodeLifeline(idx.Lifeline)

	return &reply.ObjectIndex{Index: buf}, nil
}

func (h *Handler) handleGotHotConfirmation(ctx context.Context, meta payload.Meta) {
	logger := inslogger.FromContext(ctx)
	confirm := payload.GotHotConfirmation{}
	err := confirm.Unmarshal(meta.Payload)
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to unmarshal to GotHotConfirmation"))
		return
	}

	logger.Debug("handleGotHotConfirmation. pulse: ", confirm.Pulse, ". jet: ", confirm.JetID.DebugString())

	err = h.JetModifier.Update(ctx, confirm.Pulse, true, confirm.JetID)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to update jet %s", confirm.JetID.DebugString()))
		return
	}

	err = h.JetKeeper.AddHotConfirmation(ctx, confirm.Pulse, confirm.JetID, confirm.Split)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to add hot confitmation to JetKeeper jet=%v", confirm.String()))
	} else {
		logger.Debug("got confirmation: ", confirm.String())
	}
}
