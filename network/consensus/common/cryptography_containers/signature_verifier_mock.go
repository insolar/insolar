package cryptography_containers

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "SignatureVerifier" can be found in github.com/insolar/insolar/network/consensus/common/cryptography_containers
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//SignatureVerifierMock implements github.com/insolar/insolar/network/consensus/common/cryptography_containers.SignatureVerifier
type SignatureVerifierMock struct {
	t minimock.Tester

	IsDigestMethodSupportedFunc       func(p DigestMethod) (r bool)
	IsDigestMethodSupportedCounter    uint64
	IsDigestMethodSupportedPreCounter uint64
	IsDigestMethodSupportedMock       mSignatureVerifierMockIsDigestMethodSupported

	IsSignMethodSupportedFunc       func(p SignMethod) (r bool)
	IsSignMethodSupportedCounter    uint64
	IsSignMethodSupportedPreCounter uint64
	IsSignMethodSupportedMock       mSignatureVerifierMockIsSignMethodSupported

	IsSignOfSignatureMethodSupportedFunc       func(p SignatureMethod) (r bool)
	IsSignOfSignatureMethodSupportedCounter    uint64
	IsSignOfSignatureMethodSupportedPreCounter uint64
	IsSignOfSignatureMethodSupportedMock       mSignatureVerifierMockIsSignOfSignatureMethodSupported

	IsValidDataSignatureFunc       func(p io.Reader, p1 SignatureHolder) (r bool)
	IsValidDataSignatureCounter    uint64
	IsValidDataSignaturePreCounter uint64
	IsValidDataSignatureMock       mSignatureVerifierMockIsValidDataSignature

	IsValidDigestSignatureFunc       func(p DigestHolder, p1 SignatureHolder) (r bool)
	IsValidDigestSignatureCounter    uint64
	IsValidDigestSignaturePreCounter uint64
	IsValidDigestSignatureMock       mSignatureVerifierMockIsValidDigestSignature
}

//NewSignatureVerifierMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptography_containers.SignatureVerifier
func NewSignatureVerifierMock(t minimock.Tester) *SignatureVerifierMock {
	m := &SignatureVerifierMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IsDigestMethodSupportedMock = mSignatureVerifierMockIsDigestMethodSupported{mock: m}
	m.IsSignMethodSupportedMock = mSignatureVerifierMockIsSignMethodSupported{mock: m}
	m.IsSignOfSignatureMethodSupportedMock = mSignatureVerifierMockIsSignOfSignatureMethodSupported{mock: m}
	m.IsValidDataSignatureMock = mSignatureVerifierMockIsValidDataSignature{mock: m}
	m.IsValidDigestSignatureMock = mSignatureVerifierMockIsValidDigestSignature{mock: m}

	return m
}

type mSignatureVerifierMockIsDigestMethodSupported struct {
	mock              *SignatureVerifierMock
	mainExpectation   *SignatureVerifierMockIsDigestMethodSupportedExpectation
	expectationSeries []*SignatureVerifierMockIsDigestMethodSupportedExpectation
}

type SignatureVerifierMockIsDigestMethodSupportedExpectation struct {
	input  *SignatureVerifierMockIsDigestMethodSupportedInput
	result *SignatureVerifierMockIsDigestMethodSupportedResult
}

type SignatureVerifierMockIsDigestMethodSupportedInput struct {
	p DigestMethod
}

type SignatureVerifierMockIsDigestMethodSupportedResult struct {
	r bool
}

//Expect specifies that invocation of SignatureVerifier.IsDigestMethodSupported is expected from 1 to Infinity times
func (m *mSignatureVerifierMockIsDigestMethodSupported) Expect(p DigestMethod) *mSignatureVerifierMockIsDigestMethodSupported {
	m.mock.IsDigestMethodSupportedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsDigestMethodSupportedExpectation{}
	}
	m.mainExpectation.input = &SignatureVerifierMockIsDigestMethodSupportedInput{p}
	return m
}

//Return specifies results of invocation of SignatureVerifier.IsDigestMethodSupported
func (m *mSignatureVerifierMockIsDigestMethodSupported) Return(r bool) *SignatureVerifierMock {
	m.mock.IsDigestMethodSupportedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsDigestMethodSupportedExpectation{}
	}
	m.mainExpectation.result = &SignatureVerifierMockIsDigestMethodSupportedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureVerifier.IsDigestMethodSupported is expected once
