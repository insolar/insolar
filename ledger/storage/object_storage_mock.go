package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ObjectStorage" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	index "github.com/insolar/insolar/ledger/storage/index"
	record "github.com/insolar/insolar/ledger/storage/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//ObjectStorageMock implements github.com/insolar/insolar/ledger/storage.ObjectStorage
type ObjectStorageMock struct {
	t minimock.Tester

	GetBlobFunc       func(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r []byte, r1 error)
	GetBlobCounter    uint64
	GetBlobPreCounter uint64
	GetBlobMock       mObjectStorageMockGetBlob

	GetObjectIndexFunc       func(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 bool) (r *index.ObjectLifeline, r1 error)
	GetObjectIndexCounter    uint64
	GetObjectIndexPreCounter uint64
	GetObjectIndexMock       mObjectStorageMockGetObjectIndex

	GetRecordFunc       func(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r record.Record, r1 error)
	GetRecordCounter    uint64
	GetRecordPreCounter uint64
	GetRecordMock       mObjectStorageMockGetRecord

	IterateIndexIDsFunc       func(p context.Context, p1 core.RecordID, p2 func(p core.RecordID) (r error)) (r error)
	IterateIndexIDsCounter    uint64
	IterateIndexIDsPreCounter uint64
	IterateIndexIDsMock       mObjectStorageMockIterateIndexIDs

	RemoveObjectIndexFunc       func(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r error)
	RemoveObjectIndexCounter    uint64
	RemoveObjectIndexPreCounter uint64
	RemoveObjectIndexMock       mObjectStorageMockRemoveObjectIndex

	SetBlobFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *core.RecordID, r1 error)
	SetBlobCounter    uint64
	SetBlobPreCounter uint64
	SetBlobMock       mObjectStorageMockSetBlob

	SetObjectIndexFunc       func(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 *index.ObjectLifeline) (r error)
	SetObjectIndexCounter    uint64
	SetObjectIndexPreCounter uint64
	SetObjectIndexMock       mObjectStorageMockSetObjectIndex

	SetRecordFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 record.Record) (r *core.RecordID, r1 error)
	SetRecordCounter    uint64
	SetRecordPreCounter uint64
	SetRecordMock       mObjectStorageMockSetRecord
}

//NewObjectStorageMock returns a mock for github.com/insolar/insolar/ledger/storage.ObjectStorage
func NewObjectStorageMock(t minimock.Tester) *ObjectStorageMock {
	m := &ObjectStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetBlobMock = mObjectStorageMockGetBlob{mock: m}
	m.GetObjectIndexMock = mObjectStorageMockGetObjectIndex{mock: m}
	m.GetRecordMock = mObjectStorageMockGetRecord{mock: m}
	m.IterateIndexIDsMock = mObjectStorageMockIterateIndexIDs{mock: m}
	m.RemoveObjectIndexMock = mObjectStorageMockRemoveObjectIndex{mock: m}
	m.SetBlobMock = mObjectStorageMockSetBlob{mock: m}
	m.SetObjectIndexMock = mObjectStorageMockSetObjectIndex{mock: m}
	m.SetRecordMock = mObjectStorageMockSetRecord{mock: m}

	return m
}

type mObjectStorageMockGetBlob struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockGetBlobExpectation
	expectationSeries []*ObjectStorageMockGetBlobExpectation
}

type ObjectStorageMockGetBlobExpectation struct {
	input  *ObjectStorageMockGetBlobInput
	result *ObjectStorageMockGetBlobResult
}

type ObjectStorageMockGetBlobInput struct {
	p  context.Context
	p1 core.RecordID
	p2 *core.RecordID
}

type ObjectStorageMockGetBlobResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of ObjectStorage.GetBlob is expected from 1 to Infinity times
func (m *mObjectStorageMockGetBlob) Expect(p context.Context, p1 core.RecordID, p2 *core.RecordID) *mObjectStorageMockGetBlob {
	m.mock.GetBlobFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetBlobExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockGetBlobInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ObjectStorage.GetBlob
func (m *mObjectStorageMockGetBlob) Return(r []byte, r1 error) *ObjectStorageMock {
	m.mock.GetBlobFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetBlobExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockGetBlobResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.GetBlob is expected once
func (m *mObjectStorageMockGetBlob) ExpectOnce(p context.Context, p1 core.RecordID, p2 *core.RecordID) *ObjectStorageMockGetBlobExpectation {
	m.mock.GetBlobFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockGetBlobExpectation{}
	expectation.input = &ObjectStorageMockGetBlobInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockGetBlobExpectation) Return(r []byte, r1 error) {
	e.result = &ObjectStorageMockGetBlobResult{r, r1}
}

//Set uses given function f as a mock of ObjectStorage.GetBlob method
func (m *mObjectStorageMockGetBlob) Set(f func(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r []byte, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetBlobFunc = f
	return m.mock
}

//GetBlob implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) GetBlob(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.GetBlobPreCounter, 1)
	defer atomic.AddUint64(&m.GetBlobCounter, 1)

	if len(m.GetBlobMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetBlobMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.GetBlob. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetBlobMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockGetBlobInput{p, p1, p2}, "ObjectStorage.GetBlob got unexpected parameters")

		result := m.GetBlobMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.GetBlob")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetBlobMock.mainExpectation != nil {

		input := m.GetBlobMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockGetBlobInput{p, p1, p2}, "ObjectStorage.GetBlob got unexpected parameters")
		}

		result := m.GetBlobMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.GetBlob")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetBlobFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.GetBlob. %v %v %v", p, p1, p2)
		return
	}

	return m.GetBlobFunc(p, p1, p2)
}

//GetBlobMinimockCounter returns a count of ObjectStorageMock.GetBlobFunc invocations
func (m *ObjectStorageMock) GetBlobMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetBlobCounter)
}

//GetBlobMinimockPreCounter returns the value of ObjectStorageMock.GetBlob invocations
func (m *ObjectStorageMock) GetBlobMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetBlobPreCounter)
}

//GetBlobFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) GetBlobFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetBlobMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetBlobCounter) == uint64(len(m.GetBlobMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetBlobMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetBlobCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetBlobFunc != nil {
		return atomic.LoadUint64(&m.GetBlobCounter) > 0
	}

	return true
}

type mObjectStorageMockGetObjectIndex struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockGetObjectIndexExpectation
	expectationSeries []*ObjectStorageMockGetObjectIndexExpectation
}

type ObjectStorageMockGetObjectIndexExpectation struct {
	input  *ObjectStorageMockGetObjectIndexInput
	result *ObjectStorageMockGetObjectIndexResult
}

type ObjectStorageMockGetObjectIndexInput struct {
	p  context.Context
	p1 core.RecordID
	p2 *core.RecordID
	p3 bool
}

type ObjectStorageMockGetObjectIndexResult struct {
	r  *index.ObjectLifeline
	r1 error
}

//Expect specifies that invocation of ObjectStorage.GetObjectIndex is expected from 1 to Infinity times
func (m *mObjectStorageMockGetObjectIndex) Expect(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 bool) *mObjectStorageMockGetObjectIndex {
	m.mock.GetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetObjectIndexExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockGetObjectIndexInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ObjectStorage.GetObjectIndex
func (m *mObjectStorageMockGetObjectIndex) Return(r *index.ObjectLifeline, r1 error) *ObjectStorageMock {
	m.mock.GetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetObjectIndexExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockGetObjectIndexResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.GetObjectIndex is expected once
func (m *mObjectStorageMockGetObjectIndex) ExpectOnce(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 bool) *ObjectStorageMockGetObjectIndexExpectation {
	m.mock.GetObjectIndexFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockGetObjectIndexExpectation{}
	expectation.input = &ObjectStorageMockGetObjectIndexInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockGetObjectIndexExpectation) Return(r *index.ObjectLifeline, r1 error) {
	e.result = &ObjectStorageMockGetObjectIndexResult{r, r1}
}

//Set uses given function f as a mock of ObjectStorage.GetObjectIndex method
func (m *mObjectStorageMockGetObjectIndex) Set(f func(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 bool) (r *index.ObjectLifeline, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectIndexFunc = f
	return m.mock
}

//GetObjectIndex implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) GetObjectIndex(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 bool) (r *index.ObjectLifeline, r1 error) {
	counter := atomic.AddUint64(&m.GetObjectIndexPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectIndexCounter, 1)

	if len(m.GetObjectIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.GetObjectIndex. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.GetObjectIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockGetObjectIndexInput{p, p1, p2, p3}, "ObjectStorage.GetObjectIndex got unexpected parameters")

		result := m.GetObjectIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.GetObjectIndex")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetObjectIndexMock.mainExpectation != nil {

		input := m.GetObjectIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockGetObjectIndexInput{p, p1, p2, p3}, "ObjectStorage.GetObjectIndex got unexpected parameters")
		}

		result := m.GetObjectIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.GetObjectIndex")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetObjectIndexFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.GetObjectIndex. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.GetObjectIndexFunc(p, p1, p2, p3)
}

