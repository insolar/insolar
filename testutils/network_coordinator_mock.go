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

	GetCertFunc       func(p context.Context, p1 *core.RecordRef) (r core.Certificate, r1 error)
	GetCertCounter    uint64
	GetCertPreCounter uint64
	GetCertMock       mNetworkCoordinatorMockGetCert

	IsStartedFunc       func() (r bool)
	IsStartedCounter    uint64
	IsStartedPreCounter uint64
	IsStartedMock       mNetworkCoordinatorMockIsStarted

	SetPulseFunc       func(p context.Context, p1 core.Pulse) (r error)
	SetPulseCounter    uint64
	SetPulsePreCounter uint64
	SetPulseMock       mNetworkCoordinatorMockSetPulse

	ValidateCertFunc       func(p context.Context, p1 core.AuthorizationCertificate) (r bool, r1 error)
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
	m.IsStartedMock = mNetworkCoordinatorMockIsStarted{mock: m}
	m.SetPulseMock = mNetworkCoordinatorMockSetPulse{mock: m}
	m.ValidateCertMock = mNetworkCoordinatorMockValidateCert{mock: m}
	m.WriteActiveNodesMock = mNetworkCoordinatorMockWriteActiveNodes{mock: m}

	return m
}

type mNetworkCoordinatorMockGetCert struct {
	mock              *NetworkCoordinatorMock
	mainExpectation   *NetworkCoordinatorMockGetCertExpectation
	expectationSeries []*NetworkCoordinatorMockGetCertExpectation
}

type NetworkCoordinatorMockGetCertExpectation struct {
	input  *NetworkCoordinatorMockGetCertInput
	result *NetworkCoordinatorMockGetCertResult
}

type NetworkCoordinatorMockGetCertInput struct {
	p  context.Context
	p1 *core.RecordRef
}

type NetworkCoordinatorMockGetCertResult struct {
	r  core.Certificate
	r1 error
}

//Expect specifies that invocation of NetworkCoordinator.GetCert is expected from 1 to Infinity times
func (m *mNetworkCoordinatorMockGetCert) Expect(p context.Context, p1 *core.RecordRef) *mNetworkCoordinatorMockGetCert {
	m.mock.GetCertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockGetCertExpectation{}
	}
	m.mainExpectation.input = &NetworkCoordinatorMockGetCertInput{p, p1}
	return m
}

//Return specifies results of invocation of NetworkCoordinator.GetCert
func (m *mNetworkCoordinatorMockGetCert) Return(r core.Certificate, r1 error) *NetworkCoordinatorMock {
	m.mock.GetCertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockGetCertExpectation{}
	}
	m.mainExpectation.result = &NetworkCoordinatorMockGetCertResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkCoordinator.GetCert is expected once
func (m *mNetworkCoordinatorMockGetCert) ExpectOnce(p context.Context, p1 *core.RecordRef) *NetworkCoordinatorMockGetCertExpectation {
	m.mock.GetCertFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkCoordinatorMockGetCertExpectation{}
	expectation.input = &NetworkCoordinatorMockGetCertInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkCoordinatorMockGetCertExpectation) Return(r core.Certificate, r1 error) {
	e.result = &NetworkCoordinatorMockGetCertResult{r, r1}
}

//Set uses given function f as a mock of NetworkCoordinator.GetCert method
func (m *mNetworkCoordinatorMockGetCert) Set(f func(p context.Context, p1 *core.RecordRef) (r core.Certificate, r1 error)) *NetworkCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCertFunc = f
	return m.mock
}

//GetCert implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) GetCert(p context.Context, p1 *core.RecordRef) (r core.Certificate, r1 error) {
	counter := atomic.AddUint64(&m.GetCertPreCounter, 1)
	defer atomic.AddUint64(&m.GetCertCounter, 1)

	if len(m.GetCertMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCertMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.GetCert. %v %v", p, p1)
			return
		}

		input := m.GetCertMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkCoordinatorMockGetCertInput{p, p1}, "NetworkCoordinator.GetCert got unexpected parameters")

		result := m.GetCertMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.GetCert")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCertMock.mainExpectation != nil {

		input := m.GetCertMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkCoordinatorMockGetCertInput{p, p1}, "NetworkCoordinator.GetCert got unexpected parameters")
		}

		result := m.GetCertMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.GetCert")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCertFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.GetCert. %v %v", p, p1)
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

