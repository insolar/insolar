package messagebus

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.9
The original interface "sender" can be found in github.com/insolar/insolar/messagebus
*/
import (
	context "context"
	io "io"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	core "github.com/insolar/insolar/core"

	testify_assert "github.com/stretchr/testify/assert"
)

//senderMock implements github.com/insolar/insolar/messagebus.sender
type senderMock struct {
	t minimock.Tester

	CreateParcelFunc       func(p context.Context, p1 core.Message, p2 core.DelegationToken) (r core.Parcel, r1 error)
	CreateParcelCounter    uint64
	CreateParcelPreCounter uint64
	CreateParcelMock       msenderMockCreateParcel

	MustRegisterFunc       func(p core.MessageType, p1 core.MessageHandler)
	MustRegisterCounter    uint64
	MustRegisterPreCounter uint64
	MustRegisterMock       msenderMockMustRegister

	NewPlayerFunc       func(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error)
	NewPlayerCounter    uint64
	NewPlayerPreCounter uint64
	NewPlayerMock       msenderMockNewPlayer

	NewRecorderFunc       func(p context.Context) (r core.MessageBus, r1 error)
	NewRecorderCounter    uint64
	NewRecorderPreCounter uint64
	NewRecorderMock       msenderMockNewRecorder

	RegisterFunc       func(p core.MessageType, p1 core.MessageHandler) (r error)
	RegisterCounter    uint64
	RegisterPreCounter uint64
	RegisterMock       msenderMockRegister

	SendFunc       func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error)
	SendCounter    uint64
	SendPreCounter uint64
	SendMock       msenderMockSend

	SendParcelFunc       func(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) (r core.Reply, r1 error)
	SendParcelCounter    uint64
	SendParcelPreCounter uint64
	SendParcelMock       msenderMockSendParcel

	WriteTapeFunc       func(p context.Context, p1 io.Writer) (r error)
	WriteTapeCounter    uint64
	WriteTapePreCounter uint64
	WriteTapeMock       msenderMockWriteTape
}

//NewsenderMock returns a mock for github.com/insolar/insolar/messagebus.sender
func NewsenderMock(t minimock.Tester) *senderMock {
	m := &senderMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateParcelMock = msenderMockCreateParcel{mock: m}
	m.MustRegisterMock = msenderMockMustRegister{mock: m}
	m.NewPlayerMock = msenderMockNewPlayer{mock: m}
	m.NewRecorderMock = msenderMockNewRecorder{mock: m}
	m.RegisterMock = msenderMockRegister{mock: m}
	m.SendMock = msenderMockSend{mock: m}
	m.SendParcelMock = msenderMockSendParcel{mock: m}
	m.WriteTapeMock = msenderMockWriteTape{mock: m}

	return m
}

type msenderMockCreateParcel struct {
	mock              *senderMock
	mainExpectation   *senderMockCreateParcelExpectation
	expectationSeries []*senderMockCreateParcelExpectation
}

type senderMockCreateParcelExpectation struct {
	input  *senderMockCreateParcelInput
	result *senderMockCreateParcelResult
}

type senderMockCreateParcelInput struct {
	p  context.Context
	p1 core.Message
	p2 core.DelegationToken
}

type senderMockCreateParcelResult struct {
	r  core.Parcel
	r1 error
}

//Expect specifies that invocation of sender.CreateParcel is expected from 1 to Infinity times
func (m *msenderMockCreateParcel) Expect(p context.Context, p1 core.Message, p2 core.DelegationToken) *msenderMockCreateParcel {
	m.mock.CreateParcelFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockCreateParcelExpectation{}
	}
	m.mainExpectation.input = &senderMockCreateParcelInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of sender.CreateParcel
func (m *msenderMockCreateParcel) Return(r core.Parcel, r1 error) *senderMock {
	m.mock.CreateParcelFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockCreateParcelExpectation{}
	}
	m.mainExpectation.result = &senderMockCreateParcelResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of sender.CreateParcel is expected once
func (m *msenderMockCreateParcel) ExpectOnce(p context.Context, p1 core.Message, p2 core.DelegationToken) *senderMockCreateParcelExpectation {
	m.mock.CreateParcelFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockCreateParcelExpectation{}
	expectation.input = &senderMockCreateParcelInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *senderMockCreateParcelExpectation) Return(r core.Parcel, r1 error) {
	e.result = &senderMockCreateParcelResult{r, r1}
}

