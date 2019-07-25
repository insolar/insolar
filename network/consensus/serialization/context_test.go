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

package serialization

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacketContext_InContext(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.True(t, ctx.InContext(NoContext))
	require.False(t, ctx.InContext(ContextMembershipAnnouncement))
	require.False(t, ctx.InContext(ContextNeighbourAnnouncement))

	ctx.fieldContext = ContextMembershipAnnouncement
	require.True(t, ctx.InContext(ContextMembershipAnnouncement))

	ctx.fieldContext = ContextNeighbourAnnouncement
	require.True(t, ctx.InContext(ContextNeighbourAnnouncement))
}

func TestPacketContext_SetInContext(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.True(t, ctx.InContext(NoContext))

	ctx.SetInContext(ContextMembershipAnnouncement)
	require.True(t, ctx.InContext(ContextMembershipAnnouncement))

	ctx.SetInContext(ContextNeighbourAnnouncement)
	require.True(t, ctx.InContext(ContextNeighbourAnnouncement))
}

func TestPacketContext_GetNeighbourNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	ctx.neighbourNodeID = 123
	require.EqualValues(t, 123, ctx.GetNeighbourNodeID())
}

func TestPacketContext_GetNeighbourNodeID_Panics(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.Panics(t, func() { ctx.GetNeighbourNodeID() })
}

func TestPacketContext_SetNeighbourNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	ctx.SetNeighbourNodeID(123)
	require.EqualValues(t, 123, ctx.GetNeighbourNodeID())
}

func TestPacketContext_GetAnnouncedJoinerNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.EqualValues(t, 0, ctx.GetAnnouncedJoinerNodeID())

	ctx.announcedJoinerNodeID = 123
	require.EqualValues(t, 123, ctx.GetAnnouncedJoinerNodeID())
}

func TestPacketContext_SetAnnouncedJoinerNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.EqualValues(t, 0, ctx.GetAnnouncedJoinerNodeID())

	ctx.SetAnnouncedJoinerNodeID(123)
	require.EqualValues(t, 123, ctx.GetAnnouncedJoinerNodeID())
}
