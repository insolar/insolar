package census

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MandateRegistry" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/census
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	endpoints "github.com/insolar/insolar/network/consensus/common/endpoints"
	profiles "github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	proofs "github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"

	testify_assert "github.com/stretchr/testify/assert"
)

//MandateRegistryMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.MandateRegistry
type MandateRegistryMock struct {
	t minimock.Tester

	FindRegisteredProfileFunc       func(p endpoints.Inbound) (r profiles.Host)
	FindRegisteredProfileCounter    uint64
	FindRegisteredProfilePreCounter uint64
	FindRegisteredProfileMock       mMandateRegistryMockFindRegisteredProfile

	GetCloudIdentityFunc       func() (r cryptkit.DigestHolder)
	GetCloudIdentityCounter    uint64
	GetCloudIdentityPreCounter uint64
	GetCloudIdentityMock       mMandateRegistryMockGetCloudIdentity

	GetConsensusConfigurationFunc       func() (r ConsensusConfiguration)
	GetConsensusConfigurationCounter    uint64
	GetConsensusConfigurationPreCounter uint64
	GetConsensusConfigurationMock       mMandateRegistryMockGetConsensusConfiguration

	GetPrimingCloudHashFunc       func() (r proofs.CloudStateHash)
	GetPrimingCloudHashCounter    uint64
	GetPrimingCloudHashPreCounter uint64
	GetPrimingCloudHashMock       mMandateRegistryMockGetPrimingCloudHash
}

//NewMandateRegistryMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/census.MandateRegistry
func NewMandateRegistryMock(t minimock.Tester) *MandateRegistryMock {
	m := &MandateRegistryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FindRegisteredProfileMock = mMandateRegistryMockFindRegisteredProfile{mock: m}
	m.GetCloudIdentityMock = mMandateRegistryMockGetCloudIdentity{mock: m}
	m.GetConsensusConfigurationMock = mMandateRegistryMockGetConsensusConfiguration{mock: m}
	m.GetPrimingCloudHashMock = mMandateRegistryMockGetPrimingCloudHash{mock: m}

	return m
}

type mMandateRegistryMockFindRegisteredProfile struct {
	mock              *MandateRegistryMock
	mainExpectation   *MandateRegistryMockFindRegisteredProfileExpectation
	expectationSeries []*MandateRegistryMockFindRegisteredProfileExpectation
}

type MandateRegistryMockFindRegisteredProfileExpectation struct {
	input  *MandateRegistryMockFindRegisteredProfileInput
	result *MandateRegistryMockFindRegisteredProfileResult
}

type MandateRegistryMockFindRegisteredProfileInput struct {
	p endpoints.Inbound
}

type MandateRegistryMockFindRegisteredProfileResult struct {
	r profiles.Host
}

//Expect specifies that invocation of MandateRegistry.FindRegisteredProfile is expected from 1 to Infinity times
func (m *mMandateRegistryMockFindRegisteredProfile) Expect(p endpoints.Inbound) *mMandateRegistryMockFindRegisteredProfile {
	m.mock.FindRegisteredProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockFindRegisteredProfileExpectation{}
	}
	m.mainExpectation.input = &MandateRegistryMockFindRegisteredProfileInput{p}
	return m
}

//Return specifies results of invocation of MandateRegistry.FindRegisteredProfile
func (m *mMandateRegistryMockFindRegisteredProfile) Return(r profiles.Host) *MandateRegistryMock {
	m.mock.FindRegisteredProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockFindRegisteredProfileExpectation{}
	}
	m.mainExpectation.result = &MandateRegistryMockFindRegisteredProfileResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MandateRegistry.FindRegisteredProfile is expected once
func (m *mMandateRegistryMockFindRegisteredProfile) ExpectOnce(p endpoints.Inbound) *MandateRegistryMockFindRegisteredProfileExpectation {
	m.mock.FindRegisteredProfileFunc = nil
	m.mainExpectation = nil

	expectation := &MandateRegistryMockFindRegisteredProfileExpectation{}
	expectation.input = &MandateRegistryMockFindRegisteredProfileInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MandateRegistryMockFindRegisteredProfileExpectation) Return(r profiles.Host) {
	e.result = &MandateRegistryMockFindRegisteredProfileResult{r}
}

