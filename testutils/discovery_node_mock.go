package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DiscoveryNode" can be found in github.com/insolar/insolar/insolar
*/
import (
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
)

//DiscoveryNodeMock implements github.com/insolar/insolar/insolar.DiscoveryNode
type DiscoveryNodeMock struct {
	t minimock.Tester

	GetBriefDigestFunc       func() (r []byte)
	GetBriefDigestCounter    uint64
	GetBriefDigestPreCounter uint64
	GetBriefDigestMock       mDiscoveryNodeMockGetBriefDigest

	GetBriefSignFunc       func() (r []byte)
	GetBriefSignCounter    uint64
	GetBriefSignPreCounter uint64
	GetBriefSignMock       mDiscoveryNodeMockGetBriefSign

	GetHostFunc       func() (r string)
	GetHostCounter    uint64
	GetHostPreCounter uint64
	GetHostMock       mDiscoveryNodeMockGetHost

	GetNodeRefFunc       func() (r *insolar.Reference)
	GetNodeRefCounter    uint64
	GetNodeRefPreCounter uint64
	GetNodeRefMock       mDiscoveryNodeMockGetNodeRef

	GetPublicKeyFunc       func() (r crypto.PublicKey)
	GetPublicKeyCounter    uint64
	GetPublicKeyPreCounter uint64
	GetPublicKeyMock       mDiscoveryNodeMockGetPublicKey

	GetRoleFunc       func() (r insolar.StaticRole)
	GetRoleCounter    uint64
	GetRolePreCounter uint64
	GetRoleMock       mDiscoveryNodeMockGetRole
}

//NewDiscoveryNodeMock returns a mock for github.com/insolar/insolar/insolar.DiscoveryNode
func NewDiscoveryNodeMock(t minimock.Tester) *DiscoveryNodeMock {
	m := &DiscoveryNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetBriefDigestMock = mDiscoveryNodeMockGetBriefDigest{mock: m}
	m.GetBriefSignMock = mDiscoveryNodeMockGetBriefSign{mock: m}
	m.GetHostMock = mDiscoveryNodeMockGetHost{mock: m}
	m.GetNodeRefMock = mDiscoveryNodeMockGetNodeRef{mock: m}
	m.GetPublicKeyMock = mDiscoveryNodeMockGetPublicKey{mock: m}
	m.GetRoleMock = mDiscoveryNodeMockGetRole{mock: m}

	return m
}

type mDiscoveryNodeMockGetBriefDigest struct {
	mock              *DiscoveryNodeMock
	mainExpectation   *DiscoveryNodeMockGetBriefDigestExpectation
	expectationSeries []*DiscoveryNodeMockGetBriefDigestExpectation
}

type DiscoveryNodeMockGetBriefDigestExpectation struct {
	result *DiscoveryNodeMockGetBriefDigestResult
}

type DiscoveryNodeMockGetBriefDigestResult struct {
	r []byte
}

//Expect specifies that invocation of DiscoveryNode.GetBriefDigest is expected from 1 to Infinity times
func (m *mDiscoveryNodeMockGetBriefDigest) Expect() *mDiscoveryNodeMockGetBriefDigest {
	m.mock.GetBriefDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetBriefDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of DiscoveryNode.GetBriefDigest
func (m *mDiscoveryNodeMockGetBriefDigest) Return(r []byte) *DiscoveryNodeMock {
	m.mock.GetBriefDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetBriefDigestExpectation{}
	}
	m.mainExpectation.result = &DiscoveryNodeMockGetBriefDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DiscoveryNode.GetBriefDigest is expected once
func (m *mDiscoveryNodeMockGetBriefDigest) ExpectOnce() *DiscoveryNodeMockGetBriefDigestExpectation {
	m.mock.GetBriefDigestFunc = nil
	m.mainExpectation = nil

	expectation := &DiscoveryNodeMockGetBriefDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DiscoveryNodeMockGetBriefDigestExpectation) Return(r []byte) {
	e.result = &DiscoveryNodeMockGetBriefDigestResult{r}
}

//Set uses given function f as a mock of DiscoveryNode.GetBriefDigest method
func (m *mDiscoveryNodeMockGetBriefDigest) Set(f func() (r []byte)) *DiscoveryNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetBriefDigestFunc = f
	return m.mock
}

