package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SignatureKeyHolder" can be found in github.com/insolar/insolar/network/consensus/common
*/
import (
	"io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//SignatureKeyHolderMock implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder
type SignatureKeyHolderMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mSignatureKeyHolderMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mSignatureKeyHolderMockAsBytes

	EqualsFunc       func(p SignatureKeyHolder) (r bool)
	EqualsCounter    uint64
	EqualsPreCounter uint64
	EqualsMock       mSignatureKeyHolderMockEquals

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mSignatureKeyHolderMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mSignatureKeyHolderMockFoldToUint64

	GetSignMethodFunc       func() (r SignMethod)
	GetSignMethodCounter    uint64
	GetSignMethodPreCounter uint64
	GetSignMethodMock       mSignatureKeyHolderMockGetSignMethod

	GetSignatureKeyMethodFunc       func() (r SignatureMethod)
	GetSignatureKeyMethodCounter    uint64
	GetSignatureKeyMethodPreCounter uint64
	GetSignatureKeyMethodMock       mSignatureKeyHolderMockGetSignatureKeyMethod

	GetSignatureKeyTypeFunc       func() (r SignatureKeyType)
	GetSignatureKeyTypeCounter    uint64
	GetSignatureKeyTypePreCounter uint64
	GetSignatureKeyTypeMock       mSignatureKeyHolderMockGetSignatureKeyType

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mSignatureKeyHolderMockRead

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mSignatureKeyHolderMockWriteTo
}

//NewSignatureKeyHolderMock returns a mock for github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder
func NewSignatureKeyHolderMock(t minimock.Tester) *SignatureKeyHolderMock {
	m := &SignatureKeyHolderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mSignatureKeyHolderMockAsByteString{mock: m}
	m.AsBytesMock = mSignatureKeyHolderMockAsBytes{mock: m}
	m.EqualsMock = mSignatureKeyHolderMockEquals{mock: m}
	m.FixedByteSizeMock = mSignatureKeyHolderMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mSignatureKeyHolderMockFoldToUint64{mock: m}
	m.GetSignMethodMock = mSignatureKeyHolderMockGetSignMethod{mock: m}
	m.GetSignatureKeyMethodMock = mSignatureKeyHolderMockGetSignatureKeyMethod{mock: m}
	m.GetSignatureKeyTypeMock = mSignatureKeyHolderMockGetSignatureKeyType{mock: m}
	m.ReadMock = mSignatureKeyHolderMockRead{mock: m}
	m.WriteToMock = mSignatureKeyHolderMockWriteTo{mock: m}

	return m
}

type mSignatureKeyHolderMockAsByteString struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockAsByteStringExpectation
	expectationSeries []*SignatureKeyHolderMockAsByteStringExpectation
}

type SignatureKeyHolderMockAsByteStringExpectation struct {
	result *SignatureKeyHolderMockAsByteStringResult
}

type SignatureKeyHolderMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of SignatureKeyHolder.AsByteString is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockAsByteString) Expect() *mSignatureKeyHolderMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureKeyHolder.AsByteString
func (m *mSignatureKeyHolderMockAsByteString) Return(r string) *SignatureKeyHolderMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.AsByteString is expected once
func (m *mSignatureKeyHolderMockAsByteString) ExpectOnce() *SignatureKeyHolderMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockAsByteStringExpectation) Return(r string) {
	e.result = &SignatureKeyHolderMockAsByteStringResult{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.AsByteString method
func (m *mSignatureKeyHolderMockAsByteString) Set(f func() (r string)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of SignatureKeyHolderMock.AsByteStringFunc invocations
func (m *SignatureKeyHolderMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of SignatureKeyHolderMock.AsByteString invocations
func (m *SignatureKeyHolderMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) AsByteStringFinished() bool {
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

type mSignatureKeyHolderMockAsBytes struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockAsBytesExpectation
	expectationSeries []*SignatureKeyHolderMockAsBytesExpectation
}

type SignatureKeyHolderMockAsBytesExpectation struct {
	result *SignatureKeyHolderMockAsBytesResult
}

type SignatureKeyHolderMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of SignatureKeyHolder.AsBytes is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockAsBytes) Expect() *mSignatureKeyHolderMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureKeyHolder.AsBytes
func (m *mSignatureKeyHolderMockAsBytes) Return(r []byte) *SignatureKeyHolderMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.AsBytes is expected once
func (m *mSignatureKeyHolderMockAsBytes) ExpectOnce() *SignatureKeyHolderMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockAsBytesExpectation) Return(r []byte) {
	e.result = &SignatureKeyHolderMockAsBytesResult{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.AsBytes method
func (m *mSignatureKeyHolderMockAsBytes) Set(f func() (r []byte)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of SignatureKeyHolderMock.AsBytesFunc invocations
func (m *SignatureKeyHolderMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of SignatureKeyHolderMock.AsBytes invocations
func (m *SignatureKeyHolderMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) AsBytesFinished() bool {
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

type mSignatureKeyHolderMockEquals struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockEqualsExpectation
	expectationSeries []*SignatureKeyHolderMockEqualsExpectation
}

type SignatureKeyHolderMockEqualsExpectation struct {
	input  *SignatureKeyHolderMockEqualsInput
	result *SignatureKeyHolderMockEqualsResult
}

type SignatureKeyHolderMockEqualsInput struct {
	p SignatureKeyHolder
}

type SignatureKeyHolderMockEqualsResult struct {
	r bool
}

//Expect specifies that invocation of SignatureKeyHolder.Equals is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockEquals) Expect(p SignatureKeyHolder) *mSignatureKeyHolderMockEquals {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockEqualsExpectation{}
	}
	m.mainExpectation.input = &SignatureKeyHolderMockEqualsInput{p}
	return m
}

//Return specifies results of invocation of SignatureKeyHolder.Equals
func (m *mSignatureKeyHolderMockEquals) Return(r bool) *SignatureKeyHolderMock {
	m.mock.EqualsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockEqualsExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockEqualsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.Equals is expected once
func (m *mSignatureKeyHolderMockEquals) ExpectOnce(p SignatureKeyHolder) *SignatureKeyHolderMockEqualsExpectation {
	m.mock.EqualsFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockEqualsExpectation{}
	expectation.input = &SignatureKeyHolderMockEqualsInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockEqualsExpectation) Return(r bool) {
	e.result = &SignatureKeyHolderMockEqualsResult{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.Equals method
func (m *mSignatureKeyHolderMockEquals) Set(f func(p SignatureKeyHolder) (r bool)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.EqualsFunc = f
	return m.mock
}

//Equals implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) Equals(p SignatureKeyHolder) (r bool) {
	counter := atomic.AddUint64(&m.EqualsPreCounter, 1)
	defer atomic.AddUint64(&m.EqualsCounter, 1)

	if len(m.EqualsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.EqualsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.Equals. %v", p)
			return
		}

		input := m.EqualsMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureKeyHolderMockEqualsInput{p}, "SignatureKeyHolder.Equals got unexpected parameters")

		result := m.EqualsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.Equals")
			return
		}

		r = result.r

		return
	}

	if m.EqualsMock.mainExpectation != nil {

		input := m.EqualsMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureKeyHolderMockEqualsInput{p}, "SignatureKeyHolder.Equals got unexpected parameters")
		}

		result := m.EqualsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.Equals")
		}

		r = result.r

		return
	}

	if m.EqualsFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.Equals. %v", p)
		return
	}

	return m.EqualsFunc(p)
}

//EqualsMinimockCounter returns a count of SignatureKeyHolderMock.EqualsFunc invocations
func (m *SignatureKeyHolderMock) EqualsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsCounter)
}

//EqualsMinimockPreCounter returns the value of SignatureKeyHolderMock.Equals invocations
func (m *SignatureKeyHolderMock) EqualsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EqualsPreCounter)
}

//EqualsFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) EqualsFinished() bool {
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

type mSignatureKeyHolderMockFixedByteSize struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockFixedByteSizeExpectation
	expectationSeries []*SignatureKeyHolderMockFixedByteSizeExpectation
}

type SignatureKeyHolderMockFixedByteSizeExpectation struct {
	result *SignatureKeyHolderMockFixedByteSizeResult
}

type SignatureKeyHolderMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of SignatureKeyHolder.FixedByteSize is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockFixedByteSize) Expect() *mSignatureKeyHolderMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureKeyHolder.FixedByteSize
func (m *mSignatureKeyHolderMockFixedByteSize) Return(r int) *SignatureKeyHolderMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.FixedByteSize is expected once
func (m *mSignatureKeyHolderMockFixedByteSize) ExpectOnce() *SignatureKeyHolderMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockFixedByteSizeExpectation) Return(r int) {
	e.result = &SignatureKeyHolderMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.FixedByteSize method
func (m *mSignatureKeyHolderMockFixedByteSize) Set(f func() (r int)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of SignatureKeyHolderMock.FixedByteSizeFunc invocations
func (m *SignatureKeyHolderMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of SignatureKeyHolderMock.FixedByteSize invocations
func (m *SignatureKeyHolderMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) FixedByteSizeFinished() bool {
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

type mSignatureKeyHolderMockFoldToUint64 struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockFoldToUint64Expectation
	expectationSeries []*SignatureKeyHolderMockFoldToUint64Expectation
}

type SignatureKeyHolderMockFoldToUint64Expectation struct {
	result *SignatureKeyHolderMockFoldToUint64Result
}

type SignatureKeyHolderMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of SignatureKeyHolder.FoldToUint64 is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockFoldToUint64) Expect() *mSignatureKeyHolderMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureKeyHolder.FoldToUint64
func (m *mSignatureKeyHolderMockFoldToUint64) Return(r uint64) *SignatureKeyHolderMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.FoldToUint64 is expected once
func (m *mSignatureKeyHolderMockFoldToUint64) ExpectOnce() *SignatureKeyHolderMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &SignatureKeyHolderMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.FoldToUint64 method
func (m *mSignatureKeyHolderMockFoldToUint64) Set(f func() (r uint64)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of SignatureKeyHolderMock.FoldToUint64Func invocations
func (m *SignatureKeyHolderMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of SignatureKeyHolderMock.FoldToUint64 invocations
func (m *SignatureKeyHolderMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) FoldToUint64Finished() bool {
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

type mSignatureKeyHolderMockGetSignMethod struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockGetSignMethodExpectation
	expectationSeries []*SignatureKeyHolderMockGetSignMethodExpectation
}

type SignatureKeyHolderMockGetSignMethodExpectation struct {
	result *SignatureKeyHolderMockGetSignMethodResult
}

type SignatureKeyHolderMockGetSignMethodResult struct {
	r SignMethod
}

//Expect specifies that invocation of SignatureKeyHolder.GetSignMethod is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockGetSignMethod) Expect() *mSignatureKeyHolderMockGetSignMethod {
	m.mock.GetSignMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockGetSignMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureKeyHolder.GetSignMethod
func (m *mSignatureKeyHolderMockGetSignMethod) Return(r SignMethod) *SignatureKeyHolderMock {
	m.mock.GetSignMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockGetSignMethodExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockGetSignMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.GetSignMethod is expected once
func (m *mSignatureKeyHolderMockGetSignMethod) ExpectOnce() *SignatureKeyHolderMockGetSignMethodExpectation {
	m.mock.GetSignMethodFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockGetSignMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockGetSignMethodExpectation) Return(r SignMethod) {
	e.result = &SignatureKeyHolderMockGetSignMethodResult{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.GetSignMethod method
func (m *mSignatureKeyHolderMockGetSignMethod) Set(f func() (r SignMethod)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignMethodFunc = f
	return m.mock
}

//GetSignMethod implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) GetSignMethod() (r SignMethod) {
	counter := atomic.AddUint64(&m.GetSignMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignMethodCounter, 1)

	if len(m.GetSignMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.GetSignMethod.")
			return
		}

		result := m.GetSignMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.GetSignMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignMethodMock.mainExpectation != nil {

		result := m.GetSignMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.GetSignMethod")
		}

		r = result.r

		return
	}

	if m.GetSignMethodFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.GetSignMethod.")
		return
	}

	return m.GetSignMethodFunc()
}

//GetSignMethodMinimockCounter returns a count of SignatureKeyHolderMock.GetSignMethodFunc invocations
func (m *SignatureKeyHolderMock) GetSignMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignMethodCounter)
}

//GetSignMethodMinimockPreCounter returns the value of SignatureKeyHolderMock.GetSignMethod invocations
func (m *SignatureKeyHolderMock) GetSignMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignMethodPreCounter)
}

//GetSignMethodFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) GetSignMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignMethodCounter) == uint64(len(m.GetSignMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignMethodFunc != nil {
		return atomic.LoadUint64(&m.GetSignMethodCounter) > 0
	}

	return true
}

type mSignatureKeyHolderMockGetSignatureKeyMethod struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockGetSignatureKeyMethodExpectation
	expectationSeries []*SignatureKeyHolderMockGetSignatureKeyMethodExpectation
}

type SignatureKeyHolderMockGetSignatureKeyMethodExpectation struct {
	result *SignatureKeyHolderMockGetSignatureKeyMethodResult
}

type SignatureKeyHolderMockGetSignatureKeyMethodResult struct {
	r SignatureMethod
}

//Expect specifies that invocation of SignatureKeyHolder.GetSignatureKeyMethod is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockGetSignatureKeyMethod) Expect() *mSignatureKeyHolderMockGetSignatureKeyMethod {
	m.mock.GetSignatureKeyMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockGetSignatureKeyMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureKeyHolder.GetSignatureKeyMethod
func (m *mSignatureKeyHolderMockGetSignatureKeyMethod) Return(r SignatureMethod) *SignatureKeyHolderMock {
	m.mock.GetSignatureKeyMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockGetSignatureKeyMethodExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockGetSignatureKeyMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.GetSignatureKeyMethod is expected once
func (m *mSignatureKeyHolderMockGetSignatureKeyMethod) ExpectOnce() *SignatureKeyHolderMockGetSignatureKeyMethodExpectation {
	m.mock.GetSignatureKeyMethodFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockGetSignatureKeyMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockGetSignatureKeyMethodExpectation) Return(r SignatureMethod) {
	e.result = &SignatureKeyHolderMockGetSignatureKeyMethodResult{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.GetSignatureKeyMethod method
func (m *mSignatureKeyHolderMockGetSignatureKeyMethod) Set(f func() (r SignatureMethod)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureKeyMethodFunc = f
	return m.mock
}

//GetSignatureKeyMethod implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) GetSignatureKeyMethod() (r SignatureMethod) {
	counter := atomic.AddUint64(&m.GetSignatureKeyMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureKeyMethodCounter, 1)

	if len(m.GetSignatureKeyMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureKeyMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.GetSignatureKeyMethod.")
			return
		}

		result := m.GetSignatureKeyMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.GetSignatureKeyMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureKeyMethodMock.mainExpectation != nil {

		result := m.GetSignatureKeyMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.GetSignatureKeyMethod")
		}

		r = result.r

		return
	}

	if m.GetSignatureKeyMethodFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.GetSignatureKeyMethod.")
		return
	}

	return m.GetSignatureKeyMethodFunc()
}

//GetSignatureKeyMethodMinimockCounter returns a count of SignatureKeyHolderMock.GetSignatureKeyMethodFunc invocations
func (m *SignatureKeyHolderMock) GetSignatureKeyMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureKeyMethodCounter)
}

//GetSignatureKeyMethodMinimockPreCounter returns the value of SignatureKeyHolderMock.GetSignatureKeyMethod invocations
func (m *SignatureKeyHolderMock) GetSignatureKeyMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureKeyMethodPreCounter)
}

//GetSignatureKeyMethodFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) GetSignatureKeyMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureKeyMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureKeyMethodCounter) == uint64(len(m.GetSignatureKeyMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureKeyMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureKeyMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureKeyMethodFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureKeyMethodCounter) > 0
	}

	return true
}

type mSignatureKeyHolderMockGetSignatureKeyType struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockGetSignatureKeyTypeExpectation
	expectationSeries []*SignatureKeyHolderMockGetSignatureKeyTypeExpectation
}

type SignatureKeyHolderMockGetSignatureKeyTypeExpectation struct {
	result *SignatureKeyHolderMockGetSignatureKeyTypeResult
}

type SignatureKeyHolderMockGetSignatureKeyTypeResult struct {
	r SignatureKeyType
}

//Expect specifies that invocation of SignatureKeyHolder.GetSignatureKeyType is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockGetSignatureKeyType) Expect() *mSignatureKeyHolderMockGetSignatureKeyType {
	m.mock.GetSignatureKeyTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockGetSignatureKeyTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of SignatureKeyHolder.GetSignatureKeyType
