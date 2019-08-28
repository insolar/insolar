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
	"io/ioutil"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils/network"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

func TestRecordIterator_HasNext(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	t.Run("returns false, if LastKnownPosition returns error", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(0, errors.New("some error"))

		iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, nil, nil)

		hasNext := iter.HasNext(ctx)

		require.False(t, hasNext)
	})

	t.Run("returns false, if read all the count", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(156, nil)

		iter := newRecordIterator(pn, 0, 10, positionAccessor, nil, nil, nil)
		// bigger case
		iter.read = 11

		hasNext := iter.HasNext(ctx)

		require.False(t, hasNext)

		iter = newRecordIterator(pn, 0, 10, positionAccessor, nil, nil, nil)
		// equal case
		iter.read = 10

		hasNext = iter.HasNext(ctx)

		require.False(t, hasNext)
	})

	t.Run("returns true, if read not all the count", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(156, nil)

		iter := newRecordIterator(pn, 0, 10, positionAccessor, nil, nil, nil)
		iter.read = 9

		hasNext := iter.HasNext(ctx)

		require.True(t, hasNext)
	})

	t.Run("cross-pulse situations", func(t *testing.T) {
		t.Run("no data in the current.no further pulses. returns false", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{}, store.ErrNotFound)

			iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, nil, pulseCalculator)
			iter.currentPosition = 2

			hasNext := iter.HasNext(ctx)

			require.False(t, hasNext)
		})

		t.Run("no data in the current.no more synced pulses. returns false", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{PulseNumber: 100}, nil)

			jetKeeper := executor.NewJetKeeperMock(t)
			jetKeeper.TopSyncPulseMock.Return(99)

			iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, jetKeeper, pulseCalculator)
			iter.currentPosition = 2

			hasNext := iter.HasNext(ctx)

			require.False(t, hasNext)
		})

		t.Run("no data in the current. has more synced pulses. returns true", func(t *testing.T) {
			pn := insolar.PulseNumber(99)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{PulseNumber: 100}, nil)

			jetKeeper := executor.NewJetKeeperMock(t)
			jetKeeper.TopSyncPulseMock.Return(101)

			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.When(99).Then(2, nil)
			positionAccessor.LastKnownPositionMock.Expect(100).Return(1, nil)

			iter := newRecordIterator(pn, 2, 0, positionAccessor, nil, jetKeeper, pulseCalculator)
			iter.read = 10
			iter.needToRead = 100

			hasNext := iter.HasNext(ctx)

			require.True(t, hasNext)
		})

		t.Run("no data in the current. has more synce pulses. returns false, because read everything", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{PulseNumber: 100}, nil)

			jetKeeper := executor.NewJetKeeperMock(t)
			jetKeeper.TopSyncPulseMock.Return(101)

			iter := newRecordIterator(pn, 2, 0, positionAccessor, nil, jetKeeper, pulseCalculator)
			iter.read = 10
			iter.needToRead = 10

			hasNext := iter.HasNext(ctx)

			require.False(t, hasNext)
		})

	})
}

