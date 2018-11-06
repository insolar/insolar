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
	default:
		panic("unknow message type")
	}
}
