package message

import (
	"github.com/insolar/insolar/core"
)

func Extract(msg core.Message) core.RecordRef {
	switch m := msg.(type) {
	case *BootstrapRequest:
		return core.NewRefFromBase58(m.Name)
	case *CallConstructor:
		if m.SaveAs == Delegate {
			return m.ParentRef
		}
		return *core.GenRequest(m.PulseNum, MustSerializeBytes(m))
	case *CallMethod:
		return m.ObjectRef
	case *ExecutorResults:
		return m.RecordRef
	case *GetChildren:
		return m.Parent
	case *GetCode:
		return m.Code
	case *GetDelegate:
		return m.Head
	case *GetObject:
		return m.Head
	case *JetDrop:
		return m.Jet
	case *RegisterChild:
		return m.Parent
	default:
		panic("unknow message type")
	}
}

func ExtractRole(msg core.Message) core.JetRole {
	switch _ := msg.(type) {
	case *BootstrapRequest:
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
	default:
		panic("unknow message type")
	}
}
