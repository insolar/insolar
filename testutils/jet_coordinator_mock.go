package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetCoordinator" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetCoordinatorMock implements github.com/insolar/insolar/core.JetCoordinator
type JetCoordinatorMock struct {
	t minimock.Tester

	HeavyFunc       func(p context.Context, p1 core.PulseNumber) (r *core.RecordRef, r1 error)
	HeavyCounter    uint64
	HeavyPreCounter uint64
	HeavyMock       mJetCoordinatorMockHeavy

	IsAuthorizedFunc       func(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error)
	IsAuthorizedCounter    uint64
	IsAuthorizedPreCounter uint64
	IsAuthorizedMock       mJetCoordinatorMockIsAuthorized

	IsBeyondLimitFunc       func(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) (r bool, r1 error)
	IsBeyondLimitCounter    uint64
	IsBeyondLimitPreCounter uint64
	IsBeyondLimitMock       mJetCoordinatorMockIsBeyondLimit

	LightExecutorForJetFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error)
	LightExecutorForJetCounter    uint64
	LightExecutorForJetPreCounter uint64
	LightExecutorForJetMock       mJetCoordinatorMockLightExecutorForJet

	LightExecutorForObjectFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error)
	LightExecutorForObjectCounter    uint64
	LightExecutorForObjectPreCounter uint64
	LightExecutorForObjectMock       mJetCoordinatorMockLightExecutorForObject

	LightValidatorsForJetFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error)
	LightValidatorsForJetCounter    uint64
	LightValidatorsForJetPreCounter uint64
	LightValidatorsForJetMock       mJetCoordinatorMockLightValidatorsForJet

	LightValidatorsForObjectFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error)
	LightValidatorsForObjectCounter    uint64
	LightValidatorsForObjectPreCounter uint64
	LightValidatorsForObjectMock       mJetCoordinatorMockLightValidatorsForObject

	MeFunc       func() (r core.RecordRef)
	MeCounter    uint64
	MePreCounter uint64
	MeMock       mJetCoordinatorMockMe

	NodeForJetFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (r *core.RecordRef, r1 error)
	NodeForJetCounter    uint64
	NodeForJetPreCounter uint64
	NodeForJetMock       mJetCoordinatorMockNodeForJet

	NodeForObjectFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (r *core.RecordRef, r1 error)
	NodeForObjectCounter    uint64
	NodeForObjectPreCounter uint64
	NodeForObjectMock       mJetCoordinatorMockNodeForObject

	QueryRoleFunc       func(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber) (r []core.RecordRef, r1 error)
	QueryRoleCounter    uint64
	QueryRolePreCounter uint64
	QueryRoleMock       mJetCoordinatorMockQueryRole

	VirtualExecutorForObjectFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error)
	VirtualExecutorForObjectCounter    uint64
	VirtualExecutorForObjectPreCounter uint64
	VirtualExecutorForObjectMock       mJetCoordinatorMockVirtualExecutorForObject

	VirtualValidatorsForObjectFunc       func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error)
	VirtualValidatorsForObjectCounter    uint64
	VirtualValidatorsForObjectPreCounter uint64
	VirtualValidatorsForObjectMock       mJetCoordinatorMockVirtualValidatorsForObject
}

//NewJetCoordinatorMock returns a mock for github.com/insolar/insolar/core.JetCoordinator
func NewJetCoordinatorMock(t minimock.Tester) *JetCoordinatorMock {
	m := &JetCoordinatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.HeavyMock = mJetCoordinatorMockHeavy{mock: m}
	m.IsAuthorizedMock = mJetCoordinatorMockIsAuthorized{mock: m}
	m.IsBeyondLimitMock = mJetCoordinatorMockIsBeyondLimit{mock: m}
	m.LightExecutorForJetMock = mJetCoordinatorMockLightExecutorForJet{mock: m}
	m.LightExecutorForObjectMock = mJetCoordinatorMockLightExecutorForObject{mock: m}
	m.LightValidatorsForJetMock = mJetCoordinatorMockLightValidatorsForJet{mock: m}
	m.LightValidatorsForObjectMock = mJetCoordinatorMockLightValidatorsForObject{mock: m}
	m.MeMock = mJetCoordinatorMockMe{mock: m}
	m.NodeForJetMock = mJetCoordinatorMockNodeForJet{mock: m}
	m.NodeForObjectMock = mJetCoordinatorMockNodeForObject{mock: m}
	m.QueryRoleMock = mJetCoordinatorMockQueryRole{mock: m}
	m.VirtualExecutorForObjectMock = mJetCoordinatorMockVirtualExecutorForObject{mock: m}
	m.VirtualValidatorsForObjectMock = mJetCoordinatorMockVirtualValidatorsForObject{mock: m}

	return m
}

type mJetCoordinatorMockHeavy struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockHeavyParams
}