func (m *mSignatureKeyHolderMockGetSignatureKeyType) Return(r SignatureKeyType) *SignatureKeyHolderMock {
	m.mock.GetSignatureKeyTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockGetSignatureKeyTypeExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockGetSignatureKeyTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.GetSignatureKeyType is expected once
func (m *mSignatureKeyHolderMockGetSignatureKeyType) ExpectOnce() *SignatureKeyHolderMockGetSignatureKeyTypeExpectation {
	m.mock.GetSignatureKeyTypeFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockGetSignatureKeyTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockGetSignatureKeyTypeExpectation) Return(r SignatureKeyType) {
	e.result = &SignatureKeyHolderMockGetSignatureKeyTypeResult{r}
}

//Set uses given function f as a mock of SignatureKeyHolder.GetSignatureKeyType method
func (m *mSignatureKeyHolderMockGetSignatureKeyType) Set(f func() (r SignatureKeyType)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureKeyTypeFunc = f
	return m.mock
}

//GetSignatureKeyType implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) GetSignatureKeyType() (r SignatureKeyType) {
	counter := atomic.AddUint64(&m.GetSignatureKeyTypePreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureKeyTypeCounter, 1)

	if len(m.GetSignatureKeyTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureKeyTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.GetSignatureKeyType.")
			return
		}

		result := m.GetSignatureKeyTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.GetSignatureKeyType")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureKeyTypeMock.mainExpectation != nil {

		result := m.GetSignatureKeyTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.GetSignatureKeyType")
		}

		r = result.r

		return
	}

	if m.GetSignatureKeyTypeFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.GetSignatureKeyType.")
		return
	}

	return m.GetSignatureKeyTypeFunc()
}

//GetSignatureKeyTypeMinimockCounter returns a count of SignatureKeyHolderMock.GetSignatureKeyTypeFunc invocations
func (m *SignatureKeyHolderMock) GetSignatureKeyTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureKeyTypeCounter)
}

//GetSignatureKeyTypeMinimockPreCounter returns the value of SignatureKeyHolderMock.GetSignatureKeyType invocations
func (m *SignatureKeyHolderMock) GetSignatureKeyTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureKeyTypePreCounter)
}

//GetSignatureKeyTypeFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) GetSignatureKeyTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureKeyTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureKeyTypeCounter) == uint64(len(m.GetSignatureKeyTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureKeyTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureKeyTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureKeyTypeFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureKeyTypeCounter) > 0
	}

	return true
}

type mSignatureKeyHolderMockRead struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockReadExpectation
	expectationSeries []*SignatureKeyHolderMockReadExpectation
}

type SignatureKeyHolderMockReadExpectation struct {
	input  *SignatureKeyHolderMockReadInput
	result *SignatureKeyHolderMockReadResult
}

type SignatureKeyHolderMockReadInput struct {
	p []byte
}

type SignatureKeyHolderMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of SignatureKeyHolder.Read is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockRead) Expect(p []byte) *mSignatureKeyHolderMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockReadExpectation{}
	}
	m.mainExpectation.input = &SignatureKeyHolderMockReadInput{p}
	return m
}

//Return specifies results of invocation of SignatureKeyHolder.Read
func (m *mSignatureKeyHolderMockRead) Return(r int, r1 error) *SignatureKeyHolderMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockReadExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.Read is expected once
func (m *mSignatureKeyHolderMockRead) ExpectOnce(p []byte) *SignatureKeyHolderMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockReadExpectation{}
	expectation.input = &SignatureKeyHolderMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockReadExpectation) Return(r int, r1 error) {
	e.result = &SignatureKeyHolderMockReadResult{r, r1}
}

