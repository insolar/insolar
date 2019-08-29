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
	"io"
	"log"
	"net"
	"strconv"
	"testing"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
)

type suiteTest struct {
	suite.Suite

	factory1 Factory
	factory2 Factory
}

type fakeNode struct {
	component.Starter
	component.Stopper

	tcp    StreamTransport
	udp    DatagramTransport
	udpBuf chan []byte
	tcpBuf chan []byte
}

func (f *fakeNode) HandleStream(ctx context.Context, address string, stream io.ReadWriteCloser) {
	inslogger.FromContext(ctx).Infof("HandleStream from %s", address)

	b := make([]byte, 3)
	_, err := stream.Read(b)
	if err != nil {
		log.Printf("Failed to read from connection")
	}

	f.tcpBuf <- b
}

func (f *fakeNode) HandleDatagram(ctx context.Context, address string, buf []byte) {
	inslogger.FromContext(ctx).Info("HandleDatagram from %s: %v", address, buf)
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

func newFakeNode(f Factory) *fakeNode {
	n := &fakeNode{}
	n.udp, _ = f.CreateDatagramTransport(n)
	n.tcp, _ = f.CreateStreamTransport(n)

	n.udpBuf = make(chan []byte, 1)
	n.tcpBuf = make(chan []byte, 1)
	return n
}

func (s *suiteTest) TestStreamTransport() {
	ctx := context.Background()
	n1 := newFakeNode(s.factory1)
	n2 := newFakeNode(s.factory2)
	s.NotNil(n2)

	s.NoError(n1.Start(ctx))
	s.NoError(n2.Start(ctx))

	_, port1, _ := net.SplitHostPort(n1.tcp.Address())
	_, port2, _ := net.SplitHostPort(n2.tcp.Address())
	n1Port, _ := strconv.Atoi(port1)
	n2Port, _ := strconv.Atoi(port2)
	port := 5555
	for port == n1Port || port == n2Port {
		port++
	}

	_, err := n2.tcp.Dial(ctx, "127.0.0.1:"+strconv.Itoa(port))
	s.Error(err)

	_, err = n2.tcp.Dial(ctx, "127.0.0.1:555555")
	s.Error(err)

	_, err = n2.tcp.Dial(ctx, "invalid address")
	s.Error(err)

	conn, err := n1.tcp.Dial(ctx, n2.tcp.Address())
	s.Require().NoError(err)

	count, err := conn.Write([]byte{1, 2, 3})
	s.Equal(3, count)
	s.NoError(err)
	s.NoError(conn.Close())

	s.Equal([]byte{1, 2, 3}, <-n2.tcpBuf)

	s.NoError(n1.Stop(ctx))
	s.NoError(n2.Stop(ctx))
}

func (s *suiteTest) TestDatagramTransport() {
	ctx := context.Background()
	n1 := newFakeNode(s.factory1)
	n2 := newFakeNode(s.factory2)
	s.NotNil(n2)

	s.NoError(n1.Start(ctx))
	s.NoError(n2.Start(ctx))

	err := n1.udp.SendDatagram(ctx, n2.udp.Address(), []byte{1, 2, 3})
	s.NoError(err)

	err = n2.udp.SendDatagram(ctx, n1.udp.Address(), []byte{5, 4, 3})
	s.NoError(err)

	err = n2.udp.SendDatagram(ctx, "invalid address", []byte{9, 9, 9})
	s.Error(err)

	bigBuff := make([]byte, udpMaxPacketSize+1)
	err = n2.udp.SendDatagram(ctx, n1.udp.Address(), bigBuff)
	s.Error(err)

	s.Equal([]byte{1, 2, 3}, <-n2.udpBuf)
	s.Equal([]byte{5, 4, 3}, <-n1.udpBuf)

	s.NoError(n1.Stop(ctx))
	s.NoError(n2.Stop(ctx))
}

func TestFakeTransport(t *testing.T) {

	cfg1 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:8080"}
	cfg2 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:4200"}

	f1 := NewFakeFactory(cfg1)
	f2 := NewFakeFactory(cfg2)

	suite.Run(t, &suiteTest{factory1: f1, factory2: f2})
}

func TestTransport(t *testing.T) {
	cfg1 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:0"}
	cfg2 := configuration.Transport{Protocol: "TCP", Address: "127.0.0.1:0"}

	f1 := NewFactory(cfg1)
	f2 := NewFactory(cfg2)
	suite.Run(t, &suiteTest{factory1: f1, factory2: f2})
}
