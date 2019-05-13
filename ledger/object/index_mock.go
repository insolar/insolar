package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Index" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

// IndexMock implements github.com/insolar/insolar/ledger/object.Index
type IndexMock struct {
	t minimock.Tester

	LifelineForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)
	LifelineForIDCounter    uint64
	LifelineForIDPreCounter uint64
	LifelineForIDMock       mIndexMockLifelineForID

	SetBucketFunc       func(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) (r error)
	SetBucketCounter    uint64
	SetBucketPreCounter uint64
	SetBucketMock       mIndexMockSetBucket

	SetLifelineFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)
	SetLifelineCounter    uint64
	SetLifelinePreCounter uint64
	SetLifelineMock       mIndexMockSetLifeline

	SetRequestFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)
	SetRequestCounter    uint64
	SetRequestPreCounter uint64
	SetRequestMock       mIndexMockSetRequest

	SetResultRecordFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)
	SetResultRecordCounter    uint64
	SetResultRecordPreCounter uint64
	SetResultRecordMock       mIndexMockSetResultRecord
}

// NewIndexMock returns a mock for github.com/insolar/insolar/ledger/object.Index
func NewIndexMock(t minimock.Tester) *IndexMock {
	m := &IndexMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LifelineForIDMock = mIndexMockLifelineForID{mock: m}
	m.SetBucketMock = mIndexMockSetBucket{mock: m}
	m.SetLifelineMock = mIndexMockSetLifeline{mock: m}
	m.SetRequestMock = mIndexMockSetRequest{mock: m}
	m.SetResultRecordMock = mIndexMockSetResultRecord{mock: m}

	return m
}

type mIndexMockLifelineForID struct {
	mock              *IndexMock
	mainExpectation   *IndexMockLifelineForIDExpectation
	expectationSeries []*IndexMockLifelineForIDExpectation
}

type IndexMockLifelineForIDExpectation struct {
	input  *IndexMockLifelineForIDInput
	result *IndexMockLifelineForIDResult
}

type IndexMockLifelineForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type IndexMockLifelineForIDResult struct {
	r  Lifeline
	r1 error
}

// Expect specifies that invocation of Index.LifelineForID is expected from 1 to Infinity times
func (m *mIndexMockLifelineForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mIndexMockLifelineForID {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockLifelineForIDExpectation{}
	}
	m.mainExpectation.input = &IndexMockLifelineForIDInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of Index.LifelineForID
func (m *mIndexMockLifelineForID) Return(r Lifeline, r1 error) *IndexMock {
	m.mock.LifelineForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockLifelineForIDExpectation{}
	}
	m.mainExpectation.result = &IndexMockLifelineForIDResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of Index.LifelineForID is expected once
func (m *mIndexMockLifelineForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *IndexMockLifelineForIDExpectation {
	m.mock.LifelineForIDFunc = nil
	m.mainExpectation = nil

	expectation := &IndexMockLifelineForIDExpectation{}
	expectation.input = &IndexMockLifelineForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexMockLifelineForIDExpectation) Return(r Lifeline, r1 error) {
	e.result = &IndexMockLifelineForIDResult{r, r1}
}

// Set uses given function f as a mock of Index.LifelineForID method
func (m *mIndexMockLifelineForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error)) *IndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LifelineForIDFunc = f
	return m.mock
}

// LifelineForID implements github.com/insolar/insolar/ledger/object.Index interface
func (m *IndexMock) LifelineForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.LifelineForIDPreCounter, 1)
	defer atomic.AddUint64(&m.LifelineForIDCounter, 1)

	if len(m.LifelineForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LifelineForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexMock.LifelineForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.LifelineForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexMockLifelineForIDInput{p, p1, p2}, "Index.LifelineForID got unexpected parameters")

		result := m.LifelineForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.LifelineForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDMock.mainExpectation != nil {

		input := m.LifelineForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexMockLifelineForIDInput{p, p1, p2}, "Index.LifelineForID got unexpected parameters")
		}

		result := m.LifelineForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.LifelineForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LifelineForIDFunc == nil {
		m.t.Fatalf("Unexpected call to IndexMock.LifelineForID. %v %v %v", p, p1, p2)
		return
	}

	return m.LifelineForIDFunc(p, p1, p2)
}

