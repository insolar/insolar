package heavy

import "github.com/insolar/insolar/ledger/heavy/internal/handler"

func Components() []interface{} {
	return []interface{}{
		handler.New(),
	}
}