//GetObjectIndexMinimockCounter returns a count of ObjectStorageMock.GetObjectIndexFunc invocations
func (m *ObjectStorageMock) GetObjectIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectIndexCounter)
}

//GetObjectIndexMinimockPreCounter returns the value of ObjectStorageMock.GetObjectIndex invocations
func (m *ObjectStorageMock) GetObjectIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectIndexPreCounter)
}

//GetObjectIndexFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) GetObjectIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetObjectIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetObjectIndexCounter) == uint64(len(m.GetObjectIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetObjectIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetObjectIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetObjectIndexFunc != nil {
		return atomic.LoadUint64(&m.GetObjectIndexCounter) > 0
	}

	return true
}

type mObjectStorageMockGetRecord struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockGetRecordExpectation
	expectationSeries []*ObjectStorageMockGetRecordExpectation
}

type ObjectStorageMockGetRecordExpectation struct {
	input  *ObjectStorageMockGetRecordInput
	result *ObjectStorageMockGetRecordResult
}

type ObjectStorageMockGetRecordInput struct {
	p  context.Context
	p1 core.RecordID
	p2 *core.RecordID
}

type ObjectStorageMockGetRecordResult struct {
	r  record.Record
	r1 error
}

//Expect specifies that invocation of ObjectStorage.GetRecord is expected from 1 to Infinity times
func (m *mObjectStorageMockGetRecord) Expect(p context.Context, p1 core.RecordID, p2 *core.RecordID) *mObjectStorageMockGetRecord {
	m.mock.GetRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetRecordExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockGetRecordInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ObjectStorage.GetRecord
func (m *mObjectStorageMockGetRecord) Return(r record.Record, r1 error) *ObjectStorageMock {
	m.mock.GetRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetRecordExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockGetRecordResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.GetRecord is expected once
func (m *mObjectStorageMockGetRecord) ExpectOnce(p context.Context, p1 core.RecordID, p2 *core.RecordID) *ObjectStorageMockGetRecordExpectation {
	m.mock.GetRecordFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockGetRecordExpectation{}
	expectation.input = &ObjectStorageMockGetRecordInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockGetRecordExpectation) Return(r record.Record, r1 error) {
	e.result = &ObjectStorageMockGetRecordResult{r, r1}
}

//Set uses given function f as a mock of ObjectStorage.GetRecord method
func (m *mObjectStorageMockGetRecord) Set(f func(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r record.Record, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRecordFunc = f
	return m.mock
}

//GetRecord implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) GetRecord(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r record.Record, r1 error) {
	counter := atomic.AddUint64(&m.GetRecordPreCounter, 1)
	defer atomic.AddUint64(&m.GetRecordCounter, 1)

	if len(m.GetRecordMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRecordMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.GetRecord. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetRecordMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockGetRecordInput{p, p1, p2}, "ObjectStorage.GetRecord got unexpected parameters")

		result := m.GetRecordMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.GetRecord")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetRecordMock.mainExpectation != nil {

		input := m.GetRecordMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockGetRecordInput{p, p1, p2}, "ObjectStorage.GetRecord got unexpected parameters")
		}

		result := m.GetRecordMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.GetRecord")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetRecordFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.GetRecord. %v %v %v", p, p1, p2)
		return
	}

	return m.GetRecordFunc(p, p1, p2)
}

//GetRecordMinimockCounter returns a count of ObjectStorageMock.GetRecordFunc invocations
func (m *ObjectStorageMock) GetRecordMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRecordCounter)
}

//GetRecordMinimockPreCounter returns the value of ObjectStorageMock.GetRecord invocations
func (m *ObjectStorageMock) GetRecordMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRecordPreCounter)
}

//GetRecordFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) GetRecordFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRecordMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRecordCounter) == uint64(len(m.GetRecordMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRecordMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRecordCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRecordFunc != nil {
		return atomic.LoadUint64(&m.GetRecordCounter) > 0
	}

	return true
}

type mObjectStorageMockIterateIndexIDs struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockIterateIndexIDsExpectation
	expectationSeries []*ObjectStorageMockIterateIndexIDsExpectation
}

type ObjectStorageMockIterateIndexIDsExpectation struct {
	input  *ObjectStorageMockIterateIndexIDsInput
	result *ObjectStorageMockIterateIndexIDsResult
}

type ObjectStorageMockIterateIndexIDsInput struct {
	p  context.Context
	p1 core.RecordID
	p2 func(p core.RecordID) (r error)
}

type ObjectStorageMockIterateIndexIDsResult struct {
	r error
}

//Expect specifies that invocation of ObjectStorage.IterateIndexIDs is expected from 1 to Infinity times
func (m *mObjectStorageMockIterateIndexIDs) Expect(p context.Context, p1 core.RecordID, p2 func(p core.RecordID) (r error)) *mObjectStorageMockIterateIndexIDs {
	m.mock.IterateIndexIDsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockIterateIndexIDsExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockIterateIndexIDsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ObjectStorage.IterateIndexIDs
func (m *mObjectStorageMockIterateIndexIDs) Return(r error) *ObjectStorageMock {
	m.mock.IterateIndexIDsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockIterateIndexIDsExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockIterateIndexIDsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.IterateIndexIDs is expected once
func (m *mObjectStorageMockIterateIndexIDs) ExpectOnce(p context.Context, p1 core.RecordID, p2 func(p core.RecordID) (r error)) *ObjectStorageMockIterateIndexIDsExpectation {
	m.mock.IterateIndexIDsFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockIterateIndexIDsExpectation{}
	expectation.input = &ObjectStorageMockIterateIndexIDsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockIterateIndexIDsExpectation) Return(r error) {
	e.result = &ObjectStorageMockIterateIndexIDsResult{r}
}

//Set uses given function f as a mock of ObjectStorage.IterateIndexIDs method
func (m *mObjectStorageMockIterateIndexIDs) Set(f func(p context.Context, p1 core.RecordID, p2 func(p core.RecordID) (r error)) (r error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IterateIndexIDsFunc = f
	return m.mock
}

//IterateIndexIDs implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) IterateIndexIDs(p context.Context, p1 core.RecordID, p2 func(p core.RecordID) (r error)) (r error) {
	counter := atomic.AddUint64(&m.IterateIndexIDsPreCounter, 1)
	defer atomic.AddUint64(&m.IterateIndexIDsCounter, 1)

	if len(m.IterateIndexIDsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IterateIndexIDsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.IterateIndexIDs. %v %v %v", p, p1, p2)
			return
		}

		input := m.IterateIndexIDsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockIterateIndexIDsInput{p, p1, p2}, "ObjectStorage.IterateIndexIDs got unexpected parameters")

		result := m.IterateIndexIDsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.IterateIndexIDs")
			return
		}

		r = result.r

		return
	}

	if m.IterateIndexIDsMock.mainExpectation != nil {

		input := m.IterateIndexIDsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockIterateIndexIDsInput{p, p1, p2}, "ObjectStorage.IterateIndexIDs got unexpected parameters")
		}

		result := m.IterateIndexIDsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.IterateIndexIDs")
		}

		r = result.r

		return
	}

	if m.IterateIndexIDsFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.IterateIndexIDs. %v %v %v", p, p1, p2)
		return
	}

	return m.IterateIndexIDsFunc(p, p1, p2)
}

//IterateIndexIDsMinimockCounter returns a count of ObjectStorageMock.IterateIndexIDsFunc invocations
func (m *ObjectStorageMock) IterateIndexIDsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IterateIndexIDsCounter)
}

//IterateIndexIDsMinimockPreCounter returns the value of ObjectStorageMock.IterateIndexIDs invocations
func (m *ObjectStorageMock) IterateIndexIDsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IterateIndexIDsPreCounter)
}

//IterateIndexIDsFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) IterateIndexIDsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IterateIndexIDsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IterateIndexIDsCounter) == uint64(len(m.IterateIndexIDsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IterateIndexIDsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IterateIndexIDsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IterateIndexIDsFunc != nil {
		return atomic.LoadUint64(&m.IterateIndexIDsCounter) > 0
	}

	return true
}

