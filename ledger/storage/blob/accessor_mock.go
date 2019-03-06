package blob

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Accessor" can be found in github.com/insolar/insolar/ledger/storage/blob
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//AccessorMock implements github.com/insolar/insolar/ledger/storage/blob.Accessor
type AccessorMock struct {
	t minimock.Tester

	GetFunc       func(p context.Context, p1 core.RecordID) (r Blob, r1 error)
	GetCounter    uint64
	GetPreCounter uint64
	GetMock       mAccessorMockGet
}

//NewAccessorMock returns a mock for github.com/insolar/insolar/ledger/storage/blob.Accessor
func NewAccessorMock(t minimock.Tester) *AccessorMock {
	m := &AccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetMock = mAccessorMockGet{mock: m}

	return m
}

type mAccessorMockGet struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockGetExpectation
	expectationSeries []*AccessorMockGetExpectation
}

type AccessorMockGetExpectation struct {
	input  *AccessorMockGetInput
	result *AccessorMockGetResult
}

type AccessorMockGetInput struct {
	p  context.Context
	p1 core.RecordID
}

type AccessorMockGetResult struct {
	r  Blob
	r1 error
}

//Expect specifies that invocation of Accessor.Get is expected from 1 to Infinity times
func (m *mAccessorMockGet) Expect(p context.Context, p1 core.RecordID) *mAccessorMockGet {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetExpectation{}
	}
	m.mainExpectation.input = &AccessorMockGetInput{p, p1}
	return m
}

//Return specifies results of invocation of Accessor.Get
func (m *mAccessorMockGet) Return(r Blob, r1 error) *AccessorMock {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetExpectation{}
	}
	m.mainExpectation.result = &AccessorMockGetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.Get is expected once
func (m *mAccessorMockGet) ExpectOnce(p context.Context, p1 core.RecordID) *AccessorMockGetExpectation {
	m.mock.GetFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockGetExpectation{}
	expectation.input = &AccessorMockGetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockGetExpectation) Return(r Blob, r1 error) {
	e.result = &AccessorMockGetResult{r, r1}
}

//Set uses given function f as a mock of Accessor.Get method
func (m *mAccessorMockGet) Set(f func(p context.Context, p1 core.RecordID) (r Blob, r1 error)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetFunc = f
	return m.mock
}

//Get implements github.com/insolar/insolar/ledger/storage/blob.Accessor interface
func (m *AccessorMock) Get(p context.Context, p1 core.RecordID) (r Blob, r1 error) {
	counter := atomic.AddUint64(&m.GetPreCounter, 1)
	defer atomic.AddUint64(&m.GetCounter, 1)

	if len(m.GetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.Get. %v %v", p, p1)
			return
		}

		input := m.GetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockGetInput{p, p1}, "Accessor.Get got unexpected parameters")

		result := m.GetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.Get")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMock.mainExpectation != nil {

		input := m.GetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockGetInput{p, p1}, "Accessor.Get got unexpected parameters")
		}

		result := m.GetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.Get")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.Get. %v %v", p, p1)
		return
	}

	return m.GetFunc(p, p1)
}

//GetMinimockCounter returns a count of AccessorMock.GetFunc invocations
func (m *AccessorMock) GetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCounter)
}

//GetMinimockPreCounter returns the value of AccessorMock.Get invocations
func (m *AccessorMock) GetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreCounter)
}

//GetFinished returns true if mock invocations count is ok
func (m *AccessorMock) GetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCounter) == uint64(len(m.GetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetFunc != nil {
		return atomic.LoadUint64(&m.GetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) ValidateCallCounters() {

	if !m.GetFinished() {
		m.t.Fatal("Expected call to AccessorMock.Get")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AccessorMock) MinimockFinish() {

	if !m.GetFinished() {
		m.t.Fatal("Expected call to AccessorMock.Get")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetFinished() {
				m.t.Error("Expected call to AccessorMock.Get")
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
func (m *AccessorMock) AllMocksCalled() bool {

	if !m.GetFinished() {
		return false
	}

	return true
}
