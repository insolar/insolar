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
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/merkle"
)

type FirstPhaseState struct {
	PulseEntry *merkle.PulseEntry

	PulseHash  merkle.OriginHash
	PulseProof *merkle.PulseProof

	ValidProofs map[core.Node]*merkle.PulseProof
	FaultProofs map[core.RecordRef]*merkle.PulseProof

	UnsyncList network.UnsyncList
}

type SecondPhaseState struct {
	*FirstPhaseState

	GlobuleEntry *merkle.GlobuleEntry

	GlobuleHash  merkle.OriginHash
	GlobuleProof *merkle.GlobuleProof

	GlobuleProofSet map[core.Node]*merkle.GlobuleProof

	NodeListCount uint16
	NodeListHash  []byte

	DBitSet packets.BitSet
}

type ThirdPhasePulseState struct {
}

type ThirdPhaseReferendumState struct {
}
