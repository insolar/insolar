package exporter

import (
	"context"
	"time"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
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
	authCfg         configuration.Auth
}

func NewRecordServer(
	pulseCalculator insolarPulse.Calculator,
	recordIndex object.RecordPositionAccessor,
	recordAccessor object.RecordAccessor,
	jetKeeper executor.JetKeeper,
	authCfg configuration.Auth,
) *RecordServer {
	return &RecordServer{
		pulseCalculator: pulseCalculator,
		recordIndex:     recordIndex,
		recordAccessor:  recordAccessor,
		jetKeeper:       jetKeeper,
		authCfg:         authCfg,
	}
}

func (r *RecordServer) Export(getRecords *GetRecords, stream RecordExporter_ExportServer) error {
	ctx := stream.Context()
	read := 0
	exportStart := time.Now()
	logger := inslogger.FromContext(ctx)
	logger.Info("Incoming request: ", getRecords.String())

	defer func(ctx context.Context) {
		stats.Record(
			addTagsForExporterMethodTiming(r.authCfg.Required, ctx, "record-export"),
			HeavyExporterMethodTiming.M(float64(time.Since(exportStart).Nanoseconds())/1e6),
		)
		logger.Infof("exported %d record", read)
	}(ctx)

	if getRecords.Count == 0 {
		return ErrNilCount
	}

	if getRecords.PulseNumber != 0 {
		topPulse := r.jetKeeper.TopSyncPulse()
		if topPulse < getRecords.PulseNumber {
			return ErrNotFinalPulseData
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
			if stream.Context().Err() != context.Canceled {
				logger.Error(err)
			}

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
