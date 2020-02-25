// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/pulse"

	"github.com/stretchr/testify/assert"
)

func generatePsc() *insolar.PulseSenderConfirmation {
	return &insolar.PulseSenderConfirmation{
		PulseNumber:     32,
		ChosenPublicKey: "124",
		Entropy:         insolar.Entropy{123},
		Signature:       []byte("456"),
	}
}

func TestPulseToProto(t *testing.T) {
	psc := generatePsc()
	signs := map[string]insolar.PulseSenderConfirmation{}
	signs["112"] = *psc
	p := insolar.Pulse{
		PulseNumber:      32,
		PrevPulseNumber:  22,
		NextPulseNumber:  42,
		PulseTimestamp:   111112,
		EpochPulseNumber: pulse.EphemeralPulseEpoch,
		OriginID:         [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 7, 6, 5, 4, 3, 2, 1},
		Entropy:          insolar.Entropy{123},
		Signs:            signs,
	}

	proto := ToProto(&p)
	p2 := FromProto(proto)
	assert.Equal(t, p.PulseNumber, p2.PulseNumber)
	assert.Equal(t, p.PrevPulseNumber, p2.PrevPulseNumber)
	assert.Equal(t, p.NextPulseNumber, p2.NextPulseNumber)
	assert.Equal(t, p.PulseTimestamp, p2.PulseTimestamp)
	assert.Equal(t, p.EpochPulseNumber, p2.EpochPulseNumber)
	assert.Equal(t, p.OriginID, p2.OriginID)
	assert.Equal(t, p.Entropy, p2.Entropy)
	assert.Equal(t, p.Signs, p2.Signs)
}

func TestPulseSenderConfirmationToProto(t *testing.T) {
	p := generatePsc()
	proto := SenderConfirmationToProto("112", *p)
	pk, p2 := SenderConfirmationFromProto(proto)
	assert.Equal(t, "112", pk)
	assert.EqualValues(t, p.PulseNumber, p2.PulseNumber)
	assert.Equal(t, p.ChosenPublicKey, p2.ChosenPublicKey)
	assert.Equal(t, p.Entropy, p2.Entropy)
	assert.Equal(t, p.Signature, p2.Signature)
}
