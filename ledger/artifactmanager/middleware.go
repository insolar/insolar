/*
 *    Copyright 2018 Insolar
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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

const (
	fetchJetReties = 10
)

type middleware struct {
	db                                 *storage.DB
	jetCoordinator                     core.JetCoordinator
	messageBus                         core.MessageBus
	earlyRequestCircuitBreakerProvider *earlyRequestCircuitBreakerProvider
	conf                               *configuration.Ledger
	seqMutex                           sync.Mutex
	sequencer                          map[core.RecordID]*struct {
		sync.Mutex
		done bool
	}
}

func newMiddleware(
	conf *configuration.Ledger,
	db *storage.DB,
	jetCoordinator core.JetCoordinator,
	messageBus core.MessageBus,
) *middleware {
	return &middleware{
		db:                                 db,
		jetCoordinator:                     jetCoordinator,
		messageBus:                         messageBus,
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

func (m *middleware) checkJet(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		logger := inslogger.FromContext(ctx)

		msg := parcel.Message()
		if msg.DefaultTarget() == nil {
			return nil, errors.New("unexpected message")
		}

		// Check token jet.
		token := parcel.DelegationToken()
		if token != nil {
			// Calculate jet for target pulse.
			target := *msg.DefaultTarget().Record()
			tree, err := m.db.GetJetTree(ctx, target.Pulse())
			if err != nil {
				return nil, err
			}
			jetID, _ := tree.Find(target)
			return handler(contextWithJet(ctx, *jetID), parcel)
		}

		// Calculate jet for current pulse.
		var jetID core.RecordID
		if msg.DefaultTarget().Record().Pulse() == core.PulseNumberJet {
			jetID = *msg.DefaultTarget().Record()
		} else {
			j, actual, err := m.fetchJet(ctx, *msg.DefaultTarget().Record(), parcel.Pulse(), fetchJetReties)
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
		// FIXME: This probably will never happen, because we won't receive hot records.
		if *node != m.jetCoordinator.Me() {
			// TODO: sergey.morozov 2018-12-21 This is hack. Must implement correct Jet checking for HME.
			logger.Debugf("checkJet: [ HACK ] checking if I am Heavy Material")
			heavy, err := m.jetCoordinator.Heavy(ctx, parcel.Pulse())
			if err != nil {
				logger.Debugf("checkJet: [ HACK ] failed to check for Heavy role")
				return nil, errors.Wrap(err, "[ HACK ] failed to check for heavy role")
			}
			if *heavy == m.jetCoordinator.Me() {
				logger.Debugf("checkJet: [ HACK ] I am Heavy. Accept parcel.")
				return handler(contextWithJet(ctx, jetID), parcel)
			}

			return &reply.JetMiss{JetID: jetID}, nil
		}

		return handler(contextWithJet(ctx, jetID), parcel)
	}
}

func (m *middleware) saveParcel(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		jetID := jetFromContext(ctx)
		pulse, err := m.db.GetLatestPulse(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println("saveParcel, pulse - ", pulse.Pulse.PulseNumber)
		err = m.db.SetMessage(ctx, jetID, pulse.Pulse.PulseNumber, parcel)
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
	ctx context.Context, target core.RecordID, pulse core.PulseNumber, retries int,
) (*core.RecordID, bool, error) {
	// Look in the local tree. Return if the actual jet found.
	tree, err := m.db.GetJetTree(ctx, pulse)
	if err != nil {
		return nil, false, err
	}
	jetID, actual := tree.Find(target)
	if actual || retries < 0 {
		if retries < 0 {
			fmt.Printf("[not mine] %v, %v", jetID.JetIDString(), target.String())
			fmt.Println()
		}
		if actual {
			fmt.Printf("[mine] %v, %v", jetID.JetIDString(), target.String())
			fmt.Println()
		}
		return jetID, actual, nil
	}

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
		return m.fetchJet(ctx, target, pulse, retries-1)
	}
	defer func() {
		mu.done = true
		mu.Unlock()

		m.seqMutex.Lock()
		delete(m.sequencer, *jetID)
		m.seqMutex.Unlock()
	}()

	// Couldn't find the actual jet locally. Ask for the jet from the previous executor.
	prevPulse, err := m.db.GetPreviousPulse(ctx, pulse)
	if err != nil {
		return nil, false, err
	}

	nodes, err := m.db.GetActiveNodesByRole(prevPulse.Pulse.PulseNumber, core.StaticRoleLightMaterial)
	if err != nil {
		return nil, false, err
	}

	me := m.jetCoordinator.Me()
	for i, _ := range nodes {
		if nodes[i].ID() == me {
			nodes = append(nodes[:i], nodes[i+1:]...)
			break
		}
	}

	num := len(nodes)
	if num == 0 {
		inslogger.FromContext(ctx).Error(
			"This shouldn't happen. We're solo active light material and have no active jet tree",
		)

		return nil, false, errors.New("impossible situation")
	}

	wg := sync.WaitGroup{}
	wg.Add(num)

	res := make([]*reply.Jet, num)
	for i, node := range nodes {
		go func(i int, node core.Node) {
			defer wg.Done()

			nodeID := node.ID()
			rep, err := m.messageBus.Send(
				ctx,
				&message.GetJet{Object: target},
				&core.MessageSendOptions{Receiver: &nodeID},
			)
			if err != nil {
				inslogger.FromContext(ctx).Error(
					errors.Wrap(err, "couldn't get jet"),
				)
				return
			}

			r, ok := rep.(*reply.Jet)
			if !ok {
				inslogger.FromContext(ctx).Errorf("unexpected reply: %#v\n", rep)
				return
			}
			res[i] = r
		}(i, node)
	}
	wg.Wait()

	for _, r := range res {
		if r == nil {
			continue
		}
		inslogger.FromContext(ctx).Debugf("Got jet %s Actual is %s", r.ID.JetIDString(), r.Actual)
		if !r.Actual {
			continue
		}

		err = m.db.UpdateJetTree(ctx, pulse, r.Actual, r.ID)
		if err != nil {
			return nil, false, err
		}

		return &r.ID, true, nil
	}

	inslogger.FromContext(ctx).Error(
		"All active light material nodes have no actual jet tree",
	)
	inslogger.FromContext(ctx).Error(
		"My guess is: ", jetID.JetIDString(),
	)
	return nil, false, errors.New("impossible situation")
}
