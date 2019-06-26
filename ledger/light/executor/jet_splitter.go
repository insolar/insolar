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
	// SplitIntent indicates what jet has intention to do split in next pulse.
	SplitIntent bool
	// MustSplit indicates what jet should be split in current pulse.
	MustSplit bool
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
	ctx, span := instracer.StartSpan(ctx, "jets.split")
	defer span.End()
	ctx, _ = inslogger.WithField(ctx, "current_pulse", current.String())

	err := js.jetModifier.Clone(ctx, current, newpulse)
	if err != nil {
		panic("Failed to clone jets")
	}
	jets := js.prepareJetInfo(ctx, previous, current)
	return js.splitJets(ctx, jets, newpulse)
}

func (js *JetSplitterDefault) splitJets(
	ctx context.Context,
	jets []JetInfo,
	pn insolar.PulseNumber,
) ([]JetInfo, error) {
	inslog := inslogger.FromContext(ctx).WithField("split_for_pulse", pn.String())

	result := make([]JetInfo, 0, len(jets))
	for _, jetInfo := range jets {
		if !jetInfo.MustSplit {
			err := js.jetModifier.Update(ctx, pn, true, jetInfo.ID)
			if err != nil {
				panic("failed to update jets on LM-node: " + err.Error())
			}
			result = append(result, jetInfo)
			continue
		}

		leftJetID, rightJetID, err := js.jetModifier.Split(ctx, pn, jetInfo.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to split jet tree")
		}

		// Set actual because we are the last executor for jet.
		err = js.jetModifier.Update(ctx, pn, true, leftJetID, rightJetID)
		if err != nil {
			panic("failed to update jets on LM-node: " + err.Error())
		}

		inslog.WithFields(map[string]interface{}{
			"left_child":  leftJetID.DebugString(),
			"right_child": rightJetID.DebugString(),
		}).Info("jet Split performed")

		result = append(result,
			JetInfo{ID: leftJetID},
			JetInfo{ID: rightJetID},
		)
	}
	return result, nil
}

func (js *JetSplitterDefault) prepareJetInfo(ctx context.Context, previous, current insolar.PulseNumber) []JetInfo {
	var results []JetInfo
	var nextSplitCandidates []insolar.JetID
	for _, id := range js.myJetsForPulse(ctx, current) {
		shouldSplit := js.dropForJetHasSplitFlag(ctx, previous, id)
		if !shouldSplit && js.splitCount > 0 {
			nextSplitCandidates = append(nextSplitCandidates, id)
			continue
		}
		results = append(results, JetInfo{
			ID:        id,
			MustSplit: shouldSplit,
		})
	}

	// lottery of whom has got split intention
	if len(nextSplitCandidates) > 0 {
		js.splitCount--
		splitIdx := rand.Intn(len(nextSplitCandidates))
		for i, jetID := range nextSplitCandidates {
			results = append(results, JetInfo{
				ID:          jetID,
				SplitIntent: splitIdx == i,
			})
		}
	}
	return results
}

func (js *JetSplitterDefault) myJetsForPulse(ctx context.Context, pn insolar.PulseNumber) []insolar.JetID {
	me := js.jetCoordinator.Me()
	var myJets []insolar.JetID
	for _, id := range js.jetAccessor.All(ctx, pn) {
		executor, err := js.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(id), pn)
		if err != nil && err != node.ErrNoNodes {
			panic(err)
		}
		if executor == nil || err != nil {
			continue
		}
		if *executor == me {
			myJets = append(myJets, id)
		}
	}
	return myJets
}

func (js *JetSplitterDefault) dropForJetHasSplitFlag(
	ctx context.Context,
	pn insolar.PulseNumber,
	id insolar.JetID,
) bool {
	block, err := js.dropAccessor.ForPulse(ctx, id, pn)
	if err != nil {
		inslogger.FromContext(ctx).WithFields(map[string]interface{}{
			"pulse":  pn,
			"jet_id": id.DebugString(),
		}).Warn(errors.Wrapf(err, "failed to get drop by jet ID"))
		return false
	}
	return block.Split
}
