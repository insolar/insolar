package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SignatureHolder" can be found in github.com/insolar/insolar/network/consensus/common
*/
import (
	"io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//SignatureHolderMock implements github.com/insolar/insolar/network/consensus/common.SignatureHolder
type SignatureHolderMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mSignatureHolderMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mSignatureHolderMockAsBytes

	CopyOfSignatureFunc       func() (r Signature)
	CopyOfSignatureCounter    uint64
	CopyOfSignaturePreCounter uint64
	CopyOfSignatureMock       mSignatureHolderMockCopyOfSignature

	EqualsFunc       func(p SignatureHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mSignatureHolderMockEquals

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mSignatureHolderMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mSignatureHolderMockFoldToUint64

	GetSignatureMethodFunc       func() (r SignatureMethod)
	GetSignatureMethodCounter    uint64
	GetSignatureMethodPreCounter uint64
	GetSignatureMethodMock       mSignatureHolderMockGetSignatureMethod

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mSignatureHolderMockRead

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mSignatureHolderMockWriteTo
}

//NewSignatureHolderMock returns a mock for github.com/insolar/insolar/network/consensus/common.SignatureHolder
func NewSignatureHolderMock(t minimock.Tester) *SignatureHolderMock {
	m := &SignatureHolderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mSignatureHolderMockAsByteString{mock: m}
	m.AsBytesMock = mSignatureHolderMockAsBytes{mock: m}
	m.CopyOfSignatureMock = mSignatureHolderMockCopyOfSignature{mock: m}
	m.EqualsMock = mSignatureHolderMockEquals{mock: m}
	m.FixedByteSizeMock = mSignatureHolderMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mSignatureHolderMockFoldToUint64{mock: m}
	m.GetSignatureMethodMock = mSignatureHolderMockGetSignatureMethod{mock: m}
	m.ReadMock = mSignatureHolderMockRead{mock: m}
	m.WriteToMock = mSignatureHolderMockWriteTo{mock: m}

	return m
}

type mSignatureHolderMockAsByteString struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockAsByteStringExpectation
	expectationSeries []*SignatureHolderMockAsByteStringExpectation
}

type SignatureHolderMockAsByteStringExpectation struct {
	result *SignatureHolderMockAsByteStringResult
}

type SignatureHolderMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of SignatureHolder.AsByteString is expected from 1 to Infinity times
func (m *mSignatureHolderMockAsByteString) Expect() *mSignatureHolderMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureHolder.AsByteString
func (m *mSignatureHolderMockAsByteString) Return(r string) *SignatureHolderMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.AsByteString is expected once
func (m *mSignatureHolderMockAsByteString) ExpectOnce() *SignatureHolderMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockAsByteStringExpectation) Return(r string) {
	e.result = &SignatureHolderMockAsByteStringResult{r}
}

//Set uses given function f as a mock of SignatureHolder.AsByteString method
func (m *mSignatureHolderMockAsByteString) Set(f func() (r string)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of SignatureHolderMock.AsByteStringFunc invocations
func (m *SignatureHolderMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of SignatureHolderMock.AsByteString invocations
func (m *SignatureHolderMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) AsByteStringFinished() bool {
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

type mSignatureHolderMockAsBytes struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockAsBytesExpectation
	expectationSeries []*SignatureHolderMockAsBytesExpectation
}

type SignatureHolderMockAsBytesExpectation struct {
	result *SignatureHolderMockAsBytesResult
}

type SignatureHolderMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of SignatureHolder.AsBytes is expected from 1 to Infinity times
func (m *mSignatureHolderMockAsBytes) Expect() *mSignatureHolderMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureHolder.AsBytes
func (m *mSignatureHolderMockAsBytes) Return(r []byte) *SignatureHolderMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.AsBytes is expected once
func (m *mSignatureHolderMockAsBytes) ExpectOnce() *SignatureHolderMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockAsBytesExpectation) Return(r []byte) {
	e.result = &SignatureHolderMockAsBytesResult{r}
}

//Set uses given function f as a mock of SignatureHolder.AsBytes method
func (m *mSignatureHolderMockAsBytes) Set(f func() (r []byte)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of SignatureHolderMock.AsBytesFunc invocations
func (m *SignatureHolderMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of SignatureHolderMock.AsBytes invocations
func (m *SignatureHolderMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) AsBytesFinished() bool {
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

type mSignatureHolderMockCopyOfSignature struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockCopyOfSignatureExpectation
	expectationSeries []*SignatureHolderMockCopyOfSignatureExpectation
}

type SignatureHolderMockCopyOfSignatureExpectation struct {
	result *SignatureHolderMockCopyOfSignatureResult
}

type SignatureHolderMockCopyOfSignatureResult struct {
	r Signature
}

//Expect specifies that invocation of SignatureHolder.CopyOfSignature is expected from 1 to Infinity times
func (m *mSignatureHolderMockCopyOfSignature) Expect() *mSignatureHolderMockCopyOfSignature {
	m.mock.CopyOfSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockCopyOfSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureHolder.CopyOfSignature
func (m *mSignatureHolderMockCopyOfSignature) Return(r Signature) *SignatureHolderMock {
	m.mock.CopyOfSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockCopyOfSignatureExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockCopyOfSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.CopyOfSignature is expected once
func (m *mSignatureHolderMockCopyOfSignature) ExpectOnce() *SignatureHolderMockCopyOfSignatureExpectation {
	m.mock.CopyOfSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockCopyOfSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockCopyOfSignatureExpectation) Return(r Signature) {
	e.result = &SignatureHolderMockCopyOfSignatureResult{r}
}

//Set uses given function f as a mock of SignatureHolder.CopyOfSignature method
func (m *mSignatureHolderMockCopyOfSignature) Set(f func() (r Signature)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CopyOfSignatureFunc = f
	return m.mock
}

//CopyOfSignature implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) CopyOfSignature() (r Signature) {
	counter := atomic.AddUint64(&m.CopyOfSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.CopyOfSignatureCounter, 1)

	if len(m.CopyOfSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CopyOfSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.CopyOfSignature.")
			return
		}

		result := m.CopyOfSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.CopyOfSignature")
			return
		}

		r = result.r

		return
	}

	if m.CopyOfSignatureMock.mainExpectation != nil {

		result := m.CopyOfSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.CopyOfSignature")
		}

		r = result.r

		return
	}

	if m.CopyOfSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.CopyOfSignature.")
		return
	}

	return m.CopyOfSignatureFunc()
}

