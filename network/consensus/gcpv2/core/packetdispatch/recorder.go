// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package packetdispatch

import (
	"sync"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type PostponedPacketFunc func(packet transport.PacketParser, from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool

type PostponedPacket struct {
	Packet      transport.PacketParser
	From        endpoints.Inbound
	VerifyFlags coreapi.PacketVerifyFlags
}

func NewPacketRecorder(recordingSize int) PacketRecorder {
	return PacketRecorder{pr: UnsafePacketRecorder{recordingLimit: recordingSize}}
}

func NewUnsafePacketRecorder(recordingSize int) UnsafePacketRecorder {
	return UnsafePacketRecorder{recordingLimit: recordingSize}
}

type packetRecording struct {
	packets []PostponedPacket
}

func (p *packetRecording) Record(packet transport.PacketParser, from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) {
	p.packets = append(p.packets, PostponedPacket{packet, from, verifyFlags})
}

type PacketRecorder struct {
	sync sync.Mutex
	pr   UnsafePacketRecorder
}

func (p *PacketRecorder) Record(packet transport.PacketParser, from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) {
	p.sync.Lock()
	defer p.sync.Unlock()
	p.pr.Record(packet, from, verifyFlags)
}

func (p *PacketRecorder) Playback(to PostponedPacketFunc) {
	p.sync.Lock()
	defer p.sync.Unlock()
	p.pr.Playback(to)
}

type UnsafePacketRecorder struct {
	recordingLimit int
	playbackFn     PostponedPacketFunc
	recordings     []packetRecording
}

func (p *UnsafePacketRecorder) IsRecording() bool {
	return p.playbackFn == nil
}

func (p *UnsafePacketRecorder) Record(packet transport.PacketParser, from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) {
	if p.playbackFn != nil {
		go p.playbackFn(packet, from, verifyFlags)
		return
	}
	last := len(p.recordings) - 1
	if last < 0 || len(p.recordings[last].packets) >= p.recordingLimit {
		p.recordings = append(p.recordings, packetRecording{make([]PostponedPacket, 0, p.recordingLimit)})
		last++
	}
	p.recordings[last].Record(packet, from, verifyFlags)
}

func (p *UnsafePacketRecorder) Playback(to PostponedPacketFunc) {
	if p.playbackFn != nil {
		panic("illegal state")
	}
	if to == nil {
		panic("illegal value")
	}
	p.playbackFn = to

	recordings := p.recordings
	p.recordings = nil
	if len(recordings) > 0 {
		go playbackPackets(recordings, to)
	}
}

func playbackPackets(recordings []packetRecording, to PostponedPacketFunc) {
	for _, p := range recordings {
		for _, pp := range p.packets {
			if !to(pp.Packet, pp.From, pp.VerifyFlags) {
				return
			}
		}
	}
}