//GetBriefDigest implements github.com/insolar/insolar/insolar.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetBriefDigest() (r []byte) {
	counter := atomic.AddUint64(&m.GetBriefDigestPreCounter, 1)
	defer atomic.AddUint64(&m.GetBriefDigestCounter, 1)

	if len(m.GetBriefDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetBriefDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetBriefDigest.")
			return
		}

		result := m.GetBriefDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetBriefDigest")
			return
		}

		r = result.r

		return
	}

	if m.GetBriefDigestMock.mainExpectation != nil {

		result := m.GetBriefDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetBriefDigest")
		}

		r = result.r

		return
	}

	if m.GetBriefDigestFunc == nil {
		m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetBriefDigest.")
		return
	}

	return m.GetBriefDigestFunc()
}

//GetBriefDigestMinimockCounter returns a count of DiscoveryNodeMock.GetBriefDigestFunc invocations
func (m *DiscoveryNodeMock) GetBriefDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefDigestCounter)
}

//GetBriefDigestMinimockPreCounter returns the value of DiscoveryNodeMock.GetBriefDigest invocations
func (m *DiscoveryNodeMock) GetBriefDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefDigestPreCounter)
}

//GetBriefDigestFinished returns true if mock invocations count is ok
func (m *DiscoveryNodeMock) GetBriefDigestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetBriefDigestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetBriefDigestCounter) == uint64(len(m.GetBriefDigestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetBriefDigestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetBriefDigestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetBriefDigestFunc != nil {
		return atomic.LoadUint64(&m.GetBriefDigestCounter) > 0
	}

	return true
}

type mDiscoveryNodeMockGetBriefSign struct {
	mock              *DiscoveryNodeMock
	mainExpectation   *DiscoveryNodeMockGetBriefSignExpectation
	expectationSeries []*DiscoveryNodeMockGetBriefSignExpectation
}

type DiscoveryNodeMockGetBriefSignExpectation struct {
	result *DiscoveryNodeMockGetBriefSignResult
}

type DiscoveryNodeMockGetBriefSignResult struct {
	r []byte
}

//Expect specifies that invocation of DiscoveryNode.GetBriefSign is expected from 1 to Infinity times
func (m *mDiscoveryNodeMockGetBriefSign) Expect() *mDiscoveryNodeMockGetBriefSign {
	m.mock.GetBriefSignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetBriefSignExpectation{}
	}

	return m
}

//Return specifies results of invocation of DiscoveryNode.GetBriefSign
func (m *mDiscoveryNodeMockGetBriefSign) Return(r []byte) *DiscoveryNodeMock {
	m.mock.GetBriefSignFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetBriefSignExpectation{}
	}
	m.mainExpectation.result = &DiscoveryNodeMockGetBriefSignResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DiscoveryNode.GetBriefSign is expected once
func (m *mDiscoveryNodeMockGetBriefSign) ExpectOnce() *DiscoveryNodeMockGetBriefSignExpectation {
	m.mock.GetBriefSignFunc = nil
	m.mainExpectation = nil

	expectation := &DiscoveryNodeMockGetBriefSignExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DiscoveryNodeMockGetBriefSignExpectation) Return(r []byte) {
	e.result = &DiscoveryNodeMockGetBriefSignResult{r}
}

//Set uses given function f as a mock of DiscoveryNode.GetBriefSign method
func (m *mDiscoveryNodeMockGetBriefSign) Set(f func() (r []byte)) *DiscoveryNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetBriefSignFunc = f
	return m.mock
}

