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
	"fmt"
	"net/http"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/messagerouter"
	"github.com/insolar/insolar/messagerouter/message"
	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

func makeRootDomainReference() core.RecordRef {
	const ref = "1111111-1111111-11111111-1111111"
	return core.String2Ref(base58.Encode([]byte(ref)))
}

var RootDomainReference = makeRootDomainReference()

func extractCreateMemberResponse(data []byte) (*string, error) {
	var marshRes interface{}
	marshRes, err := CBORUnMarshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "[ extractCreateMemberResponse ]")
	}

	refOrig, ok := marshRes.([1]interface{})
	if !ok || len(refOrig) < 0 {
		return nil, errors.New("[ extractCreateMemberResponse ] Problem with extracting result")
	}

	reference, ok := refOrig[0].(string)
	if !ok {
		msg := fmt.Sprintf("Can't cast response to string. orig: %T", refOrig[0])
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

	memberRef, err := extractCreateMemberResponse(routResult.Result)
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
