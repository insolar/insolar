package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DigestHolder" can be found in github.com/insolar/insolar/network/consensus/common
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//DigestHolderMock implements github.com/insolar/insolar/network/consensus/common.DigestHolder
type DigestHolderMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mDigestHolderMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mDigestHolderMockAsBytes

	CopyOfDigestFunc       func() (r Digest)
	CopyOfDigestCounter    uint64
	CopyOfDigestPreCounter uint64
	CopyOfDigestMock       mDigestHolderMockCopyOfDigest

	EqualsFunc       func(p DigestHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mDigestHolderMockEquals

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mDigestHolderMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mDigestHolderMockFoldToUint64

	GetDigestMethodFunc       func() (r DigestMethod)
	GetDigestMethodCounter    uint64
	GetDigestMethodPreCounter uint64
	GetDigestMethodMock       mDigestHolderMockGetDigestMethod

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mDigestHolderMockRead

	SignWithFunc       func(p DigestSigner) (r SignedDigest)
	SignWithCounter    uint64
	SignWithPreCounter uint64
	SignWithMock       mDigestHolderMockSignWith

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mDigestHolderMockWriteTo
}

//NewDigestHolderMock returns a mock for github.com/insolar/insolar/network/consensus/common.DigestHolder
func NewDigestHolderMock(t minimock.Tester) *DigestHolderMock {
	m := &DigestHolderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mDigestHolderMockAsByteString{mock: m}
	m.AsBytesMock = mDigestHolderMockAsBytes{mock: m}
	m.CopyOfDigestMock = mDigestHolderMockCopyOfDigest{mock: m}
	m.EqualsMock = mDigestHolderMockEquals{mock: m}
	m.FixedByteSizeMock = mDigestHolderMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mDigestHolderMockFoldToUint64{mock: m}
	m.GetDigestMethodMock = mDigestHolderMockGetDigestMethod{mock: m}
	m.ReadMock = mDigestHolderMockRead{mock: m}
	m.SignWithMock = mDigestHolderMockSignWith{mock: m}
	m.WriteToMock = mDigestHolderMockWriteTo{mock: m}

	return m
}

type mDigestHolderMockAsByteString struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockAsByteStringExpectation
	expectationSeries []*DigestHolderMockAsByteStringExpectation
}

type DigestHolderMockAsByteStringExpectation struct {
	result *DigestHolderMockAsByteStringResult
}

type DigestHolderMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of DigestHolder.AsByteString is expected from 1 to Infinity times
func (m *mDigestHolderMockAsByteString) Expect() *mDigestHolderMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of DigestHolder.AsByteString
func (m *mDigestHolderMockAsByteString) Return(r string) *DigestHolderMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.AsByteString is expected once
func (m *mDigestHolderMockAsByteString) ExpectOnce() *DigestHolderMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockAsByteStringExpectation) Return(r string) {
	e.result = &DigestHolderMockAsByteStringResult{r}
}

