package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetStorage" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"

	testify_assert "github.com/stretchr/testify/assert"
)

// JetStorageMock implements github.com/insolar/insolar/ledger/storage.JetStorage
type JetStorageMock struct {
	t minimock.Tester

	AddDropSizeFunc       func(p context.Context, p1 *jet.DropSize) (r error)
	AddDropSizeCounter    uint64
	AddDropSizePreCounter uint64
	AddDropSizeMock       mJetStorageMockAddDropSize

	AddJetsFunc       func(p context.Context, p1 ...core.RecordID) (r error)
	AddJetsCounter    uint64
	AddJetsPreCounter uint64
	AddJetsMock       mJetStorageMockAddJets

	CloneJetTreeFunc       func(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) (r *jet.Tree, r1 error)
	CloneJetTreeCounter    uint64
	CloneJetTreePreCounter uint64
	CloneJetTreeMock       mJetStorageMockCloneJetTree

	CreateDropFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error)
	CreateDropCounter    uint64
	CreateDropPreCounter uint64
	CreateDropMock       mJetStorageMockCreateDrop

	GetDropFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *jet.JetDrop, r1 error)
	GetDropCounter    uint64
	GetDropPreCounter uint64
	GetDropMock       mJetStorageMockGetDrop

	GetDropSizeHistoryFunc       func(p context.Context, p1 core.RecordID) (r jet.DropSizeHistory, r1 error)
	GetDropSizeHistoryCounter    uint64
	GetDropSizeHistoryPreCounter uint64
	GetDropSizeHistoryMock       mJetStorageMockGetDropSizeHistory

	GetJetTreeFunc       func(p context.Context, p1 core.PulseNumber) (r *jet.Tree, r1 error)
	GetJetTreeCounter    uint64
	GetJetTreePreCounter uint64
	GetJetTreeMock       mJetStorageMockGetJetTree

	GetJetsFunc       func(p context.Context) (r jet.IDSet, r1 error)
	GetJetsCounter    uint64
	GetJetsPreCounter uint64
	GetJetsMock       mJetStorageMockGetJets

	SetDropFunc       func(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) (r error)
	SetDropCounter    uint64
	SetDropPreCounter uint64
	SetDropMock       mJetStorageMockSetDrop

	SetDropSizeHistoryFunc       func(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) (r error)
	SetDropSizeHistoryCounter    uint64
	SetDropSizeHistoryPreCounter uint64
	SetDropSizeHistoryMock       mJetStorageMockSetDropSizeHistory

	SplitJetTreeFunc       func(p context.Context, p1 core.PulseNumber, p2 core.RecordID) (r *core.RecordID, r1 *core.RecordID, r2 error)
	SplitJetTreeCounter    uint64
	SplitJetTreePreCounter uint64
	SplitJetTreeMock       mJetStorageMockSplitJetTree

	UpdateJetTreeFunc       func(p context.Context, p1 core.PulseNumber, p2 bool, p3 ...core.RecordID) (r error)
	UpdateJetTreeCounter    uint64
	UpdateJetTreePreCounter uint64
	UpdateJetTreeMock       mJetStorageMockUpdateJetTree
}

// NewJetStorageMock returns a mock for github.com/insolar/insolar/ledger/storage.JetStorage
func NewJetStorageMock(t minimock.Tester) *JetStorageMock {
	m := &JetStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddDropSizeMock = mJetStorageMockAddDropSize{mock: m}
	m.AddJetsMock = mJetStorageMockAddJets{mock: m}
	m.CloneJetTreeMock = mJetStorageMockCloneJetTree{mock: m}
	m.CreateDropMock = mJetStorageMockCreateDrop{mock: m}
	m.GetDropMock = mJetStorageMockGetDrop{mock: m}
	m.GetDropSizeHistoryMock = mJetStorageMockGetDropSizeHistory{mock: m}
	m.GetJetTreeMock = mJetStorageMockGetJetTree{mock: m}
	m.GetJetsMock = mJetStorageMockGetJets{mock: m}
	m.SetDropMock = mJetStorageMockSetDrop{mock: m}
	m.SetDropSizeHistoryMock = mJetStorageMockSetDropSizeHistory{mock: m}
	m.SplitJetTreeMock = mJetStorageMockSplitJetTree{mock: m}
	m.UpdateJetTreeMock = mJetStorageMockUpdateJetTree{mock: m}

	return m
}

type mJetStorageMockAddDropSize struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockAddDropSizeExpectation
	expectationSeries []*JetStorageMockAddDropSizeExpectation
}

type JetStorageMockAddDropSizeExpectation struct {
	input  *JetStorageMockAddDropSizeInput
	result *JetStorageMockAddDropSizeResult
}

type JetStorageMockAddDropSizeInput struct {
	p  context.Context
	p1 *jet.DropSize
}

type JetStorageMockAddDropSizeResult struct {
	r error
}

// Expect specifies that invocation of JetStorage.AddDropSize is expected from 1 to Infinity times
func (m *mJetStorageMockAddDropSize) Expect(p context.Context, p1 *jet.DropSize) *mJetStorageMockAddDropSize {
	m.mock.AddDropSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockAddDropSizeExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockAddDropSizeInput{p, p1}
	return m
}

// Return specifies results of invocation of JetStorage.AddDropSize
func (m *mJetStorageMockAddDropSize) Return(r error) *JetStorageMock {
	m.mock.AddDropSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockAddDropSizeExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockAddDropSizeResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.AddDropSize is expected once
func (m *mJetStorageMockAddDropSize) ExpectOnce(p context.Context, p1 *jet.DropSize) *JetStorageMockAddDropSizeExpectation {
	m.mock.AddDropSizeFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockAddDropSizeExpectation{}
	expectation.input = &JetStorageMockAddDropSizeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockAddDropSizeExpectation) Return(r error) {
	e.result = &JetStorageMockAddDropSizeResult{r}
}

// Set uses given function f as a mock of JetStorage.AddDropSize method
func (m *mJetStorageMockAddDropSize) Set(f func(p context.Context, p1 *jet.DropSize) (r error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddDropSizeFunc = f
	return m.mock
}

// AddDropSize implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) AddDropSize(p context.Context, p1 *jet.DropSize) (r error) {
	counter := atomic.AddUint64(&m.AddDropSizePreCounter, 1)
	defer atomic.AddUint64(&m.AddDropSizeCounter, 1)

	if len(m.AddDropSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddDropSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.AddDropSize. %v %v", p, p1)
			return
		}

		input := m.AddDropSizeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockAddDropSizeInput{p, p1}, "JetStorage.AddDropSize got unexpected parameters")

		result := m.AddDropSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.AddDropSize")
			return
		}

		r = result.r

		return
	}

	if m.AddDropSizeMock.mainExpectation != nil {

		input := m.AddDropSizeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockAddDropSizeInput{p, p1}, "JetStorage.AddDropSize got unexpected parameters")
		}

		result := m.AddDropSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.AddDropSize")
		}

		r = result.r

		return
	}

	if m.AddDropSizeFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.AddDropSize. %v %v", p, p1)
		return
	}

	return m.AddDropSizeFunc(p, p1)
}

