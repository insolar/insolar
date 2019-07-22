package misbehavior

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Report" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	endpoints "github.com/insolar/insolar/network/consensus/common/endpoints"
	profiles "github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

//ReportMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report
type ReportMock struct {
	t minimock.Tester

	CaptureMarkFunc       func() (r interface{})
	CaptureMarkCounter    uint64
	CaptureMarkPreCounter uint64
	CaptureMarkMock       mReportMockCaptureMark

	DetailsFunc       func() (r []interface{})
	DetailsCounter    uint64
	DetailsPreCounter uint64
	DetailsMock       mReportMockDetails

	MisbehaviorTypeFunc       func() (r Type)
	MisbehaviorTypeCounter    uint64
	MisbehaviorTypePreCounter uint64
	MisbehaviorTypeMock       mReportMockMisbehaviorType

	ViolatorHostFunc       func() (r endpoints.InboundConnection)
	ViolatorHostCounter    uint64
	ViolatorHostPreCounter uint64
	ViolatorHostMock       mReportMockViolatorHost

	ViolatorNodeFunc       func() (r profiles.BaseNode)
	ViolatorNodeCounter    uint64
	ViolatorNodePreCounter uint64
	ViolatorNodeMock       mReportMockViolatorNode
}

//NewReportMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report
func NewReportMock(t minimock.Tester) *ReportMock {
	m := &ReportMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CaptureMarkMock = mReportMockCaptureMark{mock: m}
	m.DetailsMock = mReportMockDetails{mock: m}
	m.MisbehaviorTypeMock = mReportMockMisbehaviorType{mock: m}
	m.ViolatorHostMock = mReportMockViolatorHost{mock: m}
	m.ViolatorNodeMock = mReportMockViolatorNode{mock: m}

	return m
}

type mReportMockCaptureMark struct {
	mock              *ReportMock
	mainExpectation   *ReportMockCaptureMarkExpectation
	expectationSeries []*ReportMockCaptureMarkExpectation
}

type ReportMockCaptureMarkExpectation struct {
	result *ReportMockCaptureMarkResult
}

type ReportMockCaptureMarkResult struct {
	r interface{}
}

//Expect specifies that invocation of Report.CaptureMark is expected from 1 to Infinity times
func (m *mReportMockCaptureMark) Expect() *mReportMockCaptureMark {
	m.mock.CaptureMarkFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockCaptureMarkExpectation{}
	}

	return m
}

//Return specifies results of invocation of Report.CaptureMark
func (m *mReportMockCaptureMark) Return(r interface{}) *ReportMock {
	m.mock.CaptureMarkFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockCaptureMarkExpectation{}
	}
	m.mainExpectation.result = &ReportMockCaptureMarkResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Report.CaptureMark is expected once
func (m *mReportMockCaptureMark) ExpectOnce() *ReportMockCaptureMarkExpectation {
	m.mock.CaptureMarkFunc = nil
	m.mainExpectation = nil

	expectation := &ReportMockCaptureMarkExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReportMockCaptureMarkExpectation) Return(r interface{}) {
	e.result = &ReportMockCaptureMarkResult{r}
}

//Set uses given function f as a mock of Report.CaptureMark method
func (m *mReportMockCaptureMark) Set(f func() (r interface{})) *ReportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CaptureMarkFunc = f
	return m.mock
}

//CaptureMark implements github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report interface
func (m *ReportMock) CaptureMark() (r interface{}) {
	counter := atomic.AddUint64(&m.CaptureMarkPreCounter, 1)
	defer atomic.AddUint64(&m.CaptureMarkCounter, 1)

	if len(m.CaptureMarkMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CaptureMarkMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReportMock.CaptureMark.")
			return
		}

		result := m.CaptureMarkMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.CaptureMark")
			return
		}

		r = result.r

		return
	}

	if m.CaptureMarkMock.mainExpectation != nil {

		result := m.CaptureMarkMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.CaptureMark")
		}

		r = result.r

		return
	}

	if m.CaptureMarkFunc == nil {
		m.t.Fatalf("Unexpected call to ReportMock.CaptureMark.")
		return
	}

	return m.CaptureMarkFunc()
}

