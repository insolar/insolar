package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DBContext" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	badger "github.com/dgraph-io/badger"
	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	record "github.com/insolar/insolar/ledger/storage/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//DBContextMock implements github.com/insolar/insolar/ledger/storage.DBContext
type DBContextMock struct {
	t minimock.Tester

	BeginTransactionFunc       func(p bool) (r *TransactionManager, r1 error)
	BeginTransactionCounter    uint64
	BeginTransactionPreCounter uint64
	BeginTransactionMock       mDBContextMockBeginTransaction

	CloseFunc       func() (r error)
	CloseCounter    uint64
	ClosePreCounter uint64
	CloseMock       mDBContextMockClose

	GetBadgerDBFunc       func() (r *badger.DB)
	GetBadgerDBCounter    uint64
	GetBadgerDBPreCounter uint64
	GetBadgerDBMock       mDBContextMockGetBadgerDB

	IterateRecordsOnPulseFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 func(p core.RecordID, p1 record.Record) (r error)) (r error)
	IterateRecordsOnPulseCounter    uint64
	IterateRecordsOnPulsePreCounter uint64
	IterateRecordsOnPulseMock       mDBContextMockIterateRecordsOnPulse

	StoreKeyValuesFunc       func(p context.Context, p1 []core.KV) (r error)
	StoreKeyValuesCounter    uint64
	StoreKeyValuesPreCounter uint64
	StoreKeyValuesMock       mDBContextMockStoreKeyValues

	UpdateFunc       func(p context.Context, p1 func(p *TransactionManager) (r error)) (r error)
	UpdateCounter    uint64
	UpdatePreCounter uint64
	UpdateMock       mDBContextMockUpdate

	ViewFunc       func(p context.Context, p1 func(p *TransactionManager) (r error)) (r error)
	ViewCounter    uint64
	ViewPreCounter uint64
	ViewMock       mDBContextMockView

	getFunc       func(p context.Context, p1 []byte) (r []byte, r1 error)
	getCounter    uint64
	getPreCounter uint64
	getMock       mDBContextMockget

	iterateFunc       func(p context.Context, p1 []byte, p2 func(p []byte, p1 []byte) (r error)) (r error)
	iterateCounter    uint64
	iteratePreCounter uint64
	iterateMock       mDBContextMockiterate

	setFunc       func(p context.Context, p1 []byte, p2 []byte) (r error)
	setCounter    uint64
	setPreCounter uint64
	setMock       mDBContextMockset

	waitingFlightFunc       func()
	waitingFlightCounter    uint64
	waitingFlightPreCounter uint64
	waitingFlightMock       mDBContextMockwaitingFlight
}

//NewDBContextMock returns a mock for github.com/insolar/insolar/ledger/storage.DBContext
func NewDBContextMock(t minimock.Tester) *DBContextMock {
	m := &DBContextMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BeginTransactionMock = mDBContextMockBeginTransaction{mock: m}
	m.CloseMock = mDBContextMockClose{mock: m}
	m.GetBadgerDBMock = mDBContextMockGetBadgerDB{mock: m}
	m.IterateRecordsOnPulseMock = mDBContextMockIterateRecordsOnPulse{mock: m}
	m.StoreKeyValuesMock = mDBContextMockStoreKeyValues{mock: m}
	m.UpdateMock = mDBContextMockUpdate{mock: m}
	m.ViewMock = mDBContextMockView{mock: m}
	m.getMock = mDBContextMockget{mock: m}
	m.iterateMock = mDBContextMockiterate{mock: m}
	m.setMock = mDBContextMockset{mock: m}
	m.waitingFlightMock = mDBContextMockwaitingFlight{mock: m}

	return m
}

type mDBContextMockBeginTransaction struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockBeginTransactionExpectation
	expectationSeries []*DBContextMockBeginTransactionExpectation
}

type DBContextMockBeginTransactionExpectation struct {
	input  *DBContextMockBeginTransactionInput
	result *DBContextMockBeginTransactionResult
}

type DBContextMockBeginTransactionInput struct {
	p bool
}

type DBContextMockBeginTransactionResult struct {
	r  *TransactionManager
	r1 error
}

//Expect specifies that invocation of DBContext.BeginTransaction is expected from 1 to Infinity times
func (m *mDBContextMockBeginTransaction) Expect(p bool) *mDBContextMockBeginTransaction {
	m.mock.BeginTransactionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockBeginTransactionExpectation{}
	}
	m.mainExpectation.input = &DBContextMockBeginTransactionInput{p}
	return m
}

//Return specifies results of invocation of DBContext.BeginTransaction
func (m *mDBContextMockBeginTransaction) Return(r *TransactionManager, r1 error) *DBContextMock {
	m.mock.BeginTransactionFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockBeginTransactionExpectation{}
	}
	m.mainExpectation.result = &DBContextMockBeginTransactionResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.BeginTransaction is expected once
func (m *mDBContextMockBeginTransaction) ExpectOnce(p bool) *DBContextMockBeginTransactionExpectation {
	m.mock.BeginTransactionFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockBeginTransactionExpectation{}
	expectation.input = &DBContextMockBeginTransactionInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockBeginTransactionExpectation) Return(r *TransactionManager, r1 error) {
	e.result = &DBContextMockBeginTransactionResult{r, r1}
}

