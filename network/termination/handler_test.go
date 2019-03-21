package termination

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/core"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

type CommonTestSuite struct {
	suite.Suite

	mc           *minimock.Controller
	ctx          context.Context
	handler      *terminationHandler
	network      *testutils.NetworkMock
	pulseStorage *testutils.PulseStorageMock
}

func TestBasics(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

func (s *CommonTestSuite) BeforeTest(suiteName, testName string) {
	s.mc = minimock.NewController(s.T())
	s.ctx = inslogger.TestContext(s.T())
	s.network = testutils.NewNetworkMock(s.T())
	s.pulseStorage = testutils.NewPulseStorageMock(s.T())
	s.handler = &terminationHandler{Network: s.network, PulseStorage: s.pulseStorage}

}

func (s *CommonTestSuite) AfterTest(suiteName, testName string) {
	s.mc.Wait(10 * time.Second)
	s.mc.Finish()
}

func (s *CommonTestSuite) TestHandlerInitialState() {
	s.Equal(0, cap(s.handler.done))
	s.Equal(false, s.handler.terminating)
}

func (s *CommonTestSuite) HandlerIsTerminating() {
	s.Equal(true, s.handler.terminating)
	s.Equal(1, cap(s.handler.done))
}

func TestLeave(t *testing.T) {
	suite.Run(t, new(LeaveTestSuite))
}

type LeaveTestSuite struct {
	CommonTestSuite
}

func (s *LeaveTestSuite) TestLeaveNow() {
	s.network.LeaveMock.Expect(s.ctx, 0)
	s.handler.Leave(s.ctx, 0)

	s.HandlerIsTerminating()
}

func (s *LeaveTestSuite) TestLeaveEta() {
	mockPulseNumber := core.PulseNumber(2000000000)
	testPulse := &core.Pulse{PulseNumber: core.PulseNumber(mockPulseNumber)}
	pulseDelta := testPulse.NextPulseNumber - testPulse.PulseNumber
	leaveAfter := core.PulseNumber(5)

	//s.pulseStorage.CurrentMock.Return(testPulse, nil)
	s.pulseStorage.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		return testPulse, nil
	}
	s.network.LeaveMock.Expect(s.ctx, mockPulseNumber+leaveAfter*pulseDelta)
	s.handler.Leave(s.ctx, leaveAfter)

	s.HandlerIsTerminating()
}

func TestOnLeaveApproved(t *testing.T) {
	suite.Run(t, new(OnLeaveApprovedTestSuite))
}

type OnLeaveApprovedTestSuite struct {
	CommonTestSuite
}

func (s *OnLeaveApprovedTestSuite) TestBasicUsage() {
	s.handler.terminating = true
	s.handler.done = make(chan core.LeaveApproved, 1)

	s.handler.OnLeaveApproved(s.ctx)

	select {
	case <-s.handler.done:
		s.Equal(false, s.handler.terminating)
	case <-time.After(time.Second):
		s.Fail("done chanel doesn't close")
	}
}
