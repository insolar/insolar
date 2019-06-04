package simplerequester

import (
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/elliptic"
	"github.com/insolar/x-crypto/rand"
	"github.com/stretchr/testify/require"
	"testing"
)

const HOST = "http://localhost:19101"
const TestUrl = HOST + "/api"

var memRef string
var privateKey *ecdsa.PrivateKey

func TestCreateMember(t *testing.T) {

	var err error
	privateKey, err = ecdsa.GenerateKey(elliptic.P256K(), rand.Reader)
	require.NoError(t, err)
	seed, err := requester.GetSeed(TestUrl)
	require.NoError(t, err)
	t.Log("seed:" + string(seed))

	datas := DataToSign{
		Reference: "1tJDL5m9pKyq2mbanYfgwQ5rSQdrpsXbzc1Dk7a53d.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		Method:    "CreateMember",
		Seed:      seed,
		Params:    `{"name":"name"}`,
	}

	jws, jwk, err := createSignedData(privateKey, &datas)
	require.NoError(t, err)
	params := requester.PostParams{
		"jws": jws,
		"jwk": jwk,
	}

	body, err := requester.GetResponseBody(TestUrl+"/call", params)
	require.NoError(t, err)

	t.Log(string(body))
	response, err := getResponse(body)
	require.NoError(t, err)
	require.NotNil(t, response)

	memRef = response.Result.(string)
}

func TestGetBalance(t *testing.T) {

	seed, err := requester.GetSeed(TestUrl)
	require.NoError(t, err)
	t.Log("seed:" + string(seed))

	datas := DataToSign{
		Reference: memRef,
		Method:    "GetBalance",
		Seed:      seed,
		Params:    `{"reference":"` + memRef + `"}`,
	}

	jws, jwk, err := createSignedData(privateKey, &datas)
	require.NoError(t, err)
	params := requester.PostParams{
		"jws": jws,
		"jwk": jwk,
	}

	body, err := requester.GetResponseBody(TestUrl+"/call", params)
	require.NoError(t, err)

	t.Log(string(body))
	response, err := getResponse(body)
	require.NoError(t, err)
	require.NotNil(t, response)
}
