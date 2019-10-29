//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package appfoundation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar/gen"
)

func TestIsEthereumAddress(t *testing.T) {
	tests := []struct {
		name     string
		addr     string
		expected bool
	}{
		{"good address 1", "0x35567Abc4Fa54fe30d200F76A4868A70383e7938", true},
		{"good address 2", "0xA9BfF538A906154c80A8dBccd229F3DEddFa52D6", true},
		{"good address 3", "87a0edA943f79C31a21f123e2946726c5Dbd1F75", true},
		{"bad address 1", "39m5Wvn9ZqyhYmCYpsyHuGMt5YYw4Vmh1ZddFa52", false},
		{"bad address 2", gen.Reference().String(), false},
		{"short address", "0xA9BfF538A906154c80A8dBccd229F3DEddFa52", false},
		{"long address", "0x35567Abc4Fa54fe30d200F76A4868A70383e7938c8", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsEthereumAddress(test.addr))
		})
	}
}
