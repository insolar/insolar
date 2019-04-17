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

	testify_assert "github.com/stretchr/testify/assert"
)

//ClientMock implements github.com/insolar/insolar/logicrunner/artifacts.Client
type ClientMock struct {
	t minimock.Tester

	ActivateObjectFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) (r ObjectDescriptor, r1 error)
	ActivateObjectCounter    uint64
	ActivateObjectPreCounter uint64
	ActivateObjectMock       mClientMockActivateObject

	ActivatePrototypeFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 []byte) (r ObjectDescriptor, r1 error)
	ActivatePrototypeCounter    uint64
	ActivatePrototypePreCounter uint64
	ActivatePrototypeMock       mClientMockActivatePrototype

	DeactivateObjectFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor) (r *insolar.ID, r1 error)
	DeactivateObjectCounter    uint64
	DeactivateObjectPreCounter uint64
	DeactivateObjectMock       mClientMockDeactivateObject

	DeclareTypeFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error)
	DeclareTypeCounter    uint64
	DeclareTypePreCounter uint64
	DeclareTypeMock       mClientMockDeclareType

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

	GetObjectFunc       func(p context.Context, p1 insolar.Reference, p2 *insolar.ID, p3 bool) (r ObjectDescriptor, r1 error)
	GetObjectCounter    uint64
	GetObjectPreCounter uint64
	GetObjectMock       mClientMockGetObject

	GetPendingRequestFunc       func(p context.Context, p1 insolar.ID) (r insolar.Parcel, r1 error)
	GetPendingRequestCounter    uint64
	GetPendingRequestPreCounter uint64
	GetPendingRequestMock       mClientMockGetPendingRequest

	HasPendingRequestsFunc       func(p context.Context, p1 insolar.Reference) (r bool, r1 error)
	HasPendingRequestsCounter    uint64
	HasPendingRequestsPreCounter uint64
	HasPendingRequestsMock       mClientMockHasPendingRequests

	RegisterRequestFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error)
	RegisterRequestCounter    uint64
	RegisterRequestPreCounter uint64
	RegisterRequestMock       mClientMockRegisterRequest

	RegisterResultFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error)
	RegisterResultCounter    uint64
	RegisterResultPreCounter uint64
	RegisterResultMock       mClientMockRegisterResult

	RegisterValidationFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.ID, p3 bool, p4 []insolar.Message) (r error)
	RegisterValidationCounter    uint64
	RegisterValidationPreCounter uint64
	RegisterValidationMock       mClientMockRegisterValidation

	StateFunc       func() (r []byte, r1 error)
	StateCounter    uint64
	StatePreCounter uint64
	StateMock       mClientMockState

	UpdateObjectFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) (r ObjectDescriptor, r1 error)
	UpdateObjectCounter    uint64
	UpdateObjectPreCounter uint64
	UpdateObjectMock       mClientMockUpdateObject

	UpdatePrototypeFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte, p5 *insolar.Reference) (r ObjectDescriptor, r1 error)
	UpdatePrototypeCounter    uint64
	UpdatePrototypePreCounter uint64
	UpdatePrototypeMock       mClientMockUpdatePrototype
}

//NewClientMock returns a mock for github.com/insolar/insolar/logicrunner/artifacts.Client
func NewClientMock(t minimock.Tester) *ClientMock {
	m := &ClientMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ActivateObjectMock = mClientMockActivateObject{mock: m}
	m.ActivatePrototypeMock = mClientMockActivatePrototype{mock: m}
	m.DeactivateObjectMock = mClientMockDeactivateObject{mock: m}
	m.DeclareTypeMock = mClientMockDeclareType{mock: m}
	m.DeployCodeMock = mClientMockDeployCode{mock: m}
	m.GetChildrenMock = mClientMockGetChildren{mock: m}
	m.GetCodeMock = mClientMockGetCode{mock: m}
	m.GetDelegateMock = mClientMockGetDelegate{mock: m}
	m.GetObjectMock = mClientMockGetObject{mock: m}
	m.GetPendingRequestMock = mClientMockGetPendingRequest{mock: m}
	m.HasPendingRequestsMock = mClientMockHasPendingRequests{mock: m}
	m.RegisterRequestMock = mClientMockRegisterRequest{mock: m}
	m.RegisterResultMock = mClientMockRegisterResult{mock: m}
	m.RegisterValidationMock = mClientMockRegisterValidation{mock: m}
	m.StateMock = mClientMockState{mock: m}
	m.UpdateObjectMock = mClientMockUpdateObject{mock: m}
	m.UpdatePrototypeMock = mClientMockUpdatePrototype{mock: m}

	return m
}

