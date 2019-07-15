package longbits

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FoldableReader" can be found in github.com/insolar/insolar/network/consensus/common
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//FoldableReaderMock implements github.com/insolar/insolar/network/consensus/common.FoldableReader
type FoldableReaderMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mFoldableReaderMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mFoldableReaderMockAsBytes

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mFoldableReaderMockFixedByteSize

	FoldToUint64Func       func() (r uint64)
	FoldToUint64Counter    uint64
	FoldToUint64PreCounter uint64
	FoldToUint64Mock       mFoldableReaderMockFoldToUint64

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mFoldableReaderMockRead

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mFoldableReaderMockWriteTo
}

//NewFoldableReaderMock returns a mock for github.com/insolar/insolar/network/consensus/common.FoldableReader
func NewFoldableReaderMock(t minimock.Tester) *FoldableReaderMock {
	m := &FoldableReaderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mFoldableReaderMockAsByteString{mock: m}
	m.AsBytesMock = mFoldableReaderMockAsBytes{mock: m}
	m.FixedByteSizeMock = mFoldableReaderMockFixedByteSize{mock: m}
	m.FoldToUint64Mock = mFoldableReaderMockFoldToUint64{mock: m}
	m.ReadMock = mFoldableReaderMockRead{mock: m}
	m.WriteToMock = mFoldableReaderMockWriteTo{mock: m}

	return m
}

type mFoldableReaderMockAsByteString struct {
	mock              *FoldableReaderMock
	mainExpectation   *FoldableReaderMockAsByteStringExpectation
	expectationSeries []*FoldableReaderMockAsByteStringExpectation
}

type FoldableReaderMockAsByteStringExpectation struct {
	result *FoldableReaderMockAsByteStringResult
}

type FoldableReaderMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of FoldableReader.AsByteString is expected from 1 to Infinity times
func (m *mFoldableReaderMockAsByteString) Expect() *mFoldableReaderMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of FoldableReader.AsByteString
func (m *mFoldableReaderMockAsByteString) Return(r string) *FoldableReaderMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &FoldableReaderMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FoldableReader.AsByteString is expected once
func (m *mFoldableReaderMockAsByteString) ExpectOnce() *FoldableReaderMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &FoldableReaderMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FoldableReaderMockAsByteStringExpectation) Return(r string) {
	e.result = &FoldableReaderMockAsByteStringResult{r}
}

