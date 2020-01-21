// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package tests

import (
	"math/rand"
	"time"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
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

func (dns *delayNetStrategy) GetLinkStrategy(hostAddress endpoints.Name) LinkStrategy {
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
