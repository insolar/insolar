package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetFetcher" can be found in github.com/insolar/insolar/ledger/light/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetFetcherMock implements github.com/insolar/insolar/ledger/light/executor.JetFetcher
type JetFetcherMock struct {
	t minimock.Tester

	FetchFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error)
	FetchCounter    uint64
	FetchPreCounter uint64
	FetchMock       mJetFetcherMockFetch

	ReleaseFunc       func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber)
	ReleaseCounter    uint64
	ReleasePreCounter uint64
	ReleaseMock       mJetFetcherMockRelease
}

//NewJetFetcherMock returns a mock for github.com/insolar/insolar/ledger/light/executor.JetFetcher
func NewJetFetcherMock(t minimock.Tester) *JetFetcherMock {
	m := &JetFetcherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FetchMock = mJetFetcherMockFetch{mock: m}
	m.ReleaseMock = mJetFetcherMockRelease{mock: m}

	return m
}

type mJetFetcherMockFetch struct {
	mock              *JetFetcherMock
	mainExpectation   *JetFetcherMockFetchExpectation
	expectationSeries []*JetFetcherMockFetchExpectation
}

type JetFetcherMockFetchExpectation struct {
	input  *JetFetcherMockFetchInput
	result *JetFetcherMockFetchResult
}

type JetFetcherMockFetchInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type JetFetcherMockFetchResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of JetFetcher.Fetch is expected from 1 to Infinity times
func (m *mJetFetcherMockFetch) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mJetFetcherMockFetch {
	m.mock.FetchFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetFetcherMockFetchExpectation{}
	}
	m.mainExpectation.input = &JetFetcherMockFetchInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of JetFetcher.Fetch
func (m *mJetFetcherMockFetch) Return(r *insolar.ID, r1 error) *JetFetcherMock {
	m.mock.FetchFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetFetcherMockFetchExpectation{}
	}
	m.mainExpectation.result = &JetFetcherMockFetchResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of JetFetcher.Fetch is expected once
func (m *mJetFetcherMockFetch) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *JetFetcherMockFetchExpectation {
	m.mock.FetchFunc = nil
	m.mainExpectation = nil

	expectation := &JetFetcherMockFetchExpectation{}
	expectation.input = &JetFetcherMockFetchInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetFetcherMockFetchExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &JetFetcherMockFetchResult{r, r1}
}

//Set uses given function f as a mock of JetFetcher.Fetch method
func (m *mJetFetcherMockFetch) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error)) *JetFetcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FetchFunc = f
	return m.mock
}

//Fetch implements github.com/insolar/insolar/ledger/light/executor.JetFetcher interface
func (m *JetFetcherMock) Fetch(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.FetchPreCounter, 1)
	defer atomic.AddUint64(&m.FetchCounter, 1)

	if len(m.FetchMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FetchMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetFetcherMock.Fetch. %v %v %v", p, p1, p2)
			return
		}

		input := m.FetchMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetFetcherMockFetchInput{p, p1, p2}, "JetFetcher.Fetch got unexpected parameters")

		result := m.FetchMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetFetcherMock.Fetch")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FetchMock.mainExpectation != nil {

		input := m.FetchMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetFetcherMockFetchInput{p, p1, p2}, "JetFetcher.Fetch got unexpected parameters")
		}

		result := m.FetchMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetFetcherMock.Fetch")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FetchFunc == nil {
		m.t.Fatalf("Unexpected call to JetFetcherMock.Fetch. %v %v %v", p, p1, p2)
		return
	}

	return m.FetchFunc(p, p1, p2)
}

//FetchMinimockCounter returns a count of JetFetcherMock.FetchFunc invocations
func (m *JetFetcherMock) FetchMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FetchCounter)
}

//FetchMinimockPreCounter returns the value of JetFetcherMock.Fetch invocations
func (m *JetFetcherMock) FetchMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FetchPreCounter)
}

