package transport

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "PacketParser" can be found in github.com/insolar/insolar/network/consensus/gcpv2/api/transport
*/
import (
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"
	cryptkit "github.com/insolar/insolar/network/consensus/common/cryptkit"
	pulse "github.com/insolar/insolar/network/consensus/common/pulse"
	phases "github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
)

//PacketParserMock implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser
type PacketParserMock struct {
	t minimock.Tester

	GetMemberPacketFunc       func() (r MemberPacketReader)
	GetMemberPacketCounter    uint64
	GetMemberPacketPreCounter uint64
	GetMemberPacketMock       mPacketParserMockGetMemberPacket

	GetPacketSignatureFunc       func() (r cryptkit.SignedDigest)
	GetPacketSignatureCounter    uint64
	GetPacketSignaturePreCounter uint64
	GetPacketSignatureMock       mPacketParserMockGetPacketSignature

	GetPacketTypeFunc       func() (r phases.PacketType)
	GetPacketTypeCounter    uint64
	GetPacketTypePreCounter uint64
	GetPacketTypeMock       mPacketParserMockGetPacketType

	GetPulseNumberFunc       func() (r pulse.Number)
	GetPulseNumberCounter    uint64
	GetPulseNumberPreCounter uint64
	GetPulseNumberMock       mPacketParserMockGetPulseNumber

	GetPulsePacketFunc       func() (r PulsePacketReader)
	GetPulsePacketCounter    uint64
	GetPulsePacketPreCounter uint64
	GetPulsePacketMock       mPacketParserMockGetPulsePacket

	GetReceiverIDFunc       func() (r insolar.ShortNodeID)
	GetReceiverIDCounter    uint64
	GetReceiverIDPreCounter uint64
	GetReceiverIDMock       mPacketParserMockGetReceiverID

	GetSourceIDFunc       func() (r insolar.ShortNodeID)
	GetSourceIDCounter    uint64
	GetSourceIDPreCounter uint64
	GetSourceIDMock       mPacketParserMockGetSourceID

	GetTargetIDFunc       func() (r insolar.ShortNodeID)
	GetTargetIDCounter    uint64
	GetTargetIDPreCounter uint64
	GetTargetIDMock       mPacketParserMockGetTargetID

	IsRelayForbiddenFunc       func() (r bool)
	IsRelayForbiddenCounter    uint64
	IsRelayForbiddenPreCounter uint64
	IsRelayForbiddenMock       mPacketParserMockIsRelayForbidden

	ParsePacketBodyFunc       func() (r PacketParser, r1 error)
	ParsePacketBodyCounter    uint64
	ParsePacketBodyPreCounter uint64
	ParsePacketBodyMock       mPacketParserMockParsePacketBody
}

//NewPacketParserMock returns a mock for github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser
func NewPacketParserMock(t minimock.Tester) *PacketParserMock {
	m := &PacketParserMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetMemberPacketMock = mPacketParserMockGetMemberPacket{mock: m}
	m.GetPacketSignatureMock = mPacketParserMockGetPacketSignature{mock: m}
	m.GetPacketTypeMock = mPacketParserMockGetPacketType{mock: m}
	m.GetPulseNumberMock = mPacketParserMockGetPulseNumber{mock: m}
	m.GetPulsePacketMock = mPacketParserMockGetPulsePacket{mock: m}
	m.GetReceiverIDMock = mPacketParserMockGetReceiverID{mock: m}
	m.GetSourceIDMock = mPacketParserMockGetSourceID{mock: m}
	m.GetTargetIDMock = mPacketParserMockGetTargetID{mock: m}
	m.IsRelayForbiddenMock = mPacketParserMockIsRelayForbidden{mock: m}
	m.ParsePacketBodyMock = mPacketParserMockParsePacketBody{mock: m}

	return m
}

type mPacketParserMockGetMemberPacket struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetMemberPacketExpectation
	expectationSeries []*PacketParserMockGetMemberPacketExpectation
}

type PacketParserMockGetMemberPacketExpectation struct {
	result *PacketParserMockGetMemberPacketResult
}

type PacketParserMockGetMemberPacketResult struct {
	r MemberPacketReader
}

