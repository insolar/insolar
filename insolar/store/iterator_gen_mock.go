package store

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
)

// IteratorMock implements store.Iterator
type IteratorMock struct {
	t minimock.Tester

	funcClose          func()
	inspectFuncClose   func()
	afterCloseCounter  uint64
	beforeCloseCounter uint64
	CloseMock          mIteratorMockClose

	funcKey          func() (ba1 []byte)
	inspectFuncKey   func()
	afterKeyCounter  uint64
	beforeKeyCounter uint64
	KeyMock          mIteratorMockKey

	funcNext          func() (b1 bool)
	inspectFuncNext   func()
	afterNextCounter  uint64
	beforeNextCounter uint64
	NextMock          mIteratorMockNext

	funcValue          func() (ba1 []byte, err error)
	inspectFuncValue   func()
	afterValueCounter  uint64
	beforeValueCounter uint64
	ValueMock          mIteratorMockValue
}

// NewIteratorMock returns a mock for store.Iterator
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
	mock               *IteratorMock
	defaultExpectation *IteratorMockCloseExpectation
	expectations       []*IteratorMockCloseExpectation
}

// IteratorMockCloseExpectation specifies expectation struct of the Iterator.Close
type IteratorMockCloseExpectation struct {
	mock *IteratorMock

	Counter uint64
}

// Expect sets up expected params for Iterator.Close
func (mmClose *mIteratorMockClose) Expect() *mIteratorMockClose {
	if mmClose.mock.funcClose != nil {
		mmClose.mock.t.Fatalf("IteratorMock.Close mock is already set by Set")
	}

	if mmClose.defaultExpectation == nil {
		mmClose.defaultExpectation = &IteratorMockCloseExpectation{}
	}

	return mmClose
}

// Inspect accepts an inspector function that has same arguments as the Iterator.Close
func (mmClose *mIteratorMockClose) Inspect(f func()) *mIteratorMockClose {
	if mmClose.mock.inspectFuncClose != nil {
		mmClose.mock.t.Fatalf("Inspect function is already set for IteratorMock.Close")
	}

	mmClose.mock.inspectFuncClose = f

	return mmClose
}

// Return sets up results that will be returned by Iterator.Close
func (mmClose *mIteratorMockClose) Return() *IteratorMock {
	if mmClose.mock.funcClose != nil {
		mmClose.mock.t.Fatalf("IteratorMock.Close mock is already set by Set")
	}

	if mmClose.defaultExpectation == nil {
		mmClose.defaultExpectation = &IteratorMockCloseExpectation{mock: mmClose.mock}
	}

	return mmClose.mock
}

//Set uses given function f to mock the Iterator.Close method
func (mmClose *mIteratorMockClose) Set(f func()) *IteratorMock {
	if mmClose.defaultExpectation != nil {
		mmClose.mock.t.Fatalf("Default expectation is already set for the Iterator.Close method")
	}

	if len(mmClose.expectations) > 0 {
		mmClose.mock.t.Fatalf("Some expectations are already set for the Iterator.Close method")
	}

	mmClose.mock.funcClose = f
	return mmClose.mock
}

// Close implements store.Iterator
func (mmClose *IteratorMock) Close() {
	mm_atomic.AddUint64(&mmClose.beforeCloseCounter, 1)
	defer mm_atomic.AddUint64(&mmClose.afterCloseCounter, 1)

	if mmClose.inspectFuncClose != nil {
		mmClose.inspectFuncClose()
	}

	if mmClose.CloseMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmClose.CloseMock.defaultExpectation.Counter, 1)

		return

	}
	if mmClose.funcClose != nil {
		mmClose.funcClose()
		return
	}
	mmClose.t.Fatalf("Unexpected call to IteratorMock.Close.")

}

// CloseAfterCounter returns a count of finished IteratorMock.Close invocations
func (mmClose *IteratorMock) CloseAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmClose.afterCloseCounter)
}

// CloseBeforeCounter returns a count of IteratorMock.Close invocations
func (mmClose *IteratorMock) CloseBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmClose.beforeCloseCounter)
}

// MinimockCloseDone returns true if the count of the Close invocations corresponds
// the number of defined expectations
func (m *IteratorMock) MinimockCloseDone() bool {
	for _, e := range m.CloseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CloseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcClose != nil && mm_atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		return false
	}
	return true
}

