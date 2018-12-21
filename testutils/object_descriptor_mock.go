package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ObjectDescriptor" can be found in github.com/insolar/insolar/core
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//ObjectDescriptorMock implements github.com/insolar/insolar/core.ObjectDescriptor
type ObjectDescriptorMock struct {
	t minimock.Tester

	ChildPointerFunc       func() (r *core.RecordID)
	ChildPointerCounter    uint64
	ChildPointerPreCounter uint64
	ChildPointerMock       mObjectDescriptorMockChildPointer

	ChildrenFunc       func(p *core.PulseNumber) (r core.RefIterator, r1 error)
	ChildrenCounter    uint64
	ChildrenPreCounter uint64
	ChildrenMock       mObjectDescriptorMockChildren

	CodeFunc       func() (r *core.RecordRef, r1 error)
	CodeCounter    uint64
	CodePreCounter uint64
	CodeMock       mObjectDescriptorMockCode

	HeadRefFunc       func() (r *core.RecordRef)
	HeadRefCounter    uint64
	HeadRefPreCounter uint64
	HeadRefMock       mObjectDescriptorMockHeadRef

	IsPrototypeFunc       func() (r bool)
	IsPrototypeCounter    uint64
	IsPrototypePreCounter uint64
	IsPrototypeMock       mObjectDescriptorMockIsPrototype

	MemoryFunc       func() (r []byte)
	MemoryCounter    uint64
	MemoryPreCounter uint64
	MemoryMock       mObjectDescriptorMockMemory

	ParentFunc       func() (r *core.RecordRef)
	ParentCounter    uint64
	ParentPreCounter uint64
	ParentMock       mObjectDescriptorMockParent

	PrototypeFunc       func() (r *core.RecordRef, r1 error)
	PrototypeCounter    uint64
	PrototypePreCounter uint64
	PrototypeMock       mObjectDescriptorMockPrototype

	StateIDFunc       func() (r *core.RecordID)
	StateIDCounter    uint64
	StateIDPreCounter uint64
	StateIDMock       mObjectDescriptorMockStateID
}

//NewObjectDescriptorMock returns a mock for github.com/insolar/insolar/core.ObjectDescriptor
func NewObjectDescriptorMock(t minimock.Tester) *ObjectDescriptorMock {
	m := &ObjectDescriptorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ChildPointerMock = mObjectDescriptorMockChildPointer{mock: m}
	m.ChildrenMock = mObjectDescriptorMockChildren{mock: m}
	m.CodeMock = mObjectDescriptorMockCode{mock: m}
	m.HeadRefMock = mObjectDescriptorMockHeadRef{mock: m}
	m.IsPrototypeMock = mObjectDescriptorMockIsPrototype{mock: m}
	m.MemoryMock = mObjectDescriptorMockMemory{mock: m}
	m.ParentMock = mObjectDescriptorMockParent{mock: m}
	m.PrototypeMock = mObjectDescriptorMockPrototype{mock: m}
	m.StateIDMock = mObjectDescriptorMockStateID{mock: m}

	return m
}

type mObjectDescriptorMockChildPointer struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockChildPointerExpectation
	expectationSeries []*ObjectDescriptorMockChildPointerExpectation
}

type ObjectDescriptorMockChildPointerExpectation struct {
	result *ObjectDescriptorMockChildPointerResult
}

type ObjectDescriptorMockChildPointerResult struct {
	r *core.RecordID
}

//Expect specifies that invocation of ObjectDescriptor.ChildPointer is expected from 1 to Infinity times
func (m *mObjectDescriptorMockChildPointer) Expect() *mObjectDescriptorMockChildPointer {
	m.mock.ChildPointerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockChildPointerExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.ChildPointer
func (m *mObjectDescriptorMockChildPointer) Return(r *core.RecordID) *ObjectDescriptorMock {
	m.mock.ChildPointerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockChildPointerExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockChildPointerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.ChildPointer is expected once
func (m *mObjectDescriptorMockChildPointer) ExpectOnce() *ObjectDescriptorMockChildPointerExpectation {
	m.mock.ChildPointerFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockChildPointerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockChildPointerExpectation) Return(r *core.RecordID) {
	e.result = &ObjectDescriptorMockChildPointerResult{r}
}

//Set uses given function f as a mock of ObjectDescriptor.ChildPointer method
func (m *mObjectDescriptorMockChildPointer) Set(f func() (r *core.RecordID)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ChildPointerFunc = f
	return m.mock
}

//ChildPointer implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) ChildPointer() (r *core.RecordID) {
	counter := atomic.AddUint64(&m.ChildPointerPreCounter, 1)
	defer atomic.AddUint64(&m.ChildPointerCounter, 1)

	if len(m.ChildPointerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ChildPointerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.ChildPointer.")
			return
		}

		result := m.ChildPointerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.ChildPointer")
			return
		}

		r = result.r

		return
	}

	if m.ChildPointerMock.mainExpectation != nil {

		result := m.ChildPointerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.ChildPointer")
		}

		r = result.r

		return
	}

	if m.ChildPointerFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.ChildPointer.")
		return
	}

	return m.ChildPointerFunc()
}

