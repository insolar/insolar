//
// Copyright 2019 Insolar Technologies GmbH
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
//

package main

import (
	"fmt"

	"github.com/dustin/go-humanize"

	"github.com/insolar/insolar/insolar/record"
)

func prettyPrintVirtual(v *record.Virtual) string {
	switch r := v.Union.(type) {
	case *record.Virtual_Genesis:
		return pairsToString(20, "Type", "Genesis")
	case *record.Virtual_IncomingRequest:
		return pairsToString(20, "Type", "IncomingRequest")
	case *record.Virtual_OutgoingRequest:
		return pairsToString(20, "Type", "OutgoingRequest")
	case *record.Virtual_Result:
		return pairsToString(20, "Type", "Result")
	case *record.Virtual_Code:
		return pairsToString(20, "Type", "Code")
	case *record.Virtual_Activate:
		return pairsToString(20, "Type", "Activate")
	case *record.Virtual_Amend:
		return amendPrettyPrint(r)
	case *record.Virtual_Deactivate:
		return pairsToString(20, "Type", "Deactivate")
	case *record.Virtual_PendingFilament:
		return pairsToString(20, "Type", "PendingFilament")
	case nil:
		return "nil"
	default:
		panic(fmt.Sprintf("%T virtual record unknown type", r))
	}
}

func amendPrettyPrint(virtualRecord *record.Virtual_Amend) string {
	rec := virtualRecord.Amend
	return pairsToString(20,
		"Type", "*record.Amend",
		"request", rec.Request.String(),
		"memory", humanize.Bytes(uint64(len(rec.Memory))),
		"image", rec.Image.String(),
		"isPrototype", fmt.Sprint(rec.IsPrototype),
		"prevState", rec.PrevState.String(),
	)
}
