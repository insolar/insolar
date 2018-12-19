package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ArtifactManager" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ArtifactManagerMock implements github.com/insolar/insolar/core.ArtifactManager
type ArtifactManagerMock struct {
	t minimock.Tester

	ActivateObjectFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error)
	ActivateObjectCounter    uint64
	ActivateObjectPreCounter uint64
	ActivateObjectMock       mArtifactManagerMockActivateObject

	ActivatePrototypeFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) (r core.ObjectDescriptor, r1 error)
	ActivatePrototypeCounter    uint64
	ActivatePrototypePreCounter uint64
	ActivatePrototypeMock       mArtifactManagerMockActivatePrototype

	DeactivateObjectFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) (r *core.RecordID, r1 error)
	DeactivateObjectCounter    uint64
	DeactivateObjectPreCounter uint64
	DeactivateObjectMock       mArtifactManagerMockDeactivateObject

	DeclareTypeFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error)
	DeclareTypeCounter    uint64
	DeclareTypePreCounter uint64
	DeclareTypeMock       mArtifactManagerMockDeclareType

	DeployCodeFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) (r *core.RecordID, r1 error)
	DeployCodeCounter    uint64
	DeployCodePreCounter uint64
	DeployCodeMock       mArtifactManagerMockDeployCode

	GenesisRefFunc       func() (r *core.RecordRef)
	GenesisRefCounter    uint64
	GenesisRefPreCounter uint64
	GenesisRefMock       mArtifactManagerMockGenesisRef

	GetChildrenFunc       func(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) (r core.RefIterator, r1 error)
	GetChildrenCounter    uint64
	GetChildrenPreCounter uint64
	GetChildrenMock       mArtifactManagerMockGetChildren

	GetCodeFunc       func(p context.Context, p1 core.RecordRef) (r core.CodeDescriptor, r1 error)
	GetCodeCounter    uint64
	GetCodePreCounter uint64
	GetCodeMock       mArtifactManagerMockGetCode

	GetDelegateFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef) (r *core.RecordRef, r1 error)
	GetDelegateCounter    uint64
	GetDelegatePreCounter uint64
	GetDelegateMock       mArtifactManagerMockGetDelegate

	GetObjectFunc       func(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) (r core.ObjectDescriptor, r1 error)
	GetObjectCounter    uint64
	GetObjectPreCounter uint64
	GetObjectMock       mArtifactManagerMockGetObject

	HasPendingRequestsFunc       func(p context.Context, p1 core.RecordRef) (r bool, r1 error)
	HasPendingRequestsCounter    uint64
	HasPendingRequestsPreCounter uint64
	HasPendingRequestsMock       mArtifactManagerMockHasPendingRequests

	RegisterRequestFunc       func(p context.Context, p1 core.RecordRef, p2 core.Parcel) (r *core.RecordID, r1 error)
	RegisterRequestCounter    uint64
	RegisterRequestPreCounter uint64
	RegisterRequestMock       mArtifactManagerMockRegisterRequest

	RegisterResultFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error)
	RegisterResultCounter    uint64
	RegisterResultPreCounter uint64
	RegisterResultMock       mArtifactManagerMockRegisterResult

	RegisterValidationFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) (r error)
	RegisterValidationCounter    uint64
	RegisterValidationPreCounter uint64
	RegisterValidationMock       mArtifactManagerMockRegisterValidation

	StateFunc       func() (r []byte, r1 error)
	StateCounter    uint64
	StatePreCounter uint64
	StateMock       mArtifactManagerMockState

	UpdateObjectFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) (r core.ObjectDescriptor, r1 error)
	UpdateObjectCounter    uint64
	UpdateObjectPreCounter uint64
	UpdateObjectMock       mArtifactManagerMockUpdateObject

	UpdatePrototypeFunc       func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) (r core.ObjectDescriptor, r1 error)
	UpdatePrototypeCounter    uint64
	UpdatePrototypePreCounter uint64
	UpdatePrototypeMock       mArtifactManagerMockUpdatePrototype
}

//NewArtifactManagerMock returns a mock for github.com/insolar/insolar/core.ArtifactManager
func NewArtifactManagerMock(t minimock.Tester) *ArtifactManagerMock {
	m := &ArtifactManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ActivateObjectMock = mArtifactManagerMockActivateObject{mock: m}
	m.ActivatePrototypeMock = mArtifactManagerMockActivatePrototype{mock: m}
	m.DeactivateObjectMock = mArtifactManagerMockDeactivateObject{mock: m}
	m.DeclareTypeMock = mArtifactManagerMockDeclareType{mock: m}
	m.DeployCodeMock = mArtifactManagerMockDeployCode{mock: m}
	m.GenesisRefMock = mArtifactManagerMockGenesisRef{mock: m}
	m.GetChildrenMock = mArtifactManagerMockGetChildren{mock: m}
	m.GetCodeMock = mArtifactManagerMockGetCode{mock: m}
	m.GetDelegateMock = mArtifactManagerMockGetDelegate{mock: m}
	m.GetObjectMock = mArtifactManagerMockGetObject{mock: m}
	m.HasPendingRequestsMock = mArtifactManagerMockHasPendingRequests{mock: m}
	m.RegisterRequestMock = mArtifactManagerMockRegisterRequest{mock: m}
	m.RegisterResultMock = mArtifactManagerMockRegisterResult{mock: m}
	m.RegisterValidationMock = mArtifactManagerMockRegisterValidation{mock: m}
	m.StateMock = mArtifactManagerMockState{mock: m}
	m.UpdateObjectMock = mArtifactManagerMockUpdateObject{mock: m}
	m.UpdatePrototypeMock = mArtifactManagerMockUpdatePrototype{mock: m}

	return m
}

type mArtifactManagerMockActivateObject struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockActivateObjectExpectation
	expectationSeries []*ArtifactManagerMockActivateObjectExpectation
}

type ArtifactManagerMockActivateObjectExpectation struct {
	input  *ArtifactManagerMockActivateObjectInput
	result *ArtifactManagerMockActivateObjectResult
}

type ArtifactManagerMockActivateObjectInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.RecordRef
	p4 core.RecordRef
	p5 bool
	p6 []byte
}

type ArtifactManagerMockActivateObjectResult struct {
	r  core.ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of ArtifactManager.ActivateObject is expected from 1 to Infinity times
func (m *mArtifactManagerMockActivateObject) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) *mArtifactManagerMockActivateObject {
	m.mock.ActivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockActivateObjectExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}
	return m
}

//Return specifies results of invocation of ArtifactManager.ActivateObject
func (m *mArtifactManagerMockActivateObject) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.ActivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockActivateObjectExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockActivateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.ActivateObject is expected once
func (m *mArtifactManagerMockActivateObject) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) *ArtifactManagerMockActivateObjectExpectation {
	m.mock.ActivateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockActivateObjectExpectation{}
	expectation.input = &ArtifactManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockActivateObjectExpectation) Return(r core.ObjectDescriptor, r1 error) {
	e.result = &ArtifactManagerMockActivateObjectResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.ActivateObject method
func (m *mArtifactManagerMockActivateObject) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ActivateObjectFunc = f
	return m.mock
}

//ActivateObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) ActivateObject(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.ActivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.ActivateObjectCounter, 1)

	if len(m.ActivateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ActivateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.ActivateObject. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
			return
		}

		input := m.ActivateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}, "ArtifactManager.ActivateObject got unexpected parameters")

		result := m.ActivateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.ActivateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivateObjectMock.mainExpectation != nil {

		input := m.ActivateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}, "ArtifactManager.ActivateObject got unexpected parameters")
		}

		result := m.ActivateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.ActivateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.ActivateObject. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
		return
	}

	return m.ActivateObjectFunc(p, p1, p2, p3, p4, p5, p6)
}

//ActivateObjectMinimockCounter returns a count of ArtifactManagerMock.ActivateObjectFunc invocations
func (m *ArtifactManagerMock) ActivateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ActivateObjectCounter)
}

//ActivateObjectMinimockPreCounter returns the value of ArtifactManagerMock.ActivateObject invocations
func (m *ArtifactManagerMock) ActivateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ActivateObjectPreCounter)
}

//ActivateObjectFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) ActivateObjectFinished() bool {
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

type mArtifactManagerMockActivatePrototype struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockActivatePrototypeExpectation
	expectationSeries []*ArtifactManagerMockActivatePrototypeExpectation
}

type ArtifactManagerMockActivatePrototypeExpectation struct {
	input  *ArtifactManagerMockActivatePrototypeInput
	result *ArtifactManagerMockActivatePrototypeResult
}

type ArtifactManagerMockActivatePrototypeInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.RecordRef
	p4 core.RecordRef
	p5 []byte
}

type ArtifactManagerMockActivatePrototypeResult struct {
	r  core.ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of ArtifactManager.ActivatePrototype is expected from 1 to Infinity times
func (m *mArtifactManagerMockActivatePrototype) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) *mArtifactManagerMockActivatePrototype {
	m.mock.ActivatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockActivatePrototypeExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}
	return m
}

//Return specifies results of invocation of ArtifactManager.ActivatePrototype
func (m *mArtifactManagerMockActivatePrototype) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.ActivatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockActivatePrototypeExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockActivatePrototypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.ActivatePrototype is expected once
func (m *mArtifactManagerMockActivatePrototype) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) *ArtifactManagerMockActivatePrototypeExpectation {
	m.mock.ActivatePrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockActivatePrototypeExpectation{}
	expectation.input = &ArtifactManagerMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockActivatePrototypeExpectation) Return(r core.ObjectDescriptor, r1 error) {
	e.result = &ArtifactManagerMockActivatePrototypeResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.ActivatePrototype method
func (m *mArtifactManagerMockActivatePrototype) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ActivatePrototypeFunc = f
	return m.mock
}

//ActivatePrototype implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) ActivatePrototype(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) (r core.ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.ActivatePrototypePreCounter, 1)
	defer atomic.AddUint64(&m.ActivatePrototypeCounter, 1)

	if len(m.ActivatePrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ActivatePrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.ActivatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
			return
		}

		input := m.ActivatePrototypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}, "ArtifactManager.ActivatePrototype got unexpected parameters")

		result := m.ActivatePrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.ActivatePrototype")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivatePrototypeMock.mainExpectation != nil {

		input := m.ActivatePrototypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockActivatePrototypeInput{p, p1, p2, p3, p4, p5}, "ArtifactManager.ActivatePrototype got unexpected parameters")
		}

		result := m.ActivatePrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.ActivatePrototype")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivatePrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.ActivatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
		return
	}

	return m.ActivatePrototypeFunc(p, p1, p2, p3, p4, p5)
}

//ActivatePrototypeMinimockCounter returns a count of ArtifactManagerMock.ActivatePrototypeFunc invocations
func (m *ArtifactManagerMock) ActivatePrototypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ActivatePrototypeCounter)
}

//ActivatePrototypeMinimockPreCounter returns the value of ArtifactManagerMock.ActivatePrototype invocations
func (m *ArtifactManagerMock) ActivatePrototypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ActivatePrototypePreCounter)
}

//ActivatePrototypeFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) ActivatePrototypeFinished() bool {
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

type mArtifactManagerMockDeactivateObject struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockDeactivateObjectExpectation
	expectationSeries []*ArtifactManagerMockDeactivateObjectExpectation
}

type ArtifactManagerMockDeactivateObjectExpectation struct {
	input  *ArtifactManagerMockDeactivateObjectInput
	result *ArtifactManagerMockDeactivateObjectResult
}

type ArtifactManagerMockDeactivateObjectInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.ObjectDescriptor
}

type ArtifactManagerMockDeactivateObjectResult struct {
	r  *core.RecordID
	r1 error
}

//Expect specifies that invocation of ArtifactManager.DeactivateObject is expected from 1 to Infinity times
func (m *mArtifactManagerMockDeactivateObject) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) *mArtifactManagerMockDeactivateObject {
	m.mock.DeactivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockDeactivateObjectExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockDeactivateObjectInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ArtifactManager.DeactivateObject
func (m *mArtifactManagerMockDeactivateObject) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.DeactivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockDeactivateObjectExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockDeactivateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.DeactivateObject is expected once
func (m *mArtifactManagerMockDeactivateObject) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) *ArtifactManagerMockDeactivateObjectExpectation {
	m.mock.DeactivateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockDeactivateObjectExpectation{}
	expectation.input = &ArtifactManagerMockDeactivateObjectInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockDeactivateObjectExpectation) Return(r *core.RecordID, r1 error) {
	e.result = &ArtifactManagerMockDeactivateObjectResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.DeactivateObject method
func (m *mArtifactManagerMockDeactivateObject) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeactivateObjectFunc = f
	return m.mock
}

//DeactivateObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) DeactivateObject(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) (r *core.RecordID, r1 error) {
	counter := atomic.AddUint64(&m.DeactivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.DeactivateObjectCounter, 1)

	if len(m.DeactivateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeactivateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.DeactivateObject. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.DeactivateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockDeactivateObjectInput{p, p1, p2, p3}, "ArtifactManager.DeactivateObject got unexpected parameters")

		result := m.DeactivateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.DeactivateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeactivateObjectMock.mainExpectation != nil {

		input := m.DeactivateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockDeactivateObjectInput{p, p1, p2, p3}, "ArtifactManager.DeactivateObject got unexpected parameters")
		}

		result := m.DeactivateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.DeactivateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeactivateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.DeactivateObject. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.DeactivateObjectFunc(p, p1, p2, p3)
}

//DeactivateObjectMinimockCounter returns a count of ArtifactManagerMock.DeactivateObjectFunc invocations
func (m *ArtifactManagerMock) DeactivateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateObjectCounter)
}

//DeactivateObjectMinimockPreCounter returns the value of ArtifactManagerMock.DeactivateObject invocations
func (m *ArtifactManagerMock) DeactivateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateObjectPreCounter)
}

//DeactivateObjectFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) DeactivateObjectFinished() bool {
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

type mArtifactManagerMockDeclareType struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockDeclareTypeExpectation
	expectationSeries []*ArtifactManagerMockDeclareTypeExpectation
}

