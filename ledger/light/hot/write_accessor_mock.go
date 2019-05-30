package hot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "WriteAccessor" can be found in github.com/insolar/insolar/ledger/light/hot
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//WriteAccessorMock implements github.com/insolar/insolar/ledger/light/hot.WriteAccessor
type WriteAccessorMock struct {
	t minimock.Tester

	BeginFunc       func(p context.Context, p1 insolar.PulseNumber) (r func(), r1 error)
	BeginCounter    uint64
	BeginPreCounter uint64
	BeginMock       mWriteAccessorMockBegin
}

//NewWriteAccessorMock returns a mock for github.com/insolar/insolar/ledger/light/hot.WriteAccessor
func NewWriteAccessorMock(t minimock.Tester) *WriteAccessorMock {
	m := &WriteAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BeginMock = mWriteAccessorMockBegin{mock: m}

	return m
}

type mWriteAccessorMockBegin struct {
	mock              *WriteAccessorMock
	mainExpectation   *WriteAccessorMockBeginExpectation
	expectationSeries []*WriteAccessorMockBeginExpectation
}

type WriteAccessorMockBeginExpectation struct {
	input  *WriteAccessorMockBeginInput
	result *WriteAccessorMockBeginResult
}

type WriteAccessorMockBeginInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type WriteAccessorMockBeginResult struct {
	r  func()
	r1 error
}

//Expect specifies that invocation of WriteAccessor.Begin is expected from 1 to Infinity times
func (m *mWriteAccessorMockBegin) Expect(p context.Context, p1 insolar.PulseNumber) *mWriteAccessorMockBegin {
	m.mock.BeginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WriteAccessorMockBeginExpectation{}
	}
	m.mainExpectation.input = &WriteAccessorMockBeginInput{p, p1}
	return m
}

//Return specifies results of invocation of WriteAccessor.Begin
func (m *mWriteAccessorMockBegin) Return(r func(), r1 error) *WriteAccessorMock {
	m.mock.BeginFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WriteAccessorMockBeginExpectation{}
	}
	m.mainExpectation.result = &WriteAccessorMockBeginResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of WriteAccessor.Begin is expected once
func (m *mWriteAccessorMockBegin) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *WriteAccessorMockBeginExpectation {
	m.mock.BeginFunc = nil
	m.mainExpectation = nil

	expectation := &WriteAccessorMockBeginExpectation{}
	expectation.input = &WriteAccessorMockBeginInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *WriteAccessorMockBeginExpectation) Return(r func(), r1 error) {
	e.result = &WriteAccessorMockBeginResult{r, r1}
}

//Set uses given function f as a mock of WriteAccessor.Begin method
func (m *mWriteAccessorMockBegin) Set(f func(p context.Context, p1 insolar.PulseNumber) (r func(), r1 error)) *WriteAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BeginFunc = f
	return m.mock
}

//Begin implements github.com/insolar/insolar/ledger/light/hot.WriteAccessor interface
func (m *WriteAccessorMock) Begin(p context.Context, p1 insolar.PulseNumber) (r func(), r1 error) {
	counter := atomic.AddUint64(&m.BeginPreCounter, 1)
	defer atomic.AddUint64(&m.BeginCounter, 1)

	if len(m.BeginMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BeginMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to WriteAccessorMock.Begin. %v %v", p, p1)
			return
		}

		input := m.BeginMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, WriteAccessorMockBeginInput{p, p1}, "WriteAccessor.Begin got unexpected parameters")

		result := m.BeginMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the WriteAccessorMock.Begin")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BeginMock.mainExpectation != nil {

		input := m.BeginMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, WriteAccessorMockBeginInput{p, p1}, "WriteAccessor.Begin got unexpected parameters")
		}

		result := m.BeginMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the WriteAccessorMock.Begin")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BeginFunc == nil {
		m.t.Fatalf("Unexpected call to WriteAccessorMock.Begin. %v %v", p, p1)
		return
	}

	return m.BeginFunc(p, p1)
}

//BeginMinimockCounter returns a count of WriteAccessorMock.BeginFunc invocations
func (m *WriteAccessorMock) BeginMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BeginCounter)
}

//BeginMinimockPreCounter returns the value of WriteAccessorMock.Begin invocations
func (m *WriteAccessorMock) BeginMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BeginPreCounter)
}

//BeginFinished returns true if mock invocations count is ok
func (m *WriteAccessorMock) BeginFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.BeginMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.BeginCounter) == uint64(len(m.BeginMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.BeginMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.BeginCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.BeginFunc != nil {
		return atomic.LoadUint64(&m.BeginCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *WriteAccessorMock) ValidateCallCounters() {

	if !m.BeginFinished() {
		m.t.Fatal("Expected call to WriteAccessorMock.Begin")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *WriteAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *WriteAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *WriteAccessorMock) MinimockFinish() {

	if !m.BeginFinished() {
		m.t.Fatal("Expected call to WriteAccessorMock.Begin")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *WriteAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *WriteAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.BeginFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.BeginFinished() {
				m.t.Error("Expected call to WriteAccessorMock.Begin")
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
func (m *WriteAccessorMock) AllMocksCalled() bool {

	if !m.BeginFinished() {
		return false
	}

	return true
}
