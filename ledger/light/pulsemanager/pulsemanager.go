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
	"fmt"
	"math/rand"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/artifactmanager"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/light/replication"
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/pulsemanager.ActiveListSwapper -o ../../../testutils -s _mock.go
type ActiveListSwapper interface {
	MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error
}

// PulseManager implements insolar.PulseManager.
type PulseManager struct {
	Bus                        insolar.MessageBus                 `inject:""`
	NodeNet                    insolar.NodeNetwork                `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	GIL                        insolar.GlobalInsolarLock          `inject:""`
	CryptographyService        insolar.CryptographyService        `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	RecentStorageProvider      recentstorage.Provider             `inject:""`
	ActiveListSwapper          ActiveListSwapper                  `inject:""`
	MessageHandler             *artifactmanager.MessageHandler

	JetReleaser hot.JetReleaser `inject:""`

	JetAccessor jet.Accessor `inject:""`
	JetModifier jet.Modifier `inject:""`

	IndexAccessor           object.IndexAccessor `inject:""`
	IndexModifier           object.ExtendedIndexModifier
	CollectionIndexAccessor object.IndexCollectionAccessor
	IndexCleaner            object.IndexCleaner

	NodeSetter node.Modifier `inject:""`
	Nodes      node.Accessor `inject:""`

	DropModifier drop.Modifier `inject:""`
	DropAccessor drop.Accessor `inject:""`
	DropCleaner  drop.Cleaner

	PulseAccessor   pulse.Accessor   `inject:""`
	PulseCalculator pulse.Calculator `inject:""`
	PulseAppender   pulse.Appender   `inject:""`
	PulseShifter    pulse.Shifter

	BlobSyncAccessor blob.CollectionAccessor
	BlobCleaner      blob.Cleaner

	RecSyncAccessor object.RecordCollectionAccessor
	RecCleaner      object.RecordCleaner

	LightReplicator replication.LightReplicator

	currentPulse insolar.Pulse

	// setLock locks Set method call.
	setLock sync.RWMutex
	// saves PM stopping mode
	stopped bool

	// stores pulse manager options
	options pmOptions
}

type jetInfo struct {
	id       insolar.JetID
	mineNext bool
	left     *jetInfo
	right    *jetInfo
}

// Just store ledger configuration in PM. This is not required.
type pmOptions struct {
	// enableSync            bool
	splitThreshold   uint64
	storeLightPulses int
	// heavySyncMessageLimit int
	lightChainLimit int
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(
	conf configuration.Ledger,
	dropCleaner drop.Cleaner,
	blobCleaner blob.Cleaner,
	blobSyncAccessor blob.CollectionAccessor,
	pulseShifter pulse.Shifter,
	recCleaner object.RecordCleaner,
	recSyncAccessor object.RecordCollectionAccessor,
	idxCollectionAccessor object.IndexCollectionAccessor,
	indexCleaner object.IndexCleaner,
	lightToHeavySyncer replication.LightReplicator,
) *PulseManager {
	pmconf := conf.PulseManager

	pm := &PulseManager{
		currentPulse: *insolar.GenesisPulse,
		options: pmOptions{
			splitThreshold:   pmconf.SplitThreshold,
			storeLightPulses: conf.LightChainLimit,
			lightChainLimit:  conf.LightChainLimit,
		},
		DropCleaner:             dropCleaner,
		BlobCleaner:             blobCleaner,
		BlobSyncAccessor:        blobSyncAccessor,
		PulseShifter:            pulseShifter,
		RecCleaner:              recCleaner,
		RecSyncAccessor:         recSyncAccessor,
		CollectionIndexAccessor: idxCollectionAccessor,
		IndexCleaner:            indexCleaner,
		LightReplicator:         lightToHeavySyncer,
	}
	return pm
}

func (m *PulseManager) processEndPulse(
	ctx context.Context,
	jets []jetInfo,
	prevPulseNumber insolar.PulseNumber,
	currentPulse, newPulse insolar.Pulse,
) error {
	var g errgroup.Group
	ctx, span := instracer.StartSpan(ctx, "pulse.process_end")
	defer span.End()

	logger := inslogger.FromContext(ctx)
	for _, i := range jets {
		info := i

		g.Go(func() error {
			drop, dropSerialized, _, err := m.createDrop(ctx, insolar.ID(info.id), prevPulseNumber, currentPulse.PulseNumber)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", currentPulse.PulseNumber)
			}

			sender := func(msg message.HotData, jetID insolar.JetID) {
				ctx, span := instracer.StartSpan(ctx, "pulse.send_hot")
				defer span.End()
				msg.Jet = *insolar.NewReference(insolar.DomainID, insolar.ID(jetID))
				genericRep, err := m.Bus.Send(ctx, &msg, nil)
				if err != nil {
					logger.WithField("err", err).Error("failed to send hot data")
					return
				}
				if _, ok := genericRep.(*reply.OK); !ok {
					logger.WithField(
						"err",
						fmt.Sprintf("unexpected reply: %T", genericRep),
					).Error("failed to send hot data")
					return
				}
			}

			if info.left == nil && info.right == nil {
				msg, err := m.getExecutorHotData(
					ctx, insolar.ID(info.id), newPulse.PulseNumber, drop, dropSerialized,
				)
				if err != nil {
					return errors.Wrapf(err, "getExecutorData failed for jet id %v", info.id)
				}
				// No split happened.
				if !info.mineNext {
					go sender(*msg, info.id)
				}
			} else {
				msg, err := m.getExecutorHotData(
					ctx, insolar.ID(info.id), newPulse.PulseNumber, drop, dropSerialized,
				)
				if err != nil {
					return errors.Wrapf(err, "getExecutorData failed for jet id %v", info.id)
				}
				// Split happened.
				if !info.left.mineNext {
					go sender(*msg, info.left.id)
				}
				if !info.right.mineNext {
					go sender(*msg, info.right.id)
				}
			}

			m.RecentStorageProvider.RemovePendingStorage(ctx, insolar.ID(info.id))

			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return errors.Wrap(err, "got error on jets sync")
	}

	return nil
}

func (m *PulseManager) createDrop(
	ctx context.Context,
	jetID insolar.ID,
	prevPulse, currentPulse insolar.PulseNumber,
) (
	block *drop.Drop,
	dropSerialized []byte,
	messages [][]byte,
	err error,
) {
	block = &drop.Drop{
		Pulse: currentPulse,
		JetID: insolar.JetID(jetID),
	}

	err = m.DropModifier.Set(ctx, *block)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "[ createDrop ] Can't SetDrop")
	}

	dropSerialized = drop.MustEncode(block)
	return
}

