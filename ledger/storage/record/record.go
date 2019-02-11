/*
 *    Copyright 2019 Insolar Technologies
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

// TypeID encodes a record object type.
//go:generate go run gen/type.go
type TypeID uint32

// TypeIDSize is a size of TypeID type.
const TypeIDSize = 4

// Record is base interface for all records.
type Record interface {
	// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
	WriteHashData(w io.Writer) (int, error)
}

func init() {
	// Never change id constants.
	register(10, &GenesisRecord{})
	register(11, &ChildRecord{})
	register(12, &JetRecord{})

	register(20, &RequestRecord{})

	register(30, &ResultRecord{})
	register(31, &TypeRecord{})
	register(32, &CodeRecord{})
	register(33, &ObjectActivateRecord{})
	register(34, &ObjectAmendRecord{})
	register(35, &DeactivationRecord{})
}
