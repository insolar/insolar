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

package pulsemanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/artifactmanager"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/replication"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/pulsemanager.ActiveListSwapper -o ../../../testutils -s _mock.go

type ActiveListSwapper interface {
	MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error
}

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	Bus               insolar.MessageBus        `inject:""`
	NodeNet           insolar.NodeNetwork       `inject:""`
	GIL               insolar.GlobalInsolarLock `inject:""`
	ActiveListSwapper ActiveListSwapper         `inject:""`
	MessageHandler    *artifactmanager.MessageHandler

	JetReleaser hot.JetReleaser `inject:""`

	JetModifier jet.Modifier `inject:""`
	JetSplitter executor.JetSplitter

	NodeSetter node.Modifier `inject:""`
	Nodes      node.Accessor `inject:""`

	PulseAccessor   pulse.Accessor   `inject:""`
	PulseCalculator pulse.Calculator `inject:""`
	PulseAppender   pulse.Appender   `inject:""`

	LightReplicator replication.LightReplicator
	HotSender       executor.HotSender

	WriteManager hot.WriteManager

	// setLock locks Set method call.
	setLock sync.RWMutex
	// saves PM stopping mode
	stopped bool
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(
	jetSplitter executor.JetSplitter,
	lightToHeavySyncer replication.LightReplicator,
	writeManager hot.WriteManager,
	hotSender executor.HotSender,
) *PulseManager {
	pm := &PulseManager{
		JetSplitter:     jetSplitter,
		LightReplicator: lightToHeavySyncer,
		HotSender:       hotSender,
		WriteManager:    writeManager,
	}
	return pm
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(ctx context.Context, newPulse insolar.Pulse, persist bool) error {
	m.setLock.Lock()
	defer m.setLock.Unlock()
	if m.stopped {
		return errors.New("can't call Set method on PulseManager after stop")
	}

	ctx, span := instracer.StartSpan(
		ctx, "pulse.process", trace.WithSampler(trace.AlwaysSample()),
	)
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(newPulse.PulseNumber)),
	)
	defer span.End()

	jets, oldPulse, prevPN, err := m.setUnderGilSection(ctx, newPulse, persist)
	if err != nil {
		return err
	}

	if !persist {
		return nil
	}

	logger := inslogger.FromContext(ctx)

	if oldPulse != nil && prevPN != nil {
		err = m.WriteManager.CloseAndWait(ctx, oldPulse.PulseNumber)
		if err != nil {
			logger.Error("can't close pulse for writing", err)
		}
		err = m.HotSender.SendHot(ctx, jets, oldPulse.PulseNumber, newPulse.PulseNumber)
		if err != nil {
			return err
		}
		go m.LightReplicator.NotifyAboutPulse(ctx, newPulse.PulseNumber)
	}

	err = m.WriteManager.Open(ctx, newPulse.PulseNumber)
	if err != nil {
		logger.Error("can't open pulse for writing", err)
	}

	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
	}

	if m.MessageHandler != nil {
		m.MessageHandler.OnPulse(ctx, newPulse)
	}

	return nil
}

func (m *PulseManager) setUnderGilSection(
	ctx context.Context, newPulse insolar.Pulse, persist bool,
) (
	[]executor.JetInfo, *insolar.Pulse, *insolar.PulseNumber, error,
) {
	var (
		oldPulse *insolar.Pulse
		prevPN   *insolar.PulseNumber
	)

	m.GIL.Acquire(ctx)
	ctx, span := instracer.StartSpan(ctx, "pulse.gil_locked")
	defer span.End()
	defer m.GIL.Release(ctx)

	// FIXME: @andreyromancev. 17.12.18. return insolar.Pulse here.
	storagePulse, err := m.PulseAccessor.Latest(ctx)
	if err != nil && err != pulse.ErrNotFound {
		return nil, nil, nil, errors.Wrap(err, "call of GetLatestPulseNumber failed")
	}

	if err != pulse.ErrNotFound {
		oldPulse = &storagePulse
		pp, err := m.PulseCalculator.Backwards(ctx, oldPulse.PulseNumber, 1)
		if err == nil {
			prevPN = &pp.PulseNumber
		} else {
			prevPN = &insolar.GenesisPulse.PulseNumber
		}
		ctx, _ = inslogger.WithField(ctx, "current_pulse", fmt.Sprintf("%d", oldPulse.PulseNumber))
	}

	logger := inslogger.FromContext(ctx)
	logger.WithFields(map[string]interface{}{
		"new_pulse": newPulse.PulseNumber,
		"persist":   persist,
	}).Debugf("received pulse")

	// swap active nodes
	err = m.ActiveListSwapper.MoveSyncToActive(ctx, newPulse.PulseNumber)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to apply new active node list")
	}
	if persist {
		if err := m.PulseAppender.Append(ctx, newPulse); err != nil {
			return nil, nil, nil, errors.Wrap(err, "call of AddPulse failed")
		}
		fromNetwork := m.NodeNet.GetWorkingNodes()
		toSet := make([]insolar.Node, 0, len(fromNetwork))
		for _, node := range fromNetwork {
			toSet = append(toSet, insolar.Node{ID: node.ID(), Role: node.Role()})
		}
		err = m.NodeSetter.Set(newPulse.PulseNumber, toSet)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "call of SetActiveNodes failed")
		}
	}

	var jets []executor.JetInfo
	if persist && prevPN != nil && oldPulse != nil {
		jets, err = m.JetSplitter.Do(ctx, *prevPN, oldPulse.PulseNumber, newPulse.PulseNumber)

		// We just joined to network
		if errors.Cause(err) == node.ErrNoNodes {
			return jets, oldPulse, prevPN, nil
		}
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if oldPulse != nil && prevPN != nil {
		m.prepareArtifactManagerMessageHandlerForNextPulse(ctx, newPulse)
	}

	if persist && oldPulse != nil {
		nodes, err := m.Nodes.All(oldPulse.PulseNumber)
		if err != nil {
			return nil, nil, nil, err
		}
		// No active nodes for pulse. It means there was no processing (network start).
		if len(nodes) == 0 {
			// Activate zero jet for jet tree and unlock jet waiter.
			zeroJet := insolar.NewJetID(0, nil)
			m.JetModifier.Update(ctx, newPulse.PulseNumber, true, *zeroJet)
			err := m.JetReleaser.Unlock(ctx, insolar.ID(*zeroJet))
			if err != nil {
				if err == artifactmanager.ErrWaiterNotLocked {
					inslogger.FromContext(ctx).Error(err)
				} else {
					return nil, nil, nil, errors.Wrap(err, "failed to unlock zero jet")
				}
			}
		}
	}

	return jets, oldPulse, prevPN, nil
}

func (m *PulseManager) prepareArtifactManagerMessageHandlerForNextPulse(ctx context.Context, newPulse insolar.Pulse) {
	ctx, span := instracer.StartSpan(ctx, "early.close")
	defer span.End()

	m.JetReleaser.ThrowTimeout(ctx, newPulse.PulseNumber)
}

// Start starts pulse manager
func (m *PulseManager) Start(ctx context.Context) error {
	origin := m.NodeNet.GetOrigin()
	err := m.NodeSetter.Set(insolar.FirstPulseNumber, []insolar.Node{{ID: origin.ID(), Role: origin.Role()}})
	if err != nil && err != node.ErrOverride {
		return err
	}

	return nil
}

// Stop stops PulseManager
func (m *PulseManager) Stop(ctx context.Context) error {
	// There should not to be any Set call after Stop call
	m.setLock.Lock()
	m.stopped = true
	m.setLock.Unlock()

	return nil
}
