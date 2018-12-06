package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "KeyProcessor" can be found in github.com/insolar/insolar/core
*/
import (
	crypto "crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//KeyProcessorMock implements github.com/insolar/insolar/core.KeyProcessor
type KeyProcessorMock struct {
	t minimock.Tester

	ExportPrivateKeyFunc       func(p crypto.PrivateKey) (r []byte, r1 error)
	ExportPrivateKeyCounter    uint64
	ExportPrivateKeyPreCounter uint64
	ExportPrivateKeyMock       mKeyProcessorMockExportPrivateKey

	ExportPublicKeyFunc       func(p crypto.PublicKey) (r []byte, r1 error)
	ExportPublicKeyCounter    uint64
	ExportPublicKeyPreCounter uint64
	ExportPublicKeyMock       mKeyProcessorMockExportPublicKey

	ExtractPublicKeyFunc       func(p crypto.PrivateKey) (r crypto.PublicKey)
	ExtractPublicKeyCounter    uint64
	ExtractPublicKeyPreCounter uint64
	ExtractPublicKeyMock       mKeyProcessorMockExtractPublicKey

	GeneratePrivateKeyFunc       func() (r crypto.PrivateKey, r1 error)
	GeneratePrivateKeyCounter    uint64
	GeneratePrivateKeyPreCounter uint64
	GeneratePrivateKeyMock       mKeyProcessorMockGeneratePrivateKey

	ImportPrivateKeyFunc       func(p []byte) (r crypto.PrivateKey, r1 error)
	ImportPrivateKeyCounter    uint64
	ImportPrivateKeyPreCounter uint64
	ImportPrivateKeyMock       mKeyProcessorMockImportPrivateKey

	ImportPublicKeyFunc       func(p []byte) (r crypto.PublicKey, r1 error)
	ImportPublicKeyCounter    uint64
	ImportPublicKeyPreCounter uint64
	ImportPublicKeyMock       mKeyProcessorMockImportPublicKey
}

//NewKeyProcessorMock returns a mock for github.com/insolar/insolar/core.KeyProcessor
func NewKeyProcessorMock(t minimock.Tester) *KeyProcessorMock {
	m := &KeyProcessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExportPrivateKeyMock = mKeyProcessorMockExportPrivateKey{mock: m}
	m.ExportPublicKeyMock = mKeyProcessorMockExportPublicKey{mock: m}
	m.ExtractPublicKeyMock = mKeyProcessorMockExtractPublicKey{mock: m}
	m.GeneratePrivateKeyMock = mKeyProcessorMockGeneratePrivateKey{mock: m}
	m.ImportPrivateKeyMock = mKeyProcessorMockImportPrivateKey{mock: m}
	m.ImportPublicKeyMock = mKeyProcessorMockImportPublicKey{mock: m}

	return m
}

type mKeyProcessorMockExportPrivateKey struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockExportPrivateKeyExpectation
	expectationSeries []*KeyProcessorMockExportPrivateKeyExpectation
}

type KeyProcessorMockExportPrivateKeyExpectation struct {
	input  *KeyProcessorMockExportPrivateKeyInput
	result *KeyProcessorMockExportPrivateKeyResult
}

type KeyProcessorMockExportPrivateKeyInput struct {
	p crypto.PrivateKey
}

type KeyProcessorMockExportPrivateKeyResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ExportPrivateKey is expected from 1 to Infinity times
func (m *mKeyProcessorMockExportPrivateKey) Expect(p crypto.PrivateKey) *mKeyProcessorMockExportPrivateKey {
	m.mock.ExportPrivateKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPrivateKeyExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockExportPrivateKeyInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ExportPrivateKey
func (m *mKeyProcessorMockExportPrivateKey) Return(r []byte, r1 error) *KeyProcessorMock {
	m.mock.ExportPrivateKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPrivateKeyExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockExportPrivateKeyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ExportPrivateKey is expected once
func (m *mKeyProcessorMockExportPrivateKey) ExpectOnce(p crypto.PrivateKey) *KeyProcessorMockExportPrivateKeyExpectation {
	m.mock.ExportPrivateKeyFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockExportPrivateKeyExpectation{}
	expectation.input = &KeyProcessorMockExportPrivateKeyInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockExportPrivateKeyExpectation) Return(r []byte, r1 error) {
	e.result = &KeyProcessorMockExportPrivateKeyResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ExportPrivateKey method
func (m *mKeyProcessorMockExportPrivateKey) Set(f func(p crypto.PrivateKey) (r []byte, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExportPrivateKeyFunc = f
	return m.mock
}

//ExportPrivateKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExportPrivateKey(p crypto.PrivateKey) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.ExportPrivateKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ExportPrivateKeyCounter, 1)

	if len(m.ExportPrivateKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExportPrivateKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPrivateKey. %v", p)
			return
		}

		input := m.ExportPrivateKeyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockExportPrivateKeyInput{p}, "KeyProcessor.ExportPrivateKey got unexpected parameters")

		result := m.ExportPrivateKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPrivateKey")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPrivateKeyMock.mainExpectation != nil {

		input := m.ExportPrivateKeyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockExportPrivateKeyInput{p}, "KeyProcessor.ExportPrivateKey got unexpected parameters")
		}

		result := m.ExportPrivateKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPrivateKey")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPrivateKeyFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPrivateKey. %v", p)
		return
	}

	return m.ExportPrivateKeyFunc(p)
}

//ExportPrivateKeyMinimockCounter returns a count of KeyProcessorMock.ExportPrivateKeyFunc invocations
func (m *KeyProcessorMock) ExportPrivateKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPrivateKeyCounter)
}

//ExportPrivateKeyMinimockPreCounter returns the value of KeyProcessorMock.ExportPrivateKey invocations
func (m *KeyProcessorMock) ExportPrivateKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPrivateKeyPreCounter)
}

//ExportPrivateKeyFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ExportPrivateKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExportPrivateKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExportPrivateKeyCounter) == uint64(len(m.ExportPrivateKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExportPrivateKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExportPrivateKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExportPrivateKeyFunc != nil {
		return atomic.LoadUint64(&m.ExportPrivateKeyCounter) > 0
	}

	return true
}

type mKeyProcessorMockExportPublicKey struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockExportPublicKeyExpectation
	expectationSeries []*KeyProcessorMockExportPublicKeyExpectation
}

type KeyProcessorMockExportPublicKeyExpectation struct {
	input  *KeyProcessorMockExportPublicKeyInput
	result *KeyProcessorMockExportPublicKeyResult
}

type KeyProcessorMockExportPublicKeyInput struct {
	p crypto.PublicKey
}

type KeyProcessorMockExportPublicKeyResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ExportPublicKey is expected from 1 to Infinity times
func (m *mKeyProcessorMockExportPublicKey) Expect(p crypto.PublicKey) *mKeyProcessorMockExportPublicKey {
	m.mock.ExportPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPublicKeyExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockExportPublicKeyInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ExportPublicKey
func (m *mKeyProcessorMockExportPublicKey) Return(r []byte, r1 error) *KeyProcessorMock {
	m.mock.ExportPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPublicKeyExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockExportPublicKeyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ExportPublicKey is expected once
func (m *mKeyProcessorMockExportPublicKey) ExpectOnce(p crypto.PublicKey) *KeyProcessorMockExportPublicKeyExpectation {
	m.mock.ExportPublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockExportPublicKeyExpectation{}
	expectation.input = &KeyProcessorMockExportPublicKeyInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockExportPublicKeyExpectation) Return(r []byte, r1 error) {
	e.result = &KeyProcessorMockExportPublicKeyResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ExportPublicKey method
func (m *mKeyProcessorMockExportPublicKey) Set(f func(p crypto.PublicKey) (r []byte, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExportPublicKeyFunc = f
	return m.mock
}

//ExportPublicKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExportPublicKey(p crypto.PublicKey) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.ExportPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ExportPublicKeyCounter, 1)

	if len(m.ExportPublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExportPublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPublicKey. %v", p)
			return
		}

		input := m.ExportPublicKeyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockExportPublicKeyInput{p}, "KeyProcessor.ExportPublicKey got unexpected parameters")

		result := m.ExportPublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPublicKey")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPublicKeyMock.mainExpectation != nil {

		input := m.ExportPublicKeyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockExportPublicKeyInput{p}, "KeyProcessor.ExportPublicKey got unexpected parameters")
		}

		result := m.ExportPublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPublicKey")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPublicKey. %v", p)
		return
	}

	return m.ExportPublicKeyFunc(p)
}

