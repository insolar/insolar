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
	"github.com/insolar/insolar/network"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/heavy/executor"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

var (
	errZeroNodes = errors.New("zero nodes from network")
)

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	Bus                insolar.MessageBus          `inject:""`
	NodeNet            network.NodeNetwork         `inject:""`
	GIL                insolar.GlobalInsolarLock   `inject:""`
	NodeSetter         node.Modifier               `inject:""`
	Nodes              node.Accessor               `inject:""`
	PulseAppender      pulse.Appender              `inject:""`
	PulseAccessor      pulse.Accessor              `inject:""`
	FinalizationKeeper executor.FinalizationKeeper `inject:""`
	JetModifier        jet.Modifier                `inject:""`

	currentPulse insolar.Pulse
	StartPulse   pulse.StartPulse

	// setLock locks Set method call.
	setLock sync.RWMutex
	// saves PM stopping mode
	stopped bool
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager() *PulseManager {
	pm := &PulseManager{
		currentPulse: *insolar.GenesisPulse,
	}
	return pm
}

// Set set's new pulse.
func (m *PulseManager) Set(ctx context.Context, newPulse insolar.Pulse) error {
	logger := inslogger.FromContext(ctx)
	logger.WithFields(map[string]interface{}{
		"new_pulse": newPulse.PulseNumber,
	}).Info("trying to set new pulse")

	m.setLock.Lock()
	defer m.setLock.Unlock()

	logger.Debug("behind set lock")

	if m.stopped {
		logger.Error(errors.New("can't call Set method on PulseManager after stop"))
		return errors.New("can't call Set method on PulseManager after stop")
	}

	ctx, span := instracer.StartSpan(
		ctx, "PulseManager.Set", trace.WithSampler(trace.AlwaysSample()),
	)
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(newPulse.PulseNumber)),
	)
	defer span.End()

	logger.Debug("before calling to setUnderGilSection")
	err := m.setUnderGilSection(ctx, newPulse)
	if err != nil {
		logger.Error(err)
		instracer.AddError(span, err)
		if err == errZeroNodes {
			logger.Info("setUnderGilSection return error: ", err)
			return nil
		}
		return err
	}
	logger.Debug("after calling to setUnderGilSection")

	logger.Debug("before calling to FinalizationKeeper.OnPulse")
	err = m.FinalizationKeeper.OnPulse(ctx, newPulse.PulseNumber)
	if err != nil {
		logger.Error(err)
		instracer.AddError(span, err)
		return errors.Wrap(err, "got error calling FinalizationKeeper.OnPulse")
	}

	logger.Debug("before calling to StartPulse.SetStartPulse")
	m.StartPulse.SetStartPulse(ctx, newPulse)

	logger.Debug("before calling to Bus.OnPulse")
	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		instracer.AddError(span, err)
		logger.Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
	}

	logger.Info("new pulse is set")
	return nil
}

func (m *PulseManager) setUnderGilSection(ctx context.Context, newPulse insolar.Pulse) error {
	var (
		oldPulse *insolar.Pulse
	)
	logger := inslogger.FromContext(ctx)

	logger.Debug("before calling to GIL.Acquire")
	m.GIL.Acquire(ctx)

	ctx, span := instracer.StartSpan(ctx, "PulseManager.setUnderGilSection")
	defer span.End()

	defer m.GIL.Release(ctx)

	logger.Debug("before calling to PulseAccessor.Latest")
	storagePulse, err := m.PulseAccessor.Latest(ctx)
	if err != nil && err != pulse.ErrNotFound {
		instracer.AddError(span, err)
		logger.Error(err)
		return errors.Wrap(err, "call of m.PulseAccessor.Latest failed")
	}

	if err != pulse.ErrNotFound {
		oldPulse = &storagePulse
	}

	logger.Debug("set currentPulse")
	// swap pulse
	m.currentPulse = newPulse

	logger.Debug("calling to GetWorkingNodes")
	fromNetwork := m.NodeNet.GetAccessor(m.currentPulse.PulseNumber).GetWorkingNodes()
	toSet := make([]insolar.Node, 0, len(fromNetwork))
	if len(fromNetwork) == 0 {
		logger.Errorf("received zero nodes for pulse %d", newPulse.PulseNumber)
		return errZeroNodes
	}
	for _, node := range fromNetwork {
		toSet = append(toSet, insolar.Node{ID: node.ID(), Role: node.Role()})
	}
	logger.Debug("calling to NodeSetter.Set")
	err = m.NodeSetter.Set(newPulse.PulseNumber, toSet)
	if err != nil {
		logger.Error(err)
		instracer.AddError(span, err)
		return errors.Wrap(err, "call of SetActiveNodes failed")
	}

	logger.Debug("calling to JetModifier.Clone")
	err = m.JetModifier.Clone(ctx, storagePulse.PulseNumber, newPulse.PulseNumber, true)
	if err != nil {
		logger.Error(err)
		instracer.AddError(span, err)
		return errors.Wrapf(err, "failed to clone jet.Tree fromPulse=%v toPulse=%v", storagePulse.PulseNumber, newPulse.PulseNumber)
	}

	if oldPulse != nil {
		logger.Debug("calling to Nodes.All")
		nodes, err := m.Nodes.All(oldPulse.PulseNumber)
		if err != nil {
			logger.Info("oldPulse isn't nil. append new")
			if err := m.PulseAppender.Append(ctx, newPulse); err != nil {
				logger.Error(err)
				return errors.Wrap(err, "call of AddPulse failed")
			}
			return nil
		}
		// No active nodes for pulse. It means there was no processing (network start).
		if len(nodes) == 0 {
			// Activate zero jet for jet tree.
			logger.Debug("calling to JetModifier.Update")
			err := m.JetModifier.Update(ctx, newPulse.PulseNumber, false, insolar.ZeroJetID)
			if err != nil {
				logger.Error(err)
				return errors.Wrapf(err, "failed to update zeroJet")
			}
			logger.Infof("[PulseManager] activate zeroJet pulse=%v", newPulse.PulseNumber)
		}
	}

	logger.Info("save pulse to storage")
	if err := m.PulseAppender.Append(ctx, newPulse); err != nil {
		instracer.AddError(span, err)
		logger.Error(err)
		return errors.Wrap(err, "call of AddPulse failed")
	}
	logger.Info("pulse is saved to storage")
	return nil
}
