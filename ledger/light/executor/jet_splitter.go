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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

// JetSplitter provides method for processing and splitting jets.
type JetSplitter interface {
	// Do performs jets processing, it decides which jets to split and returns list of resulting jets).
	Do(ctx context.Context, ended, new insolar.PulseNumber) ([]insolar.JetID, error)
}

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
	cfg configuration.JetSplit

	jetCalculator   JetCalculator
	jetAccessor     jet.Accessor
	jetModifier     jet.Modifier
	dropAccessor    drop.Accessor
	dropModifier    drop.Modifier
	pulseCalculator pulse.Calculator
	recordsAccessor object.RecordCollectionAccessor
}

// NewJetSplitter returns a new instance of a default jet splitter implementation.
func NewJetSplitter(
	cfg configuration.JetSplit,
	jetCalculator JetCalculator,
	jetAccessor jet.Accessor,
	jetModifier jet.Modifier,
	dropAccessor drop.Accessor,
	dropModifier drop.Modifier,
	pulseCalculator pulse.Calculator,
	recordsAccessor object.RecordCollectionAccessor,
) *JetSplitterDefault {
	return &JetSplitterDefault{
		cfg: cfg,

		jetCalculator:   jetCalculator,
		jetAccessor:     jetAccessor,
		jetModifier:     jetModifier,
		dropAccessor:    dropAccessor,
		dropModifier:    dropModifier,
		pulseCalculator: pulseCalculator,
		recordsAccessor: recordsAccessor,
	}
}

// Do performs jets processing, it decides which jets to split and returns list of resulting jets.
func (js *JetSplitterDefault) Do(
	ctx context.Context,
	endedPulse, newPulse insolar.PulseNumber,
) ([]insolar.JetID, error) {
	ctx, span := instracer.StartSpan(ctx, "jets.split")
	defer span.End()
	ctx, _ = inslogger.WithField(ctx, "ended_pulse", endedPulse.String())
	inslog := inslogger.FromContext(ctx).WithField("new_pulse", newPulse.String())

	// copy current jet tree for new pulse, for further modification of jets owned in ended pulse.
	err := js.jetModifier.Clone(ctx, endedPulse, newPulse, false)
	if err != nil {
		panic("Failed to clone jets")
	}

	all := js.jetCalculator.MineForPulse(ctx, endedPulse)
	result := make([]insolar.JetID, 0, len(all)*2)
	for _, jetID := range all {
		exceed, err := js.createDrop(ctx, jetID, endedPulse)
		if err != nil {
			return nil, errors.Wrapf(err, "failed create drop for pulse=%v, jet=%v",
				endedPulse, jetID.DebugString())
		}

		if !exceed {
			// no split, just mark jet as actual for new pulse
			if err := js.jetModifier.Update(ctx, newPulse, true, jetID); err != nil {
				panic("failed to update jets on LM-node: " + err.Error())
			}
			result = append(result, jetID)

			inslogger.FromContext(ctx).Debug(">>>>>>>>>>>>>>>: DON'T SPLIT: pulse: ", endedPulse, ". JET: ", jetID.DebugString())

			continue
		}

		inslogger.FromContext(ctx).Debug(">>>>>>>>>>>>>>>: YYYY SPLIT: pulse: ", endedPulse, ". JET: ", jetID.DebugString())

		// split jet for new pulse
		leftJetID, rightJetID, err := js.jetModifier.Split(ctx, newPulse, jetID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to split jet tree")
		}
		result = append(result, leftJetID, rightJetID)

		inslog.WithFields(map[string]interface{}{
			"jet_left": leftJetID.DebugString(), "jet_right": rightJetID.DebugString(),
		}).Info("jet split performed")
	}

	return result, nil
}

func (js *JetSplitterDefault) createDrop(
	ctx context.Context,
	jetID insolar.JetID,
	pn insolar.PulseNumber,
) (bool, error) {
	block := drop.Drop{
		Pulse: pn,
		JetID: jetID,
	}

	// skip any thresholds calculation for split if jet depth for jetID reached limit.
	if jetID.Depth() >= js.cfg.DepthLimit {
		return false, js.dropModifier.Set(ctx, block)
	}

	threshold := js.getPreviousDropThreshold(ctx, jetID, pn)
	// reset threshold counter, if split is happened
	if threshold > js.cfg.ThresholdOverflowCount {
		inslogger.FromContext(ctx).Debug(">>>>>>>>>>>>>>>: RESET threshold: pulse: ", pn, ". JET: ", jetID.DebugString())
		threshold = 0
	}
	// if records count reached threshold increase counter (instead it reset)
	recordsCount := len(js.recordsAccessor.ForPulse(ctx, jetID, pn))
	if recordsCount > js.cfg.ThresholdRecordsCount {
		block.SplitThresholdExceeded = threshold + 1
		inslogger.FromContext(ctx).Debug(">>>>>>>>>>>>>>>: INCREASE threshold:", block.SplitThresholdExceeded, ". pulse: ", pn, ". JET: ", jetID.DebugString())
	}
	inslogger.FromContext(ctx).Debug(">>>>>>>>>>>>>>>: JET_SPLITTER threshold:", block.SplitThresholdExceeded, ". pulse: ", pn, ". JET: ", jetID.DebugString(), ". OVERFLOW: ", js.cfg.ThresholdOverflowCount)

	// first return value is split needed
	return block.SplitThresholdExceeded > js.cfg.ThresholdOverflowCount, js.dropModifier.Set(ctx, block)
}

func (js *JetSplitterDefault) getPreviousDropThreshold(
	ctx context.Context,
	jetID insolar.JetID,
	pn insolar.PulseNumber,
) int {
	prevPulse, err := js.pulseCalculator.Backwards(ctx, pn, 1)
	if err != nil {
		if err == pulse.ErrNotFound {
			return 0
		}
		panic("failed to fetch previous pulse")
	}
	return js.getDropThreshold(ctx, jetID, prevPulse.PulseNumber)
}

func (js *JetSplitterDefault) getDropThreshold(
	ctx context.Context,
	jetID insolar.JetID,
	pn insolar.PulseNumber,
) int {
	block, err := js.dropAccessor.ForPulse(ctx, jetID, pn)
	if err != nil {
		if err == drop.ErrNotFound {
			// it could happen in two cases:
			// 1) Previous drop does not exist for first pulse after (re)start.
			// 2) Previous drop was split in the previous pulse, hence has different jet.
			//    Returning 0 because we starting from 0 after split.
			return 0
		}
		panic(errors.Wrapf(err, "failed to get drop for pulse=%v and jetID=%v", pn, jetID.DebugString()))
	}
	return block.SplitThresholdExceeded
}