// LifelineForIDMinimockCounter returns a count of IndexMock.LifelineForIDFunc invocations
func (m *IndexMock) LifelineForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDCounter)
}

// LifelineForIDMinimockPreCounter returns the value of IndexMock.LifelineForID invocations
func (m *IndexMock) LifelineForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LifelineForIDPreCounter)
}

// LifelineForIDFinished returns true if mock invocations count is ok
func (m *IndexMock) LifelineForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LifelineForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LifelineForIDCounter) == uint64(len(m.LifelineForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LifelineForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LifelineForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LifelineForIDFunc != nil {
		return atomic.LoadUint64(&m.LifelineForIDCounter) > 0
	}

	return true
}

type mIndexMockSetBucket struct {
	mock              *IndexMock
	mainExpectation   *IndexMockSetBucketExpectation
	expectationSeries []*IndexMockSetBucketExpectation
}

type IndexMockSetBucketExpectation struct {
	input  *IndexMockSetBucketInput
	result *IndexMockSetBucketResult
}

type IndexMockSetBucketInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 IndexBucket
}

type IndexMockSetBucketResult struct {
	r error
}

// Expect specifies that invocation of Index.SetBucket is expected from 1 to Infinity times
func (m *mIndexMockSetBucket) Expect(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) *mIndexMockSetBucket {
	m.mock.SetBucketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetBucketExpectation{}
	}
	m.mainExpectation.input = &IndexMockSetBucketInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of Index.SetBucket
func (m *mIndexMockSetBucket) Return(r error) *IndexMock {
	m.mock.SetBucketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetBucketExpectation{}
	}
	m.mainExpectation.result = &IndexMockSetBucketResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of Index.SetBucket is expected once
func (m *mIndexMockSetBucket) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) *IndexMockSetBucketExpectation {
	m.mock.SetBucketFunc = nil
	m.mainExpectation = nil

	expectation := &IndexMockSetBucketExpectation{}
	expectation.input = &IndexMockSetBucketInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexMockSetBucketExpectation) Return(r error) {
	e.result = &IndexMockSetBucketResult{r}
}

// Set uses given function f as a mock of Index.SetBucket method
func (m *mIndexMockSetBucket) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) (r error)) *IndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetBucketFunc = f
	return m.mock
}

// SetBucket implements github.com/insolar/insolar/ledger/object.Index interface
func (m *IndexMock) SetBucket(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) (r error) {
	counter := atomic.AddUint64(&m.SetBucketPreCounter, 1)
	defer atomic.AddUint64(&m.SetBucketCounter, 1)

	if len(m.SetBucketMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetBucketMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexMock.SetBucket. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetBucketMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexMockSetBucketInput{p, p1, p2}, "Index.SetBucket got unexpected parameters")

		result := m.SetBucketMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetBucket")
			return
		}

		r = result.r

		return
	}

	if m.SetBucketMock.mainExpectation != nil {

		input := m.SetBucketMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexMockSetBucketInput{p, p1, p2}, "Index.SetBucket got unexpected parameters")
		}

		result := m.SetBucketMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetBucket")
		}

		r = result.r

		return
	}

	if m.SetBucketFunc == nil {
		m.t.Fatalf("Unexpected call to IndexMock.SetBucket. %v %v %v", p, p1, p2)
		return
	}

	return m.SetBucketFunc(p, p1, p2)
}

// SetBucketMinimockCounter returns a count of IndexMock.SetBucketFunc invocations
func (m *IndexMock) SetBucketMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetBucketCounter)
}

// SetBucketMinimockPreCounter returns the value of IndexMock.SetBucket invocations
func (m *IndexMock) SetBucketMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetBucketPreCounter)
}

