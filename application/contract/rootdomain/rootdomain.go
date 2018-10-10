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
	"crypto/rand"
	"encoding/json"

	"github.com/insolar/insolar/application/proxy/member"
	"github.com/insolar/insolar/application/proxy/nodedomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	cryptoHelper "github.com/insolar/insolar/cryptohelpers/ecdsa"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// RootDomain is smart contract representing entrance point to system
type RootDomain struct {
	foundation.BaseContract
	RootMember core.RecordRef
}

// RegisterNode processes register node request
func (rd *RootDomain) RegisterNode(publicKey string, role string) string {
	domainRefs, err := rd.GetChildrenTyped(nodedomain.ClassReference)
	if err != nil {
		panic(err)
	}

	if len(domainRefs) == 0 {
		panic("No NodeDomain references")
	}
	nd := nodedomain.GetObject(domainRefs[0])

	return nd.RegisterNode(publicKey, role).String()
}

func makeSeed() []byte {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	return seed
}

// IsAuthorized checks is node authorized
func (rd *RootDomain) IsAuthorized() bool {
	privateKey, err := cryptoHelper.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	// Make signature
	seed := makeSeed()
	signature, err := cryptoHelper.Sign(seed, privateKey)
	if err != nil {
		panic(err)
	}

	// Register node
	serPubKey, err := cryptoHelper.ExportPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	nodeRef := rd.RegisterNode(serPubKey, "virtual")

	// Validate
	domainRefs, err := rd.GetChildrenTyped(nodedomain.ClassReference)
	if err != nil {
		panic(err)
	}
	nd := nodedomain.GetObject(domainRefs[0])

	return nd.IsAuthorized(core.NewRefFromBase58(nodeRef), seed, signature)
}

// CreateMember processes create member request
func (rd *RootDomain) CreateMember(name string, key string) string {
	//if rd.GetContext().Caller != nil && *rd.GetContext().Caller == *rd.RootMember {
	memberHolder := member.New(name, key)
	m := memberHolder.AsChild(rd.GetReference())
	wHolder := wallet.New(1000)
	wHolder.AsDelegate(m.GetReference())
	return m.GetReference().String()
	//}
	//return ""
}

// GetBalance processes get balance request
func (rd *RootDomain) GetBalance(reference string) uint {
	w := wallet.GetImplementationFrom(core.NewRefFromBase58(reference))
	return w.GetTotalBalance()
}

// SendMoney processes send money request
func (rd *RootDomain) SendMoney(from string, to string, amount uint) bool {
	walletFrom := wallet.GetImplementationFrom(core.NewRefFromBase58(from))
	v := core.NewRefFromBase58(to)
	walletFrom.Transfer(amount, &v)
	return true
}

func (rd *RootDomain) getUserInfoMap(m *member.Member) map[string]interface{} {
	w := wallet.GetImplementationFrom(m.GetReference())
	res := map[string]interface{}{
		"member": m.GetName(),
		"wallet": w.GetTotalBalance(),
	}
	return res
}

// DumpUserInfo processes dump user info request
func (rd *RootDomain) DumpUserInfo(reference string) []byte {
	m := member.GetObject(core.NewRefFromBase58(reference))
	res := rd.getUserInfoMap(m)
	resJSON, _ := json.Marshal(res)
	return resJSON
}

// DumpAllUsers processes dump all users request
func (rd *RootDomain) DumpAllUsers() []byte {
	res := []map[string]interface{}{}
	crefs, err := rd.GetChildrenTyped(member.ClassReference)
	if err != nil {
		panic(err)
	}
	for _, cref := range crefs {
		m := member.GetObject(cref)
		userInfo := rd.getUserInfoMap(m)
		res = append(res, userInfo)
	}
	resJSON, _ := json.Marshal(res)
	return resJSON
}

// NewRootDomain creates new RootDomain
func NewRootDomain() *RootDomain {
	return &RootDomain{}
}
