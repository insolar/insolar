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

/*
Package instracer contains helpers for opencesus tracer.

Example:

	import (
		"github.com/insolar/insolar/configuration"
		"github.com/insolar/insolar/log"
	)

	// on client side
	//
	entryvalue := "entryvalue"
	ctx := context.Background()
	jaegerflush = instracer.ShouldRegisterJaeger(ctx, "insolard", "localhost:6831", "")
	ctx = instracer.SetBaggage(ctx, instracer.Entry{Key: "someentry", Value: entryvalue})
	defer jaegerflush() // wait all trace data to send on jaeger server

	// serialize clientctx
	spanbindata := instracer.MustSerialize(ctx)

	// send spanbindata on wire with request
	// someSendMethod(ctxdata, request)

	// on server side
	//
	// deserialized from wire
	// spanbindata := someRecieverMethod()
	ctx := context.Background()
	instracer.MustDeserialize(spanbindata)

	ctx = instracer.WithParentSpan(ctx, parentspan)
	donefn := instracer.ShouldRegisterJaeger(ctx, "server", "localhost:6831", "")
	defer donefn()

	servctx, servspan := instracer.StartSpan(ctx, "server")
	defer servspan.End()
	// call subrequests with servctx, and use instracer.StartSpan


Hints:

Use environment variables for log level setup:

	INSOLAR_TRACER_JAEGER_AGENTENDPOINT="localhost:6831"

How to run Jaeger locally:

	docker run --rm --name jaeger \
		-p 6831:6831/udp \
		-p 16686:16686 \
		jaegertracing/all-in-one:1.7 --log-level=debug
*/
package instracer