//Set uses given function f as a mock of DigestHolder.AsByteString method
func (m *mDigestHolderMockAsByteString) Set(f func() (r string)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of DigestHolderMock.AsByteStringFunc invocations
func (m *DigestHolderMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of DigestHolderMock.AsByteString invocations
func (m *DigestHolderMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) AsByteStringFinished() bool {
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

type mDigestHolderMockAsBytes struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockAsBytesExpectation
	expectationSeries []*DigestHolderMockAsBytesExpectation
}

type DigestHolderMockAsBytesExpectation struct {
	result *DigestHolderMockAsBytesResult
}

type DigestHolderMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of DigestHolder.AsBytes is expected from 1 to Infinity times
func (m *mDigestHolderMockAsBytes) Expect() *mDigestHolderMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of DigestHolder.AsBytes
func (m *mDigestHolderMockAsBytes) Return(r []byte) *DigestHolderMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.AsBytes is expected once
func (m *mDigestHolderMockAsBytes) ExpectOnce() *DigestHolderMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockAsBytesExpectation) Return(r []byte) {
	e.result = &DigestHolderMockAsBytesResult{r}
}

//Set uses given function f as a mock of DigestHolder.AsBytes method
func (m *mDigestHolderMockAsBytes) Set(f func() (r []byte)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of DigestHolderMock.AsBytesFunc invocations
func (m *DigestHolderMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of DigestHolderMock.AsBytes invocations
func (m *DigestHolderMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) AsBytesFinished() bool {
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

type mDigestHolderMockCopyOfDigest struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockCopyOfDigestExpectation
	expectationSeries []*DigestHolderMockCopyOfDigestExpectation
}

type DigestHolderMockCopyOfDigestExpectation struct {
	result *DigestHolderMockCopyOfDigestResult
}

type DigestHolderMockCopyOfDigestResult struct {
	r Digest
}

//Expect specifies that invocation of DigestHolder.CopyOfDigest is expected from 1 to Infinity times
func (m *mDigestHolderMockCopyOfDigest) Expect() *mDigestHolderMockCopyOfDigest {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockCopyOfDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of DigestHolder.CopyOfDigest
func (m *mDigestHolderMockCopyOfDigest) Return(r Digest) *DigestHolderMock {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockCopyOfDigestExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockCopyOfDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.CopyOfDigest is expected once
func (m *mDigestHolderMockCopyOfDigest) ExpectOnce() *DigestHolderMockCopyOfDigestExpectation {
	m.mock.CopyOfDigestFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockCopyOfDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockCopyOfDigestExpectation) Return(r Digest) {
	e.result = &DigestHolderMockCopyOfDigestResult{r}
}

//Set uses given function f as a mock of DigestHolder.CopyOfDigest method
func (m *mDigestHolderMockCopyOfDigest) Set(f func() (r Digest)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CopyOfDigestFunc = f
	return m.mock
}

//CopyOfDigest implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) CopyOfDigest() (r Digest) {
	counter := atomic.AddUint64(&m.CopyOfDigestPreCounter, 1)
	defer atomic.AddUint64(&m.CopyOfDigestCounter, 1)

	if len(m.CopyOfDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CopyOfDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.CopyOfDigest.")
			return
		}

		result := m.CopyOfDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.CopyOfDigest")
			return
		}

		r = result.r

		return
	}

	if m.CopyOfDigestMock.mainExpectation != nil {

		result := m.CopyOfDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.CopyOfDigest")
		}

		r = result.r

		return
	}

	if m.CopyOfDigestFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.CopyOfDigest.")
		return
	}

	return m.CopyOfDigestFunc()
}

//CopyOfDigestMinimockCounter returns a count of DigestHolderMock.CopyOfDigestFunc invocations
func (m *DigestHolderMock) CopyOfDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestCounter)
}

//CopyOfDigestMinimockPreCounter returns the value of DigestHolderMock.CopyOfDigest invocations
func (m *DigestHolderMock) CopyOfDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestPreCounter)
}

//CopyOfDigestFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) CopyOfDigestFinished() bool {
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

type mDigestHolderMockEquals struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockEqualsExpectation
	expectationSeries []*DigestHolderMockEqualsExpectation
}

type DigestHolderMockEqualsExpectation struct {
	input  *DigestHolderMockEqualsInput
	result *DigestHolderMockEqualsResult
}

type DigestHolderMockEqualsInput struct {
	p DigestHolder
}

type DigestHolderMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of DigestHolder.Equals is expected from 1 to Infinity times
func (m *mDigestHolderMockEquals) Expect(p DigestHolder) *mDigestHolderMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockEqualsExpectation{}
	}
	m.mainExpectation.input = &DigestHolderMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of DigestHolder.Equals
func (m *mDigestHolderMockEquals) Return(r bool) *DigestHolderMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockEqualsExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.Equals is expected once
func (m *mDigestHolderMockEquals) ExpectOnce(p DigestHolder) *DigestHolderMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockEqualsExpectation{}
	expectation.input = &DigestHolderMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockEqualsExpectation) Return(r bool) {
	e.result = &DigestHolderMockEqualsResult{r}
}

