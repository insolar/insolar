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

package ledger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/hash"
)

func TestGenRequestRecordID(t *testing.T) {
	tt := map[string]hash.Writer{
		"constructor": &message.CallConstructor{},
		"call":        &message.CallMethod{},
	}
	for name := range tt {
		t.Run(name, func(t *testing.T) {
			h := tt[name]
			res := GenRequestRecordID(0, h)
			fmt.Printf("%T %+v => %x\n", h, h, res)
			assert.NotNil(t, res)
		})
	}
}
