package proofs

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "GlobulaStateHash" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/proofs
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"

	testify_assert "github.com/stretchr/testify/assert"
)

//GlobulaStateHashMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash
type GlobulaStateHashMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mGlobulaStateHashMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mGlobulaStateHashMockAsBytes

	CopyOfDigestFunc       func() (r cryptkit.Digest)
	CopyOfDigestCounter    uint64
	CopyOfDigestPreCounter uint64
	CopyOfDigestMock       mGlobulaStateHashMockCopyOfDigest

	EqualsFunc       func(p cryptkit.DigestHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mGlobulaStateHashMockEquals

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mGlobulaStateHashMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mGlobulaStateHashMockFoldToUint64

	GetDigestMethodFunc       func() (r cryptkit.DigestMethod)
	GetDigestMethodCounter    uint64
	GetDigestMethodPreCounter uint64
	GetDigestMethodMock       mGlobulaStateHashMockGetDigestMethod

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mGlobulaStateHashMockRead

	SignWithFunc       func(p cryptkit.DigestSigner) (r cryptkit.SignedDigest)
	SignWithCounter    uint64
	SignWithPreCounter uint64
	SignWithMock       mGlobulaStateHashMockSignWith

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mGlobulaStateHashMockWriteTo
}

//NewGlobulaStateHashMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash
func NewGlobulaStateHashMock(t minimock.Tester) *GlobulaStateHashMock {
	m := &GlobulaStateHashMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mGlobulaStateHashMockAsByteString{mock: m}
	m.AsBytesMock = mGlobulaStateHashMockAsBytes{mock: m}
	m.CopyOfDigestMock = mGlobulaStateHashMockCopyOfDigest{mock: m}
	m.EqualsMock = mGlobulaStateHashMockEquals{mock: m}
	m.FixedByteSizeMock = mGlobulaStateHashMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mGlobulaStateHashMockFoldToUint64{mock: m}
	m.GetDigestMethodMock = mGlobulaStateHashMockGetDigestMethod{mock: m}
	m.ReadMock = mGlobulaStateHashMockRead{mock: m}
	m.SignWithMock = mGlobulaStateHashMockSignWith{mock: m}
	m.WriteToMock = mGlobulaStateHashMockWriteTo{mock: m}

	return m
}

type mGlobulaStateHashMockAsByteString struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockAsByteStringExpectation
	expectationSeries []*GlobulaStateHashMockAsByteStringExpectation
}

type GlobulaStateHashMockAsByteStringExpectation struct {
	result *GlobulaStateHashMockAsByteStringResult
}

type GlobulaStateHashMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of GlobulaStateHash.AsByteString is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockAsByteString) Expect() *mGlobulaStateHashMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of GlobulaStateHash.AsByteString
func (m *mGlobulaStateHashMockAsByteString) Return(r string) *GlobulaStateHashMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.AsByteString is expected once
func (m *mGlobulaStateHashMockAsByteString) ExpectOnce() *GlobulaStateHashMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockAsByteStringExpectation) Return(r string) {
	e.result = &GlobulaStateHashMockAsByteStringResult{r}
}

//Set uses given function f as a mock of GlobulaStateHash.AsByteString method
func (m *mGlobulaStateHashMockAsByteString) Set(f func() (r string)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of GlobulaStateHashMock.AsByteStringFunc invocations
func (m *GlobulaStateHashMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of GlobulaStateHashMock.AsByteString invocations
func (m *GlobulaStateHashMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) AsByteStringFinished() bool {
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

type mGlobulaStateHashMockAsBytes struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockAsBytesExpectation
	expectationSeries []*GlobulaStateHashMockAsBytesExpectation
}

type GlobulaStateHashMockAsBytesExpectation struct {
	result *GlobulaStateHashMockAsBytesResult
}

type GlobulaStateHashMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of GlobulaStateHash.AsBytes is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockAsBytes) Expect() *mGlobulaStateHashMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of GlobulaStateHash.AsBytes
func (m *mGlobulaStateHashMockAsBytes) Return(r []byte) *GlobulaStateHashMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.AsBytes is expected once
func (m *mGlobulaStateHashMockAsBytes) ExpectOnce() *GlobulaStateHashMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockAsBytesExpectation) Return(r []byte) {
	e.result = &GlobulaStateHashMockAsBytesResult{r}
}

