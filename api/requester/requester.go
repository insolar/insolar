//
// Copyright 2019 Insolar Technologies GmbH
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
//

package requester

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

var httpClient *http.Client

const (
	RequestTimeout = 15 * time.Second
)

func init() {
	httpClient = createHTTPClient()
}

func SetTimeout(timeout uint) {
	if timeout > 0 {
		httpClient.Timeout = time.Duration(timeout) * time.Second
	} else {
		httpClient.Timeout = RequestTimeout
	}
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   RequestTimeout,
	}

	return client
}

// verbose switches on verbose mode
var verbose = false
var scheme = platformpolicy.NewPlatformCryptographyScheme()

func verboseInfo(ctx context.Context, msg string) {
	if verbose {
		inslogger.FromContext(ctx).Info(msg)
	}
}

// SetVerbose switches on verbose mode
func SetVerbose(verb bool) {
	verbose = verb
}

// PlatformRequest represents params struct
type PlatformRequest struct {
	JSONRPC        string      `json:"jsonrpc"`
	ID             int         `json:"id"`
	Method         string      `json:"method"`
	PlatformParams interface{} `json:"params"`
	LogLevel       string      `json:"logLevel,omitempty"`
}

// GetResponseBodyContract makes request to contract and extracts body
func GetResponseBodyContract(url string, postP Request, signature string) ([]byte, error) {
	jsonValue, err := json.Marshal(postP)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyContract ] Problem with marshaling params")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyContract ] Problem with creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	h := sha256.New()
	_, err = h.Write(jsonValue)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyContract ] Cant get hash")
	}
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	req.Header.Set("Digest", "SHA-256="+sha)
	req.Header.Set("Signature", "keyId=\"member-pub-key\", algorithm=\"ecdsa\", headers=\"digest\", signature="+signature)
	postResp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyContract ] Problem with sending request")
	}

	if postResp == nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyContract ] Reponse is nil")
	}

	defer postResp.Body.Close()
	if http.StatusOK != postResp.StatusCode {
		return nil, errors.New("[ getResponseBodyContract ] Bad http response code: " + strconv.Itoa(postResp.StatusCode))
	}

	body, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyContract ] Problem with reading body")
	}

	return body, nil
}

// GetResponseBodyContract makes request to platform and extracts body
func GetResponseBodyPlatform(url string, postP PlatformRequest) ([]byte, error) {
	jsonValue, err := json.Marshal(postP)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyPlatform ] Problem with marshaling params")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyPlatform ] Problem with creating request")
	}
	req.Header.Set("Content-Type", "application/json")
	postResp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyPlatform ] Problem with sending request")
	}

	if postResp == nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyPlatform ] Reponse is nil")
	}

	defer postResp.Body.Close()
	if http.StatusOK != postResp.StatusCode {
		return nil, errors.New("[ GetResponseBodyPlatform ] Bad http response code: " + strconv.Itoa(postResp.StatusCode))
	}

	body, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyPlatform ] Problem with reading body")
	}

	return body, nil
}

// GetSeed makes rpc request to node.GetSeed method and extracts it
func GetSeed(url string) (string, error) {
	body, err := GetResponseBodyPlatform(url+"/rpc", PlatformRequest{
		JSONRPC: "2.0",
		Method:  "node.GetSeed",
		ID:      1,
	})
	if err != nil {
		return "", errors.Wrap(err, "[ GetSeed ] seed request")
	}

	seedResp := rpcSeedResponse{}

	err = json.Unmarshal(body, &seedResp)
	if err != nil {
		return "", errors.Wrap(err, "[ GetSeed ] Can't unmarshal")
	}
	if seedResp.Error != nil {
		return "", errors.New("[ GetSeed ] Field 'error' is not nil: " + fmt.Sprint(seedResp.Error))
	}
	if len(seedResp.Result.Seed) == 0 {
		return "", errors.New("[ GetSeed ] Field seed is empty")
	}

	return seedResp.Result.Seed, nil
}

