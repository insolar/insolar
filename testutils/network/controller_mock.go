package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Controller" can be found in github.com/insolar/insolar/network
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	network "github.com/insolar/insolar/network"

	testify_assert "github.com/stretchr/testify/assert"
)

//ControllerMock implements github.com/insolar/insolar/network.Controller
type ControllerMock struct {
	t minimock.Tester

	AuthenticateToDiscoveryNodeFunc       func(p context.Context, p1 insolar.DiscoveryNode) (r error)
	AuthenticateToDiscoveryNodeCounter    uint64
	AuthenticateToDiscoveryNodePreCounter uint64
	AuthenticateToDiscoveryNodeMock       mControllerMockAuthenticateToDiscoveryNode

	BootstrapFunc       func(p context.Context) (r *network.BootstrapResult, r1 error)
	BootstrapCounter    uint64
	BootstrapPreCounter uint64
	BootstrapMock       mControllerMockBootstrap

	GetLastIgnoredPulseFunc       func() (r insolar.PulseNumber)
	GetLastIgnoredPulseCounter    uint64
	GetLastIgnoredPulsePreCounter uint64
	GetLastIgnoredPulseMock       mControllerMockGetLastIgnoredPulse

	InitFunc       func(p context.Context) (r error)
	InitCounter    uint64
	InitPreCounter uint64
	InitMock       mControllerMockInit

	RemoteProcedureRegisterFunc       func(p string, p1 insolar.RemoteProcedure)
	RemoteProcedureRegisterCounter    uint64
	RemoteProcedureRegisterPreCounter uint64
	RemoteProcedureRegisterMock       mControllerMockRemoteProcedureRegister

	SendBytesFunc       func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error)
	SendBytesCounter    uint64
	SendBytesPreCounter uint64
	SendBytesMock       mControllerMockSendBytes

	SendCascadeMessageFunc       func(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error)
	SendCascadeMessageCounter    uint64
	SendCascadeMessagePreCounter uint64
	SendCascadeMessageMock       mControllerMockSendCascadeMessage

	SendMessageFunc       func(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error)
	SendMessageCounter    uint64
	SendMessagePreCounter uint64
	SendMessageMock       mControllerMockSendMessage

	SetLastIgnoredPulseFunc       func(p insolar.PulseNumber)
	SetLastIgnoredPulseCounter    uint64
	SetLastIgnoredPulsePreCounter uint64
	SetLastIgnoredPulseMock       mControllerMockSetLastIgnoredPulse
}

//NewControllerMock returns a mock for github.com/insolar/insolar/network.Controller
func NewControllerMock(t minimock.Tester) *ControllerMock {
	m := &ControllerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AuthenticateToDiscoveryNodeMock = mControllerMockAuthenticateToDiscoveryNode{mock: m}
	m.BootstrapMock = mControllerMockBootstrap{mock: m}
	m.GetLastIgnoredPulseMock = mControllerMockGetLastIgnoredPulse{mock: m}
	m.InitMock = mControllerMockInit{mock: m}
	m.RemoteProcedureRegisterMock = mControllerMockRemoteProcedureRegister{mock: m}
	m.SendBytesMock = mControllerMockSendBytes{mock: m}
	m.SendCascadeMessageMock = mControllerMockSendCascadeMessage{mock: m}
	m.SendMessageMock = mControllerMockSendMessage{mock: m}
	m.SetLastIgnoredPulseMock = mControllerMockSetLastIgnoredPulse{mock: m}

	return m
}

type mControllerMockAuthenticateToDiscoveryNode struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockAuthenticateToDiscoveryNodeExpectation
	expectationSeries []*ControllerMockAuthenticateToDiscoveryNodeExpectation
}

type ControllerMockAuthenticateToDiscoveryNodeExpectation struct {
	input  *ControllerMockAuthenticateToDiscoveryNodeInput
	result *ControllerMockAuthenticateToDiscoveryNodeResult
}

type ControllerMockAuthenticateToDiscoveryNodeInput struct {
	p  context.Context
	p1 insolar.DiscoveryNode
}

type ControllerMockAuthenticateToDiscoveryNodeResult struct {
	r error
}

//Expect specifies that invocation of Controller.AuthenticateToDiscoveryNode is expected from 1 to Infinity times
func (m *mControllerMockAuthenticateToDiscoveryNode) Expect(p context.Context, p1 insolar.DiscoveryNode) *mControllerMockAuthenticateToDiscoveryNode {
	m.mock.AuthenticateToDiscoveryNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockAuthenticateToDiscoveryNodeExpectation{}
	}
	m.mainExpectation.input = &ControllerMockAuthenticateToDiscoveryNodeInput{p, p1}
	return m
}

//Return specifies results of invocation of Controller.AuthenticateToDiscoveryNode
func (m *mControllerMockAuthenticateToDiscoveryNode) Return(r error) *ControllerMock {
	m.mock.AuthenticateToDiscoveryNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockAuthenticateToDiscoveryNodeExpectation{}
	}
	m.mainExpectation.result = &ControllerMockAuthenticateToDiscoveryNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Controller.AuthenticateToDiscoveryNode is expected once
