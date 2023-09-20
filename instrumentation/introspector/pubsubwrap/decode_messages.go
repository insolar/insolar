package pubsubwrap

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/payload"
)

type decodeError struct {
	metadataType string
	err          error
}

func (de decodeError) Error() string {
	return fmt.Sprintf("can't decode message type: %v, error: %v",
		de.metadataType, de.err.Error())
}

// decodeType tries to decode message.Message as protobuf, return annotated error with type of legacy message.
// ignore protobuf decoding errors, it will happen until legacy messages exist
func decodeType(m *message.Message) (payload.Type, error) {
	var meta payload.Meta
	err := meta.Unmarshal(m.Payload)
	if err != nil {
		return payload.TypeUnknown, decodeError{
			metadataType: m.Metadata["type"],
			err:          err,
		}
	}

	typ, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		return payload.TypeUnknown, decodeError{
			metadataType: m.Metadata["type"],
			err:          err,
		}
	}

	return typ, nil
}
