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

	ExportPrivateKeyPEMFunc       func(p crypto.PrivateKey) (r []byte, r1 error)
	ExportPrivateKeyPEMCounter    uint64
	ExportPrivateKeyPEMPreCounter uint64
	ExportPrivateKeyPEMMock       mKeyProcessorMockExportPrivateKeyPEM

	ExportPublicKeyBinaryFunc       func(p crypto.PublicKey) (r []byte, r1 error)
	ExportPublicKeyBinaryCounter    uint64
	ExportPublicKeyBinaryPreCounter uint64
	ExportPublicKeyBinaryMock       mKeyProcessorMockExportPublicKeyBinary

	ExportPublicKeyPEMFunc       func(p crypto.PublicKey) (r []byte, r1 error)
	ExportPublicKeyPEMCounter    uint64
	ExportPublicKeyPEMPreCounter uint64
	ExportPublicKeyPEMMock       mKeyProcessorMockExportPublicKeyPEM

	ExtractPublicKeyFunc       func(p crypto.PrivateKey) (r crypto.PublicKey)
	ExtractPublicKeyCounter    uint64
	ExtractPublicKeyPreCounter uint64
	ExtractPublicKeyMock       mKeyProcessorMockExtractPublicKey

	GeneratePrivateKeyFunc       func() (r crypto.PrivateKey, r1 error)
	GeneratePrivateKeyCounter    uint64
	GeneratePrivateKeyPreCounter uint64
	GeneratePrivateKeyMock       mKeyProcessorMockGeneratePrivateKey

	ImportPrivateKeyPEMFunc       func(p []byte) (r crypto.PrivateKey, r1 error)
	ImportPrivateKeyPEMCounter    uint64
	ImportPrivateKeyPEMPreCounter uint64
	ImportPrivateKeyPEMMock       mKeyProcessorMockImportPrivateKeyPEM

	ImportPublicKeyBinaryFunc       func(p []byte) (r crypto.PublicKey, r1 error)
	ImportPublicKeyBinaryCounter    uint64
	ImportPublicKeyBinaryPreCounter uint64
	ImportPublicKeyBinaryMock       mKeyProcessorMockImportPublicKeyBinary

	ImportPublicKeyPEMFunc       func(p []byte) (r crypto.PublicKey, r1 error)
	ImportPublicKeyPEMCounter    uint64
	ImportPublicKeyPEMPreCounter uint64
	ImportPublicKeyPEMMock       mKeyProcessorMockImportPublicKeyPEM
}

//NewKeyProcessorMock returns a mock for github.com/insolar/insolar/core.KeyProcessor
func NewKeyProcessorMock(t minimock.Tester) *KeyProcessorMock {
	m := &KeyProcessorMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExportPrivateKeyPEMMock = mKeyProcessorMockExportPrivateKeyPEM{mock: m}
	m.ExportPublicKeyBinaryMock = mKeyProcessorMockExportPublicKeyBinary{mock: m}
	m.ExportPublicKeyPEMMock = mKeyProcessorMockExportPublicKeyPEM{mock: m}
	m.ExtractPublicKeyMock = mKeyProcessorMockExtractPublicKey{mock: m}
	m.GeneratePrivateKeyMock = mKeyProcessorMockGeneratePrivateKey{mock: m}
	m.ImportPrivateKeyPEMMock = mKeyProcessorMockImportPrivateKeyPEM{mock: m}
	m.ImportPublicKeyBinaryMock = mKeyProcessorMockImportPublicKeyBinary{mock: m}
	m.ImportPublicKeyPEMMock = mKeyProcessorMockImportPublicKeyPEM{mock: m}

	return m
}

type mKeyProcessorMockExportPrivateKeyPEM struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockExportPrivateKeyPEMExpectation
	expectationSeries []*KeyProcessorMockExportPrivateKeyPEMExpectation
}

type KeyProcessorMockExportPrivateKeyPEMExpectation struct {
	input  *KeyProcessorMockExportPrivateKeyPEMInput
	result *KeyProcessorMockExportPrivateKeyPEMResult
}

