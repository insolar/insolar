package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DropStorage" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	jet "github.com/insolar/insolar/ledger/storage/jet"

	testify_assert "github.com/stretchr/testify/assert"
)

//DropStorageMock implements github.com/insolar/insolar/ledger/storage.DropStorage
type DropStorageMock struct {
	t minimock.Tester

	AddDropSizeFunc       func(p context.Context, p1 *jet.DropSize) (r error)
	AddDropSizeCounter    uint64
	AddDropSizePreCounter uint64
	AddDropSizeMock       mDropStorageMockAddDropSize

	CreateDropFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error)
	CreateDropCounter    uint64
	CreateDropPreCounter uint64
	CreateDropMock       mDropStorageMockCreateDrop

	GetDropFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *jet.JetDrop, r1 error)
	GetDropCounter    uint64
	GetDropPreCounter uint64
	GetDropMock       mDropStorageMockGetDrop

	GetDropSizeHistoryFunc       func(p context.Context, p1 core.RecordID) (r jet.DropSizeHistory, r1 error)
	GetDropSizeHistoryCounter    uint64
	GetDropSizeHistoryPreCounter uint64
	GetDropSizeHistoryMock       mDropStorageMockGetDropSizeHistory

	GetJetSizesHistoryDepthFunc       func() (r int)
	GetJetSizesHistoryDepthCounter    uint64
	GetJetSizesHistoryDepthPreCounter uint64
	GetJetSizesHistoryDepthMock       mDropStorageMockGetJetSizesHistoryDepth

	SetDropFunc       func(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) (r error)
	SetDropCounter    uint64
	SetDropPreCounter uint64
	SetDropMock       mDropStorageMockSetDrop

	SetDropSizeHistoryFunc       func(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) (r error)
	SetDropSizeHistoryCounter    uint64
	SetDropSizeHistoryPreCounter uint64
	SetDropSizeHistoryMock       mDropStorageMockSetDropSizeHistory
}

//NewDropStorageMock returns a mock for github.com/insolar/insolar/ledger/storage.DropStorage
func NewDropStorageMock(t minimock.Tester) *DropStorageMock {
	m := &DropStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddDropSizeMock = mDropStorageMockAddDropSize{mock: m}
	m.CreateDropMock = mDropStorageMockCreateDrop{mock: m}
	m.GetDropMock = mDropStorageMockGetDrop{mock: m}
	m.GetDropSizeHistoryMock = mDropStorageMockGetDropSizeHistory{mock: m}
	m.GetJetSizesHistoryDepthMock = mDropStorageMockGetJetSizesHistoryDepth{mock: m}
	m.SetDropMock = mDropStorageMockSetDrop{mock: m}
	m.SetDropSizeHistoryMock = mDropStorageMockSetDropSizeHistory{mock: m}

	return m
}

type mDropStorageMockAddDropSize struct {
	mock              *DropStorageMock
	mainExpectation   *DropStorageMockAddDropSizeExpectation
	expectationSeries []*DropStorageMockAddDropSizeExpectation
}

type DropStorageMockAddDropSizeExpectation struct {
	input  *DropStorageMockAddDropSizeInput
	result *DropStorageMockAddDropSizeResult
}

type DropStorageMockAddDropSizeInput struct {
	p  context.Context
	p1 *jet.DropSize
}

type DropStorageMockAddDropSizeResult struct {
	r error
}

//Expect specifies that invocation of DropStorage.AddDropSize is expected from 1 to Infinity times
func (m *mDropStorageMockAddDropSize) Expect(p context.Context, p1 *jet.DropSize) *mDropStorageMockAddDropSize {
	m.mock.AddDropSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockAddDropSizeExpectation{}
	}
	m.mainExpectation.input = &DropStorageMockAddDropSizeInput{p, p1}
	return m
}

//Return specifies results of invocation of DropStorage.AddDropSize
func (m *mDropStorageMockAddDropSize) Return(r error) *DropStorageMock {
	m.mock.AddDropSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockAddDropSizeExpectation{}
	}
	m.mainExpectation.result = &DropStorageMockAddDropSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DropStorage.AddDropSize is expected once
