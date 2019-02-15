package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Node" can be found in github.com/insolar/insolar/core
*/
import (
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
)

//NodeMock implements github.com/insolar/insolar/core.Node
type NodeMock struct {
	t minimock.Tester

	AddressFunc       func() (r string)
	AddressCounter    uint64
	AddressPreCounter uint64
	AddressMock       mNodeMockAddress

	ConsensusAddressFunc       func() (r string)
	ConsensusAddressCounter    uint64
	ConsensusAddressPreCounter uint64
	ConsensusAddressMock       mNodeMockConsensusAddress

	GetGlobuleIDFunc       func() (r core.GlobuleID)
	GetGlobuleIDCounter    uint64
	GetGlobuleIDPreCounter uint64
	GetGlobuleIDMock       mNodeMockGetGlobuleID

	IDFunc       func() (r core.RecordRef)
	IDCounter    uint64
	IDPreCounter uint64
	IDMock       mNodeMockID

	IsActiveFunc       func() (r bool)
	IsActiveCounter    uint64
	IsActivePreCounter uint64
	IsActiveMock       mNodeMockIsActive

	LeavingFunc       func() (r bool)
	LeavingCounter    uint64
	LeavingPreCounter uint64
	LeavingMock       mNodeMockLeaving

	LeavingETAFunc       func() (r core.PulseNumber)
	LeavingETACounter    uint64
	LeavingETAPreCounter uint64
	LeavingETAMock       mNodeMockLeavingETA

	PublicKeyFunc       func() (r crypto.PublicKey)
	PublicKeyCounter    uint64
	PublicKeyPreCounter uint64
	PublicKeyMock       mNodeMockPublicKey

	RoleFunc       func() (r core.StaticRole)
	RoleCounter    uint64
	RolePreCounter uint64
	RoleMock       mNodeMockRole

	ShortIDFunc       func() (r core.ShortNodeID)
	ShortIDCounter    uint64
	ShortIDPreCounter uint64
	ShortIDMock       mNodeMockShortID

	VersionFunc       func() (r string)
	VersionCounter    uint64
	VersionPreCounter uint64
	VersionMock       mNodeMockVersion
}

//NewNodeMock returns a mock for github.com/insolar/insolar/core.Node
func NewNodeMock(t minimock.Tester) *NodeMock {
	m := &NodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddressMock = mNodeMockAddress{mock: m}
	m.ConsensusAddressMock = mNodeMockConsensusAddress{mock: m}
	m.GetGlobuleIDMock = mNodeMockGetGlobuleID{mock: m}
	m.IDMock = mNodeMockID{mock: m}
	m.IsActiveMock = mNodeMockIsActive{mock: m}
	m.LeavingMock = mNodeMockLeaving{mock: m}
	m.LeavingETAMock = mNodeMockLeavingETA{mock: m}
	m.PublicKeyMock = mNodeMockPublicKey{mock: m}
	m.RoleMock = mNodeMockRole{mock: m}
	m.ShortIDMock = mNodeMockShortID{mock: m}
	m.VersionMock = mNodeMockVersion{mock: m}

	return m
}

type mNodeMockAddress struct {
	mock              *NodeMock
	mainExpectation   *NodeMockAddressExpectation
	expectationSeries []*NodeMockAddressExpectation
}

type NodeMockAddressExpectation struct {
	result *NodeMockAddressResult
}

type NodeMockAddressResult struct {
	r string
}

//Expect specifies that invocation of Node.Address is expected from 1 to Infinity times
func (m *mNodeMockAddress) Expect() *mNodeMockAddress {
	m.mock.AddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.Address
func (m *mNodeMockAddress) Return(r string) *NodeMock {
	m.mock.AddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockAddressExpectation{}
	}
	m.mainExpectation.result = &NodeMockAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.Address is expected once
func (m *mNodeMockAddress) ExpectOnce() *NodeMockAddressExpectation {
	m.mock.AddressFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockAddressExpectation) Return(r string) {
	e.result = &NodeMockAddressResult{r}
}

//Set uses given function f as a mock of Node.Address method
func (m *mNodeMockAddress) Set(f func() (r string)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddressFunc = f
	return m.mock
}

//Address implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) Address() (r string) {
	counter := atomic.AddUint64(&m.AddressPreCounter, 1)
	defer atomic.AddUint64(&m.AddressCounter, 1)

	if len(m.AddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.Address.")
			return
		}

		result := m.AddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Address")
			return
		}

		r = result.r

		return
	}

	if m.AddressMock.mainExpectation != nil {

		result := m.AddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Address")
		}

		r = result.r

		return
	}

	if m.AddressFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.Address.")
		return
	}

	return m.AddressFunc()
}

