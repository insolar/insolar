package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Accessor" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//AccessorMock implements github.com/insolar/insolar/network.Accessor
type AccessorMock struct {
	t minimock.Tester

	GetActiveNodeFunc       func(p insolar.Reference) (r insolar.NetworkNode)
	GetActiveNodeCounter    uint64
	GetActiveNodePreCounter uint64
	GetActiveNodeMock       mAccessorMockGetActiveNode

	GetActiveNodeByShortIDFunc       func(p insolar.ShortNodeID) (r insolar.NetworkNode)
	GetActiveNodeByShortIDCounter    uint64
	GetActiveNodeByShortIDPreCounter uint64
	GetActiveNodeByShortIDMock       mAccessorMockGetActiveNodeByShortID

	GetActiveNodesFunc       func() (r []insolar.NetworkNode)
	GetActiveNodesCounter    uint64
	GetActiveNodesPreCounter uint64
	GetActiveNodesMock       mAccessorMockGetActiveNodes

	GetWorkingNodeFunc       func(p insolar.Reference) (r insolar.NetworkNode)
	GetWorkingNodeCounter    uint64
	GetWorkingNodePreCounter uint64
	GetWorkingNodeMock       mAccessorMockGetWorkingNode

	GetWorkingNodesFunc       func() (r []insolar.NetworkNode)
	GetWorkingNodesCounter    uint64
	GetWorkingNodesPreCounter uint64
	GetWorkingNodesMock       mAccessorMockGetWorkingNodes

	GetWorkingNodesByRoleFunc       func(p insolar.DynamicRole) (r []insolar.Reference)
	GetWorkingNodesByRoleCounter    uint64
	GetWorkingNodesByRolePreCounter uint64
	GetWorkingNodesByRoleMock       mAccessorMockGetWorkingNodesByRole
}

//NewAccessorMock returns a mock for github.com/insolar/insolar/network.Accessor
func NewAccessorMock(t minimock.Tester) *AccessorMock {
	m := &AccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetActiveNodeMock = mAccessorMockGetActiveNode{mock: m}
	m.GetActiveNodeByShortIDMock = mAccessorMockGetActiveNodeByShortID{mock: m}
	m.GetActiveNodesMock = mAccessorMockGetActiveNodes{mock: m}
	m.GetWorkingNodeMock = mAccessorMockGetWorkingNode{mock: m}
	m.GetWorkingNodesMock = mAccessorMockGetWorkingNodes{mock: m}
	m.GetWorkingNodesByRoleMock = mAccessorMockGetWorkingNodesByRole{mock: m}

	return m
}

type mAccessorMockGetActiveNode struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockGetActiveNodeExpectation
	expectationSeries []*AccessorMockGetActiveNodeExpectation
}

type AccessorMockGetActiveNodeExpectation struct {
	input  *AccessorMockGetActiveNodeInput
	result *AccessorMockGetActiveNodeResult
}

type AccessorMockGetActiveNodeInput struct {
	p insolar.Reference
}

type AccessorMockGetActiveNodeResult struct {
	r insolar.NetworkNode
}

//Expect specifies that invocation of Accessor.GetActiveNode is expected from 1 to Infinity times
func (m *mAccessorMockGetActiveNode) Expect(p insolar.Reference) *mAccessorMockGetActiveNode {
	m.mock.GetActiveNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetActiveNodeExpectation{}
	}
	m.mainExpectation.input = &AccessorMockGetActiveNodeInput{p}
	return m
}

//Return specifies results of invocation of Accessor.GetActiveNode
func (m *mAccessorMockGetActiveNode) Return(r insolar.NetworkNode) *AccessorMock {
	m.mock.GetActiveNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetActiveNodeExpectation{}
	}
	m.mainExpectation.result = &AccessorMockGetActiveNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.GetActiveNode is expected once
func (m *mAccessorMockGetActiveNode) ExpectOnce(p insolar.Reference) *AccessorMockGetActiveNodeExpectation {
	m.mock.GetActiveNodeFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockGetActiveNodeExpectation{}
	expectation.input = &AccessorMockGetActiveNodeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockGetActiveNodeExpectation) Return(r insolar.NetworkNode) {
	e.result = &AccessorMockGetActiveNodeResult{r}
}

