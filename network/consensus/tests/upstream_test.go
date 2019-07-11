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
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"math/rand"
	"time"
)

func NewEmuUpstreamPulseController(ctx context.Context, nshDelay time.Duration) *EmuUpstreamPulseController {
	return &EmuUpstreamPulseController{ctx: ctx, nshDelay: nshDelay}
}

var _ api.UpstreamPulseController = &EmuUpstreamPulseController{}

type EmuUpstreamPulseController struct {
	ctx      context.Context
	nshDelay time.Duration
}

func (r *EmuUpstreamPulseController) PreparePulseChange(report api.MembershipUpstreamReport) <-chan gcp_types.NodeStateHash {
	c := make(chan gcp_types.NodeStateHash, 1)
	nsh := NewEmuNodeStateHash(rand.Uint64())
	if r.nshDelay == 0 {
		c <- &nsh
		close(c)
	} else {
		time.AfterFunc(r.nshDelay, func() {
			c <- &nsh
			close(c)
		})
	}
	return c
}

func (*EmuUpstreamPulseController) CommitPulseChange(report api.MembershipUpstreamReport, pd pulse_data.PulseData, activeCensus api.OperationalCensus) {
}

func (*EmuUpstreamPulseController) CancelPulseChange() {
}

func (*EmuUpstreamPulseController) ConsensusFinished(report api.MembershipUpstreamReport, expectedCensus api.OperationalCensus) {
}

func NewEmuNodeStateHash(v uint64) EmuNodeStateHash {
	return EmuNodeStateHash{Bits64: long_bits.NewBits64(v)}
}

var _ gcp_types.NodeStateHash = &EmuNodeStateHash{}

type EmuNodeStateHash struct {
	long_bits.Bits64
}

func (r *EmuNodeStateHash) CopyOfDigest() cryptography_containers.Digest {
	return cryptography_containers.NewDigest(&r.Bits64, r.GetDigestMethod())
}

func (r *EmuNodeStateHash) SignWith(signer cryptography_containers.DigestSigner) cryptography_containers.SignedDigest {
	d := r.CopyOfDigest()
	return d.SignWith(signer)
}

func (r *EmuNodeStateHash) GetDigestMethod() cryptography_containers.DigestMethod {
	return "uint64"
}

func (r *EmuNodeStateHash) Equals(o cryptography_containers.DigestHolder) bool {
	return long_bits.EqualFixedLenWriterTo(r, o)
}
