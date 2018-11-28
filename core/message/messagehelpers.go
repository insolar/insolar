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
		return *genRequest(t.PulseNum, MustSerializeBytes(t))
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
	case *HeavyPayload:
		return core.RecordRef{}
	case *GetObjectIndex:
		return t.Object
	case *Parcel:
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
	case
		*HeavyStartStop,
		*HeavyPayload,
		*GetObjectIndex:
		return core.RoleHeavyExecutor
	case *Parcel:
		return ExtractRole(t.Msg)
	default:
		panic(fmt.Sprintf("unknow message type - %v", t))
	}
}

// ExtractAllowedSenderObjectAndRole extracts information from message
// verify sender required to 's "caller" for sender
// verification purpose. If nil then check of sender's role is not
// provided by the message bus
func ExtractAllowedSenderObjectAndRole(msg core.Message) (*core.RecordRef, core.JetRole) {
	switch t := msg.(type) {
	case *GenesisRequest:
		return nil, 0
	case *CallConstructor:
		c := t.GetCaller()
		if c.IsEmpty() {
			return nil, 0
		}
		return c, core.RoleVirtualExecutor
	case *CallMethod:
		c := t.GetCaller()
		if c.IsEmpty() {
			return nil, 0
		}
		return c, core.RoleVirtualExecutor
	case *ExecutorResults:
		return nil, 0
	case *GetChildren:
		return &t.Parent, core.RoleVirtualExecutor
	case *GetCode:
		return &t.Code, core.RoleVirtualExecutor
	case *GetDelegate:
		return &t.Head, core.RoleVirtualExecutor
	case *GetObject:
		return &t.Head, core.RoleVirtualExecutor
	case *JetDrop:
		// This check is not needed, because JetDrop sender is explicitly checked in handler.
		return nil, core.RoleUndefined
	case *RegisterChild:
		return &t.Child, core.RoleVirtualExecutor
	case *SetBlob:
		return &t.TargetRef, core.RoleVirtualExecutor
	case *SetRecord:
		return &t.TargetRef, core.RoleVirtualExecutor
	case *UpdateObject:
		return &t.Object, core.RoleVirtualExecutor
	case *ValidateCaseBind:
		return &t.RecordRef, core.RoleVirtualExecutor
	case *ValidateRecord:
		return &t.Object, core.RoleVirtualExecutor
	case *ValidationResults:
		return &t.RecordRef, core.RoleVirtualValidator
	case *GetObjectIndex:
		return &t.Object, core.RoleLightExecutor
	case *Parcel:
		return ExtractAllowedSenderObjectAndRole(t.Msg)
	default:
		panic(fmt.Sprintf("unknown message type - %v", t))
	}
}