//JetCoordinatorMockHeavyParams represents input parameters of the JetCoordinator.Heavy
type JetCoordinatorMockHeavyParams struct {
	p  context.Context
	p1 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.Heavy
func (m *mJetCoordinatorMockHeavy) Expect(p context.Context, p1 core.PulseNumber) *mJetCoordinatorMockHeavy {
	m.mockExpectations = &JetCoordinatorMockHeavyParams{p, p1}
	return m
}

//Return sets up a mock for JetCoordinator.Heavy to return Return's arguments
func (m *mJetCoordinatorMockHeavy) Return(r *core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.HeavyFunc = func(p context.Context, p1 core.PulseNumber) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.Heavy method
func (m *mJetCoordinatorMockHeavy) Set(f func(p context.Context, p1 core.PulseNumber) (r *core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.HeavyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Heavy implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) Heavy(p context.Context, p1 core.PulseNumber) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.HeavyPreCounter, 1)
	defer atomic.AddUint64(&m.HeavyCounter, 1)

	if m.HeavyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.HeavyMock.mockExpectations, JetCoordinatorMockHeavyParams{p, p1},
			"JetCoordinator.Heavy got unexpected parameters")

		if m.HeavyFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.Heavy")

			return
		}
	}

	if m.HeavyFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.Heavy")
		return
	}

	return m.HeavyFunc(p, p1)
}

//HeavyMinimockCounter returns a count of JetCoordinatorMock.HeavyFunc invocations
func (m *JetCoordinatorMock) HeavyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HeavyCounter)
}

//HeavyMinimockPreCounter returns the value of JetCoordinatorMock.Heavy invocations
func (m *JetCoordinatorMock) HeavyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HeavyPreCounter)
}

type mJetCoordinatorMockIsAuthorized struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockIsAuthorizedParams
}

//JetCoordinatorMockIsAuthorizedParams represents input parameters of the JetCoordinator.IsAuthorized
type JetCoordinatorMockIsAuthorizedParams struct {
	p  context.Context
	p1 core.DynamicRole
	p2 core.RecordID
	p3 core.PulseNumber
	p4 core.RecordRef
}

//Expect sets up expected params for the JetCoordinator.IsAuthorized
func (m *mJetCoordinatorMockIsAuthorized) Expect(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber, p4 core.RecordRef) *mJetCoordinatorMockIsAuthorized {
	m.mockExpectations = &JetCoordinatorMockIsAuthorizedParams{p, p1, p2, p3, p4}
	return m
}

//Return sets up a mock for JetCoordinator.IsAuthorized to return Return's arguments
func (m *mJetCoordinatorMockIsAuthorized) Return(r bool, r1 error) *JetCoordinatorMock {
	m.mock.IsAuthorizedFunc = func(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber, p4 core.RecordRef) (bool, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.IsAuthorized method
func (m *mJetCoordinatorMockIsAuthorized) Set(f func(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error)) *JetCoordinatorMock {
	m.mock.IsAuthorizedFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IsAuthorized implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) IsAuthorized(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber, p4 core.RecordRef) (r bool, r1 error) {
	atomic.AddUint64(&m.IsAuthorizedPreCounter, 1)
	defer atomic.AddUint64(&m.IsAuthorizedCounter, 1)

	if m.IsAuthorizedMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IsAuthorizedMock.mockExpectations, JetCoordinatorMockIsAuthorizedParams{p, p1, p2, p3, p4},
			"JetCoordinator.IsAuthorized got unexpected parameters")

		if m.IsAuthorizedFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.IsAuthorized")

			return
		}
	}

	if m.IsAuthorizedFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.IsAuthorized")
		return
	}

	return m.IsAuthorizedFunc(p, p1, p2, p3, p4)
}

//IsAuthorizedMinimockCounter returns a count of JetCoordinatorMock.IsAuthorizedFunc invocations
func (m *JetCoordinatorMock) IsAuthorizedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsAuthorizedCounter)
}

//IsAuthorizedMinimockPreCounter returns the value of JetCoordinatorMock.IsAuthorized invocations
func (m *JetCoordinatorMock) IsAuthorizedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsAuthorizedPreCounter)
}

type mJetCoordinatorMockIsBeyondLimit struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockIsBeyondLimitParams
}

