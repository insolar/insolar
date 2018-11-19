package pulsar

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RPCClientWrapper" can be found in github.com/insolar/insolar/pulsar
*/
import (
	"net/rpc"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/configuration"

	testify_assert "github.com/stretchr/testify/assert"
)

//RPCClientWrapperMock implements github.com/insolar/insolar/pulsar.RPCClientWrapper
type RPCClientWrapperMock struct {
	t minimock.Tester

	CloseFunc       func() (r error)
	CloseCounter    uint64
	ClosePreCounter uint64
	CloseMock       mRPCClientWrapperMockClose

	CreateConnectionFunc       func(p configuration.ConnectionType, p1 string) (r error)
	CreateConnectionCounter    uint64
	CreateConnectionPreCounter uint64
	CreateConnectionMock       mRPCClientWrapperMockCreateConnection

	GoFunc       func(p string, p1 interface{}, p2 interface{}, p3 chan *rpc.Call) (r *rpc.Call)
	GoCounter    uint64
	GoPreCounter uint64
	GoMock       mRPCClientWrapperMockGo

	IsInitialisedFunc       func() (r bool)
	IsInitialisedCounter    uint64
	IsInitialisedPreCounter uint64
	IsInitialisedMock       mRPCClientWrapperMockIsInitialised

	LockFunc       func()
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mRPCClientWrapperMockLock

	ResetClientFunc       func()
	ResetClientCounter    uint64
	ResetClientPreCounter uint64
	ResetClientMock       mRPCClientWrapperMockResetClient

	UnlockFunc       func()
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mRPCClientWrapperMockUnlock
}

//NewRPCClientWrapperMock returns a mock for github.com/insolar/insolar/pulsar.RPCClientWrapper
func NewRPCClientWrapperMock(t minimock.Tester) *RPCClientWrapperMock {
	m := &RPCClientWrapperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseMock = mRPCClientWrapperMockClose{mock: m}
	m.CreateConnectionMock = mRPCClientWrapperMockCreateConnection{mock: m}
	m.GoMock = mRPCClientWrapperMockGo{mock: m}
	m.IsInitialisedMock = mRPCClientWrapperMockIsInitialised{mock: m}
	m.LockMock = mRPCClientWrapperMockLock{mock: m}
	m.ResetClientMock = mRPCClientWrapperMockResetClient{mock: m}
	m.UnlockMock = mRPCClientWrapperMockUnlock{mock: m}

	return m
}

type mRPCClientWrapperMockClose struct {
	mock *RPCClientWrapperMock
}

//Return sets up a mock for RPCClientWrapper.Close to return Return's arguments
func (m *mRPCClientWrapperMockClose) Return(r error) *RPCClientWrapperMock {
	m.mock.CloseFunc = func() error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapper.Close method
func (m *mRPCClientWrapperMockClose) Set(f func() (r error)) *RPCClientWrapperMock {
	m.mock.CloseFunc = f

	return m.mock
}

//Close implements github.com/insolar/insolar/pulsar.RPCClientWrapper interface
func (m *RPCClientWrapperMock) Close() (r error) {
	atomic.AddUint64(&m.ClosePreCounter, 1)
	defer atomic.AddUint64(&m.CloseCounter, 1)

	if m.CloseFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperMock.Close")
		return
	}

	return m.CloseFunc()
}

//CloseMinimockCounter returns a count of RPCClientWrapperMock.CloseFunc invocations
func (m *RPCClientWrapperMock) CloseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseCounter)
}

//CloseMinimockPreCounter returns the value of RPCClientWrapperMock.Close invocations
func (m *RPCClientWrapperMock) CloseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClosePreCounter)
}

type mRPCClientWrapperMockCreateConnection struct {
	mock             *RPCClientWrapperMock
	mockExpectations *RPCClientWrapperMockCreateConnectionParams
}

//RPCClientWrapperMockCreateConnectionParams represents input parameters of the RPCClientWrapper.CreateConnection
type RPCClientWrapperMockCreateConnectionParams struct {
	p  configuration.ConnectionType
	p1 string
}

//Expect sets up expected params for the RPCClientWrapper.CreateConnection
func (m *mRPCClientWrapperMockCreateConnection) Expect(p configuration.ConnectionType, p1 string) *mRPCClientWrapperMockCreateConnection {
	m.mockExpectations = &RPCClientWrapperMockCreateConnectionParams{p, p1}
	return m
}

