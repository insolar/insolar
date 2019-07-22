package logicrunner

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "ProxyImplementation" can be found in github.com/insolar/insolar/logicrunner
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	rpctypes "github.com/insolar/insolar/logicrunner/goplugin/rpctypes"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProxyImplementationMock implements github.com/insolar/insolar/logicrunner.ProxyImplementation
type ProxyImplementationMock struct {
	t minimock.Tester

	DeactivateObjectFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpDeactivateObjectReq, p3 *rpctypes.UpDeactivateObjectResp) (r error)
	DeactivateObjectCounter    uint64
	DeactivateObjectPreCounter uint64
	DeactivateObjectMock       mProxyImplementationMockDeactivateObject

	GetCodeFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpGetCodeReq, p3 *rpctypes.UpGetCodeResp) (r error)
	GetCodeCounter    uint64
	GetCodePreCounter uint64
	GetCodeMock       mProxyImplementationMockGetCode

	GetDelegateFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpGetDelegateReq, p3 *rpctypes.UpGetDelegateResp) (r error)
	GetDelegateCounter    uint64
	GetDelegatePreCounter uint64
	GetDelegateMock       mProxyImplementationMockGetDelegate

	GetObjChildrenIteratorFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpGetObjChildrenIteratorReq, p3 *rpctypes.UpGetObjChildrenIteratorResp) (r error)
	GetObjChildrenIteratorCounter    uint64
	GetObjChildrenIteratorPreCounter uint64
	GetObjChildrenIteratorMock       mProxyImplementationMockGetObjChildrenIterator

	RouteCallFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpRouteReq, p3 *rpctypes.UpRouteResp) (r error)
	RouteCallCounter    uint64
	RouteCallPreCounter uint64
	RouteCallMock       mProxyImplementationMockRouteCall

	SaveFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveReq, p3 *rpctypes.UpSaveResp) (r error)
	SaveCounter    uint64
	SavePreCounter uint64
	SaveMock       mProxyImplementationMockSave

	SaveAsChildFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsChildReq, p3 *rpctypes.UpSaveAsChildResp) (r error)
	SaveAsChildCounter    uint64
	SaveAsChildPreCounter uint64
	SaveAsChildMock       mProxyImplementationMockSaveAsChild

	SaveAsDelegateFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsDelegateReq, p3 *rpctypes.UpSaveAsDelegateResp) (r error)
	SaveAsDelegateCounter    uint64
	SaveAsDelegatePreCounter uint64
	SaveAsDelegateMock       mProxyImplementationMockSaveAsDelegate
}

//NewProxyImplementationMock returns a mock for github.com/insolar/insolar/logicrunner.ProxyImplementation
func NewProxyImplementationMock(t minimock.Tester) *ProxyImplementationMock {
	m := &ProxyImplementationMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeactivateObjectMock = mProxyImplementationMockDeactivateObject{mock: m}
	m.GetCodeMock = mProxyImplementationMockGetCode{mock: m}
	m.GetDelegateMock = mProxyImplementationMockGetDelegate{mock: m}
	m.GetObjChildrenIteratorMock = mProxyImplementationMockGetObjChildrenIterator{mock: m}
	m.RouteCallMock = mProxyImplementationMockRouteCall{mock: m}
	m.SaveMock = mProxyImplementationMockSave{mock: m}
	m.SaveAsChildMock = mProxyImplementationMockSaveAsChild{mock: m}
	m.SaveAsDelegateMock = mProxyImplementationMockSaveAsDelegate{mock: m}

	return m
}

type mProxyImplementationMockDeactivateObject struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockDeactivateObjectExpectation
	expectationSeries []*ProxyImplementationMockDeactivateObjectExpectation
}

type ProxyImplementationMockDeactivateObjectExpectation struct {
	input  *ProxyImplementationMockDeactivateObjectInput
	result *ProxyImplementationMockDeactivateObjectResult
}

type ProxyImplementationMockDeactivateObjectInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpDeactivateObjectReq
	p3 *rpctypes.UpDeactivateObjectResp
}

type ProxyImplementationMockDeactivateObjectResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.DeactivateObject is expected from 1 to Infinity times
func (m *mProxyImplementationMockDeactivateObject) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpDeactivateObjectReq, p3 *rpctypes.UpDeactivateObjectResp) *mProxyImplementationMockDeactivateObject {
	m.mock.DeactivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockDeactivateObjectExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockDeactivateObjectInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.DeactivateObject
func (m *mProxyImplementationMockDeactivateObject) Return(r error) *ProxyImplementationMock {
	m.mock.DeactivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockDeactivateObjectExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockDeactivateObjectResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.DeactivateObject is expected once
func (m *mProxyImplementationMockDeactivateObject) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpDeactivateObjectReq, p3 *rpctypes.UpDeactivateObjectResp) *ProxyImplementationMockDeactivateObjectExpectation {
	m.mock.DeactivateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockDeactivateObjectExpectation{}
	expectation.input = &ProxyImplementationMockDeactivateObjectInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockDeactivateObjectExpectation) Return(r error) {
	e.result = &ProxyImplementationMockDeactivateObjectResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.DeactivateObject method
func (m *mProxyImplementationMockDeactivateObject) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpDeactivateObjectReq, p3 *rpctypes.UpDeactivateObjectResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DeactivateObjectFunc = f
	return m.mock
}

//DeactivateObject implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) DeactivateObject(p context.Context, p1 *Transcript, p2 rpctypes.UpDeactivateObjectReq, p3 *rpctypes.UpDeactivateObjectResp) (r error) {
	counter := atomic.AddUint64(&m.DeactivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.DeactivateObjectCounter, 1)

	if len(m.DeactivateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DeactivateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.DeactivateObject. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.DeactivateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockDeactivateObjectInput{p, p1, p2, p3}, "ProxyImplementation.DeactivateObject got unexpected parameters")

		result := m.DeactivateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.DeactivateObject")
			return
		}

		r = result.r

		return
	}

	if m.DeactivateObjectMock.mainExpectation != nil {

		input := m.DeactivateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockDeactivateObjectInput{p, p1, p2, p3}, "ProxyImplementation.DeactivateObject got unexpected parameters")
		}

		result := m.DeactivateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.DeactivateObject")
		}

		r = result.r

		return
	}

	if m.DeactivateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.DeactivateObject. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.DeactivateObjectFunc(p, p1, p2, p3)
}

//DeactivateObjectMinimockCounter returns a count of ProxyImplementationMock.DeactivateObjectFunc invocations
func (m *ProxyImplementationMock) DeactivateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateObjectCounter)
}

//DeactivateObjectMinimockPreCounter returns the value of ProxyImplementationMock.DeactivateObject invocations
func (m *ProxyImplementationMock) DeactivateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeactivateObjectPreCounter)
}

//DeactivateObjectFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) DeactivateObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DeactivateObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DeactivateObjectCounter) == uint64(len(m.DeactivateObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DeactivateObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DeactivateObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DeactivateObjectFunc != nil {
		return atomic.LoadUint64(&m.DeactivateObjectCounter) > 0
	}

	return true
}

type mProxyImplementationMockGetCode struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockGetCodeExpectation
	expectationSeries []*ProxyImplementationMockGetCodeExpectation
}

type ProxyImplementationMockGetCodeExpectation struct {
	input  *ProxyImplementationMockGetCodeInput
	result *ProxyImplementationMockGetCodeResult
}

type ProxyImplementationMockGetCodeInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpGetCodeReq
	p3 *rpctypes.UpGetCodeResp
}

type ProxyImplementationMockGetCodeResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.GetCode is expected from 1 to Infinity times
func (m *mProxyImplementationMockGetCode) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpGetCodeReq, p3 *rpctypes.UpGetCodeResp) *mProxyImplementationMockGetCode {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockGetCodeExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockGetCodeInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.GetCode
func (m *mProxyImplementationMockGetCode) Return(r error) *ProxyImplementationMock {
	m.mock.GetCodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockGetCodeExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockGetCodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.GetCode is expected once
func (m *mProxyImplementationMockGetCode) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpGetCodeReq, p3 *rpctypes.UpGetCodeResp) *ProxyImplementationMockGetCodeExpectation {
	m.mock.GetCodeFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockGetCodeExpectation{}
	expectation.input = &ProxyImplementationMockGetCodeInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockGetCodeExpectation) Return(r error) {
	e.result = &ProxyImplementationMockGetCodeResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.GetCode method
func (m *mProxyImplementationMockGetCode) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpGetCodeReq, p3 *rpctypes.UpGetCodeResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetCodeFunc = f
	return m.mock
}

//GetCode implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) GetCode(p context.Context, p1 *Transcript, p2 rpctypes.UpGetCodeReq, p3 *rpctypes.UpGetCodeResp) (r error) {
	counter := atomic.AddUint64(&m.GetCodePreCounter, 1)
	defer atomic.AddUint64(&m.GetCodeCounter, 1)

	if len(m.GetCodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetCodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.GetCode. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.GetCodeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockGetCodeInput{p, p1, p2, p3}, "ProxyImplementation.GetCode got unexpected parameters")

		result := m.GetCodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.GetCode")
			return
		}

		r = result.r

		return
	}

	if m.GetCodeMock.mainExpectation != nil {

		input := m.GetCodeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockGetCodeInput{p, p1, p2, p3}, "ProxyImplementation.GetCode got unexpected parameters")
		}

		result := m.GetCodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.GetCode")
		}

		r = result.r

		return
	}

	if m.GetCodeFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.GetCode. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.GetCodeFunc(p, p1, p2, p3)
}

