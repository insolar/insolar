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
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/network"

	"github.com/pkg/errors"
)

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	NodeNet            network.NodeNetwork         `inject:""`
	NodeSetter         node.Modifier               `inject:""`
	Nodes              node.Accessor               `inject:""`
	PulseAppender      pulse.Appender              `inject:""`
	PulseAccessor      pulse.Accessor              `inject:""`
	FinalizationKeeper executor.FinalizationKeeper `inject:""`
	JetModifier        jet.Modifier                `inject:""`

	dispatcher dispatcher.Dispatcher

	currentPulse insolar.Pulse
	StartPulse   pulse.StartPulse

	// setLock locks Set method call.
	setLock sync.RWMutex
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(disp dispatcher.Dispatcher) *PulseManager {
	pm := &PulseManager{
		currentPulse: *insolar.GenesisPulse,
		dispatcher:   disp,
	}
	return pm
}

// Set set's new pulse.
func (m *PulseManager) Set(ctx context.Context, newPulse insolar.Pulse) error {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"new_pulse": newPulse.PulseNumber,
	})

	logger.Info("PulseManager.Set is about to acquire the lock")

	// In Go the goroutine which first tries to acquire the lock will get it first
	// (fairness property). See https://play.golang.org/p/Vkj7parznba
	m.setLock.Lock()
	defer m.setLock.Unlock()

	logger.Info("PulseManager.Set acquired the lock")

	ctx, span := instracer.StartSpan(ctx, "PulseManager.Set")
	span.SetTag("pulse.PulseNumber", int64(newPulse.PulseNumber))
	defer span.Finish()

	if m.dispatcher != nil {
		logger.Info("PulseManager.Set calls dispatcher.ClosePulse")
		m.dispatcher.ClosePulse(ctx, newPulse)
		logger.Info("PulseManager.Set returned from dispatcher.ClosePulse")
	}

	// Dealing with node lists.
	{
		logger.Info("PulseManager.Set deals with node list")

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
			logger.Panic(errors.Wrap(err, "call of SetActiveNodes failed"))
		}

		logger.Info("PulseManager.Set finished to deal with node list")
	}

	logger.Info("PulseManager.Set calls PulseAppender.Append")
	if err := m.PulseAppender.Append(ctx, newPulse); err != nil {
		instracer.AddError(span, err)
		logger.Error(err)
		return errors.Wrap(err, "call of AddPulse failed")
	}
	logger.Info("PulseManager.Set returned from PulseAppender.Append, about to call FinalizationKeeper.OnPulse")

	err := m.FinalizationKeeper.OnPulse(ctx, newPulse.PulseNumber)
	if err != nil {
		logger.Error(err)
		instracer.AddError(span, err)
		return errors.Wrap(err, "got error calling FinalizationKeeper.OnPulse")
	}

	logger.Info("PulseManager.Set returned from FinalizationKeeper.OnPulse, about to call StartPulse.SetStartPulse")
	m.StartPulse.SetStartPulse(ctx, newPulse)
	logger.Info("PulseManager.Set returned from StartPulse.SetStartPulse")
	if m.dispatcher != nil {
		logger.Info("PulseManager.Set about to call dispatcher.BeginPulse")
		m.dispatcher.BeginPulse(ctx, newPulse)
		logger.Info("PulseManager.Set returned from dispatcher.BeginPulse")
	}

	logger.Info("PulseManager.Set - All OK!")
	return nil
}