//ChildPointerMinimockCounter returns a count of ObjectDescriptorMock.ChildPointerFunc invocations
func (m *ObjectDescriptorMock) ChildPointerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ChildPointerCounter)
}

//ChildPointerMinimockPreCounter returns the value of ObjectDescriptorMock.ChildPointer invocations
func (m *ObjectDescriptorMock) ChildPointerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ChildPointerPreCounter)
}

//ChildPointerFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) ChildPointerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ChildPointerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ChildPointerCounter) == uint64(len(m.ChildPointerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ChildPointerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ChildPointerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ChildPointerFunc != nil {
		return atomic.LoadUint64(&m.ChildPointerCounter) > 0
	}

	return true
}

type mObjectDescriptorMockChildren struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockChildrenExpectation
	expectationSeries []*ObjectDescriptorMockChildrenExpectation
}

type ObjectDescriptorMockChildrenExpectation struct {
	input  *ObjectDescriptorMockChildrenInput
	result *ObjectDescriptorMockChildrenResult
}

type ObjectDescriptorMockChildrenInput struct {
	p *core.PulseNumber
}

type ObjectDescriptorMockChildrenResult struct {
	r  core.RefIterator
	r1 error
}

//Expect specifies that invocation of ObjectDescriptor.Children is expected from 1 to Infinity times
func (m *mObjectDescriptorMockChildren) Expect(p *core.PulseNumber) *mObjectDescriptorMockChildren {
	m.mock.ChildrenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockChildrenExpectation{}
	}
	m.mainExpectation.input = &ObjectDescriptorMockChildrenInput{p}
	return m
}

//Return specifies results of invocation of ObjectDescriptor.Children
func (m *mObjectDescriptorMockChildren) Return(r core.RefIterator, r1 error) *ObjectDescriptorMock {
	m.mock.ChildrenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockChildrenExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockChildrenResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.Children is expected once
func (m *mObjectDescriptorMockChildren) ExpectOnce(p *core.PulseNumber) *ObjectDescriptorMockChildrenExpectation {
	m.mock.ChildrenFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockChildrenExpectation{}
	expectation.input = &ObjectDescriptorMockChildrenInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockChildrenExpectation) Return(r core.RefIterator, r1 error) {
	e.result = &ObjectDescriptorMockChildrenResult{r, r1}
}

//Set uses given function f as a mock of ObjectDescriptor.Children method
func (m *mObjectDescriptorMockChildren) Set(f func(p *core.PulseNumber) (r core.RefIterator, r1 error)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ChildrenFunc = f
	return m.mock
}

//Children implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) Children(p *core.PulseNumber) (r core.RefIterator, r1 error) {
	counter := atomic.AddUint64(&m.ChildrenPreCounter, 1)
	defer atomic.AddUint64(&m.ChildrenCounter, 1)

	if len(m.ChildrenMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ChildrenMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Children. %v", p)
			return
		}

		input := m.ChildrenMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ObjectDescriptorMockChildrenInput{p}, "ObjectDescriptor.Children got unexpected parameters")

		result := m.ChildrenMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Children")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ChildrenMock.mainExpectation != nil {

		input := m.ChildrenMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ObjectDescriptorMockChildrenInput{p}, "ObjectDescriptor.Children got unexpected parameters")
		}

		result := m.ChildrenMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Children")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ChildrenFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Children. %v", p)
		return
	}

	return m.ChildrenFunc(p)
}

//ChildrenMinimockCounter returns a count of ObjectDescriptorMock.ChildrenFunc invocations
func (m *ObjectDescriptorMock) ChildrenMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ChildrenCounter)
}

