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

package node

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/commoncrypto"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSnapshotEncodeDecode(t *testing.T) {

	ks := commoncrypto.NewKeyProcessor()
	p1, err := ks.GeneratePrivateKey()
	p2, err := ks.GeneratePrivateKey()
	assert.NoError(t, err)

	n1 := newMutableNode(testutils.RandomRef(), insolar.StaticRoleVirtual, ks.ExtractPublicKey(p1), insolar.NodeReady, "127.0.0.1:22", "ver2")
	n2 := newMutableNode(testutils.RandomRef(), insolar.StaticRoleHeavyMaterial, ks.ExtractPublicKey(p2), insolar.NodeLeaving, "127.0.0.1:33", "ver5")

	s := Snapshot{}
	s.pulse = 22
	s.state = insolar.CompleteNetworkState
	s.nodeList[ListLeaving] = []insolar.NetworkNode{n1, n2}
	s.nodeList[ListJoiner] = []insolar.NetworkNode{n2}

	buff, err := s.Encode()
	assert.NoError(t, err)
	assert.NotEmptyf(t, buff, "should not be empty")

	s2 := Snapshot{}
	err = s2.Decode(buff)
	assert.NoError(t, err)
	assert.True(t, s.Equal(&s2))
}
