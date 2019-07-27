package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ExecutionArchive" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ExecutionArchiveMock implements github.com/insolar/insolar/logicrunner.ExecutionArchive
type ExecutionArchiveMock struct {
	t minimock.Tester

	ArchiveFunc       func(p *Transcript)
	ArchiveCounter    uint64
	ArchivePreCounter uint64
	ArchiveMock       mExecutionArchiveMockArchive

	DoneFunc       func(p *Transcript) (r bool)
	DoneCounter    uint64
	DonePreCounter uint64
	DoneMock       mExecutionArchiveMockDone

	FindRequestLoopFunc       func(p context.Context, p1 string) (r bool)
	FindRequestLoopCounter    uint64
	FindRequestLoopPreCounter uint64
	FindRequestLoopMock       mExecutionArchiveMockFindRequestLoop

	GetActiveTranscriptFunc       func(p insolar.Reference) (r *Transcript)
	GetActiveTranscriptCounter    uint64
	GetActiveTranscriptPreCounter uint64
	GetActiveTranscriptMock       mExecutionArchiveMockGetActiveTranscript

	IsEmptyFunc       func() (r bool)
	IsEmptyCounter    uint64
	IsEmptyPreCounter uint64
	IsEmptyMock       mExecutionArchiveMockIsEmpty

	OnPulseFunc       func(p context.Context) (r []insolar.Message)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mExecutionArchiveMockOnPulse
}

//NewExecutionArchiveMock returns a mock for github.com/insolar/insolar/logicrunner.ExecutionArchive
func NewExecutionArchiveMock(t minimock.Tester) *ExecutionArchiveMock {
	m := &ExecutionArchiveMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ArchiveMock = mExecutionArchiveMockArchive{mock: m}
	m.DoneMock = mExecutionArchiveMockDone{mock: m}
	m.FindRequestLoopMock = mExecutionArchiveMockFindRequestLoop{mock: m}
	m.GetActiveTranscriptMock = mExecutionArchiveMockGetActiveTranscript{mock: m}
	m.IsEmptyMock = mExecutionArchiveMockIsEmpty{mock: m}
	m.OnPulseMock = mExecutionArchiveMockOnPulse{mock: m}

	return m
}

type mExecutionArchiveMockArchive struct {
	mock              *ExecutionArchiveMock
	mainExpectation   *ExecutionArchiveMockArchiveExpectation
	expectationSeries []*ExecutionArchiveMockArchiveExpectation
}

type ExecutionArchiveMockArchiveExpectation struct {
	input *ExecutionArchiveMockArchiveInput
}

type ExecutionArchiveMockArchiveInput struct {
	p *Transcript
}

//Expect specifies that invocation of ExecutionArchive.Archive is expected from 1 to Infinity times
func (m *mExecutionArchiveMockArchive) Expect(p *Transcript) *mExecutionArchiveMockArchive {
	m.mock.ArchiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockArchiveExpectation{}
	}
	m.mainExpectation.input = &ExecutionArchiveMockArchiveInput{p}
	return m
}

//Return specifies results of invocation of ExecutionArchive.Archive
func (m *mExecutionArchiveMockArchive) Return() *ExecutionArchiveMock {
	m.mock.ArchiveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockArchiveExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionArchive.Archive is expected once
func (m *mExecutionArchiveMockArchive) ExpectOnce(p *Transcript) *ExecutionArchiveMockArchiveExpectation {
	m.mock.ArchiveFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionArchiveMockArchiveExpectation{}
	expectation.input = &ExecutionArchiveMockArchiveInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ExecutionArchive.Archive method
func (m *mExecutionArchiveMockArchive) Set(f func(p *Transcript)) *ExecutionArchiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ArchiveFunc = f
	return m.mock
}

//Archive implements github.com/insolar/insolar/logicrunner.ExecutionArchive interface
func (m *ExecutionArchiveMock) Archive(p *Transcript) {
	counter := atomic.AddUint64(&m.ArchivePreCounter, 1)
	defer atomic.AddUint64(&m.ArchiveCounter, 1)

	if len(m.ArchiveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ArchiveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionArchiveMock.Archive. %v", p)
			return
		}

		input := m.ArchiveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionArchiveMockArchiveInput{p}, "ExecutionArchive.Archive got unexpected parameters")

		return
	}

	if m.ArchiveMock.mainExpectation != nil {

		input := m.ArchiveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionArchiveMockArchiveInput{p}, "ExecutionArchive.Archive got unexpected parameters")
		}

		return
	}

	if m.ArchiveFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionArchiveMock.Archive. %v", p)
		return
	}

	m.ArchiveFunc(p)
}