//Set uses given function f as a mock of FoldableReader.AsByteString method
func (m *mFoldableReaderMockAsByteString) Set(f func() (r string)) *FoldableReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/common.FoldableReader interface
func (m *FoldableReaderMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FoldableReaderMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to FoldableReaderMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of FoldableReaderMock.AsByteStringFunc invocations
func (m *FoldableReaderMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of FoldableReaderMock.AsByteString invocations
func (m *FoldableReaderMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *FoldableReaderMock) AsByteStringFinished() bool {
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

type mFoldableReaderMockAsBytes struct {
	mock              *FoldableReaderMock
	mainExpectation   *FoldableReaderMockAsBytesExpectation
	expectationSeries []*FoldableReaderMockAsBytesExpectation
}

type FoldableReaderMockAsBytesExpectation struct {
	result *FoldableReaderMockAsBytesResult
}

type FoldableReaderMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of FoldableReader.AsBytes is expected from 1 to Infinity times
func (m *mFoldableReaderMockAsBytes) Expect() *mFoldableReaderMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of FoldableReader.AsBytes
func (m *mFoldableReaderMockAsBytes) Return(r []byte) *FoldableReaderMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &FoldableReaderMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FoldableReader.AsBytes is expected once
func (m *mFoldableReaderMockAsBytes) ExpectOnce() *FoldableReaderMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &FoldableReaderMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FoldableReaderMockAsBytesExpectation) Return(r []byte) {
	e.result = &FoldableReaderMockAsBytesResult{r}
}

//Set uses given function f as a mock of FoldableReader.AsBytes method
func (m *mFoldableReaderMockAsBytes) Set(f func() (r []byte)) *FoldableReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/common.FoldableReader interface
func (m *FoldableReaderMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FoldableReaderMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to FoldableReaderMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of FoldableReaderMock.AsBytesFunc invocations
func (m *FoldableReaderMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of FoldableReaderMock.AsBytes invocations
func (m *FoldableReaderMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *FoldableReaderMock) AsBytesFinished() bool {
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

type mFoldableReaderMockFixedByteSize struct {
	mock              *FoldableReaderMock
	mainExpectation   *FoldableReaderMockFixedByteSizeExpectation
	expectationSeries []*FoldableReaderMockFixedByteSizeExpectation
}

type FoldableReaderMockFixedByteSizeExpectation struct {
	result *FoldableReaderMockFixedByteSizeResult
}

type FoldableReaderMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of FoldableReader.FixedByteSize is expected from 1 to Infinity times
func (m *mFoldableReaderMockFixedByteSize) Expect() *mFoldableReaderMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of FoldableReader.FixedByteSize
func (m *mFoldableReaderMockFixedByteSize) Return(r int) *FoldableReaderMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &FoldableReaderMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FoldableReader.FixedByteSize is expected once
func (m *mFoldableReaderMockFixedByteSize) ExpectOnce() *FoldableReaderMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &FoldableReaderMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FoldableReaderMockFixedByteSizeExpectation) Return(r int) {
	e.result = &FoldableReaderMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of FoldableReader.FixedByteSize method
func (m *mFoldableReaderMockFixedByteSize) Set(f func() (r int)) *FoldableReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/common.FoldableReader interface
func (m *FoldableReaderMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FoldableReaderMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to FoldableReaderMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of FoldableReaderMock.FixedByteSizeFunc invocations
func (m *FoldableReaderMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of FoldableReaderMock.FixedByteSize invocations
func (m *FoldableReaderMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *FoldableReaderMock) FixedByteSizeFinished() bool {
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

type mFoldableReaderMockFoldToUint64 struct {
	mock              *FoldableReaderMock
	mainExpectation   *FoldableReaderMockFoldToUint64Expectation
	expectationSeries []*FoldableReaderMockFoldToUint64Expectation
}

type FoldableReaderMockFoldToUint64Expectation struct {
	result *FoldableReaderMockFoldToUint64Result
}

type FoldableReaderMockFoldToUint64Result struct {
	r uint64
}

//Expect specifies that invocation of FoldableReader.FoldToUint64 is expected from 1 to Infinity times
func (m *mFoldableReaderMockFoldToUint64) Expect() *mFoldableReaderMockFoldToUint64 {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockFoldToUint64Expectation{}
	}

	return m
}

//Return specifies results of invocation of FoldableReader.FoldToUint64
func (m *mFoldableReaderMockFoldToUint64) Return(r uint64) *FoldableReaderMock {
	m.mock.FoldToUint64Func = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockFoldToUint64Expectation{}
	}
	m.mainExpectation.result = &FoldableReaderMockFoldToUint64Result{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FoldableReader.FoldToUint64 is expected once
func (m *mFoldableReaderMockFoldToUint64) ExpectOnce() *FoldableReaderMockFoldToUint64Expectation {
	m.mock.FoldToUint64Func = nil
	m.mainExpectation = nil

	expectation := &FoldableReaderMockFoldToUint64Expectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FoldableReaderMockFoldToUint64Expectation) Return(r uint64) {
	e.result = &FoldableReaderMockFoldToUint64Result{r}
}

//Set uses given function f as a mock of FoldableReader.FoldToUint64 method
func (m *mFoldableReaderMockFoldToUint64) Set(f func() (r uint64)) *FoldableReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FoldToUint64Func = f
	return m.mock
}

//FoldToUint64 implements github.com/insolar/insolar/network/consensus/common.FoldableReader interface
func (m *FoldableReaderMock) FoldToUint64() (r uint64) {
	counter := atomic.AddUint64(&m.FoldToUint64PreCounter, 1)
	defer atomic.AddUint64(&m.FoldToUint64Counter, 1)

	if len(m.FoldToUint64Mock.expectationSeries) > 0 {
		if counter > uint64(len(m.FoldToUint64Mock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FoldableReaderMock.FoldToUint64.")
			return
		}

		result := m.FoldToUint64Mock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.FoldToUint64")
			return
		}

		r = result.r

		return
	}

	if m.FoldToUint64Mock.mainExpectation != nil {

		result := m.FoldToUint64Mock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.FoldToUint64")
		}

		r = result.r

		return
	}

	if m.FoldToUint64Func == nil {
		m.t.Fatalf("Unexpected call to FoldableReaderMock.FoldToUint64.")
		return
	}

	return m.FoldToUint64Func()
}

//FoldToUint64MinimockCounter returns a count of FoldableReaderMock.FoldToUint64Func invocations
func (m *FoldableReaderMock) FoldToUint64MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64Counter)
}

//FoldToUint64MinimockPreCounter returns the value of FoldableReaderMock.FoldToUint64 invocations
func (m *FoldableReaderMock) FoldToUint64MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FoldToUint64PreCounter)
}

//FoldToUint64Finished returns true if mock invocations count is ok
func (m *FoldableReaderMock) FoldToUint64Finished() bool {
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

type mFoldableReaderMockRead struct {
	mock              *FoldableReaderMock
	mainExpectation   *FoldableReaderMockReadExpectation
	expectationSeries []*FoldableReaderMockReadExpectation
}

type FoldableReaderMockReadExpectation struct {
	input  *FoldableReaderMockReadInput
	result *FoldableReaderMockReadResult
}

type FoldableReaderMockReadInput struct {
	p []byte
}

type FoldableReaderMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of FoldableReader.Read is expected from 1 to Infinity times
func (m *mFoldableReaderMockRead) Expect(p []byte) *mFoldableReaderMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockReadExpectation{}
	}
	m.mainExpectation.input = &FoldableReaderMockReadInput{p}
	return m
}

//Return specifies results of invocation of FoldableReader.Read
func (m *mFoldableReaderMockRead) Return(r int, r1 error) *FoldableReaderMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockReadExpectation{}
	}
	m.mainExpectation.result = &FoldableReaderMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FoldableReader.Read is expected once
