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
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/utils"
)

const udpMaxPacketSize = 1400

type udpTransport struct {
	baseTransport
	conn    net.PacketConn
	address string
}

type udpSerializer struct{}

func (b *udpSerializer) SerializePacket(q *packet.Packet) ([]byte, error) {
	data, ok := q.Data.(packets.ConsensusPacket)
	if !ok {
		return nil, errors.New("could not convert packet to ConsensusPacket type")
	}
	return data.Serialize()
}

func (b *udpSerializer) DeserializePacket(conn io.Reader) (*packet.Packet, error) {
	data, err := packets.ExtractPacket(conn)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert network datagram to ConsensusPacket")
	}
	p := &packet.Packet{}
	p.Data = data
	return p, nil
}

func newUDPTransport(listenAddress, fixedPublicAddress string) (*udpTransport, string, error) {
	conn, err := net.ListenPacket("udp", listenAddress)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to listen UDP")
	}
	publicAddress, err := Resolve(fixedPublicAddress, conn.LocalAddr().String())
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to resolve public address")
	}

	transport := &udpTransport{baseTransport: newBaseTransport(publicAddress), conn: conn}
	transport.sendFunc = transport.send
	transport.serializer = &udpSerializer{}

	return transport, publicAddress, nil
}

func (t *udpTransport) send(ctx context.Context, recvAddress string, data []byte) error {
	logger := inslogger.FromContext(ctx)

	logger.Debug("Sending PURE_UDP request")
	if len(data) > udpMaxPacketSize {
		return errors.New(fmt.Sprintf("udpTransport.send: too big input data. Maximum: %d. Current: %d",
			udpMaxPacketSize, len(data)))
	}

	// TODO: may be try to send second time if error
	// TODO: skip resolving every time by caching result
	udpAddr, err := net.ResolveUDPAddr("udp", recvAddress)
	if err != nil {
		return errors.Wrap(err, "udpTransport.send")
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return errors.Wrap(err, "udpTransport.send")
	}
	defer utils.CloseVerbose(udpConn)

	logger.Debug("udpTransport.send: len = ", len(data))
	n, err := udpConn.Write(data)
	stats.Record(context.Background(), consensus.SentSize.M(int64(n)))
	return errors.Wrap(err, "Failed to write data")
}

func (t *udpTransport) prepareListen() (net.PacketConn, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.disconnectStarted = make(chan bool, 1)
	t.disconnectFinished = make(chan bool, 1)

	if t.conn != nil {
		t.address = t.conn.LocalAddr().String()
	} else {
		var err error
		t.conn, err = net.ListenPacket("udp", t.address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to listen UDP")
		}
	}

	return t.conn, nil
}

// Start starts networking.
func (t *udpTransport) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Start ] Start UDP transport")

	conn, err := t.prepareListen()
	if err != nil {
		logger.Infof("[ Start ] Failed to prepare UDP transport: " + err.Error())
		return err
	}

	go func() {
		for {
			buf := make([]byte, udpMaxPacketSize)
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				<-t.disconnectFinished
				logger.Error("failed to read UDP: ", err.Error())
				return // TODO: we probably shouldn't return here
			}

			stats.Record(ctx, consensus.RecvSize.M(int64(n)))
			go t.handleAcceptedConnection(buf[:n], addr)
		}
	}()

	return nil
}

// Stop stops networking.
func (t *udpTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Info("Stop UDP transport")
	t.prepareDisconnect()

	if t.conn != nil {
		utils.CloseVerbose(t.conn)
		t.conn = nil
	}
}

func (t *udpTransport) handleAcceptedConnection(data []byte, addr net.Addr) {
	r := bytes.NewReader(data)
	msg, err := t.serializer.DeserializePacket(r)
	if err != nil {
		log.Error("[ handleAcceptedConnection ] ", err)
		return
	}
	log.Debug("[ handleAcceptedConnection ] Packet processed. size: ", len(data), ". Address: ", addr)

	go t.packetHandler.Handle(context.TODO(), msg)
}
