package census

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "EvictedPopulation" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/census
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	profiles "github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"

	testify_assert "github.com/stretchr/testify/assert"
)

//EvictedPopulationMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.EvictedPopulation
type EvictedPopulationMock struct {
	t minimock.Tester

	FindProfileFunc       func(p insolar.ShortNodeID) (r profiles.EvictedNode)
	FindProfileCounter    uint64
	FindProfilePreCounter uint64
	FindProfileMock       mEvictedPopulationMockFindProfile

	GetCountFunc       func() (r int)
	GetCountCounter    uint64
	GetCountPreCounter uint64
	GetCountMock       mEvictedPopulationMockGetCount

	GetProfilesFunc       func() (r []profiles.EvictedNode)
	GetProfilesCounter    uint64
	GetProfilesPreCounter uint64
	GetProfilesMock       mEvictedPopulationMockGetProfiles
}

//NewEvictedPopulationMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/census.EvictedPopulation
func NewEvictedPopulationMock(t minimock.Tester) *EvictedPopulationMock {
	m := &EvictedPopulationMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FindProfileMock = mEvictedPopulationMockFindProfile{mock: m}
	m.GetCountMock = mEvictedPopulationMockGetCount{mock: m}
	m.GetProfilesMock = mEvictedPopulationMockGetProfiles{mock: m}

	return m
}

type mEvictedPopulationMockFindProfile struct {
	mock              *EvictedPopulationMock
	mainExpectation   *EvictedPopulationMockFindProfileExpectation
	expectationSeries []*EvictedPopulationMockFindProfileExpectation
}

type EvictedPopulationMockFindProfileExpectation struct {
	input  *EvictedPopulationMockFindProfileInput
	result *EvictedPopulationMockFindProfileResult
}

type EvictedPopulationMockFindProfileInput struct {
	p insolar.ShortNodeID
}

type EvictedPopulationMockFindProfileResult struct {
	r profiles.EvictedNode
}

//Expect specifies that invocation of EvictedPopulation.FindProfile is expected from 1 to Infinity times
func (m *mEvictedPopulationMockFindProfile) Expect(p insolar.ShortNodeID) *mEvictedPopulationMockFindProfile {
	m.mock.FindProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &EvictedPopulationMockFindProfileExpectation{}
	}
	m.mainExpectation.input = &EvictedPopulationMockFindProfileInput{p}
	return m
}

//Return specifies results of invocation of EvictedPopulation.FindProfile
func (m *mEvictedPopulationMockFindProfile) Return(r profiles.EvictedNode) *EvictedPopulationMock {
	m.mock.FindProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &EvictedPopulationMockFindProfileExpectation{}
	}
	m.mainExpectation.result = &EvictedPopulationMockFindProfileResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of EvictedPopulation.FindProfile is expected once
func (m *mEvictedPopulationMockFindProfile) ExpectOnce(p insolar.ShortNodeID) *EvictedPopulationMockFindProfileExpectation {
	m.mock.FindProfileFunc = nil
	m.mainExpectation = nil

	expectation := &EvictedPopulationMockFindProfileExpectation{}
	expectation.input = &EvictedPopulationMockFindProfileInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *EvictedPopulationMockFindProfileExpectation) Return(r profiles.EvictedNode) {
	e.result = &EvictedPopulationMockFindProfileResult{r}
}

//Set uses given function f as a mock of EvictedPopulation.FindProfile method
func (m *mEvictedPopulationMockFindProfile) Set(f func(p insolar.ShortNodeID) (r profiles.EvictedNode)) *EvictedPopulationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FindProfileFunc = f
	return m.mock
}

//FindProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.EvictedPopulation interface
func (m *EvictedPopulationMock) FindProfile(p insolar.ShortNodeID) (r profiles.EvictedNode) {
	counter := atomic.AddUint64(&m.FindProfilePreCounter, 1)
	defer atomic.AddUint64(&m.FindProfileCounter, 1)

	if len(m.FindProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FindProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to EvictedPopulationMock.FindProfile. %v", p)
			return
		}

		input := m.FindProfileMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, EvictedPopulationMockFindProfileInput{p}, "EvictedPopulation.FindProfile got unexpected parameters")

		result := m.FindProfileMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the EvictedPopulationMock.FindProfile")
			return
		}

		r = result.r

		return
	}

	if m.FindProfileMock.mainExpectation != nil {

		input := m.FindProfileMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, EvictedPopulationMockFindProfileInput{p}, "EvictedPopulation.FindProfile got unexpected parameters")
		}

		result := m.FindProfileMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the EvictedPopulationMock.FindProfile")
		}

		r = result.r

		return
	}

	if m.FindProfileFunc == nil {
		m.t.Fatalf("Unexpected call to EvictedPopulationMock.FindProfile. %v", p)
		return
	}

	return m.FindProfileFunc(p)
}

