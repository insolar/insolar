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
	"sync"

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
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/light/replication"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
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
	JetSplitter executor.JetSplitter

	IndexBucketAccessor object.IndexBucketAccessor
	PendingAccessor     object.PendingAccessor

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

	WriteManager hot.WriteManager

	currentPulse insolar.Pulse

	// setLock locks Set method call.
	setLock sync.RWMutex
	// saves PM stopping mode
	stopped bool

	// stores pulse manager options
	options pmOptions
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
	jetSplitter executor.JetSplitter,
	idxReplicaAccessor object.IndexBucketAccessor,
	lightToHeavySyncer replication.LightReplicator,
	writeManager hot.WriteManager,
	pendingAccessor object.PendingAccessor,
) *PulseManager {
	pmconf := conf.PulseManager

	pm := &PulseManager{
		currentPulse: *insolar.GenesisPulse,
		options: pmOptions{
			splitThreshold:   pmconf.SplitThreshold,
			storeLightPulses: conf.LightChainLimit,
			lightChainLimit:  conf.LightChainLimit,
		},
		DropCleaner:         dropCleaner,
		BlobCleaner:         blobCleaner,
		BlobSyncAccessor:    blobSyncAccessor,
		PulseShifter:        pulseShifter,
		RecCleaner:          recCleaner,
		RecSyncAccessor:     recSyncAccessor,
		JetSplitter:         jetSplitter,
		IndexBucketAccessor: idxReplicaAccessor,
		LightReplicator:     lightToHeavySyncer,
		WriteManager:        writeManager,
		PendingAccessor:     pendingAccessor,
	}
	return pm
}

func (m *PulseManager) processEndPulse(
	ctx context.Context,
	jets []executor.JetInfo,
	currentPulse, newPulse insolar.Pulse,
) error {
	var g errgroup.Group
	ctx, span := instracer.StartSpan(ctx, "pulse.process_end")
	defer span.End()

	logger := inslogger.FromContext(ctx)
	for _, i := range jets {
		info := i

		g.Go(func() error {
			drop, err := m.createDrop(ctx, info, currentPulse.PulseNumber)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", currentPulse.PulseNumber)
			}

			sender := func(msg message.HotData, jetID insolar.JetID) {
				ctx, span := instracer.StartSpan(ctx, "pulse.send_hot")
				defer span.End()
				msg.Jet = *insolar.NewReference(insolar.ID(jetID))
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

			if !info.SplitPerformed {
				msg, err := m.getExecutorHotData(
					ctx, info.ID, currentPulse.PulseNumber, newPulse.PulseNumber, drop,
				)
				if err != nil {
					return errors.Wrapf(err, "getExecutorData failed for jet ID %v", info.ID)
				}
				// No Split happened.
				go sender(*msg, info.ID)
			} else {
				msg, err := m.getExecutorHotData(
					ctx, info.ID, currentPulse.PulseNumber, newPulse.PulseNumber, drop,
				)
				if err != nil {
					return errors.Wrapf(err, "getExecutorData failed for jet ID %v", info.ID)
				}
				// SplitIntent happened.
				left, right := jet.Siblings(info.ID)
				go sender(*msg, left)
				go sender(*msg, right)
			}

			m.RecentStorageProvider.RemovePendingStorage(ctx, insolar.ID(info.ID))

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
	info executor.JetInfo,
	currentPulse insolar.PulseNumber,
) (
	block *drop.Drop,
	err error,
) {
	block = &drop.Drop{
		Pulse: currentPulse,
		JetID: info.ID,
		Split: info.SplitIntent,
	}

	err = m.DropModifier.Set(ctx, *block)
	if err != nil {
		return nil, errors.Wrap(err, "[ createDrop ] Can't SetDrop")
	}

	return block, nil
}

func (m *PulseManager) getExecutorHotData(
	ctx context.Context,
	jetID insolar.JetID,
	currentPN insolar.PulseNumber,
	newPulsePN insolar.PulseNumber,
	drop *drop.Drop,
) (*message.HotData, error) {
	ctx, span := instracer.StartSpan(ctx, "pulse.prepare_hot_data")
	defer span.End()

	pendingRequests := map[insolar.ID]recentstorage.PendingObjectContext{}

	bucks := m.IndexBucketAccessor.ForPNAndJet(ctx, currentPN, jetID)
	limitPN, err := m.PulseCalculator.Backwards(ctx, currentPN, m.options.lightChainLimit)
	if err == pulse.ErrNotFound {
		limitPN = *insolar.GenesisPulse
	} else if err != nil {
		inslogger.FromContext(ctx).Errorf("failed to fetch limit %v", err)
		return nil, err
	}

	hotIndexes := []message.HotIndex{}
	for _, meta := range bucks {
		encoded, err := meta.Lifeline.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).WithField("id", meta.ObjID.DebugString()).Error("failed to marshal lifeline")
			continue
		}
		if meta.LifelineLastUsed < limitPN.PulseNumber {
			continue
		}

		hotIndexes = append(hotIndexes, message.HotIndex{
			LifelineLastUsed: meta.LifelineLastUsed,
			ObjID:            meta.ObjID,
			Index:            encoded,
		})
	}

	pendingStorage := m.RecentStorageProvider.GetPendingStorage(ctx, insolar.ID(jetID))
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
		PulseNumber:     newPulsePN,
		HotIndexes:      hotIndexes,
		PendingRequests: pendingRequests,
	}
	return msg, nil
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

	logger := inslogger.FromContext(ctx)

	if oldPulse != nil && prevPN != nil {
		err = m.WriteManager.CloseAndWait(ctx, oldPulse.PulseNumber)
		if err != nil {
			logger.Error("can't close pulse for writing", err)
		}
		err = m.processEndPulse(ctx, jets, *oldPulse, newPulse)
		if err != nil {
			return err
		}
		m.postProcessJets(ctx, jets)
		go m.LightReplicator.NotifyAboutPulse(ctx, newPulse.PulseNumber)
	}

	err = m.WriteManager.Open(ctx, newPulse.PulseNumber)
	if err != nil {
		logger.Error("can't open pulse for writing", err)
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
	[]executor.JetInfo, *insolar.Pulse, *insolar.PulseNumber, error,
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

	var jets []executor.JetInfo
	if persist && prevPN != nil && oldPulse != nil {
		jets, err = m.JetSplitter.Do(ctx, *prevPN, oldPulse.PulseNumber, newPulse.PulseNumber)

		// We just joined to network
		if errors.Cause(err) == node.ErrNoNodes {
			return jets, oldPulse, prevPN, nil
		}
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if oldPulse != nil && prevPN != nil {
		m.prepareArtifactManagerMessageHandlerForNextPulse(ctx, newPulse)
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

func (m *PulseManager) postProcessJets(ctx context.Context, jets []executor.JetInfo) {
	ctx, span := instracer.StartSpan(ctx, "jets.post_process")
	defer span.End()

	for _, jetInfo := range jets {
		if !jetInfo.MineNext {
			m.RecentStorageProvider.RemovePendingStorage(ctx, insolar.ID(jetInfo.ID))
		}
	}
}

func (m *PulseManager) prepareArtifactManagerMessageHandlerForNextPulse(ctx context.Context, newPulse insolar.Pulse) {
	ctx, span := instracer.StartSpan(ctx, "early.close")
	defer span.End()

	m.JetReleaser.ThrowTimeout(ctx, newPulse.PulseNumber)
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
