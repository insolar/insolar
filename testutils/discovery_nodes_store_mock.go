package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DiscoveryNodesStore" can be found in github.com/insolar/insolar/insolar
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//DiscoveryNodesStoreMock implements github.com/insolar/insolar/insolar.DiscoveryNodesStore
type DiscoveryNodesStoreMock struct {
	t minimock.Tester

	StoreDiscoveryNodesFunc       func(p context.Context, p1 []insolar.NetworkNode)
	StoreDiscoveryNodesCounter    uint64
	StoreDiscoveryNodesPreCounter uint64
	StoreDiscoveryNodesMock       mDiscoveryNodesStoreMockStoreDiscoveryNodes
}

//NewDiscoveryNodesStoreMock returns a mock for github.com/insolar/insolar/insolar.DiscoveryNodesStore
func NewDiscoveryNodesStoreMock(t minimock.Tester) *DiscoveryNodesStoreMock {
	m := &DiscoveryNodesStoreMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.StoreDiscoveryNodesMock = mDiscoveryNodesStoreMockStoreDiscoveryNodes{mock: m}

	return m
}

type mDiscoveryNodesStoreMockStoreDiscoveryNodes struct {
	mock              *DiscoveryNodesStoreMock
	mainExpectation   *DiscoveryNodesStoreMockStoreDiscoveryNodesExpectation
	expectationSeries []*DiscoveryNodesStoreMockStoreDiscoveryNodesExpectation
}

type DiscoveryNodesStoreMockStoreDiscoveryNodesExpectation struct {
	input *DiscoveryNodesStoreMockStoreDiscoveryNodesInput
}

type DiscoveryNodesStoreMockStoreDiscoveryNodesInput struct {
	p  context.Context
	p1 []insolar.NetworkNode
}

//Expect specifies that invocation of DiscoveryNodesStore.StoreDiscoveryNodes is expected from 1 to Infinity times
func (m *mDiscoveryNodesStoreMockStoreDiscoveryNodes) Expect(p context.Context, p1 []insolar.NetworkNode) *mDiscoveryNodesStoreMockStoreDiscoveryNodes {
	m.mock.StoreDiscoveryNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodesStoreMockStoreDiscoveryNodesExpectation{}
	}
	m.mainExpectation.input = &DiscoveryNodesStoreMockStoreDiscoveryNodesInput{p, p1}
	return m
}

//Return specifies results of invocation of DiscoveryNodesStore.StoreDiscoveryNodes
func (m *mDiscoveryNodesStoreMockStoreDiscoveryNodes) Return() *DiscoveryNodesStoreMock {
	m.mock.StoreDiscoveryNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodesStoreMockStoreDiscoveryNodesExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of DiscoveryNodesStore.StoreDiscoveryNodes is expected once
func (m *mDiscoveryNodesStoreMockStoreDiscoveryNodes) ExpectOnce(p context.Context, p1 []insolar.NetworkNode) *DiscoveryNodesStoreMockStoreDiscoveryNodesExpectation {
	m.mock.StoreDiscoveryNodesFunc = nil
	m.mainExpectation = nil

	expectation := &DiscoveryNodesStoreMockStoreDiscoveryNodesExpectation{}
	expectation.input = &DiscoveryNodesStoreMockStoreDiscoveryNodesInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of DiscoveryNodesStore.StoreDiscoveryNodes method
func (m *mDiscoveryNodesStoreMockStoreDiscoveryNodes) Set(f func(p context.Context, p1 []insolar.NetworkNode)) *DiscoveryNodesStoreMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StoreDiscoveryNodesFunc = f
	return m.mock
}

//StoreDiscoveryNodes implements github.com/insolar/insolar/insolar.DiscoveryNodesStore interface
func (m *DiscoveryNodesStoreMock) StoreDiscoveryNodes(p context.Context, p1 []insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.StoreDiscoveryNodesPreCounter, 1)
	defer atomic.AddUint64(&m.StoreDiscoveryNodesCounter, 1)

	if len(m.StoreDiscoveryNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StoreDiscoveryNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DiscoveryNodesStoreMock.StoreDiscoveryNodes. %v %v", p, p1)
			return
		}

		input := m.StoreDiscoveryNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DiscoveryNodesStoreMockStoreDiscoveryNodesInput{p, p1}, "DiscoveryNodesStore.StoreDiscoveryNodes got unexpected parameters")

		return
	}

	if m.StoreDiscoveryNodesMock.mainExpectation != nil {

		input := m.StoreDiscoveryNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DiscoveryNodesStoreMockStoreDiscoveryNodesInput{p, p1}, "DiscoveryNodesStore.StoreDiscoveryNodes got unexpected parameters")
		}

		return
	}

	if m.StoreDiscoveryNodesFunc == nil {
		m.t.Fatalf("Unexpected call to DiscoveryNodesStoreMock.StoreDiscoveryNodes. %v %v", p, p1)
		return
	}

	m.StoreDiscoveryNodesFunc(p, p1)
}

//StoreDiscoveryNodesMinimockCounter returns a count of DiscoveryNodesStoreMock.StoreDiscoveryNodesFunc invocations
func (m *DiscoveryNodesStoreMock) StoreDiscoveryNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StoreDiscoveryNodesCounter)
}

//StoreDiscoveryNodesMinimockPreCounter returns the value of DiscoveryNodesStoreMock.StoreDiscoveryNodes invocations
func (m *DiscoveryNodesStoreMock) StoreDiscoveryNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StoreDiscoveryNodesPreCounter)
}

//StoreDiscoveryNodesFinished returns true if mock invocations count is ok
func (m *DiscoveryNodesStoreMock) StoreDiscoveryNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StoreDiscoveryNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StoreDiscoveryNodesCounter) == uint64(len(m.StoreDiscoveryNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StoreDiscoveryNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StoreDiscoveryNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StoreDiscoveryNodesFunc != nil {
		return atomic.LoadUint64(&m.StoreDiscoveryNodesCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DiscoveryNodesStoreMock) ValidateCallCounters() {

	if !m.StoreDiscoveryNodesFinished() {
		m.t.Fatal("Expected call to DiscoveryNodesStoreMock.StoreDiscoveryNodes")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DiscoveryNodesStoreMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DiscoveryNodesStoreMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DiscoveryNodesStoreMock) MinimockFinish() {

	if !m.StoreDiscoveryNodesFinished() {
		m.t.Fatal("Expected call to DiscoveryNodesStoreMock.StoreDiscoveryNodes")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DiscoveryNodesStoreMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DiscoveryNodesStoreMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.StoreDiscoveryNodesFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.StoreDiscoveryNodesFinished() {
				m.t.Error("Expected call to DiscoveryNodesStoreMock.StoreDiscoveryNodes")
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
func (m *DiscoveryNodesStoreMock) AllMocksCalled() bool {

	if !m.StoreDiscoveryNodesFinished() {
		return false
	}

	return true
}
