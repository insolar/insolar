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
	"crypto"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	network2 "github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getController(t *testing.T) *pulseController {
	proc := platformpolicy.NewKeyProcessor()
	key, err := proc.GeneratePrivateKey()
	require.NoError(t, err)

	pulseHandler := network.NewPulseHandlerMock(t)
	pulseHandler.HandlePulseMock.Set(func(context.Context, insolar.Pulse, network2.ReceivedPacket) {})
	net := network.NewHostNetworkMock(t)
	net.BuildResponseMock.Return(packet.NewPacket(nil, nil, types.Pulse, 1))

	return &pulseController{
		PulseHandler:        pulseHandler,
		Network:             net,
		CryptographyScheme:  platformpolicy.NewPlatformCryptographyScheme(),
		KeyProcessor:        proc,
		CryptographyService: cryptography.NewKeyBoundCryptographyService(key),
	}
}

func getKeys(t *testing.T) (public string, private crypto.PrivateKey) {
	proc := platformpolicy.NewKeyProcessor()
	privKey, err := proc.GeneratePrivateKey()
	require.NoError(t, err)
	key := proc.ExtractPublicKey(privKey)
	pubKey, err := proc.ExportPublicKeyPEM(key)
	require.NoError(t, err)

	return string(pubKey), privKey
}

func signPulsePayload(t *testing.T, psc insolar.PulseSenderConfirmation, key crypto.PrivateKey) []byte {
	payload := pulsar.PulseSenderConfirmationPayload{PulseSenderConfirmation: psc}
	hasher := platformpolicy.NewPlatformCryptographyScheme().IntegrityHasher()
	hash, err := payload.Hash(hasher)
	require.NoError(t, err)
	service := cryptography.NewKeyBoundCryptographyService(key)
	sign, err := service.Sign(hash)
	require.NoError(t, err)
	return sign.Bytes()
}

func TestVerifyPulseSignTrue(t *testing.T) {
	controller := getController(t)
	keyStr, privateKey := getKeys(t)

	psc := insolar.PulseSenderConfirmation{
		PulseNumber:     1,
		ChosenPublicKey: keyStr,
		Entropy:         randomEntropy(),
	}

	psc.Signature = signPulsePayload(t, psc, privateKey)

	pulse := pulsar.NewPulse(1, 0, &entropygenerator.StandardEntropyGenerator{})
	pulse.Signs = make(map[string]insolar.PulseSenderConfirmation)
	pulse.Signs[keyStr] = psc

	err := controller.verifyPulseSign(*pulse)
	assert.NoError(t, err)

	psc.ChosenPublicKey = "some_other_key"
	psc.Signature = signPulsePayload(t, psc, privateKey)
	pulse.Signs[keyStr] = psc
	err = controller.verifyPulseSign(*pulse)
	assert.NoError(t, err)
}

func TestVerifyPulseSignFalse(t *testing.T) {
	controller := getController(t)
	keyStr, _ := getKeys(t)

	psc := insolar.PulseSenderConfirmation{
		PulseNumber:     1,
		ChosenPublicKey: keyStr,
		Entropy:         randomEntropy(),
	}

	psc.Signature = []byte("test")

	pulse := pulsar.NewPulse(1, 0, &entropygenerator.StandardEntropyGenerator{})
	pulse.Signs = make(map[string]insolar.PulseSenderConfirmation)
	pulse.Signs[keyStr] = psc

	err := controller.verifyPulseSign(*pulse)
	assert.Error(t, err)

	pulse.Signs = make(map[string]insolar.PulseSenderConfirmation)
	pulse.Signs["tratata"] = psc

	err = controller.verifyPulseSign(*pulse)
	assert.Error(t, err)
}

func TestVerifyPulseSignEmpty(t *testing.T) {
	controller := getController(t)
	pulse := pulsar.NewPulse(1, 0, &entropygenerator.StandardEntropyGenerator{})
	err := controller.verifyPulseSign(*pulse)
	assert.Error(t, err)
}

func randomEntropy() [64]byte {
	var buf [64]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(buf)
	}
	return buf
}

func newPulsePacket(t *testing.T) *packet.ReceivedPacket {
	refs := gen.UniqueReferences(2)
	sender, err := host.NewHostN("127.0.0.1:3344", refs[0])
	require.NoError(t, err)
	receiver, err := host.NewHostN("127.0.0.1:3345", refs[1])
	require.NoError(t, err)
	return packet.NewReceivedPacket(packet.NewPacket(sender, receiver, types.Pulse, 1), nil)
}

func TestProcessIncorrectPacket(t *testing.T) {
	controller := &pulseController{}
	request := newPulsePacket(t)
	request.SetRequest(&packet.Ping{})
	_, err := controller.processPulse(context.Background(), request)
	assert.Error(t, err)
	request.SetResponse(&packet.BasicResponse{Success: true})
	_, err = controller.processPulse(context.Background(), request)
	assert.Error(t, err)
}

func TestProcessPulseVerifyFailure(t *testing.T) {
	controller := &pulseController{}
	request := newPulsePacket(t)
	request.SetRequest(&packet.PulseRequest{
		Pulse: pulse.ToProto(pulsar.NewPulse(10, 140,
			&entropygenerator.StandardEntropyGenerator{}))},
	)
	_, err := controller.processPulse(context.Background(), request)
	assert.Error(t, err)
}

func newSignedPulse(t *testing.T) *insolar.Pulse {
	keyStr, privateKey := getKeys(t)

	psc := insolar.PulseSenderConfirmation{
		PulseNumber:     1,
		ChosenPublicKey: keyStr,
		Entropy:         randomEntropy(),
	}

	psc.Signature = signPulsePayload(t, psc, privateKey)

	pulse := pulsar.NewPulse(10, 140, &entropygenerator.StandardEntropyGenerator{})
	pulse.Signs = make(map[string]insolar.PulseSenderConfirmation)
	pulse.Signs[keyStr] = psc

	return pulse
}

func TestProcessPulseHappyPath(t *testing.T) {
	controller := getController(t)

	request := newPulsePacket(t)
	request.SetRequest(&packet.PulseRequest{
		Pulse: pulse.ToProto(newSignedPulse(t)),
	})
	_, err := controller.processPulse(context.Background(), request)
	assert.NoError(t, err)
}
