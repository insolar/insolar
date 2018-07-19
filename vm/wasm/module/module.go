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
// Package module represents wasm module
package module

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/insolar/insolar/vm/wasm/modulereader"
	"github.com/insolar/insolar/vm/wasm/types"
)

// Misc constants
const (
	Magic   uint32 = 0x6d736100
	Version uint32 = 0x1
)

// Module is a struct ...
type Module struct {
	Version uint32

	Types    *SectionTypes
	Import   *SectionImports
	Function *SectionFunctions
	Table    *SectionTables
	Memory   *SectionMemories
	Global   *SectionGlobals
	Export   *SectionExports
	Start    *SectionStartFunction
	Elements *SectionElements
	Code     *SectionCode
	Data     *SectionData
}

// SectionID is a type of binary section
type SectionID byte

// Exported constants for sections
const (
	SectionIDCustom   SectionID = 0
	SectionIDType     SectionID = 1
	SectionIDImport   SectionID = 2
	SectionIDFunction SectionID = 3
	SectionIDTable    SectionID = 4
	SectionIDMemory   SectionID = 5
	SectionIDGlobal   SectionID = 6
	SectionIDExport   SectionID = 7
	SectionIDStart    SectionID = 8
	SectionIDElement  SectionID = 9
	SectionIDCode     SectionID = 10
	SectionIDData     SectionID = 11
)

// Section is a basic info about section
type Section struct {
	Begin uint64
	End   uint64
	ID    SectionID
	// Size of this section in bytes
	Len   uint32
	Bytes []byte
}

// Next types represents actual sections data

// SectionCustom SectionCustom
type SectionCustom struct {
	Section
}

// SectionTypes SectionTypes
type SectionTypes struct {
	Section
	Entries []types.FunctionSig
}

// SectionImports SectionImports
type SectionImports struct {
	Section
	Entries []types.Import
}

// SectionFunctions SectionImports
type SectionFunctions struct {
	Section
	Types []uint32
}

// SectionTables SectionTables
type SectionTables struct {
	Section
	Entries []types.Table
}

// SectionMemories SectionMemories
type SectionMemories struct {
	Section
	Entries []types.Memory
}

// SectionGlobals SectionGlobals
type SectionGlobals struct {
	Section
	Globals []types.GlobalEntry
}

// SectionExports SectionExports
type SectionExports struct {
	Section
	Entries map[string]types.ExportEntry
}

// SectionStartFunction SectionStartFunction
type SectionStartFunction struct {
	Section
	Index uint32
}

// SectionElements SectionElements
type SectionElements struct {
	Section
	Entries []types.ElementSegment
}

// SectionCode SectionCode
type SectionCode struct {
	Section
	Bodies []types.FunctionBody
}

// SectionData SectionData
type SectionData struct {
	Section
	Entries []types.DataSegment
}

// NewModule creates empty module
func NewModule() *Module {
	return &Module{
		Types:    &SectionTypes{},
		Import:   &SectionImports{},
		Table:    &SectionTables{},
		Memory:   &SectionMemories{},
		Global:   &SectionGlobals{},
		Export:   &SectionExports{},
		Start:    &SectionStartFunction{},
		Elements: &SectionElements{},
		Data:     &SectionData{},
	}
}

