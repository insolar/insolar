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

import "io"

type CompositeDropStorage interface {
	// identified by a first pulse

	// Pulses & PulseData
	// Jets
	// Cabinet -> StorageCabinet
}

//type JetPrefix uint32
type ShortJetId uint32 // JetPrefix + 5bit length
type FullJetId uint64  // ShortJetId + LastSplitPulse

type DropStorage interface {
	// FullJetId
	// PulseNumber

	// indication of special properties of a drop: - active summary, archive summary, archived drop

	// === one for all jets === hash ref to StorageCabinet?
	// PulseData
	// JetTree
	// NodeList

	// FindByKey - accross sections?
	// DropSection
}

type DropSection interface {
	// Record listing
}

type DropLifeline interface {
	// consists of 1xDropOpening, 1xDropClosing, 1+ DropRevisions
	// DropRevision keeps info on rearrangements

	// TODO appending summary info updates - needs something simple and cheap
}

// ==================

type EntryStorageCabinet interface {
	// Has ControlSection that keeps DropLifelines etc
	// Jet trees
	// Node lists

}

type EntryStorageShelf interface {
	// one per section type per EntryStorageCabinet
	// Consists of EntryCollections:
	// - directory (keys + brief info) partitioned by drops // can have one storage file per partition when is written, can be combined later
	// - alt_directory (keys + brief info but by using an alternative cryptography scheme)
	// - content (record data) // one per shelf, accessed by index+ofs+len

	// hides differences for read-only collections and collections being written
}

type EntryCollection interface {
	// index-based access

	// hides differences:
	// - lazy and non-lazy read implementation of an indexed set
	// - set being built
}

type EntryStorageAdapter interface {
	// provides support lazy / packed / open-read-close access to physical storage
}

// ================== file specific implementation

type StorageFileFolder interface {
}

type StorageFile interface {
}

type StorageFileReader interface {
}

type StorageFormatAdapter interface {
	// checks individual entry CRC
	// checks file on reopening
	// facilitates read of lazy entries -> need to know format
}

type StorageURI string
type StorageReadAdapter interface {
	GetURI() StorageURI

	OpenForSeqRead() StorageSeqReader
	OpenForBlockRead() StorageBlockReader
}

type StorageSeqReader interface {
	io.ByteReader
	io.Reader
	io.Closer
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