func (m *mSignatureVerifierMockIsDigestMethodSupported) ExpectOnce(p DigestMethod) *SignatureVerifierMockIsDigestMethodSupportedExpectation {
	m.mock.IsDigestMethodSupportedFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureVerifierMockIsDigestMethodSupportedExpectation{}
	expectation.input = &SignatureVerifierMockIsDigestMethodSupportedInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureVerifierMockIsDigestMethodSupportedExpectation) Return(r bool) {
	e.result = &SignatureVerifierMockIsDigestMethodSupportedResult{r}
}

//Set uses given function f as a mock of SignatureVerifier.IsDigestMethodSupported method
func (m *mSignatureVerifierMockIsDigestMethodSupported) Set(f func(p DigestMethod) (r bool)) *SignatureVerifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsDigestMethodSupportedFunc = f
	return m.mock
}

//IsDigestMethodSupported implements github.com/insolar/insolar/network/consensus/common/cryptography_containers.SignatureVerifier interface
func (m *SignatureVerifierMock) IsDigestMethodSupported(p DigestMethod) (r bool) {
	counter := atomic.AddUint64(&m.IsDigestMethodSupportedPreCounter, 1)
	defer atomic.AddUint64(&m.IsDigestMethodSupportedCounter, 1)

	if len(m.IsDigestMethodSupportedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsDigestMethodSupportedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsDigestMethodSupported. %v", p)
			return
		}

		input := m.IsDigestMethodSupportedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureVerifierMockIsDigestMethodSupportedInput{p}, "SignatureVerifier.IsDigestMethodSupported got unexpected parameters")

		result := m.IsDigestMethodSupportedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsDigestMethodSupported")
			return
		}

		r = result.r

		return
	}

	if m.IsDigestMethodSupportedMock.mainExpectation != nil {

		input := m.IsDigestMethodSupportedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureVerifierMockIsDigestMethodSupportedInput{p}, "SignatureVerifier.IsDigestMethodSupported got unexpected parameters")
		}

		result := m.IsDigestMethodSupportedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsDigestMethodSupported")
		}

		r = result.r

		return
	}

	if m.IsDigestMethodSupportedFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsDigestMethodSupported. %v", p)
		return
	}

	return m.IsDigestMethodSupportedFunc(p)
}

//IsDigestMethodSupportedMinimockCounter returns a count of SignatureVerifierMock.IsDigestMethodSupportedFunc invocations
func (m *SignatureVerifierMock) IsDigestMethodSupportedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsDigestMethodSupportedCounter)
}

//IsDigestMethodSupportedMinimockPreCounter returns the value of SignatureVerifierMock.IsDigestMethodSupported invocations
func (m *SignatureVerifierMock) IsDigestMethodSupportedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsDigestMethodSupportedPreCounter)
}

//IsDigestMethodSupportedFinished returns true if mock invocations count is ok
func (m *SignatureVerifierMock) IsDigestMethodSupportedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsDigestMethodSupportedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsDigestMethodSupportedCounter) == uint64(len(m.IsDigestMethodSupportedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsDigestMethodSupportedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsDigestMethodSupportedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsDigestMethodSupportedFunc != nil {
		return atomic.LoadUint64(&m.IsDigestMethodSupportedCounter) > 0
	}

	return true
}

type mSignatureVerifierMockIsSignMethodSupported struct {
	mock              *SignatureVerifierMock
	mainExpectation   *SignatureVerifierMockIsSignMethodSupportedExpectation
	expectationSeries []*SignatureVerifierMockIsSignMethodSupportedExpectation
}

type SignatureVerifierMockIsSignMethodSupportedExpectation struct {
	input  *SignatureVerifierMockIsSignMethodSupportedInput
	result *SignatureVerifierMockIsSignMethodSupportedResult
}

type SignatureVerifierMockIsSignMethodSupportedInput struct {
	p SignMethod
}

type SignatureVerifierMockIsSignMethodSupportedResult struct {
	r bool
}

//Expect specifies that invocation of SignatureVerifier.IsSignMethodSupported is expected from 1 to Infinity times
func (m *mSignatureVerifierMockIsSignMethodSupported) Expect(p SignMethod) *mSignatureVerifierMockIsSignMethodSupported {
	m.mock.IsSignMethodSupportedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsSignMethodSupportedExpectation{}
	}
	m.mainExpectation.input = &SignatureVerifierMockIsSignMethodSupportedInput{p}
	return m
}

