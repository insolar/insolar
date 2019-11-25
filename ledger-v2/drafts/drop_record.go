///
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
///

package drafts

import "github.com/insolar/insolar/longbits"

type DropEntry struct { // byteSize=324
	RegistrantSign longbits.Bits512 // byteSize=64,  sign of DropEntryRegistration

	RegistryRecord DropEntryRegistration // byteSize=132
	RecordBody     DropRecordBody        // byteSize=128
	RecordPayload  DropRecordPayload
}

type DropEntryRegistration struct { // byteSize=132
	EntryTypeAndFlags DropEntryType

	EntryLocator DropEntryLocator // PrimarySection indicator
	EntryHash    longbits.Bits224 // hash of DropRecordBody

	PredecessorRef ShortRef

	ProducerSign longbits.Bits512 // sign of EntryHash + pulse + ?
}

type ShortRef longbits.Bits256
type FullRef [2]ShortRef

type DropRecordBody struct { // byteSize=~256, minProtoSize=2*(3 + 28) + 2*(3 + 32) + (2 + 8) = 142
	RecordTypeAndFlags  DropEntryType // DelegationFlag, PayloadIsHash
	PayloadTypeAndFlags uint32

	ProducerNodeHash   longbits.Bits224 // hash part of NodeRef
	RegistrantNodeHash longbits.Bits224 // hash part of NodeRef - optional, depends on Flags

	LifelineRef []ShortRef // 0, 1 or 2 - depends on Flags

	//	PayloadHash        longbits.Bits256 // optional - depends on Flags
	SmallPayload []byte //0-256-512, hash-or-content - depends on Flags
}

type DropEntryType uint32
type DropEntryLocator uint32

type DropRecordPayload struct {
	DelegationTokens []byte

	CustomPayload []byte
}

type CallRequestRecord struct {
	PayloadTypeAndFlags uint32
	// SmallPayload
	// FullPayload
	OutgoingReq DropEntry // compacted version that doesn't include calculable fields
}

type CustomCryptographySectionEntry struct {
	_ map[int] /* offset of a field in the original binary */ CustomCryptographyItem
}

type CustomCryptographyItem []byte