//Expect specifies that invocation of PacketParser.GetMemberPacket is expected from 1 to Infinity times
func (m *mPacketParserMockGetMemberPacket) Expect() *mPacketParserMockGetMemberPacket {
	m.mock.GetMemberPacketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetMemberPacketExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetMemberPacket
func (m *mPacketParserMockGetMemberPacket) Return(r MemberPacketReader) *PacketParserMock {
	m.mock.GetMemberPacketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetMemberPacketExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetMemberPacketResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetMemberPacket is expected once
func (m *mPacketParserMockGetMemberPacket) ExpectOnce() *PacketParserMockGetMemberPacketExpectation {
	m.mock.GetMemberPacketFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetMemberPacketExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetMemberPacketExpectation) Return(r MemberPacketReader) {
	e.result = &PacketParserMockGetMemberPacketResult{r}
}

//Set uses given function f as a mock of PacketParser.GetMemberPacket method
func (m *mPacketParserMockGetMemberPacket) Set(f func() (r MemberPacketReader)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetMemberPacketFunc = f
	return m.mock
}

//GetMemberPacket implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetMemberPacket() (r MemberPacketReader) {
	counter := atomic.AddUint64(&m.GetMemberPacketPreCounter, 1)
	defer atomic.AddUint64(&m.GetMemberPacketCounter, 1)

	if len(m.GetMemberPacketMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetMemberPacketMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetMemberPacket.")
			return
		}

		result := m.GetMemberPacketMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetMemberPacket")
			return
		}

		r = result.r

		return
	}

	if m.GetMemberPacketMock.mainExpectation != nil {

		result := m.GetMemberPacketMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetMemberPacket")
		}

		r = result.r

		return
	}

	if m.GetMemberPacketFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetMemberPacket.")
		return
	}

	return m.GetMemberPacketFunc()
}

//GetMemberPacketMinimockCounter returns a count of PacketParserMock.GetMemberPacketFunc invocations
func (m *PacketParserMock) GetMemberPacketMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetMemberPacketCounter)
}

//GetMemberPacketMinimockPreCounter returns the value of PacketParserMock.GetMemberPacket invocations
func (m *PacketParserMock) GetMemberPacketMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetMemberPacketPreCounter)
}

//GetMemberPacketFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetMemberPacketFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetMemberPacketMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetMemberPacketCounter) == uint64(len(m.GetMemberPacketMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetMemberPacketMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetMemberPacketCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetMemberPacketFunc != nil {
		return atomic.LoadUint64(&m.GetMemberPacketCounter) > 0
	}

	return true
}

type mPacketParserMockGetPacketSignature struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetPacketSignatureExpectation
	expectationSeries []*PacketParserMockGetPacketSignatureExpectation
}

type PacketParserMockGetPacketSignatureExpectation struct {
	result *PacketParserMockGetPacketSignatureResult
}

type PacketParserMockGetPacketSignatureResult struct {
	r cryptkit.SignedDigest
}

//Expect specifies that invocation of PacketParser.GetPacketSignature is expected from 1 to Infinity times
func (m *mPacketParserMockGetPacketSignature) Expect() *mPacketParserMockGetPacketSignature {
	m.mock.GetPacketSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPacketSignatureExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetPacketSignature
func (m *mPacketParserMockGetPacketSignature) Return(r cryptkit.SignedDigest) *PacketParserMock {
	m.mock.GetPacketSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPacketSignatureExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetPacketSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetPacketSignature is expected once
func (m *mPacketParserMockGetPacketSignature) ExpectOnce() *PacketParserMockGetPacketSignatureExpectation {
	m.mock.GetPacketSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetPacketSignatureExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetPacketSignatureExpectation) Return(r cryptkit.SignedDigest) {
	e.result = &PacketParserMockGetPacketSignatureResult{r}
}

//Set uses given function f as a mock of PacketParser.GetPacketSignature method
func (m *mPacketParserMockGetPacketSignature) Set(f func() (r cryptkit.SignedDigest)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPacketSignatureFunc = f
	return m.mock
}

//GetPacketSignature implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetPacketSignature() (r cryptkit.SignedDigest) {
	counter := atomic.AddUint64(&m.GetPacketSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.GetPacketSignatureCounter, 1)

	if len(m.GetPacketSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPacketSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetPacketSignature.")
			return
		}

		result := m.GetPacketSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPacketSignature")
			return
		}

		r = result.r

		return
	}

	if m.GetPacketSignatureMock.mainExpectation != nil {

		result := m.GetPacketSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPacketSignature")
		}

		r = result.r

		return
	}

	if m.GetPacketSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetPacketSignature.")
		return
	}

	return m.GetPacketSignatureFunc()
}

//GetPacketSignatureMinimockCounter returns a count of PacketParserMock.GetPacketSignatureFunc invocations
func (m *PacketParserMock) GetPacketSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPacketSignatureCounter)
}

//GetPacketSignatureMinimockPreCounter returns the value of PacketParserMock.GetPacketSignature invocations
func (m *PacketParserMock) GetPacketSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPacketSignaturePreCounter)
}

//GetPacketSignatureFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetPacketSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPacketSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPacketSignatureCounter) == uint64(len(m.GetPacketSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPacketSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPacketSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPacketSignatureFunc != nil {
		return atomic.LoadUint64(&m.GetPacketSignatureCounter) > 0
	}

	return true
}

type mPacketParserMockGetPacketType struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetPacketTypeExpectation
	expectationSeries []*PacketParserMockGetPacketTypeExpectation
}

type PacketParserMockGetPacketTypeExpectation struct {
	result *PacketParserMockGetPacketTypeResult
}

type PacketParserMockGetPacketTypeResult struct {
	r phases.PacketType
}

//Expect specifies that invocation of PacketParser.GetPacketType is expected from 1 to Infinity times
func (m *mPacketParserMockGetPacketType) Expect() *mPacketParserMockGetPacketType {
	m.mock.GetPacketTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPacketTypeExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetPacketType
func (m *mPacketParserMockGetPacketType) Return(r phases.PacketType) *PacketParserMock {
	m.mock.GetPacketTypeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPacketTypeExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetPacketTypeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetPacketType is expected once
func (m *mPacketParserMockGetPacketType) ExpectOnce() *PacketParserMockGetPacketTypeExpectation {
	m.mock.GetPacketTypeFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetPacketTypeExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetPacketTypeExpectation) Return(r phases.PacketType) {
	e.result = &PacketParserMockGetPacketTypeResult{r}
}

//Set uses given function f as a mock of PacketParser.GetPacketType method
func (m *mPacketParserMockGetPacketType) Set(f func() (r phases.PacketType)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPacketTypeFunc = f
	return m.mock
}

//GetPacketType implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetPacketType() (r phases.PacketType) {
	counter := atomic.AddUint64(&m.GetPacketTypePreCounter, 1)
	defer atomic.AddUint64(&m.GetPacketTypeCounter, 1)

	if len(m.GetPacketTypeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPacketTypeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetPacketType.")
			return
		}

		result := m.GetPacketTypeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPacketType")
			return
		}

		r = result.r

		return
	}

	if m.GetPacketTypeMock.mainExpectation != nil {

		result := m.GetPacketTypeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPacketType")
		}

		r = result.r

		return
	}

	if m.GetPacketTypeFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetPacketType.")
		return
	}

	return m.GetPacketTypeFunc()
}

//GetPacketTypeMinimockCounter returns a count of PacketParserMock.GetPacketTypeFunc invocations
func (m *PacketParserMock) GetPacketTypeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPacketTypeCounter)
}

//GetPacketTypeMinimockPreCounter returns the value of PacketParserMock.GetPacketType invocations
func (m *PacketParserMock) GetPacketTypeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPacketTypePreCounter)
}

//GetPacketTypeFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetPacketTypeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPacketTypeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPacketTypeCounter) == uint64(len(m.GetPacketTypeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPacketTypeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPacketTypeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPacketTypeFunc != nil {
		return atomic.LoadUint64(&m.GetPacketTypeCounter) > 0
	}

	return true
}

type mPacketParserMockGetPulseNumber struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetPulseNumberExpectation
	expectationSeries []*PacketParserMockGetPulseNumberExpectation
}

type PacketParserMockGetPulseNumberExpectation struct {
	result *PacketParserMockGetPulseNumberResult
}

type PacketParserMockGetPulseNumberResult struct {
	r pulse.Number
}

//Expect specifies that invocation of PacketParser.GetPulseNumber is expected from 1 to Infinity times
func (m *mPacketParserMockGetPulseNumber) Expect() *mPacketParserMockGetPulseNumber {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPulseNumberExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetPulseNumber
func (m *mPacketParserMockGetPulseNumber) Return(r pulse.Number) *PacketParserMock {
	m.mock.GetPulseNumberFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPulseNumberExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetPulseNumberResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetPulseNumber is expected once
func (m *mPacketParserMockGetPulseNumber) ExpectOnce() *PacketParserMockGetPulseNumberExpectation {
	m.mock.GetPulseNumberFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetPulseNumberExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetPulseNumberExpectation) Return(r pulse.Number) {
	e.result = &PacketParserMockGetPulseNumberResult{r}
}

//Set uses given function f as a mock of PacketParser.GetPulseNumber method
func (m *mPacketParserMockGetPulseNumber) Set(f func() (r pulse.Number)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulseNumberFunc = f
	return m.mock
}

//GetPulseNumber implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetPulseNumber() (r pulse.Number) {
	counter := atomic.AddUint64(&m.GetPulseNumberPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulseNumberCounter, 1)

	if len(m.GetPulseNumberMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulseNumberMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetPulseNumber.")
			return
		}

		result := m.GetPulseNumberMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPulseNumber")
			return
		}

		r = result.r

		return
	}

	if m.GetPulseNumberMock.mainExpectation != nil {

		result := m.GetPulseNumberMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPulseNumber")
		}

		r = result.r

		return
	}

	if m.GetPulseNumberFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetPulseNumber.")
		return
	}

	return m.GetPulseNumberFunc()
}