type KeyProcessorMockExportPrivateKeyPEMInput struct {
	p crypto.PrivateKey
}

type KeyProcessorMockExportPrivateKeyPEMResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ExportPrivateKeyPEM is expected from 1 to Infinity times
func (m *mKeyProcessorMockExportPrivateKeyPEM) Expect(p crypto.PrivateKey) *mKeyProcessorMockExportPrivateKeyPEM {
	m.mock.ExportPrivateKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPrivateKeyPEMExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockExportPrivateKeyPEMInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ExportPrivateKeyPEM
func (m *mKeyProcessorMockExportPrivateKeyPEM) Return(r []byte, r1 error) *KeyProcessorMock {
	m.mock.ExportPrivateKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPrivateKeyPEMExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockExportPrivateKeyPEMResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ExportPrivateKeyPEM is expected once
func (m *mKeyProcessorMockExportPrivateKeyPEM) ExpectOnce(p crypto.PrivateKey) *KeyProcessorMockExportPrivateKeyPEMExpectation {
	m.mock.ExportPrivateKeyPEMFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockExportPrivateKeyPEMExpectation{}
	expectation.input = &KeyProcessorMockExportPrivateKeyPEMInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockExportPrivateKeyPEMExpectation) Return(r []byte, r1 error) {
	e.result = &KeyProcessorMockExportPrivateKeyPEMResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ExportPrivateKeyPEM method
func (m *mKeyProcessorMockExportPrivateKeyPEM) Set(f func(p crypto.PrivateKey) (r []byte, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExportPrivateKeyPEMFunc = f
	return m.mock
}

//ExportPrivateKeyPEM implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExportPrivateKeyPEM(p crypto.PrivateKey) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.ExportPrivateKeyPEMPreCounter, 1)
	defer atomic.AddUint64(&m.ExportPrivateKeyPEMCounter, 1)

	if len(m.ExportPrivateKeyPEMMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExportPrivateKeyPEMMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPrivateKeyPEM. %v", p)
			return
		}

		input := m.ExportPrivateKeyPEMMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockExportPrivateKeyPEMInput{p}, "KeyProcessor.ExportPrivateKeyPEM got unexpected parameters")

		result := m.ExportPrivateKeyPEMMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPrivateKeyPEM")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPrivateKeyPEMMock.mainExpectation != nil {

		input := m.ExportPrivateKeyPEMMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockExportPrivateKeyPEMInput{p}, "KeyProcessor.ExportPrivateKeyPEM got unexpected parameters")
		}

		result := m.ExportPrivateKeyPEMMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPrivateKeyPEM")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPrivateKeyPEMFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPrivateKeyPEM. %v", p)
		return
	}

	return m.ExportPrivateKeyPEMFunc(p)
}

//ExportPrivateKeyPEMMinimockCounter returns a count of KeyProcessorMock.ExportPrivateKeyPEMFunc invocations
func (m *KeyProcessorMock) ExportPrivateKeyPEMMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPrivateKeyPEMCounter)
}

//ExportPrivateKeyPEMMinimockPreCounter returns the value of KeyProcessorMock.ExportPrivateKeyPEM invocations
func (m *KeyProcessorMock) ExportPrivateKeyPEMMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPrivateKeyPEMPreCounter)
}

//ExportPrivateKeyPEMFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ExportPrivateKeyPEMFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExportPrivateKeyPEMMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExportPrivateKeyPEMCounter) == uint64(len(m.ExportPrivateKeyPEMMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExportPrivateKeyPEMMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExportPrivateKeyPEMCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExportPrivateKeyPEMFunc != nil {
		return atomic.LoadUint64(&m.ExportPrivateKeyPEMCounter) > 0
	}

	return true
}

type mKeyProcessorMockExportPublicKeyBinary struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockExportPublicKeyBinaryExpectation
	expectationSeries []*KeyProcessorMockExportPublicKeyBinaryExpectation
}

type KeyProcessorMockExportPublicKeyBinaryExpectation struct {
	input  *KeyProcessorMockExportPublicKeyBinaryInput
	result *KeyProcessorMockExportPublicKeyBinaryResult
}

