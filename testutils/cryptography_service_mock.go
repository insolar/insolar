package testutils

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "CryptographyService" can be found in github.com/insolar/insolar/core
*/
import (
	"crypto"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"
	testify_assert "github.com/stretchr/testify/assert"
)

//CryptographyServiceMock implements github.com/insolar/insolar/core.CryptographyService
type CryptographyServiceMock struct {
	t minimock.Tester

	GetPublicKeyFunc       func() (r crypto.PublicKey, r1 error)
	GetPublicKeyCounter    uint64
	GetPublicKeyPreCounter uint64
	GetPublicKeyMock       mCryptographyServiceMockGetPublicKey

	SignFunc       func(p []byte) (r *core.Signature, r1 error)
	SignCounter    uint64
	SignPreCounter uint64
	SignMock       mCryptographyServiceMockSign

	VerifyFunc       func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool)
	VerifyCounter    uint64
	VerifyPreCounter uint64
	VerifyMock       mCryptographyServiceMockVerify
}

//NewCryptographyServiceMock returns a mock for github.com/insolar/insolar/core.CryptographyService
func NewCryptographyServiceMock(t minimock.Tester) *CryptographyServiceMock {
	m := &CryptographyServiceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetPublicKeyMock = mCryptographyServiceMockGetPublicKey{mock: m}
	m.SignMock = mCryptographyServiceMockSign{mock: m}
	m.VerifyMock = mCryptographyServiceMockVerify{mock: m}

	return m
}

type mCryptographyServiceMockGetPublicKey struct {
	mock *CryptographyServiceMock
}

//Return sets up a mock for CryptographyService.GetPublicKey to return Return's arguments
func (m *mCryptographyServiceMockGetPublicKey) Return(r crypto.PublicKey, r1 error) *CryptographyServiceMock {
	m.mock.GetPublicKeyFunc = func() (crypto.PublicKey, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of CryptographyService.GetPublicKey method
func (m *mCryptographyServiceMockGetPublicKey) Set(f func() (r crypto.PublicKey, r1 error)) *CryptographyServiceMock {
	m.mock.GetPublicKeyFunc = f

	return m.mock
}

//GetPublicKey implements github.com/insolar/insolar/core.CryptographyService interface
func (m *CryptographyServiceMock) GetPublicKey() (r crypto.PublicKey, r1 error) {
	atomic.AddUint64(&m.GetPublicKeyPreCounter, 1)
	defer atomic.AddUint64(&m.GetPublicKeyCounter, 1)

	if m.GetPublicKeyFunc == nil {
		m.t.Fatal("Unexpected call to CryptographyServiceMock.GetPublicKey")
		return
	}

	return m.GetPublicKeyFunc()
}

//GetPublicKeyMinimockCounter returns a count of CryptographyServiceMock.GetPublicKeyFunc invocations
func (m *CryptographyServiceMock) GetPublicKeyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyCounter)
}

//GetPublicKeyMinimockPreCounter returns the value of CryptographyServiceMock.GetPublicKey invocations
func (m *CryptographyServiceMock) GetPublicKeyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetPublicKeyPreCounter)
}

type mCryptographyServiceMockSign struct {
	mock             *CryptographyServiceMock
	mockExpectations *CryptographyServiceMockSignParams
}

//CryptographyServiceMockSignParams represents input parameters of the CryptographyService.Sign
type CryptographyServiceMockSignParams struct {
	p []byte
}

//Expect sets up expected params for the CryptographyService.Sign
func (m *mCryptographyServiceMockSign) Expect(p []byte) *mCryptographyServiceMockSign {
	m.mockExpectations = &CryptographyServiceMockSignParams{p}
	return m
}

//Return sets up a mock for CryptographyService.Sign to return Return's arguments
func (m *mCryptographyServiceMockSign) Return(r *core.Signature, r1 error) *CryptographyServiceMock {
	m.mock.SignFunc = func(p []byte) (*core.Signature, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of CryptographyService.Sign method
func (m *mCryptographyServiceMockSign) Set(f func(p []byte) (r *core.Signature, r1 error)) *CryptographyServiceMock {
	m.mock.SignFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Sign implements github.com/insolar/insolar/core.CryptographyService interface
func (m *CryptographyServiceMock) Sign(p []byte) (r *core.Signature, r1 error) {
	atomic.AddUint64(&m.SignPreCounter, 1)
	defer atomic.AddUint64(&m.SignCounter, 1)

	if m.SignMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.SignMock.mockExpectations, CryptographyServiceMockSignParams{p},
			"CryptographyService.Sign got unexpected parameters")

		if m.SignFunc == nil {

			m.t.Fatal("No results are set for the CryptographyServiceMock.Sign")

			return
		}
	}

	if m.SignFunc == nil {
		m.t.Fatal("Unexpected call to CryptographyServiceMock.Sign")
		return
	}

	return m.SignFunc(p)
}

//SignMinimockCounter returns a count of CryptographyServiceMock.SignFunc invocations
func (m *CryptographyServiceMock) SignMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignCounter)
}

