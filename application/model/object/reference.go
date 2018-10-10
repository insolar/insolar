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

package object

import (
	"fmt"
)

// ScopeType represent type of scope for references.
type ScopeType int

// ChildScope, ContextScope and GlobalScope represents types of scope for references.
const (
	ChildScope = ScopeType(iota + 1)
	ContextScope
	GlobalScope
)

// Reference represents address of object.
type Reference interface {
	String() string
	GetRecord() string
	GetDomain() string
	GetScope() ScopeType
}

type reference struct {
	BaseObject
	domain string
	record string
	scope  ScopeType
}

// NewReference creates new reference instance.
func NewReference(domain string, record string, scope ScopeType) (Reference, error) {
	switch scope {
	case GlobalScope, ContextScope, ChildScope:
		return &reference{
			domain: domain,
			record: record,
			scope:  scope,
		}, nil
	default:
		return nil, fmt.Errorf("unknown scope type: %d", scope)
	}
}

// GetRecord return record value for current reference.
func (r *reference) GetRecord() string {
	return r.record
}

// GetDomain return domain value for current reference.
func (r *reference) GetDomain() string {
	return r.domain
}

// GetScope return scope value for current reference.
func (r *reference) GetScope() ScopeType {
	return r.scope
}

// String return string representation of reference
func (r *reference) String() string {
	return fmt.Sprintf("#%s.#%s", r.domain, r.record)
}
