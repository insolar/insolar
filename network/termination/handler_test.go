//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package termination

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

type CommonTestSuite struct {
	suite.Suite

	mc            *minimock.Controller
	ctx           context.Context
	handler       *terminationHandler
	leaver        *testutils.LeaverMock
	pulseAccessor *pulse.AccessorMock
}

func TestBasics(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

func (s *CommonTestSuite) BeforeTest(suiteName, testName string) {
	s.mc = minimock.NewController(s.T())
	s.ctx = inslogger.TestContext(s.T())
	s.leaver = testutils.NewLeaverMock(s.T())
	s.pulseAccessor = pulse.NewAccessorMock(s.T())
	s.handler = &terminationHandler{Leaver: s.leaver, PulseAccessor: s.pulseAccessor}

}

func (s *CommonTestSuite) AfterTest(suiteName, testName string) {
	s.mc.Wait(time.Minute)
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
	s.leaver.LeaveMock.Expect(s.ctx, 0)
	s.handler.leave(s.ctx, 0)

	s.HandlerIsTerminating()
}

func (s *LeaveTestSuite) TestLeaveEta() {
	mockPulseNumber := insolar.PulseNumber(2000000000)
	testPulse := &insolar.Pulse{PulseNumber: mockPulseNumber}
	pulseDelta := testPulse.NextPulseNumber - testPulse.PulseNumber
	leaveAfter := insolar.PulseNumber(5)

	s.pulseAccessor.LatestMock.Return(*testPulse, nil)
	s.leaver.LeaveMock.Expect(s.ctx, mockPulseNumber+leaveAfter*pulseDelta)
	s.handler.leave(s.ctx, leaveAfter)

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
	s.handler.done = make(chan insolar.LeaveApproved, 1)

	s.handler.OnLeaveApproved(s.ctx)

	select {
	case <-s.handler.done:
		s.Equal(false, s.handler.terminating)
	case <-time.After(time.Second):
		s.Fail("done chanel doesn't close")
	}
}

func TestAbort(t *testing.T) {
	suite.Run(t, new(AbortTestSuite))
}

type AbortTestSuite struct {
	CommonTestSuite
}

func (s *AbortTestSuite) TestBasicUsage() {
	defer func() {
		if r := recover(); r == nil {
			s.Fail("did not catch panic")
		}
	}()

	s.handler.Abort("abort")
}
