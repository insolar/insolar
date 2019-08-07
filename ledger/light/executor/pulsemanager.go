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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	setLock sync.RWMutex

	nodeNet         insolar.NodeNetwork
	dispatcher      dispatcher.Dispatcher
	nodeSetter      node.Modifier
	pulseAccessor   pulse.Accessor
	pulseAppender   pulse.Appender
	jetReleaser     JetReleaser
	jetSplitter     JetSplitter
	lightReplicator LightReplicator
	hotSender       HotSender
	writeManager    WriteManager
	stateIniter     StateIniter
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(
	nodeNet insolar.NodeNetwork,
	disp dispatcher.Dispatcher,
	nodeSetter node.Modifier,
	pulseAccessor pulse.Accessor,
	pulseAppender pulse.Appender,
	jetReleaser JetReleaser,
	jetSplitter JetSplitter,
	lightReplicator LightReplicator,
	hotSender HotSender,
	writeManager WriteManager,
	stateIniter StateIniter,
) *PulseManager {
	pm := &PulseManager{
		nodeNet:         nodeNet,
		dispatcher:      disp,
		jetSplitter:     jetSplitter,
		nodeSetter:      nodeSetter,
		pulseAccessor:   pulseAccessor,
		pulseAppender:   pulseAppender,
		jetReleaser:     jetReleaser,
		lightReplicator: lightReplicator,
		hotSender:       hotSender,
		writeManager:    writeManager,
		stateIniter:     stateIniter,
	}
	return pm
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(ctx context.Context, newPulse insolar.Pulse) error {
	logger := inslogger.FromContext(ctx)

	m.setLock.Lock()
	defer m.setLock.Unlock()

	ctx, span := instracer.StartSpan(
		ctx, "PulseManager.Set", trace.WithSampler(trace.AlwaysSample()),
	)
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(newPulse.PulseNumber)),
	)
	defer span.End()

	logger.WithFields(map[string]interface{}{
		"new_pulse": newPulse.PulseNumber,
	}).Debugf("received pulse")

	// Dealing with node lists.
	{
		fromNetwork := m.nodeNet.GetWorkingNodes()
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
			panic(errors.Wrap(err, "call of SetActiveNodes failed"))
		}
	}

	justJoined, jets, err := m.stateIniter.PrepareState(ctx, newPulse.PulseNumber)
	if err != nil {
		panic(errors.Wrap(err, "failed to prepare light for start"))
	}
	endedPulse, err := m.pulseAccessor.Latest(ctx)
	if err != nil {
		panic(errors.Wrap(err, "failed to fetch ended pulse"))
	}

	// Changing pulse.
	{
		m.dispatcher.ClosePulse(ctx, newPulse)

		createDrops := !justJoined
		jets, err = m.jetSplitter.Do(ctx, endedPulse.PulseNumber, newPulse.PulseNumber, jets, createDrops)
		if err != nil {
			panic(errors.Wrap(err, "failed to split jets"))
		}

		m.jetReleaser.CloseAllUntil(ctx, endedPulse.PulseNumber)

		err = m.writeManager.CloseAndWait(ctx, endedPulse.PulseNumber)
		if err != nil {
			panic(errors.Wrap(err, "can't close pulse for writing"))
		}

		err = m.writeManager.Open(ctx, newPulse.PulseNumber)
		if err != nil {
			panic(errors.Wrap(err, "failed to open pulse for writing"))
		}

		if err := m.pulseAppender.Append(ctx, newPulse); err != nil {
			panic(errors.Wrap(err, "failed to add pulse"))
		}

		m.dispatcher.BeginPulse(ctx, newPulse)
	}

	if !justJoined {
		go func() {
			err = m.hotSender.SendHot(ctx, endedPulse.PulseNumber, newPulse.PulseNumber, jets)
			if err != nil {
				logger.Error("send Hot failed: ", err)
			}
		}()
		go m.lightReplicator.NotifyAboutPulse(ctx, newPulse.PulseNumber)
	}

	return nil
}
