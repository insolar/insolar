package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Factory" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//FactoryMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Factory
type FactoryMock struct {
	t minimock.Tester

	CreateBriefIntroProfileFunc       func(p BriefCandidateProfile) (r StaticProfile)
	CreateBriefIntroProfileCounter    uint64
	CreateBriefIntroProfilePreCounter uint64
	CreateBriefIntroProfileMock       mFactoryMockCreateBriefIntroProfile

	CreateFullIntroProfileFunc       func(p CandidateProfile) (r StaticProfile)
	CreateFullIntroProfileCounter    uint64
	CreateFullIntroProfilePreCounter uint64
	CreateFullIntroProfileMock       mFactoryMockCreateFullIntroProfile

	CreateUpgradableIntroProfileFunc       func(p BriefCandidateProfile) (r StaticProfile)
	CreateUpgradableIntroProfileCounter    uint64
	CreateUpgradableIntroProfilePreCounter uint64
	CreateUpgradableIntroProfileMock       mFactoryMockCreateUpgradableIntroProfile

	TryConvertUpgradableIntroProfileFunc       func(p StaticProfile) (r StaticProfile, r1 bool)
	TryConvertUpgradableIntroProfileCounter    uint64
	TryConvertUpgradableIntroProfilePreCounter uint64
	TryConvertUpgradableIntroProfileMock       mFactoryMockTryConvertUpgradableIntroProfile
}

//NewFactoryMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Factory
func NewFactoryMock(t minimock.Tester) *FactoryMock {
	m := &FactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateBriefIntroProfileMock = mFactoryMockCreateBriefIntroProfile{mock: m}
	m.CreateFullIntroProfileMock = mFactoryMockCreateFullIntroProfile{mock: m}
	m.CreateUpgradableIntroProfileMock = mFactoryMockCreateUpgradableIntroProfile{mock: m}
	m.TryConvertUpgradableIntroProfileMock = mFactoryMockTryConvertUpgradableIntroProfile{mock: m}

	return m
}

type mFactoryMockCreateBriefIntroProfile struct {
	mock              *FactoryMock
	mainExpectation   *FactoryMockCreateBriefIntroProfileExpectation
	expectationSeries []*FactoryMockCreateBriefIntroProfileExpectation
}

type FactoryMockCreateBriefIntroProfileExpectation struct {
	input  *FactoryMockCreateBriefIntroProfileInput
	result *FactoryMockCreateBriefIntroProfileResult
}

type FactoryMockCreateBriefIntroProfileInput struct {
	p BriefCandidateProfile
}

type FactoryMockCreateBriefIntroProfileResult struct {
	r StaticProfile
}

//Expect specifies that invocation of Factory.CreateBriefIntroProfile is expected from 1 to Infinity times
func (m *mFactoryMockCreateBriefIntroProfile) Expect(p BriefCandidateProfile) *mFactoryMockCreateBriefIntroProfile {
	m.mock.CreateBriefIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockCreateBriefIntroProfileExpectation{}
	}
	m.mainExpectation.input = &FactoryMockCreateBriefIntroProfileInput{p}
	return m
}

//Return specifies results of invocation of Factory.CreateBriefIntroProfile
func (m *mFactoryMockCreateBriefIntroProfile) Return(r StaticProfile) *FactoryMock {
	m.mock.CreateBriefIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockCreateBriefIntroProfileExpectation{}
	}
	m.mainExpectation.result = &FactoryMockCreateBriefIntroProfileResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Factory.CreateBriefIntroProfile is expected once
func (m *mFactoryMockCreateBriefIntroProfile) ExpectOnce(p BriefCandidateProfile) *FactoryMockCreateBriefIntroProfileExpectation {
	m.mock.CreateBriefIntroProfileFunc = nil
	m.mainExpectation = nil

	expectation := &FactoryMockCreateBriefIntroProfileExpectation{}
	expectation.input = &FactoryMockCreateBriefIntroProfileInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FactoryMockCreateBriefIntroProfileExpectation) Return(r StaticProfile) {
	e.result = &FactoryMockCreateBriefIntroProfileResult{r}
}

