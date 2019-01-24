package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Parcel" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ParcelMock implements github.com/insolar/insolar/core.Parcel
type ParcelMock struct {
	t minimock.Tester

	AllowedSenderObjectAndRoleFunc       func() (r *core.RecordRef, r1 core.DynamicRole)
	AllowedSenderObjectAndRoleCounter    uint64
	AllowedSenderObjectAndRolePreCounter uint64
	AllowedSenderObjectAndRoleMock       mParcelMockAllowedSenderObjectAndRole

	ContextFunc       func(p context.Context) (r context.Context)
	ContextCounter    uint64
	ContextPreCounter uint64
	ContextMock       mParcelMockContext

	DefaultRoleFunc       func() (r core.DynamicRole)
	DefaultRoleCounter    uint64
	DefaultRolePreCounter uint64
	DefaultRoleMock       mParcelMockDefaultRole

	DefaultTargetFunc       func() (r *core.RecordRef)
	DefaultTargetCounter    uint64
	DefaultTargetPreCounter uint64
	DefaultTargetMock       mParcelMockDefaultTarget

	DelegationTokenFunc       func() (r core.DelegationToken)
	DelegationTokenCounter    uint64
	DelegationTokenPreCounter uint64
	DelegationTokenMock       mParcelMockDelegationToken

	GetCallerFunc       func() (r *core.RecordRef)
	GetCallerCounter    uint64
	GetCallerPreCounter uint64
	GetCallerMock       mParcelMockGetCaller

	GetSenderFunc       func() (r core.RecordRef)
	GetSenderCounter    uint64
	GetSenderPreCounter uint64
	GetSenderMock       mParcelMockGetSender

	GetSignFunc       func() (r []byte)
	GetSignCounter    uint64
	GetSignPreCounter uint64
	GetSignMock       mParcelMockGetSign

	MessageFunc       func() (r core.Message)
	MessageCounter    uint64
	MessagePreCounter uint64
	MessageMock       mParcelMockMessage

	PulseFunc       func() (r core.PulseNumber)
	PulseCounter    uint64
	PulsePreCounter uint64
	PulseMock       mParcelMockPulse

	TypeFunc       func() (r core.MessageType)
	TypeCounter    uint64
	TypePreCounter uint64
	TypeMock       mParcelMockType
}

//NewParcelMock returns a mock for github.com/insolar/insolar/core.Parcel
func NewParcelMock(t minimock.Tester) *ParcelMock {
	m := &ParcelMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AllowedSenderObjectAndRoleMock = mParcelMockAllowedSenderObjectAndRole{mock: m}
	m.ContextMock = mParcelMockContext{mock: m}
	m.DefaultRoleMock = mParcelMockDefaultRole{mock: m}
	m.DefaultTargetMock = mParcelMockDefaultTarget{mock: m}
	m.DelegationTokenMock = mParcelMockDelegationToken{mock: m}
	m.GetCallerMock = mParcelMockGetCaller{mock: m}
	m.GetSenderMock = mParcelMockGetSender{mock: m}
	m.GetSignMock = mParcelMockGetSign{mock: m}
	m.MessageMock = mParcelMockMessage{mock: m}
	m.PulseMock = mParcelMockPulse{mock: m}
	m.TypeMock = mParcelMockType{mock: m}

	return m
}

type mParcelMockAllowedSenderObjectAndRole struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockAllowedSenderObjectAndRoleExpectation
	expectationSeries []*ParcelMockAllowedSenderObjectAndRoleExpectation
}

type ParcelMockAllowedSenderObjectAndRoleExpectation struct {
	result *ParcelMockAllowedSenderObjectAndRoleResult
}

type ParcelMockAllowedSenderObjectAndRoleResult struct {
	r  *core.RecordRef
	r1 core.DynamicRole
}

