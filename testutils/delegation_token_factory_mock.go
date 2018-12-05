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

	IssueGetChildrenRedirectFunc       func(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error)
	IssueGetChildrenRedirectCounter    uint64
	IssueGetChildrenRedirectPreCounter uint64
	IssueGetChildrenRedirectMock       mDelegationTokenFactoryMockIssueGetChildrenRedirect

	IssueGetCodeRedirectFunc       func(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error)
	IssueGetCodeRedirectCounter    uint64
	IssueGetCodeRedirectPreCounter uint64
	IssueGetCodeRedirectMock       mDelegationTokenFactoryMockIssueGetCodeRedirect

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

	m.IssueGetChildrenRedirectMock = mDelegationTokenFactoryMockIssueGetChildrenRedirect{mock: m}
	m.IssueGetCodeRedirectMock = mDelegationTokenFactoryMockIssueGetCodeRedirect{mock: m}
	m.IssueGetObjectRedirectMock = mDelegationTokenFactoryMockIssueGetObjectRedirect{mock: m}
	m.IssuePendingExecutionMock = mDelegationTokenFactoryMockIssuePendingExecution{mock: m}
	m.VerifyMock = mDelegationTokenFactoryMockVerify{mock: m}

	return m
}

type mDelegationTokenFactoryMockIssueGetChildrenRedirect struct {
	mock              *DelegationTokenFactoryMock
	mainExpectation   *DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation
	expectationSeries []*DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation
}

type DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation struct {
	input  *DelegationTokenFactoryMockIssueGetChildrenRedirectInput
	result *DelegationTokenFactoryMockIssueGetChildrenRedirectResult
}

type DelegationTokenFactoryMockIssueGetChildrenRedirectInput struct {
	p  *core.RecordRef
	p1 core.Message
}

type DelegationTokenFactoryMockIssueGetChildrenRedirectResult struct {
	r  core.DelegationToken
	r1 error
}

//Expect specifies that invocation of DelegationTokenFactory.IssueGetChildrenRedirect is expected from 1 to Infinity times
func (m *mDelegationTokenFactoryMockIssueGetChildrenRedirect) Expect(p *core.RecordRef, p1 core.Message) *mDelegationTokenFactoryMockIssueGetChildrenRedirect {
	m.mock.IssueGetChildrenRedirectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation{}
	}
	m.mainExpectation.input = &DelegationTokenFactoryMockIssueGetChildrenRedirectInput{p, p1}
	return m
}

//Return specifies results of invocation of DelegationTokenFactory.IssueGetChildrenRedirect
func (m *mDelegationTokenFactoryMockIssueGetChildrenRedirect) Return(r core.DelegationToken, r1 error) *DelegationTokenFactoryMock {
	m.mock.IssueGetChildrenRedirectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation{}
	}
	m.mainExpectation.result = &DelegationTokenFactoryMockIssueGetChildrenRedirectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DelegationTokenFactory.IssueGetChildrenRedirect is expected once
func (m *mDelegationTokenFactoryMockIssueGetChildrenRedirect) ExpectOnce(p *core.RecordRef, p1 core.Message) *DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation {
	m.mock.IssueGetChildrenRedirectFunc = nil
	m.mainExpectation = nil

	expectation := &DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation{}
	expectation.input = &DelegationTokenFactoryMockIssueGetChildrenRedirectInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DelegationTokenFactoryMockIssueGetChildrenRedirectExpectation) Return(r core.DelegationToken, r1 error) {
	e.result = &DelegationTokenFactoryMockIssueGetChildrenRedirectResult{r, r1}
}

//Set uses given function f as a mock of DelegationTokenFactory.IssueGetChildrenRedirect method
func (m *mDelegationTokenFactoryMockIssueGetChildrenRedirect) Set(f func(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error)) *DelegationTokenFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IssueGetChildrenRedirectFunc = f
	return m.mock
}

