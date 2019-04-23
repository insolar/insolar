package node

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Accessor" can be found in github.com/insolar/insolar/insolar/node
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//AccessorMock implements github.com/insolar/insolar/insolar/node.Accessor
type AccessorMock struct {
	t minimock.Tester

	AllFunc       func(p insolar.PulseNumber) (r []insolar.Node, r1 error)
	AllCounter    uint64
	AllPreCounter uint64
	AllMock       mAccessorMockAll

	InRoleFunc       func(p insolar.PulseNumber, p1 insolar.StaticRole) (r []insolar.Node, r1 error)
	InRoleCounter    uint64
	InRolePreCounter uint64
	InRoleMock       mAccessorMockInRole
}

//NewAccessorMock returns a mock for github.com/insolar/insolar/insolar/node.Accessor
func NewAccessorMock(t minimock.Tester) *AccessorMock {
	m := &AccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AllMock = mAccessorMockAll{mock: m}
	m.InRoleMock = mAccessorMockInRole{mock: m}

	return m
}

type mAccessorMockAll struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockAllExpectation
	expectationSeries []*AccessorMockAllExpectation
}

type AccessorMockAllExpectation struct {
	input  *AccessorMockAllInput
	result *AccessorMockAllResult
}

type AccessorMockAllInput struct {
	p insolar.PulseNumber
}

type AccessorMockAllResult struct {
	r  []insolar.Node
	r1 error
}

//Expect specifies that invocation of Accessor.All is expected from 1 to Infinity times
func (m *mAccessorMockAll) Expect(p insolar.PulseNumber) *mAccessorMockAll {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockAllExpectation{}
	}
	m.mainExpectation.input = &AccessorMockAllInput{p}
	return m
}

//Return specifies results of invocation of Accessor.All
func (m *mAccessorMockAll) Return(r []insolar.Node, r1 error) *AccessorMock {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockAllExpectation{}
	}
	m.mainExpectation.result = &AccessorMockAllResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.All is expected once
func (m *mAccessorMockAll) ExpectOnce(p insolar.PulseNumber) *AccessorMockAllExpectation {
	m.mock.AllFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockAllExpectation{}
	expectation.input = &AccessorMockAllInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockAllExpectation) Return(r []insolar.Node, r1 error) {
	e.result = &AccessorMockAllResult{r, r1}
}

//Set uses given function f as a mock of Accessor.All method
func (m *mAccessorMockAll) Set(f func(p insolar.PulseNumber) (r []insolar.Node, r1 error)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AllFunc = f
	return m.mock
}

//All implements github.com/insolar/insolar/insolar/node.Accessor interface
func (m *AccessorMock) All(p insolar.PulseNumber) (r []insolar.Node, r1 error) {
	counter := atomic.AddUint64(&m.AllPreCounter, 1)
	defer atomic.AddUint64(&m.AllCounter, 1)

	if len(m.AllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.All. %v", p)
			return
		}

		input := m.AllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockAllInput{p}, "Accessor.All got unexpected parameters")

		result := m.AllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.All")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AllMock.mainExpectation != nil {

		input := m.AllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockAllInput{p}, "Accessor.All got unexpected parameters")
		}

		result := m.AllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.All")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AllFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.All. %v", p)
		return
	}

	return m.AllFunc(p)
}

//AllMinimockCounter returns a count of AccessorMock.AllFunc invocations
func (m *AccessorMock) AllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AllCounter)
}

//AllMinimockPreCounter returns the value of AccessorMock.All invocations
func (m *AccessorMock) AllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AllPreCounter)
}

//AllFinished returns true if mock invocations count is ok
func (m *AccessorMock) AllFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AllMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AllCounter) == uint64(len(m.AllMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AllMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AllCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AllFunc != nil {
		return atomic.LoadUint64(&m.AllCounter) > 0
	}

	return true
}

type mAccessorMockInRole struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockInRoleExpectation
	expectationSeries []*AccessorMockInRoleExpectation
}

