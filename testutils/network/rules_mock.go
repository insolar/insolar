package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Rules" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
)

//RulesMock implements github.com/insolar/insolar/network.Rules
type RulesMock struct {
	t minimock.Tester

	CheckMajorityRuleFunc       func() (r bool, r1 int)
	CheckMajorityRuleCounter    uint64
	CheckMajorityRulePreCounter uint64
	CheckMajorityRuleMock       mRulesMockCheckMajorityRule

	CheckMinRoleFunc       func() (r bool)
	CheckMinRoleCounter    uint64
	CheckMinRolePreCounter uint64
	CheckMinRoleMock       mRulesMockCheckMinRole
}

//NewRulesMock returns a mock for github.com/insolar/insolar/network.Rules
func NewRulesMock(t minimock.Tester) *RulesMock {
	m := &RulesMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CheckMajorityRuleMock = mRulesMockCheckMajorityRule{mock: m}
	m.CheckMinRoleMock = mRulesMockCheckMinRole{mock: m}

	return m
}

type mRulesMockCheckMajorityRule struct {
	mock              *RulesMock
	mainExpectation   *RulesMockCheckMajorityRuleExpectation
	expectationSeries []*RulesMockCheckMajorityRuleExpectation
}

type RulesMockCheckMajorityRuleExpectation struct {
	result *RulesMockCheckMajorityRuleResult
}

type RulesMockCheckMajorityRuleResult struct {
	r  bool
	r1 int
}

//Expect specifies that invocation of Rules.CheckMajorityRule is expected from 1 to Infinity times
func (m *mRulesMockCheckMajorityRule) Expect() *mRulesMockCheckMajorityRule {
	m.mock.CheckMajorityRuleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RulesMockCheckMajorityRuleExpectation{}
	}

	return m
}

//Return specifies results of invocation of Rules.CheckMajorityRule
func (m *mRulesMockCheckMajorityRule) Return(r bool, r1 int) *RulesMock {
	m.mock.CheckMajorityRuleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RulesMockCheckMajorityRuleExpectation{}
	}
	m.mainExpectation.result = &RulesMockCheckMajorityRuleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Rules.CheckMajorityRule is expected once
func (m *mRulesMockCheckMajorityRule) ExpectOnce() *RulesMockCheckMajorityRuleExpectation {
	m.mock.CheckMajorityRuleFunc = nil
	m.mainExpectation = nil

	expectation := &RulesMockCheckMajorityRuleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RulesMockCheckMajorityRuleExpectation) Return(r bool, r1 int) {
	e.result = &RulesMockCheckMajorityRuleResult{r, r1}
}

//Set uses given function f as a mock of Rules.CheckMajorityRule method
func (m *mRulesMockCheckMajorityRule) Set(f func() (r bool, r1 int)) *RulesMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CheckMajorityRuleFunc = f
	return m.mock
}

//CheckMajorityRule implements github.com/insolar/insolar/network.Rules interface
func (m *RulesMock) CheckMajorityRule() (r bool, r1 int) {
	counter := atomic.AddUint64(&m.CheckMajorityRulePreCounter, 1)
	defer atomic.AddUint64(&m.CheckMajorityRuleCounter, 1)

	if len(m.CheckMajorityRuleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CheckMajorityRuleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RulesMock.CheckMajorityRule.")
			return
		}

		result := m.CheckMajorityRuleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RulesMock.CheckMajorityRule")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CheckMajorityRuleMock.mainExpectation != nil {

		result := m.CheckMajorityRuleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RulesMock.CheckMajorityRule")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CheckMajorityRuleFunc == nil {
		m.t.Fatalf("Unexpected call to RulesMock.CheckMajorityRule.")
		return
	}

	return m.CheckMajorityRuleFunc()
}

//CheckMajorityRuleMinimockCounter returns a count of RulesMock.CheckMajorityRuleFunc invocations
func (m *RulesMock) CheckMajorityRuleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CheckMajorityRuleCounter)
}

//CheckMajorityRuleMinimockPreCounter returns the value of RulesMock.CheckMajorityRule invocations
func (m *RulesMock) CheckMajorityRuleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CheckMajorityRulePreCounter)
}