//IssueGetChildrenRedirect implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) IssueGetChildrenRedirect(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error) {
	counter := atomic.AddUint64(&m.IssueGetChildrenRedirectPreCounter, 1)
	defer atomic.AddUint64(&m.IssueGetChildrenRedirectCounter, 1)

	if len(m.IssueGetChildrenRedirectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IssueGetChildrenRedirectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssueGetChildrenRedirect. %v %v", p, p1)
			return
		}

		input := m.IssueGetChildrenRedirectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssueGetChildrenRedirectInput{p, p1}, "DelegationTokenFactory.IssueGetChildrenRedirect got unexpected parameters")

		result := m.IssueGetChildrenRedirectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssueGetChildrenRedirect")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssueGetChildrenRedirectMock.mainExpectation != nil {

		input := m.IssueGetChildrenRedirectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssueGetChildrenRedirectInput{p, p1}, "DelegationTokenFactory.IssueGetChildrenRedirect got unexpected parameters")
		}

		result := m.IssueGetChildrenRedirectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssueGetChildrenRedirect")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssueGetChildrenRedirectFunc == nil {
		m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssueGetChildrenRedirect. %v %v", p, p1)
		return
	}

	return m.IssueGetChildrenRedirectFunc(p, p1)
}

//IssueGetChildrenRedirectMinimockCounter returns a count of DelegationTokenFactoryMock.IssueGetChildrenRedirectFunc invocations
func (m *DelegationTokenFactoryMock) IssueGetChildrenRedirectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IssueGetChildrenRedirectCounter)
}

//IssueGetChildrenRedirectMinimockPreCounter returns the value of DelegationTokenFactoryMock.IssueGetChildrenRedirect invocations
func (m *DelegationTokenFactoryMock) IssueGetChildrenRedirectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IssueGetChildrenRedirectPreCounter)
}

//IssueGetChildrenRedirectFinished returns true if mock invocations count is ok
func (m *DelegationTokenFactoryMock) IssueGetChildrenRedirectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IssueGetChildrenRedirectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IssueGetChildrenRedirectCounter) == uint64(len(m.IssueGetChildrenRedirectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IssueGetChildrenRedirectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IssueGetChildrenRedirectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IssueGetChildrenRedirectFunc != nil {
		return atomic.LoadUint64(&m.IssueGetChildrenRedirectCounter) > 0
	}

	return true
}

type mDelegationTokenFactoryMockIssueGetCodeRedirect struct {
	mock              *DelegationTokenFactoryMock
	mainExpectation   *DelegationTokenFactoryMockIssueGetCodeRedirectExpectation
	expectationSeries []*DelegationTokenFactoryMockIssueGetCodeRedirectExpectation
}

type DelegationTokenFactoryMockIssueGetCodeRedirectExpectation struct {
	input  *DelegationTokenFactoryMockIssueGetCodeRedirectInput
	result *DelegationTokenFactoryMockIssueGetCodeRedirectResult
}

type DelegationTokenFactoryMockIssueGetCodeRedirectInput struct {
	p  *core.RecordRef
	p1 core.Message
}

type DelegationTokenFactoryMockIssueGetCodeRedirectResult struct {
	r  core.DelegationToken
	r1 error
}

//Expect specifies that invocation of DelegationTokenFactory.IssueGetCodeRedirect is expected from 1 to Infinity times
func (m *mDelegationTokenFactoryMockIssueGetCodeRedirect) Expect(p *core.RecordRef, p1 core.Message) *mDelegationTokenFactoryMockIssueGetCodeRedirect {
	m.mock.IssueGetCodeRedirectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssueGetCodeRedirectExpectation{}
	}
	m.mainExpectation.input = &DelegationTokenFactoryMockIssueGetCodeRedirectInput{p, p1}
	return m
}

//Return specifies results of invocation of DelegationTokenFactory.IssueGetCodeRedirect
func (m *mDelegationTokenFactoryMockIssueGetCodeRedirect) Return(r core.DelegationToken, r1 error) *DelegationTokenFactoryMock {
	m.mock.IssueGetCodeRedirectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssueGetCodeRedirectExpectation{}
	}
	m.mainExpectation.result = &DelegationTokenFactoryMockIssueGetCodeRedirectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DelegationTokenFactory.IssueGetCodeRedirect is expected once
