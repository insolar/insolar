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

package hash

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const aRecTypeID = 1

type aRec struct {
	A1 int32
	A2 string
}

func (rec *aRec) WriteHashData(w io.Writer) {
	var data = []interface{}{
		int16(aRecTypeID),
		rec.A1,
		[]byte(rec.A2),
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

const bRecTypeID = 2

type bRec struct {
	aRec
	B1 []byte
}

func (rec *bRec) WriteHashData(w io.Writer) {
	rec.aRec.WriteHashData(w)
	var data = []interface{}{
		int16(bRecTypeID),
		rec.B1,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

const cRecTypeID = 3

type cRec struct {
	bRec
	C1 uint64
}

func (rec *cRec) WriteHashData(w io.Writer) {
	rec.bRec.WriteHashData(w)
	var data = []interface{}{
		int16(cRecTypeID),
		rec.C1,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

func TestMain(m *testing.M) {
	rand.Seed(42)
	retCode := m.Run()
	os.Exit(retCode)
}

var sha3hash224tests = []struct {
	name       string
	hw         Writer
	expectHash string
}{
	{"aRecEmpty", &aRec{},
		"1836d23c35f95c2a63bd91a9c88ca4efe35c4765e276b43cbbf36cff"},
	{"bRecEmpty", &bRec{},
		"ce24330ce23d1bf4913693dfe4f573e4e952d95cac4577b7ecf7ca5c"},
	{"cRecEmpty", &cRec{},
		"8557377f3fadd8f94eda9f6c2c391e294ad201acf3df76e16842b9a9"},
	{"aRecNonEmpty1", &aRec{A1: 100500, A2: "100500"},
		"55863f02f4cb31ec5394d8dee80a2f0555c4b9b6c89d50f57d3624dc"},
	{"bRecNonEmpty1", &bRec{B1: []byte("100500")},
		"f974703c97ea174fbd257c87a50d03db8cbe5f8960a209d7615cb245"},
	{"cRecNonEmpty1", &cRec{C1: 100500},
		"aabaf2a1d32d81bde809694ea2da94645940e2d1cf2135d9b1f67171"},
}

func Test_emptyHashesNotTheSame(t *testing.T) {
	found := make(map[string]string)
	for _, tt := range sha3hash224tests {
		testname, ok := found[tt.expectHash]
		if !ok {
			found[tt.expectHash] = testname
			continue
		}
		t.Errorf("found %s hash for \"%s\" test, should not repeats in sha3hash224tests", tt.expectHash, tt.name)
	}
}

func Test_sha3hash224(t *testing.T) {
	for _, tt := range sha3hash224tests {
		t.Run(tt.name+"_hashtest", func(t *testing.T) {
			// s := fmtSprintf(tt.in, &flagprinter)
			hashBytes := SHA3hash224(tt.hw)
			hashHexStr := fmt.Sprintf("%x", hashBytes)

			assert.Equal(t, 28, len(hashBytes))
			assert.Equal(t, tt.expectHash, hashHexStr)
			// if s != tt.out {
			// t.Errorf("got %q, want %q", s, tt.out)
			// }
		})
	}
	// check is empty hashes is not the same
}

var sha3hash224testsNotTheSame = []struct {
	name string
	hw1  Writer
	hw2  Writer
}{
	{"aRec-A1", &aRec{}, &aRec{A1: 10}},
	{"aRec-A2", &aRec{}, &aRec{A2: "test"}},
	{"bRec-B1", &bRec{}, &bRec{B1: []byte("test")}},
	{"bRec-level1", &bRec{}, &bRec{aRec: aRec{A1: 100}}},
	{"cRec-level1", &cRec{}, &cRec{bRec: bRec{B1: []byte("test")}}},
	{"cRec-level2", &cRec{}, &cRec{bRec: bRec{aRec: aRec{A2: "hi"}}}},
	{"cRec-level2", &cRec{}, &cRec{bRec: bRec{aRec: aRec{A1: 5}}}},
}

// check if changed embed hashes is not the same
func Test_sha3hash224_IfChangedNotTheSame(t *testing.T) {
	for _, tt := range sha3hash224testsNotTheSame {
		t.Run(tt.name, func(t *testing.T) {
			// s := fmtSprintf(tt.in, &flagprinter)
			hashBytes1 := SHA3hash224(tt.hw1)
			hashBytes2 := SHA3hash224(tt.hw2)
			hashHexStr1 := fmt.Sprintf("%x", hashBytes1)
			hashHexStr2 := fmt.Sprintf("%x", hashBytes2)

			if hashHexStr1 == hashHexStr2 {
				t.Errorf("struct should not be with the same hash:\n%s", hashHexStr1)
			}
		})
	}
}
