package store

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Iterator" can be found in github.com/insolar/insolar/internal/ledger/store
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//IteratorMock implements github.com/insolar/insolar/internal/ledger/store.Iterator
type IteratorMock struct {
	t minimock.Tester

	CloseFunc       func()
	CloseCounter    uint64
	ClosePreCounter uint64
	CloseMock       mIteratorMockClose

	KeyFunc       func() (r []byte)
	KeyCounter    uint64
	KeyPreCounter uint64
	KeyMock       mIteratorMockKey

	NextFunc       func() (r bool)
	NextCounter    uint64
	NextPreCounter uint64
	NextMock       mIteratorMockNext

	ValueFunc       func() (r []byte)
	ValueCounter    uint64
	ValuePreCounter uint64
	ValueMock       mIteratorMockValue
}

//NewIteratorMock returns a mock for github.com/insolar/insolar/internal/ledger/store.Iterator
func NewIteratorMock(t minimock.Tester) *IteratorMock {
	m := &IteratorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseMock = mIteratorMockClose{mock: m}
	m.KeyMock = mIteratorMockKey{mock: m}
	m.NextMock = mIteratorMockNext{mock: m}
	m.ValueMock = mIteratorMockValue{mock: m}

	return m
}

type mIteratorMockClose struct {
	mock              *IteratorMock
	mainExpectation   *IteratorMockCloseExpectation
	expectationSeries []*IteratorMockCloseExpectation
}

type IteratorMockCloseExpectation struct {
}

//Expect specifies that invocation of Iterator.Close is expected from 1 to Infinity times
func (m *mIteratorMockClose) Expect() *mIteratorMockClose {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockCloseExpectation{}
	}

	return m
}

//Return specifies results of invocation of Iterator.Close
func (m *mIteratorMockClose) Return() *IteratorMock {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockCloseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Iterator.Close is expected once
func (m *mIteratorMockClose) ExpectOnce() *IteratorMockCloseExpectation {
	m.mock.CloseFunc = nil
	m.mainExpectation = nil

	expectation := &IteratorMockCloseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Iterator.Close method
func (m *mIteratorMockClose) Set(f func()) *IteratorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloseFunc = f
	return m.mock
}

//Close implements github.com/insolar/insolar/internal/ledger/store.Iterator interface
func (m *IteratorMock) Close() {
	counter := atomic.AddUint64(&m.ClosePreCounter, 1)
	defer atomic.AddUint64(&m.CloseCounter, 1)

	if len(m.CloseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IteratorMock.Close.")
			return
		}

		return
	}

	if m.CloseMock.mainExpectation != nil {

		return
	}

	if m.CloseFunc == nil {
		m.t.Fatalf("Unexpected call to IteratorMock.Close.")
		return
	}

	m.CloseFunc()
}

//CloseMinimockCounter returns a count of IteratorMock.CloseFunc invocations
func (m *IteratorMock) CloseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseCounter)
}

//CloseMinimockPreCounter returns the value of IteratorMock.Close invocations
func (m *IteratorMock) CloseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClosePreCounter)
}

//CloseFinished returns true if mock invocations count is ok
func (m *IteratorMock) CloseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloseCounter) == uint64(len(m.CloseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloseFunc != nil {
		return atomic.LoadUint64(&m.CloseCounter) > 0
	}

	return true
}

type mIteratorMockKey struct {
	mock              *IteratorMock
	mainExpectation   *IteratorMockKeyExpectation
	expectationSeries []*IteratorMockKeyExpectation
}

type IteratorMockKeyExpectation struct {
	result *IteratorMockKeyResult
}

type IteratorMockKeyResult struct {
	r []byte
}

//Expect specifies that invocation of Iterator.Key is expected from 1 to Infinity times
func (m *mIteratorMockKey) Expect() *mIteratorMockKey {
	m.mock.KeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of Iterator.Key
func (m *mIteratorMockKey) Return(r []byte) *IteratorMock {
	m.mock.KeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockKeyExpectation{}
	}
	m.mainExpectation.result = &IteratorMockKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Iterator.Key is expected once
func (m *mIteratorMockKey) ExpectOnce() *IteratorMockKeyExpectation {
	m.mock.KeyFunc = nil
	m.mainExpectation = nil

	expectation := &IteratorMockKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IteratorMockKeyExpectation) Return(r []byte) {
	e.result = &IteratorMockKeyResult{r}
}

//Set uses given function f as a mock of Iterator.Key method
func (m *mIteratorMockKey) Set(f func() (r []byte)) *IteratorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.KeyFunc = f
	return m.mock
}

//Key implements github.com/insolar/insolar/internal/ledger/store.Iterator interface
func (m *IteratorMock) Key() (r []byte) {
	counter := atomic.AddUint64(&m.KeyPreCounter, 1)
	defer atomic.AddUint64(&m.KeyCounter, 1)

	if len(m.KeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.KeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IteratorMock.Key.")
			return
		}

		result := m.KeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IteratorMock.Key")
			return
		}

		r = result.r

		return
	}

	if m.KeyMock.mainExpectation != nil {

		result := m.KeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IteratorMock.Key")
		}

		r = result.r

		return
	}

	if m.KeyFunc == nil {
		m.t.Fatalf("Unexpected call to IteratorMock.Key.")
		return
	}

	return m.KeyFunc()
}

