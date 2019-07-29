package cryptkit

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "DataSigner" can be found in github.com/insolar/insolar/network/consensus/common/cryptkit
*/
import (
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"

	testify_assert "github.com/stretchr/testify/assert"
)

//DataSignerMock implements github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner
type DataSignerMock struct {
	t minimock.Tester

	GetDigestMethodFunc       func() (r DigestMethod)
	GetDigestMethodCounter    uint64
	GetDigestMethodPreCounter uint64
	GetDigestMethodMock       mDataSignerMockGetDigestMethod

	GetDigestOfFunc       func(p io.Reader) (r Digest)
	GetDigestOfCounter    uint64
	GetDigestOfPreCounter uint64
	GetDigestOfMock       mDataSignerMockGetDigestOf

	GetSignMethodFunc       func() (r SignMethod)
	GetSignMethodCounter    uint64
	GetSignMethodPreCounter uint64
	GetSignMethodMock       mDataSignerMockGetSignMethod

	GetSignatureMethodFunc       func() (r SignatureMethod)
	GetSignatureMethodCounter    uint64
	GetSignatureMethodPreCounter uint64
	GetSignatureMethodMock       mDataSignerMockGetSignatureMethod

	SignDataFunc       func(p io.Reader) (r SignedDigest)
	SignDataCounter    uint64
	SignDataPreCounter uint64
	SignDataMock       mDataSignerMockSignData

	SignDigestFunc       func(p Digest) (r Signature)
	SignDigestCounter    uint64
	SignDigestPreCounter uint64
	SignDigestMock       mDataSignerMockSignDigest
}

//NewDataSignerMock returns a mock for github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner
func NewDataSignerMock(t minimock.Tester) *DataSignerMock {
	m := &DataSignerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetDigestMethodMock = mDataSignerMockGetDigestMethod{mock: m}
	m.GetDigestOfMock = mDataSignerMockGetDigestOf{mock: m}
	m.GetSignMethodMock = mDataSignerMockGetSignMethod{mock: m}
	m.GetSignatureMethodMock = mDataSignerMockGetSignatureMethod{mock: m}
	m.SignDataMock = mDataSignerMockSignData{mock: m}
	m.SignDigestMock = mDataSignerMockSignDigest{mock: m}

	return m
}

type mDataSignerMockGetDigestMethod struct {
	mock              *DataSignerMock
	mainExpectation   *DataSignerMockGetDigestMethodExpectation
	expectationSeries []*DataSignerMockGetDigestMethodExpectation
}

type DataSignerMockGetDigestMethodExpectation struct {
	result *DataSignerMockGetDigestMethodResult
}

type DataSignerMockGetDigestMethodResult struct {
	r DigestMethod
}

//Expect specifies that invocation of DataSigner.GetDigestMethod is expected from 1 to Infinity times
func (m *mDataSignerMockGetDigestMethod) Expect() *mDataSignerMockGetDigestMethod {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetDigestMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of DataSigner.GetDigestMethod
func (m *mDataSignerMockGetDigestMethod) Return(r DigestMethod) *DataSignerMock {
	m.mock.GetDigestMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetDigestMethodExpectation{}
	}
	m.mainExpectation.result = &DataSignerMockGetDigestMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DataSigner.GetDigestMethod is expected once
func (m *mDataSignerMockGetDigestMethod) ExpectOnce() *DataSignerMockGetDigestMethodExpectation {
	m.mock.GetDigestMethodFunc = nil
	m.mainExpectation = nil

	expectation := &DataSignerMockGetDigestMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DataSignerMockGetDigestMethodExpectation) Return(r DigestMethod) {
	e.result = &DataSignerMockGetDigestMethodResult{r}
}

//Set uses given function f as a mock of DataSigner.GetDigestMethod method
func (m *mDataSignerMockGetDigestMethod) Set(f func() (r DigestMethod)) *DataSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDigestMethodFunc = f
	return m.mock
}

//GetDigestMethod implements github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner interface
func (m *DataSignerMock) GetDigestMethod() (r DigestMethod) {
	counter := atomic.AddUint64(&m.GetDigestMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetDigestMethodCounter, 1)

	if len(m.GetDigestMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDigestMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DataSignerMock.GetDigestMethod.")
			return
		}

		result := m.GetDigestMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetDigestMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetDigestMethodMock.mainExpectation != nil {

		result := m.GetDigestMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetDigestMethod")
		}

		r = result.r

		return
	}

	if m.GetDigestMethodFunc == nil {
		m.t.Fatalf("Unexpected call to DataSignerMock.GetDigestMethod.")
		return
	}

	return m.GetDigestMethodFunc()
}

