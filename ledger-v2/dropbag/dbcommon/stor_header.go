//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package dbcommon

import (
	"github.com/insolar/insolar/ledger-v2/protokit"
	"io"
	"math"
)

type PayloadFactory interface {
	CreatePayloadBuilder(format FileFormat, sr StorageSeqReader) PayloadBuilder
}

type ChapterDetails struct {
	// [1..)
	Index   int
	Options uint32
	// [0..2024]
	Type uint16
}

type StorageEntryPosition struct {
	StorageOffset int64
	ByteOffset    int
}

type PayloadBuilder interface {
	AddPrelude(bytes []byte, pos StorageEntryPosition) error
	AddConclude(bytes []byte, pos StorageEntryPosition, totalCount uint32) error
	AddChapter(bytes []byte, pos StorageEntryPosition, details ChapterDetails) error
	NeedsNextChapter() bool

	Finished() error
	Failed(error) error
}

type PayloadWriter interface {
	io.Closer
	WritePrelude(bytes []byte, concludeMaxLength int) (StorageEntryPosition, error)
	WriteConclude(bytes []byte) (StorageEntryPosition, error)
	WriteChapter(bytes []byte, details ChapterDetails) (StorageEntryPosition, error)
}

type FileFormat uint16
type FormatOptions uint64

const formatFieldId = 16

var formatField = protokit.WireFixed64.Tag(formatFieldId)

type ReadConfig struct {
	ReadAllEntries bool
	AlwaysCopy     bool
	//StorageOptions FormatOptions
}

func ReadFormatAndOptions(sr StorageSeqReader) (FileFormat, FormatOptions, error) {
	if v, err := formatField.DecodeFrom(sr); err != nil {
		return 0, 0, err
	} else {
		return FileFormat(v & math.MaxUint16), FormatOptions(v >> 16), nil
	}
}

func WriteFormatAndOptions(sw StorageSeqWriter, format FileFormat, options FormatOptions) error {
	return formatField.EncodeTo(sw, uint64(format)|uint64(options)<<16)
}

type StorageSeqReader interface {
	io.ByteReader
	io.Reader
	io.Seeker
	CanSeek() bool
	CanReadMapped() bool
	Offset() int64
	ReadMapped(n int64) ([]byte, error)
}

type StorageBlockReader interface {
	io.ReaderAt
	io.Closer
	ReadAtMapped(n int64, off int64) ([]byte, error)
}

type StorageSeqWriter interface {
	io.ByteWriter
	io.Writer
	Offset() int64
}
