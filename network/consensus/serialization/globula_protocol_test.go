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
	"context"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"

	"github.com/stretchr/testify/require"
)

func TestEmbeddedPulsarData_SerializeTo(t *testing.T) {
	pd := EmbeddedPulsarData{
		Data: make([]byte, 10),
	}

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := pd.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 10, buf.Len())
}

func TestEmbeddedPulsarData_DeserializeFrom(t *testing.T) {
	p := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: ProtocolTypePulsar.NewBody(),
	}
	p.Header.setProtocolType(ProtocolTypePulsar)

	b := make([]byte, 64)
	_, _ = rand.Read(b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	_, err := p.SerializeTo(context.Background(), buf, digester, signer)
	require.NoError(t, err)

	pd := EmbeddedPulsarData{}
	err = pd.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	require.Equal(t, p.Header, pd.Header)
	require.Equal(t, *p.EncryptableBody.(*PulsarPacketBody), pd.PulsarPacketBody)
	require.Equal(t, p.PacketSignature, pd.PulsarSignature)
}

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
	_, _ = rand.Read(b)

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
	_, _ = rand.Read(b)

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

func TestGlobulaConsensusPacket_SerializeTo_EmptyPacket(t *testing.T) {
	p := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: &GlobulaConsensusPacketBody{},
	}
	p.Header.setProtocolType(ProtocolTypeGlobulaConsensus)
	p.Header.setPacketType(phases.PacketType(phases.PacketTypeCount)) // To emulate empty packet

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	s, err := p.SerializeTo(context.Background(), buf, digester, signer)
	require.NoError(t, err)
	require.EqualValues(t, 84, s)

	require.NotEmpty(t, p.PacketSignature)
}

func TestGlobulaConsensusPacket_DeserializeFrom(t *testing.T) {
	p1 := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: &GlobulaConsensusPacketBody{},
	}
	p1.Header.setProtocolType(ProtocolTypeGlobulaConsensus)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	_, err := p1.SerializeTo(context.Background(), buf, digester, signer)
	require.NoError(t, err)

	p2 := Packet{}

	_, err = p2.DeserializeFrom(context.Background(), buf)
	require.NoError(t, err)

	require.Equal(t, p1, p2)
}

func TestGlobulaConsensusPacketBody_Phases(t *testing.T) {
	tests := []struct {
		name       string
		packetType phases.PacketType
		size       int
	}{
		{
			"phase0",
			phases.PacketPhase0,
			88,
		},
		{
			"phase1",
			phases.PacketPhase1,
			90,
		},
		{
			"phase2",
			phases.PacketPhase2,
			89,
		},
		{
			"phase3",
			phases.PacketPhase3,
			219,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := Packet{
				Header: Header{
					SourceID:   123,
					TargetID:   456,
					ReceiverID: 789,
				},
				EncryptableBody: &GlobulaConsensusPacketBody{},
			}
			p.Header.setProtocolType(ProtocolTypeGlobulaConsensus)
			p.Header.setPacketType(test.packetType)

			buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
			s, err := p.SerializeTo(context.Background(), buf, digester, signer)
			require.NoError(t, err)
			require.EqualValues(t, test.size, s)

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			require.Equal(t, p, p2)
		})
	}
}

func TestGlobulaConsensusPacketBody_Phases_Flag0(t *testing.T) {
	pp := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: ProtocolTypePulsar.NewBody(),
	}
	pp.Header.setProtocolType(ProtocolTypePulsar)

	b := make([]byte, 64)
	_, _ = rand.Read(b)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
	_, err := pp.SerializeTo(context.Background(), buf, digester, signer)
	require.NoError(t, err)

	p := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: &GlobulaConsensusPacketBody{},
	}
	p.Header.setProtocolType(ProtocolTypeGlobulaConsensus)

	phase1p := p
	phase1p.EncryptableBody = &GlobulaConsensusPacketBody{
		PulsarPacket: EmbeddedPulsarData{
			Data: buf.Bytes(),
		},
	}

	tests := []struct {
		name       string
		packetType phases.PacketType
		size       int
		packet     Packet
	}{
		{
			"phase0",
			phases.PacketPhase0,
			216,
			phase1p,
		},
		{
			"phase1",
			phases.PacketPhase1,
			218,
			phase1p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			89,
			p,
		},
		{
			"phase3",
			phases.PacketPhase3,
			219,
			p,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.packet

			p.Header.setPacketType(test.packetType)
			p.Header.SetFlag(0)

			s, err := p.SerializeTo(context.Background(), buf, digester, signer)
			require.NoError(t, err)
			require.EqualValues(t, test.size, s)

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)
		})
	}
}

