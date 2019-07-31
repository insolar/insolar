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

package coreapi

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

func VerifyPacketAuthenticityBy(packetSignature cryptkit.SignedDigest, nr profiles.Host, sf cryptkit.SignatureVerifier,
	from endpoints.Inbound, strictFrom bool) error {

	if strictFrom && !nr.IsAcceptableHost(from) {
		return fmt.Errorf("host is not allowed by node registration: node=%v, host=%v", nr, from)
	}
	if !packetSignature.IsVerifiableBy(sf) {
		return fmt.Errorf("unable to verify packet signature from sender: %v", from)
	}
	if !packetSignature.VerifyWith(sf) {
		return fmt.Errorf("packet signature doesn't match for sender: %v", from)
	}

	return nil
}

func FindHostProfile(memberID insolar.ShortNodeID, from endpoints.Inbound, initialCensus census.Operational) profiles.Host {

	if np := initialCensus.GetOnlinePopulation().FindProfile(memberID); np != nil {
		return np.GetStatic()
	}
	if nr := initialCensus.GetOfflinePopulation().FindRegisteredProfile(from); nr != nil {
		return nr
	}
	if nr := initialCensus.GetMandateRegistry().FindRegisteredProfile(from); nr != nil {
		return nr
	}
	return nil
}

func VerifyPacketRoute(ctx context.Context, packet transport.PacketParser, selfID insolar.ShortNodeID, from endpoints.Inbound) (bool, error) {

	sid := packet.GetSourceID()
	if sid.IsAbsent() {
		return false, fmt.Errorf("invalid sourceID(0): from=%v", from)
	}
	if sid == selfID {
		return false, fmt.Errorf("loopback, SourceID(%v) == thisNodeID(%v): from=%v", sid, selfID, from)
	}

	rid := packet.GetReceiverID()
	if rid.IsAbsent() {
		return false, fmt.Errorf("invalid receiverID(0): from=%v", from)
	}
	if rid != selfID {
		return false, fmt.Errorf("receiverID(%v) != thisNodeID(%v): from=%v", rid, selfID, from)
	}

	tid := packet.GetTargetID()
	if rid.IsAbsent() {
		return false, fmt.Errorf("invalid targetID(0): from=%v", from)
	}

	if tid != selfID {
		// Relaying
		if packet.IsRelayForbidden() {
			return false, fmt.Errorf("sender doesn't allow relaying for targetID(%v)", tid)
		}

		// TODO relay support
		err := fmt.Errorf("unsupported: relay is required for targetID(%v)", tid)
		inslogger.FromContext(ctx).Errorf(err.Error())
		// allow sender to be different from source
		return false, err
	}

	// sender must be source
	return packet.IsRelayForbidden(), nil
}

func LazyPacketParse(packet transport.PacketParser) (transport.PacketParser, error) {

	// this enables lazy parsing - packet is fully parsed AFTER validation, hence makes it less prone to exploits for non-members
	newPacket, err := packet.ParsePacketBody()
	if err != nil {
		return packet, err
	}
	if newPacket == nil {
		return packet, nil
	}
	return newPacket, nil
}