//Set uses given function f as a mock of Accessor.GetActiveNode method
func (m *mAccessorMockGetActiveNode) Set(f func(p insolar.Reference) (r insolar.NetworkNode)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodeFunc = f
	return m.mock
}

//GetActiveNode implements github.com/insolar/insolar/network.Accessor interface
func (m *AccessorMock) GetActiveNode(p insolar.Reference) (r insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetActiveNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeCounter, 1)

	if len(m.GetActiveNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.GetActiveNode. %v", p)
			return
		}

		input := m.GetActiveNodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockGetActiveNodeInput{p}, "Accessor.GetActiveNode got unexpected parameters")

		result := m.GetActiveNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetActiveNode")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodeMock.mainExpectation != nil {

		input := m.GetActiveNodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockGetActiveNodeInput{p}, "Accessor.GetActiveNode got unexpected parameters")
		}

		result := m.GetActiveNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetActiveNode")
		}

		r = result.r

		return
	}

	if m.GetActiveNodeFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.GetActiveNode. %v", p)
		return
	}

	return m.GetActiveNodeFunc(p)
}

//GetActiveNodeMinimockCounter returns a count of AccessorMock.GetActiveNodeFunc invocations
func (m *AccessorMock) GetActiveNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeCounter)
}

//GetActiveNodeMinimockPreCounter returns the value of AccessorMock.GetActiveNode invocations
func (m *AccessorMock) GetActiveNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodePreCounter)
}

//GetActiveNodeFinished returns true if mock invocations count is ok
func (m *AccessorMock) GetActiveNodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodeCounter) == uint64(len(m.GetActiveNodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodeFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodeCounter) > 0
	}

	return true
}

type mAccessorMockGetActiveNodeByShortID struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockGetActiveNodeByShortIDExpectation
	expectationSeries []*AccessorMockGetActiveNodeByShortIDExpectation
}

type AccessorMockGetActiveNodeByShortIDExpectation struct {
	input  *AccessorMockGetActiveNodeByShortIDInput
	result *AccessorMockGetActiveNodeByShortIDResult
}

type AccessorMockGetActiveNodeByShortIDInput struct {
	p insolar.ShortNodeID
}

type AccessorMockGetActiveNodeByShortIDResult struct {
	r insolar.NetworkNode
}

//Expect specifies that invocation of Accessor.GetActiveNodeByShortID is expected from 1 to Infinity times
func (m *mAccessorMockGetActiveNodeByShortID) Expect(p insolar.ShortNodeID) *mAccessorMockGetActiveNodeByShortID {
	m.mock.GetActiveNodeByShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetActiveNodeByShortIDExpectation{}
	}
	m.mainExpectation.input = &AccessorMockGetActiveNodeByShortIDInput{p}
	return m
}

//Return specifies results of invocation of Accessor.GetActiveNodeByShortID
func (m *mAccessorMockGetActiveNodeByShortID) Return(r insolar.NetworkNode) *AccessorMock {
	m.mock.GetActiveNodeByShortIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetActiveNodeByShortIDExpectation{}
	}
	m.mainExpectation.result = &AccessorMockGetActiveNodeByShortIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.GetActiveNodeByShortID is expected once
func (m *mAccessorMockGetActiveNodeByShortID) ExpectOnce(p insolar.ShortNodeID) *AccessorMockGetActiveNodeByShortIDExpectation {
	m.mock.GetActiveNodeByShortIDFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockGetActiveNodeByShortIDExpectation{}
	expectation.input = &AccessorMockGetActiveNodeByShortIDInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockGetActiveNodeByShortIDExpectation) Return(r insolar.NetworkNode) {
	e.result = &AccessorMockGetActiveNodeByShortIDResult{r}
}

//Set uses given function f as a mock of Accessor.GetActiveNodeByShortID method
func (m *mAccessorMockGetActiveNodeByShortID) Set(f func(p insolar.ShortNodeID) (r insolar.NetworkNode)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodeByShortIDFunc = f
	return m.mock
}

