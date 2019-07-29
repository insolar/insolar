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

package exporter

import (
	"context"
	"errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
)

type RecordServer struct {
	pulseCalculator pulse.Calculator
	recordIndex     object.RecordPositionAccessor
	recordAccessor  object.RecordAccessor
}

type recordIterator struct {
	currentPosition uint32
	currentPulse    insolar.PulseNumber

	read       uint32
	needToRead uint32

	recordIndex     object.RecordPositionAccessor
	recordAccessor  object.RecordAccessor
	jetKeeper       executor.JetKeeper
	pulseCalculator pulse.Calculator
}

func newRecordIterator(pn insolar.PulseNumber, lastPosition uint32, takeCount uint32) *recordIterator {
	return &recordIterator{
		needToRead:      takeCount,
		currentPosition: lastPosition,
		currentPulse:    pn,
	}
}

func (r *recordIterator) HasNext(ctx context.Context) bool {
	lastKnown, err := r.recordIndex.LastKnownPosition(r.currentPulse)
	if err != nil {
		return false
	}

	if lastKnown < r.currentPosition+1 {
		return r.checkNextPulse(ctx)
	}

	return r.read < r.needToRead
}

func (r *recordIterator) Next(ctx context.Context) (*Record, error) {
	r.currentPosition++

	lastKnown, err := r.recordIndex.LastKnownPosition(r.currentPulse)
	if err != nil {
		return nil, err
	}

	if lastKnown < r.currentPosition {
		err := r.setNextPulse(ctx)
		if err != nil {
			return nil, err
		}
	}

	id, err := r.recordIndex.AtPosition(r.currentPulse, r.currentPosition)
	if err != nil {
		return nil, err
	}

	rec, err := r.recordAccessor.ForID(ctx, id)
	if err != nil {
		return nil, err
	}

	r.read++

	return &Record{
		PulseNumber:  r.currentPulse,
		RecordNumber: r.currentPosition,
		RecordID:     id,
		Record:       rec,
	}, nil
}

func (r *recordIterator) checkNextPulse(ctx context.Context) bool {
	nextPulse, err := r.pulseCalculator.Forwards(ctx, r.currentPulse, 1)
	if err != nil {
		return false
	}
	topPulse := r.jetKeeper.TopSyncPulse()
	if topPulse < nextPulse.PulseNumber {
		return false
	}

	return true
}

func (r *recordIterator) setNextPulse(ctx context.Context) error {
	nextPulse, err := r.pulseCalculator.Forwards(ctx, r.currentPulse, 1)
	if err != nil {
		return err
	}

	r.currentPulse = nextPulse.PulseNumber
	r.currentPosition = 1

	return nil
}

func (r *RecordServer) Export(ctx context.Context, getRecords *GetRecords) (*Records, error) {
	if getRecords.Count == 0 {
		return nil, errors.New("count can't be 0")
	}

	if getRecords.PulseNumber == 0 {
		return r.fetchSince(ctx, insolar.FirstPulseNumber, getRecords.RecordNumber, getRecords.Count)
	}
}

func (r *RecordServer) fetchSince(ctx context.Context, pn insolar.PulseNumber, lastPosition uint32, takeCount uint32) (*Records, error) {
	var ids []object.OrderedID
	var err error
	currentPN := pn

	for uint32(len(ids)) < takeCount {
		takeDelta := takeCount - uint32(len(ids))

		ids, err = r.recordIndex.IDs(currentPN, lastPosition, takeDelta)
		if err != nil {
			return nil, err
		}

		if uint32(len(ids)) != takeCount {
			nextPulse, err := r.pulseCalculator.Forwards(ctx, currentPN, 1)
			if err != nil {
				inslogger.FromContext(ctx).Error(err)
				return r.loadRecords(ctx, ids)
			}
		}

	}
}

func (r *RecordServer) loadRecords(ctx context.Context, ids []object.OrderedID) (*Records, error) {
	result := &Records{}
	for _, id := range ids {
		rec, err := r.recordAccessor.ForID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		result.Records = append(result.Records, Record{
			PulseNumber:  id.ID.Pulse(),
			RecordNumber: id.Position,
			RecordID:     id.ID,
			Record:       rec,
		})
	}

	return result, nil
}
