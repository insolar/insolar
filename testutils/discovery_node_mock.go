package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DiscoveryNode" can be found in github.com/insolar/insolar/core
*/
import (
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
)

//DiscoveryNodeMock implements github.com/insolar/insolar/core.DiscoveryNode
type DiscoveryNodeMock struct {
	t minimock.Tester

	GetHostFunc       func() (r string)
	GetHostCounter    uint64
	GetHostPreCounter uint64
	GetHostMock       mDiscoveryNodeMockGetHost

	GetNodeRefFunc       func() (r *core.RecordRef)
	GetNodeRefCounter    uint64
	GetNodeRefPreCounter uint64
	GetNodeRefMock       mDiscoveryNodeMockGetNodeRef

	GetPublicKeyFunc       func() (r crypto.PublicKey)
	GetPublicKeyCounter    uint64
	GetPublicKeyPreCounter uint64
	GetPublicKeyMock       mDiscoveryNodeMockGetPublicKey
}

//NewDiscoveryNodeMock returns a mock for github.com/insolar/insolar/core.DiscoveryNode
func NewDiscoveryNodeMock(t minimock.Tester) *DiscoveryNodeMock {
	m := &DiscoveryNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetHostMock = mDiscoveryNodeMockGetHost{mock: m}
	m.GetNodeRefMock = mDiscoveryNodeMockGetNodeRef{mock: m}
	m.GetPublicKeyMock = mDiscoveryNodeMockGetPublicKey{mock: m}

	return m
}

type mDiscoveryNodeMockGetHost struct {
	mock *DiscoveryNodeMock
}

//Return sets up a mock for DiscoveryNode.GetHost to return Return's arguments
func (m *mDiscoveryNodeMockGetHost) Return(r string) *DiscoveryNodeMock {
	m.mock.GetHostFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of DiscoveryNode.GetHost method
func (m *mDiscoveryNodeMockGetHost) Set(f func() (r string)) *DiscoveryNodeMock {
	m.mock.GetHostFunc = f

	return m.mock
}

//GetHost implements github.com/insolar/insolar/core.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetHost() (r string) {
	atomic.AddUint64(&m.GetHostPreCounter, 1)
	defer atomic.AddUint64(&m.GetHostCounter, 1)

	if m.GetHostFunc == nil {
		m.t.Fatal("Unexpected call to DiscoveryNodeMock.GetHost")
		return
	}

	return m.GetHostFunc()
}

//GetHostMinimockCounter returns a count of DiscoveryNodeMock.GetHostFunc invocations
func (m *DiscoveryNodeMock) GetHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHostCounter)
}

//GetHostMinimockPreCounter returns the value of DiscoveryNodeMock.GetHost invocations
func (m *DiscoveryNodeMock) GetHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHostPreCounter)
}

type mDiscoveryNodeMockGetNodeRef struct {
	mock *DiscoveryNodeMock
}

//Return sets up a mock for DiscoveryNode.GetNodeRef to return Return's arguments
func (m *mDiscoveryNodeMockGetNodeRef) Return(r *core.RecordRef) *DiscoveryNodeMock {
	m.mock.GetNodeRefFunc = func() *core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of DiscoveryNode.GetNodeRef method
func (m *mDiscoveryNodeMockGetNodeRef) Set(f func() (r *core.RecordRef)) *DiscoveryNodeMock {
	m.mock.GetNodeRefFunc = f

	return m.mock
}

//GetNodeRef implements github.com/insolar/insolar/core.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetNodeRef() (r *core.RecordRef) {
	atomic.AddUint64(&m.GetNodeRefPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeRefCounter, 1)

	if m.GetNodeRefFunc == nil {
		m.t.Fatal("Unexpected call to DiscoveryNodeMock.GetNodeRef")
		return
	}

	return m.GetNodeRefFunc()
}

//GetNodeRefMinimockCounter returns a count of DiscoveryNodeMock.GetNodeRefFunc invocations
func (m *DiscoveryNodeMock) GetNodeRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefCounter)
}

//GetNodeRefMinimockPreCounter returns the value of DiscoveryNodeMock.GetNodeRef invocations
func (m *DiscoveryNodeMock) GetNodeRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefPreCounter)
}

type mDiscoveryNodeMockGetPublicKey struct {
	mock *DiscoveryNodeMock
}

//Return sets up a mock for DiscoveryNode.GetPublicKey to return Return's arguments
func (m *mDiscoveryNodeMockGetPublicKey) Return(r crypto.PublicKey) *DiscoveryNodeMock {
	m.mock.GetPublicKeyFunc = func() crypto.PublicKey {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of DiscoveryNode.GetPublicKey method
func (m *mDiscoveryNodeMockGetPublicKey) Set(f func() (r crypto.PublicKey)) *DiscoveryNodeMock {
	m.mock.GetPublicKeyFunc = f

	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/core.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetPublicKey() (r crypto.PublicKey) {
	atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if m.GetPublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to DiscoveryNodeMock.GetPublicKey")
		return
	}

	return m.GetPublicKeyFunc()
}

//GetPublicKeyMinimockCounter returns a count of DiscoveryNodeMock.GetPublicKeyFunc invocations
func (m *DiscoveryNodeMock) GetPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyCounter)
}

//GetPublicKeyMinimockPreCounter returns the value of DiscoveryNodeMock.GetPublicKey invocations
func (m *DiscoveryNodeMock) GetPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DiscoveryNodeMock) ValidateCallCounters() {

	if m.GetHostFunc != nil && atomic.LoadUint64(&m.GetHostCounter) == 0 {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetHost")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetNodeRef")
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetPublicKey")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DiscoveryNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DiscoveryNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DiscoveryNodeMock) MinimockFinish() {

	if m.GetHostFunc != nil && atomic.LoadUint64(&m.GetHostCounter) == 0 {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetHost")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetNodeRef")
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetPublicKey")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DiscoveryNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DiscoveryNodeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetHostFunc == nil || atomic.LoadUint64(&m.GetHostCounter) > 0)
		ok = ok && (m.GetNodeRefFunc == nil || atomic.LoadUint64(&m.GetNodeRefCounter) > 0)
		ok = ok && (m.GetPublicKeyFunc == nil || atomic.LoadUint64(&m.GetPublicKeyCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetHostFunc != nil && atomic.LoadUint64(&m.GetHostCounter) == 0 {
				m.t.Error("Expected call to DiscoveryNodeMock.GetHost")
			}

			if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
				m.t.Error("Expected call to DiscoveryNodeMock.GetNodeRef")
			}

			if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
				m.t.Error("Expected call to DiscoveryNodeMock.GetPublicKey")
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
func (m *DiscoveryNodeMock) AllMocksCalled() bool {

	if m.GetHostFunc != nil && atomic.LoadUint64(&m.GetHostCounter) == 0 {
		return false
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		return false
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		return false
	}

	return true
}
