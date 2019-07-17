package endpoints

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Inbound" can be found in github.com/insolar/insolar/network/consensus/common/endpoints
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
)

//InboundMock implements github.com/insolar/insolar/network/consensus/common/endpoints.Inbound
type InboundMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mInboundMockAsByteString

	GetNameAddressFunc       func() (r Name)
	GetNameAddressCounter    uint64
	GetNameAddressPreCounter uint64
	GetNameAddressMock       mInboundMockGetNameAddress

	GetTransportCertFunc       func() (r cryptkit.CertificateHolder)
	GetTransportCertCounter    uint64
	GetTransportCertPreCounter uint64
	GetTransportCertMock       mInboundMockGetTransportCert

	GetTransportKeyFunc       func() (r cryptkit.SignatureKeyHolder)
	GetTransportKeyCounter    uint64
	GetTransportKeyPreCounter uint64
	GetTransportKeyMock       mInboundMockGetTransportKey
}

//NewInboundMock returns a mock for github.com/insolar/insolar/network/consensus/common/endpoints.Inbound
func NewInboundMock(t minimock.Tester) *InboundMock {
	m := &InboundMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mInboundMockAsByteString{mock: m}
	m.GetNameAddressMock = mInboundMockGetNameAddress{mock: m}
	m.GetTransportCertMock = mInboundMockGetTransportCert{mock: m}
	m.GetTransportKeyMock = mInboundMockGetTransportKey{mock: m}

	return m
}

type mInboundMockAsByteString struct {
	mock              *InboundMock
	mainExpectation   *InboundMockAsByteStringExpectation
	expectationSeries []*InboundMockAsByteStringExpectation
}

type InboundMockAsByteStringExpectation struct {
	result *InboundMockAsByteStringResult
}

type InboundMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of Inbound.AsByteString is expected from 1 to Infinity times
func (m *mInboundMockAsByteString) Expect() *mInboundMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of Inbound.AsByteString
func (m *mInboundMockAsByteString) Return(r string) *InboundMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &InboundMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Inbound.AsByteString is expected once
func (m *mInboundMockAsByteString) ExpectOnce() *InboundMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &InboundMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *InboundMockAsByteStringExpectation) Return(r string) {
	e.result = &InboundMockAsByteStringResult{r}
}

//Set uses given function f as a mock of Inbound.AsByteString method
func (m *mInboundMockAsByteString) Set(f func() (r string)) *InboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/common/endpoints.Inbound interface
func (m *InboundMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to InboundMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to InboundMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of InboundMock.AsByteStringFunc invocations
func (m *InboundMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of InboundMock.AsByteString invocations
func (m *InboundMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *InboundMock) AsByteStringFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.AsByteStringMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.AsByteStringCounter) == uint64(len(m.AsByteStringMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.AsByteStringMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.AsByteStringCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.AsByteStringFunc != nil {
		return atomic.LoadUint64(&m.AsByteStringCounter) > 0
	}

	return true
}

type mInboundMockGetNameAddress struct {
	mock              *InboundMock
	mainExpectation   *InboundMockGetNameAddressExpectation
	expectationSeries []*InboundMockGetNameAddressExpectation
}

type InboundMockGetNameAddressExpectation struct {
	result *InboundMockGetNameAddressResult
}

type InboundMockGetNameAddressResult struct {
	r Name
}

//Expect specifies that invocation of Inbound.GetNameAddress is expected from 1 to Infinity times
func (m *mInboundMockGetNameAddress) Expect() *mInboundMockGetNameAddress {
	m.mock.GetNameAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockGetNameAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of Inbound.GetNameAddress
func (m *mInboundMockGetNameAddress) Return(r Name) *InboundMock {
	m.mock.GetNameAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockGetNameAddressExpectation{}
	}
	m.mainExpectation.result = &InboundMockGetNameAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Inbound.GetNameAddress is expected once
func (m *mInboundMockGetNameAddress) ExpectOnce() *InboundMockGetNameAddressExpectation {
	m.mock.GetNameAddressFunc = nil
	m.mainExpectation = nil

	expectation := &InboundMockGetNameAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *InboundMockGetNameAddressExpectation) Return(r Name) {
	e.result = &InboundMockGetNameAddressResult{r}
}

//Set uses given function f as a mock of Inbound.GetNameAddress method
func (m *mInboundMockGetNameAddress) Set(f func() (r Name)) *InboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNameAddressFunc = f
	return m.mock
}

//GetNameAddress implements github.com/insolar/insolar/network/consensus/common/endpoints.Inbound interface
func (m *InboundMock) GetNameAddress() (r Name) {
	counter := atomic.AddUint64(&m.GetNameAddressPreCounter, 1)
	defer atomic.AddUint64(&m.GetNameAddressCounter, 1)

	if len(m.GetNameAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNameAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to InboundMock.GetNameAddress.")
			return
		}

		result := m.GetNameAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.GetNameAddress")
			return
		}

		r = result.r

		return
	}

	if m.GetNameAddressMock.mainExpectation != nil {

		result := m.GetNameAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.GetNameAddress")
		}

		r = result.r

		return
	}

	if m.GetNameAddressFunc == nil {
		m.t.Fatalf("Unexpected call to InboundMock.GetNameAddress.")
		return
	}

	return m.GetNameAddressFunc()
}

