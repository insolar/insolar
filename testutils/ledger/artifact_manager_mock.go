package ledger

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

	RegisterRequestFunc       func(p context.Context, p1 core.Message) (r *core.RecordID, r1 error)
	RegisterRequestCounter    uint64
	RegisterRequestPreCounter uint64
	RegisterRequestMock       mArtifactManagerMockRegisterRequest

	RegisterResultFunc       func(p context.Context, p1 core.RecordRef, p2 []byte) (r *core.RecordID, r1 error)
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
	m.RegisterRequestMock = mArtifactManagerMockRegisterRequest{mock: m}
	m.RegisterResultMock = mArtifactManagerMockRegisterResult{mock: m}
	m.RegisterValidationMock = mArtifactManagerMockRegisterValidation{mock: m}
	m.StateMock = mArtifactManagerMockState{mock: m}
	m.UpdateObjectMock = mArtifactManagerMockUpdateObject{mock: m}
	m.UpdatePrototypeMock = mArtifactManagerMockUpdatePrototype{mock: m}

	return m
}

type mArtifactManagerMockActivateObject struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockActivateObjectParams
}

//ArtifactManagerMockActivateObjectParams represents input parameters of the ArtifactManager.ActivateObject
type ArtifactManagerMockActivateObjectParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.RecordRef
	p4 core.RecordRef
	p5 bool
	p6 []byte
}

//Expect sets up expected params for the ArtifactManager.ActivateObject
func (m *mArtifactManagerMockActivateObject) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) *mArtifactManagerMockActivateObject {
	m.mockExpectations = &ArtifactManagerMockActivateObjectParams{p, p1, p2, p3, p4, p5, p6}
	return m
}

//Return sets up a mock for ArtifactManager.ActivateObject to return Return's arguments
func (m *mArtifactManagerMockActivateObject) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.ActivateObjectFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (core.ObjectDescriptor, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.ActivateObject method
func (m *mArtifactManagerMockActivateObject) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mock.ActivateObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ActivateObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) ActivateObject(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error) {
	atomic.AddUint64(&m.ActivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.ActivateObjectCounter, 1)

	if m.ActivateObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ActivateObjectMock.mockExpectations, ArtifactManagerMockActivateObjectParams{p, p1, p2, p3, p4, p5, p6},
			"ArtifactManager.ActivateObject got unexpected parameters")

		if m.ActivateObjectFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.ActivateObject")

			return
		}
	}

	if m.ActivateObjectFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.ActivateObject")
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

type mArtifactManagerMockActivatePrototype struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockActivatePrototypeParams
}

//ArtifactManagerMockActivatePrototypeParams represents input parameters of the ArtifactManager.ActivatePrototype
type ArtifactManagerMockActivatePrototypeParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.RecordRef
	p4 core.RecordRef
	p5 []byte
}

//Expect sets up expected params for the ArtifactManager.ActivatePrototype
func (m *mArtifactManagerMockActivatePrototype) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) *mArtifactManagerMockActivatePrototype {
	m.mockExpectations = &ArtifactManagerMockActivatePrototypeParams{p, p1, p2, p3, p4, p5}
	return m
}

//Return sets up a mock for ArtifactManager.ActivatePrototype to return Return's arguments
func (m *mArtifactManagerMockActivatePrototype) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.ActivatePrototypeFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) (core.ObjectDescriptor, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.ActivatePrototype method
func (m *mArtifactManagerMockActivatePrototype) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mock.ActivatePrototypeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ActivatePrototype implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) ActivatePrototype(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 []byte) (r core.ObjectDescriptor, r1 error) {
	atomic.AddUint64(&m.ActivatePrototypePreCounter, 1)
	defer atomic.AddUint64(&m.ActivatePrototypeCounter, 1)

	if m.ActivatePrototypeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ActivatePrototypeMock.mockExpectations, ArtifactManagerMockActivatePrototypeParams{p, p1, p2, p3, p4, p5},
			"ArtifactManager.ActivatePrototype got unexpected parameters")

		if m.ActivatePrototypeFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.ActivatePrototype")

			return
		}
	}

	if m.ActivatePrototypeFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.ActivatePrototype")
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

type mArtifactManagerMockDeactivateObject struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockDeactivateObjectParams
}

//ArtifactManagerMockDeactivateObjectParams represents input parameters of the ArtifactManager.DeactivateObject
type ArtifactManagerMockDeactivateObjectParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.ObjectDescriptor
}

