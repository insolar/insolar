//
// Copyright 2019 Insolar Technologies GmbH
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
//

package main

import (
	"fmt"
	"os"

	"github.com/insolar/insolar/api/sdk"
	"github.com/insolar/insolar/log"
	"github.com/spf13/pflag"
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

	insSDK, err := sdk.NewSDK([]string{apiAdminURL}, []string{apiPublicURL}, memberKeys)
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
