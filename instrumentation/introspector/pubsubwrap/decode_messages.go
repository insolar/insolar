// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
