package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "BaseNode" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	member "github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

//BaseNodeMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.BaseNode
type BaseNodeMock struct {
	t minimock.Tester

	GetNodeIDFunc       func() (r insolar.ShortNodeID)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mBaseNodeMockGetNodeID

	GetOpModeFunc       func() (r member.OpMode)
	GetOpModeCounter    uint64
	GetOpModePreCounter uint64
	GetOpModeMock       mBaseNodeMockGetOpMode

	GetSignatureVerifierFunc       func() (r cryptkit.SignatureVerifier)
	GetSignatureVerifierCounter    uint64
	GetSignatureVerifierPreCounter uint64
	GetSignatureVerifierMock       mBaseNodeMockGetSignatureVerifier

	GetStaticFunc       func() (r StaticProfile)
	GetStaticCounter    uint64
	GetStaticPreCounter uint64
	GetStaticMock       mBaseNodeMockGetStatic
}

//NewBaseNodeMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.BaseNode
func NewBaseNodeMock(t minimock.Tester) *BaseNodeMock {
	m := &BaseNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetNodeIDMock = mBaseNodeMockGetNodeID{mock: m}
	m.GetOpModeMock = mBaseNodeMockGetOpMode{mock: m}
	m.GetSignatureVerifierMock = mBaseNodeMockGetSignatureVerifier{mock: m}
	m.GetStaticMock = mBaseNodeMockGetStatic{mock: m}

	return m
}

type mBaseNodeMockGetNodeID struct {
	mock              *BaseNodeMock
	mainExpectation   *BaseNodeMockGetNodeIDExpectation
	expectationSeries []*BaseNodeMockGetNodeIDExpectation
}

type BaseNodeMockGetNodeIDExpectation struct {
	result *BaseNodeMockGetNodeIDResult
}

type BaseNodeMockGetNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of BaseNode.GetNodeID is expected from 1 to Infinity times
func (m *mBaseNodeMockGetNodeID) Expect() *mBaseNodeMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of BaseNode.GetNodeID
func (m *mBaseNodeMockGetNodeID) Return(r insolar.ShortNodeID) *BaseNodeMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &BaseNodeMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of BaseNode.GetNodeID is expected once
func (m *mBaseNodeMockGetNodeID) ExpectOnce() *BaseNodeMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &BaseNodeMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BaseNodeMockGetNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &BaseNodeMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of BaseNode.GetNodeID method
func (m *mBaseNodeMockGetNodeID) Set(f func() (r insolar.ShortNodeID)) *BaseNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.BaseNode interface
func (m *BaseNodeMock) GetNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BaseNodeMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to BaseNodeMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of BaseNodeMock.GetNodeIDFunc invocations
func (m *BaseNodeMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of BaseNodeMock.GetNodeID invocations
func (m *BaseNodeMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *BaseNodeMock) GetNodeIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNodeIDCounter) == uint64(len(m.GetNodeIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNodeIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNodeIDFunc != nil {
		return atomic.LoadUint64(&m.GetNodeIDCounter) > 0
	}

	return true
}

type mBaseNodeMockGetOpMode struct {
	mock              *BaseNodeMock
	mainExpectation   *BaseNodeMockGetOpModeExpectation
	expectationSeries []*BaseNodeMockGetOpModeExpectation
}

type BaseNodeMockGetOpModeExpectation struct {
	result *BaseNodeMockGetOpModeResult
}

type BaseNodeMockGetOpModeResult struct {
	r member.OpMode
}

//Expect specifies that invocation of BaseNode.GetOpMode is expected from 1 to Infinity times
func (m *mBaseNodeMockGetOpMode) Expect() *mBaseNodeMockGetOpMode {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetOpModeExpectation{}
	}

	return m
}

//Return specifies results of invocation of BaseNode.GetOpMode
func (m *mBaseNodeMockGetOpMode) Return(r member.OpMode) *BaseNodeMock {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetOpModeExpectation{}
	}
	m.mainExpectation.result = &BaseNodeMockGetOpModeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of BaseNode.GetOpMode is expected once
func (m *mBaseNodeMockGetOpMode) ExpectOnce() *BaseNodeMockGetOpModeExpectation {
	m.mock.GetOpModeFunc = nil
	m.mainExpectation = nil

	expectation := &BaseNodeMockGetOpModeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BaseNodeMockGetOpModeExpectation) Return(r member.OpMode) {
	e.result = &BaseNodeMockGetOpModeResult{r}
}

//Set uses given function f as a mock of BaseNode.GetOpMode method
func (m *mBaseNodeMockGetOpMode) Set(f func() (r member.OpMode)) *BaseNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOpModeFunc = f
	return m.mock
}

