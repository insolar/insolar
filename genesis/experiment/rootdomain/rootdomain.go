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

package rootdomain

import (
	"encoding/json"

	"contract-proxy/member"
	"contract-proxy/wallet"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/toolkit/go/foundation"
)

type RootDomain struct {
	foundation.BaseContract
}

func (rd *RootDomain) CreateMember(name string) string {
	memberHolder := member.NewMember(name)
	m := memberHolder.AsChild(rd.GetReference())
	wHolder := wallet.NewWallet(1000)
	wHolder.AsDelegate(m.GetReference())
	return m.GetReference().String()
}

func (rd *RootDomain) GetBalance(reference string) uint {
	memberAsWallet := foundation.GetImplementationFor(wallet.ClassReference, core.String2Ref(reference))
	w, _ := memberAsWallet.(*wallet.Wallet)
	return w.GetTotalBalance()
}

func (rd *RootDomain) SendMoney(from string, to string, amount uint) bool {
	memberFrom := foundation.GetImplementationFor(wallet.ClassReference, core.String2Ref(from))
	walletfrom, ok := memberFrom.(*wallet.Wallet)
	if !ok {
		return false
	}
	memberTo := foundation.GetImplementationFor(wallet.ClassReference, core.String2Ref(to))
	walletTo, ok := memberTo.(*wallet.Wallet)
	if !ok {
		return false
	}

	walletfrom.Transfer(amount, walletTo.GetReference())

	return true
}

func (rd *RootDomain) getUserInfoMap(m *member.Member) map[string]interface{} {
	memberAsWallet := foundation.GetImplementationFor(wallet.ClassReference, m.GetReference())
	w, _ := memberAsWallet.(*wallet.Wallet)
	res := map[string]interface{}{}
	res["member"] = m
	res["wallet"] = w
	return res
}

func (rd *RootDomain) DumpUserInfo(reference string) []byte {
	m := member.GetObject(core.String2Ref(reference))
	res := rd.getUserInfoMap(m)
	resJson, _ := json.Marshal(res)
	return resJson
}

func (rd *RootDomain) DumpAllUsers() []byte {
	res := []map[string]interface{}{}
	for _, c := range rd.GetChildrenTyped(member.ClassReference) {
		m := c.(*member.Member)
		userInfo := rd.getUserInfoMap(m)
		res = append(res, userInfo)
	}
	resJson, _ := json.Marshal(res)
	return resJson
}

func NewRootDomain() *RootDomain {
	return &RootDomain{}
}
