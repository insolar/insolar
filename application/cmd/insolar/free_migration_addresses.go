// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/sdk"
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
