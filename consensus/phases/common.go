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

package phases

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
)

func validateProofs(
	calculator merkle.Calculator,
	unsyncList network.UnsyncList,
	pulseHash merkle.OriginHash,
	proofs map[core.RecordRef]*merkle.PulseProof,
) (valid map[core.Node]*merkle.PulseProof, fault map[core.RecordRef]*merkle.PulseProof) {

	validProofs := make(map[core.Node]*merkle.PulseProof)
	faultProofs := make(map[core.RecordRef]*merkle.PulseProof)
	for nodeID, proof := range proofs {
		valid := validateProof(calculator, unsyncList, pulseHash, nodeID, proof)
		if valid {
			validProofs[unsyncList.GetActiveNode(nodeID)] = proof
		} else {
			faultProofs[nodeID] = proof
		}
	}
	return validProofs, faultProofs
}
func validateProof(
	calculator merkle.Calculator,
	unsyncList network.UnsyncList,
	pulseHash merkle.OriginHash,
	nodeID core.RecordRef,
	proof *merkle.PulseProof) bool {

	node := unsyncList.GetActiveNode(nodeID)
	if node == nil {
		return false
	}
	return calculator.IsValid(proof, pulseHash, node.PublicKey())
}
