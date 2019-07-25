package endpoints

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Outbound" can be found in github.com/insolar/insolar/network/consensus/common/endpoints
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//OutboundMock implements github.com/insolar/insolar/network/consensus/common/endpoints.Outbound
type OutboundMock struct {
	t minimock.Tester

	AsByteStringFunc       func() (r string)
	AsByteStringCounter    uint64
	AsByteStringPreCounter uint64
	AsByteStringMock       mOutboundMockAsByteString

	CanAcceptFunc       func(p Inbound) (r bool)
	CanAcceptCounter    uint64
	CanAcceptPreCounter uint64
	CanAcceptMock       mOutboundMockCanAccept

	GetEndpointTypeFunc       func() (r NodeEndpointType)
	GetEndpointTypeCounter    uint64
	GetEndpointTypePreCounter uint64
	GetEndpointTypeMock       mOutboundMockGetEndpointType

	GetIPAddressFunc       func() (r IPAddress)
	GetIPAddressCounter    uint64
	GetIPAddressPreCounter uint64
	GetIPAddressMock       mOutboundMockGetIPAddress

	GetNameAddressFunc       func() (r Name)
	GetNameAddressCounter    uint64
	GetNameAddressPreCounter uint64
	GetNameAddressMock       mOutboundMockGetNameAddress

	GetRelayIDFunc       func() (r insolar.ShortNodeID)
	GetRelayIDCounter    uint64
	GetRelayIDPreCounter uint64
	GetRelayIDMock       mOutboundMockGetRelayID
}

//NewOutboundMock returns a mock for github.com/insolar/insolar/network/consensus/common/endpoints.Outbound
func NewOutboundMock(t minimock.Tester) *OutboundMock {
	m := &OutboundMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mOutboundMockAsByteString{mock: m}
	m.CanAcceptMock = mOutboundMockCanAccept{mock: m}
	m.GetEndpointTypeMock = mOutboundMockGetEndpointType{mock: m}
	m.GetIPAddressMock = mOutboundMockGetIPAddress{mock: m}
	m.GetNameAddressMock = mOutboundMockGetNameAddress{mock: m}
	m.GetRelayIDMock = mOutboundMockGetRelayID{mock: m}

	return m
}

type mOutboundMockAsByteString struct {
	mock              *OutboundMock
	mainExpectation   *OutboundMockAsByteStringExpectation
	expectationSeries []*OutboundMockAsByteStringExpectation
}

type OutboundMockAsByteStringExpectation struct {
	result *OutboundMockAsByteStringResult
}

type OutboundMockAsByteStringResult struct {
	r string
}

//Expect specifies that invocation of Outbound.AsByteString is expected from 1 to Infinity times
func (m *mOutboundMockAsByteString) Expect() *mOutboundMockAsByteString {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockAsByteStringExpectation{}
	}

	return m
}

//Return specifies results of invocation of Outbound.AsByteString
func (m *mOutboundMockAsByteString) Return(r string) *OutboundMock {
	m.mock.AsByteStringFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockAsByteStringExpectation{}
	}
	m.mainExpectation.result = &OutboundMockAsByteStringResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Outbound.AsByteString is expected once
func (m *mOutboundMockAsByteString) ExpectOnce() *OutboundMockAsByteStringExpectation {
	m.mock.AsByteStringFunc = nil
	m.mainExpectation = nil

	expectation := &OutboundMockAsByteStringExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OutboundMockAsByteStringExpectation) Return(r string) {
	e.result = &OutboundMockAsByteStringResult{r}
}

//Set uses given function f as a mock of Outbound.AsByteString method
func (m *mOutboundMockAsByteString) Set(f func() (r string)) *OutboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.AsByteStringFunc = f
	return m.mock
}

