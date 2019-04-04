package artifact

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Manager" can be found in github.com/insolar/insolar/internal/ledger/artifact
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	insolar "github.com/insolar/insolar/insolar"

	testify_assert "github.com/stretchr/testify/assert"
)

//ManagerMock implements github.com/insolar/insolar/internal/ledger/artifact.Manager
type ManagerMock struct {
	t minimock.Tester

	ActivateObjectFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) (r ObjectDescriptor, r1 error)
	ActivateObjectCounter    uint64
	ActivateObjectPreCounter uint64
	ActivateObjectMock       mManagerMockActivateObject

	RegisterRequestFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error)
	RegisterRequestCounter    uint64
	RegisterRequestPreCounter uint64
	RegisterRequestMock       mManagerMockRegisterRequest

	RegisterResultFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error)
	RegisterResultCounter    uint64
	RegisterResultPreCounter uint64
	RegisterResultMock       mManagerMockRegisterResult

	UpdateObjectFunc       func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) (r ObjectDescriptor, r1 error)
	UpdateObjectCounter    uint64
	UpdateObjectPreCounter uint64
	UpdateObjectMock       mManagerMockUpdateObject
}

//NewManagerMock returns a mock for github.com/insolar/insolar/internal/ledger/artifact.Manager
func NewManagerMock(t minimock.Tester) *ManagerMock {
	m := &ManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ActivateObjectMock = mManagerMockActivateObject{mock: m}
	m.RegisterRequestMock = mManagerMockRegisterRequest{mock: m}
	m.RegisterResultMock = mManagerMockRegisterResult{mock: m}
	m.UpdateObjectMock = mManagerMockUpdateObject{mock: m}

	return m
}

type mManagerMockActivateObject struct {
	mock              *ManagerMock
	mainExpectation   *ManagerMockActivateObjectExpectation
	expectationSeries []*ManagerMockActivateObjectExpectation
}

type ManagerMockActivateObjectExpectation struct {
	input  *ManagerMockActivateObjectInput
	result *ManagerMockActivateObjectResult
}

type ManagerMockActivateObjectInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 insolar.Reference
	p4 insolar.Reference
	p5 bool
	p6 []byte
}

type ManagerMockActivateObjectResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Manager.ActivateObject is expected from 1 to Infinity times
func (m *mManagerMockActivateObject) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) *mManagerMockActivateObject {
	m.mock.ActivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockActivateObjectExpectation{}
	}
	m.mainExpectation.input = &ManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}
	return m
}

//Return specifies results of invocation of Manager.ActivateObject
func (m *mManagerMockActivateObject) Return(r ObjectDescriptor, r1 error) *ManagerMock {
	m.mock.ActivateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockActivateObjectExpectation{}
	}
	m.mainExpectation.result = &ManagerMockActivateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Manager.ActivateObject is expected once
func (m *mManagerMockActivateObject) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) *ManagerMockActivateObjectExpectation {
	m.mock.ActivateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ManagerMockActivateObjectExpectation{}
	expectation.input = &ManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ManagerMockActivateObjectExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ManagerMockActivateObjectResult{r, r1}
}

//Set uses given function f as a mock of Manager.ActivateObject method
func (m *mManagerMockActivateObject) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) (r ObjectDescriptor, r1 error)) *ManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.ActivateObjectFunc = f
	return m.mock
}

//ActivateObject implements github.com/insolar/insolar/internal/ledger/artifact.Manager interface
func (m *ManagerMock) ActivateObject(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 insolar.Reference, p4 insolar.Reference, p5 bool, p6 []byte) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.ActivateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.ActivateObjectCounter, 1)

	if len(m.ActivateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.ActivateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ManagerMock.ActivateObject. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
			return
		}

		input := m.ActivateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}, "Manager.ActivateObject got unexpected parameters")

		result := m.ActivateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.ActivateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivateObjectMock.mainExpectation != nil {

		input := m.ActivateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ManagerMockActivateObjectInput{p, p1, p2, p3, p4, p5, p6}, "Manager.ActivateObject got unexpected parameters")
		}

		result := m.ActivateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.ActivateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.ActivateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ManagerMock.ActivateObject. %v %v %v %v %v %v %v", p, p1, p2, p3, p4, p5, p6)
		return
	}

	return m.ActivateObjectFunc(p, p1, p2, p3, p4, p5, p6)
}

//ActivateObjectMinimockCounter returns a count of ManagerMock.ActivateObjectFunc invocations
func (m *ManagerMock) ActivateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ActivateObjectCounter)
}

//ActivateObjectMinimockPreCounter returns the value of ManagerMock.ActivateObject invocations
func (m *ManagerMock) ActivateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ActivateObjectPreCounter)
}

