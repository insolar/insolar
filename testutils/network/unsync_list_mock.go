package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "UnsyncList" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	packets "github.com/insolar/insolar/consensus/packets"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//UnsyncListMock implements github.com/insolar/insolar/network.UnsyncList
type UnsyncListMock struct {
	t minimock.Tester

	AddClaimsFunc       func(p map[core.RecordRef][]packets.ReferendumClaim, p1 map[core.RecordRef]string)
	AddClaimsCounter    uint64
	AddClaimsPreCounter uint64
	AddClaimsMock       mUnsyncListMockAddClaims

	CalculateHashFunc       func() (r []byte, r1 error)
	CalculateHashCounter    uint64
	CalculateHashPreCounter uint64
	CalculateHashMock       mUnsyncListMockCalculateHash

	GetActiveNodeFunc       func(p core.RecordRef) (r core.Node)
	GetActiveNodeCounter    uint64
	GetActiveNodePreCounter uint64
	GetActiveNodeMock       mUnsyncListMockGetActiveNode

	GetActiveNodesFunc       func() (r []core.Node)
	GetActiveNodesCounter    uint64
	GetActiveNodesPreCounter uint64
	GetActiveNodesMock       mUnsyncListMockGetActiveNodes

	IndexToRefFunc       func(p int) (r core.RecordRef, r1 error)
	IndexToRefCounter    uint64
	IndexToRefPreCounter uint64
	IndexToRefMock       mUnsyncListMockIndexToRef

	LengthFunc       func() (r int)
	LengthCounter    uint64
	LengthPreCounter uint64
	LengthMock       mUnsyncListMockLength

	RefToIndexFunc       func(p core.RecordRef) (r int, r1 error)
	RefToIndexCounter    uint64
	RefToIndexPreCounter uint64
	RefToIndexMock       mUnsyncListMockRefToIndex

	RemoveClaimsFunc       func(p core.RecordRef)
	RemoveClaimsCounter    uint64
	RemoveClaimsPreCounter uint64
	RemoveClaimsMock       mUnsyncListMockRemoveClaims
}

//NewUnsyncListMock returns a mock for github.com/insolar/insolar/network.UnsyncList
func NewUnsyncListMock(t minimock.Tester) *UnsyncListMock {
	m := &UnsyncListMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddClaimsMock = mUnsyncListMockAddClaims{mock: m}
	m.CalculateHashMock = mUnsyncListMockCalculateHash{mock: m}
	m.GetActiveNodeMock = mUnsyncListMockGetActiveNode{mock: m}
	m.GetActiveNodesMock = mUnsyncListMockGetActiveNodes{mock: m}
	m.IndexToRefMock = mUnsyncListMockIndexToRef{mock: m}
	m.LengthMock = mUnsyncListMockLength{mock: m}
	m.RefToIndexMock = mUnsyncListMockRefToIndex{mock: m}
	m.RemoveClaimsMock = mUnsyncListMockRemoveClaims{mock: m}

	return m
}

type mUnsyncListMockAddClaims struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockAddClaimsExpectation
	expectationSeries []*UnsyncListMockAddClaimsExpectation
}

type UnsyncListMockAddClaimsExpectation struct {
	input *UnsyncListMockAddClaimsInput
}

type UnsyncListMockAddClaimsInput struct {
	p  map[core.RecordRef][]packets.ReferendumClaim
	p1 map[core.RecordRef]string
}

//Expect specifies that invocation of UnsyncList.AddClaims is expected from 1 to Infinity times
func (m *mUnsyncListMockAddClaims) Expect(p map[core.RecordRef][]packets.ReferendumClaim, p1 map[core.RecordRef]string) *mUnsyncListMockAddClaims {
	m.mock.AddClaimsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockAddClaimsExpectation{}
	}
	m.mainExpectation.input = &UnsyncListMockAddClaimsInput{p, p1}
	return m
}

