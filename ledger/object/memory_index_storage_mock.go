package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MemoryIndexStorage" can be found in github.com/insolar/insolar/ledger/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//MemoryIndexStorageMock implements github.com/insolar/insolar/ledger/object.MemoryIndexStorage
type MemoryIndexStorageMock struct {
	t minimock.Tester

	ForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r record.Index, r1 error)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mMemoryIndexStorageMockForID

	ForPulseFunc       func(p context.Context, p1 insolar.PulseNumber) (r []record.Index)
	ForPulseCounter    uint64
	ForPulsePreCounter uint64
	ForPulseMock       mMemoryIndexStorageMockForPulse

	SetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 record.Index)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mMemoryIndexStorageMockSet
}

//NewMemoryIndexStorageMock returns a mock for github.com/insolar/insolar/ledger/object.MemoryIndexStorage
func NewMemoryIndexStorageMock(t minimock.Tester) *MemoryIndexStorageMock {
	m := &MemoryIndexStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForIDMock = mMemoryIndexStorageMockForID{mock: m}
	m.ForPulseMock = mMemoryIndexStorageMockForPulse{mock: m}
	m.SetMock = mMemoryIndexStorageMockSet{mock: m}

	return m
}

type mMemoryIndexStorageMockForID struct {
	mock              *MemoryIndexStorageMock
	mainExpectation   *MemoryIndexStorageMockForIDExpectation
	expectationSeries []*MemoryIndexStorageMockForIDExpectation
}

type MemoryIndexStorageMockForIDExpectation struct {
	input  *MemoryIndexStorageMockForIDInput
	result *MemoryIndexStorageMockForIDResult
}

type MemoryIndexStorageMockForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type MemoryIndexStorageMockForIDResult struct {
	r  record.Index
	r1 error
}

//Expect specifies that invocation of MemoryIndexStorage.ForID is expected from 1 to Infinity times
func (m *mMemoryIndexStorageMockForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mMemoryIndexStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexStorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &MemoryIndexStorageMockForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of MemoryIndexStorage.ForID
func (m *mMemoryIndexStorageMockForID) Return(r record.Index, r1 error) *MemoryIndexStorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexStorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &MemoryIndexStorageMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of MemoryIndexStorage.ForID is expected once
func (m *mMemoryIndexStorageMockForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *MemoryIndexStorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &MemoryIndexStorageMockForIDExpectation{}
	expectation.input = &MemoryIndexStorageMockForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemoryIndexStorageMockForIDExpectation) Return(r record.Index, r1 error) {
	e.result = &MemoryIndexStorageMockForIDResult{r, r1}
}

//Set uses given function f as a mock of MemoryIndexStorage.ForID method
func (m *mMemoryIndexStorageMockForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r record.Index, r1 error)) *MemoryIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/ledger/object.MemoryIndexStorage interface
func (m *MemoryIndexStorageMock) ForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r record.Index, r1 error) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemoryIndexStorageMock.ForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MemoryIndexStorageMockForIDInput{p, p1, p2}, "MemoryIndexStorage.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemoryIndexStorageMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MemoryIndexStorageMockForIDInput{p, p1, p2}, "MemoryIndexStorage.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemoryIndexStorageMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to MemoryIndexStorageMock.ForID. %v %v %v", p, p1, p2)
		return
	}

	return m.ForIDFunc(p, p1, p2)
}

//ForIDMinimockCounter returns a count of MemoryIndexStorageMock.ForIDFunc invocations
func (m *MemoryIndexStorageMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of MemoryIndexStorageMock.ForID invocations
func (m *MemoryIndexStorageMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *MemoryIndexStorageMock) ForIDFinished() bool {
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

type mMemoryIndexStorageMockForPulse struct {
	mock              *MemoryIndexStorageMock
	mainExpectation   *MemoryIndexStorageMockForPulseExpectation
	expectationSeries []*MemoryIndexStorageMockForPulseExpectation
}

type MemoryIndexStorageMockForPulseExpectation struct {
	input  *MemoryIndexStorageMockForPulseInput
	result *MemoryIndexStorageMockForPulseResult
}

type MemoryIndexStorageMockForPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type MemoryIndexStorageMockForPulseResult struct {
	r []record.Index
}

//Expect specifies that invocation of MemoryIndexStorage.ForPulse is expected from 1 to Infinity times
func (m *mMemoryIndexStorageMockForPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mMemoryIndexStorageMockForPulse {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexStorageMockForPulseExpectation{}
	}
	m.mainExpectation.input = &MemoryIndexStorageMockForPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of MemoryIndexStorage.ForPulse
func (m *mMemoryIndexStorageMockForPulse) Return(r []record.Index) *MemoryIndexStorageMock {
	m.mock.ForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexStorageMockForPulseExpectation{}
	}
	m.mainExpectation.result = &MemoryIndexStorageMockForPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemoryIndexStorage.ForPulse is expected once
func (m *mMemoryIndexStorageMockForPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *MemoryIndexStorageMockForPulseExpectation {
	m.mock.ForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &MemoryIndexStorageMockForPulseExpectation{}
	expectation.input = &MemoryIndexStorageMockForPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemoryIndexStorageMockForPulseExpectation) Return(r []record.Index) {
	e.result = &MemoryIndexStorageMockForPulseResult{r}
}

//Set uses given function f as a mock of MemoryIndexStorage.ForPulse method
func (m *mMemoryIndexStorageMockForPulse) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []record.Index)) *MemoryIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/object.MemoryIndexStorage interface
func (m *MemoryIndexStorageMock) ForPulse(p context.Context, p1 insolar.PulseNumber) (r []record.Index) {
	counter := atomic.AddUint64(&m.ForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseCounter, 1)

	if len(m.ForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemoryIndexStorageMock.ForPulse. %v %v", p, p1)
			return
		}

		input := m.ForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MemoryIndexStorageMockForPulseInput{p, p1}, "MemoryIndexStorage.ForPulse got unexpected parameters")

		result := m.ForPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemoryIndexStorageMock.ForPulse")
			return
		}

		r = result.r

		return
	}

	if m.ForPulseMock.mainExpectation != nil {

		input := m.ForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MemoryIndexStorageMockForPulseInput{p, p1}, "MemoryIndexStorage.ForPulse got unexpected parameters")
		}

		result := m.ForPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemoryIndexStorageMock.ForPulse")
		}

		r = result.r

		return
	}

	if m.ForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to MemoryIndexStorageMock.ForPulse. %v %v", p, p1)
		return
	}

	return m.ForPulseFunc(p, p1)
}

