package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "HostProfile" can be found in github.com/insolar/insolar/network/consensus/gcpv2/common
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	common "github.com/insolar/insolar/network/consensus/common"

	testify_assert "github.com/stretchr/testify/assert"
)

//HostProfileMock implements github.com/insolar/insolar/network/consensus/gcpv2/common.HostProfile
type HostProfileMock struct {
	t minimock.Tester

	GetDefaultEndpointFunc       func() (r common.NodeEndpoint)
	GetDefaultEndpointCounter    uint64
	GetDefaultEndpointPreCounter uint64
	GetDefaultEndpointMock       mHostProfileMockGetDefaultEndpoint

	GetNodePublicKeyStoreFunc       func() (r common.PublicKeyStore)
	GetNodePublicKeyStoreCounter    uint64
	GetNodePublicKeyStorePreCounter uint64
	GetNodePublicKeyStoreMock       mHostProfileMockGetNodePublicKeyStore

	IsAcceptableHostFunc       func(p common.HostIdentityHolder) (r bool)
	IsAcceptableHostCounter    uint64
	IsAcceptableHostPreCounter uint64
	IsAcceptableHostMock       mHostProfileMockIsAcceptableHost
}

//NewHostProfileMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/common.HostProfile
func NewHostProfileMock(t minimock.Tester) *HostProfileMock {
	m := &HostProfileMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetDefaultEndpointMock = mHostProfileMockGetDefaultEndpoint{mock: m}
	m.GetNodePublicKeyStoreMock = mHostProfileMockGetNodePublicKeyStore{mock: m}
	m.IsAcceptableHostMock = mHostProfileMockIsAcceptableHost{mock: m}

	return m
}

type mHostProfileMockGetDefaultEndpoint struct {
	mock              *HostProfileMock
	mainExpectation   *HostProfileMockGetDefaultEndpointExpectation
	expectationSeries []*HostProfileMockGetDefaultEndpointExpectation
}

type HostProfileMockGetDefaultEndpointExpectation struct {
	result *HostProfileMockGetDefaultEndpointResult
}

type HostProfileMockGetDefaultEndpointResult struct {
	r common.NodeEndpoint
}

//Expect specifies that invocation of HostProfile.GetDefaultEndpoint is expected from 1 to Infinity times
func (m *mHostProfileMockGetDefaultEndpoint) Expect() *mHostProfileMockGetDefaultEndpoint {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostProfileMockGetDefaultEndpointExpectation{}
	}

	return m
}

//Return specifies results of invocation of HostProfile.GetDefaultEndpoint
func (m *mHostProfileMockGetDefaultEndpoint) Return(r common.NodeEndpoint) *HostProfileMock {
	m.mock.GetDefaultEndpointFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostProfileMockGetDefaultEndpointExpectation{}
	}
	m.mainExpectation.result = &HostProfileMockGetDefaultEndpointResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostProfile.GetDefaultEndpoint is expected once
func (m *mHostProfileMockGetDefaultEndpoint) ExpectOnce() *HostProfileMockGetDefaultEndpointExpectation {
	m.mock.GetDefaultEndpointFunc = nil
	m.mainExpectation = nil

	expectation := &HostProfileMockGetDefaultEndpointExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostProfileMockGetDefaultEndpointExpectation) Return(r common.NodeEndpoint) {
	e.result = &HostProfileMockGetDefaultEndpointResult{r}
}

//Set uses given function f as a mock of HostProfile.GetDefaultEndpoint method
func (m *mHostProfileMockGetDefaultEndpoint) Set(f func() (r common.NodeEndpoint)) *HostProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDefaultEndpointFunc = f
	return m.mock
}

//GetDefaultEndpoint implements github.com/insolar/insolar/network/consensus/gcpv2/common.HostProfile interface
func (m *HostProfileMock) GetDefaultEndpoint() (r common.NodeEndpoint) {
	counter := atomic.AddUint64(&m.GetDefaultEndpointPreCounter, 1)
	defer atomic.AddUint64(&m.GetDefaultEndpointCounter, 1)

	if len(m.GetDefaultEndpointMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDefaultEndpointMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostProfileMock.GetDefaultEndpoint.")
			return
		}

		result := m.GetDefaultEndpointMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostProfileMock.GetDefaultEndpoint")
			return
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointMock.mainExpectation != nil {

		result := m.GetDefaultEndpointMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostProfileMock.GetDefaultEndpoint")
		}

		r = result.r

		return
	}

	if m.GetDefaultEndpointFunc == nil {
		m.t.Fatalf("Unexpected call to HostProfileMock.GetDefaultEndpoint.")
		return
	}

	return m.GetDefaultEndpointFunc()
}

//GetDefaultEndpointMinimockCounter returns a count of HostProfileMock.GetDefaultEndpointFunc invocations
func (m *HostProfileMock) GetDefaultEndpointMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointCounter)
}

