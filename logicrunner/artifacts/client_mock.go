package artifacts

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Client" can be found in github.com/insolar/insolar/logicrunner/artifacts
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	record "github.com/insolar/insolar/insolar/record"

	testify_assert "github.com/stretchr/testify/assert"
)

//ClientMock implements github.com/insolar/insolar/logicrunner/artifacts.Client
type ClientMock struct {
	t minimock.Tester

	ActivatePrototypeFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 []byte) (r error)
	ActivatePrototypeCounter    uint64
	ActivatePrototypePreCounter uint64
	ActivatePrototypeMock       mClientMockActivatePrototype

	DeployCodeFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte, p4 insolar.MachineType) (r *insolar.ID, r1 error)
	DeployCodeCounter    uint64
	DeployCodePreCounter uint64
	DeployCodeMock       mClientMockDeployCode

	GetChildrenFunc       func(p context.Context, p1 insolar.Reference, p2 *insolar.PulseNumber) (r RefIterator, r1 error)
	GetChildrenCounter    uint64
	GetChildrenPreCounter uint64
	GetChildrenMock       mClientMockGetChildren

	GetCodeFunc       func(p context.Context, p1 insolar.Reference) (r CodeDescriptor, r1 error)
	GetCodeCounter    uint64
	GetCodePreCounter uint64
	GetCodeMock       mClientMockGetCode

	GetDelegateFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference) (r *insolar.Reference, r1 error)
	GetDelegateCounter    uint64
	GetDelegatePreCounter uint64
	GetDelegateMock       mClientMockGetDelegate

	GetObjectFunc       func(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 error)
	GetObjectCounter    uint64
	GetObjectPreCounter uint64
	GetObjectMock       mClientMockGetObject

	GetPendingRequestFunc       func(p context.Context, p1 insolar.ID) (r *insolar.Reference, r1 *record.IncomingRequest, r2 error)
	GetPendingRequestCounter    uint64
	GetPendingRequestPreCounter uint64
	GetPendingRequestMock       mClientMockGetPendingRequest

	HasPendingRequestsFunc       func(p context.Context, p1 insolar.Reference) (r bool, r1 error)
	HasPendingRequestsCounter    uint64
	HasPendingRequestsPreCounter uint64
	HasPendingRequestsMock       mClientMockHasPendingRequests

	InjectCodeDescriptorFunc       func(p insolar.Reference, p1 CodeDescriptor)
	InjectCodeDescriptorCounter    uint64
	InjectCodeDescriptorPreCounter uint64
	InjectCodeDescriptorMock       mClientMockInjectCodeDescriptor

	InjectFinishFunc       func()
	InjectFinishCounter    uint64
	InjectFinishPreCounter uint64
	InjectFinishMock       mClientMockInjectFinish

	InjectObjectDescriptorFunc       func(p insolar.Reference, p1 ObjectDescriptor)
	InjectObjectDescriptorCounter    uint64
	InjectObjectDescriptorPreCounter uint64
	InjectObjectDescriptorMock       mClientMockInjectObjectDescriptor

	RegisterIncomingRequestFunc       func(p context.Context, p1 *record.IncomingRequest) (r *insolar.ID, r1 error)
	RegisterIncomingRequestCounter    uint64
	RegisterIncomingRequestPreCounter uint64
	RegisterIncomingRequestMock       mClientMockRegisterIncomingRequest

	RegisterOutgoingRequestFunc       func(p context.Context, p1 *record.OutgoingRequest) (r *insolar.ID, r1 error)
	RegisterOutgoingRequestCounter    uint64
	RegisterOutgoingRequestPreCounter uint64
	RegisterOutgoingRequestMock       mClientMockRegisterOutgoingRequest

	RegisterResultFunc       func(p context.Context, p1 insolar.Reference, p2 RequestResult) (r error)
	RegisterResultCounter    uint64
	RegisterResultPreCounter uint64
	RegisterResultMock       mClientMockRegisterResult

	RegisterValidationFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.ID, p3 bool, p4 []insolar.Message) (r error)
	RegisterValidationCounter    uint64
	RegisterValidationPreCounter uint64
	RegisterValidationMock       mClientMockRegisterValidation

	StateFunc       func() (r []byte)
	StateCounter    uint64
	StatePreCounter uint64
	StateMock       mClientMockState
}

//NewClientMock returns a mock for github.com/insolar/insolar/logicrunner/artifacts.Client
func NewClientMock(t minimock.Tester) *ClientMock {
	m := &ClientMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ActivatePrototypeMock = mClientMockActivatePrototype{mock: m}
	m.DeployCodeMock = mClientMockDeployCode{mock: m}
	m.GetChildrenMock = mClientMockGetChildren{mock: m}
	m.GetCodeMock = mClientMockGetCode{mock: m}
	m.GetDelegateMock = mClientMockGetDelegate{mock: m}
	m.GetObjectMock = mClientMockGetObject{mock: m}
	m.GetPendingRequestMock = mClientMockGetPendingRequest{mock: m}
	m.HasPendingRequestsMock = mClientMockHasPendingRequests{mock: m}
	m.InjectCodeDescriptorMock = mClientMockInjectCodeDescriptor{mock: m}
	m.InjectFinishMock = mClientMockInjectFinish{mock: m}
	m.InjectObjectDescriptorMock = mClientMockInjectObjectDescriptor{mock: m}
	m.RegisterIncomingRequestMock = mClientMockRegisterIncomingRequest{mock: m}
	m.RegisterOutgoingRequestMock = mClientMockRegisterOutgoingRequest{mock: m}
	m.RegisterResultMock = mClientMockRegisterResult{mock: m}
	m.RegisterValidationMock = mClientMockRegisterValidation{mock: m}
	m.StateMock = mClientMockState{mock: m}

	return m
}

type mClientMockActivatePrototype struct {
	mock              *ClientMock
	mainExpectation   *ClientMockActivatePrototypeExpectation
	expectationSeries []*ClientMockActivatePrototypeExpectation
}

type ClientMockActivatePrototypeExpectation struct {
	input  *ClientMockActivatePrototypeInput
	result *ClientMockActivatePrototypeResult
}

type ClientMockActivatePrototypeInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 insolar.Reference
	p4 []byte
}

type ClientMockActivatePrototypeResult struct {
	r error
}

//Expect specifies that invocation of Client.ActivatePrototype is expected from 1 to Infinity times
func (m *mClientMockActivatePrototype) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 []byte) *mClientMockActivatePrototype {
	m.mock.ActivatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockActivatePrototypeExpectation{}
	}
	m.mainExpectation.input = &ClientMockActivatePrototypeInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of Client.ActivatePrototype
func (m *mClientMockActivatePrototype) Return(r error) *ClientMock {
	m.mock.ActivatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockActivatePrototypeExpectation{}
	}
	m.mainExpectation.result = &ClientMockActivatePrototypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.ActivatePrototype is expected once
func (m *mClientMockActivatePrototype) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 []byte) *ClientMockActivatePrototypeExpectation {
	m.mock.ActivatePrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockActivatePrototypeExpectation{}
	expectation.input = &ClientMockActivatePrototypeInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockActivatePrototypeExpectation) Return(r error) {
	e.result = &ClientMockActivatePrototypeResult{r}
}

//Set uses given function f as a mock of Client.ActivatePrototype method
func (m *mClientMockActivatePrototype) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 []byte) (r error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ActivatePrototypeFunc = f
	return m.mock
}

//ActivatePrototype implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) ActivatePrototype(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 []byte) (r error) {
	counter := atomic.AddUint64(&m.ActivatePrototypePreCounter, 1)
	defer atomic.AddUint64(&m.ActivatePrototypeCounter, 1)

	if len(m.ActivatePrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ActivatePrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.ActivatePrototype. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.ActivatePrototypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockActivatePrototypeInput{p, p1, p2, p3, p4}, "Client.ActivatePrototype got unexpected parameters")

		result := m.ActivatePrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.ActivatePrototype")
			return
		}

		r = result.r

		return
	}

	if m.ActivatePrototypeMock.mainExpectation != nil {

		input := m.ActivatePrototypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockActivatePrototypeInput{p, p1, p2, p3, p4}, "Client.ActivatePrototype got unexpected parameters")
		}

		result := m.ActivatePrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.ActivatePrototype")
		}

		r = result.r

		return
	}

	if m.ActivatePrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.ActivatePrototype. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.ActivatePrototypeFunc(p, p1, p2, p3, p4)
}

