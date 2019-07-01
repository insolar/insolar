package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexStorage" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexStorageMock implements github.com/insolar/insolar/ledger/object.IndexStorage
type IndexStorageMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r FilamentIndex, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mIndexStorageMockForID

	ForPulseFunc       func(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mIndexStorageMockForPulse

	SetIndexFunc       func(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) (r error)
	SetIndexCounter    uint64
	SetIndexPreCounter uint64
	SetIndexMock       mIndexStorageMockSetIndex
}

//NewIndexStorageMock returns a mock for github.com/insolar/insolar/ledger/object.IndexStorage
func NewIndexStorageMock(t minimock.Tester) *IndexStorageMock {
	m := &IndexStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mIndexStorageMockForID{mock: m}
	m.ForPulseMock = mIndexStorageMockForPulse{mock: m}
	m.SetIndexMock = mIndexStorageMockSetIndex{mock: m}

	return m
}

type mIndexStorageMockForID struct {
	mock              *IndexStorageMock
	mainExpectation   *IndexStorageMockForIDExpectation
	expectationSeries []*IndexStorageMockForIDExpectation
}

type IndexStorageMockForIDExpectation struct {
	input  *IndexStorageMockForIDInput
	result *IndexStorageMockForIDResult
}

type IndexStorageMockForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type IndexStorageMockForIDResult struct {
	r  FilamentIndex
	r1 error
}

//Expect specifies that invocation of IndexStorage.ForID is expected from 1 to Infinity times
func (m *mIndexStorageMockForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mIndexStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &IndexStorageMockForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexStorage.ForID
func (m *mIndexStorageMockForID) Return(r FilamentIndex, r1 error) *IndexStorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &IndexStorageMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexStorage.ForID is expected once
func (m *mIndexStorageMockForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *IndexStorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &IndexStorageMockForIDExpectation{}
	expectation.input = &IndexStorageMockForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexStorageMockForIDExpectation) Return(r FilamentIndex, r1 error) {
	e.result = &IndexStorageMockForIDResult{r, r1}
}

//Set uses given function f as a mock of IndexStorage.ForID method
func (m *mIndexStorageMockForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r FilamentIndex, r1 error)) *IndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.IndexStorage interface
func (m *IndexStorageMock) ForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r FilamentIndex, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexStorageMock.ForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexStorageMockForIDInput{p, p1, p2}, "IndexStorage.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexStorageMockForIDInput{p, p1, p2}, "IndexStorage.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to IndexStorageMock.ForID. %v %v %v", p, p1, p2)
		return
	}

	return m.ForIDFunc(p, p1, p2)
}

//ForIDMinimockCounter returns a count of IndexStorageMock.ForIDFunc invocations
func (m *IndexStorageMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of IndexStorageMock.ForID invocations
func (m *IndexStorageMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *IndexStorageMock) ForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForIDCounter) == uint64(len(m.ForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForIDFunc != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	return true
}

type mIndexStorageMockForPulse struct {
	mock              *IndexStorageMock
	mainExpectation   *IndexStorageMockForPulseExpectation
	expectationSeries []*IndexStorageMockForPulseExpectation
}

type IndexStorageMockForPulseExpectation struct {
	input  *IndexStorageMockForPulseInput
	result *IndexStorageMockForPulseResult
}

type IndexStorageMockForPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type IndexStorageMockForPulseResult struct {
	r []FilamentIndex
}

//Expect specifies that invocation of IndexStorage.ForPulse is expected from 1 to Infinity times
func (m *mIndexStorageMockForPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mIndexStorageMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockForPulseExpectation{}
	}
	m.mainExpectation.input = &IndexStorageMockForPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of IndexStorage.ForPulse
func (m *mIndexStorageMockForPulse) Return(r []FilamentIndex) *IndexStorageMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockForPulseExpectation{}
	}
	m.mainExpectation.result = &IndexStorageMockForPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexStorage.ForPulse is expected once
func (m *mIndexStorageMockForPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *IndexStorageMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &IndexStorageMockForPulseExpectation{}
	expectation.input = &IndexStorageMockForPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexStorageMockForPulseExpectation) Return(r []FilamentIndex) {
	e.result = &IndexStorageMockForPulseResult{r}
}

//Set uses given function f as a mock of IndexStorage.ForPulse method
func (m *mIndexStorageMockForPulse) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex)) *IndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/object.IndexStorage interface
func (m *IndexStorageMock) ForPulse(p context.Context, p1 insolar.PulseNumber) (r []FilamentIndex) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexStorageMock.ForPulse. %v %v", p, p1)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexStorageMockForPulseInput{p, p1}, "IndexStorage.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.ForPulse")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexStorageMockForPulseInput{p, p1}, "IndexStorage.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.ForPulse")
		}

		r = result.r

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to IndexStorageMock.ForPulse. %v %v", p, p1)
		return
	}

	return m.ForPulseFunc(p, p1)
}

