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

package genesisdataprovider

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/pkg/errors"
)

// GenesisDataProvider gives access to basic infotmation about genesis objects
type GenesisDataProvider struct {
	Certificate   core.Certificate `inject:""`
	MessageBus    core.MessageBus  `inject:""`
	nodeDomainRef *core.RecordRef
	rootDomainRef *core.RecordRef
	rootMemberRef *core.RecordRef
}

// New creates new GenesisDataProvider
func New() (*GenesisDataProvider, error) {
	return &GenesisDataProvider{}, nil
}

// RandomUint64 generates random uint64
func RandomUint64() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(buf)
}

func (gdp *GenesisDataProvider) routeCall(ctx context.Context, ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if gdp.MessageBus == nil {
		return nil, errors.New("[ GenesisDataProvider::routeCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		BaseLogicMessage: message.BaseLogicMessage{Nonce: RandomUint64()},
		ObjectRef:        ref,
		Method:           method,
		Arguments:        args,
	}

	res, err := gdp.MessageBus.Send(ctx, e)
	if err != nil {
		return nil, errors.Wrap(err, "[ GenesisDataProvider::routeCall ] couldn't send message: "+ref.String())
	}

	return res, nil
}

func extractInfoResponse(data []byte) (map[string]interface{}, error) {
	var infoMap interface{}
	var infoError *foundation.Error
	_, err := core.UnMarshalResponse(data, []interface{}{&infoMap, &infoError})
	if err != nil {
		return nil, errors.Wrap(err, "[ extractInfoResponse ] Can't unmarshal")
	}
	if infoError != nil {
		return nil, errors.Wrap(infoError, "[ extractInfoResponse ] Has error in response")
	}

	var info map[string]interface{}
	data = infoMap.([]byte)
	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, errors.Wrap(err, "[ extractInfoResponse ] Can't unmarshal response ")
	}

	return info, nil
}

func (gdp *GenesisDataProvider) sendRequest(ctx context.Context, ref *core.RecordRef, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ GenesisDataProvider::sendRequest ]")
	}

	routResult, err := gdp.routeCall(ctx, *ref, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ GenesisDataProvider::sendRequest ]")
	}

	return routResult, nil
}

func (gdp *GenesisDataProvider) setInfo(ctx context.Context) error {
	routResult, err := gdp.sendRequest(ctx, gdp.GetRootDomain(ctx), "Info", []interface{}{})
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] Can't send request")
	}

	info, err := extractInfoResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] Can't extract response")
	}
	rootMemberRef := core.NewRefFromBase58(info["root_member"].(string))
	gdp.rootMemberRef = &rootMemberRef
	nodeDomainRef := core.NewRefFromBase58(info["node_domain"].(string))
	gdp.nodeDomainRef = &nodeDomainRef

	return nil
}

// GetRootDomain returns reference to RootDomain
func (gdp *GenesisDataProvider) GetRootDomain(ctx context.Context) *core.RecordRef {
	if gdp.rootDomainRef == nil {
		gdp.rootDomainRef = gdp.Certificate.GetRootDomainReference()
	}
	return gdp.rootDomainRef
}

// GetNodeDomain returns reference to NodeDomain
func (gdp *GenesisDataProvider) GetNodeDomain(ctx context.Context) *core.RecordRef {
	if gdp.nodeDomainRef == nil {
		gdp.setInfo(ctx)
	}
	return gdp.nodeDomainRef
}

// GetRootMember returns reference to RootMember
func (gdp *GenesisDataProvider) GetRootMember(ctx context.Context) *core.RecordRef {
	if gdp.rootMemberRef == nil {
		gdp.setInfo(ctx)
	}
	return gdp.rootMemberRef
}
