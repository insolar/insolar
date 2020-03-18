// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/x-crypto/sha256"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

const CallUrl = "http://localhost:19192/api/rpc"

func TestTimeoutSuite(t *testing.T) {
	defer testutils.LeakTester(t)

	ctx, _ := inslogger.WithTraceField(context.Background(), "APItests")
	mc := minimock.NewController(t)
	defer mc.Wait(17 * time.Second)
	defer mc.Finish()

	sKey, err := secrets.GeneratePrivateKeyEthereum()
	require.NoError(t, err)
	sKeyString, err := secrets.ExportPrivateKeyPEM(sKey)
	require.NoError(t, err)
	pKey := secrets.ExtractPublicKey(sKey)
	pKeyString, err := secrets.ExportPublicKeyPEM(pKey)
	require.NoError(t, err)

	userRef := gen.Reference().String()
	user, err := requester.CreateUserConfig(userRef, string(sKeyString), string(pKeyString))

	cr := testutils.NewContractRequesterMock(mc)
	cr.CallMock.Set(func(p context.Context, p1 *insolar.Reference, method string, p3 []interface{}, p4 insolar.PulseNumber) (insolar.Reply, *insolar.Reference, error) {
		requestReference, _ := insolar.NewReferenceFromString("insolar:1MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI")
		switch method {
		case "Call":
			var result = "OK"
			data, _ := foundation.MarshalMethodResult(result, nil)
			return &reply.CallMethod{
				Result: data,
			}, requestReference, nil
		default:
			return nil, nil, errors.New("Unknown method: " + method)
		}
	})

	checker := testutils.NewAvailabilityCheckerMock(mc)
	checker.IsAvailableMock.Return(true)

	http.DefaultServeMux = new(http.ServeMux)
	cfg := configuration.NewAPIRunner(false)
	cfg.Address = "localhost:19192"
	cfg.SwaggerPath = "testdata/api-exported.yaml"
	api, err := NewRunner(
		&cfg,
		nil,
		cr,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		checker,
		Options{
			ContractMethods: map[string]bool{
				"contract.registerNode": true,
				"schema.get":            true,
			},
		},
	)
	require.NoError(t, err)
	defer api.Stop(ctx)
	seed, err := api.SeedGenerator.Next()
	require.NoError(t, err)

	api.SeedManager.Add(*seed, 0)

	seedString := base64.StdEncoding.EncodeToString(seed[:])

	requester.SetTimeout(25)
	req, err := requester.MakeRequestWithSeed(
		ctx,
		CallUrl,
		user,
		&requester.Params{CallSite: "contract.registerNode", CallParams: map[string]interface{}{}, PublicKey: user.PublicKey},
		seedString,
	)
	require.NoError(t, err, "make request with seed error")

	rr := httptest.NewRecorder()
	api.Handler().ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code, "got StatusOK http code")

	var result requester.ContractResponse
	// fmt.Println("response:", rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	require.NoError(t, err, "json unmarshal error")
	require.Nil(t, result.Error, "error should be nil in result")
	require.Equal(t, "OK", result.Result.CallResult, "call result is OK")
}

func TestDigestParser(t *testing.T) {
	invalidDigest := ""
	_, err := parseDigest(invalidDigest)
	require.Error(t, err)

	validDigest := "SHA-256=foo"
	_, err = parseDigest(validDigest)
	require.NoError(t, err)
}

func TestSignatureParser(t *testing.T) {
	invalidSignature := ""
	_, err := parseSignature(invalidSignature)

	validSignature := `keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=bar`
	_, err = parseSignature(validSignature)
	require.NoError(t, err)
}

func TestValidateRequestHeaders(t *testing.T) {
	body := []byte("foobar")
	h := sha256.New()
	_, err := h.Write(body)
	require.NoError(t, err)

	digest := h.Sum(nil)
	calculatedDigest := `SHA-256=` + base64.URLEncoding.EncodeToString(digest)
	signature := `keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=bar`
	sig, err := validateRequestHeaders(calculatedDigest, signature, body)
	require.NoError(t, err)
	require.Equal(t, "bar", sig)
}