//SignMinimockPreCounter returns the value of CryptographyServiceMock.Sign invocations
func (m *CryptographyServiceMock) SignMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignPreCounter)
}

type mCryptographyServiceMockVerify struct {
	mock             *CryptographyServiceMock
	mockExpectations *CryptographyServiceMockVerifyParams
}

//CryptographyServiceMockVerifyParams represents input parameters of the CryptographyService.Verify
type CryptographyServiceMockVerifyParams struct {
	p  crypto.PublicKey
	p1 core.Signature
	p2 []byte
}

//Expect sets up expected params for the CryptographyService.Verify
func (m *mCryptographyServiceMockVerify) Expect(p crypto.PublicKey, p1 core.Signature, p2 []byte) *mCryptographyServiceMockVerify {
	m.mockExpectations = &CryptographyServiceMockVerifyParams{p, p1, p2}
	return m
}

//Return sets up a mock for CryptographyService.Verify to return Return's arguments
func (m *mCryptographyServiceMockVerify) Return(r bool) *CryptographyServiceMock {
	m.mock.VerifyFunc = func(p crypto.PublicKey, p1 core.Signature, p2 []byte) bool {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of CryptographyService.Verify method
func (m *mCryptographyServiceMockVerify) Set(f func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool)) *CryptographyServiceMock {
	m.mock.VerifyFunc = f
	m.mockExpectations = nil
	return m.mock
}

//Verify implements github.com/insolar/insolar/core.CryptographyService interface
func (m *CryptographyServiceMock) Verify(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool) {
	atomic.AddUint64(&m.VerifyPreCounter, 1)
	defer atomic.AddUint64(&m.VerifyCounter, 1)

	if m.VerifyMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.VerifyMock.mockExpectations, CryptographyServiceMockVerifyParams{p, p1, p2},
			"CryptographyService.Verify got unexpected parameters")

		if m.VerifyFunc == nil {

			m.t.Fatal("No results are set for the CryptographyServiceMock.Verify")

			return
		}
	}

	if m.VerifyFunc == nil {
		m.t.Fatal("Unexpected call to CryptographyServiceMock.Verify")
		return
	}

	return m.VerifyFunc(p, p1, p2)
}

//VerifyMinimockCounter returns a count of CryptographyServiceMock.VerifyFunc invocations
func (m *CryptographyServiceMock) VerifyMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyCounter)
}

//VerifyMinimockPreCounter returns the value of CryptographyServiceMock.Verify invocations
func (m *CryptographyServiceMock) VerifyMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.VerifyPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CryptographyServiceMock) ValidateCallCounters() {

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to CryptographyServiceMock.GetPublicKey")
	}

	if m.SignFunc != nil && atomic.LoadUint64(&m.SignCounter) == 0 {
		m.t.Fatal("Expected call to CryptographyServiceMock.Sign")
	}

	if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
		m.t.Fatal("Expected call to CryptographyServiceMock.Verify")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *CryptographyServiceMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *CryptographyServiceMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *CryptographyServiceMock) MinimockFinish() {

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		m.t.Fatal("Expected call to CryptographyServiceMock.GetPublicKey")
	}

	if m.SignFunc != nil && atomic.LoadUint64(&m.SignCounter) == 0 {
		m.t.Fatal("Expected call to CryptographyServiceMock.Sign")
	}

	if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
		m.t.Fatal("Expected call to CryptographyServiceMock.Verify")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *CryptographyServiceMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *CryptographyServiceMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.GetPublicKeyFunc == nil || atomic.LoadUint64(&m.GetPublicKeyCounter) > 0)
		ok = ok && (m.SignFunc == nil || atomic.LoadUint64(&m.SignCounter) > 0)
		ok = ok && (m.VerifyFunc == nil || atomic.LoadUint64(&m.VerifyCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
				m.t.Error("Expected call to CryptographyServiceMock.GetPublicKey")
			}

			if m.SignFunc != nil && atomic.LoadUint64(&m.SignCounter) == 0 {
				m.t.Error("Expected call to CryptographyServiceMock.Sign")
			}

			if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
				m.t.Error("Expected call to CryptographyServiceMock.Verify")
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
func (m *CryptographyServiceMock) AllMocksCalled() bool {

	if m.GetPublicKeyFunc != nil && atomic.LoadUint64(&m.GetPublicKeyCounter) == 0 {
		return false
	}

	if m.SignFunc != nil && atomic.LoadUint64(&m.SignCounter) == 0 {
		return false
	}

	if m.VerifyFunc != nil && atomic.LoadUint64(&m.VerifyCounter) == 0 {
		return false
	}

	return true
}