//ActivatePrototypeMinimockCounter returns a count of ClientMock.ActivatePrototypeFunc invocations
func (m *ClientMock) ActivatePrototypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ActivatePrototypeCounter)
}

//ActivatePrototypeMinimockPreCounter returns the value of ClientMock.ActivatePrototype invocations
func (m *ClientMock) ActivatePrototypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ActivatePrototypePreCounter)
}

//ActivatePrototypeFinished returns true if mock invocations count is ok
func (m *ClientMock) ActivatePrototypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ActivatePrototypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ActivatePrototypeCounter) == uint64(len(m.ActivatePrototypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ActivatePrototypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ActivatePrototypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ActivatePrototypeFunc != nil {
		return atomic.LoadUint64(&m.ActivatePrototypeCounter) > 0
	}

	return true
}

type mClientMockDeployCode struct {
	mock              *ClientMock
	mainExpectation   *ClientMockDeployCodeExpectation
	expectationSeries []*ClientMockDeployCodeExpectation
}

type ClientMockDeployCodeExpectation struct {
	input  *ClientMockDeployCodeInput
	result *ClientMockDeployCodeResult
}

type ClientMockDeployCodeInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 []byte
	p4 insolar.MachineType
}

type ClientMockDeployCodeResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Client.DeployCode is expected from 1 to Infinity times
func (m *mClientMockDeployCode) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte, p4 insolar.MachineType) *mClientMockDeployCode {
	m.mock.DeployCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockDeployCodeExpectation{}
	}
	m.mainExpectation.input = &ClientMockDeployCodeInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of Client.DeployCode
func (m *mClientMockDeployCode) Return(r *insolar.ID, r1 error) *ClientMock {
	m.mock.DeployCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockDeployCodeExpectation{}
	}
	m.mainExpectation.result = &ClientMockDeployCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.DeployCode is expected once
func (m *mClientMockDeployCode) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte, p4 insolar.MachineType) *ClientMockDeployCodeExpectation {
	m.mock.DeployCodeFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockDeployCodeExpectation{}
	expectation.input = &ClientMockDeployCodeInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockDeployCodeExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ClientMockDeployCodeResult{r, r1}
}

//Set uses given function f as a mock of Client.DeployCode method
func (m *mClientMockDeployCode) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte, p4 insolar.MachineType) (r *insolar.ID, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeployCodeFunc = f
	return m.mock
}

//DeployCode implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) DeployCode(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte, p4 insolar.MachineType) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.DeployCodePreCounter, 1)
	defer atomic.AddUint64(&m.DeployCodeCounter, 1)

	if len(m.DeployCodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeployCodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.DeployCode. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.DeployCodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockDeployCodeInput{p, p1, p2, p3, p4}, "Client.DeployCode got unexpected parameters")

		result := m.DeployCodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.DeployCode")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeployCodeMock.mainExpectation != nil {

		input := m.DeployCodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockDeployCodeInput{p, p1, p2, p3, p4}, "Client.DeployCode got unexpected parameters")
		}

		result := m.DeployCodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.DeployCode")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeployCodeFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.DeployCode. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.DeployCodeFunc(p, p1, p2, p3, p4)
}

//DeployCodeMinimockCounter returns a count of ClientMock.DeployCodeFunc invocations
func (m *ClientMock) DeployCodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeployCodeCounter)
}

//DeployCodeMinimockPreCounter returns the value of ClientMock.DeployCode invocations
func (m *ClientMock) DeployCodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeployCodePreCounter)
}

//DeployCodeFinished returns true if mock invocations count is ok
func (m *ClientMock) DeployCodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeployCodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeployCodeCounter) == uint64(len(m.DeployCodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeployCodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeployCodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeployCodeFunc != nil {
		return atomic.LoadUint64(&m.DeployCodeCounter) > 0
	}

	return true
}

type mClientMockGetChildren struct {
	mock              *ClientMock
	mainExpectation   *ClientMockGetChildrenExpectation
	expectationSeries []*ClientMockGetChildrenExpectation
}

type ClientMockGetChildrenExpectation struct {
	input  *ClientMockGetChildrenInput
	result *ClientMockGetChildrenResult
}

type ClientMockGetChildrenInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 *insolar.PulseNumber
}

type ClientMockGetChildrenResult struct {
	r  RefIterator
	r1 error
}

//Expect specifies that invocation of Client.GetChildren is expected from 1 to Infinity times
func (m *mClientMockGetChildren) Expect(p context.Context, p1 insolar.Reference, p2 *insolar.PulseNumber) *mClientMockGetChildren {
	m.mock.GetChildrenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetChildrenExpectation{}
	}
	m.mainExpectation.input = &ClientMockGetChildrenInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Client.GetChildren
func (m *mClientMockGetChildren) Return(r RefIterator, r1 error) *ClientMock {
	m.mock.GetChildrenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetChildrenExpectation{}
	}
	m.mainExpectation.result = &ClientMockGetChildrenResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.GetChildren is expected once
func (m *mClientMockGetChildren) ExpectOnce(p context.Context, p1 insolar.Reference, p2 *insolar.PulseNumber) *ClientMockGetChildrenExpectation {
	m.mock.GetChildrenFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockGetChildrenExpectation{}
	expectation.input = &ClientMockGetChildrenInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockGetChildrenExpectation) Return(r RefIterator, r1 error) {
	e.result = &ClientMockGetChildrenResult{r, r1}
}

//Set uses given function f as a mock of Client.GetChildren method
func (m *mClientMockGetChildren) Set(f func(p context.Context, p1 insolar.Reference, p2 *insolar.PulseNumber) (r RefIterator, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetChildrenFunc = f
	return m.mock
}

//GetChildren implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) GetChildren(p context.Context, p1 insolar.Reference, p2 *insolar.PulseNumber) (r RefIterator, r1 error) {
	counter := atomic.AddUint64(&m.GetChildrenPreCounter, 1)
	defer atomic.AddUint64(&m.GetChildrenCounter, 1)

	if len(m.GetChildrenMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetChildrenMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.GetChildren. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetChildrenMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockGetChildrenInput{p, p1, p2}, "Client.GetChildren got unexpected parameters")

		result := m.GetChildrenMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetChildren")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetChildrenMock.mainExpectation != nil {

		input := m.GetChildrenMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockGetChildrenInput{p, p1, p2}, "Client.GetChildren got unexpected parameters")
		}

		result := m.GetChildrenMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetChildren")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetChildrenFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.GetChildren. %v %v %v", p, p1, p2)
		return
	}

	return m.GetChildrenFunc(p, p1, p2)
}

//GetChildrenMinimockCounter returns a count of ClientMock.GetChildrenFunc invocations
func (m *ClientMock) GetChildrenMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetChildrenCounter)
}

//GetChildrenMinimockPreCounter returns the value of ClientMock.GetChildren invocations
func (m *ClientMock) GetChildrenMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetChildrenPreCounter)
}

