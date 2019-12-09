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

package dbsv1

import (
	"bytes"
	"fmt"
	"hash"
	"hash/crc32"
	"io"

	"github.com/insolar/insolar/ledger-v2/dropbag/dbcommon"
	"github.com/insolar/insolar/ledger-v2/protokit"
)

type StorageFileV1Writer struct {
	StorageFileV1
	w           dbcommon.StorageSeqWriter
	totalCrc    hash.Hash32
	nextChapter entryNo
	tailLength  uint64
}

func (p *StorageFileV1Writer) Close() error {
	if c, ok := p.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func initMagic() uint32 {
	return 1
	//t := uint64(0)
	//for {
	//	t += uint64(time.Now().UnixNano())
	//	t ^= t >> 1
	//	v := uint32(t) ^ uint32(t>>32)
	//	if v != 0 {
	//		return v
	//	}
	//}
}

func NewStorageFileV1Writer(sw dbcommon.StorageSeqWriter, sf StorageFileV1) (*StorageFileV1Writer, error) {
	if sw == nil {
		panic("illegal value")
	}
	if err := sf.prepare(); err != nil {
		return nil, err
	}
	if sf.MagicZero == 0 {
		sf.MagicZero = initMagic()
	}
	return &StorageFileV1Writer{StorageFileV1: sf, w: sw}, nil
}

func (p *StorageFileV1Writer) WritePrelude(b []byte, tailMaxPayloadLength int) (dbcommon.StorageEntryPosition, error) {
	switch {
	case p.w == nil:
		panic("illegal state")
	case p.nextChapter != headEntry:
		panic("illegal state")
	case len(b) > MaxHeadPayloadLength:
		panic("illegal value - MaxHeadPayloadLength exceeded")
	case tailMaxPayloadLength < 0:
		panic("illegal value")
	case tailMaxPayloadLength > MaxTailPayloadLength:
		panic("illegal value - MaxTailPayloadLength exceeded")
	}

	p.tailLength = uint64(tailMaxPayloadLength + minTailLength)
	p.totalCrc = crc32.New(p.CrcTable)

	return p.writeField(headField, func(w *bytes.Buffer, _ int) error {
		if err := magicStringField.EncodeTo(w, uint64(len(magicStrHead))); err != nil {
			return err
		}
		if _, err := w.Write(([]byte)(magicStrHead)); err != nil {
			return err
		}
		if err := declaredTailLenField.EncodeTo(w, p.tailLength); err != nil {
			return err
		}
		p.nextChapter++
		return nil
	}, b, maxHeadLength-MaxHeadPayloadLength, true, maxHeadLength)
}

func (p *StorageFileV1Writer) WriteConclude(b []byte) (dbcommon.StorageEntryPosition, error) {
	switch {
	case p.nextChapter <= headEntry:
		panic("illegal state")
	case uint64(len(b)+minTailLength) > p.tailLength:
		panic("illegal value - conclude length exceeds declared limit")
	}
	totalCount := uint64(p.nextChapter) - 1
	p.nextChapter = tailEntry

	// write alignment padding
	if paddingBeforeTailLength := int(p.w.Offset() % int64(p.TailAlign)); paddingBeforeTailLength != 0 {

		paddingBeforeTailLength += minPaddingBeforeTail

		if paddingBeforeTailLength > int(p.TailAlign) {
			paddingBeforeTailLength %= int(p.TailAlign)
		}
		paddingBeforeTailLength = int(p.TailAlign) - paddingBeforeTailLength
		paddingBeforeTailLength -= protokit.SizeVarint32(uint32(paddingBeforeTailLength)) - 1

		switch {
		case paddingBeforeTailLength > maxPaddingBeforeTail:
			panic("unexpected - overflow before-tail padding")
		case paddingBeforeTailLength < 0:
			panic("unexpected - negative before-tail padding")
		}

		if err := p.writePadding(p.w, paddingField, uint64(paddingBeforeTailLength)); err != nil {
			return dbcommon.StorageEntryPosition{}, err
		}
		if checkAlignment := p.w.Offset() % int64(p.TailAlign); checkAlignment != 0 {
			panic("unexpected - invalid before-tail padding")
		}
	}

	return p.writeField(tailField, func(w *bytes.Buffer, knownLength int) error {
		if err := magicStringField.EncodeTo(w, uint64(len(magicStrTail))); err != nil {
			return err
		}
		if _, err := w.Write(([]byte)(magicStrTail)); err != nil {
			return err
		}
		if err := totalCountAndCrcField.EncodeTo(w, totalCount|uint64(p.totalCrc.Sum32())<<32); err != nil {
			return err
		}

		// TODO Padding with one field is unable to handle some cases, e.g. to pad 128 bytes, as it will cause 2 byte increment due to varint
		paddingLength := int(p.tailLength) - w.Len() - knownLength - int(tailInnerPadding.FieldSize(0))
		if paddingLength < 0 {
			panic("unexpected - invalid inner-tail padding")
		}
		if err := p.writePadding(w, tailInnerPadding, uint64(paddingLength)); err != nil {
			return err
		}

		p.nextChapter++
		return nil
	}, b, int(p.tailLength)-len(b), true, maxTailLength)
}

func (p *StorageFileV1Writer) WriteChapter(b []byte, details dbcommon.ChapterDetails) (dbcommon.StorageEntryPosition, error) {
	switch {
	case p.nextChapter <= headEntry:
		panic("illegal state")
	case len(b) > MaxChapterPayloadLength:
		panic("illegal value - MaxChapterPayloadLength exceeded")
	case details.Type > maxChapterFieldId-minChapterFieldId:
		panic("illegal value - ChapterDetails.Type limit exceeded")
	case details.Index != int(p.nextChapter):
		panic("illegal value - ChapterDetails.Index")
	}

	chapterField := protokit.WireBytes.Tag(int(details.Type) + minChapterFieldId)
	chapterOptions := uint64(details.Options)
	hasSelfLen := p.StorageOptions&ChapterWithoutSelfCheckOption != 0

	return p.writeField(chapterField, func(w *bytes.Buffer, _ int) error {
		if err := chapterOptionsField.EncodeTo(w, chapterOptions); err != nil {
			return err
		}
		p.nextChapter++
		return nil
	}, b, int(chapterOptionsField.FieldSize(chapterOptions)), hasSelfLen, maxChapterLength)
}

func (p *StorageFileV1Writer) writeField(fieldTag protokit.WireTag, prefixFn func(w *bytes.Buffer, knownLength int) error,
	b []byte, preBufSize int, hasSelfCheck bool, maxFieldSize int) (dbcommon.StorageEntryPosition, error) {

	entryNo := p.nextChapter

	preBuf := bytes.Buffer{}
	if preBufSize > 0 {
		preBuf.Grow(preBufSize)
	}

	entryLength := magicCRCFieldLength + len(b)
	if hasSelfCheck {
		entryLength += selfChkFieldLength
	}

	dataStartPos := int64(0)

	// avoids repeated: return dbcommon.StorageEntryPosition{}, err
	if err := func() error {
		if err := prefixFn(&preBuf, entryLength); err != nil {
			return err
		}
		entryLength += preBuf.Len()

		if entryLength > maxFieldSize {
			return fmt.Errorf("internal length limit error: entry=%v length=%d limit=%d", entryNo, entryLength, maxFieldSize)
		}

		if err := fieldTag.EncodeTo(p.w, uint64(entryLength)); err != nil {
			return err
		}

		dataStartPos = p.w.Offset()

		postBuf := bytes.Buffer{}
		if hasSelfCheck {
			switch {
			case entryLength > maxSelfCheckLength:
				panic("unexpected - file entry is too big")
			case dataStartPos > maxSelfCheckEntryPos:
				panic("unexpected - file is too big")
			}

			postBuf.Grow(selfChkFieldLength)
			if err := selfChkField.EncodeTo(&postBuf, uint64(entryLength)|uint64(dataStartPos)<<bitsSelfCheckLength); err != nil {
				return err
			}
			if postBuf.Len() != selfChkFieldLength {
				panic("unexpected, mismatched selfChkFieldLength")
			}
		}

		crcValue := calcCrc32(calcCrc32(
			crc32.New(p.CrcTable), preBuf.Bytes()), b,
		).Sum32()

		magic := p.MagicZero
		switch {
		case entryNo > 0:
			magic += uint32(entryNo)
			fallthrough
		case entryNo == headEntry:
			addCrc32(p.totalCrc, crcValue)
		}

		if err := magicCrcField.EncodeTo(p.w, uint64(magic)|uint64(crcValue)<<32); err != nil {
			return err
		}

		if _, err := p.w.Write(preBuf.Bytes()); err != nil {
			return err
		}
		if _, err := p.w.Write(b); err != nil {
			return err
		}
		if _, err := p.w.Write(postBuf.Bytes()); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		return dbcommon.StorageEntryPosition{}, err
	}

	dataEndPos := p.w.Offset()
	if dataStartPos+int64(entryLength) != dataEndPos {
		panic(fmt.Errorf("internal error, entry length mismatched: entry=%v start=%d length=%d, end=%d", entryNo, dataStartPos, entryLength, dataEndPos))
	}

	return dbcommon.StorageEntryPosition{dataStartPos, preBuf.Len()}, nil
}

type byteWriter interface {
	io.ByteWriter
	io.Writer
}

func (p *StorageFileV1Writer) writePadding(w byteWriter, fieldTag protokit.WireTag, paddingLength uint64) error {
	switch err := fieldTag.EncodeTo(w, paddingLength); {
	case err != nil:
		return err
	case paddingLength == 0:
		return nil
	case paddingLength < skipPortion:
		_, err = w.Write(make([]byte, paddingLength))
		return err
	}

	padBuf := make([]byte, skipPortion)
	for {
		switch n, err := w.Write(padBuf); {
		case err != nil:
			return err
		case paddingLength == uint64(n):
			return nil
		default:
			paddingLength -= uint64(n)

			if paddingLength < skipPortion {
				_, err = w.Write(padBuf[:paddingLength])
				return err
			}
		}
	}
}