type mClientMockActivateObject struct {
	mock              *ClientMock
	mainExpectation   *ClientMockActivateObjectExpectation
	expectationSeries []*ClientMockActivateObjectExpectation
}

type ClientMockActivateObjectExpectation struct {
	input  *ClientMockActivateObjectInput
	result *ClientMockActivateObjectResult
}

type ClientMockActivateObjectInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 insolar.Reference
	p4 insolar.Reference
	p5 bool
	p6 []byte
}

type ClientMockActivateObjectResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Client.ActivateObject is expected from 1 to Infinity times
func (m *mClientMockActivateObject) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) *mClientMockActivateObject {
	m.mock.ActivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockActivateObjectExpectation{}
	}
	m.mainExpectation.input = &ClientMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}
	return m
}

//Return specifies results of invocation of Client.ActivateObject
func (m *mClientMockActivateObject) Return(r ObjectDescriptor, r1 error) *ClientMock {
	m.mock.ActivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockActivateObjectExpectation{}
	}
	m.mainExpectation.result = &ClientMockActivateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.ActivateObject is expected once
func (m *mClientMockActivateObject) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) *ClientMockActivateObjectExpectation {
	m.mock.ActivateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockActivateObjectExpectation{}
	expectation.input = &ClientMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockActivateObjectExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ClientMockActivateObjectResult{r, r1}
}

//Set uses given function f as a mock of Client.ActivateObject method
func (m *mClientMockActivateObject) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) (r ObjectDescriptor, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ActivateObjectFunc = f
	return m.mock
}

//ActivateObject implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) ActivateObject(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.ActivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.ActivateObjectCounter, 1)

	if len(m.ActivateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ActivateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.ActivateObject. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
			return
		}

		input := m.ActivateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}, "Client.ActivateObject got unexpected parameters")

		result := m.ActivateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.ActivateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivateObjectMock.mainExpectation != nil {

		input := m.ActivateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}, "Client.ActivateObject got unexpected parameters")
		}

		result := m.ActivateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.ActivateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.ActivateObject. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
		return
	}

	return m.ActivateObjectFunc(p, p1, p2, p3, p4, p5, p6)
}

//ActivateObjectMinimockCounter returns a count of ClientMock.ActivateObjectFunc invocations
func (m *ClientMock) ActivateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ActivateObjectCounter)
}

//ActivateObjectMinimockPreCounter returns the value of ClientMock.ActivateObject invocations
func (m *ClientMock) ActivateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ActivateObjectPreCounter)
}

//ActivateObjectFinished returns true if mock invocations count is ok
func (m *ClientMock) ActivateObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ActivateObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ActivateObjectCounter) == uint64(len(m.ActivateObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ActivateObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ActivateObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ActivateObjectFunc != nil {
		return atomic.LoadUint64(&m.ActivateObjectCounter) > 0
	}

	return true
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
	p4 insolar.Reference
	p5 []byte
}

type ClientMockActivatePrototypeResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Client.ActivatePrototype is expected from 1 to Infinity times
func (m *mClientMockActivatePrototype) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 []byte) *mClientMockActivatePrototype {
	m.mock.ActivatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockActivatePrototypeExpectation{}
	}
	m.mainExpectation.input = &ClientMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}
	return m
}

//Return specifies results of invocation of Client.ActivatePrototype
func (m *mClientMockActivatePrototype) Return(r ObjectDescriptor, r1 error) *ClientMock {
	m.mock.ActivatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockActivatePrototypeExpectation{}
	}
	m.mainExpectation.result = &ClientMockActivatePrototypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.ActivatePrototype is expected once
