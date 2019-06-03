package pulsartestutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PulsarStorage" can be found in github.com/insolar/insolar/pulsar/storage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//PulsarStorageMock implements github.com/insolar/insolar/pulsar/storage.PulsarStorage
type PulsarStorageMock struct {
	t minimock.Tester

	CloseFunc       func() (r error)
	CloseCounter    uint64
	ClosePreCounter uint64
	CloseMock       mPulsarStorageMockClose

	GetLastPulseFunc       func() (r *insolar.Pulse, r1 error)
	GetLastPulseCounter    uint64
	GetLastPulsePreCounter uint64
	GetLastPulseMock       mPulsarStorageMockGetLastPulse

	SavePulseFunc       func(p *insolar.Pulse) (r error)
	SavePulseCounter    uint64
	SavePulsePreCounter uint64
	SavePulseMock       mPulsarStorageMockSavePulse

	SetLastPulseFunc       func(p *insolar.Pulse) (r error)
	SetLastPulseCounter    uint64
	SetLastPulsePreCounter uint64
	SetLastPulseMock       mPulsarStorageMockSetLastPulse
}

//NewPulsarStorageMock returns a mock for github.com/insolar/insolar/pulsar/storage.PulsarStorage
func NewPulsarStorageMock(t minimock.Tester) *PulsarStorageMock {
	m := &PulsarStorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseMock = mPulsarStorageMockClose{mock: m}
	m.GetLastPulseMock = mPulsarStorageMockGetLastPulse{mock: m}
	m.SavePulseMock = mPulsarStorageMockSavePulse{mock: m}
	m.SetLastPulseMock = mPulsarStorageMockSetLastPulse{mock: m}

	return m
}

type mPulsarStorageMockClose struct {
	mock              *PulsarStorageMock
	mainExpectation   *PulsarStorageMockCloseExpectation
	expectationSeries []*PulsarStorageMockCloseExpectation
}

type PulsarStorageMockCloseExpectation struct {
	result *PulsarStorageMockCloseResult
}

type PulsarStorageMockCloseResult struct {
	r error
}

//Expect specifies that invocation of PulsarStorage.Close is expected from 1 to Infinity times
func (m *mPulsarStorageMockClose) Expect() *mPulsarStorageMockClose {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockCloseExpectation{}
	}

	return m
}

//Return specifies results of invocation of PulsarStorage.Close
func (m *mPulsarStorageMockClose) Return(r error) *PulsarStorageMock {
	m.mock.CloseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockCloseExpectation{}
	}
	m.mainExpectation.result = &PulsarStorageMockCloseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulsarStorage.Close is expected once
func (m *mPulsarStorageMockClose) ExpectOnce() *PulsarStorageMockCloseExpectation {
	m.mock.CloseFunc = nil
	m.mainExpectation = nil

	expectation := &PulsarStorageMockCloseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulsarStorageMockCloseExpectation) Return(r error) {
	e.result = &PulsarStorageMockCloseResult{r}
}

//Set uses given function f as a mock of PulsarStorage.Close method
func (m *mPulsarStorageMockClose) Set(f func() (r error)) *PulsarStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloseFunc = f
	return m.mock
}

//Close implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) Close() (r error) {
	counter := atomic.AddUint64(&m.ClosePreCounter, 1)
	defer atomic.AddUint64(&m.CloseCounter, 1)

	if len(m.CloseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulsarStorageMock.Close.")
			return
		}

		result := m.CloseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.Close")
			return
		}

		r = result.r

		return
	}

	if m.CloseMock.mainExpectation != nil {

		result := m.CloseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.Close")
		}

		r = result.r

		return
	}

	if m.CloseFunc == nil {
		m.t.Fatalf("Unexpected call to PulsarStorageMock.Close.")
		return
	}

	return m.CloseFunc()
}

//CloseMinimockCounter returns a count of PulsarStorageMock.CloseFunc invocations
func (m *PulsarStorageMock) CloseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseCounter)
}

//CloseMinimockPreCounter returns the value of PulsarStorageMock.Close invocations
func (m *PulsarStorageMock) CloseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClosePreCounter)
}

