package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Conveyor" can be found in github.com/insolar/insolar/insolar
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	queue "github.com/insolar/insolar/conveyor/queue"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ConveyorMock implements github.com/insolar/insolar/insolar.Conveyor
type ConveyorMock struct {
	t minimock.Tester

	ActivatePulseFunc       func() (r error)
	ActivatePulseCounter    uint64
	ActivatePulsePreCounter uint64
	ActivatePulseMock       mConveyorMockActivatePulse

	GetStateFunc       func() (r insolar.ConveyorState)
	GetStateCounter    uint64
	GetStatePreCounter uint64
	GetStateMock       mConveyorMockGetState

	InitiateShutdownFunc       func(p bool)
	InitiateShutdownCounter    uint64
	InitiateShutdownPreCounter uint64
	InitiateShutdownMock       mConveyorMockInitiateShutdown

	IsOperationalFunc       func() (r bool)
	IsOperationalCounter    uint64
	IsOperationalPreCounter uint64
	IsOperationalMock       mConveyorMockIsOperational

	PreparePulseFunc       func(p insolar.Pulse, p1 queue.SyncDone) (r error)
	PreparePulseCounter    uint64
	PreparePulsePreCounter uint64
	PreparePulseMock       mConveyorMockPreparePulse

	SinkPushFunc       func(p insolar.PulseNumber, p1 interface{}) (r error)
	SinkPushCounter    uint64
	SinkPushPreCounter uint64
	SinkPushMock       mConveyorMockSinkPush

	SinkPushAllFunc       func(p insolar.PulseNumber, p1 []interface{}) (r error)
	SinkPushAllCounter    uint64
	SinkPushAllPreCounter uint64
	SinkPushAllMock       mConveyorMockSinkPushAll
}

//NewConveyorMock returns a mock for github.com/insolar/insolar/insolar.Conveyor
func NewConveyorMock(t minimock.Tester) *ConveyorMock {
	m := &ConveyorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ActivatePulseMock = mConveyorMockActivatePulse{mock: m}
	m.GetStateMock = mConveyorMockGetState{mock: m}
	m.InitiateShutdownMock = mConveyorMockInitiateShutdown{mock: m}
	m.IsOperationalMock = mConveyorMockIsOperational{mock: m}
	m.PreparePulseMock = mConveyorMockPreparePulse{mock: m}
	m.SinkPushMock = mConveyorMockSinkPush{mock: m}
	m.SinkPushAllMock = mConveyorMockSinkPushAll{mock: m}

	return m
}

type mConveyorMockActivatePulse struct {
	mock              *ConveyorMock
	mainExpectation   *ConveyorMockActivatePulseExpectation
	expectationSeries []*ConveyorMockActivatePulseExpectation
}

type ConveyorMockActivatePulseExpectation struct {
	result *ConveyorMockActivatePulseResult
}

type ConveyorMockActivatePulseResult struct {
	r error
}

//Expect specifies that invocation of Conveyor.ActivatePulse is expected from 1 to Infinity times
func (m *mConveyorMockActivatePulse) Expect() *mConveyorMockActivatePulse {
	m.mock.ActivatePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockActivatePulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of Conveyor.ActivatePulse
func (m *mConveyorMockActivatePulse) Return(r error) *ConveyorMock {
	m.mock.ActivatePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockActivatePulseExpectation{}
	}
	m.mainExpectation.result = &ConveyorMockActivatePulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Conveyor.ActivatePulse is expected once
func (m *mConveyorMockActivatePulse) ExpectOnce() *ConveyorMockActivatePulseExpectation {
	m.mock.ActivatePulseFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorMockActivatePulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorMockActivatePulseExpectation) Return(r error) {
	e.result = &ConveyorMockActivatePulseResult{r}
}

//Set uses given function f as a mock of Conveyor.ActivatePulse method
func (m *mConveyorMockActivatePulse) Set(f func() (r error)) *ConveyorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ActivatePulseFunc = f
	return m.mock
}

