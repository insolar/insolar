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

	pubKey, err := cryptoService.GetPublicKey()
	require.NoError(t, err)
	pubKeyRaw, err := keyProcessor.ExportPublicKeyPEM(pubKey)
	require.NoError(t, err)

	pulsar := Pulsar{
		CryptographyService:        cryptoService,
		PlatformCryptographyScheme: scheme,
		KeyProcessor:               keyProcessor,
		PublicKeyRaw:               string(pubKeyRaw),
	}

	t.Run("HandshakePayload", func(t *testing.T) {
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		handshakePayload := &HandshakePayload{Entropy: entropyGenerator.GenerateEntropy()}

		// Act
		payload, firstError := pulsar.preparePayload(handshakePayload)
		require.NoError(t, firstError)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)
		require.NoError(t, secondError)

		// Assert
		require.Equal(t, true, isVerified)
		require.Equal(t, handshakePayload, payload.Body.(*HandshakePayload))
	})

	t.Run("EntropySignaturePayload", func(t *testing.T) {
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		entropy := entropyGenerator.GenerateEntropy()
		entropySignPayload := &EntropySignaturePayload{EntropySignature: entropy[:]}

		// Act
		payload, firstError := pulsar.preparePayload(entropySignPayload)
		require.NoError(t, firstError)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)
		require.NoError(t, secondError)

		// Assert
		require.Equal(t, true, isVerified)
		require.Equal(t, entropySignPayload, payload.Body.(*EntropySignaturePayload))
	})

	t.Run("EntropyPayload", func(t *testing.T) {
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		entropyPayload := &EntropyPayload{Entropy: entropyGenerator.GenerateEntropy()}

		// Act
		payload, firstError := pulsar.preparePayload(entropyPayload)
		require.NoError(t, firstError)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)
		require.NoError(t, secondError)

		// Assert
		require.Equal(t, true, isVerified)
		require.Equal(t, entropyPayload, payload.Body.(*EntropyPayload))
	})

	t.Run("VectorPayload", func(t *testing.T) {
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		firstEntropy := entropyGenerator.GenerateEntropy()
		secondEntropy := entropyGenerator.GenerateEntropy()
		firstVector := &VectorPayload{Vector: map[string]*BftCell{
			"first":  &BftCell{Entropy: firstEntropy},
			"second": &BftCell{Entropy: secondEntropy},
		}}
		secondVector := &VectorPayload{Vector: map[string]*BftCell{
			"first":  &BftCell{Entropy: firstEntropy},
			"second": &BftCell{Entropy: secondEntropy},
		}}

		t.Run("preparePayload works for VectorPayload", func(t *testing.T) {
			// Act
			payload, firstError := pulsar.preparePayload(firstVector)
			require.NoError(t, firstError)
			require.NotNil(t, payload)
			isVerified, secondError := pulsar.checkPayloadSignature(payload)
			require.NoError(t, secondError)

			// Assert
			require.Equal(t, true, isVerified)
			require.Equal(t, firstVector, payload.Body.(*VectorPayload))
		})

		t.Run("checkPayloadSignature work for maps", func(t *testing.T) {
			// Act
			payload, firstError := pulsar.preparePayload(firstVector)
			require.NoError(t, firstError)
			require.NotNil(t, payload)
			payload.Body = secondVector
			isVerified, secondError := pulsar.checkPayloadSignature(payload)
			require.NoError(t, secondError)

			// Assert
			require.Equal(t, true, isVerified)
			require.Equal(t, secondVector, payload.Body.(*VectorPayload))
		})
	})

	t.Run("PulsePayload", func(t *testing.T) {
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		firstEntropy := entropyGenerator.GenerateEntropy()
		secondEntropy := entropyGenerator.GenerateEntropy()
		pulseEntropy := entropyGenerator.GenerateEntropy()
		pulsePayload := &PulsePayload{
			Pulse: core.Pulse{
				Entropy: pulseEntropy,
				Signs: map[string]core.PulseSenderConfirmation{
					"first":  core.PulseSenderConfirmation{Entropy: firstEntropy},
					"second": core.PulseSenderConfirmation{Entropy: secondEntropy},
				},
			}}
		secondPulsePayload := &PulsePayload{
			Pulse: core.Pulse{
				Entropy: pulseEntropy,
				Signs: map[string]core.PulseSenderConfirmation{
					"second": core.PulseSenderConfirmation{Entropy: secondEntropy},
					"first":  core.PulseSenderConfirmation{Entropy: firstEntropy},
				},
			}}

		t.Run("preparePayload works for PulsePayload", func(t *testing.T) {
			// Act
			payload, firstError := pulsar.preparePayload(pulsePayload)
			require.NoError(t, firstError)
			require.NotNil(t, payload)
			isVerified, secondError := pulsar.checkPayloadSignature(payload)
			require.NoError(t, secondError)

			// Assert
			require.Equal(t, true, isVerified)
			require.Equal(t, pulsePayload, payload.Body.(*PulsePayload))
		})

		t.Run("checkPayloadSignature work for maps", func(t *testing.T) {
			// Act
			payload, firstError := pulsar.preparePayload(pulsePayload)
			require.NoError(t, firstError)
			require.NotNil(t, payload)
			payload.Body = secondPulsePayload
			isVerified, secondError := pulsar.checkPayloadSignature(payload)
			require.NoError(t, secondError)

			// Assert
			require.Equal(t, true, isVerified)
			require.Equal(t, secondPulsePayload, payload.Body.(*PulsePayload))
		})
	})

	t.Run("PulseSenderConfirmationPayload", func(t *testing.T) {
		// Arrange
		entropyGenerator := entropygenerator.StandardEntropyGenerator{}
		payloadBody := &PulseSenderConfirmationPayload{core.PulseSenderConfirmation{Entropy: entropyGenerator.GenerateEntropy()}}

		// Act
		payload, firstError := pulsar.preparePayload(payloadBody)
		require.NoError(t, firstError)
		require.NotNil(t, payload)
		isVerified, secondError := pulsar.checkPayloadSignature(payload)
		require.NoError(t, secondError)

		// Assert
		require.Equal(t, true, isVerified)
		require.Equal(t, payloadBody, payload.Body.(*PulseSenderConfirmationPayload))

	})
}