type mObjectStorageMockRemoveObjectIndex struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockRemoveObjectIndexExpectation
	expectationSeries []*ObjectStorageMockRemoveObjectIndexExpectation
}

type ObjectStorageMockRemoveObjectIndexExpectation struct {
	input  *ObjectStorageMockRemoveObjectIndexInput
	result *ObjectStorageMockRemoveObjectIndexResult
}

type ObjectStorageMockRemoveObjectIndexInput struct {
	p  context.Context
	p1 core.RecordID
	p2 *core.RecordID
}

type ObjectStorageMockRemoveObjectIndexResult struct {
	r error
}

//Expect specifies that invocation of ObjectStorage.RemoveObjectIndex is expected from 1 to Infinity times
func (m *mObjectStorageMockRemoveObjectIndex) Expect(p context.Context, p1 core.RecordID, p2 *core.RecordID) *mObjectStorageMockRemoveObjectIndex {
	m.mock.RemoveObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockRemoveObjectIndexExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockRemoveObjectIndexInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ObjectStorage.RemoveObjectIndex
func (m *mObjectStorageMockRemoveObjectIndex) Return(r error) *ObjectStorageMock {
	m.mock.RemoveObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockRemoveObjectIndexExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockRemoveObjectIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.RemoveObjectIndex is expected once
func (m *mObjectStorageMockRemoveObjectIndex) ExpectOnce(p context.Context, p1 core.RecordID, p2 *core.RecordID) *ObjectStorageMockRemoveObjectIndexExpectation {
	m.mock.RemoveObjectIndexFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockRemoveObjectIndexExpectation{}
	expectation.input = &ObjectStorageMockRemoveObjectIndexInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockRemoveObjectIndexExpectation) Return(r error) {
	e.result = &ObjectStorageMockRemoveObjectIndexResult{r}
}

//Set uses given function f as a mock of ObjectStorage.RemoveObjectIndex method
func (m *mObjectStorageMockRemoveObjectIndex) Set(f func(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveObjectIndexFunc = f
	return m.mock
}

//RemoveObjectIndex implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) RemoveObjectIndex(p context.Context, p1 core.RecordID, p2 *core.RecordID) (r error) {
	counter := atomic.AddUint64(&m.RemoveObjectIndexPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveObjectIndexCounter, 1)

	if len(m.RemoveObjectIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveObjectIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.RemoveObjectIndex. %v %v %v", p, p1, p2)
			return
		}

		input := m.RemoveObjectIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockRemoveObjectIndexInput{p, p1, p2}, "ObjectStorage.RemoveObjectIndex got unexpected parameters")

		result := m.RemoveObjectIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.RemoveObjectIndex")
			return
		}

		r = result.r

		return
	}

	if m.RemoveObjectIndexMock.mainExpectation != nil {

		input := m.RemoveObjectIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockRemoveObjectIndexInput{p, p1, p2}, "ObjectStorage.RemoveObjectIndex got unexpected parameters")
		}

		result := m.RemoveObjectIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.RemoveObjectIndex")
		}

		r = result.r

		return
	}

	if m.RemoveObjectIndexFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.RemoveObjectIndex. %v %v %v", p, p1, p2)
		return
	}

	return m.RemoveObjectIndexFunc(p, p1, p2)
}

//RemoveObjectIndexMinimockCounter returns a count of ObjectStorageMock.RemoveObjectIndexFunc invocations
func (m *ObjectStorageMock) RemoveObjectIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveObjectIndexCounter)
}

//RemoveObjectIndexMinimockPreCounter returns the value of ObjectStorageMock.RemoveObjectIndex invocations
func (m *ObjectStorageMock) RemoveObjectIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveObjectIndexPreCounter)
}

//RemoveObjectIndexFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) RemoveObjectIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveObjectIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveObjectIndexCounter) == uint64(len(m.RemoveObjectIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveObjectIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveObjectIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveObjectIndexFunc != nil {
		return atomic.LoadUint64(&m.RemoveObjectIndexCounter) > 0
	}

	return true
}

type mObjectStorageMockSetBlob struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockSetBlobExpectation
	expectationSeries []*ObjectStorageMockSetBlobExpectation
}