//GetBriefSign implements github.com/insolar/insolar/insolar.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetBriefSign() (r []byte) {
	counter := atomic.AddUint64(&m.GetBriefSignPreCounter, 1)
	defer atomic.AddUint64(&m.GetBriefSignCounter, 1)

	if len(m.GetBriefSignMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetBriefSignMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetBriefSign.")
			return
		}

		result := m.GetBriefSignMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetBriefSign")
			return
		}

		r = result.r

		return
	}

	if m.GetBriefSignMock.mainExpectation != nil {

		result := m.GetBriefSignMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetBriefSign")
		}

		r = result.r

		return
	}

	if m.GetBriefSignFunc == nil {
		m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetBriefSign.")
		return
	}

	return m.GetBriefSignFunc()
}

//GetBriefSignMinimockCounter returns a count of DiscoveryNodeMock.GetBriefSignFunc invocations
func (m *DiscoveryNodeMock) GetBriefSignMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefSignCounter)
}

//GetBriefSignMinimockPreCounter returns the value of DiscoveryNodeMock.GetBriefSign invocations
func (m *DiscoveryNodeMock) GetBriefSignMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetBriefSignPreCounter)
}

//GetBriefSignFinished returns true if mock invocations count is ok
func (m *DiscoveryNodeMock) GetBriefSignFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetBriefSignMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetBriefSignCounter) == uint64(len(m.GetBriefSignMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetBriefSignMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetBriefSignCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetBriefSignFunc != nil {
		return atomic.LoadUint64(&m.GetBriefSignCounter) > 0
	}

	return true
}

type mDiscoveryNodeMockGetHost struct {
	mock              *DiscoveryNodeMock
	mainExpectation   *DiscoveryNodeMockGetHostExpectation
	expectationSeries []*DiscoveryNodeMockGetHostExpectation
}

type DiscoveryNodeMockGetHostExpectation struct {
	result *DiscoveryNodeMockGetHostResult
}

type DiscoveryNodeMockGetHostResult struct {
	r string
}

//Expect specifies that invocation of DiscoveryNode.GetHost is expected from 1 to Infinity times
func (m *mDiscoveryNodeMockGetHost) Expect() *mDiscoveryNodeMockGetHost {
	m.mock.GetHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetHostExpectation{}
	}

	return m
}

//Return specifies results of invocation of DiscoveryNode.GetHost
func (m *mDiscoveryNodeMockGetHost) Return(r string) *DiscoveryNodeMock {
	m.mock.GetHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetHostExpectation{}
	}
	m.mainExpectation.result = &DiscoveryNodeMockGetHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DiscoveryNode.GetHost is expected once
func (m *mDiscoveryNodeMockGetHost) ExpectOnce() *DiscoveryNodeMockGetHostExpectation {
	m.mock.GetHostFunc = nil
	m.mainExpectation = nil

	expectation := &DiscoveryNodeMockGetHostExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DiscoveryNodeMockGetHostExpectation) Return(r string) {
	e.result = &DiscoveryNodeMockGetHostResult{r}
}

//Set uses given function f as a mock of DiscoveryNode.GetHost method
func (m *mDiscoveryNodeMockGetHost) Set(f func() (r string)) *DiscoveryNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetHostFunc = f
	return m.mock
}

//GetHost implements github.com/insolar/insolar/insolar.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetHost() (r string) {
	counter := atomic.AddUint64(&m.GetHostPreCounter, 1)
	defer atomic.AddUint64(&m.GetHostCounter, 1)

	if len(m.GetHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetHost.")
			return
		}

		result := m.GetHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetHost")
			return
		}

		r = result.r

		return
	}

	if m.GetHostMock.mainExpectation != nil {

		result := m.GetHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetHost")
		}

		r = result.r

		return
	}

	if m.GetHostFunc == nil {
		m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetHost.")
		return
	}

	return m.GetHostFunc()
}

//GetHostMinimockCounter returns a count of DiscoveryNodeMock.GetHostFunc invocations
func (m *DiscoveryNodeMock) GetHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetHostCounter)
}

//GetHostMinimockPreCounter returns the value of DiscoveryNodeMock.GetHost invocations
func (m *DiscoveryNodeMock) GetHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetHostPreCounter)
}