//CaptureMarkMinimockCounter returns a count of ReportMock.CaptureMarkFunc invocations
func (m *ReportMock) CaptureMarkMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CaptureMarkCounter)
}

//CaptureMarkMinimockPreCounter returns the value of ReportMock.CaptureMark invocations
func (m *ReportMock) CaptureMarkMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CaptureMarkPreCounter)
}

//CaptureMarkFinished returns true if mock invocations count is ok
func (m *ReportMock) CaptureMarkFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CaptureMarkMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CaptureMarkCounter) == uint64(len(m.CaptureMarkMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CaptureMarkMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CaptureMarkCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CaptureMarkFunc != nil {
		return atomic.LoadUint64(&m.CaptureMarkCounter) > 0
	}

	return true
}

type mReportMockDetails struct {
	mock              *ReportMock
	mainExpectation   *ReportMockDetailsExpectation
	expectationSeries []*ReportMockDetailsExpectation
}

type ReportMockDetailsExpectation struct {
	result *ReportMockDetailsResult
}

type ReportMockDetailsResult struct {
	r []interface{}
}

//Expect specifies that invocation of Report.Details is expected from 1 to Infinity times
func (m *mReportMockDetails) Expect() *mReportMockDetails {
	m.mock.DetailsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockDetailsExpectation{}
	}

	return m
}

//Return specifies results of invocation of Report.Details
func (m *mReportMockDetails) Return(r []interface{}) *ReportMock {
	m.mock.DetailsFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockDetailsExpectation{}
	}
	m.mainExpectation.result = &ReportMockDetailsResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Report.Details is expected once
func (m *mReportMockDetails) ExpectOnce() *ReportMockDetailsExpectation {
	m.mock.DetailsFunc = nil
	m.mainExpectation = nil

	expectation := &ReportMockDetailsExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReportMockDetailsExpectation) Return(r []interface{}) {
	e.result = &ReportMockDetailsResult{r}
}

//Set uses given function f as a mock of Report.Details method
func (m *mReportMockDetails) Set(f func() (r []interface{})) *ReportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.DetailsFunc = f
	return m.mock
}

//Details implements github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report interface
func (m *ReportMock) Details() (r []interface{}) {
	counter := atomic.AddUint64(&m.DetailsPreCounter, 1)
	defer atomic.AddUint64(&m.DetailsCounter, 1)

	if len(m.DetailsMock.expectationSeries) > 0 {
		if counter > uint64(len(m.DetailsMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReportMock.Details.")
			return
		}

		result := m.DetailsMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.Details")
			return
		}

		r = result.r

		return
	}

	if m.DetailsMock.mainExpectation != nil {

		result := m.DetailsMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.Details")
		}

		r = result.r

		return
	}

	if m.DetailsFunc == nil {
		m.t.Fatalf("Unexpected call to ReportMock.Details.")
		return
	}

	return m.DetailsFunc()
}

//DetailsMinimockCounter returns a count of ReportMock.DetailsFunc invocations
func (m *ReportMock) DetailsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DetailsCounter)
}

//DetailsMinimockPreCounter returns the value of ReportMock.Details invocations
func (m *ReportMock) DetailsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DetailsPreCounter)
}

