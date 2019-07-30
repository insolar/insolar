package proofs

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CloudStateHash" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/proofs
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"

	testify_assert "github.com/stretchr/testify/assert"
)

//CloudStateHashMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash
type CloudStateHashMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mCloudStateHashMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mCloudStateHashMockAsBytes

	CopyOfDigestFunc       func() (r cryptkit.Digest)
	CopyOfDigestCounter    uint64
	CopyOfDigestPreCounter uint64
	CopyOfDigestMock       mCloudStateHashMockCopyOfDigest

	EqualsFunc       func(p cryptkit.DigestHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mCloudStateHashMockEquals

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mCloudStateHashMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mCloudStateHashMockFoldToUint64

	GetDigestMethodFunc       func() (r cryptkit.DigestMethod)
	GetDigestMethodCounter    uint64
	GetDigestMethodPreCounter uint64
	GetDigestMethodMock       mCloudStateHashMockGetDigestMethod

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mCloudStateHashMockRead

	SignWithFunc       func(p cryptkit.DigestSigner) (r cryptkit.SignedDigestHolder)
	SignWithCounter    uint64
	SignWithPreCounter uint64
	SignWithMock       mCloudStateHashMockSignWith

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mCloudStateHashMockWriteTo
}

//NewCloudStateHashMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash
func NewCloudStateHashMock(t minimock.Tester) *CloudStateHashMock {
	m := &CloudStateHashMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mCloudStateHashMockAsByteString{mock: m}
	m.AsBytesMock = mCloudStateHashMockAsBytes{mock: m}
	m.CopyOfDigestMock = mCloudStateHashMockCopyOfDigest{mock: m}
	m.EqualsMock = mCloudStateHashMockEquals{mock: m}
	m.FixedByteSizeMock = mCloudStateHashMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mCloudStateHashMockFoldToUint64{mock: m}
	m.GetDigestMethodMock = mCloudStateHashMockGetDigestMethod{mock: m}
	m.ReadMock = mCloudStateHashMockRead{mock: m}
	m.SignWithMock = mCloudStateHashMockSignWith{mock: m}
	m.WriteToMock = mCloudStateHashMockWriteTo{mock: m}

	return m
}

type mCloudStateHashMockAsByteString struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockAsByteStringExpectation
	expectationSeries []*CloudStateHashMockAsByteStringExpectation
}

type CloudStateHashMockAsByteStringExpectation struct {
	result *CloudStateHashMockAsByteStringResult
}

type CloudStateHashMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of CloudStateHash.AsByteString is expected from 1 to Infinity times
func (m *mCloudStateHashMockAsByteString) Expect() *mCloudStateHashMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of CloudStateHash.AsByteString
func (m *mCloudStateHashMockAsByteString) Return(r string) *CloudStateHashMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.AsByteString is expected once
func (m *mCloudStateHashMockAsByteString) ExpectOnce() *CloudStateHashMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockAsByteStringExpectation) Return(r string) {
	e.result = &CloudStateHashMockAsByteStringResult{r}
}

