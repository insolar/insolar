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

package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
)

type nodeCalculator interface {
	IsBeyondLimit(ctx context.Context, currentPN, targetPN core.PulseNumber) (bool, error)
	NodeForJet(ctx context.Context, jetID core.RecordID, rootPN, targetPN core.PulseNumber) (*core.RecordRef, error)
}

type nodeCalculatorConcrete struct {
	storage.PulseTracker `inject:""`
	core.JetCoordinator  `inject:""`
	LightChainLimit      int
}

// NewNodeCalculatorConcrete is a constructor for a helper-service nodeCalculator
func NewNodeCalculatorConcrete(lightChainLimit int) *nodeCalculatorConcrete {
	return &nodeCalculatorConcrete{LightChainLimit: lightChainLimit}
}

// IsBeyondLimit calculates if target pulse is behind clean-up limit
func (n *nodeCalculatorConcrete) IsBeyondLimit(ctx context.Context, currentPN, targetPN core.PulseNumber) (bool, error) {
	currentPulse, err := n.PulseTracker.GetPulse(ctx, currentPN)
	if err != nil {
		return false, errors.Wrapf(err, "failed to fetch pulse %v", currentPN)
	}
	targetPulse, err := n.PulseTracker.GetPulse(ctx, targetPN)
	if err != nil {
		return false, errors.Wrapf(err, "failed to fetch pulse %v", targetPN)
	}

	if currentPulse.SerialNumber-targetPulse.SerialNumber < n.LightChainLimit {
		return false, nil
	}

	return true, nil
}

// NodeForJet calculates a node for a specific jet for a specific pulseNumber
func (n *nodeCalculatorConcrete) NodeForJet(ctx context.Context, jetID core.RecordID, rootPN, targetPN core.PulseNumber) (*core.RecordRef, error) {
	toHeavy, err := n.IsBeyondLimit(ctx, rootPN, targetPN)
	if err != nil {
		return nil, err
	}

	if toHeavy {
		return n.JetCoordinator.Heavy(ctx, rootPN)
	}
	return n.JetCoordinator.LightExecutorForJet(ctx, jetID, targetPN)
}