func (m *PulseManager) getExecutorHotData(
	ctx context.Context,
	jetID insolar.ID,
	pulse insolar.PulseNumber,
	drop *drop.Drop,
	dropSerialized []byte,
) (*message.HotData, error) {
	ctx, span := instracer.StartSpan(ctx, "pulse.prepare_hot_data")
	defer span.End()

	hotIndexes := map[insolar.ID]message.HotIndex{}
	pendingRequests := map[insolar.ID]recentstorage.PendingObjectContext{}

	idxs := m.CollectionIndexAccessor.ForJet(ctx, insolar.JetID(jetID))

	for id, meta := range idxs {
		encoded := object.EncodeIndex(meta.Index)
		hotIndexes[id] = message.HotIndex{
			LastUsed: meta.LastUsed,
			Index:    encoded,
		}
	}

	pendingStorage := m.RecentStorageProvider.GetPendingStorage(ctx, jetID)
	requestCount := 0
	for objID, objContext := range pendingStorage.GetRequests() {
		if len(objContext.Requests) > 0 {
			pendingRequests[objID] = objContext
			requestCount += len(objContext.Requests)
		}
	}

	stats.Record(
		ctx,
		statHotObjectsSent.M(int64(len(hotIndexes))),
		statPendingSent.M(int64(requestCount)),
	)

	msg := &message.HotData{
		Drop:            *drop,
		PulseNumber:     pulse,
		HotIndexes:      hotIndexes,
		PendingRequests: pendingRequests,
	}
	return msg, nil
}

var splitCount = 5

