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

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaseBindReplayBasics(t *testing.T) {
	t.Parallel()

	t.Run("blank", func(t *testing.T) {
		t.Parallel()

		replay := &CaseBindReplay{}
		assert.NotNil(t, replay)

		record, step := replay.NextStep()
		assert.Nil(t, record)
		assert.Equal(t, 0, step)

		record, step = replay.NextStep()
		assert.Nil(t, record)
		assert.Equal(t, 0, step)
	})

	t.Run("one request, no records", func(t *testing.T) {
		t.Parallel()

		replay := &CaseBindReplay{
			CaseBind: CaseBind{
				Requests: []CaseRequest{
					CaseRequest{
						Request: CaseRecord{},
						Records: []CaseRecord{},
					},
				},
			},
			Request: 0,
			Record:  -1,
		}
		assert.NotNil(t, replay)

		record, step := replay.NextStep()
		assert.NotNil(t, record)
		assert.Equal(t, 1, step)

		record, step = replay.NextStep()
		assert.Nil(t, record)
		assert.Equal(t, 1, step)
	})

	t.Run("one request, one record", func(t *testing.T) {
		t.Parallel()

		replay := &CaseBindReplay{
			CaseBind: CaseBind{
				Requests: []CaseRequest{
					CaseRequest{
						Request: CaseRecord{Type: CaseRecordTypeStart},
						Records: []CaseRecord{
							CaseRecord{Type: CaseRecordTypeResult},
						},
					},
				},
			},
			Request: 0,
			Record:  -1,
		}
		assert.NotNil(t, replay)

		record, step := replay.NextStep()
		assert.Equal(t, CaseRecordTypeStart, record.Type)
		assert.Equal(t, 1, step)

		record, step = replay.NextStep()
		assert.Equal(t, CaseRecordTypeResult, record.Type)
		assert.Equal(t, 2, step)

		record, step = replay.NextStep()
		assert.Nil(t, record)
		assert.Equal(t, 2, step)
	})

	t.Run("one request, two Records", func(t *testing.T) {
		t.Parallel()

		replay := &CaseBindReplay{
			CaseBind: CaseBind{
				Requests: []CaseRequest{
					CaseRequest{
						Request: CaseRecord{Type: CaseRecordTypeStart},
						Records: []CaseRecord{
							CaseRecord{Type: CaseRecordTypeTraceID},
							CaseRecord{Type: CaseRecordTypeResult},
						},
					},
				},
			},
			Request: 0,
			Record:  -1,
		}
		assert.NotNil(t, replay)

		record, step := replay.NextStep()
		assert.Equal(t, CaseRecordTypeStart, record.Type)
		assert.Equal(t, 1, step)

		record, step = replay.NextStep()
		assert.Equal(t, CaseRecordTypeTraceID, record.Type)
		assert.Equal(t, 2, step)

		record, step = replay.NextStep()
		assert.Equal(t, CaseRecordTypeResult, record.Type)
		assert.Equal(t, 3, step)

		record, step = replay.NextStep()
		assert.Nil(t, record)
		assert.Equal(t, 3, step)
	})

	t.Run("two requests, no records", func(t *testing.T) {
		t.Parallel()

		replay := &CaseBindReplay{
			CaseBind: CaseBind{
				Requests: []CaseRequest{
					CaseRequest{
						Request: CaseRecord{Type: CaseRecordTypeStart},
						Records: []CaseRecord{},
					},
					CaseRequest{
						Request: CaseRecord{Type: CaseRecordTypeResult},
						Records: []CaseRecord{},
					},
				},
			},
			Request: 0,
			Record:  -1,
		}
		assert.NotNil(t, replay)

		record, step := replay.NextStep()
		assert.Equal(t, CaseRecordTypeStart, record.Type)
		assert.Equal(t, 1, step)

		record, step = replay.NextStep()
		assert.Equal(t, CaseRecordTypeResult, record.Type)
		assert.Equal(t, 2, step)

		record, step = replay.NextStep()
		assert.Nil(t, record)
		assert.Equal(t, 2, step)
	})

}
