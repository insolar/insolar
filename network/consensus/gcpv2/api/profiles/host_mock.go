package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Host" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	endpoints "github.com/insolar/insolar/network/consensus/common/endpoints"

	testify_assert "github.com/stretchr/testify/assert"
)

//HostMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Host
type HostMock struct {
	t minimock.Tester

	GetDefaultEndpointFunc       func() (r endpoints.Outbound)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mHostMockGetDefaultEndpoint

	GetPublicKeyStoreFunc       func() (r cryptkit.PublicKeyStore)
	GetPublicKeyStoreCounter    uint64
	GetPublicKeyStorePreCounter uint64
	GetPublicKeyStoreMock       mHostMockGetPublicKeyStore

	IsAcceptableHostFunc       func(p endpoints.Inbound) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mHostMockIsAcceptableHost
}

//NewHostMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Host
func NewHostMock(t minimock.Tester) *HostMock {
	m := &HostMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetDefaultEndpointMock = mHostMockGetDefaultEndpoint{mock: m}
	m.GetPublicKeyStoreMock = mHostMockGetPublicKeyStore{mock: m}
	m.IsAcceptableHostMock = mHostMockIsAcceptableHost{mock: m}

	return m
}

type mHostMockGetDefaultEndpoint struct {
	mock              *HostMock
	mainExpectation   *HostMockGetDefaultEndpointExpectation
	expectationSeries []*HostMockGetDefaultEndpointExpectation
}

type HostMockGetDefaultEndpointExpectation struct {
	result *HostMockGetDefaultEndpointResult
}

type HostMockGetDefaultEndpointResult struct {
	r endpoints.Outbound
}

//Expect specifies that invocation of Host.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mHostMockGetDefaultEndpoint) Expect() *mHostMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of Host.GetDefaultEndpoint
func (m *mHostMockGetDefaultEndpoint) Return(r endpoints.Outbound) *HostMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &HostMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Host.GetDefaultEndpoint is expected once
func (m *mHostMockGetDefaultEndpoint) ExpectOnce() *HostMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &HostMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostMockGetDefaultEndpointExpectation) Return(r endpoints.Outbound) {
	e.result = &HostMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of Host.GetDefaultEndpoint method
func (m *mHostMockGetDefaultEndpoint) Set(f func() (r endpoints.Outbound)) *HostMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Host interface
func (m *HostMock) GetDefaultEndpoint() (r endpoints.Outbound) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to HostMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of HostMock.GetDefaultEndpointFunc invocations
func (m *HostMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of HostMock.GetDefaultEndpoint invocations
func (m *HostMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *HostMock) GetDefaultEndpointFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDefaultEndpointCounter) == uint64(len(m.GetDefaultEndpointMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDefaultEndpointMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDefaultEndpointCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDefaultEndpointFunc != nil {
		return atomic.LoadUint64(&m.GetDefaultEndpointCounter) > 0
	}

	return true
}

type mHostMockGetPublicKeyStore struct {
	mock              *HostMock
	mainExpectation   *HostMockGetPublicKeyStoreExpectation
	expectationSeries []*HostMockGetPublicKeyStoreExpectation
}

type HostMockGetPublicKeyStoreExpectation struct {
	result *HostMockGetPublicKeyStoreResult
}

type HostMockGetPublicKeyStoreResult struct {
	r cryptkit.PublicKeyStore
}

//Expect specifies that invocation of Host.CreatePublicKeyStore is expected from 1 to Infinity times
func (m *mHostMockGetPublicKeyStore) Expect() *mHostMockGetPublicKeyStore {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostMockGetPublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of Host.CreatePublicKeyStore
func (m *mHostMockGetPublicKeyStore) Return(r cryptkit.PublicKeyStore) *HostMock {
	m.mock.GetPublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostMockGetPublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &HostMockGetPublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Host.CreatePublicKeyStore is expected once
func (m *mHostMockGetPublicKeyStore) ExpectOnce() *HostMockGetPublicKeyStoreExpectation {
	m.mock.GetPublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &HostMockGetPublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostMockGetPublicKeyStoreExpectation) Return(r cryptkit.PublicKeyStore) {
	e.result = &HostMockGetPublicKeyStoreResult{r}
}

//Set uses given function f as a mock of Host.CreatePublicKeyStore method
func (m *mHostMockGetPublicKeyStore) Set(f func() (r cryptkit.PublicKeyStore)) *HostMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyStoreFunc = f
	return m.mock
}

//CreatePublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Host interface
func (m *HostMock) GetPublicKeyStore() (r cryptkit.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetPublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyStoreCounter, 1)

	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostMock.CreatePublicKeyStore.")
			return
		}

		result := m.GetPublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostMock.CreatePublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreMock.mainExpectation != nil {

		result := m.GetPublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostMock.CreatePublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to HostMock.CreatePublicKeyStore.")
		return
	}

	return m.GetPublicKeyStoreFunc()
}

//GetPublicKeyStoreMinimockCounter returns a count of HostMock.GetPublicKeyStoreFunc invocations
func (m *HostMock) GetPublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStoreCounter)
}