//DetailsFinished returns true if mock invocations count is ok
func (m *ReportMock) DetailsFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.DetailsMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.DetailsCounter) == uint64(len(m.DetailsMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.DetailsMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.DetailsCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.DetailsFunc != nil {
		return atomic.LoadUint64(&m.DetailsCounter) > 0
	}

	return true
}

type mReportMockMisbehaviorType struct {
	mock              *ReportMock
	mainExpectation   *ReportMockMisbehaviorTypeExpectation
	expectationSeries []*ReportMockMisbehaviorTypeExpectation
}

type ReportMockMisbehaviorTypeExpectation struct {
	result *ReportMockMisbehaviorTypeResult
}

type ReportMockMisbehaviorTypeResult struct {
	r Type
}

//Expect specifies that invocation of Report.MisbehaviorType is expected from 1 to Infinity times
func (m *mReportMockMisbehaviorType) Expect() *mReportMockMisbehaviorType {
	m.mock.MisbehaviorTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockMisbehaviorTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Report.MisbehaviorType
func (m *mReportMockMisbehaviorType) Return(r Type) *ReportMock {
	m.mock.MisbehaviorTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockMisbehaviorTypeExpectation{}
	}
	m.mainExpectation.result = &ReportMockMisbehaviorTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Report.MisbehaviorType is expected once
func (m *mReportMockMisbehaviorType) ExpectOnce() *ReportMockMisbehaviorTypeExpectation {
	m.mock.MisbehaviorTypeFunc = nil
	m.mainExpectation = nil

	expectation := &ReportMockMisbehaviorTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReportMockMisbehaviorTypeExpectation) Return(r Type) {
	e.result = &ReportMockMisbehaviorTypeResult{r}
}

//Set uses given function f as a mock of Report.MisbehaviorType method
func (m *mReportMockMisbehaviorType) Set(f func() (r Type)) *ReportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MisbehaviorTypeFunc = f
	return m.mock
}

//MisbehaviorType implements github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report interface
func (m *ReportMock) MisbehaviorType() (r Type) {
	counter := atomic.AddUint64(&m.MisbehaviorTypePreCounter, 1)
	defer atomic.AddUint64(&m.MisbehaviorTypeCounter, 1)

	if len(m.MisbehaviorTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MisbehaviorTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReportMock.MisbehaviorType.")
			return
		}

		result := m.MisbehaviorTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.MisbehaviorType")
			return
		}

		r = result.r

		return
	}

	if m.MisbehaviorTypeMock.mainExpectation != nil {

		result := m.MisbehaviorTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.MisbehaviorType")
		}

		r = result.r

		return
	}

	if m.MisbehaviorTypeFunc == nil {
		m.t.Fatalf("Unexpected call to ReportMock.MisbehaviorType.")
		return
	}

	return m.MisbehaviorTypeFunc()
}

//MisbehaviorTypeMinimockCounter returns a count of ReportMock.MisbehaviorTypeFunc invocations
func (m *ReportMock) MisbehaviorTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MisbehaviorTypeCounter)
}

//MisbehaviorTypeMinimockPreCounter returns the value of ReportMock.MisbehaviorType invocations
func (m *ReportMock) MisbehaviorTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MisbehaviorTypePreCounter)
}

//MisbehaviorTypeFinished returns true if mock invocations count is ok
func (m *ReportMock) MisbehaviorTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MisbehaviorTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MisbehaviorTypeCounter) == uint64(len(m.MisbehaviorTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MisbehaviorTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MisbehaviorTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MisbehaviorTypeFunc != nil {
		return atomic.LoadUint64(&m.MisbehaviorTypeCounter) > 0
	}

	return true
}

type mReportMockViolatorHost struct {
	mock              *ReportMock
	mainExpectation   *ReportMockViolatorHostExpectation
	expectationSeries []*ReportMockViolatorHostExpectation
}

type ReportMockViolatorHostExpectation struct {
	result *ReportMockViolatorHostResult
}

type ReportMockViolatorHostResult struct {
	r endpoints.InboundConnection
}

//Expect specifies that invocation of Report.ViolatorHost is expected from 1 to Infinity times
func (m *mReportMockViolatorHost) Expect() *mReportMockViolatorHost {
	m.mock.ViolatorHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockViolatorHostExpectation{}
	}

	return m
}

//Return specifies results of invocation of Report.ViolatorHost
func (m *mReportMockViolatorHost) Return(r endpoints.InboundConnection) *ReportMock {
	m.mock.ViolatorHostFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockViolatorHostExpectation{}
	}
	m.mainExpectation.result = &ReportMockViolatorHostResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Report.ViolatorHost is expected once
