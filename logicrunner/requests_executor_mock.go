package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RequestsExecutor" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	artifacts "github.com/insolar/insolar/logicrunner/artifacts"

	testify_assert "github.com/stretchr/testify/assert"
)

//RequestsExecutorMock implements github.com/insolar/insolar/logicrunner.RequestsExecutor
type RequestsExecutorMock struct {
	t minimock.Tester

	ExecuteFunc       func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)
	ExecuteCounter    uint64
	ExecutePreCounter uint64
	ExecuteMock       mRequestsExecutorMockExecute

	ExecuteAndSaveFunc       func(p context.Context, p1 *Transcript) (r insolar.Reply, r1 error)
	ExecuteAndSaveCounter    uint64
	ExecuteAndSavePreCounter uint64
	ExecuteAndSaveMock       mRequestsExecutorMockExecuteAndSave

	SaveFunc       func(p context.Context, p1 *Transcript, p2 artifacts.RequestResult) (r insolar.Reply, r1 error)
	SaveCounter    uint64
	SavePreCounter uint64
	SaveMock       mRequestsExecutorMockSave

	SendReplyFunc       func(p context.Context, p1 *Transcript, p2 insolar.Reply, p3 error)
	SendReplyCounter    uint64
	SendReplyPreCounter uint64
	SendReplyMock       mRequestsExecutorMockSendReply
}

//NewRequestsExecutorMock returns a mock for github.com/insolar/insolar/logicrunner.RequestsExecutor
func NewRequestsExecutorMock(t minimock.Tester) *RequestsExecutorMock {
	m := &RequestsExecutorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExecuteMock = mRequestsExecutorMockExecute{mock: m}
	m.ExecuteAndSaveMock = mRequestsExecutorMockExecuteAndSave{mock: m}
	m.SaveMock = mRequestsExecutorMockSave{mock: m}
	m.SendReplyMock = mRequestsExecutorMockSendReply{mock: m}

	return m
}

type mRequestsExecutorMockExecute struct {
	mock              *RequestsExecutorMock
	mainExpectation   *RequestsExecutorMockExecuteExpectation
	expectationSeries []*RequestsExecutorMockExecuteExpectation
}

type RequestsExecutorMockExecuteExpectation struct {
	input  *RequestsExecutorMockExecuteInput
	result *RequestsExecutorMockExecuteResult
}

type RequestsExecutorMockExecuteInput struct {
	p  context.Context
	p1 *Transcript
}

type RequestsExecutorMockExecuteResult struct {
	r  artifacts.RequestResult
	r1 error
}

//Expect specifies that invocation of RequestsExecutor.Execute is expected from 1 to Infinity times
func (m *mRequestsExecutorMockExecute) Expect(p context.Context, p1 *Transcript) *mRequestsExecutorMockExecute {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockExecuteExpectation{}
	}
	m.mainExpectation.input = &RequestsExecutorMockExecuteInput{p, p1}
	return m
}

//Return specifies results of invocation of RequestsExecutor.Execute
func (m *mRequestsExecutorMockExecute) Return(r artifacts.RequestResult, r1 error) *RequestsExecutorMock {
	m.mock.ExecuteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockExecuteExpectation{}
	}
	m.mainExpectation.result = &RequestsExecutorMockExecuteResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RequestsExecutor.Execute is expected once
func (m *mRequestsExecutorMockExecute) ExpectOnce(p context.Context, p1 *Transcript) *RequestsExecutorMockExecuteExpectation {
	m.mock.ExecuteFunc = nil
	m.mainExpectation = nil

	expectation := &RequestsExecutorMockExecuteExpectation{}
	expectation.input = &RequestsExecutorMockExecuteInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestsExecutorMockExecuteExpectation) Return(r artifacts.RequestResult, r1 error) {
	e.result = &RequestsExecutorMockExecuteResult{r, r1}
}

//Set uses given function f as a mock of RequestsExecutor.Execute method
func (m *mRequestsExecutorMockExecute) Set(f func(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error)) *RequestsExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecuteFunc = f
	return m.mock
}

