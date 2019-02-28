package drop

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Builder" can be found in github.com/insolar/insolar/ledger/storage/jet/drop
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	jet "github.com/insolar/insolar/ledger/storage/jet"

	testify_assert "github.com/stretchr/testify/assert"
)

//BuilderMock implements github.com/insolar/insolar/ledger/storage/jet/drop.Builder
type BuilderMock struct {
	t minimock.Tester

	AppendFunc       func(p Hashable) (r error)
	AppendCounter    uint64
	AppendPreCounter uint64
	AppendMock       mBuilderMockAppend

	BuildFunc       func() (r jet.Drop, r1 error)
	BuildCounter    uint64
	BuildPreCounter uint64
	BuildMock       mBuilderMockBuild

	PrevHashFunc       func(p []byte)
	PrevHashCounter    uint64
	PrevHashPreCounter uint64
	PrevHashMock       mBuilderMockPrevHash

	PulseFunc       func(p core.PulseNumber)
	PulseCounter    uint64
	PulsePreCounter uint64
	PulseMock       mBuilderMockPulse

	SizeFunc       func(p uint64)
	SizeCounter    uint64
	SizePreCounter uint64
	SizeMock       mBuilderMockSize
}

//NewBuilderMock returns a mock for github.com/insolar/insolar/ledger/storage/jet/drop.Builder
func NewBuilderMock(t minimock.Tester) *BuilderMock {
	m := &BuilderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AppendMock = mBuilderMockAppend{mock: m}
	m.BuildMock = mBuilderMockBuild{mock: m}
	m.PrevHashMock = mBuilderMockPrevHash{mock: m}
	m.PulseMock = mBuilderMockPulse{mock: m}
	m.SizeMock = mBuilderMockSize{mock: m}

	return m
}

type mBuilderMockAppend struct {
	mock              *BuilderMock
	mainExpectation   *BuilderMockAppendExpectation
	expectationSeries []*BuilderMockAppendExpectation
}

type BuilderMockAppendExpectation struct {
	input  *BuilderMockAppendInput
	result *BuilderMockAppendResult
}

type BuilderMockAppendInput struct {
	p Hashable
}

type BuilderMockAppendResult struct {
	r error
}

//Expect specifies that invocation of Builder.Append is expected from 1 to Infinity times
func (m *mBuilderMockAppend) Expect(p Hashable) *mBuilderMockAppend {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockAppendExpectation{}
	}
	m.mainExpectation.input = &BuilderMockAppendInput{p}
	return m
}

//Return specifies results of invocation of Builder.Append
func (m *mBuilderMockAppend) Return(r error) *BuilderMock {
	m.mock.AppendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockAppendExpectation{}
	}
	m.mainExpectation.result = &BuilderMockAppendResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Builder.Append is expected once
func (m *mBuilderMockAppend) ExpectOnce(p Hashable) *BuilderMockAppendExpectation {
	m.mock.AppendFunc = nil
	m.mainExpectation = nil

	expectation := &BuilderMockAppendExpectation{}
	expectation.input = &BuilderMockAppendInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BuilderMockAppendExpectation) Return(r error) {
	e.result = &BuilderMockAppendResult{r}
}

//Set uses given function f as a mock of Builder.Append method
func (m *mBuilderMockAppend) Set(f func(p Hashable) (r error)) *BuilderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AppendFunc = f
	return m.mock
}

//Append implements github.com/insolar/insolar/ledger/storage/jet/drop.Builder interface
func (m *BuilderMock) Append(p Hashable) (r error) {
	counter := atomic.AddUint64(&m.AppendPreCounter, 1)
	defer atomic.AddUint64(&m.AppendCounter, 1)

	if len(m.AppendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AppendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BuilderMock.Append. %v", p)
			return
		}

		input := m.AppendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, BuilderMockAppendInput{p}, "Builder.Append got unexpected parameters")

		result := m.AppendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BuilderMock.Append")
			return
		}

		r = result.r

		return
	}

	if m.AppendMock.mainExpectation != nil {

		input := m.AppendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, BuilderMockAppendInput{p}, "Builder.Append got unexpected parameters")
		}

		result := m.AppendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BuilderMock.Append")
		}

		r = result.r

		return
	}

	if m.AppendFunc == nil {
		m.t.Fatalf("Unexpected call to BuilderMock.Append. %v", p)
		return
	}

	return m.AppendFunc(p)
}

//AppendMinimockCounter returns a count of BuilderMock.AppendFunc invocations
func (m *BuilderMock) AppendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AppendCounter)
}

//AppendMinimockPreCounter returns the value of BuilderMock.Append invocations
func (m *BuilderMock) AppendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AppendPreCounter)
}

