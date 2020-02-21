// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package examples_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"go.opencensus.io/trace"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

// run jaeger locally:
//  docker run --rm --name jaeger -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:1.7 --log-level=debug
//
// run example:
//  INSOLAR_TRACER_JAEGER_AGENTENDPOINT="localhost:6831" INSOLAR_LOG_LEVEL=debug go test -v ./instrumentation/instracer/examples/. -run=TestClientServerExample
//
// check generated traces:
//  http://localhost:16686

var testNameForExampleClientServer = "TestClientServerExample"

func TestClientServerExample(t *testing.T) {
	if os.Getenv("INSOLAR_TRACER_JAEGER_AGENTENDPOINT") == "" {
		t.Skip()
	}
	Example_ClientServer()
}

func Example_ClientServer() {
	cfgHolder := configuration.NewHolderLight("testdata/insolard-light.yaml").MustLoad()
	jconf := cfgHolder.Configuration.Tracer.Jaeger
	ctx := context.Background()
	inslogger.FromContext(ctx).Infof("jconf => %+v", jconf)

	if tracedata := os.Getenv("TRACE_SERVER"); len(tracedata) > 0 {
		servctx, tracespanbin := dataForServer(tracedata)
		donefn := instracer.ShouldRegisterJaeger(
			servctx, "server", "nodeRef", jconf.AgentEndpoint, jconf.CollectorEndpoint, jconf.ProbabilityRate)
		defer donefn()
		time.Sleep(time.Millisecond * 10)
		serverHandler(servctx, tracespanbin)
		time.Sleep(time.Millisecond * 50)
		return
	}

	traceid := fmt.Sprintf("%v", time.Now().Unix())
	ctx = inslogger.ContextWithTrace(ctx, traceid)
	donefn := instracer.ShouldRegisterJaeger(
		ctx, "client", "nodeRef", jconf.AgentEndpoint, jconf.CollectorEndpoint, jconf.ProbabilityRate)
	defer donefn()

	ctx, span := instracer.StartSpan(ctx, "root")
	defer span.Finish()

	fmt.Println("A> start")
	_, cSpan1 := instracer.StartSpan(ctx, "client-1")
	defer cSpan1.Finish()
	time.Sleep(time.Millisecond * 15)

	ctx2, cSpan2 := instracer.StartSpan(ctx, "client-2")

	requestServer(ctx2, traceid)
	cSpan2.Finish()

	time.Sleep(time.Millisecond * 150)
	fmt.Println("A> end")
}

func requestServer(ctx context.Context, traceid string) {
	fmt.Println(" A> call requestServer")
	cCtx, cSpan := instracer.StartSpan(ctx, "clientrequest")
	defer cSpan.Finish()

	cSC := trace.FromContext(cCtx).SpanContext()
	cmd := exec.Command(os.Args[0], "-test.run="+testNameForExampleClientServer)
	tracefields := strings.Join(
		[]string{
			traceid,
			string(cSC.TraceID[:]),
			string(cSC.SpanID[:]),
		}, "__")

	cmd.Env = append(os.Environ(), "TRACE_SERVER="+tracefields)

	done := make(chan error)
	go func() {
		out, err := cmd.CombinedOutput()
		fmt.Println("serverrequest output>\n", string(out))
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			done <- err
		}
		close(done)
	}()

	srverr := <-done
	if srverr != nil {
		fmt.Println("Server failed during run:", srverr)
	}
	time.Sleep(time.Millisecond * 15)
	fmt.Println(" A> end requestServer")
}

func dataForServer(tracestring string) (context.Context, []byte) {
	args := strings.SplitN(tracestring, "__", 3)
	instraceid, traceid, spanid := args[0], args[1], args[2]
	tracespan := instracer.TraceSpan{
		TraceID: []byte(traceid),
		SpanID:  []byte(spanid),
	}
	ctx := inslogger.ContextWithTrace(context.Background(), instraceid)
	b, err := tracespan.Serialize()
	if err != nil {
		panic(err)
	}
	return ctx, b
}

func serverHandler(ctx context.Context, spanbin []byte) {
	fmt.Println(" B> call serverHandler")

	parentSpan := instracer.MustDeserialize(spanbin)
	ctx = instracer.WithParentSpan(ctx, parentSpan)

	fmt.Println("  ... instracer.StartSpan")
	spanctx, servspan := instracer.StartSpan(ctx, "server")
	defer servspan.Finish()
	time.Sleep(time.Millisecond * 20)
	serverSubCalls(spanctx)
	time.Sleep(time.Millisecond * 40)
}

func serverSubCalls(ctx context.Context) {
	fmt.Println("  B> SubCall-1")

	_, servspan1 := instracer.StartSpan(ctx, "subcall-1")
	time.Sleep(time.Millisecond * 13)
	servspan1.Finish()

	fmt.Println("  B> SubCall-2")
	_, servspan2 := instracer.StartSpan(ctx, "subcall-2")
	time.Sleep(time.Millisecond * 27)
	servspan2.Finish()
}
