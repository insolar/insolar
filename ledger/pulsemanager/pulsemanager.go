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
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/utils/backoff"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/pulsemanager.ActiveListSwapper -o ../../testutils -s _mock.go
type ActiveListSwapper interface {
	MoveSyncToActive()
}

// PulseManager implements core.PulseManager.
type PulseManager struct {
	LR                    core.LogicRunner       `inject:""`
	Bus                   core.MessageBus        `inject:""`
	NodeNet               core.NodeNetwork       `inject:""`
	JetCoordinator        core.JetCoordinator    `inject:""`
	GIL                   core.GlobalInsolarLock `inject:""`
	RecentStorageProvider recentstorage.Provider `inject:""`
	ActiveListSwapper     ActiveListSwapper      `inject:""`

	currentPulse core.Pulse

	// internal stuff
	db *storage.DB
	// setLock locks Set method call.
	setLock sync.RWMutex
	stopped bool

	// Heavy sync stuff:
	//
	// is sync enabled at all
	enableSync bool
	// syncstates *jetSyncStates
	syncClientsPool *syncClientsPool
}

func backoffFromConfig(bconf configuration.Backoff) *backoff.Backoff {
	return &backoff.Backoff{
		Jitter: bconf.Jitter,
		Min:    bconf.Min,
		Max:    bconf.Max,
		Factor: bconf.Factor,
	}
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(db *storage.DB, conf configuration.Ledger) *PulseManager {
	pmconf := conf.PulseManager
	pm := &PulseManager{
		db:           db,
		currentPulse: *core.GenesisPulse,
		enableSync:   pmconf.HeavySyncEnabled,
	}
	// TODO: untie this circular dependency after moving sync client to separate component - 17.Dec.2018 @nordicdyno
	heavySyncPool := newSyncClientsPool(
		pm,
		clientOptions{
			syncMessageLimit: pmconf.HeavySyncMessageLimit,
			pulsesDeltaLimit: conf.LightChainLimit,
		})
	pm.syncClientsPool = heavySyncPool
	return pm
}

// Current returns copy (for concurrency safety) of current pulse structure.
func (m *PulseManager) Current(ctx context.Context) (*core.Pulse, error) {
	m.setLock.RLock()
	defer m.setLock.RUnlock()

	p := m.currentPulse
	return &p, nil
}

func (m *PulseManager) dropAllJets(ctx context.Context, pulse *storage.Pulse) error {
	jetIDs, err := m.db.GetJets(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get jets from storage")
	}
	var g errgroup.Group
	for jetID := range jetIDs {
		jetID := jetID
		g.Go(func() error {
			drop, dropSerialized, messages, err := m.createDrop(ctx, jetID, pulse)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", pulse)
			}

			hotRecordsError := m.processRecentObjects(
				ctx, jetID, pulse, &m.currentPulse, drop, dropSerialized)
			if hotRecordsError != nil {
				return errors.Wrap(err, "processRecentObjects failed")
			}

			dropErr := m.processDrop(ctx, jetID, pulse, &m.currentPulse, dropSerialized, messages)
			if dropErr != nil {
				return errors.Wrap(dropErr, "processDrop failed")
			}
			return nil
		})
	}
	return g.Wait()
}

func (m *PulseManager) createDrop(
	ctx context.Context,
	jetID core.RecordID,
	lastSlotPulse *storage.Pulse,
) (
	drop *jet.JetDrop,
	dropSerialized []byte,
	messages [][]byte,
	err error,
) {
	prevDrop, err := m.db.GetDrop(ctx, jetID, *lastSlotPulse.Prev)
	if err != nil {
		return nil, nil, nil, err
	}
	drop, messages, err = m.db.CreateDrop(ctx, jetID, lastSlotPulse.Pulse.PulseNumber, prevDrop.Hash)
	if err != nil {
		return nil, nil, nil, err
	}
	err = m.db.SetDrop(ctx, jetID, drop)
	if err != nil {
		return nil, nil, nil, err
	}

	dropSerialized, err = jet.Encode(drop)
	if err != nil {
		return nil, nil, nil, err
	}

	return
}