type ObjectStorageMockSetBlobExpectation struct {
	input  *ObjectStorageMockSetBlobInput
	result *ObjectStorageMockSetBlobResult
}

type ObjectStorageMockSetBlobInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 []byte
}

type ObjectStorageMockSetBlobResult struct {
	r  *core.RecordID
	r1 error
}

//Expect specifies that invocation of ObjectStorage.SetBlob is expected from 1 to Infinity times
func (m *mObjectStorageMockSetBlob) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) *mObjectStorageMockSetBlob {
	m.mock.SetBlobFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetBlobExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockSetBlobInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ObjectStorage.SetBlob
func (m *mObjectStorageMockSetBlob) Return(r *core.RecordID, r1 error) *ObjectStorageMock {
	m.mock.SetBlobFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetBlobExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockSetBlobResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.SetBlob is expected once
func (m *mObjectStorageMockSetBlob) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) *ObjectStorageMockSetBlobExpectation {
	m.mock.SetBlobFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockSetBlobExpectation{}
	expectation.input = &ObjectStorageMockSetBlobInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockSetBlobExpectation) Return(r *core.RecordID, r1 error) {
	e.result = &ObjectStorageMockSetBlobResult{r, r1}
}

//Set uses given function f as a mock of ObjectStorage.SetBlob method
func (m *mObjectStorageMockSetBlob) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *core.RecordID, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetBlobFunc = f
	return m.mock
}

//SetBlob implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) SetBlob(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *core.RecordID, r1 error) {
	counter := atomic.AddUint64(&m.SetBlobPreCounter, 1)
	defer atomic.AddUint64(&m.SetBlobCounter, 1)

	if len(m.SetBlobMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetBlobMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.SetBlob. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetBlobMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockSetBlobInput{p, p1, p2, p3}, "ObjectStorage.SetBlob got unexpected parameters")

		result := m.SetBlobMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.SetBlob")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetBlobMock.mainExpectation != nil {

		input := m.SetBlobMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockSetBlobInput{p, p1, p2, p3}, "ObjectStorage.SetBlob got unexpected parameters")
		}

		result := m.SetBlobMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.SetBlob")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetBlobFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.SetBlob. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetBlobFunc(p, p1, p2, p3)
}

//SetBlobMinimockCounter returns a count of ObjectStorageMock.SetBlobFunc invocations
func (m *ObjectStorageMock) SetBlobMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetBlobCounter)
}

//SetBlobMinimockPreCounter returns the value of ObjectStorageMock.SetBlob invocations
func (m *ObjectStorageMock) SetBlobMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetBlobPreCounter)
}

//SetBlobFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) SetBlobFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetBlobMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetBlobCounter) == uint64(len(m.SetBlobMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetBlobMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetBlobCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetBlobFunc != nil {
		return atomic.LoadUint64(&m.SetBlobCounter) > 0
	}

	return true
}

type mObjectStorageMockSetObjectIndex struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockSetObjectIndexExpectation
	expectationSeries []*ObjectStorageMockSetObjectIndexExpectation
}

type ObjectStorageMockSetObjectIndexExpectation struct {
	input  *ObjectStorageMockSetObjectIndexInput
	result *ObjectStorageMockSetObjectIndexResult
}

type ObjectStorageMockSetObjectIndexInput struct {
	p  context.Context
	p1 core.RecordID
	p2 *core.RecordID
	p3 *index.ObjectLifeline
}

type ObjectStorageMockSetObjectIndexResult struct {
	r error
}