//GetDigestMethodMinimockCounter returns a count of DataSignerMock.GetDigestMethodFunc invocations
func (m *DataSignerMock) GetDigestMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodCounter)
}

//GetDigestMethodMinimockPreCounter returns the value of DataSignerMock.GetDigestMethod invocations
func (m *DataSignerMock) GetDigestMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestMethodPreCounter)
}

//GetDigestMethodFinished returns true if mock invocations count is ok
func (m *DataSignerMock) GetDigestMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDigestMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDigestMethodCounter) == uint64(len(m.GetDigestMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDigestMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDigestMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDigestMethodFunc != nil {
		return atomic.LoadUint64(&m.GetDigestMethodCounter) > 0
	}

	return true
}

type mDataSignerMockGetDigestOf struct {
	mock              *DataSignerMock
	mainExpectation   *DataSignerMockGetDigestOfExpectation
	expectationSeries []*DataSignerMockGetDigestOfExpectation
}

type DataSignerMockGetDigestOfExpectation struct {
	input  *DataSignerMockGetDigestOfInput
	result *DataSignerMockGetDigestOfResult
}

type DataSignerMockGetDigestOfInput struct {
	p io.Reader
}

type DataSignerMockGetDigestOfResult struct {
	r Digest
}

//Expect specifies that invocation of DataSigner.GetDigestOf is expected from 1 to Infinity times
func (m *mDataSignerMockGetDigestOf) Expect(p io.Reader) *mDataSignerMockGetDigestOf {
	m.mock.GetDigestOfFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetDigestOfExpectation{}
	}
	m.mainExpectation.input = &DataSignerMockGetDigestOfInput{p}
	return m
}

//Return specifies results of invocation of DataSigner.GetDigestOf
func (m *mDataSignerMockGetDigestOf) Return(r Digest) *DataSignerMock {
	m.mock.GetDigestOfFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetDigestOfExpectation{}
	}
	m.mainExpectation.result = &DataSignerMockGetDigestOfResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DataSigner.GetDigestOf is expected once
func (m *mDataSignerMockGetDigestOf) ExpectOnce(p io.Reader) *DataSignerMockGetDigestOfExpectation {
	m.mock.GetDigestOfFunc = nil
	m.mainExpectation = nil

	expectation := &DataSignerMockGetDigestOfExpectation{}
	expectation.input = &DataSignerMockGetDigestOfInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DataSignerMockGetDigestOfExpectation) Return(r Digest) {
	e.result = &DataSignerMockGetDigestOfResult{r}
}

//Set uses given function f as a mock of DataSigner.GetDigestOf method
func (m *mDataSignerMockGetDigestOf) Set(f func(p io.Reader) (r Digest)) *DataSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetDigestOfFunc = f
	return m.mock
}

//GetDigestOf implements github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner interface
func (m *DataSignerMock) GetDigestOf(p io.Reader) (r Digest) {
	counter := atomic.AddUint64(&m.GetDigestOfPreCounter, 1)
	defer atomic.AddUint64(&m.GetDigestOfCounter, 1)

	if len(m.GetDigestOfMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetDigestOfMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DataSignerMock.GetDigestOf. %v", p)
			return
		}

		input := m.GetDigestOfMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DataSignerMockGetDigestOfInput{p}, "DataSigner.GetDigestOf got unexpected parameters")

		result := m.GetDigestOfMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetDigestOf")
			return
		}

		r = result.r

		return
	}

	if m.GetDigestOfMock.mainExpectation != nil {

		input := m.GetDigestOfMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DataSignerMockGetDigestOfInput{p}, "DataSigner.GetDigestOf got unexpected parameters")
		}

		result := m.GetDigestOfMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetDigestOf")
		}

		r = result.r

		return
	}

	if m.GetDigestOfFunc == nil {
		m.t.Fatalf("Unexpected call to DataSignerMock.GetDigestOf. %v", p)
		return
	}

	return m.GetDigestOfFunc(p)
}