//Set uses given function f as a mock of DigestHolder.Equals method
func (m *mDigestHolderMockEquals) Set(f func(p DigestHolder) (r bool)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) Equals(p DigestHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DigestHolderMockEqualsInput{p}, "DigestHolder.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DigestHolderMockEqualsInput{p}, "DigestHolder.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of DigestHolderMock.EqualsFunc invocations
func (m *DigestHolderMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of DigestHolderMock.Equals invocations
func (m *DigestHolderMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) EqualsFinished() bool {
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

type mDigestHolderMockFixedByteSize struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockFixedByteSizeExpectation
	expectationSeries []*DigestHolderMockFixedByteSizeExpectation
}

type DigestHolderMockFixedByteSizeExpectation struct {
	result *DigestHolderMockFixedByteSizeResult
}

type DigestHolderMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of DigestHolder.FixedByteSize is expected from 1 to Infinity times
func (m *mDigestHolderMockFixedByteSize) Expect() *mDigestHolderMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of DigestHolder.FixedByteSize
func (m *mDigestHolderMockFixedByteSize) Return(r int) *DigestHolderMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.FixedByteSize is expected once
func (m *mDigestHolderMockFixedByteSize) ExpectOnce() *DigestHolderMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockFixedByteSizeExpectation) Return(r int) {
	e.result = &DigestHolderMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of DigestHolder.FixedByteSize method
func (m *mDigestHolderMockFixedByteSize) Set(f func() (r int)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of DigestHolderMock.FixedByteSizeFunc invocations
func (m *DigestHolderMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of DigestHolderMock.FixedByteSize invocations
func (m *DigestHolderMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) FixedByteSizeFinished() bool {
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

type mDigestHolderMockFoldToUint64 struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockFoldToUint64Expectation
	expectationSeries []*DigestHolderMockFoldToUint64Expectation
}

type DigestHolderMockFoldToUint64Expectation struct {
	result *DigestHolderMockFoldToUint64Result
}

type DigestHolderMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of DigestHolder.FoldToUint64 is expected from 1 to Infinity times
func (m *mDigestHolderMockFoldToUint64) Expect() *mDigestHolderMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of DigestHolder.FoldToUint64
func (m *mDigestHolderMockFoldToUint64) Return(r uint64) *DigestHolderMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &DigestHolderMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.FoldToUint64 is expected once
func (m *mDigestHolderMockFoldToUint64) ExpectOnce() *DigestHolderMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &DigestHolderMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of DigestHolder.FoldToUint64 method
func (m *mDigestHolderMockFoldToUint64) Set(f func() (r uint64)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of DigestHolderMock.FoldToUint64Func invocations
func (m *DigestHolderMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of DigestHolderMock.FoldToUint64 invocations
func (m *DigestHolderMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *DigestHolderMock) FoldToUint64Finished() bool {
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

type mDigestHolderMockGetDigestMethod struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockGetDigestMethodExpectation
	expectationSeries []*DigestHolderMockGetDigestMethodExpectation
}

type DigestHolderMockGetDigestMethodExpectation struct {
	result *DigestHolderMockGetDigestMethodResult
}

type DigestHolderMockGetDigestMethodResult struct {
	r DigestMethod
}

//Expect specifies that invocation of DigestHolder.GetDigestMethod is expected from 1 to Infinity times
func (m *mDigestHolderMockGetDigestMethod) Expect() *mDigestHolderMockGetDigestMethod {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockGetDigestMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of DigestHolder.GetDigestMethod
func (m *mDigestHolderMockGetDigestMethod) Return(r DigestMethod) *DigestHolderMock {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockGetDigestMethodExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockGetDigestMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.GetDigestMethod is expected once
func (m *mDigestHolderMockGetDigestMethod) ExpectOnce() *DigestHolderMockGetDigestMethodExpectation {
	m.mock.GetDigestMethodFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockGetDigestMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockGetDigestMethodExpectation) Return(r DigestMethod) {
	e.result = &DigestHolderMockGetDigestMethodResult{r}
}

//Set uses given function f as a mock of DigestHolder.GetDigestMethod method
func (m *mDigestHolderMockGetDigestMethod) Set(f func() (r DigestMethod)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDigestMethodFunc = f
	return m.mock
}

//GetDigestMethod implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) GetDigestMethod() (r DigestMethod) {
	counter := atomic.AddUint64(&m.GetDigestMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetDigestMethodCounter, 1)

	if len(m.GetDigestMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDigestMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.GetDigestMethod.")
			return
		}

		result := m.GetDigestMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.GetDigestMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetDigestMethodMock.mainExpectation != nil {

		result := m.GetDigestMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.GetDigestMethod")
		}

		r = result.r

		return
	}

	if m.GetDigestMethodFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.GetDigestMethod.")
		return
	}

	return m.GetDigestMethodFunc()
}

//GetDigestMethodMinimockCounter returns a count of DigestHolderMock.GetDigestMethodFunc invocations
func (m *DigestHolderMock) GetDigestMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodCounter)
}

//GetDigestMethodMinimockPreCounter returns the value of DigestHolderMock.GetDigestMethod invocations
func (m *DigestHolderMock) GetDigestMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodPreCounter)
}

//GetDigestMethodFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) GetDigestMethodFinished() bool {
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

type mDigestHolderMockRead struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockReadExpectation
	expectationSeries []*DigestHolderMockReadExpectation
}