//CloseFinished returns true if mock invocations count is ok
func (m *PulsarStorageMock) CloseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloseCounter) == uint64(len(m.CloseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloseFunc != nil {
		return atomic.LoadUint64(&m.CloseCounter) > 0
	}

	return true
}

type mPulsarStorageMockGetLastPulse struct {
	mock              *PulsarStorageMock
	mainExpectation   *PulsarStorageMockGetLastPulseExpectation
	expectationSeries []*PulsarStorageMockGetLastPulseExpectation
}

type PulsarStorageMockGetLastPulseExpectation struct {
	result *PulsarStorageMockGetLastPulseResult
}

type PulsarStorageMockGetLastPulseResult struct {
	r  *insolar.Pulse
	r1 error
}

//Expect specifies that invocation of PulsarStorage.GetLastPulse is expected from 1 to Infinity times
func (m *mPulsarStorageMockGetLastPulse) Expect() *mPulsarStorageMockGetLastPulse {
	m.mock.GetLastPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockGetLastPulseExpectation{}
	}

	return m
}

//Return specifies results of invocation of PulsarStorage.GetLastPulse
func (m *mPulsarStorageMockGetLastPulse) Return(r *insolar.Pulse, r1 error) *PulsarStorageMock {
	m.mock.GetLastPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockGetLastPulseExpectation{}
	}
	m.mainExpectation.result = &PulsarStorageMockGetLastPulseResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PulsarStorage.GetLastPulse is expected once
func (m *mPulsarStorageMockGetLastPulse) ExpectOnce() *PulsarStorageMockGetLastPulseExpectation {
	m.mock.GetLastPulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulsarStorageMockGetLastPulseExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulsarStorageMockGetLastPulseExpectation) Return(r *insolar.Pulse, r1 error) {
	e.result = &PulsarStorageMockGetLastPulseResult{r, r1}
}

//Set uses given function f as a mock of PulsarStorage.GetLastPulse method
func (m *mPulsarStorageMockGetLastPulse) Set(f func() (r *insolar.Pulse, r1 error)) *PulsarStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetLastPulseFunc = f
	return m.mock
}

//GetLastPulse implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) GetLastPulse() (r *insolar.Pulse, r1 error) {
	counter := atomic.AddUint64(&m.GetLastPulsePreCounter, 1)
	defer atomic.AddUint64(&m.GetLastPulseCounter, 1)

	if len(m.GetLastPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetLastPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulsarStorageMock.GetLastPulse.")
			return
		}

		result := m.GetLastPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.GetLastPulse")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetLastPulseMock.mainExpectation != nil {

		result := m.GetLastPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.GetLastPulse")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetLastPulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulsarStorageMock.GetLastPulse.")
		return
	}

	return m.GetLastPulseFunc()
}

//GetLastPulseMinimockCounter returns a count of PulsarStorageMock.GetLastPulseFunc invocations
func (m *PulsarStorageMock) GetLastPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetLastPulseCounter)
}

//GetLastPulseMinimockPreCounter returns the value of PulsarStorageMock.GetLastPulse invocations
func (m *PulsarStorageMock) GetLastPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetLastPulsePreCounter)
}

//GetLastPulseFinished returns true if mock invocations count is ok
func (m *PulsarStorageMock) GetLastPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetLastPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetLastPulseCounter) == uint64(len(m.GetLastPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetLastPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetLastPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetLastPulseFunc != nil {
		return atomic.LoadUint64(&m.GetLastPulseCounter) > 0
	}

	return true
}

type mPulsarStorageMockSavePulse struct {
	mock              *PulsarStorageMock
	mainExpectation   *PulsarStorageMockSavePulseExpectation
	expectationSeries []*PulsarStorageMockSavePulseExpectation
}

type PulsarStorageMockSavePulseExpectation struct {
	input  *PulsarStorageMockSavePulseInput
	result *PulsarStorageMockSavePulseResult
}

type PulsarStorageMockSavePulseInput struct {
	p *insolar.Pulse
}

type PulsarStorageMockSavePulseResult struct {
	r error
}

//Expect specifies that invocation of PulsarStorage.SavePulse is expected from 1 to Infinity times
func (m *mPulsarStorageMockSavePulse) Expect(p *insolar.Pulse) *mPulsarStorageMockSavePulse {
	m.mock.SavePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockSavePulseExpectation{}
	}
	m.mainExpectation.input = &PulsarStorageMockSavePulseInput{p}
	return m
}

