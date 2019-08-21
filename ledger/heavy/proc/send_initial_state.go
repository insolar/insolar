/*
 *    Copyright 2019 Insolar Technologies
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

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
)

type SendInitialState struct {
	meta payload.Meta

	dep struct {
		startPulse    pulse.StartPulse
		jetKeeper     executor.JetKeeper
		initialState  executor.InitialStateAccessor
		pulseAccessor pulse.Accessor
		sender        bus.Sender
	}
}

func (p *SendInitialState) Dep(
	startPulse pulse.StartPulse,
	jetKeeper executor.JetKeeper,
	initialState executor.InitialStateAccessor,
	pulseAccessor pulse.Accessor,
	sender bus.Sender,
) {
	p.dep.startPulse = startPulse
	p.dep.jetKeeper = jetKeeper
	p.dep.initialState = initialState
	p.dep.pulseAccessor = pulseAccessor
	p.dep.sender = sender
}

func NewSendInitialState(meta payload.Meta) *SendInitialState {
	return &SendInitialState{
		meta: meta,
	}
}

func (p *SendInitialState) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("SendInitialState starts working")

	startPulse, err := p.dep.startPulse.PulseNumber()

	if err != nil {
		errStr := "Couldn't get start pulse"
		msg, newErr := payload.NewMessage(&payload.Error{Text: errStr, Code: uint32(payload.CodeNoStartPulse)})
		if newErr != nil {
			logger.Fatal("failed to reply error", err)
		}
		p.dep.sender.Reply(ctx, p.meta, msg)
		return nil
	}
	logger = logger.WithField("startPulse", startPulse)

	msg, err := payload.Unmarshal(p.meta.Payload)
	if err != nil {
		logger.Fatal("Couldn't unmarshall request", err)
	}

	req, ok := msg.(*payload.GetLightInitialState)
	if !ok {
		return fmt.Errorf("unexpected payload type %T", msg)
	}

	topSyncPulseNumber := p.dep.jetKeeper.TopSyncPulse()
	topSyncPulse, err := p.dep.pulseAccessor.ForPulseNumber(ctx, topSyncPulseNumber)
	if err != nil {
		logger.Fatal("Couldn't get pulse for topSyncPulse: ", topSyncPulseNumber, " ", err)
	}

	switch {
	case req.Pulse == startPulse:
		p.sendForNetworkStart(ctx, req, topSyncPulse)
	case req.Pulse > startPulse:
		p.sendForJoiner(ctx, topSyncPulse)
	default:
		logger.Fatal("received initial state request from the past")
	}

	logger.Info("SendInitialState finishes working")
	return nil
}

func (p *SendInitialState) sendForNetworkStart(
	ctx context.Context,
	req *payload.GetLightInitialState,
	topSyncPulse insolar.Pulse,
) {
	logger := inslogger.FromContext(ctx)
	state := p.dep.initialState.Get(ctx, p.meta.Sender, req.Pulse)

	msg, err := payload.NewMessage(&payload.LightInitialState{
		NetworkStart: true,
		JetIDs:       state.JetIDs,
		Drops:        state.Drops,
		Indexes:      state.Indexes,
		Pulse:        *pulse.ToProto(&topSyncPulse),
	})
	if err != nil {
		logger.Fatal("Couldn't make message", err)
	}
	p.dep.sender.Reply(ctx, p.meta, msg)
}

func (p *SendInitialState) sendForJoiner(ctx context.Context, topSyncPulse insolar.Pulse) {
	logger := inslogger.FromContext(ctx)

	msg, err := payload.NewMessage(&payload.LightInitialState{
		NetworkStart: false,
		Pulse:        *pulse.ToProto(&topSyncPulse),
	})
	if err != nil {
		logger.Fatal("Couldn't make message", err)
	}
	p.dep.sender.Reply(ctx, p.meta, msg)
}
