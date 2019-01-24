package merkle

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Calculator" can be found in github.com/insolar/insolar/network/merkle
*/
import (
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

	GetCloudProofFunc       func(p *merkle.CloudEntry) (r merkle.OriginHash, r1 *merkle.CloudProof, r2 error)
	GetCloudProofCounter    uint64
	GetCloudProofPreCounter uint64
	GetCloudProofMock       mCalculatorMockGetCloudProof

	GetGlobuleProofFunc       func(p *merkle.GlobuleEntry) (r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error)
	GetGlobuleProofCounter    uint64
	GetGlobuleProofPreCounter uint64
	GetGlobuleProofMock       mCalculatorMockGetGlobuleProof

	GetPulseProofFunc       func(p *merkle.PulseEntry) (r merkle.OriginHash, r1 *merkle.PulseProof, r2 error)
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
	mock              *CalculatorMock
	mainExpectation   *CalculatorMockGetCloudProofExpectation
	expectationSeries []*CalculatorMockGetCloudProofExpectation
}

type CalculatorMockGetCloudProofExpectation struct {
	input  *CalculatorMockGetCloudProofInput
	result *CalculatorMockGetCloudProofResult
}

type CalculatorMockGetCloudProofInput struct {
	p *merkle.CloudEntry
}

type CalculatorMockGetCloudProofResult struct {
	r  merkle.OriginHash
	r1 *merkle.CloudProof
	r2 error
}

//Expect specifies that invocation of Calculator.GetCloudProof is expected from 1 to Infinity times
func (m *mCalculatorMockGetCloudProof) Expect(p *merkle.CloudEntry) *mCalculatorMockGetCloudProof {
	m.mock.GetCloudProofFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockGetCloudProofExpectation{}
	}
	m.mainExpectation.input = &CalculatorMockGetCloudProofInput{p}
	return m
}

//Return specifies results of invocation of Calculator.GetCloudProof
func (m *mCalculatorMockGetCloudProof) Return(r merkle.OriginHash, r1 *merkle.CloudProof, r2 error) *CalculatorMock {
	m.mock.GetCloudProofFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockGetCloudProofExpectation{}
	}
	m.mainExpectation.result = &CalculatorMockGetCloudProofResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of Calculator.GetCloudProof is expected once
func (m *mCalculatorMockGetCloudProof) ExpectOnce(p *merkle.CloudEntry) *CalculatorMockGetCloudProofExpectation {
	m.mock.GetCloudProofFunc = nil
	m.mainExpectation = nil

	expectation := &CalculatorMockGetCloudProofExpectation{}
	expectation.input = &CalculatorMockGetCloudProofInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CalculatorMockGetCloudProofExpectation) Return(r merkle.OriginHash, r1 *merkle.CloudProof, r2 error) {
	e.result = &CalculatorMockGetCloudProofResult{r, r1, r2}
}

//Set uses given function f as a mock of Calculator.GetCloudProof method
func (m *mCalculatorMockGetCloudProof) Set(f func(p *merkle.CloudEntry) (r merkle.OriginHash, r1 *merkle.CloudProof, r2 error)) *CalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCloudProofFunc = f
	return m.mock
}

//GetCloudProof implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) GetCloudProof(p *merkle.CloudEntry) (r merkle.OriginHash, r1 *merkle.CloudProof, r2 error) {
	counter := atomic.AddUint64(&m.GetCloudProofPreCounter, 1)
	defer atomic.AddUint64(&m.GetCloudProofCounter, 1)

	if len(m.GetCloudProofMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCloudProofMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CalculatorMock.GetCloudProof. %v", p)
			return
		}

		input := m.GetCloudProofMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CalculatorMockGetCloudProofInput{p}, "Calculator.GetCloudProof got unexpected parameters")

		result := m.GetCloudProofMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.GetCloudProof")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetCloudProofMock.mainExpectation != nil {

		input := m.GetCloudProofMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CalculatorMockGetCloudProofInput{p}, "Calculator.GetCloudProof got unexpected parameters")
		}

		result := m.GetCloudProofMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.GetCloudProof")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetCloudProofFunc == nil {
		m.t.Fatalf("Unexpected call to CalculatorMock.GetCloudProof. %v", p)
		return
	}

	return m.GetCloudProofFunc(p)
}

//GetCloudProofMinimockCounter returns a count of CalculatorMock.GetCloudProofFunc invocations
func (m *CalculatorMock) GetCloudProofMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudProofCounter)
}

//GetCloudProofMinimockPreCounter returns the value of CalculatorMock.GetCloudProof invocations
func (m *CalculatorMock) GetCloudProofMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCloudProofPreCounter)
}