func (m *mControllerMockAuthenticateToDiscoveryNode) ExpectOnce(p context.Context, p1 insolar.DiscoveryNode) *ControllerMockAuthenticateToDiscoveryNodeExpectation {
	m.mock.AuthenticateToDiscoveryNodeFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockAuthenticateToDiscoveryNodeExpectation{}
	expectation.input = &ControllerMockAuthenticateToDiscoveryNodeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ControllerMockAuthenticateToDiscoveryNodeExpectation) Return(r error) {
	e.result = &ControllerMockAuthenticateToDiscoveryNodeResult{r}
}

//Set uses given function f as a mock of Controller.AuthenticateToDiscoveryNode method
func (m *mControllerMockAuthenticateToDiscoveryNode) Set(f func(p context.Context, p1 insolar.DiscoveryNode) (r error)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AuthenticateToDiscoveryNodeFunc = f
	return m.mock
}

//AuthenticateToDiscoveryNode implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) AuthenticateToDiscoveryNode(p context.Context, p1 insolar.DiscoveryNode) (r error) {
	counter := atomic.AddUint64(&m.AuthenticateToDiscoveryNodePreCounter, 1)
	defer atomic.AddUint64(&m.AuthenticateToDiscoveryNodeCounter, 1)

	if len(m.AuthenticateToDiscoveryNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AuthenticateToDiscoveryNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.AuthenticateToDiscoveryNode. %v %v", p, p1)
			return
		}

		input := m.AuthenticateToDiscoveryNodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockAuthenticateToDiscoveryNodeInput{p, p1}, "Controller.AuthenticateToDiscoveryNode got unexpected parameters")

		result := m.AuthenticateToDiscoveryNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.AuthenticateToDiscoveryNode")
			return
		}

		r = result.r

		return
	}

	if m.AuthenticateToDiscoveryNodeMock.mainExpectation != nil {

		input := m.AuthenticateToDiscoveryNodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockAuthenticateToDiscoveryNodeInput{p, p1}, "Controller.AuthenticateToDiscoveryNode got unexpected parameters")
		}

		result := m.AuthenticateToDiscoveryNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.AuthenticateToDiscoveryNode")
		}

		r = result.r

		return
	}

	if m.AuthenticateToDiscoveryNodeFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.AuthenticateToDiscoveryNode. %v %v", p, p1)
		return
	}

	return m.AuthenticateToDiscoveryNodeFunc(p, p1)
}

//AuthenticateToDiscoveryNodeMinimockCounter returns a count of ControllerMock.AuthenticateToDiscoveryNodeFunc invocations
func (m *ControllerMock) AuthenticateToDiscoveryNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AuthenticateToDiscoveryNodeCounter)
}

//AuthenticateToDiscoveryNodeMinimockPreCounter returns the value of ControllerMock.AuthenticateToDiscoveryNode invocations
func (m *ControllerMock) AuthenticateToDiscoveryNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AuthenticateToDiscoveryNodePreCounter)
}

//AuthenticateToDiscoveryNodeFinished returns true if mock invocations count is ok
func (m *ControllerMock) AuthenticateToDiscoveryNodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AuthenticateToDiscoveryNodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AuthenticateToDiscoveryNodeCounter) == uint64(len(m.AuthenticateToDiscoveryNodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AuthenticateToDiscoveryNodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AuthenticateToDiscoveryNodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AuthenticateToDiscoveryNodeFunc != nil {
		return atomic.LoadUint64(&m.AuthenticateToDiscoveryNodeCounter) > 0
	}

	return true
}

type mControllerMockBootstrap struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockBootstrapExpectation
	expectationSeries []*ControllerMockBootstrapExpectation
}

type ControllerMockBootstrapExpectation struct {
	input  *ControllerMockBootstrapInput
	result *ControllerMockBootstrapResult
}

type ControllerMockBootstrapInput struct {
	p context.Context
}

type ControllerMockBootstrapResult struct {
	r  *network.BootstrapResult
	r1 error
}

//Expect specifies that invocation of Controller.Bootstrap is expected from 1 to Infinity times
func (m *mControllerMockBootstrap) Expect(p context.Context) *mControllerMockBootstrap {
	m.mock.BootstrapFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockBootstrapExpectation{}
	}
	m.mainExpectation.input = &ControllerMockBootstrapInput{p}
	return m
}

//Return specifies results of invocation of Controller.Bootstrap
func (m *mControllerMockBootstrap) Return(r *network.BootstrapResult, r1 error) *ControllerMock {
	m.mock.BootstrapFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockBootstrapExpectation{}
	}
	m.mainExpectation.result = &ControllerMockBootstrapResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Controller.Bootstrap is expected once
func (m *mControllerMockBootstrap) ExpectOnce(p context.Context) *ControllerMockBootstrapExpectation {
	m.mock.BootstrapFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockBootstrapExpectation{}
	expectation.input = &ControllerMockBootstrapInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ControllerMockBootstrapExpectation) Return(r *network.BootstrapResult, r1 error) {
	e.result = &ControllerMockBootstrapResult{r, r1}
}

