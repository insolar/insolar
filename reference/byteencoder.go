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
	"bytes"
	"encoding/base64"
	"io"
	"strings"

	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

type ByteEncodeFunc func(source io.ByteReader, builder *strings.Builder) error

func byteEncodeBase58(source io.ByteReader, builder *strings.Builder) error {
	buff := bytes.Buffer{}
	for b, err := source.ReadByte(); err == nil; b, err = source.ReadByte() {
		err := buff.WriteByte(b)
		if err != nil {
			return errors.Wrap(err, "failed to write base58 encoded data to string builder")
		}
	}
	_, err := builder.Write([]byte(base58.Encode(buff.Bytes())))
	return err
}

func byteEncodeBase64(source io.ByteReader, builder *strings.Builder) error {
	buff := bytes.Buffer{}
	for b, err := source.ReadByte(); err == nil; b, err = source.ReadByte() {
		buff.WriteByte(b)
	}
	encoder := base64.NewEncoder(base64.RawURLEncoding, builder)
	_, err := encoder.Write(buff.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to write base64 encoded data to string builder")
	}
	err = encoder.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close string builder")
	}
	return nil
}
