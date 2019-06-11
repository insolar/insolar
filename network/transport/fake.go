//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
		go h.HandleDatagram(f.address, data)
	}

	return nil
}

func (f *fakeDatagramTransport) Address() string {
	return f.address
}

type fakeStreamTransport struct {
	address string
	handler StreamHandler
}

func (f *fakeStreamTransport) Start(ctx context.Context) error {
	tcpMutex.Lock()
	defer tcpMutex.Unlock()

	streamHandlers[f.address] = f.handler
	return nil
}

func (f *fakeStreamTransport) Stop(ctx context.Context) error {
	tcpMutex.Lock()
	defer tcpMutex.Unlock()

	streamHandlers[f.address] = nil
	return nil
}

func (f *fakeStreamTransport) Dial(ctx context.Context, address string) (io.ReadWriteCloser, error) {
	log.Debug("fakeStreamTransport Dial from %s to %s", f.address, address)

	tcpMutex.RLock()
	defer tcpMutex.RUnlock()

	h := streamHandlers[address]

	if h == nil {
		return nil, errors.New("dial failed")
	}

	conn1, conn2 := net.Pipe()
	go h.HandleStream(ctx, f.address, conn2)

	return conn1, nil
}

func (f *fakeStreamTransport) Address() string {
	return f.address
}
