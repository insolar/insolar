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

import "github.com/insolar/insolar/core"

func NewRecordIDFromRecord(scheme core.PlatformCryptographyScheme, pulse core.PulseNumber, rec Record) *core.RecordID {
	hasher := scheme.ReferenceHasher()
	_, err := rec.WriteHashData(hasher)
	if err != nil {
		panic(err)
	}
	return core.NewRecordID(pulse, hasher.Sum(nil))
}
