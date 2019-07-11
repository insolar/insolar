package gcp_types

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeStateHashEvidence" can be found in github.com/insolar/insolar/network/consensus/gcpv2/common
*/
import (
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//NodeStateHashEvidenceMock implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeStateHashEvidence
type NodeStateHashEvidenceMock struct {
	t minimock.Tester

	GetGlobulaNodeStateSignatureFunc       func() (r cryptography_containers.SignatureHolder)
	GetGlobulaNodeStateSignatureCounter    uint64
	GetGlobulaNodeStateSignaturePreCounter uint64
	GetGlobulaNodeStateSignatureMock       mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature

	GetNodeStateHashFunc       func() (r NodeStateHash)
	GetNodeStateHashCounter    uint64
	GetNodeStateHashPreCounter uint64
	GetNodeStateHashMock       mNodeStateHashEvidenceMockGetNodeStateHash
}

//NewNodeStateHashEvidenceMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/common.NodeStateHashEvidence
func NewNodeStateHashEvidenceMock(t minimock.Tester) *NodeStateHashEvidenceMock {
	m := &NodeStateHashEvidenceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetGlobulaNodeStateSignatureMock = mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature{mock: m}
	m.GetNodeStateHashMock = mNodeStateHashEvidenceMockGetNodeStateHash{mock: m}

	return m
}

type mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation
	expectationSeries []*NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation
}

type NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation struct {
	result *NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureResult
}

type NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureResult struct {
	r cryptography_containers.SignatureHolder
}

//Expect specifies that invocation of NodeStateHashEvidence.GetGlobulaNodeStateSignature is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature) Expect() *mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature {
	m.mock.GetGlobulaNodeStateSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.GetGlobulaNodeStateSignature
func (m *mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature) Return(r cryptography_containers.SignatureHolder) *NodeStateHashEvidenceMock {
	m.mock.GetGlobulaNodeStateSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.GetGlobulaNodeStateSignature is expected once
func (m *mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature) ExpectOnce() *NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation {
	m.mock.GetGlobulaNodeStateSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureExpectation) Return(r cryptography_containers.SignatureHolder) {
	e.result = &NodeStateHashEvidenceMockGetGlobulaNodeStateSignatureResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.GetGlobulaNodeStateSignature method
func (m *mNodeStateHashEvidenceMockGetGlobulaNodeStateSignature) Set(f func() (r cryptography_containers.SignatureHolder)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetGlobulaNodeStateSignatureFunc = f
	return m.mock
}

//GetGlobulaNodeStateSignature implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) GetGlobulaNodeStateSignature() (r cryptography_containers.SignatureHolder) {
	counter := atomic.AddUint64(&m.GetGlobulaNodeStateSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetGlobulaNodeStateSignatureCounter, 1)

	if len(m.GetGlobulaNodeStateSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetGlobulaNodeStateSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature.")
			return
		}

		result := m.GetGlobulaNodeStateSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetGlobulaNodeStateSignatureMock.mainExpectation != nil {

		result := m.GetGlobulaNodeStateSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature")
		}

		r = result.r

		return
	}

	if m.GetGlobulaNodeStateSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature.")
		return
	}

	return m.GetGlobulaNodeStateSignatureFunc()
}

//GetGlobulaNodeStateSignatureMinimockCounter returns a count of NodeStateHashEvidenceMock.GetGlobulaNodeStateSignatureFunc invocations
func (m *NodeStateHashEvidenceMock) GetGlobulaNodeStateSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobulaNodeStateSignatureCounter)
}

//GetGlobulaNodeStateSignatureMinimockPreCounter returns the value of NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature invocations
func (m *NodeStateHashEvidenceMock) GetGlobulaNodeStateSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetGlobulaNodeStateSignaturePreCounter)
}

//GetGlobulaNodeStateSignatureFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) GetGlobulaNodeStateSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetGlobulaNodeStateSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetGlobulaNodeStateSignatureCounter) == uint64(len(m.GetGlobulaNodeStateSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetGlobulaNodeStateSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetGlobulaNodeStateSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetGlobulaNodeStateSignatureFunc != nil {
		return atomic.LoadUint64(&m.GetGlobulaNodeStateSignatureCounter) > 0
	}

	return true
}

type mNodeStateHashEvidenceMockGetNodeStateHash struct {
	mock              *NodeStateHashEvidenceMock
	mainExpectation   *NodeStateHashEvidenceMockGetNodeStateHashExpectation
	expectationSeries []*NodeStateHashEvidenceMockGetNodeStateHashExpectation
}

