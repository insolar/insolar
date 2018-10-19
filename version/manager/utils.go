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
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

func ProcessVersionConsensus(nodes []*core.ActiveNode) error {
	if len(nodes) == 0 {
		return errors.New("List of nodes is empty")
	}

	mapOfVersions := getMapOfVersion(nodes)
	topVersion := getMaxVersion(getRequired(len(nodes)), mapOfVersions)
	if topVersion != nil {
		GetVM().AgreedVersion = *topVersion
	}
	return nil
}

func getMapOfVersion(nodes []*core.ActiveNode) *map[string]int {
	mapOfVersions := make(map[string]int)

	// ToDo: I Need a version from the ActiveNodeList
	for _ = range nodes { // node
		ver := "v0.5.0" // node.Version
		if _, ok := mapOfVersions[ver]; ok {
			mapOfVersions[ver]++
		} else {
			mapOfVersions[ver] = 1
		}
	}
	return &mapOfVersions
}

func getMaxVersion(required int, mapOfVersions *map[string]int) *string {
	for key, count := range *mapOfVersions {
		if count >= required {
			return &key
		}
	}
	return nil
}

func Verify(key string) bool {
	vm := GetVM()
	return vm.Verify(key)
}

func getRequired(count int) int {
	return count/2 + 1
}
