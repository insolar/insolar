/*
*    Copyright 2018 INS Ecosystem
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
package types

// Value is a type of simple value
type Value byte

// simple value types
const (
	I32 Value = 0x7F
	I64 Value = 0x7E
	F32 Value = 0x7C
	F64 Value = 0x7D
)

// FunctionSig is a signature of function.
type FunctionSig struct {
	Form    byte
	Params  []Value
	Returns []Value
}

// External is a type of external object.
type External byte

// External entries types
const (
	ExternalFunction External = 0
	ExternalTable    External = 1
	ExternalMemory   External = 2
	ExternalGlobal   External = 3
)

// Importable is an interface of object that we want to import.
type Importable interface{}

// Import is an imported object
type Import struct {
	Module string
	Field  string
	Type   Importable
}

// ImportFunc is an imported function.
type ImportFunc struct {
	Type uint32
}

// ImportTable is an imported table.
type ImportTable struct {
	Type Table
}

// ImportMemory is an imported memory
type ImportMemory struct {
	Type Memory
}

// ImportGlobalVar is an imported global var.
type ImportGlobalVar struct {
	Type GlobalVar
}

// ResizableLimits is a limit boundaries of other objects
type ResizableLimits struct {
	Flags   uint32 // 1 if the Maximum field is valid
	Initial uint32 // initial length (in units of table elements or wasm pages)
	Maximum uint32 // If flags is 1, it describes the maximum size of the table or memory
}

// Table represents a table.
type Table struct {
	ElementType byte
	Limits      ResizableLimits
}

// Memory represents a memory.
type Memory struct {
	Limits ResizableLimits
}

// GlobalVar represents a global var.
type GlobalVar struct {
	Type    Value // Type of the value stored by the variable
	Mutable bool  // Whether the value of the variable can be changed by the set_global operator
}

// GlobalEntry is an entry in globals
type GlobalEntry struct {
	Type *GlobalVar // Type holds information about the value type and mutability of the variable
	Init []byte     // Init is an initializer expression that computes the initial value of the variable
}

// LocalEntry is an entry in locals.
type LocalEntry struct {
	Count uint32 // The total number of local variables of the given Type used in the function body
	Type  Value  // The type of value stored by the variable
}

// ExportEntry is an export entry.
type ExportEntry struct {
	Name  string
	Kind  External
	Index uint32
}

// ElementSegment is an element.
type ElementSegment struct {
	Index  uint32 // The index into the global table space, should always be 0 in the MVP.
	Offset []byte // initializer expression for computing the offset for placing elements, should return an i32 value
	Elems  []uint32
}

// FunctionBody is an bytes and args of a function.
type FunctionBody struct {
	Locals []LocalEntry
	Code   []byte
}

// DataSegment an entry in data section.
type DataSegment struct {
	Index  uint32 // The index into the global linear memory space, should always be 0 in the MVP.
	Offset []byte // initializer expression for computing the offset for placing elements, should return an i32 value
	Data   []byte
}
