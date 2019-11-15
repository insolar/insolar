///
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
///

package merkler

import (
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

var _ cryptkit.SequenceDigester = &StackedCalculator{}

func NewStackedCalculator(digester cryptkit.PairDigester, unbalancedStub cryptkit.Digest) StackedCalculator {
	if digester == nil {
		panic("illegal value")
	}
	return StackedCalculator{digester: digester, unbalancedStub: unbalancedStub}
}

// A calculator to do streaming calculation of Merkle-tree by using provided PairDigester.
// The resulting merkle will have same depth for all branches except for the rightmost branch.
//
// When unbalancedStub == nil, then FinishSequence() will create the rightmost branch by recursively
// applying the same rule - all sub-branches will have same depth except for the rightmost sub-branch.
//
// When unbalancedStub != nil, then FinishSequence() will create a perfect tree by using unbalancedStub
// once per level when a value for the rightmost sub-branch is missing.
//
// When AddNext() was never called then FinishSequence() will return a non-nil unbalancedStub otherwise will panic.
//
// Complexity (n - a number of added hashes):
//  - AddNext() is O(1), it does only upto 2 calls to PairDigester.DigestPair()
//  - FinishSequence() is O(log n), it does log(n) calls to PairDigester.DigestPair()
//  - ForkSequence() is O(log n), but only copies memory
//  - Memory is O(log n)
//
type StackedCalculator struct {
	digester       cryptkit.PairDigester
	unbalancedStub cryptkit.Digest
	prevAdded      longbits.FoldableReader
	count          uint
	treeLevels     []treeLevel
	finished       bool
}

type treeLevel struct {
	digest0 cryptkit.Digest
	digest1 cryptkit.Digest
}

func (p *StackedCalculator) GetDigestMethod() cryptkit.DigestMethod {
	return p.digester.GetDigestMethod() + "/merkle"
}

func (p *StackedCalculator) GetDigestSize() int {
	return p.digester.GetDigestSize()
}

func (p *StackedCalculator) AddNext(addDigest longbits.FoldableReader) {

	if p.finished || p.digester == nil {
		panic("illegal state")
	}

	/*
		Here we use position of grey-code transition bit to choose the level that requires pair-hashing
		Tree level is counted from a leaf to the root, leaf level is considered as -1
	*/

	p.count++
	pairPosition := args.GreyIncBit(p.count) // level + 1

	var bottomDigest cryptkit.Digest

	if pairPosition == 0 {
		// Level -1 (leafs) is special - it only stores a previous value
		bottomDigest = p.digester.DigestPair(p.prevAdded, addDigest)
		p.prevAdded = nil
	} else {
		if p.prevAdded != nil {
			panic("illegal state")
		}
		p.prevAdded = addDigest

		if int(pairPosition) > len(p.treeLevels) {
			return
		}
		pairLevel := &p.treeLevels[pairPosition-1]
		bottomDigest = p.digester.DigestPair(pairLevel.digest0, pairLevel.digest1)
		pairLevel.digest0, pairLevel.digest1 = cryptkit.Digest{}, cryptkit.Digest{}
	}

	if int(pairPosition) == len(p.treeLevels) {
		p.treeLevels = append(p.treeLevels, treeLevel{digest0: bottomDigest})
		return
	}

	var d *cryptkit.Digest
	if p.count&(uint(2)<<pairPosition) != 0 {
		d = &p.treeLevels[pairPosition].digest0
	} else {
		d = &p.treeLevels[pairPosition].digest1
	}
	if !d.IsEmpty() {
		panic("illegal state")
	}
	*d = bottomDigest
}

func (p *StackedCalculator) ForkSequence() cryptkit.SequenceDigester {

	if p.finished || p.digester == nil {
		panic("illegal state")
	}

	cp := *p
	cp.treeLevels = append(make([]treeLevel, 0, cap(p.treeLevels)), p.treeLevels...)
	return &cp
}

func (p *StackedCalculator) Count() int {
	return int(p.count)
}

func (p *StackedCalculator) FinishSequence() cryptkit.Digest {

	if p.finished || p.digester == nil {
		panic("illegal state")
	}
	p.finished = true

	hasStub := !p.unbalancedStub.IsEmpty()
	if p.count == 0 {
		if hasStub {
			return p.unbalancedStub
		}
		panic("illegal state - empty")
	}

	var bottomDigest cryptkit.Digest
	if p.prevAdded != nil {
		if hasStub {
			bottomDigest = p.digester.DigestPair(p.prevAdded, p.unbalancedStub)
		} else {
			bottomDigest = cryptkit.NewDigest(p.prevAdded, p.digester.GetDigestMethod())
		}
	}

	for i := 0; i < len(p.treeLevels); i++ { // DONT USE range as treeLevels can be appended inside the loop!
		curLevel := &p.treeLevels[i]

		switch {
		case !curLevel.digest1.IsEmpty():
			// both are present
			if curLevel.digest0.IsEmpty() {
				panic("illegal state")
			}
			levelDigest := p.digester.DigestPair(curLevel.digest0, curLevel.digest1)

			if bottomDigest.IsEmpty() {
				bottomDigest = levelDigest
				continue
			}

			if i+1 == len(p.treeLevels) {
				p.treeLevels = append(p.treeLevels, treeLevel{digest0: levelDigest})
			} else {
				nextLevel := &p.treeLevels[i+1]
				switch {
				case nextLevel.digest0.IsEmpty():
					if !nextLevel.digest1.IsEmpty() {
						panic("illegal state")
					}
					nextLevel.digest0 = levelDigest
				case nextLevel.digest1.IsEmpty():
					nextLevel.digest1 = levelDigest
				default:
					panic("illegal state - only one dual-hash can be present in the stack")
				}
			}
			if hasStub {
				bottomDigest = p.digester.DigestPair(bottomDigest, p.unbalancedStub)
			}
			// or leave as is
		case !curLevel.digest0.IsEmpty():
			switch {
			case !bottomDigest.IsEmpty():
				bottomDigest = p.digester.DigestPair(curLevel.digest0, bottomDigest)
			case hasStub:
				bottomDigest = p.digester.DigestPair(curLevel.digest0, p.unbalancedStub)
			default:
				bottomDigest = curLevel.digest0
			}
		case !bottomDigest.IsEmpty() && hasStub:
			bottomDigest = p.digester.DigestPair(bottomDigest, p.unbalancedStub)
		}
	}

	return bottomDigest
}