func (m *mFoldableReaderMockRead) ExpectOnce(p []byte) *FoldableReaderMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &FoldableReaderMockReadExpectation{}
	expectation.input = &FoldableReaderMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FoldableReaderMockReadExpectation) Return(r int, r1 error) {
	e.result = &FoldableReaderMockReadResult{r, r1}
}

//Set uses given function f as a mock of FoldableReader.Read method
func (m *mFoldableReaderMockRead) Set(f func(p []byte) (r int, r1 error)) *FoldableReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/common.FoldableReader interface
func (m *FoldableReaderMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FoldableReaderMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FoldableReaderMockReadInput{p}, "FoldableReader.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FoldableReaderMockReadInput{p}, "FoldableReader.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to FoldableReaderMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of FoldableReaderMock.ReadFunc invocations
func (m *FoldableReaderMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of FoldableReaderMock.Read invocations
func (m *FoldableReaderMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *FoldableReaderMock) ReadFinished() bool {
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

type mFoldableReaderMockWriteTo struct {
	mock              *FoldableReaderMock
	mainExpectation   *FoldableReaderMockWriteToExpectation
	expectationSeries []*FoldableReaderMockWriteToExpectation
}

type FoldableReaderMockWriteToExpectation struct {
	input  *FoldableReaderMockWriteToInput
	result *FoldableReaderMockWriteToResult
}

type FoldableReaderMockWriteToInput struct {
	p io.Writer
}

type FoldableReaderMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of FoldableReader.WriteTo is expected from 1 to Infinity times
func (m *mFoldableReaderMockWriteTo) Expect(p io.Writer) *mFoldableReaderMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockWriteToExpectation{}
	}
	m.mainExpectation.input = &FoldableReaderMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of FoldableReader.WriteTo
func (m *mFoldableReaderMockWriteTo) Return(r int64, r1 error) *FoldableReaderMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FoldableReaderMockWriteToExpectation{}
	}
	m.mainExpectation.result = &FoldableReaderMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FoldableReader.WriteTo is expected once
func (m *mFoldableReaderMockWriteTo) ExpectOnce(p io.Writer) *FoldableReaderMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &FoldableReaderMockWriteToExpectation{}
	expectation.input = &FoldableReaderMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FoldableReaderMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &FoldableReaderMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of FoldableReader.WriteTo method
func (m *mFoldableReaderMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *FoldableReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/common.FoldableReader interface
func (m *FoldableReaderMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FoldableReaderMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FoldableReaderMockWriteToInput{p}, "FoldableReader.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FoldableReaderMockWriteToInput{p}, "FoldableReader.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FoldableReaderMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to FoldableReaderMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of FoldableReaderMock.WriteToFunc invocations
func (m *FoldableReaderMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of FoldableReaderMock.WriteTo invocations
func (m *FoldableReaderMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *FoldableReaderMock) WriteToFinished() bool {
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
func (m *FoldableReaderMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.AsBytes")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to FoldableReaderMock.FoldToUint64")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FoldableReaderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FoldableReaderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FoldableReaderMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.AsBytes")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.FixedByteSize")
	}

	if !m.FoldToUint64Finished() {
		m.t.Fatal("Expected call to FoldableReaderMock.FoldToUint64")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to FoldableReaderMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FoldableReaderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FoldableReaderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AsByteStringFinished()
		ok = ok && m.AsBytesFinished()
		ok = ok && m.FixedByteSizeFinished()
		ok = ok && m.FoldToUint64Finished()
		ok = ok && m.ReadFinished()
		ok = ok && m.WriteToFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AsByteStringFinished() {
				m.t.Error("Expected call to FoldableReaderMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to FoldableReaderMock.AsBytes")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to FoldableReaderMock.FixedByteSize")
			}

			if !m.FoldToUint64Finished() {
				m.t.Error("Expected call to FoldableReaderMock.FoldToUint64")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to FoldableReaderMock.Read")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to FoldableReaderMock.WriteTo")
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
func (m *FoldableReaderMock) AllMocksCalled() bool {

	if !m.AsByteStringFinished() {
		return false
	}

	if !m.AsBytesFinished() {
		return false
	}

	if !m.FixedByteSizeFinished() {
		return false
	}

	if !m.FoldToUint64Finished() {
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
