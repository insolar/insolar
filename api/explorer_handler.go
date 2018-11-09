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

package api

import (
	"context"
	"encoding/json"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
)

type explorerRecord struct {
	data  interface{}
	pulse string
}

type explorerObject struct {
	records []explorerRecord
	head    string
	parent  string
}

// ProcessGetHistory processes get history request
func (rh *RequestHandler) ProcessGetHistory(ctx context.Context) (map[string]interface{}, error) {

	result := make(map[string]interface{})
	reference := core.NewRefFromBase58(rh.params.Reference)
	routResult, err := rh.sendRequestHistory(ctx, reference)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetHistory ]")
	}
	response, err := extractHistoryResponse(routResult)
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetHistory ]")
	}
	result["history"] = response
	return result, nil
}

func (rh *RequestHandler) routeCallHistory(ctx context.Context, rootRef core.RecordRef, object core.RecordRef) (core.Reply, error) {
	if rh.messageBus == nil {
		return nil, errors.New("[ RouteCallHistory ] message bus was not set during initialization")
	}

	e := &message.GetHistory{
		Object: object,
	}

	res, err := rh.messageBus.Send(ctx, e)
	if err != nil {
		return nil, errors.Wrap(err, "[ RouteCallHistory ] couldn't send message")
	}

	return res, nil
}

func (rh *RequestHandler) sendRequestHistory(ctx context.Context, object core.RecordRef) (core.Reply, error) {
	routResult, err := rh.routeCallHistory(ctx, rh.rootDomainReference, object)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendRequest ]")
	}
	return routResult, nil
}

func extractHistoryResponse(routResult core.Reply) (string, error) {
	refs := routResult.(*reply.ExplorerList).Refs
	firstSt := true
	list := explorerObject{}
	for _, ref := range refs {
		resp, _ := core.UnMarshalResponse(ref.Memory, []interface{}{})
		elem := explorerRecord{
			data:  resp,
			pulse: string(ref.State.Pulse().Bytes()),
		}
		if firstSt {
			list.head = ref.Head.String()
			list.parent = ref.Head.String()
			firstSt = false
		}
		list.records = append(list.records, elem)
	}

	result, err := json.Marshal(list)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
