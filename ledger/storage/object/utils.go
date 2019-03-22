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

import "github.com/insolar/insolar/insolar"

func NewRecordIDFromRecord(scheme insolar.PlatformCryptographyScheme, pulse insolar.PulseNumber, rec Record) *insolar.RecordID {
	hasher := scheme.ReferenceHasher()
	_, err := rec.WriteHashData(hasher)
	if err != nil {
		panic(err)
	}
	return insolar.NewRecordID(pulse, hasher.Sum(nil))
}
