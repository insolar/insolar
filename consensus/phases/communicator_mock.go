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

	testify_assert "github.com/stretchr/testify/assert"
)

//CommunicatorMock implements github.com/insolar/insolar/consensus/phases.Communicator
type CommunicatorMock struct {
	t minimock.Tester

	ExchangePhase1Func       func(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error)
	ExchangePhase1Counter    uint64
	ExchangePhase1PreCounter uint64
	ExchangePhase1Mock       mCommunicatorMockExchangePhase1

	ExchangePhase2Func       func(p context.Context, p1 []core.Node, p2 packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error)
	ExchangePhase2Counter    uint64
	ExchangePhase2PreCounter uint64
	ExchangePhase2Mock       mCommunicatorMockExchangePhase2

	ExchangePhase3Func       func(p context.Context, p1 []core.Node, p2 packets.Phase3Packet) (r map[core.RecordRef]*packets.Phase3Packet, r1 error)
	ExchangePhase3Counter    uint64
	ExchangePhase3PreCounter uint64
	ExchangePhase3Mock       mCommunicatorMockExchangePhase3
}

//NewCommunicatorMock returns a mock for github.com/insolar/insolar/consensus/phases.Communicator
func NewCommunicatorMock(t minimock.Tester) *CommunicatorMock {
	m := &CommunicatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExchangePhase1Mock = mCommunicatorMockExchangePhase1{mock: m}
	m.ExchangePhase2Mock = mCommunicatorMockExchangePhase2{mock: m}
	m.ExchangePhase3Mock = mCommunicatorMockExchangePhase3{mock: m}

	return m
}

type mCommunicatorMockExchangePhase1 struct {
	mock             *CommunicatorMock
	mockExpectations *CommunicatorMockExchangePhase1Params
}

//CommunicatorMockExchangePhase1Params represents input parameters of the Communicator.ExchangePhase1
type CommunicatorMockExchangePhase1Params struct {
	p  context.Context
	p1 []core.Node
	p2 packets.Phase1Packet
}

//Expect sets up expected params for the Communicator.ExchangePhase1
func (m *mCommunicatorMockExchangePhase1) Expect(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) *mCommunicatorMockExchangePhase1 {
	m.mockExpectations = &CommunicatorMockExchangePhase1Params{p, p1, p2}
	return m
}

//Return sets up a mock for Communicator.ExchangePhase1 to return Return's arguments
func (m *mCommunicatorMockExchangePhase1) Return(r map[core.RecordRef]*packets.Phase1Packet, r1 error) *CommunicatorMock {
	m.mock.ExchangePhase1Func = func(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (map[core.RecordRef]*packets.Phase1Packet, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Communicator.ExchangePhase1 method
func (m *mCommunicatorMockExchangePhase1) Set(f func(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error)) *CommunicatorMock {
	m.mock.ExchangePhase1Func = f
	m.mockExpectations = nil
	return m.mock
}

//ExchangePhase1 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase1(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error) {
	atomic.AddUint64(&m.ExchangePhase1PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase1Counter, 1)

	if m.ExchangePhase1Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExchangePhase1Mock.mockExpectations, CommunicatorMockExchangePhase1Params{p, p1, p2},
			"Communicator.ExchangePhase1 got unexpected parameters")

		if m.ExchangePhase1Func == nil {

			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase1")

			return
		}
	}

	if m.ExchangePhase1Func == nil {
		m.t.Fatal("Unexpected call to CommunicatorMock.ExchangePhase1")
		return
	}

	return m.ExchangePhase1Func(p, p1, p2)
}

//ExchangePhase1MinimockCounter returns a count of CommunicatorMock.ExchangePhase1Func invocations
func (m *CommunicatorMock) ExchangePhase1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase1Counter)
}

//ExchangePhase1MinimockPreCounter returns the value of CommunicatorMock.ExchangePhase1 invocations
func (m *CommunicatorMock) ExchangePhase1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase1PreCounter)
}

type mCommunicatorMockExchangePhase2 struct {
	mock             *CommunicatorMock
	mockExpectations *CommunicatorMockExchangePhase2Params
}

//CommunicatorMockExchangePhase2Params represents input parameters of the Communicator.ExchangePhase2
type CommunicatorMockExchangePhase2Params struct {
	p  context.Context
	p1 []core.Node
	p2 packets.Phase2Packet
}

//Expect sets up expected params for the Communicator.ExchangePhase2
func (m *mCommunicatorMockExchangePhase2) Expect(p context.Context, p1 []core.Node, p2 packets.Phase2Packet) *mCommunicatorMockExchangePhase2 {
	m.mockExpectations = &CommunicatorMockExchangePhase2Params{p, p1, p2}
	return m
}

//Return sets up a mock for Communicator.ExchangePhase2 to return Return's arguments
func (m *mCommunicatorMockExchangePhase2) Return(r map[core.RecordRef]*packets.Phase2Packet, r1 error) *CommunicatorMock {
	m.mock.ExchangePhase2Func = func(p context.Context, p1 []core.Node, p2 packets.Phase2Packet) (map[core.RecordRef]*packets.Phase2Packet, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Communicator.ExchangePhase2 method
func (m *mCommunicatorMockExchangePhase2) Set(f func(p context.Context, p1 []core.Node, p2 packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error)) *CommunicatorMock {
	m.mock.ExchangePhase2Func = f
	m.mockExpectations = nil
	return m.mock
}

//ExchangePhase2 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase2(p context.Context, p1 []core.Node, p2 packets.Phase2Packet) (r map[core.RecordRef]*packets.Phase2Packet, r1 error) {
	atomic.AddUint64(&m.ExchangePhase2PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase2Counter, 1)

	if m.ExchangePhase2Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExchangePhase2Mock.mockExpectations, CommunicatorMockExchangePhase2Params{p, p1, p2},
			"Communicator.ExchangePhase2 got unexpected parameters")

		if m.ExchangePhase2Func == nil {

			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase2")

			return
		}
	}

	if m.ExchangePhase2Func == nil {
		m.t.Fatal("Unexpected call to CommunicatorMock.ExchangePhase2")
		return
	}

	return m.ExchangePhase2Func(p, p1, p2)
}