//ActivatePulse implements github.com/insolar/insolar/insolar.Conveyor interface
func (m *ConveyorMock) ActivatePulse() (r error) {
	counter := atomic.AddUint64(&m.ActivatePulsePreCounter, 1)
	defer atomic.AddUint64(&m.ActivatePulseCounter, 1)

	if len(m.ActivatePulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ActivatePulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorMock.ActivatePulse.")
			return
		}

		result := m.ActivatePulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.ActivatePulse")
			return
		}

		r = result.r

		return
	}

	if m.ActivatePulseMock.mainExpectation != nil {

		result := m.ActivatePulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.ActivatePulse")
		}

		r = result.r

		return
	}

	if m.ActivatePulseFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorMock.ActivatePulse.")
		return
	}

	return m.ActivatePulseFunc()
}

//ActivatePulseMinimockCounter returns a count of ConveyorMock.ActivatePulseFunc invocations
func (m *ConveyorMock) ActivatePulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ActivatePulseCounter)
}

//ActivatePulseMinimockPreCounter returns the value of ConveyorMock.ActivatePulse invocations
func (m *ConveyorMock) ActivatePulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ActivatePulsePreCounter)
}

//ActivatePulseFinished returns true if mock invocations count is ok
func (m *ConveyorMock) ActivatePulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ActivatePulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ActivatePulseCounter) == uint64(len(m.ActivatePulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ActivatePulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ActivatePulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ActivatePulseFunc != nil {
		return atomic.LoadUint64(&m.ActivatePulseCounter) > 0
	}

	return true
}

type mConveyorMockGetState struct {
	mock              *ConveyorMock
	mainExpectation   *ConveyorMockGetStateExpectation
	expectationSeries []*ConveyorMockGetStateExpectation
}

type ConveyorMockGetStateExpectation struct {
	result *ConveyorMockGetStateResult
}

type ConveyorMockGetStateResult struct {
	r insolar.ConveyorState
}

//Expect specifies that invocation of Conveyor.GetState is expected from 1 to Infinity times
func (m *mConveyorMockGetState) Expect() *mConveyorMockGetState {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockGetStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of Conveyor.GetState
func (m *mConveyorMockGetState) Return(r insolar.ConveyorState) *ConveyorMock {
	m.mock.GetStateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockGetStateExpectation{}
	}
	m.mainExpectation.result = &ConveyorMockGetStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Conveyor.GetState is expected once
func (m *mConveyorMockGetState) ExpectOnce() *ConveyorMockGetStateExpectation {
	m.mock.GetStateFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorMockGetStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorMockGetStateExpectation) Return(r insolar.ConveyorState) {
	e.result = &ConveyorMockGetStateResult{r}
}

//Set uses given function f as a mock of Conveyor.GetState method
func (m *mConveyorMockGetState) Set(f func() (r insolar.ConveyorState)) *ConveyorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStateFunc = f
	return m.mock
}

//GetState implements github.com/insolar/insolar/insolar.Conveyor interface
func (m *ConveyorMock) GetState() (r insolar.ConveyorState) {
	counter := atomic.AddUint64(&m.GetStatePreCounter, 1)
	defer atomic.AddUint64(&m.GetStateCounter, 1)

	if len(m.GetStateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorMock.GetState.")
			return
		}

		result := m.GetStateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.GetState")
			return
		}

		r = result.r

		return
	}

	if m.GetStateMock.mainExpectation != nil {

		result := m.GetStateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.GetState")
		}

		r = result.r

		return
	}

	if m.GetStateFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorMock.GetState.")
		return
	}

	return m.GetStateFunc()
}

//GetStateMinimockCounter returns a count of ConveyorMock.GetStateFunc invocations
func (m *ConveyorMock) GetStateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStateCounter)
}

//GetStateMinimockPreCounter returns the value of ConveyorMock.GetState invocations
func (m *ConveyorMock) GetStateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStatePreCounter)
}