//Expect specifies that invocation of ObjectStorage.SetObjectIndex is expected from 1 to Infinity times
func (m *mObjectStorageMockSetObjectIndex) Expect(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 *index.ObjectLifeline) *mObjectStorageMockSetObjectIndex {
	m.mock.SetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetObjectIndexExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockSetObjectIndexInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ObjectStorage.SetObjectIndex
func (m *mObjectStorageMockSetObjectIndex) Return(r error) *ObjectStorageMock {
	m.mock.SetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetObjectIndexExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockSetObjectIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.SetObjectIndex is expected once
func (m *mObjectStorageMockSetObjectIndex) ExpectOnce(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 *index.ObjectLifeline) *ObjectStorageMockSetObjectIndexExpectation {
	m.mock.SetObjectIndexFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockSetObjectIndexExpectation{}
	expectation.input = &ObjectStorageMockSetObjectIndexInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockSetObjectIndexExpectation) Return(r error) {
	e.result = &ObjectStorageMockSetObjectIndexResult{r}
}

//Set uses given function f as a mock of ObjectStorage.SetObjectIndex method
func (m *mObjectStorageMockSetObjectIndex) Set(f func(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 *index.ObjectLifeline) (r error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetObjectIndexFunc = f
	return m.mock
}

//SetObjectIndex implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) SetObjectIndex(p context.Context, p1 core.RecordID, p2 *core.RecordID, p3 *index.ObjectLifeline) (r error) {
	counter := atomic.AddUint64(&m.SetObjectIndexPreCounter, 1)
	defer atomic.AddUint64(&m.SetObjectIndexCounter, 1)

	if len(m.SetObjectIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetObjectIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.SetObjectIndex. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetObjectIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockSetObjectIndexInput{p, p1, p2, p3}, "ObjectStorage.SetObjectIndex got unexpected parameters")

		result := m.SetObjectIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.SetObjectIndex")
			return
		}

		r = result.r

		return
	}

	if m.SetObjectIndexMock.mainExpectation != nil {

		input := m.SetObjectIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockSetObjectIndexInput{p, p1, p2, p3}, "ObjectStorage.SetObjectIndex got unexpected parameters")
		}

		result := m.SetObjectIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.SetObjectIndex")
		}

		r = result.r

		return
	}

	if m.SetObjectIndexFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.SetObjectIndex. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetObjectIndexFunc(p, p1, p2, p3)
}

//SetObjectIndexMinimockCounter returns a count of ObjectStorageMock.SetObjectIndexFunc invocations
func (m *ObjectStorageMock) SetObjectIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetObjectIndexCounter)
}

//SetObjectIndexMinimockPreCounter returns the value of ObjectStorageMock.SetObjectIndex invocations
func (m *ObjectStorageMock) SetObjectIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetObjectIndexPreCounter)
}

//SetObjectIndexFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) SetObjectIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetObjectIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetObjectIndexCounter) == uint64(len(m.SetObjectIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetObjectIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetObjectIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetObjectIndexFunc != nil {
		return atomic.LoadUint64(&m.SetObjectIndexCounter) > 0
	}

	return true
}

type mObjectStorageMockSetRecord struct {
	mock              *ObjectStorageMock
	mainExpectation   *ObjectStorageMockSetRecordExpectation
	expectationSeries []*ObjectStorageMockSetRecordExpectation
}

type ObjectStorageMockSetRecordExpectation struct {
	input  *ObjectStorageMockSetRecordInput
	result *ObjectStorageMockSetRecordResult
}

type ObjectStorageMockSetRecordInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 record.Record
}

type ObjectStorageMockSetRecordResult struct {
	r  *core.RecordID
	r1 error
}

//Expect specifies that invocation of ObjectStorage.SetRecord is expected from 1 to Infinity times
func (m *mObjectStorageMockSetRecord) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 record.Record) *mObjectStorageMockSetRecord {
	m.mock.SetRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetRecordExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockSetRecordInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ObjectStorage.SetRecord
func (m *mObjectStorageMockSetRecord) Return(r *core.RecordID, r1 error) *ObjectStorageMock {
	m.mock.SetRecordFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetRecordExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockSetRecordResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.SetRecord is expected once
func (m *mObjectStorageMockSetRecord) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 record.Record) *ObjectStorageMockSetRecordExpectation {
	m.mock.SetRecordFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockSetRecordExpectation{}
	expectation.input = &ObjectStorageMockSetRecordInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockSetRecordExpectation) Return(r *core.RecordID, r1 error) {
	e.result = &ObjectStorageMockSetRecordResult{r, r1}
}

//Set uses given function f as a mock of ObjectStorage.SetRecord method
func (m *mObjectStorageMockSetRecord) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 record.Record) (r *core.RecordID, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetRecordFunc = f
	return m.mock
}

