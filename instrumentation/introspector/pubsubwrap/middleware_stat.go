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
	"context"
	"sort"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/instrumentation/introspector/introproto"
)

// MessageStatByType holds publish statistic per message type.
type MessageStatByType struct {
	sync.Mutex
	stat map[string]int64
}

// NewMessageStatByType is a constructor for MessageStatByType.
func NewMessageStatByType() *MessageStatByType {
	return &MessageStatByType{
		stat: map[string]int64{},
	}
}

// Filter counts published messages by type.
func (ms *MessageStatByType) Filter(m *message.Message) (*message.Message, error) {
	typ, err := decodeType(m)
	key := typ.String()
	if err != nil {
		if de, ok := err.(decodeError); ok {
			key = "legacy." + de.metadataType
		}
	}
	ms.Lock()
	defer ms.Unlock()

	ms.stat[key]++
	return m, err
}

// GetMessagesStat returns publish statistic per message type.
func (ms *MessageStatByType) GetMessagesStat(context.Context, *introproto.EmptyArgs) (*introproto.AllMessageStatByType, error) {
	ms.Lock()
	all := make([]*introproto.MessageStatByType, 0, len(ms.stat))
	for name, count := range ms.stat {
		all = append(all, &introproto.MessageStatByType{
			Name:  name,
			Count: count,
		})
	}
	ms.Unlock()

	// frequent messages goes first
	sort.Slice(all, func(i, j int) bool {
		return all[i].Count > all[j].Count
	})

	return &introproto.AllMessageStatByType{
		Counters: all,
	}, nil
}
