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

package utils

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

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

func SetTraceID(ctx context.Context, traceid string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceid)
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

func SendGracefulStopSignal() error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Signal(os.Interrupt)
}

// TimestampMs returns current timestamp in milliseconds.
func TimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// We have to use a lock since all IO is not thread safe by default
var measurementsLock sync.Mutex
var measurementsWriter *bufio.Writer
var measurementsEnabled = false

// write one measure to the log
func writeMeasure(format string, args ...interface{}) error {
	measurementsLock.Lock()
	_, err := fmt.Fprintf(measurementsWriter, format, args...)
	if err != nil { // very unlikely to happen
		measurementsEnabled = false
		measurementsLock.Unlock()
		return err
	}

	err = measurementsWriter.Flush()
	measurementsLock.Unlock()
	if err != nil {
		measurementsEnabled = false
	}
	return err
}

// EnableExecutionTimeMeasurement enables execution time measurement
// and uses `fname` to write measurements.
func EnableExecutionTimeMeasurement(fname string) (func(), error) {
	if measurementsEnabled {
		// already enabled
		return func() {}, nil
	}
	// if the file doesn't exist, create it, or append to the file
	mfile, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	measurementsWriter = bufio.NewWriter(mfile)
	measurementsEnabled = true

	cleanup := func() {
		if !measurementsEnabled {
			// an error has occurred during the execution of the program
			// don't flush any buffers in this case
			return
		}

		measurementsLock.Lock()
		_ = measurementsWriter.Flush()
		_ = mfile.Sync()
		_ = mfile.Close()
		measurementsLock.Unlock()

		measurementsEnabled = false
	}
	return cleanup, nil
}

// MeasureExecutionTime writes execution time of given function to
// the profile log (if profile logging is enabled).
func MeasureExecutionTime(ctx context.Context, comment string, thefunction func()) {
	if !measurementsEnabled {
		thefunction()
		return
	}

	traceID := TraceID(ctx)

	start := TimestampMs()
	err := writeMeasure("%v %s STARTED %s\n", start, traceID, comment)
	if err != nil {
		return
	}

	thefunction()

	end := TimestampMs()
	delta := end - start
	_ = writeMeasure("%v %s ENDED %s, took: %v ms\n", end, traceID, comment, delta)
}
