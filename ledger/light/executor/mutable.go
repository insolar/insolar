// Copyright 2020 Insolar Network Ltd.
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

package executor

import (
	"github.com/insolar/insolar/insolar/record"
)

// OldestMutable searches for a oldest mutable request through a provided list of open requests
// openRequests MUST be time-ascending order
func OldestMutable(openRequests []record.CompositeFilamentRecord) *record.CompositeFilamentRecord {
	isMutableIncoming := func(rec record.CompositeFilamentRecord) bool {
		req := record.Unwrap(&rec.Record.Virtual).(record.Request)
		inReq, isIn := req.(*record.IncomingRequest)
		return isIn && !inReq.Immutable
	}

	if len(openRequests) == 0 {
		return nil
	}

	for i := 0; i < len(openRequests); i++ {
		if isMutableIncoming(openRequests[i]) {
			return &openRequests[i]
		}
	}

	return nil
}
