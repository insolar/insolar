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
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloudIntro_SerializeTo(t *testing.T) {
	ci := CloudIntro{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := ci.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 128, buf.Len())
}

func TestCloudIntro_DeserializeFrom(t *testing.T) {
	ci1 := CloudIntro{}

	b := make([]byte, 64)
	rand.Read(b)

	copy(ci1.CloudIdentity[:], b)
	copy(ci1.LastCloudStateHash[:], b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ci1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ci2 := CloudIntro{}
	err = ci2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, ci1, ci2)
}

func TestCompactGlobulaNodeState_SerializeTo(t *testing.T) {
	s := CompactGlobulaNodeState{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := s.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 128, buf.Len())
}

func TestCompactGlobulaNodeState_DeserializeFrom(t *testing.T) {
	s1 := CompactGlobulaNodeState{}

	b := make([]byte, 64)
	rand.Read(b)

	copy(s1.NodeStateHash[:], b)
	copy(s1.GlobulaNodeStateSignature[:], b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := s1.SerializeTo(nil, buf)
	require.NoError(t, err)

	s2 := CompactGlobulaNodeState{}
	err = s2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, s1, s2)
}

func TestLeaveAnnouncement_SerializeTo(t *testing.T) {
	la := LeaveAnnouncement{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := la.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 4, buf.Len())
}

func TestLeaveAnnouncement_DeserializeFrom(t *testing.T) {
	la1 := LeaveAnnouncement{
		LeaveReason: 123,
	}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := la1.SerializeTo(nil, buf)
	require.NoError(t, err)

	la2 := LeaveAnnouncement{}
	err = la2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, la1, la2)
}
