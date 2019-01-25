package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ActiveNodesStorage" can be found in github.com/insolar/insolar/ledger/storage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ActiveNodesStorageMock implements github.com/insolar/insolar/ledger/storage.ActiveNodesStorage
type ActiveNodesStorageMock struct {
	t minimock.Tester

	GetActiveNodesFunc       func(p core.PulseNumber) (r []core.Node, r1 error)
	GetActiveNodesCounter    uint64
	GetActiveNodesPreCounter uint64
	GetActiveNodesMock       mActiveNodesStorageMockGetActiveNodes

	GetActiveNodesByRoleFunc       func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error)
	GetActiveNodesByRoleCounter    uint64
	GetActiveNodesByRolePreCounter uint64
	GetActiveNodesByRoleMock       mActiveNodesStorageMockGetActiveNodesByRole

	RemoveActiveNodesUntilFunc       func(p core.PulseNumber)
	RemoveActiveNodesUntilCounter    uint64
	RemoveActiveNodesUntilPreCounter uint64
	RemoveActiveNodesUntilMock       mActiveNodesStorageMockRemoveActiveNodesUntil

	SetActiveNodesFunc       func(p core.PulseNumber, p1 []core.Node) (r error)
	SetActiveNodesCounter    uint64
	SetActiveNodesPreCounter uint64
	SetActiveNodesMock       mActiveNodesStorageMockSetActiveNodes
}

//NewActiveNodesStorageMock returns a mock for github.com/insolar/insolar/ledger/storage.ActiveNodesStorage
func NewActiveNodesStorageMock(t minimock.Tester) *ActiveNodesStorageMock {
	m := &ActiveNodesStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetActiveNodesMock = mActiveNodesStorageMockGetActiveNodes{mock: m}
	m.GetActiveNodesByRoleMock = mActiveNodesStorageMockGetActiveNodesByRole{mock: m}
	m.RemoveActiveNodesUntilMock = mActiveNodesStorageMockRemoveActiveNodesUntil{mock: m}
	m.SetActiveNodesMock = mActiveNodesStorageMockSetActiveNodes{mock: m}

	return m
}

type mActiveNodesStorageMockGetActiveNodes struct {
	mock              *ActiveNodesStorageMock
	mainExpectation   *ActiveNodesStorageMockGetActiveNodesExpectation
	expectationSeries []*ActiveNodesStorageMockGetActiveNodesExpectation
}

type ActiveNodesStorageMockGetActiveNodesExpectation struct {
	input  *ActiveNodesStorageMockGetActiveNodesInput
	result *ActiveNodesStorageMockGetActiveNodesResult
}

type ActiveNodesStorageMockGetActiveNodesInput struct {
	p core.PulseNumber
}

type ActiveNodesStorageMockGetActiveNodesResult struct {
	r  []core.Node
	r1 error
}

//Expect specifies that invocation of ActiveNodesStorage.GetActiveNodes is expected from 1 to Infinity times
func (m *mActiveNodesStorageMockGetActiveNodes) Expect(p core.PulseNumber) *mActiveNodesStorageMockGetActiveNodes {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.input = &ActiveNodesStorageMockGetActiveNodesInput{p}
	return m
}

