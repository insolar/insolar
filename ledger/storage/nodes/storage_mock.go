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

package nodes

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Storage" can be found in github.com/insolar/insolar/ledger/storage/nodes
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//StorageMock implements github.com/insolar/insolar/ledger/storage/nodes.Storage
type StorageMock struct {
	t minimock.Tester

	AllFunc       func(p core.PulseNumber) (r []core.Node, r1 error)
	AllCounter    uint64
	AllPreCounter uint64
	AllMock       mStorageMockAll

	InRoleFunc       func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error)
	InRoleCounter    uint64
	InRolePreCounter uint64
	InRoleMock       mStorageMockInRole

	RemoveActiveNodesUntilFunc       func(p core.PulseNumber)
	RemoveActiveNodesUntilCounter    uint64
	RemoveActiveNodesUntilPreCounter uint64
	RemoveActiveNodesUntilMock       mStorageMockRemoveActiveNodesUntil

	SetFunc       func(p core.PulseNumber, p1 []core.Node) (r error)
	SetCounter    uint64
	SetPreCounter uint64
	SetMock       mStorageMockSet
}

//NewStorageMock returns a mock for github.com/insolar/insolar/ledger/storage/nodes.Storage
func NewStorageMock(t minimock.Tester) *StorageMock {
	m := &StorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AllMock = mStorageMockAll{mock: m}
	m.InRoleMock = mStorageMockInRole{mock: m}
	m.RemoveActiveNodesUntilMock = mStorageMockRemoveActiveNodesUntil{mock: m}
	m.SetMock = mStorageMockSet{mock: m}

	return m
}

type mStorageMockAll struct {
	mock              *StorageMock
	mainExpectation   *StorageMockAllExpectation
	expectationSeries []*StorageMockAllExpectation
}

type StorageMockAllExpectation struct {
	input  *StorageMockAllInput
	result *StorageMockAllResult
}

type StorageMockAllInput struct {
	p core.PulseNumber
}

type StorageMockAllResult struct {
	r  []core.Node
	r1 error
}

//Expect specifies that invocation of Storage.All is expected from 1 to Infinity times
func (m *mStorageMockAll) Expect(p core.PulseNumber) *mStorageMockAll {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockAllExpectation{}
	}
	m.mainExpectation.input = &StorageMockAllInput{p}
	return m
}

//Return specifies results of invocation of Storage.All
func (m *mStorageMockAll) Return(r []core.Node, r1 error) *StorageMock {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockAllExpectation{}
	}
	m.mainExpectation.result = &StorageMockAllResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.All is expected once
func (m *mStorageMockAll) ExpectOnce(p core.PulseNumber) *StorageMockAllExpectation {
	m.mock.AllFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockAllExpectation{}
	expectation.input = &StorageMockAllInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockAllExpectation) Return(r []core.Node, r1 error) {
	e.result = &StorageMockAllResult{r, r1}
}

//Set uses given function f as a mock of Storage.All method
func (m *mStorageMockAll) Set(f func(p core.PulseNumber) (r []core.Node, r1 error)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AllFunc = f
	return m.mock
}

//All implements github.com/insolar/insolar/ledger/storage/nodes.Storage interface
func (m *StorageMock) All(p core.PulseNumber) (r []core.Node, r1 error) {
	counter := atomic.AddUint64(&m.AllPreCounter, 1)
	defer atomic.AddUint64(&m.AllCounter, 1)

	if len(m.AllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.All. %v", p)
			return
		}

		input := m.AllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockAllInput{p}, "Storage.All got unexpected parameters")

		result := m.AllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.All")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AllMock.mainExpectation != nil {

		input := m.AllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockAllInput{p}, "Storage.All got unexpected parameters")
		}

		result := m.AllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.All")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.AllFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.All. %v", p)
		return
	}

	return m.AllFunc(p)
}

//AllMinimockCounter returns a count of StorageMock.AllFunc invocations
func (m *StorageMock) AllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AllCounter)
}