//Set uses given function f as a mock of SignatureKeyHolder.Read method
func (m *mSignatureKeyHolderMockRead) Set(f func(p []byte) (r int, r1 error)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureKeyHolderMockReadInput{p}, "SignatureKeyHolder.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureKeyHolderMockReadInput{p}, "SignatureKeyHolder.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of SignatureKeyHolderMock.ReadFunc invocations
func (m *SignatureKeyHolderMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of SignatureKeyHolderMock.Read invocations
func (m *SignatureKeyHolderMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) ReadFinished() bool {
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

type mSignatureKeyHolderMockWriteTo struct {
	mock              *SignatureKeyHolderMock
	mainExpectation   *SignatureKeyHolderMockWriteToExpectation
	expectationSeries []*SignatureKeyHolderMockWriteToExpectation
}

type SignatureKeyHolderMockWriteToExpectation struct {
	input  *SignatureKeyHolderMockWriteToInput
	result *SignatureKeyHolderMockWriteToResult
}

type SignatureKeyHolderMockWriteToInput struct {
	p io.Writer
}

type SignatureKeyHolderMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of SignatureKeyHolder.WriteTo is expected from 1 to Infinity times
func (m *mSignatureKeyHolderMockWriteTo) Expect(p io.Writer) *mSignatureKeyHolderMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockWriteToExpectation{}
	}
	m.mainExpectation.input = &SignatureKeyHolderMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of SignatureKeyHolder.WriteTo
func (m *mSignatureKeyHolderMockWriteTo) Return(r int64, r1 error) *SignatureKeyHolderMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureKeyHolderMockWriteToExpectation{}
	}
	m.mainExpectation.result = &SignatureKeyHolderMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureKeyHolder.WriteTo is expected once
func (m *mSignatureKeyHolderMockWriteTo) ExpectOnce(p io.Writer) *SignatureKeyHolderMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureKeyHolderMockWriteToExpectation{}
	expectation.input = &SignatureKeyHolderMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureKeyHolderMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &SignatureKeyHolderMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of SignatureKeyHolder.WriteTo method
func (m *mSignatureKeyHolderMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *SignatureKeyHolderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/common.SignatureKeyHolder interface
func (m *SignatureKeyHolderMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureKeyHolderMockWriteToInput{p}, "SignatureKeyHolder.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureKeyHolderMockWriteToInput{p}, "SignatureKeyHolder.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureKeyHolderMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureKeyHolderMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of SignatureKeyHolderMock.WriteToFunc invocations
func (m *SignatureKeyHolderMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of SignatureKeyHolderMock.WriteTo invocations
func (m *SignatureKeyHolderMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *SignatureKeyHolderMock) WriteToFinished() bool {
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
func (m *SignatureKeyHolderMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.AsBytes")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.FoldToUint64")
	}

	if !m.GetSignMethodFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.GetSignMethod")
	}

	if !m.GetSignatureKeyMethodFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.GetSignatureKeyMethod")
	}

	if !m.GetSignatureKeyTypeFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.GetSignatureKeyType")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SignatureKeyHolderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SignatureKeyHolderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SignatureKeyHolderMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.AsBytes")
	}

	if !m.EqualsFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.Equals")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.FoldToUint64")
	}

	if !m.GetSignMethodFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.GetSignMethod")
	}

	if !m.GetSignatureKeyMethodFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.GetSignatureKeyMethod")
	}

	if !m.GetSignatureKeyTypeFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.GetSignatureKeyType")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to SignatureKeyHolderMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SignatureKeyHolderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SignatureKeyHolderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AsByteStringFinished()
		ok = ok && m.AsBytesFinished()
		ok = ok && m.EqualsFinished()
		ok = ok && m.FixedByteSizeFinished()
		ok = ok && m.FoldToUint64Finished()
		ok = ok && m.GetSignMethodFinished()
		ok = ok && m.GetSignatureKeyMethodFinished()
		ok = ok && m.GetSignatureKeyTypeFinished()
		ok = ok && m.ReadFinished()
		ok = ok && m.WriteToFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AsByteStringFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.AsBytes")
			}

			if !m.EqualsFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.Equals")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.FoldToUint64")
			}

			if !m.GetSignMethodFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.GetSignMethod")
			}

			if !m.GetSignatureKeyMethodFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.GetSignatureKeyMethod")
			}

			if !m.GetSignatureKeyTypeFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.GetSignatureKeyType")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.Read")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to SignatureKeyHolderMock.WriteTo")
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
func (m *SignatureKeyHolderMock) AllMocksCalled() bool {

	if !m.AsByteStringFinished() {
		return false
	}

	if !m.AsBytesFinished() {
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

	if !m.GetSignMethodFinished() {
		return false
	}

	if !m.GetSignatureKeyMethodFinished() {
		return false
	}

	if !m.GetSignatureKeyTypeFinished() {
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