//ForPulseMinimockCounter returns a count of MemoryIndexStorageMock.ForPulseFunc invocations
func (m *MemoryIndexStorageMock) ForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseCounter)
}

//ForPulseMinimockPreCounter returns the value of MemoryIndexStorageMock.ForPulse invocations
func (m *MemoryIndexStorageMock) ForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulsePreCounter)
}

//ForPulseFinished returns true if mock invocations count is ok
func (m *MemoryIndexStorageMock) ForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseCounter) == uint64(len(m.ForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseFunc != nil {
		return atomic.LoadUint64(&m.ForPulseCounter) > 0
	}

	return true
}

type mMemoryIndexStorageMockSet struct {
	mock              *MemoryIndexStorageMock
	mainExpectation   *MemoryIndexStorageMockSetExpectation
	expectationSeries []*MemoryIndexStorageMockSetExpectation
}

type MemoryIndexStorageMockSetExpectation struct {
	input *MemoryIndexStorageMockSetInput
}

type MemoryIndexStorageMockSetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 record.Index
}

//Expect specifies that invocation of MemoryIndexStorage.Set is expected from 1 to Infinity times
func (m *mMemoryIndexStorageMockSet) Expect(p context.Context, p1 insolar.PulseNumber, p2 record.Index) *mMemoryIndexStorageMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexStorageMockSetExpectation{}
	}
	m.mainExpectation.input = &MemoryIndexStorageMockSetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of MemoryIndexStorage.Set
func (m *mMemoryIndexStorageMockSet) Return() *MemoryIndexStorageMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemoryIndexStorageMockSetExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of MemoryIndexStorage.Set is expected once
func (m *mMemoryIndexStorageMockSet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 record.Index) *MemoryIndexStorageMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &MemoryIndexStorageMockSetExpectation{}
	expectation.input = &MemoryIndexStorageMockSetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of MemoryIndexStorage.Set method
func (m *mMemoryIndexStorageMockSet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 record.Index)) *MemoryIndexStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/object.MemoryIndexStorage interface
func (m *MemoryIndexStorageMock) Set(p context.Context, p1 insolar.PulseNumber, p2 record.Index) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemoryIndexStorageMock.Set. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MemoryIndexStorageMockSetInput{p, p1, p2}, "MemoryIndexStorage.Set got unexpected parameters")

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MemoryIndexStorageMockSetInput{p, p1, p2}, "MemoryIndexStorage.Set got unexpected parameters")
		}

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to MemoryIndexStorageMock.Set. %v %v %v", p, p1, p2)
		return
	}

	m.SetFunc(p, p1, p2)
}

//SetMinimockCounter returns a count of MemoryIndexStorageMock.SetFunc invocations
func (m *MemoryIndexStorageMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of MemoryIndexStorageMock.Set invocations
func (m *MemoryIndexStorageMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *MemoryIndexStorageMock) SetFinished() bool {
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
func (m *MemoryIndexStorageMock) ValidateCallCounters() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to MemoryIndexStorageMock.ForID")
	}

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to MemoryIndexStorageMock.ForPulse")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to MemoryIndexStorageMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MemoryIndexStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MemoryIndexStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MemoryIndexStorageMock) MinimockFinish() {

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to MemoryIndexStorageMock.ForID")
	}

	if !m.ForPulseFinished() {
		m.t.Fatal("Expected call to MemoryIndexStorageMock.ForPulse")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to MemoryIndexStorageMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MemoryIndexStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MemoryIndexStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForIDFinished()
		ok = ok && m.ForPulseFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForIDFinished() {
				m.t.Error("Expected call to MemoryIndexStorageMock.ForID")
			}

			if !m.ForPulseFinished() {
				m.t.Error("Expected call to MemoryIndexStorageMock.ForPulse")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to MemoryIndexStorageMock.Set")
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
func (m *MemoryIndexStorageMock) AllMocksCalled() bool {

	if !m.ForIDFinished() {
		return false
	}

	if !m.ForPulseFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}
