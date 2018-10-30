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

package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeSigned(t *testing.T) {
	msg := &SetRecord{
		Record: []byte{0x0A},
	}
	signMsgIn := &SignedMessage{
		Msg:       msg,
		Signature: nil,
	}

	buff, err := SignedToBytes(signMsgIn)
	assert.NoError(t, err)

	signMsgOut, err := DeserializeSigned(bytes.NewBuffer(buff))
	assert.NoError(t, err)

	assert.Equal(t, signMsgIn, signMsgOut)
	assert.Equal(t, signMsgIn.Message(), signMsgOut.Message())
}

func TestSerializeSignedFail(t *testing.T) {
	msg := &SetRecord{
		Record: []byte{0x0A},
	}

	signMsgIn := &SignedMessage{
		Msg:       msg,
		Signature: nil,
	}
	buff, err := ToBytes(signMsgIn)
	assert.NoError(t, err)

	signMsgOut, err := Deserialize(bytes.NewBuffer(buff))
	assert.Error(t, err)
	assert.Nil(t, signMsgOut)
}
