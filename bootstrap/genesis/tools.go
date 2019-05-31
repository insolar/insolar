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

package genesis

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
)

func refByName(name string) insolar.Reference {
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	vRec := record.Wrap(record.Request{
		CallType: record.CTGenesis,
		Method:   name,
	})
	id := insolar.NewID(insolar.FirstPulseNumber, record.HashVirtual(pcs.ReferenceHasher(), vRec))
	return *insolar.NewReference(*id)
}
