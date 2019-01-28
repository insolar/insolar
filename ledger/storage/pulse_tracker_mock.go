package storage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseTracker" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseTrackerMock implements github.com/insolar/insolar/ledger/storage.PulseTracker
type PulseTrackerMock struct {
	t minimock.Tester

	AddPulseFunc       func(p context.Context, p1 core.Pulse) (r error)
	AddPulseCounter    uint64
	AddPulsePreCounter uint64
	AddPulseMock       mPulseTrackerMockAddPulse

	GetLatestPulseFunc       func(p context.Context) (r *Pulse, r1 error)
	GetLatestPulseCounter    uint64
	GetLatestPulsePreCounter uint64
	GetLatestPulseMock       mPulseTrackerMockGetLatestPulse

	GetPreviousPulseFunc       func(p context.Context, p1 core.PulseNumber) (r *Pulse, r1 error)
	GetPreviousPulseCounter    uint64
	GetPreviousPulsePreCounter uint64
	GetPreviousPulseMock       mPulseTrackerMockGetPreviousPulse

	GetPulseFunc       func(p context.Context, p1 core.PulseNumber) (r *Pulse, r1 error)
	GetPulseCounter    uint64
	GetPulsePreCounter uint64
	GetPulseMock       mPulseTrackerMockGetPulse
}

//NewPulseTrackerMock returns a mock for github.com/insolar/insolar/ledger/storage.PulseTracker
func NewPulseTrackerMock(t minimock.Tester) *PulseTrackerMock {
	m := &PulseTrackerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddPulseMock = mPulseTrackerMockAddPulse{mock: m}
	m.GetLatestPulseMock = mPulseTrackerMockGetLatestPulse{mock: m}
	m.GetPreviousPulseMock = mPulseTrackerMockGetPreviousPulse{mock: m}
	m.GetPulseMock = mPulseTrackerMockGetPulse{mock: m}

	return m
}

type mPulseTrackerMockAddPulse struct {
	mock              *PulseTrackerMock
	mainExpectation   *PulseTrackerMockAddPulseExpectation
	expectationSeries []*PulseTrackerMockAddPulseExpectation
}

type PulseTrackerMockAddPulseExpectation struct {
	input  *PulseTrackerMockAddPulseInput
	result *PulseTrackerMockAddPulseResult
}

type PulseTrackerMockAddPulseInput struct {
	p  context.Context
	p1 core.Pulse
}

type PulseTrackerMockAddPulseResult struct {
	r error
}

//Expect specifies that invocation of PulseTracker.AddPulse is expected from 1 to Infinity times
func (m *mPulseTrackerMockAddPulse) Expect(p context.Context, p1 core.Pulse) *mPulseTrackerMockAddPulse {
	m.mock.AddPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockAddPulseExpectation{}
	}
	m.mainExpectation.input = &PulseTrackerMockAddPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseTracker.AddPulse
func (m *mPulseTrackerMockAddPulse) Return(r error) *PulseTrackerMock {
	m.mock.AddPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockAddPulseExpectation{}
	}
	m.mainExpectation.result = &PulseTrackerMockAddPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseTracker.AddPulse is expected once
func (m *mPulseTrackerMockAddPulse) ExpectOnce(p context.Context, p1 core.Pulse) *PulseTrackerMockAddPulseExpectation {
	m.mock.AddPulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulseTrackerMockAddPulseExpectation{}
	expectation.input = &PulseTrackerMockAddPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseTrackerMockAddPulseExpectation) Return(r error) {
	e.result = &PulseTrackerMockAddPulseResult{r}
}

//Set uses given function f as a mock of PulseTracker.AddPulse method
func (m *mPulseTrackerMockAddPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *PulseTrackerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddPulseFunc = f
	return m.mock
}

//AddPulse implements github.com/insolar/insolar/ledger/storage.PulseTracker interface
func (m *PulseTrackerMock) AddPulse(p context.Context, p1 core.Pulse) (r error) {
	counter := atomic.AddUint64(&m.AddPulsePreCounter, 1)
	defer atomic.AddUint64(&m.AddPulseCounter, 1)

	if len(m.AddPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseTrackerMock.AddPulse. %v %v", p, p1)
			return
		}

		input := m.AddPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseTrackerMockAddPulseInput{p, p1}, "PulseTracker.AddPulse got unexpected parameters")

		result := m.AddPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.AddPulse")
			return
		}

		r = result.r

		return
	}

	if m.AddPulseMock.mainExpectation != nil {

		input := m.AddPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseTrackerMockAddPulseInput{p, p1}, "PulseTracker.AddPulse got unexpected parameters")
		}

		result := m.AddPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.AddPulse")
		}

		r = result.r

		return
	}

	if m.AddPulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulseTrackerMock.AddPulse. %v %v", p, p1)
		return
	}

	return m.AddPulseFunc(p, p1)
}

