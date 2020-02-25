// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
