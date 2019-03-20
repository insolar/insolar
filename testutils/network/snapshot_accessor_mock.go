package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SnapshotAccessor" can be found in github.com/insolar/insolar/network/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	node "github.com/insolar/insolar/network/node"

	testify_assert "github.com/stretchr/testify/assert"
)

//SnapshotAccessorMock implements github.com/insolar/insolar/network/storage.SnapshotAccessor
type SnapshotAccessorMock struct {
	t minimock.Tester

	ForPulseNumberFunc       func(p context.Context, p1 core.PulseNumber) (r *node.Snapshot, r1 error)
	ForPulseNumberCounter    uint64
	ForPulseNumberPreCounter uint64
	ForPulseNumberMock       mSnapshotAccessorMockForPulseNumber

	LatestFunc       func(p context.Context) (r *node.Snapshot, r1 error)
	LatestCounter    uint64
	LatestPreCounter uint64
	LatestMock       mSnapshotAccessorMockLatest
}

//NewSnapshotAccessorMock returns a mock for github.com/insolar/insolar/network/storage.SnapshotAccessor
func NewSnapshotAccessorMock(t minimock.Tester) *SnapshotAccessorMock {
	m := &SnapshotAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseNumberMock = mSnapshotAccessorMockForPulseNumber{mock: m}
	m.LatestMock = mSnapshotAccessorMockLatest{mock: m}

	return m
}

type mSnapshotAccessorMockForPulseNumber struct {
	mock              *SnapshotAccessorMock
	mainExpectation   *SnapshotAccessorMockForPulseNumberExpectation
	expectationSeries []*SnapshotAccessorMockForPulseNumberExpectation
}

type SnapshotAccessorMockForPulseNumberExpectation struct {
	input  *SnapshotAccessorMockForPulseNumberInput
	result *SnapshotAccessorMockForPulseNumberResult
}

type SnapshotAccessorMockForPulseNumberInput struct {
	p  context.Context
	p1 core.PulseNumber
}

type SnapshotAccessorMockForPulseNumberResult struct {
	r  *node.Snapshot
	r1 error
}

//Expect specifies that invocation of SnapshotAccessor.ForPulseNumber is expected from 1 to Infinity times
func (m *mSnapshotAccessorMockForPulseNumber) Expect(p context.Context, p1 core.PulseNumber) *mSnapshotAccessorMockForPulseNumber {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SnapshotAccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.input = &SnapshotAccessorMockForPulseNumberInput{p, p1}
	return m
}

//Return specifies results of invocation of SnapshotAccessor.ForPulseNumber
func (m *mSnapshotAccessorMockForPulseNumber) Return(r *node.Snapshot, r1 error) *SnapshotAccessorMock {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SnapshotAccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.result = &SnapshotAccessorMockForPulseNumberResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of SnapshotAccessor.ForPulseNumber is expected once
func (m *mSnapshotAccessorMockForPulseNumber) ExpectOnce(p context.Context, p1 core.PulseNumber) *SnapshotAccessorMockForPulseNumberExpectation {
	m.mock.ForPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &SnapshotAccessorMockForPulseNumberExpectation{}
	expectation.input = &SnapshotAccessorMockForPulseNumberInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SnapshotAccessorMockForPulseNumberExpectation) Return(r *node.Snapshot, r1 error) {
	e.result = &SnapshotAccessorMockForPulseNumberResult{r, r1}
}

//Set uses given function f as a mock of SnapshotAccessor.ForPulseNumber method
func (m *mSnapshotAccessorMockForPulseNumber) Set(f func(p context.Context, p1 core.PulseNumber) (r *node.Snapshot, r1 error)) *SnapshotAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseNumberFunc = f
	return m.mock
}

//ForPulseNumber implements github.com/insolar/insolar/network/storage.SnapshotAccessor interface
func (m *SnapshotAccessorMock) ForPulseNumber(p context.Context, p1 core.PulseNumber) (r *node.Snapshot, r1 error) {
	counter := atomic.AddUint64(&m.ForPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseNumberCounter, 1)

	if len(m.ForPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SnapshotAccessorMock.ForPulseNumber. %v %v", p, p1)
			return
		}

		input := m.ForPulseNumberMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SnapshotAccessorMockForPulseNumberInput{p, p1}, "SnapshotAccessor.ForPulseNumber got unexpected parameters")

		result := m.ForPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SnapshotAccessorMock.ForPulseNumber")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberMock.mainExpectation != nil {

		input := m.ForPulseNumberMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SnapshotAccessorMockForPulseNumberInput{p, p1}, "SnapshotAccessor.ForPulseNumber got unexpected parameters")
		}

		result := m.ForPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SnapshotAccessorMock.ForPulseNumber")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to SnapshotAccessorMock.ForPulseNumber. %v %v", p, p1)
		return
	}

	return m.ForPulseNumberFunc(p, p1)
}

//ForPulseNumberMinimockCounter returns a count of SnapshotAccessorMock.ForPulseNumberFunc invocations
func (m *SnapshotAccessorMock) ForPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberCounter)
}

