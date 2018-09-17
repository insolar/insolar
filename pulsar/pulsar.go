/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package pulsar

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/insolar/insolar/configuration"
)

type Pulsar struct {
	Sock       net.Listener
	Neighbours map[string]*Neighbour
	PrivateKey *rsa.PrivateKey
}

// Creation new pulsar-node
func NewPulsar(configuration configuration.Pulsar, listener func(string, string) (net.Listener, error)) (*Pulsar, error) {
	// Listen for incoming connections.
	l, err := listener(configuration.ConnectionType.String(), configuration.ListenAddress)
	if err != nil {
		return nil, err
	}

	reader := rand.Reader
	bitSize := 2048

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}
	pulsar := &Pulsar{Sock: l, Neighbours: map[string]*Neighbour{}}
	pulsar.PrivateKey = privateKey

	for _, neighbour := range configuration.NodesAddresses {
		pulsar.Neighbours[neighbour.Address] = &Neighbour{ConnectionType: neighbour.ConnectionType}
	}

	gob.Register(Message{})
	gob.Register(HandshakeMessageBody{})

	return pulsar, nil
}

// Listen port for input connections
func (pulsar *Pulsar) Listen() {
	for {
		// Listen for an incoming connection.
		conn, err := pulsar.Sock.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println(conn.RemoteAddr().String())
		if _, ok := pulsar.Neighbours[conn.RemoteAddr().String()]; !ok {
			conn.Close()
			return
		}
		// Handle connections in a new goroutine.
		go pulsar.handleNewRequest(conn)
	}
}

// Connect to all known nodes
func (pulsar *Pulsar) ConnectToAllNeighbours() error {
	for key, neighbour := range pulsar.Neighbours {
		err := pulsar.ConnectToNeighbour(key, neighbour.ConnectionType.String())
		if err != nil {
			return err
		}
	}

	return nil
}

// Connect to the concrete member
func (pulsar *Pulsar) ConnectToNeighbour(address string, connectionType string) error {
	conn, err := net.Dial(connectionType, address)
	if err != nil {
		return err
	}
	conn.(*net.TCPConn).SetKeepAlive(true)
	pulsar.Neighbours[address].Connection = conn
	pulsar.Send(address, &HandshakeMessageBody{PublicKey: pulsar.PrivateKey.PublicKey})
	go pulsar.handleNewRequest(conn)

	return nil
}

func (pulsar *Pulsar) Send(address string, data interface{}) error {
	return gob.NewEncoder(pulsar.Neighbours[address].Connection).Encode(data)
}

// Close all connections
func (pulsar *Pulsar) Close() {
	for _, neighbour := range pulsar.Neighbours {
		if neighbour.Connection != nil {
			neighbour.Connection.Close()
		}
	}

	pulsar.Sock.Close()
}

// Handles incoming requests.
func (pulsar *Pulsar) handleNewRequest(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	message := &Message{}
	err := dec.Decode(message)
	if err == io.EOF {
		remoteAddr := conn.RemoteAddr().String()
		pulsar.ConnectToNeighbour(remoteAddr, pulsar.Neighbours[remoteAddr].ConnectionType.String())
		return
	}

	switch message.Type {
	case Handshake:
		{
			remoteAddr := conn.RemoteAddr()
			if savedConn, ok := pulsar.Neighbours[remoteAddr.String()]; ok {
				if savedConn.Connection != nil {
					savedConn.Connection.Close()
				}
				conn.(*net.TCPConn).SetKeepAlive(true)
				pulsar.Neighbours[remoteAddr.String()].Connection = conn
				messageBody := message.Data.(HandshakeMessageBody)
				pulsar.Neighbours[remoteAddr.String()].PublicKey = &messageBody.PublicKey
				err := pulsar.Send(remoteAddr.String(), Message{Type: Handshake, Data: &HandshakeMessageBody{PublicKey: pulsar.PrivateKey.PublicKey}})
				if err != nil {
					fmt.Println("Error accepting: ", err.Error())
				}
				go pulsar.handleNewRequest(conn)
			}
		}
	}
}
