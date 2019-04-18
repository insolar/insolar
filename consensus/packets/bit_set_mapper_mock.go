package packets

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "BitSetMapper" can be found in github.com/insolar/insolar/consensus/packets
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//BitSetMapperMock implements github.com/insolar/insolar/consensus/packets.BitSetMapper
type BitSetMapperMock struct {
	t minimock.Tester

	IndexToRefFunc       func(p int) (r insolar.Reference, r1 error)
	IndexToRefCounter    uint64
	IndexToRefPreCounter uint64
	IndexToRefMock       mBitSetMapperMockIndexToRef

	LengthFunc       func() (r int)
	LengthCounter    uint64
	LengthPreCounter uint64
	LengthMock       mBitSetMapperMockLength

	RefToIndexFunc       func(p insolar.Reference) (r int, r1 error)
	RefToIndexCounter    uint64
	RefToIndexPreCounter uint64
	RefToIndexMock       mBitSetMapperMockRefToIndex
}

//NewBitSetMapperMock returns a mock for github.com/insolar/insolar/consensus/packets.BitSetMapper
func NewBitSetMapperMock(t minimock.Tester) *BitSetMapperMock {
	m := &BitSetMapperMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IndexToRefMock = mBitSetMapperMockIndexToRef{mock: m}
	m.LengthMock = mBitSetMapperMockLength{mock: m}
	m.RefToIndexMock = mBitSetMapperMockRefToIndex{mock: m}

	return m
}

type mBitSetMapperMockIndexToRef struct {
	mock              *BitSetMapperMock
	mainExpectation   *BitSetMapperMockIndexToRefExpectation
	expectationSeries []*BitSetMapperMockIndexToRefExpectation
}

type BitSetMapperMockIndexToRefExpectation struct {
	input  *BitSetMapperMockIndexToRefInput
	result *BitSetMapperMockIndexToRefResult
}

type BitSetMapperMockIndexToRefInput struct {
	p int
}

type BitSetMapperMockIndexToRefResult struct {
	r  insolar.Reference
	r1 error
}

//Expect specifies that invocation of BitSetMapper.IndexToRef is expected from 1 to Infinity times
func (m *mBitSetMapperMockIndexToRef) Expect(p int) *mBitSetMapperMockIndexToRef {
	m.mock.IndexToRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BitSetMapperMockIndexToRefExpectation{}
	}
	m.mainExpectation.input = &BitSetMapperMockIndexToRefInput{p}
	return m
}

//Return specifies results of invocation of BitSetMapper.IndexToRef
func (m *mBitSetMapperMockIndexToRef) Return(r insolar.Reference, r1 error) *BitSetMapperMock {
	m.mock.IndexToRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BitSetMapperMockIndexToRefExpectation{}
	}
	m.mainExpectation.result = &BitSetMapperMockIndexToRefResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of BitSetMapper.IndexToRef is expected once
func (m *mBitSetMapperMockIndexToRef) ExpectOnce(p int) *BitSetMapperMockIndexToRefExpectation {
	m.mock.IndexToRefFunc = nil
	m.mainExpectation = nil

	expectation := &BitSetMapperMockIndexToRefExpectation{}
	expectation.input = &BitSetMapperMockIndexToRefInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BitSetMapperMockIndexToRefExpectation) Return(r insolar.Reference, r1 error) {
	e.result = &BitSetMapperMockIndexToRefResult{r, r1}
}

//Set uses given function f as a mock of BitSetMapper.IndexToRef method
func (m *mBitSetMapperMockIndexToRef) Set(f func(p int) (r insolar.Reference, r1 error)) *BitSetMapperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IndexToRefFunc = f
	return m.mock
}

//IndexToRef implements github.com/insolar/insolar/consensus/packets.BitSetMapper interface
func (m *BitSetMapperMock) IndexToRef(p int) (r insolar.Reference, r1 error) {
	counter := atomic.AddUint64(&m.IndexToRefPreCounter, 1)
	defer atomic.AddUint64(&m.IndexToRefCounter, 1)

	if len(m.IndexToRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IndexToRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BitSetMapperMock.IndexToRef. %v", p)
			return
		}

		input := m.IndexToRefMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, BitSetMapperMockIndexToRefInput{p}, "BitSetMapper.IndexToRef got unexpected parameters")

		result := m.IndexToRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BitSetMapperMock.IndexToRef")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IndexToRefMock.mainExpectation != nil {

		input := m.IndexToRefMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, BitSetMapperMockIndexToRefInput{p}, "BitSetMapper.IndexToRef got unexpected parameters")
		}

		result := m.IndexToRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BitSetMapperMock.IndexToRef")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.IndexToRefFunc == nil {
		m.t.Fatalf("Unexpected call to BitSetMapperMock.IndexToRef. %v", p)
		return
	}

	return m.IndexToRefFunc(p)
}

