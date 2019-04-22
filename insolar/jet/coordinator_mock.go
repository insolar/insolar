package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Coordinator" can be found in github.com/insolar/insolar/insolar/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//CoordinatorMock implements github.com/insolar/insolar/insolar/jet.Coordinator
type CoordinatorMock struct {
	t minimock.Tester

	HeavyFunc       func(p context.Context, p1 insolar.PulseNumber) (r *insolar.Reference, r1 error)
	HeavyCounter    uint64
	HeavyPreCounter uint64
	HeavyMock       mCoordinatorMockHeavy

	IsAuthorizedFunc       func(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.Reference) (r bool, r1 error)
	IsAuthorizedCounter    uint64
	IsAuthorizedPreCounter uint64
	IsAuthorizedMock       mCoordinatorMockIsAuthorized

	IsBeyondLimitFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) (r bool, r1 error)
	IsBeyondLimitCounter    uint64
	IsBeyondLimitPreCounter uint64
	IsBeyondLimitMock       mCoordinatorMockIsBeyondLimit

	LightExecutorForJetFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error)
	LightExecutorForJetCounter    uint64
	LightExecutorForJetPreCounter uint64
	LightExecutorForJetMock       mCoordinatorMockLightExecutorForJet

	LightExecutorForObjectFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error)
	LightExecutorForObjectCounter    uint64
	LightExecutorForObjectPreCounter uint64
	LightExecutorForObjectMock       mCoordinatorMockLightExecutorForObject

	LightValidatorsForJetFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error)
	LightValidatorsForJetCounter    uint64
	LightValidatorsForJetPreCounter uint64
	LightValidatorsForJetMock       mCoordinatorMockLightValidatorsForJet

	LightValidatorsForObjectFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error)
	LightValidatorsForObjectCounter    uint64
	LightValidatorsForObjectPreCounter uint64
	LightValidatorsForObjectMock       mCoordinatorMockLightValidatorsForObject

	MeFunc       func() (r insolar.Reference)
	MeCounter    uint64
	MePreCounter uint64
	MeMock       mCoordinatorMockMe

	NodeForJetFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) (r *insolar.Reference, r1 error)
	NodeForJetCounter    uint64
	NodeForJetPreCounter uint64
	NodeForJetMock       mCoordinatorMockNodeForJet

	NodeForObjectFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) (r *insolar.Reference, r1 error)
	NodeForObjectCounter    uint64
	NodeForObjectPreCounter uint64
	NodeForObjectMock       mCoordinatorMockNodeForObject

	QueryRoleFunc       func(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber) (r []insolar.Reference, r1 error)
	QueryRoleCounter    uint64
	QueryRolePreCounter uint64
	QueryRoleMock       mCoordinatorMockQueryRole

	VirtualExecutorForObjectFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error)
	VirtualExecutorForObjectCounter    uint64
	VirtualExecutorForObjectPreCounter uint64
	VirtualExecutorForObjectMock       mCoordinatorMockVirtualExecutorForObject

	VirtualValidatorsForObjectFunc       func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error)
	VirtualValidatorsForObjectCounter    uint64
	VirtualValidatorsForObjectPreCounter uint64
	VirtualValidatorsForObjectMock       mCoordinatorMockVirtualValidatorsForObject
}

//NewCoordinatorMock returns a mock for github.com/insolar/insolar/insolar/jet.Coordinator
func NewCoordinatorMock(t minimock.Tester) *CoordinatorMock {
	m := &CoordinatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.HeavyMock = mCoordinatorMockHeavy{mock: m}
	m.IsAuthorizedMock = mCoordinatorMockIsAuthorized{mock: m}
	m.IsBeyondLimitMock = mCoordinatorMockIsBeyondLimit{mock: m}
	m.LightExecutorForJetMock = mCoordinatorMockLightExecutorForJet{mock: m}
	m.LightExecutorForObjectMock = mCoordinatorMockLightExecutorForObject{mock: m}
	m.LightValidatorsForJetMock = mCoordinatorMockLightValidatorsForJet{mock: m}
	m.LightValidatorsForObjectMock = mCoordinatorMockLightValidatorsForObject{mock: m}
	m.MeMock = mCoordinatorMockMe{mock: m}
	m.NodeForJetMock = mCoordinatorMockNodeForJet{mock: m}
	m.NodeForObjectMock = mCoordinatorMockNodeForObject{mock: m}
	m.QueryRoleMock = mCoordinatorMockQueryRole{mock: m}
	m.VirtualExecutorForObjectMock = mCoordinatorMockVirtualExecutorForObject{mock: m}
	m.VirtualValidatorsForObjectMock = mCoordinatorMockVirtualValidatorsForObject{mock: m}

	return m
}

type mCoordinatorMockHeavy struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockHeavyExpectation
	expectationSeries []*CoordinatorMockHeavyExpectation
}

type CoordinatorMockHeavyExpectation struct {
	input  *CoordinatorMockHeavyInput
	result *CoordinatorMockHeavyResult
}

type CoordinatorMockHeavyInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type CoordinatorMockHeavyResult struct {
	r  *insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.Heavy is expected from 1 to Infinity times
func (m *mCoordinatorMockHeavy) Expect(p context.Context, p1 insolar.PulseNumber) *mCoordinatorMockHeavy {
	m.mock.HeavyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockHeavyExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockHeavyInput{p, p1}
	return m
}

//Return specifies results of invocation of Coordinator.Heavy
func (m *mCoordinatorMockHeavy) Return(r *insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.HeavyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockHeavyExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockHeavyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.Heavy is expected once
func (m *mCoordinatorMockHeavy) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *CoordinatorMockHeavyExpectation {
	m.mock.HeavyFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockHeavyExpectation{}
	expectation.input = &CoordinatorMockHeavyInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockHeavyExpectation) Return(r *insolar.Reference, r1 error) {
	e.result = &CoordinatorMockHeavyResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.Heavy method
func (m *mCoordinatorMockHeavy) Set(f func(p context.Context, p1 insolar.PulseNumber) (r *insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HeavyFunc = f
	return m.mock
}

//Heavy implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) Heavy(p context.Context, p1 insolar.PulseNumber) (r *insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.HeavyPreCounter, 1)
	defer atomic.AddUint64(&m.HeavyCounter, 1)

	if len(m.HeavyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HeavyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.Heavy. %v %v", p, p1)
			return
		}

		input := m.HeavyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockHeavyInput{p, p1}, "Coordinator.Heavy got unexpected parameters")

		result := m.HeavyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.Heavy")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HeavyMock.mainExpectation != nil {

		input := m.HeavyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockHeavyInput{p, p1}, "Coordinator.Heavy got unexpected parameters")
		}

		result := m.HeavyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.Heavy")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HeavyFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.Heavy. %v %v", p, p1)
		return
	}

	return m.HeavyFunc(p, p1)
}