//ExportPublicKeyMinimockCounter returns a count of KeyProcessorMock.ExportPublicKeyFunc invocations
func (m *KeyProcessorMock) ExportPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPublicKeyCounter)
}

//ExportPublicKeyMinimockPreCounter returns the value of KeyProcessorMock.ExportPublicKey invocations
func (m *KeyProcessorMock) ExportPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPublicKeyPreCounter)
}

//ExportPublicKeyFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ExportPublicKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExportPublicKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExportPublicKeyCounter) == uint64(len(m.ExportPublicKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExportPublicKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExportPublicKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExportPublicKeyFunc != nil {
		return atomic.LoadUint64(&m.ExportPublicKeyCounter) > 0
	}

	return true
}

type mKeyProcessorMockExtractPublicKey struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockExtractPublicKeyExpectation
	expectationSeries []*KeyProcessorMockExtractPublicKeyExpectation
}

type KeyProcessorMockExtractPublicKeyExpectation struct {
	input  *KeyProcessorMockExtractPublicKeyInput
	result *KeyProcessorMockExtractPublicKeyResult
}

type KeyProcessorMockExtractPublicKeyInput struct {
	p crypto.PrivateKey
}

type KeyProcessorMockExtractPublicKeyResult struct {
	r crypto.PublicKey
}

//Expect specifies that invocation of KeyProcessor.ExtractPublicKey is expected from 1 to Infinity times
func (m *mKeyProcessorMockExtractPublicKey) Expect(p crypto.PrivateKey) *mKeyProcessorMockExtractPublicKey {
	m.mock.ExtractPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExtractPublicKeyExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockExtractPublicKeyInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ExtractPublicKey
func (m *mKeyProcessorMockExtractPublicKey) Return(r crypto.PublicKey) *KeyProcessorMock {
	m.mock.ExtractPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExtractPublicKeyExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockExtractPublicKeyResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ExtractPublicKey is expected once
func (m *mKeyProcessorMockExtractPublicKey) ExpectOnce(p crypto.PrivateKey) *KeyProcessorMockExtractPublicKeyExpectation {
	m.mock.ExtractPublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockExtractPublicKeyExpectation{}
	expectation.input = &KeyProcessorMockExtractPublicKeyInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockExtractPublicKeyExpectation) Return(r crypto.PublicKey) {
	e.result = &KeyProcessorMockExtractPublicKeyResult{r}
}

//Set uses given function f as a mock of KeyProcessor.ExtractPublicKey method
func (m *mKeyProcessorMockExtractPublicKey) Set(f func(p crypto.PrivateKey) (r crypto.PublicKey)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExtractPublicKeyFunc = f
	return m.mock
}

//ExtractPublicKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExtractPublicKey(p crypto.PrivateKey) (r crypto.PublicKey) {
	counter := atomic.AddUint64(&m.ExtractPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ExtractPublicKeyCounter, 1)

	if len(m.ExtractPublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExtractPublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ExtractPublicKey. %v", p)
			return
		}

		input := m.ExtractPublicKeyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockExtractPublicKeyInput{p}, "KeyProcessor.ExtractPublicKey got unexpected parameters")

		result := m.ExtractPublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExtractPublicKey")
			return
		}

		r = result.r

		return
	}

	if m.ExtractPublicKeyMock.mainExpectation != nil {

		input := m.ExtractPublicKeyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockExtractPublicKeyInput{p}, "KeyProcessor.ExtractPublicKey got unexpected parameters")
		}

		result := m.ExtractPublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExtractPublicKey")
		}

		r = result.r

		return
	}

	if m.ExtractPublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ExtractPublicKey. %v", p)
		return
	}

	return m.ExtractPublicKeyFunc(p)
}

//ExtractPublicKeyMinimockCounter returns a count of KeyProcessorMock.ExtractPublicKeyFunc invocations
func (m *KeyProcessorMock) ExtractPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExtractPublicKeyCounter)
}

//ExtractPublicKeyMinimockPreCounter returns the value of KeyProcessorMock.ExtractPublicKey invocations
func (m *KeyProcessorMock) ExtractPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExtractPublicKeyPreCounter)
}

//ExtractPublicKeyFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ExtractPublicKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExtractPublicKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExtractPublicKeyCounter) == uint64(len(m.ExtractPublicKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExtractPublicKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExtractPublicKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExtractPublicKeyFunc != nil {
		return atomic.LoadUint64(&m.ExtractPublicKeyCounter) > 0
	}

	return true
}

type mKeyProcessorMockGeneratePrivateKey struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockGeneratePrivateKeyExpectation
	expectationSeries []*KeyProcessorMockGeneratePrivateKeyExpectation
}

type KeyProcessorMockGeneratePrivateKeyExpectation struct {
	result *KeyProcessorMockGeneratePrivateKeyResult
}

type KeyProcessorMockGeneratePrivateKeyResult struct {
	r  crypto.PrivateKey
	r1 error
}

//Expect specifies that invocation of KeyProcessor.GeneratePrivateKey is expected from 1 to Infinity times
func (m *mKeyProcessorMockGeneratePrivateKey) Expect() *mKeyProcessorMockGeneratePrivateKey {
	m.mock.GeneratePrivateKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockGeneratePrivateKeyExpectation{}
	}

	return m
}

//Return specifies results of invocation of KeyProcessor.GeneratePrivateKey
func (m *mKeyProcessorMockGeneratePrivateKey) Return(r crypto.PrivateKey, r1 error) *KeyProcessorMock {
	m.mock.GeneratePrivateKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockGeneratePrivateKeyExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockGeneratePrivateKeyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.GeneratePrivateKey is expected once
func (m *mKeyProcessorMockGeneratePrivateKey) ExpectOnce() *KeyProcessorMockGeneratePrivateKeyExpectation {
	m.mock.GeneratePrivateKeyFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockGeneratePrivateKeyExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockGeneratePrivateKeyExpectation) Return(r crypto.PrivateKey, r1 error) {
	e.result = &KeyProcessorMockGeneratePrivateKeyResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.GeneratePrivateKey method
func (m *mKeyProcessorMockGeneratePrivateKey) Set(f func() (r crypto.PrivateKey, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GeneratePrivateKeyFunc = f
	return m.mock
}

//GeneratePrivateKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) GeneratePrivateKey() (r crypto.PrivateKey, r1 error) {
	counter := atomic.AddUint64(&m.GeneratePrivateKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GeneratePrivateKeyCounter, 1)

	if len(m.GeneratePrivateKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GeneratePrivateKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.GeneratePrivateKey.")
			return
		}

		result := m.GeneratePrivateKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.GeneratePrivateKey")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GeneratePrivateKeyMock.mainExpectation != nil {

		result := m.GeneratePrivateKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.GeneratePrivateKey")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.GeneratePrivateKeyFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.GeneratePrivateKey.")
		return
	}

	return m.GeneratePrivateKeyFunc()
}

//GeneratePrivateKeyMinimockCounter returns a count of KeyProcessorMock.GeneratePrivateKeyFunc invocations
func (m *KeyProcessorMock) GeneratePrivateKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GeneratePrivateKeyCounter)
}

//GeneratePrivateKeyMinimockPreCounter returns the value of KeyProcessorMock.GeneratePrivateKey invocations
func (m *KeyProcessorMock) GeneratePrivateKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GeneratePrivateKeyPreCounter)
}

//GeneratePrivateKeyFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) GeneratePrivateKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GeneratePrivateKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GeneratePrivateKeyCounter) == uint64(len(m.GeneratePrivateKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GeneratePrivateKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GeneratePrivateKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GeneratePrivateKeyFunc != nil {
		return atomic.LoadUint64(&m.GeneratePrivateKeyCounter) > 0
	}

	return true
}

type mKeyProcessorMockImportPrivateKey struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockImportPrivateKeyExpectation
	expectationSeries []*KeyProcessorMockImportPrivateKeyExpectation
}

type KeyProcessorMockImportPrivateKeyExpectation struct {
	input  *KeyProcessorMockImportPrivateKeyInput
	result *KeyProcessorMockImportPrivateKeyResult
}

type KeyProcessorMockImportPrivateKeyInput struct {
	p []byte
}

