package phases

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Communicator" can be found in github.com/insolar/insolar/consensus/phases
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	packets "github.com/insolar/insolar/consensus/packets"
	core "github.com/insolar/insolar/core"
	network "github.com/insolar/insolar/network"

	testify_assert "github.com/stretchr/testify/assert"
)

//CommunicatorMock implements github.com/insolar/insolar/consensus/phases.Communicator
type CommunicatorMock struct {
	t minimock.Tester

	ExchangePhase1Func       func(p context.Context, p1 *packets.NodeAnnounceClaim, p2 []core.Node, p3 *packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error)
	ExchangePhase1Counter    uint64
	ExchangePhase1PreCounter uint64
	ExchangePhase1Mock       mCommunicatorMockExchangePhase1

	ExchangePhase2Func       func(p context.Context, p1 network.UnsyncList, p2 []core.Node, p3 *packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error)
	ExchangePhase2Counter    uint64
	ExchangePhase2PreCounter uint64
	ExchangePhase2Mock       mCommunicatorMockExchangePhase2

	ExchangePhase21Func       func(p context.Context, p1 network.UnsyncList, p2 *packets.Phase2Packet, p3 []*AdditionalRequest) (r []packets.ReferendumVote, r1 error)
	ExchangePhase21Counter    uint64
	ExchangePhase21PreCounter uint64
	ExchangePhase21Mock       mCommunicatorMockExchangePhase21

	ExchangePhase3Func       func(p context.Context, p1 []core.Node, p2 *packets.Phase3Packet) (r map[core.RecordRef]*packets.Phase3Packet, r1 error)
	ExchangePhase3Counter    uint64
	ExchangePhase3PreCounter uint64
	ExchangePhase3Mock       mCommunicatorMockExchangePhase3

	StartFunc       func(p context.Context) (r error)
	StartCounter    uint64
	StartPreCounter uint64
	StartMock       mCommunicatorMockStart
}

//NewCommunicatorMock returns a mock for github.com/insolar/insolar/consensus/phases.Communicator
func NewCommunicatorMock(t minimock.Tester) *CommunicatorMock {
	m := &CommunicatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExchangePhase1Mock = mCommunicatorMockExchangePhase1{mock: m}
	m.ExchangePhase2Mock = mCommunicatorMockExchangePhase2{mock: m}
	m.ExchangePhase21Mock = mCommunicatorMockExchangePhase21{mock: m}
	m.ExchangePhase3Mock = mCommunicatorMockExchangePhase3{mock: m}
	m.StartMock = mCommunicatorMockStart{mock: m}

	return m
}

type mCommunicatorMockExchangePhase1 struct {
	mock              *CommunicatorMock
	mainExpectation   *CommunicatorMockExchangePhase1Expectation
	expectationSeries []*CommunicatorMockExchangePhase1Expectation
}

type CommunicatorMockExchangePhase1Expectation struct {
	input  *CommunicatorMockExchangePhase1Input
	result *CommunicatorMockExchangePhase1Result
}

type CommunicatorMockExchangePhase1Input struct {
	p  context.Context
	p1 *packets.NodeAnnounceClaim
	p2 []core.Node
	p3 *packets.Phase1Packet
}

type CommunicatorMockExchangePhase1Result struct {
	r  map[core.RecordRef]*packets.Phase1Packet
	r1 error
}

