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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/cryptohelpers/hash"
)

const reserved = 0xDEADBEEF

func pulseHash(pulse *core.Pulse) []byte {
	var result []byte

	pulseNumberHash := hash.IntegrityHasher(pulse.PulseNumber.Bytes())
	result = append(result, pulseNumberHash...)

	entropyHash := hash.IntegrityHasher(pulse.Entropy[:])
	result = append(result, entropyHash...)

	return hash.IntegrityHasher(result)
}

func nodeInfoHash(pulseHash, stateHash []byte) []byte {
	var result []byte

	result = append(result, pulseHash...)
	result = append(result, stateHash...)

	return hash.IntegrityHasher(result)
}

func nodeHash(nodeSignature, nodeInfoHash []byte) []byte {
	var result []byte

	nodeSignatureHash := hash.IntegrityHasher(nodeSignature)
	result = append(result, nodeSignatureHash...)

	result = append(result, nodeInfoHash...)

	return hash.IntegrityHasher(result)
}

func bucketEntryHash(entryIndex uint32, nodeHash []byte) []byte {
	var result []byte

	entryIndexHash := hash.IntegrityHasher(utils.UInt32ToBytes(entryIndex))
	result = append(result, entryIndexHash...)

	result = append(result, nodeHash...)

	return hash.IntegrityHasher(result)
}

func bucketInfoHash(role core.NodeRole, nodeCount uint32) []byte {
	var result []byte

	roleHash := hash.IntegrityHasher(utils.UInt32ToBytes(uint32(role)))
	result = append(result, roleHash...)

	nodeCountHash := hash.IntegrityHasher(utils.UInt32ToBytes(nodeCount))
	result = append(result, nodeCountHash...)

	return hash.IntegrityHasher(result)
}

func bucketHash(bucketInfoHash, bucketEntryHash []byte) []byte {
	var result []byte

	result = append(result, bucketInfoHash...)
	result = append(result, bucketEntryHash...)

	return hash.IntegrityHasher(result)
}

func globuleInfoHash(prevCloudHash []byte, gobuleIndex, nodeCount uint32) []byte {
	reservedHash := hash.IntegrityHasher(utils.UInt32ToBytes(reserved))

	var tmpResult1 []byte

	tmpResult1 = append(tmpResult1, reservedHash...)
	tmpResult1 = append(tmpResult1, prevCloudHash...)

	var tmpResult2 []byte

	globuleIndexHash := hash.IntegrityHasher(utils.UInt32ToBytes(gobuleIndex))
	tmpResult2 = append(tmpResult2, globuleIndexHash...)

	nodeCountHash := hash.IntegrityHasher(utils.UInt32ToBytes(nodeCount))
	tmpResult2 = append(tmpResult2, nodeCountHash...)

	var tmpResult3 []byte

	tmpResult1Hash := hash.IntegrityHasher(tmpResult1)
	tmpResult3 = append(tmpResult3, tmpResult1Hash...)

	tmpResult2Hash := hash.IntegrityHasher(tmpResult2)
	tmpResult3 = append(tmpResult3, tmpResult2Hash...)

	return hash.IntegrityHasher(tmpResult3)
}

func globuleHash(globuleInfoHash, globuleNodeRoot []byte) []byte {
	var result []byte

	result = append(result, globuleInfoHash...)
	result = append(result, globuleNodeRoot...)

	return hash.IntegrityHasher(result)
}
