package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "HeavySync" can be found in github.com/insolar/insolar/insolar
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//HeavySyncMock implements github.com/insolar/insolar/insolar.HeavySync
type HeavySyncMock struct {
	t minimock.Tester

	StoreBlobsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 [][]byte) (r error)
	StoreBlobsCounter    uint64
	StoreBlobsPreCounter uint64
	StoreBlobsMock       mHeavySyncMockStoreBlobs

	StoreDropFunc       func(p context.Context, p1 insolar.JetID, p2 []byte) (r error)
	StoreDropCounter    uint64
	StoreDropPreCounter uint64
	StoreDropMock       mHeavySyncMockStoreDrop

	StoreIndexesFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 map[insolar.ID][]byte) (r error)
	StoreIndexesCounter    uint64
	StoreIndexesPreCounter uint64
	StoreIndexesMock       mHeavySyncMockStoreIndexes

	StoreRecordsFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 [][]byte)
	StoreRecordsCounter    uint64
	StoreRecordsPreCounter uint64
	StoreRecordsMock       mHeavySyncMockStoreRecords
}

//NewHeavySyncMock returns a mock for github.com/insolar/insolar/insolar.HeavySync
func NewHeavySyncMock(t minimock.Tester) *HeavySyncMock {
	m := &HeavySyncMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.StoreBlobsMock = mHeavySyncMockStoreBlobs{mock: m}
	m.StoreDropMock = mHeavySyncMockStoreDrop{mock: m}
	m.StoreIndexesMock = mHeavySyncMockStoreIndexes{mock: m}
	m.StoreRecordsMock = mHeavySyncMockStoreRecords{mock: m}

	return m
}

type mHeavySyncMockStoreBlobs struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStoreBlobsExpectation
	expectationSeries []*HeavySyncMockStoreBlobsExpectation
}

type HeavySyncMockStoreBlobsExpectation struct {
	input  *HeavySyncMockStoreBlobsInput
	result *HeavySyncMockStoreBlobsResult
}

type HeavySyncMockStoreBlobsInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 [][]byte
}

type HeavySyncMockStoreBlobsResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.StoreBlobs is expected from 1 to Infinity times
func (m *mHeavySyncMockStoreBlobs) Expect(p context.Context, p1 insolar.PulseNumber, p2 [][]byte) *mHeavySyncMockStoreBlobs {
	m.mock.StoreBlobsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreBlobsExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStoreBlobsInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of HeavySync.StoreBlobs
func (m *mHeavySyncMockStoreBlobs) Return(r error) *HeavySyncMock {
	m.mock.StoreBlobsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreBlobsExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockStoreBlobsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.StoreBlobs is expected once
func (m *mHeavySyncMockStoreBlobs) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 [][]byte) *HeavySyncMockStoreBlobsExpectation {
	m.mock.StoreBlobsFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStoreBlobsExpectation{}
	expectation.input = &HeavySyncMockStoreBlobsInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockStoreBlobsExpectation) Return(r error) {
	e.result = &HeavySyncMockStoreBlobsResult{r}
}

//Set uses given function f as a mock of HeavySync.StoreBlobs method
func (m *mHeavySyncMockStoreBlobs) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 [][]byte) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreBlobsFunc = f
	return m.mock
}

//StoreBlobs implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) StoreBlobs(p context.Context, p1 insolar.PulseNumber, p2 [][]byte) (r error) {
	counter := atomic.AddUint64(&m.StoreBlobsPreCounter, 1)
	defer atomic.AddUint64(&m.StoreBlobsCounter, 1)

	if len(m.StoreBlobsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreBlobsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.StoreBlobs. %v %v %v", p, p1, p2)
			return
		}

		input := m.StoreBlobsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStoreBlobsInput{p, p1, p2}, "HeavySync.StoreBlobs got unexpected parameters")

		result := m.StoreBlobsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreBlobs")
			return
		}

		r = result.r

		return
	}

	if m.StoreBlobsMock.mainExpectation != nil {

		input := m.StoreBlobsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStoreBlobsInput{p, p1, p2}, "HeavySync.StoreBlobs got unexpected parameters")
		}

		result := m.StoreBlobsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreBlobs")
		}

		r = result.r

		return
	}

	if m.StoreBlobsFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.StoreBlobs. %v %v %v", p, p1, p2)
		return
	}

	return m.StoreBlobsFunc(p, p1, p2)
}