//ChildrenMinimockPreCounter returns the value of ObjectDescriptorMock.Children invocations
func (m *ObjectDescriptorMock) ChildrenMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ChildrenPreCounter)
}

//ChildrenFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) ChildrenFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ChildrenMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ChildrenCounter) == uint64(len(m.ChildrenMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ChildrenMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ChildrenCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ChildrenFunc != nil {
		return atomic.LoadUint64(&m.ChildrenCounter) > 0
	}

	return true
}

type mObjectDescriptorMockCode struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockCodeExpectation
	expectationSeries []*ObjectDescriptorMockCodeExpectation
}

type ObjectDescriptorMockCodeExpectation struct {
	result *ObjectDescriptorMockCodeResult
}

type ObjectDescriptorMockCodeResult struct {
	r  *core.RecordRef
	r1 error
}

//Expect specifies that invocation of ObjectDescriptor.Code is expected from 1 to Infinity times
func (m *mObjectDescriptorMockCode) Expect() *mObjectDescriptorMockCode {
	m.mock.CodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockCodeExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.Code
func (m *mObjectDescriptorMockCode) Return(r *core.RecordRef, r1 error) *ObjectDescriptorMock {
	m.mock.CodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockCodeExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.Code is expected once
func (m *mObjectDescriptorMockCode) ExpectOnce() *ObjectDescriptorMockCodeExpectation {
	m.mock.CodeFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockCodeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockCodeExpectation) Return(r *core.RecordRef, r1 error) {
	e.result = &ObjectDescriptorMockCodeResult{r, r1}
}

//Set uses given function f as a mock of ObjectDescriptor.Code method
func (m *mObjectDescriptorMockCode) Set(f func() (r *core.RecordRef, r1 error)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CodeFunc = f
	return m.mock
}

//Code implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) Code() (r *core.RecordRef, r1 error) {
	counter := atomic.AddUint64(&m.CodePreCounter, 1)
	defer atomic.AddUint64(&m.CodeCounter, 1)

	if len(m.CodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Code.")
			return
		}

		result := m.CodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Code")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CodeMock.mainExpectation != nil {

		result := m.CodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Code")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CodeFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Code.")
		return
	}

	return m.CodeFunc()
}

//CodeMinimockCounter returns a count of ObjectDescriptorMock.CodeFunc invocations
func (m *ObjectDescriptorMock) CodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CodeCounter)
}

//CodeMinimockPreCounter returns the value of ObjectDescriptorMock.Code invocations
func (m *ObjectDescriptorMock) CodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CodePreCounter)
}

//CodeFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) CodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CodeCounter) == uint64(len(m.CodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CodeFunc != nil {
		return atomic.LoadUint64(&m.CodeCounter) > 0
	}

	return true
}

type mObjectDescriptorMockHeadRef struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockHeadRefExpectation
	expectationSeries []*ObjectDescriptorMockHeadRefExpectation
}

type ObjectDescriptorMockHeadRefExpectation struct {
	result *ObjectDescriptorMockHeadRefResult
}

type ObjectDescriptorMockHeadRefResult struct {
	r *core.RecordRef
}

//Expect specifies that invocation of ObjectDescriptor.HeadRef is expected from 1 to Infinity times
func (m *mObjectDescriptorMockHeadRef) Expect() *mObjectDescriptorMockHeadRef {
	m.mock.HeadRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockHeadRefExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.HeadRef
func (m *mObjectDescriptorMockHeadRef) Return(r *core.RecordRef) *ObjectDescriptorMock {
	m.mock.HeadRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockHeadRefExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockHeadRefResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.HeadRef is expected once
func (m *mObjectDescriptorMockHeadRef) ExpectOnce() *ObjectDescriptorMockHeadRefExpectation {
	m.mock.HeadRefFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockHeadRefExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockHeadRefExpectation) Return(r *core.RecordRef) {
	e.result = &ObjectDescriptorMockHeadRefResult{r}
}

//Set uses given function f as a mock of ObjectDescriptor.HeadRef method
func (m *mObjectDescriptorMockHeadRef) Set(f func() (r *core.RecordRef)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.HeadRefFunc = f
	return m.mock
}

//HeadRef implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) HeadRef() (r *core.RecordRef) {
	counter := atomic.AddUint64(&m.HeadRefPreCounter, 1)
	defer atomic.AddUint64(&m.HeadRefCounter, 1)

	if len(m.HeadRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.HeadRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.HeadRef.")
			return
		}

		result := m.HeadRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.HeadRef")
			return
		}

		r = result.r

		return
	}

	if m.HeadRefMock.mainExpectation != nil {

		result := m.HeadRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.HeadRef")
		}

		r = result.r

		return
	}

	if m.HeadRefFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.HeadRef.")
		return
	}

	return m.HeadRefFunc()
}