type ArtifactManagerMockDeclareTypeExpectation struct {
	input  *ArtifactManagerMockDeclareTypeInput
	result *ArtifactManagerMockDeclareTypeResult
}

type ArtifactManagerMockDeclareTypeInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 []byte
}

type ArtifactManagerMockDeclareTypeResult struct {
	r  *core.RecordID
	r1 error
}

//Expect specifies that invocation of ArtifactManager.DeclareType is expected from 1 to Infinity times
func (m *mArtifactManagerMockDeclareType) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) *mArtifactManagerMockDeclareType {
	m.mock.DeclareTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockDeclareTypeExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockDeclareTypeInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ArtifactManager.DeclareType
func (m *mArtifactManagerMockDeclareType) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.DeclareTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockDeclareTypeExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockDeclareTypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.DeclareType is expected once
func (m *mArtifactManagerMockDeclareType) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) *ArtifactManagerMockDeclareTypeExpectation {
	m.mock.DeclareTypeFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockDeclareTypeExpectation{}
	expectation.input = &ArtifactManagerMockDeclareTypeInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockDeclareTypeExpectation) Return(r *core.RecordID, r1 error) {
	e.result = &ArtifactManagerMockDeclareTypeResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.DeclareType method
func (m *mArtifactManagerMockDeclareType) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeclareTypeFunc = f
	return m.mock
}

//DeclareType implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) DeclareType(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error) {
	counter := atomic.AddUint64(&m.DeclareTypePreCounter, 1)
	defer atomic.AddUint64(&m.DeclareTypeCounter, 1)

	if len(m.DeclareTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeclareTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.DeclareType. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.DeclareTypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockDeclareTypeInput{p, p1, p2, p3}, "ArtifactManager.DeclareType got unexpected parameters")

		result := m.DeclareTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.DeclareType")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeclareTypeMock.mainExpectation != nil {

		input := m.DeclareTypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockDeclareTypeInput{p, p1, p2, p3}, "ArtifactManager.DeclareType got unexpected parameters")
		}

		result := m.DeclareTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.DeclareType")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeclareTypeFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.DeclareType. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.DeclareTypeFunc(p, p1, p2, p3)
}

//DeclareTypeMinimockCounter returns a count of ArtifactManagerMock.DeclareTypeFunc invocations
func (m *ArtifactManagerMock) DeclareTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeclareTypeCounter)
}

//DeclareTypeMinimockPreCounter returns the value of ArtifactManagerMock.DeclareType invocations
func (m *ArtifactManagerMock) DeclareTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeclareTypePreCounter)
}

//DeclareTypeFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) DeclareTypeFinished() bool {
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

type mArtifactManagerMockDeployCode struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockDeployCodeExpectation
	expectationSeries []*ArtifactManagerMockDeployCodeExpectation
}

type ArtifactManagerMockDeployCodeExpectation struct {
	input  *ArtifactManagerMockDeployCodeInput
	result *ArtifactManagerMockDeployCodeResult
}

type ArtifactManagerMockDeployCodeInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 []byte
	p4 core.MachineType
}

type ArtifactManagerMockDeployCodeResult struct {
	r  *core.RecordID
	r1 error
}

//Expect specifies that invocation of ArtifactManager.DeployCode is expected from 1 to Infinity times
func (m *mArtifactManagerMockDeployCode) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) *mArtifactManagerMockDeployCode {
	m.mock.DeployCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockDeployCodeExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockDeployCodeInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of ArtifactManager.DeployCode
func (m *mArtifactManagerMockDeployCode) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.DeployCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockDeployCodeExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockDeployCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.DeployCode is expected once
func (m *mArtifactManagerMockDeployCode) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) *ArtifactManagerMockDeployCodeExpectation {
	m.mock.DeployCodeFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockDeployCodeExpectation{}
	expectation.input = &ArtifactManagerMockDeployCodeInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockDeployCodeExpectation) Return(r *core.RecordID, r1 error) {
	e.result = &ArtifactManagerMockDeployCodeResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.DeployCode method
func (m *mArtifactManagerMockDeployCode) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeployCodeFunc = f
	return m.mock
}

//DeployCode implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) DeployCode(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) (r *core.RecordID, r1 error) {
	counter := atomic.AddUint64(&m.DeployCodePreCounter, 1)
	defer atomic.AddUint64(&m.DeployCodeCounter, 1)

	if len(m.DeployCodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeployCodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.DeployCode. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.DeployCodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockDeployCodeInput{p, p1, p2, p3, p4}, "ArtifactManager.DeployCode got unexpected parameters")

		result := m.DeployCodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.DeployCode")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeployCodeMock.mainExpectation != nil {

		input := m.DeployCodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockDeployCodeInput{p, p1, p2, p3, p4}, "ArtifactManager.DeployCode got unexpected parameters")
		}

		result := m.DeployCodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.DeployCode")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.DeployCodeFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.DeployCode. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.DeployCodeFunc(p, p1, p2, p3, p4)
}

//DeployCodeMinimockCounter returns a count of ArtifactManagerMock.DeployCodeFunc invocations
func (m *ArtifactManagerMock) DeployCodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeployCodeCounter)
}

//DeployCodeMinimockPreCounter returns the value of ArtifactManagerMock.DeployCode invocations
func (m *ArtifactManagerMock) DeployCodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeployCodePreCounter)
}

//DeployCodeFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) DeployCodeFinished() bool {
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

type mArtifactManagerMockGenesisRef struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockGenesisRefExpectation
	expectationSeries []*ArtifactManagerMockGenesisRefExpectation
}

type ArtifactManagerMockGenesisRefExpectation struct {
	result *ArtifactManagerMockGenesisRefResult
}

type ArtifactManagerMockGenesisRefResult struct {
	r *core.RecordRef
}

//Expect specifies that invocation of ArtifactManager.GenesisRef is expected from 1 to Infinity times
func (m *mArtifactManagerMockGenesisRef) Expect() *mArtifactManagerMockGenesisRef {
	m.mock.GenesisRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGenesisRefExpectation{}
	}

	return m
}

//Return specifies results of invocation of ArtifactManager.GenesisRef
func (m *mArtifactManagerMockGenesisRef) Return(r *core.RecordRef) *ArtifactManagerMock {
	m.mock.GenesisRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGenesisRefExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockGenesisRefResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.GenesisRef is expected once
func (m *mArtifactManagerMockGenesisRef) ExpectOnce() *ArtifactManagerMockGenesisRefExpectation {
	m.mock.GenesisRefFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockGenesisRefExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockGenesisRefExpectation) Return(r *core.RecordRef) {
	e.result = &ArtifactManagerMockGenesisRefResult{r}
}

//Set uses given function f as a mock of ArtifactManager.GenesisRef method
func (m *mArtifactManagerMockGenesisRef) Set(f func() (r *core.RecordRef)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GenesisRefFunc = f
	return m.mock
}

//GenesisRef implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GenesisRef() (r *core.RecordRef) {
	counter := atomic.AddUint64(&m.GenesisRefPreCounter, 1)
	defer atomic.AddUint64(&m.GenesisRefCounter, 1)

	if len(m.GenesisRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GenesisRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.GenesisRef.")
			return
		}

		result := m.GenesisRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GenesisRef")
			return
		}

		r = result.r

		return
	}

	if m.GenesisRefMock.mainExpectation != nil {

		result := m.GenesisRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GenesisRef")
		}

		r = result.r

		return
	}

	if m.GenesisRefFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.GenesisRef.")
		return
	}

	return m.GenesisRefFunc()
}

