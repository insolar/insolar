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
package updater

import (
	"github.com/insolar/insolar/updater/request"
	"github.com/insolar/insolar/updateserv"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/version"
	"github.com/stretchr/testify/assert"
)

// Just to make Goland happy
func TestStub(t *testing.T) {
	updater := NewUpdater()
	assert.NotNil(t, updater)
	assert.Equal(t, updater.CurrentVer, version.Version, "Version verify success")
	assert.Equal(t, updater.BinariesList, []string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "insupdater"})
	assert.Equal(t, updater.ServersList, []string{"http://localhost:2345"})
	assert.NotEqual(t, updater.LastSuccessServer, "http://localhost:2345")
}

func TestStubSameVersion(t *testing.T) {
	updater := NewUpdater()
	updater.ServersList = []string{""}
	assert.NotNil(t, updater)
	b, s, e := updater.IsSameVersion("v0.0.0")
	assert.Error(t, e)
	assert.Equal(t, s, "")
	assert.Equal(t, updater.CurrentVer, "v0.0.0")
	assert.Equal(t, b, true)
	assert.Equal(t, updater.DownloadFiles("v0.0.0"), false)
}

func TestHttp(t *testing.T) {
	ver := request.NewVersion("v0.3.1")

	us := updateserv.NewUpdateServer("2346", "./data")
	us.LatestVersion = "v0.3.1"
	assert.Equal(t, us.LatestVersion, "v0.3.1")

	err := us.Start()
	assert.NoError(t, err)
	assert.NotNil(t, us)
	assert.Equal(t, us.UploadPath, "./data")
	assert.Equal(t, us.Port, "2346")

	response, err := http.Get("http://localhost:2346/latest")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"latest":"v0.3.1"`)

	addr, v, e := request.ReqCurrentVer([]string{"http://localhost:2346"})
	assert.NoError(t, e)
	assert.Equal(t, addr, "http://localhost:2346")
	assert.Equal(t, v, ver)

	up := NewUpdater()
	b, s, e := up.IsSameVersion("v0.3.0")
	assert.Error(t, e)
	assert.Equal(t, s, "")
	assert.Equal(t, b, true)

	up.LastSuccessServer = "http://localhost:2346"
	b, s, e = up.IsSameVersion("v0.3.0")
	assert.NoError(t, e)
	assert.Equal(t, s, "v0.3.1")
	assert.Equal(t, b, false)

	us.Stop()
}