func (m *mClientMockActivatePrototype) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 []byte) *ClientMockActivatePrototypeExpectation {
	m.mock.ActivatePrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockActivatePrototypeExpectation{}
	expectation.input = &ClientMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockActivatePrototypeExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ClientMockActivatePrototypeResult{r, r1}
}

//Set uses given function f as a mock of Client.ActivatePrototype method
func (m *mClientMockActivatePrototype) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 []byte) (r ObjectDescriptor, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ActivatePrototypeFunc = f
	return m.mock
}

//ActivatePrototype implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) ActivatePrototype(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 []byte) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.ActivatePrototypePreCounter, 1)
	defer atomic.AddUint64(&m.ActivatePrototypeCounter, 1)

	if len(m.ActivatePrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ActivatePrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.ActivatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
			return
		}

		input := m.ActivatePrototypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}, "Client.ActivatePrototype got unexpected parameters")

		result := m.ActivatePrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.ActivatePrototype")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivatePrototypeMock.mainExpectation != nil {

		input := m.ActivatePrototypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}, "Client.ActivatePrototype got unexpected parameters")
		}

		result := m.ActivatePrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.ActivatePrototype")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivatePrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.ActivatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
		return
	}

	return m.ActivatePrototypeFunc(p, p1, p2, p3, p4, p5)
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

type mClientMockDeactivateObject struct {
	mock              *ClientMock
	mainExpectation   *ClientMockDeactivateObjectExpectation
	expectationSeries []*ClientMockDeactivateObjectExpectation
}

type ClientMockDeactivateObjectExpectation struct {
	input  *ClientMockDeactivateObjectInput
	result *ClientMockDeactivateObjectResult
}

type ClientMockDeactivateObjectInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 ObjectDescriptor
}

type ClientMockDeactivateObjectResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Client.DeactivateObject is expected from 1 to Infinity times
func (m *mClientMockDeactivateObject) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor) *mClientMockDeactivateObject {
	m.mock.DeactivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockDeactivateObjectExpectation{}
	}
	m.mainExpectation.input = &ClientMockDeactivateObjectInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Client.DeactivateObject
func (m *mClientMockDeactivateObject) Return(r *insolar.ID, r1 error) *ClientMock {
	m.mock.DeactivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockDeactivateObjectExpectation{}
	}
	m.mainExpectation.result = &ClientMockDeactivateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.DeactivateObject is expected once
func (m *mClientMockDeactivateObject) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor) *ClientMockDeactivateObjectExpectation {
	m.mock.DeactivateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockDeactivateObjectExpectation{}
	expectation.input = &ClientMockDeactivateObjectInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockDeactivateObjectExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ClientMockDeactivateObjectResult{r, r1}
}

//Set uses given function f as a mock of Client.DeactivateObject method
func (m *mClientMockDeactivateObject) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor) (r *insolar.ID, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeactivateObjectFunc = f
	return m.mock
}

//DeactivateObject implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) DeactivateObject(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.DeactivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.DeactivateObjectCounter, 1)

	if len(m.DeactivateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeactivateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.DeactivateObject. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.DeactivateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockDeactivateObjectInput{p, p1, p2, p3}, "Client.DeactivateObject got unexpected parameters")

		result := m.DeactivateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.DeactivateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeactivateObjectMock.mainExpectation != nil {

		input := m.DeactivateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockDeactivateObjectInput{p, p1, p2, p3}, "Client.DeactivateObject got unexpected parameters")
		}

		result := m.DeactivateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.DeactivateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeactivateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.DeactivateObject. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.DeactivateObjectFunc(p, p1, p2, p3)
}

//DeactivateObjectMinimockCounter returns a count of ClientMock.DeactivateObjectFunc invocations
func (m *ClientMock) DeactivateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateObjectCounter)
}

//DeactivateObjectMinimockPreCounter returns the value of ClientMock.DeactivateObject invocations
func (m *ClientMock) DeactivateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateObjectPreCounter)
}

