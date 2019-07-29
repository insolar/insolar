package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "AtomicRecordStorage" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//AtomicRecordStorageMock implements github.com/insolar/insolar/ledger/object.AtomicRecordStorage
type AtomicRecordStorageMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.ID) (r record.Material, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mAtomicRecordStorageMockForID

	SetAtomicFunc       func(p context.Context, p1 ...record.Material) (r error)
	SetAtomicCounter    uint64
	SetAtomicPreCounter uint64
	SetAtomicMock       mAtomicRecordStorageMockSetAtomic
}

//NewAtomicRecordStorageMock returns a mock for github.com/insolar/insolar/ledger/object.AtomicRecordStorage
func NewAtomicRecordStorageMock(t minimock.Tester) *AtomicRecordStorageMock {
	m := &AtomicRecordStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mAtomicRecordStorageMockForID{mock: m}
	m.SetAtomicMock = mAtomicRecordStorageMockSetAtomic{mock: m}

	return m
}

type mAtomicRecordStorageMockForID struct {
	mock              *AtomicRecordStorageMock
	mainExpectation   *AtomicRecordStorageMockForIDExpectation
	expectationSeries []*AtomicRecordStorageMockForIDExpectation
}

type AtomicRecordStorageMockForIDExpectation struct {
	input  *AtomicRecordStorageMockForIDInput
	result *AtomicRecordStorageMockForIDResult
}

type AtomicRecordStorageMockForIDInput struct {
	p  context.Context
	p1 insolar.ID
}

type AtomicRecordStorageMockForIDResult struct {
	r  record.Material
	r1 error
}

//Expect specifies that invocation of AtomicRecordStorage.ForID is expected from 1 to Infinity times
func (m *mAtomicRecordStorageMockForID) Expect(p context.Context, p1 insolar.ID) *mAtomicRecordStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AtomicRecordStorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &AtomicRecordStorageMockForIDInput{p, p1}
	return m
}

//Return specifies results of invocation of AtomicRecordStorage.ForID
func (m *mAtomicRecordStorageMockForID) Return(r record.Material, r1 error) *AtomicRecordStorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AtomicRecordStorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &AtomicRecordStorageMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of AtomicRecordStorage.ForID is expected once
func (m *mAtomicRecordStorageMockForID) ExpectOnce(p context.Context, p1 insolar.ID) *AtomicRecordStorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &AtomicRecordStorageMockForIDExpectation{}
	expectation.input = &AtomicRecordStorageMockForIDInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AtomicRecordStorageMockForIDExpectation) Return(r record.Material, r1 error) {
	e.result = &AtomicRecordStorageMockForIDResult{r, r1}
}

//Set uses given function f as a mock of AtomicRecordStorage.ForID method
func (m *mAtomicRecordStorageMockForID) Set(f func(p context.Context, p1 insolar.ID) (r record.Material, r1 error)) *AtomicRecordStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.AtomicRecordStorage interface
func (m *AtomicRecordStorageMock) ForID(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AtomicRecordStorageMock.ForID. %v %v", p, p1)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AtomicRecordStorageMockForIDInput{p, p1}, "AtomicRecordStorage.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AtomicRecordStorageMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AtomicRecordStorageMockForIDInput{p, p1}, "AtomicRecordStorage.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AtomicRecordStorageMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to AtomicRecordStorageMock.ForID. %v %v", p, p1)
		return
	}

	return m.ForIDFunc(p, p1)
}