//HeavyMinimockCounter returns a count of CoordinatorMock.HeavyFunc invocations
func (m *CoordinatorMock) HeavyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HeavyCounter)
}

//HeavyMinimockPreCounter returns the value of CoordinatorMock.Heavy invocations
func (m *CoordinatorMock) HeavyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HeavyPreCounter)
}

//HeavyFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) HeavyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HeavyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HeavyCounter) == uint64(len(m.HeavyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HeavyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HeavyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HeavyFunc != nil {
		return atomic.LoadUint64(&m.HeavyCounter) > 0
	}

	return true
}

type mCoordinatorMockIsAuthorized struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockIsAuthorizedExpectation
	expectationSeries []*CoordinatorMockIsAuthorizedExpectation
}

type CoordinatorMockIsAuthorizedExpectation struct {
	input  *CoordinatorMockIsAuthorizedInput
	result *CoordinatorMockIsAuthorizedResult
}

type CoordinatorMockIsAuthorizedInput struct {
	p  context.Context
	p1 insolar.DynamicRole
	p2 insolar.ID
	p3 insolar.PulseNumber
	p4 insolar.Reference
}

type CoordinatorMockIsAuthorizedResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of Coordinator.IsAuthorized is expected from 1 to Infinity times
func (m *mCoordinatorMockIsAuthorized) Expect(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.Reference) *mCoordinatorMockIsAuthorized {
	m.mock.IsAuthorizedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockIsAuthorizedExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of Coordinator.IsAuthorized
func (m *mCoordinatorMockIsAuthorized) Return(r bool, r1 error) *CoordinatorMock {
	m.mock.IsAuthorizedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockIsAuthorizedExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockIsAuthorizedResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.IsAuthorized is expected once
func (m *mCoordinatorMockIsAuthorized) ExpectOnce(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.Reference) *CoordinatorMockIsAuthorizedExpectation {
	m.mock.IsAuthorizedFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockIsAuthorizedExpectation{}
	expectation.input = &CoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockIsAuthorizedExpectation) Return(r bool, r1 error) {
	e.result = &CoordinatorMockIsAuthorizedResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.IsAuthorized method
func (m *mCoordinatorMockIsAuthorized) Set(f func(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.Reference) (r bool, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsAuthorizedFunc = f
	return m.mock
}

//IsAuthorized implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) IsAuthorized(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber, p4 insolar.Reference) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.IsAuthorizedPreCounter, 1)
	defer atomic.AddUint64(&m.IsAuthorizedCounter, 1)

	if len(m.IsAuthorizedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsAuthorizedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.IsAuthorized. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.IsAuthorizedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}, "Coordinator.IsAuthorized got unexpected parameters")

		result := m.IsAuthorizedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.IsAuthorized")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IsAuthorizedMock.mainExpectation != nil {

		input := m.IsAuthorizedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockIsAuthorizedInput{p, p1, p2, p3, p4}, "Coordinator.IsAuthorized got unexpected parameters")
		}

		result := m.IsAuthorizedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.IsAuthorized")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IsAuthorizedFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.IsAuthorized. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.IsAuthorizedFunc(p, p1, p2, p3, p4)
}

//IsAuthorizedMinimockCounter returns a count of CoordinatorMock.IsAuthorizedFunc invocations
func (m *CoordinatorMock) IsAuthorizedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAuthorizedCounter)
}

//IsAuthorizedMinimockPreCounter returns the value of CoordinatorMock.IsAuthorized invocations
func (m *CoordinatorMock) IsAuthorizedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAuthorizedPreCounter)
}

//IsAuthorizedFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) IsAuthorizedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsAuthorizedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsAuthorizedCounter) == uint64(len(m.IsAuthorizedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsAuthorizedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsAuthorizedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsAuthorizedFunc != nil {
		return atomic.LoadUint64(&m.IsAuthorizedCounter) > 0
	}

	return true
}

type mCoordinatorMockIsBeyondLimit struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockIsBeyondLimitExpectation
	expectationSeries []*CoordinatorMockIsBeyondLimitExpectation
}

type CoordinatorMockIsBeyondLimitExpectation struct {
	input  *CoordinatorMockIsBeyondLimitInput
	result *CoordinatorMockIsBeyondLimitResult
}

type CoordinatorMockIsBeyondLimitInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.PulseNumber
}

type CoordinatorMockIsBeyondLimitResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of Coordinator.IsBeyondLimit is expected from 1 to Infinity times
func (m *mCoordinatorMockIsBeyondLimit) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) *mCoordinatorMockIsBeyondLimit {
	m.mock.IsBeyondLimitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockIsBeyondLimitExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockIsBeyondLimitInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Coordinator.IsBeyondLimit
func (m *mCoordinatorMockIsBeyondLimit) Return(r bool, r1 error) *CoordinatorMock {
	m.mock.IsBeyondLimitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockIsBeyondLimitExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockIsBeyondLimitResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.IsBeyondLimit is expected once
func (m *mCoordinatorMockIsBeyondLimit) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) *CoordinatorMockIsBeyondLimitExpectation {
	m.mock.IsBeyondLimitFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockIsBeyondLimitExpectation{}
	expectation.input = &CoordinatorMockIsBeyondLimitInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockIsBeyondLimitExpectation) Return(r bool, r1 error) {
	e.result = &CoordinatorMockIsBeyondLimitResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.IsBeyondLimit method
func (m *mCoordinatorMockIsBeyondLimit) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) (r bool, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsBeyondLimitFunc = f
	return m.mock
}

//IsBeyondLimit implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) IsBeyondLimit(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.IsBeyondLimitPreCounter, 1)
	defer atomic.AddUint64(&m.IsBeyondLimitCounter, 1)

	if len(m.IsBeyondLimitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsBeyondLimitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.IsBeyondLimit. %v %v %v", p, p1, p2)
			return
		}

		input := m.IsBeyondLimitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockIsBeyondLimitInput{p, p1, p2}, "Coordinator.IsBeyondLimit got unexpected parameters")

		result := m.IsBeyondLimitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.IsBeyondLimit")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IsBeyondLimitMock.mainExpectation != nil {

		input := m.IsBeyondLimitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockIsBeyondLimitInput{p, p1, p2}, "Coordinator.IsBeyondLimit got unexpected parameters")
		}

		result := m.IsBeyondLimitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.IsBeyondLimit")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IsBeyondLimitFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.IsBeyondLimit. %v %v %v", p, p1, p2)
		return
	}

	return m.IsBeyondLimitFunc(p, p1, p2)
}

