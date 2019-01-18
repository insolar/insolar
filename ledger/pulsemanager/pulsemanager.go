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
	"github.com/insolar/insolar/core/reply"
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
	LR                            core.LogicRunner                   `inject:""`
	Bus                           core.MessageBus                    `inject:""`
	NodeNet                       core.NodeNetwork                   `inject:""`
	JetCoordinator                core.JetCoordinator                `inject:""`
	GIL                           core.GlobalInsolarLock             `inject:""`
	CryptographyService           core.CryptographyService           `inject:""`
	PlatformCryptographyScheme    core.PlatformCryptographyScheme    `inject:""`
	RecentStorageProvider         recentstorage.Provider             `inject:""`
	ActiveListSwapper             ActiveListSwapper                  `inject:""`
	PulseStorage                  pulseStoragePm                     `inject:""`
	ArtifactManagerMessageHandler core.ArtifactManagerMessageHandler `inject:""`
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

type jetInfo struct {
	id       core.RecordID
	mineNext bool
	left     *jetInfo
	right    *jetInfo
}

// TODO: @andreyromancev. 15.01.19. Just store ledger configuration in PM. This is not required.
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
	jets []jetInfo,
	prevPulseNumber core.PulseNumber,
	currentPulse, newPulse *core.Pulse,
) error {
	var g errgroup.Group
	for _, i := range jets {
		info := i
		g.Go(func() error {
			drop, dropSerialized, _, err := m.createDrop(ctx, info.id, prevPulseNumber, currentPulse.PulseNumber)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", currentPulse.PulseNumber)
			}

			msg, err := m.getExecutorHotData(
				ctx, info.id, newPulse.PulseNumber, drop, dropSerialized)
			if err != nil {
				return errors.Wrapf(err, "getExecutorData failed for jet id %v", info.id)
			}

			if info.left == nil && info.right == nil {
				fmt.Printf("No split. jet: %v, mine next: %v", info.id, info.mineNext)

				// TODO: @andreyromancev. 12.01.19. uncomment when heavy ready.
				// m.RecentStorageProvider.GetStorage(info.id).ClearZeroTTLObjects()

				// No split happened.
				if !info.mineNext {
					msg.Jet = *core.NewRecordRef(core.DomainID, info.id)
					genericRep, err := m.Bus.Send(ctx, msg, nil)
					if err != nil {
						return errors.Wrap(err, "failed to send executor data")
					}
					if rep, ok := genericRep.(*reply.OK); !ok {
						return fmt.Errorf("unexpected reply: %#v", rep)
					}
					fmt.Printf("sent drop. pulse: %v, jet: %v\n", msg.Drop.Pulse, msg.DropJet.JetIDString())
				}
			} else {
				fmt.Printf("Split. jet: %v, left mine next: %v", info.id, info.left.mineNext)
				// Split happened.

				// TODO: @andreyromancev. 12.01.19. uncomment when heavy ready.
				// m.RecentStorageProvider.GetStorage(info.left.id).ClearZeroTTLObjects()
				// m.RecentStorageProvider.GetStorage(info.right.id).ClearZeroTTLObjects()

				if !info.left.mineNext {
					leftMsg := msg
					leftMsg.Jet = *core.NewRecordRef(core.DomainID, info.left.id)
					genericRep, err := m.Bus.Send(ctx, leftMsg, nil)
					if err != nil {
						return errors.Wrap(err, "failed to send executor data")
					}
					if rep, ok := genericRep.(*reply.OK); !ok {
						return fmt.Errorf("unexpected reply: %#v", rep)
					}
					fmt.Printf("sent drop. pulse: %v, jet: %v\n", msg.Drop.Pulse, msg.DropJet.JetIDString())
				}
				fmt.Printf("Split. jet: %v, right mine next: %v", info.id, info.right.mineNext)
				if !info.right.mineNext {
					rightMsg := msg
					rightMsg.Jet = *core.NewRecordRef(core.DomainID, info.right.id)

					genericRep, err := m.Bus.Send(ctx, rightMsg, nil)
					if err != nil {
						return errors.Wrap(err, "failed to send executor data")
					}
					if rep, ok := genericRep.(*reply.OK); !ok {
						return fmt.Errorf("unexpected reply: %#v", rep)
					}
					fmt.Printf("sent drop. pulse: %v, jet: %v\n", msg.Drop.Pulse, msg.DropJet.JetIDString())
				}
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
	err := g.Wait()
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
		parentJet := jet.Parent(jetID)
		fmt.Printf(
			"failed to fetch jet. pulse: %v, current jet: %v, parent jet: %v \n",
			prevPulse,
			jetID.JetIDString(),
			parentJet.String(),
		)
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
	fmt.Printf("saved drop. pulse: %v, jet: %v\n", drop.Pulse, jetID.JetIDString())

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
		fmt.Printf("[send id] %v\n", id.String())
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

func (m *PulseManager) processJets(ctx context.Context, currentPulse, newPulse core.PulseNumber) ([]jetInfo, error) {
	tree, err := m.db.CloneJetTree(ctx, currentPulse, newPulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to clone jet tree into a new pulse")
	}

	if m.NodeNet.GetOrigin().Role() != core.StaticRoleLightMaterial {
		return nil, nil
	}

	var results []jetInfo
	jetIDs := tree.LeafIDs()
	me := m.JetCoordinator.Me()
	for _, jetID := range jetIDs {
		fmt.Printf("I processed. jet: %v\n", jetID.JetIDString())
		executor, err := m.JetCoordinator.LightExecutorForJet(ctx, jetID, currentPulse)
		if err != nil {
			return nil, err
		}
		if *executor != me {
			continue
		}

		fmt.Printf("I am executor. jet: %v\n", jetID.JetIDString())

		info := jetInfo{id: jetID}
		if split {
			split = false

			leftJetID, rightJetID, err := m.db.SplitJetTree(
				ctx,
				newPulse,
				jetID,
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to split jet tree")
			}
			err = m.db.AddJets(ctx, *leftJetID, *rightJetID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to add jets")
			}
			// Set actual because we are the last executor for jet.
			err = m.db.UpdateJetTree(ctx, newPulse, true, *leftJetID, *rightJetID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update tree")
			}

			info.left = &jetInfo{id: *leftJetID}
			info.right = &jetInfo{id: *rightJetID}
			nextLeftExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, *leftJetID, newPulse)
			if err != nil {
				return nil, err
			}
			if *nextLeftExecutor == me {
				info.left.mineNext = true
				err := m.rewriteHotData(ctx, jetID, *leftJetID)
				if err != nil {
					return nil, err
				}
				fmt.Printf("I am next executor for left. jet: %v\n", info.left.id.JetIDString())

			}
			nextRightExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, *rightJetID, newPulse)
			if err != nil {
				return nil, err
			}
			if *nextRightExecutor == me {
				info.right.mineNext = true
				err := m.rewriteHotData(ctx, jetID, *rightJetID)
				if err != nil {
					return nil, err
				}
				fmt.Printf("I am next executor for right. jet: %v\n", info.right.id.JetIDString())
			}

			inslogger.FromContext(ctx).Debugf(
				"SPLIT HAPPENED parent: %v, left: %v, right: %v\n",
				jetID.JetIDString(),
				leftJetID.JetIDString(),
				rightJetID.JetIDString(),
			)
		} else {
			// Set actual because we are the last executor for jet.
			err = m.db.UpdateJetTree(ctx, newPulse, true, jetID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update tree")
			}
			nextExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, jetID, newPulse)
			if err != nil {
				return nil, err
			}
			if *nextExecutor == me {
				info.mineNext = true
				fmt.Printf("I am next executor. jet: %v\n", info.id.JetIDString())
			}
		}
		results = append(results, info)
	}

	return results, nil
}

func (m *PulseManager) rewriteHotData(ctx context.Context, fromJetID, toJetID core.RecordID) error {
	recentStorage := m.RecentStorageProvider.GetStorage(fromJetID)

	for id := range recentStorage.GetObjects() {
		idx, err := m.db.GetObjectIndex(ctx, fromJetID, &id, false)
		if err != nil {
			return errors.Wrap(err, "failed to rewrite index")
		}
		err = m.db.SetObjectIndex(ctx, toJetID, &id, idx)
		if err != nil {
			return errors.Wrap(err, "failed to rewrite index")
		}
	}

	for _, requests := range recentStorage.GetRequests() {
		for fromReqID := range requests {
			request, err := m.db.GetRecord(ctx, fromJetID, &fromReqID)
			if err != nil {
				return errors.Wrap(err, "failed to rewrite pending request")
			}
			toReqID, err := m.db.SetRecord(ctx, toJetID, fromReqID.Pulse(), request)
			if err == storage.ErrOverride {
				continue
			}
			if err != nil {
				return errors.Wrap(err, "failed to rewrite pending request")
			}
			if !fromReqID.Equal(toReqID) {
				return errors.New("failed to rewrite pending request (wrong ID generated)")
			}
		}
	}

	m.RecentStorageProvider.CloneStorage(fromJetID, toJetID)

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

	jets, err := m.processJets(ctx, currentPulse.PulseNumber, newPulse.PulseNumber)
	if err != nil {
		m.GIL.Release(ctx)
		return errors.Wrap(err, "failed to process jets")
	}

	m.prepareArtifactManagerMessageHandlerForNextPulse(ctx, newPulse, jets)

	m.GIL.Release(ctx)

	if !persist {
		return nil
	}

	// Run only on material executor.
	// execute only on material executor
	// TODO: do as much as possible async.
	if m.NodeNet.GetOrigin().Role() == core.StaticRoleLightMaterial {
		err = m.processEndPulse(ctx, jets, prevPulseNumber, &currentPulse, &newPulse)
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

	fmt.Printf(
		"Finished pulse %v, current: %v, time: %v",
		newPulse.PulseNumber,
		currentPulse.PulseNumber,
		time.Now(),
	)
	fmt.Println()

	// TODO: @andreyromancev. 12.01.19. uncomment when heavy ready.
	// m.postProcessJets(ctx, newPulse, jets)

	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
	}

	return m.LR.OnPulse(ctx, newPulse)
}

func (m *PulseManager) postProcessJets(ctx context.Context, newPulse core.Pulse, jets []jetInfo) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[postProcessJets] post-process jets, pulse number - %v", newPulse.PulseNumber)

	for _, jetInfo := range jets {
		if jetInfo.left == nil && jetInfo.right == nil {
			// No split happened.
			if !jetInfo.mineNext {
				logger.Debugf("[postProcessJets] clear recent storage for root jet - %v, pulse - %v", jetInfo.id, newPulse.PulseNumber)
				m.RecentStorageProvider.GetStorage(jetInfo.id).ClearObjects()
			}
		} else {
			// Split happened.
			if !jetInfo.left.mineNext {
				logger.Debugf("[postProcessJets] clear recent storage for left jet - %v, pulse - %v", jetInfo.left.id, newPulse.PulseNumber)
				m.RecentStorageProvider.GetStorage(jetInfo.left.id).ClearObjects()
			}
			if !jetInfo.right.mineNext {
				logger.Debugf("[postProcessJets] clear recent storage for right jet - %v, pulse - %v", jetInfo.right.id, newPulse.PulseNumber)
				m.RecentStorageProvider.GetStorage(jetInfo.right.id).ClearObjects()
			}
		}
	}
}

func (m *PulseManager) prepareArtifactManagerMessageHandlerForNextPulse(ctx context.Context, newPulse core.Pulse, jets []jetInfo) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] close breakers my jets for the next pulse - %v", newPulse.PulseNumber)

	m.ArtifactManagerMessageHandler.ResetEarlyRequestCircuitBreaker(ctx)

	for _, jetInfo := range jets {

		if jetInfo.left == nil && jetInfo.right == nil {
			// No split happened.
			if jetInfo.mineNext {
				logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] fetch jetInfo root %v, pulse - %v", jetInfo.id.JetIDString(), newPulse.PulseNumber)
				m.ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet(ctx, jetInfo.id)
			}
		} else {
			// Split happened.
			if jetInfo.left.mineNext {
				logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] fetch jetInfo left %v, pulse - %v", jetInfo.left.id.JetIDString(), newPulse.PulseNumber)
				m.ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet(ctx, jetInfo.left.id)
			}
			if jetInfo.right.mineNext {
				logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] fetch jetInfo right %v, pulse - %v", jetInfo.right.id.JetIDString(), newPulse.PulseNumber)
				m.ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet(ctx, jetInfo.right.id)
			}
		}
	}
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
			fmt.Printf("[restored] id %v \n", id.String())
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
