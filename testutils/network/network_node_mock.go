package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NetworkNode" can be found in github.com/insolar/insolar/insolar
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	keys "github.com/insolar/insolar/platformpolicy/keys"
)

//NetworkNodeMock implements github.com/insolar/insolar/insolar.NetworkNode
type NetworkNodeMock struct {
	t minimock.Tester

	AddressFunc       func() (r string)
	AddressCounter    uint64
	AddressPreCounter uint64
	AddressMock       mNetworkNodeMockAddress

	GetGlobuleIDFunc       func() (r insolar.GlobuleID)
	GetGlobuleIDCounter    uint64
	GetGlobuleIDPreCounter uint64
	GetGlobuleIDMock       mNetworkNodeMockGetGlobuleID

	GetStateFunc       func() (r insolar.NodeState)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mNetworkNodeMockGetState

	IDFunc       func() (r insolar.Reference)
	IDCounter    uint64
	IDPreCounter uint64
	IDMock       mNetworkNodeMockID

	LeavingETAFunc       func() (r insolar.PulseNumber)
	LeavingETACounter    uint64
	LeavingETAPreCounter uint64
	LeavingETAMock       mNetworkNodeMockLeavingETA

	PublicKeyFunc       func() (r keys.PublicKey)
	PublicKeyCounter    uint64
	PublicKeyPreCounter uint64
	PublicKeyMock       mNetworkNodeMockPublicKey

	RoleFunc       func() (r insolar.StaticRole)
	RoleCounter    uint64
	RolePreCounter uint64
	RoleMock       mNetworkNodeMockRole

	ShortIDFunc       func() (r insolar.ShortNodeID)
	ShortIDCounter    uint64
	ShortIDPreCounter uint64
	ShortIDMock       mNetworkNodeMockShortID

	VersionFunc       func() (r string)
	VersionCounter    uint64
	VersionPreCounter uint64
	VersionMock       mNetworkNodeMockVersion
}

//NewNetworkNodeMock returns a mock for github.com/insolar/insolar/insolar.NetworkNode
func NewNetworkNodeMock(t minimock.Tester) *NetworkNodeMock {
	m := &NetworkNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddressMock = mNetworkNodeMockAddress{mock: m}
	m.GetGlobuleIDMock = mNetworkNodeMockGetGlobuleID{mock: m}
	m.GetStateMock = mNetworkNodeMockGetState{mock: m}
	m.IDMock = mNetworkNodeMockID{mock: m}
	m.LeavingETAMock = mNetworkNodeMockLeavingETA{mock: m}
	m.PublicKeyMock = mNetworkNodeMockPublicKey{mock: m}
	m.RoleMock = mNetworkNodeMockRole{mock: m}
	m.ShortIDMock = mNetworkNodeMockShortID{mock: m}
	m.VersionMock = mNetworkNodeMockVersion{mock: m}

	return m
}

type mNetworkNodeMockAddress struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockAddressExpectation
	expectationSeries []*NetworkNodeMockAddressExpectation
}

type NetworkNodeMockAddressExpectation struct {
	result *NetworkNodeMockAddressResult
}

type NetworkNodeMockAddressResult struct {
	r string
}

//Expect specifies that invocation of NetworkNode.Address is expected from 1 to Infinity times
func (m *mNetworkNodeMockAddress) Expect() *mNetworkNodeMockAddress {
	m.mock.AddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.Address
func (m *mNetworkNodeMockAddress) Return(r string) *NetworkNodeMock {
	m.mock.AddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockAddressExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.Address is expected once
func (m *mNetworkNodeMockAddress) ExpectOnce() *NetworkNodeMockAddressExpectation {
	m.mock.AddressFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockAddressExpectation) Return(r string) {
	e.result = &NetworkNodeMockAddressResult{r}
}

//Set uses given function f as a mock of NetworkNode.Address method
func (m *mNetworkNodeMockAddress) Set(f func() (r string)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddressFunc = f
	return m.mock
}

//Address implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) Address() (r string) {
	counter := atomic.AddUint64(&m.AddressPreCounter, 1)
	defer atomic.AddUint64(&m.AddressCounter, 1)

	if len(m.AddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.Address.")
			return
		}

		result := m.AddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.Address")
			return
		}

		r = result.r

		return
	}

	if m.AddressMock.mainExpectation != nil {

		result := m.AddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.Address")
		}

		r = result.r

		return
	}

	if m.AddressFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.Address.")
		return
	}

	return m.AddressFunc()
}