//Return sets up a mock for RPCClientWrapper.CreateConnection to return Return's arguments
func (m *mRPCClientWrapperMockCreateConnection) Return(r error) *RPCClientWrapperMock {
	m.mock.CreateConnectionFunc = func(p configuration.ConnectionType, p1 string) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapper.CreateConnection method
func (m *mRPCClientWrapperMockCreateConnection) Set(f func(p configuration.ConnectionType, p1 string) (r error)) *RPCClientWrapperMock {
	m.mock.CreateConnectionFunc = f
	m.mockExpectations = nil
	return m.mock
}

//CreateConnection implements github.com/insolar/insolar/pulsar.RPCClientWrapper interface
func (m *RPCClientWrapperMock) CreateConnection(p configuration.ConnectionType, p1 string) (r error) {
	atomic.AddUint64(&m.CreateConnectionPreCounter, 1)
	defer atomic.AddUint64(&m.CreateConnectionCounter, 1)

	if m.CreateConnectionMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.CreateConnectionMock.mockExpectations, RPCClientWrapperMockCreateConnectionParams{p, p1},
			"RPCClientWrapper.CreateConnection got unexpected parameters")

		if m.CreateConnectionFunc == nil {

			m.t.Fatal("No results are set for the RPCClientWrapperMock.CreateConnection")

			return
		}
	}

	if m.CreateConnectionFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperMock.CreateConnection")
		return
	}

	return m.CreateConnectionFunc(p, p1)
}

//CreateConnectionMinimockCounter returns a count of RPCClientWrapperMock.CreateConnectionFunc invocations
func (m *RPCClientWrapperMock) CreateConnectionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateConnectionCounter)
}

//CreateConnectionMinimockPreCounter returns the value of RPCClientWrapperMock.CreateConnection invocations
func (m *RPCClientWrapperMock) CreateConnectionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateConnectionPreCounter)
}

type mRPCClientWrapperMockGo struct {
	mock             *RPCClientWrapperMock
	mockExpectations *RPCClientWrapperMockGoParams
}

//RPCClientWrapperMockGoParams represents input parameters of the RPCClientWrapper.Go
type RPCClientWrapperMockGoParams struct {
	p  string
	p1 interface{}
	p2 interface{}
	p3 chan *rpc.Call
}

//Expect sets up expected params for the RPCClientWrapper.Go
func (m *mRPCClientWrapperMockGo) Expect(p string, p1 interface{}, p2 interface{}, p3 chan *rpc.Call) *mRPCClientWrapperMockGo {
	m.mockExpectations = &RPCClientWrapperMockGoParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for RPCClientWrapper.Go to return Return's arguments
func (m *mRPCClientWrapperMockGo) Return(r *rpc.Call) *RPCClientWrapperMock {
	m.mock.GoFunc = func(p string, p1 interface{}, p2 interface{}, p3 chan *rpc.Call) *rpc.Call {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapper.Go method
func (m *mRPCClientWrapperMockGo) Set(f func(p string, p1 interface{}, p2 interface{}, p3 chan *rpc.Call) (r *rpc.Call)) *RPCClientWrapperMock {
	m.mock.GoFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Go implements github.com/insolar/insolar/pulsar.RPCClientWrapper interface
func (m *RPCClientWrapperMock) Go(p string, p1 interface{}, p2 interface{}, p3 chan *rpc.Call) (r *rpc.Call) {
	atomic.AddUint64(&m.GoPreCounter, 1)
	defer atomic.AddUint64(&m.GoCounter, 1)

	if m.GoMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GoMock.mockExpectations, RPCClientWrapperMockGoParams{p, p1, p2, p3},
			"RPCClientWrapper.Go got unexpected parameters")

		if m.GoFunc == nil {

			m.t.Fatal("No results are set for the RPCClientWrapperMock.Go")

			return
		}
	}

	if m.GoFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperMock.Go")
		return
	}

	return m.GoFunc(p, p1, p2, p3)
}

//GoMinimockCounter returns a count of RPCClientWrapperMock.GoFunc invocations
func (m *RPCClientWrapperMock) GoMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GoCounter)
}

//GoMinimockPreCounter returns the value of RPCClientWrapperMock.Go invocations
func (m *RPCClientWrapperMock) GoMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GoPreCounter)
}

