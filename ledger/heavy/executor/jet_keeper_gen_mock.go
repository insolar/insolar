package executor

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetKeeper" can be found in github.com/insolar/insolar/ledger/heavy/executor
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetKeeperMock implements github.com/insolar/insolar/ledger/heavy/executor.JetKeeper
type JetKeeperMock struct {
	t minimock.Tester

	AddDropConfirmationFunc       func(p context.Context, p1 insolar.PulseNumber, p2 ...insolar.JetID) (r error)
	AddDropConfirmationCounter    uint64
	AddDropConfirmationPreCounter uint64
	AddDropConfirmationMock       mJetKeeperMockAddDropConfirmation

	AddHotConfirmationFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error)
	AddHotConfirmationCounter    uint64
	AddHotConfirmationPreCounter uint64
	AddHotConfirmationMock       mJetKeeperMockAddHotConfirmation

	TopSyncPulseFunc       func() (r insolar.PulseNumber)
	TopSyncPulseCounter    uint64
	TopSyncPulsePreCounter uint64
	TopSyncPulseMock       mJetKeeperMockTopSyncPulse
}

//NewJetKeeperMock returns a mock for github.com/insolar/insolar/ledger/heavy/executor.JetKeeper
func NewJetKeeperMock(t minimock.Tester) *JetKeeperMock {
	m := &JetKeeperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddDropConfirmationMock = mJetKeeperMockAddDropConfirmation{mock: m}
	m.AddHotConfirmationMock = mJetKeeperMockAddHotConfirmation{mock: m}
	m.TopSyncPulseMock = mJetKeeperMockTopSyncPulse{mock: m}

	return m
}

type mJetKeeperMockAddDropConfirmation struct {
	mock              *JetKeeperMock
	mainExpectation   *JetKeeperMockAddDropConfirmationExpectation
	expectationSeries []*JetKeeperMockAddDropConfirmationExpectation
}

type JetKeeperMockAddDropConfirmationExpectation struct {
	input  *JetKeeperMockAddDropConfirmationInput
	result *JetKeeperMockAddDropConfirmationResult
}

type JetKeeperMockAddDropConfirmationInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 []insolar.JetID
}

type JetKeeperMockAddDropConfirmationResult struct {
	r error
}

//Expect specifies that invocation of JetKeeper.AddDropConfirmation is expected from 1 to Infinity times
func (m *mJetKeeperMockAddDropConfirmation) Expect(p context.Context, p1 insolar.PulseNumber, p2 ...insolar.JetID) *mJetKeeperMockAddDropConfirmation {
	m.mock.AddDropConfirmationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddDropConfirmationExpectation{}
	}
	m.mainExpectation.input = &JetKeeperMockAddDropConfirmationInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of JetKeeper.AddDropConfirmation
func (m *mJetKeeperMockAddDropConfirmation) Return(r error) *JetKeeperMock {
	m.mock.AddDropConfirmationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddDropConfirmationExpectation{}
	}
	m.mainExpectation.result = &JetKeeperMockAddDropConfirmationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetKeeper.AddDropConfirmation is expected once
func (m *mJetKeeperMockAddDropConfirmation) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 ...insolar.JetID) *JetKeeperMockAddDropConfirmationExpectation {
	m.mock.AddDropConfirmationFunc = nil
	m.mainExpectation = nil

	expectation := &JetKeeperMockAddDropConfirmationExpectation{}
	expectation.input = &JetKeeperMockAddDropConfirmationInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetKeeperMockAddDropConfirmationExpectation) Return(r error) {
	e.result = &JetKeeperMockAddDropConfirmationResult{r}
}

//Set uses given function f as a mock of JetKeeper.AddDropConfirmation method
func (m *mJetKeeperMockAddDropConfirmation) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 ...insolar.JetID) (r error)) *JetKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddDropConfirmationFunc = f
	return m.mock
}

//AddDropConfirmation implements github.com/insolar/insolar/ledger/heavy/executor.JetKeeper interface
func (m *JetKeeperMock) AddDropConfirmation(p context.Context, p1 insolar.PulseNumber, p2 ...insolar.JetID) (r error) {
	counter := atomic.AddUint64(&m.AddDropConfirmationPreCounter, 1)
	defer atomic.AddUint64(&m.AddDropConfirmationCounter, 1)

	if len(m.AddDropConfirmationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddDropConfirmationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetKeeperMock.AddDropConfirmation. %v %v %v", p, p1, p2)
			return
		}

		input := m.AddDropConfirmationMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetKeeperMockAddDropConfirmationInput{p, p1, p2}, "JetKeeper.AddDropConfirmation got unexpected parameters")

		result := m.AddDropConfirmationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.AddDropConfirmation")
			return
		}

		r = result.r

		return
	}

	if m.AddDropConfirmationMock.mainExpectation != nil {

		input := m.AddDropConfirmationMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetKeeperMockAddDropConfirmationInput{p, p1, p2}, "JetKeeper.AddDropConfirmation got unexpected parameters")
		}

		result := m.AddDropConfirmationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.AddDropConfirmation")
		}

		r = result.r

		return
	}

	if m.AddDropConfirmationFunc == nil {
		m.t.Fatalf("Unexpected call to JetKeeperMock.AddDropConfirmation. %v %v %v", p, p1, p2)
		return
	}

	return m.AddDropConfirmationFunc(p, p1, p2...)
}

