package jet

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Storage" can be found in github.com/insolar/insolar/insolar/jet
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//StorageMock implements github.com/insolar/insolar/insolar/jet.Storage
type StorageMock struct {
	t minimock.Tester

	AllFunc       func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)
	AllCounter    uint64
	AllPreCounter uint64
	AllMock       mStorageMockAll

	CloneFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber)
	CloneCounter    uint64
	ClonePreCounter uint64
	CloneMock       mStorageMockClone

	DeleteForPNFunc       func(p context.Context, p1 insolar.PulseNumber)
	DeleteForPNCounter    uint64
	DeleteForPNPreCounter uint64
	DeleteForPNMock       mStorageMockDeleteForPN

	ForIDFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r insolar.JetID, r1 bool)
	ForIDCounter    uint64
	ForIDPreCounter uint64
	ForIDMock       mStorageMockForID

	SplitFunc       func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r insolar.JetID, r1 insolar.JetID, r2 error)
	SplitCounter    uint64
	SplitPreCounter uint64
	SplitMock       mStorageMockSplit

	UpdateFunc       func(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID)
	UpdateCounter    uint64
	UpdatePreCounter uint64
	UpdateMock       mStorageMockUpdate
}

//NewStorageMock returns a mock for github.com/insolar/insolar/insolar/jet.Storage
func NewStorageMock(t minimock.Tester) *StorageMock {
	m := &StorageMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AllMock = mStorageMockAll{mock: m}
	m.CloneMock = mStorageMockClone{mock: m}
	m.DeleteForPNMock = mStorageMockDeleteForPN{mock: m}
	m.ForIDMock = mStorageMockForID{mock: m}
	m.SplitMock = mStorageMockSplit{mock: m}
	m.UpdateMock = mStorageMockUpdate{mock: m}

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
	p  context.Context
	p1 insolar.PulseNumber
}

type StorageMockAllResult struct {
	r []insolar.JetID
}

//Expect specifies that invocation of Storage.All is expected from 1 to Infinity times
func (m *mStorageMockAll) Expect(p context.Context, p1 insolar.PulseNumber) *mStorageMockAll {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockAllExpectation{}
	}
	m.mainExpectation.input = &StorageMockAllInput{p, p1}
	return m
}

//Return specifies results of invocation of Storage.All
func (m *mStorageMockAll) Return(r []insolar.JetID) *StorageMock {
	m.mock.AllFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockAllExpectation{}
	}
	m.mainExpectation.result = &StorageMockAllResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.All is expected once
func (m *mStorageMockAll) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *StorageMockAllExpectation {
	m.mock.AllFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockAllExpectation{}
	expectation.input = &StorageMockAllInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockAllExpectation) Return(r []insolar.JetID) {
	e.result = &StorageMockAllResult{r}
}

//Set uses given function f as a mock of Storage.All method
func (m *mStorageMockAll) Set(f func(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AllFunc = f
	return m.mock
}

//All implements github.com/insolar/insolar/insolar/jet.Storage interface
func (m *StorageMock) All(p context.Context, p1 insolar.PulseNumber) (r []insolar.JetID) {
	counter := atomic.AddUint64(&m.AllPreCounter, 1)
	defer atomic.AddUint64(&m.AllCounter, 1)

	if len(m.AllMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AllMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.All. %v %v", p, p1)
			return
		}

		input := m.AllMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockAllInput{p, p1}, "Storage.All got unexpected parameters")

		result := m.AllMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.All")
			return
		}

		r = result.r

		return
	}

	if m.AllMock.mainExpectation != nil {

		input := m.AllMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockAllInput{p, p1}, "Storage.All got unexpected parameters")
		}

		result := m.AllMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.All")
		}

		r = result.r

		return
	}

	if m.AllFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.All. %v %v", p, p1)
		return
	}

	return m.AllFunc(p, p1)
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

