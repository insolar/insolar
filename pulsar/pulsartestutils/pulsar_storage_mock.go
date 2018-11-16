package pulsartestutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulsarStorage" can be found in github.com/insolar/insolar/pulsar/storage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulsarStorageMock implements github.com/insolar/insolar/pulsar/storage.PulsarStorage
type PulsarStorageMock struct {
	t minimock.Tester

	CloseFunc       func() (r error)
	CloseCounter    uint64
	ClosePreCounter uint64
	CloseMock       mPulsarStorageMockClose

	GetLastPulseFunc       func() (r *core.Pulse, r1 error)
	GetLastPulseCounter    uint64
	GetLastPulsePreCounter uint64
	GetLastPulseMock       mPulsarStorageMockGetLastPulse

	SavePulseFunc       func(p *core.Pulse) (r error)
	SavePulseCounter    uint64
	SavePulsePreCounter uint64
	SavePulseMock       mPulsarStorageMockSavePulse

	SetLastPulseFunc       func(p *core.Pulse) (r error)
	SetLastPulseCounter    uint64
	SetLastPulsePreCounter uint64
	SetLastPulseMock       mPulsarStorageMockSetLastPulse
}

//NewPulsarStorageMock returns a mock for github.com/insolar/insolar/pulsar/storage.PulsarStorage
func NewPulsarStorageMock(t minimock.Tester) *PulsarStorageMock {
	m := &PulsarStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseMock = mPulsarStorageMockClose{mock: m}
	m.GetLastPulseMock = mPulsarStorageMockGetLastPulse{mock: m}
	m.SavePulseMock = mPulsarStorageMockSavePulse{mock: m}
	m.SetLastPulseMock = mPulsarStorageMockSetLastPulse{mock: m}

	return m
}

type mPulsarStorageMockClose struct {
	mock *PulsarStorageMock
}

//Return sets up a mock for PulsarStorage.Close to return Return's arguments
func (m *mPulsarStorageMockClose) Return(r error) *PulsarStorageMock {
	m.mock.CloseFunc = func() error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of PulsarStorage.Close method
func (m *mPulsarStorageMockClose) Set(f func() (r error)) *PulsarStorageMock {
	m.mock.CloseFunc = f

	return m.mock
}

//Close implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) Close() (r error) {
	atomic.AddUint64(&m.ClosePreCounter, 1)
	defer atomic.AddUint64(&m.CloseCounter, 1)

	if m.CloseFunc == nil {
		m.t.Fatal("Unexpected call to PulsarStorageMock.Close")
		return
	}

	return m.CloseFunc()
}

//CloseMinimockCounter returns a count of PulsarStorageMock.CloseFunc invocations
func (m *PulsarStorageMock) CloseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseCounter)
}

//CloseMinimockPreCounter returns the value of PulsarStorageMock.Close invocations
func (m *PulsarStorageMock) CloseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClosePreCounter)
}

type mPulsarStorageMockGetLastPulse struct {
	mock *PulsarStorageMock
}

//Return sets up a mock for PulsarStorage.GetLastPulse to return Return's arguments
func (m *mPulsarStorageMockGetLastPulse) Return(r *core.Pulse, r1 error) *PulsarStorageMock {
	m.mock.GetLastPulseFunc = func() (*core.Pulse, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of PulsarStorage.GetLastPulse method
func (m *mPulsarStorageMockGetLastPulse) Set(f func() (r *core.Pulse, r1 error)) *PulsarStorageMock {
	m.mock.GetLastPulseFunc = f

	return m.mock
}

//GetLastPulse implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) GetLastPulse() (r *core.Pulse, r1 error) {
	atomic.AddUint64(&m.GetLastPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetLastPulseCounter, 1)

	if m.GetLastPulseFunc == nil {
		m.t.Fatal("Unexpected call to PulsarStorageMock.GetLastPulse")
		return
	}

	return m.GetLastPulseFunc()
}

//GetLastPulseMinimockCounter returns a count of PulsarStorageMock.GetLastPulseFunc invocations
func (m *PulsarStorageMock) GetLastPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetLastPulseCounter)
}

//GetLastPulseMinimockPreCounter returns the value of PulsarStorageMock.GetLastPulse invocations
func (m *PulsarStorageMock) GetLastPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetLastPulsePreCounter)
}

type mPulsarStorageMockSavePulse struct {
	mock             *PulsarStorageMock
	mockExpectations *PulsarStorageMockSavePulseParams
}

//PulsarStorageMockSavePulseParams represents input parameters of the PulsarStorage.SavePulse
type PulsarStorageMockSavePulseParams struct {
	p *core.Pulse
}

//Expect sets up expected params for the PulsarStorage.SavePulse
func (m *mPulsarStorageMockSavePulse) Expect(p *core.Pulse) *mPulsarStorageMockSavePulse {
	m.mockExpectations = &PulsarStorageMockSavePulseParams{p}
	return m
}

//Return sets up a mock for PulsarStorage.SavePulse to return Return's arguments
func (m *mPulsarStorageMockSavePulse) Return(r error) *PulsarStorageMock {
	m.mock.SavePulseFunc = func(p *core.Pulse) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of PulsarStorage.SavePulse method
func (m *mPulsarStorageMockSavePulse) Set(f func(p *core.Pulse) (r error)) *PulsarStorageMock {
	m.mock.SavePulseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SavePulse implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) SavePulse(p *core.Pulse) (r error) {
	atomic.AddUint64(&m.SavePulsePreCounter, 1)
	defer atomic.AddUint64(&m.SavePulseCounter, 1)

	if m.SavePulseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SavePulseMock.mockExpectations, PulsarStorageMockSavePulseParams{p},
			"PulsarStorage.SavePulse got unexpected parameters")

		if m.SavePulseFunc == nil {

			m.t.Fatal("No results are set for the PulsarStorageMock.SavePulse")

			return
		}
	}

	if m.SavePulseFunc == nil {
		m.t.Fatal("Unexpected call to PulsarStorageMock.SavePulse")
		return
	}

	return m.SavePulseFunc(p)
}

