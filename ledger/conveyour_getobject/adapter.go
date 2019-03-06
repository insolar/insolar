/*
 *    Copyright 2019 Insolar
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

package conveyour_getobject

import (
	"github.com/insolar/insolar/core"
)

type JetAdapterTask struct {
	Object core.RecordID
}

type JetAdapterResult struct {
	JetID core.RecordID
}

type JetAdapter struct {
}

func (a *JetAdapter) Process(task JetAdapterTask) JetAdapterResult {
	// Fetch jet from network...
	return JetAdapterResult{JetID: core.RecordID{}}
}

type GetObjectTask struct {
	Object core.RecordID
	JetID  core.RecordID
}

type GetObjectResult struct {
	memory []byte
}

type GetObjectAdapter struct {
}

func (a *GetObjectAdapter) Process(task GetObjectTask) GetObjectResult {
	// Fetch object from db...
	return GetObjectResult{memory: []byte{1, 2, 3}}
}
