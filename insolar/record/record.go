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

// /go:generate protoc -I./vendor -I./ --gogoslick_out=./ insolar/record/record.proto
package record

import (
	"github.com/insolar/insolar/insolar"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

type SerializableRecord interface {
	MarshalRecord() ([]byte, error)
}

// VirtualRecord is base interface for all records.
// TODO: when migrating to protobuf - embed SerializableRecord into VirtualRecord
type VirtualRecord interface {
	// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
	WriteHashData(w io.Writer) (int, error)
}

type MaterialRecord struct {
	Record VirtualRecord

	JetID insolar.JetID
}

func (m *GenesisRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_Genesis{m}
	return base.Marshal()
}
func (m *GenesisRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *ChildRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_Child{m}
	return base.Marshal()
}
func (m *ChildRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *JetRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_Jet{m}
	return base.Marshal()
}
func (m *JetRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *RequestRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_Request{m}
	return base.Marshal()
}
func (m *RequestRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}


func (m *ResultRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_Result{m}
	return base.Marshal()
}
func (m *ResultRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *TypeRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_Type{m}
	return base.Marshal()
}
func (m *TypeRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *CodeRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_Code{m}
	return base.Marshal()
}
func (m *CodeRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *ObjectActivateRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_ObjectActivate{m}
	return base.Marshal()
}
func (m *ObjectActivateRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *ObjectAmendRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_ObjectAmend{m}
	return base.Marshal()
}
func (m *ObjectAmendRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

func (m *ObjectDeactivateRecord) MarshalRecord() ([]byte, error) {
	base := Record{}
	base.Union = &Record_ObjectDeactivate{m}
	return base.Marshal()
}
func (m *ObjectDeactivateRecord) WriteHashData(w io.Writer) (int, error) {
	bytes, err := m.MarshalRecord()
	if err != nil {
		return 0, err
	}
	return w.Write(bytes)
}

// Returns any sub-record type or error
func UnmarshalRecord(data []byte) (interface{}, error) {
	base := Record{}

	if error := base.Unmarshal(data); error != nil {
		return nil, errors.Wrap(error, "Failed to unmarshal request")
	}

	union := base.GetUnion()
	if union == nil {
		return nil, errors.New("We got empty request")
	}

	// using reflection to get real value, instead of oneOf wrapper,
	// needed to implement oneOf logic. Can be written using big switch:case,
	// but it'll be uglier (i guess)
	field := reflect.ValueOf(union).Field(1).Interface()
	return field, nil
}
