package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Accessor" can be found in github.com/insolar/insolar/insolar/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//AccessorMock implements github.com/insolar/insolar/insolar/jet.Accessor
type AccessorMock struct {
	t minimock.Tester

	AllFunc       func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)
	AllCounter    uint64
	AllPreCounter uint64
	AllMock       mAccessorMockAll

	ForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r insolar.JetID, r1 bool)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mAccessorMockForID
}

//NewAccessorMock returns a mock for github.com/insolar/insolar/insolar/jet.Accessor
func NewAccessorMock(t minimock.Tester) *AccessorMock {
	m := &AccessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AllMock = mAccessorMockAll{mock: m}
	m.ForIDMock = mAccessorMockForID{mock: m}

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
	p  context.Context
	p1 insolar.PulseNumber
}

type AccessorMockAllResult struct {
	r []insolar.JetID
}

//Expect specifies that invocation of Accessor.All is expected from 1 to Infinity times
func (m *mAccessorMockAll) Expect(p context.Context, p1 insolar.PulseNumber) *mAccessorMockAll {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockAllExpectation{}
	}
	m.mainExpectation.input = &AccessorMockAllInput{p, p1}
	return m
}

//Return specifies results of invocation of Accessor.All
func (m *mAccessorMockAll) Return(r []insolar.JetID) *AccessorMock {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockAllExpectation{}
	}
	m.mainExpectation.result = &AccessorMockAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.All is expected once
func (m *mAccessorMockAll) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *AccessorMockAllExpectation {
	m.mock.AllFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockAllExpectation{}
	expectation.input = &AccessorMockAllInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockAllExpectation) Return(r []insolar.JetID) {
	e.result = &AccessorMockAllResult{r}
}

//Set uses given function f as a mock of Accessor.All method
func (m *mAccessorMockAll) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AllFunc = f
	return m.mock
}

//All implements github.com/insolar/insolar/insolar/jet.Accessor interface
func (m *AccessorMock) All(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID) {
	counter := atomic.AddUint64(&m.AllPreCounter, 1)
	defer atomic.AddUint64(&m.AllCounter, 1)

	if len(m.AllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.All. %v %v", p, p1)
			return
		}

		input := m.AllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockAllInput{p, p1}, "Accessor.All got unexpected parameters")

		result := m.AllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.All")
			return
		}

		r = result.r

		return
	}

	if m.AllMock.mainExpectation != nil {

		input := m.AllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockAllInput{p, p1}, "Accessor.All got unexpected parameters")
		}

		result := m.AllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.All")
		}

		r = result.r

		return
	}

	if m.AllFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.All. %v %v", p, p1)
		return
	}

	return m.AllFunc(p, p1)
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

type mAccessorMockForID struct {
	mock              *AccessorMock
	mainExpectation   *AccessorMockForIDExpectation
	expectationSeries []*AccessorMockForIDExpectation
}

type AccessorMockForIDExpectation struct {
	input  *AccessorMockForIDInput
	result *AccessorMockForIDResult
}

type AccessorMockForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type AccessorMockForIDResult struct {
	r  insolar.JetID
	r1 bool
}

//Expect specifies that invocation of Accessor.ForID is expected from 1 to Infinity times
func (m *mAccessorMockForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mAccessorMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockForIDExpectation{}
	}
	m.mainExpectation.input = &AccessorMockForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Accessor.ForID
func (m *mAccessorMockForID) Return(r insolar.JetID, r1 bool) *AccessorMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &AccessorMockForIDExpectation{}
	}
	m.mainExpectation.result = &AccessorMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Accessor.ForID is expected once
func (m *mAccessorMockForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *AccessorMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &AccessorMockForIDExpectation{}
	expectation.input = &AccessorMockForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *AccessorMockForIDExpectation) Return(r insolar.JetID, r1 bool) {
	e.result = &AccessorMockForIDResult{r, r1}
}

//Set uses given function f as a mock of Accessor.ForID method
func (m *mAccessorMockForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r insolar.JetID, r1 bool)) *AccessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/insolar/jet.Accessor interface
func (m *AccessorMock) ForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r insolar.JetID, r1 bool) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to AccessorMock.ForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, AccessorMockForIDInput{p, p1, p2}, "Accessor.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, AccessorMockForIDInput{p, p1, p2}, "Accessor.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the AccessorMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to AccessorMock.ForID. %v %v %v", p, p1, p2)
		return
	}

	return m.ForIDFunc(p, p1, p2)
}

//ForIDMinimockCounter returns a count of AccessorMock.ForIDFunc invocations
func (m *AccessorMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of AccessorMock.ForID invocations
func (m *AccessorMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *AccessorMock) ForIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ForIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ForIDCounter) == uint64(len(m.ForIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ForIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ForIDFunc != nil {
		return atomic.LoadUint64(&m.ForIDCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *AccessorMock) ValidateCallCounters() {

	if !m.AllFinished() {
		m.t.Fatal("Expected call to AccessorMock.All")
	}

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to AccessorMock.ForID")
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

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to AccessorMock.ForID")
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
		ok = ok && m.ForIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AllFinished() {
				m.t.Error("Expected call to AccessorMock.All")
			}

			if !m.ForIDFinished() {
				m.t.Error("Expected call to AccessorMock.ForID")
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

	if !m.ForIDFinished() {
		return false
	}

	return true
}
