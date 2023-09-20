package serialization

import (
	"bytes"
	"context"
	"crypto/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/pulse"
)

func TestEmbeddedPulsarData_SerializeTo(t *testing.T) {
	pd := EmbeddedPulsarData{}
	pd.setData(make([]byte, 10))

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	err := pd.SerializeTo(nil, buf)
	require.NoError(t, err)
	require.Equal(t, 12, buf.Len())
}

func TestEmbeddedPulsarData_DeserializeFrom(t *testing.T) {
	data := pulse.NewPulsarData(100000, 10, 10, *longbits.NewBits256FromBytes(make([]byte, 32)))

	pu := adapters.NewPulse(data)
	ph, err := host.NewHost("127.0.0.1:1")
	require.NoError(t, err)
	th, err := net.ResolveTCPAddr("tcp", "127.0.0.1:2")
	require.NoError(t, err)
	pp := pulsenetwork.NewPulsePacketWithTrace(context.Background(), &pu, ph, th, 0)

	bs, err := packet.SerializePacket(pp)
	require.NoError(t, err)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

	pd1 := EmbeddedPulsarData{}
	pd1.setData(bs)
	err = pd1.SerializeTo(nil, buf)
	require.NoError(t, err)

	pd2 := EmbeddedPulsarData{}
	err = pd2.DeserializeFrom(nil, buf)
	require.NoError(t, err)

	// require.Equal(t, p.Header, pd.Header)
	require.Equal(t, pd1.Size, pd2.Size)
	require.Equal(t, pd1.Data, pd2.Data)
	// require.Equal(t, p.PacketSignature, pd.PulsarSignature)
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
	copy(s1.NodeStateHashSignature[:], b)

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

	n, err := p1.SerializeTo(context.Background(), buf, digester, signer)
	require.EqualValues(t, n, p1.Header.getPayloadLength())
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
			91,
		},
		{
			"phase2",
			phases.PacketPhase2,
			90,
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
			require.EqualValues(t, s, p.Header.getPayloadLength())

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			require.Equal(t, p, p2)
		})
	}
}

func TestGlobulaConsensusPacketBody_Phases_Flag0(t *testing.T) {
	data := pulse.NewPulsarData(100000, 10, 10, *longbits.NewBits256FromBytes(make([]byte, 32)))

	pu := adapters.NewPulse(data)
	ph, err := host.NewHost("127.0.0.1:1")
	require.NoError(t, err)
	th, err := net.ResolveTCPAddr("tcp", "127.0.0.1:2")
	require.NoError(t, err)
	pp := pulsenetwork.NewPulsePacketWithTrace(context.Background(), &pu, ph, th, 0)

	bs, err := packet.SerializePacket(pp)
	require.NoError(t, err)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

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
	phase1p.EncryptableBody = &GlobulaConsensusPacketBody{}
	phase1p.EncryptableBody.(*GlobulaConsensusPacketBody).PulsarPacket.setData(bs)

	tests := []struct {
		name       string
		packetType phases.PacketType
		size       int64
		packet     Packet
	}{
		{
			"phase0",
			phases.PacketPhase0,
			401,
			phase1p,
		},
		{
			"phase1",
			phases.PacketPhase1,
			404,
			phase1p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			90,
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
			require.EqualValues(t, s, p.Header.getPayloadLength())

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)
		})
	}
}

func TestGlobulaConsensusPacketBody_Phases_Flag0Reset(t *testing.T) {
	data := pulse.NewPulsarData(100000, 10, 10, *longbits.NewBits256FromBytes(make([]byte, 32)))

	pu := adapters.NewPulse(data)
	ph, err := host.NewHost("127.0.0.1:1")
	require.NoError(t, err)
	th, err := net.ResolveTCPAddr("tcp", "127.0.0.1:2")
	require.NoError(t, err)
	pp := pulsenetwork.NewPulsePacketWithTrace(context.Background(), &pu, ph, th, 0)

	bs, err := packet.SerializePacket(pp)
	require.NoError(t, err)

	buf := bytes.NewBuffer(make([]byte, 0, packetMaxSize))

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
	phase1p.EncryptableBody = &GlobulaConsensusPacketBody{}
	phase1p.EncryptableBody.(*GlobulaConsensusPacketBody).PulsarPacket.setData(bs)

	p.Header.SetFlag(0)
	phase1p.Header.SetFlag(0)

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
			phase1p,
		},
		{
			"phase1",
			phases.PacketPhase1,
			91,
			phase1p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			90,
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
			p.Header.ClearFlag(0)

			s, err := p.SerializeTo(context.Background(), buf, digester, signer)
			require.NoError(t, err)
			require.EqualValues(t, test.size, s)
			require.EqualValues(t, s, p.Header.getPayloadLength())

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
			91,
			p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			239,
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
			require.EqualValues(t, s, p.Header.getPayloadLength())

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			p2.EncryptableBody.(*GlobulaConsensusPacketBody).BriefSelfIntro.JoinerData = nil
			p2.EncryptableBody.(*GlobulaConsensusPacketBody).FullSelfIntro.JoinerData = nil
			p2.EncryptableBody.(*GlobulaConsensusPacketBody).Announcement.Member.Joiner.JoinerData = nil
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
			454,
			p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			453,
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
			require.EqualValues(t, s, p.Header.getPayloadLength())

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			p2.EncryptableBody.(*GlobulaConsensusPacketBody).BriefSelfIntro.JoinerData = nil
			p2.EncryptableBody.(*GlobulaConsensusPacketBody).FullSelfIntro.JoinerData = nil
			p2.EncryptableBody.(*GlobulaConsensusPacketBody).Announcement.Member.Joiner.JoinerData = nil
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
			518,
			p,
		},
		{
			"phase2",
			phases.PacketPhase2,
			517,
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
			require.EqualValues(t, s, p.Header.getPayloadLength())

			require.NotEmpty(t, p.PacketSignature)

			p2 := Packet{}

			_, err = p2.DeserializeFrom(context.Background(), buf)
			require.NoError(t, err)

			p2.EncryptableBody.(*GlobulaConsensusPacketBody).BriefSelfIntro.JoinerData = nil
			p2.EncryptableBody.(*GlobulaConsensusPacketBody).FullSelfIntro.JoinerData = nil
			p2.EncryptableBody.(*GlobulaConsensusPacketBody).Announcement.Member.Joiner.JoinerData = nil
			require.Equal(t, p, p2)
		})
	}
}
