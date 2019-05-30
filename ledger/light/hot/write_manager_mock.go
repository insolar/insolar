package hot

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "WriteManager" can be found in github.com/insolar/insolar/ledger/light/hot
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//WriteManagerMock implements github.com/insolar/insolar/ledger/light/hot.WriteManager
type WriteManagerMock struct {
	t minimock.Tester

	CloseAndWaitFunc       func(p context.Context, p1 insolar.PulseNumber) (r error)
	CloseAndWaitCounter    uint64
	CloseAndWaitPreCounter uint64
	CloseAndWaitMock       mWriteManagerMockCloseAndWait

	OpenFunc       func(p context.Context, p1 insolar.PulseNumber) (r error)
	OpenCounter    uint64
	OpenPreCounter uint64
	OpenMock       mWriteManagerMockOpen
}

//NewWriteManagerMock returns a mock for github.com/insolar/insolar/ledger/light/hot.WriteManager
func NewWriteManagerMock(t minimock.Tester) *WriteManagerMock {
	m := &WriteManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseAndWaitMock = mWriteManagerMockCloseAndWait{mock: m}
	m.OpenMock = mWriteManagerMockOpen{mock: m}

	return m
}

type mWriteManagerMockCloseAndWait struct {
	mock              *WriteManagerMock
	mainExpectation   *WriteManagerMockCloseAndWaitExpectation
	expectationSeries []*WriteManagerMockCloseAndWaitExpectation
}

type WriteManagerMockCloseAndWaitExpectation struct {
	input  *WriteManagerMockCloseAndWaitInput
	result *WriteManagerMockCloseAndWaitResult
}

type WriteManagerMockCloseAndWaitInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type WriteManagerMockCloseAndWaitResult struct {
	r error
}

//Expect specifies that invocation of WriteManager.CloseAndWait is expected from 1 to Infinity times
func (m *mWriteManagerMockCloseAndWait) Expect(p context.Context, p1 insolar.PulseNumber) *mWriteManagerMockCloseAndWait {
	m.mock.CloseAndWaitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WriteManagerMockCloseAndWaitExpectation{}
	}
	m.mainExpectation.input = &WriteManagerMockCloseAndWaitInput{p, p1}
	return m
}

//Return specifies results of invocation of WriteManager.CloseAndWait
func (m *mWriteManagerMockCloseAndWait) Return(r error) *WriteManagerMock {
	m.mock.CloseAndWaitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WriteManagerMockCloseAndWaitExpectation{}
	}
	m.mainExpectation.result = &WriteManagerMockCloseAndWaitResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of WriteManager.CloseAndWait is expected once
func (m *mWriteManagerMockCloseAndWait) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *WriteManagerMockCloseAndWaitExpectation {
	m.mock.CloseAndWaitFunc = nil
	m.mainExpectation = nil

	expectation := &WriteManagerMockCloseAndWaitExpectation{}
	expectation.input = &WriteManagerMockCloseAndWaitInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *WriteManagerMockCloseAndWaitExpectation) Return(r error) {
	e.result = &WriteManagerMockCloseAndWaitResult{r}
}

//Set uses given function f as a mock of WriteManager.CloseAndWait method
func (m *mWriteManagerMockCloseAndWait) Set(f func(p context.Context, p1 insolar.PulseNumber) (r error)) *WriteManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloseAndWaitFunc = f
	return m.mock
}

//CloseAndWait implements github.com/insolar/insolar/ledger/light/hot.WriteManager interface
func (m *WriteManagerMock) CloseAndWait(p context.Context, p1 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.CloseAndWaitPreCounter, 1)
	defer atomic.AddUint64(&m.CloseAndWaitCounter, 1)

	if len(m.CloseAndWaitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloseAndWaitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to WriteManagerMock.CloseAndWait. %v %v", p, p1)
			return
		}

		input := m.CloseAndWaitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, WriteManagerMockCloseAndWaitInput{p, p1}, "WriteManager.CloseAndWait got unexpected parameters")

		result := m.CloseAndWaitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the WriteManagerMock.CloseAndWait")
			return
		}

		r = result.r

		return
	}

	if m.CloseAndWaitMock.mainExpectation != nil {

		input := m.CloseAndWaitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, WriteManagerMockCloseAndWaitInput{p, p1}, "WriteManager.CloseAndWait got unexpected parameters")
		}

		result := m.CloseAndWaitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the WriteManagerMock.CloseAndWait")
		}

		r = result.r

		return
	}

	if m.CloseAndWaitFunc == nil {
		m.t.Fatalf("Unexpected call to WriteManagerMock.CloseAndWait. %v %v", p, p1)
		return
	}

	return m.CloseAndWaitFunc(p, p1)
}

//CloseAndWaitMinimockCounter returns a count of WriteManagerMock.CloseAndWaitFunc invocations
func (m *WriteManagerMock) CloseAndWaitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseAndWaitCounter)
}

//CloseAndWaitMinimockPreCounter returns the value of WriteManagerMock.CloseAndWait invocations
func (m *WriteManagerMock) CloseAndWaitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CloseAndWaitPreCounter)
}