//Expect specifies that invocation of Parcel.AllowedSenderObjectAndRole is expected from 1 to Infinity times
func (m *mParcelMockAllowedSenderObjectAndRole) Expect() *mParcelMockAllowedSenderObjectAndRole {
	m.mock.AllowedSenderObjectAndRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockAllowedSenderObjectAndRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.AllowedSenderObjectAndRole
func (m *mParcelMockAllowedSenderObjectAndRole) Return(r *core.RecordRef, r1 core.DynamicRole) *ParcelMock {
	m.mock.AllowedSenderObjectAndRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockAllowedSenderObjectAndRoleExpectation{}
	}
	m.mainExpectation.result = &ParcelMockAllowedSenderObjectAndRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.AllowedSenderObjectAndRole is expected once
func (m *mParcelMockAllowedSenderObjectAndRole) ExpectOnce() *ParcelMockAllowedSenderObjectAndRoleExpectation {
	m.mock.AllowedSenderObjectAndRoleFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockAllowedSenderObjectAndRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockAllowedSenderObjectAndRoleExpectation) Return(r *core.RecordRef, r1 core.DynamicRole) {
	e.result = &ParcelMockAllowedSenderObjectAndRoleResult{r, r1}
}

//Set uses given function f as a mock of Parcel.AllowedSenderObjectAndRole method
func (m *mParcelMockAllowedSenderObjectAndRole) Set(f func() (r *core.RecordRef, r1 core.DynamicRole)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AllowedSenderObjectAndRoleFunc = f
	return m.mock
}

//AllowedSenderObjectAndRole implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) AllowedSenderObjectAndRole() (r *core.RecordRef, r1 core.DynamicRole) {
	counter := atomic.AddUint64(&m.AllowedSenderObjectAndRolePreCounter, 1)
	defer atomic.AddUint64(&m.AllowedSenderObjectAndRoleCounter, 1)

	if len(m.AllowedSenderObjectAndRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AllowedSenderObjectAndRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.AllowedSenderObjectAndRole.")
			return
		}

		result := m.AllowedSenderObjectAndRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.AllowedSenderObjectAndRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AllowedSenderObjectAndRoleMock.mainExpectation != nil {

		result := m.AllowedSenderObjectAndRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.AllowedSenderObjectAndRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AllowedSenderObjectAndRoleFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.AllowedSenderObjectAndRole.")
		return
	}

	return m.AllowedSenderObjectAndRoleFunc()
}

//AllowedSenderObjectAndRoleMinimockCounter returns a count of ParcelMock.AllowedSenderObjectAndRoleFunc invocations
func (m *ParcelMock) AllowedSenderObjectAndRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AllowedSenderObjectAndRoleCounter)
}

//AllowedSenderObjectAndRoleMinimockPreCounter returns the value of ParcelMock.AllowedSenderObjectAndRole invocations
func (m *ParcelMock) AllowedSenderObjectAndRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AllowedSenderObjectAndRolePreCounter)
}

//AllowedSenderObjectAndRoleFinished returns true if mock invocations count is ok
func (m *ParcelMock) AllowedSenderObjectAndRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AllowedSenderObjectAndRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AllowedSenderObjectAndRoleCounter) == uint64(len(m.AllowedSenderObjectAndRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AllowedSenderObjectAndRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AllowedSenderObjectAndRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AllowedSenderObjectAndRoleFunc != nil {
		return atomic.LoadUint64(&m.AllowedSenderObjectAndRoleCounter) > 0
	}

	return true
}

type mParcelMockContext struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockContextExpectation
	expectationSeries []*ParcelMockContextExpectation
}

type ParcelMockContextExpectation struct {
	input  *ParcelMockContextInput
	result *ParcelMockContextResult
}

type ParcelMockContextInput struct {
	p context.Context
}

type ParcelMockContextResult struct {
	r context.Context
}

//Expect specifies that invocation of Parcel.Context is expected from 1 to Infinity times
func (m *mParcelMockContext) Expect(p context.Context) *mParcelMockContext {
	m.mock.ContextFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockContextExpectation{}
	}
	m.mainExpectation.input = &ParcelMockContextInput{p}
	return m
}

//Return specifies results of invocation of Parcel.Context
func (m *mParcelMockContext) Return(r context.Context) *ParcelMock {
	m.mock.ContextFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockContextExpectation{}
	}
	m.mainExpectation.result = &ParcelMockContextResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.Context is expected once
func (m *mParcelMockContext) ExpectOnce(p context.Context) *ParcelMockContextExpectation {
	m.mock.ContextFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockContextExpectation{}
	expectation.input = &ParcelMockContextInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockContextExpectation) Return(r context.Context) {
	e.result = &ParcelMockContextResult{r}
}

//Set uses given function f as a mock of Parcel.Context method
func (m *mParcelMockContext) Set(f func(p context.Context) (r context.Context)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ContextFunc = f
	return m.mock
}

//Context implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) Context(p context.Context) (r context.Context) {
	counter := atomic.AddUint64(&m.ContextPreCounter, 1)
	defer atomic.AddUint64(&m.ContextCounter, 1)

	if len(m.ContextMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ContextMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.Context. %v", p)
			return
		}

		input := m.ContextMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ParcelMockContextInput{p}, "Parcel.Context got unexpected parameters")

		result := m.ContextMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Context")
			return
		}

		r = result.r

		return
	}

	if m.ContextMock.mainExpectation != nil {

		input := m.ContextMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ParcelMockContextInput{p}, "Parcel.Context got unexpected parameters")
		}

		result := m.ContextMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Context")
		}

		r = result.r

		return
	}

	if m.ContextFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.Context. %v", p)
		return
	}

	return m.ContextFunc(p)
}

//ContextMinimockCounter returns a count of ParcelMock.ContextFunc invocations
func (m *ParcelMock) ContextMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ContextCounter)
}