//GetHostFinished returns true if mock invocations count is ok
func (m *DiscoveryNodeMock) GetHostFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetHostMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetHostCounter) == uint64(len(m.GetHostMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetHostMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetHostCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetHostFunc != nil {
		return atomic.LoadUint64(&m.GetHostCounter) > 0
	}

	return true
}

type mDiscoveryNodeMockGetNodeRef struct {
	mock              *DiscoveryNodeMock
	mainExpectation   *DiscoveryNodeMockGetNodeRefExpectation
	expectationSeries []*DiscoveryNodeMockGetNodeRefExpectation
}

type DiscoveryNodeMockGetNodeRefExpectation struct {
	result *DiscoveryNodeMockGetNodeRefResult
}

type DiscoveryNodeMockGetNodeRefResult struct {
	r *insolar.Reference
}

//Expect specifies that invocation of DiscoveryNode.GetNodeRef is expected from 1 to Infinity times
func (m *mDiscoveryNodeMockGetNodeRef) Expect() *mDiscoveryNodeMockGetNodeRef {
	m.mock.GetNodeRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetNodeRefExpectation{}
	}

	return m
}

//Return specifies results of invocation of DiscoveryNode.GetNodeRef
func (m *mDiscoveryNodeMockGetNodeRef) Return(r *insolar.Reference) *DiscoveryNodeMock {
	m.mock.GetNodeRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetNodeRefExpectation{}
	}
	m.mainExpectation.result = &DiscoveryNodeMockGetNodeRefResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DiscoveryNode.GetNodeRef is expected once
func (m *mDiscoveryNodeMockGetNodeRef) ExpectOnce() *DiscoveryNodeMockGetNodeRefExpectation {
	m.mock.GetNodeRefFunc = nil
	m.mainExpectation = nil

	expectation := &DiscoveryNodeMockGetNodeRefExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DiscoveryNodeMockGetNodeRefExpectation) Return(r *insolar.Reference) {
	e.result = &DiscoveryNodeMockGetNodeRefResult{r}
}

//Set uses given function f as a mock of DiscoveryNode.GetNodeRef method
func (m *mDiscoveryNodeMockGetNodeRef) Set(f func() (r *insolar.Reference)) *DiscoveryNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeRefFunc = f
	return m.mock
}

//GetNodeRef implements github.com/insolar/insolar/insolar.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetNodeRef() (r *insolar.Reference) {
	counter := atomic.AddUint64(&m.GetNodeRefPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeRefCounter, 1)

	if len(m.GetNodeRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetNodeRef.")
			return
		}

		result := m.GetNodeRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetNodeRef")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeRefMock.mainExpectation != nil {

		result := m.GetNodeRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetNodeRef")
		}

		r = result.r

		return
	}

	if m.GetNodeRefFunc == nil {
		m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetNodeRef.")
		return
	}

	return m.GetNodeRefFunc()
}

//GetNodeRefMinimockCounter returns a count of DiscoveryNodeMock.GetNodeRefFunc invocations
func (m *DiscoveryNodeMock) GetNodeRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefCounter)
}

//GetNodeRefMinimockPreCounter returns the value of DiscoveryNodeMock.GetNodeRef invocations
func (m *DiscoveryNodeMock) GetNodeRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeRefPreCounter)
}

