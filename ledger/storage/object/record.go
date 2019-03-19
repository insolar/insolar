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

package object

import (
	"io"
)

//go:generate go run gen/type.go

// TypeID encodes a record object type.
type TypeID uint32

// TypeIDSize is a size of TypeID type.
const TypeIDSize = 4

// VirtualRecord is base interface for all records.
type VirtualRecord interface {
	// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
	WriteHashData(w io.Writer) (int, error)
}

func init() {
	// ID can be any unique int value.
	// Never change id constants. They are used for serialization.
	register(100, new(GenesisRecord))
	register(101, new(ChildRecord))
	register(102, new(JetRecord))

	register(200, new(RequestRecord))

	register(300, new(ResultRecord))
	register(301, new(TypeRecord))
	register(302, new(CodeRecord))
	register(303, new(ActivateRecord))
	register(304, new(AmendRecord))
	register(305, new(DeactivationRecord))
}
