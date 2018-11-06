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
	IsAuth
	GetSeed
)

// QTypeFromString converts string representation to enum
func QTypeFromString(strQType string) QueryType {
	switch strQType {
	case "is_auth":
		return IsAuth
	case "get_seed":
		return GetSeed
	}

	return UNDEFINED
}

// Params contains supported query params
type Params struct {
	QueryType              string   `json:"query_type"`
	Name                   string   `json:"name"`
	Reference              string   `json:"reference"`
	From                   string   `json:"from"`
	To                     string   `json:"to"`
	Method                 string   `json:"method"`
	Requester              string   `json:"requester"`
	Target                 string   `json:"target"`
	Amount                 uint     `json:"amount"`
	PublicKey              string   `json:"public_key"`
	Roles                  []string `json:"roles"`
	NumberOfBootstrapNodes uint     `json:"bootstrap_nodes_num"`
	MajorityRule           uint     `json:"majority_rule"`
	Host                   string   `json:"host"`
}