//GetNodeRefFinished returns true if mock invocations count is ok
func (m *DiscoveryNodeMock) GetNodeRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeRefCounter) == uint64(len(m.GetNodeRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeRefFunc != nil {
		return atomic.LoadUint64(&m.GetNodeRefCounter) > 0
	}

	return true
}

type mDiscoveryNodeMockGetPublicKey struct {
	mock              *DiscoveryNodeMock
	mainExpectation   *DiscoveryNodeMockGetPublicKeyExpectation
	expectationSeries []*DiscoveryNodeMockGetPublicKeyExpectation
}

type DiscoveryNodeMockGetPublicKeyExpectation struct {
	result *DiscoveryNodeMockGetPublicKeyResult
}

type DiscoveryNodeMockGetPublicKeyResult struct {
	r crypto.PublicKey
}

//Expect specifies that invocation of DiscoveryNode.GetPublicKey is expected from 1 to Infinity times
func (m *mDiscoveryNodeMockGetPublicKey) Expect() *mDiscoveryNodeMockGetPublicKey {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetPublicKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of DiscoveryNode.GetPublicKey
func (m *mDiscoveryNodeMockGetPublicKey) Return(r crypto.PublicKey) *DiscoveryNodeMock {
	m.mock.GetPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetPublicKeyExpectation{}
	}
	m.mainExpectation.result = &DiscoveryNodeMockGetPublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DiscoveryNode.GetPublicKey is expected once
func (m *mDiscoveryNodeMockGetPublicKey) ExpectOnce() *DiscoveryNodeMockGetPublicKeyExpectation {
	m.mock.GetPublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &DiscoveryNodeMockGetPublicKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DiscoveryNodeMockGetPublicKeyExpectation) Return(r crypto.PublicKey) {
	e.result = &DiscoveryNodeMockGetPublicKeyResult{r}
}

//Set uses given function f as a mock of DiscoveryNode.GetPublicKey method
func (m *mDiscoveryNodeMockGetPublicKey) Set(f func() (r crypto.PublicKey)) *DiscoveryNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPublicKeyFunc = f
	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/insolar.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetPublicKey() (r crypto.PublicKey) {
	counter := atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if len(m.GetPublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetPublicKey.")
			return
		}

		result := m.GetPublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetPublicKey")
			return
		}

		r = result.r

		return
	}

	if m.GetPublicKeyMock.mainExpectation != nil {

		result := m.GetPublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetPublicKey")
		}

		r = result.r

		return
	}

	if m.GetPublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetPublicKey.")
		return
	}

	return m.GetPublicKeyFunc()
}

//GetPublicKeyMinimockCounter returns a count of DiscoveryNodeMock.GetPublicKeyFunc invocations
func (m *DiscoveryNodeMock) GetPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyCounter)
}

//GetPublicKeyMinimockPreCounter returns the value of DiscoveryNodeMock.GetPublicKey invocations
func (m *DiscoveryNodeMock) GetPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyPreCounter)
}

//GetPublicKeyFinished returns true if mock invocations count is ok
func (m *DiscoveryNodeMock) GetPublicKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPublicKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPublicKeyCounter) == uint64(len(m.GetPublicKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPublicKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPublicKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPublicKeyFunc != nil {
		return atomic.LoadUint64(&m.GetPublicKeyCounter) > 0
	}

	return true
}

type mDiscoveryNodeMockGetRole struct {
	mock              *DiscoveryNodeMock
	mainExpectation   *DiscoveryNodeMockGetRoleExpectation
	expectationSeries []*DiscoveryNodeMockGetRoleExpectation
}

type DiscoveryNodeMockGetRoleExpectation struct {
	result *DiscoveryNodeMockGetRoleResult
}

type DiscoveryNodeMockGetRoleResult struct {
	r insolar.StaticRole
}

//Expect specifies that invocation of DiscoveryNode.GetRole is expected from 1 to Infinity times
func (m *mDiscoveryNodeMockGetRole) Expect() *mDiscoveryNodeMockGetRole {
	m.mock.GetRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of DiscoveryNode.GetRole
func (m *mDiscoveryNodeMockGetRole) Return(r insolar.StaticRole) *DiscoveryNodeMock {
	m.mock.GetRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DiscoveryNodeMockGetRoleExpectation{}
	}
	m.mainExpectation.result = &DiscoveryNodeMockGetRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DiscoveryNode.GetRole is expected once
func (m *mDiscoveryNodeMockGetRole) ExpectOnce() *DiscoveryNodeMockGetRoleExpectation {
	m.mock.GetRoleFunc = nil
	m.mainExpectation = nil

	expectation := &DiscoveryNodeMockGetRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DiscoveryNodeMockGetRoleExpectation) Return(r insolar.StaticRole) {
	e.result = &DiscoveryNodeMockGetRoleResult{r}
}

//Set uses given function f as a mock of DiscoveryNode.GetRole method
func (m *mDiscoveryNodeMockGetRole) Set(f func() (r insolar.StaticRole)) *DiscoveryNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRoleFunc = f
	return m.mock
}

//GetRole implements github.com/insolar/insolar/insolar.DiscoveryNode interface
func (m *DiscoveryNodeMock) GetRole() (r insolar.StaticRole) {
	counter := atomic.AddUint64(&m.GetRolePreCounter, 1)
	defer atomic.AddUint64(&m.GetRoleCounter, 1)

	if len(m.GetRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetRole.")
			return
		}

		result := m.GetRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetRole")
			return
		}

		r = result.r

		return
	}

	if m.GetRoleMock.mainExpectation != nil {

		result := m.GetRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DiscoveryNodeMock.GetRole")
		}

		r = result.r

		return
	}

	if m.GetRoleFunc == nil {
		m.t.Fatalf("Unexpected call to DiscoveryNodeMock.GetRole.")
		return
	}

	return m.GetRoleFunc()
}