//Set uses given function f as a mock of Factory.CreateBriefIntroProfile method
func (m *mFactoryMockCreateBriefIntroProfile) Set(f func(p BriefCandidateProfile) (r StaticProfile)) *FactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateBriefIntroProfileFunc = f
	return m.mock
}

//CreateBriefIntroProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Factory interface
func (m *FactoryMock) CreateBriefIntroProfile(p BriefCandidateProfile) (r StaticProfile) {
	counter := atomic.AddUint64(&m.CreateBriefIntroProfilePreCounter, 1)
	defer atomic.AddUint64(&m.CreateBriefIntroProfileCounter, 1)

	if len(m.CreateBriefIntroProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateBriefIntroProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FactoryMock.CreateBriefIntroProfile. %v", p)
			return
		}

		input := m.CreateBriefIntroProfileMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FactoryMockCreateBriefIntroProfileInput{p}, "Factory.CreateBriefIntroProfile got unexpected parameters")

		result := m.CreateBriefIntroProfileMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.CreateBriefIntroProfile")
			return
		}

		r = result.r

		return
	}

	if m.CreateBriefIntroProfileMock.mainExpectation != nil {

		input := m.CreateBriefIntroProfileMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FactoryMockCreateBriefIntroProfileInput{p}, "Factory.CreateBriefIntroProfile got unexpected parameters")
		}

		result := m.CreateBriefIntroProfileMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.CreateBriefIntroProfile")
		}

		r = result.r

		return
	}

	if m.CreateBriefIntroProfileFunc == nil {
		m.t.Fatalf("Unexpected call to FactoryMock.CreateBriefIntroProfile. %v", p)
		return
	}

	return m.CreateBriefIntroProfileFunc(p)
}

//CreateBriefIntroProfileMinimockCounter returns a count of FactoryMock.CreateBriefIntroProfileFunc invocations
func (m *FactoryMock) CreateBriefIntroProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateBriefIntroProfileCounter)
}

//CreateBriefIntroProfileMinimockPreCounter returns the value of FactoryMock.CreateBriefIntroProfile invocations
func (m *FactoryMock) CreateBriefIntroProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateBriefIntroProfilePreCounter)
}

//CreateBriefIntroProfileFinished returns true if mock invocations count is ok
func (m *FactoryMock) CreateBriefIntroProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreateBriefIntroProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreateBriefIntroProfileCounter) == uint64(len(m.CreateBriefIntroProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreateBriefIntroProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreateBriefIntroProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreateBriefIntroProfileFunc != nil {
		return atomic.LoadUint64(&m.CreateBriefIntroProfileCounter) > 0
	}

	return true
}

type mFactoryMockCreateFullIntroProfile struct {
	mock              *FactoryMock
	mainExpectation   *FactoryMockCreateFullIntroProfileExpectation
	expectationSeries []*FactoryMockCreateFullIntroProfileExpectation
}

type FactoryMockCreateFullIntroProfileExpectation struct {
	input  *FactoryMockCreateFullIntroProfileInput
	result *FactoryMockCreateFullIntroProfileResult
}

type FactoryMockCreateFullIntroProfileInput struct {
	p CandidateProfile
}

type FactoryMockCreateFullIntroProfileResult struct {
	r StaticProfile
}

//Expect specifies that invocation of Factory.CreateFullIntroProfile is expected from 1 to Infinity times
func (m *mFactoryMockCreateFullIntroProfile) Expect(p CandidateProfile) *mFactoryMockCreateFullIntroProfile {
	m.mock.CreateFullIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockCreateFullIntroProfileExpectation{}
	}
	m.mainExpectation.input = &FactoryMockCreateFullIntroProfileInput{p}
	return m
}

