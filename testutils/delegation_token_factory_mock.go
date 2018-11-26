package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DelegationTokenFactory" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//DelegationTokenFactoryMock implements github.com/insolar/insolar/core.DelegationTokenFactory
type DelegationTokenFactoryMock struct {
	t minimock.Tester

	IssueGetObjectRedirectFunc       func(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error)
	IssueGetObjectRedirectCounter    uint64
	IssueGetObjectRedirectPreCounter uint64
	IssueGetObjectRedirectMock       mDelegationTokenFactoryMockIssueGetObjectRedirect

	IssuePendingExecutionFunc       func(p core.Message, p1 core.PulseNumber) (r core.DelegationToken, r1 error)
	IssuePendingExecutionCounter    uint64
	IssuePendingExecutionPreCounter uint64
	IssuePendingExecutionMock       mDelegationTokenFactoryMockIssuePendingExecution

	VerifyFunc       func(p core.Parcel) (r bool, r1 error)
	VerifyCounter    uint64
	VerifyPreCounter uint64
	VerifyMock       mDelegationTokenFactoryMockVerify
}

//NewDelegationTokenFactoryMock returns a mock for github.com/insolar/insolar/core.DelegationTokenFactory
func NewDelegationTokenFactoryMock(t minimock.Tester) *DelegationTokenFactoryMock {
	m := &DelegationTokenFactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IssueGetObjectRedirectMock = mDelegationTokenFactoryMockIssueGetObjectRedirect{mock: m}
	m.IssuePendingExecutionMock = mDelegationTokenFactoryMockIssuePendingExecution{mock: m}
	m.VerifyMock = mDelegationTokenFactoryMockVerify{mock: m}

	return m
}

type mDelegationTokenFactoryMockIssueGetObjectRedirect struct {
	mock             *DelegationTokenFactoryMock
	mockExpectations *DelegationTokenFactoryMockIssueGetObjectRedirectParams
}

//DelegationTokenFactoryMockIssueGetObjectRedirectParams represents input parameters of the DelegationTokenFactory.IssueGetObjectRedirect
type DelegationTokenFactoryMockIssueGetObjectRedirectParams struct {
	p  *core.RecordRef
	p1 core.Message
}

//Expect sets up expected params for the DelegationTokenFactory.IssueGetObjectRedirect
func (m *mDelegationTokenFactoryMockIssueGetObjectRedirect) Expect(p *core.RecordRef, p1 core.Message) *mDelegationTokenFactoryMockIssueGetObjectRedirect {
	m.mockExpectations = &DelegationTokenFactoryMockIssueGetObjectRedirectParams{p, p1}
	return m
}

//Return sets up a mock for DelegationTokenFactory.IssueGetObjectRedirect to return Return's arguments
func (m *mDelegationTokenFactoryMockIssueGetObjectRedirect) Return(r core.DelegationToken, r1 error) *DelegationTokenFactoryMock {
	m.mock.IssueGetObjectRedirectFunc = func(p *core.RecordRef, p1 core.Message) (core.DelegationToken, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of DelegationTokenFactory.IssueGetObjectRedirect method
func (m *mDelegationTokenFactoryMockIssueGetObjectRedirect) Set(f func(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error)) *DelegationTokenFactoryMock {
	m.mock.IssueGetObjectRedirectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IssueGetObjectRedirect implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) IssueGetObjectRedirect(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error) {
	atomic.AddUint64(&m.IssueGetObjectRedirectPreCounter, 1)
	defer atomic.AddUint64(&m.IssueGetObjectRedirectCounter, 1)

	if m.IssueGetObjectRedirectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IssueGetObjectRedirectMock.mockExpectations, DelegationTokenFactoryMockIssueGetObjectRedirectParams{p, p1},
			"DelegationTokenFactory.IssueGetObjectRedirect got unexpected parameters")

		if m.IssueGetObjectRedirectFunc == nil {

			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssueGetObjectRedirect")

			return
		}
	}

	if m.IssueGetObjectRedirectFunc == nil {
		m.t.Fatal("Unexpected call to DelegationTokenFactoryMock.IssueGetObjectRedirect")
		return
	}

	return m.IssueGetObjectRedirectFunc(p, p1)
}

//IssueGetObjectRedirectMinimockCounter returns a count of DelegationTokenFactoryMock.IssueGetObjectRedirectFunc invocations
func (m *DelegationTokenFactoryMock) IssueGetObjectRedirectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IssueGetObjectRedirectCounter)
}

//IssueGetObjectRedirectMinimockPreCounter returns the value of DelegationTokenFactoryMock.IssueGetObjectRedirect invocations
func (m *DelegationTokenFactoryMock) IssueGetObjectRedirectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IssueGetObjectRedirectPreCounter)
}

type mDelegationTokenFactoryMockIssuePendingExecution struct {
	mock             *DelegationTokenFactoryMock
	mockExpectations *DelegationTokenFactoryMockIssuePendingExecutionParams
}

//DelegationTokenFactoryMockIssuePendingExecutionParams represents input parameters of the DelegationTokenFactory.IssuePendingExecution
type DelegationTokenFactoryMockIssuePendingExecutionParams struct {
	p  core.Message
	p1 core.PulseNumber
}

