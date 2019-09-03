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

package storage

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStorage()
	startPulse := *insolar.GenesisPulse

	ks := platformpolicy.NewKeyProcessor()
	p1, err := ks.GeneratePrivateKey()
	assert.NoError(t, err)
	n := node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, ks.ExtractPublicKey(p1), "127.0.0.1:22", "ver2")
	nodes := []insolar.NetworkNode{n}

	for i := 0; i < entriesCount+2; i++ {
		p := startPulse
		p.PulseNumber += insolar.PulseNumber(i)

		snap := node.NewSnapshot(p.PulseNumber, nodes)
		err = s.Append(p.PulseNumber, snap)
		assert.NoError(t, err)

		err = s.AppendPulse(ctx, p)
		assert.NoError(t, err)

		p1, err := s.GetLatestPulse(ctx)
		assert.NoError(t, err)
		assert.Equal(t, p, p1)

		snap1, err := s.ForPulseNumber(p1.PulseNumber)
		assert.NoError(t, err)
		assert.True(t, snap1.Equal(snap), "snapshots should be equal")
	}

	// first pulse and snapshot should be truncated
	assert.Len(t, s.entries, entriesCount)
	assert.Len(t, s.snapshotEntries, entriesCount)

	p2, err := s.GetPulse(ctx, startPulse.PulseNumber)
	assert.EqualError(t, err, ErrNotFound.Error())
	assert.Equal(t, p2, *insolar.GenesisPulse)

	snap2, err := s.ForPulseNumber(startPulse.PulseNumber)
	assert.EqualError(t, err, ErrNotFound.Error())
	assert.Nil(t, snap2)

}
