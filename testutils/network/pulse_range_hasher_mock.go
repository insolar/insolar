package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulseRangeHasher" can be found in github.com/insolar/insolar/network/storage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulseRangeHasherMock implements github.com/insolar/insolar/network/storage.PulseRangeHasher
type PulseRangeHasherMock struct {
	t minimock.Tester

	GetRangeHashFunc       func(p insolar.PulseRange) (r []byte, r1 error)
	GetRangeHashCounter    uint64
	GetRangeHashPreCounter uint64
	GetRangeHashMock       mPulseRangeHasherMockGetRangeHash

	ValidateRangeHashFunc       func(p insolar.PulseRange, p1 []byte) (r bool, r1 error)
	ValidateRangeHashCounter    uint64
	ValidateRangeHashPreCounter uint64
	ValidateRangeHashMock       mPulseRangeHasherMockValidateRangeHash
}

//NewPulseRangeHasherMock returns a mock for github.com/insolar/insolar/network/storage.PulseRangeHasher
func NewPulseRangeHasherMock(t minimock.Tester) *PulseRangeHasherMock {
	m := &PulseRangeHasherMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetRangeHashMock = mPulseRangeHasherMockGetRangeHash{mock: m}
	m.ValidateRangeHashMock = mPulseRangeHasherMockValidateRangeHash{mock: m}

	return m
}

type mPulseRangeHasherMockGetRangeHash struct {
	mock              *PulseRangeHasherMock
	mainExpectation   *PulseRangeHasherMockGetRangeHashExpectation
	expectationSeries []*PulseRangeHasherMockGetRangeHashExpectation
}

type PulseRangeHasherMockGetRangeHashExpectation struct {
	input  *PulseRangeHasherMockGetRangeHashInput
	result *PulseRangeHasherMockGetRangeHashResult
}

type PulseRangeHasherMockGetRangeHashInput struct {
	p insolar.PulseRange
}

type PulseRangeHasherMockGetRangeHashResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of PulseRangeHasher.GetRangeHash is expected from 1 to Infinity times
func (m *mPulseRangeHasherMockGetRangeHash) Expect(p insolar.PulseRange) *mPulseRangeHasherMockGetRangeHash {
	m.mock.GetRangeHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseRangeHasherMockGetRangeHashExpectation{}
	}
	m.mainExpectation.input = &PulseRangeHasherMockGetRangeHashInput{p}
	return m
}

//Return specifies results of invocation of PulseRangeHasher.GetRangeHash
func (m *mPulseRangeHasherMockGetRangeHash) Return(r []byte, r1 error) *PulseRangeHasherMock {
	m.mock.GetRangeHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseRangeHasherMockGetRangeHashExpectation{}
	}
	m.mainExpectation.result = &PulseRangeHasherMockGetRangeHashResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseRangeHasher.GetRangeHash is expected once
func (m *mPulseRangeHasherMockGetRangeHash) ExpectOnce(p insolar.PulseRange) *PulseRangeHasherMockGetRangeHashExpectation {
	m.mock.GetRangeHashFunc = nil
	m.mainExpectation = nil

	expectation := &PulseRangeHasherMockGetRangeHashExpectation{}
	expectation.input = &PulseRangeHasherMockGetRangeHashInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseRangeHasherMockGetRangeHashExpectation) Return(r []byte, r1 error) {
	e.result = &PulseRangeHasherMockGetRangeHashResult{r, r1}
}

//Set uses given function f as a mock of PulseRangeHasher.GetRangeHash method
func (m *mPulseRangeHasherMockGetRangeHash) Set(f func(p insolar.PulseRange) (r []byte, r1 error)) *PulseRangeHasherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRangeHashFunc = f
	return m.mock
}

