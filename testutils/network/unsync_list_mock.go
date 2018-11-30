package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "UnsyncList" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	packets "github.com/insolar/insolar/consensus/packets"
	core "github.com/insolar/insolar/core"
	testify_assert "github.com/stretchr/testify/assert"
)

//UnsyncListMock implements github.com/insolar/insolar/network.UnsyncList
type UnsyncListMock struct {
	t minimock.Tester

	AddClaimsFunc       func(p core.RecordRef, p1 []packets.ReferendumClaim)
	AddClaimsCounter    uint64
	AddClaimsPreCounter uint64
	AddClaimsMock       mUnsyncListMockAddClaims

	CalculateHashFunc       func() (r []byte, r1 error)
	CalculateHashCounter    uint64
	CalculateHashPreCounter uint64
	CalculateHashMock       mUnsyncListMockCalculateHash

	IndexToRefFunc       func(p int) (r core.RecordRef, r1 error)
	IndexToRefCounter    uint64
	IndexToRefPreCounter uint64
	IndexToRefMock       mUnsyncListMockIndexToRef

	LengthFunc       func() (r int)
	LengthCounter    uint64
	LengthPreCounter uint64
	LengthMock       mUnsyncListMockLength

	RefToIndexFunc       func(p core.RecordRef) (r int, r1 error)
	RefToIndexCounter    uint64
	RefToIndexPreCounter uint64
	RefToIndexMock       mUnsyncListMockRefToIndex

	RemoveClaimsFunc       func(p core.RecordRef)
	RemoveClaimsCounter    uint64
	RemoveClaimsPreCounter uint64
	RemoveClaimsMock       mUnsyncListMockRemoveClaims
}

//NewUnsyncListMock returns a mock for github.com/insolar/insolar/network.UnsyncList
func NewUnsyncListMock(t minimock.Tester) *UnsyncListMock {
	m := &UnsyncListMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddClaimsMock = mUnsyncListMockAddClaims{mock: m}
	m.CalculateHashMock = mUnsyncListMockCalculateHash{mock: m}
	m.IndexToRefMock = mUnsyncListMockIndexToRef{mock: m}
	m.LengthMock = mUnsyncListMockLength{mock: m}
	m.RefToIndexMock = mUnsyncListMockRefToIndex{mock: m}
	m.RemoveClaimsMock = mUnsyncListMockRemoveClaims{mock: m}

	return m
}

type mUnsyncListMockAddClaims struct {
	mock             *UnsyncListMock
	mockExpectations *UnsyncListMockAddClaimsParams
}

//UnsyncListMockAddClaimsParams represents input parameters of the UnsyncList.AddClaims
type UnsyncListMockAddClaimsParams struct {
	p  core.RecordRef
	p1 []packets.ReferendumClaim
}

//Expect sets up expected params for the UnsyncList.AddClaims
func (m *mUnsyncListMockAddClaims) Expect(p core.RecordRef, p1 []packets.ReferendumClaim) *mUnsyncListMockAddClaims {
	m.mockExpectations = &UnsyncListMockAddClaimsParams{p, p1}
	return m
}

//Return sets up a mock for UnsyncList.AddClaims to return Return's arguments
func (m *mUnsyncListMockAddClaims) Return() *UnsyncListMock {
	m.mock.AddClaimsFunc = func(p core.RecordRef, p1 []packets.ReferendumClaim) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of UnsyncList.AddClaims method
func (m *mUnsyncListMockAddClaims) Set(f func(p core.RecordRef, p1 []packets.ReferendumClaim)) *UnsyncListMock {
	m.mock.AddClaimsFunc = f
	m.mockExpectations = nil
	return m.mock
}

//AddClaims implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) AddClaims(p core.RecordRef, p1 []packets.ReferendumClaim) {
	atomic.AddUint64(&m.AddClaimsPreCounter, 1)
	defer atomic.AddUint64(&m.AddClaimsCounter, 1)

	if m.AddClaimsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.AddClaimsMock.mockExpectations, UnsyncListMockAddClaimsParams{p, p1},
			"UnsyncList.AddClaims got unexpected parameters")

		if m.AddClaimsFunc == nil {

			m.t.Fatal("No results are set for the UnsyncListMock.AddClaims")

			return
		}
	}

	if m.AddClaimsFunc == nil {
		m.t.Fatal("Unexpected call to UnsyncListMock.AddClaims")
		return
	}

	m.AddClaimsFunc(p, p1)
}

//AddClaimsMinimockCounter returns a count of UnsyncListMock.AddClaimsFunc invocations
func (m *UnsyncListMock) AddClaimsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddClaimsCounter)
}