//GetStateFinished returns true if mock invocations count is ok
func (m *ConveyorMock) GetStateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStateCounter) == uint64(len(m.GetStateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStateFunc != nil {
		return atomic.LoadUint64(&m.GetStateCounter) > 0
	}

	return true
}

type mConveyorMockInitiateShutdown struct {
	mock              *ConveyorMock
	mainExpectation   *ConveyorMockInitiateShutdownExpectation
	expectationSeries []*ConveyorMockInitiateShutdownExpectation
}

type ConveyorMockInitiateShutdownExpectation struct {
	input *ConveyorMockInitiateShutdownInput
}

type ConveyorMockInitiateShutdownInput struct {
	p bool
}

//Expect specifies that invocation of Conveyor.InitiateShutdown is expected from 1 to Infinity times
func (m *mConveyorMockInitiateShutdown) Expect(p bool) *mConveyorMockInitiateShutdown {
	m.mock.InitiateShutdownFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockInitiateShutdownExpectation{}
	}
	m.mainExpectation.input = &ConveyorMockInitiateShutdownInput{p}
	return m
}

//Return specifies results of invocation of Conveyor.InitiateShutdown
func (m *mConveyorMockInitiateShutdown) Return() *ConveyorMock {
	m.mock.InitiateShutdownFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockInitiateShutdownExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Conveyor.InitiateShutdown is expected once
func (m *mConveyorMockInitiateShutdown) ExpectOnce(p bool) *ConveyorMockInitiateShutdownExpectation {
	m.mock.InitiateShutdownFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorMockInitiateShutdownExpectation{}
	expectation.input = &ConveyorMockInitiateShutdownInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Conveyor.InitiateShutdown method
func (m *mConveyorMockInitiateShutdown) Set(f func(p bool)) *ConveyorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InitiateShutdownFunc = f
	return m.mock
}

//InitiateShutdown implements github.com/insolar/insolar/insolar.Conveyor interface
func (m *ConveyorMock) InitiateShutdown(p bool) {
	counter := atomic.AddUint64(&m.InitiateShutdownPreCounter, 1)
	defer atomic.AddUint64(&m.InitiateShutdownCounter, 1)

	if len(m.InitiateShutdownMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InitiateShutdownMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorMock.InitiateShutdown. %v", p)
			return
		}

		input := m.InitiateShutdownMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConveyorMockInitiateShutdownInput{p}, "Conveyor.InitiateShutdown got unexpected parameters")

		return
	}

	if m.InitiateShutdownMock.mainExpectation != nil {

		input := m.InitiateShutdownMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConveyorMockInitiateShutdownInput{p}, "Conveyor.InitiateShutdown got unexpected parameters")
		}

		return
	}

	if m.InitiateShutdownFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorMock.InitiateShutdown. %v", p)
		return
	}

	m.InitiateShutdownFunc(p)
}

//InitiateShutdownMinimockCounter returns a count of ConveyorMock.InitiateShutdownFunc invocations
func (m *ConveyorMock) InitiateShutdownMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InitiateShutdownCounter)
}

//InitiateShutdownMinimockPreCounter returns the value of ConveyorMock.InitiateShutdown invocations
func (m *ConveyorMock) InitiateShutdownMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InitiateShutdownPreCounter)
}

//InitiateShutdownFinished returns true if mock invocations count is ok
func (m *ConveyorMock) InitiateShutdownFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InitiateShutdownMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InitiateShutdownCounter) == uint64(len(m.InitiateShutdownMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InitiateShutdownMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InitiateShutdownCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InitiateShutdownFunc != nil {
		return atomic.LoadUint64(&m.InitiateShutdownCounter) > 0
	}

	return true
}

type mConveyorMockIsOperational struct {
	mock              *ConveyorMock
	mainExpectation   *ConveyorMockIsOperationalExpectation
	expectationSeries []*ConveyorMockIsOperationalExpectation
}

type ConveyorMockIsOperationalExpectation struct {
	result *ConveyorMockIsOperationalResult
}

type ConveyorMockIsOperationalResult struct {
	r bool
}

//Expect specifies that invocation of Conveyor.IsOperational is expected from 1 to Infinity times
func (m *mConveyorMockIsOperational) Expect() *mConveyorMockIsOperational {
	m.mock.IsOperationalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockIsOperationalExpectation{}
	}

	return m
}

//Return specifies results of invocation of Conveyor.IsOperational
func (m *mConveyorMockIsOperational) Return(r bool) *ConveyorMock {
	m.mock.IsOperationalFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockIsOperationalExpectation{}
	}
	m.mainExpectation.result = &ConveyorMockIsOperationalResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Conveyor.IsOperational is expected once