//ActivateObjectFinished returns true if mock invocations count is ok
func (m *ManagerMock) ActivateObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.ActivateObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.ActivateObjectCounter) == uint64(len(m.ActivateObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.ActivateObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.ActivateObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.ActivateObjectFunc != nil {
		return atomic.LoadUint64(&m.ActivateObjectCounter) > 0
	}

	return true
}

type mManagerMockRegisterRequest struct {
	mock              *ManagerMock
	mainExpectation   *ManagerMockRegisterRequestExpectation
	expectationSeries []*ManagerMockRegisterRequestExpectation
}

type ManagerMockRegisterRequestExpectation struct {
	input  *ManagerMockRegisterRequestInput
	result *ManagerMockRegisterRequestResult
}

type ManagerMockRegisterRequestInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Parcel
}

type ManagerMockRegisterRequestResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Manager.RegisterRequest is expected from 1 to Infinity times
func (m *mManagerMockRegisterRequest) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) *mManagerMockRegisterRequest {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockRegisterRequestExpectation{}
	}
	m.mainExpectation.input = &ManagerMockRegisterRequestInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of Manager.RegisterRequest
func (m *mManagerMockRegisterRequest) Return(r *insolar.ID, r1 error) *ManagerMock {
	m.mock.RegisterRequestFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockRegisterRequestExpectation{}
	}
	m.mainExpectation.result = &ManagerMockRegisterRequestResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Manager.RegisterRequest is expected once
func (m *mManagerMockRegisterRequest) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) *ManagerMockRegisterRequestExpectation {
	m.mock.RegisterRequestFunc = nil
	m.mainExpectation = nil

	expectation := &ManagerMockRegisterRequestExpectation{}
	expectation.input = &ManagerMockRegisterRequestInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ManagerMockRegisterRequestExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ManagerMockRegisterRequestResult{r, r1}
}

//Set uses given function f as a mock of Manager.RegisterRequest method
func (m *mManagerMockRegisterRequest) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error)) *ManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterRequestFunc = f
	return m.mock
}

//RegisterRequest implements github.com/insolar/insolar/internal/ledger/artifact.Manager interface
func (m *ManagerMock) RegisterRequest(p context.Context, p1 insolar.Reference, p2 insolar.Parcel) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterRequestPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterRequestCounter, 1)

	if len(m.RegisterRequestMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterRequestMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ManagerMock.RegisterRequest. %v %v %v", p, p1, p2)
			return
		}

		input := m.RegisterRequestMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ManagerMockRegisterRequestInput{p, p1, p2}, "Manager.RegisterRequest got unexpected parameters")

		result := m.RegisterRequestMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.RegisterRequest")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestMock.mainExpectation != nil {

		input := m.RegisterRequestMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ManagerMockRegisterRequestInput{p, p1, p2}, "Manager.RegisterRequest got unexpected parameters")
		}

		result := m.RegisterRequestMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.RegisterRequest")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterRequestFunc == nil {
		m.t.Fatalf("Unexpected call to ManagerMock.RegisterRequest. %v %v %v", p, p1, p2)
		return
	}

	return m.RegisterRequestFunc(p, p1, p2)
}

//RegisterRequestMinimockCounter returns a count of ManagerMock.RegisterRequestFunc invocations
func (m *ManagerMock) RegisterRequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestCounter)
}

//RegisterRequestMinimockPreCounter returns the value of ManagerMock.RegisterRequest invocations
func (m *ManagerMock) RegisterRequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterRequestPreCounter)
}

//RegisterRequestFinished returns true if mock invocations count is ok
func (m *ManagerMock) RegisterRequestFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterRequestMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterRequestCounter) == uint64(len(m.RegisterRequestMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterRequestMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterRequestCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterRequestFunc != nil {
		return atomic.LoadUint64(&m.RegisterRequestCounter) > 0
	}

	return true
}

type mManagerMockRegisterResult struct {
	mock              *ManagerMock
	mainExpectation   *ManagerMockRegisterResultExpectation
	expectationSeries []*ManagerMockRegisterResultExpectation
}

type ManagerMockRegisterResultExpectation struct {
	input  *ManagerMockRegisterResultInput
	result *ManagerMockRegisterResultResult
}

type ManagerMockRegisterResultInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 []byte
}

type ManagerMockRegisterResultResult struct {
	r  *insolar.ID
	r1 error
}

//Expect specifies that invocation of Manager.RegisterResult is expected from 1 to Infinity times
func (m *mManagerMockRegisterResult) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) *mManagerMockRegisterResult {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockRegisterResultExpectation{}
	}
	m.mainExpectation.input = &ManagerMockRegisterResultInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Manager.RegisterResult