//ArchiveMinimockCounter returns a count of ExecutionArchiveMock.ArchiveFunc invocations
func (m *ExecutionArchiveMock) ArchiveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ArchiveCounter)
}

//ArchiveMinimockPreCounter returns the value of ExecutionArchiveMock.Archive invocations
func (m *ExecutionArchiveMock) ArchiveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ArchivePreCounter)
}

//ArchiveFinished returns true if mock invocations count is ok
func (m *ExecutionArchiveMock) ArchiveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ArchiveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ArchiveCounter) == uint64(len(m.ArchiveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ArchiveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ArchiveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ArchiveFunc != nil {
		return atomic.LoadUint64(&m.ArchiveCounter) > 0
	}

	return true
}

type mExecutionArchiveMockDone struct {
	mock              *ExecutionArchiveMock
	mainExpectation   *ExecutionArchiveMockDoneExpectation
	expectationSeries []*ExecutionArchiveMockDoneExpectation
}

type ExecutionArchiveMockDoneExpectation struct {
	input  *ExecutionArchiveMockDoneInput
	result *ExecutionArchiveMockDoneResult
}

type ExecutionArchiveMockDoneInput struct {
	p *Transcript
}

type ExecutionArchiveMockDoneResult struct {
	r bool
}

//Expect specifies that invocation of ExecutionArchive.Done is expected from 1 to Infinity times
func (m *mExecutionArchiveMockDone) Expect(p *Transcript) *mExecutionArchiveMockDone {
	m.mock.DoneFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockDoneExpectation{}
	}
	m.mainExpectation.input = &ExecutionArchiveMockDoneInput{p}
	return m
}

//Return specifies results of invocation of ExecutionArchive.Done
func (m *mExecutionArchiveMockDone) Return(r bool) *ExecutionArchiveMock {
	m.mock.DoneFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockDoneExpectation{}
	}
	m.mainExpectation.result = &ExecutionArchiveMockDoneResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionArchive.Done is expected once
func (m *mExecutionArchiveMockDone) ExpectOnce(p *Transcript) *ExecutionArchiveMockDoneExpectation {
	m.mock.DoneFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionArchiveMockDoneExpectation{}
	expectation.input = &ExecutionArchiveMockDoneInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionArchiveMockDoneExpectation) Return(r bool) {
	e.result = &ExecutionArchiveMockDoneResult{r}
}

//Set uses given function f as a mock of ExecutionArchive.Done method
func (m *mExecutionArchiveMockDone) Set(f func(p *Transcript) (r bool)) *ExecutionArchiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DoneFunc = f
	return m.mock
}

//Done implements github.com/insolar/insolar/logicrunner.ExecutionArchive interface
func (m *ExecutionArchiveMock) Done(p *Transcript) (r bool) {
	counter := atomic.AddUint64(&m.DonePreCounter, 1)
	defer atomic.AddUint64(&m.DoneCounter, 1)

	if len(m.DoneMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DoneMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionArchiveMock.Done. %v", p)
			return
		}

		input := m.DoneMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionArchiveMockDoneInput{p}, "ExecutionArchive.Done got unexpected parameters")

		result := m.DoneMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.Done")
			return
		}

		r = result.r

		return
	}

	if m.DoneMock.mainExpectation != nil {

		input := m.DoneMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionArchiveMockDoneInput{p}, "ExecutionArchive.Done got unexpected parameters")
		}

		result := m.DoneMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.Done")
		}

		r = result.r

		return
	}

	if m.DoneFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionArchiveMock.Done. %v", p)
		return
	}

	return m.DoneFunc(p)
}

