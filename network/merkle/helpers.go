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
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
)

const reserved = 0xDEADBEEF

type merkleHelper struct {
	scheme     core.PlatformCryptographyScheme
	leafHasher core.Hasher
}

func newMerkleHelper(scheme core.PlatformCryptographyScheme) *merkleHelper {
	return &merkleHelper{
		scheme:     scheme,
		leafHasher: scheme.IntegrityHasher(),
	}
}

func (mh *merkleHelper) doubleSliceHash(slice1, slice2 []byte) []byte {
	hasher := mh.scheme.IntegrityHasher()
	var err error

	_, err = hasher.Write(slice1)
	if err != nil {
		panic(fmt.Sprintf("[ doubleSliceHash ] Hash write error: %s", err.Error()))
	}
	_, err = hasher.Write(slice2)
	if err != nil {
		panic(fmt.Sprintf("[ doubleSliceHash ] Hash write error: %s", err.Error()))
	}

	return hasher.Sum(nil)
}

func (mh *merkleHelper) pulseHash(pulse *core.Pulse) []byte {
	pulseNumberHash := mh.leafHasher.Hash(pulse.PulseNumber.Bytes())
	entropyHash := mh.leafHasher.Hash(pulse.Entropy[:])

	return mh.doubleSliceHash(pulseNumberHash, entropyHash)
}

func (mh *merkleHelper) nodeInfoHash(pulseHash, stateHash []byte) []byte {
	return mh.doubleSliceHash(pulseHash, stateHash)
}

func (mh *merkleHelper) nodeHash(nodeSignature, nodeInfoHash []byte) []byte {
	nodeSignatureHash := mh.leafHasher.Hash(nodeSignature)
	return mh.doubleSliceHash(nodeSignatureHash, nodeInfoHash)
}

func (mh *merkleHelper) bucketEntryHash(entryIndex uint32, nodeHash []byte) []byte {
	entryIndexHash := mh.leafHasher.Hash(utils.UInt32ToBytes(entryIndex))
	return mh.doubleSliceHash(entryIndexHash, nodeHash)
}

func (mh *merkleHelper) bucketInfoHash(role core.StaticRole, nodeCount uint32) []byte {
	roleHash := mh.leafHasher.Hash(utils.UInt32ToBytes(uint32(role)))
	nodeCountHash := mh.leafHasher.Hash(utils.UInt32ToBytes(nodeCount))
	return mh.doubleSliceHash(roleHash, nodeCountHash)
}

func (mh *merkleHelper) bucketHash(bucketInfoHash, bucketEntryHash []byte) []byte {
	return mh.doubleSliceHash(bucketInfoHash, bucketEntryHash)
}

func (mh *merkleHelper) globuleInfoHash(prevCloudHash []byte, globuleID, nodeCount uint32) []byte {
	reservedHash := mh.leafHasher.Hash(utils.UInt32ToBytes(reserved))
	globuleIDHash := mh.leafHasher.Hash(utils.UInt32ToBytes(globuleID))
	nodeCountHash := mh.leafHasher.Hash(utils.UInt32ToBytes(nodeCount))

	return mh.doubleSliceHash(
		mh.doubleSliceHash(reservedHash, prevCloudHash),
		mh.doubleSliceHash(globuleIDHash, nodeCountHash),
	)
}

func (mh *merkleHelper) globuleHash(globuleInfoHash, globuleNodeRoot []byte) []byte {
	return mh.doubleSliceHash(globuleInfoHash, globuleNodeRoot)
}