//GetOpMode implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.BaseNode interface
func (m *BaseNodeMock) GetOpMode() (r member.OpMode) {
	counter := atomic.AddUint64(&m.GetOpModePreCounter, 1)
	defer atomic.AddUint64(&m.GetOpModeCounter, 1)

	if len(m.GetOpModeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOpModeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BaseNodeMock.GetOpMode.")
			return
		}

		result := m.GetOpModeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetOpMode")
			return
		}

		r = result.r

		return
	}

	if m.GetOpModeMock.mainExpectation != nil {

		result := m.GetOpModeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetOpMode")
		}

		r = result.r

		return
	}

	if m.GetOpModeFunc == nil {
		m.t.Fatalf("Unexpected call to BaseNodeMock.GetOpMode.")
		return
	}

	return m.GetOpModeFunc()
}

//GetOpModeMinimockCounter returns a count of BaseNodeMock.GetOpModeFunc invocations
func (m *BaseNodeMock) GetOpModeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModeCounter)
}

//GetOpModeMinimockPreCounter returns the value of BaseNodeMock.GetOpMode invocations
func (m *BaseNodeMock) GetOpModeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModePreCounter)
}

//GetOpModeFinished returns true if mock invocations count is ok
func (m *BaseNodeMock) GetOpModeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetOpModeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetOpModeCounter) == uint64(len(m.GetOpModeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetOpModeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetOpModeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetOpModeFunc != nil {
		return atomic.LoadUint64(&m.GetOpModeCounter) > 0
	}

	return true
}

type mBaseNodeMockGetSignatureVerifier struct {
	mock              *BaseNodeMock
	mainExpectation   *BaseNodeMockGetSignatureVerifierExpectation
	expectationSeries []*BaseNodeMockGetSignatureVerifierExpectation
}

type BaseNodeMockGetSignatureVerifierExpectation struct {
	result *BaseNodeMockGetSignatureVerifierResult
}

type BaseNodeMockGetSignatureVerifierResult struct {
	r cryptkit.SignatureVerifier
}

//Expect specifies that invocation of BaseNode.GetSignatureVerifier is expected from 1 to Infinity times
func (m *mBaseNodeMockGetSignatureVerifier) Expect() *mBaseNodeMockGetSignatureVerifier {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetSignatureVerifierExpectation{}
	}

	return m
}

//Return specifies results of invocation of BaseNode.GetSignatureVerifier
func (m *mBaseNodeMockGetSignatureVerifier) Return(r cryptkit.SignatureVerifier) *BaseNodeMock {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetSignatureVerifierExpectation{}
	}
	m.mainExpectation.result = &BaseNodeMockGetSignatureVerifierResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of BaseNode.GetSignatureVerifier is expected once
func (m *mBaseNodeMockGetSignatureVerifier) ExpectOnce() *BaseNodeMockGetSignatureVerifierExpectation {
	m.mock.GetSignatureVerifierFunc = nil
	m.mainExpectation = nil

	expectation := &BaseNodeMockGetSignatureVerifierExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BaseNodeMockGetSignatureVerifierExpectation) Return(r cryptkit.SignatureVerifier) {
	e.result = &BaseNodeMockGetSignatureVerifierResult{r}
}

//Set uses given function f as a mock of BaseNode.GetSignatureVerifier method
func (m *mBaseNodeMockGetSignatureVerifier) Set(f func() (r cryptkit.SignatureVerifier)) *BaseNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureVerifierFunc = f
	return m.mock
}

//GetSignatureVerifier implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.BaseNode interface
func (m *BaseNodeMock) GetSignatureVerifier() (r cryptkit.SignatureVerifier) {
	counter := atomic.AddUint64(&m.GetSignatureVerifierPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureVerifierCounter, 1)

	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureVerifierMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BaseNodeMock.GetSignatureVerifier.")
			return
		}

		result := m.GetSignatureVerifierMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetSignatureVerifier")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierMock.mainExpectation != nil {

		result := m.GetSignatureVerifierMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetSignatureVerifier")
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierFunc == nil {
		m.t.Fatalf("Unexpected call to BaseNodeMock.GetSignatureVerifier.")
		return
	}

	return m.GetSignatureVerifierFunc()
}

//GetSignatureVerifierMinimockCounter returns a count of BaseNodeMock.GetSignatureVerifierFunc invocations
func (m *BaseNodeMock) GetSignatureVerifierMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierCounter)
}

//GetSignatureVerifierMinimockPreCounter returns the value of BaseNodeMock.GetSignatureVerifier invocations
func (m *BaseNodeMock) GetSignatureVerifierMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierPreCounter)
}