//GetActiveNodeByShortID implements github.com/insolar/insolar/network.Accessor interface
func (m *AccessorMock) GetActiveNodeByShortID(p insolar.ShortNodeID) (r insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetActiveNodeByShortIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeByShortIDCounter, 1)

	if len(m.GetActiveNodeByShortIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodeByShortIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.GetActiveNodeByShortID. %v", p)
			return
		}

		input := m.GetActiveNodeByShortIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockGetActiveNodeByShortIDInput{p}, "Accessor.GetActiveNodeByShortID got unexpected parameters")

		result := m.GetActiveNodeByShortIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetActiveNodeByShortID")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodeByShortIDMock.mainExpectation != nil {

		input := m.GetActiveNodeByShortIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockGetActiveNodeByShortIDInput{p}, "Accessor.GetActiveNodeByShortID got unexpected parameters")
		}

		result := m.GetActiveNodeByShortIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetActiveNodeByShortID")
		}

		r = result.r

		return
	}

	if m.GetActiveNodeByShortIDFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.GetActiveNodeByShortID. %v", p)
		return
	}

	return m.GetActiveNodeByShortIDFunc(p)
}

//GetActiveNodeByShortIDMinimockCounter returns a count of AccessorMock.GetActiveNodeByShortIDFunc invocations
func (m *AccessorMock) GetActiveNodeByShortIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter)
}

//GetActiveNodeByShortIDMinimockPreCounter returns the value of AccessorMock.GetActiveNodeByShortID invocations
func (m *AccessorMock) GetActiveNodeByShortIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeByShortIDPreCounter)
}

//GetActiveNodeByShortIDFinished returns true if mock invocations count is ok
func (m *AccessorMock) GetActiveNodeByShortIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetActiveNodeByShortIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) == uint64(len(m.GetActiveNodeByShortIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetActiveNodeByShortIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetActiveNodeByShortIDFunc != nil {
		return atomic.LoadUint64(&m.GetActiveNodeByShortIDCounter) > 0
	}

	return true
}

type mAccessorMockGetActiveNodes struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockGetActiveNodesExpectation
	expectationSeries []*AccessorMockGetActiveNodesExpectation
}

type AccessorMockGetActiveNodesExpectation struct {
	result *AccessorMockGetActiveNodesResult
}

type AccessorMockGetActiveNodesResult struct {
	r []insolar.NetworkNode
}

//Expect specifies that invocation of Accessor.GetActiveNodes is expected from 1 to Infinity times
func (m *mAccessorMockGetActiveNodes) Expect() *mAccessorMockGetActiveNodes {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetActiveNodesExpectation{}
	}

	return m
}

//Return specifies results of invocation of Accessor.GetActiveNodes
func (m *mAccessorMockGetActiveNodes) Return(r []insolar.NetworkNode) *AccessorMock {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &AccessorMockGetActiveNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.GetActiveNodes is expected once
func (m *mAccessorMockGetActiveNodes) ExpectOnce() *AccessorMockGetActiveNodesExpectation {
	m.mock.GetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockGetActiveNodesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockGetActiveNodesExpectation) Return(r []insolar.NetworkNode) {
	e.result = &AccessorMockGetActiveNodesResult{r}
}

//Set uses given function f as a mock of Accessor.GetActiveNodes method
func (m *mAccessorMockGetActiveNodes) Set(f func() (r []insolar.NetworkNode)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesFunc = f
	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/network.Accessor interface
func (m *AccessorMock) GetActiveNodes() (r []insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.GetActiveNodes.")
			return
		}

		result := m.GetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetActiveNodes")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodesMock.mainExpectation != nil {

		result := m.GetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetActiveNodes")
		}

		r = result.r

		return
	}

	if m.GetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.GetActiveNodes.")
		return
	}

	return m.GetActiveNodesFunc()
}

//GetActiveNodesMinimockCounter returns a count of AccessorMock.GetActiveNodesFunc invocations
func (m *AccessorMock) GetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesCounter)
}

//GetActiveNodesMinimockPreCounter returns the value of AccessorMock.GetActiveNodes invocations
func (m *AccessorMock) GetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesPreCounter)
}

//GetActiveNodesFinished returns true if mock invocations count is ok
func (m *AccessorMock) GetActiveNodesFinished() bool {
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

type mAccessorMockGetWorkingNode struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockGetWorkingNodeExpectation
	expectationSeries []*AccessorMockGetWorkingNodeExpectation
}

type AccessorMockGetWorkingNodeExpectation struct {
	input  *AccessorMockGetWorkingNodeInput
	result *AccessorMockGetWorkingNodeResult
}

type AccessorMockGetWorkingNodeInput struct {
	p insolar.Reference
}