//DeactivateObjectFinished returns true if mock invocations count is ok
func (m *ClientMock) DeactivateObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeactivateObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeactivateObjectCounter) == uint64(len(m.DeactivateObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeactivateObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeactivateObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeactivateObjectFunc != nil {
		return atomic.LoadUint64(&m.DeactivateObjectCounter) > 0
	}

	return true
}

type mClientMockDeclareType struct {
	mock              *ClientMock
	mainExpectation   *ClientMockDeclareTypeExpectation
	expectationSeries []*ClientMockDeclareTypeExpectation
}

type ClientMockDeclareTypeExpectation struct {
	input  *ClientMockDeclareTypeInput
	result *ClientMockDeclareTypeResult
}

type ClientMockDeclareTypeInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 []byte
}

type ClientMockDeclareTypeResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Client.DeclareType is expected from 1 to Infinity times
func (m *mClientMockDeclareType) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) *mClientMockDeclareType {
	m.mock.DeclareTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockDeclareTypeExpectation{}
	}
	m.mainExpectation.input = &ClientMockDeclareTypeInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Client.DeclareType
func (m *mClientMockDeclareType) Return(r *insolar.ID, r1 error) *ClientMock {
	m.mock.DeclareTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockDeclareTypeExpectation{}
	}
	m.mainExpectation.result = &ClientMockDeclareTypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.DeclareType is expected once
func (m *mClientMockDeclareType) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) *ClientMockDeclareTypeExpectation {
	m.mock.DeclareTypeFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockDeclareTypeExpectation{}
	expectation.input = &ClientMockDeclareTypeInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockDeclareTypeExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ClientMockDeclareTypeResult{r, r1}
}

//Set uses given function f as a mock of Client.DeclareType method
func (m *mClientMockDeclareType) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeclareTypeFunc = f
	return m.mock
}

//DeclareType implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) DeclareType(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.DeclareTypePreCounter, 1)
	defer atomic.AddUint64(&m.DeclareTypeCounter, 1)

	if len(m.DeclareTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeclareTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.DeclareType. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.DeclareTypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockDeclareTypeInput{p, p1, p2, p3}, "Client.DeclareType got unexpected parameters")

		result := m.DeclareTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.DeclareType")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeclareTypeMock.mainExpectation != nil {

		input := m.DeclareTypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockDeclareTypeInput{p, p1, p2, p3}, "Client.DeclareType got unexpected parameters")
		}

		result := m.DeclareTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.DeclareType")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeclareTypeFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.DeclareType. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.DeclareTypeFunc(p, p1, p2, p3)
}

//DeclareTypeMinimockCounter returns a count of ClientMock.DeclareTypeFunc invocations
func (m *ClientMock) DeclareTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeclareTypeCounter)
}

//DeclareTypeMinimockPreCounter returns the value of ClientMock.DeclareType invocations
func (m *ClientMock) DeclareTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeclareTypePreCounter)
}

