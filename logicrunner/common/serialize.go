//
// Copyright 2019 Insolar Technologies GmbH
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
//

package common

import (
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/log"
)

type Serializer interface {
	Serialize(interface{}, *[]byte) error
	Deserialize([]byte, interface{}) error
}

type CBORSerializer struct{}

func (s *CBORSerializer) Serialize(what interface{}, to *[]byte) error {
	ch := new(codec.CborHandle)
	log.Debugf("serializing %+v", what)
	return codec.NewEncoderBytes(to, ch).Encode(what)
}

func (s *CBORSerializer) Deserialize(from []byte, into interface{}) error {
	ch := new(codec.CborHandle)
	log.Debugf("de-serializing %+v", from)
	return codec.NewDecoderBytes(from, ch).Decode(into)
}

func NewCBORSerializer() *CBORSerializer {
	return &CBORSerializer{}
}
