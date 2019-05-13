package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexModifier" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexModifierMock implements github.com/insolar/insolar/ledger/object.IndexModifier
type IndexModifierMock struct {
	t minimock.Tester

	SetBucketFunc       func(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) (r error)
	SetBucketCounter    uint64
	SetBucketPreCounter uint64
	SetBucketMock       mIndexModifierMockSetBucket

	SetLifelineFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)
	SetLifelineCounter    uint64
	SetLifelinePreCounter uint64
	SetLifelineMock       mIndexModifierMockSetLifeline

	SetRequestFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)
	SetRequestCounter    uint64
	SetRequestPreCounter uint64
	SetRequestMock       mIndexModifierMockSetRequest

	SetResultRecordFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)
	SetResultRecordCounter    uint64
	SetResultRecordPreCounter uint64
	SetResultRecordMock       mIndexModifierMockSetResultRecord
}

//NewIndexModifierMock returns a mock for github.com/insolar/insolar/ledger/object.IndexModifier
func NewIndexModifierMock(t minimock.Tester) *IndexModifierMock {
	m := &IndexModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.SetBucketMock = mIndexModifierMockSetBucket{mock: m}
	m.SetLifelineMock = mIndexModifierMockSetLifeline{mock: m}
	m.SetRequestMock = mIndexModifierMockSetRequest{mock: m}
	m.SetResultRecordMock = mIndexModifierMockSetResultRecord{mock: m}

	return m
}

type mIndexModifierMockSetBucket struct {
	mock              *IndexModifierMock
	mainExpectation   *IndexModifierMockSetBucketExpectation
	expectationSeries []*IndexModifierMockSetBucketExpectation
}

type IndexModifierMockSetBucketExpectation struct {
	input  *IndexModifierMockSetBucketInput
	result *IndexModifierMockSetBucketResult
}

type IndexModifierMockSetBucketInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 IndexBucket
}

type IndexModifierMockSetBucketResult struct {
	r error
}

//Expect specifies that invocation of IndexModifier.SetBucket is expected from 1 to Infinity times
func (m *mIndexModifierMockSetBucket) Expect(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) *mIndexModifierMockSetBucket {
	m.mock.SetBucketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetBucketExpectation{}
	}
	m.mainExpectation.input = &IndexModifierMockSetBucketInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexModifier.SetBucket
func (m *mIndexModifierMockSetBucket) Return(r error) *IndexModifierMock {
	m.mock.SetBucketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetBucketExpectation{}
	}
	m.mainExpectation.result = &IndexModifierMockSetBucketResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexModifier.SetBucket is expected once
func (m *mIndexModifierMockSetBucket) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) *IndexModifierMockSetBucketExpectation {
	m.mock.SetBucketFunc = nil
	m.mainExpectation = nil

	expectation := &IndexModifierMockSetBucketExpectation{}
	expectation.input = &IndexModifierMockSetBucketInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexModifierMockSetBucketExpectation) Return(r error) {
	e.result = &IndexModifierMockSetBucketResult{r}
}

//Set uses given function f as a mock of IndexModifier.SetBucket method
func (m *mIndexModifierMockSetBucket) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) (r error)) *IndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetBucketFunc = f
	return m.mock
}

//SetBucket implements github.com/insolar/insolar/ledger/object.IndexModifier interface
func (m *IndexModifierMock) SetBucket(p context.Context, p1 insolar.PulseNumber, p2 IndexBucket) (r error) {
	counter := atomic.AddUint64(&m.SetBucketPreCounter, 1)
	defer atomic.AddUint64(&m.SetBucketCounter, 1)

	if len(m.SetBucketMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetBucketMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexModifierMock.SetBucket. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetBucketMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexModifierMockSetBucketInput{p, p1, p2}, "IndexModifier.SetBucket got unexpected parameters")

		result := m.SetBucketMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetBucket")
			return
		}

		r = result.r

		return
	}

	if m.SetBucketMock.mainExpectation != nil {

		input := m.SetBucketMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexModifierMockSetBucketInput{p, p1, p2}, "IndexModifier.SetBucket got unexpected parameters")
		}

		result := m.SetBucketMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetBucket")
		}

		r = result.r

		return
	}

	if m.SetBucketFunc == nil {
		m.t.Fatalf("Unexpected call to IndexModifierMock.SetBucket. %v %v %v", p, p1, p2)
		return
	}

	return m.SetBucketFunc(p, p1, p2)
}

