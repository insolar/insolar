/*
 *    Copyright 2019 Insolar Technologies
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
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/artifactmanager"
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
	HotDataWaiter              artifactmanager.HotDataWaiter   `inject:""`
	JetStorage                 storage.JetStorage              `inject:""`
	DropStorage                storage.DropStorage             `inject:""`
	ObjectStorage              storage.ObjectStorage           `inject:""`
	NodeStorage                storage.NodeStorage             `inject:""`
	PulseTracker               storage.PulseTracker            `inject:""`
	ReplicaStorage             storage.ReplicaStorage          `inject:""`
	DBContext                  storage.DBContext               `inject:""`
	StorageCleaner             storage.Cleaner                 `inject:""`

	// TODO: move clients pool to component - @nordicdyno - 18.Dec.2018
	syncClientsPool *heavyclient.Pool

	currentPulse core.Pulse

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
	enableSync            bool
	splitThreshold        uint64
	dropHistorySize       int
	storeLightPulses      int
	heavySyncMessageLimit int
	lightChainLimit       int
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(conf configuration.Ledger) *PulseManager {
	pmconf := conf.PulseManager

	pm := &PulseManager{
		currentPulse: *core.GenesisPulse,
		options: pmOptions{
			enableSync:            pmconf.HeavySyncEnabled,
			splitThreshold:        pmconf.SplitThreshold,
			dropHistorySize:       conf.JetSizesHistoryDepth,
			storeLightPulses:      conf.LightChainLimit,
			heavySyncMessageLimit: pmconf.HeavySyncMessageLimit,
			lightChainLimit:       conf.LightChainLimit,
		},
	}
	return pm
}

func (m *PulseManager) processEndPulse(
	ctx context.Context,
	jets []jetInfo,
	prevPulseNumber core.PulseNumber,
	currentPulse, newPulse core.Pulse,
) error {
	var g errgroup.Group
	logger := inslogger.FromContext(ctx)
	ctx, span := instracer.StartSpan(ctx, "pulse.process_end")
	defer span.End()

	for _, i := range jets {
		info := i

		g.Go(func() error {
			drop, dropSerialized, _, err := m.createDrop(ctx, info.id, prevPulseNumber, currentPulse.PulseNumber)
			logger.Debugf("[jet]: %v create drop. Pulse: %v, Error: %s", info.id.DebugString(), currentPulse.PulseNumber, err)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", currentPulse.PulseNumber)
			}

			msg, err := m.getExecutorHotData(
				ctx, info.id, newPulse.PulseNumber, drop, dropSerialized,
			)
			if err != nil {
				return errors.Wrapf(err, "getExecutorData failed for jet id %v", info.id)
			}

			logger := inslogger.FromContext(ctx)
			sender := func(msg message.HotData, jetID core.RecordID) {
				ctx, span := instracer.StartSpan(ctx, "pulse.send_hot")
				defer span.End()
				msg.Jet = *core.NewRecordRef(core.DomainID, jetID)
				start := time.Now()
				genericRep, err := m.Bus.Send(ctx, &msg, nil)
				sendTime := time.Since(start)
				if sendTime > time.Second {
					logger.Debugf("[send] jet: %v, long send: %s. Success: %v", jetID.DebugString(), sendTime, err == nil)
				}
				if err != nil {
					logger.Debugf("[jet]: %v send hot. Pulse: %v, DropJet: %v, Error: %s", jetID.DebugString(), currentPulse.PulseNumber, msg.DropJet.DebugString(), err)
					return
				}
				if _, ok := genericRep.(*reply.OK); !ok {
					logger.Debugf("[jet]: %v send hot. Pulse: %v, DropJet: %v, Unexpected reply: %#v", jetID.DebugString(), currentPulse.PulseNumber, msg.DropJet.DebugString(), genericRep)
					return
				}
				logger.Debugf("[jet]: %v send hot. Pulse: %v, DropJet: %v, Success", jetID.DebugString(), currentPulse.PulseNumber, msg.DropJet.DebugString())
			}

			if info.left == nil && info.right == nil {
				// No split happened.
				if !info.mineNext {
					go sender(*msg, info.id)
				}
			} else {
				// Split happened.
				if !info.left.mineNext {
					go sender(*msg, info.left.id)
				}
				if !info.right.mineNext {
					go sender(*msg, info.right.id)
				}
			}

			requests := m.RecentStorageProvider.GetStorage(ctx, info.id).GetRequests()
			go func() {
				err := m.sendAbandonedRequests(
					ctx,
					newPulse,
					requests,
				)
				logger.Error(err)
			}()

			// FIXME: @andreyromancev. 09.01.2019. Temporary disabled validation. Uncomment when jet split works properly.
			// dropErr := m.processDrop(ctx, jetID, currentPulse, dropSerialized, messages)
			// if dropErr != nil {
			// 	return errors.Wrap(dropErr, "processDrop failed")
			// }

			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return errors.Wrap(err, "got error on jets sync")
	}

	return nil
}

func (m *PulseManager) sendAbandonedRequests(
	ctx context.Context,
	pulse core.Pulse,
	pendingRequests map[core.RecordID]map[core.RecordID]struct{},
) error {
	ctx, span := instracer.StartSpan(ctx, "pulse.sendAbandonedRequests")
	defer span.End()

	logger := inslogger.FromContext(ctx)
	currentDBPulse, err := m.PulseTracker.GetPulse(ctx, pulse.PulseNumber)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(pendingRequests))
	for objID, requests := range pendingRequests {
		go func(object core.RecordID, requests map[core.RecordID]struct{}) {
			defer wg.Done()
			for request := range requests {
				pulse, err := m.PulseTracker.GetPulse(ctx, request.Pulse())
				if err != nil {
					logger.Error("failed to notify about pending requests. failed to calculate pulse")
					return
				}

				if currentDBPulse.SerialNumber-pulse.SerialNumber < 2 {
					continue
				}

				rep, err := m.Bus.Send(
					ctx,
					&message.AbandonedRequestsNotification{
						Object: object,
					},
					nil,
				)
				if err != nil {
					logger.Error("failed to notify about pending requests")
					return
				}
				if _, ok := rep.(*reply.OK); !ok {
					logger.Error("received unexpected reply on pending notification")
				}

				return
			}
		}(objID, requests)
	}

	wg.Wait()
	return nil
}

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
	prevDrop, err = m.DropStorage.GetDrop(ctx, jetID, prevPulse)
	if err == storage.ErrNotFound {
		prevDrop, err = m.DropStorage.GetDrop(ctx, jet.Parent(jetID), prevPulse)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "[ createDrop ] failed to find parent")
		}
	} else if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't GetDrop")
	}

	drop, messages, dropSize, err := m.DropStorage.CreateDrop(ctx, jetID, currentPulse, prevDrop.Hash)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't CreateDrop")
	}
	err = m.DropStorage.SetDrop(ctx, jetID, drop)
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

	err = m.DropStorage.AddDropSize(ctx, dropSizeData)
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
	ctx, span := instracer.StartSpan(ctx, "pulse.prepare_hot_data")
	defer span.End()

	logger := inslogger.FromContext(ctx)
	recentStorage := m.RecentStorageProvider.GetStorage(ctx, jetID)
	recentObjectsIds := recentStorage.GetObjects()

	recentObjects := map[core.RecordID]*message.HotIndex{}
	pendingRequests := map[core.RecordID]map[core.RecordID][]byte{}

	for id, ttl := range recentObjectsIds {
		lifeline, err := m.ObjectStorage.GetObjectIndex(ctx, jetID, &id, false)
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
			pendingRecord, err := m.ObjectStorage.GetRecord(ctx, jetID, &reqID)
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

	dropSizeHistory, err := m.DropStorage.GetDropSizeHistory(ctx, jetID)
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
var splitCount = 5

func (m *PulseManager) processJets(ctx context.Context, currentPulse, newPulse core.PulseNumber) ([]jetInfo, error) {
	ctx, span := instracer.StartSpan(ctx, "jets.process")
	defer span.End()

	tree, err := m.JetStorage.CloneJetTree(ctx, currentPulse, newPulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to clone jet tree into a new pulse")
	}

	if m.NodeNet.GetOrigin().Role() != core.StaticRoleLightMaterial {
		return nil, nil
	}

	var results []jetInfo
	jetIDs := tree.LeafIDs()
	me := m.JetCoordinator.Me()
	logger := inslogger.FromContext(ctx)
	indexToSplit := rand.Intn(len(jetIDs))
	for i, jetID := range jetIDs {
		executor, err := m.JetCoordinator.LightExecutorForJet(ctx, jetID, currentPulse)
		if err != nil {
			return nil, err
		}
		imExecutor := *executor == me
		logger.Debugf("[jet]: %v process. Pulse: %v, Executor: %v", jetID.DebugString(), currentPulse, imExecutor)
		if !imExecutor {
			continue
		}

		m.RecentStorageProvider.GetStorage(ctx, jetID).DecreaseTTL(ctx)

		info := jetInfo{id: jetID}
		if indexToSplit == i && splitCount > 0 {
			splitCount--

			leftJetID, rightJetID, err := m.JetStorage.SplitJetTree(
				ctx,
				newPulse,
				jetID,
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to split jet tree")
			}
			err = m.JetStorage.AddJets(ctx, *leftJetID, *rightJetID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to add jets")
			}
			// Set actual because we are the last executor for jet.
			err = m.JetStorage.UpdateJetTree(ctx, newPulse, true, *leftJetID, *rightJetID)
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
				logger.Debugf("[jet]: %v rewrite hot left. Pulse: %v, Error: %s", info.left.id.DebugString(), currentPulse, err)
				if err != nil {
					return nil, err
				}
			}
			nextRightExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, *rightJetID, newPulse)
			if err != nil {
				return nil, err
			}
			if *nextRightExecutor == me {
				info.right.mineNext = true
				err := m.rewriteHotData(ctx, jetID, *rightJetID)
				logger.Debugf("[jet]: %v rewrite hot right. Pulse: %v, Error: %s", info.right.id.DebugString(), currentPulse, err)
				if err != nil {
					return nil, err
				}
			}

			logger.Debugf(
				"SPLIT HAPPENED parent: %v, left: %v, right: %v",
				jetID.DebugString(),
				leftJetID.DebugString(),
				rightJetID.DebugString(),
			)
		} else {
			// Set actual because we are the last executor for jet.
			err = m.JetStorage.UpdateJetTree(ctx, newPulse, true, jetID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update tree")
			}
			nextExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, jetID, newPulse)
			if err != nil {
				return nil, err
			}
			if *nextExecutor == me {
				info.mineNext = true
				logger.Debugf("[jet]: %v preserve hot. Pulse: %v", info.id.DebugString(), currentPulse)
			}
		}
		results = append(results, info)
	}

	return results, nil
}

func (m *PulseManager) rewriteHotData(ctx context.Context, fromJetID, toJetID core.RecordID) error {
	recentStorage := m.RecentStorageProvider.GetStorage(ctx, fromJetID)

	for id := range recentStorage.GetObjects() {
		idx, err := m.ObjectStorage.GetObjectIndex(ctx, fromJetID, &id, false)
		if err != nil {
			return errors.Wrap(err, "failed to rewrite index")
		}
		err = m.ObjectStorage.SetObjectIndex(ctx, toJetID, &id, idx)
		if err != nil {
			return errors.Wrap(err, "failed to rewrite index")
		}
	}

	for _, requests := range recentStorage.GetRequests() {
		for fromReqID := range requests {
			request, err := m.ObjectStorage.GetRecord(ctx, fromJetID, &fromReqID)
			if err != nil {
				return errors.Wrap(err, "failed to rewrite pending request")
			}
			toReqID, err := m.ObjectStorage.SetRecord(ctx, toJetID, fromReqID.Pulse(), request)
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

	inslogger.FromContext(ctx).Debugf("{LEAK} CloneStorage from - %v, to - %v", fromJetID, toJetID)
	m.RecentStorageProvider.CloneStorage(ctx, fromJetID, toJetID)

	return nil
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(ctx context.Context, newPulse core.Pulse, persist bool) error {
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

	jets, oldPulse, prevPulseNumber, err := m.setUnderGilSection(ctx, newPulse, persist)
	if err != nil {
		return err
	}

	if !persist {
		return nil
	}

	// Run only on material executor.
	// execute only on material executor
	// TODO: do as much as possible async.
	if m.NodeNet.GetOrigin().Role() == core.StaticRoleLightMaterial {
		err = m.processEndPulse(ctx, jets, *prevPulseNumber, *oldPulse, newPulse)
		if err != nil {
			return err
		}
		m.postProcessJets(ctx, newPulse, jets)
		m.addSync(ctx, jets, oldPulse.PulseNumber)
		go m.cleanLightData(ctx, newPulse)
	}

	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
	}

	return m.LR.OnPulse(ctx, newPulse)
}

func (m *PulseManager) setUnderGilSection(
	ctx context.Context, newPulse core.Pulse, persist bool,
) (
	[]jetInfo, *core.Pulse, *core.PulseNumber, error,
) {
	m.GIL.Acquire(ctx)
	ctx, span := instracer.StartSpan(ctx, "pulse.gil_locked")
	defer span.End()
	defer m.GIL.Release(ctx)

	m.PulseStorage.Lock()
	// FIXME: @andreyromancev. 17.12.18. return core.Pulse here.
	storagePulse, err := m.PulseTracker.GetLatestPulse(ctx)
	if err != nil {
		m.PulseStorage.Unlock()
		return nil, nil, nil, errors.Wrap(err, "call of GetLatestPulseNumber failed")
	}

	oldPulse := storagePulse.Pulse
	prevPulseNumber := storagePulse.Prev

	logger := inslogger.FromContext(ctx)
	logger.WithFields(map[string]interface{}{
		"new_pulse":     newPulse.PulseNumber,
		"current_pulse": oldPulse.PulseNumber,
		"persist":       persist,
	}).Debugf("received pulse")

	// swap pulse
	m.currentPulse = newPulse

	// swap active nodes
	// TODO: fix network consensus and uncomment this (after NETD18-74)
	// m.ActiveListSwapper.MoveSyncToActive()
	if persist {
		if err := m.PulseTracker.AddPulse(ctx, newPulse); err != nil {
			m.PulseStorage.Unlock()
			return nil, nil, nil, errors.Wrap(err, "call of AddPulse failed")
		}
		err = m.NodeStorage.SetActiveNodes(newPulse.PulseNumber, m.NodeNet.GetActiveNodes())
		if err != nil {
			m.PulseStorage.Unlock()
			return nil, nil, nil, errors.Wrap(err, "call of SetActiveNodes failed")
		}
	}

	m.PulseStorage.Set(&newPulse)
	m.PulseStorage.Unlock()

	var jets []jetInfo
	if persist {
		jets, err = m.processJets(ctx, oldPulse.PulseNumber, newPulse.PulseNumber)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to process jets")
		}
	}

	if m.NodeNet.GetOrigin().Role() == core.StaticRoleLightMaterial {
		m.prepareArtifactManagerMessageHandlerForNextPulse(ctx, newPulse, jets)
	}

	return jets, &oldPulse, prevPulseNumber, nil
}

func (m *PulseManager) addSync(ctx context.Context, jets []jetInfo, pulse core.PulseNumber) {
	ctx, span := instracer.StartSpan(ctx, "pulse.add_sync")
	defer span.End()

	if !m.options.enableSync {
		return
	}

	for _, jInfo := range jets {
		m.syncClientsPool.AddPulsesToSyncClient(ctx, jInfo.id, true, pulse)
	}
}

func (m *PulseManager) postProcessJets(ctx context.Context, newPulse core.Pulse, jets []jetInfo) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[postProcessJets] post-process jets, pulse number - %v", newPulse.PulseNumber)

	ctx, span := instracer.StartSpan(ctx, "jets.post_process")
	defer span.End()

	for _, jetInfo := range jets {
		if jetInfo.left == nil && jetInfo.right == nil {
			// No split happened.
			if !jetInfo.mineNext {
				logger.Debugf("[postProcessJets] clear recent storage for root jet - %v, pulse - %v", jetInfo.id, newPulse.PulseNumber)
				m.RecentStorageProvider.RemoveStorage(ctx, jetInfo.id)
			}
		} else {
			// Split happened.
			m.RecentStorageProvider.RemoveStorage(ctx, jetInfo.id)
			if !jetInfo.left.mineNext {
				logger.Debugf("[postProcessJets] clear recent storage for left jet - %v, pulse - %v", jetInfo.left.id, newPulse.PulseNumber)
				m.RecentStorageProvider.RemoveStorage(ctx, jetInfo.left.id)
			}
			if !jetInfo.right.mineNext {
				logger.Debugf("[postProcessJets] clear recent storage for right jet - %v, pulse - %v", jetInfo.right.id, newPulse.PulseNumber)
				m.RecentStorageProvider.RemoveStorage(ctx, jetInfo.right.id)
			}
		}
	}
}

func (m *PulseManager) cleanLightData(ctx context.Context, newPulse core.Pulse) {
	startSync := time.Now()
	inslog := inslogger.FromContext(ctx)
	ctx, span := instracer.StartSpan(ctx, "pulse.clean")
	defer func() {
		latency := time.Since(startSync)
		stats.Record(ctx, statCleanLatencyTotal.M(latency.Nanoseconds()/1e6))
		span.End()
		inslog.Infof("cleanLightData all time spend=%v", latency)
	}()

	delta := m.options.storeLightPulses

	pn := newPulse.PulseNumber
	func() {
		startLookup := time.Now()
		defer func() {
			inslog.Infof("cleanLightData pulse lookup time spend=%v", time.Since(startLookup))
		}()
		for i := 0; i <= delta; i++ {
			prevPulse, err := m.PulseTracker.GetPreviousPulse(ctx, pn)
			if err != nil {
				inslogger.FromContext(ctx).Errorf("Can't get previous Nth %v pulse by pulse number: %v", i, pn)
				return
			}

			pn = prevPulse.Pulse.PulseNumber
			if pn <= core.FirstPulseNumber {
				return
			}
		}
	}()

	m.NodeStorage.RemoveActiveNodesUntil(pn)
	m.JetStorage.DeleteJetTree(ctx, pn)

	err := m.syncClientsPool.LightCleanup(ctx, pn, m.RecentStorageProvider)
	if err != nil {
		inslogger.FromContext(ctx).Errorf(
			"Error on light cleanup, until pulse = %v, singlefligt err = %v", pn, err)
	}
}

func (m *PulseManager) prepareArtifactManagerMessageHandlerForNextPulse(ctx context.Context, newPulse core.Pulse, jets []jetInfo) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] close breakers my jets for the next pulse - %v", newPulse.PulseNumber)

	ctx, span := instracer.StartSpan(ctx, "early.close")
	defer span.End()

	m.HotDataWaiter.ThrowTimeout(ctx)

	for _, jetInfo := range jets {

		if jetInfo.left == nil && jetInfo.right == nil {
			// No split happened.
			if jetInfo.mineNext {
				logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] fetch jetInfo root %v, pulse - %v", jetInfo.id.DebugString(), newPulse.PulseNumber)
				m.HotDataWaiter.Unlock(ctx, jetInfo.id)
			}
		} else {
			// Split happened.
			if jetInfo.left.mineNext {
				logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] fetch jetInfo left %v, pulse - %v", jetInfo.left.id.DebugString(), newPulse.PulseNumber)
				m.HotDataWaiter.Unlock(ctx, jetInfo.left.id)
			}
			if jetInfo.right.mineNext {
				logger.Debugf("[breakermiddleware] [prepareHandlerForNextPulse] fetch jetInfo right %v, pulse - %v", jetInfo.right.id.DebugString(), newPulse.PulseNumber)
				m.HotDataWaiter.Unlock(ctx, jetInfo.right.id)
			}
		}
	}
}

// Start starts pulse manager, spawns replication goroutine under a hood.
func (m *PulseManager) Start(ctx context.Context) error {
	// FIXME: @andreyromancev. 21.12.18. Find a proper place for me. Somewhere at the genesis.
	err := m.NodeStorage.SetActiveNodes(core.FirstPulseNumber, m.NodeNet.GetActiveNodes())
	if err != nil && err != storage.ErrOverride {
		return err
	}

	if m.options.enableSync {
		heavySyncPool := heavyclient.NewPool(
			m.Bus,
			m.PulseStorage,
			m.PulseTracker,
			m.ReplicaStorage,
			m.StorageCleaner,
			m.DBContext,
			heavyclient.Options{
				SyncMessageLimit: m.options.heavySyncMessageLimit,
				PulsesDeltaLimit: m.options.lightChainLimit,
			},
		)
		m.syncClientsPool = heavySyncPool

		err := m.initJetSyncState(ctx)
		if err != nil {
			return err
		}
	}

	return m.restoreGenesisRecentObjects(ctx)
}

func (m *PulseManager) restoreGenesisRecentObjects(ctx context.Context) error {
	if m.NodeNet.GetOrigin().Role() == core.StaticRoleHeavyMaterial {
		return nil
	}

	jetID := *jet.NewID(0, nil)
	recent := m.RecentStorageProvider.GetStorage(ctx, jetID)

	return m.ObjectStorage.IterateIndexIDs(ctx, jetID, func(id core.RecordID) error {
		if id.Pulse() == core.FirstPulseNumber {
			recent.AddObject(ctx, id)
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