//GetDigestOfMinimockCounter returns a count of DataSignerMock.GetDigestOfFunc invocations
func (m *DataSignerMock) GetDigestOfMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestOfCounter)
}

//GetDigestOfMinimockPreCounter returns the value of DataSignerMock.GetDigestOf invocations
func (m *DataSignerMock) GetDigestOfMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetDigestOfPreCounter)
}

//GetDigestOfFinished returns true if mock invocations count is ok
func (m *DataSignerMock) GetDigestOfFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetDigestOfMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetDigestOfCounter) == uint64(len(m.GetDigestOfMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetDigestOfMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetDigestOfCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetDigestOfFunc != nil {
		return atomic.LoadUint64(&m.GetDigestOfCounter) > 0
	}

	return true
}

type mDataSignerMockGetSignMethod struct {
	mock              *DataSignerMock
	mainExpectation   *DataSignerMockGetSignMethodExpectation
	expectationSeries []*DataSignerMockGetSignMethodExpectation
}

type DataSignerMockGetSignMethodExpectation struct {
	result *DataSignerMockGetSignMethodResult
}

type DataSignerMockGetSignMethodResult struct {
	r SignMethod
}

//Expect specifies that invocation of DataSigner.GetSignMethod is expected from 1 to Infinity times
func (m *mDataSignerMockGetSignMethod) Expect() *mDataSignerMockGetSignMethod {
	m.mock.GetSignMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetSignMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of DataSigner.GetSignMethod
func (m *mDataSignerMockGetSignMethod) Return(r SignMethod) *DataSignerMock {
	m.mock.GetSignMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetSignMethodExpectation{}
	}
	m.mainExpectation.result = &DataSignerMockGetSignMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DataSigner.GetSignMethod is expected once
func (m *mDataSignerMockGetSignMethod) ExpectOnce() *DataSignerMockGetSignMethodExpectation {
	m.mock.GetSignMethodFunc = nil
	m.mainExpectation = nil

	expectation := &DataSignerMockGetSignMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DataSignerMockGetSignMethodExpectation) Return(r SignMethod) {
	e.result = &DataSignerMockGetSignMethodResult{r}
}

//Set uses given function f as a mock of DataSigner.GetSignMethod method
func (m *mDataSignerMockGetSignMethod) Set(f func() (r SignMethod)) *DataSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignMethodFunc = f
	return m.mock
}

//GetSignMethod implements github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner interface
func (m *DataSignerMock) GetSignMethod() (r SignMethod) {
	counter := atomic.AddUint64(&m.GetSignMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignMethodCounter, 1)

	if len(m.GetSignMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DataSignerMock.GetSignMethod.")
			return
		}

		result := m.GetSignMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetSignMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignMethodMock.mainExpectation != nil {

		result := m.GetSignMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetSignMethod")
		}

		r = result.r

		return
	}

	if m.GetSignMethodFunc == nil {
		m.t.Fatalf("Unexpected call to DataSignerMock.GetSignMethod.")
		return
	}

	return m.GetSignMethodFunc()
}

//GetSignMethodMinimockCounter returns a count of DataSignerMock.GetSignMethodFunc invocations
func (m *DataSignerMock) GetSignMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignMethodCounter)
}

//GetSignMethodMinimockPreCounter returns the value of DataSignerMock.GetSignMethod invocations
func (m *DataSignerMock) GetSignMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignMethodPreCounter)
}

//GetSignMethodFinished returns true if mock invocations count is ok
func (m *DataSignerMock) GetSignMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignMethodCounter) == uint64(len(m.GetSignMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignMethodFunc != nil {
		return atomic.LoadUint64(&m.GetSignMethodCounter) > 0
	}

	return true
}

type mDataSignerMockGetSignatureMethod struct {
	mock              *DataSignerMock
	mainExpectation   *DataSignerMockGetSignatureMethodExpectation
	expectationSeries []*DataSignerMockGetSignatureMethodExpectation
}

type DataSignerMockGetSignatureMethodExpectation struct {
	result *DataSignerMockGetSignatureMethodResult
}

type DataSignerMockGetSignatureMethodResult struct {
	r SignatureMethod
}

//Expect specifies that invocation of DataSigner.GetSignatureMethod is expected from 1 to Infinity times
func (m *mDataSignerMockGetSignatureMethod) Expect() *mDataSignerMockGetSignatureMethod {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetSignatureMethodExpectation{}
	}

	return m
}

