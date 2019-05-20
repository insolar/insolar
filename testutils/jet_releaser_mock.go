package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetReleaser" can be found in github.com/insolar/insolar/ledger/light/hot
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetReleaserMock implements github.com/insolar/insolar/ledger/light/hot.JetReleaser
type JetReleaserMock struct {
	t minimock.Tester

	ThrowTimeoutFunc       func(p context.Context)
	ThrowTimeoutCounter    uint64
	ThrowTimeoutPreCounter uint64
	ThrowTimeoutMock       mJetReleaserMockThrowTimeout

	UnlockFunc       func(p context.Context, p1 insolar.ID) (r error)
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mJetReleaserMockUnlock
}

//NewJetReleaserMock returns a mock for github.com/insolar/insolar/ledger/light/hot.JetReleaser
func NewJetReleaserMock(t minimock.Tester) *JetReleaserMock {
	m := &JetReleaserMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ThrowTimeoutMock = mJetReleaserMockThrowTimeout{mock: m}
	m.UnlockMock = mJetReleaserMockUnlock{mock: m}

	return m
}

type mJetReleaserMockThrowTimeout struct {
	mock              *JetReleaserMock
	mainExpectation   *JetReleaserMockThrowTimeoutExpectation
	expectationSeries []*JetReleaserMockThrowTimeoutExpectation
}

type JetReleaserMockThrowTimeoutExpectation struct {
	input *JetReleaserMockThrowTimeoutInput
}

type JetReleaserMockThrowTimeoutInput struct {
	p context.Context
}

//Expect specifies that invocation of JetReleaser.ThrowTimeout is expected from 1 to Infinity times
func (m *mJetReleaserMockThrowTimeout) Expect(p context.Context) *mJetReleaserMockThrowTimeout {
	m.mock.ThrowTimeoutFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetReleaserMockThrowTimeoutExpectation{}
	}
	m.mainExpectation.input = &JetReleaserMockThrowTimeoutInput{p}
	return m
}

//Return specifies results of invocation of JetReleaser.ThrowTimeout
func (m *mJetReleaserMockThrowTimeout) Return() *JetReleaserMock {
	m.mock.ThrowTimeoutFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetReleaserMockThrowTimeoutExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of JetReleaser.ThrowTimeout is expected once
func (m *mJetReleaserMockThrowTimeout) ExpectOnce(p context.Context) *JetReleaserMockThrowTimeoutExpectation {
	m.mock.ThrowTimeoutFunc = nil
	m.mainExpectation = nil

	expectation := &JetReleaserMockThrowTimeoutExpectation{}
	expectation.input = &JetReleaserMockThrowTimeoutInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of JetReleaser.ThrowTimeout method
func (m *mJetReleaserMockThrowTimeout) Set(f func(p context.Context)) *JetReleaserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ThrowTimeoutFunc = f
	return m.mock
}

//ThrowTimeout implements github.com/insolar/insolar/ledger/light/hot.JetReleaser interface
func (m *JetReleaserMock) ThrowTimeout(p context.Context) {
	counter := atomic.AddUint64(&m.ThrowTimeoutPreCounter, 1)
	defer atomic.AddUint64(&m.ThrowTimeoutCounter, 1)

	if len(m.ThrowTimeoutMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ThrowTimeoutMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetReleaserMock.ThrowTimeout. %v", p)
			return
		}

		input := m.ThrowTimeoutMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetReleaserMockThrowTimeoutInput{p}, "JetReleaser.ThrowTimeout got unexpected parameters")

		return
	}

	if m.ThrowTimeoutMock.mainExpectation != nil {

		input := m.ThrowTimeoutMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetReleaserMockThrowTimeoutInput{p}, "JetReleaser.ThrowTimeout got unexpected parameters")
		}

		return
	}

	if m.ThrowTimeoutFunc == nil {
		m.t.Fatalf("Unexpected call to JetReleaserMock.ThrowTimeout. %v", p)
		return
	}

	m.ThrowTimeoutFunc(p)
}

//ThrowTimeoutMinimockCounter returns a count of JetReleaserMock.ThrowTimeoutFunc invocations
func (m *JetReleaserMock) ThrowTimeoutMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ThrowTimeoutCounter)
}

//ThrowTimeoutMinimockPreCounter returns the value of JetReleaserMock.ThrowTimeout invocations
func (m *JetReleaserMock) ThrowTimeoutMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ThrowTimeoutPreCounter)
}