//GetDefaultEndpointMinimockPreCounter returns the value of HostProfileMock.GetDefaultEndpoint invocations
func (m *HostProfileMock) GetDefaultEndpointMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDefaultEndpointPreCounter)
}

//GetDefaultEndpointFinished returns true if mock invocations count is ok
func (m *HostProfileMock) GetDefaultEndpointFinished() bool {
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

type mHostProfileMockGetNodePublicKeyStore struct {
	mock              *HostProfileMock
	mainExpectation   *HostProfileMockGetNodePublicKeyStoreExpectation
	expectationSeries []*HostProfileMockGetNodePublicKeyStoreExpectation
}

type HostProfileMockGetNodePublicKeyStoreExpectation struct {
	result *HostProfileMockGetNodePublicKeyStoreResult
}

type HostProfileMockGetNodePublicKeyStoreResult struct {
	r common.PublicKeyStore
}

//Expect specifies that invocation of HostProfile.GetNodePublicKeyStore is expected from 1 to Infinity times
func (m *mHostProfileMockGetNodePublicKeyStore) Expect() *mHostProfileMockGetNodePublicKeyStore {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostProfileMockGetNodePublicKeyStoreExpectation{}
	}

	return m
}

//Return specifies results of invocation of HostProfile.GetNodePublicKeyStore
func (m *mHostProfileMockGetNodePublicKeyStore) Return(r common.PublicKeyStore) *HostProfileMock {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostProfileMockGetNodePublicKeyStoreExpectation{}
	}
	m.mainExpectation.result = &HostProfileMockGetNodePublicKeyStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostProfile.GetNodePublicKeyStore is expected once
func (m *mHostProfileMockGetNodePublicKeyStore) ExpectOnce() *HostProfileMockGetNodePublicKeyStoreExpectation {
	m.mock.GetNodePublicKeyStoreFunc = nil
	m.mainExpectation = nil

	expectation := &HostProfileMockGetNodePublicKeyStoreExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostProfileMockGetNodePublicKeyStoreExpectation) Return(r common.PublicKeyStore) {
	e.result = &HostProfileMockGetNodePublicKeyStoreResult{r}
}

//Set uses given function f as a mock of HostProfile.GetNodePublicKeyStore method
func (m *mHostProfileMockGetNodePublicKeyStore) Set(f func() (r common.PublicKeyStore)) *HostProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodePublicKeyStoreFunc = f
	return m.mock
}

//GetNodePublicKeyStore implements github.com/insolar/insolar/network/consensus/gcpv2/common.HostProfile interface
func (m *HostProfileMock) GetNodePublicKeyStore() (r common.PublicKeyStore) {
	counter := atomic.AddUint64(&m.GetNodePublicKeyStorePreCounter, 1)
	defer atomic.AddUint64(&m.GetNodePublicKeyStoreCounter, 1)

	if len(m.GetNodePublicKeyStoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodePublicKeyStoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostProfileMock.GetNodePublicKeyStore.")
			return
		}

		result := m.GetNodePublicKeyStoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostProfileMock.GetNodePublicKeyStore")
			return
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreMock.mainExpectation != nil {

		result := m.GetNodePublicKeyStoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostProfileMock.GetNodePublicKeyStore")
		}

		r = result.r

		return
	}

	if m.GetNodePublicKeyStoreFunc == nil {
		m.t.Fatalf("Unexpected call to HostProfileMock.GetNodePublicKeyStore.")
		return
	}

	return m.GetNodePublicKeyStoreFunc()
}

//GetNodePublicKeyStoreMinimockCounter returns a count of HostProfileMock.GetNodePublicKeyStoreFunc invocations
func (m *HostProfileMock) GetNodePublicKeyStoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter)
}

//GetNodePublicKeyStoreMinimockPreCounter returns the value of HostProfileMock.GetNodePublicKeyStore invocations
func (m *HostProfileMock) GetNodePublicKeyStoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodePublicKeyStorePreCounter)
}