//Expect sets up expected params for the ArtifactManager.DeactivateObject
func (m *mArtifactManagerMockDeactivateObject) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) *mArtifactManagerMockDeactivateObject {
	m.mockExpectations = &ArtifactManagerMockDeactivateObjectParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for ArtifactManager.DeactivateObject to return Return's arguments
func (m *mArtifactManagerMockDeactivateObject) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.DeactivateObjectFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) (*core.RecordID, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.DeactivateObject method
func (m *mArtifactManagerMockDeactivateObject) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mock.DeactivateObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//DeactivateObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) DeactivateObject(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor) (r *core.RecordID, r1 error) {
	atomic.AddUint64(&m.DeactivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.DeactivateObjectCounter, 1)

	if m.DeactivateObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.DeactivateObjectMock.mockExpectations, ArtifactManagerMockDeactivateObjectParams{p, p1, p2, p3},
			"ArtifactManager.DeactivateObject got unexpected parameters")

		if m.DeactivateObjectFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.DeactivateObject")

			return
		}
	}

	if m.DeactivateObjectFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.DeactivateObject")
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

type mArtifactManagerMockDeclareType struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockDeclareTypeParams
}

//ArtifactManagerMockDeclareTypeParams represents input parameters of the ArtifactManager.DeclareType
type ArtifactManagerMockDeclareTypeParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 []byte
}

//Expect sets up expected params for the ArtifactManager.DeclareType
func (m *mArtifactManagerMockDeclareType) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) *mArtifactManagerMockDeclareType {
	m.mockExpectations = &ArtifactManagerMockDeclareTypeParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for ArtifactManager.DeclareType to return Return's arguments
func (m *mArtifactManagerMockDeclareType) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.DeclareTypeFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (*core.RecordID, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.DeclareType method
func (m *mArtifactManagerMockDeclareType) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mock.DeclareTypeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//DeclareType implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) DeclareType(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error) {
	atomic.AddUint64(&m.DeclareTypePreCounter, 1)
	defer atomic.AddUint64(&m.DeclareTypeCounter, 1)

	if m.DeclareTypeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.DeclareTypeMock.mockExpectations, ArtifactManagerMockDeclareTypeParams{p, p1, p2, p3},
			"ArtifactManager.DeclareType got unexpected parameters")

		if m.DeclareTypeFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.DeclareType")

			return
		}
	}

	if m.DeclareTypeFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.DeclareType")
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

type mArtifactManagerMockDeployCode struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockDeployCodeParams
}

//ArtifactManagerMockDeployCodeParams represents input parameters of the ArtifactManager.DeployCode
type ArtifactManagerMockDeployCodeParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 []byte
	p4 core.MachineType
}

//Expect sets up expected params for the ArtifactManager.DeployCode
func (m *mArtifactManagerMockDeployCode) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) *mArtifactManagerMockDeployCode {
	m.mockExpectations = &ArtifactManagerMockDeployCodeParams{p, p1, p2, p3, p4}
	return m
}

//Return sets up a mock for ArtifactManager.DeployCode to return Return's arguments
func (m *mArtifactManagerMockDeployCode) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.DeployCodeFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) (*core.RecordID, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.DeployCode method
func (m *mArtifactManagerMockDeployCode) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mock.DeployCodeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//DeployCode implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) DeployCode(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte, p4 core.MachineType) (r *core.RecordID, r1 error) {
	atomic.AddUint64(&m.DeployCodePreCounter, 1)
	defer atomic.AddUint64(&m.DeployCodeCounter, 1)

	if m.DeployCodeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.DeployCodeMock.mockExpectations, ArtifactManagerMockDeployCodeParams{p, p1, p2, p3, p4},
			"ArtifactManager.DeployCode got unexpected parameters")

		if m.DeployCodeFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.DeployCode")

			return
		}
	}

	if m.DeployCodeFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.DeployCode")
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

type mArtifactManagerMockGenesisRef struct {
	mock *ArtifactManagerMock
}

//Return sets up a mock for ArtifactManager.GenesisRef to return Return's arguments
func (m *mArtifactManagerMockGenesisRef) Return(r *core.RecordRef) *ArtifactManagerMock {
	m.mock.GenesisRefFunc = func() *core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.GenesisRef method
func (m *mArtifactManagerMockGenesisRef) Set(f func() (r *core.RecordRef)) *ArtifactManagerMock {
	m.mock.GenesisRefFunc = f

	return m.mock
}

//GenesisRef implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GenesisRef() (r *core.RecordRef) {
	atomic.AddUint64(&m.GenesisRefPreCounter, 1)
	defer atomic.AddUint64(&m.GenesisRefCounter, 1)

	if m.GenesisRefFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.GenesisRef")
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

type mArtifactManagerMockGetChildren struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockGetChildrenParams
}

