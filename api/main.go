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
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/messagerouter"
	"github.com/insolar/insolar/messagerouter/message"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
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

func (rh *RequestHandler) RouteCall(ref core.RecordRef, method string, args core.Arguments) (*core.Response, error) {
	if rh.messageRouter == nil {
		return nil, errors.New("[ RouteCall ] message router was not set during initialization")
	}

	msg := &message.CallMethodMessage{
		ObjectRef: ref,
		Method:    method,
		Arguments: args,
	}

	res, err := rh.messageRouter.Route(msg)
	if err != nil {
		return nil, errors.Wrap(err, "[ RouteCall ] couldn't route message")
	}
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "[ RouteCall ] couldn't route message (error in response)")
	}

	return &res, nil
}

func CBORMarshal(o interface{}) ([]byte, error) {
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(o)
	return data, err
}

// CBORUnMarshal - testing deserialize helper
func CBORUnMarshal(data []byte) (interface{}, error) {
	ch := new(codec.CborHandle)
	var ret interface{}
	err := codec.NewDecoderBytes(data, ch).Decode(&ret)
	return ret, errors.Wrap(err, "[ CBORUnMarshal ]")
}

func MarshalArgs(args ...interface{}) (core.Arguments, error) {
	var argsSerialized []byte

	argsSerialized, err := CBORMarshal(args)
	if err != nil {
		return nil, errors.Wrap(err, "[ MarshalArgs ]")
	}

	result := core.Arguments(argsSerialized)

	return result, nil
}

func ExtractCreateMemberResponse(data []byte) (*string, error) {
	var marshRes interface{}
	marshRes, err := CBORUnMarshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ ExtractCreateMemberResponse ]")
	}

	refOrig, ok := marshRes.([1]interface{})
	if !ok || len(refOrig) < 0 {
		return nil, errors.New("[ ExtractCreateMemberResponse ] Problem with extracting result")
	}

	reference, ok := refOrig[0].(string)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to string. orig: %T", refOrig[0])
		return nil, errors.New(msg)
	}

	return &reference, nil
}

func (rh *RequestHandler) ProcessCreateMember() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["CreateUser"] = true
	result["reference"] = "123123-234234234-345345345"

	name := rh.req.FormValue("name")
	if len(name) == 0 {
		return nil, errors.New("field 'name' is required")
	}

	args, err := MarshalArgs(name)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessCreateMember ]")
	}

	routResult, err := rh.RouteCall(rh.rootDomainReference, "CreateMember", args)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessCreateMember ]")
	}

	memberRef, err := ExtractCreateMemberResponse(routResult.Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessCreateMember ]")
	}

	if len(*memberRef) != 0 {
		result["reference"] = memberRef
	}

	return result, nil
}

func (rh *RequestHandler) ProcessDumpUserInfo() map[string]interface{} {
	result := make(map[string]interface{})
	result["DumpUserInfo"] = true
	result["QQ"] = 222

	return result
}

func (rh *RequestHandler) ProcessGetBalance() map[string]interface{} {
	result := make(map[string]interface{})
	result["GetBalance"] = true
	result["amount"] = 333
	result["currency"] = "RUB"

	return result
}

func (rh *RequestHandler) ProcessSendMoney() map[string]interface{} {
	result := make(map[string]interface{})
	result["SendMoney"] = true
	result["success"] = true

	return result
}

func (rh *RequestHandler) ProcessDumpAllUsers() map[string]interface{} {
	result := make(map[string]interface{})
	result["AllUsers"] = true
	result["QQQ"] = 555

	return result
}

func MakeHandlerMarshalErrorJson() []byte {
	jsonErr := WriteError("Invalid data from handler", -1)
	serJson, err := json.Marshal(jsonErr)
	if err != nil {
		log.Fatal("Can't marshal base error")
	}
	return serJson
}

var handlerMarshalErrorJson = MakeHandlerMarshalErrorJson()

type RequestHandler struct {
	qid                 string
	req                 *http.Request
	messageRouter       *messagerouter.MessageRouter
	rootDomainReference core.RecordRef
}

func MakeRootDomainReference() core.RecordRef {
	const ref = "1111111-1111111-11111111-1111111"
	return core.String2Ref(ref)
}

var RootDomainReference = MakeRootDomainReference()

func NewRequestHandler(r *http.Request, router *messagerouter.MessageRouter) *RequestHandler {
	return &RequestHandler{
		qid:                 r.FormValue("qid"),
		req:                 r,
		messageRouter:       router,
		rootDomainReference: RootDomainReference,
	}
}

func WrapApiV1Handler(router *messagerouter.MessageRouter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		answer := make(map[string]interface{})
		qid := GenQID()
		rh := NewRequestHandler(req, router)
		defer func() {
			if answer == nil {
				answer = make(map[string]interface{})
			}
			answer["qid"] = qid
			serJson, err := json.MarshalIndent(answer, "", "    ")
			if err != nil {
				serJson = handlerMarshalErrorJson
			}
			var newLine byte = '\n'
			w.Write(append(serJson, newLine))
			log.Printf("[QID=%s] Request completed\n", qid)
		}()

		req.ParseForm()
		req.Form.Add("qid", qid)
		log.Printf("[QID=%s] Query: %s\n", qid, req.RequestURI)

		qTypeStr := req.FormValue("query_type")
		qtype := QTypeFromString(qTypeStr)

		var handlerError error
		switch qtype {
		case CreateMember:
			answer, handlerError = rh.ProcessCreateMember()
		case DumpUserInfo:
			answer = rh.ProcessDumpUserInfo()
		case GetBalance:
			answer = rh.ProcessGetBalance()
		case SendMoney:
			answer = rh.ProcessSendMoney()
		case DumpAllUsers:
			answer = rh.ProcessDumpAllUsers()
		default:
			msg := fmt.Sprintf("Wrong query parameter 'query_type' = '%s'", qTypeStr)
			answer = WriteError(msg, -2)
			log.Printf("[QID=%s] %s\n", qid, msg)
			return
		}
		if handlerError != nil {
			errMsg := "Handler error: " + handlerError.Error()
			log.Printf("[QID=%s] %s\n", qid, errMsg)
			answer = WriteError(errMsg, -3)
		}
	}
}

type ApiRunner struct {
	messageRouter *messagerouter.MessageRouter
	server        *http.Server
}

func (ar *ApiRunner) Start(c core.Components) error {

	//ar.messageRouter = c["core.MessageRouter"].(*messagerouter.MessageRouter)

	ar.server = &http.Server{Addr: ":" + fmt.Sprint(cmdParams.port)}
	fw := WrapApiV1Handler(ar.messageRouter)
	http.HandleFunc("/api/v1", fw)
	go func() {
		if err := ar.server.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s\n", err)
		}
	}()
	return nil
}

func (ar *ApiRunner) Stop() error {
	const timeOut = 5
	log.Printf("Shutting down server gracefully ...(waiting for %d seconds)\n", timeOut)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	err := ar.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop API server")
	}

	return nil
}

func main() {
	ParseInputParams()
	api := ApiRunner{}
	cs := core.Components{}
	api.Start(cs)
	time.Sleep(60 * time.Second)
	api.Stop()
}