//Return specifies results of invocation of SignatureVerifier.IsSignMethodSupported
func (m *mSignatureVerifierMockIsSignMethodSupported) Return(r bool) *SignatureVerifierMock {
	m.mock.IsSignMethodSupportedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsSignMethodSupportedExpectation{}
	}
	m.mainExpectation.result = &SignatureVerifierMockIsSignMethodSupportedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureVerifier.IsSignMethodSupported is expected once
func (m *mSignatureVerifierMockIsSignMethodSupported) ExpectOnce(p SignMethod) *SignatureVerifierMockIsSignMethodSupportedExpectation {
	m.mock.IsSignMethodSupportedFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureVerifierMockIsSignMethodSupportedExpectation{}
	expectation.input = &SignatureVerifierMockIsSignMethodSupportedInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureVerifierMockIsSignMethodSupportedExpectation) Return(r bool) {
	e.result = &SignatureVerifierMockIsSignMethodSupportedResult{r}
}

//Set uses given function f as a mock of SignatureVerifier.IsSignMethodSupported method
func (m *mSignatureVerifierMockIsSignMethodSupported) Set(f func(p SignMethod) (r bool)) *SignatureVerifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsSignMethodSupportedFunc = f
	return m.mock
}

//IsSignMethodSupported implements github.com/insolar/insolar/network/consensus/common/cryptography_containers.SignatureVerifier interface
func (m *SignatureVerifierMock) IsSignMethodSupported(p SignMethod) (r bool) {
	counter := atomic.AddUint64(&m.IsSignMethodSupportedPreCounter, 1)
	defer atomic.AddUint64(&m.IsSignMethodSupportedCounter, 1)

	if len(m.IsSignMethodSupportedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsSignMethodSupportedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsSignMethodSupported. %v", p)
			return
		}

		input := m.IsSignMethodSupportedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureVerifierMockIsSignMethodSupportedInput{p}, "SignatureVerifier.IsSignMethodSupported got unexpected parameters")

		result := m.IsSignMethodSupportedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsSignMethodSupported")
			return
		}

		r = result.r

		return
	}

	if m.IsSignMethodSupportedMock.mainExpectation != nil {

		input := m.IsSignMethodSupportedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureVerifierMockIsSignMethodSupportedInput{p}, "SignatureVerifier.IsSignMethodSupported got unexpected parameters")
		}

		result := m.IsSignMethodSupportedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsSignMethodSupported")
		}

		r = result.r

		return
	}

	if m.IsSignMethodSupportedFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsSignMethodSupported. %v", p)
		return
	}

	return m.IsSignMethodSupportedFunc(p)
}

//IsSignMethodSupportedMinimockCounter returns a count of SignatureVerifierMock.IsSignMethodSupportedFunc invocations
func (m *SignatureVerifierMock) IsSignMethodSupportedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsSignMethodSupportedCounter)
}

//IsSignMethodSupportedMinimockPreCounter returns the value of SignatureVerifierMock.IsSignMethodSupported invocations
func (m *SignatureVerifierMock) IsSignMethodSupportedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsSignMethodSupportedPreCounter)
}

//IsSignMethodSupportedFinished returns true if mock invocations count is ok
func (m *SignatureVerifierMock) IsSignMethodSupportedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsSignMethodSupportedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsSignMethodSupportedCounter) == uint64(len(m.IsSignMethodSupportedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsSignMethodSupportedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsSignMethodSupportedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsSignMethodSupportedFunc != nil {
		return atomic.LoadUint64(&m.IsSignMethodSupportedCounter) > 0
	}

	return true
}

type mSignatureVerifierMockIsSignOfSignatureMethodSupported struct {
	mock              *SignatureVerifierMock
	mainExpectation   *SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation
	expectationSeries []*SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation
}

type SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation struct {
	input  *SignatureVerifierMockIsSignOfSignatureMethodSupportedInput
	result *SignatureVerifierMockIsSignOfSignatureMethodSupportedResult
}

type SignatureVerifierMockIsSignOfSignatureMethodSupportedInput struct {
	p SignatureMethod
}

type SignatureVerifierMockIsSignOfSignatureMethodSupportedResult struct {
	r bool
}

//Expect specifies that invocation of SignatureVerifier.IsSignOfSignatureMethodSupported is expected from 1 to Infinity times
func (m *mSignatureVerifierMockIsSignOfSignatureMethodSupported) Expect(p SignatureMethod) *mSignatureVerifierMockIsSignOfSignatureMethodSupported {
	m.mock.IsSignOfSignatureMethodSupportedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation{}
	}
	m.mainExpectation.input = &SignatureVerifierMockIsSignOfSignatureMethodSupportedInput{p}
	return m
}

