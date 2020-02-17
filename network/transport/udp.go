// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package transport

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
	"github.com/pkg/errors"
)

const (
	udpMaxPacketSize = 1450
)

type udpTransport struct {
	mutex              sync.RWMutex
	conn               net.PacketConn
	handler            DatagramHandler
	started            uint32
	fixedPublicAddress string
	cancel             context.CancelFunc
	address            string
}

func newUDPTransport(listenAddress, fixedPublicAddress string, handler DatagramHandler) *udpTransport {
	return &udpTransport{address: listenAddress, fixedPublicAddress: fixedPublicAddress, handler: handler}
}

// SendDatagram sends datagram to remote host
func (t *udpTransport) SendDatagram(ctx context.Context, address string, data []byte) error {
	if atomic.LoadUint32(&t.started) != 1 {
		return errors.New("failed to send datagram: transport is not started")
	}

	if len(data) > udpMaxPacketSize {
		return fmt.Errorf(
			"failed to send datagram: too big input data. Maximum: %d. Current: %d",
			udpMaxPacketSize,
			len(data),
		)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return errors.Wrap(err, "failed to resolve UDP address")
	}

	_, err = t.conn.WriteTo(data, udpAddr)
	if err != nil {
		// TODO: may be try to send second time if error
		return errors.Wrap(err, "failed to write data")
	}
	return nil
}

func (t *udpTransport) Address() string {
	return t.address
}

// Start starts networking.
func (t *udpTransport) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	if atomic.CompareAndSwapUint32(&t.started, 0, 1) {

		t.mutex.Lock()
		defer t.mutex.Unlock()

		var err error
		t.conn, err = net.ListenPacket("udp", t.address)
		if err != nil {
			return errors.Wrap(err, "failed to listen UDP")
		}

		t.address, err = resolver.Resolve(t.fixedPublicAddress, t.conn.LocalAddr().String())
		if err != nil {
			return errors.Wrap(err, "failed to resolve public address")
		}

		logger.Info("[ Start ] Start UDP transport")
		ctx, t.cancel = context.WithCancel(ctx)
		go t.loop(ctx)
	}

	return nil
}

func (t *udpTransport) loop(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		buf := make([]byte, udpMaxPacketSize)
		n, addr, err := t.conn.ReadFrom(buf)

		if err != nil {
			if network.IsConnectionClosed(err) {
				logger.Info("[ loop ] Connection closed, quiting ReadFrom loop")
				return
			}

			logger.Warn("[ loop ] failed to read UDP: ", err)
			continue
		}

		go t.handler.HandleDatagram(ctx, addr.String(), buf[:n])
	}
}

// Stop stops networking.
func (t *udpTransport) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	if atomic.CompareAndSwapUint32(&t.started, 1, 0) {
		logger.Info("[ Stop ] Stop UDP transport")

		t.cancel()
		err := t.conn.Close()
		if err != nil {
			if !network.IsConnectionClosed(err) {
				return err
			}
			logger.Warn("[ Stop ] Connection already closed")
		}
	}
	return nil
}