//IsBeyondLimitMinimockCounter returns a count of CoordinatorMock.IsBeyondLimitFunc invocations
func (m *CoordinatorMock) IsBeyondLimitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsBeyondLimitCounter)
}

//IsBeyondLimitMinimockPreCounter returns the value of CoordinatorMock.IsBeyondLimit invocations
func (m *CoordinatorMock) IsBeyondLimitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsBeyondLimitPreCounter)
}

//IsBeyondLimitFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) IsBeyondLimitFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsBeyondLimitMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsBeyondLimitCounter) == uint64(len(m.IsBeyondLimitMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsBeyondLimitMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsBeyondLimitCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsBeyondLimitFunc != nil {
		return atomic.LoadUint64(&m.IsBeyondLimitCounter) > 0
	}

	return true
}

type mCoordinatorMockLightExecutorForJet struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockLightExecutorForJetExpectation
	expectationSeries []*CoordinatorMockLightExecutorForJetExpectation
}

type CoordinatorMockLightExecutorForJetExpectation struct {
	input  *CoordinatorMockLightExecutorForJetInput
	result *CoordinatorMockLightExecutorForJetResult
}

type CoordinatorMockLightExecutorForJetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type CoordinatorMockLightExecutorForJetResult struct {
	r  *insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.LightExecutorForJet is expected from 1 to Infinity times
func (m *mCoordinatorMockLightExecutorForJet) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mCoordinatorMockLightExecutorForJet {
	m.mock.LightExecutorForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightExecutorForJetExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockLightExecutorForJetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Coordinator.LightExecutorForJet
func (m *mCoordinatorMockLightExecutorForJet) Return(r *insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.LightExecutorForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightExecutorForJetExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockLightExecutorForJetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.LightExecutorForJet is expected once
func (m *mCoordinatorMockLightExecutorForJet) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *CoordinatorMockLightExecutorForJetExpectation {
	m.mock.LightExecutorForJetFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockLightExecutorForJetExpectation{}
	expectation.input = &CoordinatorMockLightExecutorForJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockLightExecutorForJetExpectation) Return(r *insolar.Reference, r1 error) {
	e.result = &CoordinatorMockLightExecutorForJetResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.LightExecutorForJet method
func (m *mCoordinatorMockLightExecutorForJet) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LightExecutorForJetFunc = f
	return m.mock
}

//LightExecutorForJet implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) LightExecutorForJet(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.LightExecutorForJetPreCounter, 1)
	defer atomic.AddUint64(&m.LightExecutorForJetCounter, 1)

	if len(m.LightExecutorForJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LightExecutorForJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.LightExecutorForJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.LightExecutorForJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockLightExecutorForJetInput{p, p1, p2}, "Coordinator.LightExecutorForJet got unexpected parameters")

		result := m.LightExecutorForJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightExecutorForJet")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightExecutorForJetMock.mainExpectation != nil {

		input := m.LightExecutorForJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockLightExecutorForJetInput{p, p1, p2}, "Coordinator.LightExecutorForJet got unexpected parameters")
		}

		result := m.LightExecutorForJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightExecutorForJet")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightExecutorForJetFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.LightExecutorForJet. %v %v %v", p, p1, p2)
		return
	}

	return m.LightExecutorForJetFunc(p, p1, p2)
}

//LightExecutorForJetMinimockCounter returns a count of CoordinatorMock.LightExecutorForJetFunc invocations
func (m *CoordinatorMock) LightExecutorForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForJetCounter)
}

//LightExecutorForJetMinimockPreCounter returns the value of CoordinatorMock.LightExecutorForJet invocations
func (m *CoordinatorMock) LightExecutorForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForJetPreCounter)
}

//LightExecutorForJetFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) LightExecutorForJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LightExecutorForJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LightExecutorForJetCounter) == uint64(len(m.LightExecutorForJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LightExecutorForJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LightExecutorForJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LightExecutorForJetFunc != nil {
		return atomic.LoadUint64(&m.LightExecutorForJetCounter) > 0
	}

	return true
}

type mCoordinatorMockLightExecutorForObject struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockLightExecutorForObjectExpectation
	expectationSeries []*CoordinatorMockLightExecutorForObjectExpectation
}

type CoordinatorMockLightExecutorForObjectExpectation struct {
	input  *CoordinatorMockLightExecutorForObjectInput
	result *CoordinatorMockLightExecutorForObjectResult
}

type CoordinatorMockLightExecutorForObjectInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type CoordinatorMockLightExecutorForObjectResult struct {
	r  *insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.LightExecutorForObject is expected from 1 to Infinity times
func (m *mCoordinatorMockLightExecutorForObject) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mCoordinatorMockLightExecutorForObject {
	m.mock.LightExecutorForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightExecutorForObjectExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockLightExecutorForObjectInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Coordinator.LightExecutorForObject
func (m *mCoordinatorMockLightExecutorForObject) Return(r *insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.LightExecutorForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightExecutorForObjectExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockLightExecutorForObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.LightExecutorForObject is expected once
func (m *mCoordinatorMockLightExecutorForObject) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *CoordinatorMockLightExecutorForObjectExpectation {
	m.mock.LightExecutorForObjectFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockLightExecutorForObjectExpectation{}
	expectation.input = &CoordinatorMockLightExecutorForObjectInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockLightExecutorForObjectExpectation) Return(r *insolar.Reference, r1 error) {
	e.result = &CoordinatorMockLightExecutorForObjectResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.LightExecutorForObject method
func (m *mCoordinatorMockLightExecutorForObject) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LightExecutorForObjectFunc = f
	return m.mock
}

//LightExecutorForObject implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) LightExecutorForObject(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.LightExecutorForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.LightExecutorForObjectCounter, 1)

	if len(m.LightExecutorForObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LightExecutorForObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.LightExecutorForObject. %v %v %v", p, p1, p2)
			return
		}

		input := m.LightExecutorForObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockLightExecutorForObjectInput{p, p1, p2}, "Coordinator.LightExecutorForObject got unexpected parameters")

		result := m.LightExecutorForObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightExecutorForObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightExecutorForObjectMock.mainExpectation != nil {

		input := m.LightExecutorForObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockLightExecutorForObjectInput{p, p1, p2}, "Coordinator.LightExecutorForObject got unexpected parameters")
		}

		result := m.LightExecutorForObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightExecutorForObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightExecutorForObjectFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.LightExecutorForObject. %v %v %v", p, p1, p2)
		return
	}

	return m.LightExecutorForObjectFunc(p, p1, p2)
}