type mRPCClientWrapperMockIsInitialised struct {
	mock *RPCClientWrapperMock
}

//Return sets up a mock for RPCClientWrapper.IsInitialised to return Return's arguments
func (m *mRPCClientWrapperMockIsInitialised) Return(r bool) *RPCClientWrapperMock {
	m.mock.IsInitialisedFunc = func() bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapper.IsInitialised method
func (m *mRPCClientWrapperMockIsInitialised) Set(f func() (r bool)) *RPCClientWrapperMock {
	m.mock.IsInitialisedFunc = f

	return m.mock
}

//IsInitialised implements github.com/insolar/insolar/pulsar.RPCClientWrapper interface
func (m *RPCClientWrapperMock) IsInitialised() (r bool) {
	atomic.AddUint64(&m.IsInitialisedPreCounter, 1)
	defer atomic.AddUint64(&m.IsInitialisedCounter, 1)

	if m.IsInitialisedFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperMock.IsInitialised")
		return
	}

	return m.IsInitialisedFunc()
}

//IsInitialisedMinimockCounter returns a count of RPCClientWrapperMock.IsInitialisedFunc invocations
func (m *RPCClientWrapperMock) IsInitialisedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsInitialisedCounter)
}

//IsInitialisedMinimockPreCounter returns the value of RPCClientWrapperMock.IsInitialised invocations
func (m *RPCClientWrapperMock) IsInitialisedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsInitialisedPreCounter)
}

type mRPCClientWrapperMockLock struct {
	mock *RPCClientWrapperMock
}

//Return sets up a mock for RPCClientWrapper.Lock to return Return's arguments
func (m *mRPCClientWrapperMockLock) Return() *RPCClientWrapperMock {
	m.mock.LockFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapper.Lock method
func (m *mRPCClientWrapperMockLock) Set(f func()) *RPCClientWrapperMock {
	m.mock.LockFunc = f

	return m.mock
}

//Lock implements github.com/insolar/insolar/pulsar.RPCClientWrapper interface
func (m *RPCClientWrapperMock) Lock() {
	atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if m.LockFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperMock.Lock")
		return
	}

	m.LockFunc()
}

//LockMinimockCounter returns a count of RPCClientWrapperMock.LockFunc invocations
func (m *RPCClientWrapperMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of RPCClientWrapperMock.Lock invocations
func (m *RPCClientWrapperMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

type mRPCClientWrapperMockResetClient struct {
	mock *RPCClientWrapperMock
}

//Return sets up a mock for RPCClientWrapper.ResetClient to return Return's arguments
func (m *mRPCClientWrapperMockResetClient) Return() *RPCClientWrapperMock {
	m.mock.ResetClientFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapper.ResetClient method
func (m *mRPCClientWrapperMockResetClient) Set(f func()) *RPCClientWrapperMock {
	m.mock.ResetClientFunc = f

	return m.mock
}

//ResetClient implements github.com/insolar/insolar/pulsar.RPCClientWrapper interface
func (m *RPCClientWrapperMock) ResetClient() {
	atomic.AddUint64(&m.ResetClientPreCounter, 1)
	defer atomic.AddUint64(&m.ResetClientCounter, 1)

	if m.ResetClientFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperMock.ResetClient")
		return
	}

	m.ResetClientFunc()
}

//ResetClientMinimockCounter returns a count of RPCClientWrapperMock.ResetClientFunc invocations
func (m *RPCClientWrapperMock) ResetClientMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResetClientCounter)
}

//ResetClientMinimockPreCounter returns the value of RPCClientWrapperMock.ResetClient invocations
func (m *RPCClientWrapperMock) ResetClientMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResetClientPreCounter)
}

type mRPCClientWrapperMockUnlock struct {
	mock *RPCClientWrapperMock
}

