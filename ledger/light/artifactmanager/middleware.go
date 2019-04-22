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

package artifactmanager

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/hot"
)

type middleware struct {
	jetAccessor    jet.Accessor
	jetCoordinator jet.Coordinator
	messageBus     insolar.MessageBus
	jetReleaser    hot.JetReleaser
	jetWaiter      hot.JetWaiter
	conf           *configuration.Ledger
	handler        *MessageHandler
}

func newMiddleware(
	h *MessageHandler,
) *middleware {
	return &middleware{
		jetAccessor:    h.JetStorage,
		jetCoordinator: h.JetCoordinator,
		messageBus:     h.Bus,
		jetReleaser:    h.JetReleaser,
		jetWaiter:      h.HotDataWaiter,
		handler:        h,
		conf:           h.conf,
	}
}

func (m *middleware) addFieldsToLogger(handler insolar.MessageHandler) insolar.MessageHandler {
	return func(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
		ctx, _ = inslogger.WithField(ctx, "targetid", parcel.DefaultTarget().String())

		return handler(ctx, parcel)
	}
}

type jetKey struct{}

func contextWithJet(ctx context.Context, jetID insolar.ID) context.Context {
	return context.WithValue(ctx, jetKey{}, jetID)
}

func jetFromContext(ctx context.Context) insolar.ID {
	val := ctx.Value(jetKey{})
	j, ok := val.(insolar.ID)
	if !ok {
		panic("failed to extract jet from context")
	}

	return j
}

func (m *middleware) zeroJetForHeavy(handler insolar.MessageHandler) insolar.MessageHandler {
	return func(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
		return handler(contextWithJet(ctx, insolar.ID(*insolar.NewJetID(0, nil))), parcel)
	}
}

func addJetIDToLogger(ctx context.Context, jetID insolar.ID) context.Context {
	ctx, _ = inslogger.WithField(ctx, "jetid", jetID.DebugString())

	return ctx
}

func (m *middleware) checkJet(handler insolar.MessageHandler) insolar.MessageHandler {
	return func(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
		msg := parcel.Message()
		if msg.DefaultTarget() == nil {
			return nil, errors.New("unexpected message")
		}

		// Hack to temporary allow any genesis request.
		if parcel.Pulse() <= insolar.FirstPulseNumber {
			return handler(contextWithJet(ctx, insolar.ID(*insolar.NewJetID(0, nil))), parcel)
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
			jetID, actual := m.jetAccessor.ForID(ctx, pulse, target)
			if !actual {
				inslogger.FromContext(ctx).WithFields(map[string]interface{}{
					"msg":   msg.Type().String(),
					"jet":   jetID.DebugString(),
					"pulse": pulse,
				}).Error("jet is not actual")
			}

			return handler(contextWithJet(ctx, insolar.ID(jetID)), parcel)
		}

		// Calculate jet and pulse.
		var jetID *insolar.ID
		var pulse insolar.PulseNumber
		if msg.DefaultTarget().Record().Pulse() == insolar.PulseNumberJet {
			jetID = msg.DefaultTarget().Record()
		} else {
			if gr, ok := msg.(*message.GetRequest); ok {
				pulse = gr.Request.Pulse()
			} else {
				pulse = parcel.Pulse()
			}

			var err error
			jetID, err = m.handler.jetTreeUpdater.Fetch(ctx, *msg.DefaultTarget().Record(), pulse)
			if err != nil {
				return nil, errors.Wrap(err, "failed to fetch jet tree")
			}
		}

		// Check if jet is ours.
		node, err := m.jetCoordinator.LightExecutorForJet(ctx, *jetID, pulse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate executor for jet")
		}
		if *node != m.jetCoordinator.Me() {
			return &reply.JetMiss{JetID: *jetID, Pulse: pulse}, nil
		}

		ctx = addJetIDToLogger(ctx, *jetID)

		return handler(contextWithJet(ctx, *jetID), parcel)
	}
}

func (m *middleware) waitForHotData(handler insolar.MessageHandler) insolar.MessageHandler {
	return func(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
		// Hack is needed for genesis:
		// because we don't have hot data on first pulse and without this we would stale.
		if parcel.Pulse() <= insolar.FirstPulseNumber {
			return handler(ctx, parcel)
		}

		// If the call is a call in redirect-chain
		// skip waiting for the hot records
		if parcel.DelegationToken() != nil {
			return handler(ctx, parcel)
		}

		jetID := jetFromContext(ctx)
		err := m.jetWaiter.Wait(ctx, jetID)
		if err != nil {
			return &reply.Error{ErrType: reply.ErrHotDataTimeout}, nil
		}
		return handler(ctx, parcel)
	}
}

func (m *middleware) releaseHotDataWaiters(handler insolar.MessageHandler) insolar.MessageHandler {
	return func(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
		rep, err := handler(ctx, parcel)

		hotDataMessage := parcel.Message().(*message.HotData)
		jetID := hotDataMessage.Jet.Record()
		unlockErr := m.jetReleaser.Unlock(ctx, *jetID)
		if unlockErr != nil {
			inslogger.FromContext(ctx).Error(err)
		}

		return rep, err
	}
}
