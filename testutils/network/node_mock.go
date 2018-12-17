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

	GetGlobuleIDFunc       func() (r core.GlobuleID)
	GetGlobuleIDCounter    uint64
	GetGlobuleIDPreCounter uint64
	GetGlobuleIDMock       mNodeMockGetGlobuleID

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

	m.GetGlobuleIDMock = mNodeMockGetGlobuleID{mock: m}
	m.IDMock = mNodeMockID{mock: m}
	m.PhysicalAddressMock = mNodeMockPhysicalAddress{mock: m}
	m.PublicKeyMock = mNodeMockPublicKey{mock: m}
	m.RoleMock = mNodeMockRole{mock: m}
	m.ShortIDMock = mNodeMockShortID{mock: m}
	m.VersionMock = mNodeMockVersion{mock: m}

	return m
}

type mNodeMockGetGlobuleID struct {
	mock *NodeMock
}

//Return sets up a mock for Node.GetGlobuleID to return Return's arguments
func (m *mNodeMockGetGlobuleID) Return(r core.GlobuleID) *NodeMock {
	m.mock.GetGlobuleIDFunc = func() core.GlobuleID {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Node.GetGlobuleID method
func (m *mNodeMockGetGlobuleID) Set(f func() (r core.GlobuleID)) *NodeMock {
	m.mock.GetGlobuleIDFunc = f

	return m.mock
}

//GetGlobuleID implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) GetGlobuleID() (r core.GlobuleID) {
	atomic.AddUint64(&m.GetGlobuleIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobuleIDCounter, 1)

	if m.GetGlobuleIDFunc == nil {
		m.t.Fatal("Unexpected call to NodeMock.GetGlobuleID")
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

type mNodeMockID struct {
	mock *NodeMock
}

//Return sets up a mock for Node.ID to return Return's arguments
func (m *mNodeMockID) Return(r core.RecordRef) *NodeMock {
	m.mock.IDFunc = func() core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Node.ID method
func (m *mNodeMockID) Set(f func() (r core.RecordRef)) *NodeMock {
	m.mock.IDFunc = f

	return m.mock
}

//ID implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) ID() (r core.RecordRef) {
	atomic.AddUint64(&m.IDPreCounter, 1)
	defer atomic.AddUint64(&m.IDCounter, 1)

	if m.IDFunc == nil {
		m.t.Fatal("Unexpected call to NodeMock.ID")
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

type mNodeMockPhysicalAddress struct {
	mock *NodeMock
}

//Return sets up a mock for Node.PhysicalAddress to return Return's arguments
func (m *mNodeMockPhysicalAddress) Return(r string) *NodeMock {
	m.mock.PhysicalAddressFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Node.PhysicalAddress method
func (m *mNodeMockPhysicalAddress) Set(f func() (r string)) *NodeMock {
	m.mock.PhysicalAddressFunc = f

	return m.mock
}

//PhysicalAddress implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) PhysicalAddress() (r string) {
	atomic.AddUint64(&m.PhysicalAddressPreCounter, 1)
	defer atomic.AddUint64(&m.PhysicalAddressCounter, 1)

	if m.PhysicalAddressFunc == nil {
		m.t.Fatal("Unexpected call to NodeMock.PhysicalAddress")
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

type mNodeMockPublicKey struct {
	mock *NodeMock
}

//Return sets up a mock for Node.PublicKey to return Return's arguments
func (m *mNodeMockPublicKey) Return(r crypto.PublicKey) *NodeMock {
	m.mock.PublicKeyFunc = func() crypto.PublicKey {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Node.PublicKey method
func (m *mNodeMockPublicKey) Set(f func() (r crypto.PublicKey)) *NodeMock {
	m.mock.PublicKeyFunc = f

	return m.mock
}

//PublicKey implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) PublicKey() (r crypto.PublicKey) {
	atomic.AddUint64(&m.PublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.PublicKeyCounter, 1)

	if m.PublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to NodeMock.PublicKey")
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

type mNodeMockRole struct {
	mock *NodeMock
}

//Return sets up a mock for Node.Role to return Return's arguments
func (m *mNodeMockRole) Return(r core.StaticRole) *NodeMock {
	m.mock.RoleFunc = func() core.StaticRole {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Node.Role method
func (m *mNodeMockRole) Set(f func() (r core.StaticRole)) *NodeMock {
	m.mock.RoleFunc = f

	return m.mock
}

//Role implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) Role() (r core.StaticRole) {
	atomic.AddUint64(&m.RolePreCounter, 1)
	defer atomic.AddUint64(&m.RoleCounter, 1)

	if m.RoleFunc == nil {
		m.t.Fatal("Unexpected call to NodeMock.Role")
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

type mNodeMockShortID struct {
	mock *NodeMock
}

//Return sets up a mock for Node.ShortID to return Return's arguments
func (m *mNodeMockShortID) Return(r core.ShortNodeID) *NodeMock {
	m.mock.ShortIDFunc = func() core.ShortNodeID {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Node.ShortID method
func (m *mNodeMockShortID) Set(f func() (r core.ShortNodeID)) *NodeMock {
	m.mock.ShortIDFunc = f

	return m.mock
}

//ShortID implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) ShortID() (r core.ShortNodeID) {
	atomic.AddUint64(&m.ShortIDPreCounter, 1)
	defer atomic.AddUint64(&m.ShortIDCounter, 1)

	if m.ShortIDFunc == nil {
		m.t.Fatal("Unexpected call to NodeMock.ShortID")
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

type mNodeMockVersion struct {
	mock *NodeMock
}

//Return sets up a mock for Node.Version to return Return's arguments
func (m *mNodeMockVersion) Return(r string) *NodeMock {
	m.mock.VersionFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Node.Version method
func (m *mNodeMockVersion) Set(f func() (r string)) *NodeMock {
	m.mock.VersionFunc = f

	return m.mock
}

//Version implements github.com/insolar/insolar/core.Node interface
func (m *NodeMock) Version() (r string) {
	atomic.AddUint64(&m.VersionPreCounter, 1)
	defer atomic.AddUint64(&m.VersionCounter, 1)

	if m.VersionFunc == nil {
		m.t.Fatal("Unexpected call to NodeMock.Version")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeMock) ValidateCallCounters() {

	if m.GetGlobuleIDFunc != nil && atomic.LoadUint64(&m.GetGlobuleIDCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.GetGlobuleID")
	}

	if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.ID")
	}

	if m.PhysicalAddressFunc != nil && atomic.LoadUint64(&m.PhysicalAddressCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.PhysicalAddress")
	}

	if m.PublicKeyFunc != nil && atomic.LoadUint64(&m.PublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.PublicKey")
	}

	if m.RoleFunc != nil && atomic.LoadUint64(&m.RoleCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.Role")
	}

	if m.ShortIDFunc != nil && atomic.LoadUint64(&m.ShortIDCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.ShortID")
	}

	if m.VersionFunc != nil && atomic.LoadUint64(&m.VersionCounter) == 0 {
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

	if m.GetGlobuleIDFunc != nil && atomic.LoadUint64(&m.GetGlobuleIDCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.GetGlobuleID")
	}

	if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.ID")
	}

	if m.PhysicalAddressFunc != nil && atomic.LoadUint64(&m.PhysicalAddressCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.PhysicalAddress")
	}

	if m.PublicKeyFunc != nil && atomic.LoadUint64(&m.PublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.PublicKey")
	}

	if m.RoleFunc != nil && atomic.LoadUint64(&m.RoleCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.Role")
	}

	if m.ShortIDFunc != nil && atomic.LoadUint64(&m.ShortIDCounter) == 0 {
		m.t.Fatal("Expected call to NodeMock.ShortID")
	}

	if m.VersionFunc != nil && atomic.LoadUint64(&m.VersionCounter) == 0 {
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
		ok = ok && (m.GetGlobuleIDFunc == nil || atomic.LoadUint64(&m.GetGlobuleIDCounter) > 0)
		ok = ok && (m.IDFunc == nil || atomic.LoadUint64(&m.IDCounter) > 0)
		ok = ok && (m.PhysicalAddressFunc == nil || atomic.LoadUint64(&m.PhysicalAddressCounter) > 0)
		ok = ok && (m.PublicKeyFunc == nil || atomic.LoadUint64(&m.PublicKeyCounter) > 0)
		ok = ok && (m.RoleFunc == nil || atomic.LoadUint64(&m.RoleCounter) > 0)
		ok = ok && (m.ShortIDFunc == nil || atomic.LoadUint64(&m.ShortIDCounter) > 0)
		ok = ok && (m.VersionFunc == nil || atomic.LoadUint64(&m.VersionCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetGlobuleIDFunc != nil && atomic.LoadUint64(&m.GetGlobuleIDCounter) == 0 {
				m.t.Error("Expected call to NodeMock.GetGlobuleID")
			}

			if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
				m.t.Error("Expected call to NodeMock.ID")
			}

			if m.PhysicalAddressFunc != nil && atomic.LoadUint64(&m.PhysicalAddressCounter) == 0 {
				m.t.Error("Expected call to NodeMock.PhysicalAddress")
			}

			if m.PublicKeyFunc != nil && atomic.LoadUint64(&m.PublicKeyCounter) == 0 {
				m.t.Error("Expected call to NodeMock.PublicKey")
			}

			if m.RoleFunc != nil && atomic.LoadUint64(&m.RoleCounter) == 0 {
				m.t.Error("Expected call to NodeMock.Role")
			}

			if m.ShortIDFunc != nil && atomic.LoadUint64(&m.ShortIDCounter) == 0 {
				m.t.Error("Expected call to NodeMock.ShortID")
			}

			if m.VersionFunc != nil && atomic.LoadUint64(&m.VersionCounter) == 0 {
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

	if m.GetGlobuleIDFunc != nil && atomic.LoadUint64(&m.GetGlobuleIDCounter) == 0 {
		return false
	}

	if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
		return false
	}

	if m.PhysicalAddressFunc != nil && atomic.LoadUint64(&m.PhysicalAddressCounter) == 0 {
		return false
	}

	if m.PublicKeyFunc != nil && atomic.LoadUint64(&m.PublicKeyCounter) == 0 {
		return false
	}

	if m.RoleFunc != nil && atomic.LoadUint64(&m.RoleCounter) == 0 {
		return false
	}

	if m.ShortIDFunc != nil && atomic.LoadUint64(&m.ShortIDCounter) == 0 {
		return false
	}

	if m.VersionFunc != nil && atomic.LoadUint64(&m.VersionCounter) == 0 {
		return false
	}

	return true
}