type DigestHolderMockReadExpectation struct {
	input  *DigestHolderMockReadInput
	result *DigestHolderMockReadResult
}

type DigestHolderMockReadInput struct {
	p []byte
}

type DigestHolderMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of DigestHolder.Read is expected from 1 to Infinity times
func (m *mDigestHolderMockRead) Expect(p []byte) *mDigestHolderMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockReadExpectation{}
	}
	m.mainExpectation.input = &DigestHolderMockReadInput{p}
	return m
}

//Return specifies results of invocation of DigestHolder.Read
func (m *mDigestHolderMockRead) Return(r int, r1 error) *DigestHolderMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockReadExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.Read is expected once
func (m *mDigestHolderMockRead) ExpectOnce(p []byte) *DigestHolderMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockReadExpectation{}
	expectation.input = &DigestHolderMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockReadExpectation) Return(r int, r1 error) {
	e.result = &DigestHolderMockReadResult{r, r1}
}

//Set uses given function f as a mock of DigestHolder.Read method
func (m *mDigestHolderMockRead) Set(f func(p []byte) (r int, r1 error)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DigestHolderMockReadInput{p}, "DigestHolder.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DigestHolderMockReadInput{p}, "DigestHolder.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of DigestHolderMock.ReadFunc invocations
func (m *DigestHolderMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of DigestHolderMock.Read invocations
func (m *DigestHolderMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) ReadFinished() bool {
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

type mDigestHolderMockSignWith struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockSignWithExpectation
	expectationSeries []*DigestHolderMockSignWithExpectation
}

type DigestHolderMockSignWithExpectation struct {
	input  *DigestHolderMockSignWithInput
	result *DigestHolderMockSignWithResult
}

type DigestHolderMockSignWithInput struct {
	p DigestSigner
}

type DigestHolderMockSignWithResult struct {
	r SignedDigest
}

//Expect specifies that invocation of DigestHolder.SignWith is expected from 1 to Infinity times
func (m *mDigestHolderMockSignWith) Expect(p DigestSigner) *mDigestHolderMockSignWith {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockSignWithExpectation{}
	}
	m.mainExpectation.input = &DigestHolderMockSignWithInput{p}
	return m
}

//Return specifies results of invocation of DigestHolder.SignWith
func (m *mDigestHolderMockSignWith) Return(r SignedDigest) *DigestHolderMock {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockSignWithExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockSignWithResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.SignWith is expected once
func (m *mDigestHolderMockSignWith) ExpectOnce(p DigestSigner) *DigestHolderMockSignWithExpectation {
	m.mock.SignWithFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockSignWithExpectation{}
	expectation.input = &DigestHolderMockSignWithInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockSignWithExpectation) Return(r SignedDigest) {
	e.result = &DigestHolderMockSignWithResult{r}
}

//Set uses given function f as a mock of DigestHolder.SignWith method
func (m *mDigestHolderMockSignWith) Set(f func(p DigestSigner) (r SignedDigest)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignWithFunc = f
	return m.mock
}

//SignWith implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) SignWith(p DigestSigner) (r SignedDigest) {
	counter := atomic.AddUint64(&m.SignWithPreCounter, 1)
	defer atomic.AddUint64(&m.SignWithCounter, 1)

	if len(m.SignWithMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignWithMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.SignWith. %v", p)
			return
		}

		input := m.SignWithMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DigestHolderMockSignWithInput{p}, "DigestHolder.SignWith got unexpected parameters")

		result := m.SignWithMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.SignWith")
			return
		}

		r = result.r

		return
	}

	if m.SignWithMock.mainExpectation != nil {

		input := m.SignWithMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DigestHolderMockSignWithInput{p}, "DigestHolder.SignWith got unexpected parameters")
		}

		result := m.SignWithMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.SignWith")
		}

		r = result.r

		return
	}

	if m.SignWithFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.SignWith. %v", p)
		return
	}

	return m.SignWithFunc(p)
}

