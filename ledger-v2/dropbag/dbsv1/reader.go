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

package dbsv1

import (
	"bytes"
	"fmt"
	"github.com/insolar/insolar/ledger-v2/dropbag/dbcommon"
	"hash/crc32"
	"io"
	"math"

	"github.com/insolar/insolar/ledger-v2/protokit"
)

const (
	ChapterWithoutSelfCheckOption dbcommon.FormatOptions = 1 << iota

	maskV1Options = (1 << iota) - 1
)

const magicCRCFieldLength = storageTagSize + 8

var magicCrcField = protokit.WireFixed64.Tag(16).EnsureFixedFieldSize(magicCRCFieldLength)

const selfChkFieldLength = storageTagSize + 8

var selfChkField = protokit.WireFixed64.Tag(2047).EnsureFixedFieldSize(selfChkFieldLength)

const (
	bitsSelfCheckLength  = 8 * 3 // 3 bytes
	selfCheckLengthMask  = 1<<bitsSelfCheckLength - 1
	maxSelfCheckLength   = selfCheckLengthMask
	maxSelfCheckEntryPos = math.MaxUint64 >> bitsSelfCheckLength

	storageTagSize     = 2
	minVarintFieldSize = storageTagSize + 1

	magicStrHead  = "insolar-head"
	minHeadLength = magicCRCFieldLength + // StorageFileScheme.Head.MagicAndCRC
		minVarintFieldSize + len(magicStrHead) + // StorageFileScheme.Head.HeadMagicStr
		minVarintFieldSize + // StorageFileScheme.Head.TailLen
		selfChkFieldLength // StorageFileScheme.Head.SelfChk

	maxHeadLength        = minHeadLength + 1<<17             // MUST be less than maxSelfCheckLength
	MaxHeadPayloadLength = maxHeadLength - minHeadLength - 3 // SizeVarint32(maxTailLength) - 1

	magicStrTail  = "insolar-tail"
	minTailLength = magicCRCFieldLength + // StorageFileScheme.Tail.MagicAndCRC
		minVarintFieldSize + len(magicStrTail) + // StorageFileScheme.Tail.TailMagicStr
		totalCountAndCrcFieldLength + // StorageFileScheme.Tail.EntryCountAndCRC
		minVarintFieldSize + // StorageFileScheme.Tail.Padding
		selfChkFieldLength // StorageFileScheme.Tail.SelfChk

	maxTailLength        = maxSelfCheckLength
	MaxTailPayloadLength = maxTailLength - minTailLength

	minChapterLength = minVarintFieldSize + // StorageFileScheme.Entry.EntryOptions
		magicCRCFieldLength // StorageFileScheme.Entry.MagicAndCRC

	maxChapterLength        = maxSelfCheckLength
	MaxChapterPayloadLength = maxChapterLength - minChapterLength - (protokit.MaxVarintSize - 1) // for EntryOptions

	maxTailOffsetAlignment = 1 << 17 // must be less than MaxInt32
	minTailOffsetAlignment = 1 << 6

	minPaddingBeforeTail = minVarintFieldSize
	maxPaddingBeforeTail = maxTailOffsetAlignment
)

const (
	headFieldId = 20

	minChapterFieldId = 21
	maxChapterFieldId = 2045

	paddingFieldId = 2046
	tailFieldId    = 2047
)

//maxChapterFieldId

var (
	headField        = protokit.WireBytes.Tag(headFieldId)
	paddingField     = protokit.WireBytes.Tag(paddingFieldId)
	tailField        = protokit.WireBytes.Tag(tailFieldId)
	magicStringField = protokit.WireBytes.Tag(17) // head and tail only
)

const totalCountAndCrcFieldLength = storageTagSize + 8

var totalCountAndCrcField = protokit.WireFixed64.Tag(18).EnsureFixedFieldSize(totalCountAndCrcFieldLength) // tail only

var tailInnerPadding = protokit.WireVarint.Tag(19)     // tail only
var declaredTailLenField = protokit.WireVarint.Tag(19) // head only
var chapterOptionsField = protokit.WireVarint.Tag(19)  // entry only

type entryNo int

const (
	headEntry entryNo = 0
	tailEntry entryNo = -1
)

func (v entryNo) String() string {
	switch {
	case v > 0:
	case v == headEntry:
		return "head"
	case v == tailEntry:
		return "tail"
	}
	return fmt.Sprintf("%d", v)
}

type StorageFileV1 struct {
	StorageOptions dbcommon.FormatOptions
	MagicZero      uint32
	TailAlign      uint32
	CrcTable       *crc32.Table
}

