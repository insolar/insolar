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
	"sync"
	"time"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"

	"github.com/pkg/errors"
)

type OneRequestLimiter struct {
	durationBetweenReq time.Duration
	lastRequest        *time.Time
	lock               sync.Mutex
}

func NewOneRequestLimiter(durationBetweenReq time.Duration) *OneRequestLimiter {
	return &OneRequestLimiter{
		durationBetweenReq: durationBetweenReq,
	}
}

func (l *OneRequestLimiter) Take(ctx context.Context) {
	l.lock.Lock()
	defer l.lock.Unlock()
	t := time.Now()

	if l.lastRequest == nil {
		l.lastRequest = &t
		return
	}

	nextReq := l.lastRequest.Add(l.durationBetweenReq)
	if time.Now().After(nextReq) {
		l.lastRequest = &t
		return
	}

	<-time.After(l.lastRequest.Add(l.durationBetweenReq).Sub(t))
}

type RecordServer struct {
	pulseCalculator insolarPulse.Calculator
	recordIndex     object.RecordPositionAccessor
	recordAccessor  object.RecordAccessor
	jetKeeper       executor.JetKeeper
	limiter         *OneRequestLimiter
	// Number of pulses after which client can see finalized pulse
	exportDelay int
}

func NewRecordServer(
	pulseCalculator insolarPulse.Calculator,
	recordIndex object.RecordPositionAccessor,
	recordAccessor object.RecordAccessor,
	jetKeeper executor.JetKeeper,
	limiter *OneRequestLimiter,
	exportDelay int,
) *RecordServer {
	return &RecordServer{
		pulseCalculator: pulseCalculator,
		recordIndex:     recordIndex,
		recordAccessor:  recordAccessor,
		jetKeeper:       jetKeeper,
		limiter:         limiter,
		exportDelay:     exportDelay,
	}
}

func (r *RecordServer) Export(getRecords *GetRecords, stream RecordExporter_ExportServer) error {
	r.limiter.Take(stream.Context())

	ctx := stream.Context()

	exportStart := time.Now()
	defer func(ctx context.Context) {
		stats.Record(
			insmetrics.InsertTag(ctx, TagHeavyExporterMethodName, "record-export"),
			HeavyExporterMethodTiming.M(float64(time.Since(exportStart).Nanoseconds())/1e6),
		)
	}(ctx)

	logger := inslogger.FromContext(ctx)
	logger.Info("Incoming request: ", getRecords.String())

	if getRecords.Count == 0 {
		return errors.New("count can't be 0")
	}

	if getRecords.PulseNumber != 0 {
		topPulse := r.jetKeeper.TopSyncPulse()
		lastPossiblePulseToExport, err := r.pulseCalculator.Backwards(ctx, topPulse, r.exportDelay)
		if err != nil {
			if err == insolarPulse.ErrNotFound {
				return errors.Wrap(err, "trying to get a non-finalized pulse data")
			}
			return errors.Wrap(err, "failed to get previous pulse")
		}
		if lastPossiblePulseToExport.PulseNumber < getRecords.PulseNumber {
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
		r.exportDelay,
	)
	read := 0

	var numSent int
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
		read++
	}

	if read == 0 {
		topPulse := r.jetKeeper.TopSyncPulse()
		err := stream.Send(&Record{
			ShouldIterateFrom: &topPulse,
		})
		if err != nil {
			logger.Error(err)
			return err
		}
		numSent++
	}
	logger.Info("exported %d record", numSent)

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
	exportDelay     int
}

func newRecordIterator(
	pn insolar.PulseNumber,
	lastPosition uint32,
	takeCount uint32,
	recordIndex object.RecordPositionAccessor,
	recordAccessor object.RecordAccessor,
	jetKeeper executor.JetKeeper,
	pulseCalculator insolarPulse.Calculator,
	exportDelay int,
) *recordIterator {
	return &recordIterator{
		needToRead:      takeCount,
		currentPosition: lastPosition,
		currentPulse:    pn,
		recordIndex:     recordIndex,
		recordAccessor:  recordAccessor,
		jetKeeper:       jetKeeper,
		pulseCalculator: pulseCalculator,
		exportDelay:     exportDelay,
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

	topPulse := r.jetKeeper.TopSyncPulse()
	lastPossiblePulseToExport, err := r.pulseCalculator.Backwards(ctx, topPulse, r.exportDelay)
	if err != nil {
		return false
	}

	for {
		nextPulse, err := r.pulseCalculator.Forwards(ctx, currentPulse, 1)
		if err != nil {
			return false
		}

		if lastPossiblePulseToExport.PulseNumber < nextPulse.PulseNumber {
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

	topPulse := r.jetKeeper.TopSyncPulse()
	lastPossiblePulseToExport, err := r.pulseCalculator.Backwards(ctx, topPulse, r.exportDelay)
	if err != nil {
		return err
	}

	for {
		nextPulse, err := r.pulseCalculator.Forwards(ctx, currentPulse, 1)
		if err != nil {
			return err
		}

		if lastPossiblePulseToExport.PulseNumber < nextPulse.PulseNumber {
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