//ArtifactManagerMockGetChildrenParams represents input parameters of the ArtifactManager.GetChildren
type ArtifactManagerMockGetChildrenParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 *core.PulseNumber
}

//Expect sets up expected params for the ArtifactManager.GetChildren
func (m *mArtifactManagerMockGetChildren) Expect(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) *mArtifactManagerMockGetChildren {
	m.mockExpectations = &ArtifactManagerMockGetChildrenParams{p, p1, p2}
	return m
}

//Return sets up a mock for ArtifactManager.GetChildren to return Return's arguments
func (m *mArtifactManagerMockGetChildren) Return(r core.RefIterator, r1 error) *ArtifactManagerMock {
	m.mock.GetChildrenFunc = func(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) (core.RefIterator, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.GetChildren method
func (m *mArtifactManagerMockGetChildren) Set(f func(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) (r core.RefIterator, r1 error)) *ArtifactManagerMock {
	m.mock.GetChildrenFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetChildren implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetChildren(p context.Context, p1 core.RecordRef, p2 *core.PulseNumber) (r core.RefIterator, r1 error) {
	atomic.AddUint64(&m.GetChildrenPreCounter, 1)
	defer atomic.AddUint64(&m.GetChildrenCounter, 1)

	if m.GetChildrenMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetChildrenMock.mockExpectations, ArtifactManagerMockGetChildrenParams{p, p1, p2},
			"ArtifactManager.GetChildren got unexpected parameters")

		if m.GetChildrenFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.GetChildren")

			return
		}
	}

	if m.GetChildrenFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.GetChildren")
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

type mArtifactManagerMockGetCode struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockGetCodeParams
}

//ArtifactManagerMockGetCodeParams represents input parameters of the ArtifactManager.GetCode
type ArtifactManagerMockGetCodeParams struct {
	p  context.Context
	p1 core.RecordRef
}

//Expect sets up expected params for the ArtifactManager.GetCode
func (m *mArtifactManagerMockGetCode) Expect(p context.Context, p1 core.RecordRef) *mArtifactManagerMockGetCode {
	m.mockExpectations = &ArtifactManagerMockGetCodeParams{p, p1}
	return m
}

//Return sets up a mock for ArtifactManager.GetCode to return Return's arguments
func (m *mArtifactManagerMockGetCode) Return(r core.CodeDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.GetCodeFunc = func(p context.Context, p1 core.RecordRef) (core.CodeDescriptor, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.GetCode method
func (m *mArtifactManagerMockGetCode) Set(f func(p context.Context, p1 core.RecordRef) (r core.CodeDescriptor, r1 error)) *ArtifactManagerMock {
	m.mock.GetCodeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetCode implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetCode(p context.Context, p1 core.RecordRef) (r core.CodeDescriptor, r1 error) {
	atomic.AddUint64(&m.GetCodePreCounter, 1)
	defer atomic.AddUint64(&m.GetCodeCounter, 1)

	if m.GetCodeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetCodeMock.mockExpectations, ArtifactManagerMockGetCodeParams{p, p1},
			"ArtifactManager.GetCode got unexpected parameters")

		if m.GetCodeFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.GetCode")

			return
		}
	}

	if m.GetCodeFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.GetCode")
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

type mArtifactManagerMockGetDelegate struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockGetDelegateParams
}

//ArtifactManagerMockGetDelegateParams represents input parameters of the ArtifactManager.GetDelegate
type ArtifactManagerMockGetDelegateParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
}

//Expect sets up expected params for the ArtifactManager.GetDelegate
func (m *mArtifactManagerMockGetDelegate) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef) *mArtifactManagerMockGetDelegate {
	m.mockExpectations = &ArtifactManagerMockGetDelegateParams{p, p1, p2}
	return m
}

//Return sets up a mock for ArtifactManager.GetDelegate to return Return's arguments
func (m *mArtifactManagerMockGetDelegate) Return(r *core.RecordRef, r1 error) *ArtifactManagerMock {
	m.mock.GetDelegateFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.GetDelegate method
func (m *mArtifactManagerMockGetDelegate) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef) (r *core.RecordRef, r1 error)) *ArtifactManagerMock {
	m.mock.GetDelegateFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetDelegate implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetDelegate(p context.Context, p1 core.RecordRef, p2 core.RecordRef) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.GetDelegatePreCounter, 1)
	defer atomic.AddUint64(&m.GetDelegateCounter, 1)

	if m.GetDelegateMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetDelegateMock.mockExpectations, ArtifactManagerMockGetDelegateParams{p, p1, p2},
			"ArtifactManager.GetDelegate got unexpected parameters")

		if m.GetDelegateFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.GetDelegate")

			return
		}
	}

	if m.GetDelegateFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.GetDelegate")
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