//Return specifies results of invocation of UnsyncList.AddClaims
func (m *mUnsyncListMockAddClaims) Return() *UnsyncListMock {
	m.mock.AddClaimsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockAddClaimsExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.AddClaims is expected once
func (m *mUnsyncListMockAddClaims) ExpectOnce(p map[core.RecordRef][]packets.ReferendumClaim, p1 map[core.RecordRef]string) *UnsyncListMockAddClaimsExpectation {
	m.mock.AddClaimsFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockAddClaimsExpectation{}
	expectation.input = &UnsyncListMockAddClaimsInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of UnsyncList.AddClaims method
func (m *mUnsyncListMockAddClaims) Set(f func(p map[core.RecordRef][]packets.ReferendumClaim, p1 map[core.RecordRef]string)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddClaimsFunc = f
	return m.mock
}

//AddClaims implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) AddClaims(p map[core.RecordRef][]packets.ReferendumClaim, p1 map[core.RecordRef]string) {
	counter := atomic.AddUint64(&m.AddClaimsPreCounter, 1)
	defer atomic.AddUint64(&m.AddClaimsCounter, 1)

	if len(m.AddClaimsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddClaimsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.AddClaims. %v %v", p, p1)
			return
		}

		input := m.AddClaimsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, UnsyncListMockAddClaimsInput{p, p1}, "UnsyncList.AddClaims got unexpected parameters")

		return
	}

	if m.AddClaimsMock.mainExpectation != nil {

		input := m.AddClaimsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, UnsyncListMockAddClaimsInput{p, p1}, "UnsyncList.AddClaims got unexpected parameters")
		}

		return
	}

	if m.AddClaimsFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.AddClaims. %v %v", p, p1)
		return
	}

	m.AddClaimsFunc(p, p1)
}

//AddClaimsMinimockCounter returns a count of UnsyncListMock.AddClaimsFunc invocations
func (m *UnsyncListMock) AddClaimsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddClaimsCounter)
}

//AddClaimsMinimockPreCounter returns the value of UnsyncListMock.AddClaims invocations
func (m *UnsyncListMock) AddClaimsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddClaimsPreCounter)
}

//AddClaimsFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) AddClaimsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddClaimsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddClaimsCounter) == uint64(len(m.AddClaimsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddClaimsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddClaimsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddClaimsFunc != nil {
		return atomic.LoadUint64(&m.AddClaimsCounter) > 0
	}

	return true
}

type mUnsyncListMockCalculateHash struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockCalculateHashExpectation
	expectationSeries []*UnsyncListMockCalculateHashExpectation
}

type UnsyncListMockCalculateHashExpectation struct {
	result *UnsyncListMockCalculateHashResult
}

type UnsyncListMockCalculateHashResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of UnsyncList.CalculateHash is expected from 1 to Infinity times
func (m *mUnsyncListMockCalculateHash) Expect() *mUnsyncListMockCalculateHash {
	m.mock.CalculateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockCalculateHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of UnsyncList.CalculateHash
func (m *mUnsyncListMockCalculateHash) Return(r []byte, r1 error) *UnsyncListMock {
	m.mock.CalculateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockCalculateHashExpectation{}
	}
	m.mainExpectation.result = &UnsyncListMockCalculateHashResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.CalculateHash is expected once
func (m *mUnsyncListMockCalculateHash) ExpectOnce() *UnsyncListMockCalculateHashExpectation {
	m.mock.CalculateHashFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockCalculateHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *UnsyncListMockCalculateHashExpectation) Return(r []byte, r1 error) {
	e.result = &UnsyncListMockCalculateHashResult{r, r1}
}

//Set uses given function f as a mock of UnsyncList.CalculateHash method
func (m *mUnsyncListMockCalculateHash) Set(f func() (r []byte, r1 error)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CalculateHashFunc = f
	return m.mock
}

//CalculateHash implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) CalculateHash() (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.CalculateHashPreCounter, 1)
	defer atomic.AddUint64(&m.CalculateHashCounter, 1)

	if len(m.CalculateHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CalculateHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.CalculateHash.")
			return
		}

		result := m.CalculateHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.CalculateHash")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CalculateHashMock.mainExpectation != nil {

		result := m.CalculateHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.CalculateHash")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CalculateHashFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.CalculateHash.")
		return
	}

	return m.CalculateHashFunc()
}

//CalculateHashMinimockCounter returns a count of UnsyncListMock.CalculateHashFunc invocations
func (m *UnsyncListMock) CalculateHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CalculateHashCounter)
}

//CalculateHashMinimockPreCounter returns the value of UnsyncListMock.CalculateHash invocations
func (m *UnsyncListMock) CalculateHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CalculateHashPreCounter)
}

//CalculateHashFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) CalculateHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CalculateHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CalculateHashCounter) == uint64(len(m.CalculateHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CalculateHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CalculateHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CalculateHashFunc != nil {
		return atomic.LoadUint64(&m.CalculateHashCounter) > 0
	}

	return true
}