//FetchFinished returns true if mock invocations count is ok
func (m *JetFetcherMock) FetchFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FetchMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FetchCounter) == uint64(len(m.FetchMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FetchMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FetchCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FetchFunc != nil {
		return atomic.LoadUint64(&m.FetchCounter) > 0
	}

	return true
}

type mJetFetcherMockRelease struct {
	mock              *JetFetcherMock
	mainExpectation   *JetFetcherMockReleaseExpectation
	expectationSeries []*JetFetcherMockReleaseExpectation
}

type JetFetcherMockReleaseExpectation struct {
	input *JetFetcherMockReleaseInput
}

type JetFetcherMockReleaseInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 insolar.PulseNumber
}

//Expect specifies that invocation of JetFetcher.Release is expected from 1 to Infinity times
func (m *mJetFetcherMockRelease) Expect(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *mJetFetcherMockRelease {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetFetcherMockReleaseExpectation{}
	}
	m.mainExpectation.input = &JetFetcherMockReleaseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of JetFetcher.Release
func (m *mJetFetcherMockRelease) Return() *JetFetcherMock {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetFetcherMockReleaseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of JetFetcher.Release is expected once
func (m *mJetFetcherMockRelease) ExpectOnce(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *JetFetcherMockReleaseExpectation {
	m.mock.ReleaseFunc = nil
	m.mainExpectation = nil

	expectation := &JetFetcherMockReleaseExpectation{}
	expectation.input = &JetFetcherMockReleaseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of JetFetcher.Release method
func (m *mJetFetcherMockRelease) Set(f func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber)) *JetFetcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReleaseFunc = f
	return m.mock
}

//Release implements github.com/insolar/insolar/ledger/light/executor.JetFetcher interface
func (m *JetFetcherMock) Release(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.ReleasePreCounter, 1)
	defer atomic.AddUint64(&m.ReleaseCounter, 1)

	if len(m.ReleaseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReleaseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetFetcherMock.Release. %v %v %v", p, p1, p2)
			return
		}

		input := m.ReleaseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetFetcherMockReleaseInput{p, p1, p2}, "JetFetcher.Release got unexpected parameters")

		return
	}

	if m.ReleaseMock.mainExpectation != nil {

		input := m.ReleaseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetFetcherMockReleaseInput{p, p1, p2}, "JetFetcher.Release got unexpected parameters")
		}

		return
	}

	if m.ReleaseFunc == nil {
		m.t.Fatalf("Unexpected call to JetFetcherMock.Release. %v %v %v", p, p1, p2)
		return
	}

	m.ReleaseFunc(p, p1, p2)
}

//ReleaseMinimockCounter returns a count of JetFetcherMock.ReleaseFunc invocations
func (m *JetFetcherMock) ReleaseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseCounter)
}

//ReleaseMinimockPreCounter returns the value of JetFetcherMock.Release invocations
func (m *JetFetcherMock) ReleaseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReleasePreCounter)
}

//ReleaseFinished returns true if mock invocations count is ok
func (m *JetFetcherMock) ReleaseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReleaseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReleaseCounter) == uint64(len(m.ReleaseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReleaseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReleaseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReleaseFunc != nil {
		return atomic.LoadUint64(&m.ReleaseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetFetcherMock) ValidateCallCounters() {

	if !m.FetchFinished() {
		m.t.Fatal("Expected call to JetFetcherMock.Fetch")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to JetFetcherMock.Release")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetFetcherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetFetcherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetFetcherMock) MinimockFinish() {

	if !m.FetchFinished() {
		m.t.Fatal("Expected call to JetFetcherMock.Fetch")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to JetFetcherMock.Release")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetFetcherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetFetcherMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FetchFinished()
		ok = ok && m.ReleaseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FetchFinished() {
				m.t.Error("Expected call to JetFetcherMock.Fetch")
			}

			if !m.ReleaseFinished() {
				m.t.Error("Expected call to JetFetcherMock.Release")
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
func (m *JetFetcherMock) AllMocksCalled() bool {

	if !m.FetchFinished() {
		return false
	}

	if !m.ReleaseFinished() {
		return false
	}

	return true
}