type AccessorMockInRoleExpectation struct {
	input  *AccessorMockInRoleInput
	result *AccessorMockInRoleResult
}

type AccessorMockInRoleInput struct {
	p  insolar.PulseNumber
	p1 insolar.StaticRole
}

type AccessorMockInRoleResult struct {
	r  []insolar.Node
	r1 error
}

//Expect specifies that invocation of Accessor.InRole is expected from 1 to Infinity times
func (m *mAccessorMockInRole) Expect(p insolar.PulseNumber, p1 insolar.StaticRole) *mAccessorMockInRole {
	m.mock.InRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockInRoleExpectation{}
	}
	m.mainExpectation.input = &AccessorMockInRoleInput{p, p1}
	return m
}

//Return specifies results of invocation of Accessor.InRole
func (m *mAccessorMockInRole) Return(r []insolar.Node, r1 error) *AccessorMock {
	m.mock.InRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockInRoleExpectation{}
	}
	m.mainExpectation.result = &AccessorMockInRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.InRole is expected once
func (m *mAccessorMockInRole) ExpectOnce(p insolar.PulseNumber, p1 insolar.StaticRole) *AccessorMockInRoleExpectation {
	m.mock.InRoleFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockInRoleExpectation{}
	expectation.input = &AccessorMockInRoleInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockInRoleExpectation) Return(r []insolar.Node, r1 error) {
	e.result = &AccessorMockInRoleResult{r, r1}
}

//Set uses given function f as a mock of Accessor.InRole method
func (m *mAccessorMockInRole) Set(f func(p insolar.PulseNumber, p1 insolar.StaticRole) (r []insolar.Node, r1 error)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InRoleFunc = f
	return m.mock
}

//InRole implements github.com/insolar/insolar/insolar/node.Accessor interface
func (m *AccessorMock) InRole(p insolar.PulseNumber, p1 insolar.StaticRole) (r []insolar.Node, r1 error) {
	counter := atomic.AddUint64(&m.InRolePreCounter, 1)
	defer atomic.AddUint64(&m.InRoleCounter, 1)

	if len(m.InRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.InRole. %v %v", p, p1)
			return
		}

		input := m.InRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockInRoleInput{p, p1}, "Accessor.InRole got unexpected parameters")

		result := m.InRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.InRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.InRoleMock.mainExpectation != nil {

		input := m.InRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockInRoleInput{p, p1}, "Accessor.InRole got unexpected parameters")
		}

		result := m.InRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.InRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.InRoleFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.InRole. %v %v", p, p1)
		return
	}

	return m.InRoleFunc(p, p1)
}

//InRoleMinimockCounter returns a count of AccessorMock.InRoleFunc invocations
func (m *AccessorMock) InRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InRoleCounter)
}

//InRoleMinimockPreCounter returns the value of AccessorMock.InRole invocations
func (m *AccessorMock) InRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InRolePreCounter)
}

//InRoleFinished returns true if mock invocations count is ok
func (m *AccessorMock) InRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InRoleCounter) == uint64(len(m.InRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InRoleFunc != nil {
		return atomic.LoadUint64(&m.InRoleCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) ValidateCallCounters() {

	if !m.AllFinished() {
		m.t.Fatal("Expected call to AccessorMock.All")
	}

	if !m.InRoleFinished() {
		m.t.Fatal("Expected call to AccessorMock.InRole")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *AccessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *AccessorMock) MinimockFinish() {

	if !m.AllFinished() {
		m.t.Fatal("Expected call to AccessorMock.All")
	}

	if !m.InRoleFinished() {
		m.t.Fatal("Expected call to AccessorMock.InRole")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *AccessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *AccessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AllFinished()
		ok = ok && m.InRoleFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AllFinished() {
				m.t.Error("Expected call to AccessorMock.All")
			}

			if !m.InRoleFinished() {
				m.t.Error("Expected call to AccessorMock.InRole")
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
func (m *AccessorMock) AllMocksCalled() bool {

	if !m.AllFinished() {
		return false
	}

	if !m.InRoleFinished() {
		return false
	}

	return true
}
