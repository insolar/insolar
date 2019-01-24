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

	GetFunc       func(p context.Context, p1 []byte) (r *TapeItem, r1 error)
	GetCounter    uint64
	GetPreCounter uint64
	GetMock       mtapeMockGet

	SetFunc       func(p context.Context, p1 []byte, p2 core.Reply, p3 error) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mtapeMockSet

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

	m.GetMock = mtapeMockGet{mock: m}
	m.SetMock = mtapeMockSet{mock: m}
	m.WriteMock = mtapeMockWrite{mock: m}

	return m
}

type mtapeMockGet struct {
	mock              *tapeMock
	mainExpectation   *tapeMockGetExpectation
	expectationSeries []*tapeMockGetExpectation
}

type tapeMockGetExpectation struct {
	input  *tapeMockGetInput
	result *tapeMockGetResult
}

type tapeMockGetInput struct {
	p  context.Context
	p1 []byte
}

type tapeMockGetResult struct {
	r  *TapeItem
	r1 error
}

//Expect specifies that invocation of tape.Get is expected from 1 to Infinity times
func (m *mtapeMockGet) Expect(p context.Context, p1 []byte) *mtapeMockGet {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockGetExpectation{}
	}
	m.mainExpectation.input = &tapeMockGetInput{p, p1}
	return m
}

//Return specifies results of invocation of tape.Get
func (m *mtapeMockGet) Return(r *TapeItem, r1 error) *tapeMock {
	m.mock.GetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockGetExpectation{}
	}
	m.mainExpectation.result = &tapeMockGetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of tape.Get is expected once
func (m *mtapeMockGet) ExpectOnce(p context.Context, p1 []byte) *tapeMockGetExpectation {
	m.mock.GetFunc = nil
	m.mainExpectation = nil

	expectation := &tapeMockGetExpectation{}
	expectation.input = &tapeMockGetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *tapeMockGetExpectation) Return(r *TapeItem, r1 error) {
	e.result = &tapeMockGetResult{r, r1}
}

//Set uses given function f as a mock of tape.Get method
func (m *mtapeMockGet) Set(f func(p context.Context, p1 []byte) (r *TapeItem, r1 error)) *tapeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetFunc = f
	return m.mock
}

//Get implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) Get(p context.Context, p1 []byte) (r *TapeItem, r1 error) {
	counter := atomic.AddUint64(&m.GetPreCounter, 1)
	defer atomic.AddUint64(&m.GetCounter, 1)

	if len(m.GetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to tapeMock.Get. %v %v", p, p1)
			return
		}

		input := m.GetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, tapeMockGetInput{p, p1}, "tape.Get got unexpected parameters")

		result := m.GetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.Get")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMock.mainExpectation != nil {

		input := m.GetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, tapeMockGetInput{p, p1}, "tape.Get got unexpected parameters")
		}

		result := m.GetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.Get")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetFunc == nil {
		m.t.Fatalf("Unexpected call to tapeMock.Get. %v %v", p, p1)
		return
	}

	return m.GetFunc(p, p1)
}

//GetMinimockCounter returns a count of tapeMock.GetFunc invocations
func (m *tapeMock) GetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCounter)
}

//GetMinimockPreCounter returns the value of tapeMock.Get invocations
func (m *tapeMock) GetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreCounter)
}

//GetFinished returns true if mock invocations count is ok
func (m *tapeMock) GetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCounter) == uint64(len(m.GetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetFunc != nil {
		return atomic.LoadUint64(&m.GetCounter) > 0
	}

	return true
}

type mtapeMockSet struct {
	mock              *tapeMock
	mainExpectation   *tapeMockSetExpectation
	expectationSeries []*tapeMockSetExpectation
}

type tapeMockSetExpectation struct {
	input  *tapeMockSetInput
	result *tapeMockSetResult
}

type tapeMockSetInput struct {
	p  context.Context
	p1 []byte
	p2 core.Reply
	p3 error
}

type tapeMockSetResult struct {
	r error
}

//Expect specifies that invocation of tape.Set is expected from 1 to Infinity times
func (m *mtapeMockSet) Expect(p context.Context, p1 []byte, p2 core.Reply, p3 error) *mtapeMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockSetExpectation{}
	}
	m.mainExpectation.input = &tapeMockSetInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of tape.Set
func (m *mtapeMockSet) Return(r error) *tapeMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &tapeMockSetExpectation{}
	}
	m.mainExpectation.result = &tapeMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of tape.Set is expected once
func (m *mtapeMockSet) ExpectOnce(p context.Context, p1 []byte, p2 core.Reply, p3 error) *tapeMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &tapeMockSetExpectation{}
	expectation.input = &tapeMockSetInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *tapeMockSetExpectation) Return(r error) {
	e.result = &tapeMockSetResult{r}
}

//Set uses given function f as a mock of tape.Set method
func (m *mtapeMockSet) Set(f func(p context.Context, p1 []byte, p2 core.Reply, p3 error) (r error)) *tapeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) Set(p context.Context, p1 []byte, p2 core.Reply, p3 error) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to tapeMock.Set. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, tapeMockSetInput{p, p1, p2, p3}, "tape.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, tapeMockSetInput{p, p1, p2, p3}, "tape.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the tapeMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to tapeMock.Set. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SetFunc(p, p1, p2, p3)
}

//SetMinimockCounter returns a count of tapeMock.SetFunc invocations
func (m *tapeMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of tapeMock.Set invocations
func (m *tapeMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *tapeMock) SetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCounter) == uint64(len(m.SetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFunc != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
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

	if !m.GetFinished() {
		m.t.Fatal("Expected call to tapeMock.Get")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to tapeMock.Set")
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

	if !m.GetFinished() {
		m.t.Fatal("Expected call to tapeMock.Get")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to tapeMock.Set")
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
		ok = ok && m.GetFinished()
		ok = ok && m.SetFinished()
		ok = ok && m.WriteFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetFinished() {
				m.t.Error("Expected call to tapeMock.Get")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to tapeMock.Set")
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

	if !m.GetFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	if !m.WriteFinished() {
		return false
	}

	return true
}