//Set uses given function f as a mock of DBContext.BeginTransaction method
func (m *mDBContextMockBeginTransaction) Set(f func(p bool) (r *TransactionManager, r1 error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BeginTransactionFunc = f
	return m.mock
}

//BeginTransaction implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) BeginTransaction(p bool) (r *TransactionManager, r1 error) {
	counter := atomic.AddUint64(&m.BeginTransactionPreCounter, 1)
	defer atomic.AddUint64(&m.BeginTransactionCounter, 1)

	if len(m.BeginTransactionMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BeginTransactionMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.BeginTransaction. %v", p)
			return
		}

		input := m.BeginTransactionMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMockBeginTransactionInput{p}, "DBContext.BeginTransaction got unexpected parameters")

		result := m.BeginTransactionMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.BeginTransaction")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BeginTransactionMock.mainExpectation != nil {

		input := m.BeginTransactionMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMockBeginTransactionInput{p}, "DBContext.BeginTransaction got unexpected parameters")
		}

		result := m.BeginTransactionMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.BeginTransaction")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BeginTransactionFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.BeginTransaction. %v", p)
		return
	}

	return m.BeginTransactionFunc(p)
}

//BeginTransactionMinimockCounter returns a count of DBContextMock.BeginTransactionFunc invocations
func (m *DBContextMock) BeginTransactionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BeginTransactionCounter)
}

//BeginTransactionMinimockPreCounter returns the value of DBContextMock.BeginTransaction invocations
func (m *DBContextMock) BeginTransactionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BeginTransactionPreCounter)
}

//BeginTransactionFinished returns true if mock invocations count is ok
func (m *DBContextMock) BeginTransactionFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.BeginTransactionMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.BeginTransactionCounter) == uint64(len(m.BeginTransactionMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.BeginTransactionMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.BeginTransactionCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.BeginTransactionFunc != nil {
		return atomic.LoadUint64(&m.BeginTransactionCounter) > 0
	}

	return true
}

type mDBContextMockClose struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockCloseExpectation
	expectationSeries []*DBContextMockCloseExpectation
}

type DBContextMockCloseExpectation struct {
	result *DBContextMockCloseResult
}

type DBContextMockCloseResult struct {
	r error
}

//Expect specifies that invocation of DBContext.Close is expected from 1 to Infinity times
func (m *mDBContextMockClose) Expect() *mDBContextMockClose {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockCloseExpectation{}
	}

	return m
}

//Return specifies results of invocation of DBContext.Close
func (m *mDBContextMockClose) Return(r error) *DBContextMock {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockCloseExpectation{}
	}
	m.mainExpectation.result = &DBContextMockCloseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.Close is expected once
func (m *mDBContextMockClose) ExpectOnce() *DBContextMockCloseExpectation {
	m.mock.CloseFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockCloseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockCloseExpectation) Return(r error) {
	e.result = &DBContextMockCloseResult{r}
}

//Set uses given function f as a mock of DBContext.Close method
func (m *mDBContextMockClose) Set(f func() (r error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloseFunc = f
	return m.mock
}

//Close implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) Close() (r error) {
	counter := atomic.AddUint64(&m.ClosePreCounter, 1)
	defer atomic.AddUint64(&m.CloseCounter, 1)

	if len(m.CloseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.Close.")
			return
		}

		result := m.CloseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.Close")
			return
		}

		r = result.r

		return
	}

	if m.CloseMock.mainExpectation != nil {

		result := m.CloseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.Close")
		}

		r = result.r

		return
	}

	if m.CloseFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.Close.")
		return
	}

	return m.CloseFunc()
}

//CloseMinimockCounter returns a count of DBContextMock.CloseFunc invocations
func (m *DBContextMock) CloseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseCounter)
}

//CloseMinimockPreCounter returns the value of DBContextMock.Close invocations
func (m *DBContextMock) CloseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClosePreCounter)
}

//CloseFinished returns true if mock invocations count is ok
func (m *DBContextMock) CloseFinished() bool {
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

type mDBContextMockGetBadgerDB struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockGetBadgerDBExpectation
	expectationSeries []*DBContextMockGetBadgerDBExpectation
}

type DBContextMockGetBadgerDBExpectation struct {
	result *DBContextMockGetBadgerDBResult
}

type DBContextMockGetBadgerDBResult struct {
	r *badger.DB
}

//Expect specifies that invocation of DBContext.GetBadgerDB is expected from 1 to Infinity times
func (m *mDBContextMockGetBadgerDB) Expect() *mDBContextMockGetBadgerDB {
	m.mock.GetBadgerDBFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockGetBadgerDBExpectation{}
	}

	return m
}

//Return specifies results of invocation of DBContext.GetBadgerDB
func (m *mDBContextMockGetBadgerDB) Return(r *badger.DB) *DBContextMock {
	m.mock.GetBadgerDBFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockGetBadgerDBExpectation{}
	}
	m.mainExpectation.result = &DBContextMockGetBadgerDBResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.GetBadgerDB is expected once
func (m *mDBContextMockGetBadgerDB) ExpectOnce() *DBContextMockGetBadgerDBExpectation {
	m.mock.GetBadgerDBFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockGetBadgerDBExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockGetBadgerDBExpectation) Return(r *badger.DB) {
	e.result = &DBContextMockGetBadgerDBResult{r}
}

//Set uses given function f as a mock of DBContext.GetBadgerDB method
func (m *mDBContextMockGetBadgerDB) Set(f func() (r *badger.DB)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetBadgerDBFunc = f
	return m.mock
}

