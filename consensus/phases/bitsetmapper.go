/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package phases

import (
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
)

type BitsetMapper struct {
	length     int
	origin     insolar.NetworkNode
	refToIndex map[insolar.Reference]int
	indexToRef map[int]insolar.Reference
}

func (bm *BitsetMapper) AddNode(node insolar.NetworkNode, bitsetIndex uint16) {
	bm.addNode(node, int(bitsetIndex))
}

func (bm *BitsetMapper) addNode(node insolar.NetworkNode, index int) {
	bm.indexToRef[index] = node.ID()
	bm.refToIndex[node.ID()] = index
}

func (bm *BitsetMapper) IndexToRef(index int) (insolar.Reference, error) {
	if index < 0 || index >= bm.length {
		return insolar.Reference{}, packets.ErrBitSetOutOfRange
	}
	result, ok := bm.indexToRef[index]
	if !ok {
		return insolar.Reference{}, packets.ErrBitSetNodeIsMissing
	}
	return result, nil
}

func (bm *BitsetMapper) RefToIndex(nodeID insolar.Reference) (int, error) {
	index, ok := bm.refToIndex[nodeID]
	if !ok {
		return 0, packets.ErrBitSetIncorrectNode
	}
	return index, nil
}

func (bm *BitsetMapper) Length() int {
	return bm.length
}