//Return specifies results of invocation of Factory.CreateFullIntroProfile
func (m *mFactoryMockCreateFullIntroProfile) Return(r StaticProfile) *FactoryMock {
	m.mock.CreateFullIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockCreateFullIntroProfileExpectation{}
	}
	m.mainExpectation.result = &FactoryMockCreateFullIntroProfileResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Factory.CreateFullIntroProfile is expected once
func (m *mFactoryMockCreateFullIntroProfile) ExpectOnce(p CandidateProfile) *FactoryMockCreateFullIntroProfileExpectation {
	m.mock.CreateFullIntroProfileFunc = nil
	m.mainExpectation = nil

	expectation := &FactoryMockCreateFullIntroProfileExpectation{}
	expectation.input = &FactoryMockCreateFullIntroProfileInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FactoryMockCreateFullIntroProfileExpectation) Return(r StaticProfile) {
	e.result = &FactoryMockCreateFullIntroProfileResult{r}
}

//Set uses given function f as a mock of Factory.CreateFullIntroProfile method
func (m *mFactoryMockCreateFullIntroProfile) Set(f func(p CandidateProfile) (r StaticProfile)) *FactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateFullIntroProfileFunc = f
	return m.mock
}

//CreateFullIntroProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Factory interface
func (m *FactoryMock) CreateFullIntroProfile(p CandidateProfile) (r StaticProfile) {
	counter := atomic.AddUint64(&m.CreateFullIntroProfilePreCounter, 1)
	defer atomic.AddUint64(&m.CreateFullIntroProfileCounter, 1)

	if len(m.CreateFullIntroProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateFullIntroProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FactoryMock.CreateFullIntroProfile. %v", p)
			return
		}

		input := m.CreateFullIntroProfileMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FactoryMockCreateFullIntroProfileInput{p}, "Factory.CreateFullIntroProfile got unexpected parameters")

		result := m.CreateFullIntroProfileMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.CreateFullIntroProfile")
			return
		}

		r = result.r

		return
	}

	if m.CreateFullIntroProfileMock.mainExpectation != nil {

		input := m.CreateFullIntroProfileMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FactoryMockCreateFullIntroProfileInput{p}, "Factory.CreateFullIntroProfile got unexpected parameters")
		}

		result := m.CreateFullIntroProfileMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.CreateFullIntroProfile")
		}

		r = result.r

		return
	}

	if m.CreateFullIntroProfileFunc == nil {
		m.t.Fatalf("Unexpected call to FactoryMock.CreateFullIntroProfile. %v", p)
		return
	}

	return m.CreateFullIntroProfileFunc(p)
}

//CreateFullIntroProfileMinimockCounter returns a count of FactoryMock.CreateFullIntroProfileFunc invocations
func (m *FactoryMock) CreateFullIntroProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateFullIntroProfileCounter)
}

//CreateFullIntroProfileMinimockPreCounter returns the value of FactoryMock.CreateFullIntroProfile invocations
func (m *FactoryMock) CreateFullIntroProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateFullIntroProfilePreCounter)
}

//CreateFullIntroProfileFinished returns true if mock invocations count is ok
func (m *FactoryMock) CreateFullIntroProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreateFullIntroProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreateFullIntroProfileCounter) == uint64(len(m.CreateFullIntroProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreateFullIntroProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreateFullIntroProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreateFullIntroProfileFunc != nil {
		return atomic.LoadUint64(&m.CreateFullIntroProfileCounter) > 0
	}

	return true
}

type mFactoryMockCreateUpgradableIntroProfile struct {
	mock              *FactoryMock
	mainExpectation   *FactoryMockCreateUpgradableIntroProfileExpectation
	expectationSeries []*FactoryMockCreateUpgradableIntroProfileExpectation
}

type FactoryMockCreateUpgradableIntroProfileExpectation struct {
	input  *FactoryMockCreateUpgradableIntroProfileInput
	result *FactoryMockCreateUpgradableIntroProfileResult
}

type FactoryMockCreateUpgradableIntroProfileInput struct {
	p BriefCandidateProfile
}

type FactoryMockCreateUpgradableIntroProfileResult struct {
	r StaticProfile
}

//Expect specifies that invocation of Factory.CreateUpgradableIntroProfile is expected from 1 to Infinity times
func (m *mFactoryMockCreateUpgradableIntroProfile) Expect(p BriefCandidateProfile) *mFactoryMockCreateUpgradableIntroProfile {
	m.mock.CreateUpgradableIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockCreateUpgradableIntroProfileExpectation{}
	}
	m.mainExpectation.input = &FactoryMockCreateUpgradableIntroProfileInput{p}
	return m
}