type KeyProcessorMockExportPublicKeyBinaryInput struct {
	p crypto.PublicKey
}

type KeyProcessorMockExportPublicKeyBinaryResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ExportPublicKeyBinary is expected from 1 to Infinity times
func (m *mKeyProcessorMockExportPublicKeyBinary) Expect(p crypto.PublicKey) *mKeyProcessorMockExportPublicKeyBinary {
	m.mock.ExportPublicKeyBinaryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPublicKeyBinaryExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockExportPublicKeyBinaryInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ExportPublicKeyBinary
func (m *mKeyProcessorMockExportPublicKeyBinary) Return(r []byte, r1 error) *KeyProcessorMock {
	m.mock.ExportPublicKeyBinaryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPublicKeyBinaryExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockExportPublicKeyBinaryResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ExportPublicKeyBinary is expected once
func (m *mKeyProcessorMockExportPublicKeyBinary) ExpectOnce(p crypto.PublicKey) *KeyProcessorMockExportPublicKeyBinaryExpectation {
	m.mock.ExportPublicKeyBinaryFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockExportPublicKeyBinaryExpectation{}
	expectation.input = &KeyProcessorMockExportPublicKeyBinaryInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockExportPublicKeyBinaryExpectation) Return(r []byte, r1 error) {
	e.result = &KeyProcessorMockExportPublicKeyBinaryResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ExportPublicKeyBinary method
func (m *mKeyProcessorMockExportPublicKeyBinary) Set(f func(p crypto.PublicKey) (r []byte, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExportPublicKeyBinaryFunc = f
	return m.mock
}

//ExportPublicKeyBinary implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExportPublicKeyBinary(p crypto.PublicKey) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.ExportPublicKeyBinaryPreCounter, 1)
	defer atomic.AddUint64(&m.ExportPublicKeyBinaryCounter, 1)

	if len(m.ExportPublicKeyBinaryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExportPublicKeyBinaryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPublicKeyBinary. %v", p)
			return
		}

		input := m.ExportPublicKeyBinaryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockExportPublicKeyBinaryInput{p}, "KeyProcessor.ExportPublicKeyBinary got unexpected parameters")

		result := m.ExportPublicKeyBinaryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPublicKeyBinary")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPublicKeyBinaryMock.mainExpectation != nil {

		input := m.ExportPublicKeyBinaryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockExportPublicKeyBinaryInput{p}, "KeyProcessor.ExportPublicKeyBinary got unexpected parameters")
		}

		result := m.ExportPublicKeyBinaryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPublicKeyBinary")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPublicKeyBinaryFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPublicKeyBinary. %v", p)
		return
	}

	return m.ExportPublicKeyBinaryFunc(p)
}

//ExportPublicKeyBinaryMinimockCounter returns a count of KeyProcessorMock.ExportPublicKeyBinaryFunc invocations
func (m *KeyProcessorMock) ExportPublicKeyBinaryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPublicKeyBinaryCounter)
}

//ExportPublicKeyBinaryMinimockPreCounter returns the value of KeyProcessorMock.ExportPublicKeyBinary invocations
func (m *KeyProcessorMock) ExportPublicKeyBinaryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPublicKeyBinaryPreCounter)
}

//ExportPublicKeyBinaryFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ExportPublicKeyBinaryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExportPublicKeyBinaryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExportPublicKeyBinaryCounter) == uint64(len(m.ExportPublicKeyBinaryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExportPublicKeyBinaryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExportPublicKeyBinaryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExportPublicKeyBinaryFunc != nil {
		return atomic.LoadUint64(&m.ExportPublicKeyBinaryCounter) > 0
	}

	return true
}

type mKeyProcessorMockExportPublicKeyPEM struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockExportPublicKeyPEMExpectation
	expectationSeries []*KeyProcessorMockExportPublicKeyPEMExpectation
}

type KeyProcessorMockExportPublicKeyPEMExpectation struct {
	input  *KeyProcessorMockExportPublicKeyPEMInput
	result *KeyProcessorMockExportPublicKeyPEMResult
}

