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

package internal

import (
	"context"
	"math/rand"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/constestus/cloud"
	"github.com/insolar/insolar/network/transport"
)

type NetworkStrategy interface {
	GetLink(datagramTransport transport.DatagramTransport) transport.DatagramTransport
}

type delayNetStrategy struct {
	conf cloud.Delays
}

func NewDelayNetStrategy(conf cloud.Delays) NetworkStrategy {
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

func (dns *delayNetStrategy) GetLink(datagramTransport transport.DatagramTransport) transport.DatagramTransport {
	return newDelayLinkStrategy(
		datagramTransport,
		dns.getDelay(),
		dns.conf.SpikeDelay,
		dns.conf.Variance,
		dns.conf.SpikeProbability,
	)
}

type delayLinkStrategy struct {
	transport.DatagramTransport

	normalDelay time.Duration
	spikeDelay  time.Duration

	normalDelayMaxVariance int
	spikeDelayMaxVariance  int

	spikeProbability float32
}

func newDelayLinkStrategy(transport transport.DatagramTransport, normalDelay, spikeDelay time.Duration, variance, spikeProbability float32) *delayLinkStrategy {
	return &delayLinkStrategy{
		DatagramTransport: transport,

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

func (dls *delayLinkStrategy) delay(f func()) {
	if delay := dls.calculateDelay(); delay > 0 {
		time.AfterFunc(delay, f)
	} else {
		f()
	}
}

func (dls *delayLinkStrategy) SendDatagram(ctx context.Context, address string, data []byte) error {
	dls.delay(func() {
		if err := dls.DatagramTransport.SendDatagram(ctx, address, data); err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
	})
	return nil
}
