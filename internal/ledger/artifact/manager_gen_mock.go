package artifact

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Manager" can be found in github.com/insolar/insolar/internal/ledger/artifact
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ManagerMock implements github.com/insolar/insolar/internal/ledger/artifact.Manager
type ManagerMock struct {
	t minimock.Tester

	RegisterRequestFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error)
	RegisterRequestCounter    uint64
	RegisterRequestPreCounter uint64
	RegisterRequestMock       mManagerMockRegisterRequest
}

//NewManagerMock returns a mock for github.com/insolar/insolar/internal/ledger/artifact.Manager
func NewManagerMock(t minimock.Tester) *ManagerMock {
	m := &ManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RegisterRequestMock = mManagerMockRegisterRequest{mock: m}

	return m
}

type mManagerMockRegisterRequest struct {
	mock              *ManagerMock
	mainExpectation   *ManagerMockRegisterRequestExpectation
	expectationSeries []*ManagerMockRegisterRequestExpectation
}

type ManagerMockRegisterRequestExpectation struct {
	input  *ManagerMockRegisterRequestInput
	result *ManagerMockRegisterRequestResult
}

type ManagerMockRegisterRequestInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Parcel
}

type ManagerMockRegisterRequestResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Manager.RegisterRequest is expected from 1 to Infinity times
func (m *mManagerMockRegisterRequest) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) *mManagerMockRegisterRequest {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockRegisterRequestExpectation{}
	}
	m.mainExpectation.input = &ManagerMockRegisterRequestInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Manager.RegisterRequest
func (m *mManagerMockRegisterRequest) Return(r *insolar.ID, r1 error) *ManagerMock {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockRegisterRequestExpectation{}
	}
	m.mainExpectation.result = &ManagerMockRegisterRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Manager.RegisterRequest is expected once
func (m *mManagerMockRegisterRequest) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) *ManagerMockRegisterRequestExpectation {
	m.mock.RegisterRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ManagerMockRegisterRequestExpectation{}
	expectation.input = &ManagerMockRegisterRequestInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ManagerMockRegisterRequestExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ManagerMockRegisterRequestResult{r, r1}
}

//Set uses given function f as a mock of Manager.RegisterRequest method
func (m *mManagerMockRegisterRequest) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error)) *ManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterRequestFunc = f
	return m.mock
}

//RegisterRequest implements github.com/insolar/insolar/internal/ledger/artifact.Manager interface
func (m *ManagerMock) RegisterRequest(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestCounter, 1)

	if len(m.RegisterRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ManagerMock.RegisterRequest. %v %v %v", p, p1, p2)
			return
		}

		input := m.RegisterRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ManagerMockRegisterRequestInput{p, p1, p2}, "Manager.RegisterRequest got unexpected parameters")

		result := m.RegisterRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.RegisterRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestMock.mainExpectation != nil {

		input := m.RegisterRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ManagerMockRegisterRequestInput{p, p1, p2}, "Manager.RegisterRequest got unexpected parameters")
		}

		result := m.RegisterRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.RegisterRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ManagerMock.RegisterRequest. %v %v %v", p, p1, p2)
		return
	}

	return m.RegisterRequestFunc(p, p1, p2)
}

//RegisterRequestMinimockCounter returns a count of ManagerMock.RegisterRequestFunc invocations
func (m *ManagerMock) RegisterRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestCounter)
}

//RegisterRequestMinimockPreCounter returns the value of ManagerMock.RegisterRequest invocations
func (m *ManagerMock) RegisterRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestPreCounter)
}

//RegisterRequestFinished returns true if mock invocations count is ok
func (m *ManagerMock) RegisterRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterRequestCounter) == uint64(len(m.RegisterRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterRequestFunc != nil {
		return atomic.LoadUint64(&m.RegisterRequestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ManagerMock) ValidateCallCounters() {

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ManagerMock.RegisterRequest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ManagerMock) MinimockFinish() {

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ManagerMock.RegisterRequest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.RegisterRequestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.RegisterRequestFinished() {
				m.t.Error("Expected call to ManagerMock.RegisterRequest")
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
func (m *ManagerMock) AllMocksCalled() bool {

	if !m.RegisterRequestFinished() {
		return false
	}

	return true
}
