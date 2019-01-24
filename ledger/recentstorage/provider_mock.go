package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Provider" can be found in github.com/insolar/insolar/ledger/recentstorage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProviderMock implements github.com/insolar/insolar/ledger/recentstorage.Provider
type ProviderMock struct {
	t minimock.Tester

	CloneStorageFunc       func(p context.Context, p1 core.RecordID, p2 core.RecordID)
	CloneStorageCounter    uint64
	CloneStoragePreCounter uint64
	CloneStorageMock       mProviderMockCloneStorage

	GetStorageFunc       func(p context.Context, p1 core.RecordID) (r RecentStorage)
	GetStorageCounter    uint64
	GetStoragePreCounter uint64
	GetStorageMock       mProviderMockGetStorage

	RemoveStorageFunc       func(p context.Context, p1 core.RecordID)
	RemoveStorageCounter    uint64
	RemoveStoragePreCounter uint64
	RemoveStorageMock       mProviderMockRemoveStorage
}

//NewProviderMock returns a mock for github.com/insolar/insolar/ledger/recentstorage.Provider
func NewProviderMock(t minimock.Tester) *ProviderMock {
	m := &ProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloneStorageMock = mProviderMockCloneStorage{mock: m}
	m.GetStorageMock = mProviderMockGetStorage{mock: m}
	m.RemoveStorageMock = mProviderMockRemoveStorage{mock: m}

	return m
}

type mProviderMockCloneStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockCloneStorageExpectation
	expectationSeries []*ProviderMockCloneStorageExpectation
}

type ProviderMockCloneStorageExpectation struct {
	input *ProviderMockCloneStorageInput
}

type ProviderMockCloneStorageInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.RecordID
}

//Expect specifies that invocation of Provider.CloneStorage is expected from 1 to Infinity times
func (m *mProviderMockCloneStorage) Expect(p context.Context, p1 core.RecordID, p2 core.RecordID) *mProviderMockCloneStorage {
	m.mock.CloneStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCloneStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockCloneStorageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Provider.CloneStorage
func (m *mProviderMockCloneStorage) Return() *ProviderMock {
	m.mock.CloneStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCloneStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.CloneStorage is expected once
func (m *mProviderMockCloneStorage) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.RecordID) *ProviderMockCloneStorageExpectation {
	m.mock.CloneStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockCloneStorageExpectation{}
	expectation.input = &ProviderMockCloneStorageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.CloneStorage method
func (m *mProviderMockCloneStorage) Set(f func(p context.Context, p1 core.RecordID, p2 core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloneStorageFunc = f
	return m.mock
}

//CloneStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) CloneStorage(p context.Context, p1 core.RecordID, p2 core.RecordID) {
	counter := atomic.AddUint64(&m.CloneStoragePreCounter, 1)
	defer atomic.AddUint64(&m.CloneStorageCounter, 1)

	if len(m.CloneStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloneStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.CloneStorage. %v %v %v", p, p1, p2)
			return
		}

		input := m.CloneStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockCloneStorageInput{p, p1, p2}, "Provider.CloneStorage got unexpected parameters")

		return
	}

	if m.CloneStorageMock.mainExpectation != nil {

		input := m.CloneStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockCloneStorageInput{p, p1, p2}, "Provider.CloneStorage got unexpected parameters")
		}

		return
	}

	if m.CloneStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.CloneStorage. %v %v %v", p, p1, p2)
		return
	}

	m.CloneStorageFunc(p, p1, p2)
}

//CloneStorageMinimockCounter returns a count of ProviderMock.CloneStorageFunc invocations
func (m *ProviderMock) CloneStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloneStorageCounter)
}

//CloneStorageMinimockPreCounter returns the value of ProviderMock.CloneStorage invocations
func (m *ProviderMock) CloneStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CloneStoragePreCounter)
}

//CloneStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) CloneStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloneStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloneStorageCounter) == uint64(len(m.CloneStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloneStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloneStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloneStorageFunc != nil {
		return atomic.LoadUint64(&m.CloneStorageCounter) > 0
	}

	return true
}

type mProviderMockGetStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockGetStorageExpectation
	expectationSeries []*ProviderMockGetStorageExpectation
}

type ProviderMockGetStorageExpectation struct {
	input  *ProviderMockGetStorageInput
	result *ProviderMockGetStorageResult
}

type ProviderMockGetStorageInput struct {
	p  context.Context
	p1 core.RecordID
}

type ProviderMockGetStorageResult struct {
	r RecentStorage
}

