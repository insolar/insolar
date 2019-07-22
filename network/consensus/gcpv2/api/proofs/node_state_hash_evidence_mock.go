package proofs

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeStateHashEvidence" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/proofs
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeStateHashEvidenceMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence
type NodeStateHashEvidenceMock struct {
	t minimock.Tester

	CopyOfSignedDigestFunc       func() (r cryptkit.SignedDigest)
	CopyOfSignedDigestCounter    uint64
	CopyOfSignedDigestPreCounter uint64
	CopyOfSignedDigestMock       mNodeStateHashEvidenceMockCopyOfSignedDigest

	EqualsFunc       func(p cryptkit.SignedDigestHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mNodeStateHashEvidenceMockEquals

	GetDigestHolderFunc       func() (r cryptkit.DigestHolder)
	GetDigestHolderCounter    uint64
	GetDigestHolderPreCounter uint64
	GetDigestHolderMock       mNodeStateHashEvidenceMockGetDigestHolder

	GetSignatureHolderFunc       func() (r cryptkit.SignatureHolder)
	GetSignatureHolderCounter    uint64
	GetSignatureHolderPreCounter uint64
	GetSignatureHolderMock       mNodeStateHashEvidenceMockGetSignatureHolder

	GetSignatureMethodFunc       func() (r cryptkit.SignatureMethod)
	GetSignatureMethodCounter    uint64
	GetSignatureMethodPreCounter uint64
	GetSignatureMethodMock       mNodeStateHashEvidenceMockGetSignatureMethod

	IsVerifiableByFunc       func(p cryptkit.SignatureVerifier) (r bool)
	IsVerifiableByCounter    uint64
	IsVerifiableByPreCounter uint64
	IsVerifiableByMock       mNodeStateHashEvidenceMockIsVerifiableBy

	VerifyWithFunc       func(p cryptkit.SignatureVerifier) (r bool)
	VerifyWithCounter    uint64
	VerifyWithPreCounter uint64
	VerifyWithMock       mNodeStateHashEvidenceMockVerifyWith
}

//NewNodeStateHashEvidenceMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence
func NewNodeStateHashEvidenceMock(t minimock.Tester) *NodeStateHashEvidenceMock {
	m := &NodeStateHashEvidenceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CopyOfSignedDigestMock = mNodeStateHashEvidenceMockCopyOfSignedDigest{mock: m}
	m.EqualsMock = mNodeStateHashEvidenceMockEquals{mock: m}
	m.GetDigestHolderMock = mNodeStateHashEvidenceMockGetDigestHolder{mock: m}
	m.GetSignatureHolderMock = mNodeStateHashEvidenceMockGetSignatureHolder{mock: m}
	m.GetSignatureMethodMock = mNodeStateHashEvidenceMockGetSignatureMethod{mock: m}
	m.IsVerifiableByMock = mNodeStateHashEvidenceMockIsVerifiableBy{mock: m}
	m.VerifyWithMock = mNodeStateHashEvidenceMockVerifyWith{mock: m}

	return m
}

type mNodeStateHashEvidenceMockCopyOfSignedDigest struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockCopyOfSignedDigestExpectation
	expectationSeries []*NodeStateHashEvidenceMockCopyOfSignedDigestExpectation
}

type NodeStateHashEvidenceMockCopyOfSignedDigestExpectation struct {
	result *NodeStateHashEvidenceMockCopyOfSignedDigestResult
}

type NodeStateHashEvidenceMockCopyOfSignedDigestResult struct {
	r cryptkit.SignedDigest
}

//Expect specifies that invocation of NodeStateHashEvidence.CopyOfSignedDigest is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockCopyOfSignedDigest) Expect() *mNodeStateHashEvidenceMockCopyOfSignedDigest {
	m.mock.CopyOfSignedDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockCopyOfSignedDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.CopyOfSignedDigest
