package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Modifier" can be found in github.com/insolar/insolar/insolar/jet
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ModifierMock implements github.com/insolar/insolar/insolar/jet.Modifier
type ModifierMock struct {
	t minimock.Tester

	CloneFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber)
	CloneCounter    uint64
	ClonePreCounter uint64
	CloneMock       mModifierMockClone

	DeleteFunc       func(p context.Context, p1 insolar.PulseNumber)
	DeleteCounter    uint64
	DeletePreCounter uint64
	DeleteMock       mModifierMockDelete

	SplitFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r insolar.JetID, r1 insolar.JetID, r2 error)
	SplitCounter    uint64
	SplitPreCounter uint64
	SplitMock       mModifierMockSplit

	UpdateFunc       func(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID)
	UpdateCounter    uint64
	UpdatePreCounter uint64
	UpdateMock       mModifierMockUpdate
}

//NewModifierMock returns a mock for github.com/insolar/insolar/insolar/jet.Modifier
func NewModifierMock(t minimock.Tester) *ModifierMock {
	m := &ModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloneMock = mModifierMockClone{mock: m}
	m.DeleteMock = mModifierMockDelete{mock: m}
	m.SplitMock = mModifierMockSplit{mock: m}
	m.UpdateMock = mModifierMockUpdate{mock: m}

	return m
}

type mModifierMockClone struct {
	mock              *ModifierMock
	mainExpectation   *ModifierMockCloneExpectation
	expectationSeries []*ModifierMockCloneExpectation
}

type ModifierMockCloneExpectation struct {
	input *ModifierMockCloneInput
}

type ModifierMockCloneInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.PulseNumber
}

//Expect specifies that invocation of Modifier.Clone is expected from 1 to Infinity times
func (m *mModifierMockClone) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) *mModifierMockClone {
	m.mock.CloneFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockCloneExpectation{}
	}
	m.mainExpectation.input = &ModifierMockCloneInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Modifier.Clone
func (m *mModifierMockClone) Return() *ModifierMock {
	m.mock.CloneFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockCloneExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Modifier.Clone is expected once
func (m *mModifierMockClone) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) *ModifierMockCloneExpectation {
	m.mock.CloneFunc = nil
	m.mainExpectation = nil

	expectation := &ModifierMockCloneExpectation{}
	expectation.input = &ModifierMockCloneInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Modifier.Clone method
func (m *mModifierMockClone) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber)) *ModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloneFunc = f
	return m.mock
}

//Clone implements github.com/insolar/insolar/insolar/jet.Modifier interface
func (m *ModifierMock) Clone(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.ClonePreCounter, 1)
	defer atomic.AddUint64(&m.CloneCounter, 1)

	if len(m.CloneMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloneMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ModifierMock.Clone. %v %v %v", p, p1, p2)
			return
		}

		input := m.CloneMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ModifierMockCloneInput{p, p1, p2}, "Modifier.Clone got unexpected parameters")

		return
	}

	if m.CloneMock.mainExpectation != nil {

		input := m.CloneMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ModifierMockCloneInput{p, p1, p2}, "Modifier.Clone got unexpected parameters")
		}

		return
	}

	if m.CloneFunc == nil {
		m.t.Fatalf("Unexpected call to ModifierMock.Clone. %v %v %v", p, p1, p2)
		return
	}

	m.CloneFunc(p, p1, p2)
}

//CloneMinimockCounter returns a count of ModifierMock.CloneFunc invocations
func (m *ModifierMock) CloneMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloneCounter)
}

//CloneMinimockPreCounter returns the value of ModifierMock.Clone invocations
func (m *ModifierMock) CloneMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClonePreCounter)
}

//CloneFinished returns true if mock invocations count is ok
func (m *ModifierMock) CloneFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloneMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloneCounter) == uint64(len(m.CloneMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloneMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloneCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloneFunc != nil {
		return atomic.LoadUint64(&m.CloneCounter) > 0
	}

	return true
}

type mModifierMockDelete struct {
	mock              *ModifierMock
	mainExpectation   *ModifierMockDeleteExpectation
	expectationSeries []*ModifierMockDeleteExpectation
}

type ModifierMockDeleteExpectation struct {
	input *ModifierMockDeleteInput
}

type ModifierMockDeleteInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of Modifier.Delete is expected from 1 to Infinity times
func (m *mModifierMockDelete) Expect(p context.Context, p1 insolar.PulseNumber) *mModifierMockDelete {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockDeleteExpectation{}
	}
	m.mainExpectation.input = &ModifierMockDeleteInput{p, p1}
	return m
}

