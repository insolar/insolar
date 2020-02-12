// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/insolar/insolar/insolar/record"
)

func prettyPrintVirtual(v *record.Virtual) string {
	pf := pairFormatter{width: 20}
	switch r := v.Union.(type) {
	case *record.Virtual_Genesis:
		return pf.Pairs("Type", "Genesis")
	case *record.Virtual_IncomingRequest:
		return pf.Pairs("Type", "IncomingRequest")
	case *record.Virtual_OutgoingRequest:
		return pf.Pairs("Type", "OutgoingRequest")
	case *record.Virtual_Result:
		return pf.Pairs("Type", "Result")
	case *record.Virtual_Code:
		return pf.Pairs("Type", "Code")
	case *record.Virtual_Activate:
		return pf.Pairs("Type", "Activate")
	case *record.Virtual_Amend:
		return amendPretty(r)
	case *record.Virtual_Deactivate:
		return pf.Pairs("Type", "Deactivate")
	case *record.Virtual_PendingFilament:
		return pf.Pairs("Type", "PendingFilament")
	case nil:
		return "nil"
	default:
		panic(fmt.Sprintf("%T virtual record unknown type", r))
	}
}

func amendPretty(virtualRecord *record.Virtual_Amend) string {
	pf := pairFormatter{width: 20}
	rec := virtualRecord.Amend
	return pf.Pairs(
		"Type", "*record.Amend",
		"request", rec.Request.String(),
		"memory", humanize.Bytes(uint64(len(rec.Memory))),
		"image", rec.Image.String(),
		"isPrototype", fmt.Sprint(rec.IsPrototype),
		"prevState", rec.PrevState.String(),
	)
}
