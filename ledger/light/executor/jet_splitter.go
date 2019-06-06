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

package executor

import (
	"context"
	"math/rand"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/pkg/errors"
)

// JetSplitter provides methods for calculating jets.
type JetSplitter interface {
	// Do performs jets processing (decides which jets to split and returns result)
	Do(ctx context.Context, previous, current, new insolar.PulseNumber) ([]jet.Info, error)
}

// TODO: move to JetSplitterDefault
var splitCount = 5

type JetSplitterDefault struct {
	JetCoordinator jet.Coordinator
	JetAccessor    jet.Accessor
	JetModifier    jet.Modifier

	DropAccessor          drop.Accessor
	RecentStorageProvider recentstorage.Provider
}

func (js *JetSplitterDefault) Do(
	ctx context.Context, previous, current, new insolar.PulseNumber,
) ([]jet.Info, error) {
	jets, err := js.processJets(ctx, previous, current, new)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process jets")
	}

	jets, err = js.splitJets(ctx, jets, previous, current, new)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Split jets")
	}

	return jets, nil
}

func (m *JetSplitterDefault) splitJets(ctx context.Context, jets []jet.Info, previous, current, new insolar.PulseNumber) ([]jet.Info, error) {
	me := m.JetCoordinator.Me()
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"current_pulse": current,
		"new_pulse":     new,
	})

	for i, jetInfo := range jets {
		newInfo := jet.Info{ID: jetInfo.ID}
		if m.hasSplitIntention(ctx, previous, jetInfo.ID) {
			leftJetID, rightJetID, err := m.JetModifier.Split(
				ctx,
				new,
				jetInfo.ID,
			)

			if err != nil {
				return nil, errors.Wrap(err, "failed to Split jet tree")
			}

			// Set actual because we are the last executor for jet.
			m.JetModifier.Update(ctx, new, true, leftJetID, rightJetID)
			newInfo.Left = &jet.Info{ID: leftJetID}
			newInfo.Right = &jet.Info{ID: rightJetID}

			nextLeftExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(leftJetID), new)
			if err != nil {
				return nil, err
			}
			if *nextLeftExecutor == me {
				newInfo.Left.MineNext = true
				m.RecentStorageProvider.ClonePendingStorage(ctx, insolar.ID(jetInfo.ID), insolar.ID(leftJetID))
			}
			nextRightExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(rightJetID), new)
			if err != nil {
				return nil, err
			}
			if *nextRightExecutor == me {
				newInfo.Right.MineNext = true
				m.RecentStorageProvider.ClonePendingStorage(ctx, insolar.ID(jetInfo.ID), insolar.ID(rightJetID))
			}

			logger.WithFields(map[string]interface{}{
				"left_child":  leftJetID.DebugString(),
				"right_child": rightJetID.DebugString(),
			}).Info("jet Split performed")

			jets[i] = newInfo
		} else {
			// Set actual because we are the last executor for jet.
			m.JetModifier.Update(ctx, new, true, jetInfo.ID)
			nextExecutor, err := m.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(jetInfo.ID), new)
			if err != nil {
				return nil, err
			}
			if *nextExecutor == me {
				newInfo.MineNext = true
			}
		}
	}
	return jets, nil
}

func (js *JetSplitterDefault) processJets(ctx context.Context, previous, current, new insolar.PulseNumber) ([]jet.Info, error) {
	ctx, span := instracer.StartSpan(ctx, "jets.process")
	defer span.End()

	js.JetModifier.Clone(ctx, current, new)

	ids := js.JetAccessor.All(ctx, current)
	ids, err := js.filterOtherExecutors(ctx, current, ids)
	if err != nil {
		return nil, err
	}

	var results []jet.Info                    // nolint: prealloc
	var withoutSplitIntention []insolar.JetID // nolint: prealloc
	for _, id := range ids {
		if js.hasSplitIntention(ctx, previous, id) {
			results = append(results, jet.Info{ID: id})
		} else {
			withoutSplitIntention = append(withoutSplitIntention, id)
		}
	}

	if len(withoutSplitIntention) == 0 {
		return results, nil
	}

	indexToSplit := rand.Intn(len(withoutSplitIntention))
	for i, jetID := range withoutSplitIntention {
		info := jet.Info{ID: jetID}
		if indexToSplit == i && splitCount > 0 {
			splitCount--
			info.Split = true
		}
		results = append(results, info)
	}
	return results, nil
}

func (js *JetSplitterDefault) filterOtherExecutors(ctx context.Context, pulse insolar.PulseNumber, ids []insolar.JetID) ([]insolar.JetID, error) {
	me := js.JetCoordinator.Me()
	result := []insolar.JetID{}
	for _, id := range ids {
		executor, err := js.JetCoordinator.LightExecutorForJet(ctx, insolar.ID(id), pulse)
		if err != nil && err != node.ErrNoNodes {
			return nil, err
		}
		if executor == nil || err != nil {
			continue
		}

		if *executor == me {
			result = append(result, id)
		}
	}
	return result, nil
}

func (js *JetSplitterDefault) hasSplitIntention(
	ctx context.Context,
	previous insolar.PulseNumber,
	id insolar.JetID,
) bool {
	drop, err := js.DropAccessor.ForPulse(ctx, id, previous)
	if err != nil {
		inslogger.FromContext(ctx).WithFields(map[string]interface{}{
			"previous_pulse": previous,
			"jet_id":         id,
		}).Warn(errors.Wrapf(err, "failed to get drop by jet.ID=%v previous_pulse=%v", id.DebugString(), previous))
		return false
	}
	return drop.Split
}