//JetCoordinatorMockIsBeyondLimitParams represents input parameters of the JetCoordinator.IsBeyondLimit
type JetCoordinatorMockIsBeyondLimitParams struct {
	p  context.Context
	p1 core.PulseNumber
	p2 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.IsBeyondLimit
func (m *mJetCoordinatorMockIsBeyondLimit) Expect(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) *mJetCoordinatorMockIsBeyondLimit {
	m.mockExpectations = &JetCoordinatorMockIsBeyondLimitParams{p, p1, p2}
	return m
}

//Return sets up a mock for JetCoordinator.IsBeyondLimit to return Return's arguments
func (m *mJetCoordinatorMockIsBeyondLimit) Return(r bool, r1 error) *JetCoordinatorMock {
	m.mock.IsBeyondLimitFunc = func(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) (bool, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.IsBeyondLimit method
func (m *mJetCoordinatorMockIsBeyondLimit) Set(f func(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) (r bool, r1 error)) *JetCoordinatorMock {
	m.mock.IsBeyondLimitFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IsBeyondLimit implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) IsBeyondLimit(p context.Context, p1 core.PulseNumber, p2 core.PulseNumber) (r bool, r1 error) {
	atomic.AddUint64(&m.IsBeyondLimitPreCounter, 1)
	defer atomic.AddUint64(&m.IsBeyondLimitCounter, 1)

	if m.IsBeyondLimitMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IsBeyondLimitMock.mockExpectations, JetCoordinatorMockIsBeyondLimitParams{p, p1, p2},
			"JetCoordinator.IsBeyondLimit got unexpected parameters")

		if m.IsBeyondLimitFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.IsBeyondLimit")

			return
		}
	}

	if m.IsBeyondLimitFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.IsBeyondLimit")
		return
	}

	return m.IsBeyondLimitFunc(p, p1, p2)
}

//IsBeyondLimitMinimockCounter returns a count of JetCoordinatorMock.IsBeyondLimitFunc invocations
func (m *JetCoordinatorMock) IsBeyondLimitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsBeyondLimitCounter)
}

//IsBeyondLimitMinimockPreCounter returns the value of JetCoordinatorMock.IsBeyondLimit invocations
func (m *JetCoordinatorMock) IsBeyondLimitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsBeyondLimitPreCounter)
}

type mJetCoordinatorMockLightExecutorForJet struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockLightExecutorForJetParams
}

//JetCoordinatorMockLightExecutorForJetParams represents input parameters of the JetCoordinator.LightExecutorForJet
type JetCoordinatorMockLightExecutorForJetParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.LightExecutorForJet
func (m *mJetCoordinatorMockLightExecutorForJet) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mJetCoordinatorMockLightExecutorForJet {
	m.mockExpectations = &JetCoordinatorMockLightExecutorForJetParams{p, p1, p2}
	return m
}

//Return sets up a mock for JetCoordinator.LightExecutorForJet to return Return's arguments
func (m *mJetCoordinatorMockLightExecutorForJet) Return(r *core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.LightExecutorForJetFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.LightExecutorForJet method
func (m *mJetCoordinatorMockLightExecutorForJet) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.LightExecutorForJetFunc = f
	m.mockExpectations = nil
	return m.mock
}

//LightExecutorForJet implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) LightExecutorForJet(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.LightExecutorForJetPreCounter, 1)
	defer atomic.AddUint64(&m.LightExecutorForJetCounter, 1)

	if m.LightExecutorForJetMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LightExecutorForJetMock.mockExpectations, JetCoordinatorMockLightExecutorForJetParams{p, p1, p2},
			"JetCoordinator.LightExecutorForJet got unexpected parameters")

		if m.LightExecutorForJetFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.LightExecutorForJet")

			return
		}
	}

	if m.LightExecutorForJetFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.LightExecutorForJet")
		return
	}

	return m.LightExecutorForJetFunc(p, p1, p2)
}

//LightExecutorForJetMinimockCounter returns a count of JetCoordinatorMock.LightExecutorForJetFunc invocations
func (m *JetCoordinatorMock) LightExecutorForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForJetCounter)
}

//LightExecutorForJetMinimockPreCounter returns the value of JetCoordinatorMock.LightExecutorForJet invocations
func (m *JetCoordinatorMock) LightExecutorForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForJetPreCounter)
}

type mJetCoordinatorMockLightExecutorForObject struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockLightExecutorForObjectParams
}

//JetCoordinatorMockLightExecutorForObjectParams represents input parameters of the JetCoordinator.LightExecutorForObject
type JetCoordinatorMockLightExecutorForObjectParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.LightExecutorForObject
func (m *mJetCoordinatorMockLightExecutorForObject) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mJetCoordinatorMockLightExecutorForObject {
	m.mockExpectations = &JetCoordinatorMockLightExecutorForObjectParams{p, p1, p2}
	return m
}

//Return sets up a mock for JetCoordinator.LightExecutorForObject to return Return's arguments
func (m *mJetCoordinatorMockLightExecutorForObject) Return(r *core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.LightExecutorForObjectFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.LightExecutorForObject method
func (m *mJetCoordinatorMockLightExecutorForObject) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.LightExecutorForObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//LightExecutorForObject implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) LightExecutorForObject(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.LightExecutorForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.LightExecutorForObjectCounter, 1)

	if m.LightExecutorForObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LightExecutorForObjectMock.mockExpectations, JetCoordinatorMockLightExecutorForObjectParams{p, p1, p2},
			"JetCoordinator.LightExecutorForObject got unexpected parameters")

		if m.LightExecutorForObjectFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.LightExecutorForObject")

			return
		}
	}

	if m.LightExecutorForObjectFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.LightExecutorForObject")
		return
	}

	return m.LightExecutorForObjectFunc(p, p1, p2)
}

