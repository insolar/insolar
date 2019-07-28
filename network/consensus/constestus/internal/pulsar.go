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
	"bytes"
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/constestus/cloud"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/pulsenetwork"
)

type Pulsar struct {
	mu    *sync.Mutex
	pulse *pulse.Data
}

func NewPulsar(config cloud.Network) Pulsar {
	var firstPulse *pulse.Data

	if config.Consensus.EphemeralPulses {
		firstPulse = pulse.NewFirstEphemeralData()
	} else {
		firstPulse = pulse.NewFirstPulsarData(config.Consensus.PulseDelta, randomEntropy())
	}

	return Pulsar{
		mu:    &sync.Mutex{},
		pulse: firstPulse,
	}
}

func (p Pulsar) Pulse(ctx context.Context, activeNodes ActiveNodes, attempts int) error {
	p.mu.Lock()
	defer time.AfterFunc(time.Duration(p.pulse.NextPulseDelta)*time.Second, func() {
		p.mu.Unlock()
	})

	p.pulse = p.pulse.CreateNextPulse(randomEntropy)

	newPulse := adapters.NewPulse(*p.pulse)
	pulsePacket, err := getReceivedPulsarPacket(ctx, newPulse)
	if err != nil {
		return err
	}

	for i := 0; i < attempts; i++ {
		node := activeNodes[rand.Intn(len(activeNodes))]
		pulseHandler := node.components.pulseHandler

		go pulseHandler.HandlePulse(ctx, newPulse, pulsePacket)
	}

	return nil
}

func getReceivedPulsarPacket(ctx context.Context, pu insolar.Pulse) (*packet.ReceivedPacket, error) {
	pulsarHost, err := host.NewHost("127.0.0.1:1")
	if err != nil {
		return nil, err
	}

	targetHost, err := host.NewHost("127.0.0.1:2")
	if err != nil {
		return nil, err
	}

	pulsePacket := pulsenetwork.NewPulsePacket(ctx, &pu, pulsarHost, targetHost, 0)
	pulsePacketData, err := packet.SerializePacket(pulsePacket)
	if err != nil {
		return nil, err
	}

	rp, err := packet.DeserializePacketRaw(bytes.NewReader(pulsePacketData))
	if err != nil {
		return nil, err
	}

	return rp, nil
}

func randomEntropy() longbits.Bits256 {
	v := longbits.Bits256{}
	_, _ = rand.Read(v[:])
	return v
}
