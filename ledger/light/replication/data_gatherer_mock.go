package replication

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DataGatherer" can be found in github.com/insolar/insolar/ledger/light/replication
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	message "github.com/insolar/insolar/insolar/message"

	testify_assert "github.com/stretchr/testify/assert"
)

//DataGathererMock implements github.com/insolar/insolar/ledger/light/replication.DataGatherer
type DataGathererMock struct {
	t minimock.Tester

	ForPulseAndJetFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r *message.HeavyPayload, r1 error)
	ForPulseAndJetCounter    uint64
	ForPulseAndJetPreCounter uint64
	ForPulseAndJetMock       mDataGathererMockForPulseAndJet
}

//NewDataGathererMock returns a mock for github.com/insolar/insolar/ledger/light/replication.DataGatherer
func NewDataGathererMock(t minimock.Tester) *DataGathererMock {
	m := &DataGathererMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseAndJetMock = mDataGathererMockForPulseAndJet{mock: m}

	return m
}

type mDataGathererMockForPulseAndJet struct {
	mock              *DataGathererMock
	mainExpectation   *DataGathererMockForPulseAndJetExpectation
	expectationSeries []*DataGathererMockForPulseAndJetExpectation
}

type DataGathererMockForPulseAndJetExpectation struct {
	input  *DataGathererMockForPulseAndJetInput
	result *DataGathererMockForPulseAndJetResult
}

type DataGathererMockForPulseAndJetInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type DataGathererMockForPulseAndJetResult struct {
	r  *message.HeavyPayload
	r1 error
}

//Expect specifies that invocation of DataGatherer.ForPulseAndJet is expected from 1 to Infinity times
func (m *mDataGathererMockForPulseAndJet) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mDataGathererMockForPulseAndJet {
	m.mock.ForPulseAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataGathererMockForPulseAndJetExpectation{}
	}
	m.mainExpectation.input = &DataGathererMockForPulseAndJetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of DataGatherer.ForPulseAndJet
func (m *mDataGathererMockForPulseAndJet) Return(r *message.HeavyPayload, r1 error) *DataGathererMock {
	m.mock.ForPulseAndJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataGathererMockForPulseAndJetExpectation{}
	}
	m.mainExpectation.result = &DataGathererMockForPulseAndJetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DataGatherer.ForPulseAndJet is expected once
func (m *mDataGathererMockForPulseAndJet) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *DataGathererMockForPulseAndJetExpectation {
	m.mock.ForPulseAndJetFunc = nil
	m.mainExpectation = nil

	expectation := &DataGathererMockForPulseAndJetExpectation{}
	expectation.input = &DataGathererMockForPulseAndJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DataGathererMockForPulseAndJetExpectation) Return(r *message.HeavyPayload, r1 error) {
	e.result = &DataGathererMockForPulseAndJetResult{r, r1}
}

//Set uses given function f as a mock of DataGatherer.ForPulseAndJet method
func (m *mDataGathererMockForPulseAndJet) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r *message.HeavyPayload, r1 error)) *DataGathererMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseAndJetFunc = f
	return m.mock
}

//ForPulseAndJet implements github.com/insolar/insolar/ledger/light/replication.DataGatherer interface
func (m *DataGathererMock) ForPulseAndJet(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r *message.HeavyPayload, r1 error) {
	counter := atomic.AddUint64(&m.ForPulseAndJetPreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseAndJetCounter, 1)

	if len(m.ForPulseAndJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseAndJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DataGathererMock.ForPulseAndJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPulseAndJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DataGathererMockForPulseAndJetInput{p, p1, p2}, "DataGatherer.ForPulseAndJet got unexpected parameters")

		result := m.ForPulseAndJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DataGathererMock.ForPulseAndJet")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseAndJetMock.mainExpectation != nil {

		input := m.ForPulseAndJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DataGathererMockForPulseAndJetInput{p, p1, p2}, "DataGatherer.ForPulseAndJet got unexpected parameters")
		}

		result := m.ForPulseAndJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DataGathererMock.ForPulseAndJet")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseAndJetFunc == nil {
		m.t.Fatalf("Unexpected call to DataGathererMock.ForPulseAndJet. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPulseAndJetFunc(p, p1, p2)
}

//ForPulseAndJetMinimockCounter returns a count of DataGathererMock.ForPulseAndJetFunc invocations
func (m *DataGathererMock) ForPulseAndJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseAndJetCounter)
}

//ForPulseAndJetMinimockPreCounter returns the value of DataGathererMock.ForPulseAndJet invocations
func (m *DataGathererMock) ForPulseAndJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseAndJetPreCounter)
}

//ForPulseAndJetFinished returns true if mock invocations count is ok
func (m *DataGathererMock) ForPulseAndJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseAndJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) == uint64(len(m.ForPulseAndJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseAndJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseAndJetFunc != nil {
		return atomic.LoadUint64(&m.ForPulseAndJetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DataGathererMock) ValidateCallCounters() {

	if !m.ForPulseAndJetFinished() {
		m.t.Fatal("Expected call to DataGathererMock.ForPulseAndJet")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DataGathererMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DataGathererMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DataGathererMock) MinimockFinish() {

	if !m.ForPulseAndJetFinished() {
		m.t.Fatal("Expected call to DataGathererMock.ForPulseAndJet")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DataGathererMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DataGathererMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseAndJetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseAndJetFinished() {
				m.t.Error("Expected call to DataGathererMock.ForPulseAndJet")
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
func (m *DataGathererMock) AllMocksCalled() bool {

	if !m.ForPulseAndJetFinished() {
		return false
	}

	return true
}