type KeyProcessorMockExportPublicKeyPEMInput struct {
	p crypto.PublicKey
}

type KeyProcessorMockExportPublicKeyPEMResult struct {
	r  []byte
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ExportPublicKeyPEM is expected from 1 to Infinity times
func (m *mKeyProcessorMockExportPublicKeyPEM) Expect(p crypto.PublicKey) *mKeyProcessorMockExportPublicKeyPEM {
	m.mock.ExportPublicKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPublicKeyPEMExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockExportPublicKeyPEMInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ExportPublicKeyPEM
func (m *mKeyProcessorMockExportPublicKeyPEM) Return(r []byte, r1 error) *KeyProcessorMock {
	m.mock.ExportPublicKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockExportPublicKeyPEMExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockExportPublicKeyPEMResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ExportPublicKeyPEM is expected once
func (m *mKeyProcessorMockExportPublicKeyPEM) ExpectOnce(p crypto.PublicKey) *KeyProcessorMockExportPublicKeyPEMExpectation {
	m.mock.ExportPublicKeyPEMFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockExportPublicKeyPEMExpectation{}
	expectation.input = &KeyProcessorMockExportPublicKeyPEMInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockExportPublicKeyPEMExpectation) Return(r []byte, r1 error) {
	e.result = &KeyProcessorMockExportPublicKeyPEMResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ExportPublicKeyPEM method
func (m *mKeyProcessorMockExportPublicKeyPEM) Set(f func(p crypto.PublicKey) (r []byte, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ExportPublicKeyPEMFunc = f
	return m.mock
}

//ExportPublicKeyPEM implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExportPublicKeyPEM(p crypto.PublicKey) (r []byte, r1 error) {
	counter := atomic.AddUint64(&m.ExportPublicKeyPEMPreCounter, 1)
	defer atomic.AddUint64(&m.ExportPublicKeyPEMCounter, 1)

	if len(m.ExportPublicKeyPEMMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ExportPublicKeyPEMMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPublicKeyPEM. %v", p)
			return
		}

		input := m.ExportPublicKeyPEMMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockExportPublicKeyPEMInput{p}, "KeyProcessor.ExportPublicKeyPEM got unexpected parameters")

		result := m.ExportPublicKeyPEMMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPublicKeyPEM")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPublicKeyPEMMock.mainExpectation != nil {

		input := m.ExportPublicKeyPEMMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockExportPublicKeyPEMInput{p}, "KeyProcessor.ExportPublicKeyPEM got unexpected parameters")
		}

		result := m.ExportPublicKeyPEMMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPublicKeyPEM")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ExportPublicKeyPEMFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ExportPublicKeyPEM. %v", p)
		return
	}

	return m.ExportPublicKeyPEMFunc(p)
}

//ExportPublicKeyPEMMinimockCounter returns a count of KeyProcessorMock.ExportPublicKeyPEMFunc invocations
func (m *KeyProcessorMock) ExportPublicKeyPEMMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPublicKeyPEMCounter)
}

//ExportPublicKeyPEMMinimockPreCounter returns the value of KeyProcessorMock.ExportPublicKeyPEM invocations
func (m *KeyProcessorMock) ExportPublicKeyPEMMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ExportPublicKeyPEMPreCounter)
}

//ExportPublicKeyPEMFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ExportPublicKeyPEMFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ExportPublicKeyPEMMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ExportPublicKeyPEMCounter) == uint64(len(m.ExportPublicKeyPEMMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ExportPublicKeyPEMMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ExportPublicKeyPEMCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ExportPublicKeyPEMFunc != nil {
		return atomic.LoadUint64(&m.ExportPublicKeyPEMCounter) > 0
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

type mKeyProcessorMockImportPrivateKeyPEM struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockImportPrivateKeyPEMExpectation
	expectationSeries []*KeyProcessorMockImportPrivateKeyPEMExpectation
}

type KeyProcessorMockImportPrivateKeyPEMExpectation struct {
	input  *KeyProcessorMockImportPrivateKeyPEMInput
	result *KeyProcessorMockImportPrivateKeyPEMResult
}