//AddressMinimockCounter returns a count of NetworkNodeMock.AddressFunc invocations
func (m *NetworkNodeMock) AddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddressCounter)
}

//AddressMinimockPreCounter returns the value of NetworkNodeMock.Address invocations
func (m *NetworkNodeMock) AddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddressPreCounter)
}

//AddressFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) AddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddressCounter) == uint64(len(m.AddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddressFunc != nil {
		return atomic.LoadUint64(&m.AddressCounter) > 0
	}

	return true
}

type mNetworkNodeMockGetGlobuleID struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockGetGlobuleIDExpectation
	expectationSeries []*NetworkNodeMockGetGlobuleIDExpectation
}

type NetworkNodeMockGetGlobuleIDExpectation struct {
	result *NetworkNodeMockGetGlobuleIDResult
}

type NetworkNodeMockGetGlobuleIDResult struct {
	r insolar.GlobuleID
}

//Expect specifies that invocation of NetworkNode.GetGlobuleID is expected from 1 to Infinity times
func (m *mNetworkNodeMockGetGlobuleID) Expect() *mNetworkNodeMockGetGlobuleID {
	m.mock.GetGlobuleIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockGetGlobuleIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.GetGlobuleID
func (m *mNetworkNodeMockGetGlobuleID) Return(r insolar.GlobuleID) *NetworkNodeMock {
	m.mock.GetGlobuleIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockGetGlobuleIDExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockGetGlobuleIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.GetGlobuleID is expected once
func (m *mNetworkNodeMockGetGlobuleID) ExpectOnce() *NetworkNodeMockGetGlobuleIDExpectation {
	m.mock.GetGlobuleIDFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockGetGlobuleIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockGetGlobuleIDExpectation) Return(r insolar.GlobuleID) {
	e.result = &NetworkNodeMockGetGlobuleIDResult{r}
}

//Set uses given function f as a mock of NetworkNode.GetGlobuleID method
func (m *mNetworkNodeMockGetGlobuleID) Set(f func() (r insolar.GlobuleID)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetGlobuleIDFunc = f
	return m.mock
}

//GetGlobuleID implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) GetGlobuleID() (r insolar.GlobuleID) {
	counter := atomic.AddUint64(&m.GetGlobuleIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobuleIDCounter, 1)

	if len(m.GetGlobuleIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetGlobuleIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.GetGlobuleID.")
			return
		}

		result := m.GetGlobuleIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.GetGlobuleID")
			return
		}

		r = result.r

		return
	}

	if m.GetGlobuleIDMock.mainExpectation != nil {

		result := m.GetGlobuleIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.GetGlobuleID")
		}

		r = result.r

		return
	}

	if m.GetGlobuleIDFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.GetGlobuleID.")
		return
	}

	return m.GetGlobuleIDFunc()
}

//GetGlobuleIDMinimockCounter returns a count of NetworkNodeMock.GetGlobuleIDFunc invocations
func (m *NetworkNodeMock) GetGlobuleIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleIDCounter)
}

//GetGlobuleIDMinimockPreCounter returns the value of NetworkNodeMock.GetGlobuleID invocations
func (m *NetworkNodeMock) GetGlobuleIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleIDPreCounter)
}

//GetGlobuleIDFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) GetGlobuleIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetGlobuleIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetGlobuleIDCounter) == uint64(len(m.GetGlobuleIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetGlobuleIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetGlobuleIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetGlobuleIDFunc != nil {
		return atomic.LoadUint64(&m.GetGlobuleIDCounter) > 0
	}

	return true
}

type mNetworkNodeMockGetState struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockGetStateExpectation
	expectationSeries []*NetworkNodeMockGetStateExpectation
}

type NetworkNodeMockGetStateExpectation struct {
	result *NetworkNodeMockGetStateResult
}

type NetworkNodeMockGetStateResult struct {
	r insolar.NodeState
}

//Expect specifies that invocation of NetworkNode.GetState is expected from 1 to Infinity times
func (m *mNetworkNodeMockGetState) Expect() *mNetworkNodeMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.GetState
func (m *mNetworkNodeMockGetState) Return(r insolar.NodeState) *NetworkNodeMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockGetStateExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.GetState is expected once
func (m *mNetworkNodeMockGetState) ExpectOnce() *NetworkNodeMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockGetStateExpectation) Return(r insolar.NodeState) {
	e.result = &NetworkNodeMockGetStateResult{r}
}

