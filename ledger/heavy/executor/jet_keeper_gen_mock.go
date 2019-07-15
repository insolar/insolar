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

	AddFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error)
	AddCounter    uint64
	AddPreCounter uint64
	AddMock       mJetKeeperMockAdd

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

	m.AddMock = mJetKeeperMockAdd{mock: m}
	m.TopSyncPulseMock = mJetKeeperMockTopSyncPulse{mock: m}

	return m
}

type mJetKeeperMockAdd struct {
	mock              *JetKeeperMock
	mainExpectation   *JetKeeperMockAddExpectation
	expectationSeries []*JetKeeperMockAddExpectation
}

type JetKeeperMockAddExpectation struct {
	input  *JetKeeperMockAddInput
	result *JetKeeperMockAddResult
}

type JetKeeperMockAddInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type JetKeeperMockAddResult struct {
	r error
}

//Expect specifies that invocation of JetKeeper.Add is expected from 1 to Infinity times
func (m *mJetKeeperMockAdd) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mJetKeeperMockAdd {
	m.mock.AddFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddExpectation{}
	}
	m.mainExpectation.input = &JetKeeperMockAddInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of JetKeeper.Add
func (m *mJetKeeperMockAdd) Return(r error) *JetKeeperMock {
	m.mock.AddFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetKeeperMockAddExpectation{}
	}
	m.mainExpectation.result = &JetKeeperMockAddResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of JetKeeper.Add is expected once
func (m *mJetKeeperMockAdd) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *JetKeeperMockAddExpectation {
	m.mock.AddFunc = nil
	m.mainExpectation = nil

	expectation := &JetKeeperMockAddExpectation{}
	expectation.input = &JetKeeperMockAddInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *JetKeeperMockAddExpectation) Return(r error) {
	e.result = &JetKeeperMockAddResult{r}
}

//Set uses given function f as a mock of JetKeeper.Add method
func (m *mJetKeeperMockAdd) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error)) *JetKeeperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddFunc = f
	return m.mock
}

//Add implements github.com/insolar/insolar/ledger/heavy/executor.JetKeeper interface
func (m *JetKeeperMock) Add(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r error) {
	counter := atomic.AddUint64(&m.AddPreCounter, 1)
	defer atomic.AddUint64(&m.AddCounter, 1)

	if len(m.AddMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetKeeperMock.Add. %v %v %v", p, p1, p2)
			return
		}

		input := m.AddMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetKeeperMockAddInput{p, p1, p2}, "JetKeeper.Add got unexpected parameters")

		result := m.AddMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.Add")
			return
		}

		r = result.r

		return
	}

	if m.AddMock.mainExpectation != nil {

		input := m.AddMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetKeeperMockAddInput{p, p1, p2}, "JetKeeper.Add got unexpected parameters")
		}

		result := m.AddMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the JetKeeperMock.Add")
		}

		r = result.r

		return
	}

	if m.AddFunc == nil {
		m.t.Fatalf("Unexpected call to JetKeeperMock.Add. %v %v %v", p, p1, p2)
		return
	}

	return m.AddFunc(p, p1, p2)
}

//AddMinimockCounter returns a count of JetKeeperMock.AddFunc invocations
func (m *JetKeeperMock) AddMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddCounter)
}

//AddMinimockPreCounter returns the value of JetKeeperMock.Add invocations
func (m *JetKeeperMock) AddMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddPreCounter)
}

//AddFinished returns true if mock invocations count is ok
func (m *JetKeeperMock) AddFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddCounter) == uint64(len(m.AddMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddFunc != nil {
		return atomic.LoadUint64(&m.AddCounter) > 0
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

	if !m.AddFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.Add")
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

	if !m.AddFinished() {
		m.t.Fatal("Expected call to JetKeeperMock.Add")
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
		ok = ok && m.AddFinished()
		ok = ok && m.TopSyncPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddFinished() {
				m.t.Error("Expected call to JetKeeperMock.Add")
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

	if !m.AddFinished() {
		return false
	}

	if !m.TopSyncPulseFinished() {
		return false
	}

	return true
}