type KeyProcessorMockImportPrivateKeyPEMInput struct {
	p []byte
}

type KeyProcessorMockImportPrivateKeyPEMResult struct {
	r  crypto.PrivateKey
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ImportPrivateKeyPEM is expected from 1 to Infinity times
func (m *mKeyProcessorMockImportPrivateKeyPEM) Expect(p []byte) *mKeyProcessorMockImportPrivateKeyPEM {
	m.mock.ImportPrivateKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPrivateKeyPEMExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockImportPrivateKeyPEMInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ImportPrivateKeyPEM
func (m *mKeyProcessorMockImportPrivateKeyPEM) Return(r crypto.PrivateKey, r1 error) *KeyProcessorMock {
	m.mock.ImportPrivateKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPrivateKeyPEMExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockImportPrivateKeyPEMResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ImportPrivateKeyPEM is expected once
func (m *mKeyProcessorMockImportPrivateKeyPEM) ExpectOnce(p []byte) *KeyProcessorMockImportPrivateKeyPEMExpectation {
	m.mock.ImportPrivateKeyPEMFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockImportPrivateKeyPEMExpectation{}
	expectation.input = &KeyProcessorMockImportPrivateKeyPEMInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockImportPrivateKeyPEMExpectation) Return(r crypto.PrivateKey, r1 error) {
	e.result = &KeyProcessorMockImportPrivateKeyPEMResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ImportPrivateKeyPEM method
func (m *mKeyProcessorMockImportPrivateKeyPEM) Set(f func(p []byte) (r crypto.PrivateKey, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ImportPrivateKeyPEMFunc = f
	return m.mock
}

//ImportPrivateKeyPEM implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ImportPrivateKeyPEM(p []byte) (r crypto.PrivateKey, r1 error) {
	counter := atomic.AddUint64(&m.ImportPrivateKeyPEMPreCounter, 1)
	defer atomic.AddUint64(&m.ImportPrivateKeyPEMCounter, 1)

	if len(m.ImportPrivateKeyPEMMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ImportPrivateKeyPEMMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPrivateKeyPEM. %v", p)
			return
		}

		input := m.ImportPrivateKeyPEMMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockImportPrivateKeyPEMInput{p}, "KeyProcessor.ImportPrivateKeyPEM got unexpected parameters")

		result := m.ImportPrivateKeyPEMMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPrivateKeyPEM")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPrivateKeyPEMMock.mainExpectation != nil {

		input := m.ImportPrivateKeyPEMMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockImportPrivateKeyPEMInput{p}, "KeyProcessor.ImportPrivateKeyPEM got unexpected parameters")
		}

		result := m.ImportPrivateKeyPEMMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPrivateKeyPEM")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPrivateKeyPEMFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPrivateKeyPEM. %v", p)
		return
	}

	return m.ImportPrivateKeyPEMFunc(p)
}

//ImportPrivateKeyPEMMinimockCounter returns a count of KeyProcessorMock.ImportPrivateKeyPEMFunc invocations
func (m *KeyProcessorMock) ImportPrivateKeyPEMMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPrivateKeyPEMCounter)
}

//ImportPrivateKeyPEMMinimockPreCounter returns the value of KeyProcessorMock.ImportPrivateKeyPEM invocations
func (m *KeyProcessorMock) ImportPrivateKeyPEMMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPrivateKeyPEMPreCounter)
}

//ImportPrivateKeyPEMFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ImportPrivateKeyPEMFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ImportPrivateKeyPEMMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ImportPrivateKeyPEMCounter) == uint64(len(m.ImportPrivateKeyPEMMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ImportPrivateKeyPEMMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ImportPrivateKeyPEMCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ImportPrivateKeyPEMFunc != nil {
		return atomic.LoadUint64(&m.ImportPrivateKeyPEMCounter) > 0
	}

	return true
}

type mKeyProcessorMockImportPublicKeyBinary struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockImportPublicKeyBinaryExpectation
	expectationSeries []*KeyProcessorMockImportPublicKeyBinaryExpectation
}

