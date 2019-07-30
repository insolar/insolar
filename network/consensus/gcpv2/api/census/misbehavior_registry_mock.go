package census

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MisbehaviorRegistry" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/census
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	misbehavior "github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"

	testify_assert "github.com/stretchr/testify/assert"
)

//MisbehaviorRegistryMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.MisbehaviorRegistry
type MisbehaviorRegistryMock struct {
	t minimock.Tester

	AddReportFunc       func(p misbehavior.Report)
	AddReportCounter    uint64
	AddReportPreCounter uint64
	AddReportMock       mMisbehaviorRegistryMockAddReport
}

//NewMisbehaviorRegistryMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/census.MisbehaviorRegistry
func NewMisbehaviorRegistryMock(t minimock.Tester) *MisbehaviorRegistryMock {
	m := &MisbehaviorRegistryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddReportMock = mMisbehaviorRegistryMockAddReport{mock: m}

	return m
}

type mMisbehaviorRegistryMockAddReport struct {
	mock              *MisbehaviorRegistryMock
	mainExpectation   *MisbehaviorRegistryMockAddReportExpectation
	expectationSeries []*MisbehaviorRegistryMockAddReportExpectation
}

type MisbehaviorRegistryMockAddReportExpectation struct {
	input *MisbehaviorRegistryMockAddReportInput
}

type MisbehaviorRegistryMockAddReportInput struct {
	p misbehavior.Report
}

//Expect specifies that invocation of MisbehaviorRegistry.AddReport is expected from 1 to Infinity times
func (m *mMisbehaviorRegistryMockAddReport) Expect(p misbehavior.Report) *mMisbehaviorRegistryMockAddReport {
	m.mock.AddReportFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MisbehaviorRegistryMockAddReportExpectation{}
	}
	m.mainExpectation.input = &MisbehaviorRegistryMockAddReportInput{p}
	return m
}

//Return specifies results of invocation of MisbehaviorRegistry.AddReport
func (m *mMisbehaviorRegistryMockAddReport) Return() *MisbehaviorRegistryMock {
	m.mock.AddReportFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MisbehaviorRegistryMockAddReportExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of MisbehaviorRegistry.AddReport is expected once
func (m *mMisbehaviorRegistryMockAddReport) ExpectOnce(p misbehavior.Report) *MisbehaviorRegistryMockAddReportExpectation {
	m.mock.AddReportFunc = nil
	m.mainExpectation = nil

	expectation := &MisbehaviorRegistryMockAddReportExpectation{}
	expectation.input = &MisbehaviorRegistryMockAddReportInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of MisbehaviorRegistry.AddReport method
func (m *mMisbehaviorRegistryMockAddReport) Set(f func(p misbehavior.Report)) *MisbehaviorRegistryMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddReportFunc = f
	return m.mock
}

//AddReport implements github.com/insolar/insolar/network/consensus/gcpv2/api/census.MisbehaviorRegistry interface
func (m *MisbehaviorRegistryMock) AddReport(p misbehavior.Report) {
	counter := atomic.AddUint64(&m.AddReportPreCounter, 1)
	defer atomic.AddUint64(&m.AddReportCounter, 1)

	if len(m.AddReportMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddReportMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MisbehaviorRegistryMock.AddReport. %v", p)
			return
		}

		input := m.AddReportMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MisbehaviorRegistryMockAddReportInput{p}, "MisbehaviorRegistry.AddReport got unexpected parameters")

		return
	}

	if m.AddReportMock.mainExpectation != nil {

		input := m.AddReportMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MisbehaviorRegistryMockAddReportInput{p}, "MisbehaviorRegistry.AddReport got unexpected parameters")
		}

		return
	}

	if m.AddReportFunc == nil {
		m.t.Fatalf("Unexpected call to MisbehaviorRegistryMock.AddReport. %v", p)
		return
	}

	m.AddReportFunc(p)
}

//AddReportMinimockCounter returns a count of MisbehaviorRegistryMock.AddReportFunc invocations
func (m *MisbehaviorRegistryMock) AddReportMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddReportCounter)
}

//AddReportMinimockPreCounter returns the value of MisbehaviorRegistryMock.AddReport invocations
func (m *MisbehaviorRegistryMock) AddReportMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddReportPreCounter)
}

//AddReportFinished returns true if mock invocations count is ok
func (m *MisbehaviorRegistryMock) AddReportFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddReportMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddReportCounter) == uint64(len(m.AddReportMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddReportMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddReportCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddReportFunc != nil {
		return atomic.LoadUint64(&m.AddReportCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MisbehaviorRegistryMock) ValidateCallCounters() {

	if !m.AddReportFinished() {
		m.t.Fatal("Expected call to MisbehaviorRegistryMock.AddReport")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MisbehaviorRegistryMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MisbehaviorRegistryMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MisbehaviorRegistryMock) MinimockFinish() {

	if !m.AddReportFinished() {
		m.t.Fatal("Expected call to MisbehaviorRegistryMock.AddReport")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MisbehaviorRegistryMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MisbehaviorRegistryMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddReportFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddReportFinished() {
				m.t.Error("Expected call to MisbehaviorRegistryMock.AddReport")
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
func (m *MisbehaviorRegistryMock) AllMocksCalled() bool {

	if !m.AddReportFinished() {
		return false
	}

	return true
}