//Set uses given function f as a mock of CloudStateHash.AsByteString method
func (m *mCloudStateHashMockAsByteString) Set(f func() (r string)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of CloudStateHashMock.AsByteStringFunc invocations
func (m *CloudStateHashMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of CloudStateHashMock.AsByteString invocations
func (m *CloudStateHashMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) AsByteStringFinished() bool {
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

type mCloudStateHashMockAsBytes struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockAsBytesExpectation
	expectationSeries []*CloudStateHashMockAsBytesExpectation
}

type CloudStateHashMockAsBytesExpectation struct {
	result *CloudStateHashMockAsBytesResult
}

type CloudStateHashMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of CloudStateHash.AsBytes is expected from 1 to Infinity times
func (m *mCloudStateHashMockAsBytes) Expect() *mCloudStateHashMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of CloudStateHash.AsBytes
func (m *mCloudStateHashMockAsBytes) Return(r []byte) *CloudStateHashMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.AsBytes is expected once
func (m *mCloudStateHashMockAsBytes) ExpectOnce() *CloudStateHashMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockAsBytesExpectation) Return(r []byte) {
	e.result = &CloudStateHashMockAsBytesResult{r}
}

//Set uses given function f as a mock of CloudStateHash.AsBytes method
func (m *mCloudStateHashMockAsBytes) Set(f func() (r []byte)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of CloudStateHashMock.AsBytesFunc invocations
func (m *CloudStateHashMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of CloudStateHashMock.AsBytes invocations
func (m *CloudStateHashMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) AsBytesFinished() bool {
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

type mCloudStateHashMockCopyOfDigest struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockCopyOfDigestExpectation
	expectationSeries []*CloudStateHashMockCopyOfDigestExpectation
}

type CloudStateHashMockCopyOfDigestExpectation struct {
	result *CloudStateHashMockCopyOfDigestResult
}

type CloudStateHashMockCopyOfDigestResult struct {
	r cryptkit.Digest
}

//Expect specifies that invocation of CloudStateHash.CopyOfDigest is expected from 1 to Infinity times
func (m *mCloudStateHashMockCopyOfDigest) Expect() *mCloudStateHashMockCopyOfDigest {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockCopyOfDigestExpectation{}
	}

	return m
}

//Return specifies results of invocation of CloudStateHash.CopyOfDigest
func (m *mCloudStateHashMockCopyOfDigest) Return(r cryptkit.Digest) *CloudStateHashMock {
	m.mock.CopyOfDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockCopyOfDigestExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockCopyOfDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.CopyOfDigest is expected once
func (m *mCloudStateHashMockCopyOfDigest) ExpectOnce() *CloudStateHashMockCopyOfDigestExpectation {
	m.mock.CopyOfDigestFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockCopyOfDigestExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockCopyOfDigestExpectation) Return(r cryptkit.Digest) {
	e.result = &CloudStateHashMockCopyOfDigestResult{r}
}

//Set uses given function f as a mock of CloudStateHash.CopyOfDigest method
func (m *mCloudStateHashMockCopyOfDigest) Set(f func() (r cryptkit.Digest)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CopyOfDigestFunc = f
	return m.mock
}

//CopyOfDigest implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) CopyOfDigest() (r cryptkit.Digest) {
	counter := atomic.AddUint64(&m.CopyOfDigestPreCounter, 1)
	defer atomic.AddUint64(&m.CopyOfDigestCounter, 1)

	if len(m.CopyOfDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CopyOfDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.CopyOfDigest.")
			return
		}

		result := m.CopyOfDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.CopyOfDigest")
			return
		}

		r = result.r

		return
	}

	if m.CopyOfDigestMock.mainExpectation != nil {

		result := m.CopyOfDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.CopyOfDigest")
		}

		r = result.r

		return
	}

	if m.CopyOfDigestFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.CopyOfDigest.")
		return
	}

	return m.CopyOfDigestFunc()
}

//CopyOfDigestMinimockCounter returns a count of CloudStateHashMock.CopyOfDigestFunc invocations
func (m *CloudStateHashMock) CopyOfDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestCounter)
}

//CopyOfDigestMinimockPreCounter returns the value of CloudStateHashMock.CopyOfDigest invocations
func (m *CloudStateHashMock) CopyOfDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfDigestPreCounter)
}

//CopyOfDigestFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) CopyOfDigestFinished() bool {
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

type mCloudStateHashMockEquals struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockEqualsExpectation
	expectationSeries []*CloudStateHashMockEqualsExpectation
}

type CloudStateHashMockEqualsExpectation struct {
	input  *CloudStateHashMockEqualsInput
	result *CloudStateHashMockEqualsResult
}

type CloudStateHashMockEqualsInput struct {
	p cryptkit.DigestHolder
}

type CloudStateHashMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of CloudStateHash.Equals is expected from 1 to Infinity times
func (m *mCloudStateHashMockEquals) Expect(p cryptkit.DigestHolder) *mCloudStateHashMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockEqualsExpectation{}
	}
	m.mainExpectation.input = &CloudStateHashMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of CloudStateHash.Equals
func (m *mCloudStateHashMockEquals) Return(r bool) *CloudStateHashMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockEqualsExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.Equals is expected once
func (m *mCloudStateHashMockEquals) ExpectOnce(p cryptkit.DigestHolder) *CloudStateHashMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockEqualsExpectation{}
	expectation.input = &CloudStateHashMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockEqualsExpectation) Return(r bool) {
	e.result = &CloudStateHashMockEqualsResult{r}
}