//DoneMinimockCounter returns a count of ExecutionArchiveMock.DoneFunc invocations
func (m *ExecutionArchiveMock) DoneMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DoneCounter)
}

//DoneMinimockPreCounter returns the value of ExecutionArchiveMock.Done invocations
func (m *ExecutionArchiveMock) DoneMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DonePreCounter)
}

//DoneFinished returns true if mock invocations count is ok
func (m *ExecutionArchiveMock) DoneFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DoneMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DoneCounter) == uint64(len(m.DoneMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DoneMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DoneCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DoneFunc != nil {
		return atomic.LoadUint64(&m.DoneCounter) > 0
	}

	return true
}

type mExecutionArchiveMockFindRequestLoop struct {
	mock              *ExecutionArchiveMock
	mainExpectation   *ExecutionArchiveMockFindRequestLoopExpectation
	expectationSeries []*ExecutionArchiveMockFindRequestLoopExpectation
}

type ExecutionArchiveMockFindRequestLoopExpectation struct {
	input  *ExecutionArchiveMockFindRequestLoopInput
	result *ExecutionArchiveMockFindRequestLoopResult
}

type ExecutionArchiveMockFindRequestLoopInput struct {
	p  context.Context
	p1 string
}

type ExecutionArchiveMockFindRequestLoopResult struct {
	r bool
}

//Expect specifies that invocation of ExecutionArchive.FindRequestLoop is expected from 1 to Infinity times
func (m *mExecutionArchiveMockFindRequestLoop) Expect(p context.Context, p1 string) *mExecutionArchiveMockFindRequestLoop {
	m.mock.FindRequestLoopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockFindRequestLoopExpectation{}
	}
	m.mainExpectation.input = &ExecutionArchiveMockFindRequestLoopInput{p, p1}
	return m
}

//Return specifies results of invocation of ExecutionArchive.FindRequestLoop
func (m *mExecutionArchiveMockFindRequestLoop) Return(r bool) *ExecutionArchiveMock {
	m.mock.FindRequestLoopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockFindRequestLoopExpectation{}
	}
	m.mainExpectation.result = &ExecutionArchiveMockFindRequestLoopResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionArchive.FindRequestLoop is expected once
func (m *mExecutionArchiveMockFindRequestLoop) ExpectOnce(p context.Context, p1 string) *ExecutionArchiveMockFindRequestLoopExpectation {
	m.mock.FindRequestLoopFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionArchiveMockFindRequestLoopExpectation{}
	expectation.input = &ExecutionArchiveMockFindRequestLoopInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionArchiveMockFindRequestLoopExpectation) Return(r bool) {
	e.result = &ExecutionArchiveMockFindRequestLoopResult{r}
}

//Set uses given function f as a mock of ExecutionArchive.FindRequestLoop method
func (m *mExecutionArchiveMockFindRequestLoop) Set(f func(p context.Context, p1 string) (r bool)) *ExecutionArchiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FindRequestLoopFunc = f
	return m.mock
}

