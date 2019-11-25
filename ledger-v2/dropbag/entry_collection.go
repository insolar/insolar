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

import "github.com/insolar/insolar/longbits"

type FieldId uint16
type FieldOffset int

const (
	InvalidField FieldOffset = 0 - iota
	AbsentField
	SearchField
)

type FieldOffsetMap interface {
	FindField(int FieldId) FieldOffset
}

type StorageLocId uint32

type EntryPosition struct {
	StorageLoc StorageLocId
	SequenceId uint32
	StorageOfs int64
}

func (v EntryPosition) IsZero() bool {
	return v.StorageOfs == 0
}

func (v EntryPosition) NotFound() bool {
	return v.IsZero()
}

type KeyEntryCollection struct {
	entries  map[longbits.ByteString]keyEntry
	complete bool
}

type keyEntry struct {
	EntryPosition
}

func (p *KeyEntryCollection) FindKey(k longbits.ByteString) (bool, EntryPosition) {
	return p.complete, p.entries[k].EntryPosition
}

//type CompleteEntryCollection struct {
//	storageLoc StorageLocId
//	entries    []StructEntry
//}
//
//func (p CompleteEntryCollection) GetEntryBySeq(seq uint32) (bool, StructEntry) {
//	if seq >= uint32(len(p.entries)) {
//		return true, StructEntry{}
//	}
//	return true, p.entries[seq]
//}
//
//func (p CompleteEntryCollection) GetEntryByPos(pos EntryPosition) (bool, StructEntry) {
//	switch {
//	case pos.NotFound():
//		panic("illegal value")
//	case pos.StorageLoc != p.storageLoc:
//		panic("illegal value")
//	}
//	return p.GetEntryBySeq(pos.SequenceId)
//}
//
type LazyEntryCollection struct {
	storageLoc  StorageLocId
	entries     []StructEntry
	loadInfo    []lazyEntry
	totalCount  uint32
	loadedCount uint32
}

type StructEntry struct {
	content  []byte
	fieldMap FieldOffsetMap
}

func (v StructEntry) NotFound() bool {
	return v.content == nil && v.fieldMap == nil
}

type lazyEntry struct {
}

func (p *LazyEntryCollection) isComplete() bool {
	return p.totalCount == p.loadedCount
}

func (p *LazyEntryCollection) GetEntryBySeq(seq uint32) (bool, StructEntry) {
	if seq >= uint32(len(p.entries)) {
		return p.isComplete(), StructEntry{}
	}
	return p.isComplete(), p.entries[seq]
}

func (p *LazyEntryCollection) GetEntryByPos(pos EntryPosition) (bool, StructEntry) {
	switch {
	case pos.NotFound():
		panic("illegal value")
	case pos.StorageLoc != p.storageLoc:
		panic("illegal value")
	}
	return p.GetEntryBySeq(pos.SequenceId)
}

func (p *LazyEntryCollection) LoadEntries(posList []EntryPosition, loadFn func(EntryPosition) StructEntry) {

	maxSeq := uint32(0)
	for _, pos := range posList {
		switch {
		case pos.NotFound():
			panic("illegal value")
		case pos.StorageLoc != p.storageLoc:
			panic("illegal value")
		case pos.SequenceId >= p.totalCount:
			panic("illegal value")
		case pos.SequenceId > maxSeq:
			maxSeq = pos.SequenceId
		}
	}

	if p.isComplete() {
		return
	}

	switch n := uint32(len(p.entries)); {
	case n < maxSeq:
	case n == maxSeq:
		p.entries = append(p.entries, StructEntry{})
		p.loadInfo = append(p.loadInfo, lazyEntry{})
	default:
		n = maxSeq - n + 1
		p.entries = append(p.entries, make([]StructEntry, n)...)
		p.loadInfo = append(p.loadInfo, make([]lazyEntry, n)...)
	}

	for _, pos := range posList {
		le := &p.entries[pos.SequenceId]
		if !le.NotFound() {
			continue
		}
		le.StructEntry = loadFn(pos)
	}

	if p.totalCount == p.loadedCount {
		p.loadInfo = nil
	}
}

func (p LazyEntryCollection) loadEntry(pos EntryPosition, loadFn func(EntryPosition) StructEntry) {

}
