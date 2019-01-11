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
	"github.com/spf13/pflag"
)

const defaultURL = "http://localhost:19191/"

var (
	role        string
	apiHost     string
	keysFileOut string
	certFileOut string
)

func parseInputParams() {
	pflag.StringVarP(&role, "role", "r", "virtual", "The role of the new node")
	pflag.StringVarP(&apiHost, "api_host", "h", defaultURL, "HTTP base host that serves insolar API requests")
	pflag.StringVarP(&keysFileOut, "keys_file", "k", "keys.json", "The OUT file for public/private keys of the node")
	pflag.StringVarP(&certFileOut, "cert_file", "c", "cert.json", "The OUT file the node certificate")
	pflag.Parse()
}

func main() {
	parseInputParams()
}