//CloseAndWaitFinished returns true if mock invocations count is ok
func (m *WriteManagerMock) CloseAndWaitFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloseAndWaitMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloseAndWaitCounter) == uint64(len(m.CloseAndWaitMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloseAndWaitMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloseAndWaitCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloseAndWaitFunc != nil {
		return atomic.LoadUint64(&m.CloseAndWaitCounter) > 0
	}

	return true
}

type mWriteManagerMockOpen struct {
	mock              *WriteManagerMock
	mainExpectation   *WriteManagerMockOpenExpectation
	expectationSeries []*WriteManagerMockOpenExpectation
}

type WriteManagerMockOpenExpectation struct {
	input  *WriteManagerMockOpenInput
	result *WriteManagerMockOpenResult
}

type WriteManagerMockOpenInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type WriteManagerMockOpenResult struct {
	r error
}

//Expect specifies that invocation of WriteManager.Open is expected from 1 to Infinity times
func (m *mWriteManagerMockOpen) Expect(p context.Context, p1 insolar.PulseNumber) *mWriteManagerMockOpen {
	m.mock.OpenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WriteManagerMockOpenExpectation{}
	}
	m.mainExpectation.input = &WriteManagerMockOpenInput{p, p1}
	return m
}

//Return specifies results of invocation of WriteManager.Open
func (m *mWriteManagerMockOpen) Return(r error) *WriteManagerMock {
	m.mock.OpenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &WriteManagerMockOpenExpectation{}
	}
	m.mainExpectation.result = &WriteManagerMockOpenResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of WriteManager.Open is expected once
func (m *mWriteManagerMockOpen) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *WriteManagerMockOpenExpectation {
	m.mock.OpenFunc = nil
	m.mainExpectation = nil

	expectation := &WriteManagerMockOpenExpectation{}
	expectation.input = &WriteManagerMockOpenInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *WriteManagerMockOpenExpectation) Return(r error) {
	e.result = &WriteManagerMockOpenResult{r}
}

//Set uses given function f as a mock of WriteManager.Open method
func (m *mWriteManagerMockOpen) Set(f func(p context.Context, p1 insolar.PulseNumber) (r error)) *WriteManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OpenFunc = f
	return m.mock
}

//Open implements github.com/insolar/insolar/ledger/light/hot.WriteManager interface
func (m *WriteManagerMock) Open(p context.Context, p1 insolar.PulseNumber) (r error) {
	counter := atomic.AddUint64(&m.OpenPreCounter, 1)
	defer atomic.AddUint64(&m.OpenCounter, 1)

	if len(m.OpenMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OpenMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to WriteManagerMock.Open. %v %v", p, p1)
			return
		}

		input := m.OpenMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, WriteManagerMockOpenInput{p, p1}, "WriteManager.Open got unexpected parameters")

		result := m.OpenMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the WriteManagerMock.Open")
			return
		}

		r = result.r

		return
	}

	if m.OpenMock.mainExpectation != nil {

		input := m.OpenMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, WriteManagerMockOpenInput{p, p1}, "WriteManager.Open got unexpected parameters")
		}

		result := m.OpenMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the WriteManagerMock.Open")
		}

		r = result.r

		return
	}

	if m.OpenFunc == nil {
		m.t.Fatalf("Unexpected call to WriteManagerMock.Open. %v %v", p, p1)
		return
	}

	return m.OpenFunc(p, p1)
}

//OpenMinimockCounter returns a count of WriteManagerMock.OpenFunc invocations
func (m *WriteManagerMock) OpenMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OpenCounter)
}

//OpenMinimockPreCounter returns the value of WriteManagerMock.Open invocations
func (m *WriteManagerMock) OpenMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OpenPreCounter)
}

//OpenFinished returns true if mock invocations count is ok
func (m *WriteManagerMock) OpenFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.OpenMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.OpenCounter) == uint64(len(m.OpenMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.OpenMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.OpenCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.OpenFunc != nil {
		return atomic.LoadUint64(&m.OpenCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *WriteManagerMock) ValidateCallCounters() {

	if !m.CloseAndWaitFinished() {
		m.t.Fatal("Expected call to WriteManagerMock.CloseAndWait")
	}

	if !m.OpenFinished() {
		m.t.Fatal("Expected call to WriteManagerMock.Open")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *WriteManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *WriteManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *WriteManagerMock) MinimockFinish() {

	if !m.CloseAndWaitFinished() {
		m.t.Fatal("Expected call to WriteManagerMock.CloseAndWait")
	}

	if !m.OpenFinished() {
		m.t.Fatal("Expected call to WriteManagerMock.Open")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *WriteManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *WriteManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloseAndWaitFinished()
		ok = ok && m.OpenFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloseAndWaitFinished() {
				m.t.Error("Expected call to WriteManagerMock.CloseAndWait")
			}

			if !m.OpenFinished() {
				m.t.Error("Expected call to WriteManagerMock.Open")
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
func (m *WriteManagerMock) AllMocksCalled() bool {

	if !m.CloseAndWaitFinished() {
		return false
	}

	if !m.OpenFinished() {
		return false
	}

	return true
}
