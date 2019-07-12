package long_bits

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "FixedReader" can be found in github.com/insolar/insolar/network/consensus/common/long_bits
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//FixedReaderMock implements github.com/insolar/insolar/network/consensus/common/long_bits.FixedReader
type FixedReaderMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mFixedReaderMockAsByteString

	AsBytesFunc       func() (r []byte)
	AsBytesCounter    uint64
	AsBytesPreCounter uint64
	AsBytesMock       mFixedReaderMockAsBytes

	FixedByteSizeFunc       func() (r int)
	FixedByteSizeCounter    uint64
	FixedByteSizePreCounter uint64
	FixedByteSizeMock       mFixedReaderMockFixedByteSize

	ReadFunc       func(p []byte) (r int, r1 error)
	ReadCounter    uint64
	ReadPreCounter uint64
	ReadMock       mFixedReaderMockRead

	WriteToFunc       func(p io.Writer) (r int64, r1 error)
	WriteToCounter    uint64
	WriteToPreCounter uint64
	WriteToMock       mFixedReaderMockWriteTo
}

//NewFixedReaderMock returns a mock for github.com/insolar/insolar/network/consensus/common/long_bits.FixedReader
func NewFixedReaderMock(t minimock.Tester) *FixedReaderMock {
	m := &FixedReaderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mFixedReaderMockAsByteString{mock: m}
	m.AsBytesMock = mFixedReaderMockAsBytes{mock: m}
	m.FixedByteSizeMock = mFixedReaderMockFixedByteSize{mock: m}
	m.ReadMock = mFixedReaderMockRead{mock: m}
	m.WriteToMock = mFixedReaderMockWriteTo{mock: m}

	return m
}

type mFixedReaderMockAsByteString struct {
	mock              *FixedReaderMock
	mainExpectation   *FixedReaderMockAsByteStringExpectation
	expectationSeries []*FixedReaderMockAsByteStringExpectation
}

type FixedReaderMockAsByteStringExpectation struct {
	result *FixedReaderMockAsByteStringResult
}

type FixedReaderMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of FixedReader.AsByteString is expected from 1 to Infinity times
func (m *mFixedReaderMockAsByteString) Expect() *mFixedReaderMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of FixedReader.AsByteString
func (m *mFixedReaderMockAsByteString) Return(r string) *FixedReaderMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &FixedReaderMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FixedReader.AsByteString is expected once
func (m *mFixedReaderMockAsByteString) ExpectOnce() *FixedReaderMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &FixedReaderMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FixedReaderMockAsByteStringExpectation) Return(r string) {
	e.result = &FixedReaderMockAsByteStringResult{r}
}

//Set uses given function f as a mock of FixedReader.AsByteString method
func (m *mFixedReaderMockAsByteString) Set(f func() (r string)) *FixedReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/common/long_bits.FixedReader interface
func (m *FixedReaderMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FixedReaderMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to FixedReaderMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of FixedReaderMock.AsByteStringFunc invocations
func (m *FixedReaderMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of FixedReaderMock.AsByteString invocations
func (m *FixedReaderMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *FixedReaderMock) AsByteStringFinished() bool {
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

type mFixedReaderMockAsBytes struct {
	mock              *FixedReaderMock
	mainExpectation   *FixedReaderMockAsBytesExpectation
	expectationSeries []*FixedReaderMockAsBytesExpectation
}

type FixedReaderMockAsBytesExpectation struct {
	result *FixedReaderMockAsBytesResult
}

type FixedReaderMockAsBytesResult struct {
	r []byte
}

//Expect specifies that invocation of FixedReader.AsBytes is expected from 1 to Infinity times
func (m *mFixedReaderMockAsBytes) Expect() *mFixedReaderMockAsBytes {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockAsBytesExpectation{}
	}

	return m
}

//Return specifies results of invocation of FixedReader.AsBytes
func (m *mFixedReaderMockAsBytes) Return(r []byte) *FixedReaderMock {
	m.mock.AsBytesFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockAsBytesExpectation{}
	}
	m.mainExpectation.result = &FixedReaderMockAsBytesResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FixedReader.AsBytes is expected once
func (m *mFixedReaderMockAsBytes) ExpectOnce() *FixedReaderMockAsBytesExpectation {
	m.mock.AsBytesFunc = nil
	m.mainExpectation = nil

	expectation := &FixedReaderMockAsBytesExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FixedReaderMockAsBytesExpectation) Return(r []byte) {
	e.result = &FixedReaderMockAsBytesResult{r}
}

//Set uses given function f as a mock of FixedReader.AsBytes method
func (m *mFixedReaderMockAsBytes) Set(f func() (r []byte)) *FixedReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsBytesFunc = f
	return m.mock
}

//AsBytes implements github.com/insolar/insolar/network/consensus/common/long_bits.FixedReader interface
func (m *FixedReaderMock) AsBytes() (r []byte) {
	counter := atomic.AddUint64(&m.AsBytesPreCounter, 1)
	defer atomic.AddUint64(&m.AsBytesCounter, 1)

	if len(m.AsBytesMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsBytesMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FixedReaderMock.AsBytes.")
			return
		}

		result := m.AsBytesMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.AsBytes")
			return
		}

		r = result.r

		return
	}

	if m.AsBytesMock.mainExpectation != nil {

		result := m.AsBytesMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.AsBytes")
		}

		r = result.r

		return
	}

	if m.AsBytesFunc == nil {
		m.t.Fatalf("Unexpected call to FixedReaderMock.AsBytes.")
		return
	}

	return m.AsBytesFunc()
}