//GetBadgerDB implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) GetBadgerDB() (r *badger.DB) {
	counter := atomic.AddUint64(&m.GetBadgerDBPreCounter, 1)
	defer atomic.AddUint64(&m.GetBadgerDBCounter, 1)

	if len(m.GetBadgerDBMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetBadgerDBMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.GetBadgerDB.")
			return
		}

		result := m.GetBadgerDBMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.GetBadgerDB")
			return
		}

		r = result.r

		return
	}

	if m.GetBadgerDBMock.mainExpectation != nil {

		result := m.GetBadgerDBMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.GetBadgerDB")
		}

		r = result.r

		return
	}

	if m.GetBadgerDBFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.GetBadgerDB.")
		return
	}

	return m.GetBadgerDBFunc()
}

//GetBadgerDBMinimockCounter returns a count of DBContextMock.GetBadgerDBFunc invocations
func (m *DBContextMock) GetBadgerDBMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetBadgerDBCounter)
}

//GetBadgerDBMinimockPreCounter returns the value of DBContextMock.GetBadgerDB invocations
func (m *DBContextMock) GetBadgerDBMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetBadgerDBPreCounter)
}

//GetBadgerDBFinished returns true if mock invocations count is ok
func (m *DBContextMock) GetBadgerDBFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetBadgerDBMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetBadgerDBCounter) == uint64(len(m.GetBadgerDBMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetBadgerDBMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetBadgerDBCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetBadgerDBFunc != nil {
		return atomic.LoadUint64(&m.GetBadgerDBCounter) > 0
	}

	return true
}

type mDBContextMockIterateRecordsOnPulse struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockIterateRecordsOnPulseExpectation
	expectationSeries []*DBContextMockIterateRecordsOnPulseExpectation
}

type DBContextMockIterateRecordsOnPulseExpectation struct {
	input  *DBContextMockIterateRecordsOnPulseInput
	result *DBContextMockIterateRecordsOnPulseResult
}

type DBContextMockIterateRecordsOnPulseInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 func(p core.RecordID, p1 record.Record) (r error)
}

type DBContextMockIterateRecordsOnPulseResult struct {
	r error
}

//Expect specifies that invocation of DBContext.IterateRecordsOnPulse is expected from 1 to Infinity times
func (m *mDBContextMockIterateRecordsOnPulse) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 func(p core.RecordID, p1 record.Record) (r error)) *mDBContextMockIterateRecordsOnPulse {
	m.mock.IterateRecordsOnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockIterateRecordsOnPulseExpectation{}
	}
	m.mainExpectation.input = &DBContextMockIterateRecordsOnPulseInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of DBContext.IterateRecordsOnPulse
func (m *mDBContextMockIterateRecordsOnPulse) Return(r error) *DBContextMock {
	m.mock.IterateRecordsOnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockIterateRecordsOnPulseExpectation{}
	}
	m.mainExpectation.result = &DBContextMockIterateRecordsOnPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.IterateRecordsOnPulse is expected once
func (m *mDBContextMockIterateRecordsOnPulse) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 func(p core.RecordID, p1 record.Record) (r error)) *DBContextMockIterateRecordsOnPulseExpectation {
	m.mock.IterateRecordsOnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockIterateRecordsOnPulseExpectation{}
	expectation.input = &DBContextMockIterateRecordsOnPulseInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockIterateRecordsOnPulseExpectation) Return(r error) {
	e.result = &DBContextMockIterateRecordsOnPulseResult{r}
}

//Set uses given function f as a mock of DBContext.IterateRecordsOnPulse method
func (m *mDBContextMockIterateRecordsOnPulse) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 func(p core.RecordID, p1 record.Record) (r error)) (r error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IterateRecordsOnPulseFunc = f
	return m.mock
}

//IterateRecordsOnPulse implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) IterateRecordsOnPulse(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 func(p core.RecordID, p1 record.Record) (r error)) (r error) {
	counter := atomic.AddUint64(&m.IterateRecordsOnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.IterateRecordsOnPulseCounter, 1)

	if len(m.IterateRecordsOnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IterateRecordsOnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.IterateRecordsOnPulse. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.IterateRecordsOnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMockIterateRecordsOnPulseInput{p, p1, p2, p3}, "DBContext.IterateRecordsOnPulse got unexpected parameters")

		result := m.IterateRecordsOnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.IterateRecordsOnPulse")
			return
		}

		r = result.r

		return
	}

	if m.IterateRecordsOnPulseMock.mainExpectation != nil {

		input := m.IterateRecordsOnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMockIterateRecordsOnPulseInput{p, p1, p2, p3}, "DBContext.IterateRecordsOnPulse got unexpected parameters")
		}

		result := m.IterateRecordsOnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.IterateRecordsOnPulse")
		}

		r = result.r

		return
	}

	if m.IterateRecordsOnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.IterateRecordsOnPulse. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.IterateRecordsOnPulseFunc(p, p1, p2, p3)
}

//IterateRecordsOnPulseMinimockCounter returns a count of DBContextMock.IterateRecordsOnPulseFunc invocations
func (m *DBContextMock) IterateRecordsOnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IterateRecordsOnPulseCounter)
}

//IterateRecordsOnPulseMinimockPreCounter returns the value of DBContextMock.IterateRecordsOnPulse invocations
func (m *DBContextMock) IterateRecordsOnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IterateRecordsOnPulsePreCounter)
}

