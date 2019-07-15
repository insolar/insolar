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

	AddHotConfirmationFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error)
	AddHotConfirmationCounter    uint64
	AddHotConfirmationPreCounter uint64
	AddHotConfirmationMock       mJetKeeperMockAddHotConfirmation

	AddJetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error)
	AddJetCounter    uint64
	AddJetPreCounter uint64
	AddJetMock       mJetKeeperMockAddJet

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

	m.AddHotConfirmationMock = mJetKeeperMockAddHotConfirmation{mock: m}
	m.AddJetMock = mJetKeeperMockAddJet{mock: m}
	m.TopSyncPulseMock = mJetKeeperMockTopSyncPulse{mock: m}

	return m
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

type mJetKeeperMockAddJet struct {
	mock              *JetKeeperMock
	mainExpectation   *JetKeeperMockAddJetExpectation
	expectationSeries []*JetKeeperMockAddJetExpectation
}

type JetKeeperMockAddJetExpectation struct {
	input  *JetKeeperMockAddJetInput
	result *JetKeeperMockAddJetResult
}

type JetKeeperMockAddJetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type JetKeeperMockAddJetResult struct {
	r error
}

//Expect specifies that invocation of JetKeeper.AddJet is expected from 1 to Infinity times
func (m *mJetKeeperMockAddJet) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mJetKeeperMockAddJet {
	m.mock.AddJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddJetExpectation{}
	}
	m.mainExpectation.input = &JetKeeperMockAddJetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of JetKeeper.AddJet
func (m *mJetKeeperMockAddJet) Return(r error) *JetKeeperMock {
	m.mock.AddJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddJetExpectation{}
	}
	m.mainExpectation.result = &JetKeeperMockAddJetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetKeeper.AddJet is expected once
func (m *mJetKeeperMockAddJet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *JetKeeperMockAddJetExpectation {
	m.mock.AddJetFunc = nil
	m.mainExpectation = nil

	expectation := &JetKeeperMockAddJetExpectation{}
	expectation.input = &JetKeeperMockAddJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetKeeperMockAddJetExpectation) Return(r error) {
	e.result = &JetKeeperMockAddJetResult{r}
}

//Set uses given function f as a mock of JetKeeper.AddJet method
func (m *mJetKeeperMockAddJet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error)) *JetKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddJetFunc = f
	return m.mock
}

//AddJet implements github.com/insolar/insolar/ledger/heavy/executor.JetKeeper interface
func (m *JetKeeperMock) AddJet(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error) {
	counter := atomic.AddUint64(&m.AddJetPreCounter, 1)
	defer atomic.AddUint64(&m.AddJetCounter, 1)

	if len(m.AddJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetKeeperMock.AddJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.AddJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetKeeperMockAddJetInput{p, p1, p2}, "JetKeeper.AddJet got unexpected parameters")

		result := m.AddJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.AddJet")
			return
		}

		r = result.r

		return
	}

	if m.AddJetMock.mainExpectation != nil {

		input := m.AddJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetKeeperMockAddJetInput{p, p1, p2}, "JetKeeper.AddJet got unexpected parameters")
		}

		result := m.AddJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.AddJet")
		}

		r = result.r

		return
	}

	if m.AddJetFunc == nil {
		m.t.Fatalf("Unexpected call to JetKeeperMock.AddJet. %v %v %v", p, p1, p2)
		return
	}

	return m.AddJetFunc(p, p1, p2)
}

//AddJetMinimockCounter returns a count of JetKeeperMock.AddJetFunc invocations
func (m *JetKeeperMock) AddJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddJetCounter)
}

//AddJetMinimockPreCounter returns the value of JetKeeperMock.AddJet invocations
func (m *JetKeeperMock) AddJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddJetPreCounter)
}

//AddJetFinished returns true if mock invocations count is ok
func (m *JetKeeperMock) AddJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddJetCounter) == uint64(len(m.AddJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddJetFunc != nil {
		return atomic.LoadUint64(&m.AddJetCounter) > 0
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

	if !m.AddHotConfirmationFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddHotConfirmation")
	}

	if !m.AddJetFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddJet")
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

	if !m.AddHotConfirmationFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddHotConfirmation")
	}

	if !m.AddJetFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.AddJet")
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
		ok = ok && m.AddHotConfirmationFinished()
		ok = ok && m.AddJetFinished()
		ok = ok && m.TopSyncPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddHotConfirmationFinished() {
				m.t.Error("Expected call to JetKeeperMock.AddHotConfirmation")
			}

			if !m.AddJetFinished() {
				m.t.Error("Expected call to JetKeeperMock.AddJet")
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

	if !m.AddHotConfirmationFinished() {
		return false
	}

	if !m.AddJetFinished() {
		return false
	}

	if !m.TopSyncPulseFinished() {
		return false
	}

	return true
}