// AddDropSizeMinimockCounter returns a count of JetStorageMock.AddDropSizeFunc invocations
func (m *JetStorageMock) AddDropSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddDropSizeCounter)
}

// AddDropSizeMinimockPreCounter returns the value of JetStorageMock.AddDropSize invocations
func (m *JetStorageMock) AddDropSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddDropSizePreCounter)
}

// AddDropSizeFinished returns true if mock invocations count is ok
func (m *JetStorageMock) AddDropSizeFinished() bool {
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

type mJetStorageMockAddJets struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockAddJetsExpectation
	expectationSeries []*JetStorageMockAddJetsExpectation
}

type JetStorageMockAddJetsExpectation struct {
	input  *JetStorageMockAddJetsInput
	result *JetStorageMockAddJetsResult
}

type JetStorageMockAddJetsInput struct {
	p  context.Context
	p1 []core.RecordID
}

type JetStorageMockAddJetsResult struct {
	r error
}

// Expect specifies that invocation of JetStorage.AddJets is expected from 1 to Infinity times
func (m *mJetStorageMockAddJets) Expect(p context.Context, p1 ...core.RecordID) *mJetStorageMockAddJets {
	m.mock.AddJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockAddJetsExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockAddJetsInput{p, p1}
	return m
}

// Return specifies results of invocation of JetStorage.AddJets
func (m *mJetStorageMockAddJets) Return(r error) *JetStorageMock {
	m.mock.AddJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockAddJetsExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockAddJetsResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.AddJets is expected once
func (m *mJetStorageMockAddJets) ExpectOnce(p context.Context, p1 ...core.RecordID) *JetStorageMockAddJetsExpectation {
	m.mock.AddJetsFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockAddJetsExpectation{}
	expectation.input = &JetStorageMockAddJetsInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockAddJetsExpectation) Return(r error) {
	e.result = &JetStorageMockAddJetsResult{r}
}

// Set uses given function f as a mock of JetStorage.AddJets method
func (m *mJetStorageMockAddJets) Set(f func(p context.Context, p1 ...core.RecordID) (r error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddJetsFunc = f
	return m.mock
}

// AddJets implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) AddJets(p context.Context, p1 ...core.RecordID) (r error) {
	counter := atomic.AddUint64(&m.AddJetsPreCounter, 1)
	defer atomic.AddUint64(&m.AddJetsCounter, 1)

	if len(m.AddJetsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddJetsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.AddJets. %v %v", p, p1)
			return
		}

		input := m.AddJetsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockAddJetsInput{p, p1}, "JetStorage.AddJets got unexpected parameters")

		result := m.AddJetsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.AddJets")
			return
		}

		r = result.r

		return
	}

	if m.AddJetsMock.mainExpectation != nil {

		input := m.AddJetsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockAddJetsInput{p, p1}, "JetStorage.AddJets got unexpected parameters")
		}

		result := m.AddJetsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.AddJets")
		}

		r = result.r

		return
	}

	if m.AddJetsFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.AddJets. %v %v", p, p1)
		return
	}

	return m.AddJetsFunc(p, p1...)
}

// AddJetsMinimockCounter returns a count of JetStorageMock.AddJetsFunc invocations
func (m *JetStorageMock) AddJetsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddJetsCounter)
}

// AddJetsMinimockPreCounter returns the value of JetStorageMock.AddJets invocations
func (m *JetStorageMock) AddJetsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddJetsPreCounter)
}

// AddJetsFinished returns true if mock invocations count is ok
func (m *JetStorageMock) AddJetsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddJetsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddJetsCounter) == uint64(len(m.AddJetsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddJetsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddJetsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddJetsFunc != nil {
		return atomic.LoadUint64(&m.AddJetsCounter) > 0
	}

	return true
}

type mJetStorageMockCloneJetTree struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockCloneJetTreeExpectation
	expectationSeries []*JetStorageMockCloneJetTreeExpectation
}

type JetStorageMockCloneJetTreeExpectation struct {
	input  *JetStorageMockCloneJetTreeInput
	result *JetStorageMockCloneJetTreeResult
}

type JetStorageMockCloneJetTreeInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 core.PulseNumber
}

type JetStorageMockCloneJetTreeResult struct {
	r  *jet.Tree
	r1 error
}

// Expect specifies that invocation of JetStorage.CloneJetTree is expected from 1 to Infinity times
func (m *mJetStorageMockCloneJetTree) Expect(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) *mJetStorageMockCloneJetTree {
	m.mock.CloneJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockCloneJetTreeExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockCloneJetTreeInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of JetStorage.CloneJetTree
func (m *mJetStorageMockCloneJetTree) Return(r *jet.Tree, r1 error) *JetStorageMock {
	m.mock.CloneJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockCloneJetTreeExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockCloneJetTreeResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.CloneJetTree is expected once
func (m *mJetStorageMockCloneJetTree) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) *JetStorageMockCloneJetTreeExpectation {
	m.mock.CloneJetTreeFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockCloneJetTreeExpectation{}
	expectation.input = &JetStorageMockCloneJetTreeInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockCloneJetTreeExpectation) Return(r *jet.Tree, r1 error) {
	e.result = &JetStorageMockCloneJetTreeResult{r, r1}
}

// Set uses given function f as a mock of JetStorage.CloneJetTree method
func (m *mJetStorageMockCloneJetTree) Set(f func(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) (r *jet.Tree, r1 error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloneJetTreeFunc = f
	return m.mock
}

// CloneJetTree implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) CloneJetTree(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) (r *jet.Tree, r1 error) {
	counter := atomic.AddUint64(&m.CloneJetTreePreCounter, 1)
	defer atomic.AddUint64(&m.CloneJetTreeCounter, 1)

	if len(m.CloneJetTreeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloneJetTreeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.CloneJetTree. %v %v %v", p, p1, p2)
			return
		}

		input := m.CloneJetTreeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockCloneJetTreeInput{p, p1, p2}, "JetStorage.CloneJetTree got unexpected parameters")

		result := m.CloneJetTreeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.CloneJetTree")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CloneJetTreeMock.mainExpectation != nil {

		input := m.CloneJetTreeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockCloneJetTreeInput{p, p1, p2}, "JetStorage.CloneJetTree got unexpected parameters")
		}

		result := m.CloneJetTreeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.CloneJetTree")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CloneJetTreeFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.CloneJetTree. %v %v %v", p, p1, p2)
		return
	}

	return m.CloneJetTreeFunc(p, p1, p2)
}