//LightExecutorForObjectMinimockCounter returns a count of CoordinatorMock.LightExecutorForObjectFunc invocations
func (m *CoordinatorMock) LightExecutorForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForObjectCounter)
}

//LightExecutorForObjectMinimockPreCounter returns the value of CoordinatorMock.LightExecutorForObject invocations
func (m *CoordinatorMock) LightExecutorForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForObjectPreCounter)
}

//LightExecutorForObjectFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) LightExecutorForObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LightExecutorForObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LightExecutorForObjectCounter) == uint64(len(m.LightExecutorForObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LightExecutorForObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LightExecutorForObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LightExecutorForObjectFunc != nil {
		return atomic.LoadUint64(&m.LightExecutorForObjectCounter) > 0
	}

	return true
}

type mCoordinatorMockLightValidatorsForJet struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockLightValidatorsForJetExpectation
	expectationSeries []*CoordinatorMockLightValidatorsForJetExpectation
}

type CoordinatorMockLightValidatorsForJetExpectation struct {
	input  *CoordinatorMockLightValidatorsForJetInput
	result *CoordinatorMockLightValidatorsForJetResult
}

type CoordinatorMockLightValidatorsForJetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type CoordinatorMockLightValidatorsForJetResult struct {
	r  []insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.LightValidatorsForJet is expected from 1 to Infinity times
func (m *mCoordinatorMockLightValidatorsForJet) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mCoordinatorMockLightValidatorsForJet {
	m.mock.LightValidatorsForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightValidatorsForJetExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockLightValidatorsForJetInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Coordinator.LightValidatorsForJet
func (m *mCoordinatorMockLightValidatorsForJet) Return(r []insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.LightValidatorsForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightValidatorsForJetExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockLightValidatorsForJetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.LightValidatorsForJet is expected once
func (m *mCoordinatorMockLightValidatorsForJet) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *CoordinatorMockLightValidatorsForJetExpectation {
	m.mock.LightValidatorsForJetFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockLightValidatorsForJetExpectation{}
	expectation.input = &CoordinatorMockLightValidatorsForJetInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockLightValidatorsForJetExpectation) Return(r []insolar.Reference, r1 error) {
	e.result = &CoordinatorMockLightValidatorsForJetResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.LightValidatorsForJet method
func (m *mCoordinatorMockLightValidatorsForJet) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LightValidatorsForJetFunc = f
	return m.mock
}

//LightValidatorsForJet implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) LightValidatorsForJet(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.LightValidatorsForJetPreCounter, 1)
	defer atomic.AddUint64(&m.LightValidatorsForJetCounter, 1)

	if len(m.LightValidatorsForJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LightValidatorsForJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.LightValidatorsForJet. %v %v %v", p, p1, p2)
			return
		}

		input := m.LightValidatorsForJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockLightValidatorsForJetInput{p, p1, p2}, "Coordinator.LightValidatorsForJet got unexpected parameters")

		result := m.LightValidatorsForJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightValidatorsForJet")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightValidatorsForJetMock.mainExpectation != nil {

		input := m.LightValidatorsForJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockLightValidatorsForJetInput{p, p1, p2}, "Coordinator.LightValidatorsForJet got unexpected parameters")
		}

		result := m.LightValidatorsForJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightValidatorsForJet")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightValidatorsForJetFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.LightValidatorsForJet. %v %v %v", p, p1, p2)
		return
	}

	return m.LightValidatorsForJetFunc(p, p1, p2)
}

//LightValidatorsForJetMinimockCounter returns a count of CoordinatorMock.LightValidatorsForJetFunc invocations
func (m *CoordinatorMock) LightValidatorsForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForJetCounter)
}

//LightValidatorsForJetMinimockPreCounter returns the value of CoordinatorMock.LightValidatorsForJet invocations
func (m *CoordinatorMock) LightValidatorsForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForJetPreCounter)
}

//LightValidatorsForJetFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) LightValidatorsForJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LightValidatorsForJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LightValidatorsForJetCounter) == uint64(len(m.LightValidatorsForJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LightValidatorsForJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LightValidatorsForJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LightValidatorsForJetFunc != nil {
		return atomic.LoadUint64(&m.LightValidatorsForJetCounter) > 0
	}

	return true
}

type mCoordinatorMockLightValidatorsForObject struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockLightValidatorsForObjectExpectation
	expectationSeries []*CoordinatorMockLightValidatorsForObjectExpectation
}

type CoordinatorMockLightValidatorsForObjectExpectation struct {
	input  *CoordinatorMockLightValidatorsForObjectInput
	result *CoordinatorMockLightValidatorsForObjectResult
}

type CoordinatorMockLightValidatorsForObjectInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type CoordinatorMockLightValidatorsForObjectResult struct {
	r  []insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.LightValidatorsForObject is expected from 1 to Infinity times
func (m *mCoordinatorMockLightValidatorsForObject) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mCoordinatorMockLightValidatorsForObject {
	m.mock.LightValidatorsForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightValidatorsForObjectExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockLightValidatorsForObjectInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Coordinator.LightValidatorsForObject
func (m *mCoordinatorMockLightValidatorsForObject) Return(r []insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.LightValidatorsForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockLightValidatorsForObjectExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockLightValidatorsForObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.LightValidatorsForObject is expected once
func (m *mCoordinatorMockLightValidatorsForObject) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *CoordinatorMockLightValidatorsForObjectExpectation {
	m.mock.LightValidatorsForObjectFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockLightValidatorsForObjectExpectation{}
	expectation.input = &CoordinatorMockLightValidatorsForObjectInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockLightValidatorsForObjectExpectation) Return(r []insolar.Reference, r1 error) {
	e.result = &CoordinatorMockLightValidatorsForObjectResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.LightValidatorsForObject method
func (m *mCoordinatorMockLightValidatorsForObject) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LightValidatorsForObjectFunc = f
	return m.mock
}

//LightValidatorsForObject implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) LightValidatorsForObject(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.LightValidatorsForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.LightValidatorsForObjectCounter, 1)

	if len(m.LightValidatorsForObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LightValidatorsForObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.LightValidatorsForObject. %v %v %v", p, p1, p2)
			return
		}

		input := m.LightValidatorsForObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockLightValidatorsForObjectInput{p, p1, p2}, "Coordinator.LightValidatorsForObject got unexpected parameters")

		result := m.LightValidatorsForObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightValidatorsForObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightValidatorsForObjectMock.mainExpectation != nil {

		input := m.LightValidatorsForObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockLightValidatorsForObjectInput{p, p1, p2}, "Coordinator.LightValidatorsForObject got unexpected parameters")
		}

		result := m.LightValidatorsForObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.LightValidatorsForObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LightValidatorsForObjectFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.LightValidatorsForObject. %v %v %v", p, p1, p2)
		return
	}

	return m.LightValidatorsForObjectFunc(p, p1, p2)
}