//AsBytesMinimockCounter returns a count of FixedReaderMock.AsBytesFunc invocations
func (m *FixedReaderMock) AsBytesMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesCounter)
}

//AsBytesMinimockPreCounter returns the value of FixedReaderMock.AsBytes invocations
func (m *FixedReaderMock) AsBytesMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsBytesPreCounter)
}

//AsBytesFinished returns true if mock invocations count is ok
func (m *FixedReaderMock) AsBytesFinished() bool {
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

type mFixedReaderMockFixedByteSize struct {
	mock              *FixedReaderMock
	mainExpectation   *FixedReaderMockFixedByteSizeExpectation
	expectationSeries []*FixedReaderMockFixedByteSizeExpectation
}

type FixedReaderMockFixedByteSizeExpectation struct {
	result *FixedReaderMockFixedByteSizeResult
}

type FixedReaderMockFixedByteSizeResult struct {
	r int
}

//Expect specifies that invocation of FixedReader.FixedByteSize is expected from 1 to Infinity times
func (m *mFixedReaderMockFixedByteSize) Expect() *mFixedReaderMockFixedByteSize {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockFixedByteSizeExpectation{}
	}

	return m
}

//Return specifies results of invocation of FixedReader.FixedByteSize
func (m *mFixedReaderMockFixedByteSize) Return(r int) *FixedReaderMock {
	m.mock.FixedByteSizeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockFixedByteSizeExpectation{}
	}
	m.mainExpectation.result = &FixedReaderMockFixedByteSizeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of FixedReader.FixedByteSize is expected once
func (m *mFixedReaderMockFixedByteSize) ExpectOnce() *FixedReaderMockFixedByteSizeExpectation {
	m.mock.FixedByteSizeFunc = nil
	m.mainExpectation = nil

	expectation := &FixedReaderMockFixedByteSizeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FixedReaderMockFixedByteSizeExpectation) Return(r int) {
	e.result = &FixedReaderMockFixedByteSizeResult{r}
}

//Set uses given function f as a mock of FixedReader.FixedByteSize method
func (m *mFixedReaderMockFixedByteSize) Set(f func() (r int)) *FixedReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FixedByteSizeFunc = f
	return m.mock
}

//FixedByteSize implements github.com/insolar/insolar/network/consensus/common/long_bits.FixedReader interface
func (m *FixedReaderMock) FixedByteSize() (r int) {
	counter := atomic.AddUint64(&m.FixedByteSizePreCounter, 1)
	defer atomic.AddUint64(&m.FixedByteSizeCounter, 1)

	if len(m.FixedByteSizeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FixedByteSizeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FixedReaderMock.FixedByteSize.")
			return
		}

		result := m.FixedByteSizeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.FixedByteSize")
			return
		}

		r = result.r

		return
	}

	if m.FixedByteSizeMock.mainExpectation != nil {

		result := m.FixedByteSizeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.FixedByteSize")
		}

		r = result.r

		return
	}

	if m.FixedByteSizeFunc == nil {
		m.t.Fatalf("Unexpected call to FixedReaderMock.FixedByteSize.")
		return
	}

	return m.FixedByteSizeFunc()
}