//ContextMinimockPreCounter returns the value of ParcelMock.Context invocations
func (m *ParcelMock) ContextMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ContextPreCounter)
}

//ContextFinished returns true if mock invocations count is ok
func (m *ParcelMock) ContextFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ContextMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ContextCounter) == uint64(len(m.ContextMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ContextMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ContextCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ContextFunc != nil {
		return atomic.LoadUint64(&m.ContextCounter) > 0
	}

	return true
}

type mParcelMockDefaultRole struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockDefaultRoleExpectation
	expectationSeries []*ParcelMockDefaultRoleExpectation
}

type ParcelMockDefaultRoleExpectation struct {
	result *ParcelMockDefaultRoleResult
}

type ParcelMockDefaultRoleResult struct {
	r core.DynamicRole
}

//Expect specifies that invocation of Parcel.DefaultRole is expected from 1 to Infinity times
func (m *mParcelMockDefaultRole) Expect() *mParcelMockDefaultRole {
	m.mock.DefaultRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockDefaultRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.DefaultRole
func (m *mParcelMockDefaultRole) Return(r core.DynamicRole) *ParcelMock {
	m.mock.DefaultRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockDefaultRoleExpectation{}
	}
	m.mainExpectation.result = &ParcelMockDefaultRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.DefaultRole is expected once
func (m *mParcelMockDefaultRole) ExpectOnce() *ParcelMockDefaultRoleExpectation {
	m.mock.DefaultRoleFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockDefaultRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockDefaultRoleExpectation) Return(r core.DynamicRole) {
	e.result = &ParcelMockDefaultRoleResult{r}
}

//Set uses given function f as a mock of Parcel.DefaultRole method
func (m *mParcelMockDefaultRole) Set(f func() (r core.DynamicRole)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DefaultRoleFunc = f
	return m.mock
}

//DefaultRole implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) DefaultRole() (r core.DynamicRole) {
	counter := atomic.AddUint64(&m.DefaultRolePreCounter, 1)
	defer atomic.AddUint64(&m.DefaultRoleCounter, 1)

	if len(m.DefaultRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DefaultRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.DefaultRole.")
			return
		}

		result := m.DefaultRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.DefaultRole")
			return
		}

		r = result.r

		return
	}

	if m.DefaultRoleMock.mainExpectation != nil {

		result := m.DefaultRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.DefaultRole")
		}

		r = result.r

		return
	}

	if m.DefaultRoleFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.DefaultRole.")
		return
	}

	return m.DefaultRoleFunc()
}

//DefaultRoleMinimockCounter returns a count of ParcelMock.DefaultRoleFunc invocations
func (m *ParcelMock) DefaultRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DefaultRoleCounter)
}

//DefaultRoleMinimockPreCounter returns the value of ParcelMock.DefaultRole invocations
func (m *ParcelMock) DefaultRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DefaultRolePreCounter)
}

//DefaultRoleFinished returns true if mock invocations count is ok
func (m *ParcelMock) DefaultRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DefaultRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DefaultRoleCounter) == uint64(len(m.DefaultRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DefaultRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DefaultRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DefaultRoleFunc != nil {
		return atomic.LoadUint64(&m.DefaultRoleCounter) > 0
	}

	return true
}

type mParcelMockDefaultTarget struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockDefaultTargetExpectation
	expectationSeries []*ParcelMockDefaultTargetExpectation
}

type ParcelMockDefaultTargetExpectation struct {
	result *ParcelMockDefaultTargetResult
}

type ParcelMockDefaultTargetResult struct {
	r *core.RecordRef
}

//Expect specifies that invocation of Parcel.DefaultTarget is expected from 1 to Infinity times
func (m *mParcelMockDefaultTarget) Expect() *mParcelMockDefaultTarget {
	m.mock.DefaultTargetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockDefaultTargetExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.DefaultTarget
func (m *mParcelMockDefaultTarget) Return(r *core.RecordRef) *ParcelMock {
	m.mock.DefaultTargetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockDefaultTargetExpectation{}
	}
	m.mainExpectation.result = &ParcelMockDefaultTargetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.DefaultTarget is expected once
func (m *mParcelMockDefaultTarget) ExpectOnce() *ParcelMockDefaultTargetExpectation {
	m.mock.DefaultTargetFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockDefaultTargetExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockDefaultTargetExpectation) Return(r *core.RecordRef) {
	e.result = &ParcelMockDefaultTargetResult{r}
}

//Set uses given function f as a mock of Parcel.DefaultTarget method
func (m *mParcelMockDefaultTarget) Set(f func() (r *core.RecordRef)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DefaultTargetFunc = f
	return m.mock
}

//DefaultTarget implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) DefaultTarget() (r *core.RecordRef) {
	counter := atomic.AddUint64(&m.DefaultTargetPreCounter, 1)
	defer atomic.AddUint64(&m.DefaultTargetCounter, 1)

	if len(m.DefaultTargetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DefaultTargetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.DefaultTarget.")
			return
		}

		result := m.DefaultTargetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.DefaultTarget")
			return
		}

		r = result.r

		return
	}

	if m.DefaultTargetMock.mainExpectation != nil {

		result := m.DefaultTargetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.DefaultTarget")
		}

		r = result.r

		return
	}

	if m.DefaultTargetFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.DefaultTarget.")
		return
	}

	return m.DefaultTargetFunc()
}

