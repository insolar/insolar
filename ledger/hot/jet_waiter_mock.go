package hot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetWaiter" can be found in github.com/insolar/insolar/ledger/hot
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetWaiterMock implements github.com/insolar/insolar/ledger/hot.JetWaiter
type JetWaiterMock struct {
	t minimock.Tester

	WaitFunc       func(p context.Context, p1 insolar.ID) (r error)
	WaitCounter    uint64
	WaitPreCounter uint64
	WaitMock       mJetWaiterMockWait
}

//NewJetWaiterMock returns a mock for github.com/insolar/insolar/ledger/hot.JetWaiter
func NewJetWaiterMock(t minimock.Tester) *JetWaiterMock {
	m := &JetWaiterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.WaitMock = mJetWaiterMockWait{mock: m}

	return m
}

type mJetWaiterMockWait struct {
	mock              *JetWaiterMock
	mainExpectation   *JetWaiterMockWaitExpectation
	expectationSeries []*JetWaiterMockWaitExpectation
}

type JetWaiterMockWaitExpectation struct {
	input  *JetWaiterMockWaitInput
	result *JetWaiterMockWaitResult
}

type JetWaiterMockWaitInput struct {
	p  context.Context
	p1 insolar.ID
}

type JetWaiterMockWaitResult struct {
	r error
}

//Expect specifies that invocation of JetWaiter.Wait is expected from 1 to Infinity times
func (m *mJetWaiterMockWait) Expect(p context.Context, p1 insolar.ID) *mJetWaiterMockWait {
	m.mock.WaitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetWaiterMockWaitExpectation{}
	}
	m.mainExpectation.input = &JetWaiterMockWaitInput{p, p1}
	return m
}

//Return specifies results of invocation of JetWaiter.Wait
func (m *mJetWaiterMockWait) Return(r error) *JetWaiterMock {
	m.mock.WaitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetWaiterMockWaitExpectation{}
	}
	m.mainExpectation.result = &JetWaiterMockWaitResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetWaiter.Wait is expected once
func (m *mJetWaiterMockWait) ExpectOnce(p context.Context, p1 insolar.ID) *JetWaiterMockWaitExpectation {
	m.mock.WaitFunc = nil
	m.mainExpectation = nil

	expectation := &JetWaiterMockWaitExpectation{}
	expectation.input = &JetWaiterMockWaitInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetWaiterMockWaitExpectation) Return(r error) {
	e.result = &JetWaiterMockWaitResult{r}
}

//Set uses given function f as a mock of JetWaiter.Wait method
func (m *mJetWaiterMockWait) Set(f func(p context.Context, p1 insolar.ID) (r error)) *JetWaiterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WaitFunc = f
	return m.mock
}

//Wait implements github.com/insolar/insolar/ledger/hot.JetWaiter interface
func (m *JetWaiterMock) Wait(p context.Context, p1 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.WaitPreCounter, 1)
	defer atomic.AddUint64(&m.WaitCounter, 1)

	if len(m.WaitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WaitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetWaiterMock.Wait. %v %v", p, p1)
			return
		}

		input := m.WaitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetWaiterMockWaitInput{p, p1}, "JetWaiter.Wait got unexpected parameters")

		result := m.WaitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetWaiterMock.Wait")
			return
		}

		r = result.r

		return
	}

	if m.WaitMock.mainExpectation != nil {

		input := m.WaitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetWaiterMockWaitInput{p, p1}, "JetWaiter.Wait got unexpected parameters")
		}

		result := m.WaitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetWaiterMock.Wait")
		}

		r = result.r

		return
	}

	if m.WaitFunc == nil {
		m.t.Fatalf("Unexpected call to JetWaiterMock.Wait. %v %v", p, p1)
		return
	}

	return m.WaitFunc(p, p1)
}

//WaitMinimockCounter returns a count of JetWaiterMock.WaitFunc invocations
func (m *JetWaiterMock) WaitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WaitCounter)
}

//WaitMinimockPreCounter returns the value of JetWaiterMock.Wait invocations
func (m *JetWaiterMock) WaitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WaitPreCounter)
}

//WaitFinished returns true if mock invocations count is ok
func (m *JetWaiterMock) WaitFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.WaitMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.WaitCounter) == uint64(len(m.WaitMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.WaitMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.WaitCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.WaitFunc != nil {
		return atomic.LoadUint64(&m.WaitCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetWaiterMock) ValidateCallCounters() {

	if !m.WaitFinished() {
		m.t.Fatal("Expected call to JetWaiterMock.Wait")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetWaiterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetWaiterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetWaiterMock) MinimockFinish() {

	if !m.WaitFinished() {
		m.t.Fatal("Expected call to JetWaiterMock.Wait")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetWaiterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetWaiterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.WaitFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.WaitFinished() {
				m.t.Error("Expected call to JetWaiterMock.Wait")
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
func (m *JetWaiterMock) AllMocksCalled() bool {

	if !m.WaitFinished() {
		return false
	}

	return true
}