//SavePulseMinimockCounter returns a count of PulsarStorageMock.SavePulseFunc invocations
func (m *PulsarStorageMock) SavePulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SavePulseCounter)
}

//SavePulseMinimockPreCounter returns the value of PulsarStorageMock.SavePulse invocations
func (m *PulsarStorageMock) SavePulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SavePulsePreCounter)
}

type mPulsarStorageMockSetLastPulse struct {
	mock             *PulsarStorageMock
	mockExpectations *PulsarStorageMockSetLastPulseParams
}

//PulsarStorageMockSetLastPulseParams represents input parameters of the PulsarStorage.SetLastPulse
type PulsarStorageMockSetLastPulseParams struct {
	p *core.Pulse
}

//Expect sets up expected params for the PulsarStorage.SetLastPulse
func (m *mPulsarStorageMockSetLastPulse) Expect(p *core.Pulse) *mPulsarStorageMockSetLastPulse {
	m.mockExpectations = &PulsarStorageMockSetLastPulseParams{p}
	return m
}

//Return sets up a mock for PulsarStorage.SetLastPulse to return Return's arguments
func (m *mPulsarStorageMockSetLastPulse) Return(r error) *PulsarStorageMock {
	m.mock.SetLastPulseFunc = func(p *core.Pulse) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of PulsarStorage.SetLastPulse method
func (m *mPulsarStorageMockSetLastPulse) Set(f func(p *core.Pulse) (r error)) *PulsarStorageMock {
	m.mock.SetLastPulseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetLastPulse implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) SetLastPulse(p *core.Pulse) (r error) {
	atomic.AddUint64(&m.SetLastPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetLastPulseCounter, 1)

	if m.SetLastPulseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetLastPulseMock.mockExpectations, PulsarStorageMockSetLastPulseParams{p},
			"PulsarStorage.SetLastPulse got unexpected parameters")

		if m.SetLastPulseFunc == nil {

			m.t.Fatal("No results are set for the PulsarStorageMock.SetLastPulse")

			return
		}
	}

	if m.SetLastPulseFunc == nil {
		m.t.Fatal("Unexpected call to PulsarStorageMock.SetLastPulse")
		return
	}

	return m.SetLastPulseFunc(p)
}

//SetLastPulseMinimockCounter returns a count of PulsarStorageMock.SetLastPulseFunc invocations
func (m *PulsarStorageMock) SetLastPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLastPulseCounter)
}

//SetLastPulseMinimockPreCounter returns the value of PulsarStorageMock.SetLastPulse invocations
func (m *PulsarStorageMock) SetLastPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLastPulsePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulsarStorageMock) ValidateCallCounters() {

	if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.Close")
	}

	if m.GetLastPulseFunc != nil && atomic.LoadUint64(&m.GetLastPulseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.GetLastPulse")
	}

	if m.SavePulseFunc != nil && atomic.LoadUint64(&m.SavePulseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.SavePulse")
	}

	if m.SetLastPulseFunc != nil && atomic.LoadUint64(&m.SetLastPulseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.SetLastPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
//noinspection GoDeprecation
func (m *PulsarStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulsarStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulsarStorageMock) MinimockFinish() {

	if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.Close")
	}

	if m.GetLastPulseFunc != nil && atomic.LoadUint64(&m.GetLastPulseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.GetLastPulse")
	}

	if m.SavePulseFunc != nil && atomic.LoadUint64(&m.SavePulseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.SavePulse")
	}

	if m.SetLastPulseFunc != nil && atomic.LoadUint64(&m.SetLastPulseCounter) == 0 {
		m.t.Fatal("Expected call to PulsarStorageMock.SetLastPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulsarStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulsarStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.CloseFunc == nil || atomic.LoadUint64(&m.CloseCounter) > 0)
		ok = ok && (m.GetLastPulseFunc == nil || atomic.LoadUint64(&m.GetLastPulseCounter) > 0)
		ok = ok && (m.SavePulseFunc == nil || atomic.LoadUint64(&m.SavePulseCounter) > 0)
		ok = ok && (m.SetLastPulseFunc == nil || atomic.LoadUint64(&m.SetLastPulseCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
				m.t.Error("Expected call to PulsarStorageMock.Close")
			}

			if m.GetLastPulseFunc != nil && atomic.LoadUint64(&m.GetLastPulseCounter) == 0 {
				m.t.Error("Expected call to PulsarStorageMock.GetLastPulse")
			}

			if m.SavePulseFunc != nil && atomic.LoadUint64(&m.SavePulseCounter) == 0 {
				m.t.Error("Expected call to PulsarStorageMock.SavePulse")
			}

			if m.SetLastPulseFunc != nil && atomic.LoadUint64(&m.SetLastPulseCounter) == 0 {
				m.t.Error("Expected call to PulsarStorageMock.SetLastPulse")
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
func (m *PulsarStorageMock) AllMocksCalled() bool {

	if m.CloseFunc != nil && atomic.LoadUint64(&m.CloseCounter) == 0 {
		return false
	}

	if m.GetLastPulseFunc != nil && atomic.LoadUint64(&m.GetLastPulseCounter) == 0 {
		return false
	}

	if m.SavePulseFunc != nil && atomic.LoadUint64(&m.SavePulseCounter) == 0 {
		return false
	}

	if m.SetLastPulseFunc != nil && atomic.LoadUint64(&m.SetLastPulseCounter) == 0 {
		return false
	}

	return true
}
