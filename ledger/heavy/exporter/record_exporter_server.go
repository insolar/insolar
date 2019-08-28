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

	"github.com/insolar/insolar/insolar"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"

	"github.com/pkg/errors"
)

type RecordServer struct {
	pulseCalculator insolarPulse.Calculator
	recordIndex     object.RecordPositionAccessor
	recordAccessor  object.RecordAccessor
	jetKeeper       executor.JetKeeper
}

func NewRecordServer(
	pulseCalculator insolarPulse.Calculator,
	recordIndex object.RecordPositionAccessor,
	recordAccessor object.RecordAccessor,
	jetKeeper executor.JetKeeper,
) *RecordServer {
	return &RecordServer{
		pulseCalculator: pulseCalculator,
		recordIndex:     recordIndex,
		recordAccessor:  recordAccessor,
		jetKeeper:       jetKeeper,
	}
}

func (r *RecordServer) Export(getRecords *GetRecords, stream RecordExporter_ExportServer) error {
	ctx := stream.Context()
	logger := inslogger.FromContext(ctx)

	if getRecords.Count == 0 {
		return errors.New("count can't be 0")
	}

	if getRecords.PulseNumber != 0 {
		topPulse := r.jetKeeper.TopSyncPulse()
		if topPulse < getRecords.PulseNumber {
			return errors.New("trying to get a non-finalized pulse data")
		}
	} else {
		getRecords.PulseNumber = pulse.MinTimePulse
	}

	iter := newRecordIterator(
		getRecords.PulseNumber,
		getRecords.RecordNumber,
		getRecords.Count,
		r.recordIndex,
		r.recordAccessor,
		r.jetKeeper,
		r.pulseCalculator,
	)

	for iter.HasNext(stream.Context()) {
		record, err := iter.Next(stream.Context())
		if err != nil {
			logger.Error(err)
			return err
		}

		err = stream.Send(record)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

type recordIterator struct {
	currentPosition uint32
	currentPulse    insolar.PulseNumber

	read       uint32
	needToRead uint32

	recordIndex     object.RecordPositionAccessor
	recordAccessor  object.RecordAccessor
	jetKeeper       executor.JetKeeper
	pulseCalculator insolarPulse.Calculator
}

func newRecordIterator(
	pn insolar.PulseNumber,
	lastPosition uint32,
	takeCount uint32,
	recordIndex object.RecordPositionAccessor,
	recordAccessor object.RecordAccessor,
	jetKeeper executor.JetKeeper,
	pulseCalculator insolarPulse.Calculator,
) *recordIterator {
	return &recordIterator{
		needToRead:      takeCount,
		currentPosition: lastPosition,
		currentPulse:    pn,
		recordIndex:     recordIndex,
		recordAccessor:  recordAccessor,
		jetKeeper:       jetKeeper,
		pulseCalculator: pulseCalculator,
	}
}

func (r *recordIterator) HasNext(ctx context.Context) bool {
	if r.read >= r.needToRead {
		return false
	}

	lastKnown, err := r.recordIndex.LastKnownPosition(r.currentPulse)
	if err != nil {
		return r.checkNextPulse(ctx)
	}

	if lastKnown < r.currentPosition+1 {
		return r.checkNextPulse(ctx)
	}

	return true
}

func (r *recordIterator) checkNextPulse(ctx context.Context) bool {
	currentPulse := r.currentPulse

	for {
		nextPulse, err := r.pulseCalculator.Forwards(ctx, currentPulse, 1)
		if err != nil {
			return false
		}
		topPulse := r.jetKeeper.TopSyncPulse()
		if topPulse < nextPulse.PulseNumber {
			return false
		}
		_, err = r.recordIndex.LastKnownPosition(nextPulse.PulseNumber)
		if err != nil {
			currentPulse = nextPulse.PulseNumber
		} else {
			return true
		}
	}
}

func (r *recordIterator) Next(ctx context.Context) (*Record, error) {
	r.currentPosition++

	lastKnown, err := r.recordIndex.LastKnownPosition(r.currentPulse)
	if err != nil || lastKnown < r.currentPosition {
		err := r.setNextPulse(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "iterator failed to change pulse")
		}
	}

	id, err := r.recordIndex.AtPosition(r.currentPulse, r.currentPosition)
	if err != nil {
		return nil, errors.Wrap(err, "iterator failed to find record's position")
	}

	rec, err := r.recordAccessor.ForID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "iterator failed to find record")
	}

	r.read++

	return &Record{
		RecordNumber: r.currentPosition,
		Record:       rec,
	}, nil
}

func (r *recordIterator) setNextPulse(ctx context.Context) error {
	currentPulse := r.currentPulse

	for {
		nextPulse, err := r.pulseCalculator.Forwards(ctx, currentPulse, 1)
		if err != nil {
			return err
		}
		topPulse := r.jetKeeper.TopSyncPulse()
		if topPulse < nextPulse.PulseNumber {
			return errors.New("there are no synced pulses")
		}
		_, err = r.recordIndex.LastKnownPosition(nextPulse.PulseNumber)
		if err != nil {
			currentPulse = nextPulse.PulseNumber
		} else {
			r.currentPulse = nextPulse.PulseNumber
			r.currentPosition = 1
			return nil
		}
	}
}