//Set uses given function f as a mock of Controller.Bootstrap method
func (m *mControllerMockBootstrap) Set(f func(p context.Context) (r *network.BootstrapResult, r1 error)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.BootstrapFunc = f
	return m.mock
}

//Bootstrap implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) Bootstrap(p context.Context) (r *network.BootstrapResult, r1 error) {
	counter := atomic.AddUint64(&m.BootstrapPreCounter, 1)
	defer atomic.AddUint64(&m.BootstrapCounter, 1)

	if len(m.BootstrapMock.expectationSeries) > 0 {
		if counter > uint64(len(m.BootstrapMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.Bootstrap. %v", p)
			return
		}

		input := m.BootstrapMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockBootstrapInput{p}, "Controller.Bootstrap got unexpected parameters")

		result := m.BootstrapMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.Bootstrap")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BootstrapMock.mainExpectation != nil {

		input := m.BootstrapMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockBootstrapInput{p}, "Controller.Bootstrap got unexpected parameters")
		}

		result := m.BootstrapMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.Bootstrap")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.BootstrapFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.Bootstrap. %v", p)
		return
	}

	return m.BootstrapFunc(p)
}

//BootstrapMinimockCounter returns a count of ControllerMock.BootstrapFunc invocations
func (m *ControllerMock) BootstrapMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.BootstrapCounter)
}

//BootstrapMinimockPreCounter returns the value of ControllerMock.Bootstrap invocations
func (m *ControllerMock) BootstrapMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.BootstrapPreCounter)
}

//BootstrapFinished returns true if mock invocations count is ok
func (m *ControllerMock) BootstrapFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.BootstrapMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.BootstrapCounter) == uint64(len(m.BootstrapMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.BootstrapMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.BootstrapCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.BootstrapFunc != nil {
		return atomic.LoadUint64(&m.BootstrapCounter) > 0
	}

	return true
}

type mControllerMockGetLastIgnoredPulse struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockGetLastIgnoredPulseExpectation
	expectationSeries []*ControllerMockGetLastIgnoredPulseExpectation
}

type ControllerMockGetLastIgnoredPulseExpectation struct {
	result *ControllerMockGetLastIgnoredPulseResult
}

type ControllerMockGetLastIgnoredPulseResult struct {
	r insolar.PulseNumber
}

//Expect specifies that invocation of Controller.GetLastIgnoredPulse is expected from 1 to Infinity times
func (m *mControllerMockGetLastIgnoredPulse) Expect() *mControllerMockGetLastIgnoredPulse {
	m.mock.GetLastIgnoredPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockGetLastIgnoredPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of Controller.GetLastIgnoredPulse
func (m *mControllerMockGetLastIgnoredPulse) Return(r insolar.PulseNumber) *ControllerMock {
	m.mock.GetLastIgnoredPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockGetLastIgnoredPulseExpectation{}
	}
	m.mainExpectation.result = &ControllerMockGetLastIgnoredPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Controller.GetLastIgnoredPulse is expected once
func (m *mControllerMockGetLastIgnoredPulse) ExpectOnce() *ControllerMockGetLastIgnoredPulseExpectation {
	m.mock.GetLastIgnoredPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockGetLastIgnoredPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ControllerMockGetLastIgnoredPulseExpectation) Return(r insolar.PulseNumber) {
	e.result = &ControllerMockGetLastIgnoredPulseResult{r}
}

//Set uses given function f as a mock of Controller.GetLastIgnoredPulse method
func (m *mControllerMockGetLastIgnoredPulse) Set(f func() (r insolar.PulseNumber)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetLastIgnoredPulseFunc = f
	return m.mock
}

//GetLastIgnoredPulse implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) GetLastIgnoredPulse() (r insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.GetLastIgnoredPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetLastIgnoredPulseCounter, 1)

	if len(m.GetLastIgnoredPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetLastIgnoredPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.GetLastIgnoredPulse.")
			return
		}

		result := m.GetLastIgnoredPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.GetLastIgnoredPulse")
			return
		}

		r = result.r

		return
	}

	if m.GetLastIgnoredPulseMock.mainExpectation != nil {

		result := m.GetLastIgnoredPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.GetLastIgnoredPulse")
		}

		r = result.r

		return
	}

	if m.GetLastIgnoredPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.GetLastIgnoredPulse.")
		return
	}

	return m.GetLastIgnoredPulseFunc()
}

//GetLastIgnoredPulseMinimockCounter returns a count of ControllerMock.GetLastIgnoredPulseFunc invocations
func (m *ControllerMock) GetLastIgnoredPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetLastIgnoredPulseCounter)
}

//GetLastIgnoredPulseMinimockPreCounter returns the value of ControllerMock.GetLastIgnoredPulse invocations
func (m *ControllerMock) GetLastIgnoredPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetLastIgnoredPulsePreCounter)
}

