package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "dropTruncater" can be found in github.com/insolar/insolar/ledger/heavy/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//dropTruncaterMock implements github.com/insolar/insolar/ledger/heavy/executor.dropTruncater
type dropTruncaterMock struct {
	t minimock.Tester

	TruncateHeadFunc       func(p context.Context, p1 insolar.PulseNumber) (r error)
	TruncateHeadCounter    uint64
	TruncateHeadPreCounter uint64
	TruncateHeadMock       mdropTruncaterMockTruncateHead
}

//NewdropTruncaterMock returns a mock for github.com/insolar/insolar/ledger/heavy/executor.dropTruncater
func NewdropTruncaterMock(t minimock.Tester) *dropTruncaterMock {
	m := &dropTruncaterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.TruncateHeadMock = mdropTruncaterMockTruncateHead{mock: m}

	return m
}

type mdropTruncaterMockTruncateHead struct {
	mock              *dropTruncaterMock
	mainExpectation   *dropTruncaterMockTruncateHeadExpectation
	expectationSeries []*dropTruncaterMockTruncateHeadExpectation
}

type dropTruncaterMockTruncateHeadExpectation struct {
	input  *dropTruncaterMockTruncateHeadInput
	result *dropTruncaterMockTruncateHeadResult
}

type dropTruncaterMockTruncateHeadInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type dropTruncaterMockTruncateHeadResult struct {
	r error
}

//Expect specifies that invocation of dropTruncater.TruncateHead is expected from 1 to Infinity times
func (m *mdropTruncaterMockTruncateHead) Expect(p context.Context, p1 insolar.PulseNumber) *mdropTruncaterMockTruncateHead {
	m.mock.TruncateHeadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &dropTruncaterMockTruncateHeadExpectation{}
	}
	m.mainExpectation.input = &dropTruncaterMockTruncateHeadInput{p, p1}
	return m
}

//Return specifies results of invocation of dropTruncater.TruncateHead
func (m *mdropTruncaterMockTruncateHead) Return(r error) *dropTruncaterMock {
	m.mock.TruncateHeadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &dropTruncaterMockTruncateHeadExpectation{}
	}
	m.mainExpectation.result = &dropTruncaterMockTruncateHeadResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of dropTruncater.TruncateHead is expected once
func (m *mdropTruncaterMockTruncateHead) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *dropTruncaterMockTruncateHeadExpectation {
	m.mock.TruncateHeadFunc = nil
	m.mainExpectation = nil

	expectation := &dropTruncaterMockTruncateHeadExpectation{}
	expectation.input = &dropTruncaterMockTruncateHeadInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *dropTruncaterMockTruncateHeadExpectation) Return(r error) {
	e.result = &dropTruncaterMockTruncateHeadResult{r}
}

//Set uses given function f as a mock of dropTruncater.TruncateHead method
func (m *mdropTruncaterMockTruncateHead) Set(f func(p context.Context, p1 insolar.PulseNumber) (r error)) *dropTruncaterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.TruncateHeadFunc = f
	return m.mock
}

//TruncateHead implements github.com/insolar/insolar/ledger/heavy/executor.dropTruncater interface
func (m *dropTruncaterMock) TruncateHead(p context.Context, p1 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.TruncateHeadPreCounter, 1)
	defer atomic.AddUint64(&m.TruncateHeadCounter, 1)

	if len(m.TruncateHeadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.TruncateHeadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to dropTruncaterMock.TruncateHead. %v %v", p, p1)
			return
		}

		input := m.TruncateHeadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, dropTruncaterMockTruncateHeadInput{p, p1}, "dropTruncater.TruncateHead got unexpected parameters")

		result := m.TruncateHeadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the dropTruncaterMock.TruncateHead")
			return
		}

		r = result.r

		return
	}

	if m.TruncateHeadMock.mainExpectation != nil {

		input := m.TruncateHeadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, dropTruncaterMockTruncateHeadInput{p, p1}, "dropTruncater.TruncateHead got unexpected parameters")
		}

		result := m.TruncateHeadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the dropTruncaterMock.TruncateHead")
		}

		r = result.r

		return
	}

	if m.TruncateHeadFunc == nil {
		m.t.Fatalf("Unexpected call to dropTruncaterMock.TruncateHead. %v %v", p, p1)
		return
	}

	return m.TruncateHeadFunc(p, p1)
}

//TruncateHeadMinimockCounter returns a count of dropTruncaterMock.TruncateHeadFunc invocations
func (m *dropTruncaterMock) TruncateHeadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.TruncateHeadCounter)
}

//TruncateHeadMinimockPreCounter returns the value of dropTruncaterMock.TruncateHead invocations
func (m *dropTruncaterMock) TruncateHeadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.TruncateHeadPreCounter)
}

//TruncateHeadFinished returns true if mock invocations count is ok
func (m *dropTruncaterMock) TruncateHeadFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.TruncateHeadMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.TruncateHeadCounter) == uint64(len(m.TruncateHeadMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.TruncateHeadMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.TruncateHeadCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.TruncateHeadFunc != nil {
		return atomic.LoadUint64(&m.TruncateHeadCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *dropTruncaterMock) ValidateCallCounters() {

	if !m.TruncateHeadFinished() {
		m.t.Fatal("Expected call to dropTruncaterMock.TruncateHead")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *dropTruncaterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *dropTruncaterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *dropTruncaterMock) MinimockFinish() {

	if !m.TruncateHeadFinished() {
		m.t.Fatal("Expected call to dropTruncaterMock.TruncateHead")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *dropTruncaterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *dropTruncaterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.TruncateHeadFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.TruncateHeadFinished() {
				m.t.Error("Expected call to dropTruncaterMock.TruncateHead")
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
func (m *dropTruncaterMock) AllMocksCalled() bool {

	if !m.TruncateHeadFinished() {
		return false
	}

	return true
}
