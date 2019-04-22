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

	ForIDFunc       func(p context.Context, p1 insolar.ID) (r Lifeline, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mIndexStorageMockForID

	SetFunc       func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mIndexStorageMockSet
}

//NewIndexStorageMock returns a mock for github.com/insolar/insolar/ledger/object.IndexStorage
func NewIndexStorageMock(t minimock.Tester) *IndexStorageMock {
	m := &IndexStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mIndexStorageMockForID{mock: m}
	m.SetMock = mIndexStorageMockSet{mock: m}

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
	p1 insolar.ID
}

type IndexStorageMockForIDResult struct {
	r  Lifeline
	r1 error
}

//Expect specifies that invocation of IndexStorage.ForID is expected from 1 to Infinity times
func (m *mIndexStorageMockForID) Expect(p context.Context, p1 insolar.ID) *mIndexStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &IndexStorageMockForIDInput{p, p1}
	return m
}

//Return specifies results of invocation of IndexStorage.ForID
func (m *mIndexStorageMockForID) Return(r Lifeline, r1 error) *IndexStorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &IndexStorageMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexStorage.ForID is expected once
func (m *mIndexStorageMockForID) ExpectOnce(p context.Context, p1 insolar.ID) *IndexStorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &IndexStorageMockForIDExpectation{}
	expectation.input = &IndexStorageMockForIDInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexStorageMockForIDExpectation) Return(r Lifeline, r1 error) {
	e.result = &IndexStorageMockForIDResult{r, r1}
}

//Set uses given function f as a mock of IndexStorage.ForID method
func (m *mIndexStorageMockForID) Set(f func(p context.Context, p1 insolar.ID) (r Lifeline, r1 error)) *IndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.IndexStorage interface
func (m *IndexStorageMock) ForID(p context.Context, p1 insolar.ID) (r Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexStorageMock.ForID. %v %v", p, p1)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexStorageMockForIDInput{p, p1}, "IndexStorage.ForID got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, IndexStorageMockForIDInput{p, p1}, "IndexStorage.ForID got unexpected parameters")
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
		m.t.Fatalf("Unexpected call to IndexStorageMock.ForID. %v %v", p, p1)
		return
	}

	return m.ForIDFunc(p, p1)
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

type mIndexStorageMockSet struct {
	mock              *IndexStorageMock
	mainExpectation   *IndexStorageMockSetExpectation
	expectationSeries []*IndexStorageMockSetExpectation
}

type IndexStorageMockSetExpectation struct {
	input  *IndexStorageMockSetInput
	result *IndexStorageMockSetResult
}

type IndexStorageMockSetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 Lifeline
}

type IndexStorageMockSetResult struct {
	r error
}

//Expect specifies that invocation of IndexStorage.Set is expected from 1 to Infinity times
func (m *mIndexStorageMockSet) Expect(p context.Context, p1 insolar.ID, p2 Lifeline) *mIndexStorageMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockSetExpectation{}
	}
	m.mainExpectation.input = &IndexStorageMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of IndexStorage.Set
func (m *mIndexStorageMockSet) Return(r error) *IndexStorageMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexStorageMockSetExpectation{}
	}
	m.mainExpectation.result = &IndexStorageMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of IndexStorage.Set is expected once
func (m *mIndexStorageMockSet) ExpectOnce(p context.Context, p1 insolar.ID, p2 Lifeline) *IndexStorageMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &IndexStorageMockSetExpectation{}
	expectation.input = &IndexStorageMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IndexStorageMockSetExpectation) Return(r error) {
	e.result = &IndexStorageMockSetResult{r}
}

//Set uses given function f as a mock of IndexStorage.Set method
func (m *mIndexStorageMockSet) Set(f func(p context.Context, p1 insolar.ID, p2 Lifeline) (r error)) *IndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/object.IndexStorage interface
func (m *IndexStorageMock) Set(p context.Context, p1 insolar.ID, p2 Lifeline) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexStorageMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexStorageMockSetInput{p, p1, p2}, "IndexStorage.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexStorageMockSetInput{p, p1, p2}, "IndexStorage.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IndexStorageMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to IndexStorageMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of IndexStorageMock.SetFunc invocations
func (m *IndexStorageMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of IndexStorageMock.Set invocations
func (m *IndexStorageMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *IndexStorageMock) SetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCounter) == uint64(len(m.SetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFunc != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexStorageMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.Set")
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

	if !m.SetFinished() {
		m.t.Fatal("Expected call to IndexStorageMock.Set")
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
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to IndexStorageMock.ForID")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to IndexStorageMock.Set")
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

	if !m.SetFinished() {
		return false
	}

	return true
}
