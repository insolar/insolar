///
// Copyright 2020 Insolar Technologies GmbH
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
///

// +build functest

package functest

import (
	"bytes"
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/instrumentation/insmetrics"
)

// This test file contains tests what always must be last in the package.

// TestMnt_DumpMetrics saves metrics values to files in launchnet logs dir.
func TestMnt_DumpMetrics(t *testing.T) {
	if !launchnet.DumpMetricsEnabled() {
		t.Skip("dump metrics disabled")
	}

	res, err := launchnet.FetchAndSaveMetrics(functestCount-1, AppPath)
	if err != nil {
		t.Errorf("metrics save failed: %v", err.Error())
	}
	var inc float64
	for _, b := range res {
		inc += insmetrics.SumMetricsValueByNamePrefix(bytes.NewReader(b), "insolar_requests_abandoned")
	}
	t.Logf("Abandons sum: %v (functestCount=%v)", inc, functestCount)
}
