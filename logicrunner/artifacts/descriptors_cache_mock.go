package artifacts

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DescriptorsCache" can be found in github.com/insolar/insolar/logicrunner/artifacts
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//DescriptorsCacheMock implements github.com/insolar/insolar/logicrunner/artifacts.DescriptorsCache
type DescriptorsCacheMock struct {
	t minimock.Tester

	ByObjectDescriptorFunc       func(p context.Context, p1 ObjectDescriptor) (r ObjectDescriptor, r1 CodeDescriptor, r2 error)
	ByObjectDescriptorCounter    uint64
	ByObjectDescriptorPreCounter uint64
	ByObjectDescriptorMock       mDescriptorsCacheMockByObjectDescriptor

	ByPrototypeRefFunc       func(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 CodeDescriptor, r2 error)
	ByPrototypeRefCounter    uint64
	ByPrototypeRefPreCounter uint64
	ByPrototypeRefMock       mDescriptorsCacheMockByPrototypeRef

	GetCodeFunc       func(p context.Context, p1 insolar.Reference) (r CodeDescriptor, r1 error)
	GetCodeCounter    uint64
	GetCodePreCounter uint64
	GetCodeMock       mDescriptorsCacheMockGetCode

	GetPrototypeFunc       func(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 error)
	GetPrototypeCounter    uint64
	GetPrototypePreCounter uint64
	GetPrototypeMock       mDescriptorsCacheMockGetPrototype
}

//NewDescriptorsCacheMock returns a mock for github.com/insolar/insolar/logicrunner/artifacts.DescriptorsCache
func NewDescriptorsCacheMock(t minimock.Tester) *DescriptorsCacheMock {
	m := &DescriptorsCacheMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ByObjectDescriptorMock = mDescriptorsCacheMockByObjectDescriptor{mock: m}
	m.ByPrototypeRefMock = mDescriptorsCacheMockByPrototypeRef{mock: m}
	m.GetCodeMock = mDescriptorsCacheMockGetCode{mock: m}
	m.GetPrototypeMock = mDescriptorsCacheMockGetPrototype{mock: m}

	return m
}

type mDescriptorsCacheMockByObjectDescriptor struct {
	mock              *DescriptorsCacheMock
	mainExpectation   *DescriptorsCacheMockByObjectDescriptorExpectation
	expectationSeries []*DescriptorsCacheMockByObjectDescriptorExpectation
}

type DescriptorsCacheMockByObjectDescriptorExpectation struct {
	input  *DescriptorsCacheMockByObjectDescriptorInput
	result *DescriptorsCacheMockByObjectDescriptorResult
}

type DescriptorsCacheMockByObjectDescriptorInput struct {
	p  context.Context
	p1 ObjectDescriptor
}

type DescriptorsCacheMockByObjectDescriptorResult struct {
	r  ObjectDescriptor
	r1 CodeDescriptor
	r2 error
}

//Expect specifies that invocation of DescriptorsCache.ByObjectDescriptor is expected from 1 to Infinity times
func (m *mDescriptorsCacheMockByObjectDescriptor) Expect(p context.Context, p1 ObjectDescriptor) *mDescriptorsCacheMockByObjectDescriptor {
	m.mock.ByObjectDescriptorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockByObjectDescriptorExpectation{}
	}
	m.mainExpectation.input = &DescriptorsCacheMockByObjectDescriptorInput{p, p1}
	return m
}

//Return specifies results of invocation of DescriptorsCache.ByObjectDescriptor
func (m *mDescriptorsCacheMockByObjectDescriptor) Return(r ObjectDescriptor, r1 CodeDescriptor, r2 error) *DescriptorsCacheMock {
	m.mock.ByObjectDescriptorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockByObjectDescriptorExpectation{}
	}
	m.mainExpectation.result = &DescriptorsCacheMockByObjectDescriptorResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of DescriptorsCache.ByObjectDescriptor is expected once
func (m *mDescriptorsCacheMockByObjectDescriptor) ExpectOnce(p context.Context, p1 ObjectDescriptor) *DescriptorsCacheMockByObjectDescriptorExpectation {
	m.mock.ByObjectDescriptorFunc = nil
	m.mainExpectation = nil

	expectation := &DescriptorsCacheMockByObjectDescriptorExpectation{}
	expectation.input = &DescriptorsCacheMockByObjectDescriptorInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DescriptorsCacheMockByObjectDescriptorExpectation) Return(r ObjectDescriptor, r1 CodeDescriptor, r2 error) {
	e.result = &DescriptorsCacheMockByObjectDescriptorResult{r, r1, r2}
}