//AllMinimockPreCounter returns the value of StorageMock.All invocations
func (m *StorageMock) AllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AllPreCounter)
}

//AllFinished returns true if mock invocations count is ok
func (m *StorageMock) AllFinished() bool {
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

type mStorageMockInRole struct {
	mock              *StorageMock
	mainExpectation   *StorageMockInRoleExpectation
	expectationSeries []*StorageMockInRoleExpectation
}

type StorageMockInRoleExpectation struct {
	input  *StorageMockInRoleInput
	result *StorageMockInRoleResult
}

type StorageMockInRoleInput struct {
	p  core.PulseNumber
	p1 core.StaticRole
}

type StorageMockInRoleResult struct {
	r  []core.Node
	r1 error
}

//Expect specifies that invocation of Storage.InRole is expected from 1 to Infinity times
func (m *mStorageMockInRole) Expect(p core.PulseNumber, p1 core.StaticRole) *mStorageMockInRole {
	m.mock.InRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockInRoleExpectation{}
	}
	m.mainExpectation.input = &StorageMockInRoleInput{p, p1}
	return m
}

//Return specifies results of invocation of Storage.InRole
func (m *mStorageMockInRole) Return(r []core.Node, r1 error) *StorageMock {
	m.mock.InRoleFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockInRoleExpectation{}
	}
	m.mainExpectation.result = &StorageMockInRoleResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.InRole is expected once
func (m *mStorageMockInRole) ExpectOnce(p core.PulseNumber, p1 core.StaticRole) *StorageMockInRoleExpectation {
	m.mock.InRoleFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockInRoleExpectation{}
	expectation.input = &StorageMockInRoleInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockInRoleExpectation) Return(r []core.Node, r1 error) {
	e.result = &StorageMockInRoleResult{r, r1}
}

//Set uses given function f as a mock of Storage.InRole method
func (m *mStorageMockInRole) Set(f func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.InRoleFunc = f
	return m.mock
}

//InRole implements github.com/insolar/insolar/ledger/storage/nodes.Storage interface
func (m *StorageMock) InRole(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error) {
	counter := atomic.AddUint64(&m.InRolePreCounter, 1)
	defer atomic.AddUint64(&m.InRoleCounter, 1)

	if len(m.InRoleMock.expectationSeries) > 0 {
		if counter > uint64(len(m.InRoleMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.InRole. %v %v", p, p1)
			return
		}

		input := m.InRoleMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockInRoleInput{p, p1}, "Storage.InRole got unexpected parameters")

		result := m.InRoleMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.InRole")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.InRoleMock.mainExpectation != nil {

		input := m.InRoleMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockInRoleInput{p, p1}, "Storage.InRole got unexpected parameters")
		}

		result := m.InRoleMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.InRole")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.InRoleFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.InRole. %v %v", p, p1)
		return
	}

	return m.InRoleFunc(p, p1)
}

//InRoleMinimockCounter returns a count of StorageMock.InRoleFunc invocations
func (m *StorageMock) InRoleMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.InRoleCounter)
}

//InRoleMinimockPreCounter returns the value of StorageMock.InRole invocations
func (m *StorageMock) InRoleMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.InRolePreCounter)
}

//InRoleFinished returns true if mock invocations count is ok
func (m *StorageMock) InRoleFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.InRoleMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.InRoleCounter) == uint64(len(m.InRoleMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.InRoleMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.InRoleCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.InRoleFunc != nil {
		return atomic.LoadUint64(&m.InRoleCounter) > 0
	}

	return true
}

type mStorageMockRemoveActiveNodesUntil struct {
	mock              *StorageMock
	mainExpectation   *StorageMockRemoveActiveNodesUntilExpectation
	expectationSeries []*StorageMockRemoveActiveNodesUntilExpectation
}

type StorageMockRemoveActiveNodesUntilExpectation struct {
	input *StorageMockRemoveActiveNodesUntilInput
}

type StorageMockRemoveActiveNodesUntilInput struct {
	p core.PulseNumber
}

//Expect specifies that invocation of Storage.RemoveActiveNodesUntil is expected from 1 to Infinity times
func (m *mStorageMockRemoveActiveNodesUntil) Expect(p core.PulseNumber) *mStorageMockRemoveActiveNodesUntil {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockRemoveActiveNodesUntilExpectation{}
	}
	m.mainExpectation.input = &StorageMockRemoveActiveNodesUntilInput{p}
	return m
}