//GetLastIgnoredPulseFinished returns true if mock invocations count is ok
func (m *ControllerMock) GetLastIgnoredPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetLastIgnoredPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetLastIgnoredPulseCounter) == uint64(len(m.GetLastIgnoredPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetLastIgnoredPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetLastIgnoredPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetLastIgnoredPulseFunc != nil {
		return atomic.LoadUint64(&m.GetLastIgnoredPulseCounter) > 0
	}

	return true
}

type mControllerMockInit struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockInitExpectation
	expectationSeries []*ControllerMockInitExpectation
}

type ControllerMockInitExpectation struct {
	input  *ControllerMockInitInput
	result *ControllerMockInitResult
}

type ControllerMockInitInput struct {
	p context.Context
}

type ControllerMockInitResult struct {
	r error
}

//Expect specifies that invocation of Controller.Init is expected from 1 to Infinity times
func (m *mControllerMockInit) Expect(p context.Context) *mControllerMockInit {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockInitExpectation{}
	}
	m.mainExpectation.input = &ControllerMockInitInput{p}
	return m
}

//Return specifies results of invocation of Controller.Init
func (m *mControllerMockInit) Return(r error) *ControllerMock {
	m.mock.InitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockInitExpectation{}
	}
	m.mainExpectation.result = &ControllerMockInitResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Controller.Init is expected once
func (m *mControllerMockInit) ExpectOnce(p context.Context) *ControllerMockInitExpectation {
	m.mock.InitFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockInitExpectation{}
	expectation.input = &ControllerMockInitInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ControllerMockInitExpectation) Return(r error) {
	e.result = &ControllerMockInitResult{r}
}

//Set uses given function f as a mock of Controller.Init method
func (m *mControllerMockInit) Set(f func(p context.Context) (r error)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InitFunc = f
	return m.mock
}

//Init implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) Init(p context.Context) (r error) {
	counter := atomic.AddUint64(&m.InitPreCounter, 1)
	defer atomic.AddUint64(&m.InitCounter, 1)

	if len(m.InitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.Init. %v", p)
			return
		}

		input := m.InitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockInitInput{p}, "Controller.Init got unexpected parameters")

		result := m.InitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.Init")
			return
		}

		r = result.r

		return
	}

	if m.InitMock.mainExpectation != nil {

		input := m.InitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockInitInput{p}, "Controller.Init got unexpected parameters")
		}

		result := m.InitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.Init")
		}

		r = result.r

		return
	}

	if m.InitFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.Init. %v", p)
		return
	}

	return m.InitFunc(p)
}

//InitMinimockCounter returns a count of ControllerMock.InitFunc invocations
func (m *ControllerMock) InitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InitCounter)
}

//InitMinimockPreCounter returns the value of ControllerMock.Init invocations
func (m *ControllerMock) InitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InitPreCounter)
}

//InitFinished returns true if mock invocations count is ok
func (m *ControllerMock) InitFinished() bool {
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

type mControllerMockRemoteProcedureRegister struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockRemoteProcedureRegisterExpectation
	expectationSeries []*ControllerMockRemoteProcedureRegisterExpectation
}

type ControllerMockRemoteProcedureRegisterExpectation struct {
	input *ControllerMockRemoteProcedureRegisterInput
}

type ControllerMockRemoteProcedureRegisterInput struct {
	p  string
	p1 insolar.RemoteProcedure
}

//Expect specifies that invocation of Controller.RemoteProcedureRegister is expected from 1 to Infinity times
func (m *mControllerMockRemoteProcedureRegister) Expect(p string, p1 insolar.RemoteProcedure) *mControllerMockRemoteProcedureRegister {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockRemoteProcedureRegisterExpectation{}
	}
	m.mainExpectation.input = &ControllerMockRemoteProcedureRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of Controller.RemoteProcedureRegister
func (m *mControllerMockRemoteProcedureRegister) Return() *ControllerMock {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockRemoteProcedureRegisterExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Controller.RemoteProcedureRegister is expected once
func (m *mControllerMockRemoteProcedureRegister) ExpectOnce(p string, p1 insolar.RemoteProcedure) *ControllerMockRemoteProcedureRegisterExpectation {
	m.mock.RemoteProcedureRegisterFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockRemoteProcedureRegisterExpectation{}
	expectation.input = &ControllerMockRemoteProcedureRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Controller.RemoteProcedureRegister method
func (m *mControllerMockRemoteProcedureRegister) Set(f func(p string, p1 insolar.RemoteProcedure)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoteProcedureRegisterFunc = f
	return m.mock
}

//RemoteProcedureRegister implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) RemoteProcedureRegister(p string, p1 insolar.RemoteProcedure) {
	counter := atomic.AddUint64(&m.RemoteProcedureRegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RemoteProcedureRegisterCounter, 1)

	if len(m.RemoteProcedureRegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoteProcedureRegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.RemoteProcedureRegister. %v %v", p, p1)
			return
		}

		input := m.RemoteProcedureRegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockRemoteProcedureRegisterInput{p, p1}, "Controller.RemoteProcedureRegister got unexpected parameters")

		return
	}

	if m.RemoteProcedureRegisterMock.mainExpectation != nil {

		input := m.RemoteProcedureRegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockRemoteProcedureRegisterInput{p, p1}, "Controller.RemoteProcedureRegister got unexpected parameters")
		}

		return
	}

	if m.RemoteProcedureRegisterFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.RemoteProcedureRegister. %v %v", p, p1)
		return
	}

	m.RemoteProcedureRegisterFunc(p, p1)
}