type AccessorMockGetWorkingNodeResult struct {
	r insolar.NetworkNode
}

//Expect specifies that invocation of Accessor.GetWorkingNode is expected from 1 to Infinity times
func (m *mAccessorMockGetWorkingNode) Expect(p insolar.Reference) *mAccessorMockGetWorkingNode {
	m.mock.GetWorkingNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetWorkingNodeExpectation{}
	}
	m.mainExpectation.input = &AccessorMockGetWorkingNodeInput{p}
	return m
}

//Return specifies results of invocation of Accessor.GetWorkingNode
func (m *mAccessorMockGetWorkingNode) Return(r insolar.NetworkNode) *AccessorMock {
	m.mock.GetWorkingNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetWorkingNodeExpectation{}
	}
	m.mainExpectation.result = &AccessorMockGetWorkingNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.GetWorkingNode is expected once
func (m *mAccessorMockGetWorkingNode) ExpectOnce(p insolar.Reference) *AccessorMockGetWorkingNodeExpectation {
	m.mock.GetWorkingNodeFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockGetWorkingNodeExpectation{}
	expectation.input = &AccessorMockGetWorkingNodeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockGetWorkingNodeExpectation) Return(r insolar.NetworkNode) {
	e.result = &AccessorMockGetWorkingNodeResult{r}
}

//Set uses given function f as a mock of Accessor.GetWorkingNode method
func (m *mAccessorMockGetWorkingNode) Set(f func(p insolar.Reference) (r insolar.NetworkNode)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodeFunc = f
	return m.mock
}

//GetWorkingNode implements github.com/insolar/insolar/network.Accessor interface
func (m *AccessorMock) GetWorkingNode(p insolar.Reference) (r insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetWorkingNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodeCounter, 1)

	if len(m.GetWorkingNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.GetWorkingNode. %v", p)
			return
		}

		input := m.GetWorkingNodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockGetWorkingNodeInput{p}, "Accessor.GetWorkingNode got unexpected parameters")

		result := m.GetWorkingNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetWorkingNode")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodeMock.mainExpectation != nil {

		input := m.GetWorkingNodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockGetWorkingNodeInput{p}, "Accessor.GetWorkingNode got unexpected parameters")
		}

		result := m.GetWorkingNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetWorkingNode")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodeFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.GetWorkingNode. %v", p)
		return
	}

	return m.GetWorkingNodeFunc(p)
}

//GetWorkingNodeMinimockCounter returns a count of AccessorMock.GetWorkingNodeFunc invocations
func (m *AccessorMock) GetWorkingNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodeCounter)
}

//GetWorkingNodeMinimockPreCounter returns the value of AccessorMock.GetWorkingNode invocations
func (m *AccessorMock) GetWorkingNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodePreCounter)
}

//GetWorkingNodeFinished returns true if mock invocations count is ok
func (m *AccessorMock) GetWorkingNodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetWorkingNodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetWorkingNodeCounter) == uint64(len(m.GetWorkingNodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetWorkingNodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetWorkingNodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetWorkingNodeFunc != nil {
		return atomic.LoadUint64(&m.GetWorkingNodeCounter) > 0
	}

	return true
}

type mAccessorMockGetWorkingNodes struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockGetWorkingNodesExpectation
	expectationSeries []*AccessorMockGetWorkingNodesExpectation
}

type AccessorMockGetWorkingNodesExpectation struct {
	result *AccessorMockGetWorkingNodesResult
}

type AccessorMockGetWorkingNodesResult struct {
	r []insolar.NetworkNode
}

//Expect specifies that invocation of Accessor.GetWorkingNodes is expected from 1 to Infinity times
func (m *mAccessorMockGetWorkingNodes) Expect() *mAccessorMockGetWorkingNodes {
	m.mock.GetWorkingNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetWorkingNodesExpectation{}
	}

	return m
}

//Return specifies results of invocation of Accessor.GetWorkingNodes
func (m *mAccessorMockGetWorkingNodes) Return(r []insolar.NetworkNode) *AccessorMock {
	m.mock.GetWorkingNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetWorkingNodesExpectation{}
	}
	m.mainExpectation.result = &AccessorMockGetWorkingNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.GetWorkingNodes is expected once
