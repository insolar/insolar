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

package index

import (
	"github.com/insolar/insolar/ledger/record"
)

// ClassLifeline represents meta information for record object
type ClassLifeline struct {
	LatestStateID record.ID   // Amend or activate record
	AmendIDs      []record.ID // ClassAmendRecord
}

// ObjectLifeline represents meta information for record object
type ObjectLifeline struct {
	ClassID       record.ID
	LatestStateID record.ID   // Amend or activate record
	StateIDs      []record.ID // ObjectAppendRecord or ObjectAmendRecord
	AppendIDs     []record.ID // ObjectAppendRecord
}
