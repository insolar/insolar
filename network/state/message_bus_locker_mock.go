/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package state

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "messageBusLocker" can be found in github.com/insolar/insolar/network/state
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//messageBusLockerMock implements github.com/insolar/insolar/network/state.messageBusLocker
type messageBusLockerMock struct {
	t minimock.Tester

	LockFunc       func(p context.Context)
	LockCounter    uint64
	LockPreCounter uint64
	LockMock       mmessageBusLockerMockLock

	UnlockFunc       func(p context.Context)
	UnlockCounter    uint64
	UnlockPreCounter uint64
	UnlockMock       mmessageBusLockerMockUnlock
}

//NewmessageBusLockerMock returns a mock for github.com/insolar/insolar/network/state.messageBusLocker
func NewmessageBusLockerMock(t minimock.Tester) *messageBusLockerMock {
	m := &messageBusLockerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.LockMock = mmessageBusLockerMockLock{mock: m}
	m.UnlockMock = mmessageBusLockerMockUnlock{mock: m}

	return m
}

type mmessageBusLockerMockLock struct {
	mock              *messageBusLockerMock
	mainExpectation   *messageBusLockerMockLockExpectation
	expectationSeries []*messageBusLockerMockLockExpectation
}

type messageBusLockerMockLockExpectation struct {
	input *messageBusLockerMockLockInput
}

type messageBusLockerMockLockInput struct {
	p context.Context
}

//Expect specifies that invocation of messageBusLocker.Lock is expected from 1 to Infinity times
func (m *mmessageBusLockerMockLock) Expect(p context.Context) *mmessageBusLockerMockLock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockLockExpectation{}
	}
	m.mainExpectation.input = &messageBusLockerMockLockInput{p}
	return m
}

//Return specifies results of invocation of messageBusLocker.Lock
func (m *mmessageBusLockerMockLock) Return() *messageBusLockerMock {
	m.mock.LockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockLockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of messageBusLocker.Lock is expected once
func (m *mmessageBusLockerMockLock) ExpectOnce(p context.Context) *messageBusLockerMockLockExpectation {
	m.mock.LockFunc = nil
	m.mainExpectation = nil

	expectation := &messageBusLockerMockLockExpectation{}
	expectation.input = &messageBusLockerMockLockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of messageBusLocker.Lock method
func (m *mmessageBusLockerMockLock) Set(f func(p context.Context)) *messageBusLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LockFunc = f
	return m.mock
}

//Lock implements github.com/insolar/insolar/network/state.messageBusLocker interface
func (m *messageBusLockerMock) Lock(p context.Context) {
	counter := atomic.AddUint64(&m.LockPreCounter, 1)
	defer atomic.AddUint64(&m.LockCounter, 1)

	if len(m.LockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to messageBusLockerMock.Lock. %v", p)
			return
		}

		input := m.LockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, messageBusLockerMockLockInput{p}, "messageBusLocker.Lock got unexpected parameters")

		return
	}

	if m.LockMock.mainExpectation != nil {

		input := m.LockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, messageBusLockerMockLockInput{p}, "messageBusLocker.Lock got unexpected parameters")
		}

		return
	}

	if m.LockFunc == nil {
		m.t.Fatalf("Unexpected call to messageBusLockerMock.Lock. %v", p)
		return
	}

	m.LockFunc(p)
}

//LockMinimockCounter returns a count of messageBusLockerMock.LockFunc invocations
func (m *messageBusLockerMock) LockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LockCounter)
}

//LockMinimockPreCounter returns the value of messageBusLockerMock.Lock invocations
func (m *messageBusLockerMock) LockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LockPreCounter)
}