//FindRequestLoop implements github.com/insolar/insolar/logicrunner.ExecutionArchive interface
func (m *ExecutionArchiveMock) FindRequestLoop(p context.Context, p1 string) (r bool) {
	counter := atomic.AddUint64(&m.FindRequestLoopPreCounter, 1)
	defer atomic.AddUint64(&m.FindRequestLoopCounter, 1)

	if len(m.FindRequestLoopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FindRequestLoopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionArchiveMock.FindRequestLoop. %v %v", p, p1)
			return
		}

		input := m.FindRequestLoopMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionArchiveMockFindRequestLoopInput{p, p1}, "ExecutionArchive.FindRequestLoop got unexpected parameters")

		result := m.FindRequestLoopMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.FindRequestLoop")
			return
		}

		r = result.r

		return
	}

	if m.FindRequestLoopMock.mainExpectation != nil {

		input := m.FindRequestLoopMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionArchiveMockFindRequestLoopInput{p, p1}, "ExecutionArchive.FindRequestLoop got unexpected parameters")
		}

		result := m.FindRequestLoopMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.FindRequestLoop")
		}

		r = result.r

		return
	}

	if m.FindRequestLoopFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionArchiveMock.FindRequestLoop. %v %v", p, p1)
		return
	}

	return m.FindRequestLoopFunc(p, p1)
}

//FindRequestLoopMinimockCounter returns a count of ExecutionArchiveMock.FindRequestLoopFunc invocations
func (m *ExecutionArchiveMock) FindRequestLoopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FindRequestLoopCounter)
}

//FindRequestLoopMinimockPreCounter returns the value of ExecutionArchiveMock.FindRequestLoop invocations
func (m *ExecutionArchiveMock) FindRequestLoopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FindRequestLoopPreCounter)
}

//FindRequestLoopFinished returns true if mock invocations count is ok
func (m *ExecutionArchiveMock) FindRequestLoopFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FindRequestLoopMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FindRequestLoopCounter) == uint64(len(m.FindRequestLoopMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FindRequestLoopMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FindRequestLoopCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FindRequestLoopFunc != nil {
		return atomic.LoadUint64(&m.FindRequestLoopCounter) > 0
	}

	return true
}

type mExecutionArchiveMockGetActiveTranscript struct {
	mock              *ExecutionArchiveMock
	mainExpectation   *ExecutionArchiveMockGetActiveTranscriptExpectation
	expectationSeries []*ExecutionArchiveMockGetActiveTranscriptExpectation
}

type ExecutionArchiveMockGetActiveTranscriptExpectation struct {
	input  *ExecutionArchiveMockGetActiveTranscriptInput
	result *ExecutionArchiveMockGetActiveTranscriptResult
}

type ExecutionArchiveMockGetActiveTranscriptInput struct {
	p insolar.Reference
}

type ExecutionArchiveMockGetActiveTranscriptResult struct {
	r *Transcript
}

//Expect specifies that invocation of ExecutionArchive.GetActiveTranscript is expected from 1 to Infinity times
func (m *mExecutionArchiveMockGetActiveTranscript) Expect(p insolar.Reference) *mExecutionArchiveMockGetActiveTranscript {
	m.mock.GetActiveTranscriptFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockGetActiveTranscriptExpectation{}
	}
	m.mainExpectation.input = &ExecutionArchiveMockGetActiveTranscriptInput{p}
	return m
}

//Return specifies results of invocation of ExecutionArchive.GetActiveTranscript
func (m *mExecutionArchiveMockGetActiveTranscript) Return(r *Transcript) *ExecutionArchiveMock {
	m.mock.GetActiveTranscriptFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockGetActiveTranscriptExpectation{}
	}
	m.mainExpectation.result = &ExecutionArchiveMockGetActiveTranscriptResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionArchive.GetActiveTranscript is expected once
func (m *mExecutionArchiveMockGetActiveTranscript) ExpectOnce(p insolar.Reference) *ExecutionArchiveMockGetActiveTranscriptExpectation {
	m.mock.GetActiveTranscriptFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionArchiveMockGetActiveTranscriptExpectation{}
	expectation.input = &ExecutionArchiveMockGetActiveTranscriptInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionArchiveMockGetActiveTranscriptExpectation) Return(r *Transcript) {
	e.result = &ExecutionArchiveMockGetActiveTranscriptResult{r}
}