//Return specifies results of invocation of Storage.RemoveActiveNodesUntil
func (m *mStorageMockRemoveActiveNodesUntil) Return() *StorageMock {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockRemoveActiveNodesUntilExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Storage.RemoveActiveNodesUntil is expected once
func (m *mStorageMockRemoveActiveNodesUntil) ExpectOnce(p core.PulseNumber) *StorageMockRemoveActiveNodesUntilExpectation {
	m.mock.RemoveActiveNodesUntilFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockRemoveActiveNodesUntilExpectation{}
	expectation.input = &StorageMockRemoveActiveNodesUntilInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Storage.RemoveActiveNodesUntil method
func (m *mStorageMockRemoveActiveNodesUntil) Set(f func(p core.PulseNumber)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RemoveActiveNodesUntilFunc = f
	return m.mock
}

//RemoveActiveNodesUntil implements github.com/insolar/insolar/ledger/storage/nodes.Storage interface
func (m *StorageMock) RemoveActiveNodesUntil(p core.PulseNumber) {
	counter := atomic.AddUint64(&m.RemoveActiveNodesUntilPreCounter, 1)
	defer atomic.AddUint64(&m.RemoveActiveNodesUntilCounter, 1)

	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.RemoveActiveNodesUntil. %v", p)
			return
		}

		input := m.RemoveActiveNodesUntilMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockRemoveActiveNodesUntilInput{p}, "Storage.RemoveActiveNodesUntil got unexpected parameters")

		return
	}

	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {

		input := m.RemoveActiveNodesUntilMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockRemoveActiveNodesUntilInput{p}, "Storage.RemoveActiveNodesUntil got unexpected parameters")
		}

		return
	}

	if m.RemoveActiveNodesUntilFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.RemoveActiveNodesUntil. %v", p)
		return
	}

	m.RemoveActiveNodesUntilFunc(p)
}

//RemoveActiveNodesUntilMinimockCounter returns a count of StorageMock.RemoveActiveNodesUntilFunc invocations
func (m *StorageMock) RemoveActiveNodesUntilMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter)
}

//RemoveActiveNodesUntilMinimockPreCounter returns the value of StorageMock.RemoveActiveNodesUntil invocations
func (m *StorageMock) RemoveActiveNodesUntilMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RemoveActiveNodesUntilPreCounter)
}

//RemoveActiveNodesUntilFinished returns true if mock invocations count is ok
func (m *StorageMock) RemoveActiveNodesUntilFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RemoveActiveNodesUntilMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) == uint64(len(m.RemoveActiveNodesUntilMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RemoveActiveNodesUntilFunc != nil {
		return atomic.LoadUint64(&m.RemoveActiveNodesUntilCounter) > 0
	}

	return true
}

type mStorageMockSet struct {
	mock              *StorageMock
	mainExpectation   *StorageMockSetExpectation
	expectationSeries []*StorageMockSetExpectation
}

type StorageMockSetExpectation struct {
	input  *StorageMockSetInput
	result *StorageMockSetResult
}

type StorageMockSetInput struct {
	p  core.PulseNumber
	p1 []core.Node
}

type StorageMockSetResult struct {
	r error
}

//Expect specifies that invocation of Storage.Set is expected from 1 to Infinity times
func (m *mStorageMockSet) Expect(p core.PulseNumber, p1 []core.Node) *mStorageMockSet {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockSetExpectation{}
	}
	m.mainExpectation.input = &StorageMockSetInput{p, p1}
	return m
}

