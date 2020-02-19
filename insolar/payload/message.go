// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package payload

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar/bus/meta"
)

func NewMessage(pl Payload) (*message.Message, error) {
	buf, err := Marshal(pl)
	if err != nil {
		return nil, err
	}
	return message.NewMessage(watermill.NewUUID(), buf), nil
}

func MustNewMessage(pl Payload) *message.Message {
	msg, err := NewMessage(pl)
	if err != nil {
		panic(err)
	}
	return msg
}

func NewResultMessage(pl Payload) (*message.Message, error) {
	msg, err := NewMessage(pl)
	if err != nil {
		return nil, err
	}
	msg.Metadata.Set(meta.Type, meta.TypeReturnResults)
	return msg, nil
}
