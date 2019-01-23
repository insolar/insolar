/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
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