//GetCloudProofFinished returns true if mock invocations count is ok
func (m *CalculatorMock) GetCloudProofFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCloudProofMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCloudProofCounter) == uint64(len(m.GetCloudProofMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCloudProofMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCloudProofCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCloudProofFunc != nil {
		return atomic.LoadUint64(&m.GetCloudProofCounter) > 0
	}

	return true
}

type mCalculatorMockGetGlobuleProof struct {
	mock              *CalculatorMock
	mainExpectation   *CalculatorMockGetGlobuleProofExpectation
	expectationSeries []*CalculatorMockGetGlobuleProofExpectation
}

type CalculatorMockGetGlobuleProofExpectation struct {
	input  *CalculatorMockGetGlobuleProofInput
	result *CalculatorMockGetGlobuleProofResult
}

type CalculatorMockGetGlobuleProofInput struct {
	p *merkle.GlobuleEntry
}

type CalculatorMockGetGlobuleProofResult struct {
	r  merkle.OriginHash
	r1 *merkle.GlobuleProof
	r2 error
}

//Expect specifies that invocation of Calculator.GetGlobuleProof is expected from 1 to Infinity times
func (m *mCalculatorMockGetGlobuleProof) Expect(p *merkle.GlobuleEntry) *mCalculatorMockGetGlobuleProof {
	m.mock.GetGlobuleProofFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockGetGlobuleProofExpectation{}
	}
	m.mainExpectation.input = &CalculatorMockGetGlobuleProofInput{p}
	return m
}

//Return specifies results of invocation of Calculator.GetGlobuleProof
func (m *mCalculatorMockGetGlobuleProof) Return(r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error) *CalculatorMock {
	m.mock.GetGlobuleProofFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockGetGlobuleProofExpectation{}
	}
	m.mainExpectation.result = &CalculatorMockGetGlobuleProofResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of Calculator.GetGlobuleProof is expected once
func (m *mCalculatorMockGetGlobuleProof) ExpectOnce(p *merkle.GlobuleEntry) *CalculatorMockGetGlobuleProofExpectation {
	m.mock.GetGlobuleProofFunc = nil
	m.mainExpectation = nil

	expectation := &CalculatorMockGetGlobuleProofExpectation{}
	expectation.input = &CalculatorMockGetGlobuleProofInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CalculatorMockGetGlobuleProofExpectation) Return(r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error) {
	e.result = &CalculatorMockGetGlobuleProofResult{r, r1, r2}
}

//Set uses given function f as a mock of Calculator.GetGlobuleProof method
func (m *mCalculatorMockGetGlobuleProof) Set(f func(p *merkle.GlobuleEntry) (r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error)) *CalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetGlobuleProofFunc = f
	return m.mock
}

//GetGlobuleProof implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) GetGlobuleProof(p *merkle.GlobuleEntry) (r merkle.OriginHash, r1 *merkle.GlobuleProof, r2 error) {
	counter := atomic.AddUint64(&m.GetGlobuleProofPreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobuleProofCounter, 1)

	if len(m.GetGlobuleProofMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetGlobuleProofMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CalculatorMock.GetGlobuleProof. %v", p)
			return
		}

		input := m.GetGlobuleProofMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CalculatorMockGetGlobuleProofInput{p}, "Calculator.GetGlobuleProof got unexpected parameters")

		result := m.GetGlobuleProofMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.GetGlobuleProof")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetGlobuleProofMock.mainExpectation != nil {

		input := m.GetGlobuleProofMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CalculatorMockGetGlobuleProofInput{p}, "Calculator.GetGlobuleProof got unexpected parameters")
		}

		result := m.GetGlobuleProofMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.GetGlobuleProof")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetGlobuleProofFunc == nil {
		m.t.Fatalf("Unexpected call to CalculatorMock.GetGlobuleProof. %v", p)
		return
	}

	return m.GetGlobuleProofFunc(p)
}

//GetGlobuleProofMinimockCounter returns a count of CalculatorMock.GetGlobuleProofFunc invocations
func (m *CalculatorMock) GetGlobuleProofMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleProofCounter)
}

//GetGlobuleProofMinimockPreCounter returns the value of CalculatorMock.GetGlobuleProof invocations
func (m *CalculatorMock) GetGlobuleProofMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobuleProofPreCounter)
}