//GetChildrenFinished returns true if mock invocations count is ok
func (m *ClientMock) GetChildrenFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetChildrenMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetChildrenCounter) == uint64(len(m.GetChildrenMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetChildrenMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetChildrenCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetChildrenFunc != nil {
		return atomic.LoadUint64(&m.GetChildrenCounter) > 0
	}

	return true
}

type mClientMockGetCode struct {
	mock              *ClientMock
	mainExpectation   *ClientMockGetCodeExpectation
	expectationSeries []*ClientMockGetCodeExpectation
}

type ClientMockGetCodeExpectation struct {
	input  *ClientMockGetCodeInput
	result *ClientMockGetCodeResult
}

type ClientMockGetCodeInput struct {
	p  context.Context
	p1 insolar.Reference
}

type ClientMockGetCodeResult struct {
	r  CodeDescriptor
	r1 error
}

//Expect specifies that invocation of Client.GetCode is expected from 1 to Infinity times
func (m *mClientMockGetCode) Expect(p context.Context, p1 insolar.Reference) *mClientMockGetCode {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetCodeExpectation{}
	}
	m.mainExpectation.input = &ClientMockGetCodeInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.GetCode
func (m *mClientMockGetCode) Return(r CodeDescriptor, r1 error) *ClientMock {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetCodeExpectation{}
	}
	m.mainExpectation.result = &ClientMockGetCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.GetCode is expected once
func (m *mClientMockGetCode) ExpectOnce(p context.Context, p1 insolar.Reference) *ClientMockGetCodeExpectation {
	m.mock.GetCodeFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockGetCodeExpectation{}
	expectation.input = &ClientMockGetCodeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockGetCodeExpectation) Return(r CodeDescriptor, r1 error) {
	e.result = &ClientMockGetCodeResult{r, r1}
}

//Set uses given function f as a mock of Client.GetCode method
func (m *mClientMockGetCode) Set(f func(p context.Context, p1 insolar.Reference) (r CodeDescriptor, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCodeFunc = f
	return m.mock
}

//GetCode implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) GetCode(p context.Context, p1 insolar.Reference) (r CodeDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.GetCodePreCounter, 1)
	defer atomic.AddUint64(&m.GetCodeCounter, 1)

	if len(m.GetCodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.GetCode. %v %v", p, p1)
			return
		}

		input := m.GetCodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockGetCodeInput{p, p1}, "Client.GetCode got unexpected parameters")

		result := m.GetCodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetCode")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeMock.mainExpectation != nil {

		input := m.GetCodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockGetCodeInput{p, p1}, "Client.GetCode got unexpected parameters")
		}

		result := m.GetCodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetCode")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.GetCode. %v %v", p, p1)
		return
	}

	return m.GetCodeFunc(p, p1)
}

//GetCodeMinimockCounter returns a count of ClientMock.GetCodeFunc invocations
func (m *ClientMock) GetCodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodeCounter)
}

//GetCodeMinimockPreCounter returns the value of ClientMock.GetCode invocations
func (m *ClientMock) GetCodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodePreCounter)
}

//GetCodeFinished returns true if mock invocations count is ok
func (m *ClientMock) GetCodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCodeCounter) == uint64(len(m.GetCodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCodeFunc != nil {
		return atomic.LoadUint64(&m.GetCodeCounter) > 0
	}

	return true
}

type mClientMockGetDelegate struct {
	mock              *ClientMock
	mainExpectation   *ClientMockGetDelegateExpectation
	expectationSeries []*ClientMockGetDelegateExpectation
}

type ClientMockGetDelegateExpectation struct {
	input  *ClientMockGetDelegateInput
	result *ClientMockGetDelegateResult
}

type ClientMockGetDelegateInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
}

type ClientMockGetDelegateResult struct {
	r  *insolar.Reference
	r1 error
}

//Expect specifies that invocation of Client.GetDelegate is expected from 1 to Infinity times
func (m *mClientMockGetDelegate) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference) *mClientMockGetDelegate {
	m.mock.GetDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetDelegateExpectation{}
	}
	m.mainExpectation.input = &ClientMockGetDelegateInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Client.GetDelegate
func (m *mClientMockGetDelegate) Return(r *insolar.Reference, r1 error) *ClientMock {
	m.mock.GetDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetDelegateExpectation{}
	}
	m.mainExpectation.result = &ClientMockGetDelegateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.GetDelegate is expected once
func (m *mClientMockGetDelegate) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference) *ClientMockGetDelegateExpectation {
	m.mock.GetDelegateFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockGetDelegateExpectation{}
	expectation.input = &ClientMockGetDelegateInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockGetDelegateExpectation) Return(r *insolar.Reference, r1 error) {
	e.result = &ClientMockGetDelegateResult{r, r1}
}

//Set uses given function f as a mock of Client.GetDelegate method
func (m *mClientMockGetDelegate) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference) (r *insolar.Reference, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDelegateFunc = f
	return m.mock
}

//GetDelegate implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) GetDelegate(p context.Context, p1 insolar.Reference, p2 insolar.Reference) (r *insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.GetDelegatePreCounter, 1)
	defer atomic.AddUint64(&m.GetDelegateCounter, 1)

	if len(m.GetDelegateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDelegateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.GetDelegate. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetDelegateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockGetDelegateInput{p, p1, p2}, "Client.GetDelegate got unexpected parameters")

		result := m.GetDelegateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetDelegate")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDelegateMock.mainExpectation != nil {

		input := m.GetDelegateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockGetDelegateInput{p, p1, p2}, "Client.GetDelegate got unexpected parameters")
		}

		result := m.GetDelegateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetDelegate")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDelegateFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.GetDelegate. %v %v %v", p, p1, p2)
		return
	}

	return m.GetDelegateFunc(p, p1, p2)
}

//GetDelegateMinimockCounter returns a count of ClientMock.GetDelegateFunc invocations
func (m *ClientMock) GetDelegateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDelegateCounter)
}

//GetDelegateMinimockPreCounter returns the value of ClientMock.GetDelegate invocations
func (m *ClientMock) GetDelegateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDelegatePreCounter)
}

//GetDelegateFinished returns true if mock invocations count is ok
func (m *ClientMock) GetDelegateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDelegateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDelegateCounter) == uint64(len(m.GetDelegateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDelegateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDelegateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDelegateFunc != nil {
		return atomic.LoadUint64(&m.GetDelegateCounter) > 0
	}

	return true
}

type mClientMockGetObject struct {
	mock              *ClientMock
	mainExpectation   *ClientMockGetObjectExpectation
	expectationSeries []*ClientMockGetObjectExpectation
}

type ClientMockGetObjectExpectation struct {
	input  *ClientMockGetObjectInput
	result *ClientMockGetObjectResult
}

type ClientMockGetObjectInput struct {
	p  context.Context
	p1 insolar.Reference
}

type ClientMockGetObjectResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Client.GetObject is expected from 1 to Infinity times
func (m *mClientMockGetObject) Expect(p context.Context, p1 insolar.Reference) *mClientMockGetObject {
	m.mock.GetObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetObjectExpectation{}
	}
	m.mainExpectation.input = &ClientMockGetObjectInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.GetObject
func (m *mClientMockGetObject) Return(r ObjectDescriptor, r1 error) *ClientMock {
	m.mock.GetObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetObjectExpectation{}
	}
	m.mainExpectation.result = &ClientMockGetObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.GetObject is expected once
func (m *mClientMockGetObject) ExpectOnce(p context.Context, p1 insolar.Reference) *ClientMockGetObjectExpectation {
	m.mock.GetObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockGetObjectExpectation{}
	expectation.input = &ClientMockGetObjectInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockGetObjectExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ClientMockGetObjectResult{r, r1}
}

//Set uses given function f as a mock of Client.GetObject method
func (m *mClientMockGetObject) Set(f func(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectFunc = f
	return m.mock
}

//GetObject implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) GetObject(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.GetObjectPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectCounter, 1)

	if len(m.GetObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.GetObject. %v %v", p, p1)
			return
		}

		input := m.GetObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockGetObjectInput{p, p1}, "Client.GetObject got unexpected parameters")

		result := m.GetObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetObjectMock.mainExpectation != nil {

		input := m.GetObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockGetObjectInput{p, p1}, "Client.GetObject got unexpected parameters")
		}

		result := m.GetObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.GetObject. %v %v", p, p1)
		return
	}

	return m.GetObjectFunc(p, p1)
}

//GetObjectMinimockCounter returns a count of ClientMock.GetObjectFunc invocations
func (m *ClientMock) GetObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectCounter)
}

//GetObjectMinimockPreCounter returns the value of ClientMock.GetObject invocations
func (m *ClientMock) GetObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectPreCounter)
}