//Set uses given function f as a mock of NetworkNode.GetState method
func (m *mNetworkNodeMockGetState) Set(f func() (r insolar.NodeState)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) GetState() (r insolar.NodeState) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of NetworkNodeMock.GetStateFunc invocations
func (m *NetworkNodeMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of NetworkNodeMock.GetState invocations
func (m *NetworkNodeMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) GetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStateCounter) == uint64(len(m.GetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStateFunc != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	return true
}

type mNetworkNodeMockID struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockIDExpectation
	expectationSeries []*NetworkNodeMockIDExpectation
}

type NetworkNodeMockIDExpectation struct {
	result *NetworkNodeMockIDResult
}

type NetworkNodeMockIDResult struct {
	r insolar.Reference
}

//Expect specifies that invocation of NetworkNode.ID is expected from 1 to Infinity times
func (m *mNetworkNodeMockID) Expect() *mNetworkNodeMockID {
	m.mock.IDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.ID
func (m *mNetworkNodeMockID) Return(r insolar.Reference) *NetworkNodeMock {
	m.mock.IDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockIDExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.ID is expected once
func (m *mNetworkNodeMockID) ExpectOnce() *NetworkNodeMockIDExpectation {
	m.mock.IDFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockIDExpectation) Return(r insolar.Reference) {
	e.result = &NetworkNodeMockIDResult{r}
}

//Set uses given function f as a mock of NetworkNode.ID method
func (m *mNetworkNodeMockID) Set(f func() (r insolar.Reference)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IDFunc = f
	return m.mock
}

//ID implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) ID() (r insolar.Reference) {
	counter := atomic.AddUint64(&m.IDPreCounter, 1)
	defer atomic.AddUint64(&m.IDCounter, 1)

	if len(m.IDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.ID.")
			return
		}

		result := m.IDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.ID")
			return
		}

		r = result.r

		return
	}

	if m.IDMock.mainExpectation != nil {

		result := m.IDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.ID")
		}

		r = result.r

		return
	}

	if m.IDFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.ID.")
		return
	}

	return m.IDFunc()
}

//IDMinimockCounter returns a count of NetworkNodeMock.IDFunc invocations
func (m *NetworkNodeMock) IDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IDCounter)
}

//IDMinimockPreCounter returns the value of NetworkNodeMock.ID invocations
func (m *NetworkNodeMock) IDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IDPreCounter)
}

//IDFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) IDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IDCounter) == uint64(len(m.IDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IDFunc != nil {
		return atomic.LoadUint64(&m.IDCounter) > 0
	}

	return true
}

type mNetworkNodeMockLeavingETA struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockLeavingETAExpectation
	expectationSeries []*NetworkNodeMockLeavingETAExpectation
}

type NetworkNodeMockLeavingETAExpectation struct {
	result *NetworkNodeMockLeavingETAResult
}

type NetworkNodeMockLeavingETAResult struct {
	r insolar.PulseNumber
}

//Expect specifies that invocation of NetworkNode.LeavingETA is expected from 1 to Infinity times
func (m *mNetworkNodeMockLeavingETA) Expect() *mNetworkNodeMockLeavingETA {
	m.mock.LeavingETAFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockLeavingETAExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.LeavingETA
func (m *mNetworkNodeMockLeavingETA) Return(r insolar.PulseNumber) *NetworkNodeMock {
	m.mock.LeavingETAFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockLeavingETAExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockLeavingETAResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.LeavingETA is expected once
func (m *mNetworkNodeMockLeavingETA) ExpectOnce() *NetworkNodeMockLeavingETAExpectation {
	m.mock.LeavingETAFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockLeavingETAExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockLeavingETAExpectation) Return(r insolar.PulseNumber) {
	e.result = &NetworkNodeMockLeavingETAResult{r}
}

//Set uses given function f as a mock of NetworkNode.LeavingETA method
func (m *mNetworkNodeMockLeavingETA) Set(f func() (r insolar.PulseNumber)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LeavingETAFunc = f
	return m.mock
}

//LeavingETA implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) LeavingETA() (r insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.LeavingETAPreCounter, 1)
	defer atomic.AddUint64(&m.LeavingETACounter, 1)

	if len(m.LeavingETAMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LeavingETAMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.LeavingETA.")
			return
		}

		result := m.LeavingETAMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.LeavingETA")
			return
		}

		r = result.r

		return
	}

	if m.LeavingETAMock.mainExpectation != nil {

		result := m.LeavingETAMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.LeavingETA")
		}

		r = result.r

		return
	}

	if m.LeavingETAFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.LeavingETA.")
		return
	}

	return m.LeavingETAFunc()
}

