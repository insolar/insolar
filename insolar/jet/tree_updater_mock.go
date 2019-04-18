package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "TreeUpdater" can be found in github.com/insolar/insolar/insolar/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//TreeUpdaterMock implements github.com/insolar/insolar/insolar/jet.TreeUpdater
type TreeUpdaterMock struct {
	t minimock.Tester

	FetchJetFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error)
	FetchJetCounter    uint64
	FetchJetPreCounter uint64
	FetchJetMock       mTreeUpdaterMockFetchJet

	ReleaseJetFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber)
	ReleaseJetCounter    uint64
	ReleaseJetPreCounter uint64
	ReleaseJetMock       mTreeUpdaterMockReleaseJet
}

//NewTreeUpdaterMock returns a mock for github.com/insolar/insolar/insolar/jet.TreeUpdater
func NewTreeUpdaterMock(t minimock.Tester) *TreeUpdaterMock {
	m := &TreeUpdaterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FetchJetMock = mTreeUpdaterMockFetchJet{mock: m}
	m.ReleaseJetMock = mTreeUpdaterMockReleaseJet{mock: m}

	return m
}

type mTreeUpdaterMockFetchJet struct {
	mock              *TreeUpdaterMock
	mainExpectation   *TreeUpdaterMockFetchJetExpectation
	expectationSeries []*TreeUpdaterMockFetchJetExpectation
}

type TreeUpdaterMockFetchJetExpectation struct {
	input  *TreeUpdaterMockFetchJetInput
	result *TreeUpdaterMockFetchJetResult
}

type TreeUpdaterMockFetchJetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type TreeUpdaterMockFetchJetResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of TreeUpdater.FetchJet is expected from 1 to Infinity times
func (m *mTreeUpdaterMockFetchJet) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mTreeUpdaterMockFetchJet {
	m.mock.FetchJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TreeUpdaterMockFetchJetExpectation{}
	}
	m.mainExpectation.input = &TreeUpdaterMockFetchJetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of TreeUpdater.FetchJet
func (m *mTreeUpdaterMockFetchJet) Return(r *insolar.ID, r1 error) *TreeUpdaterMock {
	m.mock.FetchJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TreeUpdaterMockFetchJetExpectation{}
	}
	m.mainExpectation.result = &TreeUpdaterMockFetchJetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of TreeUpdater.FetchJet is expected once
func (m *mTreeUpdaterMockFetchJet) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *TreeUpdaterMockFetchJetExpectation {
	m.mock.FetchJetFunc = nil
	m.mainExpectation = nil

	expectation := &TreeUpdaterMockFetchJetExpectation{}
	expectation.input = &TreeUpdaterMockFetchJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *TreeUpdaterMockFetchJetExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &TreeUpdaterMockFetchJetResult{r, r1}
}

//Set uses given function f as a mock of TreeUpdater.FetchJet method
func (m *mTreeUpdaterMockFetchJet) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error)) *TreeUpdaterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FetchJetFunc = f
	return m.mock
}

//FetchJet implements github.com/insolar/insolar/insolar/jet.TreeUpdater interface
func (m *TreeUpdaterMock) FetchJet(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.FetchJetPreCounter, 1)
	defer atomic.AddUint64(&m.FetchJetCounter, 1)

	if len(m.FetchJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FetchJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TreeUpdaterMock.FetchJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.FetchJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TreeUpdaterMockFetchJetInput{p, p1, p2}, "TreeUpdater.FetchJet got unexpected parameters")

		result := m.FetchJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the TreeUpdaterMock.FetchJet")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FetchJetMock.mainExpectation != nil {

		input := m.FetchJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TreeUpdaterMockFetchJetInput{p, p1, p2}, "TreeUpdater.FetchJet got unexpected parameters")
		}

		result := m.FetchJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the TreeUpdaterMock.FetchJet")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.FetchJetFunc == nil {
		m.t.Fatalf("Unexpected call to TreeUpdaterMock.FetchJet. %v %v %v", p, p1, p2)
		return
	}

	return m.FetchJetFunc(p, p1, p2)
}

//FetchJetMinimockCounter returns a count of TreeUpdaterMock.FetchJetFunc invocations
func (m *TreeUpdaterMock) FetchJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FetchJetCounter)
}

//FetchJetMinimockPreCounter returns the value of TreeUpdaterMock.FetchJet invocations
func (m *TreeUpdaterMock) FetchJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FetchJetPreCounter)
}