type mArtifactManagerMockGetObject struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockGetObjectParams
}

//ArtifactManagerMockGetObjectParams represents input parameters of the ArtifactManager.GetObject
type ArtifactManagerMockGetObjectParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 *core.RecordID
	p3 bool
}

//Expect sets up expected params for the ArtifactManager.GetObject
func (m *mArtifactManagerMockGetObject) Expect(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) *mArtifactManagerMockGetObject {
	m.mockExpectations = &ArtifactManagerMockGetObjectParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for ArtifactManager.GetObject to return Return's arguments
func (m *mArtifactManagerMockGetObject) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.GetObjectFunc = func(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) (core.ObjectDescriptor, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.GetObject method
func (m *mArtifactManagerMockGetObject) Set(f func(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mock.GetObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) GetObject(p context.Context, p1 core.RecordRef, p2 *core.RecordID, p3 bool) (r core.ObjectDescriptor, r1 error) {
	atomic.AddUint64(&m.GetObjectPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjectCounter, 1)

	if m.GetObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetObjectMock.mockExpectations, ArtifactManagerMockGetObjectParams{p, p1, p2, p3},
			"ArtifactManager.GetObject got unexpected parameters")

		if m.GetObjectFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.GetObject")

			return
		}
	}

	if m.GetObjectFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.GetObject")
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

type mArtifactManagerMockRegisterRequest struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockRegisterRequestParams
}

//ArtifactManagerMockRegisterRequestParams represents input parameters of the ArtifactManager.RegisterRequest
type ArtifactManagerMockRegisterRequestParams struct {
	p  context.Context
	p1 core.Message
}

//Expect sets up expected params for the ArtifactManager.RegisterRequest
func (m *mArtifactManagerMockRegisterRequest) Expect(p context.Context, p1 core.Message) *mArtifactManagerMockRegisterRequest {
	m.mockExpectations = &ArtifactManagerMockRegisterRequestParams{p, p1}
	return m
}

//Return sets up a mock for ArtifactManager.RegisterRequest to return Return's arguments
func (m *mArtifactManagerMockRegisterRequest) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.RegisterRequestFunc = func(p context.Context, p1 core.Message) (*core.RecordID, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.RegisterRequest method
func (m *mArtifactManagerMockRegisterRequest) Set(f func(p context.Context, p1 core.Message) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mock.RegisterRequestFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RegisterRequest implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) RegisterRequest(p context.Context, p1 core.Message) (r *core.RecordID, r1 error) {
	atomic.AddUint64(&m.RegisterRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestCounter, 1)

	if m.RegisterRequestMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RegisterRequestMock.mockExpectations, ArtifactManagerMockRegisterRequestParams{p, p1},
			"ArtifactManager.RegisterRequest got unexpected parameters")

		if m.RegisterRequestFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterRequest")

			return
		}
	}

	if m.RegisterRequestFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.RegisterRequest")
		return
	}

	return m.RegisterRequestFunc(p, p1)
}

//RegisterRequestMinimockCounter returns a count of ArtifactManagerMock.RegisterRequestFunc invocations
func (m *ArtifactManagerMock) RegisterRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestCounter)
}

//RegisterRequestMinimockPreCounter returns the value of ArtifactManagerMock.RegisterRequest invocations
func (m *ArtifactManagerMock) RegisterRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestPreCounter)
}

type mArtifactManagerMockRegisterResult struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockRegisterResultParams
}

//ArtifactManagerMockRegisterResultParams represents input parameters of the ArtifactManager.RegisterResult
type ArtifactManagerMockRegisterResultParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 []byte
}

//Expect sets up expected params for the ArtifactManager.RegisterResult
func (m *mArtifactManagerMockRegisterResult) Expect(p context.Context, p1 core.RecordRef, p2 []byte) *mArtifactManagerMockRegisterResult {
	m.mockExpectations = &ArtifactManagerMockRegisterResultParams{p, p1, p2}
	return m
}