func (m *mConveyorMockIsOperational) ExpectOnce() *ConveyorMockIsOperationalExpectation {
	m.mock.IsOperationalFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorMockIsOperationalExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorMockIsOperationalExpectation) Return(r bool) {
	e.result = &ConveyorMockIsOperationalResult{r}
}

//Set uses given function f as a mock of Conveyor.IsOperational method
func (m *mConveyorMockIsOperational) Set(f func() (r bool)) *ConveyorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsOperationalFunc = f
	return m.mock
}

//IsOperational implements github.com/insolar/insolar/insolar.Conveyor interface
func (m *ConveyorMock) IsOperational() (r bool) {
	counter := atomic.AddUint64(&m.IsOperationalPreCounter, 1)
	defer atomic.AddUint64(&m.IsOperationalCounter, 1)

	if len(m.IsOperationalMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsOperationalMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorMock.IsOperational.")
			return
		}

		result := m.IsOperationalMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.IsOperational")
			return
		}

		r = result.r

		return
	}

	if m.IsOperationalMock.mainExpectation != nil {

		result := m.IsOperationalMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.IsOperational")
		}

		r = result.r

		return
	}

	if m.IsOperationalFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorMock.IsOperational.")
		return
	}

	return m.IsOperationalFunc()
}

//IsOperationalMinimockCounter returns a count of ConveyorMock.IsOperationalFunc invocations
func (m *ConveyorMock) IsOperationalMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsOperationalCounter)
}

//IsOperationalMinimockPreCounter returns the value of ConveyorMock.IsOperational invocations
func (m *ConveyorMock) IsOperationalMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsOperationalPreCounter)
}

//IsOperationalFinished returns true if mock invocations count is ok
func (m *ConveyorMock) IsOperationalFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsOperationalMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsOperationalCounter) == uint64(len(m.IsOperationalMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsOperationalMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsOperationalCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsOperationalFunc != nil {
		return atomic.LoadUint64(&m.IsOperationalCounter) > 0
	}

	return true
}

type mConveyorMockPreparePulse struct {
	mock              *ConveyorMock
	mainExpectation   *ConveyorMockPreparePulseExpectation
	expectationSeries []*ConveyorMockPreparePulseExpectation
}

type ConveyorMockPreparePulseExpectation struct {
	input  *ConveyorMockPreparePulseInput
	result *ConveyorMockPreparePulseResult
}

type ConveyorMockPreparePulseInput struct {
	p  insolar.Pulse
	p1 queue.SyncDone
}

type ConveyorMockPreparePulseResult struct {
	r error
}

//Expect specifies that invocation of Conveyor.PreparePulse is expected from 1 to Infinity times
func (m *mConveyorMockPreparePulse) Expect(p insolar.Pulse, p1 queue.SyncDone) *mConveyorMockPreparePulse {
	m.mock.PreparePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockPreparePulseExpectation{}
	}
	m.mainExpectation.input = &ConveyorMockPreparePulseInput{p, p1}
	return m
}

//Return specifies results of invocation of Conveyor.PreparePulse
func (m *mConveyorMockPreparePulse) Return(r error) *ConveyorMock {
	m.mock.PreparePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockPreparePulseExpectation{}
	}
	m.mainExpectation.result = &ConveyorMockPreparePulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Conveyor.PreparePulse is expected once
func (m *mConveyorMockPreparePulse) ExpectOnce(p insolar.Pulse, p1 queue.SyncDone) *ConveyorMockPreparePulseExpectation {
	m.mock.PreparePulseFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorMockPreparePulseExpectation{}
	expectation.input = &ConveyorMockPreparePulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorMockPreparePulseExpectation) Return(r error) {
	e.result = &ConveyorMockPreparePulseResult{r}
}

//Set uses given function f as a mock of Conveyor.PreparePulse method
func (m *mConveyorMockPreparePulse) Set(f func(p insolar.Pulse, p1 queue.SyncDone) (r error)) *ConveyorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PreparePulseFunc = f
	return m.mock
}

