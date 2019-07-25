package profiles

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "LocalNode" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/profiles
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	member "github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

//LocalNodeMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode
type LocalNodeMock struct {
	t minimock.Tester

	GetDeclaredPowerFunc       func() (r member.Power)
	GetDeclaredPowerCounter    uint64
	GetDeclaredPowerPreCounter uint64
	GetDeclaredPowerMock       mLocalNodeMockGetDeclaredPower

	GetIndexFunc       func() (r member.Index)
	GetIndexCounter    uint64
	GetIndexPreCounter uint64
	GetIndexMock       mLocalNodeMockGetIndex

	GetNodeIDFunc       func() (r insolar.ShortNodeID)
	GetNodeIDCounter    uint64
	GetNodeIDPreCounter uint64
	GetNodeIDMock       mLocalNodeMockGetNodeID

	GetOpModeFunc       func() (r member.OpMode)
	GetOpModeCounter    uint64
	GetOpModePreCounter uint64
	GetOpModeMock       mLocalNodeMockGetOpMode

	GetSignatureVerifierFunc       func() (r cryptkit.SignatureVerifier)
	GetSignatureVerifierCounter    uint64
	GetSignatureVerifierPreCounter uint64
	GetSignatureVerifierMock       mLocalNodeMockGetSignatureVerifier

	GetStaticFunc       func() (r StaticProfile)
	GetStaticCounter    uint64
	GetStaticPreCounter uint64
	GetStaticMock       mLocalNodeMockGetStatic

	IsJoinerFunc       func() (r bool)
	IsJoinerCounter    uint64
	IsJoinerPreCounter uint64
	IsJoinerMock       mLocalNodeMockIsJoiner

	LocalNodeProfileFunc       func()
	LocalNodeProfileCounter    uint64
	LocalNodeProfilePreCounter uint64
	LocalNodeProfileMock       mLocalNodeMockLocalNodeProfile
}

//NewLocalNodeMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode
func NewLocalNodeMock(t minimock.Tester) *LocalNodeMock {
	m := &LocalNodeMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetDeclaredPowerMock = mLocalNodeMockGetDeclaredPower{mock: m}
	m.GetIndexMock = mLocalNodeMockGetIndex{mock: m}
	m.GetNodeIDMock = mLocalNodeMockGetNodeID{mock: m}
	m.GetOpModeMock = mLocalNodeMockGetOpMode{mock: m}
	m.GetSignatureVerifierMock = mLocalNodeMockGetSignatureVerifier{mock: m}
	m.GetStaticMock = mLocalNodeMockGetStatic{mock: m}
	m.IsJoinerMock = mLocalNodeMockIsJoiner{mock: m}
	m.LocalNodeProfileMock = mLocalNodeMockLocalNodeProfile{mock: m}

	return m
}

type mLocalNodeMockGetDeclaredPower struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetDeclaredPowerExpectation
	expectationSeries []*LocalNodeMockGetDeclaredPowerExpectation
}

type LocalNodeMockGetDeclaredPowerExpectation struct {
	result *LocalNodeMockGetDeclaredPowerResult
}

type LocalNodeMockGetDeclaredPowerResult struct {
	r member.Power
}

//Expect specifies that invocation of LocalNode.GetDeclaredPower is expected from 1 to Infinity times
func (m *mLocalNodeMockGetDeclaredPower) Expect() *mLocalNodeMockGetDeclaredPower {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetDeclaredPowerExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetDeclaredPower
func (m *mLocalNodeMockGetDeclaredPower) Return(r member.Power) *LocalNodeMock {
	m.mock.GetDeclaredPowerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetDeclaredPowerExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetDeclaredPowerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetDeclaredPower is expected once
func (m *mLocalNodeMockGetDeclaredPower) ExpectOnce() *LocalNodeMockGetDeclaredPowerExpectation {
	m.mock.GetDeclaredPowerFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetDeclaredPowerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetDeclaredPowerExpectation) Return(r member.Power) {
	e.result = &LocalNodeMockGetDeclaredPowerResult{r}
}

//Set uses given function f as a mock of LocalNode.GetDeclaredPower method
func (m *mLocalNodeMockGetDeclaredPower) Set(f func() (r member.Power)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDeclaredPowerFunc = f
	return m.mock
}

//GetDeclaredPower implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetDeclaredPower() (r member.Power) {
	counter := atomic.AddUint64(&m.GetDeclaredPowerPreCounter, 1)
	defer atomic.AddUint64(&m.GetDeclaredPowerCounter, 1)

	if len(m.GetDeclaredPowerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDeclaredPowerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetDeclaredPower.")
			return
		}

		result := m.GetDeclaredPowerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetDeclaredPower")
			return
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerMock.mainExpectation != nil {

		result := m.GetDeclaredPowerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetDeclaredPower")
		}

		r = result.r

		return
	}

	if m.GetDeclaredPowerFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetDeclaredPower.")
		return
	}

	return m.GetDeclaredPowerFunc()
}