//IterateRecordsOnPulseFinished returns true if mock invocations count is ok
func (m *DBContextMock) IterateRecordsOnPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IterateRecordsOnPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IterateRecordsOnPulseCounter) == uint64(len(m.IterateRecordsOnPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IterateRecordsOnPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IterateRecordsOnPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IterateRecordsOnPulseFunc != nil {
		return atomic.LoadUint64(&m.IterateRecordsOnPulseCounter) > 0
	}

	return true
}

type mDBContextMockStoreKeyValues struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockStoreKeyValuesExpectation
	expectationSeries []*DBContextMockStoreKeyValuesExpectation
}

type DBContextMockStoreKeyValuesExpectation struct {
	input  *DBContextMockStoreKeyValuesInput
	result *DBContextMockStoreKeyValuesResult
}

type DBContextMockStoreKeyValuesInput struct {
	p  context.Context
	p1 []core.KV
}

type DBContextMockStoreKeyValuesResult struct {
	r error
}

//Expect specifies that invocation of DBContext.StoreKeyValues is expected from 1 to Infinity times
func (m *mDBContextMockStoreKeyValues) Expect(p context.Context, p1 []core.KV) *mDBContextMockStoreKeyValues {
	m.mock.StoreKeyValuesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockStoreKeyValuesExpectation{}
	}
	m.mainExpectation.input = &DBContextMockStoreKeyValuesInput{p, p1}
	return m
}

//Return specifies results of invocation of DBContext.StoreKeyValues
func (m *mDBContextMockStoreKeyValues) Return(r error) *DBContextMock {
	m.mock.StoreKeyValuesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockStoreKeyValuesExpectation{}
	}
	m.mainExpectation.result = &DBContextMockStoreKeyValuesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.StoreKeyValues is expected once
func (m *mDBContextMockStoreKeyValues) ExpectOnce(p context.Context, p1 []core.KV) *DBContextMockStoreKeyValuesExpectation {
	m.mock.StoreKeyValuesFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockStoreKeyValuesExpectation{}
	expectation.input = &DBContextMockStoreKeyValuesInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockStoreKeyValuesExpectation) Return(r error) {
	e.result = &DBContextMockStoreKeyValuesResult{r}
}

//Set uses given function f as a mock of DBContext.StoreKeyValues method
func (m *mDBContextMockStoreKeyValues) Set(f func(p context.Context, p1 []core.KV) (r error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreKeyValuesFunc = f
	return m.mock
}

//StoreKeyValues implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) StoreKeyValues(p context.Context, p1 []core.KV) (r error) {
	counter := atomic.AddUint64(&m.StoreKeyValuesPreCounter, 1)
	defer atomic.AddUint64(&m.StoreKeyValuesCounter, 1)

	if len(m.StoreKeyValuesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreKeyValuesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.StoreKeyValues. %v %v", p, p1)
			return
		}

		input := m.StoreKeyValuesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMockStoreKeyValuesInput{p, p1}, "DBContext.StoreKeyValues got unexpected parameters")

		result := m.StoreKeyValuesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.StoreKeyValues")
			return
		}

		r = result.r

		return
	}

	if m.StoreKeyValuesMock.mainExpectation != nil {

		input := m.StoreKeyValuesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMockStoreKeyValuesInput{p, p1}, "DBContext.StoreKeyValues got unexpected parameters")
		}

		result := m.StoreKeyValuesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.StoreKeyValues")
		}

		r = result.r

		return
	}

	if m.StoreKeyValuesFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.StoreKeyValues. %v %v", p, p1)
		return
	}

	return m.StoreKeyValuesFunc(p, p1)
}

//StoreKeyValuesMinimockCounter returns a count of DBContextMock.StoreKeyValuesFunc invocations
func (m *DBContextMock) StoreKeyValuesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreKeyValuesCounter)
}

//StoreKeyValuesMinimockPreCounter returns the value of DBContextMock.StoreKeyValues invocations
func (m *DBContextMock) StoreKeyValuesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StoreKeyValuesPreCounter)
}

//StoreKeyValuesFinished returns true if mock invocations count is ok
func (m *DBContextMock) StoreKeyValuesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreKeyValuesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreKeyValuesCounter) == uint64(len(m.StoreKeyValuesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreKeyValuesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreKeyValuesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreKeyValuesFunc != nil {
		return atomic.LoadUint64(&m.StoreKeyValuesCounter) > 0
	}

	return true
}

type mDBContextMockUpdate struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockUpdateExpectation
	expectationSeries []*DBContextMockUpdateExpectation
}

type DBContextMockUpdateExpectation struct {
	input  *DBContextMockUpdateInput
	result *DBContextMockUpdateResult
}

type DBContextMockUpdateInput struct {
	p  context.Context
	p1 func(p *TransactionManager) (r error)
}

type DBContextMockUpdateResult struct {
	r error
}

//Expect specifies that invocation of DBContext.Update is expected from 1 to Infinity times
func (m *mDBContextMockUpdate) Expect(p context.Context, p1 func(p *TransactionManager) (r error)) *mDBContextMockUpdate {
	m.mock.UpdateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockUpdateExpectation{}
	}
	m.mainExpectation.input = &DBContextMockUpdateInput{p, p1}
	return m
}

//Return specifies results of invocation of DBContext.Update
func (m *mDBContextMockUpdate) Return(r error) *DBContextMock {
	m.mock.UpdateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockUpdateExpectation{}
	}
	m.mainExpectation.result = &DBContextMockUpdateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.Update is expected once