//LeavingETAMinimockCounter returns a count of NetworkNodeMock.LeavingETAFunc invocations
func (m *NetworkNodeMock) LeavingETAMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LeavingETACounter)
}

//LeavingETAMinimockPreCounter returns the value of NetworkNodeMock.LeavingETA invocations
func (m *NetworkNodeMock) LeavingETAMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LeavingETAPreCounter)
}

//LeavingETAFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) LeavingETAFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LeavingETAMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LeavingETACounter) == uint64(len(m.LeavingETAMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LeavingETAMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LeavingETACounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LeavingETAFunc != nil {
		return atomic.LoadUint64(&m.LeavingETACounter) > 0
	}

	return true
}

type mNetworkNodeMockPublicKey struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockPublicKeyExpectation
	expectationSeries []*NetworkNodeMockPublicKeyExpectation
}

type NetworkNodeMockPublicKeyExpectation struct {
	result *NetworkNodeMockPublicKeyResult
}

type NetworkNodeMockPublicKeyResult struct {
	r keys.PublicKey
}

//Expect specifies that invocation of NetworkNode.PublicKey is expected from 1 to Infinity times
func (m *mNetworkNodeMockPublicKey) Expect() *mNetworkNodeMockPublicKey {
	m.mock.PublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockPublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.PublicKey
func (m *mNetworkNodeMockPublicKey) Return(r keys.PublicKey) *NetworkNodeMock {
	m.mock.PublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockPublicKeyExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockPublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.PublicKey is expected once
func (m *mNetworkNodeMockPublicKey) ExpectOnce() *NetworkNodeMockPublicKeyExpectation {
	m.mock.PublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockPublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockPublicKeyExpectation) Return(r keys.PublicKey) {
	e.result = &NetworkNodeMockPublicKeyResult{r}
}

//Set uses given function f as a mock of NetworkNode.PublicKey method
func (m *mNetworkNodeMockPublicKey) Set(f func() (r keys.PublicKey)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PublicKeyFunc = f
	return m.mock
}

//PublicKey implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) PublicKey() (r keys.PublicKey) {
	counter := atomic.AddUint64(&m.PublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.PublicKeyCounter, 1)

	if len(m.PublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.PublicKey.")
			return
		}

		result := m.PublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.PublicKey")
			return
		}

		r = result.r

		return
	}

	if m.PublicKeyMock.mainExpectation != nil {

		result := m.PublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.PublicKey")
		}

		r = result.r

		return
	}

	if m.PublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.PublicKey.")
		return
	}

	return m.PublicKeyFunc()
}

//PublicKeyMinimockCounter returns a count of NetworkNodeMock.PublicKeyFunc invocations
func (m *NetworkNodeMock) PublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PublicKeyCounter)
}

//PublicKeyMinimockPreCounter returns the value of NetworkNodeMock.PublicKey invocations
func (m *NetworkNodeMock) PublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PublicKeyPreCounter)
}

//PublicKeyFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) PublicKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PublicKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PublicKeyCounter) == uint64(len(m.PublicKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PublicKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PublicKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PublicKeyFunc != nil {
		return atomic.LoadUint64(&m.PublicKeyCounter) > 0
	}

	return true
}

type mNetworkNodeMockRole struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockRoleExpectation
	expectationSeries []*NetworkNodeMockRoleExpectation
}

type NetworkNodeMockRoleExpectation struct {
	result *NetworkNodeMockRoleResult
}

type NetworkNodeMockRoleResult struct {
	r insolar.StaticRole
}

//Expect specifies that invocation of NetworkNode.Role is expected from 1 to Infinity times
func (m *mNetworkNodeMockRole) Expect() *mNetworkNodeMockRole {
	m.mock.RoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.Role
func (m *mNetworkNodeMockRole) Return(r insolar.StaticRole) *NetworkNodeMock {
	m.mock.RoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockRoleExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.Role is expected once
func (m *mNetworkNodeMockRole) ExpectOnce() *NetworkNodeMockRoleExpectation {
	m.mock.RoleFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockRoleExpectation) Return(r insolar.StaticRole) {
	e.result = &NetworkNodeMockRoleResult{r}
}

//Set uses given function f as a mock of NetworkNode.Role method
func (m *mNetworkNodeMockRole) Set(f func() (r insolar.StaticRole)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RoleFunc = f
	return m.mock
}

//Role implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) Role() (r insolar.StaticRole) {
	counter := atomic.AddUint64(&m.RolePreCounter, 1)
	defer atomic.AddUint64(&m.RoleCounter, 1)

	if len(m.RoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.Role.")
			return
		}

		result := m.RoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.Role")
			return
		}

		r = result.r

		return
	}

	if m.RoleMock.mainExpectation != nil {

		result := m.RoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.Role")
		}

		r = result.r

		return
	}

	if m.RoleFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.Role.")
		return
	}

	return m.RoleFunc()
}