//Set uses given function f as a mock of GlobulaStateHash.AsBytes method
func (m *mGlobulaStateHashMockAsBytes) Set(f func() (r []byte)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of GlobulaStateHashMock.AsBytesFunc invocations
func (m *GlobulaStateHashMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of GlobulaStateHashMock.AsBytes invocations
func (m *GlobulaStateHashMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) AsBytesFinished() bool {
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

type mGlobulaStateHashMockCopyOfDigest struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockCopyOfDigestExpectation
	expectationSeries []*GlobulaStateHashMockCopyOfDigestExpectation
}

type GlobulaStateHashMockCopyOfDigestExpectation struct {
	result *GlobulaStateHashMockCopyOfDigestResult
}

type GlobulaStateHashMockCopyOfDigestResult struct {
	r cryptkit.Digest
}

//Expect specifies that invocation of GlobulaStateHash.CopyOfDigest is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockCopyOfDigest) Expect() *mGlobulaStateHashMockCopyOfDigest {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockCopyOfDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of GlobulaStateHash.CopyOfDigest
func (m *mGlobulaStateHashMockCopyOfDigest) Return(r cryptkit.Digest) *GlobulaStateHashMock {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockCopyOfDigestExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockCopyOfDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.CopyOfDigest is expected once
func (m *mGlobulaStateHashMockCopyOfDigest) ExpectOnce() *GlobulaStateHashMockCopyOfDigestExpectation {
	m.mock.CopyOfDigestFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockCopyOfDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockCopyOfDigestExpectation) Return(r cryptkit.Digest) {
	e.result = &GlobulaStateHashMockCopyOfDigestResult{r}
}

//Set uses given function f as a mock of GlobulaStateHash.CopyOfDigest method
func (m *mGlobulaStateHashMockCopyOfDigest) Set(f func() (r cryptkit.Digest)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CopyOfDigestFunc = f
	return m.mock
}

//CopyOfDigest implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) CopyOfDigest() (r cryptkit.Digest) {
	counter := atomic.AddUint64(&m.CopyOfDigestPreCounter, 1)
	defer atomic.AddUint64(&m.CopyOfDigestCounter, 1)

	if len(m.CopyOfDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CopyOfDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.CopyOfDigest.")
			return
		}

		result := m.CopyOfDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.CopyOfDigest")
			return
		}

		r = result.r

		return
	}

	if m.CopyOfDigestMock.mainExpectation != nil {

		result := m.CopyOfDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.CopyOfDigest")
		}

		r = result.r

		return
	}

	if m.CopyOfDigestFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.CopyOfDigest.")
		return
	}

	return m.CopyOfDigestFunc()
}

//CopyOfDigestMinimockCounter returns a count of GlobulaStateHashMock.CopyOfDigestFunc invocations
func (m *GlobulaStateHashMock) CopyOfDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestCounter)
}

//CopyOfDigestMinimockPreCounter returns the value of GlobulaStateHashMock.CopyOfDigest invocations
func (m *GlobulaStateHashMock) CopyOfDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestPreCounter)
}

//CopyOfDigestFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) CopyOfDigestFinished() bool {
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

type mGlobulaStateHashMockEquals struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockEqualsExpectation
	expectationSeries []*GlobulaStateHashMockEqualsExpectation
}

type GlobulaStateHashMockEqualsExpectation struct {
	input  *GlobulaStateHashMockEqualsInput
	result *GlobulaStateHashMockEqualsResult
}

type GlobulaStateHashMockEqualsInput struct {
	p cryptkit.DigestHolder
}

type GlobulaStateHashMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of GlobulaStateHash.Equals is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockEquals) Expect(p cryptkit.DigestHolder) *mGlobulaStateHashMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockEqualsExpectation{}
	}
	m.mainExpectation.input = &GlobulaStateHashMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of GlobulaStateHash.Equals
func (m *mGlobulaStateHashMockEquals) Return(r bool) *GlobulaStateHashMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockEqualsExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.Equals is expected once
func (m *mGlobulaStateHashMockEquals) ExpectOnce(p cryptkit.DigestHolder) *GlobulaStateHashMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockEqualsExpectation{}
	expectation.input = &GlobulaStateHashMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockEqualsExpectation) Return(r bool) {
	e.result = &GlobulaStateHashMockEqualsResult{r}
}