//GetRangeHash implements github.com/insolar/insolar/network/storage.PulseRangeHasher interface
func (m *PulseRangeHasherMock) GetRangeHash(p insolar.PulseRange) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.GetRangeHashPreCounter, 1)
	defer atomic.AddUint64(&m.GetRangeHashCounter, 1)

	if len(m.GetRangeHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRangeHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseRangeHasherMock.GetRangeHash. %v", p)
			return
		}

		input := m.GetRangeHashMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseRangeHasherMockGetRangeHashInput{p}, "PulseRangeHasher.GetRangeHash got unexpected parameters")

		result := m.GetRangeHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseRangeHasherMock.GetRangeHash")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetRangeHashMock.mainExpectation != nil {

		input := m.GetRangeHashMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseRangeHasherMockGetRangeHashInput{p}, "PulseRangeHasher.GetRangeHash got unexpected parameters")
		}

		result := m.GetRangeHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseRangeHasherMock.GetRangeHash")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetRangeHashFunc == nil {
		m.t.Fatalf("Unexpected call to PulseRangeHasherMock.GetRangeHash. %v", p)
		return
	}

	return m.GetRangeHashFunc(p)
}

//GetRangeHashMinimockCounter returns a count of PulseRangeHasherMock.GetRangeHashFunc invocations
func (m *PulseRangeHasherMock) GetRangeHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRangeHashCounter)
}

//GetRangeHashMinimockPreCounter returns the value of PulseRangeHasherMock.GetRangeHash invocations
func (m *PulseRangeHasherMock) GetRangeHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRangeHashPreCounter)
}

//GetRangeHashFinished returns true if mock invocations count is ok
func (m *PulseRangeHasherMock) GetRangeHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRangeHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRangeHashCounter) == uint64(len(m.GetRangeHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRangeHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRangeHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRangeHashFunc != nil {
		return atomic.LoadUint64(&m.GetRangeHashCounter) > 0
	}

	return true
}

type mPulseRangeHasherMockValidateRangeHash struct {
	mock              *PulseRangeHasherMock
	mainExpectation   *PulseRangeHasherMockValidateRangeHashExpectation
	expectationSeries []*PulseRangeHasherMockValidateRangeHashExpectation
}

type PulseRangeHasherMockValidateRangeHashExpectation struct {
	input  *PulseRangeHasherMockValidateRangeHashInput
	result *PulseRangeHasherMockValidateRangeHashResult
}

type PulseRangeHasherMockValidateRangeHashInput struct {
	p  insolar.PulseRange
	p1 []byte
}

type PulseRangeHasherMockValidateRangeHashResult struct {
	r  bool
	r1 error
}

//Expect specifies that invocation of PulseRangeHasher.ValidateRangeHash is expected from 1 to Infinity times
func (m *mPulseRangeHasherMockValidateRangeHash) Expect(p insolar.PulseRange, p1 []byte) *mPulseRangeHasherMockValidateRangeHash {
	m.mock.ValidateRangeHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseRangeHasherMockValidateRangeHashExpectation{}
	}
	m.mainExpectation.input = &PulseRangeHasherMockValidateRangeHashInput{p, p1}
	return m
}

//Return specifies results of invocation of PulseRangeHasher.ValidateRangeHash
func (m *mPulseRangeHasherMockValidateRangeHash) Return(r bool, r1 error) *PulseRangeHasherMock {
	m.mock.ValidateRangeHashFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulseRangeHasherMockValidateRangeHashExpectation{}
	}
	m.mainExpectation.result = &PulseRangeHasherMockValidateRangeHashResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulseRangeHasher.ValidateRangeHash is expected once
func (m *mPulseRangeHasherMockValidateRangeHash) ExpectOnce(p insolar.PulseRange, p1 []byte) *PulseRangeHasherMockValidateRangeHashExpectation {
	m.mock.ValidateRangeHashFunc = nil
	m.mainExpectation = nil

	expectation := &PulseRangeHasherMockValidateRangeHashExpectation{}
	expectation.input = &PulseRangeHasherMockValidateRangeHashInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulseRangeHasherMockValidateRangeHashExpectation) Return(r bool, r1 error) {
	e.result = &PulseRangeHasherMockValidateRangeHashResult{r, r1}
}

