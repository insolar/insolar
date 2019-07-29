package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DelegationTokenFactory" can be found in github.com/insolar/insolar/insolar
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//DelegationTokenFactoryMock implements github.com/insolar/insolar/insolar.DelegationTokenFactory
type DelegationTokenFactoryMock struct {
	t minimock.Tester

	IssuePendingExecutionFunc       func(p insolar.Message, p1 insolar.PulseNumber) (r insolar.DelegationToken, r1 error)
	IssuePendingExecutionCounter    uint64
	IssuePendingExecutionPreCounter uint64
	IssuePendingExecutionMock       mDelegationTokenFactoryMockIssuePendingExecution

	VerifyFunc       func(p insolar.Parcel) (r bool, r1 error)
	VerifyCounter    uint64
	VerifyPreCounter uint64
	VerifyMock       mDelegationTokenFactoryMockVerify
}

//NewDelegationTokenFactoryMock returns a mock for github.com/insolar/insolar/insolar.DelegationTokenFactory
func NewDelegationTokenFactoryMock(t minimock.Tester) *DelegationTokenFactoryMock {
	m := &DelegationTokenFactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IssuePendingExecutionMock = mDelegationTokenFactoryMockIssuePendingExecution{mock: m}
	m.VerifyMock = mDelegationTokenFactoryMockVerify{mock: m}

	return m
}

type mDelegationTokenFactoryMockIssuePendingExecution struct {
	mock              *DelegationTokenFactoryMock
	mainExpectation   *DelegationTokenFactoryMockIssuePendingExecutionExpectation
	expectationSeries []*DelegationTokenFactoryMockIssuePendingExecutionExpectation
}

type DelegationTokenFactoryMockIssuePendingExecutionExpectation struct {
	input  *DelegationTokenFactoryMockIssuePendingExecutionInput
	result *DelegationTokenFactoryMockIssuePendingExecutionResult
}

type DelegationTokenFactoryMockIssuePendingExecutionInput struct {
	p  insolar.Message
	p1 insolar.PulseNumber
}

type DelegationTokenFactoryMockIssuePendingExecutionResult struct {
	r  insolar.DelegationToken
	r1 error
}

//Expect specifies that invocation of DelegationTokenFactory.IssuePendingExecution is expected from 1 to Infinity times
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Expect(p insolar.Message, p1 insolar.PulseNumber) *mDelegationTokenFactoryMockIssuePendingExecution {
	m.mock.IssuePendingExecutionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssuePendingExecutionExpectation{}
	}
	m.mainExpectation.input = &DelegationTokenFactoryMockIssuePendingExecutionInput{p, p1}
	return m
}

//Return specifies results of invocation of DelegationTokenFactory.IssuePendingExecution
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Return(r insolar.DelegationToken, r1 error) *DelegationTokenFactoryMock {
	m.mock.IssuePendingExecutionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssuePendingExecutionExpectation{}
	}
	m.mainExpectation.result = &DelegationTokenFactoryMockIssuePendingExecutionResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DelegationTokenFactory.IssuePendingExecution is expected once
func (m *mDelegationTokenFactoryMockIssuePendingExecution) ExpectOnce(p insolar.Message, p1 insolar.PulseNumber) *DelegationTokenFactoryMockIssuePendingExecutionExpectation {
	m.mock.IssuePendingExecutionFunc = nil
	m.mainExpectation = nil

	expectation := &DelegationTokenFactoryMockIssuePendingExecutionExpectation{}
	expectation.input = &DelegationTokenFactoryMockIssuePendingExecutionInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DelegationTokenFactoryMockIssuePendingExecutionExpectation) Return(r insolar.DelegationToken, r1 error) {
	e.result = &DelegationTokenFactoryMockIssuePendingExecutionResult{r, r1}
}

//Set uses given function f as a mock of DelegationTokenFactory.IssuePendingExecution method
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Set(f func(p insolar.Message, p1 insolar.PulseNumber) (r insolar.DelegationToken, r1 error)) *DelegationTokenFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IssuePendingExecutionFunc = f
	return m.mock
}

//IssuePendingExecution implements github.com/insolar/insolar/insolar.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) IssuePendingExecution(p insolar.Message, p1 insolar.PulseNumber) (r insolar.DelegationToken, r1 error) {
	counter := atomic.AddUint64(&m.IssuePendingExecutionPreCounter, 1)
	defer atomic.AddUint64(&m.IssuePendingExecutionCounter, 1)

	if len(m.IssuePendingExecutionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IssuePendingExecutionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssuePendingExecution. %v %v", p, p1)
			return
		}

		input := m.IssuePendingExecutionMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssuePendingExecutionInput{p, p1}, "DelegationTokenFactory.IssuePendingExecution got unexpected parameters")

		result := m.IssuePendingExecutionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssuePendingExecution")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssuePendingExecutionMock.mainExpectation != nil {

		input := m.IssuePendingExecutionMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssuePendingExecutionInput{p, p1}, "DelegationTokenFactory.IssuePendingExecution got unexpected parameters")
		}

		result := m.IssuePendingExecutionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssuePendingExecution")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssuePendingExecutionFunc == nil {
		m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssuePendingExecution. %v %v", p, p1)
		return
	}

	return m.IssuePendingExecutionFunc(p, p1)
}

//IssuePendingExecutionMinimockCounter returns a count of DelegationTokenFactoryMock.IssuePendingExecutionFunc invocations
func (m *DelegationTokenFactoryMock) IssuePendingExecutionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IssuePendingExecutionCounter)
}

