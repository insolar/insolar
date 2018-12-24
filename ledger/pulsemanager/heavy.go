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

package pulsemanager

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
)

func (m *PulseManager) initJetSyncState(ctx context.Context) error {
	allJets, err := m.db.GetJets(ctx)
	if err != nil {
		return err
	}
	// not so effective, because we rescan pulses
	// but for now it is easier to do this in this way
	for jetID := range allJets {
		pulseNums, err := m.NextSyncPulses(ctx, jetID)
		if err != nil {
			return err
		}
		m.syncClientsPool.AddPulsesToSyncClient(ctx, jetID, false, pulseNums...)
	}
	return nil
}

// NextSyncPulses returns next pulse numbers for syncing to heavy node.
// If nothing to sync it returns nil, nil.
func (m *PulseManager) NextSyncPulses(ctx context.Context, jetID core.RecordID) ([]core.PulseNumber, error) {
	var (
		replicated core.PulseNumber
		err        error
	)
	if replicated, err = m.db.GetReplicatedPulse(ctx, jetID); err != nil {
		return nil, err
	}

	if replicated == 0 {
		return m.findAllCompleted(ctx, jetID, core.FirstPulseNumber)
	}
	next, nexterr := m.findnext(ctx, replicated)
	if nexterr != nil {
		return nil, nexterr
	}
	if next == nil {
		return nil, nil
	}
	return m.findAllCompleted(ctx, jetID, *next)
}

func (m *PulseManager) findAllCompleted(ctx context.Context, jetID core.RecordID, from core.PulseNumber) ([]core.PulseNumber, error) {
	node, err := m.JetCoordinator.LightExecutorForJet(ctx, jetID, from)
	if err != nil {
		return nil, errors.Wrapf(err, "check 'am I light' for pulse num %v failed", from)
	}
	next, err := m.findnext(ctx, from)
	if err != nil {
		return nil, err
	}
	if next == nil {
		// if next is not found, we haven't got next pulse
		// in such case we don't want to replicate unfinished pulse
		return nil, nil
	}

	var found []core.PulseNumber
	if *node == m.JetCoordinator.Me() {
		found = append(found, from)
	}
	extra, err := m.findAllCompleted(ctx, jetID, *next)
	if err != nil {
		return nil, err
	}
	return append(found, extra...), nil
}

func (m *PulseManager) findnext(ctx context.Context, from core.PulseNumber) (*core.PulseNumber, error) {
	pulse, err := m.db.GetPulse(ctx, from)
	if err != nil {
		return nil, errors.Wrapf(err, "GetPulse with pulse num %v failed", from)
	}
	return pulse.Next, nil
}
