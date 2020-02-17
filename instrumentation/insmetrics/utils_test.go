// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
