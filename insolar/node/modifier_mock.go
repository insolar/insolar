package node

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Modifier" can be found in github.com/insolar/insolar/insolar/node
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ModifierMock implements github.com/insolar/insolar/insolar/node.Modifier
type ModifierMock struct {
	t minimock.Tester

	DeleteForPNFunc       func(p insolar.PulseNumber)
	DeleteForPNCounter    uint64
	DeleteForPNPreCounter uint64
	DeleteForPNMock       mModifierMockDeleteForPN

	SetFunc       func(p insolar.PulseNumber, p1 []insolar.Node) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mModifierMockSet
}

//NewModifierMock returns a mock for github.com/insolar/insolar/insolar/node.Modifier
func NewModifierMock(t minimock.Tester) *ModifierMock {
	m := &ModifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteForPNMock = mModifierMockDeleteForPN{mock: m}
	m.SetMock = mModifierMockSet{mock: m}

	return m
}

type mModifierMockDeleteForPN struct {
	mock              *ModifierMock
	mainExpectation   *ModifierMockDeleteForPNExpectation
	expectationSeries []*ModifierMockDeleteForPNExpectation
}

type ModifierMockDeleteForPNExpectation struct {
	input *ModifierMockDeleteForPNInput
}

type ModifierMockDeleteForPNInput struct {
	p insolar.PulseNumber
}

//Expect specifies that invocation of Modifier.DeleteForPN is expected from 1 to Infinity times
func (m *mModifierMockDeleteForPN) Expect(p insolar.PulseNumber) *mModifierMockDeleteForPN {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockDeleteForPNExpectation{}
	}
	m.mainExpectation.input = &ModifierMockDeleteForPNInput{p}
	return m
}

//Return specifies results of invocation of Modifier.DeleteForPN
func (m *mModifierMockDeleteForPN) Return() *ModifierMock {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockDeleteForPNExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Modifier.DeleteForPN is expected once
func (m *mModifierMockDeleteForPN) ExpectOnce(p insolar.PulseNumber) *ModifierMockDeleteForPNExpectation {
	m.mock.DeleteForPNFunc = nil
	m.mainExpectation = nil

	expectation := &ModifierMockDeleteForPNExpectation{}
	expectation.input = &ModifierMockDeleteForPNInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Modifier.DeleteForPN method
func (m *mModifierMockDeleteForPN) Set(f func(p insolar.PulseNumber)) *ModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteForPNFunc = f
	return m.mock
}

//DeleteForPN implements github.com/insolar/insolar/insolar/node.Modifier interface
func (m *ModifierMock) DeleteForPN(p insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.DeleteForPNPreCounter, 1)
	defer atomic.AddUint64(&m.DeleteForPNCounter, 1)

	if len(m.DeleteForPNMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteForPNMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ModifierMock.DeleteForPN. %v", p)
			return
		}

		input := m.DeleteForPNMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ModifierMockDeleteForPNInput{p}, "Modifier.DeleteForPN got unexpected parameters")

		return
	}

	if m.DeleteForPNMock.mainExpectation != nil {

		input := m.DeleteForPNMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ModifierMockDeleteForPNInput{p}, "Modifier.DeleteForPN got unexpected parameters")
		}

		return
	}

	if m.DeleteForPNFunc == nil {
		m.t.Fatalf("Unexpected call to ModifierMock.DeleteForPN. %v", p)
		return
	}

	m.DeleteForPNFunc(p)
}

//DeleteForPNMinimockCounter returns a count of ModifierMock.DeleteForPNFunc invocations
func (m *ModifierMock) DeleteForPNMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNCounter)
}

//DeleteForPNMinimockPreCounter returns the value of ModifierMock.DeleteForPN invocations
func (m *ModifierMock) DeleteForPNMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNPreCounter)
}

//DeleteForPNFinished returns true if mock invocations count is ok
func (m *ModifierMock) DeleteForPNFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeleteForPNMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeleteForPNCounter) == uint64(len(m.DeleteForPNMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeleteForPNMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeleteForPNCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeleteForPNFunc != nil {
		return atomic.LoadUint64(&m.DeleteForPNCounter) > 0
	}

	return true
}

type mModifierMockSet struct {
	mock              *ModifierMock
	mainExpectation   *ModifierMockSetExpectation
	expectationSeries []*ModifierMockSetExpectation
}

type ModifierMockSetExpectation struct {
	input  *ModifierMockSetInput
	result *ModifierMockSetResult
}

type ModifierMockSetInput struct {
	p  insolar.PulseNumber
	p1 []insolar.Node
}

type ModifierMockSetResult struct {
	r error
}

//Expect specifies that invocation of Modifier.Set is expected from 1 to Infinity times
func (m *mModifierMockSet) Expect(p insolar.PulseNumber, p1 []insolar.Node) *mModifierMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockSetExpectation{}
	}
	m.mainExpectation.input = &ModifierMockSetInput{p, p1}
	return m
}

//Return specifies results of invocation of Modifier.Set
func (m *mModifierMockSet) Return(r error) *ModifierMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ModifierMockSetExpectation{}
	}
	m.mainExpectation.result = &ModifierMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Modifier.Set is expected once
func (m *mModifierMockSet) ExpectOnce(p insolar.PulseNumber, p1 []insolar.Node) *ModifierMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &ModifierMockSetExpectation{}
	expectation.input = &ModifierMockSetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ModifierMockSetExpectation) Return(r error) {
	e.result = &ModifierMockSetResult{r}
}

//Set uses given function f as a mock of Modifier.Set method
func (m *mModifierMockSet) Set(f func(p insolar.PulseNumber, p1 []insolar.Node) (r error)) *ModifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/insolar/node.Modifier interface
func (m *ModifierMock) Set(p insolar.PulseNumber, p1 []insolar.Node) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ModifierMock.Set. %v %v", p, p1)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ModifierMockSetInput{p, p1}, "Modifier.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ModifierMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ModifierMockSetInput{p, p1}, "Modifier.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ModifierMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to ModifierMock.Set. %v %v", p, p1)
		return
	}

	return m.SetFunc(p, p1)
}

//SetMinimockCounter returns a count of ModifierMock.SetFunc invocations
func (m *ModifierMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of ModifierMock.Set invocations
func (m *ModifierMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *ModifierMock) SetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SetCounter) == uint64(len(m.SetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SetFunc != nil {
		return atomic.LoadUint64(&m.SetCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ModifierMock) ValidateCallCounters() {

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to ModifierMock.DeleteForPN")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to ModifierMock.Set")
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

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to ModifierMock.DeleteForPN")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to ModifierMock.Set")
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
		ok = ok && m.DeleteForPNFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeleteForPNFinished() {
				m.t.Error("Expected call to ModifierMock.DeleteForPN")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to ModifierMock.Set")
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

	if !m.DeleteForPNFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}