//PreparePulse implements github.com/insolar/insolar/insolar.Conveyor interface
func (m *ConveyorMock) PreparePulse(p insolar.Pulse, p1 queue.SyncDone) (r error) {
	counter := atomic.AddUint64(&m.PreparePulsePreCounter, 1)
	defer atomic.AddUint64(&m.PreparePulseCounter, 1)

	if len(m.PreparePulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PreparePulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorMock.PreparePulse. %v %v", p, p1)
			return
		}

		input := m.PreparePulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConveyorMockPreparePulseInput{p, p1}, "Conveyor.PreparePulse got unexpected parameters")

		result := m.PreparePulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.PreparePulse")
			return
		}

		r = result.r

		return
	}

	if m.PreparePulseMock.mainExpectation != nil {

		input := m.PreparePulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConveyorMockPreparePulseInput{p, p1}, "Conveyor.PreparePulse got unexpected parameters")
		}

		result := m.PreparePulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.PreparePulse")
		}

		r = result.r

		return
	}

	if m.PreparePulseFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorMock.PreparePulse. %v %v", p, p1)
		return
	}

	return m.PreparePulseFunc(p, p1)
}

//PreparePulseMinimockCounter returns a count of ConveyorMock.PreparePulseFunc invocations
func (m *ConveyorMock) PreparePulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PreparePulseCounter)
}

//PreparePulseMinimockPreCounter returns the value of ConveyorMock.PreparePulse invocations
func (m *ConveyorMock) PreparePulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PreparePulsePreCounter)
}

//PreparePulseFinished returns true if mock invocations count is ok
func (m *ConveyorMock) PreparePulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PreparePulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PreparePulseCounter) == uint64(len(m.PreparePulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PreparePulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PreparePulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PreparePulseFunc != nil {
		return atomic.LoadUint64(&m.PreparePulseCounter) > 0
	}

	return true
}

type mConveyorMockSinkPush struct {
	mock              *ConveyorMock
	mainExpectation   *ConveyorMockSinkPushExpectation
	expectationSeries []*ConveyorMockSinkPushExpectation
}

type ConveyorMockSinkPushExpectation struct {
	input  *ConveyorMockSinkPushInput
	result *ConveyorMockSinkPushResult
}

type ConveyorMockSinkPushInput struct {
	p  insolar.PulseNumber
	p1 interface{}
}

type ConveyorMockSinkPushResult struct {
	r error
}

//Expect specifies that invocation of Conveyor.SinkPush is expected from 1 to Infinity times
func (m *mConveyorMockSinkPush) Expect(p insolar.PulseNumber, p1 interface{}) *mConveyorMockSinkPush {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockSinkPushExpectation{}
	}
	m.mainExpectation.input = &ConveyorMockSinkPushInput{p, p1}
	return m
}

//Return specifies results of invocation of Conveyor.SinkPush
func (m *mConveyorMockSinkPush) Return(r error) *ConveyorMock {
	m.mock.SinkPushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockSinkPushExpectation{}
	}
	m.mainExpectation.result = &ConveyorMockSinkPushResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Conveyor.SinkPush is expected once
func (m *mConveyorMockSinkPush) ExpectOnce(p insolar.PulseNumber, p1 interface{}) *ConveyorMockSinkPushExpectation {
	m.mock.SinkPushFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorMockSinkPushExpectation{}
	expectation.input = &ConveyorMockSinkPushInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorMockSinkPushExpectation) Return(r error) {
	e.result = &ConveyorMockSinkPushResult{r}
}

//Set uses given function f as a mock of Conveyor.SinkPush method
func (m *mConveyorMockSinkPush) Set(f func(p insolar.PulseNumber, p1 interface{}) (r error)) *ConveyorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushFunc = f
	return m.mock
}

//SinkPush implements github.com/insolar/insolar/insolar.Conveyor interface
func (m *ConveyorMock) SinkPush(p insolar.PulseNumber, p1 interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushCounter, 1)

	if len(m.SinkPushMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorMock.SinkPush. %v %v", p, p1)
			return
		}

		input := m.SinkPushMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConveyorMockSinkPushInput{p, p1}, "Conveyor.SinkPush got unexpected parameters")

		result := m.SinkPushMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.SinkPush")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushMock.mainExpectation != nil {

		input := m.SinkPushMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConveyorMockSinkPushInput{p, p1}, "Conveyor.SinkPush got unexpected parameters")
		}

		result := m.SinkPushMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.SinkPush")
		}

		r = result.r

		return
	}

	if m.SinkPushFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorMock.SinkPush. %v %v", p, p1)
		return
	}

	return m.SinkPushFunc(p, p1)
}