//Set uses given function f as a mock of sender.CreateParcel method
func (m *msenderMockCreateParcel) Set(f func(p context.Context, p1 core.Message, p2 core.DelegationToken) (r core.Parcel, r1 error)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.CreateParcelFunc = f
	return m.mock
}

//CreateParcel implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) CreateParcel(p context.Context, p1 core.Message, p2 core.DelegationToken) (r core.Parcel, r1 error) {
	counter := atomic.AddUint64(&m.CreateParcelPreCounter, 1)
	defer atomic.AddUint64(&m.CreateParcelCounter, 1)

	if len(m.CreateParcelMock.expectationSeries) > 0 {
		if counter > uint64(len(m.CreateParcelMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.CreateParcel. %v %v %v", p, p1, p2)
			return
		}

		input := m.CreateParcelMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockCreateParcelInput{p, p1, p2}, "sender.CreateParcel got unexpected parameters")

		result := m.CreateParcelMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.CreateParcel")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CreateParcelMock.mainExpectation != nil {

		input := m.CreateParcelMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockCreateParcelInput{p, p1, p2}, "sender.CreateParcel got unexpected parameters")
		}

		result := m.CreateParcelMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.CreateParcel")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.CreateParcelFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.CreateParcel. %v %v %v", p, p1, p2)
		return
	}

	return m.CreateParcelFunc(p, p1, p2)
}

//CreateParcelMinimockCounter returns a count of senderMock.CreateParcelFunc invocations
func (m *senderMock) CreateParcelMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.CreateParcelCounter)
}

//CreateParcelMinimockPreCounter returns the value of senderMock.CreateParcel invocations
func (m *senderMock) CreateParcelMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.CreateParcelPreCounter)
}

//CreateParcelFinished returns true if mock invocations count is ok
func (m *senderMock) CreateParcelFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.CreateParcelMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.CreateParcelCounter) == uint64(len(m.CreateParcelMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.CreateParcelMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.CreateParcelCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.CreateParcelFunc != nil {
		return atomic.LoadUint64(&m.CreateParcelCounter) > 0
	}

	return true
}

type msenderMockMustRegister struct {
	mock              *senderMock
	mainExpectation   *senderMockMustRegisterExpectation
	expectationSeries []*senderMockMustRegisterExpectation
}

type senderMockMustRegisterExpectation struct {
	input *senderMockMustRegisterInput
}

type senderMockMustRegisterInput struct {
	p  core.MessageType
	p1 core.MessageHandler
}

//Expect specifies that invocation of sender.MustRegister is expected from 1 to Infinity times
func (m *msenderMockMustRegister) Expect(p core.MessageType, p1 core.MessageHandler) *msenderMockMustRegister {
	m.mock.MustRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockMustRegisterExpectation{}
	}
	m.mainExpectation.input = &senderMockMustRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of sender.MustRegister
func (m *msenderMockMustRegister) Return() *senderMock {
	m.mock.MustRegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockMustRegisterExpectation{}
	}

	return m.mock
}

//ExpectOnce specifies that invocation of sender.MustRegister is expected once
func (m *msenderMockMustRegister) ExpectOnce(p core.MessageType, p1 core.MessageHandler) *senderMockMustRegisterExpectation {
	m.mock.MustRegisterFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockMustRegisterExpectation{}
	expectation.input = &senderMockMustRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

//Set uses given function f as a mock of sender.MustRegister method
func (m *msenderMockMustRegister) Set(f func(p core.MessageType, p1 core.MessageHandler)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.MustRegisterFunc = f
	return m.mock
}

//MustRegister implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) MustRegister(p core.MessageType, p1 core.MessageHandler) {
	counter := atomic.AddUint64(&m.MustRegisterPreCounter, 1)
	defer atomic.AddUint64(&m.MustRegisterCounter, 1)

	if len(m.MustRegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.MustRegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.MustRegister. %v %v", p, p1)
			return
		}

		input := m.MustRegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockMustRegisterInput{p, p1}, "sender.MustRegister got unexpected parameters")

		return
	}

	if m.MustRegisterMock.mainExpectation != nil {

		input := m.MustRegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockMustRegisterInput{p, p1}, "sender.MustRegister got unexpected parameters")
		}

		return
	}

	if m.MustRegisterFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.MustRegister. %v %v", p, p1)
		return
	}

	m.MustRegisterFunc(p, p1)
}

