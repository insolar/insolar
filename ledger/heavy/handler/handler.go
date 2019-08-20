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
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"

	"github.com/insolar/insolar/ledger/heavy/proc"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
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
	BackupMaker   executor.BackupMaker

	Sender          bus.Sender
	StartPulse      pulse.StartPulse
	PulseCalculator pulse.Calculator
	JetTree         jet.Storage
	DropDB          *drop.DB

	Replicator executor.HeavyReplicator

	dep *proc.Dependencies
}

// New creates a new handler.
func New(cfg configuration.Ledger) *Handler {
	h := &Handler{
		cfg: cfg,
	}
	dep := proc.Dependencies{
		PassState: func(p *proc.PassState) {
			p.Dep.Records = h.RecordAccessor
			p.Dep.Sender = h.Sender
			p.Dep.Pulses = h.PulseAccessor
		},
		SendCode: func(p *proc.SendCode) {
			p.Dep.Sender = h.Sender
			p.Dep.RecordAccessor = h.RecordAccessor
		},
		SendRequests: func(p *proc.SendRequests) {
			p.Dep(h.Sender, h.RecordAccessor, h.IndexAccessor)
		},
		SendRequest: func(p *proc.SendRequest) {
			p.Dep(h.RecordAccessor, h.Sender)
		},
		Replication: func(p *proc.Replication) {
			p.Dep(
				h.Replicator,
			)
		},
		SendJet: func(p *proc.SendJet) {
			p.Dep(
				h.JetAccessor,
				h.Sender)
		},
		SendIndex: func(p *proc.SendIndex) {
			p.Dep(
				h.IndexAccessor,
				h.Sender,
			)
		},
		SendInitialState: func(p *proc.SendInitialState) {
			p.Dep(
				h.StartPulse,
				h.JetKeeper,
				h.JetTree,
				h.JetCoordinator,
				h.DropDB,
				h.PulseAccessor,
				h.Sender,
			)
		},
	}
	h.dep = &dep
	return h
}

func (h *Handler) Process(msg *watermillMsg.Message) error {
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

	return nil
}

func (h *Handler) handle(ctx context.Context, msg *watermillMsg.Message) error {
	var err error
	logger := inslogger.FromContext(ctx)

	meta := payload.Meta{}
	err = meta.Unmarshal(msg.Payload)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "failed to unmarshal meta")
	}
	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "failed to unmarshal payload type")
	}
	ctx, _ = inslogger.WithField(ctx, "msg_type", payloadType.String())

	ctx, span := instracer.StartSpan(ctx, payloadType.String())
	defer span.End()

	switch payloadType {
	case payload.TypeGetRequest:
		p := proc.NewSendRequest(meta)
		h.dep.SendRequest(p)
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
		p := proc.NewSendCode(meta)
		h.dep.SendCode(p)
		err = p.Proceed(ctx)
	case payload.TypeReplication:
		p := proc.NewReplication(meta, h.cfg)
		h.dep.Replication(p)
		err = p.Proceed(ctx)
	case payload.TypeGetJet:
		p := proc.NewSendJet(meta)
		h.dep.SendJet(p)
		err = p.Proceed(ctx)
	case payload.TypeGetIndex:
		p := proc.NewSendIndex(meta)
		h.dep.SendIndex(p)
		err = p.Proceed(ctx)
	case payload.TypePass:
		err = h.handlePass(ctx, meta)
	case payload.TypeError:
		h.handleError(ctx, meta)
	case payload.TypeGotHotConfirmation:
		h.handleGotHotConfirmation(ctx, meta)
	case payload.TypeGetLightInitialState:
		p := proc.NewSendInitialState(meta)
		h.dep.SendInitialState(p)
		err = p.Proceed(ctx)
	default:
		err = fmt.Errorf("no handler for message type %s", payloadType.String())
	}
	if err != nil {
		instracer.AddError(span, err)
		h.replyError(ctx, meta, err)
	}
	return err
}

func (h *Handler) handleError(ctx context.Context, msg payload.Meta) {
	logger := inslogger.FromContext(ctx)

	pl := payload.Error{}
	err := pl.Unmarshal(msg.Payload)
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to unmarshal error"))
		return
	}

	logger.Error("received error: ", pl.Text)
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
		p := proc.NewSendCode(originMeta)
		h.dep.SendCode(p)
		err = p.Proceed(ctx)
	case payload.TypeGetRequest:
		p := proc.NewSendRequest(originMeta)
		h.dep.SendRequest(p)
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

func (h *Handler) handleGotHotConfirmation(ctx context.Context, meta payload.Meta) {
	logger := inslogger.FromContext(ctx)
	logger.Info("handleGotHotConfirmation got new message")

	confirm := payload.GotHotConfirmation{}
	err := confirm.Unmarshal(meta.Payload)
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to unmarshal to GotHotConfirmation"))
		return
	}

	logger.Info("handleGotHotConfirmation. pulse: ", confirm.Pulse, ". jet: ", confirm.JetID.DebugString(), ". Split: ", confirm.Split)

	err = h.JetKeeper.AddHotConfirmation(ctx, confirm.Pulse, confirm.JetID, confirm.Split)
	if err != nil {
		logger.Fatalf("failed to add hot confirmation jet=%v: %v", confirm.String(), err.Error())
	}

	executor.FinalizePulse(ctx, h.PulseCalculator, h.BackupMaker, h.JetKeeper, h.IndexModifier, confirm.Pulse)
	logger.Info("handleGotHotConfirmation finish. pulse: ", confirm.Pulse, ". jet: ", confirm.JetID.DebugString(), ". Split: ", confirm.Split)
}
