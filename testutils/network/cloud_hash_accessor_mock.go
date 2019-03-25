package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CloudHashAccessor" can be found in github.com/insolar/insolar/network/storage
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//CloudHashAccessorMock implements github.com/insolar/insolar/network/storage.CloudHashAccessor
type CloudHashAccessorMock struct {
	t minimock.Tester

	ForPulseNumberFunc       func(p context.Context, p1 insolar.PulseNumber) (r []byte, r1 error)
	ForPulseNumberCounter    uint64
	ForPulseNumberPreCounter uint64
	ForPulseNumberMock       mCloudHashAccessorMockForPulseNumber

	LatestFunc       func(p context.Context) (r []byte, r1 error)
	LatestCounter    uint64
	LatestPreCounter uint64
	LatestMock       mCloudHashAccessorMockLatest
}

//NewCloudHashAccessorMock returns a mock for github.com/insolar/insolar/network/storage.CloudHashAccessor
func NewCloudHashAccessorMock(t minimock.Tester) *CloudHashAccessorMock {
	m := &CloudHashAccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ForPulseNumberMock = mCloudHashAccessorMockForPulseNumber{mock: m}
	m.LatestMock = mCloudHashAccessorMockLatest{mock: m}

	return m
}

type mCloudHashAccessorMockForPulseNumber struct {
	mock              *CloudHashAccessorMock
	mainExpectation   *CloudHashAccessorMockForPulseNumberExpectation
	expectationSeries []*CloudHashAccessorMockForPulseNumberExpectation
}

type CloudHashAccessorMockForPulseNumberExpectation struct {
	input  *CloudHashAccessorMockForPulseNumberInput
	result *CloudHashAccessorMockForPulseNumberResult
}

type CloudHashAccessorMockForPulseNumberInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

type CloudHashAccessorMockForPulseNumberResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of CloudHashAccessor.ForPulseNumber is expected from 1 to Infinity times
func (m *mCloudHashAccessorMockForPulseNumber) Expect(p context.Context, p1 insolar.PulseNumber) *mCloudHashAccessorMockForPulseNumber {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudHashAccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.input = &CloudHashAccessorMockForPulseNumberInput{p, p1}
	return m
}

//Return specifies results of invocation of CloudHashAccessor.ForPulseNumber
func (m *mCloudHashAccessorMockForPulseNumber) Return(r []byte, r1 error) *CloudHashAccessorMock {
	m.mock.ForPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudHashAccessorMockForPulseNumberExpectation{}
	}
	m.mainExpectation.result = &CloudHashAccessorMockForPulseNumberResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudHashAccessor.ForPulseNumber is expected once
func (m *mCloudHashAccessorMockForPulseNumber) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *CloudHashAccessorMockForPulseNumberExpectation {
	m.mock.ForPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &CloudHashAccessorMockForPulseNumberExpectation{}
	expectation.input = &CloudHashAccessorMockForPulseNumberInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudHashAccessorMockForPulseNumberExpectation) Return(r []byte, r1 error) {
	e.result = &CloudHashAccessorMockForPulseNumberResult{r, r1}
}

//Set uses given function f as a mock of CloudHashAccessor.ForPulseNumber method
func (m *mCloudHashAccessorMockForPulseNumber) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []byte, r1 error)) *CloudHashAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForPulseNumberFunc = f
	return m.mock
}

//ForPulseNumber implements github.com/insolar/insolar/network/storage.CloudHashAccessor interface
func (m *CloudHashAccessorMock) ForPulseNumber(p context.Context, p1 insolar.PulseNumber) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.ForPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.ForPulseNumberCounter, 1)

	if len(m.ForPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudHashAccessorMock.ForPulseNumber. %v %v", p, p1)
			return
		}

		input := m.ForPulseNumberMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CloudHashAccessorMockForPulseNumberInput{p, p1}, "CloudHashAccessor.ForPulseNumber got unexpected parameters")

		result := m.ForPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudHashAccessorMock.ForPulseNumber")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberMock.mainExpectation != nil {

		input := m.ForPulseNumberMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CloudHashAccessorMockForPulseNumberInput{p, p1}, "CloudHashAccessor.ForPulseNumber got unexpected parameters")
		}

		result := m.ForPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudHashAccessorMock.ForPulseNumber")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to CloudHashAccessorMock.ForPulseNumber. %v %v", p, p1)
		return
	}

	return m.ForPulseNumberFunc(p, p1)
}

//ForPulseNumberMinimockCounter returns a count of CloudHashAccessorMock.ForPulseNumberFunc invocations
func (m *CloudHashAccessorMock) ForPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberCounter)
}

//ForPulseNumberMinimockPreCounter returns the value of CloudHashAccessorMock.ForPulseNumber invocations
func (m *CloudHashAccessorMock) ForPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForPulseNumberPreCounter)
}

