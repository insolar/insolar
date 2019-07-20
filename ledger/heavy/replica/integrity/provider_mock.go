package integrity

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Provider" can be found in github.com/insolar/insolar/ledger/heavy/replica/integrity
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	sequence "github.com/insolar/insolar/ledger/heavy/sequence"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProviderMock implements github.com/insolar/insolar/ledger/heavy/replica/integrity.Provider
type ProviderMock struct {
	t minimock.Tester

	WrapFunc       func(p []sequence.Item) (r []byte)
	WrapCounter    uint64
	WrapPreCounter uint64
	WrapMock       mProviderMockWrap
}

//NewProviderMock returns a mock for github.com/insolar/insolar/ledger/heavy/replica/integrity.Provider
func NewProviderMock(t minimock.Tester) *ProviderMock {
	m := &ProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.WrapMock = mProviderMockWrap{mock: m}

	return m
}

type mProviderMockWrap struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockWrapExpectation
	expectationSeries []*ProviderMockWrapExpectation
}

type ProviderMockWrapExpectation struct {
	input  *ProviderMockWrapInput
	result *ProviderMockWrapResult
}

type ProviderMockWrapInput struct {
	p []sequence.Item
}

type ProviderMockWrapResult struct {
	r []byte
}

//Expect specifies that invocation of Provider.Wrap is expected from 1 to Infinity times
func (m *mProviderMockWrap) Expect(p []sequence.Item) *mProviderMockWrap {
	m.mock.WrapFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockWrapExpectation{}
	}
	m.mainExpectation.input = &ProviderMockWrapInput{p}
	return m
}

//Return specifies results of invocation of Provider.Wrap
func (m *mProviderMockWrap) Return(r []byte) *ProviderMock {
	m.mock.WrapFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockWrapExpectation{}
	}
	m.mainExpectation.result = &ProviderMockWrapResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.Wrap is expected once
func (m *mProviderMockWrap) ExpectOnce(p []sequence.Item) *ProviderMockWrapExpectation {
	m.mock.WrapFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockWrapExpectation{}
	expectation.input = &ProviderMockWrapInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockWrapExpectation) Return(r []byte) {
	e.result = &ProviderMockWrapResult{r}
}

//Set uses given function f as a mock of Provider.Wrap method
func (m *mProviderMockWrap) Set(f func(p []sequence.Item) (r []byte)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WrapFunc = f
	return m.mock
}

//Wrap implements github.com/insolar/insolar/ledger/heavy/replica/integrity.Provider interface
func (m *ProviderMock) Wrap(p []sequence.Item) (r []byte) {
	counter := atomic.AddUint64(&m.WrapPreCounter, 1)
	defer atomic.AddUint64(&m.WrapCounter, 1)

	if len(m.WrapMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WrapMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.Wrap. %v", p)
			return
		}

		input := m.WrapMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockWrapInput{p}, "Provider.Wrap got unexpected parameters")

		result := m.WrapMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.Wrap")
			return
		}

		r = result.r

		return
	}

	if m.WrapMock.mainExpectation != nil {

		input := m.WrapMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockWrapInput{p}, "Provider.Wrap got unexpected parameters")
		}

		result := m.WrapMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.Wrap")
		}

		r = result.r

		return
	}

	if m.WrapFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.Wrap. %v", p)
		return
	}

	return m.WrapFunc(p)
}

//WrapMinimockCounter returns a count of ProviderMock.WrapFunc invocations
func (m *ProviderMock) WrapMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WrapCounter)
}

//WrapMinimockPreCounter returns the value of ProviderMock.Wrap invocations
func (m *ProviderMock) WrapMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WrapPreCounter)
}

//WrapFinished returns true if mock invocations count is ok
func (m *ProviderMock) WrapFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.WrapMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.WrapCounter) == uint64(len(m.WrapMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.WrapMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.WrapCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.WrapFunc != nil {
		return atomic.LoadUint64(&m.WrapCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) ValidateCallCounters() {

	if !m.WrapFinished() {
		m.t.Fatal("Expected call to ProviderMock.Wrap")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ProviderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ProviderMock) MinimockFinish() {

	if !m.WrapFinished() {
		m.t.Fatal("Expected call to ProviderMock.Wrap")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ProviderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ProviderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.WrapFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.WrapFinished() {
				m.t.Error("Expected call to ProviderMock.Wrap")
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
func (m *ProviderMock) AllMocksCalled() bool {

	if !m.WrapFinished() {
		return false
	}

	return true
}
