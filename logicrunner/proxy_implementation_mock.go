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

	RouteCallFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpRouteReq, p3 *rpctypes.UpRouteResp) (r error)
	RouteCallCounter    uint64
	RouteCallPreCounter uint64
	RouteCallMock       mProxyImplementationMockRouteCall

	SaveAsChildFunc       func(p context.Context, p1 *Transcript, p2 rpctypes.UpSaveAsChildReq, p3 *rpctypes.UpSaveAsChildResp) (r error)
	SaveAsChildCounter    uint64
	SaveAsChildPreCounter uint64
	SaveAsChildMock       mProxyImplementationMockSaveAsChild
}

//NewProxyImplementationMock returns a mock for github.com/insolar/insolar/logicrunner.ProxyImplementation
func NewProxyImplementationMock(t minimock.Tester) *ProxyImplementationMock {
	m := &ProxyImplementationMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeactivateObjectMock = mProxyImplementationMockDeactivateObject{mock: m}
	m.GetCodeMock = mProxyImplementationMockGetCode{mock: m}
	m.RouteCallMock = mProxyImplementationMockRouteCall{mock: m}
	m.SaveAsChildMock = mProxyImplementationMockSaveAsChild{mock: m}

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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProxyImplementationMock) ValidateCallCounters() {

	if !m.DeactivateObjectFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.DeactivateObject")
	}

	if !m.GetCodeFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.GetCode")
	}

	if !m.RouteCallFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.RouteCall")
	}

	if !m.SaveAsChildFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.SaveAsChild")
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

	if !m.RouteCallFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.RouteCall")
	}

	if !m.SaveAsChildFinished() {
		m.t.Fatal("Expected call to ProxyImplementationMock.SaveAsChild")
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
		ok = ok && m.RouteCallFinished()
		ok = ok && m.SaveAsChildFinished()

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

			if !m.RouteCallFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.RouteCall")
			}

			if !m.SaveAsChildFinished() {
				m.t.Error("Expected call to ProxyImplementationMock.SaveAsChild")
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

	if !m.RouteCallFinished() {
		return false
	}

	if !m.SaveAsChildFinished() {
		return false
	}

	return true
}