//Set uses given function f as a mock of MandateRegistry.FindRegisteredProfile method
func (m *mMandateRegistryMockFindRegisteredProfile) Set(f func(p endpoints.Inbound) (r profiles.Host)) *MandateRegistryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FindRegisteredProfileFunc = f
	return m.mock
}

//FindRegisteredProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.MandateRegistry interface
func (m *MandateRegistryMock) FindRegisteredProfile(p endpoints.Inbound) (r profiles.Host) {
	counter := atomic.AddUint64(&m.FindRegisteredProfilePreCounter, 1)
	defer atomic.AddUint64(&m.FindRegisteredProfileCounter, 1)

	if len(m.FindRegisteredProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FindRegisteredProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MandateRegistryMock.FindRegisteredProfile. %v", p)
			return
		}

		input := m.FindRegisteredProfileMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MandateRegistryMockFindRegisteredProfileInput{p}, "MandateRegistry.FindRegisteredProfile got unexpected parameters")

		result := m.FindRegisteredProfileMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.FindRegisteredProfile")
			return
		}

		r = result.r

		return
	}

	if m.FindRegisteredProfileMock.mainExpectation != nil {

		input := m.FindRegisteredProfileMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MandateRegistryMockFindRegisteredProfileInput{p}, "MandateRegistry.FindRegisteredProfile got unexpected parameters")
		}

		result := m.FindRegisteredProfileMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.FindRegisteredProfile")
		}

		r = result.r

		return
	}

	if m.FindRegisteredProfileFunc == nil {
		m.t.Fatalf("Unexpected call to MandateRegistryMock.FindRegisteredProfile. %v", p)
		return
	}

	return m.FindRegisteredProfileFunc(p)
}

//FindRegisteredProfileMinimockCounter returns a count of MandateRegistryMock.FindRegisteredProfileFunc invocations
func (m *MandateRegistryMock) FindRegisteredProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FindRegisteredProfileCounter)
}

//FindRegisteredProfileMinimockPreCounter returns the value of MandateRegistryMock.FindRegisteredProfile invocations
func (m *MandateRegistryMock) FindRegisteredProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FindRegisteredProfilePreCounter)
}

//FindRegisteredProfileFinished returns true if mock invocations count is ok
func (m *MandateRegistryMock) FindRegisteredProfileFinished() bool {
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

type mMandateRegistryMockGetCloudIdentity struct {
	mock              *MandateRegistryMock
	mainExpectation   *MandateRegistryMockGetCloudIdentityExpectation
	expectationSeries []*MandateRegistryMockGetCloudIdentityExpectation
}

type MandateRegistryMockGetCloudIdentityExpectation struct {
	result *MandateRegistryMockGetCloudIdentityResult
}

type MandateRegistryMockGetCloudIdentityResult struct {
	r cryptkit.DigestHolder
}

//Expect specifies that invocation of MandateRegistry.GetCloudIdentity is expected from 1 to Infinity times
func (m *mMandateRegistryMockGetCloudIdentity) Expect() *mMandateRegistryMockGetCloudIdentity {
	m.mock.GetCloudIdentityFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockGetCloudIdentityExpectation{}
	}

	return m
}

//Return specifies results of invocation of MandateRegistry.GetCloudIdentity
func (m *mMandateRegistryMockGetCloudIdentity) Return(r cryptkit.DigestHolder) *MandateRegistryMock {
	m.mock.GetCloudIdentityFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockGetCloudIdentityExpectation{}
	}
	m.mainExpectation.result = &MandateRegistryMockGetCloudIdentityResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MandateRegistry.GetCloudIdentity is expected once
func (m *mMandateRegistryMockGetCloudIdentity) ExpectOnce() *MandateRegistryMockGetCloudIdentityExpectation {
	m.mock.GetCloudIdentityFunc = nil
	m.mainExpectation = nil

	expectation := &MandateRegistryMockGetCloudIdentityExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MandateRegistryMockGetCloudIdentityExpectation) Return(r cryptkit.DigestHolder) {
	e.result = &MandateRegistryMockGetCloudIdentityResult{r}
}