//LightExecutorForObjectMinimockCounter returns a count of JetCoordinatorMock.LightExecutorForObjectFunc invocations
func (m *JetCoordinatorMock) LightExecutorForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForObjectCounter)
}

//LightExecutorForObjectMinimockPreCounter returns the value of JetCoordinatorMock.LightExecutorForObject invocations
func (m *JetCoordinatorMock) LightExecutorForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightExecutorForObjectPreCounter)
}

type mJetCoordinatorMockLightValidatorsForJet struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockLightValidatorsForJetParams
}

//JetCoordinatorMockLightValidatorsForJetParams represents input parameters of the JetCoordinator.LightValidatorsForJet
type JetCoordinatorMockLightValidatorsForJetParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.LightValidatorsForJet
func (m *mJetCoordinatorMockLightValidatorsForJet) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mJetCoordinatorMockLightValidatorsForJet {
	m.mockExpectations = &JetCoordinatorMockLightValidatorsForJetParams{p, p1, p2}
	return m
}

//Return sets up a mock for JetCoordinator.LightValidatorsForJet to return Return's arguments
func (m *mJetCoordinatorMockLightValidatorsForJet) Return(r []core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.LightValidatorsForJetFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) ([]core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.LightValidatorsForJet method
func (m *mJetCoordinatorMockLightValidatorsForJet) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.LightValidatorsForJetFunc = f
	m.mockExpectations = nil
	return m.mock
}

//LightValidatorsForJet implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) LightValidatorsForJet(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error) {
	atomic.AddUint64(&m.LightValidatorsForJetPreCounter, 1)
	defer atomic.AddUint64(&m.LightValidatorsForJetCounter, 1)

	if m.LightValidatorsForJetMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LightValidatorsForJetMock.mockExpectations, JetCoordinatorMockLightValidatorsForJetParams{p, p1, p2},
			"JetCoordinator.LightValidatorsForJet got unexpected parameters")

		if m.LightValidatorsForJetFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.LightValidatorsForJet")

			return
		}
	}

	if m.LightValidatorsForJetFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.LightValidatorsForJet")
		return
	}

	return m.LightValidatorsForJetFunc(p, p1, p2)
}

//LightValidatorsForJetMinimockCounter returns a count of JetCoordinatorMock.LightValidatorsForJetFunc invocations
func (m *JetCoordinatorMock) LightValidatorsForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForJetCounter)
}

//LightValidatorsForJetMinimockPreCounter returns the value of JetCoordinatorMock.LightValidatorsForJet invocations
func (m *JetCoordinatorMock) LightValidatorsForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForJetPreCounter)
}

type mJetCoordinatorMockLightValidatorsForObject struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockLightValidatorsForObjectParams
}

//JetCoordinatorMockLightValidatorsForObjectParams represents input parameters of the JetCoordinator.LightValidatorsForObject
type JetCoordinatorMockLightValidatorsForObjectParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.LightValidatorsForObject
func (m *mJetCoordinatorMockLightValidatorsForObject) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mJetCoordinatorMockLightValidatorsForObject {
	m.mockExpectations = &JetCoordinatorMockLightValidatorsForObjectParams{p, p1, p2}
	return m
}

//Return sets up a mock for JetCoordinator.LightValidatorsForObject to return Return's arguments
func (m *mJetCoordinatorMockLightValidatorsForObject) Return(r []core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.LightValidatorsForObjectFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) ([]core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.LightValidatorsForObject method
func (m *mJetCoordinatorMockLightValidatorsForObject) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.LightValidatorsForObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//LightValidatorsForObject implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) LightValidatorsForObject(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error) {
	atomic.AddUint64(&m.LightValidatorsForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.LightValidatorsForObjectCounter, 1)

	if m.LightValidatorsForObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LightValidatorsForObjectMock.mockExpectations, JetCoordinatorMockLightValidatorsForObjectParams{p, p1, p2},
			"JetCoordinator.LightValidatorsForObject got unexpected parameters")

		if m.LightValidatorsForObjectFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.LightValidatorsForObject")

			return
		}
	}

	if m.LightValidatorsForObjectFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.LightValidatorsForObject")
		return
	}

	return m.LightValidatorsForObjectFunc(p, p1, p2)
}

//LightValidatorsForObjectMinimockCounter returns a count of JetCoordinatorMock.LightValidatorsForObjectFunc invocations
func (m *JetCoordinatorMock) LightValidatorsForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForObjectCounter)
}

//LightValidatorsForObjectMinimockPreCounter returns the value of JetCoordinatorMock.LightValidatorsForObject invocations
func (m *JetCoordinatorMock) LightValidatorsForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LightValidatorsForObjectPreCounter)
}

type mJetCoordinatorMockMe struct {
	mock *JetCoordinatorMock
}