//GetGlobuleProofFinished returns true if mock invocations count is ok
func (m *CalculatorMock) GetGlobuleProofFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetGlobuleProofMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetGlobuleProofCounter) == uint64(len(m.GetGlobuleProofMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetGlobuleProofMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetGlobuleProofCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetGlobuleProofFunc != nil {
		return atomic.LoadUint64(&m.GetGlobuleProofCounter) > 0
	}

	return true
}

type mCalculatorMockGetPulseProof struct {
	mock              *CalculatorMock
	mainExpectation   *CalculatorMockGetPulseProofExpectation
	expectationSeries []*CalculatorMockGetPulseProofExpectation
}

type CalculatorMockGetPulseProofExpectation struct {
	input  *CalculatorMockGetPulseProofInput
	result *CalculatorMockGetPulseProofResult
}

type CalculatorMockGetPulseProofInput struct {
	p *merkle.PulseEntry
}

type CalculatorMockGetPulseProofResult struct {
	r  merkle.OriginHash
	r1 *merkle.PulseProof
	r2 error
}

//Expect specifies that invocation of Calculator.GetPulseProof is expected from 1 to Infinity times
func (m *mCalculatorMockGetPulseProof) Expect(p *merkle.PulseEntry) *mCalculatorMockGetPulseProof {
	m.mock.GetPulseProofFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockGetPulseProofExpectation{}
	}
	m.mainExpectation.input = &CalculatorMockGetPulseProofInput{p}
	return m
}

//Return specifies results of invocation of Calculator.GetPulseProof
func (m *mCalculatorMockGetPulseProof) Return(r merkle.OriginHash, r1 *merkle.PulseProof, r2 error) *CalculatorMock {
	m.mock.GetPulseProofFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockGetPulseProofExpectation{}
	}
	m.mainExpectation.result = &CalculatorMockGetPulseProofResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of Calculator.GetPulseProof is expected once
func (m *mCalculatorMockGetPulseProof) ExpectOnce(p *merkle.PulseEntry) *CalculatorMockGetPulseProofExpectation {
	m.mock.GetPulseProofFunc = nil
	m.mainExpectation = nil

	expectation := &CalculatorMockGetPulseProofExpectation{}
	expectation.input = &CalculatorMockGetPulseProofInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CalculatorMockGetPulseProofExpectation) Return(r merkle.OriginHash, r1 *merkle.PulseProof, r2 error) {
	e.result = &CalculatorMockGetPulseProofResult{r, r1, r2}
}

//Set uses given function f as a mock of Calculator.GetPulseProof method
func (m *mCalculatorMockGetPulseProof) Set(f func(p *merkle.PulseEntry) (r merkle.OriginHash, r1 *merkle.PulseProof, r2 error)) *CalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseProofFunc = f
	return m.mock
}

//GetPulseProof implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) GetPulseProof(p *merkle.PulseEntry) (r merkle.OriginHash, r1 *merkle.PulseProof, r2 error) {
	counter := atomic.AddUint64(&m.GetPulseProofPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseProofCounter, 1)

	if len(m.GetPulseProofMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseProofMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CalculatorMock.GetPulseProof. %v", p)
			return
		}

		input := m.GetPulseProofMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CalculatorMockGetPulseProofInput{p}, "Calculator.GetPulseProof got unexpected parameters")

		result := m.GetPulseProofMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.GetPulseProof")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetPulseProofMock.mainExpectation != nil {

		input := m.GetPulseProofMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CalculatorMockGetPulseProofInput{p}, "Calculator.GetPulseProof got unexpected parameters")
		}

		result := m.GetPulseProofMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.GetPulseProof")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.GetPulseProofFunc == nil {
		m.t.Fatalf("Unexpected call to CalculatorMock.GetPulseProof. %v", p)
		return
	}

	return m.GetPulseProofFunc(p)
}

//GetPulseProofMinimockCounter returns a count of CalculatorMock.GetPulseProofFunc invocations
func (m *CalculatorMock) GetPulseProofMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseProofCounter)
}

//GetPulseProofMinimockPreCounter returns the value of CalculatorMock.GetPulseProof invocations
func (m *CalculatorMock) GetPulseProofMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseProofPreCounter)
}