//Set uses given function f as a mock of CloudStateHash.Equals method
func (m *mCloudStateHashMockEquals) Set(f func(p cryptkit.DigestHolder) (r bool)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) Equals(p cryptkit.DigestHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CloudStateHashMockEqualsInput{p}, "CloudStateHash.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CloudStateHashMockEqualsInput{p}, "CloudStateHash.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of CloudStateHashMock.EqualsFunc invocations
func (m *CloudStateHashMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of CloudStateHashMock.Equals invocations
func (m *CloudStateHashMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) EqualsFinished() bool {
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

type mCloudStateHashMockFixedByteSize struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockFixedByteSizeExpectation
	expectationSeries []*CloudStateHashMockFixedByteSizeExpectation
}

type CloudStateHashMockFixedByteSizeExpectation struct {
	result *CloudStateHashMockFixedByteSizeResult
}

type CloudStateHashMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of CloudStateHash.FixedByteSize is expected from 1 to Infinity times
func (m *mCloudStateHashMockFixedByteSize) Expect() *mCloudStateHashMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of CloudStateHash.FixedByteSize
func (m *mCloudStateHashMockFixedByteSize) Return(r int) *CloudStateHashMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.FixedByteSize is expected once
func (m *mCloudStateHashMockFixedByteSize) ExpectOnce() *CloudStateHashMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockFixedByteSizeExpectation) Return(r int) {
	e.result = &CloudStateHashMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of CloudStateHash.FixedByteSize method
func (m *mCloudStateHashMockFixedByteSize) Set(f func() (r int)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of CloudStateHashMock.FixedByteSizeFunc invocations
func (m *CloudStateHashMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of CloudStateHashMock.FixedByteSize invocations
func (m *CloudStateHashMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) FixedByteSizeFinished() bool {
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

type mCloudStateHashMockFoldToUint64 struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockFoldToUint64Expectation
	expectationSeries []*CloudStateHashMockFoldToUint64Expectation
}

type CloudStateHashMockFoldToUint64Expectation struct {
	result *CloudStateHashMockFoldToUint64Result
}

type CloudStateHashMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of CloudStateHash.FoldToUint64 is expected from 1 to Infinity times
func (m *mCloudStateHashMockFoldToUint64) Expect() *mCloudStateHashMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of CloudStateHash.FoldToUint64
func (m *mCloudStateHashMockFoldToUint64) Return(r uint64) *CloudStateHashMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.FoldToUint64 is expected once
func (m *mCloudStateHashMockFoldToUint64) ExpectOnce() *CloudStateHashMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &CloudStateHashMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of CloudStateHash.FoldToUint64 method
func (m *mCloudStateHashMockFoldToUint64) Set(f func() (r uint64)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of CloudStateHashMock.FoldToUint64Func invocations
func (m *CloudStateHashMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of CloudStateHashMock.FoldToUint64 invocations
func (m *CloudStateHashMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *CloudStateHashMock) FoldToUint64Finished() bool {
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

type mCloudStateHashMockGetDigestMethod struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockGetDigestMethodExpectation
	expectationSeries []*CloudStateHashMockGetDigestMethodExpectation
}

type CloudStateHashMockGetDigestMethodExpectation struct {
	result *CloudStateHashMockGetDigestMethodResult
}

type CloudStateHashMockGetDigestMethodResult struct {
	r cryptkit.DigestMethod
}

//Expect specifies that invocation of CloudStateHash.GetDigestMethod is expected from 1 to Infinity times
func (m *mCloudStateHashMockGetDigestMethod) Expect() *mCloudStateHashMockGetDigestMethod {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockGetDigestMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of CloudStateHash.GetDigestMethod
func (m *mCloudStateHashMockGetDigestMethod) Return(r cryptkit.DigestMethod) *CloudStateHashMock {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockGetDigestMethodExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockGetDigestMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.GetDigestMethod is expected once
func (m *mCloudStateHashMockGetDigestMethod) ExpectOnce() *CloudStateHashMockGetDigestMethodExpectation {
	m.mock.GetDigestMethodFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockGetDigestMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockGetDigestMethodExpectation) Return(r cryptkit.DigestMethod) {
	e.result = &CloudStateHashMockGetDigestMethodResult{r}
}

//Set uses given function f as a mock of CloudStateHash.GetDigestMethod method
func (m *mCloudStateHashMockGetDigestMethod) Set(f func() (r cryptkit.DigestMethod)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDigestMethodFunc = f
	return m.mock
}

//GetDigestMethod implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) GetDigestMethod() (r cryptkit.DigestMethod) {
	counter := atomic.AddUint64(&m.GetDigestMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetDigestMethodCounter, 1)

	if len(m.GetDigestMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDigestMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.GetDigestMethod.")
			return
		}

		result := m.GetDigestMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.GetDigestMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetDigestMethodMock.mainExpectation != nil {

		result := m.GetDigestMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.GetDigestMethod")
		}

		r = result.r

		return
	}

	if m.GetDigestMethodFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.GetDigestMethod.")
		return
	}

	return m.GetDigestMethodFunc()
}

//GetDigestMethodMinimockCounter returns a count of CloudStateHashMock.GetDigestMethodFunc invocations
func (m *CloudStateHashMock) GetDigestMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodCounter)
}

//GetDigestMethodMinimockPreCounter returns the value of CloudStateHashMock.GetDigestMethod invocations
func (m *CloudStateHashMock) GetDigestMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodPreCounter)
}

//GetDigestMethodFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) GetDigestMethodFinished() bool {
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

type mCloudStateHashMockRead struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockReadExpectation
	expectationSeries []*CloudStateHashMockReadExpectation
}