//DeclareTypeFinished returns true if mock invocations count is ok
func (m *ClientMock) DeclareTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeclareTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeclareTypeCounter) == uint64(len(m.DeclareTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeclareTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeclareTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeclareTypeFunc != nil {
		return atomic.LoadUint64(&m.DeclareTypeCounter) > 0
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
	p2 *insolar.ID
	p3 bool
}

type ClientMockGetObjectResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Client.GetObject is expected from 1 to Infinity times
func (m *mClientMockGetObject) Expect(p context.Context, p1 insolar.Reference, p2 *insolar.ID, p3 bool) *mClientMockGetObject {
	m.mock.GetObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetObjectExpectation{}
	}
	m.mainExpectation.input = &ClientMockGetObjectInput{p, p1, p2, p3}
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
func (m *mClientMockGetObject) ExpectOnce(p context.Context, p1 insolar.Reference, p2 *insolar.ID, p3 bool) *ClientMockGetObjectExpectation {
	m.mock.GetObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockGetObjectExpectation{}
	expectation.input = &ClientMockGetObjectInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockGetObjectExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ClientMockGetObjectResult{r, r1}
}

//Set uses given function f as a mock of Client.GetObject method
func (m *mClientMockGetObject) Set(f func(p context.Context, p1 insolar.Reference, p2 *insolar.ID, p3 bool) (r ObjectDescriptor, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectFunc = f
	return m.mock
}

//GetObject implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) GetObject(p context.Context, p1 insolar.Reference, p2 *insolar.ID, p3 bool) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.GetObjectPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectCounter, 1)

	if len(m.GetObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.GetObject. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.GetObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockGetObjectInput{p, p1, p2, p3}, "Client.GetObject got unexpected parameters")

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
			testify_assert.Equal(m.t, *input, ClientMockGetObjectInput{p, p1, p2, p3}, "Client.GetObject got unexpected parameters")
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
		m.t.Fatalf("Unexpected call to ClientMock.GetObject. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.GetObjectFunc(p, p1, p2, p3)
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
	r  insolar.Parcel
	r1 error
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
func (m *mClientMockGetPendingRequest) Return(r insolar.Parcel, r1 error) *ClientMock {
	m.mock.GetPendingRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockGetPendingRequestExpectation{}
	}
	m.mainExpectation.result = &ClientMockGetPendingRequestResult{r, r1}
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

func (e *ClientMockGetPendingRequestExpectation) Return(r insolar.Parcel, r1 error) {
	e.result = &ClientMockGetPendingRequestResult{r, r1}
}

//Set uses given function f as a mock of Client.GetPendingRequest method
func (m *mClientMockGetPendingRequest) Set(f func(p context.Context, p1 insolar.ID) (r insolar.Parcel, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPendingRequestFunc = f
	return m.mock
}

//GetPendingRequest implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) GetPendingRequest(p context.Context, p1 insolar.ID) (r insolar.Parcel, r1 error) {
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

type mClientMockRegisterRequest struct {
	mock              *ClientMock
	mainExpectation   *ClientMockRegisterRequestExpectation
	expectationSeries []*ClientMockRegisterRequestExpectation
}

type ClientMockRegisterRequestExpectation struct {
	input  *ClientMockRegisterRequestInput
	result *ClientMockRegisterRequestResult
}

type ClientMockRegisterRequestInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Parcel
}

type ClientMockRegisterRequestResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Client.RegisterRequest is expected from 1 to Infinity times
func (m *mClientMockRegisterRequest) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) *mClientMockRegisterRequest {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterRequestExpectation{}
	}
	m.mainExpectation.input = &ClientMockRegisterRequestInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Client.RegisterRequest
func (m *mClientMockRegisterRequest) Return(r *insolar.ID, r1 error) *ClientMock {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterRequestExpectation{}
	}
	m.mainExpectation.result = &ClientMockRegisterRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.RegisterRequest is expected once
func (m *mClientMockRegisterRequest) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) *ClientMockRegisterRequestExpectation {
	m.mock.RegisterRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockRegisterRequestExpectation{}
	expectation.input = &ClientMockRegisterRequestInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockRegisterRequestExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ClientMockRegisterRequestResult{r, r1}
}

//Set uses given function f as a mock of Client.RegisterRequest method
func (m *mClientMockRegisterRequest) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterRequestFunc = f
	return m.mock
}

//RegisterRequest implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) RegisterRequest(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestCounter, 1)

	if len(m.RegisterRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.RegisterRequest. %v %v %v", p, p1, p2)
			return
		}

		input := m.RegisterRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockRegisterRequestInput{p, p1, p2}, "Client.RegisterRequest got unexpected parameters")

		result := m.RegisterRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestMock.mainExpectation != nil {

		input := m.RegisterRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockRegisterRequestInput{p, p1, p2}, "Client.RegisterRequest got unexpected parameters")
		}

		result := m.RegisterRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.RegisterRequest. %v %v %v", p, p1, p2)
		return
	}

	return m.RegisterRequestFunc(p, p1, p2)
}

//RegisterRequestMinimockCounter returns a count of ClientMock.RegisterRequestFunc invocations
func (m *ClientMock) RegisterRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestCounter)
}

//RegisterRequestMinimockPreCounter returns the value of ClientMock.RegisterRequest invocations
func (m *ClientMock) RegisterRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestPreCounter)
}

