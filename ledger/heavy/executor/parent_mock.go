package replica

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Parent" can be found in github.com/insolar/insolar/ledger/heavy/replica
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//ParentMock implements github.com/insolar/insolar/ledger/heavy/replica.Parent
type ParentMock struct {
	t minimock.Tester

	PullFunc       func(p context.Context, p1 Page) (r []byte, r1 uint32, r2 error)
	PullCounter    uint64
	PullPreCounter uint64
	PullMock       mParentMockPull

	SubscribeFunc       func(p context.Context, p1 Target, p2 Page) (r error)
	SubscribeCounter    uint64
	SubscribePreCounter uint64
	SubscribeMock       mParentMockSubscribe
}

//NewParentMock returns a mock for github.com/insolar/insolar/ledger/heavy/replica.Parent
func NewParentMock(t minimock.Tester) *ParentMock {
	m := &ParentMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PullMock = mParentMockPull{mock: m}
	m.SubscribeMock = mParentMockSubscribe{mock: m}

	return m
}

type mParentMockPull struct {
	mock              *ParentMock
	mainExpectation   *ParentMockPullExpectation
	expectationSeries []*ParentMockPullExpectation
}

type ParentMockPullExpectation struct {
	input  *ParentMockPullInput
	result *ParentMockPullResult
}

type ParentMockPullInput struct {
	p  context.Context
	p1 Page
}

type ParentMockPullResult struct {
	r  []byte
	r1 uint32
	r2 error
}

//Expect specifies that invocation of Parent.Pull is expected from 1 to Infinity times
func (m *mParentMockPull) Expect(p context.Context, p1 Page) *mParentMockPull {
	m.mock.PullFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParentMockPullExpectation{}
	}
	m.mainExpectation.input = &ParentMockPullInput{p, p1}
	return m
}

//Return specifies results of invocation of Parent.Pull
func (m *mParentMockPull) Return(r []byte, r1 uint32, r2 error) *ParentMock {
	m.mock.PullFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParentMockPullExpectation{}
	}
	m.mainExpectation.result = &ParentMockPullResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of Parent.Pull is expected once
func (m *mParentMockPull) ExpectOnce(p context.Context, p1 Page) *ParentMockPullExpectation {
	m.mock.PullFunc = nil
	m.mainExpectation = nil

	expectation := &ParentMockPullExpectation{}
	expectation.input = &ParentMockPullInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParentMockPullExpectation) Return(r []byte, r1 uint32, r2 error) {
	e.result = &ParentMockPullResult{r, r1, r2}
}

//Set uses given function f as a mock of Parent.Pull method
func (m *mParentMockPull) Set(f func(p context.Context, p1 Page) (r []byte, r1 uint32, r2 error)) *ParentMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PullFunc = f
	return m.mock
}

//Pull implements github.com/insolar/insolar/ledger/heavy/replica.Parent interface
func (m *ParentMock) Pull(p context.Context, p1 Page) (r []byte, r1 uint32, r2 error) {
	counter := atomic.AddUint64(&m.PullPreCounter, 1)
	defer atomic.AddUint64(&m.PullCounter, 1)

	if len(m.PullMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PullMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParentMock.Pull. %v %v", p, p1)
			return
		}

		input := m.PullMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ParentMockPullInput{p, p1}, "Parent.Pull got unexpected parameters")

		result := m.PullMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParentMock.Pull")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.PullMock.mainExpectation != nil {

		input := m.PullMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ParentMockPullInput{p, p1}, "Parent.Pull got unexpected parameters")
		}

		result := m.PullMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParentMock.Pull")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.PullFunc == nil {
		m.t.Fatalf("Unexpected call to ParentMock.Pull. %v %v", p, p1)
		return
	}

	return m.PullFunc(p, p1)
}

//PullMinimockCounter returns a count of ParentMock.PullFunc invocations
func (m *ParentMock) PullMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PullCounter)
}

//PullMinimockPreCounter returns the value of ParentMock.Pull invocations
func (m *ParentMock) PullMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PullPreCounter)
}