//DefaultTargetMinimockCounter returns a count of ParcelMock.DefaultTargetFunc invocations
func (m *ParcelMock) DefaultTargetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DefaultTargetCounter)
}

//DefaultTargetMinimockPreCounter returns the value of ParcelMock.DefaultTarget invocations
func (m *ParcelMock) DefaultTargetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DefaultTargetPreCounter)
}

//DefaultTargetFinished returns true if mock invocations count is ok
func (m *ParcelMock) DefaultTargetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DefaultTargetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DefaultTargetCounter) == uint64(len(m.DefaultTargetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DefaultTargetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DefaultTargetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DefaultTargetFunc != nil {
		return atomic.LoadUint64(&m.DefaultTargetCounter) > 0
	}

	return true
}

type mParcelMockDelegationToken struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockDelegationTokenExpectation
	expectationSeries []*ParcelMockDelegationTokenExpectation
}

type ParcelMockDelegationTokenExpectation struct {
	result *ParcelMockDelegationTokenResult
}

type ParcelMockDelegationTokenResult struct {
	r core.DelegationToken
}

//Expect specifies that invocation of Parcel.DelegationToken is expected from 1 to Infinity times
func (m *mParcelMockDelegationToken) Expect() *mParcelMockDelegationToken {
	m.mock.DelegationTokenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockDelegationTokenExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.DelegationToken
func (m *mParcelMockDelegationToken) Return(r core.DelegationToken) *ParcelMock {
	m.mock.DelegationTokenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockDelegationTokenExpectation{}
	}
	m.mainExpectation.result = &ParcelMockDelegationTokenResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.DelegationToken is expected once
func (m *mParcelMockDelegationToken) ExpectOnce() *ParcelMockDelegationTokenExpectation {
	m.mock.DelegationTokenFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockDelegationTokenExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockDelegationTokenExpectation) Return(r core.DelegationToken) {
	e.result = &ParcelMockDelegationTokenResult{r}
}

//Set uses given function f as a mock of Parcel.DelegationToken method
func (m *mParcelMockDelegationToken) Set(f func() (r core.DelegationToken)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DelegationTokenFunc = f
	return m.mock
}

//DelegationToken implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) DelegationToken() (r core.DelegationToken) {
	counter := atomic.AddUint64(&m.DelegationTokenPreCounter, 1)
	defer atomic.AddUint64(&m.DelegationTokenCounter, 1)

	if len(m.DelegationTokenMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DelegationTokenMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.DelegationToken.")
			return
		}

		result := m.DelegationTokenMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.DelegationToken")
			return
		}

		r = result.r

		return
	}

	if m.DelegationTokenMock.mainExpectation != nil {

		result := m.DelegationTokenMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.DelegationToken")
		}

		r = result.r

		return
	}

	if m.DelegationTokenFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.DelegationToken.")
		return
	}

	return m.DelegationTokenFunc()
}

//DelegationTokenMinimockCounter returns a count of ParcelMock.DelegationTokenFunc invocations
func (m *ParcelMock) DelegationTokenMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DelegationTokenCounter)
}

//DelegationTokenMinimockPreCounter returns the value of ParcelMock.DelegationToken invocations
func (m *ParcelMock) DelegationTokenMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DelegationTokenPreCounter)
}

//DelegationTokenFinished returns true if mock invocations count is ok
func (m *ParcelMock) DelegationTokenFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DelegationTokenMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DelegationTokenCounter) == uint64(len(m.DelegationTokenMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DelegationTokenMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DelegationTokenCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DelegationTokenFunc != nil {
		return atomic.LoadUint64(&m.DelegationTokenCounter) > 0
	}

	return true
}

type mParcelMockGetCaller struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockGetCallerExpectation
	expectationSeries []*ParcelMockGetCallerExpectation
}

type ParcelMockGetCallerExpectation struct {
	result *ParcelMockGetCallerResult
}

type ParcelMockGetCallerResult struct {
	r *core.RecordRef
}

