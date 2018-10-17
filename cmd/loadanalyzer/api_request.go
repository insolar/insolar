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

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

const TestURL = "http://localhost:19191/api/v1"

type postParams map[string]interface{}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type baseResponse struct {
	Qid string         `json:"qid"`
	Err *errorResponse `json:"error"`
}

type createMemberResponse struct {
	baseResponse
	Reference string `json:"reference"`
}

func getResponseBody(postParams map[string]interface{}) []byte {
	jsonValue, err := json.Marshal(postParams)
	check("Problems with marshal request:", err)
	postResp, err := http.Post(TestURL, "application/json", bytes.NewBuffer(jsonValue))
	check("Problems with post:", err)
	body, err := ioutil.ReadAll(postResp.Body)
	check("Problems with reading from response body:", err)
	return body
}

func transfer(amount int, from string, to string) string {
	body := getResponseBody(postParams{
		"query_type": "send_money",
		"from":       from,
		"to":         to,
		"amount":     amount,
	})

	response := &baseResponse{}
	err := json.Unmarshal(body, &response)
	check("Problems with unmarshal response:", err)
	if response.Err != nil {
		return response.Err.Message
	}
	return "success"
}

func createMembers(concurrent int, repetitions int) ([]string, error) {
	var members []string
	for i := 0; i < concurrent*repetitions*2; i++ {
		body := getResponseBody(postParams{
			"query_type": "create_member",
			"name":       testutils.RandomString(),
			"public_key": "000",
		})

		memberResponse := &createMemberResponse{}
		err := json.Unmarshal(body, &memberResponse)
		check("Problems with unmarshal response:", err)

		if memberResponse.Err != nil {
			return nil, errors.New(memberResponse.Err.Message)
		}
		firstMemberRef := memberResponse.Reference
		members = append(members, firstMemberRef)
	}
	return members, nil
}