//Return specifies results of invocation of ActiveNodesStorage.GetActiveNodes
func (m *mActiveNodesStorageMockGetActiveNodes) Return(r []core.Node, r1 error) *ActiveNodesStorageMock {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &ActiveNodesStorageMockGetActiveNodesResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNodesStorage.GetActiveNodes is expected once
func (m *mActiveNodesStorageMockGetActiveNodes) ExpectOnce(p core.PulseNumber) *ActiveNodesStorageMockGetActiveNodesExpectation {
	m.mock.GetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodesStorageMockGetActiveNodesExpectation{}
	expectation.input = &ActiveNodesStorageMockGetActiveNodesInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodesStorageMockGetActiveNodesExpectation) Return(r []core.Node, r1 error) {
	e.result = &ActiveNodesStorageMockGetActiveNodesResult{r, r1}
}

//Set uses given function f as a mock of ActiveNodesStorage.GetActiveNodes method
func (m *mActiveNodesStorageMockGetActiveNodes) Set(f func(p core.PulseNumber) (r []core.Node, r1 error)) *ActiveNodesStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesFunc = f
	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/ledger/storage.ActiveNodesStorage interface
func (m *ActiveNodesStorageMock) GetActiveNodes(p core.PulseNumber) (r []core.Node, r1 error) {
	counter := atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.GetActiveNodes. %v", p)
			return
		}

		input := m.GetActiveNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ActiveNodesStorageMockGetActiveNodesInput{p}, "ActiveNodesStorage.GetActiveNodes got unexpected parameters")

		result := m.GetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodesStorageMock.GetActiveNodes")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesMock.mainExpectation != nil {

		input := m.GetActiveNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ActiveNodesStorageMockGetActiveNodesInput{p}, "ActiveNodesStorage.GetActiveNodes got unexpected parameters")
		}

		result := m.GetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodesStorageMock.GetActiveNodes")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.GetActiveNodes. %v", p)
		return
	}

	return m.GetActiveNodesFunc(p)
}

//GetActiveNodesMinimockCounter returns a count of ActiveNodesStorageMock.GetActiveNodesFunc invocations
func (m *ActiveNodesStorageMock) GetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesCounter)
}

//GetActiveNodesMinimockPreCounter returns the value of ActiveNodesStorageMock.GetActiveNodes invocations
func (m *ActiveNodesStorageMock) GetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesPreCounter)
}

//GetActiveNodesFinished returns true if mock invocations count is ok
func (m *ActiveNodesStorageMock) GetActiveNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodesCounter) == uint64(len(m.GetActiveNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodesFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodesCounter) > 0
	}

	return true
}

type mActiveNodesStorageMockGetActiveNodesByRole struct {
	mock              *ActiveNodesStorageMock
	mainExpectation   *ActiveNodesStorageMockGetActiveNodesByRoleExpectation
	expectationSeries []*ActiveNodesStorageMockGetActiveNodesByRoleExpectation
}

type ActiveNodesStorageMockGetActiveNodesByRoleExpectation struct {
	input  *ActiveNodesStorageMockGetActiveNodesByRoleInput
	result *ActiveNodesStorageMockGetActiveNodesByRoleResult
}

type ActiveNodesStorageMockGetActiveNodesByRoleInput struct {
	p  core.PulseNumber
	p1 core.StaticRole
}

type ActiveNodesStorageMockGetActiveNodesByRoleResult struct {
	r  []core.Node
	r1 error
}

//Expect specifies that invocation of ActiveNodesStorage.GetActiveNodesByRole is expected from 1 to Infinity times
func (m *mActiveNodesStorageMockGetActiveNodesByRole) Expect(p core.PulseNumber, p1 core.StaticRole) *mActiveNodesStorageMockGetActiveNodesByRole {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockGetActiveNodesByRoleExpectation{}
	}
	m.mainExpectation.input = &ActiveNodesStorageMockGetActiveNodesByRoleInput{p, p1}
	return m
}

//Return specifies results of invocation of ActiveNodesStorage.GetActiveNodesByRole
func (m *mActiveNodesStorageMockGetActiveNodesByRole) Return(r []core.Node, r1 error) *ActiveNodesStorageMock {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockGetActiveNodesByRoleExpectation{}
	}
	m.mainExpectation.result = &ActiveNodesStorageMockGetActiveNodesByRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNodesStorage.GetActiveNodesByRole is expected once
func (m *mActiveNodesStorageMockGetActiveNodesByRole) ExpectOnce(p core.PulseNumber, p1 core.StaticRole) *ActiveNodesStorageMockGetActiveNodesByRoleExpectation {
	m.mock.GetActiveNodesByRoleFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodesStorageMockGetActiveNodesByRoleExpectation{}
	expectation.input = &ActiveNodesStorageMockGetActiveNodesByRoleInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodesStorageMockGetActiveNodesByRoleExpectation) Return(r []core.Node, r1 error) {
	e.result = &ActiveNodesStorageMockGetActiveNodesByRoleResult{r, r1}
}

