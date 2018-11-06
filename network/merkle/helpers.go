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

package merkle

import (
	"context"
	ecdsa2 "crypto/ecdsa"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func verifySignature(ctx context.Context, data, signature []byte, publicKey *ecdsa2.PublicKey) bool {
	log := inslogger.FromContext(ctx)

	key, err := ecdsa.ExportPublicKey(publicKey)
	if err != nil {
		log.Error("Failed to export a public key: ", err)
		return false
	}

	verified, err := ecdsa.Verify(data, signature, key)
	if err != nil {
		log.Error("Failed to verify signature: ", err)
		return false
	}

	return verified
}

func pulseHash(pulse *core.Pulse) []byte {
	var result []byte

	pulseNumberHash := hash.SHA3Bytes256(pulse.PulseNumber.Bytes())
	result = append(result, pulseNumberHash...)

	entropyHash := hash.SHA3Bytes256(pulse.Entropy[:])
	result = append(result, entropyHash...)

	return hash.SHA3Bytes256(result)
}

func nodeInfoHash(pulseHash, stateHash []byte) []byte {
	var result []byte

	result = append(result, pulseHash...)
	result = append(result, stateHash...)

	return hash.SHA3Bytes256(result)
}

func nodeHash(nodeSignature, nodeInfoHash []byte) []byte {
	var result []byte

	pulseNumberHash := hash.SHA3Bytes256(nodeSignature)
	result = append(result, pulseNumberHash...)

	result = append(result, nodeInfoHash...)

	return hash.SHA3Bytes256(result)
}

func globuleInfoHash(prevCloudHash []byte, gobuleIndex, nodeCount uint32) []byte {
	reservedHash := hash.SHA3Bytes256(utils.UInt32ToBytes(0xDEADBEEF))

	var tmpResult1 []byte

	tmpResult1 = append(tmpResult1, reservedHash...)
	tmpResult1 = append(tmpResult1, prevCloudHash...)

	var tmpResult2 []byte

	globuleIndexHash := hash.SHA3Bytes256(utils.UInt32ToBytes(gobuleIndex))
	tmpResult2 = append(tmpResult2, globuleIndexHash...)

	nodeCountHash := hash.SHA3Bytes256(utils.UInt32ToBytes(nodeCount))
	tmpResult2 = append(tmpResult2, nodeCountHash...)

	var tmpResult3 []byte

	tmpResult1Hash := hash.SHA3Bytes256(tmpResult1)
	tmpResult3 = append(tmpResult3, tmpResult1Hash...)

	tmpResult2Hash := hash.SHA3Bytes256(tmpResult2)
	tmpResult3 = append(tmpResult3, tmpResult2Hash...)

	return hash.SHA3Bytes256(tmpResult3)
}

func globuleHash(globuleInfoHash, globuleNodeRoot []byte) []byte {
	var result []byte

	result = append(result, globuleInfoHash...)
	result = append(result, globuleNodeRoot...)

	return hash.SHA3Bytes256(result)
}