//Set uses given function f as a mock of MandateRegistry.GetCloudIdentity method
func (m *mMandateRegistryMockGetCloudIdentity) Set(f func() (r cryptkit.DigestHolder)) *MandateRegistryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCloudIdentityFunc = f
	return m.mock
}

//GetCloudIdentity implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.MandateRegistry interface
func (m *MandateRegistryMock) GetCloudIdentity() (r cryptkit.DigestHolder) {
	counter := atomic.AddUint64(&m.GetCloudIdentityPreCounter, 1)
	defer atomic.AddUint64(&m.GetCloudIdentityCounter, 1)

	if len(m.GetCloudIdentityMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCloudIdentityMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MandateRegistryMock.GetCloudIdentity.")
			return
		}

		result := m.GetCloudIdentityMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.GetCloudIdentity")
			return
		}

		r = result.r

		return
	}

	if m.GetCloudIdentityMock.mainExpectation != nil {

		result := m.GetCloudIdentityMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.GetCloudIdentity")
		}

		r = result.r

		return
	}

	if m.GetCloudIdentityFunc == nil {
		m.t.Fatalf("Unexpected call to MandateRegistryMock.GetCloudIdentity.")
		return
	}

	return m.GetCloudIdentityFunc()
}

//GetCloudIdentityMinimockCounter returns a count of MandateRegistryMock.GetCloudIdentityFunc invocations
func (m *MandateRegistryMock) GetCloudIdentityMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudIdentityCounter)
}

//GetCloudIdentityMinimockPreCounter returns the value of MandateRegistryMock.GetCloudIdentity invocations
func (m *MandateRegistryMock) GetCloudIdentityMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudIdentityPreCounter)
}

//GetCloudIdentityFinished returns true if mock invocations count is ok
func (m *MandateRegistryMock) GetCloudIdentityFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCloudIdentityMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCloudIdentityCounter) == uint64(len(m.GetCloudIdentityMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCloudIdentityMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCloudIdentityCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCloudIdentityFunc != nil {
		return atomic.LoadUint64(&m.GetCloudIdentityCounter) > 0
	}

	return true
}

type mMandateRegistryMockGetConsensusConfiguration struct {
	mock              *MandateRegistryMock
	mainExpectation   *MandateRegistryMockGetConsensusConfigurationExpectation
	expectationSeries []*MandateRegistryMockGetConsensusConfigurationExpectation
}

type MandateRegistryMockGetConsensusConfigurationExpectation struct {
	result *MandateRegistryMockGetConsensusConfigurationResult
}

type MandateRegistryMockGetConsensusConfigurationResult struct {
	r ConsensusConfiguration
}

//Expect specifies that invocation of MandateRegistry.GetConsensusConfiguration is expected from 1 to Infinity times
func (m *mMandateRegistryMockGetConsensusConfiguration) Expect() *mMandateRegistryMockGetConsensusConfiguration {
	m.mock.GetConsensusConfigurationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockGetConsensusConfigurationExpectation{}
	}

	return m
}

//Return specifies results of invocation of MandateRegistry.GetConsensusConfiguration
func (m *mMandateRegistryMockGetConsensusConfiguration) Return(r ConsensusConfiguration) *MandateRegistryMock {
	m.mock.GetConsensusConfigurationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockGetConsensusConfigurationExpectation{}
	}
	m.mainExpectation.result = &MandateRegistryMockGetConsensusConfigurationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MandateRegistry.GetConsensusConfiguration is expected once
func (m *mMandateRegistryMockGetConsensusConfiguration) ExpectOnce() *MandateRegistryMockGetConsensusConfigurationExpectation {
	m.mock.GetConsensusConfigurationFunc = nil
	m.mainExpectation = nil

	expectation := &MandateRegistryMockGetConsensusConfigurationExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MandateRegistryMockGetConsensusConfigurationExpectation) Return(r ConsensusConfiguration) {
	e.result = &MandateRegistryMockGetConsensusConfigurationResult{r}
}