//RegisterRequestFinished returns true if mock invocations count is ok
func (m *ClientMock) RegisterRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterRequestCounter) == uint64(len(m.RegisterRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterRequestFunc != nil {
		return atomic.LoadUint64(&m.RegisterRequestCounter) > 0
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
	p2 insolar.Reference
	p3 []byte
}

type ClientMockRegisterResultResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Client.RegisterResult is expected from 1 to Infinity times
func (m *mClientMockRegisterResult) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) *mClientMockRegisterResult {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterResultExpectation{}
	}
	m.mainExpectation.input = &ClientMockRegisterResultInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Client.RegisterResult
func (m *mClientMockRegisterResult) Return(r *insolar.ID, r1 error) *ClientMock {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockRegisterResultExpectation{}
	}
	m.mainExpectation.result = &ClientMockRegisterResultResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.RegisterResult is expected once
func (m *mClientMockRegisterResult) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) *ClientMockRegisterResultExpectation {
	m.mock.RegisterResultFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockRegisterResultExpectation{}
	expectation.input = &ClientMockRegisterResultInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockRegisterResultExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ClientMockRegisterResultResult{r, r1}
}

//Set uses given function f as a mock of Client.RegisterResult method
func (m *mClientMockRegisterResult) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterResultFunc = f
	return m.mock
}

//RegisterResult implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) RegisterResult(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterResultPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterResultCounter, 1)

	if len(m.RegisterResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.RegisterResult. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RegisterResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockRegisterResultInput{p, p1, p2, p3}, "Client.RegisterResult got unexpected parameters")

		result := m.RegisterResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterResult")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterResultMock.mainExpectation != nil {

		input := m.RegisterResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockRegisterResultInput{p, p1, p2, p3}, "Client.RegisterResult got unexpected parameters")
		}

		result := m.RegisterResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.RegisterResult")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterResultFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.RegisterResult. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RegisterResultFunc(p, p1, p2, p3)
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
	r  []byte
	r1 error
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
func (m *mClientMockState) Return(r []byte, r1 error) *ClientMock {
	m.mock.StateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockStateExpectation{}
	}
	m.mainExpectation.result = &ClientMockStateResult{r, r1}
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

func (e *ClientMockStateExpectation) Return(r []byte, r1 error) {
	e.result = &ClientMockStateResult{r, r1}
}

//Set uses given function f as a mock of Client.State method
func (m *mClientMockState) Set(f func() (r []byte, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StateFunc = f
	return m.mock
}

//State implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) State() (r []byte, r1 error) {
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
		r1 = result.r1

		return
	}

	if m.StateMock.mainExpectation != nil {

		result := m.StateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.State")
		}

		r = result.r
		r1 = result.r1

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

type mClientMockUpdateObject struct {
	mock              *ClientMock
	mainExpectation   *ClientMockUpdateObjectExpectation
	expectationSeries []*ClientMockUpdateObjectExpectation
}

type ClientMockUpdateObjectExpectation struct {
	input  *ClientMockUpdateObjectInput
	result *ClientMockUpdateObjectResult
}

type ClientMockUpdateObjectInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 ObjectDescriptor
	p4 []byte
}

type ClientMockUpdateObjectResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Client.UpdateObject is expected from 1 to Infinity times
func (m *mClientMockUpdateObject) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) *mClientMockUpdateObject {
	m.mock.UpdateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockUpdateObjectExpectation{}
	}
	m.mainExpectation.input = &ClientMockUpdateObjectInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of Client.UpdateObject
func (m *mClientMockUpdateObject) Return(r ObjectDescriptor, r1 error) *ClientMock {
	m.mock.UpdateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockUpdateObjectExpectation{}
	}
	m.mainExpectation.result = &ClientMockUpdateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.UpdateObject is expected once
func (m *mClientMockUpdateObject) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) *ClientMockUpdateObjectExpectation {
	m.mock.UpdateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockUpdateObjectExpectation{}
	expectation.input = &ClientMockUpdateObjectInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockUpdateObjectExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ClientMockUpdateObjectResult{r, r1}
}

//Set uses given function f as a mock of Client.UpdateObject method
func (m *mClientMockUpdateObject) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) (r ObjectDescriptor, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateObjectFunc = f
	return m.mock
}

