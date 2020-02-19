// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package servicenetwork

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/require"
)

func TestSerializeDeserialize(t *testing.T) {
	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.Metadata.Set("testKey", "testValue")

	serializedMsg, err := serializeMessage(msg)

	require.NoError(t, err)
	require.NotEmpty(t, serializedMsg)

	msgOut, err := deserializeMessage(serializedMsg)
	require.NoError(t, err)
	require.NotEmpty(t, msgOut)

	require.Equal(t, msg.Payload, msgOut.Payload)
	require.Equal(t, msg.Metadata, msgOut.Metadata)
}
