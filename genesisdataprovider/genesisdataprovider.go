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
	"sync"

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/pkg/errors"
)

// GenesisDataProvider gives access to basic information about genesis objects
type GenesisDataProvider struct {
	Certificate       core.Certificate       `inject:""`
	ContractRequester core.ContractRequester `inject:""`

	nodeDomainRef     *core.RecordRef
	nodeDomainRefLock sync.RWMutex

	rootDomainRef     *core.RecordRef
	rootDomainRefLock sync.RWMutex

	rootMemberRef     *core.RecordRef
	rootMemberRefLock sync.RWMutex
}

// New creates new GenesisDataProvider
func New() (*GenesisDataProvider, error) {
	return &GenesisDataProvider{}, nil
}

func (gdp *GenesisDataProvider) setInfo(ctx context.Context) error {
	routResult, err := gdp.ContractRequester.SendRequest(ctx, gdp.GetRootDomain(ctx), "Info", []interface{}{})
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] Can't send request")
	}

	info, err := extractor.InfoResponse(routResult.(*reply.CallMethod).Result)
	if err != nil {
		return errors.Wrap(err, "[ setInfo ] Can't extract response")
	}

	rootMemberRef := core.NewRefFromBase58(info.RootMember)
	gdp.rootMemberRefLock.Lock()
	gdp.rootMemberRef = &rootMemberRef
	gdp.rootMemberRefLock.Unlock()

	nodeDomainRef := core.NewRefFromBase58(info.NodeDomain)
	gdp.nodeDomainRefLock.Lock()
	gdp.nodeDomainRef = &nodeDomainRef
	gdp.nodeDomainRefLock.Unlock()

	return nil
}

// GetRootDomain returns reference to RootDomain
func (gdp *GenesisDataProvider) GetRootDomain(ctx context.Context) *core.RecordRef {
	gdp.rootDomainRefLock.RLock()
	defer gdp.rootDomainRefLock.RUnlock()
	if gdp.rootDomainRef == nil {
		gdp.rootDomainRefLock.RUnlock()
		gdp.rootDomainRefLock.Lock()
		gdp.rootDomainRef = gdp.Certificate.GetRootDomainReference()
		gdp.rootDomainRefLock.Unlock()
		gdp.rootDomainRefLock.RLock()
	}
	return gdp.rootDomainRef
}

// GetNodeDomain returns reference to NodeDomain
func (gdp *GenesisDataProvider) GetNodeDomain(ctx context.Context) (*core.RecordRef, error) {
	gdp.nodeDomainRefLock.RLock()
	defer gdp.nodeDomainRefLock.RUnlock()
	if gdp.nodeDomainRef == nil {
		gdp.nodeDomainRefLock.RUnlock()
		err := gdp.setInfo(ctx)
		gdp.nodeDomainRefLock.RLock()
		if err != nil {
			return nil, errors.Wrap(err, "[ GenesisDataProvider::GetNodeDomain ] Can't get info")
		}
	}
	return gdp.nodeDomainRef, nil
}

// GetRootMember returns reference to RootMember
func (gdp *GenesisDataProvider) GetRootMember(ctx context.Context) (*core.RecordRef, error) {
	gdp.rootMemberRefLock.RLock()
	defer gdp.rootMemberRefLock.RUnlock()
	if gdp.rootMemberRef == nil {
		gdp.rootMemberRefLock.RUnlock()
		err := gdp.setInfo(ctx)
		gdp.rootMemberRefLock.RLock()
		if err != nil {
			return nil, errors.Wrap(err, "[ GenesisDataProvider::GetRootMember ] Can't get info")
		}
	}
	return gdp.rootMemberRef, nil
}