type KeyProcessorMockImportPublicKeyBinaryExpectation struct {
	input  *KeyProcessorMockImportPublicKeyBinaryInput
	result *KeyProcessorMockImportPublicKeyBinaryResult
}

type KeyProcessorMockImportPublicKeyBinaryInput struct {
	p []byte
}

type KeyProcessorMockImportPublicKeyBinaryResult struct {
	r  crypto.PublicKey
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ImportPublicKeyBinary is expected from 1 to Infinity times
func (m *mKeyProcessorMockImportPublicKeyBinary) Expect(p []byte) *mKeyProcessorMockImportPublicKeyBinary {
	m.mock.ImportPublicKeyBinaryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPublicKeyBinaryExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockImportPublicKeyBinaryInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ImportPublicKeyBinary
func (m *mKeyProcessorMockImportPublicKeyBinary) Return(r crypto.PublicKey, r1 error) *KeyProcessorMock {
	m.mock.ImportPublicKeyBinaryFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPublicKeyBinaryExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockImportPublicKeyBinaryResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ImportPublicKeyBinary is expected once
func (m *mKeyProcessorMockImportPublicKeyBinary) ExpectOnce(p []byte) *KeyProcessorMockImportPublicKeyBinaryExpectation {
	m.mock.ImportPublicKeyBinaryFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockImportPublicKeyBinaryExpectation{}
	expectation.input = &KeyProcessorMockImportPublicKeyBinaryInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockImportPublicKeyBinaryExpectation) Return(r crypto.PublicKey, r1 error) {
	e.result = &KeyProcessorMockImportPublicKeyBinaryResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ImportPublicKeyBinary method
func (m *mKeyProcessorMockImportPublicKeyBinary) Set(f func(p []byte) (r crypto.PublicKey, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ImportPublicKeyBinaryFunc = f
	return m.mock
}

//ImportPublicKeyBinary implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ImportPublicKeyBinary(p []byte) (r crypto.PublicKey, r1 error) {
	counter := atomic.AddUint64(&m.ImportPublicKeyBinaryPreCounter, 1)
	defer atomic.AddUint64(&m.ImportPublicKeyBinaryCounter, 1)

	if len(m.ImportPublicKeyBinaryMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ImportPublicKeyBinaryMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPublicKeyBinary. %v", p)
			return
		}

		input := m.ImportPublicKeyBinaryMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockImportPublicKeyBinaryInput{p}, "KeyProcessor.ImportPublicKeyBinary got unexpected parameters")

		result := m.ImportPublicKeyBinaryMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPublicKeyBinary")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPublicKeyBinaryMock.mainExpectation != nil {

		input := m.ImportPublicKeyBinaryMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockImportPublicKeyBinaryInput{p}, "KeyProcessor.ImportPublicKeyBinary got unexpected parameters")
		}

		result := m.ImportPublicKeyBinaryMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPublicKeyBinary")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPublicKeyBinaryFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPublicKeyBinary. %v", p)
		return
	}

	return m.ImportPublicKeyBinaryFunc(p)
}

//ImportPublicKeyBinaryMinimockCounter returns a count of KeyProcessorMock.ImportPublicKeyBinaryFunc invocations
func (m *KeyProcessorMock) ImportPublicKeyBinaryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPublicKeyBinaryCounter)
}

//ImportPublicKeyBinaryMinimockPreCounter returns the value of KeyProcessorMock.ImportPublicKeyBinary invocations
func (m *KeyProcessorMock) ImportPublicKeyBinaryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPublicKeyBinaryPreCounter)
}

//ImportPublicKeyBinaryFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ImportPublicKeyBinaryFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ImportPublicKeyBinaryMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ImportPublicKeyBinaryCounter) == uint64(len(m.ImportPublicKeyBinaryMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ImportPublicKeyBinaryMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ImportPublicKeyBinaryCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ImportPublicKeyBinaryFunc != nil {
		return atomic.LoadUint64(&m.ImportPublicKeyBinaryCounter) > 0
	}

	return true
}

