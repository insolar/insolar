package object

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "IndexCleaner" can be found in github.com/insolar/insolar/ledger/storage/object
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//IndexCleanerMock implements github.com/insolar/insolar/ledger/storage/object.IndexCleaner
type IndexCleanerMock struct {
	t minimock.Tester

	RemoveForPulseFunc       func(p context.Context, p1 insolar.PulseNumber)
	RemoveForPulseCounter    uint64
	RemoveForPulsePreCounter uint64
	RemoveForPulseMock       mIndexCleanerMockRemoveForPulse
}

//NewIndexCleanerMock returns a mock for github.com/insolar/insolar/ledger/storage/object.IndexCleaner
func NewIndexCleanerMock(t minimock.Tester) *IndexCleanerMock {
	m := &IndexCleanerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveForPulseMock = mIndexCleanerMockRemoveForPulse{mock: m}

	return m
}

type mIndexCleanerMockRemoveForPulse struct {
	mock              *IndexCleanerMock
	mainExpectation   *IndexCleanerMockRemoveForPulseExpectation
	expectationSeries []*IndexCleanerMockRemoveForPulseExpectation
}

type IndexCleanerMockRemoveForPulseExpectation struct {
	input *IndexCleanerMockRemoveForPulseInput
}

type IndexCleanerMockRemoveForPulseInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of IndexCleaner.RemoveForPulse is expected from 1 to Infinity times
func (m *mIndexCleanerMockRemoveForPulse) Expect(p context.Context, p1 insolar.PulseNumber) *mIndexCleanerMockRemoveForPulse {
	m.mock.RemoveForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCleanerMockRemoveForPulseExpectation{}
	}
	m.mainExpectation.input = &IndexCleanerMockRemoveForPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of IndexCleaner.RemoveForPulse
func (m *mIndexCleanerMockRemoveForPulse) Return() *IndexCleanerMock {
	m.mock.RemoveForPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &IndexCleanerMockRemoveForPulseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of IndexCleaner.RemoveForPulse is expected once
func (m *mIndexCleanerMockRemoveForPulse) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *IndexCleanerMockRemoveForPulseExpectation {
	m.mock.RemoveForPulseFunc = nil
	m.mainExpectation = nil

	expectation := &IndexCleanerMockRemoveForPulseExpectation{}
	expectation.input = &IndexCleanerMockRemoveForPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of IndexCleaner.RemoveForPulse method
func (m *mIndexCleanerMockRemoveForPulse) Set(f func(p context.Context, p1 insolar.PulseNumber)) *IndexCleanerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveForPulseFunc = f
	return m.mock
}

//RemoveForPulse implements github.com/insolar/insolar/ledger/storage/object.IndexCleaner interface
func (m *IndexCleanerMock) RemoveForPulse(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.RemoveForPulsePreCounter, 1)
	defer atomic.AddUint64(&m.RemoveForPulseCounter, 1)

	if len(m.RemoveForPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveForPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to IndexCleanerMock.RemoveForPulse. %v %v", p, p1)
			return
		}

		input := m.RemoveForPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, IndexCleanerMockRemoveForPulseInput{p, p1}, "IndexCleaner.RemoveForPulse got unexpected parameters")

		return
	}

	if m.RemoveForPulseMock.mainExpectation != nil {

		input := m.RemoveForPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, IndexCleanerMockRemoveForPulseInput{p, p1}, "IndexCleaner.RemoveForPulse got unexpected parameters")
		}

		return
	}

	if m.RemoveForPulseFunc == nil {
		m.t.Fatalf("Unexpected call to IndexCleanerMock.RemoveForPulse. %v %v", p, p1)
		return
	}

	m.RemoveForPulseFunc(p, p1)
}

//RemoveForPulseMinimockCounter returns a count of IndexCleanerMock.RemoveForPulseFunc invocations
func (m *IndexCleanerMock) RemoveForPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveForPulseCounter)
}

//RemoveForPulseMinimockPreCounter returns the value of IndexCleanerMock.RemoveForPulse invocations
func (m *IndexCleanerMock) RemoveForPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveForPulsePreCounter)
}

//RemoveForPulseFinished returns true if mock invocations count is ok
func (m *IndexCleanerMock) RemoveForPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveForPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveForPulseCounter) == uint64(len(m.RemoveForPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveForPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveForPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveForPulseFunc != nil {
		return atomic.LoadUint64(&m.RemoveForPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCleanerMock) ValidateCallCounters() {

	if !m.RemoveForPulseFinished() {
		m.t.Fatal("Expected call to IndexCleanerMock.RemoveForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *IndexCleanerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *IndexCleanerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *IndexCleanerMock) MinimockFinish() {

	if !m.RemoveForPulseFinished() {
		m.t.Fatal("Expected call to IndexCleanerMock.RemoveForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *IndexCleanerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *IndexCleanerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RemoveForPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RemoveForPulseFinished() {
				m.t.Error("Expected call to IndexCleanerMock.RemoveForPulse")
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
func (m *IndexCleanerMock) AllMocksCalled() bool {

	if !m.RemoveForPulseFinished() {
		return false
	}

	return true
}
