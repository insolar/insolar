package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ArtifactManagerMessageHandler" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ArtifactManagerMessageHandlerMock implements github.com/insolar/insolar/core.ArtifactManagerMessageHandler
type ArtifactManagerMessageHandlerMock struct {
	t minimock.Tester

	OnPulseFunc       func(p context.Context, p1 core.Pulse) (r error)
	OnPulseCounter    uint64
	OnPulsePreCounter uint64
	OnPulseMock       mArtifactManagerMessageHandlerMockOnPulse
}

//NewArtifactManagerMessageHandlerMock returns a mock for github.com/insolar/insolar/core.ArtifactManagerMessageHandler
func NewArtifactManagerMessageHandlerMock(t minimock.Tester) *ArtifactManagerMessageHandlerMock {
	m := &ArtifactManagerMessageHandlerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.OnPulseMock = mArtifactManagerMessageHandlerMockOnPulse{mock: m}

	return m
}

type mArtifactManagerMessageHandlerMockOnPulse struct {
	mock              *ArtifactManagerMessageHandlerMock
	mainExpectation   *ArtifactManagerMessageHandlerMockOnPulseExpectation
	expectationSeries []*ArtifactManagerMessageHandlerMockOnPulseExpectation
}

type ArtifactManagerMessageHandlerMockOnPulseExpectation struct {
	input  *ArtifactManagerMessageHandlerMockOnPulseInput
	result *ArtifactManagerMessageHandlerMockOnPulseResult
}

type ArtifactManagerMessageHandlerMockOnPulseInput struct {
	p  context.Context
	p1 core.Pulse
}

type ArtifactManagerMessageHandlerMockOnPulseResult struct {
	r error
}

//Expect specifies that invocation of ArtifactManagerMessageHandler.OnPulse is expected from 1 to Infinity times
func (m *mArtifactManagerMessageHandlerMockOnPulse) Expect(p context.Context, p1 core.Pulse) *mArtifactManagerMessageHandlerMockOnPulse {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMessageHandlerMockOnPulseExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMessageHandlerMockOnPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of ArtifactManagerMessageHandler.OnPulse
func (m *mArtifactManagerMessageHandlerMockOnPulse) Return(r error) *ArtifactManagerMessageHandlerMock {
	m.mock.OnPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMessageHandlerMockOnPulseExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMessageHandlerMockOnPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManagerMessageHandler.OnPulse is expected once
func (m *mArtifactManagerMessageHandlerMockOnPulse) ExpectOnce(p context.Context, p1 core.Pulse) *ArtifactManagerMessageHandlerMockOnPulseExpectation {
	m.mock.OnPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMessageHandlerMockOnPulseExpectation{}
	expectation.input = &ArtifactManagerMessageHandlerMockOnPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMessageHandlerMockOnPulseExpectation) Return(r error) {
	e.result = &ArtifactManagerMessageHandlerMockOnPulseResult{r}
}

//Set uses given function f as a mock of ArtifactManagerMessageHandler.OnPulse method
func (m *mArtifactManagerMessageHandlerMockOnPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *ArtifactManagerMessageHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.OnPulseFunc = f
	return m.mock
}

//OnPulse implements github.com/insolar/insolar/core.ArtifactManagerMessageHandler interface
func (m *ArtifactManagerMessageHandlerMock) OnPulse(p context.Context, p1 core.Pulse) (r error) {
	counter := atomic.AddUint64(&m.OnPulsePreCounter, 1)
	defer atomic.AddUint64(&m.OnPulseCounter, 1)

	if len(m.OnPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.OnPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMessageHandlerMock.OnPulse. %v %v", p, p1)
			return
		}

		input := m.OnPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMessageHandlerMockOnPulseInput{p, p1}, "ArtifactManagerMessageHandler.OnPulse got unexpected parameters")

		result := m.OnPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMessageHandlerMock.OnPulse")
			return
		}

		r = result.r

		return
	}

	if m.OnPulseMock.mainExpectation != nil {

		input := m.OnPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMessageHandlerMockOnPulseInput{p, p1}, "ArtifactManagerMessageHandler.OnPulse got unexpected parameters")
		}

		result := m.OnPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMessageHandlerMock.OnPulse")
		}

		r = result.r

		return
	}

	if m.OnPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMessageHandlerMock.OnPulse. %v %v", p, p1)
		return
	}

	return m.OnPulseFunc(p, p1)
}

//OnPulseMinimockCounter returns a count of ArtifactManagerMessageHandlerMock.OnPulseFunc invocations
func (m *ArtifactManagerMessageHandlerMock) OnPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulseCounter)
}

//OnPulseMinimockPreCounter returns the value of ArtifactManagerMessageHandlerMock.OnPulse invocations
func (m *ArtifactManagerMessageHandlerMock) OnPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.OnPulsePreCounter)
}

//OnPulseFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMessageHandlerMock) OnPulseFinished() bool {
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
func (m *ArtifactManagerMessageHandlerMock) ValidateCallCounters() {

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMessageHandlerMock.OnPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ArtifactManagerMessageHandlerMock) MinimockFinish() {

	if !m.OnPulseFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMessageHandlerMock.OnPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.OnPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.OnPulseFinished() {
				m.t.Error("Expected call to ArtifactManagerMessageHandlerMock.OnPulse")
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
func (m *ArtifactManagerMessageHandlerMock) AllMocksCalled() bool {

	if !m.OnPulseFinished() {
		return false
	}

	return true
}