type mStorageMockClone struct {
	mock              *StorageMock
	mainExpectation   *StorageMockCloneExpectation
	expectationSeries []*StorageMockCloneExpectation
}

type StorageMockCloneExpectation struct {
	input *StorageMockCloneInput
}

type StorageMockCloneInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.PulseNumber
}

//Expect specifies that invocation of Storage.Clone is expected from 1 to Infinity times
func (m *mStorageMockClone) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) *mStorageMockClone {
	m.mock.CloneFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockCloneExpectation{}
	}
	m.mainExpectation.input = &StorageMockCloneInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Storage.Clone
func (m *mStorageMockClone) Return() *StorageMock {
	m.mock.CloneFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockCloneExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Storage.Clone is expected once
func (m *mStorageMockClone) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) *StorageMockCloneExpectation {
	m.mock.CloneFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockCloneExpectation{}
	expectation.input = &StorageMockCloneInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Storage.Clone method
func (m *mStorageMockClone) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CloneFunc = f
	return m.mock
}

//Clone implements github.com/insolar/insolar/insolar/jet.Storage interface
func (m *StorageMock) Clone(p context.Context, p1 insolar.PulseNumber, p2 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.ClonePreCounter, 1)
	defer atomic.AddUint64(&m.CloneCounter, 1)

	if len(m.CloneMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CloneMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.Clone. %v %v %v", p, p1, p2)
			return
		}

		input := m.CloneMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockCloneInput{p, p1, p2}, "Storage.Clone got unexpected parameters")

		return
	}

	if m.CloneMock.mainExpectation != nil {

		input := m.CloneMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockCloneInput{p, p1, p2}, "Storage.Clone got unexpected parameters")
		}

		return
	}

	if m.CloneFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.Clone. %v %v %v", p, p1, p2)
		return
	}

	m.CloneFunc(p, p1, p2)
}

//CloneMinimockCounter returns a count of StorageMock.CloneFunc invocations
func (m *StorageMock) CloneMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CloneCounter)
}

//CloneMinimockPreCounter returns the value of StorageMock.Clone invocations
func (m *StorageMock) CloneMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ClonePreCounter)
}

//CloneFinished returns true if mock invocations count is ok
func (m *StorageMock) CloneFinished() bool {
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

type mStorageMockDeleteForPN struct {
	mock              *StorageMock
	mainExpectation   *StorageMockDeleteForPNExpectation
	expectationSeries []*StorageMockDeleteForPNExpectation
}

type StorageMockDeleteForPNExpectation struct {
	input *StorageMockDeleteForPNInput
}

type StorageMockDeleteForPNInput struct {
	p  context.Context
	p1 insolar.PulseNumber
}

//Expect specifies that invocation of Storage.DeleteForPN is expected from 1 to Infinity times
func (m *mStorageMockDeleteForPN) Expect(p context.Context, p1 insolar.PulseNumber) *mStorageMockDeleteForPN {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockDeleteForPNExpectation{}
	}
	m.mainExpectation.input = &StorageMockDeleteForPNInput{p, p1}
	return m
}

//Return specifies results of invocation of Storage.DeleteForPN
func (m *mStorageMockDeleteForPN) Return() *StorageMock {
	m.mock.DeleteForPNFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockDeleteForPNExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Storage.DeleteForPN is expected once
func (m *mStorageMockDeleteForPN) ExpectOnce(p context.Context, p1 insolar.PulseNumber) *StorageMockDeleteForPNExpectation {
	m.mock.DeleteForPNFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockDeleteForPNExpectation{}
	expectation.input = &StorageMockDeleteForPNInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Storage.DeleteForPN method
func (m *mStorageMockDeleteForPN) Set(f func(p context.Context, p1 insolar.PulseNumber)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeleteForPNFunc = f
	return m.mock
}

//DeleteForPN implements github.com/insolar/insolar/insolar/jet.Storage interface
func (m *StorageMock) DeleteForPN(p context.Context, p1 insolar.PulseNumber) {
	counter := atomic.AddUint64(&m.DeleteForPNPreCounter, 1)
	defer atomic.AddUint64(&m.DeleteForPNCounter, 1)

	if len(m.DeleteForPNMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeleteForPNMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.DeleteForPN. %v %v", p, p1)
			return
		}

		input := m.DeleteForPNMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockDeleteForPNInput{p, p1}, "Storage.DeleteForPN got unexpected parameters")

		return
	}

	if m.DeleteForPNMock.mainExpectation != nil {

		input := m.DeleteForPNMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockDeleteForPNInput{p, p1}, "Storage.DeleteForPN got unexpected parameters")
		}

		return
	}

	if m.DeleteForPNFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.DeleteForPN. %v %v", p, p1)
		return
	}

	m.DeleteForPNFunc(p, p1)
}