func (m *mDropStorageMockAddDropSize) ExpectOnce(p context.Context, p1 *jet.DropSize) *DropStorageMockAddDropSizeExpectation {
	m.mock.AddDropSizeFunc = nil
	m.mainExpectation = nil

	expectation := &DropStorageMockAddDropSizeExpectation{}
	expectation.input = &DropStorageMockAddDropSizeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropStorageMockAddDropSizeExpectation) Return(r error) {
	e.result = &DropStorageMockAddDropSizeResult{r}
}

//Set uses given function f as a mock of DropStorage.AddDropSize method
func (m *mDropStorageMockAddDropSize) Set(f func(p context.Context, p1 *jet.DropSize) (r error)) *DropStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddDropSizeFunc = f
	return m.mock
}

//AddDropSize implements github.com/insolar/insolar/ledger/storage.DropStorage interface
func (m *DropStorageMock) AddDropSize(p context.Context, p1 *jet.DropSize) (r error) {
	counter := atomic.AddUint64(&m.AddDropSizePreCounter, 1)
	defer atomic.AddUint64(&m.AddDropSizeCounter, 1)

	if len(m.AddDropSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddDropSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropStorageMock.AddDropSize. %v %v", p, p1)
			return
		}

		input := m.AddDropSizeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropStorageMockAddDropSizeInput{p, p1}, "DropStorage.AddDropSize got unexpected parameters")

		result := m.AddDropSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.AddDropSize")
			return
		}

		r = result.r

		return
	}

	if m.AddDropSizeMock.mainExpectation != nil {

		input := m.AddDropSizeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropStorageMockAddDropSizeInput{p, p1}, "DropStorage.AddDropSize got unexpected parameters")
		}

		result := m.AddDropSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.AddDropSize")
		}

		r = result.r

		return
	}

	if m.AddDropSizeFunc == nil {
		m.t.Fatalf("Unexpected call to DropStorageMock.AddDropSize. %v %v", p, p1)
		return
	}

	return m.AddDropSizeFunc(p, p1)
}

//AddDropSizeMinimockCounter returns a count of DropStorageMock.AddDropSizeFunc invocations
func (m *DropStorageMock) AddDropSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddDropSizeCounter)
}

//AddDropSizeMinimockPreCounter returns the value of DropStorageMock.AddDropSize invocations
func (m *DropStorageMock) AddDropSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddDropSizePreCounter)
}

//AddDropSizeFinished returns true if mock invocations count is ok
func (m *DropStorageMock) AddDropSizeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddDropSizeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddDropSizeCounter) == uint64(len(m.AddDropSizeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddDropSizeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddDropSizeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddDropSizeFunc != nil {
		return atomic.LoadUint64(&m.AddDropSizeCounter) > 0
	}

	return true
}

type mDropStorageMockCreateDrop struct {
	mock              *DropStorageMock
	mainExpectation   *DropStorageMockCreateDropExpectation
	expectationSeries []*DropStorageMockCreateDropExpectation
}

type DropStorageMockCreateDropExpectation struct {
	input  *DropStorageMockCreateDropInput
	result *DropStorageMockCreateDropResult
}

type DropStorageMockCreateDropInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 []byte
}

type DropStorageMockCreateDropResult struct {
	r  *jet.JetDrop
	r1 [][]byte
	r2 uint64
	r3 error
}

//Expect specifies that invocation of DropStorage.CreateDrop is expected from 1 to Infinity times
func (m *mDropStorageMockCreateDrop) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) *mDropStorageMockCreateDrop {
	m.mock.CreateDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockCreateDropExpectation{}
	}
	m.mainExpectation.input = &DropStorageMockCreateDropInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of DropStorage.CreateDrop
func (m *mDropStorageMockCreateDrop) Return(r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error) *DropStorageMock {
	m.mock.CreateDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockCreateDropExpectation{}
	}
	m.mainExpectation.result = &DropStorageMockCreateDropResult{r, r1, r2, r3}
	return m.mock
}

//ExpectOnce specifies that invocation of DropStorage.CreateDrop is expected once
func (m *mDropStorageMockCreateDrop) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) *DropStorageMockCreateDropExpectation {
	m.mock.CreateDropFunc = nil
	m.mainExpectation = nil

	expectation := &DropStorageMockCreateDropExpectation{}
	expectation.input = &DropStorageMockCreateDropInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropStorageMockCreateDropExpectation) Return(r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error) {
	e.result = &DropStorageMockCreateDropResult{r, r1, r2, r3}
}

