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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/updater/request"
	"github.com/insolar/insolar/updateserv"
	"github.com/insolar/insolar/version"
	"github.com/stretchr/testify/assert"
)

func TestStubSameVersion(t *testing.T) {
	cfg := configuration.NewUpdater()
	updater, err := NewUpdater(&cfg)
	assert.NoError(t, err)
	err = updater.Start(core.Components{})
	assert.NoError(t, err)
	updater.ServersList = []string{""}
	assert.NotNil(t, updater)
	b, s, e := updater.IsSameVersion("v0.0.0")
	assert.Error(t, e)
	assert.Equal(t, s, "")
	assert.Equal(t, updater.CurrentVer, "v0.0.0")
	assert.Equal(t, b, true)
	assert.Equal(t, updater.DownloadFiles("v0.0.0"), false)
	err = updater.Stop()
	assert.NoError(t, err)
	RemoveBinariesFolder("v0.0.0")
}

func RemoveBinariesFolder(version string) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	pathToSave := path.Join(pwd, version)
	if err := os.RemoveAll(pathToSave); err != nil {
		fmt.Println(err)
	}
}

func TestHttp(t *testing.T) {
	us := updateserv.NewUpdateServer("2345", "../")
	us.LatestVersion = "v0.3.1"
	assert.Equal(t, us.LatestVersion, "v0.3.1")

	err := us.Start()
	assert.NoError(t, err)
	assert.NotNil(t, us)
	assert.Equal(t, us.UploadPath, "../")
	assert.Equal(t, us.Port, "2345")

	response, err := http.Get("http://localhost:2345/latest")
	defer response.Body.Close()
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"latest":"v0.3.1"`)
	ver, err := request.NewVersion("v0.3.1")
	assert.NoError(t, err)
	us.LatestVersion = "v0.3.1"
	assert.Equal(t, us.LatestVersion, "v0.3.1")

	addr, v, e := request.ReqCurrentVer([]string{"http://localhost:2345"})
	assert.NoError(t, e)
	assert.Equal(t, addr, "http://localhost:2345")
	assert.Equal(t, v, ver)

	cfg := configuration.NewUpdater()
	updater, err := NewUpdater(&cfg)
	assert.NoError(t, err)
	err = updater.Start(core.Components{})
	assert.NoError(t, err)

	assert.Equal(t, updater.CurrentVer, version.Version)
	assert.Equal(t, updater.BinariesList, []string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "insupdater"})
	assert.NotEqual(t, updater.ServersList, []string{""})
	assert.Equal(t, updater.LastSuccessServer, "")

	if version.Version == "unset" {
		err = updater.verifyAndUpdate()
		assert.Error(t, err)
		updater.CurrentVer = "v0.3.1"
		err = updater.verifyAndUpdate()
		assert.Error(t, err)
	} else {
		err = updater.verifyAndUpdate()
		assert.NoError(t, err)
		updater.CurrentVer = "v0.3.1"
		err = updater.verifyAndUpdate()
		assert.NoError(t, err)
	}
	updater.ServersList = []string{}
	updater.LastSuccessServer = ""
	err = updater.verifyAndUpdate()
	assert.Error(t, err)

	b, s, e := updater.IsSameVersion("v0.3.0")
	assert.Error(t, e)
	assert.Equal(t, s, "")
	assert.Equal(t, b, true)

	updater.LastSuccessServer = "http://localhost:2345"
	b, s, e = updater.IsSameVersion("v0.3.0")
	assert.NoError(t, e)
	assert.Equal(t, s, "v0.3.1")
	assert.Equal(t, b, false)

	err = updater.Stop()
	assert.NoError(t, err)
	err = us.Stop()
	assert.NoError(t, err)
	RemoveBinariesFolder("v0.3.1")
}
