package artifacts

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "cache" can be found in github.com/insolar/insolar/logicrunner/artifacts
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//cacheMock implements github.com/insolar/insolar/logicrunner/artifacts.cache
type cacheMock struct {
	t minimock.Tester

	getFunc       func(p insolar.Reference, p1 func() (r interface{}, r1 error)) (r interface{}, r1 error)
	getCounter    uint64
	getPreCounter uint64
	getMock       mcacheMockget
}

//NewcacheMock returns a mock for github.com/insolar/insolar/logicrunner/artifacts.cache
func NewcacheMock(t minimock.Tester) *cacheMock {
	m := &cacheMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.getMock = mcacheMockget{mock: m}

	return m
}

type mcacheMockget struct {
	mock              *cacheMock
	mainExpectation   *cacheMockgetExpectation
	expectationSeries []*cacheMockgetExpectation
}

type cacheMockgetExpectation struct {
	input  *cacheMockgetInput
	result *cacheMockgetResult
}

type cacheMockgetInput struct {
	p  insolar.Reference
	p1 func() (r interface{}, r1 error)
}

type cacheMockgetResult struct {
	r  interface{}
	r1 error
}

//Expect specifies that invocation of cache.get is expected from 1 to Infinity times
func (m *mcacheMockget) Expect(p insolar.Reference, p1 func() (r interface{}, r1 error)) *mcacheMockget {
	m.mock.getFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &cacheMockgetExpectation{}
	}
	m.mainExpectation.input = &cacheMockgetInput{p, p1}
	return m
}

//Return specifies results of invocation of cache.get
func (m *mcacheMockget) Return(r interface{}, r1 error) *cacheMock {
	m.mock.getFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &cacheMockgetExpectation{}
	}
	m.mainExpectation.result = &cacheMockgetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of cache.get is expected once
func (m *mcacheMockget) ExpectOnce(p insolar.Reference, p1 func() (r interface{}, r1 error)) *cacheMockgetExpectation {
	m.mock.getFunc = nil
	m.mainExpectation = nil

	expectation := &cacheMockgetExpectation{}
	expectation.input = &cacheMockgetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *cacheMockgetExpectation) Return(r interface{}, r1 error) {
	e.result = &cacheMockgetResult{r, r1}
}

//Set uses given function f as a mock of cache.get method
func (m *mcacheMockget) Set(f func(p insolar.Reference, p1 func() (r interface{}, r1 error)) (r interface{}, r1 error)) *cacheMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.getFunc = f
	return m.mock
}

//get implements github.com/insolar/insolar/logicrunner/artifacts.cache interface
func (m *cacheMock) get(p insolar.Reference, p1 func() (r interface{}, r1 error)) (r interface{}, r1 error) {
	counter := atomic.AddUint64(&m.getPreCounter, 1)
	defer atomic.AddUint64(&m.getCounter, 1)

	if len(m.getMock.expectationSeries) > 0 {
		if counter > uint64(len(m.getMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to cacheMock.get. %v %v", p, p1)
			return
		}

		input := m.getMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, cacheMockgetInput{p, p1}, "cache.get got unexpected parameters")

		result := m.getMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the cacheMock.get")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.getMock.mainExpectation != nil {

		input := m.getMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, cacheMockgetInput{p, p1}, "cache.get got unexpected parameters")
		}

		result := m.getMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the cacheMock.get")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.getFunc == nil {
		m.t.Fatalf("Unexpected call to cacheMock.get. %v %v", p, p1)
		return
	}

	return m.getFunc(p, p1)
}

//getMinimockCounter returns a count of cacheMock.getFunc invocations
func (m *cacheMock) getMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.getCounter)
}

//getMinimockPreCounter returns the value of cacheMock.get invocations
func (m *cacheMock) getMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.getPreCounter)
}

//getFinished returns true if mock invocations count is ok
func (m *cacheMock) getFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.getMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.getCounter) == uint64(len(m.getMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.getMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.getCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.getFunc != nil {
		return atomic.LoadUint64(&m.getCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *cacheMock) ValidateCallCounters() {

	if !m.getFinished() {
		m.t.Fatal("Expected call to cacheMock.get")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *cacheMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *cacheMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *cacheMock) MinimockFinish() {

	if !m.getFinished() {
		m.t.Fatal("Expected call to cacheMock.get")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *cacheMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *cacheMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.getFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.getFinished() {
				m.t.Error("Expected call to cacheMock.get")
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
func (m *cacheMock) AllMocksCalled() bool {

	if !m.getFinished() {
		return false
	}

	return true
}