//AddPulseMinimockCounter returns a count of PulseTrackerMock.AddPulseFunc invocations
func (m *PulseTrackerMock) AddPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddPulseCounter)
}

//AddPulseMinimockPreCounter returns the value of PulseTrackerMock.AddPulse invocations
func (m *PulseTrackerMock) AddPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddPulsePreCounter)
}

//AddPulseFinished returns true if mock invocations count is ok
func (m *PulseTrackerMock) AddPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddPulseCounter) == uint64(len(m.AddPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddPulseFunc != nil {
		return atomic.LoadUint64(&m.AddPulseCounter) > 0
	}

	return true
}

type mPulseTrackerMockGetLatestPulse struct {
	mock              *PulseTrackerMock
	mainExpectation   *PulseTrackerMockGetLatestPulseExpectation
	expectationSeries []*PulseTrackerMockGetLatestPulseExpectation
}

type PulseTrackerMockGetLatestPulseExpectation struct {
	input  *PulseTrackerMockGetLatestPulseInput
	result *PulseTrackerMockGetLatestPulseResult
}

type PulseTrackerMockGetLatestPulseInput struct {
	p context.Context
}

type PulseTrackerMockGetLatestPulseResult struct {
	r  *Pulse
	r1 error
}

//Expect specifies that invocation of PulseTracker.GetLatestPulse is expected from 1 to Infinity times
func (m *mPulseTrackerMockGetLatestPulse) Expect(p context.Context) *mPulseTrackerMockGetLatestPulse {
	m.mock.GetLatestPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockGetLatestPulseExpectation{}
	}
	m.mainExpectation.input = &PulseTrackerMockGetLatestPulseInput{p}
	return m
}

//Return specifies results of invocation of PulseTracker.GetLatestPulse
func (m *mPulseTrackerMockGetLatestPulse) Return(r *Pulse, r1 error) *PulseTrackerMock {
	m.mock.GetLatestPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockGetLatestPulseExpectation{}
	}
	m.mainExpectation.result = &PulseTrackerMockGetLatestPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseTracker.GetLatestPulse is expected once
func (m *mPulseTrackerMockGetLatestPulse) ExpectOnce(p context.Context) *PulseTrackerMockGetLatestPulseExpectation {
	m.mock.GetLatestPulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulseTrackerMockGetLatestPulseExpectation{}
	expectation.input = &PulseTrackerMockGetLatestPulseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseTrackerMockGetLatestPulseExpectation) Return(r *Pulse, r1 error) {
	e.result = &PulseTrackerMockGetLatestPulseResult{r, r1}
}

//Set uses given function f as a mock of PulseTracker.GetLatestPulse method
func (m *mPulseTrackerMockGetLatestPulse) Set(f func(p context.Context) (r *Pulse, r1 error)) *PulseTrackerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetLatestPulseFunc = f
	return m.mock
}

//GetLatestPulse implements github.com/insolar/insolar/ledger/storage.PulseTracker interface
func (m *PulseTrackerMock) GetLatestPulse(p context.Context) (r *Pulse, r1 error) {
	counter := atomic.AddUint64(&m.GetLatestPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetLatestPulseCounter, 1)

	if len(m.GetLatestPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetLatestPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseTrackerMock.GetLatestPulse. %v", p)
			return
		}

		input := m.GetLatestPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseTrackerMockGetLatestPulseInput{p}, "PulseTracker.GetLatestPulse got unexpected parameters")

		result := m.GetLatestPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.GetLatestPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetLatestPulseMock.mainExpectation != nil {

		input := m.GetLatestPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseTrackerMockGetLatestPulseInput{p}, "PulseTracker.GetLatestPulse got unexpected parameters")
		}

		result := m.GetLatestPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.GetLatestPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetLatestPulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulseTrackerMock.GetLatestPulse. %v", p)
		return
	}

	return m.GetLatestPulseFunc(p)
}

//GetLatestPulseMinimockCounter returns a count of PulseTrackerMock.GetLatestPulseFunc invocations
func (m *PulseTrackerMock) GetLatestPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetLatestPulseCounter)
}

//GetLatestPulseMinimockPreCounter returns the value of PulseTrackerMock.GetLatestPulse invocations
func (m *PulseTrackerMock) GetLatestPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetLatestPulsePreCounter)
}

