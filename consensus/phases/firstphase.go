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
	"context"
	"math"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/merkle"
	"github.com/pkg/errors"
)

const ConsensusAtPercents = 0.66

func consensusReached(resultLen, participanstLen int) bool {
	minParticipants := int(math.Floor(ConsensusAtPercents*float64(participanstLen))) + 1

	return resultLen >= minParticipants
}

// FirstPhase is a first phase.
type FirstPhase struct {
	NodeNetwork  core.NodeNetwork         `inject:""`
	Calculator   merkle.Calculator        `inject:""`
	Communicator Communicator             `inject:""`
	Cryptography core.CryptographyService `inject:""`
	State        *FirstPhaseState
}

// Execute do first phase
func (fp *FirstPhase) Execute(ctx context.Context, pulse *core.Pulse) (*FirstPhaseState, error) {
	entry := &merkle.PulseEntry{Pulse: pulse}
	pulseHash, pulseProof, err := fp.Calculator.GetPulseProof(entry)

	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to calculate pulse proof.")
	}

	packet := packets.Phase1Packet{}
	err = packet.SetPulseProof(pulseProof.StateHash, pulseProof.Signature.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to set pulse proof in Phase1Packet.")
	}

	activeNodes := fp.NodeNetwork.GetActiveNodes()
	err = fp.signPhase1Packet(&packet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign a packet")
	}
	proofSet, err := fp.Communicator.ExchangePhase1(ctx, activeNodes, packet)
	if err != nil {
		return nil, errors.Wrap(err, "[ Execute ] Failed to exchange results.")
	}
	nodeProofs := make(map[core.Node]*merkle.PulseProof)
	for ref, packet := range proofSet {
		signIsCorrect, err := fp.isSignPhase1PacketRight(packet, ref)
		if err != nil {
			log.Warn("failed to check a sign: ", err.Error())
		} else if !signIsCorrect {
			log.Warn("recieved a bad sign packet: ", err.Error())
		}
		node := fp.NodeNetwork.GetActiveNode(ref)
		rawProof := packet.GetPulseProof()
		proof := &merkle.PulseProof{
			BaseProof: merkle.BaseProof{
				Signature: core.SignatureFromBytes(rawProof.Signature()),
			},
			StateHash: rawProof.StateHash(),
		}

		if !fp.Calculator.IsValid(proof, pulseHash, node.PublicKey()) {
			nodeProofs[node] = proof
		}
	}

	if !consensusReached(len(nodeProofs), len(activeNodes)) {
		return nil, errors.New("[ Execute ] Consensus not reached")
	}

	return &FirstPhaseState{
		PulseEntry:    entry,
		PulseHash:     pulseHash,
		PulseProof:    pulseProof,
		PulseProofSet: nodeProofs,
	}, nil
}

func (fp *FirstPhase) signPhase1Packet(packet *packets.Phase1Packet) error {
	data, err := packet.RawBytes()
	if err != nil {
		return errors.Wrap(err, "failed to get raw bytes")
	}
	sign, err := fp.Cryptography.Sign(data)
	if err != nil {
		return errors.Wrap(err, "failed to sign a phase 2 packet")
	}
	packet.Signature = sign.Bytes()
	return nil
}

func (fp *FirstPhase) isSignPhase1PacketRight(packet *packets.Phase1Packet, recordRef core.RecordRef) (bool, error) {
	key := fp.NodeNetwork.GetActiveNode(recordRef).PublicKey()
	raw, err := packet.RawBytes()

	if err != nil {
		return false, errors.Wrap(err, "failed to serialize packet")
	}
	return fp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}
