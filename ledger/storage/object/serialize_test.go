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

package object

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
)

func Test_RecordByTypeIDPanic(t *testing.T) {
	assert.Panics(t, func() { RecordFromType(0) })
}

func TestSerializeDeserializeRecord(t *testing.T) {
	cs := platformpolicy.NewPlatformCryptographyScheme()

	rec := ObjectActivateRecord{
		ObjectStateRecord: ObjectStateRecord{
			Memory: CalculateIDForBlob(cs, core.GenesisPulse.PulseNumber, []byte{1, 2, 3}),
		},
	}
	serialized := SerializeRecord(&rec)
	deserialized := DeserializeRecord(serialized)
	assert.Equal(t, rec, *deserialized.(*ObjectActivateRecord))
}
