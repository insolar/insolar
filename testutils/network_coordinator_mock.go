package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NetworkCoordinator" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//NetworkCoordinatorMock implements github.com/insolar/insolar/core.NetworkCoordinator
type NetworkCoordinatorMock struct {
	t minimock.Tester

	GetCertFunc       func(p context.Context, p1 core.RecordRef) (r core.Certificate, r1 error)
	GetCertCounter    uint64
	GetCertPreCounter uint64
	GetCertMock       mNetworkCoordinatorMockGetCert

	SetPulseFunc       func(p context.Context, p1 core.Pulse) (r error)
	SetPulseCounter    uint64
	SetPulsePreCounter uint64
	SetPulseMock       mNetworkCoordinatorMockSetPulse

	ValidateCertFunc       func(p context.Context, p1 core.Certificate) (r bool, r1 error)
	ValidateCertCounter    uint64
	ValidateCertPreCounter uint64
	ValidateCertMock       mNetworkCoordinatorMockValidateCert

	WriteActiveNodesFunc       func(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error)
	WriteActiveNodesCounter    uint64
	WriteActiveNodesPreCounter uint64
	WriteActiveNodesMock       mNetworkCoordinatorMockWriteActiveNodes
}

//NewNetworkCoordinatorMock returns a mock for github.com/insolar/insolar/core.NetworkCoordinator
func NewNetworkCoordinatorMock(t minimock.Tester) *NetworkCoordinatorMock {
	m := &NetworkCoordinatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCertMock = mNetworkCoordinatorMockGetCert{mock: m}
	m.SetPulseMock = mNetworkCoordinatorMockSetPulse{mock: m}
	m.ValidateCertMock = mNetworkCoordinatorMockValidateCert{mock: m}
	m.WriteActiveNodesMock = mNetworkCoordinatorMockWriteActiveNodes{mock: m}

	return m
}

type mNetworkCoordinatorMockGetCert struct {
	mock             *NetworkCoordinatorMock
	mockExpectations *NetworkCoordinatorMockGetCertParams
}

//NetworkCoordinatorMockGetCertParams represents input parameters of the NetworkCoordinator.GetCert
type NetworkCoordinatorMockGetCertParams struct {
	p  context.Context
	p1 core.RecordRef
}

//Expect sets up expected params for the NetworkCoordinator.GetCert
func (m *mNetworkCoordinatorMockGetCert) Expect(p context.Context, p1 core.RecordRef) *mNetworkCoordinatorMockGetCert {
	m.mockExpectations = &NetworkCoordinatorMockGetCertParams{p, p1}
	return m
}

//Return sets up a mock for NetworkCoordinator.GetCert to return Return's arguments
func (m *mNetworkCoordinatorMockGetCert) Return(r core.Certificate, r1 error) *NetworkCoordinatorMock {
	m.mock.GetCertFunc = func(p context.Context, p1 core.RecordRef) (core.Certificate, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of NetworkCoordinator.GetCert method
func (m *mNetworkCoordinatorMockGetCert) Set(f func(p context.Context, p1 core.RecordRef) (r core.Certificate, r1 error)) *NetworkCoordinatorMock {
	m.mock.GetCertFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetCert implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) GetCert(p context.Context, p1 core.RecordRef) (r core.Certificate, r1 error) {
	atomic.AddUint64(&m.GetCertPreCounter, 1)
	defer atomic.AddUint64(&m.GetCertCounter, 1)

	if m.GetCertMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetCertMock.mockExpectations, NetworkCoordinatorMockGetCertParams{p, p1},
			"NetworkCoordinator.GetCert got unexpected parameters")

		if m.GetCertFunc == nil {

			m.t.Fatal("No results are set for the NetworkCoordinatorMock.GetCert")

			return
		}
	}

	if m.GetCertFunc == nil {
		m.t.Fatal("Unexpected call to NetworkCoordinatorMock.GetCert")
		return
	}

	return m.GetCertFunc(p, p1)
}

//GetCertMinimockCounter returns a count of NetworkCoordinatorMock.GetCertFunc invocations
func (m *NetworkCoordinatorMock) GetCertMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCertCounter)
}

//GetCertMinimockPreCounter returns the value of NetworkCoordinatorMock.GetCert invocations
func (m *NetworkCoordinatorMock) GetCertMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCertPreCounter)
}

type mNetworkCoordinatorMockSetPulse struct {
	mock             *NetworkCoordinatorMock
	mockExpectations *NetworkCoordinatorMockSetPulseParams
}

//NetworkCoordinatorMockSetPulseParams represents input parameters of the NetworkCoordinator.SetPulse
type NetworkCoordinatorMockSetPulseParams struct {
	p  context.Context
	p1 core.Pulse
}

//Expect sets up expected params for the NetworkCoordinator.SetPulse
func (m *mNetworkCoordinatorMockSetPulse) Expect(p context.Context, p1 core.Pulse) *mNetworkCoordinatorMockSetPulse {
	m.mockExpectations = &NetworkCoordinatorMockSetPulseParams{p, p1}
	return m
}