//AppendFinished returns true if mock invocations count is ok
func (m *BuilderMock) AppendFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AppendMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AppendCounter) == uint64(len(m.AppendMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AppendMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AppendCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AppendFunc != nil {
		return atomic.LoadUint64(&m.AppendCounter) > 0
	}

	return true
}

type mBuilderMockBuild struct {
	mock              *BuilderMock
	mainExpectation   *BuilderMockBuildExpectation
	expectationSeries []*BuilderMockBuildExpectation
}

type BuilderMockBuildExpectation struct {
	result *BuilderMockBuildResult
}

type BuilderMockBuildResult struct {
	r  jet.Drop
	r1 error
}

//Expect specifies that invocation of Builder.Build is expected from 1 to Infinity times
func (m *mBuilderMockBuild) Expect() *mBuilderMockBuild {
	m.mock.BuildFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockBuildExpectation{}
	}

	return m
}

//Return specifies results of invocation of Builder.Build
func (m *mBuilderMockBuild) Return(r jet.Drop, r1 error) *BuilderMock {
	m.mock.BuildFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockBuildExpectation{}
	}
	m.mainExpectation.result = &BuilderMockBuildResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Builder.Build is expected once
func (m *mBuilderMockBuild) ExpectOnce() *BuilderMockBuildExpectation {
	m.mock.BuildFunc = nil
	m.mainExpectation = nil

	expectation := &BuilderMockBuildExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BuilderMockBuildExpectation) Return(r jet.Drop, r1 error) {
	e.result = &BuilderMockBuildResult{r, r1}
}

//Set uses given function f as a mock of Builder.Build method
func (m *mBuilderMockBuild) Set(f func() (r jet.Drop, r1 error)) *BuilderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BuildFunc = f
	return m.mock
}

//Build implements github.com/insolar/insolar/ledger/storage/jet/drop.Builder interface
func (m *BuilderMock) Build() (r jet.Drop, r1 error) {
	counter := atomic.AddUint64(&m.BuildPreCounter, 1)
	defer atomic.AddUint64(&m.BuildCounter, 1)

	if len(m.BuildMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BuildMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BuilderMock.Build.")
			return
		}

		result := m.BuildMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BuilderMock.Build")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BuildMock.mainExpectation != nil {

		result := m.BuildMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BuilderMock.Build")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BuildFunc == nil {
		m.t.Fatalf("Unexpected call to BuilderMock.Build.")
		return
	}

	return m.BuildFunc()
}

//BuildMinimockCounter returns a count of BuilderMock.BuildFunc invocations
func (m *BuilderMock) BuildMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BuildCounter)
}

//BuildMinimockPreCounter returns the value of BuilderMock.Build invocations
func (m *BuilderMock) BuildMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BuildPreCounter)
}

//BuildFinished returns true if mock invocations count is ok
func (m *BuilderMock) BuildFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.BuildMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.BuildCounter) == uint64(len(m.BuildMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.BuildMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.BuildCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.BuildFunc != nil {
		return atomic.LoadUint64(&m.BuildCounter) > 0
	}

	return true
}

type mBuilderMockPrevHash struct {
	mock              *BuilderMock
	mainExpectation   *BuilderMockPrevHashExpectation
	expectationSeries []*BuilderMockPrevHashExpectation
}

type BuilderMockPrevHashExpectation struct {
	input *BuilderMockPrevHashInput
}

type BuilderMockPrevHashInput struct {
	p []byte
}

//Expect specifies that invocation of Builder.PrevHash is expected from 1 to Infinity times
func (m *mBuilderMockPrevHash) Expect(p []byte) *mBuilderMockPrevHash {
	m.mock.PrevHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockPrevHashExpectation{}
	}
	m.mainExpectation.input = &BuilderMockPrevHashInput{p}
	return m
}