//GetCertFinished returns true if mock invocations count is ok
func (m *NetworkCoordinatorMock) GetCertFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCertMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCertCounter) == uint64(len(m.GetCertMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCertMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCertCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCertFunc != nil {
		return atomic.LoadUint64(&m.GetCertCounter) > 0
	}

	return true
}

type mNetworkCoordinatorMockIsStarted struct {
	mock              *NetworkCoordinatorMock
	mainExpectation   *NetworkCoordinatorMockIsStartedExpectation
	expectationSeries []*NetworkCoordinatorMockIsStartedExpectation
}

type NetworkCoordinatorMockIsStartedExpectation struct {
	result *NetworkCoordinatorMockIsStartedResult
}

type NetworkCoordinatorMockIsStartedResult struct {
	r bool
}

//Expect specifies that invocation of NetworkCoordinator.IsStarted is expected from 1 to Infinity times
func (m *mNetworkCoordinatorMockIsStarted) Expect() *mNetworkCoordinatorMockIsStarted {
	m.mock.IsStartedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockIsStartedExpectation{}
	}

	return m
}

//Return specifies results of invocation of NetworkCoordinator.IsStarted
func (m *mNetworkCoordinatorMockIsStarted) Return(r bool) *NetworkCoordinatorMock {
	m.mock.IsStartedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockIsStartedExpectation{}
	}
	m.mainExpectation.result = &NetworkCoordinatorMockIsStartedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkCoordinator.IsStarted is expected once
func (m *mNetworkCoordinatorMockIsStarted) ExpectOnce() *NetworkCoordinatorMockIsStartedExpectation {
	m.mock.IsStartedFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkCoordinatorMockIsStartedExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkCoordinatorMockIsStartedExpectation) Return(r bool) {
	e.result = &NetworkCoordinatorMockIsStartedResult{r}
}

//Set uses given function f as a mock of NetworkCoordinator.IsStarted method
func (m *mNetworkCoordinatorMockIsStarted) Set(f func() (r bool)) *NetworkCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsStartedFunc = f
	return m.mock
}

//IsStarted implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) IsStarted() (r bool) {
	counter := atomic.AddUint64(&m.IsStartedPreCounter, 1)
	defer atomic.AddUint64(&m.IsStartedCounter, 1)

	if len(m.IsStartedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsStartedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.IsStarted.")
			return
		}

		result := m.IsStartedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.IsStarted")
			return
		}

		r = result.r

		return
	}

	if m.IsStartedMock.mainExpectation != nil {

		result := m.IsStartedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.IsStarted")
		}

		r = result.r

		return
	}

	if m.IsStartedFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.IsStarted.")
		return
	}

	return m.IsStartedFunc()
}

//IsStartedMinimockCounter returns a count of NetworkCoordinatorMock.IsStartedFunc invocations
func (m *NetworkCoordinatorMock) IsStartedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsStartedCounter)
}

//IsStartedMinimockPreCounter returns the value of NetworkCoordinatorMock.IsStarted invocations
func (m *NetworkCoordinatorMock) IsStartedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsStartedPreCounter)
}

//IsStartedFinished returns true if mock invocations count is ok
func (m *NetworkCoordinatorMock) IsStartedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsStartedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsStartedCounter) == uint64(len(m.IsStartedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsStartedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsStartedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsStartedFunc != nil {
		return atomic.LoadUint64(&m.IsStartedCounter) > 0
	}

	return true
}

type mNetworkCoordinatorMockSetPulse struct {
	mock              *NetworkCoordinatorMock
	mainExpectation   *NetworkCoordinatorMockSetPulseExpectation
	expectationSeries []*NetworkCoordinatorMockSetPulseExpectation
}

type NetworkCoordinatorMockSetPulseExpectation struct {
	input  *NetworkCoordinatorMockSetPulseInput
	result *NetworkCoordinatorMockSetPulseResult
}

type NetworkCoordinatorMockSetPulseInput struct {
	p  context.Context
	p1 core.Pulse
}

type NetworkCoordinatorMockSetPulseResult struct {
	r error
}

//Expect specifies that invocation of NetworkCoordinator.SetPulse is expected from 1 to Infinity times
func (m *mNetworkCoordinatorMockSetPulse) Expect(p context.Context, p1 core.Pulse) *mNetworkCoordinatorMockSetPulse {
	m.mock.SetPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockSetPulseExpectation{}
	}
	m.mainExpectation.input = &NetworkCoordinatorMockSetPulseInput{p, p1}
	return m
}

