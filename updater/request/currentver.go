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
	"regexp"

	"github.com/insolar/insolar/log"
)

type Version struct {
	Value    string `json:"latest"`
	Major    int    `json:"major"`
	Minor    int    `json:"minor"`
	Revision int    `json:"revision"`
}

// Create new Version object
func NewVersion(ver string) *Version {
	v := Version{}
	v.Value = ver
	re := regexp.MustCompile("[0-9]+")
	arr := re.FindAllString(ver, -1)
	v.Major = extractIntValue(arr, 0)
	v.Minor = extractIntValue(arr, 1)
	v.Revision = extractIntValue(arr, 2)
	return &v
}

// Get active update server from address list and get latest Version
func ReqCurrentVer(addresses []string) (string, *Version, error) {
	log.Debug("Found update server addresses: ", addresses)

	for _, address := range addresses {
		if address != "" {
			log.Info("Found update server address: ", address)
			ver, err := ReqCurrentVerFromAddress(GetProtocol(address), address)

			if err == nil && ver != "" {
				currentVer := ExtractVersion(ver)
				return address, currentVer, err
			}
		}
	}
	log.Warn("No Update Servers available")
	return "", nil, errors.New("No Update Servers available")
}

// Get latest Version from remote server
func ReqCurrentVerFromAddress(request UpdateNode, address string) (string, error) {
	log.Debug("Check latest version info from remote server: ", address)
	if request == nil {
		return "", errors.New("Unknown protocol")
	}
	return request.getCurrentVer(address)
}
