// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package testutils

import (
	"go.uber.org/goleak"
)

func LeakTester(t goleak.TestingT) {
	goleak.VerifyNone(t,
		// Iterate Unit Tests uses memprofile
		goleak.IgnoreTopFunction("runtime/pprof.readProfile"),
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
		// sometimes stack has full import path
		goleak.IgnoreTopFunction("github.com/insolar/insolar/vendor/go.opencensus.io/stats/view.(*worker).start"),
		goleak.IgnoreTopFunction("github.com/insolar/insolar/log/critlog.(*internalBackpressureBuffer).worker"))
}
