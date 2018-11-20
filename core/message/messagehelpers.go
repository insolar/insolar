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
	case *HeavyRecords:
		return core.RecordRef{}
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
	case *HeavyRecords:
		return core.RoleHeavyExecutor
	case *Parcel:
		return ExtractRole(t.Msg)
	default:
		panic(fmt.Sprintf("unknow message type - %v", t))
	}
}

// ExtractAllowedSenderObjectAndRole extracts information from message
// verify senderrequired to 's "caller" for sender
// verification purpose. If nil then check of sender's role is not
// provided by the message bus
func ExtractAllowedSenderObjectAndRole(msg core.Message) (*core.RecordRef, core.JetRole) {
	switch t := msg.(type) {
	case *GenesisRequest:
		return nil, 0
	case *CallConstructor:
		return t.GetCaller(), core.RoleVirtualExecutor
	case *CallMethod:
		return t.GetCaller(), core.RoleVirtualExecutor
	case *ExecutorResults:
		return nil, 0
	case *GetChildren:
		return nil, 0
	case *GetCode:
		return nil, 0
	case *GetDelegate:
		return nil, 0
	case *GetObject:
		return nil, 0
	case *JetDrop:
		return nil, 0
	case *RegisterChild:
		return nil, 0
	case *SetBlob:
		return nil, 0
	case *SetRecord:
		return nil, 0
	case *UpdateObject:
		return nil, 0
	case *ValidateCaseBind:
		return &t.RecordRef, core.RoleVirtualExecutor
	case *ValidateRecord:
		return nil, 0
	case *ValidationResults:
		return &t.RecordRef, core.RoleVirtualValidator
	case *Parcel:
		return ExtractAllowedSenderObjectAndRole(t.Msg)
	default:
		panic(fmt.Sprintf("unknow message type - %v", t))
	}
}
