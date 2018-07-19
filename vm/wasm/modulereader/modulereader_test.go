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

package modulereader

import (
	"bytes"
	"testing"
)

func TestReadVarUint32(t *testing.T) {
	r := Reader{R: bytes.NewReader([]byte{0x80, 0x7f})}
	n, err := r.ReadVarUint32()
	if err != nil {
		t.Fatal(err)
	}
	if n != uint32(16256) {
		t.Fatalf("got = %d; want = %d", n, 16256)
	}
}

func TestReadVarint32(t *testing.T) {
	r := Reader{R: bytes.NewReader([]byte{0xFF, 0x7e})}
	n, err := r.ReadVarint32()
	if err != nil {
		t.Fatal(err)
	}
	if n != int32(-129) {
		t.Fatalf("got = %d; want = %d", n, -129)
	}
}
