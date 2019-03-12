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

package controller

import (
	"context"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
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

	hostNetwork  network.HostNetwork
	routingTable network.RoutingTable
}

func (pc *pulseController) Init(ctx context.Context) error {
	pc.hostNetwork.RegisterRequestHandler(types.Pulse, pc.processPulse)
	pc.hostNetwork.RegisterRequestHandler(types.GetRandomHosts, pc.processGetRandomHosts)
	return nil
}

func (pc *pulseController) processPulse(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestPulse)
	// we should not process pulses in Waiting state because network can be unready to join current node,
	// so we should wait for pulse from consensus phase1 packet
	// verified, err := pc.verifyPulseSign(data.Pulse)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "[ ServiceNetwork ] HandlePulse: ")
	// }
	// if !verified {
	// 	return nil, errors.New("[ ServiceNetwork ] HandlePulse: failed to verify a pulse sign")
	// }
	if pc.NodeKeeper.GetState() != core.WaitingNodeNetworkState {
		go pc.PulseHandler.HandlePulse(context.Background(), data.Pulse)
	}
	return pc.hostNetwork.BuildResponse(ctx, request, &packet.ResponsePulse{Success: true, Error: ""}), nil
}

func (pc *pulseController) processGetRandomHosts(ctx context.Context, request network.Request) (network.Response, error) {
	data := request.GetData().(*packet.RequestGetRandomHosts)
	randomHosts := pc.routingTable.GetRandomNodes(data.HostsNumber)
	return pc.hostNetwork.BuildResponse(ctx, request, &packet.ResponseGetRandomHosts{Hosts: randomHosts}), nil
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

func NewPulseController(hostNetwork network.HostNetwork, routingTable network.RoutingTable) PulseController {
	return &pulseController{hostNetwork: hostNetwork, routingTable: routingTable}
}
