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
	"math/rand"
	"os"
	"testing"
)

type aRec struct {
	A1 int32
	A2 string
}

type bRec struct {
	aRec
	B1 []byte
}

type cRec struct {
	bRec
	C1 uint64
}

func TestMain(m *testing.M) {
	rand.Seed(42)
	retCode := m.Run()
	os.Exit(retCode)
}

var sha3hash224tests = []struct {
	name       string
	hw         interface{}
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
