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