func (m *mNodeStateHashEvidenceMockCopyOfSignedDigest) Return(r cryptkit.SignedDigest) *NodeStateHashEvidenceMock {
	m.mock.CopyOfSignedDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockCopyOfSignedDigestExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockCopyOfSignedDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.CopyOfSignedDigest is expected once
func (m *mNodeStateHashEvidenceMockCopyOfSignedDigest) ExpectOnce() *NodeStateHashEvidenceMockCopyOfSignedDigestExpectation {
	m.mock.CopyOfSignedDigestFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockCopyOfSignedDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockCopyOfSignedDigestExpectation) Return(r cryptkit.SignedDigest) {
	e.result = &NodeStateHashEvidenceMockCopyOfSignedDigestResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.CopyOfSignedDigest method
func (m *mNodeStateHashEvidenceMockCopyOfSignedDigest) Set(f func() (r cryptkit.SignedDigest)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CopyOfSignedDigestFunc = f
	return m.mock
}

//CopyOfSignedDigest implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) CopyOfSignedDigest() (r cryptkit.SignedDigest) {
	counter := atomic.AddUint64(&m.CopyOfSignedDigestPreCounter, 1)
	defer atomic.AddUint64(&m.CopyOfSignedDigestCounter, 1)

	if len(m.CopyOfSignedDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CopyOfSignedDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.CopyOfSignedDigest.")
			return
		}

		result := m.CopyOfSignedDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.CopyOfSignedDigest")
			return
		}

		r = result.r

		return
	}

	if m.CopyOfSignedDigestMock.mainExpectation != nil {

		result := m.CopyOfSignedDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.CopyOfSignedDigest")
		}

		r = result.r

		return
	}

	if m.CopyOfSignedDigestFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.CopyOfSignedDigest.")
		return
	}

	return m.CopyOfSignedDigestFunc()
}

//CopyOfSignedDigestMinimockCounter returns a count of NodeStateHashEvidenceMock.CopyOfSignedDigestFunc invocations
func (m *NodeStateHashEvidenceMock) CopyOfSignedDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfSignedDigestCounter)
}

//CopyOfSignedDigestMinimockPreCounter returns the value of NodeStateHashEvidenceMock.CopyOfSignedDigest invocations
func (m *NodeStateHashEvidenceMock) CopyOfSignedDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfSignedDigestPreCounter)
}

//CopyOfSignedDigestFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) CopyOfSignedDigestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CopyOfSignedDigestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CopyOfSignedDigestCounter) == uint64(len(m.CopyOfSignedDigestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CopyOfSignedDigestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CopyOfSignedDigestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CopyOfSignedDigestFunc != nil {
		return atomic.LoadUint64(&m.CopyOfSignedDigestCounter) > 0
	}

	return true
}

type mNodeStateHashEvidenceMockEquals struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockEqualsExpectation
	expectationSeries []*NodeStateHashEvidenceMockEqualsExpectation
}

type NodeStateHashEvidenceMockEqualsExpectation struct {
	input  *NodeStateHashEvidenceMockEqualsInput
	result *NodeStateHashEvidenceMockEqualsResult
}

type NodeStateHashEvidenceMockEqualsInput struct {
	p cryptkit.SignedDigestHolder
}

type NodeStateHashEvidenceMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of NodeStateHashEvidence.Equals is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockEquals) Expect(p cryptkit.SignedDigestHolder) *mNodeStateHashEvidenceMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockEqualsExpectation{}
	}
	m.mainExpectation.input = &NodeStateHashEvidenceMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.Equals