//ForPulseMinimockCounter returns a count of IndexStorageMock.ForPulseFunc invocations
func (m *IndexStorageMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

//ForPulseMinimockPreCounter returns the value of IndexStorageMock.ForPulse invocations
func (m *IndexStorageMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

//ForPulseFinished returns true if mock invocations count is ok
func (m *IndexStorageMock) ForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseCounter) == uint64(len(m.ForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseFunc != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	return true
}

type mIndexStorageMockSetIndex struct {
	mock              *IndexStorageMock
	mainExpectation   *IndexStorageMockSetIndexExpectation
	expectationSeries []*IndexStorageMockSetIndexExpectation
}

type IndexStorageMockSetIndexExpectation struct {
	input  *IndexStorageMockSetIndexInput
	result *IndexStorageMockSetIndexResult
}

type IndexStorageMockSetIndexInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 FilamentIndex
}

type IndexStorageMockSetIndexResult struct {
	r error
}

//Expect specifies that invocation of IndexStorage.SetIndex is expected from 1 to Infinity times
func (m *mIndexStorageMockSetIndex) Expect(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) *mIndexStorageMockSetIndex {
	m.mock.SetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockSetIndexExpectation{}
	}
	m.mainExpectation.input = &IndexStorageMockSetIndexInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexStorage.SetIndex
func (m *mIndexStorageMockSetIndex) Return(r error) *IndexStorageMock {
	m.mock.SetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockSetIndexExpectation{}
	}
	m.mainExpectation.result = &IndexStorageMockSetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexStorage.SetIndex is expected once
func (m *mIndexStorageMockSetIndex) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) *IndexStorageMockSetIndexExpectation {
	m.mock.SetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &IndexStorageMockSetIndexExpectation{}
	expectation.input = &IndexStorageMockSetIndexInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexStorageMockSetIndexExpectation) Return(r error) {
	e.result = &IndexStorageMockSetIndexResult{r}
}

//Set uses given function f as a mock of IndexStorage.SetIndex method
func (m *mIndexStorageMockSetIndex) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) (r error)) *IndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetIndexFunc = f
	return m.mock
}

//SetIndex implements github.com/insolar/insolar/ledger/object.IndexStorage interface
func (m *IndexStorageMock) SetIndex(p context.Context, p1 insolar.PulseNumber, p2 FilamentIndex) (r error) {
	counter := atomic.AddUint64(&m.SetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.SetIndexCounter, 1)

	if len(m.SetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexStorageMock.SetIndex. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexStorageMockSetIndexInput{p, p1, p2}, "IndexStorage.SetIndex got unexpected parameters")

		result := m.SetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.SetIndex")
			return
		}

		r = result.r

		return
	}

	if m.SetIndexMock.mainExpectation != nil {

		input := m.SetIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexStorageMockSetIndexInput{p, p1, p2}, "IndexStorage.SetIndex got unexpected parameters")
		}

		result := m.SetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.SetIndex")
		}

		r = result.r

		return
	}

	if m.SetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to IndexStorageMock.SetIndex. %v %v %v", p, p1, p2)
		return
	}

	return m.SetIndexFunc(p, p1, p2)
}

//SetIndexMinimockCounter returns a count of IndexStorageMock.SetIndexFunc invocations
func (m *IndexStorageMock) SetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetIndexCounter)
}

//SetIndexMinimockPreCounter returns the value of IndexStorageMock.SetIndex invocations
func (m *IndexStorageMock) SetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetIndexPreCounter)
}

//SetIndexFinished returns true if mock invocations count is ok
func (m *IndexStorageMock) SetIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetIndexCounter) == uint64(len(m.SetIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetIndexFunc != nil {
		return atomic.LoadUint64(&m.SetIndexCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexStorageMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.ForID")
	}

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.ForPulse")
	}

	if !m.SetIndexFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.SetIndex")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexStorageMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.ForID")
	}

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.ForPulse")
	}

	if !m.SetIndexFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.SetIndex")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()
		ok = ok && m.ForPulseFinished()
		ok = ok && m.SetIndexFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to IndexStorageMock.ForID")
			}

			if !m.ForPulseFinished() {
				m.t.Error("Expected call to IndexStorageMock.ForPulse")
			}

			if !m.SetIndexFinished() {
				m.t.Error("Expected call to IndexStorageMock.SetIndex")
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
func (m *IndexStorageMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.ForPulseFinished() {
		return false
	}

	if !m.SetIndexFinished() {
		return false
	}

	return true
}
