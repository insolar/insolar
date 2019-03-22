package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CodeDescriptor" can be found in github.com/insolar/insolar/insolar
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
)

//CodeDescriptorMock implements github.com/insolar/insolar/insolar.CodeDescriptor
type CodeDescriptorMock struct {
	t minimock.Tester

	CodeFunc       func() (r []byte, r1 error)
	CodeCounter    uint64
	CodePreCounter uint64
	CodeMock       mCodeDescriptorMockCode

	MachineTypeFunc       func() (r insolar.MachineType)
	MachineTypeCounter    uint64
	MachineTypePreCounter uint64
	MachineTypeMock       mCodeDescriptorMockMachineType

	RefFunc       func() (r *insolar.RecordRef)
	RefCounter    uint64
	RefPreCounter uint64
	RefMock       mCodeDescriptorMockRef
}

//NewCodeDescriptorMock returns a mock for github.com/insolar/insolar/insolar.CodeDescriptor
func NewCodeDescriptorMock(t minimock.Tester) *CodeDescriptorMock {
	m := &CodeDescriptorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CodeMock = mCodeDescriptorMockCode{mock: m}
	m.MachineTypeMock = mCodeDescriptorMockMachineType{mock: m}
	m.RefMock = mCodeDescriptorMockRef{mock: m}

	return m
}

type mCodeDescriptorMockCode struct {
	mock              *CodeDescriptorMock
	mainExpectation   *CodeDescriptorMockCodeExpectation
	expectationSeries []*CodeDescriptorMockCodeExpectation
}

type CodeDescriptorMockCodeExpectation struct {
	result *CodeDescriptorMockCodeResult
}

type CodeDescriptorMockCodeResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of CodeDescriptor.Code is expected from 1 to Infinity times
func (m *mCodeDescriptorMockCode) Expect() *mCodeDescriptorMockCode {
	m.mock.CodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CodeDescriptorMockCodeExpectation{}
	}

	return m
}

//Return specifies results of invocation of CodeDescriptor.Code
func (m *mCodeDescriptorMockCode) Return(r []byte, r1 error) *CodeDescriptorMock {
	m.mock.CodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CodeDescriptorMockCodeExpectation{}
	}
	m.mainExpectation.result = &CodeDescriptorMockCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of CodeDescriptor.Code is expected once
func (m *mCodeDescriptorMockCode) ExpectOnce() *CodeDescriptorMockCodeExpectation {
	m.mock.CodeFunc = nil
	m.mainExpectation = nil

	expectation := &CodeDescriptorMockCodeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CodeDescriptorMockCodeExpectation) Return(r []byte, r1 error) {
	e.result = &CodeDescriptorMockCodeResult{r, r1}
}

//Set uses given function f as a mock of CodeDescriptor.Code method
func (m *mCodeDescriptorMockCode) Set(f func() (r []byte, r1 error)) *CodeDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CodeFunc = f
	return m.mock
}

//Code implements github.com/insolar/insolar/insolar.CodeDescriptor interface
func (m *CodeDescriptorMock) Code() (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.CodePreCounter, 1)
	defer atomic.AddUint64(&m.CodeCounter, 1)

	if len(m.CodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CodeDescriptorMock.Code.")
			return
		}

		result := m.CodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CodeDescriptorMock.Code")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CodeMock.mainExpectation != nil {

		result := m.CodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CodeDescriptorMock.Code")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CodeFunc == nil {
		m.t.Fatalf("Unexpected call to CodeDescriptorMock.Code.")
		return
	}

	return m.CodeFunc()
}

//CodeMinimockCounter returns a count of CodeDescriptorMock.CodeFunc invocations
func (m *CodeDescriptorMock) CodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CodeCounter)
}

//CodeMinimockPreCounter returns the value of CodeDescriptorMock.Code invocations
func (m *CodeDescriptorMock) CodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CodePreCounter)
}

