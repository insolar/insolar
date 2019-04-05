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

	ResetFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)
	ResetCounter    uint64
	ResetPreCounter uint64
	ResetMock       mHeavySyncMockReset

	StartFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mHeavySyncMockStart

	StopFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)
	StopCounter    uint64
	StopPreCounter uint64
	StopMock       mHeavySyncMockStop

	StoreBlobsFunc       func(p context.Context, p1 insolar.PulseNumber, p2 [][]byte) (r error)
	StoreBlobsCounter    uint64
	StoreBlobsPreCounter uint64
	StoreBlobsMock       mHeavySyncMockStoreBlobs

	StoreDropFunc       func(p context.Context, p1 insolar.JetID, p2 []byte) (r error)
	StoreDropCounter    uint64
	StoreDropPreCounter uint64
	StoreDropMock       mHeavySyncMockStoreDrop

	StoreIndicesFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []insolar.KV) (r error)
	StoreIndicesCounter    uint64
	StoreIndicesPreCounter uint64
	StoreIndicesMock       mHeavySyncMockStoreIndices

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

	m.ResetMock = mHeavySyncMockReset{mock: m}
	m.StartMock = mHeavySyncMockStart{mock: m}
	m.StopMock = mHeavySyncMockStop{mock: m}
	m.StoreBlobsMock = mHeavySyncMockStoreBlobs{mock: m}
	m.StoreDropMock = mHeavySyncMockStoreDrop{mock: m}
	m.StoreIndicesMock = mHeavySyncMockStoreIndices{mock: m}
	m.StoreRecordsMock = mHeavySyncMockStoreRecords{mock: m}

	return m
}

type mHeavySyncMockReset struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockResetExpectation
	expectationSeries []*HeavySyncMockResetExpectation
}

type HeavySyncMockResetExpectation struct {
	input  *HeavySyncMockResetInput
	result *HeavySyncMockResetResult
}

type HeavySyncMockResetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type HeavySyncMockResetResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.Reset is expected from 1 to Infinity times
func (m *mHeavySyncMockReset) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mHeavySyncMockReset {
	m.mock.ResetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockResetExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockResetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of HeavySync.Reset
func (m *mHeavySyncMockReset) Return(r error) *HeavySyncMock {
	m.mock.ResetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockResetExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockResetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.Reset is expected once
func (m *mHeavySyncMockReset) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *HeavySyncMockResetExpectation {
	m.mock.ResetFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockResetExpectation{}
	expectation.input = &HeavySyncMockResetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockResetExpectation) Return(r error) {
	e.result = &HeavySyncMockResetResult{r}
}

//Set uses given function f as a mock of HeavySync.Reset method
func (m *mHeavySyncMockReset) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResetFunc = f
	return m.mock
}

//Reset implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) Reset(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.ResetPreCounter, 1)
	defer atomic.AddUint64(&m.ResetCounter, 1)

	if len(m.ResetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ResetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.Reset. %v %v %v", p, p1, p2)
			return
		}

		input := m.ResetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockResetInput{p, p1, p2}, "HeavySync.Reset got unexpected parameters")

		result := m.ResetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Reset")
			return
		}

		r = result.r

		return
	}

	if m.ResetMock.mainExpectation != nil {

		input := m.ResetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockResetInput{p, p1, p2}, "HeavySync.Reset got unexpected parameters")
		}

		result := m.ResetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Reset")
		}

		r = result.r

		return
	}

	if m.ResetFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.Reset. %v %v %v", p, p1, p2)
		return
	}

	return m.ResetFunc(p, p1, p2)
}

//ResetMinimockCounter returns a count of HeavySyncMock.ResetFunc invocations
func (m *HeavySyncMock) ResetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResetCounter)
}

//ResetMinimockPreCounter returns the value of HeavySyncMock.Reset invocations
func (m *HeavySyncMock) ResetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResetPreCounter)
}

//ResetFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) ResetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ResetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ResetCounter) == uint64(len(m.ResetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ResetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ResetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ResetFunc != nil {
		return atomic.LoadUint64(&m.ResetCounter) > 0
	}

	return true
}

type mHeavySyncMockStart struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStartExpectation
	expectationSeries []*HeavySyncMockStartExpectation
}

type HeavySyncMockStartExpectation struct {
	input  *HeavySyncMockStartInput
	result *HeavySyncMockStartResult
}

type HeavySyncMockStartInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type HeavySyncMockStartResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.Start is expected from 1 to Infinity times
func (m *mHeavySyncMockStart) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mHeavySyncMockStart {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStartExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStartInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of HeavySync.Start
func (m *mHeavySyncMockStart) Return(r error) *HeavySyncMock {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStartExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockStartResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.Start is expected once
func (m *mHeavySyncMockStart) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *HeavySyncMockStartExpectation {
	m.mock.StartFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStartExpectation{}
	expectation.input = &HeavySyncMockStartInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockStartExpectation) Return(r error) {
	e.result = &HeavySyncMockStartResult{r}
}

//Set uses given function f as a mock of HeavySync.Start method
func (m *mHeavySyncMockStart) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StartFunc = f
	return m.mock
}

//Start implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) Start(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if len(m.StartMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StartMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.Start. %v %v %v", p, p1, p2)
			return
		}

		input := m.StartMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStartInput{p, p1, p2}, "HeavySync.Start got unexpected parameters")

		result := m.StartMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Start")
			return
		}

		r = result.r

		return
	}

	if m.StartMock.mainExpectation != nil {

		input := m.StartMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStartInput{p, p1, p2}, "HeavySync.Start got unexpected parameters")
		}

		result := m.StartMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Start")
		}

		r = result.r

		return
	}

	if m.StartFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.Start. %v %v %v", p, p1, p2)
		return
	}

	return m.StartFunc(p, p1, p2)
}

//StartMinimockCounter returns a count of HeavySyncMock.StartFunc invocations
func (m *HeavySyncMock) StartMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StartCounter)
}

//StartMinimockPreCounter returns the value of HeavySyncMock.Start invocations
func (m *HeavySyncMock) StartMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StartPreCounter)
}

//StartFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StartFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StartMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StartCounter) == uint64(len(m.StartMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StartMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StartFunc != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	return true
}

type mHeavySyncMockStop struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStopExpectation
	expectationSeries []*HeavySyncMockStopExpectation
}

type HeavySyncMockStopExpectation struct {
	input  *HeavySyncMockStopInput
	result *HeavySyncMockStopResult
}

type HeavySyncMockStopInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type HeavySyncMockStopResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.Stop is expected from 1 to Infinity times
func (m *mHeavySyncMockStop) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mHeavySyncMockStop {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStopExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStopInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of HeavySync.Stop
func (m *mHeavySyncMockStop) Return(r error) *HeavySyncMock {
	m.mock.StopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStopExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockStopResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.Stop is expected once
func (m *mHeavySyncMockStop) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *HeavySyncMockStopExpectation {
	m.mock.StopFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStopExpectation{}
	expectation.input = &HeavySyncMockStopInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockStopExpectation) Return(r error) {
	e.result = &HeavySyncMockStopResult{r}
}

//Set uses given function f as a mock of HeavySync.Stop method
func (m *mHeavySyncMockStop) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StopFunc = f
	return m.mock
}

//Stop implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) Stop(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.StopPreCounter, 1)
	defer atomic.AddUint64(&m.StopCounter, 1)

	if len(m.StopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.Stop. %v %v %v", p, p1, p2)
			return
		}

		input := m.StopMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStopInput{p, p1, p2}, "HeavySync.Stop got unexpected parameters")

		result := m.StopMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Stop")
			return
		}

		r = result.r

		return
	}

	if m.StopMock.mainExpectation != nil {

		input := m.StopMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStopInput{p, p1, p2}, "HeavySync.Stop got unexpected parameters")
		}

		result := m.StopMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Stop")
		}

		r = result.r

		return
	}

	if m.StopFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.Stop. %v %v %v", p, p1, p2)
		return
	}

	return m.StopFunc(p, p1, p2)
}

//StopMinimockCounter returns a count of HeavySyncMock.StopFunc invocations
func (m *HeavySyncMock) StopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StopCounter)
}

//StopMinimockPreCounter returns the value of HeavySyncMock.Stop invocations
func (m *HeavySyncMock) StopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StopPreCounter)
}

//StopFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StopFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StopMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StopCounter) == uint64(len(m.StopMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StopMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StopCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StopFunc != nil {
		return atomic.LoadUint64(&m.StopCounter) > 0
	}

	return true
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

type mHeavySyncMockStoreIndices struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStoreIndicesExpectation
	expectationSeries []*HeavySyncMockStoreIndicesExpectation
}

type HeavySyncMockStoreIndicesExpectation struct {
	input  *HeavySyncMockStoreIndicesInput
	result *HeavySyncMockStoreIndicesResult
}

type HeavySyncMockStoreIndicesInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 []insolar.KV
}

type HeavySyncMockStoreIndicesResult struct {
	r error
}