//HeadRefMinimockCounter returns a count of ObjectDescriptorMock.HeadRefFunc invocations
func (m *ObjectDescriptorMock) HeadRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.HeadRefCounter)
}

//HeadRefMinimockPreCounter returns the value of ObjectDescriptorMock.HeadRef invocations
func (m *ObjectDescriptorMock) HeadRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.HeadRefPreCounter)
}

//HeadRefFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) HeadRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.HeadRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.HeadRefCounter) == uint64(len(m.HeadRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.HeadRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.HeadRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.HeadRefFunc != nil {
		return atomic.LoadUint64(&m.HeadRefCounter) > 0
	}

	return true
}

type mObjectDescriptorMockIsPrototype struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockIsPrototypeExpectation
	expectationSeries []*ObjectDescriptorMockIsPrototypeExpectation
}

type ObjectDescriptorMockIsPrototypeExpectation struct {
	result *ObjectDescriptorMockIsPrototypeResult
}

type ObjectDescriptorMockIsPrototypeResult struct {
	r bool
}

//Expect specifies that invocation of ObjectDescriptor.IsPrototype is expected from 1 to Infinity times
func (m *mObjectDescriptorMockIsPrototype) Expect() *mObjectDescriptorMockIsPrototype {
	m.mock.IsPrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockIsPrototypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.IsPrototype
func (m *mObjectDescriptorMockIsPrototype) Return(r bool) *ObjectDescriptorMock {
	m.mock.IsPrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockIsPrototypeExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockIsPrototypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.IsPrototype is expected once
func (m *mObjectDescriptorMockIsPrototype) ExpectOnce() *ObjectDescriptorMockIsPrototypeExpectation {
	m.mock.IsPrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockIsPrototypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockIsPrototypeExpectation) Return(r bool) {
	e.result = &ObjectDescriptorMockIsPrototypeResult{r}
}

//Set uses given function f as a mock of ObjectDescriptor.IsPrototype method
func (m *mObjectDescriptorMockIsPrototype) Set(f func() (r bool)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsPrototypeFunc = f
	return m.mock
}

//IsPrototype implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) IsPrototype() (r bool) {
	counter := atomic.AddUint64(&m.IsPrototypePreCounter, 1)
	defer atomic.AddUint64(&m.IsPrototypeCounter, 1)

	if len(m.IsPrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsPrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.IsPrototype.")
			return
		}

		result := m.IsPrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.IsPrototype")
			return
		}

		r = result.r

		return
	}

	if m.IsPrototypeMock.mainExpectation != nil {

		result := m.IsPrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.IsPrototype")
		}

		r = result.r

		return
	}

	if m.IsPrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.IsPrototype.")
		return
	}

	return m.IsPrototypeFunc()
}

//IsPrototypeMinimockCounter returns a count of ObjectDescriptorMock.IsPrototypeFunc invocations
func (m *ObjectDescriptorMock) IsPrototypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsPrototypeCounter)
}

//IsPrototypeMinimockPreCounter returns the value of ObjectDescriptorMock.IsPrototype invocations
func (m *ObjectDescriptorMock) IsPrototypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsPrototypePreCounter)
}

//IsPrototypeFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) IsPrototypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsPrototypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsPrototypeCounter) == uint64(len(m.IsPrototypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsPrototypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsPrototypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsPrototypeFunc != nil {
		return atomic.LoadUint64(&m.IsPrototypeCounter) > 0
	}

	return true
}

type mObjectDescriptorMockMemory struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockMemoryExpectation
	expectationSeries []*ObjectDescriptorMockMemoryExpectation
}

type ObjectDescriptorMockMemoryExpectation struct {
	result *ObjectDescriptorMockMemoryResult
}

type ObjectDescriptorMockMemoryResult struct {
	r []byte
}

