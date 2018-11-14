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
	mock             *KeyProcessorMock
	mockExpectations *KeyProcessorMockExportPrivateKeyParams
}

//KeyProcessorMockExportPrivateKeyParams represents input parameters of the KeyProcessor.ExportPrivateKey
type KeyProcessorMockExportPrivateKeyParams struct {
	p crypto.PrivateKey
}

//Expect sets up expected params for the KeyProcessor.ExportPrivateKey
func (m *mKeyProcessorMockExportPrivateKey) Expect(p crypto.PrivateKey) *mKeyProcessorMockExportPrivateKey {
	m.mockExpectations = &KeyProcessorMockExportPrivateKeyParams{p}
	return m
}

//Return sets up a mock for KeyProcessor.ExportPrivateKey to return Return's arguments
func (m *mKeyProcessorMockExportPrivateKey) Return(r []byte, r1 error) *KeyProcessorMock {
	m.mock.ExportPrivateKeyFunc = func(p crypto.PrivateKey) ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of KeyProcessor.ExportPrivateKey method
func (m *mKeyProcessorMockExportPrivateKey) Set(f func(p crypto.PrivateKey) (r []byte, r1 error)) *KeyProcessorMock {
	m.mock.ExportPrivateKeyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ExportPrivateKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExportPrivateKey(p crypto.PrivateKey) (r []byte, r1 error) {
	atomic.AddUint64(&m.ExportPrivateKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ExportPrivateKeyCounter, 1)

	if m.ExportPrivateKeyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExportPrivateKeyMock.mockExpectations, KeyProcessorMockExportPrivateKeyParams{p},
			"KeyProcessor.ExportPrivateKey got unexpected parameters")

		if m.ExportPrivateKeyFunc == nil {

			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPrivateKey")

			return
		}
	}

	if m.ExportPrivateKeyFunc == nil {
		m.t.Fatal("Unexpected call to KeyProcessorMock.ExportPrivateKey")
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

type mKeyProcessorMockExportPublicKey struct {
	mock             *KeyProcessorMock
	mockExpectations *KeyProcessorMockExportPublicKeyParams
}

//KeyProcessorMockExportPublicKeyParams represents input parameters of the KeyProcessor.ExportPublicKey
type KeyProcessorMockExportPublicKeyParams struct {
	p crypto.PublicKey
}

//Expect sets up expected params for the KeyProcessor.ExportPublicKey
func (m *mKeyProcessorMockExportPublicKey) Expect(p crypto.PublicKey) *mKeyProcessorMockExportPublicKey {
	m.mockExpectations = &KeyProcessorMockExportPublicKeyParams{p}
	return m
}

//Return sets up a mock for KeyProcessor.ExportPublicKey to return Return's arguments
func (m *mKeyProcessorMockExportPublicKey) Return(r []byte, r1 error) *KeyProcessorMock {
	m.mock.ExportPublicKeyFunc = func(p crypto.PublicKey) ([]byte, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of KeyProcessor.ExportPublicKey method
func (m *mKeyProcessorMockExportPublicKey) Set(f func(p crypto.PublicKey) (r []byte, r1 error)) *KeyProcessorMock {
	m.mock.ExportPublicKeyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ExportPublicKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExportPublicKey(p crypto.PublicKey) (r []byte, r1 error) {
	atomic.AddUint64(&m.ExportPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ExportPublicKeyCounter, 1)

	if m.ExportPublicKeyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExportPublicKeyMock.mockExpectations, KeyProcessorMockExportPublicKeyParams{p},
			"KeyProcessor.ExportPublicKey got unexpected parameters")

		if m.ExportPublicKeyFunc == nil {

			m.t.Fatal("No results are set for the KeyProcessorMock.ExportPublicKey")

			return
		}
	}

	if m.ExportPublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to KeyProcessorMock.ExportPublicKey")
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

type mKeyProcessorMockExtractPublicKey struct {
	mock             *KeyProcessorMock
	mockExpectations *KeyProcessorMockExtractPublicKeyParams
}

//KeyProcessorMockExtractPublicKeyParams represents input parameters of the KeyProcessor.ExtractPublicKey
type KeyProcessorMockExtractPublicKeyParams struct {
	p crypto.PrivateKey
}

//Expect sets up expected params for the KeyProcessor.ExtractPublicKey
func (m *mKeyProcessorMockExtractPublicKey) Expect(p crypto.PrivateKey) *mKeyProcessorMockExtractPublicKey {
	m.mockExpectations = &KeyProcessorMockExtractPublicKeyParams{p}
	return m
}

//Return sets up a mock for KeyProcessor.ExtractPublicKey to return Return's arguments
func (m *mKeyProcessorMockExtractPublicKey) Return(r crypto.PublicKey) *KeyProcessorMock {
	m.mock.ExtractPublicKeyFunc = func(p crypto.PrivateKey) crypto.PublicKey {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of KeyProcessor.ExtractPublicKey method
func (m *mKeyProcessorMockExtractPublicKey) Set(f func(p crypto.PrivateKey) (r crypto.PublicKey)) *KeyProcessorMock {
	m.mock.ExtractPublicKeyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ExtractPublicKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ExtractPublicKey(p crypto.PrivateKey) (r crypto.PublicKey) {
	atomic.AddUint64(&m.ExtractPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ExtractPublicKeyCounter, 1)

	if m.ExtractPublicKeyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ExtractPublicKeyMock.mockExpectations, KeyProcessorMockExtractPublicKeyParams{p},
			"KeyProcessor.ExtractPublicKey got unexpected parameters")

		if m.ExtractPublicKeyFunc == nil {

			m.t.Fatal("No results are set for the KeyProcessorMock.ExtractPublicKey")

			return
		}
	}

	if m.ExtractPublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to KeyProcessorMock.ExtractPublicKey")
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

type mKeyProcessorMockGeneratePrivateKey struct {
	mock *KeyProcessorMock
}

//Return sets up a mock for KeyProcessor.GeneratePrivateKey to return Return's arguments
func (m *mKeyProcessorMockGeneratePrivateKey) Return(r crypto.PrivateKey, r1 error) *KeyProcessorMock {
	m.mock.GeneratePrivateKeyFunc = func() (crypto.PrivateKey, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of KeyProcessor.GeneratePrivateKey method
func (m *mKeyProcessorMockGeneratePrivateKey) Set(f func() (r crypto.PrivateKey, r1 error)) *KeyProcessorMock {
	m.mock.GeneratePrivateKeyFunc = f

	return m.mock
}

//GeneratePrivateKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) GeneratePrivateKey() (r crypto.PrivateKey, r1 error) {
	atomic.AddUint64(&m.GeneratePrivateKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GeneratePrivateKeyCounter, 1)

	if m.GeneratePrivateKeyFunc == nil {
		m.t.Fatal("Unexpected call to KeyProcessorMock.GeneratePrivateKey")
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

type mKeyProcessorMockImportPrivateKey struct {
	mock             *KeyProcessorMock
	mockExpectations *KeyProcessorMockImportPrivateKeyParams
}

//KeyProcessorMockImportPrivateKeyParams represents input parameters of the KeyProcessor.ImportPrivateKey
type KeyProcessorMockImportPrivateKeyParams struct {
	p []byte
}

//Expect sets up expected params for the KeyProcessor.ImportPrivateKey
func (m *mKeyProcessorMockImportPrivateKey) Expect(p []byte) *mKeyProcessorMockImportPrivateKey {
	m.mockExpectations = &KeyProcessorMockImportPrivateKeyParams{p}
	return m
}

//Return sets up a mock for KeyProcessor.ImportPrivateKey to return Return's arguments
func (m *mKeyProcessorMockImportPrivateKey) Return(r crypto.PrivateKey, r1 error) *KeyProcessorMock {
	m.mock.ImportPrivateKeyFunc = func(p []byte) (crypto.PrivateKey, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of KeyProcessor.ImportPrivateKey method
func (m *mKeyProcessorMockImportPrivateKey) Set(f func(p []byte) (r crypto.PrivateKey, r1 error)) *KeyProcessorMock {
	m.mock.ImportPrivateKeyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ImportPrivateKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ImportPrivateKey(p []byte) (r crypto.PrivateKey, r1 error) {
	atomic.AddUint64(&m.ImportPrivateKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ImportPrivateKeyCounter, 1)

	if m.ImportPrivateKeyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ImportPrivateKeyMock.mockExpectations, KeyProcessorMockImportPrivateKeyParams{p},
			"KeyProcessor.ImportPrivateKey got unexpected parameters")

		if m.ImportPrivateKeyFunc == nil {

			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPrivateKey")

			return
		}
	}

	if m.ImportPrivateKeyFunc == nil {
		m.t.Fatal("Unexpected call to KeyProcessorMock.ImportPrivateKey")
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

type mKeyProcessorMockImportPublicKey struct {
	mock             *KeyProcessorMock
	mockExpectations *KeyProcessorMockImportPublicKeyParams
}

//KeyProcessorMockImportPublicKeyParams represents input parameters of the KeyProcessor.ImportPublicKey
type KeyProcessorMockImportPublicKeyParams struct {
	p []byte
}

//Expect sets up expected params for the KeyProcessor.ImportPublicKey
func (m *mKeyProcessorMockImportPublicKey) Expect(p []byte) *mKeyProcessorMockImportPublicKey {
	m.mockExpectations = &KeyProcessorMockImportPublicKeyParams{p}
	return m
}

//Return sets up a mock for KeyProcessor.ImportPublicKey to return Return's arguments
func (m *mKeyProcessorMockImportPublicKey) Return(r crypto.PublicKey, r1 error) *KeyProcessorMock {
	m.mock.ImportPublicKeyFunc = func(p []byte) (crypto.PublicKey, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of KeyProcessor.ImportPublicKey method
func (m *mKeyProcessorMockImportPublicKey) Set(f func(p []byte) (r crypto.PublicKey, r1 error)) *KeyProcessorMock {
	m.mock.ImportPublicKeyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//ImportPublicKey implements github.com/insolar/insolar/core.KeyProcessor interface
func (m *KeyProcessorMock) ImportPublicKey(p []byte) (r crypto.PublicKey, r1 error) {
	atomic.AddUint64(&m.ImportPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.ImportPublicKeyCounter, 1)

	if m.ImportPublicKeyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ImportPublicKeyMock.mockExpectations, KeyProcessorMockImportPublicKeyParams{p},
			"KeyProcessor.ImportPublicKey got unexpected parameters")

		if m.ImportPublicKeyFunc == nil {

			m.t.Fatal("No results are set for the KeyProcessorMock.ImportPublicKey")

			return
		}
	}

	if m.ImportPublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to KeyProcessorMock.ImportPublicKey")
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

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *KeyProcessorMock) ValidateCallCounters() {

	if m.ExportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ExportPrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPrivateKey")
	}

	if m.ExportPublicKeyFunc != nil && atomic.LoadUint64(&m.ExportPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKey")
	}

	if m.ExtractPublicKeyFunc != nil && atomic.LoadUint64(&m.ExtractPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ExtractPublicKey")
	}

	if m.GeneratePrivateKeyFunc != nil && atomic.LoadUint64(&m.GeneratePrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.GeneratePrivateKey")
	}

	if m.ImportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ImportPrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPrivateKey")
	}

	if m.ImportPublicKeyFunc != nil && atomic.LoadUint64(&m.ImportPublicKeyCounter) == 0 {
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

	if m.ExportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ExportPrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPrivateKey")
	}

	if m.ExportPublicKeyFunc != nil && atomic.LoadUint64(&m.ExportPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ExportPublicKey")
	}

	if m.ExtractPublicKeyFunc != nil && atomic.LoadUint64(&m.ExtractPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ExtractPublicKey")
	}

	if m.GeneratePrivateKeyFunc != nil && atomic.LoadUint64(&m.GeneratePrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.GeneratePrivateKey")
	}

	if m.ImportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ImportPrivateKeyCounter) == 0 {
		m.t.Fatal("Expected call to KeyProcessorMock.ImportPrivateKey")
	}

	if m.ImportPublicKeyFunc != nil && atomic.LoadUint64(&m.ImportPublicKeyCounter) == 0 {
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
		ok = ok && (m.ExportPrivateKeyFunc == nil || atomic.LoadUint64(&m.ExportPrivateKeyCounter) > 0)
		ok = ok && (m.ExportPublicKeyFunc == nil || atomic.LoadUint64(&m.ExportPublicKeyCounter) > 0)
		ok = ok && (m.ExtractPublicKeyFunc == nil || atomic.LoadUint64(&m.ExtractPublicKeyCounter) > 0)
		ok = ok && (m.GeneratePrivateKeyFunc == nil || atomic.LoadUint64(&m.GeneratePrivateKeyCounter) > 0)
		ok = ok && (m.ImportPrivateKeyFunc == nil || atomic.LoadUint64(&m.ImportPrivateKeyCounter) > 0)
		ok = ok && (m.ImportPublicKeyFunc == nil || atomic.LoadUint64(&m.ImportPublicKeyCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ExportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ExportPrivateKeyCounter) == 0 {
				m.t.Error("Expected call to KeyProcessorMock.ExportPrivateKey")
			}

			if m.ExportPublicKeyFunc != nil && atomic.LoadUint64(&m.ExportPublicKeyCounter) == 0 {
				m.t.Error("Expected call to KeyProcessorMock.ExportPublicKey")
			}

			if m.ExtractPublicKeyFunc != nil && atomic.LoadUint64(&m.ExtractPublicKeyCounter) == 0 {
				m.t.Error("Expected call to KeyProcessorMock.ExtractPublicKey")
			}

			if m.GeneratePrivateKeyFunc != nil && atomic.LoadUint64(&m.GeneratePrivateKeyCounter) == 0 {
				m.t.Error("Expected call to KeyProcessorMock.GeneratePrivateKey")
			}

			if m.ImportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ImportPrivateKeyCounter) == 0 {
				m.t.Error("Expected call to KeyProcessorMock.ImportPrivateKey")
			}

			if m.ImportPublicKeyFunc != nil && atomic.LoadUint64(&m.ImportPublicKeyCounter) == 0 {
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

	if m.ExportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ExportPrivateKeyCounter) == 0 {
		return false
	}

	if m.ExportPublicKeyFunc != nil && atomic.LoadUint64(&m.ExportPublicKeyCounter) == 0 {
		return false
	}

	if m.ExtractPublicKeyFunc != nil && atomic.LoadUint64(&m.ExtractPublicKeyCounter) == 0 {
		return false
	}

	if m.GeneratePrivateKeyFunc != nil && atomic.LoadUint64(&m.GeneratePrivateKeyCounter) == 0 {
		return false
	}

	if m.ImportPrivateKeyFunc != nil && atomic.LoadUint64(&m.ImportPrivateKeyCounter) == 0 {
		return false
	}

	if m.ImportPublicKeyFunc != nil && atomic.LoadUint64(&m.ImportPublicKeyCounter) == 0 {
		return false
	}

	return true
}