//Return sets up a mock for JetCoordinator.Me to return Return's arguments
func (m *mJetCoordinatorMockMe) Return(r core.RecordRef) *JetCoordinatorMock {
	m.mock.MeFunc = func() core.RecordRef {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.Me method
func (m *mJetCoordinatorMockMe) Set(f func() (r core.RecordRef)) *JetCoordinatorMock {
	m.mock.MeFunc = f

	return m.mock
}

//Me implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) Me() (r core.RecordRef) {
	atomic.AddUint64(&m.MePreCounter, 1)
	defer atomic.AddUint64(&m.MeCounter, 1)

	if m.MeFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.Me")
		return
	}

	return m.MeFunc()
}

//MeMinimockCounter returns a count of JetCoordinatorMock.MeFunc invocations
func (m *JetCoordinatorMock) MeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MeCounter)
}

//MeMinimockPreCounter returns the value of JetCoordinatorMock.Me invocations
func (m *JetCoordinatorMock) MeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MePreCounter)
}

type mJetCoordinatorMockNodeForJet struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockNodeForJetParams
}

//JetCoordinatorMockNodeForJetParams represents input parameters of the JetCoordinator.NodeForJet
type JetCoordinatorMockNodeForJetParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.NodeForJet
func (m *mJetCoordinatorMockNodeForJet) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) *mJetCoordinatorMockNodeForJet {
	m.mockExpectations = &JetCoordinatorMockNodeForJetParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for JetCoordinator.NodeForJet to return Return's arguments
func (m *mJetCoordinatorMockNodeForJet) Return(r *core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.NodeForJetFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.NodeForJet method
func (m *mJetCoordinatorMockNodeForJet) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (r *core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.NodeForJetFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NodeForJet implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) NodeForJet(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.NodeForJetPreCounter, 1)
	defer atomic.AddUint64(&m.NodeForJetCounter, 1)

	if m.NodeForJetMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NodeForJetMock.mockExpectations, JetCoordinatorMockNodeForJetParams{p, p1, p2, p3},
			"JetCoordinator.NodeForJet got unexpected parameters")

		if m.NodeForJetFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.NodeForJet")

			return
		}
	}

	if m.NodeForJetFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.NodeForJet")
		return
	}

	return m.NodeForJetFunc(p, p1, p2, p3)
}

//NodeForJetMinimockCounter returns a count of JetCoordinatorMock.NodeForJetFunc invocations
func (m *JetCoordinatorMock) NodeForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForJetCounter)
}

//NodeForJetMinimockPreCounter returns the value of JetCoordinatorMock.NodeForJet invocations
func (m *JetCoordinatorMock) NodeForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForJetPreCounter)
}

type mJetCoordinatorMockNodeForObject struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockNodeForObjectParams
}

//JetCoordinatorMockNodeForObjectParams represents input parameters of the JetCoordinator.NodeForObject
type JetCoordinatorMockNodeForObjectParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
	p3 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.NodeForObject
func (m *mJetCoordinatorMockNodeForObject) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) *mJetCoordinatorMockNodeForObject {
	m.mockExpectations = &JetCoordinatorMockNodeForObjectParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for JetCoordinator.NodeForObject to return Return's arguments
func (m *mJetCoordinatorMockNodeForObject) Return(r *core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.NodeForObjectFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.NodeForObject method
func (m *mJetCoordinatorMockNodeForObject) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (r *core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.NodeForObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//NodeForObject implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) NodeForObject(p context.Context, p1 core.RecordID, p2 core.PulseNumber, p3 core.PulseNumber) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.NodeForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.NodeForObjectCounter, 1)

	if m.NodeForObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.NodeForObjectMock.mockExpectations, JetCoordinatorMockNodeForObjectParams{p, p1, p2, p3},
			"JetCoordinator.NodeForObject got unexpected parameters")

		if m.NodeForObjectFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.NodeForObject")

			return
		}
	}

	if m.NodeForObjectFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.NodeForObject")
		return
	}

	return m.NodeForObjectFunc(p, p1, p2, p3)
}

//NodeForObjectMinimockCounter returns a count of JetCoordinatorMock.NodeForObjectFunc invocations
func (m *JetCoordinatorMock) NodeForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForObjectCounter)
}

//NodeForObjectMinimockPreCounter returns the value of JetCoordinatorMock.NodeForObject invocations
func (m *JetCoordinatorMock) NodeForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NodeForObjectPreCounter)
}

type mJetCoordinatorMockQueryRole struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockQueryRoleParams
}