//LightValidatorsForObjectMinimockCounter returns a count of CoordinatorMock.LightValidatorsForObjectFunc invocations
func (m *CoordinatorMock) LightValidatorsForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForObjectCounter)
}

//LightValidatorsForObjectMinimockPreCounter returns the value of CoordinatorMock.LightValidatorsForObject invocations
func (m *CoordinatorMock) LightValidatorsForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForObjectPreCounter)
}

//LightValidatorsForObjectFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) LightValidatorsForObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LightValidatorsForObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LightValidatorsForObjectCounter) == uint64(len(m.LightValidatorsForObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LightValidatorsForObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LightValidatorsForObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LightValidatorsForObjectFunc != nil {
		return atomic.LoadUint64(&m.LightValidatorsForObjectCounter) > 0
	}

	return true
}

type mCoordinatorMockMe struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockMeExpectation
	expectationSeries []*CoordinatorMockMeExpectation
}

type CoordinatorMockMeExpectation struct {
	result *CoordinatorMockMeResult
}

type CoordinatorMockMeResult struct {
	r insolar.Reference
}

//Expect specifies that invocation of Coordinator.Me is expected from 1 to Infinity times
func (m *mCoordinatorMockMe) Expect() *mCoordinatorMockMe {
	m.mock.MeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockMeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Coordinator.Me
func (m *mCoordinatorMockMe) Return(r insolar.Reference) *CoordinatorMock {
	m.mock.MeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockMeExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockMeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.Me is expected once
func (m *mCoordinatorMockMe) ExpectOnce() *CoordinatorMockMeExpectation {
	m.mock.MeFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockMeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockMeExpectation) Return(r insolar.Reference) {
	e.result = &CoordinatorMockMeResult{r}
}

//Set uses given function f as a mock of Coordinator.Me method
func (m *mCoordinatorMockMe) Set(f func() (r insolar.Reference)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MeFunc = f
	return m.mock
}

//Me implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) Me() (r insolar.Reference) {
	counter := atomic.AddUint64(&m.MePreCounter, 1)
	defer atomic.AddUint64(&m.MeCounter, 1)

	if len(m.MeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.Me.")
			return
		}

		result := m.MeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.Me")
			return
		}

		r = result.r

		return
	}

	if m.MeMock.mainExpectation != nil {

		result := m.MeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.Me")
		}

		r = result.r

		return
	}

	if m.MeFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.Me.")
		return
	}

	return m.MeFunc()
}

//MeMinimockCounter returns a count of CoordinatorMock.MeFunc invocations
func (m *CoordinatorMock) MeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MeCounter)
}

//MeMinimockPreCounter returns the value of CoordinatorMock.Me invocations
func (m *CoordinatorMock) MeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MePreCounter)
}

//MeFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) MeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MeCounter) == uint64(len(m.MeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MeFunc != nil {
		return atomic.LoadUint64(&m.MeCounter) > 0
	}

	return true
}

type mCoordinatorMockNodeForJet struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockNodeForJetExpectation
	expectationSeries []*CoordinatorMockNodeForJetExpectation
}

type CoordinatorMockNodeForJetExpectation struct {
	input  *CoordinatorMockNodeForJetInput
	result *CoordinatorMockNodeForJetResult
}

type CoordinatorMockNodeForJetInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 insolar.PulseNumber
}

type CoordinatorMockNodeForJetResult struct {
	r  *insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.NodeForJet is expected from 1 to Infinity times
func (m *mCoordinatorMockNodeForJet) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) *mCoordinatorMockNodeForJet {
	m.mock.NodeForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockNodeForJetExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockNodeForJetInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Coordinator.NodeForJet
func (m *mCoordinatorMockNodeForJet) Return(r *insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.NodeForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockNodeForJetExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockNodeForJetResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.NodeForJet is expected once
func (m *mCoordinatorMockNodeForJet) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) *CoordinatorMockNodeForJetExpectation {
	m.mock.NodeForJetFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockNodeForJetExpectation{}
	expectation.input = &CoordinatorMockNodeForJetInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockNodeForJetExpectation) Return(r *insolar.Reference, r1 error) {
	e.result = &CoordinatorMockNodeForJetResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.NodeForJet method
func (m *mCoordinatorMockNodeForJet) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) (r *insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NodeForJetFunc = f
	return m.mock
}

//NodeForJet implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) NodeForJet(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) (r *insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.NodeForJetPreCounter, 1)
	defer atomic.AddUint64(&m.NodeForJetCounter, 1)

	if len(m.NodeForJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NodeForJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.NodeForJet. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.NodeForJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockNodeForJetInput{p, p1, p2, p3}, "Coordinator.NodeForJet got unexpected parameters")

		result := m.NodeForJetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.NodeForJet")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NodeForJetMock.mainExpectation != nil {

		input := m.NodeForJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockNodeForJetInput{p, p1, p2, p3}, "Coordinator.NodeForJet got unexpected parameters")
		}

		result := m.NodeForJetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.NodeForJet")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NodeForJetFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.NodeForJet. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.NodeForJetFunc(p, p1, p2, p3)
}

//NodeForJetMinimockCounter returns a count of CoordinatorMock.NodeForJetFunc invocations
func (m *CoordinatorMock) NodeForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForJetCounter)
}

//NodeForJetMinimockPreCounter returns the value of CoordinatorMock.NodeForJet invocations
func (m *CoordinatorMock) NodeForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForJetPreCounter)
}

//NodeForJetFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) NodeForJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NodeForJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NodeForJetCounter) == uint64(len(m.NodeForJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NodeForJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NodeForJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NodeForJetFunc != nil {
		return atomic.LoadUint64(&m.NodeForJetCounter) > 0
	}

	return true
}