//MustRegisterMinimockCounter returns a count of senderMock.MustRegisterFunc invocations
func (m *senderMock) MustRegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustRegisterCounter)
}

//MustRegisterMinimockPreCounter returns the value of senderMock.MustRegister invocations
func (m *senderMock) MustRegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustRegisterPreCounter)
}

//MustRegisterFinished returns true if mock invocations count is ok
func (m *senderMock) MustRegisterFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.MustRegisterMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.MustRegisterCounter) == uint64(len(m.MustRegisterMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.MustRegisterMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.MustRegisterCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.MustRegisterFunc != nil {
		return atomic.LoadUint64(&m.MustRegisterCounter) > 0
	}

	return true
}

type msenderMockNewPlayer struct {
	mock              *senderMock
	mainExpectation   *senderMockNewPlayerExpectation
	expectationSeries []*senderMockNewPlayerExpectation
}

type senderMockNewPlayerExpectation struct {
	input  *senderMockNewPlayerInput
	result *senderMockNewPlayerResult
}

type senderMockNewPlayerInput struct {
	p  context.Context
	p1 io.Reader
}

type senderMockNewPlayerResult struct {
	r  core.MessageBus
	r1 error
}

//Expect specifies that invocation of sender.NewPlayer is expected from 1 to Infinity times
func (m *msenderMockNewPlayer) Expect(p context.Context, p1 io.Reader) *msenderMockNewPlayer {
	m.mock.NewPlayerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockNewPlayerExpectation{}
	}
	m.mainExpectation.input = &senderMockNewPlayerInput{p, p1}
	return m
}

//Return specifies results of invocation of sender.NewPlayer
func (m *msenderMockNewPlayer) Return(r core.MessageBus, r1 error) *senderMock {
	m.mock.NewPlayerFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockNewPlayerExpectation{}
	}
	m.mainExpectation.result = &senderMockNewPlayerResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of sender.NewPlayer is expected once
func (m *msenderMockNewPlayer) ExpectOnce(p context.Context, p1 io.Reader) *senderMockNewPlayerExpectation {
	m.mock.NewPlayerFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockNewPlayerExpectation{}
	expectation.input = &senderMockNewPlayerInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *senderMockNewPlayerExpectation) Return(r core.MessageBus, r1 error) {
	e.result = &senderMockNewPlayerResult{r, r1}
}

//Set uses given function f as a mock of sender.NewPlayer method
func (m *msenderMockNewPlayer) Set(f func(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NewPlayerFunc = f
	return m.mock
}

//NewPlayer implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) NewPlayer(p context.Context, p1 io.Reader) (r core.MessageBus, r1 error) {
	counter := atomic.AddUint64(&m.NewPlayerPreCounter, 1)
	defer atomic.AddUint64(&m.NewPlayerCounter, 1)

	if len(m.NewPlayerMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NewPlayerMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.NewPlayer. %v %v", p, p1)
			return
		}

		input := m.NewPlayerMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockNewPlayerInput{p, p1}, "sender.NewPlayer got unexpected parameters")

		result := m.NewPlayerMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.NewPlayer")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewPlayerMock.mainExpectation != nil {

		input := m.NewPlayerMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockNewPlayerInput{p, p1}, "sender.NewPlayer got unexpected parameters")
		}

		result := m.NewPlayerMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.NewPlayer")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewPlayerFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.NewPlayer. %v %v", p, p1)
		return
	}

	return m.NewPlayerFunc(p, p1)
}

//NewPlayerMinimockCounter returns a count of senderMock.NewPlayerFunc invocations
func (m *senderMock) NewPlayerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewPlayerCounter)
}

//NewPlayerMinimockPreCounter returns the value of senderMock.NewPlayer invocations
func (m *senderMock) NewPlayerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewPlayerPreCounter)
}

