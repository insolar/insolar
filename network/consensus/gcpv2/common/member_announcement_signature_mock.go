package common

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "MemberAnnouncementSignature" can be found in github.com/insolar/insolar/network/consensus/gcpv2/common
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	common "github.com/insolar/insolar/network/consensus/common"

	testify_assert "github.com/stretchr/testify/assert"
)

//MemberAnnouncementSignatureMock implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature
type MemberAnnouncementSignatureMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mMemberAnnouncementSignatureMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mMemberAnnouncementSignatureMockAsBytes

	CopyOfSignatureFunc       func() (r common.Signature)
	CopyOfSignatureCounter    uint64
	CopyOfSignaturePreCounter uint64
	CopyOfSignatureMock       mMemberAnnouncementSignatureMockCopyOfSignature

	EqualsFunc       func(p common.SignatureHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mMemberAnnouncementSignatureMockEquals

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mMemberAnnouncementSignatureMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mMemberAnnouncementSignatureMockFoldToUint64

	GetSignatureMethodFunc       func() (r common.SignatureMethod)
	GetSignatureMethodCounter    uint64
	GetSignatureMethodPreCounter uint64
	GetSignatureMethodMock       mMemberAnnouncementSignatureMockGetSignatureMethod

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mMemberAnnouncementSignatureMockRead

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mMemberAnnouncementSignatureMockWriteTo
}

//NewMemberAnnouncementSignatureMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature
func NewMemberAnnouncementSignatureMock(t minimock.Tester) *MemberAnnouncementSignatureMock {
	m := &MemberAnnouncementSignatureMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mMemberAnnouncementSignatureMockAsByteString{mock: m}
	m.AsBytesMock = mMemberAnnouncementSignatureMockAsBytes{mock: m}
	m.CopyOfSignatureMock = mMemberAnnouncementSignatureMockCopyOfSignature{mock: m}
	m.EqualsMock = mMemberAnnouncementSignatureMockEquals{mock: m}
	m.FixedByteSizeMock = mMemberAnnouncementSignatureMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mMemberAnnouncementSignatureMockFoldToUint64{mock: m}
	m.GetSignatureMethodMock = mMemberAnnouncementSignatureMockGetSignatureMethod{mock: m}
	m.ReadMock = mMemberAnnouncementSignatureMockRead{mock: m}
	m.WriteToMock = mMemberAnnouncementSignatureMockWriteTo{mock: m}

	return m
}

type mMemberAnnouncementSignatureMockAsByteString struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockAsByteStringExpectation
	expectationSeries []*MemberAnnouncementSignatureMockAsByteStringExpectation
}

type MemberAnnouncementSignatureMockAsByteStringExpectation struct {
	result *MemberAnnouncementSignatureMockAsByteStringResult
}

type MemberAnnouncementSignatureMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of MemberAnnouncementSignature.AsByteString is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockAsByteString) Expect() *mMemberAnnouncementSignatureMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.AsByteString
func (m *mMemberAnnouncementSignatureMockAsByteString) Return(r string) *MemberAnnouncementSignatureMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.AsByteString is expected once
func (m *mMemberAnnouncementSignatureMockAsByteString) ExpectOnce() *MemberAnnouncementSignatureMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockAsByteStringExpectation) Return(r string) {
	e.result = &MemberAnnouncementSignatureMockAsByteStringResult{r}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.AsByteString method
func (m *mMemberAnnouncementSignatureMockAsByteString) Set(f func() (r string)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of MemberAnnouncementSignatureMock.AsByteStringFunc invocations
func (m *MemberAnnouncementSignatureMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.AsByteString invocations
func (m *MemberAnnouncementSignatureMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) AsByteStringFinished() bool {
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

type mMemberAnnouncementSignatureMockAsBytes struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockAsBytesExpectation
	expectationSeries []*MemberAnnouncementSignatureMockAsBytesExpectation
}

type MemberAnnouncementSignatureMockAsBytesExpectation struct {
	result *MemberAnnouncementSignatureMockAsBytesResult
}

type MemberAnnouncementSignatureMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of MemberAnnouncementSignature.AsBytes is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockAsBytes) Expect() *mMemberAnnouncementSignatureMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.AsBytes
func (m *mMemberAnnouncementSignatureMockAsBytes) Return(r []byte) *MemberAnnouncementSignatureMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.AsBytes is expected once
func (m *mMemberAnnouncementSignatureMockAsBytes) ExpectOnce() *MemberAnnouncementSignatureMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockAsBytesExpectation) Return(r []byte) {
	e.result = &MemberAnnouncementSignatureMockAsBytesResult{r}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.AsBytes method
func (m *mMemberAnnouncementSignatureMockAsBytes) Set(f func() (r []byte)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of MemberAnnouncementSignatureMock.AsBytesFunc invocations
func (m *MemberAnnouncementSignatureMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.AsBytes invocations
func (m *MemberAnnouncementSignatureMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) AsBytesFinished() bool {
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

type mMemberAnnouncementSignatureMockCopyOfSignature struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockCopyOfSignatureExpectation
	expectationSeries []*MemberAnnouncementSignatureMockCopyOfSignatureExpectation
}

type MemberAnnouncementSignatureMockCopyOfSignatureExpectation struct {
	result *MemberAnnouncementSignatureMockCopyOfSignatureResult
}

type MemberAnnouncementSignatureMockCopyOfSignatureResult struct {
	r common.Signature
}

//Expect specifies that invocation of MemberAnnouncementSignature.CopyOfSignature is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockCopyOfSignature) Expect() *mMemberAnnouncementSignatureMockCopyOfSignature {
	m.mock.CopyOfSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockCopyOfSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.CopyOfSignature
func (m *mMemberAnnouncementSignatureMockCopyOfSignature) Return(r common.Signature) *MemberAnnouncementSignatureMock {
	m.mock.CopyOfSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockCopyOfSignatureExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockCopyOfSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.CopyOfSignature is expected once
func (m *mMemberAnnouncementSignatureMockCopyOfSignature) ExpectOnce() *MemberAnnouncementSignatureMockCopyOfSignatureExpectation {
	m.mock.CopyOfSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockCopyOfSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockCopyOfSignatureExpectation) Return(r common.Signature) {
	e.result = &MemberAnnouncementSignatureMockCopyOfSignatureResult{r}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.CopyOfSignature method
func (m *mMemberAnnouncementSignatureMockCopyOfSignature) Set(f func() (r common.Signature)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CopyOfSignatureFunc = f
	return m.mock
}

//CopyOfSignature implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) CopyOfSignature() (r common.Signature) {
	counter := atomic.AddUint64(&m.CopyOfSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.CopyOfSignatureCounter, 1)

	if len(m.CopyOfSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CopyOfSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.CopyOfSignature.")
			return
		}

		result := m.CopyOfSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.CopyOfSignature")
			return
		}

		r = result.r

		return
	}

	if m.CopyOfSignatureMock.mainExpectation != nil {

		result := m.CopyOfSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.CopyOfSignature")
		}

		r = result.r

		return
	}

	if m.CopyOfSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.CopyOfSignature.")
		return
	}

	return m.CopyOfSignatureFunc()
}

//CopyOfSignatureMinimockCounter returns a count of MemberAnnouncementSignatureMock.CopyOfSignatureFunc invocations
func (m *MemberAnnouncementSignatureMock) CopyOfSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfSignatureCounter)
}

//CopyOfSignatureMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.CopyOfSignature invocations
func (m *MemberAnnouncementSignatureMock) CopyOfSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CopyOfSignaturePreCounter)
}

//CopyOfSignatureFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) CopyOfSignatureFinished() bool {
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

type mMemberAnnouncementSignatureMockEquals struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockEqualsExpectation
	expectationSeries []*MemberAnnouncementSignatureMockEqualsExpectation
}

type MemberAnnouncementSignatureMockEqualsExpectation struct {
	input  *MemberAnnouncementSignatureMockEqualsInput
	result *MemberAnnouncementSignatureMockEqualsResult
}