//GenesisRefMinimockCounter returns a count of ArtifactManagerMock.GenesisRefFunc invocations
func (m *ArtifactManagerMock) GenesisRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GenesisRefCounter)
}

//GenesisRefMinimockPreCounter returns the value of ArtifactManagerMock.GenesisRef invocations
func (m *ArtifactManagerMock) GenesisRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GenesisRefPreCounter)
}

//GenesisRefFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) GenesisRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GenesisRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GenesisRefCounter) == uint64(len(m.GenesisRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GenesisRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GenesisRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GenesisRefFunc != nil {
		return atomic.LoadUint64(&m.GenesisRefCounter) > 0
	}

	return true
}

type mArtifactManagerMockGetChildren struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockGetChildrenExpectation
	expectationSeries []*ArtifactManagerMockGetChildrenExpectation
}

type ArtifactManagerMockGetChildrenExpectation struct {
	input  *ArtifactManagerMockGetChildrenInput
	result *ArtifactManagerMockGetChildrenResult
}

type ArtifactManagerMockGetChildrenInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 *core.PulseNumber
}

type ArtifactManagerMockGetChildrenResult struct {
	r  core.RefIterator
	r1 error
}

//Expect specifies that invocation of ArtifactManager.GetChildren is expected from 1 to Infinity times
func (m *mArtifactManagerMockGetChildren) Expect(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) *mArtifactManagerMockGetChildren {
	m.mock.GetChildrenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetChildrenExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockGetChildrenInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ArtifactManager.GetChildren
func (m *mArtifactManagerMockGetChildren) Return(r core.RefIterator, r1 error) *ArtifactManagerMock {
	m.mock.GetChildrenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetChildrenExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockGetChildrenResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.GetChildren is expected once
func (m *mArtifactManagerMockGetChildren) ExpectOnce(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) *ArtifactManagerMockGetChildrenExpectation {
	m.mock.GetChildrenFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockGetChildrenExpectation{}
	expectation.input = &ArtifactManagerMockGetChildrenInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockGetChildrenExpectation) Return(r core.RefIterator, r1 error) {
	e.result = &ArtifactManagerMockGetChildrenResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.GetChildren method
func (m *mArtifactManagerMockGetChildren) Set(f func(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) (r core.RefIterator, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetChildrenFunc = f
	return m.mock
}

//GetChildren implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetChildren(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) (r core.RefIterator, r1 error) {
	counter := atomic.AddUint64(&m.GetChildrenPreCounter, 1)
	defer atomic.AddUint64(&m.GetChildrenCounter, 1)

	if len(m.GetChildrenMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetChildrenMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetChildren. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetChildrenMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockGetChildrenInput{p, p1, p2}, "ArtifactManager.GetChildren got unexpected parameters")

		result := m.GetChildrenMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetChildren")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetChildrenMock.mainExpectation != nil {

		input := m.GetChildrenMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockGetChildrenInput{p, p1, p2}, "ArtifactManager.GetChildren got unexpected parameters")
		}

		result := m.GetChildrenMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetChildren")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetChildrenFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetChildren. %v %v %v", p, p1, p2)
		return
	}

	return m.GetChildrenFunc(p, p1, p2)
}

//GetChildrenMinimockCounter returns a count of ArtifactManagerMock.GetChildrenFunc invocations
func (m *ArtifactManagerMock) GetChildrenMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetChildrenCounter)
}

//GetChildrenMinimockPreCounter returns the value of ArtifactManagerMock.GetChildren invocations
func (m *ArtifactManagerMock) GetChildrenMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetChildrenPreCounter)
}

//GetChildrenFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) GetChildrenFinished() bool {
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

type mArtifactManagerMockGetCode struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockGetCodeExpectation
	expectationSeries []*ArtifactManagerMockGetCodeExpectation
}

type ArtifactManagerMockGetCodeExpectation struct {
	input  *ArtifactManagerMockGetCodeInput
	result *ArtifactManagerMockGetCodeResult
}

type ArtifactManagerMockGetCodeInput struct {
	p  context.Context
	p1 core.RecordRef
}

type ArtifactManagerMockGetCodeResult struct {
	r  core.CodeDescriptor
	r1 error
}

//Expect specifies that invocation of ArtifactManager.GetCode is expected from 1 to Infinity times
func (m *mArtifactManagerMockGetCode) Expect(p context.Context, p1 core.RecordRef) *mArtifactManagerMockGetCode {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetCodeExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockGetCodeInput{p, p1}
	return m
}

//Return specifies results of invocation of ArtifactManager.GetCode
func (m *mArtifactManagerMockGetCode) Return(r core.CodeDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetCodeExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockGetCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.GetCode is expected once
func (m *mArtifactManagerMockGetCode) ExpectOnce(p context.Context, p1 core.RecordRef) *ArtifactManagerMockGetCodeExpectation {
	m.mock.GetCodeFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockGetCodeExpectation{}
	expectation.input = &ArtifactManagerMockGetCodeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockGetCodeExpectation) Return(r core.CodeDescriptor, r1 error) {
	e.result = &ArtifactManagerMockGetCodeResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.GetCode method
func (m *mArtifactManagerMockGetCode) Set(f func(p context.Context, p1 core.RecordRef) (r core.CodeDescriptor, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCodeFunc = f
	return m.mock
}

//GetCode implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetCode(p context.Context, p1 core.RecordRef) (r core.CodeDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.GetCodePreCounter, 1)
	defer atomic.AddUint64(&m.GetCodeCounter, 1)

	if len(m.GetCodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetCode. %v %v", p, p1)
			return
		}

		input := m.GetCodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockGetCodeInput{p, p1}, "ArtifactManager.GetCode got unexpected parameters")

		result := m.GetCodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetCode")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeMock.mainExpectation != nil {

		input := m.GetCodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockGetCodeInput{p, p1}, "ArtifactManager.GetCode got unexpected parameters")
		}

		result := m.GetCodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetCode")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetCode. %v %v", p, p1)
		return
	}

	return m.GetCodeFunc(p, p1)
}

//GetCodeMinimockCounter returns a count of ArtifactManagerMock.GetCodeFunc invocations
func (m *ArtifactManagerMock) GetCodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodeCounter)
}

//GetCodeMinimockPreCounter returns the value of ArtifactManagerMock.GetCode invocations
func (m *ArtifactManagerMock) GetCodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodePreCounter)
}

//GetCodeFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) GetCodeFinished() bool {
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

type mArtifactManagerMockGetDelegate struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockGetDelegateExpectation
	expectationSeries []*ArtifactManagerMockGetDelegateExpectation
}

type ArtifactManagerMockGetDelegateExpectation struct {
	input  *ArtifactManagerMockGetDelegateInput
	result *ArtifactManagerMockGetDelegateResult
}

type ArtifactManagerMockGetDelegateInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
}

type ArtifactManagerMockGetDelegateResult struct {
	r  *core.RecordRef
	r1 error
}

//Expect specifies that invocation of ArtifactManager.GetDelegate is expected from 1 to Infinity times
func (m *mArtifactManagerMockGetDelegate) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef) *mArtifactManagerMockGetDelegate {
	m.mock.GetDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetDelegateExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockGetDelegateInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ArtifactManager.GetDelegate
func (m *mArtifactManagerMockGetDelegate) Return(r *core.RecordRef, r1 error) *ArtifactManagerMock {
	m.mock.GetDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetDelegateExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockGetDelegateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.GetDelegate is expected once
func (m *mArtifactManagerMockGetDelegate) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef) *ArtifactManagerMockGetDelegateExpectation {
	m.mock.GetDelegateFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockGetDelegateExpectation{}
	expectation.input = &ArtifactManagerMockGetDelegateInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockGetDelegateExpectation) Return(r *core.RecordRef, r1 error) {
	e.result = &ArtifactManagerMockGetDelegateResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.GetDelegate method
func (m *mArtifactManagerMockGetDelegate) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef) (r *core.RecordRef, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDelegateFunc = f
	return m.mock
}

//GetDelegate implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetDelegate(p context.Context, p1 core.RecordRef, p2 core.RecordRef) (r *core.RecordRef, r1 error) {
	counter := atomic.AddUint64(&m.GetDelegatePreCounter, 1)
	defer atomic.AddUint64(&m.GetDelegateCounter, 1)

	if len(m.GetDelegateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDelegateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetDelegate. %v %v %v", p, p1, p2)
			return
		}

		input := m.GetDelegateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockGetDelegateInput{p, p1, p2}, "ArtifactManager.GetDelegate got unexpected parameters")

		result := m.GetDelegateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetDelegate")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDelegateMock.mainExpectation != nil {

		input := m.GetDelegateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockGetDelegateInput{p, p1, p2}, "ArtifactManager.GetDelegate got unexpected parameters")
		}

		result := m.GetDelegateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetDelegate")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetDelegateFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetDelegate. %v %v %v", p, p1, p2)
		return
	}

	return m.GetDelegateFunc(p, p1, p2)
}

//GetDelegateMinimockCounter returns a count of ArtifactManagerMock.GetDelegateFunc invocations
func (m *ArtifactManagerMock) GetDelegateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDelegateCounter)
}

//GetDelegateMinimockPreCounter returns the value of ArtifactManagerMock.GetDelegate invocations
func (m *ArtifactManagerMock) GetDelegateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDelegatePreCounter)
}

//GetDelegateFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) GetDelegateFinished() bool {
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

type mArtifactManagerMockGetObject struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockGetObjectExpectation
	expectationSeries []*ArtifactManagerMockGetObjectExpectation
}

type ArtifactManagerMockGetObjectExpectation struct {
	input  *ArtifactManagerMockGetObjectInput
	result *ArtifactManagerMockGetObjectResult
}

type ArtifactManagerMockGetObjectInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 *core.RecordID
	p3 bool
}

type ArtifactManagerMockGetObjectResult struct {
	r  core.ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of ArtifactManager.GetObject is expected from 1 to Infinity times
func (m *mArtifactManagerMockGetObject) Expect(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) *mArtifactManagerMockGetObject {
	m.mock.GetObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetObjectExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockGetObjectInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ArtifactManager.GetObject
func (m *mArtifactManagerMockGetObject) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.GetObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockGetObjectExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockGetObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.GetObject is expected once
func (m *mArtifactManagerMockGetObject) ExpectOnce(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) *ArtifactManagerMockGetObjectExpectation {
	m.mock.GetObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockGetObjectExpectation{}
	expectation.input = &ArtifactManagerMockGetObjectInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockGetObjectExpectation) Return(r core.ObjectDescriptor, r1 error) {
	e.result = &ArtifactManagerMockGetObjectResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.GetObject method
func (m *mArtifactManagerMockGetObject) Set(f func(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjectFunc = f
	return m.mock
}

//GetObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetObject(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) (r core.ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.GetObjectPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectCounter, 1)

	if len(m.GetObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetObject. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.GetObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockGetObjectInput{p, p1, p2, p3}, "ArtifactManager.GetObject got unexpected parameters")

		result := m.GetObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetObjectMock.mainExpectation != nil {

		input := m.GetObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockGetObjectInput{p, p1, p2, p3}, "ArtifactManager.GetObject got unexpected parameters")
		}

		result := m.GetObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.GetObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.GetObject. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.GetObjectFunc(p, p1, p2, p3)
}

//GetObjectMinimockCounter returns a count of ArtifactManagerMock.GetObjectFunc invocations
func (m *ArtifactManagerMock) GetObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectCounter)
}

//GetObjectMinimockPreCounter returns the value of ArtifactManagerMock.GetObject invocations
func (m *ArtifactManagerMock) GetObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjectPreCounter)
}

//GetObjectFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) GetObjectFinished() bool {
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

type mArtifactManagerMockHasPendingRequests struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockHasPendingRequestsExpectation
	expectationSeries []*ArtifactManagerMockHasPendingRequestsExpectation
}

type ArtifactManagerMockHasPendingRequestsExpectation struct {
	input  *ArtifactManagerMockHasPendingRequestsInput
	result *ArtifactManagerMockHasPendingRequestsResult
}

type ArtifactManagerMockHasPendingRequestsInput struct {
	p  context.Context
	p1 core.RecordRef
}

type ArtifactManagerMockHasPendingRequestsResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of ArtifactManager.HasPendingRequests is expected from 1 to Infinity times
func (m *mArtifactManagerMockHasPendingRequests) Expect(p context.Context, p1 core.RecordRef) *mArtifactManagerMockHasPendingRequests {
	m.mock.HasPendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockHasPendingRequestsExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockHasPendingRequestsInput{p, p1}
	return m
}

//Return specifies results of invocation of ArtifactManager.HasPendingRequests
func (m *mArtifactManagerMockHasPendingRequests) Return(r bool, r1 error) *ArtifactManagerMock {
	m.mock.HasPendingRequestsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockHasPendingRequestsExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockHasPendingRequestsResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.HasPendingRequests is expected once
func (m *mArtifactManagerMockHasPendingRequests) ExpectOnce(p context.Context, p1 core.RecordRef) *ArtifactManagerMockHasPendingRequestsExpectation {
	m.mock.HasPendingRequestsFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockHasPendingRequestsExpectation{}
	expectation.input = &ArtifactManagerMockHasPendingRequestsInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockHasPendingRequestsExpectation) Return(r bool, r1 error) {
	e.result = &ArtifactManagerMockHasPendingRequestsResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.HasPendingRequests method
func (m *mArtifactManagerMockHasPendingRequests) Set(f func(p context.Context, p1 core.RecordRef) (r bool, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HasPendingRequestsFunc = f
	return m.mock
}

//HasPendingRequests implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) HasPendingRequests(p context.Context, p1 core.RecordRef) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.HasPendingRequestsPreCounter, 1)
	defer atomic.AddUint64(&m.HasPendingRequestsCounter, 1)

	if len(m.HasPendingRequestsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HasPendingRequestsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.HasPendingRequests. %v %v", p, p1)
			return
		}

		input := m.HasPendingRequestsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockHasPendingRequestsInput{p, p1}, "ArtifactManager.HasPendingRequests got unexpected parameters")

		result := m.HasPendingRequestsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.HasPendingRequests")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HasPendingRequestsMock.mainExpectation != nil {

		input := m.HasPendingRequestsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockHasPendingRequestsInput{p, p1}, "ArtifactManager.HasPendingRequests got unexpected parameters")
		}

		result := m.HasPendingRequestsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.HasPendingRequests")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.HasPendingRequestsFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.HasPendingRequests. %v %v", p, p1)
		return
	}

	return m.HasPendingRequestsFunc(p, p1)
}