func (m *mNodeStateHashEvidenceMockEquals) Return(r bool) *NodeStateHashEvidenceMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockEqualsExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.Equals is expected once
func (m *mNodeStateHashEvidenceMockEquals) ExpectOnce(p cryptkit.SignedDigestHolder) *NodeStateHashEvidenceMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockEqualsExpectation{}
	expectation.input = &NodeStateHashEvidenceMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockEqualsExpectation) Return(r bool) {
	e.result = &NodeStateHashEvidenceMockEqualsResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.Equals method
func (m *mNodeStateHashEvidenceMockEquals) Set(f func(p cryptkit.SignedDigestHolder) (r bool)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) Equals(p cryptkit.SignedDigestHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStateHashEvidenceMockEqualsInput{p}, "NodeStateHashEvidence.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStateHashEvidenceMockEqualsInput{p}, "NodeStateHashEvidence.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of NodeStateHashEvidenceMock.EqualsFunc invocations
func (m *NodeStateHashEvidenceMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of NodeStateHashEvidenceMock.Equals invocations
func (m *NodeStateHashEvidenceMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) EqualsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.EqualsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.EqualsCounter) == uint64(len(m.EqualsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.EqualsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.EqualsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.EqualsFunc != nil {
		return atomic.LoadUint64(&m.EqualsCounter) > 0
	}

	return true
}

type mNodeStateHashEvidenceMockGetDigestHolder struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockGetDigestHolderExpectation
	expectationSeries []*NodeStateHashEvidenceMockGetDigestHolderExpectation
}

type NodeStateHashEvidenceMockGetDigestHolderExpectation struct {
	result *NodeStateHashEvidenceMockGetDigestHolderResult
}

type NodeStateHashEvidenceMockGetDigestHolderResult struct {
	r cryptkit.DigestHolder
}

//Expect specifies that invocation of NodeStateHashEvidence.GetDigestHolder is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockGetDigestHolder) Expect() *mNodeStateHashEvidenceMockGetDigestHolder {
	m.mock.GetDigestHolderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetDigestHolderExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.GetDigestHolder
func (m *mNodeStateHashEvidenceMockGetDigestHolder) Return(r cryptkit.DigestHolder) *NodeStateHashEvidenceMock {
	m.mock.GetDigestHolderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetDigestHolderExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockGetDigestHolderResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.GetDigestHolder is expected once
func (m *mNodeStateHashEvidenceMockGetDigestHolder) ExpectOnce() *NodeStateHashEvidenceMockGetDigestHolderExpectation {
	m.mock.GetDigestHolderFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockGetDigestHolderExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockGetDigestHolderExpectation) Return(r cryptkit.DigestHolder) {
	e.result = &NodeStateHashEvidenceMockGetDigestHolderResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.GetDigestHolder method
func (m *mNodeStateHashEvidenceMockGetDigestHolder) Set(f func() (r cryptkit.DigestHolder)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDigestHolderFunc = f
	return m.mock
}

//GetDigestHolder implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) GetDigestHolder() (r cryptkit.DigestHolder) {
	counter := atomic.AddUint64(&m.GetDigestHolderPreCounter, 1)
	defer atomic.AddUint64(&m.GetDigestHolderCounter, 1)

	if len(m.GetDigestHolderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDigestHolderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetDigestHolder.")
			return
		}

		result := m.GetDigestHolderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetDigestHolder")
			return
		}

		r = result.r

		return
	}

	if m.GetDigestHolderMock.mainExpectation != nil {

		result := m.GetDigestHolderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetDigestHolder")
		}

		r = result.r

		return
	}

	if m.GetDigestHolderFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetDigestHolder.")
		return
	}

	return m.GetDigestHolderFunc()
}

//GetDigestHolderMinimockCounter returns a count of NodeStateHashEvidenceMock.GetDigestHolderFunc invocations
func (m *NodeStateHashEvidenceMock) GetDigestHolderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestHolderCounter)
}

//GetDigestHolderMinimockPreCounter returns the value of NodeStateHashEvidenceMock.GetDigestHolder invocations
func (m *NodeStateHashEvidenceMock) GetDigestHolderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestHolderPreCounter)
}

//GetDigestHolderFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) GetDigestHolderFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDigestHolderMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDigestHolderCounter) == uint64(len(m.GetDigestHolderMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDigestHolderMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDigestHolderCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDigestHolderFunc != nil {
		return atomic.LoadUint64(&m.GetDigestHolderCounter) > 0
	}

	return true
}

