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

	"github.com/insolar/insolar/testutils"
)

func makeNodeBroadCast() *NodeBroadcast {
	nodeBroadcast := &NodeBroadcast{}
	nodeBroadcast.length = uint16(3)
	nodeBroadcast.EmergencyLevel = uint8(4)

	return nodeBroadcast
}

func TestNodeBroadcast(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeBroadCast())
}

func TestNodeBroadcast_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodeBroadCast(),
		"[ NodeBroadcast.Deserialize ] Can't read length: unexpected EOF")
}

func makeCapabilityPoolingAndActivation() *CapabilityPoolingAndActivation {
	capabilityPoolingAndActivation := &CapabilityPoolingAndActivation{}
	capabilityPoolingAndActivation.PollingFlags = uint16(10)
	capabilityPoolingAndActivation.length = uint16(7)
	capabilityPoolingAndActivation.CapabilityType = uint16(7)
	capabilityPoolingAndActivation.CapabilityRef = randomArray64()

	return capabilityPoolingAndActivation
}

func TestCapabilityPoolingAndActivation(t *testing.T) {
	checkSerializationDeserialization(t, makeCapabilityPoolingAndActivation())
}

func TestCapabilityPoolingAndActivation_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeCapabilityPoolingAndActivation(),
		"[ CapabilityPoolingAndActivation.Deserialize ] Can't read length: unexpected EOF")
}

func makeNodeViolationBlame() *NodeViolationBlame {
	nodeViolationBlame := &NodeViolationBlame{}
	nodeViolationBlame.length = uint16(2)
	nodeViolationBlame.claimType = TypeNodeViolationBlame
	nodeViolationBlame.TypeViolation = uint8(4)

	return nodeViolationBlame
}

func TestNodeViolationBlame(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeViolationBlame())
}

func TestNodeViolationBlame_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodeViolationBlame(),
		"[ NodeViolationBlame.Deserialize ] Can't read length: unexpected EOF")
}

func makeNodeJoinClaim() *NodeJoinClaim {
	nodeJoinClaim := &NodeJoinClaim{}
	nodeJoinClaim.NodeID = uint32(77)
	nodeJoinClaim.RelayNodeID = uint32(26)
	nodeJoinClaim.ProtocolVersionAndFlags = uint32(99)
	nodeJoinClaim.JoinsAfter = uint32(67)
	nodeJoinClaim.NodeRoleRecID = uint32(32)
	nodeJoinClaim.NodeRef = testutils.RandomRef()
	nodeJoinClaim.NodePK = randomArray64()

	return nodeJoinClaim
}

func TestNodeJoinClaim(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeJoinClaim())
}

func TestNodeJoinClaim_BadData(t *testing.T) {
	checkBadDataSerializationDeserialization(t, makeNodeJoinClaim(),
		"[ NodeJoinClaim.Deserialize ] Can't read NodePK: unexpected EOF")
}

func TestNodeLeaveClaim(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	nodeLeaveClaim.length = uint16(333)
	checkSerializationDeserialization(t, nodeLeaveClaim)
}

func TestNodeLeaveClaim_BadData(t *testing.T) {
	nodeLeaveClaim := &NodeLeaveClaim{}
	nodeLeaveClaim.length = uint16(333)
	checkBadDataSerializationDeserialization(t, nodeLeaveClaim,
		"[ NodeLeaveClaim.Deserialize ] Can't read length: unexpected EOF")
}
