package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SwitcherWorkAround" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//SwitcherWorkAroundMock implements github.com/insolar/insolar/core.SwitcherWorkAround
type SwitcherWorkAroundMock struct {
	t minimock.Tester

	IsBootstrappedFunc       func() (r bool)
	IsBootstrappedCounter    uint64
	IsBootstrappedPreCounter uint64
	IsBootstrappedMock       mSwitcherWorkAroundMockIsBootstrapped

	SetIsBootstrappedFunc       func(p bool)
	SetIsBootstrappedCounter    uint64
	SetIsBootstrappedPreCounter uint64
	SetIsBootstrappedMock       mSwitcherWorkAroundMockSetIsBootstrapped
}

//NewSwitcherWorkAroundMock returns a mock for github.com/insolar/insolar/core.SwitcherWorkAround
func NewSwitcherWorkAroundMock(t minimock.Tester) *SwitcherWorkAroundMock {
	m := &SwitcherWorkAroundMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IsBootstrappedMock = mSwitcherWorkAroundMockIsBootstrapped{mock: m}
	m.SetIsBootstrappedMock = mSwitcherWorkAroundMockSetIsBootstrapped{mock: m}

	return m
}

type mSwitcherWorkAroundMockIsBootstrapped struct {
	mock *SwitcherWorkAroundMock
}

//Return sets up a mock for SwitcherWorkAround.IsBootstrapped to return Return's arguments
func (m *mSwitcherWorkAroundMockIsBootstrapped) Return(r bool) *SwitcherWorkAroundMock {
	m.mock.IsBootstrappedFunc = func() bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of SwitcherWorkAround.IsBootstrapped method
func (m *mSwitcherWorkAroundMockIsBootstrapped) Set(f func() (r bool)) *SwitcherWorkAroundMock {
	m.mock.IsBootstrappedFunc = f

	return m.mock
}

//IsBootstrapped implements github.com/insolar/insolar/core.SwitcherWorkAround interface
func (m *SwitcherWorkAroundMock) IsBootstrapped() (r bool) {
	atomic.AddUint64(&m.IsBootstrappedPreCounter, 1)
	defer atomic.AddUint64(&m.IsBootstrappedCounter, 1)

	if m.IsBootstrappedFunc == nil {
		m.t.Fatal("Unexpected call to SwitcherWorkAroundMock.IsBootstrapped")
		return
	}

	return m.IsBootstrappedFunc()
}

//IsBootstrappedMinimockCounter returns a count of SwitcherWorkAroundMock.IsBootstrappedFunc invocations
func (m *SwitcherWorkAroundMock) IsBootstrappedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsBootstrappedCounter)
}

//IsBootstrappedMinimockPreCounter returns the value of SwitcherWorkAroundMock.IsBootstrapped invocations
func (m *SwitcherWorkAroundMock) IsBootstrappedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsBootstrappedPreCounter)
}

type mSwitcherWorkAroundMockSetIsBootstrapped struct {
	mock             *SwitcherWorkAroundMock
	mockExpectations *SwitcherWorkAroundMockSetIsBootstrappedParams
}

//SwitcherWorkAroundMockSetIsBootstrappedParams represents input parameters of the SwitcherWorkAround.SetIsBootstrapped
type SwitcherWorkAroundMockSetIsBootstrappedParams struct {
	p bool
}

//Expect sets up expected params for the SwitcherWorkAround.SetIsBootstrapped
func (m *mSwitcherWorkAroundMockSetIsBootstrapped) Expect(p bool) *mSwitcherWorkAroundMockSetIsBootstrapped {
	m.mockExpectations = &SwitcherWorkAroundMockSetIsBootstrappedParams{p}
	return m
}

