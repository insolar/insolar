package proofs

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "NodeStateHash" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/proofs
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"

	testify_assert "github.com/stretchr/testify/assert"
)

//NodeStateHashMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash
type NodeStateHashMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mNodeStateHashMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mNodeStateHashMockAsBytes

	CopyOfDigestFunc       func() (r cryptkit.Digest)
	CopyOfDigestCounter    uint64
	CopyOfDigestPreCounter uint64
	CopyOfDigestMock       mNodeStateHashMockCopyOfDigest

	EqualsFunc       func(p cryptkit.DigestHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mNodeStateHashMockEquals

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mNodeStateHashMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mNodeStateHashMockFoldToUint64

	GetDigestMethodFunc       func() (r cryptkit.DigestMethod)
	GetDigestMethodCounter    uint64
	GetDigestMethodPreCounter uint64
	GetDigestMethodMock       mNodeStateHashMockGetDigestMethod

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mNodeStateHashMockRead

	SignWithFunc       func(p cryptkit.DigestSigner) (r cryptkit.SignedDigest)
	SignWithCounter    uint64
	SignWithPreCounter uint64
	SignWithMock       mNodeStateHashMockSignWith

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mNodeStateHashMockWriteTo
}

//NewNodeStateHashMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash
func NewNodeStateHashMock(t minimock.Tester) *NodeStateHashMock {
	m := &NodeStateHashMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mNodeStateHashMockAsByteString{mock: m}
	m.AsBytesMock = mNodeStateHashMockAsBytes{mock: m}
	m.CopyOfDigestMock = mNodeStateHashMockCopyOfDigest{mock: m}
	m.EqualsMock = mNodeStateHashMockEquals{mock: m}
	m.FixedByteSizeMock = mNodeStateHashMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mNodeStateHashMockFoldToUint64{mock: m}
	m.GetDigestMethodMock = mNodeStateHashMockGetDigestMethod{mock: m}
	m.ReadMock = mNodeStateHashMockRead{mock: m}
	m.SignWithMock = mNodeStateHashMockSignWith{mock: m}
	m.WriteToMock = mNodeStateHashMockWriteTo{mock: m}

	return m
}

type mNodeStateHashMockAsByteString struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockAsByteStringExpectation
	expectationSeries []*NodeStateHashMockAsByteStringExpectation
}

type NodeStateHashMockAsByteStringExpectation struct {
	result *NodeStateHashMockAsByteStringResult
}

type NodeStateHashMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of NodeStateHash.AsByteString is expected from 1 to Infinity times
func (m *mNodeStateHashMockAsByteString) Expect() *mNodeStateHashMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHash.AsByteString
func (m *mNodeStateHashMockAsByteString) Return(r string) *NodeStateHashMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.AsByteString is expected once
func (m *mNodeStateHashMockAsByteString) ExpectOnce() *NodeStateHashMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockAsByteStringExpectation) Return(r string) {
	e.result = &NodeStateHashMockAsByteStringResult{r}
}

//Set uses given function f as a mock of NodeStateHash.AsByteString method
func (m *mNodeStateHashMockAsByteString) Set(f func() (r string)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of NodeStateHashMock.AsByteStringFunc invocations
func (m *NodeStateHashMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of NodeStateHashMock.AsByteString invocations
func (m *NodeStateHashMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) AsByteStringFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AsByteStringMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AsByteStringCounter) == uint64(len(m.AsByteStringMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AsByteStringMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AsByteStringCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AsByteStringFunc != nil {
		return atomic.LoadUint64(&m.AsByteStringCounter) > 0
	}

	return true
}

type mNodeStateHashMockAsBytes struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockAsBytesExpectation
	expectationSeries []*NodeStateHashMockAsBytesExpectation
}

type NodeStateHashMockAsBytesExpectation struct {
	result *NodeStateHashMockAsBytesResult
}

type NodeStateHashMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of NodeStateHash.AsBytes is expected from 1 to Infinity times
func (m *mNodeStateHashMockAsBytes) Expect() *mNodeStateHashMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHash.AsBytes
func (m *mNodeStateHashMockAsBytes) Return(r []byte) *NodeStateHashMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.AsBytes is expected once
func (m *mNodeStateHashMockAsBytes) ExpectOnce() *NodeStateHashMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockAsBytesExpectation) Return(r []byte) {
	e.result = &NodeStateHashMockAsBytesResult{r}
}

