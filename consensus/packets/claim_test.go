/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
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

func makeNodeJoinClaim() *NodeJoinClaim {
	nodeJoinClaim := &NodeJoinClaim{}
	nodeJoinClaim.ShortNodeID = core.ShortNodeID(77)
	nodeJoinClaim.RelayNodeID = core.ShortNodeID(26)
	nodeJoinClaim.ProtocolVersionAndFlags = uint32(99)
	nodeJoinClaim.JoinsAfter = uint32(67)
	nodeJoinClaim.NodeRoleRecID = 32
	nodeJoinClaim.NodeRef = testutils.RandomRef()
	nodeJoinClaim.NodePK = randomArray64()
	nodeJoinClaim.Signature = randomArray71()

	return nodeJoinClaim
}

func TestNodeJoinClaim(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeJoinClaim())
}

func TestNodeJoinClaim_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodeJoinClaim(), "unexpected EOF")
}

func TestNodeLeaveClaim(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	checkSerializationDeserialization(t, nodeLeaveClaim)
}

func TestMakeClaimHeader(t *testing.T) {

}

func makeNodeAnnounceClaim() *NodeAnnounceClaim {
	nodeAnnounceClaim := &NodeAnnounceClaim{}
	nodeAnnounceClaim.NodeJoinClaim = *makeNodeJoinClaim()
	nodeAnnounceClaim.NodeCount = 266
	nodeAnnounceClaim.NodeIndex = 37
	return nodeAnnounceClaim
}

func TestNodeAnnounceClaim(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeAnnounceClaim())
}