//ForPulseNumberFinished returns true if mock invocations count is ok
func (m *CloudHashAccessorMock) ForPulseNumberFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForPulseNumberMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) == uint64(len(m.ForPulseNumberMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForPulseNumberMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForPulseNumberFunc != nil {
		return atomic.LoadUint64(&m.ForPulseNumberCounter) > 0
	}

	return true
}

type mCloudHashAccessorMockLatest struct {
	mock              *CloudHashAccessorMock
	mainExpectation   *CloudHashAccessorMockLatestExpectation
	expectationSeries []*CloudHashAccessorMockLatestExpectation
}

type CloudHashAccessorMockLatestExpectation struct {
	input  *CloudHashAccessorMockLatestInput
	result *CloudHashAccessorMockLatestResult
}

type CloudHashAccessorMockLatestInput struct {
	p context.Context
}

type CloudHashAccessorMockLatestResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of CloudHashAccessor.Latest is expected from 1 to Infinity times
func (m *mCloudHashAccessorMockLatest) Expect(p context.Context) *mCloudHashAccessorMockLatest {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudHashAccessorMockLatestExpectation{}
	}
	m.mainExpectation.input = &CloudHashAccessorMockLatestInput{p}
	return m
}

//Return specifies results of invocation of CloudHashAccessor.Latest
func (m *mCloudHashAccessorMockLatest) Return(r []byte, r1 error) *CloudHashAccessorMock {
	m.mock.LatestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CloudHashAccessorMockLatestExpectation{}
	}
	m.mainExpectation.result = &CloudHashAccessorMockLatestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CloudHashAccessor.Latest is expected once
func (m *mCloudHashAccessorMockLatest) ExpectOnce(p context.Context) *CloudHashAccessorMockLatestExpectation {
	m.mock.LatestFunc = nil
	m.mainExpectation = nil

	expectation := &CloudHashAccessorMockLatestExpectation{}
	expectation.input = &CloudHashAccessorMockLatestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CloudHashAccessorMockLatestExpectation) Return(r []byte, r1 error) {
	e.result = &CloudHashAccessorMockLatestResult{r, r1}
}

//Set uses given function f as a mock of CloudHashAccessor.Latest method
func (m *mCloudHashAccessorMockLatest) Set(f func(p context.Context) (r []byte, r1 error)) *CloudHashAccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LatestFunc = f
	return m.mock
}

//Latest implements github.com/insolar/insolar/network/storage.CloudHashAccessor interface
func (m *CloudHashAccessorMock) Latest(p context.Context) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.LatestPreCounter, 1)
	defer atomic.AddUint64(&m.LatestCounter, 1)

	if len(m.LatestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LatestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CloudHashAccessorMock.Latest. %v", p)
			return
		}

		input := m.LatestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, CloudHashAccessorMockLatestInput{p}, "CloudHashAccessor.Latest got unexpected parameters")

		result := m.LatestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CloudHashAccessorMock.Latest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestMock.mainExpectation != nil {

		input := m.LatestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, CloudHashAccessorMockLatestInput{p}, "CloudHashAccessor.Latest got unexpected parameters")
		}

		result := m.LatestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CloudHashAccessorMock.Latest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.LatestFunc == nil {
		m.t.Fatalf("Unexpected call to CloudHashAccessorMock.Latest. %v", p)
		return
	}

	return m.LatestFunc(p)
}

//LatestMinimockCounter returns a count of CloudHashAccessorMock.LatestFunc invocations
func (m *CloudHashAccessorMock) LatestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LatestCounter)
}

//LatestMinimockPreCounter returns the value of CloudHashAccessorMock.Latest invocations
func (m *CloudHashAccessorMock) LatestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LatestPreCounter)
}

//LatestFinished returns true if mock invocations count is ok
func (m *CloudHashAccessorMock) LatestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LatestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LatestCounter) == uint64(len(m.LatestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LatestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LatestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LatestFunc != nil {
		return atomic.LoadUint64(&m.LatestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CloudHashAccessorMock) ValidateCallCounters() {

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to CloudHashAccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to CloudHashAccessorMock.Latest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CloudHashAccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CloudHashAccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CloudHashAccessorMock) MinimockFinish() {

	if !m.ForPulseNumberFinished() {
		m.t.Fatal("Expected call to CloudHashAccessorMock.ForPulseNumber")
	}

	if !m.LatestFinished() {
		m.t.Fatal("Expected call to CloudHashAccessorMock.Latest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CloudHashAccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CloudHashAccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ForPulseNumberFinished()
		ok = ok && m.LatestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ForPulseNumberFinished() {
				m.t.Error("Expected call to CloudHashAccessorMock.ForPulseNumber")
			}

			if !m.LatestFinished() {
				m.t.Error("Expected call to CloudHashAccessorMock.Latest")
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
func (m *CloudHashAccessorMock) AllMocksCalled() bool {

	if !m.ForPulseNumberFinished() {
		return false
	}

	if !m.LatestFinished() {
		return false
	}

	return true
}
