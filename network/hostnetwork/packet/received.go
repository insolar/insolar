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
