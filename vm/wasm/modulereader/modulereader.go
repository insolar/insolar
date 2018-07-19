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
// Package modulereader implements binary wasm reader statemachine
// that simple reads binary code into internal representation
package modulereader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gitlab.com/insecosystem/vm/iwasm/types"
	"io"
	"io/ioutil"
)

// Reader is generalised module reader context and tokenizer
type Reader struct {
	P uint64
	R io.Reader
}

// Read implements io.Reader
func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.R.Read(p)
	r.P += uint64(n)
	return n, err
}

// ReadByte just one byte
func (r *Reader) ReadByte() (byte, error) {
	p := make([]byte, 1)
	_, err := r.R.Read(p)
	r.P++
	return p[0], err
}

// ReadBytes reads n bytes
func (r *Reader) ReadBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	return b, err
}

// ReadString reads n bytes as string
func (r *Reader) ReadString(n int) (string, error) {
	b, err := r.ReadBytes(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadU32 reads 4 bytes as uint32
func (r *Reader) ReadU32() (uint32, error) {
	var buf [4]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

//
//
// leb128
//

// ReadVarUint32Size reads leb128 uint
func (r *Reader) ReadVarUint32Size() (res uint32, size uint, err error) {
	b := make([]byte, 1)
	var shift uint
	for {
		if _, err = io.ReadFull(r, b); err != nil {
			return
		}
		size++
		cur := uint32(b[0])
		res |= (cur & 0x7f) << (shift)
		if cur&0x80 == 0 {
			return res, size, nil
		}
		shift += 7
	}
}

// ReadVarint64Size reads int64
func (r *Reader) ReadVarint64Size() (res int64, size uint, err error) {
	var shift uint
	var sign int64 = -1
	b := make([]byte, 1)

	for {
		if _, err = io.ReadFull(r, b); err != nil {
			return
		}
		size++

		cur := int64(b[0])
		res |= (cur & 0x7f) << shift
		shift += 7
		sign <<= 7
		if cur&0x80 == 0 {
			break
		}
	}

	if ((sign >> 1) & res) != 0 {
		res |= sign
	}
	return res, size, nil
}

// ReadVarint32Size reads int
func (r *Reader) ReadVarint32Size() (res int32, size uint, err error) {
	res64, size, err := r.ReadVarint64Size()
	res = int32(res64)
	return
}

// ReadVarUint32 reads uint32
func (r *Reader) ReadVarUint32() (uint32, error) {
	n, _, err := r.ReadVarUint32Size()
	return n, err
}

// ReadVarint32 reads int32
func (r *Reader) ReadVarint32() (int32, error) {
	n, _, err := r.ReadVarint32Size()
	return n, err
}

//
//
/// wasm part
//

// ReadValueType reads valuetype
func (r *Reader) ReadValueType() (types.Value, error) {
	v, err := r.ReadVarint32()
	return types.Value(v), err
}

// ReadFunction reads function signature
func (r *Reader) ReadFunction() (types.FunctionSig, error) {
	f := types.FunctionSig{}

	form, err := r.ReadByte()
	if err != nil {
		return f, err
	}

	f.Form = form

	paramCount, err := r.ReadVarUint32()
	if err != nil {
		return f, err
	}
	f.Params = make([]types.Value, paramCount)

	for i := range f.Params {
		f.Params[i], err = r.ReadValueType()
		if err != nil {
			return f, err
		}
	}

	returnCount, err := r.ReadVarUint32()
	if err != nil {
		return f, err
	}

	f.Returns = make([]types.Value, returnCount)
	for i := range f.Returns {
		vt, err := r.ReadValueType()
		if err != nil {
			return f, err
		}
		f.Returns[i] = vt
	}

	return f, nil
}

// ReadTable reads table stub
func (r *Reader) ReadTable() (*types.Table, error) {
	eltype, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	lims, err := r.ReadResizableLimits()
	if err != nil {
		return nil, err
	}

	return &types.Table{
		ElementType: eltype,
		Limits:      *lims,
	}, err
}

// ReadMemory reads memory stub
func (r *Reader) ReadMemory() (*types.Memory, error) {
	lim, err := r.ReadResizableLimits()
	if err != nil {
		return nil, err
	}
	return &types.Memory{Limits: *lim}, nil

}

// ReadResizableLimits reads ResizableLimits object
func (r *Reader) ReadResizableLimits() (*types.ResizableLimits, error) {
	lim := &types.ResizableLimits{
		Maximum: 0,
	}
	f, err := r.ReadVarUint32()
	if err != nil {
		return nil, err
	}

	lim.Flags = f
	lim.Initial, err = r.ReadVarUint32()
	if err != nil {
		return nil, err
	}

	if lim.Flags&0x1 != 0 {
		m, err := r.ReadVarUint32()
		if err != nil {
			return nil, err
		}
		lim.Maximum = m

	}
	return lim, nil
}

// ReadGlobalVar reads global var
func (r *Reader) ReadGlobalVar() (*types.GlobalVar, error) {
	t, err := r.ReadValueType()
	if err != nil {
		return nil, err
	}

	m, err := r.ReadVarUint32()
	if err != nil {
		return nil, err
	}
	return &types.GlobalVar{
		Type:    t,
		Mutable: m == 1,
	}, nil
}

// ReadImportEntry reads one generic import entry
func (r *Reader) ReadImportEntry() (types.Import, error) { // nolint: gocyclo
	i := types.Import{}
	modLen, err := r.ReadVarUint32()
	if err != nil {
		return i, err
	}
	if i.Module, err = r.ReadString(int(modLen)); err != nil {
		return i, err
	}

	fieldLen, err := r.ReadVarUint32()
	if err != nil {
		return i, err
	}

	if i.Field, err = r.ReadString(int(fieldLen)); err != nil {
		return i, err
	}

	v, err := r.ReadByte()
	if err != nil {
		return i, err
	}

	switch types.External(v) {
	case types.ExternalFunction:
		var t uint32
		t, err = r.ReadVarUint32()
		if err != nil {
			return i, err
		}
		i.Type = types.ImportFunc{Type: t}
	case types.ExternalTable:
		var table *types.Table
		table, err = r.ReadTable()
		if err != nil {
			return i, err
		} else if table != nil {
			i.Type = types.ImportTable{Type: *table}
		}
	case types.ExternalMemory:
		var mem *types.Memory
		mem, err = r.ReadMemory()
		if err != nil {
			return i, err
		} else if mem != nil {
			i.Type = types.ImportMemory{Type: *mem}
		}
	case types.ExternalGlobal:
		var gl *types.GlobalVar
		gl, err = r.ReadGlobalVar()
		if err != nil {
			return i, err
		} else if gl != nil {
			i.Type = types.ImportGlobalVar{Type: *gl}
		}
	}
	return i, nil
}

var emptyge = types.GlobalEntry{}

// ReadGlobalEntry reads one entry of globals
func (r *Reader) ReadGlobalEntry() (types.GlobalEntry, error) {
	t, err := r.ReadGlobalVar()
	if err != nil {
		return emptyge, err
	}
	init, err := r.ReadInitExpr()
	if err != nil {
		return emptyge, err
	}

	return types.GlobalEntry{
		Type: t,
		Init: init,
	}, nil
}

// ReadInitExpr reads init code
func (r *Reader) ReadInitExpr() ([]byte, error) {
	buf := new(bytes.Buffer)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == 0x0B { // TODO use instruction constant here and more complicated end finder
			// also ... bad way, 0x0b may be a part of expression ...
			return buf.Bytes(), nil
		}

		if err := buf.WriteByte(b); err != nil {
			return nil, err
		}
	}
}

// ReadExportEntry reads one export entry
func (r *Reader) ReadExportEntry() (types.ExportEntry, error) {
	e := types.ExportEntry{}
	l, err := r.ReadVarUint32()
	if err != nil {
		return e, err
	}
	if e.Name, err = r.ReadString(int(l)); err != nil {
		return e, err
	}

	k, err := r.ReadByte()
	if err != nil {
		return e, err
	}
	e.Kind = types.External(k)

	e.Index, err = r.ReadVarUint32()

	return e, err
}

// ReadElementSegment reads one element
func (r *Reader) ReadElementSegment() (types.ElementSegment, error) {
	s := types.ElementSegment{}
	var err error
	if s.Index, err = r.ReadVarUint32(); err != nil {
		return s, err
	}
	if s.Offset, err = r.ReadInitExpr(); err != nil {
		return s, err
	}

	cnt, err := r.ReadVarUint32()
	if err != nil {
		return s, err
	}
	s.Elems = make([]uint32, cnt)

	for i := range s.Elems {
		e, err := r.ReadVarUint32()
		if err != nil {
			return s, err
		}
		s.Elems[i] = e
	}

	return s, nil
}

// ReadFunctionBody reads whole function
func (r *Reader) ReadFunctionBody() (types.FunctionBody, error) {
	f := types.FunctionBody{}

	bodySize, err := r.ReadVarUint32()
	if err != nil {
		return f, err
	}

	body, err := r.ReadBytes(int(bodySize))
	if err != nil {
		return f, err
	}
	b := &Reader{R: bytes.NewBuffer(body)}

	lcnt, err := b.ReadVarUint32()
	if err != nil {
		return f, err
	}
	f.Locals = make([]types.LocalEntry, lcnt)

	for i := range f.Locals {
		if f.Locals[i], err = b.ReadLocalEntry(); err != nil {
			return f, err
		}
	}

	code, err := ioutil.ReadAll(b)
	if err != nil {
		return f, err
	}
	if code[len(code)-1] != 0x0b { //todo use constants
		return f, errors.New("function have no end")
	}

	f.Code = code

	return f, nil
}

// ReadLocalEntry reads one locals entry
func (r Reader) ReadLocalEntry() (types.LocalEntry, error) {
	l := types.LocalEntry{}
	var err error

	l.Count, err = r.ReadVarUint32()
	if err != nil {
		return l, err
	}

	l.Type, err = r.ReadValueType()
	return l, err
}

// ReadDataSegment reads one segment
func (r *Reader) ReadDataSegment() (types.DataSegment, error) {
	s := types.DataSegment{}
	var err error

	if s.Index, err = r.ReadVarUint32(); err != nil {
		return s, err
	}
	if s.Offset, err = r.ReadInitExpr(); err != nil {
		return s, err
	}

	size, err := r.ReadVarUint32()
	if err != nil {
		return s, err
	}
	s.Data, err = r.ReadBytes(int(size))

	return s, err
}
