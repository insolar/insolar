// Copyright 2020 Insolar Network Ltd.
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

// +build ignore

package logbuf

import (
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/insolar"
)

func NewJoiningLevelBuffer(buffer PlainBuffer, filler []byte) LevelBuffer {
	if len(filler) < serviceHeaderSizeMin {
		panic("illegal value")
	}

	return LevelBuffer{
		buffer:              buffer,
		serviceHeaderFiller: filler,
		serviceHeaderSize:   uint32(len(filler)),
	}
}

func NewChunkingLevelBuffer(buffer PlainBuffer) LevelBuffer {
	return LevelBuffer{
		buffer:            buffer,
		serviceHeaderSize: serviceHeaderSizeMin,
	}
}

//var _ insolar.LogLevelWriter = &LevelBuffer{}

type LevelBuffer struct {
	buffer PlainBuffer

	serviceHeaderFiller []byte
	serviceHeaderSize   uint32
}

const serviceHeaderSizeMin = 5

var byteOrder = binary.LittleEndian

func (p *LevelBuffer) LogLevelWrite(level insolar.LogLevel, b []byte) (int, error) {
	segmentLen := p.serviceHeaderSize + uint32(len(b))

	pg, buf := p.buffer.buffer.allocateBuffer(segmentLen)

	byteOrder.PutUint32(buf, segmentLen)
	buf[p.serviceHeaderSize-1] = byte(level)
	copy(buf[p.serviceHeaderSize:], b)
	pg.stopAccess()
	return len(b), nil
}

func (p *LevelBuffer) Write(b []byte) (n int, err error) {
	return p.LogLevelWrite(insolar.NoLevel, b)
}

func (p *LevelBuffer) LevelWriteTo(w insolar.LogLevelWriter) (int64, error) {
	return p.buffer.WriteTo(chunkingLevelWriter{w, p.serviceHeaderSize})
}

func (p *LevelBuffer) WriteTo(w io.Writer) (int64, error) {
	if p.serviceHeaderFiller != nil {
		return p.buffer.WriteTo(joiningWriter{w, p.serviceHeaderFiller})
	}
	return p.buffer.WriteTo(chunkingWriter{w, p.serviceHeaderSize})
}

/* ============================ */

type chunkingLevelWriter struct {
	w                 insolar.LogLevelWriter
	serviceHeaderSize uint32
}

func (w chunkingLevelWriter) Write(b []byte) (int, error) {
	totalN := 0
	pos := uint32(0)
	max := uint32(len(b))
	for pos < max {
		chunkLen := byteOrder.Uint32(b[pos:])
		level := insolar.LogLevel(b[pos+w.serviceHeaderSize-1])
		n, err := w.w.LogLevelWrite(level, b[pos+w.serviceHeaderSize:pos+chunkLen])
		totalN += n
		if err != nil {
			return totalN, err
		}
		pos += chunkLen
	}
	return totalN, nil
}

type chunkingWriter struct {
	w                 io.Writer
	serviceHeaderSize uint32
}

func (w chunkingWriter) Write(b []byte) (int, error) {
	totalN := 0
	pos := uint32(0)
	max := uint32(len(b))
	for pos < max {
		chunkLen := byteOrder.Uint32(b[pos:])
		n, err := w.w.Write(b[pos+w.serviceHeaderSize : pos+chunkLen])
		totalN += n
		if err != nil {
			return totalN, err
		}
		pos += chunkLen
	}
	return totalN, nil
}

type joiningWriter struct {
	w      io.Writer
	filler []byte
}

func (w joiningWriter) Write(b []byte) (int, error) {
	pos := uint32(0)
	max := uint32(len(b))
	for pos < max {
		chunkLen := byteOrder.Uint32(b[pos:])
		copy(b[pos:], w.filler)
		pos += chunkLen
	}
	return w.w.Write(b)
}
