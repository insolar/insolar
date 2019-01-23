/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ArtifactManagerMessageHandler" can be found in github.com/insolar/insolar/core
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ArtifactManagerMessageHandlerMock implements github.com/insolar/insolar/core.ArtifactManagerMessageHandler
type ArtifactManagerMessageHandlerMock struct {
	t minimock.Tester

	CloseEarlyRequestCircuitBreakerForJetFunc       func(p context.Context, p1 core.RecordID)
	CloseEarlyRequestCircuitBreakerForJetCounter    uint64
	CloseEarlyRequestCircuitBreakerForJetPreCounter uint64
	CloseEarlyRequestCircuitBreakerForJetMock       mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet

	ResetEarlyRequestCircuitBreakerFunc       func(p context.Context)
	ResetEarlyRequestCircuitBreakerCounter    uint64
	ResetEarlyRequestCircuitBreakerPreCounter uint64
	ResetEarlyRequestCircuitBreakerMock       mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker
}

//NewArtifactManagerMessageHandlerMock returns a mock for github.com/insolar/insolar/core.ArtifactManagerMessageHandler
func NewArtifactManagerMessageHandlerMock(t minimock.Tester) *ArtifactManagerMessageHandlerMock {
	m := &ArtifactManagerMessageHandlerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloseEarlyRequestCircuitBreakerForJetMock = mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet{mock: m}
	m.ResetEarlyRequestCircuitBreakerMock = mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker{mock: m}

	return m
}

type mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet struct {
	mock              *ArtifactManagerMessageHandlerMock
	mainExpectation   *ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetExpectation
	expectationSeries []*ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetExpectation
}

type ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetExpectation struct {
	input *ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetInput
}

type ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetInput struct {
	p  context.Context
	p1 core.RecordID
}

//Expect specifies that invocation of ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet is expected from 1 to Infinity times
func (m *mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet) Expect(p context.Context, p1 core.RecordID) *mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet {
	m.mock.CloseEarlyRequestCircuitBreakerForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetInput{p, p1}
	return m
}

//Return specifies results of invocation of ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet
func (m *mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet) Return() *ArtifactManagerMessageHandlerMock {
	m.mock.CloseEarlyRequestCircuitBreakerForJetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet is expected once
func (m *mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet) ExpectOnce(p context.Context, p1 core.RecordID) *ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetExpectation {
	m.mock.CloseEarlyRequestCircuitBreakerForJetFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetExpectation{}
	expectation.input = &ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet method
func (m *mArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJet) Set(f func(p context.Context, p1 core.RecordID)) *ArtifactManagerMessageHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloseEarlyRequestCircuitBreakerForJetFunc = f
	return m.mock
}

//CloseEarlyRequestCircuitBreakerForJet implements github.com/insolar/insolar/core.ArtifactManagerMessageHandler interface
func (m *ArtifactManagerMessageHandlerMock) CloseEarlyRequestCircuitBreakerForJet(p context.Context, p1 core.RecordID) {
	counter := atomic.AddUint64(&m.CloseEarlyRequestCircuitBreakerForJetPreCounter, 1)
	defer atomic.AddUint64(&m.CloseEarlyRequestCircuitBreakerForJetCounter, 1)

	if len(m.CloseEarlyRequestCircuitBreakerForJetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloseEarlyRequestCircuitBreakerForJetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJet. %v %v", p, p1)
			return
		}

		input := m.CloseEarlyRequestCircuitBreakerForJetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetInput{p, p1}, "ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet got unexpected parameters")

		return
	}

	if m.CloseEarlyRequestCircuitBreakerForJetMock.mainExpectation != nil {

		input := m.CloseEarlyRequestCircuitBreakerForJetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMessageHandlerMockCloseEarlyRequestCircuitBreakerForJetInput{p, p1}, "ArtifactManagerMessageHandler.CloseEarlyRequestCircuitBreakerForJet got unexpected parameters")
		}

		return
	}

	if m.CloseEarlyRequestCircuitBreakerForJetFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJet. %v %v", p, p1)
		return
	}

	m.CloseEarlyRequestCircuitBreakerForJetFunc(p, p1)
}

//CloseEarlyRequestCircuitBreakerForJetMinimockCounter returns a count of ArtifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJetFunc invocations
func (m *ArtifactManagerMessageHandlerMock) CloseEarlyRequestCircuitBreakerForJetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloseEarlyRequestCircuitBreakerForJetCounter)
}

//CloseEarlyRequestCircuitBreakerForJetMinimockPreCounter returns the value of ArtifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJet invocations
func (m *ArtifactManagerMessageHandlerMock) CloseEarlyRequestCircuitBreakerForJetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CloseEarlyRequestCircuitBreakerForJetPreCounter)
}

//CloseEarlyRequestCircuitBreakerForJetFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMessageHandlerMock) CloseEarlyRequestCircuitBreakerForJetFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloseEarlyRequestCircuitBreakerForJetMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloseEarlyRequestCircuitBreakerForJetCounter) == uint64(len(m.CloseEarlyRequestCircuitBreakerForJetMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloseEarlyRequestCircuitBreakerForJetMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloseEarlyRequestCircuitBreakerForJetCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloseEarlyRequestCircuitBreakerForJetFunc != nil {
		return atomic.LoadUint64(&m.CloseEarlyRequestCircuitBreakerForJetCounter) > 0
	}

	return true
}

type mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker struct {
	mock              *ArtifactManagerMessageHandlerMock
	mainExpectation   *ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerExpectation
	expectationSeries []*ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerExpectation
}

type ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerExpectation struct {
	input *ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerInput
}

type ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerInput struct {
	p context.Context
}

//Expect specifies that invocation of ArtifactManagerMessageHandler.ResetEarlyRequestCircuitBreaker is expected from 1 to Infinity times
func (m *mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker) Expect(p context.Context) *mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker {
	m.mock.ResetEarlyRequestCircuitBreakerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerExpectation{}
	}
	m.mainExpectation.input = &ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerInput{p}
	return m
}

//Return specifies results of invocation of ArtifactManagerMessageHandler.ResetEarlyRequestCircuitBreaker
func (m *mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker) Return() *ArtifactManagerMessageHandlerMock {
	m.mock.ResetEarlyRequestCircuitBreakerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of ArtifactManagerMessageHandler.ResetEarlyRequestCircuitBreaker is expected once
func (m *mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker) ExpectOnce(p context.Context) *ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerExpectation {
	m.mock.ResetEarlyRequestCircuitBreakerFunc = nil
	m.mainExpectation = nil

	expectation := &ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerExpectation{}
	expectation.input = &ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of ArtifactManagerMessageHandler.ResetEarlyRequestCircuitBreaker method
func (m *mArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreaker) Set(f func(p context.Context)) *ArtifactManagerMessageHandlerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ResetEarlyRequestCircuitBreakerFunc = f
	return m.mock
}

//ResetEarlyRequestCircuitBreaker implements github.com/insolar/insolar/core.ArtifactManagerMessageHandler interface
func (m *ArtifactManagerMessageHandlerMock) ResetEarlyRequestCircuitBreaker(p context.Context) {
	counter := atomic.AddUint64(&m.ResetEarlyRequestCircuitBreakerPreCounter, 1)
	defer atomic.AddUint64(&m.ResetEarlyRequestCircuitBreakerCounter, 1)

	if len(m.ResetEarlyRequestCircuitBreakerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ResetEarlyRequestCircuitBreakerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ArtifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreaker. %v", p)
			return
		}

		input := m.ResetEarlyRequestCircuitBreakerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerInput{p}, "ArtifactManagerMessageHandler.ResetEarlyRequestCircuitBreaker got unexpected parameters")

		return
	}

	if m.ResetEarlyRequestCircuitBreakerMock.mainExpectation != nil {

		input := m.ResetEarlyRequestCircuitBreakerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ArtifactManagerMessageHandlerMockResetEarlyRequestCircuitBreakerInput{p}, "ArtifactManagerMessageHandler.ResetEarlyRequestCircuitBreaker got unexpected parameters")
		}

		return
	}

	if m.ResetEarlyRequestCircuitBreakerFunc == nil {
		m.t.Fatalf("Unexpected call to ArtifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreaker. %v", p)
		return
	}

	m.ResetEarlyRequestCircuitBreakerFunc(p)
}

//ResetEarlyRequestCircuitBreakerMinimockCounter returns a count of ArtifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreakerFunc invocations
func (m *ArtifactManagerMessageHandlerMock) ResetEarlyRequestCircuitBreakerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ResetEarlyRequestCircuitBreakerCounter)
}

//ResetEarlyRequestCircuitBreakerMinimockPreCounter returns the value of ArtifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreaker invocations
func (m *ArtifactManagerMessageHandlerMock) ResetEarlyRequestCircuitBreakerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ResetEarlyRequestCircuitBreakerPreCounter)
}

//ResetEarlyRequestCircuitBreakerFinished returns true if mock invocations count is ok
func (m *ArtifactManagerMessageHandlerMock) ResetEarlyRequestCircuitBreakerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ResetEarlyRequestCircuitBreakerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ResetEarlyRequestCircuitBreakerCounter) == uint64(len(m.ResetEarlyRequestCircuitBreakerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ResetEarlyRequestCircuitBreakerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ResetEarlyRequestCircuitBreakerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ResetEarlyRequestCircuitBreakerFunc != nil {
		return atomic.LoadUint64(&m.ResetEarlyRequestCircuitBreakerCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) ValidateCallCounters() {

	if !m.CloseEarlyRequestCircuitBreakerForJetFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJet")
	}

	if !m.ResetEarlyRequestCircuitBreakerFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreaker")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ArtifactManagerMessageHandlerMock) MinimockFinish() {

	if !m.CloseEarlyRequestCircuitBreakerForJetFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJet")
	}

	if !m.ResetEarlyRequestCircuitBreakerFinished() {
		m.t.Fatal("Expected call to ArtifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreaker")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ArtifactManagerMessageHandlerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloseEarlyRequestCircuitBreakerForJetFinished()
		ok = ok && m.ResetEarlyRequestCircuitBreakerFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloseEarlyRequestCircuitBreakerForJetFinished() {
				m.t.Error("Expected call to ArtifactManagerMessageHandlerMock.CloseEarlyRequestCircuitBreakerForJet")
			}

			if !m.ResetEarlyRequestCircuitBreakerFinished() {
				m.t.Error("Expected call to ArtifactManagerMessageHandlerMock.ResetEarlyRequestCircuitBreaker")
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
func (m *ArtifactManagerMessageHandlerMock) AllMocksCalled() bool {

	if !m.CloseEarlyRequestCircuitBreakerForJetFinished() {
		return false
	}

	if !m.ResetEarlyRequestCircuitBreakerFinished() {
		return false
	}

	return true
}
