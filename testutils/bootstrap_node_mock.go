package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "BootstrapNode" can be found in github.com/insolar/insolar/core
*/
import (
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
)

//BootstrapNodeMock implements github.com/insolar/insolar/core.BootstrapNode
type BootstrapNodeMock struct {
	t minimock.Tester

	GetHostFunc       func() (r string)
	GetHostCounter    uint64
	GetHostPreCounter uint64
	GetHostMock       mBootstrapNodeMockGetHost

	GetNodeRefFunc       func() (r *core.RecordRef)
	GetNodeRefCounter    uint64
	GetNodeRefPreCounter uint64
	GetNodeRefMock       mBootstrapNodeMockGetNodeRef

	GetPublicKeyFunc       func() (r crypto.PublicKey)
	GetPublicKeyCounter    uint64
	GetPublicKeyPreCounter uint64
	GetPublicKeyMock       mBootstrapNodeMockGetPublicKey
}

//NewBootstrapNodeMock returns a mock for github.com/insolar/insolar/core.BootstrapNode
func NewBootstrapNodeMock(t minimock.Tester) *BootstrapNodeMock {
	m := &BootstrapNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetHostMock = mBootstrapNodeMockGetHost{mock: m}
	m.GetNodeRefMock = mBootstrapNodeMockGetNodeRef{mock: m}
	m.GetPublicKeyMock = mBootstrapNodeMockGetPublicKey{mock: m}

	return m
}

type mBootstrapNodeMockGetHost struct {
	mock *BootstrapNodeMock
}

//Return sets up a mock for BootstrapNode.GetHost to return Return's arguments
func (m *mBootstrapNodeMockGetHost) Return(r string) *BootstrapNodeMock {
	m.mock.GetHostFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of BootstrapNode.GetHost method
func (m *mBootstrapNodeMockGetHost) Set(f func() (r string)) *BootstrapNodeMock {
	m.mock.GetHostFunc = f

	return m.mock
}

//GetHost implements github.com/insolar/insolar/core.BootstrapNode interface
func (m *BootstrapNodeMock) GetHost() (r string) {
	atomic.AddUint64(&m.GetHostPreCounter, 1)
	defer atomic.AddUint64(&m.GetHostCounter, 1)

	if m.GetHostFunc == nil {
		m.t.Fatal("Unexpected call to BootstrapNodeMock.GetHost")
		return
	}

	return m.GetHostFunc()
}

//GetHostMinimockCounter returns a count of BootstrapNodeMock.GetHostFunc invocations
func (m *BootstrapNodeMock) GetHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHostCounter)
}

//GetHostMinimockPreCounter returns the value of BootstrapNodeMock.GetHost invocations
func (m *BootstrapNodeMock) GetHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHostPreCounter)
}

type mBootstrapNodeMockGetNodeRef struct {
	mock *BootstrapNodeMock
}

//Return sets up a mock for BootstrapNode.GetNodeRef to return Return's arguments
func (m *mBootstrapNodeMockGetNodeRef) Return(r *core.RecordRef) *BootstrapNodeMock {
	m.mock.GetNodeRefFunc = func() *core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of BootstrapNode.GetNodeRef method
func (m *mBootstrapNodeMockGetNodeRef) Set(f func() (r *core.RecordRef)) *BootstrapNodeMock {
	m.mock.GetNodeRefFunc = f

	return m.mock
}

//GetNodeRef implements github.com/insolar/insolar/core.BootstrapNode interface
func (m *BootstrapNodeMock) GetNodeRef() (r *core.RecordRef) {
	atomic.AddUint64(&m.GetNodeRefPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeRefCounter, 1)

	if m.GetNodeRefFunc == nil {
		m.t.Fatal("Unexpected call to BootstrapNodeMock.GetNodeRef")
		return
	}

	return m.GetNodeRefFunc()
}

//GetNodeRefMinimockCounter returns a count of BootstrapNodeMock.GetNodeRefFunc invocations
func (m *BootstrapNodeMock) GetNodeRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefCounter)
}

//GetNodeRefMinimockPreCounter returns the value of BootstrapNodeMock.GetNodeRef invocations
func (m *BootstrapNodeMock) GetNodeRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefPreCounter)
}

type mBootstrapNodeMockGetPublicKey struct {
	mock *BootstrapNodeMock
}

//Return sets up a mock for BootstrapNode.GetPublicKey to return Return's arguments
func (m *mBootstrapNodeMockGetPublicKey) Return(r crypto.PublicKey) *BootstrapNodeMock {
	m.mock.GetPublicKeyFunc = func() crypto.PublicKey {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of BootstrapNode.GetPublicKey method
func (m *mBootstrapNodeMockGetPublicKey) Set(f func() (r crypto.PublicKey)) *BootstrapNodeMock {
	m.mock.GetPublicKeyFunc = f

	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/core.BootstrapNode interface
func (m *BootstrapNodeMock) GetPublicKey() (r crypto.PublicKey) {
	atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if m.GetPublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to BootstrapNodeMock.GetPublicKey")
		return
	}

	return m.GetPublicKeyFunc()
}

//GetPublicKeyMinimockCounter returns a count of BootstrapNodeMock.GetPublicKeyFunc invocations
func (m *BootstrapNodeMock) GetPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyCounter)
}

//GetPublicKeyMinimockPreCounter returns the value of BootstrapNodeMock.GetPublicKey invocations
func (m *BootstrapNodeMock) GetPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BootstrapNodeMock) ValidateCallCounters() {

	if m.GetHostFunc != nil && atomic.LoadUint64(&m.GetHostCounter) == 0 {
		m.t.Fatal("Expected call to BootstrapNodeMock.GetHost")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to BootstrapNodeMock.GetNodeRef")
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to BootstrapNodeMock.GetPublicKey")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BootstrapNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *BootstrapNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *BootstrapNodeMock) MinimockFinish() {

	if m.GetHostFunc != nil && atomic.LoadUint64(&m.GetHostCounter) == 0 {
		m.t.Fatal("Expected call to BootstrapNodeMock.GetHost")
	}

	if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
		m.t.Fatal("Expected call to BootstrapNodeMock.GetNodeRef")
	}

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to BootstrapNodeMock.GetPublicKey")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *BootstrapNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *BootstrapNodeMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to BootstrapNodeMock.GetHost")
			}

			if m.GetNodeRefFunc != nil && atomic.LoadUint64(&m.GetNodeRefCounter) == 0 {
				m.t.Error("Expected call to BootstrapNodeMock.GetNodeRef")
			}

			if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
				m.t.Error("Expected call to BootstrapNodeMock.GetPublicKey")
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
func (m *BootstrapNodeMock) AllMocksCalled() bool {

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