//JetCoordinatorMockQueryRoleParams represents input parameters of the JetCoordinator.QueryRole
type JetCoordinatorMockQueryRoleParams struct {
	p  context.Context
	p1 core.DynamicRole
	p2 core.RecordID
	p3 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.QueryRole
func (m *mJetCoordinatorMockQueryRole) Expect(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber) *mJetCoordinatorMockQueryRole {
	m.mockExpectations = &JetCoordinatorMockQueryRoleParams{p, p1, p2, p3}
	return m
}

//Return sets up a mock for JetCoordinator.QueryRole to return Return's arguments
func (m *mJetCoordinatorMockQueryRole) Return(r []core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.QueryRoleFunc = func(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber) ([]core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.QueryRole method
func (m *mJetCoordinatorMockQueryRole) Set(f func(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber) (r []core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.QueryRoleFunc = f
	m.mockExpectations = nil
	return m.mock
}

//QueryRole implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) QueryRole(p context.Context, p1 core.DynamicRole, p2 core.RecordID, p3 core.PulseNumber) (r []core.RecordRef, r1 error) {
	atomic.AddUint64(&m.QueryRolePreCounter, 1)
	defer atomic.AddUint64(&m.QueryRoleCounter, 1)

	if m.QueryRoleMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.QueryRoleMock.mockExpectations, JetCoordinatorMockQueryRoleParams{p, p1, p2, p3},
			"JetCoordinator.QueryRole got unexpected parameters")

		if m.QueryRoleFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.QueryRole")

			return
		}
	}

	if m.QueryRoleFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.QueryRole")
		return
	}

	return m.QueryRoleFunc(p, p1, p2, p3)
}

//QueryRoleMinimockCounter returns a count of JetCoordinatorMock.QueryRoleFunc invocations
func (m *JetCoordinatorMock) QueryRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.QueryRoleCounter)
}

//QueryRoleMinimockPreCounter returns the value of JetCoordinatorMock.QueryRole invocations
func (m *JetCoordinatorMock) QueryRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.QueryRolePreCounter)
}

type mJetCoordinatorMockVirtualExecutorForObject struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockVirtualExecutorForObjectParams
}

//JetCoordinatorMockVirtualExecutorForObjectParams represents input parameters of the JetCoordinator.VirtualExecutorForObject
type JetCoordinatorMockVirtualExecutorForObjectParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.VirtualExecutorForObject
func (m *mJetCoordinatorMockVirtualExecutorForObject) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mJetCoordinatorMockVirtualExecutorForObject {
	m.mockExpectations = &JetCoordinatorMockVirtualExecutorForObjectParams{p, p1, p2}
	return m
}

//Return sets up a mock for JetCoordinator.VirtualExecutorForObject to return Return's arguments
func (m *mJetCoordinatorMockVirtualExecutorForObject) Return(r *core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.VirtualExecutorForObjectFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (*core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.VirtualExecutorForObject method
func (m *mJetCoordinatorMockVirtualExecutorForObject) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.VirtualExecutorForObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//VirtualExecutorForObject implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) VirtualExecutorForObject(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r *core.RecordRef, r1 error) {
	atomic.AddUint64(&m.VirtualExecutorForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.VirtualExecutorForObjectCounter, 1)

	if m.VirtualExecutorForObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.VirtualExecutorForObjectMock.mockExpectations, JetCoordinatorMockVirtualExecutorForObjectParams{p, p1, p2},
			"JetCoordinator.VirtualExecutorForObject got unexpected parameters")

		if m.VirtualExecutorForObjectFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.VirtualExecutorForObject")

			return
		}
	}

	if m.VirtualExecutorForObjectFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.VirtualExecutorForObject")
		return
	}

	return m.VirtualExecutorForObjectFunc(p, p1, p2)
}

//VirtualExecutorForObjectMinimockCounter returns a count of JetCoordinatorMock.VirtualExecutorForObjectFunc invocations
func (m *JetCoordinatorMock) VirtualExecutorForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualExecutorForObjectCounter)
}

//VirtualExecutorForObjectMinimockPreCounter returns the value of JetCoordinatorMock.VirtualExecutorForObject invocations
func (m *JetCoordinatorMock) VirtualExecutorForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualExecutorForObjectPreCounter)
}

type mJetCoordinatorMockVirtualValidatorsForObject struct {
	mock             *JetCoordinatorMock
	mockExpectations *JetCoordinatorMockVirtualValidatorsForObjectParams
}

//JetCoordinatorMockVirtualValidatorsForObjectParams represents input parameters of the JetCoordinator.VirtualValidatorsForObject
type JetCoordinatorMockVirtualValidatorsForObjectParams struct {
	p  context.Context
	p1 core.RecordID
	p2 core.PulseNumber
}

//Expect sets up expected params for the JetCoordinator.VirtualValidatorsForObject
func (m *mJetCoordinatorMockVirtualValidatorsForObject) Expect(p context.Context, p1 core.RecordID, p2 core.PulseNumber) *mJetCoordinatorMockVirtualValidatorsForObject {
	m.mockExpectations = &JetCoordinatorMockVirtualValidatorsForObjectParams{p, p1, p2}
	return m
}