func (m *mManagerMockRegisterResult) Return(r *insolar.ID, r1 error) *ManagerMock {
	m.mock.RegisterResultFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockRegisterResultExpectation{}
	}
	m.mainExpectation.result = &ManagerMockRegisterResultResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Manager.RegisterResult is expected once
func (m *mManagerMockRegisterResult) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) *ManagerMockRegisterResultExpectation {
	m.mock.RegisterResultFunc = nil
	m.mainExpectation = nil

	expectation := &ManagerMockRegisterResultExpectation{}
	expectation.input = &ManagerMockRegisterResultInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ManagerMockRegisterResultExpectation) Return(r *insolar.ID, r1 error) {
	e.result = &ManagerMockRegisterResultResult{r, r1}
}

//Set uses given function f as a mock of Manager.RegisterResult method
func (m *mManagerMockRegisterResult) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error)) *ManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterResultFunc = f
	return m.mock
}

//RegisterResult implements github.com/insolar/insolar/internal/ledger/artifact.Manager interface
func (m *ManagerMock) RegisterResult(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 []byte) (r *insolar.ID, r1 error) {
	counter := atomic.AddUint64(&m.RegisterResultPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterResultCounter, 1)

	if len(m.RegisterResultMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterResultMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ManagerMock.RegisterResult. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.RegisterResultMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ManagerMockRegisterResultInput{p, p1, p2, p3}, "Manager.RegisterResult got unexpected parameters")

		result := m.RegisterResultMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.RegisterResult")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterResultMock.mainExpectation != nil {

		input := m.RegisterResultMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ManagerMockRegisterResultInput{p, p1, p2, p3}, "Manager.RegisterResult got unexpected parameters")
		}

		result := m.RegisterResultMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.RegisterResult")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.RegisterResultFunc == nil {
		m.t.Fatalf("Unexpected call to ManagerMock.RegisterResult. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.RegisterResultFunc(p, p1, p2, p3)
}

//RegisterResultMinimockCounter returns a count of ManagerMock.RegisterResultFunc invocations
func (m *ManagerMock) RegisterResultMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultCounter)
}

//RegisterResultMinimockPreCounter returns the value of ManagerMock.RegisterResult invocations
func (m *ManagerMock) RegisterResultMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterResultPreCounter)
}

//RegisterResultFinished returns true if mock invocations count is ok
func (m *ManagerMock) RegisterResultFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterResultMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterResultCounter) == uint64(len(m.RegisterResultMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterResultMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterResultCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterResultFunc != nil {
		return atomic.LoadUint64(&m.RegisterResultCounter) > 0
	}

	return true
}

type mManagerMockUpdateObject struct {
	mock              *ManagerMock
	mainExpectation   *ManagerMockUpdateObjectExpectation
	expectationSeries []*ManagerMockUpdateObjectExpectation
}

type ManagerMockUpdateObjectExpectation struct {
	input  *ManagerMockUpdateObjectInput
	result *ManagerMockUpdateObjectResult
}

type ManagerMockUpdateObjectInput struct {
	p  context.Context
	p1 insolar.Reference
	p2 insolar.Reference
	p3 ObjectDescriptor
	p4 []byte
}

type ManagerMockUpdateObjectResult struct {
	r  ObjectDescriptor
	r1 error
}

//Expect specifies that invocation of Manager.UpdateObject is expected from 1 to Infinity times
func (m *mManagerMockUpdateObject) Expect(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) *mManagerMockUpdateObject {
	m.mock.UpdateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockUpdateObjectExpectation{}
	}
	m.mainExpectation.input = &ManagerMockUpdateObjectInput{p, p1, p2, p3, p4}
	return m
}

//Return specifies results of invocation of Manager.UpdateObject
func (m *mManagerMockUpdateObject) Return(r ObjectDescriptor, r1 error) *ManagerMock {
	m.mock.UpdateObjectFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &ManagerMockUpdateObjectExpectation{}
	}
	m.mainExpectation.result = &ManagerMockUpdateObjectResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Manager.UpdateObject is expected once
func (m *mManagerMockUpdateObject) ExpectOnce(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) *ManagerMockUpdateObjectExpectation {
	m.mock.UpdateObjectFunc = nil
	m.mainExpectation = nil

	expectation := &ManagerMockUpdateObjectExpectation{}
	expectation.input = &ManagerMockUpdateObjectInput{p, p1, p2, p3, p4}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *ManagerMockUpdateObjectExpectation) Return(r ObjectDescriptor, r1 error) {
	e.result = &ManagerMockUpdateObjectResult{r, r1}
}