//AsByteString implements github.com/insolar/insolar/network/consensus/common/endpoints.Outbound interface
func (m *OutboundMock) AsByteString() (r string) {
	counter := atomic.AddUint64(&m.AsByteStringPreCounter, 1)
	defer atomic.AddUint64(&m.AsByteStringCounter, 1)

	if len(m.AsByteStringMock.expectationSeries) > 0 {
		if counter > uint64(len(m.AsByteStringMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OutboundMock.AsByteString.")
			return
		}

		result := m.AsByteStringMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.AsByteString")
			return
		}

		r = result.r

		return
	}

	if m.AsByteStringMock.mainExpectation != nil {

		result := m.AsByteStringMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.AsByteString")
		}

		r = result.r

		return
	}

	if m.AsByteStringFunc == nil {
		m.t.Fatalf("Unexpected call to OutboundMock.AsByteString.")
		return
	}

	return m.AsByteStringFunc()
}

//AsByteStringMinimockCounter returns a count of OutboundMock.AsByteStringFunc invocations
func (m *OutboundMock) AsByteStringMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringCounter)
}

//AsByteStringMinimockPreCounter returns the value of OutboundMock.AsByteString invocations
func (m *OutboundMock) AsByteStringMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.AsByteStringPreCounter)
}

//AsByteStringFinished returns true if mock invocations count is ok
func (m *OutboundMock) AsByteStringFinished() bool {
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

type mOutboundMockCanAccept struct {
	mock              *OutboundMock
	mainExpectation   *OutboundMockCanAcceptExpectation
	expectationSeries []*OutboundMockCanAcceptExpectation
}

type OutboundMockCanAcceptExpectation struct {
	input  *OutboundMockCanAcceptInput
	result *OutboundMockCanAcceptResult
}

type OutboundMockCanAcceptInput struct {
	p Inbound
}

type OutboundMockCanAcceptResult struct {
	r bool
}

//Expect specifies that invocation of Outbound.CanAccept is expected from 1 to Infinity times
func (m *mOutboundMockCanAccept) Expect(p Inbound) *mOutboundMockCanAccept {
	m.mock.CanAcceptFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockCanAcceptExpectation{}
	}
	m.mainExpectation.input = &OutboundMockCanAcceptInput{p}
	return m
}

//Return specifies results of invocation of Outbound.CanAccept
func (m *mOutboundMockCanAccept) Return(r bool) *OutboundMock {
	m.mock.CanAcceptFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockCanAcceptExpectation{}
	}
	m.mainExpectation.result = &OutboundMockCanAcceptResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Outbound.CanAccept is expected once
func (m *mOutboundMockCanAccept) ExpectOnce(p Inbound) *OutboundMockCanAcceptExpectation {
	m.mock.CanAcceptFunc = nil
	m.mainExpectation = nil

	expectation := &OutboundMockCanAcceptExpectation{}
	expectation.input = &OutboundMockCanAcceptInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OutboundMockCanAcceptExpectation) Return(r bool) {
	e.result = &OutboundMockCanAcceptResult{r}
}

//Set uses given function f as a mock of Outbound.CanAccept method
func (m *mOutboundMockCanAccept) Set(f func(p Inbound) (r bool)) *OutboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CanAcceptFunc = f
	return m.mock
}

//CanAccept implements github.com/insolar/insolar/network/consensus/common/endpoints.Outbound interface
func (m *OutboundMock) CanAccept(p Inbound) (r bool) {
	counter := atomic.AddUint64(&m.CanAcceptPreCounter, 1)
	defer atomic.AddUint64(&m.CanAcceptCounter, 1)

	if len(m.CanAcceptMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CanAcceptMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OutboundMock.CanAccept. %v", p)
			return
		}

		input := m.CanAcceptMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, OutboundMockCanAcceptInput{p}, "Outbound.CanAccept got unexpected parameters")

		result := m.CanAcceptMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.CanAccept")
			return
		}

		r = result.r

		return
	}

	if m.CanAcceptMock.mainExpectation != nil {

		input := m.CanAcceptMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, OutboundMockCanAcceptInput{p}, "Outbound.CanAccept got unexpected parameters")
		}

		result := m.CanAcceptMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.CanAccept")
		}

		r = result.r

		return
	}

	if m.CanAcceptFunc == nil {
		m.t.Fatalf("Unexpected call to OutboundMock.CanAccept. %v", p)
		return
	}

	return m.CanAcceptFunc(p)
}