//Return sets up a mock for JetCoordinator.VirtualValidatorsForObject to return Return's arguments
func (m *mJetCoordinatorMockVirtualValidatorsForObject) Return(r []core.RecordRef, r1 error) *JetCoordinatorMock {
	m.mock.VirtualValidatorsForObjectFunc = func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) ([]core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of JetCoordinator.VirtualValidatorsForObject method
func (m *mJetCoordinatorMockVirtualValidatorsForObject) Set(f func(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error)) *JetCoordinatorMock {
	m.mock.VirtualValidatorsForObjectFunc = f
	m.mockExpectations = nil
	return m.mock
}

//VirtualValidatorsForObject implements github.com/insolar/insolar/core.JetCoordinator interface
func (m *JetCoordinatorMock) VirtualValidatorsForObject(p context.Context, p1 core.RecordID, p2 core.PulseNumber) (r []core.RecordRef, r1 error) {
	atomic.AddUint64(&m.VirtualValidatorsForObjectPreCounter, 1)
	defer atomic.AddUint64(&m.VirtualValidatorsForObjectCounter, 1)

	if m.VirtualValidatorsForObjectMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.VirtualValidatorsForObjectMock.mockExpectations, JetCoordinatorMockVirtualValidatorsForObjectParams{p, p1, p2},
			"JetCoordinator.VirtualValidatorsForObject got unexpected parameters")

		if m.VirtualValidatorsForObjectFunc == nil {

			m.t.Fatal("No results are set for the JetCoordinatorMock.VirtualValidatorsForObject")

			return
		}
	}

	if m.VirtualValidatorsForObjectFunc == nil {
		m.t.Fatal("Unexpected call to JetCoordinatorMock.VirtualValidatorsForObject")
		return
	}

	return m.VirtualValidatorsForObjectFunc(p, p1, p2)
}

//VirtualValidatorsForObjectMinimockCounter returns a count of JetCoordinatorMock.VirtualValidatorsForObjectFunc invocations
func (m *JetCoordinatorMock) VirtualValidatorsForObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter)
}

