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

package utils

import (
	"context"
	"encoding/binary"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type traceIDKey struct{}

// TraceID returns traceid provided by WithTraceField and ContextWithTrace helpers.
func TraceID(ctx context.Context) string {
	val := ctx.Value(traceIDKey{})
	if val == nil {
		return ""
	}
	return val.(string)
}

func SetInsTraceID(ctx context.Context, traceid string) (context.Context, error) {
	if TraceID(ctx) != "" {
		return context.WithValue(ctx, traceIDKey{}, traceid),
			errors.Errorf("TraceID already set: old: %s new: %s", TraceID(ctx), traceid)
	}
	return context.WithValue(ctx, traceIDKey{}, traceid), nil
}

// RandTraceID returns random traceID in uuid format.
func RandTraceID() string {
	traceID, err := uuid.NewV4()
	if err != nil {
		return "createRandomTraceIDFailed:" + err.Error()
	}
	return traceID.String()
}

func UInt32ToBytes(n uint32) []byte {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, n)
	return buff
}

// CircleXOR performs XOR for 'value' and 'src'. The result is returned as new byte slice.
// If 'value' is smaller than 'dst', XOR starts from the beginning of 'src'.
func CircleXOR(value, src []byte) []byte {
	result := make([]byte, len(value))
	srcLen := len(src)
	for i := range result {
		result[i] = value[i] ^ src[i%srcLen]
	}
	return result
}
