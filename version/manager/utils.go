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

package manager

import (
	"github.com/blang/semver"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
)

func ProcessVersionConsensus(nodes []core.Node) error {
	if len(nodes) == 0 {
		return errors.New("List of nodes is empty")
	}
	mapOfVersions := getMapOfVersion(nodes)
	topVersion, err := getMaxVersion(getRequired(len(nodes)), mapOfVersions)
	if err != nil {
		return err
	}
	vm, err := GetVersionManager()
	if err != nil {
		return err
	}
	currentVersion, err := ParseVersion(version.Version)
	if err != nil {
		return errors.New("current version is invalid: " + err.Error())
	}
	if currentVersion.Compare(*topVersion) != 0 {
		log.Warn("WARNING! Current version: " + StringVersion(currentVersion) + ", must go to version: " + StringVersion(topVersion))
	}
	vm.AgreedVersion = topVersion
	return nil
}

func getMapOfVersion(nodes []core.Node) *map[string]int {
	mapOfVersions := make(map[string]int)

	for _, node := range nodes {
		ver := node.Version()
		if _, ok := mapOfVersions[ver]; ok {
			mapOfVersions[ver]++
		} else {
			mapOfVersions[ver] = 1
		}
	}
	return &mapOfVersions
}

func getMaxVersion(required int, mapOfVersions *map[string]int) (*semver.Version, error) {
	versions := make([]semver.Version, 0, len(*mapOfVersions))
	for version := range *mapOfVersions {
		semVer, err := ParseVersion(version)
		if err != nil {
			log.Warn("Error parsing version: " + version)
		}
		if semVer != nil {
			versions = append(versions, *semVer)
		}
	}
	semver.Sort(versions)
	count := 0
	for i := len(versions) - 1; i >= 0; i-- {
		key := StringVersion(&versions[i])
		count += (*mapOfVersions)[key]
		if count >= required {
			return &versions[i], nil
		}
	}
	return nil, errors.New("Version consensus is not reached")
}

func Verify(key string) bool {
	vm, err := GetVersionManager()
	if err != nil {
		return false
	}
	return vm.IsAvailable(key)
}

func getRequired(count int) int {
	return count/2 + 1
}

func ParseVersion(ver string) (*semver.Version, error) {
	if ver == "unset" {
		return semver.New("0.0.0")
	}
	version, err := semver.ParseTolerant(ver)
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func StringVersion(ver *semver.Version) string {
	return "v" + ver.String()
}