//GetPulseProofFinished returns true if mock invocations count is ok
func (m *CalculatorMock) GetPulseProofFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulseProofMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulseProofCounter) == uint64(len(m.GetPulseProofMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulseProofMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulseProofCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulseProofFunc != nil {
		return atomic.LoadUint64(&m.GetPulseProofCounter) > 0
	}

	return true
}

type mCalculatorMockIsValid struct {
	mock              *CalculatorMock
	mainExpectation   *CalculatorMockIsValidExpectation
	expectationSeries []*CalculatorMockIsValidExpectation
}

type CalculatorMockIsValidExpectation struct {
	input  *CalculatorMockIsValidInput
	result *CalculatorMockIsValidResult
}

type CalculatorMockIsValidInput struct {
	p  merkle.Proof
	p1 merkle.OriginHash
	p2 crypto.PublicKey
}

type CalculatorMockIsValidResult struct {
	r bool
}

//Expect specifies that invocation of Calculator.IsValid is expected from 1 to Infinity times
func (m *mCalculatorMockIsValid) Expect(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) *mCalculatorMockIsValid {
	m.mock.IsValidFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockIsValidExpectation{}
	}
	m.mainExpectation.input = &CalculatorMockIsValidInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Calculator.IsValid
func (m *mCalculatorMockIsValid) Return(r bool) *CalculatorMock {
	m.mock.IsValidFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CalculatorMockIsValidExpectation{}
	}
	m.mainExpectation.result = &CalculatorMockIsValidResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Calculator.IsValid is expected once
func (m *mCalculatorMockIsValid) ExpectOnce(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) *CalculatorMockIsValidExpectation {
	m.mock.IsValidFunc = nil
	m.mainExpectation = nil

	expectation := &CalculatorMockIsValidExpectation{}
	expectation.input = &CalculatorMockIsValidInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CalculatorMockIsValidExpectation) Return(r bool) {
	e.result = &CalculatorMockIsValidResult{r}
}

//Set uses given function f as a mock of Calculator.IsValid method
func (m *mCalculatorMockIsValid) Set(f func(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) (r bool)) *CalculatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsValidFunc = f
	return m.mock
}

//IsValid implements github.com/insolar/insolar/network/merkle.Calculator interface
func (m *CalculatorMock) IsValid(p merkle.Proof, p1 merkle.OriginHash, p2 crypto.PublicKey) (r bool) {
	counter := atomic.AddUint64(&m.IsValidPreCounter, 1)
	defer atomic.AddUint64(&m.IsValidCounter, 1)

	if len(m.IsValidMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsValidMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CalculatorMock.IsValid. %v %v %v", p, p1, p2)
			return
		}

		input := m.IsValidMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CalculatorMockIsValidInput{p, p1, p2}, "Calculator.IsValid got unexpected parameters")

		result := m.IsValidMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.IsValid")
			return
		}

		r = result.r

		return
	}

	if m.IsValidMock.mainExpectation != nil {

		input := m.IsValidMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CalculatorMockIsValidInput{p, p1, p2}, "Calculator.IsValid got unexpected parameters")
		}

		result := m.IsValidMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CalculatorMock.IsValid")
		}

		r = result.r

		return
	}

	if m.IsValidFunc == nil {
		m.t.Fatalf("Unexpected call to CalculatorMock.IsValid. %v %v %v", p, p1, p2)
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

//IsValidFinished returns true if mock invocations count is ok
func (m *CalculatorMock) IsValidFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsValidMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsValidCounter) == uint64(len(m.IsValidMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsValidMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsValidCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsValidFunc != nil {
		return atomic.LoadUint64(&m.IsValidCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CalculatorMock) ValidateCallCounters() {

	if !m.GetCloudProofFinished() {
		m.t.Fatal("Expected call to CalculatorMock.GetCloudProof")
	}

	if !m.GetGlobuleProofFinished() {
		m.t.Fatal("Expected call to CalculatorMock.GetGlobuleProof")
	}

	if !m.GetPulseProofFinished() {
		m.t.Fatal("Expected call to CalculatorMock.GetPulseProof")
	}

	if !m.IsValidFinished() {
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

	if !m.GetCloudProofFinished() {
		m.t.Fatal("Expected call to CalculatorMock.GetCloudProof")
	}

	if !m.GetGlobuleProofFinished() {
		m.t.Fatal("Expected call to CalculatorMock.GetGlobuleProof")
	}

	if !m.GetPulseProofFinished() {
		m.t.Fatal("Expected call to CalculatorMock.GetPulseProof")
	}

	if !m.IsValidFinished() {
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
		ok = ok && m.GetCloudProofFinished()
		ok = ok && m.GetGlobuleProofFinished()
		ok = ok && m.GetPulseProofFinished()
		ok = ok && m.IsValidFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetCloudProofFinished() {
				m.t.Error("Expected call to CalculatorMock.GetCloudProof")
			}

			if !m.GetGlobuleProofFinished() {
				m.t.Error("Expected call to CalculatorMock.GetGlobuleProof")
			}

			if !m.GetPulseProofFinished() {
				m.t.Error("Expected call to CalculatorMock.GetPulseProof")
			}

			if !m.IsValidFinished() {
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

	if !m.GetCloudProofFinished() {
		return false
	}

	if !m.GetGlobuleProofFinished() {
		return false
	}

	if !m.GetPulseProofFinished() {
		return false
	}

	if !m.IsValidFinished() {
		return false
	}

	return true
}
