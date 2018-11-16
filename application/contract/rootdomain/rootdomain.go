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

package rootdomain

import (
	"encoding/json"
	"fmt"

	"github.com/insolar/insolar/application/proxy/member"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// RootDomain is smart contract representing entrance point to system
type RootDomain struct {
	foundation.BaseContract
	RootMember    core.RecordRef
	NodeDomainRef core.RecordRef
}

// CreateMember processes create member request
func (rd *RootDomain) CreateMember(name string, key string) (string, error) {
	if *rd.GetContext().Caller != rd.RootMember {
		return "", fmt.Errorf("[ CreateMember ] Only Root member can create members")
	}
	memberHolder := member.New(name, key)
	m, err := memberHolder.AsChild(rd.GetReference())
	if err != nil {
		return "", fmt.Errorf("[ CreateMember ] Can't save as child: %s", err.Error())
	}

	wHolder := wallet.New(1000)
	_, err = wHolder.AsDelegate(m.GetReference())
	if err != nil {
		return "", fmt.Errorf("[ CreateMember ] Can't save as delegate: %s", err.Error())
	}

	return m.GetReference().String(), nil
}

func (rd *RootDomain) getUserInfoMap(m *member.Member) (map[string]interface{}, error) {
	w, err := wallet.GetImplementationFrom(m.GetReference())
	if err != nil {
		return nil, fmt.Errorf("[ getUserInfoMap ] Can't get implementation: %s", err.Error())
	}

	name, err := m.GetName()
	if err != nil {
		return nil, fmt.Errorf("[ getUserInfoMap ] Can't get name: %s", err.Error())
	}

	balance, err := w.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("[ getUserInfoMap ] Can't get total balance: %s", err.Error())
	}
	return map[string]interface{}{
		"member": name,
		"wallet": balance,
	}, nil
}

// DumpUserInfo processes dump user info request
func (rd *RootDomain) DumpUserInfo(reference string) ([]byte, error) {
	caller := *rd.GetContext().Caller
	ref := core.NewRefFromBase58(reference)
	if ref != caller && caller != rd.RootMember {
		return nil, fmt.Errorf("[ DumpUserInfo ] You can dump only yourself")
	}
	m := member.GetObject(ref)

	res, err := rd.getUserInfoMap(m)
	if err != nil {
		return nil, fmt.Errorf("[ DumpUserInfo ] Problem with making request: %s", err.Error())
	}

	return json.Marshal(res)
}

// DumpAllUsers processes dump all users request
func (rd *RootDomain) DumpAllUsers() ([]byte, error) {
	if *rd.GetContext().Caller != rd.RootMember {
		return nil, fmt.Errorf("[ DumpUserInfo ] Only root can call this method")
	}
	res := []map[string]interface{}{}
	crefs, err := rd.GetChildrenTyped(member.PrototypeReference)
	if err != nil {
		return nil, fmt.Errorf("[ DumpUserInfo ] Can't get children: %s", err.Error())
	}
	for _, cref := range crefs {
		if cref == rd.RootMember {
			continue
		}
		m := member.GetObject(cref)
		userInfo, err := rd.getUserInfoMap(m)
		if err != nil {
			return nil, fmt.Errorf("[ DumpAllUsers ] Problem with making request: %s", err.Error())
		}
		res = append(res, userInfo)
	}
	resJSON, _ := json.Marshal(res)
	return resJSON, nil
}

var INSATTR_GetNodeDomainRef_API = true

// GetNodeDomainRef returns reference of NodeDomain instance
func (rd *RootDomain) GetNodeDomainRef() (core.RecordRef, error) {
	return rd.NodeDomainRef, nil
}

// NewRootDomain creates new RootDomain
func NewRootDomain() (*RootDomain, error) {
	return &RootDomain{}, nil
}
