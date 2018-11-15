package network

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

	ExchangeDataFunc       func(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error)
	ExchangeDataCounter    uint64
	ExchangeDataPreCounter uint64
	ExchangeDataMock       mCommunicatorMockExchangeData
}

//NewCommunicatorMock returns a mock for github.com/insolar/insolar/consensus/phases.Communicator
func NewCommunicatorMock(t minimock.Tester) *CommunicatorMock {
	m := &CommunicatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExchangeDataMock = mCommunicatorMockExchangeData{mock: m}

	return m
}

type mCommunicatorMockExchangeData struct {
	mock             *CommunicatorMock
	mockExpectations *CommunicatorMockExchangeDataParams
}

//CommunicatorMockExchangeDataParams represents input parameters of the Communicator.ExchangeData
type CommunicatorMockExchangeDataParams struct {
	p  context.Context
	p1 []core.Node
	p2 packets.Phase1Packet
}

//Expect sets up expected params for the Communicator.ExchangeData
func (m *mCommunicatorMockExchangeData) Expect(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) *mCommunicatorMockExchangeData {
	m.mockExpectations = &CommunicatorMockExchangeDataParams{p, p1, p2}
	return m
}

//Return sets up a mock for Communicator.ExchangeData to return Return's arguments
func (m *mCommunicatorMockExchangeData) Return(r map[core.RecordRef]*packets.Phase1Packet, r1 error) *CommunicatorMock {
	m.mock.ExchangeDataFunc = func(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (map[core.RecordRef]*packets.Phase1Packet, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Communicator.ExchangeData method
func (m *mCommunicatorMockExchangeData) Set(f func(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error)) *CommunicatorMock {
	m.mock.ExchangeDataFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ExchangeData implements github.com/insolar/insolar/consensus/phases.Communicator interface
func (m *CommunicatorMock) ExchangeData(p context.Context, p1 []core.Node, p2 packets.Phase1Packet) (r map[core.RecordRef]*packets.Phase1Packet, r1 error) {
	atomic.AddUint64(&m.ExchangeDataPreCounter, 1)
	defer atomic.AddUint64(&m.ExchangeDataCounter, 1)

	if m.ExchangeDataMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExchangeDataMock.mockExpectations, CommunicatorMockExchangeDataParams{p, p1, p2},
			"Communicator.ExchangeData got unexpected parameters")

		if m.ExchangeDataFunc == nil {

			m.t.Fatal("No results are set for the CommunicatorMock.ExchangeData")

			return
		}
	}

	if m.ExchangeDataFunc == nil {
		m.t.Fatal("Unexpected call to CommunicatorMock.ExchangeData")
		return
	}

	return m.ExchangeDataFunc(p, p1, p2)
}

//ExchangeDataMinimockCounter returns a count of CommunicatorMock.ExchangeDataFunc invocations
func (m *CommunicatorMock) ExchangeDataMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangeDataCounter)
}

//ExchangeDataMinimockPreCounter returns the value of CommunicatorMock.ExchangeData invocations
func (m *CommunicatorMock) ExchangeDataMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExchangeDataPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CommunicatorMock) ValidateCallCounters() {

	if m.ExchangeDataFunc != nil && atomic.LoadUint64(&m.ExchangeDataCounter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangeData")
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

	if m.ExchangeDataFunc != nil && atomic.LoadUint64(&m.ExchangeDataCounter) == 0 {
		m.t.Fatal("Expected call to CommunicatorMock.ExchangeData")
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
		ok = ok && (m.ExchangeDataFunc == nil || atomic.LoadUint64(&m.ExchangeDataCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ExchangeDataFunc != nil && atomic.LoadUint64(&m.ExchangeDataCounter) == 0 {
				m.t.Error("Expected call to CommunicatorMock.ExchangeData")
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

	if m.ExchangeDataFunc != nil && atomic.LoadUint64(&m.ExchangeDataCounter) == 0 {
		return false
	}

	return true
}
