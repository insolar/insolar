// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package tests

import (
	"context"
	"math/rand"
	"time"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

func NewEmuUpstreamPulseController(ctx context.Context, nshDelay time.Duration) *EmuUpstreamPulseController {
	return &EmuUpstreamPulseController{ctx: ctx, nshDelay: nshDelay}
}

var _ api.UpstreamController = &EmuUpstreamPulseController{}

type EmuUpstreamPulseController struct {
	ctx      context.Context
	nshDelay time.Duration
}

func (*EmuUpstreamPulseController) ConsensusAborted() {
}

func (r *EmuUpstreamPulseController) PreparePulseChange(report api.UpstreamReport, c chan<- api.UpstreamState) {
	fn := func() {
		nsh := NewEmuNodeStateHash(rand.Uint64())
		c <- api.UpstreamState{NodeState: nsh}
		close(c)
	}
	if r.nshDelay == 0 {
		fn()
	} else {
		time.AfterFunc(r.nshDelay, fn)
	}
}

func (*EmuUpstreamPulseController) CommitPulseChange(report api.UpstreamReport, pd pulse.Data, activeCensus census.Operational) {
}

func (*EmuUpstreamPulseController) CancelPulseChange() {
}

func (*EmuUpstreamPulseController) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
}

func NewEmuNodeStateHash(v uint64) *EmuNodeStateHash {
	return &EmuNodeStateHash{Bits64: longbits.NewBits64(v)}
}

var _ proofs.NodeStateHash = &EmuNodeStateHash{}

type EmuNodeStateHash struct {
	longbits.Bits64
}

func (r *EmuNodeStateHash) CopyOfDigest() cryptkit.Digest {
	return cryptkit.NewDigest(&r.Bits64, r.GetDigestMethod())
}

func (r *EmuNodeStateHash) SignWith(signer cryptkit.DigestSigner) cryptkit.SignedDigestHolder {
	d := r.CopyOfDigest()
	return d.SignWith(signer)
}

func (r *EmuNodeStateHash) GetDigestMethod() cryptkit.DigestMethod {
	return "uint64"
}

func (r *EmuNodeStateHash) Equals(o cryptkit.DigestHolder) bool {
	return longbits.EqualFixedLenWriterTo(r, o)
}
