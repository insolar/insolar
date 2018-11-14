package merkle

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Calculator" can be found in github.com/insolar/insolar/network/merkle
*/
import (
	context "context"
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	merkle "github.com/insolar/insolar/network/merkle"

	testify_assert "github.com/stretchr/testify/assert"
)

//CalculatorMock implements github.com/insolar/insolar/network/merkle.Calculator
type CalculatorMock struct {
	t minimock.Tester

	GetCloudProofFunc       func(p context.Context, p1 *merkle.CloudEntry) (r merkle.OriginHash, r1 *merkle.CloudProof, r2 error)
	GetCloudProofCounter    uint64
	GetCloudProofPreCounter uint64
	GetCloudProofMock       mCalculatorMockGetCloudProof

	GetGlobuleProofFunc       func(p context.Context, p1 *merkle.GlobuleEntry) (r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error)
	GetGlobuleProofCounter    uint64
	GetGlobuleProofPreCounter uint64
	GetGlobuleProofMock       mCalculatorMockGetGlobuleProof

	GetPulseProofFunc       func(p context.Context, p1 *merkle.PulseEntry) (r merkle.OriginHash, r1 *merkle.PulseProof, r2 error)
	GetPulseProofCounter    uint64
	GetPulseProofPreCounter uint64
	GetPulseProofMock       mCalculatorMockGetPulseProof

	IsValidFunc       func(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) (r bool)
	IsValidCounter    uint64
	IsValidPreCounter uint64
	IsValidMock       mCalculatorMockIsValid
}

//NewCalculatorMock returns a mock for github.com/insolar/insolar/network/merkle.Calculator
func NewCalculatorMock(t minimock.Tester) *CalculatorMock {
	m := &CalculatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCloudProofMock = mCalculatorMockGetCloudProof{mock: m}
	m.GetGlobuleProofMock = mCalculatorMockGetGlobuleProof{mock: m}
	m.GetPulseProofMock = mCalculatorMockGetPulseProof{mock: m}
	m.IsValidMock = mCalculatorMockIsValid{mock: m}

	return m
}

type mCalculatorMockGetCloudProof struct {
	mock             *CalculatorMock
	mockExpectations *CalculatorMockGetCloudProofParams
}

//CalculatorMockGetCloudProofParams represents input parameters of the Calculator.GetCloudProof
type CalculatorMockGetCloudProofParams struct {
	p  context.Context
	p1 *merkle.CloudEntry
}

//Expect sets up expected params for the Calculator.GetCloudProof
func (m *mCalculatorMockGetCloudProof) Expect(p context.Context, p1 *merkle.CloudEntry) *mCalculatorMockGetCloudProof {
	m.mockExpectations = &CalculatorMockGetCloudProofParams{p, p1}
	return m
}

//Return sets up a mock for Calculator.GetCloudProof to return Return's arguments
func (m *mCalculatorMockGetCloudProof) Return(r merkle.OriginHash, r1 *merkle.CloudProof, r2 error) *CalculatorMock {
	m.mock.GetCloudProofFunc = func(p context.Context, p1 *merkle.CloudEntry) (merkle.OriginHash, *merkle.CloudProof, error) {
		return r, r1, r2
	}
	return m.mock
}

//Set uses given function f as a mock of Calculator.GetCloudProof method
func (m *mCalculatorMockGetCloudProof) Set(f func(p context.Context, p1 *merkle.CloudEntry) (r merkle.OriginHash, r1 *merkle.CloudProof, r2 error)) *CalculatorMock {
	m.mock.GetCloudProofFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetCloudProof implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) GetCloudProof(p context.Context, p1 *merkle.CloudEntry) (r merkle.OriginHash, r1 *merkle.CloudProof, r2 error) {
	atomic.AddUint64(&m.GetCloudProofPreCounter, 1)
	defer atomic.AddUint64(&m.GetCloudProofCounter, 1)

	if m.GetCloudProofMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetCloudProofMock.mockExpectations, CalculatorMockGetCloudProofParams{p, p1},
			"Calculator.GetCloudProof got unexpected parameters")

		if m.GetCloudProofFunc == nil {

			m.t.Fatal("No results are set for the CalculatorMock.GetCloudProof")

			return
		}
	}

	if m.GetCloudProofFunc == nil {
		m.t.Fatal("Unexpected call to CalculatorMock.GetCloudProof")
		return
	}

	return m.GetCloudProofFunc(p, p1)
}

//GetCloudProofMinimockCounter returns a count of CalculatorMock.GetCloudProofFunc invocations
func (m *CalculatorMock) GetCloudProofMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudProofCounter)
}

//GetCloudProofMinimockPreCounter returns the value of CalculatorMock.GetCloudProof invocations
func (m *CalculatorMock) GetCloudProofMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudProofPreCounter)
}

type mCalculatorMockGetGlobuleProof struct {
	mock             *CalculatorMock
	mockExpectations *CalculatorMockGetGlobuleProofParams
}

//CalculatorMockGetGlobuleProofParams represents input parameters of the Calculator.GetGlobuleProof
type CalculatorMockGetGlobuleProofParams struct {
	p  context.Context
	p1 *merkle.GlobuleEntry
}