//Return specifies results of invocation of PulsarStorage.SavePulse
func (m *mPulsarStorageMockSavePulse) Return(r error) *PulsarStorageMock {
	m.mock.SavePulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockSavePulseExpectation{}
	}
	m.mainExpectation.result = &PulsarStorageMockSavePulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulsarStorage.SavePulse is expected once
func (m *mPulsarStorageMockSavePulse) ExpectOnce(p *insolar.Pulse) *PulsarStorageMockSavePulseExpectation {
	m.mock.SavePulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulsarStorageMockSavePulseExpectation{}
	expectation.input = &PulsarStorageMockSavePulseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulsarStorageMockSavePulseExpectation) Return(r error) {
	e.result = &PulsarStorageMockSavePulseResult{r}
}

//Set uses given function f as a mock of PulsarStorage.SavePulse method
func (m *mPulsarStorageMockSavePulse) Set(f func(p *insolar.Pulse) (r error)) *PulsarStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SavePulseFunc = f
	return m.mock
}

//SavePulse implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) SavePulse(p *insolar.Pulse) (r error) {
	counter := atomic.AddUint64(&m.SavePulsePreCounter, 1)
	defer atomic.AddUint64(&m.SavePulseCounter, 1)

	if len(m.SavePulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SavePulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulsarStorageMock.SavePulse. %v", p)
			return
		}

		input := m.SavePulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulsarStorageMockSavePulseInput{p}, "PulsarStorage.SavePulse got unexpected parameters")

		result := m.SavePulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.SavePulse")
			return
		}

		r = result.r

		return
	}

	if m.SavePulseMock.mainExpectation != nil {

		input := m.SavePulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulsarStorageMockSavePulseInput{p}, "PulsarStorage.SavePulse got unexpected parameters")
		}

		result := m.SavePulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.SavePulse")
		}

		r = result.r

		return
	}

	if m.SavePulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulsarStorageMock.SavePulse. %v", p)
		return
	}

	return m.SavePulseFunc(p)
}

//SavePulseMinimockCounter returns a count of PulsarStorageMock.SavePulseFunc invocations
func (m *PulsarStorageMock) SavePulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SavePulseCounter)
}

//SavePulseMinimockPreCounter returns the value of PulsarStorageMock.SavePulse invocations
func (m *PulsarStorageMock) SavePulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SavePulsePreCounter)
}

//SavePulseFinished returns true if mock invocations count is ok
func (m *PulsarStorageMock) SavePulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SavePulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SavePulseCounter) == uint64(len(m.SavePulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SavePulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SavePulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SavePulseFunc != nil {
		return atomic.LoadUint64(&m.SavePulseCounter) > 0
	}

	return true
}

type mPulsarStorageMockSetLastPulse struct {
	mock              *PulsarStorageMock
	mainExpectation   *PulsarStorageMockSetLastPulseExpectation
	expectationSeries []*PulsarStorageMockSetLastPulseExpectation
}

type PulsarStorageMockSetLastPulseExpectation struct {
	input  *PulsarStorageMockSetLastPulseInput
	result *PulsarStorageMockSetLastPulseResult
}

type PulsarStorageMockSetLastPulseInput struct {
	p *insolar.Pulse
}

type PulsarStorageMockSetLastPulseResult struct {
	r error
}

//Expect specifies that invocation of PulsarStorage.SetLastPulse is expected from 1 to Infinity times
func (m *mPulsarStorageMockSetLastPulse) Expect(p *insolar.Pulse) *mPulsarStorageMockSetLastPulse {
	m.mock.SetLastPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockSetLastPulseExpectation{}
	}
	m.mainExpectation.input = &PulsarStorageMockSetLastPulseInput{p}
	return m
}

//Return specifies results of invocation of PulsarStorage.SetLastPulse
func (m *mPulsarStorageMockSetLastPulse) Return(r error) *PulsarStorageMock {
	m.mock.SetLastPulseFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PulsarStorageMockSetLastPulseExpectation{}
	}
	m.mainExpectation.result = &PulsarStorageMockSetLastPulseResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PulsarStorage.SetLastPulse is expected once