func (p *StorageFileV1) CheckOptions() error {
	if p.StorageOptions&^maskV1Options != 0 {
		return fmt.Errorf("unsupported storage options: format=1 options=%x", p.StorageOptions)
	}
	switch {
	case p.TailAlign == 0:
	case p.TailAlign > maxTailOffsetAlignment || p.TailAlign < minTailOffsetAlignment:
		return fmt.Errorf("unsupported tail alignment: format=1 align=%x", p.TailAlign)
	}
	return nil
}

func (p *StorageFileV1) prepare() error {
	if err := p.CheckOptions(); err != nil {
		return err
	}
	if p.CrcTable == nil {
		p.CrcTable = crc32.MakeTable(crc32.Castagnoli)
	}
	if p.TailAlign == 0 {
		p.TailAlign = 1
	}
	return nil
}

type StorageFileV1Reader struct {
	StorageFileV1
	Config  dbcommon.ReadConfig
	Builder dbcommon.PayloadBuilder
}

func (p *StorageFileV1Reader) Read(sr dbcommon.StorageSeqReader) error {
	if p.Builder == nil {
		panic("illegal state")
	}
	if err := p.prepare(); err != nil {
		return err
	}

	var expectedTailLength uint64
	totalCrc := crc32.New(p.CrcTable)

	if headLength, err := headField.DecodeFrom(sr); err != nil {
		return err
	} else if err := p.readStorageEntry(sr, headEntry, headLength, minHeadLength-selfChkFieldLength, maxHeadLength, true,
		func(ofs int64, b []byte, magicSeq, crc uint32) error {
			if magicSeq == 0 {
				return fmt.Errorf("illegal content, invalid magic: entry=%v actual=%x", headEntry, magicSeq)
			}
			addCrc32(totalCrc, crc)
			p.MagicZero = magicSeq
			etl, err := p.readPrelude(b, ofs)
			expectedTailLength = etl
			return err
		},
	); err != nil {
		return err
	}

	canSeek := sr.CanSeek()
	hasSelfLen := p.StorageOptions&ChapterWithoutSelfCheckOption != 0
	tailLength := uint64(0)
	entryCount := uint64(0)
	partialCount := true
	lastFieldStart := int64(0)

outer:
	for entryNo := headEntry + 1; ; entryNo++ {

		lastFieldStart = sr.Offset()
		wt, entryLength, err := protokit.WireBytes.DecodeFrom(sr)
		if err != nil {
			return err
		}

		entryId := wt.FieldId()
		switch {
		case entryId >= minChapterFieldId && entryId <= maxChapterFieldId:
			//
		case entryId == tailFieldId:
			tailLength = entryLength
			partialCount = false
			break outer
		case entryId == paddingFieldId:
			if err = p.skipPadding(sr, entryLength, maxPaddingBeforeTail); err != nil {
				return err
			}

			lastFieldStart = sr.Offset()
			if tailLength, err = tailField.DecodeFrom(sr); err != nil {
				return err
			}
			partialCount = false
			break outer
		default:
			return fmt.Errorf("illegal entry content: entry=%v entryId=%d", entryNo, entryId)
		}

		if canSeek && !p.Builder.NeedsNextChapter() && !p.Config.ReadAllEntries {
			adjustedTailLength := int64(expectedTailLength)
			adjustedTailLength += storageTagSize
			adjustedTailLength += int64(protokit.SizeVarint64(expectedTailLength))

			lastFieldStart, err = sr.Seek(-adjustedTailLength, io.SeekEnd)
			if err != nil {
				return err
			}
			if tailLength, err = tailField.DecodeFrom(sr); err != nil {
				return err
			}
			break outer
		}

		if err := p.readStorageEntry(sr, entryNo, entryLength, minChapterLength, maxChapterLength, hasSelfLen,
			func(ofs int64, b []byte, magicSeq, crc uint32) error {
				if p.MagicZero+uint32(entryNo) != magicSeq {
					return fmt.Errorf("corrupted content, magic mismatch: entry=%v expected=%x actual=%x",
						entryNo, p.MagicZero+uint32(entryNo), magicSeq)
				}
				addCrc32(totalCrc, crc)
				return p.readEntry(b, ofs, entryNo, entryId)
			},
		); err != nil {
			return err
		}
		entryCount++
	}

	entryNo := tailEntry
	switch tailLength {
	case 0:
		return fmt.Errorf("illegal content, empty tail tag: entry=%v length=0", entryNo)
	case expectedTailLength:
		//
	default:
		return fmt.Errorf("corrupted content, length mismatch: entry=%v expected=%d actual=%d",
			entryNo, expectedTailLength, tailLength)
	}

	if lastFieldStart%int64(p.TailAlign) != 0 {
		return fmt.Errorf("illegal content, tail is unaligned: position=%d", lastFieldStart)
	}

	if err := p.readStorageEntry(sr, entryNo, tailLength, minTailLength-selfChkFieldLength, maxTailLength, true,
		func(ofs int64, b []byte, magicSeq, _ uint32) error {
			if p.MagicZero != magicSeq {
				return fmt.Errorf("corrupted content, magic mismatch: entry=%v expected=%x actual=%x",
					entryNo, p.MagicZero, magicSeq)
			}
			totalCrcSum := uint32(0)
			if !partialCount {
				totalCrcSum = totalCrc.Sum32()
			}
			err := p.readConclude(b, ofs, entryCount, partialCount, totalCrcSum)
			return err
		},
	); err != nil {
		return err
	}

	if canSeek {
		lastPos := sr.Offset()
		switch seekEnd, err := sr.Seek(0, io.SeekEnd); {
		case err != nil:
			return err
		case lastPos != seekEnd:
			return fmt.Errorf("illegal content, data beyond eof: expectedLen=%x actualLen=%x", lastPos, seekEnd)
		}
	}

	return nil
}