//Expect specifies that invocation of Parcel.GetCaller is expected from 1 to Infinity times
func (m *mParcelMockGetCaller) Expect() *mParcelMockGetCaller {
	m.mock.GetCallerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockGetCallerExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.GetCaller
func (m *mParcelMockGetCaller) Return(r *core.RecordRef) *ParcelMock {
	m.mock.GetCallerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockGetCallerExpectation{}
	}
	m.mainExpectation.result = &ParcelMockGetCallerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.GetCaller is expected once
func (m *mParcelMockGetCaller) ExpectOnce() *ParcelMockGetCallerExpectation {
	m.mock.GetCallerFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockGetCallerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockGetCallerExpectation) Return(r *core.RecordRef) {
	e.result = &ParcelMockGetCallerResult{r}
}

//Set uses given function f as a mock of Parcel.GetCaller method
func (m *mParcelMockGetCaller) Set(f func() (r *core.RecordRef)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCallerFunc = f
	return m.mock
}

//GetCaller implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) GetCaller() (r *core.RecordRef) {
	counter := atomic.AddUint64(&m.GetCallerPreCounter, 1)
	defer atomic.AddUint64(&m.GetCallerCounter, 1)

	if len(m.GetCallerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCallerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.GetCaller.")
			return
		}

		result := m.GetCallerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.GetCaller")
			return
		}

		r = result.r

		return
	}

	if m.GetCallerMock.mainExpectation != nil {

		result := m.GetCallerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.GetCaller")
		}

		r = result.r

		return
	}

	if m.GetCallerFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.GetCaller.")
		return
	}

	return m.GetCallerFunc()
}

//GetCallerMinimockCounter returns a count of ParcelMock.GetCallerFunc invocations
func (m *ParcelMock) GetCallerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCallerCounter)
}

//GetCallerMinimockPreCounter returns the value of ParcelMock.GetCaller invocations
func (m *ParcelMock) GetCallerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCallerPreCounter)
}

//GetCallerFinished returns true if mock invocations count is ok
func (m *ParcelMock) GetCallerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCallerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCallerCounter) == uint64(len(m.GetCallerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCallerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCallerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCallerFunc != nil {
		return atomic.LoadUint64(&m.GetCallerCounter) > 0
	}

	return true
}

type mParcelMockGetSender struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockGetSenderExpectation
	expectationSeries []*ParcelMockGetSenderExpectation
}

type ParcelMockGetSenderExpectation struct {
	result *ParcelMockGetSenderResult
}

type ParcelMockGetSenderResult struct {
	r core.RecordRef
}

//Expect specifies that invocation of Parcel.GetSender is expected from 1 to Infinity times
func (m *mParcelMockGetSender) Expect() *mParcelMockGetSender {
	m.mock.GetSenderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockGetSenderExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.GetSender
func (m *mParcelMockGetSender) Return(r core.RecordRef) *ParcelMock {
	m.mock.GetSenderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockGetSenderExpectation{}
	}
	m.mainExpectation.result = &ParcelMockGetSenderResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.GetSender is expected once
func (m *mParcelMockGetSender) ExpectOnce() *ParcelMockGetSenderExpectation {
	m.mock.GetSenderFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockGetSenderExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockGetSenderExpectation) Return(r core.RecordRef) {
	e.result = &ParcelMockGetSenderResult{r}
}

//Set uses given function f as a mock of Parcel.GetSender method
func (m *mParcelMockGetSender) Set(f func() (r core.RecordRef)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSenderFunc = f
	return m.mock
}

//GetSender implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) GetSender() (r core.RecordRef) {
	counter := atomic.AddUint64(&m.GetSenderPreCounter, 1)
	defer atomic.AddUint64(&m.GetSenderCounter, 1)

	if len(m.GetSenderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSenderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.GetSender.")
			return
		}

		result := m.GetSenderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.GetSender")
			return
		}

		r = result.r

		return
	}

	if m.GetSenderMock.mainExpectation != nil {

		result := m.GetSenderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.GetSender")
		}

		r = result.r

		return
	}

	if m.GetSenderFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.GetSender.")
		return
	}

	return m.GetSenderFunc()
}

//GetSenderMinimockCounter returns a count of ParcelMock.GetSenderFunc invocations
func (m *ParcelMock) GetSenderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSenderCounter)
}

//GetSenderMinimockPreCounter returns the value of ParcelMock.GetSender invocations
func (m *ParcelMock) GetSenderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSenderPreCounter)
}

//GetSenderFinished returns true if mock invocations count is ok
func (m *ParcelMock) GetSenderFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSenderMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSenderCounter) == uint64(len(m.GetSenderMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSenderMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSenderCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSenderFunc != nil {
		return atomic.LoadUint64(&m.GetSenderCounter) > 0
	}

	return true
}

type mParcelMockGetSign struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockGetSignExpectation
	expectationSeries []*ParcelMockGetSignExpectation
}

type ParcelMockGetSignExpectation struct {
	result *ParcelMockGetSignResult
}