//RoleMinimockCounter returns a count of NetworkNodeMock.RoleFunc invocations
func (m *NetworkNodeMock) RoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RoleCounter)
}

//RoleMinimockPreCounter returns the value of NetworkNodeMock.Role invocations
func (m *NetworkNodeMock) RoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RolePreCounter)
}

//RoleFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) RoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RoleCounter) == uint64(len(m.RoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RoleFunc != nil {
		return atomic.LoadUint64(&m.RoleCounter) > 0
	}

	return true
}

type mNetworkNodeMockShortID struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockShortIDExpectation
	expectationSeries []*NetworkNodeMockShortIDExpectation
}

type NetworkNodeMockShortIDExpectation struct {
	result *NetworkNodeMockShortIDResult
}

type NetworkNodeMockShortIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of NetworkNode.ShortID is expected from 1 to Infinity times
func (m *mNetworkNodeMockShortID) Expect() *mNetworkNodeMockShortID {
	m.mock.ShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockShortIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.ShortID
func (m *mNetworkNodeMockShortID) Return(r insolar.ShortNodeID) *NetworkNodeMock {
	m.mock.ShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockShortIDExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockShortIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.ShortID is expected once
func (m *mNetworkNodeMockShortID) ExpectOnce() *NetworkNodeMockShortIDExpectation {
	m.mock.ShortIDFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockShortIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockShortIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &NetworkNodeMockShortIDResult{r}
}

//Set uses given function f as a mock of NetworkNode.ShortID method
func (m *mNetworkNodeMockShortID) Set(f func() (r insolar.ShortNodeID)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ShortIDFunc = f
	return m.mock
}

//ShortID implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) ShortID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.ShortIDPreCounter, 1)
	defer atomic.AddUint64(&m.ShortIDCounter, 1)

	if len(m.ShortIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ShortIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.ShortID.")
			return
		}

		result := m.ShortIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.ShortID")
			return
		}

		r = result.r

		return
	}

	if m.ShortIDMock.mainExpectation != nil {

		result := m.ShortIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.ShortID")
		}

		r = result.r

		return
	}

	if m.ShortIDFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.ShortID.")
		return
	}

	return m.ShortIDFunc()
}

//ShortIDMinimockCounter returns a count of NetworkNodeMock.ShortIDFunc invocations
func (m *NetworkNodeMock) ShortIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ShortIDCounter)
}

//ShortIDMinimockPreCounter returns the value of NetworkNodeMock.ShortID invocations
func (m *NetworkNodeMock) ShortIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ShortIDPreCounter)
}

//ShortIDFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) ShortIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ShortIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ShortIDCounter) == uint64(len(m.ShortIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ShortIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ShortIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ShortIDFunc != nil {
		return atomic.LoadUint64(&m.ShortIDCounter) > 0
	}

	return true
}

type mNetworkNodeMockVersion struct {
	mock              *NetworkNodeMock
	mainExpectation   *NetworkNodeMockVersionExpectation
	expectationSeries []*NetworkNodeMockVersionExpectation
}

type NetworkNodeMockVersionExpectation struct {
	result *NetworkNodeMockVersionResult
}

type NetworkNodeMockVersionResult struct {
	r string
}

//Expect specifies that invocation of NetworkNode.Version is expected from 1 to Infinity times
func (m *mNetworkNodeMockVersion) Expect() *mNetworkNodeMockVersion {
	m.mock.VersionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockVersionExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkNode.Version
func (m *mNetworkNodeMockVersion) Return(r string) *NetworkNodeMock {
	m.mock.VersionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkNodeMockVersionExpectation{}
	}
	m.mainExpectation.result = &NetworkNodeMockVersionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkNode.Version is expected once
func (m *mNetworkNodeMockVersion) ExpectOnce() *NetworkNodeMockVersionExpectation {
	m.mock.VersionFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkNodeMockVersionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkNodeMockVersionExpectation) Return(r string) {
	e.result = &NetworkNodeMockVersionResult{r}
}

//Set uses given function f as a mock of NetworkNode.Version method
func (m *mNetworkNodeMockVersion) Set(f func() (r string)) *NetworkNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VersionFunc = f
	return m.mock
}