//IndexToRefMinimockCounter returns a count of BitSetMapperMock.IndexToRefFunc invocations
func (m *BitSetMapperMock) IndexToRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IndexToRefCounter)
}

//IndexToRefMinimockPreCounter returns the value of BitSetMapperMock.IndexToRef invocations
func (m *BitSetMapperMock) IndexToRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IndexToRefPreCounter)
}

//IndexToRefFinished returns true if mock invocations count is ok
func (m *BitSetMapperMock) IndexToRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IndexToRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IndexToRefCounter) == uint64(len(m.IndexToRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IndexToRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IndexToRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IndexToRefFunc != nil {
		return atomic.LoadUint64(&m.IndexToRefCounter) > 0
	}

	return true
}

type mBitSetMapperMockLength struct {
	mock              *BitSetMapperMock
	mainExpectation   *BitSetMapperMockLengthExpectation
	expectationSeries []*BitSetMapperMockLengthExpectation
}

type BitSetMapperMockLengthExpectation struct {
	result *BitSetMapperMockLengthResult
}

type BitSetMapperMockLengthResult struct {
	r int
}

//Expect specifies that invocation of BitSetMapper.Length is expected from 1 to Infinity times
func (m *mBitSetMapperMockLength) Expect() *mBitSetMapperMockLength {
	m.mock.LengthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BitSetMapperMockLengthExpectation{}
	}

	return m
}

//Return specifies results of invocation of BitSetMapper.Length
func (m *mBitSetMapperMockLength) Return(r int) *BitSetMapperMock {
	m.mock.LengthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BitSetMapperMockLengthExpectation{}
	}
	m.mainExpectation.result = &BitSetMapperMockLengthResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of BitSetMapper.Length is expected once
func (m *mBitSetMapperMockLength) ExpectOnce() *BitSetMapperMockLengthExpectation {
	m.mock.LengthFunc = nil
	m.mainExpectation = nil

	expectation := &BitSetMapperMockLengthExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BitSetMapperMockLengthExpectation) Return(r int) {
	e.result = &BitSetMapperMockLengthResult{r}
}

//Set uses given function f as a mock of BitSetMapper.Length method
func (m *mBitSetMapperMockLength) Set(f func() (r int)) *BitSetMapperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LengthFunc = f
	return m.mock
}

//Length implements github.com/insolar/insolar/consensus/packets.BitSetMapper interface
func (m *BitSetMapperMock) Length() (r int) {
	counter := atomic.AddUint64(&m.LengthPreCounter, 1)
	defer atomic.AddUint64(&m.LengthCounter, 1)

	if len(m.LengthMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LengthMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BitSetMapperMock.Length.")
			return
		}

		result := m.LengthMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BitSetMapperMock.Length")
			return
		}

		r = result.r

		return
	}

	if m.LengthMock.mainExpectation != nil {

		result := m.LengthMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BitSetMapperMock.Length")
		}

		r = result.r

		return
	}

	if m.LengthFunc == nil {
		m.t.Fatalf("Unexpected call to BitSetMapperMock.Length.")
		return
	}

	return m.LengthFunc()
}

//LengthMinimockCounter returns a count of BitSetMapperMock.LengthFunc invocations
func (m *BitSetMapperMock) LengthMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LengthCounter)
}

//LengthMinimockPreCounter returns the value of BitSetMapperMock.Length invocations
func (m *BitSetMapperMock) LengthMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LengthPreCounter)
}

//LengthFinished returns true if mock invocations count is ok
func (m *BitSetMapperMock) LengthFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LengthMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LengthCounter) == uint64(len(m.LengthMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LengthMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LengthCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LengthFunc != nil {
		return atomic.LoadUint64(&m.LengthCounter) > 0
	}

	return true
}

type mBitSetMapperMockRefToIndex struct {
	mock              *BitSetMapperMock
	mainExpectation   *BitSetMapperMockRefToIndexExpectation
	expectationSeries []*BitSetMapperMockRefToIndexExpectation
}

type BitSetMapperMockRefToIndexExpectation struct {
	input  *BitSetMapperMockRefToIndexInput
	result *BitSetMapperMockRefToIndexResult
}

type BitSetMapperMockRefToIndexInput struct {
	p insolar.Reference
}

type BitSetMapperMockRefToIndexResult struct {
	r  int
	r1 error
}

//Expect specifies that invocation of BitSetMapper.RefToIndex is expected from 1 to Infinity times
func (m *mBitSetMapperMockRefToIndex) Expect(p insolar.Reference) *mBitSetMapperMockRefToIndex {
	m.mock.RefToIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BitSetMapperMockRefToIndexExpectation{}
	}
	m.mainExpectation.input = &BitSetMapperMockRefToIndexInput{p}
	return m
}

