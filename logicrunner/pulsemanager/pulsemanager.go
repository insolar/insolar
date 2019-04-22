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
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ActiveListSwapper is required by network to swap active list.
type ActiveListSwapper interface {
	MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error
}

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	LR                insolar.LogicRunner       `inject:""`
	Bus               insolar.MessageBus        `inject:""`
	NodeNet           insolar.NodeNetwork       `inject:""`
	GIL               insolar.GlobalInsolarLock `inject:""`
	ActiveListSwapper ActiveListSwapper         `inject:""`
	NodeSetter        node.Modifier             `inject:""`
	Nodes             node.Accessor             `inject:""`
	PulseAccessor     pulse.Accessor            `inject:""`
	PulseAppender     pulse.Appender            `inject:""`
	JetModifier       jet.Modifier              `inject:""`

	currentPulse insolar.Pulse

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

	err := m.setUnderGilSection(ctx, newPulse, persist)
	if err != nil {
		return err
	}

	if !persist {
		return nil
	}

	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
	}

	err = m.LR.OnPulse(ctx, newPulse)
	if err != nil {
		return err
	}

	return nil
}

func (m *PulseManager) setUnderGilSection(ctx context.Context, newPulse insolar.Pulse, persist bool) error {
	m.GIL.Acquire(ctx)
	ctx, span := instracer.StartSpan(ctx, "pulse.gil_locked")

	storagePulse, err := m.PulseAccessor.Latest(ctx)
	if err == pulse.ErrNotFound {
		storagePulse = *insolar.GenesisPulse
	} else if err != nil {
		return errors.Wrap(err, "call of GetLatestPulseNumber failed")
	}

	defer span.End()
	defer m.GIL.Release(ctx)

	logger := inslogger.FromContext(ctx)
	logger.WithFields(map[string]interface{}{
		"new_pulse": newPulse.PulseNumber,
		"persist":   persist,
	}).Debugf("received pulse")

	// swap pulse
	m.currentPulse = newPulse

	// swap active nodes
	err = m.ActiveListSwapper.MoveSyncToActive(ctx, newPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to apply new active node list")
	}
	if persist {
		if err := m.PulseAppender.Append(ctx, newPulse); err != nil {
			return errors.Wrap(err, "call of AddPulse failed")
		}
		fromNetwork := m.NodeNet.GetWorkingNodes()
		toSet := make([]insolar.Node, 0, len(fromNetwork))
		for _, node := range fromNetwork {
			toSet = append(toSet, insolar.Node{ID: node.ID(), Role: node.Role()})
		}
		err = m.NodeSetter.Set(newPulse.PulseNumber, toSet)
		if err != nil {
			return errors.Wrap(err, "call of SetActiveNodes failed")
		}
	}

	m.JetModifier.Clone(ctx, storagePulse.PulseNumber, newPulse.PulseNumber)
	return nil
}

// Start starts pulse manager.
func (m *PulseManager) Start(ctx context.Context) error {
	origin := m.NodeNet.GetOrigin()
	err := m.NodeSetter.Set(insolar.FirstPulseNumber, []insolar.Node{{ID: origin.ID(), Role: origin.Role()}})
	if err != nil && err != node.ErrOverride {
		return err
	}

	return nil
}

// Stop stops PulseManager.
func (m *PulseManager) Stop(ctx context.Context) error {
	// There should not to be any Set call after Stop call
	m.setLock.Lock()
	m.stopped = true
	m.setLock.Unlock()

	return nil
}
