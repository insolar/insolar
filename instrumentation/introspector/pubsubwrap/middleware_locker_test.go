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

package pubsubwrap

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/introspector/introproto"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareLocker(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mi := NewMessageLockerByType(ctx)

	expected := []struct {
		payloadType payload.Type
		payload     payload.Payload
		locks       int
		total       int
	}{
		{
			payload.TypeGetObject,
			&payload.GetObject{Polymorph: uint32(payload.TypeGetObject)},
			4,
			10,
		},
		{
			payload.TypeGetCode,
			&payload.GetCode{Polymorph: uint32(payload.TypeGetCode)},
			10,
			10,
		},
	}

	for _, ex := range expected {
		b, err := ex.payload.Marshal()
		require.NoError(t, err, "payload should be marshaled w/o errors")
		var meta payload.Meta
		meta.Payload = b
		metaBytes, err := meta.Marshal()
		require.NoError(t, err, "meta should be marshaled w/o errors")
		msg := &message.Message{
			Payload: metaBytes,
		}

		filter := &introproto.MessageFilterByType{
			Name:   ex.payloadType.String(),
			Enable: true,
		}
		for i := 0; i < ex.total; i++ {
			filter.Enable = i < ex.locks
			_, setErr := mi.SetMessagesFilter(ctx, filter)
			require.NoError(t, setErr, "SetMessagesFilter should not return error")

			mi.Filter(msg)
		}
	}

	stats, err := mi.GetMessagesFilters(nil, nil)
	require.NoError(t, err, "GetMessagesFilters should not failed")

	for _, ex := range expected {
		typeName := ex.payloadType.String()
		for _, filter := range stats.Filters {
			if filter.Name == typeName {
				require.Equalf(t, int64(ex.locks), filter.Filtered, "check %v stat", typeName)
			}
		}
	}
}