//Return specifies results of invocation of NetworkCoordinator.SetPulse
func (m *mNetworkCoordinatorMockSetPulse) Return(r error) *NetworkCoordinatorMock {
	m.mock.SetPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockSetPulseExpectation{}
	}
	m.mainExpectation.result = &NetworkCoordinatorMockSetPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkCoordinator.SetPulse is expected once
func (m *mNetworkCoordinatorMockSetPulse) ExpectOnce(p context.Context, p1 core.Pulse) *NetworkCoordinatorMockSetPulseExpectation {
	m.mock.SetPulseFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkCoordinatorMockSetPulseExpectation{}
	expectation.input = &NetworkCoordinatorMockSetPulseInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkCoordinatorMockSetPulseExpectation) Return(r error) {
	e.result = &NetworkCoordinatorMockSetPulseResult{r}
}

//Set uses given function f as a mock of NetworkCoordinator.SetPulse method
func (m *mNetworkCoordinatorMockSetPulse) Set(f func(p context.Context, p1 core.Pulse) (r error)) *NetworkCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetPulseFunc = f
	return m.mock
}

//SetPulse implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) SetPulse(p context.Context, p1 core.Pulse) (r error) {
	counter := atomic.AddUint64(&m.SetPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetPulseCounter, 1)

	if len(m.SetPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.SetPulse. %v %v", p, p1)
			return
		}

		input := m.SetPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkCoordinatorMockSetPulseInput{p, p1}, "NetworkCoordinator.SetPulse got unexpected parameters")

		result := m.SetPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.SetPulse")
			return
		}

		r = result.r

		return
	}

	if m.SetPulseMock.mainExpectation != nil {

		input := m.SetPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkCoordinatorMockSetPulseInput{p, p1}, "NetworkCoordinator.SetPulse got unexpected parameters")
		}

		result := m.SetPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.SetPulse")
		}

		r = result.r

		return
	}

	if m.SetPulseFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.SetPulse. %v %v", p, p1)
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

//SetPulseFinished returns true if mock invocations count is ok
func (m *NetworkCoordinatorMock) SetPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetPulseCounter) == uint64(len(m.SetPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetPulseFunc != nil {
		return atomic.LoadUint64(&m.SetPulseCounter) > 0
	}

	return true
}

type mNetworkCoordinatorMockValidateCert struct {
	mock              *NetworkCoordinatorMock
	mainExpectation   *NetworkCoordinatorMockValidateCertExpectation
	expectationSeries []*NetworkCoordinatorMockValidateCertExpectation
}

type NetworkCoordinatorMockValidateCertExpectation struct {
	input  *NetworkCoordinatorMockValidateCertInput
	result *NetworkCoordinatorMockValidateCertResult
}

type NetworkCoordinatorMockValidateCertInput struct {
	p  context.Context
	p1 core.AuthorizationCertificate
}

type NetworkCoordinatorMockValidateCertResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of NetworkCoordinator.ValidateCert is expected from 1 to Infinity times
func (m *mNetworkCoordinatorMockValidateCert) Expect(p context.Context, p1 core.AuthorizationCertificate) *mNetworkCoordinatorMockValidateCert {
	m.mock.ValidateCertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockValidateCertExpectation{}
	}
	m.mainExpectation.input = &NetworkCoordinatorMockValidateCertInput{p, p1}
	return m
}

//Return specifies results of invocation of NetworkCoordinator.ValidateCert
func (m *mNetworkCoordinatorMockValidateCert) Return(r bool, r1 error) *NetworkCoordinatorMock {
	m.mock.ValidateCertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockValidateCertExpectation{}
	}
	m.mainExpectation.result = &NetworkCoordinatorMockValidateCertResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkCoordinator.ValidateCert is expected once