//VirtualValidatorsForObjectMinimockPreCounter returns the value of JetCoordinatorMock.VirtualValidatorsForObject invocations
func (m *JetCoordinatorMock) VirtualValidatorsForObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VirtualValidatorsForObjectPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetCoordinatorMock) ValidateCallCounters() {

	if m.HeavyFunc != nil && atomic.LoadUint64(&m.HeavyCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.Heavy")
	}

	if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsAuthorized")
	}

	if m.IsBeyondLimitFunc != nil && atomic.LoadUint64(&m.IsBeyondLimitCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsBeyondLimit")
	}

	if m.LightExecutorForJetFunc != nil && atomic.LoadUint64(&m.LightExecutorForJetCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightExecutorForJet")
	}

	if m.LightExecutorForObjectFunc != nil && atomic.LoadUint64(&m.LightExecutorForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightExecutorForObject")
	}

	if m.LightValidatorsForJetFunc != nil && atomic.LoadUint64(&m.LightValidatorsForJetCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightValidatorsForJet")
	}

	if m.LightValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.LightValidatorsForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightValidatorsForObject")
	}

	if m.MeFunc != nil && atomic.LoadUint64(&m.MeCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.Me")
	}

	if m.NodeForJetFunc != nil && atomic.LoadUint64(&m.NodeForJetCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.NodeForJet")
	}

	if m.NodeForObjectFunc != nil && atomic.LoadUint64(&m.NodeForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.NodeForObject")
	}

	if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.QueryRole")
	}

	if m.VirtualExecutorForObjectFunc != nil && atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.VirtualExecutorForObject")
	}

	if m.VirtualValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.VirtualValidatorsForObject")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetCoordinatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetCoordinatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetCoordinatorMock) MinimockFinish() {

	if m.HeavyFunc != nil && atomic.LoadUint64(&m.HeavyCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.Heavy")
	}

	if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsAuthorized")
	}

	if m.IsBeyondLimitFunc != nil && atomic.LoadUint64(&m.IsBeyondLimitCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.IsBeyondLimit")
	}

	if m.LightExecutorForJetFunc != nil && atomic.LoadUint64(&m.LightExecutorForJetCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightExecutorForJet")
	}

	if m.LightExecutorForObjectFunc != nil && atomic.LoadUint64(&m.LightExecutorForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightExecutorForObject")
	}

	if m.LightValidatorsForJetFunc != nil && atomic.LoadUint64(&m.LightValidatorsForJetCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightValidatorsForJet")
	}

	if m.LightValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.LightValidatorsForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.LightValidatorsForObject")
	}

	if m.MeFunc != nil && atomic.LoadUint64(&m.MeCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.Me")
	}

	if m.NodeForJetFunc != nil && atomic.LoadUint64(&m.NodeForJetCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.NodeForJet")
	}

	if m.NodeForObjectFunc != nil && atomic.LoadUint64(&m.NodeForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.NodeForObject")
	}

	if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.QueryRole")
	}

	if m.VirtualExecutorForObjectFunc != nil && atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.VirtualExecutorForObject")
	}

	if m.VirtualValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) == 0 {
		m.t.Fatal("Expected call to JetCoordinatorMock.VirtualValidatorsForObject")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetCoordinatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetCoordinatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.HeavyFunc == nil || atomic.LoadUint64(&m.HeavyCounter) > 0)
		ok = ok && (m.IsAuthorizedFunc == nil || atomic.LoadUint64(&m.IsAuthorizedCounter) > 0)
		ok = ok && (m.IsBeyondLimitFunc == nil || atomic.LoadUint64(&m.IsBeyondLimitCounter) > 0)
		ok = ok && (m.LightExecutorForJetFunc == nil || atomic.LoadUint64(&m.LightExecutorForJetCounter) > 0)
		ok = ok && (m.LightExecutorForObjectFunc == nil || atomic.LoadUint64(&m.LightExecutorForObjectCounter) > 0)
		ok = ok && (m.LightValidatorsForJetFunc == nil || atomic.LoadUint64(&m.LightValidatorsForJetCounter) > 0)
		ok = ok && (m.LightValidatorsForObjectFunc == nil || atomic.LoadUint64(&m.LightValidatorsForObjectCounter) > 0)
		ok = ok && (m.MeFunc == nil || atomic.LoadUint64(&m.MeCounter) > 0)
		ok = ok && (m.NodeForJetFunc == nil || atomic.LoadUint64(&m.NodeForJetCounter) > 0)
		ok = ok && (m.NodeForObjectFunc == nil || atomic.LoadUint64(&m.NodeForObjectCounter) > 0)
		ok = ok && (m.QueryRoleFunc == nil || atomic.LoadUint64(&m.QueryRoleCounter) > 0)
		ok = ok && (m.VirtualExecutorForObjectFunc == nil || atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) > 0)
		ok = ok && (m.VirtualValidatorsForObjectFunc == nil || atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.HeavyFunc != nil && atomic.LoadUint64(&m.HeavyCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.Heavy")
			}

			if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.IsAuthorized")
			}

			if m.IsBeyondLimitFunc != nil && atomic.LoadUint64(&m.IsBeyondLimitCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.IsBeyondLimit")
			}

			if m.LightExecutorForJetFunc != nil && atomic.LoadUint64(&m.LightExecutorForJetCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.LightExecutorForJet")
			}

			if m.LightExecutorForObjectFunc != nil && atomic.LoadUint64(&m.LightExecutorForObjectCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.LightExecutorForObject")
			}

			if m.LightValidatorsForJetFunc != nil && atomic.LoadUint64(&m.LightValidatorsForJetCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.LightValidatorsForJet")
			}

			if m.LightValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.LightValidatorsForObjectCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.LightValidatorsForObject")
			}

			if m.MeFunc != nil && atomic.LoadUint64(&m.MeCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.Me")
			}

			if m.NodeForJetFunc != nil && atomic.LoadUint64(&m.NodeForJetCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.NodeForJet")
			}

			if m.NodeForObjectFunc != nil && atomic.LoadUint64(&m.NodeForObjectCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.NodeForObject")
			}

			if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.QueryRole")
			}

			if m.VirtualExecutorForObjectFunc != nil && atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.VirtualExecutorForObject")
			}

			if m.VirtualValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) == 0 {
				m.t.Error("Expected call to JetCoordinatorMock.VirtualValidatorsForObject")
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
func (m *JetCoordinatorMock) AllMocksCalled() bool {

	if m.HeavyFunc != nil && atomic.LoadUint64(&m.HeavyCounter) == 0 {
		return false
	}

	if m.IsAuthorizedFunc != nil && atomic.LoadUint64(&m.IsAuthorizedCounter) == 0 {
		return false
	}

	if m.IsBeyondLimitFunc != nil && atomic.LoadUint64(&m.IsBeyondLimitCounter) == 0 {
		return false
	}

	if m.LightExecutorForJetFunc != nil && atomic.LoadUint64(&m.LightExecutorForJetCounter) == 0 {
		return false
	}

	if m.LightExecutorForObjectFunc != nil && atomic.LoadUint64(&m.LightExecutorForObjectCounter) == 0 {
		return false
	}

	if m.LightValidatorsForJetFunc != nil && atomic.LoadUint64(&m.LightValidatorsForJetCounter) == 0 {
		return false
	}

	if m.LightValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.LightValidatorsForObjectCounter) == 0 {
		return false
	}

	if m.MeFunc != nil && atomic.LoadUint64(&m.MeCounter) == 0 {
		return false
	}

	if m.NodeForJetFunc != nil && atomic.LoadUint64(&m.NodeForJetCounter) == 0 {
		return false
	}

	if m.NodeForObjectFunc != nil && atomic.LoadUint64(&m.NodeForObjectCounter) == 0 {
		return false
	}

	if m.QueryRoleFunc != nil && atomic.LoadUint64(&m.QueryRoleCounter) == 0 {
		return false
	}

	if m.VirtualExecutorForObjectFunc != nil && atomic.LoadUint64(&m.VirtualExecutorForObjectCounter) == 0 {
		return false
	}

	if m.VirtualValidatorsForObjectFunc != nil && atomic.LoadUint64(&m.VirtualValidatorsForObjectCounter) == 0 {
		return false
	}

	return true
}