type ParcelMockGetSignResult struct {
	r []byte
}

//Expect specifies that invocation of Parcel.GetSign is expected from 1 to Infinity times
func (m *mParcelMockGetSign) Expect() *mParcelMockGetSign {
	m.mock.GetSignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockGetSignExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.GetSign
func (m *mParcelMockGetSign) Return(r []byte) *ParcelMock {
	m.mock.GetSignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockGetSignExpectation{}
	}
	m.mainExpectation.result = &ParcelMockGetSignResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.GetSign is expected once
func (m *mParcelMockGetSign) ExpectOnce() *ParcelMockGetSignExpectation {
	m.mock.GetSignFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockGetSignExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockGetSignExpectation) Return(r []byte) {
	e.result = &ParcelMockGetSignResult{r}
}

//Set uses given function f as a mock of Parcel.GetSign method
func (m *mParcelMockGetSign) Set(f func() (r []byte)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignFunc = f
	return m.mock
}

//GetSign implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) GetSign() (r []byte) {
	counter := atomic.AddUint64(&m.GetSignPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignCounter, 1)

	if len(m.GetSignMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.GetSign.")
			return
		}

		result := m.GetSignMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.GetSign")
			return
		}

		r = result.r

		return
	}

	if m.GetSignMock.mainExpectation != nil {

		result := m.GetSignMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.GetSign")
		}

		r = result.r

		return
	}

	if m.GetSignFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.GetSign.")
		return
	}

	return m.GetSignFunc()
}

//GetSignMinimockCounter returns a count of ParcelMock.GetSignFunc invocations
func (m *ParcelMock) GetSignMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignCounter)
}

//GetSignMinimockPreCounter returns the value of ParcelMock.GetSign invocations
func (m *ParcelMock) GetSignMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignPreCounter)
}

//GetSignFinished returns true if mock invocations count is ok
func (m *ParcelMock) GetSignFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignCounter) == uint64(len(m.GetSignMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignFunc != nil {
		return atomic.LoadUint64(&m.GetSignCounter) > 0
	}

	return true
}

type mParcelMockMessage struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockMessageExpectation
	expectationSeries []*ParcelMockMessageExpectation
}

type ParcelMockMessageExpectation struct {
	result *ParcelMockMessageResult
}

type ParcelMockMessageResult struct {
	r core.Message
}

//Expect specifies that invocation of Parcel.Message is expected from 1 to Infinity times
func (m *mParcelMockMessage) Expect() *mParcelMockMessage {
	m.mock.MessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockMessageExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.Message
func (m *mParcelMockMessage) Return(r core.Message) *ParcelMock {
	m.mock.MessageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockMessageExpectation{}
	}
	m.mainExpectation.result = &ParcelMockMessageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.Message is expected once
func (m *mParcelMockMessage) ExpectOnce() *ParcelMockMessageExpectation {
	m.mock.MessageFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockMessageExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockMessageExpectation) Return(r core.Message) {
	e.result = &ParcelMockMessageResult{r}
}

//Set uses given function f as a mock of Parcel.Message method
func (m *mParcelMockMessage) Set(f func() (r core.Message)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MessageFunc = f
	return m.mock
}

//Message implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) Message() (r core.Message) {
	counter := atomic.AddUint64(&m.MessagePreCounter, 1)
	defer atomic.AddUint64(&m.MessageCounter, 1)

	if len(m.MessageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MessageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.Message.")
			return
		}

		result := m.MessageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Message")
			return
		}

		r = result.r

		return
	}

	if m.MessageMock.mainExpectation != nil {

		result := m.MessageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Message")
		}

		r = result.r

		return
	}

	if m.MessageFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.Message.")
		return
	}

	return m.MessageFunc()
}

//MessageMinimockCounter returns a count of ParcelMock.MessageFunc invocations
func (m *ParcelMock) MessageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MessageCounter)
}

//MessageMinimockPreCounter returns the value of ParcelMock.Message invocations
func (m *ParcelMock) MessageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MessagePreCounter)
}

//MessageFinished returns true if mock invocations count is ok
func (m *ParcelMock) MessageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MessageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MessageCounter) == uint64(len(m.MessageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MessageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MessageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MessageFunc != nil {
		return atomic.LoadUint64(&m.MessageCounter) > 0
	}

	return true
}

type mParcelMockPulse struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockPulseExpectation
	expectationSeries []*ParcelMockPulseExpectation
}

type ParcelMockPulseExpectation struct {
	result *ParcelMockPulseResult
}

type ParcelMockPulseResult struct {
	r core.PulseNumber
}