func TestGlobulaConsensusPacketBody_Phases_Flag1(t *testing.T) {
	p := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: &GlobulaConsensusPacketBody{},
	}
	p.Header.setProtocolType(ProtocolTypeGlobulaConsensus)

	phase3p := p
	phase3p.EncryptableBody = &GlobulaConsensusPacketBody{
		Vectors: NodeVectors{
			AdditionalStateVectors: make([]GlobulaStateVector, 1),
		},
	}

	tests := []struct {
		name       string
		packetType phases.PacketType
		size       int
		packet     Packet
	}{
		{
			"phase0",
			phases.PacketPhase0,
			88,
			p,
		},
		{
			"phase1",
			phases.PacketPhase1,
			90,
			p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			226,
			p,
		},
		{
			"phase3",
			phases.PacketPhase3,
			351,
			phase3p,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.packet

			p.Header.setPacketType(test.packetType)
			p.Header.SetFlag(1)

			buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
			s, err := p.SerializeTo(context.Background(), buf, digester, signer)
			require.NoError(t, err)
			require.EqualValues(t, test.size, s)

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			require.Equal(t, p, p2)
		})
	}
}

func TestGlobulaConsensusPacketBody_Phases_Flag2(t *testing.T) {
	p := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: &GlobulaConsensusPacketBody{},
	}
	p.Header.setProtocolType(ProtocolTypeGlobulaConsensus)

	phase3p := p
	phase3p.EncryptableBody = &GlobulaConsensusPacketBody{
		Vectors: NodeVectors{
			AdditionalStateVectors: make([]GlobulaStateVector, 2),
		},
	}

	tests := []struct {
		name       string
		packetType phases.PacketType
		size       int
		packet     Packet
	}{
		{
			"phase0",
			phases.PacketPhase0,
			88,
			p,
		},
		{
			"phase1",
			phases.PacketPhase1,
			441,
			p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			440,
			p,
		},
		{
			"phase3",
			phases.PacketPhase3,
			483,
			phase3p,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.packet

			p.Header.setPacketType(test.packetType)
			p.Header.SetFlag(2)

			buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
			s, err := p.SerializeTo(context.Background(), buf, digester, signer)
			require.NoError(t, err)
			require.EqualValues(t, test.size, s)

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			require.Equal(t, p, p2)
		})
	}
}

func TestGlobulaConsensusPacketBody_Phases_Flag12(t *testing.T) {
	p := Packet{
		Header: Header{
			SourceID:   123,
			TargetID:   456,
			ReceiverID: 789,
		},
		EncryptableBody: &GlobulaConsensusPacketBody{},
	}
	p.Header.setProtocolType(ProtocolTypeGlobulaConsensus)

	phase3p := p
	phase3p.EncryptableBody = &GlobulaConsensusPacketBody{
		Vectors: NodeVectors{
			AdditionalStateVectors: make([]GlobulaStateVector, 3),
		},
	}

	tests := []struct {
		name       string
		packetType phases.PacketType
		size       int
		packet     Packet
	}{
		{
			"phase0",
			phases.PacketPhase0,
			88,
			p,
		},
		{
			"phase1",
			phases.PacketPhase1,
			505,
			p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			504,
			p,
		},
		{
			"phase3",
			phases.PacketPhase3,
			615,
			phase3p,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := test.packet

			p.Header.setPacketType(test.packetType)
			p.Header.SetFlag(1)
			p.Header.SetFlag(2)

			buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))
			s, err := p.SerializeTo(context.Background(), buf, digester, signer)
			require.NoError(t, err)
			require.EqualValues(t, test.size, s)

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			require.Equal(t, p, p2)
		})
	}
}
