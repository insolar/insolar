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

package updateserv

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Just to make Goland happy
func TestStub(t *testing.T) {
	us := NewUpdateServer("2345", "../")
	us.LatestVersion = "v0.3.1"
	assert.Equal(t, us.LatestVersion, "v0.3.1")

	err := us.Start()
	assert.NoError(t, err)
	assert.NotNil(t, us)
	assert.Equal(t, us.UploadPath, "../")
	assert.Equal(t, us.Port, "2345")

	response, err := http.Get("http://localhost:2345/latest")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"latest":"v0.3.1"`)
	us.Stop()
}
