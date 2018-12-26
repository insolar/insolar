/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package servicenetwork

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

func (s *testSuite) TestNetworkConsensus3Times() {
	s.waitForConsensus(3)
}

func (s *testSuite) TestNodeConnect() {
	s.preInitNode(s.testNode, Disable)

	s.InitTestNode()
	s.StartTestNode()
	defer s.StopTestNode()

	s.waitForConsensus(1)

	activeNodes := s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount(), len(activeNodes))

	//log.Warn("-------=-=-=-=-=-=-================")
	err := s.bootstrapNodes[0].serviceNetwork.NodeKeeper.MoveSyncToActive()
	s.NoError(err)

	//s.waitForConsensus(1)
	//s.waitForConsensus(1)

	activeNodes = s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount()+1, len(activeNodes))

	// teardown
	// <-time.After(time.Second * 5)
}

func (s *testSuite) TestNodeLeave() {
	s.T().Skip("tmp 123")
	phasesResult := make(chan error)

	s.preInitNode(s.testNode, Disable)

	s.InitTestNode()
	s.bootstrapNodes[0].serviceNetwork.PhaseManager = &phaseManagerWrapper{original: s.bootstrapNodes[0].serviceNetwork.PhaseManager, result: phasesResult}

	s.StartTestNode()

	res := <-phasesResult
	s.NoError(res)
	activeNodes := s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount(), len(activeNodes))
	err := s.bootstrapNodes[0].serviceNetwork.NodeKeeper.MoveSyncToActive()
	s.NoError(err)
	activeNodes = s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount()+1, len(activeNodes))

	s.testNode.serviceNetwork.GracefulStop(context.Background())

	res = <-phasesResult
	s.NoError(res)
	activeNodes = s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount()+1, len(activeNodes))
	err = s.bootstrapNodes[0].serviceNetwork.NodeKeeper.MoveSyncToActive()
	s.NoError(err)
	activeNodes = s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount(), len(activeNodes))

	// teardown
	<-time.After(time.Second * 3)
	s.StopTestNode()
}

func TestServiceNetworkIntegration(t *testing.T) {
	s := NewTestSuite(1, 0)
	suite.Run(t, s)
	defer func() {
		p := recover()
		fmt.Println(p)
	}()

}

func TestServiceNetworkManyBootstraps(t *testing.T) {
	s := NewTestSuite(15, 0)
	suite.Run(t, s)
}

/*
func TestServiceNetworkManyNodes(t *testing.T) {
	t.Skip("tmp 123")

	s := NewTestSuite(3, 20)
	suite.Run(t, s)
}
*/
// Full timeout test
type FullTimeoutPhaseManager struct {
}

func (ftpm *FullTimeoutPhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	return nil
}

func (s *testSuite) TestFullTimeOut() {
	s.T().Skip("will be available after phase result fix !")
	phasesResult := make(chan error)

	s.preInitNode(s.testNode, Full)

	s.InitTestNode()
	s.StartTestNode()
	res := <-phasesResult
	s.NoError(res)
	activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(1, len(activeNodes))
	// teardown
	<-time.After(time.Second * 5)
	s.StopTestNode()
}

// Partial timeout

func (s *testSuite) TestPartialTimeOut() {
	s.T().Skip("fix me")
	phasesResult := make(chan error)

	s.preInitNode(s.testNode, Partial)

	s.InitTestNode()
	s.StartTestNode()
	res := <-phasesResult
	s.NoError(res)
	// activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	// s.Equal(1, len(activeNodes))	// TODO: do test check
	// teardown
	<-time.After(time.Second * 5)
	s.StopTestNode()
}

type PartialTimeoutPhaseManager struct {
	FirstPhase  *phases.FirstPhase
	SecondPhase *phases.SecondPhase
	ThirdPhase  *phases.ThirdPhase
	Keeper      network.NodeKeeper `inject:""`
}

func (ftpm *PartialTimeoutPhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	var err error

	pulseDuration, err := getPulseDuration(pulse)
	if err != nil {
		return errors.Wrap(err, "[ OnPulse ] Failed to get pulse duration")
	}

	var tctx context.Context
	var cancel context.CancelFunc

	tctx, cancel = contextTimeout(ctx, *pulseDuration, 0.2)
	defer cancel()

	firstPhaseState, err := ftpm.FirstPhase.Execute(tctx, pulse)

	if err != nil {
		return errors.Wrap(err, "[ TestCase.OnPulse ] failed to execute a phase")
	}

	tctx, cancel = contextTimeout(ctx, *pulseDuration, 0.2)
	defer cancel()

	secondPhaseState, err := ftpm.SecondPhase.Execute(tctx, firstPhaseState)
	checkError(err)

	_, err = ftpm.ThirdPhase.Execute(ctx, secondPhaseState)
	checkError(err)

	return nil
}

func contextTimeout(ctx context.Context, duration time.Duration, k float64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(k * float64(duration))
	timedCtx, cancelFund := context.WithTimeout(ctx, timeout)
	return timedCtx, cancelFund
}

func getPulseDuration(pulse *core.Pulse) (*time.Duration, error) {
	duration := time.Duration(pulse.PulseNumber-pulse.PrevPulseNumber) * time.Second
	return &duration, nil
}

func checkError(err error) {
	if err != nil {
		log.Error(err)
	}
}
