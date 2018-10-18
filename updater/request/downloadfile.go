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

package request

import (
	"errors"
	"path"

	"github.com/insolar/insolar/log"
)

// Start download binary files with version="version"
func DownloadFiles(version string, binariesList []string, url string) (success bool) {
	success = false
	errs := 0
	total := 0

	pathToSave, err := createCurrentPath(version)
	if err != nil {
		log.Error(err)
		return
	}
	request := GetProtocol(url)
	log.Info("Download updates from remote server: ", url)
	for _, binary := range binariesList {
		log.Info("Download file : ", binary)
		err := downloadFromAddress(request, path.Join(pathToSave, binary), url+"/"+version+"/"+binary)
		total++
		if err != nil {
			log.Error(err)
			errs++
		} else {
			log.Info("SUCCESS")
		}
	}
	log.Info("Download complete, TOTAL:", total, ", ERRORS: ", errs)
	if errs == 0 && total == len(binariesList) {
		success = true
	}
	return
}

func downloadFromAddress(request UpdateNode, filePath string, url string) error {
	if request == nil {
		return errors.New("Unknown protocol")
	}
	return request.downloadFile(filePath, url)
}