// SetBucketFinished returns true if mock invocations count is ok
func (m *IndexMock) SetBucketFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetBucketMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetBucketCounter) == uint64(len(m.SetBucketMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetBucketMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetBucketCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetBucketFunc != nil {
		return atomic.LoadUint64(&m.SetBucketCounter) > 0
	}

	return true
}

type mIndexMockSetLifeline struct {
	mock              *IndexMock
	mainExpectation   *IndexMockSetLifelineExpectation
	expectationSeries []*IndexMockSetLifelineExpectation
}

type IndexMockSetLifelineExpectation struct {
	input  *IndexMockSetLifelineInput
	result *IndexMockSetLifelineResult
}

type IndexMockSetLifelineInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 Lifeline
}

type IndexMockSetLifelineResult struct {
	r error
}

// Expect specifies that invocation of Index.SetLifeline is expected from 1 to Infinity times
func (m *mIndexMockSetLifeline) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *mIndexMockSetLifeline {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetLifelineExpectation{}
	}
	m.mainExpectation.input = &IndexMockSetLifelineInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of Index.SetLifeline
func (m *mIndexMockSetLifeline) Return(r error) *IndexMock {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetLifelineExpectation{}
	}
	m.mainExpectation.result = &IndexMockSetLifelineResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of Index.SetLifeline is expected once
func (m *mIndexMockSetLifeline) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *IndexMockSetLifelineExpectation {
	m.mock.SetLifelineFunc = nil
	m.mainExpectation = nil

	expectation := &IndexMockSetLifelineExpectation{}
	expectation.input = &IndexMockSetLifelineInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexMockSetLifelineExpectation) Return(r error) {
	e.result = &IndexMockSetLifelineResult{r}
}

// Set uses given function f as a mock of Index.SetLifeline method
func (m *mIndexMockSetLifeline) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)) *IndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLifelineFunc = f
	return m.mock
}

// SetLifeline implements github.com/insolar/insolar/ledger/object.Index interface
func (m *IndexMock) SetLifeline(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetLifelinePreCounter, 1)
	defer atomic.AddUint64(&m.SetLifelineCounter, 1)

	if len(m.SetLifelineMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLifelineMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetLifelineMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexMockSetLifelineInput{p, p1, p2, p3}, "Index.SetLifeline got unexpected parameters")

		result := m.SetLifelineMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetLifeline")
			return
		}

		r = result.r

		return
	}

	if m.SetLifelineMock.mainExpectation != nil {

		input := m.SetLifelineMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexMockSetLifelineInput{p, p1, p2, p3}, "Index.SetLifeline got unexpected parameters")
		}

		result := m.SetLifelineMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetLifeline")
		}

		r = result.r

		return
	}

	if m.SetLifelineFunc == nil {
		m.t.Fatalf("Unexpected call to IndexMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetLifelineFunc(p, p1, p2, p3)
}

// SetLifelineMinimockCounter returns a count of IndexMock.SetLifelineFunc invocations
func (m *IndexMock) SetLifelineMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineCounter)
}

// SetLifelineMinimockPreCounter returns the value of IndexMock.SetLifeline invocations
func (m *IndexMock) SetLifelineMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelinePreCounter)
}

// SetLifelineFinished returns true if mock invocations count is ok
func (m *IndexMock) SetLifelineFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetLifelineMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetLifelineCounter) == uint64(len(m.SetLifelineMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetLifelineMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetLifelineCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetLifelineFunc != nil {
		return atomic.LoadUint64(&m.SetLifelineCounter) > 0
	}

	return true
}

type mIndexMockSetRequest struct {
	mock              *IndexMock
	mainExpectation   *IndexMockSetRequestExpectation
	expectationSeries []*IndexMockSetRequestExpectation
}

type IndexMockSetRequestExpectation struct {
	input  *IndexMockSetRequestInput
	result *IndexMockSetRequestResult
}

type IndexMockSetRequestInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
}

type IndexMockSetRequestResult struct {
	r error
}

