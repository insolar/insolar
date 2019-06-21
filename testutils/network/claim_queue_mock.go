package network

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ClaimQueue" can be found in github.com/insolar/insolar/network
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	packets "github.com/insolar/insolar/network/consensusv1/packets"

	testify_assert "github.com/stretchr/testify/assert"
)

//ClaimQueueMock implements github.com/insolar/insolar/network.ClaimQueue
type ClaimQueueMock struct {
	t minimock.Tester

	ClearFunc       func()
	ClearCounter    uint64
	ClearPreCounter uint64
	ClearMock       mClaimQueueMockClear

	FrontFunc       func() (r packets.ReferendumClaim)
	FrontCounter    uint64
	FrontPreCounter uint64
	FrontMock       mClaimQueueMockFront

	LengthFunc       func() (r int)
	LengthCounter    uint64
	LengthPreCounter uint64
	LengthMock       mClaimQueueMockLength

	PopFunc       func() (r packets.ReferendumClaim)
	PopCounter    uint64
	PopPreCounter uint64
	PopMock       mClaimQueueMockPop

	PushFunc       func(p packets.ReferendumClaim)
	PushCounter    uint64
	PushPreCounter uint64
	PushMock       mClaimQueueMockPush
}

//NewClaimQueueMock returns a mock for github.com/insolar/insolar/network.ClaimQueue
func NewClaimQueueMock(t minimock.Tester) *ClaimQueueMock {
	m := &ClaimQueueMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ClearMock = mClaimQueueMockClear{mock: m}
	m.FrontMock = mClaimQueueMockFront{mock: m}
	m.LengthMock = mClaimQueueMockLength{mock: m}
	m.PopMock = mClaimQueueMockPop{mock: m}
	m.PushMock = mClaimQueueMockPush{mock: m}

	return m
}

type mClaimQueueMockClear struct {
	mock              *ClaimQueueMock
	mainExpectation   *ClaimQueueMockClearExpectation
	expectationSeries []*ClaimQueueMockClearExpectation
}

type ClaimQueueMockClearExpectation struct {
}

//Expect specifies that invocation of ClaimQueue.Clear is expected from 1 to Infinity times
func (m *mClaimQueueMockClear) Expect() *mClaimQueueMockClear {
	m.mock.ClearFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockClearExpectation{}
	}

	return m
}

//Return specifies results of invocation of ClaimQueue.Clear
func (m *mClaimQueueMockClear) Return() *ClaimQueueMock {
	m.mock.ClearFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockClearExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ClaimQueue.Clear is expected once
func (m *mClaimQueueMockClear) ExpectOnce() *ClaimQueueMockClearExpectation {
	m.mock.ClearFunc = nil
	m.mainExpectation = nil

	expectation := &ClaimQueueMockClearExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ClaimQueue.Clear method
func (m *mClaimQueueMockClear) Set(f func()) *ClaimQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ClearFunc = f
	return m.mock
}

//Clear implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Clear() {
	counter := atomic.AddUint64(&m.ClearPreCounter, 1)
	defer atomic.AddUint64(&m.ClearCounter, 1)

	if len(m.ClearMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ClearMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClaimQueueMock.Clear.")
			return
		}

		return
	}

	if m.ClearMock.mainExpectation != nil {

		return
	}

	if m.ClearFunc == nil {
		m.t.Fatalf("Unexpected call to ClaimQueueMock.Clear.")
		return
	}

	m.ClearFunc()
}

//ClearMinimockCounter returns a count of ClaimQueueMock.ClearFunc invocations
func (m *ClaimQueueMock) ClearMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ClearCounter)
}

//ClearMinimockPreCounter returns the value of ClaimQueueMock.Clear invocations
func (m *ClaimQueueMock) ClearMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClearPreCounter)
}

//ClearFinished returns true if mock invocations count is ok
func (m *ClaimQueueMock) ClearFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ClearMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ClearCounter) == uint64(len(m.ClearMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ClearMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ClearCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ClearFunc != nil {
		return atomic.LoadUint64(&m.ClearCounter) > 0
	}

	return true
}

type mClaimQueueMockFront struct {
	mock              *ClaimQueueMock
	mainExpectation   *ClaimQueueMockFrontExpectation
	expectationSeries []*ClaimQueueMockFrontExpectation
}

type ClaimQueueMockFrontExpectation struct {
	result *ClaimQueueMockFrontResult
}

type ClaimQueueMockFrontResult struct {
	r packets.ReferendumClaim
}

//Expect specifies that invocation of ClaimQueue.Front is expected from 1 to Infinity times
func (m *mClaimQueueMockFront) Expect() *mClaimQueueMockFront {
	m.mock.FrontFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockFrontExpectation{}
	}

	return m
}

