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
	claimhandler "github.com/insolar/insolar/consensus/claimhandler"
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

	ExchangePhase2Func       func(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 []core.Node, p4 *packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error)
	ExchangePhase2Counter    uint64
	ExchangePhase2PreCounter uint64
	ExchangePhase2Mock       mCommunicatorMockExchangePhase2

	ExchangePhase21Func       func(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 *packets.Phase2Packet, p4 []*AdditionalRequest) (r []packets.ReferendumVote, r1 error)
	ExchangePhase21Counter    uint64
	ExchangePhase21PreCounter uint64
	ExchangePhase21Mock       mCommunicatorMockExchangePhase21

	ExchangePhase3Func       func(p context.Context, p1 []core.Node, p2 *packets.Phase3Packet) (r map[core.RecordRef]*packets.Phase3Packet, r1 error)
	ExchangePhase3Counter    uint64
	ExchangePhase3PreCounter uint64
	ExchangePhase3Mock       mCommunicatorMockExchangePhase3

	InitFunc       func(p context.Context) (r error)
	InitCounter    uint64
	InitPreCounter uint64
	InitMock       mCommunicatorMockInit
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
	m.InitMock = mCommunicatorMockInit{mock: m}

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
	p2 *claimhandler.ClaimHandler
	p3 []core.Node
	p4 *packets.Phase2Packet
}

type CommunicatorMockExchangePhase2Result struct {
	r  map[core.RecordRef]*packets.Phase2Packet
	r1 error
}