//Set uses given function f as a mock of ExecutionArchive.GetActiveTranscript method
func (m *mExecutionArchiveMockGetActiveTranscript) Set(f func(p insolar.Reference) (r *Transcript)) *ExecutionArchiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveTranscriptFunc = f
	return m.mock
}

//GetActiveTranscript implements github.com/insolar/insolar/logicrunner.ExecutionArchive interface
func (m *ExecutionArchiveMock) GetActiveTranscript(p insolar.Reference) (r *Transcript) {
	counter := atomic.AddUint64(&m.GetActiveTranscriptPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveTranscriptCounter, 1)

	if len(m.GetActiveTranscriptMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveTranscriptMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionArchiveMock.GetActiveTranscript. %v", p)
			return
		}

		input := m.GetActiveTranscriptMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionArchiveMockGetActiveTranscriptInput{p}, "ExecutionArchive.GetActiveTranscript got unexpected parameters")

		result := m.GetActiveTranscriptMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.GetActiveTranscript")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveTranscriptMock.mainExpectation != nil {

		input := m.GetActiveTranscriptMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionArchiveMockGetActiveTranscriptInput{p}, "ExecutionArchive.GetActiveTranscript got unexpected parameters")
		}

		result := m.GetActiveTranscriptMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.GetActiveTranscript")
		}

		r = result.r

		return
	}

	if m.GetActiveTranscriptFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionArchiveMock.GetActiveTranscript. %v", p)
		return
	}

	return m.GetActiveTranscriptFunc(p)
}

//GetActiveTranscriptMinimockCounter returns a count of ExecutionArchiveMock.GetActiveTranscriptFunc invocations
func (m *ExecutionArchiveMock) GetActiveTranscriptMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveTranscriptCounter)
}

//GetActiveTranscriptMinimockPreCounter returns the value of ExecutionArchiveMock.GetActiveTranscript invocations
func (m *ExecutionArchiveMock) GetActiveTranscriptMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveTranscriptPreCounter)
}

//GetActiveTranscriptFinished returns true if mock invocations count is ok
func (m *ExecutionArchiveMock) GetActiveTranscriptFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveTranscriptMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveTranscriptCounter) == uint64(len(m.GetActiveTranscriptMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveTranscriptMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveTranscriptCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveTranscriptFunc != nil {
		return atomic.LoadUint64(&m.GetActiveTranscriptCounter) > 0
	}

	return true
}

type mExecutionArchiveMockIsEmpty struct {
	mock              *ExecutionArchiveMock
	mainExpectation   *ExecutionArchiveMockIsEmptyExpectation
	expectationSeries []*ExecutionArchiveMockIsEmptyExpectation
}

type ExecutionArchiveMockIsEmptyExpectation struct {
	result *ExecutionArchiveMockIsEmptyResult
}

type ExecutionArchiveMockIsEmptyResult struct {
	r bool
}

//Expect specifies that invocation of ExecutionArchive.IsEmpty is expected from 1 to Infinity times
func (m *mExecutionArchiveMockIsEmpty) Expect() *mExecutionArchiveMockIsEmpty {
	m.mock.IsEmptyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockIsEmptyExpectation{}
	}

	return m
}

//Return specifies results of invocation of ExecutionArchive.IsEmpty
func (m *mExecutionArchiveMockIsEmpty) Return(r bool) *ExecutionArchiveMock {
	m.mock.IsEmptyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockIsEmptyExpectation{}
	}
	m.mainExpectation.result = &ExecutionArchiveMockIsEmptyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionArchive.IsEmpty is expected once
func (m *mExecutionArchiveMockIsEmpty) ExpectOnce() *ExecutionArchiveMockIsEmptyExpectation {
	m.mock.IsEmptyFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionArchiveMockIsEmptyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionArchiveMockIsEmptyExpectation) Return(r bool) {
	e.result = &ExecutionArchiveMockIsEmptyResult{r}
}

