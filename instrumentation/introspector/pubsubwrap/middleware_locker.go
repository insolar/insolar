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

package pubsubwrap

import (
	"context"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/introspector/introproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stat struct {
	sync.Mutex
	counters map[payload.Type]int64
}

func (s *stat) get(typ payload.Type) int64 {
	s.Lock()
	defer s.Unlock()
	return s.counters[typ]
}

func (s *stat) inc(typ payload.Type) {
	s.Lock()
	s.counters[typ]++
	s.Unlock()
}

type MessageLockerByType struct {
	sync.Mutex
	types map[payload.Type]struct{}
	stat  *stat
	log   insolar.Logger
}

// NewMessageLockerByType is a constructor for MessageLockerByType.
func NewMessageLockerByType(ctx context.Context) *MessageLockerByType {
	inslog := inslogger.FromContext(ctx)
	return &MessageLockerByType{
		types: map[payload.Type]struct{}{},
		stat: &stat{
			counters: map[payload.Type]int64{},
		},
		log: inslog,
	}
}

func (ml *MessageLockerByType) typeIsFiltered(pt payload.Type) bool {
	ml.Lock()
	defer ml.Unlock()
	_, ok := ml.types[pt]
	return ok
}

// Filter returns nil for filtered message types.
func (ml *MessageLockerByType) Filter(m *message.Message) *message.Message {
	var meta payload.Meta
	err := meta.Unmarshal(m.Payload)
	if err != nil {
		panic(err)
	}

	typ, err := payload.UnmarshalType(meta.Payload)
	if err != nil {
		panic(err)
	}

	if ml.typeIsFiltered(typ) {
		ml.stat.inc(typ)
		ml.log.Debugf("MessageLocker filtered '%v'", typ.String())
		return nil
	}
	return m
}

// SetMessagesFilter sets filter for provided message type.
func (ml *MessageLockerByType) SetMessagesFilter(ctx context.Context, in *introproto.MessageFilterByType) (*introproto.MessageFilterByType, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name shouldn't be empty")
	}
	if _, ok := payload.TypesMap[in.Name]; !ok {
		return nil, status.Errorf(codes.InvalidArgument, "'%v' unknown message payload type", in.Name)
	}

	ml.set(in.Name, in.Enable)
	return in, nil
}

func (ml *MessageLockerByType) set(name string, enable bool) {
	ml.Lock()
	defer ml.Unlock()

	typ := payload.TypesMap[name]
	if enable {
		ml.types[typ] = struct{}{}
	} else {
		delete(ml.types, typ)
	}
}

// GetMessagesFilters returns filter state and statistic per message type (as map).
func (ml *MessageLockerByType) GetMessagesFilters(ctx context.Context, in *introproto.EmptyArgs) (*introproto.AllMessageFilterStats, error) {
	ml.Lock()
	defer ml.Unlock()

	filters := make(map[string]*introproto.MessageFilterWithStat)
	for name, typ := range payload.TypesMap {
		_, ok := ml.types[typ]
		filters[name] = &introproto.MessageFilterWithStat{
			Enable:   ok,
			Filtered: ml.stat.get(typ),
		}
	}

	return &introproto.AllMessageFilterStats{
		Filters: filters,
	}, nil
}