type KeyProcessorMockImportPrivateKeyResult struct {
	r  crypto.PrivateKey
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ImportPrivateKey is expected from 1 to Infinity times
func (m *mKeyProcessorMockImportPrivateKey) Expect(p []byte) *mKeyProcessorMockImportPrivateKey {
	m.mock.ImportPrivateKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPrivateKeyExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockImportPrivateKeyInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ImportPrivateKey
func (m *mKeyProcessorMockImportPrivateKey) Return(r crypto.PrivateKey, r1 error) *KeyProcessorMock {
	m.mock.ImportPrivateKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPrivateKeyExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockImportPrivateKeyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ImportPrivateKey is expected once
func (m *mKeyProcessorMockImportPrivateKey) ExpectOnce(p []byte) *KeyProcessorMockImportPrivateKeyExpectation {
	m.mock.ImportPrivateKeyFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockImportPrivateKeyExpectation{}
	expectation.input = &KeyProcessorMockImportPrivateKeyInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockImportPrivateKeyExpectation) Return(r crypto.PrivateKey, r1 error) {
	e.result = &KeyProcessorMockImportPrivateKeyResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ImportPrivateKey method
func (m *mKeyProcessorMockImportPrivateKey) Set(f func(p []byte) (r crypto.PrivateKey, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ImportPrivateKeyFunc = f
	return m.mock
}

//ImportPrivateKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ImportPrivateKey(p []byte) (r crypto.PrivateKey, r1 error) {
	counter := atomic.AddUint64(&m.ImportPrivateKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ImportPrivateKeyCounter, 1)

	if len(m.ImportPrivateKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ImportPrivateKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPrivateKey. %v", p)
			return
		}

		input := m.ImportPrivateKeyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockImportPrivateKeyInput{p}, "KeyProcessor.ImportPrivateKey got unexpected parameters")

		result := m.ImportPrivateKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPrivateKey")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPrivateKeyMock.mainExpectation != nil {

		input := m.ImportPrivateKeyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockImportPrivateKeyInput{p}, "KeyProcessor.ImportPrivateKey got unexpected parameters")
		}

		result := m.ImportPrivateKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPrivateKey")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPrivateKeyFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPrivateKey. %v", p)
		return
	}

	return m.ImportPrivateKeyFunc(p)
}

//ImportPrivateKeyMinimockCounter returns a count of KeyProcessorMock.ImportPrivateKeyFunc invocations
func (m *KeyProcessorMock) ImportPrivateKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPrivateKeyCounter)
}

//ImportPrivateKeyMinimockPreCounter returns the value of KeyProcessorMock.ImportPrivateKey invocations
func (m *KeyProcessorMock) ImportPrivateKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPrivateKeyPreCounter)
}

//ImportPrivateKeyFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ImportPrivateKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ImportPrivateKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ImportPrivateKeyCounter) == uint64(len(m.ImportPrivateKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ImportPrivateKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ImportPrivateKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ImportPrivateKeyFunc != nil {
		return atomic.LoadUint64(&m.ImportPrivateKeyCounter) > 0
	}

	return true
}

type mKeyProcessorMockImportPublicKey struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockImportPublicKeyExpectation
	expectationSeries []*KeyProcessorMockImportPublicKeyExpectation
}

type KeyProcessorMockImportPublicKeyExpectation struct {
	input  *KeyProcessorMockImportPublicKeyInput
	result *KeyProcessorMockImportPublicKeyResult
}

type KeyProcessorMockImportPublicKeyInput struct {
	p []byte
}