//Return specifies results of invocation of SignatureVerifier.IsSignOfSignatureMethodSupported
func (m *mSignatureVerifierMockIsSignOfSignatureMethodSupported) Return(r bool) *SignatureVerifierMock {
	m.mock.IsSignOfSignatureMethodSupportedFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation{}
	}
	m.mainExpectation.result = &SignatureVerifierMockIsSignOfSignatureMethodSupportedResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureVerifier.IsSignOfSignatureMethodSupported is expected once
func (m *mSignatureVerifierMockIsSignOfSignatureMethodSupported) ExpectOnce(p SignatureMethod) *SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation {
	m.mock.IsSignOfSignatureMethodSupportedFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation{}
	expectation.input = &SignatureVerifierMockIsSignOfSignatureMethodSupportedInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureVerifierMockIsSignOfSignatureMethodSupportedExpectation) Return(r bool) {
	e.result = &SignatureVerifierMockIsSignOfSignatureMethodSupportedResult{r}
}

//Set uses given function f as a mock of SignatureVerifier.IsSignOfSignatureMethodSupported method
func (m *mSignatureVerifierMockIsSignOfSignatureMethodSupported) Set(f func(p SignatureMethod) (r bool)) *SignatureVerifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsSignOfSignatureMethodSupportedFunc = f
	return m.mock
}

//IsSignOfSignatureMethodSupported implements github.com/insolar/insolar/network/consensus/common/cryptography_containers.SignatureVerifier interface
func (m *SignatureVerifierMock) IsSignOfSignatureMethodSupported(p SignatureMethod) (r bool) {
	counter := atomic.AddUint64(&m.IsSignOfSignatureMethodSupportedPreCounter, 1)
	defer atomic.AddUint64(&m.IsSignOfSignatureMethodSupportedCounter, 1)

	if len(m.IsSignOfSignatureMethodSupportedMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsSignOfSignatureMethodSupportedMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsSignOfSignatureMethodSupported. %v", p)
			return
		}

		input := m.IsSignOfSignatureMethodSupportedMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureVerifierMockIsSignOfSignatureMethodSupportedInput{p}, "SignatureVerifier.IsSignOfSignatureMethodSupported got unexpected parameters")

		result := m.IsSignOfSignatureMethodSupportedMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsSignOfSignatureMethodSupported")
			return
		}

		r = result.r

		return
	}

	if m.IsSignOfSignatureMethodSupportedMock.mainExpectation != nil {

		input := m.IsSignOfSignatureMethodSupportedMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureVerifierMockIsSignOfSignatureMethodSupportedInput{p}, "SignatureVerifier.IsSignOfSignatureMethodSupported got unexpected parameters")
		}

		result := m.IsSignOfSignatureMethodSupportedMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsSignOfSignatureMethodSupported")
		}

		r = result.r

		return
	}

	if m.IsSignOfSignatureMethodSupportedFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsSignOfSignatureMethodSupported. %v", p)
		return
	}

	return m.IsSignOfSignatureMethodSupportedFunc(p)
}

//IsSignOfSignatureMethodSupportedMinimockCounter returns a count of SignatureVerifierMock.IsSignOfSignatureMethodSupportedFunc invocations
func (m *SignatureVerifierMock) IsSignOfSignatureMethodSupportedMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsSignOfSignatureMethodSupportedCounter)
}

//IsSignOfSignatureMethodSupportedMinimockPreCounter returns the value of SignatureVerifierMock.IsSignOfSignatureMethodSupported invocations
func (m *SignatureVerifierMock) IsSignOfSignatureMethodSupportedMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsSignOfSignatureMethodSupportedPreCounter)
}

