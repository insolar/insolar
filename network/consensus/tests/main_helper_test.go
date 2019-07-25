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

package tests

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func newEmuNetworkBuilder(ctx context.Context, netStrategy NetStrategy,
	roundStrategyFactory core.RoundStrategyFactory) emuNetworkBuilder {

	r := emuNetworkBuilder{}
	r.network = NewEmuNetwork(netStrategy, ctx)
	r.config = NewEmuLocalConfig(ctx)
	r.primingCloudStateHash = EmuPrimingHash
	r.strategyFactory = roundStrategyFactory
	return r
}

type emuNetworkBuilder struct {
	network               *EmuNetwork
	config                api.LocalNodeConfiguration
	primingCloudStateHash cryptkit.DigestHolder
	strategyFactory       core.RoundStrategyFactory
}

func (p *emuNetworkBuilder) StartNetwork(ctx context.Context) {
	p.network.Start(ctx)
}

func (p *emuNetworkBuilder) StartPulsar(pulseCount int, pulseDelta uint16, pulsarAddr endpoints.Name,
	nodes []profiles.StaticProfile) {

	attempts := 4 + len(nodes)/10

	senderChan := make(chan interface{})
	go func() {
		for {
			payload, ok := <-senderChan
			if !ok {
				return
			}
			for i := 0; i < attempts; i++ {
				sendTo := nodes[rand.Intn(len(nodes))].GetDefaultEndpoint().GetNameAddress()
				p.network.SendToHost(sendTo, payload, pulsarAddr)
			}
		}
	}()

	go CreateGenerator(pulseCount, pulseDelta, senderChan)
}

const fmtNodeName = "%s%04d"

func (p *emuNetworkBuilder) connectEmuNode(nodes []profiles.StaticProfile, selfIndex int) {
	p._connectEmuNode(nodes, selfIndex, false)
}

func (p *emuNetworkBuilder) _connectEmuNode(nodes []profiles.StaticProfile, selfIndex int, asJoiner bool) {

	controlFeeder := &EmuControlFeeder{}
	candidateFeeder := &core.SequentialCandidateFeeder{}
	switch {
	case !asJoiner: // && selfIndex%3 == 1
		for i := 5000; i < 8000; i += 1000 {
			introID := i + selfIndex
			intro := NewEmuNodeIntroByName(introID, fmt.Sprintf(fmtNodeName, "V", introID))
			candidateFeeder.AddJoinCandidate(intro)
			p._connectEmuNode([]profiles.StaticProfile{intro}, 0, true)
		}

		//case selfIndex%5 == 2:
		//	controlFeeder.leaveReason = uint32(selfIndex) // simulate leave
	}

	chronicles := NewEmuChronicles(nodes, selfIndex, asJoiner, p.primingCloudStateHash)
	self := nodes[selfIndex]
	node := NewConsensusHost(self.GetDefaultEndpoint().GetNameAddress())
	node.ConnectTo(chronicles, p.network, p.strategyFactory, candidateFeeder, controlFeeder, p.config)
}

func generateNameList(countNeutral, countHeavy, countLight, countVirtual int) []string {
	r := make([]string, 0, countNeutral+countHeavy+countLight+countVirtual)

	r = appendNameList(len(r), r, fmtNodeName, "N", countNeutral)
	r = appendNameList(len(r), r, fmtNodeName, "H", countHeavy)
	r = appendNameList(len(r), r, fmtNodeName, "L", countLight)
	r = appendNameList(len(r), r, fmtNodeName, "V", countVirtual)

	return r
}

func appendNameList(baseNum int, r []string, fmtS, pfx string, count int) []string {
	for i := 0; i < count; i++ {
		r = append(r, fmt.Sprintf(fmtS, pfx, baseNum+i))
	}
	return r
}