type mNodeStateHashEvidenceMockGetSignatureHolder struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockGetSignatureHolderExpectation
	expectationSeries []*NodeStateHashEvidenceMockGetSignatureHolderExpectation
}

type NodeStateHashEvidenceMockGetSignatureHolderExpectation struct {
	result *NodeStateHashEvidenceMockGetSignatureHolderResult
}

type NodeStateHashEvidenceMockGetSignatureHolderResult struct {
	r cryptkit.SignatureHolder
}

//Expect specifies that invocation of NodeStateHashEvidence.GetSignatureHolder is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockGetSignatureHolder) Expect() *mNodeStateHashEvidenceMockGetSignatureHolder {
	m.mock.GetSignatureHolderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetSignatureHolderExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.GetSignatureHolder
func (m *mNodeStateHashEvidenceMockGetSignatureHolder) Return(r cryptkit.SignatureHolder) *NodeStateHashEvidenceMock {
	m.mock.GetSignatureHolderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetSignatureHolderExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockGetSignatureHolderResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.GetSignatureHolder is expected once
func (m *mNodeStateHashEvidenceMockGetSignatureHolder) ExpectOnce() *NodeStateHashEvidenceMockGetSignatureHolderExpectation {
	m.mock.GetSignatureHolderFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockGetSignatureHolderExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockGetSignatureHolderExpectation) Return(r cryptkit.SignatureHolder) {
	e.result = &NodeStateHashEvidenceMockGetSignatureHolderResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.GetSignatureHolder method
func (m *mNodeStateHashEvidenceMockGetSignatureHolder) Set(f func() (r cryptkit.SignatureHolder)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureHolderFunc = f
	return m.mock
}

//GetSignatureHolder implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) GetSignatureHolder() (r cryptkit.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetSignatureHolderPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureHolderCounter, 1)

	if len(m.GetSignatureHolderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureHolderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetSignatureHolder.")
			return
		}

		result := m.GetSignatureHolderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetSignatureHolder")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureHolderMock.mainExpectation != nil {

		result := m.GetSignatureHolderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetSignatureHolder")
		}

		r = result.r

		return
	}

	if m.GetSignatureHolderFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetSignatureHolder.")
		return
	}

	return m.GetSignatureHolderFunc()
}

//GetSignatureHolderMinimockCounter returns a count of NodeStateHashEvidenceMock.GetSignatureHolderFunc invocations
func (m *NodeStateHashEvidenceMock) GetSignatureHolderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureHolderCounter)
}

//GetSignatureHolderMinimockPreCounter returns the value of NodeStateHashEvidenceMock.GetSignatureHolder invocations
func (m *NodeStateHashEvidenceMock) GetSignatureHolderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureHolderPreCounter)
}

//GetSignatureHolderFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) GetSignatureHolderFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureHolderMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureHolderCounter) == uint64(len(m.GetSignatureHolderMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureHolderMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureHolderCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureHolderFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureHolderCounter) > 0
	}

	return true
}

type mNodeStateHashEvidenceMockGetSignatureMethod struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockGetSignatureMethodExpectation
	expectationSeries []*NodeStateHashEvidenceMockGetSignatureMethodExpectation
}

type NodeStateHashEvidenceMockGetSignatureMethodExpectation struct {
	result *NodeStateHashEvidenceMockGetSignatureMethodResult
}

type NodeStateHashEvidenceMockGetSignatureMethodResult struct {
	r cryptkit.SignatureMethod
}

