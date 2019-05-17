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

package bus

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

func serializeError(e *serializableError) (io.Reader, error) {
	buff := &bytes.Buffer{}
	enc := gob.NewEncoder(buff)
	err := enc.Encode(e)
	return buff, err
}

// DeserializeError returns decoded error.
func DeserializeError(buff io.Reader) (error, error) {
	var e serializableError
	enc := gob.NewDecoder(buff)
	err := enc.Decode(&e)
	return &e, err
}

// ErrorToBytes deserialize error to bytes.
func ErrorToBytes(e error) ([]byte, error) {
	errMsg := &serializableError{
		S: e.Error(),
	}
	reqBuff, err := serializeError(errMsg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize message")
	}
	buf, err := ioutil.ReadAll(reqBuff)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from buffer")
	}
	return buf, nil
}