//GetObjectFinished returns true if mock invocations count is ok
func (m *ClientMock) GetObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetObjectCounter) == uint64(len(m.GetObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetObjectFunc != nil {
		return atomic.LoadUint64(&m.GetObjectCounter) > 0
	}

	return true
}

type mClientMockGetPendingRequest struct {
	mock              *ClientMock
	mainExpectation   *ClientMockGetPendingRequestExpectation
	expectationSeries []*ClientMockGetPendingRequestExpectation
}

type ClientMockGetPendingRequestExpectation struct {
	input  *ClientMockGetPendingRequestInput
	result *ClientMockGetPendingRequestResult
}

type ClientMockGetPendingRequestInput struct {
	p  context.Context
	p1 insolar.ID
}

type ClientMockGetPendingRequestResult struct {
	r  *insolar.Reference
	r1 *record.IncomingRequest
	r2 error
}

//Expect specifies that invocation of Client.GetPendingRequest is expected from 1 to Infinity times
func (m *mClientMockGetPendingRequest) Expect(p context.Context, p1 insolar.ID) *mClientMockGetPendingRequest {
	m.mock.GetPendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetPendingRequestExpectation{}
	}
	m.mainExpectation.input = &ClientMockGetPendingRequestInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.GetPendingRequest
func (m *mClientMockGetPendingRequest) Return(r *insolar.Reference, r1 *record.IncomingRequest, r2 error) *ClientMock {
	m.mock.GetPendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetPendingRequestExpectation{}
	}
	m.mainExpectation.result = &ClientMockGetPendingRequestResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.GetPendingRequest is expected once
func (m *mClientMockGetPendingRequest) ExpectOnce(p context.Context, p1 insolar.ID) *ClientMockGetPendingRequestExpectation {
	m.mock.GetPendingRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockGetPendingRequestExpectation{}
	expectation.input = &ClientMockGetPendingRequestInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockGetPendingRequestExpectation) Return(r *insolar.Reference, r1 *record.IncomingRequest, r2 error) {
	e.result = &ClientMockGetPendingRequestResult{r, r1, r2}
}

//Set uses given function f as a mock of Client.GetPendingRequest method
func (m *mClientMockGetPendingRequest) Set(f func(p context.Context, p1 insolar.ID) (r *insolar.Reference, r1 *record.IncomingRequest, r2 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPendingRequestFunc = f
	return m.mock
}

//GetPendingRequest implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) GetPendingRequest(p context.Context, p1 insolar.ID) (r *insolar.Reference, r1 *record.IncomingRequest, r2 error) {
	counter := atomic.AddUint64(&m.GetPendingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.GetPendingRequestCounter, 1)

	if len(m.GetPendingRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPendingRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.GetPendingRequest. %v %v", p, p1)
			return
		}

		input := m.GetPendingRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockGetPendingRequestInput{p, p1}, "Client.GetPendingRequest got unexpected parameters")

		result := m.GetPendingRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetPendingRequest")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetPendingRequestMock.mainExpectation != nil {

		input := m.GetPendingRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockGetPendingRequestInput{p, p1}, "Client.GetPendingRequest got unexpected parameters")
		}

		result := m.GetPendingRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.GetPendingRequest")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetPendingRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.GetPendingRequest. %v %v", p, p1)
		return
	}

	return m.GetPendingRequestFunc(p, p1)
}

//GetPendingRequestMinimockCounter returns a count of ClientMock.GetPendingRequestFunc invocations
func (m *ClientMock) GetPendingRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPendingRequestCounter)
}

//GetPendingRequestMinimockPreCounter returns the value of ClientMock.GetPendingRequest invocations
func (m *ClientMock) GetPendingRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPendingRequestPreCounter)
}

//GetPendingRequestFinished returns true if mock invocations count is ok
func (m *ClientMock) GetPendingRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPendingRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPendingRequestCounter) == uint64(len(m.GetPendingRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPendingRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPendingRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPendingRequestFunc != nil {
		return atomic.LoadUint64(&m.GetPendingRequestCounter) > 0
	}

	return true
}

type mClientMockHasPendingRequests struct {
	mock              *ClientMock
	mainExpectation   *ClientMockHasPendingRequestsExpectation
	expectationSeries []*ClientMockHasPendingRequestsExpectation
}

type ClientMockHasPendingRequestsExpectation struct {
	input  *ClientMockHasPendingRequestsInput
	result *ClientMockHasPendingRequestsResult
}

type ClientMockHasPendingRequestsInput struct {
	p  context.Context
	p1 insolar.Reference
}

type ClientMockHasPendingRequestsResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of Client.HasPendingRequests is expected from 1 to Infinity times
func (m *mClientMockHasPendingRequests) Expect(p context.Context, p1 insolar.Reference) *mClientMockHasPendingRequests {
	m.mock.HasPendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockHasPendingRequestsExpectation{}
	}
	m.mainExpectation.input = &ClientMockHasPendingRequestsInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.HasPendingRequests
func (m *mClientMockHasPendingRequests) Return(r bool, r1 error) *ClientMock {
	m.mock.HasPendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockHasPendingRequestsExpectation{}
	}
	m.mainExpectation.result = &ClientMockHasPendingRequestsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.HasPendingRequests is expected once
func (m *mClientMockHasPendingRequests) ExpectOnce(p context.Context, p1 insolar.Reference) *ClientMockHasPendingRequestsExpectation {
	m.mock.HasPendingRequestsFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockHasPendingRequestsExpectation{}
	expectation.input = &ClientMockHasPendingRequestsInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockHasPendingRequestsExpectation) Return(r bool, r1 error) {
	e.result = &ClientMockHasPendingRequestsResult{r, r1}
}

//Set uses given function f as a mock of Client.HasPendingRequests method
func (m *mClientMockHasPendingRequests) Set(f func(p context.Context, p1 insolar.Reference) (r bool, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasPendingRequestsFunc = f
	return m.mock
}

//HasPendingRequests implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) HasPendingRequests(p context.Context, p1 insolar.Reference) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.HasPendingRequestsPreCounter, 1)
	defer atomic.AddUint64(&m.HasPendingRequestsCounter, 1)

	if len(m.HasPendingRequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasPendingRequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.HasPendingRequests. %v %v", p, p1)
			return
		}

		input := m.HasPendingRequestsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockHasPendingRequestsInput{p, p1}, "Client.HasPendingRequests got unexpected parameters")

		result := m.HasPendingRequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.HasPendingRequests")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HasPendingRequestsMock.mainExpectation != nil {

		input := m.HasPendingRequestsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockHasPendingRequestsInput{p, p1}, "Client.HasPendingRequests got unexpected parameters")
		}

		result := m.HasPendingRequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.HasPendingRequests")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HasPendingRequestsFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.HasPendingRequests. %v %v", p, p1)
		return
	}

	return m.HasPendingRequestsFunc(p, p1)
}

//HasPendingRequestsMinimockCounter returns a count of ClientMock.HasPendingRequestsFunc invocations
func (m *ClientMock) HasPendingRequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasPendingRequestsCounter)
}

//HasPendingRequestsMinimockPreCounter returns the value of ClientMock.HasPendingRequests invocations
func (m *ClientMock) HasPendingRequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasPendingRequestsPreCounter)
}