//GetPulseNumberMinimockCounter returns a count of PacketParserMock.GetPulseNumberFunc invocations
func (m *PacketParserMock) GetPulseNumberMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberCounter)
}

//GetPulseNumberMinimockPreCounter returns the value of PacketParserMock.GetPulseNumber invocations
func (m *PacketParserMock) GetPulseNumberMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulseNumberPreCounter)
}

//GetPulseNumberFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetPulseNumberFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulseNumberMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) == uint64(len(m.GetPulseNumberMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulseNumberMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulseNumberFunc != nil {
		return atomic.LoadUint64(&m.GetPulseNumberCounter) > 0
	}

	return true
}

type mPacketParserMockGetPulsePacket struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetPulsePacketExpectation
	expectationSeries []*PacketParserMockGetPulsePacketExpectation
}

type PacketParserMockGetPulsePacketExpectation struct {
	result *PacketParserMockGetPulsePacketResult
}

type PacketParserMockGetPulsePacketResult struct {
	r PulsePacketReader
}

//Expect specifies that invocation of PacketParser.GetPulsePacket is expected from 1 to Infinity times
func (m *mPacketParserMockGetPulsePacket) Expect() *mPacketParserMockGetPulsePacket {
	m.mock.GetPulsePacketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPulsePacketExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetPulsePacket
func (m *mPacketParserMockGetPulsePacket) Return(r PulsePacketReader) *PacketParserMock {
	m.mock.GetPulsePacketFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetPulsePacketExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetPulsePacketResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetPulsePacket is expected once
func (m *mPacketParserMockGetPulsePacket) ExpectOnce() *PacketParserMockGetPulsePacketExpectation {
	m.mock.GetPulsePacketFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetPulsePacketExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetPulsePacketExpectation) Return(r PulsePacketReader) {
	e.result = &PacketParserMockGetPulsePacketResult{r}
}

//Set uses given function f as a mock of PacketParser.GetPulsePacket method
func (m *mPacketParserMockGetPulsePacket) Set(f func() (r PulsePacketReader)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetPulsePacketFunc = f
	return m.mock
}

//GetPulsePacket implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetPulsePacket() (r PulsePacketReader) {
	counter := atomic.AddUint64(&m.GetPulsePacketPreCounter, 1)
	defer atomic.AddUint64(&m.GetPulsePacketCounter, 1)

	if len(m.GetPulsePacketMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetPulsePacketMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetPulsePacket.")
			return
		}

		result := m.GetPulsePacketMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPulsePacket")
			return
		}

		r = result.r

		return
	}

	if m.GetPulsePacketMock.mainExpectation != nil {

		result := m.GetPulsePacketMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetPulsePacket")
		}

		r = result.r

		return
	}

	if m.GetPulsePacketFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetPulsePacket.")
		return
	}

	return m.GetPulsePacketFunc()
}

//GetPulsePacketMinimockCounter returns a count of PacketParserMock.GetPulsePacketFunc invocations
func (m *PacketParserMock) GetPulsePacketMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulsePacketCounter)
}

//GetPulsePacketMinimockPreCounter returns the value of PacketParserMock.GetPulsePacket invocations
func (m *PacketParserMock) GetPulsePacketMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPulsePacketPreCounter)
}

//GetPulsePacketFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetPulsePacketFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetPulsePacketMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetPulsePacketCounter) == uint64(len(m.GetPulsePacketMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetPulsePacketMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetPulsePacketCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetPulsePacketFunc != nil {
		return atomic.LoadUint64(&m.GetPulsePacketCounter) > 0
	}

	return true
}

type mPacketParserMockGetReceiverID struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetReceiverIDExpectation
	expectationSeries []*PacketParserMockGetReceiverIDExpectation
}

type PacketParserMockGetReceiverIDExpectation struct {
	result *PacketParserMockGetReceiverIDResult
}

type PacketParserMockGetReceiverIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of PacketParser.GetReceiverID is expected from 1 to Infinity times
func (m *mPacketParserMockGetReceiverID) Expect() *mPacketParserMockGetReceiverID {
	m.mock.GetReceiverIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetReceiverIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetReceiverID
func (m *mPacketParserMockGetReceiverID) Return(r insolar.ShortNodeID) *PacketParserMock {
	m.mock.GetReceiverIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetReceiverIDExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetReceiverIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetReceiverID is expected once
func (m *mPacketParserMockGetReceiverID) ExpectOnce() *PacketParserMockGetReceiverIDExpectation {
	m.mock.GetReceiverIDFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetReceiverIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetReceiverIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &PacketParserMockGetReceiverIDResult{r}
}

//Set uses given function f as a mock of PacketParser.GetReceiverID method
func (m *mPacketParserMockGetReceiverID) Set(f func() (r insolar.ShortNodeID)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetReceiverIDFunc = f
	return m.mock
}

//GetReceiverID implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetReceiverID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetReceiverIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetReceiverIDCounter, 1)

	if len(m.GetReceiverIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetReceiverIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetReceiverID.")
			return
		}

		result := m.GetReceiverIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetReceiverID")
			return
		}

		r = result.r

		return
	}

	if m.GetReceiverIDMock.mainExpectation != nil {

		result := m.GetReceiverIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetReceiverID")
		}

		r = result.r

		return
	}

	if m.GetReceiverIDFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetReceiverID.")
		return
	}

	return m.GetReceiverIDFunc()
}

//GetReceiverIDMinimockCounter returns a count of PacketParserMock.GetReceiverIDFunc invocations
func (m *PacketParserMock) GetReceiverIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetReceiverIDCounter)
}

//GetReceiverIDMinimockPreCounter returns the value of PacketParserMock.GetReceiverID invocations
func (m *PacketParserMock) GetReceiverIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetReceiverIDPreCounter)
}

//GetReceiverIDFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetReceiverIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetReceiverIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetReceiverIDCounter) == uint64(len(m.GetReceiverIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetReceiverIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetReceiverIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetReceiverIDFunc != nil {
		return atomic.LoadUint64(&m.GetReceiverIDCounter) > 0
	}

	return true
}

type mPacketParserMockGetSourceID struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetSourceIDExpectation
	expectationSeries []*PacketParserMockGetSourceIDExpectation
}

type PacketParserMockGetSourceIDExpectation struct {
	result *PacketParserMockGetSourceIDResult
}

type PacketParserMockGetSourceIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of PacketParser.GetSourceID is expected from 1 to Infinity times
func (m *mPacketParserMockGetSourceID) Expect() *mPacketParserMockGetSourceID {
	m.mock.GetSourceIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetSourceIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetSourceID
func (m *mPacketParserMockGetSourceID) Return(r insolar.ShortNodeID) *PacketParserMock {
	m.mock.GetSourceIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetSourceIDExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetSourceIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetSourceID is expected once
func (m *mPacketParserMockGetSourceID) ExpectOnce() *PacketParserMockGetSourceIDExpectation {
	m.mock.GetSourceIDFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetSourceIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetSourceIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &PacketParserMockGetSourceIDResult{r}
}

//Set uses given function f as a mock of PacketParser.GetSourceID method
func (m *mPacketParserMockGetSourceID) Set(f func() (r insolar.ShortNodeID)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSourceIDFunc = f
	return m.mock
}

//GetSourceID implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetSourceID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetSourceIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetSourceIDCounter, 1)

	if len(m.GetSourceIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSourceIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetSourceID.")
			return
		}

		result := m.GetSourceIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetSourceID")
			return
		}

		r = result.r

		return
	}

	if m.GetSourceIDMock.mainExpectation != nil {

		result := m.GetSourceIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetSourceID")
		}

		r = result.r

		return
	}

	if m.GetSourceIDFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetSourceID.")
		return
	}

	return m.GetSourceIDFunc()
}

//GetSourceIDMinimockCounter returns a count of PacketParserMock.GetSourceIDFunc invocations
func (m *PacketParserMock) GetSourceIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSourceIDCounter)
}

//GetSourceIDMinimockPreCounter returns the value of PacketParserMock.GetSourceID invocations
func (m *PacketParserMock) GetSourceIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSourceIDPreCounter)
}

//GetSourceIDFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetSourceIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSourceIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSourceIDCounter) == uint64(len(m.GetSourceIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSourceIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSourceIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSourceIDFunc != nil {
		return atomic.LoadUint64(&m.GetSourceIDCounter) > 0
	}

	return true
}

type mPacketParserMockGetTargetID struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockGetTargetIDExpectation
	expectationSeries []*PacketParserMockGetTargetIDExpectation
}