//Set uses given function f as a mock of ActiveNodesStorage.GetActiveNodesByRole method
func (m *mActiveNodesStorageMockGetActiveNodesByRole) Set(f func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error)) *ActiveNodesStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesByRoleFunc = f
	return m.mock
}

//GetActiveNodesByRole implements github.com/insolar/insolar/ledger/storage.ActiveNodesStorage interface
func (m *ActiveNodesStorageMock) GetActiveNodesByRole(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error) {
	counter := atomic.AddUint64(&m.GetActiveNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesByRoleCounter, 1)

	if len(m.GetActiveNodesByRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesByRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.GetActiveNodesByRole. %v %v", p, p1)
			return
		}

		input := m.GetActiveNodesByRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ActiveNodesStorageMockGetActiveNodesByRoleInput{p, p1}, "ActiveNodesStorage.GetActiveNodesByRole got unexpected parameters")

		result := m.GetActiveNodesByRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodesStorageMock.GetActiveNodesByRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesByRoleMock.mainExpectation != nil {

		input := m.GetActiveNodesByRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ActiveNodesStorageMockGetActiveNodesByRoleInput{p, p1}, "ActiveNodesStorage.GetActiveNodesByRole got unexpected parameters")
		}

		result := m.GetActiveNodesByRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodesStorageMock.GetActiveNodesByRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetActiveNodesByRoleFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.GetActiveNodesByRole. %v %v", p, p1)
		return
	}

	return m.GetActiveNodesByRoleFunc(p, p1)
}

//GetActiveNodesByRoleMinimockCounter returns a count of ActiveNodesStorageMock.GetActiveNodesByRoleFunc invocations
func (m *ActiveNodesStorageMock) GetActiveNodesByRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter)
}

//GetActiveNodesByRoleMinimockPreCounter returns the value of ActiveNodesStorageMock.GetActiveNodesByRole invocations
func (m *ActiveNodesStorageMock) GetActiveNodesByRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesByRolePreCounter)
}

//GetActiveNodesByRoleFinished returns true if mock invocations count is ok
func (m *ActiveNodesStorageMock) GetActiveNodesByRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodesByRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) == uint64(len(m.GetActiveNodesByRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodesByRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodesByRoleFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodesByRoleCounter) > 0
	}

	return true
}

type mActiveNodesStorageMockRemoveActiveNodesUntil struct {
	mock              *ActiveNodesStorageMock
	mainExpectation   *ActiveNodesStorageMockRemoveActiveNodesUntilExpectation
	expectationSeries []*ActiveNodesStorageMockRemoveActiveNodesUntilExpectation
}

type ActiveNodesStorageMockRemoveActiveNodesUntilExpectation struct {
	input *ActiveNodesStorageMockRemoveActiveNodesUntilInput
}

type ActiveNodesStorageMockRemoveActiveNodesUntilInput struct {
	p core.PulseNumber
}

//Expect specifies that invocation of ActiveNodesStorage.RemoveActiveNodesUntil is expected from 1 to Infinity times
func (m *mActiveNodesStorageMockRemoveActiveNodesUntil) Expect(p core.PulseNumber) *mActiveNodesStorageMockRemoveActiveNodesUntil {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockRemoveActiveNodesUntilExpectation{}
	}
	m.mainExpectation.input = &ActiveNodesStorageMockRemoveActiveNodesUntilInput{p}
	return m
}