//Expect specifies that invocation of Communicator.ExchangePhase1 is expected from 1 to Infinity times
func (m *mCommunicatorMockExchangePhase1) Expect(p context.Context, p1 *packets.NodeAnnounceClaim, p2 []core.Node, p3 *packets.Phase1Packet) *mCommunicatorMockExchangePhase1 {
	m.mock.ExchangePhase1Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase1Expectation{}
	}
	m.mainExpectation.input = &CommunicatorMockExchangePhase1Input{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Communicator.ExchangePhase1
func (m *mCommunicatorMockExchangePhase1) Return(r map[core.RecordRef]*packets.Phase1Packet, r1 error) *CommunicatorMock {
	m.mock.ExchangePhase1Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase1Expectation{}
	}
	m.mainExpectation.result = &CommunicatorMockExchangePhase1Result{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Communicator.ExchangePhase1 is expected once
func (m *mCommunicatorMockExchangePhase1) ExpectOnce(p context.Context, p1 *packets.NodeAnnounceClaim, p2 []core.Node, p3 *packets.Phase1Packet) *CommunicatorMockExchangePhase1Expectation {
	m.mock.ExchangePhase1Func = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockExchangePhase1Expectation{}
	expectation.input = &CommunicatorMockExchangePhase1Input{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockExchangePhase1Expectation) Return(r map[core.RecordRef]*packets.Phase1Packet, r1 error) {
	e.result = &CommunicatorMockExchangePhase1Result{r, r1}
}

//Set uses given function f as a mock of Communicator.ExchangePhase1 method
func (m *mCommunicatorMockExchangePhase1) Set(f func(p context.Context, p1 *packets.NodeAnnounceClaim, p2 []core.Node, p3 *packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExchangePhase1Func = f
	return m.mock
}

//ExchangePhase1 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase1(p context.Context, p1 *packets.NodeAnnounceClaim, p2 []core.Node, p3 *packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error) {
	counter := atomic.AddUint64(&m.ExchangePhase1PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase1Counter, 1)

	if len(m.ExchangePhase1Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExchangePhase1Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase1. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.ExchangePhase1Mock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase1Input{p, p1, p2, p3}, "Communicator.ExchangePhase1 got unexpected parameters")

		result := m.ExchangePhase1Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase1")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase1Mock.mainExpectation != nil {

		input := m.ExchangePhase1Mock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase1Input{p, p1, p2, p3}, "Communicator.ExchangePhase1 got unexpected parameters")
		}

		result := m.ExchangePhase1Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase1")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase1Func == nil {
		m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase1. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.ExchangePhase1Func(p, p1, p2, p3)
}

//ExchangePhase1MinimockCounter returns a count of CommunicatorMock.ExchangePhase1Func invocations
func (m *CommunicatorMock) ExchangePhase1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase1Counter)
}

//ExchangePhase1MinimockPreCounter returns the value of CommunicatorMock.ExchangePhase1 invocations
func (m *CommunicatorMock) ExchangePhase1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase1PreCounter)
}