//Version implements github.com/insolar/insolar/insolar.NetworkNode interface
func (m *NetworkNodeMock) Version() (r string) {
	counter := atomic.AddUint64(&m.VersionPreCounter, 1)
	defer atomic.AddUint64(&m.VersionCounter, 1)

	if len(m.VersionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VersionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkNodeMock.Version.")
			return
		}

		result := m.VersionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.Version")
			return
		}

		r = result.r

		return
	}

	if m.VersionMock.mainExpectation != nil {

		result := m.VersionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkNodeMock.Version")
		}

		r = result.r

		return
	}

	if m.VersionFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkNodeMock.Version.")
		return
	}

	return m.VersionFunc()
}

//VersionMinimockCounter returns a count of NetworkNodeMock.VersionFunc invocations
func (m *NetworkNodeMock) VersionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VersionCounter)
}

//VersionMinimockPreCounter returns the value of NetworkNodeMock.Version invocations
func (m *NetworkNodeMock) VersionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VersionPreCounter)
}

//VersionFinished returns true if mock invocations count is ok
func (m *NetworkNodeMock) VersionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.VersionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.VersionCounter) == uint64(len(m.VersionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.VersionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.VersionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.VersionFunc != nil {
		return atomic.LoadUint64(&m.VersionCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkNodeMock) ValidateCallCounters() {

	if !m.AddressFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.Address")
	}

	if !m.GetGlobuleIDFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.GetGlobuleID")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.GetState")
	}

	if !m.IDFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.ID")
	}

	if !m.LeavingETAFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.LeavingETA")
	}

	if !m.PublicKeyFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.PublicKey")
	}

	if !m.RoleFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.Role")
	}

	if !m.ShortIDFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.ShortID")
	}

	if !m.VersionFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.Version")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NetworkNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NetworkNodeMock) MinimockFinish() {

	if !m.AddressFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.Address")
	}

	if !m.GetGlobuleIDFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.GetGlobuleID")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.GetState")
	}

	if !m.IDFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.ID")
	}

	if !m.LeavingETAFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.LeavingETA")
	}

	if !m.PublicKeyFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.PublicKey")
	}

	if !m.RoleFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.Role")
	}

	if !m.ShortIDFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.ShortID")
	}

	if !m.VersionFinished() {
		m.t.Fatal("Expected call to NetworkNodeMock.Version")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NetworkNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NetworkNodeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddressFinished()
		ok = ok && m.GetGlobuleIDFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.IDFinished()
		ok = ok && m.LeavingETAFinished()
		ok = ok && m.PublicKeyFinished()
		ok = ok && m.RoleFinished()
		ok = ok && m.ShortIDFinished()
		ok = ok && m.VersionFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddressFinished() {
				m.t.Error("Expected call to NetworkNodeMock.Address")
			}

			if !m.GetGlobuleIDFinished() {
				m.t.Error("Expected call to NetworkNodeMock.GetGlobuleID")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to NetworkNodeMock.GetState")
			}

			if !m.IDFinished() {
				m.t.Error("Expected call to NetworkNodeMock.ID")
			}

			if !m.LeavingETAFinished() {
				m.t.Error("Expected call to NetworkNodeMock.LeavingETA")
			}

			if !m.PublicKeyFinished() {
				m.t.Error("Expected call to NetworkNodeMock.PublicKey")
			}

			if !m.RoleFinished() {
				m.t.Error("Expected call to NetworkNodeMock.Role")
			}

			if !m.ShortIDFinished() {
				m.t.Error("Expected call to NetworkNodeMock.ShortID")
			}

			if !m.VersionFinished() {
				m.t.Error("Expected call to NetworkNodeMock.Version")
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
func (m *NetworkNodeMock) AllMocksCalled() bool {

	if !m.AddressFinished() {
		return false
	}

	if !m.GetGlobuleIDFinished() {
		return false
	}

	if !m.GetStateFinished() {
		return false
	}

	if !m.IDFinished() {
		return false
	}

	if !m.LeavingETAFinished() {
		return false
	}

	if !m.PublicKeyFinished() {
		return false
	}

	if !m.RoleFinished() {
		return false
	}

	if !m.ShortIDFinished() {
		return false
	}

	if !m.VersionFinished() {
		return false
	}

	return true
}