type PacketParserMockGetTargetIDExpectation struct {
	result *PacketParserMockGetTargetIDResult
}

type PacketParserMockGetTargetIDResult struct {
	r insolar.ShortNodeID
}

//Expect specifies that invocation of PacketParser.GetTargetID is expected from 1 to Infinity times
func (m *mPacketParserMockGetTargetID) Expect() *mPacketParserMockGetTargetID {
	m.mock.GetTargetIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetTargetIDExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.GetTargetID
func (m *mPacketParserMockGetTargetID) Return(r insolar.ShortNodeID) *PacketParserMock {
	m.mock.GetTargetIDFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockGetTargetIDExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockGetTargetIDResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.GetTargetID is expected once
func (m *mPacketParserMockGetTargetID) ExpectOnce() *PacketParserMockGetTargetIDExpectation {
	m.mock.GetTargetIDFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockGetTargetIDExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockGetTargetIDExpectation) Return(r insolar.ShortNodeID) {
	e.result = &PacketParserMockGetTargetIDResult{r}
}

//Set uses given function f as a mock of PacketParser.GetTargetID method
func (m *mPacketParserMockGetTargetID) Set(f func() (r insolar.ShortNodeID)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetTargetIDFunc = f
	return m.mock
}

//GetTargetID implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) GetTargetID() (r insolar.ShortNodeID) {
	counter := atomic.AddUint64(&m.GetTargetIDPreCounter, 1)
	defer atomic.AddUint64(&m.GetTargetIDCounter, 1)

	if len(m.GetTargetIDMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetTargetIDMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.GetTargetID.")
			return
		}

		result := m.GetTargetIDMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetTargetID")
			return
		}

		r = result.r

		return
	}

	if m.GetTargetIDMock.mainExpectation != nil {

		result := m.GetTargetIDMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.GetTargetID")
		}

		r = result.r

		return
	}

	if m.GetTargetIDFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.GetTargetID.")
		return
	}

	return m.GetTargetIDFunc()
}

//GetTargetIDMinimockCounter returns a count of PacketParserMock.GetTargetIDFunc invocations
func (m *PacketParserMock) GetTargetIDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetTargetIDCounter)
}

//GetTargetIDMinimockPreCounter returns the value of PacketParserMock.GetTargetID invocations
func (m *PacketParserMock) GetTargetIDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetTargetIDPreCounter)
}

//GetTargetIDFinished returns true if mock invocations count is ok
func (m *PacketParserMock) GetTargetIDFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetTargetIDMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetTargetIDCounter) == uint64(len(m.GetTargetIDMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetTargetIDMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetTargetIDCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetTargetIDFunc != nil {
		return atomic.LoadUint64(&m.GetTargetIDCounter) > 0
	}

	return true
}

type mPacketParserMockIsRelayForbidden struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockIsRelayForbiddenExpectation
	expectationSeries []*PacketParserMockIsRelayForbiddenExpectation
}

type PacketParserMockIsRelayForbiddenExpectation struct {
	result *PacketParserMockIsRelayForbiddenResult
}

type PacketParserMockIsRelayForbiddenResult struct {
	r bool
}

//Expect specifies that invocation of PacketParser.IsRelayForbidden is expected from 1 to Infinity times
func (m *mPacketParserMockIsRelayForbidden) Expect() *mPacketParserMockIsRelayForbidden {
	m.mock.IsRelayForbiddenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockIsRelayForbiddenExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.IsRelayForbidden
func (m *mPacketParserMockIsRelayForbidden) Return(r bool) *PacketParserMock {
	m.mock.IsRelayForbiddenFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockIsRelayForbiddenExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockIsRelayForbiddenResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.IsRelayForbidden is expected once
func (m *mPacketParserMockIsRelayForbidden) ExpectOnce() *PacketParserMockIsRelayForbiddenExpectation {
	m.mock.IsRelayForbiddenFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockIsRelayForbiddenExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockIsRelayForbiddenExpectation) Return(r bool) {
	e.result = &PacketParserMockIsRelayForbiddenResult{r}
}

//Set uses given function f as a mock of PacketParser.IsRelayForbidden method
func (m *mPacketParserMockIsRelayForbidden) Set(f func() (r bool)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsRelayForbiddenFunc = f
	return m.mock
}

//IsRelayForbidden implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) IsRelayForbidden() (r bool) {
	counter := atomic.AddUint64(&m.IsRelayForbiddenPreCounter, 1)
	defer atomic.AddUint64(&m.IsRelayForbiddenCounter, 1)

	if len(m.IsRelayForbiddenMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsRelayForbiddenMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.IsRelayForbidden.")
			return
		}

		result := m.IsRelayForbiddenMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.IsRelayForbidden")
			return
		}

		r = result.r

		return
	}

	if m.IsRelayForbiddenMock.mainExpectation != nil {

		result := m.IsRelayForbiddenMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.IsRelayForbidden")
		}

		r = result.r

		return
	}

	if m.IsRelayForbiddenFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.IsRelayForbidden.")
		return
	}

	return m.IsRelayForbiddenFunc()
}

