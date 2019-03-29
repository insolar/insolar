package blob

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CollectionAccessor" can be found in github.com/insolar/insolar/ledger/storage/blob
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//SyncAccessorMock implements github.com/insolar/insolar/ledger/storage/blob.CollectionAccessor
type SyncAccessorMock struct {
	t minimock.Tester

	ForPNFunc       func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []Blob)
	ForPNCounter    uint64
	ForPNPreCounter uint64
	ForPNMock       mSyncAccessorMockForPN
}

//NewSyncAccessorMock returns a mock for github.com/insolar/insolar/ledger/storage/blob.CollectionAccessor
func NewSyncAccessorMock(t minimock.Tester) *SyncAccessorMock {
	m := &SyncAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPNMock = mSyncAccessorMockForPN{mock: m}

	return m
}

type mSyncAccessorMockForPN struct {
	mock              *SyncAccessorMock
	mainExpectation   *SyncAccessorMockForPNExpectation
	expectationSeries []*SyncAccessorMockForPNExpectation
}

type SyncAccessorMockForPNExpectation struct {
	input  *SyncAccessorMockForPNInput
	result *SyncAccessorMockForPNResult
}

type SyncAccessorMockForPNInput struct {
	p  context.Context
	p1 insolar.JetID
	p2 insolar.PulseNumber
}

type SyncAccessorMockForPNResult struct {
	r []Blob
}

//Expect specifies that invocation of CollectionAccessor.ForPulse is expected from 1 to Infinity times
func (m *mSyncAccessorMockForPN) Expect(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *mSyncAccessorMockForPN {
	m.mock.ForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SyncAccessorMockForPNExpectation{}
	}
	m.mainExpectation.input = &SyncAccessorMockForPNInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of CollectionAccessor.ForPulse
func (m *mSyncAccessorMockForPN) Return(r []Blob) *SyncAccessorMock {
	m.mock.ForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SyncAccessorMockForPNExpectation{}
	}
	m.mainExpectation.result = &SyncAccessorMockForPNResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CollectionAccessor.ForPulse is expected once
func (m *mSyncAccessorMockForPN) ExpectOnce(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) *SyncAccessorMockForPNExpectation {
	m.mock.ForPNFunc = nil
	m.mainExpectation = nil

	expectation := &SyncAccessorMockForPNExpectation{}
	expectation.input = &SyncAccessorMockForPNInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SyncAccessorMockForPNExpectation) Return(r []Blob) {
	e.result = &SyncAccessorMockForPNResult{r}
}

//Set uses given function f as a mock of CollectionAccessor.ForPulse method
func (m *mSyncAccessorMockForPN) Set(f func(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []Blob)) *SyncAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPNFunc = f
	return m.mock
}

//ForPulse implements github.com/insolar/insolar/ledger/storage/blob.CollectionAccessor interface
func (m *SyncAccessorMock) ForPN(p context.Context, p1 insolar.JetID, p2 insolar.PulseNumber) (r []Blob) {
	counter := atomic.AddUint64(&m.ForPNPreCounter, 1)
	defer atomic.AddUint64(&m.ForPNCounter, 1)

	if len(m.ForPNMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPNMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SyncAccessorMock.ForPulse. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForPNMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SyncAccessorMockForPNInput{p, p1, p2}, "CollectionAccessor.ForPulse got unexpected parameters")

		result := m.ForPNMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SyncAccessorMock.ForPulse")
			return
		}

		r = result.r

		return
	}

	if m.ForPNMock.mainExpectation != nil {

		input := m.ForPNMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SyncAccessorMockForPNInput{p, p1, p2}, "CollectionAccessor.ForPulse got unexpected parameters")
		}

		result := m.ForPNMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SyncAccessorMock.ForPulse")
		}

		r = result.r

		return
	}

	if m.ForPNFunc == nil {
		m.t.Fatalf("Unexpected call to SyncAccessorMock.ForPulse. %v %v %v", p, p1, p2)
		return
	}

	return m.ForPNFunc(p, p1, p2)
}

//ForPNMinimockCounter returns a count of SyncAccessorMock.ForPNFunc invocations
func (m *SyncAccessorMock) ForPNMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNCounter)
}

//ForPNMinimockPreCounter returns the value of SyncAccessorMock.ForPulse invocations
func (m *SyncAccessorMock) ForPNMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPNPreCounter)
}

//ForPNFinished returns true if mock invocations count is ok
func (m *SyncAccessorMock) ForPNFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPNMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPNCounter) == uint64(len(m.ForPNMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPNMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPNCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPNFunc != nil {
		return atomic.LoadUint64(&m.ForPNCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SyncAccessorMock) ValidateCallCounters() {

	if !m.ForPNFinished() {
		m.t.Fatal("Expected call to SyncAccessorMock.ForPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SyncAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SyncAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SyncAccessorMock) MinimockFinish() {

	if !m.ForPNFinished() {
		m.t.Fatal("Expected call to SyncAccessorMock.ForPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SyncAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SyncAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPNFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPNFinished() {
				m.t.Error("Expected call to SyncAccessorMock.ForPulse")
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
func (m *SyncAccessorMock) AllMocksCalled() bool {

	if !m.ForPNFinished() {
		return false
	}

	return true
}
