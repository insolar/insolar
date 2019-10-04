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
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/insolar/insolar/insolar/record"
)

func prettyPrintVirtual(v *record.Virtual) string {
	switch r := v.Union.(type) {
	case *record.Virtual_Genesis:
		return "Genesis"
	case *record.Virtual_IncomingRequest:
		return "IncomingRequest"
	case *record.Virtual_OutgoingRequest:
		return "OutgoingRequest"
	case *record.Virtual_Result:
		return "Result"
	case *record.Virtual_Code:
		return "Code"
	case *record.Virtual_Activate:
		return "Activate"
	case *record.Virtual_Amend:
		return amendPrettyPrint(r)
	case *record.Virtual_Deactivate:
		return "Deactivate"
	case *record.Virtual_PendingFilament:
		return "PendingFilament"
		// return r.PendingFilament
	case nil:
		return "nil"
	default:
		panic(fmt.Sprintf("%T virtual record unknown type", r))
	}
}

func amendPrettyPrint(virtualRecord *record.Virtual_Amend) string {
	rec := virtualRecord.Amend
	lines := []string{
		"request: " + rec.Request.String(),
		"memory: " + humanize.Bytes(uint64(len(rec.Memory))),
		"image: " + rec.Image.String(),
		"isPrototype: " + fmt.Sprint(rec.IsPrototype),
		"prevState: " + rec.PrevState.String(),
	}
	return strings.Join(lines, "\n")
}