type mUnsyncListMockGetActiveNode struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockGetActiveNodeExpectation
	expectationSeries []*UnsyncListMockGetActiveNodeExpectation
}

type UnsyncListMockGetActiveNodeExpectation struct {
	input  *UnsyncListMockGetActiveNodeInput
	result *UnsyncListMockGetActiveNodeResult
}

type UnsyncListMockGetActiveNodeInput struct {
	p core.RecordRef
}

type UnsyncListMockGetActiveNodeResult struct {
	r core.Node
}

//Expect specifies that invocation of UnsyncList.GetActiveNode is expected from 1 to Infinity times
func (m *mUnsyncListMockGetActiveNode) Expect(p core.RecordRef) *mUnsyncListMockGetActiveNode {
	m.mock.GetActiveNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockGetActiveNodeExpectation{}
	}
	m.mainExpectation.input = &UnsyncListMockGetActiveNodeInput{p}
	return m
}

//Return specifies results of invocation of UnsyncList.GetActiveNode
func (m *mUnsyncListMockGetActiveNode) Return(r core.Node) *UnsyncListMock {
	m.mock.GetActiveNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockGetActiveNodeExpectation{}
	}
	m.mainExpectation.result = &UnsyncListMockGetActiveNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.GetActiveNode is expected once
func (m *mUnsyncListMockGetActiveNode) ExpectOnce(p core.RecordRef) *UnsyncListMockGetActiveNodeExpectation {
	m.mock.GetActiveNodeFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockGetActiveNodeExpectation{}
	expectation.input = &UnsyncListMockGetActiveNodeInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *UnsyncListMockGetActiveNodeExpectation) Return(r core.Node) {
	e.result = &UnsyncListMockGetActiveNodeResult{r}
}

//Set uses given function f as a mock of UnsyncList.GetActiveNode method
func (m *mUnsyncListMockGetActiveNode) Set(f func(p core.RecordRef) (r core.Node)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodeFunc = f
	return m.mock
}

//GetActiveNode implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) GetActiveNode(p core.RecordRef) (r core.Node) {
	counter := atomic.AddUint64(&m.GetActiveNodePreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodeCounter, 1)

	if len(m.GetActiveNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.GetActiveNode. %v", p)
			return
		}

		input := m.GetActiveNodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, UnsyncListMockGetActiveNodeInput{p}, "UnsyncList.GetActiveNode got unexpected parameters")

		result := m.GetActiveNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.GetActiveNode")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodeMock.mainExpectation != nil {

		input := m.GetActiveNodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, UnsyncListMockGetActiveNodeInput{p}, "UnsyncList.GetActiveNode got unexpected parameters")
		}

		result := m.GetActiveNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.GetActiveNode")
		}

		r = result.r

		return
	}

	if m.GetActiveNodeFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.GetActiveNode. %v", p)
		return
	}

	return m.GetActiveNodeFunc(p)
}

//GetActiveNodeMinimockCounter returns a count of UnsyncListMock.GetActiveNodeFunc invocations
func (m *UnsyncListMock) GetActiveNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodeCounter)
}

//GetActiveNodeMinimockPreCounter returns the value of UnsyncListMock.GetActiveNode invocations
func (m *UnsyncListMock) GetActiveNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodePreCounter)
}

//GetActiveNodeFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) GetActiveNodeFinished() bool {
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

type mUnsyncListMockGetActiveNodes struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockGetActiveNodesExpectation
	expectationSeries []*UnsyncListMockGetActiveNodesExpectation
}

type UnsyncListMockGetActiveNodesExpectation struct {
	result *UnsyncListMockGetActiveNodesResult
}

type UnsyncListMockGetActiveNodesResult struct {
	r []core.Node
}

//Expect specifies that invocation of UnsyncList.GetActiveNodes is expected from 1 to Infinity times
func (m *mUnsyncListMockGetActiveNodes) Expect() *mUnsyncListMockGetActiveNodes {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockGetActiveNodesExpectation{}
	}

	return m
}

//Return specifies results of invocation of UnsyncList.GetActiveNodes
func (m *mUnsyncListMockGetActiveNodes) Return(r []core.Node) *UnsyncListMock {
	m.mock.GetActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockGetActiveNodesExpectation{}
	}
	m.mainExpectation.result = &UnsyncListMockGetActiveNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.GetActiveNodes is expected once