// CloneJetTreeMinimockCounter returns a count of JetStorageMock.CloneJetTreeFunc invocations
func (m *JetStorageMock) CloneJetTreeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloneJetTreeCounter)
}

// CloneJetTreeMinimockPreCounter returns the value of JetStorageMock.CloneJetTree invocations
func (m *JetStorageMock) CloneJetTreeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CloneJetTreePreCounter)
}

// CloneJetTreeFinished returns true if mock invocations count is ok
func (m *JetStorageMock) CloneJetTreeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloneJetTreeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloneJetTreeCounter) == uint64(len(m.CloneJetTreeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloneJetTreeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloneJetTreeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloneJetTreeFunc != nil {
		return atomic.LoadUint64(&m.CloneJetTreeCounter) > 0
	}

	return true
}

type mJetStorageMockCreateDrop struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockCreateDropExpectation
	expectationSeries []*JetStorageMockCreateDropExpectation
}

type JetStorageMockCreateDropExpectation struct {
	input  *JetStorageMockCreateDropInput
	result *JetStorageMockCreateDropResult
}

type JetStorageMockCreateDropInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 []byte
}

type JetStorageMockCreateDropResult struct {
	r  *jet.JetDrop
	r1 [][]byte
	r2 uint64
	r3 error
}

// Expect specifies that invocation of JetStorage.CreateDrop is expected from 1 to Infinity times
func (m *mJetStorageMockCreateDrop) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) *mJetStorageMockCreateDrop {
	m.mock.CreateDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockCreateDropExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockCreateDropInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of JetStorage.CreateDrop
func (m *mJetStorageMockCreateDrop) Return(r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error) *JetStorageMock {
	m.mock.CreateDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockCreateDropExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockCreateDropResult{r, r1, r2, r3}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.CreateDrop is expected once
func (m *mJetStorageMockCreateDrop) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) *JetStorageMockCreateDropExpectation {
	m.mock.CreateDropFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockCreateDropExpectation{}
	expectation.input = &JetStorageMockCreateDropInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockCreateDropExpectation) Return(r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error) {
	e.result = &JetStorageMockCreateDropResult{r, r1, r2, r3}
}

// Set uses given function f as a mock of JetStorage.CreateDrop method
func (m *mJetStorageMockCreateDrop) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateDropFunc = f
	return m.mock
}

// CreateDrop implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) CreateDrop(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 []byte) (r *jet.JetDrop, r1 [][]byte, r2 uint64, r3 error) {
	counter := atomic.AddUint64(&m.CreateDropPreCounter, 1)
	defer atomic.AddUint64(&m.CreateDropCounter, 1)

	if len(m.CreateDropMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateDropMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.CreateDrop. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.CreateDropMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockCreateDropInput{p, p1, p2, p3}, "JetStorage.CreateDrop got unexpected parameters")

		result := m.CreateDropMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.CreateDrop")
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
			testify_assert.Equal(m.t, *input, JetStorageMockCreateDropInput{p, p1, p2, p3}, "JetStorage.CreateDrop got unexpected parameters")
		}

		result := m.CreateDropMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.CreateDrop")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2
		r3 = result.r3

		return
	}

	if m.CreateDropFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.CreateDrop. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.CreateDropFunc(p, p1, p2, p3)
}

// CreateDropMinimockCounter returns a count of JetStorageMock.CreateDropFunc invocations
func (m *JetStorageMock) CreateDropMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateDropCounter)
}

// CreateDropMinimockPreCounter returns the value of JetStorageMock.CreateDrop invocations
func (m *JetStorageMock) CreateDropMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateDropPreCounter)
}

// CreateDropFinished returns true if mock invocations count is ok
func (m *JetStorageMock) CreateDropFinished() bool {
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

type mJetStorageMockGetDrop struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockGetDropExpectation
	expectationSeries []*JetStorageMockGetDropExpectation
}

type JetStorageMockGetDropExpectation struct {
	input  *JetStorageMockGetDropInput
	result *JetStorageMockGetDropResult
}

type JetStorageMockGetDropInput struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

type JetStorageMockGetDropResult struct {
	r  *jet.JetDrop
	r1 error
}

// Expect specifies that invocation of JetStorage.GetDrop is expected from 1 to Infinity times
func (m *mJetStorageMockGetDrop) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mJetStorageMockGetDrop {
	m.mock.GetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetDropExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockGetDropInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of JetStorage.GetDrop
func (m *mJetStorageMockGetDrop) Return(r *jet.JetDrop, r1 error) *JetStorageMock {
	m.mock.GetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetDropExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockGetDropResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.GetDrop is expected once
func (m *mJetStorageMockGetDrop) ExpectOnce(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *JetStorageMockGetDropExpectation {
	m.mock.GetDropFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockGetDropExpectation{}
	expectation.input = &JetStorageMockGetDropInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockGetDropExpectation) Return(r *jet.JetDrop, r1 error) {
	e.result = &JetStorageMockGetDropResult{r, r1}
}

// Set uses given function f as a mock of JetStorage.GetDrop method
func (m *mJetStorageMockGetDrop) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *jet.JetDrop, r1 error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDropFunc = f
	return m.mock
}

// GetDrop implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) GetDrop(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *jet.JetDrop, r1 error) {
	counter := atomic.AddUint64(&m.GetDropPreCounter, 1)
	defer atomic.AddUint64(&m.GetDropCounter, 1)

	if len(m.GetDropMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDropMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.GetDrop. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetDropMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockGetDropInput{p, p1, p2}, "JetStorage.GetDrop got unexpected parameters")

		result := m.GetDropMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetDrop")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropMock.mainExpectation != nil {

		input := m.GetDropMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockGetDropInput{p, p1, p2}, "JetStorage.GetDrop got unexpected parameters")
		}

		result := m.GetDropMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetDrop")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.GetDrop. %v %v %v", p, p1, p2)
		return
	}

	return m.GetDropFunc(p, p1, p2)
}