//Execute implements github.com/insolar/insolar/logicrunner.RequestsExecutor interface
func (m *RequestsExecutorMock) Execute(p context.Context, p1 *Transcript) (r artifacts.RequestResult, r1 error) {
	counter := atomic.AddUint64(&m.ExecutePreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteCounter, 1)

	if len(m.ExecuteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecuteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestsExecutorMock.Execute. %v %v", p, p1)
			return
		}

		input := m.ExecuteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RequestsExecutorMockExecuteInput{p, p1}, "RequestsExecutor.Execute got unexpected parameters")

		result := m.ExecuteMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestsExecutorMock.Execute")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteMock.mainExpectation != nil {

		input := m.ExecuteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RequestsExecutorMockExecuteInput{p, p1}, "RequestsExecutor.Execute got unexpected parameters")
		}

		result := m.ExecuteMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestsExecutorMock.Execute")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteFunc == nil {
		m.t.Fatalf("Unexpected call to RequestsExecutorMock.Execute. %v %v", p, p1)
		return
	}

	return m.ExecuteFunc(p, p1)
}

//ExecuteMinimockCounter returns a count of RequestsExecutorMock.ExecuteFunc invocations
func (m *RequestsExecutorMock) ExecuteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteCounter)
}

//ExecuteMinimockPreCounter returns the value of RequestsExecutorMock.Execute invocations
func (m *RequestsExecutorMock) ExecuteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecutePreCounter)
}

//ExecuteFinished returns true if mock invocations count is ok
func (m *RequestsExecutorMock) ExecuteFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExecuteMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExecuteCounter) == uint64(len(m.ExecuteMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExecuteMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExecuteCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExecuteFunc != nil {
		return atomic.LoadUint64(&m.ExecuteCounter) > 0
	}

	return true
}

type mRequestsExecutorMockExecuteAndSave struct {
	mock              *RequestsExecutorMock
	mainExpectation   *RequestsExecutorMockExecuteAndSaveExpectation
	expectationSeries []*RequestsExecutorMockExecuteAndSaveExpectation
}

type RequestsExecutorMockExecuteAndSaveExpectation struct {
	input  *RequestsExecutorMockExecuteAndSaveInput
	result *RequestsExecutorMockExecuteAndSaveResult
}

type RequestsExecutorMockExecuteAndSaveInput struct {
	p  context.Context
	p1 *Transcript
}

type RequestsExecutorMockExecuteAndSaveResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of RequestsExecutor.ExecuteAndSave is expected from 1 to Infinity times
func (m *mRequestsExecutorMockExecuteAndSave) Expect(p context.Context, p1 *Transcript) *mRequestsExecutorMockExecuteAndSave {
	m.mock.ExecuteAndSaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockExecuteAndSaveExpectation{}
	}
	m.mainExpectation.input = &RequestsExecutorMockExecuteAndSaveInput{p, p1}
	return m
}

//Return specifies results of invocation of RequestsExecutor.ExecuteAndSave
func (m *mRequestsExecutorMockExecuteAndSave) Return(r insolar.Reply, r1 error) *RequestsExecutorMock {
	m.mock.ExecuteAndSaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockExecuteAndSaveExpectation{}
	}
	m.mainExpectation.result = &RequestsExecutorMockExecuteAndSaveResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RequestsExecutor.ExecuteAndSave is expected once
func (m *mRequestsExecutorMockExecuteAndSave) ExpectOnce(p context.Context, p1 *Transcript) *RequestsExecutorMockExecuteAndSaveExpectation {
	m.mock.ExecuteAndSaveFunc = nil
	m.mainExpectation = nil

	expectation := &RequestsExecutorMockExecuteAndSaveExpectation{}
	expectation.input = &RequestsExecutorMockExecuteAndSaveInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestsExecutorMockExecuteAndSaveExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &RequestsExecutorMockExecuteAndSaveResult{r, r1}
}