//Return specifies results of invocation of Factory.CreateUpgradableIntroProfile
func (m *mFactoryMockCreateUpgradableIntroProfile) Return(r StaticProfile) *FactoryMock {
	m.mock.CreateUpgradableIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockCreateUpgradableIntroProfileExpectation{}
	}
	m.mainExpectation.result = &FactoryMockCreateUpgradableIntroProfileResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Factory.CreateUpgradableIntroProfile is expected once
func (m *mFactoryMockCreateUpgradableIntroProfile) ExpectOnce(p BriefCandidateProfile) *FactoryMockCreateUpgradableIntroProfileExpectation {
	m.mock.CreateUpgradableIntroProfileFunc = nil
	m.mainExpectation = nil

	expectation := &FactoryMockCreateUpgradableIntroProfileExpectation{}
	expectation.input = &FactoryMockCreateUpgradableIntroProfileInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FactoryMockCreateUpgradableIntroProfileExpectation) Return(r StaticProfile) {
	e.result = &FactoryMockCreateUpgradableIntroProfileResult{r}
}

//Set uses given function f as a mock of Factory.CreateUpgradableIntroProfile method
func (m *mFactoryMockCreateUpgradableIntroProfile) Set(f func(p BriefCandidateProfile) (r StaticProfile)) *FactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateUpgradableIntroProfileFunc = f
	return m.mock
}

//CreateUpgradableIntroProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Factory interface
func (m *FactoryMock) CreateUpgradableIntroProfile(p BriefCandidateProfile) (r StaticProfile) {
	counter := atomic.AddUint64(&m.CreateUpgradableIntroProfilePreCounter, 1)
	defer atomic.AddUint64(&m.CreateUpgradableIntroProfileCounter, 1)

	if len(m.CreateUpgradableIntroProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateUpgradableIntroProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FactoryMock.CreateUpgradableIntroProfile. %v", p)
			return
		}

		input := m.CreateUpgradableIntroProfileMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FactoryMockCreateUpgradableIntroProfileInput{p}, "Factory.CreateUpgradableIntroProfile got unexpected parameters")

		result := m.CreateUpgradableIntroProfileMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.CreateUpgradableIntroProfile")
			return
		}

		r = result.r

		return
	}

	if m.CreateUpgradableIntroProfileMock.mainExpectation != nil {

		input := m.CreateUpgradableIntroProfileMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FactoryMockCreateUpgradableIntroProfileInput{p}, "Factory.CreateUpgradableIntroProfile got unexpected parameters")
		}

		result := m.CreateUpgradableIntroProfileMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.CreateUpgradableIntroProfile")
		}

		r = result.r

		return
	}

	if m.CreateUpgradableIntroProfileFunc == nil {
		m.t.Fatalf("Unexpected call to FactoryMock.CreateUpgradableIntroProfile. %v", p)
		return
	}

	return m.CreateUpgradableIntroProfileFunc(p)
}

//CreateUpgradableIntroProfileMinimockCounter returns a count of FactoryMock.CreateUpgradableIntroProfileFunc invocations
func (m *FactoryMock) CreateUpgradableIntroProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateUpgradableIntroProfileCounter)
}

//CreateUpgradableIntroProfileMinimockPreCounter returns the value of FactoryMock.CreateUpgradableIntroProfile invocations
func (m *FactoryMock) CreateUpgradableIntroProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateUpgradableIntroProfilePreCounter)
}