// MinimockCloseInspect logs each unmet expectation
func (m *IteratorMock) MinimockCloseInspect() {
	for _, e := range m.CloseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to IteratorMock.Close")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CloseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Close")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcClose != nil && mm_atomic.LoadUint64(&m.afterCloseCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Close")
	}
}

type mIteratorMockKey struct {
	mock               *IteratorMock
	defaultExpectation *IteratorMockKeyExpectation
	expectations       []*IteratorMockKeyExpectation
}

// IteratorMockKeyExpectation specifies expectation struct of the Iterator.Key
type IteratorMockKeyExpectation struct {
	mock *IteratorMock

	results *IteratorMockKeyResults
	Counter uint64
}

// IteratorMockKeyResults contains results of the Iterator.Key
type IteratorMockKeyResults struct {
	ba1 []byte
}

// Expect sets up expected params for Iterator.Key
func (mmKey *mIteratorMockKey) Expect() *mIteratorMockKey {
	if mmKey.mock.funcKey != nil {
		mmKey.mock.t.Fatalf("IteratorMock.Key mock is already set by Set")
	}

	if mmKey.defaultExpectation == nil {
		mmKey.defaultExpectation = &IteratorMockKeyExpectation{}
	}

	return mmKey
}

// Inspect accepts an inspector function that has same arguments as the Iterator.Key
func (mmKey *mIteratorMockKey) Inspect(f func()) *mIteratorMockKey {
	if mmKey.mock.inspectFuncKey != nil {
		mmKey.mock.t.Fatalf("Inspect function is already set for IteratorMock.Key")
	}

	mmKey.mock.inspectFuncKey = f

	return mmKey
}

// Return sets up results that will be returned by Iterator.Key
func (mmKey *mIteratorMockKey) Return(ba1 []byte) *IteratorMock {
	if mmKey.mock.funcKey != nil {
		mmKey.mock.t.Fatalf("IteratorMock.Key mock is already set by Set")
	}

	if mmKey.defaultExpectation == nil {
		mmKey.defaultExpectation = &IteratorMockKeyExpectation{mock: mmKey.mock}
	}
	mmKey.defaultExpectation.results = &IteratorMockKeyResults{ba1}
	return mmKey.mock
}

//Set uses given function f to mock the Iterator.Key method
func (mmKey *mIteratorMockKey) Set(f func() (ba1 []byte)) *IteratorMock {
	if mmKey.defaultExpectation != nil {
		mmKey.mock.t.Fatalf("Default expectation is already set for the Iterator.Key method")
	}

	if len(mmKey.expectations) > 0 {
		mmKey.mock.t.Fatalf("Some expectations are already set for the Iterator.Key method")
	}

	mmKey.mock.funcKey = f
	return mmKey.mock
}

// Key implements store.Iterator
func (mmKey *IteratorMock) Key() (ba1 []byte) {
	mm_atomic.AddUint64(&mmKey.beforeKeyCounter, 1)
	defer mm_atomic.AddUint64(&mmKey.afterKeyCounter, 1)

	if mmKey.inspectFuncKey != nil {
		mmKey.inspectFuncKey()
	}

	if mmKey.KeyMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmKey.KeyMock.defaultExpectation.Counter, 1)

		results := mmKey.KeyMock.defaultExpectation.results
		if results == nil {
			mmKey.t.Fatal("No results are set for the IteratorMock.Key")
		}
		return (*results).ba1
	}
	if mmKey.funcKey != nil {
		return mmKey.funcKey()
	}
	mmKey.t.Fatalf("Unexpected call to IteratorMock.Key.")
	return
}

// KeyAfterCounter returns a count of finished IteratorMock.Key invocations
func (mmKey *IteratorMock) KeyAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmKey.afterKeyCounter)
}

// KeyBeforeCounter returns a count of IteratorMock.Key invocations
func (mmKey *IteratorMock) KeyBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmKey.beforeKeyCounter)
}

// MinimockKeyDone returns true if the count of the Key invocations corresponds
// the number of defined expectations
func (m *IteratorMock) MinimockKeyDone() bool {
	for _, e := range m.KeyMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.KeyMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterKeyCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcKey != nil && mm_atomic.LoadUint64(&m.afterKeyCounter) < 1 {
		return false
	}
	return true
}

// MinimockKeyInspect logs each unmet expectation
func (m *IteratorMock) MinimockKeyInspect() {
	for _, e := range m.KeyMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to IteratorMock.Key")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.KeyMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterKeyCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Key")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcKey != nil && mm_atomic.LoadUint64(&m.afterKeyCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Key")
	}
}