//Return sets up a mock for ArtifactManager.RegisterResult to return Return's arguments
func (m *mArtifactManagerMockRegisterResult) Return(r *core.RecordID, r1 error) *ArtifactManagerMock {
	m.mock.RegisterResultFunc = func(p context.Context, p1 core.RecordRef, p2 []byte) (*core.RecordID, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.RegisterResult method
func (m *mArtifactManagerMockRegisterResult) Set(f func(p context.Context, p1 core.RecordRef, p2 []byte) (r *core.RecordID, r1 error)) *ArtifactManagerMock {
	m.mock.RegisterResultFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RegisterResult implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) RegisterResult(p context.Context, p1 core.RecordRef, p2 []byte) (r *core.RecordID, r1 error) {
	atomic.AddUint64(&m.RegisterResultPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterResultCounter, 1)

	if m.RegisterResultMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RegisterResultMock.mockExpectations, ArtifactManagerMockRegisterResultParams{p, p1, p2},
			"ArtifactManager.RegisterResult got unexpected parameters")

		if m.RegisterResultFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterResult")

			return
		}
	}

	if m.RegisterResultFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.RegisterResult")
		return
	}

	return m.RegisterResultFunc(p, p1, p2)
}

//RegisterResultMinimockCounter returns a count of ArtifactManagerMock.RegisterResultFunc invocations
func (m *ArtifactManagerMock) RegisterResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultCounter)
}

//RegisterResultMinimockPreCounter returns the value of ArtifactManagerMock.RegisterResult invocations
func (m *ArtifactManagerMock) RegisterResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultPreCounter)
}

type mArtifactManagerMockRegisterValidation struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockRegisterValidationParams
}

//ArtifactManagerMockRegisterValidationParams represents input parameters of the ArtifactManager.RegisterValidation
type ArtifactManagerMockRegisterValidationParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordID
	p3 bool
	p4 []core.Message
}

//Expect sets up expected params for the ArtifactManager.RegisterValidation
func (m *mArtifactManagerMockRegisterValidation) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) *mArtifactManagerMockRegisterValidation {
	m.mockExpectations = &ArtifactManagerMockRegisterValidationParams{p, p1, p2, p3, p4}
	return m
}

//Return sets up a mock for ArtifactManager.RegisterValidation to return Return's arguments
func (m *mArtifactManagerMockRegisterValidation) Return(r error) *ArtifactManagerMock {
	m.mock.RegisterValidationFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.RegisterValidation method
func (m *mArtifactManagerMockRegisterValidation) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) (r error)) *ArtifactManagerMock {
	m.mock.RegisterValidationFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RegisterValidation implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) RegisterValidation(p context.Context, p1 core.RecordRef, p2 core.RecordID, p3 bool, p4 []core.Message) (r error) {
	atomic.AddUint64(&m.RegisterValidationPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterValidationCounter, 1)

	if m.RegisterValidationMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RegisterValidationMock.mockExpectations, ArtifactManagerMockRegisterValidationParams{p, p1, p2, p3, p4},
			"ArtifactManager.RegisterValidation got unexpected parameters")

		if m.RegisterValidationFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.RegisterValidation")

			return
		}
	}

	if m.RegisterValidationFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.RegisterValidation")
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

type mArtifactManagerMockState struct {
	mock *ArtifactManagerMock
}

//Return sets up a mock for ArtifactManager.State to return Return's arguments
func (m *mArtifactManagerMockState) Return(r []byte, r1 error) *ArtifactManagerMock {
	m.mock.StateFunc = func() ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.State method
func (m *mArtifactManagerMockState) Set(f func() (r []byte, r1 error)) *ArtifactManagerMock {
	m.mock.StateFunc = f

	return m.mock
}

//State implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) State() (r []byte, r1 error) {
	atomic.AddUint64(&m.StatePreCounter, 1)
	defer atomic.AddUint64(&m.StateCounter, 1)

	if m.StateFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.State")
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

type mArtifactManagerMockUpdateObject struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockUpdateObjectParams
}

//ArtifactManagerMockUpdateObjectParams represents input parameters of the ArtifactManager.UpdateObject
type ArtifactManagerMockUpdateObjectParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.ObjectDescriptor
	p4 []byte
}

//Expect sets up expected params for the ArtifactManager.UpdateObject
func (m *mArtifactManagerMockUpdateObject) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) *mArtifactManagerMockUpdateObject {
	m.mockExpectations = &ArtifactManagerMockUpdateObjectParams{p, p1, p2, p3, p4}
	return m
}