//SinkPushMinimockCounter returns a count of ConveyorMock.SinkPushFunc invocations
func (m *ConveyorMock) SinkPushMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushCounter)
}

//SinkPushMinimockPreCounter returns the value of ConveyorMock.SinkPush invocations
func (m *ConveyorMock) SinkPushMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushPreCounter)
}

//SinkPushFinished returns true if mock invocations count is ok
func (m *ConveyorMock) SinkPushFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SinkPushMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SinkPushCounter) == uint64(len(m.SinkPushMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SinkPushMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SinkPushCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SinkPushFunc != nil {
		return atomic.LoadUint64(&m.SinkPushCounter) > 0
	}

	return true
}

type mConveyorMockSinkPushAll struct {
	mock              *ConveyorMock
	mainExpectation   *ConveyorMockSinkPushAllExpectation
	expectationSeries []*ConveyorMockSinkPushAllExpectation
}

type ConveyorMockSinkPushAllExpectation struct {
	input  *ConveyorMockSinkPushAllInput
	result *ConveyorMockSinkPushAllResult
}

type ConveyorMockSinkPushAllInput struct {
	p  insolar.PulseNumber
	p1 []interface{}
}

type ConveyorMockSinkPushAllResult struct {
	r error
}

//Expect specifies that invocation of Conveyor.SinkPushAll is expected from 1 to Infinity times
func (m *mConveyorMockSinkPushAll) Expect(p insolar.PulseNumber, p1 []interface{}) *mConveyorMockSinkPushAll {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockSinkPushAllExpectation{}
	}
	m.mainExpectation.input = &ConveyorMockSinkPushAllInput{p, p1}
	return m
}

//Return specifies results of invocation of Conveyor.SinkPushAll
func (m *mConveyorMockSinkPushAll) Return(r error) *ConveyorMock {
	m.mock.SinkPushAllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ConveyorMockSinkPushAllExpectation{}
	}
	m.mainExpectation.result = &ConveyorMockSinkPushAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Conveyor.SinkPushAll is expected once
func (m *mConveyorMockSinkPushAll) ExpectOnce(p insolar.PulseNumber, p1 []interface{}) *ConveyorMockSinkPushAllExpectation {
	m.mock.SinkPushAllFunc = nil
	m.mainExpectation = nil

	expectation := &ConveyorMockSinkPushAllExpectation{}
	expectation.input = &ConveyorMockSinkPushAllInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ConveyorMockSinkPushAllExpectation) Return(r error) {
	e.result = &ConveyorMockSinkPushAllResult{r}
}

//Set uses given function f as a mock of Conveyor.SinkPushAll method
func (m *mConveyorMockSinkPushAll) Set(f func(p insolar.PulseNumber, p1 []interface{}) (r error)) *ConveyorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SinkPushAllFunc = f
	return m.mock
}

//SinkPushAll implements github.com/insolar/insolar/insolar.Conveyor interface
func (m *ConveyorMock) SinkPushAll(p insolar.PulseNumber, p1 []interface{}) (r error) {
	counter := atomic.AddUint64(&m.SinkPushAllPreCounter, 1)
	defer atomic.AddUint64(&m.SinkPushAllCounter, 1)

	if len(m.SinkPushAllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SinkPushAllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ConveyorMock.SinkPushAll. %v %v", p, p1)
			return
		}

		input := m.SinkPushAllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ConveyorMockSinkPushAllInput{p, p1}, "Conveyor.SinkPushAll got unexpected parameters")

		result := m.SinkPushAllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.SinkPushAll")
			return
		}

		r = result.r

		return
	}

	if m.SinkPushAllMock.mainExpectation != nil {

		input := m.SinkPushAllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ConveyorMockSinkPushAllInput{p, p1}, "Conveyor.SinkPushAll got unexpected parameters")
		}

		result := m.SinkPushAllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ConveyorMock.SinkPushAll")
		}

		r = result.r

		return
	}

	if m.SinkPushAllFunc == nil {
		m.t.Fatalf("Unexpected call to ConveyorMock.SinkPushAll. %v %v", p, p1)
		return
	}

	return m.SinkPushAllFunc(p, p1)
}