//HasPendingRequestsFinished returns true if mock invocations count is ok
func (m *ClientMock) HasPendingRequestsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HasPendingRequestsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HasPendingRequestsCounter) == uint64(len(m.HasPendingRequestsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HasPendingRequestsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HasPendingRequestsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HasPendingRequestsFunc != nil {
		return atomic.LoadUint64(&m.HasPendingRequestsCounter) > 0
	}

	return true
}

type mClientMockInjectCodeDescriptor struct {
	mock              *ClientMock
	mainExpectation   *ClientMockInjectCodeDescriptorExpectation
	expectationSeries []*ClientMockInjectCodeDescriptorExpectation
}

type ClientMockInjectCodeDescriptorExpectation struct {
	input *ClientMockInjectCodeDescriptorInput
}

type ClientMockInjectCodeDescriptorInput struct {
	p  insolar.Reference
	p1 CodeDescriptor
}

//Expect specifies that invocation of Client.InjectCodeDescriptor is expected from 1 to Infinity times
func (m *mClientMockInjectCodeDescriptor) Expect(p insolar.Reference, p1 CodeDescriptor) *mClientMockInjectCodeDescriptor {
	m.mock.InjectCodeDescriptorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockInjectCodeDescriptorExpectation{}
	}
	m.mainExpectation.input = &ClientMockInjectCodeDescriptorInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.InjectCodeDescriptor
func (m *mClientMockInjectCodeDescriptor) Return() *ClientMock {
	m.mock.InjectCodeDescriptorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockInjectCodeDescriptorExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Client.InjectCodeDescriptor is expected once
func (m *mClientMockInjectCodeDescriptor) ExpectOnce(p insolar.Reference, p1 CodeDescriptor) *ClientMockInjectCodeDescriptorExpectation {
	m.mock.InjectCodeDescriptorFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockInjectCodeDescriptorExpectation{}
	expectation.input = &ClientMockInjectCodeDescriptorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Client.InjectCodeDescriptor method
func (m *mClientMockInjectCodeDescriptor) Set(f func(p insolar.Reference, p1 CodeDescriptor)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InjectCodeDescriptorFunc = f
	return m.mock
}

//InjectCodeDescriptor implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) InjectCodeDescriptor(p insolar.Reference, p1 CodeDescriptor) {
	counter := atomic.AddUint64(&m.InjectCodeDescriptorPreCounter, 1)
	defer atomic.AddUint64(&m.InjectCodeDescriptorCounter, 1)

	if len(m.InjectCodeDescriptorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InjectCodeDescriptorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.InjectCodeDescriptor. %v %v", p, p1)
			return
		}

		input := m.InjectCodeDescriptorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockInjectCodeDescriptorInput{p, p1}, "Client.InjectCodeDescriptor got unexpected parameters")

		return
	}

	if m.InjectCodeDescriptorMock.mainExpectation != nil {

		input := m.InjectCodeDescriptorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockInjectCodeDescriptorInput{p, p1}, "Client.InjectCodeDescriptor got unexpected parameters")
		}

		return
	}

	if m.InjectCodeDescriptorFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.InjectCodeDescriptor. %v %v", p, p1)
		return
	}

	m.InjectCodeDescriptorFunc(p, p1)
}

//InjectCodeDescriptorMinimockCounter returns a count of ClientMock.InjectCodeDescriptorFunc invocations
func (m *ClientMock) InjectCodeDescriptorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InjectCodeDescriptorCounter)
}

//InjectCodeDescriptorMinimockPreCounter returns the value of ClientMock.InjectCodeDescriptor invocations
func (m *ClientMock) InjectCodeDescriptorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InjectCodeDescriptorPreCounter)
}

//InjectCodeDescriptorFinished returns true if mock invocations count is ok
func (m *ClientMock) InjectCodeDescriptorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InjectCodeDescriptorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InjectCodeDescriptorCounter) == uint64(len(m.InjectCodeDescriptorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InjectCodeDescriptorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InjectCodeDescriptorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InjectCodeDescriptorFunc != nil {
		return atomic.LoadUint64(&m.InjectCodeDescriptorCounter) > 0
	}

	return true
}

type mClientMockInjectFinish struct {
	mock              *ClientMock
	mainExpectation   *ClientMockInjectFinishExpectation
	expectationSeries []*ClientMockInjectFinishExpectation
}

type ClientMockInjectFinishExpectation struct {
}

//Expect specifies that invocation of Client.InjectFinish is expected from 1 to Infinity times
func (m *mClientMockInjectFinish) Expect() *mClientMockInjectFinish {
	m.mock.InjectFinishFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockInjectFinishExpectation{}
	}

	return m
}

//Return specifies results of invocation of Client.InjectFinish
func (m *mClientMockInjectFinish) Return() *ClientMock {
	m.mock.InjectFinishFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockInjectFinishExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Client.InjectFinish is expected once
func (m *mClientMockInjectFinish) ExpectOnce() *ClientMockInjectFinishExpectation {
	m.mock.InjectFinishFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockInjectFinishExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Client.InjectFinish method
func (m *mClientMockInjectFinish) Set(f func()) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InjectFinishFunc = f
	return m.mock
}

//InjectFinish implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) InjectFinish() {
	counter := atomic.AddUint64(&m.InjectFinishPreCounter, 1)
	defer atomic.AddUint64(&m.InjectFinishCounter, 1)

	if len(m.InjectFinishMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InjectFinishMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.InjectFinish.")
			return
		}

		return
	}

	if m.InjectFinishMock.mainExpectation != nil {

		return
	}

	if m.InjectFinishFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.InjectFinish.")
		return
	}

	m.InjectFinishFunc()
}

//InjectFinishMinimockCounter returns a count of ClientMock.InjectFinishFunc invocations
func (m *ClientMock) InjectFinishMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InjectFinishCounter)
}

//InjectFinishMinimockPreCounter returns the value of ClientMock.InjectFinish invocations
func (m *ClientMock) InjectFinishMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InjectFinishPreCounter)
}

//InjectFinishFinished returns true if mock invocations count is ok
func (m *ClientMock) InjectFinishFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InjectFinishMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InjectFinishCounter) == uint64(len(m.InjectFinishMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InjectFinishMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InjectFinishCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InjectFinishFunc != nil {
		return atomic.LoadUint64(&m.InjectFinishCounter) > 0
	}

	return true
}

type mClientMockInjectObjectDescriptor struct {
	mock              *ClientMock
	mainExpectation   *ClientMockInjectObjectDescriptorExpectation
	expectationSeries []*ClientMockInjectObjectDescriptorExpectation
}

type ClientMockInjectObjectDescriptorExpectation struct {
	input *ClientMockInjectObjectDescriptorInput
}

type ClientMockInjectObjectDescriptorInput struct {
	p  insolar.Reference
	p1 ObjectDescriptor
}

//Expect specifies that invocation of Client.InjectObjectDescriptor is expected from 1 to Infinity times
func (m *mClientMockInjectObjectDescriptor) Expect(p insolar.Reference, p1 ObjectDescriptor) *mClientMockInjectObjectDescriptor {
	m.mock.InjectObjectDescriptorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockInjectObjectDescriptorExpectation{}
	}
	m.mainExpectation.input = &ClientMockInjectObjectDescriptorInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.InjectObjectDescriptor
func (m *mClientMockInjectObjectDescriptor) Return() *ClientMock {
	m.mock.InjectObjectDescriptorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockInjectObjectDescriptorExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Client.InjectObjectDescriptor is expected once
func (m *mClientMockInjectObjectDescriptor) ExpectOnce(p insolar.Reference, p1 ObjectDescriptor) *ClientMockInjectObjectDescriptorExpectation {
	m.mock.InjectObjectDescriptorFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockInjectObjectDescriptorExpectation{}
	expectation.input = &ClientMockInjectObjectDescriptorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Client.InjectObjectDescriptor method
func (m *mClientMockInjectObjectDescriptor) Set(f func(p insolar.Reference, p1 ObjectDescriptor)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InjectObjectDescriptorFunc = f
	return m.mock
}

//InjectObjectDescriptor implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) InjectObjectDescriptor(p insolar.Reference, p1 ObjectDescriptor) {
	counter := atomic.AddUint64(&m.InjectObjectDescriptorPreCounter, 1)
	defer atomic.AddUint64(&m.InjectObjectDescriptorCounter, 1)

	if len(m.InjectObjectDescriptorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InjectObjectDescriptorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.InjectObjectDescriptor. %v %v", p, p1)
			return
		}

		input := m.InjectObjectDescriptorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockInjectObjectDescriptorInput{p, p1}, "Client.InjectObjectDescriptor got unexpected parameters")

		return
	}

	if m.InjectObjectDescriptorMock.mainExpectation != nil {

		input := m.InjectObjectDescriptorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockInjectObjectDescriptorInput{p, p1}, "Client.InjectObjectDescriptor got unexpected parameters")
		}

		return
	}

	if m.InjectObjectDescriptorFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.InjectObjectDescriptor. %v %v", p, p1)
		return
	}

	m.InjectObjectDescriptorFunc(p, p1)
}