//Set uses given function f as a mock of GlobulaStateHash.Equals method
func (m *mGlobulaStateHashMockEquals) Set(f func(p cryptkit.DigestHolder) (r bool)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) Equals(p cryptkit.DigestHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GlobulaStateHashMockEqualsInput{p}, "GlobulaStateHash.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GlobulaStateHashMockEqualsInput{p}, "GlobulaStateHash.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of GlobulaStateHashMock.EqualsFunc invocations
func (m *GlobulaStateHashMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of GlobulaStateHashMock.Equals invocations
func (m *GlobulaStateHashMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) EqualsFinished() bool {
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

type mGlobulaStateHashMockFixedByteSize struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockFixedByteSizeExpectation
	expectationSeries []*GlobulaStateHashMockFixedByteSizeExpectation
}

type GlobulaStateHashMockFixedByteSizeExpectation struct {
	result *GlobulaStateHashMockFixedByteSizeResult
}

type GlobulaStateHashMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of GlobulaStateHash.FixedByteSize is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockFixedByteSize) Expect() *mGlobulaStateHashMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of GlobulaStateHash.FixedByteSize
func (m *mGlobulaStateHashMockFixedByteSize) Return(r int) *GlobulaStateHashMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.FixedByteSize is expected once
func (m *mGlobulaStateHashMockFixedByteSize) ExpectOnce() *GlobulaStateHashMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockFixedByteSizeExpectation) Return(r int) {
	e.result = &GlobulaStateHashMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of GlobulaStateHash.FixedByteSize method
func (m *mGlobulaStateHashMockFixedByteSize) Set(f func() (r int)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of GlobulaStateHashMock.FixedByteSizeFunc invocations
func (m *GlobulaStateHashMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of GlobulaStateHashMock.FixedByteSize invocations
func (m *GlobulaStateHashMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) FixedByteSizeFinished() bool {
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

type mGlobulaStateHashMockFoldToUint64 struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockFoldToUint64Expectation
	expectationSeries []*GlobulaStateHashMockFoldToUint64Expectation
}

type GlobulaStateHashMockFoldToUint64Expectation struct {
	result *GlobulaStateHashMockFoldToUint64Result
}

type GlobulaStateHashMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of GlobulaStateHash.FoldToUint64 is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockFoldToUint64) Expect() *mGlobulaStateHashMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of GlobulaStateHash.FoldToUint64
func (m *mGlobulaStateHashMockFoldToUint64) Return(r uint64) *GlobulaStateHashMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.FoldToUint64 is expected once
func (m *mGlobulaStateHashMockFoldToUint64) ExpectOnce() *GlobulaStateHashMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &GlobulaStateHashMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of GlobulaStateHash.FoldToUint64 method
func (m *mGlobulaStateHashMockFoldToUint64) Set(f func() (r uint64)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of GlobulaStateHashMock.FoldToUint64Func invocations
func (m *GlobulaStateHashMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of GlobulaStateHashMock.FoldToUint64 invocations
func (m *GlobulaStateHashMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) FoldToUint64Finished() bool {
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

type mGlobulaStateHashMockGetDigestMethod struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockGetDigestMethodExpectation
	expectationSeries []*GlobulaStateHashMockGetDigestMethodExpectation
}

type GlobulaStateHashMockGetDigestMethodExpectation struct {
	result *GlobulaStateHashMockGetDigestMethodResult
}

type GlobulaStateHashMockGetDigestMethodResult struct {
	r cryptkit.DigestMethod
}

//Expect specifies that invocation of GlobulaStateHash.GetDigestMethod is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockGetDigestMethod) Expect() *mGlobulaStateHashMockGetDigestMethod {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockGetDigestMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of GlobulaStateHash.GetDigestMethod
func (m *mGlobulaStateHashMockGetDigestMethod) Return(r cryptkit.DigestMethod) *GlobulaStateHashMock {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockGetDigestMethodExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockGetDigestMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.GetDigestMethod is expected once
func (m *mGlobulaStateHashMockGetDigestMethod) ExpectOnce() *GlobulaStateHashMockGetDigestMethodExpectation {
	m.mock.GetDigestMethodFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockGetDigestMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockGetDigestMethodExpectation) Return(r cryptkit.DigestMethod) {
	e.result = &GlobulaStateHashMockGetDigestMethodResult{r}
}

//Set uses given function f as a mock of GlobulaStateHash.GetDigestMethod method
func (m *mGlobulaStateHashMockGetDigestMethod) Set(f func() (r cryptkit.DigestMethod)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDigestMethodFunc = f
	return m.mock
}

//GetDigestMethod implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) GetDigestMethod() (r cryptkit.DigestMethod) {
	counter := atomic.AddUint64(&m.GetDigestMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetDigestMethodCounter, 1)

	if len(m.GetDigestMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDigestMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.GetDigestMethod.")
			return
		}

		result := m.GetDigestMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.GetDigestMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetDigestMethodMock.mainExpectation != nil {

		result := m.GetDigestMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.GetDigestMethod")
		}

		r = result.r

		return
	}

	if m.GetDigestMethodFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.GetDigestMethod.")
		return
	}

	return m.GetDigestMethodFunc()
}

//GetDigestMethodMinimockCounter returns a count of GlobulaStateHashMock.GetDigestMethodFunc invocations
func (m *GlobulaStateHashMock) GetDigestMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodCounter)
}

