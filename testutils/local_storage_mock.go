package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LocalStorage" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

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
	mock              *LocalStorageMock
	mainExpectation   *LocalStorageMockGetExpectation
	expectationSeries []*LocalStorageMockGetExpectation
}

type LocalStorageMockGetExpectation struct {
	input  *LocalStorageMockGetInput
	result *LocalStorageMockGetResult
}

type LocalStorageMockGetInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []byte
}

type LocalStorageMockGetResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of LocalStorage.Get is expected from 1 to Infinity times
func (m *mLocalStorageMockGet) Expect(p context.Context, p1 core.PulseNumber, p2 []byte) *mLocalStorageMockGet {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalStorageMockGetExpectation{}
	}
	m.mainExpectation.input = &LocalStorageMockGetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of LocalStorage.Get
func (m *mLocalStorageMockGet) Return(r []byte, r1 error) *LocalStorageMock {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalStorageMockGetExpectation{}
	}
	m.mainExpectation.result = &LocalStorageMockGetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalStorage.Get is expected once
func (m *mLocalStorageMockGet) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 []byte) *LocalStorageMockGetExpectation {
	m.mock.GetFunc = nil
	m.mainExpectation = nil

	expectation := &LocalStorageMockGetExpectation{}
	expectation.input = &LocalStorageMockGetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalStorageMockGetExpectation) Return(r []byte, r1 error) {
	e.result = &LocalStorageMockGetResult{r, r1}
}

//Set uses given function f as a mock of LocalStorage.Get method
func (m *mLocalStorageMockGet) Set(f func(p context.Context, p1 core.PulseNumber, p2 []byte) (r []byte, r1 error)) *LocalStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetFunc = f
	return m.mock
}

//Get implements github.com/insolar/insolar/core.LocalStorage interface
func (m *LocalStorageMock) Get(p context.Context, p1 core.PulseNumber, p2 []byte) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.GetPreCounter, 1)
	defer atomic.AddUint64(&m.GetCounter, 1)

	if len(m.GetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalStorageMock.Get. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LocalStorageMockGetInput{p, p1, p2}, "LocalStorage.Get got unexpected parameters")

		result := m.GetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalStorageMock.Get")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMock.mainExpectation != nil {

		input := m.GetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LocalStorageMockGetInput{p, p1, p2}, "LocalStorage.Get got unexpected parameters")
		}

		result := m.GetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalStorageMock.Get")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetFunc == nil {
		m.t.Fatalf("Unexpected call to LocalStorageMock.Get. %v %v %v", p, p1, p2)
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

//GetFinished returns true if mock invocations count is ok
func (m *LocalStorageMock) GetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCounter) == uint64(len(m.GetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetFunc != nil {
		return atomic.LoadUint64(&m.GetCounter) > 0
	}

	return true
}

type mLocalStorageMockIterate struct {
	mock              *LocalStorageMock
	mainExpectation   *LocalStorageMockIterateExpectation
	expectationSeries []*LocalStorageMockIterateExpectation
}

type LocalStorageMockIterateExpectation struct {
	input  *LocalStorageMockIterateInput
	result *LocalStorageMockIterateResult
}

type LocalStorageMockIterateInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []byte
	p3 func(p []byte, p1 []byte) (r error)
}

type LocalStorageMockIterateResult struct {
	r error
}

//Expect specifies that invocation of LocalStorage.Iterate is expected from 1 to Infinity times
func (m *mLocalStorageMockIterate) Expect(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) *mLocalStorageMockIterate {
	m.mock.IterateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalStorageMockIterateExpectation{}
	}
	m.mainExpectation.input = &LocalStorageMockIterateInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of LocalStorage.Iterate
func (m *mLocalStorageMockIterate) Return(r error) *LocalStorageMock {
	m.mock.IterateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalStorageMockIterateExpectation{}
	}
	m.mainExpectation.result = &LocalStorageMockIterateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalStorage.Iterate is expected once
func (m *mLocalStorageMockIterate) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) *LocalStorageMockIterateExpectation {
	m.mock.IterateFunc = nil
	m.mainExpectation = nil

	expectation := &LocalStorageMockIterateExpectation{}
	expectation.input = &LocalStorageMockIterateInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalStorageMockIterateExpectation) Return(r error) {
	e.result = &LocalStorageMockIterateResult{r}
}

//Set uses given function f as a mock of LocalStorage.Iterate method
func (m *mLocalStorageMockIterate) Set(f func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) (r error)) *LocalStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IterateFunc = f
	return m.mock
}