//CheckMajorityRuleFinished returns true if mock invocations count is ok
func (m *RulesMock) CheckMajorityRuleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CheckMajorityRuleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CheckMajorityRuleCounter) == uint64(len(m.CheckMajorityRuleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CheckMajorityRuleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CheckMajorityRuleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CheckMajorityRuleFunc != nil {
		return atomic.LoadUint64(&m.CheckMajorityRuleCounter) > 0
	}

	return true
}

type mRulesMockCheckMinRole struct {
	mock              *RulesMock
	mainExpectation   *RulesMockCheckMinRoleExpectation
	expectationSeries []*RulesMockCheckMinRoleExpectation
}

type RulesMockCheckMinRoleExpectation struct {
	result *RulesMockCheckMinRoleResult
}

type RulesMockCheckMinRoleResult struct {
	r bool
}

//Expect specifies that invocation of Rules.CheckMinRole is expected from 1 to Infinity times
func (m *mRulesMockCheckMinRole) Expect() *mRulesMockCheckMinRole {
	m.mock.CheckMinRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RulesMockCheckMinRoleExpectation{}
	}

	return m
}

//Return specifies results of invocation of Rules.CheckMinRole
func (m *mRulesMockCheckMinRole) Return(r bool) *RulesMock {
	m.mock.CheckMinRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &RulesMockCheckMinRoleExpectation{}
	}
	m.mainExpectation.result = &RulesMockCheckMinRoleResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Rules.CheckMinRole is expected once
func (m *mRulesMockCheckMinRole) ExpectOnce() *RulesMockCheckMinRoleExpectation {
	m.mock.CheckMinRoleFunc = nil
	m.mainExpectation = nil

	expectation := &RulesMockCheckMinRoleExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *RulesMockCheckMinRoleExpectation) Return(r bool) {
	e.result = &RulesMockCheckMinRoleResult{r}
}

//Set uses given function f as a mock of Rules.CheckMinRole method
func (m *mRulesMockCheckMinRole) Set(f func() (r bool)) *RulesMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CheckMinRoleFunc = f
	return m.mock
}

//CheckMinRole implements github.com/insolar/insolar/network.Rules interface
func (m *RulesMock) CheckMinRole() (r bool) {
	counter := atomic.AddUint64(&m.CheckMinRolePreCounter, 1)
	defer atomic.AddUint64(&m.CheckMinRoleCounter, 1)

	if len(m.CheckMinRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CheckMinRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to RulesMock.CheckMinRole.")
			return
		}

		result := m.CheckMinRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the RulesMock.CheckMinRole")
			return
		}

		r = result.r

		return
	}

	if m.CheckMinRoleMock.mainExpectation != nil {

		result := m.CheckMinRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the RulesMock.CheckMinRole")
		}

		r = result.r

		return
	}

	if m.CheckMinRoleFunc == nil {
		m.t.Fatalf("Unexpected call to RulesMock.CheckMinRole.")
		return
	}

	return m.CheckMinRoleFunc()
}

//CheckMinRoleMinimockCounter returns a count of RulesMock.CheckMinRoleFunc invocations
func (m *RulesMock) CheckMinRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CheckMinRoleCounter)
}

//CheckMinRoleMinimockPreCounter returns the value of RulesMock.CheckMinRole invocations
func (m *RulesMock) CheckMinRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CheckMinRolePreCounter)
}

//CheckMinRoleFinished returns true if mock invocations count is ok
func (m *RulesMock) CheckMinRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CheckMinRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CheckMinRoleCounter) == uint64(len(m.CheckMinRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CheckMinRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CheckMinRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CheckMinRoleFunc != nil {
		return atomic.LoadUint64(&m.CheckMinRoleCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RulesMock) ValidateCallCounters() {

	if !m.CheckMajorityRuleFinished() {
		m.t.Fatal("Expected call to RulesMock.CheckMajorityRule")
	}

	if !m.CheckMinRoleFinished() {
		m.t.Fatal("Expected call to RulesMock.CheckMinRole")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RulesMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RulesMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RulesMock) MinimockFinish() {

	if !m.CheckMajorityRuleFinished() {
		m.t.Fatal("Expected call to RulesMock.CheckMajorityRule")
	}

	if !m.CheckMinRoleFinished() {
		m.t.Fatal("Expected call to RulesMock.CheckMinRole")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RulesMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RulesMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CheckMajorityRuleFinished()
		ok = ok && m.CheckMinRoleFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CheckMajorityRuleFinished() {
				m.t.Error("Expected call to RulesMock.CheckMajorityRule")
			}

			if !m.CheckMinRoleFinished() {
				m.t.Error("Expected call to RulesMock.CheckMinRole")
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
func (m *RulesMock) AllMocksCalled() bool {

	if !m.CheckMajorityRuleFinished() {
		return false
	}

	if !m.CheckMinRoleFinished() {
		return false
	}

	return true
}