//Set uses given function f as a mock of NodeStateHash.AsBytes method
func (m *mNodeStateHashMockAsBytes) Set(f func() (r []byte)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of NodeStateHashMock.AsBytesFunc invocations
func (m *NodeStateHashMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of NodeStateHashMock.AsBytes invocations
func (m *NodeStateHashMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) AsBytesFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AsBytesMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AsBytesCounter) == uint64(len(m.AsBytesMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AsBytesMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AsBytesCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AsBytesFunc != nil {
		return atomic.LoadUint64(&m.AsBytesCounter) > 0
	}

	return true
}

type mNodeStateHashMockCopyOfDigest struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockCopyOfDigestExpectation
	expectationSeries []*NodeStateHashMockCopyOfDigestExpectation
}

type NodeStateHashMockCopyOfDigestExpectation struct {
	result *NodeStateHashMockCopyOfDigestResult
}

type NodeStateHashMockCopyOfDigestResult struct {
	r cryptkit.Digest
}

//Expect specifies that invocation of NodeStateHash.CopyOfDigest is expected from 1 to Infinity times
func (m *mNodeStateHashMockCopyOfDigest) Expect() *mNodeStateHashMockCopyOfDigest {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockCopyOfDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHash.CopyOfDigest
func (m *mNodeStateHashMockCopyOfDigest) Return(r cryptkit.Digest) *NodeStateHashMock {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockCopyOfDigestExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockCopyOfDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.CopyOfDigest is expected once
func (m *mNodeStateHashMockCopyOfDigest) ExpectOnce() *NodeStateHashMockCopyOfDigestExpectation {
	m.mock.CopyOfDigestFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockCopyOfDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockCopyOfDigestExpectation) Return(r cryptkit.Digest) {
	e.result = &NodeStateHashMockCopyOfDigestResult{r}
}

//Set uses given function f as a mock of NodeStateHash.CopyOfDigest method
func (m *mNodeStateHashMockCopyOfDigest) Set(f func() (r cryptkit.Digest)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CopyOfDigestFunc = f
	return m.mock
}

//CopyOfDigest implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) CopyOfDigest() (r cryptkit.Digest) {
	counter := atomic.AddUint64(&m.CopyOfDigestPreCounter, 1)
	defer atomic.AddUint64(&m.CopyOfDigestCounter, 1)

	if len(m.CopyOfDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CopyOfDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.CopyOfDigest.")
			return
		}

		result := m.CopyOfDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.CopyOfDigest")
			return
		}

		r = result.r

		return
	}

	if m.CopyOfDigestMock.mainExpectation != nil {

		result := m.CopyOfDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.CopyOfDigest")
		}

		r = result.r

		return
	}

	if m.CopyOfDigestFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.CopyOfDigest.")
		return
	}

	return m.CopyOfDigestFunc()
}

//CopyOfDigestMinimockCounter returns a count of NodeStateHashMock.CopyOfDigestFunc invocations
func (m *NodeStateHashMock) CopyOfDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestCounter)
}

//CopyOfDigestMinimockPreCounter returns the value of NodeStateHashMock.CopyOfDigest invocations
func (m *NodeStateHashMock) CopyOfDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestPreCounter)
}

//CopyOfDigestFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) CopyOfDigestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CopyOfDigestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CopyOfDigestCounter) == uint64(len(m.CopyOfDigestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CopyOfDigestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CopyOfDigestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CopyOfDigestFunc != nil {
		return atomic.LoadUint64(&m.CopyOfDigestCounter) > 0
	}

	return true
}

type mNodeStateHashMockEquals struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockEqualsExpectation
	expectationSeries []*NodeStateHashMockEqualsExpectation
}

type NodeStateHashMockEqualsExpectation struct {
	input  *NodeStateHashMockEqualsInput
	result *NodeStateHashMockEqualsResult
}