//Set uses given function f as a mock of ExecutionArchive.IsEmpty method
func (m *mExecutionArchiveMockIsEmpty) Set(f func() (r bool)) *ExecutionArchiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsEmptyFunc = f
	return m.mock
}

//IsEmpty implements github.com/insolar/insolar/logicrunner.ExecutionArchive interface
func (m *ExecutionArchiveMock) IsEmpty() (r bool) {
	counter := atomic.AddUint64(&m.IsEmptyPreCounter, 1)
	defer atomic.AddUint64(&m.IsEmptyCounter, 1)

	if len(m.IsEmptyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsEmptyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionArchiveMock.IsEmpty.")
			return
		}

		result := m.IsEmptyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.IsEmpty")
			return
		}

		r = result.r

		return
	}

	if m.IsEmptyMock.mainExpectation != nil {

		result := m.IsEmptyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.IsEmpty")
		}

		r = result.r

		return
	}

	if m.IsEmptyFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionArchiveMock.IsEmpty.")
		return
	}

	return m.IsEmptyFunc()
}

//IsEmptyMinimockCounter returns a count of ExecutionArchiveMock.IsEmptyFunc invocations
func (m *ExecutionArchiveMock) IsEmptyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsEmptyCounter)
}

//IsEmptyMinimockPreCounter returns the value of ExecutionArchiveMock.IsEmpty invocations
func (m *ExecutionArchiveMock) IsEmptyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsEmptyPreCounter)
}

//IsEmptyFinished returns true if mock invocations count is ok
func (m *ExecutionArchiveMock) IsEmptyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsEmptyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsEmptyCounter) == uint64(len(m.IsEmptyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsEmptyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsEmptyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsEmptyFunc != nil {
		return atomic.LoadUint64(&m.IsEmptyCounter) > 0
	}

	return true
}

type mExecutionArchiveMockOnPulse struct {
	mock              *ExecutionArchiveMock
	mainExpectation   *ExecutionArchiveMockOnPulseExpectation
	expectationSeries []*ExecutionArchiveMockOnPulseExpectation
}

type ExecutionArchiveMockOnPulseExpectation struct {
	input  *ExecutionArchiveMockOnPulseInput
	result *ExecutionArchiveMockOnPulseResult
}

type ExecutionArchiveMockOnPulseInput struct {
	p context.Context
}

type ExecutionArchiveMockOnPulseResult struct {
	r []insolar.Message
}

//Expect specifies that invocation of ExecutionArchive.OnPulse is expected from 1 to Infinity times
func (m *mExecutionArchiveMockOnPulse) Expect(p context.Context) *mExecutionArchiveMockOnPulse {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockOnPulseExpectation{}
	}
	m.mainExpectation.input = &ExecutionArchiveMockOnPulseInput{p}
	return m
}

//Return specifies results of invocation of ExecutionArchive.OnPulse
func (m *mExecutionArchiveMockOnPulse) Return(r []insolar.Message) *ExecutionArchiveMock {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ExecutionArchiveMockOnPulseExpectation{}
	}
	m.mainExpectation.result = &ExecutionArchiveMockOnPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ExecutionArchive.OnPulse is expected once
func (m *mExecutionArchiveMockOnPulse) ExpectOnce(p context.Context) *ExecutionArchiveMockOnPulseExpectation {
	m.mock.OnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ExecutionArchiveMockOnPulseExpectation{}
	expectation.input = &ExecutionArchiveMockOnPulseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ExecutionArchiveMockOnPulseExpectation) Return(r []insolar.Message) {
	e.result = &ExecutionArchiveMockOnPulseResult{r}
}

