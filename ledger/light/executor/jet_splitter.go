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

// JetSplitter provides method for processing and splitting jets.
type JetSplitter interface {
	// Do performs jets processing, it decides which jets to split and returns list of resulting jets).
	Do(ctx context.Context, previous, current, new insolar.PulseNumber) ([]JetInfo, error)
}

// JetSplitterDefaultCount default value for initial jets splitting.
const JetSplitterDefaultCount = 5

// JetInfo holds info about jet.
type JetInfo struct {
	ID insolar.JetID
	// SplitIntent indicates what jet has intention to do split.
	SplitIntent bool
	// SplitPerformed indicates what jet was split.
	SplitPerformed bool

	// deprecated
	// MineNext  if not set pendings would be removed for this jet from recent storage.
	MineNext bool
}

// JetSplitterDefault implements JetSplitter.
type JetSplitterDefault struct {
	splitCount int

	jetCoordinator jet.Coordinator
	jetAccessor    jet.Accessor
	jetModifier    jet.Modifier

	dropAccessor          drop.Accessor
	recentStorageProvider recentstorage.Provider
}

// NewJetSplitter returns a new instance of a default jet splitter implementation.
func NewJetSplitter(
	jetCoordinator jet.Coordinator,
	jetAccessor jet.Accessor,
	jetModifier jet.Modifier,
	dropAccessor drop.Accessor,
	recentStorageProvider recentstorage.Provider,
) *JetSplitterDefault {
	return &JetSplitterDefault{
		splitCount: JetSplitterDefaultCount,

		jetCoordinator: jetCoordinator,
		jetAccessor:    jetAccessor,
		jetModifier:    jetModifier,

		dropAccessor:          dropAccessor,
		recentStorageProvider: recentStorageProvider,
	}
}

// Do performs jets processing, it decides which jets to split and returns list of resulting jets.
func (js *JetSplitterDefault) Do(
	ctx context.Context, previous, current, newpulse insolar.PulseNumber,
) ([]JetInfo, error) {
	jets, err := js.processJets(ctx, previous, current, newpulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process jets")
	}

	jets, err = js.splitJets(ctx, jets, previous, current, newpulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to split jets")
	}
	return jets, nil
}

func (js *JetSplitterDefault) splitJets(
	ctx context.Context,
	jets []JetInfo,
	previous, current, newpulse insolar.PulseNumber,
) ([]JetInfo, error) {
	me := js.jetCoordinator.Me()
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"current_pulse": current,
		"new_pulse":     newpulse,
	})

	for i, jetInfo := range jets {
		newInfo := JetInfo{ID: jetInfo.ID}
		if js.hasSplitIntention(ctx, previous, jetInfo.ID) {
			leftJetID, rightJetID, err := js.jetModifier.Split(
				ctx,
				newpulse,
				jetInfo.ID,
			)
			if err != nil {
				return nil, errors.Wrap(err, "failed to split jet tree")
			}

			// Set actual because we are the last executor for jet.
			js.jetModifier.Update(ctx, newpulse, true, leftJetID, rightJetID)
			newInfo.SplitPerformed = true

			nextLeftExecutor, err := js.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(leftJetID), newpulse)
			if err != nil {
				return nil, err
			}

			if *nextLeftExecutor == me {
				js.recentStorageProvider.ClonePendingStorage(ctx, insolar.ID(jetInfo.ID), insolar.ID(leftJetID))
			}

			nextRightExecutor, err := js.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(rightJetID), newpulse)
			if err != nil {
				return nil, err
			}

			if *nextRightExecutor == me {
				js.recentStorageProvider.ClonePendingStorage(ctx, insolar.ID(jetInfo.ID), insolar.ID(rightJetID))
			}

			logger.WithFields(map[string]interface{}{
				"left_child":  leftJetID.DebugString(),
				"right_child": rightJetID.DebugString(),
			}).Info("jet Split performed")

			jets[i] = newInfo
		} else {
			// Set actual because we are the last executor for jet.
			js.jetModifier.Update(ctx, newpulse, true, jetInfo.ID)
			nextExecutor, err := js.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(jetInfo.ID), newpulse)
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

func (js *JetSplitterDefault) processJets(ctx context.Context, previous, current, new insolar.PulseNumber) ([]JetInfo, error) {
	ctx, span := instracer.StartSpan(ctx, "jets.process")
	defer span.End()

	js.jetModifier.Clone(ctx, current, new)

	ids := js.jetAccessor.All(ctx, current)
	ids, err := js.filterOtherExecutors(ctx, current, ids)
	if err != nil {
		return nil, err
	}

	var results []JetInfo                     // nolint: prealloc
	var withoutSplitIntention []insolar.JetID // nolint: prealloc
	for _, id := range ids {
		if js.hasSplitIntention(ctx, previous, id) {
			results = append(results, JetInfo{ID: id})
		} else {
			withoutSplitIntention = append(withoutSplitIntention, id)
		}
	}

	if len(withoutSplitIntention) == 0 {
		return results, nil
	}

	indexToSplit := rand.Intn(len(withoutSplitIntention))
	for i, jetID := range withoutSplitIntention {
		info := JetInfo{ID: jetID}
		if indexToSplit == i && js.splitCount > 0 {
			js.splitCount--
			info.SplitIntent = true
		}
		results = append(results, info)
	}
	return results, nil
}

func (js *JetSplitterDefault) filterOtherExecutors(ctx context.Context, pulse insolar.PulseNumber, ids []insolar.JetID) ([]insolar.JetID, error) {
	me := js.jetCoordinator.Me()
	result := []insolar.JetID{}
	for _, id := range ids {
		executor, err := js.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(id), pulse)
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
	drop, err := js.dropAccessor.ForPulse(ctx, id, previous)
	if err != nil {
		inslogger.FromContext(ctx).WithFields(map[string]interface{}{
			"previous_pulse": previous,
			"jet_id":         id,
		}).Warn(errors.Wrapf(err, "failed to get drop by jet.ID=%v previous_pulse=%v", id.DebugString(), previous))
		return false
	}
	return drop.Split
}