//SetRecord implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) SetRecord(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 record.Record) (r *core.RecordID, r1 error) {
	counter := atomic.AddUint64(&m.SetRecordPreCounter, 1)
	defer atomic.AddUint64(&m.SetRecordCounter, 1)

	if len(m.SetRecordMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetRecordMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.SetRecord. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetRecordMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockSetRecordInput{p, p1, p2, p3}, "ObjectStorage.SetRecord got unexpected parameters")

		result := m.SetRecordMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.SetRecord")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetRecordMock.mainExpectation != nil {

		input := m.SetRecordMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectStorageMockSetRecordInput{p, p1, p2, p3}, "ObjectStorage.SetRecord got unexpected parameters")
		}

		result := m.SetRecordMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectStorageMock.SetRecord")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetRecordFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectStorageMock.SetRecord. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetRecordFunc(p, p1, p2, p3)
}

//SetRecordMinimockCounter returns a count of ObjectStorageMock.SetRecordFunc invocations
func (m *ObjectStorageMock) SetRecordMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetRecordCounter)
}

//SetRecordMinimockPreCounter returns the value of ObjectStorageMock.SetRecord invocations
func (m *ObjectStorageMock) SetRecordMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetRecordPreCounter)
}

//SetRecordFinished returns true if mock invocations count is ok
func (m *ObjectStorageMock) SetRecordFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetRecordMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetRecordCounter) == uint64(len(m.SetRecordMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetRecordMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetRecordCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetRecordFunc != nil {
		return atomic.LoadUint64(&m.SetRecordCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectStorageMock) ValidateCallCounters() {

	if !m.GetBlobFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetBlob")
	}

	if !m.GetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetObjectIndex")
	}

	if !m.GetRecordFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetRecord")
	}

	if !m.IterateIndexIDsFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.IterateIndexIDs")
	}

	if !m.RemoveObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.RemoveObjectIndex")
	}

	if !m.SetBlobFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetBlob")
	}

	if !m.SetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetObjectIndex")
	}

	if !m.SetRecordFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetRecord")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ObjectStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ObjectStorageMock) MinimockFinish() {

	if !m.GetBlobFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetBlob")
	}

	if !m.GetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetObjectIndex")
	}

	if !m.GetRecordFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetRecord")
	}

	if !m.IterateIndexIDsFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.IterateIndexIDs")
	}

	if !m.RemoveObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.RemoveObjectIndex")
	}

	if !m.SetBlobFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetBlob")
	}

	if !m.SetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetObjectIndex")
	}

	if !m.SetRecordFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetRecord")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ObjectStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ObjectStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetBlobFinished()
		ok = ok && m.GetObjectIndexFinished()
		ok = ok && m.GetRecordFinished()
		ok = ok && m.IterateIndexIDsFinished()
		ok = ok && m.RemoveObjectIndexFinished()
		ok = ok && m.SetBlobFinished()
		ok = ok && m.SetObjectIndexFinished()
		ok = ok && m.SetRecordFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetBlobFinished() {
				m.t.Error("Expected call to ObjectStorageMock.GetBlob")
			}

			if !m.GetObjectIndexFinished() {
				m.t.Error("Expected call to ObjectStorageMock.GetObjectIndex")
			}

			if !m.GetRecordFinished() {
				m.t.Error("Expected call to ObjectStorageMock.GetRecord")
			}

			if !m.IterateIndexIDsFinished() {
				m.t.Error("Expected call to ObjectStorageMock.IterateIndexIDs")
			}

			if !m.RemoveObjectIndexFinished() {
				m.t.Error("Expected call to ObjectStorageMock.RemoveObjectIndex")
			}

			if !m.SetBlobFinished() {
				m.t.Error("Expected call to ObjectStorageMock.SetBlob")
			}

			if !m.SetObjectIndexFinished() {
				m.t.Error("Expected call to ObjectStorageMock.SetObjectIndex")
			}

			if !m.SetRecordFinished() {
				m.t.Error("Expected call to ObjectStorageMock.SetRecord")
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
func (m *ObjectStorageMock) AllMocksCalled() bool {

	if !m.GetBlobFinished() {
		return false
	}

	if !m.GetObjectIndexFinished() {
		return false
	}

	if !m.GetRecordFinished() {
		return false
	}

	if !m.IterateIndexIDsFinished() {
		return false
	}

	if !m.RemoveObjectIndexFinished() {
		return false
	}

	if !m.SetBlobFinished() {
		return false
	}

	if !m.SetObjectIndexFinished() {
		return false
	}

	if !m.SetRecordFinished() {
		return false
	}

	return true
}
