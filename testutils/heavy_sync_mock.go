package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "HeavySync" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	testify_assert "github.com/stretchr/testify/assert"
)

//HeavySyncMock implements github.com/insolar/insolar/core.HeavySync
type HeavySyncMock struct {
	t minimock.Tester

	ResetFunc       func(p context.Context, p1 core.PulseNumber) (r error)
	ResetCounter    uint64
	ResetPreCounter uint64
	ResetMock       mHeavySyncMockReset

	StartFunc       func(p context.Context, p1 core.PulseNumber) (r error)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mHeavySyncMockStart

	StopFunc       func(p context.Context, p1 core.PulseNumber) (r error)
	StopCounter    uint64
	StopPreCounter uint64
	StopMock       mHeavySyncMockStop

	StoreFunc       func(p context.Context, p1 core.PulseNumber, p2 []core.KV) (r error)
	StoreCounter    uint64
	StorePreCounter uint64
	StoreMock       mHeavySyncMockStore
}

//NewHeavySyncMock returns a mock for github.com/insolar/insolar/core.HeavySync
func NewHeavySyncMock(t minimock.Tester) *HeavySyncMock {
	m := &HeavySyncMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ResetMock = mHeavySyncMockReset{mock: m}
	m.StartMock = mHeavySyncMockStart{mock: m}
	m.StopMock = mHeavySyncMockStop{mock: m}
	m.StoreMock = mHeavySyncMockStore{mock: m}

	return m
}

type mHeavySyncMockReset struct {
	mock             *HeavySyncMock
	mockExpectations *HeavySyncMockResetParams
}

//HeavySyncMockResetParams represents input parameters of the HeavySync.Reset
type HeavySyncMockResetParams struct {
	p  context.Context
	p1 core.PulseNumber
}

//Expect sets up expected params for the HeavySync.Reset
func (m *mHeavySyncMockReset) Expect(p context.Context, p1 core.PulseNumber) *mHeavySyncMockReset {
	m.mockExpectations = &HeavySyncMockResetParams{p, p1}
	return m
}

//Return sets up a mock for HeavySync.Reset to return Return's arguments
func (m *mHeavySyncMockReset) Return(r error) *HeavySyncMock {
	m.mock.ResetFunc = func(p context.Context, p1 core.PulseNumber) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HeavySync.Reset method
func (m *mHeavySyncMockReset) Set(f func(p context.Context, p1 core.PulseNumber) (r error)) *HeavySyncMock {
	m.mock.ResetFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Reset implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Reset(p context.Context, p1 core.PulseNumber) (r error) {
	atomic.AddUint64(&m.ResetPreCounter, 1)
	defer atomic.AddUint64(&m.ResetCounter, 1)

	if m.ResetMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ResetMock.mockExpectations, HeavySyncMockResetParams{p, p1},
			"HeavySync.Reset got unexpected parameters")

		if m.ResetFunc == nil {

			m.t.Fatal("No results are set for the HeavySyncMock.Reset")

			return
		}
	}

	if m.ResetFunc == nil {
		m.t.Fatal("Unexpected call to HeavySyncMock.Reset")
		return
	}

	return m.ResetFunc(p, p1)
}

//ResetMinimockCounter returns a count of HeavySyncMock.ResetFunc invocations
func (m *HeavySyncMock) ResetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResetCounter)
}

//ResetMinimockPreCounter returns the value of HeavySyncMock.Reset invocations
func (m *HeavySyncMock) ResetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResetPreCounter)
}

type mHeavySyncMockStart struct {
	mock             *HeavySyncMock
	mockExpectations *HeavySyncMockStartParams
}

//HeavySyncMockStartParams represents input parameters of the HeavySync.Start
type HeavySyncMockStartParams struct {
	p  context.Context
	p1 core.PulseNumber
}

//Expect sets up expected params for the HeavySync.Start
func (m *mHeavySyncMockStart) Expect(p context.Context, p1 core.PulseNumber) *mHeavySyncMockStart {
	m.mockExpectations = &HeavySyncMockStartParams{p, p1}
	return m
}

//Return sets up a mock for HeavySync.Start to return Return's arguments
func (m *mHeavySyncMockStart) Return(r error) *HeavySyncMock {
	m.mock.StartFunc = func(p context.Context, p1 core.PulseNumber) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HeavySync.Start method
func (m *mHeavySyncMockStart) Set(f func(p context.Context, p1 core.PulseNumber) (r error)) *HeavySyncMock {
	m.mock.StartFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Start implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Start(p context.Context, p1 core.PulseNumber) (r error) {
	atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if m.StartMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.StartMock.mockExpectations, HeavySyncMockStartParams{p, p1},
			"HeavySync.Start got unexpected parameters")

		if m.StartFunc == nil {

			m.t.Fatal("No results are set for the HeavySyncMock.Start")

			return
		}
	}

	if m.StartFunc == nil {
		m.t.Fatal("Unexpected call to HeavySyncMock.Start")
		return
	}

	return m.StartFunc(p, p1)
}

//StartMinimockCounter returns a count of HeavySyncMock.StartFunc invocations
func (m *HeavySyncMock) StartMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StartCounter)
}

//StartMinimockPreCounter returns the value of HeavySyncMock.Start invocations
func (m *HeavySyncMock) StartMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StartPreCounter)
}

type mHeavySyncMockStop struct {
	mock             *HeavySyncMock
	mockExpectations *HeavySyncMockStopParams
}

//HeavySyncMockStopParams represents input parameters of the HeavySync.Stop
type HeavySyncMockStopParams struct {
	p  context.Context
	p1 core.PulseNumber
}

