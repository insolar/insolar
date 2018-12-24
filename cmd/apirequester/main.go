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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/insolar/insolar/log"
	"github.com/spf13/pflag"
)

const defaultURL = "http://localhost:19191/api"

var (
	rootmemberkeys string
	apiurl         string

	rootMember memberInfo
)

func parseInputParams() {
	pflag.StringVarP(&rootmemberkeys, "rootmemberkeys", "k", "", "path to file with RootMember keys")
	pflag.StringVarP(&apiurl, "url", "u", defaultURL, "api url")
	pflag.Parse()
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

type memberInfo struct {
	ref        string
	privateKey string
	traceId    string
}

type memberKeys struct {
	Private string `json:"private_key"`
	Public  string `json:"public_key"`
}

func getRootMemberRef() string {
	infoResp, err := info()
	check("Can not get info:", err)
	return infoResp.RootMember
}

func getRootMemberInfo(fileName string) memberInfo {

	rawConf, err := ioutil.ReadFile(fileName)
	check("problem with reading root member keys file", err)

	keys := memberKeys{}
	err = json.Unmarshal(rawConf, &keys)
	check("problem with unmarshaling root member keys", err)

	return memberInfo{getRootMemberRef(), keys.Private, ""}
}

func main() {
	parseInputParams()

	err := log.SetLevel("error")
	check("can not set 'error' level on logger:", err)
	rootMember = getRootMemberInfo(rootmemberkeys)

	// you can modify this manual tests by commenting any of this functions or/and add some new functions if necessary

	// make one request to create new member
	oneSimpleRequest()

	// make several (10) requests to create new member (every request make call to RootMember instance)
	severalSimpleRequestToRootMember()

	// make several (10) requests to transfer money (every request make call to different members instances)
	severalSimpleRequestToDifferentMembers()

	// make several (10) requests in parallel to create new member (every request make call to RootMember instance)
	severalParallelRequestToRootMember()

	// make several (10) requests in parallel to transfer money (every request make call to different members instances)
	severalParallelRequestToDifferentMembers()
}
