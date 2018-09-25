/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

// Package caserecord is various types of records for validating contract execution
package caserecord

import (
	"github.com/insolar/insolar/core"
	"golang.org/x/crypto/sha3"
)

func Hash(in []byte) []byte {
	sh := sha3.New224()
	return sh.Sum(in)
}

// Base is abstract base record
type Base struct {
}

type Incoming struct {
	Base
	Event core.Event
}

type Result struct {
	Base
	Sig []byte
}