func (m *mDBContextMockUpdate) ExpectOnce(p context.Context, p1 func(p *TransactionManager) (r error)) *DBContextMockUpdateExpectation {
	m.mock.UpdateFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockUpdateExpectation{}
	expectation.input = &DBContextMockUpdateInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockUpdateExpectation) Return(r error) {
	e.result = &DBContextMockUpdateResult{r}
}

//Set uses given function f as a mock of DBContext.Update method
func (m *mDBContextMockUpdate) Set(f func(p context.Context, p1 func(p *TransactionManager) (r error)) (r error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateFunc = f
	return m.mock
}

//Update implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) Update(p context.Context, p1 func(p *TransactionManager) (r error)) (r error) {
	counter := atomic.AddUint64(&m.UpdatePreCounter, 1)
	defer atomic.AddUint64(&m.UpdateCounter, 1)

	if len(m.UpdateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.Update. %v %v", p, p1)
			return
		}

		input := m.UpdateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMockUpdateInput{p, p1}, "DBContext.Update got unexpected parameters")

		result := m.UpdateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.Update")
			return
		}

		r = result.r

		return
	}

	if m.UpdateMock.mainExpectation != nil {

		input := m.UpdateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMockUpdateInput{p, p1}, "DBContext.Update got unexpected parameters")
		}

		result := m.UpdateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.Update")
		}

		r = result.r

		return
	}

	if m.UpdateFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.Update. %v %v", p, p1)
		return
	}

	return m.UpdateFunc(p, p1)
}

//UpdateMinimockCounter returns a count of DBContextMock.UpdateFunc invocations
func (m *DBContextMock) UpdateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateCounter)
}

//UpdateMinimockPreCounter returns the value of DBContextMock.Update invocations
func (m *DBContextMock) UpdateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdatePreCounter)
}

//UpdateFinished returns true if mock invocations count is ok
func (m *DBContextMock) UpdateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpdateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpdateCounter) == uint64(len(m.UpdateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpdateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpdateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpdateFunc != nil {
		return atomic.LoadUint64(&m.UpdateCounter) > 0
	}

	return true
}

type mDBContextMockView struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockViewExpectation
	expectationSeries []*DBContextMockViewExpectation
}

type DBContextMockViewExpectation struct {
	input  *DBContextMockViewInput
	result *DBContextMockViewResult
}

type DBContextMockViewInput struct {
	p  context.Context
	p1 func(p *TransactionManager) (r error)
}

type DBContextMockViewResult struct {
	r error
}

//Expect specifies that invocation of DBContext.View is expected from 1 to Infinity times
func (m *mDBContextMockView) Expect(p context.Context, p1 func(p *TransactionManager) (r error)) *mDBContextMockView {
	m.mock.ViewFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockViewExpectation{}
	}
	m.mainExpectation.input = &DBContextMockViewInput{p, p1}
	return m
}

//Return specifies results of invocation of DBContext.View
func (m *mDBContextMockView) Return(r error) *DBContextMock {
	m.mock.ViewFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockViewExpectation{}
	}
	m.mainExpectation.result = &DBContextMockViewResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.View is expected once
func (m *mDBContextMockView) ExpectOnce(p context.Context, p1 func(p *TransactionManager) (r error)) *DBContextMockViewExpectation {
	m.mock.ViewFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockViewExpectation{}
	expectation.input = &DBContextMockViewInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockViewExpectation) Return(r error) {
	e.result = &DBContextMockViewResult{r}
}

//Set uses given function f as a mock of DBContext.View method
func (m *mDBContextMockView) Set(f func(p context.Context, p1 func(p *TransactionManager) (r error)) (r error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ViewFunc = f
	return m.mock
}

//View implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) View(p context.Context, p1 func(p *TransactionManager) (r error)) (r error) {
	counter := atomic.AddUint64(&m.ViewPreCounter, 1)
	defer atomic.AddUint64(&m.ViewCounter, 1)

	if len(m.ViewMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ViewMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.View. %v %v", p, p1)
			return
		}

		input := m.ViewMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMockViewInput{p, p1}, "DBContext.View got unexpected parameters")

		result := m.ViewMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.View")
			return
		}

		r = result.r

		return
	}

	if m.ViewMock.mainExpectation != nil {

		input := m.ViewMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMockViewInput{p, p1}, "DBContext.View got unexpected parameters")
		}

		result := m.ViewMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.View")
		}

		r = result.r

		return
	}

	if m.ViewFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.View. %v %v", p, p1)
		return
	}

	return m.ViewFunc(p, p1)
}

//ViewMinimockCounter returns a count of DBContextMock.ViewFunc invocations
func (m *DBContextMock) ViewMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ViewCounter)
}

//ViewMinimockPreCounter returns the value of DBContextMock.View invocations
func (m *DBContextMock) ViewMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ViewPreCounter)
}

//ViewFinished returns true if mock invocations count is ok
func (m *DBContextMock) ViewFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ViewMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ViewCounter) == uint64(len(m.ViewMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ViewMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ViewCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ViewFunc != nil {
		return atomic.LoadUint64(&m.ViewCounter) > 0
	}

	return true
}

type mDBContextMockget struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockgetExpectation
	expectationSeries []*DBContextMockgetExpectation
}

type DBContextMockgetExpectation struct {
	input  *DBContextMockgetInput
	result *DBContextMockgetResult
}

type DBContextMockgetInput struct {
	p  context.Context
	p1 []byte
}

type DBContextMockgetResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of DBContext.get is expected from 1 to Infinity times
func (m *mDBContextMockget) Expect(p context.Context, p1 []byte) *mDBContextMockget {
	m.mock.getFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockgetExpectation{}
	}
	m.mainExpectation.input = &DBContextMockgetInput{p, p1}
	return m
}

//Return specifies results of invocation of DBContext.get
func (m *mDBContextMockget) Return(r []byte, r1 error) *DBContextMock {
	m.mock.getFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockgetExpectation{}
	}
	m.mainExpectation.result = &DBContextMockgetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.get is expected once
func (m *mDBContextMockget) ExpectOnce(p context.Context, p1 []byte) *DBContextMockgetExpectation {
	m.mock.getFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockgetExpectation{}
	expectation.input = &DBContextMockgetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockgetExpectation) Return(r []byte, r1 error) {
	e.result = &DBContextMockgetResult{r, r1}
}

//Set uses given function f as a mock of DBContext.get method
func (m *mDBContextMockget) Set(f func(p context.Context, p1 []byte) (r []byte, r1 error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.getFunc = f
	return m.mock
}

//get implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) get(p context.Context, p1 []byte) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.getPreCounter, 1)
	defer atomic.AddUint64(&m.getCounter, 1)

	if len(m.getMock.expectationSeries) > 0 {
		if counter > uint64(len(m.getMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.get. %v %v", p, p1)
			return
		}

		input := m.getMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMockgetInput{p, p1}, "DBContext.get got unexpected parameters")

		result := m.getMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.get")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.getMock.mainExpectation != nil {

		input := m.getMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMockgetInput{p, p1}, "DBContext.get got unexpected parameters")
		}

		result := m.getMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.get")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.getFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.get. %v %v", p, p1)
		return
	}

	return m.getFunc(p, p1)
}

//getMinimockCounter returns a count of DBContextMock.getFunc invocations
func (m *DBContextMock) getMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.getCounter)
}

//getMinimockPreCounter returns the value of DBContextMock.get invocations
func (m *DBContextMock) getMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.getPreCounter)
}

//getFinished returns true if mock invocations count is ok
func (m *DBContextMock) getFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.getMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.getCounter) == uint64(len(m.getMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.getMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.getCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.getFunc != nil {
		return atomic.LoadUint64(&m.getCounter) > 0
	}

	return true
}

type mDBContextMockiterate struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockiterateExpectation
	expectationSeries []*DBContextMockiterateExpectation
}

type DBContextMockiterateExpectation struct {
	input  *DBContextMockiterateInput
	result *DBContextMockiterateResult
}

type DBContextMockiterateInput struct {
	p  context.Context
	p1 []byte
	p2 func(p []byte, p1 []byte) (r error)
}

type DBContextMockiterateResult struct {
	r error
}

//Expect specifies that invocation of DBContext.iterate is expected from 1 to Infinity times
func (m *mDBContextMockiterate) Expect(p context.Context, p1 []byte, p2 func(p []byte, p1 []byte) (r error)) *mDBContextMockiterate {
	m.mock.iterateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockiterateExpectation{}
	}
	m.mainExpectation.input = &DBContextMockiterateInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DBContext.iterate
func (m *mDBContextMockiterate) Return(r error) *DBContextMock {
	m.mock.iterateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockiterateExpectation{}
	}
	m.mainExpectation.result = &DBContextMockiterateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.iterate is expected once
func (m *mDBContextMockiterate) ExpectOnce(p context.Context, p1 []byte, p2 func(p []byte, p1 []byte) (r error)) *DBContextMockiterateExpectation {
	m.mock.iterateFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockiterateExpectation{}
	expectation.input = &DBContextMockiterateInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMockiterateExpectation) Return(r error) {
	e.result = &DBContextMockiterateResult{r}
}

//Set uses given function f as a mock of DBContext.iterate method
func (m *mDBContextMockiterate) Set(f func(p context.Context, p1 []byte, p2 func(p []byte, p1 []byte) (r error)) (r error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.iterateFunc = f
	return m.mock
}

//iterate implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) iterate(p context.Context, p1 []byte, p2 func(p []byte, p1 []byte) (r error)) (r error) {
	counter := atomic.AddUint64(&m.iteratePreCounter, 1)
	defer atomic.AddUint64(&m.iterateCounter, 1)

	if len(m.iterateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.iterateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.iterate. %v %v %v", p, p1, p2)
			return
		}

		input := m.iterateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMockiterateInput{p, p1, p2}, "DBContext.iterate got unexpected parameters")

		result := m.iterateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.iterate")
			return
		}

		r = result.r

		return
	}

	if m.iterateMock.mainExpectation != nil {

		input := m.iterateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMockiterateInput{p, p1, p2}, "DBContext.iterate got unexpected parameters")
		}

		result := m.iterateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.iterate")
		}

		r = result.r

		return
	}

	if m.iterateFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.iterate. %v %v %v", p, p1, p2)
		return
	}

	return m.iterateFunc(p, p1, p2)
}

//iterateMinimockCounter returns a count of DBContextMock.iterateFunc invocations
func (m *DBContextMock) iterateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.iterateCounter)
}

//iterateMinimockPreCounter returns the value of DBContextMock.iterate invocations
func (m *DBContextMock) iterateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.iteratePreCounter)
}

