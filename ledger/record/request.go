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
)

// Request extends Record interface with GetPayload method.
type Request interface {
	Record
	GetPayload() []byte
}

// CallRequest is a contract execution request.
type CallRequest struct {
	Payload []byte
}

// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
func (r *CallRequest) WriteHashData(w io.Writer) (int, error) {
	return w.Write(r.Payload)
}

// Type implementation of Record interface.
func (r *CallRequest) Type() TypeID { return typeCallRequest }

// GetPayload returns payload. Required for Record interface implementation.
func (r *CallRequest) GetPayload() []byte {
	return r.Payload
}