//Return specifies results of invocation of ClaimQueue.Front
func (m *mClaimQueueMockFront) Return(r packets.ReferendumClaim) *ClaimQueueMock {
	m.mock.FrontFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockFrontExpectation{}
	}
	m.mainExpectation.result = &ClaimQueueMockFrontResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ClaimQueue.Front is expected once
func (m *mClaimQueueMockFront) ExpectOnce() *ClaimQueueMockFrontExpectation {
	m.mock.FrontFunc = nil
	m.mainExpectation = nil

	expectation := &ClaimQueueMockFrontExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClaimQueueMockFrontExpectation) Return(r packets.ReferendumClaim) {
	e.result = &ClaimQueueMockFrontResult{r}
}

//Set uses given function f as a mock of ClaimQueue.Front method
func (m *mClaimQueueMockFront) Set(f func() (r packets.ReferendumClaim)) *ClaimQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.FrontFunc = f
	return m.mock
}

//Front implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Front() (r packets.ReferendumClaim) {
	counter := atomic.AddUint64(&m.FrontPreCounter, 1)
	defer atomic.AddUint64(&m.FrontCounter, 1)

	if len(m.FrontMock.expectationSeries) > 0 {
		if counter > uint64(len(m.FrontMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClaimQueueMock.Front.")
			return
		}

		result := m.FrontMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClaimQueueMock.Front")
			return
		}

		r = result.r

		return
	}

	if m.FrontMock.mainExpectation != nil {

		result := m.FrontMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClaimQueueMock.Front")
		}

		r = result.r

		return
	}

	if m.FrontFunc == nil {
		m.t.Fatalf("Unexpected call to ClaimQueueMock.Front.")
		return
	}

	return m.FrontFunc()
}

//FrontMinimockCounter returns a count of ClaimQueueMock.FrontFunc invocations
func (m *ClaimQueueMock) FrontMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FrontCounter)
}

//FrontMinimockPreCounter returns the value of ClaimQueueMock.Front invocations
func (m *ClaimQueueMock) FrontMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FrontPreCounter)
}

//FrontFinished returns true if mock invocations count is ok
func (m *ClaimQueueMock) FrontFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.FrontMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.FrontCounter) == uint64(len(m.FrontMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.FrontMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.FrontCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.FrontFunc != nil {
		return atomic.LoadUint64(&m.FrontCounter) > 0
	}

	return true
}

type mClaimQueueMockLength struct {
	mock              *ClaimQueueMock
	mainExpectation   *ClaimQueueMockLengthExpectation
	expectationSeries []*ClaimQueueMockLengthExpectation
}

type ClaimQueueMockLengthExpectation struct {
	result *ClaimQueueMockLengthResult
}

type ClaimQueueMockLengthResult struct {
	r int
}

//Expect specifies that invocation of ClaimQueue.Length is expected from 1 to Infinity times
func (m *mClaimQueueMockLength) Expect() *mClaimQueueMockLength {
	m.mock.LengthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockLengthExpectation{}
	}

	return m
}

//Return specifies results of invocation of ClaimQueue.Length
func (m *mClaimQueueMockLength) Return(r int) *ClaimQueueMock {
	m.mock.LengthFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockLengthExpectation{}
	}
	m.mainExpectation.result = &ClaimQueueMockLengthResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ClaimQueue.Length is expected once
func (m *mClaimQueueMockLength) ExpectOnce() *ClaimQueueMockLengthExpectation {
	m.mock.LengthFunc = nil
	m.mainExpectation = nil

	expectation := &ClaimQueueMockLengthExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClaimQueueMockLengthExpectation) Return(r int) {
	e.result = &ClaimQueueMockLengthResult{r}
}

//Set uses given function f as a mock of ClaimQueue.Length method
func (m *mClaimQueueMockLength) Set(f func() (r int)) *ClaimQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LengthFunc = f
	return m.mock
}

