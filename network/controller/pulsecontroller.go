// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package controller

import (
	"context"

	"github.com/pkg/errors"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/pulsar"
)

type PulseController interface {
	component.Initer
}

type pulseController struct {
	PulseHandler        network.PulseHandler               `inject:""`
	CryptographyScheme  insolar.PlatformCryptographyScheme `inject:""`
	KeyProcessor        insolar.KeyProcessor               `inject:""`
	CryptographyService insolar.CryptographyService        `inject:""`
	Network             network.HostNetwork                `inject:""`
}

func (pc *pulseController) Init(ctx context.Context) error {
	pc.Network.RegisterRequestHandler(types.Pulse, pc.processPulse)
	return nil
}

func (pc *pulseController) processPulse(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetPulse() == nil {
		return nil, errors.Errorf("process pulse: got invalid protobuf request message: %s", request)
	}

	logger := inslogger.FromContext(ctx)

	data := request.GetRequest().GetPulse()
	p := *pulse.FromProto(data.Pulse)
	err := pc.verifyPulseSign(p)
	if err != nil {
		logger.Error("processPulse: failed to verify p: ", err.Error())
		return nil, errors.Wrap(err, "[ pulseController ] processPulse: failed to verify pulse")
	}

	go pc.PulseHandler.HandlePulse(ctx, p, request)
	return nil, nil
}

func (pc *pulseController) verifyPulseSign(pulse insolar.Pulse) error {
	if len(pulse.Signs) == 0 {
		return errors.New("received empty pulse signs")
	}
	hashProvider := pc.CryptographyScheme.IntegrityHasher()
	for key, psc := range pulse.Signs {
		payload := pulsar.PulseSenderConfirmationPayload{PulseSenderConfirmation: psc}
		hash, err := payload.Hash(hashProvider)
		if err != nil {
			return errors.Wrap(err, "failed to get hash from pulse payload")
		}
		pk, err := pc.KeyProcessor.ImportPublicKeyPEM([]byte(key))
		if err != nil {
			return errors.Wrap(err, "failed to import public key")
		}

		verified := pc.CryptographyService.Verify(pk, insolar.SignatureFromBytes(psc.Signature), hash)

		if !verified {
			return errors.New("cryptographic signature verification failed")
		}
	}
	return nil
}

func NewPulseController() PulseController {
	return &pulseController{}
}
