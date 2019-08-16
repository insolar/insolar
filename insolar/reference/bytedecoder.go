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

package reference

import (
	"encoding/base64"
	"io"

	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type ByteDecoderFactory interface {
	GetByteDecoder(encodingName string) ByteDecodeFunc
	DefaultDecoder() ByteDecodeFunc
	LegacyDecoder() ByteDecodeFunc
}

type byteDecoderFactory struct {
	decoders map[string]ByteDecodeFunc
}

func byteDecodeBase58(s string, target io.ByteWriter) (stringRead int, err error) {
	bytes := base58.Decode(s)
	if len(s) > 0 && len(bytes) == 0 {
		return 0, errors.New("input string contains bad charachters")
	}

	for _, b := range bytes {
		err := target.WriteByte(b)
		if err != nil {
			return 0, errors.Wrap(err, "failed to write byte")
		}
	}
	return len(s), nil
}

func byteDecodeBase64(s string, target io.ByteWriter) (stringRead int, err error) {
	bytes, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}

	for _, b := range bytes {
		err := target.WriteByte(b)
		if err != nil {
			return 0, errors.Wrap(err, "failed to write byte")
		}
	}
	return len(s), nil
}

func NewByteDecoderFactory() ByteDecoderFactory {
	return &byteDecoderFactory{
		decoders: map[string]ByteDecodeFunc{
			"base58": byteDecodeBase58,
			"base64": byteDecodeBase64,
		},
	}
}

func (f *byteDecoderFactory) GetByteDecoder(encodingName string) ByteDecodeFunc {
	return f.decoders[encodingName]
}

func (f *byteDecoderFactory) DefaultDecoder() ByteDecodeFunc {
	return f.decoders["base58"]
}

func (f *byteDecoderFactory) LegacyDecoder() ByteDecodeFunc {
	return f.decoders["base58"]
}
