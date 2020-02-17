// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
