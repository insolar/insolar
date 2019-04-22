package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Provider" can be found in github.com/insolar/insolar/ledger/light/recentstorage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProviderMock implements github.com/insolar/insolar/ledger/light/recentstorage.Provider
type ProviderMock struct {
	t minimock.Tester

	ClonePendingStorageFunc       func(p context.Context, p1 insolar.ID, p2 insolar.ID)
	ClonePendingStorageCounter    uint64
	ClonePendingStoragePreCounter uint64
	ClonePendingStorageMock       mProviderMockClonePendingStorage

	CountFunc       func() (r int)
	CountCounter    uint64
	CountPreCounter uint64
	CountMock       mProviderMockCount

	GetPendingStorageFunc       func(p context.Context, p1 insolar.ID) (r PendingStorage)
	GetPendingStorageCounter    uint64
	GetPendingStoragePreCounter uint64
	GetPendingStorageMock       mProviderMockGetPendingStorage

	RemovePendingStorageFunc       func(p context.Context, p1 insolar.ID)
	RemovePendingStorageCounter    uint64
	RemovePendingStoragePreCounter uint64
	RemovePendingStorageMock       mProviderMockRemovePendingStorage
}

//NewProviderMock returns a mock for github.com/insolar/insolar/ledger/light/recentstorage.Provider
func NewProviderMock(t minimock.Tester) *ProviderMock {
	m := &ProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ClonePendingStorageMock = mProviderMockClonePendingStorage{mock: m}
	m.CountMock = mProviderMockCount{mock: m}
	m.GetPendingStorageMock = mProviderMockGetPendingStorage{mock: m}
	m.RemovePendingStorageMock = mProviderMockRemovePendingStorage{mock: m}

	return m
}

type mProviderMockClonePendingStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockClonePendingStorageExpectation
	expectationSeries []*ProviderMockClonePendingStorageExpectation
}

type ProviderMockClonePendingStorageExpectation struct {
	input *ProviderMockClonePendingStorageInput
}

type ProviderMockClonePendingStorageInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.ID
}

//Expect specifies that invocation of Provider.ClonePendingStorage is expected from 1 to Infinity times
func (m *mProviderMockClonePendingStorage) Expect(p context.Context, p1 insolar.ID, p2 insolar.ID) *mProviderMockClonePendingStorage {
	m.mock.ClonePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockClonePendingStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockClonePendingStorageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Provider.ClonePendingStorage
func (m *mProviderMockClonePendingStorage) Return() *ProviderMock {
	m.mock.ClonePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockClonePendingStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.ClonePendingStorage is expected once
func (m *mProviderMockClonePendingStorage) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.ID) *ProviderMockClonePendingStorageExpectation {
	m.mock.ClonePendingStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockClonePendingStorageExpectation{}
	expectation.input = &ProviderMockClonePendingStorageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.ClonePendingStorage method
func (m *mProviderMockClonePendingStorage) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.ID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ClonePendingStorageFunc = f
	return m.mock
}

//ClonePendingStorage implements github.com/insolar/insolar/ledger/light/recentstorage.Provider interface
func (m *ProviderMock) ClonePendingStorage(p context.Context, p1 insolar.ID, p2 insolar.ID) {
	counter := atomic.AddUint64(&m.ClonePendingStoragePreCounter, 1)
	defer atomic.AddUint64(&m.ClonePendingStorageCounter, 1)

	if len(m.ClonePendingStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ClonePendingStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.ClonePendingStorage. %v %v %v", p, p1, p2)
			return
		}

		input := m.ClonePendingStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockClonePendingStorageInput{p, p1, p2}, "Provider.ClonePendingStorage got unexpected parameters")

		return
	}

	if m.ClonePendingStorageMock.mainExpectation != nil {

		input := m.ClonePendingStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockClonePendingStorageInput{p, p1, p2}, "Provider.ClonePendingStorage got unexpected parameters")
		}

		return
	}

	if m.ClonePendingStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.ClonePendingStorage. %v %v %v", p, p1, p2)
		return
	}

	m.ClonePendingStorageFunc(p, p1, p2)
}

//ClonePendingStorageMinimockCounter returns a count of ProviderMock.ClonePendingStorageFunc invocations
func (m *ProviderMock) ClonePendingStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ClonePendingStorageCounter)
}

//ClonePendingStorageMinimockPreCounter returns the value of ProviderMock.ClonePendingStorage invocations
func (m *ProviderMock) ClonePendingStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClonePendingStoragePreCounter)
}

//ClonePendingStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) ClonePendingStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ClonePendingStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ClonePendingStorageCounter) == uint64(len(m.ClonePendingStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ClonePendingStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ClonePendingStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ClonePendingStorageFunc != nil {
		return atomic.LoadUint64(&m.ClonePendingStorageCounter) > 0
	}

	return true
}

type mProviderMockCount struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockCountExpectation
	expectationSeries []*ProviderMockCountExpectation
}

type ProviderMockCountExpectation struct {
	result *ProviderMockCountResult
}

type ProviderMockCountResult struct {
	r int
}