//GetDeclaredPowerMinimockCounter returns a count of LocalNodeMock.GetDeclaredPowerFunc invocations
func (m *LocalNodeMock) GetDeclaredPowerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerCounter)
}

//GetDeclaredPowerMinimockPreCounter returns the value of LocalNodeMock.GetDeclaredPower invocations
func (m *LocalNodeMock) GetDeclaredPowerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDeclaredPowerPreCounter)
}

//GetDeclaredPowerFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetDeclaredPowerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDeclaredPowerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDeclaredPowerCounter) == uint64(len(m.GetDeclaredPowerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDeclaredPowerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDeclaredPowerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDeclaredPowerFunc != nil {
		return atomic.LoadUint64(&m.GetDeclaredPowerCounter) > 0
	}

	return true
}

type mLocalNodeMockGetIndex struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetIndexExpectation
	expectationSeries []*LocalNodeMockGetIndexExpectation
}

type LocalNodeMockGetIndexExpectation struct {
	result *LocalNodeMockGetIndexResult
}

type LocalNodeMockGetIndexResult struct {
	r member.Index
}

//Expect specifies that invocation of LocalNode.GetIndex is expected from 1 to Infinity times
func (m *mLocalNodeMockGetIndex) Expect() *mLocalNodeMockGetIndex {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetIndexExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetIndex
func (m *mLocalNodeMockGetIndex) Return(r member.Index) *LocalNodeMock {
	m.mock.GetIndexFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetIndexExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetIndexResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetIndex is expected once
func (m *mLocalNodeMockGetIndex) ExpectOnce() *LocalNodeMockGetIndexExpectation {
	m.mock.GetIndexFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetIndexExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetIndexExpectation) Return(r member.Index) {
	e.result = &LocalNodeMockGetIndexResult{r}
}

//Set uses given function f as a mock of LocalNode.GetIndex method
func (m *mLocalNodeMockGetIndex) Set(f func() (r member.Index)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIndexFunc = f
	return m.mock
}

//GetIndex implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetIndex() (r member.Index) {
	counter := atomic.AddUint64(&m.GetIndexPreCounter, 1)
	defer atomic.AddUint64(&m.GetIndexCounter, 1)

	if len(m.GetIndexMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIndexMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetIndex.")
			return
		}

		result := m.GetIndexMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetIndex")
			return
		}

		r = result.r

		return
	}

	if m.GetIndexMock.mainExpectation != nil {

		result := m.GetIndexMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetIndex")
		}

		r = result.r

		return
	}

	if m.GetIndexFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetIndex.")
		return
	}

	return m.GetIndexFunc()
}

//GetIndexMinimockCounter returns a count of LocalNodeMock.GetIndexFunc invocations
func (m *LocalNodeMock) GetIndexMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexCounter)
}

//GetIndexMinimockPreCounter returns the value of LocalNodeMock.GetIndex invocations
func (m *LocalNodeMock) GetIndexMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIndexPreCounter)
}

//GetIndexFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetIndexFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIndexMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIndexCounter) == uint64(len(m.GetIndexMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIndexMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIndexCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIndexFunc != nil {
		return atomic.LoadUint64(&m.GetIndexCounter) > 0
	}

	return true
}

type mLocalNodeMockGetNodeID struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetNodeIDExpectation
	expectationSeries []*LocalNodeMockGetNodeIDExpectation
}

type LocalNodeMockGetNodeIDExpectation struct {
	result *LocalNodeMockGetNodeIDResult
}

type LocalNodeMockGetNodeIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of LocalNode.GetNodeID is expected from 1 to Infinity times
func (m *mLocalNodeMockGetNodeID) Expect() *mLocalNodeMockGetNodeID {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetNodeIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetNodeID
func (m *mLocalNodeMockGetNodeID) Return(r insolar.ShortNodeID) *LocalNodeMock {
	m.mock.GetNodeIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetNodeIDExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetNodeIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetNodeID is expected once
func (m *mLocalNodeMockGetNodeID) ExpectOnce() *LocalNodeMockGetNodeIDExpectation {
	m.mock.GetNodeIDFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetNodeIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetNodeIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &LocalNodeMockGetNodeIDResult{r}
}

//Set uses given function f as a mock of LocalNode.GetNodeID method
func (m *mLocalNodeMockGetNodeID) Set(f func() (r insolar.ShortNodeID)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNodeIDFunc = f
	return m.mock
}

//GetNodeID implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetNodeID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetNodeIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetNodeIDCounter, 1)

	if len(m.GetNodeIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNodeIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetNodeID.")
			return
		}

		result := m.GetNodeIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetNodeID")
			return
		}

		r = result.r

		return
	}

	if m.GetNodeIDMock.mainExpectation != nil {

		result := m.GetNodeIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetNodeID")
		}

		r = result.r

		return
	}

	if m.GetNodeIDFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetNodeID.")
		return
	}

	return m.GetNodeIDFunc()
}

//GetNodeIDMinimockCounter returns a count of LocalNodeMock.GetNodeIDFunc invocations
func (m *LocalNodeMock) GetNodeIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDCounter)
}

//GetNodeIDMinimockPreCounter returns the value of LocalNodeMock.GetNodeID invocations
func (m *LocalNodeMock) GetNodeIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNodeIDPreCounter)
}

//GetNodeIDFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetNodeIDFinished() bool {
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

type mLocalNodeMockGetOpMode struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetOpModeExpectation
	expectationSeries []*LocalNodeMockGetOpModeExpectation
}

type LocalNodeMockGetOpModeExpectation struct {
	result *LocalNodeMockGetOpModeResult
}

type LocalNodeMockGetOpModeResult struct {
	r member.OpMode
}

//Expect specifies that invocation of LocalNode.GetOpMode is expected from 1 to Infinity times
func (m *mLocalNodeMockGetOpMode) Expect() *mLocalNodeMockGetOpMode {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetOpModeExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetOpMode
func (m *mLocalNodeMockGetOpMode) Return(r member.OpMode) *LocalNodeMock {
	m.mock.GetOpModeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetOpModeExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetOpModeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetOpMode is expected once
func (m *mLocalNodeMockGetOpMode) ExpectOnce() *LocalNodeMockGetOpModeExpectation {
	m.mock.GetOpModeFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetOpModeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetOpModeExpectation) Return(r member.OpMode) {
	e.result = &LocalNodeMockGetOpModeResult{r}
}

//Set uses given function f as a mock of LocalNode.GetOpMode method
func (m *mLocalNodeMockGetOpMode) Set(f func() (r member.OpMode)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetOpModeFunc = f
	return m.mock
}

//GetOpMode implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetOpMode() (r member.OpMode) {
	counter := atomic.AddUint64(&m.GetOpModePreCounter, 1)
	defer atomic.AddUint64(&m.GetOpModeCounter, 1)

	if len(m.GetOpModeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetOpModeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetOpMode.")
			return
		}

		result := m.GetOpModeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetOpMode")
			return
		}

		r = result.r

		return
	}

	if m.GetOpModeMock.mainExpectation != nil {

		result := m.GetOpModeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetOpMode")
		}

		r = result.r

		return
	}

	if m.GetOpModeFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetOpMode.")
		return
	}

	return m.GetOpModeFunc()
}

//GetOpModeMinimockCounter returns a count of LocalNodeMock.GetOpModeFunc invocations
func (m *LocalNodeMock) GetOpModeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModeCounter)
}

//GetOpModeMinimockPreCounter returns the value of LocalNodeMock.GetOpMode invocations
func (m *LocalNodeMock) GetOpModeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetOpModePreCounter)
}

//GetOpModeFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetOpModeFinished() bool {
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

type mLocalNodeMockGetSignatureVerifier struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetSignatureVerifierExpectation
	expectationSeries []*LocalNodeMockGetSignatureVerifierExpectation
}

type LocalNodeMockGetSignatureVerifierExpectation struct {
	result *LocalNodeMockGetSignatureVerifierResult
}

type LocalNodeMockGetSignatureVerifierResult struct {
	r cryptkit.SignatureVerifier
}