func (m *PulseManager) processDrop(
	ctx context.Context,
	jetID core.RecordID,
	lastSlotPulse *storage.Pulse,
	currentSlotPulse *core.Pulse,
	dropSerialized []byte,
	messages [][]byte,
) error {
	msg := &message.JetDrop{
		JetID:       jetID,
		Drop:        dropSerialized,
		Messages:    messages,
		PulseNumber: *lastSlotPulse.Prev,
	}
	_, err := m.Bus.Send(ctx, msg, *currentSlotPulse, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *PulseManager) processRecentObjects(
	ctx context.Context,
	jetID core.RecordID,
	previousSlotPulse *storage.Pulse,
	currentSlotPulse *core.Pulse,
	drop *jet.JetDrop,
	dropSerialized []byte,
) error {
	logger := inslogger.FromContext(ctx)
	recentStorage := m.RecentStorageProvider.GetStorage(core.TODOJetID)
	recentStorage.ClearZeroTTLObjects()
	recentObjectsIds := recentStorage.GetObjects()
	pendingRequestsIds := recentStorage.GetRequests()
	defer recentStorage.ClearObjects()

	recentObjects := map[core.RecordID]*message.HotIndex{}
	pendingRequests := map[core.RecordID][]byte{}

	for id, ttl := range recentObjectsIds {
		lifeline, err := m.db.GetObjectIndex(ctx, jetID, &id, false)
		if err != nil {
			logger.Error(err)
			continue
		}
		encoded, err := index.EncodeObjectLifeline(lifeline)
		if err != nil {
			logger.Error(err)
			continue
		}
		recentObjects[id] = &message.HotIndex{
			TTL:   ttl,
			Index: encoded,
		}

		if !recentStorage.IsMine(id) {
			err := m.db.RemoveObjectIndex(ctx, jetID, &id)
			if err != nil {
				logger.Error(err)
				return err
			}
		}
	}

	for _, id := range pendingRequestsIds {
		pendingRecord, err := m.db.GetRecord(ctx, jetID, &id)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
			continue
		}
		pendingRequests[id] = record.SerializeRecord(pendingRecord)
	}

	msg := &message.HotData{
		Drop:            *drop,
		PulseNumber:     previousSlotPulse.Pulse.PulseNumber,
		RecentObjects:   recentObjects,
		PendingRequests: pendingRequests,
	}
	_, err := m.Bus.Send(ctx, msg, *currentSlotPulse, nil)
	if err != nil {
		return err
	}
	return nil
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(ctx context.Context, pulse core.Pulse, dry bool) error {
	// Ensure this does not execute in parallel.
	m.setLock.Lock()
	defer m.setLock.Unlock()
	if m.stopped {
		return errors.New("can't call Set method on PulseManager after stop")
	}

	var err error
	m.GIL.Acquire(ctx)

	// swap pulse
	m.currentPulse = pulse

	lastSlotPulse, err := m.db.GetLatestPulse(ctx)
	if err != nil {
		return errors.Wrap(err, "call of GetLatestPulseNumber failed")
	}

	// swap active nodes
	m.ActiveListSwapper.MoveSyncToActive()
	if !dry {
		if err := m.db.AddPulse(ctx, pulse); err != nil {
			return errors.Wrap(err, "call of AddPulse failed")
		}
		err = m.db.SetActiveNodes(pulse.PulseNumber, m.NodeNet.GetActiveNodes())
		if err != nil {
			return errors.Wrap(err, "call of SetActiveNodes failed")
		}
	}

	m.GIL.Release(ctx)

	if dry {
		return nil
	}

	// Run only on material executor.
	// execute only on material executor
	// TODO: do as much as possible async.
	if m.NodeNet.GetOrigin().Role() == core.StaticRoleLightMaterial {
		err = m.dropAllJets(ctx, lastSlotPulse)
		if err != nil {
			return err
		}
		if m.enableSync {
			err := m.AddPulseToSyncClients(ctx, lastSlotPulse.Pulse.PulseNumber)
			if err != nil {
				return err
			}
		}
	}

	return m.LR.OnPulse(ctx, pulse)
}

// AddPulseToSyncClients add pulse number to all sync clients in pool.
func (m *PulseManager) AddPulseToSyncClients(ctx context.Context, pn core.PulseNumber) error {
	// get all jets with drops (required sync)
	allJets, err := m.db.GetJets(ctx)
	if err != nil {
		return err
	}
	for jetID := range allJets {
		_, err := m.db.GetDrop(ctx, jetID, pn)
		if err == nil {
			m.syncClientsPool.AddPulsesToSyncClient(ctx, m, jetID, pn)
		}
	}
	return nil
}

// Start starts pulse manager, spawns replication goroutine under a hood.
func (m *PulseManager) Start(ctx context.Context) error {
	if m.enableSync {
		err := m.initJetSyncState(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stops PulseManager. Waits replication goroutine is done.
func (m *PulseManager) Stop(ctx context.Context) error {
	// There should not to be any Set call after Stop call
	m.setLock.Lock()
	m.stopped = true
	m.setLock.Unlock()

	if m.enableSync {
		inslogger.FromContext(ctx).Info("waiting finish of heavy replication client...")
		m.syncClientsPool.Stop(ctx)
	}
	return nil
}
