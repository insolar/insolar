package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "HeavySync" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//HeavySyncMock implements github.com/insolar/insolar/core.HeavySync
type HeavySyncMock struct {
	t minimock.Tester

	ResetFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)
	ResetCounter    uint64
	ResetPreCounter uint64
	ResetMock       mHeavySyncMockReset

	StartFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mHeavySyncMockStart

	StopFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)
	StopCounter    uint64
	StopPreCounter uint64
	StopMock       mHeavySyncMockStop

	StoreFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []core.KV) (r error)
	StoreCounter    uint64
	StorePreCounter uint64
	StoreMock       mHeavySyncMockStore
}

//NewHeavySyncMock returns a mock for github.com/insolar/insolar/core.HeavySync
func NewHeavySyncMock(t minimock.Tester) *HeavySyncMock {
	m := &HeavySyncMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ResetMock = mHeavySyncMockReset{mock: m}
	m.StartMock = mHeavySyncMockStart{mock: m}
	m.StopMock = mHeavySyncMockStop{mock: m}
	m.StoreMock = mHeavySyncMockStore{mock: m}

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
	p1 core.RecordID
	p2 core.PulseNumber
}

type HeavySyncMockResetResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.Reset is expected from 1 to Infinity times
func (m *mHeavySyncMockReset) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mHeavySyncMockReset {
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
func (m *mHeavySyncMockReset) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *HeavySyncMockResetExpectation {
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
func (m *mHeavySyncMockReset) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResetFunc = f
	return m.mock
}

//Reset implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Reset(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error) {
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
	p1 core.RecordID
	p2 core.PulseNumber
}

type HeavySyncMockStartResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.Start is expected from 1 to Infinity times
func (m *mHeavySyncMockStart) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mHeavySyncMockStart {
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
func (m *mHeavySyncMockStart) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *HeavySyncMockStartExpectation {
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
func (m *mHeavySyncMockStart) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StartFunc = f
	return m.mock
}

//Start implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Start(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error) {
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
	p1 core.RecordID
	p2 core.PulseNumber
}

type HeavySyncMockStopResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.Stop is expected from 1 to Infinity times
func (m *mHeavySyncMockStop) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mHeavySyncMockStop {
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
func (m *mHeavySyncMockStop) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *HeavySyncMockStopExpectation {
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
func (m *mHeavySyncMockStop) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StopFunc = f
	return m.mock
}

//Stop implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Stop(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r error) {
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

type mHeavySyncMockStore struct {
	mock              *HeavySyncMock
	mainExpectation   *HeavySyncMockStoreExpectation
	expectationSeries []*HeavySyncMockStoreExpectation
}

type HeavySyncMockStoreExpectation struct {
	input  *HeavySyncMockStoreInput
	result *HeavySyncMockStoreResult
}

type HeavySyncMockStoreInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 []core.KV
}

type HeavySyncMockStoreResult struct {
	r error
}

//Expect specifies that invocation of HeavySync.Store is expected from 1 to Infinity times
func (m *mHeavySyncMockStore) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []core.KV) *mHeavySyncMockStore {
	m.mock.StoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreExpectation{}
	}
	m.mainExpectation.input = &HeavySyncMockStoreInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of HeavySync.Store
func (m *mHeavySyncMockStore) Return(r error) *HeavySyncMock {
	m.mock.StoreFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &HeavySyncMockStoreExpectation{}
	}
	m.mainExpectation.result = &HeavySyncMockStoreResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of HeavySync.Store is expected once
func (m *mHeavySyncMockStore) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []core.KV) *HeavySyncMockStoreExpectation {
	m.mock.StoreFunc = nil
	m.mainExpectation = nil

	expectation := &HeavySyncMockStoreExpectation{}
	expectation.input = &HeavySyncMockStoreInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *HeavySyncMockStoreExpectation) Return(r error) {
	e.result = &HeavySyncMockStoreResult{r}
}

//Set uses given function f as a mock of HeavySync.Store method
func (m *mHeavySyncMockStore) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []core.KV) (r error)) *HeavySyncMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreFunc = f
	return m.mock
}

//Store implements github.com/insolar/insolar/core.HeavySync interface
func (m *HeavySyncMock) Store(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []core.KV) (r error) {
	counter := atomic.AddUint64(&m.StorePreCounter, 1)
	defer atomic.AddUint64(&m.StoreCounter, 1)

	if len(m.StoreMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to HeavySyncMock.Store. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.StoreMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, HeavySyncMockStoreInput{p, p1, p2, p3}, "HeavySync.Store got unexpected parameters")

		result := m.StoreMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Store")
			return
		}

		r = result.r

		return
	}

	if m.StoreMock.mainExpectation != nil {

		input := m.StoreMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, HeavySyncMockStoreInput{p, p1, p2, p3}, "HeavySync.Store got unexpected parameters")
		}

		result := m.StoreMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the HeavySyncMock.Store")
		}

		r = result.r

		return
	}

	if m.StoreFunc == nil {
		m.t.Fatalf("Unexpected call to HeavySyncMock.Store. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.StoreFunc(p, p1, p2, p3)
}

//StoreMinimockCounter returns a count of HeavySyncMock.StoreFunc invocations
func (m *HeavySyncMock) StoreMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreCounter)
}

//StoreMinimockPreCounter returns the value of HeavySyncMock.Store invocations
func (m *HeavySyncMock) StoreMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StorePreCounter)
}

//StoreFinished returns true if mock invocations count is ok
func (m *HeavySyncMock) StoreFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreCounter) == uint64(len(m.StoreMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreFunc != nil {
		return atomic.LoadUint64(&m.StoreCounter) > 0
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

	if !m.StoreFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Store")
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

	if !m.StoreFinished() {
		m.t.Fatal("Expected call to HeavySyncMock.Store")
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
		ok = ok && m.StoreFinished()

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

			if !m.StoreFinished() {
				m.t.Error("Expected call to HeavySyncMock.Store")
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

	if !m.StoreFinished() {
		return false
	}

	return true
}
