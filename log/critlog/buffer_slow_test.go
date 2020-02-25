// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package critlog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/stat"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logoutput"
)

func Test_BackpressureBuffer_Deviations(t *testing.T) {

	generateDelay := time.Microsecond
	writeDelay := generateDelay * 2

	threads, iterations := 10, 1000
	t.Run(fmt.Sprintf("th=%d iter=%d buf=%d - bypass", threads, iterations, 0), func(t *testing.T) {
		test_BackpressureBuffer_Deviations(t, threads, iterations, 0, 5, false, generateDelay, writeDelay)
	})

	for _, threads := range []int{1, 10, 100, 1000} {
		iterations := 5000 / threads
		if iterations > 1000 {
			iterations = 1000
		}
		for _, bufSize := range []int{10, 100, 1000} {
			for _, parWrites := range []int{1, 5, 10, 100} {
				if parWrites > threads {
					continue
				}
				t.Run(fmt.Sprintf("th=%d iter=%d buf=%d wr=%d m=unfair", threads, iterations, bufSize, parWrites), func(t *testing.T) {
					test_BackpressureBuffer_Deviations(t, threads, iterations, bufSize, parWrites, false, generateDelay, writeDelay)
				})
				t.Run(fmt.Sprintf("th=%d iter=%d buf=%d wr=%d m=fair", threads, iterations, bufSize, parWrites), func(t *testing.T) {
					test_BackpressureBuffer_Deviations(t, threads, iterations, bufSize, parWrites, true, generateDelay, writeDelay)
				})
			}
		}
	}
}

func test_BackpressureBuffer_Deviations(t *testing.T, threads, iterations, bufSize, parWrites int, fair bool,
	generateDelay, writeDelay time.Duration) {

	logStorage := NewConcurrentBuilder(threads*iterations+iterations, writeDelay)
	logger := NewTestLogger(context.Background(), logStorage, uint8(parWrites), bufSize, fair)

	generateLogs(logger, threads, iterations, generateDelay)

	err := logger.Close()
	require.NoError(t, err)
	err = logStorage.Close()
	require.EqualError(t, err, "writer is closed")
	for !logStorage.IsFlushed() {
		time.Sleep(100 * time.Millisecond)
	}

	logString := logStorage.String()

	out, err := parseOutput(logString)
	require.NoError(t, err)

	/* ============================ */

	lastVal := make(map[uint64]logOutput, threads)
	distances := make(map[uint64][]float64, threads)
	allDistances := make([]float64, 0, threads*iterations)
	displacements := make([]float64, 0)
	for _, v := range out {
		lv, ok := lastVal[v.Thread]
		if ok {
			dist := v.Time.Sub(lv.Time).Seconds()
			distances[v.Thread] = append(distances[v.Thread], dist)
			allDistances = append(allDistances, dist)

			if v.Iteration < lv.Iteration {
				displacements = append(displacements, float64(lv.Iteration-v.Iteration))
				// fmt.Printf("\tCurrent < last in thread %d: current=%d last=%d\n", v.Thread, v.Iteration, lv.Iteration)
			}
		}
		lastVal[v.Thread] = v
	}

	ttlMean, ttlStd := stat.MeanStdDev(allDistances, nil)
	ttlSum := float64(0)
	for _, d := range allDistances {
		ttlSum += d
	}

	for k, v := range distances {
		if iterations != len(v)+1 {
			assert.Equal(t, iterations, len(v)+1, "Incorrect number of log records in thread %d", k)
		}

		mean, std := stat.MeanStdDev(v, nil)
		_, _ = mean, std
		// fmt.Printf("Thread %03d: mean = %8.2f ms %+6.2f%%, stddev = %8.2f ms %+6.2f%%\n", k,
		// 	mean*1e3, 100*(mean-ttlMean)/ttlMean, std*1e3, 100*(std-ttlStd)/ttlStd)
	}

	fmt.Printf("\tTotal: sum = %8.2f s, mean = %8.2f ms, stddev = %8.2f ms\n", ttlSum, ttlMean*1e3, ttlStd*1e3)
	meanDisplace, stdDisplace := stat.MeanStdDev(displacements, nil)
	fmt.Printf("\tDisplacements: total = %d, mean = %.2f (std = %.2f)\n", len(displacements), meanDisplace, stdDisplace)
}

func generateLogs(logger insolar.LoggerOutput, threads, iterations int, generateDelay time.Duration) {
	var start, finish sync.WaitGroup
	start.Add(threads)
	finish.Add(threads)
	for i := 0; i < threads; i++ {
		threadID := i
		go func() {
			start.Wait()
			for j := 0; j < iterations; j++ {
				msg := fmt.Sprintf(`{"message":"%d %d %d"}%s`, threadID, j, time.Now().UnixNano(), "\n")

				_, _ = logger.LogLevelWrite(insolar.InfoLevel, []byte(msg))
				if generateDelay > 0 {
					time.Sleep(generateDelay)
				}
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

func NewTestLogger(ctx context.Context, w io.Writer, parWrites uint8, bufSize int, fair bool) insolar.LoggerOutput {
	adapter := logoutput.NewAdapter(
		w,
		false,
		func() error {
			return nil
		},
		func() error {
			return nil
		},
	)

	if parWrites == 0 {
		return NewFatalDirectWriter(adapter)
	}
	flags := BufferWriteDelayFairness | BufferTrackWriteDuration
	if !fair {
		flags = 0
	}

	var bp *BackpressureBuffer
	if bufSize == 0 {
		bp = NewBackpressureBufferWithBypass(adapter, 100, parWrites, flags, nil)
	} else {
		bp = NewBackpressureBuffer(adapter, bufSize, parWrites, flags, nil)
	}

	bp.StartWorker(ctx)
	return bp
}

// ConcurrentBuilder is a simple thread safe io.Writer implementation based on strings.Builder.
type ConcurrentBuilder struct {
	builder             *strings.Builder
	queue               chan []byte
	isClosed, isFlushed bool
	lock                sync.RWMutex
	writeDelay          time.Duration
}

func NewConcurrentBuilder(bufSize int, delay time.Duration) *ConcurrentBuilder {
	cb := ConcurrentBuilder{
		builder:    &strings.Builder{},
		queue:      make(chan []byte, bufSize),
		writeDelay: delay,
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
	if cb.writeDelay > 0 {
		time.Sleep(cb.writeDelay)
	}
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
