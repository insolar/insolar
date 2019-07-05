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
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

func (b *GlobulaConsensusPacketBody) DeserializeFrom(ctx DeserializeContext, data io.Reader) error {
	// todo: check packet type
	if ctx.HasFlag(0) {
		var pulsarData EmbeddedPulsarData
		err := pulsarData.DeserializeFrom(ctx, data)
		if err != nil {
			return errors.Wrap(err, "[ GlobulaConsensusPacketBody.Deserialize ] Can't deserialize EmbeddedPulsarData")
		}
	}

	// todo: check packet type
	var announcement MembershipAnnouncement
	err := announcement.DeserializeFrom(ctx, data)
	if err != nil {
		return errors.Wrap(err, "[ GlobulaConsensusPacketBody.Deserialize ] Can't read Announcement")
	}
	b.Announcement = &announcement

	panic("implement me")
}

func (p *EmbeddedPulsarData) DeserializeFrom(ctx DeserializeContext, data io.Reader) error {
	err := p.Header.DeserializeFrom(nil, data)
	if err != nil {
		return errors.Wrap(err, "[ EmbeddedPulsarData.Deserialize ] Can't deserialize Header")
	}

	// todo: PulsarPacketBody

	err = binary.Read(data, defaultByteOrder, &p.PulsarSignature)
	if err != nil {
		return errors.Wrap(err, "[ Header.Deserialize ] Can't read ReceiverID")
	}

	panic("implement me")
}

func (m *MembershipAnnouncement) DeserializeFrom(ctx DeserializeContext, data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &m.CurrentRank)
	if err != nil {
		return errors.Wrap(err, "[ MembershipAnnouncement.Deserialize ] Can't read CurrentRank")
	}

	if m.CurrentRank != 0 {
		err := m.NonJoinerMembershipAnnouncement.DeserializeFrom(ctx, data)
		if err != nil {
			return err
		}
	}

	panic("implement me")
}

func (m *NonJoinerMembershipAnnouncement) DeserializeFrom(ctx DeserializeContext, data io.Reader) error {
	err := binary.Read(data, defaultByteOrder, &m.RequestedPower)
	if err != nil {
		return errors.Wrap(err, "[ MembershipAnnouncement.Deserialize ] Can't read RequestedPower")
	}
	// todo: Member
	// todo: AnnounceSignature

	panic("implement me")
}