//CopyOfSignatureMinimockCounter returns a count of SignatureHolderMock.CopyOfSignatureFunc invocations
func (m *SignatureHolderMock) CopyOfSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfSignatureCounter)
}

//CopyOfSignatureMinimockPreCounter returns the value of SignatureHolderMock.CopyOfSignature invocations
func (m *SignatureHolderMock) CopyOfSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfSignaturePreCounter)
}

//CopyOfSignatureFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) CopyOfSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CopyOfSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CopyOfSignatureCounter) == uint64(len(m.CopyOfSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CopyOfSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CopyOfSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CopyOfSignatureFunc != nil {
		return atomic.LoadUint64(&m.CopyOfSignatureCounter) > 0
	}

	return true
}

type mSignatureHolderMockEquals struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockEqualsExpectation
	expectationSeries []*SignatureHolderMockEqualsExpectation
}

type SignatureHolderMockEqualsExpectation struct {
	input  *SignatureHolderMockEqualsInput
	result *SignatureHolderMockEqualsResult
}

type SignatureHolderMockEqualsInput struct {
	p SignatureHolder
}

type SignatureHolderMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of SignatureHolder.Equals is expected from 1 to Infinity times
func (m *mSignatureHolderMockEquals) Expect(p SignatureHolder) *mSignatureHolderMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockEqualsExpectation{}
	}
	m.mainExpectation.input = &SignatureHolderMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of SignatureHolder.Equals
func (m *mSignatureHolderMockEquals) Return(r bool) *SignatureHolderMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockEqualsExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.Equals is expected once
func (m *mSignatureHolderMockEquals) ExpectOnce(p SignatureHolder) *SignatureHolderMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockEqualsExpectation{}
	expectation.input = &SignatureHolderMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockEqualsExpectation) Return(r bool) {
	e.result = &SignatureHolderMockEqualsResult{r}
}

