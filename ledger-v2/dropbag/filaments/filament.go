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

package filaments

import (
	"github.com/insolar/insolar/ledger-v2/dropbag"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/reference"
)

type WriteEntry struct {
	Key reference.Holder

	// local or remote
	// remote may have Prev = nil if prev is not read
	Prev   *WriteEntry
	Next   AtomicEntry
	Latest *LocalSegment // FilamentSegment

	// EntryHash == Key.Local

	EventSeq     uint32 // per JetDrop
	DirectorySeq uint32 // per JetDrop, index in Merkle tree
	FilamentSeq  uint64 // sequence through the whole filament, first entry =1, zero is invalid

	// StorageLocator uint64(?) // atomic

	BodyHash    cryptkit.Digest
	BodySection dropbag.JetSectionId
	Body        *WriteEntryBody

	// AuthRecordHash etc - an additional section for custom primary cryptography

	ProducerSignature  cryptkit.Signature // over BodyRecordHash + (?)AuthRecordHash
	RegistrarSignature cryptkit.Signature // over ProducerSignature + Sequences + Prev.Key.Local (EntryHash)
}

type WriteEntryBody struct {
	//	RecordHash  cryptkit.Digest

	ProducerNode  reference.Holder // TODO make specific type for Node ref
	RegistrarNode reference.Holder // TODO make specific type for Node ref

	LifelineRoot reference.Holder
	// OtherRef
	// ReasonRef

	// other sections
}

func (p *WriteEntry) LifelineRoot() reference.Holder {
	if ll := p.Latest.lifelineRoot; ll != nil {
		return ll
	}
	return p.FilamentRoot()
}

func (p *WriteEntry) FilamentRoot() reference.Holder {
	return p.Latest.filamentRoot
}

func (p *WriteEntry) FilamentSection() dropbag.JetSectionId {
	return p.Latest.filamentSection
}