//StoreBlobsMinimockCounter returns a count of HeavySyncMock.StoreBlobsFunc invocations
func (m *HeavySyncMock) StoreBlobsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreBlobsCounter)
}

//StoreBlobsMinimockPreCounter returns the value of HeavySyncMock.StoreBlobs invocations
func (m *HeavySyncMock) StoreBlobsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StoreBlobsPreCounter)
}

//StoreBlobsFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StoreBlobsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreBlobsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreBlobsCounter) == uint64(len(m.StoreBlobsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreBlobsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreBlobsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreBlobsFunc != nil {
		return atomic.LoadUint64(&m.StoreBlobsCounter) > 0
	}

	return true
}

type mHeavySyncMockStoreDrop struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStoreDropExpectation
	expectationSeries []*HeavySyncMockStoreDropExpectation
}

type HeavySyncMockStoreDropExpectation struct {
	input  *HeavySyncMockStoreDropInput
	result *HeavySyncMockStoreDropResult
}

type HeavySyncMockStoreDropInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 []byte
}

type HeavySyncMockStoreDropResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.StoreDrop is expected from 1 to Infinity times
func (m *mHeavySyncMockStoreDrop) Expect(p context.Context, p1 insolar.JetID, p2 []byte) *mHeavySyncMockStoreDrop {
	m.mock.StoreDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreDropExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStoreDropInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of HeavySync.StoreDrop
func (m *mHeavySyncMockStoreDrop) Return(r error) *HeavySyncMock {
	m.mock.StoreDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreDropExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockStoreDropResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.StoreDrop is expected once
func (m *mHeavySyncMockStoreDrop) ExpectOnce(p context.Context, p1 insolar.JetID, p2 []byte) *HeavySyncMockStoreDropExpectation {
	m.mock.StoreDropFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStoreDropExpectation{}
	expectation.input = &HeavySyncMockStoreDropInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockStoreDropExpectation) Return(r error) {
	e.result = &HeavySyncMockStoreDropResult{r}
}

//Set uses given function f as a mock of HeavySync.StoreDrop method
func (m *mHeavySyncMockStoreDrop) Set(f func(p context.Context, p1 insolar.JetID, p2 []byte) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreDropFunc = f
	return m.mock
}

//StoreDrop implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) StoreDrop(p context.Context, p1 insolar.JetID, p2 []byte) (r error) {
	counter := atomic.AddUint64(&m.StoreDropPreCounter, 1)
	defer atomic.AddUint64(&m.StoreDropCounter, 1)

	if len(m.StoreDropMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreDropMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.StoreDrop. %v %v %v", p, p1, p2)
			return
		}

		input := m.StoreDropMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStoreDropInput{p, p1, p2}, "HeavySync.StoreDrop got unexpected parameters")

		result := m.StoreDropMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreDrop")
			return
		}

		r = result.r

		return
	}

	if m.StoreDropMock.mainExpectation != nil {

		input := m.StoreDropMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStoreDropInput{p, p1, p2}, "HeavySync.StoreDrop got unexpected parameters")
		}

		result := m.StoreDropMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreDrop")
		}

		r = result.r

		return
	}

	if m.StoreDropFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.StoreDrop. %v %v %v", p, p1, p2)
		return
	}

	return m.StoreDropFunc(p, p1, p2)
}

//StoreDropMinimockCounter returns a count of HeavySyncMock.StoreDropFunc invocations
func (m *HeavySyncMock) StoreDropMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreDropCounter)
}

//StoreDropMinimockPreCounter returns the value of HeavySyncMock.StoreDrop invocations
func (m *HeavySyncMock) StoreDropMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StoreDropPreCounter)
}

//StoreDropFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StoreDropFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreDropMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreDropCounter) == uint64(len(m.StoreDropMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreDropMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreDropCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreDropFunc != nil {
		return atomic.LoadUint64(&m.StoreDropCounter) > 0
	}

	return true
}

type mHeavySyncMockStoreIndexes struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStoreIndexesExpectation
	expectationSeries []*HeavySyncMockStoreIndexesExpectation
}

type HeavySyncMockStoreIndexesExpectation struct {
	input  *HeavySyncMockStoreIndexesInput
	result *HeavySyncMockStoreIndexesResult
}

type HeavySyncMockStoreIndexesInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 map[insolar.ID][]byte
}

type HeavySyncMockStoreIndexesResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.StoreIndexes is expected from 1 to Infinity times
func (m *mHeavySyncMockStoreIndexes) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 map[insolar.ID][]byte) *mHeavySyncMockStoreIndexes {
	m.mock.StoreIndexesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreIndexesExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStoreIndexesInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of HeavySync.StoreIndexes
func (m *mHeavySyncMockStoreIndexes) Return(r error) *HeavySyncMock {
	m.mock.StoreIndexesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreIndexesExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockStoreIndexesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.StoreIndexes is expected once
func (m *mHeavySyncMockStoreIndexes) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 map[insolar.ID][]byte) *HeavySyncMockStoreIndexesExpectation {
	m.mock.StoreIndexesFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStoreIndexesExpectation{}
	expectation.input = &HeavySyncMockStoreIndexesInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockStoreIndexesExpectation) Return(r error) {
	e.result = &HeavySyncMockStoreIndexesResult{r}
}

//Set uses given function f as a mock of HeavySync.StoreIndexes method
func (m *mHeavySyncMockStoreIndexes) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 map[insolar.ID][]byte) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreIndexesFunc = f
	return m.mock
}

//StoreIndexes implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) StoreIndexes(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 map[insolar.ID][]byte) (r error) {
	counter := atomic.AddUint64(&m.StoreIndexesPreCounter, 1)
	defer atomic.AddUint64(&m.StoreIndexesCounter, 1)

	if len(m.StoreIndexesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreIndexesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.StoreIndexes. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.StoreIndexesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStoreIndexesInput{p, p1, p2, p3}, "HeavySync.StoreIndexes got unexpected parameters")

		result := m.StoreIndexesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreIndexes")
			return
		}

		r = result.r

		return
	}

	if m.StoreIndexesMock.mainExpectation != nil {

		input := m.StoreIndexesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStoreIndexesInput{p, p1, p2, p3}, "HeavySync.StoreIndexes got unexpected parameters")
		}

		result := m.StoreIndexesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreIndexes")
		}

		r = result.r

		return
	}

	if m.StoreIndexesFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.StoreIndexes. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.StoreIndexesFunc(p, p1, p2, p3)
}

//StoreIndexesMinimockCounter returns a count of HeavySyncMock.StoreIndexesFunc invocations
func (m *HeavySyncMock) StoreIndexesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreIndexesCounter)
}

//StoreIndexesMinimockPreCounter returns the value of HeavySyncMock.StoreIndexes invocations
func (m *HeavySyncMock) StoreIndexesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StoreIndexesPreCounter)
}

//StoreIndexesFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StoreIndexesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreIndexesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreIndexesCounter) == uint64(len(m.StoreIndexesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreIndexesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreIndexesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreIndexesFunc != nil {
		return atomic.LoadUint64(&m.StoreIndexesCounter) > 0
	}

	return true
}

type mHeavySyncMockStoreRecords struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStoreRecordsExpectation
	expectationSeries []*HeavySyncMockStoreRecordsExpectation
}

type HeavySyncMockStoreRecordsExpectation struct {
	input *HeavySyncMockStoreRecordsInput
}

type HeavySyncMockStoreRecordsInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 [][]byte
}

