// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build networktest

package tests

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bootstrapSuite struct {
	testSuite
}

func (s *bootstrapSuite) Setup() {
	var err error
	s.pulsar, err = NewTestPulsar(reqTimeoutMs*10, pulseDelta*10)
	require.NoError(s.t, err)

	inslogger.FromContext(s.ctx).Info("SetupTest")

	for i := 0; i < s.bootstrapCount; i++ {
		role := insolar.StaticRoleVirtual
		if i == 0 {
			role = insolar.StaticRoleHeavyMaterial
		}

		s.bootstrapNodes = append(s.bootstrapNodes, s.newNetworkNodeWithRole(fmt.Sprintf("bootstrap_%d", i), role))
	}

	s.SetupNodesNetwork(s.bootstrapNodes)

	pulseReceivers := make([]string, 0)
	for _, node := range s.bootstrapNodes {
		pulseReceivers = append(pulseReceivers, node.host)
	}

	log.Info("Start test pulsar")
	err = s.pulsar.Start(s.ctx, pulseReceivers)
	require.NoError(s.t, err)
}

func (s *bootstrapSuite) stopBootstrapSuite() {
	inslogger.FromContext(s.ctx).Info("stopNetworkSuite")

	suiteLogger.Info("Stop bootstrap nodes")
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Stop(n.ctx)
		assert.NoError(s.t, err)
	}
}

func (s *bootstrapSuite) waitForConsensus(consensusCount int) {
	for i := 0; i < consensusCount; i++ {
		for _, n := range s.bootstrapNodes {
			<-n.consensusResult
		}
	}
}

func newBootstraptSuite(t *testing.T, bootstrapCount int) *bootstrapSuite {
	return &bootstrapSuite{
		testSuite: newTestSuite(t, bootstrapCount, 0),
	}
}

func startBootstrapSuite(t *testing.T) *bootstrapSuite {
	t.Skip("Skip until fix consensus bugs")

	s := newBootstraptSuite(t, 11)
	s.Setup()
	return s
}

func TestBootstrap(t *testing.T) {
	s := startBootstrapSuite(t)
	defer s.stopBootstrapSuite()

	s.StartNodesNetwork(s.bootstrapNodes)

	s.waitForConsensus(2)
	s.AssertActiveNodesCountDelta(0)

	s.waitForConsensus(1)
	s.AssertActiveNodesCountDelta(0)
	s.AssertWorkingNodesCountDelta(0)
}
