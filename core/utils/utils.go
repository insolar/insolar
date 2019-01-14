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
	"context"
	"encoding/binary"
	"os"

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

func MeasureExecutionTime(ctx context.Context, comment string, thefunction func()) {
	// TODO FIXME
	thefunction()
}

/*var measurementsEnabled = false

// EnableExecutionTimeMeasurement enables execution time measurement
// and uses `jaegerAddr` (host:port of Jaeger service) to write traces.
func EnableExecutionTimeMeasurement(jaegerAddr string) error {
	if measurementsEnabled {
		// already enabled
		return nil
	}

	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: jaegerAddr,
		Process: jaeger.Process{
			ServiceName: "testervice", // TODO FIXME
			Tags:        []jaeger.Tag{
				// You can specify some global tags here, e.g:
				// jaeger.StringTag("hostname", "localhost"),
			},
		},
	})
	if err != nil {
		return err
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	measurementsEnabled = true
	return nil
}*/

// MeasureExecutionTime writes execution time of given function to
// the profile log (if profile logging is enabled).
/*func MeasureExecutionTime(ctx context.Context, comment string, thefunction func()) {
	if !measurementsEnabled {
		thefunction()
		return
	}

	traceID := TraceID(ctx)

	ctx, span := trace.StartSpan(ctx, comment) // TODO: will not work like this, should pass the new ctx to thefunction
	span.AddAttributes(
		trace.StringAttribute("insolarTraceId", traceID),
	)
	thefunction()
	span.End()
}*/