//Set uses given function f as a mock of SignatureHolder.Equals method
func (m *mSignatureHolderMockEquals) Set(f func(p SignatureHolder) (r bool)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) Equals(p SignatureHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureHolderMockEqualsInput{p}, "SignatureHolder.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureHolderMockEqualsInput{p}, "SignatureHolder.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of SignatureHolderMock.EqualsFunc invocations
func (m *SignatureHolderMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of SignatureHolderMock.Equals invocations
func (m *SignatureHolderMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) EqualsFinished() bool {
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

type mSignatureHolderMockFixedByteSize struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockFixedByteSizeExpectation
	expectationSeries []*SignatureHolderMockFixedByteSizeExpectation
}

type SignatureHolderMockFixedByteSizeExpectation struct {
	result *SignatureHolderMockFixedByteSizeResult
}

type SignatureHolderMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of SignatureHolder.FixedByteSize is expected from 1 to Infinity times
func (m *mSignatureHolderMockFixedByteSize) Expect() *mSignatureHolderMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureHolder.FixedByteSize
func (m *mSignatureHolderMockFixedByteSize) Return(r int) *SignatureHolderMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.FixedByteSize is expected once
func (m *mSignatureHolderMockFixedByteSize) ExpectOnce() *SignatureHolderMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockFixedByteSizeExpectation) Return(r int) {
	e.result = &SignatureHolderMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of SignatureHolder.FixedByteSize method
func (m *mSignatureHolderMockFixedByteSize) Set(f func() (r int)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of SignatureHolderMock.FixedByteSizeFunc invocations
func (m *SignatureHolderMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of SignatureHolderMock.FixedByteSize invocations
func (m *SignatureHolderMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) FixedByteSizeFinished() bool {
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

type mSignatureHolderMockFoldToUint64 struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockFoldToUint64Expectation
	expectationSeries []*SignatureHolderMockFoldToUint64Expectation
}

type SignatureHolderMockFoldToUint64Expectation struct {
	result *SignatureHolderMockFoldToUint64Result
}

type SignatureHolderMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of SignatureHolder.FoldToUint64 is expected from 1 to Infinity times
func (m *mSignatureHolderMockFoldToUint64) Expect() *mSignatureHolderMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureHolder.FoldToUint64
func (m *mSignatureHolderMockFoldToUint64) Return(r uint64) *SignatureHolderMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.FoldToUint64 is expected once
func (m *mSignatureHolderMockFoldToUint64) ExpectOnce() *SignatureHolderMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &SignatureHolderMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of SignatureHolder.FoldToUint64 method
func (m *mSignatureHolderMockFoldToUint64) Set(f func() (r uint64)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of SignatureHolderMock.FoldToUint64Func invocations
func (m *SignatureHolderMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of SignatureHolderMock.FoldToUint64 invocations
func (m *SignatureHolderMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *SignatureHolderMock) FoldToUint64Finished() bool {
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

type mSignatureHolderMockGetSignatureMethod struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockGetSignatureMethodExpectation
	expectationSeries []*SignatureHolderMockGetSignatureMethodExpectation
}

type SignatureHolderMockGetSignatureMethodExpectation struct {
	result *SignatureHolderMockGetSignatureMethodResult
}

type SignatureHolderMockGetSignatureMethodResult struct {
	r SignatureMethod
}

//Expect specifies that invocation of SignatureHolder.GetSignatureMethod is expected from 1 to Infinity times
func (m *mSignatureHolderMockGetSignatureMethod) Expect() *mSignatureHolderMockGetSignatureMethod {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockGetSignatureMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureHolder.GetSignatureMethod
func (m *mSignatureHolderMockGetSignatureMethod) Return(r SignatureMethod) *SignatureHolderMock {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockGetSignatureMethodExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockGetSignatureMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.GetSignatureMethod is expected once
func (m *mSignatureHolderMockGetSignatureMethod) ExpectOnce() *SignatureHolderMockGetSignatureMethodExpectation {
	m.mock.GetSignatureMethodFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockGetSignatureMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockGetSignatureMethodExpectation) Return(r SignatureMethod) {
	e.result = &SignatureHolderMockGetSignatureMethodResult{r}
}

//Set uses given function f as a mock of SignatureHolder.GetSignatureMethod method
func (m *mSignatureHolderMockGetSignatureMethod) Set(f func() (r SignatureMethod)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureMethodFunc = f
	return m.mock
}

//GetSignatureMethod implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) GetSignatureMethod() (r SignatureMethod) {
	counter := atomic.AddUint64(&m.GetSignatureMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureMethodCounter, 1)

	if len(m.GetSignatureMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.GetSignatureMethod.")
			return
		}

		result := m.GetSignatureMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.GetSignatureMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodMock.mainExpectation != nil {

		result := m.GetSignatureMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.GetSignatureMethod")
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.GetSignatureMethod.")
		return
	}

	return m.GetSignatureMethodFunc()
}

//GetSignatureMethodMinimockCounter returns a count of SignatureHolderMock.GetSignatureMethodFunc invocations
func (m *SignatureHolderMock) GetSignatureMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodCounter)
}

//GetSignatureMethodMinimockPreCounter returns the value of SignatureHolderMock.GetSignatureMethod invocations
func (m *SignatureHolderMock) GetSignatureMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodPreCounter)
}

//GetSignatureMethodFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) GetSignatureMethodFinished() bool {
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

type mSignatureHolderMockRead struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockReadExpectation
	expectationSeries []*SignatureHolderMockReadExpectation
}

type SignatureHolderMockReadExpectation struct {
	input  *SignatureHolderMockReadInput
	result *SignatureHolderMockReadResult
}

type SignatureHolderMockReadInput struct {
	p []byte
}

type SignatureHolderMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of SignatureHolder.Read is expected from 1 to Infinity times
func (m *mSignatureHolderMockRead) Expect(p []byte) *mSignatureHolderMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockReadExpectation{}
	}
	m.mainExpectation.input = &SignatureHolderMockReadInput{p}
	return m
}

