package introspector

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PublisherServer" can be found in github.com/insolar/insolar/instrumentation/introspector/introproto
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	introproto "github.com/insolar/insolar/instrumentation/introspector/introproto"

	testify_assert "github.com/stretchr/testify/assert"
)

//PublisherServerMock implements github.com/insolar/insolar/instrumentation/introspector/introproto.PublisherServer
type PublisherServerMock struct {
	t minimock.Tester

	GetMessagesFiltersFunc       func(p context.Context, p1 *introproto.EmptyArgs) (r *introproto.AllMessageFilterStats, r1 error)
	GetMessagesFiltersCounter    uint64
	GetMessagesFiltersPreCounter uint64
	GetMessagesFiltersMock       mPublisherServerMockGetMessagesFilters

	GetMessagesStatFunc       func(p context.Context, p1 *introproto.EmptyArgs) (r *introproto.AllMessageStatByType, r1 error)
	GetMessagesStatCounter    uint64
	GetMessagesStatPreCounter uint64
	GetMessagesStatMock       mPublisherServerMockGetMessagesStat

	SetMessagesFilterFunc       func(p context.Context, p1 *introproto.MessageFilterByType) (r *introproto.MessageFilterByType, r1 error)
	SetMessagesFilterCounter    uint64
	SetMessagesFilterPreCounter uint64
	SetMessagesFilterMock       mPublisherServerMockSetMessagesFilter
}

//NewPublisherServerMock returns a mock for github.com/insolar/insolar/instrumentation/introspector/introproto.PublisherServer
func NewPublisherServerMock(t minimock.Tester) *PublisherServerMock {
	m := &PublisherServerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetMessagesFiltersMock = mPublisherServerMockGetMessagesFilters{mock: m}
	m.GetMessagesStatMock = mPublisherServerMockGetMessagesStat{mock: m}
	m.SetMessagesFilterMock = mPublisherServerMockSetMessagesFilter{mock: m}

	return m
}

type mPublisherServerMockGetMessagesFilters struct {
	mock              *PublisherServerMock
	mainExpectation   *PublisherServerMockGetMessagesFiltersExpectation
	expectationSeries []*PublisherServerMockGetMessagesFiltersExpectation
}

type PublisherServerMockGetMessagesFiltersExpectation struct {
	input  *PublisherServerMockGetMessagesFiltersInput
	result *PublisherServerMockGetMessagesFiltersResult
}

type PublisherServerMockGetMessagesFiltersInput struct {
	p  context.Context
	p1 *introproto.EmptyArgs
}

type PublisherServerMockGetMessagesFiltersResult struct {
	r  *introproto.AllMessageFilterStats
	r1 error
}

//Expect specifies that invocation of PublisherServer.GetMessagesFilters is expected from 1 to Infinity times
func (m *mPublisherServerMockGetMessagesFilters) Expect(p context.Context, p1 *introproto.EmptyArgs) *mPublisherServerMockGetMessagesFilters {
	m.mock.GetMessagesFiltersFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublisherServerMockGetMessagesFiltersExpectation{}
	}
	m.mainExpectation.input = &PublisherServerMockGetMessagesFiltersInput{p, p1}
	return m
}

//Return specifies results of invocation of PublisherServer.GetMessagesFilters
func (m *mPublisherServerMockGetMessagesFilters) Return(r *introproto.AllMessageFilterStats, r1 error) *PublisherServerMock {
	m.mock.GetMessagesFiltersFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublisherServerMockGetMessagesFiltersExpectation{}
	}
	m.mainExpectation.result = &PublisherServerMockGetMessagesFiltersResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PublisherServer.GetMessagesFilters is expected once
func (m *mPublisherServerMockGetMessagesFilters) ExpectOnce(p context.Context, p1 *introproto.EmptyArgs) *PublisherServerMockGetMessagesFiltersExpectation {
	m.mock.GetMessagesFiltersFunc = nil
	m.mainExpectation = nil

	expectation := &PublisherServerMockGetMessagesFiltersExpectation{}
	expectation.input = &PublisherServerMockGetMessagesFiltersInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PublisherServerMockGetMessagesFiltersExpectation) Return(r *introproto.AllMessageFilterStats, r1 error) {
	e.result = &PublisherServerMockGetMessagesFiltersResult{r, r1}
}

//Set uses given function f as a mock of PublisherServer.GetMessagesFilters method
func (m *mPublisherServerMockGetMessagesFilters) Set(f func(p context.Context, p1 *introproto.EmptyArgs) (r *introproto.AllMessageFilterStats, r1 error)) *PublisherServerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMessagesFiltersFunc = f
	return m.mock
}