//RemoteProcedureRegisterMinimockCounter returns a count of ControllerMock.RemoteProcedureRegisterFunc invocations
func (m *ControllerMock) RemoteProcedureRegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter)
}

//RemoteProcedureRegisterMinimockPreCounter returns the value of ControllerMock.RemoteProcedureRegister invocations
func (m *ControllerMock) RemoteProcedureRegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoteProcedureRegisterPreCounter)
}

//RemoteProcedureRegisterFinished returns true if mock invocations count is ok
func (m *ControllerMock) RemoteProcedureRegisterFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoteProcedureRegisterMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter) == uint64(len(m.RemoteProcedureRegisterMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoteProcedureRegisterMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoteProcedureRegisterFunc != nil {
		return atomic.LoadUint64(&m.RemoteProcedureRegisterCounter) > 0
	}

	return true
}

type mControllerMockSendBytes struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockSendBytesExpectation
	expectationSeries []*ControllerMockSendBytesExpectation
}

type ControllerMockSendBytesExpectation struct {
	input  *ControllerMockSendBytesInput
	result *ControllerMockSendBytesResult
}

type ControllerMockSendBytesInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 string
	p3 []byte
}

type ControllerMockSendBytesResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of Controller.SendBytes is expected from 1 to Infinity times
func (m *mControllerMockSendBytes) Expect(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) *mControllerMockSendBytes {
	m.mock.SendBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSendBytesExpectation{}
	}
	m.mainExpectation.input = &ControllerMockSendBytesInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Controller.SendBytes
func (m *mControllerMockSendBytes) Return(r []byte, r1 error) *ControllerMock {
	m.mock.SendBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSendBytesExpectation{}
	}
	m.mainExpectation.result = &ControllerMockSendBytesResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Controller.SendBytes is expected once
func (m *mControllerMockSendBytes) ExpectOnce(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) *ControllerMockSendBytesExpectation {
	m.mock.SendBytesFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockSendBytesExpectation{}
	expectation.input = &ControllerMockSendBytesInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ControllerMockSendBytesExpectation) Return(r []byte, r1 error) {
	e.result = &ControllerMockSendBytesResult{r, r1}
}

//Set uses given function f as a mock of Controller.SendBytes method
func (m *mControllerMockSendBytes) Set(f func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendBytesFunc = f
	return m.mock
}

//SendBytes implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) SendBytes(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.SendBytesPreCounter, 1)
	defer atomic.AddUint64(&m.SendBytesCounter, 1)

	if len(m.SendBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.SendBytes. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SendBytesMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockSendBytesInput{p, p1, p2, p3}, "Controller.SendBytes got unexpected parameters")

		result := m.SendBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.SendBytes")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendBytesMock.mainExpectation != nil {

		input := m.SendBytesMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockSendBytesInput{p, p1, p2, p3}, "Controller.SendBytes got unexpected parameters")
		}

		result := m.SendBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.SendBytes")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendBytesFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.SendBytes. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SendBytesFunc(p, p1, p2, p3)
}

//SendBytesMinimockCounter returns a count of ControllerMock.SendBytesFunc invocations
func (m *ControllerMock) SendBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendBytesCounter)
}

//SendBytesMinimockPreCounter returns the value of ControllerMock.SendBytes invocations
func (m *ControllerMock) SendBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendBytesPreCounter)
}

//SendBytesFinished returns true if mock invocations count is ok
func (m *ControllerMock) SendBytesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendBytesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendBytesCounter) == uint64(len(m.SendBytesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendBytesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendBytesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendBytesFunc != nil {
		return atomic.LoadUint64(&m.SendBytesCounter) > 0
	}

	return true
}

type mControllerMockSendCascadeMessage struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockSendCascadeMessageExpectation
	expectationSeries []*ControllerMockSendCascadeMessageExpectation
}

type ControllerMockSendCascadeMessageExpectation struct {
	input  *ControllerMockSendCascadeMessageInput
	result *ControllerMockSendCascadeMessageResult
}

type ControllerMockSendCascadeMessageInput struct {
	p  insolar.Cascade
	p1 string
	p2 insolar.Parcel
}

type ControllerMockSendCascadeMessageResult struct {
	r error
}