//Set uses given function f as a mock of DescriptorsCache.ByObjectDescriptor method
func (m *mDescriptorsCacheMockByObjectDescriptor) Set(f func(p context.Context, p1 ObjectDescriptor) (r ObjectDescriptor, r1 CodeDescriptor, r2 error)) *DescriptorsCacheMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ByObjectDescriptorFunc = f
	return m.mock
}

//ByObjectDescriptor implements github.com/insolar/insolar/logicrunner/artifacts.DescriptorsCache interface
func (m *DescriptorsCacheMock) ByObjectDescriptor(p context.Context, p1 ObjectDescriptor) (r ObjectDescriptor, r1 CodeDescriptor, r2 error) {
	counter := atomic.AddUint64(&m.ByObjectDescriptorPreCounter, 1)
	defer atomic.AddUint64(&m.ByObjectDescriptorCounter, 1)

	if len(m.ByObjectDescriptorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ByObjectDescriptorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DescriptorsCacheMock.ByObjectDescriptor. %v %v", p, p1)
			return
		}

		input := m.ByObjectDescriptorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DescriptorsCacheMockByObjectDescriptorInput{p, p1}, "DescriptorsCache.ByObjectDescriptor got unexpected parameters")

		result := m.ByObjectDescriptorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.ByObjectDescriptor")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.ByObjectDescriptorMock.mainExpectation != nil {

		input := m.ByObjectDescriptorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DescriptorsCacheMockByObjectDescriptorInput{p, p1}, "DescriptorsCache.ByObjectDescriptor got unexpected parameters")
		}

		result := m.ByObjectDescriptorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.ByObjectDescriptor")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.ByObjectDescriptorFunc == nil {
		m.t.Fatalf("Unexpected call to DescriptorsCacheMock.ByObjectDescriptor. %v %v", p, p1)
		return
	}

	return m.ByObjectDescriptorFunc(p, p1)
}

//ByObjectDescriptorMinimockCounter returns a count of DescriptorsCacheMock.ByObjectDescriptorFunc invocations
func (m *DescriptorsCacheMock) ByObjectDescriptorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ByObjectDescriptorCounter)
}

//ByObjectDescriptorMinimockPreCounter returns the value of DescriptorsCacheMock.ByObjectDescriptor invocations
func (m *DescriptorsCacheMock) ByObjectDescriptorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ByObjectDescriptorPreCounter)
}

//ByObjectDescriptorFinished returns true if mock invocations count is ok
func (m *DescriptorsCacheMock) ByObjectDescriptorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ByObjectDescriptorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ByObjectDescriptorCounter) == uint64(len(m.ByObjectDescriptorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ByObjectDescriptorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ByObjectDescriptorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ByObjectDescriptorFunc != nil {
		return atomic.LoadUint64(&m.ByObjectDescriptorCounter) > 0
	}

	return true
}

type mDescriptorsCacheMockByPrototypeRef struct {
	mock              *DescriptorsCacheMock
	mainExpectation   *DescriptorsCacheMockByPrototypeRefExpectation
	expectationSeries []*DescriptorsCacheMockByPrototypeRefExpectation
}

type DescriptorsCacheMockByPrototypeRefExpectation struct {
	input  *DescriptorsCacheMockByPrototypeRefInput
	result *DescriptorsCacheMockByPrototypeRefResult
}

type DescriptorsCacheMockByPrototypeRefInput struct {
	p  context.Context
	p1 insolar.Reference
}

type DescriptorsCacheMockByPrototypeRefResult struct {
	r  ObjectDescriptor
	r1 CodeDescriptor
	r2 error
}

//Expect specifies that invocation of DescriptorsCache.ByPrototypeRef is expected from 1 to Infinity times
func (m *mDescriptorsCacheMockByPrototypeRef) Expect(p context.Context, p1 insolar.Reference) *mDescriptorsCacheMockByPrototypeRef {
	m.mock.ByPrototypeRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockByPrototypeRefExpectation{}
	}
	m.mainExpectation.input = &DescriptorsCacheMockByPrototypeRefInput{p, p1}
	return m
}

//Return specifies results of invocation of DescriptorsCache.ByPrototypeRef
func (m *mDescriptorsCacheMockByPrototypeRef) Return(r ObjectDescriptor, r1 CodeDescriptor, r2 error) *DescriptorsCacheMock {
	m.mock.ByPrototypeRefFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockByPrototypeRefExpectation{}
	}
	m.mainExpectation.result = &DescriptorsCacheMockByPrototypeRefResult{r, r1, r2}
	return m.mock
}