//Return sets up a mock for ArtifactManager.UpdateObject to return Return's arguments
func (m *mArtifactManagerMockUpdateObject) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.UpdateObjectFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) (core.ObjectDescriptor, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.UpdateObject method
func (m *mArtifactManagerMockUpdateObject) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mock.UpdateObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//UpdateObject implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) UpdateObject(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte) (r core.ObjectDescriptor, r1 error) {
	atomic.AddUint64(&m.UpdateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.UpdateObjectCounter, 1)

	if m.UpdateObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.UpdateObjectMock.mockExpectations, ArtifactManagerMockUpdateObjectParams{p, p1, p2, p3, p4},
			"ArtifactManager.UpdateObject got unexpected parameters")

		if m.UpdateObjectFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.UpdateObject")

			return
		}
	}

	if m.UpdateObjectFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.UpdateObject")
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

type mArtifactManagerMockUpdatePrototype struct {
	mock             *ArtifactManagerMock
	mockExpectations *ArtifactManagerMockUpdatePrototypeParams
}

//ArtifactManagerMockUpdatePrototypeParams represents input parameters of the ArtifactManager.UpdatePrototype
type ArtifactManagerMockUpdatePrototypeParams struct {
	p  context.Context
	p1 core.RecordRef
	p2 core.RecordRef
	p3 core.ObjectDescriptor
	p4 []byte
	p5 *core.RecordRef
}

//Expect sets up expected params for the ArtifactManager.UpdatePrototype
func (m *mArtifactManagerMockUpdatePrototype) Expect(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) *mArtifactManagerMockUpdatePrototype {
	m.mockExpectations = &ArtifactManagerMockUpdatePrototypeParams{p, p1, p2, p3, p4, p5}
	return m
}

//Return sets up a mock for ArtifactManager.UpdatePrototype to return Return's arguments
func (m *mArtifactManagerMockUpdatePrototype) Return(r core.ObjectDescriptor, r1 error) *ArtifactManagerMock {
	m.mock.UpdatePrototypeFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) (core.ObjectDescriptor, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of ArtifactManager.UpdatePrototype method
func (m *mArtifactManagerMockUpdatePrototype) Set(f func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) (r core.ObjectDescriptor, r1 error)) *ArtifactManagerMock {
	m.mock.UpdatePrototypeFunc = f
	m.mockExpectations = nil
	return m.mock
}