//Set uses given function f as a mock of RequestsExecutor.ExecuteAndSave method
func (m *mRequestsExecutorMockExecuteAndSave) Set(f func(p context.Context, p1 *Transcript) (r insolar.Reply, r1 error)) *RequestsExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExecuteAndSaveFunc = f
	return m.mock
}

//ExecuteAndSave implements github.com/insolar/insolar/logicrunner.RequestsExecutor interface
func (m *RequestsExecutorMock) ExecuteAndSave(p context.Context, p1 *Transcript) (r insolar.Reply, r1 error) {
	counter := atomic.AddUint64(&m.ExecuteAndSavePreCounter, 1)
	defer atomic.AddUint64(&m.ExecuteAndSaveCounter, 1)

	if len(m.ExecuteAndSaveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExecuteAndSaveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestsExecutorMock.ExecuteAndSave. %v %v", p, p1)
			return
		}

		input := m.ExecuteAndSaveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RequestsExecutorMockExecuteAndSaveInput{p, p1}, "RequestsExecutor.ExecuteAndSave got unexpected parameters")

		result := m.ExecuteAndSaveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestsExecutorMock.ExecuteAndSave")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteAndSaveMock.mainExpectation != nil {

		input := m.ExecuteAndSaveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RequestsExecutorMockExecuteAndSaveInput{p, p1}, "RequestsExecutor.ExecuteAndSave got unexpected parameters")
		}

		result := m.ExecuteAndSaveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestsExecutorMock.ExecuteAndSave")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExecuteAndSaveFunc == nil {
		m.t.Fatalf("Unexpected call to RequestsExecutorMock.ExecuteAndSave. %v %v", p, p1)
		return
	}

	return m.ExecuteAndSaveFunc(p, p1)
}

//ExecuteAndSaveMinimockCounter returns a count of RequestsExecutorMock.ExecuteAndSaveFunc invocations
func (m *RequestsExecutorMock) ExecuteAndSaveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteAndSaveCounter)
}

//ExecuteAndSaveMinimockPreCounter returns the value of RequestsExecutorMock.ExecuteAndSave invocations
func (m *RequestsExecutorMock) ExecuteAndSaveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExecuteAndSavePreCounter)
}

//ExecuteAndSaveFinished returns true if mock invocations count is ok
func (m *RequestsExecutorMock) ExecuteAndSaveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExecuteAndSaveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExecuteAndSaveCounter) == uint64(len(m.ExecuteAndSaveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExecuteAndSaveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExecuteAndSaveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExecuteAndSaveFunc != nil {
		return atomic.LoadUint64(&m.ExecuteAndSaveCounter) > 0
	}

	return true
}

type mRequestsExecutorMockSave struct {
	mock              *RequestsExecutorMock
	mainExpectation   *RequestsExecutorMockSaveExpectation
	expectationSeries []*RequestsExecutorMockSaveExpectation
}

type RequestsExecutorMockSaveExpectation struct {
	input  *RequestsExecutorMockSaveInput
	result *RequestsExecutorMockSaveResult
}

type RequestsExecutorMockSaveInput struct {
	p  context.Context
	p1 *Transcript
	p2 artifacts.RequestResult
}

type RequestsExecutorMockSaveResult struct {
	r  insolar.Reply
	r1 error
}

//Expect specifies that invocation of RequestsExecutor.Save is expected from 1 to Infinity times
func (m *mRequestsExecutorMockSave) Expect(p context.Context, p1 *Transcript, p2 artifacts.RequestResult) *mRequestsExecutorMockSave {
	m.mock.SaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockSaveExpectation{}
	}
	m.mainExpectation.input = &RequestsExecutorMockSaveInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of RequestsExecutor.Save
func (m *mRequestsExecutorMockSave) Return(r insolar.Reply, r1 error) *RequestsExecutorMock {
	m.mock.SaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockSaveExpectation{}
	}
	m.mainExpectation.result = &RequestsExecutorMockSaveResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RequestsExecutor.Save is expected once
func (m *mRequestsExecutorMockSave) ExpectOnce(p context.Context, p1 *Transcript, p2 artifacts.RequestResult) *RequestsExecutorMockSaveExpectation {
	m.mock.SaveFunc = nil
	m.mainExpectation = nil

	expectation := &RequestsExecutorMockSaveExpectation{}
	expectation.input = &RequestsExecutorMockSaveInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RequestsExecutorMockSaveExpectation) Return(r insolar.Reply, r1 error) {
	e.result = &RequestsExecutorMockSaveResult{r, r1}
}

//Set uses given function f as a mock of RequestsExecutor.Save method
func (m *mRequestsExecutorMockSave) Set(f func(p context.Context, p1 *Transcript, p2 artifacts.RequestResult) (r insolar.Reply, r1 error)) *RequestsExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SaveFunc = f
	return m.mock
}

