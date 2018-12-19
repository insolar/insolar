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
}