//FetchJetFinished returns true if mock invocations count is ok
func (m *TreeUpdaterMock) FetchJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FetchJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FetchJetCounter) == uint64(len(m.FetchJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FetchJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FetchJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FetchJetFunc != nil {
		return atomic.LoadUint64(&m.FetchJetCounter) > 0
	}

	return true
}

type mTreeUpdaterMockReleaseJet struct {
	mock              *TreeUpdaterMock
	mainExpectation   *TreeUpdaterMockReleaseJetExpectation
	expectationSeries []*TreeUpdaterMockReleaseJetExpectation
}

type TreeUpdaterMockReleaseJetExpectation struct {
	input *TreeUpdaterMockReleaseJetInput
}

type TreeUpdaterMockReleaseJetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

//Expect specifies that invocation of TreeUpdater.ReleaseJet is expected from 1 to Infinity times
func (m *mTreeUpdaterMockReleaseJet) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mTreeUpdaterMockReleaseJet {
	m.mock.ReleaseJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TreeUpdaterMockReleaseJetExpectation{}
	}
	m.mainExpectation.input = &TreeUpdaterMockReleaseJetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of TreeUpdater.ReleaseJet
func (m *mTreeUpdaterMockReleaseJet) Return() *TreeUpdaterMock {
	m.mock.ReleaseJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &TreeUpdaterMockReleaseJetExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of TreeUpdater.ReleaseJet is expected once
func (m *mTreeUpdaterMockReleaseJet) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *TreeUpdaterMockReleaseJetExpectation {
	m.mock.ReleaseJetFunc = nil
	m.mainExpectation = nil

	expectation := &TreeUpdaterMockReleaseJetExpectation{}
	expectation.input = &TreeUpdaterMockReleaseJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of TreeUpdater.ReleaseJet method
func (m *mTreeUpdaterMockReleaseJet) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber)) *TreeUpdaterMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReleaseJetFunc = f
	return m.mock
}

//ReleaseJet implements github.com/insolar/insolar/insolar/jet.TreeUpdater interface
func (m *TreeUpdaterMock) ReleaseJet(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.ReleaseJetPreCounter, 1)
	defer atomic.AddUint64(&m.ReleaseJetCounter, 1)

	if len(m.ReleaseJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReleaseJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to TreeUpdaterMock.ReleaseJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.ReleaseJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, TreeUpdaterMockReleaseJetInput{p, p1, p2}, "TreeUpdater.ReleaseJet got unexpected parameters")

		return
	}

	if m.ReleaseJetMock.mainExpectation != nil {

		input := m.ReleaseJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, TreeUpdaterMockReleaseJetInput{p, p1, p2}, "TreeUpdater.ReleaseJet got unexpected parameters")
		}

		return
	}

	if m.ReleaseJetFunc == nil {
		m.t.Fatalf("Unexpected call to TreeUpdaterMock.ReleaseJet. %v %v %v", p, p1, p2)
		return
	}

	m.ReleaseJetFunc(p, p1, p2)
}

//ReleaseJetMinimockCounter returns a count of TreeUpdaterMock.ReleaseJetFunc invocations
func (m *TreeUpdaterMock) ReleaseJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseJetCounter)
}

//ReleaseJetMinimockPreCounter returns the value of TreeUpdaterMock.ReleaseJet invocations
func (m *TreeUpdaterMock) ReleaseJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReleaseJetPreCounter)
}

//ReleaseJetFinished returns true if mock invocations count is ok
func (m *TreeUpdaterMock) ReleaseJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReleaseJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReleaseJetCounter) == uint64(len(m.ReleaseJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReleaseJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReleaseJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReleaseJetFunc != nil {
		return atomic.LoadUint64(&m.ReleaseJetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TreeUpdaterMock) ValidateCallCounters() {

	if !m.FetchJetFinished() {
		m.t.Fatal("Expected call to TreeUpdaterMock.FetchJet")
	}

	if !m.ReleaseJetFinished() {
		m.t.Fatal("Expected call to TreeUpdaterMock.ReleaseJet")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *TreeUpdaterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *TreeUpdaterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *TreeUpdaterMock) MinimockFinish() {

	if !m.FetchJetFinished() {
		m.t.Fatal("Expected call to TreeUpdaterMock.FetchJet")
	}

	if !m.ReleaseJetFinished() {
		m.t.Fatal("Expected call to TreeUpdaterMock.ReleaseJet")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *TreeUpdaterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *TreeUpdaterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.FetchJetFinished()
		ok = ok && m.ReleaseJetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.FetchJetFinished() {
				m.t.Error("Expected call to TreeUpdaterMock.FetchJet")
			}

			if !m.ReleaseJetFinished() {
				m.t.Error("Expected call to TreeUpdaterMock.ReleaseJet")
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
func (m *TreeUpdaterMock) AllMocksCalled() bool {

	if !m.FetchJetFinished() {
		return false
	}

	if !m.ReleaseJetFinished() {
		return false
	}

	return true
}