//Expect specifies that invocation of Parcel.Pulse is expected from 1 to Infinity times
func (m *mParcelMockPulse) Expect() *mParcelMockPulse {
	m.mock.PulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.Pulse
func (m *mParcelMockPulse) Return(r core.PulseNumber) *ParcelMock {
	m.mock.PulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockPulseExpectation{}
	}
	m.mainExpectation.result = &ParcelMockPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.Pulse is expected once
func (m *mParcelMockPulse) ExpectOnce() *ParcelMockPulseExpectation {
	m.mock.PulseFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockPulseExpectation) Return(r core.PulseNumber) {
	e.result = &ParcelMockPulseResult{r}
}

//Set uses given function f as a mock of Parcel.Pulse method
func (m *mParcelMockPulse) Set(f func() (r core.PulseNumber)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PulseFunc = f
	return m.mock
}

//Pulse implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) Pulse() (r core.PulseNumber) {
	counter := atomic.AddUint64(&m.PulsePreCounter, 1)
	defer atomic.AddUint64(&m.PulseCounter, 1)

	if len(m.PulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.Pulse.")
			return
		}

		result := m.PulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Pulse")
			return
		}

		r = result.r

		return
	}

	if m.PulseMock.mainExpectation != nil {

		result := m.PulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Pulse")
		}

		r = result.r

		return
	}

	if m.PulseFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.Pulse.")
		return
	}

	return m.PulseFunc()
}

//PulseMinimockCounter returns a count of ParcelMock.PulseFunc invocations
func (m *ParcelMock) PulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PulseCounter)
}

//PulseMinimockPreCounter returns the value of ParcelMock.Pulse invocations
func (m *ParcelMock) PulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PulsePreCounter)
}

//PulseFinished returns true if mock invocations count is ok
func (m *ParcelMock) PulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PulseCounter) == uint64(len(m.PulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PulseFunc != nil {
		return atomic.LoadUint64(&m.PulseCounter) > 0
	}

	return true
}

type mParcelMockType struct {
	mock              *ParcelMock
	mainExpectation   *ParcelMockTypeExpectation
	expectationSeries []*ParcelMockTypeExpectation
}

type ParcelMockTypeExpectation struct {
	result *ParcelMockTypeResult
}

type ParcelMockTypeResult struct {
	r core.MessageType
}

//Expect specifies that invocation of Parcel.Type is expected from 1 to Infinity times
func (m *mParcelMockType) Expect() *mParcelMockType {
	m.mock.TypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Parcel.Type
func (m *mParcelMockType) Return(r core.MessageType) *ParcelMock {
	m.mock.TypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ParcelMockTypeExpectation{}
	}
	m.mainExpectation.result = &ParcelMockTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Parcel.Type is expected once
func (m *mParcelMockType) ExpectOnce() *ParcelMockTypeExpectation {
	m.mock.TypeFunc = nil
	m.mainExpectation = nil

	expectation := &ParcelMockTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ParcelMockTypeExpectation) Return(r core.MessageType) {
	e.result = &ParcelMockTypeResult{r}
}

//Set uses given function f as a mock of Parcel.Type method
func (m *mParcelMockType) Set(f func() (r core.MessageType)) *ParcelMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.TypeFunc = f
	return m.mock
}

//Type implements github.com/insolar/insolar/core.Parcel interface
func (m *ParcelMock) Type() (r core.MessageType) {
	counter := atomic.AddUint64(&m.TypePreCounter, 1)
	defer atomic.AddUint64(&m.TypeCounter, 1)

	if len(m.TypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.TypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ParcelMock.Type.")
			return
		}

		result := m.TypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Type")
			return
		}

		r = result.r

		return
	}

	if m.TypeMock.mainExpectation != nil {

		result := m.TypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ParcelMock.Type")
		}

		r = result.r

		return
	}

	if m.TypeFunc == nil {
		m.t.Fatalf("Unexpected call to ParcelMock.Type.")
		return
	}

	return m.TypeFunc()
}

//TypeMinimockCounter returns a count of ParcelMock.TypeFunc invocations
func (m *ParcelMock) TypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.TypeCounter)
}

//TypeMinimockPreCounter returns the value of ParcelMock.Type invocations
func (m *ParcelMock) TypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.TypePreCounter)
}