//Expect specifies that invocation of Controller.SendCascadeMessage is expected from 1 to Infinity times
func (m *mControllerMockSendCascadeMessage) Expect(p insolar.Cascade, p1 string, p2 insolar.Parcel) *mControllerMockSendCascadeMessage {
	m.mock.SendCascadeMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSendCascadeMessageExpectation{}
	}
	m.mainExpectation.input = &ControllerMockSendCascadeMessageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Controller.SendCascadeMessage
func (m *mControllerMockSendCascadeMessage) Return(r error) *ControllerMock {
	m.mock.SendCascadeMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSendCascadeMessageExpectation{}
	}
	m.mainExpectation.result = &ControllerMockSendCascadeMessageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Controller.SendCascadeMessage is expected once
func (m *mControllerMockSendCascadeMessage) ExpectOnce(p insolar.Cascade, p1 string, p2 insolar.Parcel) *ControllerMockSendCascadeMessageExpectation {
	m.mock.SendCascadeMessageFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockSendCascadeMessageExpectation{}
	expectation.input = &ControllerMockSendCascadeMessageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ControllerMockSendCascadeMessageExpectation) Return(r error) {
	e.result = &ControllerMockSendCascadeMessageResult{r}
}

//Set uses given function f as a mock of Controller.SendCascadeMessage method
func (m *mControllerMockSendCascadeMessage) Set(f func(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendCascadeMessageFunc = f
	return m.mock
}

//SendCascadeMessage implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) SendCascadeMessage(p insolar.Cascade, p1 string, p2 insolar.Parcel) (r error) {
	counter := atomic.AddUint64(&m.SendCascadeMessagePreCounter, 1)
	defer atomic.AddUint64(&m.SendCascadeMessageCounter, 1)

	if len(m.SendCascadeMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendCascadeMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.SendCascadeMessage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendCascadeMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockSendCascadeMessageInput{p, p1, p2}, "Controller.SendCascadeMessage got unexpected parameters")

		result := m.SendCascadeMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.SendCascadeMessage")
			return
		}

		r = result.r

		return
	}

	if m.SendCascadeMessageMock.mainExpectation != nil {

		input := m.SendCascadeMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockSendCascadeMessageInput{p, p1, p2}, "Controller.SendCascadeMessage got unexpected parameters")
		}

		result := m.SendCascadeMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.SendCascadeMessage")
		}

		r = result.r

		return
	}

	if m.SendCascadeMessageFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.SendCascadeMessage. %v %v %v", p, p1, p2)
		return
	}

	return m.SendCascadeMessageFunc(p, p1, p2)
}

//SendCascadeMessageMinimockCounter returns a count of ControllerMock.SendCascadeMessageFunc invocations
func (m *ControllerMock) SendCascadeMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCascadeMessageCounter)
}

//SendCascadeMessageMinimockPreCounter returns the value of ControllerMock.SendCascadeMessage invocations
func (m *ControllerMock) SendCascadeMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendCascadeMessagePreCounter)
}

//SendCascadeMessageFinished returns true if mock invocations count is ok
func (m *ControllerMock) SendCascadeMessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendCascadeMessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendCascadeMessageCounter) == uint64(len(m.SendCascadeMessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendCascadeMessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendCascadeMessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendCascadeMessageFunc != nil {
		return atomic.LoadUint64(&m.SendCascadeMessageCounter) > 0
	}

	return true
}

type mControllerMockSendMessage struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockSendMessageExpectation
	expectationSeries []*ControllerMockSendMessageExpectation
}

type ControllerMockSendMessageExpectation struct {
	input  *ControllerMockSendMessageInput
	result *ControllerMockSendMessageResult
}

type ControllerMockSendMessageInput struct {
	p  insolar.Reference
	p1 string
	p2 insolar.Parcel
}

type ControllerMockSendMessageResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of Controller.SendMessage is expected from 1 to Infinity times
func (m *mControllerMockSendMessage) Expect(p insolar.Reference, p1 string, p2 insolar.Parcel) *mControllerMockSendMessage {
	m.mock.SendMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSendMessageExpectation{}
	}
	m.mainExpectation.input = &ControllerMockSendMessageInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Controller.SendMessage
func (m *mControllerMockSendMessage) Return(r []byte, r1 error) *ControllerMock {
	m.mock.SendMessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSendMessageExpectation{}
	}
	m.mainExpectation.result = &ControllerMockSendMessageResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Controller.SendMessage is expected once
func (m *mControllerMockSendMessage) ExpectOnce(p insolar.Reference, p1 string, p2 insolar.Parcel) *ControllerMockSendMessageExpectation {
	m.mock.SendMessageFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockSendMessageExpectation{}
	expectation.input = &ControllerMockSendMessageInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ControllerMockSendMessageExpectation) Return(r []byte, r1 error) {
	e.result = &ControllerMockSendMessageResult{r, r1}
}