//AddressMinimockCounter returns a count of NodeMock.AddressFunc invocations
func (m *NodeMock) AddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddressCounter)
}

//AddressMinimockPreCounter returns the value of NodeMock.Address invocations
func (m *NodeMock) AddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddressPreCounter)
}

//AddressFinished returns true if mock invocations count is ok
func (m *NodeMock) AddressFinished() bool {
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

type mNodeMockConsensusAddress struct {
	mock              *NodeMock
	mainExpectation   *NodeMockConsensusAddressExpectation
	expectationSeries []*NodeMockConsensusAddressExpectation
}

type NodeMockConsensusAddressExpectation struct {
	result *NodeMockConsensusAddressResult
}

type NodeMockConsensusAddressResult struct {
	r string
}

//Expect specifies that invocation of Node.ConsensusAddress is expected from 1 to Infinity times
func (m *mNodeMockConsensusAddress) Expect() *mNodeMockConsensusAddress {
	m.mock.ConsensusAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockConsensusAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.ConsensusAddress
func (m *mNodeMockConsensusAddress) Return(r string) *NodeMock {
	m.mock.ConsensusAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockConsensusAddressExpectation{}
	}
	m.mainExpectation.result = &NodeMockConsensusAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.ConsensusAddress is expected once
func (m *mNodeMockConsensusAddress) ExpectOnce() *NodeMockConsensusAddressExpectation {
	m.mock.ConsensusAddressFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockConsensusAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockConsensusAddressExpectation) Return(r string) {
	e.result = &NodeMockConsensusAddressResult{r}
}

//Set uses given function f as a mock of Node.ConsensusAddress method
func (m *mNodeMockConsensusAddress) Set(f func() (r string)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ConsensusAddressFunc = f
	return m.mock
}

//ConsensusAddress implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) ConsensusAddress() (r string) {
	counter := atomic.AddUint64(&m.ConsensusAddressPreCounter, 1)
	defer atomic.AddUint64(&m.ConsensusAddressCounter, 1)

	if len(m.ConsensusAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ConsensusAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.ConsensusAddress.")
			return
		}

		result := m.ConsensusAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.ConsensusAddress")
			return
		}

		r = result.r

		return
	}

	if m.ConsensusAddressMock.mainExpectation != nil {

		result := m.ConsensusAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.ConsensusAddress")
		}

		r = result.r

		return
	}

	if m.ConsensusAddressFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.ConsensusAddress.")
		return
	}

	return m.ConsensusAddressFunc()
}

//ConsensusAddressMinimockCounter returns a count of NodeMock.ConsensusAddressFunc invocations
func (m *NodeMock) ConsensusAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ConsensusAddressCounter)
}

//ConsensusAddressMinimockPreCounter returns the value of NodeMock.ConsensusAddress invocations
func (m *NodeMock) ConsensusAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ConsensusAddressPreCounter)
}

