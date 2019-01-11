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
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavyclient"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/pulsemanager.ActiveListSwapper -o ../../testutils -s _mock.go
type ActiveListSwapper interface {
	MoveSyncToActive()
}

// PulseManager implements core.PulseManager.
type PulseManager struct {
	LR                         core.LogicRunner                `inject:""`
	Bus                        core.MessageBus                 `inject:""`
	NodeNet                    core.NodeNetwork                `inject:""`
	JetCoordinator             core.JetCoordinator             `inject:""`
	GIL                        core.GlobalInsolarLock          `inject:""`
	CryptographyService        core.CryptographyService        `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	RecentStorageProvider      recentstorage.Provider          `inject:""`
	ActiveListSwapper          ActiveListSwapper               `inject:""`
	PulseStorage               pulseStoragePm                  `inject:""`
	// TODO: move clients pool to component - @nordicdyno - 18.Dec.2018
	syncClientsPool *heavyclient.Pool

	currentPulse core.Pulse

	// internal stuff
	db *storage.DB
	// setLock locks Set method call.
	setLock sync.RWMutex
	// saves PM stopping mode
	stopped bool

	// stores pulse manager options
	options pmOptions
}

type pmOptions struct {
	enableSync       bool
	splitThreshold   uint64
	dropHistorySize  int
	storeLightPulses core.PulseNumber
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(db *storage.DB, conf configuration.Ledger) *PulseManager {
	pmconf := conf.PulseManager
	heavySyncPool := heavyclient.NewPool(
		db,
		heavyclient.Options{
			SyncMessageLimit: pmconf.HeavySyncMessageLimit,
			PulsesDeltaLimit: conf.LightChainLimit,
		},
	)
	pm := &PulseManager{
		db:           db,
		currentPulse: *core.GenesisPulse,
		options: pmOptions{
			enableSync:       pmconf.HeavySyncEnabled,
			splitThreshold:   pmconf.SplitThreshold,
			dropHistorySize:  conf.JetSizesHistoryDepth,
			storeLightPulses: conf.LightChainLimit,
		},
		syncClientsPool: heavySyncPool,
	}

	// TODO: untie this circular dependency after moving sync client to separate component - 17.Dec.2018 @nordicdyno
	return pm
}

func (m *PulseManager) processEndPulse(
	ctx context.Context,
	prevPulseNumber core.PulseNumber,
	currentPulse, newPulse *core.Pulse,
) error {
	err := m.db.CloneJetTree(ctx, currentPulse.PulseNumber, newPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to clone jet tree into a new pulse")
	}

	jetIDs, err := m.db.GetJets(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get jets from storage")
	}
	var g errgroup.Group
	for jetID := range jetIDs {
		jetID := jetID
		g.Go(func() error {
			drop, dropSerialized, _, err := m.createDrop(ctx, jetID, prevPulseNumber, currentPulse.PulseNumber)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", currentPulse.PulseNumber)
			}

			inslogger.FromContext(ctx).Debugf("[processEndPulse] before getExecutorHotData - %v", time.Now())
			msg, hotRecordsError := m.getExecutorHotData(
				ctx, jetID, currentPulse.PulseNumber, drop, dropSerialized)
			if hotRecordsError != nil {
				return errors.Wrapf(err, "getExecutorData failed for jet id %v", jetID)
			}
			inslogger.FromContext(ctx).Debugf("[processEndPulse] after getExecutorHotData - %v", time.Now())

			inslogger.FromContext(ctx).Debugf("[processEndPulse] before sendExecutorData - %v", time.Now())
			sendError := m.sendExecutorData(ctx, currentPulse, newPulse, jetID, msg)
			if sendError != nil {
				return err
			}
			inslogger.FromContext(ctx).Debugf("[processEndPulse] after sendExecutorData - %v", time.Now())

			// FIXME: @andreyromancev. 09.01.2019. Temporary disabled validation. Uncomment when jet split works properly.
			// dropErr := m.processDrop(ctx, jetID, currentPulse, dropSerialized, messages)
			// if dropErr != nil {
			// 	return errors.Wrap(dropErr, "processDrop failed")
			// }

			// TODO: @andreyromancev. 20.12.18. uncomment me when pending notifications required.
			// m.sendAbandonedRequests(ctx, newPulse, jetID)

			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return errors.Wrap(err, "got error on jets sync")
	}

	// TODO: maybe move cleanup in the above cycle or process removal in separate job - 20.Dec.2018 @nordicdyno
	untilPN := currentPulse.PulseNumber - m.options.storeLightPulses
	for jetID := range jetIDs {
		replicated, err := m.db.GetReplicatedPulse(ctx, jetID)
		if err != nil {
			return err
		}
		if untilPN >= replicated {
			inslogger.FromContext(ctx).Errorf(
				"light cleanup aborted (remove from: %v, replicated: %v)",
				untilPN,
				replicated,
			)
			return nil
		}
		if _, err := m.db.RemoveJetIndexesUntil(ctx, jetID, untilPN); err != nil {
			return err
		}
	}
	return nil
}

// TODO: @andreyromancev. 20.12.18. uncomment me when pending notifications required.
// func (m *PulseManager) sendAbandonedRequests(ctx context.Context, pulse *core.Pulse, jetID core.RecordID) {
// 	pendingRequests := m.RecentStorageProvider.GetStorage(jetID).GetRequests()
// 	wg := sync.WaitGroup{}
// 	wg.Add(len(pendingRequests))
// 	for objID, requests := range pendingRequests {
// 		go func(object core.RecordID, objectRequests map[core.RecordID]struct{}) {
// 			defer wg.Done()
//
// 			var toSend []core.RecordID
// 			for reqID := range objectRequests {
// 				toSend = append(toSend, reqID)
// 			}
// 			rep, err := m.Bus.Send(ctx, &message.AbandonedRequestsNotification{
// 				Object:   object,
// 				Requests: toSend,
// 			}, *pulse, nil)
// 			if err != nil {
// 				inslogger.FromContext(ctx).Error("failed to notify about pending requests")
// 				return
// 			}
// 			if _, ok := rep.(*reply.OK); !ok {
// 				inslogger.FromContext(ctx).Error("received unexpected reply on pending notification")
// 			}
// 		}(objID, requests)
// 	}
//
// 	wg.Wait()
// }

func (m *PulseManager) createDrop(
	ctx context.Context,
	jetID core.RecordID,
	prevPulse, currentPulse core.PulseNumber,
) (
	drop *jet.JetDrop,
	dropSerialized []byte,
	messages [][]byte,
	err error,
) {
	prevDrop, err := m.db.GetDrop(ctx, jetID, prevPulse)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't GetDrop")
	}
	drop, messages, dropSize, err := m.db.CreateDrop(ctx, jetID, currentPulse, prevDrop.Hash)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't CreateDrop")
	}
	err = m.db.SetDrop(ctx, jetID, drop)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't SetDrop")
	}

	dropSerialized, err = jet.Encode(drop)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't Encode")
	}

	dropSizeData := &jet.DropSize{
		JetID:    jetID,
		PulseNo:  currentPulse,
		DropSize: dropSize,
	}
	hasher := m.PlatformCryptographyScheme.IntegrityHasher()
	_, err = dropSizeData.WriteHashData(hasher)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't WriteHashData")
	}
	signature, err := m.CryptographyService.Sign(hasher.Sum(nil))
	dropSizeData.Signature = signature.Bytes()

	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't Sign")
	}

	err = m.db.AddDropSize(ctx, dropSizeData)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't AddDropSize")
	}

	return
}

func (m *PulseManager) processDrop(
	ctx context.Context,
	jetID core.RecordID,
	pulse *core.Pulse,
	dropSerialized []byte,
	messages [][]byte,
) error {
	msg := &message.JetDrop{
		JetID:       jetID,
		Drop:        dropSerialized,
		Messages:    messages,
		PulseNumber: pulse.PulseNumber,
	}
	_, err := m.Bus.Send(ctx, msg, nil)
	if err != nil {
		return err
	}
	return nil
}

func (m *PulseManager) getExecutorHotData(
	ctx context.Context,
	jetID core.RecordID,
	pulse core.PulseNumber,
	drop *jet.JetDrop,
	dropSerialized []byte,
) (*message.HotData, error) {
	logger := inslogger.FromContext(ctx)
	recentStorage := m.RecentStorageProvider.GetStorage(jetID)
	recentStorage.ClearZeroTTLObjects()
	recentObjectsIds := recentStorage.GetObjects()
	defer recentStorage.ClearObjects()

	recentObjects := map[core.RecordID]*message.HotIndex{}
	pendingRequests := map[core.RecordID]map[core.RecordID][]byte{}

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

	}

	for objID, requests := range recentStorage.GetRequests() {
		for reqID := range requests {
			pendingRecord, err := m.db.GetRecord(ctx, jetID, &reqID)
			if err != nil {
				inslogger.FromContext(ctx).Error(err)
				continue
			}
			if _, ok := pendingRequests[objID]; !ok {
				pendingRequests[objID] = map[core.RecordID][]byte{}
			}
			pendingRequests[objID][reqID] = record.SerializeRecord(pendingRecord)
		}
	}

	dropSizeHistory, err := m.db.GetDropSizeHistory(ctx, jetID)
	if err != nil {
		return nil, errors.Wrap(err, "[ processRecentObjects ] Can't GetDropSizeHistory")
	}

	msg := &message.HotData{
		Drop:               *drop,
		PulseNumber:        pulse,
		RecentObjects:      recentObjects,
		PendingRequests:    pendingRequests,
		JetDropSizeHistory: dropSizeHistory,
	}
	return msg, nil
}

func (m *PulseManager) sendExecutorData(
	ctx context.Context,
	currentPulse, newPulse *core.Pulse,
	jetID core.RecordID,
	msg *message.HotData,
) error {
	// shouldSplit := func() bool {
	// 	if len(msg.JetDropSizeHistory) < m.options.dropHistorySize {
	// 		return false
	// 	}
	// 	for _, info := range msg.JetDropSizeHistory {
	// 		if info.DropSize < m.options.splitThreshold {
	// 			return false
	// 		}
	// 	}
	// 	return true
	// }

	// FIXME: enable split
	// if shouldSplit() {
	if false {
		left, right, err := m.db.SplitJetTree(
			ctx,
			currentPulse.PulseNumber,
			newPulse.PulseNumber,
			jetID,
		)
		if err != nil {
			return errors.Wrap(err, "failed to split jet tree")
		}
		err = m.db.AddJets(ctx, *left, *right)
		if err != nil {
			return errors.Wrap(err, "failed to add jets")
		}
		leftMsg := *msg
		leftMsg.Jet = *core.NewRecordRef(core.DomainID, *left)
		rightMsg := *msg
		rightMsg.Jet = *core.NewRecordRef(core.DomainID, *right)
		_, err = m.Bus.Send(ctx, &leftMsg, nil)
		if err != nil {
			return errors.Wrap(err, "failed to send executor data")
		}
		_, err = m.Bus.Send(ctx, &rightMsg, nil)
		if err != nil {
			return errors.Wrap(err, "failed to send executor data")
		}
	} else {
		msg.Jet = *core.NewRecordRef(core.DomainID, jetID)
		inslogger.FromContext(ctx).Debugf("[sendExecutorData] before m.Bus.Send(ctx, msg, nil) - %v", time.Now())
		_, err := m.Bus.Send(ctx, msg, nil)
		if err != nil {
			return errors.Wrap(err, "failed to send executor data")
		}
		inslogger.FromContext(ctx).Debugf("[sendExecutorData] after m.Bus.Send(ctx, msg, nil) - %v", time.Now())
	}

	return nil
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(ctx context.Context, newPulse core.Pulse, persist bool) error {
	// Ensure this does not execute in parallel.
	m.setLock.Lock()
	defer m.setLock.Unlock()
	if m.stopped {
		return errors.New("can't call Set method on PulseManager after stop")
	}

	var err error
	m.GIL.Acquire(ctx)

	m.PulseStorage.Lock()

	// FIXME: @andreyromancev. 17.12.18. return core.Pulse here.
	storagePulse, err := m.db.GetLatestPulse(ctx)
	if err != nil {
		m.PulseStorage.Unlock()
		m.GIL.Release(ctx)
		return errors.Wrap(err, "call of GetLatestPulseNumber failed")
	}
	currentPulse := storagePulse.Pulse
	prevPulseNumber := *storagePulse.Prev

	// swap pulse
	m.currentPulse = newPulse

	// swap active nodes
	// TODO: fix network consensus and uncomment this (after NETD18-74)
	// m.ActiveListSwapper.MoveSyncToActive()
	if persist {
		if err := m.db.AddPulse(ctx, newPulse); err != nil {
			m.GIL.Release(ctx)
			m.PulseStorage.Unlock()
			return errors.Wrap(err, "call of AddPulse failed")
		}
		err = m.db.SetActiveNodes(newPulse.PulseNumber, m.NodeNet.GetActiveNodes())
		if err != nil {
			m.GIL.Release(ctx)
			m.PulseStorage.Unlock()
			return errors.Wrap(err, "call of SetActiveNodes failed")
		}
	}

	m.PulseStorage.Unlock()
	m.PulseStorage.Set(&newPulse)
	m.GIL.Release(ctx)

	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
	}

	if !persist {
		return nil
	}

	// Run only on material executor.
	// execute only on material executor
	// TODO: do as much as possible async.
	if m.NodeNet.GetOrigin().Role() == core.StaticRoleLightMaterial {
		err = m.processEndPulse(ctx, prevPulseNumber, &currentPulse, &newPulse)
		if err != nil {
			return err
		}
		if m.options.enableSync {
			err := m.AddPulseToSyncClients(ctx, storagePulse.Pulse.PulseNumber)
			if err != nil {
				return err
			}
			go m.sendTreeToHeavy(ctx, storagePulse.Pulse.PulseNumber)
		}
	}

	return m.LR.OnPulse(ctx, newPulse)
}

func (m *PulseManager) sendTreeToHeavy(ctx context.Context, pn core.PulseNumber) {
	jetTree, err := m.db.GetJetTree(ctx, pn)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
	_, err = m.Bus.Send(ctx, &message.HeavyJetTree{PulseNum: pn, JetTree: *jetTree}, nil)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
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
			m.syncClientsPool.AddPulsesToSyncClient(ctx, jetID, true, pn)
		}
	}
	return nil
}

// Start starts pulse manager, spawns replication goroutine under a hood.
func (m *PulseManager) Start(ctx context.Context) error {
	// FIXME: @andreyromancev. 21.12.18. Find a proper place for me. Somewhere at the genesis.
	err := m.db.SetActiveNodes(core.FirstPulseNumber, m.NodeNet.GetActiveNodes())
	if err != nil && err != storage.ErrOverride {
		return err
	}

	if m.options.enableSync {
		m.syncClientsPool.Bus = m.Bus
		m.syncClientsPool.PulseStorage = m.PulseStorage
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

	if m.options.enableSync {
		inslogger.FromContext(ctx).Info("waiting finish of heavy replication client...")
		m.syncClientsPool.Stop(ctx)
	}
	return nil
}