//Set uses given function f as a mock of DropStorage.CreateDrop method
func (m *mDropStorageMockCreateDrop) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error)) *DropStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateDropFunc = f
	return m.mock
}

//CreateDrop implements github.com/insolar/insolar/ledger/storage.DropStorage interface
func (m *DropStorageMock) CreateDrop(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error) {
	counter := atomic.AddUint64(&m.CreateDropPreCounter, 1)
	defer atomic.AddUint64(&m.CreateDropCounter, 1)

	if len(m.CreateDropMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateDropMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropStorageMock.CreateDrop. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.CreateDropMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropStorageMockCreateDropInput{p, p1, p2, p3}, "DropStorage.CreateDrop got unexpected parameters")

		result := m.CreateDropMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.CreateDrop")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2
		r3 = result.r3

		return
	}

	if m.CreateDropMock.mainExpectation != nil {

		input := m.CreateDropMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropStorageMockCreateDropInput{p, p1, p2, p3}, "DropStorage.CreateDrop got unexpected parameters")
		}

		result := m.CreateDropMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.CreateDrop")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2
		r3 = result.r3

		return
	}

	if m.CreateDropFunc == nil {
		m.t.Fatalf("Unexpected call to DropStorageMock.CreateDrop. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.CreateDropFunc(p, p1, p2, p3)
}

//CreateDropMinimockCounter returns a count of DropStorageMock.CreateDropFunc invocations
func (m *DropStorageMock) CreateDropMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateDropCounter)
}

//CreateDropMinimockPreCounter returns the value of DropStorageMock.CreateDrop invocations
func (m *DropStorageMock) CreateDropMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateDropPreCounter)
}

//CreateDropFinished returns true if mock invocations count is ok
func (m *DropStorageMock) CreateDropFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreateDropMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreateDropCounter) == uint64(len(m.CreateDropMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreateDropMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreateDropCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreateDropFunc != nil {
		return atomic.LoadUint64(&m.CreateDropCounter) > 0
	}

	return true
}

type mDropStorageMockGetDrop struct {
	mock              *DropStorageMock
	mainExpectation   *DropStorageMockGetDropExpectation
	expectationSeries []*DropStorageMockGetDropExpectation
}

type DropStorageMockGetDropExpectation struct {
	input  *DropStorageMockGetDropInput
	result *DropStorageMockGetDropResult
}

type DropStorageMockGetDropInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

type DropStorageMockGetDropResult struct {
	r  *jet.JetDrop
	r1 error
}

//Expect specifies that invocation of DropStorage.GetDrop is expected from 1 to Infinity times
func (m *mDropStorageMockGetDrop) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mDropStorageMockGetDrop {
	m.mock.GetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockGetDropExpectation{}
	}
	m.mainExpectation.input = &DropStorageMockGetDropInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DropStorage.GetDrop
func (m *mDropStorageMockGetDrop) Return(r *jet.JetDrop, r1 error) *DropStorageMock {
	m.mock.GetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockGetDropExpectation{}
	}
	m.mainExpectation.result = &DropStorageMockGetDropResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DropStorage.GetDrop is expected once
func (m *mDropStorageMockGetDrop) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *DropStorageMockGetDropExpectation {
	m.mock.GetDropFunc = nil
	m.mainExpectation = nil

	expectation := &DropStorageMockGetDropExpectation{}
	expectation.input = &DropStorageMockGetDropInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropStorageMockGetDropExpectation) Return(r *jet.JetDrop, r1 error) {
	e.result = &DropStorageMockGetDropResult{r, r1}
}

//Set uses given function f as a mock of DropStorage.GetDrop method
func (m *mDropStorageMockGetDrop) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *jet.JetDrop, r1 error)) *DropStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDropFunc = f
	return m.mock
}

//GetDrop implements github.com/insolar/insolar/ledger/storage.DropStorage interface
func (m *DropStorageMock) GetDrop(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *jet.JetDrop, r1 error) {
	counter := atomic.AddUint64(&m.GetDropPreCounter, 1)
	defer atomic.AddUint64(&m.GetDropCounter, 1)

	if len(m.GetDropMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDropMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropStorageMock.GetDrop. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetDropMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropStorageMockGetDropInput{p, p1, p2}, "DropStorage.GetDrop got unexpected parameters")

		result := m.GetDropMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.GetDrop")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropMock.mainExpectation != nil {

		input := m.GetDropMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropStorageMockGetDropInput{p, p1, p2}, "DropStorage.GetDrop got unexpected parameters")
		}

		result := m.GetDropMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.GetDrop")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropFunc == nil {
		m.t.Fatalf("Unexpected call to DropStorageMock.GetDrop. %v %v %v", p, p1, p2)
		return
	}

	return m.GetDropFunc(p, p1, p2)
}