//HasPendingRequestsMinimockCounter returns a count of ArtifactManagerMock.HasPendingRequestsFunc invocations
func (m *ArtifactManagerMock) HasPendingRequestsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HasPendingRequestsCounter)
}

//HasPendingRequestsMinimockPreCounter returns the value of ArtifactManagerMock.HasPendingRequests invocations
func (m *ArtifactManagerMock) HasPendingRequestsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HasPendingRequestsPreCounter)
}

//HasPendingRequestsFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) HasPendingRequestsFinished() bool {
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

type mArtifactManagerMockRegisterRequest struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockRegisterRequestExpectation
	expectationSeries []*ArtifactManagerMockRegisterRequestExpectation
}

type ArtifactManagerMockRegisterRequestExpectation struct {
	input  *ArtifactManagerMockRegisterRequestInput
	result *ArtifactManagerMockRegisterRequestResult
}

type ArtifactManagerMockRegisterRequestInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.Parcel
}

type ArtifactManagerMockRegisterRequestResult struct {
	r  *core.RecordID
	r1 error
}

//Expect specifies that invocation of ArtifactManager.RegisterRequest is expected from 1 to Infinity times
func (m *mArtifactManagerMockRegisterRequest) Expect(p context.Context, p1 core.RecordRef, p2 core.Parcel) *mArtifactManagerMockRegisterRequest {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockRegisterRequestExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockRegisterRequestInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of ArtifactManager.RegisterRequest
func (m *mArtifactManagerMockRegisterRequest) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockRegisterRequestExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockRegisterRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.RegisterRequest is expected once
func (m *mArtifactManagerMockRegisterRequest) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.Parcel) *ArtifactManagerMockRegisterRequestExpectation {
	m.mock.RegisterRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockRegisterRequestExpectation{}
	expectation.input = &ArtifactManagerMockRegisterRequestInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockRegisterRequestExpectation) Return(r *core.RecordID, r1 error) {
	e.result = &ArtifactManagerMockRegisterRequestResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.RegisterRequest method
func (m *mArtifactManagerMockRegisterRequest) Set(f func(p context.Context, p1 core.RecordRef, p2 core.Parcel) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterRequestFunc = f
	return m.mock
}

//RegisterRequest implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) RegisterRequest(p context.Context, p1 core.RecordRef, p2 core.Parcel) (r *core.RecordID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestCounter, 1)

	if len(m.RegisterRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.RegisterRequest. %v %v %v", p, p1, p2)
			return
		}

		input := m.RegisterRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockRegisterRequestInput{p, p1, p2}, "ArtifactManager.RegisterRequest got unexpected parameters")

		result := m.RegisterRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestMock.mainExpectation != nil {

		input := m.RegisterRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockRegisterRequestInput{p, p1, p2}, "ArtifactManager.RegisterRequest got unexpected parameters")
		}

		result := m.RegisterRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.RegisterRequest. %v %v %v", p, p1, p2)
		return
	}

	return m.RegisterRequestFunc(p, p1, p2)
}

//RegisterRequestMinimockCounter returns a count of ArtifactManagerMock.RegisterRequestFunc invocations
func (m *ArtifactManagerMock) RegisterRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestCounter)
}

//RegisterRequestMinimockPreCounter returns the value of ArtifactManagerMock.RegisterRequest invocations
func (m *ArtifactManagerMock) RegisterRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestPreCounter)
}

//RegisterRequestFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) RegisterRequestFinished() bool {
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

type mArtifactManagerMockRegisterResult struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockRegisterResultExpectation
	expectationSeries []*ArtifactManagerMockRegisterResultExpectation
}

type ArtifactManagerMockRegisterResultExpectation struct {
	input  *ArtifactManagerMockRegisterResultInput
	result *ArtifactManagerMockRegisterResultResult
}

type ArtifactManagerMockRegisterResultInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 []byte
}

type ArtifactManagerMockRegisterResultResult struct {
	r  *core.RecordID
	r1 error
}

//Expect specifies that invocation of ArtifactManager.RegisterResult is expected from 1 to Infinity times
func (m *mArtifactManagerMockRegisterResult) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) *mArtifactManagerMockRegisterResult {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockRegisterResultExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockRegisterResultInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ArtifactManager.RegisterResult
func (m *mArtifactManagerMockRegisterResult) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockRegisterResultExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockRegisterResultResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.RegisterResult is expected once
func (m *mArtifactManagerMockRegisterResult) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) *ArtifactManagerMockRegisterResultExpectation {
	m.mock.RegisterResultFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockRegisterResultExpectation{}
	expectation.input = &ArtifactManagerMockRegisterResultInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockRegisterResultExpectation) Return(r *core.RecordID, r1 error) {
	e.result = &ArtifactManagerMockRegisterResultResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.RegisterResult method
func (m *mArtifactManagerMockRegisterResult) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterResultFunc = f
	return m.mock
}

//RegisterResult implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) RegisterResult(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterResultPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterResultCounter, 1)

	if len(m.RegisterResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.RegisterResult. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RegisterResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockRegisterResultInput{p, p1, p2, p3}, "ArtifactManager.RegisterResult got unexpected parameters")

		result := m.RegisterResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterResult")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterResultMock.mainExpectation != nil {

		input := m.RegisterResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockRegisterResultInput{p, p1, p2, p3}, "ArtifactManager.RegisterResult got unexpected parameters")
		}

		result := m.RegisterResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterResult")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterResultFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.RegisterResult. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RegisterResultFunc(p, p1, p2, p3)
}

//RegisterResultMinimockCounter returns a count of ArtifactManagerMock.RegisterResultFunc invocations
func (m *ArtifactManagerMock) RegisterResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultCounter)
}

//RegisterResultMinimockPreCounter returns the value of ArtifactManagerMock.RegisterResult invocations
func (m *ArtifactManagerMock) RegisterResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultPreCounter)
}

//RegisterResultFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) RegisterResultFinished() bool {
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

type mArtifactManagerMockRegisterValidation struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockRegisterValidationExpectation
	expectationSeries []*ArtifactManagerMockRegisterValidationExpectation
}

type ArtifactManagerMockRegisterValidationExpectation struct {
	input  *ArtifactManagerMockRegisterValidationInput
	result *ArtifactManagerMockRegisterValidationResult
}

type ArtifactManagerMockRegisterValidationInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordID
	p3 bool
	p4 []core.Message
}

type ArtifactManagerMockRegisterValidationResult struct {
	r error
}

//Expect specifies that invocation of ArtifactManager.RegisterValidation is expected from 1 to Infinity times
func (m *mArtifactManagerMockRegisterValidation) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) *mArtifactManagerMockRegisterValidation {
	m.mock.RegisterValidationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockRegisterValidationExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockRegisterValidationInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of ArtifactManager.RegisterValidation
func (m *mArtifactManagerMockRegisterValidation) Return(r error) *ArtifactManagerMock {
	m.mock.RegisterValidationFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockRegisterValidationExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockRegisterValidationResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.RegisterValidation is expected once
func (m *mArtifactManagerMockRegisterValidation) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) *ArtifactManagerMockRegisterValidationExpectation {
	m.mock.RegisterValidationFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockRegisterValidationExpectation{}
	expectation.input = &ArtifactManagerMockRegisterValidationInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockRegisterValidationExpectation) Return(r error) {
	e.result = &ArtifactManagerMockRegisterValidationResult{r}
}

//Set uses given function f as a mock of ArtifactManager.RegisterValidation method
func (m *mArtifactManagerMockRegisterValidation) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) (r error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterValidationFunc = f
	return m.mock
}

//RegisterValidation implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) RegisterValidation(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) (r error) {
	counter := atomic.AddUint64(&m.RegisterValidationPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterValidationCounter, 1)

	if len(m.RegisterValidationMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterValidationMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.RegisterValidation. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.RegisterValidationMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockRegisterValidationInput{p, p1, p2, p3, p4}, "ArtifactManager.RegisterValidation got unexpected parameters")

		result := m.RegisterValidationMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterValidation")
			return
		}

		r = result.r

		return
	}

	if m.RegisterValidationMock.mainExpectation != nil {

		input := m.RegisterValidationMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockRegisterValidationInput{p, p1, p2, p3, p4}, "ArtifactManager.RegisterValidation got unexpected parameters")
		}

		result := m.RegisterValidationMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterValidation")
		}

		r = result.r

		return
	}

	if m.RegisterValidationFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.RegisterValidation. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.RegisterValidationFunc(p, p1, p2, p3, p4)
}

//RegisterValidationMinimockCounter returns a count of ArtifactManagerMock.RegisterValidationFunc invocations
func (m *ArtifactManagerMock) RegisterValidationMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterValidationCounter)
}

//RegisterValidationMinimockPreCounter returns the value of ArtifactManagerMock.RegisterValidation invocations
func (m *ArtifactManagerMock) RegisterValidationMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterValidationPreCounter)
}

//RegisterValidationFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) RegisterValidationFinished() bool {
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

type mArtifactManagerMockState struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockStateExpectation
	expectationSeries []*ArtifactManagerMockStateExpectation
}

type ArtifactManagerMockStateExpectation struct {
	result *ArtifactManagerMockStateResult
}

type ArtifactManagerMockStateResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of ArtifactManager.State is expected from 1 to Infinity times
func (m *mArtifactManagerMockState) Expect() *mArtifactManagerMockState {
	m.mock.StateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockStateExpectation{}
	}

	return m
}

//Return specifies results of invocation of ArtifactManager.State
func (m *mArtifactManagerMockState) Return(r []byte, r1 error) *ArtifactManagerMock {
	m.mock.StateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockStateExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockStateResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.State is expected once
func (m *mArtifactManagerMockState) ExpectOnce() *ArtifactManagerMockStateExpectation {
	m.mock.StateFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockStateExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockStateExpectation) Return(r []byte, r1 error) {
	e.result = &ArtifactManagerMockStateResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.State method
func (m *mArtifactManagerMockState) Set(f func() (r []byte, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StateFunc = f
	return m.mock
}

//State implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) State() (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.StatePreCounter, 1)
	defer atomic.AddUint64(&m.StateCounter, 1)

	if len(m.StateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.State.")
			return
		}

		result := m.StateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.State")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.StateMock.mainExpectation != nil {

		result := m.StateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.State")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.StateFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.State.")
		return
	}

	return m.StateFunc()
}

//StateMinimockCounter returns a count of ArtifactManagerMock.StateFunc invocations
func (m *ArtifactManagerMock) StateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StateCounter)
}

//StateMinimockPreCounter returns the value of ArtifactManagerMock.State invocations
func (m *ArtifactManagerMock) StateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StatePreCounter)
}

//StateFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) StateFinished() bool {
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

type mArtifactManagerMockUpdateObject struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockUpdateObjectExpectation
	expectationSeries []*ArtifactManagerMockUpdateObjectExpectation
}

type ArtifactManagerMockUpdateObjectExpectation struct {
	input  *ArtifactManagerMockUpdateObjectInput
	result *ArtifactManagerMockUpdateObjectResult
}

type ArtifactManagerMockUpdateObjectInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.ObjectDescriptor
	p4 []byte
}

type ArtifactManagerMockUpdateObjectResult struct {
	r  core.ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of ArtifactManager.UpdateObject is expected from 1 to Infinity times
func (m *mArtifactManagerMockUpdateObject) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) *mArtifactManagerMockUpdateObject {
	m.mock.UpdateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockUpdateObjectExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockUpdateObjectInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of ArtifactManager.UpdateObject
func (m *mArtifactManagerMockUpdateObject) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.UpdateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockUpdateObjectExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockUpdateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.UpdateObject is expected once
func (m *mArtifactManagerMockUpdateObject) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) *ArtifactManagerMockUpdateObjectExpectation {
	m.mock.UpdateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockUpdateObjectExpectation{}
	expectation.input = &ArtifactManagerMockUpdateObjectInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockUpdateObjectExpectation) Return(r core.ObjectDescriptor, r1 error) {
	e.result = &ArtifactManagerMockUpdateObjectResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.UpdateObject method
func (m *mArtifactManagerMockUpdateObject) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateObjectFunc = f
	return m.mock
}

//UpdateObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) UpdateObject(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) (r core.ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.UpdateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.UpdateObjectCounter, 1)

	if len(m.UpdateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.UpdateObject. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.UpdateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockUpdateObjectInput{p, p1, p2, p3, p4}, "ArtifactManager.UpdateObject got unexpected parameters")

		result := m.UpdateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.UpdateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdateObjectMock.mainExpectation != nil {

		input := m.UpdateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockUpdateObjectInput{p, p1, p2, p3, p4}, "ArtifactManager.UpdateObject got unexpected parameters")
		}

		result := m.UpdateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.UpdateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.UpdateObject. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.UpdateObjectFunc(p, p1, p2, p3, p4)
}

//UpdateObjectMinimockCounter returns a count of ArtifactManagerMock.UpdateObjectFunc invocations
func (m *ArtifactManagerMock) UpdateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateObjectCounter)
}

//UpdateObjectMinimockPreCounter returns the value of ArtifactManagerMock.UpdateObject invocations
func (m *ArtifactManagerMock) UpdateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateObjectPreCounter)
}

//UpdateObjectFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) UpdateObjectFinished() bool {
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

