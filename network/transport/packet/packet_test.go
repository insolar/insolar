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

package packet

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"testing"

	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
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

	assert.NoError(t, err)
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

	assert.NoError(t, err)
	assert.Equal(t, deserialized, msg)
}

func TestDeserializeBigPacket(t *testing.T) {
	hostOne, _ := host.NewHost("127.0.0.1:31337")

	data := make([]byte, 1024*1024*10)
	rand.Read(data)

	builder := NewBuilder(hostOne)
	msg := builder.Receiver(hostOne).Type(TestPacket).Request(&RequestTest{data}).Build()

	serialized, err := SerializePacket(msg)
	assert.NoError(t, err)

	var buffer bytes.Buffer
	buffer.Write(serialized)

	deserializedMsg, err := DeserializePacket(&buffer)
	assert.NoError(t, err)

	deserializedData := deserializedMsg.Data.(*RequestTest).Data
	assert.Equal(t, data, deserializedData)
}
