/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package pulsar

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/stretchr/testify/require"
)

func TestPreparePayloadAndCheckIt(t *testing.T) {
	t.Parallel()

	keyProcessor := platformpolicy.NewKeyProcessor()
	privateKey, err := keyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	cryptoService := cryptography.NewKeyBoundCryptographyService(privateKey)
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulsar, err := NewPulsar(
		configuration.Pulsar{},
		cryptoService,
		scheme,
		keyProcessor,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	t.Run("HandshakePayload payload", func(t *testing.T){
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		handshakePayload := &HandshakePayload{Entropy: entropyGenerator.GenerateEntropy()}

		// Act
		payload, firstError := pulsar.preparePayload(handshakePayload)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)

		// Assert
		require.NoError(t, firstError)
		require.NoError(t, secondError)
		require.Equal(t, true, isVerified)
		require.Equal(t, handshakePayload, payload.Body.(*HandshakePayload))
	})

	t.Run("EntropySignaturePayload payload", func(t *testing.T){
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		entropySignPayload := &EntropySignaturePayload{EntropySignature: entropyGenerator.GenerateEntropy()[:]}

		// Act
		payload, firstError := pulsar.preparePayload(entropySignPayload)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)

		// Assert
		require.NoError(t, firstError)
		require.NoError(t, secondError)
		require.Equal(t, true, isVerified)
		require.Equal(t, entropySignPayload, payload.Body.(*EntropySignaturePayload))
	})

	t.Run("EntropyPayload payload", func(t *testing.T){
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		entropyPayload := &EntropyPayload{Entropy: entropyGenerator.GenerateEntropy()}

		// Act
		payload, firstError := pulsar.preparePayload(entropyPayload)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)

		// Assert
		require.NoError(t, firstError)
		require.NoError(t, secondError)
		require.Equal(t, true, isVerified)
		require.Equal(t, entropyPayload, payload.Body.(*EntropyPayload))
	})

	t.Run("VectorPayload payload", func(t *testing.T){
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		firstEntropy := entropyGenerator.GenerateEntropy()
		secondEntropy := entropyGenerator.GenerateEntropy()
		firstVector := &VectorPayload{Vector: map[string]*BftCell{
			"first" : &BftCell{Entropy:firstEntropy},
			"second" : &BftCell{Entropy:secondEntropy},
		}}
		secondVector := &VectorPayload{Vector: map[string]*BftCell{
			"first" : &BftCell{Entropy:firstEntropy},
			"second" : &BftCell{Entropy:secondEntropy},
		}}

		t.Run("preparePayload works for VectorPayload", func(t *testing.T){
			// Act
			payload, firstError := pulsar.preparePayload(firstVector)
			require.NotNil(t, payload)
			isVerified, secondError := pulsar.checkPayloadSignature(payload)

			// Assert
			require.NoError(t, firstError)
			require.NoError(t, secondError)
			require.Equal(t, true, isVerified)
			require.Equal(t, firstVector, payload.Body.(*VectorPayload))
		})

		t.Run("checkPayloadSignature work for maps", func(t *testing.T){
			// Act
			payload, firstError := pulsar.preparePayload(firstVector)
			require.NotNil(t, payload)
			payload.Body = secondVector
			isVerified, secondError := pulsar.checkPayloadSignature(payload)

			// Assert
			require.NoError(t, firstError)
			require.NoError(t, secondError)
			require.Equal(t, true, isVerified)
			require.Equal(t, secondVector, payload.Body.(*VectorPayload))
		})
	})

	t.Run("PulsePayload payload", func(t *testing.T){
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		firstEntropy := entropyGenerator.GenerateEntropy()
		secondEntropy := entropyGenerator.GenerateEntropy()
		pulsePayload := &PulsePayload{
			Pulse : core.Pulse{
				Entropy: entropyGenerator.GenerateEntropy(),
				Signs: map[string]core.PulseSenderConfirmation{
					"first" : core.PulseSenderConfirmation{Entropy:firstEntropy},
					"second" : core.PulseSenderConfirmation{Entropy:secondEntropy},
				},
		}}
		secondPulsePayload := &PulsePayload{
			Pulse : core.Pulse{
				Entropy: entropyGenerator.GenerateEntropy(),
				Signs: map[string]core.PulseSenderConfirmation{
					"second" : core.PulseSenderConfirmation{Entropy:secondEntropy},
					"first" : core.PulseSenderConfirmation{Entropy:firstEntropy},
				},
			}}

		t.Run("preparePayload works for PulsePayload", func(t *testing.T){
			// Act
			payload, firstError := pulsar.preparePayload(pulsePayload)
			require.NotNil(t, payload)
			isVerified, secondError := pulsar.checkPayloadSignature(payload)

			// Assert
			require.NoError(t, firstError)
			require.NoError(t, secondError)
			require.Equal(t, true, isVerified)
			require.Equal(t, pulsePayload, payload.Body.(*PulsePayload))
		})

		t.Run("checkPayloadSignature work for maps", func(t *testing.T){
			// Act
			payload, firstError := pulsar.preparePayload(pulsePayload)
			require.NotNil(t, payload)
			payload.Body = secondPulsePayload
			isVerified, secondError := pulsar.checkPayloadSignature(payload)

			// Assert
			require.NoError(t, firstError)
			require.NoError(t, secondError)
			require.Equal(t, true, isVerified)
			require.Equal(t, secondPulsePayload, payload.Body.(*PulsePayload))
		})
	})

	t.Run("PulseSenderConfirmationPayload payload", func(t *testing.T){
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		payloadBody := &PulseSenderConfirmationPayload{core.PulseSenderConfirmation{Entropy: entropyGenerator.GenerateEntropy()}}

		// Act
		payload, firstError := pulsar.preparePayload(payloadBody)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)

		// Assert
		require.NoError(t, firstError)
		require.NoError(t, secondError)
		require.Equal(t, true, isVerified)
		require.Equal(t, payloadBody, payload.Body.(*PulseSenderConfirmationPayload))

	})
}