//GetDropMinimockCounter returns a count of DropStorageMock.GetDropFunc invocations
func (m *DropStorageMock) GetDropMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropCounter)
}

//GetDropMinimockPreCounter returns the value of DropStorageMock.GetDrop invocations
func (m *DropStorageMock) GetDropMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropPreCounter)
}

//GetDropFinished returns true if mock invocations count is ok
func (m *DropStorageMock) GetDropFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDropMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDropCounter) == uint64(len(m.GetDropMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDropMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDropCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDropFunc != nil {
		return atomic.LoadUint64(&m.GetDropCounter) > 0
	}

	return true
}

type mDropStorageMockGetDropSizeHistory struct {
	mock              *DropStorageMock
	mainExpectation   *DropStorageMockGetDropSizeHistoryExpectation
	expectationSeries []*DropStorageMockGetDropSizeHistoryExpectation
}

type DropStorageMockGetDropSizeHistoryExpectation struct {
	input  *DropStorageMockGetDropSizeHistoryInput
	result *DropStorageMockGetDropSizeHistoryResult
}

type DropStorageMockGetDropSizeHistoryInput struct {
	p  context.Context
	p1 core.RecordID
}

type DropStorageMockGetDropSizeHistoryResult struct {
	r  jet.DropSizeHistory
	r1 error
}

//Expect specifies that invocation of DropStorage.GetDropSizeHistory is expected from 1 to Infinity times
func (m *mDropStorageMockGetDropSizeHistory) Expect(p context.Context, p1 core.RecordID) *mDropStorageMockGetDropSizeHistory {
	m.mock.GetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockGetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.input = &DropStorageMockGetDropSizeHistoryInput{p, p1}
	return m
}

//Return specifies results of invocation of DropStorage.GetDropSizeHistory
func (m *mDropStorageMockGetDropSizeHistory) Return(r jet.DropSizeHistory, r1 error) *DropStorageMock {
	m.mock.GetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockGetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.result = &DropStorageMockGetDropSizeHistoryResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DropStorage.GetDropSizeHistory is expected once
func (m *mDropStorageMockGetDropSizeHistory) ExpectOnce(p context.Context, p1 core.RecordID) *DropStorageMockGetDropSizeHistoryExpectation {
	m.mock.GetDropSizeHistoryFunc = nil
	m.mainExpectation = nil

	expectation := &DropStorageMockGetDropSizeHistoryExpectation{}
	expectation.input = &DropStorageMockGetDropSizeHistoryInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropStorageMockGetDropSizeHistoryExpectation) Return(r jet.DropSizeHistory, r1 error) {
	e.result = &DropStorageMockGetDropSizeHistoryResult{r, r1}
}

//Set uses given function f as a mock of DropStorage.GetDropSizeHistory method
func (m *mDropStorageMockGetDropSizeHistory) Set(f func(p context.Context, p1 core.RecordID) (r jet.DropSizeHistory, r1 error)) *DropStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDropSizeHistoryFunc = f
	return m.mock
}

//GetDropSizeHistory implements github.com/insolar/insolar/ledger/storage.DropStorage interface
func (m *DropStorageMock) GetDropSizeHistory(p context.Context, p1 core.RecordID) (r jet.DropSizeHistory, r1 error) {
	counter := atomic.AddUint64(&m.GetDropSizeHistoryPreCounter, 1)
	defer atomic.AddUint64(&m.GetDropSizeHistoryCounter, 1)

	if len(m.GetDropSizeHistoryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDropSizeHistoryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropStorageMock.GetDropSizeHistory. %v %v", p, p1)
			return
		}

		input := m.GetDropSizeHistoryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropStorageMockGetDropSizeHistoryInput{p, p1}, "DropStorage.GetDropSizeHistory got unexpected parameters")

		result := m.GetDropSizeHistoryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.GetDropSizeHistory")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropSizeHistoryMock.mainExpectation != nil {

		input := m.GetDropSizeHistoryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropStorageMockGetDropSizeHistoryInput{p, p1}, "DropStorage.GetDropSizeHistory got unexpected parameters")
		}

		result := m.GetDropSizeHistoryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.GetDropSizeHistory")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropSizeHistoryFunc == nil {
		m.t.Fatalf("Unexpected call to DropStorageMock.GetDropSizeHistory. %v %v", p, p1)
		return
	}

	return m.GetDropSizeHistoryFunc(p, p1)
}