//GetNameAddressMinimockCounter returns a count of InboundMock.GetNameAddressFunc invocations
func (m *InboundMock) GetNameAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNameAddressCounter)
}

//GetNameAddressMinimockPreCounter returns the value of InboundMock.GetNameAddress invocations
func (m *InboundMock) GetNameAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNameAddressPreCounter)
}

//GetNameAddressFinished returns true if mock invocations count is ok
func (m *InboundMock) GetNameAddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetNameAddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetNameAddressCounter) == uint64(len(m.GetNameAddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetNameAddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetNameAddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetNameAddressFunc != nil {
		return atomic.LoadUint64(&m.GetNameAddressCounter) > 0
	}

	return true
}

type mInboundMockGetTransportCert struct {
	mock              *InboundMock
	mainExpectation   *InboundMockGetTransportCertExpectation
	expectationSeries []*InboundMockGetTransportCertExpectation
}

type InboundMockGetTransportCertExpectation struct {
	result *InboundMockGetTransportCertResult
}

type InboundMockGetTransportCertResult struct {
	r cryptkit.CertificateHolder
}

//Expect specifies that invocation of Inbound.GetTransportCert is expected from 1 to Infinity times
func (m *mInboundMockGetTransportCert) Expect() *mInboundMockGetTransportCert {
	m.mock.GetTransportCertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockGetTransportCertExpectation{}
	}

	return m
}

//Return specifies results of invocation of Inbound.GetTransportCert
func (m *mInboundMockGetTransportCert) Return(r cryptkit.CertificateHolder) *InboundMock {
	m.mock.GetTransportCertFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockGetTransportCertExpectation{}
	}
	m.mainExpectation.result = &InboundMockGetTransportCertResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Inbound.GetTransportCert is expected once
func (m *mInboundMockGetTransportCert) ExpectOnce() *InboundMockGetTransportCertExpectation {
	m.mock.GetTransportCertFunc = nil
	m.mainExpectation = nil

	expectation := &InboundMockGetTransportCertExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *InboundMockGetTransportCertExpectation) Return(r cryptkit.CertificateHolder) {
	e.result = &InboundMockGetTransportCertResult{r}
}

//Set uses given function f as a mock of Inbound.GetTransportCert method
func (m *mInboundMockGetTransportCert) Set(f func() (r cryptkit.CertificateHolder)) *InboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTransportCertFunc = f
	return m.mock
}

//GetTransportCert implements github.com/insolar/insolar/network/consensus/common/endpoints.Inbound interface
func (m *InboundMock) GetTransportCert() (r cryptkit.CertificateHolder) {
	counter := atomic.AddUint64(&m.GetTransportCertPreCounter, 1)
	defer atomic.AddUint64(&m.GetTransportCertCounter, 1)

	if len(m.GetTransportCertMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTransportCertMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to InboundMock.GetTransportCert.")
			return
		}

		result := m.GetTransportCertMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.GetTransportCert")
			return
		}

		r = result.r

		return
	}

	if m.GetTransportCertMock.mainExpectation != nil {

		result := m.GetTransportCertMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.GetTransportCert")
		}

		r = result.r

		return
	}

	if m.GetTransportCertFunc == nil {
		m.t.Fatalf("Unexpected call to InboundMock.GetTransportCert.")
		return
	}

	return m.GetTransportCertFunc()
}

//GetTransportCertMinimockCounter returns a count of InboundMock.GetTransportCertFunc invocations
func (m *InboundMock) GetTransportCertMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransportCertCounter)
}

//GetTransportCertMinimockPreCounter returns the value of InboundMock.GetTransportCert invocations
func (m *InboundMock) GetTransportCertMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransportCertPreCounter)
}