type mIteratorMockNext struct {
	mock               *IteratorMock
	defaultExpectation *IteratorMockNextExpectation
	expectations       []*IteratorMockNextExpectation
}

// IteratorMockNextExpectation specifies expectation struct of the Iterator.Next
type IteratorMockNextExpectation struct {
	mock *IteratorMock

	results *IteratorMockNextResults
	Counter uint64
}

// IteratorMockNextResults contains results of the Iterator.Next
type IteratorMockNextResults struct {
	b1 bool
}

// Expect sets up expected params for Iterator.Next
func (mmNext *mIteratorMockNext) Expect() *mIteratorMockNext {
	if mmNext.mock.funcNext != nil {
		mmNext.mock.t.Fatalf("IteratorMock.Next mock is already set by Set")
	}

	if mmNext.defaultExpectation == nil {
		mmNext.defaultExpectation = &IteratorMockNextExpectation{}
	}

	return mmNext
}

// Inspect accepts an inspector function that has same arguments as the Iterator.Next
func (mmNext *mIteratorMockNext) Inspect(f func()) *mIteratorMockNext {
	if mmNext.mock.inspectFuncNext != nil {
		mmNext.mock.t.Fatalf("Inspect function is already set for IteratorMock.Next")
	}

	mmNext.mock.inspectFuncNext = f

	return mmNext
}

// Return sets up results that will be returned by Iterator.Next
func (mmNext *mIteratorMockNext) Return(b1 bool) *IteratorMock {
	if mmNext.mock.funcNext != nil {
		mmNext.mock.t.Fatalf("IteratorMock.Next mock is already set by Set")
	}

	if mmNext.defaultExpectation == nil {
		mmNext.defaultExpectation = &IteratorMockNextExpectation{mock: mmNext.mock}
	}
	mmNext.defaultExpectation.results = &IteratorMockNextResults{b1}
	return mmNext.mock
}

//Set uses given function f to mock the Iterator.Next method
func (mmNext *mIteratorMockNext) Set(f func() (b1 bool)) *IteratorMock {
	if mmNext.defaultExpectation != nil {
		mmNext.mock.t.Fatalf("Default expectation is already set for the Iterator.Next method")
	}

	if len(mmNext.expectations) > 0 {
		mmNext.mock.t.Fatalf("Some expectations are already set for the Iterator.Next method")
	}

	mmNext.mock.funcNext = f
	return mmNext.mock
}

// Next implements store.Iterator
func (mmNext *IteratorMock) Next() (b1 bool) {
	mm_atomic.AddUint64(&mmNext.beforeNextCounter, 1)
	defer mm_atomic.AddUint64(&mmNext.afterNextCounter, 1)

	if mmNext.inspectFuncNext != nil {
		mmNext.inspectFuncNext()
	}

	if mmNext.NextMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmNext.NextMock.defaultExpectation.Counter, 1)

		results := mmNext.NextMock.defaultExpectation.results
		if results == nil {
			mmNext.t.Fatal("No results are set for the IteratorMock.Next")
		}
		return (*results).b1
	}
	if mmNext.funcNext != nil {
		return mmNext.funcNext()
	}
	mmNext.t.Fatalf("Unexpected call to IteratorMock.Next.")
	return
}

// NextAfterCounter returns a count of finished IteratorMock.Next invocations
func (mmNext *IteratorMock) NextAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmNext.afterNextCounter)
}

// NextBeforeCounter returns a count of IteratorMock.Next invocations
func (mmNext *IteratorMock) NextBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmNext.beforeNextCounter)
}

// MinimockNextDone returns true if the count of the Next invocations corresponds
// the number of defined expectations
func (m *IteratorMock) MinimockNextDone() bool {
	for _, e := range m.NextMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.NextMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterNextCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcNext != nil && mm_atomic.LoadUint64(&m.afterNextCounter) < 1 {
		return false
	}
	return true
}

// MinimockNextInspect logs each unmet expectation
func (m *IteratorMock) MinimockNextInspect() {
	for _, e := range m.NextMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to IteratorMock.Next")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.NextMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterNextCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Next")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcNext != nil && mm_atomic.LoadUint64(&m.afterNextCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Next")
	}
}

type mIteratorMockValue struct {
	mock               *IteratorMock
	defaultExpectation *IteratorMockValueExpectation
	expectations       []*IteratorMockValueExpectation
}

