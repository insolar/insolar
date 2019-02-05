/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package pinger

import (
	"context"
	"time"

	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

// Pinger is a light and stateless component that can ping remote host to receive its NodeID
type Pinger struct {
	transport network.InternalTransport
}

// PingWithTimeout ping remote host with timeout
func (p *Pinger) Ping(ctx context.Context, address string, timeout time.Duration) (*host.Host, error) {
	ctx, span := instracer.StartSpan(ctx, "Pinger.Ping")
	defer span.End()
	request := p.transport.NewRequestBuilder().Type(types.Ping).Build()
	h, err := host.NewHost(address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve address %s", address)
	}
	future, err := p.transport.SendRequestPacket(ctx, request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ping address %s", address)
	}
	result, err := future.GetResponse(timeout)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to receive ping response from address %s", address)
	}
	return result.GetSenderHost(), nil
}

func NewPinger(transport network.InternalTransport) *Pinger {
	return &Pinger{transport: transport}
}
