package message

import (
	"fmt"

	"github.com/insolar/insolar/core"
)

func ExtractTarget(msg core.Message) core.RecordRef {
	switch t := msg.(type) {
	case *GenesisRequest:
		return core.NewRefFromBase58(t.Name)
	case *CallConstructor:
		if t.SaveAs == Delegate {
			return t.ParentRef
		}
		return *core.GenRequest(t.PulseNum, MustSerializeBytes(t))
	case *CallMethod:
		return t.ObjectRef
	case *ExecutorResults:
		return t.RecordRef
	case *GetChildren:
		return t.Parent
	case *GetCode:
		return t.Code
	case *GetDelegate:
		return t.Head
	case *GetObject:
		return t.Head
	case *JetDrop:
		return t.Jet
	case *RegisterChild:
		return t.Parent
	case *SetBlob:
		return t.TargetRef
	case *SetRecord:
		return t.TargetRef
	case *UpdateObject:
		return t.Object
	case *ValidateCaseBind:
		return t.RecordRef
	case *ValidateRecord:
		return t.Object
	case *ValidationResults:
		return t.RecordRef
	case *SignedMessage:
		return ExtractTarget(t.Msg)
	default:
		panic(fmt.Sprintf("unknow message type - %v", t))
	}
}

func ExtractRole(msg core.Message) core.JetRole {
	switch t := msg.(type) {
	case *GenesisRequest:
		return core.RoleLightExecutor
	case *CallConstructor:
		return core.RoleVirtualExecutor
	case *CallMethod:
		return core.RoleVirtualExecutor
	case *ExecutorResults:
		return core.RoleVirtualExecutor
	case *GetChildren:
		return core.RoleLightExecutor
	case *GetCode:
		return core.RoleLightExecutor
	case *GetDelegate:
		return core.RoleLightExecutor
	case *GetObject:
		return core.RoleLightExecutor
	case *JetDrop:
		return core.RoleLightExecutor
	case *RegisterChild:
		return core.RoleLightExecutor
	case *SetBlob:
		return core.RoleLightExecutor
	case *SetRecord:
		return core.RoleLightExecutor
	case *UpdateObject:
		return core.RoleLightExecutor
	case *ValidateCaseBind:
		return core.RoleVirtualValidator
	case *ValidateRecord:
		return core.RoleLightExecutor
	case *ValidationResults:
		return core.RoleVirtualExecutor
	case *SignedMessage:
		return ExtractRole(t.Msg)
	default:
		panic(fmt.Sprintf("unknow message type - %v", t))
	}
}