//Expect sets up expected params for the HeavySync.Stop
func (m *mHeavySyncMockStop) Expect(p context.Context, p1 core.PulseNumber) *mHeavySyncMockStop {
	m.mockExpectations = &HeavySyncMockStopParams{p, p1}
	return m
}

//Return sets up a mock for HeavySync.Stop to return Return's arguments
func (m *mHeavySyncMockStop) Return(r error) *HeavySyncMock {
	m.mock.StopFunc = func(p context.Context, p1 core.PulseNumber) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HeavySync.Stop method
func (m *mHeavySyncMockStop) Set(f func(p context.Context, p1 core.PulseNumber) (r error)) *HeavySyncMock {
	m.mock.StopFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Stop implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Stop(p context.Context, p1 core.PulseNumber) (r error) {
	atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if m.StopMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.StopMock.mockExpectations, HeavySyncMockStopParams{p, p1},
			"HeavySync.Stop got unexpected parameters")

		if m.StopFunc == nil {

			m.t.Fatal("No results are set for the HeavySyncMock.Stop")

			return
		}
	}

	if m.StopFunc == nil {
		m.t.Fatal("Unexpected call to HeavySyncMock.Stop")
		return
	}

	return m.StopFunc(p, p1)
}

//StopMinimockCounter returns a count of HeavySyncMock.StopFunc invocations
func (m *HeavySyncMock) StopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StopCounter)
}

//StopMinimockPreCounter returns the value of HeavySyncMock.Stop invocations
func (m *HeavySyncMock) StopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StopPreCounter)
}

type mHeavySyncMockStore struct {
	mock             *HeavySyncMock
	mockExpectations *HeavySyncMockStoreParams
}

//HeavySyncMockStoreParams represents input parameters of the HeavySync.Store
type HeavySyncMockStoreParams struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []core.KV
}

//Expect sets up expected params for the HeavySync.Store
func (m *mHeavySyncMockStore) Expect(p context.Context, p1 core.PulseNumber, p2 []core.KV) *mHeavySyncMockStore {
	m.mockExpectations = &HeavySyncMockStoreParams{p, p1, p2}
	return m
}

//Return sets up a mock for HeavySync.Store to return Return's arguments
func (m *mHeavySyncMockStore) Return(r error) *HeavySyncMock {
	m.mock.StoreFunc = func(p context.Context, p1 core.PulseNumber, p2 []core.KV) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of HeavySync.Store method
func (m *mHeavySyncMockStore) Set(f func(p context.Context, p1 core.PulseNumber, p2 []core.KV) (r error)) *HeavySyncMock {
	m.mock.StoreFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Store implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Store(p context.Context, p1 core.PulseNumber, p2 []core.KV) (r error) {
	atomic.AddUint64(&m.StorePreCounter, 1)
	defer atomic.AddUint64(&m.StoreCounter, 1)

	if m.StoreMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.StoreMock.mockExpectations, HeavySyncMockStoreParams{p, p1, p2},
			"HeavySync.Store got unexpected parameters")

		if m.StoreFunc == nil {

			m.t.Fatal("No results are set for the HeavySyncMock.Store")

			return
		}
	}

	if m.StoreFunc == nil {
		m.t.Fatal("Unexpected call to HeavySyncMock.Store")
		return
	}

	return m.StoreFunc(p, p1, p2)
}

//StoreMinimockCounter returns a count of HeavySyncMock.StoreFunc invocations
func (m *HeavySyncMock) StoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreCounter)
}

//StoreMinimockPreCounter returns the value of HeavySyncMock.Store invocations
func (m *HeavySyncMock) StoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StorePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HeavySyncMock) ValidateCallCounters() {

	if m.ResetFunc != nil && atomic.LoadUint64(&m.ResetCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Reset")
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Start")
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Stop")
	}

	if m.StoreFunc != nil && atomic.LoadUint64(&m.StoreCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Store")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HeavySyncMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *HeavySyncMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *HeavySyncMock) MinimockFinish() {

	if m.ResetFunc != nil && atomic.LoadUint64(&m.ResetCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Reset")
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Start")
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Stop")
	}

	if m.StoreFunc != nil && atomic.LoadUint64(&m.StoreCounter) == 0 {
		m.t.Fatal("Expected call to HeavySyncMock.Store")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *HeavySyncMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *HeavySyncMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.ResetFunc == nil || atomic.LoadUint64(&m.ResetCounter) > 0)
		ok = ok && (m.StartFunc == nil || atomic.LoadUint64(&m.StartCounter) > 0)
		ok = ok && (m.StopFunc == nil || atomic.LoadUint64(&m.StopCounter) > 0)
		ok = ok && (m.StoreFunc == nil || atomic.LoadUint64(&m.StoreCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ResetFunc != nil && atomic.LoadUint64(&m.ResetCounter) == 0 {
				m.t.Error("Expected call to HeavySyncMock.Reset")
			}

			if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
				m.t.Error("Expected call to HeavySyncMock.Start")
			}

			if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
				m.t.Error("Expected call to HeavySyncMock.Stop")
			}

			if m.StoreFunc != nil && atomic.LoadUint64(&m.StoreCounter) == 0 {
				m.t.Error("Expected call to HeavySyncMock.Store")
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
func (m *HeavySyncMock) AllMocksCalled() bool {

	if m.ResetFunc != nil && atomic.LoadUint64(&m.ResetCounter) == 0 {
		return false
	}

	if m.StartFunc != nil && atomic.LoadUint64(&m.StartCounter) == 0 {
		return false
	}

	if m.StopFunc != nil && atomic.LoadUint64(&m.StopCounter) == 0 {
		return false
	}

	if m.StoreFunc != nil && atomic.LoadUint64(&m.StoreCounter) == 0 {
		return false
	}

	return true
}
