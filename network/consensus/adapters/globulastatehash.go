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

package adapters

import (
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type SequenceDigester struct {
	dataDigester cryptkit.DataDigester
	state        uint64
}

func NewSequenceDigester(dataDigester cryptkit.DataDigester) *SequenceDigester {
	return &SequenceDigester{
		dataDigester: dataDigester,
	}
}

func (d *SequenceDigester) AddNext(digest longbits.FoldableReader) {
	d.addNext(digest.FoldToUint64())
}

func (d *SequenceDigester) addNext(state uint64) {
	d.state ^= state
}

func (d *SequenceDigester) FinishSequence() cryptkit.Digest {
	bits64 := longbits.NewBits64(d.state)
	return d.dataDigester.GetDigestOf(&bits64)
}

func (d *SequenceDigester) GetDigestMethod() cryptkit.DigestMethod {
	return d.dataDigester.GetDigestMethod()
}

func (d *SequenceDigester) ForkSequence() cryptkit.SequenceDigester {
	return &SequenceDigester{
		dataDigester: d.dataDigester,
		state:        d.state,
	}
}

type StateDigester struct {
	sequenceDigester *SequenceDigester
	defaultDigest    longbits.FoldableReader
}

func NewStateDigester(sequenceDigester *SequenceDigester) *StateDigester {
	return &StateDigester{
		sequenceDigester: sequenceDigester,
		defaultDigest:    &longbits.Bits512{},
	}
}

func (d *StateDigester) AddNext(digest longbits.FoldableReader, fullRank member.FullRank) {
	if digest == nil {
		d.sequenceDigester.AddNext(d.defaultDigest)
	} else {
		d.sequenceDigester.AddNext(digest)
		d.sequenceDigester.addNext(uint64(fullRank.AsMembershipRank(member.MaxNodeIndex)))
	}
}

func (d *StateDigester) GetDigestMethod() cryptkit.DigestMethod {
	return d.sequenceDigester.GetDigestMethod()
}

func (d *StateDigester) ForkSequence() transport.StateDigester {
	return &StateDigester{
		sequenceDigester: d.sequenceDigester.ForkSequence().(*SequenceDigester),
		defaultDigest:    d.defaultDigest,
	}
}

func (d *StateDigester) FinishSequence() cryptkit.Digest {
	return d.sequenceDigester.FinishSequence()
}