//TypeFinished returns true if mock invocations count is ok
func (m *ParcelMock) TypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.TypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.TypeCounter) == uint64(len(m.TypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.TypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.TypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.TypeFunc != nil {
		return atomic.LoadUint64(&m.TypeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ParcelMock) ValidateCallCounters() {

	if !m.AllowedSenderObjectAndRoleFinished() {
		m.t.Fatal("Expected call to ParcelMock.AllowedSenderObjectAndRole")
	}

	if !m.ContextFinished() {
		m.t.Fatal("Expected call to ParcelMock.Context")
	}

	if !m.DefaultRoleFinished() {
		m.t.Fatal("Expected call to ParcelMock.DefaultRole")
	}

	if !m.DefaultTargetFinished() {
		m.t.Fatal("Expected call to ParcelMock.DefaultTarget")
	}

	if !m.DelegationTokenFinished() {
		m.t.Fatal("Expected call to ParcelMock.DelegationToken")
	}

	if !m.GetCallerFinished() {
		m.t.Fatal("Expected call to ParcelMock.GetCaller")
	}

	if !m.GetSenderFinished() {
		m.t.Fatal("Expected call to ParcelMock.GetSender")
	}

	if !m.GetSignFinished() {
		m.t.Fatal("Expected call to ParcelMock.GetSign")
	}

	if !m.MessageFinished() {
		m.t.Fatal("Expected call to ParcelMock.Message")
	}

	if !m.PulseFinished() {
		m.t.Fatal("Expected call to ParcelMock.Pulse")
	}

	if !m.TypeFinished() {
		m.t.Fatal("Expected call to ParcelMock.Type")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ParcelMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ParcelMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ParcelMock) MinimockFinish() {

	if !m.AllowedSenderObjectAndRoleFinished() {
		m.t.Fatal("Expected call to ParcelMock.AllowedSenderObjectAndRole")
	}

	if !m.ContextFinished() {
		m.t.Fatal("Expected call to ParcelMock.Context")
	}

	if !m.DefaultRoleFinished() {
		m.t.Fatal("Expected call to ParcelMock.DefaultRole")
	}

	if !m.DefaultTargetFinished() {
		m.t.Fatal("Expected call to ParcelMock.DefaultTarget")
	}

	if !m.DelegationTokenFinished() {
		m.t.Fatal("Expected call to ParcelMock.DelegationToken")
	}

	if !m.GetCallerFinished() {
		m.t.Fatal("Expected call to ParcelMock.GetCaller")
	}

	if !m.GetSenderFinished() {
		m.t.Fatal("Expected call to ParcelMock.GetSender")
	}

	if !m.GetSignFinished() {
		m.t.Fatal("Expected call to ParcelMock.GetSign")
	}

	if !m.MessageFinished() {
		m.t.Fatal("Expected call to ParcelMock.Message")
	}

	if !m.PulseFinished() {
		m.t.Fatal("Expected call to ParcelMock.Pulse")
	}

	if !m.TypeFinished() {
		m.t.Fatal("Expected call to ParcelMock.Type")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ParcelMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ParcelMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AllowedSenderObjectAndRoleFinished()
		ok = ok && m.ContextFinished()
		ok = ok && m.DefaultRoleFinished()
		ok = ok && m.DefaultTargetFinished()
		ok = ok && m.DelegationTokenFinished()
		ok = ok && m.GetCallerFinished()
		ok = ok && m.GetSenderFinished()
		ok = ok && m.GetSignFinished()
		ok = ok && m.MessageFinished()
		ok = ok && m.PulseFinished()
		ok = ok && m.TypeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AllowedSenderObjectAndRoleFinished() {
				m.t.Error("Expected call to ParcelMock.AllowedSenderObjectAndRole")
			}

			if !m.ContextFinished() {
				m.t.Error("Expected call to ParcelMock.Context")
			}

			if !m.DefaultRoleFinished() {
				m.t.Error("Expected call to ParcelMock.DefaultRole")
			}

			if !m.DefaultTargetFinished() {
				m.t.Error("Expected call to ParcelMock.DefaultTarget")
			}

			if !m.DelegationTokenFinished() {
				m.t.Error("Expected call to ParcelMock.DelegationToken")
			}

			if !m.GetCallerFinished() {
				m.t.Error("Expected call to ParcelMock.GetCaller")
			}

			if !m.GetSenderFinished() {
				m.t.Error("Expected call to ParcelMock.GetSender")
			}

			if !m.GetSignFinished() {
				m.t.Error("Expected call to ParcelMock.GetSign")
			}

			if !m.MessageFinished() {
				m.t.Error("Expected call to ParcelMock.Message")
			}

			if !m.PulseFinished() {
				m.t.Error("Expected call to ParcelMock.Pulse")
			}

			if !m.TypeFinished() {
				m.t.Error("Expected call to ParcelMock.Type")
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
func (m *ParcelMock) AllMocksCalled() bool {

	if !m.AllowedSenderObjectAndRoleFinished() {
		return false
	}

	if !m.ContextFinished() {
		return false
	}

	if !m.DefaultRoleFinished() {
		return false
	}

	if !m.DefaultTargetFinished() {
		return false
	}

	if !m.DelegationTokenFinished() {
		return false
	}

	if !m.GetCallerFinished() {
		return false
	}

	if !m.GetSenderFinished() {
		return false
	}

	if !m.GetSignFinished() {
		return false
	}

	if !m.MessageFinished() {
		return false
	}

	if !m.PulseFinished() {
		return false
	}

	if !m.TypeFinished() {
		return false
	}

	return true
}
