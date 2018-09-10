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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/messagerouter"
)

type CmdParams struct {
	port uint
}

var cmdParams CmdParams

func ParseInputParams() {
	flag.UintVar(&cmdParams.port, "port", 8080, "listening port")
}

type QueryType int

const (
	UNDEFINED QueryType = iota
	CreateMember
	DumpUserInfo
	GetBalance
	SendMoney
	DumpAllUsers
)

func QTypeFromString(strQType string) QueryType {
	switch strQType {
	case "create_member":
		return CreateMember
	case "dump_user_info":
		return DumpUserInfo
	case "get_balance":
		return GetBalance
	case "send_money":
		return SendMoney
	case "dump_all_user":
		return DumpAllUsers
	}

	return UNDEFINED
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "0123456789abcdef"

func GenQID() string {
	b := [16]byte{}
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b[:])
}

func WriteError(message string, code int) map[string]interface{} {
	errJson := map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"code":    code,
		},
	}
	return errJson
}

func ProcessCreateMember(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	result := make(map[string]interface{})
	result["CreateUser"] = true
	result["reference"] = "123123-234234234-345345345"

	return result
}

func ProcessDumpUserInfo(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	result := make(map[string]interface{})
	result["DumpUserInfo"] = true
	result["Putin"] = 222

	return result
}

func ProcessGetBalance(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	result := make(map[string]interface{})
	result["GetBalance"] = true
	result["amount"] = 333
	result["currency"] = "RUB"

	return result
}

func ProcessSendMoney(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	result := make(map[string]interface{})
	result["SendMoney"] = true
	result["success"] = true

	return result
}

func ProcessDumpAllUsers(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	result := make(map[string]interface{})
	result["AllUsers"] = true
	result["Putin"] = 555

	return result
}

func MakeHandlerMarshalErrorJson() []byte {
	jsonErr := WriteError("Invalid data from handler", -1)
	serJson, _ := json.Marshal(jsonErr)
	return serJson
}

var handlerMarshalErrorJson = MakeHandlerMarshalErrorJson()

func ApiV1Handler(w http.ResponseWriter, r *http.Request) {
	answer := make(map[string]interface{})
	qid := GenQID()
	defer func() {
		answer["qid"] = qid
		serJson, err := json.MarshalIndent(answer, "", "    ")
		if err != nil {
			serJson = handlerMarshalErrorJson
		}
		var newLine byte = '\n'
		w.Write(append(serJson, newLine))
	}()

	r.ParseForm()
	qTypeStr := r.FormValue("query_type")
	qtype := QTypeFromString(qTypeStr)
	switch qtype {
	case CreateMember:
		answer = ProcessCreateMember(w, r)
	case DumpUserInfo:
		answer = ProcessDumpUserInfo(w, r)
	case GetBalance:
		answer = ProcessGetBalance(w, r)
	case SendMoney:
		answer = ProcessSendMoney(w, r)
	case DumpAllUsers:
		answer = ProcessDumpAllUsers(w, r)
	default:
		answer = WriteError(fmt.Sprintf("Wrong query parameter 'query_type' = '%s'", qTypeStr), -2)
		return
	}
}

type ApiRunner struct {
	messageRouter *messagerouter.MessageRouter
	server        *http.Server
}

func (ar *ApiRunner) Start(c core.Components) error {
	//ar.messageRouter = c["core.MessageRouter"].(*messagerouter.MessageRouter)

	ar.server = &http.Server{Addr: ":" + fmt.Sprint(cmdParams.port)}
	http.HandleFunc("/api/v1", ApiV1Handler)
	go func() {
		if err := ar.server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	return nil
}

func (ar *ApiRunner) Stop() error {
	ar.server.Shutdown(nil)
	return nil
}

func main() {
	ParseInputParams()
	api := ApiRunner{}
	cs := core.Components{}
	api.Start(cs)
	time.Sleep(100 * time.Second)
}