//Save implements github.com/insolar/insolar/logicrunner.RequestsExecutor interface
func (m *RequestsExecutorMock) Save(p context.Context, p1 *Transcript, p2 artifacts.RequestResult) (r insolar.Reply, r1 error) {
	counter := atomic.AddUint64(&m.SavePreCounter, 1)
	defer atomic.AddUint64(&m.SaveCounter, 1)

	if len(m.SaveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SaveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestsExecutorMock.Save. %v %v %v", p, p1, p2)
			return
		}

		input := m.SaveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RequestsExecutorMockSaveInput{p, p1, p2}, "RequestsExecutor.Save got unexpected parameters")

		result := m.SaveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RequestsExecutorMock.Save")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SaveMock.mainExpectation != nil {

		input := m.SaveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RequestsExecutorMockSaveInput{p, p1, p2}, "RequestsExecutor.Save got unexpected parameters")
		}

		result := m.SaveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RequestsExecutorMock.Save")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SaveFunc == nil {
		m.t.Fatalf("Unexpected call to RequestsExecutorMock.Save. %v %v %v", p, p1, p2)
		return
	}

	return m.SaveFunc(p, p1, p2)
}

//SaveMinimockCounter returns a count of RequestsExecutorMock.SaveFunc invocations
func (m *RequestsExecutorMock) SaveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SaveCounter)
}

//SaveMinimockPreCounter returns the value of RequestsExecutorMock.Save invocations
func (m *RequestsExecutorMock) SaveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SavePreCounter)
}