//ConsensusAddressFinished returns true if mock invocations count is ok
func (m *NodeMock) ConsensusAddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ConsensusAddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ConsensusAddressCounter) == uint64(len(m.ConsensusAddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ConsensusAddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ConsensusAddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ConsensusAddressFunc != nil {
		return atomic.LoadUint64(&m.ConsensusAddressCounter) > 0
	}

	return true
}

type mNodeMockGetGlobuleID struct {
	mock              *NodeMock
	mainExpectation   *NodeMockGetGlobuleIDExpectation
	expectationSeries []*NodeMockGetGlobuleIDExpectation
}

type NodeMockGetGlobuleIDExpectation struct {
	result *NodeMockGetGlobuleIDResult
}

type NodeMockGetGlobuleIDResult struct {
	r core.GlobuleID
}

//Expect specifies that invocation of Node.GetGlobuleID is expected from 1 to Infinity times
func (m *mNodeMockGetGlobuleID) Expect() *mNodeMockGetGlobuleID {
	m.mock.GetGlobuleIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockGetGlobuleIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.GetGlobuleID
func (m *mNodeMockGetGlobuleID) Return(r core.GlobuleID) *NodeMock {
	m.mock.GetGlobuleIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockGetGlobuleIDExpectation{}
	}
	m.mainExpectation.result = &NodeMockGetGlobuleIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.GetGlobuleID is expected once
func (m *mNodeMockGetGlobuleID) ExpectOnce() *NodeMockGetGlobuleIDExpectation {
	m.mock.GetGlobuleIDFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockGetGlobuleIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockGetGlobuleIDExpectation) Return(r core.GlobuleID) {
	e.result = &NodeMockGetGlobuleIDResult{r}
}

//Set uses given function f as a mock of Node.GetGlobuleID method
func (m *mNodeMockGetGlobuleID) Set(f func() (r core.GlobuleID)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetGlobuleIDFunc = f
	return m.mock
}

//GetGlobuleID implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) GetGlobuleID() (r core.GlobuleID) {
	counter := atomic.AddUint64(&m.GetGlobuleIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobuleIDCounter, 1)

	if len(m.GetGlobuleIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetGlobuleIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.GetGlobuleID.")
			return
		}

		result := m.GetGlobuleIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.GetGlobuleID")
			return
		}

		r = result.r

		return
	}

	if m.GetGlobuleIDMock.mainExpectation != nil {

		result := m.GetGlobuleIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.GetGlobuleID")
		}

		r = result.r

		return
	}

	if m.GetGlobuleIDFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.GetGlobuleID.")
		return
	}

	return m.GetGlobuleIDFunc()
}

//GetGlobuleIDMinimockCounter returns a count of NodeMock.GetGlobuleIDFunc invocations
func (m *NodeMock) GetGlobuleIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleIDCounter)
}

//GetGlobuleIDMinimockPreCounter returns the value of NodeMock.GetGlobuleID invocations
func (m *NodeMock) GetGlobuleIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleIDPreCounter)
}

//GetGlobuleIDFinished returns true if mock invocations count is ok
func (m *NodeMock) GetGlobuleIDFinished() bool {
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

type mNodeMockID struct {
	mock              *NodeMock
	mainExpectation   *NodeMockIDExpectation
	expectationSeries []*NodeMockIDExpectation
}

type NodeMockIDExpectation struct {
	result *NodeMockIDResult
}

type NodeMockIDResult struct {
	r core.RecordRef
}

//Expect specifies that invocation of Node.ID is expected from 1 to Infinity times
func (m *mNodeMockID) Expect() *mNodeMockID {
	m.mock.IDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.ID
func (m *mNodeMockID) Return(r core.RecordRef) *NodeMock {
	m.mock.IDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockIDExpectation{}
	}
	m.mainExpectation.result = &NodeMockIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.ID is expected once
func (m *mNodeMockID) ExpectOnce() *NodeMockIDExpectation {
	m.mock.IDFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockIDExpectation) Return(r core.RecordRef) {
	e.result = &NodeMockIDResult{r}
}

//Set uses given function f as a mock of Node.ID method
func (m *mNodeMockID) Set(f func() (r core.RecordRef)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IDFunc = f
	return m.mock
}

//ID implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) ID() (r core.RecordRef) {
	counter := atomic.AddUint64(&m.IDPreCounter, 1)
	defer atomic.AddUint64(&m.IDCounter, 1)

	if len(m.IDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.ID.")
			return
		}

		result := m.IDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.ID")
			return
		}

		r = result.r

		return
	}

	if m.IDMock.mainExpectation != nil {

		result := m.IDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.ID")
		}

		r = result.r

		return
	}

	if m.IDFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.ID.")
		return
	}

	return m.IDFunc()
}

//IDMinimockCounter returns a count of NodeMock.IDFunc invocations
func (m *NodeMock) IDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IDCounter)
}

//IDMinimockPreCounter returns the value of NodeMock.ID invocations
func (m *NodeMock) IDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IDPreCounter)
}

