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

package record

import (
	"fmt"
)

func Wrap(record Record) Virtual {
	switch generic := record.(type) {
	case *Genesis:
		return Virtual{
			Union: &Virtual_Genesis{
				Genesis: generic,
			},
		}
	case *IncomingRequest:
		return Virtual{
			Union: &Virtual_IncomingRequest{
				IncomingRequest: generic,
			},
		}
	case *OutgoingRequest:
		return Virtual{
			Union: &Virtual_OutgoingRequest{
				OutgoingRequest: generic,
			},
		}
	case *Result:
		return Virtual{
			Union: &Virtual_Result{
				Result: generic,
			},
		}
	case *Code:
		return Virtual{
			Union: &Virtual_Code{
				Code: generic,
			},
		}
	case *Activate:
		return Virtual{
			Union: &Virtual_Activate{
				Activate: generic,
			},
		}
	case *Amend:
		return Virtual{
			Union: &Virtual_Amend{
				Amend: generic,
			},
		}
	case *Deactivate:
		return Virtual{
			Union: &Virtual_Deactivate{
				Deactivate: generic,
			},
		}
	case *PendingFilament:
		return Virtual{
			Union: &Virtual_PendingFilament{
				PendingFilament: generic,
			},
		}
	default:
		panic(fmt.Sprintf("%T record is not registered", generic))
	}
}

func Unwrap(v *Virtual) Record {
	if v == nil {
		return nil
	}
	switch r := v.Union.(type) {
	case *Virtual_Genesis:
		return r.Genesis
	case *Virtual_IncomingRequest:
		return r.IncomingRequest
	case *Virtual_OutgoingRequest:
		return r.OutgoingRequest
	case *Virtual_Result:
		return r.Result
	case *Virtual_Code:
		return r.Code
	case *Virtual_Activate:
		return r.Activate
	case *Virtual_Amend:
		return r.Amend
	case *Virtual_Deactivate:
		return r.Deactivate
	case *Virtual_PendingFilament:
		return r.PendingFilament
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("%T virtual record unknown type", r))
	}
}
