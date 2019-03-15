package db

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "JetIndexModifier" can be found in github.com/insolar/insolar/ledger/storage/db
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//JetIndexModifierMock implements github.com/insolar/insolar/ledger/storage/db.JetIndexModifier
type JetIndexModifierMock struct {
	t minimock.Tester

	AddFunc       func(p core.RecordID, p1 core.JetID)
	AddCounter    uint64
	AddPreCounter uint64
	AddMock       mJetIndexModifierMockAdd

	DeleteFunc       func(p core.RecordID, p1 core.JetID)
	DeleteCounter    uint64
	DeletePreCounter uint64
	DeleteMock       mJetIndexModifierMockDelete
}

//NewJetIndexModifierMock returns a mock for github.com/insolar/insolar/ledger/storage/db.JetIndexModifier
func NewJetIndexModifierMock(t minimock.Tester) *JetIndexModifierMock {
	m := &JetIndexModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddMock = mJetIndexModifierMockAdd{mock: m}
	m.DeleteMock = mJetIndexModifierMockDelete{mock: m}

	return m
}

type mJetIndexModifierMockAdd struct {
	mock              *JetIndexModifierMock
	mainExpectation   *JetIndexModifierMockAddExpectation
	expectationSeries []*JetIndexModifierMockAddExpectation
}

type JetIndexModifierMockAddExpectation struct {
	input *JetIndexModifierMockAddInput
}

type JetIndexModifierMockAddInput struct {
	p  core.RecordID
	p1 core.JetID
}

//Expect specifies that invocation of JetIndexModifier.Add is expected from 1 to Infinity times
func (m *mJetIndexModifierMockAdd) Expect(p core.RecordID, p1 core.JetID) *mJetIndexModifierMockAdd {
	m.mock.AddFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetIndexModifierMockAddExpectation{}
	}
	m.mainExpectation.input = &JetIndexModifierMockAddInput{p, p1}
	return m
}

//Return specifies results of invocation of JetIndexModifier.Add
func (m *mJetIndexModifierMockAdd) Return() *JetIndexModifierMock {
	m.mock.AddFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetIndexModifierMockAddExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of JetIndexModifier.Add is expected once
func (m *mJetIndexModifierMockAdd) ExpectOnce(p core.RecordID, p1 core.JetID) *JetIndexModifierMockAddExpectation {
	m.mock.AddFunc = nil
	m.mainExpectation = nil

	expectation := &JetIndexModifierMockAddExpectation{}
	expectation.input = &JetIndexModifierMockAddInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of JetIndexModifier.Add method
func (m *mJetIndexModifierMockAdd) Set(f func(p core.RecordID, p1 core.JetID)) *JetIndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AddFunc = f
	return m.mock
}

//Add implements github.com/insolar/insolar/ledger/storage/db.JetIndexModifier interface
func (m *JetIndexModifierMock) Add(p core.RecordID, p1 core.JetID) {
	counter := atomic.AddUint64(&m.AddPreCounter, 1)
	defer atomic.AddUint64(&m.AddCounter, 1)

	if len(m.AddMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AddMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetIndexModifierMock.Add. %v %v", p, p1)
			return
		}

		input := m.AddMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetIndexModifierMockAddInput{p, p1}, "JetIndexModifier.Add got unexpected parameters")

		return
	}

	if m.AddMock.mainExpectation != nil {

		input := m.AddMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetIndexModifierMockAddInput{p, p1}, "JetIndexModifier.Add got unexpected parameters")
		}

		return
	}

	if m.AddFunc == nil {
		m.t.Fatalf("Unexpected call to JetIndexModifierMock.Add. %v %v", p, p1)
		return
	}

	m.AddFunc(p, p1)
}

//AddMinimockCounter returns a count of JetIndexModifierMock.AddFunc invocations
func (m *JetIndexModifierMock) AddMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AddCounter)
}

//AddMinimockPreCounter returns the value of JetIndexModifierMock.Add invocations
func (m *JetIndexModifierMock) AddMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AddPreCounter)
}