// GetDropMinimockCounter returns a count of JetStorageMock.GetDropFunc invocations
func (m *JetStorageMock) GetDropMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropCounter)
}

// GetDropMinimockPreCounter returns the value of JetStorageMock.GetDrop invocations
func (m *JetStorageMock) GetDropMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropPreCounter)
}

// GetDropFinished returns true if mock invocations count is ok
func (m *JetStorageMock) GetDropFinished() bool {
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

type mJetStorageMockGetDropSizeHistory struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockGetDropSizeHistoryExpectation
	expectationSeries []*JetStorageMockGetDropSizeHistoryExpectation
}

type JetStorageMockGetDropSizeHistoryExpectation struct {
	input  *JetStorageMockGetDropSizeHistoryInput
	result *JetStorageMockGetDropSizeHistoryResult
}

type JetStorageMockGetDropSizeHistoryInput struct {
	p  context.Context
	p1 core.RecordID
}

type JetStorageMockGetDropSizeHistoryResult struct {
	r  jet.DropSizeHistory
	r1 error
}

// Expect specifies that invocation of JetStorage.GetDropSizeHistory is expected from 1 to Infinity times
func (m *mJetStorageMockGetDropSizeHistory) Expect(p context.Context, p1 core.RecordID) *mJetStorageMockGetDropSizeHistory {
	m.mock.GetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockGetDropSizeHistoryInput{p, p1}
	return m
}

// Return specifies results of invocation of JetStorage.GetDropSizeHistory
func (m *mJetStorageMockGetDropSizeHistory) Return(r jet.DropSizeHistory, r1 error) *JetStorageMock {
	m.mock.GetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockGetDropSizeHistoryResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.GetDropSizeHistory is expected once
func (m *mJetStorageMockGetDropSizeHistory) ExpectOnce(p context.Context, p1 core.RecordID) *JetStorageMockGetDropSizeHistoryExpectation {
	m.mock.GetDropSizeHistoryFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockGetDropSizeHistoryExpectation{}
	expectation.input = &JetStorageMockGetDropSizeHistoryInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockGetDropSizeHistoryExpectation) Return(r jet.DropSizeHistory, r1 error) {
	e.result = &JetStorageMockGetDropSizeHistoryResult{r, r1}
}

// Set uses given function f as a mock of JetStorage.GetDropSizeHistory method
func (m *mJetStorageMockGetDropSizeHistory) Set(f func(p context.Context, p1 core.RecordID) (r jet.DropSizeHistory, r1 error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDropSizeHistoryFunc = f
	return m.mock
}

// GetDropSizeHistory implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) GetDropSizeHistory(p context.Context, p1 core.RecordID) (r jet.DropSizeHistory, r1 error) {
	counter := atomic.AddUint64(&m.GetDropSizeHistoryPreCounter, 1)
	defer atomic.AddUint64(&m.GetDropSizeHistoryCounter, 1)

	if len(m.GetDropSizeHistoryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDropSizeHistoryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.GetDropSizeHistory. %v %v", p, p1)
			return
		}

		input := m.GetDropSizeHistoryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockGetDropSizeHistoryInput{p, p1}, "JetStorage.GetDropSizeHistory got unexpected parameters")

		result := m.GetDropSizeHistoryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetDropSizeHistory")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropSizeHistoryMock.mainExpectation != nil {

		input := m.GetDropSizeHistoryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockGetDropSizeHistoryInput{p, p1}, "JetStorage.GetDropSizeHistory got unexpected parameters")
		}

		result := m.GetDropSizeHistoryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetDropSizeHistory")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDropSizeHistoryFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.GetDropSizeHistory. %v %v", p, p1)
		return
	}

	return m.GetDropSizeHistoryFunc(p, p1)
}

// GetDropSizeHistoryMinimockCounter returns a count of JetStorageMock.GetDropSizeHistoryFunc invocations
func (m *JetStorageMock) GetDropSizeHistoryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropSizeHistoryCounter)
}

// GetDropSizeHistoryMinimockPreCounter returns the value of JetStorageMock.GetDropSizeHistory invocations
func (m *JetStorageMock) GetDropSizeHistoryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDropSizeHistoryPreCounter)
}

// GetDropSizeHistoryFinished returns true if mock invocations count is ok
func (m *JetStorageMock) GetDropSizeHistoryFinished() bool {
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

type mJetStorageMockGetJetTree struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockGetJetTreeExpectation
	expectationSeries []*JetStorageMockGetJetTreeExpectation
}

type JetStorageMockGetJetTreeExpectation struct {
	input  *JetStorageMockGetJetTreeInput
	result *JetStorageMockGetJetTreeResult
}

type JetStorageMockGetJetTreeInput struct {
	p  context.Context
	p1 core.PulseNumber
}

type JetStorageMockGetJetTreeResult struct {
	r  *jet.Tree
	r1 error
}

// Expect specifies that invocation of JetStorage.GetJetTree is expected from 1 to Infinity times
func (m *mJetStorageMockGetJetTree) Expect(p context.Context, p1 core.PulseNumber) *mJetStorageMockGetJetTree {
	m.mock.GetJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetJetTreeExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockGetJetTreeInput{p, p1}
	return m
}

