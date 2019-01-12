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
	"fmt"
	"sync"
	"time"

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
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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
	tree, err := m.db.CloneJetTree(ctx, currentPulse.PulseNumber, newPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to clone jet tree into a new pulse")
	}

	jetIDs := tree.LeafIDs()

	var g errgroup.Group
	for _, jetID := range jetIDs {
		jetID := jetID

		executor, err := m.JetCoordinator.LightExecutorForJet(ctx, jetID, currentPulse.PulseNumber)
		if err != nil {
			return err
		}
		if *executor != m.JetCoordinator.Me() {
			return nil
		}

		g.Go(func() error {
			drop, dropSerialized, _, err := m.createDrop(ctx, jetID, prevPulseNumber, currentPulse.PulseNumber)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", currentPulse.PulseNumber)
			}

			msg, err := m.getExecutorHotData(
				ctx, jetID, newPulse.PulseNumber, drop, dropSerialized)
			if err != nil {
				return errors.Wrapf(err, "getExecutorData failed for jet id %v", jetID)
			}

			nextExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, jetID, newPulse.PulseNumber)
			if err != nil {
				return err
			}
			err = m.sendExecutorData(ctx, currentPulse, newPulse, jetID, msg, nextExecutor)
			if err != nil {
				return err
			}

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

	// TODO: @andreyromancev. 12.01.19. Uncomment when heavy is ready.
	// untilPN := currentPulse.PulseNumber - m.options.storeLightPulses
	// for jetID := range jetIDs {
	// 	replicated, err := m.db.GetReplicatedPulse(ctx, jetID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if untilPN >= replicated {
	// 		inslogger.FromContext(ctx).Errorf(
	// 			"light cleanup aborted (remove from: %v, replicated: %v)",
	// 			untilPN,
	// 			replicated,
	// 		)
	// 		return nil
	// 	}
	// 	if _, err := m.db.RemoveJetIndexesUntil(ctx, jetID, untilPN); err != nil {
	// 		return err
	// 	}
	// }
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
	var prevDrop *jet.JetDrop
	prevDrop, err = m.db.GetDrop(ctx, jetID, prevPulse)
	if err == storage.ErrNotFound {
		prevDrop, err = m.db.GetDrop(ctx, jet.Parent(jetID), prevPulse)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "[ createDrop ] failed to find parent")
		}
		err = nil
	}
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't GetDrop")
	}
	drop, messages, dropSize, err := m.db.CreateDrop(ctx, jetID, currentPulse, prevDrop.Hash)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't CreateDrop")
	}
	err = m.db.SetDrop(ctx, jetID, drop)
	if err == storage.ErrOverride {
		err = nil
	}
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
	// TODO: @andreyromancev. 12.01.19. Uncomment to check if this doesn't delete indexes it should not.
	// recentStorage.ClearZeroTTLObjects()
	// defer recentStorage.ClearObjects()
	recentObjectsIds := recentStorage.GetObjects()

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
		fmt.Println("[send id] ", id)
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
		DropJet:            jetID,
		PulseNumber:        pulse,
		RecentObjects:      recentObjects,
		PendingRequests:    pendingRequests,
		JetDropSizeHistory: dropSizeHistory,
	}
	return msg, nil
}

// TODO: @andreyromancev. 12.01.19. Remove when dynamic split is working.
var split = true

func (m *PulseManager) sendExecutorData(
	ctx context.Context,
	currentPulse, newPulse *core.Pulse,
	jetID core.RecordID,
	msg *message.HotData,
	receiver *core.RecordRef,
) error {
	// TODO: @andreyromancev. 12.01.19. Uncomment when split checking is ready.
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

	if split {
		split = false
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
		fmt.Printf(
			"[send hot] dropPulse: %v, dropJet: %v, jet: %v",
			leftMsg.Drop.Pulse,
			leftMsg.DropJet.JetIDString(),
			leftMsg.Jet.Record().JetIDString(),
		)
		fmt.Println("")
		_, err = m.Bus.Send(ctx, &leftMsg, nil)
		if err != nil {
			return errors.Wrap(err, "failed to send executor data")
		}
		fmt.Printf(
			"[send hot] dropPulse: %v, dropJet: %v, jet: %v",
			rightMsg.Drop.Pulse,
			rightMsg.DropJet.JetIDString(),
			rightMsg.Jet.Record().JetIDString(),
		)
		fmt.Println("")
		_, err = m.Bus.Send(ctx, &rightMsg, nil)
		if err != nil {
			return errors.Wrap(err, "failed to send executor data")
		}
	} else {
		msg.Jet = *core.NewRecordRef(core.DomainID, jetID)
		_, err := m.Bus.Send(ctx, msg, &core.MessageSendOptions{Receiver: receiver})
		if err != nil {
			return errors.Wrap(err, "failed to send executor data")
		}
	}

	return nil
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(ctx context.Context, newPulse core.Pulse, persist bool) error {
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

	fmt.Printf(
		"Received pulse %v, current: %v, time: %v",
		newPulse.PulseNumber,
		currentPulse.PulseNumber,
		time.Now(),
	)
	fmt.Println()

	// swap pulse
	m.currentPulse = newPulse

	// swap active nodes
	// TODO: fix network consensus and uncomment this (after NETD18-74)
	// m.ActiveListSwapper.MoveSyncToActive()
	fmt.Printf("Persist for pulse: %v is %v\n", newPulse.PulseNumber, persist)
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

	if !persist {
		return nil
	}

	// Run only on material executor.
	// execute only on material executor
	// TODO: do as much as possible async.
	if m.NodeNet.GetOrigin().Role() == core.StaticRoleLightMaterial {
		err = m.processEndPulse(ctx, prevPulseNumber, &currentPulse, &newPulse)
		if err != nil {
			fmt.Println("process end pulse failed: ", err)
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

	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
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

	return m.restoreGenesisRecentObjects(ctx)
}

func (m *PulseManager) restoreGenesisRecentObjects(ctx context.Context) error {
	jetID := *jet.NewID(0, nil)
	recent := m.RecentStorageProvider.GetStorage(jetID)

	return m.db.IterateIndexIDs(ctx, jetID, func(id core.RecordID) error {
		if id.Pulse() == core.FirstPulseNumber {
			recent.AddObject(id)
			fmt.Println("[restored] id ", id)
		}
		return nil
	})
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
