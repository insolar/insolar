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

package servicenetwork

import (
	"crypto"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/stretchr/testify/assert"
)

func TestNewServiceNetwork_incrementPort(t *testing.T) {
	addr, err := incrementPort("0.0.0.0:8080")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0:8081", addr)

	addr, err = incrementPort("[::]:8080")
	assert.NoError(t, err)
	assert.Equal(t, "[::]:8081", addr)

	addr, err = incrementPort("0.0.0.0:0")
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0:0", addr)

	addr, err = incrementPort("invalid_address")
	assert.Error(t, err)
	assert.Equal(t, "invalid_address", addr)

	addr, err = incrementPort("127.0.0.1:port")
	assert.Error(t, err)
	assert.Equal(t, "127.0.0.1:port", addr)
}

func getServiceNetwork(t *testing.T) *ServiceNetwork {
	cfg := configuration.NewConfiguration()
	cfg.Pulsar.PulseTime = pulseTimeMs // pulse 5 sec for faster tests
	cfg.Host.Transport.Address = "127.0.0.1:50001"
	cfg.Service.Skip = 5

	componentManager := &component.Manager{}
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	proc := platformpolicy.NewKeyProcessor()
	privKey, err := proc.GeneratePrivateKey()
	assert.NoError(t, err)
	service := cryptography.NewKeyBoundCryptographyService(privKey)

	serviceNetwork, err := NewServiceNetwork(cfg, componentManager, false)
	assert.NoError(t, err)

	serviceNetwork.CryptographyScheme = scheme
	serviceNetwork.KeyProcessor = proc
	serviceNetwork.CryptographyService = service

	return serviceNetwork
}

func getKeys(t *testing.T) (public string, private crypto.PrivateKey) {
	proc := platformpolicy.NewKeyProcessor()
	privKey, err := proc.GeneratePrivateKey()
	assert.NoError(t, err)
	key := proc.ExtractPublicKey(privKey)
	pubKey, err := proc.ExportPublicKeyBinary(key)
	assert.NoError(t, err)

	return string(pubKey[:]), privKey
}

func TestVerifyPulseSign(t *testing.T) {
	network := getServiceNetwork(t)
	keyStr, privateKey := getKeys(t)

	psc := core.PulseSenderConfirmation{
		PulseNumber:     1,
		ChosenPublicKey: string(keyStr[:]),
		Entropy:         randomEntropy(),
	}

	payload := pulsar.PulseSenderConfirmationPayload{PulseSenderConfirmation: psc}
	hasher := platformpolicy.NewPlatformCryptographyScheme().IntegrityHasher()
	hash, err := payload.Hash(hasher)
	assert.NoError(t, err)
	service := cryptography.NewKeyBoundCryptographyService(privateKey)
	sign, err := service.Sign(hash)
	assert.NoError(t, err)

	psc.Signature = sign.Bytes()

	pulse := pulsar.NewPulse(1, 0, &entropygenerator.StandardEntropyGenerator{})
	pulse.Signs = make(map[string]core.PulseSenderConfirmation, 1)
	pulse.Signs["test"] = psc

	valid, err := network.verifyPulseSign(*pulse)
	assert.NoError(t, err)
	assert.True(t, valid)
}

func randomEntropy() [64]byte {
	var buf [64]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		panic(buf)
	}
	return buf
}
