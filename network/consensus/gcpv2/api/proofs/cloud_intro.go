// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proofs

import (
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

type NodeWelcomePackage struct {
	CloudIdentity      cryptkit.DigestHolder
	LastCloudStateHash cryptkit.DigestHolder
	JoinerSecret       cryptkit.DigestHolder // can be nil
}
