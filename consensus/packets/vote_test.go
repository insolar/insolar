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
)

func makeNodeListSupplementaryVote() *NodeListSupplementaryVote {
	nodeListVote := &NodeListSupplementaryVote{}
	nodeListVote.NodeListHash = randomArray32()
	nodeListVote.NodeListCount = uint16(77)

	return nodeListVote
}

func TestNodeListSupplementaryVote(t *testing.T) {
	checkSerializationDeserialization(t, makeNodeListSupplementaryVote())

	checkBadDataSerializationDeserialization(t, makeNodeListSupplementaryVote(), "unexpected EOF")
}

func TestNodeJoinSupplementaryVote(t *testing.T) {
	checkSerializationDeserialization(t, &NodeJoinSupplementaryVote{})
}

func makeStateFraudNodeSupplementaryVote() *StateFraudNodeSupplementaryVote {
	result := &StateFraudNodeSupplementaryVote{}
	result.Node1PulseProof = *makeNodePulseProof()
	result.Node2PulseProof = *makeNodePulseProof()
	result.PulseData = PulseData{PulseNumber: 1, Data: makeDefaultPulseDataExt()}

	return result
}

func TestStateFraudNodeSupplementaryVote(t *testing.T) {
	checkSerializationDeserialization(t, makeStateFraudNodeSupplementaryVote())

	checkBadDataSerializationDeserialization(t, makeStateFraudNodeSupplementaryVote(), "unexpected EOF")
}

func TestMissingNodeSupplementaryVote(t *testing.T) {
	checkSerializationDeserialization(t, &MissingNodeSupplementaryVote{*makeNodePulseProof()})

	checkBadDataSerializationDeserialization(t,
		&MissingNodeSupplementaryVote{*makeNodePulseProof()},
		"unexpected EOF")
}