//CreateUpgradableIntroProfileFinished returns true if mock invocations count is ok
func (m *FactoryMock) CreateUpgradableIntroProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreateUpgradableIntroProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreateUpgradableIntroProfileCounter) == uint64(len(m.CreateUpgradableIntroProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreateUpgradableIntroProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreateUpgradableIntroProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreateUpgradableIntroProfileFunc != nil {
		return atomic.LoadUint64(&m.CreateUpgradableIntroProfileCounter) > 0
	}

	return true
}

type mFactoryMockTryConvertUpgradableIntroProfile struct {
	mock              *FactoryMock
	mainExpectation   *FactoryMockTryConvertUpgradableIntroProfileExpectation
	expectationSeries []*FactoryMockTryConvertUpgradableIntroProfileExpectation
}

type FactoryMockTryConvertUpgradableIntroProfileExpectation struct {
	input  *FactoryMockTryConvertUpgradableIntroProfileInput
	result *FactoryMockTryConvertUpgradableIntroProfileResult
}

type FactoryMockTryConvertUpgradableIntroProfileInput struct {
	p StaticProfile
}

type FactoryMockTryConvertUpgradableIntroProfileResult struct {
	r  StaticProfile
	r1 bool
}

//Expect specifies that invocation of Factory.TryConvertUpgradableIntroProfile is expected from 1 to Infinity times
func (m *mFactoryMockTryConvertUpgradableIntroProfile) Expect(p StaticProfile) *mFactoryMockTryConvertUpgradableIntroProfile {
	m.mock.TryConvertUpgradableIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockTryConvertUpgradableIntroProfileExpectation{}
	}
	m.mainExpectation.input = &FactoryMockTryConvertUpgradableIntroProfileInput{p}
	return m
}

//Return specifies results of invocation of Factory.TryConvertUpgradableIntroProfile
func (m *mFactoryMockTryConvertUpgradableIntroProfile) Return(r StaticProfile, r1 bool) *FactoryMock {
	m.mock.TryConvertUpgradableIntroProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FactoryMockTryConvertUpgradableIntroProfileExpectation{}
	}
	m.mainExpectation.result = &FactoryMockTryConvertUpgradableIntroProfileResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Factory.TryConvertUpgradableIntroProfile is expected once
func (m *mFactoryMockTryConvertUpgradableIntroProfile) ExpectOnce(p StaticProfile) *FactoryMockTryConvertUpgradableIntroProfileExpectation {
	m.mock.TryConvertUpgradableIntroProfileFunc = nil
	m.mainExpectation = nil

	expectation := &FactoryMockTryConvertUpgradableIntroProfileExpectation{}
	expectation.input = &FactoryMockTryConvertUpgradableIntroProfileInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FactoryMockTryConvertUpgradableIntroProfileExpectation) Return(r StaticProfile, r1 bool) {
	e.result = &FactoryMockTryConvertUpgradableIntroProfileResult{r, r1}
}

//Set uses given function f as a mock of Factory.TryConvertUpgradableIntroProfile method
func (m *mFactoryMockTryConvertUpgradableIntroProfile) Set(f func(p StaticProfile) (r StaticProfile, r1 bool)) *FactoryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.TryConvertUpgradableIntroProfileFunc = f
	return m.mock
}

//TryConvertUpgradableIntroProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Factory interface
func (m *FactoryMock) TryConvertUpgradableIntroProfile(p StaticProfile) (r StaticProfile, r1 bool) {
	counter := atomic.AddUint64(&m.TryConvertUpgradableIntroProfilePreCounter, 1)
	defer atomic.AddUint64(&m.TryConvertUpgradableIntroProfileCounter, 1)

	if len(m.TryConvertUpgradableIntroProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.TryConvertUpgradableIntroProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FactoryMock.TryConvertUpgradableIntroProfile. %v", p)
			return
		}

		input := m.TryConvertUpgradableIntroProfileMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FactoryMockTryConvertUpgradableIntroProfileInput{p}, "Factory.TryConvertUpgradableIntroProfile got unexpected parameters")

		result := m.TryConvertUpgradableIntroProfileMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.TryConvertUpgradableIntroProfile")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.TryConvertUpgradableIntroProfileMock.mainExpectation != nil {

		input := m.TryConvertUpgradableIntroProfileMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FactoryMockTryConvertUpgradableIntroProfileInput{p}, "Factory.TryConvertUpgradableIntroProfile got unexpected parameters")
		}

		result := m.TryConvertUpgradableIntroProfileMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FactoryMock.TryConvertUpgradableIntroProfile")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.TryConvertUpgradableIntroProfileFunc == nil {
		m.t.Fatalf("Unexpected call to FactoryMock.TryConvertUpgradableIntroProfile. %v", p)
		return
	}

	return m.TryConvertUpgradableIntroProfileFunc(p)
}

