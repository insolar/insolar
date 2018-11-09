package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LocalStorage" can be found in github.com/insolar/insolar/core
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//LocalStorageMock implements github.com/insolar/insolar/core.LocalStorage
type LocalStorageMock struct {
	t minimock.Tester

	GetFunc       func(p context.Context, p1 core.PulseNumber, p2 []byte) (r []byte, r1 error)
	GetCounter    uint64
	GetPreCounter uint64
	GetMock       mLocalStorageMockGet

	IterateFunc       func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) (r error)
	IterateCounter    uint64
	IteratePreCounter uint64
	IterateMock       mLocalStorageMockIterate

	SetFunc       func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mLocalStorageMockSet
}

//NewLocalStorageMock returns a mock for github.com/insolar/insolar/core.LocalStorage
func NewLocalStorageMock(t minimock.Tester) *LocalStorageMock {
	m := &LocalStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetMock = mLocalStorageMockGet{mock: m}
	m.IterateMock = mLocalStorageMockIterate{mock: m}
	m.SetMock = mLocalStorageMockSet{mock: m}

	return m
}

type mLocalStorageMockGet struct {
	mock             *LocalStorageMock
	mockExpectations *LocalStorageMockGetParams
}

//LocalStorageMockGetParams represents input parameters of the LocalStorage.Get
type LocalStorageMockGetParams struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []byte
}

//Expect sets up expected params for the LocalStorage.Get
func (m *mLocalStorageMockGet) Expect(p context.Context, p1 core.PulseNumber, p2 []byte) *mLocalStorageMockGet {
	m.mockExpectations = &LocalStorageMockGetParams{p, p1, p2}
	return m
}

//Return sets up a mock for LocalStorage.Get to return Return's arguments
func (m *mLocalStorageMockGet) Return(r []byte, r1 error) *LocalStorageMock {
	m.mock.GetFunc = func(p context.Context, p1 core.PulseNumber, p2 []byte) ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of LocalStorage.Get method
func (m *mLocalStorageMockGet) Set(f func(p context.Context, p1 core.PulseNumber, p2 []byte) (r []byte, r1 error)) *LocalStorageMock {
	m.mock.GetFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Get implements github.com/insolar/insolar/core.LocalStorage interface
func (m *LocalStorageMock) Get(p context.Context, p1 core.PulseNumber, p2 []byte) (r []byte, r1 error) {
	atomic.AddUint64(&m.GetPreCounter, 1)
	defer atomic.AddUint64(&m.GetCounter, 1)

	if m.GetMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetMock.mockExpectations, LocalStorageMockGetParams{p, p1, p2},
			"LocalStorage.Get got unexpected parameters")

		if m.GetFunc == nil {

			m.t.Fatal("No results are set for the LocalStorageMock.Get")

			return
		}
	}

	if m.GetFunc == nil {
		m.t.Fatal("Unexpected call to LocalStorageMock.Get")
		return
	}

	return m.GetFunc(p, p1, p2)
}

//GetMinimockCounter returns a count of LocalStorageMock.GetFunc invocations
func (m *LocalStorageMock) GetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCounter)
}

//GetMinimockPreCounter returns the value of LocalStorageMock.Get invocations
func (m *LocalStorageMock) GetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreCounter)
}

type mLocalStorageMockIterate struct {
	mock             *LocalStorageMock
	mockExpectations *LocalStorageMockIterateParams
}

//LocalStorageMockIterateParams represents input parameters of the LocalStorage.Iterate
type LocalStorageMockIterateParams struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []byte
	p3 func(p []byte, p1 []byte) (r error)
}

//Expect sets up expected params for the LocalStorage.Iterate
func (m *mLocalStorageMockIterate) Expect(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) *mLocalStorageMockIterate {
	m.mockExpectations = &LocalStorageMockIterateParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for LocalStorage.Iterate to return Return's arguments
func (m *mLocalStorageMockIterate) Return(r error) *LocalStorageMock {
	m.mock.IterateFunc = func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of LocalStorage.Iterate method
func (m *mLocalStorageMockIterate) Set(f func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) (r error)) *LocalStorageMock {
	m.mock.IterateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Iterate implements github.com/insolar/insolar/core.LocalStorage interface
func (m *LocalStorageMock) Iterate(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) (r error) {
	atomic.AddUint64(&m.IteratePreCounter, 1)
	defer atomic.AddUint64(&m.IterateCounter, 1)

	if m.IterateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IterateMock.mockExpectations, LocalStorageMockIterateParams{p, p1, p2, p3},
			"LocalStorage.Iterate got unexpected parameters")

		if m.IterateFunc == nil {

			m.t.Fatal("No results are set for the LocalStorageMock.Iterate")

			return
		}
	}

	if m.IterateFunc == nil {
		m.t.Fatal("Unexpected call to LocalStorageMock.Iterate")
		return
	}

	return m.IterateFunc(p, p1, p2, p3)
}

