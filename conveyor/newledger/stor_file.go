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

package newledger

// This structure is intended for incremental/streamed write with some abilities to detect corruptions and to do self heal.
// MUST: Strict order of field
type StorageFile struct {
	FormatVersion uint32 `protobuf:"varint,16,opt"` // != 0

	Head struct {
		Magic     uint32 `protobuf:"fixed32,16,opt"` // != 0, a number updated each time when a file is regenerated
		HeadMagic string `protobuf:"string,20,opt"`  // = "insolar-head"
		TailLen   uint32 `protobuf:"varint,21,opt"`  // must be defined at the creation of a file

		HeadObj interface{}

		SelfLen uint32 `protobuf:"fixed32,4095,opt"` // != 0, MUST be equal to byte len of this struct (we use fixed size here to make it easier to calculate)
	} `protobuf:"bytes,20,opt"` // required, and MUST be the second field in the file

	Content struct {
		Magic  uint32 `protobuf:"fixed32,16,opt"` // != 0, MUST match Head.Magic
		Serial uint32 `protobuf:"varint,20,opt"`  // serial/directory number of the entry, BUT file entries may go unordered

		DataObj interface{}

		SelfLen uint32 `protobuf:"fixed32,4095,opt"` // != 0, MUST be equal to byte len of this struct (we use fixed size here to make it easier to calculate)
	} `protobuf:"bytes,20<N<4095,rep"` // can be multiple entries of different types

	Tail struct {
		Magic     uint32 `protobuf:"fixed32,16,opt"` // != 0, MUST match Head.Magic
		TailMagic string `protobuf:"string,20,opt"`  // = "insolar-tail"

		TailObj interface{}

		Padding []byte `protobuf:"bytes,4094,opt"`   // as the size of Tail structure must be defined at the creation of a file, padding may be needed.
		SelfLen uint32 `protobuf:"fixed32,4095,opt"` // != 0, MUST be equal to byte len of this struct (we use fixed size here to make it easier to calculate)
	} `protobuf:"bytes,4095,opt"` // required, and MUST be the last field in the file
}