//ExchangePhase1Finished returns true if mock invocations count is ok
func (m *CommunicatorMock) ExchangePhase1Finished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExchangePhase1Mock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExchangePhase1Counter) == uint64(len(m.ExchangePhase1Mock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExchangePhase1Mock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExchangePhase1Counter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExchangePhase1Func != nil {
		return atomic.LoadUint64(&m.ExchangePhase1Counter) > 0
	}

	return true
}

type mCommunicatorMockExchangePhase2 struct {
	mock              *CommunicatorMock
	mainExpectation   *CommunicatorMockExchangePhase2Expectation
	expectationSeries []*CommunicatorMockExchangePhase2Expectation
}

type CommunicatorMockExchangePhase2Expectation struct {
	input  *CommunicatorMockExchangePhase2Input
	result *CommunicatorMockExchangePhase2Result
}

type CommunicatorMockExchangePhase2Input struct {
	p  context.Context
	p1 network.UnsyncList
	p2 []core.Node
	p3 *packets.Phase2Packet
}

type CommunicatorMockExchangePhase2Result struct {
	r  map[core.RecordRef]*packets.Phase2Packet
	r1 error
}

//Expect specifies that invocation of Communicator.ExchangePhase2 is expected from 1 to Infinity times
func (m *mCommunicatorMockExchangePhase2) Expect(p context.Context, p1 network.UnsyncList, p2 []core.Node, p3 *packets.Phase2Packet) *mCommunicatorMockExchangePhase2 {
	m.mock.ExchangePhase2Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase2Expectation{}
	}
	m.mainExpectation.input = &CommunicatorMockExchangePhase2Input{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Communicator.ExchangePhase2
func (m *mCommunicatorMockExchangePhase2) Return(r map[core.RecordRef]*packets.Phase2Packet, r1 error) *CommunicatorMock {
	m.mock.ExchangePhase2Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase2Expectation{}
	}
	m.mainExpectation.result = &CommunicatorMockExchangePhase2Result{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Communicator.ExchangePhase2 is expected once
func (m *mCommunicatorMockExchangePhase2) ExpectOnce(p context.Context, p1 network.UnsyncList, p2 []core.Node, p3 *packets.Phase2Packet) *CommunicatorMockExchangePhase2Expectation {
	m.mock.ExchangePhase2Func = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockExchangePhase2Expectation{}
	expectation.input = &CommunicatorMockExchangePhase2Input{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockExchangePhase2Expectation) Return(r map[core.RecordRef]*packets.Phase2Packet, r1 error) {
	e.result = &CommunicatorMockExchangePhase2Result{r, r1}
}

//Set uses given function f as a mock of Communicator.ExchangePhase2 method
func (m *mCommunicatorMockExchangePhase2) Set(f func(p context.Context, p1 network.UnsyncList, p2 []core.Node, p3 *packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExchangePhase2Func = f
	return m.mock
}

//ExchangePhase2 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase2(p context.Context, p1 network.UnsyncList, p2 []core.Node, p3 *packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error) {
	counter := atomic.AddUint64(&m.ExchangePhase2PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase2Counter, 1)

	if len(m.ExchangePhase2Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExchangePhase2Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase2. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.ExchangePhase2Mock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase2Input{p, p1, p2, p3}, "Communicator.ExchangePhase2 got unexpected parameters")

		result := m.ExchangePhase2Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase2")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase2Mock.mainExpectation != nil {

		input := m.ExchangePhase2Mock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase2Input{p, p1, p2, p3}, "Communicator.ExchangePhase2 got unexpected parameters")
		}

		result := m.ExchangePhase2Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase2")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase2Func == nil {
		m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase2. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.ExchangePhase2Func(p, p1, p2, p3)
}

//ExchangePhase2MinimockCounter returns a count of CommunicatorMock.ExchangePhase2Func invocations
func (m *CommunicatorMock) ExchangePhase2MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase2Counter)
}

//ExchangePhase2MinimockPreCounter returns the value of CommunicatorMock.ExchangePhase2 invocations
func (m *CommunicatorMock) ExchangePhase2MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase2PreCounter)
}

