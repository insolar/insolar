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
)

const reserved = 0xDEADBEEF

type merkleHelper struct {
	hasher core.Hasher
}

func newMerkleHelper(scheme core.PlatformCryptographyScheme) *merkleHelper {
	return &merkleHelper{
		hasher: scheme.IntegrityHasher(),
	}
}

func (mh *merkleHelper) pulseHash(pulse *core.Pulse) []byte {
	var result []byte

	pulseNumberHash := mh.hasher.Hash(pulse.PulseNumber.Bytes())
	result = append(result, pulseNumberHash...)

	entropyHash := mh.hasher.Hash(pulse.Entropy[:])
	result = append(result, entropyHash...)

	return mh.hasher.Hash(result)
}

func (mh *merkleHelper) nodeInfoHash(pulseHash, stateHash []byte) []byte {
	var result []byte

	result = append(result, pulseHash...)
	result = append(result, stateHash...)

	return mh.hasher.Hash(result)
}

func (mh *merkleHelper) nodeHash(nodeSignature, nodeInfoHash []byte) []byte {
	var result []byte

	nodeSignatureHash := mh.hasher.Hash(nodeSignature)
	result = append(result, nodeSignatureHash...)

	result = append(result, nodeInfoHash...)

	return mh.hasher.Hash(result)
}

func (mh *merkleHelper) bucketEntryHash(entryIndex uint32, nodeHash []byte) []byte {
	var result []byte

	entryIndexHash := mh.hasher.Hash(utils.UInt32ToBytes(entryIndex))
	result = append(result, entryIndexHash...)

	result = append(result, nodeHash...)

	return mh.hasher.Hash(result)
}

func (mh *merkleHelper) bucketInfoHash(role core.NodeRole, nodeCount uint32) []byte {
	var result []byte

	roleHash := mh.hasher.Hash(utils.UInt32ToBytes(uint32(role)))
	result = append(result, roleHash...)

	nodeCountHash := mh.hasher.Hash(utils.UInt32ToBytes(nodeCount))
	result = append(result, nodeCountHash...)

	return mh.hasher.Hash(result)
}

func (mh *merkleHelper) bucketHash(bucketInfoHash, bucketEntryHash []byte) []byte {
	var result []byte

	result = append(result, bucketInfoHash...)
	result = append(result, bucketEntryHash...)

	return mh.hasher.Hash(result)
}

func (mh *merkleHelper) globuleInfoHash(prevCloudHash []byte, gobuleIndex, nodeCount uint32) []byte {
	reservedHash := mh.hasher.Hash(utils.UInt32ToBytes(reserved))

	var tmpResult1 []byte

	tmpResult1 = append(tmpResult1, reservedHash...)
	tmpResult1 = append(tmpResult1, prevCloudHash...)

	var tmpResult2 []byte

	globuleIndexHash := mh.hasher.Hash(utils.UInt32ToBytes(gobuleIndex))
	tmpResult2 = append(tmpResult2, globuleIndexHash...)

	nodeCountHash := mh.hasher.Hash(utils.UInt32ToBytes(nodeCount))
	tmpResult2 = append(tmpResult2, nodeCountHash...)

	var tmpResult3 []byte

	tmpResult1Hash := mh.hasher.Hash(tmpResult1)
	tmpResult3 = append(tmpResult3, tmpResult1Hash...)

	tmpResult2Hash := mh.hasher.Hash(tmpResult2)
	tmpResult3 = append(tmpResult3, tmpResult2Hash...)

	return mh.hasher.Hash(tmpResult3)
}

func (mh *merkleHelper) globuleHash(globuleInfoHash, globuleNodeRoot []byte) []byte {
	var result []byte

	result = append(result, globuleInfoHash...)
	result = append(result, globuleNodeRoot...)

	return mh.hasher.Hash(result)
}
