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
	packets "github.com/insolar/insolar/consensus/packets"
)

//ClaimQueueMock implements github.com/insolar/insolar/network.ClaimQueue
type ClaimQueueMock struct {
	t minimock.Tester

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
}

//NewClaimQueueMock returns a mock for github.com/insolar/insolar/network.ClaimQueue
func NewClaimQueueMock(t minimock.Tester) *ClaimQueueMock {
	m := &ClaimQueueMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FrontMock = mClaimQueueMockFront{mock: m}
	m.LengthMock = mClaimQueueMockLength{mock: m}
	m.PopMock = mClaimQueueMockPop{mock: m}

	return m
}

type mClaimQueueMockFront struct {
	mock *ClaimQueueMock
}

//Return sets up a mock for ClaimQueue.Front to return Return's arguments
func (m *mClaimQueueMockFront) Return(r packets.ReferendumClaim) *ClaimQueueMock {
	m.mock.FrontFunc = func() packets.ReferendumClaim {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ClaimQueue.Front method
func (m *mClaimQueueMockFront) Set(f func() (r packets.ReferendumClaim)) *ClaimQueueMock {
	m.mock.FrontFunc = f

	return m.mock
}

//Front implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Front() (r packets.ReferendumClaim) {
	atomic.AddUint64(&m.FrontPreCounter, 1)
	defer atomic.AddUint64(&m.FrontCounter, 1)

	if m.FrontFunc == nil {
		m.t.Fatal("Unexpected call to ClaimQueueMock.Front")
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

type mClaimQueueMockLength struct {
	mock *ClaimQueueMock
}

//Return sets up a mock for ClaimQueue.Length to return Return's arguments
func (m *mClaimQueueMockLength) Return(r int) *ClaimQueueMock {
	m.mock.LengthFunc = func() int {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ClaimQueue.Length method
func (m *mClaimQueueMockLength) Set(f func() (r int)) *ClaimQueueMock {
	m.mock.LengthFunc = f

	return m.mock
}

//Length implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Length() (r int) {
	atomic.AddUint64(&m.LengthPreCounter, 1)
	defer atomic.AddUint64(&m.LengthCounter, 1)

	if m.LengthFunc == nil {
		m.t.Fatal("Unexpected call to ClaimQueueMock.Length")
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

type mClaimQueueMockPop struct {
	mock *ClaimQueueMock
}

//Return sets up a mock for ClaimQueue.Pop to return Return's arguments
func (m *mClaimQueueMockPop) Return(r packets.ReferendumClaim) *ClaimQueueMock {
	m.mock.PopFunc = func() packets.ReferendumClaim {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of ClaimQueue.Pop method
func (m *mClaimQueueMockPop) Set(f func() (r packets.ReferendumClaim)) *ClaimQueueMock {
	m.mock.PopFunc = f

	return m.mock
}

//Pop implements github.com/insolar/insolar/network.ClaimQueue interface
func (m *ClaimQueueMock) Pop() (r packets.ReferendumClaim) {
	atomic.AddUint64(&m.PopPreCounter, 1)
	defer atomic.AddUint64(&m.PopCounter, 1)

	if m.PopFunc == nil {
		m.t.Fatal("Unexpected call to ClaimQueueMock.Pop")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ClaimQueueMock) ValidateCallCounters() {

	if m.FrontFunc != nil && atomic.LoadUint64(&m.FrontCounter) == 0 {
		m.t.Fatal("Expected call to ClaimQueueMock.Front")
	}

	if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
		m.t.Fatal("Expected call to ClaimQueueMock.Length")
	}

	if m.PopFunc != nil && atomic.LoadUint64(&m.PopCounter) == 0 {
		m.t.Fatal("Expected call to ClaimQueueMock.Pop")
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

	if m.FrontFunc != nil && atomic.LoadUint64(&m.FrontCounter) == 0 {
		m.t.Fatal("Expected call to ClaimQueueMock.Front")
	}

	if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
		m.t.Fatal("Expected call to ClaimQueueMock.Length")
	}

	if m.PopFunc != nil && atomic.LoadUint64(&m.PopCounter) == 0 {
		m.t.Fatal("Expected call to ClaimQueueMock.Pop")
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
		ok = ok && (m.FrontFunc == nil || atomic.LoadUint64(&m.FrontCounter) > 0)
		ok = ok && (m.LengthFunc == nil || atomic.LoadUint64(&m.LengthCounter) > 0)
		ok = ok && (m.PopFunc == nil || atomic.LoadUint64(&m.PopCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.FrontFunc != nil && atomic.LoadUint64(&m.FrontCounter) == 0 {
				m.t.Error("Expected call to ClaimQueueMock.Front")
			}

			if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
				m.t.Error("Expected call to ClaimQueueMock.Length")
			}

			if m.PopFunc != nil && atomic.LoadUint64(&m.PopCounter) == 0 {
				m.t.Error("Expected call to ClaimQueueMock.Pop")
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

	if m.FrontFunc != nil && atomic.LoadUint64(&m.FrontCounter) == 0 {
		return false
	}

	if m.LengthFunc != nil && atomic.LoadUint64(&m.LengthCounter) == 0 {
		return false
	}

	if m.PopFunc != nil && atomic.LoadUint64(&m.PopCounter) == 0 {
		return false
	}

	return true
}
