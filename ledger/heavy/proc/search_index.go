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

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"
	pulse_core "github.com/insolar/insolar/pulse"
	"github.com/pkg/errors"
)

type SearchIndex struct {
	meta payload.Meta

	dep struct {
		indexes         object.IndexAccessor
		pulseCalculator pulse.Calculator
		pulseStorage    pulse.Accessor
		sender          bus.Sender
	}
}

func (p *SearchIndex) Dep(
	indexes object.IndexAccessor,
	pulseCalculator pulse.Calculator,
	pulseStorage pulse.Accessor,
	sender bus.Sender,
) {
	p.dep.indexes = indexes
	p.dep.sender = sender
	p.dep.pulseCalculator = pulseCalculator
	p.dep.pulseStorage = pulseStorage
}

func NewSearchIndex(meta payload.Meta) *SearchIndex {
	return &SearchIndex{
		meta: meta,
	}
}

func (p *SearchIndex) Proceed(ctx context.Context) error {
	searchIndex := payload.SearchIndex{}
	err := searchIndex.Unmarshal(p.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal searchIndex message")
	}

	if searchIndex.Until < pulse_core.MinTimePulse {
		return errors.New("searching index with until less than MinTimePulse is impossible")
	}

	currentP, err := p.dep.pulseStorage.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to fetch pulse")
	}
	currentPN := currentP.PulseNumber

	// Until is above heavy's current pulse
	// It's impossible to find an index
	if currentPN < searchIndex.Until {
		msg, err := payload.NewMessage(&payload.SearchIndexInfo{})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.meta, msg)
		return nil
	}

	var idx *record.Index
	for currentPN >= searchIndex.Until {
		savedIdx, err := p.dep.indexes.ForID(ctx, currentPN, *insolar.NewID(currentPN, searchIndex.ObjectID.Hash()))
		if err != nil && err != object.ErrIndexNotFound {
			return errors.Wrapf(
				err,
				"failed to fetch object index for %v", *insolar.NewID(currentPN, searchIndex.ObjectID.Hash()),
			)
		}
		if err == nil {
			idx = &savedIdx
			break
		}
		prev, err := p.dep.pulseCalculator.Backwards(ctx, currentPN, 1)
		if err != nil {
			return errors.Wrapf(
				err,
				"failed to fetch previous pulse for %v", currentPN,
			)
		}
		currentPN = prev.PulseNumber
	}

	msg, err := payload.NewMessage(&payload.SearchIndexInfo{
		Index: idx,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.meta, msg)
	return nil
}