//GetPublicKeyStoreMinimockPreCounter returns the value of HostMock.CreatePublicKeyStore invocations
func (m *HostMock) GetPublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyStorePreCounter)
}

//GetPublicKeyStoreFinished returns true if mock invocations count is ok
func (m *HostMock) GetPublicKeyStoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPublicKeyStoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) == uint64(len(m.GetPublicKeyStoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPublicKeyStoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPublicKeyStoreFunc != nil {
		return atomic.LoadUint64(&m.GetPublicKeyStoreCounter) > 0
	}

	return true
}

type mHostMockIsAcceptableHost struct {
	mock              *HostMock
	mainExpectation   *HostMockIsAcceptableHostExpectation
	expectationSeries []*HostMockIsAcceptableHostExpectation
}

type HostMockIsAcceptableHostExpectation struct {
	input  *HostMockIsAcceptableHostInput
	result *HostMockIsAcceptableHostResult
}

type HostMockIsAcceptableHostInput struct {
	p endpoints.Inbound
}

type HostMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of Host.IsAcceptableHost is expected from 1 to Infinity times
func (m *mHostMockIsAcceptableHost) Expect(p endpoints.Inbound) *mHostMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &HostMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of Host.IsAcceptableHost
func (m *mHostMockIsAcceptableHost) Return(r bool) *HostMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &HostMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Host.IsAcceptableHost is expected once
func (m *mHostMockIsAcceptableHost) ExpectOnce(p endpoints.Inbound) *HostMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &HostMockIsAcceptableHostExpectation{}
	expectation.input = &HostMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &HostMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of Host.IsAcceptableHost method
func (m *mHostMockIsAcceptableHost) Set(f func(p endpoints.Inbound) (r bool)) *HostMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Host interface
func (m *HostMock) IsAcceptableHost(p endpoints.Inbound) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostMockIsAcceptableHostInput{p}, "Host.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostMockIsAcceptableHostInput{p}, "Host.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to HostMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of HostMock.IsAcceptableHostFunc invocations
func (m *HostMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of HostMock.IsAcceptableHost invocations
func (m *HostMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *HostMock) IsAcceptableHostFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsAcceptableHostCounter) == uint64(len(m.IsAcceptableHostMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsAcceptableHostMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsAcceptableHostCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsAcceptableHostFunc != nil {
		return atomic.LoadUint64(&m.IsAcceptableHostCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HostMock) ValidateCallCounters() {

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to HostMock.GetDefaultEndpoint")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to HostMock.CreatePublicKeyStore")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to HostMock.IsAcceptableHost")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HostMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *HostMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *HostMock) MinimockFinish() {

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to HostMock.GetDefaultEndpoint")
	}

	if !m.GetPublicKeyStoreFinished() {
		m.t.Fatal("Expected call to HostMock.CreatePublicKeyStore")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to HostMock.IsAcceptableHost")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *HostMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *HostMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetPublicKeyStoreFinished()
		ok = ok && m.IsAcceptableHostFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to HostMock.GetDefaultEndpoint")
			}

			if !m.GetPublicKeyStoreFinished() {
				m.t.Error("Expected call to HostMock.CreatePublicKeyStore")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to HostMock.IsAcceptableHost")
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
func (m *HostMock) AllMocksCalled() bool {

	if !m.GetDefaultEndpointFinished() {
		return false
	}

	if !m.GetPublicKeyStoreFinished() {
		return false
	}

	if !m.IsAcceptableHostFinished() {
		return false
	}

	return true
}
