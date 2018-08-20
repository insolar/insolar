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

package jetdrop

import (
	"golang.org/x/crypto/sha3"
)

// JetDrop is a blockchain block. It contains hashes from all records from slot.
type JetDrop struct {
	PrevHash     []byte
	RecordHashes [][]byte // TODO: this should be a byte slice that represents the merkle tree root of records
}

// Hash calculates jet drop hash. Raw data for hash should contain previous hash and merkle tree hash from records.
func (jd *JetDrop) Hash() ([]byte, error) {
	encoded, err := EncodeJetDrop(jd)
	if err != nil {
		return nil, err
	}
	h := sha3.New224()
	_, err = h.Write(encoded)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