type NodeStateHashMockEqualsInput struct {
	p cryptkit.DigestHolder
}

type NodeStateHashMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of NodeStateHash.Equals is expected from 1 to Infinity times
func (m *mNodeStateHashMockEquals) Expect(p cryptkit.DigestHolder) *mNodeStateHashMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockEqualsExpectation{}
	}
	m.mainExpectation.input = &NodeStateHashMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of NodeStateHash.Equals
func (m *mNodeStateHashMockEquals) Return(r bool) *NodeStateHashMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockEqualsExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.Equals is expected once
func (m *mNodeStateHashMockEquals) ExpectOnce(p cryptkit.DigestHolder) *NodeStateHashMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockEqualsExpectation{}
	expectation.input = &NodeStateHashMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockEqualsExpectation) Return(r bool) {
	e.result = &NodeStateHashMockEqualsResult{r}
}

//Set uses given function f as a mock of NodeStateHash.Equals method
func (m *mNodeStateHashMockEquals) Set(f func(p cryptkit.DigestHolder) (r bool)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) Equals(p cryptkit.DigestHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStateHashMockEqualsInput{p}, "NodeStateHash.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStateHashMockEqualsInput{p}, "NodeStateHash.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of NodeStateHashMock.EqualsFunc invocations
func (m *NodeStateHashMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of NodeStateHashMock.Equals invocations
func (m *NodeStateHashMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) EqualsFinished() bool {
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

type mNodeStateHashMockFixedByteSize struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockFixedByteSizeExpectation
	expectationSeries []*NodeStateHashMockFixedByteSizeExpectation
}

type NodeStateHashMockFixedByteSizeExpectation struct {
	result *NodeStateHashMockFixedByteSizeResult
}

type NodeStateHashMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of NodeStateHash.FixedByteSize is expected from 1 to Infinity times
func (m *mNodeStateHashMockFixedByteSize) Expect() *mNodeStateHashMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHash.FixedByteSize
func (m *mNodeStateHashMockFixedByteSize) Return(r int) *NodeStateHashMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.FixedByteSize is expected once
func (m *mNodeStateHashMockFixedByteSize) ExpectOnce() *NodeStateHashMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockFixedByteSizeExpectation) Return(r int) {
	e.result = &NodeStateHashMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of NodeStateHash.FixedByteSize method
func (m *mNodeStateHashMockFixedByteSize) Set(f func() (r int)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of NodeStateHashMock.FixedByteSizeFunc invocations
func (m *NodeStateHashMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of NodeStateHashMock.FixedByteSize invocations
func (m *NodeStateHashMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) FixedByteSizeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FixedByteSizeCounter) == uint64(len(m.FixedByteSizeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FixedByteSizeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FixedByteSizeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FixedByteSizeFunc != nil {
		return atomic.LoadUint64(&m.FixedByteSizeCounter) > 0
	}

	return true
}

type mNodeStateHashMockFoldToUint64 struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockFoldToUint64Expectation
	expectationSeries []*NodeStateHashMockFoldToUint64Expectation
}

type NodeStateHashMockFoldToUint64Expectation struct {
	result *NodeStateHashMockFoldToUint64Result
}

type NodeStateHashMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of NodeStateHash.FoldToUint64 is expected from 1 to Infinity times
func (m *mNodeStateHashMockFoldToUint64) Expect() *mNodeStateHashMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHash.FoldToUint64
func (m *mNodeStateHashMockFoldToUint64) Return(r uint64) *NodeStateHashMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.FoldToUint64 is expected once
func (m *mNodeStateHashMockFoldToUint64) ExpectOnce() *NodeStateHashMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &NodeStateHashMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of NodeStateHash.FoldToUint64 method
func (m *mNodeStateHashMockFoldToUint64) Set(f func() (r uint64)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of NodeStateHashMock.FoldToUint64Func invocations
func (m *NodeStateHashMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of NodeStateHashMock.FoldToUint64 invocations
func (m *NodeStateHashMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *NodeStateHashMock) FoldToUint64Finished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FoldToUint64Counter) == uint64(len(m.FoldToUint64Mock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FoldToUint64Mock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FoldToUint64Counter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FoldToUint64Func != nil {
		return atomic.LoadUint64(&m.FoldToUint64Counter) > 0
	}

	return true
}

type mNodeStateHashMockGetDigestMethod struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockGetDigestMethodExpectation
	expectationSeries []*NodeStateHashMockGetDigestMethodExpectation
}

type NodeStateHashMockGetDigestMethodExpectation struct {
	result *NodeStateHashMockGetDigestMethodResult
}

type NodeStateHashMockGetDigestMethodResult struct {
	r cryptkit.DigestMethod
}

//Expect specifies that invocation of NodeStateHash.GetDigestMethod is expected from 1 to Infinity times
func (m *mNodeStateHashMockGetDigestMethod) Expect() *mNodeStateHashMockGetDigestMethod {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockGetDigestMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of NodeStateHash.GetDigestMethod
func (m *mNodeStateHashMockGetDigestMethod) Return(r cryptkit.DigestMethod) *NodeStateHashMock {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockGetDigestMethodExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockGetDigestMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.GetDigestMethod is expected once
func (m *mNodeStateHashMockGetDigestMethod) ExpectOnce() *NodeStateHashMockGetDigestMethodExpectation {
	m.mock.GetDigestMethodFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockGetDigestMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockGetDigestMethodExpectation) Return(r cryptkit.DigestMethod) {
	e.result = &NodeStateHashMockGetDigestMethodResult{r}
}

//Set uses given function f as a mock of NodeStateHash.GetDigestMethod method
func (m *mNodeStateHashMockGetDigestMethod) Set(f func() (r cryptkit.DigestMethod)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDigestMethodFunc = f
	return m.mock
}

//GetDigestMethod implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) GetDigestMethod() (r cryptkit.DigestMethod) {
	counter := atomic.AddUint64(&m.GetDigestMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetDigestMethodCounter, 1)

	if len(m.GetDigestMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDigestMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.GetDigestMethod.")
			return
		}

		result := m.GetDigestMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.GetDigestMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetDigestMethodMock.mainExpectation != nil {

		result := m.GetDigestMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.GetDigestMethod")
		}

		r = result.r

		return
	}

	if m.GetDigestMethodFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.GetDigestMethod.")
		return
	}

	return m.GetDigestMethodFunc()
}

//GetDigestMethodMinimockCounter returns a count of NodeStateHashMock.GetDigestMethodFunc invocations
func (m *NodeStateHashMock) GetDigestMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodCounter)
}

//GetDigestMethodMinimockPreCounter returns the value of NodeStateHashMock.GetDigestMethod invocations
func (m *NodeStateHashMock) GetDigestMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodPreCounter)
}

//GetDigestMethodFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) GetDigestMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDigestMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDigestMethodCounter) == uint64(len(m.GetDigestMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDigestMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDigestMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDigestMethodFunc != nil {
		return atomic.LoadUint64(&m.GetDigestMethodCounter) > 0
	}

	return true
}

type mNodeStateHashMockRead struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockReadExpectation
	expectationSeries []*NodeStateHashMockReadExpectation
}

type NodeStateHashMockReadExpectation struct {
	input  *NodeStateHashMockReadInput
	result *NodeStateHashMockReadResult
}

type NodeStateHashMockReadInput struct {
	p []byte
}

type NodeStateHashMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of NodeStateHash.Read is expected from 1 to Infinity times
func (m *mNodeStateHashMockRead) Expect(p []byte) *mNodeStateHashMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockReadExpectation{}
	}
	m.mainExpectation.input = &NodeStateHashMockReadInput{p}
	return m
}

//Return specifies results of invocation of NodeStateHash.Read
func (m *mNodeStateHashMockRead) Return(r int, r1 error) *NodeStateHashMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockReadExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.Read is expected once
func (m *mNodeStateHashMockRead) ExpectOnce(p []byte) *NodeStateHashMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockReadExpectation{}
	expectation.input = &NodeStateHashMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockReadExpectation) Return(r int, r1 error) {
	e.result = &NodeStateHashMockReadResult{r, r1}
}