type CloudStateHashMockReadExpectation struct {
	input  *CloudStateHashMockReadInput
	result *CloudStateHashMockReadResult
}

type CloudStateHashMockReadInput struct {
	p []byte
}

type CloudStateHashMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of CloudStateHash.Read is expected from 1 to Infinity times
func (m *mCloudStateHashMockRead) Expect(p []byte) *mCloudStateHashMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockReadExpectation{}
	}
	m.mainExpectation.input = &CloudStateHashMockReadInput{p}
	return m
}

//Return specifies results of invocation of CloudStateHash.Read
func (m *mCloudStateHashMockRead) Return(r int, r1 error) *CloudStateHashMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockReadExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.Read is expected once
func (m *mCloudStateHashMockRead) ExpectOnce(p []byte) *CloudStateHashMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockReadExpectation{}
	expectation.input = &CloudStateHashMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockReadExpectation) Return(r int, r1 error) {
	e.result = &CloudStateHashMockReadResult{r, r1}
}

//Set uses given function f as a mock of CloudStateHash.Read method
func (m *mCloudStateHashMockRead) Set(f func(p []byte) (r int, r1 error)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CloudStateHashMockReadInput{p}, "CloudStateHash.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CloudStateHashMockReadInput{p}, "CloudStateHash.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of CloudStateHashMock.ReadFunc invocations
func (m *CloudStateHashMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of CloudStateHashMock.Read invocations
func (m *CloudStateHashMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) ReadFinished() bool {
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

type mCloudStateHashMockSignWith struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockSignWithExpectation
	expectationSeries []*CloudStateHashMockSignWithExpectation
}

type CloudStateHashMockSignWithExpectation struct {
	input  *CloudStateHashMockSignWithInput
	result *CloudStateHashMockSignWithResult
}

type CloudStateHashMockSignWithInput struct {
	p cryptkit.DigestSigner
}

type CloudStateHashMockSignWithResult struct {
	r cryptkit.SignedDigestHolder
}

//Expect specifies that invocation of CloudStateHash.SignWith is expected from 1 to Infinity times
func (m *mCloudStateHashMockSignWith) Expect(p cryptkit.DigestSigner) *mCloudStateHashMockSignWith {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockSignWithExpectation{}
	}
	m.mainExpectation.input = &CloudStateHashMockSignWithInput{p}
	return m
}

//Return specifies results of invocation of CloudStateHash.SignWith
func (m *mCloudStateHashMockSignWith) Return(r cryptkit.SignedDigestHolder) *CloudStateHashMock {
	m.mock.SignWithFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockSignWithExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockSignWithResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.SignWith is expected once
func (m *mCloudStateHashMockSignWith) ExpectOnce(p cryptkit.DigestSigner) *CloudStateHashMockSignWithExpectation {
	m.mock.SignWithFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockSignWithExpectation{}
	expectation.input = &CloudStateHashMockSignWithInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockSignWithExpectation) Return(r cryptkit.SignedDigestHolder) {
	e.result = &CloudStateHashMockSignWithResult{r}
}

//Set uses given function f as a mock of CloudStateHash.SignWith method
func (m *mCloudStateHashMockSignWith) Set(f func(p cryptkit.DigestSigner) (r cryptkit.SignedDigestHolder)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignWithFunc = f
	return m.mock
}

//SignWith implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) SignWith(p cryptkit.DigestSigner) (r cryptkit.SignedDigestHolder) {
	counter := atomic.AddUint64(&m.SignWithPreCounter, 1)
	defer atomic.AddUint64(&m.SignWithCounter, 1)

	if len(m.SignWithMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignWithMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.SignWith. %v", p)
			return
		}

		input := m.SignWithMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CloudStateHashMockSignWithInput{p}, "CloudStateHash.SignWith got unexpected parameters")

		result := m.SignWithMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.SignWith")
			return
		}

		r = result.r

		return
	}

	if m.SignWithMock.mainExpectation != nil {

		input := m.SignWithMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CloudStateHashMockSignWithInput{p}, "CloudStateHash.SignWith got unexpected parameters")
		}

		result := m.SignWithMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.SignWith")
		}

		r = result.r

		return
	}

	if m.SignWithFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.SignWith. %v", p)
		return
	}

	return m.SignWithFunc(p)
}