//TryConvertUpgradableIntroProfileMinimockCounter returns a count of FactoryMock.TryConvertUpgradableIntroProfileFunc invocations
func (m *FactoryMock) TryConvertUpgradableIntroProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.TryConvertUpgradableIntroProfileCounter)
}

//TryConvertUpgradableIntroProfileMinimockPreCounter returns the value of FactoryMock.TryConvertUpgradableIntroProfile invocations
func (m *FactoryMock) TryConvertUpgradableIntroProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.TryConvertUpgradableIntroProfilePreCounter)
}

//TryConvertUpgradableIntroProfileFinished returns true if mock invocations count is ok
func (m *FactoryMock) TryConvertUpgradableIntroProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.TryConvertUpgradableIntroProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.TryConvertUpgradableIntroProfileCounter) == uint64(len(m.TryConvertUpgradableIntroProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.TryConvertUpgradableIntroProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.TryConvertUpgradableIntroProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.TryConvertUpgradableIntroProfileFunc != nil {
		return atomic.LoadUint64(&m.TryConvertUpgradableIntroProfileCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FactoryMock) ValidateCallCounters() {

	if !m.CreateBriefIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.CreateBriefIntroProfile")
	}

	if !m.CreateFullIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.CreateFullIntroProfile")
	}

	if !m.CreateUpgradableIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.CreateUpgradableIntroProfile")
	}

	if !m.TryConvertUpgradableIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.TryConvertUpgradableIntroProfile")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FactoryMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FactoryMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FactoryMock) MinimockFinish() {

	if !m.CreateBriefIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.CreateBriefIntroProfile")
	}

	if !m.CreateFullIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.CreateFullIntroProfile")
	}

	if !m.CreateUpgradableIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.CreateUpgradableIntroProfile")
	}

	if !m.TryConvertUpgradableIntroProfileFinished() {
		m.t.Fatal("Expected call to FactoryMock.TryConvertUpgradableIntroProfile")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FactoryMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FactoryMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CreateBriefIntroProfileFinished()
		ok = ok && m.CreateFullIntroProfileFinished()
		ok = ok && m.CreateUpgradableIntroProfileFinished()
		ok = ok && m.TryConvertUpgradableIntroProfileFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CreateBriefIntroProfileFinished() {
				m.t.Error("Expected call to FactoryMock.CreateBriefIntroProfile")
			}

			if !m.CreateFullIntroProfileFinished() {
				m.t.Error("Expected call to FactoryMock.CreateFullIntroProfile")
			}

			if !m.CreateUpgradableIntroProfileFinished() {
				m.t.Error("Expected call to FactoryMock.CreateUpgradableIntroProfile")
			}

			if !m.TryConvertUpgradableIntroProfileFinished() {
				m.t.Error("Expected call to FactoryMock.TryConvertUpgradableIntroProfile")
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
func (m *FactoryMock) AllMocksCalled() bool {

	if !m.CreateBriefIntroProfileFinished() {
		return false
	}

	if !m.CreateFullIntroProfileFinished() {
		return false
	}

	if !m.CreateUpgradableIntroProfileFinished() {
		return false
	}

	if !m.TryConvertUpgradableIntroProfileFinished() {
		return false
	}

	return true
}
