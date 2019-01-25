/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
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
