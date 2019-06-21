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
	"math/rand"
	"time"

	"github.com/insolar/insolar/network/consensus/common"
)

type DelayStrategyConf struct {
	MinDelay         time.Duration
	MaxDelay         time.Duration
	SpikeDelay       time.Duration
	Variance         float32
	SpikeProbability float32
}

type delayNetStrategy struct {
	conf DelayStrategyConf
}

func NewDelayNetStrategy(conf DelayStrategyConf) NetStrategy {
	if conf.MinDelay > conf.MaxDelay {
		panic("MinDelay must <= MaxDelay")
	}

	if conf.Variance < 0 {
		panic("Variance must be in [0, Inf)")
	}

	if conf.SpikeProbability < 0 || conf.SpikeProbability > 1 {
		panic("SpikeProbability must be in [0, 1]")
	}

	if conf.SpikeDelay == 0 {
		conf.SpikeDelay = conf.MaxDelay
	}

	return &delayNetStrategy{
		conf: conf,
	}
}

func (dns *delayNetStrategy) getDelay() time.Duration {
	if dns.conf.MaxDelay == dns.conf.MinDelay {
		return dns.conf.MaxDelay
	}

	if dns.conf.MaxDelay > 0 {
		randomDelay := rand.Intn(int(dns.conf.MaxDelay)-int(dns.conf.MinDelay)) + int(dns.conf.MinDelay)
		return time.Duration(randomDelay)
	}

	return 0
}

func (dns *delayNetStrategy) GetLinkStrategy(hostAddress common.HostAddress) LinkStrategy {
	return newDelayLinkStrategy(
		dns.getDelay(),
		dns.conf.SpikeDelay,
		dns.conf.Variance,
		dns.conf.SpikeProbability,
	)
}

type delayLinkStrategy struct {
	normalDelay time.Duration
	spikeDelay  time.Duration

	normalDelayMaxVariance int
	spikeDelayMaxVariance  int

	spikeProbability float32
}

func newDelayLinkStrategy(normalDelay, spikeDelay time.Duration, variance, spikeProbability float32) *delayLinkStrategy {
	return &delayLinkStrategy{
		normalDelay: normalDelay,
		spikeDelay:  spikeDelay,

		normalDelayMaxVariance: int(float32(normalDelay) * variance),
		spikeDelayMaxVariance:  int(float32(spikeDelay) * variance),

		spikeProbability: spikeProbability,
	}
}

func (dls *delayLinkStrategy) calculateDelay() time.Duration {
	var (
		initialDelay     time.Duration
		delayMaxVariance int
	)

	if rand.Float32() <= dls.spikeProbability {
		initialDelay = dls.spikeDelay
		delayMaxVariance = dls.spikeDelayMaxVariance
	} else {
		initialDelay = dls.normalDelay
		delayMaxVariance = dls.normalDelayMaxVariance
	}

	delay := initialDelay

	if delayMaxVariance > 0 {
		delay += time.Duration(rand.Intn(delayMaxVariance))
	}

	return delay
}

func (dls *delayLinkStrategy) delay(tp string, packet *Packet, out PacketFunc) {
	if delay := dls.calculateDelay(); delay > 0 {
		// fmt.Printf(">>>> %s packet delay %v: %v\n", tp, delay, packet)
		time.AfterFunc(delay, func() {
			out(packet)
		})
	} else {
		out(packet)
	}
}

func (dls *delayLinkStrategy) BeforeSend(packet *Packet, out PacketFunc) {
	dls.delay("OUT", packet, out)
}

func (dls *delayLinkStrategy) BeforeReceive(packet *Packet, out PacketFunc) {
	// dls.delay(" IN", packet, out)
	out(packet)
}
