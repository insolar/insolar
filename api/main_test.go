/*
 *    Copyright 2018 INS Ecosystem
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
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func TestLaunchApi(t *testing.T) {
	cfg := configuration.NewApiRunner()
	api, err := NewApiRunner(cfg)
	assert.NoError(t, err)

	cs := core.Components{}
	err = api.Start(cs)
	assert.NoError(t, err)

	resp, err := http.Get("http://localhost:8080/api/v1?query_type=LOL")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"message": "Wrong query parameter 'query_type' = 'LOL'"`)

	api.Stop()
	assert.NoError(t, err)
}
