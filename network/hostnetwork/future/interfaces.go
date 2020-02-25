// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package future

import (
	"context"
	"errors"
	"time"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

var (
	// ErrTimeout is returned when the operation timeout is exceeded.
	ErrTimeout = errors.New("timeout")
	// ErrChannelClosed is returned when the input channel is closed.
	ErrChannelClosed = errors.New("channel closed")
)

// Future is network response future.
type Future interface {

	// ID returns packet sequence number.
	ID() types.RequestID

	// Receiver returns the initiator of the packet.
	Receiver() *host.Host

	// Request returns origin request.
	Request() network.Packet

	// Response is a channel to listen for future response.
	Response() <-chan network.ReceivedPacket

	// SetResponse makes packet to appear in response channel.
	SetResponse(network.ReceivedPacket)

	// WaitResponse gets the future response from Response() channel with a timeout set to `duration`.
	WaitResponse(duration time.Duration) (network.ReceivedPacket, error)

	// Cancel closes all channels and cleans up underlying structures.
	Cancel()
}

// CancelCallback is a callback function executed when cancelling Future.
type CancelCallback func(Future)

type Manager interface {
	Get(packet *packet.Packet) Future
	Create(packet *packet.Packet) Future
}

type PacketHandler interface {
	Handle(ctx context.Context, msg *packet.ReceivedPacket)
}