//GetNodePublicKeyStoreFinished returns true if mock invocations count is ok
func (m *HostProfileMock) GetNodePublicKeyStoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodePublicKeyStoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter) == uint64(len(m.GetNodePublicKeyStoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodePublicKeyStoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodePublicKeyStoreFunc != nil {
		return atomic.LoadUint64(&m.GetNodePublicKeyStoreCounter) > 0
	}

	return true
}

type mHostProfileMockIsAcceptableHost struct {
	mock              *HostProfileMock
	mainExpectation   *HostProfileMockIsAcceptableHostExpectation
	expectationSeries []*HostProfileMockIsAcceptableHostExpectation
}

type HostProfileMockIsAcceptableHostExpectation struct {
	input  *HostProfileMockIsAcceptableHostInput
	result *HostProfileMockIsAcceptableHostResult
}

type HostProfileMockIsAcceptableHostInput struct {
	p common.HostIdentityHolder
}

type HostProfileMockIsAcceptableHostResult struct {
	r bool
}

//Expect specifies that invocation of HostProfile.IsAcceptableHost is expected from 1 to Infinity times
func (m *mHostProfileMockIsAcceptableHost) Expect(p common.HostIdentityHolder) *mHostProfileMockIsAcceptableHost {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.input = &HostProfileMockIsAcceptableHostInput{p}
	return m
}

//Return specifies results of invocation of HostProfile.IsAcceptableHost
func (m *mHostProfileMockIsAcceptableHost) Return(r bool) *HostProfileMock {
	m.mock.IsAcceptableHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HostProfileMockIsAcceptableHostExpectation{}
	}
	m.mainExpectation.result = &HostProfileMockIsAcceptableHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HostProfile.IsAcceptableHost is expected once
func (m *mHostProfileMockIsAcceptableHost) ExpectOnce(p common.HostIdentityHolder) *HostProfileMockIsAcceptableHostExpectation {
	m.mock.IsAcceptableHostFunc = nil
	m.mainExpectation = nil

	expectation := &HostProfileMockIsAcceptableHostExpectation{}
	expectation.input = &HostProfileMockIsAcceptableHostInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HostProfileMockIsAcceptableHostExpectation) Return(r bool) {
	e.result = &HostProfileMockIsAcceptableHostResult{r}
}

//Set uses given function f as a mock of HostProfile.IsAcceptableHost method
func (m *mHostProfileMockIsAcceptableHost) Set(f func(p common.HostIdentityHolder) (r bool)) *HostProfileMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAcceptableHostFunc = f
	return m.mock
}

//IsAcceptableHost implements github.com/insolar/insolar/network/consensus/gcpv2/common.HostProfile interface
func (m *HostProfileMock) IsAcceptableHost(p common.HostIdentityHolder) (r bool) {
	counter := atomic.AddUint64(&m.IsAcceptableHostPreCounter, 1)
	defer atomic.AddUint64(&m.IsAcceptableHostCounter, 1)

	if len(m.IsAcceptableHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAcceptableHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HostProfileMock.IsAcceptableHost. %v", p)
			return
		}

		input := m.IsAcceptableHostMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HostProfileMockIsAcceptableHostInput{p}, "HostProfile.IsAcceptableHost got unexpected parameters")

		result := m.IsAcceptableHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HostProfileMock.IsAcceptableHost")
			return
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostMock.mainExpectation != nil {

		input := m.IsAcceptableHostMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HostProfileMockIsAcceptableHostInput{p}, "HostProfile.IsAcceptableHost got unexpected parameters")
		}

		result := m.IsAcceptableHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HostProfileMock.IsAcceptableHost")
		}

		r = result.r

		return
	}

	if m.IsAcceptableHostFunc == nil {
		m.t.Fatalf("Unexpected call to HostProfileMock.IsAcceptableHost. %v", p)
		return
	}

	return m.IsAcceptableHostFunc(p)
}

//IsAcceptableHostMinimockCounter returns a count of HostProfileMock.IsAcceptableHostFunc invocations
func (m *HostProfileMock) IsAcceptableHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostCounter)
}

//IsAcceptableHostMinimockPreCounter returns the value of HostProfileMock.IsAcceptableHost invocations
func (m *HostProfileMock) IsAcceptableHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAcceptableHostPreCounter)
}

//IsAcceptableHostFinished returns true if mock invocations count is ok
func (m *HostProfileMock) IsAcceptableHostFinished() bool {
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
func (m *HostProfileMock) ValidateCallCounters() {

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to HostProfileMock.GetDefaultEndpoint")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to HostProfileMock.GetNodePublicKeyStore")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to HostProfileMock.IsAcceptableHost")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HostProfileMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *HostProfileMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *HostProfileMock) MinimockFinish() {

	if !m.GetDefaultEndpointFinished() {
		m.t.Fatal("Expected call to HostProfileMock.GetDefaultEndpoint")
	}

	if !m.GetNodePublicKeyStoreFinished() {
		m.t.Fatal("Expected call to HostProfileMock.GetNodePublicKeyStore")
	}

	if !m.IsAcceptableHostFinished() {
		m.t.Fatal("Expected call to HostProfileMock.IsAcceptableHost")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *HostProfileMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *HostProfileMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetDefaultEndpointFinished()
		ok = ok && m.GetNodePublicKeyStoreFinished()
		ok = ok && m.IsAcceptableHostFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetDefaultEndpointFinished() {
				m.t.Error("Expected call to HostProfileMock.GetDefaultEndpoint")
			}

			if !m.GetNodePublicKeyStoreFinished() {
				m.t.Error("Expected call to HostProfileMock.GetNodePublicKeyStore")
			}

			if !m.IsAcceptableHostFinished() {
				m.t.Error("Expected call to HostProfileMock.IsAcceptableHost")
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
func (m *HostProfileMock) AllMocksCalled() bool {

	if !m.GetDefaultEndpointFinished() {
		return false
	}

	if !m.GetNodePublicKeyStoreFinished() {
		return false
	}

	if !m.IsAcceptableHostFinished() {
		return false
	}

	return true
}
