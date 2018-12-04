package messagebus

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "tape" can be found in github.com/insolar/insolar/messagebus
*/
import (
	context "context"
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//tapeMock implements github.com/insolar/insolar/messagebus.tape
type tapeMock struct {
	t minimock.Tester

	GetReplyFunc       func(p context.Context, p1 []byte) (r core.Reply, r1 error)
	GetReplyCounter    uint64
	GetReplyPreCounter uint64
	GetReplyMock       mtapeMockGetReply

	SetReplyFunc       func(p context.Context, p1 []byte, p2 core.Reply) (r error)
	SetReplyCounter    uint64
	SetReplyPreCounter uint64
	SetReplyMock       mtapeMockSetReply

	WriteFunc       func(p context.Context, p1 io.Writer) (r error)
	WriteCounter    uint64
	WritePreCounter uint64
	WriteMock       mtapeMockWrite
}

//NewtapeMock returns a mock for github.com/insolar/insolar/messagebus.tape
func NewtapeMock(t minimock.Tester) *tapeMock {
	m := &tapeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetReplyMock = mtapeMockGetReply{mock: m}
	m.SetReplyMock = mtapeMockSetReply{mock: m}
	m.WriteMock = mtapeMockWrite{mock: m}

	return m
}

type mtapeMockGetReply struct {
	mock              *tapeMock
	mainExpectation   *tapeMockGetReplyExpectation
	expectationSeries []*tapeMockGetReplyExpectation
}

type tapeMockGetReplyExpectation struct {
	input  *tapeMockGetReplyInput
	result *tapeMockGetReplyResult
}

type tapeMockGetReplyInput struct {
	p  context.Context
	p1 []byte
}

type tapeMockGetReplyResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of tape.GetReply is expected from 1 to Infinity times
func (m *mtapeMockGetReply) Expect(p context.Context, p1 []byte) *mtapeMockGetReply {
	m.mock.GetReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockGetReplyExpectation{}
	}
	m.mainExpectation.input = &tapeMockGetReplyInput{p, p1}
	return m
}

//Return specifies results of invocation of tape.GetReply
func (m *mtapeMockGetReply) Return(r core.Reply, r1 error) *tapeMock {
	m.mock.GetReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockGetReplyExpectation{}
	}
	m.mainExpectation.result = &tapeMockGetReplyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of tape.GetReply is expected once
func (m *mtapeMockGetReply) ExpectOnce(p context.Context, p1 []byte) *tapeMockGetReplyExpectation {
	m.mock.GetReplyFunc = nil
	m.mainExpectation = nil

	expectation := &tapeMockGetReplyExpectation{}
	expectation.input = &tapeMockGetReplyInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *tapeMockGetReplyExpectation) Return(r core.Reply, r1 error) {
	e.result = &tapeMockGetReplyResult{r, r1}
}

//Set uses given function f as a mock of tape.GetReply method
func (m *mtapeMockGetReply) Set(f func(p context.Context, p1 []byte) (r core.Reply, r1 error)) *tapeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetReplyFunc = f
	return m.mock
}

//GetReply implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) GetReply(p context.Context, p1 []byte) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.GetReplyPreCounter, 1)
	defer atomic.AddUint64(&m.GetReplyCounter, 1)

	if len(m.GetReplyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetReplyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to tapeMock.GetReply. %v %v", p, p1)
			return
		}

		input := m.GetReplyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, tapeMockGetReplyInput{p, p1}, "tape.GetReply got unexpected parameters")

		result := m.GetReplyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.GetReply")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetReplyMock.mainExpectation != nil {

		input := m.GetReplyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, tapeMockGetReplyInput{p, p1}, "tape.GetReply got unexpected parameters")
		}

		result := m.GetReplyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.GetReply")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetReplyFunc == nil {
		m.t.Fatalf("Unexpected call to tapeMock.GetReply. %v %v", p, p1)
		return
	}

	return m.GetReplyFunc(p, p1)
}

//GetReplyMinimockCounter returns a count of tapeMock.GetReplyFunc invocations
func (m *tapeMock) GetReplyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetReplyCounter)
}

//GetReplyMinimockPreCounter returns the value of tapeMock.GetReply invocations
func (m *tapeMock) GetReplyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetReplyPreCounter)
}

