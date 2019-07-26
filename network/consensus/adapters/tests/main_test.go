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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/stretchr/testify/require"
)

const (
	defaultPulseDelta   = 5
	defaultLogLevel     = insolar.DebugLevel
	defaultTestDuration = defaultPulseDelta * time.Second * 10
)

var strategy = NewDelayNetStrategy(DelayStrategyConf{
	MinDelay:         10 * time.Millisecond,
	MaxDelay:         30 * time.Millisecond,
	Variance:         0.2,
	SpikeProbability: 0.1,
})

func TestConsensusJoin(t *testing.T) {
	startedAt := time.Now()
	ctx := initLogger(defaultLogLevel)

	nodeIdentities := generateNodeIdentities(0, 1, 8, 8)
	nodeInfos := generateNodeInfos(nodeIdentities)
	nodes, discoveryNodes := nodesFromInfo(nodeInfos)

	joinIdentities := generateNodeIdentities(0, 0, 4, 4)
	joinInfos := generateNodeInfos(joinIdentities)
	joiners, _ := nodesFromInfo(joinInfos)

	controllers, pulseHandlers, _, _, _, err := initNodes(ctx, consensus.ReadyNetwork, nodes, discoveryNodes, strategy, nodeInfos)
	require.NoError(t, err)

	_, _, _, _, joinerProfiles, err := initNodes(ctx, consensus.Joiner, joiners, discoveryNodes, strategy, joinInfos)
	require.NoError(t, err)

	fmt.Println("===", len(nodes), "=================================================")

	pulsar := NewPulsar(defaultPulseDelta, pulseHandlers)
	go func() {
		for {
			pulsar.Pulse(ctx, 4+len(nodes)/10)
		}
	}()

	once := sync.Once{}

	for {
		fmt.Println("===", time.Since(startedAt), "=================================================")
		time.Sleep(time.Second)
		if time.Since(startedAt) > defaultTestDuration {
			return
		}

		if time.Since(startedAt) > 1*time.Second {
			once.Do(func() {
				type candidate struct {
					profiles.StaticProfile
					profiles.StaticProfileExtension
				}

				for i, joiner := range joinerProfiles {
					controllers[i].AddJoinCandidate(candidate{
						joiner,
						joiner.GetExtension(),
					})
				}
			})
		}
	}
}

func TestConsensusLeave(t *testing.T) {
	startedAt := time.Now()
	ctx := initLogger(defaultLogLevel)

	nodeIdentities := generateNodeIdentities(0, 1, 3, 5)
	nodeInfos := generateNodeInfos(nodeIdentities)
	nodes, discoveryNodes := nodesFromInfo(nodeInfos)

	controllers, pulseHandlers, transports, contexts, _, err := initNodes(ctx, consensus.ReadyNetwork, nodes, discoveryNodes, strategy, nodeInfos)
	require.NoError(t, err)

	fmt.Println("===", len(nodes), "=================================================")

	pulsar := NewPulsar(defaultPulseDelta, pulseHandlers)
	go func() {
		for {
			pulsar.Pulse(ctx, 4+len(nodes)/10)
		}
	}()

	once := sync.Once{}

	for {
		fmt.Println("===", time.Since(startedAt), "=================================================")
		time.Sleep(time.Second)
		if time.Since(startedAt) > defaultTestDuration {
			return
		}

		nodeIdx := 0

		if time.Since(startedAt) > 1*time.Second {
			once.Do(func() {
				<-controllers[nodeIdx].Leave(0)
				err := transports[nodeIdx].Stop(contexts[nodeIdx])
				require.NoError(t, err)
				controllers[nodeIdx].Abort()
			})
		}
	}
}

func TestConsensusDrop(t *testing.T) {
	startedAt := time.Now()
	ctx := initLogger(defaultLogLevel)

	nodeIdentities := generateNodeIdentities(0, 1, 3, 5)
	nodeInfos := generateNodeInfos(nodeIdentities)
	nodes, discoveryNodes := nodesFromInfo(nodeInfos)

	_, pulseHandlers, transports, contexts, _, err := initNodes(ctx, consensus.ReadyNetwork, nodes, discoveryNodes, strategy, nodeInfos)
	require.NoError(t, err)

	fmt.Println("===", len(nodes), "=================================================")

	pulsar := NewPulsar(defaultPulseDelta, pulseHandlers)
	go func() {
		for {
			pulsar.Pulse(ctx, 4+len(nodes)/10)
		}
	}()

	once := sync.Once{}

	for {
		fmt.Println("===", time.Since(startedAt), "=================================================")
		time.Sleep(time.Second)
		if time.Since(startedAt) > defaultTestDuration {
			return
		}

		nodeIdx := 0

		if time.Since(startedAt) > 1*time.Second {
			once.Do(func() {
				err := transports[nodeIdx].Stop(contexts[nodeIdx])
				require.NoError(t, err)
			})
		}
	}
}

func TestConsensusAll(t *testing.T) {
	startedAt := time.Now()
	ctx := initLogger(defaultLogLevel)

	nodeIdentities := generateNodeIdentities(0, 1, 3, 5)
	nodeInfos := generateNodeInfos(nodeIdentities)
	nodes, discoveryNodes := nodesFromInfo(nodeInfos)

	joinIdentities := generateNodeIdentities(0, 0, 2, 2)
	joinInfos := generateNodeInfos(joinIdentities)
	joiners, _ := nodesFromInfo(joinInfos)

	controllers, pulseHandlers, transports, contexts, _, err := initNodes(ctx, consensus.ReadyNetwork, nodes, discoveryNodes, strategy, nodeInfos)
	require.NoError(t, err)

	_, _, _, _, joinerProfiles, err := initNodes(ctx, consensus.Joiner, joiners, discoveryNodes, strategy, joinInfos)
	require.NoError(t, err)

	fmt.Println("===", len(nodes), "=================================================")

	pulsar := NewPulsar(defaultPulseDelta, pulseHandlers)
	go func() {
		for {
			pulsar.Pulse(ctx, 4+len(nodes)/10)
		}
	}()

	once1 := sync.Once{}
	once2 := sync.Once{}
	once3 := sync.Once{}

	for {
		fmt.Println("===", time.Since(startedAt), "=================================================")
		time.Sleep(time.Second)
		if time.Since(startedAt) > defaultTestDuration {
			return
		}

		if time.Since(startedAt) > 1*time.Second {
			once1.Do(func() {
				nodeIdx := 6

				<-controllers[nodeIdx].Leave(0)
				err := transports[nodeIdx].Stop(contexts[nodeIdx])
				require.NoError(t, err)
				controllers[nodeIdx].Abort()
			})

			once2.Do(func() {
				nodeIdx := 7

				err := transports[nodeIdx].Stop(contexts[nodeIdx])
				require.NoError(t, err)
			})

			once3.Do(func() {
				type candidate struct {
					profiles.StaticProfile
					profiles.StaticProfileExtension
				}

				for i, joiner := range joinerProfiles {
					controllers[i].AddJoinCandidate(candidate{
						joiner,
						joiner.GetExtension(),
					})
				}
			})
		}
	}
}