//GetTransportCertFinished returns true if mock invocations count is ok
func (m *InboundMock) GetTransportCertFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetTransportCertMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetTransportCertCounter) == uint64(len(m.GetTransportCertMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetTransportCertMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetTransportCertCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetTransportCertFunc != nil {
		return atomic.LoadUint64(&m.GetTransportCertCounter) > 0
	}

	return true
}

type mInboundMockGetTransportKey struct {
	mock              *InboundMock
	mainExpectation   *InboundMockGetTransportKeyExpectation
	expectationSeries []*InboundMockGetTransportKeyExpectation
}

type InboundMockGetTransportKeyExpectation struct {
	result *InboundMockGetTransportKeyResult
}

type InboundMockGetTransportKeyResult struct {
	r cryptkit.SignatureKeyHolder
}

//Expect specifies that invocation of Inbound.GetTransportKey is expected from 1 to Infinity times
func (m *mInboundMockGetTransportKey) Expect() *mInboundMockGetTransportKey {
	m.mock.GetTransportKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockGetTransportKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of Inbound.GetTransportKey
func (m *mInboundMockGetTransportKey) Return(r cryptkit.SignatureKeyHolder) *InboundMock {
	m.mock.GetTransportKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &InboundMockGetTransportKeyExpectation{}
	}
	m.mainExpectation.result = &InboundMockGetTransportKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Inbound.GetTransportKey is expected once
func (m *mInboundMockGetTransportKey) ExpectOnce() *InboundMockGetTransportKeyExpectation {
	m.mock.GetTransportKeyFunc = nil
	m.mainExpectation = nil

	expectation := &InboundMockGetTransportKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *InboundMockGetTransportKeyExpectation) Return(r cryptkit.SignatureKeyHolder) {
	e.result = &InboundMockGetTransportKeyResult{r}
}

//Set uses given function f as a mock of Inbound.GetTransportKey method
func (m *mInboundMockGetTransportKey) Set(f func() (r cryptkit.SignatureKeyHolder)) *InboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTransportKeyFunc = f
	return m.mock
}

//GetTransportKey implements github.com/insolar/insolar/network/consensus/common/endpoints.Inbound interface
func (m *InboundMock) GetTransportKey() (r cryptkit.SignatureKeyHolder) {
	counter := atomic.AddUint64(&m.GetTransportKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetTransportKeyCounter, 1)

	if len(m.GetTransportKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTransportKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to InboundMock.GetTransportKey.")
			return
		}

		result := m.GetTransportKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.GetTransportKey")
			return
		}

		r = result.r

		return
	}

	if m.GetTransportKeyMock.mainExpectation != nil {

		result := m.GetTransportKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the InboundMock.GetTransportKey")
		}

		r = result.r

		return
	}

	if m.GetTransportKeyFunc == nil {
		m.t.Fatalf("Unexpected call to InboundMock.GetTransportKey.")
		return
	}

	return m.GetTransportKeyFunc()
}

//GetTransportKeyMinimockCounter returns a count of InboundMock.GetTransportKeyFunc invocations
func (m *InboundMock) GetTransportKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransportKeyCounter)
}

//GetTransportKeyMinimockPreCounter returns the value of InboundMock.GetTransportKey invocations
func (m *InboundMock) GetTransportKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTransportKeyPreCounter)
}

//GetTransportKeyFinished returns true if mock invocations count is ok
func (m *InboundMock) GetTransportKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetTransportKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetTransportKeyCounter) == uint64(len(m.GetTransportKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetTransportKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetTransportKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetTransportKeyFunc != nil {
		return atomic.LoadUint64(&m.GetTransportKeyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *InboundMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to InboundMock.AsByteString")
	}

	if !m.GetNameAddressFinished() {
		m.t.Fatal("Expected call to InboundMock.GetNameAddress")
	}

	if !m.GetTransportCertFinished() {
		m.t.Fatal("Expected call to InboundMock.GetTransportCert")
	}

	if !m.GetTransportKeyFinished() {
		m.t.Fatal("Expected call to InboundMock.GetTransportKey")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *InboundMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *InboundMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *InboundMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to InboundMock.AsByteString")
	}

	if !m.GetNameAddressFinished() {
		m.t.Fatal("Expected call to InboundMock.GetNameAddress")
	}

	if !m.GetTransportCertFinished() {
		m.t.Fatal("Expected call to InboundMock.GetTransportCert")
	}

	if !m.GetTransportKeyFinished() {
		m.t.Fatal("Expected call to InboundMock.GetTransportKey")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *InboundMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *InboundMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AsByteStringFinished()
		ok = ok && m.GetNameAddressFinished()
		ok = ok && m.GetTransportCertFinished()
		ok = ok && m.GetTransportKeyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AsByteStringFinished() {
				m.t.Error("Expected call to InboundMock.AsByteString")
			}

			if !m.GetNameAddressFinished() {
				m.t.Error("Expected call to InboundMock.GetNameAddress")
			}

			if !m.GetTransportCertFinished() {
				m.t.Error("Expected call to InboundMock.GetTransportCert")
			}

			if !m.GetTransportKeyFinished() {
				m.t.Error("Expected call to InboundMock.GetTransportKey")
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
func (m *InboundMock) AllMocksCalled() bool {

	if !m.AsByteStringFinished() {
		return false
	}

	if !m.GetNameAddressFinished() {
		return false
	}

	if !m.GetTransportCertFinished() {
		return false
	}

	if !m.GetTransportKeyFinished() {
		return false
	}

	return true
}
