package transport

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
)

type fakeNode struct {
	component.Starter
	component.Stopper

	tcp    StreamTransport
	udp    DatagramTransport
	udpBuf chan []byte
	tcpBuf chan []byte
}

func (f *fakeNode) HandleStream(address string, stream io.ReadWriteCloser) {
	log.Printf("HandleStream from %s", address)

	b := make([]byte, 3)
	_, err := stream.Read(b)
	if err != nil {
		log.Printf("Failed to read from connection")
	}

	f.tcpBuf <- b
}

func (f *fakeNode) HandleDatagram(address string, buf []byte) {
	log.Printf("HandleDatagram from %s: %v", address, buf)
	f.udpBuf <- buf
}

func (f *fakeNode) Start(ctx context.Context) error {
	err1 := f.udp.Start(ctx)
	err2 := f.tcp.Start(ctx)
	if err1 != nil || err2 != nil {
		return err1
	} else {
		return nil
	}
}

func (f *fakeNode) Stop(ctx context.Context) error {
	err1 := f.udp.Stop(ctx)
	err2 := f.tcp.Stop(ctx)
	if err1 != nil || err2 != nil {
		return err1
	} else {
		return nil
	}
}

func newFakeNode(address string) *fakeNode {
	cfg := configuration.NewHostNetwork().Transport
	cfg.Address = address

	factory := NewFakeFactory(cfg)
	n := &fakeNode{}
	n.udp, _ = factory.CreateDatagramTransport(n)
	n.tcp, _ = factory.CreateStreamTransport(n)

	n.udpBuf = make(chan []byte, 1)
	n.tcpBuf = make(chan []byte, 1)
	return n
}

func TestFakeNetwork_TCP(t *testing.T) {
	ctx := context.Background()
	n1 := newFakeNode("127.0.0.1:8080")
	n2 := newFakeNode("127.0.0.1:4200")
	assert.NotNil(t, n2)

	assert.NoError(t, n1.Start(ctx))
	assert.NoError(t, n2.Start(ctx))

	_, err := n2.tcp.Dial(ctx, "127.0.0.1:5555")
	assert.Error(t, err)

	conn, err := n1.tcp.Dial(ctx, n2.tcp.Address())
	assert.NoError(t, err)

	count, err := conn.Write([]byte{1, 2, 3})
	assert.Equal(t, 3, count)
	assert.NoError(t, err)
	assert.NoError(t, conn.Close())

	assert.Equal(t, []byte{1, 2, 3}, <-n2.tcpBuf)

	assert.NoError(t, n1.Stop(ctx))
	assert.NoError(t, n2.Stop(ctx))
}

func TestFakeNetwork_UDP(t *testing.T) {
	ctx := context.Background()
	n1 := newFakeNode("127.0.0.1:8080")
	n2 := newFakeNode("127.0.0.1:4200")
	assert.NotNil(t, n2)

	assert.NoError(t, n1.Start(ctx))
	assert.NoError(t, n2.Start(ctx))

	err := n1.udp.SendDatagram(ctx, n2.udp.Address(), []byte{1, 2, 3})
	assert.NoError(t, err)

	err = n2.udp.SendDatagram(ctx, n1.udp.Address(), []byte{5, 4, 3})
	assert.NoError(t, err)

	assert.Equal(t, []byte{1, 2, 3}, <-n2.udpBuf)
	assert.Equal(t, []byte{5, 4, 3}, <-n1.udpBuf)

	assert.NoError(t, n1.Stop(ctx))
	assert.NoError(t, n2.Stop(ctx))
}
