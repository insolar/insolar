/*
 *    Copyright 2019 Insolar Technologies
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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/insolar/insolar/core"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var client http.Client

const (
	esc       = "\x1b%s"
	moveUp    = "[%dA"
	clearDown = "[0J"
)

const (
	insolarReady    = "Ready"
	insolarNotReady = "Not Ready"
)

func escape(format string, args ...interface{}) string {
	return fmt.Sprintf(esc, fmt.Sprintf(format, args...))
}

func moveBack(reader io.Reader) {
	fileScanner := bufio.NewScanner(reader)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}

	fmt.Print(escape(moveUp, lineCount))
	fmt.Print(escape(clearDown))
}

func main() {
	var configFile string
	pflag.StringVarP(&configFile, "config", "c", "", "config file")
	pflag.Parse()

	conf, err := pulsewatcher.ReadConfig(configFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "couldn't load config file"))
	}
	if len(conf.Nodes) == 0 {
		log.Fatal("couldn't find any nodes in config file")
	}
	if conf.Interval == 0 {
		conf.Interval = 100 * time.Millisecond
	}

	client = http.Client{
		Transport: &http.Transport{},
		Timeout:   conf.Timeout,
	}

	buffer := &bytes.Buffer{}

	var (
		state   bool
		errored int
	)

	fmt.Print("\n\n")

	for {
		state = true
		errored = 0
		results := make([][]string, len(conf.Nodes))
		lock := &sync.Mutex{}

		wg := &sync.WaitGroup{}
		wg.Add(len(conf.Nodes))
		for i, url := range conf.Nodes {
			go func(url string, i int) {
				res, err := client.Post("http://"+url+"/api/rpc", "application/json",
					strings.NewReader(`{"jsonrpc": "2.0", "method": "status.Get", "id": 0}`))
				if err != nil {
					lock.Lock()
					results[i] = []string{url, "", "", "", "", "", "", "", err.Error()}
					errored++
					lock.Unlock()
					wg.Done()
					return
				}
				defer res.Body.Close()
				data, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
				}
				var out struct {
					Result struct {
						PulseNumber         uint32
						NetworkState        string
						NodeState           string
						AdditionalNodeState string
						Origin              struct {
							Role string
						}
						ActiveListSize  int
						WorkingListSize int
					}
				}
				err = json.Unmarshal(data, &out)
				if err != nil {
					fmt.Println(string(data))
					log.Fatal(err)
				}
				lock.Lock()
				results[i] = []string{
					url,
					out.Result.NetworkState,
					out.Result.NodeState,
					out.Result.AdditionalNodeState,
					strconv.Itoa(int(out.Result.PulseNumber)),
					strconv.Itoa(out.Result.ActiveListSize),
					strconv.Itoa(out.Result.WorkingListSize),
					out.Result.Origin.Role,
					"",
				}
				state = state && out.Result.NetworkState == core.CompleteNetworkState.String() && out.Result.NodeState == core.ReadyNodeNetworkState.String()
				lock.Unlock()
				wg.Done()
			}(url, i)
		}
		wg.Wait()

		table := tablewriter.NewWriter(buffer)
		table.SetHeader([]string{
			"URL",
			"Network State",
			"Node State",
			"Additional Node State",
			"Pulse Number",
			"Active List Size",
			"Working List Size",
			"Role",
			"Error",
		})
		table.SetBorder(false)

		table.ClearRows()
		table.ClearFooter()

		moveBack(buffer)
		buffer.Reset()

		stateString := insolarReady
		color := tablewriter.FgHiGreenColor
		if !state || errored == len(conf.Nodes) {
			stateString = insolarNotReady
			color = tablewriter.FgHiRedColor
		}

		table.SetFooter([]string{
			"", "", "", "", "",
			"Insolar State", stateString,
			"Time", time.Now().Format(time.RFC3339),
		})
		table.SetFooterColor(
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},

			tablewriter.Colors{},
			tablewriter.Colors{color},

			tablewriter.Colors{},
			tablewriter.Colors{},
		)

		lock.Lock()
		table.AppendBulk(results)
		lock.Unlock()

		table.Render()

		fmt.Print(buffer)

		time.Sleep(conf.Interval)
	}
}
