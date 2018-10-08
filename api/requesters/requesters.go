/*
 *    Copyright 2018 Insolar
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

package requesters

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type PostParams = map[string]string

func GetResponseBody(url string, postP map[string]string) ([]byte, error) {
	jsonValue, _ := json.Marshal(postP)
	postResp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponseBody ] Problem with sending request")
	}
	if http.StatusOK != postResp.StatusCode {
		return nil, errors.Wrap(err, "[ getResponseBody ] Bad http response code: "+strconv.Itoa(postResp.StatusCode))
	}

	body, err := ioutil.ReadAll(postResp.Body)
	defer postResp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponseBody ] Problem with reading body")
	}

	return body, nil
}

func GetSeed(url string) (string, error) {
	body, err := GetResponseBody(url, PostParams{
		"query_type": "get_seed",
	})
	if err != nil {
		return "", errors.Wrap(err, "[ getSeed ]")
	}

	type seedResponse struct{ Seed string }
	seedResp := seedResponse{}

	err = json.Unmarshal(body, &seedResp)
	if err != nil {
		return "", errors.Wrap(err, "[ getSeed ]")
	}

	return seedResp.Seed, nil
}