//ExpectOnce specifies that invocation of DescriptorsCache.ByPrototypeRef is expected once
func (m *mDescriptorsCacheMockByPrototypeRef) ExpectOnce(p context.Context, p1 insolar.Reference) *DescriptorsCacheMockByPrototypeRefExpectation {
	m.mock.ByPrototypeRefFunc = nil
	m.mainExpectation = nil

	expectation := &DescriptorsCacheMockByPrototypeRefExpectation{}
	expectation.input = &DescriptorsCacheMockByPrototypeRefInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DescriptorsCacheMockByPrototypeRefExpectation) Return(r ObjectDescriptor, r1 CodeDescriptor, r2 error) {
	e.result = &DescriptorsCacheMockByPrototypeRefResult{r, r1, r2}
}

//Set uses given function f as a mock of DescriptorsCache.ByPrototypeRef method
func (m *mDescriptorsCacheMockByPrototypeRef) Set(f func(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 CodeDescriptor, r2 error)) *DescriptorsCacheMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ByPrototypeRefFunc = f
	return m.mock
}

//ByPrototypeRef implements github.com/insolar/insolar/logicrunner/artifacts.DescriptorsCache interface
func (m *DescriptorsCacheMock) ByPrototypeRef(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 CodeDescriptor, r2 error) {
	counter := atomic.AddUint64(&m.ByPrototypeRefPreCounter, 1)
	defer atomic.AddUint64(&m.ByPrototypeRefCounter, 1)

	if len(m.ByPrototypeRefMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ByPrototypeRefMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DescriptorsCacheMock.ByPrototypeRef. %v %v", p, p1)
			return
		}

		input := m.ByPrototypeRefMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DescriptorsCacheMockByPrototypeRefInput{p, p1}, "DescriptorsCache.ByPrototypeRef got unexpected parameters")

		result := m.ByPrototypeRefMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.ByPrototypeRef")
			return
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.ByPrototypeRefMock.mainExpectation != nil {

		input := m.ByPrototypeRefMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DescriptorsCacheMockByPrototypeRefInput{p, p1}, "DescriptorsCache.ByPrototypeRef got unexpected parameters")
		}

		result := m.ByPrototypeRefMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.ByPrototypeRef")
		}

		r = result.r
		r1 = result.r1
		r2 = result.r2

		return
	}

	if m.ByPrototypeRefFunc == nil {
		m.t.Fatalf("Unexpected call to DescriptorsCacheMock.ByPrototypeRef. %v %v", p, p1)
		return
	}

	return m.ByPrototypeRefFunc(p, p1)
}

//ByPrototypeRefMinimockCounter returns a count of DescriptorsCacheMock.ByPrototypeRefFunc invocations
func (m *DescriptorsCacheMock) ByPrototypeRefMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ByPrototypeRefCounter)
}

//ByPrototypeRefMinimockPreCounter returns the value of DescriptorsCacheMock.ByPrototypeRef invocations
func (m *DescriptorsCacheMock) ByPrototypeRefMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ByPrototypeRefPreCounter)
}

//ByPrototypeRefFinished returns true if mock invocations count is ok
func (m *DescriptorsCacheMock) ByPrototypeRefFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ByPrototypeRefMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ByPrototypeRefCounter) == uint64(len(m.ByPrototypeRefMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ByPrototypeRefMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ByPrototypeRefCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ByPrototypeRefFunc != nil {
		return atomic.LoadUint64(&m.ByPrototypeRefCounter) > 0
	}

	return true
}

type mDescriptorsCacheMockGetCode struct {
	mock              *DescriptorsCacheMock
	mainExpectation   *DescriptorsCacheMockGetCodeExpectation
	expectationSeries []*DescriptorsCacheMockGetCodeExpectation
}

type DescriptorsCacheMockGetCodeExpectation struct {
	input  *DescriptorsCacheMockGetCodeInput
	result *DescriptorsCacheMockGetCodeResult
}

type DescriptorsCacheMockGetCodeInput struct {
	p  context.Context
	p1 insolar.Reference
}

type DescriptorsCacheMockGetCodeResult struct {
	r  CodeDescriptor
	r1 error
}

//Expect specifies that invocation of DescriptorsCache.GetCode is expected from 1 to Infinity times
func (m *mDescriptorsCacheMockGetCode) Expect(p context.Context, p1 insolar.Reference) *mDescriptorsCacheMockGetCode {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockGetCodeExpectation{}
	}
	m.mainExpectation.input = &DescriptorsCacheMockGetCodeInput{p, p1}
	return m
}