//NewPlayerFinished returns true if mock invocations count is ok
func (m *senderMock) NewPlayerFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NewPlayerMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NewPlayerCounter) == uint64(len(m.NewPlayerMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NewPlayerMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NewPlayerCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NewPlayerFunc != nil {
		return atomic.LoadUint64(&m.NewPlayerCounter) > 0
	}

	return true
}

type msenderMockNewRecorder struct {
	mock              *senderMock
	mainExpectation   *senderMockNewRecorderExpectation
	expectationSeries []*senderMockNewRecorderExpectation
}

type senderMockNewRecorderExpectation struct {
	input  *senderMockNewRecorderInput
	result *senderMockNewRecorderResult
}

type senderMockNewRecorderInput struct {
	p context.Context
}

type senderMockNewRecorderResult struct {
	r  core.MessageBus
	r1 error
}

//Expect specifies that invocation of sender.NewRecorder is expected from 1 to Infinity times
func (m *msenderMockNewRecorder) Expect(p context.Context) *msenderMockNewRecorder {
	m.mock.NewRecorderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockNewRecorderExpectation{}
	}
	m.mainExpectation.input = &senderMockNewRecorderInput{p}
	return m
}

//Return specifies results of invocation of sender.NewRecorder
func (m *msenderMockNewRecorder) Return(r core.MessageBus, r1 error) *senderMock {
	m.mock.NewRecorderFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockNewRecorderExpectation{}
	}
	m.mainExpectation.result = &senderMockNewRecorderResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of sender.NewRecorder is expected once
func (m *msenderMockNewRecorder) ExpectOnce(p context.Context) *senderMockNewRecorderExpectation {
	m.mock.NewRecorderFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockNewRecorderExpectation{}
	expectation.input = &senderMockNewRecorderInput{p}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *senderMockNewRecorderExpectation) Return(r core.MessageBus, r1 error) {
	e.result = &senderMockNewRecorderResult{r, r1}
}

//Set uses given function f as a mock of sender.NewRecorder method
func (m *msenderMockNewRecorder) Set(f func(p context.Context) (r core.MessageBus, r1 error)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.NewRecorderFunc = f
	return m.mock
}

//NewRecorder implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) NewRecorder(p context.Context) (r core.MessageBus, r1 error) {
	counter := atomic.AddUint64(&m.NewRecorderPreCounter, 1)
	defer atomic.AddUint64(&m.NewRecorderCounter, 1)

	if len(m.NewRecorderMock.expectationSeries) > 0 {
		if counter > uint64(len(m.NewRecorderMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.NewRecorder. %v", p)
			return
		}

		input := m.NewRecorderMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockNewRecorderInput{p}, "sender.NewRecorder got unexpected parameters")

		result := m.NewRecorderMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.NewRecorder")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewRecorderMock.mainExpectation != nil {

		input := m.NewRecorderMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockNewRecorderInput{p}, "sender.NewRecorder got unexpected parameters")
		}

		result := m.NewRecorderMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.NewRecorder")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.NewRecorderFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.NewRecorder. %v", p)
		return
	}

	return m.NewRecorderFunc(p)
}

//NewRecorderMinimockCounter returns a count of senderMock.NewRecorderFunc invocations
func (m *senderMock) NewRecorderMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.NewRecorderCounter)
}

//NewRecorderMinimockPreCounter returns the value of senderMock.NewRecorder invocations
func (m *senderMock) NewRecorderMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.NewRecorderPreCounter)
}

//NewRecorderFinished returns true if mock invocations count is ok
func (m *senderMock) NewRecorderFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.NewRecorderMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.NewRecorderCounter) == uint64(len(m.NewRecorderMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.NewRecorderMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.NewRecorderCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.NewRecorderFunc != nil {
		return atomic.LoadUint64(&m.NewRecorderCounter) > 0
	}

	return true
}

type msenderMockRegister struct {
	mock              *senderMock
	mainExpectation   *senderMockRegisterExpectation
	expectationSeries []*senderMockRegisterExpectation
}

type senderMockRegisterExpectation struct {
	input  *senderMockRegisterInput
	result *senderMockRegisterResult
}

type senderMockRegisterInput struct {
	p  core.MessageType
	p1 core.MessageHandler
}

type senderMockRegisterResult struct {
	r error
}

//Expect specifies that invocation of sender.Register is expected from 1 to Infinity times
func (m *msenderMockRegister) Expect(p core.MessageType, p1 core.MessageHandler) *msenderMockRegister {
	m.mock.RegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockRegisterExpectation{}
	}
	m.mainExpectation.input = &senderMockRegisterInput{p, p1}
	return m
}