//IDFinished returns true if mock invocations count is ok
func (m *NodeMock) IDFinished() bool {
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

type mNodeMockIsActive struct {
	mock              *NodeMock
	mainExpectation   *NodeMockIsActiveExpectation
	expectationSeries []*NodeMockIsActiveExpectation
}

type NodeMockIsActiveExpectation struct {
	result *NodeMockIsActiveResult
}

type NodeMockIsActiveResult struct {
	r bool
}

//Expect specifies that invocation of Node.IsActive is expected from 1 to Infinity times
func (m *mNodeMockIsActive) Expect() *mNodeMockIsActive {
	m.mock.IsActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockIsActiveExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.IsActive
func (m *mNodeMockIsActive) Return(r bool) *NodeMock {
	m.mock.IsActiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockIsActiveExpectation{}
	}
	m.mainExpectation.result = &NodeMockIsActiveResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.IsActive is expected once
func (m *mNodeMockIsActive) ExpectOnce() *NodeMockIsActiveExpectation {
	m.mock.IsActiveFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockIsActiveExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockIsActiveExpectation) Return(r bool) {
	e.result = &NodeMockIsActiveResult{r}
}

//Set uses given function f as a mock of Node.IsActive method
func (m *mNodeMockIsActive) Set(f func() (r bool)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsActiveFunc = f
	return m.mock
}

//IsActive implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) IsActive() (r bool) {
	counter := atomic.AddUint64(&m.IsActivePreCounter, 1)
	defer atomic.AddUint64(&m.IsActiveCounter, 1)

	if len(m.IsActiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsActiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.IsActive.")
			return
		}

		result := m.IsActiveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.IsActive")
			return
		}

		r = result.r

		return
	}

	if m.IsActiveMock.mainExpectation != nil {

		result := m.IsActiveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.IsActive")
		}

		r = result.r

		return
	}

	if m.IsActiveFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.IsActive.")
		return
	}

	return m.IsActiveFunc()
}

//IsActiveMinimockCounter returns a count of NodeMock.IsActiveFunc invocations
func (m *NodeMock) IsActiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsActiveCounter)
}

//IsActiveMinimockPreCounter returns the value of NodeMock.IsActive invocations
func (m *NodeMock) IsActiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsActivePreCounter)
}

//IsActiveFinished returns true if mock invocations count is ok
func (m *NodeMock) IsActiveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsActiveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsActiveCounter) == uint64(len(m.IsActiveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsActiveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsActiveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsActiveFunc != nil {
		return atomic.LoadUint64(&m.IsActiveCounter) > 0
	}

	return true
}

type mNodeMockLeaving struct {
	mock              *NodeMock
	mainExpectation   *NodeMockLeavingExpectation
	expectationSeries []*NodeMockLeavingExpectation
}

type NodeMockLeavingExpectation struct {
	result *NodeMockLeavingResult
}

type NodeMockLeavingResult struct {
	r bool
}

//Expect specifies that invocation of Node.Leaving is expected from 1 to Infinity times
func (m *mNodeMockLeaving) Expect() *mNodeMockLeaving {
	m.mock.LeavingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockLeavingExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.Leaving
func (m *mNodeMockLeaving) Return(r bool) *NodeMock {
	m.mock.LeavingFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockLeavingExpectation{}
	}
	m.mainExpectation.result = &NodeMockLeavingResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.Leaving is expected once
func (m *mNodeMockLeaving) ExpectOnce() *NodeMockLeavingExpectation {
	m.mock.LeavingFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockLeavingExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockLeavingExpectation) Return(r bool) {
	e.result = &NodeMockLeavingResult{r}
}

//Set uses given function f as a mock of Node.Leaving method
func (m *mNodeMockLeaving) Set(f func() (r bool)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LeavingFunc = f
	return m.mock
}

//Leaving implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) Leaving() (r bool) {
	counter := atomic.AddUint64(&m.LeavingPreCounter, 1)
	defer atomic.AddUint64(&m.LeavingCounter, 1)

	if len(m.LeavingMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LeavingMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.Leaving.")
			return
		}

		result := m.LeavingMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Leaving")
			return
		}

		r = result.r

		return
	}

	if m.LeavingMock.mainExpectation != nil {

		result := m.LeavingMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Leaving")
		}

		r = result.r

		return
	}

	if m.LeavingFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.Leaving.")
		return
	}

	return m.LeavingFunc()
}