//ThrowTimeoutFinished returns true if mock invocations count is ok
func (m *JetReleaserMock) ThrowTimeoutFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ThrowTimeoutMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ThrowTimeoutCounter) == uint64(len(m.ThrowTimeoutMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ThrowTimeoutMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ThrowTimeoutCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ThrowTimeoutFunc != nil {
		return atomic.LoadUint64(&m.ThrowTimeoutCounter) > 0
	}

	return true
}

type mJetReleaserMockUnlock struct {
	mock              *JetReleaserMock
	mainExpectation   *JetReleaserMockUnlockExpectation
	expectationSeries []*JetReleaserMockUnlockExpectation
}

type JetReleaserMockUnlockExpectation struct {
	input  *JetReleaserMockUnlockInput
	result *JetReleaserMockUnlockResult
}

type JetReleaserMockUnlockInput struct {
	p  context.Context
	p1 insolar.ID
}

type JetReleaserMockUnlockResult struct {
	r error
}

//Expect specifies that invocation of JetReleaser.Unlock is expected from 1 to Infinity times
func (m *mJetReleaserMockUnlock) Expect(p context.Context, p1 insolar.ID) *mJetReleaserMockUnlock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetReleaserMockUnlockExpectation{}
	}
	m.mainExpectation.input = &JetReleaserMockUnlockInput{p, p1}
	return m
}

//Return specifies results of invocation of JetReleaser.Unlock
func (m *mJetReleaserMockUnlock) Return(r error) *JetReleaserMock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetReleaserMockUnlockExpectation{}
	}
	m.mainExpectation.result = &JetReleaserMockUnlockResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetReleaser.Unlock is expected once
func (m *mJetReleaserMockUnlock) ExpectOnce(p context.Context, p1 insolar.ID) *JetReleaserMockUnlockExpectation {
	m.mock.UnlockFunc = nil
	m.mainExpectation = nil

	expectation := &JetReleaserMockUnlockExpectation{}
	expectation.input = &JetReleaserMockUnlockInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetReleaserMockUnlockExpectation) Return(r error) {
	e.result = &JetReleaserMockUnlockResult{r}
}

//Set uses given function f as a mock of JetReleaser.Unlock method
func (m *mJetReleaserMockUnlock) Set(f func(p context.Context, p1 insolar.ID) (r error)) *JetReleaserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnlockFunc = f
	return m.mock
}

//Unlock implements github.com/insolar/insolar/ledger/light/hot.JetReleaser interface
func (m *JetReleaserMock) Unlock(p context.Context, p1 insolar.ID) (r error) {
	counter := atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if len(m.UnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetReleaserMock.Unlock. %v %v", p, p1)
			return
		}

		input := m.UnlockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetReleaserMockUnlockInput{p, p1}, "JetReleaser.Unlock got unexpected parameters")

		result := m.UnlockMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetReleaserMock.Unlock")
			return
		}

		r = result.r

		return
	}

	if m.UnlockMock.mainExpectation != nil {

		input := m.UnlockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetReleaserMockUnlockInput{p, p1}, "JetReleaser.Unlock got unexpected parameters")
		}

		result := m.UnlockMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetReleaserMock.Unlock")
		}

		r = result.r

		return
	}

	if m.UnlockFunc == nil {
		m.t.Fatalf("Unexpected call to JetReleaserMock.Unlock. %v %v", p, p1)
		return
	}

	return m.UnlockFunc(p, p1)
}

//UnlockMinimockCounter returns a count of JetReleaserMock.UnlockFunc invocations
func (m *JetReleaserMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of JetReleaserMock.Unlock invocations
func (m *JetReleaserMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//UnlockFinished returns true if mock invocations count is ok
func (m *JetReleaserMock) UnlockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UnlockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UnlockCounter) == uint64(len(m.UnlockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UnlockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UnlockFunc != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetReleaserMock) ValidateCallCounters() {

	if !m.ThrowTimeoutFinished() {
		m.t.Fatal("Expected call to JetReleaserMock.ThrowTimeout")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to JetReleaserMock.Unlock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetReleaserMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetReleaserMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetReleaserMock) MinimockFinish() {

	if !m.ThrowTimeoutFinished() {
		m.t.Fatal("Expected call to JetReleaserMock.ThrowTimeout")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to JetReleaserMock.Unlock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetReleaserMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetReleaserMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ThrowTimeoutFinished()
		ok = ok && m.UnlockFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ThrowTimeoutFinished() {
				m.t.Error("Expected call to JetReleaserMock.ThrowTimeout")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to JetReleaserMock.Unlock")
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
func (m *JetReleaserMock) AllMocksCalled() bool {

	if !m.ThrowTimeoutFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	return true
}