// SendWithSeed sends request with known seed
func SendWithSeed(ctx context.Context, url string, userCfg *UserConfigJSON, reqCfg *Request, seed string) ([]byte, error) {
	if userCfg == nil || reqCfg == nil {
		return nil, errors.New("[ SendWithSeed ] Configs must be initialized")
	}

	ks := platformpolicy.NewKeyProcessor()

	pem, err := ks.ExportPublicKeyPEM(userCfg.privateKeyObject.(*ecdsa.PrivateKey).Public())
	if err != nil {
		return nil, errors.Wrap(err, "[ SendWithSeed ] Cant export public key to PEM")
	}

	reqCfg.Params.Reference = userCfg.Caller
	reqCfg.Params.PublicKey = string(pem)
	reqCfg.Params.Seed = seed

	verboseInfo(ctx, "Signing request ...")
	dataToSign, err := json.Marshal(reqCfg)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendWithSeed ] Config request marshaling failed")
	}
	signature, err := sign(userCfg.privateKeyObject, dataToSign)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendWithSeed ] Problem with signing request")
	}
	verboseInfo(ctx, "Signing request completed")

	body, err := GetResponseBodyContract(url, *reqCfg, signature)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendWithSeed ] Problem with sending target request")
	}

	return body, nil
}

func sign(privateKey crypto.PrivateKey, data []byte) (string, error) {
	hash := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hash[:])
	if err != nil {
		return "", errors.Wrap(err, "[ sign ] Cant sign data")
	}

	return PointsToDER(r, s), nil
}

func PointsToDER(r, s *big.Int) string {
	prefixPoint := func(b []byte) []byte {
		if len(b) == 0 {
			b = []byte{0x00}
		}
		if b[0]&0x80 != 0 {
			paddedBytes := make([]byte, len(b)+1)
			copy(paddedBytes[1:], b)
			b = paddedBytes
		}
		return b
	}

	rb := prefixPoint(r.Bytes())
	sb := prefixPoint(s.Bytes())

	// DER encoding:
	// 0x30 + z + 0x02 + len(rb) + rb + 0x02 + len(sb) + sb
	length := 2 + len(rb) + 2 + len(sb)

	der := append([]byte{0x30, byte(length), 0x02, byte(len(rb))}, rb...)
	der = append(der, 0x02, byte(len(sb)))
	der = append(der, sb...)

	return base64.StdEncoding.EncodeToString(der)
}

// Send first gets seed and after that makes target request
func Send(ctx context.Context, url string, userCfg *UserConfigJSON, reqCfg *Request) ([]byte, error) {
	verboseInfo(ctx, "Sending GETSEED request ...")
	seed, err := GetSeed(url)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with getting seed")
	}
	verboseInfo(ctx, "GETSEED request completed. seed: "+seed)

	response, err := SendWithSeed(ctx, url+"/call", userCfg, reqCfg, seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ]")
	}

	return response, nil
}

func getDefaultRPCParams(method string) PlatformRequest {
	return PlatformRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  method,
	}
}

// Info makes rpc request to network.GetInfo method and extracts it
func Info(url string) (*InfoResponse, error) {
	params := getDefaultRPCParams("network.GetInfo")

	body, err := GetResponseBodyPlatform(url+"/rpc", params)
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ]")
	}

	infoResp := rpcInfoResponse{}

	err = json.Unmarshal(body, &infoResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ] Can't unmarshal")
	}
	if infoResp.Error != nil {
		return nil, errors.New("[ Info ] Field 'error' is not nil: " + fmt.Sprint(infoResp.Error))
	}

	return &infoResp.Result, nil
}

// Status makes rpc request to info.Status method and extracts it
func Status(url string) (*StatusResponse, error) {
	params := getDefaultRPCParams("node.GetStatus")

	body, err := GetResponseBodyPlatform(url+"/rpc", params)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ]")
	}

	statusResp := rpcStatusResponse{}

	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ] Can't unmarshal")
	}
	if statusResp.Error != nil {
		return nil, errors.New("[ Status ] Field 'error' is not nil: " + fmt.Sprint(statusResp.Error))
	}

	return &statusResp.Result, nil
}

// LogOff rpc request turns network state to NoNetwork to initiate reconnect sequence.
func LogOff(url string) (*StatusResponse, error) {
	params := getDefaultRPCParams("status.LogOff")

	body, err := GetResponseBodyPlatform(url+"/rpc", params)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ]")
	}

	statusResp := rpcStatusResponse{}

	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ] Can't unmarshal")
	}
	if statusResp.Error != nil {
		return nil, errors.New("[ Status ] Field 'error' is not nil: " + fmt.Sprint(statusResp.Error))
	}

	return &statusResp.Result, nil
}