// Return specifies results of invocation of JetStorage.GetJetTree
func (m *mJetStorageMockGetJetTree) Return(r *jet.Tree, r1 error) *JetStorageMock {
	m.mock.GetJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetJetTreeExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockGetJetTreeResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.GetJetTree is expected once
func (m *mJetStorageMockGetJetTree) ExpectOnce(p context.Context, p1 core.PulseNumber) *JetStorageMockGetJetTreeExpectation {
	m.mock.GetJetTreeFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockGetJetTreeExpectation{}
	expectation.input = &JetStorageMockGetJetTreeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockGetJetTreeExpectation) Return(r *jet.Tree, r1 error) {
	e.result = &JetStorageMockGetJetTreeResult{r, r1}
}

// Set uses given function f as a mock of JetStorage.GetJetTree method
func (m *mJetStorageMockGetJetTree) Set(f func(p context.Context, p1 core.PulseNumber) (r *jet.Tree, r1 error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetJetTreeFunc = f
	return m.mock
}

// GetJetTree implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) GetJetTree(p context.Context, p1 core.PulseNumber) (r *jet.Tree, r1 error) {
	counter := atomic.AddUint64(&m.GetJetTreePreCounter, 1)
	defer atomic.AddUint64(&m.GetJetTreeCounter, 1)

	if len(m.GetJetTreeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetJetTreeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.GetJetTree. %v %v", p, p1)
			return
		}

		input := m.GetJetTreeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockGetJetTreeInput{p, p1}, "JetStorage.GetJetTree got unexpected parameters")

		result := m.GetJetTreeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetJetTree")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetJetTreeMock.mainExpectation != nil {

		input := m.GetJetTreeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockGetJetTreeInput{p, p1}, "JetStorage.GetJetTree got unexpected parameters")
		}

		result := m.GetJetTreeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetJetTree")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetJetTreeFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.GetJetTree. %v %v", p, p1)
		return
	}

	return m.GetJetTreeFunc(p, p1)
}

// GetJetTreeMinimockCounter returns a count of JetStorageMock.GetJetTreeFunc invocations
func (m *JetStorageMock) GetJetTreeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetJetTreeCounter)
}

// GetJetTreeMinimockPreCounter returns the value of JetStorageMock.GetJetTree invocations
func (m *JetStorageMock) GetJetTreeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetJetTreePreCounter)
}

// GetJetTreeFinished returns true if mock invocations count is ok
func (m *JetStorageMock) GetJetTreeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetJetTreeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetJetTreeCounter) == uint64(len(m.GetJetTreeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetJetTreeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetJetTreeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetJetTreeFunc != nil {
		return atomic.LoadUint64(&m.GetJetTreeCounter) > 0
	}

	return true
}

type mJetStorageMockGetJets struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockGetJetsExpectation
	expectationSeries []*JetStorageMockGetJetsExpectation
}

type JetStorageMockGetJetsExpectation struct {
	input  *JetStorageMockGetJetsInput
	result *JetStorageMockGetJetsResult
}

type JetStorageMockGetJetsInput struct {
	p context.Context
}

type JetStorageMockGetJetsResult struct {
	r  jet.IDSet
	r1 error
}

// Expect specifies that invocation of JetStorage.GetJets is expected from 1 to Infinity times
func (m *mJetStorageMockGetJets) Expect(p context.Context) *mJetStorageMockGetJets {
	m.mock.GetJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetJetsExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockGetJetsInput{p}
	return m
}

// Return specifies results of invocation of JetStorage.GetJets
func (m *mJetStorageMockGetJets) Return(r jet.IDSet, r1 error) *JetStorageMock {
	m.mock.GetJetsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockGetJetsExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockGetJetsResult{r, r1}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.GetJets is expected once
func (m *mJetStorageMockGetJets) ExpectOnce(p context.Context) *JetStorageMockGetJetsExpectation {
	m.mock.GetJetsFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockGetJetsExpectation{}
	expectation.input = &JetStorageMockGetJetsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockGetJetsExpectation) Return(r jet.IDSet, r1 error) {
	e.result = &JetStorageMockGetJetsResult{r, r1}
}

// Set uses given function f as a mock of JetStorage.GetJets method
func (m *mJetStorageMockGetJets) Set(f func(p context.Context) (r jet.IDSet, r1 error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetJetsFunc = f
	return m.mock
}

// GetJets implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) GetJets(p context.Context) (r jet.IDSet, r1 error) {
	counter := atomic.AddUint64(&m.GetJetsPreCounter, 1)
	defer atomic.AddUint64(&m.GetJetsCounter, 1)

	if len(m.GetJetsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetJetsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.GetJets. %v", p)
			return
		}

		input := m.GetJetsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockGetJetsInput{p}, "JetStorage.GetJets got unexpected parameters")

		result := m.GetJetsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetJets")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetJetsMock.mainExpectation != nil {

		input := m.GetJetsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockGetJetsInput{p}, "JetStorage.GetJets got unexpected parameters")
		}

		result := m.GetJetsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.GetJets")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetJetsFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.GetJets. %v", p)
		return
	}

	return m.GetJetsFunc(p)
}

// GetJetsMinimockCounter returns a count of JetStorageMock.GetJetsFunc invocations
func (m *JetStorageMock) GetJetsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetJetsCounter)
}

// GetJetsMinimockPreCounter returns the value of JetStorageMock.GetJets invocations
func (m *JetStorageMock) GetJetsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetJetsPreCounter)
}

// GetJetsFinished returns true if mock invocations count is ok
func (m *JetStorageMock) GetJetsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetJetsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetJetsCounter) == uint64(len(m.GetJetsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetJetsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetJetsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetJetsFunc != nil {
		return atomic.LoadUint64(&m.GetJetsCounter) > 0
	}

	return true
}

type mJetStorageMockSetDrop struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockSetDropExpectation
	expectationSeries []*JetStorageMockSetDropExpectation
}

type JetStorageMockSetDropExpectation struct {
	input  *JetStorageMockSetDropInput
	result *JetStorageMockSetDropResult
}

type JetStorageMockSetDropInput struct {
	p  context.Context
	p1 core.RecordID
	p2 *jet.JetDrop
}

type JetStorageMockSetDropResult struct {
	r error
}

