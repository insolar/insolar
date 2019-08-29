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

package executor

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
)

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	setLock sync.RWMutex

	nodeNet          network.NodeNetwork
	dispatchers      []dispatcher.Dispatcher
	nodeSetter       node.Modifier
	pulseAccessor    pulse.Accessor
	pulseAppender    pulse.Appender
	jetReleaser      JetReleaser
	jetSplitter      JetSplitter
	lightReplicator  LightReplicator
	hotSender        HotSender
	writeManager     WriteManager
	stateIniter      StateIniter
	hotStatusChecker HotDataStatusChecker
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(
	nodeNet network.NodeNetwork,
	dispatchers []dispatcher.Dispatcher,
	nodeSetter node.Modifier,
	pulseAccessor pulse.Accessor,
	pulseAppender pulse.Appender,
	jetReleaser JetReleaser,
	jetSplitter JetSplitter,
	lightReplicator LightReplicator,
	hotSender HotSender,
	writeManager WriteManager,
	stateIniter StateIniter,
	hotStatusChecker HotDataStatusChecker,
) *PulseManager {
	pm := &PulseManager{
		nodeNet:          nodeNet,
		dispatchers:      dispatchers,
		jetSplitter:      jetSplitter,
		nodeSetter:       nodeSetter,
		pulseAccessor:    pulseAccessor,
		pulseAppender:    pulseAppender,
		jetReleaser:      jetReleaser,
		lightReplicator:  lightReplicator,
		hotSender:        hotSender,
		writeManager:     writeManager,
		stateIniter:      stateIniter,
		hotStatusChecker: hotStatusChecker,
	}
	return pm
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(ctx context.Context, newPulse insolar.Pulse) error {
	logger := inslogger.FromContext(ctx)
	logger.WithFields(map[string]interface{}{
		"new_pulse": newPulse.PulseNumber,
	}).Info("trying to set new pulse")

	m.setLock.Lock()
	defer m.setLock.Unlock()

	logger.Debug("behind set lock")
	ctx, span := instracer.StartSpan(
		ctx, "PulseManager.Set", trace.WithSampler(trace.AlwaysSample()),
	)
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(newPulse.PulseNumber)),
	)
	defer span.End()

	// Dealing with node lists.
	logger.Debug("dealing with node lists.")
	{
		fromNetwork := m.nodeNet.GetAccessor(newPulse.PulseNumber).GetWorkingNodes()
		if len(fromNetwork) == 0 {
			logger.Errorf("received zero nodes for pulse %d", newPulse.PulseNumber)
			return nil
		}
		toSet := make([]insolar.Node, 0, len(fromNetwork))
		for _, n := range fromNetwork {
			toSet = append(toSet, insolar.Node{ID: n.ID(), Role: n.Role()})
		}
		err := m.nodeSetter.Set(newPulse.PulseNumber, toSet)
		if err != nil {
			logger.Panic(errors.Wrap(err, "call of SetActiveNodes failed"))
		}
	}

	logger.Debug("before preparing state")
	justJoined, jets, err := m.stateIniter.PrepareState(ctx, newPulse.PulseNumber)
	if err != nil {
		logger.Panic(errors.Wrap(err, "failed to prepare light for start"))
	}
	stats.Record(ctx, statJets.M(int64(len(jets))))

	endedPulse, err := m.pulseAccessor.Latest(ctx)
	if err != nil {
		logger.Panic(errors.Wrap(err, "failed to fetch ended pulse"))
	}

	// Changing pulse.
	logger.Debug("before changing pulse")
	{
		logger.Debug("before dispatcher closePulse")
		for _, d := range m.dispatchers {
			d.ClosePulse(ctx, newPulse)
		}

		if !justJoined {
			logger.Debug("before parsing jets")
			for _, jet := range jets {

				logger.WithFields(map[string]interface{}{
					"jet_id":     jet.DebugString(),
					"endedPulse": endedPulse.PulseNumber,
				}).Debug("before hotStatusChecker.IsReceived")

				if !m.hotStatusChecker.IsReceived(ctx, jet, endedPulse.PulseNumber) {
					log.Fatal("hot data for jet: %s and pulse: %d wasn't received", jet.DebugString(), endedPulse.PulseNumber)
				}
			}

			logger.WithFields(map[string]interface{}{
				"newPulse":   newPulse.PulseNumber,
				"endedPulse": endedPulse.PulseNumber,
			}).Debug("before jetSplitter.Do")
			jets, err = m.jetSplitter.Do(ctx, endedPulse.PulseNumber, newPulse.PulseNumber, jets, true)
			if err != nil {
				logger.Panic(errors.Wrap(err, "failed to split jets"))
			}
		}

		logger.WithFields(map[string]interface{}{
			"endedPulse": endedPulse.PulseNumber,
		}).Debugf("before jetReleaser.CloseAllUntil")
		m.jetReleaser.CloseAllUntil(ctx, endedPulse.PulseNumber)

		logger.WithFields(map[string]interface{}{
			"endedPulse": endedPulse.PulseNumber,
		}).Debugf("before writeManager.CloseAndWait")
		err = m.writeManager.CloseAndWait(ctx, endedPulse.PulseNumber)
		if err != nil {
			logger.Panic(errors.Wrap(err, "can't close pulse for writing"))
		}

		logger.WithField("newPulse.PulseNumber", newPulse.PulseNumber).Debug("before writeManager.Open")
		err = m.writeManager.Open(ctx, newPulse.PulseNumber)
		if err != nil {
			logger.Panic(errors.Wrap(err, "failed to open pulse for writing"))
		}

		logger.WithField("newPulse.PulseNumber", newPulse.PulseNumber).Debug("before pulseAppender.Append")
		if err := m.pulseAppender.Append(ctx, newPulse); err != nil {
			logger.Panic(errors.Wrap(err, "failed to add pulse"))
		}

		logger.WithField("newPulse", newPulse.PulseNumber).Debugf("before dispatcher.BeginPulse", newPulse)
		for _, d := range m.dispatchers {
			d.BeginPulse(ctx, newPulse)
		}
	}

	if !justJoined {
		logger.Info("going to send hots")
		go func() {
			err = m.hotSender.SendHot(ctx, endedPulse.PulseNumber, newPulse.PulseNumber, jets)
			if err != nil {
				logger.Error("send Hot failed: ", err)
			}
		}()
		logger.Info("going to notify cleaner about new pulse")
		go m.lightReplicator.NotifyAboutPulse(ctx, newPulse.PulseNumber)
	}

	logger.Info("new pulse is set")
	return nil
}