//IsSignOfSignatureMethodSupportedFinished returns true if mock invocations count is ok
func (m *SignatureVerifierMock) IsSignOfSignatureMethodSupportedFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsSignOfSignatureMethodSupportedMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsSignOfSignatureMethodSupportedCounter) == uint64(len(m.IsSignOfSignatureMethodSupportedMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsSignOfSignatureMethodSupportedMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsSignOfSignatureMethodSupportedCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsSignOfSignatureMethodSupportedFunc != nil {
		return atomic.LoadUint64(&m.IsSignOfSignatureMethodSupportedCounter) > 0
	}

	return true
}

type mSignatureVerifierMockIsValidDataSignature struct {
	mock              *SignatureVerifierMock
	mainExpectation   *SignatureVerifierMockIsValidDataSignatureExpectation
	expectationSeries []*SignatureVerifierMockIsValidDataSignatureExpectation
}

type SignatureVerifierMockIsValidDataSignatureExpectation struct {
	input  *SignatureVerifierMockIsValidDataSignatureInput
	result *SignatureVerifierMockIsValidDataSignatureResult
}

type SignatureVerifierMockIsValidDataSignatureInput struct {
	p  io.Reader
	p1 SignatureHolder
}

type SignatureVerifierMockIsValidDataSignatureResult struct {
	r bool
}

//Expect specifies that invocation of SignatureVerifier.IsValidDataSignature is expected from 1 to Infinity times
func (m *mSignatureVerifierMockIsValidDataSignature) Expect(p io.Reader, p1 SignatureHolder) *mSignatureVerifierMockIsValidDataSignature {
	m.mock.IsValidDataSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsValidDataSignatureExpectation{}
	}
	m.mainExpectation.input = &SignatureVerifierMockIsValidDataSignatureInput{p, p1}
	return m
}

//Return specifies results of invocation of SignatureVerifier.IsValidDataSignature
func (m *mSignatureVerifierMockIsValidDataSignature) Return(r bool) *SignatureVerifierMock {
	m.mock.IsValidDataSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsValidDataSignatureExpectation{}
	}
	m.mainExpectation.result = &SignatureVerifierMockIsValidDataSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureVerifier.IsValidDataSignature is expected once
func (m *mSignatureVerifierMockIsValidDataSignature) ExpectOnce(p io.Reader, p1 SignatureHolder) *SignatureVerifierMockIsValidDataSignatureExpectation {
	m.mock.IsValidDataSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureVerifierMockIsValidDataSignatureExpectation{}
	expectation.input = &SignatureVerifierMockIsValidDataSignatureInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureVerifierMockIsValidDataSignatureExpectation) Return(r bool) {
	e.result = &SignatureVerifierMockIsValidDataSignatureResult{r}
}

//Set uses given function f as a mock of SignatureVerifier.IsValidDataSignature method
func (m *mSignatureVerifierMockIsValidDataSignature) Set(f func(p io.Reader, p1 SignatureHolder) (r bool)) *SignatureVerifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsValidDataSignatureFunc = f
	return m.mock
}

//IsValidDataSignature implements github.com/insolar/insolar/network/consensus/common/cryptography_containers.SignatureVerifier interface
func (m *SignatureVerifierMock) IsValidDataSignature(p io.Reader, p1 SignatureHolder) (r bool) {
	counter := atomic.AddUint64(&m.IsValidDataSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.IsValidDataSignatureCounter, 1)

	if len(m.IsValidDataSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsValidDataSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsValidDataSignature. %v %v", p, p1)
			return
		}

		input := m.IsValidDataSignatureMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureVerifierMockIsValidDataSignatureInput{p, p1}, "SignatureVerifier.IsValidDataSignature got unexpected parameters")

		result := m.IsValidDataSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsValidDataSignature")
			return
		}

		r = result.r

		return
	}

	if m.IsValidDataSignatureMock.mainExpectation != nil {

		input := m.IsValidDataSignatureMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureVerifierMockIsValidDataSignatureInput{p, p1}, "SignatureVerifier.IsValidDataSignature got unexpected parameters")
		}

		result := m.IsValidDataSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsValidDataSignature")
		}

		r = result.r

		return
	}

	if m.IsValidDataSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsValidDataSignature. %v %v", p, p1)
		return
	}

	return m.IsValidDataSignatureFunc(p, p1)
}

//IsValidDataSignatureMinimockCounter returns a count of SignatureVerifierMock.IsValidDataSignatureFunc invocations
func (m *SignatureVerifierMock) IsValidDataSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidDataSignatureCounter)
}