type mCoordinatorMockNodeForObject struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockNodeForObjectExpectation
	expectationSeries []*CoordinatorMockNodeForObjectExpectation
}

type CoordinatorMockNodeForObjectExpectation struct {
	input  *CoordinatorMockNodeForObjectInput
	result *CoordinatorMockNodeForObjectResult
}

type CoordinatorMockNodeForObjectInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
	p3 insolar.PulseNumber
}

type CoordinatorMockNodeForObjectResult struct {
	r  *insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.NodeForObject is expected from 1 to Infinity times
func (m *mCoordinatorMockNodeForObject) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) *mCoordinatorMockNodeForObject {
	m.mock.NodeForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockNodeForObjectExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockNodeForObjectInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Coordinator.NodeForObject
func (m *mCoordinatorMockNodeForObject) Return(r *insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.NodeForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockNodeForObjectExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockNodeForObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.NodeForObject is expected once
func (m *mCoordinatorMockNodeForObject) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) *CoordinatorMockNodeForObjectExpectation {
	m.mock.NodeForObjectFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockNodeForObjectExpectation{}
	expectation.input = &CoordinatorMockNodeForObjectInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockNodeForObjectExpectation) Return(r *insolar.Reference, r1 error) {
	e.result = &CoordinatorMockNodeForObjectResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.NodeForObject method
func (m *mCoordinatorMockNodeForObject) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) (r *insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NodeForObjectFunc = f
	return m.mock
}

//NodeForObject implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) NodeForObject(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber, p3 insolar.PulseNumber) (r *insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.NodeForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.NodeForObjectCounter, 1)

	if len(m.NodeForObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NodeForObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.NodeForObject. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.NodeForObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockNodeForObjectInput{p, p1, p2, p3}, "Coordinator.NodeForObject got unexpected parameters")

		result := m.NodeForObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.NodeForObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NodeForObjectMock.mainExpectation != nil {

		input := m.NodeForObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockNodeForObjectInput{p, p1, p2, p3}, "Coordinator.NodeForObject got unexpected parameters")
		}

		result := m.NodeForObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.NodeForObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NodeForObjectFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.NodeForObject. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.NodeForObjectFunc(p, p1, p2, p3)
}

//NodeForObjectMinimockCounter returns a count of CoordinatorMock.NodeForObjectFunc invocations
func (m *CoordinatorMock) NodeForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForObjectCounter)
}

//NodeForObjectMinimockPreCounter returns the value of CoordinatorMock.NodeForObject invocations
func (m *CoordinatorMock) NodeForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForObjectPreCounter)
}

//NodeForObjectFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) NodeForObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NodeForObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NodeForObjectCounter) == uint64(len(m.NodeForObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NodeForObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NodeForObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NodeForObjectFunc != nil {
		return atomic.LoadUint64(&m.NodeForObjectCounter) > 0
	}

	return true
}

type mCoordinatorMockQueryRole struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockQueryRoleExpectation
	expectationSeries []*CoordinatorMockQueryRoleExpectation
}

type CoordinatorMockQueryRoleExpectation struct {
	input  *CoordinatorMockQueryRoleInput
	result *CoordinatorMockQueryRoleResult
}

type CoordinatorMockQueryRoleInput struct {
	p  context.Context
	p1 insolar.DynamicRole
	p2 insolar.ID
	p3 insolar.PulseNumber
}

type CoordinatorMockQueryRoleResult struct {
	r  []insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.QueryRole is expected from 1 to Infinity times
func (m *mCoordinatorMockQueryRole) Expect(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber) *mCoordinatorMockQueryRole {
	m.mock.QueryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockQueryRoleExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockQueryRoleInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Coordinator.QueryRole
func (m *mCoordinatorMockQueryRole) Return(r []insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.QueryRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockQueryRoleExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockQueryRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.QueryRole is expected once
func (m *mCoordinatorMockQueryRole) ExpectOnce(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber) *CoordinatorMockQueryRoleExpectation {
	m.mock.QueryRoleFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockQueryRoleExpectation{}
	expectation.input = &CoordinatorMockQueryRoleInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockQueryRoleExpectation) Return(r []insolar.Reference, r1 error) {
	e.result = &CoordinatorMockQueryRoleResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.QueryRole method
func (m *mCoordinatorMockQueryRole) Set(f func(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber) (r []insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.QueryRoleFunc = f
	return m.mock
}

//QueryRole implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) QueryRole(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber) (r []insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.QueryRolePreCounter, 1)
	defer atomic.AddUint64(&m.QueryRoleCounter, 1)

	if len(m.QueryRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.QueryRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.QueryRole. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.QueryRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockQueryRoleInput{p, p1, p2, p3}, "Coordinator.QueryRole got unexpected parameters")

		result := m.QueryRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.QueryRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.QueryRoleMock.mainExpectation != nil {

		input := m.QueryRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockQueryRoleInput{p, p1, p2, p3}, "Coordinator.QueryRole got unexpected parameters")
		}

		result := m.QueryRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.QueryRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.QueryRoleFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.QueryRole. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.QueryRoleFunc(p, p1, p2, p3)
}

//QueryRoleMinimockCounter returns a count of CoordinatorMock.QueryRoleFunc invocations
func (m *CoordinatorMock) QueryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.QueryRoleCounter)
}

//QueryRoleMinimockPreCounter returns the value of CoordinatorMock.QueryRole invocations
func (m *CoordinatorMock) QueryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.QueryRolePreCounter)
}

//QueryRoleFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) QueryRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.QueryRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.QueryRoleCounter) == uint64(len(m.QueryRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.QueryRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.QueryRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.QueryRoleFunc != nil {
		return atomic.LoadUint64(&m.QueryRoleCounter) > 0
	}

	return true
}

type mCoordinatorMockVirtualExecutorForObject struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockVirtualExecutorForObjectExpectation
	expectationSeries []*CoordinatorMockVirtualExecutorForObjectExpectation
}

type CoordinatorMockVirtualExecutorForObjectExpectation struct {
	input  *CoordinatorMockVirtualExecutorForObjectInput
	result *CoordinatorMockVirtualExecutorForObjectResult
}

