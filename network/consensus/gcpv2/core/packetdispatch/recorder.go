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

package packetdispatch

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"sync"

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
