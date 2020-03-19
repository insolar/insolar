// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package common

import (
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

type LogicRunnerRPCStub interface {
	GetCode(rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	DeactivateObject(rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
}