//CanAcceptMinimockCounter returns a count of OutboundMock.CanAcceptFunc invocations
func (m *OutboundMock) CanAcceptMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CanAcceptCounter)
}

//CanAcceptMinimockPreCounter returns the value of OutboundMock.CanAccept invocations
func (m *OutboundMock) CanAcceptMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CanAcceptPreCounter)
}

//CanAcceptFinished returns true if mock invocations count is ok
func (m *OutboundMock) CanAcceptFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CanAcceptMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CanAcceptCounter) == uint64(len(m.CanAcceptMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CanAcceptMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CanAcceptCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CanAcceptFunc != nil {
		return atomic.LoadUint64(&m.CanAcceptCounter) > 0
	}

	return true
}

type mOutboundMockGetEndpointType struct {
	mock              *OutboundMock
	mainExpectation   *OutboundMockGetEndpointTypeExpectation
	expectationSeries []*OutboundMockGetEndpointTypeExpectation
}

type OutboundMockGetEndpointTypeExpectation struct {
	result *OutboundMockGetEndpointTypeResult
}

type OutboundMockGetEndpointTypeResult struct {
	r NodeEndpointType
}

//Expect specifies that invocation of Outbound.GetEndpointType is expected from 1 to Infinity times
func (m *mOutboundMockGetEndpointType) Expect() *mOutboundMockGetEndpointType {
	m.mock.GetEndpointTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetEndpointTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of Outbound.GetEndpointType
func (m *mOutboundMockGetEndpointType) Return(r NodeEndpointType) *OutboundMock {
	m.mock.GetEndpointTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetEndpointTypeExpectation{}
	}
	m.mainExpectation.result = &OutboundMockGetEndpointTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Outbound.GetEndpointType is expected once
func (m *mOutboundMockGetEndpointType) ExpectOnce() *OutboundMockGetEndpointTypeExpectation {
	m.mock.GetEndpointTypeFunc = nil
	m.mainExpectation = nil

	expectation := &OutboundMockGetEndpointTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OutboundMockGetEndpointTypeExpectation) Return(r NodeEndpointType) {
	e.result = &OutboundMockGetEndpointTypeResult{r}
}

//Set uses given function f as a mock of Outbound.GetEndpointType method
func (m *mOutboundMockGetEndpointType) Set(f func() (r NodeEndpointType)) *OutboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetEndpointTypeFunc = f
	return m.mock
}

//GetEndpointType implements github.com/insolar/insolar/network/consensus/common/endpoints.Outbound interface
func (m *OutboundMock) GetEndpointType() (r NodeEndpointType) {
	counter := atomic.AddUint64(&m.GetEndpointTypePreCounter, 1)
	defer atomic.AddUint64(&m.GetEndpointTypeCounter, 1)

	if len(m.GetEndpointTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetEndpointTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OutboundMock.GetEndpointType.")
			return
		}

		result := m.GetEndpointTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetEndpointType")
			return
		}

		r = result.r

		return
	}

	if m.GetEndpointTypeMock.mainExpectation != nil {

		result := m.GetEndpointTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetEndpointType")
		}

		r = result.r

		return
	}

	if m.GetEndpointTypeFunc == nil {
		m.t.Fatalf("Unexpected call to OutboundMock.GetEndpointType.")
		return
	}

	return m.GetEndpointTypeFunc()
}

//GetEndpointTypeMinimockCounter returns a count of OutboundMock.GetEndpointTypeFunc invocations
func (m *OutboundMock) GetEndpointTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetEndpointTypeCounter)
}

//GetEndpointTypeMinimockPreCounter returns the value of OutboundMock.GetEndpointType invocations
func (m *OutboundMock) GetEndpointTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetEndpointTypePreCounter)
}

