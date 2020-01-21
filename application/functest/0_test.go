// Copyright 2020 Insolar Network Ltd.
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

// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/application/testutils/launchnet"
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