//GetDropSizeHistoryMinimockCounter returns a count of DropStorageMock.GetDropSizeHistoryFunc invocations
func (m *DropStorageMock) GetDropSizeHistoryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropSizeHistoryCounter)
}

//GetDropSizeHistoryMinimockPreCounter returns the value of DropStorageMock.GetDropSizeHistory invocations
func (m *DropStorageMock) GetDropSizeHistoryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropSizeHistoryPreCounter)
}

//GetDropSizeHistoryFinished returns true if mock invocations count is ok
func (m *DropStorageMock) GetDropSizeHistoryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDropSizeHistoryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDropSizeHistoryCounter) == uint64(len(m.GetDropSizeHistoryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDropSizeHistoryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDropSizeHistoryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDropSizeHistoryFunc != nil {
		return atomic.LoadUint64(&m.GetDropSizeHistoryCounter) > 0
	}

	return true
}

type mDropStorageMockGetJetSizesHistoryDepth struct {
	mock              *DropStorageMock
	mainExpectation   *DropStorageMockGetJetSizesHistoryDepthExpectation
	expectationSeries []*DropStorageMockGetJetSizesHistoryDepthExpectation
}

type DropStorageMockGetJetSizesHistoryDepthExpectation struct {
	result *DropStorageMockGetJetSizesHistoryDepthResult
}

type DropStorageMockGetJetSizesHistoryDepthResult struct {
	r int
}

//Expect specifies that invocation of DropStorage.GetJetSizesHistoryDepth is expected from 1 to Infinity times
func (m *mDropStorageMockGetJetSizesHistoryDepth) Expect() *mDropStorageMockGetJetSizesHistoryDepth {
	m.mock.GetJetSizesHistoryDepthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockGetJetSizesHistoryDepthExpectation{}
	}

	return m
}

//Return specifies results of invocation of DropStorage.GetJetSizesHistoryDepth
func (m *mDropStorageMockGetJetSizesHistoryDepth) Return(r int) *DropStorageMock {
	m.mock.GetJetSizesHistoryDepthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockGetJetSizesHistoryDepthExpectation{}
	}
	m.mainExpectation.result = &DropStorageMockGetJetSizesHistoryDepthResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DropStorage.GetJetSizesHistoryDepth is expected once
func (m *mDropStorageMockGetJetSizesHistoryDepth) ExpectOnce() *DropStorageMockGetJetSizesHistoryDepthExpectation {
	m.mock.GetJetSizesHistoryDepthFunc = nil
	m.mainExpectation = nil

	expectation := &DropStorageMockGetJetSizesHistoryDepthExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropStorageMockGetJetSizesHistoryDepthExpectation) Return(r int) {
	e.result = &DropStorageMockGetJetSizesHistoryDepthResult{r}
}

//Set uses given function f as a mock of DropStorage.GetJetSizesHistoryDepth method
func (m *mDropStorageMockGetJetSizesHistoryDepth) Set(f func() (r int)) *DropStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetJetSizesHistoryDepthFunc = f
	return m.mock
}

//GetJetSizesHistoryDepth implements github.com/insolar/insolar/ledger/storage.DropStorage interface
func (m *DropStorageMock) GetJetSizesHistoryDepth() (r int) {
	counter := atomic.AddUint64(&m.GetJetSizesHistoryDepthPreCounter, 1)
	defer atomic.AddUint64(&m.GetJetSizesHistoryDepthCounter, 1)

	if len(m.GetJetSizesHistoryDepthMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetJetSizesHistoryDepthMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropStorageMock.GetJetSizesHistoryDepth.")
			return
		}

		result := m.GetJetSizesHistoryDepthMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.GetJetSizesHistoryDepth")
			return
		}

		r = result.r

		return
	}

	if m.GetJetSizesHistoryDepthMock.mainExpectation != nil {

		result := m.GetJetSizesHistoryDepthMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.GetJetSizesHistoryDepth")
		}

		r = result.r

		return
	}

	if m.GetJetSizesHistoryDepthFunc == nil {
		m.t.Fatalf("Unexpected call to DropStorageMock.GetJetSizesHistoryDepth.")
		return
	}

	return m.GetJetSizesHistoryDepthFunc()
}

