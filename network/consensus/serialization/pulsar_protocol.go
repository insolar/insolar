// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package serialization

import (
	"io"

	"github.com/insolar/insolar/pulse"
)

type PulsarPacketBody struct {
	// ByteSize>=108
	PulseNumber           pulse.Number  `insolar-transport:"ignore=send"`
	PulseDataExt          pulse.DataExt // ByteSize=44
	PulsarConsensusProofs []byte        // variable lengths >=0
}

func (b *PulsarPacketBody) String(ctx PacketContext) string {
	return "pulsar packet body"
}

func (b *PulsarPacketBody) SerializeTo(_ SerializeContext, writer io.Writer) error {
	// TODO: proofs
	return write(writer, b.PulseDataExt)
}

func (b *PulsarPacketBody) DeserializeFrom(_ DeserializeContext, reader io.Reader) error {
	// TODO: proofs
	return read(reader, &b.PulseDataExt)
}

func (b *PulsarPacketBody) getPulseData() pulse.Data {
	return pulse.Data{
		PulseNumber: b.PulseNumber,
		DataExt:     b.PulseDataExt,
	}
}
