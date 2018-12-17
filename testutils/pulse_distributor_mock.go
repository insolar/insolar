package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseDistributor" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseDistributorMock implements github.com/insolar/insolar/core.PulseDistributor
type PulseDistributorMock struct {
	t minimock.Tester

	DistributeFunc       func(p context.Context, p1 *core.Pulse)
	DistributeCounter    uint64
	DistributePreCounter uint64
	DistributeMock       mPulseDistributorMockDistribute
}

//NewPulseDistributorMock returns a mock for github.com/insolar/insolar/core.PulseDistributor
func NewPulseDistributorMock(t minimock.Tester) *PulseDistributorMock {
	m := &PulseDistributorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DistributeMock = mPulseDistributorMockDistribute{mock: m}

	return m
}

type mPulseDistributorMockDistribute struct {
	mock             *PulseDistributorMock
	mockExpectations *PulseDistributorMockDistributeParams
}

//PulseDistributorMockDistributeParams represents input parameters of the PulseDistributor.Distribute
type PulseDistributorMockDistributeParams struct {
	p  context.Context
	p1 *core.Pulse
}

//Expect sets up expected params for the PulseDistributor.Distribute
func (m *mPulseDistributorMockDistribute) Expect(p context.Context, p1 *core.Pulse) *mPulseDistributorMockDistribute {
	m.mockExpectations = &PulseDistributorMockDistributeParams{p, p1}
	return m
}

//Return sets up a mock for PulseDistributor.Distribute to return Return's arguments
func (m *mPulseDistributorMockDistribute) Return() *PulseDistributorMock {
	m.mock.DistributeFunc = func(p context.Context, p1 *core.Pulse) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of PulseDistributor.Distribute method
func (m *mPulseDistributorMockDistribute) Set(f func(p context.Context, p1 *core.Pulse)) *PulseDistributorMock {
	m.mock.DistributeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Distribute implements github.com/insolar/insolar/core.PulseDistributor interface
func (m *PulseDistributorMock) Distribute(p context.Context, p1 *core.Pulse) {
	atomic.AddUint64(&m.DistributePreCounter, 1)
	defer atomic.AddUint64(&m.DistributeCounter, 1)

	if m.DistributeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.DistributeMock.mockExpectations, PulseDistributorMockDistributeParams{p, p1},
			"PulseDistributor.Distribute got unexpected parameters")

		if m.DistributeFunc == nil {

			m.t.Fatal("No results are set for the PulseDistributorMock.Distribute")

			return
		}
	}

	if m.DistributeFunc == nil {
		m.t.Fatal("Unexpected call to PulseDistributorMock.Distribute")
		return
	}

	m.DistributeFunc(p, p1)
}

//DistributeMinimockCounter returns a count of PulseDistributorMock.DistributeFunc invocations
func (m *PulseDistributorMock) DistributeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DistributeCounter)
}

//DistributeMinimockPreCounter returns the value of PulseDistributorMock.Distribute invocations
func (m *PulseDistributorMock) DistributeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DistributePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseDistributorMock) ValidateCallCounters() {

	if m.DistributeFunc != nil && atomic.LoadUint64(&m.DistributeCounter) == 0 {
		m.t.Fatal("Expected call to PulseDistributorMock.Distribute")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseDistributorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseDistributorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseDistributorMock) MinimockFinish() {

	if m.DistributeFunc != nil && atomic.LoadUint64(&m.DistributeCounter) == 0 {
		m.t.Fatal("Expected call to PulseDistributorMock.Distribute")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseDistributorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseDistributorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.DistributeFunc == nil || atomic.LoadUint64(&m.DistributeCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.DistributeFunc != nil && atomic.LoadUint64(&m.DistributeCounter) == 0 {
				m.t.Error("Expected call to PulseDistributorMock.Distribute")
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
func (m *PulseDistributorMock) AllMocksCalled() bool {

	if m.DistributeFunc != nil && atomic.LoadUint64(&m.DistributeCounter) == 0 {
		return false
	}

	return true
}