//IsRelayForbiddenMinimockCounter returns a count of PacketParserMock.IsRelayForbiddenFunc invocations
func (m *PacketParserMock) IsRelayForbiddenMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsRelayForbiddenCounter)
}

//IsRelayForbiddenMinimockPreCounter returns the value of PacketParserMock.IsRelayForbidden invocations
func (m *PacketParserMock) IsRelayForbiddenMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsRelayForbiddenPreCounter)
}

//IsRelayForbiddenFinished returns true if mock invocations count is ok
func (m *PacketParserMock) IsRelayForbiddenFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsRelayForbiddenMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsRelayForbiddenCounter) == uint64(len(m.IsRelayForbiddenMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsRelayForbiddenMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsRelayForbiddenCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsRelayForbiddenFunc != nil {
		return atomic.LoadUint64(&m.IsRelayForbiddenCounter) > 0
	}

	return true
}

type mPacketParserMockParsePacketBody struct {
	mock              *PacketParserMock
	mainExpectation   *PacketParserMockParsePacketBodyExpectation
	expectationSeries []*PacketParserMockParsePacketBodyExpectation
}

type PacketParserMockParsePacketBodyExpectation struct {
	result *PacketParserMockParsePacketBodyResult
}

type PacketParserMockParsePacketBodyResult struct {
	r  PacketParser
	r1 error
}

//Expect specifies that invocation of PacketParser.ParsePacketBody is expected from 1 to Infinity times
func (m *mPacketParserMockParsePacketBody) Expect() *mPacketParserMockParsePacketBody {
	m.mock.ParsePacketBodyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockParsePacketBodyExpectation{}
	}

	return m
}

//Return specifies results of invocation of PacketParser.ParsePacketBody
func (m *mPacketParserMockParsePacketBody) Return(r PacketParser, r1 error) *PacketParserMock {
	m.mock.ParsePacketBodyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PacketParserMockParsePacketBodyExpectation{}
	}
	m.mainExpectation.result = &PacketParserMockParsePacketBodyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of PacketParser.ParsePacketBody is expected once
func (m *mPacketParserMockParsePacketBody) ExpectOnce() *PacketParserMockParsePacketBodyExpectation {
	m.mock.ParsePacketBodyFunc = nil
	m.mainExpectation = nil

	expectation := &PacketParserMockParsePacketBodyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PacketParserMockParsePacketBodyExpectation) Return(r PacketParser, r1 error) {
	e.result = &PacketParserMockParsePacketBodyResult{r, r1}
}

//Set uses given function f as a mock of PacketParser.ParsePacketBody method
func (m *mPacketParserMockParsePacketBody) Set(f func() (r PacketParser, r1 error)) *PacketParserMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ParsePacketBodyFunc = f
	return m.mock
}

//ParsePacketBody implements github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser interface
func (m *PacketParserMock) ParsePacketBody() (r PacketParser, r1 error) {
	counter := atomic.AddUint64(&m.ParsePacketBodyPreCounter, 1)
	defer atomic.AddUint64(&m.ParsePacketBodyCounter, 1)

	if len(m.ParsePacketBodyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ParsePacketBodyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PacketParserMock.ParsePacketBody.")
			return
		}

		result := m.ParsePacketBodyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.ParsePacketBody")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ParsePacketBodyMock.mainExpectation != nil {

		result := m.ParsePacketBodyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PacketParserMock.ParsePacketBody")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ParsePacketBodyFunc == nil {
		m.t.Fatalf("Unexpected call to PacketParserMock.ParsePacketBody.")
		return
	}

	return m.ParsePacketBodyFunc()
}

//ParsePacketBodyMinimockCounter returns a count of PacketParserMock.ParsePacketBodyFunc invocations
func (m *PacketParserMock) ParsePacketBodyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ParsePacketBodyCounter)
}

//ParsePacketBodyMinimockPreCounter returns the value of PacketParserMock.ParsePacketBody invocations
func (m *PacketParserMock) ParsePacketBodyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ParsePacketBodyPreCounter)
}

