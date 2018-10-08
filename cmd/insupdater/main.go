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
	"github.com/chzyer/readline"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	upd "github.com/insolar/insolar/updater"
	"github.com/insolar/insolar/version"
	jww "github.com/spf13/jwalterweatherman"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	initLogger(cfgHolder.Configuration.Log)
	fmt.Println("Updater module (Check for update and run insolar peer)")
	fmt.Println("Version: ", version.GetFullVersion())

	// send version.Version to update server
	// if version != lastVersionResult in the Update Server => download and change files
	updater := upd.NewUpdater()
	verifyAndUpdate(updater)
	service(updater)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debug("caught sig: ", sig)
		os.Exit(0)
	}()

	interactive()
}

func verifyAndUpdate(updater *upd.Updater) {
	log.Info("Start verify for update ")
	sameVersion, newVersion, err := updater.IsSameVersion(version.Version)
	if err != nil {
		onErr("Error at the Update Server access stage: ", err)
	}
	if !sameVersion {
		log.Debug("Current version: ", version.Version, ", found version: ", newVersion)
		// Run Update
		success := updater.DownloadFiles(newVersion)
		if success {
			// ToDo: send stop signal
			// ToDo: copy files from folder=./${VERSION} to current folder
		}
	} else {
		log.Info("Already updated!")
	}
	// Run peer
	executePeer()
	// ToDo: Run update service with timer
	// exit
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Errorln(err.Error())
	}
}

func onErr(text string, err error) {
	fmt.Println(text, err)
	os.Exit(1)
}

func executePeer() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Warn(err)
		pwd = "."
	}
	out, err := exec.Command(path.Join(pwd, "insolard"), "--config", path.Join(pwd, "..", "scripts", "insolard", "insolar.yaml")).CombinedOutput()
	if err != nil {
		log.Warn("Cannot run insolar deamon, verify PATH to file 'insolard'")
	} else {
		log.Info(out)
	}
}

func service(updater *upd.Updater) {
	delay := time.Duration(updater.Delay)
	go func() {
		ticker := time.NewTicker(time.Second * delay)
		defer func() {
			log.Info("Stopping update service")
			ticker.Stop()
		}()

		for range ticker.C {
			verifyAndUpdate(updater)
		}
	}()
}

func interactive() {

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer func() {
		errRlClose := rl.Close()
		if errRlClose != nil {
			panic(errRlClose)
		}
	}()
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		input := strings.Split(line, " ")

		switch input[0] {
		case "exit":
			fallthrough
		case "quit":
			{
				return
			}
		case "version":
			{
				fmt.Println(version.GetFullVersion())
			}
		}

	}

}
