package testutils

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock"
	mm_insolar "github.com/insolar/insolar/insolar"
)

// CertificateGetterMock implements insolar.CertificateGetter
type CertificateGetterMock struct {
	t minimock.Tester

	funcGetCert          func(ctx context.Context, rp1 *mm_insolar.Reference) (c2 mm_insolar.Certificate, err error)
	inspectFuncGetCert   func(ctx context.Context, rp1 *mm_insolar.Reference)
	afterGetCertCounter  uint64
	beforeGetCertCounter uint64
	GetCertMock          mCertificateGetterMockGetCert
}

// NewCertificateGetterMock returns a mock for insolar.CertificateGetter
func NewCertificateGetterMock(t minimock.Tester) *CertificateGetterMock {
	m := &CertificateGetterMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCertMock = mCertificateGetterMockGetCert{mock: m}
	m.GetCertMock.callArgs = []*CertificateGetterMockGetCertParams{}

	return m
}

type mCertificateGetterMockGetCert struct {
	mock               *CertificateGetterMock
	defaultExpectation *CertificateGetterMockGetCertExpectation
	expectations       []*CertificateGetterMockGetCertExpectation

	callArgs []*CertificateGetterMockGetCertParams
	mutex    sync.RWMutex
}

// CertificateGetterMockGetCertExpectation specifies expectation struct of the CertificateGetter.GetCert
type CertificateGetterMockGetCertExpectation struct {
	mock    *CertificateGetterMock
	params  *CertificateGetterMockGetCertParams
	results *CertificateGetterMockGetCertResults
	Counter uint64
}

// CertificateGetterMockGetCertParams contains parameters of the CertificateGetter.GetCert
type CertificateGetterMockGetCertParams struct {
	ctx context.Context
	rp1 *mm_insolar.Reference
}

// CertificateGetterMockGetCertResults contains results of the CertificateGetter.GetCert
type CertificateGetterMockGetCertResults struct {
	c2  mm_insolar.Certificate
	err error
}

// Expect sets up expected params for CertificateGetter.GetCert
func (mmGetCert *mCertificateGetterMockGetCert) Expect(ctx context.Context, rp1 *mm_insolar.Reference) *mCertificateGetterMockGetCert {
	if mmGetCert.mock.funcGetCert != nil {
		mmGetCert.mock.t.Fatalf("CertificateGetterMock.GetCert mock is already set by Set")
	}

	if mmGetCert.defaultExpectation == nil {
		mmGetCert.defaultExpectation = &CertificateGetterMockGetCertExpectation{}
	}

	mmGetCert.defaultExpectation.params = &CertificateGetterMockGetCertParams{ctx, rp1}
	for _, e := range mmGetCert.expectations {
		if minimock.Equal(e.params, mmGetCert.defaultExpectation.params) {
			mmGetCert.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetCert.defaultExpectation.params)
		}
	}

	return mmGetCert
}

// Inspect accepts an inspector function that has same arguments as the CertificateGetter.GetCert
func (mmGetCert *mCertificateGetterMockGetCert) Inspect(f func(ctx context.Context, rp1 *mm_insolar.Reference)) *mCertificateGetterMockGetCert {
	if mmGetCert.mock.inspectFuncGetCert != nil {
		mmGetCert.mock.t.Fatalf("Inspect function is already set for CertificateGetterMock.GetCert")
	}

	mmGetCert.mock.inspectFuncGetCert = f

	return mmGetCert
}

// Return sets up results that will be returned by CertificateGetter.GetCert
func (mmGetCert *mCertificateGetterMockGetCert) Return(c2 mm_insolar.Certificate, err error) *CertificateGetterMock {
	if mmGetCert.mock.funcGetCert != nil {
		mmGetCert.mock.t.Fatalf("CertificateGetterMock.GetCert mock is already set by Set")
	}

	if mmGetCert.defaultExpectation == nil {
		mmGetCert.defaultExpectation = &CertificateGetterMockGetCertExpectation{mock: mmGetCert.mock}
	}
	mmGetCert.defaultExpectation.results = &CertificateGetterMockGetCertResults{c2, err}
	return mmGetCert.mock
}

//Set uses given function f to mock the CertificateGetter.GetCert method
func (mmGetCert *mCertificateGetterMockGetCert) Set(f func(ctx context.Context, rp1 *mm_insolar.Reference) (c2 mm_insolar.Certificate, err error)) *CertificateGetterMock {
	if mmGetCert.defaultExpectation != nil {
		mmGetCert.mock.t.Fatalf("Default expectation is already set for the CertificateGetter.GetCert method")
	}

	if len(mmGetCert.expectations) > 0 {
		mmGetCert.mock.t.Fatalf("Some expectations are already set for the CertificateGetter.GetCert method")
	}

	mmGetCert.mock.funcGetCert = f
	return mmGetCert.mock
}