//IsValidDataSignatureMinimockPreCounter returns the value of SignatureVerifierMock.IsValidDataSignature invocations
func (m *SignatureVerifierMock) IsValidDataSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidDataSignaturePreCounter)
}

//IsValidDataSignatureFinished returns true if mock invocations count is ok
func (m *SignatureVerifierMock) IsValidDataSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsValidDataSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsValidDataSignatureCounter) == uint64(len(m.IsValidDataSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsValidDataSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsValidDataSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsValidDataSignatureFunc != nil {
		return atomic.LoadUint64(&m.IsValidDataSignatureCounter) > 0
	}

	return true
}

type mSignatureVerifierMockIsValidDigestSignature struct {
	mock              *SignatureVerifierMock
	mainExpectation   *SignatureVerifierMockIsValidDigestSignatureExpectation
	expectationSeries []*SignatureVerifierMockIsValidDigestSignatureExpectation
}

type SignatureVerifierMockIsValidDigestSignatureExpectation struct {
	input  *SignatureVerifierMockIsValidDigestSignatureInput
	result *SignatureVerifierMockIsValidDigestSignatureResult
}

type SignatureVerifierMockIsValidDigestSignatureInput struct {
	p  DigestHolder
	p1 SignatureHolder
}

type SignatureVerifierMockIsValidDigestSignatureResult struct {
	r bool
}

//Expect specifies that invocation of SignatureVerifier.IsValidDigestSignature is expected from 1 to Infinity times
func (m *mSignatureVerifierMockIsValidDigestSignature) Expect(p DigestHolder, p1 SignatureHolder) *mSignatureVerifierMockIsValidDigestSignature {
	m.mock.IsValidDigestSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsValidDigestSignatureExpectation{}
	}
	m.mainExpectation.input = &SignatureVerifierMockIsValidDigestSignatureInput{p, p1}
	return m
}

//Return specifies results of invocation of SignatureVerifier.IsValidDigestSignature
func (m *mSignatureVerifierMockIsValidDigestSignature) Return(r bool) *SignatureVerifierMock {
	m.mock.IsValidDigestSignatureFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &SignatureVerifierMockIsValidDigestSignatureExpectation{}
	}
	m.mainExpectation.result = &SignatureVerifierMockIsValidDigestSignatureResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of SignatureVerifier.IsValidDigestSignature is expected once
func (m *mSignatureVerifierMockIsValidDigestSignature) ExpectOnce(p DigestHolder, p1 SignatureHolder) *SignatureVerifierMockIsValidDigestSignatureExpectation {
	m.mock.IsValidDigestSignatureFunc = nil
	m.mainExpectation = nil

	expectation := &SignatureVerifierMockIsValidDigestSignatureExpectation{}
	expectation.input = &SignatureVerifierMockIsValidDigestSignatureInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *SignatureVerifierMockIsValidDigestSignatureExpectation) Return(r bool) {
	e.result = &SignatureVerifierMockIsValidDigestSignatureResult{r}
}

//Set uses given function f as a mock of SignatureVerifier.IsValidDigestSignature method
func (m *mSignatureVerifierMockIsValidDigestSignature) Set(f func(p DigestHolder, p1 SignatureHolder) (r bool)) *SignatureVerifierMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.IsValidDigestSignatureFunc = f
	return m.mock
}

//IsValidDigestSignature implements github.com/insolar/insolar/network/consensus/common/cryptography_containers.SignatureVerifier interface
func (m *SignatureVerifierMock) IsValidDigestSignature(p DigestHolder, p1 SignatureHolder) (r bool) {
	counter := atomic.AddUint64(&m.IsValidDigestSignaturePreCounter, 1)
	defer atomic.AddUint64(&m.IsValidDigestSignatureCounter, 1)

	if len(m.IsValidDigestSignatureMock.expectationSeries) > 0 {
		if counter > uint64(len(m.IsValidDigestSignatureMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsValidDigestSignature. %v %v", p, p1)
			return
		}

		input := m.IsValidDigestSignatureMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, SignatureVerifierMockIsValidDigestSignatureInput{p, p1}, "SignatureVerifier.IsValidDigestSignature got unexpected parameters")

		result := m.IsValidDigestSignatureMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsValidDigestSignature")
			return
		}

		r = result.r

		return
	}

	if m.IsValidDigestSignatureMock.mainExpectation != nil {

		input := m.IsValidDigestSignatureMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, SignatureVerifierMockIsValidDigestSignatureInput{p, p1}, "SignatureVerifier.IsValidDigestSignature got unexpected parameters")
		}

		result := m.IsValidDigestSignatureMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the SignatureVerifierMock.IsValidDigestSignature")
		}

		r = result.r

		return
	}

	if m.IsValidDigestSignatureFunc == nil {
		m.t.Fatalf("Unexpected call to SignatureVerifierMock.IsValidDigestSignature. %v %v", p, p1)
		return
	}

	return m.IsValidDigestSignatureFunc(p, p1)
}