func (m *mDelegationTokenFactoryMockIssueGetCodeRedirect) ExpectOnce(p *core.RecordRef, p1 core.Message) *DelegationTokenFactoryMockIssueGetCodeRedirectExpectation {
	m.mock.IssueGetCodeRedirectFunc = nil
	m.mainExpectation = nil

	expectation := &DelegationTokenFactoryMockIssueGetCodeRedirectExpectation{}
	expectation.input = &DelegationTokenFactoryMockIssueGetCodeRedirectInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DelegationTokenFactoryMockIssueGetCodeRedirectExpectation) Return(r core.DelegationToken, r1 error) {
	e.result = &DelegationTokenFactoryMockIssueGetCodeRedirectResult{r, r1}
}

//Set uses given function f as a mock of DelegationTokenFactory.IssueGetCodeRedirect method
func (m *mDelegationTokenFactoryMockIssueGetCodeRedirect) Set(f func(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error)) *DelegationTokenFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IssueGetCodeRedirectFunc = f
	return m.mock
}

//IssueGetCodeRedirect implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) IssueGetCodeRedirect(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error) {
	counter := atomic.AddUint64(&m.IssueGetCodeRedirectPreCounter, 1)
	defer atomic.AddUint64(&m.IssueGetCodeRedirectCounter, 1)

	if len(m.IssueGetCodeRedirectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IssueGetCodeRedirectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssueGetCodeRedirect. %v %v", p, p1)
			return
		}

		input := m.IssueGetCodeRedirectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssueGetCodeRedirectInput{p, p1}, "DelegationTokenFactory.IssueGetCodeRedirect got unexpected parameters")

		result := m.IssueGetCodeRedirectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssueGetCodeRedirect")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssueGetCodeRedirectMock.mainExpectation != nil {

		input := m.IssueGetCodeRedirectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssueGetCodeRedirectInput{p, p1}, "DelegationTokenFactory.IssueGetCodeRedirect got unexpected parameters")
		}

		result := m.IssueGetCodeRedirectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssueGetCodeRedirect")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssueGetCodeRedirectFunc == nil {
		m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssueGetCodeRedirect. %v %v", p, p1)
		return
	}

	return m.IssueGetCodeRedirectFunc(p, p1)
}

//IssueGetCodeRedirectMinimockCounter returns a count of DelegationTokenFactoryMock.IssueGetCodeRedirectFunc invocations
func (m *DelegationTokenFactoryMock) IssueGetCodeRedirectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IssueGetCodeRedirectCounter)
}

//IssueGetCodeRedirectMinimockPreCounter returns the value of DelegationTokenFactoryMock.IssueGetCodeRedirect invocations
func (m *DelegationTokenFactoryMock) IssueGetCodeRedirectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IssueGetCodeRedirectPreCounter)
}

//IssueGetCodeRedirectFinished returns true if mock invocations count is ok
func (m *DelegationTokenFactoryMock) IssueGetCodeRedirectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IssueGetCodeRedirectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IssueGetCodeRedirectCounter) == uint64(len(m.IssueGetCodeRedirectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IssueGetCodeRedirectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IssueGetCodeRedirectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IssueGetCodeRedirectFunc != nil {
		return atomic.LoadUint64(&m.IssueGetCodeRedirectCounter) > 0
	}

	return true
}

type mDelegationTokenFactoryMockIssueGetObjectRedirect struct {
	mock              *DelegationTokenFactoryMock
	mainExpectation   *DelegationTokenFactoryMockIssueGetObjectRedirectExpectation
	expectationSeries []*DelegationTokenFactoryMockIssueGetObjectRedirectExpectation
}

type DelegationTokenFactoryMockIssueGetObjectRedirectExpectation struct {
	input  *DelegationTokenFactoryMockIssueGetObjectRedirectInput
	result *DelegationTokenFactoryMockIssueGetObjectRedirectResult
}

type DelegationTokenFactoryMockIssueGetObjectRedirectInput struct {
	p  *core.RecordRef
	p1 core.Message
}

type DelegationTokenFactoryMockIssueGetObjectRedirectResult struct {
	r  core.DelegationToken
	r1 error
}

//Expect specifies that invocation of DelegationTokenFactory.IssueGetObjectRedirect is expected from 1 to Infinity times
func (m *mDelegationTokenFactoryMockIssueGetObjectRedirect) Expect(p *core.RecordRef, p1 core.Message) *mDelegationTokenFactoryMockIssueGetObjectRedirect {
	m.mock.IssueGetObjectRedirectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssueGetObjectRedirectExpectation{}
	}
	m.mainExpectation.input = &DelegationTokenFactoryMockIssueGetObjectRedirectInput{p, p1}
	return m
}

