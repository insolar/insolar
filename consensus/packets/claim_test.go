/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package packets

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
)

func makeNodeBroadCast() *NodeBroadcast {
	nodeBroadcast := &NodeBroadcast{}
	nodeBroadcast.EmergencyLevel = uint8(4)

	return nodeBroadcast
}

func TestNodeBroadcast(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeBroadCast())
}

func makeCapabilityPoolingAndActivation() *CapabilityPoolingAndActivation {
	capabilityPoolingAndActivation := &CapabilityPoolingAndActivation{}
	capabilityPoolingAndActivation.PollingFlags = uint16(10)
	capabilityPoolingAndActivation.CapabilityType = uint16(7)
	capabilityPoolingAndActivation.CapabilityRef = randomArray64()

	return capabilityPoolingAndActivation
}

func TestCapabilityPoolingAndActivation(t *testing.T) {
	checkSerializationDeserialization(t, makeCapabilityPoolingAndActivation())
}

func makeNodeViolationBlame() *NodeViolationBlame {
	nodeViolationBlame := &NodeViolationBlame{}
	nodeViolationBlame.TypeViolation = uint8(4)

	return nodeViolationBlame
}

func TestNodeViolationBlame(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeViolationBlame())
}

func makeNodeJoinClaim(withSignature bool) *NodeJoinClaim {
	nodeJoinClaim := &NodeJoinClaim{}
	nodeJoinClaim.ShortNodeID = core.ShortNodeID(77)
	nodeJoinClaim.RelayNodeID = core.ShortNodeID(26)
	nodeJoinClaim.ProtocolVersionAndFlags = uint32(99)
	nodeJoinClaim.JoinsAfter = uint32(67)
	nodeJoinClaim.NodeRoleRecID = 32
	nodeJoinClaim.NodeRef = testutils.RandomRef()
	nodeJoinClaim.NodePK = randomArray66()
	if withSignature {
		nodeJoinClaim.Signature = randomArray66()
	}
	nodeJoinClaim.NodeAddress.Set("127.0.0.1:5566")

	return nodeJoinClaim
}

func TestNodeJoinClaim(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeJoinClaim(true))
}

func TestNodeJoinClaim_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodeJoinClaim(true), "unexpected EOF")
}

func TestNodeLeaveClaim(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	checkSerializationDeserialization(t, nodeLeaveClaim)
}

func TestMakeClaimHeader(t *testing.T) {

}

func makeNodeAnnounceClaim() *NodeAnnounceClaim {
	nodeAnnounceClaim := &NodeAnnounceClaim{}
	nodeAnnounceClaim.NodeJoinClaim = *makeNodeJoinClaim(true)
	nodeAnnounceClaim.NodeCount = 266
	nodeAnnounceClaim.NodeAnnouncerIndex = 37
	nodeAnnounceClaim.NodeJoinerIndex = 38
	return nodeAnnounceClaim
}

func TestNodeAnnounceClaim(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeAnnounceClaim())
}