//SaveFinished returns true if mock invocations count is ok
func (m *RequestsExecutorMock) SaveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SaveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SaveCounter) == uint64(len(m.SaveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SaveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SaveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SaveFunc != nil {
		return atomic.LoadUint64(&m.SaveCounter) > 0
	}

	return true
}

type mRequestsExecutorMockSendReply struct {
	mock              *RequestsExecutorMock
	mainExpectation   *RequestsExecutorMockSendReplyExpectation
	expectationSeries []*RequestsExecutorMockSendReplyExpectation
}

type RequestsExecutorMockSendReplyExpectation struct {
	input *RequestsExecutorMockSendReplyInput
}

type RequestsExecutorMockSendReplyInput struct {
	p  context.Context
	p1 *Transcript
	p2 insolar.Reply
	p3 error
}

//Expect specifies that invocation of RequestsExecutor.SendReply is expected from 1 to Infinity times
func (m *mRequestsExecutorMockSendReply) Expect(p context.Context, p1 *Transcript, p2 insolar.Reply, p3 error) *mRequestsExecutorMockSendReply {
	m.mock.SendReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockSendReplyExpectation{}
	}
	m.mainExpectation.input = &RequestsExecutorMockSendReplyInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of RequestsExecutor.SendReply
func (m *mRequestsExecutorMockSendReply) Return() *RequestsExecutorMock {
	m.mock.SendReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RequestsExecutorMockSendReplyExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RequestsExecutor.SendReply is expected once
func (m *mRequestsExecutorMockSendReply) ExpectOnce(p context.Context, p1 *Transcript, p2 insolar.Reply, p3 error) *RequestsExecutorMockSendReplyExpectation {
	m.mock.SendReplyFunc = nil
	m.mainExpectation = nil

	expectation := &RequestsExecutorMockSendReplyExpectation{}
	expectation.input = &RequestsExecutorMockSendReplyInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RequestsExecutor.SendReply method
func (m *mRequestsExecutorMockSendReply) Set(f func(p context.Context, p1 *Transcript, p2 insolar.Reply, p3 error)) *RequestsExecutorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendReplyFunc = f
	return m.mock
}

//SendReply implements github.com/insolar/insolar/logicrunner.RequestsExecutor interface
func (m *RequestsExecutorMock) SendReply(p context.Context, p1 *Transcript, p2 insolar.Reply, p3 error) {
	counter := atomic.AddUint64(&m.SendReplyPreCounter, 1)
	defer atomic.AddUint64(&m.SendReplyCounter, 1)

	if len(m.SendReplyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendReplyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RequestsExecutorMock.SendReply. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendReplyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RequestsExecutorMockSendReplyInput{p, p1, p2, p3}, "RequestsExecutor.SendReply got unexpected parameters")

		return
	}

	if m.SendReplyMock.mainExpectation != nil {

		input := m.SendReplyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RequestsExecutorMockSendReplyInput{p, p1, p2, p3}, "RequestsExecutor.SendReply got unexpected parameters")
		}

		return
	}

	if m.SendReplyFunc == nil {
		m.t.Fatalf("Unexpected call to RequestsExecutorMock.SendReply. %v %v %v %v", p, p1, p2, p3)
		return
	}

	m.SendReplyFunc(p, p1, p2, p3)
}

//SendReplyMinimockCounter returns a count of RequestsExecutorMock.SendReplyFunc invocations
func (m *RequestsExecutorMock) SendReplyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendReplyCounter)
}

//SendReplyMinimockPreCounter returns the value of RequestsExecutorMock.SendReply invocations
func (m *RequestsExecutorMock) SendReplyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendReplyPreCounter)
}

//SendReplyFinished returns true if mock invocations count is ok
func (m *RequestsExecutorMock) SendReplyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendReplyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendReplyCounter) == uint64(len(m.SendReplyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendReplyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendReplyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendReplyFunc != nil {
		return atomic.LoadUint64(&m.SendReplyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestsExecutorMock) ValidateCallCounters() {

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.Execute")
	}

	if !m.ExecuteAndSaveFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.ExecuteAndSave")
	}

	if !m.SaveFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.Save")
	}

	if !m.SendReplyFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.SendReply")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RequestsExecutorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RequestsExecutorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RequestsExecutorMock) MinimockFinish() {

	if !m.ExecuteFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.Execute")
	}

	if !m.ExecuteAndSaveFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.ExecuteAndSave")
	}

	if !m.SaveFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.Save")
	}

	if !m.SendReplyFinished() {
		m.t.Fatal("Expected call to RequestsExecutorMock.SendReply")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RequestsExecutorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RequestsExecutorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ExecuteFinished()
		ok = ok && m.ExecuteAndSaveFinished()
		ok = ok && m.SaveFinished()
		ok = ok && m.SendReplyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ExecuteFinished() {
				m.t.Error("Expected call to RequestsExecutorMock.Execute")
			}

			if !m.ExecuteAndSaveFinished() {
				m.t.Error("Expected call to RequestsExecutorMock.ExecuteAndSave")
			}

			if !m.SaveFinished() {
				m.t.Error("Expected call to RequestsExecutorMock.Save")
			}

			if !m.SendReplyFinished() {
				m.t.Error("Expected call to RequestsExecutorMock.SendReply")
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
func (m *RequestsExecutorMock) AllMocksCalled() bool {

	if !m.ExecuteFinished() {
		return false
	}

	if !m.ExecuteAndSaveFinished() {
		return false
	}

	if !m.SaveFinished() {
		return false
	}

	if !m.SendReplyFinished() {
		return false
	}

	return true
}
