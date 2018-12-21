/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package record

import (
	"io"

	"github.com/insolar/insolar/core"
)

// Request extends Record interface with GetPayload method.
type Request interface {
	Record
	GetPayload() []byte
	GetObject() core.RecordID
}

// RequestRecord is a contract execution request.
type RequestRecord struct {
	Payload []byte
	Object  core.RecordID
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *RequestRecord) WriteHashData(w io.Writer) (int, error) {
	return w.Write(r.Payload)
}

// Type implementation of Record interface.
func (r *RequestRecord) Type() TypeID { return typeCallRequest }

// GetPayload returns payload. Required for Record interface implementation.
func (r *RequestRecord) GetPayload() []byte {
	return r.Payload
}

// GetObject returns request object.
func (r *RequestRecord) GetObject() core.RecordID {
	return r.Object
}