// Read reads module from stream
func Read(input io.Reader) (*Module, error) {
	m := NewModule()
	r := &modulereader.Reader{R: input}
	magic, err := r.ReadU32()
	if err != nil {
		return nil, err
	}
	if magic != Magic {
		return nil, errors.New("wasm: Invalid magic number")
	}

	if m.Version, err = r.ReadU32(); err != nil {
		return nil, err
	}

	for {
		done, err := m.readSection(r)
		if done {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Module) readSection(r *modulereader.Reader) (done bool, err error) {
	var id uint32
	if id, err = r.ReadVarUint32(); err != nil {
		done = err == io.EOF
		return done, err
	}

	s := Section{ID: SectionID(id)}
	if s.Len, err = r.ReadVarUint32(); err != nil {
		return done, err
	}

	length := s.Len
	s.Begin = r.P
	b := make([]byte, length)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return false, err
	}
	sr := modulereader.Reader{R: bytes.NewReader(b)}
	s.End = r.P
	s.Bytes = b

	if s.ID >= sectionReadersLen {
		return false, fmt.Errorf("wrong section ID %d", s.ID)
	}

	err = sectionReaders[s.ID](m, &sr, s)

	return false, err
}

var sectionReaders = []func(m *Module, r *modulereader.Reader, bs Section) error{

	// 0 - Custom
	func(m *Module, r *modulereader.Reader, bs Section) error {
		return nil
	},

	// 1 - Section Type
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionTypes{Section: bs}
		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}

		s.Entries = make([]types.FunctionSig, int(count))

		for i := range s.Entries {
			if s.Entries[i], err = r.ReadFunction(); err != nil {
				return err
			}
		}

		m.Types = s
		return nil
	},

	// 2 - import
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionImports{Section: bs}
		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}
		s.Entries = make([]types.Import, count)

		for i := range s.Entries {
			s.Entries[i], err = r.ReadImportEntry()
			if err != nil {
				return err
			}
		}

		m.Import = s
		return nil
	},

	// 3 - functions
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionFunctions{Section: bs}
		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}

		s.Types = make([]uint32, count)

		for i := range s.Types {
			t, err := r.ReadVarUint32()
			if err != nil {
				return err
			}
			s.Types[i] = t
		}

		m.Function = s
		return nil
	},

	// 4 - tables
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionTables{Section: bs}
		count, err := r.ReadVarUint32()
		if err != nil {
			return err

		}
		s.Entries = make([]types.Table, count)

		for i := range s.Entries {
			t, err := r.ReadTable()
			if err != nil {
				return err
			}
			s.Entries[i] = *t
		}

		m.Table = s
		return nil
	},

	// 5 - memories
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionMemories{Section: bs}
		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}

		s.Entries = make([]types.Memory, count)

		for i := range s.Entries {
			m, err := r.ReadMemory()
			if err != nil {
				return err
			}
			s.Entries[i] = *m
		}

		m.Memory = s
		return nil
	},

	// 6 - globals
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionGlobals{Section: bs}

		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}

		s.Globals = make([]types.GlobalEntry, count)

		for i := range s.Globals {
			s.Globals[i], err = r.ReadGlobalEntry()
			if err != nil {
				return err
			}
		}

		m.Global = s
		return nil
	},

	// 7 - exports
	func(m *Module, r *modulereader.Reader, bs Section) error {
		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}
		s := &SectionExports{Section: bs}
		s.Entries = make(map[string]types.ExportEntry, count)

		for i := uint32(0); i < count; i++ {
			entry, err := r.ReadExportEntry()
			if err != nil {
				return err
			}

			if _, exists := s.Entries[entry.Name]; exists {
				return errors.New("Duplicated export" + entry.Name) // todo use error.Wrap here
			}
			s.Entries[entry.Name] = entry
		}

		m.Export = s
		return nil
	},

	// 8 - start
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionStartFunction{Section: bs}
		var err error

		s.Index, err = r.ReadVarUint32()
		if err != nil {
			return err
		}

		m.Start = s
		return nil
	},

	// 9 - Elements
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionElements{Section: bs}
		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}
		s.Entries = make([]types.ElementSegment, count)

		for i := range s.Entries {
			s.Entries[i], err = r.ReadElementSegment()
			if err != nil {
				return err
			}
		}

		m.Elements = s
		return nil
	},

	// 10 - code
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionCode{Section: bs}

		count, err := r.ReadVarUint32()
		if err != nil {
			return err
		}
		s.Bodies = make([]types.FunctionBody, count)

		for i := range s.Bodies {
			if s.Bodies[i], err = r.ReadFunctionBody(); err != nil {
				return err
			}
		}

		m.Code = s
		if m.Function == nil || len(m.Function.Types) == 0 {
			return errors.New("empty function section")
		}

		return nil
	},

	// 11 - data
	func(m *Module, r *modulereader.Reader, bs Section) error {
		s := &SectionData{Section: bs}
		cnt, err := r.ReadVarUint32()
		if err != nil {
			return err
		}

		s.Entries = make([]types.DataSegment, cnt)

		for i := range s.Entries {
			if s.Entries[i], err = r.ReadDataSegment(); err != nil {
				return err
			}
		}

		m.Data = s
		return err
	},
}
var sectionReadersLen = SectionID(len(sectionReaders))

//

//

//

//

//

// tail for editing