func (m *mUnsyncListMockGetActiveNodes) ExpectOnce() *UnsyncListMockGetActiveNodesExpectation {
	m.mock.GetActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockGetActiveNodesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *UnsyncListMockGetActiveNodesExpectation) Return(r []core.Node) {
	e.result = &UnsyncListMockGetActiveNodesResult{r}
}

//Set uses given function f as a mock of UnsyncList.GetActiveNodes method
func (m *mUnsyncListMockGetActiveNodes) Set(f func() (r []core.Node)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetActiveNodesFunc = f
	return m.mock
}

//GetActiveNodes implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) GetActiveNodes() (r []core.Node) {
	counter := atomic.AddUint64(&m.GetActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.GetActiveNodesCounter, 1)

	if len(m.GetActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.GetActiveNodes.")
			return
		}

		result := m.GetActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.GetActiveNodes")
			return
		}

		r = result.r

		return
	}

	if m.GetActiveNodesMock.mainExpectation != nil {

		result := m.GetActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.GetActiveNodes")
		}

		r = result.r

		return
	}

	if m.GetActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.GetActiveNodes.")
		return
	}

	return m.GetActiveNodesFunc()
}

//GetActiveNodesMinimockCounter returns a count of UnsyncListMock.GetActiveNodesFunc invocations
func (m *UnsyncListMock) GetActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesCounter)
}

//GetActiveNodesMinimockPreCounter returns the value of UnsyncListMock.GetActiveNodes invocations
func (m *UnsyncListMock) GetActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetActiveNodesPreCounter)
}

//GetActiveNodesFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) GetActiveNodesFinished() bool {
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

type mUnsyncListMockIndexToRef struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockIndexToRefExpectation
	expectationSeries []*UnsyncListMockIndexToRefExpectation
}

type UnsyncListMockIndexToRefExpectation struct {
	input  *UnsyncListMockIndexToRefInput
	result *UnsyncListMockIndexToRefResult
}

type UnsyncListMockIndexToRefInput struct {
	p int
}

type UnsyncListMockIndexToRefResult struct {
	r  core.RecordRef
	r1 error
}

//Expect specifies that invocation of UnsyncList.IndexToRef is expected from 1 to Infinity times
func (m *mUnsyncListMockIndexToRef) Expect(p int) *mUnsyncListMockIndexToRef {
	m.mock.IndexToRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockIndexToRefExpectation{}
	}
	m.mainExpectation.input = &UnsyncListMockIndexToRefInput{p}
	return m
}

//Return specifies results of invocation of UnsyncList.IndexToRef
func (m *mUnsyncListMockIndexToRef) Return(r core.RecordRef, r1 error) *UnsyncListMock {
	m.mock.IndexToRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockIndexToRefExpectation{}
	}
	m.mainExpectation.result = &UnsyncListMockIndexToRefResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.IndexToRef is expected once
func (m *mUnsyncListMockIndexToRef) ExpectOnce(p int) *UnsyncListMockIndexToRefExpectation {
	m.mock.IndexToRefFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockIndexToRefExpectation{}
	expectation.input = &UnsyncListMockIndexToRefInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *UnsyncListMockIndexToRefExpectation) Return(r core.RecordRef, r1 error) {
	e.result = &UnsyncListMockIndexToRefResult{r, r1}
}

//Set uses given function f as a mock of UnsyncList.IndexToRef method
func (m *mUnsyncListMockIndexToRef) Set(f func(p int) (r core.RecordRef, r1 error)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IndexToRefFunc = f
	return m.mock
}

//IndexToRef implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) IndexToRef(p int) (r core.RecordRef, r1 error) {
	counter := atomic.AddUint64(&m.IndexToRefPreCounter, 1)
	defer atomic.AddUint64(&m.IndexToRefCounter, 1)

	if len(m.IndexToRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IndexToRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.IndexToRef. %v", p)
			return
		}

		input := m.IndexToRefMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, UnsyncListMockIndexToRefInput{p}, "UnsyncList.IndexToRef got unexpected parameters")

		result := m.IndexToRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.IndexToRef")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IndexToRefMock.mainExpectation != nil {

		input := m.IndexToRefMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, UnsyncListMockIndexToRefInput{p}, "UnsyncList.IndexToRef got unexpected parameters")
		}

		result := m.IndexToRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.IndexToRef")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IndexToRefFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.IndexToRef. %v", p)
		return
	}

	return m.IndexToRefFunc(p)
}

