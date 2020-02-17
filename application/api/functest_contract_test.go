// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalUpload(t *testing.T) {
	jsonResponse := `
{
    "jsonrpc": "2.0",
    "result": {
        "Test": "Test",
        "PrototypeRef": "6R46iNSizv7pzHrLiR8m1qtEPC9FvLtsdKoFV9w2r6V.11111111111111111111111111111111"
    },
    "id": ""
}`
	res := struct {
		Version string      `json:"jsonrpc"`
		ID      string      `json:"id"`
		Result  UploadReply `json:"result"`
	}{}

	expectedRes := struct {
		Version string      `json:"jsonrpc"`
		ID      string      `json:"id"`
		Result  UploadReply `json:"result"`
	}{
		Version: "2.0",
		ID:      "",
	}

	err := json.Unmarshal([]byte(jsonResponse), &res)
	require.NoError(t, err)

	require.Equal(t, expectedRes.Version, res.Version)
	require.Equal(t, expectedRes.ID, res.ID)
	require.NotEqual(t, "", res.Result.PrototypeRef)
}
