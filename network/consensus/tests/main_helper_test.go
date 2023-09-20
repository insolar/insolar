package tests

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"

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
	candidateFeeder := &coreapi.SequentialCandidateFeeder{}
	ephemeralFeeder := &EmuEphemeralFeeder{}

	switch {
	case !asJoiner && selfIndex == 1:
		for i := 5000; i < 8000; i += 1000 {
			introID := i + selfIndex
			intro := NewEmuNodeIntroByName(introID, fmt.Sprintf(fmtNodeName, "V", introID))
			candidateFeeder.AddJoinCandidate(intro)
			p._connectEmuNode([]profiles.StaticProfile{intro}, 0, true)
		}

		// case selfIndex%5 == 2:
		//	controlFeeder.leaveReason = uint32(selfIndex) // simulate leave
	}

	chronicles := NewEmuChronicles(nodes, selfIndex, asJoiner, p.primingCloudStateHash)
	self := nodes[selfIndex]
	node := NewConsensusHost(self.GetDefaultEndpoint().GetNameAddress())
	node.ConnectTo(chronicles, p.network, p.strategyFactory, candidateFeeder, controlFeeder, ephemeralFeeder, p.config)
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