// IteratorMockValueExpectation specifies expectation struct of the Iterator.Value
type IteratorMockValueExpectation struct {
	mock *IteratorMock

	results *IteratorMockValueResults
	Counter uint64
}

// IteratorMockValueResults contains results of the Iterator.Value
type IteratorMockValueResults struct {
	ba1 []byte
	err error
}

// Expect sets up expected params for Iterator.Value
func (mmValue *mIteratorMockValue) Expect() *mIteratorMockValue {
	if mmValue.mock.funcValue != nil {
		mmValue.mock.t.Fatalf("IteratorMock.Value mock is already set by Set")
	}

	if mmValue.defaultExpectation == nil {
		mmValue.defaultExpectation = &IteratorMockValueExpectation{}
	}

	return mmValue
}

// Inspect accepts an inspector function that has same arguments as the Iterator.Value
func (mmValue *mIteratorMockValue) Inspect(f func()) *mIteratorMockValue {
	if mmValue.mock.inspectFuncValue != nil {
		mmValue.mock.t.Fatalf("Inspect function is already set for IteratorMock.Value")
	}

	mmValue.mock.inspectFuncValue = f

	return mmValue
}

// Return sets up results that will be returned by Iterator.Value
func (mmValue *mIteratorMockValue) Return(ba1 []byte, err error) *IteratorMock {
	if mmValue.mock.funcValue != nil {
		mmValue.mock.t.Fatalf("IteratorMock.Value mock is already set by Set")
	}

	if mmValue.defaultExpectation == nil {
		mmValue.defaultExpectation = &IteratorMockValueExpectation{mock: mmValue.mock}
	}
	mmValue.defaultExpectation.results = &IteratorMockValueResults{ba1, err}
	return mmValue.mock
}

//Set uses given function f to mock the Iterator.Value method
func (mmValue *mIteratorMockValue) Set(f func() (ba1 []byte, err error)) *IteratorMock {
	if mmValue.defaultExpectation != nil {
		mmValue.mock.t.Fatalf("Default expectation is already set for the Iterator.Value method")
	}

	if len(mmValue.expectations) > 0 {
		mmValue.mock.t.Fatalf("Some expectations are already set for the Iterator.Value method")
	}

	mmValue.mock.funcValue = f
	return mmValue.mock
}

// Value implements store.Iterator
func (mmValue *IteratorMock) Value() (ba1 []byte, err error) {
	mm_atomic.AddUint64(&mmValue.beforeValueCounter, 1)
	defer mm_atomic.AddUint64(&mmValue.afterValueCounter, 1)

	if mmValue.inspectFuncValue != nil {
		mmValue.inspectFuncValue()
	}

	if mmValue.ValueMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmValue.ValueMock.defaultExpectation.Counter, 1)

		results := mmValue.ValueMock.defaultExpectation.results
		if results == nil {
			mmValue.t.Fatal("No results are set for the IteratorMock.Value")
		}
		return (*results).ba1, (*results).err
	}
	if mmValue.funcValue != nil {
		return mmValue.funcValue()
	}
	mmValue.t.Fatalf("Unexpected call to IteratorMock.Value.")
	return
}

// ValueAfterCounter returns a count of finished IteratorMock.Value invocations
func (mmValue *IteratorMock) ValueAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmValue.afterValueCounter)
}

// ValueBeforeCounter returns a count of IteratorMock.Value invocations
func (mmValue *IteratorMock) ValueBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmValue.beforeValueCounter)
}

// MinimockValueDone returns true if the count of the Value invocations corresponds
// the number of defined expectations
func (m *IteratorMock) MinimockValueDone() bool {
	for _, e := range m.ValueMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ValueMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterValueCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcValue != nil && mm_atomic.LoadUint64(&m.afterValueCounter) < 1 {
		return false
	}
	return true
}

// MinimockValueInspect logs each unmet expectation
func (m *IteratorMock) MinimockValueInspect() {
	for _, e := range m.ValueMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to IteratorMock.Value")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ValueMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterValueCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Value")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcValue != nil && mm_atomic.LoadUint64(&m.afterValueCounter) < 1 {
		m.t.Error("Expected call to IteratorMock.Value")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *IteratorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockCloseInspect()

		m.MinimockKeyInspect()

		m.MinimockNextInspect()

		m.MinimockValueInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *IteratorMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *IteratorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCloseDone() &&
		m.MinimockKeyDone() &&
		m.MinimockNextDone() &&
		m.MinimockValueDone()
}
