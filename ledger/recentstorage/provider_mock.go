/*
 *    Copyright 2019 Insolar
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

package recentstorage

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Provider" can be found in github.com/insolar/insolar/ledger/recentstorage
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProviderMock implements github.com/insolar/insolar/ledger/recentstorage.Provider
type ProviderMock struct {
	t minimock.Tester

	CloneStorageFunc       func(p core.RecordID, p1 core.RecordID)
	CloneStorageCounter    uint64
	CloneStoragePreCounter uint64
	CloneStorageMock       mProviderMockCloneStorage

	GetStorageFunc       func(p core.RecordID) (r RecentStorage)
	GetStorageCounter    uint64
	GetStoragePreCounter uint64
	GetStorageMock       mProviderMockGetStorage
}

//NewProviderMock returns a mock for github.com/insolar/insolar/ledger/recentstorage.Provider
func NewProviderMock(t minimock.Tester) *ProviderMock {
	m := &ProviderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CloneStorageMock = mProviderMockCloneStorage{mock: m}
	m.GetStorageMock = mProviderMockGetStorage{mock: m}

	return m
}

type mProviderMockCloneStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockCloneStorageExpectation
	expectationSeries []*ProviderMockCloneStorageExpectation
}

type ProviderMockCloneStorageExpectation struct {
	input *ProviderMockCloneStorageInput
}

type ProviderMockCloneStorageInput struct {
	p  core.RecordID
	p1 core.RecordID
}

//Expect specifies that invocation of Provider.CloneStorage is expected from 1 to Infinity times
func (m *mProviderMockCloneStorage) Expect(p core.RecordID, p1 core.RecordID) *mProviderMockCloneStorage {
	m.mock.CloneStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCloneStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockCloneStorageInput{p, p1}
	return m
}

//Return specifies results of invocation of Provider.CloneStorage
func (m *mProviderMockCloneStorage) Return() *ProviderMock {
	m.mock.CloneStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockCloneStorageExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Provider.CloneStorage is expected once
func (m *mProviderMockCloneStorage) ExpectOnce(p core.RecordID, p1 core.RecordID) *ProviderMockCloneStorageExpectation {
	m.mock.CloneStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockCloneStorageExpectation{}
	expectation.input = &ProviderMockCloneStorageInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Provider.CloneStorage method
func (m *mProviderMockCloneStorage) Set(f func(p core.RecordID, p1 core.RecordID)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloneStorageFunc = f
	return m.mock
}

//CloneStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) CloneStorage(p core.RecordID, p1 core.RecordID) {
	counter := atomic.AddUint64(&m.CloneStoragePreCounter, 1)
	defer atomic.AddUint64(&m.CloneStorageCounter, 1)

	if len(m.CloneStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloneStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.CloneStorage. %v %v", p, p1)
			return
		}

		input := m.CloneStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockCloneStorageInput{p, p1}, "Provider.CloneStorage got unexpected parameters")

		return
	}

	if m.CloneStorageMock.mainExpectation != nil {

		input := m.CloneStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockCloneStorageInput{p, p1}, "Provider.CloneStorage got unexpected parameters")
		}

		return
	}

	if m.CloneStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.CloneStorage. %v %v", p, p1)
		return
	}

	m.CloneStorageFunc(p, p1)
}

//CloneStorageMinimockCounter returns a count of ProviderMock.CloneStorageFunc invocations
func (m *ProviderMock) CloneStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloneStorageCounter)
}

//CloneStorageMinimockPreCounter returns the value of ProviderMock.CloneStorage invocations
func (m *ProviderMock) CloneStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CloneStoragePreCounter)
}

//CloneStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) CloneStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CloneStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CloneStorageCounter) == uint64(len(m.CloneStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CloneStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CloneStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CloneStorageFunc != nil {
		return atomic.LoadUint64(&m.CloneStorageCounter) > 0
	}

	return true
}

type mProviderMockGetStorage struct {
	mock              *ProviderMock
	mainExpectation   *ProviderMockGetStorageExpectation
	expectationSeries []*ProviderMockGetStorageExpectation
}

type ProviderMockGetStorageExpectation struct {
	input  *ProviderMockGetStorageInput
	result *ProviderMockGetStorageResult
}

type ProviderMockGetStorageInput struct {
	p core.RecordID
}

type ProviderMockGetStorageResult struct {
	r RecentStorage
}

//Expect specifies that invocation of Provider.GetStorage is expected from 1 to Infinity times
func (m *mProviderMockGetStorage) Expect(p core.RecordID) *mProviderMockGetStorage {
	m.mock.GetStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetStorageExpectation{}
	}
	m.mainExpectation.input = &ProviderMockGetStorageInput{p}
	return m
}

//Return specifies results of invocation of Provider.GetStorage
func (m *mProviderMockGetStorage) Return(r RecentStorage) *ProviderMock {
	m.mock.GetStorageFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProviderMockGetStorageExpectation{}
	}
	m.mainExpectation.result = &ProviderMockGetStorageResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Provider.GetStorage is expected once
func (m *mProviderMockGetStorage) ExpectOnce(p core.RecordID) *ProviderMockGetStorageExpectation {
	m.mock.GetStorageFunc = nil
	m.mainExpectation = nil

	expectation := &ProviderMockGetStorageExpectation{}
	expectation.input = &ProviderMockGetStorageInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProviderMockGetStorageExpectation) Return(r RecentStorage) {
	e.result = &ProviderMockGetStorageResult{r}
}

//Set uses given function f as a mock of Provider.GetStorage method
func (m *mProviderMockGetStorage) Set(f func(p core.RecordID) (r RecentStorage)) *ProviderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStorageFunc = f
	return m.mock
}

//GetStorage implements github.com/insolar/insolar/ledger/recentstorage.Provider interface
func (m *ProviderMock) GetStorage(p core.RecordID) (r RecentStorage) {
	counter := atomic.AddUint64(&m.GetStoragePreCounter, 1)
	defer atomic.AddUint64(&m.GetStorageCounter, 1)

	if len(m.GetStorageMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStorageMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProviderMock.GetStorage. %v", p)
			return
		}

		input := m.GetStorageMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProviderMockGetStorageInput{p}, "Provider.GetStorage got unexpected parameters")

		result := m.GetStorageMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetStorage")
			return
		}

		r = result.r

		return
	}

	if m.GetStorageMock.mainExpectation != nil {

		input := m.GetStorageMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProviderMockGetStorageInput{p}, "Provider.GetStorage got unexpected parameters")
		}

		result := m.GetStorageMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProviderMock.GetStorage")
		}

		r = result.r

		return
	}

	if m.GetStorageFunc == nil {
		m.t.Fatalf("Unexpected call to ProviderMock.GetStorage. %v", p)
		return
	}

	return m.GetStorageFunc(p)
}

//GetStorageMinimockCounter returns a count of ProviderMock.GetStorageFunc invocations
func (m *ProviderMock) GetStorageMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStorageCounter)
}

//GetStorageMinimockPreCounter returns the value of ProviderMock.GetStorage invocations
func (m *ProviderMock) GetStorageMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStoragePreCounter)
}

//GetStorageFinished returns true if mock invocations count is ok
func (m *ProviderMock) GetStorageFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStorageMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStorageCounter) == uint64(len(m.GetStorageMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStorageMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStorageCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStorageFunc != nil {
		return atomic.LoadUint64(&m.GetStorageCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) ValidateCallCounters() {

	if !m.CloneStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.CloneStorage")
	}

	if !m.GetStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetStorage")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProviderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ProviderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ProviderMock) MinimockFinish() {

	if !m.CloneStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.CloneStorage")
	}

	if !m.GetStorageFinished() {
		m.t.Fatal("Expected call to ProviderMock.GetStorage")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ProviderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ProviderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CloneStorageFinished()
		ok = ok && m.GetStorageFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CloneStorageFinished() {
				m.t.Error("Expected call to ProviderMock.CloneStorage")
			}

			if !m.GetStorageFinished() {
				m.t.Error("Expected call to ProviderMock.GetStorage")
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
func (m *ProviderMock) AllMocksCalled() bool {

	if !m.CloneStorageFinished() {
		return false
	}

	if !m.GetStorageFinished() {
		return false
	}

	return true
}