//AddDropConfirmationMinimockCounter returns a count of JetKeeperMock.AddDropConfirmationFunc invocations
func (m *JetKeeperMock) AddDropConfirmationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddDropConfirmationCounter)
}

//AddDropConfirmationMinimockPreCounter returns the value of JetKeeperMock.AddDropConfirmation invocations
func (m *JetKeeperMock) AddDropConfirmationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddDropConfirmationPreCounter)
}

//AddDropConfirmationFinished returns true if mock invocations count is ok
func (m *JetKeeperMock) AddDropConfirmationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddDropConfirmationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddDropConfirmationCounter) == uint64(len(m.AddDropConfirmationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddDropConfirmationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddDropConfirmationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddDropConfirmationFunc != nil {
		return atomic.LoadUint64(&m.AddDropConfirmationCounter) > 0
	}

	return true
}

type mJetKeeperMockAddHotConfirmation struct {
	mock              *JetKeeperMock
	mainExpectation   *JetKeeperMockAddHotConfirmationExpectation
	expectationSeries []*JetKeeperMockAddHotConfirmationExpectation
}

type JetKeeperMockAddHotConfirmationExpectation struct {
	input  *JetKeeperMockAddHotConfirmationInput
	result *JetKeeperMockAddHotConfirmationResult
}

type JetKeeperMockAddHotConfirmationInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type JetKeeperMockAddHotConfirmationResult struct {
	r error
}

//Expect specifies that invocation of JetKeeper.AddHotConfirmation is expected from 1 to Infinity times
func (m *mJetKeeperMockAddHotConfirmation) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mJetKeeperMockAddHotConfirmation {
	m.mock.AddHotConfirmationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddHotConfirmationExpectation{}
	}
	m.mainExpectation.input = &JetKeeperMockAddHotConfirmationInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of JetKeeper.AddHotConfirmation
func (m *mJetKeeperMockAddHotConfirmation) Return(r error) *JetKeeperMock {
	m.mock.AddHotConfirmationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddHotConfirmationExpectation{}
	}
	m.mainExpectation.result = &JetKeeperMockAddHotConfirmationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetKeeper.AddHotConfirmation is expected once
func (m *mJetKeeperMockAddHotConfirmation) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *JetKeeperMockAddHotConfirmationExpectation {
	m.mock.AddHotConfirmationFunc = nil
	m.mainExpectation = nil

	expectation := &JetKeeperMockAddHotConfirmationExpectation{}
	expectation.input = &JetKeeperMockAddHotConfirmationInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetKeeperMockAddHotConfirmationExpectation) Return(r error) {
	e.result = &JetKeeperMockAddHotConfirmationResult{r}
}

//Set uses given function f as a mock of JetKeeper.AddHotConfirmation method
func (m *mJetKeeperMockAddHotConfirmation) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error)) *JetKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddHotConfirmationFunc = f
	return m.mock
}

//AddHotConfirmation implements github.com/insolar/insolar/ledger/heavy/executor.JetKeeper interface
func (m *JetKeeperMock) AddHotConfirmation(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error) {
	counter := atomic.AddUint64(&m.AddHotConfirmationPreCounter, 1)
	defer atomic.AddUint64(&m.AddHotConfirmationCounter, 1)

	if len(m.AddHotConfirmationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddHotConfirmationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetKeeperMock.AddHotConfirmation. %v %v %v", p, p1, p2)
			return
		}

		input := m.AddHotConfirmationMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetKeeperMockAddHotConfirmationInput{p, p1, p2}, "JetKeeper.AddHotConfirmation got unexpected parameters")

		result := m.AddHotConfirmationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.AddHotConfirmation")
			return
		}

		r = result.r

		return
	}

	if m.AddHotConfirmationMock.mainExpectation != nil {

		input := m.AddHotConfirmationMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetKeeperMockAddHotConfirmationInput{p, p1, p2}, "JetKeeper.AddHotConfirmation got unexpected parameters")
		}

		result := m.AddHotConfirmationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.AddHotConfirmation")
		}

		r = result.r

		return
	}

	if m.AddHotConfirmationFunc == nil {
		m.t.Fatalf("Unexpected call to JetKeeperMock.AddHotConfirmation. %v %v %v", p, p1, p2)
		return
	}

	return m.AddHotConfirmationFunc(p, p1, p2)
}

//AddHotConfirmationMinimockCounter returns a count of JetKeeperMock.AddHotConfirmationFunc invocations
func (m *JetKeeperMock) AddHotConfirmationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddHotConfirmationCounter)
}

//AddHotConfirmationMinimockPreCounter returns the value of JetKeeperMock.AddHotConfirmation invocations
func (m *JetKeeperMock) AddHotConfirmationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddHotConfirmationPreCounter)
}