//SetBucketMinimockCounter returns a count of IndexModifierMock.SetBucketFunc invocations
func (m *IndexModifierMock) SetBucketMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetBucketCounter)
}

//SetBucketMinimockPreCounter returns the value of IndexModifierMock.SetBucket invocations
func (m *IndexModifierMock) SetBucketMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetBucketPreCounter)
}

//SetBucketFinished returns true if mock invocations count is ok
func (m *IndexModifierMock) SetBucketFinished() bool {
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

type mIndexModifierMockSetLifeline struct {
	mock              *IndexModifierMock
	mainExpectation   *IndexModifierMockSetLifelineExpectation
	expectationSeries []*IndexModifierMockSetLifelineExpectation
}

type IndexModifierMockSetLifelineExpectation struct {
	input  *IndexModifierMockSetLifelineInput
	result *IndexModifierMockSetLifelineResult
}

type IndexModifierMockSetLifelineInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 Lifeline
}

type IndexModifierMockSetLifelineResult struct {
	r error
}

//Expect specifies that invocation of IndexModifier.SetLifeline is expected from 1 to Infinity times
func (m *mIndexModifierMockSetLifeline) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *mIndexModifierMockSetLifeline {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetLifelineExpectation{}
	}
	m.mainExpectation.input = &IndexModifierMockSetLifelineInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of IndexModifier.SetLifeline
func (m *mIndexModifierMockSetLifeline) Return(r error) *IndexModifierMock {
	m.mock.SetLifelineFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetLifelineExpectation{}
	}
	m.mainExpectation.result = &IndexModifierMockSetLifelineResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexModifier.SetLifeline is expected once
func (m *mIndexModifierMockSetLifeline) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) *IndexModifierMockSetLifelineExpectation {
	m.mock.SetLifelineFunc = nil
	m.mainExpectation = nil

	expectation := &IndexModifierMockSetLifelineExpectation{}
	expectation.input = &IndexModifierMockSetLifelineInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexModifierMockSetLifelineExpectation) Return(r error) {
	e.result = &IndexModifierMockSetLifelineResult{r}
}

//Set uses given function f as a mock of IndexModifier.SetLifeline method
func (m *mIndexModifierMockSetLifeline) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error)) *IndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLifelineFunc = f
	return m.mock
}

//SetLifeline implements github.com/insolar/insolar/ledger/object.IndexModifier interface
func (m *IndexModifierMock) SetLifeline(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetLifelinePreCounter, 1)
	defer atomic.AddUint64(&m.SetLifelineCounter, 1)

	if len(m.SetLifelineMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLifelineMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexModifierMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetLifelineMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexModifierMockSetLifelineInput{p, p1, p2, p3}, "IndexModifier.SetLifeline got unexpected parameters")

		result := m.SetLifelineMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetLifeline")
			return
		}

		r = result.r

		return
	}

	if m.SetLifelineMock.mainExpectation != nil {

		input := m.SetLifelineMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexModifierMockSetLifelineInput{p, p1, p2, p3}, "IndexModifier.SetLifeline got unexpected parameters")
		}

		result := m.SetLifelineMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetLifeline")
		}

		r = result.r

		return
	}

	if m.SetLifelineFunc == nil {
		m.t.Fatalf("Unexpected call to IndexModifierMock.SetLifeline. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetLifelineFunc(p, p1, p2, p3)
}

//SetLifelineMinimockCounter returns a count of IndexModifierMock.SetLifelineFunc invocations
func (m *IndexModifierMock) SetLifelineMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelineCounter)
}

//SetLifelineMinimockPreCounter returns the value of IndexModifierMock.SetLifeline invocations
func (m *IndexModifierMock) SetLifelineMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLifelinePreCounter)
}

