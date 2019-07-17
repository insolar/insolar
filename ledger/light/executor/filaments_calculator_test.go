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

package executor

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestFilamentCalculatorDefault_PendingRequestsUnit(t *testing.T) {
	// t.Skip()
	ctx := inslogger.TestContext(t)
	mt := minimock.NewController(t)

	indexes := object.NewIndexAccessorMock(mt)
	records := object.NewRecordMemory()
	coordinator := jet.NewCoordinatorMock(mt)
	fetcher := jet.NewFetcherMock(mt)
	sender := bus.NewSenderMock(mt)
	calc := NewFilamentCalculator(indexes, records, coordinator, fetcher, sender)

	startPulse := insolar.PulseNumber(insolar.FirstPulseNumber + 1000)
	objectID := gen.IDWithPulse(startPulse)
	indexes.ForIDFunc = func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		if id != objectID {
			return record.Index{}, errors.Errorf("unexpected id %v in ForID call", id.DebugString())
		}
		// just stubs index to pass checks in PendingRequests method
		return record.Index{
			Lifeline: record.Lifeline{
				PendingPointer:      gen.IDRef(),
				EarliestOpenRequest: &startPulse,
			},
		}, nil
	}

	iterMock := newFetchIteratorMock()
	calc.iteratorMaker = fetchIteratorMockMaker(iterMock)

	// use pulse number just for uniqueness
	pn := startPulse
	inRequestWithID := func() (insolar.ID, record.Record) {
		pn++
		return gen.IDWithPulse(pn), record.IncomingRequest{}
	}
	outRequestWithID := func(reason insolar.ID, mode record.ReturnMode) (insolar.ID, record.Record) {
		pn++
		return gen.IDWithPulse(pn), record.OutgoingRequest{
			Reason:     *insolar.NewReference(reason),
			ReturnMode: mode,
		}
	}
	resultWithID := func(reqID insolar.ID) (insolar.ID, record.Record) {
		pn++
		return gen.IDWithPulse(pn), record.Result{Request: *insolar.NewReference(reqID)}
	}

	var (
		// first request
		inRequestID1, inRequest1 = inRequestWithID()
		// first not detached outgoing with reason='first request'
		outRequestID1, outRequest1 = outRequestWithID(inRequestID1, record.ReturnResult)
		// first detached (saga mode) outgoing with reason='first request'
		outRequestDetachedID1, outRequestDetached1 = outRequestWithID(inRequestID1, record.ReturnSaga)
		// first request's result
		inResultID1, inResult1 = resultWithID(inRequestID1)
	)

	t.Run("request_only", func(t *testing.T) {
		iterMock.cleanup()
		iterMock.pushRecord(inRequestID1, inRequest1)

		pendings, err := calc.PendingRequests(ctx, startPulse, objectID)
		require.NoError(t, err, "PendingRequests call without error")
		require.Equal(t, []insolar.ID{inRequestID1}, pendings, "single unclosed request")
	})

	t.Run("request_with_detached", func(t *testing.T) {
		iterMock.cleanup()
		iterMock.pushRecord(inRequestID1, inRequest1)
		// out-request should be before (while iterates) it's reason
		iterMock.pushRecord(outRequestDetachedID1, outRequestDetached1)

		pendings, err := calc.PendingRequests(ctx, startPulse, objectID)
		require.NoError(t, err, "PendingRequests call without error")
		require.Equal(t, []insolar.ID{inRequestID1, outRequestDetachedID1}, pendings,
			"single unclosed request with detached outgoing are both in pendings in reverse order")
	})

	t.Run("request_with_detached_and_result", func(t *testing.T) {
		iterMock.cleanup()
		iterMock.pushRecord(inRequestID1, inRequest1)
		// out-request should be before (while iterates) it's reason
		iterMock.pushRecord(outRequestDetachedID1, outRequestDetached1)
		// result
		iterMock.pushRecord(inResultID1, inResult1)

		pendings, err := calc.PendingRequests(ctx, startPulse, objectID)
		require.NoError(t, err, "PendingRequests call without error")
		require.Equal(t, []insolar.ID{}, pendings,
			"single closed request with detached outgoing are both not in pendings")
	})

	t.Run("request_with_result", func(t *testing.T) {
		iterMock.cleanup()
		iterMock.pushRecord(inRequestID1, inRequest1)
		iterMock.pushRecord(inResultID1, inResult1)

		pendings, err := calc.PendingRequests(ctx, startPulse, objectID)
		require.NoError(t, err, "PendingRequests call without error")
		require.Equal(t, 0, len(pendings), "check pendings")
	})

	t.Run("request_with_result_and_not_detached_outgoing", func(t *testing.T) {
		iterMock.cleanup()
		iterMock.pushRecord(inRequestID1, inRequest1)
		iterMock.pushRecord(outRequestID1, outRequest1)
		iterMock.pushRecord(inResultID1, inResult1)

		pendings, err := calc.PendingRequests(ctx, startPulse, objectID)
		require.NoError(t, err, "PendingRequests call without error")
		require.Equal(t, 0, len(pendings), "check pendings")
	})

	mt.Finish()
}

func fetchIteratorMockMaker(fim *fetchIteratorMock) fetchIteratorMaker {
	return func(
		ctx context.Context,
		cache *filamentCache,
		objectID, from insolar.ID,
		readUntil insolar.PulseNumber,
		fetcher jet.Fetcher,
		coordinator jet.Coordinator,
		sender bus.Sender,
	) fetchIterator {
		return fim
	}
}

type fetchIteratorMock struct {
	state   int
	records []record.CompositeFilamentRecord
}

func newFetchIteratorMock() *fetchIteratorMock {
	return &fetchIteratorMock{
		state: -1,
	}
}

// pushRecord adds records to records tail
func (fi *fetchIteratorMock) pushRecord(id insolar.ID, rec record.Record) {
	var fr record.CompositeFilamentRecord
	fr.RecordID = id
	virtual := record.Wrap(rec)
	fr.Record.Virtual = &virtual
	fi.records = append(fi.records, fr)
	fi.state++
}

func (fi *fetchIteratorMock) reset() {
	fi.state = len(fi.records) - 1
}

func (fi *fetchIteratorMock) cleanup() {
	fi.records = nil
	fi.reset()
}

func (fi *fetchIteratorMock) PrevID() *insolar.ID {
	panic("not implemented")
}

// check if iterator reaches head
func (fi *fetchIteratorMock) HasPrev() bool {
	return fi.state != -1
}

// Prev iterates record from tail
func (fi *fetchIteratorMock) Prev(ctx context.Context) (record.CompositeFilamentRecord, error) {
	rec := fi.records[fi.state]
	fi.state--
	return rec, nil
}