//Expect specifies that invocation of LocalNode.GetSignatureVerifier is expected from 1 to Infinity times
func (m *mLocalNodeMockGetSignatureVerifier) Expect() *mLocalNodeMockGetSignatureVerifier {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetSignatureVerifierExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetSignatureVerifier
func (m *mLocalNodeMockGetSignatureVerifier) Return(r cryptkit.SignatureVerifier) *LocalNodeMock {
	m.mock.GetSignatureVerifierFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetSignatureVerifierExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetSignatureVerifierResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetSignatureVerifier is expected once
func (m *mLocalNodeMockGetSignatureVerifier) ExpectOnce() *LocalNodeMockGetSignatureVerifierExpectation {
	m.mock.GetSignatureVerifierFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetSignatureVerifierExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetSignatureVerifierExpectation) Return(r cryptkit.SignatureVerifier) {
	e.result = &LocalNodeMockGetSignatureVerifierResult{r}
}

//Set uses given function f as a mock of LocalNode.GetSignatureVerifier method
func (m *mLocalNodeMockGetSignatureVerifier) Set(f func() (r cryptkit.SignatureVerifier)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureVerifierFunc = f
	return m.mock
}

//GetSignatureVerifier implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetSignatureVerifier() (r cryptkit.SignatureVerifier) {
	counter := atomic.AddUint64(&m.GetSignatureVerifierPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureVerifierCounter, 1)

	if len(m.GetSignatureVerifierMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureVerifierMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetSignatureVerifier.")
			return
		}

		result := m.GetSignatureVerifierMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetSignatureVerifier")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierMock.mainExpectation != nil {

		result := m.GetSignatureVerifierMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetSignatureVerifier")
		}

		r = result.r

		return
	}

	if m.GetSignatureVerifierFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetSignatureVerifier.")
		return
	}

	return m.GetSignatureVerifierFunc()
}

//GetSignatureVerifierMinimockCounter returns a count of LocalNodeMock.GetSignatureVerifierFunc invocations
func (m *LocalNodeMock) GetSignatureVerifierMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierCounter)
}

//GetSignatureVerifierMinimockPreCounter returns the value of LocalNodeMock.GetSignatureVerifier invocations
func (m *LocalNodeMock) GetSignatureVerifierMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureVerifierPreCounter)
}

//GetSignatureVerifierFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetSignatureVerifierFinished() bool {
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

type mLocalNodeMockGetStatic struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockGetStaticExpectation
	expectationSeries []*LocalNodeMockGetStaticExpectation
}

type LocalNodeMockGetStaticExpectation struct {
	result *LocalNodeMockGetStaticResult
}

type LocalNodeMockGetStaticResult struct {
	r StaticProfile
}

//Expect specifies that invocation of LocalNode.GetStatic is expected from 1 to Infinity times
func (m *mLocalNodeMockGetStatic) Expect() *mLocalNodeMockGetStatic {
	m.mock.GetStaticFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetStaticExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.GetStatic
func (m *mLocalNodeMockGetStatic) Return(r StaticProfile) *LocalNodeMock {
	m.mock.GetStaticFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockGetStaticExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockGetStaticResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.GetStatic is expected once
func (m *mLocalNodeMockGetStatic) ExpectOnce() *LocalNodeMockGetStaticExpectation {
	m.mock.GetStaticFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockGetStaticExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockGetStaticExpectation) Return(r StaticProfile) {
	e.result = &LocalNodeMockGetStaticResult{r}
}

//Set uses given function f as a mock of LocalNode.GetStatic method
func (m *mLocalNodeMockGetStatic) Set(f func() (r StaticProfile)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetStaticFunc = f
	return m.mock
}

//GetStatic implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) GetStatic() (r StaticProfile) {
	counter := atomic.AddUint64(&m.GetStaticPreCounter, 1)
	defer atomic.AddUint64(&m.GetStaticCounter, 1)

	if len(m.GetStaticMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetStaticMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.GetStatic.")
			return
		}

		result := m.GetStaticMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetStatic")
			return
		}

		r = result.r

		return
	}

	if m.GetStaticMock.mainExpectation != nil {

		result := m.GetStaticMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.GetStatic")
		}

		r = result.r

		return
	}

	if m.GetStaticFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.GetStatic.")
		return
	}

	return m.GetStaticFunc()
}