type MemberAnnouncementSignatureMockEqualsInput struct {
	p common.SignatureHolder
}

type MemberAnnouncementSignatureMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of MemberAnnouncementSignature.Equals is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockEquals) Expect(p common.SignatureHolder) *mMemberAnnouncementSignatureMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockEqualsExpectation{}
	}
	m.mainExpectation.input = &MemberAnnouncementSignatureMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.Equals
func (m *mMemberAnnouncementSignatureMockEquals) Return(r bool) *MemberAnnouncementSignatureMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockEqualsExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.Equals is expected once
func (m *mMemberAnnouncementSignatureMockEquals) ExpectOnce(p common.SignatureHolder) *MemberAnnouncementSignatureMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockEqualsExpectation{}
	expectation.input = &MemberAnnouncementSignatureMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockEqualsExpectation) Return(r bool) {
	e.result = &MemberAnnouncementSignatureMockEqualsResult{r}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.Equals method
func (m *mMemberAnnouncementSignatureMockEquals) Set(f func(p common.SignatureHolder) (r bool)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) Equals(p common.SignatureHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MemberAnnouncementSignatureMockEqualsInput{p}, "MemberAnnouncementSignature.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MemberAnnouncementSignatureMockEqualsInput{p}, "MemberAnnouncementSignature.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of MemberAnnouncementSignatureMock.EqualsFunc invocations
func (m *MemberAnnouncementSignatureMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.Equals invocations
func (m *MemberAnnouncementSignatureMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) EqualsFinished() bool {
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

type mMemberAnnouncementSignatureMockFixedByteSize struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockFixedByteSizeExpectation
	expectationSeries []*MemberAnnouncementSignatureMockFixedByteSizeExpectation
}

type MemberAnnouncementSignatureMockFixedByteSizeExpectation struct {
	result *MemberAnnouncementSignatureMockFixedByteSizeResult
}

type MemberAnnouncementSignatureMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of MemberAnnouncementSignature.FixedByteSize is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockFixedByteSize) Expect() *mMemberAnnouncementSignatureMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.FixedByteSize
func (m *mMemberAnnouncementSignatureMockFixedByteSize) Return(r int) *MemberAnnouncementSignatureMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.FixedByteSize is expected once
func (m *mMemberAnnouncementSignatureMockFixedByteSize) ExpectOnce() *MemberAnnouncementSignatureMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockFixedByteSizeExpectation) Return(r int) {
	e.result = &MemberAnnouncementSignatureMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.FixedByteSize method
func (m *mMemberAnnouncementSignatureMockFixedByteSize) Set(f func() (r int)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of MemberAnnouncementSignatureMock.FixedByteSizeFunc invocations
func (m *MemberAnnouncementSignatureMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.FixedByteSize invocations
func (m *MemberAnnouncementSignatureMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) FixedByteSizeFinished() bool {
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

type mMemberAnnouncementSignatureMockFoldToUint64 struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockFoldToUint64Expectation
	expectationSeries []*MemberAnnouncementSignatureMockFoldToUint64Expectation
}

type MemberAnnouncementSignatureMockFoldToUint64Expectation struct {
	result *MemberAnnouncementSignatureMockFoldToUint64Result
}

type MemberAnnouncementSignatureMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of MemberAnnouncementSignature.FoldToUint64 is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockFoldToUint64) Expect() *mMemberAnnouncementSignatureMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.FoldToUint64
func (m *mMemberAnnouncementSignatureMockFoldToUint64) Return(r uint64) *MemberAnnouncementSignatureMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.FoldToUint64 is expected once
func (m *mMemberAnnouncementSignatureMockFoldToUint64) ExpectOnce() *MemberAnnouncementSignatureMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &MemberAnnouncementSignatureMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.FoldToUint64 method
func (m *mMemberAnnouncementSignatureMockFoldToUint64) Set(f func() (r uint64)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of MemberAnnouncementSignatureMock.FoldToUint64Func invocations
func (m *MemberAnnouncementSignatureMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of MemberAnnouncementSignatureMock.FoldToUint64 invocations
func (m *MemberAnnouncementSignatureMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) FoldToUint64Finished() bool {
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

type mMemberAnnouncementSignatureMockGetSignatureMethod struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockGetSignatureMethodExpectation
	expectationSeries []*MemberAnnouncementSignatureMockGetSignatureMethodExpectation
}

type MemberAnnouncementSignatureMockGetSignatureMethodExpectation struct {
	result *MemberAnnouncementSignatureMockGetSignatureMethodResult
}

type MemberAnnouncementSignatureMockGetSignatureMethodResult struct {
	r common.SignatureMethod
}

//Expect specifies that invocation of MemberAnnouncementSignature.GetSignatureMethod is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockGetSignatureMethod) Expect() *mMemberAnnouncementSignatureMockGetSignatureMethod {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockGetSignatureMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.GetSignatureMethod
func (m *mMemberAnnouncementSignatureMockGetSignatureMethod) Return(r common.SignatureMethod) *MemberAnnouncementSignatureMock {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockGetSignatureMethodExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockGetSignatureMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.GetSignatureMethod is expected once
func (m *mMemberAnnouncementSignatureMockGetSignatureMethod) ExpectOnce() *MemberAnnouncementSignatureMockGetSignatureMethodExpectation {
	m.mock.GetSignatureMethodFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockGetSignatureMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockGetSignatureMethodExpectation) Return(r common.SignatureMethod) {
	e.result = &MemberAnnouncementSignatureMockGetSignatureMethodResult{r}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.GetSignatureMethod method
func (m *mMemberAnnouncementSignatureMockGetSignatureMethod) Set(f func() (r common.SignatureMethod)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureMethodFunc = f
	return m.mock
}

//GetSignatureMethod implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) GetSignatureMethod() (r common.SignatureMethod) {
	counter := atomic.AddUint64(&m.GetSignatureMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureMethodCounter, 1)

	if len(m.GetSignatureMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.GetSignatureMethod.")
			return
		}

		result := m.GetSignatureMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.GetSignatureMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodMock.mainExpectation != nil {

		result := m.GetSignatureMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.GetSignatureMethod")
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.GetSignatureMethod.")
		return
	}

	return m.GetSignatureMethodFunc()
}

//GetSignatureMethodMinimockCounter returns a count of MemberAnnouncementSignatureMock.GetSignatureMethodFunc invocations
func (m *MemberAnnouncementSignatureMock) GetSignatureMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodCounter)
}

//GetSignatureMethodMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.GetSignatureMethod invocations
func (m *MemberAnnouncementSignatureMock) GetSignatureMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodPreCounter)
}

//GetSignatureMethodFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) GetSignatureMethodFinished() bool {
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

type mMemberAnnouncementSignatureMockRead struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockReadExpectation
	expectationSeries []*MemberAnnouncementSignatureMockReadExpectation
}

type MemberAnnouncementSignatureMockReadExpectation struct {
	input  *MemberAnnouncementSignatureMockReadInput
	result *MemberAnnouncementSignatureMockReadResult
}

type MemberAnnouncementSignatureMockReadInput struct {
	p []byte
}

type MemberAnnouncementSignatureMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of MemberAnnouncementSignature.Read is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockRead) Expect(p []byte) *mMemberAnnouncementSignatureMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockReadExpectation{}
	}
	m.mainExpectation.input = &MemberAnnouncementSignatureMockReadInput{p}
	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.Read
