package pulsar

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RPCClientWrapperFactory" can be found in github.com/insolar/insolar/pulsar
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//RPCClientWrapperFactoryMock implements github.com/insolar/insolar/pulsar.RPCClientWrapperFactory
type RPCClientWrapperFactoryMock struct {
	t minimock.Tester

	CreateWrapperFunc       func() (r RPCClientWrapper)
	CreateWrapperCounter    uint64
	CreateWrapperPreCounter uint64
	CreateWrapperMock       mRPCClientWrapperFactoryMockCreateWrapper
}

//NewRPCClientWrapperFactoryMock returns a mock for github.com/insolar/insolar/pulsar.RPCClientWrapperFactory
func NewRPCClientWrapperFactoryMock(t minimock.Tester) *RPCClientWrapperFactoryMock {
	m := &RPCClientWrapperFactoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateWrapperMock = mRPCClientWrapperFactoryMockCreateWrapper{mock: m}

	return m
}

type mRPCClientWrapperFactoryMockCreateWrapper struct {
	mock *RPCClientWrapperFactoryMock
}

//Return sets up a mock for RPCClientWrapperFactory.CreateWrapper to return Return's arguments
func (m *mRPCClientWrapperFactoryMockCreateWrapper) Return(r RPCClientWrapper) *RPCClientWrapperFactoryMock {
	m.mock.CreateWrapperFunc = func() RPCClientWrapper {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapperFactory.CreateWrapper method
func (m *mRPCClientWrapperFactoryMockCreateWrapper) Set(f func() (r RPCClientWrapper)) *RPCClientWrapperFactoryMock {
	m.mock.CreateWrapperFunc = f

	return m.mock
}

//CreateWrapper implements github.com/insolar/insolar/pulsar.RPCClientWrapperFactory interface
func (m *RPCClientWrapperFactoryMock) CreateWrapper() (r RPCClientWrapper) {
	atomic.AddUint64(&m.CreateWrapperPreCounter, 1)
	defer atomic.AddUint64(&m.CreateWrapperCounter, 1)

	if m.CreateWrapperFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperFactoryMock.CreateWrapper")
		return
	}

	return m.CreateWrapperFunc()
}

//CreateWrapperMinimockCounter returns a count of RPCClientWrapperFactoryMock.CreateWrapperFunc invocations
func (m *RPCClientWrapperFactoryMock) CreateWrapperMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateWrapperCounter)
}

//CreateWrapperMinimockPreCounter returns the value of RPCClientWrapperFactoryMock.CreateWrapper invocations
func (m *RPCClientWrapperFactoryMock) CreateWrapperMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateWrapperPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RPCClientWrapperFactoryMock) ValidateCallCounters() {

	if m.CreateWrapperFunc != nil && atomic.LoadUint64(&m.CreateWrapperCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperFactoryMock.CreateWrapper")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RPCClientWrapperFactoryMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RPCClientWrapperFactoryMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RPCClientWrapperFactoryMock) MinimockFinish() {

	if m.CreateWrapperFunc != nil && atomic.LoadUint64(&m.CreateWrapperCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperFactoryMock.CreateWrapper")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RPCClientWrapperFactoryMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RPCClientWrapperFactoryMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.CreateWrapperFunc == nil || atomic.LoadUint64(&m.CreateWrapperCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CreateWrapperFunc != nil && atomic.LoadUint64(&m.CreateWrapperCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperFactoryMock.CreateWrapper")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. require.True(mock.AllMocksCalled())
func (m *RPCClientWrapperFactoryMock) AllMocksCalled() bool {

	if m.CreateWrapperFunc != nil && atomic.LoadUint64(&m.CreateWrapperCounter) == 0 {
		return false
	}

	return true
}