func (m *mNetworkCoordinatorMockValidateCert) ExpectOnce(p context.Context, p1 core.AuthorizationCertificate) *NetworkCoordinatorMockValidateCertExpectation {
	m.mock.ValidateCertFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkCoordinatorMockValidateCertExpectation{}
	expectation.input = &NetworkCoordinatorMockValidateCertInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkCoordinatorMockValidateCertExpectation) Return(r bool, r1 error) {
	e.result = &NetworkCoordinatorMockValidateCertResult{r, r1}
}

//Set uses given function f as a mock of NetworkCoordinator.ValidateCert method
func (m *mNetworkCoordinatorMockValidateCert) Set(f func(p context.Context, p1 core.AuthorizationCertificate) (r bool, r1 error)) *NetworkCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ValidateCertFunc = f
	return m.mock
}

//ValidateCert implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) ValidateCert(p context.Context, p1 core.AuthorizationCertificate) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.ValidateCertPreCounter, 1)
	defer atomic.AddUint64(&m.ValidateCertCounter, 1)

	if len(m.ValidateCertMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ValidateCertMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.ValidateCert. %v %v", p, p1)
			return
		}

		input := m.ValidateCertMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkCoordinatorMockValidateCertInput{p, p1}, "NetworkCoordinator.ValidateCert got unexpected parameters")

		result := m.ValidateCertMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.ValidateCert")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateCertMock.mainExpectation != nil {

		input := m.ValidateCertMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkCoordinatorMockValidateCertInput{p, p1}, "NetworkCoordinator.ValidateCert got unexpected parameters")
		}

		result := m.ValidateCertMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.ValidateCert")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateCertFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.ValidateCert. %v %v", p, p1)
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

//ValidateCertFinished returns true if mock invocations count is ok
func (m *NetworkCoordinatorMock) ValidateCertFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ValidateCertMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ValidateCertCounter) == uint64(len(m.ValidateCertMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ValidateCertMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ValidateCertCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ValidateCertFunc != nil {
		return atomic.LoadUint64(&m.ValidateCertCounter) > 0
	}

	return true
}

type mNetworkCoordinatorMockWriteActiveNodes struct {
	mock              *NetworkCoordinatorMock
	mainExpectation   *NetworkCoordinatorMockWriteActiveNodesExpectation
	expectationSeries []*NetworkCoordinatorMockWriteActiveNodesExpectation
}

type NetworkCoordinatorMockWriteActiveNodesExpectation struct {
	input  *NetworkCoordinatorMockWriteActiveNodesInput
	result *NetworkCoordinatorMockWriteActiveNodesResult
}

type NetworkCoordinatorMockWriteActiveNodesInput struct {
	p  context.Context
	p1 core.PulseNumber
	p2 []core.Node
}

type NetworkCoordinatorMockWriteActiveNodesResult struct {
	r error
}

//Expect specifies that invocation of NetworkCoordinator.WriteActiveNodes is expected from 1 to Infinity times
func (m *mNetworkCoordinatorMockWriteActiveNodes) Expect(p context.Context, p1 core.PulseNumber, p2 []core.Node) *mNetworkCoordinatorMockWriteActiveNodes {
	m.mock.WriteActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockWriteActiveNodesExpectation{}
	}
	m.mainExpectation.input = &NetworkCoordinatorMockWriteActiveNodesInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of NetworkCoordinator.WriteActiveNodes
func (m *mNetworkCoordinatorMockWriteActiveNodes) Return(r error) *NetworkCoordinatorMock {
	m.mock.WriteActiveNodesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NetworkCoordinatorMockWriteActiveNodesExpectation{}
	}
	m.mainExpectation.result = &NetworkCoordinatorMockWriteActiveNodesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NetworkCoordinator.WriteActiveNodes is expected once
func (m *mNetworkCoordinatorMockWriteActiveNodes) ExpectOnce(p context.Context, p1 core.PulseNumber, p2 []core.Node) *NetworkCoordinatorMockWriteActiveNodesExpectation {
	m.mock.WriteActiveNodesFunc = nil
	m.mainExpectation = nil

	expectation := &NetworkCoordinatorMockWriteActiveNodesExpectation{}
	expectation.input = &NetworkCoordinatorMockWriteActiveNodesInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NetworkCoordinatorMockWriteActiveNodesExpectation) Return(r error) {
	e.result = &NetworkCoordinatorMockWriteActiveNodesResult{r}
}

//Set uses given function f as a mock of NetworkCoordinator.WriteActiveNodes method
func (m *mNetworkCoordinatorMockWriteActiveNodes) Set(f func(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error)) *NetworkCoordinatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteActiveNodesFunc = f
	return m.mock
}