//Return specifies results of invocation of sender.Register
func (m *msenderMockRegister) Return(r error) *senderMock {
	m.mock.RegisterFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockRegisterExpectation{}
	}
	m.mainExpectation.result = &senderMockRegisterResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of sender.Register is expected once
func (m *msenderMockRegister) ExpectOnce(p core.MessageType, p1 core.MessageHandler) *senderMockRegisterExpectation {
	m.mock.RegisterFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockRegisterExpectation{}
	expectation.input = &senderMockRegisterInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *senderMockRegisterExpectation) Return(r error) {
	e.result = &senderMockRegisterResult{r}
}

//Set uses given function f as a mock of sender.Register method
func (m *msenderMockRegister) Set(f func(p core.MessageType, p1 core.MessageHandler) (r error)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.RegisterFunc = f
	return m.mock
}

//Register implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) Register(p core.MessageType, p1 core.MessageHandler) (r error) {
	counter := atomic.AddUint64(&m.RegisterPreCounter, 1)
	defer atomic.AddUint64(&m.RegisterCounter, 1)

	if len(m.RegisterMock.expectationSeries) > 0 {
		if counter > uint64(len(m.RegisterMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.Register. %v %v", p, p1)
			return
		}

		input := m.RegisterMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockRegisterInput{p, p1}, "sender.Register got unexpected parameters")

		result := m.RegisterMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.Register")
			return
		}

		r = result.r

		return
	}

	if m.RegisterMock.mainExpectation != nil {

		input := m.RegisterMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockRegisterInput{p, p1}, "sender.Register got unexpected parameters")
		}

		result := m.RegisterMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.Register")
		}

		r = result.r

		return
	}

	if m.RegisterFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.Register. %v %v", p, p1)
		return
	}

	return m.RegisterFunc(p, p1)
}

//RegisterMinimockCounter returns a count of senderMock.RegisterFunc invocations
func (m *senderMock) RegisterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterCounter)
}

//RegisterMinimockPreCounter returns the value of senderMock.Register invocations
func (m *senderMock) RegisterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RegisterPreCounter)
}

//RegisterFinished returns true if mock invocations count is ok
func (m *senderMock) RegisterFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.RegisterMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.RegisterCounter) == uint64(len(m.RegisterMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.RegisterMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.RegisterCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.RegisterFunc != nil {
		return atomic.LoadUint64(&m.RegisterCounter) > 0
	}

	return true
}

type msenderMockSend struct {
	mock              *senderMock
	mainExpectation   *senderMockSendExpectation
	expectationSeries []*senderMockSendExpectation
}

type senderMockSendExpectation struct {
	input  *senderMockSendInput
	result *senderMockSendResult
}

type senderMockSendInput struct {
	p  context.Context
	p1 core.Message
	p2 *core.MessageSendOptions
}

type senderMockSendResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of sender.Send is expected from 1 to Infinity times
func (m *msenderMockSend) Expect(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) *msenderMockSend {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockSendExpectation{}
	}
	m.mainExpectation.input = &senderMockSendInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of sender.Send
func (m *msenderMockSend) Return(r core.Reply, r1 error) *senderMock {
	m.mock.SendFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockSendExpectation{}
	}
	m.mainExpectation.result = &senderMockSendResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of sender.Send is expected once
func (m *msenderMockSend) ExpectOnce(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) *senderMockSendExpectation {
	m.mock.SendFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockSendExpectation{}
	expectation.input = &senderMockSendInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *senderMockSendExpectation) Return(r core.Reply, r1 error) {
	e.result = &senderMockSendResult{r, r1}
}

//Set uses given function f as a mock of sender.Send method
func (m *msenderMockSend) Set(f func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendFunc = f
	return m.mock
}

//Send implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) Send(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.SendPreCounter, 1)
	defer atomic.AddUint64(&m.SendCounter, 1)

	if len(m.SendMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.Send. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockSendInput{p, p1, p2}, "sender.Send got unexpected parameters")

		result := m.SendMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.Send")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendMock.mainExpectation != nil {

		input := m.SendMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockSendInput{p, p1, p2}, "sender.Send got unexpected parameters")
		}

		result := m.SendMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.Send")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.Send. %v %v %v", p, p1, p2)
		return
	}

	return m.SendFunc(p, p1, p2)
}