//IterateMinimockCounter returns a count of LocalStorageMock.IterateFunc invocations
func (m *LocalStorageMock) IterateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IterateCounter)
}

//IterateMinimockPreCounter returns the value of LocalStorageMock.Iterate invocations
func (m *LocalStorageMock) IterateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IteratePreCounter)
}

type mLocalStorageMockSet struct {
	mock             *LocalStorageMock
	mockExpectations *LocalStorageMockSetParams
}

//LocalStorageMockSetParams represents input parameters of the LocalStorage.Set
type LocalStorageMockSetParams struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []byte
	p3 []byte
}

//Expect sets up expected params for the LocalStorage.Set
func (m *mLocalStorageMockSet) Expect(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) *mLocalStorageMockSet {
	m.mockExpectations = &LocalStorageMockSetParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for LocalStorage.Set to return Return's arguments
func (m *mLocalStorageMockSet) Return(r error) *LocalStorageMock {
	m.mock.SetFunc = func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of LocalStorage.Set method
func (m *mLocalStorageMockSet) Set(f func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) (r error)) *LocalStorageMock {
	m.mock.SetFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Set implements github.com/insolar/insolar/core.LocalStorage interface
func (m *LocalStorageMock) Set(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) (r error) {
	atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if m.SetMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetMock.mockExpectations, LocalStorageMockSetParams{p, p1, p2, p3},
			"LocalStorage.Set got unexpected parameters")

		if m.SetFunc == nil {

			m.t.Fatal("No results are set for the LocalStorageMock.Set")

			return
		}
	}

	if m.SetFunc == nil {
		m.t.Fatal("Unexpected call to LocalStorageMock.Set")
		return
	}

	return m.SetFunc(p, p1, p2, p3)
}

//SetMinimockCounter returns a count of LocalStorageMock.SetFunc invocations
func (m *LocalStorageMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of LocalStorageMock.Set invocations
func (m *LocalStorageMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LocalStorageMock) ValidateCallCounters() {

	if m.GetFunc != nil && atomic.LoadUint64(&m.GetCounter) == 0 {
		m.t.Fatal("Expected call to LocalStorageMock.Get")
	}

	if m.IterateFunc != nil && atomic.LoadUint64(&m.IterateCounter) == 0 {
		m.t.Fatal("Expected call to LocalStorageMock.Iterate")
	}

	if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
		m.t.Fatal("Expected call to LocalStorageMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LocalStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LocalStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LocalStorageMock) MinimockFinish() {

	if m.GetFunc != nil && atomic.LoadUint64(&m.GetCounter) == 0 {
		m.t.Fatal("Expected call to LocalStorageMock.Get")
	}

	if m.IterateFunc != nil && atomic.LoadUint64(&m.IterateCounter) == 0 {
		m.t.Fatal("Expected call to LocalStorageMock.Iterate")
	}

	if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
		m.t.Fatal("Expected call to LocalStorageMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LocalStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LocalStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetFunc == nil || atomic.LoadUint64(&m.GetCounter) > 0)
		ok = ok && (m.IterateFunc == nil || atomic.LoadUint64(&m.IterateCounter) > 0)
		ok = ok && (m.SetFunc == nil || atomic.LoadUint64(&m.SetCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetFunc != nil && atomic.LoadUint64(&m.GetCounter) == 0 {
				m.t.Error("Expected call to LocalStorageMock.Get")
			}

			if m.IterateFunc != nil && atomic.LoadUint64(&m.IterateCounter) == 0 {
				m.t.Error("Expected call to LocalStorageMock.Iterate")
			}

			if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
				m.t.Error("Expected call to LocalStorageMock.Set")
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
func (m *LocalStorageMock) AllMocksCalled() bool {

	if m.GetFunc != nil && atomic.LoadUint64(&m.GetCounter) == 0 {
		return false
	}

	if m.IterateFunc != nil && atomic.LoadUint64(&m.IterateCounter) == 0 {
		return false
	}

	if m.SetFunc != nil && atomic.LoadUint64(&m.SetCounter) == 0 {
		return false
	}

	return true
}