type CoordinatorMockVirtualExecutorForObjectInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type CoordinatorMockVirtualExecutorForObjectResult struct {
	r  *insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.VirtualExecutorForObject is expected from 1 to Infinity times
func (m *mCoordinatorMockVirtualExecutorForObject) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mCoordinatorMockVirtualExecutorForObject {
	m.mock.VirtualExecutorForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockVirtualExecutorForObjectExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockVirtualExecutorForObjectInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Coordinator.VirtualExecutorForObject
func (m *mCoordinatorMockVirtualExecutorForObject) Return(r *insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.VirtualExecutorForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockVirtualExecutorForObjectExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockVirtualExecutorForObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.VirtualExecutorForObject is expected once
func (m *mCoordinatorMockVirtualExecutorForObject) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *CoordinatorMockVirtualExecutorForObjectExpectation {
	m.mock.VirtualExecutorForObjectFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockVirtualExecutorForObjectExpectation{}
	expectation.input = &CoordinatorMockVirtualExecutorForObjectInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockVirtualExecutorForObjectExpectation) Return(r *insolar.Reference, r1 error) {
	e.result = &CoordinatorMockVirtualExecutorForObjectResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.VirtualExecutorForObject method
func (m *mCoordinatorMockVirtualExecutorForObject) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VirtualExecutorForObjectFunc = f
	return m.mock
}

//VirtualExecutorForObject implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) VirtualExecutorForObject(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r *insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.VirtualExecutorForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.VirtualExecutorForObjectCounter, 1)

	if len(m.VirtualExecutorForObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VirtualExecutorForObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.VirtualExecutorForObject. %v %v %v", p, p1, p2)
			return
		}

		input := m.VirtualExecutorForObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockVirtualExecutorForObjectInput{p, p1, p2}, "Coordinator.VirtualExecutorForObject got unexpected parameters")

		result := m.VirtualExecutorForObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.VirtualExecutorForObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VirtualExecutorForObjectMock.mainExpectation != nil {

		input := m.VirtualExecutorForObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockVirtualExecutorForObjectInput{p, p1, p2}, "Coordinator.VirtualExecutorForObject got unexpected parameters")
		}

		result := m.VirtualExecutorForObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.VirtualExecutorForObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VirtualExecutorForObjectFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.VirtualExecutorForObject. %v %v %v", p, p1, p2)
		return
	}

	return m.VirtualExecutorForObjectFunc(p, p1, p2)
}

//VirtualExecutorForObjectMinimockCounter returns a count of CoordinatorMock.VirtualExecutorForObjectFunc invocations
func (m *CoordinatorMock) VirtualExecutorForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualExecutorForObjectCounter)
}

//VirtualExecutorForObjectMinimockPreCounter returns the value of CoordinatorMock.VirtualExecutorForObject invocations
func (m *CoordinatorMock) VirtualExecutorForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualExecutorForObjectPreCounter)
}

//VirtualExecutorForObjectFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) VirtualExecutorForObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.VirtualExecutorForObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) == uint64(len(m.VirtualExecutorForObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.VirtualExecutorForObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.VirtualExecutorForObjectFunc != nil {
		return atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) > 0
	}

	return true
}

type mCoordinatorMockVirtualValidatorsForObject struct {
	mock              *CoordinatorMock
	mainExpectation   *CoordinatorMockVirtualValidatorsForObjectExpectation
	expectationSeries []*CoordinatorMockVirtualValidatorsForObjectExpectation
}

type CoordinatorMockVirtualValidatorsForObjectExpectation struct {
	input  *CoordinatorMockVirtualValidatorsForObjectInput
	result *CoordinatorMockVirtualValidatorsForObjectResult
}

type CoordinatorMockVirtualValidatorsForObjectInput struct {
	p  context.Context
	p1 insolar.ID
	p2 insolar.PulseNumber
}

type CoordinatorMockVirtualValidatorsForObjectResult struct {
	r  []insolar.Reference
	r1 error
}

//Expect specifies that invocation of Coordinator.VirtualValidatorsForObject is expected from 1 to Infinity times
func (m *mCoordinatorMockVirtualValidatorsForObject) Expect(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *mCoordinatorMockVirtualValidatorsForObject {
	m.mock.VirtualValidatorsForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockVirtualValidatorsForObjectExpectation{}
	}
	m.mainExpectation.input = &CoordinatorMockVirtualValidatorsForObjectInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Coordinator.VirtualValidatorsForObject
func (m *mCoordinatorMockVirtualValidatorsForObject) Return(r []insolar.Reference, r1 error) *CoordinatorMock {
	m.mock.VirtualValidatorsForObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CoordinatorMockVirtualValidatorsForObjectExpectation{}
	}
	m.mainExpectation.result = &CoordinatorMockVirtualValidatorsForObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Coordinator.VirtualValidatorsForObject is expected once
func (m *mCoordinatorMockVirtualValidatorsForObject) ExpectOnce(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) *CoordinatorMockVirtualValidatorsForObjectExpectation {
	m.mock.VirtualValidatorsForObjectFunc = nil
	m.mainExpectation = nil

	expectation := &CoordinatorMockVirtualValidatorsForObjectExpectation{}
	expectation.input = &CoordinatorMockVirtualValidatorsForObjectInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CoordinatorMockVirtualValidatorsForObjectExpectation) Return(r []insolar.Reference, r1 error) {
	e.result = &CoordinatorMockVirtualValidatorsForObjectResult{r, r1}
}

//Set uses given function f as a mock of Coordinator.VirtualValidatorsForObject method
func (m *mCoordinatorMockVirtualValidatorsForObject) Set(f func(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error)) *CoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VirtualValidatorsForObjectFunc = f
	return m.mock
}

//VirtualValidatorsForObject implements github.com/insolar/insolar/insolar/jet.Coordinator interface
func (m *CoordinatorMock) VirtualValidatorsForObject(p context.Context, p1 insolar.ID, p2 insolar.PulseNumber) (r []insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.VirtualValidatorsForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.VirtualValidatorsForObjectCounter, 1)

	if len(m.VirtualValidatorsForObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VirtualValidatorsForObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CoordinatorMock.VirtualValidatorsForObject. %v %v %v", p, p1, p2)
			return
		}

		input := m.VirtualValidatorsForObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CoordinatorMockVirtualValidatorsForObjectInput{p, p1, p2}, "Coordinator.VirtualValidatorsForObject got unexpected parameters")

		result := m.VirtualValidatorsForObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.VirtualValidatorsForObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VirtualValidatorsForObjectMock.mainExpectation != nil {

		input := m.VirtualValidatorsForObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CoordinatorMockVirtualValidatorsForObjectInput{p, p1, p2}, "Coordinator.VirtualValidatorsForObject got unexpected parameters")
		}

		result := m.VirtualValidatorsForObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CoordinatorMock.VirtualValidatorsForObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.VirtualValidatorsForObjectFunc == nil {
		m.t.Fatalf("Unexpected call to CoordinatorMock.VirtualValidatorsForObject. %v %v %v", p, p1, p2)
		return
	}

	return m.VirtualValidatorsForObjectFunc(p, p1, p2)
}