//Expect specifies that invocation of Communicator.ExchangePhase2 is expected from 1 to Infinity times
func (m *mCommunicatorMockExchangePhase2) Expect(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 []core.Node, p4 *packets.Phase2Packet) *mCommunicatorMockExchangePhase2 {
	m.mock.ExchangePhase2Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase2Expectation{}
	}
	m.mainExpectation.input = &CommunicatorMockExchangePhase2Input{p, p1, p2, p3, p4}
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
func (m *mCommunicatorMockExchangePhase2) ExpectOnce(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 []core.Node, p4 *packets.Phase2Packet) *CommunicatorMockExchangePhase2Expectation {
	m.mock.ExchangePhase2Func = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockExchangePhase2Expectation{}
	expectation.input = &CommunicatorMockExchangePhase2Input{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockExchangePhase2Expectation) Return(r map[core.RecordRef]*packets.Phase2Packet, r1 error) {
	e.result = &CommunicatorMockExchangePhase2Result{r, r1}
}

//Set uses given function f as a mock of Communicator.ExchangePhase2 method
func (m *mCommunicatorMockExchangePhase2) Set(f func(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 []core.Node, p4 *packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExchangePhase2Func = f
	return m.mock
}

//ExchangePhase2 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase2(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 []core.Node, p4 *packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error) {
	counter := atomic.AddUint64(&m.ExchangePhase2PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase2Counter, 1)

	if len(m.ExchangePhase2Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExchangePhase2Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase2. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.ExchangePhase2Mock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase2Input{p, p1, p2, p3, p4}, "Communicator.ExchangePhase2 got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase2Input{p, p1, p2, p3, p4}, "Communicator.ExchangePhase2 got unexpected parameters")
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
		m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase2. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.ExchangePhase2Func(p, p1, p2, p3, p4)
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
	p2 *claimhandler.ClaimHandler
	p3 *packets.Phase2Packet
	p4 []*AdditionalRequest
}

type CommunicatorMockExchangePhase21Result struct {
	r  []packets.ReferendumVote
	r1 error
}

//Expect specifies that invocation of Communicator.ExchangePhase21 is expected from 1 to Infinity times
func (m *mCommunicatorMockExchangePhase21) Expect(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 *packets.Phase2Packet, p4 []*AdditionalRequest) *mCommunicatorMockExchangePhase21 {
	m.mock.ExchangePhase21Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockExchangePhase21Expectation{}
	}
	m.mainExpectation.input = &CommunicatorMockExchangePhase21Input{p, p1, p2, p3, p4}
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
func (m *mCommunicatorMockExchangePhase21) ExpectOnce(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 *packets.Phase2Packet, p4 []*AdditionalRequest) *CommunicatorMockExchangePhase21Expectation {
	m.mock.ExchangePhase21Func = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockExchangePhase21Expectation{}
	expectation.input = &CommunicatorMockExchangePhase21Input{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockExchangePhase21Expectation) Return(r []packets.ReferendumVote, r1 error) {
	e.result = &CommunicatorMockExchangePhase21Result{r, r1}
}

//Set uses given function f as a mock of Communicator.ExchangePhase21 method
func (m *mCommunicatorMockExchangePhase21) Set(f func(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 *packets.Phase2Packet, p4 []*AdditionalRequest) (r []packets.ReferendumVote, r1 error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExchangePhase21Func = f
	return m.mock
}

//ExchangePhase21 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase21(p context.Context, p1 network.UnsyncList, p2 *claimhandler.ClaimHandler, p3 *packets.Phase2Packet, p4 []*AdditionalRequest) (r []packets.ReferendumVote, r1 error) {
	counter := atomic.AddUint64(&m.ExchangePhase21PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase21Counter, 1)

	if len(m.ExchangePhase21Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExchangePhase21Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase21. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.ExchangePhase21Mock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase21Input{p, p1, p2, p3, p4}, "Communicator.ExchangePhase21 got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, CommunicatorMockExchangePhase21Input{p, p1, p2, p3, p4}, "Communicator.ExchangePhase21 got unexpected parameters")
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
		m.t.Fatalf("Unexpected call to CommunicatorMock.ExchangePhase21. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.ExchangePhase21Func(p, p1, p2, p3, p4)
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

type mCommunicatorMockInit struct {
	mock              *CommunicatorMock
	mainExpectation   *CommunicatorMockInitExpectation
	expectationSeries []*CommunicatorMockInitExpectation
}

type CommunicatorMockInitExpectation struct {
	input  *CommunicatorMockInitInput
	result *CommunicatorMockInitResult
}

type CommunicatorMockInitInput struct {
	p context.Context
}

type CommunicatorMockInitResult struct {
	r error
}

//Expect specifies that invocation of Communicator.Init is expected from 1 to Infinity times
func (m *mCommunicatorMockInit) Expect(p context.Context) *mCommunicatorMockInit {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockInitExpectation{}
	}
	m.mainExpectation.input = &CommunicatorMockInitInput{p}
	return m
}

//Return specifies results of invocation of Communicator.Init
func (m *mCommunicatorMockInit) Return(r error) *CommunicatorMock {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CommunicatorMockInitExpectation{}
	}
	m.mainExpectation.result = &CommunicatorMockInitResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Communicator.Init is expected once
func (m *mCommunicatorMockInit) ExpectOnce(p context.Context) *CommunicatorMockInitExpectation {
	m.mock.InitFunc = nil
	m.mainExpectation = nil

	expectation := &CommunicatorMockInitExpectation{}
	expectation.input = &CommunicatorMockInitInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CommunicatorMockInitExpectation) Return(r error) {
	e.result = &CommunicatorMockInitResult{r}
}

//Set uses given function f as a mock of Communicator.Init method
func (m *mCommunicatorMockInit) Set(f func(p context.Context) (r error)) *CommunicatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InitFunc = f
	return m.mock
}

//Init implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) Init(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.InitPreCounter, 1)
	defer atomic.AddUint64(&m.InitCounter, 1)

	if len(m.InitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CommunicatorMock.Init. %v", p)
			return
		}

		input := m.InitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CommunicatorMockInitInput{p}, "Communicator.Init got unexpected parameters")

		result := m.InitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.Init")
			return
		}

		r = result.r

		return
	}

	if m.InitMock.mainExpectation != nil {

		input := m.InitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CommunicatorMockInitInput{p}, "Communicator.Init got unexpected parameters")
		}

		result := m.InitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CommunicatorMock.Init")
		}

		r = result.r

		return
	}

	if m.InitFunc == nil {
		m.t.Fatalf("Unexpected call to CommunicatorMock.Init. %v", p)
		return
	}

	return m.InitFunc(p)
}

//InitMinimockCounter returns a count of CommunicatorMock.InitFunc invocations
func (m *CommunicatorMock) InitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InitCounter)
}

//InitMinimockPreCounter returns the value of CommunicatorMock.Init invocations
func (m *CommunicatorMock) InitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InitPreCounter)
}

//InitFinished returns true if mock invocations count is ok
func (m *CommunicatorMock) InitFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InitMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InitCounter) == uint64(len(m.InitMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InitMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InitCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InitFunc != nil {
		return atomic.LoadUint64(&m.InitCounter) > 0
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

	if !m.InitFinished() {
		m.t.Fatal("Expected call to CommunicatorMock.Init")
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

	if !m.InitFinished() {
		m.t.Fatal("Expected call to CommunicatorMock.Init")
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
		ok = ok && m.InitFinished()

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

			if !m.InitFinished() {
				m.t.Error("Expected call to CommunicatorMock.Init")
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

	if !m.InitFinished() {
		return false
	}

	return true
}