//SendMinimockCounter returns a count of senderMock.SendFunc invocations
func (m *senderMock) SendMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendCounter)
}

//SendMinimockPreCounter returns the value of senderMock.Send invocations
func (m *senderMock) SendMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendPreCounter)
}

//SendFinished returns true if mock invocations count is ok
func (m *senderMock) SendFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendCounter) == uint64(len(m.SendMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendFunc != nil {
		return atomic.LoadUint64(&m.SendCounter) > 0
	}

	return true
}

type msenderMockSendParcel struct {
	mock              *senderMock
	mainExpectation   *senderMockSendParcelExpectation
	expectationSeries []*senderMockSendParcelExpectation
}

type senderMockSendParcelExpectation struct {
	input  *senderMockSendParcelInput
	result *senderMockSendParcelResult
}

type senderMockSendParcelInput struct {
	p  context.Context
	p1 core.Parcel
	p2 *core.MessageSendOptions
}

type senderMockSendParcelResult struct {
	r  core.Reply
	r1 error
}

//Expect specifies that invocation of sender.SendParcel is expected from 1 to Infinity times
func (m *msenderMockSendParcel) Expect(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) *msenderMockSendParcel {
	m.mock.SendParcelFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockSendParcelExpectation{}
	}
	m.mainExpectation.input = &senderMockSendParcelInput{p, p1, p2}
	return m
}

//Return specifies results of invocation of sender.SendParcel
func (m *msenderMockSendParcel) Return(r core.Reply, r1 error) *senderMock {
	m.mock.SendParcelFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockSendParcelExpectation{}
	}
	m.mainExpectation.result = &senderMockSendParcelResult{r, r1}
	return m.mock
}

//ExpectOnce specifies that invocation of sender.SendParcel is expected once
func (m *msenderMockSendParcel) ExpectOnce(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) *senderMockSendParcelExpectation {
	m.mock.SendParcelFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockSendParcelExpectation{}
	expectation.input = &senderMockSendParcelInput{p, p1, p2}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *senderMockSendParcelExpectation) Return(r core.Reply, r1 error) {
	e.result = &senderMockSendParcelResult{r, r1}
}

//Set uses given function f as a mock of sender.SendParcel method
func (m *msenderMockSendParcel) Set(f func(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) (r core.Reply, r1 error)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.SendParcelFunc = f
	return m.mock
}

//SendParcel implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) SendParcel(p context.Context, p1 core.Parcel, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
	counter := atomic.AddUint64(&m.SendParcelPreCounter, 1)
	defer atomic.AddUint64(&m.SendParcelCounter, 1)

	if len(m.SendParcelMock.expectationSeries) > 0 {
		if counter > uint64(len(m.SendParcelMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.SendParcel. %v %v %v", p, p1, p2)
			return
		}

		input := m.SendParcelMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockSendParcelInput{p, p1, p2}, "sender.SendParcel got unexpected parameters")

		result := m.SendParcelMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.SendParcel")
			return
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendParcelMock.mainExpectation != nil {

		input := m.SendParcelMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockSendParcelInput{p, p1, p2}, "sender.SendParcel got unexpected parameters")
		}

		result := m.SendParcelMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.SendParcel")
		}

		r = result.r
		r1 = result.r1

		return
	}

	if m.SendParcelFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.SendParcel. %v %v %v", p, p1, p2)
		return
	}

	return m.SendParcelFunc(p, p1, p2)
}

//SendParcelMinimockCounter returns a count of senderMock.SendParcelFunc invocations
func (m *senderMock) SendParcelMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SendParcelCounter)
}

//SendParcelMinimockPreCounter returns the value of senderMock.SendParcel invocations
func (m *senderMock) SendParcelMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SendParcelPreCounter)
}