//IndexToRefMinimockCounter returns a count of UnsyncListMock.IndexToRefFunc invocations
func (m *UnsyncListMock) IndexToRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IndexToRefCounter)
}

//IndexToRefMinimockPreCounter returns the value of UnsyncListMock.IndexToRef invocations
func (m *UnsyncListMock) IndexToRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IndexToRefPreCounter)
}

//IndexToRefFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) IndexToRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IndexToRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IndexToRefCounter) == uint64(len(m.IndexToRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IndexToRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IndexToRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IndexToRefFunc != nil {
		return atomic.LoadUint64(&m.IndexToRefCounter) > 0
	}

	return true
}

type mUnsyncListMockLength struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockLengthExpectation
	expectationSeries []*UnsyncListMockLengthExpectation
}

type UnsyncListMockLengthExpectation struct {
	result *UnsyncListMockLengthResult
}

type UnsyncListMockLengthResult struct {
	r int
}

//Expect specifies that invocation of UnsyncList.Length is expected from 1 to Infinity times
func (m *mUnsyncListMockLength) Expect() *mUnsyncListMockLength {
	m.mock.LengthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockLengthExpectation{}
	}

	return m
}

//Return specifies results of invocation of UnsyncList.Length
func (m *mUnsyncListMockLength) Return(r int) *UnsyncListMock {
	m.mock.LengthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockLengthExpectation{}
	}
	m.mainExpectation.result = &UnsyncListMockLengthResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.Length is expected once
func (m *mUnsyncListMockLength) ExpectOnce() *UnsyncListMockLengthExpectation {
	m.mock.LengthFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockLengthExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *UnsyncListMockLengthExpectation) Return(r int) {
	e.result = &UnsyncListMockLengthResult{r}
}

//Set uses given function f as a mock of UnsyncList.Length method
func (m *mUnsyncListMockLength) Set(f func() (r int)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LengthFunc = f
	return m.mock
}

//Length implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) Length() (r int) {
	counter := atomic.AddUint64(&m.LengthPreCounter, 1)
	defer atomic.AddUint64(&m.LengthCounter, 1)

	if len(m.LengthMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LengthMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.Length.")
			return
		}

		result := m.LengthMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.Length")
			return
		}

		r = result.r

		return
	}

	if m.LengthMock.mainExpectation != nil {

		result := m.LengthMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.Length")
		}

		r = result.r

		return
	}

	if m.LengthFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.Length.")
		return
	}

	return m.LengthFunc()
}

//LengthMinimockCounter returns a count of UnsyncListMock.LengthFunc invocations
func (m *UnsyncListMock) LengthMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LengthCounter)
}

//LengthMinimockPreCounter returns the value of UnsyncListMock.Length invocations
func (m *UnsyncListMock) LengthMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LengthPreCounter)
}

//LengthFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) LengthFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LengthMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LengthCounter) == uint64(len(m.LengthMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LengthMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LengthCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LengthFunc != nil {
		return atomic.LoadUint64(&m.LengthCounter) > 0
	}

	return true
}

type mUnsyncListMockRefToIndex struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockRefToIndexExpectation
	expectationSeries []*UnsyncListMockRefToIndexExpectation
}

type UnsyncListMockRefToIndexExpectation struct {
	input  *UnsyncListMockRefToIndexInput
	result *UnsyncListMockRefToIndexResult
}

type UnsyncListMockRefToIndexInput struct {
	p core.RecordRef
}

type UnsyncListMockRefToIndexResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of UnsyncList.RefToIndex is expected from 1 to Infinity times
func (m *mUnsyncListMockRefToIndex) Expect(p core.RecordRef) *mUnsyncListMockRefToIndex {
	m.mock.RefToIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockRefToIndexExpectation{}
	}
	m.mainExpectation.input = &UnsyncListMockRefToIndexInput{p}
	return m
}

//Return specifies results of invocation of UnsyncList.RefToIndex
func (m *mUnsyncListMockRefToIndex) Return(r int, r1 error) *UnsyncListMock {
	m.mock.RefToIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockRefToIndexExpectation{}
	}
	m.mainExpectation.result = &UnsyncListMockRefToIndexResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.RefToIndex is expected once