//Return specifies results of invocation of DescriptorsCache.GetCode
func (m *mDescriptorsCacheMockGetCode) Return(r CodeDescriptor, r1 error) *DescriptorsCacheMock {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockGetCodeExpectation{}
	}
	m.mainExpectation.result = &DescriptorsCacheMockGetCodeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DescriptorsCache.GetCode is expected once
func (m *mDescriptorsCacheMockGetCode) ExpectOnce(p context.Context, p1 insolar.Reference) *DescriptorsCacheMockGetCodeExpectation {
	m.mock.GetCodeFunc = nil
	m.mainExpectation = nil

	expectation := &DescriptorsCacheMockGetCodeExpectation{}
	expectation.input = &DescriptorsCacheMockGetCodeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DescriptorsCacheMockGetCodeExpectation) Return(r CodeDescriptor, r1 error) {
	e.result = &DescriptorsCacheMockGetCodeResult{r, r1}
}

//Set uses given function f as a mock of DescriptorsCache.GetCode method
func (m *mDescriptorsCacheMockGetCode) Set(f func(p context.Context, p1 insolar.Reference) (r CodeDescriptor, r1 error)) *DescriptorsCacheMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCodeFunc = f
	return m.mock
}

//GetCode implements github.com/insolar/insolar/logicrunner/artifacts.DescriptorsCache interface
func (m *DescriptorsCacheMock) GetCode(p context.Context, p1 insolar.Reference) (r CodeDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.GetCodePreCounter, 1)
	defer atomic.AddUint64(&m.GetCodeCounter, 1)

	if len(m.GetCodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DescriptorsCacheMock.GetCode. %v %v", p, p1)
			return
		}

		input := m.GetCodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DescriptorsCacheMockGetCodeInput{p, p1}, "DescriptorsCache.GetCode got unexpected parameters")

		result := m.GetCodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.GetCode")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeMock.mainExpectation != nil {

		input := m.GetCodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DescriptorsCacheMockGetCodeInput{p, p1}, "DescriptorsCache.GetCode got unexpected parameters")
		}

		result := m.GetCodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.GetCode")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetCodeFunc == nil {
		m.t.Fatalf("Unexpected call to DescriptorsCacheMock.GetCode. %v %v", p, p1)
		return
	}

	return m.GetCodeFunc(p, p1)
}

//GetCodeMinimockCounter returns a count of DescriptorsCacheMock.GetCodeFunc invocations
func (m *DescriptorsCacheMock) GetCodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodeCounter)
}

//GetCodeMinimockPreCounter returns the value of DescriptorsCacheMock.GetCode invocations
func (m *DescriptorsCacheMock) GetCodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodePreCounter)
}

//GetCodeFinished returns true if mock invocations count is ok
func (m *DescriptorsCacheMock) GetCodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetCodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetCodeCounter) == uint64(len(m.GetCodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetCodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetCodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetCodeFunc != nil {
		return atomic.LoadUint64(&m.GetCodeCounter) > 0
	}

	return true
}

type mDescriptorsCacheMockGetPrototype struct {
	mock              *DescriptorsCacheMock
	mainExpectation   *DescriptorsCacheMockGetPrototypeExpectation
	expectationSeries []*DescriptorsCacheMockGetPrototypeExpectation
}

type DescriptorsCacheMockGetPrototypeExpectation struct {
	input  *DescriptorsCacheMockGetPrototypeInput
	result *DescriptorsCacheMockGetPrototypeResult
}

type DescriptorsCacheMockGetPrototypeInput struct {
	p  context.Context
	p1 insolar.Reference
}

type DescriptorsCacheMockGetPrototypeResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of DescriptorsCache.GetPrototype is expected from 1 to Infinity times
func (m *mDescriptorsCacheMockGetPrototype) Expect(p context.Context, p1 insolar.Reference) *mDescriptorsCacheMockGetPrototype {
	m.mock.GetPrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockGetPrototypeExpectation{}
	}
	m.mainExpectation.input = &DescriptorsCacheMockGetPrototypeInput{p, p1}
	return m
}

//Return specifies results of invocation of DescriptorsCache.GetPrototype
func (m *mDescriptorsCacheMockGetPrototype) Return(r ObjectDescriptor, r1 error) *DescriptorsCacheMock {
	m.mock.GetPrototypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DescriptorsCacheMockGetPrototypeExpectation{}
	}
	m.mainExpectation.result = &DescriptorsCacheMockGetPrototypeResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of DescriptorsCache.GetPrototype is expected once
