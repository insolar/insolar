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

package dropbag

import (
	"bytes"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"math"

	"github.com/insolar/insolar/ledger-v2/protokit"
)

const (
	v1optionEntryWithoutSelfLen FormatOptions = 1 << iota

	maskV1Options = (1 << iota) - 1
)

const magicCRCFieldLength = storageTagSize + 8

var magicCrcField = protokit.WireFixed64.Tag(16).EnsureFixedFieldSize(magicCRCFieldLength)

const selfLenFieldLength = storageTagSize + 4

var selfLenField = protokit.WireFixed32.Tag(2047).EnsureFixedFieldSize(selfLenFieldLength)

const (
	storageTagSize = 2
	minVarintSize  = storageTagSize + 1

	magicStrHead  = "insolar-head"
	minHeadLength = magicCRCFieldLength + // StorageFileScheme.Head.MagicAndCRC
		minVarintSize + len(magicStrHead) + // StorageFileScheme.Head.HeadMagicStr
		minVarintSize + // StorageFileScheme.Head.TailLen
		selfLenFieldLength // StorageFileScheme.Head.SelfLen

	maxHeadLength = minHeadLength + 1<<17

	magicStrTail  = "insolar-tail"
	minTailLength = magicCRCFieldLength + // StorageFileScheme.Tail.MagicAndCRC
		minVarintSize + len(magicStrTail) + // StorageFileScheme.Tail.TailMagicStr
		totalCountAndCrcFieldLength + // StorageFileScheme.Tail.EntryCountAndCRC
		selfLenFieldLength // StorageFileScheme.Tail.SelfLen

	maxTailLength = minTailLength + 1<<24

	minEntryLength = magicCRCFieldLength // StorageFileScheme.Content.MagicAndCRC

	maxEntryLength = minEntryLength + selfLenFieldLength + 1<<24

	TailOffsetAlignment  = 4096
	maxPaddingBeforeTail = TailOffsetAlignment - 1 - minTailLength
)

const (
	formatFieldId = 16
	headFieldId   = 20

	minEntryFieldId = 21
	maxEntryFieldId = 2045

	paddingFieldId = 2046
	tailFieldId    = 2047
)

var (
	formatField      = protokit.WireFixed64.Tag(formatFieldId)
	headField        = protokit.WireBytes.Tag(headFieldId)
	tailField        = protokit.WireBytes.Tag(tailFieldId)
	magicStringField = protokit.WireBytes.Tag(17) // head and tail only

)

const totalCountAndCrcFieldLength = storageTagSize + 8

var totalCountAndCrcField = protokit.WireFixed64.Tag(19).EnsureFixedFieldSize(totalCountAndCrcFieldLength) // tail only

var declaredTailLenField = protokit.WireVarint.Tag(19) // head only

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
	Config    ReadConfig
	Builder   PayloadBuilder
	MagicZero uint32
	CrcTable  *crc32.Table
}