//Return specifies results of invocation of Modifier.Delete
func (m *mModifierMockDelete) Return() *ModifierMock {
	m.mock.DeleteFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockDeleteExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Modifier.Delete is expected once
func (m *mModifierMockDelete) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *ModifierMockDeleteExpectation {
	m.mock.DeleteFunc = nil
	m.mainExpectation = nil

	expectation := &ModifierMockDeleteExpectation{}
	expectation.input = &ModifierMockDeleteInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Modifier.Delete method
func (m *mModifierMockDelete) Set(f func(p context.Context, p1 insolar.PulseNumber)) *ModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteFunc = f
	return m.mock
}

//Delete implements github.com/insolar/insolar/insolar/jet.Modifier interface
func (m *ModifierMock) DeleteForPN(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.DeletePreCounter, 1)
	defer atomic.AddUint64(&m.DeleteCounter, 1)

	if len(m.DeleteMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ModifierMock.Delete. %v %v", p, p1)
			return
		}

		input := m.DeleteMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ModifierMockDeleteInput{p, p1}, "Modifier.Delete got unexpected parameters")

		return
	}

	if m.DeleteMock.mainExpectation != nil {

		input := m.DeleteMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ModifierMockDeleteInput{p, p1}, "Modifier.Delete got unexpected parameters")
		}

		return
	}

	if m.DeleteFunc == nil {
		m.t.Fatalf("Unexpected call to ModifierMock.Delete. %v %v", p, p1)
		return
	}

	m.DeleteFunc(p, p1)
}

//DeleteMinimockCounter returns a count of ModifierMock.DeleteFunc invocations
func (m *ModifierMock) DeleteMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteCounter)
}

//DeleteMinimockPreCounter returns the value of ModifierMock.Delete invocations
func (m *ModifierMock) DeleteMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeletePreCounter)
}

//DeleteFinished returns true if mock invocations count is ok
func (m *ModifierMock) DeleteFinished() bool {
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

type mModifierMockSplit struct {
	mock              *ModifierMock
	mainExpectation   *ModifierMockSplitExpectation
	expectationSeries []*ModifierMockSplitExpectation
}

type ModifierMockSplitExpectation struct {
	input  *ModifierMockSplitInput
	result *ModifierMockSplitResult
}

type ModifierMockSplitInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type ModifierMockSplitResult struct {
	r  insolar.JetID
	r1 insolar.JetID
	r2 error
}

//Expect specifies that invocation of Modifier.Split is expected from 1 to Infinity times
func (m *mModifierMockSplit) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mModifierMockSplit {
	m.mock.SplitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockSplitExpectation{}
	}
	m.mainExpectation.input = &ModifierMockSplitInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Modifier.Split
func (m *mModifierMockSplit) Return(r insolar.JetID, r1 insolar.JetID, r2 error) *ModifierMock {
	m.mock.SplitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockSplitExpectation{}
	}
	m.mainExpectation.result = &ModifierMockSplitResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of Modifier.Split is expected once
func (m *mModifierMockSplit) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *ModifierMockSplitExpectation {
	m.mock.SplitFunc = nil
	m.mainExpectation = nil

	expectation := &ModifierMockSplitExpectation{}
	expectation.input = &ModifierMockSplitInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ModifierMockSplitExpectation) Return(r insolar.JetID, r1 insolar.JetID, r2 error) {
	e.result = &ModifierMockSplitResult{r, r1, r2}
}

//Set uses given function f as a mock of Modifier.Split method
func (m *mModifierMockSplit) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r insolar.JetID, r1 insolar.JetID, r2 error)) *ModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SplitFunc = f
	return m.mock
}

//Split implements github.com/insolar/insolar/insolar/jet.Modifier interface
func (m *ModifierMock) Split(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r insolar.JetID, r1 insolar.JetID, r2 error) {
	counter := atomic.AddUint64(&m.SplitPreCounter, 1)
	defer atomic.AddUint64(&m.SplitCounter, 1)

	if len(m.SplitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SplitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ModifierMock.Split. %v %v %v", p, p1, p2)
			return
		}

		input := m.SplitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ModifierMockSplitInput{p, p1, p2}, "Modifier.Split got unexpected parameters")

		result := m.SplitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ModifierMock.Split")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.SplitMock.mainExpectation != nil {

		input := m.SplitMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ModifierMockSplitInput{p, p1, p2}, "Modifier.Split got unexpected parameters")
		}

		result := m.SplitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ModifierMock.Split")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.SplitFunc == nil {
		m.t.Fatalf("Unexpected call to ModifierMock.Split. %v %v %v", p, p1, p2)
		return
	}

	return m.SplitFunc(p, p1, p2)
}

//SplitMinimockCounter returns a count of ModifierMock.SplitFunc invocations
func (m *ModifierMock) SplitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SplitCounter)
}

//SplitMinimockPreCounter returns the value of ModifierMock.Split invocations
func (m *ModifierMock) SplitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SplitPreCounter)
}

