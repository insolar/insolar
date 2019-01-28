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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	pulsewatcher "github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var client http.Client

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

	for {
		results := make([]string, len(conf.Nodes))
		lock := &sync.Mutex{}
		wg := &sync.WaitGroup{}
		wg.Add(len(conf.Nodes))
		for i, url := range conf.Nodes {
			go func(url string, i int) {
				res, err := client.Post("http://"+url+"/api/rpc", "application/json",
					strings.NewReader(`{"jsonrpc": "2.0", "method": "status.Get", "id": 0}`))
				if err != nil {
					lock.Lock()
					results[i] = url + " : " + err.Error()
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
						Origin       struct {
							Role string
						}
					}
				}
				err = json.Unmarshal(data, &out)
				if err != nil {
					fmt.Println(string(data))
					log.Fatal(err)
				}
				lock.Lock()
				results[i] = url + " : " + out.Result.NetworkState + " : " + strconv.Itoa(int(out.Result.PulseNumber)) + " : " + out.Result.Origin.Role
				lock.Unlock()
				wg.Done()
			}(url, i)
		}
		wg.Wait()
		fmt.Println("\033[2J")
		fmt.Printf("%v\n\n", time.Now())
		lock.Lock()
		for _, result := range results {
			fmt.Println(result)
		}
		lock.Unlock()
		time.Sleep(conf.Interval)
	}
}