func (p *StorageFileV1Reader) readStorageEntry(sr dbcommon.StorageSeqReader, entryNo entryNo, delimLen uint64, minLen, maxLen int,
	hasSelfCheck bool, processFn func(ofs int64, b []byte, magicSeq, crc uint32) error,
) error {
	if hasSelfCheck {
		minLen += selfChkFieldLength
	}

	if delimLen < uint64(minLen) || delimLen > uint64(maxLen) {
		return fmt.Errorf("illegal content: entry=%v length=%d", entryNo, delimLen)
	}

	//if err := sr.Prefetch(int64(delimLen) + storageTagSize + 16); err != nil {
	//	return err
	//}

	rawStart := sr.Offset()

	magicAndCrc, err := magicCrcField.DecodeFrom(sr)
	if err != nil {
		return err
	}
	crcValue := uint32(magicAndCrc >> 32)

	dataStart := sr.Offset()
	dataLen := int64(delimLen) - (dataStart - rawStart)
	var buf []byte

	switch {
	case dataLen == 0:
	case !p.Config.AlwaysCopy && sr.CanReadMapped():
		switch buf, err = sr.ReadMapped(dataLen); {
		case err != nil:
			return err
		case int64(len(buf)) != dataLen:
			return fmt.Errorf("inconsistent map-read: entry=%v expected=%d actual=%d", entryNo, dataLen, len(buf))
		}
	default:
		buf = make([]byte, dataLen)
		if _, err := io.ReadFull(sr, buf); err != nil {
			return err
		}
	}

	if hasSelfCheck {
		fieldOfs := len(buf) - selfChkFieldLength
		if fieldOfs < 0 {
			return fmt.Errorf("corrupted content, insufficient length: entry=%v outer=%d deficit=%d", entryNo, delimLen, fieldOfs)
		}
		rBuf := bytes.NewBuffer(buf[fieldOfs:])
		switch selfChk, err := selfChkField.DecodeFrom(rBuf); {
		case err != nil:
			return err
		case rBuf.Len() != 0:
			return fmt.Errorf("internal error, unaligned read: entry=%v field=selfLen", entryNo)
		case delimLen != selfChk&selfCheckLengthMask:
			return fmt.Errorf("corrupted content, self check failed: entry=%v outer=%d selfLength=%d", entryNo, delimLen, selfChk&selfCheckLengthMask)
		case uint64(rawStart) != selfChk>>bitsSelfCheckLength:
			return fmt.Errorf("corrupted content, self check failed: entry=%v actual=%d selfPos=%d", entryNo, rawStart, selfChk>>bitsSelfCheckLength)
		}
		buf = buf[:fieldOfs]
	}

	if crcCoder := calcCrc32(crc32.New(p.CrcTable), buf); crcValue != crcCoder.Sum32() {
		return fmt.Errorf("corrupted content, crc mismatch: entry=%v expected=%x actual=%x", entryNo, crcValue, crcCoder.Sum32())
	}

	return processFn(dataStart, buf, uint32(magicAndCrc), crcValue)
}

