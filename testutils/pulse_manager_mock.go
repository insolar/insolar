package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseManager" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseManagerMock implements github.com/insolar/insolar/core.PulseManager
type PulseManagerMock struct {
	t minimock.Tester

	CurrentFunc       func(p context.Context) (r *core.Pulse, r1 error)
	CurrentCounter    uint64
	CurrentPreCounter uint64
	CurrentMock       mPulseManagerMockCurrent

	SetFunc       func(p context.Context, p1 core.Pulse) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mPulseManagerMockSet
}

//NewPulseManagerMock returns a mock for github.com/insolar/insolar/core.PulseManager
func NewPulseManagerMock(t minimock.Tester) *PulseManagerMock {
	m := &PulseManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CurrentMock = mPulseManagerMockCurrent{mock: m}
	m.SetMock = mPulseManagerMockSet{mock: m}

	return m
}

type mPulseManagerMockCurrent struct {
	mock             *PulseManagerMock
	mockExpectations *PulseManagerMockCurrentParams
}

//PulseManagerMockCurrentParams represents input parameters of the PulseManager.Current
type PulseManagerMockCurrentParams struct {
	p context.Context
}

//Expect sets up expected params for the PulseManager.Current
func (m *mPulseManagerMockCurrent) Expect(p context.Context) *mPulseManagerMockCurrent {
	m.mockExpectations = &PulseManagerMockCurrentParams{p}
	return m
}

//Return sets up a mock for PulseManager.Current to return Return's arguments
func (m *mPulseManagerMockCurrent) Return(r *core.Pulse, r1 error) *PulseManagerMock {
	m.mock.CurrentFunc = func(p context.Context) (*core.Pulse, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of PulseManager.Current method
func (m *mPulseManagerMockCurrent) Set(f func(p context.Context) (r *core.Pulse, r1 error)) *PulseManagerMock {
	m.mock.CurrentFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Current implements github.com/insolar/insolar/core.PulseManager interface
func (m *PulseManagerMock) Current(p context.Context) (r *core.Pulse, r1 error) {
	atomic.AddUint64(&m.CurrentPreCounter, 1)
	defer atomic.AddUint64(&m.CurrentCounter, 1)

	if m.CurrentMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.CurrentMock.mockExpectations, PulseManagerMockCurrentParams{p},
			"PulseManager.Current got unexpected parameters")

		if m.CurrentFunc == nil {

			m.t.Fatal("No results are set for the PulseManagerMock.Current")

			return
		}
	}

	if m.CurrentFunc == nil {
		m.t.Fatal("Unexpected call to PulseManagerMock.Current")
		return
	}

	return m.CurrentFunc(p)
}

//CurrentMinimockCounter returns a count of PulseManagerMock.CurrentFunc invocations
func (m *PulseManagerMock) CurrentMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CurrentCounter)
}

//CurrentMinimockPreCounter returns the value of PulseManagerMock.Current invocations
func (m *PulseManagerMock) CurrentMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CurrentPreCounter)
}

type mPulseManagerMockSet struct {
	mock             *PulseManagerMock
	mockExpectations *PulseManagerMockSetParams
}

//PulseManagerMockSetParams represents input parameters of the PulseManager.Set
type PulseManagerMockSetParams struct {
	p  context.Context
	p1 core.Pulse
}

//Expect sets up expected params for the PulseManager.Set
func (m *mPulseManagerMockSet) Expect(p context.Context, p1 core.Pulse) *mPulseManagerMockSet {
	m.mockExpectations = &PulseManagerMockSetParams{p, p1}
	return m
}

//Return sets up a mock for PulseManager.Set to return Return's arguments
func (m *mPulseManagerMockSet) Return(r error) *PulseManagerMock {
	m.mock.SetFunc = func(p context.Context, p1 core.Pulse) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of PulseManager.Set method
func (m *mPulseManagerMockSet) Set(f func(p context.Context, p1 core.Pulse) (r error)) *PulseManagerMock {
	m.mock.SetFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Set implements github.com/insolar/insolar/core.PulseManager interface
func (m *PulseManagerMock) Set(p context.Context, p1 core.Pulse) (r error) {
	atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if m.SetMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetMock.mockExpectations, PulseManagerMockSetParams{p, p1},
			"PulseManager.Set got unexpected parameters")

		if m.SetFunc == nil {

			m.t.Fatal("No results are set for the PulseManagerMock.Set")

			return
		}
	}

	if m.SetFunc == nil {
		m.t.Fatal("Unexpected call to PulseManagerMock.Set")
		return
	}

	return m.SetFunc(p, p1)
}

//SetMinimockCounter returns a count of PulseManagerMock.SetFunc invocations
func (m *PulseManagerMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of PulseManagerMock.Set invocations
func (m *PulseManagerMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseManagerMock) ValidateCallCounters() {

	if m.CurrentFunc != nil && atomic.LoadUint64(&m.CurrentCounter) == 0 {
		m.t.Fatal("Expected call to PulseManagerMock.Current")
	}

	if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
		m.t.Fatal("Expected call to PulseManagerMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseManagerMock) MinimockFinish() {

	if m.CurrentFunc != nil && atomic.LoadUint64(&m.CurrentCounter) == 0 {
		m.t.Fatal("Expected call to PulseManagerMock.Current")
	}

	if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
		m.t.Fatal("Expected call to PulseManagerMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.CurrentFunc == nil || atomic.LoadUint64(&m.CurrentCounter) > 0)
		ok = ok && (m.SetFunc == nil || atomic.LoadUint64(&m.SetCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CurrentFunc != nil && atomic.LoadUint64(&m.CurrentCounter) == 0 {
				m.t.Error("Expected call to PulseManagerMock.Current")
			}

			if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
				m.t.Error("Expected call to PulseManagerMock.Set")
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
func (m *PulseManagerMock) AllMocksCalled() bool {

	if m.CurrentFunc != nil && atomic.LoadUint64(&m.CurrentCounter) == 0 {
		return false
	}

	if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
		return false
	}

	return true
}