//Return specifies results of invocation of DataSigner.GetSignatureMethod
func (m *mDataSignerMockGetSignatureMethod) Return(r SignatureMethod) *DataSignerMock {
	m.mock.GetSignatureMethodFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockGetSignatureMethodExpectation{}
	}
	m.mainExpectation.result = &DataSignerMockGetSignatureMethodResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DataSigner.GetSignatureMethod is expected once
func (m *mDataSignerMockGetSignatureMethod) ExpectOnce() *DataSignerMockGetSignatureMethodExpectation {
	m.mock.GetSignatureMethodFunc = nil
	m.mainExpectation = nil

	expectation := &DataSignerMockGetSignatureMethodExpectation{}

	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DataSignerMockGetSignatureMethodExpectation) Return(r SignatureMethod) {
	e.result = &DataSignerMockGetSignatureMethodResult{r}
}

//Set uses given function f as a mock of DataSigner.GetSignatureMethod method
func (m *mDataSignerMockGetSignatureMethod) Set(f func() (r SignatureMethod)) *DataSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.GetSignatureMethodFunc = f
	return m.mock
}

//GetSignatureMethod implements github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner interface
func (m *DataSignerMock) GetSignatureMethod() (r SignatureMethod) {
	counter := atomic.AddUint64(&m.GetSignatureMethodPreCounter, 1)
	defer atomic.AddUint64(&m.GetSignatureMethodCounter, 1)

	if len(m.GetSignatureMethodMock.expectationSeries) > 0 {
		if counter > uint64(len(m.GetSignatureMethodMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DataSignerMock.GetSignatureMethod.")
			return
		}

		result := m.GetSignatureMethodMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetSignatureMethod")
			return
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodMock.mainExpectation != nil {

		result := m.GetSignatureMethodMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.GetSignatureMethod")
		}

		r = result.r

		return
	}

	if m.GetSignatureMethodFunc == nil {
		m.t.Fatalf("Unexpected call to DataSignerMock.GetSignatureMethod.")
		return
	}

	return m.GetSignatureMethodFunc()
}

//GetSignatureMethodMinimockCounter returns a count of DataSignerMock.GetSignatureMethodFunc invocations
func (m *DataSignerMock) GetSignatureMethodMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodCounter)
}

//GetSignatureMethodMinimockPreCounter returns the value of DataSignerMock.GetSignatureMethod invocations
func (m *DataSignerMock) GetSignatureMethodMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.GetSignatureMethodPreCounter)
}

//GetSignatureMethodFinished returns true if mock invocations count is ok
func (m *DataSignerMock) GetSignatureMethodFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.GetSignatureMethodMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.GetSignatureMethodCounter) == uint64(len(m.GetSignatureMethodMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.GetSignatureMethodMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.GetSignatureMethodCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.GetSignatureMethodFunc != nil {
		return atomic.LoadUint64(&m.GetSignatureMethodCounter) > 0
	}

	return true
}

type mDataSignerMockSignData struct {
	mock              *DataSignerMock
	mainExpectation   *DataSignerMockSignDataExpectation
	expectationSeries []*DataSignerMockSignDataExpectation
}

type DataSignerMockSignDataExpectation struct {
	input  *DataSignerMockSignDataInput
	result *DataSignerMockSignDataResult
}

type DataSignerMockSignDataInput struct {
	p io.Reader
}

type DataSignerMockSignDataResult struct {
	r SignedDigest
}

//Expect specifies that invocation of DataSigner.SignData is expected from 1 to Infinity times
func (m *mDataSignerMockSignData) Expect(p io.Reader) *mDataSignerMockSignData {
	m.mock.SignDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockSignDataExpectation{}
	}
	m.mainExpectation.input = &DataSignerMockSignDataInput{p}
	return m
}

//Return specifies results of invocation of DataSigner.SignData
func (m *mDataSignerMockSignData) Return(r SignedDigest) *DataSignerMock {
	m.mock.SignDataFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockSignDataExpectation{}
	}
	m.mainExpectation.result = &DataSignerMockSignDataResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DataSigner.SignData is expected once