//SplitFinished returns true if mock invocations count is ok
func (m *ModifierMock) SplitFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SplitMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SplitCounter) == uint64(len(m.SplitMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SplitMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SplitCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SplitFunc != nil {
		return atomic.LoadUint64(&m.SplitCounter) > 0
	}

	return true
}

type mModifierMockUpdate struct {
	mock              *ModifierMock
	mainExpectation   *ModifierMockUpdateExpectation
	expectationSeries []*ModifierMockUpdateExpectation
}

type ModifierMockUpdateExpectation struct {
	input *ModifierMockUpdateInput
}

type ModifierMockUpdateInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 bool
	p3 []insolar.JetID
}

//Expect specifies that invocation of Modifier.Update is expected from 1 to Infinity times
func (m *mModifierMockUpdate) Expect(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID) *mModifierMockUpdate {
	m.mock.UpdateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockUpdateExpectation{}
	}
	m.mainExpectation.input = &ModifierMockUpdateInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Modifier.Update
func (m *mModifierMockUpdate) Return() *ModifierMock {
	m.mock.UpdateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockUpdateExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Modifier.Update is expected once
func (m *mModifierMockUpdate) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID) *ModifierMockUpdateExpectation {
	m.mock.UpdateFunc = nil
	m.mainExpectation = nil

	expectation := &ModifierMockUpdateExpectation{}
	expectation.input = &ModifierMockUpdateInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Modifier.Update method
func (m *mModifierMockUpdate) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID)) *ModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateFunc = f
	return m.mock
}

//Update implements github.com/insolar/insolar/insolar/jet.Modifier interface
func (m *ModifierMock) Update(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID) {
	counter := atomic.AddUint64(&m.UpdatePreCounter, 1)
	defer atomic.AddUint64(&m.UpdateCounter, 1)

	if len(m.UpdateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ModifierMock.Update. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.UpdateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ModifierMockUpdateInput{p, p1, p2, p3}, "Modifier.Update got unexpected parameters")

		return
	}

	if m.UpdateMock.mainExpectation != nil {

		input := m.UpdateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ModifierMockUpdateInput{p, p1, p2, p3}, "Modifier.Update got unexpected parameters")
		}

		return
	}

	if m.UpdateFunc == nil {
		m.t.Fatalf("Unexpected call to ModifierMock.Update. %v %v %v %v", p, p1, p2, p3)
		return
	}

	m.UpdateFunc(p, p1, p2, p3...)
}

//UpdateMinimockCounter returns a count of ModifierMock.UpdateFunc invocations
func (m *ModifierMock) UpdateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateCounter)
}

//UpdateMinimockPreCounter returns the value of ModifierMock.Update invocations
func (m *ModifierMock) UpdateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdatePreCounter)
}

//UpdateFinished returns true if mock invocations count is ok
func (m *ModifierMock) UpdateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpdateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpdateCounter) == uint64(len(m.UpdateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpdateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpdateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpdateFunc != nil {
		return atomic.LoadUint64(&m.UpdateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ModifierMock) ValidateCallCounters() {

	if !m.CloneFinished() {
		m.t.Fatal("Expected call to ModifierMock.Clone")
	}

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to ModifierMock.Delete")
	}

	if !m.SplitFinished() {
		m.t.Fatal("Expected call to ModifierMock.Split")
	}

	if !m.UpdateFinished() {
		m.t.Fatal("Expected call to ModifierMock.Update")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ModifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ModifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ModifierMock) MinimockFinish() {

	if !m.CloneFinished() {
		m.t.Fatal("Expected call to ModifierMock.Clone")
	}

	if !m.DeleteFinished() {
		m.t.Fatal("Expected call to ModifierMock.Delete")
	}

	if !m.SplitFinished() {
		m.t.Fatal("Expected call to ModifierMock.Split")
	}

	if !m.UpdateFinished() {
		m.t.Fatal("Expected call to ModifierMock.Update")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ModifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ModifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloneFinished()
		ok = ok && m.DeleteFinished()
		ok = ok && m.SplitFinished()
		ok = ok && m.UpdateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloneFinished() {
				m.t.Error("Expected call to ModifierMock.Clone")
			}

			if !m.DeleteFinished() {
				m.t.Error("Expected call to ModifierMock.Delete")
			}

			if !m.SplitFinished() {
				m.t.Error("Expected call to ModifierMock.Split")
			}

			if !m.UpdateFinished() {
				m.t.Error("Expected call to ModifierMock.Update")
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
func (m *ModifierMock) AllMocksCalled() bool {

	if !m.CloneFinished() {
		return false
	}

	if !m.DeleteFinished() {
		return false
	}

	if !m.SplitFinished() {
		return false
	}

	if !m.UpdateFinished() {
		return false
	}

	return true
}