//SinkPushAllMinimockCounter returns a count of ConveyorMock.SinkPushAllFunc invocations
func (m *ConveyorMock) SinkPushAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllCounter)
}

//SinkPushAllMinimockPreCounter returns the value of ConveyorMock.SinkPushAll invocations
func (m *ConveyorMock) SinkPushAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SinkPushAllPreCounter)
}

//SinkPushAllFinished returns true if mock invocations count is ok
func (m *ConveyorMock) SinkPushAllFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SinkPushAllMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SinkPushAllCounter) == uint64(len(m.SinkPushAllMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SinkPushAllMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SinkPushAllCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SinkPushAllFunc != nil {
		return atomic.LoadUint64(&m.SinkPushAllCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ConveyorMock) ValidateCallCounters() {

	if !m.ActivatePulseFinished() {
		m.t.Fatal("Expected call to ConveyorMock.ActivatePulse")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to ConveyorMock.GetState")
	}

	if !m.InitiateShutdownFinished() {
		m.t.Fatal("Expected call to ConveyorMock.InitiateShutdown")
	}

	if !m.IsOperationalFinished() {
		m.t.Fatal("Expected call to ConveyorMock.IsOperational")
	}

	if !m.PreparePulseFinished() {
		m.t.Fatal("Expected call to ConveyorMock.PreparePulse")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to ConveyorMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to ConveyorMock.SinkPushAll")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ConveyorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ConveyorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ConveyorMock) MinimockFinish() {

	if !m.ActivatePulseFinished() {
		m.t.Fatal("Expected call to ConveyorMock.ActivatePulse")
	}

	if !m.GetStateFinished() {
		m.t.Fatal("Expected call to ConveyorMock.GetState")
	}

	if !m.InitiateShutdownFinished() {
		m.t.Fatal("Expected call to ConveyorMock.InitiateShutdown")
	}

	if !m.IsOperationalFinished() {
		m.t.Fatal("Expected call to ConveyorMock.IsOperational")
	}

	if !m.PreparePulseFinished() {
		m.t.Fatal("Expected call to ConveyorMock.PreparePulse")
	}

	if !m.SinkPushFinished() {
		m.t.Fatal("Expected call to ConveyorMock.SinkPush")
	}

	if !m.SinkPushAllFinished() {
		m.t.Fatal("Expected call to ConveyorMock.SinkPushAll")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ConveyorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ConveyorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ActivatePulseFinished()
		ok = ok && m.GetStateFinished()
		ok = ok && m.InitiateShutdownFinished()
		ok = ok && m.IsOperationalFinished()
		ok = ok && m.PreparePulseFinished()
		ok = ok && m.SinkPushFinished()
		ok = ok && m.SinkPushAllFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ActivatePulseFinished() {
				m.t.Error("Expected call to ConveyorMock.ActivatePulse")
			}

			if !m.GetStateFinished() {
				m.t.Error("Expected call to ConveyorMock.GetState")
			}

			if !m.InitiateShutdownFinished() {
				m.t.Error("Expected call to ConveyorMock.InitiateShutdown")
			}

			if !m.IsOperationalFinished() {
				m.t.Error("Expected call to ConveyorMock.IsOperational")
			}

			if !m.PreparePulseFinished() {
				m.t.Error("Expected call to ConveyorMock.PreparePulse")
			}

			if !m.SinkPushFinished() {
				m.t.Error("Expected call to ConveyorMock.SinkPush")
			}

			if !m.SinkPushAllFinished() {
				m.t.Error("Expected call to ConveyorMock.SinkPushAll")
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
func (m *ConveyorMock) AllMocksCalled() bool {

	if !m.ActivatePulseFinished() {
		return false
	}

	if !m.GetStateFinished() {
		return false
	}

	if !m.InitiateShutdownFinished() {
		return false
	}

	if !m.IsOperationalFinished() {
		return false
	}

	if !m.PreparePulseFinished() {
		return false
	}

	if !m.SinkPushFinished() {
		return false
	}

	if !m.SinkPushAllFinished() {
		return false
	}

	return true
}