//SignWithMinimockCounter returns a count of CloudStateHashMock.SignWithFunc invocations
func (m *CloudStateHashMock) SignWithMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithCounter)
}

//SignWithMinimockPreCounter returns the value of CloudStateHashMock.SignWith invocations
func (m *CloudStateHashMock) SignWithMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignWithPreCounter)
}

//SignWithFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) SignWithFinished() bool {
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

type mCloudStateHashMockWriteTo struct {
	mock              *CloudStateHashMock
	mainExpectation   *CloudStateHashMockWriteToExpectation
	expectationSeries []*CloudStateHashMockWriteToExpectation
}

type CloudStateHashMockWriteToExpectation struct {
	input  *CloudStateHashMockWriteToInput
	result *CloudStateHashMockWriteToResult
}

type CloudStateHashMockWriteToInput struct {
	p io.Writer
}

type CloudStateHashMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of CloudStateHash.WriteTo is expected from 1 to Infinity times
func (m *mCloudStateHashMockWriteTo) Expect(p io.Writer) *mCloudStateHashMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockWriteToExpectation{}
	}
	m.mainExpectation.input = &CloudStateHashMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of CloudStateHash.WriteTo
func (m *mCloudStateHashMockWriteTo) Return(r int64, r1 error) *CloudStateHashMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudStateHashMockWriteToExpectation{}
	}
	m.mainExpectation.result = &CloudStateHashMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudStateHash.WriteTo is expected once
func (m *mCloudStateHashMockWriteTo) ExpectOnce(p io.Writer) *CloudStateHashMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &CloudStateHashMockWriteToExpectation{}
	expectation.input = &CloudStateHashMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudStateHashMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &CloudStateHashMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of CloudStateHash.WriteTo method
func (m *mCloudStateHashMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *CloudStateHashMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/gcpv2/api/proofs.CloudStateHash interface
func (m *CloudStateHashMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudStateHashMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CloudStateHashMockWriteToInput{p}, "CloudStateHash.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CloudStateHashMockWriteToInput{p}, "CloudStateHash.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudStateHashMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to CloudStateHashMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of CloudStateHashMock.WriteToFunc invocations
func (m *CloudStateHashMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of CloudStateHashMock.WriteTo invocations
func (m *CloudStateHashMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *CloudStateHashMock) WriteToFinished() bool {
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
func (m *CloudStateHashMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to CloudStateHashMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CloudStateHashMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CloudStateHashMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CloudStateHashMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.AsBytes")
	}

	if !m.CopyOfDigestFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.CopyOfDigest")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to CloudStateHashMock.FoldToUint64")
	}

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.GetDigestMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.Read")
	}

	if !m.SignWithFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.SignWith")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to CloudStateHashMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CloudStateHashMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CloudStateHashMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to CloudStateHashMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to CloudStateHashMock.AsBytes")
			}

			if !m.CopyOfDigestFinished() {
				m.t.Error("Expected call to CloudStateHashMock.CopyOfDigest")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to CloudStateHashMock.Equals")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to CloudStateHashMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to CloudStateHashMock.FoldToUint64")
			}

			if !m.GetDigestMethodFinished() {
				m.t.Error("Expected call to CloudStateHashMock.GetDigestMethod")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to CloudStateHashMock.Read")
			}

			if !m.SignWithFinished() {
				m.t.Error("Expected call to CloudStateHashMock.SignWith")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to CloudStateHashMock.WriteTo")
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
func (m *CloudStateHashMock) AllMocksCalled() bool {

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