func (m *mDescriptorsCacheMockGetPrototype) ExpectOnce(p context.Context, p1 insolar.Reference) *DescriptorsCacheMockGetPrototypeExpectation {
	m.mock.GetPrototypeFunc = nil
	m.mainExpectation = nil

	expectation := &DescriptorsCacheMockGetPrototypeExpectation{}
	expectation.input = &DescriptorsCacheMockGetPrototypeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DescriptorsCacheMockGetPrototypeExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &DescriptorsCacheMockGetPrototypeResult{r, r1}
}

//Set uses given function f as a mock of DescriptorsCache.GetPrototype method
func (m *mDescriptorsCacheMockGetPrototype) Set(f func(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 error)) *DescriptorsCacheMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPrototypeFunc = f
	return m.mock
}

//GetPrototype implements github.com/insolar/insolar/logicrunner/artifacts.DescriptorsCache interface
func (m *DescriptorsCacheMock) GetPrototype(p context.Context, p1 insolar.Reference) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.GetPrototypePreCounter, 1)
	defer atomic.AddUint64(&m.GetPrototypeCounter, 1)

	if len(m.GetPrototypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPrototypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DescriptorsCacheMock.GetPrototype. %v %v", p, p1)
			return
		}

		input := m.GetPrototypeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DescriptorsCacheMockGetPrototypeInput{p, p1}, "DescriptorsCache.GetPrototype got unexpected parameters")

		result := m.GetPrototypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.GetPrototype")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPrototypeMock.mainExpectation != nil {

		input := m.GetPrototypeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DescriptorsCacheMockGetPrototypeInput{p, p1}, "DescriptorsCache.GetPrototype got unexpected parameters")
		}

		result := m.GetPrototypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DescriptorsCacheMock.GetPrototype")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GetPrototypeFunc == nil {
		m.t.Fatalf("Unexpected call to DescriptorsCacheMock.GetPrototype. %v %v", p, p1)
		return
	}

	return m.GetPrototypeFunc(p, p1)
}

//GetPrototypeMinimockCounter returns a count of DescriptorsCacheMock.GetPrototypeFunc invocations
func (m *DescriptorsCacheMock) GetPrototypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrototypeCounter)
}

//GetPrototypeMinimockPreCounter returns the value of DescriptorsCacheMock.GetPrototype invocations
func (m *DescriptorsCacheMock) GetPrototypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPrototypePreCounter)
}

//GetPrototypeFinished returns true if mock invocations count is ok
func (m *DescriptorsCacheMock) GetPrototypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPrototypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPrototypeCounter) == uint64(len(m.GetPrototypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPrototypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPrototypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPrototypeFunc != nil {
		return atomic.LoadUint64(&m.GetPrototypeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DescriptorsCacheMock) ValidateCallCounters() {

	if !m.ByObjectDescriptorFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.ByObjectDescriptor")
	}

	if !m.ByPrototypeRefFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.ByPrototypeRef")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.GetCode")
	}

	if !m.GetPrototypeFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.GetPrototype")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DescriptorsCacheMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DescriptorsCacheMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DescriptorsCacheMock) MinimockFinish() {

	if !m.ByObjectDescriptorFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.ByObjectDescriptor")
	}

	if !m.ByPrototypeRefFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.ByPrototypeRef")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.GetCode")
	}

	if !m.GetPrototypeFinished() {
		m.t.Fatal("Expected call to DescriptorsCacheMock.GetPrototype")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DescriptorsCacheMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DescriptorsCacheMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ByObjectDescriptorFinished()
		ok = ok && m.ByPrototypeRefFinished()
		ok = ok && m.GetCodeFinished()
		ok = ok && m.GetPrototypeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ByObjectDescriptorFinished() {
				m.t.Error("Expected call to DescriptorsCacheMock.ByObjectDescriptor")
			}

			if !m.ByPrototypeRefFinished() {
				m.t.Error("Expected call to DescriptorsCacheMock.ByPrototypeRef")
			}

			if !m.GetCodeFinished() {
				m.t.Error("Expected call to DescriptorsCacheMock.GetCode")
			}

			if !m.GetPrototypeFinished() {
				m.t.Error("Expected call to DescriptorsCacheMock.GetPrototype")
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
func (m *DescriptorsCacheMock) AllMocksCalled() bool {

	if !m.ByObjectDescriptorFinished() {
		return false
	}

	if !m.ByPrototypeRefFinished() {
		return false
	}

	if !m.GetCodeFinished() {
		return false
	}

	if !m.GetPrototypeFinished() {
		return false
	}

	return true
}