//Set uses given function f as a mock of Controller.SendMessage method
func (m *mControllerMockSendMessage) Set(f func(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendMessageFunc = f
	return m.mock
}

//SendMessage implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) SendMessage(p insolar.Reference, p1 string, p2 insolar.Parcel) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.SendMessagePreCounter, 1)
	defer atomic.AddUint64(&m.SendMessageCounter, 1)

	if len(m.SendMessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.SendMessage. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendMessageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockSendMessageInput{p, p1, p2}, "Controller.SendMessage got unexpected parameters")

		result := m.SendMessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.SendMessage")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMessageMock.mainExpectation != nil {

		input := m.SendMessageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockSendMessageInput{p, p1, p2}, "Controller.SendMessage got unexpected parameters")
		}

		result := m.SendMessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ControllerMock.SendMessage")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMessageFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.SendMessage. %v %v %v", p, p1, p2)
		return
	}

	return m.SendMessageFunc(p, p1, p2)
}

//SendMessageMinimockCounter returns a count of ControllerMock.SendMessageFunc invocations
func (m *ControllerMock) SendMessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessageCounter)
}

//SendMessageMinimockPreCounter returns the value of ControllerMock.SendMessage invocations
func (m *ControllerMock) SendMessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendMessagePreCounter)
}

//SendMessageFinished returns true if mock invocations count is ok
func (m *ControllerMock) SendMessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendMessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendMessageCounter) == uint64(len(m.SendMessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendMessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendMessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendMessageFunc != nil {
		return atomic.LoadUint64(&m.SendMessageCounter) > 0
	}

	return true
}

type mControllerMockSetLastIgnoredPulse struct {
	mock              *ControllerMock
	mainExpectation   *ControllerMockSetLastIgnoredPulseExpectation
	expectationSeries []*ControllerMockSetLastIgnoredPulseExpectation
}

type ControllerMockSetLastIgnoredPulseExpectation struct {
	input *ControllerMockSetLastIgnoredPulseInput
}

type ControllerMockSetLastIgnoredPulseInput struct {
	p insolar.PulseNumber
}

//Expect specifies that invocation of Controller.SetLastIgnoredPulse is expected from 1 to Infinity times
func (m *mControllerMockSetLastIgnoredPulse) Expect(p insolar.PulseNumber) *mControllerMockSetLastIgnoredPulse {
	m.mock.SetLastIgnoredPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSetLastIgnoredPulseExpectation{}
	}
	m.mainExpectation.input = &ControllerMockSetLastIgnoredPulseInput{p}
	return m
}

//Return specifies results of invocation of Controller.SetLastIgnoredPulse
func (m *mControllerMockSetLastIgnoredPulse) Return() *ControllerMock {
	m.mock.SetLastIgnoredPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ControllerMockSetLastIgnoredPulseExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Controller.SetLastIgnoredPulse is expected once
func (m *mControllerMockSetLastIgnoredPulse) ExpectOnce(p insolar.PulseNumber) *ControllerMockSetLastIgnoredPulseExpectation {
	m.mock.SetLastIgnoredPulseFunc = nil
	m.mainExpectation = nil

	expectation := &ControllerMockSetLastIgnoredPulseExpectation{}
	expectation.input = &ControllerMockSetLastIgnoredPulseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Controller.SetLastIgnoredPulse method
func (m *mControllerMockSetLastIgnoredPulse) Set(f func(p insolar.PulseNumber)) *ControllerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLastIgnoredPulseFunc = f
	return m.mock
}

//SetLastIgnoredPulse implements github.com/insolar/insolar/network.Controller interface
func (m *ControllerMock) SetLastIgnoredPulse(p insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.SetLastIgnoredPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetLastIgnoredPulseCounter, 1)

	if len(m.SetLastIgnoredPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLastIgnoredPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ControllerMock.SetLastIgnoredPulse. %v", p)
			return
		}

		input := m.SetLastIgnoredPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ControllerMockSetLastIgnoredPulseInput{p}, "Controller.SetLastIgnoredPulse got unexpected parameters")

		return
	}

	if m.SetLastIgnoredPulseMock.mainExpectation != nil {

		input := m.SetLastIgnoredPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ControllerMockSetLastIgnoredPulseInput{p}, "Controller.SetLastIgnoredPulse got unexpected parameters")
		}

		return
	}

	if m.SetLastIgnoredPulseFunc == nil {
		m.t.Fatalf("Unexpected call to ControllerMock.SetLastIgnoredPulse. %v", p)
		return
	}

	m.SetLastIgnoredPulseFunc(p)
}

//SetLastIgnoredPulseMinimockCounter returns a count of ControllerMock.SetLastIgnoredPulseFunc invocations
func (m *ControllerMock) SetLastIgnoredPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLastIgnoredPulseCounter)
}

//SetLastIgnoredPulseMinimockPreCounter returns the value of ControllerMock.SetLastIgnoredPulse invocations
func (m *ControllerMock) SetLastIgnoredPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLastIgnoredPulsePreCounter)
}

