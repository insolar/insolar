package sequence

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Sequencer" can be found in github.com/insolar/insolar/ledger/heavy/sequence
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//SequencerMock implements github.com/insolar/insolar/ledger/heavy/sequence.Sequencer
type SequencerMock struct {
	t minimock.Tester

	FirstFunc       func(p byte) (r *Item)
	FirstCounter    uint64
	FirstPreCounter uint64
	FirstMock       mSequencerMockFirst

	LastFunc       func(p byte) (r *Item)
	LastCounter    uint64
	LastPreCounter uint64
	LastMock       mSequencerMockLast

	LenFunc       func(p byte, p1 insolar.PulseNumber) (r uint32)
	LenCounter    uint64
	LenPreCounter uint64
	LenMock       mSequencerMockLen

	SliceFunc       func(p byte, p1 insolar.PulseNumber, p2 uint32, p3 uint32) (r []Item)
	SliceCounter    uint64
	SlicePreCounter uint64
	SliceMock       mSequencerMockSlice

	UpsertFunc       func(p byte, p1 []Item) (r error)
	UpsertCounter    uint64
	UpsertPreCounter uint64
	UpsertMock       mSequencerMockUpsert
}

//NewSequencerMock returns a mock for github.com/insolar/insolar/ledger/heavy/sequence.Sequencer
func NewSequencerMock(t minimock.Tester) *SequencerMock {
	m := &SequencerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FirstMock = mSequencerMockFirst{mock: m}
	m.LastMock = mSequencerMockLast{mock: m}
	m.LenMock = mSequencerMockLen{mock: m}
	m.SliceMock = mSequencerMockSlice{mock: m}
	m.UpsertMock = mSequencerMockUpsert{mock: m}

	return m
}

type mSequencerMockFirst struct {
	mock              *SequencerMock
	mainExpectation   *SequencerMockFirstExpectation
	expectationSeries []*SequencerMockFirstExpectation
}

type SequencerMockFirstExpectation struct {
	input  *SequencerMockFirstInput
	result *SequencerMockFirstResult
}

type SequencerMockFirstInput struct {
	p byte
}

type SequencerMockFirstResult struct {
	r *Item
}

//Expect specifies that invocation of Sequencer.First is expected from 1 to Infinity times
func (m *mSequencerMockFirst) Expect(p byte) *mSequencerMockFirst {
	m.mock.FirstFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockFirstExpectation{}
	}
	m.mainExpectation.input = &SequencerMockFirstInput{p}
	return m
}

//Return specifies results of invocation of Sequencer.First
func (m *mSequencerMockFirst) Return(r *Item) *SequencerMock {
	m.mock.FirstFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockFirstExpectation{}
	}
	m.mainExpectation.result = &SequencerMockFirstResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Sequencer.First is expected once
func (m *mSequencerMockFirst) ExpectOnce(p byte) *SequencerMockFirstExpectation {
	m.mock.FirstFunc = nil
	m.mainExpectation = nil

	expectation := &SequencerMockFirstExpectation{}
	expectation.input = &SequencerMockFirstInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SequencerMockFirstExpectation) Return(r *Item) {
	e.result = &SequencerMockFirstResult{r}
}

//Set uses given function f as a mock of Sequencer.First method
func (m *mSequencerMockFirst) Set(f func(p byte) (r *Item)) *SequencerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FirstFunc = f
	return m.mock
}

//First implements github.com/insolar/insolar/ledger/heavy/sequence.Sequencer interface
func (m *SequencerMock) First(p byte) (r *Item) {
	counter := atomic.AddUint64(&m.FirstPreCounter, 1)
	defer atomic.AddUint64(&m.FirstCounter, 1)

	if len(m.FirstMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FirstMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SequencerMock.First. %v", p)
			return
		}

		input := m.FirstMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SequencerMockFirstInput{p}, "Sequencer.First got unexpected parameters")

		result := m.FirstMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.First")
			return
		}

		r = result.r

		return
	}

	if m.FirstMock.mainExpectation != nil {

		input := m.FirstMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SequencerMockFirstInput{p}, "Sequencer.First got unexpected parameters")
		}

		result := m.FirstMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.First")
		}

		r = result.r

		return
	}

	if m.FirstFunc == nil {
		m.t.Fatalf("Unexpected call to SequencerMock.First. %v", p)
		return
	}

	return m.FirstFunc(p)
}