//Expect specifies that invocation of NodeStateHashEvidence.GetSignatureMethod is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockGetSignatureMethod) Expect() *mNodeStateHashEvidenceMockGetSignatureMethod {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetSignatureMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.GetSignatureMethod
func (m *mNodeStateHashEvidenceMockGetSignatureMethod) Return(r cryptkit.SignatureMethod) *NodeStateHashEvidenceMock {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetSignatureMethodExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockGetSignatureMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.GetSignatureMethod is expected once
func (m *mNodeStateHashEvidenceMockGetSignatureMethod) ExpectOnce() *NodeStateHashEvidenceMockGetSignatureMethodExpectation {
	m.mock.GetSignatureMethodFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockGetSignatureMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockGetSignatureMethodExpectation) Return(r cryptkit.SignatureMethod) {
	e.result = &NodeStateHashEvidenceMockGetSignatureMethodResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.GetSignatureMethod method
func (m *mNodeStateHashEvidenceMockGetSignatureMethod) Set(f func() (r cryptkit.SignatureMethod)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureMethodFunc = f
	return m.mock
}

//GetSignatureMethod implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) GetSignatureMethod() (r cryptkit.SignatureMethod) {
	counter := atomic.AddUint64(&m.GetSignatureMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureMethodCounter, 1)

	if len(m.GetSignatureMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetSignatureMethod.")
			return
		}

		result := m.GetSignatureMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetSignatureMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodMock.mainExpectation != nil {

		result := m.GetSignatureMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetSignatureMethod")
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetSignatureMethod.")
		return
	}

	return m.GetSignatureMethodFunc()
}

//GetSignatureMethodMinimockCounter returns a count of NodeStateHashEvidenceMock.GetSignatureMethodFunc invocations
func (m *NodeStateHashEvidenceMock) GetSignatureMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodCounter)
}

//GetSignatureMethodMinimockPreCounter returns the value of NodeStateHashEvidenceMock.GetSignatureMethod invocations
func (m *NodeStateHashEvidenceMock) GetSignatureMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodPreCounter)
}

//GetSignatureMethodFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) GetSignatureMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureMethodCounter) == uint64(len(m.GetSignatureMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureMethodFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureMethodCounter) > 0
	}

	return true
}

type mNodeStateHashEvidenceMockIsVerifiableBy struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockIsVerifiableByExpectation
	expectationSeries []*NodeStateHashEvidenceMockIsVerifiableByExpectation
}

type NodeStateHashEvidenceMockIsVerifiableByExpectation struct {
	input  *NodeStateHashEvidenceMockIsVerifiableByInput
	result *NodeStateHashEvidenceMockIsVerifiableByResult
}

type NodeStateHashEvidenceMockIsVerifiableByInput struct {
	p cryptkit.SignatureVerifier
}

type NodeStateHashEvidenceMockIsVerifiableByResult struct {
	r bool
}

//Expect specifies that invocation of NodeStateHashEvidence.IsVerifiableBy is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockIsVerifiableBy) Expect(p cryptkit.SignatureVerifier) *mNodeStateHashEvidenceMockIsVerifiableBy {
	m.mock.IsVerifiableByFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockIsVerifiableByExpectation{}
	}
	m.mainExpectation.input = &NodeStateHashEvidenceMockIsVerifiableByInput{p}
	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.IsVerifiableBy
func (m *mNodeStateHashEvidenceMockIsVerifiableBy) Return(r bool) *NodeStateHashEvidenceMock {
	m.mock.IsVerifiableByFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockIsVerifiableByExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockIsVerifiableByResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.IsVerifiableBy is expected once
func (m *mNodeStateHashEvidenceMockIsVerifiableBy) ExpectOnce(p cryptkit.SignatureVerifier) *NodeStateHashEvidenceMockIsVerifiableByExpectation {
	m.mock.IsVerifiableByFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockIsVerifiableByExpectation{}
	expectation.input = &NodeStateHashEvidenceMockIsVerifiableByInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockIsVerifiableByExpectation) Return(r bool) {
	e.result = &NodeStateHashEvidenceMockIsVerifiableByResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.IsVerifiableBy method
func (m *mNodeStateHashEvidenceMockIsVerifiableBy) Set(f func(p cryptkit.SignatureVerifier) (r bool)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsVerifiableByFunc = f
	return m.mock
}

//IsVerifiableBy implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) IsVerifiableBy(p cryptkit.SignatureVerifier) (r bool) {
	counter := atomic.AddUint64(&m.IsVerifiableByPreCounter, 1)
	defer atomic.AddUint64(&m.IsVerifiableByCounter, 1)

	if len(m.IsVerifiableByMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsVerifiableByMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.IsVerifiableBy. %v", p)
			return
		}

		input := m.IsVerifiableByMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStateHashEvidenceMockIsVerifiableByInput{p}, "NodeStateHashEvidence.IsVerifiableBy got unexpected parameters")

		result := m.IsVerifiableByMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.IsVerifiableBy")
			return
		}

		r = result.r

		return
	}

	if m.IsVerifiableByMock.mainExpectation != nil {

		input := m.IsVerifiableByMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStateHashEvidenceMockIsVerifiableByInput{p}, "NodeStateHashEvidence.IsVerifiableBy got unexpected parameters")
		}

		result := m.IsVerifiableByMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.IsVerifiableBy")
		}

		r = result.r

		return
	}

	if m.IsVerifiableByFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.IsVerifiableBy. %v", p)
		return
	}

	return m.IsVerifiableByFunc(p)
}

