package servicenetwork

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

// serializeMessage returns io.Reader on buffer with encoded message.Message (from watermill).
func serializeMessage(msg *message.Message) ([]byte, error) {
	wm := &WatermillMessage{
		UUID:     msg.UUID,
		Metadata: msg.Metadata,
		Payload:  msg.Payload,
	}
	return wm.Marshal()
}

// deserializeMessage returns decoded signed message.
func deserializeMessage(data []byte) (*message.Message, error) {
	var wm WatermillMessage
	err := wm.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return &message.Message{
		UUID:     wm.UUID,
		Metadata: wm.Metadata,
		Payload:  wm.Payload,
	}, nil
}
