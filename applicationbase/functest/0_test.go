// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
)

// This test file contains tests what always must be first in the package.

var functestCount int

// TestMnt_IterationCounter counts which iteration of functest is.
// Because go test framework doesn't provide such info right now.
func TestMnt_IterationCounter(t *testing.T) {
	functestCount++
	t.Log("functest iteration:", functestCount)
}

// TestMnt_RotateLogs rotates launchnet logs (removes and reopens them).
func TestMnt_RotateLogs(t *testing.T) {
	if !launchnet.LogRotateEnabled() {
		t.Skip("log rotate disabled")
	}
	launchnet.RotateLogs(true)
}