type NodeStateHashEvidenceMockGetNodeStateHashExpectation struct {
	result *NodeStateHashEvidenceMockGetNodeStateHashResult
}

type NodeStateHashEvidenceMockGetNodeStateHashResult struct {
	r NodeStateHash
}

//Expect specifies that invocation of NodeStateHashEvidence.GetNodeStateHash is expected from 1 to Infinity times
func (m *mNodeStateHashEvidenceMockGetNodeStateHash) Expect() *mNodeStateHashEvidenceMockGetNodeStateHash {
	m.mock.GetNodeStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetNodeStateHashExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHashEvidence.GetNodeStateHash
func (m *mNodeStateHashEvidenceMockGetNodeStateHash) Return(r NodeStateHash) *NodeStateHashEvidenceMock {
	m.mock.GetNodeStateHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashEvidenceMockGetNodeStateHashExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashEvidenceMockGetNodeStateHashResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHashEvidence.GetNodeStateHash is expected once
func (m *mNodeStateHashEvidenceMockGetNodeStateHash) ExpectOnce() *NodeStateHashEvidenceMockGetNodeStateHashExpectation {
	m.mock.GetNodeStateHashFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashEvidenceMockGetNodeStateHashExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashEvidenceMockGetNodeStateHashExpectation) Return(r NodeStateHash) {
	e.result = &NodeStateHashEvidenceMockGetNodeStateHashResult{r}
}

//Set uses given function f as a mock of NodeStateHashEvidence.GetNodeStateHash method
func (m *mNodeStateHashEvidenceMockGetNodeStateHash) Set(f func() (r NodeStateHash)) *NodeStateHashEvidenceMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeStateHashFunc = f
	return m.mock
}

//GetNodeStateHash implements github.com/insolar/insolar/network/consensus/gcpv2/common.NodeStateHashEvidence interface
func (m *NodeStateHashEvidenceMock) GetNodeStateHash() (r NodeStateHash) {
	counter := atomic.AddUint64(&m.GetNodeStateHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeStateHashCounter, 1)

	if len(m.GetNodeStateHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeStateHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetNodeStateHash.")
			return
		}

		result := m.GetNodeStateHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetNodeStateHash")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeStateHashMock.mainExpectation != nil {

		result := m.GetNodeStateHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashEvidenceMock.GetNodeStateHash")
		}

		r = result.r

		return
	}

	if m.GetNodeStateHashFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashEvidenceMock.GetNodeStateHash.")
		return
	}

	return m.GetNodeStateHashFunc()
}

//GetNodeStateHashMinimockCounter returns a count of NodeStateHashEvidenceMock.GetNodeStateHashFunc invocations
func (m *NodeStateHashEvidenceMock) GetNodeStateHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeStateHashCounter)
}

//GetNodeStateHashMinimockPreCounter returns the value of NodeStateHashEvidenceMock.GetNodeStateHash invocations
func (m *NodeStateHashEvidenceMock) GetNodeStateHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeStateHashPreCounter)
}

//GetNodeStateHashFinished returns true if mock invocations count is ok
func (m *NodeStateHashEvidenceMock) GetNodeStateHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeStateHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeStateHashCounter) == uint64(len(m.GetNodeStateHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeStateHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeStateHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeStateHashFunc != nil {
		return atomic.LoadUint64(&m.GetNodeStateHashCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeStateHashEvidenceMock) ValidateCallCounters() {

	if !m.GetGlobulaNodeStateSignatureFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature")
	}

	if !m.GetNodeStateHashFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetNodeStateHash")
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

	if !m.GetGlobulaNodeStateSignatureFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature")
	}

	if !m.GetNodeStateHashFinished() {
		m.t.Fatal("Expected call to NodeStateHashEvidenceMock.GetNodeStateHash")
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
		ok = ok && m.GetGlobulaNodeStateSignatureFinished()
		ok = ok && m.GetNodeStateHashFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetGlobulaNodeStateSignatureFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.GetGlobulaNodeStateSignature")
			}

			if !m.GetNodeStateHashFinished() {
				m.t.Error("Expected call to NodeStateHashEvidenceMock.GetNodeStateHash")
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

	if !m.GetGlobulaNodeStateSignatureFinished() {
		return false
	}

	if !m.GetNodeStateHashFinished() {
		return false
	}

	return true
}
