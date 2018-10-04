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

package functest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func askApi(t *testing.T) *getSeedResponse {
	body := getResponseBody(t, postParams{
		"query_type": "get_seed",
	})

	getSeedResponse := &getSeedResponse{}
	unmarshalResponse(t, body, getSeedResponse)

	return getSeedResponse
}

func TestGetSeed(t *testing.T) {
	resp1 := askApi(t)
	resp2 := askApi(t)

	assert.NotEqual(t, resp1.Seed, resp2.Seed)

}
