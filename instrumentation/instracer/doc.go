// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