//Expect sets up expected params for the Calculator.GetGlobuleProof
func (m *mCalculatorMockGetGlobuleProof) Expect(p context.Context, p1 *merkle.GlobuleEntry) *mCalculatorMockGetGlobuleProof {
	m.mockExpectations = &CalculatorMockGetGlobuleProofParams{p, p1}
	return m
}

//Return sets up a mock for Calculator.GetGlobuleProof to return Return's arguments
func (m *mCalculatorMockGetGlobuleProof) Return(r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error) *CalculatorMock {
	m.mock.GetGlobuleProofFunc = func(p context.Context, p1 *merkle.GlobuleEntry) (merkle.OriginHash, *merkle.GlobuleProof, error) {
		return r, r1, r2
	}
	return m.mock
}

//Set uses given function f as a mock of Calculator.GetGlobuleProof method
func (m *mCalculatorMockGetGlobuleProof) Set(f func(p context.Context, p1 *merkle.GlobuleEntry) (r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error)) *CalculatorMock {
	m.mock.GetGlobuleProofFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetGlobuleProof implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) GetGlobuleProof(p context.Context, p1 *merkle.GlobuleEntry) (r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error) {
	atomic.AddUint64(&m.GetGlobuleProofPreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobuleProofCounter, 1)

	if m.GetGlobuleProofMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetGlobuleProofMock.mockExpectations, CalculatorMockGetGlobuleProofParams{p, p1},
			"Calculator.GetGlobuleProof got unexpected parameters")

		if m.GetGlobuleProofFunc == nil {

			m.t.Fatal("No results are set for the CalculatorMock.GetGlobuleProof")

			return
		}
	}

	if m.GetGlobuleProofFunc == nil {
		m.t.Fatal("Unexpected call to CalculatorMock.GetGlobuleProof")
		return
	}

	return m.GetGlobuleProofFunc(p, p1)
}

//GetGlobuleProofMinimockCounter returns a count of CalculatorMock.GetGlobuleProofFunc invocations
func (m *CalculatorMock) GetGlobuleProofMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleProofCounter)
}

//GetGlobuleProofMinimockPreCounter returns the value of CalculatorMock.GetGlobuleProof invocations
func (m *CalculatorMock) GetGlobuleProofMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleProofPreCounter)
}

type mCalculatorMockGetPulseProof struct {
	mock             *CalculatorMock
	mockExpectations *CalculatorMockGetPulseProofParams
}

//CalculatorMockGetPulseProofParams represents input parameters of the Calculator.GetPulseProof
type CalculatorMockGetPulseProofParams struct {
	p  context.Context
	p1 *merkle.PulseEntry
}

//Expect sets up expected params for the Calculator.GetPulseProof
func (m *mCalculatorMockGetPulseProof) Expect(p context.Context, p1 *merkle.PulseEntry) *mCalculatorMockGetPulseProof {
	m.mockExpectations = &CalculatorMockGetPulseProofParams{p, p1}
	return m
}

//Return sets up a mock for Calculator.GetPulseProof to return Return's arguments
func (m *mCalculatorMockGetPulseProof) Return(r merkle.OriginHash, r1 *merkle.PulseProof, r2 error) *CalculatorMock {
	m.mock.GetPulseProofFunc = func(p context.Context, p1 *merkle.PulseEntry) (merkle.OriginHash, *merkle.PulseProof, error) {
		return r, r1, r2
	}
	return m.mock
}

//Set uses given function f as a mock of Calculator.GetPulseProof method
func (m *mCalculatorMockGetPulseProof) Set(f func(p context.Context, p1 *merkle.PulseEntry) (r merkle.OriginHash, r1 *merkle.PulseProof, r2 error)) *CalculatorMock {
	m.mock.GetPulseProofFunc = f
	m.mockExpectations = nil
	return m.mock
}

//GetPulseProof implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) GetPulseProof(p context.Context, p1 *merkle.PulseEntry) (r merkle.OriginHash, r1 *merkle.PulseProof, r2 error) {
	atomic.AddUint64(&m.GetPulseProofPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseProofCounter, 1)

	if m.GetPulseProofMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.GetPulseProofMock.mockExpectations, CalculatorMockGetPulseProofParams{p, p1},
			"Calculator.GetPulseProof got unexpected parameters")

		if m.GetPulseProofFunc == nil {

			m.t.Fatal("No results are set for the CalculatorMock.GetPulseProof")

			return
		}
	}

	if m.GetPulseProofFunc == nil {
		m.t.Fatal("Unexpected call to CalculatorMock.GetPulseProof")
		return
	}

	return m.GetPulseProofFunc(p, p1)
}

//GetPulseProofMinimockCounter returns a count of CalculatorMock.GetPulseProofFunc invocations
func (m *CalculatorMock) GetPulseProofMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseProofCounter)
}

//GetPulseProofMinimockPreCounter returns the value of CalculatorMock.GetPulseProof invocations
func (m *CalculatorMock) GetPulseProofMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseProofPreCounter)
}