//GetDigestMethodMinimockPreCounter returns the value of GlobulaStateHashMock.GetDigestMethod invocations
func (m *GlobulaStateHashMock) GetDigestMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodPreCounter)
}

//GetDigestMethodFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) GetDigestMethodFinished() bool {
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

type mGlobulaStateHashMockRead struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockReadExpectation
	expectationSeries []*GlobulaStateHashMockReadExpectation
}

type GlobulaStateHashMockReadExpectation struct {
	input  *GlobulaStateHashMockReadInput
	result *GlobulaStateHashMockReadResult
}

type GlobulaStateHashMockReadInput struct {
	p []byte
}

type GlobulaStateHashMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of GlobulaStateHash.Read is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockRead) Expect(p []byte) *mGlobulaStateHashMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockReadExpectation{}
	}
	m.mainExpectation.input = &GlobulaStateHashMockReadInput{p}
	return m
}

//Return specifies results of invocation of GlobulaStateHash.Read
func (m *mGlobulaStateHashMockRead) Return(r int, r1 error) *GlobulaStateHashMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockReadExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.Read is expected once
func (m *mGlobulaStateHashMockRead) ExpectOnce(p []byte) *GlobulaStateHashMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockReadExpectation{}
	expectation.input = &GlobulaStateHashMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockReadExpectation) Return(r int, r1 error) {
	e.result = &GlobulaStateHashMockReadResult{r, r1}
}

//Set uses given function f as a mock of GlobulaStateHash.Read method
func (m *mGlobulaStateHashMockRead) Set(f func(p []byte) (r int, r1 error)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GlobulaStateHashMockReadInput{p}, "GlobulaStateHash.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GlobulaStateHashMockReadInput{p}, "GlobulaStateHash.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of GlobulaStateHashMock.ReadFunc invocations
func (m *GlobulaStateHashMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of GlobulaStateHashMock.Read invocations
func (m *GlobulaStateHashMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) ReadFinished() bool {
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

type mGlobulaStateHashMockSignWith struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockSignWithExpectation
	expectationSeries []*GlobulaStateHashMockSignWithExpectation
}

type GlobulaStateHashMockSignWithExpectation struct {
	input  *GlobulaStateHashMockSignWithInput
	result *GlobulaStateHashMockSignWithResult
}

type GlobulaStateHashMockSignWithInput struct {
	p cryptkit.DigestSigner
}

type GlobulaStateHashMockSignWithResult struct {
	r cryptkit.SignedDigest
}

//Expect specifies that invocation of GlobulaStateHash.SignWith is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockSignWith) Expect(p cryptkit.DigestSigner) *mGlobulaStateHashMockSignWith {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockSignWithExpectation{}
	}
	m.mainExpectation.input = &GlobulaStateHashMockSignWithInput{p}
	return m
}