//Return specifies results of invocation of DelegationTokenFactory.IssueGetObjectRedirect
func (m *mDelegationTokenFactoryMockIssueGetObjectRedirect) Return(r core.DelegationToken, r1 error) *DelegationTokenFactoryMock {
	m.mock.IssueGetObjectRedirectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssueGetObjectRedirectExpectation{}
	}
	m.mainExpectation.result = &DelegationTokenFactoryMockIssueGetObjectRedirectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DelegationTokenFactory.IssueGetObjectRedirect is expected once
func (m *mDelegationTokenFactoryMockIssueGetObjectRedirect) ExpectOnce(p *core.RecordRef, p1 core.Message) *DelegationTokenFactoryMockIssueGetObjectRedirectExpectation {
	m.mock.IssueGetObjectRedirectFunc = nil
	m.mainExpectation = nil

	expectation := &DelegationTokenFactoryMockIssueGetObjectRedirectExpectation{}
	expectation.input = &DelegationTokenFactoryMockIssueGetObjectRedirectInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DelegationTokenFactoryMockIssueGetObjectRedirectExpectation) Return(r core.DelegationToken, r1 error) {
	e.result = &DelegationTokenFactoryMockIssueGetObjectRedirectResult{r, r1}
}

//Set uses given function f as a mock of DelegationTokenFactory.IssueGetObjectRedirect method
func (m *mDelegationTokenFactoryMockIssueGetObjectRedirect) Set(f func(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error)) *DelegationTokenFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IssueGetObjectRedirectFunc = f
	return m.mock
}

//IssueGetObjectRedirect implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) IssueGetObjectRedirect(p *core.RecordRef, p1 core.Message) (r core.DelegationToken, r1 error) {
	counter := atomic.AddUint64(&m.IssueGetObjectRedirectPreCounter, 1)
	defer atomic.AddUint64(&m.IssueGetObjectRedirectCounter, 1)

	if len(m.IssueGetObjectRedirectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IssueGetObjectRedirectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssueGetObjectRedirect. %v %v", p, p1)
			return
		}

		input := m.IssueGetObjectRedirectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssueGetObjectRedirectInput{p, p1}, "DelegationTokenFactory.IssueGetObjectRedirect got unexpected parameters")

		result := m.IssueGetObjectRedirectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssueGetObjectRedirect")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssueGetObjectRedirectMock.mainExpectation != nil {

		input := m.IssueGetObjectRedirectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DelegationTokenFactoryMockIssueGetObjectRedirectInput{p, p1}, "DelegationTokenFactory.IssueGetObjectRedirect got unexpected parameters")
		}

		result := m.IssueGetObjectRedirectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DelegationTokenFactoryMock.IssueGetObjectRedirect")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IssueGetObjectRedirectFunc == nil {
		m.t.Fatalf("Unexpected call to DelegationTokenFactoryMock.IssueGetObjectRedirect. %v %v", p, p1)
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

//IssueGetObjectRedirectFinished returns true if mock invocations count is ok
func (m *DelegationTokenFactoryMock) IssueGetObjectRedirectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IssueGetObjectRedirectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) == uint64(len(m.IssueGetObjectRedirectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IssueGetObjectRedirectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IssueGetObjectRedirectFunc != nil {
		return atomic.LoadUint64(&m.IssueGetObjectRedirectCounter) > 0
	}

	return true
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
	p  core.Message
	p1 core.PulseNumber
}

type DelegationTokenFactoryMockIssuePendingExecutionResult struct {
	r  core.DelegationToken
	r1 error
}

//Expect specifies that invocation of DelegationTokenFactory.IssuePendingExecution is expected from 1 to Infinity times
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Expect(p core.Message, p1 core.PulseNumber) *mDelegationTokenFactoryMockIssuePendingExecution {
	m.mock.IssuePendingExecutionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssuePendingExecutionExpectation{}
	}
	m.mainExpectation.input = &DelegationTokenFactoryMockIssuePendingExecutionInput{p, p1}
	return m
}