//AddClaimsMinimockPreCounter returns the value of UnsyncListMock.AddClaims invocations
func (m *UnsyncListMock) AddClaimsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddClaimsPreCounter)
}

type mUnsyncListMockCalculateHash struct {
	mock *UnsyncListMock
}

//Return sets up a mock for UnsyncList.CalculateHash to return Return's arguments
func (m *mUnsyncListMockCalculateHash) Return(r []byte, r1 error) *UnsyncListMock {
	m.mock.CalculateHashFunc = func() ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of UnsyncList.CalculateHash method
func (m *mUnsyncListMockCalculateHash) Set(f func() (r []byte, r1 error)) *UnsyncListMock {
	m.mock.CalculateHashFunc = f

	return m.mock
}

//CalculateHash implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) CalculateHash() (r []byte, r1 error) {
	atomic.AddUint64(&m.CalculateHashPreCounter, 1)
	defer atomic.AddUint64(&m.CalculateHashCounter, 1)

	if m.CalculateHashFunc == nil {
		m.t.Fatal("Unexpected call to UnsyncListMock.CalculateHash")
		return
	}

	return m.CalculateHashFunc()
}

//CalculateHashMinimockCounter returns a count of UnsyncListMock.CalculateHashFunc invocations
func (m *UnsyncListMock) CalculateHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CalculateHashCounter)
}

//CalculateHashMinimockPreCounter returns the value of UnsyncListMock.CalculateHash invocations
func (m *UnsyncListMock) CalculateHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CalculateHashPreCounter)
}

type mUnsyncListMockIndexToRef struct {
	mock             *UnsyncListMock
	mockExpectations *UnsyncListMockIndexToRefParams
}

//UnsyncListMockIndexToRefParams represents input parameters of the UnsyncList.IndexToRef
type UnsyncListMockIndexToRefParams struct {
	p int
}

//Expect sets up expected params for the UnsyncList.IndexToRef
func (m *mUnsyncListMockIndexToRef) Expect(p int) *mUnsyncListMockIndexToRef {
	m.mockExpectations = &UnsyncListMockIndexToRefParams{p}
	return m
}

//Return sets up a mock for UnsyncList.IndexToRef to return Return's arguments
func (m *mUnsyncListMockIndexToRef) Return(r core.RecordRef, r1 error) *UnsyncListMock {
	m.mock.IndexToRefFunc = func(p int) (core.RecordRef, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of UnsyncList.IndexToRef method
func (m *mUnsyncListMockIndexToRef) Set(f func(p int) (r core.RecordRef, r1 error)) *UnsyncListMock {
	m.mock.IndexToRefFunc = f
	m.mockExpectations = nil
	return m.mock
}

//IndexToRef implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) IndexToRef(p int) (r core.RecordRef, r1 error) {
	atomic.AddUint64(&m.IndexToRefPreCounter, 1)
	defer atomic.AddUint64(&m.IndexToRefCounter, 1)

	if m.IndexToRefMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.IndexToRefMock.mockExpectations, UnsyncListMockIndexToRefParams{p},
			"UnsyncList.IndexToRef got unexpected parameters")

		if m.IndexToRefFunc == nil {

			m.t.Fatal("No results are set for the UnsyncListMock.IndexToRef")

			return
		}
	}

	if m.IndexToRefFunc == nil {
		m.t.Fatal("Unexpected call to UnsyncListMock.IndexToRef")
		return
	}

	return m.IndexToRefFunc(p)
}

//IndexToRefMinimockCounter returns a count of UnsyncListMock.IndexToRefFunc invocations
func (m *UnsyncListMock) IndexToRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IndexToRefCounter)
}

//IndexToRefMinimockPreCounter returns the value of UnsyncListMock.IndexToRef invocations
func (m *UnsyncListMock) IndexToRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IndexToRefPreCounter)
}

type mUnsyncListMockLength struct {
	mock *UnsyncListMock
}

//Return sets up a mock for UnsyncList.Length to return Return's arguments
func (m *mUnsyncListMockLength) Return(r int) *UnsyncListMock {
	m.mock.LengthFunc = func() int {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of UnsyncList.Length method
func (m *mUnsyncListMockLength) Set(f func() (r int)) *UnsyncListMock {
	m.mock.LengthFunc = f

	return m.mock
}

//Length implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) Length() (r int) {
	atomic.AddUint64(&m.LengthPreCounter, 1)
	defer atomic.AddUint64(&m.LengthCounter, 1)

	if m.LengthFunc == nil {
		m.t.Fatal("Unexpected call to UnsyncListMock.Length")
		return
	}

	return m.LengthFunc()
}

//LengthMinimockCounter returns a count of UnsyncListMock.LengthFunc invocations
func (m *UnsyncListMock) LengthMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LengthCounter)
}