//UpdateObject implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) UpdateObject(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.UpdateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.UpdateObjectCounter, 1)

	if len(m.UpdateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.UpdateObject. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.UpdateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockUpdateObjectInput{p, p1, p2, p3, p4}, "Client.UpdateObject got unexpected parameters")

		result := m.UpdateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.UpdateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdateObjectMock.mainExpectation != nil {

		input := m.UpdateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockUpdateObjectInput{p, p1, p2, p3, p4}, "Client.UpdateObject got unexpected parameters")
		}

		result := m.UpdateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.UpdateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.UpdateObject. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.UpdateObjectFunc(p, p1, p2, p3, p4)
}

//UpdateObjectMinimockCounter returns a count of ClientMock.UpdateObjectFunc invocations
func (m *ClientMock) UpdateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateObjectCounter)
}

//UpdateObjectMinimockPreCounter returns the value of ClientMock.UpdateObject invocations
func (m *ClientMock) UpdateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateObjectPreCounter)
}

//UpdateObjectFinished returns true if mock invocations count is ok
func (m *ClientMock) UpdateObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpdateObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpdateObjectCounter) == uint64(len(m.UpdateObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpdateObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpdateObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpdateObjectFunc != nil {
		return atomic.LoadUint64(&m.UpdateObjectCounter) > 0
	}

	return true
}

type mClientMockUpdatePrototype struct {
	mock              *ClientMock
	mainExpectation   *ClientMockUpdatePrototypeExpectation
	expectationSeries []*ClientMockUpdatePrototypeExpectation
}

type ClientMockUpdatePrototypeExpectation struct {
	input  *ClientMockUpdatePrototypeInput
	result *ClientMockUpdatePrototypeResult
}

type ClientMockUpdatePrototypeInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 ObjectDescriptor
	p4 []byte
	p5 *insolar.Reference
}

type ClientMockUpdatePrototypeResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Client.UpdatePrototype is expected from 1 to Infinity times
func (m *mClientMockUpdatePrototype) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte, p5 *insolar.Reference) *mClientMockUpdatePrototype {
	m.mock.UpdatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockUpdatePrototypeExpectation{}
	}
	m.mainExpectation.input = &ClientMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}
	return m
}

//Return specifies results of invocation of Client.UpdatePrototype
func (m *mClientMockUpdatePrototype) Return(r ObjectDescriptor, r1 error) *ClientMock {
	m.mock.UpdatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClientMockUpdatePrototypeExpectation{}
	}
	m.mainExpectation.result = &ClientMockUpdatePrototypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Client.UpdatePrototype is expected once
func (m *mClientMockUpdatePrototype) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte, p5 *insolar.Reference) *ClientMockUpdatePrototypeExpectation {
	m.mock.UpdatePrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &ClientMockUpdatePrototypeExpectation{}
	expectation.input = &ClientMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClientMockUpdatePrototypeExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ClientMockUpdatePrototypeResult{r, r1}
}

//Set uses given function f as a mock of Client.UpdatePrototype method
func (m *mClientMockUpdatePrototype) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte, p5 *insolar.Reference) (r ObjectDescriptor, r1 error)) *ClientMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdatePrototypeFunc = f
	return m.mock
}

//UpdatePrototype implements github.com/insolar/insolar/logicrunner/artifacts.Client interface
func (m *ClientMock) UpdatePrototype(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte, p5 *insolar.Reference) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.UpdatePrototypePreCounter, 1)
	defer atomic.AddUint64(&m.UpdatePrototypeCounter, 1)

	if len(m.UpdatePrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdatePrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClientMock.UpdatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
			return
		}

		input := m.UpdatePrototypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClientMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}, "Client.UpdatePrototype got unexpected parameters")

		result := m.UpdatePrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.UpdatePrototype")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdatePrototypeMock.mainExpectation != nil {

		input := m.UpdatePrototypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClientMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}, "Client.UpdatePrototype got unexpected parameters")
		}

		result := m.UpdatePrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClientMock.UpdatePrototype")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdatePrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to ClientMock.UpdatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
		return
	}

	return m.UpdatePrototypeFunc(p, p1, p2, p3, p4, p5)
}

