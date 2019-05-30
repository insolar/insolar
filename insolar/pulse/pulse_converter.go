//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pulse

import (
	"github.com/insolar/insolar/insolar"
)

func FromProto(p *PulseProto) *insolar.Pulse {
	result := &insolar.Pulse{
		PulseNumber:      p.PulseNumber,
		PrevPulseNumber:  p.PrevPulseNumber,
		NextPulseNumber:  p.NextPulseNumber,
		PulseTimestamp:   p.PulseTimestamp,
		EpochPulseNumber: int(p.EpochPulseNumber),
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