//InjectObjectDescriptorMinimockCounter returns a count of ClientMock.InjectObjectDescriptorFunc invocations
func (m *ClientMock) InjectObjectDescriptorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InjectObjectDescriptorCounter)
}

//InjectObjectDescriptorMinimockPreCounter returns the value of ClientMock.InjectObjectDescriptor invocations
func (m *ClientMock) InjectObjectDescriptorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InjectObjectDescriptorPreCounter)
}

//InjectObjectDescriptorFinished returns true if mock invocations count is ok
func (m *ClientMock) InjectObjectDescriptorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InjectObjectDescriptorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InjectObjectDescriptorCounter) == uint64(len(m.InjectObjectDescriptorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InjectObjectDescriptorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InjectObjectDescriptorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InjectObjectDescriptorFunc != nil {
		return atomic.LoadUint64(&m.InjectObjectDescriptorCounter) > 0
	}

	return true
}

type mClientMockRegisterIncomingRequest struct {
	mock              *ClientMock
	mainExpectation   *ClientMockRegisterIncomingRequestExpectation
	expectationSeries []*ClientMockRegisterIncomingRequestExpectation
}

type ClientMockRegisterIncomingRequestExpectation struct {
	input  *ClientMockRegisterIncomingRequestInput
	result *ClientMockRegisterIncomingRequestResult
}

type ClientMockRegisterIncomingRequestInput struct {
	p  context.Context
	p1 *record.IncomingRequest
}

type ClientMockRegisterIncomingRequestResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Client.RegisterIncomingRequest is expected from 1 to Infinity times
func (m *mClientMockRegisterIncomingRequest) Expect(p context.Context, p1 *record.IncomingRequest) *mClientMockRegisterIncomingRequest {
	m.mock.RegisterIncomingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterIncomingRequestExpectation{}
	}
	m.mainExpectation.input = &ClientMockRegisterIncomingRequestInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.RegisterIncomingRequest
func (m *mClientMockRegisterIncomingRequest) Return(r *insolar.ID, r1 error) *ClientMock {
	m.mock.RegisterIncomingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterIncomingRequestExpectation{}
	}
	m.mainExpectation.result = &ClientMockRegisterIncomingRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.RegisterIncomingRequest is expected once
func (m *mClientMockRegisterIncomingRequest) ExpectOnce(p context.Context, p1 *record.IncomingRequest) *ClientMockRegisterIncomingRequestExpectation {
	m.mock.RegisterIncomingRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockRegisterIncomingRequestExpectation{}
	expectation.input = &ClientMockRegisterIncomingRequestInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockRegisterIncomingRequestExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ClientMockRegisterIncomingRequestResult{r, r1}
}

//Set uses given function f as a mock of Client.RegisterIncomingRequest method
func (m *mClientMockRegisterIncomingRequest) Set(f func(p context.Context, p1 *record.IncomingRequest) (r *insolar.ID, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterIncomingRequestFunc = f
	return m.mock
}

//RegisterIncomingRequest implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) RegisterIncomingRequest(p context.Context, p1 *record.IncomingRequest) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterIncomingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterIncomingRequestCounter, 1)

	if len(m.RegisterIncomingRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterIncomingRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.RegisterIncomingRequest. %v %v", p, p1)
			return
		}

		input := m.RegisterIncomingRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockRegisterIncomingRequestInput{p, p1}, "Client.RegisterIncomingRequest got unexpected parameters")

		result := m.RegisterIncomingRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterIncomingRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterIncomingRequestMock.mainExpectation != nil {

		input := m.RegisterIncomingRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockRegisterIncomingRequestInput{p, p1}, "Client.RegisterIncomingRequest got unexpected parameters")
		}

		result := m.RegisterIncomingRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterIncomingRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterIncomingRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.RegisterIncomingRequest. %v %v", p, p1)
		return
	}

	return m.RegisterIncomingRequestFunc(p, p1)
}

//RegisterIncomingRequestMinimockCounter returns a count of ClientMock.RegisterIncomingRequestFunc invocations
func (m *ClientMock) RegisterIncomingRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterIncomingRequestCounter)
}

//RegisterIncomingRequestMinimockPreCounter returns the value of ClientMock.RegisterIncomingRequest invocations
func (m *ClientMock) RegisterIncomingRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterIncomingRequestPreCounter)
}

//RegisterIncomingRequestFinished returns true if mock invocations count is ok
func (m *ClientMock) RegisterIncomingRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterIncomingRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterIncomingRequestCounter) == uint64(len(m.RegisterIncomingRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterIncomingRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterIncomingRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterIncomingRequestFunc != nil {
		return atomic.LoadUint64(&m.RegisterIncomingRequestCounter) > 0
	}

	return true
}

type mClientMockRegisterOutgoingRequest struct {
	mock              *ClientMock
	mainExpectation   *ClientMockRegisterOutgoingRequestExpectation
	expectationSeries []*ClientMockRegisterOutgoingRequestExpectation
}

type ClientMockRegisterOutgoingRequestExpectation struct {
	input  *ClientMockRegisterOutgoingRequestInput
	result *ClientMockRegisterOutgoingRequestResult
}

type ClientMockRegisterOutgoingRequestInput struct {
	p  context.Context
	p1 *record.OutgoingRequest
}

type ClientMockRegisterOutgoingRequestResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Client.RegisterOutgoingRequest is expected from 1 to Infinity times
func (m *mClientMockRegisterOutgoingRequest) Expect(p context.Context, p1 *record.OutgoingRequest) *mClientMockRegisterOutgoingRequest {
	m.mock.RegisterOutgoingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterOutgoingRequestExpectation{}
	}
	m.mainExpectation.input = &ClientMockRegisterOutgoingRequestInput{p, p1}
	return m
}

//Return specifies results of invocation of Client.RegisterOutgoingRequest
func (m *mClientMockRegisterOutgoingRequest) Return(r *insolar.ID, r1 error) *ClientMock {
	m.mock.RegisterOutgoingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterOutgoingRequestExpectation{}
	}
	m.mainExpectation.result = &ClientMockRegisterOutgoingRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.RegisterOutgoingRequest is expected once
func (m *mClientMockRegisterOutgoingRequest) ExpectOnce(p context.Context, p1 *record.OutgoingRequest) *ClientMockRegisterOutgoingRequestExpectation {
	m.mock.RegisterOutgoingRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockRegisterOutgoingRequestExpectation{}
	expectation.input = &ClientMockRegisterOutgoingRequestInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockRegisterOutgoingRequestExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ClientMockRegisterOutgoingRequestResult{r, r1}
}

//Set uses given function f as a mock of Client.RegisterOutgoingRequest method
func (m *mClientMockRegisterOutgoingRequest) Set(f func(p context.Context, p1 *record.OutgoingRequest) (r *insolar.ID, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterOutgoingRequestFunc = f
	return m.mock
}

//RegisterOutgoingRequest implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) RegisterOutgoingRequest(p context.Context, p1 *record.OutgoingRequest) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterOutgoingRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterOutgoingRequestCounter, 1)

	if len(m.RegisterOutgoingRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterOutgoingRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.RegisterOutgoingRequest. %v %v", p, p1)
			return
		}

		input := m.RegisterOutgoingRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockRegisterOutgoingRequestInput{p, p1}, "Client.RegisterOutgoingRequest got unexpected parameters")

		result := m.RegisterOutgoingRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterOutgoingRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterOutgoingRequestMock.mainExpectation != nil {

		input := m.RegisterOutgoingRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockRegisterOutgoingRequestInput{p, p1}, "Client.RegisterOutgoingRequest got unexpected parameters")
		}

		result := m.RegisterOutgoingRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterOutgoingRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterOutgoingRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.RegisterOutgoingRequest. %v %v", p, p1)
		return
	}

	return m.RegisterOutgoingRequestFunc(p, p1)
}

//RegisterOutgoingRequestMinimockCounter returns a count of ClientMock.RegisterOutgoingRequestFunc invocations
func (m *ClientMock) RegisterOutgoingRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterOutgoingRequestCounter)
}

//RegisterOutgoingRequestMinimockPreCounter returns the value of ClientMock.RegisterOutgoingRequest invocations
func (m *ClientMock) RegisterOutgoingRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterOutgoingRequestPreCounter)
}

