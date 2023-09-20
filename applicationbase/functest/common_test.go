// +build functest

package functest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
)

func TestGetRequest(t *testing.T) {
	postResp, err := http.Get(launchnet.TestRPCUrl)
	require.NoError(t, err)
	defer postResp.Body.Close()
	require.Equal(t, http.StatusMethodNotAllowed, postResp.StatusCode)
}

func TestWrongUrl(t *testing.T) {
	jsonValue, _ := json.Marshal(testresponse.PostParams{})
	testURL := launchnet.AdminHostPort + "/not_api"
	postResp, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonValue))
	defer postResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, postResp.StatusCode)
}

func TestWrongJson(t *testing.T) {
	postResp, err := http.Post(launchnet.TestRPCUrl, "application/json", bytes.NewBuffer([]byte("some not json value")))
	defer postResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)

	response := &requester.ContractResponse{}
	testresponse.UnmarshalCallResponse(t, body, response)
	require.NotNil(t, response.Error)

	require.Equal(t, "The JSON received is not a valid request payload.", response.Error.Message)
	require.Nil(t, response.Result)
}