//Expect specifies that invocation of HeavySync.StoreRecords is expected from 1 to Infinity times
func (m *mHeavySyncMockStoreRecords) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 [][]byte) *mHeavySyncMockStoreRecords {
	m.mock.StoreRecordsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreRecordsExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStoreRecordsInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of HeavySync.StoreRecords
func (m *mHeavySyncMockStoreRecords) Return() *HeavySyncMock {
	m.mock.StoreRecordsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreRecordsExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.StoreRecords is expected once
func (m *mHeavySyncMockStoreRecords) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 [][]byte) *HeavySyncMockStoreRecordsExpectation {
	m.mock.StoreRecordsFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStoreRecordsExpectation{}
	expectation.input = &HeavySyncMockStoreRecordsInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of HeavySync.StoreRecords method
func (m *mHeavySyncMockStoreRecords) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 [][]byte)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreRecordsFunc = f
	return m.mock
}

//StoreRecords implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) StoreRecords(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 [][]byte) {
	counter := atomic.AddUint64(&m.StoreRecordsPreCounter, 1)
	defer atomic.AddUint64(&m.StoreRecordsCounter, 1)

	if len(m.StoreRecordsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreRecordsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.StoreRecords. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.StoreRecordsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStoreRecordsInput{p, p1, p2, p3}, "HeavySync.StoreRecords got unexpected parameters")

		return
	}

	if m.StoreRecordsMock.mainExpectation != nil {

		input := m.StoreRecordsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStoreRecordsInput{p, p1, p2, p3}, "HeavySync.StoreRecords got unexpected parameters")
		}

		return
	}

	if m.StoreRecordsFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.StoreRecords. %v %v %v %v", p, p1, p2, p3)
		return
	}

	m.StoreRecordsFunc(p, p1, p2, p3)
}

//StoreRecordsMinimockCounter returns a count of HeavySyncMock.StoreRecordsFunc invocations
func (m *HeavySyncMock) StoreRecordsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreRecordsCounter)
}

//StoreRecordsMinimockPreCounter returns the value of HeavySyncMock.StoreRecords invocations
func (m *HeavySyncMock) StoreRecordsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StoreRecordsPreCounter)
}

//StoreRecordsFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StoreRecordsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreRecordsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreRecordsCounter) == uint64(len(m.StoreRecordsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreRecordsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreRecordsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreRecordsFunc != nil {
		return atomic.LoadUint64(&m.StoreRecordsCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HeavySyncMock) ValidateCallCounters() {

	if !m.StoreBlobsFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreBlobs")
	}

	if !m.StoreDropFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreDrop")
	}

	if !m.StoreIndexesFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreIndexes")
	}

	if !m.StoreRecordsFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreRecords")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *HeavySyncMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *HeavySyncMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *HeavySyncMock) MinimockFinish() {

	if !m.StoreBlobsFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreBlobs")
	}

	if !m.StoreDropFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreDrop")
	}

	if !m.StoreIndexesFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreIndexes")
	}

	if !m.StoreRecordsFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreRecords")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *HeavySyncMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *HeavySyncMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.StoreBlobsFinished()
		ok = ok && m.StoreDropFinished()
		ok = ok && m.StoreIndexesFinished()
		ok = ok && m.StoreRecordsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.StoreBlobsFinished() {
				m.t.Error("Expected call to HeavySyncMock.StoreBlobs")
			}

			if !m.StoreDropFinished() {
				m.t.Error("Expected call to HeavySyncMock.StoreDrop")
			}

			if !m.StoreIndexesFinished() {
				m.t.Error("Expected call to HeavySyncMock.StoreIndexes")
			}

			if !m.StoreRecordsFinished() {
				m.t.Error("Expected call to HeavySyncMock.StoreRecords")
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
func (m *HeavySyncMock) AllMocksCalled() bool {

	if !m.StoreBlobsFinished() {
		return false
	}

	if !m.StoreDropFinished() {
		return false
	}

	if !m.StoreIndexesFinished() {
		return false
	}

	if !m.StoreRecordsFinished() {
		return false
	}

	return true
}