// When sets expectation for the CertificateGetter.GetCert which will trigger the result defined by the following
// Then helper
func (mmGetCert *mCertificateGetterMockGetCert) When(ctx context.Context, rp1 *mm_insolar.Reference) *CertificateGetterMockGetCertExpectation {
	if mmGetCert.mock.funcGetCert != nil {
		mmGetCert.mock.t.Fatalf("CertificateGetterMock.GetCert mock is already set by Set")
	}

	expectation := &CertificateGetterMockGetCertExpectation{
		mock:   mmGetCert.mock,
		params: &CertificateGetterMockGetCertParams{ctx, rp1},
	}
	mmGetCert.expectations = append(mmGetCert.expectations, expectation)
	return expectation
}

// Then sets up CertificateGetter.GetCert return parameters for the expectation previously defined by the When method
func (e *CertificateGetterMockGetCertExpectation) Then(c2 mm_insolar.Certificate, err error) *CertificateGetterMock {
	e.results = &CertificateGetterMockGetCertResults{c2, err}
	return e.mock
}

// GetCert implements insolar.CertificateGetter
func (mmGetCert *CertificateGetterMock) GetCert(ctx context.Context, rp1 *mm_insolar.Reference) (c2 mm_insolar.Certificate, err error) {
	mm_atomic.AddUint64(&mmGetCert.beforeGetCertCounter, 1)
	defer mm_atomic.AddUint64(&mmGetCert.afterGetCertCounter, 1)

	if mmGetCert.inspectFuncGetCert != nil {
		mmGetCert.inspectFuncGetCert(ctx, rp1)
	}

	params := &CertificateGetterMockGetCertParams{ctx, rp1}

	// Record call args
	mmGetCert.GetCertMock.mutex.Lock()
	mmGetCert.GetCertMock.callArgs = append(mmGetCert.GetCertMock.callArgs, params)
	mmGetCert.GetCertMock.mutex.Unlock()

	for _, e := range mmGetCert.GetCertMock.expectations {
		if minimock.Equal(e.params, params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.c2, e.results.err
		}
	}

	if mmGetCert.GetCertMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetCert.GetCertMock.defaultExpectation.Counter, 1)
		want := mmGetCert.GetCertMock.defaultExpectation.params
		got := CertificateGetterMockGetCertParams{ctx, rp1}
		if want != nil && !minimock.Equal(*want, got) {
			mmGetCert.t.Errorf("CertificateGetterMock.GetCert got unexpected parameters, want: %#v, got: %#v%s\n", *want, got, minimock.Diff(*want, got))
		}

		results := mmGetCert.GetCertMock.defaultExpectation.results
		if results == nil {
			mmGetCert.t.Fatal("No results are set for the CertificateGetterMock.GetCert")
		}
		return (*results).c2, (*results).err
	}
	if mmGetCert.funcGetCert != nil {
		return mmGetCert.funcGetCert(ctx, rp1)
	}
	mmGetCert.t.Fatalf("Unexpected call to CertificateGetterMock.GetCert. %v %v", ctx, rp1)
	return
}

// GetCertAfterCounter returns a count of finished CertificateGetterMock.GetCert invocations
func (mmGetCert *CertificateGetterMock) GetCertAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCert.afterGetCertCounter)
}

// GetCertBeforeCounter returns a count of CertificateGetterMock.GetCert invocations
func (mmGetCert *CertificateGetterMock) GetCertBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCert.beforeGetCertCounter)
}

// Calls returns a list of arguments used in each call to CertificateGetterMock.GetCert.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetCert *mCertificateGetterMockGetCert) Calls() []*CertificateGetterMockGetCertParams {
	mmGetCert.mutex.RLock()

	argCopy := make([]*CertificateGetterMockGetCertParams, len(mmGetCert.callArgs))
	copy(argCopy, mmGetCert.callArgs)

	mmGetCert.mutex.RUnlock()

	return argCopy
}

// MinimockGetCertDone returns true if the count of the GetCert invocations corresponds
// the number of defined expectations
func (m *CertificateGetterMock) MinimockGetCertDone() bool {
	for _, e := range m.GetCertMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCertMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCertCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCert != nil && mm_atomic.LoadUint64(&m.afterGetCertCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetCertInspect logs each unmet expectation
func (m *CertificateGetterMock) MinimockGetCertInspect() {
	for _, e := range m.GetCertMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CertificateGetterMock.GetCert with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCertMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCertCounter) < 1 {
		if m.GetCertMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CertificateGetterMock.GetCert")
		} else {
			m.t.Errorf("Expected call to CertificateGetterMock.GetCert with params: %#v", *m.GetCertMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCert != nil && mm_atomic.LoadUint64(&m.afterGetCertCounter) < 1 {
		m.t.Error("Expected call to CertificateGetterMock.GetCert")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CertificateGetterMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetCertInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CertificateGetterMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *CertificateGetterMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetCertDone()
}