//GetCodeMinimockCounter returns a count of ProxyImplementationMock.GetCodeFunc invocations
func (m *ProxyImplementationMock) GetCodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodeCounter)
}

//GetCodeMinimockPreCounter returns the value of ProxyImplementationMock.GetCode invocations
func (m *ProxyImplementationMock) GetCodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetCodePreCounter)
}

//GetCodeFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) GetCodeFinished() bool {
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

type mProxyImplementationMockGetDelegate struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockGetDelegateExpectation
	expectationSeries []*ProxyImplementationMockGetDelegateExpectation
}

type ProxyImplementationMockGetDelegateExpectation struct {
	input  *ProxyImplementationMockGetDelegateInput
	result *ProxyImplementationMockGetDelegateResult
}

type ProxyImplementationMockGetDelegateInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpGetDelegateReq
	p3 *rpctypes.UpGetDelegateResp
}

type ProxyImplementationMockGetDelegateResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.GetDelegate is expected from 1 to Infinity times
func (m *mProxyImplementationMockGetDelegate) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpGetDelegateReq, p3 *rpctypes.UpGetDelegateResp) *mProxyImplementationMockGetDelegate {
	m.mock.GetDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockGetDelegateExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockGetDelegateInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.GetDelegate
func (m *mProxyImplementationMockGetDelegate) Return(r error) *ProxyImplementationMock {
	m.mock.GetDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockGetDelegateExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockGetDelegateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.GetDelegate is expected once
func (m *mProxyImplementationMockGetDelegate) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpGetDelegateReq, p3 *rpctypes.UpGetDelegateResp) *ProxyImplementationMockGetDelegateExpectation {
	m.mock.GetDelegateFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockGetDelegateExpectation{}
	expectation.input = &ProxyImplementationMockGetDelegateInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockGetDelegateExpectation) Return(r error) {
	e.result = &ProxyImplementationMockGetDelegateResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.GetDelegate method
func (m *mProxyImplementationMockGetDelegate) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpGetDelegateReq, p3 *rpctypes.UpGetDelegateResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDelegateFunc = f
	return m.mock
}

//GetDelegate implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) GetDelegate(p context.Context, p1 *Transcript, p2 rpctypes.UpGetDelegateReq, p3 *rpctypes.UpGetDelegateResp) (r error) {
	counter := atomic.AddUint64(&m.GetDelegatePreCounter, 1)
	defer atomic.AddUint64(&m.GetDelegateCounter, 1)

	if len(m.GetDelegateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDelegateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.GetDelegate. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.GetDelegateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockGetDelegateInput{p, p1, p2, p3}, "ProxyImplementation.GetDelegate got unexpected parameters")

		result := m.GetDelegateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.GetDelegate")
			return
		}

		r = result.r

		return
	}

	if m.GetDelegateMock.mainExpectation != nil {

		input := m.GetDelegateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockGetDelegateInput{p, p1, p2, p3}, "ProxyImplementation.GetDelegate got unexpected parameters")
		}

		result := m.GetDelegateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.GetDelegate")
		}

		r = result.r

		return
	}

	if m.GetDelegateFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.GetDelegate. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.GetDelegateFunc(p, p1, p2, p3)
}

//GetDelegateMinimockCounter returns a count of ProxyImplementationMock.GetDelegateFunc invocations
func (m *ProxyImplementationMock) GetDelegateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDelegateCounter)
}

//GetDelegateMinimockPreCounter returns the value of ProxyImplementationMock.GetDelegate invocations
func (m *ProxyImplementationMock) GetDelegateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDelegatePreCounter)
}

//GetDelegateFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) GetDelegateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDelegateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDelegateCounter) == uint64(len(m.GetDelegateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDelegateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDelegateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDelegateFunc != nil {
		return atomic.LoadUint64(&m.GetDelegateCounter) > 0
	}

	return true
}

type mProxyImplementationMockGetObjChildrenIterator struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockGetObjChildrenIteratorExpectation
	expectationSeries []*ProxyImplementationMockGetObjChildrenIteratorExpectation
}

type ProxyImplementationMockGetObjChildrenIteratorExpectation struct {
	input  *ProxyImplementationMockGetObjChildrenIteratorInput
	result *ProxyImplementationMockGetObjChildrenIteratorResult
}

type ProxyImplementationMockGetObjChildrenIteratorInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpGetObjChildrenIteratorReq
	p3 *rpctypes.UpGetObjChildrenIteratorResp
}

type ProxyImplementationMockGetObjChildrenIteratorResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.GetObjChildrenIterator is expected from 1 to Infinity times
func (m *mProxyImplementationMockGetObjChildrenIterator) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpGetObjChildrenIteratorReq, p3 *rpctypes.UpGetObjChildrenIteratorResp) *mProxyImplementationMockGetObjChildrenIterator {
	m.mock.GetObjChildrenIteratorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockGetObjChildrenIteratorExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockGetObjChildrenIteratorInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.GetObjChildrenIterator
func (m *mProxyImplementationMockGetObjChildrenIterator) Return(r error) *ProxyImplementationMock {
	m.mock.GetObjChildrenIteratorFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockGetObjChildrenIteratorExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockGetObjChildrenIteratorResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.GetObjChildrenIterator is expected once
func (m *mProxyImplementationMockGetObjChildrenIterator) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpGetObjChildrenIteratorReq, p3 *rpctypes.UpGetObjChildrenIteratorResp) *ProxyImplementationMockGetObjChildrenIteratorExpectation {
	m.mock.GetObjChildrenIteratorFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockGetObjChildrenIteratorExpectation{}
	expectation.input = &ProxyImplementationMockGetObjChildrenIteratorInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockGetObjChildrenIteratorExpectation) Return(r error) {
	e.result = &ProxyImplementationMockGetObjChildrenIteratorResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.GetObjChildrenIterator method
func (m *mProxyImplementationMockGetObjChildrenIterator) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpGetObjChildrenIteratorReq, p3 *rpctypes.UpGetObjChildrenIteratorResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetObjChildrenIteratorFunc = f
	return m.mock
}

//GetObjChildrenIterator implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) GetObjChildrenIterator(p context.Context, p1 *Transcript, p2 rpctypes.UpGetObjChildrenIteratorReq, p3 *rpctypes.UpGetObjChildrenIteratorResp) (r error) {
	counter := atomic.AddUint64(&m.GetObjChildrenIteratorPreCounter, 1)
	defer atomic.AddUint64(&m.GetObjChildrenIteratorCounter, 1)

	if len(m.GetObjChildrenIteratorMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetObjChildrenIteratorMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.GetObjChildrenIterator. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.GetObjChildrenIteratorMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockGetObjChildrenIteratorInput{p, p1, p2, p3}, "ProxyImplementation.GetObjChildrenIterator got unexpected parameters")

		result := m.GetObjChildrenIteratorMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.GetObjChildrenIterator")
			return
		}

		r = result.r

		return
	}

	if m.GetObjChildrenIteratorMock.mainExpectation != nil {

		input := m.GetObjChildrenIteratorMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockGetObjChildrenIteratorInput{p, p1, p2, p3}, "ProxyImplementation.GetObjChildrenIterator got unexpected parameters")
		}

		result := m.GetObjChildrenIteratorMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.GetObjChildrenIterator")
		}

		r = result.r

		return
	}

	if m.GetObjChildrenIteratorFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.GetObjChildrenIterator. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.GetObjChildrenIteratorFunc(p, p1, p2, p3)
}

//GetObjChildrenIteratorMinimockCounter returns a count of ProxyImplementationMock.GetObjChildrenIteratorFunc invocations
func (m *ProxyImplementationMock) GetObjChildrenIteratorMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjChildrenIteratorCounter)
}

//GetObjChildrenIteratorMinimockPreCounter returns the value of ProxyImplementationMock.GetObjChildrenIterator invocations
func (m *ProxyImplementationMock) GetObjChildrenIteratorMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetObjChildrenIteratorPreCounter)
}

//GetObjChildrenIteratorFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) GetObjChildrenIteratorFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetObjChildrenIteratorMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetObjChildrenIteratorCounter) == uint64(len(m.GetObjChildrenIteratorMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetObjChildrenIteratorMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetObjChildrenIteratorCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetObjChildrenIteratorFunc != nil {
		return atomic.LoadUint64(&m.GetObjChildrenIteratorCounter) > 0
	}

	return true
}

type mProxyImplementationMockRouteCall struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockRouteCallExpectation
	expectationSeries []*ProxyImplementationMockRouteCallExpectation
}

type ProxyImplementationMockRouteCallExpectation struct {
	input  *ProxyImplementationMockRouteCallInput
	result *ProxyImplementationMockRouteCallResult
}

type ProxyImplementationMockRouteCallInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpRouteReq
	p3 *rpctypes.UpRouteResp
}

type ProxyImplementationMockRouteCallResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.RouteCall is expected from 1 to Infinity times
func (m *mProxyImplementationMockRouteCall) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpRouteReq, p3 *rpctypes.UpRouteResp) *mProxyImplementationMockRouteCall {
	m.mock.RouteCallFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockRouteCallExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockRouteCallInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.RouteCall
func (m *mProxyImplementationMockRouteCall) Return(r error) *ProxyImplementationMock {
	m.mock.RouteCallFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockRouteCallExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockRouteCallResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.RouteCall is expected once
func (m *mProxyImplementationMockRouteCall) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpRouteReq, p3 *rpctypes.UpRouteResp) *ProxyImplementationMockRouteCallExpectation {
	m.mock.RouteCallFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockRouteCallExpectation{}
	expectation.input = &ProxyImplementationMockRouteCallInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockRouteCallExpectation) Return(r error) {
	e.result = &ProxyImplementationMockRouteCallResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.RouteCall method
func (m *mProxyImplementationMockRouteCall) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpRouteReq, p3 *rpctypes.UpRouteResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RouteCallFunc = f
	return m.mock
}

//RouteCall implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) RouteCall(p context.Context, p1 *Transcript, p2 rpctypes.UpRouteReq, p3 *rpctypes.UpRouteResp) (r error) {
	counter := atomic.AddUint64(&m.RouteCallPreCounter, 1)
	defer atomic.AddUint64(&m.RouteCallCounter, 1)

	if len(m.RouteCallMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RouteCallMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.RouteCall. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RouteCallMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockRouteCallInput{p, p1, p2, p3}, "ProxyImplementation.RouteCall got unexpected parameters")

		result := m.RouteCallMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.RouteCall")
			return
		}

		r = result.r

		return
	}

	if m.RouteCallMock.mainExpectation != nil {

		input := m.RouteCallMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockRouteCallInput{p, p1, p2, p3}, "ProxyImplementation.RouteCall got unexpected parameters")
		}

		result := m.RouteCallMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.RouteCall")
		}

		r = result.r

		return
	}

	if m.RouteCallFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.RouteCall. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RouteCallFunc(p, p1, p2, p3)
}

//RouteCallMinimockCounter returns a count of ProxyImplementationMock.RouteCallFunc invocations
func (m *ProxyImplementationMock) RouteCallMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RouteCallCounter)
}

//RouteCallMinimockPreCounter returns the value of ProxyImplementationMock.RouteCall invocations
func (m *ProxyImplementationMock) RouteCallMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RouteCallPreCounter)
}

//RouteCallFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) RouteCallFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RouteCallMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RouteCallCounter) == uint64(len(m.RouteCallMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RouteCallMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RouteCallCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RouteCallFunc != nil {
		return atomic.LoadUint64(&m.RouteCallCounter) > 0
	}

	return true
}

type mProxyImplementationMockSave struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockSaveExpectation
	expectationSeries []*ProxyImplementationMockSaveExpectation
}

type ProxyImplementationMockSaveExpectation struct {
	input  *ProxyImplementationMockSaveInput
	result *ProxyImplementationMockSaveResult
}

type ProxyImplementationMockSaveInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpSaveReq
	p3 *rpctypes.UpSaveResp
}

type ProxyImplementationMockSaveResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.Save is expected from 1 to Infinity times
func (m *mProxyImplementationMockSave) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveReq, p3 *rpctypes.UpSaveResp) *mProxyImplementationMockSave {
	m.mock.SaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockSaveExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockSaveInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.Save
func (m *mProxyImplementationMockSave) Return(r error) *ProxyImplementationMock {
	m.mock.SaveFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockSaveExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockSaveResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.Save is expected once
func (m *mProxyImplementationMockSave) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveReq, p3 *rpctypes.UpSaveResp) *ProxyImplementationMockSaveExpectation {
	m.mock.SaveFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockSaveExpectation{}
	expectation.input = &ProxyImplementationMockSaveInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockSaveExpectation) Return(r error) {
	e.result = &ProxyImplementationMockSaveResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.Save method
func (m *mProxyImplementationMockSave) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveReq, p3 *rpctypes.UpSaveResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SaveFunc = f
	return m.mock
}

//Save implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) Save(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveReq, p3 *rpctypes.UpSaveResp) (r error) {
	counter := atomic.AddUint64(&m.SavePreCounter, 1)
	defer atomic.AddUint64(&m.SaveCounter, 1)

	if len(m.SaveMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SaveMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.Save. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SaveMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockSaveInput{p, p1, p2, p3}, "ProxyImplementation.Save got unexpected parameters")

		result := m.SaveMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.Save")
			return
		}

		r = result.r

		return
	}

	if m.SaveMock.mainExpectation != nil {

		input := m.SaveMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockSaveInput{p, p1, p2, p3}, "ProxyImplementation.Save got unexpected parameters")
		}

		result := m.SaveMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.Save")
		}

		r = result.r

		return
	}

	if m.SaveFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.Save. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SaveFunc(p, p1, p2, p3)
}

//SaveMinimockCounter returns a count of ProxyImplementationMock.SaveFunc invocations
func (m *ProxyImplementationMock) SaveMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SaveCounter)
}

//SaveMinimockPreCounter returns the value of ProxyImplementationMock.Save invocations
func (m *ProxyImplementationMock) SaveMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SavePreCounter)
}

//SaveFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) SaveFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SaveMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SaveCounter) == uint64(len(m.SaveMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SaveMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SaveCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SaveFunc != nil {
		return atomic.LoadUint64(&m.SaveCounter) > 0
	}

	return true
}

type mProxyImplementationMockSaveAsChild struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockSaveAsChildExpectation
	expectationSeries []*ProxyImplementationMockSaveAsChildExpectation
}

type ProxyImplementationMockSaveAsChildExpectation struct {
	input  *ProxyImplementationMockSaveAsChildInput
	result *ProxyImplementationMockSaveAsChildResult
}

type ProxyImplementationMockSaveAsChildInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpSaveAsChildReq
	p3 *rpctypes.UpSaveAsChildResp
}

type ProxyImplementationMockSaveAsChildResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.SaveAsChild is expected from 1 to Infinity times
func (m *mProxyImplementationMockSaveAsChild) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsChildReq, p3 *rpctypes.UpSaveAsChildResp) *mProxyImplementationMockSaveAsChild {
	m.mock.SaveAsChildFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockSaveAsChildExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockSaveAsChildInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.SaveAsChild
func (m *mProxyImplementationMockSaveAsChild) Return(r error) *ProxyImplementationMock {
	m.mock.SaveAsChildFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockSaveAsChildExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockSaveAsChildResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.SaveAsChild is expected once
func (m *mProxyImplementationMockSaveAsChild) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsChildReq, p3 *rpctypes.UpSaveAsChildResp) *ProxyImplementationMockSaveAsChildExpectation {
	m.mock.SaveAsChildFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockSaveAsChildExpectation{}
	expectation.input = &ProxyImplementationMockSaveAsChildInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockSaveAsChildExpectation) Return(r error) {
	e.result = &ProxyImplementationMockSaveAsChildResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.SaveAsChild method
func (m *mProxyImplementationMockSaveAsChild) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsChildReq, p3 *rpctypes.UpSaveAsChildResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SaveAsChildFunc = f
	return m.mock
}

//SaveAsChild implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) SaveAsChild(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsChildReq, p3 *rpctypes.UpSaveAsChildResp) (r error) {
	counter := atomic.AddUint64(&m.SaveAsChildPreCounter, 1)
	defer atomic.AddUint64(&m.SaveAsChildCounter, 1)

	if len(m.SaveAsChildMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SaveAsChildMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.SaveAsChild. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SaveAsChildMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockSaveAsChildInput{p, p1, p2, p3}, "ProxyImplementation.SaveAsChild got unexpected parameters")

		result := m.SaveAsChildMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.SaveAsChild")
			return
		}

		r = result.r

		return
	}

	if m.SaveAsChildMock.mainExpectation != nil {

		input := m.SaveAsChildMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockSaveAsChildInput{p, p1, p2, p3}, "ProxyImplementation.SaveAsChild got unexpected parameters")
		}

		result := m.SaveAsChildMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.SaveAsChild")
		}

		r = result.r

		return
	}

	if m.SaveAsChildFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.SaveAsChild. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SaveAsChildFunc(p, p1, p2, p3)
}