func (m *mDataSignerMockSignData) ExpectOnce(p io.Reader) *DataSignerMockSignDataExpectation {
	m.mock.SignDataFunc = nil
	m.mainExpectation = nil

	expectation := &DataSignerMockSignDataExpectation{}
	expectation.input = &DataSignerMockSignDataInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DataSignerMockSignDataExpectation) Return(r SignedDigest) {
	e.result = &DataSignerMockSignDataResult{r}
}

//Set uses given function f as a mock of DataSigner.SignData method
func (m *mDataSignerMockSignData) Set(f func(p io.Reader) (r SignedDigest)) *DataSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignDataFunc = f
	return m.mock
}

//SignData implements github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner interface
func (m *DataSignerMock) SignData(p io.Reader) (r SignedDigest) {
	counter := atomic.AddUint64(&m.SignDataPreCounter, 1)
	defer atomic.AddUint64(&m.SignDataCounter, 1)

	if len(m.SignDataMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignDataMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DataSignerMock.SignData. %v", p)
			return
		}

		input := m.SignDataMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DataSignerMockSignDataInput{p}, "DataSigner.SignData got unexpected parameters")

		result := m.SignDataMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.SignData")
			return
		}

		r = result.r

		return
	}

	if m.SignDataMock.mainExpectation != nil {

		input := m.SignDataMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DataSignerMockSignDataInput{p}, "DataSigner.SignData got unexpected parameters")
		}

		result := m.SignDataMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.SignData")
		}

		r = result.r

		return
	}

	if m.SignDataFunc == nil {
		m.t.Fatalf("Unexpected call to DataSignerMock.SignData. %v", p)
		return
	}

	return m.SignDataFunc(p)
}

//SignDataMinimockCounter returns a count of DataSignerMock.SignDataFunc invocations
func (m *DataSignerMock) SignDataMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignDataCounter)
}

//SignDataMinimockPreCounter returns the value of DataSignerMock.SignData invocations
func (m *DataSignerMock) SignDataMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignDataPreCounter)
}

//SignDataFinished returns true if mock invocations count is ok
func (m *DataSignerMock) SignDataFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SignDataMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SignDataCounter) == uint64(len(m.SignDataMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SignDataMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SignDataCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SignDataFunc != nil {
		return atomic.LoadUint64(&m.SignDataCounter) > 0
	}

	return true
}

type mDataSignerMockSignDigest struct {
	mock              *DataSignerMock
	mainExpectation   *DataSignerMockSignDigestExpectation
	expectationSeries []*DataSignerMockSignDigestExpectation
}

type DataSignerMockSignDigestExpectation struct {
	input  *DataSignerMockSignDigestInput
	result *DataSignerMockSignDigestResult
}

type DataSignerMockSignDigestInput struct {
	p Digest
}

type DataSignerMockSignDigestResult struct {
	r Signature
}

//Expect specifies that invocation of DataSigner.SignDigest is expected from 1 to Infinity times
func (m *mDataSignerMockSignDigest) Expect(p Digest) *mDataSignerMockSignDigest {
	m.mock.SignDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockSignDigestExpectation{}
	}
	m.mainExpectation.input = &DataSignerMockSignDigestInput{p}
	return m
}

//Return specifies results of invocation of DataSigner.SignDigest
func (m *mDataSignerMockSignDigest) Return(r Signature) *DataSignerMock {
	m.mock.SignDigestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &DataSignerMockSignDigestExpectation{}
	}
	m.mainExpectation.result = &DataSignerMockSignDigestResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of DataSigner.SignDigest is expected once
func (m *mDataSignerMockSignDigest) ExpectOnce(p Digest) *DataSignerMockSignDigestExpectation {
	m.mock.SignDigestFunc = nil
	m.mainExpectation = nil

	expectation := &DataSignerMockSignDigestExpectation{}
	expectation.input = &DataSignerMockSignDigestInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *DataSignerMockSignDigestExpectation) Return(r Signature) {
	e.result = &DataSignerMockSignDigestResult{r}
}

//Set uses given function f as a mock of DataSigner.SignDigest method
func (m *mDataSignerMockSignDigest) Set(f func(p Digest) (r Signature)) *DataSignerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SignDigestFunc = f
	return m.mock
}