//Expect specifies that invocation of ObjectDescriptor.Memory is expected from 1 to Infinity times
func (m *mObjectDescriptorMockMemory) Expect() *mObjectDescriptorMockMemory {
	m.mock.MemoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockMemoryExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.Memory
func (m *mObjectDescriptorMockMemory) Return(r []byte) *ObjectDescriptorMock {
	m.mock.MemoryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockMemoryExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockMemoryResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.Memory is expected once
func (m *mObjectDescriptorMockMemory) ExpectOnce() *ObjectDescriptorMockMemoryExpectation {
	m.mock.MemoryFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockMemoryExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockMemoryExpectation) Return(r []byte) {
	e.result = &ObjectDescriptorMockMemoryResult{r}
}

//Set uses given function f as a mock of ObjectDescriptor.Memory method
func (m *mObjectDescriptorMockMemory) Set(f func() (r []byte)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MemoryFunc = f
	return m.mock
}

//Memory implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) Memory() (r []byte) {
	counter := atomic.AddUint64(&m.MemoryPreCounter, 1)
	defer atomic.AddUint64(&m.MemoryCounter, 1)

	if len(m.MemoryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MemoryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Memory.")
			return
		}

		result := m.MemoryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Memory")
			return
		}

		r = result.r

		return
	}

	if m.MemoryMock.mainExpectation != nil {

		result := m.MemoryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Memory")
		}

		r = result.r

		return
	}

	if m.MemoryFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Memory.")
		return
	}

	return m.MemoryFunc()
}

//MemoryMinimockCounter returns a count of ObjectDescriptorMock.MemoryFunc invocations
func (m *ObjectDescriptorMock) MemoryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MemoryCounter)
}

//MemoryMinimockPreCounter returns the value of ObjectDescriptorMock.Memory invocations
func (m *ObjectDescriptorMock) MemoryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MemoryPreCounter)
}

//MemoryFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) MemoryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MemoryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MemoryCounter) == uint64(len(m.MemoryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MemoryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MemoryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MemoryFunc != nil {
		return atomic.LoadUint64(&m.MemoryCounter) > 0
	}

	return true
}

type mObjectDescriptorMockParent struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockParentExpectation
	expectationSeries []*ObjectDescriptorMockParentExpectation
}

type ObjectDescriptorMockParentExpectation struct {
	result *ObjectDescriptorMockParentResult
}

type ObjectDescriptorMockParentResult struct {
	r *core.RecordRef
}

//Expect specifies that invocation of ObjectDescriptor.Parent is expected from 1 to Infinity times
func (m *mObjectDescriptorMockParent) Expect() *mObjectDescriptorMockParent {
	m.mock.ParentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockParentExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.Parent
func (m *mObjectDescriptorMockParent) Return(r *core.RecordRef) *ObjectDescriptorMock {
	m.mock.ParentFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockParentExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockParentResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.Parent is expected once
func (m *mObjectDescriptorMockParent) ExpectOnce() *ObjectDescriptorMockParentExpectation {
	m.mock.ParentFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockParentExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockParentExpectation) Return(r *core.RecordRef) {
	e.result = &ObjectDescriptorMockParentResult{r}
}

//Set uses given function f as a mock of ObjectDescriptor.Parent method
func (m *mObjectDescriptorMockParent) Set(f func() (r *core.RecordRef)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ParentFunc = f
	return m.mock
}

//Parent implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) Parent() (r *core.RecordRef) {
	counter := atomic.AddUint64(&m.ParentPreCounter, 1)
	defer atomic.AddUint64(&m.ParentCounter, 1)

	if len(m.ParentMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ParentMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Parent.")
			return
		}

		result := m.ParentMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Parent")
			return
		}

		r = result.r

		return
	}

	if m.ParentMock.mainExpectation != nil {

		result := m.ParentMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Parent")
		}

		r = result.r

		return
	}

	if m.ParentFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Parent.")
		return
	}

	return m.ParentFunc()
}

//ParentMinimockCounter returns a count of ObjectDescriptorMock.ParentFunc invocations
func (m *ObjectDescriptorMock) ParentMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ParentCounter)
}

//ParentMinimockPreCounter returns the value of ObjectDescriptorMock.Parent invocations
func (m *ObjectDescriptorMock) ParentMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ParentPreCounter)
}

//ParentFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) ParentFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ParentMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ParentCounter) == uint64(len(m.ParentMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ParentMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ParentCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ParentFunc != nil {
		return atomic.LoadUint64(&m.ParentCounter) > 0
	}

	return true
}

type mObjectDescriptorMockPrototype struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockPrototypeExpectation
	expectationSeries []*ObjectDescriptorMockPrototypeExpectation
}

type ObjectDescriptorMockPrototypeExpectation struct {
	result *ObjectDescriptorMockPrototypeResult
}

type ObjectDescriptorMockPrototypeResult struct {
	r  *core.RecordRef
	r1 error
}

//Expect specifies that invocation of ObjectDescriptor.Prototype is expected from 1 to Infinity times
func (m *mObjectDescriptorMockPrototype) Expect() *mObjectDescriptorMockPrototype {
	m.mock.PrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockPrototypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.Prototype
func (m *mObjectDescriptorMockPrototype) Return(r *core.RecordRef, r1 error) *ObjectDescriptorMock {
	m.mock.PrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockPrototypeExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockPrototypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.Prototype is expected once
func (m *mObjectDescriptorMockPrototype) ExpectOnce() *ObjectDescriptorMockPrototypeExpectation {
	m.mock.PrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockPrototypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockPrototypeExpectation) Return(r *core.RecordRef, r1 error) {
	e.result = &ObjectDescriptorMockPrototypeResult{r, r1}
}

//Set uses given function f as a mock of ObjectDescriptor.Prototype method
func (m *mObjectDescriptorMockPrototype) Set(f func() (r *core.RecordRef, r1 error)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PrototypeFunc = f
	return m.mock
}

//Prototype implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) Prototype() (r *core.RecordRef, r1 error) {
	counter := atomic.AddUint64(&m.PrototypePreCounter, 1)
	defer atomic.AddUint64(&m.PrototypeCounter, 1)

	if len(m.PrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Prototype.")
			return
		}

		result := m.PrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Prototype")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PrototypeMock.mainExpectation != nil {

		result := m.PrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.Prototype")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.Prototype.")
		return
	}

	return m.PrototypeFunc()
}

//PrototypeMinimockCounter returns a count of ObjectDescriptorMock.PrototypeFunc invocations
func (m *ObjectDescriptorMock) PrototypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PrototypeCounter)
}

//PrototypeMinimockPreCounter returns the value of ObjectDescriptorMock.Prototype invocations
func (m *ObjectDescriptorMock) PrototypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PrototypePreCounter)
}

//PrototypeFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) PrototypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PrototypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PrototypeCounter) == uint64(len(m.PrototypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PrototypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PrototypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PrototypeFunc != nil {
		return atomic.LoadUint64(&m.PrototypeCounter) > 0
	}

	return true
}

type mObjectDescriptorMockStateID struct {
	mock              *ObjectDescriptorMock
	mainExpectation   *ObjectDescriptorMockStateIDExpectation
	expectationSeries []*ObjectDescriptorMockStateIDExpectation
}

type ObjectDescriptorMockStateIDExpectation struct {
	result *ObjectDescriptorMockStateIDResult
}

type ObjectDescriptorMockStateIDResult struct {
	r *core.RecordID
}

//Expect specifies that invocation of ObjectDescriptor.StateID is expected from 1 to Infinity times
func (m *mObjectDescriptorMockStateID) Expect() *mObjectDescriptorMockStateID {
	m.mock.StateIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockStateIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of ObjectDescriptor.StateID
func (m *mObjectDescriptorMockStateID) Return(r *core.RecordID) *ObjectDescriptorMock {
	m.mock.StateIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ObjectDescriptorMockStateIDExpectation{}
	}
	m.mainExpectation.result = &ObjectDescriptorMockStateIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ObjectDescriptor.StateID is expected once
func (m *mObjectDescriptorMockStateID) ExpectOnce() *ObjectDescriptorMockStateIDExpectation {
	m.mock.StateIDFunc = nil
	m.mainExpectation = nil

	expectation := &ObjectDescriptorMockStateIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ObjectDescriptorMockStateIDExpectation) Return(r *core.RecordID) {
	e.result = &ObjectDescriptorMockStateIDResult{r}
}

//Set uses given function f as a mock of ObjectDescriptor.StateID method
func (m *mObjectDescriptorMockStateID) Set(f func() (r *core.RecordID)) *ObjectDescriptorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.StateIDFunc = f
	return m.mock
}