func (p *StorageFileV1) Read(sr StorageSeqReader) error {
	if p.Builder == nil {
		panic("illegal state")
	}
	if p.Config.StorageOptions&^maskV1Options != 0 {
		return fmt.Errorf("unsupported storage options: format=1 options=%x", p.Config.StorageOptions)
	}

	var expectedTailLength uint64
	totalCrc := crc32.New(p.CrcTable)

	if headLength, err := headField.DecodeFrom(sr); err != nil {
		return err
	} else if err := p.readStorageEntry(sr, headEntry, headLength, minHeadLength, maxHeadLength, true,
		func(ofs int64, b []byte, magicSeq, crc uint32) error {
			if magicSeq == 0 {
				return fmt.Errorf("illegal content, invalid magic: entry=%v actual=%x", headEntry, magicSeq)
			}
			AddCrc32(totalCrc, crc)
			p.MagicZero = magicSeq
			etl, err := p.readPreamble(b, ofs)
			expectedTailLength = etl
			return err
		},
	); err != nil {
		return err
	}

	canSeek := sr.CanSeek()
	hasSelfLen := p.Config.StorageOptions&v1optionEntryWithoutSelfLen != 0
	tailLength := uint64(0)
	entryCount := uint64(0)
	partialCount := true

outer:
	for entryNo := headEntry + 1; ; entryNo++ {

		wt, entryLength, err := protokit.WireBytes.DecodeFrom(sr)
		if err != nil {
			return err
		}

		entryId := wt.FieldId()
		switch {
		case entryId >= minEntryFieldId && entryId <= maxEntryFieldId:
			//
		case entryId == tailFieldId:
			tailLength = entryLength
			partialCount = false
			break outer
		case entryId == paddingFieldId:
			if err = p.skipPadding(sr, entryLength, maxPaddingBeforeTail); err != nil {
				return err
			}
			if tailLength, err = tailField.DecodeFrom(sr); err != nil {
				return err
			}
			partialCount = false
			break outer
		default:
			return fmt.Errorf("illegal entry content: entry=%v entryId=%d", entryNo, entryId)
		}

		if canSeek && !p.Builder.NeedsNextEntry() && !p.Config.ReadAllEntries {
			if _, err := sr.Seek(-int64(expectedTailLength), io.SeekEnd); err != nil {
				return err
			}
			break outer
		}

		if err := p.readStorageEntry(sr, entryNo, entryLength, minEntryLength, maxEntryLength, hasSelfLen,
			func(ofs int64, b []byte, magicSeq, crc uint32) error {
				if p.MagicZero+uint32(entryNo) != magicSeq {
					return fmt.Errorf("corrupted content, magic mismatch: entry=%v expected=%x actual=%x",
						entryNo, p.MagicZero+uint32(entryNo), magicSeq)
				}
				AddCrc32(totalCrc, crc)
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

	if tailOffset := sr.Offset(); tailOffset%TailOffsetAlignment != 0 {
		return fmt.Errorf("illegal content, tail is unaligned")
	}

	if err := p.readStorageEntry(sr, entryNo, tailLength, minTailLength, maxTailLength, true,
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

func AddCrc32(hash hash.Hash32, x uint32) {
	// byte order is according to crc32.appendUint32
	if n, err := hash.Write([]byte{byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)}); err != nil || n != 4 {
		panic(fmt.Errorf("crc calc failure: %d, %v", n, err))
	}
}

func (p *StorageFileV1) readStorageEntry(sr StorageSeqReader, entryNo entryNo, delimLen uint64, minLen, maxLen int,
	hasSelfLen bool, processFn func(ofs int64, b []byte, magicSeq, crc uint32) error,
) error {
	if hasSelfLen {
		minLen += selfLenFieldLength
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

	if hasSelfLen {
		fieldOfs := len(buf) - selfLenFieldLength
		if fieldOfs < 0 {
			return fmt.Errorf("corrupted content, insufficient length: entry=%v outer=%d deficit=%d", entryNo, delimLen, fieldOfs)
		}
		rBuf := bytes.NewBuffer(buf[fieldOfs:])
		switch selfLen, err := selfLenField.DecodeFrom(rBuf); {
		case err != nil:
			return err
		case delimLen != selfLen:
			return fmt.Errorf("corrupted content, length mismatch: entry=%v outer=%d self=%d", entryNo, delimLen, selfLen)
		case rBuf.Len() != 0:
			return fmt.Errorf("internal error, unaligned read: entry=%v field=selfLen", entryNo)
		}
		buf = buf[:fieldOfs]
	}

	crcCoder := crc32.New(p.CrcTable)
	switch n, err := crcCoder.Write(buf); {
	case n != len(buf) || err != nil:
		panic(fmt.Errorf("crc calc failure: entry=%v n=%d len=%d", entryNo, n, len(buf)))
	case crcValue != crcCoder.Sum32():
		return fmt.Errorf("corrupted content, crc mismatch: entry=%v expected=%x actual=%x", entryNo, crcValue, crcCoder.Sum32())
	}

	return processFn(dataStart, buf, uint32(magicAndCrc), crcValue)
}

func (p *StorageFileV1) readPreamble(buf []byte, ofs int64) (uint64, error) {
	rBuf := bytes.NewBuffer(buf)

	switch strLen, err := magicStringField.DecodeFrom(rBuf); {
	case err != nil:
		return 0, err
	case uint64(len(magicStrHead)) != strLen:
		return 0, fmt.Errorf("illegal content: entry=%v len(headMagic)=%d expected=%d", headEntry, strLen, len(magicStrHead))
	default:
		strBuf := rBuf.Bytes()[:len(magicStrHead)]
		if magicStrHead != string(strBuf) {
			return 0, fmt.Errorf("illegal content: entry=%v headMagic='%s'", headEntry, strBuf)
		}
	}

	if declaredTailLen, err := declaredTailLenField.DecodeFrom(rBuf); err != nil {
		return 0, err
	} else if declaredTailLen < uint64(minTailLength) || declaredTailLen > uint64(maxTailLength) {
		return 0, fmt.Errorf("illegal content: entry=%v declaredTailLen=%d", headEntry, declaredTailLen)
	} else {
		err = p.Builder.AddPreamble(buf, ofs, len(buf)-rBuf.Len())
		return declaredTailLen, err
	}
}

func (p *StorageFileV1) readEntry(buf []byte, ofs int64, entryNo entryNo, entryId int) error {
	return p.Builder.AddEntry(int(entryNo), entryId, buf, ofs, 0)
}

func (p *StorageFileV1) readConclude(buf []byte, ofs int64, totalCount uint64, partialCount bool, totalCrc uint32) error {
	rBuf := bytes.NewBuffer(buf)

	switch strLen, err := magicStringField.DecodeFrom(rBuf); {
	case err != nil:
		return err
	case uint64(len(magicStrTail)) != strLen:
		return fmt.Errorf("illegal content: entry=%v len(tailMagic)=%d expected=%d", tailEntry, strLen, len(magicStrTail))
	default:
		strBuf := rBuf.Bytes()[:len(magicStrTail)]
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
	case uint32(totalCount) != uint32(totals):
		return fmt.Errorf("illegal content: entry=%v totalCount=%d expected=%d", tailEntry, totalCount, uint32(totals))
	case totalCrc != uint32(totals>>32):
		return fmt.Errorf("corrupted content, crc mismatch: entry=%v expected=%x actual=%x", tailEntry, uint32(totals>>32), totalCrc)
	}

	return p.Builder.AddConclude(buf, ofs, len(buf)-rBuf.Len())
}

const skipPortion = 4096

func (p *StorageFileV1) skipPadding(sr StorageSeqReader, paddingLength uint64, maxLength int) error {
	switch {
	case paddingLength == 0:
		return nil
	case paddingLength > uint64(maxLength):
		return fmt.Errorf("unsupported storage format: entry=padding length=%d", paddingLength)
	case sr.CanSeek():
		_, err := sr.Seek(int64(paddingLength), io.SeekCurrent)
		return err
	case paddingLength <= skipPortion:
		_, err := io.ReadFull(sr, make([]byte, paddingLength))
		return err
	}

	skipBuf := make([]byte, skipPortion)
	for {
		switch n, err := io.ReadAtLeast(sr, skipBuf, 1); {
		case err != nil:
			return err
		case paddingLength == uint64(n):
			return nil
		default:
			paddingLength -= uint64(n)
		}
	}
}
