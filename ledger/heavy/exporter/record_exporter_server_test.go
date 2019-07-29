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
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"
)

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
		iter.read = 11

		hasNext := iter.HasNext(ctx)

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

		t.Run("no data in the current. has more synce pulses. returns true", func(t *testing.T) {
			pn := gen.PulseNumber()
			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(pn).Return(1, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, pn, 1).Return(insolar.Pulse{PulseNumber: 100}, nil)

			jetKeeper := executor.NewJetKeeperMock(t)
			jetKeeper.TopSyncPulseMock.Return(101)

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

	t.Run("returns err, if LastKnownPosition returns err", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(0, errors.New("some error"))

		iter := newRecordIterator(pn, 0, 0, positionAccessor, nil, nil, nil)

		_, err := iter.Next(ctx)

		require.Error(t, err)
	})

	t.Run("returns err, if AtPosition returns err", func(t *testing.T) {
		pn := gen.PulseNumber()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(10, nil)
		positionAccessor.AtPositionMock.Expect(pn, uint32(2)).Return(insolar.ID{}, store.ErrNotFound)

		iter := newRecordIterator(pn, 1, 0, positionAccessor, nil, nil, nil)

		_, err := iter.Next(ctx)

		require.Error(t, err)
		require.Equal(t, err.Error(), store.ErrNotFound.Error())
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
		require.Equal(t, err.Error(), store.ErrNotFound.Error())
	})

	t.Run("reading data works", func(t *testing.T) {
		pn := gen.PulseNumber()
		id := gen.ID()
		positionAccessor := object.NewRecordPositionAccessorMock(t)
		positionAccessor.LastKnownPositionMock.Expect(pn).Return(10, nil)
		positionAccessor.AtPositionMock.Expect(pn, uint32(2)).Return(id, nil)

		record := record.Material{
			JetID: gen.JetID(),
		}
		recordsAccessor := object.NewRecordAccessorMock(t)
		recordsAccessor.ForIDMock.Expect(ctx, id).Return(record, nil)

		iter := newRecordIterator(pn, 1, 0, positionAccessor, recordsAccessor, nil, nil)
		next, err := iter.Next(ctx)

		require.NoError(t, err)
		require.Equal(t, uint32(1), iter.read)
		require.Equal(t, pn, next.PulseNumber)
		require.Equal(t, uint32(2), next.RecordNumber)
		require.Equal(t, id, next.RecordID)
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
			require.Equal(t, err.Error(), store.ErrNotFound.Error())
		})

		t.Run("Changing pulse works successfully", func(t *testing.T) {
			firstPN := gen.PulseNumber()
			nextPN := gen.PulseNumber()
			id := gen.ID()

			positionAccessor := object.NewRecordPositionAccessorMock(t)
			positionAccessor.LastKnownPositionMock.Expect(firstPN).Return(5, nil)
			positionAccessor.AtPositionMock.Expect(nextPN, uint32(1)).Return(id, nil)

			record := record.Material{
				JetID: gen.JetID(),
			}
			recordsAccessor := object.NewRecordAccessorMock(t)
			recordsAccessor.ForIDMock.Expect(ctx, id).Return(record, nil)

			pulseCalculator := network.NewPulseCalculatorMock(t)
			pulseCalculator.ForwardsMock.Expect(ctx, firstPN, 1).Return(insolar.Pulse{PulseNumber: nextPN}, nil)

			iter := newRecordIterator(firstPN, 10, 0, positionAccessor, recordsAccessor, nil, pulseCalculator)

			next, err := iter.Next(ctx)

			require.NoError(t, err)
			require.Equal(t, nextPN, iter.currentPulse)
			require.Equal(t, uint32(1), iter.read)
			require.Equal(t, nextPN, next.PulseNumber)
			require.Equal(t, uint32(1), next.RecordNumber)
			require.Equal(t, id, next.RecordID)
			require.Equal(t, record, next.Record)
		})
	})
}