//Return sets up a mock for RPCClientWrapper.Unlock to return Return's arguments
func (m *mRPCClientWrapperMockUnlock) Return() *RPCClientWrapperMock {
	m.mock.UnlockFunc = func() {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of RPCClientWrapper.Unlock method
func (m *mRPCClientWrapperMockUnlock) Set(f func()) *RPCClientWrapperMock {
	m.mock.UnlockFunc = f

	return m.mock
}

//Unlock implements github.com/insolar/insolar/pulsar.RPCClientWrapper interface
func (m *RPCClientWrapperMock) Unlock() {
	atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if m.UnlockFunc == nil {
		m.t.Fatal("Unexpected call to RPCClientWrapperMock.Unlock")
		return
	}

	m.UnlockFunc()
}

//UnlockMinimockCounter returns a count of RPCClientWrapperMock.UnlockFunc invocations
func (m *RPCClientWrapperMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of RPCClientWrapperMock.Unlock invocations
func (m *RPCClientWrapperMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RPCClientWrapperMock) ValidateCallCounters() {

	if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Close")
	}

	if m.CreateConnectionFunc != nil && atomic.LoadUint64(&m.CreateConnectionCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.CreateConnection")
	}

	if m.GoFunc != nil && atomic.LoadUint64(&m.GoCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Go")
	}

	if m.IsInitialisedFunc != nil && atomic.LoadUint64(&m.IsInitialisedCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.IsInitialised")
	}

	if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Lock")
	}

	if m.ResetClientFunc != nil && atomic.LoadUint64(&m.ResetClientCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.ResetClient")
	}

	if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Unlock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RPCClientWrapperMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RPCClientWrapperMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RPCClientWrapperMock) MinimockFinish() {

	if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Close")
	}

	if m.CreateConnectionFunc != nil && atomic.LoadUint64(&m.CreateConnectionCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.CreateConnection")
	}

	if m.GoFunc != nil && atomic.LoadUint64(&m.GoCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Go")
	}

	if m.IsInitialisedFunc != nil && atomic.LoadUint64(&m.IsInitialisedCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.IsInitialised")
	}

	if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Lock")
	}

	if m.ResetClientFunc != nil && atomic.LoadUint64(&m.ResetClientCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.ResetClient")
	}

	if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
		m.t.Fatal("Expected call to RPCClientWrapperMock.Unlock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RPCClientWrapperMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RPCClientWrapperMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.CloseFunc == nil || atomic.LoadUint64(&m.CloseCounter) > 0)
		ok = ok && (m.CreateConnectionFunc == nil || atomic.LoadUint64(&m.CreateConnectionCounter) > 0)
		ok = ok && (m.GoFunc == nil || atomic.LoadUint64(&m.GoCounter) > 0)
		ok = ok && (m.IsInitialisedFunc == nil || atomic.LoadUint64(&m.IsInitialisedCounter) > 0)
		ok = ok && (m.LockFunc == nil || atomic.LoadUint64(&m.LockCounter) > 0)
		ok = ok && (m.ResetClientFunc == nil || atomic.LoadUint64(&m.ResetClientCounter) > 0)
		ok = ok && (m.UnlockFunc == nil || atomic.LoadUint64(&m.UnlockCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperMock.Close")
			}

			if m.CreateConnectionFunc != nil && atomic.LoadUint64(&m.CreateConnectionCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperMock.CreateConnection")
			}

			if m.GoFunc != nil && atomic.LoadUint64(&m.GoCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperMock.Go")
			}

			if m.IsInitialisedFunc != nil && atomic.LoadUint64(&m.IsInitialisedCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperMock.IsInitialised")
			}

			if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperMock.Lock")
			}

			if m.ResetClientFunc != nil && atomic.LoadUint64(&m.ResetClientCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperMock.ResetClient")
			}

			if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
				m.t.Error("Expected call to RPCClientWrapperMock.Unlock")
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
func (m *RPCClientWrapperMock) AllMocksCalled() bool {

	if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
		return false
	}

	if m.CreateConnectionFunc != nil && atomic.LoadUint64(&m.CreateConnectionCounter) == 0 {
		return false
	}

	if m.GoFunc != nil && atomic.LoadUint64(&m.GoCounter) == 0 {
		return false
	}

	if m.IsInitialisedFunc != nil && atomic.LoadUint64(&m.IsInitialisedCounter) == 0 {
		return false
	}

	if m.LockFunc != nil && atomic.LoadUint64(&m.LockCounter) == 0 {
		return false
	}

	if m.ResetClientFunc != nil && atomic.LoadUint64(&m.ResetClientCounter) == 0 {
		return false
	}

	if m.UnlockFunc != nil && atomic.LoadUint64(&m.UnlockCounter) == 0 {
		return false
	}

	return true
}
