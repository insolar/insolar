// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package appfoundation

import "github.com/insolar/insolar/insolar"

type SagaAcceptInfo struct {
	Amount     string
	FromMember insolar.Reference
	Request    insolar.Reference
}