//GetMessagesFilters implements github.com/insolar/insolar/instrumentation/introspector/introproto.PublisherServer interface
func (m *PublisherServerMock) GetMessagesFilters(p context.Context, p1 *introproto.EmptyArgs) (r *introproto.AllMessageFilterStats, r1 error) {
	counter := atomic.AddUint64(&m.GetMessagesFiltersPreCounter, 1)
	defer atomic.AddUint64(&m.GetMessagesFiltersCounter, 1)

	if len(m.GetMessagesFiltersMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMessagesFiltersMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PublisherServerMock.GetMessagesFilters. %v %v", p, p1)
			return
		}

		input := m.GetMessagesFiltersMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PublisherServerMockGetMessagesFiltersInput{p, p1}, "PublisherServer.GetMessagesFilters got unexpected parameters")

		result := m.GetMessagesFiltersMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PublisherServerMock.GetMessagesFilters")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMessagesFiltersMock.mainExpectation != nil {

		input := m.GetMessagesFiltersMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PublisherServerMockGetMessagesFiltersInput{p, p1}, "PublisherServer.GetMessagesFilters got unexpected parameters")
		}

		result := m.GetMessagesFiltersMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PublisherServerMock.GetMessagesFilters")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMessagesFiltersFunc == nil {
		m.t.Fatalf("Unexpected call to PublisherServerMock.GetMessagesFilters. %v %v", p, p1)
		return
	}

	return m.GetMessagesFiltersFunc(p, p1)
}

//GetMessagesFiltersMinimockCounter returns a count of PublisherServerMock.GetMessagesFiltersFunc invocations
func (m *PublisherServerMock) GetMessagesFiltersMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMessagesFiltersCounter)
}

//GetMessagesFiltersMinimockPreCounter returns the value of PublisherServerMock.GetMessagesFilters invocations
func (m *PublisherServerMock) GetMessagesFiltersMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMessagesFiltersPreCounter)
}

//GetMessagesFiltersFinished returns true if mock invocations count is ok
func (m *PublisherServerMock) GetMessagesFiltersFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMessagesFiltersMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetMessagesFiltersCounter) == uint64(len(m.GetMessagesFiltersMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMessagesFiltersMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetMessagesFiltersCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetMessagesFiltersFunc != nil {
		return atomic.LoadUint64(&m.GetMessagesFiltersCounter) > 0
	}

	return true
}

type mPublisherServerMockGetMessagesStat struct {
	mock              *PublisherServerMock
	mainExpectation   *PublisherServerMockGetMessagesStatExpectation
	expectationSeries []*PublisherServerMockGetMessagesStatExpectation
}

type PublisherServerMockGetMessagesStatExpectation struct {
	input  *PublisherServerMockGetMessagesStatInput
	result *PublisherServerMockGetMessagesStatResult
}

type PublisherServerMockGetMessagesStatInput struct {
	p  context.Context
	p1 *introproto.EmptyArgs
}

type PublisherServerMockGetMessagesStatResult struct {
	r  *introproto.AllMessageStatByType
	r1 error
}

//Expect specifies that invocation of PublisherServer.GetMessagesStat is expected from 1 to Infinity times
func (m *mPublisherServerMockGetMessagesStat) Expect(p context.Context, p1 *introproto.EmptyArgs) *mPublisherServerMockGetMessagesStat {
	m.mock.GetMessagesStatFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublisherServerMockGetMessagesStatExpectation{}
	}
	m.mainExpectation.input = &PublisherServerMockGetMessagesStatInput{p, p1}
	return m
}

//Return specifies results of invocation of PublisherServer.GetMessagesStat
func (m *mPublisherServerMockGetMessagesStat) Return(r *introproto.AllMessageStatByType, r1 error) *PublisherServerMock {
	m.mock.GetMessagesStatFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublisherServerMockGetMessagesStatExpectation{}
	}
	m.mainExpectation.result = &PublisherServerMockGetMessagesStatResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PublisherServer.GetMessagesStat is expected once
func (m *mPublisherServerMockGetMessagesStat) ExpectOnce(p context.Context, p1 *introproto.EmptyArgs) *PublisherServerMockGetMessagesStatExpectation {
	m.mock.GetMessagesStatFunc = nil
	m.mainExpectation = nil

	expectation := &PublisherServerMockGetMessagesStatExpectation{}
	expectation.input = &PublisherServerMockGetMessagesStatInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PublisherServerMockGetMessagesStatExpectation) Return(r *introproto.AllMessageStatByType, r1 error) {
	e.result = &PublisherServerMockGetMessagesStatResult{r, r1}
}