//iterateFinished returns true if mock invocations count is ok
func (m *DBContextMock) iterateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.iterateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.iterateCounter) == uint64(len(m.iterateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.iterateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.iterateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.iterateFunc != nil {
		return atomic.LoadUint64(&m.iterateCounter) > 0
	}

	return true
}

type mDBContextMockset struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMocksetExpectation
	expectationSeries []*DBContextMocksetExpectation
}

type DBContextMocksetExpectation struct {
	input  *DBContextMocksetInput
	result *DBContextMocksetResult
}

type DBContextMocksetInput struct {
	p  context.Context
	p1 []byte
	p2 []byte
}

type DBContextMocksetResult struct {
	r error
}

//Expect specifies that invocation of DBContext.set is expected from 1 to Infinity times
func (m *mDBContextMockset) Expect(p context.Context, p1 []byte, p2 []byte) *mDBContextMockset {
	m.mock.setFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMocksetExpectation{}
	}
	m.mainExpectation.input = &DBContextMocksetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DBContext.set
func (m *mDBContextMockset) Return(r error) *DBContextMock {
	m.mock.setFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMocksetExpectation{}
	}
	m.mainExpectation.result = &DBContextMocksetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.set is expected once
func (m *mDBContextMockset) ExpectOnce(p context.Context, p1 []byte, p2 []byte) *DBContextMocksetExpectation {
	m.mock.setFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMocksetExpectation{}
	expectation.input = &DBContextMocksetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DBContextMocksetExpectation) Return(r error) {
	e.result = &DBContextMocksetResult{r}
}

//Set uses given function f as a mock of DBContext.set method
func (m *mDBContextMockset) Set(f func(p context.Context, p1 []byte, p2 []byte) (r error)) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.setFunc = f
	return m.mock
}

//set implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) set(p context.Context, p1 []byte, p2 []byte) (r error) {
	counter := atomic.AddUint64(&m.setPreCounter, 1)
	defer atomic.AddUint64(&m.setCounter, 1)

	if len(m.setMock.expectationSeries) > 0 {
		if counter > uint64(len(m.setMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.set. %v %v %v", p, p1, p2)
			return
		}

		input := m.setMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DBContextMocksetInput{p, p1, p2}, "DBContext.set got unexpected parameters")

		result := m.setMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.set")
			return
		}

		r = result.r

		return
	}

	if m.setMock.mainExpectation != nil {

		input := m.setMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DBContextMocksetInput{p, p1, p2}, "DBContext.set got unexpected parameters")
		}

		result := m.setMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DBContextMock.set")
		}

		r = result.r

		return
	}

	if m.setFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.set. %v %v %v", p, p1, p2)
		return
	}

	return m.setFunc(p, p1, p2)
}

//setMinimockCounter returns a count of DBContextMock.setFunc invocations
func (m *DBContextMock) setMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.setCounter)
}

//setMinimockPreCounter returns the value of DBContextMock.set invocations
func (m *DBContextMock) setMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.setPreCounter)
}

//setFinished returns true if mock invocations count is ok
func (m *DBContextMock) setFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.setMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.setCounter) == uint64(len(m.setMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.setMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.setCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.setFunc != nil {
		return atomic.LoadUint64(&m.setCounter) > 0
	}

	return true
}

type mDBContextMockwaitingFlight struct {
	mock              *DBContextMock
	mainExpectation   *DBContextMockwaitingFlightExpectation
	expectationSeries []*DBContextMockwaitingFlightExpectation
}

type DBContextMockwaitingFlightExpectation struct {
}

//Expect specifies that invocation of DBContext.waitingFlight is expected from 1 to Infinity times
func (m *mDBContextMockwaitingFlight) Expect() *mDBContextMockwaitingFlight {
	m.mock.waitingFlightFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockwaitingFlightExpectation{}
	}

	return m
}

//Return specifies results of invocation of DBContext.waitingFlight
func (m *mDBContextMockwaitingFlight) Return() *DBContextMock {
	m.mock.waitingFlightFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DBContextMockwaitingFlightExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of DBContext.waitingFlight is expected once
func (m *mDBContextMockwaitingFlight) ExpectOnce() *DBContextMockwaitingFlightExpectation {
	m.mock.waitingFlightFunc = nil
	m.mainExpectation = nil

	expectation := &DBContextMockwaitingFlightExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of DBContext.waitingFlight method
func (m *mDBContextMockwaitingFlight) Set(f func()) *DBContextMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.waitingFlightFunc = f
	return m.mock
}

//waitingFlight implements github.com/insolar/insolar/ledger/storage.DBContext interface
func (m *DBContextMock) waitingFlight() {
	counter := atomic.AddUint64(&m.waitingFlightPreCounter, 1)
	defer atomic.AddUint64(&m.waitingFlightCounter, 1)

	if len(m.waitingFlightMock.expectationSeries) > 0 {
		if counter > uint64(len(m.waitingFlightMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DBContextMock.waitingFlight.")
			return
		}

		return
	}

	if m.waitingFlightMock.mainExpectation != nil {

		return
	}

	if m.waitingFlightFunc == nil {
		m.t.Fatalf("Unexpected call to DBContextMock.waitingFlight.")
		return
	}

	m.waitingFlightFunc()
}

//waitingFlightMinimockCounter returns a count of DBContextMock.waitingFlightFunc invocations
func (m *DBContextMock) waitingFlightMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.waitingFlightCounter)
}

//waitingFlightMinimockPreCounter returns the value of DBContextMock.waitingFlight invocations
func (m *DBContextMock) waitingFlightMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.waitingFlightPreCounter)
}

