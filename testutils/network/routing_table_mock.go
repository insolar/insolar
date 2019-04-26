package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "RoutingTable" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	network "github.com/insolar/insolar/network"
	host "github.com/insolar/insolar/network/hostnetwork/host"

	testify_assert "github.com/stretchr/testify/assert"
)

//RoutingTableMock implements github.com/insolar/insolar/network.RoutingTable
type RoutingTableMock struct {
	t minimock.Tester

	AddToKnownHostsFunc       func(p *host.Host)
	AddToKnownHostsCounter    uint64
	AddToKnownHostsPreCounter uint64
	AddToKnownHostsMock       mRoutingTableMockAddToKnownHosts

	RebalanceFunc       func(p network.PartitionPolicy)
	RebalanceCounter    uint64
	RebalancePreCounter uint64
	RebalanceMock       mRoutingTableMockRebalance

	ResolveFunc       func(p insolar.Reference) (r *host.Host, r1 error)
	ResolveCounter    uint64
	ResolvePreCounter uint64
	ResolveMock       mRoutingTableMockResolve

	ResolveConsensusFunc       func(p insolar.ShortNodeID) (r *host.Host, r1 error)
	ResolveConsensusCounter    uint64
	ResolveConsensusPreCounter uint64
	ResolveConsensusMock       mRoutingTableMockResolveConsensus

	ResolveConsensusRefFunc       func(p insolar.Reference) (r *host.Host, r1 error)
	ResolveConsensusRefCounter    uint64
	ResolveConsensusRefPreCounter uint64
	ResolveConsensusRefMock       mRoutingTableMockResolveConsensusRef
}

//NewRoutingTableMock returns a mock for github.com/insolar/insolar/network.RoutingTable
func NewRoutingTableMock(t minimock.Tester) *RoutingTableMock {
	m := &RoutingTableMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddToKnownHostsMock = mRoutingTableMockAddToKnownHosts{mock: m}
	m.RebalanceMock = mRoutingTableMockRebalance{mock: m}
	m.ResolveMock = mRoutingTableMockResolve{mock: m}
	m.ResolveConsensusMock = mRoutingTableMockResolveConsensus{mock: m}
	m.ResolveConsensusRefMock = mRoutingTableMockResolveConsensusRef{mock: m}

	return m
}

type mRoutingTableMockAddToKnownHosts struct {
	mock              *RoutingTableMock
	mainExpectation   *RoutingTableMockAddToKnownHostsExpectation
	expectationSeries []*RoutingTableMockAddToKnownHostsExpectation
}

type RoutingTableMockAddToKnownHostsExpectation struct {
	input *RoutingTableMockAddToKnownHostsInput
}

type RoutingTableMockAddToKnownHostsInput struct {
	p *host.Host
}

//Expect specifies that invocation of RoutingTable.AddToKnownHosts is expected from 1 to Infinity times
func (m *mRoutingTableMockAddToKnownHosts) Expect(p *host.Host) *mRoutingTableMockAddToKnownHosts {
	m.mock.AddToKnownHostsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockAddToKnownHostsExpectation{}
	}
	m.mainExpectation.input = &RoutingTableMockAddToKnownHostsInput{p}
	return m
}