// Expect specifies that invocation of JetStorage.SetDrop is expected from 1 to Infinity times
func (m *mJetStorageMockSetDrop) Expect(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) *mJetStorageMockSetDrop {
	m.mock.SetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockSetDropExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockSetDropInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of JetStorage.SetDrop
func (m *mJetStorageMockSetDrop) Return(r error) *JetStorageMock {
	m.mock.SetDropFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockSetDropExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockSetDropResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.SetDrop is expected once
func (m *mJetStorageMockSetDrop) ExpectOnce(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) *JetStorageMockSetDropExpectation {
	m.mock.SetDropFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockSetDropExpectation{}
	expectation.input = &JetStorageMockSetDropInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockSetDropExpectation) Return(r error) {
	e.result = &JetStorageMockSetDropResult{r}
}

// Set uses given function f as a mock of JetStorage.SetDrop method
func (m *mJetStorageMockSetDrop) Set(f func(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) (r error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetDropFunc = f
	return m.mock
}

// SetDrop implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) SetDrop(p context.Context, p1 core.RecordID, p2 *jet.JetDrop) (r error) {
	counter := atomic.AddUint64(&m.SetDropPreCounter, 1)
	defer atomic.AddUint64(&m.SetDropCounter, 1)

	if len(m.SetDropMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetDropMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.SetDrop. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetDropMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockSetDropInput{p, p1, p2}, "JetStorage.SetDrop got unexpected parameters")

		result := m.SetDropMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.SetDrop")
			return
		}

		r = result.r

		return
	}

	if m.SetDropMock.mainExpectation != nil {

		input := m.SetDropMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockSetDropInput{p, p1, p2}, "JetStorage.SetDrop got unexpected parameters")
		}

		result := m.SetDropMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.SetDrop")
		}

		r = result.r

		return
	}

	if m.SetDropFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.SetDrop. %v %v %v", p, p1, p2)
		return
	}

	return m.SetDropFunc(p, p1, p2)
}

// SetDropMinimockCounter returns a count of JetStorageMock.SetDropFunc invocations
func (m *JetStorageMock) SetDropMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropCounter)
}

// SetDropMinimockPreCounter returns the value of JetStorageMock.SetDrop invocations
func (m *JetStorageMock) SetDropMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropPreCounter)
}

// SetDropFinished returns true if mock invocations count is ok
func (m *JetStorageMock) SetDropFinished() bool {
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

type mJetStorageMockSetDropSizeHistory struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockSetDropSizeHistoryExpectation
	expectationSeries []*JetStorageMockSetDropSizeHistoryExpectation
}

type JetStorageMockSetDropSizeHistoryExpectation struct {
	input  *JetStorageMockSetDropSizeHistoryInput
	result *JetStorageMockSetDropSizeHistoryResult
}

type JetStorageMockSetDropSizeHistoryInput struct {
	p  context.Context
	p1 core.RecordID
	p2 jet.DropSizeHistory
}

type JetStorageMockSetDropSizeHistoryResult struct {
	r error
}

// Expect specifies that invocation of JetStorage.SetDropSizeHistory is expected from 1 to Infinity times
func (m *mJetStorageMockSetDropSizeHistory) Expect(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) *mJetStorageMockSetDropSizeHistory {
	m.mock.SetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockSetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockSetDropSizeHistoryInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of JetStorage.SetDropSizeHistory
func (m *mJetStorageMockSetDropSizeHistory) Return(r error) *JetStorageMock {
	m.mock.SetDropSizeHistoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockSetDropSizeHistoryExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockSetDropSizeHistoryResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.SetDropSizeHistory is expected once
func (m *mJetStorageMockSetDropSizeHistory) ExpectOnce(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) *JetStorageMockSetDropSizeHistoryExpectation {
	m.mock.SetDropSizeHistoryFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockSetDropSizeHistoryExpectation{}
	expectation.input = &JetStorageMockSetDropSizeHistoryInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockSetDropSizeHistoryExpectation) Return(r error) {
	e.result = &JetStorageMockSetDropSizeHistoryResult{r}
}

// Set uses given function f as a mock of JetStorage.SetDropSizeHistory method
func (m *mJetStorageMockSetDropSizeHistory) Set(f func(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) (r error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetDropSizeHistoryFunc = f
	return m.mock
}

// SetDropSizeHistory implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) SetDropSizeHistory(p context.Context, p1 core.RecordID, p2 jet.DropSizeHistory) (r error) {
	counter := atomic.AddUint64(&m.SetDropSizeHistoryPreCounter, 1)
	defer atomic.AddUint64(&m.SetDropSizeHistoryCounter, 1)

	if len(m.SetDropSizeHistoryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetDropSizeHistoryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.SetDropSizeHistory. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetDropSizeHistoryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockSetDropSizeHistoryInput{p, p1, p2}, "JetStorage.SetDropSizeHistory got unexpected parameters")

		result := m.SetDropSizeHistoryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.SetDropSizeHistory")
			return
		}

		r = result.r

		return
	}

	if m.SetDropSizeHistoryMock.mainExpectation != nil {

		input := m.SetDropSizeHistoryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockSetDropSizeHistoryInput{p, p1, p2}, "JetStorage.SetDropSizeHistory got unexpected parameters")
		}

		result := m.SetDropSizeHistoryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.SetDropSizeHistory")
		}

		r = result.r

		return
	}

	if m.SetDropSizeHistoryFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.SetDropSizeHistory. %v %v %v", p, p1, p2)
		return
	}

	return m.SetDropSizeHistoryFunc(p, p1, p2)
}

// SetDropSizeHistoryMinimockCounter returns a count of JetStorageMock.SetDropSizeHistoryFunc invocations
func (m *JetStorageMock) SetDropSizeHistoryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropSizeHistoryCounter)
}

// SetDropSizeHistoryMinimockPreCounter returns the value of JetStorageMock.SetDropSizeHistory invocations
func (m *JetStorageMock) SetDropSizeHistoryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetDropSizeHistoryPreCounter)
}

// SetDropSizeHistoryFinished returns true if mock invocations count is ok
func (m *JetStorageMock) SetDropSizeHistoryFinished() bool {
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

type mJetStorageMockSplitJetTree struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockSplitJetTreeExpectation
	expectationSeries []*JetStorageMockSplitJetTreeExpectation
}

type JetStorageMockSplitJetTreeExpectation struct {
	input  *JetStorageMockSplitJetTreeInput
	result *JetStorageMockSplitJetTreeResult
}

type JetStorageMockSplitJetTreeInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 core.RecordID
}

type JetStorageMockSplitJetTreeResult struct {
	r  *core.RecordID
	r1 *core.RecordID
	r2 error
}

// Expect specifies that invocation of JetStorage.SplitJetTree is expected from 1 to Infinity times
func (m *mJetStorageMockSplitJetTree) Expect(p context.Context, p1 core.PulseNumber, p2 core.RecordID) *mJetStorageMockSplitJetTree {
	m.mock.SplitJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockSplitJetTreeExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockSplitJetTreeInput{p, p1, p2}
	return m
}