//RegisterOutgoingRequestFinished returns true if mock invocations count is ok
func (m *ClientMock) RegisterOutgoingRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterOutgoingRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterOutgoingRequestCounter) == uint64(len(m.RegisterOutgoingRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterOutgoingRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterOutgoingRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterOutgoingRequestFunc != nil {
		return atomic.LoadUint64(&m.RegisterOutgoingRequestCounter) > 0
	}

	return true
}

type mClientMockRegisterResult struct {
	mock              *ClientMock
	mainExpectation   *ClientMockRegisterResultExpectation
	expectationSeries []*ClientMockRegisterResultExpectation
}

type ClientMockRegisterResultExpectation struct {
	input  *ClientMockRegisterResultInput
	result *ClientMockRegisterResultResult
}

type ClientMockRegisterResultInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 RequestResult
}

type ClientMockRegisterResultResult struct {
	r error
}

//Expect specifies that invocation of Client.RegisterResult is expected from 1 to Infinity times
func (m *mClientMockRegisterResult) Expect(p context.Context, p1 insolar.Reference, p2 RequestResult) *mClientMockRegisterResult {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterResultExpectation{}
	}
	m.mainExpectation.input = &ClientMockRegisterResultInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Client.RegisterResult
func (m *mClientMockRegisterResult) Return(r error) *ClientMock {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterResultExpectation{}
	}
	m.mainExpectation.result = &ClientMockRegisterResultResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.RegisterResult is expected once
func (m *mClientMockRegisterResult) ExpectOnce(p context.Context, p1 insolar.Reference, p2 RequestResult) *ClientMockRegisterResultExpectation {
	m.mock.RegisterResultFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockRegisterResultExpectation{}
	expectation.input = &ClientMockRegisterResultInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockRegisterResultExpectation) Return(r error) {
	e.result = &ClientMockRegisterResultResult{r}
}

//Set uses given function f as a mock of Client.RegisterResult method
func (m *mClientMockRegisterResult) Set(f func(p context.Context, p1 insolar.Reference, p2 RequestResult) (r error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterResultFunc = f
	return m.mock
}

//RegisterResult implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) RegisterResult(p context.Context, p1 insolar.Reference, p2 RequestResult) (r error) {
	counter := atomic.AddUint64(&m.RegisterResultPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterResultCounter, 1)

	if len(m.RegisterResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.RegisterResult. %v %v %v", p, p1, p2)
			return
		}

		input := m.RegisterResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockRegisterResultInput{p, p1, p2}, "Client.RegisterResult got unexpected parameters")

		result := m.RegisterResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterResult")
			return
		}

		r = result.r

		return
	}

	if m.RegisterResultMock.mainExpectation != nil {

		input := m.RegisterResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockRegisterResultInput{p, p1, p2}, "Client.RegisterResult got unexpected parameters")
		}

		result := m.RegisterResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterResult")
		}

		r = result.r

		return
	}

	if m.RegisterResultFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.RegisterResult. %v %v %v", p, p1, p2)
		return
	}

	return m.RegisterResultFunc(p, p1, p2)
}

//RegisterResultMinimockCounter returns a count of ClientMock.RegisterResultFunc invocations
func (m *ClientMock) RegisterResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultCounter)
}

//RegisterResultMinimockPreCounter returns the value of ClientMock.RegisterResult invocations
func (m *ClientMock) RegisterResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultPreCounter)
}

//RegisterResultFinished returns true if mock invocations count is ok
func (m *ClientMock) RegisterResultFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterResultMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterResultCounter) == uint64(len(m.RegisterResultMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterResultMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterResultCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterResultFunc != nil {
		return atomic.LoadUint64(&m.RegisterResultCounter) > 0
	}

	return true
}

type mClientMockRegisterValidation struct {
	mock              *ClientMock
	mainExpectation   *ClientMockRegisterValidationExpectation
	expectationSeries []*ClientMockRegisterValidationExpectation
}

type ClientMockRegisterValidationExpectation struct {
	input  *ClientMockRegisterValidationInput
	result *ClientMockRegisterValidationResult
}

type ClientMockRegisterValidationInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.ID
	p3 bool
	p4 []insolar.Message
}

type ClientMockRegisterValidationResult struct {
	r error
}

//Expect specifies that invocation of Client.RegisterValidation is expected from 1 to Infinity times
func (m *mClientMockRegisterValidation) Expect(p context.Context, p1 insolar.Reference, p2 insolar.ID, p3 bool, p4 []insolar.Message) *mClientMockRegisterValidation {
	m.mock.RegisterValidationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterValidationExpectation{}
	}
	m.mainExpectation.input = &ClientMockRegisterValidationInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of Client.RegisterValidation
func (m *mClientMockRegisterValidation) Return(r error) *ClientMock {
	m.mock.RegisterValidationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterValidationExpectation{}
	}
	m.mainExpectation.result = &ClientMockRegisterValidationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.RegisterValidation is expected once
func (m *mClientMockRegisterValidation) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.ID, p3 bool, p4 []insolar.Message) *ClientMockRegisterValidationExpectation {
	m.mock.RegisterValidationFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockRegisterValidationExpectation{}
	expectation.input = &ClientMockRegisterValidationInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockRegisterValidationExpectation) Return(r error) {
	e.result = &ClientMockRegisterValidationResult{r}
}

//Set uses given function f as a mock of Client.RegisterValidation method
func (m *mClientMockRegisterValidation) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.ID, p3 bool, p4 []insolar.Message) (r error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterValidationFunc = f
	return m.mock
}

//RegisterValidation implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) RegisterValidation(p context.Context, p1 insolar.Reference, p2 insolar.ID, p3 bool, p4 []insolar.Message) (r error) {
	counter := atomic.AddUint64(&m.RegisterValidationPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterValidationCounter, 1)

	if len(m.RegisterValidationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterValidationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.RegisterValidation. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.RegisterValidationMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockRegisterValidationInput{p, p1, p2, p3, p4}, "Client.RegisterValidation got unexpected parameters")

		result := m.RegisterValidationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterValidation")
			return
		}

		r = result.r

		return
	}

	if m.RegisterValidationMock.mainExpectation != nil {

		input := m.RegisterValidationMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockRegisterValidationInput{p, p1, p2, p3, p4}, "Client.RegisterValidation got unexpected parameters")
		}

		result := m.RegisterValidationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterValidation")
		}

		r = result.r

		return
	}

	if m.RegisterValidationFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.RegisterValidation. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.RegisterValidationFunc(p, p1, p2, p3, p4)
}

//RegisterValidationMinimockCounter returns a count of ClientMock.RegisterValidationFunc invocations
func (m *ClientMock) RegisterValidationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterValidationCounter)
}

//RegisterValidationMinimockPreCounter returns the value of ClientMock.RegisterValidation invocations
func (m *ClientMock) RegisterValidationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterValidationPreCounter)
}

//RegisterValidationFinished returns true if mock invocations count is ok
func (m *ClientMock) RegisterValidationFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterValidationMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterValidationCounter) == uint64(len(m.RegisterValidationMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterValidationMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterValidationCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterValidationFunc != nil {
		return atomic.LoadUint64(&m.RegisterValidationCounter) > 0
	}

	return true
}

type mClientMockState struct {
	mock              *ClientMock
	mainExpectation   *ClientMockStateExpectation
	expectationSeries []*ClientMockStateExpectation
}

type ClientMockStateExpectation struct {
	result *ClientMockStateResult
}

type ClientMockStateResult struct {
	r []byte
}

//Expect specifies that invocation of Client.State is expected from 1 to Infinity times
func (m *mClientMockState) Expect() *mClientMockState {
	m.mock.StateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of Client.State
func (m *mClientMockState) Return(r []byte) *ClientMock {
	m.mock.StateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockStateExpectation{}
	}
	m.mainExpectation.result = &ClientMockStateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.State is expected once
func (m *mClientMockState) ExpectOnce() *ClientMockStateExpectation {
	m.mock.StateFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockStateExpectation) Return(r []byte) {
	e.result = &ClientMockStateResult{r}
}

//Set uses given function f as a mock of Client.State method
func (m *mClientMockState) Set(f func() (r []byte)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StateFunc = f
	return m.mock
}

//State implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) State() (r []byte) {
	counter := atomic.AddUint64(&m.StatePreCounter, 1)
	defer atomic.AddUint64(&m.StateCounter, 1)

	if len(m.StateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.State.")
			return
		}

		result := m.StateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.State")
			return
		}

		r = result.r

		return
	}

	if m.StateMock.mainExpectation != nil {

		result := m.StateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.State")
		}

		r = result.r

		return
	}

	if m.StateFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.State.")
		return
	}

	return m.StateFunc()
}