//VirtualValidatorsForObjectMinimockCounter returns a count of CoordinatorMock.VirtualValidatorsForObjectFunc invocations
func (m *CoordinatorMock) VirtualValidatorsForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter)
}

//VirtualValidatorsForObjectMinimockPreCounter returns the value of CoordinatorMock.VirtualValidatorsForObject invocations
func (m *CoordinatorMock) VirtualValidatorsForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualValidatorsForObjectPreCounter)
}

//VirtualValidatorsForObjectFinished returns true if mock invocations count is ok
func (m *CoordinatorMock) VirtualValidatorsForObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.VirtualValidatorsForObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) == uint64(len(m.VirtualValidatorsForObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.VirtualValidatorsForObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.VirtualValidatorsForObjectFunc != nil {
		return atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CoordinatorMock) ValidateCallCounters() {

	if !m.HeavyFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.Heavy")
	}

	if !m.IsAuthorizedFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.IsAuthorized")
	}

	if !m.IsBeyondLimitFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.IsBeyondLimit")
	}

	if !m.LightExecutorForJetFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightExecutorForJet")
	}

	if !m.LightExecutorForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightExecutorForObject")
	}

	if !m.LightValidatorsForJetFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightValidatorsForJet")
	}

	if !m.LightValidatorsForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightValidatorsForObject")
	}

	if !m.MeFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.Me")
	}

	if !m.NodeForJetFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.NodeForJet")
	}

	if !m.NodeForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.NodeForObject")
	}

	if !m.QueryRoleFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.QueryRole")
	}

	if !m.VirtualExecutorForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.VirtualExecutorForObject")
	}

	if !m.VirtualValidatorsForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.VirtualValidatorsForObject")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CoordinatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CoordinatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CoordinatorMock) MinimockFinish() {

	if !m.HeavyFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.Heavy")
	}

	if !m.IsAuthorizedFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.IsAuthorized")
	}

	if !m.IsBeyondLimitFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.IsBeyondLimit")
	}

	if !m.LightExecutorForJetFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightExecutorForJet")
	}

	if !m.LightExecutorForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightExecutorForObject")
	}

	if !m.LightValidatorsForJetFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightValidatorsForJet")
	}

	if !m.LightValidatorsForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.LightValidatorsForObject")
	}

	if !m.MeFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.Me")
	}

	if !m.NodeForJetFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.NodeForJet")
	}

	if !m.NodeForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.NodeForObject")
	}

	if !m.QueryRoleFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.QueryRole")
	}

	if !m.VirtualExecutorForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.VirtualExecutorForObject")
	}

	if !m.VirtualValidatorsForObjectFinished() {
		m.t.Fatal("Expected call to CoordinatorMock.VirtualValidatorsForObject")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CoordinatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CoordinatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.HeavyFinished()
		ok = ok && m.IsAuthorizedFinished()
		ok = ok && m.IsBeyondLimitFinished()
		ok = ok && m.LightExecutorForJetFinished()
		ok = ok && m.LightExecutorForObjectFinished()
		ok = ok && m.LightValidatorsForJetFinished()
		ok = ok && m.LightValidatorsForObjectFinished()
		ok = ok && m.MeFinished()
		ok = ok && m.NodeForJetFinished()
		ok = ok && m.NodeForObjectFinished()
		ok = ok && m.QueryRoleFinished()
		ok = ok && m.VirtualExecutorForObjectFinished()
		ok = ok && m.VirtualValidatorsForObjectFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.HeavyFinished() {
				m.t.Error("Expected call to CoordinatorMock.Heavy")
			}

			if !m.IsAuthorizedFinished() {
				m.t.Error("Expected call to CoordinatorMock.IsAuthorized")
			}

			if !m.IsBeyondLimitFinished() {
				m.t.Error("Expected call to CoordinatorMock.IsBeyondLimit")
			}

			if !m.LightExecutorForJetFinished() {
				m.t.Error("Expected call to CoordinatorMock.LightExecutorForJet")
			}

			if !m.LightExecutorForObjectFinished() {
				m.t.Error("Expected call to CoordinatorMock.LightExecutorForObject")
			}

			if !m.LightValidatorsForJetFinished() {
				m.t.Error("Expected call to CoordinatorMock.LightValidatorsForJet")
			}

			if !m.LightValidatorsForObjectFinished() {
				m.t.Error("Expected call to CoordinatorMock.LightValidatorsForObject")
			}

			if !m.MeFinished() {
				m.t.Error("Expected call to CoordinatorMock.Me")
			}

			if !m.NodeForJetFinished() {
				m.t.Error("Expected call to CoordinatorMock.NodeForJet")
			}

			if !m.NodeForObjectFinished() {
				m.t.Error("Expected call to CoordinatorMock.NodeForObject")
			}

			if !m.QueryRoleFinished() {
				m.t.Error("Expected call to CoordinatorMock.QueryRole")
			}

			if !m.VirtualExecutorForObjectFinished() {
				m.t.Error("Expected call to CoordinatorMock.VirtualExecutorForObject")
			}

			if !m.VirtualValidatorsForObjectFinished() {
				m.t.Error("Expected call to CoordinatorMock.VirtualValidatorsForObject")
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
func (m *CoordinatorMock) AllMocksCalled() bool {

	if !m.HeavyFinished() {
		return false
	}

	if !m.IsAuthorizedFinished() {
		return false
	}

	if !m.IsBeyondLimitFinished() {
		return false
	}

	if !m.LightExecutorForJetFinished() {
		return false
	}

	if !m.LightExecutorForObjectFinished() {
		return false
	}

	if !m.LightValidatorsForJetFinished() {
		return false
	}

	if !m.LightValidatorsForObjectFinished() {
		return false
	}

	if !m.MeFinished() {
		return false
	}

	if !m.NodeForJetFinished() {
		return false
	}

	if !m.NodeForObjectFinished() {
		return false
	}

	if !m.QueryRoleFinished() {
		return false
	}

	if !m.VirtualExecutorForObjectFinished() {
		return false
	}

	if !m.VirtualValidatorsForObjectFinished() {
		return false
	}

	return true
}
