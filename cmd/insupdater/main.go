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
	"strings"
	"syscall"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	upd "github.com/insolar/insolar/updater"
	"github.com/insolar/insolar/version"
	jww "github.com/spf13/jwalterweatherman"
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
	err := verifyAndUpdate(updater)
	if err != nil {
		log.Warn(err)
	}
	service(updater)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debug("caught sig: ", sig)
		os.Exit(0)
	}()

}

func verifyAndUpdate(updater *upd.Updater) error {
	log.Info("Start verify for update ")
	sameVersion, newVersion, err := updater.IsSameVersion(version.Version)
	if err != nil {
		return err
	}
	if !sameVersion {
		log.Debug("Current version: ", version.Version, ", found version: ", newVersion)
		// Run Update
		//if updater.DownloadFiles(newVersion) {
		//	// ToDo: send stop signal, then copy files from folder=./${VERSION} to current folder
		//}
	}
	// Run peer
	//executePeer()
	// ToDo: Run update service with timer
	// exit
	return nil
}

func initLogger(cfg configuration.Log) (err error) {
	err = log.SetLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Errorln(err.Error())
	}
	return
}

//func onErr(text string, err error) {
//	fmt.Println(text, err)
//	os.Exit(1)
//}

//func executePeer() {
//	pwd, err := os.Getwd()
//	if err != nil {
//		log.Warn(err)
//		pwd = "."
//	}
//	out, err := exec.Command(path.Join(pwd, "insolard"), "--config", path.Join(pwd, "..", "scripts", "insolard", "insolar.yaml")).CombinedOutput()
//	if err != nil {
//		log.Warn("Cannot run insolar deamon, verify PATH to file 'insolard'")
//	} else {
//		log.Info(out)
//	}
//}

func service(updater *upd.Updater) {
	delay := time.Duration(updater.Delay)
	go func() {
		ticker := time.NewTicker(time.Second * delay)
		defer func() {
			log.Info("Stopping update service")
			ticker.Stop()
		}()

		for range ticker.C {
			err := verifyAndUpdate(updater)
			if err != nil {
				log.Warn(err)
			}
		}
	}()
}