//GetStaticMinimockCounter returns a count of LocalNodeMock.GetStaticFunc invocations
func (m *LocalNodeMock) GetStaticMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticCounter)
}

//GetStaticMinimockPreCounter returns the value of LocalNodeMock.GetStatic invocations
func (m *LocalNodeMock) GetStaticMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetStaticPreCounter)
}

//GetStaticFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) GetStaticFinished() bool {
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

type mLocalNodeMockIsJoiner struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockIsJoinerExpectation
	expectationSeries []*LocalNodeMockIsJoinerExpectation
}

type LocalNodeMockIsJoinerExpectation struct {
	result *LocalNodeMockIsJoinerResult
}

type LocalNodeMockIsJoinerResult struct {
	r bool
}

//Expect specifies that invocation of LocalNode.IsJoiner is expected from 1 to Infinity times
func (m *mLocalNodeMockIsJoiner) Expect() *mLocalNodeMockIsJoiner {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockIsJoinerExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.IsJoiner
func (m *mLocalNodeMockIsJoiner) Return(r bool) *LocalNodeMock {
	m.mock.IsJoinerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockIsJoinerExpectation{}
	}
	m.mainExpectation.result = &LocalNodeMockIsJoinerResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.IsJoiner is expected once
func (m *mLocalNodeMockIsJoiner) ExpectOnce() *LocalNodeMockIsJoinerExpectation {
	m.mock.IsJoinerFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockIsJoinerExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *LocalNodeMockIsJoinerExpectation) Return(r bool) {
	e.result = &LocalNodeMockIsJoinerResult{r}
}

//Set uses given function f as a mock of LocalNode.IsJoiner method
func (m *mLocalNodeMockIsJoiner) Set(f func() (r bool)) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsJoinerFunc = f
	return m.mock
}

//IsJoiner implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) IsJoiner() (r bool) {
	counter := atomic.AddUint64(&m.IsJoinerPreCounter, 1)
	defer atomic.AddUint64(&m.IsJoinerCounter, 1)

	if len(m.IsJoinerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsJoinerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.IsJoiner.")
			return
		}

		result := m.IsJoinerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.IsJoiner")
			return
		}

		r = result.r

		return
	}

	if m.IsJoinerMock.mainExpectation != nil {

		result := m.IsJoinerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the LocalNodeMock.IsJoiner")
		}

		r = result.r

		return
	}

	if m.IsJoinerFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.IsJoiner.")
		return
	}

	return m.IsJoinerFunc()
}

//IsJoinerMinimockCounter returns a count of LocalNodeMock.IsJoinerFunc invocations
func (m *LocalNodeMock) IsJoinerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerCounter)
}

//IsJoinerMinimockPreCounter returns the value of LocalNodeMock.IsJoiner invocations
func (m *LocalNodeMock) IsJoinerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsJoinerPreCounter)
}

//IsJoinerFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) IsJoinerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsJoinerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsJoinerCounter) == uint64(len(m.IsJoinerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsJoinerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsJoinerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsJoinerFunc != nil {
		return atomic.LoadUint64(&m.IsJoinerCounter) > 0
	}

	return true
}

type mLocalNodeMockLocalNodeProfile struct {
	mock              *LocalNodeMock
	mainExpectation   *LocalNodeMockLocalNodeProfileExpectation
	expectationSeries []*LocalNodeMockLocalNodeProfileExpectation
}

type LocalNodeMockLocalNodeProfileExpectation struct {
}

//Expect specifies that invocation of LocalNode.LocalNodeProfile is expected from 1 to Infinity times
func (m *mLocalNodeMockLocalNodeProfile) Expect() *mLocalNodeMockLocalNodeProfile {
	m.mock.LocalNodeProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockLocalNodeProfileExpectation{}
	}

	return m
}

