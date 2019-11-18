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
	"fmt"
	"github.com/insolar/insolar/ledger-v2/protokit"
	"io"
	"math"
)

type StorageReader interface {
	io.ByteReader
	io.Reader
	io.Seeker
	CanSeek() bool
	Offset() int64
	ReadMapped(size int64) ([]byte, error)
	ReadCopy(size int64) ([]byte, error)
	//Prefetch(size int64) error
}

const (
	v1optionEntryWithoutSelfLen FormatOptions = 1 << iota

	maskV1Options = (1 << iota) - 1
)

const (
	storageTagSize = 2

	selfLenFieldLength = storageTagSize + 4

	magicHead     = "insolar-head"
	minHeadLength = storageTagSize + 4 + // StorageFileScheme.Head.Magic
		storageTagSize + 1 + len(magicHead) + // StorageFileScheme.Head.HeadMagic
		storageTagSize + 1 + // StorageFileScheme.Head.TailLen
		selfLenFieldLength // StorageFileScheme.Head.SelfLen

	maxHeadLength = minHeadLength + 1<<16

	magicTail     = "insolar-tail"
	minTailLength = storageTagSize + 4 + // StorageFileScheme.Tail.Magic
		storageTagSize + 1 + len(magicTail) + // StorageFileScheme.Tail.TailMagic
		selfLenFieldLength // StorageFileScheme.Tail.SelfLen

	maxTailLength = minTailLength + 1<<24

	minPreambleLength = storageTagSize + protokit.MaxVarintSize + // StorageFileScheme.FormatVersion
		storageTagSize + protokit.MaxVarintSize + // tag of StorageFileScheme.Head
		minHeadLength

	minEntryLength = storageTagSize + 4 + // StorageFileScheme.Content.Magic
		selfLenFieldLength // StorageFileScheme.Head.SelfLen

	maxEntryLength = minEntryLength + 1<<24

	tailAlignment        = 4096
	maxPaddingBeforeTail = tailAlignment - 1 - minTailLength
	maxPaddingInsideTail = math.MaxUint32
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
	formatField = protokit.WireFixed64.Tag(formatFieldId)
	headField   = protokit.WireBytes.Tag(headFieldId)
	//paddingField = protokit.WireBytes.Tag(paddingFieldId)
	tailField = protokit.WireBytes.Tag(tailFieldId)

	entryMagicField   = protokit.WireFixed32.Tag(16)
	entryOptionsField = protokit.WireVarint.Tag(19)

	magicStringField    = protokit.WireBytes.Tag(17)  // for head and tail only
	headTailLengthField = protokit.WireVarint.Tag(18) // for head only

	tailPaddingField  = protokit.WireBytes.Tag(2046)
	entrySelfLenField = protokit.WireFixed32.Tag(2047)
)

type entryNo int

const (
	headEntry    entryNo = 0
	tailEntry    entryNo = -1
	paddingEntry entryNo = -2
)

func (v entryNo) String() string {
	if v < 0 {
		switch v {
		case headEntry:
			return "head"
		case tailEntry:
			return "tail"
		case paddingEntry:
			return "padding"
		}
	}
	return fmt.Sprintf("%d", v)
}

type StorageFileV1 struct {
	Config     ReadConfig
	Builder    PayloadBuilder
	EntryMagic uint32

	//Head []byte
	//Tail []byte
	//Body [][]byte
}

func (p *StorageFileV1) Read(sr StorageReader) error {
	if p.Builder == nil {
		panic("illegal state")
	}
	if p.Config.StorageOptions&^maskV1Options != 0 {
		return fmt.Errorf("unsupported storage options: format=1 options=%x", p.Config.StorageOptions)
	}

	var expectedTailLength int64

	if headLength, err := headField.DecodeFrom(sr); err != nil {
		return err
	} else if headLength < uint64(minHeadLength) || headLength > uint64(maxHeadLength) {
		return fmt.Errorf("unsupported storage format: headLength=%d", headLength)
	} else {
		// read head
		headStart := sr.Offset()

		//if err := sr.Prefetch(int64(headLength)); err != nil {
		//	return err
		//}

		if entryMagic, err := entryMagicField.DecodeFrom(sr); err != nil {
			return err
		} else if entryMagic > math.MaxUint32 {
			return fmt.Errorf("illegal internal state: entry=head magic=%x", entryMagic)
		} else {
			p.EntryMagic = uint32(entryMagic)
		}

		if headMagicLength, err := magicStringField.DecodeFrom(sr); err != nil {
			return err
		} else if headMagicLength != uint64(len(magicHead)) {
			return fmt.Errorf("unsupported storage format: entry=head magicStrLength=%d", headMagicLength)
		}
		if err := protokit.MatchString(magicHead, sr, ": entry=head"); err != nil {
			return err
		}

		if tailLength, err := headTailLengthField.DecodeFrom(sr); err != nil {
			return err
		} else if tailLength < uint64(minTailLength) || tailLength > uint64(maxTailLength) {
			return fmt.Errorf("unsupported storage format: entry=head tailLength=%d", tailLength)
		} else {
			expectedTailLength = int64(tailLength)
		}

		if entryBytes, err := p.readCustomPayload(sr, headEntry, headStart, headLength); err != nil {
			return err
		} else if err := p.Builder.AddPreamble(entryBytes, headStart, int64(headLength)); err != nil {
			return err
		}
	}

	remainingLength := uint64(1)

	if !p.Config.NonLazyEntries && sr.CanSeek() {
		if _, err := sr.Seek(-expectedTailLength, io.SeekEnd); err != nil {
			return err
		}
		remainingLength = 0
	}

	if tailLength, err := p.readAllEntries(sr, remainingLength); err != nil {
		return err
	} else if tailBytes, tailStart, err := p.readTail(sr, tailLength); err != nil {
		return err
	} else {
		return p.Builder.AddPreamble(tailBytes, tailStart, int64(tailLength))
	}
}

