// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