//LockFinished returns true if mock invocations count is ok
func (m *messageBusLockerMock) LockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LockCounter) == uint64(len(m.LockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LockFunc != nil {
		return atomic.LoadUint64(&m.LockCounter) > 0
	}

	return true
}

type mmessageBusLockerMockUnlock struct {
	mock              *messageBusLockerMock
	mainExpectation   *messageBusLockerMockUnlockExpectation
	expectationSeries []*messageBusLockerMockUnlockExpectation
}

type messageBusLockerMockUnlockExpectation struct {
	input *messageBusLockerMockUnlockInput
}

type messageBusLockerMockUnlockInput struct {
	p context.Context
}

//Expect specifies that invocation of messageBusLocker.Unlock is expected from 1 to Infinity times
func (m *mmessageBusLockerMockUnlock) Expect(p context.Context) *mmessageBusLockerMockUnlock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockUnlockExpectation{}
	}
	m.mainExpectation.input = &messageBusLockerMockUnlockInput{p}
	return m
}

//Return specifies results of invocation of messageBusLocker.Unlock
func (m *mmessageBusLockerMockUnlock) Return() *messageBusLockerMock {
	m.mock.UnlockFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &messageBusLockerMockUnlockExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of messageBusLocker.Unlock is expected once
func (m *mmessageBusLockerMockUnlock) ExpectOnce(p context.Context) *messageBusLockerMockUnlockExpectation {
	m.mock.UnlockFunc = nil
	m.mainExpectation = nil

	expectation := &messageBusLockerMockUnlockExpectation{}
	expectation.input = &messageBusLockerMockUnlockInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of messageBusLocker.Unlock method
func (m *mmessageBusLockerMockUnlock) Set(f func(p context.Context)) *messageBusLockerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UnlockFunc = f
	return m.mock
}

//Unlock implements github.com/insolar/insolar/network/state.messageBusLocker interface
func (m *messageBusLockerMock) Unlock(p context.Context) {
	counter := atomic.AddUint64(&m.UnlockPreCounter, 1)
	defer atomic.AddUint64(&m.UnlockCounter, 1)

	if len(m.UnlockMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UnlockMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to messageBusLockerMock.Unlock. %v", p)
			return
		}

		input := m.UnlockMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, messageBusLockerMockUnlockInput{p}, "messageBusLocker.Unlock got unexpected parameters")

		return
	}

	if m.UnlockMock.mainExpectation != nil {

		input := m.UnlockMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, messageBusLockerMockUnlockInput{p}, "messageBusLocker.Unlock got unexpected parameters")
		}

		return
	}

	if m.UnlockFunc == nil {
		m.t.Fatalf("Unexpected call to messageBusLockerMock.Unlock. %v", p)
		return
	}

	m.UnlockFunc(p)
}

//UnlockMinimockCounter returns a count of messageBusLockerMock.UnlockFunc invocations
func (m *messageBusLockerMock) UnlockMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockCounter)
}

//UnlockMinimockPreCounter returns the value of messageBusLockerMock.Unlock invocations
func (m *messageBusLockerMock) UnlockMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UnlockPreCounter)
}

//UnlockFinished returns true if mock invocations count is ok
func (m *messageBusLockerMock) UnlockFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UnlockMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UnlockCounter) == uint64(len(m.UnlockMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UnlockMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UnlockFunc != nil {
		return atomic.LoadUint64(&m.UnlockCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *messageBusLockerMock) ValidateCallCounters() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Unlock")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *messageBusLockerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *messageBusLockerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *messageBusLockerMock) MinimockFinish() {

	if !m.LockFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Lock")
	}

	if !m.UnlockFinished() {
		m.t.Fatal("Expected call to messageBusLockerMock.Unlock")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *messageBusLockerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *messageBusLockerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.LockFinished()
		ok = ok && m.UnlockFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.LockFinished() {
				m.t.Error("Expected call to messageBusLockerMock.Lock")
			}

			if !m.UnlockFinished() {
				m.t.Error("Expected call to messageBusLockerMock.Unlock")
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
func (m *messageBusLockerMock) AllMocksCalled() bool {

	if !m.LockFinished() {
		return false
	}

	if !m.UnlockFinished() {
		return false
	}

	return true
}