func TestRecordIterator_Next(t *testing.T) {
	ctx := inslogger.TestContext(t)

	t.Run("returns err, if AtPosition returns err", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(10, nil)
		positionAccessor.AtPositionMock.Expect(pn, uint32(2)).Return(insolar.ID{}, store.ErrNotFound)

		iter := newRecordIterator(pn, 1, 0, positionAccessor, nil, nil, nil)

		_, err := iter.Next(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), store.ErrNotFound.Error())
	})

	t.Run("returns err, if ForID returns err", func(t *testing.T) {
		pn := gen.PulseNumber()
		id := gen.ID()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(10, nil)
		positionAccessor.AtPositionMock.Expect(pn, uint32(2)).Return(id, nil)

		recordsAccessor := object.NewRecordAccessorMock(t)
		recordsAccessor.ForIDMock.Expect(ctx, id).Return(record.Material{}, store.ErrNotFound)

		iter := newRecordIterator(pn, 1, 0, positionAccessor, recordsAccessor, nil, nil)

		_, err := iter.Next(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), store.ErrNotFound.Error())
	})

	t.Run("reading data works", func(t *testing.T) {
		pn := gen.PulseNumber()
		id := gen.IDWithPulse(pn)

		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(10, nil)
		positionAccessor.AtPositionMock.Expect(pn, uint32(2)).Return(id, nil)

		record := record.Material{
			JetID: gen.JetID(),
			ID:    id,
		}
		recordsAccessor := object.NewRecordAccessorMock(t)
		recordsAccessor.ForIDMock.Expect(ctx, id).Return(record, nil)

		iter := newRecordIterator(pn, 1, 0, positionAccessor, recordsAccessor, nil, nil)
		next, err := iter.Next(ctx)

		require.NoError(t, err)
		require.Equal(t, uint32(1), iter.read)
		require.Equal(t, pn, next.Record.ID.Pulse())
		require.Equal(t, uint32(2), next.RecordNumber)
		require.Equal(t, id, next.Record.ID)
		require.Equal(t, record, next.Record)
	})

	t.Run("cross-pulse edges", func(t *testing.T) {
		t.Run("Forwards returns error", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{}, store.ErrNotFound)

			iter := newRecordIterator(pn, 1, 0, positionAccessor, nil, nil, pulseCalculator)

			_, err := iter.Next(ctx)

			require.Error(t, err)
			require.Contains(t, err.Error(), store.ErrNotFound.Error())
		})

		t.Run("Changing pulse works successfully", func(t *testing.T) {
			firstPN := gen.PulseNumber()
			nextPN := firstPN + 10
			id := gen.IDWithPulse(nextPN)

			jetKeeper := executor.NewJetKeeperMock(t)
			jetKeeper.TopSyncPulseMock.Return(nextPN)

			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.When(firstPN).Then(5, nil)
			positionAccessor.LastKnownPositionMock.When(nextPN).Then(1, nil)

			positionAccessor.AtPositionMock.Expect(nextPN, uint32(1)).Return(id, nil)

			rec := record.Material{
				JetID: gen.JetID(),
				ID:    id,
			}
			recordsAccessor := object.NewRecordAccessorMock(t)
			recordsAccessor.ForIDMock.Expect(ctx, id).Return(rec, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, firstPN, 1).Return(insolar.Pulse{PulseNumber: nextPN}, nil)

			iter := newRecordIterator(firstPN, 10, 0, positionAccessor, recordsAccessor, jetKeeper, pulseCalculator)

			next, err := iter.Next(ctx)

			require.NoError(t, err)
			require.Equal(t, nextPN, iter.currentPulse)
			require.Equal(t, uint32(1), iter.read)
			require.Equal(t, nextPN, next.Record.ID.Pulse())
			require.Equal(t, uint32(1), next.RecordNumber)
			require.Equal(t, id, next.Record.ID)
			require.Equal(t, rec, next.Record)
		})
	})
}

type streamMock struct {
	checker func(*Record) error
}

func (s streamMock) Send(rec *Record) error {
	return s.checker(rec)
}

func (s streamMock) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (s streamMock) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (s streamMock) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (s streamMock) Context() context.Context {
	return context.Background()
}

func (s streamMock) SendMsg(m interface{}) error {
	panic("implement me")
}

func (s streamMock) RecvMsg(m interface{}) error {
	panic("implement me")
}