//Return specifies results of invocation of Builder.PrevHash
func (m *mBuilderMockPrevHash) Return() *BuilderMock {
	m.mock.PrevHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockPrevHashExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Builder.PrevHash is expected once
func (m *mBuilderMockPrevHash) ExpectOnce(p []byte) *BuilderMockPrevHashExpectation {
	m.mock.PrevHashFunc = nil
	m.mainExpectation = nil

	expectation := &BuilderMockPrevHashExpectation{}
	expectation.input = &BuilderMockPrevHashInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Builder.PrevHash method
func (m *mBuilderMockPrevHash) Set(f func(p []byte)) *BuilderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PrevHashFunc = f
	return m.mock
}

//PrevHash implements github.com/insolar/insolar/ledger/storage/jet/drop.Builder interface
func (m *BuilderMock) PrevHash(p []byte) {
	counter := atomic.AddUint64(&m.PrevHashPreCounter, 1)
	defer atomic.AddUint64(&m.PrevHashCounter, 1)

	if len(m.PrevHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PrevHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BuilderMock.PrevHash. %v", p)
			return
		}

		input := m.PrevHashMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, BuilderMockPrevHashInput{p}, "Builder.PrevHash got unexpected parameters")

		return
	}

	if m.PrevHashMock.mainExpectation != nil {

		input := m.PrevHashMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, BuilderMockPrevHashInput{p}, "Builder.PrevHash got unexpected parameters")
		}

		return
	}

	if m.PrevHashFunc == nil {
		m.t.Fatalf("Unexpected call to BuilderMock.PrevHash. %v", p)
		return
	}

	m.PrevHashFunc(p)
}

//PrevHashMinimockCounter returns a count of BuilderMock.PrevHashFunc invocations
func (m *BuilderMock) PrevHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PrevHashCounter)
}

//PrevHashMinimockPreCounter returns the value of BuilderMock.PrevHash invocations
func (m *BuilderMock) PrevHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PrevHashPreCounter)
}

//PrevHashFinished returns true if mock invocations count is ok
func (m *BuilderMock) PrevHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PrevHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PrevHashCounter) == uint64(len(m.PrevHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PrevHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PrevHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PrevHashFunc != nil {
		return atomic.LoadUint64(&m.PrevHashCounter) > 0
	}

	return true
}

type mBuilderMockPulse struct {
	mock              *BuilderMock
	mainExpectation   *BuilderMockPulseExpectation
	expectationSeries []*BuilderMockPulseExpectation
}

type BuilderMockPulseExpectation struct {
	input *BuilderMockPulseInput
}

type BuilderMockPulseInput struct {
	p core.PulseNumber
}

//Expect specifies that invocation of Builder.Pulse is expected from 1 to Infinity times
func (m *mBuilderMockPulse) Expect(p core.PulseNumber) *mBuilderMockPulse {
	m.mock.PulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockPulseExpectation{}
	}
	m.mainExpectation.input = &BuilderMockPulseInput{p}
	return m
}

//Return specifies results of invocation of Builder.Pulse
func (m *mBuilderMockPulse) Return() *BuilderMock {
	m.mock.PulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockPulseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Builder.Pulse is expected once
func (m *mBuilderMockPulse) ExpectOnce(p core.PulseNumber) *BuilderMockPulseExpectation {
	m.mock.PulseFunc = nil
	m.mainExpectation = nil

	expectation := &BuilderMockPulseExpectation{}
	expectation.input = &BuilderMockPulseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Builder.Pulse method
func (m *mBuilderMockPulse) Set(f func(p core.PulseNumber)) *BuilderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PulseFunc = f
	return m.mock
}

//Pulse implements github.com/insolar/insolar/ledger/storage/jet/drop.Builder interface
func (m *BuilderMock) Pulse(p core.PulseNumber) {
	counter := atomic.AddUint64(&m.PulsePreCounter, 1)
	defer atomic.AddUint64(&m.PulseCounter, 1)

	if len(m.PulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BuilderMock.Pulse. %v", p)
			return
		}

		input := m.PulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, BuilderMockPulseInput{p}, "Builder.Pulse got unexpected parameters")

		return
	}

	if m.PulseMock.mainExpectation != nil {

		input := m.PulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, BuilderMockPulseInput{p}, "Builder.Pulse got unexpected parameters")
		}

		return
	}

	if m.PulseFunc == nil {
		m.t.Fatalf("Unexpected call to BuilderMock.Pulse. %v", p)
		return
	}

	m.PulseFunc(p)
}

//PulseMinimockCounter returns a count of BuilderMock.PulseFunc invocations
func (m *BuilderMock) PulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PulseCounter)
}

//PulseMinimockPreCounter returns the value of BuilderMock.Pulse invocations
func (m *BuilderMock) PulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PulsePreCounter)
}

//PulseFinished returns true if mock invocations count is ok
func (m *BuilderMock) PulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PulseCounter) == uint64(len(m.PulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PulseFunc != nil {
		return atomic.LoadUint64(&m.PulseCounter) > 0
	}

	return true
}

type mBuilderMockSize struct {
	mock              *BuilderMock
	mainExpectation   *BuilderMockSizeExpectation
	expectationSeries []*BuilderMockSizeExpectation
}

type BuilderMockSizeExpectation struct {
	input *BuilderMockSizeInput
}

type BuilderMockSizeInput struct {
	p uint64
}

//Expect specifies that invocation of Builder.Size is expected from 1 to Infinity times
func (m *mBuilderMockSize) Expect(p uint64) *mBuilderMockSize {
	m.mock.SizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockSizeExpectation{}
	}
	m.mainExpectation.input = &BuilderMockSizeInput{p}
	return m
}

