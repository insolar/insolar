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

package dispatcher

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus/meta"
	busMeta "github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/flow"
	flowPulse "github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/insolar/payload"
	insPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/pulse"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/flow/dispatcher.Dispatcher -o ./ -s _mock.go -g
type Dispatcher interface {
	BeginPulse(ctx context.Context, pulse insolar.Pulse)
	ClosePulse(ctx context.Context, pulse insolar.Pulse)
	Process(msg *message.Message) error
}

type dispatcher struct {
	handles struct {
		present flow.MakeHandle
		future  flow.MakeHandle
		past    flow.MakeHandle
	}
	controller *thread.Controller
	pulses     insPulse.Accessor
}

func NewDispatcher(pulseAccessor insPulse.Accessor, present flow.MakeHandle, future flow.MakeHandle, past flow.MakeHandle) Dispatcher {
	d := &dispatcher{
		controller: thread.NewController(),
		pulses:     pulseAccessor,
	}

	d.handles.present = present
	d.handles.future = future
	d.handles.past = past

	return d
}

// BeginPulse is a handle for pulse begin event.
func (d *dispatcher) BeginPulse(ctx context.Context, pulseObject insolar.Pulse) {
	d.controller.BeginPulse()
	inslogger.FromContext(ctx).Debugf("Pulse was changed to %s in dispatcher", pulseObject.PulseNumber)
}

// ClosePulse is a handle for pulse close event.
func (d *dispatcher) ClosePulse(ctx context.Context, pulseObject insolar.Pulse) {
	d.controller.ClosePulse()
	inslogger.FromContext(ctx).Debugf("Pulse %s was closed in dispatcher", pulseObject.PulseNumber)
}

func (d *dispatcher) getHandleByPulse(ctx context.Context, msgPulseNumber insolar.PulseNumber) flow.MakeHandle {
	currentPulseNumber := insolar.PulseNumber(pulse.MinTimePulse)
	p, err := d.pulses.Latest(ctx)
	if err == nil {
		currentPulseNumber = p.PulseNumber
	} else {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to fetch pulse in dispatcher"))
	}

	if msgPulseNumber > currentPulseNumber {
		inslogger.FromContext(ctx).Debugf("Get message from future (pulse in msg %d, current pulse %d)", msgPulseNumber, currentPulseNumber)
		return d.handles.future
	}
	if msgPulseNumber < currentPulseNumber {
		return d.handles.past
	}
	return d.handles.present
}

// Process handles incoming message.
func (d *dispatcher) Process(msg *message.Message) error {
	processStart := time.Now()
	ctx := context.Background()
	ctx = inslogger.ContextWithTrace(ctx, msg.Metadata.Get(meta.TraceID))

	for k, v := range msg.Metadata {
		if k == meta.SpanData || k == meta.TraceID {
			continue
		}
		ctx, _ = inslogger.WithField(ctx, k, v)
	}
	logger := inslogger.FromContext(ctx)

	pn, err := insolar.NewPulseNumberFromStr(msg.Metadata.Get(meta.Pulse))
	if err != nil {
		logger.Error("failed to handle message: ", err)
		return nil
	}
	ctx = flowPulse.ContextWith(ctx, pn)
	parentSpan := instracer.MustDeserialize([]byte(msg.Metadata.Get(meta.SpanData)))
	ctx = instracer.WithParentSpan(ctx, parentSpan)

	msgType := messagePayloadTypeName(msg)

	go func() {
		<-d.controller.CanProcess()
		runStart := time.Now()

		f := thread.NewThread(msg, d.controller)
		handle := d.getHandleByPulse(ctx, pn)
		err := f.Run(ctx, handle(msg))

		runDuration := time.Since(runStart)
		procDuration := time.Since(processStart)
		result := "ok"
		if err != nil {
			if err == flow.ErrCancelled {
				result = "cancelled"
				logger.Info(errors.Wrap(err, "flow handling failed"))
			} else {
				result = "error"
				logger.Error(errors.Wrap(err, "flow handling failed"))
			}
		}

		ctx = insmetrics.ChangeTags(ctx,
			tag.Insert(tagMessageType, msgType),
			tag.Insert(tagResult, result),
		)
		stats.Record(ctx,
			statHandlerTime.M(float64(runDuration.Nanoseconds())/1e6),
			statProcessTime.M(float64(procDuration.Nanoseconds())/1e6))
	}()
	return nil
}

func pulseFromString(p string) (insolar.PulseNumber, error) {
	u64, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		return insolar.PulseNumber(0), errors.Wrap(err, "can't convert string value to pulse")
	}
	pInt := uint32(u64)
	return insolar.PulseNumber(pInt), nil
}

func messagePayloadTypeName(msg *message.Message) string {
	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		fmt.Println("meta decoding failed:", err)
		return "unknown"
	}
	payloadType, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		// branch for legacy messages format: INS-2973
		return msg.Metadata.Get(busMeta.Type)
	}
	return payloadType.String()
}