//LeavingMinimockCounter returns a count of NodeMock.LeavingFunc invocations
func (m *NodeMock) LeavingMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LeavingCounter)
}

//LeavingMinimockPreCounter returns the value of NodeMock.Leaving invocations
func (m *NodeMock) LeavingMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LeavingPreCounter)
}

//LeavingFinished returns true if mock invocations count is ok
func (m *NodeMock) LeavingFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LeavingMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LeavingCounter) == uint64(len(m.LeavingMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LeavingMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LeavingCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LeavingFunc != nil {
		return atomic.LoadUint64(&m.LeavingCounter) > 0
	}

	return true
}

type mNodeMockLeavingETA struct {
	mock              *NodeMock
	mainExpectation   *NodeMockLeavingETAExpectation
	expectationSeries []*NodeMockLeavingETAExpectation
}

type NodeMockLeavingETAExpectation struct {
	result *NodeMockLeavingETAResult
}

type NodeMockLeavingETAResult struct {
	r core.PulseNumber
}

//Expect specifies that invocation of Node.LeavingETA is expected from 1 to Infinity times
func (m *mNodeMockLeavingETA) Expect() *mNodeMockLeavingETA {
	m.mock.LeavingETAFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockLeavingETAExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.LeavingETA
func (m *mNodeMockLeavingETA) Return(r core.PulseNumber) *NodeMock {
	m.mock.LeavingETAFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockLeavingETAExpectation{}
	}
	m.mainExpectation.result = &NodeMockLeavingETAResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.LeavingETA is expected once
func (m *mNodeMockLeavingETA) ExpectOnce() *NodeMockLeavingETAExpectation {
	m.mock.LeavingETAFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockLeavingETAExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockLeavingETAExpectation) Return(r core.PulseNumber) {
	e.result = &NodeMockLeavingETAResult{r}
}

//Set uses given function f as a mock of Node.LeavingETA method
func (m *mNodeMockLeavingETA) Set(f func() (r core.PulseNumber)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LeavingETAFunc = f
	return m.mock
}

//LeavingETA implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) LeavingETA() (r core.PulseNumber) {
	counter := atomic.AddUint64(&m.LeavingETAPreCounter, 1)
	defer atomic.AddUint64(&m.LeavingETACounter, 1)

	if len(m.LeavingETAMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LeavingETAMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.LeavingETA.")
			return
		}

		result := m.LeavingETAMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.LeavingETA")
			return
		}

		r = result.r

		return
	}

	if m.LeavingETAMock.mainExpectation != nil {

		result := m.LeavingETAMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.LeavingETA")
		}

		r = result.r

		return
	}

	if m.LeavingETAFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.LeavingETA.")
		return
	}

	return m.LeavingETAFunc()
}

//LeavingETAMinimockCounter returns a count of NodeMock.LeavingETAFunc invocations
func (m *NodeMock) LeavingETAMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LeavingETACounter)
}

//LeavingETAMinimockPreCounter returns the value of NodeMock.LeavingETA invocations
func (m *NodeMock) LeavingETAMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LeavingETAPreCounter)
}