//FindProfileMinimockCounter returns a count of EvictedPopulationMock.FindProfileFunc invocations
func (m *EvictedPopulationMock) FindProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FindProfileCounter)
}

//FindProfileMinimockPreCounter returns the value of EvictedPopulationMock.FindProfile invocations
func (m *EvictedPopulationMock) FindProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FindProfilePreCounter)
}

//FindProfileFinished returns true if mock invocations count is ok
func (m *EvictedPopulationMock) FindProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FindProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FindProfileCounter) == uint64(len(m.FindProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FindProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FindProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FindProfileFunc != nil {
		return atomic.LoadUint64(&m.FindProfileCounter) > 0
	}

	return true
}

type mEvictedPopulationMockGetCount struct {
	mock              *EvictedPopulationMock
	mainExpectation   *EvictedPopulationMockGetCountExpectation
	expectationSeries []*EvictedPopulationMockGetCountExpectation
}

type EvictedPopulationMockGetCountExpectation struct {
	result *EvictedPopulationMockGetCountResult
}

type EvictedPopulationMockGetCountResult struct {
	r int
}

//Expect specifies that invocation of EvictedPopulation.GetCount is expected from 1 to Infinity times
func (m *mEvictedPopulationMockGetCount) Expect() *mEvictedPopulationMockGetCount {
	m.mock.GetCountFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &EvictedPopulationMockGetCountExpectation{}
	}

	return m
}

//Return specifies results of invocation of EvictedPopulation.GetCount
func (m *mEvictedPopulationMockGetCount) Return(r int) *EvictedPopulationMock {
	m.mock.GetCountFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &EvictedPopulationMockGetCountExpectation{}
	}
	m.mainExpectation.result = &EvictedPopulationMockGetCountResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of EvictedPopulation.GetCount is expected once
func (m *mEvictedPopulationMockGetCount) ExpectOnce() *EvictedPopulationMockGetCountExpectation {
	m.mock.GetCountFunc = nil
	m.mainExpectation = nil

	expectation := &EvictedPopulationMockGetCountExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *EvictedPopulationMockGetCountExpectation) Return(r int) {
	e.result = &EvictedPopulationMockGetCountResult{r}
}

//Set uses given function f as a mock of EvictedPopulation.GetCount method
func (m *mEvictedPopulationMockGetCount) Set(f func() (r int)) *EvictedPopulationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCountFunc = f
	return m.mock
}

//GetCount implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.EvictedPopulation interface
func (m *EvictedPopulationMock) GetCount() (r int) {
	counter := atomic.AddUint64(&m.GetCountPreCounter, 1)
	defer atomic.AddUint64(&m.GetCountCounter, 1)

	if len(m.GetCountMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCountMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to EvictedPopulationMock.GetCount.")
			return
		}

		result := m.GetCountMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the EvictedPopulationMock.GetCount")
			return
		}

		r = result.r

		return
	}

	if m.GetCountMock.mainExpectation != nil {

		result := m.GetCountMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the EvictedPopulationMock.GetCount")
		}

		r = result.r

		return
	}

	if m.GetCountFunc == nil {
		m.t.Fatalf("Unexpected call to EvictedPopulationMock.GetCount.")
		return
	}

	return m.GetCountFunc()
}

//GetCountMinimockCounter returns a count of EvictedPopulationMock.GetCountFunc invocations
func (m *EvictedPopulationMock) GetCountMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCountCounter)
}

//GetCountMinimockPreCounter returns the value of EvictedPopulationMock.GetCount invocations
func (m *EvictedPopulationMock) GetCountMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCountPreCounter)
}

//GetCountFinished returns true if mock invocations count is ok
func (m *EvictedPopulationMock) GetCountFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCountMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCountCounter) == uint64(len(m.GetCountMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCountMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCountCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCountFunc != nil {
		return atomic.LoadUint64(&m.GetCountCounter) > 0
	}

	return true
}

type mEvictedPopulationMockGetProfiles struct {
	mock              *EvictedPopulationMock
	mainExpectation   *EvictedPopulationMockGetProfilesExpectation
	expectationSeries []*EvictedPopulationMockGetProfilesExpectation
}

type EvictedPopulationMockGetProfilesExpectation struct {
	result *EvictedPopulationMockGetProfilesResult
}

type EvictedPopulationMockGetProfilesResult struct {
	r []profiles.EvictedNode
}

//Expect specifies that invocation of EvictedPopulation.GetProfiles is expected from 1 to Infinity times
func (m *mEvictedPopulationMockGetProfiles) Expect() *mEvictedPopulationMockGetProfiles {
	m.mock.GetProfilesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &EvictedPopulationMockGetProfilesExpectation{}
	}

	return m
}