//Set uses given function f as a mock of MandateRegistry.GetConsensusConfiguration method
func (m *mMandateRegistryMockGetConsensusConfiguration) Set(f func() (r ConsensusConfiguration)) *MandateRegistryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetConsensusConfigurationFunc = f
	return m.mock
}

//GetConsensusConfiguration implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.MandateRegistry interface
func (m *MandateRegistryMock) GetConsensusConfiguration() (r ConsensusConfiguration) {
	counter := atomic.AddUint64(&m.GetConsensusConfigurationPreCounter, 1)
	defer atomic.AddUint64(&m.GetConsensusConfigurationCounter, 1)

	if len(m.GetConsensusConfigurationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetConsensusConfigurationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MandateRegistryMock.GetConsensusConfiguration.")
			return
		}

		result := m.GetConsensusConfigurationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.GetConsensusConfiguration")
			return
		}

		r = result.r

		return
	}

	if m.GetConsensusConfigurationMock.mainExpectation != nil {

		result := m.GetConsensusConfigurationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.GetConsensusConfiguration")
		}

		r = result.r

		return
	}

	if m.GetConsensusConfigurationFunc == nil {
		m.t.Fatalf("Unexpected call to MandateRegistryMock.GetConsensusConfiguration.")
		return
	}

	return m.GetConsensusConfigurationFunc()
}

//GetConsensusConfigurationMinimockCounter returns a count of MandateRegistryMock.GetConsensusConfigurationFunc invocations
func (m *MandateRegistryMock) GetConsensusConfigurationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetConsensusConfigurationCounter)
}

//GetConsensusConfigurationMinimockPreCounter returns the value of MandateRegistryMock.GetConsensusConfiguration invocations
func (m *MandateRegistryMock) GetConsensusConfigurationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetConsensusConfigurationPreCounter)
}

//GetConsensusConfigurationFinished returns true if mock invocations count is ok
func (m *MandateRegistryMock) GetConsensusConfigurationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetConsensusConfigurationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetConsensusConfigurationCounter) == uint64(len(m.GetConsensusConfigurationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetConsensusConfigurationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetConsensusConfigurationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetConsensusConfigurationFunc != nil {
		return atomic.LoadUint64(&m.GetConsensusConfigurationCounter) > 0
	}

	return true
}

type mMandateRegistryMockGetPrimingCloudHash struct {
	mock              *MandateRegistryMock
	mainExpectation   *MandateRegistryMockGetPrimingCloudHashExpectation
	expectationSeries []*MandateRegistryMockGetPrimingCloudHashExpectation
}

type MandateRegistryMockGetPrimingCloudHashExpectation struct {
	result *MandateRegistryMockGetPrimingCloudHashResult
}

type MandateRegistryMockGetPrimingCloudHashResult struct {
	r proofs.CloudStateHash
}

//Expect specifies that invocation of MandateRegistry.GetPrimingCloudHash is expected from 1 to Infinity times
func (m *mMandateRegistryMockGetPrimingCloudHash) Expect() *mMandateRegistryMockGetPrimingCloudHash {
	m.mock.GetPrimingCloudHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockGetPrimingCloudHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of MandateRegistry.GetPrimingCloudHash
func (m *mMandateRegistryMockGetPrimingCloudHash) Return(r proofs.CloudStateHash) *MandateRegistryMock {
	m.mock.GetPrimingCloudHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MandateRegistryMockGetPrimingCloudHashExpectation{}
	}
	m.mainExpectation.result = &MandateRegistryMockGetPrimingCloudHashResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MandateRegistry.GetPrimingCloudHash is expected once
func (m *mMandateRegistryMockGetPrimingCloudHash) ExpectOnce() *MandateRegistryMockGetPrimingCloudHashExpectation {
	m.mock.GetPrimingCloudHashFunc = nil
	m.mainExpectation = nil

	expectation := &MandateRegistryMockGetPrimingCloudHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MandateRegistryMockGetPrimingCloudHashExpectation) Return(r proofs.CloudStateHash) {
	e.result = &MandateRegistryMockGetPrimingCloudHashResult{r}
}

//Set uses given function f as a mock of MandateRegistry.GetPrimingCloudHash method
func (m *mMandateRegistryMockGetPrimingCloudHash) Set(f func() (r proofs.CloudStateHash)) *MandateRegistryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrimingCloudHashFunc = f
	return m.mock
}

//GetPrimingCloudHash implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.MandateRegistry interface
func (m *MandateRegistryMock) GetPrimingCloudHash() (r proofs.CloudStateHash) {
	counter := atomic.AddUint64(&m.GetPrimingCloudHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetPrimingCloudHashCounter, 1)

	if len(m.GetPrimingCloudHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrimingCloudHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MandateRegistryMock.GetPrimingCloudHash.")
			return
		}

		result := m.GetPrimingCloudHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.GetPrimingCloudHash")
			return
		}

		r = result.r

		return
	}

	if m.GetPrimingCloudHashMock.mainExpectation != nil {

		result := m.GetPrimingCloudHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MandateRegistryMock.GetPrimingCloudHash")
		}

		r = result.r

		return
	}

	if m.GetPrimingCloudHashFunc == nil {
		m.t.Fatalf("Unexpected call to MandateRegistryMock.GetPrimingCloudHash.")
		return
	}

	return m.GetPrimingCloudHashFunc()
}