//Return specifies results of invocation of Storage.Set
func (m *mStorageMockSet) Return(r error) *StorageMock {
	m.mock.SetFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockSetExpectation{}
	}
	m.mainExpectation.result = &StorageMockSetResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.Set is expected once
func (m *mStorageMockSet) ExpectOnce(p core.PulseNumber, p1 []core.Node) *StorageMockSetExpectation {
	m.mock.SetFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockSetExpectation{}
	expectation.input = &StorageMockSetInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockSetExpectation) Return(r error) {
	e.result = &StorageMockSetResult{r}
}

//Set uses given function f as a mock of Storage.Set method
func (m *mStorageMockSet) Set(f func(p core.PulseNumber, p1 []core.Node) (r error)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SetFunc = f
	return m.mock
}

//Set implements github.com/insolar/insolar/ledger/storage/nodes.Storage interface
func (m *StorageMock) Set(p core.PulseNumber, p1 []core.Node) (r error) {
	counter := atomic.AddUint64(&m.SetPreCounter, 1)
	defer atomic.AddUint64(&m.SetCounter, 1)

	if len(m.SetMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SetMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.Set. %v %v", p, p1)
			return
		}

		input := m.SetMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockSetInput{p, p1}, "Storage.Set got unexpected parameters")

		result := m.SetMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.Set")
			return
		}

		r = result.r

		return
	}

	if m.SetMock.mainExpectation != nil {

		input := m.SetMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockSetInput{p, p1}, "Storage.Set got unexpected parameters")
		}

		result := m.SetMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.Set")
		}

		r = result.r

		return
	}

	if m.SetFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.Set. %v %v", p, p1)
		return
	}

	return m.SetFunc(p, p1)
}

//SetMinimockCounter returns a count of StorageMock.SetFunc invocations
func (m *StorageMock) SetMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SetCounter)
}

//SetMinimockPreCounter returns the value of StorageMock.Set invocations
func (m *StorageMock) SetMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SetPreCounter)
}

//SetFinished returns true if mock invocations count is ok
func (m *StorageMock) SetFinished() bool {
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
func (m *StorageMock) ValidateCallCounters() {

	if !m.AllFinished() {
		m.t.Fatal("Expected call to StorageMock.All")
	}

	if !m.InRoleFinished() {
		m.t.Fatal("Expected call to StorageMock.InRole")
	}

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to StorageMock.RemoveActiveNodesUntil")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to StorageMock.Set")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *StorageMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *StorageMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *StorageMock) MinimockFinish() {

	if !m.AllFinished() {
		m.t.Fatal("Expected call to StorageMock.All")
	}

	if !m.InRoleFinished() {
		m.t.Fatal("Expected call to StorageMock.InRole")
	}

	if !m.RemoveActiveNodesUntilFinished() {
		m.t.Fatal("Expected call to StorageMock.RemoveActiveNodesUntil")
	}

	if !m.SetFinished() {
		m.t.Fatal("Expected call to StorageMock.Set")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *StorageMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *StorageMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AllFinished()
		ok = ok && m.InRoleFinished()
		ok = ok && m.RemoveActiveNodesUntilFinished()
		ok = ok && m.SetFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AllFinished() {
				m.t.Error("Expected call to StorageMock.All")
			}

			if !m.InRoleFinished() {
				m.t.Error("Expected call to StorageMock.InRole")
			}

			if !m.RemoveActiveNodesUntilFinished() {
				m.t.Error("Expected call to StorageMock.RemoveActiveNodesUntil")
			}

			if !m.SetFinished() {
				m.t.Error("Expected call to StorageMock.Set")
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
func (m *StorageMock) AllMocksCalled() bool {

	if !m.AllFinished() {
		return false
	}

	if !m.InRoleFinished() {
		return false
	}

	if !m.RemoveActiveNodesUntilFinished() {
		return false
	}

	if !m.SetFinished() {
		return false
	}

	return true
}