type KeyProcessorMockImportPublicKeyResult struct {
	r  crypto.PublicKey
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ImportPublicKey is expected from 1 to Infinity times
func (m *mKeyProcessorMockImportPublicKey) Expect(p []byte) *mKeyProcessorMockImportPublicKey {
	m.mock.ImportPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPublicKeyExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockImportPublicKeyInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ImportPublicKey
func (m *mKeyProcessorMockImportPublicKey) Return(r crypto.PublicKey, r1 error) *KeyProcessorMock {
	m.mock.ImportPublicKeyFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPublicKeyExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockImportPublicKeyResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ImportPublicKey is expected once
func (m *mKeyProcessorMockImportPublicKey) ExpectOnce(p []byte) *KeyProcessorMockImportPublicKeyExpectation {
	m.mock.ImportPublicKeyFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockImportPublicKeyExpectation{}
	expectation.input = &KeyProcessorMockImportPublicKeyInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockImportPublicKeyExpectation) Return(r crypto.PublicKey, r1 error) {
	e.result = &KeyProcessorMockImportPublicKeyResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ImportPublicKey method
func (m *mKeyProcessorMockImportPublicKey) Set(f func(p []byte) (r crypto.PublicKey, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ImportPublicKeyFunc = f
	return m.mock
}

//ImportPublicKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ImportPublicKey(p []byte) (r crypto.PublicKey, r1 error) {
	counter := atomic.AddUint64(&m.ImportPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ImportPublicKeyCounter, 1)

	if len(m.ImportPublicKeyMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ImportPublicKeyMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPublicKey. %v", p)
			return
		}

		input := m.ImportPublicKeyMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockImportPublicKeyInput{p}, "KeyProcessor.ImportPublicKey got unexpected parameters")

		result := m.ImportPublicKeyMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPublicKey")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPublicKeyMock.mainExpectation != nil {

		input := m.ImportPublicKeyMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockImportPublicKeyInput{p}, "KeyProcessor.ImportPublicKey got unexpected parameters")
		}

		result := m.ImportPublicKeyMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPublicKey")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPublicKeyFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPublicKey. %v", p)
		return
	}

	return m.ImportPublicKeyFunc(p)
}

//ImportPublicKeyMinimockCounter returns a count of KeyProcessorMock.ImportPublicKeyFunc invocations
func (m *KeyProcessorMock) ImportPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPublicKeyCounter)
}

//ImportPublicKeyMinimockPreCounter returns the value of KeyProcessorMock.ImportPublicKey invocations
func (m *KeyProcessorMock) ImportPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPublicKeyPreCounter)
}

//ImportPublicKeyFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ImportPublicKeyFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ImportPublicKeyMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ImportPublicKeyCounter) == uint64(len(m.ImportPublicKeyMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ImportPublicKeyMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ImportPublicKeyCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ImportPublicKeyFunc != nil {
		return atomic.LoadUint64(&m.ImportPublicKeyCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *KeyProcessorMock) ValidateCallCounters() {

	if !m.ExportPrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPrivateKey")
	}

	if !m.ExportPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKey")
	}

	if !m.ExtractPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExtractPublicKey")
	}

	if !m.GeneratePrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.GeneratePrivateKey")
	}

	if !m.ImportPrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPrivateKey")
	}

	if !m.ImportPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPublicKey")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *KeyProcessorMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *KeyProcessorMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *KeyProcessorMock) MinimockFinish() {

	if !m.ExportPrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPrivateKey")
	}

	if !m.ExportPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKey")
	}

	if !m.ExtractPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExtractPublicKey")
	}

	if !m.GeneratePrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.GeneratePrivateKey")
	}

	if !m.ImportPrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPrivateKey")
	}

	if !m.ImportPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPublicKey")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *KeyProcessorMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *KeyProcessorMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ExportPrivateKeyFinished()
		ok = ok && m.ExportPublicKeyFinished()
		ok = ok && m.ExtractPublicKeyFinished()
		ok = ok && m.GeneratePrivateKeyFinished()
		ok = ok && m.ImportPrivateKeyFinished()
		ok = ok && m.ImportPublicKeyFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ExportPrivateKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ExportPrivateKey")
			}

			if !m.ExportPublicKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ExportPublicKey")
			}

			if !m.ExtractPublicKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ExtractPublicKey")
			}

			if !m.GeneratePrivateKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.GeneratePrivateKey")
			}

			if !m.ImportPrivateKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ImportPrivateKey")
			}

			if !m.ImportPublicKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ImportPublicKey")
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
func (m *KeyProcessorMock) AllMocksCalled() bool {

	if !m.ExportPrivateKeyFinished() {
		return false
	}

	if !m.ExportPublicKeyFinished() {
		return false
	}

	if !m.ExtractPublicKeyFinished() {
		return false
	}

	if !m.GeneratePrivateKeyFinished() {
		return false
	}

	if !m.ImportPrivateKeyFinished() {
		return false
	}

	if !m.ImportPublicKeyFinished() {
		return false
	}

	return true
}