//Set uses given function f as a mock of PulseRangeHasher.ValidateRangeHash method
func (m *mPulseRangeHasherMockValidateRangeHash) Set(f func(p insolar.PulseRange, p1 []byte) (r bool, r1 error)) *PulseRangeHasherMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ValidateRangeHashFunc = f
	return m.mock
}

//ValidateRangeHash implements github.com/insolar/insolar/network/storage.PulseRangeHasher interface
func (m *PulseRangeHasherMock) ValidateRangeHash(p insolar.PulseRange, p1 []byte) (r bool, r1 error) {
	counter := atomic.AddUint64(&m.ValidateRangeHashPreCounter, 1)
	defer atomic.AddUint64(&m.ValidateRangeHashCounter, 1)

	if len(m.ValidateRangeHashMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ValidateRangeHashMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulseRangeHasherMock.ValidateRangeHash. %v %v", p, p1)
			return
		}

		input := m.ValidateRangeHashMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulseRangeHasherMockValidateRangeHashInput{p, p1}, "PulseRangeHasher.ValidateRangeHash got unexpected parameters")

		result := m.ValidateRangeHashMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulseRangeHasherMock.ValidateRangeHash")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateRangeHashMock.mainExpectation != nil {

		input := m.ValidateRangeHashMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulseRangeHasherMockValidateRangeHashInput{p, p1}, "PulseRangeHasher.ValidateRangeHash got unexpected parameters")
		}

		result := m.ValidateRangeHashMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulseRangeHasherMock.ValidateRangeHash")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ValidateRangeHashFunc == nil {
		m.t.Fatalf("Unexpected call to PulseRangeHasherMock.ValidateRangeHash. %v %v", p, p1)
		return
	}

	return m.ValidateRangeHashFunc(p, p1)
}

//ValidateRangeHashMinimockCounter returns a count of PulseRangeHasherMock.ValidateRangeHashFunc invocations
func (m *PulseRangeHasherMock) ValidateRangeHashMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ValidateRangeHashCounter)
}

//ValidateRangeHashMinimockPreCounter returns the value of PulseRangeHasherMock.ValidateRangeHash invocations
func (m *PulseRangeHasherMock) ValidateRangeHashMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ValidateRangeHashPreCounter)
}

//ValidateRangeHashFinished returns true if mock invocations count is ok
func (m *PulseRangeHasherMock) ValidateRangeHashFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ValidateRangeHashMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ValidateRangeHashCounter) == uint64(len(m.ValidateRangeHashMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ValidateRangeHashMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ValidateRangeHashCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ValidateRangeHashFunc != nil {
		return atomic.LoadUint64(&m.ValidateRangeHashCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseRangeHasherMock) ValidateCallCounters() {

	if !m.GetRangeHashFinished() {
		m.t.Fatal("Expected call to PulseRangeHasherMock.GetRangeHash")
	}

	if !m.ValidateRangeHashFinished() {
		m.t.Fatal("Expected call to PulseRangeHasherMock.ValidateRangeHash")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulseRangeHasherMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulseRangeHasherMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulseRangeHasherMock) MinimockFinish() {

	if !m.GetRangeHashFinished() {
		m.t.Fatal("Expected call to PulseRangeHasherMock.GetRangeHash")
	}

	if !m.ValidateRangeHashFinished() {
		m.t.Fatal("Expected call to PulseRangeHasherMock.ValidateRangeHash")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulseRangeHasherMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulseRangeHasherMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetRangeHashFinished()
		ok = ok && m.ValidateRangeHashFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetRangeHashFinished() {
				m.t.Error("Expected call to PulseRangeHasherMock.GetRangeHash")
			}

			if !m.ValidateRangeHashFinished() {
				m.t.Error("Expected call to PulseRangeHasherMock.ValidateRangeHash")
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
func (m *PulseRangeHasherMock) AllMocksCalled() bool {

	if !m.GetRangeHashFinished() {
		return false
	}

	if !m.ValidateRangeHashFinished() {
		return false
	}

	return true
}