func TestRecordServer_Export(t *testing.T) {
	t.Parallel()

	t.Run("count can't be 0", func(t *testing.T) {
		server := &RecordServer{}

		err := server.Export(&GetRecords{Count: 0}, &streamMock{})

		require.Error(t, err)
	})

	t.Run("PulseNumber can't be more than TopSyncPulseNumber", func(t *testing.T) {
		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(insolar.PulseNumber(0))
		server := &RecordServer{
			jetKeeper: jetKeeper,
		}

		err := server.Export(&GetRecords{Count: 1, PulseNumber: pulse.MinTimePulse}, &streamMock{})

		require.Error(t, err)
	})

	t.Run("returns empty slice of records, if no records", func(t *testing.T) {
		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse)

		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordPosition := object.NewRecordDB(db)
		pulses := insolarPulse.NewDB(db)

		recordServer := NewRecordServer(pulses, recordPosition, nil, jetKeeper)

		streamMock := &streamMock{checker: func(i *Record) error {
			t.Error("it shouldn't be called")
			return nil
		}}

		err = recordServer.Export(&GetRecords{
			PulseNumber:  pulse.MinTimePulse,
			RecordNumber: 0,
			Count:        10,
		}, streamMock)

		require.NoError(t, err)
	})
}

// getVirtualRecord generates random Virtual record
func getVirtualRecord() record.Virtual {
	var requestRecord record.IncomingRequest

	obj := gen.Reference()
	requestRecord.Object = &obj

	virtualRecord := record.Virtual{
		Union: &record.Virtual_IncomingRequest{
			IncomingRequest: &requestRecord,
		},
	}

	return virtualRecord
}

// getMaterialRecord generates random Material record
func getMaterialRecord() record.Material {
	virtRec := getVirtualRecord()

	materialRecord := record.Material{
		Virtual: virtRec,
		JetID:   gen.JetID(),
	}

	return materialRecord
}

func TestRecordServer_Export_Composite(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	// Pulses
	firstPN := insolar.PulseNumber(pulse.MinTimePulse + 100)
	secondPN := insolar.PulseNumber(firstPN + 10)

	// JetKeeper
	jetKeeper := executor.NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseMock.Return(secondPN)

	// IDs and Records
	firstID := gen.IDWithPulse(firstPN)
	firstRec := getMaterialRecord()
	firstRec.ID = firstID

	secondID := gen.IDWithPulse(firstPN)
	secondRec := getMaterialRecord()
	secondRec.ID = secondID

	thirdID := gen.IDWithPulse(secondPN)
	thirdRec := getMaterialRecord()
	thirdRec.ID = thirdID

	// TempDB
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	pulseStorage := insolarPulse.NewDB(db)
	recordStorage := object.NewRecordDB(db)
	recordPosition := object.NewRecordDB(db)

	// Save records to DB
	err = recordStorage.Set(ctx, firstRec)
	require.NoError(t, err)

	err = recordStorage.Set(ctx, secondRec)
	require.NoError(t, err)

	err = recordStorage.Set(ctx, thirdRec)
	require.NoError(t, err)

	// Pulses

	// Trash pulses without data
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: pulse.MinTimePulse})
	require.NoError(t, err)
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: pulse.MinTimePulse + 10})
	require.NoError(t, err)
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: pulse.MinTimePulse + 20})
	require.NoError(t, err)

	// LegalInfo
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: firstPN})
	require.NoError(t, err)
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: secondPN})
	require.NoError(t, err)

	recordServer := NewRecordServer(pulseStorage, recordPosition, recordStorage, jetKeeper)

	t.Run("export 1 of 3. first pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 0,
			Count:        1,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 1, len(recs))

		resRecord := recs[0]
		require.Equal(t, firstPN, resRecord.Record.ID.Pulse())
		require.Equal(t, uint32(1), resRecord.RecordNumber)
		require.Equal(t, firstID, resRecord.Record.ID)
		require.Equal(t, firstRec, resRecord.Record)
	})

	t.Run("export 1 of 3. second pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  secondPN,
			RecordNumber: 0,
			Count:        1,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 1, len(recs))

		resRecord := recs[0]
		require.Equal(t, secondPN, resRecord.Record.ID.Pulse())
		require.Equal(t, uint32(1), resRecord.RecordNumber)
		require.Equal(t, thirdID, resRecord.Record.ID)
		require.Equal(t, thirdRec, resRecord.Record)
	})

	t.Run("export 3 of 3. first pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 0,
			Count:        5,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 3, len(recs))
	})

	t.Run("export 3 of 3. zero pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  0,
			RecordNumber: 0,
			Count:        5,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 3, len(recs))
	})

	t.Run("export 2d. first pulse, set previousRecordNumber", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 1,
			Count:        1,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 1, len(recs))

		resRecord := recs[0]
		require.Equal(t, firstPN, resRecord.Record.ID.Pulse())
		require.Equal(t, uint32(2), resRecord.RecordNumber)
		require.Equal(t, secondID, resRecord.Record.ID)
		require.Equal(t, secondRec, resRecord.Record)
	})

}

