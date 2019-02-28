package drop

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "Packer" can be found in github.com/insolar/insolar/ledger/storage/jet/drop
*/
import (
	context "context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"
	jet "github.com/insolar/insolar/ledger/storage/jet"

	testify_assert "github.com/stretchr/testify/assert"
)

//PackerMock implements github.com/insolar/insolar/ledger/storage/jet/drop.Packer
type PackerMock struct {
	t minimock.Tester

	PackFunc       func(p context.Context, p1 core.JetID, p2 core.PulseNumber, p3 []byte) (r jet.Drop, r1 error)
	PackCounter    uint64
	PackPreCounter uint64
	PackMock       mPackerMockPack
}

//NewPackerMock returns a mock for github.com/insolar/insolar/ledger/storage/jet/drop.Packer
func NewPackerMock(t minimock.Tester) *PackerMock {
	m := &PackerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PackMock = mPackerMockPack{mock: m}

	return m
}

type mPackerMockPack struct {
	mock              *PackerMock
	mainExpectation   *PackerMockPackExpectation
	expectationSeries []*PackerMockPackExpectation
}

type PackerMockPackExpectation struct {
	input  *PackerMockPackInput
	result *PackerMockPackResult
}

type PackerMockPackInput struct {
	p  context.Context
	p1 core.JetID
	p2 core.PulseNumber
	p3 []byte
}

type PackerMockPackResult struct {
	r  jet.Drop
	r1 error
}

//Expect specifies that invocation of Packer.Pack is expected from 1 to Infinity times
func (m *mPackerMockPack) Expect(p context.Context, p1 core.JetID, p2 core.PulseNumber, p3 []byte) *mPackerMockPack {
	m.mock.PackFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PackerMockPackExpectation{}
	}
	m.mainExpectation.input = &PackerMockPackInput{p, p1, p2, p3}
	return m
}

//Return specifies results of invocation of Packer.Pack
func (m *mPackerMockPack) Return(r jet.Drop, r1 error) *PackerMock {
	m.mock.PackFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &PackerMockPackExpectation{}
	}
	m.mainExpectation.result = &PackerMockPackResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of Packer.Pack is expected once
func (m *mPackerMockPack) ExpectOnce(p context.Context, p1 core.JetID, p2 core.PulseNumber, p3 []byte) *PackerMockPackExpectation {
	m.mock.PackFunc = nil
	m.mainExpectation = nil

	expectation := &PackerMockPackExpectation{}
	expectation.input = &PackerMockPackInput{p, p1, p2, p3}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *PackerMockPackExpectation) Return(r jet.Drop, r1 error) {
	e.result = &PackerMockPackResult{r, r1}
}

//Set uses given function f as a mock of Packer.Pack method
func (m *mPackerMockPack) Set(f func(p context.Context, p1 core.JetID, p2 core.PulseNumber, p3 []byte) (r jet.Drop, r1 error)) *PackerMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.PackFunc = f
	return m.mock
}

//Pack implements github.com/insolar/insolar/ledger/storage/jet/drop.Packer interface
func (m *PackerMock) Pack(p context.Context, p1 core.JetID, p2 core.PulseNumber, p3 []byte) (r jet.Drop, r1 error) {
	counter := atomic.AddUint64(&m.PackPreCounter, 1)
	defer atomic.AddUint64(&m.PackCounter, 1)

	if len(m.PackMock.expectationSeries) > 0 {
		if counter > uint64(len(m.PackMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to PackerMock.Pack. %v %v %v %v", p, p1, p2, p3)
			return
		}

		input := m.PackMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, PackerMockPackInput{p, p1, p2, p3}, "Packer.Pack got unexpected parameters")

		result := m.PackMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the PackerMock.Pack")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PackMock.mainExpectation != nil {

		input := m.PackMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, PackerMockPackInput{p, p1, p2, p3}, "Packer.Pack got unexpected parameters")
		}

		result := m.PackMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the PackerMock.Pack")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.PackFunc == nil {
		m.t.Fatalf("Unexpected call to PackerMock.Pack. %v %v %v %v", p, p1, p2, p3)
		return
	}

	return m.PackFunc(p, p1, p2, p3)
}

//PackMinimockCounter returns a count of PackerMock.PackFunc invocations
func (m *PackerMock) PackMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PackCounter)
}

//PackMinimockPreCounter returns the value of PackerMock.Pack invocations
func (m *PackerMock) PackMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PackPreCounter)
}

//PackFinished returns true if mock invocations count is ok
func (m *PackerMock) PackFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.PackMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.PackCounter) == uint64(len(m.PackMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.PackMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.PackCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.PackFunc != nil {
		return atomic.LoadUint64(&m.PackCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PackerMock) ValidateCallCounters() {

	if !m.PackFinished() {
		m.t.Fatal("Expected call to PackerMock.Pack")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *PackerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *PackerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *PackerMock) MinimockFinish() {

	if !m.PackFinished() {
		m.t.Fatal("Expected call to PackerMock.Pack")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *PackerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *PackerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.PackFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.PackFinished() {
				m.t.Error("Expected call to PackerMock.Pack")
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
func (m *PackerMock) AllMocksCalled() bool {

	if !m.PackFinished() {
		return false
	}

	return true
}