type mCalculatorMockIsValid struct {
	mock             *CalculatorMock
	mockExpectations *CalculatorMockIsValidParams
}

//CalculatorMockIsValidParams represents input parameters of the Calculator.IsValid
type CalculatorMockIsValidParams struct {
	p  merkle.Proof
	p1 merkle.OriginHash
	p2 crypto.PublicKey
}

//Expect sets up expected params for the Calculator.IsValid
func (m *mCalculatorMockIsValid) Expect(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) *mCalculatorMockIsValid {
	m.mockExpectations = &CalculatorMockIsValidParams{p, p1, p2}
	return m
}

//Return sets up a mock for Calculator.IsValid to return Return's arguments
func (m *mCalculatorMockIsValid) Return(r bool) *CalculatorMock {
	m.mock.IsValidFunc = func(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Calculator.IsValid method
func (m *mCalculatorMockIsValid) Set(f func(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) (r bool)) *CalculatorMock {
	m.mock.IsValidFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IsValid implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) IsValid(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) (r bool) {
	atomic.AddUint64(&m.IsValidPreCounter, 1)
	defer atomic.AddUint64(&m.IsValidCounter, 1)

	if m.IsValidMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IsValidMock.mockExpectations, CalculatorMockIsValidParams{p, p1, p2},
			"Calculator.IsValid got unexpected parameters")

		if m.IsValidFunc == nil {

			m.t.Fatal("No results are set for the CalculatorMock.IsValid")

			return
		}
	}

	if m.IsValidFunc == nil {
		m.t.Fatal("Unexpected call to CalculatorMock.IsValid")
		return
	}

	return m.IsValidFunc(p, p1, p2)
}

//IsValidMinimockCounter returns a count of CalculatorMock.IsValidFunc invocations
func (m *CalculatorMock) IsValidMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidCounter)
}

//IsValidMinimockPreCounter returns the value of CalculatorMock.IsValid invocations
func (m *CalculatorMock) IsValidMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CalculatorMock) ValidateCallCounters() {

	if m.GetCloudProofFunc != nil && atomic.LoadUint64(&m.GetCloudProofCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.GetCloudProof")
	}

	if m.GetGlobuleProofFunc != nil && atomic.LoadUint64(&m.GetGlobuleProofCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.GetGlobuleProof")
	}

	if m.GetPulseProofFunc != nil && atomic.LoadUint64(&m.GetPulseProofCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.GetPulseProof")
	}

	if m.IsValidFunc != nil && atomic.LoadUint64(&m.IsValidCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.IsValid")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CalculatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CalculatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CalculatorMock) MinimockFinish() {

	if m.GetCloudProofFunc != nil && atomic.LoadUint64(&m.GetCloudProofCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.GetCloudProof")
	}

	if m.GetGlobuleProofFunc != nil && atomic.LoadUint64(&m.GetGlobuleProofCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.GetGlobuleProof")
	}

	if m.GetPulseProofFunc != nil && atomic.LoadUint64(&m.GetPulseProofCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.GetPulseProof")
	}

	if m.IsValidFunc != nil && atomic.LoadUint64(&m.IsValidCounter) == 0 {
		m.t.Fatal("Expected call to CalculatorMock.IsValid")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CalculatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CalculatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetCloudProofFunc == nil || atomic.LoadUint64(&m.GetCloudProofCounter) > 0)
		ok = ok && (m.GetGlobuleProofFunc == nil || atomic.LoadUint64(&m.GetGlobuleProofCounter) > 0)
		ok = ok && (m.GetPulseProofFunc == nil || atomic.LoadUint64(&m.GetPulseProofCounter) > 0)
		ok = ok && (m.IsValidFunc == nil || atomic.LoadUint64(&m.IsValidCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetCloudProofFunc != nil && atomic.LoadUint64(&m.GetCloudProofCounter) == 0 {
				m.t.Error("Expected call to CalculatorMock.GetCloudProof")
			}

			if m.GetGlobuleProofFunc != nil && atomic.LoadUint64(&m.GetGlobuleProofCounter) == 0 {
				m.t.Error("Expected call to CalculatorMock.GetGlobuleProof")
			}

			if m.GetPulseProofFunc != nil && atomic.LoadUint64(&m.GetPulseProofCounter) == 0 {
				m.t.Error("Expected call to CalculatorMock.GetPulseProof")
			}

			if m.IsValidFunc != nil && atomic.LoadUint64(&m.IsValidCounter) == 0 {
				m.t.Error("Expected call to CalculatorMock.IsValid")
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
func (m *CalculatorMock) AllMocksCalled() bool {

	if m.GetCloudProofFunc != nil && atomic.LoadUint64(&m.GetCloudProofCounter) == 0 {
		return false
	}

	if m.GetGlobuleProofFunc != nil && atomic.LoadUint64(&m.GetGlobuleProofCounter) == 0 {
		return false
	}

	if m.GetPulseProofFunc != nil && atomic.LoadUint64(&m.GetPulseProofCounter) == 0 {
		return false
	}

	if m.IsValidFunc != nil && atomic.LoadUint64(&m.IsValidCounter) == 0 {
		return false
	}

	return true
}