func (m *mMemberAnnouncementSignatureMockRead) Return(r int, r1 error) *MemberAnnouncementSignatureMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockReadExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.Read is expected once
func (m *mMemberAnnouncementSignatureMockRead) ExpectOnce(p []byte) *MemberAnnouncementSignatureMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockReadExpectation{}
	expectation.input = &MemberAnnouncementSignatureMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockReadExpectation) Return(r int, r1 error) {
	e.result = &MemberAnnouncementSignatureMockReadResult{r, r1}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.Read method
func (m *mMemberAnnouncementSignatureMockRead) Set(f func(p []byte) (r int, r1 error)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MemberAnnouncementSignatureMockReadInput{p}, "MemberAnnouncementSignature.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MemberAnnouncementSignatureMockReadInput{p}, "MemberAnnouncementSignature.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of MemberAnnouncementSignatureMock.ReadFunc invocations
func (m *MemberAnnouncementSignatureMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.Read invocations
func (m *MemberAnnouncementSignatureMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) ReadFinished() bool {
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

type mMemberAnnouncementSignatureMockWriteTo struct {
	mock              *MemberAnnouncementSignatureMock
	mainExpectation   *MemberAnnouncementSignatureMockWriteToExpectation
	expectationSeries []*MemberAnnouncementSignatureMockWriteToExpectation
}

type MemberAnnouncementSignatureMockWriteToExpectation struct {
	input  *MemberAnnouncementSignatureMockWriteToInput
	result *MemberAnnouncementSignatureMockWriteToResult
}

type MemberAnnouncementSignatureMockWriteToInput struct {
	p io.Writer
}

type MemberAnnouncementSignatureMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of MemberAnnouncementSignature.WriteTo is expected from 1 to Infinity times
func (m *mMemberAnnouncementSignatureMockWriteTo) Expect(p io.Writer) *mMemberAnnouncementSignatureMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockWriteToExpectation{}
	}
	m.mainExpectation.input = &MemberAnnouncementSignatureMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of MemberAnnouncementSignature.WriteTo
func (m *mMemberAnnouncementSignatureMockWriteTo) Return(r int64, r1 error) *MemberAnnouncementSignatureMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &MemberAnnouncementSignatureMockWriteToExpectation{}
	}
	m.mainExpectation.result = &MemberAnnouncementSignatureMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of MemberAnnouncementSignature.WriteTo is expected once
func (m *mMemberAnnouncementSignatureMockWriteTo) ExpectOnce(p io.Writer) *MemberAnnouncementSignatureMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &MemberAnnouncementSignatureMockWriteToExpectation{}
	expectation.input = &MemberAnnouncementSignatureMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *MemberAnnouncementSignatureMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &MemberAnnouncementSignatureMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of MemberAnnouncementSignature.WriteTo method
func (m *mMemberAnnouncementSignatureMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *MemberAnnouncementSignatureMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/gcpv2/common.MemberAnnouncementSignature interface
func (m *MemberAnnouncementSignatureMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, MemberAnnouncementSignatureMockWriteToInput{p}, "MemberAnnouncementSignature.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, MemberAnnouncementSignatureMockWriteToInput{p}, "MemberAnnouncementSignature.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the MemberAnnouncementSignatureMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to MemberAnnouncementSignatureMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of MemberAnnouncementSignatureMock.WriteToFunc invocations
func (m *MemberAnnouncementSignatureMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of MemberAnnouncementSignatureMock.WriteTo invocations
func (m *MemberAnnouncementSignatureMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *MemberAnnouncementSignatureMock) WriteToFinished() bool {
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
func (m *MemberAnnouncementSignatureMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.AsBytes")
	}

	if !m.CopyOfSignatureFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.CopyOfSignature")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.FoldToUint64")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.GetSignatureMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *MemberAnnouncementSignatureMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *MemberAnnouncementSignatureMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *MemberAnnouncementSignatureMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.AsBytes")
	}

	if !m.CopyOfSignatureFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.CopyOfSignature")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.FoldToUint64")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.GetSignatureMethod")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to MemberAnnouncementSignatureMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *MemberAnnouncementSignatureMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *MemberAnnouncementSignatureMock) MinimockWait(timeout time.Duration) {
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
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.AsBytes")
			}

			if !m.CopyOfSignatureFinished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.CopyOfSignature")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.Equals")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.FoldToUint64")
			}

			if !m.GetSignatureMethodFinished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.GetSignatureMethod")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.Read")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to MemberAnnouncementSignatureMock.WriteTo")
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
func (m *MemberAnnouncementSignatureMock) AllMocksCalled() bool {

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
