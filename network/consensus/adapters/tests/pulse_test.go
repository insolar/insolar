// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package tests

import (
	"bytes"
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/pulse"
)

const (
	initialPulse = 100000
)

type Pulsar struct {
	pulseDelta    uint16
	pulseNumber   pulse.Number
	pulseHandlers []network.PulseHandler

	mu *sync.Mutex
}

func NewPulsar(pulseDelta uint16, pulseHandlers []network.PulseHandler) Pulsar {
	return Pulsar{
		pulseDelta:    pulseDelta,
		pulseNumber:   initialPulse,
		pulseHandlers: pulseHandlers,
		mu:            &sync.Mutex{},
	}
}

func (p *Pulsar) Pulse(ctx context.Context, attempts int) {
	p.mu.Lock()
	defer time.AfterFunc(time.Duration(p.pulseDelta)*time.Second, func() {
		p.mu.Unlock()
	})

	prevDelta := p.pulseDelta
	if p.pulseNumber == initialPulse {
		prevDelta = 0
	}

	data := pulse.NewPulsarData(p.pulseNumber, p.pulseDelta, prevDelta, randBits256())
	p.pulseNumber += pulse.Number(p.pulseDelta)

	pu := adapters.NewPulse(data)
	ph, _ := host.NewHost("127.0.0.1:1")
	th, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2")
	pp := pulsenetwork.NewPulsePacketWithTrace(ctx, &pu, ph, th, 0)

	bs, _ := packet.SerializePacket(pp)
	rp, _, _ := packet.DeserializePacketRaw(bytes.NewReader(bs))

	go func() {
		for i := 0; i < attempts; i++ {
			handler := p.pulseHandlers[rand.Intn(len(p.pulseHandlers))]
			go handler.HandlePulse(ctx, pu, rp)
		}
	}()
}

func randBits256() longbits.Bits256 {
	v := longbits.Bits256{}
	_, _ = rand.Read(v[:])
	return v
}