//waitingFlightFinished returns true if mock invocations count is ok
func (m *DBContextMock) waitingFlightFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.waitingFlightMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.waitingFlightCounter) == uint64(len(m.waitingFlightMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.waitingFlightMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.waitingFlightCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.waitingFlightFunc != nil {
		return atomic.LoadUint64(&m.waitingFlightCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DBContextMock) ValidateCallCounters() {

	if !m.BeginTransactionFinished() {
		m.t.Fatal("Expected call to DBContextMock.BeginTransaction")
	}

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to DBContextMock.Close")
	}

	if !m.GetBadgerDBFinished() {
		m.t.Fatal("Expected call to DBContextMock.GetBadgerDB")
	}

	if !m.IterateRecordsOnPulseFinished() {
		m.t.Fatal("Expected call to DBContextMock.IterateRecordsOnPulse")
	}

	if !m.StoreKeyValuesFinished() {
		m.t.Fatal("Expected call to DBContextMock.StoreKeyValues")
	}

	if !m.UpdateFinished() {
		m.t.Fatal("Expected call to DBContextMock.Update")
	}

	if !m.ViewFinished() {
		m.t.Fatal("Expected call to DBContextMock.View")
	}

	if !m.getFinished() {
		m.t.Fatal("Expected call to DBContextMock.get")
	}

	if !m.iterateFinished() {
		m.t.Fatal("Expected call to DBContextMock.iterate")
	}

	if !m.setFinished() {
		m.t.Fatal("Expected call to DBContextMock.set")
	}

	if !m.waitingFlightFinished() {
		m.t.Fatal("Expected call to DBContextMock.waitingFlight")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DBContextMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DBContextMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DBContextMock) MinimockFinish() {

	if !m.BeginTransactionFinished() {
		m.t.Fatal("Expected call to DBContextMock.BeginTransaction")
	}

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to DBContextMock.Close")
	}

	if !m.GetBadgerDBFinished() {
		m.t.Fatal("Expected call to DBContextMock.GetBadgerDB")
	}

	if !m.IterateRecordsOnPulseFinished() {
		m.t.Fatal("Expected call to DBContextMock.IterateRecordsOnPulse")
	}

	if !m.StoreKeyValuesFinished() {
		m.t.Fatal("Expected call to DBContextMock.StoreKeyValues")
	}

	if !m.UpdateFinished() {
		m.t.Fatal("Expected call to DBContextMock.Update")
	}

	if !m.ViewFinished() {
		m.t.Fatal("Expected call to DBContextMock.View")
	}

	if !m.getFinished() {
		m.t.Fatal("Expected call to DBContextMock.get")
	}

	if !m.iterateFinished() {
		m.t.Fatal("Expected call to DBContextMock.iterate")
	}

	if !m.setFinished() {
		m.t.Fatal("Expected call to DBContextMock.set")
	}

	if !m.waitingFlightFinished() {
		m.t.Fatal("Expected call to DBContextMock.waitingFlight")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DBContextMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DBContextMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.BeginTransactionFinished()
		ok = ok && m.CloseFinished()
		ok = ok && m.GetBadgerDBFinished()
		ok = ok && m.IterateRecordsOnPulseFinished()
		ok = ok && m.StoreKeyValuesFinished()
		ok = ok && m.UpdateFinished()
		ok = ok && m.ViewFinished()
		ok = ok && m.getFinished()
		ok = ok && m.iterateFinished()
		ok = ok && m.setFinished()
		ok = ok && m.waitingFlightFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.BeginTransactionFinished() {
				m.t.Error("Expected call to DBContextMock.BeginTransaction")
			}

			if !m.CloseFinished() {
				m.t.Error("Expected call to DBContextMock.Close")
			}

			if !m.GetBadgerDBFinished() {
				m.t.Error("Expected call to DBContextMock.GetBadgerDB")
			}

			if !m.IterateRecordsOnPulseFinished() {
				m.t.Error("Expected call to DBContextMock.IterateRecordsOnPulse")
			}

			if !m.StoreKeyValuesFinished() {
				m.t.Error("Expected call to DBContextMock.StoreKeyValues")
			}

			if !m.UpdateFinished() {
				m.t.Error("Expected call to DBContextMock.Update")
			}

			if !m.ViewFinished() {
				m.t.Error("Expected call to DBContextMock.View")
			}

			if !m.getFinished() {
				m.t.Error("Expected call to DBContextMock.get")
			}

			if !m.iterateFinished() {
				m.t.Error("Expected call to DBContextMock.iterate")
			}

			if !m.setFinished() {
				m.t.Error("Expected call to DBContextMock.set")
			}

			if !m.waitingFlightFinished() {
				m.t.Error("Expected call to DBContextMock.waitingFlight")
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
func (m *DBContextMock) AllMocksCalled() bool {

	if !m.BeginTransactionFinished() {
		return false
	}

	if !m.CloseFinished() {
		return false
	}

	if !m.GetBadgerDBFinished() {
		return false
	}

	if !m.IterateRecordsOnPulseFinished() {
		return false
	}

	if !m.StoreKeyValuesFinished() {
		return false
	}

	if !m.UpdateFinished() {
		return false
	}

	if !m.ViewFinished() {
		return false
	}

	if !m.getFinished() {
		return false
	}

	if !m.iterateFinished() {
		return false
	}

	if !m.setFinished() {
		return false
	}

	if !m.waitingFlightFinished() {
		return false
	}

	return true
}
