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
	"os"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater/request"
	"github.com/insolar/insolar/version"
)

func (up *Updater) verifyAndUpdate() error {
	log.Info("Try verify for update ")
	sameVersion, newVersion, err := up.IsSameVersion(version.Version)
	if err != nil {
		return err
	}
	if !sameVersion {
		log.Debug("Current version: ", version.Version, ", found version: ", newVersion)
		// Run Update
		if up.DownloadFiles(newVersion) {
			// ToDo: send stop signal, then copy files from folder=./${VERSION} to current folder

			err := os.Setenv("INS_LATEST_VER", newVersion)
			if err != nil {
				log.Warn("Can not set OS envelop value INS_LATEST_VER: ", err)
			}
		}
	}
	// Run peer
	//executePeer()
	// ToDo: Run update service with timer
	// exit
	return nil
}

func (up *Updater) IsSameVersion(currentVersion string) (bool, string, error) {
	log.Debug("Verify latest peer version from remote server")
	up.CurrentVer = currentVersion
	currentVer := request.NewVersion(currentVersion)
	if up.LastSuccessServer != "" {
		log.Debug("Latest update server was: ", up.LastSuccessServer)
		vers, err := request.ReqCurrentVerFromAddress(request.GetProtocol(up.LastSuccessServer), up.LastSuccessServer)
		if err == nil && vers != "" {
			versionFromUS := request.ExtractVersion(vers)
			return request.CompareVersion(versionFromUS, currentVer) < 0, versionFromUS.Value, nil
		}
	}
	lastSuccessServer, versionFromUS, err := request.ReqCurrentVer(up.ServersList)
	if err != nil {
		return true, "", err
	}
	log.Debug("Get version=", versionFromUS.Value, " from remote server: ", lastSuccessServer)
	up.LastSuccessServer = lastSuccessServer

	if versionFromUS == nil || up.CurrentVer == "" {
		return true, "unset", nil
	} else
	//if(updater.currentVer != versionFromUS){
	if request.CompareVersion(versionFromUS, currentVer) > 0 {
		return false, versionFromUS.Value, nil
	}
	return true, versionFromUS.Value, nil
}

func (up *Updater) DownloadFiles(version string) (success bool) {
	if up.started {
		return false
	}
	log.Info("Start download files from remote server")
	if up.LastSuccessServer == "" {
		return false
	}
	up.started = true
	success = request.DownloadFiles(version, up.BinariesList, up.LastSuccessServer)
	up.started = false
	return
}