//GetReplyFinished returns true if mock invocations count is ok
func (m *tapeMock) GetReplyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetReplyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetReplyCounter) == uint64(len(m.GetReplyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetReplyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetReplyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetReplyFunc != nil {
		return atomic.LoadUint64(&m.GetReplyCounter) > 0
	}

	return true
}

type mtapeMockSetReply struct {
	mock              *tapeMock
	mainExpectation   *tapeMockSetReplyExpectation
	expectationSeries []*tapeMockSetReplyExpectation
}

type tapeMockSetReplyExpectation struct {
	input  *tapeMockSetReplyInput
	result *tapeMockSetReplyResult
}

type tapeMockSetReplyInput struct {
	p  context.Context
	p1 []byte
	p2 core.Reply
}

type tapeMockSetReplyResult struct {
	r error
}

//Expect specifies that invocation of tape.SetReply is expected from 1 to Infinity times
func (m *mtapeMockSetReply) Expect(p context.Context, p1 []byte, p2 core.Reply) *mtapeMockSetReply {
	m.mock.SetReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockSetReplyExpectation{}
	}
	m.mainExpectation.input = &tapeMockSetReplyInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of tape.SetReply
func (m *mtapeMockSetReply) Return(r error) *tapeMock {
	m.mock.SetReplyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockSetReplyExpectation{}
	}
	m.mainExpectation.result = &tapeMockSetReplyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of tape.SetReply is expected once
func (m *mtapeMockSetReply) ExpectOnce(p context.Context, p1 []byte, p2 core.Reply) *tapeMockSetReplyExpectation {
	m.mock.SetReplyFunc = nil
	m.mainExpectation = nil

	expectation := &tapeMockSetReplyExpectation{}
	expectation.input = &tapeMockSetReplyInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *tapeMockSetReplyExpectation) Return(r error) {
	e.result = &tapeMockSetReplyResult{r}
}

//Set uses given function f as a mock of tape.SetReply method
func (m *mtapeMockSetReply) Set(f func(p context.Context, p1 []byte, p2 core.Reply) (r error)) *tapeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetReplyFunc = f
	return m.mock
}

//SetReply implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) SetReply(p context.Context, p1 []byte, p2 core.Reply) (r error) {
	counter := atomic.AddUint64(&m.SetReplyPreCounter, 1)
	defer atomic.AddUint64(&m.SetReplyCounter, 1)

	if len(m.SetReplyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetReplyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to tapeMock.SetReply. %v %v %v", p, p1, p2)
			return
		}

		input := m.SetReplyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, tapeMockSetReplyInput{p, p1, p2}, "tape.SetReply got unexpected parameters")

		result := m.SetReplyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.SetReply")
			return
		}

		r = result.r

		return
	}

	if m.SetReplyMock.mainExpectation != nil {

		input := m.SetReplyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, tapeMockSetReplyInput{p, p1, p2}, "tape.SetReply got unexpected parameters")
		}

		result := m.SetReplyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.SetReply")
		}

		r = result.r

		return
	}

	if m.SetReplyFunc == nil {
		m.t.Fatalf("Unexpected call to tapeMock.SetReply. %v %v %v", p, p1, p2)
		return
	}

	return m.SetReplyFunc(p, p1, p2)
}

//SetReplyMinimockCounter returns a count of tapeMock.SetReplyFunc invocations
func (m *tapeMock) SetReplyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetReplyCounter)
}

//SetReplyMinimockPreCounter returns the value of tapeMock.SetReply invocations
func (m *tapeMock) SetReplyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetReplyPreCounter)
}

//SetReplyFinished returns true if mock invocations count is ok
func (m *tapeMock) SetReplyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetReplyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetReplyCounter) == uint64(len(m.SetReplyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetReplyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetReplyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetReplyFunc != nil {
		return atomic.LoadUint64(&m.SetReplyCounter) > 0
	}

	return true
}

type mtapeMockWrite struct {
	mock              *tapeMock
	mainExpectation   *tapeMockWriteExpectation
	expectationSeries []*tapeMockWriteExpectation
}

type tapeMockWriteExpectation struct {
	input  *tapeMockWriteInput
	result *tapeMockWriteResult
}

type tapeMockWriteInput struct {
	p  context.Context
	p1 io.Writer
}

type tapeMockWriteResult struct {
	r error
}

//Expect specifies that invocation of tape.Write is expected from 1 to Infinity times
func (m *mtapeMockWrite) Expect(p context.Context, p1 io.Writer) *mtapeMockWrite {
	m.mock.WriteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockWriteExpectation{}
	}
	m.mainExpectation.input = &tapeMockWriteInput{p, p1}
	return m
}

