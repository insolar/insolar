// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.JetCalculator -o ./ -s _mock.go -g

// JetCalculator provides get jets method for provided pulse.
type JetCalculator interface {
	// MineForPulse returns current node's jets for a provided pulse
	MineForPulse(ctx context.Context, pn insolar.PulseNumber) ([]insolar.JetID, error)
}

// JetCalculatorDefault implements JetCalculator.
type JetCalculatorDefault struct {
	jetCoordinator jet.Coordinator
	jetAccessor    jet.Accessor
}

// NewJetCalculator returns a new instance of a default jet calculator implementation.
func NewJetCalculator(jetCoordinator jet.Coordinator, jetAccessor jet.Accessor) *JetCalculatorDefault {
	return &JetCalculatorDefault{
		jetCoordinator: jetCoordinator,
		jetAccessor:    jetAccessor,
	}
}

// MineForPulse returns current node's jets for a provided pulse.
func (c *JetCalculatorDefault) MineForPulse(ctx context.Context, pn insolar.PulseNumber) ([]insolar.JetID, error) {
	var res []insolar.JetID

	jetIDs := c.jetAccessor.All(ctx, pn)
	me := c.jetCoordinator.Me()

	for _, jetID := range jetIDs {
		executor, err := c.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(jetID), pn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate executor")
		}
		if *executor == me {
			res = append(res, jetID)
		}
	}

	return res, nil
}