func (p *StorageFileV1Reader) readPrelude(buf []byte, ofs int64) (uint64, error) {
	rBuf := bytes.NewBuffer(buf)

	switch strLen, err := magicStringField.DecodeFrom(rBuf); {
	case err != nil:
		return 0, err
	case uint64(len(magicStrHead)) != strLen:
		return 0, fmt.Errorf("illegal content: entry=%v len(headMagic)=%d expected=%d", headEntry, strLen, len(magicStrHead))
	default:
		strBuf := rBuf.Next(len(magicStrHead))
		if magicStrHead != string(strBuf) {
			return 0, fmt.Errorf("illegal content: entry=%v headMagic='%s'", headEntry, strBuf)
		}
	}

	if declaredTailLen, err := declaredTailLenField.DecodeFrom(rBuf); err != nil {
		return 0, err
	} else if declaredTailLen < uint64(minTailLength) || declaredTailLen > uint64(maxTailLength) {
		return 0, fmt.Errorf("illegal content: entry=%v declaredTailLen=%d", headEntry, declaredTailLen)
	} else {
		err = p.Builder.AddPrelude(buf, dbcommon.StorageEntryPosition{ofs, len(buf) - rBuf.Len()})
		return declaredTailLen, err
	}
}

func (p *StorageFileV1Reader) readEntry(buf []byte, ofs int64, entryNo entryNo, entryId int) error {

	rBuf := bytes.NewBuffer(buf)

	entryOptions, err := chapterOptionsField.DecodeFrom(rBuf)
	if err != nil {
		return err
	}

	return p.Builder.AddChapter(buf,
		dbcommon.StorageEntryPosition{ofs, len(buf) - rBuf.Len()},
		dbcommon.ChapterDetails{int(entryNo), uint32(entryOptions), uint16(entryId - minChapterFieldId)})
}

func (p *StorageFileV1Reader) readConclude(buf []byte, ofs int64, totalCount uint64, partialCount bool, totalCrc uint32) error {
	rBuf := bytes.NewBuffer(buf)

	switch strLen, err := magicStringField.DecodeFrom(rBuf); {
	case err != nil:
		return err
	case uint64(len(magicStrTail)) != strLen:
		return fmt.Errorf("illegal content: entry=%v len(tailMagic)=%d expected=%d", tailEntry, strLen, len(magicStrTail))
	default:
		strBuf := rBuf.Next(len(magicStrTail))
		if magicStrTail != string(strBuf) {
			return fmt.Errorf("illegal content: entry=%v tailMagic='%s'", tailEntry, strBuf)
		}
	}

	switch totals, err := totalCountAndCrcField.DecodeFrom(rBuf); {
	case err != nil:
		return err
	case totalCount > math.MaxUint32:
		return fmt.Errorf("illegal content: entry=%v totalCount=%d", tailEntry, totalCount)
	case partialCount:
		if uint32(totalCount) >= uint32(totals) {
			return fmt.Errorf("illegal content: entry=%v partialCount=%d expected=%d", tailEntry, totalCount, uint32(totals))
		}
		totalCount = totals & math.MaxUint32
	case uint32(totalCount) != uint32(totals):
		return fmt.Errorf("illegal content: entry=%v totalCount=%d expected=%d", tailEntry, totalCount, uint32(totals))
	case totalCrc != uint32(totals>>32):
		return fmt.Errorf("corrupted content, crc mismatch: entry=%v expected=%x actual=%x", tailEntry, uint32(totals>>32), totalCrc)
	}

	switch paddingLength, err := tailInnerPadding.DecodeFrom(rBuf); {
	case err != nil:
		return err
	case paddingLength >= maxTailLength:
		return fmt.Errorf("illegal content: entry=%v padding=%d", tailEntry, paddingLength)
	default:
		skippedLen := len(rBuf.Next(int(paddingLength)))
		if skippedLen != int(paddingLength) {
			return fmt.Errorf("illegal content: entry=%v padding=%d actual=%d", tailEntry, paddingLength, skippedLen)
		}
	}

	return p.Builder.AddConclude(buf, dbcommon.StorageEntryPosition{ofs, len(buf) - rBuf.Len()},
		uint32(totalCount))
}

const skipPortion = 4096

func (p *StorageFileV1Reader) skipPadding(sr dbcommon.StorageSeqReader, paddingLength uint64, maxLength int) error {
	switch {
	case paddingLength == 0:
		return nil
	case paddingLength > uint64(maxLength):
		return fmt.Errorf("illegal content: entry=padding length=%d", paddingLength)
	case sr.CanSeek():
		_, err := sr.Seek(int64(paddingLength), io.SeekCurrent)
		return err
	case paddingLength <= skipPortion:
		_, err := io.ReadFull(sr, make([]byte, paddingLength))
		return err
	}

	skipBuf := make([]byte, skipPortion)
	for {
		switch n, err := io.ReadFull(sr, skipBuf); {
		case err != nil:
			return err
		case paddingLength == uint64(n):
			return nil
		default:
			paddingLength -= uint64(n)

			if paddingLength < skipPortion {
				_, err = io.ReadFull(sr, skipBuf[:paddingLength])
				return err
			}
		}
	}
}