//Set uses given function f as a mock of NodeStateHash.Read method
func (m *mNodeStateHashMockRead) Set(f func(p []byte) (r int, r1 error)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStateHashMockReadInput{p}, "NodeStateHash.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStateHashMockReadInput{p}, "NodeStateHash.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of NodeStateHashMock.ReadFunc invocations
func (m *NodeStateHashMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of NodeStateHashMock.Read invocations
func (m *NodeStateHashMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) ReadFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ReadMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ReadCounter) == uint64(len(m.ReadMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ReadMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ReadCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ReadFunc != nil {
		return atomic.LoadUint64(&m.ReadCounter) > 0
	}

	return true
}

type mNodeStateHashMockSignWith struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockSignWithExpectation
	expectationSeries []*NodeStateHashMockSignWithExpectation
}

type NodeStateHashMockSignWithExpectation struct {
	input  *NodeStateHashMockSignWithInput
	result *NodeStateHashMockSignWithResult
}

type NodeStateHashMockSignWithInput struct {
	p cryptkit.DigestSigner
}

type NodeStateHashMockSignWithResult struct {
	r cryptkit.SignedDigest
}

//Expect specifies that invocation of NodeStateHash.SignWith is expected from 1 to Infinity times
func (m *mNodeStateHashMockSignWith) Expect(p cryptkit.DigestSigner) *mNodeStateHashMockSignWith {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockSignWithExpectation{}
	}
	m.mainExpectation.input = &NodeStateHashMockSignWithInput{p}
	return m
}

//Return specifies results of invocation of NodeStateHash.SignWith
func (m *mNodeStateHashMockSignWith) Return(r cryptkit.SignedDigest) *NodeStateHashMock {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockSignWithExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockSignWithResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.SignWith is expected once
func (m *mNodeStateHashMockSignWith) ExpectOnce(p cryptkit.DigestSigner) *NodeStateHashMockSignWithExpectation {
	m.mock.SignWithFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockSignWithExpectation{}
	expectation.input = &NodeStateHashMockSignWithInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockSignWithExpectation) Return(r cryptkit.SignedDigest) {
	e.result = &NodeStateHashMockSignWithResult{r}
}

//Set uses given function f as a mock of NodeStateHash.SignWith method
func (m *mNodeStateHashMockSignWith) Set(f func(p cryptkit.DigestSigner) (r cryptkit.SignedDigest)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignWithFunc = f
	return m.mock
}

//SignWith implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) SignWith(p cryptkit.DigestSigner) (r cryptkit.SignedDigest) {
	counter := atomic.AddUint64(&m.SignWithPreCounter, 1)
	defer atomic.AddUint64(&m.SignWithCounter, 1)

	if len(m.SignWithMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignWithMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.SignWith. %v", p)
			return
		}

		input := m.SignWithMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStateHashMockSignWithInput{p}, "NodeStateHash.SignWith got unexpected parameters")

		result := m.SignWithMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.SignWith")
			return
		}

		r = result.r

		return
	}

	if m.SignWithMock.mainExpectation != nil {

		input := m.SignWithMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStateHashMockSignWithInput{p}, "NodeStateHash.SignWith got unexpected parameters")
		}

		result := m.SignWithMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.SignWith")
		}

		r = result.r

		return
	}

	if m.SignWithFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.SignWith. %v", p)
		return
	}

	return m.SignWithFunc(p)
}

//SignWithMinimockCounter returns a count of NodeStateHashMock.SignWithFunc invocations
func (m *NodeStateHashMock) SignWithMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithCounter)
}

//SignWithMinimockPreCounter returns the value of NodeStateHashMock.SignWith invocations
func (m *NodeStateHashMock) SignWithMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithPreCounter)
}

//SignWithFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) SignWithFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SignWithMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SignWithCounter) == uint64(len(m.SignWithMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SignWithMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SignWithCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SignWithFunc != nil {
		return atomic.LoadUint64(&m.SignWithCounter) > 0
	}

	return true
}