func (m *mUnsyncListMockRefToIndex) ExpectOnce(p core.RecordRef) *UnsyncListMockRefToIndexExpectation {
	m.mock.RefToIndexFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockRefToIndexExpectation{}
	expectation.input = &UnsyncListMockRefToIndexInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *UnsyncListMockRefToIndexExpectation) Return(r int, r1 error) {
	e.result = &UnsyncListMockRefToIndexResult{r, r1}
}

//Set uses given function f as a mock of UnsyncList.RefToIndex method
func (m *mUnsyncListMockRefToIndex) Set(f func(p core.RecordRef) (r int, r1 error)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RefToIndexFunc = f
	return m.mock
}

//RefToIndex implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) RefToIndex(p core.RecordRef) (r int, r1 error) {
	counter := atomic.AddUint64(&m.RefToIndexPreCounter, 1)
	defer atomic.AddUint64(&m.RefToIndexCounter, 1)

	if len(m.RefToIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RefToIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.RefToIndex. %v", p)
			return
		}

		input := m.RefToIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, UnsyncListMockRefToIndexInput{p}, "UnsyncList.RefToIndex got unexpected parameters")

		result := m.RefToIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.RefToIndex")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RefToIndexMock.mainExpectation != nil {

		input := m.RefToIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, UnsyncListMockRefToIndexInput{p}, "UnsyncList.RefToIndex got unexpected parameters")
		}

		result := m.RefToIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the UnsyncListMock.RefToIndex")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RefToIndexFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.RefToIndex. %v", p)
		return
	}

	return m.RefToIndexFunc(p)
}

//RefToIndexMinimockCounter returns a count of UnsyncListMock.RefToIndexFunc invocations
func (m *UnsyncListMock) RefToIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RefToIndexCounter)
}

//RefToIndexMinimockPreCounter returns the value of UnsyncListMock.RefToIndex invocations
func (m *UnsyncListMock) RefToIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RefToIndexPreCounter)
}

//RefToIndexFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) RefToIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RefToIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RefToIndexCounter) == uint64(len(m.RefToIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RefToIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RefToIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RefToIndexFunc != nil {
		return atomic.LoadUint64(&m.RefToIndexCounter) > 0
	}

	return true
}

type mUnsyncListMockRemoveClaims struct {
	mock              *UnsyncListMock
	mainExpectation   *UnsyncListMockRemoveClaimsExpectation
	expectationSeries []*UnsyncListMockRemoveClaimsExpectation
}

type UnsyncListMockRemoveClaimsExpectation struct {
	input *UnsyncListMockRemoveClaimsInput
}

type UnsyncListMockRemoveClaimsInput struct {
	p core.RecordRef
}

//Expect specifies that invocation of UnsyncList.RemoveClaims is expected from 1 to Infinity times
func (m *mUnsyncListMockRemoveClaims) Expect(p core.RecordRef) *mUnsyncListMockRemoveClaims {
	m.mock.RemoveClaimsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockRemoveClaimsExpectation{}
	}
	m.mainExpectation.input = &UnsyncListMockRemoveClaimsInput{p}
	return m
}

//Return specifies results of invocation of UnsyncList.RemoveClaims
func (m *mUnsyncListMockRemoveClaims) Return() *UnsyncListMock {
	m.mock.RemoveClaimsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &UnsyncListMockRemoveClaimsExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of UnsyncList.RemoveClaims is expected once
func (m *mUnsyncListMockRemoveClaims) ExpectOnce(p core.RecordRef) *UnsyncListMockRemoveClaimsExpectation {
	m.mock.RemoveClaimsFunc = nil
	m.mainExpectation = nil

	expectation := &UnsyncListMockRemoveClaimsExpectation{}
	expectation.input = &UnsyncListMockRemoveClaimsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of UnsyncList.RemoveClaims method
func (m *mUnsyncListMockRemoveClaims) Set(f func(p core.RecordRef)) *UnsyncListMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveClaimsFunc = f
	return m.mock
}

//RemoveClaims implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) RemoveClaims(p core.RecordRef) {
	counter := atomic.AddUint64(&m.RemoveClaimsPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveClaimsCounter, 1)

	if len(m.RemoveClaimsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveClaimsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to UnsyncListMock.RemoveClaims. %v", p)
			return
		}

		input := m.RemoveClaimsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, UnsyncListMockRemoveClaimsInput{p}, "UnsyncList.RemoveClaims got unexpected parameters")

		return
	}

	if m.RemoveClaimsMock.mainExpectation != nil {

		input := m.RemoveClaimsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, UnsyncListMockRemoveClaimsInput{p}, "UnsyncList.RemoveClaims got unexpected parameters")
		}

		return
	}

	if m.RemoveClaimsFunc == nil {
		m.t.Fatalf("Unexpected call to UnsyncListMock.RemoveClaims. %v", p)
		return
	}

	m.RemoveClaimsFunc(p)
}

