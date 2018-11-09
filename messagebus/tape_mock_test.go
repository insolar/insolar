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
	mock             *tapeMock
	mockExpectations *tapeMockGetReplyParams
}

//tapeMockGetReplyParams represents input parameters of the tape.GetReply
type tapeMockGetReplyParams struct {
	p  context.Context
	p1 []byte
}

//Expect sets up expected params for the tape.GetReply
func (m *mtapeMockGetReply) Expect(p context.Context, p1 []byte) *mtapeMockGetReply {
	m.mockExpectations = &tapeMockGetReplyParams{p, p1}
	return m
}

//Return sets up a mock for tape.GetReply to return Return's arguments
func (m *mtapeMockGetReply) Return(r core.Reply, r1 error) *tapeMock {
	m.mock.GetReplyFunc = func(p context.Context, p1 []byte) (core.Reply, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of tape.GetReply method
func (m *mtapeMockGetReply) Set(f func(p context.Context, p1 []byte) (r core.Reply, r1 error)) *tapeMock {
	m.mock.GetReplyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetReply implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) GetReply(p context.Context, p1 []byte) (r core.Reply, r1 error) {
	atomic.AddUint64(&m.GetReplyPreCounter, 1)
	defer atomic.AddUint64(&m.GetReplyCounter, 1)

	if m.GetReplyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetReplyMock.mockExpectations, tapeMockGetReplyParams{p, p1},
			"tape.GetReply got unexpected parameters")

		if m.GetReplyFunc == nil {

			m.t.Fatal("No results are set for the tapeMock.GetReply")

			return
		}
	}

	if m.GetReplyFunc == nil {
		m.t.Fatal("Unexpected call to tapeMock.GetReply")
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

type mtapeMockSetReply struct {
	mock             *tapeMock
	mockExpectations *tapeMockSetReplyParams
}

//tapeMockSetReplyParams represents input parameters of the tape.SetReply
type tapeMockSetReplyParams struct {
	p  context.Context
	p1 []byte
	p2 core.Reply
}

//Expect sets up expected params for the tape.SetReply
func (m *mtapeMockSetReply) Expect(p context.Context, p1 []byte, p2 core.Reply) *mtapeMockSetReply {
	m.mockExpectations = &tapeMockSetReplyParams{p, p1, p2}
	return m
}

//Return sets up a mock for tape.SetReply to return Return's arguments
func (m *mtapeMockSetReply) Return(r error) *tapeMock {
	m.mock.SetReplyFunc = func(p context.Context, p1 []byte, p2 core.Reply) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of tape.SetReply method
func (m *mtapeMockSetReply) Set(f func(p context.Context, p1 []byte, p2 core.Reply) (r error)) *tapeMock {
	m.mock.SetReplyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetReply implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) SetReply(p context.Context, p1 []byte, p2 core.Reply) (r error) {
	atomic.AddUint64(&m.SetReplyPreCounter, 1)
	defer atomic.AddUint64(&m.SetReplyCounter, 1)

	if m.SetReplyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetReplyMock.mockExpectations, tapeMockSetReplyParams{p, p1, p2},
			"tape.SetReply got unexpected parameters")

		if m.SetReplyFunc == nil {

			m.t.Fatal("No results are set for the tapeMock.SetReply")

			return
		}
	}

	if m.SetReplyFunc == nil {
		m.t.Fatal("Unexpected call to tapeMock.SetReply")
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

type mtapeMockWrite struct {
	mock             *tapeMock
	mockExpectations *tapeMockWriteParams
}

//tapeMockWriteParams represents input parameters of the tape.Write
type tapeMockWriteParams struct {
	p  context.Context
	p1 io.Writer
}

//Expect sets up expected params for the tape.Write
func (m *mtapeMockWrite) Expect(p context.Context, p1 io.Writer) *mtapeMockWrite {
	m.mockExpectations = &tapeMockWriteParams{p, p1}
	return m
}

//Return sets up a mock for tape.Write to return Return's arguments
func (m *mtapeMockWrite) Return(r error) *tapeMock {
	m.mock.WriteFunc = func(p context.Context, p1 io.Writer) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of tape.Write method
func (m *mtapeMockWrite) Set(f func(p context.Context, p1 io.Writer) (r error)) *tapeMock {
	m.mock.WriteFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Write implements github.com/insolar/insolar/messagebus.tape interface
func (m *tapeMock) Write(p context.Context, p1 io.Writer) (r error) {
	atomic.AddUint64(&m.WritePreCounter, 1)
	defer atomic.AddUint64(&m.WriteCounter, 1)

	if m.WriteMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.WriteMock.mockExpectations, tapeMockWriteParams{p, p1},
			"tape.Write got unexpected parameters")

		if m.WriteFunc == nil {

			m.t.Fatal("No results are set for the tapeMock.Write")

			return
		}
	}

	if m.WriteFunc == nil {
		m.t.Fatal("Unexpected call to tapeMock.Write")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *tapeMock) ValidateCallCounters() {

	if m.GetReplyFunc != nil && atomic.LoadUint64(&m.GetReplyCounter) == 0 {
		m.t.Fatal("Expected call to tapeMock.GetReply")
	}

	if m.SetReplyFunc != nil && atomic.LoadUint64(&m.SetReplyCounter) == 0 {
		m.t.Fatal("Expected call to tapeMock.SetReply")
	}

	if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
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

	if m.GetReplyFunc != nil && atomic.LoadUint64(&m.GetReplyCounter) == 0 {
		m.t.Fatal("Expected call to tapeMock.GetReply")
	}

	if m.SetReplyFunc != nil && atomic.LoadUint64(&m.SetReplyCounter) == 0 {
		m.t.Fatal("Expected call to tapeMock.SetReply")
	}

	if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
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
		ok = ok && (m.GetReplyFunc == nil || atomic.LoadUint64(&m.GetReplyCounter) > 0)
		ok = ok && (m.SetReplyFunc == nil || atomic.LoadUint64(&m.SetReplyCounter) > 0)
		ok = ok && (m.WriteFunc == nil || atomic.LoadUint64(&m.WriteCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetReplyFunc != nil && atomic.LoadUint64(&m.GetReplyCounter) == 0 {
				m.t.Error("Expected call to tapeMock.GetReply")
			}

			if m.SetReplyFunc != nil && atomic.LoadUint64(&m.SetReplyCounter) == 0 {
				m.t.Error("Expected call to tapeMock.SetReply")
			}

			if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
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

	if m.GetReplyFunc != nil && atomic.LoadUint64(&m.GetReplyCounter) == 0 {
		return false
	}

	if m.SetReplyFunc != nil && atomic.LoadUint64(&m.SetReplyCounter) == 0 {
		return false
	}

	if m.WriteFunc != nil && atomic.LoadUint64(&m.WriteCounter) == 0 {
		return false
	}

	return true
}