func (m *mReportMockViolatorHost) ExpectOnce() *ReportMockViolatorHostExpectation {
	m.mock.ViolatorHostFunc = nil
	m.mainExpectation = nil

	expectation := &ReportMockViolatorHostExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReportMockViolatorHostExpectation) Return(r endpoints.InboundConnection) {
	e.result = &ReportMockViolatorHostResult{r}
}

//Set uses given function f as a mock of Report.ViolatorHost method
func (m *mReportMockViolatorHost) Set(f func() (r endpoints.InboundConnection)) *ReportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ViolatorHostFunc = f
	return m.mock
}

//ViolatorHost implements github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report interface
func (m *ReportMock) ViolatorHost() (r endpoints.InboundConnection) {
	counter := atomic.AddUint64(&m.ViolatorHostPreCounter, 1)
	defer atomic.AddUint64(&m.ViolatorHostCounter, 1)

	if len(m.ViolatorHostMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ViolatorHostMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReportMock.ViolatorHost.")
			return
		}

		result := m.ViolatorHostMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.ViolatorHost")
			return
		}

		r = result.r

		return
	}

	if m.ViolatorHostMock.mainExpectation != nil {

		result := m.ViolatorHostMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.ViolatorHost")
		}

		r = result.r

		return
	}

	if m.ViolatorHostFunc == nil {
		m.t.Fatalf("Unexpected call to ReportMock.ViolatorHost.")
		return
	}

	return m.ViolatorHostFunc()
}

//ViolatorHostMinimockCounter returns a count of ReportMock.ViolatorHostFunc invocations
func (m *ReportMock) ViolatorHostMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ViolatorHostCounter)
}

//ViolatorHostMinimockPreCounter returns the value of ReportMock.ViolatorHost invocations
func (m *ReportMock) ViolatorHostMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ViolatorHostPreCounter)
}

//ViolatorHostFinished returns true if mock invocations count is ok
func (m *ReportMock) ViolatorHostFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ViolatorHostMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ViolatorHostCounter) == uint64(len(m.ViolatorHostMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ViolatorHostMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ViolatorHostCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ViolatorHostFunc != nil {
		return atomic.LoadUint64(&m.ViolatorHostCounter) > 0
	}

	return true
}

type mReportMockViolatorNode struct {
	mock              *ReportMock
	mainExpectation   *ReportMockViolatorNodeExpectation
	expectationSeries []*ReportMockViolatorNodeExpectation
}

type ReportMockViolatorNodeExpectation struct {
	result *ReportMockViolatorNodeResult
}

type ReportMockViolatorNodeResult struct {
	r profiles.BaseNode
}

//Expect specifies that invocation of Report.ViolatorNode is expected from 1 to Infinity times
func (m *mReportMockViolatorNode) Expect() *mReportMockViolatorNode {
	m.mock.ViolatorNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockViolatorNodeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Report.ViolatorNode
func (m *mReportMockViolatorNode) Return(r profiles.BaseNode) *ReportMock {
	m.mock.ViolatorNodeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ReportMockViolatorNodeExpectation{}
	}
	m.mainExpectation.result = &ReportMockViolatorNodeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Report.ViolatorNode is expected once
func (m *mReportMockViolatorNode) ExpectOnce() *ReportMockViolatorNodeExpectation {
	m.mock.ViolatorNodeFunc = nil
	m.mainExpectation = nil

	expectation := &ReportMockViolatorNodeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ReportMockViolatorNodeExpectation) Return(r profiles.BaseNode) {
	e.result = &ReportMockViolatorNodeResult{r}
}

//Set uses given function f as a mock of Report.ViolatorNode method
func (m *mReportMockViolatorNode) Set(f func() (r profiles.BaseNode)) *ReportMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ViolatorNodeFunc = f
	return m.mock
}

//ViolatorNode implements github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior.Report interface
func (m *ReportMock) ViolatorNode() (r profiles.BaseNode) {
	counter := atomic.AddUint64(&m.ViolatorNodePreCounter, 1)
	defer atomic.AddUint64(&m.ViolatorNodeCounter, 1)

	if len(m.ViolatorNodeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ViolatorNodeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ReportMock.ViolatorNode.")
			return
		}

		result := m.ViolatorNodeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.ViolatorNode")
			return
		}

		r = result.r

		return
	}

	if m.ViolatorNodeMock.mainExpectation != nil {

		result := m.ViolatorNodeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ReportMock.ViolatorNode")
		}

		r = result.r

		return
	}

	if m.ViolatorNodeFunc == nil {
		m.t.Fatalf("Unexpected call to ReportMock.ViolatorNode.")
		return
	}

	return m.ViolatorNodeFunc()
}