//Return specifies results of invocation of Builder.Size
func (m *mBuilderMockSize) Return() *BuilderMock {
	m.mock.SizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BuilderMockSizeExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Builder.Size is expected once
func (m *mBuilderMockSize) ExpectOnce(p uint64) *BuilderMockSizeExpectation {
	m.mock.SizeFunc = nil
	m.mainExpectation = nil

	expectation := &BuilderMockSizeExpectation{}
	expectation.input = &BuilderMockSizeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Builder.Size method
func (m *mBuilderMockSize) Set(f func(p uint64)) *BuilderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SizeFunc = f
	return m.mock
}

//Size implements github.com/insolar/insolar/ledger/storage/jet/drop.Builder interface
func (m *BuilderMock) Size(p uint64) {
	counter := atomic.AddUint64(&m.SizePreCounter, 1)
	defer atomic.AddUint64(&m.SizeCounter, 1)

	if len(m.SizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BuilderMock.Size. %v", p)
			return
		}

		input := m.SizeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, BuilderMockSizeInput{p}, "Builder.Size got unexpected parameters")

		return
	}

	if m.SizeMock.mainExpectation != nil {

		input := m.SizeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, BuilderMockSizeInput{p}, "Builder.Size got unexpected parameters")
		}

		return
	}

	if m.SizeFunc == nil {
		m.t.Fatalf("Unexpected call to BuilderMock.Size. %v", p)
		return
	}

	m.SizeFunc(p)
}

//SizeMinimockCounter returns a count of BuilderMock.SizeFunc invocations
func (m *BuilderMock) SizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SizeCounter)
}

//SizeMinimockPreCounter returns the value of BuilderMock.Size invocations
func (m *BuilderMock) SizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SizePreCounter)
}

//SizeFinished returns true if mock invocations count is ok
func (m *BuilderMock) SizeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SizeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SizeCounter) == uint64(len(m.SizeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SizeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SizeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SizeFunc != nil {
		return atomic.LoadUint64(&m.SizeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BuilderMock) ValidateCallCounters() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to BuilderMock.Append")
	}

	if !m.BuildFinished() {
		m.t.Fatal("Expected call to BuilderMock.Build")
	}

	if !m.PrevHashFinished() {
		m.t.Fatal("Expected call to BuilderMock.PrevHash")
	}

	if !m.PulseFinished() {
		m.t.Fatal("Expected call to BuilderMock.Pulse")
	}

	if !m.SizeFinished() {
		m.t.Fatal("Expected call to BuilderMock.Size")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BuilderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *BuilderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *BuilderMock) MinimockFinish() {

	if !m.AppendFinished() {
		m.t.Fatal("Expected call to BuilderMock.Append")
	}

	if !m.BuildFinished() {
		m.t.Fatal("Expected call to BuilderMock.Build")
	}

	if !m.PrevHashFinished() {
		m.t.Fatal("Expected call to BuilderMock.PrevHash")
	}

	if !m.PulseFinished() {
		m.t.Fatal("Expected call to BuilderMock.Pulse")
	}

	if !m.SizeFinished() {
		m.t.Fatal("Expected call to BuilderMock.Size")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *BuilderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *BuilderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AppendFinished()
		ok = ok && m.BuildFinished()
		ok = ok && m.PrevHashFinished()
		ok = ok && m.PulseFinished()
		ok = ok && m.SizeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AppendFinished() {
				m.t.Error("Expected call to BuilderMock.Append")
			}

			if !m.BuildFinished() {
				m.t.Error("Expected call to BuilderMock.Build")
			}

			if !m.PrevHashFinished() {
				m.t.Error("Expected call to BuilderMock.PrevHash")
			}

			if !m.PulseFinished() {
				m.t.Error("Expected call to BuilderMock.Pulse")
			}

			if !m.SizeFinished() {
				m.t.Error("Expected call to BuilderMock.Size")
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
func (m *BuilderMock) AllMocksCalled() bool {

	if !m.AppendFinished() {
		return false
	}

	if !m.BuildFinished() {
		return false
	}

	if !m.PrevHashFinished() {
		return false
	}

	if !m.PulseFinished() {
		return false
	}

	if !m.SizeFinished() {
		return false
	}

	return true
}
