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
	insolar "github.com/insolar/insolar/insolar"
	object "github.com/insolar/insolar/ledger/storage/object"

	testify_assert "github.com/stretchr/testify/assert"
)

//ObjectStorageMock implements github.com/insolar/insolar/ledger/storage.ObjectStorage
type ObjectStorageMock struct {
	t minimock.Tester

	GetBlobFunc       func(p context.Context, p1 insolar.ID, p2 *insolar.ID) (r []byte, r1 error)
	GetBlobCounter    uint64
	GetBlobPreCounter uint64
	GetBlobMock       mObjectStorageMockGetBlob

	GetObjectIndexFunc       func(p context.Context, p1 insolar.ID, p2 *insolar.ID) (r *object.Lifeline, r1 error)
	GetObjectIndexCounter    uint64
	GetObjectIndexPreCounter uint64
	GetObjectIndexMock       mObjectStorageMockGetObjectIndex

	IterateIndexIDsFunc       func(p context.Context, p1 insolar.ID, p2 func(p insolar.ID) (r error)) (r error)
	IterateIndexIDsCounter    uint64
	IterateIndexIDsPreCounter uint64
	IterateIndexIDsMock       mObjectStorageMockIterateIndexIDs

	SetBlobFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []byte) (r *insolar.ID, r1 error)
	SetBlobCounter    uint64
	SetBlobPreCounter uint64
	SetBlobMock       mObjectStorageMockSetBlob

	SetObjectIndexFunc       func(p context.Context, p1 insolar.ID, p2 *insolar.ID, p3 *object.Lifeline) (r error)
	SetObjectIndexCounter    uint64
	SetObjectIndexPreCounter uint64
	SetObjectIndexMock       mObjectStorageMockSetObjectIndex
}

//NewObjectStorageMock returns a mock for github.com/insolar/insolar/ledger/storage.ObjectStorage
func NewObjectStorageMock(t minimock.Tester) *ObjectStorageMock {
	m := &ObjectStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetBlobMock = mObjectStorageMockGetBlob{mock: m}
	m.GetObjectIndexMock = mObjectStorageMockGetObjectIndex{mock: m}
	m.IterateIndexIDsMock = mObjectStorageMockIterateIndexIDs{mock: m}
	m.SetBlobMock = mObjectStorageMockSetBlob{mock: m}
	m.SetObjectIndexMock = mObjectStorageMockSetObjectIndex{mock: m}

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
	p1 insolar.ID
	p2 *insolar.ID
}

type ObjectStorageMockGetBlobResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of ObjectStorage.GetBlob is expected from 1 to Infinity times
func (m *mObjectStorageMockGetBlob) Expect(p context.Context, p1 insolar.ID, p2 *insolar.ID) *mObjectStorageMockGetBlob {
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
func (m *mObjectStorageMockGetBlob) ExpectOnce(p context.Context, p1 insolar.ID, p2 *insolar.ID) *ObjectStorageMockGetBlobExpectation {
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
func (m *mObjectStorageMockGetBlob) Set(f func(p context.Context, p1 insolar.ID, p2 *insolar.ID) (r []byte, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetBlobFunc = f
	return m.mock
}

//GetBlob implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) GetBlob(p context.Context, p1 insolar.ID, p2 *insolar.ID) (r []byte, r1 error) {
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
	p1 insolar.ID
	p2 *insolar.ID
}

type ObjectStorageMockGetObjectIndexResult struct {
	r  *object.Lifeline
	r1 error
}

//Expect specifies that invocation of ObjectStorage.GetObjectIndex is expected from 1 to Infinity times
func (m *mObjectStorageMockGetObjectIndex) Expect(p context.Context, p1 insolar.ID, p2 *insolar.ID) *mObjectStorageMockGetObjectIndex {
	m.mock.GetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetObjectIndexExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockGetObjectIndexInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ObjectStorage.GetObjectIndex
func (m *mObjectStorageMockGetObjectIndex) Return(r *object.Lifeline, r1 error) *ObjectStorageMock {
	m.mock.GetObjectIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockGetObjectIndexExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockGetObjectIndexResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.GetObjectIndex is expected once
func (m *mObjectStorageMockGetObjectIndex) ExpectOnce(p context.Context, p1 insolar.ID, p2 *insolar.ID) *ObjectStorageMockGetObjectIndexExpectation {
	m.mock.GetObjectIndexFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockGetObjectIndexExpectation{}
	expectation.input = &ObjectStorageMockGetObjectIndexInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockGetObjectIndexExpectation) Return(r *object.Lifeline, r1 error) {
	e.result = &ObjectStorageMockGetObjectIndexResult{r, r1}
}

//Set uses given function f as a mock of ObjectStorage.GetObjectIndex method
func (m *mObjectStorageMockGetObjectIndex) Set(f func(p context.Context, p1 insolar.ID, p2 *insolar.ID) (r *object.Lifeline, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectIndexFunc = f
	return m.mock
}

//GetObjectIndex implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) GetObjectIndex(p context.Context, p1 insolar.ID, p2 *insolar.ID) (r *object.Lifeline, r1 error) {
	counter := atomic.AddUint64(&m.GetObjectIndexPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectIndexCounter, 1)

	if len(m.GetObjectIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectStorageMock.GetObjectIndex. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetObjectIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectStorageMockGetObjectIndexInput{p, p1, p2}, "ObjectStorage.GetObjectIndex got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, ObjectStorageMockGetObjectIndexInput{p, p1, p2}, "ObjectStorage.GetObjectIndex got unexpected parameters")
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
		m.t.Fatalf("Unexpected call to ObjectStorageMock.GetObjectIndex. %v %v %v", p, p1, p2)
		return
	}

	return m.GetObjectIndexFunc(p, p1, p2)
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
	p1 insolar.ID
	p2 func(p insolar.ID) (r error)
}

type ObjectStorageMockIterateIndexIDsResult struct {
	r error
}

//Expect specifies that invocation of ObjectStorage.IterateIndexIDs is expected from 1 to Infinity times
func (m *mObjectStorageMockIterateIndexIDs) Expect(p context.Context, p1 insolar.ID, p2 func(p insolar.ID) (r error)) *mObjectStorageMockIterateIndexIDs {
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
func (m *mObjectStorageMockIterateIndexIDs) ExpectOnce(p context.Context, p1 insolar.ID, p2 func(p insolar.ID) (r error)) *ObjectStorageMockIterateIndexIDsExpectation {
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
func (m *mObjectStorageMockIterateIndexIDs) Set(f func(p context.Context, p1 insolar.ID, p2 func(p insolar.ID) (r error)) (r error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IterateIndexIDsFunc = f
	return m.mock
}

//IterateIndexIDs implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) IterateIndexIDs(p context.Context, p1 insolar.ID, p2 func(p insolar.ID) (r error)) (r error) {
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
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 []byte
}

type ObjectStorageMockSetBlobResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of ObjectStorage.SetBlob is expected from 1 to Infinity times
func (m *mObjectStorageMockSetBlob) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []byte) *mObjectStorageMockSetBlob {
	m.mock.SetBlobFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetBlobExpectation{}
	}
	m.mainExpectation.input = &ObjectStorageMockSetBlobInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ObjectStorage.SetBlob
func (m *mObjectStorageMockSetBlob) Return(r *insolar.ID, r1 error) *ObjectStorageMock {
	m.mock.SetBlobFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectStorageMockSetBlobExpectation{}
	}
	m.mainExpectation.result = &ObjectStorageMockSetBlobResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectStorage.SetBlob is expected once
func (m *mObjectStorageMockSetBlob) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []byte) *ObjectStorageMockSetBlobExpectation {
	m.mock.SetBlobFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectStorageMockSetBlobExpectation{}
	expectation.input = &ObjectStorageMockSetBlobInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectStorageMockSetBlobExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ObjectStorageMockSetBlobResult{r, r1}
}

//Set uses given function f as a mock of ObjectStorage.SetBlob method
func (m *mObjectStorageMockSetBlob) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []byte) (r *insolar.ID, r1 error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetBlobFunc = f
	return m.mock
}

//SetBlob implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) SetBlob(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []byte) (r *insolar.ID, r1 error) {
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
	p1 insolar.ID
	p2 *insolar.ID
	p3 *object.Lifeline
}

type ObjectStorageMockSetObjectIndexResult struct {
	r error
}

//Expect specifies that invocation of ObjectStorage.SetObjectIndex is expected from 1 to Infinity times
func (m *mObjectStorageMockSetObjectIndex) Expect(p context.Context, p1 insolar.ID, p2 *insolar.ID, p3 *object.Lifeline) *mObjectStorageMockSetObjectIndex {
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
func (m *mObjectStorageMockSetObjectIndex) ExpectOnce(p context.Context, p1 insolar.ID, p2 *insolar.ID, p3 *object.Lifeline) *ObjectStorageMockSetObjectIndexExpectation {
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
func (m *mObjectStorageMockSetObjectIndex) Set(f func(p context.Context, p1 insolar.ID, p2 *insolar.ID, p3 *object.Lifeline) (r error)) *ObjectStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetObjectIndexFunc = f
	return m.mock
}

//SetObjectIndex implements github.com/insolar/insolar/ledger/storage.ObjectStorage interface
func (m *ObjectStorageMock) SetObjectIndex(p context.Context, p1 insolar.ID, p2 *insolar.ID, p3 *object.Lifeline) (r error) {
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectStorageMock) ValidateCallCounters() {

	if !m.GetBlobFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetBlob")
	}

	if !m.GetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.GetObjectIndex")
	}

	if !m.IterateIndexIDsFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.IterateIndexIDs")
	}

	if !m.SetBlobFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetBlob")
	}

	if !m.SetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetObjectIndex")
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

	if !m.IterateIndexIDsFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.IterateIndexIDs")
	}

	if !m.SetBlobFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetBlob")
	}

	if !m.SetObjectIndexFinished() {
		m.t.Fatal("Expected call to ObjectStorageMock.SetObjectIndex")
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
		ok = ok && m.IterateIndexIDsFinished()
		ok = ok && m.SetBlobFinished()
		ok = ok && m.SetObjectIndexFinished()

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

			if !m.IterateIndexIDsFinished() {
				m.t.Error("Expected call to ObjectStorageMock.IterateIndexIDs")
			}

			if !m.SetBlobFinished() {
				m.t.Error("Expected call to ObjectStorageMock.SetBlob")
			}

			if !m.SetObjectIndexFinished() {
				m.t.Error("Expected call to ObjectStorageMock.SetObjectIndex")
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

	if !m.IterateIndexIDsFinished() {
		return false
	}

	if !m.SetBlobFinished() {
		return false
	}

	if !m.SetObjectIndexFinished() {
		return false
	}

	return true
}