//KeyMinimockCounter returns a count of IteratorMock.KeyFunc invocations
func (m *IteratorMock) KeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.KeyCounter)
}

//KeyMinimockPreCounter returns the value of IteratorMock.Key invocations
func (m *IteratorMock) KeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.KeyPreCounter)
}

//KeyFinished returns true if mock invocations count is ok
func (m *IteratorMock) KeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.KeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.KeyCounter) == uint64(len(m.KeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.KeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.KeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.KeyFunc != nil {
		return atomic.LoadUint64(&m.KeyCounter) > 0
	}

	return true
}

type mIteratorMockNext struct {
	mock              *IteratorMock
	mainExpectation   *IteratorMockNextExpectation
	expectationSeries []*IteratorMockNextExpectation
}

type IteratorMockNextExpectation struct {
	result *IteratorMockNextResult
}

type IteratorMockNextResult struct {
	r bool
}

//Expect specifies that invocation of Iterator.Next is expected from 1 to Infinity times
func (m *mIteratorMockNext) Expect() *mIteratorMockNext {
	m.mock.NextFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockNextExpectation{}
	}

	return m
}

//Return specifies results of invocation of Iterator.Next
func (m *mIteratorMockNext) Return(r bool) *IteratorMock {
	m.mock.NextFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockNextExpectation{}
	}
	m.mainExpectation.result = &IteratorMockNextResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Iterator.Next is expected once
func (m *mIteratorMockNext) ExpectOnce() *IteratorMockNextExpectation {
	m.mock.NextFunc = nil
	m.mainExpectation = nil

	expectation := &IteratorMockNextExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IteratorMockNextExpectation) Return(r bool) {
	e.result = &IteratorMockNextResult{r}
}

//Set uses given function f as a mock of Iterator.Next method
func (m *mIteratorMockNext) Set(f func() (r bool)) *IteratorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NextFunc = f
	return m.mock
}

//Next implements github.com/insolar/insolar/internal/ledger/store.Iterator interface
func (m *IteratorMock) Next() (r bool) {
	counter := atomic.AddUint64(&m.NextPreCounter, 1)
	defer atomic.AddUint64(&m.NextCounter, 1)

	if len(m.NextMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NextMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IteratorMock.Next.")
			return
		}

		result := m.NextMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IteratorMock.Next")
			return
		}

		r = result.r

		return
	}

	if m.NextMock.mainExpectation != nil {

		result := m.NextMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IteratorMock.Next")
		}

		r = result.r

		return
	}

	if m.NextFunc == nil {
		m.t.Fatalf("Unexpected call to IteratorMock.Next.")
		return
	}

	return m.NextFunc()
}

//NextMinimockCounter returns a count of IteratorMock.NextFunc invocations
func (m *IteratorMock) NextMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NextCounter)
}

//NextMinimockPreCounter returns the value of IteratorMock.Next invocations
func (m *IteratorMock) NextMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NextPreCounter)
}