//Expect specifies that invocation of Provider.Count is expected from 1 to Infinity times
func (m *mProviderMockCount) Expect() *mProviderMockCount {
	m.mock.CountFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCountExpectation{}
	}

	return m
}

//Return specifies results of invocation of Provider.Count
func (m *mProviderMockCount) Return(r int) *ProviderMock {
	m.mock.CountFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCountExpectation{}
	}
	m.mainExpectation.result = &ProviderMockCountResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.Count is expected once
func (m *mProviderMockCount) ExpectOnce() *ProviderMockCountExpectation {
	m.mock.CountFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockCountExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockCountExpectation) Return(r int) {
	e.result = &ProviderMockCountResult{r}
}

//Set uses given function f as a mock of Provider.Count method
func (m *mProviderMockCount) Set(f func() (r int)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CountFunc = f
	return m.mock
}

//Count implements github.com/insolar/insolar/ledger/light/recentstorage.Provider interface
func (m *ProviderMock) Count() (r int) {
	counter := atomic.AddUint64(&m.CountPreCounter, 1)
	defer atomic.AddUint64(&m.CountCounter, 1)

	if len(m.CountMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CountMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.Count.")
			return
		}

		result := m.CountMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.Count")
			return
		}

		r = result.r

		return
	}

	if m.CountMock.mainExpectation != nil {

		result := m.CountMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.Count")
		}

		r = result.r

		return
	}

	if m.CountFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.Count.")
		return
	}

	return m.CountFunc()
}

//CountMinimockCounter returns a count of ProviderMock.CountFunc invocations
func (m *ProviderMock) CountMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CountCounter)
}

//CountMinimockPreCounter returns the value of ProviderMock.Count invocations
func (m *ProviderMock) CountMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CountPreCounter)
}

//CountFinished returns true if mock invocations count is ok
func (m *ProviderMock) CountFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CountMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CountCounter) == uint64(len(m.CountMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CountMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CountCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CountFunc != nil {
		return atomic.LoadUint64(&m.CountCounter) > 0
	}

	return true
}

type mProviderMockGetPendingStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockGetPendingStorageExpectation
	expectationSeries []*ProviderMockGetPendingStorageExpectation
}

type ProviderMockGetPendingStorageExpectation struct {
	input  *ProviderMockGetPendingStorageInput
	result *ProviderMockGetPendingStorageResult
}

type ProviderMockGetPendingStorageInput struct {
	p  context.Context
	p1 insolar.ID
}

type ProviderMockGetPendingStorageResult struct {
	r PendingStorage
}

//Expect specifies that invocation of Provider.GetPendingStorage is expected from 1 to Infinity times
func (m *mProviderMockGetPendingStorage) Expect(p context.Context, p1 insolar.ID) *mProviderMockGetPendingStorage {
	m.mock.GetPendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetPendingStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockGetPendingStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.GetPendingStorage
func (m *mProviderMockGetPendingStorage) Return(r PendingStorage) *ProviderMock {
	m.mock.GetPendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetPendingStorageExpectation{}
	}
	m.mainExpectation.result = &ProviderMockGetPendingStorageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.GetPendingStorage is expected once
func (m *mProviderMockGetPendingStorage) ExpectOnce(p context.Context, p1 insolar.ID) *ProviderMockGetPendingStorageExpectation {
	m.mock.GetPendingStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockGetPendingStorageExpectation{}
	expectation.input = &ProviderMockGetPendingStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockGetPendingStorageExpectation) Return(r PendingStorage) {
	e.result = &ProviderMockGetPendingStorageResult{r}
}

//Set uses given function f as a mock of Provider.GetPendingStorage method
func (m *mProviderMockGetPendingStorage) Set(f func(p context.Context, p1 insolar.ID) (r PendingStorage)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPendingStorageFunc = f
	return m.mock
}

//GetPendingStorage implements github.com/insolar/insolar/ledger/light/recentstorage.Provider interface
func (m *ProviderMock) GetPendingStorage(p context.Context, p1 insolar.ID) (r PendingStorage) {
	counter := atomic.AddUint64(&m.GetPendingStoragePreCounter, 1)
	defer atomic.AddUint64(&m.GetPendingStorageCounter, 1)

	if len(m.GetPendingStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPendingStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.GetPendingStorage. %v %v", p, p1)
			return
		}

		input := m.GetPendingStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockGetPendingStorageInput{p, p1}, "Provider.GetPendingStorage got unexpected parameters")

		result := m.GetPendingStorageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetPendingStorage")
			return
		}

		r = result.r

		return
	}

	if m.GetPendingStorageMock.mainExpectation != nil {

		input := m.GetPendingStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockGetPendingStorageInput{p, p1}, "Provider.GetPendingStorage got unexpected parameters")
		}

		result := m.GetPendingStorageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetPendingStorage")
		}

		r = result.r

		return
	}

	if m.GetPendingStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.GetPendingStorage. %v %v", p, p1)
		return
	}

	return m.GetPendingStorageFunc(p, p1)
}

//GetPendingStorageMinimockCounter returns a count of ProviderMock.GetPendingStorageFunc invocations
func (m *ProviderMock) GetPendingStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPendingStorageCounter)
}