//WriteActiveNodes implements github.com/insolar/insolar/core.NetworkCoordinator interface
func (m *NetworkCoordinatorMock) WriteActiveNodes(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error) {
	counter := atomic.AddUint64(&m.WriteActiveNodesPreCounter, 1)
	defer atomic.AddUint64(&m.WriteActiveNodesCounter, 1)

	if len(m.WriteActiveNodesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteActiveNodesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.WriteActiveNodes. %v %v %v", p, p1, p2)
			return
		}

		input := m.WriteActiveNodesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NetworkCoordinatorMockWriteActiveNodesInput{p, p1, p2}, "NetworkCoordinator.WriteActiveNodes got unexpected parameters")

		result := m.WriteActiveNodesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.WriteActiveNodes")
			return
		}

		r = result.r

		return
	}

	if m.WriteActiveNodesMock.mainExpectation != nil {

		input := m.WriteActiveNodesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NetworkCoordinatorMockWriteActiveNodesInput{p, p1, p2}, "NetworkCoordinator.WriteActiveNodes got unexpected parameters")
		}

		result := m.WriteActiveNodesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NetworkCoordinatorMock.WriteActiveNodes")
		}

		r = result.r

		return
	}

	if m.WriteActiveNodesFunc == nil {
		m.t.Fatalf("Unexpected call to NetworkCoordinatorMock.WriteActiveNodes. %v %v %v", p, p1, p2)
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

//WriteActiveNodesFinished returns true if mock invocations count is ok
func (m *NetworkCoordinatorMock) WriteActiveNodesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.WriteActiveNodesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.WriteActiveNodesCounter) == uint64(len(m.WriteActiveNodesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.WriteActiveNodesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.WriteActiveNodesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.WriteActiveNodesFunc != nil {
		return atomic.LoadUint64(&m.WriteActiveNodesCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NetworkCoordinatorMock) ValidateCallCounters() {

	if !m.GetCertFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.GetCert")
	}

	if !m.IsStartedFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.IsStarted")
	}

	if !m.SetPulseFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.SetPulse")
	}

	if !m.ValidateCertFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.ValidateCert")
	}

	if !m.WriteActiveNodesFinished() {
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

	if !m.GetCertFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.GetCert")
	}

	if !m.IsStartedFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.IsStarted")
	}

	if !m.SetPulseFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.SetPulse")
	}

	if !m.ValidateCertFinished() {
		m.t.Fatal("Expected call to NetworkCoordinatorMock.ValidateCert")
	}

	if !m.WriteActiveNodesFinished() {
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
		ok = ok && m.GetCertFinished()
		ok = ok && m.IsStartedFinished()
		ok = ok && m.SetPulseFinished()
		ok = ok && m.ValidateCertFinished()
		ok = ok && m.WriteActiveNodesFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetCertFinished() {
				m.t.Error("Expected call to NetworkCoordinatorMock.GetCert")
			}

			if !m.IsStartedFinished() {
				m.t.Error("Expected call to NetworkCoordinatorMock.IsStarted")
			}

			if !m.SetPulseFinished() {
				m.t.Error("Expected call to NetworkCoordinatorMock.SetPulse")
			}

			if !m.ValidateCertFinished() {
				m.t.Error("Expected call to NetworkCoordinatorMock.ValidateCert")
			}

			if !m.WriteActiveNodesFinished() {
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

	if !m.GetCertFinished() {
		return false
	}

	if !m.IsStartedFinished() {
		return false
	}

	if !m.SetPulseFinished() {
		return false
	}

	if !m.ValidateCertFinished() {
		return false
	}

	if !m.WriteActiveNodesFinished() {
		return false
	}

	return true
}
