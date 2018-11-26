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
	"github.com/pkg/errors"
)

type ThirdPhase struct {
	Cryptography core.CryptographyService `inject:""`
	NodeNetwork  core.NodeNetwork         `inject:""`
}

func (tpr *ThirdPhase) Execute(ctx context.Context, state *SecondPhaseState) error {
	// TODO: do something here
	return nil
}

func (tp *ThirdPhase) signPhase2Packet(p *packets.Phase3Packet) error {
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

func (tp *ThirdPhase) isSignPhase2PacketRight(packet *packets.Phase3Packet, recordRef core.RecordRef) (bool, error) {
	key := tp.NodeNetwork.GetActiveNode(recordRef).PublicKey()

	raw, err := packet.RawBytes()
	if err != nil {
		return false, errors.Wrap(err, "failed to serialize")
	}

	return tp.Cryptography.Verify(key, core.SignatureFromBytes(raw), raw), nil
}
