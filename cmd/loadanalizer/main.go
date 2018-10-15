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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const defaultStdoutPath = "-"

func chooseOutput(path string) (io.Writer, error) {
	var res io.Writer
	if path == defaultStdoutPath {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't open file for writing")
		}
	}
	return res, nil
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func writeToOutput(out io.Writer, data string) {
	_, err := out.Write([]byte(data))
	check("Can't write data to output", err)
}

func getMembersRef(fileName string) ([]string, error) {
	var members []string

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open file for reading")
	}
	defer file.Close() //nolint: errcheck

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		members = append(members, scanner.Text())
	}

	return members, nil
}

func runScenarios(out io.Writer, members []string, concurrent int, repetitions int) {
	result := transferMoneyWithDifferentMember(members, concurrent, repetitions)
	writeToOutput(out, result)
}

const TestURL = "http://localhost:19191/api/v1"

type postParams map[string]interface{}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type baseResponse struct {
	Qid string         `json:"qid"`
	Err *errorResponse `json:"error"`
}

func (r *baseResponse) getError() *errorResponse {
	return r.Err
}

func getResponseBody(postParams map[string]interface{}) []byte {
	jsonValue, err := json.Marshal(postParams)
	check("Problems with marshal request:", err)
	postResp, err := http.Post(TestURL, "application/json", bytes.NewBuffer(jsonValue))
	check("Problems with post:", err)
	body, err := ioutil.ReadAll(postResp.Body)
	check("Problems with reading from response body:", err)
	return body
}

func transfer(amount int, from string, to string) string {
	body := getResponseBody(postParams{
		"query_type": "send_money",
		"from":       from,
		"to":         to,
		"amount":     amount,
	})

	response := &baseResponse{}
	json.Unmarshal(body, &response)
	if response.Err != nil {
		return response.Err.Message
	}
	return "success"
}

func transferMoneyWithDifferentMember(members []string, concurrent int, repetitions int) string {
	var result string
	var wg sync.WaitGroup

	if len(members) < concurrent*repetitions*2 {
		return "Not enough members"
	}
	fmt.Println("Start to transfer")

	start := time.Now()
	for i := 0; i < concurrent*repetitions*2; i = i + repetitions*2 {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for j := 0; j < repetitions*2; j = j + 2 {
				from := members[index+j]
				to := members[index+j+1]
				response := transfer(1, from, to)
				fmt.Printf("[Member №%d] Transfer from %s to %s. Response: %s.\n", index, from, to, response)
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Transfering took %s \n", elapsed)
	fmt.Printf("Speed - %f tr/s \n", float64(concurrent*repetitions)*1000000000/float64(elapsed))

	return result
}

type createMemberResponse struct {
	baseResponse
	Reference string `json:"reference"`
}

func createMembers(concurrent int, repetitions int) ([]string, error) {
	var members []string
	for i := 0; i < concurrent*repetitions*2; i++ {
		body := getResponseBody(postParams{
			"query_type": "create_membera",
			"name":       testutils.RandomString(),
			"public_key": "000",
		})

		memberResponse := &createMemberResponse{}
		json.Unmarshal(body, &memberResponse)

		if memberResponse.Err != nil {
			return nil, errors.New(memberResponse.Err.Message)
		}
		firstMemberRef := memberResponse.Reference
		members = append(members, firstMemberRef)
	}
	return members, nil
}

func main() {
	input := pflag.StringP("input", "i", "", "path to file with initial data for loads")
	output := pflag.StringP("output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	concurrent := pflag.IntP("concurrent", "c", 1, "concurrent users")
	repetitions := pflag.IntP("repetitions", "r", 1, "repetitions for one user")
	withInit := pflag.Bool("with_init", false, "do initialization before run load")
	pflag.Parse()

	out, err := chooseOutput(*output)
	check("Problems with output file:", err)

	var members []string

	if *withInit == true {
		members, err = createMembers(*concurrent, *repetitions)
		check("Problems with create members. One of creating request ended with error: ", err)
	}

	if *input != "" {
		members, err = getMembersRef(*input)
		check("Problems with parsing input:", err)
	}

	runScenarios(out, members, *concurrent, *repetitions)
}