//IssuePendingExecutionMinimockPreCounter returns the value of DelegationTokenFactoryMock.IssuePendingExecution invocations
func (m *DelegationTokenFactoryMock) IssuePendingExecutionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IssuePendingExecutionPreCounter)
}

//IssuePendingExecutionFinished returns true if mock invocations count is ok
func (m *DelegationTokenFactoryMock) IssuePendingExecutionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IssuePendingExecutionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IssuePendingExecutionCounter) == uint64(len(m.IssuePendingExecutionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IssuePendingExecutionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IssuePendingExecutionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IssuePendingExecutionFunc != nil {
		return atomic.LoadUint64(&m.IssuePendingExecutionCounter) > 0
	}

	return true
}

type mDelegationTokenFactoryMockVerify struct {
	mock              *DelegationTokenFactoryMock
	mainExpectation   *DelegationTokenFactoryMockVerifyExpectation
	expectationSeries []*DelegationTokenFactoryMockVerifyExpectation
}

type DelegationTokenFactoryMockVerifyExpectation struct {
	input  *DelegationTokenFactoryMockVerifyInput
	result *DelegationTokenFactoryMockVerifyResult
}

type DelegationTokenFactoryMockVerifyInput struct {
	p insolar.Parcel
}

type DelegationTokenFactoryMockVerifyResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of DelegationTokenFactory.Verify is expected from 1 to Infinity times
func (m *mDelegationTokenFactoryMockVerify) Expect(p insolar.Parcel) *mDelegationTokenFactoryMockVerify {
	m.mock.VerifyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockVerifyExpectation{}
	}
	m.mainExpectation.input = &DelegationTokenFactoryMockVerifyInput{p}
	return m
}

//Return specifies results of invocation of DelegationTokenFactory.Verify
func (m *mDelegationTokenFactoryMockVerify) Return(r bool, r1 error) *DelegationTokenFactoryMock {
	m.mock.VerifyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockVerifyExpectation{}
	}
	m.mainExpectation.result = &DelegationTokenFactoryMockVerifyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DelegationTokenFactory.Verify is expected once
func (m *mDelegationTokenFactoryMockVerify) ExpectOnce(p insolar.Parcel) *DelegationTokenFactoryMockVerifyExpectation {
	m.mock.VerifyFunc = nil
	m.mainExpectation = nil

	expectation := &DelegationTokenFactoryMockVerifyExpectation{}
	expectation.input = &DelegationTokenFactoryMockVerifyInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DelegationTokenFactoryMockVerifyExpectation) Return(r bool, r1 error) {
	e.result = &DelegationTokenFactoryMockVerifyResult{r, r1}
}

//Set uses given function f as a mock of DelegationTokenFactory.Verify method
func (m *mDelegationTokenFactoryMockVerify) Set(f func(p insolar.Parcel) (r bool, r1 error)) *DelegationTokenFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VerifyFunc = f
	return m.mock
}

//Verify implements github.com/insolar/insolar/insolar.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) Verify(p insolar.Parcel) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.VerifyPreCounter, 1)
	defer atomic.AddUint64(&m.VerifyCounter, 1)

	if len(m.VerifyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VerifyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.Verify. %v", p)
			return
		}

		input := m.VerifyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockVerifyInput{p}, "DelegationTokenFactory.Verify got unexpected parameters")

		result := m.VerifyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.Verify")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VerifyMock.mainExpectation != nil {

		input := m.VerifyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockVerifyInput{p}, "DelegationTokenFactory.Verify got unexpected parameters")
		}

		result := m.VerifyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.Verify")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VerifyFunc == nil {
		m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.Verify. %v", p)
		return
	}

	return m.VerifyFunc(p)
}

//VerifyMinimockCounter returns a count of DelegationTokenFactoryMock.VerifyFunc invocations
func (m *DelegationTokenFactoryMock) VerifyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyCounter)
}

//VerifyMinimockPreCounter returns the value of DelegationTokenFactoryMock.Verify invocations
func (m *DelegationTokenFactoryMock) VerifyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyPreCounter)
}

//VerifyFinished returns true if mock invocations count is ok
func (m *DelegationTokenFactoryMock) VerifyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.VerifyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.VerifyCounter) == uint64(len(m.VerifyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.VerifyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.VerifyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.VerifyFunc != nil {
		return atomic.LoadUint64(&m.VerifyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DelegationTokenFactoryMock) ValidateCallCounters() {

	if !m.IssuePendingExecutionFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssuePendingExecution")
	}

	if !m.VerifyFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.Verify")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DelegationTokenFactoryMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DelegationTokenFactoryMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DelegationTokenFactoryMock) MinimockFinish() {

	if !m.IssuePendingExecutionFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssuePendingExecution")
	}

	if !m.VerifyFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.Verify")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DelegationTokenFactoryMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DelegationTokenFactoryMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.IssuePendingExecutionFinished()
		ok = ok && m.VerifyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IssuePendingExecutionFinished() {
				m.t.Error("Expected call to DelegationTokenFactoryMock.IssuePendingExecution")
			}

			if !m.VerifyFinished() {
				m.t.Error("Expected call to DelegationTokenFactoryMock.Verify")
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
func (m *DelegationTokenFactoryMock) AllMocksCalled() bool {

	if !m.IssuePendingExecutionFinished() {
		return false
	}

	if !m.VerifyFinished() {
		return false
	}

	return true
}
