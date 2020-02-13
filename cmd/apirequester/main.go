// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/insolar/insolar/application/api/sdk"
	"github.com/insolar/insolar/log"
)

const defaultAdminURL = "http://localhost:19001/admin-api/rpc"
const defaultPublicURL = "http://localhost:19101/api/rpc"

var (
	memberKeys   string
	apiAdminURL  string
	apiPublicURL string
)

func parseInputParams() {
	pflag.StringVarP(&memberKeys, "memberkeys", "k", "", "path to dir with members keys")
	pflag.StringVarP(&apiAdminURL, "adminurls", "a", defaultAdminURL, "admin api url")
	pflag.StringVarP(&apiPublicURL, "publicurls", "p", defaultPublicURL, "public api url")
	pflag.Parse()
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	parseInputParams()

	err := log.SetLevel("error")
	check("can't set 'error' level on logger: ", err)

	insSDK, err := sdk.NewSDK([]string{apiAdminURL}, []string{apiPublicURL}, memberKeys, sdk.DefaultOptions)
	check("can't create SDK: ", err)

	// you can modify this manual tests by commenting any of this functions or/and add some new functions if necessary

	// make one request to create new member
	oneSimpleRequest(insSDK)

	// make several (10) requests to create new member (every request make call to RootMember instance)
	severalSimpleRequestToRootMember(insSDK)

	// make several (10) requests to transfer money (every request make call to different members instances)
	severalSimpleRequestToDifferentMembers(insSDK)

	// make several (10) requests in parallel to create new member (every request make call to RootMember instance)
	severalParallelRequestToRootMember(insSDK)

	// make several (10) requests in parallel to transfer money (every request make call to different members instances)
	severalParallelRequestToDifferentMembers(insSDK)
}
