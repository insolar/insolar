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
	"fmt"
	"net/http"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/messagerouter"
	"github.com/insolar/insolar/messagerouter/message"
	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

func makeRootDomainReference() core.RecordRef {
	const ref = "1111111-1111111-11111111-1111111"
	return core.String2Ref(base58.Encode([]byte(ref)))
}

var RootDomainReference = makeRootDomainReference()

func extractCreateMemberResponse(data []byte) (*string, error) {
	refOrig, err := UnMarshalResponse(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ extractCreateMemberResponse ]")
	}

	reference, ok := refOrig[0].(string)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to string. orig: %T", refOrig)
		return nil, errors.New(msg)
	}

	return &reference, nil
}

type RequestHandler struct {
	qid                 string
	req                 *http.Request
	messageRouter       *messagerouter.MessageRouter
	rootDomainReference core.RecordRef
}

func NewRequestHandler(r *http.Request, router *messagerouter.MessageRouter) *RequestHandler {
	return &RequestHandler{
		qid:                 r.FormValue("qid"),
		req:                 r,
		messageRouter:       router,
		rootDomainReference: RootDomainReference,
	}
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

func (rh *RequestHandler) ProcessCreateMember() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["CreateUser"] = true
	result["reference"] = "123123-234234234-345345345"

	const NameField = "name"
	name := rh.req.FormValue(NameField)
	if len(name) == 0 {
		return nil, errors.New("field is required: " + NameField)
	}

	routResult, err := rh.SendRequest("CreateMember", []interface{}{name})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessCreateMember ]")
	}

	memberRef, err := extractCreateMemberResponse(routResult.Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessCreateMember ]")
	}

	result["reference"] = memberRef

	return result, nil
}

func (rh *RequestHandler) ProcessDumpUserInfo() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["DumpUserInfo"] = true
	result["QQ"] = 222

	return result, nil
}

func extractGetBalanceResponse(data []byte) (uint, error) {
	dataUnmarsh, err := UnMarshalResponse(data)
	if err != nil {
		return 0, errors.Wrap(err, "[ extractGetBalanceResponse ]")
	}

	balance, ok := dataUnmarsh[0].(uint)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to uint. orig: %T", dataUnmarsh)
		return 0, errors.New(msg)
	}

	return balance, nil
}

func (rh *RequestHandler) SendRequest(method string, argsIn []interface{}) (*core.Response, error) {
	args, err := MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendRequest ]")
	}

	routResult, err := rh.RouteCall(rh.rootDomainReference, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendRequest ]")
	}

	return routResult, nil
}

func (rh *RequestHandler) ProcessGetBalance() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["GetBalance"] = true
	result["amount"] = 333
	result["currency"] = "RUB"

	const ReferenceField = "reference"
	name := rh.req.FormValue(ReferenceField)
	if len(name) == 0 {
		return nil, errors.New("field is required: " + ReferenceField)
	}

	routResult, err := rh.SendRequest("GetBalance", []interface{}{name})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetBalance ]")
	}

	amount, err := extractGetBalanceResponse(routResult.Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetBalance ]")
	}

	result["amount"] = amount

	return result, nil
}

func extractSendMoneyResponse(data []byte) (bool, error) {
	dataUnmarsh, err := UnMarshalResponse(data)
	if err != nil {
		return false, errors.Wrap(err, "[ extractSendMoneyResponse ]")
	}

	isSent, ok := dataUnmarsh[0].(bool)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to bool. orig: %T", dataUnmarsh)
		return false, errors.New(msg)
	}

	return isSent, nil
}

func (rh *RequestHandler) ProcessSendMoney() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["SendMoney"] = true
	result["success"] = true

	const FromField = "from"
	const ToField = "to"
	const AmountField = "to"
	from := rh.req.FormValue(FromField)
	if len(from) == 0 {
		return nil, errors.New("field is required: " + FromField)
	}
	to := rh.req.FormValue(ToField)
	if len(from) == 0 {
		return nil, errors.New("field is required: " + ToField)
	}
	amount := rh.req.FormValue(AmountField)
	if len(from) == 0 {
		return nil, errors.New("field is required: " + AmountField)
	}

	routResult, err := rh.SendRequest("SendMoney", []interface{}{from, to, amount})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessSendMoney ]")
	}

	isSent, err := extractSendMoneyResponse(routResult.Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessSendMoney ]")
	}

	result["success"] = isSent

	return result, nil
}

func extractDumpAllUsersResponse(data []byte) ([]byte, error) {
	dataUnmarsh, err := UnMarshalResponse(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ extractDumpAllUsersResponse ]")
	}

	dumpJson, ok := dataUnmarsh[0].([]byte)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to string. orig: %T", dataUnmarsh)
		return nil, errors.New(msg)
	}

	return dumpJson, nil
}

func (rh *RequestHandler) ProcessDumpAllUsers() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["CreateUser"] = true
	result["reference"] = "123123-234234234-345345345"

	const ReferenceField = "reference"
	name := rh.req.FormValue(ReferenceField)
	if len(name) == 0 {
		return nil, errors.New("field is required: " + ReferenceField)
	}

	routResult, err := rh.SendRequest("DumpAllUsers", []interface{}{name})
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessDumpAllUsers ]")
	}

	serJsonDump, err := extractDumpAllUsersResponse(routResult.Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessDumpAllUsers ]")
	}

	var dd interface{}
	json.Unmarshal(serJsonDump, &dd)
	result["dump_info"] = dd

	return result, nil
}
