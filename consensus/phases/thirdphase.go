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

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/pkg/errors"
)

type ThirdPhase struct {
	Cryptography core.CryptographyService `inject:""`
	NodeNetwork  core.NodeNetwork         `inject:""`
	Communicator Communicator             `inject:""`
	NodeKeeper   network.NodeKeeper       `inject:""`

	newActiveNodeList []core.Node
	// TODO: insert it from somewhere
	mapper packets.BitSetMapper
}

func (tp *ThirdPhase) Execute(ctx context.Context, state *SecondPhaseState) error {
	var gSign [packets.SignatureLength]byte
	copy(gSign[:], state.GlobuleProof.Signature.Bytes()[:packets.SignatureLength])
	packet := packets.NewPhase3Packet(gSign, state.DBitSet)

	err := tp.signPhase3Packet(&packet)

	if err != nil {
		return errors.Wrap(err, "[ Execute ] failed to sign a phase 3 packet")
	}

	nodes := tp.NodeKeeper.GetActiveNodes()
	answers, err := tp.Communicator.ExchangePhase3(ctx, nodes, &packet)
	if err != nil {
		return errors.Wrap(err, "[ Execute ] failed to get answers on phase 3")
	}

	for ref, packet := range answers {
		signed, err := tp.isSignPhase3PacketRight(packet, ref)
		if err != nil {
			return errors.Wrap(err, "[ Execute ] failed to check a packet sign")
		} else if !signed {
			return errors.New("recv not signed packet")
		}
		cells, err := packet.GetBitset().GetCells(tp.mapper)
		if err != nil {
			return errors.Wrap(err, "[ Execute ] failed to get a cells")
		}
		for _, cell := range cells {
			if cell.State == packets.Legit {
				node, err := getNode(cell.NodeID, nodes)
				if err != nil {
					return errors.Wrap(err, "[ Execute ] failed to find a node on phase 3")
				}
				tp.newActiveNodeList = append(tp.newActiveNodeList, node)
			}
		}
	}

	return nil
}

func getNode(ref core.RecordRef, nodes []core.Node) (core.Node, error) {
	for _, node := range nodes {
		if ref == node.ID() {
			return node, nil
		}
	}
	return nil, errors.New("[ getNode] failed to find a node on phase 3")
}

func (tp *ThirdPhase) signPhase3Packet(p *packets.Phase3Packet) error {
	data, err := p.RawBytes()
	if err != nil {
		return errors.Wrap(err, "failed to get raw bytes")
	}
	sign, err := tp.Cryptography.Sign(data)
	if err != nil {
		return errors.Wrap(err, "failed to sign a phase 2 packet")
	}

	copy(p.SignatureHeaderSection1[:], sign.Bytes())
	// TODO: sign a second part after claim addition
	return nil
}

func (tp *ThirdPhase) isSignPhase3PacketRight(packet *packets.Phase3Packet, recordRef core.RecordRef) (bool, error) {
	key := tp.NodeNetwork.GetActiveNode(recordRef).PublicKey()

	raw, err := packet.RawBytes()
	if err != nil {
		return false, errors.Wrap(err, "failed to serialize")
	}

	return tp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}
