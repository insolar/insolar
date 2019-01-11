/*
 *    Copyright 2019 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package inssdk

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/insolar/insolar/api/requester"
	"github.com/pkg/errors"
)

type ringBuffer struct {
	sync.Mutex
	urls   []string
	cursor int
}

func (rb *ringBuffer) Next() string {
	rb.Lock()
	defer rb.Unlock()
	rb.cursor++
	if rb.cursor >= len(rb.urls) {
		rb.cursor = 0
	}
	return rb.urls[rb.cursor]
}

type response struct {
	Error   string
	Result  interface{}
	TraceID string
}

type SDK struct {
	apiUrls *ringBuffer
}

func NewSDK(urls []string) (*SDK, error) {

}

func (sdk *SDK) sendRequest(ctx context.Context, method string, params []interface{}, userCfg *requester.UserConfigJSON) ([]byte, error) {
	reqCfg := &requester.RequestConfigJSON{
		Params: params,
		Method: method,
	}

	body, err := requester.Send(ctx, sdk.apiUrls.Next(), userCfg, reqCfg)
	if err != nil {
		errors.Wrap(err, "can not send request")
	}

	return body, nil
}

func (sdk *SDK) Transfer(ctx context.Context, amount uint, from memberInfo, to memberInfo) error {
	params := []interface{}{amount, to.ref}
	body, err := sendRequest(ctx, "Transfer", params, from)
	if err != nil {
		return err
	}
	transferResponse := getResponse(body)

	if transferResponse.Error != "" {
		return errors.New(transferResponse.Error)
	}

	return nil
}

func (sdk *SDK) CreateMember() (*memberInfo, error) {

}

func (sdk *SDK) getResponse(body []byte) (*response, error) {
	res := &response{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "Problems with unmarshal response: ")
	}

	return res, nil
}

// Constructor
func (sdk *SDK) Info() (*requester.InfoResponse, error) {
	info, err := requester.Info(apiurls.Next())
	check("problem with request to info:", err)
	return info, nil
}