//FirstMinimockCounter returns a count of SequencerMock.FirstFunc invocations
func (m *SequencerMock) FirstMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FirstCounter)
}

//FirstMinimockPreCounter returns the value of SequencerMock.First invocations
func (m *SequencerMock) FirstMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FirstPreCounter)
}

//FirstFinished returns true if mock invocations count is ok
func (m *SequencerMock) FirstFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FirstMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FirstCounter) == uint64(len(m.FirstMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FirstMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FirstCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FirstFunc != nil {
		return atomic.LoadUint64(&m.FirstCounter) > 0
	}

	return true
}

type mSequencerMockLast struct {
	mock              *SequencerMock
	mainExpectation   *SequencerMockLastExpectation
	expectationSeries []*SequencerMockLastExpectation
}

type SequencerMockLastExpectation struct {
	input  *SequencerMockLastInput
	result *SequencerMockLastResult
}

type SequencerMockLastInput struct {
	p byte
}

type SequencerMockLastResult struct {
	r *Item
}

//Expect specifies that invocation of Sequencer.Last is expected from 1 to Infinity times
func (m *mSequencerMockLast) Expect(p byte) *mSequencerMockLast {
	m.mock.LastFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockLastExpectation{}
	}
	m.mainExpectation.input = &SequencerMockLastInput{p}
	return m
}

//Return specifies results of invocation of Sequencer.Last
func (m *mSequencerMockLast) Return(r *Item) *SequencerMock {
	m.mock.LastFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockLastExpectation{}
	}
	m.mainExpectation.result = &SequencerMockLastResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Sequencer.Last is expected once
func (m *mSequencerMockLast) ExpectOnce(p byte) *SequencerMockLastExpectation {
	m.mock.LastFunc = nil
	m.mainExpectation = nil

	expectation := &SequencerMockLastExpectation{}
	expectation.input = &SequencerMockLastInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SequencerMockLastExpectation) Return(r *Item) {
	e.result = &SequencerMockLastResult{r}
}

//Set uses given function f as a mock of Sequencer.Last method
func (m *mSequencerMockLast) Set(f func(p byte) (r *Item)) *SequencerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LastFunc = f
	return m.mock
}

//Last implements github.com/insolar/insolar/ledger/heavy/sequence.Sequencer interface
func (m *SequencerMock) Last(p byte) (r *Item) {
	counter := atomic.AddUint64(&m.LastPreCounter, 1)
	defer atomic.AddUint64(&m.LastCounter, 1)

	if len(m.LastMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LastMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SequencerMock.Last. %v", p)
			return
		}

		input := m.LastMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SequencerMockLastInput{p}, "Sequencer.Last got unexpected parameters")

		result := m.LastMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Last")
			return
		}

		r = result.r

		return
	}

	if m.LastMock.mainExpectation != nil {

		input := m.LastMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SequencerMockLastInput{p}, "Sequencer.Last got unexpected parameters")
		}

		result := m.LastMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Last")
		}

		r = result.r

		return
	}

	if m.LastFunc == nil {
		m.t.Fatalf("Unexpected call to SequencerMock.Last. %v", p)
		return
	}

	return m.LastFunc(p)
}

//LastMinimockCounter returns a count of SequencerMock.LastFunc invocations
func (m *SequencerMock) LastMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LastCounter)
}

//LastMinimockPreCounter returns the value of SequencerMock.Last invocations
func (m *SequencerMock) LastMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LastPreCounter)
}