//ViolatorNodeMinimockCounter returns a count of ReportMock.ViolatorNodeFunc invocations
func (m *ReportMock) ViolatorNodeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ViolatorNodeCounter)
}

//ViolatorNodeMinimockPreCounter returns the value of ReportMock.ViolatorNode invocations
func (m *ReportMock) ViolatorNodeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ViolatorNodePreCounter)
}

//ViolatorNodeFinished returns true if mock invocations count is ok
func (m *ReportMock) ViolatorNodeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ViolatorNodeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ViolatorNodeCounter) == uint64(len(m.ViolatorNodeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ViolatorNodeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ViolatorNodeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ViolatorNodeFunc != nil {
		return atomic.LoadUint64(&m.ViolatorNodeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ReportMock) ValidateCallCounters() {

	if !m.CaptureMarkFinished() {
		m.t.Fatal("Expected call to ReportMock.CaptureMark")
	}

	if !m.DetailsFinished() {
		m.t.Fatal("Expected call to ReportMock.Details")
	}

	if !m.MisbehaviorTypeFinished() {
		m.t.Fatal("Expected call to ReportMock.MisbehaviorType")
	}

	if !m.ViolatorHostFinished() {
		m.t.Fatal("Expected call to ReportMock.ViolatorHost")
	}

	if !m.ViolatorNodeFinished() {
		m.t.Fatal("Expected call to ReportMock.ViolatorNode")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ReportMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ReportMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ReportMock) MinimockFinish() {

	if !m.CaptureMarkFinished() {
		m.t.Fatal("Expected call to ReportMock.CaptureMark")
	}

	if !m.DetailsFinished() {
		m.t.Fatal("Expected call to ReportMock.Details")
	}

	if !m.MisbehaviorTypeFinished() {
		m.t.Fatal("Expected call to ReportMock.MisbehaviorType")
	}

	if !m.ViolatorHostFinished() {
		m.t.Fatal("Expected call to ReportMock.ViolatorHost")
	}

	if !m.ViolatorNodeFinished() {
		m.t.Fatal("Expected call to ReportMock.ViolatorNode")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ReportMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ReportMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CaptureMarkFinished()
		ok = ok && m.DetailsFinished()
		ok = ok && m.MisbehaviorTypeFinished()
		ok = ok && m.ViolatorHostFinished()
		ok = ok && m.ViolatorNodeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CaptureMarkFinished() {
				m.t.Error("Expected call to ReportMock.CaptureMark")
			}

			if !m.DetailsFinished() {
				m.t.Error("Expected call to ReportMock.Details")
			}

			if !m.MisbehaviorTypeFinished() {
				m.t.Error("Expected call to ReportMock.MisbehaviorType")
			}

			if !m.ViolatorHostFinished() {
				m.t.Error("Expected call to ReportMock.ViolatorHost")
			}

			if !m.ViolatorNodeFinished() {
				m.t.Error("Expected call to ReportMock.ViolatorNode")
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
func (m *ReportMock) AllMocksCalled() bool {

	if !m.CaptureMarkFinished() {
		return false
	}

	if !m.DetailsFinished() {
		return false
	}

	if !m.MisbehaviorTypeFinished() {
		return false
	}

	if !m.ViolatorHostFinished() {
		return false
	}

	if !m.ViolatorNodeFinished() {
		return false
	}

	return true
}