//Return specifies results of invocation of RoutingTable.AddToKnownHosts
func (m *mRoutingTableMockAddToKnownHosts) Return() *RoutingTableMock {
	m.mock.AddToKnownHostsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockAddToKnownHostsExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RoutingTable.AddToKnownHosts is expected once
func (m *mRoutingTableMockAddToKnownHosts) ExpectOnce(p *host.Host) *RoutingTableMockAddToKnownHostsExpectation {
	m.mock.AddToKnownHostsFunc = nil
	m.mainExpectation = nil

	expectation := &RoutingTableMockAddToKnownHostsExpectation{}
	expectation.input = &RoutingTableMockAddToKnownHostsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RoutingTable.AddToKnownHosts method
func (m *mRoutingTableMockAddToKnownHosts) Set(f func(p *host.Host)) *RoutingTableMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddToKnownHostsFunc = f
	return m.mock
}

//AddToKnownHosts implements github.com/insolar/insolar/network.RoutingTable interface
func (m *RoutingTableMock) AddToKnownHosts(p *host.Host) {
	counter := atomic.AddUint64(&m.AddToKnownHostsPreCounter, 1)
	defer atomic.AddUint64(&m.AddToKnownHostsCounter, 1)

	if len(m.AddToKnownHostsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddToKnownHostsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RoutingTableMock.AddToKnownHosts. %v", p)
			return
		}

		input := m.AddToKnownHostsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RoutingTableMockAddToKnownHostsInput{p}, "RoutingTable.AddToKnownHosts got unexpected parameters")

		return
	}

	if m.AddToKnownHostsMock.mainExpectation != nil {

		input := m.AddToKnownHostsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RoutingTableMockAddToKnownHostsInput{p}, "RoutingTable.AddToKnownHosts got unexpected parameters")
		}

		return
	}

	if m.AddToKnownHostsFunc == nil {
		m.t.Fatalf("Unexpected call to RoutingTableMock.AddToKnownHosts. %v", p)
		return
	}

	m.AddToKnownHostsFunc(p)
}

//AddToKnownHostsMinimockCounter returns a count of RoutingTableMock.AddToKnownHostsFunc invocations
func (m *RoutingTableMock) AddToKnownHostsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddToKnownHostsCounter)
}

//AddToKnownHostsMinimockPreCounter returns the value of RoutingTableMock.AddToKnownHosts invocations
func (m *RoutingTableMock) AddToKnownHostsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddToKnownHostsPreCounter)
}

//AddToKnownHostsFinished returns true if mock invocations count is ok
func (m *RoutingTableMock) AddToKnownHostsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddToKnownHostsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddToKnownHostsCounter) == uint64(len(m.AddToKnownHostsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddToKnownHostsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddToKnownHostsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddToKnownHostsFunc != nil {
		return atomic.LoadUint64(&m.AddToKnownHostsCounter) > 0
	}

	return true
}

type mRoutingTableMockRebalance struct {
	mock              *RoutingTableMock
	mainExpectation   *RoutingTableMockRebalanceExpectation
	expectationSeries []*RoutingTableMockRebalanceExpectation
}

type RoutingTableMockRebalanceExpectation struct {
	input *RoutingTableMockRebalanceInput
}

type RoutingTableMockRebalanceInput struct {
	p network.PartitionPolicy
}

//Expect specifies that invocation of RoutingTable.Rebalance is expected from 1 to Infinity times
func (m *mRoutingTableMockRebalance) Expect(p network.PartitionPolicy) *mRoutingTableMockRebalance {
	m.mock.RebalanceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockRebalanceExpectation{}
	}
	m.mainExpectation.input = &RoutingTableMockRebalanceInput{p}
	return m
}

//Return specifies results of invocation of RoutingTable.Rebalance
func (m *mRoutingTableMockRebalance) Return() *RoutingTableMock {
	m.mock.RebalanceFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockRebalanceExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of RoutingTable.Rebalance is expected once
func (m *mRoutingTableMockRebalance) ExpectOnce(p network.PartitionPolicy) *RoutingTableMockRebalanceExpectation {
	m.mock.RebalanceFunc = nil
	m.mainExpectation = nil

	expectation := &RoutingTableMockRebalanceExpectation{}
	expectation.input = &RoutingTableMockRebalanceInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of RoutingTable.Rebalance method
func (m *mRoutingTableMockRebalance) Set(f func(p network.PartitionPolicy)) *RoutingTableMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RebalanceFunc = f
	return m.mock
}

//Rebalance implements github.com/insolar/insolar/network.RoutingTable interface
func (m *RoutingTableMock) Rebalance(p network.PartitionPolicy) {
	counter := atomic.AddUint64(&m.RebalancePreCounter, 1)
	defer atomic.AddUint64(&m.RebalanceCounter, 1)

	if len(m.RebalanceMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RebalanceMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RoutingTableMock.Rebalance. %v", p)
			return
		}

		input := m.RebalanceMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RoutingTableMockRebalanceInput{p}, "RoutingTable.Rebalance got unexpected parameters")

		return
	}

	if m.RebalanceMock.mainExpectation != nil {

		input := m.RebalanceMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RoutingTableMockRebalanceInput{p}, "RoutingTable.Rebalance got unexpected parameters")
		}

		return
	}

	if m.RebalanceFunc == nil {
		m.t.Fatalf("Unexpected call to RoutingTableMock.Rebalance. %v", p)
		return
	}

	m.RebalanceFunc(p)
}