//ForPulseNumberMinimockPreCounter returns the value of SnapshotAccessorMock.ForPulseNumber invocations
func (m *SnapshotAccessorMock) ForPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberPreCounter)
}

//ForPulseNumberFinished returns true if mock invocations count is ok
func (m *SnapshotAccessorMock) ForPulseNumberFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseNumberMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) == uint64(len(m.ForPulseNumberMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseNumberMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseNumberFunc != nil {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) > 0
	}

	return true
}

type mSnapshotAccessorMockLatest struct {
	mock              *SnapshotAccessorMock
	mainExpectation   *SnapshotAccessorMockLatestExpectation
	expectationSeries []*SnapshotAccessorMockLatestExpectation
}

type SnapshotAccessorMockLatestExpectation struct {
	input  *SnapshotAccessorMockLatestInput
	result *SnapshotAccessorMockLatestResult
}

type SnapshotAccessorMockLatestInput struct {
	p context.Context
}

type SnapshotAccessorMockLatestResult struct {
	r  *node.Snapshot
	r1 error
}

//Expect specifies that invocation of SnapshotAccessor.Latest is expected from 1 to Infinity times
func (m *mSnapshotAccessorMockLatest) Expect(p context.Context) *mSnapshotAccessorMockLatest {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SnapshotAccessorMockLatestExpectation{}
	}
	m.mainExpectation.input = &SnapshotAccessorMockLatestInput{p}
	return m
}

//Return specifies results of invocation of SnapshotAccessor.Latest
func (m *mSnapshotAccessorMockLatest) Return(r *node.Snapshot, r1 error) *SnapshotAccessorMock {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SnapshotAccessorMockLatestExpectation{}
	}
	m.mainExpectation.result = &SnapshotAccessorMockLatestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of SnapshotAccessor.Latest is expected once
func (m *mSnapshotAccessorMockLatest) ExpectOnce(p context.Context) *SnapshotAccessorMockLatestExpectation {
	m.mock.LatestFunc = nil
	m.mainExpectation = nil

	expectation := &SnapshotAccessorMockLatestExpectation{}
	expectation.input = &SnapshotAccessorMockLatestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SnapshotAccessorMockLatestExpectation) Return(r *node.Snapshot, r1 error) {
	e.result = &SnapshotAccessorMockLatestResult{r, r1}
}

//Set uses given function f as a mock of SnapshotAccessor.Latest method
func (m *mSnapshotAccessorMockLatest) Set(f func(p context.Context) (r *node.Snapshot, r1 error)) *SnapshotAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LatestFunc = f
	return m.mock
}

//Latest implements github.com/insolar/insolar/network/storage.SnapshotAccessor interface
func (m *SnapshotAccessorMock) Latest(p context.Context) (r *node.Snapshot, r1 error) {
	counter := atomic.AddUint64(&m.LatestPreCounter, 1)
	defer atomic.AddUint64(&m.LatestCounter, 1)

	if len(m.LatestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LatestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SnapshotAccessorMock.Latest. %v", p)
			return
		}

		input := m.LatestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SnapshotAccessorMockLatestInput{p}, "SnapshotAccessor.Latest got unexpected parameters")

		result := m.LatestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SnapshotAccessorMock.Latest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestMock.mainExpectation != nil {

		input := m.LatestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SnapshotAccessorMockLatestInput{p}, "SnapshotAccessor.Latest got unexpected parameters")
		}

		result := m.LatestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SnapshotAccessorMock.Latest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestFunc == nil {
		m.t.Fatalf("Unexpected call to SnapshotAccessorMock.Latest. %v", p)
		return
	}

	return m.LatestFunc(p)
}

//LatestMinimockCounter returns a count of SnapshotAccessorMock.LatestFunc invocations
func (m *SnapshotAccessorMock) LatestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LatestCounter)
}

//LatestMinimockPreCounter returns the value of SnapshotAccessorMock.Latest invocations
func (m *SnapshotAccessorMock) LatestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LatestPreCounter)
}

//LatestFinished returns true if mock invocations count is ok
func (m *SnapshotAccessorMock) LatestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LatestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LatestCounter) == uint64(len(m.LatestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LatestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LatestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LatestFunc != nil {
		return atomic.LoadUint64(&m.LatestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SnapshotAccessorMock) ValidateCallCounters() {

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to SnapshotAccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to SnapshotAccessorMock.Latest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SnapshotAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SnapshotAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SnapshotAccessorMock) MinimockFinish() {

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to SnapshotAccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to SnapshotAccessorMock.Latest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SnapshotAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SnapshotAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseNumberFinished()
		ok = ok && m.LatestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseNumberFinished() {
				m.t.Error("Expected call to SnapshotAccessorMock.ForPulseNumber")
			}

			if !m.LatestFinished() {
				m.t.Error("Expected call to SnapshotAccessorMock.Latest")
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
func (m *SnapshotAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseNumberFinished() {
		return false
	}

	if !m.LatestFinished() {
		return false
	}

	return true
}