//Iterate implements github.com/insolar/insolar/core.LocalStorage interface
func (m *LocalStorageMock) Iterate(p context.Context, p1 core.PulseNumber, p2 []byte, p3 func(p []byte, p1 []byte) (r error)) (r error) {
	counter := atomic.AddUint64(&m.IteratePreCounter, 1)
	defer atomic.AddUint64(&m.IterateCounter, 1)

	if len(m.IterateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IterateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalStorageMock.Iterate. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.IterateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LocalStorageMockIterateInput{p, p1, p2, p3}, "LocalStorage.Iterate got unexpected parameters")

		result := m.IterateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalStorageMock.Iterate")
			return
		}

		r = result.r

		return
	}

	if m.IterateMock.mainExpectation != nil {

		input := m.IterateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LocalStorageMockIterateInput{p, p1, p2, p3}, "LocalStorage.Iterate got unexpected parameters")
		}

		result := m.IterateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalStorageMock.Iterate")
		}

		r = result.r

		return
	}

	if m.IterateFunc == nil {
		m.t.Fatalf("Unexpected call to LocalStorageMock.Iterate. %v %v %v %v", p, p1, p2, p3)
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

//IterateFinished returns true if mock invocations count is ok
func (m *LocalStorageMock) IterateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IterateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IterateCounter) == uint64(len(m.IterateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IterateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IterateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IterateFunc != nil {
		return atomic.LoadUint64(&m.IterateCounter) > 0
	}

	return true
}

type mLocalStorageMockSet struct {
	mock              *LocalStorageMock
	mainExpectation   *LocalStorageMockSetExpectation
	expectationSeries []*LocalStorageMockSetExpectation
}

type LocalStorageMockSetExpectation struct {
	input  *LocalStorageMockSetInput
	result *LocalStorageMockSetResult
}

type LocalStorageMockSetInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []byte
	p3 []byte
}

type LocalStorageMockSetResult struct {
	r error
}

//Expect specifies that invocation of LocalStorage.Set is expected from 1 to Infinity times
func (m *mLocalStorageMockSet) Expect(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) *mLocalStorageMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalStorageMockSetExpectation{}
	}
	m.mainExpectation.input = &LocalStorageMockSetInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of LocalStorage.Set
func (m *mLocalStorageMockSet) Return(r error) *LocalStorageMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalStorageMockSetExpectation{}
	}
	m.mainExpectation.result = &LocalStorageMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalStorage.Set is expected once
func (m *mLocalStorageMockSet) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) *LocalStorageMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &LocalStorageMockSetExpectation{}
	expectation.input = &LocalStorageMockSetInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalStorageMockSetExpectation) Return(r error) {
	e.result = &LocalStorageMockSetResult{r}
}

//Set uses given function f as a mock of LocalStorage.Set method
func (m *mLocalStorageMockSet) Set(f func(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) (r error)) *LocalStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/core.LocalStorage interface
func (m *LocalStorageMock) Set(p context.Context, p1 core.PulseNumber, p2 []byte, p3 []byte) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalStorageMock.Set. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, LocalStorageMockSetInput{p, p1, p2, p3}, "LocalStorage.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalStorageMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, LocalStorageMockSetInput{p, p1, p2, p3}, "LocalStorage.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalStorageMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to LocalStorageMock.Set. %v %v %v %v", p, p1, p2, p3)
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

//SetFinished returns true if mock invocations count is ok
func (m *LocalStorageMock) SetFinished() bool {
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
func (m *LocalStorageMock) ValidateCallCounters() {

	if !m.GetFinished() {
		m.t.Fatal("Expected call to LocalStorageMock.Get")
	}

	if !m.IterateFinished() {
		m.t.Fatal("Expected call to LocalStorageMock.Iterate")
	}

	if !m.SetFinished() {
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

	if !m.GetFinished() {
		m.t.Fatal("Expected call to LocalStorageMock.Get")
	}

	if !m.IterateFinished() {
		m.t.Fatal("Expected call to LocalStorageMock.Iterate")
	}

	if !m.SetFinished() {
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
		ok = ok && m.GetFinished()
		ok = ok && m.IterateFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetFinished() {
				m.t.Error("Expected call to LocalStorageMock.Get")
			}

			if !m.IterateFinished() {
				m.t.Error("Expected call to LocalStorageMock.Iterate")
			}

			if !m.SetFinished() {
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

	if !m.GetFinished() {
		return false
	}

	if !m.IterateFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}