//FixedByteSizeMinimockCounter returns a count of FixedReaderMock.FixedByteSizeFunc invocations
func (m *FixedReaderMock) FixedByteSizeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizeCounter)
}

//FixedByteSizeMinimockPreCounter returns the value of FixedReaderMock.FixedByteSize invocations
func (m *FixedReaderMock) FixedByteSizeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FixedByteSizePreCounter)
}

//FixedByteSizeFinished returns true if mock invocations count is ok
func (m *FixedReaderMock) FixedByteSizeFinished() bool {
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

type mFixedReaderMockRead struct {
	mock              *FixedReaderMock
	mainExpectation   *FixedReaderMockReadExpectation
	expectationSeries []*FixedReaderMockReadExpectation
}

type FixedReaderMockReadExpectation struct {
	input  *FixedReaderMockReadInput
	result *FixedReaderMockReadResult
}

type FixedReaderMockReadInput struct {
	p []byte
}

type FixedReaderMockReadResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of FixedReader.Read is expected from 1 to Infinity times
func (m *mFixedReaderMockRead) Expect(p []byte) *mFixedReaderMockRead {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockReadExpectation{}
	}
	m.mainExpectation.input = &FixedReaderMockReadInput{p}
	return m
}

//Return specifies results of invocation of FixedReader.Read
func (m *mFixedReaderMockRead) Return(r int, r1 error) *FixedReaderMock {
	m.mock.ReadFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockReadExpectation{}
	}
	m.mainExpectation.result = &FixedReaderMockReadResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FixedReader.Read is expected once
func (m *mFixedReaderMockRead) ExpectOnce(p []byte) *FixedReaderMockReadExpectation {
	m.mock.ReadFunc = nil
	m.mainExpectation = nil

	expectation := &FixedReaderMockReadExpectation{}
	expectation.input = &FixedReaderMockReadInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FixedReaderMockReadExpectation) Return(r int, r1 error) {
	e.result = &FixedReaderMockReadResult{r, r1}
}

//Set uses given function f as a mock of FixedReader.Read method
func (m *mFixedReaderMockRead) Set(f func(p []byte) (r int, r1 error)) *FixedReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ReadFunc = f
	return m.mock
}

//Read implements github.com/insolar/insolar/network/consensus/common/long_bits.FixedReader interface
func (m *FixedReaderMock) Read(p []byte) (r int, r1 error) {
	counter := atomic.AddUint64(&m.ReadPreCounter, 1)
	defer atomic.AddUint64(&m.ReadCounter, 1)

	if len(m.ReadMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ReadMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FixedReaderMock.Read. %v", p)
			return
		}

		input := m.ReadMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FixedReaderMockReadInput{p}, "FixedReader.Read got unexpected parameters")

		result := m.ReadMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.Read")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadMock.mainExpectation != nil {

		input := m.ReadMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FixedReaderMockReadInput{p}, "FixedReader.Read got unexpected parameters")
		}

		result := m.ReadMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.Read")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ReadFunc == nil {
		m.t.Fatalf("Unexpected call to FixedReaderMock.Read. %v", p)
		return
	}

	return m.ReadFunc(p)
}

//ReadMinimockCounter returns a count of FixedReaderMock.ReadFunc invocations
func (m *FixedReaderMock) ReadMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ReadCounter)
}

//ReadMinimockPreCounter returns the value of FixedReaderMock.Read invocations
func (m *FixedReaderMock) ReadMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ReadPreCounter)
}

//ReadFinished returns true if mock invocations count is ok
func (m *FixedReaderMock) ReadFinished() bool {
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

type mFixedReaderMockWriteTo struct {
	mock              *FixedReaderMock
	mainExpectation   *FixedReaderMockWriteToExpectation
	expectationSeries []*FixedReaderMockWriteToExpectation
}

type FixedReaderMockWriteToExpectation struct {
	input  *FixedReaderMockWriteToInput
	result *FixedReaderMockWriteToResult
}

type FixedReaderMockWriteToInput struct {
	p io.Writer
}

type FixedReaderMockWriteToResult struct {
	r  int64
	r1 error
}

//Expect specifies that invocation of FixedReader.WriteTo is expected from 1 to Infinity times
func (m *mFixedReaderMockWriteTo) Expect(p io.Writer) *mFixedReaderMockWriteTo {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockWriteToExpectation{}
	}
	m.mainExpectation.input = &FixedReaderMockWriteToInput{p}
	return m
}

