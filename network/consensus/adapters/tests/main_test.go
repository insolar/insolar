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

// +build never_run

package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultLogLevel       = insolar.DebugLevel
	defaultPulseDelta     = 2
	defaultTestDuration   = defaultPulseDelta * time.Second * 10
	defaultStartCaseAfter = 1 * time.Second
)

var strategy = NewDelayNetStrategy(DelayStrategyConf{
	MinDelay:         10 * time.Millisecond,
	MaxDelay:         30 * time.Millisecond,
	Variance:         0.2,
	SpikeProbability: 0.1,
})

var ctx = initLogger(defaultLogLevel)

func TestConsensusJoin(t *testing.T) {
	nodes, err := generateNodes(0, 1, 3, 5, nil)
	require.NoError(t, err)

	joiners, err := generateNodes(0, 0, 6, 1, nodes.discoveryNodes)
	require.NoError(t, err)

	ns, err := initNodes(ctx, consensus.ReadyNetwork, *nodes, strategy)
	require.NoError(t, err)

	js, err := initNodes(ctx, consensus.Joiner, *joiners, strategy)
	require.NoError(t, err)

	initPulsar(ctx, defaultPulseDelta, *ns)

	testCase(defaultTestDuration, defaultStartCaseAfter, func() {
		for i, joiner := range js.staticProfiles {
			err := ns.controllers[i].AddJoinCandidate(candidate{
				joiner,
				joiner.GetExtension(),
			})

			require.NoError(t, err)
		}
	})

	// require.Len(t, ns.nodeKeepers[0].GetAccessor().GetActiveNodes(), len(nodes.nodes)+len(joiners.nodes))
}

func TestConsensusLeave(t *testing.T) {
	nodes, err := generateNodes(0, 1, 3, 5, nil)
	require.NoError(t, err)

	ns, err := initNodes(ctx, consensus.ReadyNetwork, *nodes, strategy)
	require.NoError(t, err)

	initPulsar(ctx, defaultPulseDelta, *ns)

	testCase(defaultTestDuration, defaultStartCaseAfter, func() {
		nodeIdx := 1

		<-ns.controllers[nodeIdx].Leave(0)
		err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
		require.NoError(t, err)
		ns.controllers[nodeIdx].Abort()
	})

	// require.Len(t, ns.nodeKeepers[0].GetAccessor().GetActiveNodes(), len(nodes.nodes)-1)
}

func TestConsensusDrop(t *testing.T) {
	nodes, err := generateNodes(0, 1, 3, 5, nil)
	require.NoError(t, err)

	ns, err := initNodes(ctx, consensus.ReadyNetwork, *nodes, strategy)
	require.NoError(t, err)

	initPulsar(ctx, defaultPulseDelta, *ns)

	testCase(defaultTestDuration, defaultStartCaseAfter, func() {
		nodeIdx := 1

		err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
		require.NoError(t, err)
	})

	// require.Len(t, ns.nodeKeepers[0].GetAccessor().GetActiveNodes(), len(nodes.nodes)-1)
}

func TestConsensusJoinLeave(t *testing.T) {
	t.Skip("Until phase 4 ready")

	nodes, err := generateNodes(0, 1, 3, 5, nil)
	require.NoError(t, err)

	joiners, err := generateNodes(0, 0, 0, 1, nodes.discoveryNodes)
	require.NoError(t, err)

	ns, err := initNodes(ctx, consensus.ReadyNetwork, *nodes, strategy)
	require.NoError(t, err)

	js, err := initNodes(ctx, consensus.Joiner, *joiners, strategy)
	require.NoError(t, err)

	initPulsar(ctx, defaultPulseDelta, *ns)

	testCase(defaultTestDuration, defaultStartCaseAfter, func() {
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			nodeIdx := len(joiners.nodes) + 1

			<-ns.controllers[nodeIdx].Leave(0)
			err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
			assert.NoError(t, err)
			ns.controllers[nodeIdx].Abort()

			wg.Done()
		}()

		go func() {
			for i, joiner := range js.staticProfiles {
				err := ns.controllers[i].AddJoinCandidate(candidate{
					joiner,
					joiner.GetExtension(),
				})

				require.NoError(t, err)
			}

			wg.Done()
		}()

		wg.Wait()
	})

	// require.Len(t, ns.nodeKeepers[0].GetAccessor().GetActiveNodes(), len(nodes.nodes)+len(joiners.nodes)-1)
}

