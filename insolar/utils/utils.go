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
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
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

type SyncT struct {
	*testing.T

	mu sync.Mutex
}

var _ testing.TB = (*SyncT)(nil)

func (t *SyncT) Error(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Error(args...)
}
func (t *SyncT) Errorf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Errorf(format, args...)
}
func (t *SyncT) Fail() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Fail()
}
func (t *SyncT) FailNow() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.FailNow()
}
func (t *SyncT) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.T.Failed()
}
func (t *SyncT) Fatal(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Fatal(args...)
}
func (t *SyncT) Fatalf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Fatalf(format, args...)
}
func (t *SyncT) Log(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Log(args...)
}
func (t *SyncT) Logf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Logf(format, args...)
}
func (t *SyncT) Name() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.T.Name()
}
func (t *SyncT) Skip(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Skip(args...)
}
func (t *SyncT) SkipNow() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.SkipNow()
}
func (t *SyncT) Skipf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.T.Skipf(format, args...)
}
func (t *SyncT) Skipped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.T.Skipped()
}
func (t *SyncT) Helper() {}