//IsVerifiableByMinimockCounter returns a count of NodeStateHashEvidenceMock.IsVerifiableByFunc invocations
func (m *NodeStateHashEvidenceMock) IsVerifiableByMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsVerifiableByCounter)
}

//IsVerifiableByMinimockPreCounter returns the value of NodeStateHashEvidenceMock.IsVerifiableBy invocations
func (m *NodeStateHashEvidenceMock) IsVerifiableByMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsVerifiableByPreCounter)
}

//IsVerifiableByFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) IsVerifiableByFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsVerifiableByMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsVerifiableByCounter) == uint64(len(m.IsVerifiableByMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsVerifiableByMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsVerifiableByCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsVerifiableByFunc != nil {
		return atomic.LoadUint64(&m.IsVerifiableByCounter) > 0
	}

	return true
}

type mNodeStateHashEvidenceMockVerifyWith struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockVerifyWithExpectation
	expectationSeries []*NodeStateHashEvidenceMockVerifyWithExpectation
}

type NodeStateHashEvidenceMockVerifyWithExpectation struct {
	input  *NodeStateHashEvidenceMockVerifyWithInput
	result *NodeStateHashEvidenceMockVerifyWithResult
}

type NodeStateHashEvidenceMockVerifyWithInput struct {
	p cryptkit.SignatureVerifier
}

type NodeStateHashEvidenceMockVerifyWithResult struct {
	r bool
}

//Expect specifies that invocation of NodeStateHashEvidence.VerifyWith is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockVerifyWith) Expect(p cryptkit.SignatureVerifier) *mNodeStateHashEvidenceMockVerifyWith {
	m.mock.VerifyWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockVerifyWithExpectation{}
	}
	m.mainExpectation.input = &NodeStateHashEvidenceMockVerifyWithInput{p}
	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.VerifyWith
func (m *mNodeStateHashEvidenceMockVerifyWith) Return(r bool) *NodeStateHashEvidenceMock {
	m.mock.VerifyWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockVerifyWithExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockVerifyWithResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.VerifyWith is expected once
func (m *mNodeStateHashEvidenceMockVerifyWith) ExpectOnce(p cryptkit.SignatureVerifier) *NodeStateHashEvidenceMockVerifyWithExpectation {
	m.mock.VerifyWithFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockVerifyWithExpectation{}
	expectation.input = &NodeStateHashEvidenceMockVerifyWithInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockVerifyWithExpectation) Return(r bool) {
	e.result = &NodeStateHashEvidenceMockVerifyWithResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.VerifyWith method
func (m *mNodeStateHashEvidenceMockVerifyWith) Set(f func(p cryptkit.SignatureVerifier) (r bool)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.VerifyWithFunc = f
	return m.mock
}

//VerifyWith implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) VerifyWith(p cryptkit.SignatureVerifier) (r bool) {
	counter := atomic.AddUint64(&m.VerifyWithPreCounter, 1)
	defer atomic.AddUint64(&m.VerifyWithCounter, 1)

	if len(m.VerifyWithMock.expectationSeries) > 0 {
		if counter > uint64(len(m.VerifyWithMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.VerifyWith. %v", p)
			return
		}

		input := m.VerifyWithMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStateHashEvidenceMockVerifyWithInput{p}, "NodeStateHashEvidence.VerifyWith got unexpected parameters")

		result := m.VerifyWithMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.VerifyWith")
			return
		}

		r = result.r

		return
	}

	if m.VerifyWithMock.mainExpectation != nil {

		input := m.VerifyWithMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStateHashEvidenceMockVerifyWithInput{p}, "NodeStateHashEvidence.VerifyWith got unexpected parameters")
		}

		result := m.VerifyWithMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.VerifyWith")
		}

		r = result.r

		return
	}

	if m.VerifyWithFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.VerifyWith. %v", p)
		return
	}

	return m.VerifyWithFunc(p)
}

