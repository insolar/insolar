package adapters

import (
	"context"

	transport2 "github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/transport"
)

type PacketSender struct {
	datagramTransport transport.DatagramTransport
}

func NewPacketSender(datagramTransport transport.DatagramTransport) *PacketSender {
	return &PacketSender{
		datagramTransport: datagramTransport,
	}
}

func (ps *PacketSender) SendPacketToTransport(ctx context.Context, to transport2.TargetProfile, sendOptions transport2.PacketSendOptions, payload interface{}) {
	addr := to.GetStatic().GetDefaultEndpoint().GetIPAddress().String()

	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"receiver_addr": addr,
	})

	err := ps.datagramTransport.SendDatagram(ctx, addr, payload.([]byte))
	if err != nil {
		logger.Error("Failed to send datagram: ", err)
	}
}