//SetLastIgnoredPulseFinished returns true if mock invocations count is ok
func (m *ControllerMock) SetLastIgnoredPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetLastIgnoredPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetLastIgnoredPulseCounter) == uint64(len(m.SetLastIgnoredPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetLastIgnoredPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetLastIgnoredPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetLastIgnoredPulseFunc != nil {
		return atomic.LoadUint64(&m.SetLastIgnoredPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ControllerMock) ValidateCallCounters() {

	if !m.AuthenticateToDiscoveryNodeFinished() {
		m.t.Fatal("Expected call to ControllerMock.AuthenticateToDiscoveryNode")
	}

	if !m.BootstrapFinished() {
		m.t.Fatal("Expected call to ControllerMock.Bootstrap")
	}

	if !m.GetLastIgnoredPulseFinished() {
		m.t.Fatal("Expected call to ControllerMock.GetLastIgnoredPulse")
	}

	if !m.InitFinished() {
		m.t.Fatal("Expected call to ControllerMock.Init")
	}

	if !m.RemoteProcedureRegisterFinished() {
		m.t.Fatal("Expected call to ControllerMock.RemoteProcedureRegister")
	}

	if !m.SendBytesFinished() {
		m.t.Fatal("Expected call to ControllerMock.SendBytes")
	}

	if !m.SendCascadeMessageFinished() {
		m.t.Fatal("Expected call to ControllerMock.SendCascadeMessage")
	}

	if !m.SendMessageFinished() {
		m.t.Fatal("Expected call to ControllerMock.SendMessage")
	}

	if !m.SetLastIgnoredPulseFinished() {
		m.t.Fatal("Expected call to ControllerMock.SetLastIgnoredPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ControllerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ControllerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ControllerMock) MinimockFinish() {

	if !m.AuthenticateToDiscoveryNodeFinished() {
		m.t.Fatal("Expected call to ControllerMock.AuthenticateToDiscoveryNode")
	}

	if !m.BootstrapFinished() {
		m.t.Fatal("Expected call to ControllerMock.Bootstrap")
	}

	if !m.GetLastIgnoredPulseFinished() {
		m.t.Fatal("Expected call to ControllerMock.GetLastIgnoredPulse")
	}

	if !m.InitFinished() {
		m.t.Fatal("Expected call to ControllerMock.Init")
	}

	if !m.RemoteProcedureRegisterFinished() {
		m.t.Fatal("Expected call to ControllerMock.RemoteProcedureRegister")
	}

	if !m.SendBytesFinished() {
		m.t.Fatal("Expected call to ControllerMock.SendBytes")
	}

	if !m.SendCascadeMessageFinished() {
		m.t.Fatal("Expected call to ControllerMock.SendCascadeMessage")
	}

	if !m.SendMessageFinished() {
		m.t.Fatal("Expected call to ControllerMock.SendMessage")
	}

	if !m.SetLastIgnoredPulseFinished() {
		m.t.Fatal("Expected call to ControllerMock.SetLastIgnoredPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ControllerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ControllerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AuthenticateToDiscoveryNodeFinished()
		ok = ok && m.BootstrapFinished()
		ok = ok && m.GetLastIgnoredPulseFinished()
		ok = ok && m.InitFinished()
		ok = ok && m.RemoteProcedureRegisterFinished()
		ok = ok && m.SendBytesFinished()
		ok = ok && m.SendCascadeMessageFinished()
		ok = ok && m.SendMessageFinished()
		ok = ok && m.SetLastIgnoredPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AuthenticateToDiscoveryNodeFinished() {
				m.t.Error("Expected call to ControllerMock.AuthenticateToDiscoveryNode")
			}

			if !m.BootstrapFinished() {
				m.t.Error("Expected call to ControllerMock.Bootstrap")
			}

			if !m.GetLastIgnoredPulseFinished() {
				m.t.Error("Expected call to ControllerMock.GetLastIgnoredPulse")
			}

			if !m.InitFinished() {
				m.t.Error("Expected call to ControllerMock.Init")
			}

			if !m.RemoteProcedureRegisterFinished() {
				m.t.Error("Expected call to ControllerMock.RemoteProcedureRegister")
			}

			if !m.SendBytesFinished() {
				m.t.Error("Expected call to ControllerMock.SendBytes")
			}

			if !m.SendCascadeMessageFinished() {
				m.t.Error("Expected call to ControllerMock.SendCascadeMessage")
			}

			if !m.SendMessageFinished() {
				m.t.Error("Expected call to ControllerMock.SendMessage")
			}

			if !m.SetLastIgnoredPulseFinished() {
				m.t.Error("Expected call to ControllerMock.SetLastIgnoredPulse")
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
func (m *ControllerMock) AllMocksCalled() bool {

	if !m.AuthenticateToDiscoveryNodeFinished() {
		return false
	}

	if !m.BootstrapFinished() {
		return false
	}

	if !m.GetLastIgnoredPulseFinished() {
		return false
	}

	if !m.InitFinished() {
		return false
	}

	if !m.RemoteProcedureRegisterFinished() {
		return false
	}

	if !m.SendBytesFinished() {
		return false
	}

	if !m.SendCascadeMessageFinished() {
		return false
	}

	if !m.SendMessageFinished() {
		return false
	}

	if !m.SetLastIgnoredPulseFinished() {
		return false
	}

	return true
}