//CodeFinished returns true if mock invocations count is ok
func (m *CodeDescriptorMock) CodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CodeCounter) == uint64(len(m.CodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CodeFunc != nil {
		return atomic.LoadUint64(&m.CodeCounter) > 0
	}

	return true
}

type mCodeDescriptorMockMachineType struct {
	mock              *CodeDescriptorMock
	mainExpectation   *CodeDescriptorMockMachineTypeExpectation
	expectationSeries []*CodeDescriptorMockMachineTypeExpectation
}

type CodeDescriptorMockMachineTypeExpectation struct {
	result *CodeDescriptorMockMachineTypeResult
}

type CodeDescriptorMockMachineTypeResult struct {
	r insolar.MachineType
}

//Expect specifies that invocation of CodeDescriptor.MachineType is expected from 1 to Infinity times
func (m *mCodeDescriptorMockMachineType) Expect() *mCodeDescriptorMockMachineType {
	m.mock.MachineTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CodeDescriptorMockMachineTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of CodeDescriptor.MachineType
func (m *mCodeDescriptorMockMachineType) Return(r insolar.MachineType) *CodeDescriptorMock {
	m.mock.MachineTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CodeDescriptorMockMachineTypeExpectation{}
	}
	m.mainExpectation.result = &CodeDescriptorMockMachineTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CodeDescriptor.MachineType is expected once
func (m *mCodeDescriptorMockMachineType) ExpectOnce() *CodeDescriptorMockMachineTypeExpectation {
	m.mock.MachineTypeFunc = nil
	m.mainExpectation = nil

	expectation := &CodeDescriptorMockMachineTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CodeDescriptorMockMachineTypeExpectation) Return(r insolar.MachineType) {
	e.result = &CodeDescriptorMockMachineTypeResult{r}
}

//Set uses given function f as a mock of CodeDescriptor.MachineType method
func (m *mCodeDescriptorMockMachineType) Set(f func() (r insolar.MachineType)) *CodeDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MachineTypeFunc = f
	return m.mock
}

//MachineType implements github.com/insolar/insolar/insolar.CodeDescriptor interface
func (m *CodeDescriptorMock) MachineType() (r insolar.MachineType) {
	counter := atomic.AddUint64(&m.MachineTypePreCounter, 1)
	defer atomic.AddUint64(&m.MachineTypeCounter, 1)

	if len(m.MachineTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MachineTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CodeDescriptorMock.MachineType.")
			return
		}

		result := m.MachineTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CodeDescriptorMock.MachineType")
			return
		}

		r = result.r

		return
	}

	if m.MachineTypeMock.mainExpectation != nil {

		result := m.MachineTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CodeDescriptorMock.MachineType")
		}

		r = result.r

		return
	}

	if m.MachineTypeFunc == nil {
		m.t.Fatalf("Unexpected call to CodeDescriptorMock.MachineType.")
		return
	}

	return m.MachineTypeFunc()
}

//MachineTypeMinimockCounter returns a count of CodeDescriptorMock.MachineTypeFunc invocations
func (m *CodeDescriptorMock) MachineTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MachineTypeCounter)
}

//MachineTypeMinimockPreCounter returns the value of CodeDescriptorMock.MachineType invocations
func (m *CodeDescriptorMock) MachineTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MachineTypePreCounter)
}

//MachineTypeFinished returns true if mock invocations count is ok
func (m *CodeDescriptorMock) MachineTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MachineTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MachineTypeCounter) == uint64(len(m.MachineTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MachineTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MachineTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MachineTypeFunc != nil {
		return atomic.LoadUint64(&m.MachineTypeCounter) > 0
	}

	return true
}

type mCodeDescriptorMockRef struct {
	mock              *CodeDescriptorMock
	mainExpectation   *CodeDescriptorMockRefExpectation
	expectationSeries []*CodeDescriptorMockRefExpectation
}

type CodeDescriptorMockRefExpectation struct {
	result *CodeDescriptorMockRefResult
}

type CodeDescriptorMockRefResult struct {
	r *insolar.RecordRef
}

//Expect specifies that invocation of CodeDescriptor.Ref is expected from 1 to Infinity times
func (m *mCodeDescriptorMockRef) Expect() *mCodeDescriptorMockRef {
	m.mock.RefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CodeDescriptorMockRefExpectation{}
	}

	return m
}

