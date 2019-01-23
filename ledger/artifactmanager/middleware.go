/*
 *    Copyright 2019 Insolar
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

package artifactmanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
)

type middleware struct {
	db                                 *storage.DB
	jetCoordinator                     core.JetCoordinator
	messageBus                         core.MessageBus
	pulseStorage                       core.PulseStorage
	earlyRequestCircuitBreakerProvider *earlyRequestCircuitBreakerProvider
	conf                               *configuration.Ledger
	handler                            *MessageHandler
	seqMutex                           sync.Mutex
	sequencer                          map[core.RecordID]*struct {
		sync.Mutex
		done bool
	}
}

func newMiddleware(
	conf *configuration.Ledger,
	db *storage.DB,
	h *MessageHandler,
) *middleware {
	return &middleware{
		db:                                 db,
		handler:                            h,
		jetCoordinator:                     h.JetCoordinator,
		messageBus:                         h.Bus,
		pulseStorage:                       h.PulseStorage,
		earlyRequestCircuitBreakerProvider: &earlyRequestCircuitBreakerProvider{breakers: map[core.RecordID]*requestCircuitBreakerProvider{}},
		conf:                               conf,
		sequencer: map[core.RecordID]*struct {
			sync.Mutex
			done bool
		}{},
	}
}

type jetKey struct{}

func contextWithJet(ctx context.Context, jetID core.RecordID) context.Context {
	return context.WithValue(ctx, jetKey{}, jetID)
}

func jetFromContext(ctx context.Context) core.RecordID {
	val := ctx.Value(jetKey{})
	j, ok := val.(core.RecordID)
	if !ok {
		panic("failed to extract jet from context")
	}

	return j
}

func (m *middleware) zeroJetForHeavy(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		return handler(contextWithJet(ctx, *jet.NewID(0, nil)), parcel)
	}
}

func (m *middleware) checkJet(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		msg := parcel.Message()
		if msg.DefaultTarget() == nil {
			return nil, errors.New("unexpected message")
		}

		logger := inslogger.FromContext(ctx)
		logger.Debugf("checking jet for %v", parcel.Type().String())

		// FIXME: @andreyromancev. 17.01.19. Temporary allow any genesis request. Remove it.
		if parcel.Pulse() == core.FirstPulseNumber {
			logger.Debugf("genesis pulse shortcut")
			return handler(contextWithJet(ctx, *jet.NewID(0, nil)), parcel)
		}

		// Check token jet.
		token := parcel.DelegationToken()
		if token != nil {
			logger.Debugf("received token. returning any jet")
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
			}
			tree, err := m.db.GetJetTree(ctx, pulse)
			if err != nil {
				return nil, err
			}

			jetID, actual := tree.Find(target)
			if !actual {
				inslogger.FromContext(ctx).Errorf(
					"got message of type %s with redirect token,"+
						" but jet %s for pulse %d is not actual",
					msg.Type(), jetID.DebugString(), pulse,
				)
			}

			return handler(contextWithJet(ctx, *jetID), parcel)
		}

		// Calculate jet for current pulse.
		var jetID core.RecordID
		if msg.DefaultTarget().Record().Pulse() == core.PulseNumberJet {
			logger.Debugf("special pulse number (jet). returning jet from message")
			jetID = *msg.DefaultTarget().Record()
		} else {
			j, actual, err := m.fetchJet(ctx, *msg.DefaultTarget().Record(), parcel.Pulse())
			if err != nil {
				return nil, errors.Wrap(err, "failed to fetch jet tree")
			}
			if !actual {
				return &reply.JetMiss{JetID: *j}, nil
			}
			jetID = *j
		}

		// Check if jet is ours.
		node, err := m.jetCoordinator.LightExecutorForJet(ctx, jetID, parcel.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate executor for jet")
		}
		if *node != m.jetCoordinator.Me() {
			return &reply.JetMiss{JetID: jetID}, nil
		}

		return handler(contextWithJet(ctx, jetID), parcel)
	}
}

func (m *middleware) saveParcel(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		jetID := jetFromContext(ctx)
		pulse, err := m.pulseStorage.Current(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println("saveParcel, pulse - ", pulse.PulseNumber)
		err = m.db.SetMessage(ctx, jetID, pulse.PulseNumber, parcel)
		if err != nil {
			return nil, err
		}

		return handler(ctx, parcel)
	}
}

func (m *middleware) checkHeavySync(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		// TODO: @andreyromancev. 10.01.2019. Uncomment to enable backpressure for writing requests.
		// Currently disabled due to big initial difference in pulse numbers, which prevents requests from being accepted.
		// jetID := jetFromContext(ctx)
		// replicated, err := m.db.GetReplicatedPulse(ctx, jetID)
		// if err != nil {
		// 	return nil, err
		// }
		// if parcel.Pulse()-replicated >= m.conf.LightChainLimit {
		// 	return nil, errors.New("failed to write data (waiting for heavy replication)")
		// }

		return handler(ctx, parcel)
	}
}

func (m *middleware) fetchJet(
	ctx context.Context, target core.RecordID, pulse core.PulseNumber,
) (*core.RecordID, bool, error) {
	// Look in the local tree. Return if the actual jet found.
	tree, err := m.db.GetJetTree(ctx, pulse)
	if err != nil {
		return nil, false, err
	}
	jetID, actual := tree.Find(target)
	if actual {
		inslogger.FromContext(ctx).Debugf(
			"we believe object %s is in JET %s", target.String(), jetID.DebugString(),
		)
		return jetID, actual, nil
	}

	inslogger.FromContext(ctx).Debugf(
		"jet %s is not actual in our tree, asking neighbors for jet of object %s",
		jetID.DebugString(), target.String(),
	)

	m.seqMutex.Lock()
	if _, ok := m.sequencer[*jetID]; !ok {
		m.sequencer[*jetID] = &struct {
			sync.Mutex
			done bool
		}{}
	}
	mu := m.sequencer[*jetID]
	m.seqMutex.Unlock()

	mu.Lock()
	if mu.done {
		mu.Unlock()
		inslogger.FromContext(ctx).Debugf(
			"somebody else updated actuality of jet %s, rechecking our DB",
			jetID.DebugString(),
		)
		return m.fetchJet(ctx, target, pulse)
	}
	defer func() {
		inslogger.FromContext(ctx).Debugf("done fetching jet, cleaning")

		mu.done = true
		mu.Unlock()

		m.seqMutex.Lock()
		inslogger.FromContext(ctx).Debugf("deleting sequencer for jet %s", jetID.DebugString())
		delete(m.sequencer, *jetID)
		m.seqMutex.Unlock()
	}()

	resJet, err := m.handler.fetchActualJetFromOtherNodes(ctx, target, pulse)
	if err != nil {
		return nil, false, err
	}

	err = m.db.UpdateJetTree(ctx, pulse, true, *resJet)
	if err != nil {
		inslogger.FromContext(ctx).Error(
			errors.Wrapf(err, "couldn't actualize jet %s", resJet.DebugString()),
		)
	}

	return resJet, true, nil
}