func (m *PulseManager) processJets(ctx context.Context, currentPulse, newPulse insolar.PulseNumber) ([]jetInfo, error) {
	ctx, span := instracer.StartSpan(ctx, "jets.process")
	defer span.End()

	m.JetModifier.Clone(ctx, currentPulse, newPulse)

	if m.NodeNet.GetOrigin().Role() != insolar.StaticRoleLightMaterial {
		return nil, nil
	}

	var results []jetInfo
	jetIDs := m.JetAccessor.All(ctx, newPulse)
	me := m.JetCoordinator.Me()
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"current_pulse": currentPulse,
		"new_pulse":     newPulse,
	})
	indexToSplit := rand.Intn(len(jetIDs))
	for i, jetID := range jetIDs {
		wasExecutor := false
		executor, err := m.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(jetID), currentPulse)
		if err != nil && err != node.ErrNoNodes {
			return nil, err
		}
		if err == nil {
			wasExecutor = *executor == me
		}

		logger = logger.WithField("jetid", jetID.DebugString())
		inslogger.SetLogger(ctx, logger)
		logger.WithField("i_was_executor", wasExecutor).Debug("process jet")
		if !wasExecutor {
			continue
		}

		info := jetInfo{id: jetID}
		if indexToSplit == i && splitCount > 0 {
			splitCount--

			leftJetID, rightJetID, err := m.JetModifier.Split(
				ctx,
				newPulse,
				jetID,
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to split jet tree")
			}

			// Set actual because we are the last executor for jet.
			m.JetModifier.Update(ctx, newPulse, true, leftJetID, rightJetID)

			info.left = &jetInfo{id: leftJetID}
			info.right = &jetInfo{id: rightJetID}
			nextLeftExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(leftJetID), newPulse)
			if err != nil {
				return nil, err
			}
			if *nextLeftExecutor == me {
				info.left.mineNext = true
				m.rewriteHotData(ctx, jetID, leftJetID)
			}
			nextRightExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(rightJetID), newPulse)
			if err != nil {
				return nil, err
			}
			if *nextRightExecutor == me {
				info.right.mineNext = true
				m.rewriteHotData(ctx, jetID, rightJetID)
			}

			logger.WithFields(map[string]interface{}{
				"left_child":  leftJetID.DebugString(),
				"right_child": rightJetID.DebugString(),
			}).Info("jet split performed")
		} else {
			// Set actual because we are the last executor for jet.
			m.JetModifier.Update(ctx, newPulse, true, jetID)
			nextExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(jetID), newPulse)
			if err != nil {
				return nil, err
			}
			if *nextExecutor == me {
				info.mineNext = true
			}
		}
		results = append(results, info)
	}

	return results, nil
}

func (m *PulseManager) rewriteHotData(ctx context.Context, fromJetID, toJetID insolar.JetID) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"from_jet": fromJetID.DebugString(),
		"to_jet":   toJetID.DebugString(),
	})

	idxs := m.CollectionIndexAccessor.ForJet(ctx, fromJetID)
	for id, meta := range idxs {
		meta.Index.JetID = toJetID
		err := m.IndexModifier.SetWithMeta(ctx, id, meta.LastUsed, meta.Index)
		if err == object.ErrIndexNotFound {
			logger.WithField("id", id.DebugString()).Error("failed to rewrite index")
			continue
		}
	}

	m.RecentStorageProvider.ClonePendingStorage(ctx, insolar.ID(fromJetID), insolar.ID(toJetID))
}

// Set set's new pulse and closes current jet drop.
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

	jets, oldPulse, prevPN, err := m.setUnderGilSection(ctx, newPulse, persist)
	if err != nil {
		return err
	}

	if !persist {
		return nil
	}

	// Run only on material executor.
	// execute only on material executor
	if m.NodeNet.GetOrigin().Role() == insolar.StaticRoleLightMaterial && oldPulse != nil && prevPN != nil {
		err = m.processEndPulse(ctx, jets, *prevPN, *oldPulse, newPulse)
		if err != nil {
			return err
		}
		m.postProcessJets(ctx, newPulse, jets)
		// m.addSync(ctx, jets, oldPulse.PulseNumber)
		// go m.cleanLightData(ctx, newPulse, jetIndexesRemoved)
		go m.LightReplicator.NotifyAboutPulse(ctx, newPulse.PulseNumber)
	}

	err = m.Bus.OnPulse(ctx, newPulse)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "MessageBus OnPulse() returns error"))
	}

	if m.MessageHandler != nil {
		m.MessageHandler.OnPulse(ctx, newPulse)
	}

	return nil
}

