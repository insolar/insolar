/*
 *    Copyright 2018 Insolar
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

package storage_test

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/insolar/insolar/core"
)

func zerohash() []byte {
	b := make([]byte, core.RecordHashSize)
	return b
}

func randhash() []byte {
	b := zerohash()
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func hexhash(hash string) []byte {
	b := zerohash()
	if len(hash)%2 == 1 {
		hash = "0" + hash
	}
	h, err := hex.DecodeString(hash)
	if err != nil {
		panic(err)
	}
	_ = copy(b, h)
	return b
}

func referenceWithHashes(domainhash, recordhash string) core.RecordRef {
	dh := hexhash(domainhash)
	rh := hexhash(recordhash)

	return *core.NewRecordRef(*core.NewRecordID(0, dh), *core.NewRecordID(0, rh))
}