//RemoveClaimsMinimockCounter returns a count of UnsyncListMock.RemoveClaimsFunc invocations
func (m *UnsyncListMock) RemoveClaimsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveClaimsCounter)
}

//RemoveClaimsMinimockPreCounter returns the value of UnsyncListMock.RemoveClaims invocations
func (m *UnsyncListMock) RemoveClaimsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveClaimsPreCounter)
}

//RemoveClaimsFinished returns true if mock invocations count is ok
func (m *UnsyncListMock) RemoveClaimsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveClaimsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveClaimsCounter) == uint64(len(m.RemoveClaimsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveClaimsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveClaimsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveClaimsFunc != nil {
		return atomic.LoadUint64(&m.RemoveClaimsCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *UnsyncListMock) ValidateCallCounters() {

	if !m.AddClaimsFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.AddClaims")
	}

	if !m.CalculateHashFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.CalculateHash")
	}

	if !m.GetActiveNodeFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.GetActiveNode")
	}

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.GetActiveNodes")
	}

	if !m.IndexToRefFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.IndexToRef")
	}

	if !m.LengthFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.Length")
	}

	if !m.RefToIndexFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.RefToIndex")
	}

	if !m.RemoveClaimsFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.RemoveClaims")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *UnsyncListMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *UnsyncListMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *UnsyncListMock) MinimockFinish() {

	if !m.AddClaimsFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.AddClaims")
	}

	if !m.CalculateHashFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.CalculateHash")
	}

	if !m.GetActiveNodeFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.GetActiveNode")
	}

	if !m.GetActiveNodesFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.GetActiveNodes")
	}

	if !m.IndexToRefFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.IndexToRef")
	}

	if !m.LengthFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.Length")
	}

	if !m.RefToIndexFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.RefToIndex")
	}

	if !m.RemoveClaimsFinished() {
		m.t.Fatal("Expected call to UnsyncListMock.RemoveClaims")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *UnsyncListMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *UnsyncListMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddClaimsFinished()
		ok = ok && m.CalculateHashFinished()
		ok = ok && m.GetActiveNodeFinished()
		ok = ok && m.GetActiveNodesFinished()
		ok = ok && m.IndexToRefFinished()
		ok = ok && m.LengthFinished()
		ok = ok && m.RefToIndexFinished()
		ok = ok && m.RemoveClaimsFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddClaimsFinished() {
				m.t.Error("Expected call to UnsyncListMock.AddClaims")
			}

			if !m.CalculateHashFinished() {
				m.t.Error("Expected call to UnsyncListMock.CalculateHash")
			}

			if !m.GetActiveNodeFinished() {
				m.t.Error("Expected call to UnsyncListMock.GetActiveNode")
			}

			if !m.GetActiveNodesFinished() {
				m.t.Error("Expected call to UnsyncListMock.GetActiveNodes")
			}

			if !m.IndexToRefFinished() {
				m.t.Error("Expected call to UnsyncListMock.IndexToRef")
			}

			if !m.LengthFinished() {
				m.t.Error("Expected call to UnsyncListMock.Length")
			}

			if !m.RefToIndexFinished() {
				m.t.Error("Expected call to UnsyncListMock.RefToIndex")
			}

			if !m.RemoveClaimsFinished() {
				m.t.Error("Expected call to UnsyncListMock.RemoveClaims")
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
func (m *UnsyncListMock) AllMocksCalled() bool {

	if !m.AddClaimsFinished() {
		return false
	}

	if !m.CalculateHashFinished() {
		return false
	}

	if !m.GetActiveNodeFinished() {
		return false
	}

	if !m.GetActiveNodesFinished() {
		return false
	}

	if !m.IndexToRefFinished() {
		return false
	}

	if !m.LengthFinished() {
		return false
	}

	if !m.RefToIndexFinished() {
		return false
	}

	if !m.RemoveClaimsFinished() {
		return false
	}

	return true
}