//UpdatePrototype implements github.com/insolar/insolar/core.ArtifactManager interface
func (m *ArtifactManagerMock) UpdatePrototype(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.ObjectDescriptor, p4 []byte, p5 *core.RecordRef) (r core.ObjectDescriptor, r1 error) {
	atomic.AddUint64(&m.UpdatePrototypePreCounter, 1)
	defer atomic.AddUint64(&m.UpdatePrototypeCounter, 1)

	if m.UpdatePrototypeMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.UpdatePrototypeMock.mockExpectations, ArtifactManagerMockUpdatePrototypeParams{p, p1, p2, p3, p4, p5},
			"ArtifactManager.UpdatePrototype got unexpected parameters")

		if m.UpdatePrototypeFunc == nil {

			m.t.Fatal("No results are set for the ArtifactManagerMock.UpdatePrototype")

			return
		}
	}

	if m.UpdatePrototypeFunc == nil {
		m.t.Fatal("Unexpected call to ArtifactManagerMock.UpdatePrototype")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ArtifactManagerMock) ValidateCallCounters() {

	if m.ActivateObjectFunc != nil && atomic.LoadUint64(&m.ActivateObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivateObject")
	}

	if m.ActivatePrototypeFunc != nil && atomic.LoadUint64(&m.ActivatePrototypeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivatePrototype")
	}

	if m.DeactivateObjectFunc != nil && atomic.LoadUint64(&m.DeactivateObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeactivateObject")
	}

	if m.DeclareTypeFunc != nil && atomic.LoadUint64(&m.DeclareTypeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeclareType")
	}

	if m.DeployCodeFunc != nil && atomic.LoadUint64(&m.DeployCodeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeployCode")
	}

	if m.GenesisRefFunc != nil && atomic.LoadUint64(&m.GenesisRefCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GenesisRef")
	}

	if m.GetChildrenFunc != nil && atomic.LoadUint64(&m.GetChildrenCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetChildren")
	}

	if m.GetCodeFunc != nil && atomic.LoadUint64(&m.GetCodeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetCode")
	}

	if m.GetDelegateFunc != nil && atomic.LoadUint64(&m.GetDelegateCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetDelegate")
	}

	if m.GetObjectFunc != nil && atomic.LoadUint64(&m.GetObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetObject")
	}

	if m.RegisterRequestFunc != nil && atomic.LoadUint64(&m.RegisterRequestCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterRequest")
	}

	if m.RegisterResultFunc != nil && atomic.LoadUint64(&m.RegisterResultCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterResult")
	}

	if m.RegisterValidationFunc != nil && atomic.LoadUint64(&m.RegisterValidationCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterValidation")
	}

	if m.StateFunc != nil && atomic.LoadUint64(&m.StateCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.State")
	}

	if m.UpdateObjectFunc != nil && atomic.LoadUint64(&m.UpdateObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.UpdateObject")
	}

	if m.UpdatePrototypeFunc != nil && atomic.LoadUint64(&m.UpdatePrototypeCounter) == 0 {
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

	if m.ActivateObjectFunc != nil && atomic.LoadUint64(&m.ActivateObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivateObject")
	}

	if m.ActivatePrototypeFunc != nil && atomic.LoadUint64(&m.ActivatePrototypeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.ActivatePrototype")
	}

	if m.DeactivateObjectFunc != nil && atomic.LoadUint64(&m.DeactivateObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeactivateObject")
	}

	if m.DeclareTypeFunc != nil && atomic.LoadUint64(&m.DeclareTypeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeclareType")
	}

	if m.DeployCodeFunc != nil && atomic.LoadUint64(&m.DeployCodeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.DeployCode")
	}

	if m.GenesisRefFunc != nil && atomic.LoadUint64(&m.GenesisRefCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GenesisRef")
	}

	if m.GetChildrenFunc != nil && atomic.LoadUint64(&m.GetChildrenCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetChildren")
	}

	if m.GetCodeFunc != nil && atomic.LoadUint64(&m.GetCodeCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetCode")
	}

	if m.GetDelegateFunc != nil && atomic.LoadUint64(&m.GetDelegateCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetDelegate")
	}

	if m.GetObjectFunc != nil && atomic.LoadUint64(&m.GetObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.GetObject")
	}

	if m.RegisterRequestFunc != nil && atomic.LoadUint64(&m.RegisterRequestCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterRequest")
	}

	if m.RegisterResultFunc != nil && atomic.LoadUint64(&m.RegisterResultCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterResult")
	}

	if m.RegisterValidationFunc != nil && atomic.LoadUint64(&m.RegisterValidationCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.RegisterValidation")
	}

	if m.StateFunc != nil && atomic.LoadUint64(&m.StateCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.State")
	}

	if m.UpdateObjectFunc != nil && atomic.LoadUint64(&m.UpdateObjectCounter) == 0 {
		m.t.Fatal("Expected call to ArtifactManagerMock.UpdateObject")
	}

	if m.UpdatePrototypeFunc != nil && atomic.LoadUint64(&m.UpdatePrototypeCounter) == 0 {
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
		ok = ok && (m.ActivateObjectFunc == nil || atomic.LoadUint64(&m.ActivateObjectCounter) > 0)
		ok = ok && (m.ActivatePrototypeFunc == nil || atomic.LoadUint64(&m.ActivatePrototypeCounter) > 0)
		ok = ok && (m.DeactivateObjectFunc == nil || atomic.LoadUint64(&m.DeactivateObjectCounter) > 0)
		ok = ok && (m.DeclareTypeFunc == nil || atomic.LoadUint64(&m.DeclareTypeCounter) > 0)
		ok = ok && (m.DeployCodeFunc == nil || atomic.LoadUint64(&m.DeployCodeCounter) > 0)
		ok = ok && (m.GenesisRefFunc == nil || atomic.LoadUint64(&m.GenesisRefCounter) > 0)
		ok = ok && (m.GetChildrenFunc == nil || atomic.LoadUint64(&m.GetChildrenCounter) > 0)
		ok = ok && (m.GetCodeFunc == nil || atomic.LoadUint64(&m.GetCodeCounter) > 0)
		ok = ok && (m.GetDelegateFunc == nil || atomic.LoadUint64(&m.GetDelegateCounter) > 0)
		ok = ok && (m.GetObjectFunc == nil || atomic.LoadUint64(&m.GetObjectCounter) > 0)
		ok = ok && (m.RegisterRequestFunc == nil || atomic.LoadUint64(&m.RegisterRequestCounter) > 0)
		ok = ok && (m.RegisterResultFunc == nil || atomic.LoadUint64(&m.RegisterResultCounter) > 0)
		ok = ok && (m.RegisterValidationFunc == nil || atomic.LoadUint64(&m.RegisterValidationCounter) > 0)
		ok = ok && (m.StateFunc == nil || atomic.LoadUint64(&m.StateCounter) > 0)
		ok = ok && (m.UpdateObjectFunc == nil || atomic.LoadUint64(&m.UpdateObjectCounter) > 0)
		ok = ok && (m.UpdatePrototypeFunc == nil || atomic.LoadUint64(&m.UpdatePrototypeCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ActivateObjectFunc != nil && atomic.LoadUint64(&m.ActivateObjectCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.ActivateObject")
			}

			if m.ActivatePrototypeFunc != nil && atomic.LoadUint64(&m.ActivatePrototypeCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.ActivatePrototype")
			}

			if m.DeactivateObjectFunc != nil && atomic.LoadUint64(&m.DeactivateObjectCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.DeactivateObject")
			}

			if m.DeclareTypeFunc != nil && atomic.LoadUint64(&m.DeclareTypeCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.DeclareType")
			}

			if m.DeployCodeFunc != nil && atomic.LoadUint64(&m.DeployCodeCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.DeployCode")
			}

			if m.GenesisRefFunc != nil && atomic.LoadUint64(&m.GenesisRefCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.GenesisRef")
			}

			if m.GetChildrenFunc != nil && atomic.LoadUint64(&m.GetChildrenCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.GetChildren")
			}

			if m.GetCodeFunc != nil && atomic.LoadUint64(&m.GetCodeCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.GetCode")
			}

			if m.GetDelegateFunc != nil && atomic.LoadUint64(&m.GetDelegateCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.GetDelegate")
			}

			if m.GetObjectFunc != nil && atomic.LoadUint64(&m.GetObjectCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.GetObject")
			}

			if m.RegisterRequestFunc != nil && atomic.LoadUint64(&m.RegisterRequestCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.RegisterRequest")
			}

			if m.RegisterResultFunc != nil && atomic.LoadUint64(&m.RegisterResultCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.RegisterResult")
			}

			if m.RegisterValidationFunc != nil && atomic.LoadUint64(&m.RegisterValidationCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.RegisterValidation")
			}

			if m.StateFunc != nil && atomic.LoadUint64(&m.StateCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.State")
			}

			if m.UpdateObjectFunc != nil && atomic.LoadUint64(&m.UpdateObjectCounter) == 0 {
				m.t.Error("Expected call to ArtifactManagerMock.UpdateObject")
			}

			if m.UpdatePrototypeFunc != nil && atomic.LoadUint64(&m.UpdatePrototypeCounter) == 0 {
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

	if m.ActivateObjectFunc != nil && atomic.LoadUint64(&m.ActivateObjectCounter) == 0 {
		return false
	}

	if m.ActivatePrototypeFunc != nil && atomic.LoadUint64(&m.ActivatePrototypeCounter) == 0 {
		return false
	}

	if m.DeactivateObjectFunc != nil && atomic.LoadUint64(&m.DeactivateObjectCounter) == 0 {
		return false
	}

	if m.DeclareTypeFunc != nil && atomic.LoadUint64(&m.DeclareTypeCounter) == 0 {
		return false
	}

	if m.DeployCodeFunc != nil && atomic.LoadUint64(&m.DeployCodeCounter) == 0 {
		return false
	}

	if m.GenesisRefFunc != nil && atomic.LoadUint64(&m.GenesisRefCounter) == 0 {
		return false
	}

	if m.GetChildrenFunc != nil && atomic.LoadUint64(&m.GetChildrenCounter) == 0 {
		return false
	}

	if m.GetCodeFunc != nil && atomic.LoadUint64(&m.GetCodeCounter) == 0 {
		return false
	}

	if m.GetDelegateFunc != nil && atomic.LoadUint64(&m.GetDelegateCounter) == 0 {
		return false
	}

	if m.GetObjectFunc != nil && atomic.LoadUint64(&m.GetObjectCounter) == 0 {
		return false
	}

	if m.RegisterRequestFunc != nil && atomic.LoadUint64(&m.RegisterRequestCounter) == 0 {
		return false
	}

	if m.RegisterResultFunc != nil && atomic.LoadUint64(&m.RegisterResultCounter) == 0 {
		return false
	}

	if m.RegisterValidationFunc != nil && atomic.LoadUint64(&m.RegisterValidationCounter) == 0 {
		return false
	}

	if m.StateFunc != nil && atomic.LoadUint64(&m.StateCounter) == 0 {
		return false
	}

	if m.UpdateObjectFunc != nil && atomic.LoadUint64(&m.UpdateObjectCounter) == 0 {
		return false
	}

	if m.UpdatePrototypeFunc != nil && atomic.LoadUint64(&m.UpdatePrototypeCounter) == 0 {
		return false
	}

	return true
}