//Return sets up a mock for SwitcherWorkAround.SetIsBootstrapped to return Return's arguments
func (m *mSwitcherWorkAroundMockSetIsBootstrapped) Return() *SwitcherWorkAroundMock {
	m.mock.SetIsBootstrappedFunc = func(p bool) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of SwitcherWorkAround.SetIsBootstrapped method
func (m *mSwitcherWorkAroundMockSetIsBootstrapped) Set(f func(p bool)) *SwitcherWorkAroundMock {
	m.mock.SetIsBootstrappedFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetIsBootstrapped implements github.com/insolar/insolar/core.SwitcherWorkAround interface
func (m *SwitcherWorkAroundMock) SetIsBootstrapped(p bool) {
	atomic.AddUint64(&m.SetIsBootstrappedPreCounter, 1)
	defer atomic.AddUint64(&m.SetIsBootstrappedCounter, 1)

	if m.SetIsBootstrappedMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetIsBootstrappedMock.mockExpectations, SwitcherWorkAroundMockSetIsBootstrappedParams{p},
			"SwitcherWorkAround.SetIsBootstrapped got unexpected parameters")

		if m.SetIsBootstrappedFunc == nil {

			m.t.Fatal("No results are set for the SwitcherWorkAroundMock.SetIsBootstrapped")

			return
		}
	}

	if m.SetIsBootstrappedFunc == nil {
		m.t.Fatal("Unexpected call to SwitcherWorkAroundMock.SetIsBootstrapped")
		return
	}

	m.SetIsBootstrappedFunc(p)
}

//SetIsBootstrappedMinimockCounter returns a count of SwitcherWorkAroundMock.SetIsBootstrappedFunc invocations
func (m *SwitcherWorkAroundMock) SetIsBootstrappedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetIsBootstrappedCounter)
}

//SetIsBootstrappedMinimockPreCounter returns the value of SwitcherWorkAroundMock.SetIsBootstrapped invocations
func (m *SwitcherWorkAroundMock) SetIsBootstrappedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetIsBootstrappedPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SwitcherWorkAroundMock) ValidateCallCounters() {

	if m.IsBootstrappedFunc != nil && atomic.LoadUint64(&m.IsBootstrappedCounter) == 0 {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.IsBootstrapped")
	}

	if m.SetIsBootstrappedFunc != nil && atomic.LoadUint64(&m.SetIsBootstrappedCounter) == 0 {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.SetIsBootstrapped")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SwitcherWorkAroundMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SwitcherWorkAroundMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SwitcherWorkAroundMock) MinimockFinish() {

	if m.IsBootstrappedFunc != nil && atomic.LoadUint64(&m.IsBootstrappedCounter) == 0 {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.IsBootstrapped")
	}

	if m.SetIsBootstrappedFunc != nil && atomic.LoadUint64(&m.SetIsBootstrappedCounter) == 0 {
		m.t.Fatal("Expected call to SwitcherWorkAroundMock.SetIsBootstrapped")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SwitcherWorkAroundMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SwitcherWorkAroundMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.IsBootstrappedFunc == nil || atomic.LoadUint64(&m.IsBootstrappedCounter) > 0)
		ok = ok && (m.SetIsBootstrappedFunc == nil || atomic.LoadUint64(&m.SetIsBootstrappedCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.IsBootstrappedFunc != nil && atomic.LoadUint64(&m.IsBootstrappedCounter) == 0 {
				m.t.Error("Expected call to SwitcherWorkAroundMock.IsBootstrapped")
			}

			if m.SetIsBootstrappedFunc != nil && atomic.LoadUint64(&m.SetIsBootstrappedCounter) == 0 {
				m.t.Error("Expected call to SwitcherWorkAroundMock.SetIsBootstrapped")
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
func (m *SwitcherWorkAroundMock) AllMocksCalled() bool {

	if m.IsBootstrappedFunc != nil && atomic.LoadUint64(&m.IsBootstrappedCounter) == 0 {
		return false
	}

	if m.SetIsBootstrappedFunc != nil && atomic.LoadUint64(&m.SetIsBootstrappedCounter) == 0 {
		return false
	}

	return true
}