//Set uses given function f as a mock of PublisherServer.GetMessagesStat method
func (m *mPublisherServerMockGetMessagesStat) Set(f func(p context.Context, p1 *introproto.EmptyArgs) (r *introproto.AllMessageStatByType, r1 error)) *PublisherServerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMessagesStatFunc = f
	return m.mock
}

//GetMessagesStat implements github.com/insolar/insolar/instrumentation/introspector/introproto.PublisherServer interface
func (m *PublisherServerMock) GetMessagesStat(p context.Context, p1 *introproto.EmptyArgs) (r *introproto.AllMessageStatByType, r1 error) {
	counter := atomic.AddUint64(&m.GetMessagesStatPreCounter, 1)
	defer atomic.AddUint64(&m.GetMessagesStatCounter, 1)

	if len(m.GetMessagesStatMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMessagesStatMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PublisherServerMock.GetMessagesStat. %v %v", p, p1)
			return
		}

		input := m.GetMessagesStatMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PublisherServerMockGetMessagesStatInput{p, p1}, "PublisherServer.GetMessagesStat got unexpected parameters")

		result := m.GetMessagesStatMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PublisherServerMock.GetMessagesStat")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMessagesStatMock.mainExpectation != nil {

		input := m.GetMessagesStatMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PublisherServerMockGetMessagesStatInput{p, p1}, "PublisherServer.GetMessagesStat got unexpected parameters")
		}

		result := m.GetMessagesStatMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PublisherServerMock.GetMessagesStat")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetMessagesStatFunc == nil {
		m.t.Fatalf("Unexpected call to PublisherServerMock.GetMessagesStat. %v %v", p, p1)
		return
	}

	return m.GetMessagesStatFunc(p, p1)
}

//GetMessagesStatMinimockCounter returns a count of PublisherServerMock.GetMessagesStatFunc invocations
func (m *PublisherServerMock) GetMessagesStatMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMessagesStatCounter)
}

//GetMessagesStatMinimockPreCounter returns the value of PublisherServerMock.GetMessagesStat invocations
func (m *PublisherServerMock) GetMessagesStatMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMessagesStatPreCounter)
}

//GetMessagesStatFinished returns true if mock invocations count is ok
func (m *PublisherServerMock) GetMessagesStatFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMessagesStatMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetMessagesStatCounter) == uint64(len(m.GetMessagesStatMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMessagesStatMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetMessagesStatCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetMessagesStatFunc != nil {
		return atomic.LoadUint64(&m.GetMessagesStatCounter) > 0
	}

	return true
}

type mPublisherServerMockSetMessagesFilter struct {
	mock              *PublisherServerMock
	mainExpectation   *PublisherServerMockSetMessagesFilterExpectation
	expectationSeries []*PublisherServerMockSetMessagesFilterExpectation
}

type PublisherServerMockSetMessagesFilterExpectation struct {
	input  *PublisherServerMockSetMessagesFilterInput
	result *PublisherServerMockSetMessagesFilterResult
}

type PublisherServerMockSetMessagesFilterInput struct {
	p  context.Context
	p1 *introproto.MessageFilterByType
}

type PublisherServerMockSetMessagesFilterResult struct {
	r  *introproto.MessageFilterByType
	r1 error
}

//Expect specifies that invocation of PublisherServer.SetMessagesFilter is expected from 1 to Infinity times
func (m *mPublisherServerMockSetMessagesFilter) Expect(p context.Context, p1 *introproto.MessageFilterByType) *mPublisherServerMockSetMessagesFilter {
	m.mock.SetMessagesFilterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublisherServerMockSetMessagesFilterExpectation{}
	}
	m.mainExpectation.input = &PublisherServerMockSetMessagesFilterInput{p, p1}
	return m
}

//Return specifies results of invocation of PublisherServer.SetMessagesFilter
func (m *mPublisherServerMockSetMessagesFilter) Return(r *introproto.MessageFilterByType, r1 error) *PublisherServerMock {
	m.mock.SetMessagesFilterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PublisherServerMockSetMessagesFilterExpectation{}
	}
	m.mainExpectation.result = &PublisherServerMockSetMessagesFilterResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PublisherServer.SetMessagesFilter is expected once
func (m *mPublisherServerMockSetMessagesFilter) ExpectOnce(p context.Context, p1 *introproto.MessageFilterByType) *PublisherServerMockSetMessagesFilterExpectation {
	m.mock.SetMessagesFilterFunc = nil
	m.mainExpectation = nil

	expectation := &PublisherServerMockSetMessagesFilterExpectation{}
	expectation.input = &PublisherServerMockSetMessagesFilterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PublisherServerMockSetMessagesFilterExpectation) Return(r *introproto.MessageFilterByType, r1 error) {
	e.result = &PublisherServerMockSetMessagesFilterResult{r, r1}
}