//UpdatePrototypeMinimockCounter returns a count of ClientMock.UpdatePrototypeFunc invocations
func (m *ClientMock) UpdatePrototypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdatePrototypeCounter)
}

//UpdatePrototypeMinimockPreCounter returns the value of ClientMock.UpdatePrototype invocations
func (m *ClientMock) UpdatePrototypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdatePrototypePreCounter)
}

//UpdatePrototypeFinished returns true if mock invocations count is ok
func (m *ClientMock) UpdatePrototypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpdatePrototypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpdatePrototypeCounter) == uint64(len(m.UpdatePrototypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpdatePrototypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpdatePrototypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpdatePrototypeFunc != nil {
		return atomic.LoadUint64(&m.UpdatePrototypeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ClientMock) ValidateCallCounters() {

	if !m.ActivateObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.ActivateObject")
	}

	if !m.ActivatePrototypeFinished() {
		m.t.Fatal("Expected call to ClientMock.ActivatePrototype")
	}

	if !m.DeactivateObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.DeactivateObject")
	}

	if !m.DeclareTypeFinished() {
		m.t.Fatal("Expected call to ClientMock.DeclareType")
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

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterRequest")
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

	if !m.UpdateObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.UpdateObject")
	}

	if !m.UpdatePrototypeFinished() {
		m.t.Fatal("Expected call to ClientMock.UpdatePrototype")
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

	if !m.ActivateObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.ActivateObject")
	}

	if !m.ActivatePrototypeFinished() {
		m.t.Fatal("Expected call to ClientMock.ActivatePrototype")
	}

	if !m.DeactivateObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.DeactivateObject")
	}

	if !m.DeclareTypeFinished() {
		m.t.Fatal("Expected call to ClientMock.DeclareType")
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

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ClientMock.RegisterRequest")
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

	if !m.UpdateObjectFinished() {
		m.t.Fatal("Expected call to ClientMock.UpdateObject")
	}

	if !m.UpdatePrototypeFinished() {
		m.t.Fatal("Expected call to ClientMock.UpdatePrototype")
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
		ok = ok && m.ActivateObjectFinished()
		ok = ok && m.ActivatePrototypeFinished()
		ok = ok && m.DeactivateObjectFinished()
		ok = ok && m.DeclareTypeFinished()
		ok = ok && m.DeployCodeFinished()
		ok = ok && m.GetChildrenFinished()
		ok = ok && m.GetCodeFinished()
		ok = ok && m.GetDelegateFinished()
		ok = ok && m.GetObjectFinished()
		ok = ok && m.GetPendingRequestFinished()
		ok = ok && m.HasPendingRequestsFinished()
		ok = ok && m.RegisterRequestFinished()
		ok = ok && m.RegisterResultFinished()
		ok = ok && m.RegisterValidationFinished()
		ok = ok && m.StateFinished()
		ok = ok && m.UpdateObjectFinished()
		ok = ok && m.UpdatePrototypeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ActivateObjectFinished() {
				m.t.Error("Expected call to ClientMock.ActivateObject")
			}

			if !m.ActivatePrototypeFinished() {
				m.t.Error("Expected call to ClientMock.ActivatePrototype")
			}

			if !m.DeactivateObjectFinished() {
				m.t.Error("Expected call to ClientMock.DeactivateObject")
			}

			if !m.DeclareTypeFinished() {
				m.t.Error("Expected call to ClientMock.DeclareType")
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

			if !m.RegisterRequestFinished() {
				m.t.Error("Expected call to ClientMock.RegisterRequest")
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

			if !m.UpdateObjectFinished() {
				m.t.Error("Expected call to ClientMock.UpdateObject")
			}

			if !m.UpdatePrototypeFinished() {
				m.t.Error("Expected call to ClientMock.UpdatePrototype")
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

	if !m.ActivateObjectFinished() {
		return false
	}

	if !m.ActivatePrototypeFinished() {
		return false
	}

	if !m.DeactivateObjectFinished() {
		return false
	}

	if !m.DeclareTypeFinished() {
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

	if !m.RegisterRequestFinished() {
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

	if !m.UpdateObjectFinished() {
		return false
	}

	if !m.UpdatePrototypeFinished() {
		return false
	}

	return true
}
