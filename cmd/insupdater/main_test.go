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

	"github.com/insolar/insolar/configuration"
	upd "github.com/insolar/insolar/updater"
	"github.com/insolar/insolar/updateserv"
	"github.com/insolar/insolar/version"
	"github.com/stretchr/testify/assert"
)

// Just to make Goland happy
func TestStub(t *testing.T) {
	cfgHolder := configuration.NewHolder()
	cfgHolder.Configuration.Log.Level = "wewqewq"
	assert.Error(t, initLogger(cfgHolder.Configuration.Log))
	cfgHolder = configuration.NewHolder()
	assert.NoError(t, initLogger(cfgHolder.Configuration.Log))

	us := updateserv.NewUpdateServer("2345", "../")
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

	updater := upd.NewUpdater()
	assert.NotNil(t, updater)
	assert.Equal(t, updater.CurrentVer, version.Version)
	assert.Equal(t, updater.BinariesList, []string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "insupdater"})
	assert.NotEqual(t, updater.ServersList, []string{""})
	assert.Equal(t, updater.LastSuccessServer, "")
	err = verifyAndUpdate(updater)
	assert.NoError(t, err)
	updater.CurrentVer = "v0.3.1"
	err = verifyAndUpdate(updater)
	assert.NoError(t, err)
	service(updater)

	updater.ServersList = []string{}
	updater.LastSuccessServer = ""
	err = verifyAndUpdate(updater)
	assert.Error(t, err)
	us.Stop()
}