//DeleteForPNMinimockCounter returns a count of StorageMock.DeleteForPNFunc invocations
func (m *StorageMock) DeleteForPNMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNCounter)
}

//DeleteForPNMinimockPreCounter returns the value of StorageMock.DeleteForPN invocations
func (m *StorageMock) DeleteForPNMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeleteForPNPreCounter)
}

//DeleteForPNFinished returns true if mock invocations count is ok
func (m *StorageMock) DeleteForPNFinished() bool {
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

type mStorageMockForID struct {
	mock              *StorageMock
	mainExpectation   *StorageMockForIDExpectation
	expectationSeries []*StorageMockForIDExpectation
}

type StorageMockForIDExpectation struct {
	input  *StorageMockForIDInput
	result *StorageMockForIDResult
}

type StorageMockForIDInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.ID
}

type StorageMockForIDResult struct {
	r  insolar.JetID
	r1 bool
}

//Expect specifies that invocation of Storage.ForID is expected from 1 to Infinity times
func (m *mStorageMockForID) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *mStorageMockForID {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockForIDExpectation{}
	}
	m.mainExpectation.input = &StorageMockForIDInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Storage.ForID
func (m *mStorageMockForID) Return(r insolar.JetID, r1 bool) *StorageMock {
	m.mock.ForIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockForIDExpectation{}
	}
	m.mainExpectation.result = &StorageMockForIDResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.ForID is expected once
func (m *mStorageMockForID) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) *StorageMockForIDExpectation {
	m.mock.ForIDFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockForIDExpectation{}
	expectation.input = &StorageMockForIDInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockForIDExpectation) Return(r insolar.JetID, r1 bool) {
	e.result = &StorageMockForIDResult{r, r1}
}

//Set uses given function f as a mock of Storage.ForID method
func (m *mStorageMockForID) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r insolar.JetID, r1 bool)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ForIDFunc = f
	return m.mock
}

//ForID implements github.com/insolar/insolar/insolar/jet.Storage interface
func (m *StorageMock) ForID(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID) (r insolar.JetID, r1 bool) {
	counter := atomic.AddUint64(&m.ForIDPreCounter, 1)
	defer atomic.AddUint64(&m.ForIDCounter, 1)

	if len(m.ForIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ForIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.ForID. %v %v %v", p, p1, p2)
			return
		}

		input := m.ForIDMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockForIDInput{p, p1, p2}, "Storage.ForID got unexpected parameters")

		result := m.ForIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.ForID")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDMock.mainExpectation != nil {

		input := m.ForIDMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockForIDInput{p, p1, p2}, "Storage.ForID got unexpected parameters")
		}

		result := m.ForIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.ForID")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ForIDFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.ForID. %v %v %v", p, p1, p2)
		return
	}

	return m.ForIDFunc(p, p1, p2)
}

//ForIDMinimockCounter returns a count of StorageMock.ForIDFunc invocations
func (m *StorageMock) ForIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDCounter)
}

//ForIDMinimockPreCounter returns the value of StorageMock.ForID invocations
func (m *StorageMock) ForIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ForIDPreCounter)
}