type mArtifactManagerMockUpdatePrototype struct {
	mock              *ArtifactManagerMock
	mainExpectation   *ArtifactManagerMockUpdatePrototypeExpectation
	expectationSeries []*ArtifactManagerMockUpdatePrototypeExpectation
}

type ArtifactManagerMockUpdatePrototypeExpectation struct {
	input  *ArtifactManagerMockUpdatePrototypeInput
	result *ArtifactManagerMockUpdatePrototypeResult
}

type ArtifactManagerMockUpdatePrototypeInput struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.ObjectDescriptor
	p4 []byte
	p5 *core.RecordRef
}

type ArtifactManagerMockUpdatePrototypeResult struct {
	r  core.ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of ArtifactManager.UpdatePrototype is expected from 1 to Infinity times
func (m *mArtifactManagerMockUpdatePrototype) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) *mArtifactManagerMockUpdatePrototype {
	m.mock.UpdatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockUpdatePrototypeExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}
	return m
}

//Return specifies results of invocation of ArtifactManager.UpdatePrototype
func (m *mArtifactManagerMockUpdatePrototype) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.UpdatePrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMockUpdatePrototypeExpectation{}
	}
	m.mainExpectation.result = &ArtifactManagerMockUpdatePrototypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManager.UpdatePrototype is expected once
func (m *mArtifactManagerMockUpdatePrototype) ExpectOnce(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) *ArtifactManagerMockUpdatePrototypeExpectation {
	m.mock.UpdatePrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMockUpdatePrototypeExpectation{}
	expectation.input = &ArtifactManagerMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ArtifactManagerMockUpdatePrototypeExpectation) Return(r core.ObjectDescriptor, r1 error) {
	e.result = &ArtifactManagerMockUpdatePrototypeResult{r, r1}
}

//Set uses given function f as a mock of ArtifactManager.UpdatePrototype method
func (m *mArtifactManagerMockUpdatePrototype) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdatePrototypeFunc = f
	return m.mock
}

//UpdatePrototype implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) UpdatePrototype(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) (r core.ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.UpdatePrototypePreCounter, 1)
	defer atomic.AddUint64(&m.UpdatePrototypeCounter, 1)

	if len(m.UpdatePrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdatePrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMock.UpdatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
			return
		}

		input := m.UpdatePrototypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}, "ArtifactManager.UpdatePrototype got unexpected parameters")

		result := m.UpdatePrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.UpdatePrototype")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdatePrototypeMock.mainExpectation != nil {

		input := m.UpdatePrototypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMockUpdatePrototypeInput{p, p1, p2, p3, p4, p5}, "ArtifactManager.UpdatePrototype got unexpected parameters")
		}

		result := m.UpdatePrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ArtifactManagerMock.UpdatePrototype")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdatePrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMock.UpdatePrototype. %v %v %v %v %v %v", p, p1, p2, p3, p4, p5)
		return
	}

	return m.UpdatePrototypeFunc(p, p1, p2, p3, p4, p5)
}

//UpdatePrototypeMinimockCounter returns a count of ArtifactManagerMock.UpdatePrototypeFunc invocations
func (m *ArtifactManagerMock) UpdatePrototypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdatePrototypeCounter)
}

//UpdatePrototypeMinimockPreCounter returns the value of ArtifactManagerMock.UpdatePrototype invocations
func (m *ArtifactManagerMock) UpdatePrototypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdatePrototypePreCounter)
}

//UpdatePrototypeFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMock) UpdatePrototypeFinished() bool {
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
func (m *ArtifactManagerMock) ValidateCallCounters() {

	if !m.ActivateObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivateObject")
	}

	if !m.ActivatePrototypeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivatePrototype")
	}

	if !m.DeactivateObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeactivateObject")
	}

	if !m.DeclareTypeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeclareType")
	}

	if !m.DeployCodeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeployCode")
	}

	if !m.GenesisRefFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GenesisRef")
	}

	if !m.GetChildrenFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetChildren")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetCode")
	}

	if !m.GetDelegateFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetDelegate")
	}

	if !m.GetObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetObject")
	}

	if !m.HasPendingRequestsFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.HasPendingRequests")
	}

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterRequest")
	}

	if !m.RegisterResultFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterResult")
	}

	if !m.RegisterValidationFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterValidation")
	}

	if !m.StateFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.State")
	}

	if !m.UpdateObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.UpdateObject")
	}

	if !m.UpdatePrototypeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.UpdatePrototype")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ArtifactManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ArtifactManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ArtifactManagerMock) MinimockFinish() {

	if !m.ActivateObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivateObject")
	}

	if !m.ActivatePrototypeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivatePrototype")
	}

	if !m.DeactivateObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeactivateObject")
	}

	if !m.DeclareTypeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeclareType")
	}

	if !m.DeployCodeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeployCode")
	}

	if !m.GenesisRefFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GenesisRef")
	}

	if !m.GetChildrenFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetChildren")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetCode")
	}

	if !m.GetDelegateFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetDelegate")
	}

	if !m.GetObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetObject")
	}

	if !m.HasPendingRequestsFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.HasPendingRequests")
	}

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterRequest")
	}

	if !m.RegisterResultFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterResult")
	}

	if !m.RegisterValidationFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterValidation")
	}

	if !m.StateFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.State")
	}

	if !m.UpdateObjectFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.UpdateObject")
	}

	if !m.UpdatePrototypeFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMock.UpdatePrototype")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ArtifactManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ArtifactManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ActivateObjectFinished()
		ok = ok && m.ActivatePrototypeFinished()
		ok = ok && m.DeactivateObjectFinished()
		ok = ok && m.DeclareTypeFinished()
		ok = ok && m.DeployCodeFinished()
		ok = ok && m.GenesisRefFinished()
		ok = ok && m.GetChildrenFinished()
		ok = ok && m.GetCodeFinished()
		ok = ok && m.GetDelegateFinished()
		ok = ok && m.GetObjectFinished()
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
				m.t.Error("Expected call to ArtifactManagerMock.ActivateObject")
			}

			if !m.ActivatePrototypeFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.ActivatePrototype")
			}

			if !m.DeactivateObjectFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.DeactivateObject")
			}

			if !m.DeclareTypeFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.DeclareType")
			}

			if !m.DeployCodeFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.DeployCode")
			}

			if !m.GenesisRefFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.GenesisRef")
			}

			if !m.GetChildrenFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.GetChildren")
			}

			if !m.GetCodeFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.GetCode")
			}

			if !m.GetDelegateFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.GetDelegate")
			}

			if !m.GetObjectFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.GetObject")
			}

			if !m.HasPendingRequestsFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.HasPendingRequests")
			}

			if !m.RegisterRequestFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.RegisterRequest")
			}

			if !m.RegisterResultFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.RegisterResult")
			}

			if !m.RegisterValidationFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.RegisterValidation")
			}

			if !m.StateFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.State")
			}

			if !m.UpdateObjectFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.UpdateObject")
			}

			if !m.UpdatePrototypeFinished() {
				m.t.Error("Expected call to ArtifactManagerMock.UpdatePrototype")
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
func (m *ArtifactManagerMock) AllMocksCalled() bool {

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

	if !m.GenesisRefFinished() {
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