//SignDigest implements github.com/insolar/insolar/network/consensus/common/cryptkit.DataSigner interface
func (m *DataSignerMock) SignDigest(p Digest) (r Signature) {
	counter := atomic.AddUint64(&m.SignDigestPreCounter, 1)
	defer atomic.AddUint64(&m.SignDigestCounter, 1)

	if len(m.SignDigestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SignDigestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to DataSignerMock.SignDigest. %v", p)
			return
		}

		input := m.SignDigestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, DataSignerMockSignDigestInput{p}, "DataSigner.SignDigest got unexpected parameters")

		result := m.SignDigestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.SignDigest")
			return
		}

		r = result.r

		return
	}

	if m.SignDigestMock.mainExpectation != nil {

		input := m.SignDigestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, DataSignerMockSignDigestInput{p}, "DataSigner.SignDigest got unexpected parameters")
		}

		result := m.SignDigestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the DataSignerMock.SignDigest")
		}

		r = result.r

		return
	}

	if m.SignDigestFunc == nil {
		m.t.Fatalf("Unexpected call to DataSignerMock.SignDigest. %v", p)
		return
	}

	return m.SignDigestFunc(p)
}

//SignDigestMinimockCounter returns a count of DataSignerMock.SignDigestFunc invocations
func (m *DataSignerMock) SignDigestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SignDigestCounter)
}

//SignDigestMinimockPreCounter returns the value of DataSignerMock.SignDigest invocations
func (m *DataSignerMock) SignDigestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SignDigestPreCounter)
}

//SignDigestFinished returns true if mock invocations count is ok
func (m *DataSignerMock) SignDigestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SignDigestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SignDigestCounter) == uint64(len(m.SignDigestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SignDigestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SignDigestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SignDigestFunc != nil {
		return atomic.LoadUint64(&m.SignDigestCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DataSignerMock) ValidateCallCounters() {

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetDigestMethod")
	}

	if !m.GetDigestOfFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetDigestOf")
	}

	if !m.GetSignMethodFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetSignMethod")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetSignatureMethod")
	}

	if !m.SignDataFinished() {
		m.t.Fatal("Expected call to DataSignerMock.SignData")
	}

	if !m.SignDigestFinished() {
		m.t.Fatal("Expected call to DataSignerMock.SignDigest")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *DataSignerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *DataSignerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *DataSignerMock) MinimockFinish() {

	if !m.GetDigestMethodFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetDigestMethod")
	}

	if !m.GetDigestOfFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetDigestOf")
	}

	if !m.GetSignMethodFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetSignMethod")
	}

	if !m.GetSignatureMethodFinished() {
		m.t.Fatal("Expected call to DataSignerMock.GetSignatureMethod")
	}

	if !m.SignDataFinished() {
		m.t.Fatal("Expected call to DataSignerMock.SignData")
	}

	if !m.SignDigestFinished() {
		m.t.Fatal("Expected call to DataSignerMock.SignDigest")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *DataSignerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *DataSignerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.GetDigestMethodFinished()
		ok = ok && m.GetDigestOfFinished()
		ok = ok && m.GetSignMethodFinished()
		ok = ok && m.GetSignatureMethodFinished()
		ok = ok && m.SignDataFinished()
		ok = ok && m.SignDigestFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.GetDigestMethodFinished() {
				m.t.Error("Expected call to DataSignerMock.GetDigestMethod")
			}

			if !m.GetDigestOfFinished() {
				m.t.Error("Expected call to DataSignerMock.GetDigestOf")
			}

			if !m.GetSignMethodFinished() {
				m.t.Error("Expected call to DataSignerMock.GetSignMethod")
			}

			if !m.GetSignatureMethodFinished() {
				m.t.Error("Expected call to DataSignerMock.GetSignatureMethod")
			}

			if !m.SignDataFinished() {
				m.t.Error("Expected call to DataSignerMock.SignData")
			}

			if !m.SignDigestFinished() {
				m.t.Error("Expected call to DataSignerMock.SignDigest")
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
func (m *DataSignerMock) AllMocksCalled() bool {

	if !m.GetDigestMethodFinished() {
		return false
	}

	if !m.GetDigestOfFinished() {
		return false
	}

	if !m.GetSignMethodFinished() {
		return false
	}

	if !m.GetSignatureMethodFinished() {
		return false
	}

	if !m.SignDataFinished() {
		return false
	}

	if !m.SignDigestFinished() {
		return false
	}

	return true
}
