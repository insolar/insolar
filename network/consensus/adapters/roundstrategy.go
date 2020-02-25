// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package adapters

import (
	"context"
	"math/rand"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/pulse"
)

type RoundStrategy struct {
	localConfig api.LocalNodeConfiguration
}

func NewRoundStrategy(
	localConfig api.LocalNodeConfiguration,
) *RoundStrategy {
	return &RoundStrategy{
		localConfig: localConfig,
	}
}

func (rs *RoundStrategy) ConfigureRoundContext(ctx context.Context, expectedPulse pulse.Number, self profiles.LocalNode) context.Context {
	ctx, _ = inslogger.WithFields(ctx, map[string]interface{}{
		"is_joiner":   self.IsJoiner(),
		"round_pulse": expectedPulse,
	})
	return ctx
}

func (rs *RoundStrategy) GetBaselineWeightForNeighbours() uint32 {
	return rand.Uint32()
}

func (rs *RoundStrategy) ShuffleNodeSequence(n int, swap func(i, j int)) {
	rand.Shuffle(n, swap)
}