//Return specifies results of invocation of FixedReader.WriteTo
func (m *mFixedReaderMockWriteTo) Return(r int64, r1 error) *FixedReaderMock {
	m.mock.WriteToFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &FixedReaderMockWriteToExpectation{}
	}
	m.mainExpectation.result = &FixedReaderMockWriteToResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of FixedReader.WriteTo is expected once
func (m *mFixedReaderMockWriteTo) ExpectOnce(p io.Writer) *FixedReaderMockWriteToExpectation {
	m.mock.WriteToFunc = nil
	m.mainExpectation = nil

	expectation := &FixedReaderMockWriteToExpectation{}
	expectation.input = &FixedReaderMockWriteToInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *FixedReaderMockWriteToExpectation) Return(r int64, r1 error) {
	e.result = &FixedReaderMockWriteToResult{r, r1}
}

//Set uses given function f as a mock of FixedReader.WriteTo method
func (m *mFixedReaderMockWriteTo) Set(f func(p io.Writer) (r int64, r1 error)) *FixedReaderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteToFunc = f
	return m.mock
}

//WriteTo implements github.com/insolar/insolar/network/consensus/common/long_bits.FixedReader interface
func (m *FixedReaderMock) WriteTo(p io.Writer) (r int64, r1 error) {
	counter := atomic.AddUint64(&m.WriteToPreCounter, 1)
	defer atomic.AddUint64(&m.WriteToCounter, 1)

	if len(m.WriteToMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteToMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to FixedReaderMock.WriteTo. %v", p)
			return
		}

		input := m.WriteToMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, FixedReaderMockWriteToInput{p}, "FixedReader.WriteTo got unexpected parameters")

		result := m.WriteToMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.WriteTo")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToMock.mainExpectation != nil {

		input := m.WriteToMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, FixedReaderMockWriteToInput{p}, "FixedReader.WriteTo got unexpected parameters")
		}

		result := m.WriteToMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the FixedReaderMock.WriteTo")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.WriteToFunc == nil {
		m.t.Fatalf("Unexpected call to FixedReaderMock.WriteTo. %v", p)
		return
	}

	return m.WriteToFunc(p)
}

//WriteToMinimockCounter returns a count of FixedReaderMock.WriteToFunc invocations
func (m *FixedReaderMock) WriteToMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToCounter)
}

//WriteToMinimockPreCounter returns the value of FixedReaderMock.WriteTo invocations
func (m *FixedReaderMock) WriteToMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteToPreCounter)
}

//WriteToFinished returns true if mock invocations count is ok
func (m *FixedReaderMock) WriteToFinished() bool {
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
func (m *FixedReaderMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.AsBytes")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.FixedByteSize")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.WriteTo")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *FixedReaderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *FixedReaderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FixedReaderMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.AsByteString")
	}

	if !m.AsBytesFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.AsBytes")
	}

	if !m.FixedByteSizeFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.FixedByteSize")
	}

	if !m.ReadFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.Read")
	}

	if !m.WriteToFinished() {
		m.t.Fatal("Expected call to FixedReaderMock.WriteTo")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *FixedReaderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FixedReaderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AsByteStringFinished()
		ok = ok && m.AsBytesFinished()
		ok = ok && m.FixedByteSizeFinished()
		ok = ok && m.ReadFinished()
		ok = ok && m.WriteToFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AsByteStringFinished() {
				m.t.Error("Expected call to FixedReaderMock.AsByteString")
			}

			if !m.AsBytesFinished() {
				m.t.Error("Expected call to FixedReaderMock.AsBytes")
			}

			if !m.FixedByteSizeFinished() {
				m.t.Error("Expected call to FixedReaderMock.FixedByteSize")
			}

			if !m.ReadFinished() {
				m.t.Error("Expected call to FixedReaderMock.Read")
			}

			if !m.WriteToFinished() {
				m.t.Error("Expected call to FixedReaderMock.WriteTo")
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
func (m *FixedReaderMock) AllMocksCalled() bool {

	if !m.AsByteStringFinished() {
		return false
	}

	if !m.AsBytesFinished() {
		return false
	}

	if !m.FixedByteSizeFinished() {
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