//ForIDFinished returns true if mock invocations count is ok
func (m *StorageMock) ForIDFinished() bool {
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

type mStorageMockSplit struct {
	mock              *StorageMock
	mainExpectation   *StorageMockSplitExpectation
	expectationSeries []*StorageMockSplitExpectation
}

type StorageMockSplitExpectation struct {
	input  *StorageMockSplitInput
	result *StorageMockSplitResult
}

type StorageMockSplitInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 insolar.JetID
}

type StorageMockSplitResult struct {
	r  insolar.JetID
	r1 insolar.JetID
	r2 error
}

//Expect specifies that invocation of Storage.Split is expected from 1 to Infinity times
func (m *mStorageMockSplit) Expect(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *mStorageMockSplit {
	m.mock.SplitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockSplitExpectation{}
	}
	m.mainExpectation.input = &StorageMockSplitInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Storage.Split
func (m *mStorageMockSplit) Return(r insolar.JetID, r1 insolar.JetID, r2 error) *StorageMock {
	m.mock.SplitFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockSplitExpectation{}
	}
	m.mainExpectation.result = &StorageMockSplitResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of Storage.Split is expected once
func (m *mStorageMockSplit) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) *StorageMockSplitExpectation {
	m.mock.SplitFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockSplitExpectation{}
	expectation.input = &StorageMockSplitInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *StorageMockSplitExpectation) Return(r insolar.JetID, r1 insolar.JetID, r2 error) {
	e.result = &StorageMockSplitResult{r, r1, r2}
}

//Set uses given function f as a mock of Storage.Split method
func (m *mStorageMockSplit) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r insolar.JetID, r1 insolar.JetID, r2 error)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SplitFunc = f
	return m.mock
}

//Split implements github.com/insolar/insolar/insolar/jet.Storage interface
func (m *StorageMock) Split(p context.Context, p1 insolar.PulseNumber, p2 insolar.JetID) (r insolar.JetID, r1 insolar.JetID, r2 error) {
	counter := atomic.AddUint64(&m.SplitPreCounter, 1)
	defer atomic.AddUint64(&m.SplitCounter, 1)

	if len(m.SplitMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SplitMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.Split. %v %v %v", p, p1, p2)
			return
		}

		input := m.SplitMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockSplitInput{p, p1, p2}, "Storage.Split got unexpected parameters")

		result := m.SplitMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.Split")
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
			testify_assert.Equal(m.t, *input, StorageMockSplitInput{p, p1, p2}, "Storage.Split got unexpected parameters")
		}

		result := m.SplitMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the StorageMock.Split")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.SplitFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.Split. %v %v %v", p, p1, p2)
		return
	}

	return m.SplitFunc(p, p1, p2)
}

//SplitMinimockCounter returns a count of StorageMock.SplitFunc invocations
func (m *StorageMock) SplitMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SplitCounter)
}

//SplitMinimockPreCounter returns the value of StorageMock.Split invocations
func (m *StorageMock) SplitMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SplitPreCounter)
}