//Return specifies results of invocation of tape.Write
func (m *mtapeMockWrite) Return(r error) *tapeMock {
	m.mock.WriteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockWriteExpectation{}
	}
	m.mainExpectation.result = &tapeMockWriteResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of tape.Write is expected once
func (m *mtapeMockWrite) ExpectOnce(p context.Context, p1 io.Writer) *tapeMockWriteExpectation {
	m.mock.WriteFunc = nil
	m.mainExpectation = nil

	expectation := &tapeMockWriteExpectation{}
	expectation.input = &tapeMockWriteInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *tapeMockWriteExpectation) Return(r error) {
	e.result = &tapeMockWriteResult{r}
}

//Set uses given function f as a mock of tape.Write method
func (m *mtapeMockWrite) Set(f func(p context.Context, p1 io.Writer) (r error)) *tapeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteFunc = f
	return m.mock
}

//Write implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) Write(p context.Context, p1 io.Writer) (r error) {
	counter := atomic.AddUint64(&m.WritePreCounter, 1)
	defer atomic.AddUint64(&m.WriteCounter, 1)

	if len(m.WriteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to tapeMock.Write. %v %v", p, p1)
			return
		}

		input := m.WriteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, tapeMockWriteInput{p, p1}, "tape.Write got unexpected parameters")

		result := m.WriteMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.Write")
			return
		}

		r = result.r

		return
	}

	if m.WriteMock.mainExpectation != nil {

		input := m.WriteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, tapeMockWriteInput{p, p1}, "tape.Write got unexpected parameters")
		}

		result := m.WriteMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.Write")
		}

		r = result.r

		return
	}

	if m.WriteFunc == nil {
		m.t.Fatalf("Unexpected call to tapeMock.Write. %v %v", p, p1)
		return
	}

	return m.WriteFunc(p, p1)
}

//WriteMinimockCounter returns a count of tapeMock.WriteFunc invocations
func (m *tapeMock) WriteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteCounter)
}

//WriteMinimockPreCounter returns the value of tapeMock.Write invocations
func (m *tapeMock) WriteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WritePreCounter)
}

//WriteFinished returns true if mock invocations count is ok
func (m *tapeMock) WriteFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.WriteMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.WriteCounter) == uint64(len(m.WriteMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.WriteMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.WriteCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.WriteFunc != nil {
		return atomic.LoadUint64(&m.WriteCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *tapeMock) ValidateCallCounters() {

	if !m.GetReplyFinished() {
		m.t.Fatal("Expected call to tapeMock.GetReply")
	}

	if !m.SetReplyFinished() {
		m.t.Fatal("Expected call to tapeMock.SetReply")
	}

	if !m.WriteFinished() {
		m.t.Fatal("Expected call to tapeMock.Write")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *tapeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *tapeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *tapeMock) MinimockFinish() {

	if !m.GetReplyFinished() {
		m.t.Fatal("Expected call to tapeMock.GetReply")
	}

	if !m.SetReplyFinished() {
		m.t.Fatal("Expected call to tapeMock.SetReply")
	}

	if !m.WriteFinished() {
		m.t.Fatal("Expected call to tapeMock.Write")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *tapeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *tapeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetReplyFinished()
		ok = ok && m.SetReplyFinished()
		ok = ok && m.WriteFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetReplyFinished() {
				m.t.Error("Expected call to tapeMock.GetReply")
			}

			if !m.SetReplyFinished() {
				m.t.Error("Expected call to tapeMock.SetReply")
			}

			if !m.WriteFinished() {
				m.t.Error("Expected call to tapeMock.Write")
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
func (m *tapeMock) AllMocksCalled() bool {

	if !m.GetReplyFinished() {
		return false
	}

	if !m.SetReplyFinished() {
		return false
	}

	if !m.WriteFinished() {
		return false
	}

	return true
}