type mNodeStateHashMockWriteTo struct {
	mock              *NodeStateHashMock
	mainExpectation   *NodeStateHashMockWriteToExpectation
	expectationSeries []*NodeStateHashMockWriteToExpectation
}

type NodeStateHashMockWriteToExpectation struct {
	input  *NodeStateHashMockWriteToInput
	result *NodeStateHashMockWriteToResult
}

type NodeStateHashMockWriteToInput struct {
	p io.Writer
}

type NodeStateHashMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of NodeStateHash.WriteTo is expected from 1 to Infinity times
func (m *mNodeStateHashMockWriteTo) Expect(p io.Writer) *mNodeStateHashMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockWriteToExpectation{}
	}
	m.mainExpectation.input = &NodeStateHashMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of NodeStateHash.WriteTo
func (m *mNodeStateHashMockWriteTo) Return(r int64, r1 error) *NodeStateHashMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &NodeStateHashMockWriteToExpectation{}
	}
	m.mainExpectation.result = &NodeStateHashMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of NodeStateHash.WriteTo is expected once
func (m *mNodeStateHashMockWriteTo) ExpectOnce(p io.Writer) *NodeStateHashMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &NodeStateHashMockWriteToExpectation{}
	expectation.input = &NodeStateHashMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *NodeStateHashMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &NodeStateHashMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of NodeStateHash.WriteTo method
func (m *mNodeStateHashMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *NodeStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.NodeStateHash interface
func (m *NodeStateHashMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to NodeStateHashMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, NodeStateHashMockWriteToInput{p}, "NodeStateHash.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, NodeStateHashMockWriteToInput{p}, "NodeStateHash.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the NodeStateHashMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to NodeStateHashMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of NodeStateHashMock.WriteToFunc invocations
func (m *NodeStateHashMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of NodeStateHashMock.WriteTo invocations
func (m *NodeStateHashMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *NodeStateHashMock) WriteToFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.WriteToMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.WriteToCounter) == uint64(len(m.WriteToMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.WriteToMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.WriteToCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.WriteToFunc != nil {
		return atomic.LoadUint64(&m.WriteToCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeStateHashMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to NodeStateHashMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *NodeStateHashMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *NodeStateHashMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *NodeStateHashMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to NodeStateHashMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to NodeStateHashMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *NodeStateHashMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *NodeStateHashMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AsByteStringFinished()
		ok = ok && m.AsBytesFinished()
		ok = ok && m.CopyOfDigestFinished()
		ok = ok && m.EqualsFinished()
		ok = ok && m.FixedByteSizeFinished()
		ok = ok && m.FoldToUint64Finished()
		ok = ok && m.GetDigestMethodFinished()
		ok = ok && m.ReadFinished()
		ok = ok && m.SignWithFinished()
		ok = ok && m.WriteToFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AsByteStringFinished() {
				m.t.Error("Expected call to NodeStateHashMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to NodeStateHashMock.AsBytes")
			}

			if !m.CopyOfDigestFinished() {
				m.t.Error("Expected call to NodeStateHashMock.CopyOfDigest")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to NodeStateHashMock.Equals")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to NodeStateHashMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to NodeStateHashMock.FoldToUint64")
			}

			if !m.GetDigestMethodFinished() {
				m.t.Error("Expected call to NodeStateHashMock.GetDigestMethod")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to NodeStateHashMock.Read")
			}

			if !m.SignWithFinished() {
				m.t.Error("Expected call to NodeStateHashMock.SignWith")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to NodeStateHashMock.WriteTo")
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
func (m *NodeStateHashMock) AllMocksCalled() bool {

	if !m.AsByteStringFinished() {
		return false
	}

	if !m.AsBytesFinished() {
		return false
	}

	if !m.CopyOfDigestFinished() {
		return false
	}

	if !m.EqualsFinished() {
		return false
	}

	if !m.FixedByteSizeFinished() {
		return false
	}

	if !m.FoldToUint64Finished() {
		return false
	}

	if !m.GetDigestMethodFinished() {
		return false
	}

	if !m.ReadFinished() {
		return false
	}

	if !m.SignWithFinished() {
		return false
	}

	if !m.WriteToFinished() {
		return false
	}

	return true
}
