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
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestUrl = "http://localhost:2345/latest"
)

func TestRequest(t *testing.T) {
	resp, err := http.Get(TestUrl)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotNil(t, body)
}

// Just to make Goland happy
func TestStub(t *testing.T) {
	us := newUpdateServer()
	assert.NotNil(t, us)
	assert.Equal(t, us.uploadPath, "./data")
	assert.Equal(t, us.port, "2345")
	ver := us.getLatestVersion()
	handler := us.versionHandler(ver)
	assert.NotNil(t, handler)
	assert.Nil(t, ver)
}