//Set uses given function f as a mock of ExecutionArchive.OnPulse method
func (m *mExecutionArchiveMockOnPulse) Set(f func(p context.Context) (r []insolar.Message)) *ExecutionArchiveMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/logicrunner.ExecutionArchive interface
func (m *ExecutionArchiveMock) OnPulse(p context.Context) (r []insolar.Message) {
	counter := atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if len(m.OnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ExecutionArchiveMock.OnPulse. %v", p)
			return
		}

		input := m.OnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ExecutionArchiveMockOnPulseInput{p}, "ExecutionArchive.OnPulse got unexpected parameters")

		result := m.OnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.OnPulse")
			return
		}

		r = result.r

		return
	}

	if m.OnPulseMock.mainExpectation != nil {

		input := m.OnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ExecutionArchiveMockOnPulseInput{p}, "ExecutionArchive.OnPulse got unexpected parameters")
		}

		result := m.OnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ExecutionArchiveMock.OnPulse")
		}

		r = result.r

		return
	}

	if m.OnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ExecutionArchiveMock.OnPulse. %v", p)
		return
	}

	return m.OnPulseFunc(p)
}

//OnPulseMinimockCounter returns a count of ExecutionArchiveMock.OnPulseFunc invocations
func (m *ExecutionArchiveMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of ExecutionArchiveMock.OnPulse invocations
func (m *ExecutionArchiveMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//OnPulseFinished returns true if mock invocations count is ok
func (m *ExecutionArchiveMock) OnPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OnPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OnPulseCounter) == uint64(len(m.OnPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OnPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OnPulseFunc != nil {
		return atomic.LoadUint64(&m.OnPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExecutionArchiveMock) ValidateCallCounters() {

	if !m.ArchiveFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.Archive")
	}

	if !m.DoneFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.Done")
	}

	if !m.FindRequestLoopFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.FindRequestLoop")
	}

	if !m.GetActiveTranscriptFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.GetActiveTranscript")
	}

	if !m.IsEmptyFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.IsEmpty")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.OnPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ExecutionArchiveMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ExecutionArchiveMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ExecutionArchiveMock) MinimockFinish() {

	if !m.ArchiveFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.Archive")
	}

	if !m.DoneFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.Done")
	}

	if !m.FindRequestLoopFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.FindRequestLoop")
	}

	if !m.GetActiveTranscriptFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.GetActiveTranscript")
	}

	if !m.IsEmptyFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.IsEmpty")
	}

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to ExecutionArchiveMock.OnPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ExecutionArchiveMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ExecutionArchiveMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ArchiveFinished()
		ok = ok && m.DoneFinished()
		ok = ok && m.FindRequestLoopFinished()
		ok = ok && m.GetActiveTranscriptFinished()
		ok = ok && m.IsEmptyFinished()
		ok = ok && m.OnPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ArchiveFinished() {
				m.t.Error("Expected call to ExecutionArchiveMock.Archive")
			}

			if !m.DoneFinished() {
				m.t.Error("Expected call to ExecutionArchiveMock.Done")
			}

			if !m.FindRequestLoopFinished() {
				m.t.Error("Expected call to ExecutionArchiveMock.FindRequestLoop")
			}

			if !m.GetActiveTranscriptFinished() {
				m.t.Error("Expected call to ExecutionArchiveMock.GetActiveTranscript")
			}

			if !m.IsEmptyFinished() {
				m.t.Error("Expected call to ExecutionArchiveMock.IsEmpty")
			}

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to ExecutionArchiveMock.OnPulse")
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
func (m *ExecutionArchiveMock) AllMocksCalled() bool {

	if !m.ArchiveFinished() {
		return false
	}

	if !m.DoneFinished() {
		return false
	}

	if !m.FindRequestLoopFinished() {
		return false
	}

	if !m.GetActiveTranscriptFinished() {
		return false
	}

	if !m.IsEmptyFinished() {
		return false
	}

	if !m.OnPulseFinished() {
		return false
	}

	return true
}
