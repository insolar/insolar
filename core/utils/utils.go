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
	"encoding/binary"
	"log"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

// RandTraceID returns random traceID in uuid format
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

// Returns current timestamp in milliseconds
func TimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Writes execution time of given function to the profile log (if profile logging is enabled)
// TODO: use seperate log file! + enable/disable flag
func MeasureExecutionTime(comment string, f func()) {
	start := TimestampMs()
	log.Printf("[PROFILE] %s - STARTED @ %v\n", comment, start)
	f()
	end := TimestampMs()
	delta := end - start
	log.Printf("[PROFILE] %s - ENDED @ %v, delta: %v ms\n", comment, end, delta)
}