//SaveAsChildMinimockCounter returns a count of ProxyImplementationMock.SaveAsChildFunc invocations
func (m *ProxyImplementationMock) SaveAsChildMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SaveAsChildCounter)
}

//SaveAsChildMinimockPreCounter returns the value of ProxyImplementationMock.SaveAsChild invocations
func (m *ProxyImplementationMock) SaveAsChildMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SaveAsChildPreCounter)
}

//SaveAsChildFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) SaveAsChildFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SaveAsChildMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SaveAsChildCounter) == uint64(len(m.SaveAsChildMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SaveAsChildMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SaveAsChildCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SaveAsChildFunc != nil {
		return atomic.LoadUint64(&m.SaveAsChildCounter) > 0
	}

	return true
}

type mProxyImplementationMockSaveAsDelegate struct {
	mock              *ProxyImplementationMock
	mainExpectation   *ProxyImplementationMockSaveAsDelegateExpectation
	expectationSeries []*ProxyImplementationMockSaveAsDelegateExpectation
}

type ProxyImplementationMockSaveAsDelegateExpectation struct {
	input  *ProxyImplementationMockSaveAsDelegateInput
	result *ProxyImplementationMockSaveAsDelegateResult
}

type ProxyImplementationMockSaveAsDelegateInput struct {
	p  context.Context
	p1 *Transcript
	p2 rpctypes.UpSaveAsDelegateReq
	p3 *rpctypes.UpSaveAsDelegateResp
}

type ProxyImplementationMockSaveAsDelegateResult struct {
	r error
}

//Expect specifies that invocation of ProxyImplementation.SaveAsDelegate is expected from 1 to Infinity times
func (m *mProxyImplementationMockSaveAsDelegate) Expect(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsDelegateReq, p3 *rpctypes.UpSaveAsDelegateResp) *mProxyImplementationMockSaveAsDelegate {
	m.mock.SaveAsDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockSaveAsDelegateExpectation{}
	}
	m.mainExpectation.input = &ProxyImplementationMockSaveAsDelegateInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of ProxyImplementation.SaveAsDelegate
func (m *mProxyImplementationMockSaveAsDelegate) Return(r error) *ProxyImplementationMock {
	m.mock.SaveAsDelegateFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ProxyImplementationMockSaveAsDelegateExpectation{}
	}
	m.mainExpectation.result = &ProxyImplementationMockSaveAsDelegateResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of ProxyImplementation.SaveAsDelegate is expected once
func (m *mProxyImplementationMockSaveAsDelegate) ExpectOnce(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsDelegateReq, p3 *rpctypes.UpSaveAsDelegateResp) *ProxyImplementationMockSaveAsDelegateExpectation {
	m.mock.SaveAsDelegateFunc = nil
	m.mainExpectation = nil

	expectation := &ProxyImplementationMockSaveAsDelegateExpectation{}
	expectation.input = &ProxyImplementationMockSaveAsDelegateInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ProxyImplementationMockSaveAsDelegateExpectation) Return(r error) {
	e.result = &ProxyImplementationMockSaveAsDelegateResult{r}
}

//Set uses given function f as a mock of ProxyImplementation.SaveAsDelegate method
func (m *mProxyImplementationMockSaveAsDelegate) Set(f func(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsDelegateReq, p3 *rpctypes.UpSaveAsDelegateResp) (r error)) *ProxyImplementationMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SaveAsDelegateFunc = f
	return m.mock
}

