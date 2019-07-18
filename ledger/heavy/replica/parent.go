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
	Subscribe(context.Context, Target, Page) error
	Pull(context.Context, Page) ([]byte, uint32, error)
}

func NewParent() Parent {
	return &parent{}
}

type parent struct {
	Sequencer sequence.Sequencer  `inject:""`
	JetKeeper JetKeeper           `inject:""`
	Provider  intergrity.Provider `inject:""`
}

func (p *parent) Subscribe(ctx context.Context, target Target, at Page) error {
	highest := p.JetKeeper.TopSyncPulse()

	if at.Pulse > highest {
		p.JetKeeper.Subscribe(at.Pulse, func(highest insolar.PulseNumber) {
			p.notify(target, highest)
		})
		return nil
	}
	p.notify(target, highest)
	return nil
}

func (p *parent) Pull(ctx context.Context, page Page) ([]byte, uint32, error) {
	highest := p.JetKeeper.TopSyncPulse()
	if page.Pulse > highest {
		packet := p.Provider.Wrap([]sequence.Item{})
		return packet, 0, nil
	}
	items := p.Sequencer.Slice(page.Scope, page.Pulse, page.Skip, page.Limit)
	packet := p.Provider.Wrap(items)
	total := p.Sequencer.Len(page.Scope, page.Pulse)
	return packet, total, nil
}

func (p *parent) notify(target Target, pn insolar.PulseNumber) {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx)
	err := target.Notify(ctx, pn)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"target": target,
		}).Error(errors.Wrapf(err, "failed to notify target"))
		return
	}
}