//Length implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Length() (r int) {
	counter := atomic.AddUint64(&m.LengthPreCounter, 1)
	defer atomic.AddUint64(&m.LengthCounter, 1)

	if len(m.LengthMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LengthMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClaimQueueMock.Length.")
			return
		}

		result := m.LengthMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClaimQueueMock.Length")
			return
		}

		r = result.r

		return
	}

	if m.LengthMock.mainExpectation != nil {

		result := m.LengthMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClaimQueueMock.Length")
		}

		r = result.r

		return
	}

	if m.LengthFunc == nil {
		m.t.Fatalf("Unexpected call to ClaimQueueMock.Length.")
		return
	}

	return m.LengthFunc()
}

//LengthMinimockCounter returns a count of ClaimQueueMock.LengthFunc invocations
func (m *ClaimQueueMock) LengthMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LengthCounter)
}

//LengthMinimockPreCounter returns the value of ClaimQueueMock.Length invocations
func (m *ClaimQueueMock) LengthMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LengthPreCounter)
}

//LengthFinished returns true if mock invocations count is ok
func (m *ClaimQueueMock) LengthFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LengthMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LengthCounter) == uint64(len(m.LengthMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LengthMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LengthCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LengthFunc != nil {
		return atomic.LoadUint64(&m.LengthCounter) > 0
	}

	return true
}

type mClaimQueueMockPop struct {
	mock              *ClaimQueueMock
	mainExpectation   *ClaimQueueMockPopExpectation
	expectationSeries []*ClaimQueueMockPopExpectation
}

type ClaimQueueMockPopExpectation struct {
	result *ClaimQueueMockPopResult
}

type ClaimQueueMockPopResult struct {
	r packets.ReferendumClaim
}

//Expect specifies that invocation of ClaimQueue.Pop is expected from 1 to Infinity times
func (m *mClaimQueueMockPop) Expect() *mClaimQueueMockPop {
	m.mock.PopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockPopExpectation{}
	}

	return m
}

//Return specifies results of invocation of ClaimQueue.Pop
func (m *mClaimQueueMockPop) Return(r packets.ReferendumClaim) *ClaimQueueMock {
	m.mock.PopFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockPopExpectation{}
	}
	m.mainExpectation.result = &ClaimQueueMockPopResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ClaimQueue.Pop is expected once
func (m *mClaimQueueMockPop) ExpectOnce() *ClaimQueueMockPopExpectation {
	m.mock.PopFunc = nil
	m.mainExpectation = nil

	expectation := &ClaimQueueMockPopExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ClaimQueueMockPopExpectation) Return(r packets.ReferendumClaim) {
	e.result = &ClaimQueueMockPopResult{r}
}

//Set uses given function f as a mock of ClaimQueue.Pop method
func (m *mClaimQueueMockPop) Set(f func() (r packets.ReferendumClaim)) *ClaimQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PopFunc = f
	return m.mock
}

//Pop implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Pop() (r packets.ReferendumClaim) {
	counter := atomic.AddUint64(&m.PopPreCounter, 1)
	defer atomic.AddUint64(&m.PopCounter, 1)

	if len(m.PopMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PopMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClaimQueueMock.Pop.")
			return
		}

		result := m.PopMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ClaimQueueMock.Pop")
			return
		}

		r = result.r

		return
	}

	if m.PopMock.mainExpectation != nil {

		result := m.PopMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ClaimQueueMock.Pop")
		}

		r = result.r

		return
	}

	if m.PopFunc == nil {
		m.t.Fatalf("Unexpected call to ClaimQueueMock.Pop.")
		return
	}

	return m.PopFunc()
}

//PopMinimockCounter returns a count of ClaimQueueMock.PopFunc invocations
func (m *ClaimQueueMock) PopMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PopCounter)
}

//PopMinimockPreCounter returns the value of ClaimQueueMock.Pop invocations
func (m *ClaimQueueMock) PopMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PopPreCounter)
}