// Expect specifies that invocation of Index.SetRequest is expected from 1 to Infinity times
func (m *mIndexMockSetRequest) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *mIndexMockSetRequest {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetRequestExpectation{}
	}
	m.mainExpectation.input = &IndexMockSetRequestInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of Index.SetRequest
func (m *mIndexMockSetRequest) Return(r error) *IndexMock {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetRequestExpectation{}
	}
	m.mainExpectation.result = &IndexMockSetRequestResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of Index.SetRequest is expected once
func (m *mIndexMockSetRequest) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *IndexMockSetRequestExpectation {
	m.mock.SetRequestFunc = nil
	m.mainExpectation = nil

	expectation := &IndexMockSetRequestExpectation{}
	expectation.input = &IndexMockSetRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexMockSetRequestExpectation) Return(r error) {
	e.result = &IndexMockSetRequestResult{r}
}

// Set uses given function f as a mock of Index.SetRequest method
func (m *mIndexMockSetRequest) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)) *IndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetRequestFunc = f
	return m.mock
}

// SetRequest implements github.com/insolar/insolar/ledger/object.Index interface
func (m *IndexMock) SetRequest(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SetRequestCounter, 1)

	if len(m.SetRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexMockSetRequestInput{p, p1, p2, p3}, "Index.SetRequest got unexpected parameters")

		result := m.SetRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetRequest")
			return
		}

		r = result.r

		return
	}

	if m.SetRequestMock.mainExpectation != nil {

		input := m.SetRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexMockSetRequestInput{p, p1, p2, p3}, "Index.SetRequest got unexpected parameters")
		}

		result := m.SetRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetRequest")
		}

		r = result.r

		return
	}

	if m.SetRequestFunc == nil {
		m.t.Fatalf("Unexpected call to IndexMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetRequestFunc(p, p1, p2, p3)
}

// SetRequestMinimockCounter returns a count of IndexMock.SetRequestFunc invocations
func (m *IndexMock) SetRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestCounter)
}

// SetRequestMinimockPreCounter returns the value of IndexMock.SetRequest invocations
func (m *IndexMock) SetRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestPreCounter)
}

// SetRequestFinished returns true if mock invocations count is ok
func (m *IndexMock) SetRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetRequestCounter) == uint64(len(m.SetRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetRequestFunc != nil {
		return atomic.LoadUint64(&m.SetRequestCounter) > 0
	}

	return true
}

type mIndexMockSetResultRecord struct {
	mock              *IndexMock
	mainExpectation   *IndexMockSetResultRecordExpectation
	expectationSeries []*IndexMockSetResultRecordExpectation
}

type IndexMockSetResultRecordExpectation struct {
	input  *IndexMockSetResultRecordInput
	result *IndexMockSetResultRecordResult
}

type IndexMockSetResultRecordInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
}

type IndexMockSetResultRecordResult struct {
	r error
}

// Expect specifies that invocation of Index.SetResultRecord is expected from 1 to Infinity times
func (m *mIndexMockSetResultRecord) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *mIndexMockSetResultRecord {
	m.mock.SetResultRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetResultRecordExpectation{}
	}
	m.mainExpectation.input = &IndexMockSetResultRecordInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of Index.SetResultRecord
func (m *mIndexMockSetResultRecord) Return(r error) *IndexMock {
	m.mock.SetResultRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexMockSetResultRecordExpectation{}
	}
	m.mainExpectation.result = &IndexMockSetResultRecordResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of Index.SetResultRecord is expected once
func (m *mIndexMockSetResultRecord) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *IndexMockSetResultRecordExpectation {
	m.mock.SetResultRecordFunc = nil
	m.mainExpectation = nil

	expectation := &IndexMockSetResultRecordExpectation{}
	expectation.input = &IndexMockSetResultRecordInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexMockSetResultRecordExpectation) Return(r error) {
	e.result = &IndexMockSetResultRecordResult{r}
}

