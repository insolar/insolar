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
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/pkg/errors"
)

// JetSplitter provides method for processing and splitting jets.
type JetSplitter interface {
	// Do performs jets processing, it decides which jets to split and returns list of resulting jets).
	Do(ctx context.Context, ended, new insolar.PulseNumber) ([]insolar.JetID, error)
}

// JetSplitLimitDefault default value of how many jets LM-node should mark as intended for split.
const JetSplitLimitDefault = 5

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
	jetCalculator   JetCalculator
	jetAccessor     jet.Accessor
	jetModifier     jet.Modifier
	dropAccessor    drop.Accessor
	dropModifier    drop.Modifier
	pulseCalculator pulse.Calculator

	splitsLimit int
}

// NewJetSplitter returns a new instance of a default jet splitter implementation.
func NewJetSplitter(
	jetCalculator JetCalculator,
	jetAccessor jet.Accessor,
	jetModifier jet.Modifier,
	dropAccessor drop.Accessor,
	dropModifier drop.Modifier,
	pulseCalculator pulse.Calculator,
) *JetSplitterDefault {
	return &JetSplitterDefault{
		jetCalculator:   jetCalculator,
		jetAccessor:     jetAccessor,
		jetModifier:     jetModifier,
		dropAccessor:    dropAccessor,
		dropModifier:    dropModifier,
		pulseCalculator: pulseCalculator,

		splitsLimit: JetSplitLimitDefault,
	}
}

// Do performs jets processing, it decides which jets to split and returns list of resulting jets.
func (js *JetSplitterDefault) Do(
	ctx context.Context,
	endedPulse, newPulse insolar.PulseNumber,
) ([]insolar.JetID, error) {
	ctx, span := instracer.StartSpan(ctx, "jets.split")
	defer span.End()
	ctx, _ = inslogger.WithField(ctx, "current_pulse", endedPulse.String())
	inslog := inslogger.FromContext(ctx).WithField("split_for_pulse", newPulse.String())

	// copy current jets for new pulse, for further jets modification in new pulse.
	err := js.jetModifier.Clone(ctx, endedPulse, newPulse)
	if err != nil {
		panic("Failed to clone jets")
	}

	all := js.jetCalculator.MineForPulse(ctx, endedPulse)
	// result at least the same size
	result := make([]insolar.JetID, 0, len(all))

	var splitCandidatesIndexes []int
	for i, jetID := range all {
		// if no split intention, add to next pulse split candidates if splitsLimit counter greater than zero
		if !js.prevDropExistsAndHasSplitFlag(ctx, endedPulse, jetID) {
			// mark jet as actual for new pulse
			if err := js.jetModifier.Update(ctx, newPulse, true, jetID); err != nil {
				panic("failed to update jets on LM-node: " + err.Error())
			}
			if js.splitsLimit > 0 {
				splitCandidatesIndexes = append(splitCandidatesIndexes, i)
			}
			result = append(result, jetID)
			continue
		}

		// split jet for new pulse if it got a split intention on previous pulse.
		leftJetID, rightJetID, err := js.jetModifier.Split(ctx, newPulse, jetID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to split jet tree")
		}

		inslog.WithFields(map[string]interface{}{
			"left_child":  leftJetID.DebugString(),
			"right_child": rightJetID.DebugString(),
		}).Info("jet split performed")
		result = append(result, leftJetID, rightJetID)
	}

	// split intent lottery (for split at new pulse)
	var intentIdx *int
	if len(splitCandidatesIndexes) > 0 {
		// some jet is lucky and got a split intent for new pulse
		intentIdx = &splitCandidatesIndexes[rand.Intn(len(splitCandidatesIndexes))]
		js.splitsLimit--
	}

	// save jet drops for current pulse
	for i, jetID := range all {
		intent := intentIdx != nil && *intentIdx == i
		block := drop.Drop{
			Pulse: endedPulse,
			JetID: jetID,
			Split: intent,
		}
		if err := js.dropModifier.Set(ctx, block); err != nil {
			panic(errors.Wrapf(err, "failed create drop for pulse=%v, jet=%v",
				endedPulse, jetID.DebugString()))
		}
	}

	return result, nil
}

func (js *JetSplitterDefault) prevDropExistsAndHasSplitFlag(
	ctx context.Context,
	pn insolar.PulseNumber,
	jetID insolar.JetID,
) bool {
	prevPulse, err := js.pulseCalculator.Backwards(ctx, pn, 1)
	if err != nil {
		if err == pulse.ErrNotFound {
			return false
		}
		panic("failed to fetch previous pulse")
	}
	block, err := js.dropAccessor.ForPulse(ctx, jetID, prevPulse.PulseNumber)
	if err != nil {
		if err == drop.ErrNotFound {
			// it could happen in two cases:
			// 1) Previous drop does not exist for first pulse after (re)start.
			// 2) Previous drop was split in the previous pulse, hence has different jet.
			//    Returning false because it cannot be split again.
			return false
		}
		panic(errors.Wrapf(err, "failed to get drop for pulse=%v and jetID=%v", pn, jetID.DebugString()))
	}
	return block.Split
}
