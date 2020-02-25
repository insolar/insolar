// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/pulse"
)

func FromProto(p *PulseProto) *insolar.Pulse {
	result := &insolar.Pulse{
		PulseNumber:      p.PulseNumber,
		PrevPulseNumber:  p.PrevPulseNumber,
		NextPulseNumber:  p.NextPulseNumber,
		PulseTimestamp:   p.PulseTimestamp,
		EpochPulseNumber: pulse.Epoch(p.EpochPulseNumber),
		Entropy:          p.Entropy,
		Signs:            map[string]insolar.PulseSenderConfirmation{},
	}
	copy(result.OriginID[:], p.OriginID)
	for _, sign := range p.Signs {
		pk, confirmation := SenderConfirmationFromProto(sign)
		result.Signs[pk] = confirmation
	}
	return result
}

func ToProto(p *insolar.Pulse) *PulseProto {
	result := &PulseProto{
		PulseNumber:      p.PulseNumber,
		PrevPulseNumber:  p.PrevPulseNumber,
		NextPulseNumber:  p.NextPulseNumber,
		PulseTimestamp:   p.PulseTimestamp,
		EpochPulseNumber: int32(p.EpochPulseNumber),
		OriginID:         p.OriginID[:],
		Entropy:          p.Entropy,
	}
	for pk, sign := range p.Signs {
		result.Signs = append(result.Signs, SenderConfirmationToProto(pk, sign))
	}
	return result
}

func SenderConfirmationToProto(publicKey string, p insolar.PulseSenderConfirmation) *PulseSenderConfirmationProto {
	return &PulseSenderConfirmationProto{
		PublicKey:       publicKey,
		PulseNumber:     p.PulseNumber,
		ChosenPublicKey: p.ChosenPublicKey,
		Entropy:         p.Entropy,
		Signature:       p.Signature,
	}
}

func SenderConfirmationFromProto(p *PulseSenderConfirmationProto) (string, insolar.PulseSenderConfirmation) {
	return p.PublicKey, insolar.PulseSenderConfirmation{
		PulseNumber:     p.PulseNumber,
		ChosenPublicKey: p.ChosenPublicKey,
		Entropy:         p.Entropy,
		Signature:       p.Signature,
	}
}
