// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package packet

type ReceivedPacket struct {
	*Packet
	data []byte
}

func NewReceivedPacket(p *Packet, data []byte) *ReceivedPacket {
	return &ReceivedPacket{
		Packet: p,
		data:   data,
	}
}

func (p *ReceivedPacket) Bytes() []byte {
	return p.data
}