//ParsePacketBodyFinished returns true if mock invocations count is ok
func (m *PacketParserMock) ParsePacketBodyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ParsePacketBodyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ParsePacketBodyCounter) == uint64(len(m.ParsePacketBodyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ParsePacketBodyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ParsePacketBodyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ParsePacketBodyFunc != nil {
		return atomic.LoadUint64(&m.ParsePacketBodyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PacketParserMock) ValidateCallCounters() {

	if !m.GetMemberPacketFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetMemberPacket")
	}

	if !m.GetPacketSignatureFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPacketSignature")
	}

	if !m.GetPacketTypeFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPacketType")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPulseNumber")
	}

	if !m.GetPulsePacketFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPulsePacket")
	}

	if !m.GetReceiverIDFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetReceiverID")
	}

	if !m.GetSourceIDFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetSourceID")
	}

	if !m.GetTargetIDFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetTargetID")
	}

	if !m.IsRelayForbiddenFinished() {
		m.t.Fatal("Expected call to PacketParserMock.IsRelayForbidden")
	}

	if !m.ParsePacketBodyFinished() {
		m.t.Fatal("Expected call to PacketParserMock.ParsePacketBody")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PacketParserMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PacketParserMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PacketParserMock) MinimockFinish() {

	if !m.GetMemberPacketFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetMemberPacket")
	}

	if !m.GetPacketSignatureFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPacketSignature")
	}

	if !m.GetPacketTypeFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPacketType")
	}

	if !m.GetPulseNumberFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPulseNumber")
	}

	if !m.GetPulsePacketFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetPulsePacket")
	}

	if !m.GetReceiverIDFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetReceiverID")
	}

	if !m.GetSourceIDFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetSourceID")
	}

	if !m.GetTargetIDFinished() {
		m.t.Fatal("Expected call to PacketParserMock.GetTargetID")
	}

	if !m.IsRelayForbiddenFinished() {
		m.t.Fatal("Expected call to PacketParserMock.IsRelayForbidden")
	}

	if !m.ParsePacketBodyFinished() {
		m.t.Fatal("Expected call to PacketParserMock.ParsePacketBody")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PacketParserMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PacketParserMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetMemberPacketFinished()
		ok = ok && m.GetPacketSignatureFinished()
		ok = ok && m.GetPacketTypeFinished()
		ok = ok && m.GetPulseNumberFinished()
		ok = ok && m.GetPulsePacketFinished()
		ok = ok && m.GetReceiverIDFinished()
		ok = ok && m.GetSourceIDFinished()
		ok = ok && m.GetTargetIDFinished()
		ok = ok && m.IsRelayForbiddenFinished()
		ok = ok && m.ParsePacketBodyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetMemberPacketFinished() {
				m.t.Error("Expected call to PacketParserMock.GetMemberPacket")
			}

			if !m.GetPacketSignatureFinished() {
				m.t.Error("Expected call to PacketParserMock.GetPacketSignature")
			}

			if !m.GetPacketTypeFinished() {
				m.t.Error("Expected call to PacketParserMock.GetPacketType")
			}

			if !m.GetPulseNumberFinished() {
				m.t.Error("Expected call to PacketParserMock.GetPulseNumber")
			}

			if !m.GetPulsePacketFinished() {
				m.t.Error("Expected call to PacketParserMock.GetPulsePacket")
			}

			if !m.GetReceiverIDFinished() {
				m.t.Error("Expected call to PacketParserMock.GetReceiverID")
			}

			if !m.GetSourceIDFinished() {
				m.t.Error("Expected call to PacketParserMock.GetSourceID")
			}

			if !m.GetTargetIDFinished() {
				m.t.Error("Expected call to PacketParserMock.GetTargetID")
			}

			if !m.IsRelayForbiddenFinished() {
				m.t.Error("Expected call to PacketParserMock.IsRelayForbidden")
			}

			if !m.ParsePacketBodyFinished() {
				m.t.Error("Expected call to PacketParserMock.ParsePacketBody")
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
func (m *PacketParserMock) AllMocksCalled() bool {

	if !m.GetMemberPacketFinished() {
		return false
	}

	if !m.GetPacketSignatureFinished() {
		return false
	}

	if !m.GetPacketTypeFinished() {
		return false
	}

	if !m.GetPulseNumberFinished() {
		return false
	}

	if !m.GetPulsePacketFinished() {
		return false
	}

	if !m.GetReceiverIDFinished() {
		return false
	}

	if !m.GetSourceIDFinished() {
		return false
	}

	if !m.GetTargetIDFinished() {
		return false
	}

	if !m.IsRelayForbiddenFinished() {
		return false
	}

	if !m.ParsePacketBodyFinished() {
		return false
	}

	return true
}
