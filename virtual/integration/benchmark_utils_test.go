package integration

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/testutils"
)

type ServerHelper struct {
	s *Server
}

func (h *ServerHelper) createUser(ctx context.Context) (*User, error) {
	user, err := NewUserWithKeys()
	if err != nil {
		return nil, errors.Errorf("failed to create new user: " + err.Error())
	}

	{
		callMethodReply, _, err := h.s.BasicAPICall(ctx, "member.create", nil, genesisrefs.ContractRootMember, user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to call member.create")
		}

		var result map[string]interface{}
		if cm, ok := callMethodReply.(*reply.CallMethod); !ok {
			return nil, errors.Wrapf(err, "unexpected type of return value %T", callMethodReply)
		} else if err := insolar.Deserialize(cm.Result, &result); err != nil {
			return nil, errors.Wrap(err, "failed to deserialize result")
		}

		r0, ok := result["Returns"]
		if ok && r0 != nil {
			if r1, ok := r0.([]interface{}); !ok {
				return nil, errors.Errorf("bad response: bad type of 'Returns' [%#v]", r0)
			} else if len(r1) != 2 {
				return nil, errors.Errorf("bad response: bad length of 'Returns' [%#v]", r0)
			} else if r2, ok := r1[0].(map[string]interface{}); !ok {
				return nil, errors.Errorf("bad response: bad type of first value [%#v]", r1)
			} else if r3, ok := r2["reference"]; !ok {
				return nil, errors.Errorf("bad response: absent reference field [%#v]", r2)
			} else if walletReferenceString, ok := r3.(string); !ok {
				return nil, errors.Errorf("bad response: reference field expected to be a string [%#v]", r3)
			} else if walletReference, err := insolar.NewReferenceFromString(walletReferenceString); err != nil {
				return nil, errors.Wrap(err, "bad response: got bad reference")
			} else {
				user.Reference = *walletReference
			}

			return user, nil
		}

		r0, ok = result["Error"]
		if ok && r0 != nil {
			return nil, errors.Errorf("%T: %#v", r0, r0)
		}

		panic("unreachable")
	}
}

func (h *ServerHelper) transferMoney(ctx context.Context, from User, to User, amount int64) (int64, error) {
	callParams := map[string]interface{}{
		"amount":            strconv.FormatInt(amount, 10),
		"toMemberReference": to.Reference.String(),
	}
	callMethodReply, _, err := h.s.BasicAPICall(ctx, "member.transfer", callParams, from.Reference, &from)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call member.transfer")
	}

	var result map[string]interface{}
	if cm, ok := callMethodReply.(*reply.CallMethod); !ok {
		return 0, errors.Wrapf(err, "unexpected type of return value %T", callMethodReply)
	} else if err := insolar.Deserialize(cm.Result, &result); err != nil {
		return 0, errors.Wrap(err, "failed to deserialize result")
	}

	r0, ok := result["Returns"]
	if ok && r0 != nil {
		if r1, ok := r0.([]interface{}); !ok {
			return 0, errors.Errorf("bad response: bad type of 'Returns' [%#v]", r0)
		} else if len(r1) != 2 {
			return 0, errors.Errorf("bad response: bad length of 'Returns' [%#v]", r0)
		} else if r2, ok := r1[0].(map[string]interface{}); !ok {
			return 0, errors.Errorf("bad response: bad type of first value [%#v]", r1)
		} else if r3, ok := r2["fee"]; !ok {
			return 0, errors.Errorf("bad response: absent fee field [%#v]", r2)
		} else if feeRaw, ok := r3.(string); !ok {
			return 0, errors.Errorf("bad response: Fee field expected to be a string [%#v]", r3)
		} else if fee, err := strconv.ParseInt(feeRaw, 10, 0); err != nil {
			return 0, errors.Wrapf(err, "failed to parse fee [%#v]", feeRaw)
		} else {
			return fee, nil
		}
	}

	r0, ok = result["Error"]
	if ok && r0 != nil {
		return 0, errors.Errorf("%T: %#v", r0, r0)
	}

	panic("unreachable")
}

func (h *ServerHelper) getBalance(ctx context.Context, user User) (int64, error) {
	callParams := map[string]interface{}{
		"reference": user.Reference.String(),
	}
	callMethodReply, _, err := h.s.BasicAPICall(ctx, "member.getBalance", callParams, user.Reference, &user)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call member.getBalance")
	}

	var result map[string]interface{}
	if cm, ok := callMethodReply.(*reply.CallMethod); !ok {
		return 0, errors.Wrapf(err, "unexpected type of return value %T", callMethodReply)
	} else if err := insolar.Deserialize(cm.Result, &result); err != nil {
		return 0, errors.Wrap(err, "failed to deserialize result")
	}

	r0, ok := result["Returns"]
	if ok && r0 != nil {
		if r1, ok := r0.([]interface{}); !ok {
			return 0, errors.Errorf("bad response: bad type of 'Returns' [%#v]", r0)
		} else if len(r1) != 2 {
			return 0, errors.Errorf("bad response: bad length of 'Returns' [%#v]", r0)
		} else if r2, ok := r1[0].(map[string]interface{}); !ok {
			return 0, errors.Errorf("bad response: bad type of first value [%#v]", r1)
		} else if r3, ok := r2["balance"]; !ok {
			return 0, errors.Errorf("bad response: absent balance field [%#v]", r2)
		} else if balanceRaw, ok := r3.(string); !ok {
			return 0, errors.Errorf("bad response: balance field expected to be a string [%#v]", r3)
		} else if balance, err := strconv.ParseInt(balanceRaw, 10, 0); err != nil {
			return 0, errors.Wrapf(err, "failed to parse balance [%#v]", balanceRaw)
		} else {
			return balance, nil
		}
	}

	r0, ok = result["Error"]
	if ok && r0 != nil {
		return 0, errors.Errorf("%T: %#v", r0, r0)
	}

	panic("unreachable")
}

func (h *ServerHelper) waitBalance(ctx context.Context, user User, balance int64) error {
	doneWaiting := false
	for i := 0; i < 100; i++ {
		balance, err := h.getBalance(ctx, user)
		if err != nil {
			return err
		}
		if balance == balance {
			doneWaiting = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !doneWaiting {
		return errors.New("failed to wait until balance match")
	}
	return nil
}

func NewTSAssert(tb testing.TB) *assert.Assertions {
	return assert.New(&testutils.SyncT{TB: tb})
}

func NewTSRequire(tb testing.TB) *require.Assertions {
	return require.New(&testutils.SyncT{TB: tb})
}