//Return sets up a mock for NetworkCoordinator.SetPulse to return Return's arguments
func (m *mNetworkCoordinatorMockSetPulse) Return(r error) *NetworkCoordinatorMock {
	m.mock.SetPulseFunc = func(p context.Context, p1 core.Pulse) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NetworkCoordinator.SetPulse method
func (m *mNetworkCoordinatorMockSetPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *NetworkCoordinatorMock {
	m.mock.SetPulseFunc = f
	m.mockExpectations = nil
	return m.mock
}

//SetPulse implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) SetPulse(p context.Context, p1 core.Pulse) (r error) {
	atomic.AddUint64(&m.SetPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetPulseCounter, 1)

	if m.SetPulseMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SetPulseMock.mockExpectations, NetworkCoordinatorMockSetPulseParams{p, p1},
			"NetworkCoordinator.SetPulse got unexpected parameters")

		if m.SetPulseFunc == nil {

			m.t.Fatal("No results are set for the NetworkCoordinatorMock.SetPulse")

			return
		}
	}

	if m.SetPulseFunc == nil {
		m.t.Fatal("Unexpected call to NetworkCoordinatorMock.SetPulse")
		return
	}

	return m.SetPulseFunc(p, p1)
}

//SetPulseMinimockCounter returns a count of NetworkCoordinatorMock.SetPulseFunc invocations
func (m *NetworkCoordinatorMock) SetPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetPulseCounter)
}

//SetPulseMinimockPreCounter returns the value of NetworkCoordinatorMock.SetPulse invocations
func (m *NetworkCoordinatorMock) SetPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPulsePreCounter)
}

type mNetworkCoordinatorMockValidateCert struct {
	mock             *NetworkCoordinatorMock
	mockExpectations *NetworkCoordinatorMockValidateCertParams
}

//NetworkCoordinatorMockValidateCertParams represents input parameters of the NetworkCoordinator.ValidateCert
type NetworkCoordinatorMockValidateCertParams struct {
	p  context.Context
	p1 core.Certificate
}

//Expect sets up expected params for the NetworkCoordinator.ValidateCert
func (m *mNetworkCoordinatorMockValidateCert) Expect(p context.Context, p1 core.Certificate) *mNetworkCoordinatorMockValidateCert {
	m.mockExpectations = &NetworkCoordinatorMockValidateCertParams{p, p1}
	return m
}

//Return sets up a mock for NetworkCoordinator.ValidateCert to return Return's arguments
func (m *mNetworkCoordinatorMockValidateCert) Return(r bool, r1 error) *NetworkCoordinatorMock {
	m.mock.ValidateCertFunc = func(p context.Context, p1 core.Certificate) (bool, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of NetworkCoordinator.ValidateCert method
func (m *mNetworkCoordinatorMockValidateCert) Set(f func(p context.Context, p1 core.Certificate) (r bool, r1 error)) *NetworkCoordinatorMock {
	m.mock.ValidateCertFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ValidateCert implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) ValidateCert(p context.Context, p1 core.Certificate) (r bool, r1 error) {
	atomic.AddUint64(&m.ValidateCertPreCounter, 1)
	defer atomic.AddUint64(&m.ValidateCertCounter, 1)

	if m.ValidateCertMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ValidateCertMock.mockExpectations, NetworkCoordinatorMockValidateCertParams{p, p1},
			"NetworkCoordinator.ValidateCert got unexpected parameters")

		if m.ValidateCertFunc == nil {

			m.t.Fatal("No results are set for the NetworkCoordinatorMock.ValidateCert")

			return
		}
	}

	if m.ValidateCertFunc == nil {
		m.t.Fatal("Unexpected call to NetworkCoordinatorMock.ValidateCert")
		return
	}

	return m.ValidateCertFunc(p, p1)
}

//ValidateCertMinimockCounter returns a count of NetworkCoordinatorMock.ValidateCertFunc invocations
func (m *NetworkCoordinatorMock) ValidateCertMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ValidateCertCounter)
}

//ValidateCertMinimockPreCounter returns the value of NetworkCoordinatorMock.ValidateCert invocations
func (m *NetworkCoordinatorMock) ValidateCertMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ValidateCertPreCounter)
}

type mNetworkCoordinatorMockWriteActiveNodes struct {
	mock             *NetworkCoordinatorMock
	mockExpectations *NetworkCoordinatorMockWriteActiveNodesParams
}

//NetworkCoordinatorMockWriteActiveNodesParams represents input parameters of the NetworkCoordinator.WriteActiveNodes
type NetworkCoordinatorMockWriteActiveNodesParams struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []core.Node
}

//Expect sets up expected params for the NetworkCoordinator.WriteActiveNodes
func (m *mNetworkCoordinatorMockWriteActiveNodes) Expect(p context.Context, p1 core.PulseNumber, p2 []core.Node) *mNetworkCoordinatorMockWriteActiveNodes {
	m.mockExpectations = &NetworkCoordinatorMockWriteActiveNodesParams{p, p1, p2}
	return m
}

