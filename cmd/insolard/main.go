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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/cmd/insolard/componentmanager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updater"
	"github.com/insolar/insolar/version"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var (
	configPath string
)

func parseInputParams() {
	var rootCmd = &cobra.Command{Use: "insolard"}
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Wrong input params:", err)
	}
}

func main() {
	parseInputParams()

	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	var err error
	if len(configPath) != 0 {
		err = cfgHolder.LoadFromFile(configPath)
	} else {
		err = cfgHolder.Load()
	}
	if err != nil {
		log.Warnln("failed to load configuration from file: ", err.Error())
	}

	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("failed to load configuration from env:", err.Error())
	}
	initLogger(cfgHolder.Configuration.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))
	cm, nw := componentmanager.New(cfgHolder)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debugln("caught sig: ", sig)

		cm.StopAll()
		os.Exit(0)
	}()

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Running interactive mode:")

	go func() {
		runUpdateService(cm.Components.Updater.(*updater.Updater))
	}()
	repl(nw)
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(cfg.Level)
	if err != nil {
		log.Errorln(err.Error())
	}
}