//GetLatestPulseFinished returns true if mock invocations count is ok
func (m *PulseTrackerMock) GetLatestPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetLatestPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetLatestPulseCounter) == uint64(len(m.GetLatestPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetLatestPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetLatestPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetLatestPulseFunc != nil {
		return atomic.LoadUint64(&m.GetLatestPulseCounter) > 0
	}

	return true
}

type mPulseTrackerMockGetPreviousPulse struct {
	mock              *PulseTrackerMock
	mainExpectation   *PulseTrackerMockGetPreviousPulseExpectation
	expectationSeries []*PulseTrackerMockGetPreviousPulseExpectation
}

type PulseTrackerMockGetPreviousPulseExpectation struct {
	input  *PulseTrackerMockGetPreviousPulseInput
	result *PulseTrackerMockGetPreviousPulseResult
}

type PulseTrackerMockGetPreviousPulseInput struct {
	p  context.Context
	p1 core.PulseNumber
}

type PulseTrackerMockGetPreviousPulseResult struct {
	r  *Pulse
	r1 error
}

//Expect specifies that invocation of PulseTracker.GetPreviousPulse is expected from 1 to Infinity times
func (m *mPulseTrackerMockGetPreviousPulse) Expect(p context.Context, p1 core.PulseNumber) *mPulseTrackerMockGetPreviousPulse {
	m.mock.GetPreviousPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockGetPreviousPulseExpectation{}
	}
	m.mainExpectation.input = &PulseTrackerMockGetPreviousPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseTracker.GetPreviousPulse
func (m *mPulseTrackerMockGetPreviousPulse) Return(r *Pulse, r1 error) *PulseTrackerMock {
	m.mock.GetPreviousPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockGetPreviousPulseExpectation{}
	}
	m.mainExpectation.result = &PulseTrackerMockGetPreviousPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseTracker.GetPreviousPulse is expected once
func (m *mPulseTrackerMockGetPreviousPulse) ExpectOnce(p context.Context, p1 core.PulseNumber) *PulseTrackerMockGetPreviousPulseExpectation {
	m.mock.GetPreviousPulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulseTrackerMockGetPreviousPulseExpectation{}
	expectation.input = &PulseTrackerMockGetPreviousPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseTrackerMockGetPreviousPulseExpectation) Return(r *Pulse, r1 error) {
	e.result = &PulseTrackerMockGetPreviousPulseResult{r, r1}
}

//Set uses given function f as a mock of PulseTracker.GetPreviousPulse method
func (m *mPulseTrackerMockGetPreviousPulse) Set(f func(p context.Context, p1 core.PulseNumber) (r *Pulse, r1 error)) *PulseTrackerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPreviousPulseFunc = f
	return m.mock
}

//GetPreviousPulse implements github.com/insolar/insolar/ledger/storage.PulseTracker interface
func (m *PulseTrackerMock) GetPreviousPulse(p context.Context, p1 core.PulseNumber) (r *Pulse, r1 error) {
	counter := atomic.AddUint64(&m.GetPreviousPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetPreviousPulseCounter, 1)

	if len(m.GetPreviousPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPreviousPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseTrackerMock.GetPreviousPulse. %v %v", p, p1)
			return
		}

		input := m.GetPreviousPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseTrackerMockGetPreviousPulseInput{p, p1}, "PulseTracker.GetPreviousPulse got unexpected parameters")

		result := m.GetPreviousPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.GetPreviousPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPreviousPulseMock.mainExpectation != nil {

		input := m.GetPreviousPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseTrackerMockGetPreviousPulseInput{p, p1}, "PulseTracker.GetPreviousPulse got unexpected parameters")
		}

		result := m.GetPreviousPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.GetPreviousPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPreviousPulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulseTrackerMock.GetPreviousPulse. %v %v", p, p1)
		return
	}

	return m.GetPreviousPulseFunc(p, p1)
}

//GetPreviousPulseMinimockCounter returns a count of PulseTrackerMock.GetPreviousPulseFunc invocations
func (m *PulseTrackerMock) GetPreviousPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreviousPulseCounter)
}

//GetPreviousPulseMinimockPreCounter returns the value of PulseTrackerMock.GetPreviousPulse invocations
func (m *PulseTrackerMock) GetPreviousPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPreviousPulsePreCounter)
}

