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
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

type HTTPUpdateNode struct {
	UpdateNode
}

func (request HTTPUpdateNode) getCurrentVer(address string) (string, error) {
	response, err := http.Get(address + "/latest")
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (request HTTPUpdateNode) downloadFile(filepath string, url string) error {

	//Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Error("OS Create file error: ", err)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Error("HTTP server error: ", err)
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Error("HTTP bad status: ", resp.Status)
		return errors.Errorf("HTTP error: %s ", resp.Status)
	}

	// Writer the body to file
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Error("OS write file error: ", err)
		return err
	}
	log.Info("Downloaded file: "+url+", save to: "+filepath+", total bytes: ", written)
	return nil
}
