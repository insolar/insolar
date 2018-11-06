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
	default:
		panic("unknow message type")
	}
}
