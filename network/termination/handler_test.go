package termination

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

type CommonTestSuite struct {
	suite.Suite

	mc      *minimock.Controller
	ctx     context.Context
	handler *terminationHandler
	network *testutils.NetworkMock
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

func (s *CommonTestSuite) BeforeTest(suiteName, testName string) {
	s.mc = minimock.NewController(s.T())
	s.ctx = inslogger.TestContext(s.T())
	s.network = testutils.NewNetworkMock(s.T())
	s.handler = &terminationHandler{Network: s.network}

}

func (s *CommonTestSuite) AfterTest(suiteName, testName string) {
	s.mc.Wait(10 * time.Second)
	s.mc.Finish()
}

func (s *CommonTestSuite) TestLeave() {
	s.network.LeaveMock.Expect(s.ctx, 0)
	s.handler.Leave(s.ctx, 0)
	s.Equal(s.handler.terminating, true)
}