//SetLifelineFinished returns true if mock invocations count is ok
func (m *IndexModifierMock) SetLifelineFinished() bool {
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

type mIndexModifierMockSetRequest struct {
	mock              *IndexModifierMock
	mainExpectation   *IndexModifierMockSetRequestExpectation
	expectationSeries []*IndexModifierMockSetRequestExpectation
}

type IndexModifierMockSetRequestExpectation struct {
	input  *IndexModifierMockSetRequestInput
	result *IndexModifierMockSetRequestResult
}

type IndexModifierMockSetRequestInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
}

type IndexModifierMockSetRequestResult struct {
	r error
}

//Expect specifies that invocation of IndexModifier.SetRequest is expected from 1 to Infinity times
func (m *mIndexModifierMockSetRequest) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *mIndexModifierMockSetRequest {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.input = &IndexModifierMockSetRequestInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of IndexModifier.SetRequest
func (m *mIndexModifierMockSetRequest) Return(r error) *IndexModifierMock {
	m.mock.SetRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetRequestExpectation{}
	}
	m.mainExpectation.result = &IndexModifierMockSetRequestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexModifier.SetRequest is expected once
func (m *mIndexModifierMockSetRequest) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *IndexModifierMockSetRequestExpectation {
	m.mock.SetRequestFunc = nil
	m.mainExpectation = nil

	expectation := &IndexModifierMockSetRequestExpectation{}
	expectation.input = &IndexModifierMockSetRequestInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexModifierMockSetRequestExpectation) Return(r error) {
	e.result = &IndexModifierMockSetRequestResult{r}
}

//Set uses given function f as a mock of IndexModifier.SetRequest method
func (m *mIndexModifierMockSetRequest) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)) *IndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetRequestFunc = f
	return m.mock
}

//SetRequest implements github.com/insolar/insolar/ledger/object.IndexModifier interface
func (m *IndexModifierMock) SetRequest(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetRequestPreCounter, 1)
	defer atomic.AddUint64(&m.SetRequestCounter, 1)

	if len(m.SetRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexModifierMockSetRequestInput{p, p1, p2, p3}, "IndexModifier.SetRequest got unexpected parameters")

		result := m.SetRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetRequest")
			return
		}

		r = result.r

		return
	}

	if m.SetRequestMock.mainExpectation != nil {

		input := m.SetRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexModifierMockSetRequestInput{p, p1, p2, p3}, "IndexModifier.SetRequest got unexpected parameters")
		}

		result := m.SetRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetRequest")
		}

		r = result.r

		return
	}

	if m.SetRequestFunc == nil {
		m.t.Fatalf("Unexpected call to IndexModifierMock.SetRequest. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetRequestFunc(p, p1, p2, p3)
}

//SetRequestMinimockCounter returns a count of IndexModifierMock.SetRequestFunc invocations
func (m *IndexModifierMock) SetRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestCounter)
}

//SetRequestMinimockPreCounter returns the value of IndexModifierMock.SetRequest invocations
func (m *IndexModifierMock) SetRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRequestPreCounter)
}

//SetRequestFinished returns true if mock invocations count is ok
func (m *IndexModifierMock) SetRequestFinished() bool {
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

type mIndexModifierMockSetResultRecord struct {
	mock              *IndexModifierMock
	mainExpectation   *IndexModifierMockSetResultRecordExpectation
	expectationSeries []*IndexModifierMockSetResultRecordExpectation
}

type IndexModifierMockSetResultRecordExpectation struct {
	input  *IndexModifierMockSetResultRecordInput
	result *IndexModifierMockSetResultRecordResult
}

type IndexModifierMockSetResultRecordInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
	p3 insolar.ID
}

type IndexModifierMockSetResultRecordResult struct {
	r error
}

//Expect specifies that invocation of IndexModifier.SetResultRecord is expected from 1 to Infinity times
func (m *mIndexModifierMockSetResultRecord) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *mIndexModifierMockSetResultRecord {
	m.mock.SetResultRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetResultRecordExpectation{}
	}
	m.mainExpectation.input = &IndexModifierMockSetResultRecordInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of IndexModifier.SetResultRecord
