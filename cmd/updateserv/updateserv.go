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
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/updateserv"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	err := initLogger(cfgHolder.Configuration.Log)
	if err != nil {
		log.Warn("Can't initialize logger: ", err)
	}
	//port = pflag.StringP("port", "p", port, "port to listen")

	server := updateserv.NewUpdateServer(getPort(), getUploadPath())
	latestVersion := os.Getenv("BUILD_VERSION")
	server.LatestVersion = latestVersion
	err = server.Start()
	if err != nil {
		log.Warn("Can't start server: ", err)
		os.Exit(1)
	}
	defer func() {
		err := server.Stop()
		if err != nil {
			log.Warn(err)
		}
	}()
	repl()
}

func initLogger(cfg configuration.Log) (err error) {
	return log.SetLevel(cfg.Level)
}

func getPort() (port string) {
	port = "2345"
	if portValue := os.Getenv("updateserver_port"); portValue != "" {
		port = portValue
	}
	return
}

func getUploadPath() (uploadPath string) {
	uploadPath = "./data"
	if uploadPathValue := os.Getenv("upload_path"); uploadPathValue != "" {
		uploadPath = uploadPathValue
	}
	return
}

func repl() {
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
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		input := strings.Split(line, " ")

		switch input[0] {
		case "exit":
			fallthrough
		case "quit":
			return
		}
	}
}