//Return sets up a mock for NetworkCoordinator.WriteActiveNodes to return Return's arguments
func (m *mNetworkCoordinatorMockWriteActiveNodes) Return(r error) *NetworkCoordinatorMock {
	m.mock.WriteActiveNodesFunc = func(p context.Context, p1 core.PulseNumber, p2 []core.Node) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of NetworkCoordinator.WriteActiveNodes method
func (m *mNetworkCoordinatorMockWriteActiveNodes) Set(f func(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error)) *NetworkCoordinatorMock {
	m.mock.WriteActiveNodesFunc = f
	m.mockExpectations = nil
	return m.mock
}

//WriteActiveNodes implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) WriteActiveNodes(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error) {
	atomic.AddUint64(&m.WriteActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.WriteActiveNodesCounter, 1)

	if m.WriteActiveNodesMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.WriteActiveNodesMock.mockExpectations, NetworkCoordinatorMockWriteActiveNodesParams{p, p1, p2},
			"NetworkCoordinator.WriteActiveNodes got unexpected parameters")

		if m.WriteActiveNodesFunc == nil {

			m.t.Fatal("No results are set for the NetworkCoordinatorMock.WriteActiveNodes")

			return
		}
	}

	if m.WriteActiveNodesFunc == nil {
		m.t.Fatal("Unexpected call to NetworkCoordinatorMock.WriteActiveNodes")
		return
	}

	return m.WriteActiveNodesFunc(p, p1, p2)
}

//WriteActiveNodesMinimockCounter returns a count of NetworkCoordinatorMock.WriteActiveNodesFunc invocations
func (m *NetworkCoordinatorMock) WriteActiveNodesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteActiveNodesCounter)
}

//WriteActiveNodesMinimockPreCounter returns the value of NetworkCoordinatorMock.WriteActiveNodes invocations
func (m *NetworkCoordinatorMock) WriteActiveNodesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteActiveNodesPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkCoordinatorMock) ValidateCallCounters() {

	if m.GetCertFunc != nil && atomic.LoadUint64(&m.GetCertCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.GetCert")
	}

	if m.SetPulseFunc != nil && atomic.LoadUint64(&m.SetPulseCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.SetPulse")
	}

	if m.ValidateCertFunc != nil && atomic.LoadUint64(&m.ValidateCertCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.ValidateCert")
	}

	if m.WriteActiveNodesFunc != nil && atomic.LoadUint64(&m.WriteActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.WriteActiveNodes")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkCoordinatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NetworkCoordinatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NetworkCoordinatorMock) MinimockFinish() {

	if m.GetCertFunc != nil && atomic.LoadUint64(&m.GetCertCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.GetCert")
	}

	if m.SetPulseFunc != nil && atomic.LoadUint64(&m.SetPulseCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.SetPulse")
	}

	if m.ValidateCertFunc != nil && atomic.LoadUint64(&m.ValidateCertCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.ValidateCert")
	}

	if m.WriteActiveNodesFunc != nil && atomic.LoadUint64(&m.WriteActiveNodesCounter) == 0 {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.WriteActiveNodes")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NetworkCoordinatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NetworkCoordinatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetCertFunc == nil || atomic.LoadUint64(&m.GetCertCounter) > 0)
		ok = ok && (m.SetPulseFunc == nil || atomic.LoadUint64(&m.SetPulseCounter) > 0)
		ok = ok && (m.ValidateCertFunc == nil || atomic.LoadUint64(&m.ValidateCertCounter) > 0)
		ok = ok && (m.WriteActiveNodesFunc == nil || atomic.LoadUint64(&m.WriteActiveNodesCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetCertFunc != nil && atomic.LoadUint64(&m.GetCertCounter) == 0 {
				m.t.Error("Expected call to NetworkCoordinatorMock.GetCert")
			}

			if m.SetPulseFunc != nil && atomic.LoadUint64(&m.SetPulseCounter) == 0 {
				m.t.Error("Expected call to NetworkCoordinatorMock.SetPulse")
			}

			if m.ValidateCertFunc != nil && atomic.LoadUint64(&m.ValidateCertCounter) == 0 {
				m.t.Error("Expected call to NetworkCoordinatorMock.ValidateCert")
			}

			if m.WriteActiveNodesFunc != nil && atomic.LoadUint64(&m.WriteActiveNodesCounter) == 0 {
				m.t.Error("Expected call to NetworkCoordinatorMock.WriteActiveNodes")
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
func (m *NetworkCoordinatorMock) AllMocksCalled() bool {

	if m.GetCertFunc != nil && atomic.LoadUint64(&m.GetCertCounter) == 0 {
		return false
	}

	if m.SetPulseFunc != nil && atomic.LoadUint64(&m.SetPulseCounter) == 0 {
		return false
	}

	if m.ValidateCertFunc != nil && atomic.LoadUint64(&m.ValidateCertCounter) == 0 {
		return false
	}

	if m.WriteActiveNodesFunc != nil && atomic.LoadUint64(&m.WriteActiveNodesCounter) == 0 {
		return false
	}

	return true
}
