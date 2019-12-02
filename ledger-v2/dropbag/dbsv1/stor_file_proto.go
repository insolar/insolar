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

// This structure is intended for incremental/streamed write with some abilities to detect corruptions and to do self heal.
// MUST: Strict order of field
type StorageFileScheme struct {
	FormatVersionAndOptions uint64 `protobuf:"fixed64,16,req"` // != 0

	Head struct {
		MagicAndCRC  uint64 `protobuf:"fixed64,16,req"` // != 0, a number updated each time when a file is regenerated
		HeadMagicStr string `protobuf:"string,17,req"`  // = "insolar-head"

		TailLen uint32 `protobuf:"varint,19,req"` // must be defined at creation of a file

		HeadObj interface{}

		SelfChk uint64 `protobuf:"fixed64,2047,req"` // != 0, MUST be equal to byte len of this struct (we use fixed size here to make it easier to calculate)
	} `protobuf:"bytes,20,opt"` // required, and MUST be the second field in the file

	Entry struct {
		MagicAndCRC uint64 `protobuf:"fixed64,16,req"` // != 0, MUST match Head.Magic+EntrySeqNo

		EntryOptions uint32 `protobuf:"varint,19,req"`

		DataObj interface{}

		SelfChk uint64 `protobuf:"fixed64,2047,opt"` // != 0, MUST be equal to byte len of this struct (we use fixed size here to make it easier to calculate)
	} `protobuf:"bytes,20<N<2046,rep"` // can be multiple entries of different types

	AlignPadding []byte `protobuf:"bytes,2046,opt"` // optional, and MUST be the 2nd from the end

	Tail struct {
		MagicAndCRC      uint64 `protobuf:"fixed64,16,req"` // != 0, MUST match Head.Magic
		TailMagicStr     string `protobuf:"string,17,req"`  // = "insolar-tail"
		EntryCountAndCRC uint64 `protobuf:"fixed64,18,req"` //
		Padding          []byte `protobuf:"bytes,19,opt"`   // as the size of Tail structure must be defined at the creation of a file, padding is needed.

		TailObj interface{}

		SelfChk uint64 `protobuf:"fixed64,2047,req"` // != 0, MUST be equal to byte len of this struct (we use fixed size here to make it easier to calculate)
	} `protobuf:"bytes,2047,opt"` // required, and MUST be the last field in the file
}
