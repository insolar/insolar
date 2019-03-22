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