func (m *mPulsarStorageMockSetLastPulse) ExpectOnce(p *insolar.Pulse) *PulsarStorageMockSetLastPulseExpectation {
	m.mock.SetLastPulseFunc = nil
	m.mainExpectation = nil

	expectation := &PulsarStorageMockSetLastPulseExpectation{}
	expectation.input = &PulsarStorageMockSetLastPulseInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PulsarStorageMockSetLastPulseExpectation) Return(r error) {
	e.result = &PulsarStorageMockSetLastPulseResult{r}
}

//Set uses given function f as a mock of PulsarStorage.SetLastPulse method
func (m *mPulsarStorageMockSetLastPulse) Set(f func(p *insolar.Pulse) (r error)) *PulsarStorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetLastPulseFunc = f
	return m.mock
}

//SetLastPulse implements github.com/insolar/insolar/pulsar/storage.PulsarStorage interface
func (m *PulsarStorageMock) SetLastPulse(p *insolar.Pulse) (r error) {
	counter := atomic.AddUint64(&m.SetLastPulsePreCounter, 1)
	defer atomic.AddUint64(&m.SetLastPulseCounter, 1)

	if len(m.SetLastPulseMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetLastPulseMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PulsarStorageMock.SetLastPulse. %v", p)
			return
		}

		input := m.SetLastPulseMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PulsarStorageMockSetLastPulseInput{p}, "PulsarStorage.SetLastPulse got unexpected parameters")

		result := m.SetLastPulseMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.SetLastPulse")
			return
		}

		r = result.r

		return
	}

	if m.SetLastPulseMock.mainExpectation != nil {

		input := m.SetLastPulseMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PulsarStorageMockSetLastPulseInput{p}, "PulsarStorage.SetLastPulse got unexpected parameters")
		}

		result := m.SetLastPulseMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PulsarStorageMock.SetLastPulse")
		}

		r = result.r

		return
	}

	if m.SetLastPulseFunc == nil {
		m.t.Fatalf("Unexpected call to PulsarStorageMock.SetLastPulse. %v", p)
		return
	}

	return m.SetLastPulseFunc(p)
}

//SetLastPulseMinimockCounter returns a count of PulsarStorageMock.SetLastPulseFunc invocations
func (m *PulsarStorageMock) SetLastPulseMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetLastPulseCounter)
}

//SetLastPulseMinimockPreCounter returns the value of PulsarStorageMock.SetLastPulse invocations
func (m *PulsarStorageMock) SetLastPulseMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetLastPulsePreCounter)
}

//SetLastPulseFinished returns true if mock invocations count is ok
func (m *PulsarStorageMock) SetLastPulseFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetLastPulseMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetLastPulseCounter) == uint64(len(m.SetLastPulseMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetLastPulseMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetLastPulseCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetLastPulseFunc != nil {
		return atomic.LoadUint64(&m.SetLastPulseCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulsarStorageMock) ValidateCallCounters() {

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.Close")
	}

	if !m.GetLastPulseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.GetLastPulse")
	}

	if !m.SavePulseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.SavePulse")
	}

	if !m.SetLastPulseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.SetLastPulse")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PulsarStorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PulsarStorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PulsarStorageMock) MinimockFinish() {

	if !m.CloseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.Close")
	}

	if !m.GetLastPulseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.GetLastPulse")
	}

	if !m.SavePulseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.SavePulse")
	}

	if !m.SetLastPulseFinished() {
		m.t.Fatal("Expected call to PulsarStorageMock.SetLastPulse")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PulsarStorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PulsarStorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloseFinished()
		ok = ok && m.GetLastPulseFinished()
		ok = ok && m.SavePulseFinished()
		ok = ok && m.SetLastPulseFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloseFinished() {
				m.t.Error("Expected call to PulsarStorageMock.Close")
			}

			if !m.GetLastPulseFinished() {
				m.t.Error("Expected call to PulsarStorageMock.GetLastPulse")
			}

			if !m.SavePulseFinished() {
				m.t.Error("Expected call to PulsarStorageMock.SavePulse")
			}

			if !m.SetLastPulseFinished() {
				m.t.Error("Expected call to PulsarStorageMock.SetLastPulse")
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
func (m *PulsarStorageMock) AllMocksCalled() bool {

	if !m.CloseFinished() {
		return false
	}

	if !m.GetLastPulseFinished() {
		return false
	}

	if !m.SavePulseFinished() {
		return false
	}

	if !m.SetLastPulseFinished() {
		return false
	}

	return true
}
