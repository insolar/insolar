package api

import (
	"context"
	"encoding/json"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
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
	log.Info("Send HISTORY request, ref: " + rh.params.Reference)
	log.Info("RouteResult: ", routResult)
	log.Info("RouteResult type: ", routResult.Type())

	refs := routResult.(*reply.ExplorerList).Refs
	log.Info("refs: ", refs)

	response, err := extractHistoryResponse(refs)
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

func extractHistoryResponse(refs []reply.Object) (string, error) {
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