//IsValidDigestSignatureMinimockCounter returns a count of SignatureVerifierMock.IsValidDigestSignatureFunc invocations
func (m *SignatureVerifierMock) IsValidDigestSignatureMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidDigestSignatureCounter)
}

//IsValidDigestSignatureMinimockPreCounter returns the value of SignatureVerifierMock.IsValidDigestSignature invocations
func (m *SignatureVerifierMock) IsValidDigestSignatureMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IsValidDigestSignaturePreCounter)
}

//IsValidDigestSignatureFinished returns true if mock invocations count is ok
func (m *SignatureVerifierMock) IsValidDigestSignatureFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.IsValidDigestSignatureMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.IsValidDigestSignatureCounter) == uint64(len(m.IsValidDigestSignatureMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.IsValidDigestSignatureMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.IsValidDigestSignatureCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.IsValidDigestSignatureFunc != nil {
		return atomic.LoadUint64(&m.IsValidDigestSignatureCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SignatureVerifierMock) ValidateCallCounters() {

	if !m.IsDigestMethodSupportedFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsDigestMethodSupported")
	}

	if !m.IsSignMethodSupportedFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsSignMethodSupported")
	}

	if !m.IsSignOfSignatureMethodSupportedFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsSignOfSignatureMethodSupported")
	}

	if !m.IsValidDataSignatureFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsValidDataSignature")
	}

	if !m.IsValidDigestSignatureFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsValidDigestSignature")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *SignatureVerifierMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *SignatureVerifierMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *SignatureVerifierMock) MinimockFinish() {

	if !m.IsDigestMethodSupportedFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsDigestMethodSupported")
	}

	if !m.IsSignMethodSupportedFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsSignMethodSupported")
	}

	if !m.IsSignOfSignatureMethodSupportedFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsSignOfSignatureMethodSupported")
	}

	if !m.IsValidDataSignatureFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsValidDataSignature")
	}

	if !m.IsValidDigestSignatureFinished() {
		m.t.Fatal("Expected call to SignatureVerifierMock.IsValidDigestSignature")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *SignatureVerifierMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *SignatureVerifierMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.IsDigestMethodSupportedFinished()
		ok = ok && m.IsSignMethodSupportedFinished()
		ok = ok && m.IsSignOfSignatureMethodSupportedFinished()
		ok = ok && m.IsValidDataSignatureFinished()
		ok = ok && m.IsValidDigestSignatureFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.IsDigestMethodSupportedFinished() {
				m.t.Error("Expected call to SignatureVerifierMock.IsDigestMethodSupported")
			}

			if !m.IsSignMethodSupportedFinished() {
				m.t.Error("Expected call to SignatureVerifierMock.IsSignMethodSupported")
			}

			if !m.IsSignOfSignatureMethodSupportedFinished() {
				m.t.Error("Expected call to SignatureVerifierMock.IsSignOfSignatureMethodSupported")
			}

			if !m.IsValidDataSignatureFinished() {
				m.t.Error("Expected call to SignatureVerifierMock.IsValidDataSignature")
			}

			if !m.IsValidDigestSignatureFinished() {
				m.t.Error("Expected call to SignatureVerifierMock.IsValidDigestSignature")
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
func (m *SignatureVerifierMock) AllMocksCalled() bool {

	if !m.IsDigestMethodSupportedFinished() {
		return false
	}

	if !m.IsSignMethodSupportedFinished() {
		return false
	}

	if !m.IsSignOfSignatureMethodSupportedFinished() {
		return false
	}

	if !m.IsValidDataSignatureFinished() {
		return false
	}

	if !m.IsValidDigestSignatureFinished() {
		return false
	}

	return true
}