// Set uses given function f as a mock of Index.SetResultRecord method
func (m *mIndexMockSetResultRecord) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)) *IndexMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetResultRecordFunc = f
	return m.mock
}

// SetResultRecord implements github.com/insolar/insolar/ledger/object.Index interface
func (m *IndexMock) SetResultRecord(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetResultRecordPreCounter, 1)
	defer atomic.AddUint64(&m.SetResultRecordCounter, 1)

	if len(m.SetResultRecordMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetResultRecordMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexMock.SetResultRecord. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetResultRecordMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexMockSetResultRecordInput{p, p1, p2, p3}, "Index.SetResultRecord got unexpected parameters")

		result := m.SetResultRecordMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetResultRecord")
			return
		}

		r = result.r

		return
	}

	if m.SetResultRecordMock.mainExpectation != nil {

		input := m.SetResultRecordMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexMockSetResultRecordInput{p, p1, p2, p3}, "Index.SetResultRecord got unexpected parameters")
		}

		result := m.SetResultRecordMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexMock.SetResultRecord")
		}

		r = result.r

		return
	}

	if m.SetResultRecordFunc == nil {
		m.t.Fatalf("Unexpected call to IndexMock.SetResultRecord. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetResultRecordFunc(p, p1, p2, p3)
}

// SetResultRecordMinimockCounter returns a count of IndexMock.SetResultRecordFunc invocations
func (m *IndexMock) SetResultRecordMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultRecordCounter)
}

// SetResultRecordMinimockPreCounter returns the value of IndexMock.SetResultRecord invocations
func (m *IndexMock) SetResultRecordMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultRecordPreCounter)
}

// SetResultRecordFinished returns true if mock invocations count is ok
func (m *IndexMock) SetResultRecordFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetResultRecordMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetResultRecordCounter) == uint64(len(m.SetResultRecordMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetResultRecordMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetResultRecordCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetResultRecordFunc != nil {
		return atomic.LoadUint64(&m.SetResultRecordCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexMock) ValidateCallCounters() {

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to IndexMock.LifelineForID")
	}

	if !m.SetBucketFinished() {
		m.t.Fatal("Expected call to IndexMock.SetBucket")
	}

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to IndexMock.SetLifeline")
	}

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to IndexMock.SetRequest")
	}

	if !m.SetResultRecordFinished() {
		m.t.Fatal("Expected call to IndexMock.SetResultRecord")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexMock) MinimockFinish() {

	if !m.LifelineForIDFinished() {
		m.t.Fatal("Expected call to IndexMock.LifelineForID")
	}

	if !m.SetBucketFinished() {
		m.t.Fatal("Expected call to IndexMock.SetBucket")
	}

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to IndexMock.SetLifeline")
	}

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to IndexMock.SetRequest")
	}

	if !m.SetResultRecordFinished() {
		m.t.Fatal("Expected call to IndexMock.SetResultRecord")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *IndexMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.LifelineForIDFinished()
		ok = ok && m.SetBucketFinished()
		ok = ok && m.SetLifelineFinished()
		ok = ok && m.SetRequestFinished()
		ok = ok && m.SetResultRecordFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.LifelineForIDFinished() {
				m.t.Error("Expected call to IndexMock.LifelineForID")
			}

			if !m.SetBucketFinished() {
				m.t.Error("Expected call to IndexMock.SetBucket")
			}

			if !m.SetLifelineFinished() {
				m.t.Error("Expected call to IndexMock.SetLifeline")
			}

			if !m.SetRequestFinished() {
				m.t.Error("Expected call to IndexMock.SetRequest")
			}

			if !m.SetResultRecordFinished() {
				m.t.Error("Expected call to IndexMock.SetResultRecord")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *IndexMock) AllMocksCalled() bool {

	if !m.LifelineForIDFinished() {
		return false
	}

	if !m.SetBucketFinished() {
		return false
	}

	if !m.SetLifelineFinished() {
		return false
	}

	if !m.SetRequestFinished() {
		return false
	}

	if !m.SetResultRecordFinished() {
		return false
	}

	return true
}