//Return specifies results of invocation of CodeDescriptor.Ref
func (m *mCodeDescriptorMockRef) Return(r *insolar.RecordRef) *CodeDescriptorMock {
	m.mock.RefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &CodeDescriptorMockRefExpectation{}
	}
	m.mainExpectation.result = &CodeDescriptorMockRefResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of CodeDescriptor.Ref is expected once
func (m *mCodeDescriptorMockRef) ExpectOnce() *CodeDescriptorMockRefExpectation {
	m.mock.RefFunc = nil
	m.mainExpectation = nil

	expectation := &CodeDescriptorMockRefExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *CodeDescriptorMockRefExpectation) Return(r *insolar.RecordRef) {
	e.result = &CodeDescriptorMockRefResult{r}
}

//Set uses given function f as a mock of CodeDescriptor.Ref method
func (m *mCodeDescriptorMockRef) Set(f func() (r *insolar.RecordRef)) *CodeDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RefFunc = f
	return m.mock
}

//Ref implements github.com/insolar/insolar/insolar.CodeDescriptor interface
func (m *CodeDescriptorMock) Ref() (r *insolar.RecordRef) {
	counter := atomic.AddUint64(&m.RefPreCounter, 1)
	defer atomic.AddUint64(&m.RefCounter, 1)

	if len(m.RefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to CodeDescriptorMock.Ref.")
			return
		}

		result := m.RefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the CodeDescriptorMock.Ref")
			return
		}

		r = result.r

		return
	}

	if m.RefMock.mainExpectation != nil {

		result := m.RefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the CodeDescriptorMock.Ref")
		}

		r = result.r

		return
	}

	if m.RefFunc == nil {
		m.t.Fatalf("Unexpected call to CodeDescriptorMock.Ref.")
		return
	}

	return m.RefFunc()
}

//RefMinimockCounter returns a count of CodeDescriptorMock.RefFunc invocations
func (m *CodeDescriptorMock) RefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RefCounter)
}

//RefMinimockPreCounter returns the value of CodeDescriptorMock.Ref invocations
func (m *CodeDescriptorMock) RefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RefPreCounter)
}

//RefFinished returns true if mock invocations count is ok
func (m *CodeDescriptorMock) RefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RefCounter) == uint64(len(m.RefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RefFunc != nil {
		return atomic.LoadUint64(&m.RefCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CodeDescriptorMock) ValidateCallCounters() {

	if !m.CodeFinished() {
		m.t.Fatal("Expected call to CodeDescriptorMock.Code")
	}

	if !m.MachineTypeFinished() {
		m.t.Fatal("Expected call to CodeDescriptorMock.MachineType")
	}

	if !m.RefFinished() {
		m.t.Fatal("Expected call to CodeDescriptorMock.Ref")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CodeDescriptorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CodeDescriptorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CodeDescriptorMock) MinimockFinish() {

	if !m.CodeFinished() {
		m.t.Fatal("Expected call to CodeDescriptorMock.Code")
	}

	if !m.MachineTypeFinished() {
		m.t.Fatal("Expected call to CodeDescriptorMock.MachineType")
	}

	if !m.RefFinished() {
		m.t.Fatal("Expected call to CodeDescriptorMock.Ref")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CodeDescriptorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CodeDescriptorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CodeFinished()
		ok = ok && m.MachineTypeFinished()
		ok = ok && m.RefFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CodeFinished() {
				m.t.Error("Expected call to CodeDescriptorMock.Code")
			}

			if !m.MachineTypeFinished() {
				m.t.Error("Expected call to CodeDescriptorMock.MachineType")
			}

			if !m.RefFinished() {
				m.t.Error("Expected call to CodeDescriptorMock.Ref")
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
func (m *CodeDescriptorMock) AllMocksCalled() bool {

	if !m.CodeFinished() {
		return false
	}

	if !m.MachineTypeFinished() {
		return false
	}

	if !m.RefFinished() {
		return false
	}

	return true
}
