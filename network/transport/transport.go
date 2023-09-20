package transport

import (
	"context"
	"io"

	component "github.com/insolar/component-manager"
)

// DatagramHandler interface provides callback method to process received datagrams
type DatagramHandler interface {
	HandleDatagram(ctx context.Context, address string, buf []byte)
}

// DatagramTransport interface provides methods to send and receive datagrams
type DatagramTransport interface {
	component.Starter
	component.Stopper

	SendDatagram(ctx context.Context, address string, data []byte) error
	Address() string
}

// StreamHandler interface provides callback method to process data stream
type StreamHandler interface {
	HandleStream(ctx context.Context, address string, stream io.ReadWriteCloser)
}

//go:generate minimock -i github.com/insolar/insolar/network/transport.StreamTransport -o ../../testutils/network -s _mock.go -g

// StreamTransport interface provides methods to send and receive data streams
type StreamTransport interface {
	component.Starter
	component.Stopper

	Dial(ctx context.Context, address string) (io.ReadWriteCloser, error)
	Address() string
}
