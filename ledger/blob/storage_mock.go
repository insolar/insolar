package blob

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Storage" can be found in github.com/insolar/insolar/ledger/blob
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//StorageMock implements github.com/insolar/insolar/ledger/blob.Storage
type StorageMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.ID) (r Blob, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mStorageMockForID

	SetFunc       func(p context.Context, p1 insolar.ID, p2 Blob) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mStorageMockSet
}

//NewStorageMock returns a mock for github.com/insolar/insolar/ledger/blob.Storage
func NewStorageMock(t minimock.Tester) *StorageMock {
	m := &StorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mStorageMockForID{mock: m}
	m.SetMock = mStorageMockSet{mock: m}

	return m
}

type mStorageMockForID struct {
	mock              *StorageMock
	mainExpectation   *StorageMockForIDExpectation
	expectationSeries []*StorageMockForIDExpectation
}

type StorageMockForIDExpectation struct {
	input  *StorageMockForIDInput
	result *StorageMockForIDResult
}

type StorageMockForIDInput struct {
	p  context.Context
	p1 insolar.ID
}

type StorageMockForIDResult struct {
	r  Blob
	r1 error
}

//Expect specifies that invocation of Storage.ForID is expected from 1 to Infinity times
func (m *mStorageMockForID) Expect(p context.Context, p1 insolar.ID) *mStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &StorageMockForIDInput{p, p1}
	return m
}

//Return specifies results of invocation of Storage.ForID
func (m *mStorageMockForID) Return(r Blob, r1 error) *StorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &StorageMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.ForID is expected once
func (m *mStorageMockForID) ExpectOnce(p context.Context, p1 insolar.ID) *StorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockForIDExpectation{}
	expectation.input = &StorageMockForIDInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockForIDExpectation) Return(r Blob, r1 error) {
	e.result = &StorageMockForIDResult{r, r1}
}

//Set uses given function f as a mock of Storage.ForID method
func (m *mStorageMockForID) Set(f func(p context.Context, p1 insolar.ID) (r Blob, r1 error)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/blob.Storage interface
func (m *StorageMock) ForID(p context.Context, p1 insolar.ID) (r Blob, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.ForID. %v %v", p, p1)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockForIDInput{p, p1}, "Storage.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockForIDInput{p, p1}, "Storage.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.ForID. %v %v", p, p1)
		return
	}

	return m.ForIDFunc(p, p1)
}

//ForIDMinimockCounter returns a count of StorageMock.ForIDFunc invocations
func (m *StorageMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of StorageMock.ForID invocations
func (m *StorageMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *StorageMock) ForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForIDCounter) == uint64(len(m.ForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForIDFunc != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	return true
}

type mStorageMockSet struct {
	mock              *StorageMock
	mainExpectation   *StorageMockSetExpectation
	expectationSeries []*StorageMockSetExpectation
}

type StorageMockSetExpectation struct {
	input  *StorageMockSetInput
	result *StorageMockSetResult
}

type StorageMockSetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 Blob
}

type StorageMockSetResult struct {
	r error
}

//Expect specifies that invocation of Storage.Set is expected from 1 to Infinity times
func (m *mStorageMockSet) Expect(p context.Context, p1 insolar.ID, p2 Blob) *mStorageMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockSetExpectation{}
	}
	m.mainExpectation.input = &StorageMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Storage.Set
func (m *mStorageMockSet) Return(r error) *StorageMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockSetExpectation{}
	}
	m.mainExpectation.result = &StorageMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.Set is expected once
func (m *mStorageMockSet) ExpectOnce(p context.Context, p1 insolar.ID, p2 Blob) *StorageMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockSetExpectation{}
	expectation.input = &StorageMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockSetExpectation) Return(r error) {
	e.result = &StorageMockSetResult{r}
}

//Set uses given function f as a mock of Storage.Set method
func (m *mStorageMockSet) Set(f func(p context.Context, p1 insolar.ID, p2 Blob) (r error)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/blob.Storage interface
func (m *StorageMock) Set(p context.Context, p1 insolar.ID, p2 Blob) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockSetInput{p, p1, p2}, "Storage.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockSetInput{p, p1, p2}, "Storage.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.Set. %v %v %v", p, p1, p2)
		return
	}

	return m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of StorageMock.SetFunc invocations
func (m *StorageMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of StorageMock.Set invocations
func (m *StorageMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *StorageMock) SetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCounter) == uint64(len(m.SetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFunc != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StorageMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to StorageMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to StorageMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StorageMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to StorageMock.ForID")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to StorageMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to StorageMock.ForID")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to StorageMock.Set")
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
func (m *StorageMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}
