package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "OriginProvider" can be found in github.com/insolar/insolar/insolar
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
)

//OriginProviderMock implements github.com/insolar/insolar/insolar.OriginProvider
type OriginProviderMock struct {
	t minimock.Tester

	GetOriginFunc       func() (r insolar.NetworkNode)
	GetOriginCounter    uint64
	GetOriginPreCounter uint64
	GetOriginMock       mOriginProviderMockGetOrigin
}

//NewOriginProviderMock returns a mock for github.com/insolar/insolar/insolar.OriginProvider
func NewOriginProviderMock(t minimock.Tester) *OriginProviderMock {
	m := &OriginProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetOriginMock = mOriginProviderMockGetOrigin{mock: m}

	return m
}

type mOriginProviderMockGetOrigin struct {
	mock              *OriginProviderMock
	mainExpectation   *OriginProviderMockGetOriginExpectation
	expectationSeries []*OriginProviderMockGetOriginExpectation
}

type OriginProviderMockGetOriginExpectation struct {
	result *OriginProviderMockGetOriginResult
}

type OriginProviderMockGetOriginResult struct {
	r insolar.NetworkNode
}

//Expect specifies that invocation of OriginProvider.GetOrigin is expected from 1 to Infinity times
func (m *mOriginProviderMockGetOrigin) Expect() *mOriginProviderMockGetOrigin {
	m.mock.GetOriginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OriginProviderMockGetOriginExpectation{}
	}

	return m
}

//Return specifies results of invocation of OriginProvider.GetOrigin
func (m *mOriginProviderMockGetOrigin) Return(r insolar.NetworkNode) *OriginProviderMock {
	m.mock.GetOriginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OriginProviderMockGetOriginExpectation{}
	}
	m.mainExpectation.result = &OriginProviderMockGetOriginResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of OriginProvider.GetOrigin is expected once
func (m *mOriginProviderMockGetOrigin) ExpectOnce() *OriginProviderMockGetOriginExpectation {
	m.mock.GetOriginFunc = nil
	m.mainExpectation = nil

	expectation := &OriginProviderMockGetOriginExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OriginProviderMockGetOriginExpectation) Return(r insolar.NetworkNode) {
	e.result = &OriginProviderMockGetOriginResult{r}
}

//Set uses given function f as a mock of OriginProvider.GetOrigin method
func (m *mOriginProviderMockGetOrigin) Set(f func() (r insolar.NetworkNode)) *OriginProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOriginFunc = f
	return m.mock
}

//GetOrigin implements github.com/insolar/insolar/insolar.OriginProvider interface
func (m *OriginProviderMock) GetOrigin() (r insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetOriginPreCounter, 1)
	defer atomic.AddUint64(&m.GetOriginCounter, 1)

	if len(m.GetOriginMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOriginMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OriginProviderMock.GetOrigin.")
			return
		}

		result := m.GetOriginMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OriginProviderMock.GetOrigin")
			return
		}

		r = result.r

		return
	}

	if m.GetOriginMock.mainExpectation != nil {

		result := m.GetOriginMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OriginProviderMock.GetOrigin")
		}

		r = result.r

		return
	}

	if m.GetOriginFunc == nil {
		m.t.Fatalf("Unexpected call to OriginProviderMock.GetOrigin.")
		return
	}

	return m.GetOriginFunc()
}

//GetOriginMinimockCounter returns a count of OriginProviderMock.GetOriginFunc invocations
func (m *OriginProviderMock) GetOriginMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginCounter)
}

//GetOriginMinimockPreCounter returns the value of OriginProviderMock.GetOrigin invocations
func (m *OriginProviderMock) GetOriginMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOriginPreCounter)
}

//GetOriginFinished returns true if mock invocations count is ok
func (m *OriginProviderMock) GetOriginFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOriginMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOriginCounter) == uint64(len(m.GetOriginMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOriginMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOriginCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOriginFunc != nil {
		return atomic.LoadUint64(&m.GetOriginCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *OriginProviderMock) ValidateCallCounters() {

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to OriginProviderMock.GetOrigin")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *OriginProviderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *OriginProviderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *OriginProviderMock) MinimockFinish() {

	if !m.GetOriginFinished() {
		m.t.Fatal("Expected call to OriginProviderMock.GetOrigin")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *OriginProviderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *OriginProviderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetOriginFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetOriginFinished() {
				m.t.Error("Expected call to OriginProviderMock.GetOrigin")
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
func (m *OriginProviderMock) AllMocksCalled() bool {

	if !m.GetOriginFinished() {
		return false
	}

	return true
}
