package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LedgerLogic" can be found in github.com/insolar/insolar/ledger/artifactmanager
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//LedgerLogicMock implements github.com/insolar/insolar/ledger/artifactmanager.LedgerLogic
type LedgerLogicMock struct {
	t minimock.Tester

	GetCodeFunc       func(p context.Context, p1 insolar.Parcel) (r insolar.Reply, r1 error)
	GetCodeCounter    uint64
	GetCodePreCounter uint64
	GetCodeMock       mLedgerLogicMockGetCode
}

//NewLedgerLogicMock returns a mock for github.com/insolar/insolar/ledger/artifactmanager.LedgerLogic
func NewLedgerLogicMock(t minimock.Tester) *LedgerLogicMock {
	m := &LedgerLogicMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCodeMock = mLedgerLogicMockGetCode{mock: m}

	return m
}

type mLedgerLogicMockGetCode struct {
	mock              *LedgerLogicMock
	mainExpectation   *LedgerLogicMockGetCodeExpectation
	expectationSeries []*LedgerLogicMockGetCodeExpectation
}

type LedgerLogicMockGetCodeExpectation struct {
	input  *LedgerLogicMockGetCodeInput
	result *LedgerLogicMockGetCodeResult
}

type LedgerLogicMockGetCodeInput struct {
	p  context.Context
	p1 insolar.Parcel
}

type LedgerLogicMockGetCodeResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of LedgerLogic.GetCode is expected from 1 to Infinity times
func (m *mLedgerLogicMockGetCode) Expect(p context.Context, p1 insolar.Parcel) *mLedgerLogicMockGetCode {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LedgerLogicMockGetCodeExpectation{}
	}
	m.mainExpectation.input = &LedgerLogicMockGetCodeInput{p, p1}
	return m
}

//Return specifies results of invocation of LedgerLogic.GetCode
func (m *mLedgerLogicMockGetCode) Return(r insolar.Reply, r1 error) *LedgerLogicMock {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LedgerLogicMockGetCodeExpectation{}
	}
	m.mainExpectation.result = &LedgerLogicMockGetCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LedgerLogic.GetCode is expected once
func (m *mLedgerLogicMockGetCode) ExpectOnce(p context.Context, p1 insolar.Parcel) *LedgerLogicMockGetCodeExpectation {
	m.mock.GetCodeFunc = nil
	m.mainExpectation = nil

	expectation := &LedgerLogicMockGetCodeExpectation{}
	expectation.input = &LedgerLogicMockGetCodeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LedgerLogicMockGetCodeExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &LedgerLogicMockGetCodeResult{r, r1}
}

//Set uses given function f as a mock of LedgerLogic.GetCode method
func (m *mLedgerLogicMockGetCode) Set(f func(p context.Context, p1 insolar.Parcel) (r insolar.Reply, r1 error)) *LedgerLogicMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCodeFunc = f
	return m.mock
}

//GetCode implements github.com/insolar/insolar/ledger/artifactmanager.LedgerLogic interface
func (m *LedgerLogicMock) GetCode(p context.Context, p1 insolar.Parcel) (r insolar.Reply, r1 error) {
	counter := atomic.AddUint64(&m.GetCodePreCounter, 1)
	defer atomic.AddUint64(&m.GetCodeCounter, 1)

	if len(m.GetCodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LedgerLogicMock.GetCode. %v %v", p, p1)
			return
		}

		input := m.GetCodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LedgerLogicMockGetCodeInput{p, p1}, "LedgerLogic.GetCode got unexpected parameters")

		result := m.GetCodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LedgerLogicMock.GetCode")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeMock.mainExpectation != nil {

		input := m.GetCodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LedgerLogicMockGetCodeInput{p, p1}, "LedgerLogic.GetCode got unexpected parameters")
		}

		result := m.GetCodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LedgerLogicMock.GetCode")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeFunc == nil {
		m.t.Fatalf("Unexpected call to LedgerLogicMock.GetCode. %v %v", p, p1)
		return
	}

	return m.GetCodeFunc(p, p1)
}

//GetCodeMinimockCounter returns a count of LedgerLogicMock.GetCodeFunc invocations
func (m *LedgerLogicMock) GetCodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodeCounter)
}

//GetCodeMinimockPreCounter returns the value of LedgerLogicMock.GetCode invocations
func (m *LedgerLogicMock) GetCodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodePreCounter)
}

//GetCodeFinished returns true if mock invocations count is ok
func (m *LedgerLogicMock) GetCodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCodeCounter) == uint64(len(m.GetCodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCodeFunc != nil {
		return atomic.LoadUint64(&m.GetCodeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LedgerLogicMock) ValidateCallCounters() {

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to LedgerLogicMock.GetCode")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LedgerLogicMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LedgerLogicMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LedgerLogicMock) MinimockFinish() {

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to LedgerLogicMock.GetCode")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LedgerLogicMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LedgerLogicMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetCodeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetCodeFinished() {
				m.t.Error("Expected call to LedgerLogicMock.GetCode")
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
func (m *LedgerLogicMock) AllMocksCalled() bool {

	if !m.GetCodeFinished() {
		return false
	}

	return true
}
