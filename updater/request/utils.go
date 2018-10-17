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
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// Verify current protocol version from URL, if http://... -> HTTPUpdateNode
func GetProtocol(address string) UpdateNode {
	protocol, err := getProtocolFromAddress(address)
	if err != nil {
		return nil
	}
	if protocol != "" {
		switch protocol {
		case "http":
			{
				return HTTPUpdateNode{}
			}
		default:
			{
				log.Warn("Unknown protocol ", protocol[0])
				return nil
			}
		}
	}
	return nil
}

func getProtocolFromAddress(urlString string) (string, error) {
	if !strings.Contains(urlString, "://") {
		return "", errors.New("Protocol is not set, use 'http://' before address")
	}
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	return u.Scheme, nil
}

func createCurrentPath(version string) string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	pathToSave := path.Join(pwd, version)
	if err := os.Mkdir(pathToSave, 0750); err != nil {
		log.Warn("Error while create folder: ", err)
	}
	return pathToSave
}

// Unmarshal JSON Version
func ExtractVersion(ver string) *Version {
	latestVersion := Version{}
	err := json.Unmarshal([]byte(ver), &latestVersion)
	if err != nil {
		log.Warn("Error parsing data: ", err)
		return nil
	}
	return &latestVersion
}

// Compare two Version objects
// if ver1 < ver2   return -1
// if ver1 == ver2  return 0
// if ver1 > ver2   return 1
func CompareVersion(ver1 *Version, ver2 *Version) (result int) {
	result = 0
	if result = compare(ver1.Major, ver2.Major); result == 0 {
		if result = compare(ver1.Minor, ver2.Minor); result == 0 {
			result = compare(ver1.Revision, ver2.Revision)
		}
	}
	return
}

// Compare two Version objects and return MAX
func GetMaxVersion(ver1 *Version, ver2 *Version) *Version {
	if ver1 == nil {
		return ver2
	}
	if ver2 == nil {
		return ver1
	}
	resultCompare := CompareVersion(ver1, ver2)
	if resultCompare == 1 {
		return ver1
	}
	return ver2
}

func compare(x int, y int) int {
	if x < y {
		return -1
	} else if x > y {
		return 1
	} else {
		return 0
	}
}

func extractIntValue(arr []string, index int) int {
	if len(arr) >= index+1 && arr[index] != "" {
		value, err := strconv.Atoi(arr[index])
		if err == nil {
			return value
		}
	}
	return 0
}