//GetPreviousPulseFinished returns true if mock invocations count is ok
func (m *PulseTrackerMock) GetPreviousPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPreviousPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPreviousPulseCounter) == uint64(len(m.GetPreviousPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPreviousPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPreviousPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPreviousPulseFunc != nil {
		return atomic.LoadUint64(&m.GetPreviousPulseCounter) > 0
	}

	return true
}

type mPulseTrackerMockGetPulse struct {
	mock              *PulseTrackerMock
	mainExpectation   *PulseTrackerMockGetPulseExpectation
	expectationSeries []*PulseTrackerMockGetPulseExpectation
}

type PulseTrackerMockGetPulseExpectation struct {
	input  *PulseTrackerMockGetPulseInput
	result *PulseTrackerMockGetPulseResult
}

type PulseTrackerMockGetPulseInput struct {
	p  context.Context
	p1 core.PulseNumber
}

type PulseTrackerMockGetPulseResult struct {
	r  *Pulse
	r1 error
}

//Expect specifies that invocation of PulseTracker.GetPulse is expected from 1 to Infinity times
func (m *mPulseTrackerMockGetPulse) Expect(p context.Context, p1 core.PulseNumber) *mPulseTrackerMockGetPulse {
	m.mock.GetPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockGetPulseExpectation{}
	}
	m.mainExpectation.input = &PulseTrackerMockGetPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseTracker.GetPulse
func (m *mPulseTrackerMockGetPulse) Return(r *Pulse, r1 error) *PulseTrackerMock {
	m.mock.GetPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseTrackerMockGetPulseExpectation{}
	}
	m.mainExpectation.result = &PulseTrackerMockGetPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseTracker.GetPulse is expected once
func (m *mPulseTrackerMockGetPulse) ExpectOnce(p context.Context, p1 core.PulseNumber) *PulseTrackerMockGetPulseExpectation {
	m.mock.GetPulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulseTrackerMockGetPulseExpectation{}
	expectation.input = &PulseTrackerMockGetPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseTrackerMockGetPulseExpectation) Return(r *Pulse, r1 error) {
	e.result = &PulseTrackerMockGetPulseResult{r, r1}
}

//Set uses given function f as a mock of PulseTracker.GetPulse method
func (m *mPulseTrackerMockGetPulse) Set(f func(p context.Context, p1 core.PulseNumber) (r *Pulse, r1 error)) *PulseTrackerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseFunc = f
	return m.mock
}

//GetPulse implements github.com/insolar/insolar/ledger/storage.PulseTracker interface
func (m *PulseTrackerMock) GetPulse(p context.Context, p1 core.PulseNumber) (r *Pulse, r1 error) {
	counter := atomic.AddUint64(&m.GetPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseCounter, 1)

	if len(m.GetPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseTrackerMock.GetPulse. %v %v", p, p1)
			return
		}

		input := m.GetPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseTrackerMockGetPulseInput{p, p1}, "PulseTracker.GetPulse got unexpected parameters")

		result := m.GetPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.GetPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPulseMock.mainExpectation != nil {

		input := m.GetPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseTrackerMockGetPulseInput{p, p1}, "PulseTracker.GetPulse got unexpected parameters")
		}

		result := m.GetPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseTrackerMock.GetPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulseTrackerMock.GetPulse. %v %v", p, p1)
		return
	}

	return m.GetPulseFunc(p, p1)
}

//GetPulseMinimockCounter returns a count of PulseTrackerMock.GetPulseFunc invocations
func (m *PulseTrackerMock) GetPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseCounter)
}

//GetPulseMinimockPreCounter returns the value of PulseTrackerMock.GetPulse invocations
func (m *PulseTrackerMock) GetPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulsePreCounter)
}

//GetPulseFinished returns true if mock invocations count is ok
func (m *PulseTrackerMock) GetPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulseCounter) == uint64(len(m.GetPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulseFunc != nil {
		return atomic.LoadUint64(&m.GetPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseTrackerMock) ValidateCallCounters() {

	if !m.AddPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.AddPulse")
	}

	if !m.GetLatestPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.GetLatestPulse")
	}

	if !m.GetPreviousPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.GetPreviousPulse")
	}

	if !m.GetPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.GetPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseTrackerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseTrackerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseTrackerMock) MinimockFinish() {

	if !m.AddPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.AddPulse")
	}

	if !m.GetLatestPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.GetLatestPulse")
	}

	if !m.GetPreviousPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.GetPreviousPulse")
	}

	if !m.GetPulseFinished() {
		m.t.Fatal("Expected call to PulseTrackerMock.GetPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseTrackerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseTrackerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddPulseFinished()
		ok = ok && m.GetLatestPulseFinished()
		ok = ok && m.GetPreviousPulseFinished()
		ok = ok && m.GetPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddPulseFinished() {
				m.t.Error("Expected call to PulseTrackerMock.AddPulse")
			}

			if !m.GetLatestPulseFinished() {
				m.t.Error("Expected call to PulseTrackerMock.GetLatestPulse")
			}

			if !m.GetPreviousPulseFinished() {
				m.t.Error("Expected call to PulseTrackerMock.GetPreviousPulse")
			}

			if !m.GetPulseFinished() {
				m.t.Error("Expected call to PulseTrackerMock.GetPulse")
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
func (m *PulseTrackerMock) AllMocksCalled() bool {

	if !m.AddPulseFinished() {
		return false
	}

	if !m.GetLatestPulseFinished() {
		return false
	}

	if !m.GetPreviousPulseFinished() {
		return false
	}

	if !m.GetPulseFinished() {
		return false
	}

	return true
}