//ExchangePhase2Finished returns true if mock invocations count is ok
func (m *CommunicatorMock) ExchangePhase2Finished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExchangePhase2Mock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExchangePhase2Counter) == uint64(len(m.ExchangePhase2Mock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExchangePhase2Mock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExchangePhase2Counter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExchangePhase2Func != nil {
		return atomic.LoadUint64(&m.ExchangePhase2Counter) > 0
	}

	return true
}

type mCommunicatorMockExchangePhase21 struct {
	mock              *CommunicatorMock
	mainExpectation   *CommunicatorMockExchangePhase21Expectation
	expectationSeries []*CommunicatorMockExchangePhase21Expectation
}

type CommunicatorMockExchangePhase21Expectation struct {
	input  *CommunicatorMockExchangePhase21Input
	result *CommunicatorMockExchangePhase21Result
}

type CommunicatorMockExchangePhase21Input struct {
	p  context.Context
	p1 network.UnsyncList
	p2 *packets.Phase2Packet
	p3 []*AdditionalRequest
}

type CommunicatorMockExchangePhase21Result struct {
	r  []packets.ReferendumVote
	r1 error
}

//Expect specifies that invocation of Communicator.ExchangePhase21 is expected from 1 to Infinity times
func (m *mCommunicatorMockExchangePhase21) Expect(p context.Context, p1 network.UnsyncList, p2 *packets.Phase2Packet, p3 []*AdditionalRequest) *mCommunicatorMockExchangePhase21 {
	m.mock.ExchangePhase21Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase21Expectation{}
	}
	m.mainExpectation.input = &CommunicatorMockExchangePhase21Input{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Communicator.ExchangePhase21
func (m *mCommunicatorMockExchangePhase21) Return(r []packets.ReferendumVote, r1 error) *CommunicatorMock {
	m.mock.ExchangePhase21Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase21Expectation{}
	}
	m.mainExpectation.result = &CommunicatorMockExchangePhase21Result{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Communicator.ExchangePhase21 is expected once
func (m *mCommunicatorMockExchangePhase21) ExpectOnce(p context.Context, p1 network.UnsyncList, p2 *packets.Phase2Packet, p3 []*AdditionalRequest) *CommunicatorMockExchangePhase21Expectation {
	m.mock.ExchangePhase21Func = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockExchangePhase21Expectation{}
	expectation.input = &CommunicatorMockExchangePhase21Input{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockExchangePhase21Expectation) Return(r []packets.ReferendumVote, r1 error) {
	e.result = &CommunicatorMockExchangePhase21Result{r, r1}
}

//Set uses given function f as a mock of Communicator.ExchangePhase21 method
func (m *mCommunicatorMockExchangePhase21) Set(f func(p context.Context, p1 network.UnsyncList, p2 *packets.Phase2Packet, p3 []*AdditionalRequest) (r []packets.ReferendumVote, r1 error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExchangePhase21Func = f
	return m.mock
}

//ExchangePhase21 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase21(p context.Context, p1 network.UnsyncList, p2 *packets.Phase2Packet, p3 []*AdditionalRequest) (r []packets.ReferendumVote, r1 error) {
	counter := atomic.AddUint64(&m.ExchangePhase21PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase21Counter, 1)

	if len(m.ExchangePhase21Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExchangePhase21Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase21. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.ExchangePhase21Mock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase21Input{p, p1, p2, p3}, "Communicator.ExchangePhase21 got unexpected parameters")

		result := m.ExchangePhase21Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase21")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase21Mock.mainExpectation != nil {

		input := m.ExchangePhase21Mock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase21Input{p, p1, p2, p3}, "Communicator.ExchangePhase21 got unexpected parameters")
		}

		result := m.ExchangePhase21Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase21")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase21Func == nil {
		m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase21. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.ExchangePhase21Func(p, p1, p2, p3)
}

//ExchangePhase21MinimockCounter returns a count of CommunicatorMock.ExchangePhase21Func invocations
func (m *CommunicatorMock) ExchangePhase21MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase21Counter)
}

//ExchangePhase21MinimockPreCounter returns the value of CommunicatorMock.ExchangePhase21 invocations
func (m *CommunicatorMock) ExchangePhase21MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase21PreCounter)
}

//ExchangePhase21Finished returns true if mock invocations count is ok
func (m *CommunicatorMock) ExchangePhase21Finished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExchangePhase21Mock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExchangePhase21Counter) == uint64(len(m.ExchangePhase21Mock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExchangePhase21Mock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExchangePhase21Counter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExchangePhase21Func != nil {
		return atomic.LoadUint64(&m.ExchangePhase21Counter) > 0
	}

	return true
}

type mCommunicatorMockExchangePhase3 struct {
	mock              *CommunicatorMock
	mainExpectation   *CommunicatorMockExchangePhase3Expectation
	expectationSeries []*CommunicatorMockExchangePhase3Expectation
}

type CommunicatorMockExchangePhase3Expectation struct {
	input  *CommunicatorMockExchangePhase3Input
	result *CommunicatorMockExchangePhase3Result
}

type CommunicatorMockExchangePhase3Input struct {
	p  context.Context
	p1 []core.Node
	p2 *packets.Phase3Packet
}

type CommunicatorMockExchangePhase3Result struct {
	r  map[core.RecordRef]*packets.Phase3Packet
	r1 error
}

//Expect specifies that invocation of Communicator.ExchangePhase3 is expected from 1 to Infinity times
func (m *mCommunicatorMockExchangePhase3) Expect(p context.Context, p1 []core.Node, p2 *packets.Phase3Packet) *mCommunicatorMockExchangePhase3 {
	m.mock.ExchangePhase3Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase3Expectation{}
	}
	m.mainExpectation.input = &CommunicatorMockExchangePhase3Input{p, p1, p2}
	return m
}

