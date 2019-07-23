package integrity

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Validator" can be found in github.com/insolar/insolar/ledger/heavy/replica/integrity
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	sequence "github.com/insolar/insolar/ledger/heavy/sequence"

	testify_assert "github.com/stretchr/testify/assert"
)

//ValidatorMock implements github.com/insolar/insolar/ledger/heavy/replica/integrity.Validator
type ValidatorMock struct {
	t minimock.Tester

	UnwrapAndValidateFunc       func(p []byte) (r []sequence.Item)
	UnwrapAndValidateCounter    uint64
	UnwrapAndValidatePreCounter uint64
	UnwrapAndValidateMock       mValidatorMockUnwrapAndValidate
}

//NewValidatorMock returns a mock for github.com/insolar/insolar/ledger/heavy/replica/integrity.Validator
func NewValidatorMock(t minimock.Tester) *ValidatorMock {
	m := &ValidatorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.UnwrapAndValidateMock = mValidatorMockUnwrapAndValidate{mock: m}

	return m
}

type mValidatorMockUnwrapAndValidate struct {
	mock              *ValidatorMock
	mainExpectation   *ValidatorMockUnwrapAndValidateExpectation
	expectationSeries []*ValidatorMockUnwrapAndValidateExpectation
}

type ValidatorMockUnwrapAndValidateExpectation struct {
	input  *ValidatorMockUnwrapAndValidateInput
	result *ValidatorMockUnwrapAndValidateResult
}

type ValidatorMockUnwrapAndValidateInput struct {
	p []byte
}

type ValidatorMockUnwrapAndValidateResult struct {
	r []sequence.Item
}

//Expect specifies that invocation of Validator.UnwrapAndValidate is expected from 1 to Infinity times
func (m *mValidatorMockUnwrapAndValidate) Expect(p []byte) *mValidatorMockUnwrapAndValidate {
	m.mock.UnwrapAndValidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ValidatorMockUnwrapAndValidateExpectation{}
	}
	m.mainExpectation.input = &ValidatorMockUnwrapAndValidateInput{p}
	return m
}

//Return specifies results of invocation of Validator.UnwrapAndValidate
func (m *mValidatorMockUnwrapAndValidate) Return(r []sequence.Item) *ValidatorMock {
	m.mock.UnwrapAndValidateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ValidatorMockUnwrapAndValidateExpectation{}
	}
	m.mainExpectation.result = &ValidatorMockUnwrapAndValidateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Validator.UnwrapAndValidate is expected once
func (m *mValidatorMockUnwrapAndValidate) ExpectOnce(p []byte) *ValidatorMockUnwrapAndValidateExpectation {
	m.mock.UnwrapAndValidateFunc = nil
	m.mainExpectation = nil

	expectation := &ValidatorMockUnwrapAndValidateExpectation{}
	expectation.input = &ValidatorMockUnwrapAndValidateInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ValidatorMockUnwrapAndValidateExpectation) Return(r []sequence.Item) {
	e.result = &ValidatorMockUnwrapAndValidateResult{r}
}

//Set uses given function f as a mock of Validator.UnwrapAndValidate method
func (m *mValidatorMockUnwrapAndValidate) Set(f func(p []byte) (r []sequence.Item)) *ValidatorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnwrapAndValidateFunc = f
	return m.mock
}

//UnwrapAndValidate implements github.com/insolar/insolar/ledger/heavy/replica/integrity.Validator interface
func (m *ValidatorMock) UnwrapAndValidate(p []byte) (r []sequence.Item) {
	counter := atomic.AddUint64(&m.UnwrapAndValidatePreCounter, 1)
	defer atomic.AddUint64(&m.UnwrapAndValidateCounter, 1)

	if len(m.UnwrapAndValidateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnwrapAndValidateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ValidatorMock.UnwrapAndValidate. %v", p)
			return
		}

		input := m.UnwrapAndValidateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ValidatorMockUnwrapAndValidateInput{p}, "Validator.UnwrapAndValidate got unexpected parameters")

		result := m.UnwrapAndValidateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ValidatorMock.UnwrapAndValidate")
			return
		}

		r = result.r

		return
	}

	if m.UnwrapAndValidateMock.mainExpectation != nil {

		input := m.UnwrapAndValidateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ValidatorMockUnwrapAndValidateInput{p}, "Validator.UnwrapAndValidate got unexpected parameters")
		}

		result := m.UnwrapAndValidateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ValidatorMock.UnwrapAndValidate")
		}

		r = result.r

		return
	}

	if m.UnwrapAndValidateFunc == nil {
		m.t.Fatalf("Unexpected call to ValidatorMock.UnwrapAndValidate. %v", p)
		return
	}

	return m.UnwrapAndValidateFunc(p)
}

//UnwrapAndValidateMinimockCounter returns a count of ValidatorMock.UnwrapAndValidateFunc invocations
func (m *ValidatorMock) UnwrapAndValidateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnwrapAndValidateCounter)
}

//UnwrapAndValidateMinimockPreCounter returns the value of ValidatorMock.UnwrapAndValidate invocations
func (m *ValidatorMock) UnwrapAndValidateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnwrapAndValidatePreCounter)
}

//UnwrapAndValidateFinished returns true if mock invocations count is ok
func (m *ValidatorMock) UnwrapAndValidateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UnwrapAndValidateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UnwrapAndValidateCounter) == uint64(len(m.UnwrapAndValidateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UnwrapAndValidateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UnwrapAndValidateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UnwrapAndValidateFunc != nil {
		return atomic.LoadUint64(&m.UnwrapAndValidateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ValidatorMock) ValidateCallCounters() {

	if !m.UnwrapAndValidateFinished() {
		m.t.Fatal("Expected call to ValidatorMock.UnwrapAndValidate")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ValidatorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ValidatorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ValidatorMock) MinimockFinish() {

	if !m.UnwrapAndValidateFinished() {
		m.t.Fatal("Expected call to ValidatorMock.UnwrapAndValidate")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ValidatorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ValidatorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.UnwrapAndValidateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.UnwrapAndValidateFinished() {
				m.t.Error("Expected call to ValidatorMock.UnwrapAndValidate")
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
func (m *ValidatorMock) AllMocksCalled() bool {

	if !m.UnwrapAndValidateFinished() {
		return false
	}

	return true
}
