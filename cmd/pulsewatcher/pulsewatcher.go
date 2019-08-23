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

	"github.com/insolar/insolar/api/requester"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	pulsewatcher "github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/insolar/insolar/insolar"
)

var client http.Client
var emoji *Emoji
var startTime time.Time

const (
	esc        = "\x1b%s"
	moveUp     = "[%dA"
	clearDown  = "[0J"
	timeFormat = "15:04:05.999999"
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

func displayResultsTable(results []nodeStatus, ready bool, buffer *bytes.Buffer) {
	table := tablewriter.NewWriter(buffer)
	table.SetHeader([]string{
		"URL",
		"State",
		"ID",
		"Network Pulse",
		"Pulse",
		"Active",
		"Working",
		"Role",
		"Timestamp",
		"Uptime",
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
		"Time", time.Now().Format(timeFormat),
		"Insolar Uptime", time.Since(startTime).Round(time.Second).String(), "",
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
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
	)
	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},

		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgHiRedColor},
	)

	intToString := func(n int) string {
		if n == 0 {
			return ""
		}
		return strconv.Itoa(n)
	}

	shortRole := func(r string) string {
		switch r {
		case "virtual":
			return "Virtual"
		case "heavy_material":
			return "Heavy"
		case "light_material":
			return "Light"
		default:
			return r
		}
	}

	for _, row := range results {
		emoji.RegisterNode(row.url, row.reply.Origin)
	}

	for _, row := range results {
		var activeNodeEmoji string
		for _, n := range row.reply.Nodes {
			activeNodeEmoji += emoji.GetEmoji(n)
		}

		var uptime string
		var timestamp string
		if row.errStr == "" {
			uptime = time.Since(row.reply.StartTime).Round(time.Second).String()
			timestamp = row.reply.Timestamp.Format(timeFormat)
		}

		table.Append([]string{
			row.url,
			row.reply.NetworkState,
			fmt.Sprintf(" %s %s", emoji.GetEmoji(row.reply.Origin), intToString(int(row.reply.Origin.ID))),
			intToString(int(row.reply.NetworkPulseNumber)),
			intToString(int(row.reply.PulseNumber)),
			fmt.Sprintf("%d %s", row.reply.ActiveListSize, activeNodeEmoji),
			intToString(row.reply.WorkingListSize),
			shortRole(row.reply.Origin.Role),
			timestamp,
			uptime,
			row.errStr,
		})
	}
	table.Render()
	fmt.Print(buffer)
}

func displayResultsJSON(results []nodeStatus) {
	type DocumentItem struct {
		URL                string
		NetworkState       string
		ID                 uint32
		NetworkPulseNumber uint32
		PulseNumber        uint32
		ActiveListSize     int
		WorkingListSize    int
		Role               string
		Timestamp          string
		Error              string
	}

	doc := make([]DocumentItem, len(results))

	for i, res := range results {
		doc[i].URL = res.url
		doc[i].NetworkState = res.reply.NetworkState
		doc[i].ID = res.reply.Origin.ID
		doc[i].NetworkPulseNumber = res.reply.NetworkPulseNumber
		doc[i].PulseNumber = res.reply.PulseNumber
		doc[i].ActiveListSize = res.reply.ActiveListSize
		doc[i].WorkingListSize = res.reply.WorkingListSize
		doc[i].Role = res.reply.Origin.Role
		doc[i].Timestamp = res.reply.Timestamp.Format(timeFormat)
		doc[i].Error = res.errStr
	}

	jsonDoc, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		panic(err) // should never happen
	}
	fmt.Print(string(jsonDoc))
	fmt.Print("\n\n")
}

func collectNodesStatuses(conf *pulsewatcher.Config, lastResults []nodeStatus) ([]nodeStatus, bool) {
	state := true
	errored := 0
	results := make([]nodeStatus, len(conf.Nodes))
	lock := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	wg.Add(len(conf.Nodes))
	for i, url := range conf.Nodes {
		go func(url string, i int) {
			res, err := client.Post("http://"+url+"/api/rpc", "application/json",
				strings.NewReader(`{"jsonrpc": "2.0", "method": "node.getStatus", "id": 0}`))

			url = strings.TrimPrefix(url, "127.0.0.1")

			if err != nil {
				errStr := err.Error()
				if strings.Contains(errStr, "connection refused") ||
					strings.Contains(errStr, "request canceled while waiting for connection") ||
					strings.Contains(errStr, "no such host") {
					// Print compact error string when node is down.
					// This prevents table distortion on small screens.
					errStr = "NODE IS DOWN"
				}
				if strings.Contains(errStr, "exceeded while awaiting headers") {
					errStr = "TIMEOUT"
				}

				lock.Lock()
				if len(lastResults) > i {
					results[i] = lastResults[i]
					results[i].errStr = errStr
				} else {
					results[i] = nodeStatus{url, requester.StatusResponse{}, errStr}
				}
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
				Result requester.StatusResponse
			}
			err = json.Unmarshal(data, &out)
			if err != nil {
				fmt.Println(string(data))
				log.Fatal(err)
			}
			lock.Lock()

			results[i] = nodeStatus{url, out.Result, ""}
			state = state && out.Result.NetworkState == insolar.CompleteNetworkState.String()
			lock.Unlock()
			wg.Done()
		}(url, i)
	}
	wg.Wait()

	ready := state && errored != len(conf.Nodes)
	return results, ready
}

type nodeStatus struct {
	url    string
	reply  requester.StatusResponse
	errStr string
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

	emoji = NewEmoji()
	var results []nodeStatus
	var ready bool
	startTime = time.Now()
	for {
		results, ready = collectNodesStatuses(conf, results)
		if useJSONFormat {
			displayResultsJSON(results)
		} else {
			displayResultsTable(results, ready, buffer)
		}

		if singleOutput {
			break
		}

		time.Sleep(conf.Interval)
	}
}
