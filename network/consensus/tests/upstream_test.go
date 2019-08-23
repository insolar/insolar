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