func (m *mIndexModifierMockSetResultRecord) Return(r error) *IndexModifierMock {
	m.mock.SetResultRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexModifierMockSetResultRecordExpectation{}
	}
	m.mainExpectation.result = &IndexModifierMockSetResultRecordResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexModifier.SetResultRecord is expected once
func (m *mIndexModifierMockSetResultRecord) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) *IndexModifierMockSetResultRecordExpectation {
	m.mock.SetResultRecordFunc = nil
	m.mainExpectation = nil

	expectation := &IndexModifierMockSetResultRecordExpectation{}
	expectation.input = &IndexModifierMockSetResultRecordInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexModifierMockSetResultRecordExpectation) Return(r error) {
	e.result = &IndexModifierMockSetResultRecordResult{r}
}

//Set uses given function f as a mock of IndexModifier.SetResultRecord method
func (m *mIndexModifierMockSetResultRecord) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error)) *IndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetResultRecordFunc = f
	return m.mock
}

//SetResultRecord implements github.com/insolar/insolar/ledger/object.IndexModifier interface
func (m *IndexModifierMock) SetResultRecord(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.SetResultRecordPreCounter, 1)
	defer atomic.AddUint64(&m.SetResultRecordCounter, 1)

	if len(m.SetResultRecordMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetResultRecordMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexModifierMock.SetResultRecord. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetResultRecordMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexModifierMockSetResultRecordInput{p, p1, p2, p3}, "IndexModifier.SetResultRecord got unexpected parameters")

		result := m.SetResultRecordMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetResultRecord")
			return
		}

		r = result.r

		return
	}

	if m.SetResultRecordMock.mainExpectation != nil {

		input := m.SetResultRecordMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexModifierMockSetResultRecordInput{p, p1, p2, p3}, "IndexModifier.SetResultRecord got unexpected parameters")
		}

		result := m.SetResultRecordMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexModifierMock.SetResultRecord")
		}

		r = result.r

		return
	}

	if m.SetResultRecordFunc == nil {
		m.t.Fatalf("Unexpected call to IndexModifierMock.SetResultRecord. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetResultRecordFunc(p, p1, p2, p3)
}

//SetResultRecordMinimockCounter returns a count of IndexModifierMock.SetResultRecordFunc invocations
func (m *IndexModifierMock) SetResultRecordMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultRecordCounter)
}

//SetResultRecordMinimockPreCounter returns the value of IndexModifierMock.SetResultRecord invocations
func (m *IndexModifierMock) SetResultRecordMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetResultRecordPreCounter)
}

//SetResultRecordFinished returns true if mock invocations count is ok
func (m *IndexModifierMock) SetResultRecordFinished() bool {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexModifierMock) ValidateCallCounters() {

	if !m.SetBucketFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetBucket")
	}

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetLifeline")
	}

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetRequest")
	}

	if !m.SetResultRecordFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetResultRecord")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexModifierMock) MinimockFinish() {

	if !m.SetBucketFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetBucket")
	}

	if !m.SetLifelineFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetLifeline")
	}

	if !m.SetRequestFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetRequest")
	}

	if !m.SetResultRecordFinished() {
		m.t.Fatal("Expected call to IndexModifierMock.SetResultRecord")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.SetBucketFinished()
		ok = ok && m.SetLifelineFinished()
		ok = ok && m.SetRequestFinished()
		ok = ok && m.SetResultRecordFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.SetBucketFinished() {
				m.t.Error("Expected call to IndexModifierMock.SetBucket")
			}

			if !m.SetLifelineFinished() {
				m.t.Error("Expected call to IndexModifierMock.SetLifeline")
			}

			if !m.SetRequestFinished() {
				m.t.Error("Expected call to IndexModifierMock.SetRequest")
			}

			if !m.SetResultRecordFinished() {
				m.t.Error("Expected call to IndexModifierMock.SetResultRecord")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *IndexModifierMock) AllMocksCalled() bool {

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
