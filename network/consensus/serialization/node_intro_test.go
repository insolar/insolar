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
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeBriefIntro_getPrimaryRole(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, gcp_types.PrimaryRoleInactive, ni.getPrimaryRole())

	ni.PrimaryRoleAndFlags = 1
	require.Equal(t, gcp_types.PrimaryRoleNeutral, ni.getPrimaryRole())

	ni.PrimaryRoleAndFlags = 2
	require.Equal(t, gcp_types.PrimaryRoleHeavyMaterial, ni.getPrimaryRole())
}

func TestNodeBriefIntro_setPrimaryRole(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, gcp_types.PrimaryRoleInactive, ni.getPrimaryRole())

	ni.setPrimaryRole(gcp_types.PrimaryRoleVirtual)
	require.Equal(t, gcp_types.PrimaryRoleVirtual, ni.getPrimaryRole())
}

func TestNodeBriefIntro_setPrimaryRole_Panic(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Panics(t, func() { ni.setPrimaryRole(primaryRoleMax + 1) })
}

func TestNodeBriefIntro_getAddrMode(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, endpoints.IPEndpoint, ni.getAddrMode())

	ni.PrimaryRoleAndFlags = 64 // 0b01000000
	require.Equal(t, endpoints.NameEndpoint, ni.getAddrMode())
}

func TestNodeBriefIntro_setAddrMode(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, endpoints.IPEndpoint, ni.getAddrMode())

	ni.setAddrMode(endpoints.RelayEndpoint)
	require.Equal(t, endpoints.RelayEndpoint, ni.getAddrMode())
}

func TestNodeBriefIntro_setAddrMode_Panic(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Panics(t, func() { ni.setAddrMode(addrModeMax + 1) })
}

func TestNodeBriefIntro_SerializeTo(t *testing.T) {
	ni := NodeBriefIntro{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := ni.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 137, buf.Len())
}

func TestNodeBriefIntro_DeserializeFrom(t *testing.T) {
	ni1 := NodeBriefIntro{
		PrimaryRoleAndFlags: 64,
		SpecialRoles:        gcp_types.SpecialRoleDiscovery,
		StartPower:          10,
		BasePort:            1400,
		PrimaryIPv4:         123123412,
	}

	b := make([]byte, 64)
	rand.Read(b)

	copy(ni1.JoinerSignature[:], b)
	copy(ni1.NodePK[:], b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ni1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ni2 := NodeBriefIntro{}
	err = ni2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, ni1, ni2)
}

func TestNodeBriefIntro_DeserializeFrom_NoShortID(t *testing.T) {
	ni1 := NodeBriefIntro{
		ShortID:             123,
		PrimaryRoleAndFlags: 64,
		SpecialRoles:        gcp_types.SpecialRoleDiscovery,
		StartPower:          10,
		BasePort:            1400,
		PrimaryIPv4:         123123412,
	}

	b := make([]byte, 64)
	rand.Read(b)

	copy(ni1.JoinerSignature[:], b)
	copy(ni1.NodePK[:], b)

	ni2 := ni1 // Copy and reset short ic
	ni2.ShortID = 0

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ni1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ni3 := NodeBriefIntro{}
	err = ni3.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, ni2, ni3)
	require.EqualValues(t, 0, ni3.ShortID)
}

func TestNodeFullIntro_SerializeTo(t *testing.T) {
	ni := NodeFullIntro{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := ni.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 223, buf.Len())
}

func TestNodeFullIntro_DeserializeFrom(t *testing.T) {
	ni1 := NodeFullIntro{
		NodeBriefIntro: NodeBriefIntro{
			PrimaryRoleAndFlags: 64,
			SpecialRoles:        gcp_types.SpecialRoleDiscovery,
			StartPower:          10,
			BasePort:            1400,
			PrimaryIPv4:         123123412,
		},
	}

	b := make([]byte, 64)
	rand.Read(b)

	copy(ni1.JoinerSignature[:], b)
	copy(ni1.NodePK[:], b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ni1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ni2 := NodeFullIntro{}
	err = ni2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, ni1, ni2)
}

func TestNodeFullIntro_DeserializeFrom_NoShortID(t *testing.T) {
	ni1 := NodeFullIntro{
		NodeBriefIntro: NodeBriefIntro{
			PrimaryRoleAndFlags: 64,
			SpecialRoles:        gcp_types.SpecialRoleDiscovery,
			StartPower:          10,
			BasePort:            1400,
			PrimaryIPv4:         123123412,
		},
	}

	b := make([]byte, 64)
	rand.Read(b)

	copy(ni1.JoinerSignature[:], b)
	copy(ni1.NodePK[:], b)

	ni2 := ni1 // Copy and reset short ic
	ni2.ShortID = 0

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ni1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ni3 := NodeFullIntro{}
	err = ni3.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, ni2, ni3)
	require.EqualValues(t, 0, ni3.ShortID)
}

func TestNodeFullIntro_DeserializeFrom_Slices(t *testing.T) {
	ni1 := NodeFullIntro{
		NodeBriefIntro: NodeBriefIntro{
			PrimaryRoleAndFlags: 64,
			SpecialRoles:        gcp_types.SpecialRoleDiscovery,
			StartPower:          10,
			BasePort:            1400,
			PrimaryIPv4:         123123412,
		},
		EndpointLen:    2,
		ExtraEndpoints: make([]uint16, 2),
		ProofLen:       2,
		NodeRefProof:   make([]long_bits.Bits512, 2),
	}

	b := make([]byte, 64)
	rand.Read(b)

	copy(ni1.JoinerSignature[:], b)
	copy(ni1.NodePK[:], b)

	ni2 := ni1 // Copy and reset short ic
	ni2.ShortID = 0

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ni1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ni3 := NodeFullIntro{}
	err = ni3.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, ni2, ni3)
	require.EqualValues(t, 0, ni3.ShortID)
}
