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

	IDFunc       func() (r core.RecordRef)
	IDCounter    uint64
	IDPreCounter uint64
	IDMock       mNodeMockID

	PhysicalAddressFunc       func() (r string)
	PhysicalAddressCounter    uint64
	PhysicalAddressPreCounter uint64
	PhysicalAddressMock       mNodeMockPhysicalAddress

	PublicKeyFunc       func() (r crypto.PublicKey)
	PublicKeyCounter    uint64
	PublicKeyPreCounter uint64
	PublicKeyMock       mNodeMockPublicKey

	PulseFunc       func() (r core.PulseNumber)
	PulseCounter    uint64
	PulsePreCounter uint64
	PulseMock       mNodeMockPulse

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

	m.IDMock = mNodeMockID{mock: m}
	m.PhysicalAddressMock = mNodeMockPhysicalAddress{mock: m}
	m.PublicKeyMock = mNodeMockPublicKey{mock: m}
	m.PulseMock = mNodeMockPulse{mock: m}
	m.RoleMock = mNodeMockRole{mock: m}
	m.ShortIDMock = mNodeMockShortID{mock: m}
	m.VersionMock = mNodeMockVersion{mock: m}

	return m
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

type mNodeMockPhysicalAddress struct {
	mock              *NodeMock
	mainExpectation   *NodeMockPhysicalAddressExpectation
	expectationSeries []*NodeMockPhysicalAddressExpectation
}

type NodeMockPhysicalAddressExpectation struct {
	result *NodeMockPhysicalAddressResult
}

type NodeMockPhysicalAddressResult struct {
	r string
}

//Expect specifies that invocation of Node.PhysicalAddress is expected from 1 to Infinity times
func (m *mNodeMockPhysicalAddress) Expect() *mNodeMockPhysicalAddress {
	m.mock.PhysicalAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockPhysicalAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.PhysicalAddress
func (m *mNodeMockPhysicalAddress) Return(r string) *NodeMock {
	m.mock.PhysicalAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockPhysicalAddressExpectation{}
	}
	m.mainExpectation.result = &NodeMockPhysicalAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.PhysicalAddress is expected once
func (m *mNodeMockPhysicalAddress) ExpectOnce() *NodeMockPhysicalAddressExpectation {
	m.mock.PhysicalAddressFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockPhysicalAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockPhysicalAddressExpectation) Return(r string) {
	e.result = &NodeMockPhysicalAddressResult{r}
}

//Set uses given function f as a mock of Node.PhysicalAddress method
func (m *mNodeMockPhysicalAddress) Set(f func() (r string)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PhysicalAddressFunc = f
	return m.mock
}

//PhysicalAddress implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) PhysicalAddress() (r string) {
	counter := atomic.AddUint64(&m.PhysicalAddressPreCounter, 1)
	defer atomic.AddUint64(&m.PhysicalAddressCounter, 1)

	if len(m.PhysicalAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PhysicalAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.PhysicalAddress.")
			return
		}

		result := m.PhysicalAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.PhysicalAddress")
			return
		}

		r = result.r

		return
	}

	if m.PhysicalAddressMock.mainExpectation != nil {

		result := m.PhysicalAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.PhysicalAddress")
		}

		r = result.r

		return
	}

	if m.PhysicalAddressFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.PhysicalAddress.")
		return
	}

	return m.PhysicalAddressFunc()
}

//PhysicalAddressMinimockCounter returns a count of NodeMock.PhysicalAddressFunc invocations
func (m *NodeMock) PhysicalAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PhysicalAddressCounter)
}

//PhysicalAddressMinimockPreCounter returns the value of NodeMock.PhysicalAddress invocations
func (m *NodeMock) PhysicalAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PhysicalAddressPreCounter)
}