//LengthMinimockPreCounter returns the value of UnsyncListMock.Length invocations
func (m *UnsyncListMock) LengthMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LengthPreCounter)
}

type mUnsyncListMockRefToIndex struct {
	mock             *UnsyncListMock
	mockExpectations *UnsyncListMockRefToIndexParams
}

//UnsyncListMockRefToIndexParams represents input parameters of the UnsyncList.RefToIndex
type UnsyncListMockRefToIndexParams struct {
	p core.RecordRef
}

//Expect sets up expected params for the UnsyncList.RefToIndex
func (m *mUnsyncListMockRefToIndex) Expect(p core.RecordRef) *mUnsyncListMockRefToIndex {
	m.mockExpectations = &UnsyncListMockRefToIndexParams{p}
	return m
}

//Return sets up a mock for UnsyncList.RefToIndex to return Return's arguments
func (m *mUnsyncListMockRefToIndex) Return(r int, r1 error) *UnsyncListMock {
	m.mock.RefToIndexFunc = func(p core.RecordRef) (int, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of UnsyncList.RefToIndex method
func (m *mUnsyncListMockRefToIndex) Set(f func(p core.RecordRef) (r int, r1 error)) *UnsyncListMock {
	m.mock.RefToIndexFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RefToIndex implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) RefToIndex(p core.RecordRef) (r int, r1 error) {
	atomic.AddUint64(&m.RefToIndexPreCounter, 1)
	defer atomic.AddUint64(&m.RefToIndexCounter, 1)

	if m.RefToIndexMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RefToIndexMock.mockExpectations, UnsyncListMockRefToIndexParams{p},
			"UnsyncList.RefToIndex got unexpected parameters")

		if m.RefToIndexFunc == nil {

			m.t.Fatal("No results are set for the UnsyncListMock.RefToIndex")

			return
		}
	}

	if m.RefToIndexFunc == nil {
		m.t.Fatal("Unexpected call to UnsyncListMock.RefToIndex")
		return
	}

	return m.RefToIndexFunc(p)
}

//RefToIndexMinimockCounter returns a count of UnsyncListMock.RefToIndexFunc invocations
func (m *UnsyncListMock) RefToIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RefToIndexCounter)
}

//RefToIndexMinimockPreCounter returns the value of UnsyncListMock.RefToIndex invocations
func (m *UnsyncListMock) RefToIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RefToIndexPreCounter)
}

type mUnsyncListMockRemoveClaims struct {
	mock             *UnsyncListMock
	mockExpectations *UnsyncListMockRemoveClaimsParams
}

//UnsyncListMockRemoveClaimsParams represents input parameters of the UnsyncList.RemoveClaims
type UnsyncListMockRemoveClaimsParams struct {
	p core.RecordRef
}

//Expect sets up expected params for the UnsyncList.RemoveClaims
func (m *mUnsyncListMockRemoveClaims) Expect(p core.RecordRef) *mUnsyncListMockRemoveClaims {
	m.mockExpectations = &UnsyncListMockRemoveClaimsParams{p}
	return m
}

//Return sets up a mock for UnsyncList.RemoveClaims to return Return's arguments
func (m *mUnsyncListMockRemoveClaims) Return() *UnsyncListMock {
	m.mock.RemoveClaimsFunc = func(p core.RecordRef) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of UnsyncList.RemoveClaims method
func (m *mUnsyncListMockRemoveClaims) Set(f func(p core.RecordRef)) *UnsyncListMock {
	m.mock.RemoveClaimsFunc = f
	m.mockExpectations = nil
	return m.mock
}

//RemoveClaims implements github.com/insolar/insolar/network.UnsyncList interface
func (m *UnsyncListMock) RemoveClaims(p core.RecordRef) {
	atomic.AddUint64(&m.RemoveClaimsPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveClaimsCounter, 1)

	if m.RemoveClaimsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.RemoveClaimsMock.mockExpectations, UnsyncListMockRemoveClaimsParams{p},
			"UnsyncList.RemoveClaims got unexpected parameters")

		if m.RemoveClaimsFunc == nil {

			m.t.Fatal("No results are set for the UnsyncListMock.RemoveClaims")

			return
		}
	}

	if m.RemoveClaimsFunc == nil {
		m.t.Fatal("Unexpected call to UnsyncListMock.RemoveClaims")
		return
	}

	m.RemoveClaimsFunc(p)
}