func (m *PulseManager) setUnderGilSection(
	ctx context.Context, newPulse insolar.Pulse, persist bool,
) (
	[]jetInfo, *insolar.Pulse, *insolar.PulseNumber, error,
) {
	var (
		oldPulse *insolar.Pulse
		prevPN   *insolar.PulseNumber
	)

	m.GIL.Acquire(ctx)
	ctx, span := instracer.StartSpan(ctx, "pulse.gil_locked")
	defer span.End()
	defer m.GIL.Release(ctx)

	// FIXME: @andreyromancev. 17.12.18. return insolar.Pulse here.
	storagePulse, err := m.PulseAccessor.Latest(ctx)
	if err != nil && err != pulse.ErrNotFound {
		return nil, nil, nil, errors.Wrap(err, "call of GetLatestPulseNumber failed")
	}

	if err != pulse.ErrNotFound {
		oldPulse = &storagePulse
		pp, err := m.PulseCalculator.Backwards(ctx, oldPulse.PulseNumber, 1)
		if err == nil {
			prevPN = &pp.PulseNumber
		} else {
			prevPN = &insolar.GenesisPulse.PulseNumber
		}
		ctx, _ = inslogger.WithField(ctx, "current_pulse", fmt.Sprintf("%d", oldPulse.PulseNumber))
	}

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
		return nil, nil, nil, errors.Wrap(err, "failed to apply new active node list")
	}
	if persist {
		if err := m.PulseAppender.Append(ctx, newPulse); err != nil {
			return nil, nil, nil, errors.Wrap(err, "call of AddPulse failed")
		}
		fromNetwork := m.NodeNet.GetWorkingNodes()
		toSet := make([]insolar.Node, 0, len(fromNetwork))
		for _, node := range fromNetwork {
			toSet = append(toSet, insolar.Node{ID: node.ID(), Role: node.Role()})
		}
		err = m.NodeSetter.Set(newPulse.PulseNumber, toSet)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "call of SetActiveNodes failed")
		}
	}

	if m.NodeNet.GetOrigin().Role() == insolar.StaticRoleHeavyMaterial {
		return nil, nil, nil, nil
	}

	var jets []jetInfo
	if persist && oldPulse != nil {
		jets, err = m.processJets(ctx, oldPulse.PulseNumber, newPulse.PulseNumber)
		// We just joined to network
		if err == node.ErrNoNodes {
			return jets, oldPulse, prevPN, nil
		}
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to process jets")
		}
	}

	if oldPulse != nil && prevPN != nil {
		if m.NodeNet.GetOrigin().Role() == insolar.StaticRoleLightMaterial {
			m.prepareArtifactManagerMessageHandlerForNextPulse(ctx, newPulse, jets)
		}
	}

	if persist && oldPulse != nil {
		nodes, err := m.Nodes.All(oldPulse.PulseNumber)
		if err != nil {
			return nil, nil, nil, err
		}
		// No active nodes for pulse. It means there was no processing (network start).
		if len(nodes) == 0 {
			// Activate zero jet for jet tree and unlock jet waiter.
			zeroJet := insolar.NewJetID(0, nil)
			m.JetModifier.Update(ctx, newPulse.PulseNumber, true, *zeroJet)
			err := m.JetReleaser.Unlock(ctx, insolar.ID(*zeroJet))
			if err != nil {
				if err == artifactmanager.ErrWaiterNotLocked {
					inslogger.FromContext(ctx).Error(err)
				} else {
					return nil, nil, nil, errors.Wrap(err, "failed to unlock zero jet")
				}
			}
		}
	}

	return jets, oldPulse, prevPN, nil
}

func (m *PulseManager) postProcessJets(ctx context.Context, newPulse insolar.Pulse, jets []jetInfo) {
	ctx, span := instracer.StartSpan(ctx, "jets.post_process")
	defer span.End()

	for _, jetInfo := range jets {
		if !jetInfo.mineNext {
			m.RecentStorageProvider.RemovePendingStorage(ctx, insolar.ID(jetInfo.id))
		}
	}
}

func (m *PulseManager) prepareArtifactManagerMessageHandlerForNextPulse(ctx context.Context, newPulse insolar.Pulse, jets []jetInfo) {
	ctx, span := instracer.StartSpan(ctx, "early.close")
	defer span.End()

	m.JetReleaser.ThrowTimeout(ctx)

	logger := inslogger.FromContext(ctx)
	for _, jetInfo := range jets {
		if jetInfo.left == nil && jetInfo.right == nil {
			// No split happened.
			if jetInfo.mineNext {
				err := m.JetReleaser.Unlock(ctx, insolar.ID(jetInfo.id))
				if err != nil {
					logger.Error(err)
				}
			}
		} else {
			// Split happened.
			if jetInfo.left.mineNext {
				err := m.JetReleaser.Unlock(ctx, insolar.ID(jetInfo.left.id))
				if err != nil {
					logger.Error(err)
				}
			}
			if jetInfo.right.mineNext {
				err := m.JetReleaser.Unlock(ctx, insolar.ID(jetInfo.right.id))
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
}

// Start starts pulse manager
func (m *PulseManager) Start(ctx context.Context) error {
	origin := m.NodeNet.GetOrigin()
	err := m.NodeSetter.Set(insolar.FirstPulseNumber, []insolar.Node{{ID: origin.ID(), Role: origin.Role()}})
	if err != nil && err != node.ErrOverride {
		return err
	}

	return nil
}

// Stop stops PulseManager
func (m *PulseManager) Stop(ctx context.Context) error {
	// There should not to be any Set call after Stop call
	m.setLock.Lock()
	m.stopped = true
	m.setLock.Unlock()

	return nil
}
