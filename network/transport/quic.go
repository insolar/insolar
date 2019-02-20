/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package transport

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/network/utils"
	quic "github.com/lucas-clemente/quic-go"
	"github.com/pkg/errors"
)

type quicConnection struct {
	session quic.Session
	stream  quic.Stream
}

type quicTransport struct {
	baseTransport
	l           quic.Listener
	conn        net.PacketConn
	connections map[string]quicConnection
}

func newQuicTransport(conn net.PacketConn, proxy relay.Proxy, publicAddress string) (*quicTransport, error) {
	listener, err := quic.Listen(conn, generateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}

	transport := &quicTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		l:             listener,
		conn:          conn,
		connections:   make(map[string]quicConnection),
	}

	transport.sendFunc = transport.send
	return transport, nil
}

func (t *quicTransport) send(recvAddress string, data []byte) error {
	conn, ok := t.connections[recvAddress]
	var stream quic.Stream
	var err error
	if !ok {
		var session quic.Session
		session, stream, err = createConnection(recvAddress)
		if err != nil {
			return errors.Wrap(err, "[ send ] failed to create a connection")
		}
		t.connections[recvAddress] = quicConnection{session, stream}
	} else {
		stream = conn.stream
	}

	n, err := stream.Write(data)
	if err != nil {
		return errors.Wrap(err, "[ send ] failed to write to a stream")
	}

	if n != len(data) {
		return errors.New("[ send ] sent a part of data")
	}

	return nil
}

// Start starts networking.
func (t *quicTransport) Listen(ctx context.Context, started chan struct{}) error {
	log.Debug("Start QUIC transport")
	started <- struct{}{}
	for {
		session, err := t.l.Accept()
		if err != nil {
			<-t.disconnectFinished
			return err
		}

		log.Debugf("accept from: %s", session.RemoteAddr().String())
		go t.handleAcceptedConnection(session)
	}
}

// Stop stops networking.
func (t *quicTransport) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	log.Debug("[ Stop ] Stop QUIC transport")
	t.prepareDisconnect()

	utils.CloseVerbose(t.l)

	for _, conn := range t.connections {
		utils.CloseVerbose(conn.stream)
		utils.CloseVerbose(conn.session)
	}

	utils.CloseVerbose(t.conn)
}

func (t *quicTransport) handleAcceptedConnection(session quic.Session) {
	stream, err := session.AcceptStream()
	if err != nil {
		log.Error(err, "[ handleAcceptedConnection ] failed to get a stream")
	}

	msg, err := t.serializer.DeserializePacket(stream)
	if err != nil {
		log.Error(err, "[ handleAcceptedConnection ] failed to deserialize a packet")
	}

	go t.packetHandler.Handle(context.TODO(), msg)

	utils.CloseVerbose(stream)
}

func createConnection(addr string) (quic.Session, quic.Stream, error) {
	// TODO: NETD18-78
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, nil) //nolint: gosec
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ createConnection ] failed to create a session")
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ createConnection ] failed to open a stream")
	}
	log.Debug("connected to: %s", session.RemoteAddr().String())
	return session, stream, nil
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