//GetRoleMinimockCounter returns a count of DiscoveryNodeMock.GetRoleFunc invocations
func (m *DiscoveryNodeMock) GetRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRoleCounter)
}

//GetRoleMinimockPreCounter returns the value of DiscoveryNodeMock.GetRole invocations
func (m *DiscoveryNodeMock) GetRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRolePreCounter)
}

//GetRoleFinished returns true if mock invocations count is ok
func (m *DiscoveryNodeMock) GetRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRoleCounter) == uint64(len(m.GetRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRoleFunc != nil {
		return atomic.LoadUint64(&m.GetRoleCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DiscoveryNodeMock) ValidateCallCounters() {

	if !m.GetBriefDigestFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetBriefDigest")
	}

	if !m.GetBriefSignFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetBriefSign")
	}

	if !m.GetHostFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetHost")
	}

	if !m.GetNodeRefFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetNodeRef")
	}

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetPublicKey")
	}

	if !m.GetRoleFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetRole")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DiscoveryNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DiscoveryNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DiscoveryNodeMock) MinimockFinish() {

	if !m.GetBriefDigestFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetBriefDigest")
	}

	if !m.GetBriefSignFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetBriefSign")
	}

	if !m.GetHostFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetHost")
	}

	if !m.GetNodeRefFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetNodeRef")
	}

	if !m.GetPublicKeyFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetPublicKey")
	}

	if !m.GetRoleFinished() {
		m.t.Fatal("Expected call to DiscoveryNodeMock.GetRole")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DiscoveryNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DiscoveryNodeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetBriefDigestFinished()
		ok = ok && m.GetBriefSignFinished()
		ok = ok && m.GetHostFinished()
		ok = ok && m.GetNodeRefFinished()
		ok = ok && m.GetPublicKeyFinished()
		ok = ok && m.GetRoleFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetBriefDigestFinished() {
				m.t.Error("Expected call to DiscoveryNodeMock.GetBriefDigest")
			}

			if !m.GetBriefSignFinished() {
				m.t.Error("Expected call to DiscoveryNodeMock.GetBriefSign")
			}

			if !m.GetHostFinished() {
				m.t.Error("Expected call to DiscoveryNodeMock.GetHost")
			}

			if !m.GetNodeRefFinished() {
				m.t.Error("Expected call to DiscoveryNodeMock.GetNodeRef")
			}

			if !m.GetPublicKeyFinished() {
				m.t.Error("Expected call to DiscoveryNodeMock.GetPublicKey")
			}

			if !m.GetRoleFinished() {
				m.t.Error("Expected call to DiscoveryNodeMock.GetRole")
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
func (m *DiscoveryNodeMock) AllMocksCalled() bool {

	if !m.GetBriefDigestFinished() {
		return false
	}

	if !m.GetBriefSignFinished() {
		return false
	}

	if !m.GetHostFinished() {
		return false
	}

	if !m.GetNodeRefFinished() {
		return false
	}

	if !m.GetPublicKeyFinished() {
		return false
	}

	if !m.GetRoleFinished() {
		return false
	}

	return true
}
