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
	case *HotData:
		return t.Jet
	case *GetObjectIndex:
		return t.Object
	case *Parcel:
		return ExtractTarget(t.Msg)
	case *NodeSignPayload:
		return *t.NodeRef
	default:
		panic(fmt.Sprintf("[ ExtractTarget ] unknow message type - %s", t.Type().String()))
	}
}

func ExtractRole(msg core.Message) core.DynamicRole {
	switch t := msg.(type) {
	case
		*ExecutorResults,
		*CallMethod,
		*CallConstructor,
		*GenesisRequest,
		*ValidationResults:
		return core.DynamicRoleVirtualExecutor
	case
		*GetChildren,
		*GetCode,
		*GetDelegate,
		*GetObject,
		*JetDrop,
		*RegisterChild,
		*SetBlob,
		*SetRecord,
		*UpdateObject,
		*ValidateRecord,
		*HotData:
		return core.DynamicRoleLightExecutor
	case *ValidateCaseBind:
		return core.DynamicRoleVirtualValidator
	case
		*HeavyStartStop,
		*HeavyPayload,
		*GetObjectIndex:
		return core.DynamicRoleHeavyExecutor
	case *Parcel:
		return ExtractRole(t.Msg)
	case *NodeSignPayload:
		return core.DynamicRoleUndefined
	default:
		panic(fmt.Sprintf("[ ExtractRole ] unknow message type - %s", t.Type().String()))
	}
}

// ExtractAllowedSenderObjectAndRole extracts information from message
// verify sender required to 's "caller" for sender
// verification purpose. If nil then check of sender's role is not
// provided by the message bus
func ExtractAllowedSenderObjectAndRole(msg core.Message) (*core.RecordRef, core.DynamicRole) {
	switch t := msg.(type) {
	case *GenesisRequest:
		return nil, 0
	case *CallConstructor:
		c := t.GetCaller()
		if c.IsEmpty() {
			return nil, 0
		}
		return c, core.DynamicRoleVirtualExecutor
	case *CallMethod:
		c := t.GetCaller()
		if c.IsEmpty() {
			return nil, 0
		}
		return c, core.DynamicRoleVirtualExecutor
	case *ExecutorResults:
		return nil, 0
	case *GetChildren:
		return &t.Parent, core.DynamicRoleVirtualExecutor
	case *GetCode:
		return &t.Code, core.DynamicRoleVirtualExecutor
	case *GetDelegate:
		return &t.Head, core.DynamicRoleVirtualExecutor
	case *GetObject:
		return &t.Head, core.DynamicRoleVirtualExecutor
	case *JetDrop:
		// This check is not needed, because JetDrop sender is explicitly checked in handler.
		return nil, core.DynamicRoleUndefined
	case *RegisterChild:
		return &t.Child, core.DynamicRoleVirtualExecutor
	case *SetBlob:
		return &t.TargetRef, core.DynamicRoleVirtualExecutor
	case *SetRecord:
		return &t.TargetRef, core.DynamicRoleVirtualExecutor
	case *UpdateObject:
		return &t.Object, core.DynamicRoleVirtualExecutor
	case *ValidateCaseBind:
		return &t.RecordRef, core.DynamicRoleVirtualExecutor
	case *ValidateRecord:
		return &t.Object, core.DynamicRoleVirtualExecutor
	case *ValidationResults:
		return &t.RecordRef, core.DynamicRoleVirtualValidator
	case *GetObjectIndex:
		return &t.Object, core.DynamicRoleLightExecutor
	case *HotData:
		return &t.Jet, core.DynamicRoleLightExecutor
	case *Parcel:
		return ExtractAllowedSenderObjectAndRole(t.Msg)
	case *NodeSignPayload:
		return nil, core.DynamicRoleUndefined
	default:
		panic(fmt.Sprintf("unknown message type - %s", t.Type().String()))
	}
}