//LastFinished returns true if mock invocations count is ok
func (m *SequencerMock) LastFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LastMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LastCounter) == uint64(len(m.LastMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LastMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LastCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LastFunc != nil {
		return atomic.LoadUint64(&m.LastCounter) > 0
	}

	return true
}

type mSequencerMockLen struct {
	mock              *SequencerMock
	mainExpectation   *SequencerMockLenExpectation
	expectationSeries []*SequencerMockLenExpectation
}

type SequencerMockLenExpectation struct {
	input  *SequencerMockLenInput
	result *SequencerMockLenResult
}

type SequencerMockLenInput struct {
	p  byte
	p1 insolar.PulseNumber
}

type SequencerMockLenResult struct {
	r uint32
}

//Expect specifies that invocation of Sequencer.Len is expected from 1 to Infinity times
func (m *mSequencerMockLen) Expect(p byte, p1 insolar.PulseNumber) *mSequencerMockLen {
	m.mock.LenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockLenExpectation{}
	}
	m.mainExpectation.input = &SequencerMockLenInput{p, p1}
	return m
}

//Return specifies results of invocation of Sequencer.Len
func (m *mSequencerMockLen) Return(r uint32) *SequencerMock {
	m.mock.LenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockLenExpectation{}
	}
	m.mainExpectation.result = &SequencerMockLenResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Sequencer.Len is expected once
func (m *mSequencerMockLen) ExpectOnce(p byte, p1 insolar.PulseNumber) *SequencerMockLenExpectation {
	m.mock.LenFunc = nil
	m.mainExpectation = nil

	expectation := &SequencerMockLenExpectation{}
	expectation.input = &SequencerMockLenInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SequencerMockLenExpectation) Return(r uint32) {
	e.result = &SequencerMockLenResult{r}
}

//Set uses given function f as a mock of Sequencer.Len method
func (m *mSequencerMockLen) Set(f func(p byte, p1 insolar.PulseNumber) (r uint32)) *SequencerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LenFunc = f
	return m.mock
}

//Len implements github.com/insolar/insolar/ledger/heavy/sequence.Sequencer interface
func (m *SequencerMock) Len(p byte, p1 insolar.PulseNumber) (r uint32) {
	counter := atomic.AddUint64(&m.LenPreCounter, 1)
	defer atomic.AddUint64(&m.LenCounter, 1)

	if len(m.LenMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LenMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SequencerMock.Len. %v %v", p, p1)
			return
		}

		input := m.LenMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SequencerMockLenInput{p, p1}, "Sequencer.Len got unexpected parameters")

		result := m.LenMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Len")
			return
		}

		r = result.r

		return
	}

	if m.LenMock.mainExpectation != nil {

		input := m.LenMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SequencerMockLenInput{p, p1}, "Sequencer.Len got unexpected parameters")
		}

		result := m.LenMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Len")
		}

		r = result.r

		return
	}

	if m.LenFunc == nil {
		m.t.Fatalf("Unexpected call to SequencerMock.Len. %v %v", p, p1)
		return
	}

	return m.LenFunc(p, p1)
}

//LenMinimockCounter returns a count of SequencerMock.LenFunc invocations
func (m *SequencerMock) LenMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LenCounter)
}

//LenMinimockPreCounter returns the value of SequencerMock.Len invocations
func (m *SequencerMock) LenMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LenPreCounter)
}

//LenFinished returns true if mock invocations count is ok
func (m *SequencerMock) LenFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LenMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LenCounter) == uint64(len(m.LenMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LenMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LenCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LenFunc != nil {
		return atomic.LoadUint64(&m.LenCounter) > 0
	}

	return true
}

type mSequencerMockSlice struct {
	mock              *SequencerMock
	mainExpectation   *SequencerMockSliceExpectation
	expectationSeries []*SequencerMockSliceExpectation
}

type SequencerMockSliceExpectation struct {
	input  *SequencerMockSliceInput
	result *SequencerMockSliceResult
}

type SequencerMockSliceInput struct {
	p  byte
	p1 insolar.PulseNumber
	p2 uint32
	p3 uint32
}

type SequencerMockSliceResult struct {
	r []Item
}

//Expect specifies that invocation of Sequencer.Slice is expected from 1 to Infinity times
func (m *mSequencerMockSlice) Expect(p byte, p1 insolar.PulseNumber, p2 uint32, p3 uint32) *mSequencerMockSlice {
	m.mock.SliceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockSliceExpectation{}
	}
	m.mainExpectation.input = &SequencerMockSliceInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Sequencer.Slice
func (m *mSequencerMockSlice) Return(r []Item) *SequencerMock {
	m.mock.SliceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockSliceExpectation{}
	}
	m.mainExpectation.result = &SequencerMockSliceResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Sequencer.Slice is expected once
func (m *mSequencerMockSlice) ExpectOnce(p byte, p1 insolar.PulseNumber, p2 uint32, p3 uint32) *SequencerMockSliceExpectation {
	m.mock.SliceFunc = nil
	m.mainExpectation = nil

	expectation := &SequencerMockSliceExpectation{}
	expectation.input = &SequencerMockSliceInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SequencerMockSliceExpectation) Return(r []Item) {
	e.result = &SequencerMockSliceResult{r}
}

//Set uses given function f as a mock of Sequencer.Slice method
func (m *mSequencerMockSlice) Set(f func(p byte, p1 insolar.PulseNumber, p2 uint32, p3 uint32) (r []Item)) *SequencerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SliceFunc = f
	return m.mock
}