//Return specifies results of invocation of ActiveNodesStorage.RemoveActiveNodesUntil
func (m *mActiveNodesStorageMockRemoveActiveNodesUntil) Return() *ActiveNodesStorageMock {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockRemoveActiveNodesUntilExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNodesStorage.RemoveActiveNodesUntil is expected once
func (m *mActiveNodesStorageMockRemoveActiveNodesUntil) ExpectOnce(p core.PulseNumber) *ActiveNodesStorageMockRemoveActiveNodesUntilExpectation {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodesStorageMockRemoveActiveNodesUntilExpectation{}
	expectation.input = &ActiveNodesStorageMockRemoveActiveNodesUntilInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ActiveNodesStorage.RemoveActiveNodesUntil method
func (m *mActiveNodesStorageMockRemoveActiveNodesUntil) Set(f func(p core.PulseNumber)) *ActiveNodesStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveActiveNodesUntilFunc = f
	return m.mock
}

//RemoveActiveNodesUntil implements github.com/insolar/insolar/ledger/storage.ActiveNodesStorage interface
func (m *ActiveNodesStorageMock) RemoveActiveNodesUntil(p core.PulseNumber) {
	counter := atomic.AddUint64(&m.RemoveActiveNodesUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveActiveNodesUntilCounter, 1)

	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.RemoveActiveNodesUntil. %v", p)
			return
		}

		input := m.RemoveActiveNodesUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ActiveNodesStorageMockRemoveActiveNodesUntilInput{p}, "ActiveNodesStorage.RemoveActiveNodesUntil got unexpected parameters")

		return
	}

	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {

		input := m.RemoveActiveNodesUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ActiveNodesStorageMockRemoveActiveNodesUntilInput{p}, "ActiveNodesStorage.RemoveActiveNodesUntil got unexpected parameters")
		}

		return
	}

	if m.RemoveActiveNodesUntilFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.RemoveActiveNodesUntil. %v", p)
		return
	}

	m.RemoveActiveNodesUntilFunc(p)
}

//RemoveActiveNodesUntilMinimockCounter returns a count of ActiveNodesStorageMock.RemoveActiveNodesUntilFunc invocations
func (m *ActiveNodesStorageMock) RemoveActiveNodesUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter)
}

//RemoveActiveNodesUntilMinimockPreCounter returns the value of ActiveNodesStorageMock.RemoveActiveNodesUntil invocations
func (m *ActiveNodesStorageMock) RemoveActiveNodesUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilPreCounter)
}

//RemoveActiveNodesUntilFinished returns true if mock invocations count is ok
func (m *ActiveNodesStorageMock) RemoveActiveNodesUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) == uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	return true
}

type mActiveNodesStorageMockSetActiveNodes struct {
	mock              *ActiveNodesStorageMock
	mainExpectation   *ActiveNodesStorageMockSetActiveNodesExpectation
	expectationSeries []*ActiveNodesStorageMockSetActiveNodesExpectation
}

type ActiveNodesStorageMockSetActiveNodesExpectation struct {
	input  *ActiveNodesStorageMockSetActiveNodesInput
	result *ActiveNodesStorageMockSetActiveNodesResult
}

type ActiveNodesStorageMockSetActiveNodesInput struct {
	p  core.PulseNumber
	p1 []core.Node
}

type ActiveNodesStorageMockSetActiveNodesResult struct {
	r error
}

//Expect specifies that invocation of ActiveNodesStorage.SetActiveNodes is expected from 1 to Infinity times
func (m *mActiveNodesStorageMockSetActiveNodes) Expect(p core.PulseNumber, p1 []core.Node) *mActiveNodesStorageMockSetActiveNodes {
	m.mock.SetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockSetActiveNodesExpectation{}
	}
	m.mainExpectation.input = &ActiveNodesStorageMockSetActiveNodesInput{p, p1}
	return m
}

//Return specifies results of invocation of ActiveNodesStorage.SetActiveNodes
func (m *mActiveNodesStorageMockSetActiveNodes) Return(r error) *ActiveNodesStorageMock {
	m.mock.SetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ActiveNodesStorageMockSetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &ActiveNodesStorageMockSetActiveNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ActiveNodesStorage.SetActiveNodes is expected once
func (m *mActiveNodesStorageMockSetActiveNodes) ExpectOnce(p core.PulseNumber, p1 []core.Node) *ActiveNodesStorageMockSetActiveNodesExpectation {
	m.mock.SetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &ActiveNodesStorageMockSetActiveNodesExpectation{}
	expectation.input = &ActiveNodesStorageMockSetActiveNodesInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ActiveNodesStorageMockSetActiveNodesExpectation) Return(r error) {
	e.result = &ActiveNodesStorageMockSetActiveNodesResult{r}
}

//Set uses given function f as a mock of ActiveNodesStorage.SetActiveNodes method
func (m *mActiveNodesStorageMockSetActiveNodes) Set(f func(p core.PulseNumber, p1 []core.Node) (r error)) *ActiveNodesStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetActiveNodesFunc = f
	return m.mock
}

