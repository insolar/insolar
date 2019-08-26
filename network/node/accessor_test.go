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
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/assert"
)

func TestAccessor(t *testing.T) {
	t.Skip("FIXME")

	node := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeReady, "127.0.0.1:0", "")

	node2 := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeJoining, "127.0.0.1:0", "")
	node2.SetShortID(11)

	node3 := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeLeaving, "127.0.0.1:0", "")
	node3.SetShortID(10)

	node4 := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeUndefined, "127.0.0.1:0", "")

	snapshot := NewSnapshot(insolar.FirstPulseNumber, []insolar.NetworkNode{node, node2, node3, node4})
	accessor := NewAccessor(snapshot)
	assert.Equal(t, 4, len(accessor.GetActiveNodes()))
	assert.Equal(t, 1, len(accessor.GetWorkingNodes()))
	assert.NotNil(t, accessor.GetWorkingNode(node.ID()))
	assert.Nil(t, accessor.GetWorkingNode(node2.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node2.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node3.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node4.ID()))

	assert.NotNil(t, accessor.GetActiveNodeByShortID(10))
	assert.NotNil(t, accessor.GetActiveNodeByShortID(11))
	assert.Nil(t, accessor.GetActiveNodeByShortID(12))
}