//GetSignatureVerifierFinished returns true if mock invocations count is ok
func (m *BaseNodeMock) GetSignatureVerifierFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureVerifierCounter) == uint64(len(m.GetSignatureVerifierMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureVerifierMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureVerifierCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureVerifierFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureVerifierCounter) > 0
	}

	return true
}

type mBaseNodeMockGetStatic struct {
	mock              *BaseNodeMock
	mainExpectation   *BaseNodeMockGetStaticExpectation
	expectationSeries []*BaseNodeMockGetStaticExpectation
}

type BaseNodeMockGetStaticExpectation struct {
	result *BaseNodeMockGetStaticResult
}

type BaseNodeMockGetStaticResult struct {
	r StaticProfile
}

//Expect specifies that invocation of BaseNode.GetStatic is expected from 1 to Infinity times
func (m *mBaseNodeMockGetStatic) Expect() *mBaseNodeMockGetStatic {
	m.mock.GetStaticFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetStaticExpectation{}
	}

	return m
}

//Return specifies results of invocation of BaseNode.GetStatic
func (m *mBaseNodeMockGetStatic) Return(r StaticProfile) *BaseNodeMock {
	m.mock.GetStaticFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &BaseNodeMockGetStaticExpectation{}
	}
	m.mainExpectation.result = &BaseNodeMockGetStaticResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of BaseNode.GetStatic is expected once
func (m *mBaseNodeMockGetStatic) ExpectOnce() *BaseNodeMockGetStaticExpectation {
	m.mock.GetStaticFunc = nil
	m.mainExpectation = nil

	expectation := &BaseNodeMockGetStaticExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *BaseNodeMockGetStaticExpectation) Return(r StaticProfile) {
	e.result = &BaseNodeMockGetStaticResult{r}
}

//Set uses given function f as a mock of BaseNode.GetStatic method
func (m *mBaseNodeMockGetStatic) Set(f func() (r StaticProfile)) *BaseNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStaticFunc = f
	return m.mock
}

//GetStatic implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.BaseNode interface
func (m *BaseNodeMock) GetStatic() (r StaticProfile) {
	counter := atomic.AddUint64(&m.GetStaticPreCounter, 1)
	defer atomic.AddUint64(&m.GetStaticCounter, 1)

	if len(m.GetStaticMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStaticMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to BaseNodeMock.GetStatic.")
			return
		}

		result := m.GetStaticMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetStatic")
			return
		}

		r = result.r

		return
	}

	if m.GetStaticMock.mainExpectation != nil {

		result := m.GetStaticMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the BaseNodeMock.GetStatic")
		}

		r = result.r

		return
	}

	if m.GetStaticFunc == nil {
		m.t.Fatalf("Unexpected call to BaseNodeMock.GetStatic.")
		return
	}

	return m.GetStaticFunc()
}

//GetStaticMinimockCounter returns a count of BaseNodeMock.GetStaticFunc invocations
func (m *BaseNodeMock) GetStaticMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticCounter)
}

//GetStaticMinimockPreCounter returns the value of BaseNodeMock.GetStatic invocations
func (m *BaseNodeMock) GetStaticMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticPreCounter)
}

//GetStaticFinished returns true if mock invocations count is ok
func (m *BaseNodeMock) GetStaticFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetStaticMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetStaticCounter) == uint64(len(m.GetStaticMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetStaticMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetStaticCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetStaticFunc != nil {
		return atomic.LoadUint64(&m.GetStaticCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BaseNodeMock) ValidateCallCounters() {

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetNodeID")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetOpMode")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetSignatureVerifier")
	}

	if !m.GetStaticFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetStatic")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *BaseNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *BaseNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *BaseNodeMock) MinimockFinish() {

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetNodeID")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetOpMode")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetSignatureVerifier")
	}

	if !m.GetStaticFinished() {
		m.t.Fatal("Expected call to BaseNodeMock.GetStatic")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *BaseNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *BaseNodeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetOpModeFinished()
		ok = ok && m.GetSignatureVerifierFinished()
		ok = ok && m.GetStaticFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to BaseNodeMock.GetNodeID")
			}

			if !m.GetOpModeFinished() {
				m.t.Error("Expected call to BaseNodeMock.GetOpMode")
			}

			if !m.GetSignatureVerifierFinished() {
				m.t.Error("Expected call to BaseNodeMock.GetSignatureVerifier")
			}

			if !m.GetStaticFinished() {
				m.t.Error("Expected call to BaseNodeMock.GetStatic")
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
func (m *BaseNodeMock) AllMocksCalled() bool {

	if !m.GetNodeIDFinished() {
		return false
	}

	if !m.GetOpModeFinished() {
		return false
	}

	if !m.GetSignatureVerifierFinished() {
		return false
	}

	if !m.GetStaticFinished() {
		return false
	}

	return true
}