//LeavingETAFinished returns true if mock invocations count is ok
func (m *NodeMock) LeavingETAFinished() bool {
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

type mNodeMockPublicKey struct {
	mock              *NodeMock
	mainExpectation   *NodeMockPublicKeyExpectation
	expectationSeries []*NodeMockPublicKeyExpectation
}

type NodeMockPublicKeyExpectation struct {
	result *NodeMockPublicKeyResult
}

type NodeMockPublicKeyResult struct {
	r crypto.PublicKey
}

//Expect specifies that invocation of Node.PublicKey is expected from 1 to Infinity times
func (m *mNodeMockPublicKey) Expect() *mNodeMockPublicKey {
	m.mock.PublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockPublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.PublicKey
func (m *mNodeMockPublicKey) Return(r crypto.PublicKey) *NodeMock {
	m.mock.PublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockPublicKeyExpectation{}
	}
	m.mainExpectation.result = &NodeMockPublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.PublicKey is expected once
func (m *mNodeMockPublicKey) ExpectOnce() *NodeMockPublicKeyExpectation {
	m.mock.PublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockPublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockPublicKeyExpectation) Return(r crypto.PublicKey) {
	e.result = &NodeMockPublicKeyResult{r}
}

//Set uses given function f as a mock of Node.PublicKey method
func (m *mNodeMockPublicKey) Set(f func() (r crypto.PublicKey)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PublicKeyFunc = f
	return m.mock
}

//PublicKey implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) PublicKey() (r crypto.PublicKey) {
	counter := atomic.AddUint64(&m.PublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.PublicKeyCounter, 1)

	if len(m.PublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.PublicKey.")
			return
		}

		result := m.PublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.PublicKey")
			return
		}

		r = result.r

		return
	}

	if m.PublicKeyMock.mainExpectation != nil {

		result := m.PublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.PublicKey")
		}

		r = result.r

		return
	}

	if m.PublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.PublicKey.")
		return
	}

	return m.PublicKeyFunc()
}

//PublicKeyMinimockCounter returns a count of NodeMock.PublicKeyFunc invocations
func (m *NodeMock) PublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PublicKeyCounter)
}

//PublicKeyMinimockPreCounter returns the value of NodeMock.PublicKey invocations
func (m *NodeMock) PublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PublicKeyPreCounter)
}

//PublicKeyFinished returns true if mock invocations count is ok
func (m *NodeMock) PublicKeyFinished() bool {
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

type mNodeMockRole struct {
	mock              *NodeMock
	mainExpectation   *NodeMockRoleExpectation
	expectationSeries []*NodeMockRoleExpectation
}

type NodeMockRoleExpectation struct {
	result *NodeMockRoleResult
}

type NodeMockRoleResult struct {
	r core.StaticRole
}

//Expect specifies that invocation of Node.Role is expected from 1 to Infinity times
func (m *mNodeMockRole) Expect() *mNodeMockRole {
	m.mock.RoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.Role
func (m *mNodeMockRole) Return(r core.StaticRole) *NodeMock {
	m.mock.RoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockRoleExpectation{}
	}
	m.mainExpectation.result = &NodeMockRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.Role is expected once
func (m *mNodeMockRole) ExpectOnce() *NodeMockRoleExpectation {
	m.mock.RoleFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockRoleExpectation) Return(r core.StaticRole) {
	e.result = &NodeMockRoleResult{r}
}

//Set uses given function f as a mock of Node.Role method
func (m *mNodeMockRole) Set(f func() (r core.StaticRole)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RoleFunc = f
	return m.mock
}

//Role implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) Role() (r core.StaticRole) {
	counter := atomic.AddUint64(&m.RolePreCounter, 1)
	defer atomic.AddUint64(&m.RoleCounter, 1)

	if len(m.RoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.Role.")
			return
		}

		result := m.RoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Role")
			return
		}

		r = result.r

		return
	}

	if m.RoleMock.mainExpectation != nil {

		result := m.RoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Role")
		}

		r = result.r

		return
	}

	if m.RoleFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.Role.")
		return
	}

	return m.RoleFunc()
}

//RoleMinimockCounter returns a count of NodeMock.RoleFunc invocations
func (m *NodeMock) RoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RoleCounter)
}

//RoleMinimockPreCounter returns the value of NodeMock.Role invocations
func (m *NodeMock) RoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RolePreCounter)
}

