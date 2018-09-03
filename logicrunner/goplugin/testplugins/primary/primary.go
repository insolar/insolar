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

package primary

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type INFHwRunner interface {
	Run() string
}

// nolint
type HwRunner struct {
	foundation.BaseContract
	Runned int
}

//
type Hw struct {
	Reference core.RecordRef
}

func (Hw) GetNewInstance(r core.RecordRef) Hw {
	return Hw{Reference: r}
}

//func (_self *Hw) Echo(s string) string {
//	foundation.APICall(_self.Reference, "Echo", s)
//	return s
//}