type mKeyProcessorMockImportPublicKeyPEM struct {
	mock              *KeyProcessorMock
	mainExpectation   *KeyProcessorMockImportPublicKeyPEMExpectation
	expectationSeries []*KeyProcessorMockImportPublicKeyPEMExpectation
}

type KeyProcessorMockImportPublicKeyPEMExpectation struct {
	input  *KeyProcessorMockImportPublicKeyPEMInput
	result *KeyProcessorMockImportPublicKeyPEMResult
}

type KeyProcessorMockImportPublicKeyPEMInput struct {
	p []byte
}

type KeyProcessorMockImportPublicKeyPEMResult struct {
	r  crypto.PublicKey
	r1 error
}

//Expect specifies that invocation of KeyProcessor.ImportPublicKeyPEM is expected from 1 to Infinity times
func (m *mKeyProcessorMockImportPublicKeyPEM) Expect(p []byte) *mKeyProcessorMockImportPublicKeyPEM {
	m.mock.ImportPublicKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPublicKeyPEMExpectation{}
	}
	m.mainExpectation.input = &KeyProcessorMockImportPublicKeyPEMInput{p}
	return m
}

//Return specifies results of invocation of KeyProcessor.ImportPublicKeyPEM
func (m *mKeyProcessorMockImportPublicKeyPEM) Return(r crypto.PublicKey, r1 error) *KeyProcessorMock {
	m.mock.ImportPublicKeyPEMFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &KeyProcessorMockImportPublicKeyPEMExpectation{}
	}
	m.mainExpectation.result = &KeyProcessorMockImportPublicKeyPEMResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of KeyProcessor.ImportPublicKeyPEM is expected once
func (m *mKeyProcessorMockImportPublicKeyPEM) ExpectOnce(p []byte) *KeyProcessorMockImportPublicKeyPEMExpectation {
	m.mock.ImportPublicKeyPEMFunc = nil
	m.mainExpectation = nil

	expectation := &KeyProcessorMockImportPublicKeyPEMExpectation{}
	expectation.input = &KeyProcessorMockImportPublicKeyPEMInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *KeyProcessorMockImportPublicKeyPEMExpectation) Return(r crypto.PublicKey, r1 error) {
	e.result = &KeyProcessorMockImportPublicKeyPEMResult{r, r1}
}

//Set uses given function f as a mock of KeyProcessor.ImportPublicKeyPEM method
func (m *mKeyProcessorMockImportPublicKeyPEM) Set(f func(p []byte) (r crypto.PublicKey, r1 error)) *KeyProcessorMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ImportPublicKeyPEMFunc = f
	return m.mock
}

//ImportPublicKeyPEM implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ImportPublicKeyPEM(p []byte) (r crypto.PublicKey, r1 error) {
	counter := atomic.AddUint64(&m.ImportPublicKeyPEMPreCounter, 1)
	defer atomic.AddUint64(&m.ImportPublicKeyPEMCounter, 1)

	if len(m.ImportPublicKeyPEMMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ImportPublicKeyPEMMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPublicKeyPEM. %v", p)
			return
		}

		input := m.ImportPublicKeyPEMMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, KeyProcessorMockImportPublicKeyPEMInput{p}, "KeyProcessor.ImportPublicKeyPEM got unexpected parameters")

		result := m.ImportPublicKeyPEMMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPublicKeyPEM")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPublicKeyPEMMock.mainExpectation != nil {

		input := m.ImportPublicKeyPEMMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, KeyProcessorMockImportPublicKeyPEMInput{p}, "KeyProcessor.ImportPublicKeyPEM got unexpected parameters")
		}

		result := m.ImportPublicKeyPEMMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPublicKeyPEM")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ImportPublicKeyPEMFunc == nil {
		m.t.Fatalf("Unexpected call to KeyProcessorMock.ImportPublicKeyPEM. %v", p)
		return
	}

	return m.ImportPublicKeyPEMFunc(p)
}