//SendParcelFinished returns true if mock invocations count is ok
func (m *senderMock) SendParcelFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.SendParcelMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.SendParcelCounter) == uint64(len(m.SendParcelMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.SendParcelMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.SendParcelCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.SendParcelFunc != nil {
		return atomic.LoadUint64(&m.SendParcelCounter) > 0
	}

	return true
}

type msenderMockWriteTape struct {
	mock              *senderMock
	mainExpectation   *senderMockWriteTapeExpectation
	expectationSeries []*senderMockWriteTapeExpectation
}

type senderMockWriteTapeExpectation struct {
	input  *senderMockWriteTapeInput
	result *senderMockWriteTapeResult
}

type senderMockWriteTapeInput struct {
	p  context.Context
	p1 io.Writer
}

type senderMockWriteTapeResult struct {
	r error
}

//Expect specifies that invocation of sender.WriteTape is expected from 1 to Infinity times
func (m *msenderMockWriteTape) Expect(p context.Context, p1 io.Writer) *msenderMockWriteTape {
	m.mock.WriteTapeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockWriteTapeExpectation{}
	}
	m.mainExpectation.input = &senderMockWriteTapeInput{p, p1}
	return m
}

//Return specifies results of invocation of sender.WriteTape
func (m *msenderMockWriteTape) Return(r error) *senderMock {
	m.mock.WriteTapeFunc = nil
	m.expectationSeries = nil

	if m.mainExpectation == nil {
		m.mainExpectation = &senderMockWriteTapeExpectation{}
	}
	m.mainExpectation.result = &senderMockWriteTapeResult{r}
	return m.mock
}

//ExpectOnce specifies that invocation of sender.WriteTape is expected once
func (m *msenderMockWriteTape) ExpectOnce(p context.Context, p1 io.Writer) *senderMockWriteTapeExpectation {
	m.mock.WriteTapeFunc = nil
	m.mainExpectation = nil

	expectation := &senderMockWriteTapeExpectation{}
	expectation.input = &senderMockWriteTapeInput{p, p1}
	m.expectationSeries = append(m.expectationSeries, expectation)
	return expectation
}

func (e *senderMockWriteTapeExpectation) Return(r error) {
	e.result = &senderMockWriteTapeResult{r}
}

//Set uses given function f as a mock of sender.WriteTape method
func (m *msenderMockWriteTape) Set(f func(p context.Context, p1 io.Writer) (r error)) *senderMock {
	m.mainExpectation = nil
	m.expectationSeries = nil

	m.mock.WriteTapeFunc = f
	return m.mock
}

//WriteTape implements github.com/insolar/insolar/messagebus.sender interface
func (m *senderMock) WriteTape(p context.Context, p1 io.Writer) (r error) {
	counter := atomic.AddUint64(&m.WriteTapePreCounter, 1)
	defer atomic.AddUint64(&m.WriteTapeCounter, 1)

	if len(m.WriteTapeMock.expectationSeries) > 0 {
		if counter > uint64(len(m.WriteTapeMock.expectationSeries)) {
			m.t.Fatalf("Unexpected call to senderMock.WriteTape. %v %v", p, p1)
			return
		}

		input := m.WriteTapeMock.expectationSeries[counter-1].input
		testify_assert.Equal(m.t, *input, senderMockWriteTapeInput{p, p1}, "sender.WriteTape got unexpected parameters")

		result := m.WriteTapeMock.expectationSeries[counter-1].result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.WriteTape")
			return
		}

		r = result.r

		return
	}

	if m.WriteTapeMock.mainExpectation != nil {

		input := m.WriteTapeMock.mainExpectation.input
		if input != nil {
			testify_assert.Equal(m.t, *input, senderMockWriteTapeInput{p, p1}, "sender.WriteTape got unexpected parameters")
		}

		result := m.WriteTapeMock.mainExpectation.result
		if result == nil {
			m.t.Fatal("No results are set for the senderMock.WriteTape")
		}

		r = result.r

		return
	}

	if m.WriteTapeFunc == nil {
		m.t.Fatalf("Unexpected call to senderMock.WriteTape. %v %v", p, p1)
		return
	}

	return m.WriteTapeFunc(p, p1)
}

//WriteTapeMinimockCounter returns a count of senderMock.WriteTapeFunc invocations
func (m *senderMock) WriteTapeMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WriteTapeCounter)
}

//WriteTapeMinimockPreCounter returns the value of senderMock.WriteTape invocations
func (m *senderMock) WriteTapeMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WriteTapePreCounter)
}