//VerifyWithMinimockCounter returns a count of NodeStateHashEvidenceMock.VerifyWithFunc invocations
func (m *NodeStateHashEvidenceMock) VerifyWithMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyWithCounter)
}

//VerifyWithMinimockPreCounter returns the value of NodeStateHashEvidenceMock.VerifyWith invocations
func (m *NodeStateHashEvidenceMock) VerifyWithMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyWithPreCounter)
}

//VerifyWithFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) VerifyWithFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.VerifyWithMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.VerifyWithCounter) == uint64(len(m.VerifyWithMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.VerifyWithMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.VerifyWithCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.VerifyWithFunc != nil {
		return atomic.LoadUint64(&m.VerifyWithCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeStateHashEvidenceMock) ValidateCallCounters() {

	if !m.CopyOfSignedDigestFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.CopyOfSignedDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.Equals")
	}

	if !m.GetDigestHolderFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetDigestHolder")
	}

	if !m.GetSignatureHolderFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetSignatureHolder")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetSignatureMethod")
	}

	if !m.IsVerifiableByFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.IsVerifiableBy")
	}

	if !m.VerifyWithFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.VerifyWith")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeStateHashEvidenceMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeStateHashEvidenceMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeStateHashEvidenceMock) MinimockFinish() {

	if !m.CopyOfSignedDigestFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.CopyOfSignedDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.Equals")
	}

	if !m.GetDigestHolderFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetDigestHolder")
	}

	if !m.GetSignatureHolderFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetSignatureHolder")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetSignatureMethod")
	}

	if !m.IsVerifiableByFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.IsVerifiableBy")
	}

	if !m.VerifyWithFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.VerifyWith")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeStateHashEvidenceMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeStateHashEvidenceMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CopyOfSignedDigestFinished()
		ok = ok && m.EqualsFinished()
		ok = ok && m.GetDigestHolderFinished()
		ok = ok && m.GetSignatureHolderFinished()
		ok = ok && m.GetSignatureMethodFinished()
		ok = ok && m.IsVerifiableByFinished()
		ok = ok && m.VerifyWithFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CopyOfSignedDigestFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.CopyOfSignedDigest")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.Equals")
			}

			if !m.GetDigestHolderFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.GetDigestHolder")
			}

			if !m.GetSignatureHolderFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.GetSignatureHolder")
			}

			if !m.GetSignatureMethodFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.GetSignatureMethod")
			}

			if !m.IsVerifiableByFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.IsVerifiableBy")
			}

			if !m.VerifyWithFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.VerifyWith")
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
func (m *NodeStateHashEvidenceMock) AllMocksCalled() bool {

	if !m.CopyOfSignedDigestFinished() {
		return false
	}

	if !m.EqualsFinished() {
		return false
	}

	if !m.GetDigestHolderFinished() {
		return false
	}

	if !m.GetSignatureHolderFinished() {
		return false
	}

	if !m.GetSignatureMethodFinished() {
		return false
	}

	if !m.IsVerifiableByFinished() {
		return false
	}

	if !m.VerifyWithFinished() {
		return false
	}

	return true
}
