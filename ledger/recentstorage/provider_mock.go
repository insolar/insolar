package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Provider" can be found in github.com/insolar/insolar/ledger/recentstorage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProviderMock implements github.com/insolar/insolar/ledger/recentstorage.Provider
type ProviderMock struct {
	t minimock.Tester

	GetStorageFunc       func(p core.RecordID) (r RecentStorage)
	GetStorageCounter    uint64
	GetStoragePreCounter uint64
	GetStorageMock       mProviderMockGetStorage
}

//NewProviderMock returns a mock for github.com/insolar/insolar/ledger/recentstorage.Provider
func NewProviderMock(t minimock.Tester) *ProviderMock {
	m := &ProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetStorageMock = mProviderMockGetStorage{mock: m}

	return m
}

type mProviderMockGetStorage struct {
	mock             *ProviderMock
	mockExpectations *ProviderMockGetStorageParams
}

//ProviderMockGetStorageParams represents input parameters of the Provider.GetStorage
type ProviderMockGetStorageParams struct {
	p core.RecordID
}

//Expect sets up expected params for the Provider.GetStorage
func (m *mProviderMockGetStorage) Expect(p core.RecordID) *mProviderMockGetStorage {
	m.mockExpectations = &ProviderMockGetStorageParams{p}
	return m
}

//Return sets up a mock for Provider.GetStorage to return Return's arguments
func (m *mProviderMockGetStorage) Return(r RecentStorage) *ProviderMock {
	m.mock.GetStorageFunc = func(p core.RecordID) RecentStorage {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Provider.GetStorage method
func (m *mProviderMockGetStorage) Set(f func(p core.RecordID) (r RecentStorage)) *ProviderMock {
	m.mock.GetStorageFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) GetStorage(p core.RecordID) (r RecentStorage) {
	atomic.AddUint64(&m.GetStoragePreCounter, 1)
	defer atomic.AddUint64(&m.GetStorageCounter, 1)

	if m.GetStorageMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetStorageMock.mockExpectations, ProviderMockGetStorageParams{p},
			"Provider.GetStorage got unexpected parameters")

		if m.GetStorageFunc == nil {

			m.t.Fatal("No results are set for the ProviderMock.GetStorage")

			return
		}
	}

	if m.GetStorageFunc == nil {
		m.t.Fatal("Unexpected call to ProviderMock.GetStorage")
		return
	}

	return m.GetStorageFunc(p)
}

//GetStorageMinimockCounter returns a count of ProviderMock.GetStorageFunc invocations
func (m *ProviderMock) GetStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStorageCounter)
}

//GetStorageMinimockPreCounter returns the value of ProviderMock.GetStorage invocations
func (m *ProviderMock) GetStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStoragePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) ValidateCallCounters() {

	if m.GetStorageFunc != nil && atomic.LoadUint64(&m.GetStorageCounter) == 0 {
		m.t.Fatal("Expected call to ProviderMock.GetStorage")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ProviderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ProviderMock) MinimockFinish() {

	if m.GetStorageFunc != nil && atomic.LoadUint64(&m.GetStorageCounter) == 0 {
		m.t.Fatal("Expected call to ProviderMock.GetStorage")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ProviderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ProviderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetStorageFunc == nil || atomic.LoadUint64(&m.GetStorageCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetStorageFunc != nil && atomic.LoadUint64(&m.GetStorageCounter) == 0 {
				m.t.Error("Expected call to ProviderMock.GetStorage")
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
func (m *ProviderMock) AllMocksCalled() bool {

	if m.GetStorageFunc != nil && atomic.LoadUint64(&m.GetStorageCounter) == 0 {
		return false
	}

	return true
}