//Expect specifies that invocation of Provider.GetStorage is expected from 1 to Infinity times
func (m *mProviderMockGetStorage) Expect(p context.Context, p1 core.RecordID) *mProviderMockGetStorage {
	m.mock.GetStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockGetStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.GetStorage
func (m *mProviderMockGetStorage) Return(r RecentStorage) *ProviderMock {
	m.mock.GetStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetStorageExpectation{}
	}
	m.mainExpectation.result = &ProviderMockGetStorageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.GetStorage is expected once
func (m *mProviderMockGetStorage) ExpectOnce(p context.Context, p1 core.RecordID) *ProviderMockGetStorageExpectation {
	m.mock.GetStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockGetStorageExpectation{}
	expectation.input = &ProviderMockGetStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockGetStorageExpectation) Return(r RecentStorage) {
	e.result = &ProviderMockGetStorageResult{r}
}

//Set uses given function f as a mock of Provider.GetStorage method
func (m *mProviderMockGetStorage) Set(f func(p context.Context, p1 core.RecordID) (r RecentStorage)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStorageFunc = f
	return m.mock
}

//GetStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) GetStorage(p context.Context, p1 core.RecordID) (r RecentStorage) {
	counter := atomic.AddUint64(&m.GetStoragePreCounter, 1)
	defer atomic.AddUint64(&m.GetStorageCounter, 1)

	if len(m.GetStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.GetStorage. %v %v", p, p1)
			return
		}

		input := m.GetStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockGetStorageInput{p, p1}, "Provider.GetStorage got unexpected parameters")

		result := m.GetStorageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetStorage")
			return
		}

		r = result.r

		return
	}

	if m.GetStorageMock.mainExpectation != nil {

		input := m.GetStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockGetStorageInput{p, p1}, "Provider.GetStorage got unexpected parameters")
		}

		result := m.GetStorageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetStorage")
		}

		r = result.r

		return
	}

	if m.GetStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.GetStorage. %v %v", p, p1)
		return
	}

	return m.GetStorageFunc(p, p1)
}

//GetStorageMinimockCounter returns a count of ProviderMock.GetStorageFunc invocations
func (m *ProviderMock) GetStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStorageCounter)
}

//GetStorageMinimockPreCounter returns the value of ProviderMock.GetStorage invocations
func (m *ProviderMock) GetStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStoragePreCounter)
}

//GetStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) GetStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStorageCounter) == uint64(len(m.GetStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStorageFunc != nil {
		return atomic.LoadUint64(&m.GetStorageCounter) > 0
	}

	return true
}

type mProviderMockRemoveStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockRemoveStorageExpectation
	expectationSeries []*ProviderMockRemoveStorageExpectation
}

type ProviderMockRemoveStorageExpectation struct {
	input *ProviderMockRemoveStorageInput
}

type ProviderMockRemoveStorageInput struct {
	p  context.Context
	p1 core.RecordID
}

//Expect specifies that invocation of Provider.RemoveStorage is expected from 1 to Infinity times
func (m *mProviderMockRemoveStorage) Expect(p context.Context, p1 core.RecordID) *mProviderMockRemoveStorage {
	m.mock.RemoveStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemoveStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockRemoveStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.RemoveStorage
func (m *mProviderMockRemoveStorage) Return() *ProviderMock {
	m.mock.RemoveStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemoveStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.RemoveStorage is expected once
func (m *mProviderMockRemoveStorage) ExpectOnce(p context.Context, p1 core.RecordID) *ProviderMockRemoveStorageExpectation {
	m.mock.RemoveStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockRemoveStorageExpectation{}
	expectation.input = &ProviderMockRemoveStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.RemoveStorage method
func (m *mProviderMockRemoveStorage) Set(f func(p context.Context, p1 core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveStorageFunc = f
	return m.mock
}

//RemoveStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) RemoveStorage(p context.Context, p1 core.RecordID) {
	counter := atomic.AddUint64(&m.RemoveStoragePreCounter, 1)
	defer atomic.AddUint64(&m.RemoveStorageCounter, 1)

	if len(m.RemoveStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.RemoveStorage. %v %v", p, p1)
			return
		}

		input := m.RemoveStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockRemoveStorageInput{p, p1}, "Provider.RemoveStorage got unexpected parameters")

		return
	}

	if m.RemoveStorageMock.mainExpectation != nil {

		input := m.RemoveStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockRemoveStorageInput{p, p1}, "Provider.RemoveStorage got unexpected parameters")
		}

		return
	}

	if m.RemoveStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.RemoveStorage. %v %v", p, p1)
		return
	}

	m.RemoveStorageFunc(p, p1)
}

//RemoveStorageMinimockCounter returns a count of ProviderMock.RemoveStorageFunc invocations
func (m *ProviderMock) RemoveStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveStorageCounter)
}

//RemoveStorageMinimockPreCounter returns the value of ProviderMock.RemoveStorage invocations
func (m *ProviderMock) RemoveStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveStoragePreCounter)
}

//RemoveStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) RemoveStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveStorageCounter) == uint64(len(m.RemoveStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveStorageFunc != nil {
		return atomic.LoadUint64(&m.RemoveStorageCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) ValidateCallCounters() {

	if !m.CloneStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.CloneStorage")
	}

	if !m.GetStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetStorage")
	}

	if !m.RemoveStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemoveStorage")
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

	if !m.CloneStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.CloneStorage")
	}

	if !m.GetStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetStorage")
	}

	if !m.RemoveStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemoveStorage")
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
		ok = ok && m.CloneStorageFinished()
		ok = ok && m.GetStorageFinished()
		ok = ok && m.RemoveStorageFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloneStorageFinished() {
				m.t.Error("Expected call to ProviderMock.CloneStorage")
			}

			if !m.GetStorageFinished() {
				m.t.Error("Expected call to ProviderMock.GetStorage")
			}

			if !m.RemoveStorageFinished() {
				m.t.Error("Expected call to ProviderMock.RemoveStorage")
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

	if !m.CloneStorageFinished() {
		return false
	}

	if !m.GetStorageFinished() {
		return false
	}

	if !m.RemoveStorageFinished() {
		return false
	}

	return true
}
