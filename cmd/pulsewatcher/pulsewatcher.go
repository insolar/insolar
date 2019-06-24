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

	pulsewatcher "github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/insolar/insolar/insolar"
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

func displayResultsTable(results [][]string, ready bool, buffer *bytes.Buffer) {
	table := tablewriter.NewWriter(buffer)
	table.SetHeader([]string{
		"URL",
		"Network State",
		"NetworkNode State",
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
	if !ready {
		stateString = insolarNotReady
		color = tablewriter.FgHiRedColor
	}

	table.SetFooter([]string{
		"", "", "", "",
		"Insolar State", stateString,
		"Time", time.Now().Format(time.RFC3339),
	})
	table.SetFooterColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},

		tablewriter.Colors{},
		tablewriter.Colors{color},

		tablewriter.Colors{},
		tablewriter.Colors{},
	)

	table.AppendBulk(results)
	table.Render()
	fmt.Print(buffer)
}

func parseInt64(str string) int64 {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		res = -1
	}
	return res
}

func displayResultsJSON(results [][]string, _ bool, _ *bytes.Buffer) {
	type DocumentItem struct {
		URL             string
		NetworkState    string
		NodeState       string
		PulseNumber     int64
		ActiveListSize  int64
		WorkingListSize int64
		Role            string
		Error           string
	}

	doc := make([]DocumentItem, len(results))

	for i, res := range results {
		doc[i].URL = res[0]
		doc[i].NetworkState = res[1]
		doc[i].NodeState = res[2]
		doc[i].PulseNumber = parseInt64(res[3])
		doc[i].ActiveListSize = parseInt64(res[4])
		doc[i].WorkingListSize = parseInt64(res[5])
		doc[i].Role = res[6]
		doc[i].Error = res[7]
	}

	jsonDoc, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		panic(err) // should never happen
	}
	fmt.Print(string(jsonDoc))
	fmt.Print("\n\n")
}

func collectNodesStatuses(conf *pulsewatcher.Config) ([][]string, bool) {
	state := true
	errored := 0
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
				results[i] = []string{url, "", "", "", "", "", "", err.Error()}
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
					PulseNumber  uint32
					NetworkState string
					NodeState    string
					Origin       struct {
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
				strconv.Itoa(int(out.Result.PulseNumber)),
				strconv.Itoa(out.Result.ActiveListSize),
				strconv.Itoa(out.Result.WorkingListSize),
				out.Result.Origin.Role,
				"",
			}
			state = state && out.Result.NetworkState == insolar.CompleteNetworkState.String() &&
				out.Result.NodeState == insolar.NodeReady.String()
			lock.Unlock()
			wg.Done()
		}(url, i)
	}
	wg.Wait()

	ready := state && errored != len(conf.Nodes)
	return results, ready
}

func main() {
	var configFile string
	var useJSONFormat bool
	var singleOutput bool
	pflag.StringVarP(&configFile, "config", "c", "", "config file")
	pflag.BoolVarP(&useJSONFormat, "json", "j", false, "use JSON format")
	pflag.BoolVarP(&singleOutput, "single", "s", false, "single output")
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

	buffer := &bytes.Buffer{}
	fmt.Print("\n\n")

	client = http.Client{
		Transport: &http.Transport{},
		Timeout:   conf.Timeout,
	}

	for {
		results, ready := collectNodesStatuses(conf)
		if useJSONFormat {
			displayResultsJSON(results, ready, buffer)
		} else {
			displayResultsTable(results, ready, buffer)
		}

		if singleOutput {
			break
		}

		time.Sleep(conf.Interval)
	}
}