//AddHotConfirmationFinished returns true if mock invocations count is ok
func (m *JetKeeperMock) AddHotConfirmationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddHotConfirmationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddHotConfirmationCounter) == uint64(len(m.AddHotConfirmationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddHotConfirmationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddHotConfirmationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddHotConfirmationFunc != nil {
		return atomic.LoadUint64(&m.AddHotConfirmationCounter) > 0
	}

	return true
}

type mJetKeeperMockTopSyncPulse struct {
	mock              *JetKeeperMock
	mainExpectation   *JetKeeperMockTopSyncPulseExpectation
	expectationSeries []*JetKeeperMockTopSyncPulseExpectation
}

type JetKeeperMockTopSyncPulseExpectation struct {
	result *JetKeeperMockTopSyncPulseResult
}

type JetKeeperMockTopSyncPulseResult struct {
	r insolar.PulseNumber
}

//Expect specifies that invocation of JetKeeper.TopSyncPulse is expected from 1 to Infinity times
func (m *mJetKeeperMockTopSyncPulse) Expect() *mJetKeeperMockTopSyncPulse {
	m.mock.TopSyncPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockTopSyncPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of JetKeeper.TopSyncPulse
func (m *mJetKeeperMockTopSyncPulse) Return(r insolar.PulseNumber) *JetKeeperMock {
	m.mock.TopSyncPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockTopSyncPulseExpectation{}
	}
	m.mainExpectation.result = &JetKeeperMockTopSyncPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetKeeper.TopSyncPulse is expected once
func (m *mJetKeeperMockTopSyncPulse) ExpectOnce() *JetKeeperMockTopSyncPulseExpectation {
	m.mock.TopSyncPulseFunc = nil
	m.mainExpectation = nil

	expectation := &JetKeeperMockTopSyncPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetKeeperMockTopSyncPulseExpectation) Return(r insolar.PulseNumber) {
	e.result = &JetKeeperMockTopSyncPulseResult{r}
}

//Set uses given function f as a mock of JetKeeper.TopSyncPulse method
func (m *mJetKeeperMockTopSyncPulse) Set(f func() (r insolar.PulseNumber)) *JetKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.TopSyncPulseFunc = f
	return m.mock
}

//TopSyncPulse implements github.com/insolar/insolar/ledger/heavy/executor.JetKeeper interface
func (m *JetKeeperMock) TopSyncPulse() (r insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.TopSyncPulsePreCounter, 1)
	defer atomic.AddUint64(&m.TopSyncPulseCounter, 1)

	if len(m.TopSyncPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.TopSyncPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetKeeperMock.TopSyncPulse.")
			return
		}

		result := m.TopSyncPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.TopSyncPulse")
			return
		}

		r = result.r

		return
	}

	if m.TopSyncPulseMock.mainExpectation != nil {

		result := m.TopSyncPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.TopSyncPulse")
		}

		r = result.r

		return
	}

	if m.TopSyncPulseFunc == nil {
		m.t.Fatalf("Unexpected call to JetKeeperMock.TopSyncPulse.")
		return
	}

	return m.TopSyncPulseFunc()
}

//TopSyncPulseMinimockCounter returns a count of JetKeeperMock.TopSyncPulseFunc invocations
func (m *JetKeeperMock) TopSyncPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.TopSyncPulseCounter)
}

//TopSyncPulseMinimockPreCounter returns the value of JetKeeperMock.TopSyncPulse invocations
func (m *JetKeeperMock) TopSyncPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.TopSyncPulsePreCounter)
}

//TopSyncPulseFinished returns true if mock invocations count is ok
func (m *JetKeeperMock) TopSyncPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.TopSyncPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.TopSyncPulseCounter) == uint64(len(m.TopSyncPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.TopSyncPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.TopSyncPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.TopSyncPulseFunc != nil {
		return atomic.LoadUint64(&m.TopSyncPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetKeeperMock) ValidateCallCounters() {

	if !m.AddDropConfirmationFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddDropConfirmation")
	}

	if !m.AddHotConfirmationFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddHotConfirmation")
	}

	if !m.TopSyncPulseFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.TopSyncPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetKeeperMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetKeeperMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetKeeperMock) MinimockFinish() {

	if !m.AddDropConfirmationFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddDropConfirmation")
	}

	if !m.AddHotConfirmationFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddHotConfirmation")
	}

	if !m.TopSyncPulseFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.TopSyncPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetKeeperMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetKeeperMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddDropConfirmationFinished()
		ok = ok && m.AddHotConfirmationFinished()
		ok = ok && m.TopSyncPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddDropConfirmationFinished() {
				m.t.Error("Expected call to JetKeeperMock.AddDropConfirmation")
			}

			if !m.AddHotConfirmationFinished() {
				m.t.Error("Expected call to JetKeeperMock.AddHotConfirmation")
			}

			if !m.TopSyncPulseFinished() {
				m.t.Error("Expected call to JetKeeperMock.TopSyncPulse")
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
func (m *JetKeeperMock) AllMocksCalled() bool {

	if !m.AddDropConfirmationFinished() {
		return false
	}

	if !m.AddHotConfirmationFinished() {
		return false
	}

	if !m.TopSyncPulseFinished() {
		return false
	}

	return true
}