//AddFinished returns true if mock invocations count is ok
func (m *JetIndexModifierMock) AddFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AddMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AddCounter) == uint64(len(m.AddMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AddMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AddCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AddFunc != nil {
		return atomic.LoadUint64(&m.AddCounter) > 0
	}

	return true
}

type mJetIndexModifierMockDelete struct {
	mock              *JetIndexModifierMock
	mainExpectation   *JetIndexModifierMockDeleteExpectation
	expectationSeries []*JetIndexModifierMockDeleteExpectation
}

type JetIndexModifierMockDeleteExpectation struct {
	input *JetIndexModifierMockDeleteInput
}

type JetIndexModifierMockDeleteInput struct {
	p  core.RecordID
	p1 core.JetID
}

//Expect specifies that invocation of JetIndexModifier.Delete is expected from 1 to Infinity times
func (m *mJetIndexModifierMockDelete) Expect(p core.RecordID, p1 core.JetID) *mJetIndexModifierMockDelete {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetIndexModifierMockDeleteExpectation{}
	}
	m.mainExpectation.input = &JetIndexModifierMockDeleteInput{p, p1}
	return m
}

//Return specifies results of invocation of JetIndexModifier.Delete
func (m *mJetIndexModifierMockDelete) Return() *JetIndexModifierMock {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &JetIndexModifierMockDeleteExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of JetIndexModifier.Delete is expected once
func (m *mJetIndexModifierMockDelete) ExpectOnce(p core.RecordID, p1 core.JetID) *JetIndexModifierMockDeleteExpectation {
	m.mock.DeleteFunc = nil
	m.mainExpectation = nil

	expectation := &JetIndexModifierMockDeleteExpectation{}
	expectation.input = &JetIndexModifierMockDeleteInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of JetIndexModifier.Delete method
func (m *mJetIndexModifierMockDelete) Set(f func(p core.RecordID, p1 core.JetID)) *JetIndexModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteFunc = f
	return m.mock
}

//Delete implements github.com/insolar/insolar/ledger/storage/db.JetIndexModifier interface
func (m *JetIndexModifierMock) Delete(p core.RecordID, p1 core.JetID) {
	counter := atomic.AddUint64(&m.DeletePreCounter, 1)
	defer atomic.AddUint64(&m.DeleteCounter, 1)

	if len(m.DeleteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to JetIndexModifierMock.Delete. %v %v", p, p1)
			return
		}

		input := m.DeleteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, JetIndexModifierMockDeleteInput{p, p1}, "JetIndexModifier.Delete got unexpected parameters")

		return
	}

	if m.DeleteMock.mainExpectation != nil {

		input := m.DeleteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, JetIndexModifierMockDeleteInput{p, p1}, "JetIndexModifier.Delete got unexpected parameters")
		}

		return
	}

	if m.DeleteFunc == nil {
		m.t.Fatalf("Unexpected call to JetIndexModifierMock.Delete. %v %v", p, p1)
		return
	}

	m.DeleteFunc(p, p1)
}

//DeleteMinimockCounter returns a count of JetIndexModifierMock.DeleteFunc invocations
func (m *JetIndexModifierMock) DeleteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteCounter)
}

//DeleteMinimockPreCounter returns the value of JetIndexModifierMock.Delete invocations
func (m *JetIndexModifierMock) DeleteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeletePreCounter)
}

//DeleteFinished returns true if mock invocations count is ok
func (m *JetIndexModifierMock) DeleteFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeleteMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeleteCounter) == uint64(len(m.DeleteMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeleteMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeleteCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeleteFunc != nil {
		return atomic.LoadUint64(&m.DeleteCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetIndexModifierMock) ValidateCallCounters() {

	if !m.AddFinished() {
		m.t.Fatal("Expected call to JetIndexModifierMock.Add")
	}

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to JetIndexModifierMock.Delete")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *JetIndexModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *JetIndexModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *JetIndexModifierMock) MinimockFinish() {

	if !m.AddFinished() {
		m.t.Fatal("Expected call to JetIndexModifierMock.Add")
	}

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to JetIndexModifierMock.Delete")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *JetIndexModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *JetIndexModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AddFinished()
		ok = ok && m.DeleteFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AddFinished() {
				m.t.Error("Expected call to JetIndexModifierMock.Add")
			}

			if !m.DeleteFinished() {
				m.t.Error("Expected call to JetIndexModifierMock.Delete")
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
func (m *JetIndexModifierMock) AllMocksCalled() bool {

	if !m.AddFinished() {
		return false
	}

	if !m.DeleteFinished() {
		return false
	}

	return true
}