//GetJetSizesHistoryDepthMinimockCounter returns a count of DropStorageMock.GetJetSizesHistoryDepthFunc invocations
func (m *DropStorageMock) GetJetSizesHistoryDepthMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetJetSizesHistoryDepthCounter)
}

//GetJetSizesHistoryDepthMinimockPreCounter returns the value of DropStorageMock.GetJetSizesHistoryDepth invocations
func (m *DropStorageMock) GetJetSizesHistoryDepthMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetJetSizesHistoryDepthPreCounter)
}

//GetJetSizesHistoryDepthFinished returns true if mock invocations count is ok
func (m *DropStorageMock) GetJetSizesHistoryDepthFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetJetSizesHistoryDepthMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetJetSizesHistoryDepthCounter) == uint64(len(m.GetJetSizesHistoryDepthMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetJetSizesHistoryDepthMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetJetSizesHistoryDepthCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetJetSizesHistoryDepthFunc != nil {
		return atomic.LoadUint64(&m.GetJetSizesHistoryDepthCounter) > 0
	}

	return true
}

type mDropStorageMockSetDrop struct {
	mock              *DropStorageMock
	mainExpectation   *DropStorageMockSetDropExpectation
	expectationSeries []*DropStorageMockSetDropExpectation
}

type DropStorageMockSetDropExpectation struct {
	input  *DropStorageMockSetDropInput
	result *DropStorageMockSetDropResult
}

type DropStorageMockSetDropInput struct {
	p  context.Context
	p1 core.RecordID
	p2 *jet.JetDrop
}

type DropStorageMockSetDropResult struct {
	r error
}

//Expect specifies that invocation of DropStorage.SetDrop is expected from 1 to Infinity times
func (m *mDropStorageMockSetDrop) Expect(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) *mDropStorageMockSetDrop {
	m.mock.SetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockSetDropExpectation{}
	}
	m.mainExpectation.input = &DropStorageMockSetDropInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DropStorage.SetDrop
func (m *mDropStorageMockSetDrop) Return(r error) *DropStorageMock {
	m.mock.SetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockSetDropExpectation{}
	}
	m.mainExpectation.result = &DropStorageMockSetDropResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DropStorage.SetDrop is expected once
func (m *mDropStorageMockSetDrop) ExpectOnce(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) *DropStorageMockSetDropExpectation {
	m.mock.SetDropFunc = nil
	m.mainExpectation = nil

	expectation := &DropStorageMockSetDropExpectation{}
	expectation.input = &DropStorageMockSetDropInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropStorageMockSetDropExpectation) Return(r error) {
	e.result = &DropStorageMockSetDropResult{r}
}

//Set uses given function f as a mock of DropStorage.SetDrop method
func (m *mDropStorageMockSetDrop) Set(f func(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) (r error)) *DropStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetDropFunc = f
	return m.mock
}

//SetDrop implements github.com/insolar/insolar/ledger/storage.DropStorage interface
func (m *DropStorageMock) SetDrop(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) (r error) {
	counter := atomic.AddUint64(&m.SetDropPreCounter, 1)
	defer atomic.AddUint64(&m.SetDropCounter, 1)

	if len(m.SetDropMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetDropMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropStorageMock.SetDrop. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetDropMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropStorageMockSetDropInput{p, p1, p2}, "DropStorage.SetDrop got unexpected parameters")

		result := m.SetDropMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.SetDrop")
			return
		}

		r = result.r

		return
	}

	if m.SetDropMock.mainExpectation != nil {

		input := m.SetDropMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropStorageMockSetDropInput{p, p1, p2}, "DropStorage.SetDrop got unexpected parameters")
		}

		result := m.SetDropMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.SetDrop")
		}

		r = result.r

		return
	}

	if m.SetDropFunc == nil {
		m.t.Fatalf("Unexpected call to DropStorageMock.SetDrop. %v %v %v", p, p1, p2)
		return
	}

	return m.SetDropFunc(p, p1, p2)
}

