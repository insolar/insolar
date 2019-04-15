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

package packet

import (
	"bytes"
	"encoding/gob"
	"github.com/insolar/x-crypto/rand"
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func init() {
	gob.Register(&RequestTest{})
}

func TestSerializePacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	msg := builder.Receiver(receiver).Type(TestPacket).Request(&RequestTest{[]byte{0, 1, 2, 3}}).Build()

	_, err := SerializePacket(msg)

	require.NoError(t, err)
}

func TestDeserializePacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	msg := builder.Receiver(receiver).Type(TestPacket).Request(&RequestTest{[]byte{0, 1, 2, 3}}).Build()

	serialized, _ := SerializePacket(msg)

	var buffer bytes.Buffer

	buffer.Write(serialized)

	deserialized, err := DeserializePacket(&buffer)

	require.NoError(t, err)
	require.Equal(t, deserialized, msg)
}

func TestDeserializeBigPacket(t *testing.T) {
	hostOne, _ := host.NewHost("127.0.0.1:31337")

	data := make([]byte, 1024*1024*10)
	rand.Read(data)

	builder := NewBuilder(hostOne)
	msg := builder.Receiver(hostOne).Type(TestPacket).Request(&RequestTest{data}).Build()

	serialized, err := SerializePacket(msg)
	require.NoError(t, err)

	var buffer bytes.Buffer
	buffer.Write(serialized)

	deserializedMsg, err := DeserializePacket(&buffer)
	require.NoError(t, err)

	deserializedData := deserializedMsg.Data.(*RequestTest).Data
	require.Equal(t, data, deserializedData)
}