//SplitFinished returns true if mock invocations count is ok
func (m *StorageMock) SplitFinished() bool {
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

type mStorageMockUpdate struct {
	mock              *StorageMock
	mainExpectation   *StorageMockUpdateExpectation
	expectationSeries []*StorageMockUpdateExpectation
}

type StorageMockUpdateExpectation struct {
	input *StorageMockUpdateInput
}

type StorageMockUpdateInput struct {
	p  context.Context
	p1 insolar.PulseNumber
	p2 bool
	p3 []insolar.JetID
}

//Expect specifies that invocation of Storage.Update is expected from 1 to Infinity times
func (m *mStorageMockUpdate) Expect(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID) *mStorageMockUpdate {
	m.mock.UpdateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockUpdateExpectation{}
	}
	m.mainExpectation.input = &StorageMockUpdateInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Storage.Update
func (m *mStorageMockUpdate) Return() *StorageMock {
	m.mock.UpdateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &StorageMockUpdateExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of Storage.Update is expected once
func (m *mStorageMockUpdate) ExpectOnce(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID) *StorageMockUpdateExpectation {
	m.mock.UpdateFunc = nil
	m.mainExpectation = nil

	expectation := &StorageMockUpdateExpectation{}
	expectation.input = &StorageMockUpdateInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of Storage.Update method
func (m *mStorageMockUpdate) Set(f func(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID)) *StorageMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateFunc = f
	return m.mock
}

//Update implements github.com/insolar/insolar/insolar/jet.Storage interface
func (m *StorageMock) Update(p context.Context, p1 insolar.PulseNumber, p2 bool, p3 ...insolar.JetID) {
	counter := atomic.AddUint64(&m.UpdatePreCounter, 1)
	defer atomic.AddUint64(&m.UpdateCounter, 1)

	if len(m.UpdateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to StorageMock.Update. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.UpdateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, StorageMockUpdateInput{p, p1, p2, p3}, "Storage.Update got unexpected parameters")

		return
	}

	if m.UpdateMock.mainExpectation != nil {

		input := m.UpdateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, StorageMockUpdateInput{p, p1, p2, p3}, "Storage.Update got unexpected parameters")
		}

		return
	}

	if m.UpdateFunc == nil {
		m.t.Fatalf("Unexpected call to StorageMock.Update. %v %v %v %v", p, p1, p2, p3)
		return
	}

	m.UpdateFunc(p, p1, p2, p3...)
}

//UpdateMinimockCounter returns a count of StorageMock.UpdateFunc invocations
func (m *StorageMock) UpdateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateCounter)
}

//UpdateMinimockPreCounter returns the value of StorageMock.Update invocations
func (m *StorageMock) UpdateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdatePreCounter)
}

//UpdateFinished returns true if mock invocations count is ok
func (m *StorageMock) UpdateFinished() bool {
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
func (m *StorageMock) ValidateCallCounters() {

	if !m.AllFinished() {
		m.t.Fatal("Expected call to StorageMock.All")
	}

	if !m.CloneFinished() {
		m.t.Fatal("Expected call to StorageMock.Clone")
	}

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to StorageMock.DeleteForPN")
	}

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to StorageMock.ForID")
	}

	if !m.SplitFinished() {
		m.t.Fatal("Expected call to StorageMock.Split")
	}

	if !m.UpdateFinished() {
		m.t.Fatal("Expected call to StorageMock.Update")
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

	if !m.CloneFinished() {
		m.t.Fatal("Expected call to StorageMock.Clone")
	}

	if !m.DeleteForPNFinished() {
		m.t.Fatal("Expected call to StorageMock.DeleteForPN")
	}

	if !m.ForIDFinished() {
		m.t.Fatal("Expected call to StorageMock.ForID")
	}

	if !m.SplitFinished() {
		m.t.Fatal("Expected call to StorageMock.Split")
	}

	if !m.UpdateFinished() {
		m.t.Fatal("Expected call to StorageMock.Update")
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
		ok = ok && m.CloneFinished()
		ok = ok && m.DeleteForPNFinished()
		ok = ok && m.ForIDFinished()
		ok = ok && m.SplitFinished()
		ok = ok && m.UpdateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AllFinished() {
				m.t.Error("Expected call to StorageMock.All")
			}

			if !m.CloneFinished() {
				m.t.Error("Expected call to StorageMock.Clone")
			}

			if !m.DeleteForPNFinished() {
				m.t.Error("Expected call to StorageMock.DeleteForPN")
			}

			if !m.ForIDFinished() {
				m.t.Error("Expected call to StorageMock.ForID")
			}

			if !m.SplitFinished() {
				m.t.Error("Expected call to StorageMock.Split")
			}

			if !m.UpdateFinished() {
				m.t.Error("Expected call to StorageMock.Update")
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

	if !m.CloneFinished() {
		return false
	}

	if !m.DeleteForPNFinished() {
		return false
	}

	if !m.ForIDFinished() {
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