//SetActiveNodes implements github.com/insolar/insolar/ledger/storage.ActiveNodesStorage interface
func (m *ActiveNodesStorageMock) SetActiveNodes(p core.PulseNumber, p1 []core.Node) (r error) {
	counter := atomic.AddUint64(&m.SetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.SetActiveNodesCounter, 1)

	if len(m.SetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.SetActiveNodes. %v %v", p, p1)
			return
		}

		input := m.SetActiveNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ActiveNodesStorageMockSetActiveNodesInput{p, p1}, "ActiveNodesStorage.SetActiveNodes got unexpected parameters")

		result := m.SetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodesStorageMock.SetActiveNodes")
			return
		}

		r = result.r

		return
	}

	if m.SetActiveNodesMock.mainExpectation != nil {

		input := m.SetActiveNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ActiveNodesStorageMockSetActiveNodesInput{p, p1}, "ActiveNodesStorage.SetActiveNodes got unexpected parameters")
		}

		result := m.SetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ActiveNodesStorageMock.SetActiveNodes")
		}

		r = result.r

		return
	}

	if m.SetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to ActiveNodesStorageMock.SetActiveNodes. %v %v", p, p1)
		return
	}

	return m.SetActiveNodesFunc(p, p1)
}

//SetActiveNodesMinimockCounter returns a count of ActiveNodesStorageMock.SetActiveNodesFunc invocations
func (m *ActiveNodesStorageMock) SetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetActiveNodesCounter)
}

//SetActiveNodesMinimockPreCounter returns the value of ActiveNodesStorageMock.SetActiveNodes invocations
func (m *ActiveNodesStorageMock) SetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetActiveNodesPreCounter)
}

//SetActiveNodesFinished returns true if mock invocations count is ok
func (m *ActiveNodesStorageMock) SetActiveNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetActiveNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetActiveNodesCounter) == uint64(len(m.SetActiveNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetActiveNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetActiveNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetActiveNodesFunc != nil {
		return atomic.LoadUint64(&m.SetActiveNodesCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ActiveNodesStorageMock) ValidateCallCounters() {

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.GetActiveNodes")
	}

	if !m.GetActiveNodesByRoleFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.GetActiveNodesByRole")
	}

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.RemoveActiveNodesUntil")
	}

	if !m.SetActiveNodesFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.SetActiveNodes")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ActiveNodesStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ActiveNodesStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ActiveNodesStorageMock) MinimockFinish() {

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.GetActiveNodes")
	}

	if !m.GetActiveNodesByRoleFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.GetActiveNodesByRole")
	}

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.RemoveActiveNodesUntil")
	}

	if !m.SetActiveNodesFinished() {
		m.t.Fatal("Expected call to ActiveNodesStorageMock.SetActiveNodes")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ActiveNodesStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ActiveNodesStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetActiveNodesFinished()
		ok = ok && m.GetActiveNodesByRoleFinished()
		ok = ok && m.RemoveActiveNodesUntilFinished()
		ok = ok && m.SetActiveNodesFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetActiveNodesFinished() {
				m.t.Error("Expected call to ActiveNodesStorageMock.GetActiveNodes")
			}

			if !m.GetActiveNodesByRoleFinished() {
				m.t.Error("Expected call to ActiveNodesStorageMock.GetActiveNodesByRole")
			}

			if !m.RemoveActiveNodesUntilFinished() {
				m.t.Error("Expected call to ActiveNodesStorageMock.RemoveActiveNodesUntil")
			}

			if !m.SetActiveNodesFinished() {
				m.t.Error("Expected call to ActiveNodesStorageMock.SetActiveNodes")
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
func (m *ActiveNodesStorageMock) AllMocksCalled() bool {

	if !m.GetActiveNodesFinished() {
		return false
	}

	if !m.GetActiveNodesByRoleFinished() {
		return false
	}

	if !m.RemoveActiveNodesUntilFinished() {
		return false
	}

	if !m.SetActiveNodesFinished() {
		return false
	}

	return true
}
