/*
 *    Copyright 2018 INS Ecosystem
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
	"fmt"
	"github.com/insolar/insolar/configuration"
	log "github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
)

func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	err := cfgHolder.Load()
	if err != nil {
		log.Warnln("Failed to load configuration from file: ", err.Error())
	}

	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("Failed to load configuration from env:", err.Error())
	}

	initLogger(cfgHolder.Configuration.Log)
	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))
}

func initLogger(cfg configuration.Log) {

	cfg.Level = "debug"
	level, err := log.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Warnln(err.Error())
	}
	jww.SetLogOutput(log.StandardLogger().Out)
	log.SetLevel(level)
}
