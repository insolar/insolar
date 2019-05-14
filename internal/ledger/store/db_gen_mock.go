package store

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DB" can be found in github.com/insolar/insolar/internal/ledger/store
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//DBMock implements github.com/insolar/insolar/internal/ledger/store.DB
type DBMock struct {
	t minimock.Tester

	GetFunc       func(p Key) (r []byte, r1 error)
	GetCounter    uint64
	GetPreCounter uint64
	GetMock       mDBMockGet

	NewIteratorFunc       func(p Scope) (r Iterator)
	NewIteratorCounter    uint64
	NewIteratorPreCounter uint64
	NewIteratorMock       mDBMockNewIterator

	SetFunc       func(p Key, p1 []byte) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mDBMockSet
}

//NewDBMock returns a mock for github.com/insolar/insolar/internal/ledger/store.DB
func NewDBMock(t minimock.Tester) *DBMock {
	m := &DBMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetMock = mDBMockGet{mock: m}
	m.NewIteratorMock = mDBMockNewIterator{mock: m}
	m.SetMock = mDBMockSet{mock: m}

	return m
}

type mDBMockGet struct {
	mock              *DBMock
	mainExpectation   *DBMockGetExpectation
	expectationSeries []*DBMockGetExpectation
}

type DBMockGetExpectation struct {
	input  *DBMockGetInput
	result *DBMockGetResult
}

type DBMockGetInput struct {
	p Key
}

type DBMockGetResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of DB.Get is expected from 1 to Infinity times
func (m *mDBMockGet) Expect(p Key) *mDBMockGet {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBMockGetExpectation{}
	}
	m.mainExpectation.input = &DBMockGetInput{p}
	return m
}

//Return specifies results of invocation of DB.Get
func (m *mDBMockGet) Return(r []byte, r1 error) *DBMock {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBMockGetExpectation{}
	}
	m.mainExpectation.result = &DBMockGetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DB.Get is expected once
func (m *mDBMockGet) ExpectOnce(p Key) *DBMockGetExpectation {
	m.mock.GetFunc = nil
	m.mainExpectation = nil

	expectation := &DBMockGetExpectation{}
	expectation.input = &DBMockGetInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBMockGetExpectation) Return(r []byte, r1 error) {
	e.result = &DBMockGetResult{r, r1}
}

//Set uses given function f as a mock of DB.Get method
func (m *mDBMockGet) Set(f func(p Key) (r []byte, r1 error)) *DBMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetFunc = f
	return m.mock
}

//Get implements github.com/insolar/insolar/internal/ledger/store.DB interface
func (m *DBMock) Get(p Key) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.GetPreCounter, 1)
	defer atomic.AddUint64(&m.GetCounter, 1)

	if len(m.GetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBMock.Get. %v", p)
			return
		}

		input := m.GetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBMockGetInput{p}, "DB.Get got unexpected parameters")

		result := m.GetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBMock.Get")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMock.mainExpectation != nil {

		input := m.GetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBMockGetInput{p}, "DB.Get got unexpected parameters")
		}

		result := m.GetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBMock.Get")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetFunc == nil {
		m.t.Fatalf("Unexpected call to DBMock.Get. %v", p)
		return
	}

	return m.GetFunc(p)
}

//GetMinimockCounter returns a count of DBMock.GetFunc invocations
func (m *DBMock) GetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCounter)
}

//GetMinimockPreCounter returns the value of DBMock.Get invocations
func (m *DBMock) GetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreCounter)
}

//GetFinished returns true if mock invocations count is ok
func (m *DBMock) GetFinished() bool {
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

type mDBMockNewIterator struct {
	mock              *DBMock
	mainExpectation   *DBMockNewIteratorExpectation
	expectationSeries []*DBMockNewIteratorExpectation
}

type DBMockNewIteratorExpectation struct {
	input  *DBMockNewIteratorInput
	result *DBMockNewIteratorResult
}

type DBMockNewIteratorInput struct {
	p Scope
}

type DBMockNewIteratorResult struct {
	r Iterator
}

//Expect specifies that invocation of DB.NewIterator is expected from 1 to Infinity times
func (m *mDBMockNewIterator) Expect(p Scope) *mDBMockNewIterator {
	m.mock.NewIteratorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBMockNewIteratorExpectation{}
	}
	m.mainExpectation.input = &DBMockNewIteratorInput{p}
	return m
}

//Return specifies results of invocation of DB.NewIterator
func (m *mDBMockNewIterator) Return(r Iterator) *DBMock {
	m.mock.NewIteratorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBMockNewIteratorExpectation{}
	}
	m.mainExpectation.result = &DBMockNewIteratorResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DB.NewIterator is expected once
func (m *mDBMockNewIterator) ExpectOnce(p Scope) *DBMockNewIteratorExpectation {
	m.mock.NewIteratorFunc = nil
	m.mainExpectation = nil

	expectation := &DBMockNewIteratorExpectation{}
	expectation.input = &DBMockNewIteratorInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBMockNewIteratorExpectation) Return(r Iterator) {
	e.result = &DBMockNewIteratorResult{r}
}

