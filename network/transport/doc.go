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

/*
Package transport provides network transport interface. It allows to abstract our network from physical transport.
It can either be IP based network or any other kind of packet courier (e.g. an industrial packet bus).

Package exports simple interfaces for easily defining new transports.

For now we provide two implementations of transport.
The default is UTPTransport which using BitTorrent µTP protocol.

Usage:

	var conn net.PacketConn
	// get udp connection anywhere

	tp, _ := transport.NewUTPTransport(conn)
	msg := &packet.Packet{}

	// Send the async queries and wait for a future
	future, err := tp.SendRequest(msg)
	if err != nil {
		panic(err)
	}

	select {
	case response := <-future.Result():
		// Channel was closed
		if response == nil {
			panic("chanel closed unexpectedly")
		}

		// do something with response

	case <-time.After(1 * time.Second):
		future.Cancel()
	}

*/
package transport
