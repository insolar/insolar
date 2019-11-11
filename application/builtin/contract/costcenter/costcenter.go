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

package costcenter

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type CostCenter struct {
	foundation.BaseContract
	FeeMember *insolar.Reference
	Fee       string
}

// New creates new CostCenter.
func New(feeMember *insolar.Reference, fee string) (*CostCenter, error) {
	return &CostCenter{
		FeeMember: feeMember,
		Fee:       fee,
	}, nil
}

// GetFeeMember gets fee member reference.
// ins:immutable
func (cc *CostCenter) GetFeeMember() (*insolar.Reference, error) {
	return cc.FeeMember, nil
}

// CalcFee calculates fee for amount. Returns fee.
// ins:immutable
func (cc *CostCenter) CalcFee(amountStr string) (string, error) {
	return cc.Fee, nil
}