//Return specifies results of invocation of GlobulaStateHash.SignWith
func (m *mGlobulaStateHashMockSignWith) Return(r cryptkit.SignedDigest) *GlobulaStateHashMock {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockSignWithExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockSignWithResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.SignWith is expected once
func (m *mGlobulaStateHashMockSignWith) ExpectOnce(p cryptkit.DigestSigner) *GlobulaStateHashMockSignWithExpectation {
	m.mock.SignWithFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockSignWithExpectation{}
	expectation.input = &GlobulaStateHashMockSignWithInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockSignWithExpectation) Return(r cryptkit.SignedDigest) {
	e.result = &GlobulaStateHashMockSignWithResult{r}
}

//Set uses given function f as a mock of GlobulaStateHash.SignWith method
func (m *mGlobulaStateHashMockSignWith) Set(f func(p cryptkit.DigestSigner) (r cryptkit.SignedDigest)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignWithFunc = f
	return m.mock
}

//SignWith implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) SignWith(p cryptkit.DigestSigner) (r cryptkit.SignedDigest) {
	counter := atomic.AddUint64(&m.SignWithPreCounter, 1)
	defer atomic.AddUint64(&m.SignWithCounter, 1)

	if len(m.SignWithMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignWithMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.SignWith. %v", p)
			return
		}

		input := m.SignWithMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GlobulaStateHashMockSignWithInput{p}, "GlobulaStateHash.SignWith got unexpected parameters")

		result := m.SignWithMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.SignWith")
			return
		}

		r = result.r

		return
	}

	if m.SignWithMock.mainExpectation != nil {

		input := m.SignWithMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GlobulaStateHashMockSignWithInput{p}, "GlobulaStateHash.SignWith got unexpected parameters")
		}

		result := m.SignWithMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.SignWith")
		}

		r = result.r

		return
	}

	if m.SignWithFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.SignWith. %v", p)
		return
	}

	return m.SignWithFunc(p)
}

//SignWithMinimockCounter returns a count of GlobulaStateHashMock.SignWithFunc invocations
func (m *GlobulaStateHashMock) SignWithMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithCounter)
}

//SignWithMinimockPreCounter returns the value of GlobulaStateHashMock.SignWith invocations
func (m *GlobulaStateHashMock) SignWithMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithPreCounter)
}

//SignWithFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) SignWithFinished() bool {
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

type mGlobulaStateHashMockWriteTo struct {
	mock              *GlobulaStateHashMock
	mainExpectation   *GlobulaStateHashMockWriteToExpectation
	expectationSeries []*GlobulaStateHashMockWriteToExpectation
}

type GlobulaStateHashMockWriteToExpectation struct {
	input  *GlobulaStateHashMockWriteToInput
	result *GlobulaStateHashMockWriteToResult
}

type GlobulaStateHashMockWriteToInput struct {
	p io.Writer
}

type GlobulaStateHashMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of GlobulaStateHash.WriteTo is expected from 1 to Infinity times
func (m *mGlobulaStateHashMockWriteTo) Expect(p io.Writer) *mGlobulaStateHashMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockWriteToExpectation{}
	}
	m.mainExpectation.input = &GlobulaStateHashMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of GlobulaStateHash.WriteTo
func (m *mGlobulaStateHashMockWriteTo) Return(r int64, r1 error) *GlobulaStateHashMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &GlobulaStateHashMockWriteToExpectation{}
	}
	m.mainExpectation.result = &GlobulaStateHashMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of GlobulaStateHash.WriteTo is expected once
func (m *mGlobulaStateHashMockWriteTo) ExpectOnce(p io.Writer) *GlobulaStateHashMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &GlobulaStateHashMockWriteToExpectation{}
	expectation.input = &GlobulaStateHashMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *GlobulaStateHashMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &GlobulaStateHashMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of GlobulaStateHash.WriteTo method
func (m *mGlobulaStateHashMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *GlobulaStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.GlobulaStateHash interface
func (m *GlobulaStateHashMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to GlobulaStateHashMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, GlobulaStateHashMockWriteToInput{p}, "GlobulaStateHash.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, GlobulaStateHashMockWriteToInput{p}, "GlobulaStateHash.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the GlobulaStateHashMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to GlobulaStateHashMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of GlobulaStateHashMock.WriteToFunc invocations
func (m *GlobulaStateHashMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of GlobulaStateHashMock.WriteTo invocations
func (m *GlobulaStateHashMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *GlobulaStateHashMock) WriteToFinished() bool {
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
func (m *GlobulaStateHashMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *GlobulaStateHashMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *GlobulaStateHashMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *GlobulaStateHashMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to GlobulaStateHashMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *GlobulaStateHashMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *GlobulaStateHashMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to GlobulaStateHashMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.AsBytes")
			}

			if !m.CopyOfDigestFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.CopyOfDigest")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.Equals")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to GlobulaStateHashMock.FoldToUint64")
			}

			if !m.GetDigestMethodFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.GetDigestMethod")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.Read")
			}

			if !m.SignWithFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.SignWith")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to GlobulaStateHashMock.WriteTo")
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
func (m *GlobulaStateHashMock) AllMocksCalled() bool {

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