//GetEndpointTypeFinished returns true if mock invocations count is ok
func (m *OutboundMock) GetEndpointTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetEndpointTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetEndpointTypeCounter) == uint64(len(m.GetEndpointTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetEndpointTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetEndpointTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetEndpointTypeFunc != nil {
		return atomic.LoadUint64(&m.GetEndpointTypeCounter) > 0
	}

	return true
}

type mOutboundMockGetIPAddress struct {
	mock              *OutboundMock
	mainExpectation   *OutboundMockGetIPAddressExpectation
	expectationSeries []*OutboundMockGetIPAddressExpectation
}

type OutboundMockGetIPAddressExpectation struct {
	result *OutboundMockGetIPAddressResult
}

type OutboundMockGetIPAddressResult struct {
	r IPAddress
}

//Expect specifies that invocation of Outbound.GetIPAddress is expected from 1 to Infinity times
func (m *mOutboundMockGetIPAddress) Expect() *mOutboundMockGetIPAddress {
	m.mock.GetIPAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetIPAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of Outbound.GetIPAddress
func (m *mOutboundMockGetIPAddress) Return(r IPAddress) *OutboundMock {
	m.mock.GetIPAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetIPAddressExpectation{}
	}
	m.mainExpectation.result = &OutboundMockGetIPAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Outbound.GetIPAddress is expected once
func (m *mOutboundMockGetIPAddress) ExpectOnce() *OutboundMockGetIPAddressExpectation {
	m.mock.GetIPAddressFunc = nil
	m.mainExpectation = nil

	expectation := &OutboundMockGetIPAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OutboundMockGetIPAddressExpectation) Return(r IPAddress) {
	e.result = &OutboundMockGetIPAddressResult{r}
}

//Set uses given function f as a mock of Outbound.GetIPAddress method
func (m *mOutboundMockGetIPAddress) Set(f func() (r IPAddress)) *OutboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetIPAddressFunc = f
	return m.mock
}

//GetIPAddress implements github.com/insolar/insolar/network/consensus/common/endpoints.Outbound interface
func (m *OutboundMock) GetIPAddress() (r IPAddress) {
	counter := atomic.AddUint64(&m.GetIPAddressPreCounter, 1)
	defer atomic.AddUint64(&m.GetIPAddressCounter, 1)

	if len(m.GetIPAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetIPAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OutboundMock.GetIPAddress.")
			return
		}

		result := m.GetIPAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetIPAddress")
			return
		}

		r = result.r

		return
	}

	if m.GetIPAddressMock.mainExpectation != nil {

		result := m.GetIPAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetIPAddress")
		}

		r = result.r

		return
	}

	if m.GetIPAddressFunc == nil {
		m.t.Fatalf("Unexpected call to OutboundMock.GetIPAddress.")
		return
	}

	return m.GetIPAddressFunc()
}

//GetIPAddressMinimockCounter returns a count of OutboundMock.GetIPAddressFunc invocations
func (m *OutboundMock) GetIPAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetIPAddressCounter)
}

//GetIPAddressMinimockPreCounter returns the value of OutboundMock.GetIPAddress invocations
func (m *OutboundMock) GetIPAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetIPAddressPreCounter)
}

//GetIPAddressFinished returns true if mock invocations count is ok
func (m *OutboundMock) GetIPAddressFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetIPAddressMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetIPAddressCounter) == uint64(len(m.GetIPAddressMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetIPAddressMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetIPAddressCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetIPAddressFunc != nil {
		return atomic.LoadUint64(&m.GetIPAddressCounter) > 0
	}

	return true
}

type mOutboundMockGetNameAddress struct {
	mock              *OutboundMock
	mainExpectation   *OutboundMockGetNameAddressExpectation
	expectationSeries []*OutboundMockGetNameAddressExpectation
}

type OutboundMockGetNameAddressExpectation struct {
	result *OutboundMockGetNameAddressResult
}

type OutboundMockGetNameAddressResult struct {
	r Name
}

//Expect specifies that invocation of Outbound.GetNameAddress is expected from 1 to Infinity times
func (m *mOutboundMockGetNameAddress) Expect() *mOutboundMockGetNameAddress {
	m.mock.GetNameAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetNameAddressExpectation{}
	}

	return m
}