//SetDropMinimockCounter returns a count of DropStorageMock.SetDropFunc invocations
func (m *DropStorageMock) SetDropMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropCounter)
}

//SetDropMinimockPreCounter returns the value of DropStorageMock.SetDrop invocations
func (m *DropStorageMock) SetDropMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropPreCounter)
}

//SetDropFinished returns true if mock invocations count is ok
func (m *DropStorageMock) SetDropFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetDropMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetDropCounter) == uint64(len(m.SetDropMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetDropMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetDropCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetDropFunc != nil {
		return atomic.LoadUint64(&m.SetDropCounter) > 0
	}

	return true
}

type mDropStorageMockSetDropSizeHistory struct {
	mock              *DropStorageMock
	mainExpectation   *DropStorageMockSetDropSizeHistoryExpectation
	expectationSeries []*DropStorageMockSetDropSizeHistoryExpectation
}

type DropStorageMockSetDropSizeHistoryExpectation struct {
	input  *DropStorageMockSetDropSizeHistoryInput
	result *DropStorageMockSetDropSizeHistoryResult
}

type DropStorageMockSetDropSizeHistoryInput struct {
	p  context.Context
	p1 core.RecordID
	p2 jet.DropSizeHistory
}

type DropStorageMockSetDropSizeHistoryResult struct {
	r error
}

//Expect specifies that invocation of DropStorage.SetDropSizeHistory is expected from 1 to Infinity times
func (m *mDropStorageMockSetDropSizeHistory) Expect(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) *mDropStorageMockSetDropSizeHistory {
	m.mock.SetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockSetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.input = &DropStorageMockSetDropSizeHistoryInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DropStorage.SetDropSizeHistory
func (m *mDropStorageMockSetDropSizeHistory) Return(r error) *DropStorageMock {
	m.mock.SetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DropStorageMockSetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.result = &DropStorageMockSetDropSizeHistoryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DropStorage.SetDropSizeHistory is expected once
func (m *mDropStorageMockSetDropSizeHistory) ExpectOnce(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) *DropStorageMockSetDropSizeHistoryExpectation {
	m.mock.SetDropSizeHistoryFunc = nil
	m.mainExpectation = nil

	expectation := &DropStorageMockSetDropSizeHistoryExpectation{}
	expectation.input = &DropStorageMockSetDropSizeHistoryInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DropStorageMockSetDropSizeHistoryExpectation) Return(r error) {
	e.result = &DropStorageMockSetDropSizeHistoryResult{r}
}

//Set uses given function f as a mock of DropStorage.SetDropSizeHistory method
func (m *mDropStorageMockSetDropSizeHistory) Set(f func(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) (r error)) *DropStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetDropSizeHistoryFunc = f
	return m.mock
}

//SetDropSizeHistory implements github.com/insolar/insolar/ledger/storage.DropStorage interface
func (m *DropStorageMock) SetDropSizeHistory(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) (r error) {
	counter := atomic.AddUint64(&m.SetDropSizeHistoryPreCounter, 1)
	defer atomic.AddUint64(&m.SetDropSizeHistoryCounter, 1)

	if len(m.SetDropSizeHistoryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetDropSizeHistoryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DropStorageMock.SetDropSizeHistory. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetDropSizeHistoryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DropStorageMockSetDropSizeHistoryInput{p, p1, p2}, "DropStorage.SetDropSizeHistory got unexpected parameters")

		result := m.SetDropSizeHistoryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.SetDropSizeHistory")
			return
		}

		r = result.r

		return
	}

	if m.SetDropSizeHistoryMock.mainExpectation != nil {

		input := m.SetDropSizeHistoryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DropStorageMockSetDropSizeHistoryInput{p, p1, p2}, "DropStorage.SetDropSizeHistory got unexpected parameters")
		}

		result := m.SetDropSizeHistoryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DropStorageMock.SetDropSizeHistory")
		}

		r = result.r

		return
	}

	if m.SetDropSizeHistoryFunc == nil {
		m.t.Fatalf("Unexpected call to DropStorageMock.SetDropSizeHistory. %v %v %v", p, p1, p2)
		return
	}

	return m.SetDropSizeHistoryFunc(p, p1, p2)
}

