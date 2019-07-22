package census

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "OfflinePopulation" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/census
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	endpoints "github.com/insolar/insolar/network/consensus/common/endpoints"
	profiles "github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"

	testify_assert "github.com/stretchr/testify/assert"
)

//OfflinePopulationMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.OfflinePopulation
type OfflinePopulationMock struct {
	t minimock.Tester

	FindRegisteredProfileFunc       func(p endpoints.Inbound) (r profiles.Host)
	FindRegisteredProfileCounter    uint64
	FindRegisteredProfilePreCounter uint64
	FindRegisteredProfileMock       mOfflinePopulationMockFindRegisteredProfile
}

//NewOfflinePopulationMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/census.OfflinePopulation
func NewOfflinePopulationMock(t minimock.Tester) *OfflinePopulationMock {
	m := &OfflinePopulationMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FindRegisteredProfileMock = mOfflinePopulationMockFindRegisteredProfile{mock: m}

	return m
}

type mOfflinePopulationMockFindRegisteredProfile struct {
	mock              *OfflinePopulationMock
	mainExpectation   *OfflinePopulationMockFindRegisteredProfileExpectation
	expectationSeries []*OfflinePopulationMockFindRegisteredProfileExpectation
}

type OfflinePopulationMockFindRegisteredProfileExpectation struct {
	input  *OfflinePopulationMockFindRegisteredProfileInput
	result *OfflinePopulationMockFindRegisteredProfileResult
}

type OfflinePopulationMockFindRegisteredProfileInput struct {
	p endpoints.Inbound
}

type OfflinePopulationMockFindRegisteredProfileResult struct {
	r profiles.Host
}

//Expect specifies that invocation of OfflinePopulation.FindRegisteredProfile is expected from 1 to Infinity times
func (m *mOfflinePopulationMockFindRegisteredProfile) Expect(p endpoints.Inbound) *mOfflinePopulationMockFindRegisteredProfile {
	m.mock.FindRegisteredProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OfflinePopulationMockFindRegisteredProfileExpectation{}
	}
	m.mainExpectation.input = &OfflinePopulationMockFindRegisteredProfileInput{p}
	return m
}

//Return specifies results of invocation of OfflinePopulation.FindRegisteredProfile
func (m *mOfflinePopulationMockFindRegisteredProfile) Return(r profiles.Host) *OfflinePopulationMock {
	m.mock.FindRegisteredProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OfflinePopulationMockFindRegisteredProfileExpectation{}
	}
	m.mainExpectation.result = &OfflinePopulationMockFindRegisteredProfileResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of OfflinePopulation.FindRegisteredProfile is expected once
func (m *mOfflinePopulationMockFindRegisteredProfile) ExpectOnce(p endpoints.Inbound) *OfflinePopulationMockFindRegisteredProfileExpectation {
	m.mock.FindRegisteredProfileFunc = nil
	m.mainExpectation = nil

	expectation := &OfflinePopulationMockFindRegisteredProfileExpectation{}
	expectation.input = &OfflinePopulationMockFindRegisteredProfileInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OfflinePopulationMockFindRegisteredProfileExpectation) Return(r profiles.Host) {
	e.result = &OfflinePopulationMockFindRegisteredProfileResult{r}
}

//Set uses given function f as a mock of OfflinePopulation.FindRegisteredProfile method
func (m *mOfflinePopulationMockFindRegisteredProfile) Set(f func(p endpoints.Inbound) (r profiles.Host)) *OfflinePopulationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FindRegisteredProfileFunc = f
	return m.mock
}

//FindRegisteredProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.OfflinePopulation interface
func (m *OfflinePopulationMock) FindRegisteredProfile(p endpoints.Inbound) (r profiles.Host) {
	counter := atomic.AddUint64(&m.FindRegisteredProfilePreCounter, 1)
	defer atomic.AddUint64(&m.FindRegisteredProfileCounter, 1)

	if len(m.FindRegisteredProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FindRegisteredProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OfflinePopulationMock.FindRegisteredProfile. %v", p)
			return
		}

		input := m.FindRegisteredProfileMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, OfflinePopulationMockFindRegisteredProfileInput{p}, "OfflinePopulation.FindRegisteredProfile got unexpected parameters")

		result := m.FindRegisteredProfileMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OfflinePopulationMock.FindRegisteredProfile")
			return
		}

		r = result.r

		return
	}

	if m.FindRegisteredProfileMock.mainExpectation != nil {

		input := m.FindRegisteredProfileMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, OfflinePopulationMockFindRegisteredProfileInput{p}, "OfflinePopulation.FindRegisteredProfile got unexpected parameters")
		}

		result := m.FindRegisteredProfileMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OfflinePopulationMock.FindRegisteredProfile")
		}

		r = result.r

		return
	}

	if m.FindRegisteredProfileFunc == nil {
		m.t.Fatalf("Unexpected call to OfflinePopulationMock.FindRegisteredProfile. %v", p)
		return
	}

	return m.FindRegisteredProfileFunc(p)
}

//FindRegisteredProfileMinimockCounter returns a count of OfflinePopulationMock.FindRegisteredProfileFunc invocations
func (m *OfflinePopulationMock) FindRegisteredProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FindRegisteredProfileCounter)
}

//FindRegisteredProfileMinimockPreCounter returns the value of OfflinePopulationMock.FindRegisteredProfile invocations
func (m *OfflinePopulationMock) FindRegisteredProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FindRegisteredProfilePreCounter)
}

//FindRegisteredProfileFinished returns true if mock invocations count is ok
func (m *OfflinePopulationMock) FindRegisteredProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FindRegisteredProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FindRegisteredProfileCounter) == uint64(len(m.FindRegisteredProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FindRegisteredProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FindRegisteredProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FindRegisteredProfileFunc != nil {
		return atomic.LoadUint64(&m.FindRegisteredProfileCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *OfflinePopulationMock) ValidateCallCounters() {

	if !m.FindRegisteredProfileFinished() {
		m.t.Fatal("Expected call to OfflinePopulationMock.FindRegisteredProfile")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *OfflinePopulationMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *OfflinePopulationMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *OfflinePopulationMock) MinimockFinish() {

	if !m.FindRegisteredProfileFinished() {
		m.t.Fatal("Expected call to OfflinePopulationMock.FindRegisteredProfile")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *OfflinePopulationMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *OfflinePopulationMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FindRegisteredProfileFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FindRegisteredProfileFinished() {
				m.t.Error("Expected call to OfflinePopulationMock.FindRegisteredProfile")
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
func (m *OfflinePopulationMock) AllMocksCalled() bool {

	if !m.FindRegisteredProfileFinished() {
		return false
	}

	return true
}