//Return specifies results of invocation of EvictedPopulation.GetProfiles
func (m *mEvictedPopulationMockGetProfiles) Return(r []profiles.EvictedNode) *EvictedPopulationMock {
	m.mock.GetProfilesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &EvictedPopulationMockGetProfilesExpectation{}
	}
	m.mainExpectation.result = &EvictedPopulationMockGetProfilesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of EvictedPopulation.GetProfiles is expected once
func (m *mEvictedPopulationMockGetProfiles) ExpectOnce() *EvictedPopulationMockGetProfilesExpectation {
	m.mock.GetProfilesFunc = nil
	m.mainExpectation = nil

	expectation := &EvictedPopulationMockGetProfilesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *EvictedPopulationMockGetProfilesExpectation) Return(r []profiles.EvictedNode) {
	e.result = &EvictedPopulationMockGetProfilesResult{r}
}

//Set uses given function f as a mock of EvictedPopulation.GetProfiles method
func (m *mEvictedPopulationMockGetProfiles) Set(f func() (r []profiles.EvictedNode)) *EvictedPopulationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetProfilesFunc = f
	return m.mock
}

//GetProfiles implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.EvictedPopulation interface
func (m *EvictedPopulationMock) GetProfiles() (r []profiles.EvictedNode) {
	counter := atomic.AddUint64(&m.GetProfilesPreCounter, 1)
	defer atomic.AddUint64(&m.GetProfilesCounter, 1)

	if len(m.GetProfilesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetProfilesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to EvictedPopulationMock.GetProfiles.")
			return
		}

		result := m.GetProfilesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the EvictedPopulationMock.GetProfiles")
			return
		}

		r = result.r

		return
	}

	if m.GetProfilesMock.mainExpectation != nil {

		result := m.GetProfilesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the EvictedPopulationMock.GetProfiles")
		}

		r = result.r

		return
	}

	if m.GetProfilesFunc == nil {
		m.t.Fatalf("Unexpected call to EvictedPopulationMock.GetProfiles.")
		return
	}

	return m.GetProfilesFunc()
}

//GetProfilesMinimockCounter returns a count of EvictedPopulationMock.GetProfilesFunc invocations
func (m *EvictedPopulationMock) GetProfilesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetProfilesCounter)
}

//GetProfilesMinimockPreCounter returns the value of EvictedPopulationMock.GetProfiles invocations
func (m *EvictedPopulationMock) GetProfilesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetProfilesPreCounter)
}

//GetProfilesFinished returns true if mock invocations count is ok
func (m *EvictedPopulationMock) GetProfilesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetProfilesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetProfilesCounter) == uint64(len(m.GetProfilesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetProfilesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetProfilesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetProfilesFunc != nil {
		return atomic.LoadUint64(&m.GetProfilesCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *EvictedPopulationMock) ValidateCallCounters() {

	if !m.FindProfileFinished() {
		m.t.Fatal("Expected call to EvictedPopulationMock.FindProfile")
	}

	if !m.GetCountFinished() {
		m.t.Fatal("Expected call to EvictedPopulationMock.GetCount")
	}

	if !m.GetProfilesFinished() {
		m.t.Fatal("Expected call to EvictedPopulationMock.GetProfiles")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *EvictedPopulationMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *EvictedPopulationMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *EvictedPopulationMock) MinimockFinish() {

	if !m.FindProfileFinished() {
		m.t.Fatal("Expected call to EvictedPopulationMock.FindProfile")
	}

	if !m.GetCountFinished() {
		m.t.Fatal("Expected call to EvictedPopulationMock.GetCount")
	}

	if !m.GetProfilesFinished() {
		m.t.Fatal("Expected call to EvictedPopulationMock.GetProfiles")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *EvictedPopulationMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *EvictedPopulationMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FindProfileFinished()
		ok = ok && m.GetCountFinished()
		ok = ok && m.GetProfilesFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FindProfileFinished() {
				m.t.Error("Expected call to EvictedPopulationMock.FindProfile")
			}

			if !m.GetCountFinished() {
				m.t.Error("Expected call to EvictedPopulationMock.GetCount")
			}

			if !m.GetProfilesFinished() {
				m.t.Error("Expected call to EvictedPopulationMock.GetProfiles")
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
func (m *EvictedPopulationMock) AllMocksCalled() bool {

	if !m.FindProfileFinished() {
		return false
	}

	if !m.GetCountFinished() {
		return false
	}

	if !m.GetProfilesFinished() {
		return false
	}

	return true
}
