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

package insmetrics

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtils_ValueByNamePrefix(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/openmetrics.txt")
	require.NoError(t, err, "metrics file open")

	r := func() io.Reader { return bytes.NewReader(b) }
	v1 := SumMetricsValueByNamePrefix(r(), "insolar_bus_sent_milliseconds_bucket")
	assert.Equal(t, float64(65795), v1, "check bucket sum")

	v2 := SumMetricsValueByNamePrefix(r(), "insolar_bus_sent_milliseconds_sum")
	assert.Equal(t, 1560.873845, v2, "check bucket sum")
}
