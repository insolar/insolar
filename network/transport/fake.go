// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package transport

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
)

var (
	udpMutex         sync.RWMutex
	datagramHandlers = make(map[string]DatagramHandler)
	tcpMutex         sync.RWMutex
	streamHandlers   = make(map[string]StreamHandler)
)

// NewFakeFactory constructor creates new fake transport factory
func NewFakeFactory(cfg configuration.Transport) Factory {
	return &fakeFactory{cfg: cfg}
}

type fakeFactory struct {
	cfg configuration.Transport
}

// CreateStreamTransport creates fake StreamTransport for tests
func (f *fakeFactory) CreateStreamTransport(handler StreamHandler) (StreamTransport, error) {
	return &fakeStreamTransport{address: f.cfg.Address, handler: handler}, nil
}

// CreateDatagramTransport creates fake DatagramTransport for tests
func (f *fakeFactory) CreateDatagramTransport(handler DatagramHandler) (DatagramTransport, error) {
	return &fakeDatagramTransport{address: f.cfg.Address, handler: handler}, nil
}

type fakeDatagramTransport struct {
	address string
	handler DatagramHandler
}

func (f *fakeDatagramTransport) Start(ctx context.Context) error {
	udpMutex.Lock()
	defer udpMutex.Unlock()

	datagramHandlers[f.address] = f.handler
	return nil
}

func (f *fakeDatagramTransport) Stop(ctx context.Context) error {
	udpMutex.Lock()
	defer udpMutex.Unlock()

	datagramHandlers[f.address] = nil
	return nil
}

func (f *fakeDatagramTransport) SendDatagram(ctx context.Context, address string, data []byte) error {
	log.Debugf("fakeDatagramTransport SendDatagram to %s : %v", address, data)

	if len(data) > udpMaxPacketSize {
		return errors.New(fmt.Sprintf("udpTransport.send: too big input data. Maximum: %d. Current: %d",
			udpMaxPacketSize, len(data)))
	}
	_, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return errors.Wrap(err, "Failed to resolve UDP address")
	}

	udpMutex.RLock()
	defer udpMutex.RUnlock()

	h := datagramHandlers[address]
	if h != nil {
		cpData := make([]byte, len(data))
		copy(cpData, data)
		go h.HandleDatagram(ctx, f.address, cpData)
	}

	return nil
}

func (f *fakeDatagramTransport) Address() string {
	return f.address
}

type fakeStreamTransport struct {
	address string
	handler StreamHandler
	cancel  context.CancelFunc
	ctx     context.Context
}

func (f *fakeStreamTransport) Start(ctx context.Context) error {
	tcpMutex.Lock()
	defer tcpMutex.Unlock()

	f.ctx, f.cancel = context.WithCancel(ctx)
	streamHandlers[f.address] = f.handler
	return nil
}

func (f *fakeStreamTransport) Stop(ctx context.Context) error {
	tcpMutex.Lock()
	defer tcpMutex.Unlock()

	f.cancel()
	streamHandlers[f.address] = nil
	return nil
}

func (f *fakeStreamTransport) Dial(ctx context.Context, address string) (io.ReadWriteCloser, error) {
	log.Debugf("fakeStreamTransport Dial from %s to %s", f.address, address)

	tcpMutex.RLock()
	defer tcpMutex.RUnlock()

	h := streamHandlers[address]

	if h == nil {
		return nil, errors.New("fakeStreamTransport: dial failed")
	}

	conn1, conn2 := net.Pipe()
	go h.HandleStream(f.ctx, f.address, conn2)

	return conn1, nil
}

func (f *fakeStreamTransport) Address() string {
	return f.address
}
