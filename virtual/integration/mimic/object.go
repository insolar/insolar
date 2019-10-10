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

package mimic

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
)

func recordStateToStateID(state record.State) record.StateID {
	switch state.(type) {
	case *record.Activate:
		return record.StateActivation
	case *record.Amend:
		return record.StateAmend
	case *record.Deactivate:
		return record.StateDeactivation
	default:
		return record.StateUndefined
	}
}

type ObjectState struct {
	State     record.State
	RequestID insolar.ID
}

type ObjectEntity struct {
	ObjectChanges []ObjectState
	RequestsMap   map[insolar.ID]*RequestEntity
	RequestsList  []*RequestEntity
}

func (e *ObjectEntity) isDeactivated() bool {
	return e.getLatestStateID() == record.StateDeactivation
}

func (e *ObjectEntity) addIncomingRequest(entity *RequestEntity) {
	e.RequestsMap[entity.ID] = entity
	e.RequestsList = append(e.RequestsList, entity)
}

func (e ObjectEntity) getLatestStateID() record.StateID {
	if len(e.ObjectChanges) == 0 {
		return record.StateUndefined
	}

	return recordStateToStateID(e.ObjectChanges[len(e.ObjectChanges)-1].State)
}

// getRequestInfo returns:
// * count of opened requests
// * earliest request
// * latest request
func (e ObjectEntity) getRequestsInfo() (uint32, *RequestEntity, *RequestEntity) {
	var (
		openRequestCount uint32
		firstRequest     *RequestEntity
		lastRequest      *RequestEntity
	)

	for _, request := range e.RequestsList {
		if request.Status == RequestRegistered {
			openRequestCount++
			if firstRequest == nil {
				firstRequest = request
			}
			lastRequest = request
		}
	}
	return openRequestCount, firstRequest, lastRequest
}

type CodeEntity struct {
	Code []byte
}