//Return specifies results of invocation of SignatureHolder.Read
func (m *mSignatureHolderMockRead) Return(r int, r1 error) *SignatureHolderMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockReadExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.Read is expected once
func (m *mSignatureHolderMockRead) ExpectOnce(p []byte) *SignatureHolderMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockReadExpectation{}
	expectation.input = &SignatureHolderMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockReadExpectation) Return(r int, r1 error) {
	e.result = &SignatureHolderMockReadResult{r, r1}
}

//Set uses given function f as a mock of SignatureHolder.Read method
func (m *mSignatureHolderMockRead) Set(f func(p []byte) (r int, r1 error)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureHolderMockReadInput{p}, "SignatureHolder.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureHolderMockReadInput{p}, "SignatureHolder.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of SignatureHolderMock.ReadFunc invocations
func (m *SignatureHolderMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of SignatureHolderMock.Read invocations
func (m *SignatureHolderMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) ReadFinished() bool {
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

type mSignatureHolderMockWriteTo struct {
	mock              *SignatureHolderMock
	mainExpectation   *SignatureHolderMockWriteToExpectation
	expectationSeries []*SignatureHolderMockWriteToExpectation
}

type SignatureHolderMockWriteToExpectation struct {
	input  *SignatureHolderMockWriteToInput
	result *SignatureHolderMockWriteToResult
}

type SignatureHolderMockWriteToInput struct {
	p io.Writer
}

type SignatureHolderMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of SignatureHolder.WriteTo is expected from 1 to Infinity times
func (m *mSignatureHolderMockWriteTo) Expect(p io.Writer) *mSignatureHolderMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockWriteToExpectation{}
	}
	m.mainExpectation.input = &SignatureHolderMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of SignatureHolder.WriteTo
func (m *mSignatureHolderMockWriteTo) Return(r int64, r1 error) *SignatureHolderMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureHolderMockWriteToExpectation{}
	}
	m.mainExpectation.result = &SignatureHolderMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureHolder.WriteTo is expected once
func (m *mSignatureHolderMockWriteTo) ExpectOnce(p io.Writer) *SignatureHolderMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureHolderMockWriteToExpectation{}
	expectation.input = &SignatureHolderMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureHolderMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &SignatureHolderMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of SignatureHolder.WriteTo method
func (m *mSignatureHolderMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *SignatureHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/common.SignatureHolder interface
func (m *SignatureHolderMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureHolderMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureHolderMockWriteToInput{p}, "SignatureHolder.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureHolderMockWriteToInput{p}, "SignatureHolder.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureHolderMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureHolderMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of SignatureHolderMock.WriteToFunc invocations
func (m *SignatureHolderMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of SignatureHolderMock.WriteTo invocations
func (m *SignatureHolderMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *SignatureHolderMock) WriteToFinished() bool {
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
func (m *SignatureHolderMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.AsBytes")
	}

	if !m.CopyOfSignatureFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.CopyOfSignature")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to SignatureHolderMock.FoldToUint64")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.GetSignatureMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SignatureHolderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SignatureHolderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SignatureHolderMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.AsBytes")
	}

	if !m.CopyOfSignatureFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.CopyOfSignature")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to SignatureHolderMock.FoldToUint64")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.GetSignatureMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to SignatureHolderMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SignatureHolderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SignatureHolderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AsByteStringFinished()
		ok = ok && m.AsBytesFinished()
		ok = ok && m.CopyOfSignatureFinished()
		ok = ok && m.EqualsFinished()
		ok = ok && m.FixedByteSizeFinished()
		ok = ok && m.FoldToUint64Finished()
		ok = ok && m.GetSignatureMethodFinished()
		ok = ok && m.ReadFinished()
		ok = ok && m.WriteToFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AsByteStringFinished() {
				m.t.Error("Expected call to SignatureHolderMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to SignatureHolderMock.AsBytes")
			}

			if !m.CopyOfSignatureFinished() {
				m.t.Error("Expected call to SignatureHolderMock.CopyOfSignature")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to SignatureHolderMock.Equals")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to SignatureHolderMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to SignatureHolderMock.FoldToUint64")
			}

			if !m.GetSignatureMethodFinished() {
				m.t.Error("Expected call to SignatureHolderMock.GetSignatureMethod")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to SignatureHolderMock.Read")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to SignatureHolderMock.WriteTo")
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
func (m *SignatureHolderMock) AllMocksCalled() bool {

	if !m.AsByteStringFinished() {
		return false
	}

	if !m.AsBytesFinished() {
		return false
	}

	if !m.CopyOfSignatureFinished() {
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

	if !m.GetSignatureMethodFinished() {
		return false
	}

	if !m.ReadFinished() {
		return false
	}

	if !m.WriteToFinished() {
		return false
	}

	return true
}
