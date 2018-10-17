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

package record

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// This ensures serialized reference has Record prefix and Domain suffix.
// It's required for selecting records by record pulse
func TestReference_Key(t *testing.T) {
	ref := Reference{
		Domain: ID{Pulse: 1},
		Record: ID{Pulse: 2},
	}
	assert.Equal(t, []byte{0, 0, 0, 2}, ref.CoreRef()[:4])
}