//NextFinished returns true if mock invocations count is ok
func (m *IteratorMock) NextFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NextMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NextCounter) == uint64(len(m.NextMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NextMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NextCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NextFunc != nil {
		return atomic.LoadUint64(&m.NextCounter) > 0
	}

	return true
}

type mIteratorMockValue struct {
	mock              *IteratorMock
	mainExpectation   *IteratorMockValueExpectation
	expectationSeries []*IteratorMockValueExpectation
}

type IteratorMockValueExpectation struct {
	result *IteratorMockValueResult
}

type IteratorMockValueResult struct {
	r []byte
}

//Expect specifies that invocation of Iterator.Value is expected from 1 to Infinity times
func (m *mIteratorMockValue) Expect() *mIteratorMockValue {
	m.mock.ValueFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockValueExpectation{}
	}

	return m
}

//Return specifies results of invocation of Iterator.Value
func (m *mIteratorMockValue) Return(r []byte) *IteratorMock {
	m.mock.ValueFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IteratorMockValueExpectation{}
	}
	m.mainExpectation.result = &IteratorMockValueResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Iterator.Value is expected once
func (m *mIteratorMockValue) ExpectOnce() *IteratorMockValueExpectation {
	m.mock.ValueFunc = nil
	m.mainExpectation = nil

	expectation := &IteratorMockValueExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *IteratorMockValueExpectation) Return(r []byte) {
	e.result = &IteratorMockValueResult{r}
}

//Set uses given function f as a mock of Iterator.Value method
func (m *mIteratorMockValue) Set(f func() (r []byte)) *IteratorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ValueFunc = f
	return m.mock
}

//Value implements github.com/insolar/insolar/internal/ledger/store.Iterator interface
func (m *IteratorMock) Value() (r []byte) {
	counter := atomic.AddUint64(&m.ValuePreCounter, 1)
	defer atomic.AddUint64(&m.ValueCounter, 1)

	if len(m.ValueMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ValueMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IteratorMock.Value.")
			return
		}

		result := m.ValueMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the IteratorMock.Value")
			return
		}

		r = result.r

		return
	}

	if m.ValueMock.mainExpectation != nil {

		result := m.ValueMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the IteratorMock.Value")
		}

		r = result.r

		return
	}

	if m.ValueFunc == nil {
		m.t.Fatalf("Unexpected call to IteratorMock.Value.")
		return
	}

	return m.ValueFunc()
}

//ValueMinimockCounter returns a count of IteratorMock.ValueFunc invocations
func (m *IteratorMock) ValueMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ValueCounter)
}

//ValueMinimockPreCounter returns the value of IteratorMock.Value invocations
func (m *IteratorMock) ValueMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ValuePreCounter)
}

//ValueFinished returns true if mock invocations count is ok
func (m *IteratorMock) ValueFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ValueMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ValueCounter) == uint64(len(m.ValueMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ValueMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ValueCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ValueFunc != nil {
		return atomic.LoadUint64(&m.ValueCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IteratorMock) ValidateCallCounters() {

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to IteratorMock.Close")
	}

	if !m.KeyFinished() {
		m.t.Fatal("Expected call to IteratorMock.Key")
	}

	if !m.NextFinished() {
		m.t.Fatal("Expected call to IteratorMock.Next")
	}

	if !m.ValueFinished() {
		m.t.Fatal("Expected call to IteratorMock.Value")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IteratorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IteratorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IteratorMock) MinimockFinish() {

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to IteratorMock.Close")
	}

	if !m.KeyFinished() {
		m.t.Fatal("Expected call to IteratorMock.Key")
	}

	if !m.NextFinished() {
		m.t.Fatal("Expected call to IteratorMock.Next")
	}

	if !m.ValueFinished() {
		m.t.Fatal("Expected call to IteratorMock.Value")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IteratorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IteratorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloseFinished()
		ok = ok && m.KeyFinished()
		ok = ok && m.NextFinished()
		ok = ok && m.ValueFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloseFinished() {
				m.t.Error("Expected call to IteratorMock.Close")
			}

			if !m.KeyFinished() {
				m.t.Error("Expected call to IteratorMock.Key")
			}

			if !m.NextFinished() {
				m.t.Error("Expected call to IteratorMock.Next")
			}

			if !m.ValueFinished() {
				m.t.Error("Expected call to IteratorMock.Value")
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
func (m *IteratorMock) AllMocksCalled() bool {

	if !m.CloseFinished() {
		return false
	}

	if !m.KeyFinished() {
		return false
	}

	if !m.NextFinished() {
		return false
	}

	if !m.ValueFinished() {
		return false
	}

	return true
}