func (m *mAccessorMockGetWorkingNodes) ExpectOnce() *AccessorMockGetWorkingNodesExpectation {
	m.mock.GetWorkingNodesFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockGetWorkingNodesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockGetWorkingNodesExpectation) Return(r []insolar.NetworkNode) {
	e.result = &AccessorMockGetWorkingNodesResult{r}
}

//Set uses given function f as a mock of Accessor.GetWorkingNodes method
func (m *mAccessorMockGetWorkingNodes) Set(f func() (r []insolar.NetworkNode)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodesFunc = f
	return m.mock
}

//GetWorkingNodes implements github.com/insolar/insolar/network.Accessor interface
func (m *AccessorMock) GetWorkingNodes() (r []insolar.NetworkNode) {
	counter := atomic.AddUint64(&m.GetWorkingNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodesCounter, 1)

	if len(m.GetWorkingNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.GetWorkingNodes.")
			return
		}

		result := m.GetWorkingNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetWorkingNodes")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesMock.mainExpectation != nil {

		result := m.GetWorkingNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetWorkingNodes")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.GetWorkingNodes.")
		return
	}

	return m.GetWorkingNodesFunc()
}

//GetWorkingNodesMinimockCounter returns a count of AccessorMock.GetWorkingNodesFunc invocations
func (m *AccessorMock) GetWorkingNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesCounter)
}

//GetWorkingNodesMinimockPreCounter returns the value of AccessorMock.GetWorkingNodes invocations
func (m *AccessorMock) GetWorkingNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesPreCounter)
}

//GetWorkingNodesFinished returns true if mock invocations count is ok
func (m *AccessorMock) GetWorkingNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetWorkingNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetWorkingNodesCounter) == uint64(len(m.GetWorkingNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetWorkingNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetWorkingNodesFunc != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesCounter) > 0
	}

	return true
}

type mAccessorMockGetWorkingNodesByRole struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockGetWorkingNodesByRoleExpectation
	expectationSeries []*AccessorMockGetWorkingNodesByRoleExpectation
}

type AccessorMockGetWorkingNodesByRoleExpectation struct {
	input  *AccessorMockGetWorkingNodesByRoleInput
	result *AccessorMockGetWorkingNodesByRoleResult
}

type AccessorMockGetWorkingNodesByRoleInput struct {
	p insolar.DynamicRole
}

type AccessorMockGetWorkingNodesByRoleResult struct {
	r []insolar.Reference
}

//Expect specifies that invocation of Accessor.GetWorkingNodesByRole is expected from 1 to Infinity times
func (m *mAccessorMockGetWorkingNodesByRole) Expect(p insolar.DynamicRole) *mAccessorMockGetWorkingNodesByRole {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetWorkingNodesByRoleExpectation{}
	}
	m.mainExpectation.input = &AccessorMockGetWorkingNodesByRoleInput{p}
	return m
}

//Return specifies results of invocation of Accessor.GetWorkingNodesByRole
func (m *mAccessorMockGetWorkingNodesByRole) Return(r []insolar.Reference) *AccessorMock {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockGetWorkingNodesByRoleExpectation{}
	}
	m.mainExpectation.result = &AccessorMockGetWorkingNodesByRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.GetWorkingNodesByRole is expected once
func (m *mAccessorMockGetWorkingNodesByRole) ExpectOnce(p insolar.DynamicRole) *AccessorMockGetWorkingNodesByRoleExpectation {
	m.mock.GetWorkingNodesByRoleFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockGetWorkingNodesByRoleExpectation{}
	expectation.input = &AccessorMockGetWorkingNodesByRoleInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockGetWorkingNodesByRoleExpectation) Return(r []insolar.Reference) {
	e.result = &AccessorMockGetWorkingNodesByRoleResult{r}
}

//Set uses given function f as a mock of Accessor.GetWorkingNodesByRole method
func (m *mAccessorMockGetWorkingNodesByRole) Set(f func(p insolar.DynamicRole) (r []insolar.Reference)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetWorkingNodesByRoleFunc = f
	return m.mock
}