//GetPrimingCloudHashMinimockCounter returns a count of MandateRegistryMock.GetPrimingCloudHashFunc invocations
func (m *MandateRegistryMock) GetPrimingCloudHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimingCloudHashCounter)
}

//GetPrimingCloudHashMinimockPreCounter returns the value of MandateRegistryMock.GetPrimingCloudHash invocations
func (m *MandateRegistryMock) GetPrimingCloudHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrimingCloudHashPreCounter)
}

//GetPrimingCloudHashFinished returns true if mock invocations count is ok
func (m *MandateRegistryMock) GetPrimingCloudHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPrimingCloudHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPrimingCloudHashCounter) == uint64(len(m.GetPrimingCloudHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPrimingCloudHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPrimingCloudHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPrimingCloudHashFunc != nil {
		return atomic.LoadUint64(&m.GetPrimingCloudHashCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MandateRegistryMock) ValidateCallCounters() {

	if !m.FindRegisteredProfileFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.FindRegisteredProfile")
	}

	if !m.GetCloudIdentityFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.GetCloudIdentity")
	}

	if !m.GetConsensusConfigurationFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.GetConsensusConfiguration")
	}

	if !m.GetPrimingCloudHashFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.GetPrimingCloudHash")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MandateRegistryMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MandateRegistryMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MandateRegistryMock) MinimockFinish() {

	if !m.FindRegisteredProfileFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.FindRegisteredProfile")
	}

	if !m.GetCloudIdentityFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.GetCloudIdentity")
	}

	if !m.GetConsensusConfigurationFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.GetConsensusConfiguration")
	}

	if !m.GetPrimingCloudHashFinished() {
		m.t.Fatal("Expected call to MandateRegistryMock.GetPrimingCloudHash")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MandateRegistryMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MandateRegistryMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FindRegisteredProfileFinished()
		ok = ok && m.GetCloudIdentityFinished()
		ok = ok && m.GetConsensusConfigurationFinished()
		ok = ok && m.GetPrimingCloudHashFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FindRegisteredProfileFinished() {
				m.t.Error("Expected call to MandateRegistryMock.FindRegisteredProfile")
			}

			if !m.GetCloudIdentityFinished() {
				m.t.Error("Expected call to MandateRegistryMock.GetCloudIdentity")
			}

			if !m.GetConsensusConfigurationFinished() {
				m.t.Error("Expected call to MandateRegistryMock.GetConsensusConfiguration")
			}

			if !m.GetPrimingCloudHashFinished() {
				m.t.Error("Expected call to MandateRegistryMock.GetPrimingCloudHash")
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
func (m *MandateRegistryMock) AllMocksCalled() bool {

	if !m.FindRegisteredProfileFinished() {
		return false
	}

	if !m.GetCloudIdentityFinished() {
		return false
	}

	if !m.GetConsensusConfigurationFinished() {
		return false
	}

	if !m.GetPrimingCloudHashFinished() {
		return false
	}

	return true
}
