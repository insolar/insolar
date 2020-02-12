// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/api/sdk"
)

var shardsAtOneTime = 10

func getfreeMigrationCount(adminUrls []string, publicUrls []string, memberKeysDirPath string, shardsCount int, alertCount int) {
	insSDK, err := sdk.NewSDK(adminUrls, publicUrls, memberKeysDirPath, sdk.DefaultOptions)
	check("SDK is not initialized: ", err)

	shoudAlert := map[int]int{}
	freeAdressesInShards := map[int]int{}

	for i := 0; i < shardsCount; i += shardsAtOneTime {
		part, _, err := insSDK.GetAddressCount(i)
		check(fmt.Sprintf("Error while getting addresses from shards %d to %d: ", i, i+shardsAtOneTime), err)
		partSliced, ok := part.([]interface{})
		if !ok {
			check(
				fmt.Sprintf("Error while getting addresses from shards %d to %d: ", i, i+shardsAtOneTime),
				errors.New("error while converting result"),
			)
		}

		for _, r := range partSliced {
			rMap := r.(map[string]interface{})
			s, ok := rMap["shardIndex"].(float64)
			if !ok {
				check(
					fmt.Sprintf("Error while getting addresses from shards %d to %d: ", i, i+shardsAtOneTime),
					errors.New("error while converting shardIndex"),
				)
			}
			shardIndex := int(s)

			f, ok := rMap["freeCount"].(float64)
			if !ok {
				check(
					fmt.Sprintf("Error while getting addresses from shard %d: ", shardIndex),
					errors.New("error while converting freeCount"),
				)
			}
			freeCount := int(f)

			freeAdressesInShards[shardIndex] = freeCount
			if freeCount <= alertCount {
				shoudAlert[shardIndex] = freeCount
				fmt.Printf("Shard: %4d. Count: %d\n", shardIndex, freeCount)
			}
		}
	}
	if len(shoudAlert) > 0 {
		fmt.Println("ALERT: too little free addresses in shards!")
	} else {
		fmt.Println("All shards have enough free addresses")
	}
}