//Return specifies results of invocation of LocalNode.LocalNodeProfile
func (m *mLocalNodeMockLocalNodeProfile) Return() *LocalNodeMock {
	m.mock.LocalNodeProfileFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &LocalNodeMockLocalNodeProfileExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of LocalNode.LocalNodeProfile is expected once
func (m *mLocalNodeMockLocalNodeProfile) ExpectOnce() *LocalNodeMockLocalNodeProfileExpectation {
	m.mock.LocalNodeProfileFunc = nil
	m.mainExpectation = nil

	expectation := &LocalNodeMockLocalNodeProfileExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of LocalNode.LocalNodeProfile method
func (m *mLocalNodeMockLocalNodeProfile) Set(f func()) *LocalNodeMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.LocalNodeProfileFunc = f
	return m.mock
}

//LocalNodeProfile implements github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode interface
func (m *LocalNodeMock) LocalNodeProfile() {
	counter := atomic.AddUint64(&m.LocalNodeProfilePreCounter, 1)
	defer atomic.AddUint64(&m.LocalNodeProfileCounter, 1)

	if len(m.LocalNodeProfileMock.expectationSeries) > 0 {
		if counter > uint64(len(m.LocalNodeProfileMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to LocalNodeMock.LocalNodeProfile.")
			return
		}

		return
	}

	if m.LocalNodeProfileMock.mainExpectation != nil {

		return
	}

	if m.LocalNodeProfileFunc == nil {
		m.t.Fatalf("Unexpected call to LocalNodeMock.LocalNodeProfile.")
		return
	}

	m.LocalNodeProfileFunc()
}

//LocalNodeProfileMinimockCounter returns a count of LocalNodeMock.LocalNodeProfileFunc invocations
func (m *LocalNodeMock) LocalNodeProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LocalNodeProfileCounter)
}

//LocalNodeProfileMinimockPreCounter returns the value of LocalNodeMock.LocalNodeProfile invocations
func (m *LocalNodeMock) LocalNodeProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LocalNodeProfilePreCounter)
}

//LocalNodeProfileFinished returns true if mock invocations count is ok
func (m *LocalNodeMock) LocalNodeProfileFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.LocalNodeProfileMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.LocalNodeProfileCounter) == uint64(len(m.LocalNodeProfileMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.LocalNodeProfileMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.LocalNodeProfileCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.LocalNodeProfileFunc != nil {
		return atomic.LoadUint64(&m.LocalNodeProfileCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LocalNodeMock) ValidateCallCounters() {

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetDeclaredPower")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetIndex")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetNodeID")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetOpMode")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetSignatureVerifier")
	}

	if !m.GetStaticFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetStatic")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.IsJoiner")
	}

	if !m.LocalNodeProfileFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.LocalNodeProfile")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *LocalNodeMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *LocalNodeMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *LocalNodeMock) MinimockFinish() {

	if !m.GetDeclaredPowerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetDeclaredPower")
	}

	if !m.GetIndexFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetIndex")
	}

	if !m.GetNodeIDFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetNodeID")
	}

	if !m.GetOpModeFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetOpMode")
	}

	if !m.GetSignatureVerifierFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetSignatureVerifier")
	}

	if !m.GetStaticFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.GetStatic")
	}

	if !m.IsJoinerFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.IsJoiner")
	}

	if !m.LocalNodeProfileFinished() {
		m.t.Fatal("Expected call to LocalNodeMock.LocalNodeProfile")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *LocalNodeMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *LocalNodeMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetDeclaredPowerFinished()
		ok = ok && m.GetIndexFinished()
		ok = ok && m.GetNodeIDFinished()
		ok = ok && m.GetOpModeFinished()
		ok = ok && m.GetSignatureVerifierFinished()
		ok = ok && m.GetStaticFinished()
		ok = ok && m.IsJoinerFinished()
		ok = ok && m.LocalNodeProfileFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetDeclaredPowerFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetDeclaredPower")
			}

			if !m.GetIndexFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetIndex")
			}

			if !m.GetNodeIDFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetNodeID")
			}

			if !m.GetOpModeFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetOpMode")
			}

			if !m.GetSignatureVerifierFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetSignatureVerifier")
			}

			if !m.GetStaticFinished() {
				m.t.Error("Expected call to LocalNodeMock.GetStatic")
			}

			if !m.IsJoinerFinished() {
				m.t.Error("Expected call to LocalNodeMock.IsJoiner")
			}

			if !m.LocalNodeProfileFinished() {
				m.t.Error("Expected call to LocalNodeMock.LocalNodeProfile")
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
func (m *LocalNodeMock) AllMocksCalled() bool {

	if !m.GetDeclaredPowerFinished() {
		return false
	}

	if !m.GetIndexFinished() {
		return false
	}

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

	if !m.IsJoinerFinished() {
		return false
	}

	if !m.LocalNodeProfileFinished() {
		return false
	}

	return true
}