//ForIDMinimockCounter returns a count of AtomicRecordStorageMock.ForIDFunc invocations
func (m *AtomicRecordStorageMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of AtomicRecordStorageMock.ForID invocations
func (m *AtomicRecordStorageMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *AtomicRecordStorageMock) ForIDFinished() bool {
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

type mAtomicRecordStorageMockSetAtomic struct {
	mock              *AtomicRecordStorageMock
	mainExpectation   *AtomicRecordStorageMockSetAtomicExpectation
	expectationSeries []*AtomicRecordStorageMockSetAtomicExpectation
}

type AtomicRecordStorageMockSetAtomicExpectation struct {
	input  *AtomicRecordStorageMockSetAtomicInput
	result *AtomicRecordStorageMockSetAtomicResult
}

type AtomicRecordStorageMockSetAtomicInput struct {
	p  context.Context
	p1 []record.Material
}

type AtomicRecordStorageMockSetAtomicResult struct {
	r error
}

//Expect specifies that invocation of AtomicRecordStorage.SetAtomic is expected from 1 to Infinity times
func (m *mAtomicRecordStorageMockSetAtomic) Expect(p context.Context, p1 ...record.Material) *mAtomicRecordStorageMockSetAtomic {
	m.mock.SetAtomicFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AtomicRecordStorageMockSetAtomicExpectation{}
	}
	m.mainExpectation.input = &AtomicRecordStorageMockSetAtomicInput{p, p1}
	return m
}

//Return specifies results of invocation of AtomicRecordStorage.SetAtomic
func (m *mAtomicRecordStorageMockSetAtomic) Return(r error) *AtomicRecordStorageMock {
	m.mock.SetAtomicFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AtomicRecordStorageMockSetAtomicExpectation{}
	}
	m.mainExpectation.result = &AtomicRecordStorageMockSetAtomicResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of AtomicRecordStorage.SetAtomic is expected once
func (m *mAtomicRecordStorageMockSetAtomic) ExpectOnce(p context.Context, p1 ...record.Material) *AtomicRecordStorageMockSetAtomicExpectation {
	m.mock.SetAtomicFunc = nil
	m.mainExpectation = nil

	expectation := &AtomicRecordStorageMockSetAtomicExpectation{}
	expectation.input = &AtomicRecordStorageMockSetAtomicInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AtomicRecordStorageMockSetAtomicExpectation) Return(r error) {
	e.result = &AtomicRecordStorageMockSetAtomicResult{r}
}

//Set uses given function f as a mock of AtomicRecordStorage.SetAtomic method
func (m *mAtomicRecordStorageMockSetAtomic) Set(f func(p context.Context, p1 ...record.Material) (r error)) *AtomicRecordStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetAtomicFunc = f
	return m.mock
}

//SetAtomic implements github.com/insolar/insolar/ledger/object.AtomicRecordStorage interface
func (m *AtomicRecordStorageMock) SetAtomic(p context.Context, p1 ...record.Material) (r error) {
	counter := atomic.AddUint64(&m.SetAtomicPreCounter, 1)
	defer atomic.AddUint64(&m.SetAtomicCounter, 1)

	if len(m.SetAtomicMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetAtomicMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AtomicRecordStorageMock.SetAtomic. %v %v", p, p1)
			return
		}

		input := m.SetAtomicMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AtomicRecordStorageMockSetAtomicInput{p, p1}, "AtomicRecordStorage.SetAtomic got unexpected parameters")

		result := m.SetAtomicMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AtomicRecordStorageMock.SetAtomic")
			return
		}

		r = result.r

		return
	}

	if m.SetAtomicMock.mainExpectation != nil {

		input := m.SetAtomicMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AtomicRecordStorageMockSetAtomicInput{p, p1}, "AtomicRecordStorage.SetAtomic got unexpected parameters")
		}

		result := m.SetAtomicMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AtomicRecordStorageMock.SetAtomic")
		}

		r = result.r

		return
	}

	if m.SetAtomicFunc == nil {
		m.t.Fatalf("Unexpected call to AtomicRecordStorageMock.SetAtomic. %v %v", p, p1)
		return
	}

	return m.SetAtomicFunc(p, p1...)
}

//SetAtomicMinimockCounter returns a count of AtomicRecordStorageMock.SetAtomicFunc invocations
func (m *AtomicRecordStorageMock) SetAtomicMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetAtomicCounter)
}

//SetAtomicMinimockPreCounter returns the value of AtomicRecordStorageMock.SetAtomic invocations
func (m *AtomicRecordStorageMock) SetAtomicMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetAtomicPreCounter)
}

//SetAtomicFinished returns true if mock invocations count is ok
func (m *AtomicRecordStorageMock) SetAtomicFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetAtomicMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetAtomicCounter) == uint64(len(m.SetAtomicMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetAtomicMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetAtomicCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetAtomicFunc != nil {
		return atomic.LoadUint64(&m.SetAtomicCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AtomicRecordStorageMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to AtomicRecordStorageMock.ForID")
	}

	if !m.SetAtomicFinished() {
		m.t.Fatal("Expected call to AtomicRecordStorageMock.SetAtomic")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AtomicRecordStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AtomicRecordStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AtomicRecordStorageMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to AtomicRecordStorageMock.ForID")
	}

	if !m.SetAtomicFinished() {
		m.t.Fatal("Expected call to AtomicRecordStorageMock.SetAtomic")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AtomicRecordStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AtomicRecordStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()
		ok = ok && m.SetAtomicFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to AtomicRecordStorageMock.ForID")
			}

			if !m.SetAtomicFinished() {
				m.t.Error("Expected call to AtomicRecordStorageMock.SetAtomic")
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
func (m *AtomicRecordStorageMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.SetAtomicFinished() {
		return false
	}

	return true
}