//StateMinimockCounter returns a count of ClientMock.StateFunc invocations
func (m *ClientMock) StateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StateCounter)
}

//StateMinimockPreCounter returns the value of ClientMock.State invocations
func (m *ClientMock) StateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StatePreCounter)
}

//StateFinished returns true if mock invocations count is ok
func (m *ClientMock) StateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StateCounter) == uint64(len(m.StateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StateFunc != nil {
		return atomic.LoadUint64(&m.StateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ClientMock) ValidateCallCounters() {

	if !m.ActivatePrototypeFinished() {
		m.t.Fatal("Expected call to ClientMock.ActivatePrototype")
	}

	if !m.DeployCodeFinished() {
		m.t.Fatal("Expected call to ClientMock.DeployCode")
	}

	if !m.GetChildrenFinished() {
		m.t.Fatal("Expected call to ClientMock.GetChildren")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to ClientMock.GetCode")
	}

	if !m.GetDelegateFinished() {
		m.t.Fatal("Expected call to ClientMock.GetDelegate")
	}

	if !m.GetObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.GetObject")
	}

	if !m.GetPendingRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.GetPendingRequest")
	}

	if !m.HasPendingRequestsFinished() {
		m.t.Fatal("Expected call to ClientMock.HasPendingRequests")
	}

	if !m.InjectCodeDescriptorFinished() {
		m.t.Fatal("Expected call to ClientMock.InjectCodeDescriptor")
	}

	if !m.InjectFinishFinished() {
		m.t.Fatal("Expected call to ClientMock.InjectFinish")
	}

	if !m.InjectObjectDescriptorFinished() {
		m.t.Fatal("Expected call to ClientMock.InjectObjectDescriptor")
	}

	if !m.RegisterIncomingRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterIncomingRequest")
	}

	if !m.RegisterOutgoingRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterOutgoingRequest")
	}

	if !m.RegisterResultFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterResult")
	}

	if !m.RegisterValidationFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterValidation")
	}

	if !m.StateFinished() {
		m.t.Fatal("Expected call to ClientMock.State")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ClientMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ClientMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ClientMock) MinimockFinish() {

	if !m.ActivatePrototypeFinished() {
		m.t.Fatal("Expected call to ClientMock.ActivatePrototype")
	}

	if !m.DeployCodeFinished() {
		m.t.Fatal("Expected call to ClientMock.DeployCode")
	}

	if !m.GetChildrenFinished() {
		m.t.Fatal("Expected call to ClientMock.GetChildren")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to ClientMock.GetCode")
	}

	if !m.GetDelegateFinished() {
		m.t.Fatal("Expected call to ClientMock.GetDelegate")
	}

	if !m.GetObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.GetObject")
	}

	if !m.GetPendingRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.GetPendingRequest")
	}

	if !m.HasPendingRequestsFinished() {
		m.t.Fatal("Expected call to ClientMock.HasPendingRequests")
	}

	if !m.InjectCodeDescriptorFinished() {
		m.t.Fatal("Expected call to ClientMock.InjectCodeDescriptor")
	}

	if !m.InjectFinishFinished() {
		m.t.Fatal("Expected call to ClientMock.InjectFinish")
	}

	if !m.InjectObjectDescriptorFinished() {
		m.t.Fatal("Expected call to ClientMock.InjectObjectDescriptor")
	}

	if !m.RegisterIncomingRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterIncomingRequest")
	}

	if !m.RegisterOutgoingRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterOutgoingRequest")
	}

	if !m.RegisterResultFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterResult")
	}

	if !m.RegisterValidationFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterValidation")
	}

	if !m.StateFinished() {
		m.t.Fatal("Expected call to ClientMock.State")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ClientMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ClientMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ActivatePrototypeFinished()
		ok = ok && m.DeployCodeFinished()
		ok = ok && m.GetChildrenFinished()
		ok = ok && m.GetCodeFinished()
		ok = ok && m.GetDelegateFinished()
		ok = ok && m.GetObjectFinished()
		ok = ok && m.GetPendingRequestFinished()
		ok = ok && m.HasPendingRequestsFinished()
		ok = ok && m.InjectCodeDescriptorFinished()
		ok = ok && m.InjectFinishFinished()
		ok = ok && m.InjectObjectDescriptorFinished()
		ok = ok && m.RegisterIncomingRequestFinished()
		ok = ok && m.RegisterOutgoingRequestFinished()
		ok = ok && m.RegisterResultFinished()
		ok = ok && m.RegisterValidationFinished()
		ok = ok && m.StateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ActivatePrototypeFinished() {
				m.t.Error("Expected call to ClientMock.ActivatePrototype")
			}

			if !m.DeployCodeFinished() {
				m.t.Error("Expected call to ClientMock.DeployCode")
			}

			if !m.GetChildrenFinished() {
				m.t.Error("Expected call to ClientMock.GetChildren")
			}

			if !m.GetCodeFinished() {
				m.t.Error("Expected call to ClientMock.GetCode")
			}

			if !m.GetDelegateFinished() {
				m.t.Error("Expected call to ClientMock.GetDelegate")
			}

			if !m.GetObjectFinished() {
				m.t.Error("Expected call to ClientMock.GetObject")
			}

			if !m.GetPendingRequestFinished() {
				m.t.Error("Expected call to ClientMock.GetPendingRequest")
			}

			if !m.HasPendingRequestsFinished() {
				m.t.Error("Expected call to ClientMock.HasPendingRequests")
			}

			if !m.InjectCodeDescriptorFinished() {
				m.t.Error("Expected call to ClientMock.InjectCodeDescriptor")
			}

			if !m.InjectFinishFinished() {
				m.t.Error("Expected call to ClientMock.InjectFinish")
			}

			if !m.InjectObjectDescriptorFinished() {
				m.t.Error("Expected call to ClientMock.InjectObjectDescriptor")
			}

			if !m.RegisterIncomingRequestFinished() {
				m.t.Error("Expected call to ClientMock.RegisterIncomingRequest")
			}

			if !m.RegisterOutgoingRequestFinished() {
				m.t.Error("Expected call to ClientMock.RegisterOutgoingRequest")
			}

			if !m.RegisterResultFinished() {
				m.t.Error("Expected call to ClientMock.RegisterResult")
			}

			if !m.RegisterValidationFinished() {
				m.t.Error("Expected call to ClientMock.RegisterValidation")
			}

			if !m.StateFinished() {
				m.t.Error("Expected call to ClientMock.State")
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
func (m *ClientMock) AllMocksCalled() bool {

	if !m.ActivatePrototypeFinished() {
		return false
	}

	if !m.DeployCodeFinished() {
		return false
	}

	if !m.GetChildrenFinished() {
		return false
	}

	if !m.GetCodeFinished() {
		return false
	}

	if !m.GetDelegateFinished() {
		return false
	}

	if !m.GetObjectFinished() {
		return false
	}

	if !m.GetPendingRequestFinished() {
		return false
	}

	if !m.HasPendingRequestsFinished() {
		return false
	}

	if !m.InjectCodeDescriptorFinished() {
		return false
	}

	if !m.InjectFinishFinished() {
		return false
	}

	if !m.InjectObjectDescriptorFinished() {
		return false
	}

	if !m.RegisterIncomingRequestFinished() {
		return false
	}

	if !m.RegisterOutgoingRequestFinished() {
		return false
	}

	if !m.RegisterResultFinished() {
		return false
	}

	if !m.RegisterValidationFinished() {
		return false
	}

	if !m.StateFinished() {
		return false
	}

	return true
}
