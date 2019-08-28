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

// Package message represents message that messagebus can route
package message

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"

	"github.com/insolar/insolar/insolar"
)

// MustSerialize serialize a insolar.Message to bytes.
func MustSerialize(msg insolar.Message) []byte {
	r := insolar.MustSerialize(msg)
	r = append([]byte{byte(msg.Type())}, r...)
	return r
}

// SerializeParcel returns io.Reader on buffer with encoded insolar.Parcel.
func SerializeParcel(parcel insolar.Parcel) (io.Reader, error) {
	buff := &bytes.Buffer{}
	enc := gob.NewEncoder(buff)
	err := enc.Encode(parcel)
	return buff, err
}

// ParcelToBytes deserialize a insolar.Parcel to bytes.
func ParcelToBytes(msg insolar.Parcel) []byte {
	reqBuff, err := SerializeParcel(msg)
	if err != nil {
		panic("failed to serialize message: " + err.Error())
	}
	buf, err := ioutil.ReadAll(reqBuff)
	if err != nil {
		panic("failed to serialize message: " + err.Error())
	}
	return buf
}

func init() {
	gob.Register(&Parcel{})
	gob.Register(insolar.Reference{})
}
