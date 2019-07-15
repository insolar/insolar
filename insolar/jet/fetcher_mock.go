package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Fetcher" can be found in github.com/insolar/insolar/insolar/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//FetcherMock implements github.com/insolar/insolar/insolar/jet.Fetcher
type FetcherMock struct {
	t minimock.Tester

	FetchFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error)
	FetchCounter    uint64
	FetchPreCounter uint64
	FetchMock       mFetcherMockFetch

	ReleaseFunc       func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber)
	ReleaseCounter    uint64
	ReleasePreCounter uint64
	ReleaseMock       mFetcherMockRelease
}

//NewFetcherMock returns a mock for github.com/insolar/insolar/insolar/jet.Fetcher
func NewFetcherMock(t minimock.Tester) *FetcherMock {
	m := &FetcherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FetchMock = mFetcherMockFetch{mock: m}
	m.ReleaseMock = mFetcherMockRelease{mock: m}

	return m
}

type mFetcherMockFetch struct {
	mock              *FetcherMock
	mainExpectation   *FetcherMockFetchExpectation
	expectationSeries []*FetcherMockFetchExpectation
}

type FetcherMockFetchExpectation struct {
	input  *FetcherMockFetchInput
	result *FetcherMockFetchResult
}

type FetcherMockFetchInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type FetcherMockFetchResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Fetcher.Fetch is expected from 1 to Infinity times
func (m *mFetcherMockFetch) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mFetcherMockFetch {
	m.mock.FetchFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FetcherMockFetchExpectation{}
	}
	m.mainExpectation.input = &FetcherMockFetchInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Fetcher.Fetch
func (m *mFetcherMockFetch) Return(r *insolar.ID, r1 error) *FetcherMock {
	m.mock.FetchFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FetcherMockFetchExpectation{}
	}
	m.mainExpectation.result = &FetcherMockFetchResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Fetcher.Fetch is expected once
func (m *mFetcherMockFetch) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *FetcherMockFetchExpectation {
	m.mock.FetchFunc = nil
	m.mainExpectation = nil

	expectation := &FetcherMockFetchExpectation{}
	expectation.input = &FetcherMockFetchInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FetcherMockFetchExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &FetcherMockFetchResult{r, r1}
}

//Set uses given function f as a mock of Fetcher.Fetch method
func (m *mFetcherMockFetch) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error)) *FetcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FetchFunc = f
	return m.mock
}

//Fetch implements github.com/insolar/insolar/insolar/jet.Fetcher interface
func (m *FetcherMock) Fetch(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.FetchPreCounter, 1)
	defer atomic.AddUint64(&m.FetchCounter, 1)

	if len(m.FetchMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FetchMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FetcherMock.Fetch. %v %v %v", p, p1, p2)
			return
		}

		input := m.FetchMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FetcherMockFetchInput{p, p1, p2}, "Fetcher.Fetch got unexpected parameters")

		result := m.FetchMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FetcherMock.Fetch")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FetchMock.mainExpectation != nil {

		input := m.FetchMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FetcherMockFetchInput{p, p1, p2}, "Fetcher.Fetch got unexpected parameters")
		}

		result := m.FetchMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FetcherMock.Fetch")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FetchFunc == nil {
		m.t.Fatalf("Unexpected call to FetcherMock.Fetch. %v %v %v", p, p1, p2)
		return
	}

	return m.FetchFunc(p, p1, p2)
}

//FetchMinimockCounter returns a count of FetcherMock.FetchFunc invocations
func (m *FetcherMock) FetchMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FetchCounter)
}

//FetchMinimockPreCounter returns the value of FetcherMock.Fetch invocations
func (m *FetcherMock) FetchMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FetchPreCounter)
}

//FetchFinished returns true if mock invocations count is ok
func (m *FetcherMock) FetchFinished() bool {
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

type mFetcherMockRelease struct {
	mock              *FetcherMock
	mainExpectation   *FetcherMockReleaseExpectation
	expectationSeries []*FetcherMockReleaseExpectation
}

type FetcherMockReleaseExpectation struct {
	input *FetcherMockReleaseInput
}

type FetcherMockReleaseInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 insolar.PulseNumber
}

//Expect specifies that invocation of Fetcher.Release is expected from 1 to Infinity times
func (m *mFetcherMockRelease) Expect(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *mFetcherMockRelease {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FetcherMockReleaseExpectation{}
	}
	m.mainExpectation.input = &FetcherMockReleaseInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Fetcher.Release
func (m *mFetcherMockRelease) Return() *FetcherMock {
	m.mock.ReleaseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FetcherMockReleaseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Fetcher.Release is expected once
func (m *mFetcherMockRelease) ExpectOnce(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *FetcherMockReleaseExpectation {
	m.mock.ReleaseFunc = nil
	m.mainExpectation = nil

	expectation := &FetcherMockReleaseExpectation{}
	expectation.input = &FetcherMockReleaseInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Fetcher.Release method
func (m *mFetcherMockRelease) Set(f func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber)) *FetcherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReleaseFunc = f
	return m.mock
}

//Release implements github.com/insolar/insolar/insolar/jet.Fetcher interface
func (m *FetcherMock) Release(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.ReleasePreCounter, 1)
	defer atomic.AddUint64(&m.ReleaseCounter, 1)

	if len(m.ReleaseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReleaseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FetcherMock.Release. %v %v %v", p, p1, p2)
			return
		}

		input := m.ReleaseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FetcherMockReleaseInput{p, p1, p2}, "Fetcher.Release got unexpected parameters")

		return
	}

	if m.ReleaseMock.mainExpectation != nil {

		input := m.ReleaseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FetcherMockReleaseInput{p, p1, p2}, "Fetcher.Release got unexpected parameters")
		}

		return
	}

	if m.ReleaseFunc == nil {
		m.t.Fatalf("Unexpected call to FetcherMock.Release. %v %v %v", p, p1, p2)
		return
	}

	m.ReleaseFunc(p, p1, p2)
}

//ReleaseMinimockCounter returns a count of FetcherMock.ReleaseFunc invocations
func (m *FetcherMock) ReleaseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseCounter)
}

//ReleaseMinimockPreCounter returns the value of FetcherMock.Release invocations
func (m *FetcherMock) ReleaseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReleasePreCounter)
}

//ReleaseFinished returns true if mock invocations count is ok
func (m *FetcherMock) ReleaseFinished() bool {
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
func (m *FetcherMock) ValidateCallCounters() {

	if !m.FetchFinished() {
		m.t.Fatal("Expected call to FetcherMock.Fetch")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to FetcherMock.Release")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FetcherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FetcherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FetcherMock) MinimockFinish() {

	if !m.FetchFinished() {
		m.t.Fatal("Expected call to FetcherMock.Fetch")
	}

	if !m.ReleaseFinished() {
		m.t.Fatal("Expected call to FetcherMock.Release")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FetcherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FetcherMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to FetcherMock.Fetch")
			}

			if !m.ReleaseFinished() {
				m.t.Error("Expected call to FetcherMock.Release")
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
func (m *FetcherMock) AllMocksCalled() bool {

	if !m.FetchFinished() {
		return false
	}

	if !m.ReleaseFinished() {
		return false
	}

	return true
}