func TestConsensusJoinDrop(t *testing.T) {
	t.Skip("Until phase 4 ready")

	nodes, err := generateNodes(0, 1, 3, 5, nil)
	require.NoError(t, err)

	joiners, err := generateNodes(0, 0, 0, 1, nodes.discoveryNodes)
	require.NoError(t, err)

	ns, err := initNodes(ctx, consensus.ReadyNetwork, *nodes, strategy)
	require.NoError(t, err)

	js, err := initNodes(ctx, consensus.Joiner, *joiners, strategy)
	require.NoError(t, err)

	initPulsar(ctx, defaultPulseDelta, *ns)

	testCase(defaultTestDuration, defaultStartCaseAfter, func() {
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			nodeIdx := len(joiners.nodes) + 1

			err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
			assert.NoError(t, err)

			wg.Done()
		}()

		go func() {
			for i, joiner := range js.staticProfiles {
				err := ns.controllers[i].AddJoinCandidate(candidate{
					joiner,
					joiner.GetExtension(),
				})

				require.NoError(t, err)
			}

			wg.Done()
		}()

		wg.Wait()
	})

	// require.Len(t, ns.nodeKeepers[0].GetAccessor().GetActiveNodes(), len(nodes.nodes)+len(joiners.nodes)-1)
}

func TestConsensusDropLeave(t *testing.T) {
	nodes, err := generateNodes(0, 1, 3, 5, nil)
	require.NoError(t, err)

	ns, err := initNodes(ctx, consensus.ReadyNetwork, *nodes, strategy)
	require.NoError(t, err)

	initPulsar(ctx, defaultPulseDelta, *ns)

	testCase(defaultTestDuration, defaultStartCaseAfter, func() {
		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			nodeIdx := 6

			<-ns.controllers[nodeIdx].Leave(0)
			err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
			assert.NoError(t, err)
			ns.controllers[nodeIdx].Abort()

			wg.Done()
		}()

		go func() {
			nodeIdx := 7

			err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
			assert.NoError(t, err)

			wg.Done()
		}()

		wg.Wait()
	})

	// require.Len(t, ns.nodeKeepers[0].GetAccessor().GetActiveNodes(), len(nodes.nodes)-2)
}

func TestConsensusAll(t *testing.T) {
	t.Skip("Until phase 4 ready")

	nodes, err := generateNodes(0, 1, 3, 5, nil)
	require.NoError(t, err)

	joiners, err := generateNodes(0, 0, 1, 1, nodes.discoveryNodes)
	require.NoError(t, err)

	ns, err := initNodes(ctx, consensus.ReadyNetwork, *nodes, strategy)
	require.NoError(t, err)

	js, err := initNodes(ctx, consensus.Joiner, *joiners, strategy)
	require.NoError(t, err)

	initPulsar(ctx, defaultPulseDelta, *ns)

	testCase(defaultTestDuration, defaultStartCaseAfter, func() {
		wg := &sync.WaitGroup{}
		wg.Add(3)

		go func() {
			nodeIdx := len(joiners.nodes) + 1

			<-ns.controllers[nodeIdx].Leave(0)
			err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
			assert.NoError(t, err)
			ns.controllers[nodeIdx].Abort()

			wg.Done()
		}()

		go func() {
			nodeIdx := len(joiners.nodes) + 2

			err := ns.transports[nodeIdx].Stop(ns.contexts[nodeIdx])
			assert.NoError(t, err)

			wg.Done()
		}()

		go func() {
			for i, joiner := range js.staticProfiles {
				err := ns.controllers[i].AddJoinCandidate(candidate{
					joiner,
					joiner.GetExtension(),
				})

				require.NoError(t, err)
			}

			wg.Done()
		}()

		wg.Wait()
	})

	// require.Len(t, ns.nodeKeepers[0].GetAccessor().GetActiveNodes(), len(nodes.nodes)+len(joiners.nodes)-2)
}
