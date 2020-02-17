// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
