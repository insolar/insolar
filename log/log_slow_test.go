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

// +build slowtest

package log

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

func Test_Log_Deviations(t *testing.T) {
	logStorage := NewConcurrentBuilder(500000)
	logger, err := NewTestLogger(logStorage)
	require.NoError(t, err)

	generateLogs(logger, 1000, 1000)

	err = logStorage.Close()
	require.NoError(t, err)
	for !logStorage.IsFlushed() {
		time.Sleep(100 * time.Millisecond)
	}

	logString := logStorage.String()
	out, err := parseOutput(logString)
	require.NoError(t, err)

	t.Run("log sequence", func(t *testing.T) {
		checkLogSequence(t, out)
	})
}

func checkLogSequence(t *testing.T, out []logOutput) {
	lastVal := make(map[uint64]logOutput)
	for _, o := range out {
		lv, ok := lastVal[o.Thread]
		if ok {
			assert.Equal(t, int(lv.Iteration+1), int(o.Iteration),
				"Bad sequence in thread %d. Last iteration: %d, current iteration: %d.",
				o.Thread, lv.Iteration, o.Iteration,
			)
		}
		lastVal[o.Thread] = o
	}
}

func generateLogs(logger insolar.Logger, threads, iterations int) {
	var start, finish sync.WaitGroup
	start.Add(threads)
	finish.Add(threads)
	for i := 0; i < threads; i++ {
		threadID := i
		l, _ := logger.Copy().Build()
		go func() {
			start.Wait()
			for j := 0; j < iterations; j++ {
				l.Warn(threadID, j, time.Now().UnixNano())
			}
			finish.Done()
		}()
		start.Done()
	}
	finish.Wait()
}

func parseOutput(o string) ([]logOutput, error) {
	out := make([]logOutput, 0)
	for _, s := range strings.Split(o, "\n") {
		if len(s) == 0 {
			continue
		}
		ll := logLine{}
		err := json.Unmarshal([]byte(s), &ll)
		if err != nil {
			return nil, err
		}

		fields := strings.Fields(ll.Message)
		threadID, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return nil, err
		}
		iteration, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, err
		}
		tsInt, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return nil, err
		}
		ts := time.Unix(0, int64(tsInt))
		// fmt.Printf("%02d %02d %d\n", threadID, iteration, ts.UnixNano())
		out = append(out, logOutput{threadID, iteration, ts})
	}
	return out, nil
}

type logLine struct {
	Time    time.Time
	Message string
}

type logOutput struct {
	Thread, Iteration uint64
	Time              time.Time
}

func NewTestLogger(w io.Writer) (insolar.Logger, error) {
	logger, err := NewLog(configuration.Log{Level: "info", Adapter: "zerolog", Formatter: "json"})
	if err != nil {
		return nil, err
	}
	logger, err = logger.Copy().WithOutput(w).Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// ConcurrentBuilder is a simple thread safe io.Writer implementation based on strings.Builder.
type ConcurrentBuilder struct {
	builder             *strings.Builder
	queue               chan []byte
	isClosed, isFlushed bool
	lock                sync.RWMutex
}

func NewConcurrentBuilder(bufSize int) *ConcurrentBuilder {
	cb := ConcurrentBuilder{
		builder: &strings.Builder{},
		queue:   make(chan []byte, bufSize),
	}
	go cb.loop()
	return &cb
}

func (cb *ConcurrentBuilder) loop() {
	for data := range cb.queue {
		cb.builder.Write(data)
	}
	cb.lock.Lock()
	defer cb.lock.Unlock()
	cb.isFlushed = true
}

func (cb *ConcurrentBuilder) Write(p []byte) (int, error) {
	cb.lock.RLock()
	defer cb.lock.RUnlock()
	if cb.isClosed {
		return 0, errors.New("writer is closed")
	}
	data := make([]byte, len(p))
	n := copy(data, p)
	cb.queue <- data
	return n, nil
}

func (cb *ConcurrentBuilder) Close() error {
	cb.lock.Lock()
	defer cb.lock.Unlock()
	if cb.isClosed {
		return errors.New("writer is closed")
	}
	close(cb.queue)
	cb.isClosed = true
	return nil
}

func (cb *ConcurrentBuilder) String() string {
	return cb.builder.String()
}

func (cb *ConcurrentBuilder) IsFlushed() bool {
	cb.lock.RLock()
	defer cb.lock.RUnlock()
	return cb.isFlushed
}
