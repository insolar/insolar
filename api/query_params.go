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

// QueryType represents type of query
type QueryType int

// Supported query types
const (
	UNDEFINED QueryType = iota
	CreateMember
	DumpUserInfo
	GetBalance
	SendMoney
	DumpAllUsers
	RegisterNode
	IsAuth
)

// QTypeFromString converts string representation to enum
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
	case "dump_all_users":
		return DumpAllUsers
	case "register_node":
		return RegisterNode
	case "is_auth":
		return IsAuth
	}

	return UNDEFINED
}

// Params contains supported query params
type Params struct {
	QType     string `json:"query_type"`
	Name      string `json:"name"`
	Reference string `json:"reference"`
	From      string `json:"from"`
	To        string `json:"to"`
	QID       string `json:"qid"`
	Amount    uint   `json:"amount"`
	PublicKey string `json:"public_key"`
	Role      string `json:"role"`
}