//RebalanceMinimockCounter returns a count of RoutingTableMock.RebalanceFunc invocations
func (m *RoutingTableMock) RebalanceMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RebalanceCounter)
}

//RebalanceMinimockPreCounter returns the value of RoutingTableMock.Rebalance invocations
func (m *RoutingTableMock) RebalanceMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RebalancePreCounter)
}

//RebalanceFinished returns true if mock invocations count is ok
func (m *RoutingTableMock) RebalanceFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RebalanceMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RebalanceCounter) == uint64(len(m.RebalanceMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RebalanceMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RebalanceCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RebalanceFunc != nil {
		return atomic.LoadUint64(&m.RebalanceCounter) > 0
	}

	return true
}

type mRoutingTableMockResolve struct {
	mock              *RoutingTableMock
	mainExpectation   *RoutingTableMockResolveExpectation
	expectationSeries []*RoutingTableMockResolveExpectation
}

type RoutingTableMockResolveExpectation struct {
	input  *RoutingTableMockResolveInput
	result *RoutingTableMockResolveResult
}

type RoutingTableMockResolveInput struct {
	p insolar.Reference
}

type RoutingTableMockResolveResult struct {
	r  *host.Host
	r1 error
}

//Expect specifies that invocation of RoutingTable.Resolve is expected from 1 to Infinity times
func (m *mRoutingTableMockResolve) Expect(p insolar.Reference) *mRoutingTableMockResolve {
	m.mock.ResolveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockResolveExpectation{}
	}
	m.mainExpectation.input = &RoutingTableMockResolveInput{p}
	return m
}

//Return specifies results of invocation of RoutingTable.Resolve
func (m *mRoutingTableMockResolve) Return(r *host.Host, r1 error) *RoutingTableMock {
	m.mock.ResolveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockResolveExpectation{}
	}
	m.mainExpectation.result = &RoutingTableMockResolveResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RoutingTable.Resolve is expected once
func (m *mRoutingTableMockResolve) ExpectOnce(p insolar.Reference) *RoutingTableMockResolveExpectation {
	m.mock.ResolveFunc = nil
	m.mainExpectation = nil

	expectation := &RoutingTableMockResolveExpectation{}
	expectation.input = &RoutingTableMockResolveInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RoutingTableMockResolveExpectation) Return(r *host.Host, r1 error) {
	e.result = &RoutingTableMockResolveResult{r, r1}
}

//Set uses given function f as a mock of RoutingTable.Resolve method
func (m *mRoutingTableMockResolve) Set(f func(p insolar.Reference) (r *host.Host, r1 error)) *RoutingTableMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResolveFunc = f
	return m.mock
}

//Resolve implements github.com/insolar/insolar/network.RoutingTable interface
func (m *RoutingTableMock) Resolve(p insolar.Reference) (r *host.Host, r1 error) {
	counter := atomic.AddUint64(&m.ResolvePreCounter, 1)
	defer atomic.AddUint64(&m.ResolveCounter, 1)

	if len(m.ResolveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ResolveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RoutingTableMock.Resolve. %v", p)
			return
		}

		input := m.ResolveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RoutingTableMockResolveInput{p}, "RoutingTable.Resolve got unexpected parameters")

		result := m.ResolveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RoutingTableMock.Resolve")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResolveMock.mainExpectation != nil {

		input := m.ResolveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RoutingTableMockResolveInput{p}, "RoutingTable.Resolve got unexpected parameters")
		}

		result := m.ResolveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RoutingTableMock.Resolve")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResolveFunc == nil {
		m.t.Fatalf("Unexpected call to RoutingTableMock.Resolve. %v", p)
		return
	}

	return m.ResolveFunc(p)
}

