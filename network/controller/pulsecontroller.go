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

package controller

import (
	"context"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/pulsar"
	"github.com/pkg/errors"
)

type PulseController interface {
	component.Initer
}

type pulseController struct {
	PulseHandler        network.PulseHandler            `inject:""`
	NodeKeeper          network.NodeKeeper              `inject:""`
	CryptographyScheme  core.PlatformCryptographyScheme `inject:""`
	KeyProcessor        core.KeyProcessor               `inject:""`
	CryptographyService core.CryptographyService        `inject:""`
	Resolver            network.RoutingTable            `inject:""`
	Network             network.HostNetwork             `inject:""`
}

func (pc *pulseController) Init(ctx context.Context) error {
	pc.Network.RegisterRequestHandler(types.Pulse, pc.processPulse)
	pc.Network.RegisterRequestHandler(types.GetRandomHosts, pc.processGetRandomHosts)
	return nil
}

func (pc *pulseController) processPulse(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestPulse)
	verified, err := pc.verifyPulseSign(data.Pulse)
	if err != nil {
		return nil, errors.Wrap(err, "[ pulseController ] processPulse: error to verify a pulse sign")
	}
	if !verified {
		return nil, errors.New("[ pulseController ] processPulse: failed to verify a pulse sign")
	}
	// if we are a joiner node, we should receive pulse from phase1 packet and ignore pulse from pulsar
	if !pc.NodeKeeper.GetConsensusInfo().IsJoiner() {
		go pc.PulseHandler.HandlePulse(context.Background(), data.Pulse)
	} else {
		log.Debugf("Ignore pulse %v from pulsar, waiting for consensus phase1 packet", data.Pulse)
	}
	return pc.Network.BuildResponse(ctx, request, &packet.ResponsePulse{Success: true, Error: ""}), nil
}

func (pc *pulseController) processGetRandomHosts(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestGetRandomHosts)
	randomHosts := pc.Resolver.GetRandomNodes(data.HostsNumber)
	return pc.Network.BuildResponse(ctx, request, &packet.ResponseGetRandomHosts{Hosts: randomHosts}), nil
}

func (pc *pulseController) verifyPulseSign(pulse core.Pulse) (bool, error) {
	hashProvider := pc.CryptographyScheme.IntegrityHasher()
	if len(pulse.Signs) == 0 {
		return false, errors.New("[ verifyPulseSign ] received empty pulse signs")
	}
	for _, psc := range pulse.Signs {
		payload := pulsar.PulseSenderConfirmationPayload{PulseSenderConfirmation: psc}
		hash, err := payload.Hash(hashProvider)
		if err != nil {
			return false, errors.Wrap(err, "[ verifyPulseSign ] error to get a hash from pulse payload")
		}
		key, err := pc.KeyProcessor.ImportPublicKeyPEM([]byte(psc.ChosenPublicKey))
		if err != nil {
			return false, errors.Wrap(err, "[ verifyPulseSign ] error to import a public key")
		}

		verified := pc.CryptographyService.Verify(key, core.SignatureFromBytes(psc.Signature), hash)

		if !verified {
			return false, errors.New("[ verifyPulseSign ] error to verify a pulse")
		}
	}
	return true, nil
}

func NewPulseController() PulseController {
	return &pulseController{}
}
