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

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/instracer"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	LogicRunner    insolar.LogicRunner `inject:""`
	NodeNet        network.NodeNetwork `inject:""` //nolint:staticcheck
	NodeSetter     node.Modifier       `inject:""`
	PulseAccessor  pulse.Accessor      `inject:""`
	PulseAppender  pulse.Appender      `inject:""`
	JetModifier    jet.Modifier        `inject:""`
	FlowDispatcher dispatcher.Dispatcher

	// setLock locks Set method call.
	setLock sync.RWMutex
	// saves PM stopping mode
	stopped bool
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager() *PulseManager {
	return &PulseManager{}
}

// Set set's new pulse.
func (m *PulseManager) Set(ctx context.Context, newPulse insolar.Pulse) error {
	m.setLock.Lock()
	defer m.setLock.Unlock()
	if m.stopped {
		return errors.New("can't call Set method on PulseManager after stop")
	}

	ctx, logger := inslogger.WithField(ctx, "new_pulse", newPulse.PulseNumber.String())
	logger.Debug("received pulse")

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
		fromNetwork := m.NodeNet.GetAccessor(newPulse.PulseNumber).GetWorkingNodes()
		if len(fromNetwork) == 0 {
			logger.Errorf("received zero nodes for pulse %d", newPulse.PulseNumber)
			return nil
		}
		toSet := make([]insolar.Node, 0, len(fromNetwork))
		for _, n := range fromNetwork {
			toSet = append(toSet, insolar.Node{ID: n.ID(), Role: n.Role()})
		}
		err := m.NodeSetter.Set(newPulse.PulseNumber, toSet)
		if err != nil {
			panic(errors.Wrap(err, "call of SetActiveNodes failed"))
		}
	}

	storagePulse, err := m.PulseAccessor.Latest(ctx)
	if err == pulse.ErrNotFound {
		storagePulse = *insolar.GenesisPulse
	} else if err != nil {
		return errors.Wrap(err, "call of GetLatestPulseNumber failed")
	}

	m.FlowDispatcher.ClosePulse(ctx, storagePulse)

	err = m.JetModifier.Clone(ctx, storagePulse.PulseNumber, newPulse.PulseNumber, false)
	if err != nil {
		return errors.Wrapf(err, "failed to clone jet.Tree fromPulse=%v toPulse=%v", storagePulse.PulseNumber, newPulse.PulseNumber)
	}

	if err := m.PulseAppender.Append(ctx, newPulse); err != nil {
		return errors.Wrap(err, "call of AddPulse failed")
	}

	err = m.LogicRunner.OnPulse(ctx, storagePulse, newPulse)
	if err != nil {
		return err
	}

	m.FlowDispatcher.BeginPulse(ctx, newPulse)

	return nil
}

// Start starts pulse manager.
func (m *PulseManager) Start(ctx context.Context) error {
	return nil
}

// Stop stops PulseManager.
func (m *PulseManager) Stop(ctx context.Context) error {
	// There should not to be any Set call after Stop call
	m.setLock.Lock()
	defer m.setLock.Unlock()

	m.stopped = true
	return nil
}
