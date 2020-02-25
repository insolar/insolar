// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	pulse_core "github.com/insolar/insolar/pulse"
	"github.com/pkg/errors"
)

type SearchIndex struct {
	meta payload.Meta

	dep struct {
		indexes         object.IndexAccessor
		records         object.RecordAccessor
		pulseCalculator pulse.Calculator
		pulseStorage    pulse.Accessor
		sender          bus.Sender
	}
}

func (p *SearchIndex) Dep(
	indexes object.IndexAccessor,
	pulseCalculator pulse.Calculator,
	pulseStorage pulse.Accessor,
	records object.RecordAccessor,
	sender bus.Sender,
) {
	p.dep.indexes = indexes
	p.dep.sender = sender
	p.dep.pulseCalculator = pulseCalculator
	p.dep.records = records
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

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id": searchIndex.ObjectID.DebugString(),
		"until":     searchIndex.Until,
	})
	logger.Debug("search index. start to search index")

	if searchIndex.Until < pulse_core.MinTimePulse {
		return errors.New("searching index with until less than MinTimePulse is impossible")
	}

	currentP, err := p.dep.pulseStorage.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to fetch pulse")
	}
	currentPN := currentP.PulseNumber
	logger.Debug("search index. currentPN:", currentPN)

	// Until is above heavy's current pulse
	// It's impossible to find an index
	if currentPN < searchIndex.Until {
		logger.Warn("search index. currentPN < searchIndex.Until")

		msg, err := payload.NewMessage(&payload.SearchIndexInfo{})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.meta, msg)
		return nil
	}

	var idx *record.Index
	for currentPN >= searchIndex.Until {
		// searching for creation request
		reqID := *insolar.NewID(currentPN, searchIndex.ObjectID.Hash())
		_, err := p.dep.records.ForID(ctx, reqID)
		if err != nil && err != object.ErrNotFound {
			return errors.Wrapf(
				err,
				"failed to fetch object index for %v", *insolar.NewID(currentPN, searchIndex.ObjectID.Hash()),
			)
		}
		if err == nil {
			savedIdx, err := p.dep.indexes.LastKnownForID(ctx, reqID)
			if err != nil {
				return errors.Wrapf(
					err,
					"failed to fetch index for record %v", reqID.DebugString(),
				)
			}
			idx = &savedIdx
			break
		}
		logger.Debug("search index. didn't find for", currentPN)
		prev, err := p.dep.pulseCalculator.Backwards(ctx, currentPN, 1)
		if err != nil {
			return errors.Wrapf(
				err,
				"failed to fetch previous pulse for %v", currentPN,
			)
		}
		currentPN = prev.PulseNumber
		logger.Debug("search index. nextPN", currentPN)
	}

	if idx != nil {
		logger.Debug("search index. found index:", idx.ObjID.DebugString())
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
