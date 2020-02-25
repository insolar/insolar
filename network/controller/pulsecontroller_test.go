// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
	request.SetRequest(&packet.RPCRequest{})
	_, err := controller.processPulse(context.Background(), request)
	assert.Error(t, err)
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
