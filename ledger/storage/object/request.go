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

package object

import (
	"io"

	"github.com/insolar/insolar/insolar/record"

	"github.com/insolar/insolar/insolar"
)

// Request extends VirtualRecord interface with GetPayload method.
type Request interface {
	record.VirtualRecord
	GetPayload() []byte
	GetObject() insolar.ID
}

// RequestRecord is a contract execution request.
type RequestRecord struct {
	Parcel      []byte
	MessageHash []byte
	Object      insolar.ID
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *RequestRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(r.MessageHash)
}

// GetPayload returns payload. Required for VirtualRecord interface implementation.
func (r *RequestRecord) GetPayload() []byte {
	return r.Parcel
}

// GetObject returns request object.
func (r *RequestRecord) GetObject() insolar.ID {
	return r.Object
}