func (p *StorageFileV1) readAllEntries(sr StorageReader, remainingLength uint64) (uint64, error) {
	entryNo := headEntry
	for {
		entryNo++
		if wt, entryLength, err := protokit.WireBytes.DecodeFrom(sr); err != nil {
			return 0, err
		} else {
			entryId := wt.FieldId()
			switch {
			case entryId >= minEntryFieldId && entryId <= maxEntryFieldId:
				if entryLength < uint64(minEntryLength) || entryLength > uint64(maxEntryLength) || entryLength > remainingLength {
					return 0, fmt.Errorf("unsupported storage format: entry=%d length=%d remainingLength=%d", entryNo, entryLength, remainingLength)
				}
			case entryId == paddingFieldId:
				if err := p.skipPadding(sr, entryLength, uint64(maxPaddingBeforeTail)); err != nil {
					return 0, err
				}
				// ensures that Tail comes immediately after the padding
				if entryLength, err = tailField.DecodeFrom(sr); err != nil {
					return 0, err
				}
				return entryLength, nil
			case entryId == tailFieldId:
				return entryLength, nil
			default:
				return 0, fmt.Errorf("unsupported storage format: entry=%d entryNo=%d", entryNo, entryId)
			}

			entryStart := sr.Offset()

			if entryMagic, err := entryMagicField.DecodeFrom(sr); err != nil {
				return 0, err
			} else if entryMagic > math.MaxUint32 || uint32(entryMagic) != p.EntryMagic {
				return 0, fmt.Errorf("illegal internal state: entry=%v magic=%x expected=%x", entryNo, entryMagic, p.EntryMagic)
			}

			if customPayload, err := p.readCustomPayload(sr, entryNo, entryStart, entryLength); err != nil {
				return 0, err
			} else if err := p.Builder.AddEntry(int(entryNo), customPayload, entryStart, int64(entryLength)); err != nil {
				return 0, err
			}
		}
	}
}

func (p *StorageFileV1) readTail(sr StorageReader, tailLength uint64) ([]byte, int64, error) {

	if tailLength < uint64(minTailLength) || tailLength > uint64(maxTailLength) {
		return nil, 0, fmt.Errorf("unsupported storage format: entry=tail length=%d", tailLength)
	}

	tailStart := sr.Offset()

	if err := func() error {
		if entryMagic, err := entryMagicField.DecodeFrom(sr); err != nil {
			return err
		} else if entryMagic > math.MaxUint32 || uint32(entryMagic) != p.EntryMagic {
			return fmt.Errorf("illegal internal state: entry=tail magic=%x expected=%x", entryMagic, p.EntryMagic)
		}

		if tailMagicLength, err := magicStringField.DecodeFrom(sr); err != nil {
			return err
		} else if tailMagicLength != uint64(len(magicTail)) {
			return fmt.Errorf("unsupported storage format: entry=tail magicStrLength=%d", tailMagicLength)
		}
		return protokit.MatchString(magicTail, sr, ": entry=tail")
	}(); err != nil {
		return nil, 0, err
	}

	if customPayload, err := p.readCustomPayload(sr, tailEntry, tailStart, tailLength); err != nil {
		return nil, 0, err
	} else {
		return customPayload, tailStart, nil
	}
}

func (p *StorageFileV1) readCustomPayload(sr StorageReader, entryNo entryNo, entryStart int64, entryLength uint64) (entryBytes []byte, err error) {

	var entryOptions uint64
	if entryOptions, err = entryOptionsField.DecodeFrom(sr); err != nil {
		return nil, err
	}

	customEntryStart := sr.Offset()
	customEntryLength := int64(entryLength) - selfLenFieldLength
	customEntryLength -= customEntryStart - entryStart

	switch {
	case customEntryLength > 0:
	case customEntryLength == 0:
		if entryNo <= 0 {
			return nil, nil
		}
		fallthrough
	default:
		return nil, fmt.Errorf("illegal internal state: entry=%v customLength=%d", entryNo, customEntryLength)
	}

	if p.Config.AlwaysCopy {
		entryBytes, err = sr.ReadCopy(customEntryLength)
	} else {
		entryBytes, err = sr.ReadMapped(customEntryLength)
	}

	switch {
	case err != nil:
		return nil, err
	case entryNo == tailEntry:
		if paddingLen, err := tailPaddingField.DecodeFrom(sr); err != nil {
			return nil, err
		} else if err := p.skipPadding(sr, paddingLen, uint64(maxPaddingInsideTail)); err != nil {
			return nil, err
		}
	case entryNo <= 0:
		panic("illegal value")
	case p.Config.StorageOptions&v1optionEntryWithoutSelfLen != 0:
		return entryBytes, nil
	}

	if selfLen, err := entrySelfLenField.DecodeFrom(sr); err != nil {
		return nil, err
	} else if selfLen != entryLength {
		return nil, fmt.Errorf("integrity failure: entry=%v selfLen=%d expected=%d", entryNo, selfLen, entryLength)
	}
	return entryBytes, nil
}

const skipPortion = 4096

func (p *StorageFileV1) skipPadding(sr StorageReader, paddingLength, maxLength uint64) error {
	switch {
	case paddingLength == 0:
		return nil
	case paddingLength > maxLength:
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