//Set uses given function f as a mock of DB.NewIterator method
func (m *mDBMockNewIterator) Set(f func(p Scope) (r Iterator)) *DBMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NewIteratorFunc = f
	return m.mock
}

//NewIterator implements github.com/insolar/insolar/internal/ledger/store.DB interface
func (m *DBMock) NewIterator(p Scope) (r Iterator) {
	counter := atomic.AddUint64(&m.NewIteratorPreCounter, 1)
	defer atomic.AddUint64(&m.NewIteratorCounter, 1)

	if len(m.NewIteratorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NewIteratorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBMock.NewIterator. %v", p)
			return
		}

		input := m.NewIteratorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBMockNewIteratorInput{p}, "DB.NewIterator got unexpected parameters")

		result := m.NewIteratorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBMock.NewIterator")
			return
		}

		r = result.r

		return
	}

	if m.NewIteratorMock.mainExpectation != nil {

		input := m.NewIteratorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBMockNewIteratorInput{p}, "DB.NewIterator got unexpected parameters")
		}

		result := m.NewIteratorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBMock.NewIterator")
		}

		r = result.r

		return
	}

	if m.NewIteratorFunc == nil {
		m.t.Fatalf("Unexpected call to DBMock.NewIterator. %v", p)
		return
	}

	return m.NewIteratorFunc(p)
}

//NewIteratorMinimockCounter returns a count of DBMock.NewIteratorFunc invocations
func (m *DBMock) NewIteratorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewIteratorCounter)
}

//NewIteratorMinimockPreCounter returns the value of DBMock.NewIterator invocations
func (m *DBMock) NewIteratorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewIteratorPreCounter)
}

//NewIteratorFinished returns true if mock invocations count is ok
func (m *DBMock) NewIteratorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NewIteratorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NewIteratorCounter) == uint64(len(m.NewIteratorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NewIteratorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NewIteratorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NewIteratorFunc != nil {
		return atomic.LoadUint64(&m.NewIteratorCounter) > 0
	}

	return true
}

type mDBMockSet struct {
	mock              *DBMock
	mainExpectation   *DBMockSetExpectation
	expectationSeries []*DBMockSetExpectation
}

type DBMockSetExpectation struct {
	input  *DBMockSetInput
	result *DBMockSetResult
}

type DBMockSetInput struct {
	p  Key
	p1 []byte
}

type DBMockSetResult struct {
	r error
}

//Expect specifies that invocation of DB.Set is expected from 1 to Infinity times
func (m *mDBMockSet) Expect(p Key, p1 []byte) *mDBMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBMockSetExpectation{}
	}
	m.mainExpectation.input = &DBMockSetInput{p, p1}
	return m
}

//Return specifies results of invocation of DB.Set
func (m *mDBMockSet) Return(r error) *DBMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBMockSetExpectation{}
	}
	m.mainExpectation.result = &DBMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DB.Set is expected once
func (m *mDBMockSet) ExpectOnce(p Key, p1 []byte) *DBMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &DBMockSetExpectation{}
	expectation.input = &DBMockSetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBMockSetExpectation) Return(r error) {
	e.result = &DBMockSetResult{r}
}

//Set uses given function f as a mock of DB.Set method
func (m *mDBMockSet) Set(f func(p Key, p1 []byte) (r error)) *DBMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/internal/ledger/store.DB interface
func (m *DBMock) Set(p Key, p1 []byte) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBMock.Set. %v %v", p, p1)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBMockSetInput{p, p1}, "DB.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBMockSetInput{p, p1}, "DB.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to DBMock.Set. %v %v", p, p1)
		return
	}

	return m.SetFunc(p, p1)
}

//SetMinimockCounter returns a count of DBMock.SetFunc invocations
func (m *DBMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of DBMock.Set invocations
func (m *DBMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *DBMock) SetFinished() bool {
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
func (m *DBMock) ValidateCallCounters() {

	if !m.GetFinished() {
		m.t.Fatal("Expected call to DBMock.Get")
	}

	if !m.NewIteratorFinished() {
		m.t.Fatal("Expected call to DBMock.NewIterator")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to DBMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DBMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DBMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DBMock) MinimockFinish() {

	if !m.GetFinished() {
		m.t.Fatal("Expected call to DBMock.Get")
	}

	if !m.NewIteratorFinished() {
		m.t.Fatal("Expected call to DBMock.NewIterator")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to DBMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DBMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DBMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetFinished()
		ok = ok && m.NewIteratorFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetFinished() {
				m.t.Error("Expected call to DBMock.Get")
			}

			if !m.NewIteratorFinished() {
				m.t.Error("Expected call to DBMock.NewIterator")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to DBMock.Set")
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
func (m *DBMock) AllMocksCalled() bool {

	if !m.GetFinished() {
		return false
	}

	if !m.NewIteratorFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}