//Set uses given function f as a mock of PublisherServer.SetMessagesFilter method
func (m *mPublisherServerMockSetMessagesFilter) Set(f func(p context.Context, p1 *introproto.MessageFilterByType) (r *introproto.MessageFilterByType, r1 error)) *PublisherServerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetMessagesFilterFunc = f
	return m.mock
}

//SetMessagesFilter implements github.com/insolar/insolar/instrumentation/introspector/introproto.PublisherServer interface
func (m *PublisherServerMock) SetMessagesFilter(p context.Context, p1 *introproto.MessageFilterByType) (r *introproto.MessageFilterByType, r1 error) {
	counter := atomic.AddUint64(&m.SetMessagesFilterPreCounter, 1)
	defer atomic.AddUint64(&m.SetMessagesFilterCounter, 1)

	if len(m.SetMessagesFilterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMessagesFilterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PublisherServerMock.SetMessagesFilter. %v %v", p, p1)
			return
		}

		input := m.SetMessagesFilterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PublisherServerMockSetMessagesFilterInput{p, p1}, "PublisherServer.SetMessagesFilter got unexpected parameters")

		result := m.SetMessagesFilterMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PublisherServerMock.SetMessagesFilter")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetMessagesFilterMock.mainExpectation != nil {

		input := m.SetMessagesFilterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PublisherServerMockSetMessagesFilterInput{p, p1}, "PublisherServer.SetMessagesFilter got unexpected parameters")
		}

		result := m.SetMessagesFilterMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PublisherServerMock.SetMessagesFilter")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SetMessagesFilterFunc == nil {
		m.t.Fatalf("Unexpected call to PublisherServerMock.SetMessagesFilter. %v %v", p, p1)
		return
	}

	return m.SetMessagesFilterFunc(p, p1)
}

//SetMessagesFilterMinimockCounter returns a count of PublisherServerMock.SetMessagesFilterFunc invocations
func (m *PublisherServerMock) SetMessagesFilterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetMessagesFilterCounter)
}

//SetMessagesFilterMinimockPreCounter returns the value of PublisherServerMock.SetMessagesFilter invocations
func (m *PublisherServerMock) SetMessagesFilterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetMessagesFilterPreCounter)
}

//SetMessagesFilterFinished returns true if mock invocations count is ok
func (m *PublisherServerMock) SetMessagesFilterFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMessagesFilterMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetMessagesFilterCounter) == uint64(len(m.SetMessagesFilterMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMessagesFilterMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetMessagesFilterCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetMessagesFilterFunc != nil {
		return atomic.LoadUint64(&m.SetMessagesFilterCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PublisherServerMock) ValidateCallCounters() {

	if !m.GetMessagesFiltersFinished() {
		m.t.Fatal("Expected call to PublisherServerMock.GetMessagesFilters")
	}

	if !m.GetMessagesStatFinished() {
		m.t.Fatal("Expected call to PublisherServerMock.GetMessagesStat")
	}

	if !m.SetMessagesFilterFinished() {
		m.t.Fatal("Expected call to PublisherServerMock.SetMessagesFilter")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PublisherServerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PublisherServerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PublisherServerMock) MinimockFinish() {

	if !m.GetMessagesFiltersFinished() {
		m.t.Fatal("Expected call to PublisherServerMock.GetMessagesFilters")
	}

	if !m.GetMessagesStatFinished() {
		m.t.Fatal("Expected call to PublisherServerMock.GetMessagesStat")
	}

	if !m.SetMessagesFilterFinished() {
		m.t.Fatal("Expected call to PublisherServerMock.SetMessagesFilter")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PublisherServerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PublisherServerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetMessagesFiltersFinished()
		ok = ok && m.GetMessagesStatFinished()
		ok = ok && m.SetMessagesFilterFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetMessagesFiltersFinished() {
				m.t.Error("Expected call to PublisherServerMock.GetMessagesFilters")
			}

			if !m.GetMessagesStatFinished() {
				m.t.Error("Expected call to PublisherServerMock.GetMessagesStat")
			}

			if !m.SetMessagesFilterFinished() {
				m.t.Error("Expected call to PublisherServerMock.SetMessagesFilter")
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
func (m *PublisherServerMock) AllMocksCalled() bool {

	if !m.GetMessagesFiltersFinished() {
		return false
	}

	if !m.GetMessagesStatFinished() {
		return false
	}

	if !m.SetMessagesFilterFinished() {
		return false
	}

	return true
}