//PopFinished returns true if mock invocations count is ok
func (m *ClaimQueueMock) PopFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PopMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PopCounter) == uint64(len(m.PopMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PopMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PopCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PopFunc != nil {
		return atomic.LoadUint64(&m.PopCounter) > 0
	}

	return true
}

type mClaimQueueMockPush struct {
	mock              *ClaimQueueMock
	mainExpectation   *ClaimQueueMockPushExpectation
	expectationSeries []*ClaimQueueMockPushExpectation
}

type ClaimQueueMockPushExpectation struct {
	input *ClaimQueueMockPushInput
}

type ClaimQueueMockPushInput struct {
	p packets.ReferendumClaim
}

//Expect specifies that invocation of ClaimQueue.Push is expected from 1 to Infinity times
func (m *mClaimQueueMockPush) Expect(p packets.ReferendumClaim) *mClaimQueueMockPush {
	m.mock.PushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockPushExpectation{}
	}
	m.mainExpectation.input = &ClaimQueueMockPushInput{p}
	return m
}

//Return specifies results of invocation of ClaimQueue.Push
func (m *mClaimQueueMockPush) Return() *ClaimQueueMock {
	m.mock.PushFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ClaimQueueMockPushExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ClaimQueue.Push is expected once
func (m *mClaimQueueMockPush) ExpectOnce(p packets.ReferendumClaim) *ClaimQueueMockPushExpectation {
	m.mock.PushFunc = nil
	m.mainExpectation = nil

	expectation := &ClaimQueueMockPushExpectation{}
	expectation.input = &ClaimQueueMockPushInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ClaimQueue.Push method
func (m *mClaimQueueMockPush) Set(f func(p packets.ReferendumClaim)) *ClaimQueueMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PushFunc = f
	return m.mock
}

//Push implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Push(p packets.ReferendumClaim) {
	counter := atomic.AddUint64(&m.PushPreCounter, 1)
	defer atomic.AddUint64(&m.PushCounter, 1)

	if len(m.PushMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PushMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ClaimQueueMock.Push. %v", p)
			return
		}

		input := m.PushMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ClaimQueueMockPushInput{p}, "ClaimQueue.Push got unexpected parameters")

		return
	}

	if m.PushMock.mainExpectation != nil {

		input := m.PushMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ClaimQueueMockPushInput{p}, "ClaimQueue.Push got unexpected parameters")
		}

		return
	}

	if m.PushFunc == nil {
		m.t.Fatalf("Unexpected call to ClaimQueueMock.Push. %v", p)
		return
	}

	m.PushFunc(p)
}

//PushMinimockCounter returns a count of ClaimQueueMock.PushFunc invocations
func (m *ClaimQueueMock) PushMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PushCounter)
}

//PushMinimockPreCounter returns the value of ClaimQueueMock.Push invocations
func (m *ClaimQueueMock) PushMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PushPreCounter)
}

//PushFinished returns true if mock invocations count is ok
func (m *ClaimQueueMock) PushFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PushMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PushCounter) == uint64(len(m.PushMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PushMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PushCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PushFunc != nil {
		return atomic.LoadUint64(&m.PushCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ClaimQueueMock) ValidateCallCounters() {

	if !m.ClearFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Clear")
	}

	if !m.FrontFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Front")
	}

	if !m.LengthFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Length")
	}

	if !m.PopFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Pop")
	}

	if !m.PushFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Push")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ClaimQueueMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ClaimQueueMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ClaimQueueMock) MinimockFinish() {

	if !m.ClearFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Clear")
	}

	if !m.FrontFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Front")
	}

	if !m.LengthFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Length")
	}

	if !m.PopFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Pop")
	}

	if !m.PushFinished() {
		m.t.Fatal("Expected call to ClaimQueueMock.Push")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ClaimQueueMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ClaimQueueMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ClearFinished()
		ok = ok && m.FrontFinished()
		ok = ok && m.LengthFinished()
		ok = ok && m.PopFinished()
		ok = ok && m.PushFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ClearFinished() {
				m.t.Error("Expected call to ClaimQueueMock.Clear")
			}

			if !m.FrontFinished() {
				m.t.Error("Expected call to ClaimQueueMock.Front")
			}

			if !m.LengthFinished() {
				m.t.Error("Expected call to ClaimQueueMock.Length")
			}

			if !m.PopFinished() {
				m.t.Error("Expected call to ClaimQueueMock.Pop")
			}

			if !m.PushFinished() {
				m.t.Error("Expected call to ClaimQueueMock.Push")
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
func (m *ClaimQueueMock) AllMocksCalled() bool {

	if !m.ClearFinished() {
		return false
	}

	if !m.FrontFinished() {
		return false
	}

	if !m.LengthFinished() {
		return false
	}

	if !m.PopFinished() {
		return false
	}

	if !m.PushFinished() {
		return false
	}

	return true
}