//Return specifies results of invocation of Outbound.GetNameAddress
func (m *mOutboundMockGetNameAddress) Return(r Name) *OutboundMock {
	m.mock.GetNameAddressFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetNameAddressExpectation{}
	}
	m.mainExpectation.result = &OutboundMockGetNameAddressResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Outbound.GetNameAddress is expected once
func (m *mOutboundMockGetNameAddress) ExpectOnce() *OutboundMockGetNameAddressExpectation {
	m.mock.GetNameAddressFunc = nil
	m.mainExpectation = nil

	expectation := &OutboundMockGetNameAddressExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OutboundMockGetNameAddressExpectation) Return(r Name) {
	e.result = &OutboundMockGetNameAddressResult{r}
}

//Set uses given function f as a mock of Outbound.GetNameAddress method
func (m *mOutboundMockGetNameAddress) Set(f func() (r Name)) *OutboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetNameAddressFunc = f
	return m.mock
}

//GetNameAddress implements github.com/insolar/insolar/network/consensus/common/endpoints.Outbound interface
func (m *OutboundMock) GetNameAddress() (r Name) {
	counter := atomic.AddUint64(&m.GetNameAddressPreCounter, 1)
	defer atomic.AddUint64(&m.GetNameAddressCounter, 1)

	if len(m.GetNameAddressMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetNameAddressMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OutboundMock.GetNameAddress.")
			return
		}

		result := m.GetNameAddressMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetNameAddress")
			return
		}

		r = result.r

		return
	}

	if m.GetNameAddressMock.mainExpectation != nil {

		result := m.GetNameAddressMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetNameAddress")
		}

		r = result.r

		return
	}

	if m.GetNameAddressFunc == nil {
		m.t.Fatalf("Unexpected call to OutboundMock.GetNameAddress.")
		return
	}

	return m.GetNameAddressFunc()
}

//GetNameAddressMinimockCounter returns a count of OutboundMock.GetNameAddressFunc invocations
func (m *OutboundMock) GetNameAddressMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetNameAddressCounter)
}

//GetNameAddressMinimockPreCounter returns the value of OutboundMock.GetNameAddress invocations
func (m *OutboundMock) GetNameAddressMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetNameAddressPreCounter)
}

//GetNameAddressFinished returns true if mock invocations count is ok
func (m *OutboundMock) GetNameAddressFinished() bool {
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

type mOutboundMockGetRelayID struct {
	mock              *OutboundMock
	mainExpectation   *OutboundMockGetRelayIDExpectation
	expectationSeries []*OutboundMockGetRelayIDExpectation
}

type OutboundMockGetRelayIDExpectation struct {
	result *OutboundMockGetRelayIDResult
}

type OutboundMockGetRelayIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of Outbound.GetRelayID is expected from 1 to Infinity times
func (m *mOutboundMockGetRelayID) Expect() *mOutboundMockGetRelayID {
	m.mock.GetRelayIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetRelayIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of Outbound.GetRelayID
func (m *mOutboundMockGetRelayID) Return(r insolar.ShortNodeID) *OutboundMock {
	m.mock.GetRelayIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &OutboundMockGetRelayIDExpectation{}
	}
	m.mainExpectation.result = &OutboundMockGetRelayIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of Outbound.GetRelayID is expected once
func (m *mOutboundMockGetRelayID) ExpectOnce() *OutboundMockGetRelayIDExpectation {
	m.mock.GetRelayIDFunc = nil
	m.mainExpectation = nil

	expectation := &OutboundMockGetRelayIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *OutboundMockGetRelayIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &OutboundMockGetRelayIDResult{r}
}

//Set uses given function f as a mock of Outbound.GetRelayID method
func (m *mOutboundMockGetRelayID) Set(f func() (r insolar.ShortNodeID)) *OutboundMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetRelayIDFunc = f
	return m.mock
}

//GetRelayID implements github.com/insolar/insolar/network/consensus/common/endpoints.Outbound interface
func (m *OutboundMock) GetRelayID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetRelayIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetRelayIDCounter, 1)

	if len(m.GetRelayIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetRelayIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to OutboundMock.GetRelayID.")
			return
		}

		result := m.GetRelayIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetRelayID")
			return
		}

		r = result.r

		return
	}

	if m.GetRelayIDMock.mainExpectation != nil {

		result := m.GetRelayIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the OutboundMock.GetRelayID")
		}

		r = result.r

		return
	}

	if m.GetRelayIDFunc == nil {
		m.t.Fatalf("Unexpected call to OutboundMock.GetRelayID.")
		return
	}

	return m.GetRelayIDFunc()
}