// Return specifies results of invocation of JetStorage.SplitJetTree
func (m *mJetStorageMockSplitJetTree) Return(r *core.RecordID, r1 *core.RecordID, r2 error) *JetStorageMock {
	m.mock.SplitJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockSplitJetTreeExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockSplitJetTreeResult{r, r1, r2}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.SplitJetTree is expected once
func (m *mJetStorageMockSplitJetTree) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 core.RecordID) *JetStorageMockSplitJetTreeExpectation {
	m.mock.SplitJetTreeFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockSplitJetTreeExpectation{}
	expectation.input = &JetStorageMockSplitJetTreeInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockSplitJetTreeExpectation) Return(r *core.RecordID, r1 *core.RecordID, r2 error) {
	e.result = &JetStorageMockSplitJetTreeResult{r, r1, r2}
}

// Set uses given function f as a mock of JetStorage.SplitJetTree method
func (m *mJetStorageMockSplitJetTree) Set(f func(p context.Context, p1 core.PulseNumber, p2 core.RecordID) (r *core.RecordID, r1 *core.RecordID, r2 error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SplitJetTreeFunc = f
	return m.mock
}

// SplitJetTree implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) SplitJetTree(p context.Context, p1 core.PulseNumber, p2 core.RecordID) (r *core.RecordID, r1 *core.RecordID, r2 error) {
	counter := atomic.AddUint64(&m.SplitJetTreePreCounter, 1)
	defer atomic.AddUint64(&m.SplitJetTreeCounter, 1)

	if len(m.SplitJetTreeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SplitJetTreeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.SplitJetTree. %v %v %v", p, p1, p2)
			return
		}

		input := m.SplitJetTreeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockSplitJetTreeInput{p, p1, p2}, "JetStorage.SplitJetTree got unexpected parameters")

		result := m.SplitJetTreeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.SplitJetTree")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.SplitJetTreeMock.mainExpectation != nil {

		input := m.SplitJetTreeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockSplitJetTreeInput{p, p1, p2}, "JetStorage.SplitJetTree got unexpected parameters")
		}

		result := m.SplitJetTreeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.SplitJetTree")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.SplitJetTreeFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.SplitJetTree. %v %v %v", p, p1, p2)
		return
	}

	return m.SplitJetTreeFunc(p, p1, p2)
}

// SplitJetTreeMinimockCounter returns a count of JetStorageMock.SplitJetTreeFunc invocations
func (m *JetStorageMock) SplitJetTreeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SplitJetTreeCounter)
}

// SplitJetTreeMinimockPreCounter returns the value of JetStorageMock.SplitJetTree invocations
func (m *JetStorageMock) SplitJetTreeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SplitJetTreePreCounter)
}

// SplitJetTreeFinished returns true if mock invocations count is ok
func (m *JetStorageMock) SplitJetTreeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SplitJetTreeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SplitJetTreeCounter) == uint64(len(m.SplitJetTreeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SplitJetTreeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SplitJetTreeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SplitJetTreeFunc != nil {
		return atomic.LoadUint64(&m.SplitJetTreeCounter) > 0
	}

	return true
}

type mJetStorageMockUpdateJetTree struct {
	mock              *JetStorageMock
	mainExpectation   *JetStorageMockUpdateJetTreeExpectation
	expectationSeries []*JetStorageMockUpdateJetTreeExpectation
}

type JetStorageMockUpdateJetTreeExpectation struct {
	input  *JetStorageMockUpdateJetTreeInput
	result *JetStorageMockUpdateJetTreeResult
}

type JetStorageMockUpdateJetTreeInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 bool
	p3 []core.RecordID
}

type JetStorageMockUpdateJetTreeResult struct {
	r error
}

// Expect specifies that invocation of JetStorage.UpdateJetTree is expected from 1 to Infinity times
func (m *mJetStorageMockUpdateJetTree) Expect(p context.Context, p1 core.PulseNumber, p2 bool, p3 ...core.RecordID) *mJetStorageMockUpdateJetTree {
	m.mock.UpdateJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockUpdateJetTreeExpectation{}
	}
	m.mainExpectation.input = &JetStorageMockUpdateJetTreeInput{p, p1, p2, p3}
	return m
}

// Return specifies results of invocation of JetStorage.UpdateJetTree
func (m *mJetStorageMockUpdateJetTree) Return(r error) *JetStorageMock {
	m.mock.UpdateJetTreeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetStorageMockUpdateJetTreeExpectation{}
	}
	m.mainExpectation.result = &JetStorageMockUpdateJetTreeResult{r}
	return m.mock
}

// ExpectOnce specifies that invocation of JetStorage.UpdateJetTree is expected once
func (m *mJetStorageMockUpdateJetTree) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 bool, p3 ...core.RecordID) *JetStorageMockUpdateJetTreeExpectation {
	m.mock.UpdateJetTreeFunc = nil
	m.mainExpectation = nil

	expectation := &JetStorageMockUpdateJetTreeExpectation{}
	expectation.input = &JetStorageMockUpdateJetTreeInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetStorageMockUpdateJetTreeExpectation) Return(r error) {
	e.result = &JetStorageMockUpdateJetTreeResult{r}
}

// Set uses given function f as a mock of JetStorage.UpdateJetTree method
func (m *mJetStorageMockUpdateJetTree) Set(f func(p context.Context, p1 core.PulseNumber, p2 bool, p3 ...core.RecordID) (r error)) *JetStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateJetTreeFunc = f
	return m.mock
}

// UpdateJetTree implements github.com/insolar/insolar/ledger/storage.JetStorage interface
func (m *JetStorageMock) UpdateJetTree(p context.Context, p1 core.PulseNumber, p2 bool, p3 ...core.RecordID) (r error) {
	counter := atomic.AddUint64(&m.UpdateJetTreePreCounter, 1)
	defer atomic.AddUint64(&m.UpdateJetTreeCounter, 1)

	if len(m.UpdateJetTreeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateJetTreeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetStorageMock.UpdateJetTree. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.UpdateJetTreeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetStorageMockUpdateJetTreeInput{p, p1, p2, p3}, "JetStorage.UpdateJetTree got unexpected parameters")

		result := m.UpdateJetTreeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.UpdateJetTree")
			return
		}

		r = result.r

		return
	}

	if m.UpdateJetTreeMock.mainExpectation != nil {

		input := m.UpdateJetTreeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetStorageMockUpdateJetTreeInput{p, p1, p2, p3}, "JetStorage.UpdateJetTree got unexpected parameters")
		}

		result := m.UpdateJetTreeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetStorageMock.UpdateJetTree")
		}

		r = result.r

		return
	}

	if m.UpdateJetTreeFunc == nil {
		m.t.Fatalf("Unexpected call to JetStorageMock.UpdateJetTree. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.UpdateJetTreeFunc(p, p1, p2, p3...)
}

// UpdateJetTreeMinimockCounter returns a count of JetStorageMock.UpdateJetTreeFunc invocations
func (m *JetStorageMock) UpdateJetTreeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateJetTreeCounter)
}