//PullFinished returns true if mock invocations count is ok
func (m *ParentMock) PullFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PullMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PullCounter) == uint64(len(m.PullMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PullMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PullCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PullFunc != nil {
		return atomic.LoadUint64(&m.PullCounter) > 0
	}

	return true
}

type mParentMockSubscribe struct {
	mock              *ParentMock
	mainExpectation   *ParentMockSubscribeExpectation
	expectationSeries []*ParentMockSubscribeExpectation
}

type ParentMockSubscribeExpectation struct {
	input  *ParentMockSubscribeInput
	result *ParentMockSubscribeResult
}

type ParentMockSubscribeInput struct {
	p  context.Context
	p1 Target
	p2 Page
}

type ParentMockSubscribeResult struct {
	r error
}

//Expect specifies that invocation of Parent.Subscribe is expected from 1 to Infinity times
func (m *mParentMockSubscribe) Expect(p context.Context, p1 Target, p2 Page) *mParentMockSubscribe {
	m.mock.SubscribeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParentMockSubscribeExpectation{}
	}
	m.mainExpectation.input = &ParentMockSubscribeInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Parent.Subscribe
func (m *mParentMockSubscribe) Return(r error) *ParentMock {
	m.mock.SubscribeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParentMockSubscribeExpectation{}
	}
	m.mainExpectation.result = &ParentMockSubscribeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parent.Subscribe is expected once
func (m *mParentMockSubscribe) ExpectOnce(p context.Context, p1 Target, p2 Page) *ParentMockSubscribeExpectation {
	m.mock.SubscribeFunc = nil
	m.mainExpectation = nil

	expectation := &ParentMockSubscribeExpectation{}
	expectation.input = &ParentMockSubscribeInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParentMockSubscribeExpectation) Return(r error) {
	e.result = &ParentMockSubscribeResult{r}
}

//Set uses given function f as a mock of Parent.Subscribe method
func (m *mParentMockSubscribe) Set(f func(p context.Context, p1 Target, p2 Page) (r error)) *ParentMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SubscribeFunc = f
	return m.mock
}

//Subscribe implements github.com/insolar/insolar/ledger/heavy/replica.Parent interface
func (m *ParentMock) Subscribe(p context.Context, p1 Target, p2 Page) (r error) {
	counter := atomic.AddUint64(&m.SubscribePreCounter, 1)
	defer atomic.AddUint64(&m.SubscribeCounter, 1)

	if len(m.SubscribeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SubscribeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParentMock.Subscribe. %v %v %v", p, p1, p2)
			return
		}

		input := m.SubscribeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ParentMockSubscribeInput{p, p1, p2}, "Parent.Subscribe got unexpected parameters")

		result := m.SubscribeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParentMock.Subscribe")
			return
		}

		r = result.r

		return
	}

	if m.SubscribeMock.mainExpectation != nil {

		input := m.SubscribeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ParentMockSubscribeInput{p, p1, p2}, "Parent.Subscribe got unexpected parameters")
		}

		result := m.SubscribeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParentMock.Subscribe")
		}

		r = result.r

		return
	}

	if m.SubscribeFunc == nil {
		m.t.Fatalf("Unexpected call to ParentMock.Subscribe. %v %v %v", p, p1, p2)
		return
	}

	return m.SubscribeFunc(p, p1, p2)
}

//SubscribeMinimockCounter returns a count of ParentMock.SubscribeFunc invocations
func (m *ParentMock) SubscribeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SubscribeCounter)
}

//SubscribeMinimockPreCounter returns the value of ParentMock.Subscribe invocations
func (m *ParentMock) SubscribeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SubscribePreCounter)
}

//SubscribeFinished returns true if mock invocations count is ok
func (m *ParentMock) SubscribeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SubscribeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SubscribeCounter) == uint64(len(m.SubscribeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SubscribeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SubscribeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SubscribeFunc != nil {
		return atomic.LoadUint64(&m.SubscribeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ParentMock) ValidateCallCounters() {

	if !m.PullFinished() {
		m.t.Fatal("Expected call to ParentMock.Pull")
	}

	if !m.SubscribeFinished() {
		m.t.Fatal("Expected call to ParentMock.Subscribe")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ParentMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ParentMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ParentMock) MinimockFinish() {

	if !m.PullFinished() {
		m.t.Fatal("Expected call to ParentMock.Pull")
	}

	if !m.SubscribeFinished() {
		m.t.Fatal("Expected call to ParentMock.Subscribe")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ParentMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ParentMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.PullFinished()
		ok = ok && m.SubscribeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.PullFinished() {
				m.t.Error("Expected call to ParentMock.Pull")
			}

			if !m.SubscribeFinished() {
				m.t.Error("Expected call to ParentMock.Subscribe")
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
func (m *ParentMock) AllMocksCalled() bool {

	if !m.PullFinished() {
		return false
	}

	if !m.SubscribeFinished() {
		return false
	}

	return true
}