//Slice implements github.com/insolar/insolar/ledger/heavy/sequence.Sequencer interface
func (m *SequencerMock) Slice(p byte, p1 insolar.PulseNumber, p2 uint32, p3 uint32) (r []Item) {
	counter := atomic.AddUint64(&m.SlicePreCounter, 1)
	defer atomic.AddUint64(&m.SliceCounter, 1)

	if len(m.SliceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SliceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SequencerMock.Slice. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SliceMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SequencerMockSliceInput{p, p1, p2, p3}, "Sequencer.Slice got unexpected parameters")

		result := m.SliceMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Slice")
			return
		}

		r = result.r

		return
	}

	if m.SliceMock.mainExpectation != nil {

		input := m.SliceMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SequencerMockSliceInput{p, p1, p2, p3}, "Sequencer.Slice got unexpected parameters")
		}

		result := m.SliceMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Slice")
		}

		r = result.r

		return
	}

	if m.SliceFunc == nil {
		m.t.Fatalf("Unexpected call to SequencerMock.Slice. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SliceFunc(p, p1, p2, p3)
}

//SliceMinimockCounter returns a count of SequencerMock.SliceFunc invocations
func (m *SequencerMock) SliceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SliceCounter)
}

//SliceMinimockPreCounter returns the value of SequencerMock.Slice invocations
func (m *SequencerMock) SliceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SlicePreCounter)
}

//SliceFinished returns true if mock invocations count is ok
func (m *SequencerMock) SliceFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SliceMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SliceCounter) == uint64(len(m.SliceMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SliceMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SliceCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SliceFunc != nil {
		return atomic.LoadUint64(&m.SliceCounter) > 0
	}

	return true
}

type mSequencerMockUpsert struct {
	mock              *SequencerMock
	mainExpectation   *SequencerMockUpsertExpectation
	expectationSeries []*SequencerMockUpsertExpectation
}

type SequencerMockUpsertExpectation struct {
	input  *SequencerMockUpsertInput
	result *SequencerMockUpsertResult
}

type SequencerMockUpsertInput struct {
	p  byte
	p1 []Item
}

type SequencerMockUpsertResult struct {
	r error
}

//Expect specifies that invocation of Sequencer.Upsert is expected from 1 to Infinity times
func (m *mSequencerMockUpsert) Expect(p byte, p1 []Item) *mSequencerMockUpsert {
	m.mock.UpsertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockUpsertExpectation{}
	}
	m.mainExpectation.input = &SequencerMockUpsertInput{p, p1}
	return m
}

//Return specifies results of invocation of Sequencer.Upsert
func (m *mSequencerMockUpsert) Return(r error) *SequencerMock {
	m.mock.UpsertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SequencerMockUpsertExpectation{}
	}
	m.mainExpectation.result = &SequencerMockUpsertResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Sequencer.Upsert is expected once
func (m *mSequencerMockUpsert) ExpectOnce(p byte, p1 []Item) *SequencerMockUpsertExpectation {
	m.mock.UpsertFunc = nil
	m.mainExpectation = nil

	expectation := &SequencerMockUpsertExpectation{}
	expectation.input = &SequencerMockUpsertInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SequencerMockUpsertExpectation) Return(r error) {
	e.result = &SequencerMockUpsertResult{r}
}