//Return specifies results of invocation of BitSetMapper.RefToIndex
func (m *mBitSetMapperMockRefToIndex) Return(r int, r1 error) *BitSetMapperMock {
	m.mock.RefToIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BitSetMapperMockRefToIndexExpectation{}
	}
	m.mainExpectation.result = &BitSetMapperMockRefToIndexResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of BitSetMapper.RefToIndex is expected once
func (m *mBitSetMapperMockRefToIndex) ExpectOnce(p insolar.Reference) *BitSetMapperMockRefToIndexExpectation {
	m.mock.RefToIndexFunc = nil
	m.mainExpectation = nil

	expectation := &BitSetMapperMockRefToIndexExpectation{}
	expectation.input = &BitSetMapperMockRefToIndexInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BitSetMapperMockRefToIndexExpectation) Return(r int, r1 error) {
	e.result = &BitSetMapperMockRefToIndexResult{r, r1}
}

//Set uses given function f as a mock of BitSetMapper.RefToIndex method
func (m *mBitSetMapperMockRefToIndex) Set(f func(p insolar.Reference) (r int, r1 error)) *BitSetMapperMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RefToIndexFunc = f
	return m.mock
}

//RefToIndex implements github.com/insolar/insolar/consensus/packets.BitSetMapper interface
func (m *BitSetMapperMock) RefToIndex(p insolar.Reference) (r int, r1 error) {
	counter := atomic.AddUint64(&m.RefToIndexPreCounter, 1)
	defer atomic.AddUint64(&m.RefToIndexCounter, 1)

	if len(m.RefToIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RefToIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BitSetMapperMock.RefToIndex. %v", p)
			return
		}

		input := m.RefToIndexMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, BitSetMapperMockRefToIndexInput{p}, "BitSetMapper.RefToIndex got unexpected parameters")

		result := m.RefToIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BitSetMapperMock.RefToIndex")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RefToIndexMock.mainExpectation != nil {

		input := m.RefToIndexMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, BitSetMapperMockRefToIndexInput{p}, "BitSetMapper.RefToIndex got unexpected parameters")
		}

		result := m.RefToIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BitSetMapperMock.RefToIndex")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RefToIndexFunc == nil {
		m.t.Fatalf("Unexpected call to BitSetMapperMock.RefToIndex. %v", p)
		return
	}

	return m.RefToIndexFunc(p)
}

//RefToIndexMinimockCounter returns a count of BitSetMapperMock.RefToIndexFunc invocations
func (m *BitSetMapperMock) RefToIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RefToIndexCounter)
}

//RefToIndexMinimockPreCounter returns the value of BitSetMapperMock.RefToIndex invocations
func (m *BitSetMapperMock) RefToIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RefToIndexPreCounter)
}

//RefToIndexFinished returns true if mock invocations count is ok
func (m *BitSetMapperMock) RefToIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RefToIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RefToIndexCounter) == uint64(len(m.RefToIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RefToIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RefToIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RefToIndexFunc != nil {
		return atomic.LoadUint64(&m.RefToIndexCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BitSetMapperMock) ValidateCallCounters() {

	if !m.IndexToRefFinished() {
		m.t.Fatal("Expected call to BitSetMapperMock.IndexToRef")
	}

	if !m.LengthFinished() {
		m.t.Fatal("Expected call to BitSetMapperMock.Length")
	}

	if !m.RefToIndexFinished() {
		m.t.Fatal("Expected call to BitSetMapperMock.RefToIndex")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BitSetMapperMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *BitSetMapperMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *BitSetMapperMock) MinimockFinish() {

	if !m.IndexToRefFinished() {
		m.t.Fatal("Expected call to BitSetMapperMock.IndexToRef")
	}

	if !m.LengthFinished() {
		m.t.Fatal("Expected call to BitSetMapperMock.Length")
	}

	if !m.RefToIndexFinished() {
		m.t.Fatal("Expected call to BitSetMapperMock.RefToIndex")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *BitSetMapperMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *BitSetMapperMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.IndexToRefFinished()
		ok = ok && m.LengthFinished()
		ok = ok && m.RefToIndexFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IndexToRefFinished() {
				m.t.Error("Expected call to BitSetMapperMock.IndexToRef")
			}

			if !m.LengthFinished() {
				m.t.Error("Expected call to BitSetMapperMock.Length")
			}

			if !m.RefToIndexFinished() {
				m.t.Error("Expected call to BitSetMapperMock.RefToIndex")
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
func (m *BitSetMapperMock) AllMocksCalled() bool {

	if !m.IndexToRefFinished() {
		return false
	}

	if !m.LengthFinished() {
		return false
	}

	if !m.RefToIndexFinished() {
		return false
	}

	return true
}