//ImportPublicKeyPEMMinimockCounter returns a count of KeyProcessorMock.ImportPublicKeyPEMFunc invocations
func (m *KeyProcessorMock) ImportPublicKeyPEMMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPublicKeyPEMCounter)
}

//ImportPublicKeyPEMMinimockPreCounter returns the value of KeyProcessorMock.ImportPublicKeyPEM invocations
func (m *KeyProcessorMock) ImportPublicKeyPEMMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ImportPublicKeyPEMPreCounter)
}

//ImportPublicKeyPEMFinished returns true if mock invocations count is ok
func (m *KeyProcessorMock) ImportPublicKeyPEMFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ImportPublicKeyPEMMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ImportPublicKeyPEMCounter) == uint64(len(m.ImportPublicKeyPEMMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ImportPublicKeyPEMMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ImportPublicKeyPEMCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ImportPublicKeyPEMFunc != nil {
		return atomic.LoadUint64(&m.ImportPublicKeyPEMCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *KeyProcessorMock) ValidateCallCounters() {

	if !m.ExportPrivateKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPrivateKeyPEM")
	}

	if !m.ExportPublicKeyBinaryFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKeyBinary")
	}

	if !m.ExportPublicKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKeyPEM")
	}

	if !m.ExtractPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExtractPublicKey")
	}

	if !m.GeneratePrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.GeneratePrivateKey")
	}

	if !m.ImportPrivateKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPrivateKeyPEM")
	}

	if !m.ImportPublicKeyBinaryFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPublicKeyBinary")
	}

	if !m.ImportPublicKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPublicKeyPEM")
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

	if !m.ExportPrivateKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPrivateKeyPEM")
	}

	if !m.ExportPublicKeyBinaryFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKeyBinary")
	}

	if !m.ExportPublicKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKeyPEM")
	}

	if !m.ExtractPublicKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ExtractPublicKey")
	}

	if !m.GeneratePrivateKeyFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.GeneratePrivateKey")
	}

	if !m.ImportPrivateKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPrivateKeyPEM")
	}

	if !m.ImportPublicKeyBinaryFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPublicKeyBinary")
	}

	if !m.ImportPublicKeyPEMFinished() {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPublicKeyPEM")
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
		ok = ok && m.ExportPrivateKeyPEMFinished()
		ok = ok && m.ExportPublicKeyBinaryFinished()
		ok = ok && m.ExportPublicKeyPEMFinished()
		ok = ok && m.ExtractPublicKeyFinished()
		ok = ok && m.GeneratePrivateKeyFinished()
		ok = ok && m.ImportPrivateKeyPEMFinished()
		ok = ok && m.ImportPublicKeyBinaryFinished()
		ok = ok && m.ImportPublicKeyPEMFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ExportPrivateKeyPEMFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ExportPrivateKeyPEM")
			}

			if !m.ExportPublicKeyBinaryFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ExportPublicKeyBinary")
			}

			if !m.ExportPublicKeyPEMFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ExportPublicKeyPEM")
			}

			if !m.ExtractPublicKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ExtractPublicKey")
			}

			if !m.GeneratePrivateKeyFinished() {
				m.t.Error("Expected call to KeyProcessorMock.GeneratePrivateKey")
			}

			if !m.ImportPrivateKeyPEMFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ImportPrivateKeyPEM")
			}

			if !m.ImportPublicKeyBinaryFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ImportPublicKeyBinary")
			}

			if !m.ImportPublicKeyPEMFinished() {
				m.t.Error("Expected call to KeyProcessorMock.ImportPublicKeyPEM")
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

	if !m.ExportPrivateKeyPEMFinished() {
		return false
	}

	if !m.ExportPublicKeyBinaryFinished() {
		return false
	}

	if !m.ExportPublicKeyPEMFinished() {
		return false
	}

	if !m.ExtractPublicKeyFinished() {
		return false
	}

	if !m.GeneratePrivateKeyFinished() {
		return false
	}

	if !m.ImportPrivateKeyPEMFinished() {
		return false
	}

	if !m.ImportPublicKeyBinaryFinished() {
		return false
	}

	if !m.ImportPublicKeyPEMFinished() {
		return false
	}

	return true
}
