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

package replica

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/replica/intergrity"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

type Parent interface {
	Subscribe(Target, Position) error
	Pull(byte, Position, uint32) ([]byte, error)
}

func NewParent(Sequencer sequence.Sequencer, keeper JetKeeper, cryptoService insolar.CryptographyService) Parent {
	provider := intergrity.NewProvider(cryptoService)
	return &localParent{Sequencer: Sequencer, jetKeeper: keeper, provider: provider}
}

type localParent struct {
	Sequencer sequence.Sequencer
	jetKeeper JetKeeper
	provider  intergrity.Provider
}

func (p *localParent) Subscribe(child Target, at Position) error {
	logger := inslogger.FromContext(context.Background())
	current := p.jetKeeper.TopSyncPulse()

	if current < at.Pulse {
		// TODO: register handler on sync pulse update
		logger.Warn("I'm replicaroot. Current pulse less than requested position.")
	}
	err := child.Notify()
	if err != nil {
		// TODO: retry notify
		return errors.Wrapf(err, "failed to notify child")
	}
	return nil
}

func (p *localParent) Pull(scope byte, from Position, limit uint32) ([]byte, error) {
	logger := inslogger.FromContext(context.Background())
	highestPulse := p.jetKeeper.TopSyncPulse()
	items := p.Sequencer.Slice(scope, from.Pulse, from.Skip, highestPulse, limit)
	logger.Warnf("PULL_BATCH slicing scope: %v len(items): %v from: %v skip: %v highest: %v limit: %v", scope, len(items), from.Pulse, from.Skip, highestPulse, limit)
	packet := p.provider.Wrap(items)
	return packet, nil
}
