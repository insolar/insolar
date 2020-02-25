// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsemanager

import (
	"context"
	"sync"
	"time"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/network"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/instracer"

	"github.com/pkg/errors"
)

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	LogicRunner   insolar.LogicRunner `inject:""`
	NodeNet       network.NodeNetwork `inject:""` //nolint:staticcheck
	NodeSetter    node.Modifier       `inject:""`
	PulseAccessor pulse.Accessor      `inject:""`
	PulseAppender pulse.Appender      `inject:""`
	JetModifier   jet.Modifier        `inject:""`
	dispatchers   []dispatcher.Dispatcher

	// setLock locks Set method call.
	setLock sync.RWMutex
	// saves PM stopping mode
	stopped bool
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager() *PulseManager {
	return &PulseManager{}
}

// AddDispatcher adds dispatchers to handling
// that could be done only when Set is not happening
func (m *PulseManager) AddDispatcher(d ...dispatcher.Dispatcher) {
	m.setLock.Lock()
	defer m.setLock.Unlock()

	m.dispatchers = append(m.dispatchers, d...)
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

	ctx, span := instracer.StartSpan(ctx, "PulseManager.Set")
	span.SetTag("pulse.PulseNumber", int64(newPulse.PulseNumber))

	onPulseStart := time.Now()
	defer func() {
		stats.Record(ctx, metrics.PulseManagerOnPulseTiming.M(float64(time.Since(onPulseStart).Nanoseconds())/1e6))
		span.Finish()
	}()

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

	for _, d := range m.dispatchers {
		d.ClosePulse(ctx, storagePulse)
	}

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

	for _, d := range m.dispatchers {
		d.BeginPulse(ctx, newPulse)
	}

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