//GetWorkingNodesByRole implements github.com/insolar/insolar/network.Accessor interface
func (m *AccessorMock) GetWorkingNodesByRole(p insolar.DynamicRole) (r []insolar.Reference) {
	counter := atomic.AddUint64(&m.GetWorkingNodesByRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetWorkingNodesByRoleCounter, 1)

	if len(m.GetWorkingNodesByRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetWorkingNodesByRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.GetWorkingNodesByRole. %v", p)
			return
		}

		input := m.GetWorkingNodesByRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockGetWorkingNodesByRoleInput{p}, "Accessor.GetWorkingNodesByRole got unexpected parameters")

		result := m.GetWorkingNodesByRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetWorkingNodesByRole")
			return
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesByRoleMock.mainExpectation != nil {

		input := m.GetWorkingNodesByRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockGetWorkingNodesByRoleInput{p}, "Accessor.GetWorkingNodesByRole got unexpected parameters")
		}

		result := m.GetWorkingNodesByRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.GetWorkingNodesByRole")
		}

		r = result.r

		return
	}

	if m.GetWorkingNodesByRoleFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.GetWorkingNodesByRole. %v", p)
		return
	}

	return m.GetWorkingNodesByRoleFunc(p)
}

//GetWorkingNodesByRoleMinimockCounter returns a count of AccessorMock.GetWorkingNodesByRoleFunc invocations
func (m *AccessorMock) GetWorkingNodesByRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter)
}

//GetWorkingNodesByRoleMinimockPreCounter returns the value of AccessorMock.GetWorkingNodesByRole invocations
func (m *AccessorMock) GetWorkingNodesByRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetWorkingNodesByRolePreCounter)
}

//GetWorkingNodesByRoleFinished returns true if mock invocations count is ok
func (m *AccessorMock) GetWorkingNodesByRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetWorkingNodesByRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter) == uint64(len(m.GetWorkingNodesByRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetWorkingNodesByRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetWorkingNodesByRoleFunc != nil {
		return atomic.LoadUint64(&m.GetWorkingNodesByRoleCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) ValidateCallCounters() {

	if !m.GetActiveNodeFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetActiveNode")
	}

	if !m.GetActiveNodeByShortIDFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetActiveNodeByShortID")
	}

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetActiveNodes")
	}

	if !m.GetWorkingNodeFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetWorkingNode")
	}

	if !m.GetWorkingNodesFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetWorkingNodes")
	}

	if !m.GetWorkingNodesByRoleFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetWorkingNodesByRole")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AccessorMock) MinimockFinish() {

	if !m.GetActiveNodeFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetActiveNode")
	}

	if !m.GetActiveNodeByShortIDFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetActiveNodeByShortID")
	}

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetActiveNodes")
	}

	if !m.GetWorkingNodeFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetWorkingNode")
	}

	if !m.GetWorkingNodesFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetWorkingNodes")
	}

	if !m.GetWorkingNodesByRoleFinished() {
		m.t.Fatal("Expected call to AccessorMock.GetWorkingNodesByRole")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetActiveNodeFinished()
		ok = ok && m.GetActiveNodeByShortIDFinished()
		ok = ok && m.GetActiveNodesFinished()
		ok = ok && m.GetWorkingNodeFinished()
		ok = ok && m.GetWorkingNodesFinished()
		ok = ok && m.GetWorkingNodesByRoleFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetActiveNodeFinished() {
				m.t.Error("Expected call to AccessorMock.GetActiveNode")
			}

			if !m.GetActiveNodeByShortIDFinished() {
				m.t.Error("Expected call to AccessorMock.GetActiveNodeByShortID")
			}

			if !m.GetActiveNodesFinished() {
				m.t.Error("Expected call to AccessorMock.GetActiveNodes")
			}

			if !m.GetWorkingNodeFinished() {
				m.t.Error("Expected call to AccessorMock.GetWorkingNode")
			}

			if !m.GetWorkingNodesFinished() {
				m.t.Error("Expected call to AccessorMock.GetWorkingNodes")
			}

			if !m.GetWorkingNodesByRoleFinished() {
				m.t.Error("Expected call to AccessorMock.GetWorkingNodesByRole")
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
func (m *AccessorMock) AllMocksCalled() bool {

	if !m.GetActiveNodeFinished() {
		return false
	}

	if !m.GetActiveNodeByShortIDFinished() {
		return false
	}

	if !m.GetActiveNodesFinished() {
		return false
	}

	if !m.GetWorkingNodeFinished() {
		return false
	}

	if !m.GetWorkingNodesFinished() {
		return false
	}

	if !m.GetWorkingNodesByRoleFinished() {
		return false
	}

	return true
}
