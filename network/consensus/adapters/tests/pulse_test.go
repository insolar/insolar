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
	"bytes"
	"context"
	"math/rand"
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
	th, _ := host.NewHost("127.0.0.1:2")
	pp := pulsenetwork.NewPulsePacketWithTrace(ctx, &pu, ph, th, 0)

	bs, _ := packet.SerializePacket(pp)
	rp, _ := packet.DeserializePacketRaw(bytes.NewReader(bs))

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