//Return specifies results of invocation of Communicator.ExchangePhase3
func (m *mCommunicatorMockExchangePhase3) Return(r map[core.RecordRef]*packets.Phase3Packet, r1 error) *CommunicatorMock {
	m.mock.ExchangePhase3Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase3Expectation{}
	}
	m.mainExpectation.result = &CommunicatorMockExchangePhase3Result{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Communicator.ExchangePhase3 is expected once
func (m *mCommunicatorMockExchangePhase3) ExpectOnce(p context.Context, p1 []core.Node, p2 *packets.Phase3Packet) *CommunicatorMockExchangePhase3Expectation {
	m.mock.ExchangePhase3Func = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockExchangePhase3Expectation{}
	expectation.input = &CommunicatorMockExchangePhase3Input{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockExchangePhase3Expectation) Return(r map[core.RecordRef]*packets.Phase3Packet, r1 error) {
	e.result = &CommunicatorMockExchangePhase3Result{r, r1}
}

//Set uses given function f as a mock of Communicator.ExchangePhase3 method
func (m *mCommunicatorMockExchangePhase3) Set(f func(p context.Context, p1 []core.Node, p2 *packets.Phase3Packet) (r map[core.RecordRef]*packets.Phase3Packet, r1 error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExchangePhase3Func = f
	return m.mock
}

//ExchangePhase3 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase3(p context.Context, p1 []core.Node, p2 *packets.Phase3Packet) (r map[core.RecordRef]*packets.Phase3Packet, r1 error) {
	counter := atomic.AddUint64(&m.ExchangePhase3PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase3Counter, 1)

	if len(m.ExchangePhase3Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExchangePhase3Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase3. %v %v %v", p, p1, p2)
			return
		}

		input := m.ExchangePhase3Mock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase3Input{p, p1, p2}, "Communicator.ExchangePhase3 got unexpected parameters")

		result := m.ExchangePhase3Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase3")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase3Mock.mainExpectation != nil {

		input := m.ExchangePhase3Mock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase3Input{p, p1, p2}, "Communicator.ExchangePhase3 got unexpected parameters")
		}

		result := m.ExchangePhase3Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase3")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExchangePhase3Func == nil {
		m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase3. %v %v %v", p, p1, p2)
		return
	}

	return m.ExchangePhase3Func(p, p1, p2)
}

//ExchangePhase3MinimockCounter returns a count of CommunicatorMock.ExchangePhase3Func invocations
func (m *CommunicatorMock) ExchangePhase3MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase3Counter)
}

//ExchangePhase3MinimockPreCounter returns the value of CommunicatorMock.ExchangePhase3 invocations
func (m *CommunicatorMock) ExchangePhase3MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase3PreCounter)
}

//ExchangePhase3Finished returns true if mock invocations count is ok
func (m *CommunicatorMock) ExchangePhase3Finished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExchangePhase3Mock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExchangePhase3Counter) == uint64(len(m.ExchangePhase3Mock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExchangePhase3Mock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExchangePhase3Counter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExchangePhase3Func != nil {
		return atomic.LoadUint64(&m.ExchangePhase3Counter) > 0
	}

	return true
}

type mCommunicatorMockStart struct {
	mock              *CommunicatorMock
	mainExpectation   *CommunicatorMockStartExpectation
	expectationSeries []*CommunicatorMockStartExpectation
}

type CommunicatorMockStartExpectation struct {
	input  *CommunicatorMockStartInput
	result *CommunicatorMockStartResult
}

type CommunicatorMockStartInput struct {
	p context.Context
}

type CommunicatorMockStartResult struct {
	r error
}

//Expect specifies that invocation of Communicator.Start is expected from 1 to Infinity times
func (m *mCommunicatorMockStart) Expect(p context.Context) *mCommunicatorMockStart {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockStartExpectation{}
	}
	m.mainExpectation.input = &CommunicatorMockStartInput{p}
	return m
}