// UpdateJetTreeMinimockPreCounter returns the value of JetStorageMock.UpdateJetTree invocations
func (m *JetStorageMock) UpdateJetTreeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateJetTreePreCounter)
}

// UpdateJetTreeFinished returns true if mock invocations count is ok
func (m *JetStorageMock) UpdateJetTreeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpdateJetTreeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpdateJetTreeCounter) == uint64(len(m.UpdateJetTreeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpdateJetTreeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpdateJetTreeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpdateJetTreeFunc != nil {
		return atomic.LoadUint64(&m.UpdateJetTreeCounter) > 0
	}

	return true
}

// ValidateCallCounters checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetStorageMock) ValidateCallCounters() {

	if !m.AddDropSizeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.AddDropSize")
	}

	if !m.AddJetsFinished() {
		m.t.Fatal("Expected call to JetStorageMock.AddJets")
	}

	if !m.CloneJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.CloneJetTree")
	}

	if !m.CreateDropFinished() {
		m.t.Fatal("Expected call to JetStorageMock.CreateDrop")
	}

	if !m.GetDropFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetDrop")
	}

	if !m.GetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetDropSizeHistory")
	}

	if !m.GetJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetJetTree")
	}

	if !m.GetJetsFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetJets")
	}

	if !m.SetDropFinished() {
		m.t.Fatal("Expected call to JetStorageMock.SetDrop")
	}

	if !m.SetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to JetStorageMock.SetDropSizeHistory")
	}

	if !m.SplitJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.SplitJetTree")
	}

	if !m.UpdateJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.UpdateJetTree")
	}

}

// CheckMocksCalled checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetStorageMock) CheckMocksCalled() {
	m.Finish()
}

// Finish checks that all mocked methods of the interface have been called at least once
// Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetStorageMock) Finish() {
	m.MinimockFinish()
}

// MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetStorageMock) MinimockFinish() {

	if !m.AddDropSizeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.AddDropSize")
	}

	if !m.AddJetsFinished() {
		m.t.Fatal("Expected call to JetStorageMock.AddJets")
	}

	if !m.CloneJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.CloneJetTree")
	}

	if !m.CreateDropFinished() {
		m.t.Fatal("Expected call to JetStorageMock.CreateDrop")
	}

	if !m.GetDropFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetDrop")
	}

	if !m.GetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetDropSizeHistory")
	}

	if !m.GetJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetJetTree")
	}

	if !m.GetJetsFinished() {
		m.t.Fatal("Expected call to JetStorageMock.GetJets")
	}

	if !m.SetDropFinished() {
		m.t.Fatal("Expected call to JetStorageMock.SetDrop")
	}

	if !m.SetDropSizeHistoryFinished() {
		m.t.Fatal("Expected call to JetStorageMock.SetDropSizeHistory")
	}

	if !m.SplitJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.SplitJetTree")
	}

	if !m.UpdateJetTreeFinished() {
		m.t.Fatal("Expected call to JetStorageMock.UpdateJetTree")
	}

}

// Wait waits for all mocked methods to be called at least once
// Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

// MinimockWait waits for all mocked methods to be called at least once
// this method is called by minimock.Controller
func (m *JetStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddDropSizeFinished()
		ok = ok && m.AddJetsFinished()
		ok = ok && m.CloneJetTreeFinished()
		ok = ok && m.CreateDropFinished()
		ok = ok && m.GetDropFinished()
		ok = ok && m.GetDropSizeHistoryFinished()
		ok = ok && m.GetJetTreeFinished()
		ok = ok && m.GetJetsFinished()
		ok = ok && m.SetDropFinished()
		ok = ok && m.SetDropSizeHistoryFinished()
		ok = ok && m.SplitJetTreeFinished()
		ok = ok && m.UpdateJetTreeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddDropSizeFinished() {
				m.t.Error("Expected call to JetStorageMock.AddDropSize")
			}

			if !m.AddJetsFinished() {
				m.t.Error("Expected call to JetStorageMock.AddJets")
			}

			if !m.CloneJetTreeFinished() {
				m.t.Error("Expected call to JetStorageMock.CloneJetTree")
			}

			if !m.CreateDropFinished() {
				m.t.Error("Expected call to JetStorageMock.CreateDrop")
			}

			if !m.GetDropFinished() {
				m.t.Error("Expected call to JetStorageMock.GetDrop")
			}

			if !m.GetDropSizeHistoryFinished() {
				m.t.Error("Expected call to JetStorageMock.GetDropSizeHistory")
			}

			if !m.GetJetTreeFinished() {
				m.t.Error("Expected call to JetStorageMock.GetJetTree")
			}

			if !m.GetJetsFinished() {
				m.t.Error("Expected call to JetStorageMock.GetJets")
			}

			if !m.SetDropFinished() {
				m.t.Error("Expected call to JetStorageMock.SetDrop")
			}

			if !m.SetDropSizeHistoryFinished() {
				m.t.Error("Expected call to JetStorageMock.SetDropSizeHistory")
			}

			if !m.SplitJetTreeFinished() {
				m.t.Error("Expected call to JetStorageMock.SplitJetTree")
			}

			if !m.UpdateJetTreeFinished() {
				m.t.Error("Expected call to JetStorageMock.UpdateJetTree")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
// it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *JetStorageMock) AllMocksCalled() bool {

	if !m.AddDropSizeFinished() {
		return false
	}

	if !m.AddJetsFinished() {
		return false
	}

	if !m.CloneJetTreeFinished() {
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

	if !m.GetJetTreeFinished() {
		return false
	}

	if !m.GetJetsFinished() {
		return false
	}

	if !m.SetDropFinished() {
		return false
	}

	if !m.SetDropSizeHistoryFinished() {
		return false
	}

	if !m.SplitJetTreeFinished() {
		return false
	}

	if !m.UpdateJetTreeFinished() {
		return false
	}

	return true
}