//PhysicalAddressFinished returns true if mock invocations count is ok
func (m *NodeMock) PhysicalAddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PhysicalAddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PhysicalAddressCounter) == uint64(len(m.PhysicalAddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PhysicalAddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PhysicalAddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PhysicalAddressFunc != nil {
		return atomic.LoadUint64(&m.PhysicalAddressCounter) > 0
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

type mNodeMockPulse struct {
	mock              *NodeMock
	mainExpectation   *NodeMockPulseExpectation
	expectationSeries []*NodeMockPulseExpectation
}

type NodeMockPulseExpectation struct {
	result *NodeMockPulseResult
}

type NodeMockPulseResult struct {
	r core.PulseNumber
}

//Expect specifies that invocation of Node.Pulse is expected from 1 to Infinity times
func (m *mNodeMockPulse) Expect() *mNodeMockPulse {
	m.mock.PulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of Node.Pulse
func (m *mNodeMockPulse) Return(r core.PulseNumber) *NodeMock {
	m.mock.PulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeMockPulseExpectation{}
	}
	m.mainExpectation.result = &NodeMockPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Node.Pulse is expected once
func (m *mNodeMockPulse) ExpectOnce() *NodeMockPulseExpectation {
	m.mock.PulseFunc = nil
	m.mainExpectation = nil

	expectation := &NodeMockPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeMockPulseExpectation) Return(r core.PulseNumber) {
	e.result = &NodeMockPulseResult{r}
}

//Set uses given function f as a mock of Node.Pulse method
func (m *mNodeMockPulse) Set(f func() (r core.PulseNumber)) *NodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PulseFunc = f
	return m.mock
}

//Pulse implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) Pulse() (r core.PulseNumber) {
	counter := atomic.AddUint64(&m.PulsePreCounter, 1)
	defer atomic.AddUint64(&m.PulseCounter, 1)

	if len(m.PulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeMock.Pulse.")
			return
		}

		result := m.PulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Pulse")
			return
		}

		r = result.r

		return
	}

	if m.PulseMock.mainExpectation != nil {

		result := m.PulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeMock.Pulse")
		}

		r = result.r

		return
	}

	if m.PulseFunc == nil {
		m.t.Fatalf("Unexpected call to NodeMock.Pulse.")
		return
	}

	return m.PulseFunc()
}

//PulseMinimockCounter returns a count of NodeMock.PulseFunc invocations
func (m *NodeMock) PulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PulseCounter)
}

//PulseMinimockPreCounter returns the value of NodeMock.Pulse invocations
func (m *NodeMock) PulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PulsePreCounter)
}

//PulseFinished returns true if mock invocations count is ok
func (m *NodeMock) PulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PulseCounter) == uint64(len(m.PulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PulseFunc != nil {
		return atomic.LoadUint64(&m.PulseCounter) > 0
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

	if !m.IDFinished() {
		m.t.Fatal("Expected call to NodeMock.ID")
	}

	if !m.PhysicalAddressFinished() {
		m.t.Fatal("Expected call to NodeMock.PhysicalAddress")
	}

	if !m.PublicKeyFinished() {
		m.t.Fatal("Expected call to NodeMock.PublicKey")
	}

	if !m.PulseFinished() {
		m.t.Fatal("Expected call to NodeMock.Pulse")
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

	if !m.IDFinished() {
		m.t.Fatal("Expected call to NodeMock.ID")
	}

	if !m.PhysicalAddressFinished() {
		m.t.Fatal("Expected call to NodeMock.PhysicalAddress")
	}

	if !m.PublicKeyFinished() {
		m.t.Fatal("Expected call to NodeMock.PublicKey")
	}

	if !m.PulseFinished() {
		m.t.Fatal("Expected call to NodeMock.Pulse")
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
		ok = ok && m.IDFinished()
		ok = ok && m.PhysicalAddressFinished()
		ok = ok && m.PublicKeyFinished()
		ok = ok && m.PulseFinished()
		ok = ok && m.RoleFinished()
		ok = ok && m.ShortIDFinished()
		ok = ok && m.VersionFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IDFinished() {
				m.t.Error("Expected call to NodeMock.ID")
			}

			if !m.PhysicalAddressFinished() {
				m.t.Error("Expected call to NodeMock.PhysicalAddress")
			}

			if !m.PublicKeyFinished() {
				m.t.Error("Expected call to NodeMock.PublicKey")
			}

			if !m.PulseFinished() {
				m.t.Error("Expected call to NodeMock.Pulse")
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

	if !m.IDFinished() {
		return false
	}

	if !m.PhysicalAddressFinished() {
		return false
	}

	if !m.PublicKeyFinished() {
		return false
	}

	if !m.PulseFinished() {
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