// Expect specifies that invocation of HeavySync.StoreIndexes is expected from 1 to Infinity times
func (m *mHeavySyncMockStoreIndices) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []insolar.KV) *mHeavySyncMockStoreIndices {
	m.mock.StoreIndicesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreIndicesExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStoreIndicesInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of HeavySync.StoreIndexes
func (m *mHeavySyncMockStoreIndices) Return(r error) *HeavySyncMock {
	m.mock.StoreIndicesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreIndicesExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockStoreIndicesResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of HeavySync.StoreIndexes is expected once
func (m *mHeavySyncMockStoreIndices) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []insolar.KV) *HeavySyncMockStoreIndicesExpectation {
	m.mock.StoreIndicesFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStoreIndicesExpectation{}
	expectation.input = &HeavySyncMockStoreIndicesInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockStoreIndicesExpectation) Return(r error) {
	e.result = &HeavySyncMockStoreIndicesResult{r}
}

// Set uses given function f as a mock of HeavySync.StoreIndexes method
func (m *mHeavySyncMockStoreIndices) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []insolar.KV) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreIndicesFunc = f
	return m.mock
}

// StoreIndexes implements github.com/insolar/insolar/insolar.HeavySync interface
func (m *HeavySyncMock) StoreIndexes(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 []insolar.KV) (r error) {
	counter := atomic.AddUint64(&m.StoreIndicesPreCounter, 1)
	defer atomic.AddUint64(&m.StoreIndicesCounter, 1)

	if len(m.StoreIndicesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreIndicesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.StoreIndexes. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.StoreIndicesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStoreIndicesInput{p, p1, p2, p3}, "HeavySync.StoreIndexes got unexpected parameters")

		result := m.StoreIndicesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreIndexes")
			return
		}

		r = result.r

		return
	}

	if m.StoreIndicesMock.mainExpectation != nil {

		input := m.StoreIndicesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStoreIndicesInput{p, p1, p2, p3}, "HeavySync.StoreIndexes got unexpected parameters")
		}

		result := m.StoreIndicesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.StoreIndexes")
		}

		r = result.r

		return
	}

	if m.StoreIndicesFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.StoreIndexes. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.StoreIndicesFunc(p, p1, p2, p3)
}

//StoreIndicesMinimockCounter returns a count of HeavySyncMock.StoreIndicesFunc invocations
func (m *HeavySyncMock) StoreIndicesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreIndicesCounter)
}

// StoreIndicesMinimockPreCounter returns the value of HeavySyncMock.StoreIndexes invocations
func (m *HeavySyncMock) StoreIndicesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StoreIndicesPreCounter)
}

//StoreIndicesFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StoreIndicesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreIndicesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreIndicesCounter) == uint64(len(m.StoreIndicesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreIndicesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreIndicesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreIndicesFunc != nil {
		return atomic.LoadUint64(&m.StoreIndicesCounter) > 0
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

	if !m.ResetFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Reset")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Stop")
	}

	if !m.StoreBlobsFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreBlobs")
	}

	if !m.StoreDropFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreDrop")
	}

	if !m.StoreIndicesFinished() {
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

	if !m.ResetFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Reset")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Start")
	}

	if !m.StopFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Stop")
	}

	if !m.StoreBlobsFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreBlobs")
	}

	if !m.StoreDropFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.StoreDrop")
	}

	if !m.StoreIndicesFinished() {
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
		ok = ok && m.ResetFinished()
		ok = ok && m.StartFinished()
		ok = ok && m.StopFinished()
		ok = ok && m.StoreBlobsFinished()
		ok = ok && m.StoreDropFinished()
		ok = ok && m.StoreIndicesFinished()
		ok = ok && m.StoreRecordsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ResetFinished() {
				m.t.Error("Expected call to HeavySyncMock.Reset")
			}

			if !m.StartFinished() {
				m.t.Error("Expected call to HeavySyncMock.Start")
			}

			if !m.StopFinished() {
				m.t.Error("Expected call to HeavySyncMock.Stop")
			}

			if !m.StoreBlobsFinished() {
				m.t.Error("Expected call to HeavySyncMock.StoreBlobs")
			}

			if !m.StoreDropFinished() {
				m.t.Error("Expected call to HeavySyncMock.StoreDrop")
			}

			if !m.StoreIndicesFinished() {
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

	if !m.ResetFinished() {
		return false
	}

	if !m.StartFinished() {
		return false
	}

	if !m.StopFinished() {
		return false
	}

	if !m.StoreBlobsFinished() {
		return false
	}

	if !m.StoreDropFinished() {
		return false
	}

	if !m.StoreIndicesFinished() {
		return false
	}

	if !m.StoreRecordsFinished() {
		return false
	}

	return true
}