//ResolveMinimockCounter returns a count of RoutingTableMock.ResolveFunc invocations
func (m *RoutingTableMock) ResolveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResolveCounter)
}

//ResolveMinimockPreCounter returns the value of RoutingTableMock.Resolve invocations
func (m *RoutingTableMock) ResolveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResolvePreCounter)
}

//ResolveFinished returns true if mock invocations count is ok
func (m *RoutingTableMock) ResolveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ResolveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ResolveCounter) == uint64(len(m.ResolveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ResolveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ResolveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ResolveFunc != nil {
		return atomic.LoadUint64(&m.ResolveCounter) > 0
	}

	return true
}

type mRoutingTableMockResolveConsensus struct {
	mock              *RoutingTableMock
	mainExpectation   *RoutingTableMockResolveConsensusExpectation
	expectationSeries []*RoutingTableMockResolveConsensusExpectation
}

type RoutingTableMockResolveConsensusExpectation struct {
	input  *RoutingTableMockResolveConsensusInput
	result *RoutingTableMockResolveConsensusResult
}

type RoutingTableMockResolveConsensusInput struct {
	p insolar.ShortNodeID
}

type RoutingTableMockResolveConsensusResult struct {
	r  *host.Host
	r1 error
}

//Expect specifies that invocation of RoutingTable.ResolveConsensus is expected from 1 to Infinity times
func (m *mRoutingTableMockResolveConsensus) Expect(p insolar.ShortNodeID) *mRoutingTableMockResolveConsensus {
	m.mock.ResolveConsensusFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockResolveConsensusExpectation{}
	}
	m.mainExpectation.input = &RoutingTableMockResolveConsensusInput{p}
	return m
}

//Return specifies results of invocation of RoutingTable.ResolveConsensus
func (m *mRoutingTableMockResolveConsensus) Return(r *host.Host, r1 error) *RoutingTableMock {
	m.mock.ResolveConsensusFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockResolveConsensusExpectation{}
	}
	m.mainExpectation.result = &RoutingTableMockResolveConsensusResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RoutingTable.ResolveConsensus is expected once
func (m *mRoutingTableMockResolveConsensus) ExpectOnce(p insolar.ShortNodeID) *RoutingTableMockResolveConsensusExpectation {
	m.mock.ResolveConsensusFunc = nil
	m.mainExpectation = nil

	expectation := &RoutingTableMockResolveConsensusExpectation{}
	expectation.input = &RoutingTableMockResolveConsensusInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RoutingTableMockResolveConsensusExpectation) Return(r *host.Host, r1 error) {
	e.result = &RoutingTableMockResolveConsensusResult{r, r1}
}

//Set uses given function f as a mock of RoutingTable.ResolveConsensus method
func (m *mRoutingTableMockResolveConsensus) Set(f func(p insolar.ShortNodeID) (r *host.Host, r1 error)) *RoutingTableMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResolveConsensusFunc = f
	return m.mock
}

//ResolveConsensus implements github.com/insolar/insolar/network.RoutingTable interface
func (m *RoutingTableMock) ResolveConsensus(p insolar.ShortNodeID) (r *host.Host, r1 error) {
	counter := atomic.AddUint64(&m.ResolveConsensusPreCounter, 1)
	defer atomic.AddUint64(&m.ResolveConsensusCounter, 1)

	if len(m.ResolveConsensusMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ResolveConsensusMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RoutingTableMock.ResolveConsensus. %v", p)
			return
		}

		input := m.ResolveConsensusMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RoutingTableMockResolveConsensusInput{p}, "RoutingTable.ResolveConsensus got unexpected parameters")

		result := m.ResolveConsensusMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RoutingTableMock.ResolveConsensus")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResolveConsensusMock.mainExpectation != nil {

		input := m.ResolveConsensusMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RoutingTableMockResolveConsensusInput{p}, "RoutingTable.ResolveConsensus got unexpected parameters")
		}

		result := m.ResolveConsensusMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RoutingTableMock.ResolveConsensus")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResolveConsensusFunc == nil {
		m.t.Fatalf("Unexpected call to RoutingTableMock.ResolveConsensus. %v", p)
		return
	}

	return m.ResolveConsensusFunc(p)
}