//SetDropSizeHistoryMinimockCounter returns a count of DropStorageMock.SetDropSizeHistoryFunc invocations
func (m *DropStorageMock) SetDropSizeHistoryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropSizeHistoryCounter)
}

//SetDropSizeHistoryMinimockPreCounter returns the value of DropStorageMock.SetDropSizeHistory invocations
func (m *DropStorageMock) SetDropSizeHistoryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropSizeHistoryPreCounter)
}

//SetDropSizeHistoryFinished returns true if mock invocations count is ok
func (m *DropStorageMock) SetDropSizeHistoryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetDropSizeHistoryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetDropSizeHistoryCounter) == uint64(len(m.SetDropSizeHistoryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetDropSizeHistoryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetDropSizeHistoryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetDropSizeHistoryFunc != nil {
		return atomic.LoadUint64(&m.SetDropSizeHistoryCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DropStorageMock) ValidateCallCounters() {

	if !m.AddDropSizeFinished() {
		m.t.Fatal("Expected call to DropStorageMock.AddDropSize")
	}

	if !m.CreateDropFinished() {
		m.t.Fatal("Expected call to DropStorageMock.CreateDrop")
	}

	if !m.GetDropFinished() {
		m.t.Fatal("Expected call to DropStorageMock.GetDrop")
	}

	if !m.GetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to DropStorageMock.GetDropSizeHistory")
	}

	if !m.GetJetSizesHistoryDepthFinished() {
		m.t.Fatal("Expected call to DropStorageMock.GetJetSizesHistoryDepth")
	}

	if !m.SetDropFinished() {
		m.t.Fatal("Expected call to DropStorageMock.SetDrop")
	}

	if !m.SetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to DropStorageMock.SetDropSizeHistory")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DropStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DropStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DropStorageMock) MinimockFinish() {

	if !m.AddDropSizeFinished() {
		m.t.Fatal("Expected call to DropStorageMock.AddDropSize")
	}

	if !m.CreateDropFinished() {
		m.t.Fatal("Expected call to DropStorageMock.CreateDrop")
	}

	if !m.GetDropFinished() {
		m.t.Fatal("Expected call to DropStorageMock.GetDrop")
	}

	if !m.GetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to DropStorageMock.GetDropSizeHistory")
	}

	if !m.GetJetSizesHistoryDepthFinished() {
		m.t.Fatal("Expected call to DropStorageMock.GetJetSizesHistoryDepth")
	}

	if !m.SetDropFinished() {
		m.t.Fatal("Expected call to DropStorageMock.SetDrop")
	}

	if !m.SetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to DropStorageMock.SetDropSizeHistory")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DropStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DropStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddDropSizeFinished()
		ok = ok && m.CreateDropFinished()
		ok = ok && m.GetDropFinished()
		ok = ok && m.GetDropSizeHistoryFinished()
		ok = ok && m.GetJetSizesHistoryDepthFinished()
		ok = ok && m.SetDropFinished()
		ok = ok && m.SetDropSizeHistoryFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddDropSizeFinished() {
				m.t.Error("Expected call to DropStorageMock.AddDropSize")
			}

			if !m.CreateDropFinished() {
				m.t.Error("Expected call to DropStorageMock.CreateDrop")
			}

			if !m.GetDropFinished() {
				m.t.Error("Expected call to DropStorageMock.GetDrop")
			}

			if !m.GetDropSizeHistoryFinished() {
				m.t.Error("Expected call to DropStorageMock.GetDropSizeHistory")
			}

			if !m.GetJetSizesHistoryDepthFinished() {
				m.t.Error("Expected call to DropStorageMock.GetJetSizesHistoryDepth")
			}

			if !m.SetDropFinished() {
				m.t.Error("Expected call to DropStorageMock.SetDrop")
			}

			if !m.SetDropSizeHistoryFinished() {
				m.t.Error("Expected call to DropStorageMock.SetDropSizeHistory")
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
func (m *DropStorageMock) AllMocksCalled() bool {

	if !m.AddDropSizeFinished() {
		return false
	}

	if !m.CreateDropFinished() {
		return false
	}

	if !m.GetDropFinished() {
		return false
	}

	if !m.GetDropSizeHistoryFinished() {
		return false
	}

	if !m.GetJetSizesHistoryDepthFinished() {
		return false
	}

	if !m.SetDropFinished() {
		return false
	}

	if !m.SetDropSizeHistoryFinished() {
		return false
	}

	return true
}
