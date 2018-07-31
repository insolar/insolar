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
	"testing"

	"github.com/stretchr/testify/assert"
)

var convertTests = []struct {
	name string
	key  Key
	id   ID
}{
	{
		key: Key{Pulse: 10, Hash: str2Bytes("21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193")},
		id:  str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
	},
}

func Test_KeyIDConversion(t *testing.T) {
	for _, tt := range convertTests {
		t.Run(tt.name, func(t *testing.T) {
			gotID := Key2ID(tt.key)
			gotKey := ID2Key(gotID)
			assert.Equal(t, tt.key, gotKey)
			assert.Equal(t, tt.id, gotID)
		})
	}
}