//ResolveConsensusMinimockCounter returns a count of RoutingTableMock.ResolveConsensusFunc invocations
func (m *RoutingTableMock) ResolveConsensusMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResolveConsensusCounter)
}

//ResolveConsensusMinimockPreCounter returns the value of RoutingTableMock.ResolveConsensus invocations
func (m *RoutingTableMock) ResolveConsensusMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResolveConsensusPreCounter)
}

//ResolveConsensusFinished returns true if mock invocations count is ok
func (m *RoutingTableMock) ResolveConsensusFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ResolveConsensusMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ResolveConsensusCounter) == uint64(len(m.ResolveConsensusMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ResolveConsensusMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ResolveConsensusCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ResolveConsensusFunc != nil {
		return atomic.LoadUint64(&m.ResolveConsensusCounter) > 0
	}

	return true
}

type mRoutingTableMockResolveConsensusRef struct {
	mock              *RoutingTableMock
	mainExpectation   *RoutingTableMockResolveConsensusRefExpectation
	expectationSeries []*RoutingTableMockResolveConsensusRefExpectation
}

type RoutingTableMockResolveConsensusRefExpectation struct {
	input  *RoutingTableMockResolveConsensusRefInput
	result *RoutingTableMockResolveConsensusRefResult
}

type RoutingTableMockResolveConsensusRefInput struct {
	p insolar.Reference
}

type RoutingTableMockResolveConsensusRefResult struct {
	r  *host.Host
	r1 error
}

//Expect specifies that invocation of RoutingTable.ResolveConsensusRef is expected from 1 to Infinity times
func (m *mRoutingTableMockResolveConsensusRef) Expect(p insolar.Reference) *mRoutingTableMockResolveConsensusRef {
	m.mock.ResolveConsensusRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockResolveConsensusRefExpectation{}
	}
	m.mainExpectation.input = &RoutingTableMockResolveConsensusRefInput{p}
	return m
}

//Return specifies results of invocation of RoutingTable.ResolveConsensusRef
func (m *mRoutingTableMockResolveConsensusRef) Return(r *host.Host, r1 error) *RoutingTableMock {
	m.mock.ResolveConsensusRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RoutingTableMockResolveConsensusRefExpectation{}
	}
	m.mainExpectation.result = &RoutingTableMockResolveConsensusRefResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of RoutingTable.ResolveConsensusRef is expected once
func (m *mRoutingTableMockResolveConsensusRef) ExpectOnce(p insolar.Reference) *RoutingTableMockResolveConsensusRefExpectation {
	m.mock.ResolveConsensusRefFunc = nil
	m.mainExpectation = nil

	expectation := &RoutingTableMockResolveConsensusRefExpectation{}
	expectation.input = &RoutingTableMockResolveConsensusRefInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RoutingTableMockResolveConsensusRefExpectation) Return(r *host.Host, r1 error) {
	e.result = &RoutingTableMockResolveConsensusRefResult{r, r1}
}

//Set uses given function f as a mock of RoutingTable.ResolveConsensusRef method
func (m *mRoutingTableMockResolveConsensusRef) Set(f func(p insolar.Reference) (r *host.Host, r1 error)) *RoutingTableMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResolveConsensusRefFunc = f
	return m.mock
}

