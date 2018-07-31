/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package record

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

// Raw struct contains raw serialized record.
// We need raw blob to not have dependency on record structure changes in future,
// and have ability of consistent hash checking on old records.
type Raw struct {
	Type TypeID
	Data []byte
}

// Hash returns 28 bytes of SHA3 hash on Data field.
func (s *Raw) Hash() Hash {
	return sha3.Sum224(s.Data)
}

// Key2ID converts Key with PulseNum and Hash pair to binary representation (record.ID).
func Key2ID(k Key) ID {
	var id ID
	var err error
	buf := bytes.NewBuffer(id[:0])

	err = binary.Write(buf, binary.BigEndian, k.Pulse)
	if err != nil {
		panic("binary.Write failed to write PulseNum:" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, k.Hash)
	if err != nil {
		panic("binary.Write failed to write Hash:" + err.Error())
	}
	return id
}