//GetRelayIDMinimockCounter returns a count of OutboundMock.GetRelayIDFunc invocations
func (m *OutboundMock) GetRelayIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetRelayIDCounter)
}

//GetRelayIDMinimockPreCounter returns the value of OutboundMock.GetRelayID invocations
func (m *OutboundMock) GetRelayIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetRelayIDPreCounter)
}

//GetRelayIDFinished returns true if mock invocations count is ok
func (m *OutboundMock) GetRelayIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetRelayIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetRelayIDCounter) == uint64(len(m.GetRelayIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetRelayIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetRelayIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetRelayIDFunc != nil {
		return atomic.LoadUint64(&m.GetRelayIDCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *OutboundMock) ValidateCallCounters() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to OutboundMock.AsByteString")
	}

	if !m.CanAcceptFinished() {
		m.t.Fatal("Expected call to OutboundMock.CanAccept")
	}

	if !m.GetEndpointTypeFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetEndpointType")
	}

	if !m.GetIPAddressFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetIPAddress")
	}

	if !m.GetNameAddressFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetNameAddress")
	}

	if !m.GetRelayIDFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetRelayID")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *OutboundMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *OutboundMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *OutboundMock) MinimockFinish() {

	if !m.AsByteStringFinished() {
		m.t.Fatal("Expected call to OutboundMock.AsByteString")
	}

	if !m.CanAcceptFinished() {
		m.t.Fatal("Expected call to OutboundMock.CanAccept")
	}

	if !m.GetEndpointTypeFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetEndpointType")
	}

	if !m.GetIPAddressFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetIPAddress")
	}

	if !m.GetNameAddressFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetNameAddress")
	}

	if !m.GetRelayIDFinished() {
		m.t.Fatal("Expected call to OutboundMock.GetRelayID")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *OutboundMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *OutboundMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.AsByteStringFinished()
		ok = ok && m.CanAcceptFinished()
		ok = ok && m.GetEndpointTypeFinished()
		ok = ok && m.GetIPAddressFinished()
		ok = ok && m.GetNameAddressFinished()
		ok = ok && m.GetRelayIDFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.AsByteStringFinished() {
				m.t.Error("Expected call to OutboundMock.AsByteString")
			}

			if !m.CanAcceptFinished() {
				m.t.Error("Expected call to OutboundMock.CanAccept")
			}

			if !m.GetEndpointTypeFinished() {
				m.t.Error("Expected call to OutboundMock.GetEndpointType")
			}

			if !m.GetIPAddressFinished() {
				m.t.Error("Expected call to OutboundMock.GetIPAddress")
			}

			if !m.GetNameAddressFinished() {
				m.t.Error("Expected call to OutboundMock.GetNameAddress")
			}

			if !m.GetRelayIDFinished() {
				m.t.Error("Expected call to OutboundMock.GetRelayID")
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
func (m *OutboundMock) AllMocksCalled() bool {

	if !m.AsByteStringFinished() {
		return false
	}

	if !m.CanAcceptFinished() {
		return false
	}

	if !m.GetEndpointTypeFinished() {
		return false
	}

	if !m.GetIPAddressFinished() {
		return false
	}

	if !m.GetNameAddressFinished() {
		return false
	}

	if !m.GetRelayIDFinished() {
		return false
	}

	return true
}