func TestRecordServer_Export_Composite_BatchVersion(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	// Pulses
	firstPN := insolar.PulseNumber(pulse.MinTimePulse + 100)
	secondPN := insolar.PulseNumber(firstPN + 10)

	// JetKeeper
	jetKeeper := executor.NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseMock.Return(secondPN)

	// IDs and Records
	firstID := *insolar.NewID(firstPN, []byte{1})
	firstRec := getMaterialRecord()
	firstRec.ID = firstID

	secondID := *insolar.NewID(firstPN, []byte{2})
	secondRec := getMaterialRecord()
	secondRec.ID = secondID

	thirdID := *insolar.NewID(secondPN, []byte{1})
	thirdRec := getMaterialRecord()
	thirdRec.ID = thirdID

	// TempDB
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	pulseStorage := insolarPulse.NewDB(db)
	recordStorage := object.NewRecordDB(db)
	recordPosition := object.NewRecordDB(db)

	// Save records to DB
	err = recordStorage.BatchSet(ctx, []record.Material{firstRec, secondRec, thirdRec})
	require.NoError(t, err)

	// Pulses

	// Trash pulses without data
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: pulse.MinTimePulse})
	require.NoError(t, err)
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: pulse.MinTimePulse + 10})
	require.NoError(t, err)
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: pulse.MinTimePulse + 20})
	require.NoError(t, err)

	// LegalInfo
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: firstPN})
	require.NoError(t, err)
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: secondPN})
	require.NoError(t, err)

	recordServer := NewRecordServer(pulseStorage, recordPosition, recordStorage, jetKeeper)

	t.Run("export 1 of 3. first pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 0,
			Count:        1,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 1, len(recs))

		resRecord := recs[0]
		require.Equal(t, firstPN, resRecord.Record.ID.Pulse())
		require.Equal(t, uint32(1), resRecord.RecordNumber)
		require.Equal(t, firstID, resRecord.Record.ID)
		require.Equal(t, firstRec, resRecord.Record)
	})

	t.Run("export 1 of 3. second pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  secondPN,
			RecordNumber: 0,
			Count:        1,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 1, len(recs))

		resRecord := recs[0]
		require.Equal(t, secondPN, resRecord.Record.ID.Pulse())
		require.Equal(t, uint32(1), resRecord.RecordNumber)
		require.Equal(t, thirdID, resRecord.Record.ID)
		require.Equal(t, thirdRec, resRecord.Record)
	})

	t.Run("export 3 of 3. first pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 0,
			Count:        5,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 3, len(recs))
	})

	t.Run("export 3 of 3. zero pulse", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  0,
			RecordNumber: 0,
			Count:        5,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 3, len(recs))
	})

	t.Run("export 2d. first pulse, set previousRecordNumber", func(t *testing.T) {
		var recs []*Record
		streamMock := &streamMock{checker: func(i *Record) error {
			recs = append(recs, i)
			return nil
		}}

		err := recordServer.Export(&GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 1,
			Count:        1,
		}, streamMock)
		require.NoError(t, err)
		require.Equal(t, 1, len(recs))

		resRecord := recs[0]
		require.Equal(t, firstPN, resRecord.Record.ID.Pulse())
		require.Equal(t, uint32(2), resRecord.RecordNumber)
		require.Equal(t, secondID, resRecord.Record.ID)
		require.Equal(t, secondRec, resRecord.Record)
	})

}