//Return specifies results of invocation of DelegationTokenFactory.IssuePendingExecution
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Return(r core.DelegationToken, r1 error) *DelegationTokenFactoryMock {
	m.mock.IssuePendingExecutionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DelegationTokenFactoryMockIssuePendingExecutionExpectation{}
	}
	m.mainExpectation.result = &DelegationTokenFactoryMockIssuePendingExecutionResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DelegationTokenFactory.IssuePendingExecution is expected once
func (m *mDelegationTokenFactoryMockIssuePendingExecution) ExpectOnce(p core.Message, p1 core.PulseNumber) *DelegationTokenFactoryMockIssuePendingExecutionExpectation {
	m.mock.IssuePendingExecutionFunc = nil
	m.mainExpectation = nil

	expectation := &DelegationTokenFactoryMockIssuePendingExecutionExpectation{}
	expectation.input = &DelegationTokenFactoryMockIssuePendingExecutionInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DelegationTokenFactoryMockIssuePendingExecutionExpectation) Return(r core.DelegationToken, r1 error) {
	e.result = &DelegationTokenFactoryMockIssuePendingExecutionResult{r, r1}
}

//Set uses given function f as a mock of DelegationTokenFactory.IssuePendingExecution method
func (m *mDelegationTokenFactoryMockIssuePendingExecution) Set(f func(p core.Message, p1 core.PulseNumber) (r core.DelegationToken, r1 error)) *DelegationTokenFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IssuePendingExecutionFunc = f
	return m.mock
}

//IssuePendingExecution implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) IssuePendingExecution(p core.Message, p1 core.PulseNumber) (r core.DelegationToken, r1 error) {
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
	p core.Parcel
}

type DelegationTokenFactoryMockVerifyResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of DelegationTokenFactory.Verify is expected from 1 to Infinity times
func (m *mDelegationTokenFactoryMockVerify) Expect(p core.Parcel) *mDelegationTokenFactoryMockVerify {
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
func (m *mDelegationTokenFactoryMockVerify) ExpectOnce(p core.Parcel) *DelegationTokenFactoryMockVerifyExpectation {
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
func (m *mDelegationTokenFactoryMockVerify) Set(f func(p core.Parcel) (r bool, r1 error)) *DelegationTokenFactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VerifyFunc = f
	return m.mock
}

//Verify implements github.com/insolar/insolar/core.DelegationTokenFactory interface
func (m *DelegationTokenFactoryMock) Verify(p core.Parcel) (r bool, r1 error) {
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

	if !m.IssueGetChildrenRedirectFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetChildrenRedirect")
	}

	if !m.IssueGetCodeRedirectFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetCodeRedirect")
	}

	if !m.IssueGetObjectRedirectFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetObjectRedirect")
	}

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

	if !m.IssueGetChildrenRedirectFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetChildrenRedirect")
	}

	if !m.IssueGetCodeRedirectFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetCodeRedirect")
	}

	if !m.IssueGetObjectRedirectFinished() {
		m.t.Fatal("Expected call to DelegationTokenFactoryMock.IssueGetObjectRedirect")
	}

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
		ok = ok && m.IssueGetChildrenRedirectFinished()
		ok = ok && m.IssueGetCodeRedirectFinished()
		ok = ok && m.IssueGetObjectRedirectFinished()
		ok = ok && m.IssuePendingExecutionFinished()
		ok = ok && m.VerifyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IssueGetChildrenRedirectFinished() {
				m.t.Error("Expected call to DelegationTokenFactoryMock.IssueGetChildrenRedirect")
			}

			if !m.IssueGetCodeRedirectFinished() {
				m.t.Error("Expected call to DelegationTokenFactoryMock.IssueGetCodeRedirect")
			}

			if !m.IssueGetObjectRedirectFinished() {
				m.t.Error("Expected call to DelegationTokenFactoryMock.IssueGetObjectRedirect")
			}

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

	if !m.IssueGetChildrenRedirectFinished() {
		return false
	}

	if !m.IssueGetCodeRedirectFinished() {
		return false
	}

	if !m.IssueGetObjectRedirectFinished() {
		return false
	}

	if !m.IssuePendingExecutionFinished() {
		return false
	}

	if !m.VerifyFinished() {
		return false
	}

	return true
}
