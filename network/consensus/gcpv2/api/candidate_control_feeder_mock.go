package api

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CandidateControlFeeder" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	profiles "github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"

	testify_assert "github.com/stretchr/testify/assert"
)

//CandidateControlFeederMock implements github.com/insolar/insolar/network/consensus/gcpv2/api.CandidateControlFeeder
type CandidateControlFeederMock struct {
	t minimock.Tester

	PickNextJoinCandidateFunc       func() (r profiles.CandidateProfile, r1 cryptkit.DigestHolder)
	PickNextJoinCandidateCounter    uint64
	PickNextJoinCandidatePreCounter uint64
	PickNextJoinCandidateMock       mCandidateControlFeederMockPickNextJoinCandidate

	RemoveJoinCandidateFunc       func(p bool, p1 insolar.ShortNodeID) (r bool)
	RemoveJoinCandidateCounter    uint64
	RemoveJoinCandidatePreCounter uint64
	RemoveJoinCandidateMock       mCandidateControlFeederMockRemoveJoinCandidate
}

//NewCandidateControlFeederMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api.CandidateControlFeeder
func NewCandidateControlFeederMock(t minimock.Tester) *CandidateControlFeederMock {
	m := &CandidateControlFeederMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PickNextJoinCandidateMock = mCandidateControlFeederMockPickNextJoinCandidate{mock: m}
	m.RemoveJoinCandidateMock = mCandidateControlFeederMockRemoveJoinCandidate{mock: m}

	return m
}

type mCandidateControlFeederMockPickNextJoinCandidate struct {
	mock              *CandidateControlFeederMock
	mainExpectation   *CandidateControlFeederMockPickNextJoinCandidateExpectation
	expectationSeries []*CandidateControlFeederMockPickNextJoinCandidateExpectation
}

type CandidateControlFeederMockPickNextJoinCandidateExpectation struct {
	result *CandidateControlFeederMockPickNextJoinCandidateResult
}

type CandidateControlFeederMockPickNextJoinCandidateResult struct {
	r  profiles.CandidateProfile
	r1 cryptkit.DigestHolder
}

//Expect specifies that invocation of CandidateControlFeeder.PickNextJoinCandidate is expected from 1 to Infinity times
func (m *mCandidateControlFeederMockPickNextJoinCandidate) Expect() *mCandidateControlFeederMockPickNextJoinCandidate {
	m.mock.PickNextJoinCandidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateControlFeederMockPickNextJoinCandidateExpectation{}
	}

	return m
}

//Return specifies results of invocation of CandidateControlFeeder.PickNextJoinCandidate
func (m *mCandidateControlFeederMockPickNextJoinCandidate) Return(r profiles.CandidateProfile, r1 cryptkit.DigestHolder) *CandidateControlFeederMock {
	m.mock.PickNextJoinCandidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateControlFeederMockPickNextJoinCandidateExpectation{}
	}
	m.mainExpectation.result = &CandidateControlFeederMockPickNextJoinCandidateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateControlFeeder.PickNextJoinCandidate is expected once
func (m *mCandidateControlFeederMockPickNextJoinCandidate) ExpectOnce() *CandidateControlFeederMockPickNextJoinCandidateExpectation {
	m.mock.PickNextJoinCandidateFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateControlFeederMockPickNextJoinCandidateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateControlFeederMockPickNextJoinCandidateExpectation) Return(r profiles.CandidateProfile, r1 cryptkit.DigestHolder) {
	e.result = &CandidateControlFeederMockPickNextJoinCandidateResult{r, r1}
}

//Set uses given function f as a mock of CandidateControlFeeder.PickNextJoinCandidate method
func (m *mCandidateControlFeederMockPickNextJoinCandidate) Set(f func() (r profiles.CandidateProfile, r1 cryptkit.DigestHolder)) *CandidateControlFeederMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PickNextJoinCandidateFunc = f
	return m.mock
}