//RoleFinished returns true if mock invocations count is ok
func (m *NodeMock) RoleFinished() bool {
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

type mNodeMockShortID struct {
	mock              *NodeMock
	mainExpectation   *NodeMockShortIDExpectation
	expectationSeries []*NodeMockShortIDExpectation
}

type NodeMockShortIDExpectation struct {
	result *NodeMockShortIDResult
}

type NodeMockShortIDResult struct {
	r core.ShortNodeID
}

//Expect specifies that invocation of Node.ShortID is expected from 1 to Infinity times
func (m *mNodeMockShortID) Expect() *mNodeMockShortID {
	m.mock.ShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockShortIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.ShortID
func (m *mNodeMockShortID) Return(r core.ShortNodeID) *NodeMock {
	m.mock.ShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockShortIDExpectation{}
	}
	m.mainExpectation.result = &NodeMockShortIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.ShortID is expected once
func (m *mNodeMockShortID) ExpectOnce() *NodeMockShortIDExpectation {
	m.mock.ShortIDFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockShortIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockShortIDExpectation) Return(r core.ShortNodeID) {
	e.result = &NodeMockShortIDResult{r}
}

//Set uses given function f as a mock of Node.ShortID method
func (m *mNodeMockShortID) Set(f func() (r core.ShortNodeID)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ShortIDFunc = f
	return m.mock
}

//ShortID implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) ShortID() (r core.ShortNodeID) {
	counter := atomic.AddUint64(&m.ShortIDPreCounter, 1)
	defer atomic.AddUint64(&m.ShortIDCounter, 1)

	if len(m.ShortIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ShortIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.ShortID.")
			return
		}

		result := m.ShortIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.ShortID")
			return
		}

		r = result.r

		return
	}

	if m.ShortIDMock.mainExpectation != nil {

		result := m.ShortIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.ShortID")
		}

		r = result.r

		return
	}

	if m.ShortIDFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.ShortID.")
		return
	}

	return m.ShortIDFunc()
}

//ShortIDMinimockCounter returns a count of NodeMock.ShortIDFunc invocations
func (m *NodeMock) ShortIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ShortIDCounter)
}

//ShortIDMinimockPreCounter returns the value of NodeMock.ShortID invocations
func (m *NodeMock) ShortIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ShortIDPreCounter)
}

//ShortIDFinished returns true if mock invocations count is ok
func (m *NodeMock) ShortIDFinished() bool {
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

type mNodeMockVersion struct {
	mock              *NodeMock
	mainExpectation   *NodeMockVersionExpectation
	expectationSeries []*NodeMockVersionExpectation
}

type NodeMockVersionExpectation struct {
	result *NodeMockVersionResult
}

type NodeMockVersionResult struct {
	r string
}

//Expect specifies that invocation of Node.Version is expected from 1 to Infinity times
func (m *mNodeMockVersion) Expect() *mNodeMockVersion {
	m.mock.VersionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockVersionExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.Version
func (m *mNodeMockVersion) Return(r string) *NodeMock {
	m.mock.VersionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockVersionExpectation{}
	}
	m.mainExpectation.result = &NodeMockVersionResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.Version is expected once
func (m *mNodeMockVersion) ExpectOnce() *NodeMockVersionExpectation {
	m.mock.VersionFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockVersionExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockVersionExpectation) Return(r string) {
	e.result = &NodeMockVersionResult{r}
}

//Set uses given function f as a mock of Node.Version method
func (m *mNodeMockVersion) Set(f func() (r string)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VersionFunc = f
	return m.mock
}

//Version implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) Version() (r string) {
	counter := atomic.AddUint64(&m.VersionPreCounter, 1)
	defer atomic.AddUint64(&m.VersionCounter, 1)

	if len(m.VersionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VersionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.Version.")
			return
		}

		result := m.VersionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Version")
			return
		}

		r = result.r

		return
	}

	if m.VersionMock.mainExpectation != nil {

		result := m.VersionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Version")
		}

		r = result.r

		return
	}

	if m.VersionFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.Version.")
		return
	}

	return m.VersionFunc()
}

//VersionMinimockCounter returns a count of NodeMock.VersionFunc invocations
func (m *NodeMock) VersionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VersionCounter)
}

//VersionMinimockPreCounter returns the value of NodeMock.Version invocations
func (m *NodeMock) VersionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VersionPreCounter)
}