//Set uses given function f as a mock of Sequencer.Upsert method
func (m *mSequencerMockUpsert) Set(f func(p byte, p1 []Item) (r error)) *SequencerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpsertFunc = f
	return m.mock
}

//Upsert implements github.com/insolar/insolar/ledger/heavy/sequence.Sequencer interface
func (m *SequencerMock) Upsert(p byte, p1 []Item) (r error) {
	counter := atomic.AddUint64(&m.UpsertPreCounter, 1)
	defer atomic.AddUint64(&m.UpsertCounter, 1)

	if len(m.UpsertMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpsertMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SequencerMock.Upsert. %v %v", p, p1)
			return
		}

		input := m.UpsertMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SequencerMockUpsertInput{p, p1}, "Sequencer.Upsert got unexpected parameters")

		result := m.UpsertMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Upsert")
			return
		}

		r = result.r

		return
	}

	if m.UpsertMock.mainExpectation != nil {

		input := m.UpsertMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SequencerMockUpsertInput{p, p1}, "Sequencer.Upsert got unexpected parameters")
		}

		result := m.UpsertMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SequencerMock.Upsert")
		}

		r = result.r

		return
	}

	if m.UpsertFunc == nil {
		m.t.Fatalf("Unexpected call to SequencerMock.Upsert. %v %v", p, p1)
		return
	}

	return m.UpsertFunc(p, p1)
}

//UpsertMinimockCounter returns a count of SequencerMock.UpsertFunc invocations
func (m *SequencerMock) UpsertMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertCounter)
}

//UpsertMinimockPreCounter returns the value of SequencerMock.Upsert invocations
func (m *SequencerMock) UpsertMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpsertPreCounter)
}

//UpsertFinished returns true if mock invocations count is ok
func (m *SequencerMock) UpsertFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpsertMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpsertCounter) == uint64(len(m.UpsertMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpsertMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpsertCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpsertFunc != nil {
		return atomic.LoadUint64(&m.UpsertCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SequencerMock) ValidateCallCounters() {

	if !m.FirstFinished() {
		m.t.Fatal("Expected call to SequencerMock.First")
	}

	if !m.LastFinished() {
		m.t.Fatal("Expected call to SequencerMock.Last")
	}

	if !m.LenFinished() {
		m.t.Fatal("Expected call to SequencerMock.Len")
	}

	if !m.SliceFinished() {
		m.t.Fatal("Expected call to SequencerMock.Slice")
	}

	if !m.UpsertFinished() {
		m.t.Fatal("Expected call to SequencerMock.Upsert")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SequencerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SequencerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SequencerMock) MinimockFinish() {

	if !m.FirstFinished() {
		m.t.Fatal("Expected call to SequencerMock.First")
	}

	if !m.LastFinished() {
		m.t.Fatal("Expected call to SequencerMock.Last")
	}

	if !m.LenFinished() {
		m.t.Fatal("Expected call to SequencerMock.Len")
	}

	if !m.SliceFinished() {
		m.t.Fatal("Expected call to SequencerMock.Slice")
	}

	if !m.UpsertFinished() {
		m.t.Fatal("Expected call to SequencerMock.Upsert")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SequencerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SequencerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FirstFinished()
		ok = ok && m.LastFinished()
		ok = ok && m.LenFinished()
		ok = ok && m.SliceFinished()
		ok = ok && m.UpsertFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FirstFinished() {
				m.t.Error("Expected call to SequencerMock.First")
			}

			if !m.LastFinished() {
				m.t.Error("Expected call to SequencerMock.Last")
			}

			if !m.LenFinished() {
				m.t.Error("Expected call to SequencerMock.Len")
			}

			if !m.SliceFinished() {
				m.t.Error("Expected call to SequencerMock.Slice")
			}

			if !m.UpsertFinished() {
				m.t.Error("Expected call to SequencerMock.Upsert")
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
func (m *SequencerMock) AllMocksCalled() bool {

	if !m.FirstFinished() {
		return false
	}

	if !m.LastFinished() {
		return false
	}

	if !m.LenFinished() {
		return false
	}

	if !m.SliceFinished() {
		return false
	}

	if !m.UpsertFinished() {
		return false
	}

	return true
}