//SignWithMinimockCounter returns a count of DigestHolderMock.SignWithFunc invocations
func (m *DigestHolderMock) SignWithMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithCounter)
}

//SignWithMinimockPreCounter returns the value of DigestHolderMock.SignWith invocations
func (m *DigestHolderMock) SignWithMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithPreCounter)
}

//SignWithFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) SignWithFinished() bool {
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

type mDigestHolderMockWriteTo struct {
	mock              *DigestHolderMock
	mainExpectation   *DigestHolderMockWriteToExpectation
	expectationSeries []*DigestHolderMockWriteToExpectation
}

type DigestHolderMockWriteToExpectation struct {
	input  *DigestHolderMockWriteToInput
	result *DigestHolderMockWriteToResult
}

type DigestHolderMockWriteToInput struct {
	p io.Writer
}

type DigestHolderMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of DigestHolder.WriteTo is expected from 1 to Infinity times
func (m *mDigestHolderMockWriteTo) Expect(p io.Writer) *mDigestHolderMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockWriteToExpectation{}
	}
	m.mainExpectation.input = &DigestHolderMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of DigestHolder.WriteTo
func (m *mDigestHolderMockWriteTo) Return(r int64, r1 error) *DigestHolderMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DigestHolderMockWriteToExpectation{}
	}
	m.mainExpectation.result = &DigestHolderMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DigestHolder.WriteTo is expected once
func (m *mDigestHolderMockWriteTo) ExpectOnce(p io.Writer) *DigestHolderMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &DigestHolderMockWriteToExpectation{}
	expectation.input = &DigestHolderMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DigestHolderMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &DigestHolderMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of DigestHolder.WriteTo method
func (m *mDigestHolderMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *DigestHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/common.DigestHolder interface
func (m *DigestHolderMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DigestHolderMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DigestHolderMockWriteToInput{p}, "DigestHolder.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DigestHolderMockWriteToInput{p}, "DigestHolder.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DigestHolderMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to DigestHolderMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of DigestHolderMock.WriteToFunc invocations
func (m *DigestHolderMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of DigestHolderMock.WriteTo invocations
func (m *DigestHolderMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *DigestHolderMock) WriteToFinished() bool {
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
func (m *DigestHolderMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to DigestHolderMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DigestHolderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DigestHolderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DigestHolderMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to DigestHolderMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to DigestHolderMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DigestHolderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DigestHolderMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to DigestHolderMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to DigestHolderMock.AsBytes")
			}

			if !m.CopyOfDigestFinished() {
				m.t.Error("Expected call to DigestHolderMock.CopyOfDigest")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to DigestHolderMock.Equals")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to DigestHolderMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to DigestHolderMock.FoldToUint64")
			}

			if !m.GetDigestMethodFinished() {
				m.t.Error("Expected call to DigestHolderMock.GetDigestMethod")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to DigestHolderMock.Read")
			}

			if !m.SignWithFinished() {
				m.t.Error("Expected call to DigestHolderMock.SignWith")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to DigestHolderMock.WriteTo")
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
func (m *DigestHolderMock) AllMocksCalled() bool {

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