func TestRecordServer_Export(t *testing.T) {
	t.Parallel()

	t.Run("count can't be 0", func(t *testing.T) {
		server := &RecordServer{}

		res, err := server.Export(inslogger.TestContext(t), &GetRecords{Count: 0})

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("PulseNumber can't be more than TopSyncPulseNumber", func(t *testing.T) {
		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(insolar.PulseNumber(0))
		server := &RecordServer{
			jetKeeper: jetKeeper,
		}

		res, err := server.Export(inslogger.TestContext(t), &GetRecords{Count: 1, PulseNumber: insolar.FirstPulseNumber})

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("returns empty slice of records, if no records", func(t *testing.T) {
		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(insolar.FirstPulseNumber)

		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(tmpdir)
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordPosition := object.NewRecordPositionDB(db)

		recordServer := NewRecordServer(nil, recordPosition, nil, jetKeeper)

		res, err := recordServer.Export(inslogger.TestContext(t), &GetRecords{
			PulseNumber:  insolar.FirstPulseNumber,
			RecordNumber: 0,
			Count:        10,
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, 0, len(res.Records))
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
		Virtual: &virtRec,
		JetID:   gen.JetID(),
	}

	return materialRecord
}

func TestRecordServer_Export_Composite(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	// Pulses
	firstPN := gen.PulseNumber()
	secondPN := firstPN + 10

	// JetKeeper
	jetKeeper := executor.NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseMock.Return(secondPN)

	// IDs and Records
	firstID := gen.ID()
	firstID.SetPulse(firstPN)
	firstRec := getMaterialRecord()

	secondID := gen.ID()
	secondID.SetPulse(firstPN)
	secondRec := getMaterialRecord()

	thirdID := gen.ID()
	thirdID.SetPulse(secondPN)
	thirdRec := getMaterialRecord()

	// TempDB
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	pulseStorage := pulse.NewDB(db)
	recordStorage := object.NewRecordDB(db)
	recordPosition := object.NewRecordPositionDB(db)

	// Save records to DB
	err = recordStorage.Set(ctx, firstID, firstRec)
	require.NoError(t, err)
	err = recordPosition.IncrementPosition(firstID)
	require.NoError(t, err)

	err = recordStorage.Set(ctx, secondID, secondRec)
	require.NoError(t, err)
	err = recordPosition.IncrementPosition(secondID)
	require.NoError(t, err)

	err = recordStorage.Set(ctx, thirdID, thirdRec)
	require.NoError(t, err)
	err = recordPosition.IncrementPosition(thirdID)
	require.NoError(t, err)

	// Pulses
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: firstPN})
	require.NoError(t, err)
	err = pulseStorage.Append(ctx, insolar.Pulse{PulseNumber: secondPN})
	require.NoError(t, err)

	recordServer := NewRecordServer(pulseStorage, recordPosition, recordStorage, jetKeeper)

	t.Run("export 1 of 3. first pulse", func(t *testing.T) {
		res, err := recordServer.Export(ctx, &GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 0,
			Count:        1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Records))

		resRecord := res.Records[0]
		require.Equal(t, firstPN, resRecord.PulseNumber)
		require.Equal(t, uint32(1), resRecord.RecordNumber)
		require.Equal(t, firstID, resRecord.RecordID)
		require.Equal(t, firstRec, resRecord.Record)
	})

	t.Run("export 1 of 3. second pulse", func(t *testing.T) {
		res, err := recordServer.Export(ctx, &GetRecords{
			PulseNumber:  secondPN,
			RecordNumber: 0,
			Count:        1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Records))

		resRecord := res.Records[0]
		require.Equal(t, secondPN, resRecord.PulseNumber)
		require.Equal(t, uint32(1), resRecord.RecordNumber)
		require.Equal(t, thirdID, resRecord.RecordID)
		require.Equal(t, thirdRec, resRecord.Record)
	})

	t.Run("export 3 of 3. first pulse", func(t *testing.T) {
		res, err := recordServer.Export(ctx, &GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 0,
			Count:        5,
		})
		require.NoError(t, err)
		require.Equal(t, 3, len(res.Records))
	})

	t.Run("export 2d. first pulse, set previousRecordNumber", func(t *testing.T) {
		res, err := recordServer.Export(ctx, &GetRecords{
			PulseNumber:  firstPN,
			RecordNumber: 1,
			Count:        1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(res.Records))

		resRecord := res.Records[0]
		require.Equal(t, secondPN, resRecord.PulseNumber)
		require.Equal(t, uint32(1), resRecord.RecordNumber)
		require.Equal(t, firstID, resRecord.RecordID)
		require.Equal(t, firstRec, resRecord.Record)
	})

}