//Set uses given function f as a mock of Manager.UpdateObject method
func (m *mManagerMockUpdateObject) Set(f func(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) (r ObjectDescriptor, r1 error)) *ManagerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.UpdateObjectFunc = f
	return m.mock
}

//UpdateObject implements github.com/insolar/insolar/internal/ledger/artifact.Manager interface
func (m *ManagerMock) UpdateObject(p context.Context, p1 insolar.Reference, p2 insolar.Reference, p3 ObjectDescriptor, p4 []byte) (r ObjectDescriptor, r1 error) {
	counter := atomic.AddUint64(&m.UpdateObjectPreCounter, 1)
	defer atomic.AddUint64(&m.UpdateObjectCounter, 1)

	if len(m.UpdateObjectMock.expectationSeries) > 0 {
		if counter > uint64(len(m.UpdateObjectMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to ManagerMock.UpdateObject. %v %v %v %v %v", p, p1, p2, p3, p4)
			return
		}

		input := m.UpdateObjectMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, ManagerMockUpdateObjectInput{p, p1, p2, p3, p4}, "Manager.UpdateObject got unexpected parameters")

		result := m.UpdateObjectMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.UpdateObject")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdateObjectMock.mainExpectation != nil {

		input := m.UpdateObjectMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, ManagerMockUpdateObjectInput{p, p1, p2, p3, p4}, "Manager.UpdateObject got unexpected parameters")
		}

		result := m.UpdateObjectMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the ManagerMock.UpdateObject")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.UpdateObjectFunc == nil {
		m.t.Fatalf("Unexpected call to ManagerMock.UpdateObject. %v %v %v %v %v", p, p1, p2, p3, p4)
		return
	}

	return m.UpdateObjectFunc(p, p1, p2, p3, p4)
}

//UpdateObjectMinimockCounter returns a count of ManagerMock.UpdateObjectFunc invocations
func (m *ManagerMock) UpdateObjectMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateObjectCounter)
}

//UpdateObjectMinimockPreCounter returns the value of ManagerMock.UpdateObject invocations
func (m *ManagerMock) UpdateObjectMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.UpdateObjectPreCounter)
}

//UpdateObjectFinished returns true if mock invocations count is ok
func (m *ManagerMock) UpdateObjectFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.UpdateObjectMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.UpdateObjectCounter) == uint64(len(m.UpdateObjectMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.UpdateObjectMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.UpdateObjectCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.UpdateObjectFunc != nil {
		return atomic.LoadUint64(&m.UpdateObjectCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ManagerMock) ValidateCallCounters() {

	if !m.ActivateObjectFinished() {
		m.t.Fatal("Expected call to ManagerMock.ActivateObject")
	}

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ManagerMock.RegisterRequest")
	}

	if !m.RegisterResultFinished() {
		m.t.Fatal("Expected call to ManagerMock.RegisterResult")
	}

	if !m.UpdateObjectFinished() {
		m.t.Fatal("Expected call to ManagerMock.UpdateObject")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ManagerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ManagerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ManagerMock) MinimockFinish() {

	if !m.ActivateObjectFinished() {
		m.t.Fatal("Expected call to ManagerMock.ActivateObject")
	}

	if !m.RegisterRequestFinished() {
		m.t.Fatal("Expected call to ManagerMock.RegisterRequest")
	}

	if !m.RegisterResultFinished() {
		m.t.Fatal("Expected call to ManagerMock.RegisterResult")
	}

	if !m.UpdateObjectFinished() {
		m.t.Fatal("Expected call to ManagerMock.UpdateObject")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ManagerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ManagerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.ActivateObjectFinished()
		ok = ok && m.RegisterRequestFinished()
		ok = ok && m.RegisterResultFinished()
		ok = ok && m.UpdateObjectFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.ActivateObjectFinished() {
				m.t.Error("Expected call to ManagerMock.ActivateObject")
			}

			if !m.RegisterRequestFinished() {
				m.t.Error("Expected call to ManagerMock.RegisterRequest")
			}

			if !m.RegisterResultFinished() {
				m.t.Error("Expected call to ManagerMock.RegisterResult")
			}

			if !m.UpdateObjectFinished() {
				m.t.Error("Expected call to ManagerMock.UpdateObject")
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
func (m *ManagerMock) AllMocksCalled() bool {

	if !m.ActivateObjectFinished() {
		return false
	}

	if !m.RegisterRequestFinished() {
		return false
	}

	if !m.RegisterResultFinished() {
		return false
	}

	if !m.UpdateObjectFinished() {
		return false
	}

	return true
}