//ResolveConsensusRef implements github.com/insolar/insolar/network.RoutingTable interface
func (m *RoutingTableMock) ResolveConsensusRef(p insolar.Reference) (r *host.Host, r1 error) {
	counter := atomic.AddUint64(&m.ResolveConsensusRefPreCounter, 1)
	defer atomic.AddUint64(&m.ResolveConsensusRefCounter, 1)

	if len(m.ResolveConsensusRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ResolveConsensusRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RoutingTableMock.ResolveConsensusRef. %v", p)
			return
		}

		input := m.ResolveConsensusRefMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, RoutingTableMockResolveConsensusRefInput{p}, "RoutingTable.ResolveConsensusRef got unexpected parameters")

		result := m.ResolveConsensusRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RoutingTableMock.ResolveConsensusRef")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResolveConsensusRefMock.mainExpectation != nil {

		input := m.ResolveConsensusRefMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, RoutingTableMockResolveConsensusRefInput{p}, "RoutingTable.ResolveConsensusRef got unexpected parameters")
		}

		result := m.ResolveConsensusRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RoutingTableMock.ResolveConsensusRef")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ResolveConsensusRefFunc == nil {
		m.t.Fatalf("Unexpected call to RoutingTableMock.ResolveConsensusRef. %v", p)
		return
	}

	return m.ResolveConsensusRefFunc(p)
}

//ResolveConsensusRefMinimockCounter returns a count of RoutingTableMock.ResolveConsensusRefFunc invocations
func (m *RoutingTableMock) ResolveConsensusRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResolveConsensusRefCounter)
}

//ResolveConsensusRefMinimockPreCounter returns the value of RoutingTableMock.ResolveConsensusRef invocations
func (m *RoutingTableMock) ResolveConsensusRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResolveConsensusRefPreCounter)
}

//ResolveConsensusRefFinished returns true if mock invocations count is ok
func (m *RoutingTableMock) ResolveConsensusRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ResolveConsensusRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ResolveConsensusRefCounter) == uint64(len(m.ResolveConsensusRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ResolveConsensusRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ResolveConsensusRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ResolveConsensusRefFunc != nil {
		return atomic.LoadUint64(&m.ResolveConsensusRefCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RoutingTableMock) ValidateCallCounters() {

	if !m.AddToKnownHostsFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.AddToKnownHosts")
	}

	if !m.RebalanceFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.Rebalance")
	}

	if !m.ResolveFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.Resolve")
	}

	if !m.ResolveConsensusFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.ResolveConsensus")
	}

	if !m.ResolveConsensusRefFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.ResolveConsensusRef")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RoutingTableMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RoutingTableMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RoutingTableMock) MinimockFinish() {

	if !m.AddToKnownHostsFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.AddToKnownHosts")
	}

	if !m.RebalanceFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.Rebalance")
	}

	if !m.ResolveFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.Resolve")
	}

	if !m.ResolveConsensusFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.ResolveConsensus")
	}

	if !m.ResolveConsensusRefFinished() {
		m.t.Fatal("Expected call to RoutingTableMock.ResolveConsensusRef")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RoutingTableMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RoutingTableMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddToKnownHostsFinished()
		ok = ok && m.RebalanceFinished()
		ok = ok && m.ResolveFinished()
		ok = ok && m.ResolveConsensusFinished()
		ok = ok && m.ResolveConsensusRefFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddToKnownHostsFinished() {
				m.t.Error("Expected call to RoutingTableMock.AddToKnownHosts")
			}

			if !m.RebalanceFinished() {
				m.t.Error("Expected call to RoutingTableMock.Rebalance")
			}

			if !m.ResolveFinished() {
				m.t.Error("Expected call to RoutingTableMock.Resolve")
			}

			if !m.ResolveConsensusFinished() {
				m.t.Error("Expected call to RoutingTableMock.ResolveConsensus")
			}

			if !m.ResolveConsensusRefFinished() {
				m.t.Error("Expected call to RoutingTableMock.ResolveConsensusRef")
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
func (m *RoutingTableMock) AllMocksCalled() bool {

	if !m.AddToKnownHostsFinished() {
		return false
	}

	if !m.RebalanceFinished() {
		return false
	}

	if !m.ResolveFinished() {
		return false
	}

	if !m.ResolveConsensusFinished() {
		return false
	}

	if !m.ResolveConsensusRefFinished() {
		return false
	}

	return true
}
