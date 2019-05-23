package payload

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewMessage(pl Payload) (*message.Message, error) {
	buf, err := Marshal(pl)
	if err != nil {
		return nil, err
	}
	return message.NewMessage(watermill.NewUUID(), buf), nil
}
