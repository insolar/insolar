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

package nodenetwork

import (
	"encoding/binary"
	"fmt"
	"hash"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
)

func hashWriteChecked(hash hash.Hash, data []byte) {
	n, err := hash.Write(data)
	if n != len(data) {
		panic(fmt.Sprintf("Error writing hash. Bytes expected: %d; bytes actual: %d", len(data), n))
	}
	if err != nil {
		panic(err)
	}
}

func calculateNodeHash(scheme core.PlatformCryptographyScheme, processor core.KeyProcessor, node core.Node) []byte {
	h := scheme.IntegrityHasher()
	hashWriteChecked(h, node.ID().Bytes())

	b := [8]byte{}
	binary.LittleEndian.PutUint32(b[:4], uint32(node.ShortID()))
	hashWriteChecked(h, b[:4])
	binary.LittleEndian.PutUint32(b[:4], uint32(node.Role()))

	hashWriteChecked(h, b[:4])
	pk, err := processor.ExportPublicKeyPEM(node.PublicKey())
	if err != nil {
		panic(err)
	}
	hashWriteChecked(h, pk)
	hashWriteChecked(h, []byte(node.PhysicalAddress()))
	hashWriteChecked(h, []byte(node.Version()))
	return h.Sum(nil)
}

// CalculateHash calculates hash of active node list
func CalculateHash(scheme core.PlatformCryptographyScheme, list []core.Node) (result []byte, err error) {
	// catch possible panic from hashWriteChecked in this function and in all calculateNodeHash funcs
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, fmt.Errorf("error calculating h: %s", r)
		}
	}()

	h := scheme.IntegrityHasher()
	processor := platformpolicy.NewKeyProcessor()
	for _, node := range list {
		nodeHash := calculateNodeHash(scheme, processor, node)
		hashWriteChecked(h, nodeHash)
	}
	return h.Sum(nil), nil
}