//GetPendingStorageMinimockPreCounter returns the value of ProviderMock.GetPendingStorage invocations
func (m *ProviderMock) GetPendingStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPendingStoragePreCounter)
}

//GetPendingStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) GetPendingStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPendingStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPendingStorageCounter) == uint64(len(m.GetPendingStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPendingStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPendingStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPendingStorageFunc != nil {
		return atomic.LoadUint64(&m.GetPendingStorageCounter) > 0
	}

	return true
}

type mProviderMockRemovePendingStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockRemovePendingStorageExpectation
	expectationSeries []*ProviderMockRemovePendingStorageExpectation
}

type ProviderMockRemovePendingStorageExpectation struct {
	input *ProviderMockRemovePendingStorageInput
}

type ProviderMockRemovePendingStorageInput struct {
	p  context.Context
	p1 insolar.ID
}

//Expect specifies that invocation of Provider.RemovePendingStorage is expected from 1 to Infinity times
func (m *mProviderMockRemovePendingStorage) Expect(p context.Context, p1 insolar.ID) *mProviderMockRemovePendingStorage {
	m.mock.RemovePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemovePendingStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockRemovePendingStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.RemovePendingStorage
func (m *mProviderMockRemovePendingStorage) Return() *ProviderMock {
	m.mock.RemovePendingStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockRemovePendingStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.RemovePendingStorage is expected once
func (m *mProviderMockRemovePendingStorage) ExpectOnce(p context.Context, p1 insolar.ID) *ProviderMockRemovePendingStorageExpectation {
	m.mock.RemovePendingStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockRemovePendingStorageExpectation{}
	expectation.input = &ProviderMockRemovePendingStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.RemovePendingStorage method
func (m *mProviderMockRemovePendingStorage) Set(f func(p context.Context, p1 insolar.ID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemovePendingStorageFunc = f
	return m.mock
}

//RemovePendingStorage implements github.com/insolar/insolar/ledger/light/recentstorage.Provider interface
func (m *ProviderMock) RemovePendingStorage(p context.Context, p1 insolar.ID) {
	counter := atomic.AddUint64(&m.RemovePendingStoragePreCounter, 1)
	defer atomic.AddUint64(&m.RemovePendingStorageCounter, 1)

	if len(m.RemovePendingStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemovePendingStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.RemovePendingStorage. %v %v", p, p1)
			return
		}

		input := m.RemovePendingStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockRemovePendingStorageInput{p, p1}, "Provider.RemovePendingStorage got unexpected parameters")

		return
	}

	if m.RemovePendingStorageMock.mainExpectation != nil {

		input := m.RemovePendingStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockRemovePendingStorageInput{p, p1}, "Provider.RemovePendingStorage got unexpected parameters")
		}

		return
	}

	if m.RemovePendingStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.RemovePendingStorage. %v %v", p, p1)
		return
	}

	m.RemovePendingStorageFunc(p, p1)
}

//RemovePendingStorageMinimockCounter returns a count of ProviderMock.RemovePendingStorageFunc invocations
func (m *ProviderMock) RemovePendingStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingStorageCounter)
}

//RemovePendingStorageMinimockPreCounter returns the value of ProviderMock.RemovePendingStorage invocations
func (m *ProviderMock) RemovePendingStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemovePendingStoragePreCounter)
}

//RemovePendingStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) RemovePendingStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemovePendingStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemovePendingStorageCounter) == uint64(len(m.RemovePendingStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemovePendingStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemovePendingStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemovePendingStorageFunc != nil {
		return atomic.LoadUint64(&m.RemovePendingStorageCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) ValidateCallCounters() {

	if !m.ClonePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.ClonePendingStorage")
	}

	if !m.CountFinished() {
		m.t.Fatal("Expected call to ProviderMock.Count")
	}

	if !m.GetPendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetPendingStorage")
	}

	if !m.RemovePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemovePendingStorage")
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

	if !m.ClonePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.ClonePendingStorage")
	}

	if !m.CountFinished() {
		m.t.Fatal("Expected call to ProviderMock.Count")
	}

	if !m.GetPendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetPendingStorage")
	}

	if !m.RemovePendingStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.RemovePendingStorage")
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
		ok = ok && m.ClonePendingStorageFinished()
		ok = ok && m.CountFinished()
		ok = ok && m.GetPendingStorageFinished()
		ok = ok && m.RemovePendingStorageFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ClonePendingStorageFinished() {
				m.t.Error("Expected call to ProviderMock.ClonePendingStorage")
			}

			if !m.CountFinished() {
				m.t.Error("Expected call to ProviderMock.Count")
			}

			if !m.GetPendingStorageFinished() {
				m.t.Error("Expected call to ProviderMock.GetPendingStorage")
			}

			if !m.RemovePendingStorageFinished() {
				m.t.Error("Expected call to ProviderMock.RemovePendingStorage")
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

	if !m.ClonePendingStorageFinished() {
		return false
	}

	if !m.CountFinished() {
		return false
	}

	if !m.GetPendingStorageFinished() {
		return false
	}

	if !m.RemovePendingStorageFinished() {
		return false
	}

	return true
}