//SaveAsDelegate implements github.com/insolar/insolar/logicrunner.ProxyImplementation interface
func (m *ProxyImplementationMock) SaveAsDelegate(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsDelegateReq, p3 *rpctypes.UpSaveAsDelegateResp) (r error) {
	counter := atomic.AddUint64(&m.SaveAsDelegatePreCounter, 1)
	defer atomic.AddUint64(&m.SaveAsDelegateCounter, 1)

	if len(m.SaveAsDelegateMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SaveAsDelegateMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ProxyImplementationMock.SaveAsDelegate. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.SaveAsDelegateMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ProxyImplementationMockSaveAsDelegateInput{p, p1, p2, p3}, "ProxyImplementation.SaveAsDelegate got unexpected parameters")

		result := m.SaveAsDelegateMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.SaveAsDelegate")
			return
		}

		r = result.r

		return
	}

	if m.SaveAsDelegateMock.mainExpectation != nil {

		input := m.SaveAsDelegateMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ProxyImplementationMockSaveAsDelegateInput{p, p1, p2, p3}, "ProxyImplementation.SaveAsDelegate got unexpected parameters")
		}

		result := m.SaveAsDelegateMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ProxyImplementationMock.SaveAsDelegate")
		}

		r = result.r

		return
	}

	if m.SaveAsDelegateFunc == nil {
		m.t.Fatalf("Unexpected call to ProxyImplementationMock.SaveAsDelegate. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.SaveAsDelegateFunc(p, p1, p2, p3)
}

//SaveAsDelegateMinimockCounter returns a count of ProxyImplementationMock.SaveAsDelegateFunc invocations
func (m *ProxyImplementationMock) SaveAsDelegateMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SaveAsDelegateCounter)
}

//SaveAsDelegateMinimockPreCounter returns the value of ProxyImplementationMock.SaveAsDelegate invocations
func (m *ProxyImplementationMock) SaveAsDelegateMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SaveAsDelegatePreCounter)
}

//SaveAsDelegateFinished returns true if mock invocations count is ok
func (m *ProxyImplementationMock) SaveAsDelegateFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SaveAsDelegateMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SaveAsDelegateCounter) == uint64(len(m.SaveAsDelegateMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SaveAsDelegateMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SaveAsDelegateCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SaveAsDelegateFunc != nil {
		return atomic.LoadUint64(&m.SaveAsDelegateCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProxyImplementationMock) ValidateCallCounters() {

	if !m.DeactivateObjectFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.DeactivateObject")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.GetCode")
	}

	if !m.GetDelegateFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.GetDelegate")
	}

	if !m.GetObjChildrenIteratorFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.GetObjChildrenIterator")
	}

	if !m.RouteCallFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.RouteCall")
	}

	if !m.SaveFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.Save")
	}

	if !m.SaveAsChildFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.SaveAsChild")
	}

	if !m.SaveAsDelegateFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.SaveAsDelegate")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProxyImplementationMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ProxyImplementationMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ProxyImplementationMock) MinimockFinish() {

	if !m.DeactivateObjectFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.DeactivateObject")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.GetCode")
	}

	if !m.GetDelegateFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.GetDelegate")
	}

	if !m.GetObjChildrenIteratorFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.GetObjChildrenIterator")
	}

	if !m.RouteCallFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.RouteCall")
	}

	if !m.SaveFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.Save")
	}

	if !m.SaveAsChildFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.SaveAsChild")
	}

	if !m.SaveAsDelegateFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.SaveAsDelegate")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ProxyImplementationMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ProxyImplementationMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.DeactivateObjectFinished()
		ok = ok && m.GetCodeFinished()
		ok = ok && m.GetDelegateFinished()
		ok = ok && m.GetObjChildrenIteratorFinished()
		ok = ok && m.RouteCallFinished()
		ok = ok && m.SaveFinished()
		ok = ok && m.SaveAsChildFinished()
		ok = ok && m.SaveAsDelegateFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.DeactivateObjectFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.DeactivateObject")
			}

			if !m.GetCodeFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.GetCode")
			}

			if !m.GetDelegateFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.GetDelegate")
			}

			if !m.GetObjChildrenIteratorFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.GetObjChildrenIterator")
			}

			if !m.RouteCallFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.RouteCall")
			}

			if !m.SaveFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.Save")
			}

			if !m.SaveAsChildFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.SaveAsChild")
			}

			if !m.SaveAsDelegateFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.SaveAsDelegate")
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
func (m *ProxyImplementationMock) AllMocksCalled() bool {

	if !m.DeactivateObjectFinished() {
		return false
	}

	if !m.GetCodeFinished() {
		return false
	}

	if !m.GetDelegateFinished() {
		return false
	}

	if !m.GetObjChildrenIteratorFinished() {
		return false
	}

	if !m.RouteCallFinished() {
		return false
	}

	if !m.SaveFinished() {
		return false
	}

	if !m.SaveAsChildFinished() {
		return false
	}

	if !m.SaveAsDelegateFinished() {
		return false
	}

	return true
}