//ExchangePhase2MinimockCounter returns a count of CommunicatorMock.ExchangePhase2Func invocations
func (m *CommunicatorMock) ExchangePhase2MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase2Counter)
}

//ExchangePhase2MinimockPreCounter returns the value of CommunicatorMock.ExchangePhase2 invocations
func (m *CommunicatorMock) ExchangePhase2MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangePhase2PreCounter)
}

type mCommunicatorMockExchangePhase3 struct {
	mock             *CommunicatorMock
	mockExpectations *CommunicatorMockExchangePhase3Params
}

//CommunicatorMockExchangePhase3Params represents input parameters of the Communicator.ExchangePhase3
type CommunicatorMockExchangePhase3Params struct {
	p  context.Context
	p1 []core.Node
	p2 packets.Phase3Packet
}

//Expect sets up expected params for the Communicator.ExchangePhase3
func (m *mCommunicatorMockExchangePhase3) Expect(p context.Context, p1 []core.Node, p2 packets.Phase3Packet) *mCommunicatorMockExchangePhase3 {
	m.mockExpectations = &CommunicatorMockExchangePhase3Params{p, p1, p2}
	return m
}

//Return sets up a mock for Communicator.ExchangePhase3 to return Return's arguments
func (m *mCommunicatorMockExchangePhase3) Return(r map[core.RecordRef]*packets.Phase3Packet, r1 error) *CommunicatorMock {
	m.mock.ExchangePhase3Func = func(p context.Context, p1 []core.Node, p2 packets.Phase3Packet) (map[core.RecordRef]*packets.Phase3Packet, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Communicator.ExchangePhase3 method
func (m *mCommunicatorMockExchangePhase3) Set(f func(p context.Context, p1 []core.Node, p2 packets.Phase3Packet) (r map[core.RecordRef]*packets.Phase3Packet, r1 error)) *CommunicatorMock {
	m.mock.ExchangePhase3Func = f
	m.mockExpectations = nil
	return m.mock
}

//ExchangePhase3 implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangePhase3(p context.Context, p1 []core.Node, p2 packets.Phase3Packet) (r map[core.RecordRef]*packets.Phase3Packet, r1 error) {
	atomic.AddUint64(&m.ExchangePhase3PreCounter, 1)
	defer atomic.AddUint64(&m.ExchangePhase3Counter, 1)

	if m.ExchangePhase3Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExchangePhase3Mock.mockExpectations, CommunicatorMockExchangePhase3Params{p, p1, p2},
			"Communicator.ExchangePhase3 got unexpected parameters")

		if m.ExchangePhase3Func == nil {

			m.t.Fatal("No results are set for the CommunicatorMock.ExchangePhase3")

			return
		}
	}

	if m.ExchangePhase3Func == nil {
		m.t.Fatal("Unexpected call to CommunicatorMock.ExchangePhase3")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CommunicatorMock) ValidateCallCounters() {

	if m.ExchangePhase1Func != nil && atomic.LoadUint64(&m.ExchangePhase1Counter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase1")
	}

	if m.ExchangePhase2Func != nil && atomic.LoadUint64(&m.ExchangePhase2Counter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase2")
	}

	if m.ExchangePhase3Func != nil && atomic.LoadUint64(&m.ExchangePhase3Counter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase3")
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

	if m.ExchangePhase1Func != nil && atomic.LoadUint64(&m.ExchangePhase1Counter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase1")
	}

	if m.ExchangePhase2Func != nil && atomic.LoadUint64(&m.ExchangePhase2Counter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase2")
	}

	if m.ExchangePhase3Func != nil && atomic.LoadUint64(&m.ExchangePhase3Counter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangePhase3")
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
		ok = ok && (m.ExchangePhase1Func == nil || atomic.LoadUint64(&m.ExchangePhase1Counter) > 0)
		ok = ok && (m.ExchangePhase2Func == nil || atomic.LoadUint64(&m.ExchangePhase2Counter) > 0)
		ok = ok && (m.ExchangePhase3Func == nil || atomic.LoadUint64(&m.ExchangePhase3Counter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ExchangePhase1Func != nil && atomic.LoadUint64(&m.ExchangePhase1Counter) == 0 {
				m.t.Error("Expected call to CommunicatorMock.ExchangePhase1")
			}

			if m.ExchangePhase2Func != nil && atomic.LoadUint64(&m.ExchangePhase2Counter) == 0 {
				m.t.Error("Expected call to CommunicatorMock.ExchangePhase2")
			}

			if m.ExchangePhase3Func != nil && atomic.LoadUint64(&m.ExchangePhase3Counter) == 0 {
				m.t.Error("Expected call to CommunicatorMock.ExchangePhase3")
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

	if m.ExchangePhase1Func != nil && atomic.LoadUint64(&m.ExchangePhase1Counter) == 0 {
		return false
	}

	if m.ExchangePhase2Func != nil && atomic.LoadUint64(&m.ExchangePhase2Counter) == 0 {
		return false
	}

	if m.ExchangePhase3Func != nil && atomic.LoadUint64(&m.ExchangePhase3Counter) == 0 {
		return false
	}

	return true
}