//WriteTapeFinished returns true if mock invocations count is ok
func (m *senderMock) WriteTapeFinished() bool {
	// if expectation series were set then invocations count should be equal to expectations count
	if len(m.WriteTapeMock.expectationSeries) > 0 {
		return atomic.LoadUint64(&m.WriteTapeCounter) == uint64(len(m.WriteTapeMock.expectationSeries))
	}

	// if main expectation was set then invocations count should be greater than zero
	if m.WriteTapeMock.mainExpectation != nil {
		return atomic.LoadUint64(&m.WriteTapeCounter) > 0
	}

	// if func was set then invocations count should be greater than zero
	if m.WriteTapeFunc != nil {
		return atomic.LoadUint64(&m.WriteTapeCounter) > 0
	}

	return true
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *senderMock) ValidateCallCounters() {

	if !m.CreateParcelFinished() {
		m.t.Fatal("Expected call to senderMock.CreateParcel")
	}

	if !m.MustRegisterFinished() {
		m.t.Fatal("Expected call to senderMock.MustRegister")
	}

	if !m.NewPlayerFinished() {
		m.t.Fatal("Expected call to senderMock.NewPlayer")
	}

	if !m.NewRecorderFinished() {
		m.t.Fatal("Expected call to senderMock.NewRecorder")
	}

	if !m.RegisterFinished() {
		m.t.Fatal("Expected call to senderMock.Register")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to senderMock.Send")
	}

	if !m.SendParcelFinished() {
		m.t.Fatal("Expected call to senderMock.SendParcel")
	}

	if !m.WriteTapeFinished() {
		m.t.Fatal("Expected call to senderMock.WriteTape")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *senderMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *senderMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *senderMock) MinimockFinish() {

	if !m.CreateParcelFinished() {
		m.t.Fatal("Expected call to senderMock.CreateParcel")
	}

	if !m.MustRegisterFinished() {
		m.t.Fatal("Expected call to senderMock.MustRegister")
	}

	if !m.NewPlayerFinished() {
		m.t.Fatal("Expected call to senderMock.NewPlayer")
	}

	if !m.NewRecorderFinished() {
		m.t.Fatal("Expected call to senderMock.NewRecorder")
	}

	if !m.RegisterFinished() {
		m.t.Fatal("Expected call to senderMock.Register")
	}

	if !m.SendFinished() {
		m.t.Fatal("Expected call to senderMock.Send")
	}

	if !m.SendParcelFinished() {
		m.t.Fatal("Expected call to senderMock.SendParcel")
	}

	if !m.WriteTapeFinished() {
		m.t.Fatal("Expected call to senderMock.WriteTape")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *senderMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *senderMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && m.CreateParcelFinished()
		ok = ok && m.MustRegisterFinished()
		ok = ok && m.NewPlayerFinished()
		ok = ok && m.NewRecorderFinished()
		ok = ok && m.RegisterFinished()
		ok = ok && m.SendFinished()
		ok = ok && m.SendParcelFinished()
		ok = ok && m.WriteTapeFinished()

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if !m.CreateParcelFinished() {
				m.t.Error("Expected call to senderMock.CreateParcel")
			}

			if !m.MustRegisterFinished() {
				m.t.Error("Expected call to senderMock.MustRegister")
			}

			if !m.NewPlayerFinished() {
				m.t.Error("Expected call to senderMock.NewPlayer")
			}

			if !m.NewRecorderFinished() {
				m.t.Error("Expected call to senderMock.NewRecorder")
			}

			if !m.RegisterFinished() {
				m.t.Error("Expected call to senderMock.Register")
			}

			if !m.SendFinished() {
				m.t.Error("Expected call to senderMock.Send")
			}

			if !m.SendParcelFinished() {
				m.t.Error("Expected call to senderMock.SendParcel")
			}

			if !m.WriteTapeFinished() {
				m.t.Error("Expected call to senderMock.WriteTape")
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
func (m *senderMock) AllMocksCalled() bool {

	if !m.CreateParcelFinished() {
		return false
	}

	if !m.MustRegisterFinished() {
		return false
	}

	if !m.NewPlayerFinished() {
		return false
	}

	if !m.NewRecorderFinished() {
		return false
	}

	if !m.RegisterFinished() {
		return false
	}

	if !m.SendFinished() {
		return false
	}

	if !m.SendParcelFinished() {
		return false
	}

	if !m.WriteTapeFinished() {
		return false
	}

	return true
}
