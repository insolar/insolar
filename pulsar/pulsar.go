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
	"bytes"
	"fmt"
	"github.com/insolar/insolar/configuration"
	"net"
	"os"
	"strconv"
)

type Pulsar struct {
	Sock       net.Listener
	Neighbours map[string]net.Conn
}

func NewPulsar(configuration configuration.Pulsar) *Pulsar {
	// Listen for incoming connections.
	l, err := net.Listen(configuration.Type, configuration.ListenAddress)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	return &Pulsar{Sock: l, Neighbours: map[string]net.Conn{}}
}

func (pulsar *Pulsar) Listen() {
	for {
		// Listen for an incoming connection.
		conn, err := pulsar.Sock.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func (pulsar *Pulsar) ConnectToNeighbour(address string, connectionType string) error {
	conn, err := net.Dial(connectionType, address)
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		return err
	}
	pulsar.Neighbours[address] = conn

	return nil
}

func (pulsar *Pulsar) Send(address string, data interface{}) {
}

func (pulsar *Pulsar) Close() {
	for _, connection := range pulsar.Neighbours {
		connection.Close()
	}

	pulsar.Sock.Close()
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Builds the message.
	message := "Hi, I received your message! It was "
	message += strconv.Itoa(reqLen)
	message += " bytes long and that's what it said: \""
	n := bytes.Index(buf, []byte{0})
	message += string(buf[:n-1])
	message += "\" ! Honestly I have no clue about what to do with your messages, so Bye Bye!\n"

	// Write the message in the connection channel.
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Close the connection when you're done with it.
	conn.Close()
}