//PickNextJoinCandidate implements github.com/insolar/insolar/network/consensus/gcpv2/api.CandidateControlFeeder interface
func (m *CandidateControlFeederMock) PickNextJoinCandidate() (r profiles.CandidateProfile, r1 cryptkit.DigestHolder) {
	counter := atomic.AddUint64(&m.PickNextJoinCandidatePreCounter, 1)
	defer atomic.AddUint64(&m.PickNextJoinCandidateCounter, 1)

	if len(m.PickNextJoinCandidateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PickNextJoinCandidateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateControlFeederMock.PickNextJoinCandidate.")
			return
		}

		result := m.PickNextJoinCandidateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateControlFeederMock.PickNextJoinCandidate")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PickNextJoinCandidateMock.mainExpectation != nil {

		result := m.PickNextJoinCandidateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateControlFeederMock.PickNextJoinCandidate")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PickNextJoinCandidateFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateControlFeederMock.PickNextJoinCandidate.")
		return
	}

	return m.PickNextJoinCandidateFunc()
}

//PickNextJoinCandidateMinimockCounter returns a count of CandidateControlFeederMock.PickNextJoinCandidateFunc invocations
func (m *CandidateControlFeederMock) PickNextJoinCandidateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PickNextJoinCandidateCounter)
}

//PickNextJoinCandidateMinimockPreCounter returns the value of CandidateControlFeederMock.PickNextJoinCandidate invocations
func (m *CandidateControlFeederMock) PickNextJoinCandidateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PickNextJoinCandidatePreCounter)
}

//PickNextJoinCandidateFinished returns true if mock invocations count is ok
func (m *CandidateControlFeederMock) PickNextJoinCandidateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PickNextJoinCandidateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PickNextJoinCandidateCounter) == uint64(len(m.PickNextJoinCandidateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PickNextJoinCandidateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PickNextJoinCandidateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PickNextJoinCandidateFunc != nil {
		return atomic.LoadUint64(&m.PickNextJoinCandidateCounter) > 0
	}

	return true
}

type mCandidateControlFeederMockRemoveJoinCandidate struct {
	mock              *CandidateControlFeederMock
	mainExpectation   *CandidateControlFeederMockRemoveJoinCandidateExpectation
	expectationSeries []*CandidateControlFeederMockRemoveJoinCandidateExpectation
}

type CandidateControlFeederMockRemoveJoinCandidateExpectation struct {
	input  *CandidateControlFeederMockRemoveJoinCandidateInput
	result *CandidateControlFeederMockRemoveJoinCandidateResult
}

type CandidateControlFeederMockRemoveJoinCandidateInput struct {
	p  bool
	p1 insolar.ShortNodeID
}

type CandidateControlFeederMockRemoveJoinCandidateResult struct {
	r bool
}

//Expect specifies that invocation of CandidateControlFeeder.RemoveJoinCandidate is expected from 1 to Infinity times
func (m *mCandidateControlFeederMockRemoveJoinCandidate) Expect(p bool, p1 insolar.ShortNodeID) *mCandidateControlFeederMockRemoveJoinCandidate {
	m.mock.RemoveJoinCandidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateControlFeederMockRemoveJoinCandidateExpectation{}
	}
	m.mainExpectation.input = &CandidateControlFeederMockRemoveJoinCandidateInput{p, p1}
	return m
}

//Return specifies results of invocation of CandidateControlFeeder.RemoveJoinCandidate
func (m *mCandidateControlFeederMockRemoveJoinCandidate) Return(r bool) *CandidateControlFeederMock {
	m.mock.RemoveJoinCandidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CandidateControlFeederMockRemoveJoinCandidateExpectation{}
	}
	m.mainExpectation.result = &CandidateControlFeederMockRemoveJoinCandidateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CandidateControlFeeder.RemoveJoinCandidate is expected once
func (m *mCandidateControlFeederMockRemoveJoinCandidate) ExpectOnce(p bool, p1 insolar.ShortNodeID) *CandidateControlFeederMockRemoveJoinCandidateExpectation {
	m.mock.RemoveJoinCandidateFunc = nil
	m.mainExpectation = nil

	expectation := &CandidateControlFeederMockRemoveJoinCandidateExpectation{}
	expectation.input = &CandidateControlFeederMockRemoveJoinCandidateInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CandidateControlFeederMockRemoveJoinCandidateExpectation) Return(r bool) {
	e.result = &CandidateControlFeederMockRemoveJoinCandidateResult{r}
}