//StateID implements github.com/insolar/insolar/core.ObjectDescriptor interface
func (m *ObjectDescriptorMock) StateID() (r *core.RecordID) {
	counter := atomic.AddUint64(&m.StateIDPreCounter, 1)
	defer atomic.AddUint64(&m.StateIDCounter, 1)

	if len(m.StateIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.StateIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ObjectDescriptorMock.StateID.")
			return
		}

		result := m.StateIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.StateID")
			return
		}

		r = result.r

		return
	}

	if m.StateIDMock.mainExpectation != nil {

		result := m.StateIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ObjectDescriptorMock.StateID")
		}

		r = result.r

		return
	}

	if m.StateIDFunc == nil {
		m.t.Fatalf("Unexpected call to ObjectDescriptorMock.StateID.")
		return
	}

	return m.StateIDFunc()
}

//StateIDMinimockCounter returns a count of ObjectDescriptorMock.StateIDFunc invocations
func (m *ObjectDescriptorMock) StateIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.StateIDCounter)
}

//StateIDMinimockPreCounter returns the value of ObjectDescriptorMock.StateID invocations
func (m *ObjectDescriptorMock) StateIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.StateIDPreCounter)
}

//StateIDFinished returns true if mock invocations count is ok
func (m *ObjectDescriptorMock) StateIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.StateIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.StateIDCounter) == uint64(len(m.StateIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.StateIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.StateIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.StateIDFunc != nil {
		return atomic.LoadUint64(&m.StateIDCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectDescriptorMock) ValidateCallCounters() {

	if !m.ChildPointerFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.ChildPointer")
	}

	if !m.ChildrenFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Children")
	}

	if !m.CodeFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Code")
	}

	if !m.HeadRefFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.HeadRef")
	}

	if !m.IsPrototypeFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.IsPrototype")
	}

	if !m.MemoryFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Memory")
	}

	if !m.ParentFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Parent")
	}

	if !m.PrototypeFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Prototype")
	}

	if !m.StateIDFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.StateID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ObjectDescriptorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ObjectDescriptorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ObjectDescriptorMock) MinimockFinish() {

	if !m.ChildPointerFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.ChildPointer")
	}

	if !m.ChildrenFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Children")
	}

	if !m.CodeFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Code")
	}

	if !m.HeadRefFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.HeadRef")
	}

	if !m.IsPrototypeFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.IsPrototype")
	}

	if !m.MemoryFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Memory")
	}

	if !m.ParentFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Parent")
	}

	if !m.PrototypeFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.Prototype")
	}

	if !m.StateIDFinished() {
		m.t.Fatal("Expected call to ObjectDescriptorMock.StateID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ObjectDescriptorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ObjectDescriptorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ChildPointerFinished()
		ok = ok && m.ChildrenFinished()
		ok = ok && m.CodeFinished()
		ok = ok && m.HeadRefFinished()
		ok = ok && m.IsPrototypeFinished()
		ok = ok && m.MemoryFinished()
		ok = ok && m.ParentFinished()
		ok = ok && m.PrototypeFinished()
		ok = ok && m.StateIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ChildPointerFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.ChildPointer")
			}

			if !m.ChildrenFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.Children")
			}

			if !m.CodeFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.Code")
			}

			if !m.HeadRefFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.HeadRef")
			}

			if !m.IsPrototypeFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.IsPrototype")
			}

			if !m.MemoryFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.Memory")
			}

			if !m.ParentFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.Parent")
			}

			if !m.PrototypeFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.Prototype")
			}

			if !m.StateIDFinished() {
				m.t.Error("Expected call to ObjectDescriptorMock.StateID")
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
func (m *ObjectDescriptorMock) AllMocksCalled() bool {

	if !m.ChildPointerFinished() {
		return false
	}

	if !m.ChildrenFinished() {
		return false
	}

	if !m.CodeFinished() {
		return false
	}

	if !m.HeadRefFinished() {
		return false
	}

	if !m.IsPrototypeFinished() {
		return false
	}

	if !m.MemoryFinished() {
		return false
	}

	if !m.ParentFinished() {
		return false
	}

	if !m.PrototypeFinished() {
		return false
	}

	if !m.StateIDFinished() {
		return false
	}

	return true
}