//Expect sets up expected params for the DelegationTokenFactory.IssuePendingExecution
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Expect(p core.Message, p1 core.PulseNumber) *mDelegationTokenFactoryMockIssuePendingExecution {
	m.mockExpectations = &DelegationTokenFactoryMockIssuePendingExecutionParams{p, p1}
	return m
}

//Return sets up a mock for DelegationTokenFactory.IssuePendingExecution to return Return's arguments
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Return(r core.DelegationToken, r1 error) *DelegationTokenFactoryMock {
	m.mock.IssuePendingExecutionFunc = func(p core.Message, p1 core.PulseNumber) (core.DelegationToken, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of DelegationTokenFactory.IssuePendingExecution method
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Set(f func(p core.Message, p1 core.PulseNumber) (r core.DelegationToken, r1 error)) *DelegationTokenFactoryMock {
	m.mock.IssuePendingExecutionFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IssuePendingExecution implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) IssuePendingExecution(p core.Message, p1 core.PulseNumber) (r core.DelegationToken, r1 error) {
	atomic.AddUint64(&m.IssuePendingExecutionPreCounter, 1)
	defer atomic.AddUint64(&m.IssuePendingExecutionCounter, 1)

	if m.IssuePendingExecutionMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IssuePendingExecutionMock.mockExpectations, DelegationTokenFactoryMockIssuePendingExecutionParams{p, p1},
			"DelegationTokenFactory.IssuePendingExecution got unexpected parameters")

		if m.IssuePendingExecutionFunc == nil {

			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssuePendingExecution")

			return
		}
	}

	if m.IssuePendingExecutionFunc == nil {
		m.t.Fatal("Unexpected call to DelegationTokenFactoryMock.IssuePendingExecution")
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

type mDelegationTokenFactoryMockVerify struct {
	mock             *DelegationTokenFactoryMock
	mockExpectations *DelegationTokenFactoryMockVerifyParams
}

//DelegationTokenFactoryMockVerifyParams represents input parameters of the DelegationTokenFactory.Verify
type DelegationTokenFactoryMockVerifyParams struct {
	p core.Parcel
}

//Expect sets up expected params for the DelegationTokenFactory.Verify
func (m *mDelegationTokenFactoryMockVerify) Expect(p core.Parcel) *mDelegationTokenFactoryMockVerify {
	m.mockExpectations = &DelegationTokenFactoryMockVerifyParams{p}
	return m
}

//Return sets up a mock for DelegationTokenFactory.Verify to return Return's arguments
func (m *mDelegationTokenFactoryMockVerify) Return(r bool, r1 error) *DelegationTokenFactoryMock {
	m.mock.VerifyFunc = func(p core.Parcel) (bool, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of DelegationTokenFactory.Verify method
func (m *mDelegationTokenFactoryMockVerify) Set(f func(p core.Parcel) (r bool, r1 error)) *DelegationTokenFactoryMock {
	m.mock.VerifyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Verify implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) Verify(p core.Parcel) (r bool, r1 error) {
	atomic.AddUint64(&m.VerifyPreCounter, 1)
	defer atomic.AddUint64(&m.VerifyCounter, 1)

	if m.VerifyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.VerifyMock.mockExpectations, DelegationTokenFactoryMockVerifyParams{p},
			"DelegationTokenFactory.Verify got unexpected parameters")

		if m.VerifyFunc == nil {

			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.Verify")

			return
		}
	}

	if m.VerifyFunc == nil {
		m.t.Fatal("Unexpected call to DelegationTokenFactoryMock.Verify")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DelegationTokenFactoryMock) ValidateCallCounters() {

	if m.IssueGetObjectRedirectFunc != nil && atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) == 0 {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetObjectRedirect")
	}

	if m.IssuePendingExecutionFunc != nil && atomic.LoadUint64(&m.IssuePendingExecutionCounter) == 0 {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssuePendingExecution")
	}

	if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
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

	if m.IssueGetObjectRedirectFunc != nil && atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) == 0 {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetObjectRedirect")
	}

	if m.IssuePendingExecutionFunc != nil && atomic.LoadUint64(&m.IssuePendingExecutionCounter) == 0 {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssuePendingExecution")
	}

	if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
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
		ok = ok && (m.IssueGetObjectRedirectFunc == nil || atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) > 0)
		ok = ok && (m.IssuePendingExecutionFunc == nil || atomic.LoadUint64(&m.IssuePendingExecutionCounter) > 0)
		ok = ok && (m.VerifyFunc == nil || atomic.LoadUint64(&m.VerifyCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.IssueGetObjectRedirectFunc != nil && atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) == 0 {
				m.t.Error("Expected call to DelegationTokenFactoryMock.IssueGetObjectRedirect")
			}

			if m.IssuePendingExecutionFunc != nil && atomic.LoadUint64(&m.IssuePendingExecutionCounter) == 0 {
				m.t.Error("Expected call to DelegationTokenFactoryMock.IssuePendingExecution")
			}

			if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
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

	if m.IssueGetObjectRedirectFunc != nil && atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) == 0 {
		return false
	}

	if m.IssuePendingExecutionFunc != nil && atomic.LoadUint64(&m.IssuePendingExecutionCounter) == 0 {
		return false
	}

	if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
		return false
	}

	return true
}
