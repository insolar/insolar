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
Package packet provides network messaging protocol and serialization layer.

Packet can be created with shortcut:

	senderAddr, _ := host.NewAddress("127.0.0.1:1337")
	receiverAddr, _ := host.NewAddress("127.0.0.1:1338")

	sender := host.NewHost(senderAddr)
	receiver := host.NewHost(receiverAddr)

	msg := packet.NewPingPacket(sender, receiver)

	// do something with packet


Or with builder:

	builder := packet.NewBuilder()

	senderAddr, _ := host.NewAddress("127.0.0.1:1337")
	receiverAddr, _ := host.NewAddress("127.0.0.1:1338")

	sender := host.NewHost(senderAddr)
	receiver := host.NewHost(receiverAddr)

	msg := builder.
		Sender(sender).
		Receiver(receiver).
		Type(packet.TypeFindHost).
		Request(&packet.RequestDataFindHost{}).
		Build()

	// do something with packet


Packet may be serialized:

	msg := &packet.Packet{}
	serialized, err := packet.SerializePacket(msg)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(serialized)


And deserialized therefore:

	var buffer bytes.Buffer

	// Fill buffer somewhere

	msg, err := packet.DeserializePacket(buffer)

	if err != nil {
		panic(err.Error())
	}

	// do something with packet

*/
package packet