//Set uses given function f as a mock of CandidateControlFeeder.RemoveJoinCandidate method
func (m *mCandidateControlFeederMockRemoveJoinCandidate) Set(f func(p bool, p1 insolar.ShortNodeID) (r bool)) *CandidateControlFeederMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveJoinCandidateFunc = f
	return m.mock
}

//RemoveJoinCandidate implements github.com/insolar/insolar/network/consensus/gcpv2/api.CandidateControlFeeder interface
func (m *CandidateControlFeederMock) RemoveJoinCandidate(p bool, p1 insolar.ShortNodeID) (r bool) {
	counter := atomic.AddUint64(&m.RemoveJoinCandidatePreCounter, 1)
	defer atomic.AddUint64(&m.RemoveJoinCandidateCounter, 1)

	if len(m.RemoveJoinCandidateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveJoinCandidateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CandidateControlFeederMock.RemoveJoinCandidate. %v %v", p, p1)
			return
		}

		input := m.RemoveJoinCandidateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CandidateControlFeederMockRemoveJoinCandidateInput{p, p1}, "CandidateControlFeeder.RemoveJoinCandidate got unexpected parameters")

		result := m.RemoveJoinCandidateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateControlFeederMock.RemoveJoinCandidate")
			return
		}

		r = result.r

		return
	}

	if m.RemoveJoinCandidateMock.mainExpectation != nil {

		input := m.RemoveJoinCandidateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CandidateControlFeederMockRemoveJoinCandidateInput{p, p1}, "CandidateControlFeeder.RemoveJoinCandidate got unexpected parameters")
		}

		result := m.RemoveJoinCandidateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CandidateControlFeederMock.RemoveJoinCandidate")
		}

		r = result.r

		return
	}

	if m.RemoveJoinCandidateFunc == nil {
		m.t.Fatalf("Unexpected call to CandidateControlFeederMock.RemoveJoinCandidate. %v %v", p, p1)
		return
	}

	return m.RemoveJoinCandidateFunc(p, p1)
}

//RemoveJoinCandidateMinimockCounter returns a count of CandidateControlFeederMock.RemoveJoinCandidateFunc invocations
func (m *CandidateControlFeederMock) RemoveJoinCandidateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJoinCandidateCounter)
}

//RemoveJoinCandidateMinimockPreCounter returns the value of CandidateControlFeederMock.RemoveJoinCandidate invocations
func (m *CandidateControlFeederMock) RemoveJoinCandidateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveJoinCandidatePreCounter)
}

//RemoveJoinCandidateFinished returns true if mock invocations count is ok
func (m *CandidateControlFeederMock) RemoveJoinCandidateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveJoinCandidateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveJoinCandidateCounter) == uint64(len(m.RemoveJoinCandidateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveJoinCandidateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveJoinCandidateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveJoinCandidateFunc != nil {
		return atomic.LoadUint64(&m.RemoveJoinCandidateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CandidateControlFeederMock) ValidateCallCounters() {

	if !m.PickNextJoinCandidateFinished() {
		m.t.Fatal("Expected call to CandidateControlFeederMock.PickNextJoinCandidate")
	}

	if !m.RemoveJoinCandidateFinished() {
		m.t.Fatal("Expected call to CandidateControlFeederMock.RemoveJoinCandidate")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CandidateControlFeederMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CandidateControlFeederMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CandidateControlFeederMock) MinimockFinish() {

	if !m.PickNextJoinCandidateFinished() {
		m.t.Fatal("Expected call to CandidateControlFeederMock.PickNextJoinCandidate")
	}

	if !m.RemoveJoinCandidateFinished() {
		m.t.Fatal("Expected call to CandidateControlFeederMock.RemoveJoinCandidate")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CandidateControlFeederMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CandidateControlFeederMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.PickNextJoinCandidateFinished()
		ok = ok && m.RemoveJoinCandidateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.PickNextJoinCandidateFinished() {
				m.t.Error("Expected call to CandidateControlFeederMock.PickNextJoinCandidate")
			}

			if !m.RemoveJoinCandidateFinished() {
				m.t.Error("Expected call to CandidateControlFeederMock.RemoveJoinCandidate")
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
func (m *CandidateControlFeederMock) AllMocksCalled() bool {

	if !m.PickNextJoinCandidateFinished() {
		return false
	}

	if !m.RemoveJoinCandidateFinished() {
		return false
	}

	return true
}