//VersionFinished returns true if mock invocations count is ok
func (m *NodeMock) VersionFinished() bool {
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
func (m *NodeMock) ValidateCallCounters() {

	if !m.AddressFinished() {
		m.t.Fatal("Expected call to NodeMock.Address")
	}

	if !m.ConsensusAddressFinished() {
		m.t.Fatal("Expected call to NodeMock.ConsensusAddress")
	}

	if !m.GetGlobuleIDFinished() {
		m.t.Fatal("Expected call to NodeMock.GetGlobuleID")
	}

	if !m.IDFinished() {
		m.t.Fatal("Expected call to NodeMock.ID")
	}

	if !m.IsActiveFinished() {
		m.t.Fatal("Expected call to NodeMock.IsActive")
	}

	if !m.LeavingFinished() {
		m.t.Fatal("Expected call to NodeMock.Leaving")
	}

	if !m.LeavingETAFinished() {
		m.t.Fatal("Expected call to NodeMock.LeavingETA")
	}

	if !m.PublicKeyFinished() {
		m.t.Fatal("Expected call to NodeMock.PublicKey")
	}

	if !m.RoleFinished() {
		m.t.Fatal("Expected call to NodeMock.Role")
	}

	if !m.ShortIDFinished() {
		m.t.Fatal("Expected call to NodeMock.ShortID")
	}

	if !m.VersionFinished() {
		m.t.Fatal("Expected call to NodeMock.Version")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeMock) MinimockFinish() {

	if !m.AddressFinished() {
		m.t.Fatal("Expected call to NodeMock.Address")
	}

	if !m.ConsensusAddressFinished() {
		m.t.Fatal("Expected call to NodeMock.ConsensusAddress")
	}

	if !m.GetGlobuleIDFinished() {
		m.t.Fatal("Expected call to NodeMock.GetGlobuleID")
	}

	if !m.IDFinished() {
		m.t.Fatal("Expected call to NodeMock.ID")
	}

	if !m.IsActiveFinished() {
		m.t.Fatal("Expected call to NodeMock.IsActive")
	}

	if !m.LeavingFinished() {
		m.t.Fatal("Expected call to NodeMock.Leaving")
	}

	if !m.LeavingETAFinished() {
		m.t.Fatal("Expected call to NodeMock.LeavingETA")
	}

	if !m.PublicKeyFinished() {
		m.t.Fatal("Expected call to NodeMock.PublicKey")
	}

	if !m.RoleFinished() {
		m.t.Fatal("Expected call to NodeMock.Role")
	}

	if !m.ShortIDFinished() {
		m.t.Fatal("Expected call to NodeMock.ShortID")
	}

	if !m.VersionFinished() {
		m.t.Fatal("Expected call to NodeMock.Version")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddressFinished()
		ok = ok && m.ConsensusAddressFinished()
		ok = ok && m.GetGlobuleIDFinished()
		ok = ok && m.IDFinished()
		ok = ok && m.IsActiveFinished()
		ok = ok && m.LeavingFinished()
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
				m.t.Error("Expected call to NodeMock.Address")
			}

			if !m.ConsensusAddressFinished() {
				m.t.Error("Expected call to NodeMock.ConsensusAddress")
			}

			if !m.GetGlobuleIDFinished() {
				m.t.Error("Expected call to NodeMock.GetGlobuleID")
			}

			if !m.IDFinished() {
				m.t.Error("Expected call to NodeMock.ID")
			}

			if !m.IsActiveFinished() {
				m.t.Error("Expected call to NodeMock.IsActive")
			}

			if !m.LeavingFinished() {
				m.t.Error("Expected call to NodeMock.Leaving")
			}

			if !m.LeavingETAFinished() {
				m.t.Error("Expected call to NodeMock.LeavingETA")
			}

			if !m.PublicKeyFinished() {
				m.t.Error("Expected call to NodeMock.PublicKey")
			}

			if !m.RoleFinished() {
				m.t.Error("Expected call to NodeMock.Role")
			}

			if !m.ShortIDFinished() {
				m.t.Error("Expected call to NodeMock.ShortID")
			}

			if !m.VersionFinished() {
				m.t.Error("Expected call to NodeMock.Version")
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
func (m *NodeMock) AllMocksCalled() bool {

	if !m.AddressFinished() {
		return false
	}

	if !m.ConsensusAddressFinished() {
		return false
	}

	if !m.GetGlobuleIDFinished() {
		return false
	}

	if !m.IDFinished() {
		return false
	}

	if !m.IsActiveFinished() {
		return false
	}

	if !m.LeavingFinished() {
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