//Return specifies results of invocation of Communicator.Start
func (m *mCommunicatorMockStart) Return(r error) *CommunicatorMock {
	m.mock.StartFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockStartExpectation{}
	}
	m.mainExpectation.result = &CommunicatorMockStartResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Communicator.Start is expected once
func (m *mCommunicatorMockStart) ExpectOnce(p context.Context) *CommunicatorMockStartExpectation {
	m.mock.StartFunc = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockStartExpectation{}
	expectation.input = &CommunicatorMockStartInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockStartExpectation) Return(r error) {
	e.result = &CommunicatorMockStartResult{r}
}

//Set uses given function f as a mock of Communicator.Start method
func (m *mCommunicatorMockStart) Set(f func(p context.Context) (r error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StartFunc = f
	return m.mock
}

//Start implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) Start(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.StartPreCounter, 1)
	defer atomic.AddUint64(&m.StartCounter, 1)

	if len(m.StartMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StartMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.Start. %v", p)
			return
		}

		input := m.StartMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockStartInput{p}, "Communicator.Start got unexpected parameters")

		result := m.StartMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.Start")
			return
		}

		r = result.r

		return
	}

	if m.StartMock.mainExpectation != nil {

		input := m.StartMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CommunicatorMockStartInput{p}, "Communicator.Start got unexpected parameters")
		}

		result := m.StartMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.Start")
		}

		r = result.r

		return
	}

	if m.StartFunc == nil {
		m.t.Fatalf("Unexpected call to CommunicatorMock.Start. %v", p)
		return
	}

	return m.StartFunc(p)
}

//StartMinimockCounter returns a count of CommunicatorMock.StartFunc invocations
func (m *CommunicatorMock) StartMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StartCounter)
}

//StartMinimockPreCounter returns the value of CommunicatorMock.Start invocations
func (m *CommunicatorMock) StartMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StartPreCounter)
}

//StartFinished returns true if mock invocations count is ok
func (m *CommunicatorMock) StartFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StartMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StartCounter) == uint64(len(m.StartMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StartMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StartFunc != nil {
		return atomic.LoadUint64(&m.StartCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CommunicatorMock) ValidateCallCounters() {

	if !m.ExchangePhase1Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase1")
	}

	if !m.ExchangePhase2Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase2")
	}

	if !m.ExchangePhase21Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase21")
	}

	if !m.ExchangePhase3Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase3")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to CommunicatorMock.Start")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CommunicatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CommunicatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CommunicatorMock) MinimockFinish() {

	if !m.ExchangePhase1Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase1")
	}

	if !m.ExchangePhase2Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase2")
	}

	if !m.ExchangePhase21Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase21")
	}

	if !m.ExchangePhase3Finished() {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase3")
	}

	if !m.StartFinished() {
		m.t.Fatal("Expected call to CommunicatorMock.Start")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CommunicatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CommunicatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ExchangePhase1Finished()
		ok = ok && m.ExchangePhase2Finished()
		ok = ok && m.ExchangePhase21Finished()
		ok = ok && m.ExchangePhase3Finished()
		ok = ok && m.StartFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ExchangePhase1Finished() {
				m.t.Error("Expected call to CommunicatorMock.ExchangePhase1")
			}

			if !m.ExchangePhase2Finished() {
				m.t.Error("Expected call to CommunicatorMock.ExchangePhase2")
			}

			if !m.ExchangePhase21Finished() {
				m.t.Error("Expected call to CommunicatorMock.ExchangePhase21")
			}

			if !m.ExchangePhase3Finished() {
				m.t.Error("Expected call to CommunicatorMock.ExchangePhase3")
			}

			if !m.StartFinished() {
				m.t.Error("Expected call to CommunicatorMock.Start")
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
func (m *CommunicatorMock) AllMocksCalled() bool {

	if !m.ExchangePhase1Finished() {
		return false
	}

	if !m.ExchangePhase2Finished() {
		return false
	}

	if !m.ExchangePhase21Finished() {
		return false
	}

	if !m.ExchangePhase3Finished() {
		return false
	}

	if !m.StartFinished() {
		return false
	}

	return true
}