//RemoveClaimsMinimockCounter returns a count of UnsyncListMock.RemoveClaimsFunc invocations
func (m *UnsyncListMock) RemoveClaimsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveClaimsCounter)
}

//RemoveClaimsMinimockPreCounter returns the value of UnsyncListMock.RemoveClaims invocations
func (m *UnsyncListMock) RemoveClaimsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveClaimsPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *UnsyncListMock) ValidateCallCounters() {

	if m.AddClaimsFunc != nil && atomic.LoadUint64(&m.AddClaimsCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.AddClaims")
	}

	if m.CalculateHashFunc != nil && atomic.LoadUint64(&m.CalculateHashCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.CalculateHash")
	}

	if m.IndexToRefFunc != nil && atomic.LoadUint64(&m.IndexToRefCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.IndexToRef")
	}

	if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.Length")
	}

	if m.RefToIndexFunc != nil && atomic.LoadUint64(&m.RefToIndexCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.RefToIndex")
	}

	if m.RemoveClaimsFunc != nil && atomic.LoadUint64(&m.RemoveClaimsCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.RemoveClaims")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *UnsyncListMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *UnsyncListMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *UnsyncListMock) MinimockFinish() {

	if m.AddClaimsFunc != nil && atomic.LoadUint64(&m.AddClaimsCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.AddClaims")
	}

	if m.CalculateHashFunc != nil && atomic.LoadUint64(&m.CalculateHashCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.CalculateHash")
	}

	if m.IndexToRefFunc != nil && atomic.LoadUint64(&m.IndexToRefCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.IndexToRef")
	}

	if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.Length")
	}

	if m.RefToIndexFunc != nil && atomic.LoadUint64(&m.RefToIndexCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.RefToIndex")
	}

	if m.RemoveClaimsFunc != nil && atomic.LoadUint64(&m.RemoveClaimsCounter) == 0 {
		m.t.Fatal("Expected call to UnsyncListMock.RemoveClaims")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *UnsyncListMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *UnsyncListMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.AddClaimsFunc == nil || atomic.LoadUint64(&m.AddClaimsCounter) > 0)
		ok = ok && (m.CalculateHashFunc == nil || atomic.LoadUint64(&m.CalculateHashCounter) > 0)
		ok = ok && (m.IndexToRefFunc == nil || atomic.LoadUint64(&m.IndexToRefCounter) > 0)
		ok = ok && (m.LengthFunc == nil || atomic.LoadUint64(&m.LengthCounter) > 0)
		ok = ok && (m.RefToIndexFunc == nil || atomic.LoadUint64(&m.RefToIndexCounter) > 0)
		ok = ok && (m.RemoveClaimsFunc == nil || atomic.LoadUint64(&m.RemoveClaimsCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.AddClaimsFunc != nil && atomic.LoadUint64(&m.AddClaimsCounter) == 0 {
				m.t.Error("Expected call to UnsyncListMock.AddClaims")
			}

			if m.CalculateHashFunc != nil && atomic.LoadUint64(&m.CalculateHashCounter) == 0 {
				m.t.Error("Expected call to UnsyncListMock.CalculateHash")
			}

			if m.IndexToRefFunc != nil && atomic.LoadUint64(&m.IndexToRefCounter) == 0 {
				m.t.Error("Expected call to UnsyncListMock.IndexToRef")
			}

			if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
				m.t.Error("Expected call to UnsyncListMock.Length")
			}

			if m.RefToIndexFunc != nil && atomic.LoadUint64(&m.RefToIndexCounter) == 0 {
				m.t.Error("Expected call to UnsyncListMock.RefToIndex")
			}

			if m.RemoveClaimsFunc != nil && atomic.LoadUint64(&m.RemoveClaimsCounter) == 0 {
				m.t.Error("Expected call to UnsyncListMock.RemoveClaims")
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
func (m *UnsyncListMock) AllMocksCalled() bool {

	if m.AddClaimsFunc != nil && atomic.LoadUint64(&m.AddClaimsCounter) == 0 {
		return false
	}

	if m.CalculateHashFunc != nil && atomic.LoadUint64(&m.CalculateHashCounter) == 0 {
		return false
	}

	if m.IndexToRefFunc != nil && atomic.LoadUint64(&m.IndexToRefCounter) == 0 {
		return false
	}

	if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
		return false
	}

	if m.RefToIndexFunc != nil && atomic.LoadUint64(&m.RefToIndexCounter) == 0 {
		return false
	}

	if m.RemoveClaimsFunc != nil && atomic.LoadUint64(&m.RemoveClaimsCounter) == 0 {
		return false
	}

	return true
}
