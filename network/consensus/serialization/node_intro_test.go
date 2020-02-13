// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package serialization

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

func TestNodeBriefIntro_getPrimaryRole(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, member.PrimaryRoleInactive, ni.GetPrimaryRole())

	ni.PrimaryRoleAndFlags = 1
	require.Equal(t, member.PrimaryRoleNeutral, ni.GetPrimaryRole())

	ni.PrimaryRoleAndFlags = 2
	require.Equal(t, member.PrimaryRoleHeavyMaterial, ni.GetPrimaryRole())
}

func TestNodeBriefIntro_setPrimaryRole(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, member.PrimaryRoleInactive, ni.GetPrimaryRole())

	ni.SetPrimaryRole(member.PrimaryRoleVirtual)
	require.Equal(t, member.PrimaryRoleVirtual, ni.GetPrimaryRole())
}

func TestNodeBriefIntro_setPrimaryRole_Panic(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Panics(t, func() { ni.SetPrimaryRole(primaryRoleMax + 1) })
}

func TestNodeBriefIntro_getAddrMode(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, endpoints.IPEndpoint, ni.GetAddrMode())

	ni.PrimaryRoleAndFlags = 64 // 0b01000000
	require.Equal(t, endpoints.NameEndpoint, ni.GetAddrMode())
}

func TestNodeBriefIntro_setAddrMode(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Equal(t, endpoints.IPEndpoint, ni.GetAddrMode())

	ni.SetAddrMode(endpoints.RelayEndpoint)
	require.Equal(t, endpoints.RelayEndpoint, ni.GetAddrMode())
}

func TestNodeBriefIntro_setAddrMode_Panic(t *testing.T) {
	ni := NodeBriefIntro{}

	require.Panics(t, func() { ni.SetAddrMode(addrModeMax + 1) })
}

func TestNodeBriefIntro_SerializeTo(t *testing.T) {
	ni := NodeBriefIntro{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := ni.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 149, buf.Len())
}

func TestNodeBriefIntro_DeserializeFrom(t *testing.T) {
	ni1 := NodeBriefIntro{
		PrimaryRoleAndFlags: 64,
		SpecialRoles:        member.SpecialRoleDiscovery,
		StartPower:          10,
	}

	b := make([]byte, 64)
	_, _ = rand.Read(b)

	copy(ni1.JoinerSignature[:], b)
	copy(ni1.NodePK[:], b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ni1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ni2 := NodeBriefIntro{}
	err = ni2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	ni2.JoinerData = nil // Skip in comparison
	require.Equal(t, ni1, ni2)
}

func TestNodeBriefIntro_DeserializeFrom_NoShortID(t *testing.T) {
	ni1 := NodeBriefIntro{
		ShortID:             123,
		PrimaryRoleAndFlags: 64,
		SpecialRoles:        member.SpecialRoleDiscovery,
		StartPower:          10,
	}

	b := make([]byte, 64)
	_, _ = rand.Read(b)

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

	ni3.JoinerData = nil // Skip in comparison
	require.Equal(t, ni2, ni3)

	require.EqualValues(t, 0, ni3.ShortID)
}

func TestNodeFullIntro_SerializeTo(t *testing.T) {
	ni := NodeFullIntro{}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := ni.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 235, buf.Len())
}

func TestNodeFullIntro_DeserializeFrom(t *testing.T) {
	ni1 := NodeFullIntro{
		NodeBriefIntro: NodeBriefIntro{
			PrimaryRoleAndFlags: 64,
			SpecialRoles:        member.SpecialRoleDiscovery,
			StartPower:          10,
		},
	}

	b := make([]byte, 64)
	_, _ = rand.Read(b)

	copy(ni1.JoinerSignature[:], b)
	copy(ni1.NodePK[:], b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	err := ni1.SerializeTo(nil, buf)
	require.NoError(t, err)

	ni2 := NodeFullIntro{}
	err = ni2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	ni2.JoinerData = nil // Skip in comparison
	require.Equal(t, ni1, ni2)
}

func TestNodeFullIntro_DeserializeFrom_NoShortID(t *testing.T) {
	ni1 := NodeFullIntro{
		NodeBriefIntro: NodeBriefIntro{
			PrimaryRoleAndFlags: 64,
			SpecialRoles:        member.SpecialRoleDiscovery,
			StartPower:          10,
		},
	}

	b := make([]byte, 64)
	_, _ = rand.Read(b)

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

	ni3.JoinerData = nil // Skip in comparison
	require.Equal(t, ni2, ni3)
	require.EqualValues(t, 0, ni3.ShortID)
}

func TestNodeFullIntro_DeserializeFrom_Slices(t *testing.T) {
	ni1 := NodeFullIntro{
		NodeBriefIntro: NodeBriefIntro{
			PrimaryRoleAndFlags: 64,
			SpecialRoles:        member.SpecialRoleDiscovery,
			StartPower:          10,
		},
		NodeExtendedIntro: NodeExtendedIntro{
			EndpointLen:    2,
			ExtraEndpoints: make([]uint16, 2),
			ProofLen:       2,
			NodeRefProof:   make([]longbits.Bits512, 2),
		},
	}

	b := make([]byte, 64)
	_, _ = rand.Read(b)

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

	ni3.JoinerData = nil // Skip in comparison
	require.Equal(t, ni2, ni3)
	require.EqualValues(t, 0, ni3.ShortID)
}
